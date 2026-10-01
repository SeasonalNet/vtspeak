package text

import (
	"path/filepath"
	"testing"
)

func TestTPPNumericBytePortsFGNarrowing(t *testing.T) {
	for _, test := range []struct {
		atom TPPAtom
		want byte
	}{
		{atom: TPPAtom{Tag: 'F', Suffix: []byte("120")}, want: 120},
		{atom: TPPAtom{Tag: 'G', Suffix: []byte("95")}, want: 95},
		{atom: TPPAtom{Tag: 'F', Suffix: []byte("257")}, want: 1},
	} {
		got, ok, err := test.atom.NumericByte()
		if err != nil || !ok || got != test.want {
			t.Fatalf("NumericByte(%+v) = %d, %t, %v; want %d, true, nil", test.atom, got, ok, err, test.want)
		}
	}
	if _, ok, err := (TPPAtom{Tag: 'A', Suffix: []byte("0")}).NumericByte(); err != nil || ok {
		t.Fatalf("non-F/G numeric operation = %t, %v; want false, nil", ok, err)
	}
	if _, _, err := (TPPAtom{Tag: 'G', Suffix: []byte("x")}).NumericByte(); err == nil {
		t.Fatal("invalid numeric suffix accepted")
	}
}

func TestLoadPaul2013TPPDictionaryAndLookup(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	dictionary, err := LoadTPPDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(dictionary.records); got != paul2013TPPRecordCount {
		t.Fatalf("TPP record count = %d, want %d", got, paul2013TPPRecordCount)
	}
	for _, test := range []struct {
		surface string
		payload string
		tag     byte
		suffix  string
		numeric byte
	}{
		{surface: "ACCORD", payload: "A0 G95", tag: 'G', suffix: "95", numeric: 95},
		{surface: "ABOUT-SHIPPING", payload: "F120", tag: 'F', suffix: "120", numeric: 120},
	} {
		record, ok, err := dictionary.Lookup([]byte(test.surface), Paul2013EmbeddedKeyTables())
		if err != nil || !ok {
			t.Fatalf("Lookup(%q) = (%+v, %t, %v)", test.surface, record, ok, err)
		}
		if string(record.Payload) != test.payload {
			t.Fatalf("Lookup(%q) payload = %q, want %q", test.surface, record.Payload, test.payload)
		}
		atom := record.Atoms[len(record.Atoms)-1]
		if atom.Tag != test.tag || string(atom.Suffix) != test.suffix {
			t.Fatalf("Lookup(%q) final atom = %+v, want %c%s", test.surface, atom, test.tag, test.suffix)
		}
		value, numeric, err := atom.NumericByte()
		if err != nil || !numeric || value != test.numeric {
			t.Fatalf("%c numeric byte = %d, %t, %v; want %d, true, nil", test.tag, value, numeric, err, test.numeric)
		}
		textValue, found, err := dictionary.LookupNumericText([]byte(test.surface), test.tag, Paul2013EmbeddedKeyTables())
		if err != nil || !found || string(textValue) != test.suffix {
			t.Fatalf("LookupNumericText(%q, %c) = %q, %t, %v; want %q, true, nil", test.surface, test.tag, textValue, found, err, test.suffix)
		}
	}
	if _, found, err := dictionary.LookupNumericText([]byte("ACCORD"), 'F', Paul2013EmbeddedKeyTables()); err != nil || found {
		t.Fatalf("LookupNumericText(ACCORD, F) = found %t, err %v; want false, nil", found, err)
	}

	window, err := MatchPaul2013TPPWindowWithDictionary(Paul2013TPPWindowInput{
		Rows: []Paul2013TPPWindowRow{
			{Word2: 1, State: 'S', Text: []byte("ABOUT")},
			{Word0: 1, Word1: 0, Word2: 1, State: 'A', Text: []byte("SHIPPING")},
		},
		CharacterMap: Paul2013EmbeddedKeyTables().CharacterMap,
	}, dictionary, Paul2013EmbeddedKeyTables())
	if err != nil || !window.Matched || string(window.Compound) != "ABOUT-SHIPPING" || string(window.Replacement) != "120" {
		t.Fatalf("F window lookup = %+v, %v; want ABOUT-SHIPPING -> 120", window, err)
	}

	_, err = MatchPaul2013TPPWindowWithDictionary(Paul2013TPPWindowInput{
		Rows:            []Paul2013TPPWindowRow{{Text: []byte("ABOUT")}, {Word0: 1, Word1: 0, State: 'A', Text: []byte("SHIPPING")}},
		DictionaryIndex: 1,
	}, dictionary, Paul2013EmbeddedKeyTables())
	if err == nil {
		t.Fatal("unloaded TPP dictionary index was accepted")
	}

	classLookup := func(surface []byte) (int, bool, error) {
		switch string(surface) {
		case "ABOUT":
			return 2, true, nil
		case "SHIPPING":
			return 4, true, nil
		case "ACCORD":
			return 5, true, nil
		default:
			return 7, true, nil
		}
	}
	rowsResult, err := ApplyPaul2013TPPNumericTextRows(Paul2013TPPNumericTextInput{
		Rows: []Paul2013TPPWindowRow{
			{Word2: 1, State: 'S', Text: []byte("ABOUT")},
			{Word0: 1, Word1: 0, Word2: 1, State: 'A', Text: []byte("SHIPPING")},
		},
		Dictionary: dictionary, KeyTables: Paul2013EmbeddedKeyTables(), LookupWABClass: classLookup,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rowsResult.Operations) != 1 || rowsResult.Operations[0] != (Paul2013TPPNumericTextOperation{Tag: 'F', StartRow: 0, EndRow: 1}) {
		t.Fatalf("F consumer operations = %+v, want one rows 0-1 operation", rowsResult.Operations)
	}
	if len(rowsResult.Updates) != 2 || rowsResult.Updates[0] != (Paul2013TPPTokenUpdate{ProcessingFlagWritten: true, ProcessingFlag: 1, TypedCodeWritten: true, TypedCode: 120, ClassCodeWritten: true, ClassCode: 3}) ||
		rowsResult.Updates[1] != (Paul2013TPPTokenUpdate{TypedCodeWritten: true, TypedCode: 120, ClassCodeWritten: true, ClassCode: 5}) {
		t.Fatalf("F consumer row updates = %+v", rowsResult.Updates)
	}

	rowsResult, err = ApplyPaul2013TPPNumericTextRows(Paul2013TPPNumericTextInput{
		Rows:       []Paul2013TPPWindowRow{{Text: []byte("ACCORD")}},
		Dictionary: dictionary, KeyTables: Paul2013EmbeddedKeyTables(), LookupWABClass: classLookup,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rowsResult.Operations) != 1 || rowsResult.Operations[0].Tag != 'G' || !rowsResult.Updates[0].TypedCodeWritten || rowsResult.Updates[0].TypedCode != 95 ||
		!rowsResult.Updates[0].ClassCodeWritten || rowsResult.Updates[0].ClassCode != 6 {
		t.Fatalf("G consumer updates = %+v, operations %+v", rowsResult.Updates, rowsResult.Operations)
	}

	rowsResult, err = ApplyPaul2013TPPNumericTextRows(Paul2013TPPNumericTextInput{
		Rows:       []Paul2013TPPWindowRow{{Text: []byte("UNLISTED")}},
		Dictionary: dictionary, KeyTables: Paul2013EmbeddedKeyTables(), LookupWABClass: classLookup,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rowsResult.Operations) != 0 || rowsResult.Updates[0] != (Paul2013TPPTokenUpdate{ClassCodeWritten: true, ClassCode: 8}) {
		t.Fatalf("WAB fallback updates = %+v, operations %+v", rowsResult.Updates, rowsResult.Operations)
	}
}

func TestParseTPPAtomsRejectsUnobservedSequences(t *testing.T) {
	for _, payload := range [][]byte{
		[]byte("A0 B201"),
		[]byte("A0 G95 F12"),
		[]byte("AX G95"),
		[]byte("Gx"),
	} {
		if _, err := parseTPPAtoms(payload); err == nil {
			t.Fatalf("parseTPPAtoms(%q) succeeded", payload)
		}
	}
}

func TestApplyPaul2013TPPNumericAtom(t *testing.T) {
	for _, test := range []struct {
		name         string
		atom         TPPAtom
		classIndexes []int
		wantCode     byte
		wantFlags    []bool
	}{
		{
			name:         "compound F suffix",
			atom:         TPPAtom{Tag: 'F', Suffix: []byte("120")},
			classIndexes: []int{3, -1},
			wantCode:     120,
			wantFlags:    []bool{true, false},
		},
		{
			name:         "single G suffix",
			atom:         TPPAtom{Tag: 'G', Suffix: []byte("95")},
			classIndexes: []int{7},
			wantCode:     95,
			wantFlags:    []bool{false},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			updates, applied, err := ApplyPaul2013TPPNumericAtom(test.atom, test.classIndexes)
			if err != nil || !applied {
				t.Fatalf("ApplyPaul2013TPPNumericAtom() = (%v, %t, %v)", updates, applied, err)
			}
			if len(updates) != len(test.classIndexes) {
				t.Fatalf("received %d updates, want %d", len(updates), len(test.classIndexes))
			}
			for tokenIndex, update := range updates {
				if update.TypedCode != test.wantCode {
					t.Errorf("token %d typed code = %d, want %d", tokenIndex, update.TypedCode, test.wantCode)
				}
				if update.ProcessingFlagWritten != test.wantFlags[tokenIndex] {
					t.Errorf("token %d processing flag written = %t, want %t", tokenIndex, update.ProcessingFlagWritten, test.wantFlags[tokenIndex])
				}
				if test.classIndexes[tokenIndex] >= 0 && (!update.ClassCodeWritten || update.ClassCode != byte(test.classIndexes[tokenIndex]+1)) {
					t.Errorf("token %d class update = %+v, want class code %d", tokenIndex, update, test.classIndexes[tokenIndex]+1)
				}
				if test.classIndexes[tokenIndex] < 0 && update.ClassCodeWritten {
					t.Errorf("token %d wrote class code despite missing class lookup: %+v", tokenIndex, update)
				}
			}
		})
	}
	if updates, applied, err := ApplyPaul2013TPPNumericAtom(TPPAtom{Tag: 'A', Suffix: []byte("0")}, []int{0}); err != nil || applied || updates != nil {
		t.Fatalf("non-numeric TPP atom update = (%v, %t, %v), want no update", updates, applied, err)
	}
	if _, _, err := ApplyPaul2013TPPNumericAtom(TPPAtom{Tag: 'F', Suffix: []byte("invalid")}, []int{0}); err == nil {
		t.Fatal("malformed numeric TPP atom was accepted")
	}
}
