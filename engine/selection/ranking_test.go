package selection

import (
	"reflect"
	"testing"
)

func TestClassMismatchDistanceWeightsDifferingBytes(t *testing.T) {
	var input, candidate [5]byte
	input[0] = 1
	score, err := ClassMismatchDistance(input, candidate, 0)
	if err != nil {
		t.Fatal(err)
	}
	if score != 40 {
		t.Fatalf("key mismatch score = %d, want 40", score)
	}
	if score, err := ClassMismatchDistance(input, input, 0); err != nil || score != 0 {
		t.Fatalf("equal key score = %d, error = %v", score, err)
	}
}

func TestRankClassesPreservesOrderWhenLimitsFit(t *testing.T) {
	candidates := []ClassCandidate{
		{ID: 4, Population: 3}, {ID: 2, Population: 4}, {ID: 8, Population: 1},
	}
	ranked, err := RankClasses([5]byte{}, candidates, RankOptions{MaxClasses: 3, MaxUnits: 8})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ranked, candidates) {
		t.Fatalf("ranked classes = %+v, want source order %+v", ranked, candidates)
	}
}

func TestRankClassesScoresAndIncludesWholeThresholdCrossingClass(t *testing.T) {
	var oneByte1, oneByte0 [5]byte
	oneByte1[1] = 1
	oneByte0[0] = 1
	candidates := []ClassCandidate{
		{ID: 10, Key: oneByte0, Population: 5},
		{ID: 11, Key: oneByte1, Population: 4},
		{ID: 12, Population: 3},
	}
	ranked, err := RankClasses([5]byte{}, candidates, RankOptions{MaxClasses: 3, MaxUnits: 6})
	if err != nil {
		t.Fatal(err)
	}
	if len(ranked) != 2 || ranked[0].ID != 12 || ranked[1].ID != 11 {
		t.Fatalf("ranked class IDs = %v, want [12 11]", classIDs(ranked))
	}
	if got := ranked[0].Population + ranked[1].Population; got != 7 {
		t.Fatalf("selected population = %d, want threshold-crossing total 7", got)
	}
}

func TestRankClassesAppliesClassCountLimitAndRejectsInvalidOptions(t *testing.T) {
	candidates := []ClassCandidate{
		{ID: 1, Population: 1}, {ID: 2, Key: [5]byte{1}, Population: 1}, {ID: 3, Key: [5]byte{0, 1}, Population: 1},
	}
	ranked, err := RankClasses([5]byte{}, candidates, RankOptions{MaxClasses: 2, MaxUnits: 100})
	if err != nil {
		t.Fatal(err)
	}
	if got := classIDs(ranked); !reflect.DeepEqual(got, []uint32{1, 3}) {
		t.Fatalf("ranked class IDs = %v, want [1 3]", got)
	}
	for name, options := range map[string]RankOptions{
		"zero class limit": {MaxUnits: 1},
		"zero unit limit":  {MaxClasses: 1},
		"invalid view":     {MaxClasses: 1, MaxUnits: 1, ViewMode: 3},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := RankClasses([5]byte{}, candidates, options); err == nil {
				t.Fatal("invalid ranking options accepted")
			}
		})
	}
}

func classIDs(classes []ClassCandidate) []uint32 {
	ids := make([]uint32, len(classes))
	for i, class := range classes {
		ids[i] = class.ID
	}
	return ids
}
