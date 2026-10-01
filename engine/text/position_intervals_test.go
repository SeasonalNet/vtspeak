package text

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestBuildPaul2013PositionIntervals(t *testing.T) {
	got, err := BuildPaul2013PositionIntervals(100, []Paul2013PositionIntervalOffsets{
		{Start: 0, End: 1},
		{Start: 4, End: 9},
		{Start: 12, End: 12},
		{Start: 15, End: 10},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013PositionEventRange{
		{Minimum: 100, Maximum: 100},
		{Minimum: 104, Maximum: 108},
		{Minimum: 112, Maximum: 112},
		{Minimum: 115, Maximum: 115},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("position intervals = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013PositionIntervalsRejectsOverflow(t *testing.T) {
	if _, err := BuildPaul2013PositionIntervals(1<<31-1, []Paul2013PositionIntervalOffsets{{Start: 1}}); err == nil {
		t.Fatal("overflowing position interval succeeded")
	}
}

func TestBuildPaul2013OrdinaryWordParserOffsetSegmentsMatchesRuntimeTraces(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   [][]Paul2013PositionIntervalOffsets
	}{
		{
			name:   "repeated words",
			source: "Hello hello.",
			want: [][]Paul2013PositionIntervalOffsets{{
				{Start: 0, End: 5}, {Start: 6, End: 11},
			}},
		},
		{
			name:   "comma remains and terminal resets",
			source: "Hello, world! Hello?",
			want: [][]Paul2013PositionIntervalOffsets{
				{{Start: 0, End: 5}, {Start: 7, End: 12}},
				{{Start: 0, End: 5}},
			},
		},
		{
			name:   "period question exclamation resets",
			source: "Hello. World? Hello!",
			want: [][]Paul2013PositionIntervalOffsets{
				{{Start: 0, End: 5}},
				{{Start: 0, End: 5}},
				{{Start: 0, End: 5}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := BuildPaul2013OrdinaryWordParserOffsetSegments(test.source)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("parser offset segments = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestBuildPaul2013OrdinaryWordParserOffsetSegmentsRejectsUnrecoveredClasses(t *testing.T) {
	for _, source := range []string{"", "123", "can't", "hello; world", "hello?!"} {
		if source == "" {
			if got, err := BuildPaul2013OrdinaryWordParserOffsetSegments(source); err != nil || len(got) != 0 {
				t.Fatalf("empty parser source = %+v, %v; want no segments", got, err)
			}
			continue
		}
		if _, err := BuildPaul2013OrdinaryWordParserOffsetSegments(source); err == nil {
			t.Errorf("unsupported parser source %q succeeded", source)
		}
	}
}

func TestBuildPaul2013OrdinaryParserOffsetSegmentsMatchesRuntimeNumberRows(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   [][]Paul2013PositionIntervalOffsets
	}{
		{
			name:   "captured number widths",
			source: "1. 12? 123! 1234.",
			want: [][]Paul2013PositionIntervalOffsets{
				{{Start: 0, End: 1}},
				{{Start: 0, End: 2}},
				{{Start: 0, End: 3}, {Start: 0, End: 3}, {Start: 0, End: 3}, {Start: 0, End: 3}},
				{{Start: 0, End: 4}, {Start: 0, End: 4}, {Start: 0, End: 4}},
			},
		},
		{
			name:   "year in word context",
			source: "Hello 2024. Hello?",
			want: [][]Paul2013PositionIntervalOffsets{
				{
					{Start: 0, End: 5},
					{Start: 6, End: 10}, {Start: 6, End: 10}, {Start: 6, End: 10},
				},
				{{Start: 0, End: 5}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := BuildPaul2013OrdinaryParserOffsetSegments(test.source)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("parser offset segments = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestBuildPaul2013OrdinaryParserOffsetSegmentsRejectsUnsupportedClasses(t *testing.T) {
	for _, source := range []string{"can't", "hello; world", "hello?!", "555-123"} {
		if _, err := BuildPaul2013OrdinaryParserOffsetSegments(source); err == nil {
			t.Errorf("unsupported parser source %q succeeded", source)
		}
	}
}

func TestBuildPaul2013OrdinaryParserOffsetSegmentsMatchesNumericClassTraces(t *testing.T) {
	const source = "+12! 1.25! $5.00! 25%! 21st! 01/02/2024! 3:45 PM! 555-1234!"
	repeated := func(start, end int32, count int) []Paul2013PositionIntervalOffsets {
		rows := make([]Paul2013PositionIntervalOffsets, count)
		for index := range rows {
			rows[index] = Paul2013PositionIntervalOffsets{Start: start, End: end}
		}
		return rows
	}
	want := [][]Paul2013PositionIntervalOffsets{
		{{Start: 0, End: 1}, {Start: 1, End: 3}},
		repeated(0, 4, 4),
		repeated(0, 5, 2),
		{{Start: 0, End: 2}, {Start: 0, End: 2}, {Start: 2, End: 3}},
		repeated(0, 4, 2),
		repeated(0, 10, 5),
		{{Start: 0, End: 4}, {Start: 0, End: 4}, {Start: 0, End: 4}, {Start: 5, End: 7}},
		{
			{Start: 0, End: 3}, {Start: 0, End: 3}, {Start: 0, End: 3}, {Start: 0, End: 3},
			{Start: 3, End: 4},
			{Start: 4, End: 8}, {Start: 4, End: 8}, {Start: 4, End: 8},
		},
	}
	got, err := BuildPaul2013OrdinaryParserOffsetSegments(source)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parser offset segments = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013PositionIntervalsFromParserState(t *testing.T) {
	state := make([]byte, paul2013ParserStateRowsOffset+2*paul2013ParserStateRowStride)
	binary.LittleEndian.PutUint16(state[0:2], 2)
	row0 := paul2013ParserStateRowsOffset
	row1 := row0 + paul2013ParserStateRowStride
	binary.LittleEndian.PutUint32(state[row0:row0+4], 0xfffffffd)
	binary.LittleEndian.PutUint32(state[row0+4:row0+8], 2)
	binary.LittleEndian.PutUint32(state[row1:row1+4], 8)
	binary.LittleEndian.PutUint32(state[row1+4:row1+8], 8)

	got, err := BuildPaul2013PositionIntervalsFromParserState(100, state)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013PositionEventRange{
		{Minimum: 97, Maximum: 101},
		{Minimum: 108, Maximum: 108},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parser-state intervals = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013PositionIntervalsFromParserStateValidatesCountAndBounds(t *testing.T) {
	for name, state := range map[string][]byte{
		"short count":    {0},
		"negative count": {0xff, 0xff},
		"over limit":     {101, 0},
		"truncated row":  {1, 0},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := BuildPaul2013PositionIntervalsFromParserState(0, state); err == nil {
				t.Fatal("invalid parser state succeeded")
			}
		})
	}
	got, err := BuildPaul2013PositionIntervalsFromParserState(0, []byte{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("zero-row intervals = %+v, want empty", got)
	}
}

func TestMapPaul2013PositionStateIndexesClampsAndUsesSeparateTables(t *testing.T) {
	got, err := MapPaul2013PositionStateIndexes(
		[]int32{-1, 4},
		[]int32{2, 99},
		[]int32{10, 20, 30, 40, 50},
		[]int32{100, 200, 300, 400, 500},
		nil,
		5,
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013PositionEventRange{
		{Minimum: 10, Maximum: 300},
		{Minimum: 50, Maximum: 500},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mapped position indexes = %+v, want %+v", got, want)
	}
}

func TestMapPaul2013PositionStateIndexesAppliesFinalTable(t *testing.T) {
	got, err := MapPaul2013PositionStateIndexes(
		[]int32{1},
		[]int32{2},
		[]int32{3, 4},
		[]int32{5, 6, 7},
		[]int32{20, 21, 22, 23, 24, 25, 26, 27},
		0,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013PositionEventRange{{Minimum: 24, Maximum: 27}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("remapped position indexes = %+v, want %+v", got, want)
	}
	if _, err := MapPaul2013PositionStateIndexes(
		[]int32{2}, []int32{0}, []int32{3, 4}, []int32{5}, []int32{0, 1}, 0, true,
	); err == nil {
		t.Fatal("out-of-range remapped index succeeded")
	}
}
