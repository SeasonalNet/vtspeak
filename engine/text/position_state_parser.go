package text

import "fmt"

// ApplyPaul2013PositionStateProgramFromParserState connects the raw parser
// state consumed by FUN_10022dc0 to the recovered position-state passes. The
// parser state supplies row count and event intervals; range tables, event
// values, and native gates remain explicit in program.
func ApplyPaul2013PositionStateProgramFromParserState(
	program Paul2013PositionStateProgram,
	basePosition int32,
	parserState []byte,
) (Paul2013PositionStateProgramResult, error) {
	intervals, err := BuildPaul2013PositionIntervalsFromParserState(basePosition, parserState)
	if err != nil {
		return Paul2013PositionStateProgramResult{}, fmt.Errorf("read position-state parser intervals: %w", err)
	}
	if program.RowCount != len(intervals) {
		return Paul2013PositionStateProgramResult{}, fmt.Errorf(
			"position-state program has %d rows for %d parser intervals",
			program.RowCount,
			len(intervals),
		)
	}
	return applyPaul2013PositionStateProgramToIntervals(program, intervals)
}

// ApplyPaul2013PositionStateProgramsForOrdinarySource composes the bounded
// ordinary parser's row offsets with position-state event processing. Each
// program corresponds to one source segment split by a captured .?! boundary;
// program rows must match that segment's emitted parser rows. Unsupported
// source grammar and unresolved range/event values remain explicit errors or
// caller inputs.
func ApplyPaul2013PositionStateProgramsForOrdinarySource(
	source string,
	basePosition int32,
	programs []Paul2013PositionStateProgram,
) ([]Paul2013PositionStateProgramResult, error) {
	intervalsBySegment, err := BuildPaul2013OrdinaryParserPositionIntervals(source, basePosition)
	if err != nil {
		return nil, fmt.Errorf("build ordinary-source position intervals: %w", err)
	}
	if len(programs) != len(intervalsBySegment) {
		return nil, fmt.Errorf("received %d position-state programs for %d ordinary-source segments", len(programs), len(intervalsBySegment))
	}
	results := make([]Paul2013PositionStateProgramResult, len(programs))
	for segmentIndex, program := range programs {
		if program.RowCount != len(intervalsBySegment[segmentIndex]) {
			return nil, fmt.Errorf(
				"ordinary-source segment %d program has %d rows for %d parser intervals",
				segmentIndex,
				program.RowCount,
				len(intervalsBySegment[segmentIndex]),
			)
		}
		results[segmentIndex], err = applyPaul2013PositionStateProgramToIntervals(
			program, intervalsBySegment[segmentIndex],
		)
		if err != nil {
			return nil, fmt.Errorf("apply ordinary-source segment %d: %w", segmentIndex, err)
		}
	}
	return results, nil
}

func applyPaul2013PositionStateProgramToIntervals(
	program Paul2013PositionStateProgram,
	intervals []Paul2013PositionEventRange,
) (Paul2013PositionStateProgramResult, error) {
	if len(program.Events) > 0 {
		events := make(map[byte]Paul2013PositionEventValues, len(program.Events))
		for mode, pass := range program.Events {
			pass.Ranges = append([]Paul2013PositionEventRange(nil), intervals...)
			events[mode] = pass
		}
		program.Events = events
	}
	result, err := ApplyPaul2013PositionStateProgram(program)
	if err != nil {
		return Paul2013PositionStateProgramResult{}, err
	}
	return result, nil
}
