package text

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestModelPhoneGroupArenaWritesAndPreservesFields(t *testing.T) {
	arena := boundaryTestArena(3)
	for index, code := range []byte{0x13, 0x01, 0x01} {
		row := arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		row[0x95] = 1
		row[0x2e8] = code
		row[0x329] = '0'
	}
	before := append([]byte(nil), arena...)
	got, err := PopulatePaul2013ModelPhoneGroupArena(arena, 0x10000000)
	if err != nil {
		t.Fatal(err)
	}
	if got.Count != 2 || binary.LittleEndian.Uint16(got.Arena[4:6]) != 2 {
		t.Fatalf("group count = %d", got.Count)
	}
	for index, want := range []uint32{0x10017d4c, 0x10017d4c, 0x10017d6a} {
		row := got.Arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		if pointer := binary.LittleEndian.Uint32(row[8:12]); pointer != want {
			t.Fatalf("record %d pointer = %#x, want %#x", index, pointer, want)
		}
		if row[0x36a] != []byte{1, 2, 2}[index] || row[0x36b] != 0xa5 {
			t.Fatalf("record %d labels = %x", index, row[0x36a:0x36c])
		}
	}
	for index := 0; index < 2; index++ {
		row := got.Arena[paul2013ModelPhoneGroupArenaBase+index*0x1e:]
		if row[1] != 0 || row[2] != 1 || row[0x1c] != 0 || row[0x1d] != 1 {
			t.Fatalf("group %d = %x", index, row[:0x1e])
		}
		for field := 3; field < 0x1c; field++ {
			if row[field] != 0xa5 {
				t.Fatalf("group %d overwrote field %#x", index, field)
			}
		}
	}
	if !bytes.Equal(before, arena) {
		t.Fatal("mutated input")
	}
	prepared, err := PreparePaul2013TokenBoundaryArena(arena, 0x10000000)
	if err != nil {
		t.Fatal(err)
	}
	if len(prepared.Groups) != 1 || prepared.Groups[0].RecordCount != 3 || prepared.Groups[0].OpaqueRowSpan != 2 {
		t.Fatalf("pre-tree groups = %+v", prepared.Groups)
	}
	finished, err := FinishPaul2013TokenBoundaryArena(prepared.PhoneGroups.Arena, []int32{-1, -1, -1}, []int32{4, 4, 4}, []int32{-1, -1})
	if err != nil {
		t.Fatal(err)
	}
	if finished.Groups[0].OpaqueRowSpan != 2 || finished.State.Markers[2].Marker != 'Z' {
		t.Fatalf("post-tree handoff = %+v/%+v", finished.Groups, finished.State.Markers)
	}
}

func TestModelPhoneGroupArenaUncountedTailAndEmptyRecord(t *testing.T) {
	arena := boundaryTestArena(2)
	row := arena[paul2013TokenMarkerArenaRecords:]
	row[0x95] = 1
	row[0x2e8] = 0x13
	row[0x329] = '0'
	got, err := PopulatePaul2013ModelPhoneGroupArena(arena, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got.Count != 0 || got.Arena[paul2013ModelPhoneGroupArenaBase+1] != 4 || got.Arena[paul2013ModelPhoneGroupArenaBase+0x1d] != 1 {
		t.Fatalf("fallback = count %d row %x", got.Count, got.Arena[paul2013ModelPhoneGroupArenaBase:paul2013ModelPhoneGroupArenaBase+0x1e])
	}
	if got.Arena[paul2013TokenMarkerArenaRecords+paul2013TokenMarkerArenaStride+0x94] != 0 {
		t.Fatal("empty record group count not reset")
	}
	if _, err := PopulatePaul2013ModelPhoneGroupArena(arena, 0xffff0000); err == nil {
		t.Fatal("accepted overflowing pointer")
	}
}
