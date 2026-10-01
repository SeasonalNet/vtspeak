package text

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPaul2013CapturedVTMLSubstitutions(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{name: "plain text", source: "Hello world.", want: "Hello world."},
		{
			name:   "captured alias form",
			source: `before <vtml_sub alias="Hello world.">x</vtml_sub> after`,
			want:   "before Hello world. after",
		},
		{
			name:   "multiple aliases",
			source: `<vtml_sub alias="Hello">x</vtml_sub>, <vtml_sub alias="world">y</vtml_sub>!`,
			want:   "Hello, world!",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ExpandPaul2013CapturedVTMLSubstitutions(test.source)
			if err != nil {
				t.Fatalf("ExpandPaul2013CapturedVTMLSubstitutions() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("ExpandPaul2013CapturedVTMLSubstitutions() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestExpandPaul2013CapturedVTMLSubstitutionsRejectsUnsupportedForms(t *testing.T) {
	for _, source := range []string{
		`<vtml_sub>x</vtml_sub>`,
		`<vtml_sub alias="Hello">x`,
		`<vtml_sub alias="Hello"><b>x</b></vtml_sub>`,
		`<vtml_pause time="200"/>`,
		`<vtml_sub alias="a" extra="b">x</vtml_sub>`,
		`<vtml_sub alias="">x</vtml_sub>`,
		"<vtml_sub alias=\"caf\u00e9\">x</vtml_sub>",
		`<vtml_sub alias="a&amp;b">x</vtml_sub>`,
		"<vtml_sub alias=\"a\tb\">x</vtml_sub>",
	} {
		t.Run(source, func(t *testing.T) {
			if got, err := ExpandPaul2013CapturedVTMLSubstitutions(source); err == nil {
				t.Fatalf("ExpandPaul2013CapturedVTMLSubstitutions() = %q, want an error", got)
			}
		})
	}
}

func TestResolveTextWithPaul2013CapturedVTMLSubstitutions(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	source := `<vtml_sub alias="Hello Hello.">x</vtml_sub>`
	got, err := (LexiconFrontend{Dictionary: dictionary}).ResolveTextWithPaul2013CapturedVTMLSubstitutions(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Surface != "Hello" || got[1].Surface != "Hello" {
		t.Fatalf("resolved tokens = %+v, want two Hello tokens", got)
	}
	aliasStart := strings.Index(source, "Hello Hello.")
	wantSpans := [][2]int{{aliasStart, aliasStart + 4}, {aliasStart + 6, aliasStart + 10}}
	for index, token := range got {
		if !token.HasSourceByteSpan || token.SourceByteStart != wantSpans[index][0] || token.SourceByteEnd != wantSpans[index][1] {
			t.Errorf("token %d source span = [%d,%d], want [%d,%d]", index,
				token.SourceByteStart, token.SourceByteEnd, wantSpans[index][0], wantSpans[index][1])
		}
	}
	if _, err := SelectLexicalPronunciations(got, []int{0, 0}); err != nil {
		t.Fatalf("select resolved substitutions: %v", err)
	}
	if _, err := (LexiconFrontend{Dictionary: dictionary}).ResolveTextWithPaul2013CapturedVTMLSubstitutions(
		context.Background(), `Hel<vtml_sub alias="lo">x</vtml_sub>`,
	); err == nil {
		t.Fatal("token crossing a rewritten source boundary was accepted")
	}
}
