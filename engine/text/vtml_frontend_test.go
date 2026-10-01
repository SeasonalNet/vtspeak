package text

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTextWithPaul2013CapturedVTMLComposesSupportedForms(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	source := `Hel<vtml_mark name="inside"/>lo <vtml_pause time="1"/><vtml_sub alias="Hello.">ignored</vtml_sub>`
	got, err := (LexiconFrontend{Dictionary: dictionary}).ResolveTextWithPaul2013CapturedVTML(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tokens) != 2 || got.Tokens[0].Surface != "Hello" || got.Tokens[1].Surface != "Hello" {
		t.Fatalf("tokens = %+v, want two Hello tokens", got.Tokens)
	}
	markEnd := len(`<vtml_mark name="inside"/>`)
	if got.Tokens[0].SourceByteStart != 0 || got.Tokens[0].SourceByteEnd != markEnd+4 {
		t.Errorf("marked token source span = [%d,%d], want [0,%d]", got.Tokens[0].SourceByteStart, got.Tokens[0].SourceByteEnd, markEnd+4)
	}
	if len(got.Marks) != 1 || got.Marks[0] != (Paul2013InlineMarkEvent{
		SourceByteOffset: 3, Kind: paul2013VTMLMarkNamedKind, Name: "inside",
	}) {
		t.Errorf("marks = %+v", got.Marks)
	}
	if len(got.Pauses) != 1 || got.Pauses[0].AfterToken != 1 || got.Pauses[0].DurationMilliseconds != 1 || got.Pauses[0].OutputFramesAt16KHz != 16 {
		t.Errorf("pauses = %+v", got.Pauses)
	}
	sequence, err := SelectPaul2013CapturedVTMLPronunciations(got, []int{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(sequence.InlinePauses) != 1 || len(sequence.InlineMarks) != 1 {
		t.Fatalf("selected sequence lost VTML events: pauses=%+v marks=%+v", sequence.InlinePauses, sequence.InlineMarks)
	}
}

func TestResolveTextWithPaul2013CapturedVTMLRejectsUnknownAndAmbiguousMarkup(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	frontend := LexiconFrontend{Dictionary: dictionary}
	for _, source := range []string{
		`<vtml_unknown/>Hello`,
		`Hel<vtml_sub alias="lo">ignored</vtml_sub>`,
		`<vtml_pause time="1" extra="x"/>Hello`,
	} {
		if _, err := frontend.ResolveTextWithPaul2013CapturedVTML(context.Background(), source); err == nil {
			t.Errorf("ResolveTextWithPaul2013CapturedVTML(%q) succeeded", source)
		}
	}
}
