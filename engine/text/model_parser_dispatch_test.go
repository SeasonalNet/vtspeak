package text

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestInitializePaul2013ModelParserState(t *testing.T) {
	const stateSize = 0x14 + 0xe74*4

	for _, control := range []uint16{0, 1, 2} {
		state := make([]byte, stateSize)
		for index := range state {
			state[index] = 0xa5
		}
		binary.LittleEndian.PutUint16(state[0x0c:0x0e], control)
		if err := InitializePaul2013ModelParserState(state); err != nil {
			t.Fatalf("initialize with control %d: %v", control, err)
		}
		if got := binary.LittleEndian.Uint32(state[4:8]); got != 0 {
			t.Errorf("control %d: counter = %#x, want zero", control, got)
		}
		if got := binary.LittleEndian.Uint32(state[0:4]); got != 0 {
			t.Errorf("control %d: leading counters = %#x, want zero", control, got)
		}
		wantControl := uint16(0)
		if control == 1 {
			wantControl = 1
		}
		if got := binary.LittleEndian.Uint16(state[0x0c:0x0e]); got != wantControl {
			t.Errorf("control %d: retained control = %d, want %d", control, got, wantControl)
		}
		for index, value := range state[paul2013ModelParserInitOffset:stateSize] {
			if value != 0 {
				t.Fatalf("control %d: byte %#x after clear is %#x, want zero", control, index+paul2013ModelParserInitOffset, value)
			}
		}
	}
	if err := InitializePaul2013ModelParserState(make([]byte, stateSize-1)); err == nil {
		t.Fatal("short parser state initialized without error")
	}
}

func TestPaul2013ModeProbeHasWord(t *testing.T) {
	if got, err := Paul2013ModeProbeHasWord(nil); err != nil || got {
		t.Fatalf("null probe = %t, %v; want false, nil", got, err)
	}
	probe := make([]byte, 10)
	if got, err := Paul2013ModeProbeHasWord(probe); err != nil || got {
		t.Fatalf("zero word probe = %t, %v; want false, nil", got, err)
	}
	binary.LittleEndian.PutUint16(probe[8:10], 0x0100)
	if got, err := Paul2013ModeProbeHasWord(probe); err != nil || !got {
		t.Fatalf("nonzero word probe = %t, %v; want true, nil", got, err)
	}
	if _, err := Paul2013ModeProbeHasWord(make([]byte, 9)); err == nil {
		t.Fatal("short nonnull probe accepted")
	}
}

func TestDispatchPaul2013ModelParser(t *testing.T) {
	errParser := errors.New("parser failure")
	tests := []struct {
		name        string
		modePointer uint32
		probeWord   uint16
		finalResult int16
		wantKind    Paul2013ModelParserKind
		wantSuccess bool
	}{
		{name: "zero mode pointer", wantKind: Paul2013PrimaryModelParser, finalResult: 1, wantSuccess: true},
		{name: "nonzero probe word", modePointer: 1, probeWord: 1, wantKind: Paul2013PrimaryModelParser, finalResult: -1, wantSuccess: true},
		{name: "zero probe word", modePointer: 1, wantKind: Paul2013AlternateModelParser, finalResult: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := make([]byte, paul2013ModelParserRequiredSize())
			binary.LittleEndian.PutUint32(state[paul2013ModelParserModeOffset:paul2013ModelParserModeOffset+4], test.modePointer)
			probe := make([]byte, 10)
			binary.LittleEndian.PutUint16(probe[8:10], test.probeWord)
			primaryCalls, alternateCalls, finalizeCalls := 0, 0, 0
			checkRows := func(gotState []byte) error {
				for row := 0; row < paul2013ModelParserRowCount; row++ {
					offset := paul2013ModelParserRowsOffset + row*paul2013ModelParserRowStride
					if got := binary.LittleEndian.Uint32(gotState[offset : offset+4]); got != paul2013ModelParserRowSentinel {
						return errors.New("record sentinel was not initialized before parser call")
					}
				}
				return nil
			}
			callbacks := Paul2013ModelParserCallbacks{
				Primary: func(_ []byte, gotState []byte) (int32, error) {
					primaryCalls++
					return 17, checkRows(gotState)
				},
				Alternate: func(_ []byte, gotState []byte) (int32, error) {
					alternateCalls++
					return 19, checkRows(gotState)
				},
				Finalize: func(_ []byte, parserResult int32) (int16, error) {
					finalizeCalls++
					if parserResult != 17 && parserResult != 19 {
						return 0, errors.New("unexpected parser result")
					}
					return test.finalResult, nil
				},
			}
			got, err := DispatchPaul2013ModelParser(nil, state, probe, callbacks)
			if err != nil {
				t.Fatalf("dispatch: %v", err)
			}
			if got.Kind != test.wantKind || got.Success != test.wantSuccess {
				t.Fatalf("dispatch result = %+v, want kind %d success %t", got, test.wantKind, test.wantSuccess)
			}
			wantPrimary, wantAlternate := 1, 0
			if test.wantKind == Paul2013AlternateModelParser {
				wantPrimary, wantAlternate = 0, 1
			}
			if primaryCalls != wantPrimary || alternateCalls != wantAlternate || finalizeCalls != 1 {
				t.Fatalf("callback counts primary=%d alternate=%d finalize=%d", primaryCalls, alternateCalls, finalizeCalls)
			}
		})
	}

	state := make([]byte, paul2013ModelParserRequiredSize())
	binary.LittleEndian.PutUint32(state[paul2013ModelParserModeOffset:paul2013ModelParserModeOffset+4], 1)
	callbacks := Paul2013ModelParserCallbacks{
		Primary:   func([]byte, []byte) (int32, error) { return 0, nil },
		Alternate: func([]byte, []byte) (int32, error) { return 0, nil },
		Finalize:  func([]byte, int32) (int16, error) { return 0, nil },
	}
	if _, err := DispatchPaul2013ModelParser(nil, state, nil, callbacks); err == nil {
		t.Fatal("nonzero mode pointer without pointed-to bytes accepted")
	}
	if _, err := DispatchPaul2013ModelParser(nil, state[:paul2013ModelParserRequiredSize()-1], nil, callbacks); err == nil {
		t.Fatal("short parser state accepted")
	}
	callbacks.Primary = func([]byte, []byte) (int32, error) { return 0, errParser }
	probe := make([]byte, 10)
	binary.LittleEndian.PutUint16(probe[8:10], 1)
	if _, err := DispatchPaul2013ModelParser(nil, state[:], probe, callbacks); err == nil {
		t.Fatal("parser error was dropped")
	}
}

func TestDispatchPaul2013ModelParserRequiresOnlySelectedParser(t *testing.T) {
	state := make([]byte, paul2013ModelParserRequiredSize())
	callbacks := Paul2013ModelParserCallbacks{
		Primary:  func([]byte, []byte) (int32, error) { return 1, nil },
		Finalize: func([]byte, int32) (int16, error) { return 1, nil },
	}
	if _, err := DispatchPaul2013ModelParser(nil, state, nil, callbacks); err != nil {
		t.Fatalf("primary dispatch required the unused alternate callback: %v", err)
	}

	state = make([]byte, paul2013ModelParserRequiredSize())
	binary.LittleEndian.PutUint32(state[paul2013ModelParserModeOffset:paul2013ModelParserModeOffset+4], 1)
	probe := make([]byte, 10)
	callbacks = Paul2013ModelParserCallbacks{
		Alternate: func([]byte, []byte) (int32, error) { return 2, nil },
		Finalize:  func([]byte, int32) (int16, error) { return 1, nil },
	}
	if _, err := DispatchPaul2013ModelParser(nil, state, probe, callbacks); err != nil {
		t.Fatalf("alternate dispatch required the unused primary callback: %v", err)
	}

	state = make([]byte, paul2013ModelParserRequiredSize())
	if _, err := DispatchPaul2013ModelParser(nil, state, nil, Paul2013ModelParserCallbacks{
		Finalize: func([]byte, int32) (int16, error) { return 1, nil },
	}); err == nil {
		t.Fatal("dispatch without its selected primary callback succeeded")
	}
}
