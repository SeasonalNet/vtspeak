package text

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func boundaryTestArena(count int) []byte {
	arena := bytes.Repeat([]byte{0xa5}, paul2013ModelStateFinalModeOffset+1)
	binary.LittleEndian.PutUint16(arena[2:4], uint16(count))
	for index := 0; index < count; index++ {
		row := arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		binary.LittleEndian.PutUint16(row[0x3ac:0x3ae], 0)
		row[0x3b0], row[0x94], row[0x95] = 0, 1, 0
	}
	return arena
}

func TestTokenBoundaryArenaNativeOrder(t *testing.T) {
	arena := boundaryTestArena(3)
	arena[paul2013ModelStateFinalModeOffset] = 6
	row := func(index int) []byte {
		return arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
	}
	row(1)[0x3b0] = 8
	before := append([]byte(nil), arena...)
	initialized, err := InitializePaul2013TokenBoundaryArena(arena)
	if err != nil {
		t.Fatal(err)
	}
	if got := initialized[paul2013TokenMarkerArenaRecords+2*paul2013TokenMarkerArenaStride+0x3bd]; got != '^' {
		t.Fatalf("final mode marker = %q", got)
	}
	initializedBefore := append([]byte(nil), initialized...)
	got, err := FinishPaul2013TokenBoundaryArena(initialized, []int32{-1, -1, -1}, []int32{3, 1, 4}, []int32{-2, 7})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.State.Values, []int32{100, 7, -1}) || !reflect.DeepEqual(got.State.States, []int32{2, 1, 4}) {
		t.Fatalf("state = %+v", got.State)
	}
	want := []Paul2013TokenBoundaryMarker{{Marker: '[', PreviousStateFlag: true}, {Marker: '\\', PreviousStateFlag: true}, {Marker: '^'}}
	if !reflect.DeepEqual(got.State.Markers, want) {
		t.Fatalf("markers = %+v, want %+v", got.State.Markers, want)
	}
	for index, marker := range want {
		row := got.Arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		flag := byte(0)
		if marker.PreviousStateFlag {
			flag = 1
		}
		if row[0x3bd] != marker.Marker || row[0x3bc] != flag {
			t.Fatalf("row %d marker/flag = %x", index, row[0x3bc:0x3be])
		}
	}
	if len(got.Groups) != 2 || got.Groups[1].StartRecordIndex != 1 || got.Groups[1].RecordCount != 2 {
		t.Fatalf("groups = %+v", got.Groups)
	}
	if !bytes.Equal(arena, before) || !bytes.Equal(initialized, initializedBefore) {
		t.Fatal("mutated caller arena")
	}
}

func TestTokenBoundaryArenaFullRecordCapacityAndRescan(t *testing.T) {
	arena := boundaryTestArena(100)
	for index := 0; index < 100; index++ {
		arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride+0x95] = 250
	}
	initialized, err := InitializePaul2013TokenBoundaryArena(arena)
	if err != nil {
		t.Fatal(err)
	}
	states, values, outputs := make([]int32, 100), make([]int32, 100), make([]int32, 99)
	for index := range states {
		states[index], values[index] = 4, -1
	}
	for index := range outputs {
		outputs[index] = -1
	}
	got, err := FinishPaul2013TokenBoundaryArena(initialized, values, states, outputs)
	if err != nil {
		t.Fatal(err)
	}
	for index := range got.State.Markers {
		want := byte(']')
		if index%2 == 1 {
			want = '['
		}
		if index == 99 {
			want = 'Z'
		}
		if got.State.Markers[index].Marker != want {
			t.Fatalf("record %d = %q, want %q", index, got.State.Markers[index].Marker, want)
		}
	}
	if len(got.Groups) != 50 {
		t.Fatalf("groups = %d", len(got.Groups))
	}
	markers := bytes.Repeat([]byte{']'}, 100)
	limited, err := ApplyPaul2013TokenBoundaryDurationLimitFromModelStateArena(arena, markers)
	if err != nil || limited[1] != '[' || limited[97] != '[' {
		t.Fatalf("arena limiter = %v, %v", limited, err)
	}
}

func TestTokenBoundaryArenaModeAndFinalState(t *testing.T) {
	for _, mode := range []byte{6, 7, 0} {
		for _, state := range []int32{-1, 0, 1, 2, 3, 4} {
			arena := boundaryTestArena(1)
			arena[paul2013ModelStateFinalModeOffset] = mode
			initialized, err := InitializePaul2013TokenBoundaryArena(arena)
			if err != nil {
				t.Fatal(err)
			}
			got, err := FinishPaul2013TokenBoundaryArena(initialized, []int32{-1}, []int32{state}, nil)
			if err != nil {
				t.Fatal(err)
			}
			want := byte('Z')
			if mode == 6 {
				want = '^'
			}
			if state >= 0 && state < 3 {
				want = '['
			}
			if got.State.Markers[0].Marker != want || got.State.Markers[0].PreviousStateFlag {
				t.Fatalf("mode %d state %d = %+v", mode, state, got.State.Markers)
			}
		}
	}
}

func TestTokenBoundaryArenaRejectsInvalidInputs(t *testing.T) {
	if _, err := InitializePaul2013TokenBoundaryArena(boundaryTestArena(0)); err == nil {
		t.Fatal("accepted zero records")
	}
	if _, err := InitializePaul2013TokenBoundaryArena(boundaryTestArena(101)); err == nil {
		t.Fatal("accepted excess records")
	}
	if _, err := InitializePaul2013TokenBoundaryArena(boundaryTestArena(1)[:0x1000]); err == nil {
		t.Fatal("accepted missing mode")
	}
	if _, err := FinishPaul2013TokenBoundaryArena(boundaryTestArena(2), []int32{-1, -1}, []int32{4, 4}, nil); err == nil {
		t.Fatal("accepted missing parser outputs")
	}
}
