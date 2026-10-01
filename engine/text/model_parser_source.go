package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	paul2013ParserSourceRowStride         = 0x94
	paul2013ParserSourceRowRawClassOffset = 0x24
	paul2013ParserSourceRowTypeOffset     = 0x2c
	paul2013ParserSourceRowModeOffset     = 0x30
	paul2013ParserSourceRowLengthOffset   = 0x08
	paul2013ParserSourceRowTextOffset     = 0x34
	paul2013ParserSourceRowAuxOffset      = 0x52
)

// Paul2013ModelParserSourceRowResult is the bounded row-to-string portion of
// FUN_1000d190. SourceType is FUN_1000fd20's class byte for the transformed
// output; the later embedded lookup and token-result write remain separate.
type Paul2013ModelParserSourceRowResult struct {
	RowType        uint32
	MarkerAppended bool
	Transform      Paul2013ModelSourceTransformResult
	SourceType     byte
}

// TransformPaul2013ModelParserSourceRow reads one 0x94-byte source row and
// composes its fixed row-type marker with the recovered FUN_1000d640 source
// transform and FUN_1000fd20 classifier. Row index/count supply only the
// final-row predicate. The source lookup callback remains explicit for the
// non-final punctuation branch that needs it.
func TransformPaul2013ModelParserSourceRow(
	row []byte,
	rowIndex int,
	rowCount int,
	lookup func([]byte) (bool, error),
) (Paul2013ModelParserSourceRowResult, error) {
	if len(row) < paul2013ParserSourceRowStride {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("parser source row has %d bytes, need %d", len(row), paul2013ParserSourceRowStride)
	}
	if rowCount < 1 || rowCount > paul2013SourceRowLimit || rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("parser source row index/count %d/%d is outside the native 0..%d row range", rowIndex, rowCount, paul2013SourceRowLimit)
	}
	text, err := paul2013ParserSourceCString(row, paul2013ParserSourceRowTextOffset, paul2013ParserSourceRowAuxOffset)
	if err != nil {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("read parser source-row text: %w", err)
	}
	if len(text) > 31 {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("parser source-row text has %d bytes, exceeds FUN_1000d190's 32-byte local buffer", len(text))
	}
	auxiliaryText, err := paul2013ParserSourceCString(row, paul2013ParserSourceRowAuxOffset, paul2013ParserSourceRowStride)
	if err != nil {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("read parser source-row auxiliary text: %w", err)
	}
	rowType := binary.LittleEndian.Uint32(row[paul2013ParserSourceRowTypeOffset:])
	headerValue := int32(binary.LittleEndian.Uint32(row[paul2013ParserSourceRowLengthOffset:]))
	marked, markerAppended, err := AppendPaul2013ParserRowTypeMarker(Paul2013ParserRowMarkerInput{
		HeaderValue: headerValue,
		RowType:     rowType,
		Text:        text,
	})
	if err != nil {
		return Paul2013ModelParserSourceRowResult{}, err
	}
	class := row[paul2013ParserSourceRowRawClassOffset]
	previousCount := int16(0)
	if class == 'S' {
		previousCount = 1
	}
	transformed, err := TransformPaul2013ModelSourceStringFromPaul2013CharacterTable(
		Paul2013ModelSourceTransformInput{
			Source:         marked,
			ParserType:     int8(row[paul2013ParserSourceRowModeOffset]),
			AssociatedText: auxiliaryText,
			PreviousCount:  previousCount,
			FinalRow:       rowIndex == rowCount-1,
		},
		lookup,
	)
	if err != nil {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("transform parser source row %d: %w", rowIndex, err)
	}
	return Paul2013ModelParserSourceRowResult{
		RowType:        rowType,
		MarkerAppended: markerAppended,
		Transform:      transformed,
		SourceType:     ClassifyPaul2013ModelSourceType(transformed.Output, Paul2013ExceptionCharacterAttributes()),
	}, nil
}

// TransformPaul2013ModelParserSourceRowWithEmbeddedDictionary supplies the
// loaded embedded lexicon for FUN_10003a70's non-final punctuation lookup.
func TransformPaul2013ModelParserSourceRowWithEmbeddedDictionary(
	row []byte,
	rowIndex int,
	rowCount int,
	dictionary *EmbeddedDictionary,
) (Paul2013ModelParserSourceRowResult, error) {
	if dictionary == nil {
		return Paul2013ModelParserSourceRowResult{}, fmt.Errorf("parser source-row transform has no embedded dictionary")
	}
	lookup := func(surface []byte) (bool, error) {
		return dictionary.ContainsPaul2013Surface(surface)
	}
	return TransformPaul2013ModelParserSourceRow(row, rowIndex, rowCount, lookup)
}

func paul2013ParserSourceCString(row []byte, offset, limit int) ([]byte, error) {
	if offset < 0 || limit < offset || limit > len(row) {
		return nil, fmt.Errorf("C-string range [%#x,%#x) exceeds row size %d", offset, limit, len(row))
	}
	area := row[offset:limit]
	nul := bytes.IndexByte(area, 0)
	if nul < 0 {
		return nil, fmt.Errorf("C string at +%#x has no NUL within row", offset)
	}
	return append([]byte(nil), area[:nul]...), nil
}
