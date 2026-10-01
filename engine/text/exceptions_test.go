package text

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestNormalizePaul2013ExceptionSurface(t *testing.T) {
	attributes := Paul2013ExceptionCharacterAttributes()
	unsignedAttributes := Paul2013UnsignedCharacterAttributeTable()
	if unsignedAttributes[0xf7] != 0x08 {
		t.Fatalf("attribute at 0xf7 = %#x, want %#x", unsignedAttributes[0xf7], 0x08)
	}
	for value := 0xf8; value <= 0xff; value++ {
		if got := unsignedAttributes[byte(value)]; got != 0x40 {
			t.Fatalf("attribute at %#02x = %#x, want %#x", value, got, 0x40)
		}
	}
	normalized, category, err := NormalizePaul2013ExceptionSurface(
		[]byte("H'Expose-X"), true, attributes, Paul2013EmbeddedKeyTables().CharacterMap,
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(normalized) != "h'expose-x" || category != 2 {
		t.Fatalf("normalization = %q, category %d; want h'expose-x, 2", normalized, category)
	}
	if _, _, err := NormalizePaul2013ExceptionSurface(
		[]byte("EXPOSE-X"), false, attributes, Paul2013EmbeddedKeyTables().CharacterMap,
	); err == nil {
		t.Fatal("normalizer accepted a hyphen when disallowed")
	}
}

func TestSplitPaul2013ExceptionPhoneCodes(t *testing.T) {
	rows, err := SplitPaul2013ExceptionPhoneCodes([]byte("aadbbdcc\x00ignored"), []int{2, 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || string(rows[0]) != "aadbb" || string(rows[1]) != "cc" {
		t.Fatalf("split phone codes = %q, want [aadbb cc] across two rows", rows)
	}
	if _, err := SplitPaul2013ExceptionPhoneCodes([]byte("a d b d c"), []int{1}); err == nil {
		t.Fatal("exception phone code advanced beyond supplied destination rows")
	}
}

func TestLoadPaul2013ExceptionDictionary(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	dictionary, err := LoadExceptionDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	want := [4]int{11, 89, 22, 1}
	if got := dictionary.ExceptionGroupSizes(); !reflect.DeepEqual(got, want) {
		t.Fatalf("exception groups = %v, want %v", got, want)
	}
	first := dictionary.groups[0][0]
	value, ok := dictionary.Lookup(1, first.Key)
	if !ok || string(value) != string(first.Value) {
		t.Fatalf("first group lookup = % x, %t; want % x", value, ok, first.Value)
	}
	for _, test := range []struct {
		surface  string
		category uint32
	}{
		{surface: "H'Expose", category: 1},
		{surface: "San-Francisco", category: 2},
	} {
		if _, category, found, err := dictionary.LookupPaul2013Surface([]byte(test.surface), true); err != nil || !found || category != test.category {
			t.Fatalf("LookupPaul2013Surface(%q) = category %d, found %t, error %v; want category %d", test.surface, category, found, err, test.category)
		}
	}
	match, found, err := dictionary.LookupPaul2013SurfaceSequence([]string{"San", "Francisco", "unused"})
	if err != nil || !found || match.Category != 2 || match.TokenCount != 2 || len(match.PhoneCodes) == 0 {
		t.Fatalf("sequence lookup = %+v, found %t, error %v; want two-token category-2 match", match, found, err)
	}
	phoneRows, err := match.DecodePhoneRows([]int{1, 1})
	if err != nil {
		t.Fatalf("decode San-Francisco exception phone rows: %v", err)
	}
	if len(phoneRows) != 2 || len(phoneRows[0]) == 0 || len(phoneRows[1]) == 0 {
		t.Fatalf("decoded San-Francisco rows = %v; want two nonempty destination rows", phoneRows)
	}
	if _, err := match.DecodePhoneRows([]int{1}); err == nil {
		t.Fatal("exception match accepted a delimiter-count list with the wrong number of token rows")
	}
}
