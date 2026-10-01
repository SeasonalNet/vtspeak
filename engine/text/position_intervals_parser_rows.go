package text

import "fmt"

// BuildPaul2013PositionIntervalsFromParserRows connects the ordinary parser
// row offsets to the inclusive model-state intervals written by
// FUN_10022dc0. Parser offsets are relative to the current source segment;
// segmentBase is the absolute position added by the native caller. It does
// not produce the range/event tables consumed later by FUN_10022970.
func BuildPaul2013PositionIntervalsFromParserRows(
	segmentBase int32,
	rows []Paul2013OrdinaryParserOffsetRow,
) ([]Paul2013PositionEventRange, error) {
	if len(rows) > paul2013ParserStateRowLimit {
		return nil, fmt.Errorf("position interval parser has %d rows, native limit is %d", len(rows), paul2013ParserStateRowLimit)
	}
	offsets := make([]Paul2013PositionIntervalOffsets, len(rows))
	for index, row := range rows {
		offsets[index] = Paul2013PositionIntervalOffsets{Start: row.Start, End: row.End}
	}
	intervals, err := BuildPaul2013PositionIntervals(segmentBase, offsets)
	if err != nil {
		return nil, fmt.Errorf("convert parser row offsets to inclusive position intervals: %w", err)
	}
	return intervals, nil
}
