package text

import (
	"encoding/binary"
	"fmt"
)

// Paul2013FinalizedModelState retains both stages of the context-table to
// model-record handoff. Finalized owns the surrounding parser state needed
// to resolve the record's parser-row pointers.
type Paul2013FinalizedModelState struct {
	Finalized Paul2013ModelParserFinalizerResult
	Records   Paul2013ModelStateRecordArena
}

// FinalizePaul2013ModelParserRowsAndBuildModelStateRecords connects
// FUN_1003e240's grouped output to FUN_10016c90. The parser-row view begins
// at ParserStateOffset itself: phone strings are at +0x66 and raw source
// offsets at +0x14/+0x18, with a 0x94 stride. This differs from the source-row
// slice accepted by the earlier token parser. The count and inclusive source
// intervals are copied into model records as in FUN_10022dc0 before any
// optional position-index remapping. Per-token state values and pointer
// address mapping remain caller inputs.
func FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(
	input Paul2013ModelParserFinalizerInput,
	modelStateArena []byte,
	basePosition int32,
	perTokenStateValues []int32,
	resolveAddress Paul2013ParserRowAddressResolver,
) (Paul2013FinalizedModelState, error) {
	finalized, err := FinalizePaul2013ModelParserRows(input)
	result := Paul2013FinalizedModelState{Finalized: finalized}
	if err != nil {
		return result, fmt.Errorf("finalize model parser rows: %w", err)
	}
	if finalized.ReturnValue != 1 {
		return result, fmt.Errorf("model parser finalizer returned %d for %d groups", finalized.ReturnValue, finalized.OutputGroupCount)
	}
	count := finalized.OutputGroupCount
	if count < 1 || count > paul2013ModelStateRecordCountLimit {
		return result, fmt.Errorf("finalized model parser group count %d outside [1, %d]", count, paul2013ModelStateRecordCountLimit)
	}
	start := input.ParserStateOffset
	needed := count * paul2013ParserSourceRowStride
	if needed > len(finalized.StateArena)-start {
		return result, fmt.Errorf("finalized parser state has %d bytes after +%#x, need %d for %d record-source rows", len(finalized.StateArena)-start, start, needed, count)
	}
	parserRows := finalized.StateArena[start : start+needed]
	offsets := make([]Paul2013PositionIntervalOffsets, count)
	for index := range offsets {
		row := parserRows[index*paul2013ParserSourceRowStride:]
		offsets[index] = Paul2013PositionIntervalOffsets{
			Start: int32(binary.LittleEndian.Uint32(row[0x14:0x18])),
			End:   int32(binary.LittleEndian.Uint32(row[0x18:0x1c])),
		}
	}
	intervals, err := BuildPaul2013PositionIntervals(basePosition, offsets)
	if err != nil {
		return result, fmt.Errorf("project finalized parser intervals: %w", err)
	}
	recordEnd := paul2013ModelStateRecordBaseOffset + count*paul2013ModelStateRecordStride
	if len(modelStateArena) < recordEnd {
		return result, fmt.Errorf("model-state arena has %d bytes, need %d for finalized records", len(modelStateArena), recordEnd)
	}
	arena := append([]byte(nil), modelStateArena...)
	binary.LittleEndian.PutUint16(arena[2:4], uint16(count))
	for index, interval := range intervals {
		record := arena[paul2013ModelStateRecordBaseOffset+index*paul2013ModelStateRecordStride:]
		binary.LittleEndian.PutUint32(record[0:4], uint32(interval.Minimum))
		binary.LittleEndian.PutUint32(record[4:8], uint32(interval.Maximum))
	}
	terminalPitchValueCount := int16(binary.LittleEndian.Uint16(parserRows[2:4]))
	result.Records, err = PopulatePaul2013ModelStateRecordsFromParserRows(
		arena, parserRows, perTokenStateValues, terminalPitchValueCount, resolveAddress,
	)
	if err != nil {
		return result, fmt.Errorf("populate finalized model-state records: %w", err)
	}
	return result, nil
}
