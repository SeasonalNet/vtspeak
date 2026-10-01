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

func TestRankClassesPreservesEqualScoresInNativeInsertionPath(t *testing.T) {
	candidates := []ClassCandidate{
		{ID: 1, Key: [5]byte{1}, Population: 1},
		{ID: 2, Key: [5]byte{1}, Population: 1},
		{ID: 3, Key: [5]byte{1}, Population: 1},
	}
	ranked, err := RankClasses([5]byte{}, candidates, RankOptions{MaxClasses: 2, MaxUnits: 10})
	if err != nil {
		t.Fatal(err)
	}
	if got := classIDs(ranked); !reflect.DeepEqual(got, []uint32{1, 2}) {
		t.Fatalf("short equal-score class order = %v, want [1 2]", got)
	}
}

func TestSortPaul2013HeapOrdersAndRetainsEveryValue(t *testing.T) {
	values := []int{3, 10, 8, 7, 6, 5, 4, -2, 10, 0, 1}
	sortPaul2013Heap(values, func(left, right int) bool { return left < right })
	want := []int{-2, 0, 1, 3, 4, 5, 6, 7, 8, 10, 10}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("heap sort result = %v, want %v", values, want)
	}
}

func TestSortPaul2013NativeUsesHeapFallbackForHighlyUnbalancedPartition(t *testing.T) {
	values := make([]int, 100)
	for index := range values {
		values[index] = 3
	}
	values[0], values[len(values)/2], values[len(values)-1] = 0, 1, 2
	sortPaul2013Native(values, func(left, right int) bool { return left < right })
	if values[0] != 0 || values[1] != 1 || values[2] != 2 {
		t.Fatalf("heap-fallback prefix = %v, want [0 1 2]", values[:3])
	}
	for index := 3; index < len(values); index++ {
		if values[index] != 3 {
			t.Fatalf("heap-fallback value at %d = %d, want 3", index, values[index])
		}
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

func TestRankClassesRejectsOversizedKeyOnlyPools(t *testing.T) {
	const candidateLimit = 10_000
	candidates := make([]ClassCandidate, candidateLimit+1)
	for index := 0; index < candidateLimit; index++ {
		candidates[index] = ClassCandidate{
			ID:         uint32(index),
			Key:        [5]byte{1},
			Population: 1,
		}
	}
	// The key-only DLL path sorts the original count after scoring at most
	// 10,000 scratch entries. Reject this ambiguous tail rather than inventing
	// scores for it or dropping candidates.
	candidates[candidateLimit] = ClassCandidate{ID: candidateLimit, Population: 1}

	if _, err := RankClasses([5]byte{}, candidates, RankOptions{MaxClasses: 1, MaxUnits: candidateLimit}); err == nil {
		t.Fatal("class pool beyond the native score buffer was accepted")
	}
}

func TestRankClassesFeatureViewScoresOnlyNativeBufferPrefix(t *testing.T) {
	const candidateLimit = 10_000
	candidates := make([]ClassCandidate, candidateLimit+1)
	for index := 0; index < candidateLimit; index++ {
		candidates[index] = ClassCandidate{
			ID:         uint32(index),
			Key:        [5]byte{1},
			Population: 1,
		}
	}
	candidates[candidateLimit] = ClassCandidate{ID: candidateLimit, Population: 1}

	ranked, err := RankClasses([5]byte{}, candidates, RankOptions{
		MaxClasses: 1,
		MaxUnits:   candidateLimit,
		ViewMode:   1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ranked) != 1 || ranked[0].ID >= candidateLimit {
		t.Fatalf("ranked classes = %+v, want a class from the scored 10,000-entry prefix", ranked)
	}
}

func classIDs(classes []ClassCandidate) []uint32 {
	ids := make([]uint32, len(classes))
	for i, class := range classes {
		ids[i] = class.ID
	}
	return ids
}
