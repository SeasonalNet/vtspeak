package text

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestApplyPaul2013PositionStateProgramFromParserState(t *testing.T) {
	state := make([]byte, 0x14+2*0x94)
	binary.LittleEndian.PutUint16(state[:2], 2)
	binary.LittleEndian.PutUint32(state[0x14:], uint32(1))
	binary.LittleEndian.PutUint32(state[0x18:], uint32(4))
	binary.LittleEndian.PutUint32(state[0x14+0x94:], uint32(5))
	binary.LittleEndian.PutUint32(state[0x18+0x94:], uint32(7))

	program := Paul2013PositionStateProgram{
		RowCount: 2,
		Initial:  [3]int32{100, 100, 200},
		Events: map[byte]Paul2013PositionEventValues{
			3: {
				Mode: 3, Boundaries: []int32{11, 15}, Values: []int32{2, 3},
				Initial: []int32{-1, -1}, StartBoundary: 0, OverflowStart: 0,
				Minimum: 0, Maximum: 500,
			},
			4: {
				Mode: 4, Boundaries: []int32{11, 15}, Values: []int32{7, 9},
				Initial: []int32{-1, -1}, StartBoundary: 0, OverflowStart: 0,
			},
		},
	}
	got, err := ApplyPaul2013PositionStateProgramFromParserState(program, 10, state)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int32{2, 3}; !reflect.DeepEqual(got.Arrays[3], want) {
		t.Fatalf("mode-3 values = %v, want %v", got.Arrays[3], want)
	}
	if want := []int32{7, 9}; !reflect.DeepEqual(got.Arrays[4], want) {
		t.Fatalf("mode-4 values = %v, want %v", got.Arrays[4], want)
	}
	if len(program.Events[3].Ranges) != 0 {
		t.Fatal("adapter mutated the caller's event map")
	}
}

func TestApplyPaul2013PositionStateProgramFromParserStateRejectsRowCountMismatch(t *testing.T) {
	state := make([]byte, 2)
	binary.LittleEndian.PutUint16(state, 1)
	_, err := ApplyPaul2013PositionStateProgramFromParserState(
		Paul2013PositionStateProgram{RowCount: 0}, 0, state,
	)
	if err == nil {
		t.Fatal("expected parser/program row-count mismatch")
	}
}

func TestApplyPaul2013PositionStateProgramsForOrdinarySource(t *testing.T) {
	programs := []Paul2013PositionStateProgram{
		{
			RowCount: 2,
			Events: map[byte]Paul2013PositionEventValues{
				3: {
					Mode: 3, Boundaries: []int32{100, 103}, Values: []int32{4, 7},
					Initial: []int32{-1, -1}, Minimum: 0, Maximum: 100,
				},
			},
		},
		{
			RowCount: 1,
			Events: map[byte]Paul2013PositionEventValues{
				3: {
					Mode: 3, Boundaries: []int32{110}, Values: []int32{9},
					Initial: []int32{-1}, Minimum: 0, Maximum: 100,
				},
			},
		},
	}
	got, err := ApplyPaul2013PositionStateProgramsForOrdinarySource(
		"Hi there. Yes", 100, programs,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d segment results, want 2", len(got))
	}
	if want := []int32{4, 7}; !reflect.DeepEqual(got[0].Arrays[3], want) {
		t.Fatalf("first segment event values = %v, want %v", got[0].Arrays[3], want)
	}
	if want := []int32{9}; !reflect.DeepEqual(got[1].Arrays[3], want) {
		t.Fatalf("second segment event values = %v, want %v", got[1].Arrays[3], want)
	}
}

func TestApplyPaul2013PositionStateProgramsForOrdinarySourceRejectsProgramCount(t *testing.T) {
	if _, err := ApplyPaul2013PositionStateProgramsForOrdinarySource("Hi. Yes", 0, nil); err == nil {
		t.Fatal("expected one program per ordinary-source segment")
	}
}
