package selection

import (
	"errors"

	"vtspeak/engine/distance"
)

// ModeTwoCategoricalPenalty returns the directly observed categorical cost
// for the mode-2 branch of FUN_10018c80. It leaves the signature byte's
// linguistic meaning unspecified.
func ModeTwoCategoricalPenalty(current, previous UnitRef, currentSignature, previousSignature [7]byte) float32 {
	penalty := float32(5)
	if current == previous {
		penalty = 0
	}
	if currentSignature[2] != previousSignature[2] {
		penalty += 1000
	}
	return penalty
}

// ModeOtherCategoricalPenalty returns the observed signature-category cost
// for transition modes other than 2. The class and soft maps are kept as
// byte-code sets because the DLL only tests their values for zero/nonzero.
// They are the truthy entries at RVAs 0x7c384 and 0x7c2c4, respectively.
// Context-derived +10 and +2000 rules are added separately by callers that
// have the corresponding model-state values.
func ModeOtherCategoricalPenalty(currentSignature, previousSignature [7]byte) float32 {
	penalty := float32(0)

	currentHighBoundary := currentSignature[5]>>7 != 0
	previousLowBoundary := previousSignature[5]>>6&1 != 0
	if currentHighBoundary != previousLowBoundary {
		penalty += 1000
	} else if currentHighBoundary {
		currentClass := currentSignature[5] >> 3 & 7
		previousClass := previousSignature[5] & 7
		penalty += float32(transitionBoundaryCosts[previousClass*8+currentClass])
	}

	if currentSignature[2] != previousSignature[1] {
		if transitionClassByte(currentSignature[2]) || transitionClassByte(previousSignature[1]) {
			penalty += 1000
		} else if transitionSoftByte(currentSignature[2]) {
			penalty += 100
		}
	}

	if currentSignature[3] != previousSignature[2] {
		if transitionClassByte(currentSignature[3]) || transitionClassByte(previousSignature[2]) {
			penalty += 1000
		} else if transitionSoftByte(previousSignature[2]) {
			penalty += 100
		}
	}

	return penalty
}

// TransitionContextPenalty adds the two mode-independent penalties whose
// predicates are derived from model context arrays in FUN_10018c80. The
// caller supplies those predicates because they are not unit-signature data.
func TransitionContextPenalty(addTenPointRule, currentSpecial, previousSpecial bool) float32 {
	var penalty float32
	if addTenPointRule {
		penalty += 10
	}
	if currentSpecial != previousSpecial {
		penalty += 2000
	}
	return penalty
}

// TransitionDistanceCost combines the transition's fixed weighted-distance
// terms and mode-specific divisor. Distances and durationTerm are float32
// values already obtained from the model tables/state; durationTerm is the
// current context's scaled prosody contribution. weightRow is the local row
// selected by the DLL's context-state branch and must be in [0,4].
func TransitionDistanceCost(modeTwo bool, weightRow uint8, rawDistance, featureDistanceOne, featureDistanceTwo, durationTerm float32, currentMetricCode, previousMetricCode uint16, currentSignature, previousSignature [7]byte) (float32, bool) {
	if weightRow >= uint8(len(transitionWeights)) {
		return 0, false
	}
	weights := transitionWeights[weightRow]
	weighted := rawDistance*weights.raw + featureDistanceOne*weights.featureOne + featureDistanceTwo*weights.featureTwo + 2
	divisor := transitionDivisor(modeTwo, currentSignature, previousSignature)
	category := transitionPackedCategoryCosts[(previousMetricCode>>14&3)*3+(currentMetricCode>>14&3)]
	return weighted/divisor + float32(category) + durationTerm, true
}

// TransitionScoreInput contains current/previous candidate values plus the
// feature and context terms produced elsewhere in the synthesis state.
// WeightRow follows the local row selected by the DLL's context-state branch.
type TransitionScoreInput struct {
	Current, Previous                      UnitRef
	CurrentSignature, PreviousSignature    [7]byte
	CurrentMetricCode, PreviousMetricCode  uint16
	CurrentFeatureOne, PreviousFeatureOne  byte
	CurrentFeatureTwo, PreviousFeatureTwo  byte
	WeightRow                              uint8
	ModeTwo                                bool
	DurationTerm, PreviousCumulative       float32
	AddTenPointContextRule                 bool
	CurrentSpecialContext, PreviousSpecial bool
}

// ScoreTransition combines the recovered transition scorer when supplied the
// feature-bin codes and context-state predicates. It looks up the raw masked
// metric distance and both generated feature distances, then adds category,
// cumulative, and scaled duration/prosody costs.
func ScoreTransition(table *distance.Table, featureTable *distance.FeatureTable, input TransitionScoreInput) (float32, error) {
	if table == nil {
		return 0, errors.New("transition distance table is not loaded")
	}
	if featureTable == nil {
		return 0, errors.New("transition feature-distance table is not loaded")
	}
	rawDistance, err := table.LookupPacked(input.CurrentMetricCode, input.PreviousMetricCode)
	if err != nil {
		return 0, err
	}
	featureDistanceOne, err := featureTable.Lookup(input.CurrentFeatureOne, input.PreviousFeatureOne)
	if err != nil {
		return 0, err
	}
	featureDistanceTwo, err := featureTable.Lookup(input.CurrentFeatureTwo, input.PreviousFeatureTwo)
	if err != nil {
		return 0, err
	}
	distanceCost, ok := TransitionDistanceCost(
		input.ModeTwo,
		input.WeightRow,
		rawDistance,
		featureDistanceOne,
		featureDistanceTwo,
		input.DurationTerm,
		input.CurrentMetricCode,
		input.PreviousMetricCode,
		input.CurrentSignature,
		input.PreviousSignature,
	)
	if !ok {
		return 0, errors.New("transition weight row is outside the recovered table")
	}
	if input.ModeTwo {
		distanceCost += ModeTwoCategoricalPenalty(input.Current, input.Previous, input.CurrentSignature, input.PreviousSignature)
	} else {
		distanceCost += ModeOtherCategoricalPenalty(input.CurrentSignature, input.PreviousSignature)
	}
	distanceCost += TransitionContextPenalty(input.AddTenPointContextRule, input.CurrentSpecialContext, input.PreviousSpecial)
	return distanceCost + input.PreviousCumulative, nil
}

// transitionDivisor ports the byte-field rule from FUN_10018c80.
func transitionDivisor(modeTwo bool, currentSignature, previousSignature [7]byte) float32 {
	currentHighBoundary := currentSignature[5]>>7 != 0
	previousLowBoundary := previousSignature[5]>>6&1 != 0
	if modeTwo || !currentHighBoundary || !previousLowBoundary {
		return 1
	}
	class := currentSignature[5] >> 3 & 7
	switch {
	case class < 2:
		return 1
	case class == 2:
		return 2
	default:
		return 8
	}
}

type transitionWeightsRow struct {
	raw, featureOne, featureTwo float32
}

// Tables below are little-endian constants read from binary/vt_pau.dll at
// the corresponding RVAs. Their source fields are intentionally unnamed.
var transitionBoundaryCosts = [64]int32{
	0, 20, 500, 1000, 1000, 1000, 5000, 1000,
	20, 0, 500, 1000, 1000, 1000, 5000, 1000,
	500, 500, 0, 50, 500, 1000, 5000, 1000,
	1000, 1000, 50, 0, 10, 100, 5000, 1000,
	1000, 1000, 500, 10, 0, 100, 5000, 1000,
	1000, 1000, 1000, 500, 100, 0, 5000, 1000,
	5000, 5000, 5000, 5000, 5000, 5000, 0, 1000,
	1000, 1000, 1000, 1000, 1000, 100, 1000, 0,
}

var transitionPackedCategoryCosts = [9]int32{0, 100, 0, 100, 0, 0, 0, 0, 0}

// FUN_10018c80 reads three float32 coefficients at row strides of 12 bytes.
// The address order is raw, featureTwo, featureOne; the struct labels follow
// the pseudocode's actual multiplications.
var transitionWeights = [5]transitionWeightsRow{
	{raw: 10, featureOne: 5, featureTwo: 10},
	{raw: 10, featureOne: 2, featureTwo: 5},
	{raw: 10, featureOne: 2, featureTwo: 5},
	{raw: 10, featureOne: 2, featureTwo: 2},
	{raw: 10, featureOne: 2, featureTwo: 5},
}

func transitionClassByte(value byte) bool {
	switch value {
	case 73, 74, 139, 150:
		return true
	default:
		return false
	}
}

func transitionSoftByte(value byte) bool {
	switch value {
	case 21, 42, 53, 57, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
		99, 102, 105, 108, 111, 114, 121, 124, 127, 133, 136, 145, 148,
		157, 160:
		return true
	default:
		return false
	}
}
