package text

import (
	"encoding/binary"
	"testing"
)

func TestWritePaul2013PhoneContextRow(t *testing.T) {
	destination := make([]byte, Paul2013PhoneContextRowSize)
	for index := range destination {
		destination[index] = 0xa5
	}
	if err := WritePaul2013PhoneContextRow(
		destination, 0x1234, []byte("Hello"), 'A', []byte("AH0S1"), 'X', [4]bool{true, false, true, true},
	); err != nil {
		t.Fatal(err)
	}
	if got := binary.LittleEndian.Uint16(destination[Paul2013PhoneContextRowTokenIndex:]); got != 0x1234 {
		t.Errorf("token index = %#x, want %#x", got, 0x1234)
	}
	if destination[Paul2013PhoneContextRowStatus] != 'A' ||
		string(destination[Paul2013PhoneContextRowSurface:Paul2013PhoneContextRowSurface+6]) != "Hello\x00" ||
		string(destination[Paul2013PhoneContextRowPhoneCodes:Paul2013PhoneContextRowPhoneCodes+6]) != "AH0S1\x00" ||
		destination[Paul2013PhoneContextRowSourceMarker] != 'X' || destination[Paul2013PhoneContextRowFlags] != 1 {
		t.Fatalf("phone-context row fields are incorrect: % x", destination)
	}
	for index, want := range []uint16{1, 0, 1, 1} {
		offset := Paul2013PhoneContextRowMetadata + index*2
		if got := binary.LittleEndian.Uint16(destination[offset:]); got != want {
			t.Errorf("metadata %d = %d, want %d", index, got, want)
		}
	}
	if destination[0x04] != 0xa5 || destination[0x05] != 0xa5 || destination[0x67] != 0xa5 {
		t.Fatal("phone-context writer changed bytes that FUN_1000ea20 leaves untouched")
	}
	if err := WritePaul2013PhoneContextRow(destination, 1, []byte("Hello"), 'A', []byte("AH0S"), '0', [4]bool{}); err != nil {
		t.Fatal(err)
	}
	if destination[Paul2013PhoneContextRowFlags] != 0 {
		t.Fatal("phone-context initial flag is set when byte +4 of the phone string is NUL")
	}
}

func TestBuildPaul2013PhoneContextRows(t *testing.T) {
	sequence := LexicalPhoneSequence{Tokens: []LexicalTokenSpan{
		{
			SourceSurface: "Hello", Surface: "Hello", AlternativeIndex: 1,
			DictionaryMetadata: [4]bool{true, false, true, false},
			ModelPhoneRows: Paul2013DictionaryPhoneRows{
				ResultType: 'A', PhoneStrings: [][]byte{[]byte("HH0"), []byte("AH0S")},
			},
		},
		{
			SourceSurface: "there", Surface: "there", AlternativeIndex: 0,
			ModelPhoneRows: Paul2013DictionaryPhoneRows{
				PhoneStrings: [][]byte{[]byte("DH1ER0")},
			},
		},
	}}
	got, err := BuildPaul2013PhoneContextRows(sequence, []uint16{7, 8}, []byte{'X', 'A'})
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint16(got) != 2 || len(got) != Paul2013PhoneContextTableHeaderSize+2*Paul2013PhoneContextRowSize {
		t.Fatalf("phone-context table header/length = %d/%d", binary.LittleEndian.Uint16(got), len(got))
	}
	first := got[Paul2013PhoneContextTableHeaderSize : Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize]
	second := got[Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize:]
	if binary.LittleEndian.Uint16(first[Paul2013PhoneContextRowTokenIndex:]) != 7 ||
		string(first[Paul2013PhoneContextRowPhoneCodes:Paul2013PhoneContextRowPhoneCodes+5]) != "AH0S\x00" ||
		first[Paul2013PhoneContextRowSourceMarker] != 'X' {
		t.Fatalf("first projected row = % x", first)
	}
	if binary.LittleEndian.Uint16(second[Paul2013PhoneContextRowTokenIndex:]) != 8 ||
		string(second[Paul2013PhoneContextRowPhoneCodes:Paul2013PhoneContextRowPhoneCodes+7]) != "DH1ER0\x00" ||
		second[Paul2013PhoneContextRowSourceMarker] != '0' {
		t.Fatalf("second projected row = % x", second)
	}
	sourceParserRow := make([]byte, 0x94)
	sourceParserRow[0x23] = 'U'
	sourceParserRow[0x30] = 0xff
	gate, err := BuildPaul2013ExceptionDispatchGateFromRows(sourceParserRow, 0, first, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !ShouldDispatchPaul2013PronunciationException(gate) {
		t.Fatalf("row-built exception dispatch gate = %+v, want dispatch", gate)
	}
	if _, err := BuildPaul2013PhoneContextRows(sequence, []uint16{7}, []byte{'X', 'A'}); err == nil {
		t.Fatal("misaligned parser input arrays were accepted")
	}
}

func TestBuildPaul2013PhoneContextRowsFromParserRows(t *testing.T) {
	sequence := LexicalPhoneSequence{Tokens: []LexicalTokenSpan{
		{
			SourceSurface: "Alpha", Surface: "Alpha", AlternativeIndex: 0,
			ModelPhoneRows: Paul2013DictionaryPhoneRows{PhoneStrings: [][]byte{[]byte("AE1")}},
		},
		{
			SourceSurface: "Beta", Surface: "Beta", AlternativeIndex: 0,
			ModelPhoneRows: Paul2013DictionaryPhoneRows{PhoneStrings: [][]byte{[]byte("B EY1 T AH0")}},
		},
	}}
	parserRows := make([]byte, 3*paul2013PhoneContextParserRowStride)
	parserRows[2*paul2013PhoneContextParserRowStride+paul2013PhoneContextParserClassByte] = 'X'
	rows, err := BuildPaul2013PhoneContextRowsFromParserRows(sequence, parserRows, []uint16{0, 2})
	if err != nil {
		t.Fatal(err)
	}
	first := rows[Paul2013PhoneContextTableHeaderSize : Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize]
	second := rows[Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize:]
	if first[Paul2013PhoneContextRowSourceMarker] != '0' || second[Paul2013PhoneContextRowSourceMarker] != 'X' {
		t.Fatalf("parser-row-derived source markers = %q/%q, want '0'/'X'", first[Paul2013PhoneContextRowSourceMarker], second[Paul2013PhoneContextRowSourceMarker])
	}
	if _, err := BuildPaul2013PhoneContextRowsFromParserRows(sequence, parserRows[:len(parserRows)-1], []uint16{0, 2}); err == nil {
		t.Fatal("misaligned parser-row buffer was accepted")
	}
	if _, err := BuildPaul2013PhoneContextRowsFromParserRows(sequence, parserRows, []uint16{0, 3}); err == nil {
		t.Fatal("out-of-range parser-row index was accepted")
	}
}

func TestExtractPaul2013ExceptionRowInputs(t *testing.T) {
	table := make([]byte, Paul2013PhoneContextTableHeaderSize+2*Paul2013PhoneContextRowSize)
	binary.LittleEndian.PutUint16(table, 2)
	first := table[Paul2013PhoneContextTableHeaderSize : Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize]
	second := table[Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize:]
	copy(first[Paul2013PhoneContextRowSurface:], []byte("San\x00"))
	first[Paul2013PhoneContextRowSourceMarker] = 'X'
	copy(second[Paul2013PhoneContextRowSurface:], []byte("Francisco\x00"))
	second[Paul2013PhoneContextRowSourceMarker] = '0'

	got, err := ExtractPaul2013ExceptionRowInputs(table, []int{1, 0})
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013ExceptionRowInput{
		{RowIndex: 1, Surface: "Francisco"},
		{RowIndex: 0, Surface: "San", RetryHPrefix: true},
	}
	if len(got) != len(want) {
		t.Fatalf("extracted %d row inputs, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Errorf("row input %d = %+v, want %+v", index, got[index], want[index])
		}
	}
	if _, err := ExtractPaul2013ExceptionRowInputs(table, []int{2}); err == nil {
		t.Fatal("out-of-range phone-context row was accepted")
	}
	for index := Paul2013PhoneContextRowSurface; index < Paul2013PhoneContextRowPhoneCodes; index++ {
		second[index] = 'A'
	}
	if _, err := ExtractPaul2013ExceptionRowInputs(table, []int{1}); err == nil {
		t.Fatal("unterminated phone-context surface was accepted")
	}
}
