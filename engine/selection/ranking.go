package selection

import (
	"errors"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

// ClassCandidate is one already-built class in the legacy candidate index.
// Population is the number of unit candidates in that class.
type ClassCandidate struct {
	ID         uint32
	Key        [5]byte
	Population uint32
}

// RankOptions mirrors the observed class-count and cumulative-population
// limits. ViewMode 0 compares the five-byte key; 1 and 2 also compare the
// corresponding ten-byte feature projection.
type RankOptions struct {
	MaxClasses int
	MaxUnits   uint64
	ViewMode   byte
}

const paul2013MaxScoredClassCandidates = 10_000

// ClassMismatchDistance computes the weighted key mismatch used by the
// recovered class-ranking routine. A nonzero view mode adds the matching
// ten-byte projected mismatch.
func ClassMismatchDistance(input, candidate [5]byte, viewMode byte) (int64, error) {
	var score int64
	for position := range input {
		if input[position] == candidate[position] {
			continue
		}
		weight, ok := text.Paul2013KeyMismatchWeight(position)
		if !ok {
			return 0, fmt.Errorf("no Paul 2013 key weight for position %d", position)
		}
		score += int64(weight)
	}
	if viewMode == 0 {
		return score, nil
	}
	inputView, err := text.Paul2013FeatureViewForKey(input, viewMode)
	if err != nil {
		return 0, err
	}
	candidateView, err := text.Paul2013FeatureViewForKey(candidate, viewMode)
	if err != nil {
		return 0, err
	}
	for position := range inputView {
		if inputView[position] == candidateView[position] {
			continue
		}
		weight, ok := text.Paul2013ViewMismatchWeight(position)
		if !ok {
			return 0, fmt.Errorf("no Paul 2013 view weight for position %d", position)
		}
		score += int64(weight)
	}
	return score, nil
}

// RankClasses preserves source order when both observed candidate limits fit.
// Otherwise it scores and sorts candidates by weighted mismatch, then
// emits complete classes until either limit is reached; the final class may
// cross the unit budget. Feature-view pools above 10,000 rank only their first
// 10,000 entries. Larger key-only pools are rejected because the native path
// sorts an uncharacterized scratch tail. Scored entries use the recovered
// FUN_1001b5f0 ordering algorithm.
func RankClasses(contextKey [5]byte, candidates []ClassCandidate, options RankOptions) ([]ClassCandidate, error) {
	if options.MaxClasses <= 0 || options.MaxUnits == 0 {
		return nil, errors.New("class ranking limits must be positive")
	}
	if options.ViewMode > 2 {
		return nil, fmt.Errorf("unsupported Paul 2013 feature view mode %d", options.ViewMode)
	}
	if len(candidates) == 0 {
		return []ClassCandidate{}, nil
	}
	totalPopulation := uint64(0)
	for _, candidate := range candidates {
		totalPopulation += uint64(candidate.Population)
	}
	if len(candidates) <= options.MaxClasses && totalPopulation <= options.MaxUnits {
		return append([]ClassCandidate(nil), candidates...), nil
	}

	// Both native scoring loops fill at most 10,000 scratch entries. The
	// feature-view path sorts only that copied prefix. The key-only path passes
	// the original count to its sort helper, so its unscored scratch tail is not
	// characterized well enough to reproduce for larger pools.
	scoreCount := len(candidates)
	if scoreCount > paul2013MaxScoredClassCandidates {
		if options.ViewMode == 0 {
			return nil, fmt.Errorf("class candidate pool has %d entries; native score buffer supports at most %d", len(candidates), paul2013MaxScoredClassCandidates)
		}
		scoreCount = paul2013MaxScoredClassCandidates
	}
	type scoredClass struct {
		candidate ClassCandidate
		score     int64
	}
	scored := make([]scoredClass, scoreCount)
	for position, candidate := range candidates[:scoreCount] {
		score, err := ClassMismatchDistance(contextKey, candidate.Key, options.ViewMode)
		if err != nil {
			return nil, fmt.Errorf("score class %d: %w", candidate.ID, err)
		}
		scored[position] = scoredClass{candidate: candidate, score: score}
	}
	less := func(left, right scoredClass) bool { return left.score < right.score }
	sortPaul2013Native(scored, less)

	limit := options.MaxClasses
	if limit > len(scored) {
		limit = len(scored)
	}
	ranked := make([]ClassCandidate, 0, limit)
	population := uint64(0)
	for _, item := range scored {
		if len(ranked) >= options.MaxClasses {
			break
		}
		ranked = append(ranked, item.candidate)
		population += uint64(item.candidate.Population)
		if population >= options.MaxUnits {
			break
		}
	}
	return ranked, nil
}

// RankClassRecords applies the observed class limits to model-backed class
// records while retaining each record's unit membership.
func RankClassRecords(contextKey [5]byte, classes []voice.ClassRecord, options RankOptions) ([]voice.ClassRecord, error) {
	candidates := make([]ClassCandidate, len(classes))
	byID := make(map[uint32]voice.ClassRecord, len(classes))
	for index, class := range classes {
		if _, exists := byID[class.ID]; exists {
			return nil, fmt.Errorf("duplicate class ID %d", class.ID)
		}
		population := uint64(len(class.Members))
		if population > uint64(^uint32(0)) {
			return nil, fmt.Errorf("class %d population exceeds supported width", class.ID)
		}
		candidates[index] = ClassCandidate{ID: class.ID, Key: class.Key, Population: uint32(population)}
		byID[class.ID] = class
	}
	ranked, err := RankClasses(contextKey, candidates, options)
	if err != nil {
		return nil, err
	}
	result := make([]voice.ClassRecord, len(ranked))
	for index, candidate := range ranked {
		class, ok := byID[candidate.ID]
		if !ok {
			return nil, fmt.Errorf("ranked class ID %d has no source record", candidate.ID)
		}
		result[index] = class
	}
	return result, nil
}

// ExpandClassRecords flattens ranked classes into their unit locations in
// class order and member order. maxUnits bounds the expanded list; the final
// class may be truncated at that limit, matching the observed unit-candidate
// capacity. This does not generate neighboring-context, duration, or scoring
// metadata for the candidates.
func ExpandClassRecords(classes []voice.ClassRecord, maxUnits uint64) ([]voice.UnitLocation, error) {
	if maxUnits == 0 {
		return nil, errors.New("unit expansion limit must be positive")
	}
	units := make([]voice.UnitLocation, 0)
	for _, class := range classes {
		for _, member := range class.Members {
			if uint64(len(units)) == maxUnits {
				return units, nil
			}
			units = append(units, member)
		}
	}
	return units, nil
}

// ExactContextUnitCandidates resolves only the exact class key produced by a
// supplied seven-byte context and returns its members under maxUnits. found
// is false when the catalog has no exact key; this function does not broaden
// the lookup or construct candidate scoring metadata.
func ExactContextUnitCandidates(catalog *voice.ClassCatalog, context text.Context, maxUnits uint64) ([]UnitRef, bool, error) {
	if catalog == nil {
		return nil, false, errors.New("context candidate lookup has no class catalog")
	}
	if maxUnits == 0 {
		return nil, false, errors.New("unit expansion limit must be positive")
	}
	class, found := catalog.LookupContext(context)
	if !found {
		return []UnitRef{}, false, nil
	}
	locations, err := ExpandClassRecords([]voice.ClassRecord{class}, maxUnits)
	if err != nil {
		return nil, false, err
	}
	candidates := make([]UnitRef, len(locations))
	for index, location := range locations {
		candidates[index] = UnitRef{Bank: location.Bank, Index: location.Index}
	}
	return candidates, true, nil
}
