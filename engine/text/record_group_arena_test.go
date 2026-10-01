package text

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestRecordGroupArenaHeaderAndPointerWrites(t *testing.T) {
	arena := markerArenaFixture(4)
	arena[paul2013TokenMarkerArenaRecords+paul2013TokenMarkerArenaStride+0x3bd] = '['
	before := append([]byte(nil), arena...)
	got, groups, err := PopulatePaul2013RecordGroupDescriptorArena(arena, 0x10000000)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 || binary.LittleEndian.Uint16(got[0:2]) != 2 {
		t.Fatalf("group count = %+v", groups)
	}
	for index, want := range []struct {
		start, count, span uint16
		record, rows       uint32
	}{{0, 2, 3, 0x1000064c, 0x10017d4c}, {2, 2, 7, 0x10000dcc, 0x10017da6}} {
		base := index * 0x10
		if binary.LittleEndian.Uint16(got[base+0xc:]) != want.start || binary.LittleEndian.Uint16(got[base+0xe:]) != want.count || binary.LittleEndian.Uint16(got[base+0x10:]) != want.span || binary.LittleEndian.Uint32(got[base+0x14:]) != want.record || binary.LittleEndian.Uint32(got[base+0x18:]) != want.rows {
			t.Fatalf("descriptor %d = %x", index, got[base+0xc:base+0x1c])
		}
		if !bytes.Equal(got[base+0x12:base+0x14], []byte{0xa5, 0xa5}) {
			t.Fatal("overwrote descriptor padding")
		}
	}
	if !bytes.Equal(got[2:12], before[2:12]) || !bytes.Equal(got[0x2c:], before[0x2c:]) || !bytes.Equal(arena, before) {
		t.Fatal("overwrote unrelated bytes or caller input")
	}
	if _, _, err := PopulatePaul2013RecordGroupDescriptorArena(arena, 0xffff0000); err == nil {
		t.Fatal("accepted pointer overflow")
	}
}

func TestRecordGroupArenaEmptyOnlyClearsCount(t *testing.T) {
	arena := boundaryTestArena(0)
	got, groups, err := PopulatePaul2013RecordGroupDescriptorArena(arena, 0xffffffff)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 || binary.LittleEndian.Uint16(got[:2]) != 0 || !bytes.Equal(got[2:], arena[2:]) {
		t.Fatal("zero-record descriptor reset changed other fields")
	}
}

func TestBoundaryPipelineRewritesDescriptorsAfterDispatch(t *testing.T) {
	got, err := RunPaul2013TokenBoundaryPipeline(markerArenaFixture(3), 0x10000000, markerConstantTree(0), func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil }, []int32{-1, -1, -1}, []int32{2, 4, 4}, []int32{-1, -1})
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint16(got.Prepared.PhoneGroups.Arena[:2]) != 1 || binary.LittleEndian.Uint16(got.Final.Arena[:2]) != 2 {
		t.Fatal("final descriptors were not rebuilt")
	}
	if binary.LittleEndian.Uint16(got.Final.Arena[0xe:]) != 1 || binary.LittleEndian.Uint16(got.Final.Arena[0x1e:]) != 2 || binary.LittleEndian.Uint32(got.Final.Arena[0x24:]) != 0x10000a0c || binary.LittleEndian.Uint32(got.Final.Arena[0x28:]) != 0x10017d6a {
		t.Fatalf("final descriptors = %x", got.Final.Arena[:0x2c])
	}
}
