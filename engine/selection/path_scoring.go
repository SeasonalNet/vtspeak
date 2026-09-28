package selection

import (
	"errors"
	"fmt"
	"math"
)

// TransitionCandidate is one current-context unit admitted by candidate
// expansion and any preceding local-score filtering.
type TransitionCandidate struct {
	Unit     UnitRef
	NodeFlag uint16
}

// TransitionCandidateScorer returns the cumulative score for one current and
// previous candidate pair. The callback can use ScoreTransition and receives
// the predecessor's cumulative cost explicitly.
type TransitionCandidateScorer func(current, previous UnitRef, previousCumulative float32) (float32, error)

// InitializePathLayer creates the first path layer from caller-computed local
// candidate costs. Initial predecessor links use -1 because no preceding
// layer exists; Backtrack does not follow a predecessor from this layer.
func InitializePathLayer(candidates []ScoredUnit) ([]PathCandidate, error) {
	if len(candidates) == 0 {
		return nil, errors.New("initial path layer requires at least one candidate")
	}
	result := make([]PathCandidate, len(candidates))
	for index, candidate := range candidates {
		if math.IsNaN(float64(candidate.Cost)) || math.IsInf(float64(candidate.Cost), 0) {
			return nil, fmt.Errorf("initial candidate %d has non-finite local cost", index)
		}
		result[index] = PathCandidate{
			Unit:           candidate.Unit,
			CumulativeCost: candidate.Cost,
			PreviousIndex:  -1,
			NodeFlag:       candidate.NodeFlag,
		}
	}
	return result, nil
}

// ScoreTransitionLayer finds the minimum-cost predecessor for every current
// candidate and records its index for Backtrack. Equal costs keep the first
// predecessor, matching the observed minimum scan order.
func ScoreTransitionLayer(previous []PathCandidate, current []TransitionCandidate, score TransitionCandidateScorer) ([]PathCandidate, error) {
	if len(previous) == 0 {
		return nil, errors.New("transition scoring requires a nonempty previous layer")
	}
	if len(current) == 0 {
		return nil, errors.New("transition scoring requires a nonempty current layer")
	}
	if score == nil {
		return nil, errors.New("transition scoring requires a score function")
	}

	result := make([]PathCandidate, len(current))
	for currentIndex, candidate := range current {
		bestCost := float32(math.Inf(1))
		bestPrevious := -1
		for previousIndex, predecessor := range previous {
			if math.IsNaN(float64(predecessor.CumulativeCost)) || math.IsInf(float64(predecessor.CumulativeCost), 0) {
				return nil, fmt.Errorf("previous candidate %d has non-finite cumulative cost", previousIndex)
			}
			cost, err := score(candidate.Unit, predecessor.Unit, predecessor.CumulativeCost)
			if err != nil {
				return nil, fmt.Errorf("score current candidate %d against predecessor %d: %w", currentIndex, previousIndex, err)
			}
			if math.IsNaN(float64(cost)) || math.IsInf(float64(cost), 0) {
				return nil, fmt.Errorf("score for current candidate %d and predecessor %d is not finite", currentIndex, previousIndex)
			}
			if bestPrevious < 0 || cost < bestCost {
				bestCost = cost
				bestPrevious = previousIndex
			}
		}
		result[currentIndex] = PathCandidate{
			Unit:           candidate.Unit,
			CumulativeCost: bestCost,
			PreviousIndex:  bestPrevious,
			NodeFlag:       candidate.NodeFlag,
		}
	}
	return result, nil
}

// ScoreAndPruneTransitionLayer runs the predecessor minimum scan and then the
// observed cumulative-cost pruning rule for the resulting current layer.
func ScoreAndPruneTransitionLayer(previous []PathCandidate, current []TransitionCandidate, score TransitionCandidateScorer, contextGateClear bool, multiplier float32) ([]PathCandidate, error) {
	scored, err := ScoreTransitionLayer(previous, current, score)
	if err != nil {
		return nil, err
	}
	return PrunePathCandidates(scored, contextGateClear, multiplier)
}
