package duration

import (
	"errors"

	"vtspeak/engine/text"
)

// Paul2013ProperNameTPPLookupResult distinguishes the caller's count gate
// from a dictionary miss. LookupPerformed is false for counts that bypass the
// native selector block, including five-component E records. Accepted means
// only that the known marker/value gate passed; ContextChecksBypassed and
// ContextChecksPending describe the later FUN_10034180 branch separately.
type Paul2013ProperNameTPPLookupResult struct {
	LookupPerformed        bool
	Selector               byte
	CompositeSurface       []byte
	Pattern                text.TPPComponentPattern
	Found                  bool
	Accepted               bool
	ContextScanInputKnown  bool
	ContextChecksBypassed  bool
	ContextChecksPending   bool
	ContextChecksEvaluated bool
	ContextGateAccepted    bool
	FollowupAttempted      bool
	FollowupSucceeded      bool
	SourceArena            []byte
}

// LookupTypedComponentPattern decodes the raw A-E component count/value
// payload from the loaded shared TPP dictionary. Numeric suffix and AX
// interpretations remain limited to the observed byte mechanics.
func (engine *Engine) LookupTypedComponentPattern(
	surface []byte,
) (text.TPPComponentPattern, bool, error) {
	if engine == nil || engine.tpp == nil {
		return text.TPPComponentPattern{}, false, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	return engine.tpp.LookupTPPComponentPattern(surface, text.Paul2013EmbeddedKeyTables())
}

// LookupProperNameTPPComponentPattern applies FUN_10034180's recovered
// one-through-four-component selector choice and then performs an exact TPP
// lookup. E-family records are not reached through this caller because its
// count gate bypasses selector construction at five components.
func (engine *Engine) LookupProperNameTPPComponentPattern(
	surface []byte,
	componentCount int,
) (Paul2013ProperNameTPPLookupResult, error) {
	selector, selected := text.Paul2013ProperNameTPPSelector(componentCount)
	result := Paul2013ProperNameTPPLookupResult{Selector: selector}
	if !selected {
		return result, nil
	}
	if engine == nil || engine.tpp == nil {
		return Paul2013ProperNameTPPLookupResult{}, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	pattern, found, err := engine.tpp.LookupTPPComponentPatternForSelector(
		surface, selector, text.Paul2013EmbeddedKeyTables(),
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	result.LookupPerformed = true
	result.Pattern = pattern
	result.Found = found
	return result, nil
}

// MatchProperNameTPPComponentPattern composes the native short-count
// selector, exact dictionary lookup, and FUN_10034180's marker/bit gate.
// componentMarkers are the per-component bytes observed at native offset
// +0x1a; the frontend that produces them remains a separate stage.
func (engine *Engine) MatchProperNameTPPComponentPattern(
	surface []byte,
	componentMarkers []byte,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.LookupProperNameTPPComponentPattern(surface, len(componentMarkers))
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	if !result.LookupPerformed {
		return result, nil
	}
	result.Accepted, err = text.AcceptPaul2013ProperNameTPPComponentPattern(
		result.Pattern, componentMarkers, result.Found,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return result, nil
}

// MatchProperNameTPPComponents composes FUN_10034180's bounded component
// surface builder, A-D selector, exact TPP lookup, and marker/value gate.
// Component surfaces, stored lengths, and +0x1a marker bytes come from the
// caller's already-produced native component rows.
func (engine *Engine) MatchProperNameTPPComponents(
	components []text.Paul2013ProperNameTPPComponent,
) (Paul2013ProperNameTPPLookupResult, error) {
	surface, eligible, err := text.BuildPaul2013ProperNameTPPSurface(components)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	if !eligible {
		return Paul2013ProperNameTPPLookupResult{}, nil
	}
	markers := make([]byte, len(components))
	for index := range components {
		markers[index] = components[index].Marker
	}
	result, err := engine.MatchProperNameTPPComponentPattern(surface, markers)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	result.CompositeSurface = append([]byte(nil), surface...)
	if result.Accepted {
		last := components[len(components)-1]
		result.ContextScanInputKnown = last.HasContextScanInput
		result.ContextChecksBypassed = last.HasContextScanInput && last.ContextScanInput == 0
		result.ContextChecksPending = !result.ContextChecksBypassed
	}
	return result, nil
}

// MatchProperNameTPPComponentArenaWithContext resolves the nonzero scanner
// branch after the exact lookup and marker gate. The supplied evidence carries
// outputs from FUN_1005a350 and its helper/table predicates; component count
// is always derived from the arena itself.
func (engine *Engine) MatchProperNameTPPComponentArenaWithContext(
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.MatchProperNameTPPComponentArena(componentArena)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	if !result.Accepted || !result.ContextChecksPending {
		return result, nil
	}
	components, err := text.Paul2013ProperNameTPPComponentsFromArena(componentArena)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	evidence.ComponentCount = len(components)
	result.ContextChecksEvaluated = true
	result.ContextGateAccepted = text.EvaluatePaul2013ProperNameTPPContextGate(evidence)
	result.ContextChecksPending = false
	return result, nil
}

// MatchProperNameTPPComponentArenaWithContextTables derives the two native
// context-set outcomes from supplied FUN_100035d0 table snapshots and scanner
// output, then evaluates the post-scanner gate.
func (engine *Engine) MatchProperNameTPPComponentArenaWithContextTables(
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
	scannerOutput []byte,
	contextSetA *text.Paul2013NativeStringSearchTable,
	contextSetB *text.Paul2013NativeStringSearchTable,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.MatchProperNameTPPComponentArena(componentArena)
	if err != nil || !result.Accepted || !result.ContextChecksPending {
		return result, err
	}
	components, err := text.Paul2013ProperNameTPPComponentsFromArena(componentArena)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	evidence.ComponentCount = len(components)
	result.ContextChecksEvaluated = true
	result.ContextGateAccepted, err = text.EvaluatePaul2013ProperNameTPPContextGateWithStringTables(
		evidence, scannerOutput, contextSetA, contextSetB,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	result.ContextChecksPending = false
	return result, nil
}

// MatchProperNameTPPComponentArenaWithScannerInputAndContextTables derives
// the scanner's leading-prefix counters and status 8/9 early returns from its
// NUL-terminated source. Ordinary output and status remain caller-supplied
// until the remaining FUN_1005a350 token scanner is ported.
func (engine *Engine) MatchProperNameTPPComponentArenaWithScannerInputAndContextTables(
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
	scannerInput []byte,
	ordinaryScannerOutput []byte,
	ordinaryScannerStatus int32,
	contextSetA *text.Paul2013NativeStringSearchTable,
	contextSetB *text.Paul2013NativeStringSearchTable,
) (Paul2013ProperNameTPPLookupResult, error) {
	observation, err := text.ObservePaul2013ProperNameScannerInput(
		scannerInput, ordinaryScannerOutput, ordinaryScannerStatus,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	evidence.ScannerResultCount = observation.ResultCount
	evidence.ScannerPrefixCount = observation.PrefixCount
	evidence.ScannerStatus = observation.Status
	return engine.MatchProperNameTPPComponentArenaWithContextTables(
		componentArena, evidence, observation.Output, contextSetA, contextSetB,
	)
}

// MatchProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables uses
// the directly supported mode-0x15 ASCII scanner paths. Unsupported
// scanner input still returns the independently known lookup result, with
// contextual checks pending; the boolean and reason report that boundary.
func (engine *Engine) MatchProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables(
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
	scannerInput []byte,
	contextSetA *text.Paul2013NativeStringSearchTable,
	contextSetB *text.Paul2013NativeStringSearchTable,
) (Paul2013ProperNameTPPLookupResult, bool, string, error) {
	observation := text.ObservePaul2013ProperNameASCIIToken(scannerInput)
	if !observation.Supported {
		result, err := engine.MatchProperNameTPPComponentArena(componentArena)
		return result, false, observation.Reason, err
	}
	evidence.ScannerResultCount = observation.ResultCount
	evidence.ScannerPrefixCount = observation.PrefixCount
	evidence.ScannerStatus = observation.Status
	result, err := engine.MatchProperNameTPPComponentArenaWithContextTables(
		componentArena, evidence, observation.Output, contextSetA, contextSetB,
	)
	return result, true, observation.Reason, err
}

// MatchProperNameTPPComponentArena projects the native counted component
// rows, then runs the supported surface, selector, exact-lookup, and marker
// gate. The nonzero context branch is evaluated by its WithContext variant.
func (engine *Engine) MatchProperNameTPPComponentArena(
	componentArena []byte,
) (Paul2013ProperNameTPPLookupResult, error) {
	components, err := text.Paul2013ProperNameTPPComponentsFromArena(componentArena)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return engine.MatchProperNameTPPComponents(components)
}

// MatchAndAppendProperNameTPPComponentArena composes the fully known fast
// path from FUN_10034180 with FUN_10034110's matched-component source-row
// follow-up. It appends when lookup and marker checks pass and the final-row
// +0x18 scanner-input dword is zero. Other candidates remain pending and do
// not mutate sourceArena; the WithContext variant can evaluate that branch
// from caller-supplied scanner and helper outcomes.
func (engine *Engine) MatchAndAppendProperNameTPPComponentArena(
	sourceArena []byte,
	componentArena []byte,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.MatchProperNameTPPComponentArena(componentArena)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return engine.appendMatchedProperNameComponentRows(sourceArena, componentArena, result)
}

// MatchAndAppendProperNameTPPComponentArenaWithContext also evaluates the
// nonzero scanner-input branch from caller-supplied native scanner and helper
// outcomes, then runs FUN_10034110 only if the complete modeled gate accepts.
func (engine *Engine) MatchAndAppendProperNameTPPComponentArenaWithContext(
	sourceArena []byte,
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.MatchProperNameTPPComponentArenaWithContext(componentArena, evidence)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return engine.appendMatchedProperNameComponentRows(sourceArena, componentArena, result)
}

// MatchAndAppendProperNameTPPComponentArenaWithContextTables also derives
// the context-set lookups from caller-supplied native table snapshots before
// applying the accepted source-row follow-up.
func (engine *Engine) MatchAndAppendProperNameTPPComponentArenaWithContextTables(
	sourceArena []byte,
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
	scannerOutput []byte,
	contextSetA *text.Paul2013NativeStringSearchTable,
	contextSetB *text.Paul2013NativeStringSearchTable,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.MatchProperNameTPPComponentArenaWithContextTables(
		componentArena, evidence, scannerOutput, contextSetA, contextSetB,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return engine.appendMatchedProperNameComponentRows(sourceArena, componentArena, result)
}

// MatchAndAppendProperNameTPPComponentArenaWithScannerInputAndContextTables
// composes the derived scanner prefix/terminal paths, supplied resource-table
// snapshots, context gate, and native source-row follow-up.
func (engine *Engine) MatchAndAppendProperNameTPPComponentArenaWithScannerInputAndContextTables(
	sourceArena []byte,
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
	scannerInput []byte,
	ordinaryScannerOutput []byte,
	ordinaryScannerStatus int32,
	contextSetA *text.Paul2013NativeStringSearchTable,
	contextSetB *text.Paul2013NativeStringSearchTable,
) (Paul2013ProperNameTPPLookupResult, error) {
	result, err := engine.MatchProperNameTPPComponentArenaWithScannerInputAndContextTables(
		componentArena, evidence, scannerInput, ordinaryScannerOutput,
		ordinaryScannerStatus, contextSetA, contextSetB,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return engine.appendMatchedProperNameComponentRows(sourceArena, componentArena, result)
}

// MatchAndAppendProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables
// composes the bounded mode-0x15 ASCII scanner, the proper-name lookup and
// context gate, and FUN_10034110's source-row append. Unsupported scanner
// input remains pending unless the component arena's native +0x18 zero value
// bypasses scanner-dependent context checks.
func (engine *Engine) MatchAndAppendProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables(
	sourceArena []byte,
	componentArena []byte,
	evidence text.Paul2013ProperNameTPPContextGateInput,
	scannerInput []byte,
	contextSetA *text.Paul2013NativeStringSearchTable,
	contextSetB *text.Paul2013NativeStringSearchTable,
) (Paul2013ProperNameTPPLookupResult, bool, string, error) {
	result, supported, reason, err := engine.MatchProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables(
		componentArena, evidence, scannerInput, contextSetA, contextSetB,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, supported, reason, err
	}
	result, err = engine.appendMatchedProperNameComponentRows(sourceArena, componentArena, result)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, supported, reason, err
	}
	return result, supported, reason, nil
}

func (engine *Engine) appendMatchedProperNameComponentRows(
	sourceArena []byte,
	componentArena []byte,
	result Paul2013ProperNameTPPLookupResult,
) (Paul2013ProperNameTPPLookupResult, error) {
	contextAccepted := result.ContextChecksBypassed ||
		(result.ContextChecksEvaluated && result.ContextGateAccepted)
	if !result.Accepted || !contextAccepted {
		return result, nil
	}
	result.FollowupAttempted = true
	var err error
	result.SourceArena, result.FollowupSucceeded, err = text.AppendPaul2013MatchedComponentSourceRows(
		sourceArena, componentArena,
	)
	if err != nil {
		return Paul2013ProperNameTPPLookupResult{}, err
	}
	return result, nil
}
