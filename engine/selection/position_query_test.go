package selection

import (
	"reflect"
	"testing"

	"vtspeak/engine/dat"
	"vtspeak/engine/distance"
	"vtspeak/engine/voice"
)

func TestExpandPaul2013PositionQueryCandidatesKeepsWholePositionOrder(t *testing.T) {
	classes := []voice.ClassRecord{
		{ID: 8, Members: []voice.UnitLocation{{Bank: "gen", Index: 10}, {Bank: "gen", Index: 11}}},
		{ID: 3, Members: []voice.UnitLocation{{Bank: "num", Index: 4}}},
	}
	got, err := ExpandPaul2013PositionQueryCandidates(Paul2013PositionQueryResult{
		WholePosition: Paul2013WholePositionQueryResult{
			Selection: Paul2013WholePositionResult{Classes: classes, Accepted: true},
		},
		ReturnCount: 1,
	}, 10)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]UnitRef{{
		{Bank: "gen", Index: 10}, {Bank: "gen", Index: 11}, {Bank: "num", Index: 4},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expanded whole-position candidates = %+v, want %+v", got, want)
	}
}

func TestExpandPaul2013PositionQueryCandidatesKeepsFallbackRowsAndIDOrder(t *testing.T) {
	classA := voice.ClassRecord{ID: 7, Members: []voice.UnitLocation{{Bank: "gen", Index: 1}, {Bank: "gen", Index: 2}}}
	classB := voice.ClassRecord{ID: 4, Members: []voice.UnitLocation{{Bank: "num", Index: 9}}}
	result := Paul2013PositionQueryResult{
		UsedFallback: true,
		ReturnCount:  2,
		FallbackRows: [2]Paul2013FallbackRowResult{
			{
				RankedClasses: []voice.ClassRecord{classA},
				Candidates:    Paul2013FallbackCandidateList{IDs: []uint32{7, 4}, RankedCount: 1},
				QueryPasses:   []Paul2013CatalogQuerySequenceResult{{Candidates: []voice.ClassRecord{classB}}},
			},
			{
				RankedClasses: []voice.ClassRecord{classB},
				Candidates:    Paul2013FallbackCandidateList{IDs: []uint32{4, 7}, RankedCount: 1},
				QueryPasses:   []Paul2013CatalogQuerySequenceResult{{Candidates: []voice.ClassRecord{classA}}},
			},
		},
	}
	got, err := ExpandPaul2013PositionQueryCandidates(result, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]UnitRef{
		{{Bank: "gen", Index: 1}},
		{{Bank: "num", Index: 9}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expanded fallback candidate positions = %+v, want %+v", got, want)
	}
}

func TestExpandPaul2013PositionQueryCandidatesRejectsMissingFallbackRecord(t *testing.T) {
	_, err := ExpandPaul2013PositionQueryCandidates(Paul2013PositionQueryResult{
		UsedFallback: true,
		ReturnCount:  2,
		FallbackRows: [2]Paul2013FallbackRowResult{
			{Candidates: Paul2013FallbackCandidateList{IDs: []uint32{55}}},
			{Candidates: Paul2013FallbackCandidateList{IDs: []uint32{55}}},
		},
	}, 10)
	if err == nil {
		t.Fatal("accepted a fallback class ID without a source class record")
	}
}

func TestRankPaul2013QueriedCandidatePositionsConnectsExpansionAndContinuity(t *testing.T) {
	model := continuityFixtureModel{
		counts: map[string]uint32{"gen": 3, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{
			0: {},
			1: {},
		},
	}
	query := Paul2013PositionQueryResult{
		WholePosition: Paul2013WholePositionQueryResult{
			Selection: Paul2013WholePositionResult{
				Accepted: true,
				Classes: []voice.ClassRecord{{
					ID: 1,
					Members: []voice.UnitLocation{
						{Bank: "gen", Index: 0}, {Bank: "gen", Index: 1},
					},
				}},
			},
		},
		ReturnCount: 1,
	}
	got, err := RankPaul2013QueriedCandidatePositions(model, []Paul2013PositionQueryResult{query}, 10, []byte{0})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0]) != 2 {
		t.Fatalf("ranked queried positions = %+v, want one two-candidate position", got)
	}
	for _, candidate := range got[0] {
		if candidate.TotalSpan != 1 || candidate.WeightedSide != 1 {
			t.Fatalf("single-position continuity metadata = %+v, want span/side 1/1", candidate)
		}
	}
}

func TestRankAndScorePaul2013QueriedCandidatePositions(t *testing.T) {
	model := continuityFixtureModel{
		counts:  map[string]uint32{"gen": 2, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{0: {}, 1: {}},
	}
	query := Paul2013PositionQueryResult{
		WholePosition: Paul2013WholePositionQueryResult{
			Selection: Paul2013WholePositionResult{
				Accepted: true,
				Classes: []voice.ClassRecord{{
					ID: 1,
					Members: []voice.UnitLocation{
						{Bank: "gen", Index: 0}, {Bank: "gen", Index: 1},
					},
				}},
			},
		},
		ReturnCount: 1,
	}
	got, err := RankAndScorePaul2013QueriedCandidatePositions(
		model, distance.NewFeatureTable(), []Paul2013PositionQueryResult{query},
		10, []byte{0}, 1,
		[]Paul2013QueriedCandidateScoringInput{{Contexts: make([]Paul2013UnitScoreContext, 2)}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].Ranked) != 2 || len(got[0].Scores.Candidates) != 2 {
		t.Fatalf("queried local-score result = %+v, want one scored two-candidate position", got)
	}
	if got[0].HasProtection || got[0].Scores.ScoredPrefix != 2 {
		t.Fatalf("ordinary queried score pass = %+v, want full scoring without protection expansion", got[0])
	}
	for index, candidate := range got[0].Scores.Candidates {
		if !candidate.HasLocalCost || candidate.Continuity.Unit != got[0].Ranked[index].Unit {
			t.Errorf("scored candidate %d = %+v, want local cost aligned with ranked candidate", index, candidate)
		}
	}
}

func TestRankAndScorePaul2013QueriedCandidatePositionsUsesOversizedProtectionPath(t *testing.T) {
	members := make([]voice.UnitLocation, 31)
	for index := range members {
		members[index] = voice.UnitLocation{Bank: "gen", Index: uint32(index)}
	}
	model := continuityFixtureModel{
		counts:  map[string]uint32{"gen": 31, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{},
	}
	query := Paul2013PositionQueryResult{
		WholePosition: Paul2013WholePositionQueryResult{
			Selection: Paul2013WholePositionResult{
				Accepted: true,
				Classes:  []voice.ClassRecord{{ID: 1, Members: members}},
			},
		},
		ReturnCount: 1,
	}
	keys := make([]uint32, len(members))
	for index := range keys {
		keys[index] = uint32(index + 1)
	}
	got, err := RankAndScorePaul2013QueriedCandidatePositions(
		model, distance.NewFeatureTable(), []Paul2013PositionQueryResult{query},
		31, []byte{0}, 1,
		[]Paul2013QueriedCandidateScoringInput{{
			Contexts: make([]Paul2013UnitScoreContext, 31), CandidateKeys: keys,
		}},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || !got[0].HasProtection || len(got[0].Protection.Candidates) != 31 {
		t.Fatalf("oversized queried score pass = %+v, want protection metadata for 31 candidates", got)
	}
	if got[0].Scores.ScoredPrefix != 31 {
		t.Fatalf("oversized queried score prefix = %d, want 31", got[0].Scores.ScoredPrefix)
	}
	shortlist, err := FinalizePaul2013QueriedCandidateShortlist(
		got[0], 0, make([]byte, 31), make([]byte, 31),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(shortlist.Candidates) != 30 {
		t.Fatalf("final queried shortlist has %d candidates, want native cap 30", len(shortlist.Candidates))
	}
}

func TestSelectPaul2013QueriedCandidatePathComposesSinglePosition(t *testing.T) {
	model := continuityFixtureModel{
		counts:  map[string]uint32{"gen": 2, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{0: {}, 1: {}},
	}
	query := Paul2013PositionQueryResult{
		WholePosition: Paul2013WholePositionQueryResult{
			Selection: Paul2013WholePositionResult{
				Accepted: true,
				Classes: []voice.ClassRecord{{
					ID: 1,
					Members: []voice.UnitLocation{
						{Bank: "gen", Index: 0}, {Bank: "gen", Index: 1},
					},
				}},
			},
		},
		ReturnCount: 1,
	}
	got, err := SelectPaul2013QueriedCandidatePath(Paul2013QueriedCandidatePathInput{
		Model: model, FeatureDistanceTable: distance.NewFeatureTable(),
		Queries: []Paul2013PositionQueryResult{query}, MaxUnits: 10,
		Modes: []byte{0}, PositionCount: 1,
		ScoringInputs: []Paul2013QueriedCandidateScoringInput{{Contexts: make([]Paul2013UnitScoreContext, 2)}},
		TailStarts:    []int{0}, CandidateClasses: [][]byte{nil}, ContextClasses: [][]byte{nil},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Shortlists) != 1 || len(got.Shortlists[0].Candidates) != 2 || len(got.Path.Selected) != 1 {
		t.Fatalf("queried path result = %+v, want one two-candidate shortlist and one selected unit", got)
	}
	if got.Path.Selected[0].Bank != "gen" {
		t.Fatalf("selected queried unit = %+v, want bank gen", got.Path.Selected[0])
	}
}

func TestSelectPaul2013QueriedCandidatePathScoresAdjacentPositions(t *testing.T) {
	model := continuityFixtureModel{
		counts: map[string]uint32{"gen": 3, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{
			0: {Signature: [7]byte{0, 0, 1}},
			2: {Signature: [7]byte{0, 0, 2}},
		},
	}
	queries := make([]Paul2013PositionQueryResult, 2)
	for position, unitIndex := range []uint32{0, 2} {
		queries[position] = Paul2013PositionQueryResult{
			WholePosition: Paul2013WholePositionQueryResult{
				Selection: Paul2013WholePositionResult{
					Accepted: true,
					Classes: []voice.ClassRecord{{
						ID:      uint32(position + 1),
						Members: []voice.UnitLocation{{Bank: "gen", Index: unitIndex}},
					}},
				},
			},
			ReturnCount: 1,
		}
	}
	transitionTable, err := distance.Parse([]byte{1, 0, 0, 0, 0, 0})
	if err != nil {
		t.Fatal(err)
	}
	got, err := SelectPaul2013QueriedCandidatePath(Paul2013QueriedCandidatePathInput{
		Model: model, TransitionDistanceTable: transitionTable,
		FeatureDistanceTable: distance.NewFeatureTable(), Queries: queries,
		MaxUnits: 2, Modes: []byte{0, 0}, PositionCount: 2,
		ScoringInputs: []Paul2013QueriedCandidateScoringInput{
			{Contexts: make([]Paul2013UnitScoreContext, 1)},
			{Contexts: make([]Paul2013UnitScoreContext, 1)},
		},
		TailStarts: []int{0, 0}, CandidateClasses: [][]byte{nil, nil}, ContextClasses: [][]byte{nil, nil},
		LayerOptions: []Paul2013PathLayerOptions{{ContextGateClear: true}},
		Transitions:  []Paul2013PathTransitionContext{{RecordContext: Paul2013TransitionRecordContext{Mode: 2}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []UnitRef{{Bank: "gen", Index: 0}, {Bank: "gen", Index: 2}}
	if len(got.Path.Selected) != len(want) {
		t.Fatalf("queried adjacent path = %+v, want %d positions", got.Path.Selected, len(want))
	}
	for index := range want {
		if got.Path.Selected[index] != want[index] {
			t.Errorf("queried path position %d = %+v, want %+v", index, got.Path.Selected[index], want[index])
		}
	}
	if len(got.Path.Layers) != 2 || got.Path.Layers[1][0].PreviousIndex != 0 {
		t.Fatalf("queried path predecessor chain = %+v, want one recorded edge", got.Path.Layers)
	}
}
