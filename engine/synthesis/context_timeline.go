package synthesis

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
)

// Paul2013ContextTimelineResult describes one FUN_1002d230 row operation.
// Applied is false when the native context gate suppresses both sides; in
// that case output is left untouched and the ordinary current-row path applies.
type Paul2013ContextTimelineResult struct {
	Applied bool
	Plan    Paul2013ContextEdgePlan
}

// Paul2013ContextTimelineRowInput supplies the already-produced native inputs
// for one primary row and its independently selected side contexts.
type Paul2013ContextTimelineRowInput struct {
	Current Paul2013TimelinePCMInput
	Gate    ContextGateInput
	Left    *Paul2013TimelinePCMInput
	Right   *Paul2013TimelinePCMInput
}

// Paul2013NormalContextTimelineRowInput keeps context selection alongside
// the normal timeline rows produced by the row builder.
type Paul2013NormalContextTimelineRowInput struct {
	Current               Paul2013NormalTimelineRow
	Gate                  ContextGateInput
	Left                  *Paul2013NormalTimelineRow
	Right                 *Paul2013NormalTimelineRow
	CurrentScratchSamples []int16
	LeftScratchSamples    []int16
	RightScratchSamples   []int16
}

// BuildPaul2013NormalContextTimelineRowInput selects the final left and right
// timeline row positions from the base indexes consumed by FUN_1002d230.
// Those base indexes, the gate rows, and the timeline rows remain explicit
// caller inputs; this ports only the observed mode-dependent +1/-1 adjustment.
func BuildPaul2013NormalContextTimelineRowInput(
	rows []Paul2013NormalTimelineRow,
	gate ContextGateInput,
	scratchSamples [][]int16,
) (Paul2013NormalContextTimelineRowInput, error) {
	if len(rows) == 0 || len(gate.Rows) != len(rows) || len(gate.Modes) != len(rows) {
		return Paul2013NormalContextTimelineRowInput{}, errors.New("context timeline rows, gate rows, and modes must have the same nonzero length")
	}
	if len(scratchSamples) != 0 && len(scratchSamples) != len(rows) {
		return Paul2013NormalContextTimelineRowInput{}, errors.New("context timeline scratch rows must be omitted or match the timeline row count")
	}
	multipliers, err := Paul2013ContextMultipliers(gate)
	if err != nil {
		return Paul2013NormalContextTimelineRowInput{}, err
	}
	input := Paul2013NormalContextTimelineRowInput{Current: rows[gate.CurrentIndex], Gate: gate}
	if len(scratchSamples) != 0 {
		input.CurrentScratchSamples = scratchSamples[gate.CurrentIndex]
	}
	mode := gate.Modes[gate.CurrentIndex]
	if multipliers.Left != 0 {
		leftIndex := gate.LeftIndex
		if mode != 2 {
			leftIndex++
		}
		if leftIndex < 0 || leftIndex >= len(rows) {
			return Paul2013NormalContextTimelineRowInput{}, fmt.Errorf("adjusted left context row index %d is outside the timeline row range", leftIndex)
		}
		input.Left = &rows[leftIndex]
		if len(scratchSamples) != 0 {
			input.LeftScratchSamples = scratchSamples[leftIndex]
		}
	}
	if multipliers.Right != 0 {
		rightIndex := gate.RightIndex
		if mode != 1 {
			rightIndex--
		}
		if rightIndex < 0 || rightIndex >= len(rows) {
			return Paul2013NormalContextTimelineRowInput{}, fmt.Errorf("adjusted right context row index %d is outside the timeline row range", rightIndex)
		}
		input.Right = &rows[rightIndex]
		if len(scratchSamples) != 0 {
			input.RightScratchSamples = scratchSamples[rightIndex]
		}
	}
	return input, nil
}

// BuildPaul2013NormalContextTimelineInputs applies the recovered side-index
// adjustments to each supplied row gate. Gates remain ordered as output rows;
// their model rows and base indexes are still produced by the caller.
func BuildPaul2013NormalContextTimelineInputs(
	rows []Paul2013NormalTimelineRow,
	gates []ContextGateInput,
	scratchSamples [][]int16,
) ([]Paul2013NormalContextTimelineRowInput, error) {
	if len(gates) == 0 {
		return nil, errors.New("normal context timeline requires at least one gate")
	}
	inputs := make([]Paul2013NormalContextTimelineRowInput, len(gates))
	for index, gate := range gates {
		input, err := BuildPaul2013NormalContextTimelineRowInput(rows, gate, scratchSamples)
		if err != nil {
			return nil, fmt.Errorf("normal context timeline gate %d: %w", index, err)
		}
		inputs[index] = input
	}
	return inputs, nil
}

// BuildPaul2013ContextTimelineRowInputs maps typed normal timeline rows to
// the PCM inputs consumed by the context renderer without losing the supplied
// scratch tails. Every row must be a normal row; synthetic boundary rows are
// handled by the separate timeline output path.
func BuildPaul2013ContextTimelineRowInputs(
	rows []Paul2013NormalContextTimelineRowInput,
) ([]Paul2013ContextTimelineRowInput, error) {
	if len(rows) == 0 {
		return nil, errors.New("context timeline row input builder has no rows")
	}
	inputs := make([]Paul2013ContextTimelineRowInput, len(rows))
	for index, row := range rows {
		current, err := buildPaul2013ContextTimelinePCMInput(row.Current, row.CurrentScratchSamples)
		if err != nil {
			return nil, fmt.Errorf("context timeline row %d current: %w", index, err)
		}
		inputs[index] = Paul2013ContextTimelineRowInput{Current: current, Gate: row.Gate}
		if row.Left != nil {
			left, err := buildPaul2013ContextTimelinePCMInput(*row.Left, row.LeftScratchSamples)
			if err != nil {
				return nil, fmt.Errorf("context timeline row %d left: %w", index, err)
			}
			inputs[index].Left = &left
		}
		if row.Right != nil {
			right, err := buildPaul2013ContextTimelinePCMInput(*row.Right, row.RightScratchSamples)
			if err != nil {
				return nil, fmt.Errorf("context timeline row %d right: %w", index, err)
			}
			inputs[index].Right = &right
		}
	}
	return inputs, nil
}

func buildPaul2013ContextTimelinePCMInput(
	row Paul2013NormalTimelineRow,
	scratch []int16,
) (Paul2013TimelinePCMInput, error) {
	inputs, err := BuildPaul2013TimelinePCMInputs([]Paul2013NormalTimelineRow{row})
	if err != nil {
		return Paul2013TimelinePCMInput{}, err
	}
	inputs[0].SourceScratchSamples = scratch
	return inputs[0], nil
}

// RenderPaul2013ContextTimelineRows applies per-row context reconstruction and
// then the direct default-pitch timeline join. All primary rows, context rows,
// view modes, gate indexes, gains, and any scratch tails must already be
// supplied; the text and neighbor selectors are outside this operation.
func RenderPaul2013ContextTimelineRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013ContextTimelineRowInput,
) (PCM, error) {
	if units == nil {
		return PCM{}, errors.New("context timeline renderer has no voice model")
	}
	if len(rows) == 0 {
		return PCM{}, errors.New("context timeline renderer requires at least one row")
	}
	prepared := make([]TimelinePCMRow, len(rows))
	for index, input := range rows {
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		current, err := readPaul2013ContextTimelineSource(ctx, units, input.Current)
		if err != nil {
			return PCM{}, fmt.Errorf("prepare current timeline row %d: %w", index, err)
		}
		// FUN_1002d230 clears its destination accumulator before adding the
		// current and context contributions.
		samples := make([]int16, len(current.samples))
		mixResult, err := mixPaul2013ContextTimelineSource(
			ctx, units, samples, 0, input.Gate, current, input.Left, input.Right,
		)
		if err != nil {
			return PCM{}, fmt.Errorf("reconstruct context for timeline row %d: %w", index, err)
		}
		if mixResult.Applied {
			current.samples = samples
		}
		prepared[index] = TimelinePCMRow{
			Samples: current.samples, LeadingSpan: input.Current.LeadingSpan,
			TrailingSpan: input.Current.TrailingSpan,
		}
	}
	samples, err := MixPaul2013Timeline(prepared)
	if err != nil {
		return PCM{}, fmt.Errorf("join context-prepared timeline rows: %w", err)
	}
	return PCM{SampleRate: 16000, Samples: samples}, nil
}

// RenderPaul2013NormalContextTimelineRows carries already-built normal row
// metadata through context reconstruction and the direct timeline join. It
// does not produce rows or choose neighboring units from source text.
func RenderPaul2013NormalContextTimelineRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013NormalContextTimelineRowInput,
) (PCM, error) {
	inputs, err := BuildPaul2013ContextTimelineRowInputs(rows)
	if err != nil {
		return PCM{}, err
	}
	return RenderPaul2013ContextTimelineRows(ctx, units, inputs)
}

// RenderPaul2013NormalContextTimelineFromRows derives side row positions from
// explicit gate inputs, then composes typed-row preparation and rendering.
func RenderPaul2013NormalContextTimelineFromRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013NormalTimelineRow,
	gates []ContextGateInput,
	scratchSamples [][]int16,
) (PCM, error) {
	inputs, err := BuildPaul2013NormalContextTimelineInputs(rows, gates, scratchSamples)
	if err != nil {
		return PCM{}, err
	}
	return RenderPaul2013NormalContextTimelineRows(ctx, units, inputs)
}

// MixPaul2013ContextTimelineRows loads the already-selected current and
// neighboring timeline rows, aligns their decoded samples by UPM period, and
// applies the Stage 8 context blend. The caller supplies the native gate rows,
// selected neighbor references, row modes, spans, and gains. It does not run
// FUN_1001b200 or infer those inputs from text.
func MixPaul2013ContextTimelineRows(
	ctx context.Context,
	units UnitReader,
	output []int16,
	outputOffset int,
	gate ContextGateInput,
	current Paul2013TimelinePCMInput,
	left *Paul2013TimelinePCMInput,
	right *Paul2013TimelinePCMInput,
) (Paul2013ContextTimelineResult, error) {
	if units == nil {
		return Paul2013ContextTimelineResult{}, errors.New("context timeline mixer has no voice model")
	}
	currentSource, err := readPaul2013ContextTimelineSource(ctx, units, current)
	if err != nil {
		return Paul2013ContextTimelineResult{}, fmt.Errorf("read current context timeline row: %w", err)
	}
	return mixPaul2013ContextTimelineSource(ctx, units, output, outputOffset, gate, currentSource, left, right)
}

func mixPaul2013ContextTimelineSource(
	ctx context.Context,
	units UnitReader,
	output []int16,
	outputOffset int,
	gate ContextGateInput,
	currentSource paul2013ContextTimelineSource,
	left *Paul2013TimelinePCMInput,
	right *Paul2013TimelinePCMInput,
) (Paul2013ContextTimelineResult, error) {
	multipliers, err := Paul2013ContextMultipliers(gate)
	if err != nil {
		return Paul2013ContextTimelineResult{}, err
	}
	if !multipliers.HasEligibleSide() {
		return Paul2013ContextTimelineResult{}, nil
	}
	var leftSource, rightSource *paul2013ContextTimelineSource
	if multipliers.Left != 0 {
		if left == nil {
			return Paul2013ContextTimelineResult{}, errors.New("left context is eligible but no selected left row was supplied")
		}
		loaded, err := readPaul2013ContextTimelineSource(ctx, units, *left)
		if err != nil {
			return Paul2013ContextTimelineResult{}, fmt.Errorf("read selected left context row: %w", err)
		}
		leftSource = &loaded
	}
	if multipliers.Right != 0 {
		if right == nil {
			return Paul2013ContextTimelineResult{}, errors.New("right context is eligible but no selected right row was supplied")
		}
		loaded, err := readPaul2013ContextTimelineSource(ctx, units, *right)
		if err != nil {
			return Paul2013ContextTimelineResult{}, fmt.Errorf("read selected right context row: %w", err)
		}
		rightSource = &loaded
	}
	plan, err := BuildPaul2013ContextEdgePlan(
		gate, len(currentSource.upm), contextUPMCount(leftSource), contextUPMCount(rightSource),
	)
	if err != nil {
		return Paul2013ContextTimelineResult{}, err
	}
	periods, err := buildPaul2013ContextUPMPeriods(&currentSource, leftSource, rightSource, plan)
	if err != nil {
		return Paul2013ContextTimelineResult{}, err
	}
	if err := MixPaul2013UPMTimeline(output, outputOffset, plan.Weights, periods); err != nil {
		return Paul2013ContextTimelineResult{}, fmt.Errorf("mix selected context rows: %w", err)
	}
	return Paul2013ContextTimelineResult{Applied: true, Plan: plan}, nil
}

type paul2013ContextTimelineSource struct {
	samples []int16
	upm     []byte
}

func readPaul2013ContextTimelineSource(
	ctx context.Context,
	units UnitReader,
	row Paul2013TimelinePCMInput,
) (paul2013ContextTimelineSource, error) {
	if row.RowType != 2 {
		return paul2013ContextTimelineSource{}, fmt.Errorf("unsupported timeline row kind %d", row.RowType)
	}
	if err := ctx.Err(); err != nil {
		return paul2013ContextTimelineSource{}, err
	}
	unit, err := units.ReadUnit(row.Unit.Bank, row.Unit.Index)
	if err != nil {
		return paul2013ContextTimelineSource{}, fmt.Errorf("read %s:%d: %w", row.Unit.Bank, row.Unit.Index, err)
	}
	if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
		return paul2013ContextTimelineSource{}, fmt.Errorf("%s:%d has invalid PCM byte length %d", row.Unit.Bank, row.Unit.Index, len(unit.PCM))
	}
	view, err := BuildPaul2013TimelineUnitView(unit, row.Unit, 0, row.Mode)
	if err != nil {
		return paul2013ContextTimelineSource{}, fmt.Errorf("build selected unit view: %w", err)
	}
	if row.SampleCount != int(view.SampleCount) ||
		row.LeadingSpan != int(view.LeadingPeriod)*2 ||
		row.TrailingSpan != int(view.TrailingPeriod)*2 {
		return paul2013ContextTimelineSource{}, fmt.Errorf(
			"%s:%d timeline row metadata count/spans %d/%d/%d disagree with mode-%d unit view %d/%d/%d",
			row.Unit.Bank, row.Unit.Index, row.SampleCount, row.LeadingSpan, row.TrailingSpan,
			row.Mode, view.SampleCount, int(view.LeadingPeriod)*2, int(view.TrailingPeriod)*2,
		)
	}
	decodedCount := len(unit.PCM) / 2
	if len(row.SourceScratchSamples) > int(^uint(0)>>1)-decodedCount {
		return paul2013ContextTimelineSource{}, errors.New("row scratch tail exceeds platform integer range")
	}
	source := make([]int16, decodedCount, decodedCount+len(row.SourceScratchSamples))
	for index := range source {
		source[index] = int16(binary.LittleEndian.Uint16(unit.PCM[index*2:]))
	}
	source = append(source, row.SourceScratchSamples...)
	prepared, err := BuildPaul2013TimelinePCMRow(
		source, decodedCount, row.Mode, row.SampleCount, row.LeadingSpan, row.TrailingSpan,
	)
	if err != nil {
		return paul2013ContextTimelineSource{}, err
	}
	prepared, err = ScalePaul2013TimelinePCMRow(prepared, row.GainPercent)
	if err != nil {
		return paul2013ContextTimelineSource{}, err
	}
	first, second, err := unit.Record.UPMSides(unit.UPM)
	if err != nil {
		return paul2013ContextTimelineSource{}, fmt.Errorf("validate %s:%d UPM sides: %w", row.Unit.Bank, row.Unit.Index, err)
	}
	upm := unit.UPM
	switch row.Mode {
	case Paul2013TimelineFirstSideView:
		upm = first
	case Paul2013TimelineSecondSideView:
		upm = second
	}
	if len(upm) == 0 || len(upm) > 255 {
		return paul2013ContextTimelineSource{}, fmt.Errorf("%s:%d selected UPM view has unsupported period count %d", row.Unit.Bank, row.Unit.Index, len(upm))
	}
	var expectedSamples int
	for _, period := range upm {
		expectedSamples += int(period) * 2
	}
	if expectedSamples != len(prepared.Samples) {
		return paul2013ContextTimelineSource{}, fmt.Errorf(
			"%s:%d selected UPM view covers %d samples, row extraction produced %d",
			row.Unit.Bank, row.Unit.Index, expectedSamples, len(prepared.Samples),
		)
	}
	return paul2013ContextTimelineSource{samples: prepared.Samples, upm: append([]byte(nil), upm...)}, nil
}

func contextUPMCount(source *paul2013ContextTimelineSource) int {
	if source == nil {
		return 0
	}
	return len(source.upm)
}

func buildPaul2013ContextUPMPeriods(
	current, left, right *paul2013ContextTimelineSource,
	plan Paul2013ContextEdgePlan,
) ([]Paul2013UPMPeriod, error) {
	if current == nil || len(current.upm) == 0 {
		return nil, errors.New("context period builder has no current unit")
	}
	periodCount := len(current.upm)
	leftCount := plan.LeftContextCount
	rightCount := plan.RightContextCount
	if leftCount > 0 && (left == nil || len(left.upm) < leftCount) {
		return nil, errors.New("left context UPM view is shorter than its edge count")
	}
	if rightCount > 0 && (right == nil || len(right.upm) < rightCount) {
		return nil, errors.New("right context UPM view is shorter than its edge count")
	}
	periods := make([]Paul2013UPMPeriod, periodCount)
	currentOffset := 0
	leftOffsets := upmSampleOffsets(left)
	rightOffsets := upmSampleOffsets(right)
	for index, period := range current.upm {
		length := int(period) * 2
		currentWindow := current.samples[currentOffset : currentOffset+length]
		currentOffset += length
		windows := Paul2013UPMPeriodWindows{Current: currentWindow}
		if index < leftCount {
			leftLength := int(left.upm[index]) * 2
			start := leftOffsets[index]
			windows.Left = left.samples[start : start+leftLength]
		}
		firstRightPeriod := periodCount - rightCount
		if rightCount > 0 && index >= firstRightPeriod {
			rightIndex := len(right.upm) - rightCount + index - firstRightPeriod
			rightLength := int(right.upm[rightIndex]) * 2
			start := rightOffsets[rightIndex]
			windows.Right = right.samples[start : start+rightLength]
		}
		periods[index] = Paul2013UPMPeriod{Length: length, Windows: windows}
	}
	return periods, nil
}

func upmSampleOffsets(source *paul2013ContextTimelineSource) []int {
	if source == nil {
		return nil
	}
	offsets := make([]int, len(source.upm))
	for index := 1; index < len(source.upm); index++ {
		offsets[index] = offsets[index-1] + int(source.upm[index-1])*2
	}
	return offsets
}
