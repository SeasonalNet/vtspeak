package selection

import (
	"fmt"

	"vtspeak/engine/distance"
	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

// Paul2013PositionQueryResult combines the whole-position lookup and its
// conditional two-row fallback as observed in the caller of FUN_10024060.
// ReturnCount is the native row advance: one for an accepted whole-position
// result, two when FUN_100242a0 supplies fallback rows.
type Paul2013PositionQueryResult struct {
	WholePosition Paul2013WholePositionQueryResult
	FallbackRows  [2]Paul2013FallbackRowResult
	UsedFallback  bool
	ReturnCount   int
}

// QueryPaul2013Position runs the whole-position query first and invokes the
// two-row fallback only when its native acceptance predicate fails. The
// caller still supplies state-derived fallback inputs because their
// production is outside the recovered query path.
func QueryPaul2013Position(
	catalog Paul2013PrefixClassLookup,
	target text.Context,
	contextRow [6]byte,
	priorMetric uint32,
	modelClass byte,
	contextGate func() byte,
	fallbackBaseRow [6]byte,
	fallbackContextDelta int8,
	fallbackModelRowClass byte,
	fallbackRuntime [2]Paul2013FallbackRowRuntime,
) (Paul2013PositionQueryResult, error) {
	whole, err := QueryPaul2013WholePosition(
		catalog, target, contextRow, priorMetric, modelClass, contextGate,
	)
	result := Paul2013PositionQueryResult{WholePosition: whole}
	if err != nil {
		return result, fmt.Errorf("query whole position: %w", err)
	}
	if whole.Selection.Accepted {
		result.ReturnCount = 1
		return result, nil
	}
	fallback, err := QueryPaul2013FallbackRows(
		catalog, target, fallbackBaseRow, fallbackContextDelta,
		fallbackModelRowClass, fallbackRuntime,
	)
	if err != nil {
		return result, fmt.Errorf("query fallback rows: %w", err)
	}
	result.FallbackRows = fallback
	result.UsedFallback = true
	result.ReturnCount = 2
	return result, nil
}

// ExpandPaul2013PositionQueryCandidates converts the accepted whole-position
// shortlist or the two fallback class-ID rows into per-position unit pools.
// It retains native class/member ordering and the fallback's two-row advance.
// maxUnits is supplied because the unit-array capacity is a caller-state value
// not established by the query results. Missing fallback class records fail
// closed rather than dropping IDs or inventing memberships.
func ExpandPaul2013PositionQueryCandidates(
	result Paul2013PositionQueryResult,
	maxUnits uint64,
) ([][]UnitRef, error) {
	if result.ReturnCount == 1 && !result.UsedFallback && result.WholePosition.Selection.Accepted {
		units, err := expandPaul2013QueryClasses(result.WholePosition.Selection.Classes, maxUnits)
		if err != nil {
			return nil, fmt.Errorf("expand accepted whole-position classes: %w", err)
		}
		if len(units) == 0 {
			return nil, fmt.Errorf("accepted whole-position query has no unit candidates")
		}
		return [][]UnitRef{units}, nil
	}
	if result.ReturnCount != 2 || !result.UsedFallback || result.WholePosition.Selection.Accepted {
		return nil, fmt.Errorf("position query has inconsistent return count, fallback, and acceptance state")
	}

	positions := make([][]UnitRef, 2)
	for rowIndex, row := range result.FallbackRows {
		classesByID := make(map[uint32]voice.ClassRecord)
		for _, class := range row.RankedClasses {
			if err := addPaul2013CandidateClass(classesByID, class); err != nil {
				return nil, fmt.Errorf("fallback row %d ranked classes: %w", rowIndex, err)
			}
		}
		for passIndex, pass := range row.QueryPasses {
			for _, class := range pass.Candidates {
				if err := addPaul2013CandidateClass(classesByID, class); err != nil {
					return nil, fmt.Errorf("fallback row %d query pass %d: %w", rowIndex, passIndex, err)
				}
			}
		}
		ordered := make([]voice.ClassRecord, len(row.Candidates.IDs))
		for candidateIndex, classID := range row.Candidates.IDs {
			class, ok := classesByID[classID]
			if !ok {
				return nil, fmt.Errorf("fallback row %d candidate %d references class %d without a source record", rowIndex, candidateIndex, classID)
			}
			ordered[candidateIndex] = class
		}
		units, err := expandPaul2013QueryClasses(ordered, maxUnits)
		if err != nil {
			return nil, fmt.Errorf("expand fallback row %d classes: %w", rowIndex, err)
		}
		if len(units) == 0 {
			return nil, fmt.Errorf("fallback row %d has no unit candidates", rowIndex)
		}
		positions[rowIndex] = units
	}
	return positions, nil
}

// RankPaul2013QueriedCandidatePositions joins position-query expansion to the
// native continuity metadata pass. modes must be produced for every expanded
// position, including both rows of fallback positions; their producer remains
// an explicit upstream input.
func RankPaul2013QueriedCandidatePositions(
	model Paul2013ContinuityModel,
	queries []Paul2013PositionQueryResult,
	maxUnits uint64,
	modes []byte,
) ([][]Paul2013CandidateContinuity, error) {
	positions := make([][]UnitRef, 0, len(queries))
	for queryIndex, query := range queries {
		expanded, err := ExpandPaul2013PositionQueryCandidates(query, maxUnits)
		if err != nil {
			return nil, fmt.Errorf("expand position query %d: %w", queryIndex, err)
		}
		positions = append(positions, expanded...)
	}
	if len(positions) == 0 {
		return nil, fmt.Errorf("continuity ranking has no expanded query positions")
	}
	if len(modes) != len(positions) {
		return nil, fmt.Errorf("received %d context modes for %d expanded query positions", len(modes), len(positions))
	}
	ranked := make([][]Paul2013CandidateContinuity, len(positions))
	for position := range positions {
		candidates, err := RankPaul2013CandidatesByContinuity(model, positions, modes, position)
		if err != nil {
			return nil, fmt.Errorf("rank expanded query position %d: %w", position, err)
		}
		ranked[position] = candidates
	}
	return ranked, nil
}

// Paul2013QueriedCandidateScoringInput contains the state-derived local-score
// values aligned to one continuity-ranked position. Identity keys and the
// whole-position population are needed only for oversized candidate pools.
type Paul2013QueriedCandidateScoringInput struct {
	Contexts                    []Paul2013UnitScoreContext
	CandidateKeys               []uint32
	WholePositionKeys           []uint32
	WholePositionCandidateCount int
}

// Paul2013QueriedCandidateScorePass retains the ranked candidates, their
// local-score pass, and (for oversized pools) the native protected-prefix
// expansion that selected the candidates to score.
type Paul2013QueriedCandidateScorePass struct {
	Ranked        []Paul2013CandidateContinuity
	Scores        Paul2013CandidateLocalScorePass
	Protection    Paul2013ProtectionPrefixExpansion
	HasProtection bool
}

// Paul2013QueriedCandidatePathInput supplies the caller-owned state needed to
// carry queried positions through local scoring, shortlist finalization, and
// model-backed path selection.
type Paul2013QueriedCandidatePathInput struct {
	Model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	}
	TransitionDistanceTable *distance.Table
	FeatureDistanceTable    *distance.FeatureTable
	Queries                 []Paul2013PositionQueryResult
	MaxUnits                uint64
	Modes                   []byte
	PositionCount           int
	ScoringInputs           []Paul2013QueriedCandidateScoringInput
	TailStarts              []int
	CandidateClasses        [][]byte
	ContextClasses          [][]byte
	LayerOptions            []Paul2013PathLayerOptions
	Transitions             []Paul2013PathTransitionContext
}

// Paul2013QueriedCandidatePathResult retains each per-position shortlist and
// the selected minimum-cost path so both stages can be inspected independently.
type Paul2013QueriedCandidatePathResult struct {
	ScorePasses []Paul2013QueriedCandidateScorePass
	Shortlists  []Paul2013CandidateShortlist
	Path        Paul2013PathResult
}

// RankAndScorePaul2013QueriedCandidatePositions joins query expansion,
// continuity ranking, and the bounded local-score stage. State-derived score
// contexts and candidate-key mappings remain explicit inputs and must already
// align with each position's continuity-ranked order.
func RankAndScorePaul2013QueriedCandidatePositions(
	model Paul2013ContinuityModel,
	table *distance.FeatureTable,
	queries []Paul2013PositionQueryResult,
	maxUnits uint64,
	modes []byte,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
) ([]Paul2013QueriedCandidateScorePass, error) {
	ranked, err := RankPaul2013QueriedCandidatePositions(model, queries, maxUnits, modes)
	if err != nil {
		return nil, err
	}
	if len(scoringInputs) != len(ranked) {
		return nil, fmt.Errorf("received %d scoring input rows for %d ranked positions", len(scoringInputs), len(ranked))
	}

	result := make([]Paul2013QueriedCandidateScorePass, len(ranked))
	for position, candidates := range ranked {
		input := scoringInputs[position]
		pass := Paul2013QueriedCandidateScorePass{Ranked: candidates}
		if len(candidates) > 30 {
			expanded, scoreErr := ScorePaul2013ExpandedCandidateLocalPrefix(
				model, table, candidates, positionCount, input.Contexts,
				input.CandidateKeys, input.WholePositionKeys, input.WholePositionCandidateCount,
			)
			if scoreErr != nil {
				return nil, fmt.Errorf("score oversized queried position %d: %w", position, scoreErr)
			}
			pass.Scores = expanded.Scores
			pass.Protection = expanded.Protection
			pass.HasProtection = true
		} else {
			scores, scoreErr := ScorePaul2013CandidateLocalPrefix(
				model, table, candidates, positionCount, input.Contexts,
				input.CandidateKeys, input.WholePositionKeys,
			)
			if scoreErr != nil {
				return nil, fmt.Errorf("score queried position %d: %w", position, scoreErr)
			}
			pass.Scores = scores
		}
		result[position] = pass
	}
	return result, nil
}

// FinalizePaul2013QueriedCandidateShortlist applies the secondary tail score
// and native 30-candidate cap to a queried score pass. For oversized pools,
// tailStart and the class arrays are caller-produced state aligned to the
// scored candidate order. It fails closed when any candidate lacks a local
// cost because the native tail score depends on that value.
func FinalizePaul2013QueriedCandidateShortlist(
	pass Paul2013QueriedCandidateScorePass,
	tailStart int,
	candidateClasses, contextClasses []byte,
) (Paul2013CandidateShortlist, error) {
	count := len(pass.Scores.Candidates)
	if count == 0 {
		return Paul2013CandidateShortlist{}, fmt.Errorf("queried candidate score pass is empty")
	}
	if count > 30 && !pass.HasProtection {
		return Paul2013CandidateShortlist{}, fmt.Errorf("oversized queried candidate pass is missing protection metadata")
	}
	if count > 30 && (len(candidateClasses) != count || len(contextClasses) != count) {
		return Paul2013CandidateShortlist{}, fmt.Errorf("received %d candidate classes and %d context classes for %d oversized candidates", len(candidateClasses), len(contextClasses), count)
	}

	candidates := make([]Paul2013ScoredCandidate, count)
	totalProtected := 0
	for index, scored := range pass.Scores.Candidates {
		if !scored.HasLocalCost {
			return Paul2013CandidateShortlist{}, fmt.Errorf("candidate %d has no local cost for final shortlist ordering", index)
		}
		secondary := scored.LocalCost
		if count > 30 {
			var err error
			secondary, err = Paul2013TailOrderingScore(
				scored.LocalCost, candidateClasses[index], contextClasses[index],
			)
			if err != nil {
				return Paul2013CandidateShortlist{}, fmt.Errorf("score candidate %d tail order: %w", index, err)
			}
		}
		candidates[index] = Paul2013ScoredCandidate{
			Continuity: scored.Continuity, NodeFlag: scored.NodeFlag,
			LocalCost: scored.LocalCost, SecondaryOrderScore: secondary,
		}
		if scored.NodeFlag != 0 {
			totalProtected++
		}
	}
	if pass.HasProtection {
		totalProtected = pass.Protection.ProtectedCount
	}
	return FinalizePaul2013CandidateShortlist(candidates, tailStart, totalProtected)
}

// SelectPaul2013QueriedCandidatePath composes accepted/fallback class queries,
// unit expansion, continuity ranking, bounded local scoring, native shortlist
// finalization, transition scoring, pruning, and path backtracking. Text-derived
// query results, score contexts, class keys, and per-edge state remain caller
// inputs; incomplete oversized local-score tails fail closed during shortlist
// finalization.
func SelectPaul2013QueriedCandidatePath(
	input Paul2013QueriedCandidatePathInput,
) (Paul2013QueriedCandidatePathResult, error) {
	scorePasses, err := RankAndScorePaul2013QueriedCandidatePositions(
		input.Model, input.FeatureDistanceTable, input.Queries, input.MaxUnits,
		input.Modes, input.PositionCount, input.ScoringInputs,
	)
	if err != nil {
		return Paul2013QueriedCandidatePathResult{}, err
	}
	count := len(scorePasses)
	if len(input.TailStarts) != count || len(input.CandidateClasses) != count || len(input.ContextClasses) != count {
		return Paul2013QueriedCandidatePathResult{}, fmt.Errorf(
			"received %d tail starts, %d candidate-class rows, and %d context-class rows for %d queried positions",
			len(input.TailStarts), len(input.CandidateClasses), len(input.ContextClasses), count,
		)
	}

	shortlists := make([]Paul2013CandidateShortlist, count)
	localPasses := make([]Paul2013CandidateLocalScorePass, count)
	for position, pass := range scorePasses {
		shortlist, finalizeErr := FinalizePaul2013QueriedCandidateShortlist(
			pass, input.TailStarts[position], input.CandidateClasses[position],
			input.ContextClasses[position],
		)
		if finalizeErr != nil {
			return Paul2013QueriedCandidatePathResult{}, fmt.Errorf("finalize queried position %d: %w", position, finalizeErr)
		}
		shortlists[position] = shortlist
		localPasses[position] = Paul2013CandidateLocalScorePass{
			ScoredPrefix: len(shortlist.Candidates),
			Candidates:   make([]Paul2013CandidateLocalScore, len(shortlist.Candidates)),
		}
		for candidateIndex, candidate := range shortlist.Candidates {
			localPasses[position].Candidates[candidateIndex] = Paul2013CandidateLocalScore{
				Continuity: candidate.Continuity, NodeFlag: candidate.NodeFlag,
				LocalCost: candidate.LocalCost, HasLocalCost: true,
			}
		}
	}
	path, err := SelectPaul2013PathFromLocalScorePasses(
		input.Model, input.TransitionDistanceTable, input.FeatureDistanceTable,
		localPasses, input.LayerOptions, input.Transitions,
	)
	if err != nil {
		return Paul2013QueriedCandidatePathResult{}, err
	}
	return Paul2013QueriedCandidatePathResult{ScorePasses: scorePasses, Shortlists: shortlists, Path: path}, nil
}

// SelectPaul2013QueriedCandidatePathFromTransitionStates composes accepted
// and fallback position queries with the native mode bytes carried by their
// compact six-byte transition states. The states must be in expanded query
// order: one state for an accepted whole-position query, two for fallback.
// Text-derived score contexts, tail/class rows, and edge-transition inputs
// remain explicit.
func SelectPaul2013QueriedCandidatePathFromTransitionStates(
	input Paul2013QueriedCandidatePathInput,
	states []Paul2013TransitionContextState,
) (Paul2013QueriedCandidatePathResult, error) {
	expectedStates := 0
	for queryIndex, query := range input.Queries {
		if query.ReturnCount != 1 && query.ReturnCount != 2 {
			return Paul2013QueriedCandidatePathResult{}, fmt.Errorf("position query %d has unsupported native row advance %d", queryIndex, query.ReturnCount)
		}
		expectedStates += query.ReturnCount
	}
	if len(states) != expectedStates {
		return Paul2013QueriedCandidatePathResult{}, fmt.Errorf("received %d transition states for %d expanded query positions", len(states), expectedStates)
	}
	modes, err := Paul2013ContextModesFromTransitionStates(states)
	if err != nil {
		return Paul2013QueriedCandidatePathResult{}, err
	}
	input.Modes = modes
	return SelectPaul2013QueriedCandidatePath(input)
}

func expandPaul2013QueryClasses(classes []voice.ClassRecord, maxUnits uint64) ([]UnitRef, error) {
	locations, err := ExpandClassRecords(classes, maxUnits)
	if err != nil {
		return nil, err
	}
	units := make([]UnitRef, len(locations))
	for index, location := range locations {
		units[index] = UnitRef{Bank: location.Bank, Index: location.Index}
	}
	return units, nil
}

func addPaul2013CandidateClass(classes map[uint32]voice.ClassRecord, candidate voice.ClassRecord) error {
	if previous, exists := classes[candidate.ID]; exists {
		if previous.Key != candidate.Key || !sameUnitLocations(previous.Members, candidate.Members) {
			return fmt.Errorf("class ID %d has conflicting source records", candidate.ID)
		}
		return nil
	}
	candidate.Members = append([]voice.UnitLocation(nil), candidate.Members...)
	classes[candidate.ID] = candidate
	return nil
}
