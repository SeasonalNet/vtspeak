package selection

import (
	"errors"
	"fmt"
	"math"

	"vtspeak/engine/dat"
	"vtspeak/engine/distance"
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

type indexedTransitionCandidateScorer func(currentIndex int, current TransitionCandidate, previous PathCandidate) (float32, error)

// Paul2013PathEdgeScorer returns the cumulative score for one current unit
// and one predecessor. It receives the current candidate's local cost and
// flag so the callback can build the complete FUN_10018c80 score input.
type Paul2013PathEdgeScorer func(position int, current ScoredUnit, previous PathCandidate) (float32, error)

// Paul2013PathLayerOptions contains the context gate that controls the
// FUN_10024510 prune after one transition layer. The observed caller uses a
// pruning multiplier of 1.
type Paul2013PathLayerOptions struct {
	ContextGateClear bool
}

// Paul2013PathResult retains the transition layers and the selected minimum
// cost unit path for review or downstream timeline construction.
type Paul2013PathResult struct {
	Layers   [][]PathCandidate
	Selected []UnitRef
}

// BuildPaul2013PathCandidateLayers adapts complete native local-score passes
// to the minimum-cost path input while retaining unit order and protection
// flags. It fails closed when a candidate remains outside the recovered local
// scoring prefix because its path cost is then unknown.
func BuildPaul2013PathCandidateLayers(
	passes []Paul2013CandidateLocalScorePass,
) ([][]ScoredUnit, error) {
	if len(passes) == 0 {
		return nil, errors.New("path candidate adapter requires at least one position")
	}
	layers := make([][]ScoredUnit, len(passes))
	for position, pass := range passes {
		if len(pass.Candidates) == 0 {
			return nil, fmt.Errorf("local-score pass at position %d is empty", position)
		}
		if pass.ScoredPrefix < 0 || pass.ScoredPrefix > len(pass.Candidates) {
			return nil, fmt.Errorf("local-score pass at position %d has invalid prefix %d for %d candidates", position, pass.ScoredPrefix, len(pass.Candidates))
		}
		layers[position] = make([]ScoredUnit, len(pass.Candidates))
		for candidateIndex, candidate := range pass.Candidates {
			if !candidate.HasLocalCost {
				return nil, fmt.Errorf("candidate %d at position %d has no recovered local path cost", candidateIndex, position)
			}
			if math.IsNaN(float64(candidate.LocalCost)) || math.IsInf(float64(candidate.LocalCost), 0) {
				return nil, fmt.Errorf("candidate %d at position %d has a non-finite local path cost", candidateIndex, position)
			}
			if candidate.NodeFlag > 1 {
				return nil, fmt.Errorf("candidate %d at position %d has invalid native flag %d", candidateIndex, position, candidate.NodeFlag)
			}
			layers[position][candidateIndex] = ScoredUnit{
				Unit: candidate.Continuity.Unit, Cost: candidate.LocalCost, NodeFlag: candidate.NodeFlag,
			}
		}
	}
	return layers, nil
}

// SelectPaul2013PathFromLocalScorePasses carries complete continuity-ranked
// and local-scored candidate positions through record-backed transition
// scoring, cumulative pruning, and backtracking. Context/query producers stay
// upstream; any candidate lacking a local path cost is rejected.
func SelectPaul2013PathFromLocalScorePasses(
	model Paul2013PathModel,
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	passes []Paul2013CandidateLocalScorePass,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PathResult, error) {
	layers, err := BuildPaul2013PathCandidateLayers(passes)
	if err != nil {
		return Paul2013PathResult{}, err
	}
	return SelectPaul2013ModelMinimumCostPath(
		model, distanceTable, featureTable, layers, layerOptions, transitions,
	)
}

// Paul2013PathTransitionContext carries the explicit context-state inputs for
// one adjacent candidate-layer edge. Candidate records and cumulative costs
// are filled by SelectPaul2013ModelMinimumCostPath.
type Paul2013PathTransitionContext struct {
	RecordContext Paul2013TransitionRecordContext
	CurrentState  Paul2013TransitionContextState
	PreviousState Paul2013TransitionContextState
}

// Paul2013PathModel provides the indexed candidate records and global unit
// ordering needed by model-backed transition scoring.
type Paul2013PathModel interface {
	GlobalUnitOrdinalResolver
	ReadRecord(bank string, index uint32) (dat.UnitRecord, error)
}

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
	return scoreTransitionCandidates(previous, current, func(_ int, candidate TransitionCandidate, predecessor PathCandidate) (float32, error) {
		return score(candidate.Unit, predecessor.Unit, predecessor.CumulativeCost)
	})
}

func scoreTransitionCandidates(previous []PathCandidate, current []TransitionCandidate, score indexedTransitionCandidateScorer) ([]PathCandidate, error) {
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
			cost, err := score(currentIndex, candidate, predecessor)
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

// SelectPaul2013MinimumCostPath composes the observed multi-position search:
// initialize the first layer from local costs, score every later candidate
// against every retained predecessor, apply the native cumulative-cost prune
// with multiplier 1, then backtrack from the final layer. Candidate expansion,
// local-cost production, context gates, and transition-state production stay
// explicit inputs because their producers are not yet connected to text.
func SelectPaul2013MinimumCostPath(
	candidateLayers [][]ScoredUnit,
	layerOptions []Paul2013PathLayerOptions,
	score Paul2013PathEdgeScorer,
) (Paul2013PathResult, error) {
	if len(candidateLayers) == 0 {
		return Paul2013PathResult{}, errors.New("minimum-cost path requires at least one candidate layer")
	}
	if len(layerOptions) != len(candidateLayers)-1 {
		return Paul2013PathResult{}, fmt.Errorf(
			"received %d transition-layer options for %d later candidate layers",
			len(layerOptions), len(candidateLayers)-1,
		)
	}
	if len(candidateLayers) > 1 && score == nil {
		return Paul2013PathResult{}, errors.New("minimum-cost path requires an edge scorer for multiple layers")
	}

	first, err := InitializePathLayer(candidateLayers[0])
	if err != nil {
		return Paul2013PathResult{}, fmt.Errorf("initialize candidate layer 0: %w", err)
	}
	layers := make([][]PathCandidate, len(candidateLayers))
	layers[0] = first
	for position := 1; position < len(candidateLayers); position++ {
		currentScored := candidateLayers[position]
		current := make([]TransitionCandidate, len(currentScored))
		for index, candidate := range currentScored {
			current[index] = TransitionCandidate{Unit: candidate.Unit, NodeFlag: candidate.NodeFlag}
		}
		scored, err := scoreTransitionCandidates(layers[position-1], current,
			func(currentIndex int, _ TransitionCandidate, previous PathCandidate) (float32, error) {
				return score(position, currentScored[currentIndex], previous)
			},
		)
		if err != nil {
			return Paul2013PathResult{}, fmt.Errorf("score candidate layer %d: %w", position, err)
		}
		layers[position], err = PrunePathCandidates(scored, layerOptions[position-1].ContextGateClear, 1)
		if err != nil {
			return Paul2013PathResult{}, fmt.Errorf("prune candidate layer %d: %w", position, err)
		}
	}
	selected, err := Backtrack(layers)
	if err != nil {
		return Paul2013PathResult{}, err
	}
	return Paul2013PathResult{Layers: layers, Selected: selected}, nil
}

// SelectPaul2013ModelMinimumCostPath connects model record reads, mode-specific
// transition-input construction, transition scoring, cumulative pruning, and
// backtracking. Candidate Cost values are raw local scores from the native
// local scorer; this path applies the observed /2 normalization to each
// candidate before adding it to cumulative path cost. Per-edge state records
// still come from caller-side context production.
func SelectPaul2013ModelMinimumCostPath(
	model Paul2013PathModel,
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	candidateLayers [][]ScoredUnit,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PathResult, error) {
	if len(candidateLayers) == 0 {
		return Paul2013PathResult{}, errors.New("model-backed minimum-cost path requires at least one candidate layer")
	}
	if len(transitions) != len(candidateLayers)-1 {
		return Paul2013PathResult{}, fmt.Errorf("received %d transition contexts for %d edges", len(transitions), len(candidateLayers)-1)
	}
	if len(candidateLayers) > 1 {
		if model == nil {
			return Paul2013PathResult{}, errors.New("model-backed path has no Paul model")
		}
		if distanceTable == nil || featureTable == nil {
			return Paul2013PathResult{}, errors.New("model-backed path requires loaded transition distance tables")
		}
	}
	for edgeIndex, transition := range transitions {
		if transition.RecordContext.Mode > 2 {
			return Paul2013PathResult{}, fmt.Errorf("transition edge %d has unsupported context mode %d", edgeIndex, transition.RecordContext.Mode)
		}
	}
	if len(candidateLayers) == 1 {
		return SelectPaul2013MinimumCostPath(normalizePaul2013PathLocalCosts(candidateLayers), layerOptions, nil)
	}

	records := make(map[UnitRef]dat.UnitRecord)
	readRecord := func(unit UnitRef) (dat.UnitRecord, error) {
		if record, found := records[unit]; found {
			return record, nil
		}
		record, err := model.ReadRecord(unit.Bank, unit.Index)
		if err != nil {
			return dat.UnitRecord{}, err
		}
		records[unit] = record
		return record, nil
	}

	score := func(position int, current ScoredUnit, previous PathCandidate) (float32, error) {
		currentRecord, err := readRecord(current.Unit)
		if err != nil {
			return 0, fmt.Errorf("read current unit %s:%d: %w", current.Unit.Bank, current.Unit.Index, err)
		}
		previousRecord, err := readRecord(previous.Unit)
		if err != nil {
			return 0, fmt.Errorf("read previous unit %s:%d: %w", previous.Unit.Bank, previous.Unit.Index, err)
		}
		transition := transitions[position-1]
		context := transition.RecordContext
		// Candidate layers passed into the generic path scorer have already
		// been normalized by the float32 constant 2.0 at DLL VA 0x1007c284.
		context.DurationTerm = current.Cost
		context.PreviousCumulative = previous.CumulativeCost
		input, err := BuildPaul2013TransitionScoreInput(
			current.Unit, currentRecord, previous.Unit, previousRecord, context,
			transition.CurrentState, transition.PreviousState, model,
		)
		if err != nil {
			return 0, fmt.Errorf("build transition input at edge %d: %w", position-1, err)
		}
		return ScoreTransition(distanceTable, featureTable, input)
	}
	return SelectPaul2013MinimumCostPath(normalizePaul2013PathLocalCosts(candidateLayers), layerOptions, score)
}

func normalizePaul2013PathLocalCosts(candidateLayers [][]ScoredUnit) [][]ScoredUnit {
	normalized := make([][]ScoredUnit, len(candidateLayers))
	for layerIndex, layer := range candidateLayers {
		normalized[layerIndex] = make([]ScoredUnit, len(layer))
		for candidateIndex, candidate := range layer {
			normalized[layerIndex][candidateIndex] = candidate
			normalized[layerIndex][candidateIndex].Cost = candidate.Cost / 2
		}
	}
	return normalized
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
