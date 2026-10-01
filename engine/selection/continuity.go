package selection

import (
	"errors"
	"fmt"
	"math"

	"vtspeak/engine/dat"
	"vtspeak/engine/distance"
)

// Paul2013ContinuityModel supplies the indexed unit order and signatures used
// by the recovered candidate-span builder. Implementations need not read PCM.
type Paul2013ContinuityModel interface {
	GlobalUnitOrdinal(bank string, index uint32) (uint32, error)
	UnitAtGlobalOrdinal(ordinal uint32) (string, uint32, error)
	ReadRecord(bank string, index uint32) (dat.UnitRecord, error)
}

// Paul2013CandidateContinuity contains the three fields written by
// FUN_100230a0: right matches, total covered span, and weighted side coverage.
type Paul2013CandidateContinuity struct {
	Unit               UnitRef
	RightMatches       uint16
	TotalSpan          uint16
	WeightedSide       uint16
	InitialOrderingKey float32
}

// RankPaul2013CandidatesByContinuity ports FUN_100230a0. candidatePositions
// is the already-expanded candidate set at each phone position; modes has one
// 0, 1, or 2 context mode per position. The returned candidates at position
// are sorted by ascending initial key with FUN_1001b5f0's native ordering
// path, before local candidate costs are applied by FUN_10023350.
func RankPaul2013CandidatesByContinuity(
	model Paul2013ContinuityModel,
	candidatePositions [][]UnitRef,
	modes []byte,
	position int,
) ([]Paul2013CandidateContinuity, error) {
	if model == nil {
		return nil, errors.New("candidate continuity requires a model")
	}
	if len(candidatePositions) == 0 || len(candidatePositions) != len(modes) {
		return nil, errors.New("candidate positions and context modes must have the same nonzero length")
	}
	if position < 0 || position >= len(candidatePositions) {
		return nil, fmt.Errorf("candidate continuity position %d is outside 0..%d", position, len(candidatePositions)-1)
	}
	for index, mode := range modes {
		if mode > 2 {
			return nil, fmt.Errorf("context mode %d at position %d is outside 0..2", mode, index)
		}
	}

	ordinals := make([][]uint32, len(candidatePositions))
	sets := make([]map[uint32]struct{}, len(candidatePositions))
	for pos, candidates := range candidatePositions {
		ordinals[pos] = make([]uint32, 0, len(candidates))
		sets[pos] = make(map[uint32]struct{}, len(candidates))
		for index, candidate := range candidates {
			ordinal, err := model.GlobalUnitOrdinal(candidate.Bank, candidate.Index)
			if err != nil {
				return nil, fmt.Errorf("resolve candidate at position %d index %d: %w", pos, index, err)
			}
			ordinals[pos] = append(ordinals[pos], ordinal)
			sets[pos][ordinal] = struct{}{}
		}
	}

	records := make(map[uint32]dat.UnitRecord)
	readOrdinal := func(ordinal uint32) (dat.UnitRecord, error) {
		if record, ok := records[ordinal]; ok {
			return record, nil
		}
		bank, index, err := model.UnitAtGlobalOrdinal(ordinal)
		if err != nil {
			return dat.UnitRecord{}, err
		}
		record, err := model.ReadRecord(bank, index)
		if err != nil {
			return dat.UnitRecord{}, err
		}
		records[ordinal] = record
		return record, nil
	}

	result := make([]Paul2013CandidateContinuity, 0, len(candidatePositions[position]))
	for candidateIndex, candidate := range candidatePositions[position] {
		currentOrdinal := ordinals[position][candidateIndex]
		leftMatches, rightMatches := 0, 0
		leftWeight, rightWeight := 0, 0
		traversed := currentOrdinal

		for pos := position - 1; pos >= 0; pos-- {
			if modes[pos] != 1 {
				if traversed == 0 {
					break
				}
				traversed--
				record, err := readOrdinal(traversed)
				if err != nil {
					return nil, fmt.Errorf("read left continuity signature at ordinal %d: %w", traversed, err)
				}
				if record.Signature[6]&0x80 == 0 {
					break
				}
			}
			if _, found := sets[pos][traversed]; !found {
				break
			}
			leftMatches++
			leftWeight += paul2013ContinuityModeWeight(modes[pos])
		}

		traversed = currentOrdinal
		for pos := position + 1; pos < len(candidatePositions); pos++ {
			if modes[pos] != 2 {
				record, err := readOrdinal(traversed)
				if err != nil {
					return nil, fmt.Errorf("read right continuity signature at ordinal %d: %w", traversed, err)
				}
				if record.Signature[6]&0x80 == 0 || traversed == ^uint32(0) {
					break
				}
				traversed++
			}
			if _, found := sets[pos][traversed]; !found {
				break
			}
			rightMatches++
			rightWeight += paul2013ContinuityModeWeight(modes[pos])
		}

		totalSpan := leftMatches + 1 + rightMatches
		leftSide, rightSide := leftWeight/2, rightWeight/2
		weightedSide := 0
		switch {
		case totalSpan == len(candidatePositions):
			weightedSide = rightSide + 1 + leftSide
		case leftMatches == position:
			weightedSide = rightSide
		case rightMatches == len(candidatePositions)-position-1:
			weightedSide = leftSide
		case leftSide < rightSide:
			weightedSide = leftSide
		default:
			weightedSide = rightSide
		}
		key := -float32(totalSpan + weightedSide*100)
		result = append(result, Paul2013CandidateContinuity{
			Unit: candidate, RightMatches: uint16(rightMatches), TotalSpan: uint16(totalSpan),
			WeightedSide: uint16(weightedSide), InitialOrderingKey: key,
		})
	}

	sortPaul2013Native(result, func(left, right Paul2013CandidateContinuity) bool {
		return left.InitialOrderingKey < right.InitialOrderingKey
	})
	return result, nil
}

// Paul2013CandidateMembershipFlags marks each supplied candidate key found
// in the whole-position key set. FUN_10023350 supplies only its local-score
// prefix to this pass; mapping unit references into these key namespaces is
// caller work.
func Paul2013CandidateMembershipFlags(candidateKeys, wholePositionKeys []uint32) ([]uint16, error) {
	if len(candidateKeys) > 10000 {
		return nil, fmt.Errorf("candidate membership pass exceeds the native 10000-candidate limit: %d", len(candidateKeys))
	}
	wholePositionSet := make(map[uint32]struct{}, len(wholePositionKeys))
	for _, key := range wholePositionKeys {
		wholePositionSet[key] = struct{}{}
	}
	flags := make([]uint16, len(candidateKeys))
	for index, key := range candidateKeys {
		if _, protected := wholePositionSet[key]; protected {
			flags[index] = 1
		}
	}
	return flags, nil
}

// Paul2013ProtectionPrefixExpansion is the candidate list after the native
// whole-position protection scan. ScoredPrefixSources maps each output prefix
// entry back to its original candidate index, including entries copied from
// the tail. NodeFlags and CandidateKeys align with Candidates after those
// copies.
type Paul2013ProtectionPrefixExpansion struct {
	Candidates          []Paul2013CandidateContinuity
	CandidateKeys       []uint32
	NodeFlags           []uint16
	ScoredPrefixSources []int
	ScoredPrefix        int
	ProtectedCount      int
	Added               int
}

// ExpandPaul2013CandidateProtectionPrefix ports FUN_10023350's protection
// scan after it identifies the initial local-score prefix. It marks protected
// candidates in that prefix, then, when the whole-position pool contains more
// than nine units and fewer than ten prefix candidates are protected, copies
// matching tail candidates into the prefix. After the tenth protected
// candidate, scanning stops at the first different initial ordering key.
// Candidate-to-key mapping and the whole-position pool count are caller inputs.
func ExpandPaul2013CandidateProtectionPrefix(
	candidates []Paul2013CandidateContinuity,
	candidateKeys, wholePositionKeys []uint32,
	initialPrefix, wholePositionCandidateCount int,
) (Paul2013ProtectionPrefixExpansion, error) {
	if len(candidates) <= 30 || len(candidates) > 10000 {
		return Paul2013ProtectionPrefixExpansion{}, fmt.Errorf("candidate protection expansion requires 31 through 10000 candidates, got %d", len(candidates))
	}
	if len(candidateKeys) != len(candidates) {
		return Paul2013ProtectionPrefixExpansion{}, fmt.Errorf("received %d candidate keys for %d candidates", len(candidateKeys), len(candidates))
	}
	if initialPrefix < 0 || initialPrefix > len(candidates) {
		return Paul2013ProtectionPrefixExpansion{}, fmt.Errorf("initial scored prefix %d is outside 0..%d", initialPrefix, len(candidates))
	}
	if wholePositionCandidateCount < 0 || wholePositionCandidateCount > 10000 {
		return Paul2013ProtectionPrefixExpansion{}, fmt.Errorf("whole-position candidate count %d is outside 0..10000", wholePositionCandidateCount)
	}

	wholePositionSet := make(map[uint32]struct{}, len(wholePositionKeys))
	for _, key := range wholePositionKeys {
		wholePositionSet[key] = struct{}{}
	}
	result := Paul2013ProtectionPrefixExpansion{
		Candidates:    append([]Paul2013CandidateContinuity(nil), candidates...),
		CandidateKeys: append([]uint32(nil), candidateKeys...),
		NodeFlags:     make([]uint16, len(candidates)),
		ScoredPrefix:  initialPrefix,
	}
	for index := 0; index < initialPrefix; index++ {
		if _, protected := wholePositionSet[candidateKeys[index]]; protected {
			result.NodeFlags[index] = 1
			result.ProtectedCount++
		}
		result.ScoredPrefixSources = append(result.ScoredPrefixSources, index)
	}
	if result.ProtectedCount > 9 || wholePositionCandidateCount <= 9 || initialPrefix == len(candidates) {
		return result, nil
	}

	boundaryKey := float32(0)
	originalKeys := append([]uint32(nil), candidateKeys...)
	for sourceIndex := initialPrefix; sourceIndex < len(candidates); sourceIndex++ {
		source := candidates[sourceIndex]
		if result.ProtectedCount > 9 && source.InitialOrderingKey != boundaryKey {
			break
		}
		if _, protected := wholePositionSet[originalKeys[sourceIndex]]; !protected {
			continue
		}

		destinationIndex := result.ScoredPrefix
		destination := result.Candidates[destinationIndex]
		destination.Unit = source.Unit
		destination.RightMatches = source.RightMatches
		destination.TotalSpan = source.TotalSpan
		destination.WeightedSide = source.WeightedSide
		// The native copy omits node score and ordering-key fields, so retain
		// those values already stored in the destination slot.
		result.Candidates[destinationIndex] = destination
		result.CandidateKeys[destinationIndex] = originalKeys[sourceIndex]
		result.NodeFlags[destinationIndex] = 1
		result.ScoredPrefixSources = append(result.ScoredPrefixSources, sourceIndex)
		result.ScoredPrefix++
		result.Added++
		result.ProtectedCount++
		if result.ProtectedCount == 10 {
			boundaryKey = source.InitialOrderingKey
		}
	}
	return result, nil
}

func paul2013ContinuityModeWeight(mode byte) int {
	if mode == 0 {
		return 2
	}
	return 1
}

// Paul2013CandidateLocalScorePrefix returns how many continuity-ranked
// candidates FUN_10023350 sends through its local scorer. Full-span candidates
// at the head are scored as one prefix. Otherwise, sets smaller than 31 score
// in full; larger sets stop at the first strict key increase after 30 entries,
// retaining any tie at that boundary.
func Paul2013CandidateLocalScorePrefix(candidates []Paul2013CandidateContinuity, positionCount int) (int, error) {
	if positionCount < 1 {
		return 0, errors.New("candidate local-score prefix requires at least one position")
	}
	for index, candidate := range candidates {
		if candidate.TotalSpan == 0 || int(candidate.TotalSpan) > positionCount {
			return 0, fmt.Errorf("candidate %d span %d is outside 1..%d", index, candidate.TotalSpan, positionCount)
		}
		if index > 0 && candidates[index-1].InitialOrderingKey > candidate.InitialOrderingKey {
			return 0, fmt.Errorf("candidate continuity list is not ordered at index %d", index)
		}
	}
	if len(candidates) == 0 {
		return 0, nil
	}
	if int(candidates[0].TotalSpan) == positionCount {
		count := 0
		for _, candidate := range candidates {
			if int(candidate.TotalSpan) != positionCount {
				break
			}
			count++
		}
		return count, nil
	}
	if len(candidates) < 31 {
		return len(candidates), nil
	}
	count := 1
	for index := 1; index < len(candidates); index++ {
		if candidates[index-1].InitialOrderingKey < candidates[index].InitialOrderingKey && count > 29 {
			return count, nil
		}
		count++
	}
	return count, nil
}

// Paul2013ScoredCandidate carries the native local score and the secondary
// score used to order the tail of an oversized shortlist.
type Paul2013ScoredCandidate struct {
	Continuity          Paul2013CandidateContinuity
	NodeFlag            uint16
	LocalCost           float32
	SecondaryOrderScore float32
}

// Paul2013ProtectedCandidateReplacement records a flagged candidate restored
// from beyond the 30-entry ranked prefix and the shortlist slot it replaces.
type Paul2013ProtectedCandidateReplacement struct {
	Slot      int
	Candidate Paul2013ScoredCandidate
}

// Paul2013CandidateShortlist is the final native shortlist and the protected
// tail candidates restored into it after sorting and truncation.
type Paul2013CandidateShortlist struct {
	Candidates            []Paul2013ScoredCandidate
	ProtectedReplacements []Paul2013ProtectedCandidateReplacement
}

// Paul2013CandidateLocalScore records one continuity-ranked candidate after
// FUN_10023350's bounded local-scoring pass. LocalCost is meaningful only when
// HasLocalCost is true; candidates outside the native scoring prefix retain
// their continuity order and membership flag for the later tail pass.
type Paul2013CandidateLocalScore struct {
	Continuity   Paul2013CandidateContinuity
	NodeFlag     uint16
	LocalCost    float32
	HasLocalCost bool
}

// Paul2013CandidateLocalScorePass is the scored prefix and untouched tail
// produced by FUN_10023350 before its final context-dependent tail ordering.
type Paul2013CandidateLocalScorePass struct {
	Candidates   []Paul2013CandidateLocalScore
	ScoredPrefix int
}

// Paul2013ExpandedCandidateLocalScorePass joins the large-pool protection
// expansion to record loading and local scoring, retaining the source-index
// map needed to audit context alignment.
type Paul2013ExpandedCandidateLocalScorePass struct {
	Protection Paul2013ProtectionPrefixExpansion
	Scores     Paul2013CandidateLocalScorePass
}

// ScorePaul2013CandidateLocalPrefix connects candidate-record lookup,
// FUN_100182e0 input construction, continuity scaling, and the native
// FUN_10023350 local-score prefix. candidates must already be in the
// continuity order returned by RankPaul2013CandidatesByContinuity. Contexts
// are the caller-produced state values aligned to that order. For lists above
// 30, candidateKeys and wholePositionKeys supply the two identity namespaces
// consumed by the native protection-marker pass. Tail ordering remains a
// separate stage because its secondary context scores are produced upstream.
func ScorePaul2013CandidateLocalPrefix(
	model Paul2013ContinuityModel,
	table *distance.FeatureTable,
	candidates []Paul2013CandidateContinuity,
	positionCount int,
	contexts []Paul2013UnitScoreContext,
	candidateKeys, wholePositionKeys []uint32,
) (Paul2013CandidateLocalScorePass, error) {
	if model == nil {
		return Paul2013CandidateLocalScorePass{}, errors.New("candidate local scoring requires a model")
	}
	if table == nil {
		return Paul2013CandidateLocalScorePass{}, errors.New("candidate local scoring requires a feature-distance table")
	}
	if len(contexts) != len(candidates) {
		return Paul2013CandidateLocalScorePass{}, fmt.Errorf("received %d score contexts for %d candidates", len(contexts), len(candidates))
	}
	if len(candidates) > 30 && len(candidateKeys) != len(candidates) {
		return Paul2013CandidateLocalScorePass{}, fmt.Errorf("received %d identity keys for %d oversized candidates", len(candidateKeys), len(candidates))
	}
	prefix, err := Paul2013CandidateLocalScorePrefix(candidates, positionCount)
	if err != nil {
		return Paul2013CandidateLocalScorePass{}, fmt.Errorf("determine native local-score prefix: %w", err)
	}
	flags := make([]uint16, len(candidates))
	if len(candidates) > 30 {
		flags, err = Paul2013CandidateMembershipFlags(candidateKeys[:prefix], wholePositionKeys)
		if err != nil {
			return Paul2013CandidateLocalScorePass{}, fmt.Errorf("derive candidate protection flags: %w", err)
		}
	}
	return scorePaul2013CandidatePrefix(model, table, candidates, contexts, flags, prefix)
}

// ScorePaul2013ExpandedCandidateLocalPrefix connects the native prefix
// boundary, large-whole-position protected-tail expansion, context realignment,
// candidate-record reads, and continuity-scaled local costs. The candidate
// key namespaces, whole-position pool count, and context values remain
// caller-produced inputs.
func ScorePaul2013ExpandedCandidateLocalPrefix(
	model Paul2013ContinuityModel,
	table *distance.FeatureTable,
	candidates []Paul2013CandidateContinuity,
	positionCount int,
	contexts []Paul2013UnitScoreContext,
	candidateKeys, wholePositionKeys []uint32,
	wholePositionCandidateCount int,
) (Paul2013ExpandedCandidateLocalScorePass, error) {
	if model == nil {
		return Paul2013ExpandedCandidateLocalScorePass{}, errors.New("candidate local scoring requires a model")
	}
	if table == nil {
		return Paul2013ExpandedCandidateLocalScorePass{}, errors.New("candidate local scoring requires a feature-distance table")
	}
	if len(contexts) != len(candidates) {
		return Paul2013ExpandedCandidateLocalScorePass{}, fmt.Errorf("received %d score contexts for %d candidates", len(contexts), len(candidates))
	}
	if len(candidates) <= 30 || len(candidateKeys) != len(candidates) {
		return Paul2013ExpandedCandidateLocalScorePass{}, fmt.Errorf("expanded candidate scoring requires aligned keys for an oversized candidate list, got %d candidates and %d keys", len(candidates), len(candidateKeys))
	}
	initialPrefix, err := Paul2013CandidateLocalScorePrefix(candidates, positionCount)
	if err != nil {
		return Paul2013ExpandedCandidateLocalScorePass{}, fmt.Errorf("determine native local-score prefix: %w", err)
	}
	protection, err := ExpandPaul2013CandidateProtectionPrefix(
		candidates, candidateKeys, wholePositionKeys, initialPrefix, wholePositionCandidateCount,
	)
	if err != nil {
		return Paul2013ExpandedCandidateLocalScorePass{}, fmt.Errorf("expand protected candidate prefix: %w", err)
	}
	alignedContexts := append([]Paul2013UnitScoreContext(nil), contexts...)
	for outputIndex, sourceIndex := range protection.ScoredPrefixSources {
		alignedContexts[outputIndex] = contexts[sourceIndex]
	}
	scores, err := scorePaul2013CandidatePrefix(
		model, table, protection.Candidates, alignedContexts, protection.NodeFlags, protection.ScoredPrefix,
	)
	if err != nil {
		return Paul2013ExpandedCandidateLocalScorePass{}, err
	}
	return Paul2013ExpandedCandidateLocalScorePass{Protection: protection, Scores: scores}, nil
}

func scorePaul2013CandidatePrefix(
	model Paul2013ContinuityModel,
	table *distance.FeatureTable,
	candidates []Paul2013CandidateContinuity,
	contexts []Paul2013UnitScoreContext,
	flags []uint16,
	prefix int,
) (Paul2013CandidateLocalScorePass, error) {
	result := Paul2013CandidateLocalScorePass{
		Candidates: make([]Paul2013CandidateLocalScore, len(candidates)), ScoredPrefix: prefix,
	}
	for index, candidate := range candidates {
		result.Candidates[index] = Paul2013CandidateLocalScore{
			Continuity: candidate,
		}
		if index < len(flags) {
			result.Candidates[index].NodeFlag = flags[index]
		}
		if index >= prefix {
			continue
		}
		record, err := model.ReadRecord(candidate.Unit.Bank, candidate.Unit.Index)
		if err != nil {
			return Paul2013CandidateLocalScorePass{}, fmt.Errorf("read candidate %d (%s:%d): %w", index, candidate.Unit.Bank, candidate.Unit.Index, err)
		}
		input, err := BuildPaul2013UnitScoreInput(record, contexts[index])
		if err != nil {
			return Paul2013CandidateLocalScorePass{}, fmt.Errorf("build candidate %d unit-score input: %w", index, err)
		}
		cost, err := ScorePaul2013CandidateLocalCost(table, input, candidate)
		if err != nil {
			return Paul2013CandidateLocalScorePass{}, fmt.Errorf("score candidate %d: %w", index, err)
		}
		result.Candidates[index].LocalCost = cost
		result.Candidates[index].HasLocalCost = true
	}
	return result, nil
}

// Paul2013TailOrderingScore applies the per-candidate multiplier in
// FUN_10023350. The value at 0x1006d1a4 is a float32 loaded from raw bits
// 0x00002000 in the local DLL; it is intentionally represented as such here.
func Paul2013TailOrderingScore(localCost float32, candidateClass, contextClass byte) (float32, error) {
	if math.IsNaN(float64(localCost)) || math.IsInf(float64(localCost), 0) {
		return 0, errors.New("candidate local cost must be finite")
	}
	switch {
	case candidateClass == 'd':
		return localCost * 5, nil
	case candidateClass == 0x0c && contextClass != 0x0c:
		return localCost * math.Float32frombits(0x00002000), nil
	default:
		return localCost, nil
	}
}

// SortPaul2013CandidateTailAndCap reproduces the final oversized-list step in
// FUN_10023350. It preserves the prefix before tailStart, sorts only
// the remaining candidates by their precomputed secondary score, then retains
// 30 entries. Lists of at most 30 are returned unchanged. Larger lists use
// the recovered native partition and heap-fallback path.
func SortPaul2013CandidateTailAndCap(candidates []Paul2013ScoredCandidate, tailStart int) ([]Paul2013ScoredCandidate, error) {
	if tailStart < 0 || tailStart > len(candidates) {
		return nil, fmt.Errorf("candidate tail start %d is outside 0..%d", tailStart, len(candidates))
	}
	result := append([]Paul2013ScoredCandidate(nil), candidates...)
	if len(result) <= 30 {
		return result, nil
	}
	for index, candidate := range result {
		if math.IsNaN(float64(candidate.LocalCost)) || math.IsInf(float64(candidate.LocalCost), 0) ||
			math.IsNaN(float64(candidate.SecondaryOrderScore)) || math.IsInf(float64(candidate.SecondaryOrderScore), 0) {
			return nil, fmt.Errorf("candidate %d has a non-finite score", index)
		}
	}
	tail := result[tailStart:]
	less := func(left, right Paul2013ScoredCandidate) bool {
		return left.SecondaryOrderScore < right.SecondaryOrderScore
	}
	sortPaul2013Native(tail, less)
	return result[:30], nil
}

// FinalizePaul2013CandidateShortlist ports FUN_10023350's complete oversized
// tail step when the caller supplies its secondary scores and final protected
// flag count. It sorts only the native tail, caps the shortlist at 30, then
// replaces trailing unprotected slots with flagged candidates from the
// remainder. The native pass permits at most five protected candidates in the
// output and does not restore tail candidates when the prefix already has
// more than four flags.
func FinalizePaul2013CandidateShortlist(
	candidates []Paul2013ScoredCandidate,
	tailStart int,
	totalProtected int,
) (Paul2013CandidateShortlist, error) {
	if tailStart < 0 || tailStart > len(candidates) {
		return Paul2013CandidateShortlist{}, fmt.Errorf("candidate tail start %d is outside 0..%d", tailStart, len(candidates))
	}
	if totalProtected < 0 || totalProtected > len(candidates) {
		return Paul2013CandidateShortlist{}, fmt.Errorf("protected candidate count %d is outside 0..%d", totalProtected, len(candidates))
	}
	result := append([]Paul2013ScoredCandidate(nil), candidates...)
	for index, candidate := range result {
		if candidate.NodeFlag > 1 {
			return Paul2013CandidateShortlist{}, fmt.Errorf("candidate %d has invalid native flag %d", index, candidate.NodeFlag)
		}
		if math.IsNaN(float64(candidate.LocalCost)) || math.IsInf(float64(candidate.LocalCost), 0) ||
			math.IsNaN(float64(candidate.SecondaryOrderScore)) || math.IsInf(float64(candidate.SecondaryOrderScore), 0) {
			return Paul2013CandidateShortlist{}, fmt.Errorf("candidate %d has a non-finite score", index)
		}
	}
	if len(result) <= 30 {
		return Paul2013CandidateShortlist{Candidates: result}, nil
	}

	sortPaul2013Native(result[tailStart:], func(left, right Paul2013ScoredCandidate) bool {
		return left.SecondaryOrderScore < right.SecondaryOrderScore
	})
	shortlist := result[:30]
	protectedInPrefix := 0
	for _, candidate := range shortlist {
		if candidate.NodeFlag != 0 {
			protectedInPrefix++
		}
	}
	final := Paul2013CandidateShortlist{Candidates: append([]Paul2013ScoredCandidate(nil), shortlist...)}
	if protectedInPrefix > 4 || protectedInPrefix >= totalProtected {
		return final, nil
	}

	// The native code walks backward through the 30 output slots, skipping
	// already-protected nodes, and writes each selected tail identity there.
	slot := len(final.Candidates) - 1
	for _, candidate := range result[30:] {
		if candidate.NodeFlag == 0 {
			continue
		}
		if slot > 0 {
			for slot > 0 && final.Candidates[slot].NodeFlag == 1 {
				slot--
			}
		}
		final.Candidates[slot] = candidate
		final.ProtectedReplacements = append(final.ProtectedReplacements, Paul2013ProtectedCandidateReplacement{
			Slot: slot, Candidate: candidate,
		})
		protectedInPrefix++
		slot--
		if slot < 1 || protectedInPrefix >= totalProtected || protectedInPrefix > 4 {
			break
		}
	}
	return final, nil
}

// Paul2013ContinuityFeatureScale ports the span normalization in
// FUN_10023350 before it calls FUN_100182e0. The DLL constants are 1.0 at
// 0x1006d170 and 2.0 at 0x1006d174; the resulting multiplier applies to the
// feature-distance term, while categorical penalties are left unchanged.
func Paul2013ContinuityFeatureScale(totalSpan, weightedSide uint16) (float32, error) {
	if totalSpan == 0 {
		return 0, errors.New("candidate continuity span must be positive")
	}
	normalizedSpan := float32(totalSpan) / 2
	if normalizedSpan <= 1 {
		normalizedSpan = 1
	}
	denominator := float32(weightedSide) + normalizedSpan
	scale := float32(1) / denominator
	return scale, nil
}

// ScorePaul2013CandidateLocalCost derives the continuity multiplier and feeds
// it into the recovered per-unit scorer. Context feature selection and
// categorical state remain in input, as their producers are separate stages.
func ScorePaul2013CandidateLocalCost(
	table *distance.FeatureTable,
	input UnitScoreInput,
	continuity Paul2013CandidateContinuity,
) (float32, error) {
	scale, err := Paul2013ContinuityFeatureScale(continuity.TotalSpan, continuity.WeightedSide)
	if err != nil {
		return 0, err
	}
	input.Scale = scale
	return ScoreUnitCost(table, input)
}
