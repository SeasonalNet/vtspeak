package text

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"vtspeak/engine/tree3"
)

func TestPositionBoundaryPipelineTransfersStatesAndMappedIntervals(t *testing.T) {
	state := finalizedBoundaryWorkspaceFixture()
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[:4], 15)
	binary.LittleEndian.PutUint32(workspace[4:8], 0)
	binary.LittleEndian.PutUint32(state.Finalized.StateArena[4:8], 15)
	stateBefore := append([]byte(nil), state.Records.Bytes...)
	workspaceBefore := append([]byte(nil), workspace...)
	start, end := make([]int32, 15), make([]int32, 15)
	for index := range start {
		start[index], end[index] = int32(index+100), int32(index+200)
	}
	program := Paul2013PositionStateProgram{Ranges: map[byte]Paul2013PositionStateRanges{7: {Mode: 7, Boundaries: []int32{-1, 20}, Values: []int32{7}, Minimum: 0, Maximum: 10}}, IndexTables: &Paul2013PositionIndexTables{Start: start, End: end, MaximumIndex: 15}}
	got, err := state.ApplyPositionStatesAndRunTokenBoundaries(workspace, program, 0x10000000, 0, &tree3.Catalog{Pronunciation: markerConstantTree(0)}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Position.Program.Arrays[7], []int32{7, 7}) || got.State.Records.Bytes[0x4770a] != 7 {
		t.Fatalf("state transfer = %+v", got.Position.Program)
	}
	for index := 0; index < 2; index++ {
		row := got.State.Records.Bytes[0x64c+index*0x3c0:]
		if row[0x2df] != 12 || int32(binary.LittleEndian.Uint32(row[:4])) != int32(index*10+100) || int32(binary.LittleEndian.Uint32(row[4:8])) != int32(index*10+202) {
			t.Fatalf("mapped row %d = interval %d/%d flag %d", index, binary.LittleEndian.Uint32(row[:4]), binary.LittleEndian.Uint32(row[4:8]), row[0x2df])
		}
	}
	if !bytes.Equal(stateBefore, state.Records.Bytes) || !bytes.Equal(workspaceBefore, workspace) {
		t.Fatal("mutated caller arena/workspace")
	}
}

func TestPositionMappingRejectsRecordSliceAndShortArena(t *testing.T) {
	tables := Paul2013PositionIndexTables{Start: []int32{1}, End: []int32{2}, MaximumIndex: 1}
	for _, arena := range [][]byte{make([]byte, 0x3c0), make([]byte, 0x64c+0x3c0-1)} {
		if _, err := MapPaul2013PositionStateRecordIndexes(arena, 1, tables); err == nil {
			t.Fatal("accepted truncated full arena")
		}
	}
	arena := make([]byte, 0x64c+0x3c0)
	got, err := MapPaul2013PositionStateRecordIndexes(arena, 1, tables)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []Paul2013PositionEventRange{{Minimum: 1, Maximum: 2}}) {
		t.Fatalf("full arena mapping = %+v", got)
	}
}

func TestPositionBoundaryPipelineRecomputesFlagsAndFinalMode(t *testing.T) {
	for _, test := range []struct {
		code, pitch int16
		value       int32
		mode, flag  byte
	}{{3, 0, -1, 6, 0}, {3, 1, -1, 7, 0}, {4, 0, -1, 7, 0}, {0, 0, -1, 5, 0}, {0, 0, 7, 7, 12}} {
		state := finalizedBoundaryWorkspaceFixture()
		last := state.Records.Bytes[0x64c+0x3c0:]
		binary.LittleEndian.PutUint16(last[0x3ac:], uint16(test.code))
		binary.LittleEndian.PutUint16(state.Finalized.StateArena[2:4], uint16(test.pitch))
		workspace := stateWorkspaceFixture()
		binary.LittleEndian.PutUint32(workspace[0x121a68:], 0xffffffff)
		program := Paul2013PositionStateProgram{Ranges: map[byte]Paul2013PositionStateRanges{7: {Mode: 7, Boundaries: []int32{-1, 20}, Values: []int32{test.value}, Minimum: 0, Maximum: 10}}}
		got, err := state.ApplyPositionStatesAndRunTokenBoundaries(workspace, program, 0x10000000, 0, &tree3.Catalog{Pronunciation: markerConstantTree(0)}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
		if err != nil {
			t.Fatal(err)
		}
		if got.State.Records.Bytes[0x4770a] != test.mode || got.State.Records.Bytes[0x64c+0x3c0+0x2df] != test.flag {
			t.Fatalf("code %d pitch %d value %d -> mode/flag %d/%d want %d/%d", test.code, test.pitch, test.value, got.State.Records.Bytes[0x4770a], got.State.Records.Bytes[0x64c+0x3c0+0x2df], test.mode, test.flag)
		}
	}
}
