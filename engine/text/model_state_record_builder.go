package text

import (
	"encoding/binary"
	"fmt"
)

const (
	paul2013ModelStateRecordBaseOffset = 0x64c
	paul2013ModelStateRecordStride     = 0x3c0
	paul2013ModelStateRecordCountLimit = 100
	paul2013ModelStateFinalModeOffset  = 0x4770a
)

// Paul2013ParserRowAddressResolver maps one parser-row index/byte offset to
// the 32-bit address value stored in the native record. The caller owns the
// address space and the bytes reached by these pointers.
type Paul2013ParserRowAddressResolver func(rowIndex, rowOffset int) (uint32, error)

// Paul2013ModelStateRecordArena retains the portable counted arena and the
// code-string projection used to populate each native record.
type Paul2013ModelStateRecordArena struct {
	Bytes        []byte
	ContextCodes []Paul2013ContextCodes
}

// PopulatePaul2013ModelStateRecordsFromParserRows ports the record writes in
// FUN_10016c90 for an already-produced 0x94-stride parser-state view. Row zero
// begins at the native parser-state base, with phone codes at +0x66; it is not
// the earlier source-row slice beginning after the parser header. It writes
// the recovered phone-code string, d/c phone markers, M flags, row control
// bytes, and the three pointer fields. Per-token state values remain explicit
// because their producer is upstream of this function. The final mode byte is
// derived from the last projected row and supplied terminal pitch-value count.
// Bytes not written by the native routine are preserved from modelStateArena.
func PopulatePaul2013ModelStateRecordsFromParserRows(
	modelStateArena []byte,
	parserRows []byte,
	perTokenStateValues []int32,
	terminalPitchValueCount int16,
	resolveAddress Paul2013ParserRowAddressResolver,
) (Paul2013ModelStateRecordArena, error) {
	if len(parserRows)%paul2013ParserSourceRowStride != 0 {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("parser rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013ParserSourceRowStride)
	}
	rowCount := len(parserRows) / paul2013ParserSourceRowStride
	if rowCount > paul2013ModelStateRecordCountLimit {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("model-state record count %d exceeds %d", rowCount, paul2013ModelStateRecordCountLimit)
	}
	if len(perTokenStateValues) != rowCount {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("received %d per-token state values for %d parser rows", len(perTokenStateValues), rowCount)
	}
	if len(modelStateArena) < 4 {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("model-state arena has %d bytes, need its record count at +2", len(modelStateArena))
	}
	declaredCount := int(int16(binary.LittleEndian.Uint16(modelStateArena[2:])))
	if declaredCount != rowCount {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("model-state arena declares %d records for %d parser rows", declaredCount, rowCount)
	}
	needed := paul2013ModelStateRecordBaseOffset + rowCount*paul2013ModelStateRecordStride
	if len(modelStateArena) < needed {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("model-state arena has %d bytes, need %d for %d records", len(modelStateArena), needed, rowCount)
	}
	if rowCount > 0 && len(modelStateArena) <= paul2013ModelStateFinalModeOffset {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("model-state arena has %d bytes, need byte +%#x for final mode", len(modelStateArena), paul2013ModelStateFinalModeOffset)
	}
	if resolveAddress == nil && rowCount != 0 {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("model-state record builder has no parser-row address resolver")
	}

	const (
		phoneCountOffset   = 0x95
		mFlagsOffset       = 0x29e
		phoneCodesOffset   = 0x2e8
		phoneMarkersOffset = 0x328
		phoneFollowOffset  = 0x329
		stateFlagOffset    = 0x2df
	)
	arena := append([]byte(nil), modelStateArena...)
	result := Paul2013ModelStateRecordArena{
		Bytes:        arena,
		ContextCodes: make([]Paul2013ContextCodes, rowCount),
	}
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		rowStart := rowIndex * paul2013ParserSourceRowStride
		row := parserRows[rowStart : rowStart+paul2013ParserSourceRowStride]
		codeArea := row[0x66:]
		codeEnd := 0
		for codeEnd < len(codeArea) && codeArea[codeEnd] != 0 {
			codeEnd++
		}
		if codeEnd == len(codeArea) {
			return Paul2013ModelStateRecordArena{}, fmt.Errorf("parser row %d context-code string at +0x66 is not NUL-terminated", rowIndex)
		}
		codes, err := ParsePaul2013ContextCodes(codeArea[:codeEnd])
		if err != nil {
			return Paul2013ModelStateRecordArena{}, fmt.Errorf("parse model-state context codes for parser row %d: %w", rowIndex, err)
		}
		if len(codes.Codes) > 0xff {
			return Paul2013ModelStateRecordArena{}, fmt.Errorf("parser row %d produced %d phone codes, exceeding the native byte count", rowIndex, len(codes.Codes))
		}

		recordStart := paul2013ModelStateRecordBaseOffset + rowIndex*paul2013ModelStateRecordStride
		record := arena[recordStart : recordStart+paul2013ModelStateRecordStride]
		for _, field := range [...]struct {
			recordOffset int
			rowOffset    int
		}{{0x2e0, 0x48}, {0x2e4, 0x66}, {0x3b8, 0x44}} {
			address, err := resolveAddress(rowIndex, field.rowOffset)
			if err != nil {
				return Paul2013ModelStateRecordArena{}, fmt.Errorf("resolve parser row %d +%#x pointer: %w", rowIndex, field.rowOffset, err)
			}
			binary.LittleEndian.PutUint32(record[field.recordOffset:], address)
		}
		copy(record[0x3ac:0x3ae], row[0x40:0x42])
		copy(record[0x3ae:0x3b0], row[0x34:0x36])
		record[0x3b0] = row[0x3c]
		record[0x3b1] = row[0x37]
		record[0x3b2] = row[0x39]
		record[0x3b3] = row[0x3a]
		copy(record[0x3b4:0x3b6], row[0x24:0x26])
		record[phoneCountOffset] = byte(len(codes.Codes))
		for index, flagged := range codes.MFlags {
			record[mFlagsOffset+index] = 0
			if flagged {
				record[mFlagsOffset+index] = 1
			}
		}
		copy(record[phoneCodesOffset:], codes.Codes)
		for index, marker := range codes.PhoneMarkers {
			record[phoneFollowOffset+index] = '0'
			if marker == '1' || marker == '2' {
				record[phoneMarkersOffset+index] = marker
			}
		}
		if perTokenStateValues[rowIndex] != -1 {
			record[stateFlagOffset] = 0x0c
		} else {
			record[stateFlagOffset] = 0
		}
		result.ContextCodes[rowIndex] = codes
	}
	if rowCount > 0 {
		lastRecordStart := paul2013ModelStateRecordBaseOffset + (rowCount-1)*paul2013ModelStateRecordStride
		lastRecord := arena[lastRecordStart : lastRecordStart+paul2013ModelStateRecordStride]
		finalTokenStateCode := int16(binary.LittleEndian.Uint16(lastRecord[0x3ac:]))
		summary, err := SummarizePaul2013ContextCodeState(
			perTokenStateValues,
			finalTokenStateCode,
			lastRecord[0x2df],
			terminalPitchValueCount,
		)
		if err != nil {
			return Paul2013ModelStateRecordArena{}, fmt.Errorf("summarize model-state context code state: %w", err)
		}
		arena[paul2013ModelStateFinalModeOffset] = summary.ModeCode
	}
	return result, nil
}
