package selection

import (
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

// Paul2013WholePositionQueryResult records the query passes and final
// whole-position shortlist used by FUN_10024060.
type Paul2013WholePositionQueryResult struct {
	Passes         []Paul2013CatalogQuerySequenceResult
	Selection      Paul2013WholePositionResult
	CandidateCount int
}

// QueryPaul2013WholePosition connects FUN_10024060's query dispatch to
// FUN_10018770, FUN_10023f90, and the post-lookup shortlist. Whole-position
// rows force context byte +4 to zero. Model class 12 first queries a copy
// whose signature byte +6 has bit 0x20 cleared, then queries the original in
// mode two while retaining the first pass's candidate list. Other model
// classes query the original once in mode one.
func QueryPaul2013WholePosition(
	catalog Paul2013PrefixClassLookup,
	target text.Context,
	contextRow [6]byte,
	priorMetric uint32,
	modelClass byte,
	contextGate func() byte,
) (Paul2013WholePositionQueryResult, error) {
	queryRow := contextRow
	queryRow[4] = 0
	result := Paul2013WholePositionQueryResult{Passes: make([]Paul2013CatalogQuerySequenceResult, 0, 2)}
	signature := target.Signature
	initialCandidates := []voice.ClassRecord(nil)
	queryMode := int16(1)
	if modelClass == 12 {
		masked := signature
		masked[6] &= 0xdf
		pass, err := RunPaul2013CatalogQuerySequence(
			catalog, masked, target, queryRow, nil, priorMetric, 1,
			contextGate, nil,
		)
		result.Passes = append(result.Passes, pass)
		if err != nil {
			return result, fmt.Errorf("query masked model-class-12 signature: %w", err)
		}
		if len(pass.Candidates) == 0 {
			return result, nil
		}
		initialCandidates = pass.Candidates
		queryMode = 2
	}

	pass, err := RunPaul2013CatalogQuerySequence(
		catalog, signature, target, queryRow, initialCandidates, priorMetric,
		queryMode, contextGate, nil,
	)
	result.Passes = append(result.Passes, pass)
	if err != nil {
		return result, fmt.Errorf("query whole-position signature: %w", err)
	}
	result.CandidateCount = len(pass.Candidates)
	if result.CandidateCount == 0 {
		return result, nil
	}
	selection, err := EvaluatePaul2013WholePosition(
		target.Paul2013SelectionKey(), pass.Candidates, modelClass,
	)
	if err != nil {
		return result, fmt.Errorf("evaluate whole-position candidates: %w", err)
	}
	result.Selection = selection
	return result, nil
}
