package text

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func finalizerModelStateFixture(mode, class byte, codes string) Paul2013ModelParserFinalizerInput {
	const stateOffset = 0x80
	state := make([]byte, stateOffset+0x94)
	binary.LittleEndian.PutUint16(state[stateOffset:], 1)
	binary.LittleEndian.PutUint32(state[stateOffset+0x14:], 5)
	binary.LittleEndian.PutUint32(state[stateOffset+0x18:], 9)
	state[stateOffset+0x37] = mode
	state[stateOffset+0x38] = class
	binary.LittleEndian.PutUint16(state[stateOffset+0x40:], 3)
	copy(state[stateOffset+0x66:], "B\x00")
	table := make([]byte, 0x0b+0x70)
	copy(table[0x0b+0x1e:], codes)
	return Paul2013ModelParserFinalizerInput{
		StateArena: state, ParserStateOffset: stateOffset,
		ContextTable: table, ContextRowCount: 1,
	}
}

func TestFinalizerSurfaceCopyGateMatchesNativeBranches(t *testing.T) {
	for _, test := range []struct {
		name        string
		mode, class byte
		want        string
	}{
		{name: "U Y copies", mode: 'U', class: 'Y', want: "A"},
		{name: "U A copies", mode: 'U', class: 'A', want: "A"},
		{name: "A Y preserves", mode: 'A', class: 'Y', want: "B"},
		{name: "A A copies", mode: 'A', class: 'A', want: "A"},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := finalizerModelStateFixture(test.mode, test.class, "A\x00")
			got, err := FinalizePaul2013ModelParserRows(input)
			if err != nil {
				t.Fatal(err)
			}
			area := got.StateArena[input.ParserStateOffset+0x66:]
			end := bytes.IndexByte(area, 0)
			if end < 0 || string(area[:end]) != test.want {
				t.Fatalf("finalized phone string = %x, want %q", area[:2], test.want)
			}
		})
	}
}

func TestFinalizedContextTableBuildsCountedModelStateRecords(t *testing.T) {
	input := finalizerModelStateFixture('U', 'Y', "ABdMc\x00")
	stateBefore := append([]byte(nil), input.StateArena...)
	arena := make([]byte, 0x4770b)
	got, err := FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(
		input, arena, 100, []int32{7},
		func(rowIndex, rowOffset int) (uint32, error) {
			return uint32(0x10000000 + rowIndex*0x94 + rowOffset), nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(input.StateArena, stateBefore) || !bytes.Equal(arena, make([]byte, len(arena))) {
		t.Fatal("finalized model-state build mutated caller inputs")
	}
	if got.Finalized.ReturnValue != 1 || got.Finalized.OutputGroupCount != 1 ||
		binary.LittleEndian.Uint16(got.Records.Bytes[2:4]) != 1 {
		t.Fatalf("finalized/count outputs = %d/%d/%d, want 1/1/1", got.Finalized.ReturnValue, got.Finalized.OutputGroupCount, binary.LittleEndian.Uint16(got.Records.Bytes[2:4]))
	}
	record := got.Records.Bytes[0x64c:0xa0c]
	if start, end := int32(binary.LittleEndian.Uint32(record[0:4])), int32(binary.LittleEndian.Uint32(record[4:8])); start != 105 || end != 108 {
		t.Fatalf("model-state interval = %d..%d, want 105..108", start, end)
	}
	if record[0x95] != 3 || string(record[0x2e8:0x2eb]) != "ABM" ||
		!bytes.Equal(record[0x328:0x32c], []byte{0, '1', '2', '0'}) {
		t.Fatalf("model-state phone fields = %d/%q/%x", record[0x95], record[0x2e8:0x2eb], record[0x328:0x32c])
	}
	if got.Records.Bytes[0x4770a] != 6 || record[0x2df] != 12 {
		t.Fatalf("model-state mode/flag = %d/%d, want 6/12", got.Records.Bytes[0x4770a], record[0x2df])
	}
}

func TestFinalizedModelStateRunsNativeBoundaryPipeline(t *testing.T) {
	input := finalizerModelStateFixture('U', 'Y', "\x01\x13\x00")
	got, err := FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(input, make([]byte, 0x4770b), 0, []int32{7}, func(rowIndex, rowOffset int) (uint32, error) { return uint32(0x1000 + rowIndex*0x94 + rowOffset), nil })
	if err != nil {
		t.Fatal(err)
	}
	boundaries, err := got.RunTokenBoundaries(0x20000000, input.ParserStateOffset, nil, nil, []int32{-1}, []int32{4})
	if err != nil {
		t.Fatal(err)
	}
	if boundaries.Prepared.PhoneGroups.Count != 1 || boundaries.Final.State.Markers[0].Marker != '^' || boundaries.Final.Groups[0].OpaqueRowSpan != 1 {
		t.Fatalf("finalized boundary pipeline = %+v", boundaries.Final)
	}
}

func TestFinalizedModelStateRejectsGroupCountMismatch(t *testing.T) {
	input := finalizerModelStateFixture('U', 'Y', "A\x00")
	binary.LittleEndian.PutUint16(input.StateArena[input.ParserStateOffset:], 2)
	got, err := FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(input, make([]byte, 0x4770b), 0, []int32{7}, nil)
	if err == nil || got.Finalized.ReturnValue != -1 {
		t.Fatalf("group mismatch = return %d, error %v; want rejected -1 finalizer", got.Finalized.ReturnValue, err)
	}
}

func TestFinalizedModelStateGroupsContextRowsBeforeRecordProjection(t *testing.T) {
	input := finalizerModelStateFixture('U', 'Y', "A\x00")
	state := make([]byte, input.ParserStateOffset+2*0x94)
	copy(state, input.StateArena)
	input.StateArena = state
	binary.LittleEndian.PutUint16(state[input.ParserStateOffset:], 2)
	second := state[input.ParserStateOffset+0x94:]
	second[0x37], second[0x38] = 'U', 'Y'
	binary.LittleEndian.PutUint32(second[0x14:], 12)
	binary.LittleEndian.PutUint32(second[0x18:], 12)
	binary.LittleEndian.PutUint16(second[0x40:], 3)
	input.ContextTable = make([]byte, 0x0b+3*0x70)
	input.ContextRowCount, input.ModelResult = 3, 1
	for index, code := range []byte{'A', 'B', 'M'} {
		rowStart := 0x0b + index*0x70
		input.ContextTable[rowStart+0x1e] = code
		if index == 2 {
			binary.LittleEndian.PutUint16(input.ContextTable[rowStart-5:], 1)
		}
	}
	got, err := FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(
		input, make([]byte, 0x4770b), 100, []int32{-1, 7},
		func(rowIndex, rowOffset int) (uint32, error) { return uint32(0x1000 + rowIndex*0x94 + rowOffset), nil },
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.Finalized.OutputGroupCount != 2 || len(got.Records.ContextCodes) != 2 ||
		string(got.Records.ContextCodes[0].Codes) != "AB" || string(got.Records.ContextCodes[1].Codes) != "M" {
		t.Fatalf("grouped model records = %d/%+v, want two groups AB/M", got.Finalized.OutputGroupCount, got.Records.ContextCodes)
	}
	secondRecord := got.Records.Bytes[0x64c+0x3c0:]
	if start, end := int32(binary.LittleEndian.Uint32(secondRecord[0:4])), int32(binary.LittleEndian.Uint32(secondRecord[4:8])); start != 112 || end != 112 {
		t.Fatalf("second record interval = %d..%d, want point 112", start, end)
	}
	if got.Records.Bytes[0x4770a] != 7 || secondRecord[0x2df] != 12 || secondRecord[0x29e] != 1 {
		t.Fatalf("second record mode/state/M = %d/%d/%d, want 7/12/1", got.Records.Bytes[0x4770a], secondRecord[0x2df], secondRecord[0x29e])
	}
	initialized, err := InitializePaul2013TokenBoundaryArena(got.Records.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	// This fixture supplies an explicitly unchanged marker-tree result; it
	// proves the subsequent handoff, not production marker-tree behavior.
	binary.LittleEndian.PutUint32(got.Finalized.StateArena[input.ParserStateOffset+0x20:], 23)
	finished, err := got.FinishTokenBoundaries(initialized, input.ParserStateOffset, []int32{-1, -1}, []int32{4, 4})
	if err != nil {
		t.Fatal(err)
	}
	if finished.State.Values[0] != 23 || finished.State.Markers[1].Marker != 'Z' || finished.Arena[0x64c+0x3c0+0x3bc] != 0 {
		t.Fatalf("finalized boundary handoff = %+v", finished.State)
	}
	if _, err := got.FinishTokenBoundaries(initialized, -1, []int32{-1, -1}, []int32{4, 4}); err == nil {
		t.Fatal("accepted invalid parser base")
	}
}
