package text

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestParseTXT2AndExactLookup(t *testing.T) {
	raw := makeTXT2("M02", 10, "dr|1\nDr|2")
	table, err := ParseTXT2("abbrh_sort.txt2", raw)
	if err != nil {
		t.Fatal(err)
	}
	if table.Mode != TwoColumns || len(table.Rows) != 2 {
		t.Fatalf("unexpected table metadata: mode=%d rows=%d", table.Mode, len(table.Rows))
	}
	if got := table.FindExact("dr"); len(got) != 1 || got[0].Value != "1" {
		t.Fatalf("exact lowercase lookup = %v", got)
	}
	if got := table.FindExact("DR"); len(got) != 0 {
		t.Fatalf("exact lookup unexpectedly folded case: %v", got)
	}
}

func TestLookupPaul2013WABClassUsesMappedSortedSearch(t *testing.T) {
	table := &Table{
		Name: "wab.txt2",
		Mode: OneColumn,
		Rows: []Row{{Key: "APPLE"}, {Key: "BANANA"}, {Key: "PEAR"}},
	}
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	for _, surface := range []string{"apple", "APPLE", "Apple"} {
		index, found, err := table.LookupPaul2013WABClass([]byte(surface), characterMap)
		if err != nil || !found || index != 0 {
			t.Fatalf("WAB lookup %q = (%d, %t, %v), want (0, true, nil)", surface, index, found, err)
		}
	}
	if _, found, err := table.LookupPaul2013WABClass([]byte("APRICOT"), characterMap); err != nil || found {
		t.Fatalf("missing WAB lookup = found %t, err %v; want false, nil", found, err)
	}

	table.Rows[1], table.Rows[2] = table.Rows[2], table.Rows[1]
	if _, _, err := table.LookupPaul2013WABClass([]byte("APPLE"), characterMap); err == nil {
		t.Fatal("unsorted WAB table was accepted")
	}
}

func TestLookupPaul2013CHCFlagsUsesNativeBitOrderAndExactKeys(t *testing.T) {
	table := &Table{
		Name: "chc_sort.txt2",
		Mode: TwoColumns,
		Rows: []Row{{Key: "alpha", Value: "1001"}, {Key: "bravo", Value: "0110"}, {Key: "charlie", Value: "0001"}},
	}
	for _, test := range []struct {
		key      string
		mask     uint16
		wantRow  int
		wantOkay bool
	}{
		{key: "alpha", mask: 8, wantRow: 0, wantOkay: true},
		{key: "alpha", mask: 1, wantRow: 0, wantOkay: true},
		{key: "alpha", mask: 4, wantRow: 0, wantOkay: false},
		{key: "bravo", mask: 6, wantRow: 1, wantOkay: true},
		{key: "BRAVO", mask: 2, wantRow: 0, wantOkay: false},
	} {
		row, found, err := table.LookupPaul2013CHCFlags([]byte(test.key), test.mask)
		if err != nil {
			t.Fatalf("CHC lookup %q/%d: %v", test.key, test.mask, err)
		}
		if found != test.wantOkay || (found && row != test.wantRow) {
			t.Errorf("CHC lookup %q/%d = (%d, %t), want (%d, %t)", test.key, test.mask, row, found, test.wantRow, test.wantOkay)
		}
	}
	if row, found, err := table.LookupPaul2013CHCFlags([]byte("charlie\x00ignored"), 1); err != nil || !found || row != 2 {
		t.Fatalf("NUL-terminated CHC lookup = (%d, %t, %v), want (2, true, nil)", row, found, err)
	}
}

func TestLookupPaul2013CHCFlagsRejectsMalformedTable(t *testing.T) {
	for _, test := range []struct {
		name  string
		table *Table
	}{
		{name: "wrong table", table: &Table{Name: "wab.txt2", Mode: OneColumn}},
		{name: "unsorted", table: &Table{Name: "chc_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "b", Value: "1000"}, {Key: "a", Value: "1000"}}}},
		{name: "duplicate", table: &Table{Name: "chc_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "a", Value: "1000"}, {Key: "a", Value: "0100"}}}},
		{name: "short mask", table: &Table{Name: "chc_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "a", Value: "1"}}}},
		{name: "invalid mask", table: &Table{Name: "chc_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "a", Value: "10x0"}}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := test.table.LookupPaul2013CHCFlags([]byte("a"), 1); err == nil {
				t.Fatal("malformed CHC table accepted")
			}
		})
	}
}

func TestMatchPaul2013CHCSubstringPortsCompositeAndFallbackBranches(t *testing.T) {
	table := &Table{
		Name: "chc_sort.txt2",
		Mode: TwoColumns,
		Rows: []Row{
			{Key: "AB", Value: "0001"},
			{Key: "CD", Value: "0100"},
			{Key: "DOG", Value: "0001"},
			{Key: "EXACT", Value: "0010"},
		},
	}
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	tests := []struct {
		name        string
		word        string
		tokenLength int
		mask        uint16
		want        bool
	}{
		{name: "two-part compound", word: "ABCD", tokenLength: 4, mask: 2, want: true},
		{name: "exact mask", word: "EXACT", tokenLength: 5, mask: 2, want: true},
		{name: "mapped trailing S", word: "DOGS", tokenLength: 4, mask: 1, want: true},
		{name: "long-token short-substring fallback", word: "xyz", tokenLength: 7, mask: 1, want: true},
		{name: "mask-eight fallback excluded", word: "xyz", tokenLength: 7, mask: 8, want: false},
		{name: "candidate longer than token", word: "ABCD", tokenLength: 3, mask: 2, want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := table.MatchPaul2013CHCSubstring([]byte(test.word), test.tokenLength, test.mask, characterMap)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("CHC substring match = %t, want %t", got, test.want)
			}
		})
	}
}

func TestParseTXT2RejectsInvalidResources(t *testing.T) {
	valid := makeTXT2("L01", 10, "hello")
	markerMismatch := append([]byte(nil), valid...)
	markerMismatch[len(markerMismatch)-1] ^= 1
	wrongColumnCount := makeTXT2("L01", 10, "single|column")
	wrongRowCount := makeTXT2Declared("L02", 10, "one", 2)
	for name, raw := range map[string][]byte{
		"short":     valid[:10],
		"markers":   markerMismatch,
		"columns":   wrongColumnCount,
		"row count": wrongRowCount,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseTXT2("streetf_sort.txt2", raw); err == nil {
				t.Fatal("invalid resource accepted")
			}
		})
	}
	if _, err := ParseTXT2("unknown.txt2", valid); err == nil {
		t.Fatal("resource without a mapped column mode accepted")
	}
}

func TestLoadSharedPaulTXT2TablesWhenAssetsArePresent(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "wab.txt2")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	tables, err := LoadTXT2Tables(root)
	if err != nil {
		t.Fatal(err)
	}
	wantRows := map[string]int{
		"wab.txt2": 113, "chc_sort.txt2": 2581, "streeta_sort.txt2": 347,
		"streetf_sort.txt2": 211, "citya_sort.txt2": 117, "abbrh_sort.txt2": 55,
		"abbrt_sort.txt2": 42, "abbrc_sort.txt2": 330, "sbdw_sort.txt2": 1067,
	}
	if len(tables) != len(wantRows) {
		t.Fatalf("loaded %d tables, expected %d", len(tables), len(wantRows))
	}
	for name, want := range wantRows {
		if table := tables[name]; table == nil || len(table.Rows) != want {
			got := 0
			if table != nil {
				got = len(table.Rows)
			}
			t.Errorf("%s row count = %d, expected %d", name, got, want)
		}
	}
	chc := tables["chc_sort.txt2"]
	if len(chc.Rows) != 0 && len(chc.Rows[0].Value) >= 4 {
		mask := uint16(0)
		for index, value := range []byte(chc.Rows[0].Value[:4]) {
			if value == '1' {
				mask |= 1 << (3 - index)
			}
		}
		index, found, err := chc.LookupPaul2013CHCFlags([]byte(chc.Rows[0].Key), mask)
		if err != nil || !found || index != 0 {
			t.Fatalf("loaded CHC table first-row lookup = (%d, %t, %v), want (0, true, nil)", index, found, err)
		}
	}
}

func makeTXT2(marker string, shift byte, body string) []byte {
	rows := strings.FieldsFunc(body, func(r rune) bool { return r == '\r' || r == '\n' })
	return makeTXT2Declared(marker, shift, body, len(rows))
}

func makeTXT2Declared(marker string, shift byte, body string, rowCount int) []byte {
	encoded := func(raw []byte) []byte {
		out := make([]byte, len(raw))
		for i, value := range raw {
			out[i] = value + shift
		}
		return out
	}
	count := make([]byte, 10)
	copy(count, strconv.Itoa(rowCount))
	raw := append([]byte(marker), encoded([]byte(body))...)
	raw = append(raw, shift)
	raw = append(raw, encoded(count)...)
	raw = append(raw, []byte(marker)...)
	return raw
}
