package text

import (
	"encoding/binary"
	"fmt"
)

const (
	paul2013ParserStateRowCountOffset = 0x00
	paul2013ParserStateCursorOffset   = 0x04
	paul2013ParserStateRowCursorEnd   = 0x04
	paul2013ParserStateRowSentinel    = 0x0c
)

// RollbackPaul2013ModelParserCandidateRow restores the shared state after a
// FUN_1003d3d0 handler rejects the row at committedRows. The native loop
// restores the committed row count and its last source-end cursor, clears the
// candidate's 0x94-byte slot, and writes 0xffffffff at row +0x0c. The input
// arena is not mutated.
func RollbackPaul2013ModelParserCandidateRow(
	state []byte,
	committedRows int,
) ([]byte, error) {
	if committedRows < 0 || committedRows > paul2013ParserStateRowLimit {
		return nil, fmt.Errorf("Paul 2013 parser committed row count %d is outside 0..%d", committedRows, paul2013ParserStateRowLimit)
	}
	minimum := 0
	if committedRows < paul2013ParserStateRowLimit {
		minimum = paul2013ParserStateRowsOffset + (committedRows+1)*paul2013ParserStateRowStride
	} else {
		lastRow := paul2013ParserStateRowsOffset + (committedRows-1)*paul2013ParserStateRowStride
		minimum = lastRow + paul2013ParserStateRowCursorEnd + 4
	}
	if len(state) < minimum {
		return nil, fmt.Errorf("Paul 2013 parser state has %d bytes, need %d to roll back after %d rows", len(state), minimum, committedRows)
	}

	working := append([]byte(nil), state...)
	binary.LittleEndian.PutUint16(working[paul2013ParserStateRowCountOffset:], uint16(int16(committedRows)))
	if committedRows > 0 {
		lastRow := paul2013ParserStateRowsOffset + (committedRows-1)*paul2013ParserStateRowStride
		cursor := binary.LittleEndian.Uint32(working[lastRow+paul2013ParserStateRowCursorEnd:])
		binary.LittleEndian.PutUint32(working[paul2013ParserStateCursorOffset:], cursor)
	}
	if committedRows < paul2013ParserStateRowLimit {
		candidateRow := paul2013ParserStateRowsOffset + committedRows*paul2013ParserStateRowStride
		clear(working[candidateRow : candidateRow+paul2013ParserStateRowStride])
		binary.LittleEndian.PutUint32(
			working[candidateRow+paul2013ParserStateRowSentinel:],
			^uint32(0),
		)
	}
	return working, nil
}
