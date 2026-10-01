package selection

import "testing"

func TestReadPaul2013QueryRuntimeState(t *testing.T) {
	state := Paul2013TransitionContextState{2, 0, 5, 0, 0, 0}
	base := int(state[0])*paul2013QueryRuntimeRecordStride + int(state[2])
	arena := make([]byte, base+paul2013QueryContextGateOffset+1)
	arena[base+paul2013QueryPriorMetricOffset] = 27
	arena[base+paul2013QueryContextGateOffset] = 1

	got, err := ReadPaul2013QueryRuntimeState(arena, state)
	if err != nil {
		t.Fatal(err)
	}
	if got.PriorMetric != 27 || got.ContextGate != 1 {
		t.Fatalf("query runtime state = %+v, want metric 27 and gate 1", got)
	}
}

func TestReadPaul2013QueryRuntimeStateRejectsShortArena(t *testing.T) {
	state := Paul2013TransitionContextState{1, 0, 0, 0, 0, 0}
	if _, err := ReadPaul2013QueryRuntimeState(nil, state); err == nil {
		t.Fatal("short query state arena was accepted")
	}
}
