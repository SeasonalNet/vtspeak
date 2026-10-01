package synthesis

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/selection"
)

// CursorTimelineRenderer is an opt-in default-control renderer that applies
// the captured normal-row cursor advance by omitting each non-final unit's
// trailing UPM span. It does not reconstruct or blend the omitted boundary
// audio, so it is a timing experiment rather than a parity renderer.
type CursorTimelineRenderer struct {
	Units             UnitReader
	ApplyObservedGain bool
}

// UPMEqualSpanJoinRenderer joins complete selected-unit payloads using their
// UPM edge spans when each incoming leading span equals the prior trailing
// span. Complete payloads are diagnostic source windows; the native timeline
// builder can select subwindows that this renderer does not reconstruct.
type UPMEqualSpanJoinRenderer struct {
	Units             UnitReader
	ApplyObservedGain bool
}

// UPMTimelineJoinRenderer joins complete selected-unit payloads using the
// direct 2013 timeline path, including unequal adjacent UPM edge spans.
type UPMTimelineJoinRenderer struct {
	Units             UnitReader
	ApplyObservedGain bool
}

// Paul2013TimelinePCMInput is the subset of one already-produced native
// timeline row needed to extract and prepare its selected-unit source window.
// RowType is the row-kind byte at +0x27; Mode is the view byte at +0x26. The
// caller supplies both because their producers from text remain unresolved.
type Paul2013TimelinePCMInput struct {
	Unit                 selection.UnitRef
	RowType              uint8
	Mode                 uint8
	SampleCount          int
	LeadingSpan          int
	TrailingSpan         int
	GainPercent          int32
	SourceScratchSamples []int16
}

// BuildPaul2013TimelinePCMInputs maps recovered normal timeline rows to the
// selected-unit window inputs consumed by RenderPaul2013TimelineRows. The
// native normal-row kind is 2; its mode, sample count, doubled UPM edge spans,
// unit reference, and percentage gain are carried without reinterpretation.
func BuildPaul2013TimelinePCMInputs(
	rows []Paul2013NormalTimelineRow,
) ([]Paul2013TimelinePCMInput, error) {
	if len(rows) == 0 {
		return nil, errors.New("normal timeline input builder has no rows")
	}
	inputs := make([]Paul2013TimelinePCMInput, len(rows))
	for index, row := range rows {
		if row.RowKind != 2 {
			return nil, fmt.Errorf("normal timeline row %d has unsupported row kind %d", index, row.RowKind)
		}
		if uint64(row.SampleCount) > uint64(^uint(0)>>1) {
			return nil, fmt.Errorf("normal timeline row %d sample count exceeds platform integer range", index)
		}
		if row.Controls.Gain > uint32(^uint32(0)>>1) {
			return nil, fmt.Errorf("normal timeline row %d gain exceeds signed 32-bit range", index)
		}
		inputs[index] = Paul2013TimelinePCMInput{
			Unit: row.Unit, RowType: row.RowKind, Mode: row.Mode,
			SampleCount: int(row.SampleCount),
			LeadingSpan: int(row.LeadingSpan), TrailingSpan: int(row.TrailingSpan),
			GainPercent: int32(row.Controls.Gain),
		}
	}
	return inputs, nil
}

// RenderPaul2013NormalTimelineRows carries already-built normal timeline rows
// through selected-unit extraction, per-row gain, and the direct 2013 join.
// Row mode, controls, and selected-unit production remain caller inputs.
func RenderPaul2013NormalTimelineRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013NormalTimelineRow,
) (PCM, error) {
	inputs, err := BuildPaul2013TimelinePCMInputs(rows)
	if err != nil {
		return PCM{}, err
	}
	return RenderPaul2013TimelineRows(ctx, units, inputs)
}

// RenderPaul2013EqualSpanTimelineRows connects captured row metadata to the
// mode-dependent unit window, per-row gain, and measured equal-span join.
// It supports the default pitch/speed row path only and rejects adjacent
// unequal edge spans in MixPaul2013EqualSpanTimeline. It does not generate
// timeline rows or claim full-render parity.
func RenderPaul2013EqualSpanTimelineRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013TimelinePCMInput,
) (PCM, error) {
	return renderPaul2013TimelineRows(ctx, units, rows, true)
}

// RenderPaul2013TimelineRows connects captured row metadata to selected-unit
// windows, per-row gain, and the direct 2013 join path, including unequal
// adjacent edge spans. It supports the default pitch/speed path and does not
// generate timeline rows from text or claim full-render parity.
func RenderPaul2013TimelineRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013TimelinePCMInput,
) (PCM, error) {
	return renderPaul2013TimelineRows(ctx, units, rows, false)
}

func renderPaul2013TimelineRows(
	ctx context.Context,
	units UnitReader,
	rows []Paul2013TimelinePCMInput,
	requireEqualSpans bool,
) (PCM, error) {
	if units == nil {
		return PCM{}, errors.New("prepared timeline renderer has no voice model")
	}
	if len(rows) == 0 {
		return PCM{}, errors.New("prepared timeline renderer requires at least one row")
	}
	prepared := make([]TimelinePCMRow, len(rows))
	for index, row := range rows {
		if row.RowType != 2 {
			return PCM{}, fmt.Errorf("timeline row %d has unsupported row kind %d", index, row.RowType)
		}
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		unit, err := units.ReadUnit(row.Unit.Bank, row.Unit.Index)
		if err != nil {
			return PCM{}, fmt.Errorf("read prepared timeline unit %d (%s:%d): %w", index, row.Unit.Bank, row.Unit.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return PCM{}, fmt.Errorf("prepared timeline unit %d (%s:%d) has invalid PCM byte length %d", index, row.Unit.Bank, row.Unit.Index, len(unit.PCM))
		}
		decodedSampleCount := len(unit.PCM) / 2
		if len(row.SourceScratchSamples) > int(^uint(0)>>1)-decodedSampleCount {
			return PCM{}, fmt.Errorf("timeline row %d scratch tail exceeds platform integer range", index)
		}
		source := make([]int16, decodedSampleCount, decodedSampleCount+len(row.SourceScratchSamples))
		for sampleIndex := range source {
			source[sampleIndex] = int16(binary.LittleEndian.Uint16(unit.PCM[sampleIndex*2:]))
		}
		source = append(source, row.SourceScratchSamples...)
		window, err := BuildPaul2013TimelinePCMRow(
			source, decodedSampleCount, row.Mode, row.SampleCount,
			row.LeadingSpan, row.TrailingSpan,
		)
		if err != nil {
			return PCM{}, fmt.Errorf("prepare timeline row %d (%s:%d): %w", index, row.Unit.Bank, row.Unit.Index, err)
		}
		prepared[index], err = ScalePaul2013TimelinePCMRow(window, row.GainPercent)
		if err != nil {
			return PCM{}, fmt.Errorf("scale timeline row %d (%s:%d): %w", index, row.Unit.Bank, row.Unit.Index, err)
		}
	}
	var samples []int16
	var err error
	if requireEqualSpans {
		samples, err = MixPaul2013EqualSpanTimeline(prepared)
	} else {
		samples, err = MixPaul2013Timeline(prepared)
	}
	if err != nil {
		return PCM{}, fmt.Errorf("join prepared timeline rows: %w", err)
	}
	return PCM{SampleRate: 16000, Samples: samples}, nil
}

// Render decodes selected units and applies the equal-span join path at
// default pitch and speed. Unequal edge spans fail closed in this compatibility
// renderer; use UPMTimelineJoinRenderer for the general direct path.
func (renderer UPMEqualSpanJoinRenderer) Render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
) (PCM, error) {
	return renderer.render(ctx, units, controls, true)
}

// Render decodes selected units and applies the direct 2013 timeline join at
// default pitch and speed, including when neighboring edge spans differ.
// Complete-unit windows remain diagnostic substitutes for native selected
// subwindows; text-based window production is not implemented here.
func (renderer UPMTimelineJoinRenderer) Render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
) (PCM, error) {
	return UPMEqualSpanJoinRenderer(renderer).render(ctx, units, controls, false)
}

func (renderer UPMEqualSpanJoinRenderer) render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
	requireEqualSpans bool,
) (PCM, error) {
	if renderer.Units == nil {
		return PCM{}, errors.New("UPM timeline join renderer has no voice model")
	}
	if len(units) == 0 {
		return PCM{}, errors.New("unit sequence is empty")
	}
	controls = NormalizePaul2013Controls(controls)
	if controls.Pitch != 100 || controls.Speed != 100 {
		if requireEqualSpans {
			return PCM{}, errors.New("equal-span join renderer only supports default pitch and speed")
		}
		return PCM{}, errors.New("UPM timeline join renderer only supports default pitch and speed")
	}
	if controls.Volume != 200 && !renderer.ApplyObservedGain {
		return PCM{}, errors.New("observed gain must be enabled for non-default volume")
	}
	gainPercent := int32(100)
	if renderer.ApplyObservedGain {
		gainPercent = controls.Volume
	}

	rows := make([]TimelinePCMRow, len(units))
	for unitPosition, reference := range units {
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		unit, err := renderer.Units.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return PCM{}, fmt.Errorf("read joined unit %d (%s:%d): %w", unitPosition, reference.Bank, reference.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return PCM{}, fmt.Errorf("joined unit %d (%s:%d) has invalid PCM byte length %d", unitPosition, reference.Bank, reference.Index, len(unit.PCM))
		}
		if err := validateRendererUnitTiming(unit.UPM, len(unit.PCM)/2); err != nil {
			return PCM{}, fmt.Errorf("joined unit %d (%s:%d) has invalid timing metadata: %w", unitPosition, reference.Bank, reference.Index, err)
		}
		samples := make([]int16, len(unit.PCM)/2)
		for sampleIndex := range samples {
			sample := int16(binary.LittleEndian.Uint16(unit.PCM[sampleIndex*2:]))
			samples[sampleIndex] = scaleSample(sample, gainPercent)
		}
		rows[unitPosition] = TimelinePCMRow{
			Samples: samples, LeadingSpan: int(unit.UPM[0]) * 2,
			TrailingSpan: int(unit.UPM[len(unit.UPM)-1]) * 2,
		}
	}
	var samples []int16
	var err error
	if requireEqualSpans {
		samples, err = MixPaul2013EqualSpanTimeline(rows)
	} else {
		samples, err = MixPaul2013Timeline(rows)
	}
	if err != nil {
		return PCM{}, fmt.Errorf("join selected-unit timeline: %w", err)
	}
	return PCM{SampleRate: 16000, Samples: samples}, nil
}

// Render decodes selected units and places their samples according to the
// normal timeline cursor rule. Synthetic sentence rows, pitch/speed changes,
// and edge-window blending are handled elsewhere or remain unimplemented.
func (renderer CursorTimelineRenderer) Render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
) (PCM, error) {
	if renderer.Units == nil {
		return PCM{}, errors.New("cursor timeline renderer has no voice model")
	}
	if len(units) == 0 {
		return PCM{}, errors.New("unit sequence is empty")
	}
	controls = NormalizePaul2013Controls(controls)
	if controls.Pitch != 100 || controls.Speed != 100 {
		return PCM{}, errors.New("cursor timeline renderer only supports default pitch and speed")
	}
	if controls.Volume != 200 && !renderer.ApplyObservedGain {
		return PCM{}, errors.New("observed gain must be enabled for non-default volume")
	}
	gainPercent := int32(100)
	if renderer.ApplyObservedGain {
		gainPercent = controls.Volume
	}

	rows := make([]TimelineRow, len(units))
	output := make([]int16, 0)
	for index, reference := range units {
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		unit, err := renderer.Units.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return PCM{}, fmt.Errorf("read timeline unit %d (%s:%d): %w", index, reference.Bank, reference.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return PCM{}, fmt.Errorf("timeline unit %d (%s:%d) has invalid PCM byte length %d", index, reference.Bank, reference.Index, len(unit.PCM))
		}
		if len(unit.UPM) == 0 {
			return PCM{}, fmt.Errorf("timeline unit %d (%s:%d) has no UPM periods", index, reference.Bank, reference.Index)
		}
		sampleCount := len(unit.PCM) / 2
		var upmSampleCount int
		for periodIndex, period := range unit.UPM {
			if period == 0 {
				return PCM{}, fmt.Errorf("timeline unit %d (%s:%d) UPM period %d is zero", index, reference.Bank, reference.Index, periodIndex)
			}
			upmSampleCount += int(period) * 2
		}
		if upmSampleCount != sampleCount {
			return PCM{}, fmt.Errorf("timeline unit %d (%s:%d) UPM spans %d samples but PCM has %d", index, reference.Bank, reference.Index, upmSampleCount, sampleCount)
		}
		leadingSpan := int(unit.UPM[0]) * 2
		trailingSpan := int(unit.UPM[len(unit.UPM)-1]) * 2
		rows[index] = TimelineRow{
			Unit: reference, SampleCount: sampleCount,
			LeadingSpan: leadingSpan, TrailingSpan: trailingSpan,
		}
		limit := sampleCount
		if index+1 < len(units) {
			limit -= trailingSpan
		}
		for sampleIndex := 0; sampleIndex < limit; sampleIndex++ {
			offset := sampleIndex * 2
			sample := int16(binary.LittleEndian.Uint16(unit.PCM[offset : offset+2]))
			output = append(output, scaleSample(sample, gainPercent))
		}
	}
	frameCount, err := Paul2013TimelineOutputFrames(rows)
	if err != nil {
		return PCM{}, fmt.Errorf("plan selected-unit timeline: %w", err)
	}
	if len(output) != frameCount {
		return PCM{}, fmt.Errorf("cursor renderer emitted %d samples, planned %d", len(output), frameCount)
	}
	return PCM{SampleRate: 16000, Samples: output}, nil
}
