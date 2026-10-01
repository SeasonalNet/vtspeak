package selection

import (
	"fmt"

	"vtspeak/engine/internal/paul2013tables"
)

// These counted rule lists and pointer-table layouts are the direct byte
// contents referenced by FUN_10018770 at 0x1007c4dc and 0x1007c504. Rule
// values remain opaque and are kept in native iteration order.
var paul2013LeftQueryRules = [...][5]byte{
	{2, 0, 1},
	{3, 1, 0, 2},
	{3, 2, 1, 3},
	{4, 3, 4, 2, 1},
	{4, 4, 3, 2, 1},
	{1, 0},
	{3, 1, 0, 2},
	{3, 2, 1, 3},
	{2, 3, 4},
	{2, 4, 3},
	{1, 0},
	{3, 1, 0, 2},
	{3, 2, 1, 3},
}

var paul2013RightQueryRules = [...][5]byte{
	{1, 0},
	{3, 1, 0, 2},
	{3, 2, 1, 3},
	{3, 3, 4, 5},
	{3, 4, 3, 5},
	{1, 5},
	{1, 6},
	{2, 7, 5},
	{1, 0},
	{3, 1, 0, 2},
	{3, 2, 1, 3},
	{3, 3, 4, 5},
	{3, 4, 3, 5},
	{1, 5},
	{1, 6},
	{2, 7, 5},
}

// Paul2013QueryRulePair is one eligible pair of opaque left/right rule
// values from FUN_10018770. PackedCode is the byte expression written to the
// query signature's rule-code position after the two state flags are applied.
type Paul2013QueryRulePair struct {
	LeftOrdinal    int
	RightOrdinal   int
	LeftRule       byte
	RightRule      byte
	PackedCode     byte
	LeftGateMatch  bool
	RightGateMatch bool
}

// Paul2013QueryAttempt is one primary FUN_10023f90 input emitted by
// FUN_10018770. Signature is the exact seven-byte value passed to the lookup;
// Mode is the caller's mode for the untouched first pair and zero for packed
// rule variants.
type Paul2013QueryAttempt struct {
	Signature [7]byte
	Mode      int16
	RulePair  Paul2013QueryRulePair
}

var paul2013MappedQueryRule = [8]byte{0, 0, 0, 0x61, 0x5b, 0x5a, 0x01, 0x00}

// Paul2013QueryRulePairs selects the native left and right rule lists from
// signature byte +5, then applies the count-remainder and context-row gates
// before returning eligible pairs in native nested-loop order. priorMetric is
// the value divided into quotient/remainder by ten in FUN_10018770. The
// function only ports this rule-selection stage; accumulated candidate state
// and the producers of its inputs stay outside this helper.
func Paul2013QueryRulePairs(
	signature [7]byte,
	contextRow [6]byte,
	priorMetric uint32,
) ([]Paul2013QueryRulePair, error) {
	selector := signature[5]
	leftIndex := int(selector>>7)*5 + int((selector>>3)&7)
	if leftIndex >= len(paul2013LeftQueryRules) {
		return nil, fmt.Errorf("Paul 2013 left query rule selector %d is outside the captured table", leftIndex)
	}
	rightIndex := int(selector&7) + int((selector>>6)&1)*8
	leftRules := paul2013LeftQueryRules[leftIndex]
	rightRules := paul2013RightQueryRules[rightIndex]
	leftCount, rightCount := int(leftRules[0]), int(rightRules[0])
	if leftCount == 0 || rightCount == 0 {
		return nil, nil
	}

	leftDefault := (selector>>3)&7 <= 1
	rightDefault := selector&7 <= 1
	leftFlag := selector >> 7
	rightFlag := (selector >> 6) & 1
	quotient, remainder := priorMetric/10, priorMetric%10
	leftStateGate := quotient == 0 || contextRow[4] == 2
	rightStateGate := remainder == 0 || contextRow[4] == 1
	pairs := make([]Paul2013QueryRulePair, 0, leftCount*rightCount)
	for leftOrdinal, leftRule := range leftRules[1 : leftCount+1] {
		leftRuleClass := leftRule <= 1
		if !leftStateGate && leftRuleClass != leftDefault {
			continue
		}
		packedLeftFlag := byte(0)
		if leftRule != 0 {
			packedLeftFlag = leftFlag
		}
		for rightOrdinal, rightRule := range rightRules[1 : rightCount+1] {
			rightRuleClass := rightRule <= 1
			if !rightStateGate && rightRuleClass != rightDefault {
				continue
			}
			packedRightFlag := byte(0)
			if rightRule != 0 {
				packedRightFlag = rightFlag
			}
			packed := (((packedLeftFlag*2+packedRightFlag)*8+leftRule)*8 + rightRule)
			pairs = append(pairs, Paul2013QueryRulePair{
				LeftOrdinal: leftOrdinal, RightOrdinal: rightOrdinal,
				LeftRule: leftRule, RightRule: rightRule, PackedCode: packed,
				LeftGateMatch:  leftRuleClass == leftDefault,
				RightGateMatch: rightRuleClass == rightDefault,
			})
		}
	}
	return pairs, nil
}

// Paul2013QueryFollowup is one result-dependent category lookup from
// FUN_10018770. All such lookups use mode zero. Kind preserves the native
// branch/order for trace comparison; signature bytes not written by that
// branch retain the values in its source signature.
type Paul2013QueryFollowup struct {
	Signature [7]byte
	Kind      string
}

// Paul2013QueryCategoryFollowups emits the category-derived lookups after a
// primary lookup. sourceSignature is the original function input and
// querySignature is the possibly rule-mutated primary signature. This
// distinction preserves the native third-branch restore to the original
// signature. primaryResult is the signed 16-bit result consumed by the native
// gate. contextGate is the byte read from the per-position state table at
// +0x8ea; its producer is deliberately caller-supplied.
func Paul2013QueryCategoryFollowups(
	sourceSignature [7]byte,
	querySignature [7]byte,
	primaryResult int16,
	contextGate byte,
) []Paul2013QueryFollowup {
	selector := sourceSignature[5]
	if primaryResult == 0 && selector&0xc0 == 0 {
		return nil
	}
	mapValue := paul2013SignatureQueryMap[sourceSignature[2]]
	if mapValue == 0 {
		left := paul2013SignatureQueryMap[sourceSignature[1]]
		right := paul2013SignatureQueryMap[sourceSignature[3]]
		switch {
		case left == 0 && right == 0:
			return nil
		case left == 0:
			query := querySignature
			query[3] = right
			return []Paul2013QueryFollowup{{Signature: query, Kind: "right-category"}}
		case right == 0:
			query := querySignature
			query[1] = left
			return []Paul2013QueryFollowup{{Signature: query, Kind: "left-category"}}
		default:
			leftQuery := querySignature
			leftQuery[1] = left
			rightAfterLeft := leftQuery
			rightAfterLeft[3] = right
			rightFromSource := sourceSignature
			rightFromSource[3] = right
			return []Paul2013QueryFollowup{
				{Signature: leftQuery, Kind: "left-category"},
				{Signature: rightAfterLeft, Kind: "right-category-after-left"},
				{Signature: rightFromSource, Kind: "right-category-from-source"},
			}
		}
	}

	boundary := int8(paul2013tables.ByteTable1007BB70(sourceSignature[3]))
	if (selector&0x40 != 0 || sourceSignature[2] != 0x39 || boundary < 1) && contextGate == 0 {
		query := querySignature
		query[2] = mapValue
		return []Paul2013QueryFollowup{{Signature: query, Kind: "middle-category"}}
	}
	return nil
}

// Paul2013QueryRecursiveSignature returns the optional recursive mode-zero
// target from FUN_10018770. recursionGate corresponds to local_28, which is
// set only after a positive accumulated lookup result on a pair whose two
// default-category predicates both matched. Mode two bypasses that gate.
func Paul2013QueryRecursiveSignature(
	signature [7]byte,
	mode int16,
	recursionGate bool,
) ([7]byte, bool, error) {
	if mode < 0 || mode > 2 {
		return [7]byte{}, false, fmt.Errorf("Paul 2013 query mode %d is outside the observed range 0 through 2", mode)
	}
	if mode == 0 || (mode != 2 && !recursionGate) {
		return [7]byte{}, false, nil
	}
	value := paul2013SignatureQueryMap[0x60+int(signature[2])]
	if value == 0 {
		return [7]byte{}, false, nil
	}
	query := signature
	query[2] = value
	return query, true, nil
}

// Paul2013QueryLookup executes one native FUN_10023f90 call. Implementations
// can close over the model-state and candidate buffers that lookup mutates.
type Paul2013QueryLookup func(signature [7]byte, mode int16) (int16, error)

// Paul2013QueryExecutionAttempt records one lookup call in native order.
type Paul2013QueryExecutionAttempt struct {
	Signature [7]byte
	Mode      int16
	RulePair  Paul2013QueryRulePair
	Kind      string
	Result    int16
}

// Paul2013QuerySequenceResult describes the deterministic control flow around
// FUN_10023f90. Accumulated is the last rule pair's signed-short subtotal;
// RecursionGate is the function-wide local_28 predicate. The callback owns
// the native lookup's state mutations, which this function cannot reproduce
// without the unresolved state and candidate buffers.
type Paul2013QuerySequenceResult struct {
	Attempts      []Paul2013QueryExecutionAttempt
	Accumulated   int16
	RecursionGate bool
	Accepted      bool
	Recurred      bool
}

// RunPaul2013QuerySequence executes the source-supported query orchestration
// for FUN_10018770, including primary lookups, result-gated category lookups,
// signed-short accumulation, return gating, and its single mode-zero recursive
// pass. contextGate reads the live state byte at +0x8ea after each primary
// lookup. The rule-selection row and prior metric are the values observed on
// entry; their producers remain caller-owned.
func RunPaul2013QuerySequence(
	signature [7]byte,
	contextRow [6]byte,
	priorMetric uint32,
	mode int16,
	contextGate func() byte,
	lookup Paul2013QueryLookup,
) (Paul2013QuerySequenceResult, error) {
	if mode < 0 || mode > 2 {
		return Paul2013QuerySequenceResult{}, fmt.Errorf("Paul 2013 query mode %d is outside the observed range 0 through 2", mode)
	}
	if contextGate == nil {
		return Paul2013QuerySequenceResult{}, fmt.Errorf("Paul 2013 query sequence has no live context-state gate")
	}
	if lookup == nil {
		return Paul2013QuerySequenceResult{}, fmt.Errorf("Paul 2013 query sequence has no class-lookup callback")
	}
	return runPaul2013QuerySequence(signature, contextRow, priorMetric, mode, contextGate, lookup, true)
}

func runPaul2013QuerySequence(
	signature [7]byte,
	contextRow [6]byte,
	priorMetric uint32,
	mode int16,
	contextGate func() byte,
	lookup Paul2013QueryLookup,
	allowRecursion bool,
) (Paul2013QuerySequenceResult, error) {
	baseAttempts, err := Paul2013QueryBaseAttempts(signature, contextRow, priorMetric, mode)
	if err != nil {
		return Paul2013QuerySequenceResult{}, err
	}
	result := Paul2013QuerySequenceResult{Attempts: make([]Paul2013QueryExecutionAttempt, 0, len(baseAttempts))}
	for _, base := range baseAttempts {
		primary, err := lookup(base.Signature, base.Mode)
		if err != nil {
			return result, fmt.Errorf("primary query for rule pair (%d,%d): %w", base.RulePair.LeftOrdinal, base.RulePair.RightOrdinal, err)
		}
		attempt := Paul2013QueryExecutionAttempt{
			Signature: base.Signature, Mode: base.Mode,
			RulePair: base.RulePair, Kind: "primary", Result: primary,
		}
		result.Attempts = append(result.Attempts, attempt)
		accumulated := primary
		stateGate := byte(0)
		if primary != 0 || signature[5]&0xc0 != 0 {
			stateGate = contextGate()
		}
		followups := Paul2013QueryCategoryFollowups(signature, base.Signature, primary, stateGate)
		for _, followup := range followups {
			followupResult, err := lookup(followup.Signature, 0)
			if err != nil {
				return result, fmt.Errorf("%s query for rule pair (%d,%d): %w", followup.Kind, base.RulePair.LeftOrdinal, base.RulePair.RightOrdinal, err)
			}
			result.Attempts = append(result.Attempts, Paul2013QueryExecutionAttempt{
				Signature: followup.Signature, RulePair: base.RulePair,
				Kind: followup.Kind, Result: followupResult,
			})
			accumulated = int16(uint16(accumulated) + uint16(followupResult))
		}
		result.Accumulated = accumulated
		if accumulated > 0 && base.RulePair.LeftGateMatch && base.RulePair.RightGateMatch {
			result.RecursionGate = true
		}
	}
	if mode != 2 && !result.RecursionGate {
		return result, nil
	}
	result.Accepted = true
	if !allowRecursion || mode == 0 {
		return result, nil
	}
	recursiveSignature, recurse, err := Paul2013QueryRecursiveSignature(signature, mode, result.RecursionGate)
	if err != nil {
		return result, err
	}
	if !recurse {
		return result, nil
	}
	result.Recurred = true
	recursiveResult, err := runPaul2013QuerySequence(
		recursiveSignature, contextRow, priorMetric, 0, contextGate, lookup, false,
	)
	if err != nil {
		return result, fmt.Errorf("recursive mode-zero query: %w", err)
	}
	result.Attempts = append(result.Attempts, recursiveResult.Attempts...)
	return result, nil
}

// Paul2013QueryBaseAttempts adds the signature-byte mutations performed
// before the first lookup for each eligible rule pair in FUN_10018770. The
// untouched first pair carries lookupMode; packed pairs use mode zero.
// Result-dependent category lookups are emitted separately by
// Paul2013QueryCategoryFollowups.
func Paul2013QueryBaseAttempts(
	signature [7]byte,
	contextRow [6]byte,
	priorMetric uint32,
	lookupMode int16,
) ([]Paul2013QueryAttempt, error) {
	if lookupMode < 0 || lookupMode > 2 {
		return nil, fmt.Errorf("Paul 2013 query mode %d is outside the observed range 0 through 2", lookupMode)
	}
	pairs, err := Paul2013QueryRulePairs(signature, contextRow, priorMetric)
	if err != nil {
		return nil, err
	}
	attempts := make([]Paul2013QueryAttempt, 0, len(pairs))
	for _, pair := range pairs {
		query := signature
		mode := int16(0)
		if pair.LeftOrdinal == 0 && pair.RightOrdinal == 0 {
			mode = lookupMode
		} else {
			selector := signature[5]
			if pair.LeftOrdinal > 0 && selector&0x80 != 0 && pair.LeftRule > 2 {
				mapped := paul2013MappedQueryRule[pair.LeftRule]
				query[0], query[1] = mapped, mapped
			}
			if pair.RightOrdinal > 0 && selector&0x40 != 0 && pair.RightRule > 2 {
				mapped := paul2013MappedQueryRule[pair.RightRule]
				query[3], query[4] = mapped, mapped
			}
			query[5] = pair.PackedCode
		}
		attempts = append(attempts, Paul2013QueryAttempt{
			Signature: query, Mode: mode, RulePair: pair,
		})
	}
	return attempts, nil
}
