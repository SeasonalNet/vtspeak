package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
)

func TestBackwardContextNativeBounds(t *testing.T) {
	for _, test := range []struct {
		source   string
		count    int
		want     string
		consumed int16
	}{
		{"", 3, "", 0}, {" \t\n", 3, "", 0},
		{"abc", 3, "abc", 3}, {"abc \t", 3, "abc", 5},
		{"one two three four", 3, "two three four", 14},
		{"one  two", 3, "one  two", 8},
		{"one   two", 3, "  two", 5},
		{strings.Repeat("a", 40), 1, strings.Repeat("a", 28), 28},
	} {
		got, consumed, err := CopyPaul2013BackwardContext([]byte(test.source), len(test.source), test.count)
		if err != nil || string(got) != test.want+"\x00" || consumed != test.consumed {
			t.Fatalf("%q count %d -> %q/%d want %q/%d: %v", test.source, test.count, got, consumed, test.want, test.consumed, err)
		}
	}
	if _, _, err := CopyPaul2013BackwardContext(nil, 1, 3); err == nil {
		t.Fatal("invalid position accepted")
	}
}

// Fixture scanner deliberately makes its return differ from result +0x14.
// These tests validate native caller behavior, not ordinary scanner parity.
func fixtureTerminalWordScanner(ctx context.Context, source []byte, offset, mode int32, modePointer uint32) (Paul2013ModelScannerResult, error) {
	if err := ctx.Err(); err != nil {
		return Paul2013ModelScannerResult{}, err
	}
	prefix := ScanPaul2013ModelParserWhitespacePrefix(source).ConsumedBytes
	length := 0
	for prefix+length < len(source) && source[prefix+length] != 0 && source[prefix+length] != ' ' {
		length++
	}
	return Paul2013ModelScannerResult{Field0: int32(prefix), TextLength: int32(length), Field14: int32(prefix + length), Status: 1, Text: source[prefix : prefix+length], Advance: 999, PrefixFlag: 12, First: uint32(offset), Second: 33}, nil
}

func TestBackwardContextRollingScanner(t *testing.T) {
	records, accepted, err := ScanPaul2013BackwardContext(context.Background(), []byte("one two three four\x00"), 100, 1, fixtureTerminalWordScanner)
	if err != nil || !accepted {
		t.Fatal(accepted, err)
	}
	for index, want := range []string{"two", "three", "four"} {
		if string(records[index].Text) != want || records[index].Advance != 999 {
			t.Fatalf("record %d = %+v", index, records[index])
		}
	}
	zero := func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
		return Paul2013ModelScannerResult{}, nil
	}
	if _, _, err := ScanPaul2013BackwardContext(context.Background(), []byte("x\x00"), 0, 1, zero); err == nil {
		t.Fatal("nonprogress accepted")
	}
	if _, accepted, err := ScanPaul2013BackwardContext(context.Background(), []byte{0}, 0, 1, nil); err != nil || accepted {
		t.Fatal(accepted, err)
	}
}

func TestTerminalContextRecognizerProjectionAndModelFallback(t *testing.T) {
	full := []byte("one two word. next later\x00")
	position := int32(bytes.IndexByte(full, '.'))
	state := make([]byte, 0x39ec)
	binary.LittleEndian.PutUint32(state[0x39e8:], 123)
	current := Paul2013ModelScannerResult{Text: []byte("."), TextLength: 1, Field14: 1, Status: 3, Advance: 88}
	calls := 0
	scan := func(ctx context.Context, source []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		calls++
		if mode != 1 {
			t.Fatal("wrong mode")
		}
		if offset > position && pointer != 123 || offset < position && pointer != 0 {
			t.Fatal("wrong mode pointer")
		}
		return fixtureTerminalWordScanner(ctx, source, offset, mode, pointer)
	}
	for _, result := range []uint32{0, 1, 0x100, 0x101} {
		lookup := func(_ context.Context, window Paul2013TerminalContextWindow) (uint32, error) {
			for index, want := range []string{"one", "two", "word", ".", "next", "later"} {
				if string(window[index].Text) != want || window[index].Advance != 0 || window[index].Field14 != 0 || window[index].PrefixFlag != 0 || window[index].First != 0 || window[index].Second != 0 {
					t.Fatalf("projected record %d = %+v", index, window[index])
				}
			}
			return result, nil
		}
		calls = 0
		recognize := NewPaul2013TerminalContextRecognizer(full, scan, lookup)
		got, err := recognize(context.Background(), state, &current, full[position:], position)
		want := int32('N')
		if byte(result) != 0 {
			want = 'P'
		}
		if err != nil || got != want || calls != 5 {
			t.Fatalf("result %x got %c calls %d: %v", result, got, calls, err)
		}
	}
	recognize := NewPaul2013TerminalContextRecognizer(full, scan, nil)
	if _, err := recognize(context.Background(), state, &current, full[position:], position); err == nil {
		t.Fatal("missing model lookup accepted")
	}
	wantErr := errors.New("lookup failure")
	recognize = NewPaul2013TerminalContextRecognizer(full, scan, func(context.Context, Paul2013TerminalContextWindow) (uint32, error) { return 0, wantErr })
	if _, err := recognize(context.Background(), state, &current, full[position:], position); !errors.Is(err, wantErr) {
		t.Fatal(err)
	}
}

func TestTerminalContextRecognizerEndAndValidation(t *testing.T) {
	state := make([]byte, 0x39ec)
	current := Paul2013ModelScannerResult{Text: []byte("?"), TextLength: 1, Field14: 1, Status: 3}
	full := []byte("?\x00")
	calls := 0
	scan := func(ctx context.Context, source []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		calls++
		return fixtureTerminalWordScanner(ctx, source, offset, mode, pointer)
	}
	recognize := NewPaul2013TerminalContextRecognizer(full, scan, nil)
	full[0] = '!' // Constructor owns its copy.
	if got, err := recognize(context.Background(), state, &current, []byte("?\x00"), 0); err != nil || got != 'P' || calls != 1 {
		t.Fatal(got, calls, err)
	}
	if _, err := recognize(context.Background(), state, &current, full, 0); err == nil {
		t.Fatal("source mismatch accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := recognize(ctx, state, &current, []byte("?\x00"), 0); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	current.Field14 = 3
	if _, err := recognize(context.Background(), state, &current, []byte("?\x00"), 0); err == nil {
		t.Fatal("advance outside source accepted")
	}
}

func TestTerminalContextRecognizerThroughTerminalHandler(t *testing.T) {
	full := []byte("word?\x00")
	scan := func(ctx context.Context, source []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		if offset == 4 {
			return Paul2013ModelScannerResult{Text: []byte("?"), TextLength: 1, Field14: 1, Status: 3}, nil
		}
		return fixtureTerminalWordScanner(ctx, source, offset, mode, pointer)
	}
	state := terminalStateFixture("word", 0)
	handler := NewPaul2013TerminalParserHandler(scan, NewPaul2013TerminalContextRecognizer(full, scan, nil))
	var token Paul2013ModelScannerResult
	got, err := handler(context.Background(), state, &token, full[4:], 4, 0)
	if err != nil || got != 1 || binary.LittleEndian.Uint32(state[0x14+0x2c:]) != 3 {
		t.Fatalf("terminal handler/context result %d, %v", got, err)
	}
}
