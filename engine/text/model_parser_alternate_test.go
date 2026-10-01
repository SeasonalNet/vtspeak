package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"testing"
)

// This scanner is a fixture for caller control flow, not a native scanner port.
func alternateScannerFixture(_ context.Context, source []byte, offset, mode int32, modePointer uint32) (Paul2013ModelScannerResult, error) {
	prefix := 0
	for prefix < len(source) && (source[prefix] == ' ' || source[prefix] == '\t' || source[prefix] == '\n' || source[prefix] == '\r') {
		prefix++
	}
	result := Paul2013ModelScannerResult{First: uint32(offset) + uint32(prefix)}
	if source[prefix] == 0 {
		return result, nil
	}
	end := prefix + 1
	result.Status = 3
	if !bytes.ContainsAny(source[prefix:prefix+1], "<>[]") {
		result.Status = 1
		for end < len(source) && source[end] != 0 && !bytes.ContainsAny(source[end:end+1], " <>[]\t\n\r") {
			end++
		}
	}
	result.Text = append([]byte(nil), source[prefix:end]...)
	result.TextLength = int32(end - prefix)
	result.Advance = int32(end)
	result.Second = uint32(offset) + uint32(end)
	return result, nil
}

func TestAlternateParserOrdinaryAngleAndBracketRows(t *testing.T) {
	state := bytes.Repeat([]byte{0xa5}, 0x3a00)
	binary.LittleEndian.PutUint16(state[:2], 0)
	binary.LittleEndian.PutUint32(state[0x39e8:], 0x12345678)
	before := append([]byte(nil), state...)
	source := []byte(" \tword <say more> [skip this] end\x00")
	got, err := RunPaul2013AlternateModelParser(context.Background(), source, state, func(ctx context.Context, input []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		if mode != 0 || pointer != 0x12345678 {
			t.Fatal("wrong native scanner controls")
		}
		return alternateScannerFixture(ctx, input, offset, mode, pointer)
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ReturnValue != 1 || binary.LittleEndian.Uint16(got.State[:2]) != 4 || got.ConsumedBytes != len(source)-1 || !bytes.Equal(before, state) {
		t.Fatalf("alternate loop result %+v", got)
	}
	for index, test := range []struct {
		text       string
		class      byte
		start, end uint32
	}{{"word", 'D', 2, 6}, {"say", 'S', 8, 11}, {"more", 'S', 12, 16}, {"end", 'D', 30, 33}} {
		row := got.State[0x14+index*0x94:]
		if string(cString(row[0x34:])) != test.text || row[0x23] != 'A' || row[0x24] != test.class || row[0x30] != 0xff || binary.LittleEndian.Uint32(row[0x28:]) != 0x15 || binary.LittleEndian.Uint32(row[:4]) != test.start || binary.LittleEndian.Uint32(row[4:8]) != test.end {
			t.Fatalf("row %d = text %q coord %d/%d", index, cString(row[0x34:]), binary.LittleEndian.Uint32(row[:4]), binary.LittleEndian.Uint32(row[4:8]))
		}
		if binary.LittleEndian.Uint32(row[0x2c:]) != 0xa5a5a5a5 {
			t.Fatal("changed unproduced row type")
		}
	}
}

func TestAlternateParserTerminalInvalidAndFullRows(t *testing.T) {
	for _, test := range []struct {
		input string
		count uint16
		want  int16
	}{{" \n\t\x00", 0, -1}, {"\x00", 1, 1}, {"[skip]\x00", 0, -1}, {"<say>\x00", 0, 1}, {"]\x00", 0, -1}, {"word\x00", 100, 0}} {
		state := make([]byte, 0x3a00)
		binary.LittleEndian.PutUint16(state[:2], test.count)
		got, err := RunPaul2013AlternateModelParser(context.Background(), []byte(test.input), state, alternateScannerFixture)
		if err != nil || got.ReturnValue != test.want {
			t.Fatalf("%q count %d -> %d, %v", test.input, test.count, got.ReturnValue, err)
		}
	}
	state := make([]byte, 0x3a00)
	got, err := RunPaul2013AlternateModelParser(context.Background(), []byte("12\x00"), state, func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
		return Paul2013ModelScannerResult{Status: 2, Advance: 2, TextLength: 2, Text: []byte("12")}, nil
	})
	if err != nil || got.ReturnValue != -1 {
		t.Fatal("accepted unsupported status")
	}
}

func TestAlternateParserErrorsAndDispatcherSelection(t *testing.T) {
	state := make([]byte, 0x3a00)
	for _, input := range [][]byte{[]byte("missing terminator"), []byte("<\x00"), []byte("[\x00")} {
		if _, err := RunPaul2013AlternateModelParser(context.Background(), input, state, alternateScannerFixture); err == nil {
			t.Fatal("accepted invalid input or nonprogressing unterminated group")
		}
	}
	if _, err := RunPaul2013AlternateModelParser(context.Background(), []byte{0}, state, nil); err == nil {
		t.Fatal("accepted missing scanner")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RunPaul2013AlternateModelParser(ctx, []byte{0}, state, alternateScannerFixture); !errors.Is(err, context.Canceled) {
		t.Fatal("ignored cancellation")
	}
	binary.LittleEndian.PutUint32(state[0x39e8:], 1)
	finalized := false
	got, err := DispatchPaul2013ModelParserWithAlternateScanner(context.Background(), []byte("word\x00"), state, make([]byte, 10), Paul2013ModelParserCallbacks{Finalize: func(state []byte, result int32) (int16, error) {
		finalized = true
		if result != 1 || binary.LittleEndian.Uint16(state[:2]) != 1 {
			t.Fatal("alternate parser did not reach finalizer")
		}
		return -1, nil
	}}, alternateScannerFixture)
	if err != nil || !finalized || got.Kind != Paul2013AlternateModelParser || got.FinalResult != -1 || !got.Success {
		t.Fatalf("dispatcher = %+v, %v", got, err)
	}
}
