package text

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"
)

func terminalStateFixture(text string, rowType uint32) []byte {
	state := make([]byte, 0x3a00)
	if text != "" {
		binary.LittleEndian.PutUint16(state[:2], 1)
		row := state[0x14:]
		binary.LittleEndian.PutUint32(row[8:12], uint32(len(text)))
		binary.LittleEndian.PutUint32(row[4:8], 7)
		binary.LittleEndian.PutUint32(row[0x2c:], rowType)
		copy(row[0x34:], text)
	}
	binary.LittleEndian.PutUint32(state[0x39e8:], 0x1234)
	return state
}

func terminalScannerFixture(token Paul2013ModelScannerResult) Paul2013ModelScanner {
	return func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
		return token, nil
	}
}

func TestTerminalHandlerMultilineNativeReturnsAndPriorGate(t *testing.T) {
	for _, test := range []struct {
		text          string
		char          byte
		initial, want uint32
		advance       int32
	}{{"word", '?', 0, 3, 5}, {"word", '!', 0, 4, 5}, {"word", '.', 0, 2, 5}, {"word", '?', 9, 9, 5}, {"word.", '?', 0, 0, 5}, {"", '?', 0, 0, -5}} {
		state := terminalStateFixture(test.text, test.initial)
		token := Paul2013ModelScannerResult{Status: 99, Field14: 99}
		handler := NewPaul2013TerminalParserHandler(func(_ context.Context, source []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
			if mode != 1 || pointer != 0x1234 || offset != 4 || string(source) != "remaining\x00" {
				t.Fatal("terminal native rescan arguments differ")
			}
			return Paul2013ModelScannerResult{Status: 8, Text: []byte{test.char}, Field14: 5, Advance: 90}, nil
		}, nil)
		got, err := handler(context.Background(), state, &token, []byte("remaining\x00"), 4, 0)
		if err != nil || got != test.advance || binary.LittleEndian.Uint32(state[0x14+0x2c:]) != test.want || token.Status != 8 || token.Field14 != 5 {
			t.Fatalf("%+v -> advance %d, %v", test, got, err)
		}
	}
}

func TestTerminalHandlerPunctuationRecognitionAndFallbacks(t *testing.T) {
	for _, test := range []struct {
		text           string
		token          Paul2013ModelScannerResult
		classification int32
		flag           uint32
		wantAdvance    int32
		wantType       uint32
	}{
		{"word", Paul2013ModelScannerResult{Text: []byte("?"), TextLength: 1, Field14: 3}, 0x10050, 0, 3, 3},
		{"word.", Paul2013ModelScannerResult{Text: []byte("!"), TextLength: 1, Field14: 3}, 0x50, 0, 3, 4},
		{"", Paul2013ModelScannerResult{Text: []byte("..."), TextLength: 3, Field14: 3}, 0x50, 0, -3, 0},
		{"word", Paul2013ModelScannerResult{Text: []byte("."), TextLength: 1, Field14: 3}, 0, 0, -3, 0},
		{"word", Paul2013ModelScannerResult{Text: []byte("?"), TextLength: 1, Field14: 3, Field0: 2}, 0, 1, -3, 0},
		{"word", Paul2013ModelScannerResult{Text: []byte("?"), TextLength: 1, Field14: 3, Field0: 2}, 0, 0, 0, 0},
	} {
		state := terminalStateFixture(test.text, 0)
		binary.LittleEndian.PutUint32(state[0x14+0x28:], test.flag)
		calls := 0
		handler := NewPaul2013TerminalParserHandler(terminalScannerFixture(test.token), func(context.Context, []byte, *Paul2013ModelScannerResult, []byte, int32) (int32, error) {
			calls++
			return test.classification, nil
		})
		var token Paul2013ModelScannerResult
		got, err := handler(context.Background(), state, &token, []byte{0}, 5, 0)
		if err != nil || calls != 1 || got != test.wantAdvance || binary.LittleEndian.Uint32(state[0x14+0x2c:]) != test.wantType {
			t.Fatalf("punct %+v -> %d, %v", test, got, err)
		}
	}
}

func TestTerminalHandlerEndOfInputStructuralAndMiss(t *testing.T) {
	for _, test := range []struct {
		token Paul2013ModelScannerResult
		want  int32
	}{
		{Paul2013ModelScannerResult{Status: 9}, 1},
		{Paul2013ModelScannerResult{Status: 9, Field14: 5}, 5},
		{Paul2013ModelScannerResult{Status: 3, TextLength: 2, Text: []byte{0xa2, 0xfd}, Field14: 4}, 4},
		{Paul2013ModelScannerResult{Status: 3, TextLength: 2, Text: []byte{0xa2, 0xfe}, Field14: 4}, 0},
		{Paul2013ModelScannerResult{Status: 1, Text: []byte("word"), Field14: 5}, 0},
	} {
		for _, text := range []string{"", "word", "word."} {
			state := terminalStateFixture(text, 7)
			wantType := uint32(7)
			if text == "" {
				wantType = 0
			}
			handler := NewPaul2013TerminalParserHandler(terminalScannerFixture(test.token), nil)
			var token Paul2013ModelScannerResult
			got, err := handler(context.Background(), state, &token, []byte{0}, 0, 0)
			if err != nil || got != test.want || binary.LittleEndian.Uint32(state[0x14+0x2c:]) != wantType {
				t.Fatal("terminal completion/structural branch differs")
			}
		}
	}
}

func TestTerminalHandlerBindingErrorsAndPrimaryComposition(t *testing.T) {
	handlers := primaryHandlersFixture()
	delete(handlers, 0x10051a00)
	scanner := func(_ context.Context, _ []byte, offset, mode int32, _ uint32) (Paul2013ModelScannerResult, error) {
		if mode == 1 && offset == 1 {
			return Paul2013ModelScannerResult{Status: 9}, nil
		}
		return Paul2013ModelScannerResult{Status: 0, Text: []byte("a")}, nil
	}
	bound := BindPaul2013TerminalParserHandler(handlers, scanner, nil)
	if handlers[0x10051a00] != nil || bound[0x10051a00] == nil {
		t.Fatal("binding mutated source map or omitted terminal body")
	}
	bound[0x100344e0] = func(_ context.Context, state []byte, _ *Paul2013ModelScannerResult, _ []byte, offset int32, _ int16) (int32, error) {
		return 1, appendPrimaryFixtureRow(state, offset)
	}
	got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("a\x00"), make([]byte, 0x3a00), scanner, bound)
	if err != nil || got.ReturnValue != 1 || binary.LittleEndian.Uint16(got.State[:2]) != 1 || binary.LittleEndian.Uint32(got.State[4:8]) != 1 {
		t.Fatal("terminal body did not stop primary parser at native source end")
	}
	state := terminalStateFixture("word", 0)
	var token Paul2013ModelScannerResult
	if _, err := NewPaul2013TerminalParserHandler(terminalScannerFixture(Paul2013ModelScannerResult{Text: []byte("?")}), nil)(context.Background(), state, &token, []byte{0}, 0, 0); err == nil {
		t.Fatal("silently skipped punctuation recognizer")
	}
	expected := errors.New("scanner failed")
	if _, err := NewPaul2013TerminalParserHandler(func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
		return Paul2013ModelScannerResult{}, expected
	}, nil)(context.Background(), state, &token, []byte{0}, 0, 0); !errors.Is(err, expected) {
		t.Fatal("lost scanner error")
	}
}
