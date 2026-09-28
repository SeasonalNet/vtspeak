package selection

import (
	"errors"

	"vtspeak/engine/distance"
)

// FUN_100182e0 uses float32 divisors 2, 16, 8, and 5 from
// binary/vt_pau.dll RVAs 0x6d174, 0x6d17c, 0x6d180, and 0x6d184.

// FeaturePair identifies two observed 8-bit feature codes whose cost is read
// from the generated Paul feature-distance table.
type FeaturePair struct {
	Candidate byte
	Context   byte
}

// UnitBytePenaltyInput holds the signatures and state markers needed by the
// direct byte-1 through byte-3 comparisons in FUN_100182e0.
type UnitBytePenaltyInput struct {
	CandidateSignature [7]byte
	ContextSignature   [7]byte
	ContextMode        byte
	ByteThreeStateFive bool
	CandidateMarker    byte
	ContextMarker      byte
}

// UnitRecordBytePenalty calculates the directly observed byte-1 through
// byte-3 mismatch costs using the DLL category maps. The caller supplies the
// state-derived marker bytes; table-coded byte-5 and cross-byte costs are
// separate and are not included here.
func UnitRecordBytePenalty(input UnitBytePenaltyInput) float32 {
	candidate := input.CandidateSignature
	context := input.ContextSignature
	var penalty float32

	if candidate[2] != context[2] {
		if candidate[2] == 0x4b || context[2] == 0x4b {
			penalty += 100
		} else if unitPrimaryCategory[candidate[2]] == unitPrimaryCategory[context[2]] {
			penalty += 5
		} else {
			penalty += 1000
		}
	}

	if candidate[1] != context[1] && input.ContextMode != 2 {
		if candidate[5]&0x80 != 0 && candidate[5]&0x38 > 0x17 {
			penalty++
		} else if unitPrimaryCategory[candidate[1]] == unitPrimaryCategory[context[1]] {
			penalty += 5
		} else if unitLookupClass[candidate[1]] == unitLookupClass[context[1]] {
			penalty += 10
		} else {
			penalty += 500
		}
	}

	if candidate[3] != context[3] && input.ContextMode != 1 {
		switch {
		case input.ByteThreeStateFive && context[3] == 0x5a:
			penalty += 100
		case candidate[5]&0x40 != 0 && candidate[5]&7 > 2:
			penalty++
		case candidate[3] == 0x4b:
			penalty += 100
		case unitPrimaryCategory[candidate[3]] == unitPrimaryCategory[context[3]]:
			if context[3] != 0x5b {
				penalty += 5
			}
		case unitLookupClass[candidate[3]] == unitLookupClass[context[3]]:
			penalty += 10
		default:
			penalty += 500
		}
	}

	if input.CandidateMarker == 'd' {
		penalty += 10
	}
	if input.CandidateMarker == 0x0c && input.ContextMarker != 0x0c {
		penalty += 1000
	}
	return penalty
}

// UnitByteFivePenalty applies the two recovered field-cost tables to byte 5
// of the candidate and context signatures. The high-field table has five
// rows and columns; other observed field values fail closed instead of
// reading adjacent DLL data.
func UnitByteFivePenalty(candidateSignature, contextSignature [7]byte) (float32, error) {
	candidateHigh := candidateSignature[5] >> 3 & 7
	contextHigh := contextSignature[5] >> 3 & 7
	if candidateHigh >= 5 || contextHigh >= 5 {
		return 0, errors.New("unit byte-5 high field exceeds the recovered 5-by-5 table")
	}
	candidateLow := candidateSignature[5] & 7
	contextLow := contextSignature[5] & 7
	highCost := unitHighFieldCosts[contextHigh*5+candidateHigh]
	lowCost := unitLowFieldCosts[contextLow*8+candidateLow]
	return float32(highCost + lowCost), nil
}

// UnitCrossBytePenalty applies the recovered 3-by-3 category-pair tables. It
// rejects mapped values outside the table rather than reproducing the DLL's
// unchecked adjacent-memory reads for sentinel categories.
func UnitCrossBytePenalty(candidateSignature, contextSignature [7]byte) (float32, error) {
	candidate := candidateSignature
	context := contextSignature
	if unitPairClass[context[2]] != unitPairClass[candidate[2]] {
		return unitPairCost(unitPairPenaltyPrimary, context[2], candidate[2])
	}

	if unitContextFlags[context[1]] != 0 && (context[2] == 0x36 || context[2] == 0x2b) {
		penalty, err := unitPairCost(unitPairPenaltyPrimary, context[1], candidate[1])
		if err != nil {
			return 0, err
		}
		if context[1] != candidate[1] {
			penalty += 100
		}
		return penalty, nil
	}

	if candidate[5]&7 != 0 || unitContextFlags[context[3]] == 0 {
		return 0, nil
	}
	return unitPairCost(unitPairPenaltyAlternate, context[3], candidate[3])
}

func unitPairCost(table [9]int32, contextCode, candidateCode byte) (float32, error) {
	contextClass := unitPairClass[contextCode]
	candidateClass := unitPairClass[candidateCode]
	if contextClass >= 3 || candidateClass >= 3 {
		return 0, errors.New("unit cross-byte category exceeds the recovered 3-by-3 table")
	}
	return float32(table[contextClass*3+candidateClass]), nil
}

// KnownUnitCategoricalPenalty combines the currently recovered per-unit
// byte-1 through byte-3, byte-5, and in-range cross-byte costs.
func KnownUnitCategoricalPenalty(input UnitBytePenaltyInput) (float32, error) {
	byteFive, err := UnitByteFivePenalty(input.CandidateSignature, input.ContextSignature)
	if err != nil {
		return 0, err
	}
	crossByte, err := UnitCrossBytePenalty(input.CandidateSignature, input.ContextSignature)
	if err != nil {
		return 0, err
	}
	return UnitRecordBytePenalty(input) + byteFive + crossByte, nil
}

// UnitScoreInput contains the scalar decisions and feature codes consumed by
// FUN_100182e0 after its caller has resolved the context-specific arrays.
// The categorical penalty is caller supplied because its predicates use
// additional context maps and unit-record tables.
type UnitScoreInput struct {
	PrimaryFeature          FeaturePair
	AdditionalFeatures      [2]FeaturePair
	BytePenalty             UnitBytePenaltyInput
	ThreeFeatureMode        bool
	ScalePrimaryByEight     bool
	EqualCategories         bool
	BothSpecialFeatureClass bool
	AdditionalCategoryCost  float32
	Scale                   float32
}

// ScoreUnitCost calculates the observed per-unit feature term and combines
// it with recovered categorical penalties and any additional context cost.
// It ports the arithmetic in FUN_100182e0 while leaving context-array and
// feature-bin selection to the caller.
func ScoreUnitCost(table *distance.FeatureTable, input UnitScoreInput) (float32, error) {
	if table == nil {
		return 0, errors.New("unit feature-distance table is not loaded")
	}
	categoryCost, err := KnownUnitCategoricalPenalty(input.BytePenalty)
	if err != nil {
		return 0, err
	}
	categoryCost += input.AdditionalCategoryCost
	primary, err := table.Lookup(input.PrimaryFeature.Candidate, input.PrimaryFeature.Context)
	if err != nil {
		return 0, err
	}

	var featureCost float32
	if input.ThreeFeatureMode {
		first, err := table.Lookup(input.AdditionalFeatures[0].Candidate, input.AdditionalFeatures[0].Context)
		if err != nil {
			return 0, err
		}
		second, err := table.Lookup(input.AdditionalFeatures[1].Candidate, input.AdditionalFeatures[1].Context)
		if err != nil {
			return 0, err
		}
		primaryDivisor := float32(1)
		additionalDivisor := float32(2)
		if input.ScalePrimaryByEight {
			primaryDivisor = 8
			additionalDivisor = 8
		}
		featureCost = primary/primaryDivisor + (first+second)/additionalDivisor
	} else {
		featureCost = primary / 5
		if input.EqualCategories {
			categoryCost /= 2
		}
	}
	if input.BothSpecialFeatureClass {
		featureCost /= 16
	}
	return featureCost*input.Scale + categoryCost, nil
}
