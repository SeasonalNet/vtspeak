package text

import "fmt"

// Paul2013ProperNameTPPContextGateInput contains scanner fields and raw
// predicate outcomes consumed by the nonzero-input branch of FUN_10034180.
// The numeric scanner fields preserve local_3c, local_4c[0], and local_30.
// Boolean outcomes correspond to their named native lookup/helper calls; they
// do not assign linguistic meaning to those tables or helper results.
type Paul2013ProperNameTPPContextGateInput struct {
	ComponentCount int

	ScannerResultCount             int32
	ScannerPrefixCount             int32
	ScannerStatus                  int32
	ContextLookupAResultNotHighBit bool
	ContextLookupBResultNotHighBit bool

	SingleComponentHeaderAccepted bool
	TwoComponentFirstIsGeorge     bool
	TwoComponentLastMarkerIsD     bool
	TwoComponentKeyFound          bool
	ScannerWordPredicateAccepted  bool

	SingleComponentLateLookupBelowHighBit bool
	SingleComponentGeorgeAliasFound       bool
	SingleComponentMartinOrLuther         bool
}

// EvaluatePaul2013ProperNameTPPContextGate ports FUN_10034180's
// post-marker-gate Boolean decision tree for a nonzero final-row scanner
// input. Scanner production, pointer-backed surfaces, resource-table queries,
// and helper predicates are supplied as raw observed outcomes.
func EvaluatePaul2013ProperNameTPPContextGate(
	input Paul2013ProperNameTPPContextGateInput,
) bool {
	scannerPathNeedsExceptions := input.ScannerResultCount < 1 ||
		input.ScannerPrefixCount > 1 || input.ScannerStatus != 1 ||
		(!input.ContextLookupAResultNotHighBit && !input.ContextLookupBResultNotHighBit)
	if !scannerPathNeedsExceptions {
		return false
	}
	if input.ComponentCount == 1 && !input.SingleComponentHeaderAccepted {
		return false
	}

	// This conjunction is the only two-component path through the compound
	// rejection condition that reaches LAB_10034470 (low-word success).
	if input.ComponentCount == 2 && input.TwoComponentFirstIsGeorge &&
		input.TwoComponentLastMarkerIsD && input.TwoComponentKeyFound &&
		input.ScannerResultCount > 0 && input.ScannerPrefixCount <= 1 &&
		input.ScannerStatus == 1 && input.ScannerWordPredicateAccepted {
		return true
	}

	if input.ComponentCount == 1 && input.ScannerResultCount > 0 &&
		input.ScannerPrefixCount < 2 && input.ScannerStatus == 1 {
		if input.SingleComponentLateLookupBelowHighBit {
			return false
		}
		if (input.SingleComponentGeorgeAliasFound || input.SingleComponentMartinOrLuther) &&
			input.ScannerWordPredicateAccepted {
			return false
		}
	}
	return true
}

// EvaluatePaul2013ProperNameTPPContextGateWithStringTables derives the two
// FUN_100035d0 membership outcomes from caller-supplied resource snapshots
// and the scanner output, then evaluates the same native decision tree. A nil
// table represents the native nil/empty-table miss path.
func EvaluatePaul2013ProperNameTPPContextGateWithStringTables(
	input Paul2013ProperNameTPPContextGateInput,
	scannerOutput []byte,
	contextSetA *Paul2013NativeStringSearchTable,
	contextSetB *Paul2013NativeStringSearchTable,
) (bool, error) {
	input.SingleComponentLateLookupBelowHighBit = Paul2013ProperNameContextLateLookup(scannerOutput)
	input.ScannerWordPredicateAccepted = Paul2013ProperNameScannerWordPredicate(scannerOutput)
	var err error
	input.ContextLookupAResultNotHighBit, err = nativeContextStringTableContains(contextSetA, scannerOutput)
	if err != nil {
		return false, fmt.Errorf("search proper-name context set A: %w", err)
	}
	input.ContextLookupBResultNotHighBit, err = nativeContextStringTableContains(contextSetB, scannerOutput)
	if err != nil {
		return false, fmt.Errorf("search proper-name context set B: %w", err)
	}
	return EvaluatePaul2013ProperNameTPPContextGate(input), nil
}

func nativeContextStringTableContains(
	table *Paul2013NativeStringSearchTable,
	key []byte,
) (bool, error) {
	if table == nil {
		return false, nil
	}
	_, found, err := table.Lookup(key)
	return found, err
}
