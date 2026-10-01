package synthesis

import (
	"encoding/binary"
	"reflect"
	"testing"

	"vtspeak/engine/text"
)

func TestApplyPaul2013PositionStateProgramWithControls(t *testing.T) {
	input := Paul2013ControlOverrides{
		Defaults: Controls{Pitch: 120, Speed: 90, Volume: 180},
		Values:   Controls{Pitch: 240, Speed: 450, Volume: -1},
		Pitch:    true,
		Speed:    true,
		Volume:   false,
	}
	program := text.Paul2013PositionStateProgram{
		RowCount: 2,
		Ranges: map[byte]text.Paul2013PositionStateRanges{
			0: {Mode: 0, RowKeys: []int32{0, 1}, Boundaries: []int32{0, 2}, Values: []int32{-1}, Minimum: 0, Maximum: 250},
		},
	}
	got, controls, err := ApplyPaul2013PositionStateProgramWithControls(program, input)
	if err != nil {
		t.Fatal(err)
	}
	wantControls := Controls{Pitch: 200, Speed: 400, Volume: 180}
	if controls != wantControls {
		t.Fatalf("effective controls = %+v, want %+v", controls, wantControls)
	}
	wantInitial := []int32{200, 200}
	if !reflect.DeepEqual(got.Arrays[0], wantInitial) {
		t.Fatalf("pitch state array = %v, want %v", got.Arrays[0], wantInitial)
	}
	if !reflect.DeepEqual(got.Arrays[1], []int32{400, 400}) || !reflect.DeepEqual(got.Arrays[2], []int32{180, 180}) {
		t.Fatalf("speed/volume arrays = %v/%v", got.Arrays[1], got.Arrays[2])
	}
	if got.Arrays[3][0] != -1 || got.Arrays[4][0] != -1 || got.Arrays[7][0] != -1 {
		t.Fatalf("sentinel arrays were changed: 3=%v 4=%v 7=%v", got.Arrays[3], got.Arrays[4], got.Arrays[7])
	}
}

func TestApplyPaul2013PositionStateProgramFromParserStateWithControls(t *testing.T) {
	parserState := make([]byte, 0x14+0x94)
	binary.LittleEndian.PutUint16(parserState[:2], 1)
	binary.LittleEndian.PutUint32(parserState[0x14:], 0)
	binary.LittleEndian.PutUint32(parserState[0x18:], 1)
	program := text.Paul2013PositionStateProgram{
		RowCount: 1,
		Events: map[byte]text.Paul2013PositionEventValues{
			3: {
				Mode: 3, Boundaries: []int32{4}, Values: []int32{6},
				Initial: []int32{-1}, Minimum: 0, Maximum: 100,
			},
		},
	}
	result, effective, err := ApplyPaul2013PositionStateProgramFromParserStateWithControls(
		program, 4, parserState,
		Paul2013ControlOverrides{
			Defaults: Controls{Pitch: 90, Speed: 110, Volume: 150},
			Values:   Controls{Pitch: 250, Speed: 25, Volume: 700},
			Pitch:    true, Speed: true, Volume: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if effective != (Controls{Pitch: 200, Speed: 50, Volume: 500}) {
		t.Fatalf("effective controls = %+v", effective)
	}
	if got, want := result.Arrays[0][0], int32(200); got != want {
		t.Fatalf("pitch state = %d, want %d", got, want)
	}
	if got, want := result.Arrays[1][0], int32(50); got != want {
		t.Fatalf("speed state = %d, want %d", got, want)
	}
	if got, want := result.Arrays[2][0], int32(500); got != want {
		t.Fatalf("volume state = %d, want %d", got, want)
	}
	if got, want := result.Arrays[3][0], int32(6); got != want {
		t.Fatalf("event state = %d, want %d", got, want)
	}
}

func TestApplyPaul2013PositionStateProgramsForOrdinarySourceWithControls(t *testing.T) {
	programs := []text.Paul2013PositionStateProgram{{
		RowCount: 1,
		Events: map[byte]text.Paul2013PositionEventValues{
			3: {
				Mode: 3, Boundaries: []int32{20}, Values: []int32{8},
				Initial: []int32{-1}, Minimum: 0, Maximum: 100,
			},
		},
	}}
	results, effective, err := ApplyPaul2013PositionStateProgramsForOrdinarySourceWithControls(
		"Hi", 20, programs,
		Paul2013ControlOverrides{Defaults: Controls{Pitch: 80, Speed: 90, Volume: 120}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if effective != (Controls{Pitch: 80, Speed: 90, Volume: 120}) {
		t.Fatalf("effective controls = %+v", effective)
	}
	if len(results) != 1 || results[0].Arrays[0][0] != 80 ||
		results[0].Arrays[1][0] != 90 || results[0].Arrays[2][0] != 120 ||
		results[0].Arrays[3][0] != 8 {
		t.Fatalf("ordinary-source position states = %+v", results)
	}
}

func TestApplyPaul2013PositionStateProgramWithControlsPropagatesProgramErrors(t *testing.T) {
	program := text.Paul2013PositionStateProgram{
		RowCount: 1,
		Ranges: map[byte]text.Paul2013PositionStateRanges{
			3: {Mode: 3},
		},
	}
	_, effective, err := ApplyPaul2013PositionStateProgramWithControls(program, Paul2013ControlOverrides{
		Defaults: Controls{Pitch: 100, Speed: 100, Volume: 200},
	})
	if err == nil {
		t.Fatal("unsupported range pass succeeded")
	}
	if effective != (Controls{Pitch: 100, Speed: 100, Volume: 200}) {
		t.Fatalf("effective controls on error = %+v", effective)
	}
}
