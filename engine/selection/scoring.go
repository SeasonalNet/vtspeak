package selection

import (
	"errors"

	"vtspeak/engine/dat"
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
	return transitionDistanceCost(
		modeTwo, weightRow, rawDistance, featureDistanceOne, featureDistanceTwo,
		durationTerm, currentMetricCode, previousMetricCode, currentSignature,
		previousSignature, false,
	)
}

func transitionDistanceCost(modeTwo bool, weightRow uint8, rawDistance, featureDistanceOne, featureDistanceTwo, durationTerm float32, currentMetricCode, previousMetricCode uint16, currentSignature, previousSignature [7]byte, suppressWeightedTable bool) (float32, bool) {
	if weightRow >= uint8(len(transitionWeights)) {
		return 0, false
	}
	weights := transitionWeights[weightRow]
	weighted := float32(0)
	if !suppressWeightedTable {
		weighted = rawDistance*weights.raw + featureDistanceOne*weights.featureOne + featureDistanceTwo*weights.featureTwo + 2
	}
	divisor := transitionDivisor(modeTwo, currentSignature, previousSignature)
	category := transitionPackedCategoryCosts[(previousMetricCode>>14&3)*3+(currentMetricCode>>14&3)]
	return weighted/divisor + float32(category) + durationTerm, true
}

// TransitionScoreInput contains current/previous candidate values plus the
// feature and context terms produced elsewhere in the synthesis state.
// WeightRow follows the local row selected by the DLL's context-state branch.
// Model enables bank-crossing ordinal comparisons for continuity shortcuts.
type TransitionScoreInput struct {
	Current, Previous                      UnitRef
	Model                                  GlobalUnitOrdinalResolver
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

// Paul2013TransitionRecordContext carries state-owned transition terms that
// are not stored in the two unit index records.
type Paul2013TransitionRecordContext struct {
	Mode                                   byte
	DurationTerm, PreviousCumulative       float32
	AddTenPointContextRule                 bool
	CurrentSpecialContext, PreviousSpecial bool
}

// BuildPaul2013TransitionScoreInput selects the model columns consumed by
// FUN_10018c80 from the current context mode and two parsed unit records.
// Mode 0/1 uses current group 0 and previous group 2; mode 2 uses group 1 for
// both. Group feature pairs are kept opaque and copied in their observed
// first/second column order. Duration and context predicates remain explicit
// inputs because their producers depend on synthesis state outside the unit
// records. The model-backed path wrapper supplies the current node's /2-
// normalized local score as the duration term; direct scorer callers supply
// the already-scaled term.
func BuildPaul2013TransitionScoreInput(
	current UnitRef,
	currentRecord dat.UnitRecord,
	previous UnitRef,
	previousRecord dat.UnitRecord,
	context Paul2013TransitionRecordContext,
	currentState Paul2013TransitionContextState,
	previousState Paul2013TransitionContextState,
	model GlobalUnitOrdinalResolver,
) (TransitionScoreInput, error) {
	if context.Mode > 2 {
		return TransitionScoreInput{}, errors.New("Paul 2013 transition context mode is outside 0..2")
	}
	currentGroup, previousGroup := 0, 2
	if context.Mode == 2 {
		currentGroup, previousGroup = 1, 1
	}
	return TransitionScoreInput{
		Current:                current,
		Previous:               previous,
		Model:                  model,
		CurrentSignature:       currentRecord.Signature,
		PreviousSignature:      previousRecord.Signature,
		CurrentMetricCode:      currentRecord.MetricCodes[currentGroup],
		PreviousMetricCode:     previousRecord.MetricCodes[previousGroup],
		CurrentFeatureOne:      currentRecord.FeatureCodes[currentGroup][0],
		PreviousFeatureOne:     previousRecord.FeatureCodes[previousGroup][0],
		CurrentFeatureTwo:      currentRecord.FeatureCodes[currentGroup][1],
		PreviousFeatureTwo:     previousRecord.FeatureCodes[previousGroup][1],
		WeightRow:              Paul2013TransitionWeightRow(currentState, previousState),
		ModeTwo:                context.Mode == 2,
		DurationTerm:           context.DurationTerm,
		PreviousCumulative:     context.PreviousCumulative,
		AddTenPointContextRule: context.AddTenPointContextRule,
		CurrentSpecialContext:  context.CurrentSpecialContext,
		PreviousSpecial:        context.PreviousSpecial,
	}, nil
}

// GlobalUnitOrdinalResolver maps a bank-local reference into the selected
// model's ordered unit list. voice.Paul2013 implements this interface.
type GlobalUnitOrdinalResolver interface {
	GlobalUnitOrdinal(bank string, index uint32) (uint32, error)
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
	var globalOrdinals *[2]uint32
	if input.Model != nil {
		currentOrdinal, err := input.Model.GlobalUnitOrdinal(input.Current.Bank, input.Current.Index)
		if err != nil {
			return 0, err
		}
		previousOrdinal, err := input.Model.GlobalUnitOrdinal(input.Previous.Bank, input.Previous.Index)
		if err != nil {
			return 0, err
		}
		globalOrdinals = &[2]uint32{currentOrdinal, previousOrdinal}
	}
	tableCostSuppressed := suppressTransitionTableCostWithOrdinals(
		input.ModeTwo, input.Current, input.Previous, input.PreviousSignature, globalOrdinals,
	)
	distanceCost, ok := transitionDistanceCost(
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
		tableCostSuppressed,
	)
	if !ok {
		return 0, errors.New("transition weight row is outside the recovered table")
	}
	if input.ModeTwo {
		distanceCost += ModeTwoCategoricalPenalty(input.Current, input.Previous, input.CurrentSignature, input.PreviousSignature)
	} else if !tableCostSuppressed {
		distanceCost += ModeOtherCategoricalPenalty(input.CurrentSignature, input.PreviousSignature)
	}
	distanceCost += TransitionContextPenalty(input.AddTenPointContextRule, input.CurrentSpecialContext, input.PreviousSpecial)
	return distanceCost + input.PreviousCumulative, nil
}

// ScoreTransitionFromContextStates derives the DLL coefficient row from the
// two six-byte transition state records, then scores the edge using the same
// inputs as ScoreTransition. The producers of those state records remain the
// caller's responsibility.
func ScoreTransitionFromContextStates(
	table *distance.Table,
	featureTable *distance.FeatureTable,
	input TransitionScoreInput,
	currentState Paul2013TransitionContextState,
	previousState Paul2013TransitionContextState,
) (float32, error) {
	input.WeightRow = Paul2013TransitionWeightRow(currentState, previousState)
	return ScoreTransition(table, featureTable, input)
}

func suppressTransitionTableCost(modeTwo bool, current, previous UnitRef, previousSignature [7]byte) bool {
	return suppressTransitionTableCostWithOrdinals(modeTwo, current, previous, previousSignature, nil)
}

func suppressTransitionTableCostWithOrdinals(modeTwo bool, current, previous UnitRef, previousSignature [7]byte, globalOrdinals *[2]uint32) bool {
	if modeTwo {
		return current == previous
	}
	if previousSignature[6]&0x80 == 0 {
		return false
	}
	if globalOrdinals != nil {
		return globalOrdinals[1] < ^uint32(0) && globalOrdinals[0] == globalOrdinals[1]+1
	}
	return current.Bank == previous.Bank && previous.Index < ^uint32(0) && current.Index == previous.Index+1
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
