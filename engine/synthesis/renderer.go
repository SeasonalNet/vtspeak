// Package synthesis defines the selected-unit to PCM boundary.
package synthesis

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

// Controls carries raw API values. Negative values select engine defaults;
// NormalizePaul2013Controls resolves those sentinels and applies the observed
// bounds before values enter the synthesis path.
type Controls struct {
	Pitch  int32
	Speed  int32
	Volume int32
}

// PCM is mono signed 16-bit audio. The first supported output profile is
// 16 kHz, matching the observed Paul file path.
type PCM struct {
	SampleRate uint32
	Samples    []int16
}

type Renderer interface {
	Render(context.Context, []selection.UnitRef, Controls) (PCM, error)
}

// UPMSegment is one measured five-word timing record between adjacent UPM
// periods. The values preserve the DLL record's arithmetic roles without
// assigning physical names to its control fields.
type UPMSegment struct {
	StartSample      int32
	FirstPeriod      int32
	SecondPeriod     int32
	SpeedRatio       int32
	PitchControlWord int32
}

// UPMSegmentResampleStep records one segment chosen by the pitch renderer's
// moving-coordinate scan and the period length it will emit.
type UPMSegmentResampleStep struct {
	SegmentIndex        int
	TargetCoordinate    int32
	PriorOutputSamples  int32
	SourcePeriod        int32
	ResampledPeriodSize int32
}

// BuildPaul2013UPMSegments ports the observed non-default-pitch segment
// construction for a raw unit UPM vector. UPM bytes are doubled onto the
// 16 kHz sample grid; each adjacent pair becomes one 20-byte/five-word record.
func BuildPaul2013UPMSegments(upm []byte, pitchControlWord, speedControlWord int32) ([]UPMSegment, error) {
	if len(upm) < 2 {
		return nil, errors.New("UPM segment construction requires at least two periods")
	}
	if pitchControlWord <= 0 || speedControlWord <= 0 {
		return nil, errors.New("UPM segment controls must be positive")
	}
	speedRatio := (int64(speedControlWord>>1) + 10000) / int64(speedControlWord)
	if speedRatio <= 0 || speedRatio > int64(^uint32(0)>>1) {
		return nil, errors.New("UPM speed ratio is outside the positive signed 32-bit range")
	}
	segments := make([]UPMSegment, len(upm)-1)
	var startSample int64
	for index := range segments {
		firstPeriod := int32(upm[index]) * 2
		secondPeriod := int32(upm[index+1]) * 2
		if firstPeriod == 0 || secondPeriod == 0 {
			return nil, fmt.Errorf("UPM period %d or %d is zero", index, index+1)
		}
		if startSample > int64(^uint32(0)>>1) {
			return nil, errors.New("UPM cumulative position exceeds signed 32-bit range")
		}
		segments[index] = UPMSegment{
			StartSample:      int32(startSample),
			FirstPeriod:      firstPeriod,
			SecondPeriod:     secondPeriod,
			SpeedRatio:       int32(speedRatio),
			PitchControlWord: pitchControlWord,
		}
		startSample += int64(firstPeriod)
	}
	return segments, nil
}

// SelectPaul2013UPMSegment selects the segment whose start plus first period
// is nearest the target coordinate. The legacy scan keeps the first segment
// on equal distances and returns no segment once the target reaches the final
// segment's end coordinate.
func SelectPaul2013UPMSegment(segments []UPMSegment, target int32) (int, bool, error) {
	if len(segments) == 0 {
		return 0, false, errors.New("UPM segment selection requires at least one segment")
	}
	last := segments[len(segments)-1]
	end := int64(last.StartSample) + int64(last.FirstPeriod) + int64(last.SecondPeriod)
	if int64(target) >= end {
		return 0, false, nil
	}
	selected := 0
	bestDistance := int64(20000)
	for index, segment := range segments {
		coordinate := int64(segment.StartSample) + int64(segment.FirstPeriod)
		distance := coordinate - int64(target)
		if distance < 0 {
			distance = -distance
		}
		if distance < bestDistance {
			selected = index
			bestDistance = distance
		}
	}
	return selected, true, nil
}

// PlanPaul2013UPMSegmentResampling follows FUN_1002afb0's segment-selection
// loop without reading or interpolating audio. It derives the initial target
// coordinate from the first segment's first period and speed ratio, then
// advances using each selected segment's second period.
func PlanPaul2013UPMSegmentResampling(segments []UPMSegment) ([]UPMSegmentResampleStep, error) {
	if len(segments) == 0 {
		return nil, errors.New("UPM resampling plan requires at least one segment")
	}
	for index, segment := range segments {
		if segment.FirstPeriod <= 0 || segment.SecondPeriod <= 0 ||
			segment.SpeedRatio <= 0 || segment.PitchControlWord <= 0 {
			return nil, fmt.Errorf("UPM segment %d has a nonpositive period or control", index)
		}
	}
	initialCoordinate := int64(segments[0].FirstPeriod) * 100 / int64(segments[0].SpeedRatio)
	if initialCoordinate > int64(^uint32(0)>>1) {
		return nil, errors.New("initial UPM target coordinate exceeds signed 32-bit range")
	}
	targetCoordinate := int32(initialCoordinate)
	var priorOutputSamples int32
	steps := make([]UPMSegmentResampleStep, 0, len(segments))
	for stepIndex := 0; stepIndex < len(segments); stepIndex++ {
		segmentIndex, found, err := SelectPaul2013UPMSegment(segments, targetCoordinate)
		if err != nil {
			return nil, err
		}
		if !found {
			return steps, nil
		}
		segment := segments[segmentIndex]
		periodLength, err := Paul2013UPMSegmentTargetLength(segments, segmentIndex, priorOutputSamples)
		if err != nil {
			return nil, fmt.Errorf("calculate target length for segment %d: %w", segmentIndex, err)
		}
		steps = append(steps, UPMSegmentResampleStep{
			SegmentIndex:        segmentIndex,
			TargetCoordinate:    targetCoordinate,
			PriorOutputSamples:  priorOutputSamples,
			SourcePeriod:        segment.FirstPeriod,
			ResampledPeriodSize: periodLength,
		})
		priorOutput := int64(priorOutputSamples) + int64(periodLength)
		if priorOutput > int64(^uint32(0)>>1) {
			return nil, errors.New("cumulative UPM resampled length exceeds signed 32-bit range")
		}
		priorOutputSamples = int32(priorOutput)
		nextCoordinate := (priorOutput + int64(segment.SecondPeriod)) * 100 / int64(segment.SpeedRatio)
		if nextCoordinate <= int64(targetCoordinate) {
			return nil, errors.New("UPM target coordinate failed to advance")
		}
		if nextCoordinate > int64(^uint32(0)>>1) {
			return nil, errors.New("next UPM target coordinate exceeds signed 32-bit range")
		}
		targetCoordinate = int32(nextCoordinate)
	}
	_, found, err := SelectPaul2013UPMSegment(segments, targetCoordinate)
	if err != nil {
		return nil, err
	}
	if !found {
		return steps, nil
	}
	return nil, errors.New("UPM segment scan exceeded the segment count")
}

// Paul2013ResampledPeriodLength applies the captured pitch-period ratio and
// clamp. LimitPaul2013PeriodLength applies the separate cumulative budget cap.
func Paul2013ResampledPeriodLength(sourcePeriod, pitchControlWord int32) (int32, error) {
	if sourcePeriod <= 0 || pitchControlWord <= 0 {
		return 0, errors.New("source period and pitch control must be positive")
	}
	length := int64(sourcePeriod) * 100 / int64(pitchControlWord)
	minimum := int64(sourcePeriod) / 16
	maximum := int64(sourcePeriod) * 31 / 16
	if length < minimum {
		length = minimum
	}
	if length > maximum {
		length = maximum
	}
	if length > int64(^uint32(0)>>1) {
		return 0, errors.New("resampled period length exceeds signed 32-bit range")
	}
	return int32(length), nil
}

// LimitPaul2013PeriodLength applies the cumulative sample-budget cap in
// FUN_1002afb0. The decompile limits this period so priorOutputSamples plus
// the new length does not exceed sourcePeriod plus contextSampleBudget. The
// runtime meaning and lifecycle of that budget remain caller-owned.
func LimitPaul2013PeriodLength(
	sourcePeriod, targetLength, priorOutputSamples, contextSampleBudget int32,
) (int32, error) {
	if sourcePeriod <= 0 || targetLength <= 0 {
		return 0, errors.New("source period and target length must be positive")
	}
	if priorOutputSamples < 0 || contextSampleBudget < 0 {
		return 0, errors.New("prior output samples and context budget must be nonnegative")
	}
	maximumLength := int64(sourcePeriod) - int64(priorOutputSamples) + int64(contextSampleBudget)
	if maximumLength <= 0 {
		return 0, errors.New("context budget leaves no positive length for this period")
	}
	if int64(targetLength) > maximumLength {
		targetLength = int32(maximumLength)
	}
	return targetLength, nil
}

// Paul2013UPMSegmentTargetLength applies the segment's pitch ratio and the
// cumulative cap used by FUN_1002afb0. The decompile reads the budget from the
// first segment's FirstPeriod word and advances priorOutputSamples across the
// selected segment sequence.
func Paul2013UPMSegmentTargetLength(
	segments []UPMSegment,
	segmentIndex int,
	priorOutputSamples int32,
) (int32, error) {
	if len(segments) == 0 {
		return 0, errors.New("UPM segment target length requires at least one segment")
	}
	if segmentIndex < 0 || segmentIndex >= len(segments) {
		return 0, errors.New("UPM segment index is outside the segment list")
	}
	budget := segments[0].FirstPeriod
	segment := segments[segmentIndex]
	targetLength, err := Paul2013ResampledPeriodLength(segment.FirstPeriod, segment.PitchControlWord)
	if err != nil {
		return 0, err
	}
	return LimitPaul2013PeriodLength(segment.FirstPeriod, targetLength, priorOutputSamples, budget)
}

// ResamplePaul2013Period linearly resamples one isolated sample window using
// the recovered pitch-period length rule. The target length and clamp follow
// the DLL trace; linear interpolation within that window is an explicit
// approximation because the DLL's full window-selection and interpolation
// behavior has not been ported. This helper does not apply context joins.
func ResamplePaul2013Period(samples []int16, pitchControlWord int32) ([]int16, error) {
	if len(samples) == 0 {
		return nil, errors.New("pitch-period sample window is empty")
	}
	if uint64(len(samples)) > uint64(^uint32(0)>>1) {
		return nil, errors.New("pitch-period sample window exceeds signed 32-bit range")
	}
	targetLength, err := Paul2013ResampledPeriodLength(int32(len(samples)), pitchControlWord)
	if err != nil {
		return nil, err
	}
	if targetLength <= 0 {
		return nil, errors.New("pitch-period target length is not positive")
	}
	return resamplePaul2013SamplesToLength(samples, targetLength)
}

// ResamplePaul2013PeriodWithinBudget applies the pitch-derived length and
// cumulative context-budget cap before linearly resampling the supplied
// period. It still operates on an isolated window, not reconstructed context.
func ResamplePaul2013PeriodWithinBudget(
	samples []int16,
	pitchControlWord, priorOutputSamples, contextSampleBudget int32,
) ([]int16, error) {
	if len(samples) == 0 {
		return nil, errors.New("pitch-period sample window is empty")
	}
	if uint64(len(samples)) > uint64(^uint32(0)>>1) {
		return nil, errors.New("pitch-period sample window exceeds signed 32-bit range")
	}
	targetLength, err := Paul2013ResampledPeriodLength(int32(len(samples)), pitchControlWord)
	if err != nil {
		return nil, err
	}
	targetLength, err = LimitPaul2013PeriodLength(
		int32(len(samples)),
		targetLength,
		priorOutputSamples,
		contextSampleBudget,
	)
	if err != nil {
		return nil, err
	}
	return resamplePaul2013SamplesToLength(samples, targetLength)
}

func resamplePaul2013SamplesToLength(samples []int16, targetLength int32) ([]int16, error) {
	if len(samples) == 0 {
		return nil, errors.New("pitch-period sample window is empty")
	}
	if targetLength <= 0 {
		return nil, errors.New("pitch-period target length is not positive")
	}
	if targetLength == int32(len(samples)) {
		return append([]int16(nil), samples...), nil
	}
	resampled := make([]int16, int(targetLength))
	if len(resampled) == 1 {
		resampled[0] = samples[0]
		return resampled, nil
	}
	lastSource := int64(len(samples) - 1)
	lastTarget := int64(len(resampled) - 1)
	for targetIndex := range resampled {
		positionNumerator := int64(targetIndex) * lastSource
		leftIndex := positionNumerator / lastTarget
		remainder := positionNumerator % lastTarget
		left := int64(samples[leftIndex])
		right := left
		if leftIndex+1 < int64(len(samples)) {
			right = int64(samples[leftIndex+1])
		}
		value := left + (right-left)*remainder/lastTarget
		resampled[targetIndex] = int16(value)
	}
	return resampled, nil
}

// UnitReader supplies decoded units from an already opened voice model.
type UnitReader interface {
	ReadUnit(bank string, index uint32) (voice.Unit, error)
}

// ConcatenatingRenderer is a diagnostic renderer for explicit unit sequences.
// It concatenates complete decoded units and can optionally apply the
// normalized volume gain from 0% through 500% with 16-bit saturation to each
// sample. It does not model timeline clipping, UPM joins, or prosody.
type ConcatenatingRenderer struct {
	Units             UnitReader
	ApplyObservedGain bool
}

// Render concatenates selected units at the observed 16 kHz output rate. Pitch
// and speed must resolve to their defaults. Without gain, volume must also be
// the default; enabling observed gain applies any normalized volume value.
func (renderer ConcatenatingRenderer) Render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
) (PCM, error) {
	if renderer.Units == nil {
		return PCM{}, errors.New("unit renderer has no voice model")
	}
	controls = NormalizePaul2013Controls(controls)
	if controls.Pitch != 100 || controls.Speed != 100 {
		return PCM{}, errors.New("unit concatenation does not support pitch or speed changes")
	}
	if controls.Volume != 200 && !renderer.ApplyObservedGain {
		return PCM{}, errors.New("observed gain must be enabled for non-default volume")
	}
	if len(units) == 0 {
		return PCM{}, errors.New("unit sequence is empty")
	}
	var samples []int16
	gainPercent := int32(100)
	if renderer.ApplyObservedGain {
		gainPercent = controls.Volume
	}
	for position, reference := range units {
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		unit, err := renderer.Units.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return PCM{}, fmt.Errorf("read selected unit %d (%s:%d): %w", position, reference.Bank, reference.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return PCM{}, fmt.Errorf("selected unit %d (%s:%d) has invalid PCM byte length %d", position, reference.Bank, reference.Index, len(unit.PCM))
		}
		for offset := 0; offset < len(unit.PCM); offset += 2 {
			sample := int16(binary.LittleEndian.Uint16(unit.PCM[offset : offset+2]))
			if renderer.ApplyObservedGain {
				sample = scaleSample(sample, gainPercent)
			}
			samples = append(samples, sample)
		}
	}
	if len(samples) == 0 {
		return PCM{}, errors.New("unit sequence produced no PCM samples")
	}
	return PCM{SampleRate: 16000, Samples: samples}, nil
}

// PeriodResamplingRenderer is an opt-in diagnostic renderer. It splits each
// selected unit at its UPM period boundaries and resamples each period in
// isolation. This provides a working consumer for the observed period-length
// rule, but omits the DLL's neighboring-window selection, context blending,
// timeline clipping, and speed behavior.
type PeriodResamplingRenderer struct {
	Units             UnitReader
	ApplyObservedGain bool
}

// Render processes UPM periods independently at 16 kHz after normalizing the
// raw controls. Speed must resolve to its default. Interpolation is
// approximate and output must not be treated as DLL-parity audio.
func (renderer PeriodResamplingRenderer) Render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
) (PCM, error) {
	if renderer.Units == nil {
		return PCM{}, errors.New("period renderer has no voice model")
	}
	if len(units) == 0 {
		return PCM{}, errors.New("unit sequence is empty")
	}
	controls = NormalizePaul2013Controls(controls)
	if controls.Speed != 100 {
		return PCM{}, errors.New("period renderer does not support speed changes")
	}
	pitchControlWord := controls.Pitch
	if controls.Volume != 200 && !renderer.ApplyObservedGain {
		return PCM{}, errors.New("observed gain must be enabled for non-default volume")
	}
	gainPercent := int32(100)
	if renderer.ApplyObservedGain {
		gainPercent = controls.Volume
	}
	var output []int16
	for unitPosition, reference := range units {
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		unit, err := renderer.Units.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return PCM{}, fmt.Errorf("read selected unit %d (%s:%d): %w", unitPosition, reference.Bank, reference.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return PCM{}, fmt.Errorf("selected unit %d (%s:%d) has invalid PCM byte length %d", unitPosition, reference.Bank, reference.Index, len(unit.PCM))
		}
		samples := make([]int16, len(unit.PCM)/2)
		for index := range samples {
			samples[index] = int16(binary.LittleEndian.Uint16(unit.PCM[index*2:]))
		}
		periodStart := 0
		for periodIndex, period := range unit.UPM {
			if err := ctx.Err(); err != nil {
				return PCM{}, err
			}
			periodLength := int(period) * 2
			periodEnd := periodStart + periodLength
			if periodLength == 0 || periodEnd > len(samples) {
				return PCM{}, fmt.Errorf("selected unit %d (%s:%d) UPM period %d exceeds decoded PCM", unitPosition, reference.Bank, reference.Index, periodIndex)
			}
			periodSamples, err := ResamplePaul2013Period(samples[periodStart:periodEnd], pitchControlWord)
			if err != nil {
				return PCM{}, fmt.Errorf("resample selected unit %d (%s:%d) period %d: %w", unitPosition, reference.Bank, reference.Index, periodIndex, err)
			}
			for _, sample := range periodSamples {
				output = append(output, scaleSample(sample, gainPercent))
			}
			periodStart = periodEnd
		}
		if periodStart != len(samples) {
			return PCM{}, fmt.Errorf("selected unit %d (%s:%d) UPM periods cover %d of %d samples", unitPosition, reference.Bank, reference.Index, periodStart, len(samples))
		}
	}
	if len(output) == 0 {
		return PCM{}, errors.New("unit sequence produced no PCM samples")
	}
	return PCM{SampleRate: 16000, Samples: output}, nil
}

// UPMSegmentPlanRenderer is an opt-in experiment that follows the recovered
// pitch/speed segment scan over each selected unit's own UPM spans. It does
// not reconstruct neighboring-unit sample windows or perform context joins.
type UPMSegmentPlanRenderer struct {
	Units             UnitReader
	ApplyObservedGain bool
}

// Render executes the segment plan using source spans from each selected
// unit. It normalizes raw pitch and speed inputs using the observed API
// defaults and clamps; sample interpolation remains approximate and output is
// not DLL parity.
func (renderer UPMSegmentPlanRenderer) Render(
	ctx context.Context,
	units []selection.UnitRef,
	controls Controls,
) (PCM, error) {
	if renderer.Units == nil {
		return PCM{}, errors.New("UPM segment renderer has no voice model")
	}
	if len(units) == 0 {
		return PCM{}, errors.New("unit sequence is empty")
	}
	controls = NormalizePaul2013Controls(controls)
	pitchControlWord := controls.Pitch
	speedControlWord := controls.Speed
	if controls.Volume != 200 && !renderer.ApplyObservedGain {
		return PCM{}, errors.New("observed gain must be enabled for non-default volume")
	}
	gainPercent := int32(100)
	if renderer.ApplyObservedGain {
		gainPercent = controls.Volume
	}

	var output []int16
	for unitPosition, reference := range units {
		if err := ctx.Err(); err != nil {
			return PCM{}, err
		}
		unit, err := renderer.Units.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return PCM{}, fmt.Errorf("read selected unit %d (%s:%d): %w", unitPosition, reference.Bank, reference.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return PCM{}, fmt.Errorf("selected unit %d (%s:%d) has invalid PCM byte length %d", unitPosition, reference.Bank, reference.Index, len(unit.PCM))
		}
		segments, err := BuildPaul2013UPMSegments(unit.UPM, pitchControlWord, speedControlWord)
		if err != nil {
			return PCM{}, fmt.Errorf("build selected unit %d (%s:%d) UPM segments: %w", unitPosition, reference.Bank, reference.Index, err)
		}
		plan, err := PlanPaul2013UPMSegmentResampling(segments)
		if err != nil {
			return PCM{}, fmt.Errorf("plan selected unit %d (%s:%d) UPM segments: %w", unitPosition, reference.Bank, reference.Index, err)
		}
		samples := make([]int16, len(unit.PCM)/2)
		for index := range samples {
			samples[index] = int16(binary.LittleEndian.Uint16(unit.PCM[index*2:]))
		}
		unitOutputStart := len(output)
		priorOutputSamples := 0
		for stepIndex, step := range plan {
			if err := ctx.Err(); err != nil {
				return PCM{}, err
			}
			if int64(priorOutputSamples) != int64(step.PriorOutputSamples) {
				return PCM{}, fmt.Errorf("selected unit %d (%s:%d) plan step %d has a discontinuous output count", unitPosition, reference.Bank, reference.Index, stepIndex)
			}
			segment := segments[step.SegmentIndex]
			start := int64(segment.StartSample)
			firstPeriodEnd := start + int64(step.SourcePeriod)
			secondPeriodEnd := firstPeriodEnd + int64(segment.SecondPeriod)
			if start < 0 || firstPeriodEnd > int64(len(samples)) ||
				secondPeriodEnd > int64(len(samples)) || step.SourcePeriod != segment.FirstPeriod {
				return PCM{}, fmt.Errorf("selected unit %d (%s:%d) planned segment %d has an invalid PCM span", unitPosition, reference.Bank, reference.Index, step.SegmentIndex)
			}
			period, err := resamplePaul2013SamplesToLength(samples[int(start):int(firstPeriodEnd)], step.ResampledPeriodSize)
			if err != nil {
				return PCM{}, fmt.Errorf("resample selected unit %d (%s:%d) plan step %d: %w", unitPosition, reference.Bank, reference.Index, stepIndex, err)
			}
			outputEnd := unitOutputStart + priorOutputSamples + len(period) + int(segment.SecondPeriod)
			if outputEnd > len(output) {
				output = append(output, make([]int16, outputEnd-len(output))...)
			}
			for _, sample := range period {
				output[unitOutputStart+priorOutputSamples] = scaleSample(sample, gainPercent)
				priorOutputSamples++
			}
			tailOffset := unitOutputStart + priorOutputSamples
			for sampleIndex, sample := range samples[int(firstPeriodEnd):int(secondPeriodEnd)] {
				output[tailOffset+sampleIndex] = scaleSample(sample, gainPercent)
			}
		}
	}
	if len(output) == 0 {
		return PCM{}, errors.New("UPM segment plan produced no PCM samples")
	}
	return PCM{SampleRate: 16000, Samples: output}, nil
}

// scaleSample applies integer gain with the asymmetric limits recovered from
// FUN_1002c8b0: [-32767, 32766].
func scaleSample(sample int16, gainPercent int32) int16 {
	if gainPercent == 100 {
		return sample
	}
	scaled := int32(sample) * gainPercent / 100
	if scaled > 32766 {
		return 32766
	}
	if scaled < -32767 {
		return -32767
	}
	return int16(scaled)
}
