package selection

import (
	"testing"

	"vtspeak/engine/dat"
	"vtspeak/engine/distance"
)

func TestSelectPaul2013ModelMinimumCostPathBuildsAndScoresTransitions(t *testing.T) {
	model := continuityFixtureModel{
		counts: map[string]uint32{"gen": 3, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{
			0: {Signature: [7]byte{0, 0, 1}},
			1: {Signature: [7]byte{0, 0, 2}},
			2: {Signature: [7]byte{0, 0, 2}},
		},
	}
	table, err := distance.Parse([]byte{1, 0, 0, 0, 0, 0})
	if err != nil {
		t.Fatal(err)
	}
	passes := []Paul2013CandidateLocalScorePass{
		{ScoredPrefix: 2, Candidates: []Paul2013CandidateLocalScore{
			{Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Bank: "gen", Index: 0}}, LocalCost: 0, HasLocalCost: true},
			{Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Bank: "gen", Index: 1}}, LocalCost: 100, HasLocalCost: true},
		}},
		{ScoredPrefix: 1, Candidates: []Paul2013CandidateLocalScore{
			{Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Bank: "gen", Index: 2}}, LocalCost: 16, HasLocalCost: true},
		}},
	}
	path, err := SelectPaul2013PathFromLocalScorePasses(
		model, table, distance.NewFeatureTable(),
		passes,
		[]Paul2013PathLayerOptions{{ContextGateClear: true}},
		[]Paul2013PathTransitionContext{{RecordContext: Paul2013TransitionRecordContext{
			Mode: 2, DurationTerm: 999,
		}}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(path.Selected) != 2 || path.Selected[0] != (UnitRef{Bank: "gen", Index: 1}) ||
		path.Selected[1] != (UnitRef{Bank: "gen", Index: 2}) {
		t.Fatalf("selected model path = %+v, want gen:1 -> gen:2", path.Selected)
	}
	if got := path.Layers[1][0].CumulativeCost; got != 65 {
		t.Fatalf("selected cumulative cost = %g, want 65 (50 initial + 8 local + 7 transition)", got)
	}
}

func TestBuildPaul2013PathCandidateLayersRejectsUnscoredCandidate(t *testing.T) {
	_, err := BuildPaul2013PathCandidateLayers([]Paul2013CandidateLocalScorePass{{
		Candidates: []Paul2013CandidateLocalScore{{
			Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Bank: "gen", Index: 1}},
		}},
	}})
	if err == nil {
		t.Fatal("accepted a candidate without a recovered local path cost")
	}
}

func TestSelectPaul2013ModelMinimumCostPathValidatesEdgeInputs(t *testing.T) {
	candidates := [][]ScoredUnit{{{Unit: UnitRef{Bank: "gen", Index: 0}}}, {{Unit: UnitRef{Bank: "gen", Index: 1}}}}
	options := []Paul2013PathLayerOptions{{}}
	if _, err := SelectPaul2013ModelMinimumCostPath(nil, nil, nil, candidates, options, nil); err == nil {
		t.Fatal("accepted missing transition context")
	}
	if _, err := SelectPaul2013ModelMinimumCostPath(
		continuityFixtureModel{}, nil, nil, candidates, options,
		[]Paul2013PathTransitionContext{{RecordContext: Paul2013TransitionRecordContext{Mode: 3}}},
	); err == nil {
		t.Fatal("accepted an unsupported transition mode")
	}
}

func TestSelectPaul2013ModelMinimumCostPathNormalizesSingleLayer(t *testing.T) {
	unit := UnitRef{Bank: "gen", Index: 2}
	path, err := SelectPaul2013ModelMinimumCostPath(
		nil, nil, nil,
		[][]ScoredUnit{{{Unit: unit, Cost: 18}}},
		nil, nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(path.Selected) != 1 || path.Selected[0] != unit {
		t.Fatalf("single-layer selected path = %+v, want %s:%d", path.Selected, unit.Bank, unit.Index)
	}
	if got := path.Layers[0][0].CumulativeCost; got != 9 {
		t.Fatalf("single-layer initial cost = %g, want normalized local score 9", got)
	}
}
