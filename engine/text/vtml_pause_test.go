package text

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveTextWithPaul2013InlinePauses(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	source := `Hello<vtml_pause time="200"/>Hello<vtml_pause time="1000"/>.`
	got, err := (LexiconFrontend{Dictionary: dictionary}).ResolveTextWithPaul2013InlinePauses(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Tokens) != 2 || got.Tokens[0].Surface != "Hello" || got.Tokens[1].Surface != "Hello" {
		t.Fatalf("tokens = %+v, want two Hello tokens", got.Tokens)
	}
	for index, token := range got.Tokens {
		if token.ModelPhoneRows.AlternativeCount != len(token.Alternatives) ||
			len(token.ModelPhoneRows.PhoneStrings) != len(token.Alternatives) {
			t.Fatalf("token %d model rows = %+v for %d alternatives", index, token.ModelPhoneRows, len(token.Alternatives))
		}
		for alternativeIndex, alternative := range token.Alternatives {
			if string(token.ModelPhoneRows.PhoneStrings[alternativeIndex]) != string(alternative.Symbols) {
				t.Fatalf("token %d model phone row %d = %q, want %q", index, alternativeIndex,
					token.ModelPhoneRows.PhoneStrings[alternativeIndex], alternative.Symbols)
			}
		}
	}
	resolvedRow := make([]byte, Paul2013TokenResultRowSize)
	if err := got.Tokens[0].WritePaul2013DictionaryPhoneRow(resolvedRow, 7); err != nil {
		t.Fatal(err)
	}
	if resolvedRow[Paul2013TokenResultRowIndex] != 7 ||
		string(resolvedRow[Paul2013TokenResultRowSurface:Paul2013TokenResultRowSurface+6]) != "Hello\x00" {
		t.Fatalf("resolved token row has wrong token index or surface: % x", resolvedRow[:0x2b])
	}
	secondHelloStart := strings.LastIndex(source, "Hello")
	if got.Tokens[0].SourceByteStart != 0 || got.Tokens[0].SourceByteEnd != 4 ||
		got.Tokens[1].SourceByteStart != secondHelloStart || got.Tokens[1].SourceByteEnd != secondHelloStart+4 {
		t.Fatalf("source spans = [%d,%d], [%d,%d], want [0,4], [%d,%d]",
			got.Tokens[0].SourceByteStart, got.Tokens[0].SourceByteEnd,
			got.Tokens[1].SourceByteStart, got.Tokens[1].SourceByteEnd,
			secondHelloStart, secondHelloStart+4)
	}
	if len(got.Pauses) != 2 {
		t.Fatalf("pause count = %d, want 2", len(got.Pauses))
	}
	want := []Paul2013InlinePause{
		{SourceByteOffset: strings.Index(source, "<vtml_pause"), AfterToken: 1, DurationMilliseconds: 200, OutputFramesAt16KHz: 3200},
		{SourceByteOffset: strings.LastIndex(source, "<vtml_pause"), AfterToken: 2, DurationMilliseconds: 1000, OutputFramesAt16KHz: 16000},
	}
	for index := range want {
		if got.Pauses[index] != want[index] {
			t.Errorf("pause %d = %+v, want %+v", index, got.Pauses[index], want[index])
		}
	}
	sequence, err := SelectLexicalPronunciationsWithPauses(got, []int{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(sequence.Tokens) != 2 || len(sequence.InlinePauses) != len(got.Pauses) ||
		sequence.InlinePauses[0] != got.Pauses[0] || sequence.InlinePauses[1] != got.Pauses[1] {
		t.Fatalf("selected sequence did not retain pause boundaries: %+v", sequence)
	}
	if sequence.Tokens[0].ModelPhoneRows.AlternativeCount != got.Tokens[0].ModelPhoneRows.AlternativeCount {
		t.Fatal("selected sequence lost its model phone rows")
	}
	selectedRow := make([]byte, Paul2013TokenResultRowSize)
	if err := sequence.Tokens[0].WritePaul2013DictionaryPhoneRow(selectedRow, 7); err != nil {
		t.Fatal(err)
	}
	if string(selectedRow) != string(resolvedRow) {
		t.Fatal("selected token row differs from its resolved token row")
	}
	phoneContextRows, err := BuildPaul2013PhoneContextRows(sequence, []uint16{7, 8}, []byte{'0', 'X'})
	if err != nil {
		t.Fatal(err)
	}
	if binary.LittleEndian.Uint16(phoneContextRows) != 2 {
		t.Fatalf("phone-context row count = %d, want 2", binary.LittleEndian.Uint16(phoneContextRows))
	}
	firstContextRow := phoneContextRows[Paul2013PhoneContextTableHeaderSize : Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize]
	if firstContextRow[Paul2013PhoneContextRowSourceMarker] != '0' ||
		string(firstContextRow[Paul2013PhoneContextRowSurface:Paul2013PhoneContextRowSurface+6]) != "Hello\x00" {
		t.Fatalf("text-derived phone-context row is missing selected token fields: % x", firstContextRow)
	}
	if len(got.Tokens[0].ModelPhoneRows.PathControlBytes) != 0 {
		sequence.Tokens[0].ModelPhoneRows.PathControlBytes[0] ^= 0xff
		if sequence.Tokens[0].ModelPhoneRows.PathControlBytes[0] == got.Tokens[0].ModelPhoneRows.PathControlBytes[0] {
			t.Fatal("selected model phone rows alias the resolved token")
		}
	}
}

func TestDispatchPaul2013InlinePauseModelSourceRule(t *testing.T) {
	var events []Paul2013InlinePauseTagEvent
	rule, err := Paul2013InlinePauseModelSourceRule(func(event Paul2013InlinePauseTagEvent) error {
		events = append(events, event)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	destination := make([]byte, 32)
	copy(destination, "left")
	result, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source:            []byte(`<vtml_pause time="250"/>right`),
		Destination:       destination,
		DestinationCursor: 4,
		Rules:             []Paul2013ModelSourceRule{rule},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Matched || !result.HandlerAccepted || result.SourceCursor != len(`<vtml_pause time="250"/>`) || result.DestinationCursor != 5 {
		t.Fatalf("dispatch result = %+v", result)
	}
	if got := string(result.Destination[:result.DestinationCursor]); got != "left " {
		t.Fatalf("destination = %q, want %q", got, "left ")
	}
	want := Paul2013InlinePauseTagEvent{
		SourceByteOffset:     0,
		DurationMilliseconds: 250,
		OutputFramesAt16KHz:  4000,
	}
	if len(events) != 1 || events[0] != want {
		t.Fatalf("pause events = %+v, want [%+v]", events, want)
	}
}

func TestDispatchPaul2013InlinePauseModelSourceRuleDeclinesUnsupportedForm(t *testing.T) {
	eventCount := 0
	rule, err := Paul2013InlinePauseModelSourceRule(func(Paul2013InlinePauseTagEvent) error {
		eventCount++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	destination := []byte("unchanged")
	result, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source:      []byte(`<vtml_pause duration="250"/>`),
		Destination: destination,
		Rules:       []Paul2013ModelSourceRule{rule},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Matched || result.HandlerAccepted || result.SourceCursor != 0 || result.DestinationCursor != 0 {
		t.Fatalf("unsupported form dispatch result = %+v", result)
	}
	if eventCount != 0 || string(result.Destination) != "unchanged" {
		t.Fatalf("unsupported form changed outputs: events=%d destination=%q", eventCount, result.Destination)
	}
}

func TestStripPaul2013InlinePauseTags(t *testing.T) {
	plain, offsets, pauses, err := stripPaul2013InlinePauseTags("a<vtml_pause time=\"0\"/>b")
	if err != nil {
		t.Fatal(err)
	}
	if plain != "a b" || len(offsets) != len(plain) || len(pauses) != 1 ||
		pauses[0].SourceByteOffset != 1 || pauses[0].DurationMilliseconds != 0 {
		t.Fatalf("parse = (%q, %v, %+v), want spaced text, byte map, and zero-duration event", plain, offsets, pauses)
	}
	if offsets[0] != 0 || offsets[1] != -1 || offsets[2] != len("a<vtml_pause time=\"0\"/>") {
		t.Fatalf("source offsets = %v, want source bytes for a, synthetic space, b", offsets)
	}
	for _, source := range []string{
		"a<unknown/>b",
		"a<vtml_pause time=\"\"/>b",
		"a<vtml_pause time=\"1.5\"/>b",
		"a<vtml_pause time=\"4294967296\"/>b",
		"a<vtml_pause time=\"2\" b=\"3\"/>b",
		"a<unfinished",
	} {
		if _, _, _, err := stripPaul2013InlinePauseTags(source); err == nil {
			t.Errorf("unsupported markup accepted: %q", source)
		}
	}
}
