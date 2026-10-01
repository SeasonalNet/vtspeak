package text

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestPositionModelArenaDerivesNativeGapsAndFinalGate(t *testing.T) {
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[:4], 12)
	binary.LittleEndian.PutUint32(workspace[4:8], 0)
	arena := markerArenaFixture(2)
	for index, interval := range [][2]uint32{{2, 4}, {8, 10}} {
		row := arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		binary.LittleEndian.PutUint32(row[:4], interval[0])
		binary.LittleEndian.PutUint32(row[4:8], interval[1])
	}
	input := Paul2013PositionStateProgram{RowCount: 99, Events: map[byte]Paul2013PositionEventValues{
		3: {Mode: 3, Boundaries: []int32{1, 3, 5, 9, 11}, Values: []int32{10, 20, 30, 40, 50}, Minimum: 0, Maximum: 100},
		4: {Mode: 4, Boundaries: []int32{1, 3, 5, 9, 11}, Values: []int32{1, 2, 3, 4, 5}},
	}}
	got, err := ApplyPaul2013PositionStateProgramFromModelArena(workspace, arena, 12, input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Program.Arrays[3], []int32{10, 30}) || !reflect.DeepEqual(got.Program.Arrays[4], []int32{3, 5}) || got.Program.Terminal.Value != 50 || got.Program.Terminal.NextBoundary != 5 {
		t.Fatalf("native gap program = %+v", got.Program)
	}
	if binary.LittleEndian.Uint32(got.Workspace[0x1223f4:]) != 5 || binary.LittleEndian.Uint32(got.Workspace[0x122404:]) != 5 {
		t.Fatal("gap cursors not written")
	}
	if input.RowCount != 99 || len(input.Events[3].Ranges) != 0 {
		t.Fatal("mutated caller program")
	}
	nonfinal, err := ApplyPaul2013PositionStateProgramFromModelArena(workspace, arena, 11, input)
	if err != nil {
		t.Fatal(err)
	}
	if nonfinal.Program.Terminal.Value != -1 || binary.LittleEndian.Uint32(nonfinal.Workspace[0x1223f4:]) != 3 {
		t.Fatal("nonfinal segment consumed terminal events")
	}
}

func TestPositionModelArenaOverlappingRecordsHaveEmptyGaps(t *testing.T) {
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[:4], 20)
	binary.LittleEndian.PutUint32(workspace[4:8], 0)
	arena := markerArenaFixture(2)
	for index := 0; index < 2; index++ {
		row := arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		binary.LittleEndian.PutUint32(row[:4], 5)
		binary.LittleEndian.PutUint32(row[4:8], 9)
	}
	got, err := ApplyPaul2013PositionStateProgramFromModelArena(workspace, arena, 10, Paul2013PositionStateProgram{Events: map[byte]Paul2013PositionEventValues{3: {Mode: 3, Boundaries: []int32{7}, Values: []int32{3}, Minimum: 0, Maximum: 10}}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Program.Arrays[3], []int32{-1, -1}) || got.Program.Events[3].NextBoundary != 0 {
		t.Fatalf("overlapping gap = %+v", got.Program)
	}
}
