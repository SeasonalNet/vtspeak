package duration

import (
	"context"
	"reflect"
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

func TestNormalizePaul2013GenericTokenDispatchesSingleVowelToContextEncoder(t *testing.T) {
	engine := genericNormalizerTestEngine(nil)
	got, err := engine.NormalizePaul2013GenericToken(context.Background(), []byte("A"))
	if err != nil {
		t.Fatal(err)
	}
	wantCodes, err := text.EncodePaul2013ContextString([]byte("A"))
	if err != nil {
		t.Fatal(err)
	}
	if !got.Eligible || !got.UsedContextFallback || got.RowFlagMask != 0x20 ||
		!reflect.DeepEqual(got.ContextCodes, wantCodes) || len(got.ClassCodes) != 0 {
		t.Fatalf("generic normalization of single vowel = %#v, want direct context-code fallback %#v", got, wantCodes)
	}
}

func TestNormalizePaul2013GenericTokenUsesClassDispatchForTwoByteToken(t *testing.T) {
	engine := genericNormalizerTestEngine(classifierTestTrees(2))
	got, err := engine.NormalizePaul2013GenericToken(context.Background(), []byte("be"))
	if err != nil {
		t.Fatal(err)
	}
	if !got.Eligible || got.UsedContextFallback || got.RowFlagMask != 0x10 ||
		!reflect.DeepEqual(got.ClassCodes, []byte{1, 1}) || len(got.ContextCodes) != 0 {
		t.Fatalf("generic normalization of two-byte token = %#v, want class codes 01 01 and flag 0x10", got)
	}
}

func TestNormalizePaul2013GenericTokenExpandsMCPrefixBeforeClassification(t *testing.T) {
	engine := genericNormalizerTestEngine(classifierTestTrees(2))
	engine.txt2["chc_sort.txt2"] = &text.Table{
		Name: "chc_sort.txt2",
		Mode: text.TwoColumns,
		Rows: []text.Row{{Key: "MC", Value: "0100"}},
	}
	got, err := engine.NormalizePaul2013GenericToken(context.Background(), []byte("MCA"))
	if err != nil {
		t.Fatal(err)
	}
	if !got.Eligible || got.UsedContextFallback || got.RowFlagMask != 0x10 ||
		!reflect.DeepEqual(got.ClassCodes, []byte{1, 1, 1, 1}) || len(got.ContextCodes) != 0 {
		t.Fatalf("generic normalization of MC-prefixed token = %#v, want four classified characters from MACA", got)
	}
}

func TestNormalizePaul2013ContextTokenUsesDictionaryBeforeGenericFallback(t *testing.T) {
	engine := genericNormalizerTestEngine(nil)
	lookupCalls := 0
	got, err := engine.NormalizePaul2013ContextToken(context.Background(), []byte("A\x00ignored"),
		func(_ context.Context, source []byte) (Paul2013ContextLookupResult, error) {
			lookupCalls++
			if string(source) != "A" {
				t.Fatalf("lookup source = %q, want C-string prefix A", source)
			}
			return Paul2013ContextLookupResult{
				RecordCount: 1, GenericFallbackGate: true, DictionaryText: []byte("dict\x00tail"),
			}, nil
		})
	if err != nil {
		t.Fatal(err)
	}
	if lookupCalls != 1 || !got.Matched || !got.UsedDictionary ||
		!reflect.DeepEqual(got.DictionaryText, []byte("dict")) || got.Generic.Eligible {
		t.Fatalf("positive dictionary result = %#v (calls %d), want copied dictionary text and no generic call", got, lookupCalls)
	}
}

func TestNormalizePaul2013ContextTokenHonorsGenericGate(t *testing.T) {
	engine := genericNormalizerTestEngine(nil)
	tests := []struct {
		name string
		gate bool
		want bool
	}{
		{name: "closed gate", want: false},
		{name: "open gate", gate: true, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := engine.NormalizePaul2013ContextToken(context.Background(), []byte("A"),
				func(context.Context, []byte) (Paul2013ContextLookupResult, error) {
					return Paul2013ContextLookupResult{GenericFallbackGate: test.gate}, nil
				})
			if err != nil {
				t.Fatal(err)
			}
			if got.Matched != test.want {
				t.Fatalf("matched = %t, want %t", got.Matched, test.want)
			}
			if test.want && (!got.Generic.Eligible || !got.Generic.UsedContextFallback || got.Generic.RowFlagMask != 0x20) {
				t.Fatalf("generic branch = %#v, want eligible context fallback", got.Generic)
			}
		})
	}
}

func genericNormalizerTestEngine(trees []*tree3.Tree) *Engine {
	if trees == nil {
		trees = make([]*tree3.Tree, 27)
	}
	for index := range trees {
		if trees[index] == nil {
			trees[index] = &tree3.Tree{
				OutputWidth: 1,
				Nodes:       []tree3.Node{{Feature: 0, Operation: 'C', Threshold: 0, WhenTrue: -1, WhenFalse: -1}},
				Outputs:     [][]int16{{2}},
			}
		}
	}
	return &Engine{
		txt2: text.Tables{
			"chc_sort.txt2": {
				Name: "chc_sort.txt2",
				Mode: text.TwoColumns,
			},
		},
		atmtTrees: trees,
	}
}

func classifierTestTrees(value int16) []*tree3.Tree {
	trees := make([]*tree3.Tree, 27)
	for index := range trees {
		trees[index] = &tree3.Tree{
			OutputWidth: 1,
			Nodes:       []tree3.Node{{Feature: 0, Operation: 'C', Threshold: 0, WhenTrue: -1, WhenFalse: -1}},
			Outputs:     [][]int16{{value}},
		}
	}
	return trees
}
