package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"vtspeak/engine/tree3"
)

func positionIndexWorkspaceFixture() ([]byte, []byte, Paul2013FinalizedModelState) {
	workspace, shared := descriptorWorkspaceFixture(), make([]byte, 0x20425)
	binary.LittleEndian.PutUint32(workspace[:4], 15)
	binary.LittleEndian.PutUint32(workspace[4:8], 0)
	for offset, pointer := range map[int]uint32{0x47790: 0x1000, 0x47794: 0x2000, 0x47798: 0x3000} {
		binary.LittleEndian.PutUint32(workspace[offset:], pointer)
	}
	state := finalizedBoundaryWorkspaceFixture()
	binary.LittleEndian.PutUint32(state.Finalized.StateArena[4:8], 15)
	return workspace, shared, state
}

func TestNativeIndexWorkspaceOrdinaryAndRemappedPipeline(t *testing.T) {
	for _, flag := range []byte{0, 1, 2} {
		workspace, shared, state := positionIndexWorkspaceFixture()
		shared[0x20424] = flag
		before := append([]byte(nil), workspace...)
		calls := make(map[uint32]int)
		resolve := func(address uint32, count int) ([]int32, error) {
			calls[address] = count
			cells := make([]int32, count)
			for index := range cells {
				cells[index] = -100 // Unused final indexes must not reject remap.
			}
			switch address {
			case 0x2000:
				cells[0], cells[10] = 1, 2
			case 0x3000:
				cells[2], cells[12] = 3, 4
			case 0x1000:
				for index := range cells {
					cells[index] = int32(100 + index)
				}
			default:
				return nil, fmt.Errorf("unknown pointer %#x", address)
			}
			return cells, nil
		}
		got, err := state.RunNativePositionWorkspaceBoundaries(workspace, shared, resolve, 0x10000000, 0, &tree3.Catalog{Pronunciation: markerConstantTree(0)}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
		if err != nil {
			t.Fatal(err)
		}
		want := []Paul2013PositionEventRange{{Minimum: 1, Maximum: 3}, {Minimum: 2, Maximum: 4}}
		if flag == 1 {
			want = []Paul2013PositionEventRange{{Minimum: 101, Maximum: 103}, {Minimum: 102, Maximum: 104}}
			if calls[0x1000] != 5 {
				t.Fatalf("final prefix = %v", calls)
			}
		} else if _, ok := calls[0x1000]; ok {
			t.Fatal("ordinary branch resolved final table")
		}
		if calls[0x2000] != 11 || calls[0x3000] != 13 || !reflect.DeepEqual(got.Position.Program.MappedIntervals, want) {
			t.Fatalf("flag %d intervals %v calls %v", flag, got.Position.Program.MappedIntervals, calls)
		}
		if !bytes.Equal(workspace, before) {
			t.Fatal("mutated workspace")
		}
	}
}

func TestNativeIndexWorkspaceClampingAndOwnership(t *testing.T) {
	workspace, shared, state := positionIndexWorkspaceFixture()
	row := state.Records.Bytes[0x64c:]
	binary.LittleEndian.PutUint32(row[:4], 0xffffffff)
	binary.LittleEndian.PutUint32(row[4:8], 999)
	storage := make([]int32, 15)
	for index := range storage {
		storage[index] = int32(index + 50)
	}
	tables, err := ReadPaul2013PositionIndexWorkspaceTables(workspace, shared, state.Records.Bytes, func(_ uint32, count int) ([]int32, error) { return storage[:count], nil })
	if err != nil {
		t.Fatal(err)
	}
	storage[0] = 999
	got, err := MapPaul2013PositionStateRecordIndexes(state.Records.Bytes, 2, tables)
	if err != nil || got[0] != (Paul2013PositionEventRange{Minimum: 50, Maximum: 64}) {
		t.Fatalf("clamp = %v, %v", got, err)
	}
}

func TestNativeIndexWorkspaceRejectsInvalidMemory(t *testing.T) {
	for _, kind := range []string{"short workspace", "short shared", "short arena", "zero length", "negative remap", "huge remap", "negative final", "huge final", "resolver error", "short cells", "null pointer"} {
		t.Run(kind, func(t *testing.T) {
			workspace, shared, state := positionIndexWorkspaceFixture()
			resolve := Paul2013Int32PointerResolver(func(_ uint32, count int) ([]int32, error) { return make([]int32, count), nil })
			switch kind {
			case "short workspace":
				workspace = workspace[:0x4779b]
			case "short shared":
				shared = shared[:0x20424]
			case "short arena":
				state.Records.Bytes = state.Records.Bytes[:2]
			case "zero length":
				binary.LittleEndian.PutUint32(workspace[:4], 0)
			case "negative remap", "huge remap":
				shared[0x20424] = 1
				value := uint32(0xffffffff)
				if kind == "huge remap" {
					value = 0x7fffffff
				}
				binary.LittleEndian.PutUint32(state.Records.Bytes[0x64c:], value)
			case "negative final", "huge final":
				shared[0x20424] = 1
				resolve = func(_ uint32, count int) ([]int32, error) {
					cells := make([]int32, count)
					cells[0] = -1
					if kind == "huge final" {
						cells[0] = 0x7fffffff
					}
					return cells, nil
				}
			case "resolver error":
				resolve = func(uint32, int) ([]int32, error) { return nil, fmt.Errorf("unmapped") }
			case "short cells":
				resolve = func(uint32, int) ([]int32, error) { return nil, nil }
			case "null pointer":
				binary.LittleEndian.PutUint32(workspace[0x47794:], 0)
			}
			if _, err := ReadPaul2013PositionIndexWorkspaceTables(workspace, shared, state.Records.Bytes, resolve); err == nil {
				t.Fatal("accepted invalid mapping memory")
			}
		})
	}
}
