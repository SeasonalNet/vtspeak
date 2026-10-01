package text

import (
	"reflect"
	"testing"
)

func TestBuildPaul2013PositionIntervalsFromParserRows(t *testing.T) {
	rows := []Paul2013OrdinaryParserOffsetRow{
		{Start: 0, End: 5, RawTypeByte: 'A', Text: []byte("Hello")},
		{Start: 7, End: 12, RawTypeByte: 'A', Text: []byte("world")},
		{Start: 12, End: 12, RawTypeByte: 'D'},
	}
	got, err := BuildPaul2013PositionIntervalsFromParserRows(30, rows)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013PositionEventRange{
		{Minimum: 30, Maximum: 34},
		{Minimum: 37, Maximum: 41},
		{Minimum: 42, Maximum: 42},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parser row position intervals = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013PositionIntervalsFromParserRowsRejectsNativeLimitAndOverflow(t *testing.T) {
	tooMany := make([]Paul2013OrdinaryParserOffsetRow, paul2013ParserStateRowLimit+1)
	if _, err := BuildPaul2013PositionIntervalsFromParserRows(0, tooMany); err == nil {
		t.Fatal("parser rows above native capacity succeeded")
	}
	if _, err := BuildPaul2013PositionIntervalsFromParserRows(1<<31-1, []Paul2013OrdinaryParserOffsetRow{{Start: 1, End: 2}}); err == nil {
		t.Fatal("overflowing model-state interval succeeded")
	}
}
