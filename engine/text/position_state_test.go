package text

import (
	"reflect"
	"testing"
)

func TestApplyPaul2013PositionStateRanges(t *testing.T) {
	got, err := ApplyPaul2013PositionStateRanges(Paul2013PositionStateRanges{
		Mode:       0,
		RowKeys:    []int32{0, 2, 4, 6, 8, 10},
		Boundaries: []int32{2, 6, 10},
		Values:     []int32{5, 99, -2},
		Initial:    []int32{77, 77, 77, 77, 77, 77},
		Minimum:    0,
		Maximum:    10,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int32{77, 5, 5, 10, 10, 77}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("state ranges = %v, want %v", got, want)
	}
}

func TestApplyPaul2013PositionStateRangesCarriesAndPreservesModeSevenSentinel(t *testing.T) {
	got, err := ApplyPaul2013PositionStateRanges(Paul2013PositionStateRanges{
		Mode:       7,
		RowKeys:    []int32{1, 2, 3, 4},
		Boundaries: []int32{2, 4},
		Values:     []int32{-1, 22},
		Initial:    []int32{4, 5, 6, 7},
		Minimum:    0,
		Maximum:    10,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int32{4, 5, 6, 7}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sentinel state ranges = %v, want carried values %v", got, want)
	}
}

func TestApplyPaul2013PositionStateRangesPreservesCarriedModeSevenSentinel(t *testing.T) {
	got, err := ApplyPaul2013PositionStateRanges(Paul2013PositionStateRanges{
		Mode:       7,
		RowKeys:    []int32{0, 1, 2},
		Boundaries: []int32{0, 2},
		Values:     []int32{-2, -5},
		Initial:    []int32{-1, -1, 8},
		Minimum:    0,
		Maximum:    10,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int32{-1, -1, 8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("carried sentinel state ranges = %v, want %v", got, want)
	}
}

func TestApplyPaul2013PositionStateRangesRejectsUnsupportedInputs(t *testing.T) {
	for name, input := range map[string]Paul2013PositionStateRanges{
		"mode":           {Mode: 3},
		"row-count":      {Mode: 0, RowKeys: []int32{1}, Initial: nil},
		"value-count":    {Mode: 0, RowKeys: []int32{1}, Initial: []int32{0}, Boundaries: []int32{1, 2}},
		"bounds":         {Mode: 0, Minimum: 2, Maximum: 1},
		"row-order":      {Mode: 0, RowKeys: []int32{2, 1}, Initial: []int32{0, 0}},
		"boundary-order": {Mode: 0, Boundaries: []int32{3, 2}},
	} {
		if _, err := ApplyPaul2013PositionStateRanges(input); err == nil {
			t.Errorf("%s input unexpectedly succeeded", name)
		}
	}
}

func TestApplyPaul2013PositionEventValuesAccumulatesAndClampsSelectorThree(t *testing.T) {
	got, err := ApplyPaul2013PositionEventValues(Paul2013PositionEventValues{
		Mode: 3,
		Ranges: []Paul2013PositionEventRange{
			{Minimum: 0, Maximum: 1},
			{Minimum: 2, Maximum: 4},
		},
		Boundaries: []int32{0, 1, 2, 3, 4},
		Values:     []int32{8, 8, 7, -2, 5},
		Initial:    []int32{0, 0},
		Minimum:    0,
		Maximum:    10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Values, []int32{10, 10}) || got.NextBoundary != 5 || got.OverflowBoundary != 2 {
		t.Fatalf("selector 3 result = %#v, want values [10 10], cursors 5 and 2", got)
	}
}

func TestPositionEventUpperClampDoesNotReplayOnOverlappingRows(t *testing.T) {
	got, err := ApplyPaul2013PositionEventValues(Paul2013PositionEventValues{Mode: 3, Ranges: []Paul2013PositionEventRange{{Minimum: 0, Maximum: 1}, {Minimum: 0, Maximum: 1}}, Boundaries: []int32{0, 1}, Values: []int32{8, 8}, Initial: []int32{-1, -1}, Minimum: 0, Maximum: 10})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Values, []int32{10, -1}) || got.NextBoundary != 2 {
		t.Fatalf("upper-clamped event replayed: %+v", got)
	}
}

func TestPositionProgramTerminalStartsAfterEveryConsumedEvent(t *testing.T) {
	got, err := ApplyPaul2013PositionStateProgram(Paul2013PositionStateProgram{RowCount: 1, Events: map[byte]Paul2013PositionEventValues{3: {Mode: 3, Ranges: []Paul2013PositionEventRange{{Minimum: 0, Maximum: 3}}, Boundaries: []int32{0, 1, 2, 3, 4}, Values: []int32{8, 8, -2, 1, 7}, Minimum: 0, Maximum: 10}}, Terminal: &Paul2013PositionTerminalAccumulator{Enabled: true, MinimumBoundary: 0, MaximumBoundary: 5, Boundaries: []int32{0, 1, 2, 3, 4}, Values: []int32{8, 8, -2, 1, 7}, Minimum: 0, Maximum: 10}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Events[3].NextBoundary != 4 || got.Events[3].OverflowBoundary != 2 || got.Terminal.Value != 7 || got.Terminal.NextBoundary != 5 {
		t.Fatalf("wrong shared event cursor: %+v", got)
	}
}

func TestApplyPaul2013PositionEventValuesUsesLastMatchForSelectorFour(t *testing.T) {
	got, err := ApplyPaul2013PositionEventValues(Paul2013PositionEventValues{
		Mode: 4,
		Ranges: []Paul2013PositionEventRange{
			{Minimum: 0, Maximum: 1},
			{Minimum: 2, Maximum: 4},
		},
		Boundaries: []int32{0, 1, 2, 3, 4},
		Values:     []int32{10, 20, 30, 40, 50},
		Initial:    []int32{-1, -1},
		Minimum:    0,
		Maximum:    0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Values, []int32{20, 50}) || got.NextBoundary != 5 {
		t.Fatalf("selector 4 result = %#v, want values [20 50], cursor 5", got)
	}
}

func TestApplyPaul2013PositionEventValuesRejectsMalformedInputs(t *testing.T) {
	for name, input := range map[string]Paul2013PositionEventValues{
		"selector":       {Mode: 2},
		"row-count":      {Mode: 3, Ranges: []Paul2013PositionEventRange{{}}, Initial: nil},
		"value-count":    {Mode: 4, Boundaries: []int32{1}, Values: nil},
		"start-index":    {Mode: 3, StartBoundary: 1},
		"overflow-index": {Mode: 3, OverflowStart: 1},
		"boundary-order": {Mode: 3, Boundaries: []int32{2, 1}, Values: []int32{0, 0}},
	} {
		if _, err := ApplyPaul2013PositionEventValues(input); err == nil {
			t.Errorf("%s input unexpectedly succeeded", name)
		}
	}
}

func TestPositionEventReversedGapRetainsValueAndCursor(t *testing.T) {
	got, err := ApplyPaul2013PositionEventValues(Paul2013PositionEventValues{Mode: 4, Ranges: []Paul2013PositionEventRange{{Minimum: 2, Maximum: 1}}, Initial: []int32{9}, Boundaries: []int32{1, 2}, Values: []int32{0, 1}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Values[0] != 9 || got.NextBoundary != 0 {
		t.Fatalf("empty native gap consumed event: %+v", got)
	}
}

func TestEvaluatePaul2013PositionTerminalAccumulator(t *testing.T) {
	got, err := EvaluatePaul2013PositionTerminalAccumulator(Paul2013PositionTerminalAccumulator{
		Enabled:         true,
		MinimumBoundary: 2,
		MaximumBoundary: 4,
		Boundaries:      []int32{1, 2, 3, 4},
		Values:          []int32{100, 3, 20, 4},
		StartBoundary:   0,
		Minimum:         0,
		Maximum:         10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != 10 || got.NextBoundary != 4 {
		t.Fatalf("terminal accumulator = %#v, want clamped value 10 at cursor 4", got)
	}
}

func TestEvaluatePaul2013PositionTerminalAccumulatorHonorsGate(t *testing.T) {
	got, err := EvaluatePaul2013PositionTerminalAccumulator(Paul2013PositionTerminalAccumulator{
		Enabled: false, Boundaries: []int32{1}, Values: []int32{9}, StartBoundary: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != -1 || got.NextBoundary != 0 {
		t.Fatalf("disabled terminal accumulator = %#v, want sentinel and unchanged cursor", got)
	}
}

func TestPositionTerminalReversedIntervalIsEmpty(t *testing.T) {
	got, err := EvaluatePaul2013PositionTerminalAccumulator(Paul2013PositionTerminalAccumulator{Enabled: true, MinimumBoundary: 2, MaximumBoundary: 1, Boundaries: []int32{1, 2}, Values: []int32{7, 9}, Minimum: 0, Maximum: 10})
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != -1 || got.NextBoundary != 0 {
		t.Fatalf("reversed terminal interval consumed event: %+v", got)
	}
}
