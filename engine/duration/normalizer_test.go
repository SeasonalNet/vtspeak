package duration

import (
	"context"
	"reflect"
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

func TestEvaluatePaul2013NormalizerCharacterUsesLoadedTreeSlot(t *testing.T) {
	trees := normalizerTestTrees()
	trees[0].Outputs = [][]int16{{23}}
	engine := &Engine{atmtTrees: trees}
	got, err := engine.EvaluatePaul2013NormalizerCharacter(context.Background(), []byte("A"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := NormalizerCharacterResult{
		Character: 'A', TreeIndex: 0, LeafOrdinal: 0,
		Features: text.Paul2013NormalizerFeatureWindow{0, 0, 0, 1}, Value: 23, CategoryByte: 23,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizer character result = %+v, want %+v", got, want)
	}
}

func TestEvaluatePaul2013NormalizerCharacterSequenceCarriesReverseCategories(t *testing.T) {
	trees := normalizerTestTrees()
	trees[4].Outputs = [][]int16{{1}}
	engine := &Engine{atmtTrees: trees}
	got, err := engine.EvaluatePaul2013NormalizerCharacterSequence(context.Background(), []byte("beb"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("normalizer results = %d, want 3", len(got))
	}
	if got[1].Character != 'E' || got[1].Value != 1 || got[1].CategoryByte != 0x1c {
		t.Fatalf("middle E normalizer result = %+v, want its special category byte 0x1c", got[1])
	}
	if got[1].Features[7] != 2 {
		t.Fatalf("middle E trailing category = %d, want right neighbor category 2", got[1].Features[7])
	}
	if got[0].Character != 'B' || got[0].Features[7] != 0x1c || got[0].Features[8] != 2 {
		t.Fatalf("leftmost B normalizer result = %+v, want reversed E/B categories", got[0])
	}
	if got[2].Character != 'B' || got[2].CategoryByte != 2 {
		t.Fatalf("rightmost normalizer result = %+v", got[2])
	}
}

func TestEvaluatePaul2013NormalizerCharacterSequenceRejectsUnsupportedTokenBytes(t *testing.T) {
	engine := &Engine{atmtTrees: normalizerTestTrees()}
	if _, err := engine.EvaluatePaul2013NormalizerCharacterSequence(context.Background(), []byte("street-1")); err == nil {
		t.Fatal("unsupported token punctuation was accepted as a letter/apostrophe sequence")
	}
}

func TestClassifyPaul2013NormalizerTokenComposesClassRemapping(t *testing.T) {
	trees := normalizerTestTrees()
	trees[0].Outputs = [][]int16{{23}}
	engine := &Engine{atmtTrees: trees}
	got, err := engine.ClassifyPaul2013NormalizerToken(context.Background(), []byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Characters) != 1 || got.Characters[0].Character != 'A' {
		t.Fatalf("normalizer character trace = %+v", got.Characters)
	}
	if !got.Produced || !reflect.DeepEqual(got.ClassCodes, []byte{18}) {
		t.Fatalf("normalizer class output = (% x, %t), want (12, true)", got.ClassCodes, got.Produced)
	}
}

func TestClassifyPaul2013NormalizerTokenReportsNativeEmptyClassOutput(t *testing.T) {
	trees := normalizerTestTrees()
	trees[0].Outputs = [][]int16{{1}}
	engine := &Engine{atmtTrees: trees}
	got, err := engine.ClassifyPaul2013NormalizerToken(context.Background(), []byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Produced || len(got.ClassCodes) != 0 {
		t.Fatalf("empty normalizer class output = (% x, %t), want (empty, false)", got.ClassCodes, got.Produced)
	}
}

func TestIsPaul2013NormalizerEligibleMatchesNativeConsonantRuns(t *testing.T) {
	engine := normalizerEligibilityTestEngine()
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{name: "single consonants around vowel", source: "CAT", want: true},
		{name: "internal two-consonant class", source: "STAA", want: true},
		{name: "apostrophe selects alternate CHC mask", source: "S'TA", want: true},
		{name: "final consonant run", source: "ASTR", want: true},
		{name: "all consonants", source: "STR", want: false},
		{name: "lowercase vowel", source: "a", want: true},
		{name: "unsupported punctuation", source: "CAT-", want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := engine.IsPaul2013NormalizerEligible([]byte(test.source))
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("normalizer eligibility = %t, want %t", got, test.want)
			}
		})
	}
}

func TestIsPaul2013NormalizerEligibleUsesCStringAndBoundedNativeBuffer(t *testing.T) {
	engine := normalizerEligibilityTestEngine()
	if got, err := engine.IsPaul2013NormalizerEligible([]byte("CAT\x00-")); err != nil || !got {
		t.Fatalf("NUL-terminated eligibility = (%t, %v), want (true, nil)", got, err)
	}
	longToken := make([]byte, 32)
	for index := range longToken {
		longToken[index] = 'B'
	}
	if _, err := engine.IsPaul2013NormalizerEligible(longToken); err == nil {
		t.Fatal("token exceeding native local run buffer was not rejected")
	}
}

func normalizerTestTrees() []*tree3.Tree {
	trees := make([]*tree3.Tree, 27)
	for index := range trees {
		trees[index] = &tree3.Tree{
			Nodes:   []tree3.Node{{Feature: 0, Operation: 'L', Threshold: 32767, WhenTrue: -1, WhenFalse: -1}},
			Outputs: [][]int16{{int16(index + 1)}},
		}
	}
	return trees
}

func normalizerEligibilityTestEngine() *Engine {
	return &Engine{txt2: text.Tables{
		"chc_sort.txt2": {
			Name: "chc_sort.txt2",
			Mode: text.TwoColumns,
			Rows: []text.Row{
				{Key: "ST", Value: "0110"},
				{Key: "STR", Value: "0001"},
			},
		},
	}}
}
