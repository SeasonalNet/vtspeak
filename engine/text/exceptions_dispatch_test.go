package text

import "testing"

func TestBuildPaul2013ExceptionDispatchGateFromRows(t *testing.T) {
	const rowSize = 0x94
	sourceRows := make([]byte, 2*rowSize)
	selected := sourceRows[rowSize : 2*rowSize]
	selected[0x23] = 'U'
	selected[0x30] = 0xff
	selected[0x52] = 0
	phoneRow := make([]byte, Paul2013PhoneContextRowSize)
	phoneRow[0x29] = 0

	gate, err := BuildPaul2013ExceptionDispatchGateFromRows(sourceRows, 1, phoneRow, 0)
	if err != nil {
		t.Fatal(err)
	}
	if gate.SourceTokenStatus != -1 || gate.PhoneRowMarker != 0 || gate.SourceKind != 'U' ||
		!gate.ContextFormIsEmpty || gate.ContextPreparationCode != 0 ||
		!ShouldDispatchPaul2013PronunciationException(gate) {
		t.Fatalf("extracted dispatch gate = %+v", gate)
	}

	selected[0x30] = 1
	selected[0x52] = 'N'
	phoneRow[0x29] = 'a'
	gate, err = BuildPaul2013ExceptionDispatchGateFromRows(sourceRows, 1, phoneRow, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ShouldDispatchPaul2013PronunciationException(gate) {
		t.Fatalf("gate unexpectedly dispatches without any matching predicate: %+v", gate)
	}

	if _, err := BuildPaul2013ExceptionDispatchGateFromRows(sourceRows, 2, phoneRow, 0); err == nil {
		t.Fatal("out-of-range source token index was accepted")
	}
	if _, err := BuildPaul2013ExceptionDispatchGateFromRows(sourceRows, 1, phoneRow[:0x6f], 0); err == nil {
		t.Fatal("short phone/context row was accepted")
	}
}
