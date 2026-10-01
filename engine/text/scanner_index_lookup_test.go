package text

import (
	"context"
	"errors"
	"testing"
)

func TestScannerIndexedLookupNativeOrder(t *testing.T) {
	table := &Paul2013ScannerLookupIndex{Lower: []int32{-1, -1, 0}, Upper: []int32{-1, -1, 4}, Keys: [][]byte{[]byte("aaa"), []byte("bbb"), []byte("ccc"), []byte("ddd"), []byte("zzz")}}
	for _, test := range []struct {
		source string
		mapped bool
		want   int32
	}{{"aaa more", false, 0}, {"zzz", false, 4}, {"ccc", false, 2}, {"ddd", false, 3}, {"CCC", false, -1}, {"CCC", true, 2}, {"short", false, -1}, {"a", true, -1}} {
		got, err := table.Lookup(context.Background(), []byte(test.source), test.mapped, func(context.Context, []byte, int32) (int16, error) { return 1, nil })
		if err != nil || got != test.want {
			t.Fatal(test, got, err)
		}
	}
	calls := 0
	got, err := table.Lookup(context.Background(), []byte("ccc"), false, func(context.Context, []byte, int32) (int16, error) { calls++; return 0, nil })
	if err != nil || got != -1 || calls != 3 {
		t.Fatal(got, calls, err)
	} // midpoint, repeated lower endpoint, then final lower endpoint
	if got, err := (*Paul2013ScannerLookupIndex)(nil).Lookup(context.Background(), nil, false, nil); err != nil || got != -1 {
		t.Fatal(got, err)
	}
	if _, err := table.Lookup(context.Background(), []byte("ccc"), false, nil); err == nil {
		t.Fatal("missing reached gate accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := table.Lookup(ctx, []byte("ccc"), false, nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestScannerLookupBoundaryDecisions(t *testing.T) {
	for _, test := range []struct {
		source, text    string
		length, advance int32
		want            int16
	}{{"abc def", "", 3, 0, 1}, {"abc", "abc", 3, 3, 1}, {"abc.more", "abc.more", 3, 8, 1}, {"abc-more", "abc-more", 3, 8, 1}, {"abc's", "abc's", 3, 5, 1}, {"abc've", "abc've", 3, 6, 1}, {"abc'xyz", "abc'xyz", 3, 7, 0}, {"abcdef", "abcdef", 3, 6, 0}, {"abc.more", "abc.more", 4, 8, 0}} {
		scan := func(_ context.Context, _ []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
			if mode != 0x17 || pointer != 0 || offset != 0 {
				t.Fatal(offset, mode, pointer)
			}
			return Paul2013ModelScannerResult{Advance: test.advance, Text: []byte(test.text), TextLength: int32(len(test.text)), Status: 1}, nil
		}
		got, err := CheckPaul2013ScannerLookupBoundary(context.Background(), []byte(test.source), test.length, scan)
		if err != nil || got != test.want {
			t.Fatal(test, got, err)
		}
	}
	for _, suffix := range []string{"'s", "'D", "'m", "'ve", "'EM", "'re", "'ll"} {
		if !Paul2013ScannerApostropheBoundary([]byte(suffix)) {
			t.Fatal(suffix)
		}
	}
	for _, suffix := range []string{"'", "'n't", "'sX", "x's"} {
		if Paul2013ScannerApostropheBoundary([]byte(suffix)) {
			t.Fatal(suffix)
		}
	}
}

func TestIndexedScannerLookupBindsPostTokenFlag(t *testing.T) {
	index := &Paul2013ScannerLookupIndex{Lower: []int32{-1, -1, -1, 0}, Upper: []int32{-1, -1, -1, 0}, Keys: [][]byte{[]byte("Word")}}
	boundaryScan := NewPaul2013ModelScannerWithASCII(nil, nil)
	lookup := NewPaul2013ScannerIndexedTokenLookup(func(pointer uint32) (*Paul2013ScannerLookupIndex, *Paul2013ScannerLookupIndex, error) {
		if pointer != 123 {
			t.Fatal(pointer)
		}
		return index, index, nil
	}, boundaryScan)
	scan := NewPaul2013ModelScannerWithASCII(nil, lookup)
	got, err := scan(context.Background(), []byte("word\x00"), 0, 0, 123)
	if err != nil || got.Field18 != 0x10000 {
		t.Fatal(got, err)
	}
}
