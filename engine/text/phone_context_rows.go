package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	Paul2013PhoneContextTableHeaderSize = 4
	Paul2013PhoneContextRowSize         = 0x70
	Paul2013PhoneContextRowFlags        = 0x00
	Paul2013PhoneContextRowTokenIndex   = 0x02
	Paul2013PhoneContextRowStatus       = 0x06
	Paul2013PhoneContextRowSurface      = 0x07
	Paul2013PhoneContextRowPhoneCodes   = 0x25
	Paul2013PhoneContextRowSourceMarker = 0x66
	Paul2013PhoneContextRowMetadata     = 0x68

	paul2013PhoneContextParserRowStride = 0x94
	paul2013PhoneContextParserClassByte = 0x24

	paul2013PhoneContextSurfaceCapacity = Paul2013PhoneContextRowPhoneCodes - Paul2013PhoneContextRowSurface
	paul2013PhoneContextPhoneCapacity   = Paul2013PhoneContextRowSourceMarker - Paul2013PhoneContextRowPhoneCodes
)

// WritePaul2013PhoneContextRow projects the fields written by FUN_1000ea20
// for one already-selected token result. It does not run FUN_100091b0's later
// normalization cascade; those rule flags remain in the zero-based output
// row for a downstream caller to produce.
func WritePaul2013PhoneContextRow(
	destination []byte,
	tokenIndex uint16,
	surface []byte,
	status byte,
	phoneCodes []byte,
	sourceClass byte,
	metadata [4]bool,
) error {
	if len(destination) < Paul2013PhoneContextRowSize {
		return fmt.Errorf("phone-context row has %d bytes, need %d", len(destination), Paul2013PhoneContextRowSize)
	}
	if len(surface) >= paul2013PhoneContextSurfaceCapacity || containsNUL(surface) {
		return fmt.Errorf("phone-context surface must be NUL-free and at most %d bytes", paul2013PhoneContextSurfaceCapacity-1)
	}
	if len(phoneCodes) >= paul2013PhoneContextPhoneCapacity || containsNUL(phoneCodes) {
		return fmt.Errorf("phone-code string must be NUL-free and at most %d bytes", paul2013PhoneContextPhoneCapacity-1)
	}
	row := destination[:Paul2013PhoneContextRowSize]
	binary.LittleEndian.PutUint16(row[Paul2013PhoneContextRowTokenIndex:], tokenIndex)
	row[Paul2013PhoneContextRowStatus] = status
	copyNULTerminated(row[Paul2013PhoneContextRowSurface:], surface)
	copyNULTerminated(row[Paul2013PhoneContextRowPhoneCodes:], phoneCodes)
	if sourceClass == 'X' {
		row[Paul2013PhoneContextRowSourceMarker] = 'X'
	} else {
		row[Paul2013PhoneContextRowSourceMarker] = '0'
	}
	for index, value := range metadata {
		if value {
			binary.LittleEndian.PutUint16(row[Paul2013PhoneContextRowMetadata+index*2:], 1)
		} else {
			binary.LittleEndian.PutUint16(row[Paul2013PhoneContextRowMetadata+index*2:], 0)
		}
	}
	if len(phoneCodes) > 4 && phoneCodes[4] != 0 {
		row[Paul2013PhoneContextRowFlags] = 1
	} else {
		row[Paul2013PhoneContextRowFlags] = 0
	}
	return nil
}

// BuildPaul2013PhoneContextRows projects selected lexical token spans into
// FUN_1000ea20's counted 0x70-byte row array. tokenIndexes and sourceClasses
// are explicit parser inputs; source-class matching and the later
// FUN_100091b0/FUN_10007520 rule cascades are not inferred here.
func BuildPaul2013PhoneContextRows(
	sequence LexicalPhoneSequence,
	tokenIndexes []uint16,
	sourceClasses []byte,
) ([]byte, error) {
	if len(sequence.Tokens) != len(tokenIndexes) || len(sequence.Tokens) != len(sourceClasses) {
		return nil, fmt.Errorf("phone-context projection has %d token spans, %d token indexes, and %d source classes",
			len(sequence.Tokens), len(tokenIndexes), len(sourceClasses))
	}
	if len(sequence.Tokens) > int(^uint16(0)) {
		return nil, fmt.Errorf("phone-context row count %d exceeds u16 capacity", len(sequence.Tokens))
	}
	result := make([]byte, Paul2013PhoneContextTableHeaderSize+len(sequence.Tokens)*Paul2013PhoneContextRowSize)
	binary.LittleEndian.PutUint16(result, uint16(len(sequence.Tokens)))
	for index, token := range sequence.Tokens {
		var phoneCodes []byte
		if token.ModelPhoneRows.HasContextMarker {
			phoneCodes = token.ModelPhoneRows.SelectedPhone
		} else {
			if token.AlternativeIndex < 0 || token.AlternativeIndex >= len(token.ModelPhoneRows.PhoneStrings) {
				return nil, fmt.Errorf("token %d (%q) has no model phone row for alternative %d", index, token.SourceSurface, token.AlternativeIndex)
			}
			phoneCodes = token.ModelPhoneRows.PhoneStrings[token.AlternativeIndex]
		}
		rowStart := Paul2013PhoneContextTableHeaderSize + index*Paul2013PhoneContextRowSize
		if err := WritePaul2013PhoneContextRow(
			result[rowStart:rowStart+Paul2013PhoneContextRowSize], tokenIndexes[index], []byte(token.Surface),
			token.ModelPhoneRows.ResultType, phoneCodes, sourceClasses[index], token.DictionaryMetadata,
		); err != nil {
			return nil, fmt.Errorf("project phone-context row %d (%q): %w", index, token.SourceSurface, err)
		}
	}
	return result, nil
}

// BuildPaul2013PhoneContextRowsFromParserRows derives FUN_1000ea20's X retry
// marker from byte +0x24 of each selected 0x94-byte parser row, then builds
// the counted phone-context rows. tokenParserIndexes remain explicit because
// the mapping from resolved lexical spans to parser rows is not recovered.
func BuildPaul2013PhoneContextRowsFromParserRows(
	sequence LexicalPhoneSequence,
	parserRows []byte,
	tokenParserIndexes []uint16,
) ([]byte, error) {
	if len(sequence.Tokens) != len(tokenParserIndexes) {
		return nil, fmt.Errorf("phone-context projection has %d token spans and %d parser-row indexes",
			len(sequence.Tokens), len(tokenParserIndexes))
	}
	if len(parserRows)%paul2013PhoneContextParserRowStride != 0 {
		return nil, fmt.Errorf("parser-row buffer has %d bytes, not a multiple of 0x%x",
			len(parserRows), paul2013PhoneContextParserRowStride)
	}
	parserRowCount := len(parserRows) / paul2013PhoneContextParserRowStride
	sourceClasses := make([]byte, len(tokenParserIndexes))
	for index, tokenParserIndex := range tokenParserIndexes {
		if int(tokenParserIndex) >= parserRowCount {
			return nil, fmt.Errorf("token %d references parser row %d outside %d rows",
				index, tokenParserIndex, parserRowCount)
		}
		rowStart := int(tokenParserIndex) * paul2013PhoneContextParserRowStride
		if parserRows[rowStart+paul2013PhoneContextParserClassByte] == 'X' {
			sourceClasses[index] = 'X'
		}
	}
	return BuildPaul2013PhoneContextRows(sequence, tokenParserIndexes, sourceClasses)
}

// Paul2013PhoneContextPreparation records the low-short result returned by
// FUN_100091b0. The cascade's row mutations may remain unresolved while this
// value is known from the function's branch and tail-return structure.
type Paul2013PhoneContextPreparation struct {
	Code        int16
	Known       bool
	EarlyReturn bool
}

// Paul2013ExceptionRowInput contains the fields FUN_10008dc0 reads from one
// selected phone/context row. Row eligibility and ordering remain caller
// inputs; Surface and RetryHPrefix are copied from the row itself.
type Paul2013ExceptionRowInput struct {
	RowIndex     int
	Surface      string
	RetryHPrefix bool
}

// ExtractPaul2013ExceptionRowInputs reads surfaces at +0x07 and X retry
// markers at +0x66 from the caller-selected ordered row indexes. The native
// producer of the marker and the caller's eligibility/ordering policy remain
// outside this extractor.
func ExtractPaul2013ExceptionRowInputs(
	phoneContextTable []byte,
	rowIndexes []int,
) ([]Paul2013ExceptionRowInput, error) {
	if len(phoneContextTable) < Paul2013PhoneContextTableHeaderSize {
		return nil, fmt.Errorf("phone-context table has %d bytes, need at least %d",
			len(phoneContextTable), Paul2013PhoneContextTableHeaderSize)
	}
	rowCount := int(binary.LittleEndian.Uint16(phoneContextTable))
	rowsBytes := len(phoneContextTable) - Paul2013PhoneContextTableHeaderSize
	if rowCount > rowsBytes/Paul2013PhoneContextRowSize {
		return nil, fmt.Errorf("phone-context table declares %d rows but has space for %d",
			rowCount, rowsBytes/Paul2013PhoneContextRowSize)
	}
	result := make([]Paul2013ExceptionRowInput, len(rowIndexes))
	for index, rowIndex := range rowIndexes {
		if rowIndex < 0 || rowIndex >= rowCount {
			return nil, fmt.Errorf("exception row index %d at position %d is outside 0..%d",
				rowIndex, index, rowCount-1)
		}
		start := Paul2013PhoneContextTableHeaderSize + rowIndex*Paul2013PhoneContextRowSize
		row := phoneContextTable[start : start+Paul2013PhoneContextRowSize]
		surfaceArea := row[Paul2013PhoneContextRowSurface:Paul2013PhoneContextRowPhoneCodes]
		nul := bytes.IndexByte(surfaceArea, 0)
		if nul < 0 {
			return nil, fmt.Errorf("phone-context row %d surface is not NUL-terminated", rowIndex)
		}
		result[index] = Paul2013ExceptionRowInput{
			RowIndex: rowIndex, Surface: string(surfaceArea[:nul]),
			RetryHPrefix: row[Paul2013PhoneContextRowSourceMarker] == 'X',
		}
	}
	return result, nil
}

// DerivePaul2013PhoneContextPreparationFromRows ports every low-short return
// from FUN_100091b0. The source/parser table is a contiguous 0x94-byte row
// array; phoneContextTable has a four-byte header followed by counted
// 0x70-byte rows. The direct early return and non-Y/S tail paths return zero;
// the Y/S paths set local_c to one before entering the cascade and retain that
// low-short result regardless of which normalization rule matches.
func DerivePaul2013PhoneContextPreparationFromRows(
	sourceParserRows []byte,
	phoneContextTable []byte,
	phoneContextRowIndex int,
) (Paul2013PhoneContextPreparation, error) {
	const (
		sourceParserRowSize = 0x94
		sourceKindOffset    = 0x23
		sourceFormOffset    = 0x24
	)
	if len(phoneContextTable) < Paul2013PhoneContextTableHeaderSize {
		return Paul2013PhoneContextPreparation{}, fmt.Errorf("phone-context table has %d bytes, need at least %d", len(phoneContextTable), Paul2013PhoneContextTableHeaderSize)
	}
	rowCount := int(binary.LittleEndian.Uint16(phoneContextTable))
	if phoneContextRowIndex < 0 || phoneContextRowIndex >= rowCount {
		return Paul2013PhoneContextPreparation{}, fmt.Errorf("phone-context row index %d is outside 0..%d", phoneContextRowIndex, rowCount-1)
	}
	rowsBytes := len(phoneContextTable) - Paul2013PhoneContextTableHeaderSize
	if rowCount > rowsBytes/Paul2013PhoneContextRowSize {
		return Paul2013PhoneContextPreparation{}, fmt.Errorf("phone-context table declares %d rows but has space for %d", rowCount, rowsBytes/Paul2013PhoneContextRowSize)
	}
	rowStart := Paul2013PhoneContextTableHeaderSize + phoneContextRowIndex*Paul2013PhoneContextRowSize
	phoneRow := phoneContextTable[rowStart : rowStart+Paul2013PhoneContextRowSize]
	sourceTokenIndex := int16(binary.LittleEndian.Uint16(phoneRow[Paul2013PhoneContextRowTokenIndex:]))
	if sourceTokenIndex < 0 {
		return Paul2013PhoneContextPreparation{}, fmt.Errorf("phone-context row %d has negative source/parser index %d", phoneContextRowIndex, sourceTokenIndex)
	}
	sourceStart := int(sourceTokenIndex) * sourceParserRowSize
	if sourceStart > len(sourceParserRows) || len(sourceParserRows)-sourceStart < sourceParserRowSize {
		return Paul2013PhoneContextPreparation{}, fmt.Errorf("source/parser row %d is outside %d available bytes", sourceTokenIndex, len(sourceParserRows))
	}
	earlyReturn := sourceParserRows[sourceStart+sourceKindOffset] != 'U'
	if phoneContextRowIndex > 0 {
		previousStart := Paul2013PhoneContextTableHeaderSize + (phoneContextRowIndex-1)*Paul2013PhoneContextRowSize
		previousIndex := int16(binary.LittleEndian.Uint16(phoneContextTable[previousStart+Paul2013PhoneContextRowTokenIndex:]))
		earlyReturn = earlyReturn || previousIndex == sourceTokenIndex
	}
	if earlyReturn {
		return Paul2013PhoneContextPreparation{Code: 0, Known: true, EarlyReturn: true}, nil
	}
	if sourceForm := sourceParserRows[sourceStart+sourceFormOffset]; sourceForm == 'Y' || sourceForm == 'S' {
		return Paul2013PhoneContextPreparation{Code: 1, Known: true}, nil
	}
	return Paul2013PhoneContextPreparation{Code: 0, Known: true}, nil
}
