package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
)

var primaryHandlerOrder = []uint32{0x10051a00, 0x1004d990, 0x1003d290, 0x1005c890, 0x100344e0, 0x10036940, 0x10058490, 0x100570c0, 0x10052e80, 0x10053440, 0x1003e4a0, 0x10033760, 0x1003c1a0, 0x1003aad0, 0x100420f0, 0x10031d40, 0x10040ef0, 0x100544f0, 0x10063780, 0x1003a880}

func primaryHandlersFixture() map[uint32]Paul2013PrimaryParserHandler {
	result := make(map[uint32]Paul2013PrimaryParserHandler)
	for _, address := range append([]uint32{0x1005dd20, 0x1005deb0}, primaryHandlerOrder...) {
		value := int32(-1)
		if address == 0x10051a00 {
			value = 0
		}
		result[address] = func(context.Context, []byte, *Paul2013ModelScannerResult, []byte, int32, int16) (int32, error) {
			return value, nil
		}
	}
	return result
}

func appendPrimaryFixtureRow(state []byte, offset int32) error {
	updated, _, err := AppendPaul2013ModelSourceRow(state, Paul2013ModelSourceRowInput{First: uint32(offset), Second: uint32(offset + 1), Text: []byte("a"), Type: 'A', Class: 'A'})
	if err == nil {
		copy(state, updated)
	}
	return err
}

func TestPrimaryParserOrderedFallbackAndRollback(t *testing.T) {
	state := bytes.Repeat([]byte{0x55}, 0x3a00)
	binary.LittleEndian.PutUint16(state[:2], 0)
	binary.LittleEndian.PutUint16(state[0xc:], 0)
	before := append([]byte(nil), state...)
	handlers := primaryHandlersFixture()
	var calls []uint32
	for _, address := range primaryHandlerOrder {
		handlers[address] = func(_ context.Context, state []byte, _ *Paul2013ModelScannerResult, _ []byte, _ int32, _ int16) (int32, error) {
			calls = append(calls, address)
			if address != 0x10051a00 {
				if binary.LittleEndian.Uint16(state[:2]) != 0 || binary.LittleEndian.Uint32(state[0x20:]) != 0xffffffff || state[0x21+4] != 0 {
					t.Fatalf("handler %#x observed rejected candidate", address)
				}
			}
			if address == 0x1003a880 {
				return 0, nil
			}
			if err := appendPrimaryFixtureRow(state, 0); err != nil {
				return 0, err
			}
			if address == 0x10051a00 {
				return 0, nil
			}
			return -1, nil
		}
	}
	got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("a\x00"), state, func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
		return Paul2013ModelScannerResult{}, nil
	}, handlers)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReturnValue != -1 || !reflect.DeepEqual(calls, primaryHandlerOrder) || !bytes.Equal(before, state) || binary.LittleEndian.Uint32(got.State[0x20:]) != 0xffffffff {
		t.Fatalf("cascade = %v, result %d", calls, got.ReturnValue)
	}
}

func TestPrimaryParserAcceptedRowsAndScannerPrefix(t *testing.T) {
	state := make([]byte, 0x3a00)
	handlers := primaryHandlersFixture()
	handlers[0x100344e0] = func(_ context.Context, state []byte, _ *Paul2013ModelScannerResult, _ []byte, offset int32, _ int16) (int32, error) {
		if offset == 2 {
			return 0, nil
		}
		return 1, appendPrimaryFixtureRow(state, offset)
	}
	got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("ab\x00"), state, func(_ context.Context, _ []byte, offset, mode int32, _ uint32) (Paul2013ModelScannerResult, error) {
		return Paul2013ModelScannerResult{PrefixFlag: 1}, nil
	}, handlers)
	if err != nil {
		t.Fatal(err)
	}
	if got.ReturnValue != 1 || got.ScannerCalls != 3 || binary.LittleEndian.Uint16(got.State[:2]) != 2 || binary.LittleEndian.Uint32(got.State[4:8]) != 2 {
		t.Fatalf("accepted rows = %+v", got)
	}
	for index := 0; index < 2; index++ {
		if binary.LittleEndian.Uint32(got.State[0x14+index*0x94+0x2c:]) != 1 {
			t.Fatal("scanner prefix did not mark previous row")
		}
	}
}

func TestPrimaryParserMarkedBoundaryAndMerge(t *testing.T) {
	for _, merge := range []bool{false, true} {
		state := make([]byte, 0x3a00)
		handlers := primaryHandlersFixture()
		handlers[0x10051a00] = func(_ context.Context, state []byte, token *Paul2013ModelScannerResult, _ []byte, offset int32, _ int16) (int32, error) {
			if offset != 0 {
				return 0, nil
			}
			token.Field14 = 1
			if err := appendPrimaryFixtureRow(state, 0); err != nil {
				return 0, err
			}
			binary.LittleEndian.PutUint32(state[0x14+0x28:], 0x12)
			binary.LittleEndian.PutUint32(state[0x14+0x2c:], 2)
			return 1, nil
		}
		handlers[0x10033760] = func(_ context.Context, state []byte, _ *Paul2013ModelScannerResult, _ []byte, offset int32, marked int16) (int32, error) {
			if offset == 2 {
				return 0, nil
			}
			if marked != 1 {
				t.Fatal("lost marked count")
			}
			if err := appendPrimaryFixtureRow(state, offset); err != nil {
				return 0, err
			}
			if merge {
				binary.LittleEndian.PutUint32(state[0x14+0x94+0x28:], 0x12)
			}
			return 1, nil
		}
		got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("ab\x00"), state, func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
			return Paul2013ModelScannerResult{}, nil
		}, handlers)
		if err != nil {
			t.Fatal(err)
		}
		wantCount, wantCursor := uint16(1), uint32(1)
		if merge {
			wantCount, wantCursor = 2, 2
		}
		if binary.LittleEndian.Uint16(got.State[:2]) != wantCount || binary.LittleEndian.Uint32(got.State[4:8]) != wantCursor {
			t.Fatalf("merge %v count/cursor %d/%d", merge, binary.LittleEndian.Uint16(got.State[:2]), binary.LittleEndian.Uint32(got.State[4:8]))
		}
		if merge && (got.MarkedRows != 0 || binary.LittleEndian.Uint32(got.State[0x14+0x2c:]) != 0) {
			t.Fatal("merge failed to reset marked row")
		}
	}
}

func TestPrimaryParserModeAndTerminalBranches(t *testing.T) {
	for _, initialMode := range []uint16{0, 1} {
		state := make([]byte, 0x3a00)
		binary.LittleEndian.PutUint16(state[0xc:], initialMode)
		handlers := primaryHandlersFixture()
		address := uint32(0x1005dd20)
		if initialMode == 0 {
			address = 0x1005c890
		}
		handlers[address] = func(_ context.Context, state []byte, _ *Paul2013ModelScannerResult, _ []byte, offset int32, _ int16) (int32, error) {
			if err := appendPrimaryFixtureRow(state, offset); err != nil {
				return 0, err
			}
			binary.LittleEndian.PutUint16(state[0xc:], 1)
			return 1, nil
		}
		got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("ab\x00"), state, func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
			return Paul2013ModelScannerResult{}, nil
		}, handlers)
		if err != nil || got.MarkedRows != 1 || binary.LittleEndian.Uint16(got.State[:2]) != 1 || binary.LittleEndian.Uint32(got.State[4:8]) != 1 {
			t.Fatalf("mode %d result count/marked/cursor, error %v", initialMode, err)
		}
	}
	state := make([]byte, 0x3a00)
	handlers := primaryHandlersFixture()
	handlers[0x10051a00] = func(_ context.Context, _ []byte, _ *Paul2013ModelScannerResult, _ []byte, offset int32, _ int16) (int32, error) {
		if offset == 0 {
			return -1, nil
		}
		return 1, nil
	}
	got, err := RunPaul2013PrimaryModelParser(context.Background(), []byte("a\x00"), state, func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error) {
		return Paul2013ModelScannerResult{}, nil
	}, handlers)
	if err != nil || got.ReturnValue != -1 || got.ConsumedBytes != 1 {
		t.Fatal("negative terminal advance or empty marked exit differed")
	}
}

func TestPrimaryParserDependencyAndCancellationErrors(t *testing.T) {
	state := make([]byte, 0x3a00)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RunPaul2013PrimaryModelParser(ctx, []byte{0}, state, nil, nil); !errors.Is(err, context.Canceled) {
		t.Fatal("ignored cancellation")
	}
	if _, err := RunPaul2013PrimaryModelParser(context.Background(), []byte{0}, state, nil, nil); err == nil {
		t.Fatal("accepted missing scanner")
	}
	if _, err := RunPaul2013PrimaryModelParser(context.Background(), []byte{0}, state, alternateScannerFixture, nil); err == nil {
		t.Fatal("treated missing handler as miss")
	}
}

func TestPrimaryParserDispatcherUsesRecoveredLoop(t *testing.T) {
	state := make([]byte, 0x3a00)
	handlers := primaryHandlersFixture()
	handlers[0x1003a880] = func(context.Context, []byte, *Paul2013ModelScannerResult, []byte, int32, int16) (int32, error) {
		return 0, nil
	}
	got, err := DispatchPaul2013ModelParserWithHandlers(context.Background(), []byte{0}, state, nil, Paul2013ModelParserCallbacks{Finalize: func(state []byte, result int32) (int16, error) {
		if result != -1 || binary.LittleEndian.Uint32(state[0x20:]) != 0xffffffff {
			t.Fatal("primary parser finalization handoff differs")
		}
		return 1, nil
	}}, alternateScannerFixture, handlers)
	if err != nil || got.Kind != Paul2013PrimaryModelParser || !got.Success || got.ParserResult != -1 {
		t.Fatalf("dispatcher primary=%+v, %v", got, err)
	}
}
