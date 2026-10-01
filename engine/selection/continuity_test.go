package selection

import (
	"fmt"
	"math"
	"testing"

	"vtspeak/engine/dat"
	"vtspeak/engine/distance"
)

type continuityFixtureModel struct {
	counts  map[string]uint32
	records map[uint32]dat.UnitRecord
}

func (m continuityFixtureModel) GlobalUnitOrdinal(bank string, index uint32) (uint32, error) {
	var offset uint32
	for _, name := range []string{"gen", "num", "etc", "alp"} {
		count, ok := m.counts[name]
		if !ok || count == 0 {
			return 0, fmt.Errorf("missing bank %q", name)
		}
		if name == bank {
			if index >= count {
				return 0, fmt.Errorf("index out of range")
			}
			return offset + index, nil
		}
		offset += count
	}
	return 0, fmt.Errorf("unknown bank %q", bank)
}

func (m continuityFixtureModel) UnitAtGlobalOrdinal(ordinal uint32) (string, uint32, error) {
	var offset uint32
	for _, name := range []string{"gen", "num", "etc", "alp"} {
		count := m.counts[name]
		if ordinal < offset+count {
			return name, ordinal - offset, nil
		}
		offset += count
	}
	return "", 0, fmt.Errorf("ordinal out of range")
}

func (m continuityFixtureModel) ReadRecord(bank string, index uint32) (dat.UnitRecord, error) {
	ordinal, err := m.GlobalUnitOrdinal(bank, index)
	if err != nil {
		return dat.UnitRecord{}, err
	}
	return m.records[ordinal], nil
}

func TestRankPaul2013CandidatesByContinuityAcrossBankBoundary(t *testing.T) {
	model := continuityFixtureModel{
		counts: map[string]uint32{"gen": 10, "num": 10, "etc": 10, "alp": 10},
		records: map[uint32]dat.UnitRecord{
			9:  {Signature: [7]byte{0, 0, 0, 0, 0, 0, 0x80}},
			10: {Signature: [7]byte{0, 0, 0, 0, 0, 0, 0x80}},
		},
	}
	positions := [][]UnitRef{
		{{Bank: "gen", Index: 9}},
		{{Bank: "num", Index: 0}},
		{{Bank: "num", Index: 1}},
	}
	got, err := RankPaul2013CandidatesByContinuity(model, positions, []byte{0, 0, 0}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("ranked %d candidates, want 1", len(got))
	}
	if got[0].RightMatches != 1 || got[0].TotalSpan != 3 || got[0].WeightedSide != 3 || got[0].InitialOrderingKey != -303 {
		t.Fatalf("continuity result = %+v, want right=1 span=3 weighted=3 key=-303", got[0])
	}
}

func TestRankPaul2013CandidatesByContinuityHonorsModesAndMarkers(t *testing.T) {
	model := continuityFixtureModel{
		counts:  map[string]uint32{"gen": 10, "num": 10, "etc": 10, "alp": 10},
		records: map[uint32]dat.UnitRecord{},
	}
	positions := [][]UnitRef{
		{{Bank: "gen", Index: 5}},
		{{Bank: "gen", Index: 5}},
		{{Bank: "gen", Index: 5}},
	}
	// Modes 1 and 2 preserve the same unit ID on their respective side, so
	// they do not consult the continuity marker.
	got, err := RankPaul2013CandidatesByContinuity(model, positions, []byte{1, 0, 2}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got[0].TotalSpan != 3 || got[0].WeightedSide != 1 {
		t.Fatalf("mode-preserved continuity = %+v, want full span with bilateral bonus", got[0])
	}

	// A mode-0 traversal to a neighboring ordinal stops when its signature
	// marker is absent, even if that unit appears in the neighboring list.
	got, err = RankPaul2013CandidatesByContinuity(model, positions, []byte{0, 0, 0}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got[0].TotalSpan != 1 {
		t.Fatalf("unmarked adjacency produced total span %d, want 1", got[0].TotalSpan)
	}
}

func TestRankPaul2013CandidatesByContinuityValidatesInputs(t *testing.T) {
	model := continuityFixtureModel{
		counts: map[string]uint32{"gen": 2, "num": 2, "etc": 2, "alp": 2},
	}
	for _, test := range []struct {
		name      string
		positions [][]UnitRef
		modes     []byte
		position  int
	}{
		{name: "mode length mismatch", positions: [][]UnitRef{{}}, modes: nil, position: 0},
		{name: "position out of range", positions: [][]UnitRef{{}}, modes: []byte{0}, position: 1},
		{name: "invalid mode", positions: [][]UnitRef{{}}, modes: []byte{3}, position: 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := RankPaul2013CandidatesByContinuity(model, test.positions, test.modes, test.position); err == nil {
				t.Fatal("invalid continuity input was accepted")
			}
		})
	}
}

func TestPaul2013ContinuityFeatureScale(t *testing.T) {
	for _, test := range []struct {
		span     uint16
		weighted uint16
		want     float32
	}{
		{span: 1, weighted: 0, want: 1},
		{span: 3, weighted: 2, want: 1.0 / 3.5},
		{span: 4, weighted: 1, want: 1.0 / 3.0},
	} {
		got, err := Paul2013ContinuityFeatureScale(test.span, test.weighted)
		if err != nil {
			t.Fatalf("scale for span=%d weighted=%d: %v", test.span, test.weighted, err)
		}
		if math.Abs(float64(got-test.want)) > 0.000001 {
			t.Errorf("scale for span=%d weighted=%d = %v, want %v", test.span, test.weighted, got, test.want)
		}
	}
	if _, err := Paul2013ContinuityFeatureScale(0, 0); err == nil {
		t.Fatal("accepted zero continuity span")
	}
}

func TestPaul2013CandidateLocalScorePrefix(t *testing.T) {
	tests := []struct {
		name          string
		candidates    []Paul2013CandidateContinuity
		positionCount int
		want          int
	}{
		{
			name: "leading full-span group only",
			candidates: []Paul2013CandidateContinuity{
				{TotalSpan: 4, InitialOrderingKey: -401},
				{TotalSpan: 4, InitialOrderingKey: -400},
				{TotalSpan: 3, InitialOrderingKey: -303},
			},
			positionCount: 4,
			want:          2,
		},
		{
			name: "small non-full set scores all",
			candidates: []Paul2013CandidateContinuity{
				{TotalSpan: 3, InitialOrderingKey: -303},
				{TotalSpan: 2, InitialOrderingKey: -202},
			},
			positionCount: 4,
			want:          2,
		},
		{
			name: "large set stops after thirty",
			candidates: func() []Paul2013CandidateContinuity {
				result := make([]Paul2013CandidateContinuity, 31)
				for index := range result {
					result[index] = Paul2013CandidateContinuity{TotalSpan: 1, InitialOrderingKey: -1}
				}
				result[30].InitialOrderingKey = 0
				return result
			}(),
			positionCount: 4,
			want:          30,
		},
		{
			name: "large set retains key tie at boundary",
			candidates: func() []Paul2013CandidateContinuity {
				result := make([]Paul2013CandidateContinuity, 32)
				for index := range result {
					result[index] = Paul2013CandidateContinuity{TotalSpan: 1, InitialOrderingKey: -1}
				}
				result[31].InitialOrderingKey = 0
				return result
			}(),
			positionCount: 4,
			want:          31,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Paul2013CandidateLocalScorePrefix(test.candidates, test.positionCount)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("local-score prefix = %d, want %d", got, test.want)
			}
		})
	}
}

func TestExpandPaul2013CandidateProtectionPrefixCopiesBoundedTailMatches(t *testing.T) {
	candidates := make([]Paul2013CandidateContinuity, 40)
	keys := make([]uint32, len(candidates))
	for index := range candidates {
		keys[index] = uint32(index + 100)
		candidates[index] = Paul2013CandidateContinuity{
			Unit: UnitRef{Bank: "gen", Index: uint32(index)}, TotalSpan: 1,
			InitialOrderingKey: float32(index),
		}
	}
	for index := 32; index < 38; index++ {
		candidates[index].InitialOrderingKey = 5
	}
	for index := 38; index < len(candidates); index++ {
		candidates[index].InitialOrderingKey = 6
	}
	wholePositionKeys := []uint32{keys[0], keys[1], keys[2], keys[3]}
	for index := 32; index < len(candidates); index++ {
		wholePositionKeys = append(wholePositionKeys, keys[index])
	}

	got, err := ExpandPaul2013CandidateProtectionPrefix(candidates, keys, wholePositionKeys, 30, 10)
	if err != nil {
		t.Fatal(err)
	}
	if got.ScoredPrefix != 36 || got.Added != 6 || got.ProtectedCount != 10 {
		t.Fatalf("expanded prefix/protected/added = %d/%d/%d, want 36/10/6", got.ScoredPrefix, got.ProtectedCount, got.Added)
	}
	for offset := 0; offset < 6; offset++ {
		sourceIndex := 32 + offset
		destinationIndex := 30 + offset
		if got.ScoredPrefixSources[destinationIndex] != sourceIndex ||
			got.Candidates[destinationIndex].Unit.Index != uint32(sourceIndex) ||
			got.CandidateKeys[destinationIndex] != keys[sourceIndex] || got.NodeFlags[destinationIndex] != 1 {
			t.Fatalf("expanded candidate %d did not copy source %d into slot %d: %+v", offset, sourceIndex, destinationIndex, got)
		}
	}
	if got.Candidates[30].InitialOrderingKey != candidates[30].InitialOrderingKey {
		t.Fatalf("native copy changed destination ordering key to %g", got.Candidates[30].InitialOrderingKey)
	}
	if got.NodeFlags[36] != 0 || got.Candidates[36].Unit.Index != 36 {
		t.Fatalf("scan past the tenth flag's ordering-key boundary: flag=%d unit=%d", got.NodeFlags[36], got.Candidates[36].Unit.Index)
	}
	if candidates[30].Unit.Index != 30 || keys[30] != 130 {
		t.Fatal("protection expansion modified caller input")
	}
}

func TestExpandPaul2013CandidateProtectionPrefixSkipsAndValidates(t *testing.T) {
	candidates := make([]Paul2013CandidateContinuity, 31)
	keys := make([]uint32, len(candidates))
	for index := range candidates {
		keys[index] = uint32(index)
		candidates[index].Unit.Index = uint32(index)
	}
	whole := []uint32{0, 1}
	got, err := ExpandPaul2013CandidateProtectionPrefix(candidates, keys, whole, 2, 9)
	if err != nil {
		t.Fatal(err)
	}
	if got.Added != 0 || got.ScoredPrefix != 2 || got.ProtectedCount != 2 {
		t.Fatalf("expanded despite whole-position pool at the threshold: %+v", got)
	}
	whole = []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	got, err = ExpandPaul2013CandidateProtectionPrefix(candidates, keys, whole, 10, 12)
	if err != nil {
		t.Fatal(err)
	}
	if got.Added != 0 || got.ProtectedCount != 10 {
		t.Fatalf("expanded with ten prefix candidates already protected: %+v", got)
	}
	if _, err := ExpandPaul2013CandidateProtectionPrefix(candidates, keys[:1], whole, 2, 12); err == nil {
		t.Fatal("accepted misaligned candidate keys")
	}
	if _, err := ExpandPaul2013CandidateProtectionPrefix(candidates, keys, whole, 32, 12); err == nil {
		t.Fatal("accepted an initial prefix beyond the candidate list")
	}
}

func TestScorePaul2013ExpandedCandidateLocalPrefixRealignsContexts(t *testing.T) {
	model := continuityFixtureModel{
		counts:  map[string]uint32{"gen": 40, "num": 1, "etc": 1, "alp": 1},
		records: make(map[uint32]dat.UnitRecord),
	}
	candidates := make([]Paul2013CandidateContinuity, 40)
	contexts := make([]Paul2013UnitScoreContext, len(candidates))
	keys := make([]uint32, len(candidates))
	for index := range candidates {
		keys[index] = uint32(index + 100)
		orderingKey := float32(-1)
		if index >= 30 {
			orderingKey = 0
		}
		if index >= 38 {
			orderingKey = 1
		}
		candidates[index] = Paul2013CandidateContinuity{
			Unit: UnitRef{Bank: "gen", Index: uint32(index)}, TotalSpan: 1,
			InitialOrderingKey: orderingKey,
		}
		model.records[uint32(index)] = dat.UnitRecord{}
	}
	wholePositionKeys := []uint32{keys[0], keys[1], keys[2], keys[3]}
	for index := 32; index < len(candidates); index++ {
		wholePositionKeys = append(wholePositionKeys, keys[index])
	}

	got, err := ScorePaul2013ExpandedCandidateLocalPrefix(
		model, distance.NewFeatureTable(), candidates, 3, contexts, keys, wholePositionKeys, 10,
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.Scores.ScoredPrefix != 36 || got.Protection.ScoredPrefixSources[30] != 32 {
		t.Fatalf("expanded scored prefix/source = %d/%d, want 36/32", got.Scores.ScoredPrefix, got.Protection.ScoredPrefixSources[30])
	}
	for index := 0; index < got.Scores.ScoredPrefix; index++ {
		if !got.Scores.Candidates[index].HasLocalCost {
			t.Errorf("expanded prefix candidate %d was not locally scored", index)
		}
	}
	if got.Scores.Candidates[30].Continuity.Unit.Index != 32 || got.Scores.Candidates[30].NodeFlag != 1 {
		t.Fatalf("first copied candidate = %+v, want protected unit 32", got.Scores.Candidates[30])
	}
	if got.Scores.Candidates[36].NodeFlag != 0 || got.Scores.Candidates[36].HasLocalCost {
		t.Fatalf("candidate after the protected ordering boundary = %+v", got.Scores.Candidates[36])
	}
}

func TestPaul2013CandidateLocalScorePrefixRejectsMalformedOrderAndSpan(t *testing.T) {
	for _, candidates := range [][]Paul2013CandidateContinuity{
		{{TotalSpan: 1, InitialOrderingKey: 0}, {TotalSpan: 1, InitialOrderingKey: -1}},
		{{TotalSpan: 0, InitialOrderingKey: -1}},
		{{TotalSpan: 5, InitialOrderingKey: -1}},
	} {
		if _, err := Paul2013CandidateLocalScorePrefix(candidates, 4); err == nil {
			t.Fatalf("accepted malformed candidates: %+v", candidates)
		}
	}
	if _, err := Paul2013CandidateLocalScorePrefix(nil, 0); err == nil {
		t.Fatal("accepted zero phone positions")
	}
}

func TestPaul2013TailOrderingScore(t *testing.T) {
	for _, test := range []struct {
		name      string
		local     float32
		candidate byte
		context   byte
		want      float32
	}{
		{name: "d class multiplier", local: 7, candidate: 'd', want: 35},
		{name: "special class multiplier", local: 7, candidate: 0x0c, context: 1, want: 7 * math.Float32frombits(0x00002000)},
		{name: "matching special class unchanged", local: 7, candidate: 0x0c, context: 0x0c, want: 7},
		{name: "ordinary class unchanged", local: 7, candidate: 'a', context: 'b', want: 7},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := Paul2013TailOrderingScore(test.local, test.candidate, test.context)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("tail ordering score = %g, want %g", got, test.want)
			}
		})
	}
	if _, err := Paul2013TailOrderingScore(float32(math.Inf(1)), 0, 0); err == nil {
		t.Fatal("accepted an infinite local cost")
	}
}

func TestSortPaul2013CandidateTailAndCap(t *testing.T) {
	candidates := make([]Paul2013ScoredCandidate, 33)
	for index := range candidates {
		candidates[index] = Paul2013ScoredCandidate{
			Continuity:          Paul2013CandidateContinuity{Unit: UnitRef{Bank: "gen", Index: uint32(index)}},
			LocalCost:           float32(index),
			SecondaryOrderScore: float32(33 - index),
		}
	}
	// The native sorter leaves the prefix untouched and sorts only from this
	// boundary before truncating the candidate list to 30.
	got, err := SortPaul2013CandidateTailAndCap(candidates, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 30 {
		t.Fatalf("shortlist has %d candidates, want 30", len(got))
	}
	if got[0].Continuity.Unit.Index != 0 {
		t.Fatalf("prefix candidate moved to unit %d", got[0].Continuity.Unit.Index)
	}
	if got[1].Continuity.Unit.Index != 32 || got[29].Continuity.Unit.Index != 4 {
		t.Fatalf("sorted tail bounds = %d..%d, want 32..4", got[1].Continuity.Unit.Index, got[29].Continuity.Unit.Index)
	}
	for index := 2; index < len(got); index++ {
		if got[index-1].SecondaryOrderScore > got[index].SecondaryOrderScore {
			t.Fatalf("tail is not sorted at %d: %g > %g", index, got[index-1].SecondaryOrderScore, got[index].SecondaryOrderScore)
		}
	}
}

func TestSortPaul2013CandidateTailAndCapSkipsShortListsAndValidatesInputs(t *testing.T) {
	short := []Paul2013ScoredCandidate{
		{Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Index: 1}}, SecondaryOrderScore: 2},
		{Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Index: 2}}, SecondaryOrderScore: 1},
	}
	got, err := SortPaul2013CandidateTailAndCap(short, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Continuity.Unit.Index != 1 || got[1].Continuity.Unit.Index != 2 {
		t.Fatalf("short list was reordered: %+v", got)
	}

	equalScores := make([]Paul2013ScoredCandidate, 31)
	for index := range equalScores {
		equalScores[index] = Paul2013ScoredCandidate{
			Continuity:          Paul2013CandidateContinuity{Unit: UnitRef{Index: uint32(index)}},
			SecondaryOrderScore: 1,
		}
	}
	got, err = SortPaul2013CandidateTailAndCap(equalScores, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Continuity.Unit.Index != 0 || len(got) != 30 {
		t.Fatalf("equal-score sort changed the untouched prefix or wrong cap: len=%d first=%d", len(got), got[0].Continuity.Unit.Index)
	}
	seen := make(map[uint32]bool, len(got))
	for _, candidate := range got {
		if seen[candidate.Continuity.Unit.Index] {
			t.Fatalf("equal-score sort duplicated candidate %d", candidate.Continuity.Unit.Index)
		}
		seen[candidate.Continuity.Unit.Index] = true
	}
	if _, err := SortPaul2013CandidateTailAndCap(equalScores, 32); err == nil {
		t.Fatal("accepted a tail start beyond the candidate list")
	}
	equalScores[3].SecondaryOrderScore = float32(math.NaN())
	if _, err := SortPaul2013CandidateTailAndCap(equalScores, 0); err == nil {
		t.Fatal("accepted a non-finite secondary score")
	}
}

func TestSortPaul2013CandidateTailPreservesNativeInsertionTieOrder(t *testing.T) {
	candidates := make([]Paul2013ScoredCandidate, 31)
	for index := range candidates {
		candidates[index] = Paul2013ScoredCandidate{
			Continuity:          Paul2013CandidateContinuity{Unit: UnitRef{Index: uint32(index)}},
			SecondaryOrderScore: 1,
		}
	}
	got, err := SortPaul2013CandidateTailAndCap(candidates, 14)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 30 || got[13].Continuity.Unit.Index != 13 ||
		got[14].Continuity.Unit.Index != 14 || got[29].Continuity.Unit.Index != 29 {
		t.Fatalf("short equal-score tail order = len:%d prefix:%d tail:%d..%d, want 30, 13, 14..29",
			len(got), got[13].Continuity.Unit.Index, got[14].Continuity.Unit.Index, got[29].Continuity.Unit.Index)
	}
}

func TestFinalizePaul2013CandidateShortlistRestoresProtectedTail(t *testing.T) {
	candidates := make([]Paul2013ScoredCandidate, 35)
	for index := range candidates {
		candidates[index] = Paul2013ScoredCandidate{
			Continuity:          Paul2013CandidateContinuity{Unit: UnitRef{Index: uint32(index)}},
			LocalCost:           float32(index),
			SecondaryOrderScore: float32(index),
		}
	}
	candidates[2].NodeFlag = 1
	candidates[32].NodeFlag = 1
	candidates[34].NodeFlag = 1

	got, err := FinalizePaul2013CandidateShortlist(candidates, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Candidates) != 30 || len(got.ProtectedReplacements) != 2 {
		t.Fatalf("shortlist has %d candidates and %d replacements, want 30 and 2", len(got.Candidates), len(got.ProtectedReplacements))
	}
	if got.Candidates[0].Continuity.Unit.Index != 0 {
		t.Fatalf("untouched prefix moved to unit %d", got.Candidates[0].Continuity.Unit.Index)
	}
	if got.ProtectedReplacements[0].Slot != 29 || got.ProtectedReplacements[0].Candidate.Continuity.Unit.Index != 32 ||
		got.ProtectedReplacements[1].Slot != 28 || got.ProtectedReplacements[1].Candidate.Continuity.Unit.Index != 34 {
		t.Fatalf("protected replacements = %+v, want units 32/34 in slots 29/28", got.ProtectedReplacements)
	}
	if got.Candidates[29].Continuity.Unit.Index != 32 || got.Candidates[28].Continuity.Unit.Index != 34 {
		t.Fatalf("shortlist tail replacements = %d/%d, want 32/34", got.Candidates[29].Continuity.Unit.Index, got.Candidates[28].Continuity.Unit.Index)
	}
	if candidates[29].Continuity.Unit.Index != 29 {
		t.Fatal("finalization modified the caller's candidate slice")
	}
}

func TestFinalizePaul2013CandidateShortlistRespectsNativeProtectionLimit(t *testing.T) {
	candidates := make([]Paul2013ScoredCandidate, 40)
	for index := range candidates {
		candidates[index] = Paul2013ScoredCandidate{
			Continuity:          Paul2013CandidateContinuity{Unit: UnitRef{Index: uint32(index)}},
			SecondaryOrderScore: float32(index),
		}
	}
	for _, index := range []int{0, 1, 2, 3, 4, 37, 38, 39} {
		candidates[index].NodeFlag = 1
	}
	got, err := FinalizePaul2013CandidateShortlist(candidates, 5, 8)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ProtectedReplacements) != 0 {
		t.Fatalf("restored protected tail despite five protected prefix candidates: %+v", got.ProtectedReplacements)
	}

	for index := range candidates {
		candidates[index].NodeFlag = 0
	}
	for _, index := range []int{35, 36, 37, 38, 39} {
		candidates[index].NodeFlag = 1
	}
	got, err = FinalizePaul2013CandidateShortlist(candidates, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ProtectedReplacements) != 5 {
		t.Fatalf("restored %d protected candidates, want native maximum 5", len(got.ProtectedReplacements))
	}
}

func TestFinalizePaul2013CandidateShortlistValidatesInputsAndSkipsSmallLists(t *testing.T) {
	small := []Paul2013ScoredCandidate{{Continuity: Paul2013CandidateContinuity{Unit: UnitRef{Index: 7}}}}
	got, err := FinalizePaul2013CandidateShortlist(small, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Candidates) != 1 || got.Candidates[0].Continuity.Unit.Index != 7 {
		t.Fatalf("small shortlist changed: %+v", got)
	}
	if _, err := FinalizePaul2013CandidateShortlist(small, -1, 0); err == nil {
		t.Fatal("accepted a negative tail start")
	}
	if _, err := FinalizePaul2013CandidateShortlist(small, 0, 2); err == nil {
		t.Fatal("accepted a protected count larger than the candidate list")
	}
	large := make([]Paul2013ScoredCandidate, 31)
	large[4].NodeFlag = 2
	if _, err := FinalizePaul2013CandidateShortlist(large, 0, 1); err == nil {
		t.Fatal("accepted a non-boolean native node flag")
	}
}

func TestScorePaul2013CandidateLocalCostScalesFeatureAndPreservesCategory(t *testing.T) {
	table := distance.NewFeatureTable()
	input := UnitScoreInput{
		PrimaryFeature: FeaturePair{Candidate: 4, Context: 2},
		Scale:          1,
	}
	continuity := Paul2013CandidateContinuity{TotalSpan: 4, WeightedSide: 1}
	scored, err := ScorePaul2013CandidateLocalCost(table, input, continuity)
	if err != nil {
		t.Fatal(err)
	}

	category, err := KnownUnitCategoricalPenalty(input.BytePenalty)
	if err != nil {
		t.Fatal(err)
	}
	feature, err := table.Lookup(4, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := category + feature/5*(1.0/3.0)
	if math.Abs(float64(scored-want)) > 0.000001 {
		t.Fatalf("continuity-scaled local cost = %v, want %v", scored, want)
	}
}

func TestScorePaul2013CandidateLocalPrefixConnectsRecordScoringAndProtectionFlags(t *testing.T) {
	model := continuityFixtureModel{
		counts:  map[string]uint32{"gen": 64, "num": 1, "etc": 1, "alp": 1},
		records: map[uint32]dat.UnitRecord{},
	}
	candidates := make([]Paul2013CandidateContinuity, 31)
	contexts := make([]Paul2013UnitScoreContext, len(candidates))
	keys := make([]uint32, len(candidates))
	for index := range candidates {
		candidates[index] = Paul2013CandidateContinuity{
			Unit: UnitRef{Bank: "gen", Index: uint32(index)}, TotalSpan: 1,
			InitialOrderingKey: -1,
		}
		keys[index] = uint32(index + 1)
	}
	candidates[len(candidates)-1].InitialOrderingKey = 0
	wholePositionKeys := []uint32{keys[0], keys[30]}

	got, err := ScorePaul2013CandidateLocalPrefix(
		model, distance.NewFeatureTable(), candidates, 3, contexts, keys, wholePositionKeys,
	)
	if err != nil {
		t.Fatal(err)
	}
	if got.ScoredPrefix != 30 || len(got.Candidates) != len(candidates) {
		t.Fatalf("local-score result has prefix/list %d/%d, want 30/%d", got.ScoredPrefix, len(got.Candidates), len(candidates))
	}
	for index := 0; index < 30; index++ {
		if !got.Candidates[index].HasLocalCost {
			t.Errorf("candidate %d in native score prefix was not scored", index)
		}
	}
	if got.Candidates[0].NodeFlag != 1 || got.Candidates[29].NodeFlag != 0 || got.Candidates[30].NodeFlag != 0 {
		t.Fatalf("native prefix flags = %d/%d/%d, want 1/0/0", got.Candidates[0].NodeFlag, got.Candidates[29].NodeFlag, got.Candidates[30].NodeFlag)
	}
	if got.Candidates[30].HasLocalCost {
		t.Fatal("candidate outside the native local-score prefix was scored")
	}
}

func TestScorePaul2013CandidateLocalPrefixValidatesContextAlignment(t *testing.T) {
	_, err := ScorePaul2013CandidateLocalPrefix(
		continuityFixtureModel{}, distance.NewFeatureTable(),
		[]Paul2013CandidateContinuity{{TotalSpan: 1}}, 1, nil, nil, nil,
	)
	if err == nil {
		t.Fatal("accepted a context list that does not align with candidates")
	}
}
