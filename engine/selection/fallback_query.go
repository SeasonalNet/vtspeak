package selection

import (
	"fmt"
	"sort"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

// Paul2013FallbackRowRuntime contains values read from synthesis state while
// FUN_100242a0 constructs and ranks one fallback row. The state-table
// producers are not reconstructed, so callers provide the live query gate,
// prior metric and feature scope explicitly.
type Paul2013FallbackRowRuntime struct {
	PriorMetric uint32
	// ReadPriorMetric rereads the native state byte before each tree query.
	// When nil, PriorMetric is used as a snapshot.
	ReadPriorMetric func() uint32
	ContextGate     func() byte
	// FeatureScope overrides the class catalog's range index for controlled
	// comparisons. Leave it nil to query the reconstructed catalog directly.
	FeatureScope Paul2013FeatureScopeProvider
}

// Paul2013FallbackRowResult is the portable result for one of the two rows
// emitted by FUN_100242a0. QueryPasses records each FUN_10018770 invocation;
// RankedClasses retains FUN_10023c70 score order, while Candidates mirrors
// the native ranked-prefix plus raw-sidecar row storage.
type Paul2013FallbackRowResult struct {
	Row             Paul2013FallbackQueryRow
	QueryPasses     []Paul2013CatalogQuerySequenceResult
	QueryClassCount int
	RankedClasses   []voice.ClassRecord
	Candidates      Paul2013FallbackCandidateList
}

// QueryPaul2013FallbackRows connects the recovered row builder, query
// dispatcher, shortlist ranker, and final sidecar merge for both fallback
// rows. State-derived values that are not recovered remain explicit in
// runtime. The mode-1 query sequence captures the raw sidecar under the
// observed positive-result/empty-initial-list gate.
func QueryPaul2013FallbackRows(
	catalog Paul2013PrefixClassLookup,
	target text.Context,
	baseRow [6]byte,
	contextDelta int8,
	modelRowClass byte,
	runtime [2]Paul2013FallbackRowRuntime,
) ([2]Paul2013FallbackRowResult, error) {
	rows := BuildPaul2013FallbackRows(baseRow, contextDelta, modelRowClass, target.Signature)
	var results [2]Paul2013FallbackRowResult
	for rowIndex, row := range rows {
		current := Paul2013FallbackRowResult{
			Row:         row,
			QueryPasses: make([]Paul2013CatalogQuerySequenceResult, 0, len(row.TreeTargets)),
		}
		candidates := make([]voice.ClassRecord, 0)
		var rawCandidateIDs []uint32
		for targetIndex, treeTarget := range row.TreeTargets {
			priorMetric := runtime[rowIndex].PriorMetric
			if runtime[rowIndex].ReadPriorMetric != nil {
				priorMetric = runtime[rowIndex].ReadPriorMetric()
			}
			pass, err := RunPaul2013CatalogQuerySequence(
				catalog,
				treeTarget.Signature,
				target,
				row.Bytes,
				candidates,
				priorMetric,
				int16(treeTarget.QueryMode),
				runtime[rowIndex].ContextGate,
				runtime[rowIndex].FeatureScope,
			)
			current.QueryPasses = append(current.QueryPasses, pass)
			if err != nil {
				return results, fmt.Errorf("fallback row %d tree target %d: %w", rowIndex, targetIndex, err)
			}
			candidates = pass.Candidates
			if treeTarget.QueryMode == 1 {
				rawCandidateIDs = pass.FallbackRawCandidateIDs
			}
		}
		ordered, err := canonicalPaul2013FallbackClasses(candidates)
		if err != nil {
			return results, fmt.Errorf("canonicalize fallback row %d classes: %w", rowIndex, err)
		}
		current.QueryClassCount = len(ordered)
		ranked, err := RankClassRecords(target.Paul2013SelectionKey(), ordered, RankOptions{
			MaxClasses: paul2013WholePositionMaxClasses,
			MaxUnits:   paul2013WholePositionMaxUnits,
			ViewMode:   row.Bytes[4],
		})
		if err != nil {
			return results, fmt.Errorf("rank fallback row %d classes: %w", rowIndex, err)
		}
		current.RankedClasses = ranked
		rankedIDs := make([]uint32, len(ranked))
		for index, class := range ranked {
			rankedIDs[index] = class.ID
		}
		merged, err := MergePaul2013FallbackCandidates(rankedIDs, rawCandidateIDs)
		if err != nil {
			return results, fmt.Errorf("assemble fallback row %d candidates: %w", rowIndex, err)
		}
		current.Candidates = merged
		results[rowIndex] = current
	}
	return results, nil
}

func canonicalPaul2013FallbackClasses(classes []voice.ClassRecord) ([]voice.ClassRecord, error) {
	ordered := append([]voice.ClassRecord(nil), classes...)
	for index, class := range ordered {
		if class.ID > uint32(^uint16(0)) {
			return nil, fmt.Errorf("class ID at index %d (%d) exceeds native 16-bit range", index, class.ID)
		}
	}
	sort.Slice(ordered, func(left, right int) bool { return ordered[left].ID < ordered[right].ID })
	unique := ordered[:0]
	for index, class := range ordered {
		if len(unique) != 0 && unique[len(unique)-1].ID == class.ID {
			previous := unique[len(unique)-1]
			if previous.Key != class.Key || !sameUnitLocations(previous.Members, class.Members) {
				return nil, fmt.Errorf("class ID %d has conflicting records at candidate index %d", class.ID, index)
			}
			continue
		}
		unique = append(unique, class)
	}
	return unique, nil
}
