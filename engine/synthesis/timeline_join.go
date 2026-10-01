package synthesis

import (
	"errors"
	"fmt"
	"math"
)

// TimelinePCMRow is one prepared normal-row sample window. LeadingSpan and
// TrailingSpan are sample counts on the 16 kHz grid, not raw UPM bytes.
type TimelinePCMRow struct {
	Samples      []int16
	LeadingSpan  int
	TrailingSpan int
}

// BuildPaul2013TimelinePCMRow extracts the source window selected by one
// captured timeline row. FUN_1002c8b0 starts at sample zero for combined and
// first-side views; only view mode 2 starts at decodedSampleCount minus the
// leading span. The backing buffer may include scratch samples after the
// decoded segment, so bounds are checked against source while the mode-2
// offset is based on decodedSampleCount. Row counts and spans use 16 kHz
// samples. Scaling from the row's separate scale field is intentionally left
// to the caller so gain is applied before joining.
func BuildPaul2013TimelinePCMRow(
	source []int16,
	decodedSampleCount int,
	viewMode uint8,
	sampleCount int,
	leadingSpan int,
	trailingSpan int,
) (TimelinePCMRow, error) {
	if len(source) == 0 {
		return TimelinePCMRow{}, errors.New("timeline source unit has no samples")
	}
	if decodedSampleCount < 1 || decodedSampleCount > len(source) {
		return TimelinePCMRow{}, errors.New("decoded sample count falls outside the timeline source buffer")
	}
	if sampleCount < 1 || sampleCount > math.MaxInt16 {
		return TimelinePCMRow{}, errors.New("timeline row sample count must be in the range 1 through 32767")
	}
	if leadingSpan < 0 || leadingSpan > sampleCount || trailingSpan < 0 || trailingSpan > sampleCount {
		return TimelinePCMRow{}, errors.New("timeline row edge span exceeds its sample count")
	}
	start := 0
	if viewMode == Paul2013TimelineSecondSideView {
		if leadingSpan > decodedSampleCount {
			return TimelinePCMRow{}, errors.New("timeline row leading span exceeds its decoded sample count")
		}
		start = decodedSampleCount - leadingSpan
	}
	if sampleCount > len(source)-start {
		return TimelinePCMRow{}, fmt.Errorf(
			"timeline mode %d selects samples [%d,%d) from a %d-sample unit",
			viewMode, start, start+sampleCount, len(source),
		)
	}
	return TimelinePCMRow{
		Samples:     append([]int16(nil), source[start:start+sampleCount]...),
		LeadingSpan: leadingSpan, TrailingSpan: trailingSpan,
	}, nil
}

// ScalePaul2013TimelinePCMRow applies the row's observed percentage gain
// before timeline joining. Integer division truncates toward zero and each
// scaled sample saturates to the native [-32767, 32766] range.
func ScalePaul2013TimelinePCMRow(row TimelinePCMRow, gainPercent int32) (TimelinePCMRow, error) {
	if len(row.Samples) == 0 {
		return TimelinePCMRow{}, errors.New("timeline row has no samples")
	}
	if gainPercent < 0 || gainPercent > 500 {
		return TimelinePCMRow{}, errors.New("timeline row gain must be in the range 0 through 500 percent")
	}
	if row.LeadingSpan < 0 || row.TrailingSpan < 0 ||
		row.LeadingSpan > len(row.Samples) || row.TrailingSpan > len(row.Samples) {
		return TimelinePCMRow{}, errors.New("timeline row has invalid edge spans")
	}
	scaled := TimelinePCMRow{
		Samples:      make([]int16, len(row.Samples)),
		LeadingSpan:  row.LeadingSpan,
		TrailingSpan: row.TrailingSpan,
	}
	for index, sample := range row.Samples {
		scaled.Samples[index] = scaleSample(sample, gainPercent)
	}
	return scaled, nil
}

// MixPaul2013EqualSpanTimeline joins prepared normal rows when every incoming
// leading span matches the preceding row's pending trailing span. The first
// row fades in from an empty carry and the final row drains its tail to zero.
// It ports the equal-span branch of FUN_1002aac0; the coefficient curve is
// read from the exact extracted 2013 table. Rows must already contain the
// selected source windows.
func MixPaul2013EqualSpanTimeline(rows []TimelinePCMRow) ([]int16, error) {
	return mixPaul2013Timeline(rows, true)
}

// MixPaul2013Timeline joins prepared normal rows with arbitrary adjacent edge
// spans. FUN_1002aac0 blends over each current leading span; when that differs
// from the pending prior trailing span, the prior contribution uses the
// shorter span's coefficient schedule and is zero-padded or truncated to the
// current leading span. Rows must already contain the selected source windows.
func MixPaul2013Timeline(rows []TimelinePCMRow) ([]int16, error) {
	return mixPaul2013Timeline(rows, false)
}

func mixPaul2013Timeline(rows []TimelinePCMRow, requireEqualSpans bool) ([]int16, error) {
	if len(rows) == 0 {
		return nil, errors.New("timeline join has no rows")
	}
	for index, row := range rows {
		if len(row.Samples) == 0 {
			return nil, fmt.Errorf("timeline row %d has no samples", index)
		}
		if row.LeadingSpan < 0 || row.TrailingSpan < 0 ||
			row.LeadingSpan > len(row.Samples) || row.TrailingSpan > len(row.Samples)-row.LeadingSpan {
			return nil, fmt.Errorf("timeline row %d has invalid edge spans", index)
		}
		if requireEqualSpans && index > 0 && row.LeadingSpan != rows[index-1].TrailingSpan {
			return nil, fmt.Errorf("timeline row %d leading span %d does not match previous trailing span %d", index, row.LeadingSpan, rows[index-1].TrailingSpan)
		}
	}

	outputLength := 0
	for index, row := range rows {
		advance := len(row.Samples)
		if index+1 < len(rows) {
			advance -= row.TrailingSpan
		}
		if advance > int(^uint(0)>>1)-outputLength {
			return nil, errors.New("timeline join output exceeds the platform integer range")
		}
		outputLength += advance
	}
	output := make([]int16, 0, outputLength)
	var pending []int16
	for rowIndex, row := range rows {
		for sampleIndex := 0; sampleIndex < row.LeadingSpan; sampleIndex++ {
			output = append(output, mixTimelineLeadingSample(pending, row.Samples[sampleIndex], sampleIndex, row.LeadingSpan))
		}

		middleStart := row.LeadingSpan
		middleEnd := len(row.Samples) - row.TrailingSpan
		output = append(output, row.Samples[middleStart:middleEnd]...)
		if rowIndex+1 < len(rows) {
			pending = row.Samples[middleEnd:]
			continue
		}
		for sampleIndex, sample := range row.Samples[middleEnd:] {
			output = append(output, mixTimelineJoinSample(sample, 0, sampleIndex, row.TrailingSpan))
		}
	}
	if len(output) != outputLength {
		return nil, fmt.Errorf("timeline join emitted %d samples, planned %d", len(output), outputLength)
	}
	return output, nil
}

func mixTimelineLeadingSample(pending []int16, current int16, index, currentSpan int) int16 {
	phase := index * paul2013WindowTableHalf / currentSpan
	currentLane := timelineJoinLane(current, paul2013WindowCoefficient(phase))
	if index >= len(pending) {
		return currentLane
	}
	previousSpan := len(pending)
	if previousSpan > currentSpan {
		previousSpan = currentSpan
	}
	previousPhase := paul2013WindowTableHalf + index*paul2013WindowTableHalf/previousSpan
	previousLane := timelineJoinLane(pending[index], paul2013WindowCoefficient(previousPhase))
	return int16(uint16(int32(previousLane) + int32(currentLane)))
}

func mixTimelineJoinSample(previous, current int16, index, span int) int16 {
	if span == 0 {
		return current
	}
	phase := index * paul2013WindowTableHalf / span
	currentWeight := paul2013WindowCoefficient(phase)
	previousWeight := paul2013WindowCoefficient(paul2013WindowTableHalf + phase)
	previousLane := timelineJoinLane(previous, previousWeight)
	currentLane := timelineJoinLane(current, currentWeight)
	return int16(uint16(int32(previousLane) + int32(currentLane)))
}

func timelineJoinLane(sample int16, weight float64) int16 {
	return paul2013WeightedLaneSample(weight * float64(sample))
}
