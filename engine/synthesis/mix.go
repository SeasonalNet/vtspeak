package synthesis

import (
	"errors"
	"fmt"
)

// WeightedWindowContribution is one already-windowed sample array and its
// position in the output timeline.
type WeightedWindowContribution struct {
	Samples []int16
	Offset  int
	Weight  int32
}

// Paul2013UPMPeriodWindows carries the source samples reconstructed for the
// left neighbor, current unit, and right neighbor of one UPM interval.
type Paul2013UPMPeriodWindows struct {
	Left    []int16
	Current []int16
	Right   []int16
}

// Paul2013UPMPeriod contains the target timeline length and prepared source
// windows for one period in the current unit.
type Paul2013UPMPeriod struct {
	Length  int
	Windows Paul2013UPMPeriodWindows
}

// MixPaul2013UPMInterval prepares the leading and trailing window forms for
// one UPM interval, then accumulates current/left/right contributions using
// edge entries at periodIndex and periodIndex+1, as FUN_1002d230 does. The
// source windows must already be selected and reconstructed from the model;
// the rising curve remains analytic.
func MixPaul2013UPMInterval(
	output []int16,
	outputOffset int,
	intervalLength int,
	periodIndex int,
	weights UPMEdgeWeights,
	windows Paul2013UPMPeriodWindows,
) error {
	if outputOffset < 0 || outputOffset > len(output) || intervalLength > len(output)-outputOffset {
		return errors.New("UPM interval falls outside the output")
	}
	contributions, err := paul2013UPMIntervalContributions(
		outputOffset,
		intervalLength,
		periodIndex,
		weights,
		windows,
	)
	if err != nil {
		return err
	}
	return AddPaul2013WeightedWindows(output, weights.Normalization, contributions)
}

// MixPaul2013UPMTimeline places all of one unit's prepared UPM periods at
// consecutive output offsets. The period array must match the edge-weight
// count; all windows are validated before any output sample is changed.
func MixPaul2013UPMTimeline(
	output []int16,
	outputOffset int,
	weights UPMEdgeWeights,
	periods []Paul2013UPMPeriod,
) error {
	if len(periods) == 0 {
		return errors.New("UPM timeline requires at least one period")
	}
	if len(weights.Left) != len(periods)+1 ||
		len(weights.Current) != len(weights.Left) || len(weights.Right) != len(weights.Left) {
		return errors.New("UPM timeline periods do not match the edge-weight arrays")
	}
	if outputOffset < 0 || outputOffset > len(output) {
		return errors.New("UPM timeline offset is outside the output")
	}

	contributions := make([]WeightedWindowContribution, 0, len(periods)*6)
	currentOffset := outputOffset
	for periodIndex, period := range periods {
		if period.Length < 1 || period.Length > len(output)-currentOffset {
			return fmt.Errorf("UPM period %d length falls outside the output", periodIndex)
		}
		periodContributions, err := paul2013UPMIntervalContributions(
			currentOffset,
			period.Length,
			periodIndex,
			weights,
			period.Windows,
		)
		if err != nil {
			return fmt.Errorf("prepare UPM period %d: %w", periodIndex, err)
		}
		contributions = append(contributions, periodContributions...)
		currentOffset += period.Length
	}
	return AddPaul2013WeightedWindows(output, weights.Normalization, contributions)
}

func paul2013UPMIntervalContributions(
	outputOffset int,
	intervalLength int,
	periodIndex int,
	weights UPMEdgeWeights,
	windows Paul2013UPMPeriodWindows,
) ([]WeightedWindowContribution, error) {
	if intervalLength < 1 || intervalLength > 32767 {
		return nil, errors.New("UPM interval length must be in the range 1 through 32767")
	}
	if weights.Normalization < 1 || weights.Normalization > 32767 {
		return nil, errors.New("UPM edge-weight normalization must be in the range 1 through 32767")
	}
	if len(weights.Left) < 2 || len(weights.Current) != len(weights.Left) || len(weights.Right) != len(weights.Left) {
		return nil, errors.New("UPM edge-weight arrays must have matching lengths of at least two")
	}
	if periodIndex < 0 || periodIndex+1 >= len(weights.Left) {
		return nil, errors.New("UPM period index is outside the edge-weight arrays")
	}
	for weightIndex := periodIndex; weightIndex <= periodIndex+1; weightIndex++ {
		left := weights.Left[weightIndex]
		current := weights.Current[weightIndex]
		right := weights.Right[weightIndex]
		if left < 0 || left > 32767 || current < 0 || current > 32767 || right < 0 || right > 32767 {
			return nil, fmt.Errorf("UPM edge weights at index %d are outside the signed-short range", weightIndex)
		}
		if left+current+right != weights.Normalization {
			return nil, fmt.Errorf("UPM edge weights at index %d do not sum to normalization", weightIndex)
		}
	}

	contributions := make([]WeightedWindowContribution, 0, 6)
	for edgeIndex, alignment := range [...]WindowAlignment{WindowAlignLeft, WindowAlignRight} {
		weightIndex := periodIndex + edgeIndex
		for _, side := range [...]struct {
			name    string
			samples []int16
			weights []int32
		}{
			{name: "current", samples: windows.Current, weights: weights.Current},
			{name: "left", samples: windows.Left, weights: weights.Left},
			{name: "right", samples: windows.Right, weights: weights.Right},
		} {
			weight := side.weights[weightIndex]
			if weight == 0 {
				continue
			}
			prepared, err := ApplyPaul2013BlendWindow(side.samples, intervalLength, alignment)
			if err != nil {
				return nil, fmt.Errorf("prepare UPM %s window: %w", side.name, err)
			}
			contributions = append(contributions, WeightedWindowContribution{
				Samples: prepared,
				Offset:  outputOffset,
				Weight:  weight,
			})
		}
	}
	return contributions, nil
}

// AddPaul2013WeightedWindows applies the integer accumulation used by the
// Stage 8 context mixer to prepared windows. Each contribution is multiplied
// by its edge weight, divided by normalization with signed integer truncation,
// clipped to [-32767, 32766], then added with the DLL's 16-bit accumulator
// behavior. Window selection and coefficient preparation remain caller work.
func AddPaul2013WeightedWindows(
	output []int16,
	normalization int32,
	contributions []WeightedWindowContribution,
) error {
	if len(output) == 0 {
		return errors.New("weighted-window output is empty")
	}
	if normalization < 1 || normalization > 32767 {
		return errors.New("weighted-window normalization must be in the range 1 through 32767")
	}
	for contributionIndex, contribution := range contributions {
		if contribution.Weight < 0 || contribution.Weight > 32767 {
			return fmt.Errorf("weighted-window contribution %d has invalid weight %d", contributionIndex, contribution.Weight)
		}
		if contribution.Offset < 0 || contribution.Offset > len(output) ||
			len(contribution.Samples) > len(output)-contribution.Offset {
			return fmt.Errorf("weighted-window contribution %d falls outside the output", contributionIndex)
		}
	}
	for _, contribution := range contributions {
		if contribution.Weight == 0 {
			continue
		}
		for sampleIndex, sample := range contribution.Samples {
			value := int64(sample) * int64(contribution.Weight) / int64(normalization)
			if value > 32766 {
				value = 32766
			} else if value < -32767 {
				value = -32767
			}
			outputIndex := contribution.Offset + sampleIndex
			output[outputIndex] = int16(int32(output[outputIndex]) + int32(value))
		}
	}
	return nil
}
