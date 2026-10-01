package selection

import (
	"encoding/binary"
	"errors"
	"math"
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/distance"
	"vtspeak/engine/voice"
)

func TestTransitionDistanceCostCanSuppressWeightedTableSubtotal(t *testing.T) {
	current := UnitRef{Bank: "gen", Index: 1}
	previous := UnitRef{Bank: "gen", Index: 0}
	var previousSignature [7]byte
	previousSignature[6] = 0x80
	if !suppressTransitionTableCost(false, current, previous, previousSignature) {
		t.Fatal("consecutive row with predecessor marker did not suppress the transition table subtotal")
	}

	var currentSignature [7]byte
	cost, ok := transitionDistanceCost(
		false, 2, 0.5, 0.25, 0.5, 3, 0x4000, 0,
		currentSignature, previousSignature, true,
	)
	if !ok {
		t.Fatal("valid transition weight row was rejected")
	}
	if cost != 103 {
		t.Fatalf("suppressed transition distance cost = %v, want 103 (packed category 100 plus duration 3)", cost)
	}

	unsuppressed, ok := transitionDistanceCost(
		false, 2, 0.5, 0.25, 0.5, 3, 0x4000, 0,
		currentSignature, previousSignature, false,
	)
	if !ok || math.Abs(float64(unsuppressed-113)) > 0.0001 {
		t.Fatalf("ordinary transition cost = %v, want 113 (got ok=%v)", unsuppressed, ok)
	}
}

func TestScoreTransitionAppliesObservedContinuityShortcut(t *testing.T) {
	distanceBytes := make([]byte, 2+3*4)
	binary.LittleEndian.PutUint16(distanceBytes[:2], 2)
	binary.LittleEndian.PutUint32(distanceBytes[2+4:], math.Float32bits(0.5))
	table, err := distance.Parse(distanceBytes)
	if err != nil {
		t.Fatal(err)
	}
	features := distance.NewFeatureTable()

	current := UnitRef{Bank: "gen", Index: 64729}
	previous := UnitRef{Bank: "gen", Index: 64728}
	var currentSignature, previousSignature [7]byte
	currentSignature[2] = 73
	previousSignature[1] = 21
	previousSignature[6] = 0x80
	input := TransitionScoreInput{
		Current: current, Previous: previous,
		CurrentSignature: currentSignature, PreviousSignature: previousSignature,
		CurrentMetricCode: 1, PreviousMetricCode: 0,
		WeightRow: 2, DurationTerm: 3, PreviousCumulative: 7,
	}
	cost, err := ScoreTransition(table, features, input)
	if err != nil {
		t.Fatal(err)
	}
	if cost != 10 {
		t.Fatalf("continuity-gated transition score = %v, want 10 (cumulative 7 plus duration 3)", cost)
	}

	input.PreviousSignature[6] = 0
	withoutMarker, err := ScoreTransition(table, features, input)
	if err != nil {
		t.Fatal(err)
	}
	if withoutMarker <= cost {
		t.Fatalf("clearing predecessor marker did not restore edge costs: marked %v, clear %v", cost, withoutMarker)
	}
}

func TestScoreTransitionFromContextStatesDerivesWeightRow(t *testing.T) {
	distanceBytes := make([]byte, 2+3*4)
	binary.LittleEndian.PutUint16(distanceBytes[:2], 2)
	binary.LittleEndian.PutUint32(distanceBytes[2+4:], math.Float32bits(0.5))
	table, err := distance.Parse(distanceBytes)
	if err != nil {
		t.Fatal(err)
	}
	features := distance.NewFeatureTable()
	input := TransitionScoreInput{
		Current:           UnitRef{Bank: "gen", Index: 2},
		Previous:          UnitRef{Bank: "gen", Index: 0},
		CurrentMetricCode: 1, PreviousMetricCode: 0,
	}
	currentState := Paul2013TransitionContextState{1, 0, 0, 0, 0, 0}
	previousState := Paul2013TransitionContextState{2, 0, 0, 0, 0, 0}
	got, err := ScoreTransitionFromContextStates(table, features, input, currentState, previousState)
	if err != nil {
		t.Fatal(err)
	}

	input.WeightRow = 4
	want, err := ScoreTransition(table, features, input)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("context-derived transition score = %v, explicit row-four score = %v", got, want)
	}
}

func TestSuppressTransitionTableCostRequiresObservedIdentityAndMarker(t *testing.T) {
	previous := UnitRef{Bank: "gen", Index: 4}
	var marked [7]byte
	marked[6] = 0x80
	for name, test := range map[string]struct {
		modeTwo  bool
		current  UnitRef
		previous UnitRef
		marker   [7]byte
		want     bool
	}{
		"same candidate in mode two":   {modeTwo: true, current: previous, previous: previous, want: true},
		"consecutive marked candidate": {current: UnitRef{Bank: "gen", Index: 5}, previous: previous, marker: marked, want: true},
		"consecutive without marker":   {current: UnitRef{Bank: "gen", Index: 5}, previous: previous},
		"nonconsecutive candidate":     {current: UnitRef{Bank: "gen", Index: 6}, previous: previous, marker: marked},
		"different bank":               {current: UnitRef{Bank: "etc", Index: 5}, previous: previous, marker: marked},
	} {
		t.Run(name, func(t *testing.T) {
			got := suppressTransitionTableCost(test.modeTwo, test.current, test.previous, test.marker)
			if got != test.want {
				t.Fatalf("suppression = %v, want %v", got, test.want)
			}
		})
	}
}

func TestSuppressTransitionTableCostUsesGlobalOrdinalsAcrossBanks(t *testing.T) {
	var marked [7]byte
	marked[6] = 0x80
	for name, test := range map[string]struct {
		current  UnitRef
		previous UnitRef
		ordinals [2]uint32
		want     bool
	}{
		"adjacent across gen to num": {
			current: UnitRef{Bank: "num", Index: 0}, previous: UnitRef{Bank: "gen", Index: 9},
			ordinals: [2]uint32{10, 9}, want: true,
		},
		"adjacent across num to etc": {
			current: UnitRef{Bank: "etc", Index: 0}, previous: UnitRef{Bank: "num", Index: 3},
			ordinals: [2]uint32{14, 13}, want: true,
		},
		"not globally adjacent": {
			current: UnitRef{Bank: "etc", Index: 0}, previous: UnitRef{Bank: "num", Index: 3},
			ordinals: [2]uint32{15, 13},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := suppressTransitionTableCostWithOrdinals(false, test.current, test.previous, marked, &test.ordinals)
			if got != test.want {
				t.Fatalf("global-ordinal suppression = %v, want %v", got, test.want)
			}
		})
	}
}

func TestScoreTransitionUsesModelGlobalOrdinalAcrossBanks(t *testing.T) {
	distanceBytes := make([]byte, 2+3*4)
	binary.LittleEndian.PutUint16(distanceBytes[:2], 2)
	binary.LittleEndian.PutUint32(distanceBytes[2+4:], math.Float32bits(0.5))
	table, err := distance.Parse(distanceBytes)
	if err != nil {
		t.Fatal(err)
	}
	input := TransitionScoreInput{
		Current: UnitRef{Bank: "num", Index: 0}, Previous: UnitRef{Bank: "gen", Index: 9},
		CurrentMetricCode: 1, PreviousMetricCode: 0,
		WeightRow: 2, DurationTerm: 3, PreviousCumulative: 7,
	}
	input.PreviousSignature[6] = 0x80
	input.Model = testOrdinalResolver{{Bank: "num", Index: 0}: 10, {Bank: "gen", Index: 9}: 9}

	cost, err := ScoreTransition(table, distance.NewFeatureTable(), input)
	if err != nil {
		t.Fatal(err)
	}
	if cost != 10 {
		t.Fatalf("cross-bank consecutive transition score = %v, want 10 (cumulative 7 plus duration 3)", cost)
	}
}

func TestScoreTransitionUsesPaulBankBoundaryWhenAssetsArePresent(t *testing.T) {
	dataRoot := filepath.Join("..", "..", "data-paul", "M16")
	if _, err := os.Stat(filepath.Join(dataRoot, "dblist.idx")); err != nil {
		t.Skip("local Paul model files are unavailable")
	}
	model, err := voice.OpenPaul2013(dataRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer model.Close()

	distanceBytes := make([]byte, 2+3*4)
	binary.LittleEndian.PutUint16(distanceBytes[:2], 2)
	binary.LittleEndian.PutUint32(distanceBytes[2+4:], math.Float32bits(0.5))
	table, err := distance.Parse(distanceBytes)
	if err != nil {
		t.Fatal(err)
	}
	previous := UnitRef{Bank: "gen", Index: model.Banks["gen"].UnitCount() - 1}
	current := UnitRef{Bank: "num", Index: 0}
	previousOrdinal, err := model.GlobalUnitOrdinal(previous.Bank, previous.Index)
	if err != nil {
		t.Fatal(err)
	}
	currentOrdinal, err := model.GlobalUnitOrdinal(current.Bank, current.Index)
	if err != nil {
		t.Fatal(err)
	}
	if currentOrdinal != previousOrdinal+1 {
		t.Fatalf("loaded bank boundary ordinals = %d -> %d, want consecutive", previousOrdinal, currentOrdinal)
	}
	input := TransitionScoreInput{
		Model: model, Current: current, Previous: previous,
		CurrentMetricCode: 1, PreviousMetricCode: 0,
		WeightRow: 2, DurationTerm: 3, PreviousCumulative: 7,
	}
	input.PreviousSignature[6] = 0x80
	cost, err := ScoreTransition(table, distance.NewFeatureTable(), input)
	if err != nil {
		t.Fatal(err)
	}
	if cost != 10 {
		t.Fatalf("loaded cross-bank transition score = %v, want 10", cost)
	}
}

type testOrdinalResolver map[UnitRef]uint32

func (resolver testOrdinalResolver) GlobalUnitOrdinal(bank string, index uint32) (uint32, error) {
	ordinal, ok := resolver[UnitRef{Bank: bank, Index: index}]
	if !ok {
		return 0, errors.New("unit reference has no test ordinal")
	}
	return ordinal, nil
}
