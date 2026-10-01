package duration

import (
	"encoding/binary"
	"path/filepath"
	"testing"

	"vtspeak/engine/text"
)

func TestResolvePaul2013PronunciationExceptionSequence(t *testing.T) {
	dictionary, err := text.LoadExceptionDictionary(filepath.Join("..", "..", "data-common", "dict-eng"))
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{exceptions: dictionary}
	match, phoneRows, found, err := engine.ResolvePaul2013PronunciationExceptionSequence(
		[]string{"San", "Francisco", "unused"}, []int{1, 1},
	)
	if err != nil || !found {
		t.Fatalf("ResolvePaul2013PronunciationExceptionSequence() = (%+v, %v, %t, %v)", match, phoneRows, found, err)
	}
	if match.Category != 2 || match.TokenCount != 2 || len(phoneRows) != 2 {
		t.Fatalf("exception match = %+v with %d phone rows; want category 2 and two destination rows", match, len(phoneRows))
	}
	for rowIndex, row := range phoneRows {
		if len(row) == 0 {
			t.Errorf("exception phone row %d is empty", rowIndex)
		}
	}
	if _, _, _, err := engine.ResolvePaul2013PronunciationExceptionSequence(
		[]string{"San", "Francisco"}, []int{1},
	); err == nil {
		t.Fatal("exception sequence accepted a delimiter-count list with the wrong number of rows")
	}
}

func TestResolvePaul2013PronunciationExceptionFromPhoneContextRows(t *testing.T) {
	dictionary, err := text.LoadExceptionDictionary(filepath.Join("..", "..", "data-common", "dict-eng"))
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{exceptions: dictionary}
	const (
		parserRowSize = 0x94
		contextHeader = text.Paul2013PhoneContextTableHeaderSize
		contextStride = text.Paul2013PhoneContextRowSize
	)
	parserRows := make([]byte, 2*parserRowSize)
	parserRows[0x23] = 'U'
	parserRows[0x24] = 'N'
	parserRows[0x30] = 0xff
	parserRows[parserRowSize+0x23] = 'U'
	parserRows[parserRowSize+0x24] = 'N'
	contextRows := make([]byte, contextHeader+2*contextStride)
	binary.LittleEndian.PutUint16(contextRows, 2)
	first := contextRows[contextHeader : contextHeader+contextStride]
	second := contextRows[contextHeader+contextStride:]
	binary.LittleEndian.PutUint16(first[text.Paul2013PhoneContextRowTokenIndex:], 0)
	binary.LittleEndian.PutUint16(second[text.Paul2013PhoneContextRowTokenIndex:], 1)
	copy(first[text.Paul2013PhoneContextRowSurface:], []byte("San\x00"))
	copy(second[text.Paul2013PhoneContextRowSurface:], []byte("Francisco\x00"))
	first[text.Paul2013PhoneContextRowSourceMarker] = '0'
	second[text.Paul2013PhoneContextRowSourceMarker] = '0'

	match, phones, found, preparationKnown, err := engine.ResolvePaul2013PronunciationExceptionFromPhoneContextRows(
		parserRows, contextRows, 0, []int{0, 1},
	)
	if err != nil || !found || !preparationKnown {
		t.Fatalf("row-based exception resolution = (%+v, %v, %t, %t, %v)",
			match, phones, found, preparationKnown, err)
	}
	if match.Category != 2 || match.TokenCount != 2 || len(phones) != 2 {
		t.Fatalf("row-based exception match = %+v with %d phone rows, want category 2 and two rows", match, len(phones))
	}
}
