package text

import (
	"strings"
	"testing"
)

func TestParsePaul2013VTMLCMUPhonemeCapturedInput(t *testing.T) {
	const source = `<vtml_phoneme alphabet="x-cmu" ph="P AH0">pah</vtml_phoneme>.`
	got, err := ParsePaul2013VTMLCMUPhoneme(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Phones) != 2 || got.Phones[0] != (CMUPhone{Label: "P"}) ||
		got.Phones[1] != (CMUPhone{Label: "AH", Vowel: true}) {
		t.Fatalf("forced phones = %+v, want P AH0", got.Phones)
	}
	if len(got.Symbols) != 2 || len(got.Tokens) != 1 {
		t.Fatalf("forced sequence = %+v, want two symbol codes and one token", got)
	}
	span := got.Tokens[0]
	if span.Surface != "pah" || span.SourceSurface != "pah" || span.SourceByteStart != strings.Index(source, "pah") ||
		span.SourceByteEnd != strings.Index(source, "pah")+2 || span.SeparatorAfter != "." {
		t.Fatalf("forced token span = %+v, want the source body and sentence period", span)
	}
}

func TestParsePaul2013VTMLCMUPhonemeRejectsUnportedForms(t *testing.T) {
	for _, source := range []string{
		"",
		`<vtml_phoneme alphabet="x-sampa" ph="p">p</vtml_phoneme>`,
		`<vtml_phoneme alphabet="x-cmu" ph="P AH0" mode="x">pah</vtml_phoneme>`,
		`<vtml_phoneme alphabet="x-cmu" ph="P AH0"><b>pah</b></vtml_phoneme>`,
		`<vtml_phoneme alphabet="x-cmu" ph="P AH0">pah</vtml_phoneme>!`,
		`text <vtml_phoneme alphabet="x-cmu" ph="P AH0">pah</vtml_phoneme>`,
		`<vtml_phoneme alphabet="x-cmu" ph="P AH0">pah</vtml_phoneme><vtml_phoneme alphabet="x-cmu" ph="T">t</vtml_phoneme>`,
	} {
		if _, err := ParsePaul2013VTMLCMUPhoneme(source); err == nil {
			t.Errorf("unsupported forced-phoneme markup %q succeeded", source)
		}
	}
}
