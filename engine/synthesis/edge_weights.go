package synthesis

import (
	"errors"
)

// UPMEdgeWeights contains the integer edge ramps built by FUN_1002d230.
// Arrays contain periodCount+1 entries so the renderer can use both endpoints
// of each adjacent-period interval.
type UPMEdgeWeights struct {
	Normalization int32
	Left          []int32
	Current       []int32
	Right         []int32
}

// LimitPaul2013ContextCount applies the observed period-count and five-entry
// limits to one gathered neighbor count. The count and period count come from
// byte-sized runtime fields.
func LimitPaul2013ContextCount(requestedCount, periodCount int) (int, error) {
	if requestedCount < 0 || requestedCount > 255 {
		return 0, errors.New("requested context count must fit in an unsigned byte")
	}
	if periodCount < 1 || periodCount > 255 {
		return 0, errors.New("UPM period count must be in the range 1 through 255")
	}
	if requestedCount > periodCount {
		requestedCount = periodCount
	}
	if requestedCount > 5 {
		requestedCount = 5
	}
	return requestedCount, nil
}

// BuildPaul2013UPMEdgeWeights constructs the traced integer ramps for one
// current unit. It first limits each gathered neighbor count by the period
// count and five-entry cap; multipliers come from
// Paul2013ContextMultipliers. It uses normalization L*R, left[k]=(leftCount-k)*R,
// right[periodCount-k]=(rightCount-k)*L, and current=normalization-left-right;
// a suppressed side contributes a scale of one. This does not interpolate or
// blend audio.
func BuildPaul2013UPMEdgeWeights(
	periodCount, leftContextCount, leftMultiplier, rightContextCount, rightMultiplier int,
) (UPMEdgeWeights, error) {
	var err error
	leftContextCount, err = LimitPaul2013ContextCount(leftContextCount, periodCount)
	if err != nil {
		return UPMEdgeWeights{}, err
	}
	rightContextCount, err = LimitPaul2013ContextCount(rightContextCount, periodCount)
	if err != nil {
		return UPMEdgeWeights{}, err
	}
	if leftMultiplier < 0 || leftMultiplier > 2 || rightMultiplier < 0 || rightMultiplier > 2 {
		return UPMEdgeWeights{}, errors.New("context multipliers must be zero, one, or two")
	}
	if (leftMultiplier == 0 && leftContextCount != 0) ||
		(rightMultiplier == 0 && rightContextCount != 0) ||
		(leftMultiplier != 0 && leftContextCount == 0) ||
		(rightMultiplier != 0 && rightContextCount == 0) {
		return UPMEdgeWeights{}, errors.New("each suppressed side must have zero contexts and each active side must have contexts")
	}

	leftScale := int32(1)
	if leftMultiplier != 0 {
		leftScale = int32(leftContextCount * leftMultiplier)
	}
	rightScale := int32(1)
	if rightMultiplier != 0 {
		rightScale = int32(rightContextCount * rightMultiplier)
	}
	normalization := leftScale * rightScale
	weights := UPMEdgeWeights{
		Normalization: normalization,
		Left:          make([]int32, periodCount+1),
		Current:       make([]int32, periodCount+1),
		Right:         make([]int32, periodCount+1),
	}
	for index := 0; index <= periodCount; index++ {
		if index < leftContextCount {
			weights.Left[index] = int32(leftContextCount-index) * rightScale
		}
		if index < rightContextCount {
			weights.Right[periodCount-index] = int32(rightContextCount-index) * leftScale
		}
	}
	for index := range weights.Current {
		weights.Current[index] = normalization - weights.Left[index] - weights.Right[index]
	}
	return weights, nil
}
