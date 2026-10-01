package selection

import (
	"errors"
	"fmt"
	"sort"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

const (
	paul2013WholePositionMaxClasses = 30
	paul2013WholePositionMaxUnits   = 10_000
	paul2013PrefixQueryMaxClasses   = 300
	paul2013PrefixQueryMinUnits     = 30
	paul2013FeatureViewMaxClasses   = 10
	paul2013FeatureViewMinUnits     = 10
	paul2013FeatureViewInitialWidth = 10
	paul2013FeatureViewMinRange     = 2
)

// Paul2013WholePositionResult is the deterministic post-lookup result of
// FUN_10024060. Classes contains the ranked, deduplicated shortlist and
// MetricSubtotal is the sum of their member counts.
type Paul2013WholePositionResult struct {
	Classes              []voice.ClassRecord
	QueryHitCount        int
	LookupCandidateCount int
	MetricSubtotal       uint32
	LookupPrefixLength   int
	LookupMetricSubtotal uint64
	Accepted             bool
}

// Paul2013ContextClassLookup is implemented by voice.ClassCatalog. It is
// separated here so the query-to-shortlist composition can be tested without
// opening proprietary model files.
type Paul2013ContextClassLookup interface {
	LookupContext(text.Context) (voice.ClassRecord, bool)
}

// Paul2013PrefixClassLookup is implemented by voice.ClassCatalog. It exposes
// the variable-width key comparison used by FUN_10019450.
type Paul2013PrefixClassLookup interface {
	LookupPrefixLimited([5]byte, int, int) ([]voice.ClassRecord, error)
}

// Paul2013FeatureViewClassLookup provides a class-ID vector sorted by one of
// the native 10-byte feature views. voice.ClassCatalog reconstructs those
// vectors from the local unit indexes.
type Paul2013FeatureViewClassLookup interface {
	LookupFeatureViewPrefix([10]byte, byte, int, int) ([]voice.ClassRecord, error)
}

// Paul2013FeatureViewLookupResult contains the ranked classes and the member
// coverage used by the native ten-unit feature-view acceptance check.
type Paul2013FeatureViewLookupResult struct {
	Classes              []voice.ClassRecord
	QueryHitCount        int
	LookupCandidateCount int
	MetricSubtotal       uint64
	LookupPrefixLength   int
	Accepted             bool
}

// LookupPaul2013WholePosition looks up each already-generated query variant
// in order, then canonicalizes the found class IDs as FUN_10024010 does before
// passing them through the FUN_10024060 shortlist and acceptance path. target
// is the unmodified tree-leaf context used as the ranking key. The function
// does not produce query variants or fallback rows.
func LookupPaul2013WholePosition(
	catalog Paul2013ContextClassLookup,
	target text.Context,
	queryVariants []text.Context,
	rowClass byte,
) (Paul2013WholePositionResult, error) {
	if catalog == nil {
		return Paul2013WholePositionResult{}, errors.New("whole-position lookup has no class catalog")
	}
	if len(queryVariants) == 0 {
		return Paul2013WholePositionResult{}, errors.New("whole-position lookup has no query variants")
	}
	candidates := make([]voice.ClassRecord, 0, len(queryVariants))
	for index, query := range queryVariants {
		class, found := catalog.LookupContext(query)
		if !found {
			continue
		}
		if class.ID > uint32(^uint16(0)) {
			return Paul2013WholePositionResult{}, fmt.Errorf("query variant %d resolves to class ID %d outside native 16-bit range", index, class.ID)
		}
		candidates = append(candidates, class)
	}
	result, err := EvaluatePaul2013WholePosition(target.Paul2013SelectionKey(), candidates, rowClass)
	if err != nil {
		return Paul2013WholePositionResult{}, err
	}
	result.QueryHitCount = len(candidates)
	result.LookupCandidateCount = len(candidates)
	return result, nil
}

// LookupPaul2013WholePositionByPrefix ports FUN_10023dc0's bounded class-table
// broadening for one already-produced five-byte query key. It queries widths
// 5 through 1, appending each matching range until the accumulated member
// population reaches 30 or the native 300-class buffer fills. Existing
// candidates model entries written by an earlier state-generated query. The
// duplicated records from nested prefix ranges are retained for the native
// coverage calculation and canonicalized by EvaluatePaul2013WholePosition.
// Query-variant generation, row-state dispatch, and the feature-view lookup
// path remain caller responsibilities.
func LookupPaul2013WholePositionByPrefix(
	catalog Paul2013PrefixClassLookup,
	target text.Context,
	existing []voice.ClassRecord,
	rowClass byte,
) (Paul2013WholePositionResult, error) {
	if catalog == nil {
		return Paul2013WholePositionResult{}, errors.New("whole-position prefix lookup has no class catalog")
	}
	candidates := append([]voice.ClassRecord(nil), existing...)
	if len(candidates) > paul2013PrefixQueryMaxClasses {
		return Paul2013WholePositionResult{}, fmt.Errorf("existing candidate list has %d classes, exceeding native capacity %d", len(candidates), paul2013PrefixQueryMaxClasses)
	}
	var coverage uint64
	for _, class := range candidates {
		if class.ID > uint32(^uint16(0)) {
			return Paul2013WholePositionResult{}, fmt.Errorf("existing candidate has class ID %d outside native 16-bit range", class.ID)
		}
		population := uint64(len(class.Members))
		if ^uint64(0)-coverage < population {
			return Paul2013WholePositionResult{}, errors.New("native class metric subtotal overflows uint64")
		}
		coverage += population
	}
	lastPrefixLength := 0
	queryHits := 0
	key := target.Paul2013SelectionKey()
	for prefixLength := len(key); prefixLength >= 1 && len(candidates) < paul2013PrefixQueryMaxClasses; prefixLength-- {
		// The seven-byte context is projected to the native five-byte key.
		remaining := paul2013PrefixQueryMaxClasses - len(candidates)
		matches, err := catalog.LookupPrefixLimited(key, prefixLength, remaining)
		if err != nil {
			return Paul2013WholePositionResult{}, fmt.Errorf("query %d-byte class prefix: %w", prefixLength, err)
		}
		queryHits += len(matches)
		for _, class := range matches {
			if class.ID > uint32(^uint16(0)) {
				return Paul2013WholePositionResult{}, fmt.Errorf("prefix query resolves to class ID %d outside native 16-bit range", class.ID)
			}
			population := uint64(len(class.Members))
			if ^uint64(0)-coverage < population {
				return Paul2013WholePositionResult{}, errors.New("native class metric subtotal overflows uint64")
			}
			coverage += population
			candidates = append(candidates, class)
		}
		lastPrefixLength = prefixLength
		if coverage >= paul2013PrefixQueryMinUnits || len(candidates) >= paul2013PrefixQueryMaxClasses {
			break
		}
	}
	result, err := EvaluatePaul2013WholePosition(target.Paul2013SelectionKey(), candidates, rowClass)
	if err != nil {
		return Paul2013WholePositionResult{}, err
	}
	result.QueryHitCount = queryHits
	result.LookupCandidateCount = len(candidates)
	result.LookupPrefixLength = lastPrefixLength
	result.LookupMetricSubtotal = coverage
	return result, nil
}

// LookupPaul2013WholePositionByFeatureView is a compatibility wrapper for
// callers of the earlier helper. Its rowClass-based query-mode assumption is
// retained for that API; native callers should pass the actual mode to
// LookupPaul2013WholePositionByFeatureViewWithQueryMode.
func LookupPaul2013WholePositionByFeatureView(
	candidatePool []voice.ClassRecord,
	target, query text.Context,
	viewMode byte,
	rowClass byte,
) (Paul2013FeatureViewLookupResult, error) {
	queryMode := byte(2)
	if rowClass == 1 {
		queryMode = 1
	}
	return lookupPaul2013FeatureView(candidatePool, target, query, viewMode, queryMode)
}

// LookupPaul2013WholePositionByFeatureViewWithQueryMode ports prefix
// relaxation and shortlisting over one caller-supplied class scope from
// FUN_10019570, FUN_1002df50, and FUN_10023c70. This helper preserves the
// supplied scope order on the native fast path. Callers with a ClassCatalog
// can use LookupPaul2013WholePositionByFeatureViewCatalog to obtain ranges
// from its reconstructed full view index. Native range counts still disagree
// with captured examples, so Go key/view parity remains open. The helper
// starts at ten view bytes and drops one byte per retry. FUN_10023e70 permits
// widths down to 1 only when its query mode is 1; modes 0 and 2 stop at width
// 6. The row view selects the feature projection. Each matching range is
// ranked to at most 10 classes/10,000 members and returns once its metric
// reaches ten; scoring is capped at 10,000 candidates.
func LookupPaul2013WholePositionByFeatureViewWithQueryMode(
	candidatePool []voice.ClassRecord,
	target, query text.Context,
	rowClass byte,
	queryMode byte,
) (Paul2013FeatureViewLookupResult, error) {
	if rowClass != 1 && rowClass != 2 {
		return Paul2013FeatureViewLookupResult{}, fmt.Errorf("unsupported Paul 2013 feature row class %d", rowClass)
	}
	return lookupPaul2013FeatureView(candidatePool, target, query, rowClass, queryMode)
}

// LookupPaul2013WholePositionByFeatureViewCatalog applies the native
// ten-to-one-byte range relaxation directly to the reconstructed class view
// index. The returned query hit count is capped at 10,000 because that is the
// maximum candidate prefix FUN_10023c70 scores.
func LookupPaul2013WholePositionByFeatureViewCatalog(
	catalog Paul2013FeatureViewClassLookup,
	target, query text.Context,
	viewMode byte,
	queryMode byte,
) (Paul2013FeatureViewLookupResult, error) {
	if catalog == nil {
		return Paul2013FeatureViewLookupResult{}, errors.New("feature-view lookup has no class catalog")
	}
	if viewMode != 1 && viewMode != 2 {
		return Paul2013FeatureViewLookupResult{}, fmt.Errorf("unsupported Paul 2013 feature-view mode %d", viewMode)
	}
	if queryMode > 2 {
		return Paul2013FeatureViewLookupResult{}, fmt.Errorf("unsupported Paul 2013 query mode %d", queryMode)
	}
	view, err := query.Paul2013FeatureView(viewMode)
	if err != nil {
		return Paul2013FeatureViewLookupResult{}, err
	}
	minimumWidth := 6
	if queryMode == 1 {
		minimumWidth = 1
	}
	var result Paul2013FeatureViewLookupResult
	for prefixLength := paul2013FeatureViewInitialWidth; prefixLength >= minimumWidth; prefixLength-- {
		matches, err := catalog.LookupFeatureViewPrefix(view, viewMode, prefixLength, paul2013MaxScoredClassCandidates)
		if err != nil {
			return Paul2013FeatureViewLookupResult{}, fmt.Errorf("query %d-byte feature-view prefix: %w", prefixLength, err)
		}
		current := Paul2013FeatureViewLookupResult{
			LookupPrefixLength: prefixLength,
			QueryHitCount:      len(matches),
		}
		if len(matches) < paul2013FeatureViewMinRange {
			result = current
			continue
		}
		ranked, err := RankClassRecords(target.Paul2013SelectionKey(), matches, RankOptions{
			MaxClasses: paul2013FeatureViewMaxClasses,
			MaxUnits:   paul2013WholePositionMaxUnits,
			ViewMode:   viewMode,
		})
		if err != nil {
			return Paul2013FeatureViewLookupResult{}, err
		}
		var subtotal uint64
		for _, class := range ranked {
			population := uint64(len(class.Members))
			if ^uint64(0)-subtotal < population {
				return Paul2013FeatureViewLookupResult{}, errors.New("native feature-view metric subtotal overflows uint64")
			}
			subtotal += population
		}
		current.Classes = ranked
		current.LookupCandidateCount = len(ranked)
		current.MetricSubtotal = subtotal
		result = current
		if subtotal >= paul2013FeatureViewMinUnits {
			result.Accepted = true
			return result, nil
		}
	}
	return result, nil
}

func lookupPaul2013FeatureView(
	candidatePool []voice.ClassRecord,
	target, query text.Context,
	viewMode byte,
	queryMode byte,
) (Paul2013FeatureViewLookupResult, error) {
	if viewMode != 1 && viewMode != 2 {
		return Paul2013FeatureViewLookupResult{}, fmt.Errorf("unsupported Paul 2013 feature-view mode %d", viewMode)
	}
	if queryMode > 2 {
		return Paul2013FeatureViewLookupResult{}, fmt.Errorf("unsupported Paul 2013 query mode %d", queryMode)
	}
	view, err := query.Paul2013FeatureView(viewMode)
	if err != nil {
		return Paul2013FeatureViewLookupResult{}, err
	}
	minimumWidth := 6
	if queryMode == 1 {
		minimumWidth = 1
	}
	allViews := make([]scoredViewClass, len(candidatePool))
	for index, class := range candidatePool {
		candidateView, err := text.Paul2013FeatureViewForKey(class.Key, viewMode)
		if err != nil {
			return Paul2013FeatureViewLookupResult{}, fmt.Errorf("build feature view for class %d: %w", class.ID, err)
		}
		allViews[index] = scoredViewClass{class: class, view: candidateView}
	}
	result := Paul2013FeatureViewLookupResult{}
	for prefixLength := paul2013FeatureViewInitialWidth; prefixLength >= minimumWidth; prefixLength-- {
		matched := make([]voice.ClassRecord, 0)
		matchCount := 0
		for _, candidate := range allViews {
			if equalFeatureViewPrefix(candidate.view, view, prefixLength) {
				matchCount++
				if len(matched) < paul2013MaxScoredClassCandidates {
					matched = append(matched, candidate.class)
				}
			}
		}
		current := Paul2013FeatureViewLookupResult{
			LookupPrefixLength: prefixLength,
			QueryHitCount:      matchCount,
		}
		if matchCount < paul2013FeatureViewMinRange {
			result = current
			continue
		}
		ranked, err := RankClassRecords(target.Paul2013SelectionKey(), matched, RankOptions{
			MaxClasses: paul2013FeatureViewMaxClasses,
			MaxUnits:   paul2013WholePositionMaxUnits,
			ViewMode:   viewMode,
		})
		if err != nil {
			return Paul2013FeatureViewLookupResult{}, err
		}
		var subtotal uint64
		for _, class := range ranked {
			population := uint64(len(class.Members))
			if ^uint64(0)-subtotal < population {
				return Paul2013FeatureViewLookupResult{}, errors.New("native feature-view metric subtotal overflows uint64")
			}
			subtotal += population
		}
		current.Classes = ranked
		current.LookupCandidateCount = len(ranked)
		current.MetricSubtotal = subtotal
		result = current
		if subtotal >= paul2013FeatureViewMinUnits {
			result.Accepted = true
			return result, nil
		}
	}
	return result, nil
}

type scoredViewClass struct {
	class voice.ClassRecord
	view  [10]byte
}

func equalFeatureViewPrefix(left, right [10]byte, prefixLength int) bool {
	for position := 0; position < prefixLength; position++ {
		if left[position] != right[position] {
			return false
		}
	}
	return true
}

// EvaluatePaul2013WholePosition ports the post-query class-list and acceptance
// portion of FUN_10024060. It first reproduces FUN_10024010's numeric class-ID
// sort and duplicate removal, then ranks/shortlists those records. Native
// class weights equal member populations, so the subtotal is derived from
// member counts. Query generation and model-state-dependent lookups remain
// caller responsibilities.
func EvaluatePaul2013WholePosition(
	contextKey [5]byte,
	candidates []voice.ClassRecord,
	rowClass byte,
) (Paul2013WholePositionResult, error) {
	ordered := append([]voice.ClassRecord(nil), candidates...)
	sort.Slice(ordered, func(left, right int) bool {
		return ordered[left].ID < ordered[right].ID
	})
	unique := make([]voice.ClassRecord, 0, len(ordered))
	seen := make(map[uint32]voice.ClassRecord, len(candidates))
	for index, candidate := range ordered {
		if previous, exists := seen[candidate.ID]; exists {
			if previous.Key != candidate.Key || !sameUnitLocations(previous.Members, candidate.Members) {
				return Paul2013WholePositionResult{}, fmt.Errorf("candidate class %d has conflicting records at index %d", candidate.ID, index)
			}
			continue
		}
		seen[candidate.ID] = candidate
		unique = append(unique, candidate)
	}

	ranked, err := RankClassRecords(contextKey, unique, RankOptions{
		MaxClasses: paul2013WholePositionMaxClasses,
		MaxUnits:   paul2013WholePositionMaxUnits,
		ViewMode:   0,
	})
	if err != nil {
		return Paul2013WholePositionResult{}, err
	}
	var subtotal uint32
	for _, class := range ranked {
		if uint64(len(class.Members)) > uint64(^uint32(0)) {
			return Paul2013WholePositionResult{}, fmt.Errorf("candidate class %d population exceeds supported width", class.ID)
		}
		weight := uint32(len(class.Members))
		if ^uint32(0)-subtotal < weight {
			return Paul2013WholePositionResult{}, errors.New("native class metric subtotal overflows uint32")
		}
		subtotal += weight
	}
	accepted := subtotal > 9
	if rowClass == 8 {
		accepted = subtotal > 0
	}
	return Paul2013WholePositionResult{Classes: ranked, MetricSubtotal: subtotal, Accepted: accepted}, nil
}

func sameUnitLocations(left, right []voice.UnitLocation) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
