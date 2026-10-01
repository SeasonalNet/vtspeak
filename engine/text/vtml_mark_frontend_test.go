package text

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTextWithPaul2013InlineMarks(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	source := `Hel<vtml_mark name="inside"/>lo <vtml_mark/>Hello.`
	got, err := (LexiconFrontend{Dictionary: dictionary}).ResolveTextWithPaul2013InlineMarks(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tokens) != 2 || got.Tokens[0].Surface != "Hello" || got.Tokens[1].Surface != "Hello" {
		t.Fatalf("resolved tokens = %+v, want two Hello tokens", got.Tokens)
	}
	firstMarkEnd := len(`<vtml_mark name="inside"/>`)
	if got.Tokens[0].SourceByteStart != 0 || got.Tokens[0].SourceByteEnd != firstMarkEnd+4 {
		t.Errorf("first token source span = [%d,%d], want [0,%d]", got.Tokens[0].SourceByteStart, got.Tokens[0].SourceByteEnd, firstMarkEnd+4)
	}
	if len(got.Marks) != 2 || got.Marks[0] != (Paul2013InlineMarkEvent{
		SourceByteOffset: 3, Kind: paul2013VTMLMarkNamedKind, Name: "inside",
	}) || got.Marks[1] != (Paul2013InlineMarkEvent{
		SourceByteOffset: len(`Hel<vtml_mark name="inside"/>lo `), Kind: paul2013VTMLMarkUnnamedKind,
	}) {
		t.Errorf("mark events = %+v", got.Marks)
	}
	sequence, err := SelectLexicalPronunciationsWithMarks(got, []int{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(sequence.InlineMarks) != len(got.Marks) || sequence.InlineMarks[0] != got.Marks[0] || sequence.InlineMarks[1] != got.Marks[1] {
		t.Fatalf("selected phone sequence lost inline marks: %+v", sequence.InlineMarks)
	}
}
