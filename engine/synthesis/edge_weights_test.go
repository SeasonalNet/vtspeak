package synthesis

import (
	"reflect"
	"testing"
)

func TestBuildPaul2013ContextEdgePlanCombinesGateAndCaps(t *testing.T) {
	rows := make([][7]byte, 5)
	rows[2][6] = 0x80
	rows[3][6] = 0x80
	gate := ContextGateInput{
		Rows: rows, Modes: []byte{0, 0, 2, 0, 0}, CurrentIndex: 2,
		LeftIndex: 0, RightIndex: 4, TimelineSpeedControl: 100,
	}
	got, err := BuildPaul2013ContextEdgePlan(gate, 8, 6, 2)
	if err != nil {
		t.Fatal(err)
	}
	if got.Multipliers != (ContextMultipliers{Left: 2, Right: 2}) {
		t.Fatalf("context multipliers = %+v, want both sides at 2", got.Multipliers)
	}
	if got.LeftContextCount != 5 || got.RightContextCount != 2 {
		t.Fatalf("capped context counts = %d/%d, want 5/2", got.LeftContextCount, got.RightContextCount)
	}
	if got.Weights.Normalization != 40 || len(got.Weights.Left) != 9 {
		t.Fatalf("edge plan normalization/size = %d/%d, want 40/9", got.Weights.Normalization, len(got.Weights.Left))
	}
	for index := range got.Weights.Left {
		if got.Weights.Left[index]+got.Weights.Current[index]+got.Weights.Right[index] != got.Weights.Normalization {
			t.Fatalf("edge weights at %d do not sum to normalization", index)
		}
	}
}

func TestBuildPaul2013ContextEdgePlanZerosSuppressedAndEmptySides(t *testing.T) {
	rows := make([][7]byte, 3)
	gate := ContextGateInput{
		Rows: rows, Modes: []byte{0, 0, 0}, CurrentIndex: 1,
		LeftIndex: 0, RightIndex: 2, TimelineSpeedControl: 100,
	}
	got, err := BuildPaul2013ContextEdgePlan(gate, 4, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if got.Multipliers != (ContextMultipliers{}) || got.LeftContextCount != 0 || got.RightContextCount != 0 {
		t.Fatalf("suppressed context plan retained side state: %+v", got)
	}
	if got.Weights.Normalization != 1 {
		t.Fatalf("suppressed plan normalization = %d, want 1", got.Weights.Normalization)
	}
	for index := range got.Weights.Current {
		if got.Weights.Current[index] != 1 || got.Weights.Left[index] != 0 || got.Weights.Right[index] != 0 {
			t.Fatalf("suppressed plan weights at %d = L:%d C:%d R:%d", index, got.Weights.Left[index], got.Weights.Current[index], got.Weights.Right[index])
		}
	}

	rows[2][6] = 0x80
	gate.Rows = rows
	gate.Modes[1] = 2
	got, err = BuildPaul2013ContextEdgePlan(gate, 4, 0, 3)
	if err != nil {
		t.Fatal(err)
	}
	if got.Multipliers.Left != 0 || got.LeftContextCount != 0 {
		t.Fatalf("empty eligible side retained nonzero multiplier/count: %+v", got)
	}
}

func TestBuildPaul2013ContextEdgePlanRejectsInvalidCounts(t *testing.T) {
	gate := ContextGateInput{
		Rows: make([][7]byte, 1), Modes: []byte{0}, CurrentIndex: 0,
	}
	if _, err := BuildPaul2013ContextEdgePlan(gate, 0, 0, 0); err == nil {
		t.Fatal("accepted a zero-period current unit")
	}
	if _, err := BuildPaul2013ContextEdgePlan(gate, 3, 256, 0); err == nil {
		t.Fatal("accepted a neighbor count wider than the native byte")
	}
}

func TestMixPaul2013ContextUPMTimelineMatchesExplicitPlan(t *testing.T) {
	gate := ContextGateInput{
		Rows: make([][7]byte, 1), Modes: []byte{0}, CurrentIndex: 0,
		TimelineSpeedControl: 100,
	}
	periods := []Paul2013UPMPeriod{{
		Length: 4,
		Windows: Paul2013UPMPeriodWindows{
			Current: []int16{100, 200, 300, 400},
		},
	}}
	plan, err := BuildPaul2013ContextEdgePlan(gate, len(periods), 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]int16, 4)
	if err := MixPaul2013UPMTimeline(want, 0, plan.Weights, periods); err != nil {
		t.Fatal(err)
	}
	got := make([]int16, 4)
	composedPlan, err := MixPaul2013ContextUPMTimeline(got, 0, gate, 0, 0, periods)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(composedPlan, plan) || !reflect.DeepEqual(got, want) {
		t.Fatalf("composed mix plan/audio = %+v/%v, explicit plan/audio = %+v/%v", composedPlan, got, plan, want)
	}
}
