package selection

import (
	"errors"
	"fmt"
	"math"
)

// ScoredUnit carries the cumulative path cost and the observed 16-bit node
// flag consumed by the legacy candidate-pruning step.
type ScoredUnit struct {
	Unit     UnitRef
	Cost     float32
	NodeFlag uint16
}

// PruneCandidates applies the observed post-transition filter. It preserves
// source order, skips pruning when the context gate is set or the candidate
// count is at most ten, and retains every nonzero-flag candidate regardless
// of cost. The measured caller uses multiplier 1.
func PruneCandidates(candidates []ScoredUnit, contextGateClear bool, multiplier float32) ([]ScoredUnit, error) {
	indexes, err := retainedCandidateIndexes(candidates, contextGateClear, multiplier)
	if err != nil {
		return nil, err
	}
	retained := make([]ScoredUnit, 0, len(indexes))
	for _, index := range indexes {
		retained = append(retained, candidates[index])
	}
	return retained, nil
}

// PrunePathCandidates applies the same cutoff while preserving transition
// predecessor indexes and node flags for later layers and Backtrack.
func PrunePathCandidates(candidates []PathCandidate, contextGateClear bool, multiplier float32) ([]PathCandidate, error) {
	scored := make([]ScoredUnit, len(candidates))
	for index, candidate := range candidates {
		scored[index] = ScoredUnit{Unit: candidate.Unit, Cost: candidate.CumulativeCost, NodeFlag: candidate.NodeFlag}
	}
	indexes, err := retainedCandidateIndexes(scored, contextGateClear, multiplier)
	if err != nil {
		return nil, err
	}
	retained := make([]PathCandidate, 0, len(indexes))
	for _, index := range indexes {
		retained = append(retained, candidates[index])
	}
	return retained, nil
}

func retainedCandidateIndexes(candidates []ScoredUnit, contextGateClear bool, multiplier float32) ([]int, error) {
	if !contextGateClear || len(candidates) <= 10 {
		indexes := make([]int, len(candidates))
		for index := range candidates {
			indexes[index] = index
		}
		return indexes, nil
	}
	if math.IsNaN(float64(multiplier)) || math.IsInf(float64(multiplier), 0) || multiplier < 0 {
		return nil, errors.New("pruning multiplier must be finite and nonnegative")
	}
	var sum float32
	minimum := float32(math.Inf(1))
	for index, candidate := range candidates {
		if math.IsNaN(float64(candidate.Cost)) || math.IsInf(float64(candidate.Cost), 0) {
			return nil, fmt.Errorf("candidate %d has non-finite cumulative cost", index)
		}
		sum += candidate.Cost
		if math.IsInf(float64(sum), 0) {
			return nil, errors.New("cumulative cost sum overflowed float32")
		}
		if candidate.Cost < minimum {
			minimum = candidate.Cost
		}
	}
	average := sum / float32(len(candidates))
	quarterSpread := (average - minimum) / 4
	cutoff := multiplier*quarterSpread + average
	if math.IsNaN(float64(cutoff)) || math.IsInf(float64(cutoff), 0) {
		return nil, errors.New("pruning cutoff is not finite")
	}
	retained := make([]int, 0, len(candidates))
	for index, candidate := range candidates {
		if candidate.Cost <= cutoff || candidate.NodeFlag != 0 {
			retained = append(retained, index)
		}
	}
	return retained, nil
}
