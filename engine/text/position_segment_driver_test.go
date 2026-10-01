package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"testing"
)

func segmentDriverFixture() ([]byte, []byte, []byte, []byte) {
	workspace, shared, state := positionIndexWorkspaceFixture()
	workspace = append(workspace, make([]byte, 0x1312c8-len(workspace))...)
	binary.LittleEndian.PutUint32(workspace[8:12], 0x100000)
	contextBytes := make([]byte, 0x4cfc)
	for _, offset := range []int{0x4cf4, 0x4cf0, 0x4cf8} {
		binary.LittleEndian.PutUint32(contextBytes[offset:], 150)
	}
	return contextBytes, workspace, shared, state.Records.Bytes
}

func segmentParserFixture(count int16, consumed uint32) []byte {
	state := make([]byte, 0x14+max(0, int(count))*0x94)
	binary.LittleEndian.PutUint16(state[:2], uint16(count))
	binary.LittleEndian.PutUint32(state[4:8], consumed)
	for index := 0; index < int(count); index++ {
		row := state[0x14+index*0x94:]
		binary.LittleEndian.PutUint32(row[:4], uint32(index*4))
		binary.LittleEndian.PutUint32(row[4:8], uint32(index*4+3))
	}
	return state
}

func TestPositionSegmentDriverRetriesAndProjectsNativeRows(t *testing.T) {
	engineContext, workspace, shared, model := segmentDriverFixture()
	// Native position processing runs before these unrelated fields are valid.
	for index := 0; index < 2; index++ {
		binary.LittleEndian.PutUint16(model[0x64c+index*0x3c0+0x3ac:], 0xffff)
	}
	before := append([]byte(nil), workspace...)
	modelBefore := append([]byte(nil), model...)
	calls := 0
	var lastState []byte
	parser := func(_ context.Context, source uint32, prepared []byte) (Paul2013PositionSegmentParserOutput, error) {
		calls++
		wantAddress := uint32(0x100000)
		wantDefault := uint32(150)
		if calls > 1 {
			wantAddress += 3
			wantDefault = 175
		}
		if source != wantAddress || binary.LittleEndian.Uint32(prepared[0x1210fc:]) != wantDefault {
			t.Fatal("incorrect retry source/defaults")
		}
		if calls == 1 {
			// Parser workspace changes survive the retry; preparation repeats.
			prepared[0x100] = 0x42
			prepared[0x1210d7] = 1
			binary.LittleEndian.PutUint32(prepared[0x1210dc:], 175)
			return Paul2013PositionSegmentParserOutput{Workspace: prepared, State: segmentParserFixture(-1, 3)}, nil
		}
		if prepared[0x100] != 0x42 || binary.LittleEndian.Uint32(prepared[0x1210e8:]) != 175 {
			t.Fatal("retry lost parser workspace or override")
		}
		lastState = segmentParserFixture(2, 8)
		return Paul2013PositionSegmentParserOutput{State: lastState}, nil
	}
	got, err := RunPaul2013PositionSegmentDriver(context.Background(), engineContext, workspace, shared, model, parser, func(_ uint32, count int) ([]int32, error) {
		cells := make([]int32, count)
		for index := range cells {
			cells[index] = int32(index + 100)
		}
		return cells, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || got.ParserCalls != 2 || got.ReturnValue != 11 || got.EndOfSource || binary.LittleEndian.Uint32(got.Workspace[4:8]) != 11 || binary.LittleEndian.Uint16(got.ModelArena[2:4]) != 2 {
		t.Fatalf("retry result = %+v", got)
	}
	for index, want := range [][2]uint32{{103, 105}, {107, 109}} {
		row := got.ModelArena[0x64c+index*0x3c0:]
		if binary.LittleEndian.Uint32(row[:4]) != want[0] || binary.LittleEndian.Uint32(row[4:8]) != want[1] || binary.LittleEndian.Uint16(row[0x3ac:]) != 0xffff {
			t.Fatal("mapped intervals or untouched controls differ")
		}
	}
	if got.Position.Program.Arrays[0][0] != 175 || got.Position.Program.Arrays[1][0] != 150 || !bytes.Equal(workspace, before) || !bytes.Equal(model, modelBefore) {
		t.Fatal("state defaults or input preservation differ")
	}
	lastState[4] = 99
	if got.ParserState[4] != 8 {
		t.Fatal("retained parser output aliases callback storage")
	}
}

func TestPositionSegmentDriverCompletionAndNullSource(t *testing.T) {
	for _, count := range []int16{-1, 0, 2} {
		engineContext, workspace, _, model := segmentDriverFixture()
		binary.LittleEndian.PutUint32(workspace[4:8], 7)
		before := append([]byte(nil), workspace...)
		got, err := RunPaul2013PositionSegmentDriver(context.Background(), engineContext, workspace, nil, model, func(context.Context, uint32, []byte) (Paul2013PositionSegmentParserOutput, error) {
			return Paul2013PositionSegmentParserOutput{State: segmentParserFixture(count, 0)}, nil
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !got.EndOfSource || got.ReturnValue != 7 || got.ParserCalls != 1 || binary.LittleEndian.Uint32(got.Workspace[0x44:]) != 1 || binary.LittleEndian.Uint16(got.ModelArena[2:4]) != uint16(count) || !bytes.Equal(workspace, before) {
			t.Fatal("incorrect completion branch")
		}
		if count == 2 && (got.Position.Program.Arrays[0][0] != 150 || got.Position.Program.Arrays[7][1] != -1 || binary.LittleEndian.Uint32(got.ModelArena[0x64c:]) != 7) {
			t.Fatal("zero-byte branch omitted initial row writes")
		}
	}
	workspace := make([]byte, 12)
	got, err := RunPaul2013PositionSegmentDriver(context.Background(), nil, workspace, nil, nil, nil, nil)
	if err != nil || got.EndOfSource || got.ReturnValue != 0 || got.ParserCalls != 0 || !bytes.Equal(got.Workspace, workspace) {
		t.Fatal("null source did work")
	}
}

func TestPositionSegmentDriverCancellationAndErrors(t *testing.T) {
	engineContext, workspace, shared, model := segmentDriverFixture()
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	got, err := RunPaul2013PositionSegmentDriver(ctx, engineContext, workspace, shared, model, func(context.Context, uint32, []byte) (Paul2013PositionSegmentParserOutput, error) {
		calls++
		cancel()
		return Paul2013PositionSegmentParserOutput{State: segmentParserFixture(0, 4)}, nil
	}, nil)
	if !errors.Is(err, context.Canceled) || calls != 1 || got.ParserCalls != 1 || got.ReturnValue != 4 {
		t.Fatal("retry ignored cancellation or lost completed writes")
	}
	for _, kind := range []string{"missing parser", "parser error", "short header", "short rows", "oversized count", "workspace extent"} {
		t.Run(kind, func(t *testing.T) {
			parser := Paul2013PositionSegmentParser(func(context.Context, uint32, []byte) (Paul2013PositionSegmentParserOutput, error) {
				output := Paul2013PositionSegmentParserOutput{State: segmentParserFixture(2, 3)}
				switch kind {
				case "parser error":
					return output, errors.New("parser failed")
				case "short header":
					output.State = output.State[:7]
				case "short rows":
					output.State = output.State[:0x1b]
				case "oversized count":
					binary.LittleEndian.PutUint16(output.State[:2], 101)
				case "workspace extent":
					output.Workspace = make([]byte, 12)
				}
				return output, nil
			})
			if kind == "missing parser" {
				parser = nil
			}
			if _, err := RunPaul2013PositionSegmentDriver(context.Background(), engineContext, workspace, shared, model, parser, nil); err == nil {
				t.Fatal("accepted invalid segment")
			}
		})
	}
}
