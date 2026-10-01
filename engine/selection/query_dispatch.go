package selection

import (
	"errors"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

const (
	paul2013QueryClassCapacity    = 300
	paul2013FeatureQuerySoftLimit = 0x122
)

// Paul2013QueryDispatchResult is the class-list delta returned by one
// FUN_10023f90 call. Appended preserves native order and repeated matches;
// CandidateCount is the signed-short-compatible return value for that call.
type Paul2013QueryDispatchResult struct {
	Appended        []voice.ClassRecord
	CandidateCount  int16
	QueryHitCount   int
	MetricSubtotal  uint64
	PrefixLength    int
	CoverageReached bool
}

// Paul2013FeatureScopeProvider overrides catalog feature-range lookup for
// controlled comparisons. The normal path reads the reconstructed feature
// view index directly from the class catalog.
type Paul2013FeatureScopeProvider func(
	querySignature [7]byte,
	rowClass byte,
	queryMode byte,
) ([]voice.ClassRecord, error)

// Paul2013CatalogQuerySequenceResult combines the ordered FUN_10018770 query
// trace with the class records appended by its FUN_10023f90 lookups. On a
// mode-1 sequence that starts with no candidates and has a positive lookup, it
// also records the first 30 IDs copied to the fallback raw-candidate sidecar.
type Paul2013CatalogQuerySequenceResult struct {
	Sequence                Paul2013QuerySequenceResult
	Candidates              []voice.ClassRecord
	FallbackRawCandidateIDs []uint32
}

// RunPaul2013CatalogQuerySequence connects the query producer to the class
// catalog. Whole-position queries use the five-byte class index; feature-view
// queries use the reconstructed ten-byte catalog indexes unless featureScope
// supplies an explicit comparison pool. initialCandidates continues an
// earlier FUN_10018770 pass such as the class-12 two-mode path. target is the
// context used by the native feature ranker.
func RunPaul2013CatalogQuerySequence(
	catalog Paul2013PrefixClassLookup,
	signature [7]byte,
	target text.Context,
	contextRow [6]byte,
	initialCandidates []voice.ClassRecord,
	priorMetric uint32,
	mode int16,
	contextGate func() byte,
	featureScope Paul2013FeatureScopeProvider,
) (Paul2013CatalogQuerySequenceResult, error) {
	result := Paul2013CatalogQuerySequenceResult{
		Candidates: append([]voice.ClassRecord(nil), initialCandidates...),
	}
	lookup := func(querySignature [7]byte, queryMode int16) (int16, error) {
		rowClass := contextRow[4]
		if rowClass != 0 && featureScope == nil {
			featureCatalog, ok := catalog.(Paul2013FeatureViewClassLookup)
			if !ok {
				return 0, errors.New("feature-view query catalog has no reconstructed feature-range index")
			}
			dispatch, err := lookupPaul2013CatalogFeatureQuery(
				featureCatalog, querySignature, target, rowClass, byte(queryMode), result.Candidates,
			)
			if err != nil {
				return 0, err
			}
			result.Candidates = append(result.Candidates, dispatch.Appended...)
			return dispatch.CandidateCount, nil
		}
		var scope []voice.ClassRecord
		if rowClass != 0 {
			var err error
			scope, err = featureScope(querySignature, rowClass, byte(queryMode))
			if err != nil {
				return 0, fmt.Errorf("provide feature-query comparison scope: %w", err)
			}
		}
		dispatch, err := LookupPaul2013QueryDispatch(
			catalog,
			querySignature,
			target,
			rowClass,
			byte(queryMode),
			result.Candidates,
			scope,
		)
		if err != nil {
			return 0, err
		}
		result.Candidates = append(result.Candidates, dispatch.Appended...)
		return dispatch.CandidateCount, nil
	}

	sequence, err := RunPaul2013QuerySequence(
		signature,
		contextRow,
		priorMetric,
		mode,
		contextGate,
		lookup,
	)
	result.Sequence = sequence
	if err != nil {
		return result, err
	}
	if mode == 1 && len(initialCandidates) == 0 && paul2013QuerySequenceHasPositiveResult(sequence) {
		limit := len(result.Candidates)
		if limit > paul2013FallbackRawCandidateLimit {
			limit = paul2013FallbackRawCandidateLimit
		}
		result.FallbackRawCandidateIDs = make([]uint32, limit)
		for index := 0; index < limit; index++ {
			result.FallbackRawCandidateIDs[index] = result.Candidates[index].ID
		}
	}
	return result, nil
}

func lookupPaul2013CatalogFeatureQuery(
	catalog Paul2013FeatureViewClassLookup,
	querySignature [7]byte,
	target text.Context,
	viewMode byte,
	queryMode byte,
	currentCandidates []voice.ClassRecord,
) (Paul2013QueryDispatchResult, error) {
	if len(currentCandidates) >= paul2013FeatureQuerySoftLimit {
		return Paul2013QueryDispatchResult{}, nil
	}
	feature, err := LookupPaul2013WholePositionByFeatureViewCatalog(
		catalog, target, text.Context{Signature: querySignature}, viewMode, queryMode,
	)
	if err != nil {
		return Paul2013QueryDispatchResult{}, fmt.Errorf("lookup feature-view query in class catalog: %w", err)
	}
	if len(feature.Classes) > paul2013FeatureViewMaxClasses {
		return Paul2013QueryDispatchResult{}, fmt.Errorf("feature query returned %d classes, exceeding native per-query limit %d", len(feature.Classes), paul2013FeatureViewMaxClasses)
	}
	return Paul2013QueryDispatchResult{
		Appended:        feature.Classes,
		CandidateCount:  int16(len(feature.Classes)),
		QueryHitCount:   feature.QueryHitCount,
		MetricSubtotal:  feature.MetricSubtotal,
		PrefixLength:    feature.LookupPrefixLength,
		CoverageReached: feature.MetricSubtotal >= paul2013FeatureViewMinUnits,
	}, nil
}

func paul2013QuerySequenceHasPositiveResult(sequence Paul2013QuerySequenceResult) bool {
	for _, attempt := range sequence.Attempts {
		if attempt.Result > 0 {
			return true
		}
	}
	return false
}

// LookupPaul2013QueryDispatch ports the row-class dispatch in
// FUN_10023f90. A zero row class uses FUN_10023dc0 to append nested
// five-to-one-byte class-prefix matches, stopping at 30 population or the
// 300-class buffer. A nonzero row class uses FUN_10023e70 over the explicitly
// supplied state-selected feature scope; its mode-1 query may relax to one
// byte, while modes 0 and 2 stop at six. The source producer for that scope
// remains unresolved.
func LookupPaul2013QueryDispatch(
	catalog Paul2013PrefixClassLookup,
	querySignature [7]byte,
	target text.Context,
	rowClass byte,
	queryMode byte,
	currentCandidates []voice.ClassRecord,
	featureScope []voice.ClassRecord,
) (Paul2013QueryDispatchResult, error) {
	if rowClass == 0 {
		return lookupPaul2013WholeQuery(catalog, querySignature, currentCandidates)
	}
	if len(currentCandidates) >= paul2013FeatureQuerySoftLimit {
		return Paul2013QueryDispatchResult{}, nil
	}
	feature, err := LookupPaul2013WholePositionByFeatureViewWithQueryMode(
		featureScope,
		target,
		text.Context{Signature: querySignature},
		rowClass,
		queryMode,
	)
	if err != nil {
		return Paul2013QueryDispatchResult{}, err
	}
	if len(feature.Classes) > paul2013FeatureViewMaxClasses {
		return Paul2013QueryDispatchResult{}, fmt.Errorf("feature query returned %d classes, exceeding native per-query limit %d", len(feature.Classes), paul2013FeatureViewMaxClasses)
	}
	return Paul2013QueryDispatchResult{
		Appended:        feature.Classes,
		CandidateCount:  int16(len(feature.Classes)),
		QueryHitCount:   feature.QueryHitCount,
		MetricSubtotal:  feature.MetricSubtotal,
		PrefixLength:    feature.LookupPrefixLength,
		CoverageReached: feature.MetricSubtotal >= paul2013FeatureViewMinUnits,
	}, nil
}

func lookupPaul2013WholeQuery(
	catalog Paul2013PrefixClassLookup,
	querySignature [7]byte,
	currentCandidates []voice.ClassRecord,
) (Paul2013QueryDispatchResult, error) {
	if catalog == nil {
		return Paul2013QueryDispatchResult{}, errors.New("whole-position query has no class catalog")
	}
	if len(currentCandidates) > paul2013QueryClassCapacity {
		return Paul2013QueryDispatchResult{}, fmt.Errorf("current class list has %d entries, exceeding native capacity %d", len(currentCandidates), paul2013QueryClassCapacity)
	}
	coverage, err := paul2013ClassPopulationSubtotal(currentCandidates)
	if err != nil {
		return Paul2013QueryDispatchResult{}, fmt.Errorf("sum current query class population: %w", err)
	}
	key := (text.Context{Signature: querySignature}).Paul2013SelectionKey()
	result := Paul2013QueryDispatchResult{Appended: make([]voice.ClassRecord, 0)}
	for width := len(key); width >= 1 && len(currentCandidates)+len(result.Appended) < paul2013QueryClassCapacity; width-- {
		remaining := paul2013QueryClassCapacity - len(currentCandidates) - len(result.Appended)
		matches, err := catalog.LookupPrefixLimited(key, width, remaining)
		if err != nil {
			return result, fmt.Errorf("lookup %d-byte query prefix: %w", width, err)
		}
		result.QueryHitCount += len(matches)
		for _, class := range matches {
			if class.ID > uint32(^uint16(0)) {
				return result, fmt.Errorf("query prefix class ID %d exceeds native 16-bit range", class.ID)
			}
			population := uint64(len(class.Members))
			if ^uint64(0)-coverage < population {
				return result, errors.New("native query class metric subtotal overflows uint64")
			}
			coverage += population
			result.Appended = append(result.Appended, class)
		}
		result.PrefixLength = width
		if coverage >= paul2013PrefixQueryMinUnits {
			result.CoverageReached = true
			break
		}
	}
	result.CandidateCount = int16(len(result.Appended))
	result.MetricSubtotal = coverage
	return result, nil
}

func paul2013ClassPopulationSubtotal(classes []voice.ClassRecord) (uint64, error) {
	var subtotal uint64
	for _, class := range classes {
		population := uint64(len(class.Members))
		if ^uint64(0)-subtotal < population {
			return 0, errors.New("class population subtotal overflows uint64")
		}
		subtotal += population
	}
	return subtotal, nil
}
