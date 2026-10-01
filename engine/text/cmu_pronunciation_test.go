package text

import "testing"

func TestBuildPaul2013CMUPhoneSequence(t *testing.T) {
	sequence, err := BuildPaul2013CMUPhoneSequence("pah", "P AH0")
	if err != nil {
		t.Fatal(err)
	}
	if len(sequence.Phones) != 2 || len(sequence.Symbols) != 2 || len(sequence.Tokens) != 1 {
		t.Fatalf("CMU phone sequence = %+v, want two phones, symbols, and one token", sequence)
	}
	if sequence.Phones[0] != (CMUPhone{Label: "P"}) || sequence.Phones[1] != (CMUPhone{Label: "AH", Stress: 0, Vowel: true}) {
		t.Fatalf("CMU phones = %+v, want P AH0", sequence.Phones)
	}
	if sequence.Tokens[0].SourceByteStart != 0 || sequence.Tokens[0].SourceByteEnd != 2 || !sequence.Tokens[0].HasSourceByteSpan {
		t.Fatalf("CMU token span = %+v, want inclusive surface span 0..2", sequence.Tokens[0])
	}
	for index, phone := range sequence.Phones {
		features, err := phone.TreeFeatures()
		if err != nil {
			t.Fatal(err)
		}
		if sequence.Symbols[index] != features.SymbolCode {
			t.Errorf("phone %d symbol code = %#x, want %#x", index, sequence.Symbols[index], features.SymbolCode)
		}
	}
}

func TestParsePaul2013CMUPronunciationRejectsUnsupportedInput(t *testing.T) {
	for _, source := range []string{"", "p AH0", "P AH3", "B0", "Q AH0", "P\tAH0", "P AH0!", "P AHé"} {
		if _, err := ParsePaul2013CMUPronunciation(source); err == nil {
			t.Errorf("unsupported CMU pronunciation %q succeeded", source)
		}
	}
}

func TestBuildPaul2013CMUPhoneSequenceRejectsInvalidSurface(t *testing.T) {
	for _, surface := range []string{"", "p\x00ah", "pahé"} {
		if _, err := BuildPaul2013CMUPhoneSequence(surface, "P AH0"); err == nil {
			t.Errorf("invalid CMU pronunciation surface %q succeeded", surface)
		}
	}
}
