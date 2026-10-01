package duration

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

func TestEvaluatePronunciationFeaturesUsesAllFifteenValues(t *testing.T) {
	engine := &Engine{pronunciationTree: &tree3.Tree{
		OutputWidth: 1,
		Nodes: []tree3.Node{{
			Feature: 14, Operation: 'C', Threshold: 0,
			WhenTrue: -1, WhenFalse: -2,
		}},
		Outputs: [][]int16{{13}, {27}},
	}}

	features := [15]int16{}
	features[14] = 1
	got, err := engine.EvaluatePronunciationFeatures(features)
	if err != nil {
		t.Fatal(err)
	}
	if got != 27 {
		t.Fatalf("classifier output = %d, want 27", got)
	}
}

func TestChoosePaul2013PronunciationAlternativesEvaluatesEveryPathCode(t *testing.T) {
	features := func(class int16) text.Paul2013PronunciationFeatures {
		row := text.Paul2013PronunciationFeatures{Available: 0x7fff}
		row.Values[13] = class
		return row
	}
	analysis := text.LexicalAnalysis{
		Tokens: []text.LexicalToken{{
			SourceSurface: "read",
			Alternatives: []text.LabeledAlternative{
				{PathGroups: [][]byte{{13, 14}}, Phones: []text.CMUPhone{{Label: "R"}, {Label: "IY", Stress: 1, Vowel: true}}},
				{PathGroups: [][]byte{{15}}, Phones: []text.CMUPhone{{Label: "R"}, {Label: "EH", Stress: 1, Vowel: true}}},
			},
		}},
		PronunciationCandidates: [][]text.Paul2013PronunciationCandidate{{
			{AlternativeIndex: 0, PathGroupIndex: 0, PathCodeIndex: 0, Features: features(1)},
			{AlternativeIndex: 0, PathGroupIndex: 0, PathCodeIndex: 1, Features: features(2)},
			{AlternativeIndex: 1, PathGroupIndex: 0, PathCodeIndex: 0, Features: features(3)},
		}},
	}
	evaluations := 0
	choices, err := choosePaul2013PronunciationAlternatives(analysis, func([15]int16) (int16, error) {
		evaluations++
		return 3, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if evaluations != 3 {
		t.Fatalf("classifier evaluations = %d, want one per path code", evaluations)
	}
	if len(choices) != 1 || choices[0] != 1 {
		t.Fatalf("pronunciation choices = %v, want alternative 1", choices)
	}
	sequence, err := text.SelectLexicalPronunciations(analysis.Tokens, choices)
	if err != nil {
		t.Fatal(err)
	}
	if got := sequence.Tokens[0].AlternativeIndex; got != 1 {
		t.Fatalf("selected alternative = %d, want 1", got)
	}
}

func TestChoosePaul2013PronunciationAlternativesRejectsUnknownGroupMapping(t *testing.T) {
	analysis := text.LexicalAnalysis{
		Tokens: []text.LexicalToken{{
			SourceSurface: "ambiguous",
			Alternatives: []text.LabeledAlternative{
				{PathGroups: [][]byte{{13}, {14}}},
				{PathGroups: [][]byte{{15}}},
			},
		}},
		PronunciationCandidates: [][]text.Paul2013PronunciationCandidate{{}},
	}
	_, err := choosePaul2013PronunciationAlternatives(analysis, func([15]int16) (int16, error) {
		return 0, nil
	})
	if err == nil {
		t.Fatal("multiple path groups in one alternative accepted")
	}
}

func TestChoosePaul2013PronunciationAlternativesPropagatesClassifierFailure(t *testing.T) {
	features := text.Paul2013PronunciationFeatures{Available: 0x7fff}
	analysis := text.LexicalAnalysis{
		Tokens: []text.LexicalToken{{
			SourceSurface: "ambiguous",
			Alternatives: []text.LabeledAlternative{
				{PathGroups: [][]byte{{13}}},
				{PathGroups: [][]byte{{14}}},
			},
		}},
		PronunciationCandidates: [][]text.Paul2013PronunciationCandidate{{
			{AlternativeIndex: 0, PathGroupIndex: 0, Features: features},
		}},
	}
	wantErr := errors.New("test classifier error")
	_, err := choosePaul2013PronunciationAlternatives(analysis, func([15]int16) (int16, error) {
		return 0, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("classifier error = %v, want %v", err, wantErr)
	}
}

func TestEngineResolvesAndEvaluatesLocalAmbiguousWord(t *testing.T) {
	voiceRoot := filepath.Join("..", "..", "data-paul", "M16")
	dictionaryRoot := filepath.Join("..", "..", "data-common", "dict-eng")
	for _, path := range []string{
		filepath.Join(voiceRoot, "ttsdata", "tree3", "duration", "vshort.tree3"),
		filepath.Join(dictionaryRoot, "engttsdict_emb"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Skipf("local VoiceText test resource is unavailable: %s", path)
		}
	}
	engine, err := OpenPaul2013(voiceRoot, dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	source := "read record bow minute bass lead"
	sequence, err := engine.ResolvePronunciationSequence(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(sequence.Tokens) != 6 {
		t.Fatalf("selected pronunciation sequence has %d tokens, want 6", len(sequence.Tokens))
	}
	results, err := engine.Evaluate(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 6 {
		t.Fatalf("duration results = %d, want 6 tokens", len(results))
	}
	for tokenIndex, result := range results {
		if len(result.Phones) == 0 {
			t.Errorf("duration results for token %d (%q) are empty", tokenIndex, result.Token.Surface)
			continue
		}
		if len(result.PitchBoundaryAverages) != len(result.Phones)-1 {
			t.Errorf("token %d has %d pitch boundary averages for %d phones, want %d", tokenIndex, len(result.PitchBoundaryAverages), len(result.Phones), len(result.Phones)-1)
		}
		for phoneIndex, phone := range result.Phones {
			if phone.Pitch == nil {
				t.Errorf("duration result for token %d phone %d has no pitch result", tokenIndex, phoneIndex)
				continue
			}
			if len(phone.Pitch.VectorValues) != 12 || phone.Pitch.Input[11] != int16(phone.Pitch.StoredScalarValue) {
				t.Errorf("pitch result for token %d phone %d has an invalid vector or scalar handoff: %+v", tokenIndex, phoneIndex, phone.Pitch)
			}
		}
	}
	lastPhone := results[len(results)-1].Phones[len(results[len(results)-1].Phones)-1]
	if lastPhone.Pitch == nil || lastPhone.Pitch.ScalarTree != "sbt.tree3" || lastPhone.Pitch.VectorTree != "sbf.tree3" {
		t.Errorf("final phone pitch trees = %+v, want standard-Z sbt/sbf pair", lastPhone.Pitch)
	}
}

func TestEngineEvaluatesComposedCapturedVTML(t *testing.T) {
	voiceRoot := filepath.Join("..", "..", "data-paul", "M16")
	dictionaryRoot := filepath.Join("..", "..", "data-common", "dict-eng")
	for _, path := range []string{
		filepath.Join(voiceRoot, "ttsdata", "tree3", "duration", "vshort.tree3"),
		filepath.Join(dictionaryRoot, "engttsdict_emb"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Skipf("local VoiceText test resource is unavailable: %s", path)
		}
	}
	engine, err := OpenPaul2013(voiceRoot, dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	source := `Hel<vtml_mark name="inside"/>lo <vtml_pause time="1"/><vtml_sub alias="Hello.">ignored</vtml_sub>`
	got, err := engine.EvaluateWithCapturedVTML(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tokens) != 2 || len(got.Sequence.Tokens) != 2 {
		t.Fatalf("evaluation has %d token results and %d sequence spans, want 2", len(got.Tokens), len(got.Sequence.Tokens))
	}
	if len(got.Pauses) != 1 || got.Pauses[0].AfterToken != 1 || got.Pauses[0].OutputFramesAt16KHz != 16 {
		t.Errorf("pause events = %+v", got.Pauses)
	}
	if len(got.Marks) != 1 || got.Marks[0].Name != "inside" || len(got.Sequence.InlineMarks) != 1 {
		t.Errorf("mark events = %+v; sequence marks = %+v", got.Marks, got.Sequence.InlineMarks)
	}
	if got.Tokens[0].Token.Surface != "Hello" || got.Tokens[1].Token.Surface != "Hello" {
		t.Errorf("evaluated surfaces = %q, %q", got.Tokens[0].Token.Surface, got.Tokens[1].Token.Surface)
	}
}

func TestEngineEvaluatesCapturedPhoneSequenceWithExplicitPositionStates(t *testing.T) {
	voiceRoot := filepath.Join("..", "..", "data-paul", "M16")
	dictionaryRoot := filepath.Join("..", "..", "data-common", "dict-eng")
	for _, path := range []string{
		filepath.Join(voiceRoot, "ttsdata", "tree3", "duration", "vshort.tree3"),
		filepath.Join(dictionaryRoot, "engttsdict_emb"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Skipf("local VoiceText test resource is unavailable: %s", path)
		}
	}
	engine, err := OpenPaul2013(voiceRoot, dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	results, err := engine.EvaluateVTMLCMUPhoneme(
		`<vtml_phoneme alphabet="x-cmu" ph="P AH0">pah</vtml_phoneme>.`,
		[]byte{'0', '0'}, []byte{'Z'}, []uint8{1, 1},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := [][9]int16{
		{27, 40, 3, 0, 1, 1, 1, 1, 1},
		{3, 27, 40, 0, 1, 1, 1, 1, 2},
	}
	if len(results) != 1 || len(results[0].Phones) != len(want) {
		t.Fatalf("evaluated results = %+v, want one token with two phones", results)
	}
	for phoneIndex := range want {
		got := results[0].Phones[phoneIndex]
		if got.Input != want[phoneIndex] {
			t.Errorf("phone %d duration tree input = %v, want captured vector %v", phoneIndex, got.Input, want[phoneIndex])
		}
		if got.Pitch == nil || got.Pitch.Input[8] != 1 {
			t.Errorf("phone %d pitch result = %+v, want position state 1", phoneIndex, got.Pitch)
		}
	}
}

func TestEngineLoadsAndEvaluatesSharedATMTTree(t *testing.T) {
	voiceRoot := filepath.Join("..", "..", "data-paul", "M16")
	dictionaryRoot := filepath.Join("..", "..", "data-common", "dict-eng")
	for _, path := range []string{
		filepath.Join(voiceRoot, "ttsdata", "tree3", "duration", "vshort.tree3"),
		filepath.Join(dictionaryRoot, "engttsdict_emb"),
		filepath.Join(dictionaryRoot, "atmt.tree3"),
		filepath.Join(dictionaryRoot, "tppdict_eng"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Skipf("local VoiceText test resource is unavailable: %s", path)
		}
	}
	engine, err := OpenPaul2013(voiceRoot, dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	leaf, output, err := engine.EvaluateATMTTree(context.Background(), 0, make([]int16, 256))
	if err != nil {
		t.Fatal(err)
	}
	if leaf < 0 || len(output) == 0 {
		t.Fatalf("ATMT evaluation returned leaf %d and output %v", leaf, output)
	}
	if _, _, err := engine.EvaluateATMTTree(context.Background(), 27, make([]int16, 256)); err == nil {
		t.Fatal("ATMT tree index 27 was accepted")
	}
	for _, test := range []struct {
		surface string
		tag     byte
		want    string
	}{
		{surface: "ABOUT-SHIPPING", tag: 'F', want: "120"},
		{surface: "ACCORD", tag: 'G', want: "95"},
	} {
		got, found, err := engine.LookupTypedTextCode([]byte(test.surface), test.tag)
		if err != nil || !found || string(got) != test.want {
			t.Errorf("LookupTypedTextCode(%q, %c) = %q, %t, %v; want %q, true, nil", test.surface, test.tag, got, found, err, test.want)
		}
	}
}

func TestEvaluatePronunciationFeaturesRequiresLoadedTree(t *testing.T) {
	var engine *Engine
	if _, err := engine.EvaluatePronunciationFeatures([15]int16{}); err == nil {
		t.Fatal("nil engine accepted pronunciation features")
	}
	if _, err := (&Engine{}).EvaluatePronunciationFeatures([15]int16{}); err == nil {
		t.Fatal("engine without classifier accepted pronunciation features")
	}
}
