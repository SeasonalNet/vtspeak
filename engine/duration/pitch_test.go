package duration

import (
	"reflect"
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

func TestEvaluatePaul2013PitchTextTruncatesScalarIntoVectorInput(t *testing.T) {
	sequence := text.LexicalPhoneSequence{
		Phones: []text.CMUPhone{{Label: "B"}, {Label: "AH", Stress: 1, Vowel: true}},
		Tokens: []text.LexicalTokenSpan{{SourceSurface: "bah", Surface: "bah", PhoneStart: 0, PhoneEnd: 2}},
	}
	vectorOutput := make([]int16, 12)
	vectorOutput[0] = 42
	constantScalar := func(value int16) *tree3.Tree {
		return &tree3.Tree{
			OutputWidth: 1,
			Nodes:       []tree3.Node{{Feature: 0, Operation: 'C', Threshold: 32767, WhenTrue: -1, WhenFalse: -1}},
			Outputs:     [][]int16{{value}},
		}
	}
	constantVector := func(values []int16) *tree3.Tree {
		return &tree3.Tree{
			OutputWidth: 12,
			Nodes:       []tree3.Node{{Feature: 0, Operation: 'C', Threshold: 32767, WhenTrue: -1, WhenFalse: -1}},
			Outputs:     [][]int16{values},
		}
	}
	wideVectorTree := &tree3.Tree{
		OutputWidth: 12,
		Nodes: []tree3.Node{{
			Feature: 11, Operation: 'C', Threshold: -1,
			WhenTrue: -1, WhenFalse: -2,
		}},
		Outputs: [][]int16{vectorOutput, make([]int16, 12)},
	}
	catalog := &tree3.Catalog{Pitch: map[string]*tree3.Tree{
		"nbt.tree3": constantScalar(7),
		"nbf.tree3": constantVector(make([]int16, 12)),
		"sbt.tree3": constantScalar(255),
		"sbf.tree3": wideVectorTree,
	}}
	got, err := EvaluatePaul2013PitchText(sequence, catalog)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].Phones) != 2 {
		t.Fatalf("pitch results = %+v, want two phone results", got)
	}
	if got[0].Phones[0].ScalarTree != "nbt.tree3" || got[0].Phones[0].VectorTree != "nbf.tree3" {
		t.Fatalf("first phone pitch pair = %s/%s, want nbt/nbf", got[0].Phones[0].ScalarTree, got[0].Phones[0].VectorTree)
	}
	last := got[0].Phones[1]
	if last.ScalarTree != "sbt.tree3" || last.VectorTree != "sbf.tree3" {
		t.Fatalf("terminal phone pitch pair = %s/%s, want sbt/sbf", last.ScalarTree, last.VectorTree)
	}
	if last.ScalarValue != 255 || last.StoredScalarValue != -1 || last.Input[11] != -1 {
		t.Fatalf("terminal scalar transfer = raw %d stored %d vector input %d, want 255/-1/-1", last.ScalarValue, last.StoredScalarValue, last.Input[11])
	}
	if last.VectorLeaf != 0 || last.VectorTreeValues[0] != 42 || last.VectorValues[0] != 6 {
		t.Fatalf("terminal vector raw/processed result = leaf %d raw %v processed %v, want raw value 42 and smoothed value 6", last.VectorLeaf, last.VectorTreeValues, last.VectorValues)
	}
	if got[0].BoundaryAverages[0] != ([2][3]int16{{6, 6, 6}, {6, 6, 6}}) {
		t.Fatalf("first phone-pair boundary averages = %v, want six rounded values of 6", got[0].BoundaryAverages[0])
	}
	if got[0].Phones[0].VectorValues[9] != 6 || got[0].Phones[1].VectorValues[0] != 6 {
		t.Fatalf("smoothed adjacent vector edges = %d/%d, want 6/6", got[0].Phones[0].VectorValues[9], got[0].Phones[1].VectorValues[0])
	}
}

func TestEvaluatePaul2013PitchTextRejectsMissingTreePairs(t *testing.T) {
	sequence := text.LexicalPhoneSequence{
		Phones: []text.CMUPhone{{Label: "AH", Vowel: true}},
		Tokens: []text.LexicalTokenSpan{{SourceSurface: "a", Surface: "a", PhoneStart: 0, PhoneEnd: 1}},
	}
	if _, err := EvaluatePaul2013PitchText(sequence, &tree3.Catalog{Pitch: map[string]*tree3.Tree{}}); err == nil {
		t.Fatal("pitch evaluation accepted a catalog without pitch trees")
	}
}

func TestEvaluatePaul2013PitchTextWithMarkersSelectsTerminalTreePairs(t *testing.T) {
	sequence := text.LexicalPhoneSequence{
		Phones: []text.CMUPhone{{Label: "B"}, {Label: "D"}, {Label: "F"}, {Label: "G"}},
		Tokens: []text.LexicalTokenSpan{
			{SourceSurface: "one", Surface: "one", PhoneStart: 0, PhoneEnd: 1},
			{SourceSurface: "two", Surface: "two", PhoneStart: 1, PhoneEnd: 2},
			{SourceSurface: "three", Surface: "three", PhoneStart: 2, PhoneEnd: 3},
			{SourceSurface: "four", Surface: "four", PhoneStart: 3, PhoneEnd: 4},
		},
	}
	constantScalar := func(value int16) *tree3.Tree {
		return &tree3.Tree{
			OutputWidth: 1,
			Nodes:       []tree3.Node{{Feature: 0, Operation: 'C', Threshold: 32767, WhenTrue: -1, WhenFalse: -1}},
			Outputs:     [][]int16{{value}},
		}
	}
	constantVector := func() *tree3.Tree {
		return &tree3.Tree{
			OutputWidth: 12,
			Nodes:       []tree3.Node{{Feature: 0, Operation: 'C', Threshold: 32767, WhenTrue: -1, WhenFalse: -1}},
			Outputs:     [][]int16{make([]int16, 12)},
		}
	}
	pitchTrees := map[string]*tree3.Tree{}
	for index, pair := range [][2]string{{"nbt.tree3", "nbf.tree3"}, {"bt.tree3", "bf.tree3"}, {"sbt.tree3", "sbf.tree3"}, {"qbt.tree3", "qbf.tree3"}} {
		pitchTrees[pair[0]] = constantScalar(int16(index))
		pitchTrees[pair[1]] = constantVector()
	}
	catalog := &tree3.Catalog{Pitch: pitchTrees}
	got, err := EvaluatePaul2013PitchTextWithMarkers(sequence, []byte{'^', 'Z', '[', '.'}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	wantScalarTrees := []string{"bt.tree3", "sbt.tree3", "qbt.tree3", "nbt.tree3"}
	wantVectorTrees := []string{"bf.tree3", "sbf.tree3", "qbf.tree3", "nbf.tree3"}
	for tokenIndex, token := range got {
		if len(token.Phones) != 1 {
			t.Fatalf("token %d has %d phone results, want one", tokenIndex, len(token.Phones))
		}
		phone := token.Phones[0]
		if phone.ScalarTree != wantScalarTrees[tokenIndex] || phone.VectorTree != wantVectorTrees[tokenIndex] {
			t.Errorf("token %d marker selected %s/%s, want %s/%s", tokenIndex, phone.ScalarTree, phone.VectorTree, wantScalarTrees[tokenIndex], wantVectorTrees[tokenIndex])
		}
	}
}

func TestEvaluatePaul2013PitchTextWithMarkersRejectsCountMismatch(t *testing.T) {
	sequence := text.LexicalPhoneSequence{
		Phones: []text.CMUPhone{{Label: "B"}},
		Tokens: []text.LexicalTokenSpan{{SourceSurface: "be", Surface: "be", PhoneStart: 0, PhoneEnd: 1}},
	}
	if _, err := EvaluatePaul2013PitchTextWithMarkers(sequence, nil, &tree3.Catalog{}); err == nil {
		t.Fatal("pitch evaluation accepted a terminal-marker count mismatch")
	}
}

func TestBuildPaul2013PitchBoundaryAverages(t *testing.T) {
	var left, right [12]int16
	for index := range left {
		left[index] = int16(index)
		right[index] = int16(index + len(left))
	}
	got := BuildPaul2013PitchBoundaryAverages(left, right)
	want := [2][3]int16{{9, 10, 11}, {12, 13, 14}}
	if got != want {
		t.Fatalf("boundary averages = %v, want %v", got, want)
	}

	for index := range left {
		left[index] = -2
		right[index] = -2
	}
	negative := BuildPaul2013PitchBoundaryAverages(left, right)
	if negative != ([2][3]int16{{-1, -1, -1}, {-1, -1, -1}}) {
		t.Fatalf("negative half-step averages = %v, want signed result -1 from native division", negative)
	}
}

func TestSmoothPaul2013PitchVectorsAppliesEveryAdjacentBoundary(t *testing.T) {
	vectors := make([][12]int16, 3)
	for index := range vectors[1] {
		vectors[1][index] = 100
		vectors[2][index] = 200
	}
	got := SmoothPaul2013PitchVectors(vectors)
	want := [][2][3]int16{
		{{14, 29, 43}, {57, 71, 86}},
		{{114, 129, 143}, {157, 171, 186}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("boundary averages = %v, want %v", got, want)
	}
	if vectors[0][9] != 14 || vectors[0][10] != 29 || vectors[0][11] != 43 {
		t.Fatalf("left vector tail = %v, want [14 29 43]", vectors[0][9:])
	}
	if [3]int16(vectors[1][:3]) != [3]int16{57, 71, 86} || [3]int16(vectors[1][9:]) != [3]int16{114, 129, 143} {
		t.Fatalf("middle vector edges = %v/%v, want [57 71 86]/[114 129 143]", vectors[1][:3], vectors[1][9:])
	}
	if [3]int16(vectors[2][:3]) != [3]int16{157, 171, 186} {
		t.Fatalf("right vector head = %v, want [157 171 186]", vectors[2][:3])
	}
}
