package text

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalAmbiguousDictionaryAlternativesHaveOnePathGroupEach(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}

	ambiguousRecords := 0
	ambiguousAlternatives := 0
	for key, record := range dictionary.records {
		payload, err := ParsePhonePayload(record.Payload)
		if err != nil {
			t.Fatalf("parse %q: %v", key, err)
		}
		if len(payload.Pronunciations) < 2 {
			continue
		}
		ambiguousRecords++
		for alternativeIndex, alternative := range payload.Pronunciations {
			ambiguousAlternatives++
			groups, err := ParsePaul2013PronunciationPath(alternative.Path)
			if err != nil {
				t.Fatalf("parse %q alternative %d path: %v", key, alternativeIndex, err)
			}
			if len(groups) != 1 {
				t.Errorf("ambiguous record %q alternative %d has %d path groups, want one", key, alternativeIndex, len(groups))
				continue
			}
			if len(groups[0]) == 0 {
				t.Errorf("ambiguous record %q alternative %d has an empty path group", key, alternativeIndex)
			}
		}
	}
	if ambiguousRecords == 0 {
		t.Fatal("local embedded dictionary has no ambiguous pronunciation records")
	}
	t.Logf("validated one nonempty path group for each of %d alternatives across %d ambiguous records", ambiguousAlternatives, ambiguousRecords)
}

func TestPhonePayloadBuildPaul2013DictionaryPhoneRows(t *testing.T) {
	var codebook PhoneIDCodebook
	codebook[1] = [5]byte{'A', 'H', '0'}
	codebook[2] = [5]byte{'B', 'I', 'Y', '1'}

	for _, test := range []struct {
		name        string
		payload     PhonePayload
		wantCount   int
		wantControl []byte
		wantPhones  [][]byte
	}{
		{
			name: "single alternative suppresses path bytes",
			payload: PhonePayload{
				Flags: 0x01, ResultType: 'A', Metadata: [4]bool{true, false, true, false},
				Pronunciations: []Pronunciation{{Path: []byte{'x'}, Phone: []byte{1}}},
			},
			wantCount: 1, wantControl: []byte{0xff}, wantPhones: [][]byte{[]byte("AH0")},
		},
		{
			name: "multiple alternatives join path bytes",
			payload: PhonePayload{Pronunciations: []Pronunciation{
				{Path: []byte{'a', ')'}, Phone: []byte{1}},
				{Path: []byte{'b'}, Phone: []byte{2}},
			}},
			wantCount: 2, wantControl: []byte{'a', ')', 'd', 'b', 0xff},
			wantPhones: [][]byte{[]byte("AH0"), []byte("BIY1")},
		},
		{
			name:        "empty payload rows retain terminator",
			payload:     PhonePayload{},
			wantControl: []byte{0xff},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.payload.BuildPaul2013DictionaryPhoneRows(codebook)
			if err != nil {
				t.Fatal(err)
			}
			if got.AlternativeCount != test.wantCount || string(got.PathControlBytes) != string(test.wantControl) {
				t.Fatalf("count/path controls = %d/% x, want %d/% x", got.AlternativeCount, got.PathControlBytes, test.wantCount, test.wantControl)
			}
			if len(got.PhoneStrings) != len(test.wantPhones) {
				t.Fatalf("phone rows = %d, want %d", len(got.PhoneStrings), len(test.wantPhones))
			}
			for index := range test.wantPhones {
				if string(got.PhoneStrings[index]) != string(test.wantPhones[index]) {
					t.Errorf("phone row %d = %q, want %q", index, got.PhoneStrings[index], test.wantPhones[index])
				}
			}
			if got.ResultType != test.payload.ResultType || got.Metadata != test.payload.Metadata {
				t.Errorf("result metadata = %q/%v, want %q/%v", got.ResultType, got.Metadata, test.payload.ResultType, test.payload.Metadata)
			}
		})
	}
}

func TestPhonePayloadBuildPaul2013DictionaryPhoneRowsForState(t *testing.T) {
	var codebook PhoneIDCodebook
	codebook[1] = [5]byte{'A', 'H', '0'}
	codebook[2] = [5]byte{'B', 'I', 'Y', '1'}
	payload := PhonePayload{Pronunciations: []Pronunciation{
		{Path: []byte{'a'}, Phone: []byte{1}},
		{Path: []byte{'b'}, Phone: []byte{2}},
	}}

	conditional, err := payload.BuildPaul2013DictionaryPhoneRowsForState(codebook, Paul2013DictionaryPhoneRowState{
		HasMarker: true, Marker: 'b', SelectPathPhone: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if conditional.AlternativeCount != 0 || !conditional.HasContextMarker || conditional.ContextMarker != 'b' ||
		!conditional.HasMarkerSentinel || conditional.MarkerSentinel != 0xff ||
		string(conditional.SelectedPhone) != "BIY1" || len(conditional.PhoneStrings) != 0 {
		t.Fatalf("conditional model row = %+v", conditional)
	}

	final, err := payload.BuildPaul2013DictionaryPhoneRowsForState(codebook, Paul2013DictionaryPhoneRowState{
		FinalRow: true, HasMarker: true, Marker: 'b', SelectPathPhone: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if final.AlternativeCount != 2 || final.HasContextMarker || final.HasMarkerSentinel ||
		string(final.PathControlBytes) != string([]byte{'a', 'd', 'b', 0xff}) || len(final.PhoneStrings) != 2 {
		t.Fatalf("final model row did not use the all-alternatives branch: %+v", final)
	}
}

func TestWritePaul2013DictionaryPhoneRow(t *testing.T) {
	rows := Paul2013DictionaryPhoneRows{
		ResultType: 'A', Metadata: [4]bool{true, false, true, true},
		AlternativeCount: 2, PathControlBytes: []byte{'a', ')', 'd', 'b', 0xff},
		PhoneStrings: [][]byte{[]byte("AH0"), []byte("BIY1")},
	}
	destination := make([]byte, Paul2013TokenResultRowSize)
	for index := range destination {
		destination[index] = 0xa5
	}
	if err := WritePaul2013DictionaryPhoneRow(destination, 0x1234, []byte("Word"), rows); err != nil {
		t.Fatal(err)
	}
	for offset, want := range map[int][]byte{
		Paul2013TokenResultRowCount:                                        {2, 0},
		Paul2013TokenResultRowIndex:                                        {0x34, 0x12},
		Paul2013TokenResultRowType:                                         {'A'},
		Paul2013TokenResultRowSurface:                                      {'W', 'o', 'r', 'd', 0},
		Paul2013TokenResultPathControls:                                    {'a', ')', 'd', 'b', 0xff},
		Paul2013TokenResultPhoneStrings:                                    {'A', 'H', '0', 0},
		Paul2013TokenResultPhoneStrings + Paul2013TokenResultPhoneStride:   {'B', 'I', 'Y', '1', 0},
		Paul2013TokenResultPhoneStrings + 2*Paul2013TokenResultPhoneStride: {0},
		Paul2013TokenResultMetadata:                                        {1, 0},
		Paul2013TokenResultMetadata + 2:                                    {0, 0},
		Paul2013TokenResultMetadata + 4:                                    {1, 0},
		Paul2013TokenResultMetadata + 6:                                    {1, 0},
	} {
		if got := destination[offset : offset+len(want)]; string(got) != string(want) {
			t.Errorf("row bytes at %#x = % x, want % x", offset, got, want)
		}
	}
	if destination[0x17c] != 0xa5 {
		t.Fatalf("unwritten gap byte %#x changed for two alternatives", destination[0x17c])
	}

	five := rows
	five.AlternativeCount = 5
	five.PathControlBytes = []byte{'a', 'd', 'b', 'd', 'c', 'd', 'e', 'd', 'f', 0xff}
	five.PhoneStrings = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
	if err := WritePaul2013DictionaryPhoneRow(destination, 0, []byte("Word"), five); err != nil {
		t.Fatal(err)
	}
	if destination[0x17c] != 0 {
		t.Fatalf("five-alternative terminator at %#x = %#x, want NUL", 0x17c, destination[0x17c])
	}

	conditional := Paul2013DictionaryPhoneRows{
		ResultType: 'E', Metadata: [4]bool{true}, ContextMarker: 'x', HasContextMarker: true,
		MarkerSentinel: 0xff, HasMarkerSentinel: true, SelectedPhone: []byte("KAA1"),
	}
	if err := WritePaul2013DictionaryPhoneRow(destination, 7, []byte("Word"), conditional); err != nil {
		t.Fatal(err)
	}
	if destination[0x12] != 0xff || destination[Paul2013TokenResultPathControls] != 'x' ||
		string(destination[Paul2013TokenResultPhoneStrings:Paul2013TokenResultPhoneStrings+5]) != "KAA1\x00" {
		t.Fatalf("conditional row marker/phone fields were not written: % x", destination[:0x3c])
	}
}

func TestWritePaul2013DictionaryPhoneRowRejectsInvalidInputWithoutMutation(t *testing.T) {
	rows := Paul2013DictionaryPhoneRows{
		AlternativeCount: 1, PathControlBytes: []byte{0xff}, PhoneStrings: [][]byte{[]byte("AH0")},
	}
	destination := make([]byte, Paul2013TokenResultRowSize)
	for index := range destination {
		destination[index] = 0xa5
	}
	before := append([]byte(nil), destination...)
	if err := WritePaul2013DictionaryPhoneRow(destination, 0, make([]byte, 30), rows); err == nil {
		t.Fatal("oversized surface was accepted")
	}
	for index := range destination {
		if destination[index] != before[index] {
			t.Fatalf("invalid write changed destination byte %#x", index)
		}
	}
	if err := WritePaul2013DictionaryPhoneRow(destination[:Paul2013TokenResultRowSize-1], 0, []byte("Word"), rows); err == nil {
		t.Fatal("short destination was accepted")
	}
}
