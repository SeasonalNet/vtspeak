package text

import (
	"encoding/binary"
	"fmt"
)

const (
	paul2013ModelTokenRowCountOffset = 0
	paul2013ModelTokenRowsOffset     = 2
	paul2013ModelTokenRowLimit       = 200
)

// BuildPaul2013ModelTokenRows projects already-resolved lexical tokens into
// the model's 0x554-byte token-result rows written by FUN_1000d450. Parser row
// indexes remain explicit because source-row production and its association
// with resolved tokens are not recovered. Existing bytes that FUN_1000d450
// leaves untouched are preserved.
func BuildPaul2013ModelTokenRows(
	model []byte,
	tokens []LexicalToken,
	parserRows []byte,
	tokenParserIndexes []uint16,
) ([]byte, error) {
	if len(tokens) != len(tokenParserIndexes) {
		return nil, fmt.Errorf("model token projection has %d tokens and %d parser row indexes", len(tokens), len(tokenParserIndexes))
	}
	if len(tokens) > paul2013ModelTokenRowLimit {
		return nil, fmt.Errorf("model token projection has %d rows, native limit is %d", len(tokens), paul2013ModelTokenRowLimit)
	}
	if len(parserRows)%paul2013ModelParserOffsetRowStride != 0 {
		return nil, fmt.Errorf("parser offset rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013ModelParserOffsetRowStride)
	}
	parserRowCount := len(parserRows) / paul2013ModelParserOffsetRowStride
	for tokenIndex, parserIndex := range tokenParserIndexes {
		if int(parserIndex) >= parserRowCount {
			return nil, fmt.Errorf("model token %d references parser row %d outside %d rows", tokenIndex, parserIndex, parserRowCount)
		}
	}
	tokenBytes := len(tokens) * Paul2013TokenResultRowSize
	if tokenBytes > paul2013ModelContextCountOffset-paul2013ModelTokenRowsOffset {
		return nil, fmt.Errorf("%d token rows overlap the model context table", len(tokens))
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return nil, fmt.Errorf("model buffer has %d bytes, need at least %d for token rows and context count", len(model), paul2013ModelContextCountOffset+2)
	}

	working := append([]byte(nil), model...)
	for tokenIndex, token := range tokens {
		rowStart := paul2013ModelTokenRowsOffset + tokenIndex*Paul2013TokenResultRowSize
		rowEnd := rowStart + Paul2013TokenResultRowSize
		if err := token.WritePaul2013DictionaryPhoneRow(
			working[rowStart:rowEnd], tokenParserIndexes[tokenIndex],
		); err != nil {
			return nil, fmt.Errorf("write model token row %d (%q): %w", tokenIndex, token.Surface, err)
		}
	}
	binary.LittleEndian.PutUint16(working[paul2013ModelTokenRowCountOffset:], uint16(len(tokens)))
	return working, nil
}
