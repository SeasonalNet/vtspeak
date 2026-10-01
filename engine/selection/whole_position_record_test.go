package selection

import (
	"testing"

	"vtspeak/engine/voice"
)

func TestBuildPaul2013WholePositionRecordInput(t *testing.T) {
	const blockOrdinal = 1
	blockOffset := blockOrdinal * paul2013QueryRecordBlockStride
	arena := make([]byte, blockOffset+paul2013QueryModelClassOffset+1)
	arena[blockOffset+paul2013QueryRowClassOffset] = 12
	arena[blockOffset+paul2013QueryModelClassOffset] = 8
	state := [6]byte{byte(blockOrdinal), 4, 2, 99, 77, 6}
	wantSignature := [7]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27}
	signatureOffset := blockOffset + paul2013QuerySignatureOffset + int(state[2])*7
	copy(arena[signatureOffset:], wantSignature[:])

	got, err := BuildPaul2013WholePositionRecordInput(arena, state)
	if err != nil {
		t.Fatal(err)
	}
	if got.Target.Signature != wantSignature {
		t.Fatalf("target signature = % x, want % x", got.Target.Signature, wantSignature)
	}
	if got.ModelClass != 8 {
		t.Fatalf("model class = %d, want 8", got.ModelClass)
	}
	wantRow := state
	wantRow[3] = 12
	wantRow[4] = 0
	if got.ContextRow != wantRow {
		t.Fatalf("query context row = % x, want % x", got.ContextRow, wantRow)
	}
}

func TestBuildPaul2013WholePositionRecordInputRejectsShortArena(t *testing.T) {
	if _, err := BuildPaul2013WholePositionRecordInput(nil, [6]byte{0, 0, 0xff}); err == nil {
		t.Fatal("missing signature was accepted")
	}
	arena := make([]byte, paul2013QuerySignatureOffset+7)
	if _, err := BuildPaul2013WholePositionRecordInput(arena, [6]byte{}); err == nil {
		t.Fatal("arena without model-class byte was accepted")
	}
}

func TestQueryPaul2013WholePositionFromTransitionState(t *testing.T) {
	const blockOrdinal = 1
	blockOffset := blockOrdinal * paul2013QueryRecordBlockStride
	arena := make([]byte, blockOffset+paul2013QueryModelClassOffset+1)
	arena[blockOffset+paul2013QueryRowClassOffset] = 1
	arena[blockOffset+paul2013QueryModelClassOffset] = 8
	copy(arena[blockOffset+paul2013QuerySignatureOffset+2*7:], []byte{1, 2, 3, 4, 5, 6, 7})
	catalog := &catalogFeatureQueryFixture{featureViewRangeFixture: featureViewRangeFixture{
		classes: []voice.ClassRecord{
			{ID: 1, Key: [5]byte{1}, Members: make([]voice.UnitLocation, 5)},
			{ID: 2, Key: [5]byte{2}, Members: make([]voice.UnitLocation, 5)},
		},
	}}
	state := Paul2013TransitionContextState{blockOrdinal, 0, 2, 0, 0, 0}

	result, err := QueryPaul2013WholePositionFromTransitionState(
		catalog, arena, state, 0, func() byte { return 0 },
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Passes) == 0 || len(result.Passes[0].Sequence.Attempts) == 0 {
		t.Fatalf("transition-state whole-position query produced no query attempts: %#v", result)
	}
	wantSignature := [7]byte{1, 2, 3, 4, 5, 6, 7}
	if got := result.Passes[0].Sequence.Attempts[0].Signature; got != wantSignature {
		t.Fatalf("query signature = % x, want record selected by transition state % x", got, wantSignature)
	}
}

func TestQueryPaul2013PositionFromTransitionStateUsesRecordFallbackState(t *testing.T) {
	const blockOrdinal = 1
	blockOffset := blockOrdinal * paul2013QueryRecordBlockStride
	arena := make([]byte, blockOffset+paul2013QueryModelClassOffset+1)
	arena[blockOffset+paul2013QueryRowClassOffset] = 7
	arena[blockOffset+paul2013QueryModelClassOffset] = 8
	copy(arena[blockOffset+paul2013QuerySignatureOffset+7:], []byte{1, 2, 3, 4, 5, 6, 7})
	catalog := &catalogFeatureQueryFixture{}
	state := Paul2013TransitionContextState{blockOrdinal, 0, 1, 0, 0, 0}
	runtime := [2]Paul2013FallbackRowRuntime{
		{ContextGate: func() byte { return 0 }},
		{ContextGate: func() byte { return 0 }},
	}

	result, err := QueryPaul2013PositionFromTransitionState(
		catalog, arena, state, 0, func() byte { return 0 }, runtime,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !result.UsedFallback || result.ReturnCount != 2 {
		t.Fatalf("record-backed position query = %+v, want native two-row fallback", result)
	}
	for index, row := range result.FallbackRows {
		if row.Row.Bytes[0] != byte(blockOrdinal) || row.Row.Bytes[1] != 0 || row.Row.Bytes[2] != 1 {
			t.Errorf("fallback row %d base state = % x, want transition selectors preserved", index, row.Row.Bytes)
		}
		if row.Row.Bytes[3] != byte(7+index) || row.Row.Bytes[4] != byte(index+1) {
			t.Errorf("fallback row %d class/mode = %d/%d, want %d/%d", index, row.Row.Bytes[3], row.Row.Bytes[4], 7+index, index+1)
		}
	}
}

func TestQueryPaul2013PositionFromStateArenasDerivesRuntimeInputs(t *testing.T) {
	const blockOrdinal = 1
	state := Paul2013TransitionContextState{blockOrdinal, 0, 1, 0, 0, 0}
	blockOffset := int(state[0])*paul2013QueryRuntimeRecordStride + int(state[2])
	queryStateArena := make([]byte, blockOffset+paul2013QueryContextGateOffset+1)
	queryStateArena[blockOffset+paul2013QueryPriorMetricOffset] = 27
	queryStateArena[blockOffset+paul2013QueryContextGateOffset] = 1

	recordBlockOffset := int(state[0]) * paul2013QueryRecordBlockStride
	recordArena := make([]byte, recordBlockOffset+paul2013QueryModelClassOffset+1)
	recordArena[recordBlockOffset+paul2013QueryRowClassOffset] = 7
	recordArena[recordBlockOffset+paul2013QueryModelClassOffset] = 8
	catalog := &catalogFeatureQueryFixture{}

	result, err := QueryPaul2013PositionFromStateArenas(catalog, recordArena, queryStateArena, state)
	if err != nil {
		t.Fatal(err)
	}
	if !result.UsedFallback || result.ReturnCount != 2 {
		t.Fatalf("arena-backed position query = %+v, want conditional two-row fallback", result)
	}
	if result.FallbackRows[0].Row.Bytes[3] != 7 || result.FallbackRows[1].Row.Bytes[3] != 8 {
		t.Fatalf("arena-backed fallback row classes = %d,%d, want signed record-byte sequence 7,8", result.FallbackRows[0].Row.Bytes[3], result.FallbackRows[1].Row.Bytes[3])
	}
}
