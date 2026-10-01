package text

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTerminalAbbreviationVariantSearch(t *testing.T) {
	input := &Table{Name: "abbrh_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "AA", Value: "1"}, {Key: "Aa", Value: "1"}, {Key: "aa", Value: "1"}, {Key: "z", Value: "2"}}}
	table, err := NewPaul2013TerminalAbbreviationTable(input)
	if err != nil {
		t.Fatal(err)
	}
	input.Rows[1].Value = "0" // The compiled snapshot owns its rows.
	for _, key := range []string{"AA", "Aa", "aa"} {
		if got, err := table.Lookup([]byte(key)); err != nil || got != 1 {
			t.Fatalf("%s -> %d, %v", key, got, err)
		}
	}
	for _, key := range []string{"Z", "z"} {
		if got, err := table.Lookup([]byte(key)); err != nil || got != 3 {
			t.Fatal(key, got, err)
		}
	}
	if got, err := table.Lookup([]byte("missing")); err != nil || got != -1 {
		t.Fatal(got, err)
	}
	input.Rows = []Row{{Key: "AA", Value: "0"}, {Key: "Aa", Value: "0"}, {Key: "aa", Value: "0"}}
	table, err = NewPaul2013TerminalAbbreviationTable(input)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := table.Lookup([]byte("AA")); err != nil || got != -1 {
		t.Fatal(got, err)
	}
	if _, err := NewPaul2013TerminalAbbreviationTable(&Table{Name: "abbrh_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "z"}, {Key: "a"}}}); err == nil {
		t.Fatal("unsorted table accepted")
	}
	if _, err := (*Paul2013TerminalAbbreviationTable)(nil).Lookup(nil); err == nil {
		t.Fatal("nil table accepted")
	}
}

func TestTerminalAbbreviationMaskAndNativeFeatures(t *testing.T) {
	var tables [3]*Paul2013TerminalAbbreviationTable
	for index, name := range []string{"abbrh_sort.txt2", "abbrt_sort.txt2", "abbrc_sort.txt2"} {
		value := "0"
		if index != 1 {
			value = "2"
		}
		var err error
		tables[index], err = NewPaul2013TerminalAbbreviationTable(&Table{Name: name, Mode: TwoColumns, Rows: []Row{{Key: "word", Value: value}}})
		if err != nil {
			t.Fatal(err)
		}
	}
	if mask, err := Paul2013TerminalAbbreviationMask(tables, []byte("WORD\x00ignored")); err != nil || mask != 5 {
		t.Fatal(mask, err)
	}
	for _, test := range []struct {
		text   string
		status int32
		want   int16
	}{{"", 1, -1}, {"ABC", 1, 1}, {"aBc", 1, 5}, {"abc", 1, 2}, {"Abc", 1, 3}, {"ABC9", 1, 1}, {"9", 1, 5}, {"word", 2, 4}, {"Word", 3, 5}, {"\xff", 1, 5}} {
		got, err := ClassifyPaul2013TerminalTokenCase(Paul2013ModelScannerResult{Text: []byte(test.text), TextLength: int32(len(test.text)), Status: test.status})
		if err != nil || got != test.want {
			t.Fatalf("case %+v -> %d %v", test, got, err)
		}
	}
	if _, err := ClassifyPaul2013TerminalTokenCase(Paul2013ModelScannerResult{TextLength: 37, Status: 1}); err == nil {
		t.Fatal("scanner overread accepted")
	}
	for _, test := range []struct {
		text   string
		length int32
		want   int16
	}{{".", 1, 1}, {"?", 1, 2}, {"!", 1, 3}, {";", 1, 4}, {",", 1, 5}, {":", 1, 5}, {"\"", 1, 6}, {"`", 1, 6}, {"-", 1, 7}, {"'", 1, -1}, {"..", 2, -1}, {"", 1, -1}} {
		if got := ClassifyPaul2013TerminalTokenPunctuation(Paul2013ModelScannerResult{Text: []byte(test.text), TextLength: test.length}); got != test.want {
			t.Fatalf("punct %+v -> %d", test, got)
		}
	}
	for _, test := range []struct {
		text string
		want int16
	}{{"a.b", 1}, {"a..", 1}, {".a.b", 0}, {"ab.", 0}, {"abc", 0}, {"a\x00.b", 0}} {
		if got := Paul2013TerminalTokenHasInternalDot([]byte(test.text)); got != test.want {
			t.Fatal(test, got)
		}
	}
}

func TestTerminalAbbreviationAllLocalResourceRows(t *testing.T) {
	weights := Paul2013ContextCharacterWeights()
	for _, name := range []string{"abbrh_sort.txt2", "abbrt_sort.txt2", "abbrc_sort.txt2"} {
		raw, err := os.ReadFile(filepath.Join("..", "..", "data-common", "dict-eng", name))
		if os.IsNotExist(err) {
			t.Skip("local proprietary shared data absent")
		}
		if err != nil {
			t.Fatal(err)
		}
		input, err := ParseTXT2(name, raw)
		if err != nil {
			t.Fatal(err)
		}
		table, err := NewPaul2013TerminalAbbreviationTable(input)
		if err != nil {
			t.Fatal(name, err)
		}
		for _, row := range input.Rows {
			for _, key := range []string{row.Key, strings.ToUpper(row.Key), strings.ToLower(row.Key)} {
				want := false
				for _, candidate := range input.Rows {
					if ComparePaul2013MappedCString([]byte(candidate.Key), []byte(key), weights) == 0 && len(candidate.Value) > 0 && (candidate.Value[0] == '2' || candidate.Value[0] == '1' && bytes.Equal([]byte(candidate.Key), []byte(key))) {
						want = true
						break
					}
				}
				got, err := table.Lookup([]byte(key))
				if err != nil || (got >= 0) != want {
					t.Fatalf("%s %q -> %d want match %v, %v", name, key, got, want, err)
				}
			}
		}
	}
}
