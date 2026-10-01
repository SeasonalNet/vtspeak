package synthesis

import (
	"vtspeak/engine/text"
)

// ApplyPaul2013PositionStateProgramWithControls connects FUN_10022850's
// selected, clamped controls to FUN_10022dc0's initial position-state arrays
// before executing the known FUN_10022970 passes. Native offsets map pitch,
// speed, and volume to state arrays 0, 1, and 2 respectively. The range and
// event tables remain explicit inputs in program.
func ApplyPaul2013PositionStateProgramWithControls(
	program text.Paul2013PositionStateProgram,
	controls Paul2013ControlOverrides,
) (text.Paul2013PositionStateProgramResult, Controls, error) {
	effective := ResolvePaul2013ControlOverrides(controls)
	program.Initial = [3]int32{effective.Pitch, effective.Speed, effective.Volume}
	result, err := text.ApplyPaul2013PositionStateProgram(program)
	if err != nil {
		return text.Paul2013PositionStateProgramResult{}, effective, err
	}
	return result, effective, nil
}

// ApplyPaul2013PositionStateProgramFromParserStateWithControls composes the
// selected, clamped controls and raw parser-state intervals with the known
// position-state passes. The parser-state row count must match program; its
// range/event values and gates remain caller supplied.
func ApplyPaul2013PositionStateProgramFromParserStateWithControls(
	program text.Paul2013PositionStateProgram,
	basePosition int32,
	parserState []byte,
	controls Paul2013ControlOverrides,
) (text.Paul2013PositionStateProgramResult, Controls, error) {
	effective := ResolvePaul2013ControlOverrides(controls)
	program.Initial = [3]int32{effective.Pitch, effective.Speed, effective.Volume}
	result, err := text.ApplyPaul2013PositionStateProgramFromParserState(
		program, basePosition, parserState,
	)
	if err != nil {
		return text.Paul2013PositionStateProgramResult{}, effective, err
	}
	return result, effective, nil
}

// ApplyPaul2013PositionStateProgramsForOrdinarySourceWithControls composes
// bounded ordinary-source interval production, one state program per parsed
// segment, and the selected clamped controls. Program row counts must match
// the emitted source rows; native range/event tables and gates stay explicit.
func ApplyPaul2013PositionStateProgramsForOrdinarySourceWithControls(
	source string,
	basePosition int32,
	programs []text.Paul2013PositionStateProgram,
	controls Paul2013ControlOverrides,
) ([]text.Paul2013PositionStateProgramResult, Controls, error) {
	effective := ResolvePaul2013ControlOverrides(controls)
	controlledPrograms := append([]text.Paul2013PositionStateProgram(nil), programs...)
	for index := range controlledPrograms {
		controlledPrograms[index].Initial = [3]int32{effective.Pitch, effective.Speed, effective.Volume}
	}
	results, err := text.ApplyPaul2013PositionStateProgramsForOrdinarySource(
		source, basePosition, controlledPrograms,
	)
	if err != nil {
		return nil, effective, err
	}
	return results, effective, nil
}
