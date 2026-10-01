package text

import (
	"reflect"
	"testing"
)

func TestPaul2013PronunciationWordClassTable(t *testing.T) {
	if len(paul2013PronunciationContextWords) != 99 {
		t.Fatalf("word table length = %d, want 99", len(paul2013PronunciationContextWords))
	}
	for index, surface := range paul2013PronunciationContextWords {
		if index > 0 && paul2013PronunciationContextWords[index-1] >= surface {
			t.Fatalf("word table is not strictly sorted at %d: %q then %q", index, paul2013PronunciationContextWords[index-1], surface)
		}
		if got, want := Paul2013PronunciationWordClass(surface), int16(index+2); got != want {
			t.Errorf("class(%q) = %d, want %d", surface, got, want)
		}
	}
	for surface, want := range map[string]int16{"": 0, "not-in-table": 1, "ABOUT": 1} {
		if got := Paul2013PronunciationWordClass(surface); got != want {
			t.Errorf("class(%q) = %d, want %d", surface, got, want)
		}
	}
}

func TestPaul2013PronunciationContinuityFlag(t *testing.T) {
	for surface, want := range map[string]int16{
		"":      0,
		"A":     1,
		"Zebra": 1,
		"about": 0,
		"zebra": 0,
		"'A":    0,
		"é":     0, // The DLL sign-extends the first UTF-8 byte.
	} {
		if got := Paul2013PronunciationContinuityFlag(surface); got != want {
			t.Errorf("continuity flag(%q) = %d, want %d", surface, got, want)
		}
	}
}

func TestPaul2013PronunciationPhoneContextFeature(t *testing.T) {
	tests := []struct {
		name   string
		word   string
		result int16
	}{
		{name: "empty", word: "", result: 0},
		{name: "case insensitive class two", word: "DECISION", result: 2},
		{name: "class two suffix", word: "happiness", result: 2},
		{name: "class three suffix", word: "beautiful", result: 3},
		{name: "class four suffix", word: "carefully", result: 4},
		{name: "class five suffix", word: "morning", result: 5},
		{name: "class six suffix", word: "artist", result: 6},
		{name: "class six plural suffix", word: "artists", result: 6},
		{name: "class seven suffix", word: "walked", result: 7},
		{name: "class one residual", word: "hello", result: 1},
		{name: "class one long residual", word: "ordinary", result: 1},
		{name: "class eight final s", word: "plants", result: 8},
		{name: "class one short final s", word: "cats", result: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Paul2013PronunciationPhoneContextFeature(test.word); got != test.result {
				t.Fatalf("feature(%q) = %d, want %d", test.word, got, test.result)
			}
		})
	}
}

func TestPaul2013PronunciationPhoneContextFeatureSuffixTables(t *testing.T) {
	for _, word := range []string{
		"carelessness", "overbalance", "invention", "reactions", "mansion",
		"transmissions", "kindness", "kindnesses", "payment", "alignments",
		"friendship", "friendships", "ability", "varieties", "greetings",
	} {
		if got := Paul2013PronunciationPhoneContextFeature(word); got != 2 {
			t.Errorf("class-two feature(%q) = %d, want 2", word, got)
		}
	}
	for _, word := range []string{
		"cautious", "Italian", "endless", "expensive", "helpful", "practical",
		"scientific", "initial", "creative", "conditional",
	} {
		if got := Paul2013PronunciationPhoneContextFeature(word); got != 3 {
			t.Errorf("class-three feature(%q) = %d, want 3", word, got)
		}
	}
}

func TestPaul2013KnownPronunciationTripletFeatureFastPaths(t *testing.T) {
	for input, want := range map[[3]string]struct {
		value int16
		known bool
	}{
		{"left", "", "right"}:     {value: 0, known: true},
		{"left", "you", "right"}:  {value: 8, known: true},
		{"let", "'s", "here"}:     {value: 8, known: true},
		{"Let", "'s", "here"}:     {value: 10, known: true},
		{"left", "You", "right"}:  {value: 1, known: true},
		{"left", "your", "right"}: {value: 7, known: true},
		{"left", "wo", "job"}:     {value: 2, known: true},
		{"left", "sha", "job"}:    {value: 2, known: true},
		{"left", "ca", "job"}:     {value: 2, known: true},
		{"left", "WO", "job"}:     {value: 1, known: true},
		{"left", "wo", "work"}:    {value: 1, known: true},
		{"left", "ought", "to"}:   {value: 2, known: true},
		{"left", "have", "to"}:    {value: 2, known: true},
		{"left", "has", "to"}:     {value: 2, known: true},
		{"left", "had", "to"}:     {value: 2, known: true},
		{"left", "need", "to"}:    {value: 2, known: true},
		{"ought", "to", "right"}:  {value: 2, known: true},
		{"have", "to", "right"}:   {value: 2, known: true},
		{"has", "to", "right"}:    {value: 2, known: true},
		{"had", "to", "right"}:    {value: 2, known: true},
		{"need", "to", "right"}:   {value: 2, known: true},
		{"LEFT", "to", "right"}:   {value: 5, known: true},
		{"left", "a7!", "right"}:  {value: 9, known: true},
		{"left", "a!", "right"}:   {value: 10, known: true},
		{"left", "word", "right"}: {value: 1, known: true},
	} {
		gotValue, gotKnown := Paul2013KnownPronunciationTripletFeature(input[0], input[1], input[2])
		if gotValue != want.value || gotKnown != want.known {
			t.Errorf("triplet fast path(%q, %q, %q) = (%d, %t), want (%d, %t)", input[0], input[1], input[2], gotValue, gotKnown, want.value, want.known)
		}
	}
}

func TestPaul2013PronunciationTripletClassTables(t *testing.T) {
	wantCounts := []int{30, 10, 12, 17, 3, 12, 24, 18, 53}
	for groupIndex, group := range paul2013TripletClassGroups {
		if got := len(group.words); got != wantCounts[groupIndex] {
			t.Errorf("triplet class %d table size = %d, want %d", group.class, got, wantCounts[groupIndex])
		}
		for i, word := range group.words {
			if i > 0 && group.words[i-1] >= word {
				t.Errorf("triplet class %d table is not strictly sorted at %d: %q then %q", group.class, i, group.words[i-1], word)
			}
			got, known := Paul2013KnownPronunciationTripletFeature("left", word, "right")
			want := group.class
			for _, earlier := range paul2013TripletClassGroups {
				if earlier.class == group.class {
					break
				}
				if paul2013ContainsWord(earlier.words, word) {
					want = earlier.class
					break
				}
			}
			if !known || got != want {
				t.Errorf("triplet class(%q) = (%d, %t), want (%d, true) by table precedence", word, got, known, want)
			}
		}
	}
	if got, known := Paul2013KnownPronunciationTripletFeature("left", "about", "right"); !known || got != 5 {
		t.Errorf("99-word table class = (%d, %t), want (5, true)", got, known)
	}
	if got := len(paul2013TripletApostropheSLeftClass3Words); got != 10 {
		t.Fatalf("apostrophe-s class-3 left table size = %d, want 10", got)
	}
	for index, left := range paul2013TripletApostropheSLeftClass3Words {
		if index > 0 && paul2013TripletApostropheSLeftClass3Words[index-1] >= left {
			t.Errorf("apostrophe-s left table is not strictly sorted at %d", index)
		}
		if got, known := Paul2013KnownPronunciationTripletFeature(left, "'s", "right"); !known || got != 3 {
			t.Errorf("apostrophe-s left context %q = (%d, %t), want (3, true)", left, got, known)
		}
	}
	for _, left := range []string{"'em", "how", "whoever"} {
		if got, known := Paul2013KnownPronunciationTripletFeature(left, "'s", "right"); !known || got != 3 {
			t.Errorf("apostrophe-s class-8/11 left context %q = (%d, %t), want (3, true)", left, got, known)
		}
	}
	if got, known := Paul2013KnownPronunciationTripletFeature("let", "'s", "right"); !known || got != 8 {
		t.Errorf("let/'s context = (%d, %t), want (8, true)", got, known)
	}
	if got, known := Paul2013KnownPronunciationTripletFeature("unknown", "'s", "right"); !known || got != 10 {
		t.Errorf("apostrophe-s fallback = (%d, %t), want (10, true)", got, known)
	}
}

func TestBuildPaul2013PronunciationWordFeatures(t *testing.T) {
	tokens := []LexicalToken{{Surface: "about"}, {Surface: "Hello"}, {Surface: "without"}}
	got, err := BuildPaul2013PronunciationWordFeatures(tokens, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := [4]int16{0, 2, 100, 0}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("word window = %v, want %v", got, want)
	}
	if _, err := BuildPaul2013PronunciationWordFeatures(tokens, 3); err == nil {
		t.Fatal("out-of-range center accepted")
	}
	if _, err := BuildPaul2013PronunciationWordFeatures(nil, 0); err == nil {
		t.Fatal("empty token sequence accepted")
	}
}

func TestBuildPaul2013KnownPronunciationFeatures(t *testing.T) {
	tokens := []LexicalToken{{Surface: "about"}, {Surface: "Hello"}, {Surface: "without"}}
	got, err := BuildPaul2013KnownPronunciationFeatures(tokens, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Values[0] != 0 || got.Values[1] != 2 || got.Values[2] != 100 || got.Values[3] != 0 || got.Values[14] != 1 {
		t.Fatalf("known features = %v, want positions 0-3 [0 2 100 0] and position 14 1", got.Values)
	}
	if got.Complete() {
		t.Fatal("surface-only row reported the unavailable path feature as complete")
	}
	wantMissing := []int{13}
	if !reflect.DeepEqual(got.MissingPositions(), wantMissing) {
		t.Fatalf("missing positions = %v, want %v", got.MissingPositions(), wantMissing)
	}
}

func TestBuildPaul2013PronunciationPathFeatures(t *testing.T) {
	tokens := []LexicalToken{{Surface: "left"}, {Surface: "Hello"}, {Surface: "right"}}
	features, present, err := BuildPaul2013PronunciationPathFeatures(tokens, 1, []byte{13, 14})
	if err != nil {
		t.Fatal(err)
	}
	if !present {
		t.Fatal("nonempty path group was skipped")
	}
	if features.Values[13] != 1 || features.Available&(1<<13) == 0 {
		t.Fatalf("path feature = %d with mask %#x, want 1 and feature 13 available", features.Values[13], features.Available)
	}
	if got, want := features.MissingPositions(), []int{}; !reflect.DeepEqual(got, want) {
		t.Fatalf("missing positions = %v, want %v", got, want)
	}

	empty, present, err := BuildPaul2013PronunciationPathFeatures(tokens, 1, nil)
	if err != nil || present || empty.Available != 0x5fff {
		t.Fatalf("empty path group = (%+v, %t, %v), want surface-only row, false, nil", empty, present, err)
	}
	if _, _, err := BuildPaul2013PronunciationPathFeatures(tokens, 1, []byte{0xff}); err == nil {
		t.Fatal("out-of-range path code accepted")
	}
}

func TestBuildPaul2013PronunciationPathCodeFeatures(t *testing.T) {
	tokens := []LexicalToken{{Surface: "word"}}
	first, present, err := BuildPaul2013PronunciationPathCodeFeatures(tokens, 0, 13)
	if err != nil || !present || first.Values[13] != 1 {
		t.Fatalf("first path-code features = (%+v, %t, %v), want class 1", first, present, err)
	}
	second, present, err := BuildPaul2013PronunciationPathCodeFeatures(tokens, 0, 14)
	if err != nil || !present || second.Values[13] != 2 {
		t.Fatalf("second path-code features = (%+v, %t, %v), want class 2", second, present, err)
	}
	if _, _, err := BuildPaul2013PronunciationPathCodeFeatures(tokens, 0, 0xff); err == nil {
		t.Fatal("out-of-range path code accepted")
	}
}

func TestBuildPaul2013KnownPronunciationFeaturesMarksEmptyTripletMiddles(t *testing.T) {
	tokens := []LexicalToken{{Surface: "left"}, {Surface: ""}, {Surface: "center"}, {Surface: ""}, {Surface: "right"}}
	got, err := BuildPaul2013KnownPronunciationFeatures(tokens, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, position := range []int{5, 6} {
		if got.Values[position] != 0 || got.Available&(1<<position) == 0 {
			t.Errorf("triplet position %d = %d with mask %#x, want known zero", position, got.Values[position], got.Available)
		}
	}
	for _, position := range []int{9, 11} {
		if got.Values[position] != 0 || got.Available&(1<<position) == 0 {
			t.Errorf("empty-word feature %d = %d with mask %#x, want known zero", position, got.Values[position], got.Available)
		}
	}
	for _, position := range []int{4, 7} {
		if got.Values[position] != 1 || got.Available&(1<<position) == 0 {
			t.Errorf("nonempty triplet position %d = %d with mask %#x, want known default class 1", position, got.Values[position], got.Available)
		}
	}
	for position, want := range map[int]int16{8: 1, 9: 0, 10: 1, 11: 0, 12: 1} {
		if got.Values[position] != want || got.Available&(1<<position) == 0 {
			t.Errorf("phone context feature %d = %d with mask %#x, want known value %d", position, got.Values[position], got.Available, want)
		}
	}
}

func TestBuildPaul2013KnownPronunciationFeaturesMarksYouTriplets(t *testing.T) {
	tokens := []LexicalToken{{Surface: "left"}, {Surface: "you"}, {Surface: "center"}, {Surface: "you"}, {Surface: "right"}}
	got, err := BuildPaul2013KnownPronunciationFeatures(tokens, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, position := range []int{5, 6} {
		if got.Values[position] != 8 || got.Available&(1<<position) == 0 {
			t.Errorf("triplet position %d = %d with mask %#x, want known class 8", position, got.Values[position], got.Available)
		}
	}
	for _, position := range []int{4, 7} {
		if got.Values[position] != 1 || got.Available&(1<<position) == 0 {
			t.Errorf("nonempty triplet position %d = %d with mask %#x, want known class 1", position, got.Values[position], got.Available)
		}
	}
}

func TestBuildPaul2013KnownPronunciationFeaturesMarksLetsTriplet(t *testing.T) {
	tokens := []LexicalToken{{Surface: "let"}, {Surface: "'s"}, {Surface: "here"}, {Surface: "center"}}
	got, err := BuildPaul2013KnownPronunciationFeatures(tokens, 3)
	if err != nil {
		t.Fatal(err)
	}
	if got.Values[4] != 8 || got.Available&(1<<4) == 0 {
		t.Fatalf("let/'s triplet position 4 = %d with mask %#x, want known class 8", got.Values[4], got.Available)
	}
}

func TestBuildPaul2013KnownPronunciationFeaturesUsesTripletNeighbors(t *testing.T) {
	for _, test := range []struct {
		name   string
		tokens []LexicalToken
		want   int16
	}{
		{
			name:   "right context",
			tokens: []LexicalToken{{Surface: "left"}, {Surface: "wo"}, {Surface: "job"}, {Surface: "center"}},
			want:   2,
		},
		{
			name:   "left context",
			tokens: []LexicalToken{{Surface: "ought"}, {Surface: "to"}, {Surface: "right"}, {Surface: "center"}},
			want:   2,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := BuildPaul2013KnownPronunciationFeatures(test.tokens, 3)
			if err != nil {
				t.Fatal(err)
			}
			if got.Values[4] != test.want || got.Available&(1<<4) == 0 {
				t.Fatalf("triplet position 4 = %d with mask %#x, want known class %d", got.Values[4], got.Available, test.want)
			}
		})
	}
}
