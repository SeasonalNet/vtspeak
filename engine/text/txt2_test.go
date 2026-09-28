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
