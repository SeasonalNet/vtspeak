package text

import (
	"reflect"
	"testing"
)

func TestNormalizePaul2013PositionStateProducerSeriesCompactsAndFillsPairs(t *testing.T) {
	input := Paul2013PositionStateProducerSeries{
		Boundaries: []int32{10, 20, 30, 40, 50, 60},
		Values:     []int32{5, 6, -1, 9, -1, -1},
	}
	got, err := NormalizePaul2013PositionStateProducerSeries(input)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013PositionStateProducerSeries{
		Boundaries: []int32{10, 20, 30, 40, 50, 60},
		Values:     []int32{5, 6, 5, 9, 5, -1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalized producer series = %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(input.Values, []int32{5, 6, -1, 9, -1, -1}) {
		t.Fatalf("normalization mutated input values: %v", input.Values)
	}
}

func TestNormalizePaul2013PositionStateProducerSeriesDropsUnpairedBoundaries(t *testing.T) {
	got, err := NormalizePaul2013PositionStateProducerSeries(Paul2013PositionStateProducerSeries{
		Boundaries: []int32{10, 20, 30, 40, 50},
		Values:     []int32{5, -1, 2, -1, 9},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013PositionStateProducerSeries{
		Boundaries: []int32{10, 20, 30, 40},
		Values:     []int32{5, -1, 2, -1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalized producer series = %#v, want %#v", got, want)
	}
}

func TestNormalizePaul2013PositionStateProducerArraysRunsModesInOrder(t *testing.T) {
	input := [3]Paul2013PositionStateProducerSeries{
		{Boundaries: []int32{1, 2}, Values: []int32{3, -1}},
		{Boundaries: []int32{4}, Values: []int32{5}},
		{},
	}
	got, err := NormalizePaul2013PositionStateProducerArrays(input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got[0], input[0]) || len(got[1].Boundaries) != 0 || len(got[2].Values) != 0 {
		t.Fatalf("normalized mode arrays = %#v", got)
	}
}

func TestNormalizePaul2013PositionStateProducerSeriesRejectsMismatchedArrays(t *testing.T) {
	if _, err := NormalizePaul2013PositionStateProducerSeries(Paul2013PositionStateProducerSeries{
		Boundaries: []int32{1},
	}); err == nil {
		t.Fatal("mismatched boundary/value arrays were accepted")
	}
}

func TestApplyPaul2013PositionStateProgramWithProducerArrays(t *testing.T) {
	program := Paul2013PositionStateProgram{
		RowCount: 4,
		Ranges: map[byte]Paul2013PositionStateRanges{
			0: {
				Mode:    0,
				RowKeys: []int32{1, 2, 3, 4},
				Initial: []int32{0, 0, 0, 0},
				Minimum: 0,
				Maximum: 20,
			},
		},
	}
	got, err := ApplyPaul2013PositionStateProgramWithProducerArrays(
		program,
		[3]Paul2013PositionStateProducerSeries{
			{Boundaries: []int32{1, 3, 5}, Values: []int32{10, -1, 20}},
			{},
			{},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []int32{10, 10, 0, 0}
	if !reflect.DeepEqual(got.Arrays[0], want) {
		t.Fatalf("composed producer and range pass = %v, want %v", got.Arrays[0], want)
	}
	if len(program.Ranges[0].Boundaries) != 0 {
		t.Fatalf("composition mutated caller's range map: %#v", program.Ranges[0])
	}
}

func TestApplyPaul2013PositionStateProgramWithProducerArraysRequiresConsumers(t *testing.T) {
	_, err := ApplyPaul2013PositionStateProgramWithProducerArrays(
		Paul2013PositionStateProgram{},
		[3]Paul2013PositionStateProducerSeries{
			{Boundaries: []int32{1, 2}, Values: []int32{4, -1}},
			{},
			{},
		},
	)
	if err == nil {
		t.Fatal("producer series without its state-range consumer was accepted")
	}
}
