package text

import (
	"context"
	"reflect"
	"testing"
)

func TestModelScannerTerminalPrefixNativeFields(t *testing.T) {
	for _, test := range []struct {
		source               string
		status, prefix, flag int32
		handled              bool
	}{{"\x00", 9, 0, 0, true}, {" \t\r\x00", 9, 3, 0, true}, {"\n\x00", 9, 1, 1, true}, {"\n\nword\x00", 8, 2, 0, true}, {" \n \nword\x00", 8, 4, 0, true}, {"\nword\x00", 0, 0, 0, false}, {"word\x00", 0, 0, 0, false}} {
		got, handled, err := ScanPaul2013ModelScannerTerminalPrefix([]byte(test.source), 100)
		if err != nil || handled != test.handled {
			t.Fatal("terminal gate differs")
		}
		if handled {
			want := Paul2013ModelScannerResult{Field0: test.prefix, PrefixFlag: test.flag, First: uint32(100 + test.prefix), Second: uint32(100 + test.prefix), Field14: test.prefix, Advance: test.prefix, Status: test.status}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("%q scanner %+v want %+v", test.source, got, want)
			}
		}
	}
	if _, _, err := ScanPaul2013ModelScannerTerminalPrefix([]byte("unterminated"), 0); err == nil {
		t.Fatal("accepted unterminated scanner source")
	}
}

func TestTerminalScannerRunsPrimaryHandlerWithoutFixtureScanner(t *testing.T) {
	scan := NewPaul2013ModelScannerWithTerminalPrefixes(nil)
	handlers := BindPaul2013TerminalParserHandler(nil, scan, nil)
	got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("\n\n\x00"), make([]byte, 0x3a00), scan, handlers)
	if err != nil || got.ReturnValue != -1 || got.ConsumedBytes != 2 || got.ScannerCalls != 2 {
		t.Fatalf("native terminal-only primary %+v, %v", got, err)
	}
	if _, err := scan(context.Background(), []byte("word\x00"), 0, 1, 0); err == nil {
		t.Fatal("fabricated ordinary scanner support")
	}
	calls := 0
	wrapped := NewPaul2013ModelScannerWithTerminalPrefixes(func(_ context.Context, input []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		calls++
		if string(input) != " word\x00" || offset != 7 || mode != 15 || pointer != 0x1234 {
			t.Fatal("fallback parameters changed")
		}
		return Paul2013ModelScannerResult{Status: 1}, nil
	})
	if got, err := wrapped(context.Background(), []byte(" word\x00"), 7, 15, 0x1234); err != nil || got.Status != 1 || calls != 1 {
		t.Fatal("fallback lost")
	}
}
