package text

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Paul2013TPPTokenUpdate describes the writes made by FUN_1000e0c0 for one
// token row. The native row's +0x20 flag is written only on rows preceding
// the final row; the typed suffix is written on every row. A class code is
// written only when the caller's WAB lookup found that token.
type Paul2013TPPTokenUpdate struct {
	ProcessingFlagWritten bool
	ProcessingFlag        uint16
	TypedCodeWritten      bool
	TypedCode             byte
	ClassCodeWritten      bool
	ClassCode             byte
}

// ApplyPaul2013TPPNumericAtom ports FUN_1000e0c0 after TPP lookup and numeric
// parsing. classIndexes contains one native WAB row index per matched token,
// with -1 for a token absent from that table. The row layout and WAB match
// policy remain caller-owned inputs.
func ApplyPaul2013TPPNumericAtom(atom TPPAtom, classIndexes []int) ([]Paul2013TPPTokenUpdate, bool, error) {
	if atom.Tag != 'F' && atom.Tag != 'G' {
		return nil, false, nil
	}
	if len(classIndexes) == 0 {
		return nil, false, fmt.Errorf("TPP %c numeric operation has no matched token rows", atom.Tag)
	}
	value, numeric, err := atom.NumericByte()
	if err != nil {
		return nil, false, fmt.Errorf("parse TPP %c numeric operation: %w", atom.Tag, err)
	}
	if !numeric {
		return nil, false, fmt.Errorf("TPP %c numeric operation did not produce a byte", atom.Tag)
	}
	updates, err := ApplyPaul2013TPPTokenRange(value, classIndexes)
	if err != nil {
		return nil, false, err
	}
	return updates, true, nil
}

// ApplyPaul2013TPPTokenRange ports FUN_1000e0c0's writes for an already
// narrowed byte value and caller-selected inclusive token range. Earlier rows
// receive processing flag 1; every row receives the value at +0x25, and a
// found WAB class writes its one-based code at +0x26. The selector-specific
// value conversion and row-range production are upstream inputs.
func ApplyPaul2013TPPTokenRange(value byte, classIndexes []int) ([]Paul2013TPPTokenUpdate, error) {
	if len(classIndexes) == 0 {
		return nil, errors.New("TPP token range has no rows")
	}
	updates := make([]Paul2013TPPTokenUpdate, len(classIndexes))
	for tokenIndex, classIndex := range classIndexes {
		if classIndex < -1 || classIndex >= 0xff {
			return nil, fmt.Errorf("WAB class index %d is outside the absent-or-byte-index range", classIndex)
		}
		update := Paul2013TPPTokenUpdate{TypedCodeWritten: true, TypedCode: value}
		if tokenIndex+1 < len(classIndexes) {
			update.ProcessingFlagWritten = true
			update.ProcessingFlag = 1
		}
		if classIndex >= 0 {
			update.ClassCodeWritten = true
			update.ClassCode = byte(uint8(classIndex + 1))
		}
		updates[tokenIndex] = update
	}
	return updates, nil
}

// Paul2013TPPNumericTextInput supplies the row text and WAB resource used by
// FUN_1000dfc0. LookupWABClass can override the native table lookup for
// controlled comparisons; otherwise WABTable and KeyTables provide it.
type Paul2013TPPNumericTextInput struct {
	Rows            []Paul2013TPPWindowRow
	Dictionary      *TPPDictionary
	WABTable        *Table
	KeyTables       EmbeddedKeyTables
	DictionaryIndex int
	LookupWABClass  func(text []byte) (classIndex int, found bool, err error)
}

// Paul2013TPPNumericTextOperation records a direct F/G code applied to one or
// more lexical rows. Row indexes are inclusive, as in FUN_1000e0c0's range.
type Paul2013TPPNumericTextOperation struct {
	Tag      byte
	StartRow int
	EndRow   int
}

// Paul2013TPPNumericTextResult contains row writes planned by the recovered
// direct F/G branches and the WAB-only fallback. Input rows are never mutated.
type Paul2013TPPNumericTextResult struct {
	Updates    []Paul2013TPPTokenUpdate
	Operations []Paul2013TPPNumericTextOperation
}

const (
	paul2013TPPParserRowStride         = 0x94
	paul2013TPPParserRowLimit          = 100
	paul2013TPPParserRowProcessingFlag = 0x20
	paul2013TPPParserRowTypedCode      = 0x25
	paul2013TPPParserRowClassCode      = 0x26
)

// Paul2013TPPWindowRowsFromParserArena projects the fields read by
// FUN_1000e160 and FUN_1000dfc0 from a caller-arena parser row. The native
// rows have a 0x94-byte stride: three signed dwords at +0x00/+0x04/+0x08,
// the state byte at +0x23, and a NUL-terminated surface at +0x34.
func Paul2013TPPWindowRowsFromParserArena(parserRows []byte) ([]Paul2013TPPWindowRow, error) {
	if len(parserRows)%paul2013TPPParserRowStride != 0 {
		return nil, fmt.Errorf("TPP parser rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013TPPParserRowStride)
	}
	rowCount := len(parserRows) / paul2013TPPParserRowStride
	if rowCount > paul2013TPPParserRowLimit {
		return nil, fmt.Errorf("TPP parser arena has %d rows, native limit is %d", rowCount, paul2013TPPParserRowLimit)
	}
	rows := make([]Paul2013TPPWindowRow, rowCount)
	for rowIndex := range rows {
		row := parserRows[rowIndex*paul2013TPPParserRowStride : (rowIndex+1)*paul2013TPPParserRowStride]
		textEnd := bytes.IndexByte(row[0x34:], 0)
		if textEnd < 0 {
			return nil, fmt.Errorf("TPP parser row %d surface at +0x34 is not NUL-terminated", rowIndex)
		}
		rows[rowIndex] = Paul2013TPPWindowRow{
			Word0: int32(binary.LittleEndian.Uint32(row[0x00:])),
			Word1: int32(binary.LittleEndian.Uint32(row[0x04:])),
			Word2: int32(binary.LittleEndian.Uint32(row[0x08:])),
			State: row[0x23],
			Text:  append([]byte(nil), row[0x34:0x34+textEnd]...),
		}
	}
	return rows, nil
}

// WritePaul2013TPPNumericTextUpdates applies the fields written by
// FUN_1000e0c0 to a copy of the caller-arena parser rows. Fields whose
// Written flag is false retain their existing bytes.
func WritePaul2013TPPNumericTextUpdates(
	parserRows []byte,
	updates []Paul2013TPPTokenUpdate,
) ([]byte, error) {
	if len(parserRows)%paul2013TPPParserRowStride != 0 {
		return nil, fmt.Errorf("TPP parser rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013TPPParserRowStride)
	}
	rowCount := len(parserRows) / paul2013TPPParserRowStride
	if rowCount > paul2013TPPParserRowLimit {
		return nil, fmt.Errorf("TPP parser arena has %d rows, native limit is %d", rowCount, paul2013TPPParserRowLimit)
	}
	if len(updates) > rowCount {
		return nil, fmt.Errorf("TPP update list has %d rows, parser arena has %d", len(updates), rowCount)
	}
	result := append([]byte(nil), parserRows...)
	for rowIndex, update := range updates {
		row := result[rowIndex*paul2013TPPParserRowStride : (rowIndex+1)*paul2013TPPParserRowStride]
		if update.ProcessingFlagWritten {
			binary.LittleEndian.PutUint16(row[paul2013TPPParserRowProcessingFlag:], update.ProcessingFlag)
		}
		if update.TypedCodeWritten {
			row[paul2013TPPParserRowTypedCode] = update.TypedCode
		}
		if update.ClassCodeWritten {
			row[paul2013TPPParserRowClassCode] = update.ClassCode
		}
	}
	return result, nil
}

// ApplyPaul2013TPPNumericTextRowsToParserArena composes the supported F/G and
// WAB routing in FUN_1000dfc0 with FUN_1000e0c0's exact parser-row writes.
// parserRows must contain one 0x94-byte row for every input row. The original
// arena is not mutated.
func ApplyPaul2013TPPNumericTextRowsToParserArena(
	input Paul2013TPPNumericTextInput,
	parserRows []byte,
) ([]byte, Paul2013TPPNumericTextResult, error) {
	if len(parserRows)%paul2013TPPParserRowStride != 0 {
		return nil, Paul2013TPPNumericTextResult{}, fmt.Errorf("TPP parser rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013TPPParserRowStride)
	}
	parserRowCount := len(parserRows) / paul2013TPPParserRowStride
	if parserRowCount != len(input.Rows) {
		return nil, Paul2013TPPNumericTextResult{}, fmt.Errorf("TPP input has %d rows but parser arena has %d", len(input.Rows), parserRowCount)
	}
	result, err := ApplyPaul2013TPPNumericTextRows(input)
	if err != nil {
		return nil, Paul2013TPPNumericTextResult{}, err
	}
	updatedRows, err := WritePaul2013TPPNumericTextUpdates(parserRows, result.Updates)
	if err != nil {
		return nil, Paul2013TPPNumericTextResult{}, err
	}
	return updatedRows, result, nil
}

// ApplyPaul2013TPPNumericTextRows ports the row routing in FUN_1000dfc0 for
// exact direct F/G atoms: it tries the F compound window, then a G lookup on
// the current row, then writes only the WAB class byte when neither code
// matches. Other A-G selector lookups are available from the dictionary, but
// their caller-specific row effects and selector production remain unresolved.
func ApplyPaul2013TPPNumericTextRows(input Paul2013TPPNumericTextInput) (Paul2013TPPNumericTextResult, error) {
	result := Paul2013TPPNumericTextResult{Updates: make([]Paul2013TPPTokenUpdate, len(input.Rows))}
	if len(input.Rows) == 0 {
		return result, nil
	}
	if len(input.Rows) > 100 {
		return result, fmt.Errorf("TPP numeric-text input has %d rows, exceeds native 100-row state", len(input.Rows))
	}
	if input.Dictionary == nil {
		return result, errors.New("TPP numeric-text input has no dictionary")
	}
	if input.DictionaryIndex != 0 {
		return result, fmt.Errorf("TPP dictionary index %d is not loaded", input.DictionaryIndex)
	}
	if input.LookupWABClass == nil && input.WABTable == nil {
		return result, errors.New("TPP numeric-text input has no WAB table or class lookup")
	}

	rowIndex := 0
	for rowIndex < len(input.Rows) {
		window, err := MatchPaul2013TPPWindowWithDictionary(Paul2013TPPWindowInput{
			Rows: input.Rows, StartRow: rowIndex, CharacterMap: input.KeyTables.CharacterMap,
			DictionaryIndex: input.DictionaryIndex, InitialOutputRowIndex: -1,
		}, input.Dictionary, input.KeyTables)
		if err != nil {
			return Paul2013TPPNumericTextResult{}, fmt.Errorf("scan F-code window at row %d: %w", rowIndex, err)
		}
		if window.Matched {
			if window.OutputRowIndex < rowIndex || window.OutputRowIndex >= len(input.Rows) {
				return Paul2013TPPNumericTextResult{}, fmt.Errorf("F-code window returned row %d outside [%d, %d)", window.OutputRowIndex, rowIndex, len(input.Rows))
			}
			classIndexes, err := lookupPaul2013WABClasses(input, rowIndex, window.OutputRowIndex)
			if err != nil {
				return Paul2013TPPNumericTextResult{}, err
			}
			updates, applied, err := ApplyPaul2013TPPNumericAtom(
				TPPAtom{Tag: 'F', Suffix: window.Replacement}, classIndexes,
			)
			if err != nil {
				return Paul2013TPPNumericTextResult{}, fmt.Errorf("apply F code at rows %d-%d: %w", rowIndex, window.OutputRowIndex, err)
			}
			if !applied {
				return Paul2013TPPNumericTextResult{}, errors.New("F-code update was not applied")
			}
			copy(result.Updates[rowIndex:window.OutputRowIndex+1], updates)
			result.Operations = append(result.Operations, Paul2013TPPNumericTextOperation{
				Tag: 'F', StartRow: rowIndex, EndRow: window.OutputRowIndex,
			})
			rowIndex = window.OutputRowIndex + 1
			continue
		}

		rowText := cString(input.Rows[rowIndex].Text)
		suffix, found, err := input.Dictionary.LookupSelectedText(rowText, 'G', input.KeyTables)
		if err != nil {
			return Paul2013TPPNumericTextResult{}, fmt.Errorf("lookup G code for row %d: %w", rowIndex, err)
		}
		if found {
			classIndexes, err := lookupPaul2013WABClasses(input, rowIndex, rowIndex)
			if err != nil {
				return Paul2013TPPNumericTextResult{}, err
			}
			updates, applied, err := ApplyPaul2013TPPNumericAtom(
				TPPAtom{Tag: 'G', Suffix: suffix}, classIndexes,
			)
			if err != nil {
				return Paul2013TPPNumericTextResult{}, fmt.Errorf("apply G code at row %d: %w", rowIndex, err)
			}
			if !applied {
				return Paul2013TPPNumericTextResult{}, errors.New("G-code update was not applied")
			}
			result.Updates[rowIndex] = updates[0]
			result.Operations = append(result.Operations, Paul2013TPPNumericTextOperation{
				Tag: 'G', StartRow: rowIndex, EndRow: rowIndex,
			})
			rowIndex++
			continue
		}

		classIndex, found, err := lookupPaul2013WABClass(input, rowText)
		if err != nil {
			return Paul2013TPPNumericTextResult{}, fmt.Errorf("lookup WAB class for row %d: %w", rowIndex, err)
		}
		if found {
			if classIndex < 0 || classIndex >= 0xff {
				return Paul2013TPPNumericTextResult{}, fmt.Errorf("WAB lookup for row %d returned class index %d outside byte-code range", rowIndex, classIndex)
			}
			result.Updates[rowIndex].ClassCodeWritten = true
			result.Updates[rowIndex].ClassCode = byte(uint8(classIndex + 1))
		}
		rowIndex++
	}
	return result, nil
}

func lookupPaul2013WABClasses(input Paul2013TPPNumericTextInput, first, last int) ([]int, error) {
	indexes := make([]int, last-first+1)
	for rowIndex := first; rowIndex <= last; rowIndex++ {
		classIndex, found, err := lookupPaul2013WABClass(input, cString(input.Rows[rowIndex].Text))
		if err != nil {
			return nil, fmt.Errorf("lookup WAB class for row %d: %w", rowIndex, err)
		}
		indexes[rowIndex-first] = -1
		if found {
			if classIndex < 0 || classIndex >= 0xff {
				return nil, fmt.Errorf("WAB lookup for row %d returned class index %d outside byte-code range", rowIndex, classIndex)
			}
			indexes[rowIndex-first] = classIndex
		}
	}
	return indexes, nil
}

func lookupPaul2013WABClass(input Paul2013TPPNumericTextInput, surface []byte) (int, bool, error) {
	if input.LookupWABClass != nil {
		return input.LookupWABClass(surface)
	}
	return input.WABTable.LookupPaul2013WABClass(surface, input.KeyTables.CharacterMap)
}
