package text

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"testing"

	"vtspeak/engine/tree3"
)

func descriptorWorkspaceFixture() []byte {
	workspace := stateWorkspaceFixture()
	clear(workspace[0x1223c0:0x122440])
	return workspace
}

func setStateDescriptor(workspace []byte, mode byte, count, cursor int32, boundaries, values uint32) {
	row := workspace[0x1223c0+int(mode)*0x10:]
	for index, value := range []uint32{uint32(count), uint32(cursor), boundaries, values} {
		binary.LittleEndian.PutUint32(row[index*4:], value)
	}
}

func TestNativeStateDescriptorsCompletePipeline(t *testing.T) {
	state := finalizedBoundaryWorkspaceFixture()
	workspace := descriptorWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[:4], 15)
	binary.LittleEndian.PutUint32(workspace[4:8], 0)
	binary.LittleEndian.PutUint32(state.Finalized.StateArena[4:8], 15)
	storage := map[uint32][]int32{}
	for _, mode := range []byte{0, 1, 2, 7} {
		bp, vp := uint32(0x1000+int(mode)*0x100), uint32(0x2000+int(mode)*0x100)
		setStateDescriptor(workspace, mode, 2, 0, bp, vp)
		storage[bp], storage[vp] = []int32{-1, 20}, []int32{1000}
	}
	setStateDescriptor(workspace, 3, 3, 1, 0x3000, 0x4000)
	storage[0x3000], storage[0x4000] = []int32{0, 5, 14}, []int32{999, 70000, 7}
	setStateDescriptor(workspace, 4, 2, 0, 0x5000, 0x6000)
	storage[0x5000], storage[0x6000] = []int32{5, 14}, []int32{8, 2}
	before := append([]byte(nil), workspace...)
	resolve := func(address uint32, count int) ([]int32, error) {
		values, ok := storage[address]
		if !ok || len(values) < count {
			return nil, fmt.Errorf("missing cells at %#x", address)
		}
		return values[:count], nil
	}
	got, err := state.RunPositionDescriptorBoundaries(workspace, resolve, nil, 0x10000000, 0, &tree3.Catalog{Pronunciation: markerConstantTree(0)}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
	if err != nil {
		t.Fatal(err)
	}
	for mode, want := range map[int][]int32{0: {200, 200}, 1: {400, 400}, 2: {500, 500}, 7: {9, 9}, 3: {-1, 65535}, 4: {8, 2}} {
		if !reflect.DeepEqual(got.Position.Program.Arrays[mode], want) {
			t.Fatalf("selector %d = %v want %v", mode, got.Position.Program.Arrays[mode], want)
		}
	}
	if !got.Position.Program.HasTerminalPass || got.Position.Program.Terminal.Value != 7 || binary.LittleEndian.Uint32(got.Boundary.Workspace[0x1223f4:]) != 3 || got.State.Records.Bytes[0x4770a] != 7 {
		t.Fatalf("terminal/mode handoff = %+v", got.Position.Program)
	}
	if !bytes.Equal(before, workspace) || storage[0x4000][1] != 70000 {
		t.Fatal("mutated native input storage")
	}
}

func TestNativeStateDescriptorGatesAndOwnedCopies(t *testing.T) {
	workspace := descriptorWorkspaceFixture()
	for _, mode := range []byte{0, 1, 2, 7, 3, 4} {
		setStateDescriptor(workspace, mode, -1, -9, 1, 2)
	}
	setStateDescriptor(workspace, 0, 2, 0, 0, 2)
	setStateDescriptor(workspace, 1, 2, 0, 1, 0)
	setStateDescriptor(workspace, 3, 2, 0, 0, 0)
	setStateDescriptor(workspace, 5, 2, 0, 1, 2) // Unvisited native selector.
	got, err := ReadPaul2013PositionStateDescriptorProgram(workspace, nil, nil)
	if err != nil || len(got.Ranges) != 0 || len(got.Events) != 0 {
		t.Fatalf("inactive descriptors = %+v, %v", got, err)
	}
	setStateDescriptor(workspace, 2, 1, 0, 1, 2)
	cells := []int32{10}
	calls := 0
	got, err = ReadPaul2013PositionStateDescriptorProgram(workspace, func(address uint32, count int) ([]int32, error) {
		calls++
		if address != 1 || count != 1 {
			t.Fatalf("unnecessary read %#x/%d", address, count)
		}
		return cells, nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	cells[0] = 42
	if calls != 1 || got.Ranges[2].Boundaries[0] != 10 || len(got.Ranges[2].Values) != 0 {
		t.Fatal("singleton range read/copy differs from native loop")
	}
}

func TestNativeStateDescriptorErrors(t *testing.T) {
	if _, err := ReadPaul2013PositionStateDescriptorProgram(make([]byte, 0x12243f), nil, nil); err == nil {
		t.Fatal("accepted truncated descriptors")
	}
	for _, test := range []struct {
		mode    byte
		bp, vp  uint32
		resolve Paul2013Int32PointerResolver
	}{{0, 1, 2, nil}, {3, 1, 0, func(uint32, int) ([]int32, error) { return []int32{0}, nil }}, {4, 0xfffffffe, 2, func(uint32, int) ([]int32, error) { t.Fatal("resolved overflowing pointer"); return nil, nil }}, {7, 1, 2, func(uint32, int) ([]int32, error) { return nil, nil }}, {2, 1, 2, func(uint32, int) ([]int32, error) { return nil, fmt.Errorf("resolver failure") }}} {
		workspace := descriptorWorkspaceFixture()
		setStateDescriptor(workspace, test.mode, 1, 0, test.bp, test.vp)
		if _, err := ReadPaul2013PositionStateDescriptorProgram(workspace, test.resolve, nil); err == nil {
			t.Fatalf("accepted invalid selector %d pointers", test.mode)
		}
	}
}

func TestNativePositionClampTablesMatchLocalDLL(t *testing.T) {
	file, err := pe.Open("../../binary/vt_pau.dll")
	if os.IsNotExist(err) {
		t.Skip("local proprietary DLL absent")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	header, ok := file.OptionalHeader.(*pe.OptionalHeader32)
	if !ok || header.ImageBase != 0x10000000 {
		t.Fatal("unexpected PE image base")
	}
	for _, table := range []struct {
		va     uint32
		values [8]int32
	}{{0x1007d66c, paul2013PositionMaximum}, {0x1007d690, paul2013PositionMinimum}} {
		rva := table.va - header.ImageBase
		found := false
		for _, section := range file.Sections {
			if rva < section.VirtualAddress || uint64(rva-section.VirtualAddress)+32 > uint64(section.Size) {
				continue
			}
			raw, err := section.Data()
			if err != nil {
				t.Fatal(err)
			}
			start := int(rva - section.VirtualAddress)
			for index, value := range table.values {
				if got := int32(binary.LittleEndian.Uint32(raw[start+index*4:])); got != value {
					t.Fatalf("VA %#x[%d] = %d want %d", table.va, index, got, value)
				}
			}
			found = true
			break
		}
		if !found {
			t.Fatalf("unmapped clamp table %#x", table.va)
		}
	}
}
