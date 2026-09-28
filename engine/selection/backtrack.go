package selection

import (
	"errors"
	"fmt"
	"math"
)

// PathCandidate is one unit in a context's transition-scored candidate list.
// PreviousIndex points to its selected predecessor in the preceding context.
type PathCandidate struct {
	Unit           UnitRef
	CumulativeCost float32
	PreviousIndex  int
	NodeFlag       uint16
}

// Backtrack chooses the first minimum-cost candidate in the final context and
// follows saved predecessor indices to produce one unit per context. Equal
// costs retain the earlier candidate, matching the observed minimum scan.
func Backtrack(layers [][]PathCandidate) ([]UnitRef, error) {
	if len(layers) == 0 {
		return nil, errors.New("backtracking requires at least one context")
	}
	last := layers[len(layers)-1]
	if len(last) == 0 {
		return nil, fmt.Errorf("context %d has no candidates", len(layers)-1)
	}
	selected := 0
	minimum := last[0].CumulativeCost
	if math.IsNaN(float64(minimum)) || math.IsInf(float64(minimum), 0) {
		return nil, fmt.Errorf("final candidate 0 has non-finite cumulative cost")
	}
	for index := 1; index < len(last); index++ {
		cost := last[index].CumulativeCost
		if math.IsNaN(float64(cost)) || math.IsInf(float64(cost), 0) {
			return nil, fmt.Errorf("final candidate %d has non-finite cumulative cost", index)
		}
		if cost < minimum {
			selected = index
			minimum = cost
		}
	}
	path := make([]UnitRef, len(layers))
	for layerIndex := len(layers) - 1; layerIndex >= 0; layerIndex-- {
		layer := layers[layerIndex]
		if len(layer) == 0 {
			return nil, fmt.Errorf("context %d has no candidates", layerIndex)
		}
		if selected < 0 || selected >= len(layer) {
			return nil, fmt.Errorf("selected candidate %d is outside context %d range", selected, layerIndex)
		}
		candidate := layer[selected]
		path[layerIndex] = candidate.Unit
		if layerIndex == 0 {
			continue
		}
		previous := candidate.PreviousIndex
		if previous < 0 || previous >= len(layers[layerIndex-1]) {
			return nil, fmt.Errorf("candidate %d in context %d points to invalid predecessor %d", selected, layerIndex, previous)
		}
		selected = previous
	}
	return path, nil
}
