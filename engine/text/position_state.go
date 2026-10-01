package text

import (
	"errors"
	"fmt"
)

// Paul2013PositionStateArrays holds the eight native state-array slots used
// around FUN_10022970. Slots 5 and 6 are nil because the recovered caller does
// not initialize or process them in this pass.
type Paul2013PositionStateArrays [8][]int32

// InitializePaul2013PositionStateArrays ports FUN_10022dc0's per-row initial
// writes before it calls FUN_10022970. Arrays 0, 1, and 2 receive the three
// caller-provided scalar defaults; arrays 3, 4, and 7 receive -1. The parser
// row count is bounded by the native 100-row parser capacity.
func InitializePaul2013PositionStateArrays(
	rowCount int,
	initial [3]int32,
) (Paul2013PositionStateArrays, error) {
	if rowCount < 0 || rowCount > paul2013ParserStateRowLimit {
		return Paul2013PositionStateArrays{}, fmt.Errorf(
			"Paul 2013 position-state row count %d is outside [0, %d]",
			rowCount,
			paul2013ParserStateRowLimit,
		)
	}
	var arrays Paul2013PositionStateArrays
	for mode := 0; mode <= 4; mode++ {
		arrays[mode] = make([]int32, rowCount)
	}
	arrays[7] = make([]int32, rowCount)
	for row := 0; row < rowCount; row++ {
		arrays[0][row] = initial[0]
		arrays[1][row] = initial[1]
		arrays[2][row] = initial[2]
		arrays[3][row] = -1
		arrays[4][row] = -1
		arrays[7][row] = -1
	}
	return arrays, nil
}

// Paul2013PositionStateRanges contains the inputs consumed by one of the
// sorted-boundary passes in FUN_10022970. The vendor fields remain opaque:
// row keys and boundaries are signed 32-bit values, Values supplies the
// value written between adjacent boundaries, and Initial supplies the
// caller-initialized destination. Modes 0, 1, and 2 clamp every assigned
// value; mode 7 preserves the -1 sentinel without clamping it.
type Paul2013PositionStateRanges struct {
	Mode       byte
	RowKeys    []int32
	Boundaries []int32
	Values     []int32
	Initial    []int32
	Minimum    int32
	Maximum    int32
}

// ApplyPaul2013PositionStateRanges ports the sorted-boundary sweep used for
// state arrays 0, 1, 2, and 7 in FUN_10022970. All inputs are explicit because
// the native caller's state tables and row-key producers are not yet
// reconstructed. Rows below the first boundary and rows at or beyond the last
// boundary retain their initialized values, as in the native loop.
func ApplyPaul2013PositionStateRanges(input Paul2013PositionStateRanges) ([]int32, error) {
	if input.Mode != 0 && input.Mode != 1 && input.Mode != 2 && input.Mode != 7 {
		return nil, fmt.Errorf("unsupported Paul 2013 position-state array mode %d", input.Mode)
	}
	if len(input.RowKeys) != len(input.Initial) {
		return nil, fmt.Errorf("received %d row keys for %d initial state values", len(input.RowKeys), len(input.Initial))
	}
	if len(input.Values)+1 < len(input.Boundaries) {
		return nil, fmt.Errorf("received %d interval values for %d boundaries", len(input.Values), len(input.Boundaries))
	}
	if input.Minimum > input.Maximum {
		return nil, errors.New("position-state minimum exceeds maximum")
	}
	if !sortedInt32(input.RowKeys) {
		return nil, errors.New("position-state row keys are not nondecreasing")
	}
	if !sortedInt32(input.Boundaries) {
		return nil, errors.New("position-state boundaries are not nondecreasing")
	}

	result := append([]int32(nil), input.Initial...)
	rowCursor := 0
	valueCursor := 0
	for boundaryIndex, boundary := range input.Boundaries {
		for rowCursor < len(input.RowKeys) && input.RowKeys[rowCursor] < boundary {
			rowCursor++
		}
		if boundaryIndex == 0 {
			valueCursor = rowCursor
			continue
		}
		for rowIndex := valueCursor; rowIndex < rowCursor; rowIndex++ {
			value := input.Values[boundaryIndex-1]
			if value < 0 {
				value = result[rowIndex]
			}
			if input.Mode != 7 || value != -1 {
				if value < input.Minimum {
					value = input.Minimum
				}
				if value > input.Maximum {
					value = input.Maximum
				}
			}
			result[rowIndex] = value
		}
		valueCursor = rowCursor
	}
	return result, nil
}

func sortedInt32(values []int32) bool {
	for index := 1; index < len(values); index++ {
		if values[index] < values[index-1] {
			return false
		}
	}
	return true
}

// Paul2013PositionStateProgram contains the inputs needed to execute the
// recovered FUN_10022970 state passes. Range and event producers remain
// caller-supplied; the native pass order and initial destination values are
// handled here.
type Paul2013PositionStateProgram struct {
	RowCount       int
	Initial        [3]int32
	Ranges         map[byte]Paul2013PositionStateRanges
	Events         map[byte]Paul2013PositionEventValues
	Terminal       *Paul2013PositionTerminalAccumulator
	ModelStateRows []byte
	IndexTables    *Paul2013PositionIndexTables
}

// Paul2013PositionStateProgramResult returns the mutated state arrays and,
// when supplied, the terminal selector-3 accumulator executed between event
// passes 3 and 4.
type Paul2013PositionStateProgramResult struct {
	Arrays             Paul2013PositionStateArrays
	Events             map[byte]Paul2013PositionEventValuesResult
	Terminal           Paul2013PositionTerminalAccumulatorResult
	HasTerminalPass    bool
	MappedIntervals    []Paul2013PositionEventRange
	HasMappedIntervals bool
}

// ApplyPaul2013PositionStateProgram composes FUN_10022dc0's state-array
// initialization with FUN_10022970's recovered pass order: range passes 0, 1,
// 2, and 7, then event pass 3, its optional terminal accumulator, and event
// pass 4, then maps the saved row indexes when the caller supplies the model
// record arena and lookup tables. Missing passes leave their initialized
// arrays unchanged. Inputs that produce parser rows, range tables, event
// tables, and terminal gates remain explicit.
func ApplyPaul2013PositionStateProgram(
	input Paul2013PositionStateProgram,
) (Paul2013PositionStateProgramResult, error) {
	arrays, err := InitializePaul2013PositionStateArrays(input.RowCount, input.Initial)
	if err != nil {
		return Paul2013PositionStateProgramResult{}, err
	}
	for mode, pass := range input.Ranges {
		if mode != 0 && mode != 1 && mode != 2 && mode != 7 {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("range pass supplied for unsupported state mode %d", mode)
		}
		if pass.Mode != mode {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("range pass key %d disagrees with embedded mode %d", mode, pass.Mode)
		}
	}
	for mode, pass := range input.Events {
		if mode != 3 && mode != 4 {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("event pass supplied for unsupported state mode %d", mode)
		}
		if pass.Mode != mode {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("event pass key %d disagrees with embedded mode %d", mode, pass.Mode)
		}
	}
	for _, mode := range [...]byte{0, 1, 2, 7} {
		pass, ok := input.Ranges[mode]
		if !ok {
			continue
		}
		pass.Initial = arrays[mode]
		arrays[mode], err = ApplyPaul2013PositionStateRanges(pass)
		if err != nil {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("apply state range pass %d: %w", mode, err)
		}
	}
	result := Paul2013PositionStateProgramResult{Arrays: arrays}
	if pass, ok := input.Events[3]; ok {
		pass.Initial = result.Arrays[3]
		eventResult, eventErr := ApplyPaul2013PositionEventValues(pass)
		if eventErr != nil {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("apply state event pass 3: %w", eventErr)
		}
		result.Arrays[3] = eventResult.Values
		result.Events = map[byte]Paul2013PositionEventValuesResult{3: eventResult}
		if input.Terminal != nil {
			terminal := *input.Terminal
			terminal.StartBoundary = eventResult.NextBoundary
			result.Terminal, err = EvaluatePaul2013PositionTerminalAccumulator(terminal)
			if err != nil {
				return Paul2013PositionStateProgramResult{}, fmt.Errorf("apply terminal state event pass: %w", err)
			}
			result.HasTerminalPass = true
		}
	} else if input.Terminal != nil {
		return Paul2013PositionStateProgramResult{}, errors.New("terminal accumulator requires event pass 3")
	}
	if pass, ok := input.Events[4]; ok {
		pass.Initial = result.Arrays[4]
		eventResult, eventErr := ApplyPaul2013PositionEventValues(pass)
		if eventErr != nil {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("apply state event pass 4: %w", eventErr)
		}
		result.Arrays[4] = eventResult.Values
		if result.Events == nil {
			result.Events = make(map[byte]Paul2013PositionEventValuesResult)
		}
		result.Events[4] = eventResult
	}
	if input.IndexTables != nil {
		result.MappedIntervals, err = MapPaul2013PositionStateRecordIndexes(
			input.ModelStateRows,
			input.RowCount,
			*input.IndexTables,
		)
		if err != nil {
			return Paul2013PositionStateProgramResult{}, fmt.Errorf("map position-state record indexes: %w", err)
		}
		result.HasMappedIntervals = true
	}
	return result, nil
}

// ApplyPaul2013PositionStateProgramWithProducerArrays connects the populated
// mode 0-2 arrays from FUN_1001d5d0 to the FUN_10022970 passes. Row keys,
// initial arrays, clamps, event tables, and terminal gates remain supplied by
// the caller.
func ApplyPaul2013PositionStateProgramWithProducerArrays(
	input Paul2013PositionStateProgram,
	raw [3]Paul2013PositionStateProducerSeries,
) (Paul2013PositionStateProgramResult, error) {
	normalized, err := NormalizePaul2013PositionStateProducerArrays(raw)
	if err != nil {
		return Paul2013PositionStateProgramResult{}, err
	}
	ranges := make(map[byte]Paul2013PositionStateRanges, len(input.Ranges))
	for mode, pass := range input.Ranges {
		ranges[mode] = pass
	}
	input.Ranges = ranges
	for mode, series := range normalized {
		pass, exists := input.Ranges[byte(mode)]
		if !exists {
			if len(series.Boundaries) != 0 {
				return Paul2013PositionStateProgramResult{}, fmt.Errorf(
					"normalized producer mode %d has %d boundaries but no range pass",
					mode,
					len(series.Boundaries),
				)
			}
			continue
		}
		pass.Boundaries = series.Boundaries
		pass.Values = series.Values
		input.Ranges[byte(mode)] = pass
	}
	return ApplyPaul2013PositionStateProgram(input)
}

// Paul2013RepeatedStateValueInput supplies the scalar values consumed by
// FUN_10026750. Their native structure offsets are known, but their semantic
// names are not; the fields preserve the observed branch roles.
type Paul2013RepeatedStateValueInput struct {
	Current       int32
	Primary       int32
	Secondary     int32
	ZeroFlag      bool
	EngineDefault int32
}

// ResolvePaul2013RepeatedStateValue ports FUN_10026750's negative-sentinel
// propagation and upper clamp. The native caller invokes it after
// FUN_10022dc0 for each repeated state row. Nonnegative Current values pass
// through; unresolved negative values select Primary, Secondary, zero, or the
// engine default according to the observed gates. Positive values above
// 65534 become 65535, while negative sentinels are preserved.
func ResolvePaul2013RepeatedStateValue(input Paul2013RepeatedStateValueInput) int32 {
	value := input.Current
	if value < 0 {
		switch {
		case input.Primary >= 0:
			if input.ZeroFlag {
				value = 0
			} else {
				value = input.Primary
			}
		case input.Secondary >= 0:
			value = input.Secondary
		case input.ZeroFlag:
			value = 0
		default:
			value = input.EngineDefault
		}
	}
	if value > 0xfffe {
		value = 0xffff
	}
	return value
}

// Paul2013PositionEventRange is the inclusive boundary interval searched for
// one row by selectors 3 and 4 in FUN_10022970.
type Paul2013PositionEventRange struct {
	Minimum int32
	Maximum int32
}

// Paul2013PositionEventValues contains the explicit inputs to the selector 3
// and 4 event scans. Selector 3 sums all matching values and clamps after
// every addition; selector 4 assigns each match so the last match wins. The
// terminal selector-3 accumulator at FUN_10022970's final utterance branch is
// separate and is not included.
type Paul2013PositionEventValues struct {
	Mode          byte
	Ranges        []Paul2013PositionEventRange
	Boundaries    []int32
	Values        []int32
	Initial       []int32
	StartBoundary int
	OverflowStart int
	Minimum       int32
	Maximum       int32
}

// Paul2013PositionEventValuesResult returns the modified row values and the
// first unconsumed boundary index, corresponding to the native scan cursor.
type Paul2013PositionEventValuesResult struct {
	Values       []int32
	NextBoundary int
	// OverflowBoundary retains the last upper-clamp index as a diagnostic.
	// It is not a second native cursor; terminal processing uses NextBoundary.
	OverflowBoundary int
}

// ApplyPaul2013PositionEventValues ports the per-row event scan for selectors
// 3 and 4 in FUN_10022970. Range endpoints and boundary values are explicit
// because their model-row producers remain unresolved. The persistent cursor
// advances only after a matching event; skipped boundaries are reconsidered
// by later rows. Upper-clamp matches advance the same native cursor.
// OverflowBoundary is retained only as a diagnostic for existing callers.
func ApplyPaul2013PositionEventValues(input Paul2013PositionEventValues) (Paul2013PositionEventValuesResult, error) {
	if input.Mode != 3 && input.Mode != 4 {
		return Paul2013PositionEventValuesResult{}, fmt.Errorf("unsupported Paul 2013 position-event selector %d", input.Mode)
	}
	if len(input.Ranges) != len(input.Initial) {
		return Paul2013PositionEventValuesResult{}, fmt.Errorf("received %d row ranges for %d initial values", len(input.Ranges), len(input.Initial))
	}
	if len(input.Boundaries) != len(input.Values) {
		return Paul2013PositionEventValuesResult{}, fmt.Errorf("received %d boundaries for %d event values", len(input.Boundaries), len(input.Values))
	}
	if input.StartBoundary < 0 || input.StartBoundary > len(input.Boundaries) {
		return Paul2013PositionEventValuesResult{}, fmt.Errorf("start boundary index %d outside [0, %d]", input.StartBoundary, len(input.Boundaries))
	}
	if input.OverflowStart < 0 || input.OverflowStart > len(input.Boundaries) {
		return Paul2013PositionEventValuesResult{}, fmt.Errorf("overflow boundary index %d outside [0, %d]", input.OverflowStart, len(input.Boundaries))
	}
	if input.Mode == 3 && input.Minimum > input.Maximum {
		return Paul2013PositionEventValuesResult{}, errors.New("position-event minimum exceeds maximum")
	}
	if !sortedInt32(input.Boundaries) {
		return Paul2013PositionEventValuesResult{}, errors.New("position-event boundaries are not nondecreasing")
	}
	// Overlapping source records can produce reversed inter-record bounds.
	// Native comparisons then match no event and retain the current cursor.

	result := Paul2013PositionEventValuesResult{
		Values:           append([]int32(nil), input.Initial...),
		NextBoundary:     input.StartBoundary,
		OverflowBoundary: input.OverflowStart,
	}
	for rowIndex, interval := range input.Ranges {
		matched := false
		scanBoundary := result.NextBoundary
		for scanBoundary < len(input.Boundaries) {
			boundaryIndex := scanBoundary
			boundary := input.Boundaries[boundaryIndex]
			if boundary < interval.Minimum || boundary > interval.Maximum {
				scanBoundary++
				continue
			}
			value := input.Values[boundaryIndex]
			if input.Mode == 3 && matched {
				value += result.Values[rowIndex]
			}
			overflowed := false
			if input.Mode == 3 {
				if value < input.Minimum {
					value = input.Minimum
				}
				if value > input.Maximum {
					value = input.Maximum
					overflowed = true
				}
			}
			result.Values[rowIndex] = value
			matched = true
			scanBoundary = boundaryIndex + 1
			result.NextBoundary = scanBoundary
			if overflowed {
				result.OverflowBoundary = scanBoundary
			}
		}
	}
	return result, nil
}

// Paul2013PositionTerminalAccumulator supplies the explicit inputs to the
// final selector-3 event pass in FUN_10022970.
type Paul2013PositionTerminalAccumulator struct {
	Enabled         bool
	MinimumBoundary int32
	MaximumBoundary int32
	Boundaries      []int32
	Values          []int32
	StartBoundary   int
	Minimum         int32
	Maximum         int32
}

// Paul2013PositionTerminalAccumulatorResult contains the final accumulator
// and its persistent matching-event cursor. Value is -1 when the native
// final-row gate is false or no event matches.
type Paul2013PositionTerminalAccumulatorResult struct {
	Value        int32
	NextBoundary int
}

// EvaluatePaul2013PositionTerminalAccumulator ports the final selector-3
// scan in FUN_10022970. The caller supplies the native end-row gate, terminal
// interval, sorted event arrays, and `+0x488fd` starting cursor.
func EvaluatePaul2013PositionTerminalAccumulator(
	input Paul2013PositionTerminalAccumulator,
) (Paul2013PositionTerminalAccumulatorResult, error) {
	result := Paul2013PositionTerminalAccumulatorResult{Value: -1, NextBoundary: input.StartBoundary}
	if len(input.Boundaries) != len(input.Values) {
		return result, fmt.Errorf("received %d terminal boundaries for %d event values", len(input.Boundaries), len(input.Values))
	}
	if input.StartBoundary < 0 || input.StartBoundary > len(input.Boundaries) {
		return result, fmt.Errorf("terminal start boundary index %d outside [0, %d]", input.StartBoundary, len(input.Boundaries))
	}
	// A reversed terminal interval matches no event under the native pair
	// of inclusive comparisons; it preserves the sentinel and cursor.
	if input.Minimum > input.Maximum {
		return result, errors.New("terminal accumulator minimum exceeds maximum")
	}
	if !sortedInt32(input.Boundaries) {
		return result, errors.New("terminal boundaries are not nondecreasing")
	}
	if !input.Enabled {
		return result, nil
	}

	matched := false
	scanBoundary := input.StartBoundary
	for scanBoundary < len(input.Boundaries) {
		boundaryIndex := scanBoundary
		boundary := input.Boundaries[boundaryIndex]
		scanBoundary++
		if boundary < input.MinimumBoundary || boundary > input.MaximumBoundary {
			continue
		}
		value := input.Values[boundaryIndex]
		if matched {
			value += result.Value
		}
		if value < input.Minimum {
			value = input.Minimum
		}
		if value > input.Maximum {
			value = input.Maximum
		}
		result.Value = value
		result.NextBoundary = scanBoundary
		matched = true
	}
	return result, nil
}
