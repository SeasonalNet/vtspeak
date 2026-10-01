package selection

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/distance"
	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

const (
	paul2013ModelStateRecordBaseOffset = 0x64c
	paul2013ModelStateRecordStride     = 0x3c0
	paul2013ModelStateRecordCountLimit = 100
)

// Paul2013PackedContextCandidates contains exact-key candidates aligned to
// the flattened phone contexts from one already-selected native record group.
// A false ExactClassFound entry is an ordinary catalog miss, not a fallback.
type Paul2013PackedContextCandidates struct {
	Prepared        text.Paul2013PackedContextPhoneGroupResult
	Contexts        []text.Context
	CandidatePools  [][]UnitRef
	ExactClassFound []bool
}

// Paul2013PhoneMarkerGroupCandidates retains exact-key lookup results aligned
// to one FUN_10012df0 group. Group descriptors and prepared records are kept
// so callers can preserve the native boundary and output-short metadata.
type Paul2013PhoneMarkerGroupCandidates struct {
	Descriptor      text.Paul2013RecordGroupDescriptor
	Contexts        []text.Context
	CandidatePools  [][]UnitRef
	ExactClassFound []bool
}

// Paul2013PhoneMarkerRecordCandidates carries the composed marker-prepass,
// group projection, packed-context preparation, and exact catalog lookup for
// a supplied native record stream.
type Paul2013PhoneMarkerRecordCandidates struct {
	Prepared text.Paul2013PreparedRecordGroups
	Groups   []Paul2013PhoneMarkerGroupCandidates
}

// Paul2013PhoneMarkerContinuityGroup retains continuity-ranked candidates
// aligned to one native record group.
type Paul2013PhoneMarkerContinuityGroup struct {
	Descriptor text.Paul2013RecordGroupDescriptor
	Candidates [][]Paul2013CandidateContinuity
}

// Paul2013PhoneMarkerContinuityResult preserves prepared source records and
// the group alignment after continuity ranking.
type Paul2013PhoneMarkerContinuityResult struct {
	Prepared text.Paul2013PreparedRecordGroups
	Groups   []Paul2013PhoneMarkerContinuityGroup
}

// Paul2013PhoneMarkerScoreGroup retains local-score passes aligned to one
// native record group.
type Paul2013PhoneMarkerScoreGroup struct {
	Descriptor text.Paul2013RecordGroupDescriptor
	Passes     []Paul2013QueriedCandidateScorePass
}

// Paul2013PhoneMarkerScoreResult carries candidates through exact lookup,
// continuity ranking, and local scoring while preserving group boundaries.
type Paul2013PhoneMarkerScoreResult struct {
	Prepared text.Paul2013PreparedRecordGroups
	Groups   []Paul2013PhoneMarkerScoreGroup
}

// Paul2013PhoneMarkerShortlistGroup retains final shortlist rows aligned to
// one native record group.
type Paul2013PhoneMarkerShortlistGroup struct {
	Descriptor text.Paul2013RecordGroupDescriptor
	Shortlists []Paul2013CandidateShortlist
}

// Paul2013PhoneMarkerShortlistResult carries finalized rows and their native
// group boundaries into path scoring.
type Paul2013PhoneMarkerShortlistResult struct {
	Prepared text.Paul2013PreparedRecordGroups
	Groups   []Paul2013PhoneMarkerShortlistGroup
}

// Paul2013PhoneMarkerPathResult retains each group's shortlist rows and the
// minimum-cost path selected across their flattened phone order.
type Paul2013PhoneMarkerPathResult struct {
	Prepared text.Paul2013PreparedRecordGroups
	Groups   []Paul2013PhoneMarkerShortlistGroup
	Path     Paul2013PathResult
}

// LookupPaul2013PackedContextCandidates composes FUN_10017510's packed phone
// context producer with the catalog's exact five-byte key lookup. It returns
// one bounded unit pool per phone context in record/phone order. It does not
// run whole-position query broadening, fallback-row generation, continuity
// ranking, or path scoring.
func LookupPaul2013PackedContextCandidates(
	catalog *voice.ClassCatalog,
	records [][]byte,
	maxUnits uint64,
) (Paul2013PackedContextCandidates, error) {
	if catalog == nil {
		return Paul2013PackedContextCandidates{}, fmt.Errorf("packed-context candidate lookup has no class catalog")
	}
	if maxUnits == 0 {
		return Paul2013PackedContextCandidates{}, fmt.Errorf("packed-context candidate lookup requires a positive unit limit")
	}
	prepared, contexts, err := text.PreparePaul2013SelectionContextsFromPackedContextPhoneGroup(records)
	if err != nil {
		return Paul2013PackedContextCandidates{}, fmt.Errorf("prepare Paul 2013 packed selection contexts: %w", err)
	}
	result := Paul2013PackedContextCandidates{
		Prepared:        prepared,
		Contexts:        contexts,
		CandidatePools:  make([][]UnitRef, len(contexts)),
		ExactClassFound: make([]bool, len(contexts)),
	}
	for contextIndex, context := range contexts {
		units, found, err := ExactContextUnitCandidates(catalog, context, maxUnits)
		if err != nil {
			return Paul2013PackedContextCandidates{}, fmt.Errorf("lookup packed context %d: %w", contextIndex, err)
		}
		result.CandidatePools[contextIndex] = units
		result.ExactClassFound[contextIndex] = found
	}
	return result, nil
}

// LookupPaul2013PhoneMarkerRecordCandidates composes record-group formation,
// the FUN_10017100 marker prepass, each group's FUN_10017510 context pass, and
// exact five-byte catalog lookup. The source records and special-key table
// remain explicit inputs. It preserves per-group context ordering and treats
// a missing exact class as an empty pool; it does not broaden queries, build
// fallback rows, rank continuity, or score paths. This slice-based variant
// does not perform the arena-wide phone-state reset; use
// LookupPaul2013PhoneMarkerArenaCandidates when the complete contiguous arena
// is available.
func LookupPaul2013PhoneMarkerRecordCandidates(
	catalog *voice.ClassCatalog,
	records [][]byte,
	specialKeys [][]byte,
	characterMap [256]byte,
	maxUnits uint64,
) (Paul2013PhoneMarkerRecordCandidates, error) {
	if catalog == nil {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("phone-marker candidate lookup has no class catalog")
	}
	if maxUnits == 0 {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("phone-marker candidate lookup requires a positive unit limit")
	}
	prepared, err := text.PreparePaul2013PhoneMarkerRecordGroups(records, specialKeys, characterMap)
	if err != nil {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("prepare Paul 2013 phone-marker record groups: %w", err)
	}
	return lookupPaul2013PhoneMarkerRecordCandidates(catalog, prepared, maxUnits)
}

// LookupPaul2013PhoneMarkerArenaCandidates adds FUN_10017510's shared-state
// reset to the full-record composition. arena must contain recordCount
// contiguous 0x3c0-byte records; the returned Prepared value retains the
// reset copy for inspection.
func LookupPaul2013PhoneMarkerArenaCandidates(
	catalog *voice.ClassCatalog,
	arena []byte,
	recordCount int,
	specialKeys [][]byte,
	characterMap [256]byte,
	maxUnits uint64,
) (Paul2013PhoneMarkerRecordCandidates, error) {
	return lookupPaul2013PhoneMarkerArenaCandidates(catalog, recordCount, maxUnits, func() (text.Paul2013PreparedRecordGroups, error) {
		return text.PreparePaul2013PhoneMarkerRecordGroupsFromArena(
			arena, recordCount, specialKeys, characterMap,
		)
	})
}

// LookupPaul2013PhoneMarkerArenaCandidatesWithPointerResolver composes the
// complete arena reset, marker/group preparation, and exact class lookup,
// resolving native +0x3b8 string pointers through the supplied address space.
func LookupPaul2013PhoneMarkerArenaCandidatesWithPointerResolver(
	catalog *voice.ClassCatalog,
	arena []byte,
	recordCount int,
	resolvePointer text.Paul2013CStringPointerResolver,
	characterMap [256]byte,
	maxUnits uint64,
) (Paul2013PhoneMarkerRecordCandidates, error) {
	return lookupPaul2013PhoneMarkerArenaCandidates(catalog, recordCount, maxUnits, func() (text.Paul2013PreparedRecordGroups, error) {
		return text.PreparePaul2013PhoneMarkerRecordGroupsFromArenaWithPointerResolver(
			arena, recordCount, resolvePointer, characterMap,
		)
	})
}

// LookupPaul2013PhoneMarkerModelStateArenaCandidatesWithPointerResolver
// consumes the full counted model-state arena produced by the text pipeline.
// It reads the signed record count at +2, slices the 0x3c0-byte records from
// +0x64c, then applies the existing marker/group/context/catalog stages. The
// parser-row writer, +0x3b8 address mapping, and vendor pointer resolution
// remain explicit inputs.
func LookupPaul2013PhoneMarkerModelStateArenaCandidatesWithPointerResolver(
	catalog *voice.ClassCatalog,
	modelStateArena []byte,
	resolvePointer text.Paul2013CStringPointerResolver,
	characterMap [256]byte,
	maxUnits uint64,
) (Paul2013PhoneMarkerRecordCandidates, error) {
	if catalog == nil {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("model-state candidate lookup has no class catalog")
	}
	if maxUnits == 0 {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("model-state candidate lookup requires a positive unit limit")
	}
	if len(modelStateArena) < 4 {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("model-state arena has %d bytes, need its signed record count at +2", len(modelStateArena))
	}
	recordCount := int(int16(binary.LittleEndian.Uint16(modelStateArena[2:4])))
	if recordCount < 0 || recordCount > paul2013ModelStateRecordCountLimit {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("model-state arena record count %d outside [0, %d]", recordCount, paul2013ModelStateRecordCountLimit)
	}
	if paul2013ModelStateRecordBaseOffset > len(modelStateArena) || recordCount > (len(modelStateArena)-paul2013ModelStateRecordBaseOffset)/paul2013ModelStateRecordStride {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(modelStateArena), recordCount, paul2013ModelStateRecordBaseOffset)
	}
	start := paul2013ModelStateRecordBaseOffset
	end := start + recordCount*paul2013ModelStateRecordStride
	return LookupPaul2013PhoneMarkerArenaCandidatesWithPointerResolver(
		catalog, modelStateArena[start:end], recordCount, resolvePointer, characterMap, maxUnits,
	)
}

func lookupPaul2013PhoneMarkerArenaCandidates(
	catalog *voice.ClassCatalog,
	recordCount int,
	maxUnits uint64,
	prepare func() (text.Paul2013PreparedRecordGroups, error),
) (Paul2013PhoneMarkerRecordCandidates, error) {
	if catalog == nil {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("phone-marker candidate lookup has no class catalog")
	}
	if maxUnits == 0 {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("phone-marker candidate lookup requires a positive unit limit")
	}
	prepared, err := prepare()
	if err != nil {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("prepare Paul 2013 phone-marker arena groups: %w", err)
	}
	return lookupPaul2013PhoneMarkerRecordCandidates(catalog, prepared, maxUnits)
}

// RankPaul2013PhoneMarkerRecordContinuity carries exact-key candidate pools
// through FUN_100230a0 in flattened record/phone order, then restores each
// group's boundary in the result. modes must have one native context mode per
// phone. Exact misses fail closed because fallback rows have not been applied.
func RankPaul2013PhoneMarkerRecordContinuity(
	model Paul2013ContinuityModel,
	input Paul2013PhoneMarkerRecordCandidates,
	modes []byte,
) (Paul2013PhoneMarkerContinuityResult, error) {
	if model == nil {
		return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("phone-marker continuity ranking has no model")
	}
	if len(input.Groups) != len(input.Prepared.Descriptors) {
		return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("candidate result has %d groups for %d prepared descriptors", len(input.Groups), len(input.Prepared.Descriptors))
	}
	var positions [][]UnitRef
	groupStarts := make([]int, len(input.Groups)+1)
	for groupIndex, group := range input.Groups {
		if group.Descriptor != input.Prepared.Descriptors[groupIndex] {
			return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("candidate group %d descriptor differs from prepared group", groupIndex)
		}
		if len(group.CandidatePools) != len(group.Contexts) || len(group.ExactClassFound) != len(group.Contexts) {
			return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("candidate group %d has %d contexts, %d pools, and %d exact-match flags", groupIndex, len(group.Contexts), len(group.CandidatePools), len(group.ExactClassFound))
		}
		groupStarts[groupIndex] = len(positions)
		for position, pool := range group.CandidatePools {
			if !group.ExactClassFound[position] || len(pool) == 0 {
				return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("record group %d phone %d has no exact candidates; native fallback must run before continuity ranking", groupIndex, position)
			}
			positions = append(positions, pool)
		}
	}
	groupStarts[len(input.Groups)] = len(positions)
	if len(positions) == 0 {
		return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("phone-marker continuity ranking has no candidate positions")
	}
	if len(modes) != len(positions) {
		return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("received %d context modes for %d phone positions", len(modes), len(positions))
	}
	ranked := make([][]Paul2013CandidateContinuity, len(positions))
	for position := range positions {
		candidates, err := RankPaul2013CandidatesByContinuity(model, positions, modes, position)
		if err != nil {
			return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("rank phone position %d: %w", position, err)
		}
		ranked[position] = candidates
	}
	result := Paul2013PhoneMarkerContinuityResult{
		Prepared: input.Prepared,
		Groups:   make([]Paul2013PhoneMarkerContinuityGroup, len(input.Groups)),
	}
	for groupIndex, group := range input.Groups {
		start, end := groupStarts[groupIndex], groupStarts[groupIndex+1]
		result.Groups[groupIndex] = Paul2013PhoneMarkerContinuityGroup{
			Descriptor: group.Descriptor,
			Candidates: ranked[start:end],
		}
	}
	return result, nil
}

// RankPaul2013PhoneMarkerRecordContinuityFromTransitionStates derives the
// native per-position modes from FUN_10024680's compact six-byte states before
// running the continuity pass. State rows must already align with the
// flattened candidate positions, including both rows emitted by fallback.
func RankPaul2013PhoneMarkerRecordContinuityFromTransitionStates(
	model Paul2013ContinuityModel,
	input Paul2013PhoneMarkerRecordCandidates,
	states []Paul2013TransitionContextState,
) (Paul2013PhoneMarkerContinuityResult, error) {
	positionCount := 0
	for _, group := range input.Groups {
		positionCount += len(group.CandidatePools)
	}
	if len(states) != positionCount {
		return Paul2013PhoneMarkerContinuityResult{}, fmt.Errorf("received %d transition states for %d candidate positions", len(states), positionCount)
	}
	modes, err := Paul2013ContextModesFromTransitionStates(states)
	if err != nil {
		return Paul2013PhoneMarkerContinuityResult{}, err
	}
	return RankPaul2013PhoneMarkerRecordContinuity(model, input, modes)
}

// ScorePaul2013PhoneMarkerRecordCandidates composes the continuity ranker and
// bounded local-score pass for exact pools from one prepared native record
// stream. modes and scoringInputs must align with flattened record/phone
// order; scoring key maps and positionCount remain explicit native state.
// Exact misses still require the native fallback stage first.
func ScorePaul2013PhoneMarkerRecordCandidates(
	model Paul2013ContinuityModel,
	table *distance.FeatureTable,
	input Paul2013PhoneMarkerRecordCandidates,
	modes []byte,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
) (Paul2013PhoneMarkerScoreResult, error) {
	continuity, err := RankPaul2013PhoneMarkerRecordContinuity(model, input, modes)
	if err != nil {
		return Paul2013PhoneMarkerScoreResult{}, fmt.Errorf("rank phone-marker record continuity: %w", err)
	}
	positionCountTotal := 0
	for _, group := range continuity.Groups {
		positionCountTotal += len(group.Candidates)
	}
	if len(scoringInputs) != positionCountTotal {
		return Paul2013PhoneMarkerScoreResult{}, fmt.Errorf("received %d scoring input rows for %d phone positions", len(scoringInputs), positionCountTotal)
	}
	result := Paul2013PhoneMarkerScoreResult{
		Prepared: continuity.Prepared,
		Groups:   make([]Paul2013PhoneMarkerScoreGroup, len(continuity.Groups)),
	}
	position := 0
	for groupIndex, group := range continuity.Groups {
		scoredGroup := Paul2013PhoneMarkerScoreGroup{
			Descriptor: group.Descriptor,
			Passes:     make([]Paul2013QueriedCandidateScorePass, len(group.Candidates)),
		}
		for groupPosition, candidates := range group.Candidates {
			input := scoringInputs[position]
			pass := Paul2013QueriedCandidateScorePass{Ranked: candidates}
			if len(candidates) > 30 {
				expanded, err := ScorePaul2013ExpandedCandidateLocalPrefix(
					model, table, candidates, positionCount, input.Contexts,
					input.CandidateKeys, input.WholePositionKeys, input.WholePositionCandidateCount,
				)
				if err != nil {
					return Paul2013PhoneMarkerScoreResult{}, fmt.Errorf("score record group %d position %d oversized candidates: %w", groupIndex, groupPosition, err)
				}
				pass.Scores = expanded.Scores
				pass.Protection = expanded.Protection
				pass.HasProtection = true
			} else {
				scores, err := ScorePaul2013CandidateLocalPrefix(
					model, table, candidates, positionCount, input.Contexts,
					input.CandidateKeys, input.WholePositionKeys,
				)
				if err != nil {
					return Paul2013PhoneMarkerScoreResult{}, fmt.Errorf("score record group %d position %d candidates: %w", groupIndex, groupPosition, err)
				}
				pass.Scores = scores
			}
			scoredGroup.Passes[groupPosition] = pass
			position++
		}
		result.Groups[groupIndex] = scoredGroup
	}
	return result, nil
}

// ScorePaul2013PhoneMarkerRecordCandidatesFromTransitionStates composes the
// recovered state-byte mode extraction with continuity ranking and bounded
// local scoring. Score inputs, key maps, and positionCount remain explicit
// caller-owned state.
func ScorePaul2013PhoneMarkerRecordCandidatesFromTransitionStates(
	model Paul2013ContinuityModel,
	table *distance.FeatureTable,
	input Paul2013PhoneMarkerRecordCandidates,
	positionCount int,
	states []Paul2013TransitionContextState,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
) (Paul2013PhoneMarkerScoreResult, error) {
	positionCountFromStates := 0
	for _, group := range input.Groups {
		positionCountFromStates += len(group.CandidatePools)
	}
	if len(states) != positionCountFromStates {
		return Paul2013PhoneMarkerScoreResult{}, fmt.Errorf("received %d transition states for %d candidate positions", len(states), positionCountFromStates)
	}
	modes, err := Paul2013ContextModesFromTransitionStates(states)
	if err != nil {
		return Paul2013PhoneMarkerScoreResult{}, err
	}
	return ScorePaul2013PhoneMarkerRecordCandidates(
		model, table, input, modes, positionCount, scoringInputs,
	)
}

// FinalizePaul2013PhoneMarkerRecordShortlists applies native tail ordering,
// truncation, and protected-candidate restoration to every scored phone. The
// three state slices align to flattened record/phone order and remain caller
// inputs because their text-state producers are not yet recovered.
func FinalizePaul2013PhoneMarkerRecordShortlists(
	input Paul2013PhoneMarkerScoreResult,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
) (Paul2013PhoneMarkerShortlistResult, error) {
	positionCount := 0
	for _, group := range input.Groups {
		positionCount += len(group.Passes)
	}
	if len(input.Groups) != len(input.Prepared.Descriptors) {
		return Paul2013PhoneMarkerShortlistResult{}, fmt.Errorf("score result has %d groups for %d prepared descriptors", len(input.Groups), len(input.Prepared.Descriptors))
	}
	if len(tailStarts) != positionCount || len(candidateClasses) != positionCount || len(contextClasses) != positionCount {
		return Paul2013PhoneMarkerShortlistResult{}, fmt.Errorf("received %d tail starts, %d candidate-class rows, and %d context-class rows for %d phone positions", len(tailStarts), len(candidateClasses), len(contextClasses), positionCount)
	}
	result := Paul2013PhoneMarkerShortlistResult{
		Prepared: input.Prepared,
		Groups:   make([]Paul2013PhoneMarkerShortlistGroup, len(input.Groups)),
	}
	position := 0
	for groupIndex, group := range input.Groups {
		if group.Descriptor != input.Prepared.Descriptors[groupIndex] {
			return Paul2013PhoneMarkerShortlistResult{}, fmt.Errorf("score group %d descriptor differs from prepared group", groupIndex)
		}
		shortlistGroup := Paul2013PhoneMarkerShortlistGroup{
			Descriptor: group.Descriptor,
			Shortlists: make([]Paul2013CandidateShortlist, len(group.Passes)),
		}
		for groupPosition, pass := range group.Passes {
			shortlist, err := FinalizePaul2013QueriedCandidateShortlist(
				pass, tailStarts[position], candidateClasses[position], contextClasses[position],
			)
			if err != nil {
				return Paul2013PhoneMarkerShortlistResult{}, fmt.Errorf("finalize record group %d phone %d: %w", groupIndex, groupPosition, err)
			}
			shortlistGroup.Shortlists[groupPosition] = shortlist
			position++
		}
		result.Groups[groupIndex] = shortlistGroup
	}
	return result, nil
}

// SelectPaul2013PhoneMarkerRecordPath composes continuity ranking, bounded
// local scoring, shortlist finalization, and model-backed path selection for
// exact pools from a complete supplied record stream. All
// text-derived score/tail state and per-edge transition state remain explicit
// caller inputs; exact lookup misses still require fallback before this path.
func SelectPaul2013PhoneMarkerRecordPath(
	model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	},
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	input Paul2013PhoneMarkerRecordCandidates,
	modes []byte,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PhoneMarkerPathResult, error) {
	scores, err := ScorePaul2013PhoneMarkerRecordCandidates(
		model, featureTable, input, modes, positionCount, scoringInputs,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	shortlists, err := FinalizePaul2013PhoneMarkerRecordShortlists(
		scores, tailStarts, candidateClasses, contextClasses,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	passes := make([]Paul2013CandidateLocalScorePass, 0)
	for _, group := range shortlists.Groups {
		for _, shortlist := range group.Shortlists {
			pass := Paul2013CandidateLocalScorePass{
				ScoredPrefix: len(shortlist.Candidates),
				Candidates:   make([]Paul2013CandidateLocalScore, len(shortlist.Candidates)),
			}
			for candidateIndex, candidate := range shortlist.Candidates {
				pass.Candidates[candidateIndex] = Paul2013CandidateLocalScore{
					Continuity: candidate.Continuity, NodeFlag: candidate.NodeFlag,
					LocalCost: candidate.LocalCost, HasLocalCost: true,
				}
			}
			passes = append(passes, pass)
		}
	}
	path, err := SelectPaul2013PathFromLocalScorePasses(
		model, distanceTable, featureTable, passes, layerOptions, transitions,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	return Paul2013PhoneMarkerPathResult{Prepared: shortlists.Prepared, Groups: shortlists.Groups, Path: path}, nil
}

// SelectPaul2013PhoneMarkerRecordPathFromTransitionStates composes the
// six-byte state mode extraction with the exact-pool selection path. Transition
// states must align with flattened candidate positions; scoring, shortlist,
// and per-edge path inputs remain explicit.
func SelectPaul2013PhoneMarkerRecordPathFromTransitionStates(
	model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	},
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	input Paul2013PhoneMarkerRecordCandidates,
	states []Paul2013TransitionContextState,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PhoneMarkerPathResult, error) {
	positionCountFromStates := 0
	for _, group := range input.Groups {
		positionCountFromStates += len(group.CandidatePools)
	}
	if len(states) != positionCountFromStates {
		return Paul2013PhoneMarkerPathResult{}, fmt.Errorf("received %d transition states for %d candidate positions", len(states), positionCountFromStates)
	}
	modes, err := Paul2013ContextModesFromTransitionStates(states)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	return SelectPaul2013PhoneMarkerRecordPath(
		model, distanceTable, featureTable, input, modes, positionCount,
		scoringInputs, tailStarts, candidateClasses, contextClasses,
		layerOptions, transitions,
	)
}

// SelectPaul2013PhoneMarkerArenaPath composes the shared phone-state reset,
// record grouping, packed-context lookup, continuity scoring, shortlist
// finalization, and minimum-cost path for one complete contiguous record
// arena. Text-derived score/tail state and edge-transition state remain
// explicit caller inputs; exact lookup misses still require fallback first.
func SelectPaul2013PhoneMarkerArenaPath(
	model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	},
	catalog *voice.ClassCatalog,
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	arena []byte,
	recordCount int,
	specialKeys [][]byte,
	characterMap [256]byte,
	maxUnits uint64,
	modes []byte,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PhoneMarkerPathResult, error) {
	candidates, err := LookupPaul2013PhoneMarkerArenaCandidates(
		catalog, arena, recordCount, specialKeys, characterMap, maxUnits,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	return SelectPaul2013PhoneMarkerRecordPath(
		model, distanceTable, featureTable, candidates, modes, positionCount,
		scoringInputs, tailStarts, candidateClasses, contextClasses,
		layerOptions, transitions,
	)
}

// SelectPaul2013PhoneMarkerArenaPathWithPointerResolver is the full arena
// composition when the native record stream retains its +0x3b8 key pointers.
// Pointer resolution remains bounded and owned by the caller.
func SelectPaul2013PhoneMarkerArenaPathWithPointerResolver(
	model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	},
	catalog *voice.ClassCatalog,
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	arena []byte,
	recordCount int,
	resolvePointer text.Paul2013CStringPointerResolver,
	characterMap [256]byte,
	maxUnits uint64,
	modes []byte,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PhoneMarkerPathResult, error) {
	candidates, err := LookupPaul2013PhoneMarkerArenaCandidatesWithPointerResolver(
		catalog, arena, recordCount, resolvePointer, characterMap, maxUnits,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	return SelectPaul2013PhoneMarkerRecordPath(
		model, distanceTable, featureTable, candidates, modes, positionCount,
		scoringInputs, tailStarts, candidateClasses, contextClasses,
		layerOptions, transitions,
	)
}

// SelectPaul2013PhoneMarkerArenaPathFromTransitionStates composes the arena
// reset and exact candidate lookup with transition-state mode extraction and
// the record-level path selector. State rows must align with the flattened
// candidate positions; score contexts, shortlist metadata, and edge state
// remain explicit inputs.
func SelectPaul2013PhoneMarkerArenaPathFromTransitionStates(
	model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	},
	catalog *voice.ClassCatalog,
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	arena []byte,
	recordCount int,
	specialKeys [][]byte,
	characterMap [256]byte,
	maxUnits uint64,
	states []Paul2013TransitionContextState,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PhoneMarkerPathResult, error) {
	candidates, err := LookupPaul2013PhoneMarkerArenaCandidates(
		catalog, arena, recordCount, specialKeys, characterMap, maxUnits,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	return SelectPaul2013PhoneMarkerRecordPathFromTransitionStates(
		model, distanceTable, featureTable, candidates, states,
		positionCount, scoringInputs, tailStarts, candidateClasses,
		contextClasses, layerOptions, transitions,
	)
}

// SelectPaul2013PhoneMarkerArenaPathWithPointerResolverFromTransitionStates
// is the pointer-backed counterpart of
// SelectPaul2013PhoneMarkerArenaPathFromTransitionStates. It retains the same
// alignment and explicit downstream-state requirements.
func SelectPaul2013PhoneMarkerArenaPathWithPointerResolverFromTransitionStates(
	model interface {
		Paul2013ContinuityModel
		Paul2013PathModel
	},
	catalog *voice.ClassCatalog,
	distanceTable *distance.Table,
	featureTable *distance.FeatureTable,
	arena []byte,
	recordCount int,
	resolvePointer text.Paul2013CStringPointerResolver,
	characterMap [256]byte,
	maxUnits uint64,
	states []Paul2013TransitionContextState,
	positionCount int,
	scoringInputs []Paul2013QueriedCandidateScoringInput,
	tailStarts []int,
	candidateClasses, contextClasses [][]byte,
	layerOptions []Paul2013PathLayerOptions,
	transitions []Paul2013PathTransitionContext,
) (Paul2013PhoneMarkerPathResult, error) {
	candidates, err := LookupPaul2013PhoneMarkerArenaCandidatesWithPointerResolver(
		catalog, arena, recordCount, resolvePointer, characterMap, maxUnits,
	)
	if err != nil {
		return Paul2013PhoneMarkerPathResult{}, err
	}
	return SelectPaul2013PhoneMarkerRecordPathFromTransitionStates(
		model, distanceTable, featureTable, candidates, states,
		positionCount, scoringInputs, tailStarts, candidateClasses,
		contextClasses, layerOptions, transitions,
	)
}

func lookupPaul2013PhoneMarkerRecordCandidates(
	catalog *voice.ClassCatalog,
	prepared text.Paul2013PreparedRecordGroups,
	maxUnits uint64,
) (Paul2013PhoneMarkerRecordCandidates, error) {
	if len(prepared.ContextGroups) != len(prepared.Descriptors) {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("prepared %d context groups for %d descriptors", len(prepared.ContextGroups), len(prepared.Descriptors))
	}
	if len(prepared.SelectionContexts) != len(prepared.ContextGroups) {
		return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("prepared %d selection-context groups for %d context groups", len(prepared.SelectionContexts), len(prepared.ContextGroups))
	}
	result := Paul2013PhoneMarkerRecordCandidates{
		Prepared: prepared,
		Groups:   make([]Paul2013PhoneMarkerGroupCandidates, len(prepared.ContextGroups)),
	}
	for groupIndex, group := range prepared.ContextGroups {
		contexts := prepared.SelectionContexts[groupIndex]
		if len(contexts) != len(group.Categories.Positions) {
			return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("prepared selection-context group %d has %d contexts for %d packed phone positions", groupIndex, len(contexts), len(group.Categories.Positions))
		}
		groupResult := Paul2013PhoneMarkerGroupCandidates{
			Descriptor:      prepared.Descriptors[groupIndex],
			Contexts:        contexts,
			CandidatePools:  make([][]UnitRef, len(contexts)),
			ExactClassFound: make([]bool, len(contexts)),
		}
		for contextIndex, context := range contexts {
			units, found, err := ExactContextUnitCandidates(catalog, context, maxUnits)
			if err != nil {
				return Paul2013PhoneMarkerRecordCandidates{}, fmt.Errorf("lookup record group %d context %d: %w", groupIndex, contextIndex, err)
			}
			groupResult.CandidatePools[contextIndex] = units
			groupResult.ExactClassFound[contextIndex] = found
		}
		result.Groups[groupIndex] = groupResult
	}
	return result, nil
}
