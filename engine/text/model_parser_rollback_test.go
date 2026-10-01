package text

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestRollbackPaul2013ModelParserCandidateRow(t *testing.T) {
	const count = 3
	state := make([]byte, paul2013ParserStateRowsOffset+(count+1)*paul2013ParserStateRowStride)
	binary.LittleEndian.PutUint16(state[paul2013ParserStateRowCountOffset:], 9)
	binary.LittleEndian.PutUint32(state[paul2013ParserStateCursorOffset:], 0xdeadbeef)
	for rowIndex := 0; rowIndex < count+1; rowIndex++ {
		rowStart := paul2013ParserStateRowsOffset + rowIndex*paul2013ParserStateRowStride
		for offset := 0; offset < paul2013ParserStateRowStride; offset++ {
			state[rowStart+offset] = byte(rowIndex + 1)
		}
		binary.LittleEndian.PutUint32(state[rowStart+paul2013ParserStateRowCursorEnd:], uint32(0x100+rowIndex))
	}
	before := append([]byte(nil), state...)

	got, err := RollbackPaul2013ModelParserCandidateRow(state, count)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(state, before) {
		t.Fatal("rollback mutated its input arena")
	}
	if binary.LittleEndian.Uint16(got[paul2013ParserStateRowCountOffset:]) != count {
		t.Fatalf("row count = %d, want %d", binary.LittleEndian.Uint16(got), count)
	}
	if cursor := binary.LittleEndian.Uint32(got[paul2013ParserStateCursorOffset:]); cursor != 0x102 {
		t.Fatalf("source cursor = %#x, want %#x", cursor, 0x102)
	}
	if !bytes.Equal(got[paul2013ParserStateRowsOffset:paul2013ParserStateRowsOffset+count*paul2013ParserStateRowStride], before[paul2013ParserStateRowsOffset:paul2013ParserStateRowsOffset+count*paul2013ParserStateRowStride]) {
		t.Fatal("rollback changed a committed parser row")
	}
	candidate := paul2013ParserStateRowsOffset + count*paul2013ParserStateRowStride
	wantCandidate := make([]byte, paul2013ParserStateRowStride)
	binary.LittleEndian.PutUint32(wantCandidate[paul2013ParserStateRowSentinel:], ^uint32(0))
	if !bytes.Equal(got[candidate:candidate+paul2013ParserStateRowStride], wantCandidate) {
		t.Fatalf("candidate row = % x, want zeroed row with sentinel at +0x%x", got[candidate:candidate+paul2013ParserStateRowStride], paul2013ParserStateRowSentinel)
	}
	if !bytes.Equal(got[candidate+paul2013ParserStateRowStride:], before[candidate+paul2013ParserStateRowStride:]) {
		t.Fatal("rollback changed bytes after the candidate row")
	}
}

func TestRollbackPaul2013ModelParserCandidateRowBoundaries(t *testing.T) {
	t.Run("first row preserves cursor", func(t *testing.T) {
		state := bytes.Repeat([]byte{0xa5}, paul2013ParserStateRowsOffset+paul2013ParserStateRowStride)
		binary.LittleEndian.PutUint32(state[paul2013ParserStateCursorOffset:], 0x12345678)
		got, err := RollbackPaul2013ModelParserCandidateRow(state, 0)
		if err != nil {
			t.Fatal(err)
		}
		if binary.LittleEndian.Uint16(got) != 0 || binary.LittleEndian.Uint32(got[paul2013ParserStateCursorOffset:]) != 0x12345678 {
			t.Fatalf("row count/cursor = %d/%#x", binary.LittleEndian.Uint16(got), binary.LittleEndian.Uint32(got[paul2013ParserStateCursorOffset:]))
		}
		candidate := got[paul2013ParserStateRowsOffset:]
		if binary.LittleEndian.Uint32(candidate[paul2013ParserStateRowSentinel:]) != ^uint32(0) {
			t.Fatalf("first candidate sentinel = %#x", binary.LittleEndian.Uint32(candidate[paul2013ParserStateRowSentinel:]))
		}
	})

	t.Run("full row table restores last committed cursor", func(t *testing.T) {
		state := make([]byte, paul2013ParserStateRowsOffset+paul2013ParserStateRowLimit*paul2013ParserStateRowStride)
		lastRow := paul2013ParserStateRowsOffset + (paul2013ParserStateRowLimit-1)*paul2013ParserStateRowStride
		binary.LittleEndian.PutUint32(state[lastRow+paul2013ParserStateRowCursorEnd:], 0x76543210)
		got, err := RollbackPaul2013ModelParserCandidateRow(state, paul2013ParserStateRowLimit)
		if err != nil {
			t.Fatal(err)
		}
		if cursor := binary.LittleEndian.Uint32(got[paul2013ParserStateCursorOffset:]); cursor != 0x76543210 {
			t.Fatalf("full-table cursor = %#x, want %#x", cursor, 0x76543210)
		}
	})

	for _, count := range []int{-1, paul2013ParserStateRowLimit + 1} {
		if _, err := RollbackPaul2013ModelParserCandidateRow(make([]byte, paul2013ParserStateRowsOffset), count); err == nil {
			t.Errorf("committed row count %d was accepted", count)
		}
	}
	if _, err := RollbackPaul2013ModelParserCandidateRow(make([]byte, paul2013ParserStateRowsOffset), 0); err == nil {
		t.Fatal("state without the current candidate row was accepted")
	}
}
