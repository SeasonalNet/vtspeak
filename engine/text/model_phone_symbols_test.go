package text

import (
	"bytes"
	"reflect"
	"testing"
)

func TestModelPhoneSymbolsMatchAllCanonicalPhones(t *testing.T) {
	for symbol := byte(1); symbol <= 0x45; symbol++ {
		canonical, err := DecodeCMUPhones([]byte{symbol})
		if err != nil {
			t.Fatal(err)
		}
		model, err := DecodePaul2013ModelPhoneSymbols([]byte{symbol})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(canonical, model) {
			t.Fatalf("symbol %#x canonical %+v model %+v", symbol, canonical, model)
		}
	}
	alias, err := DecodePaul2013ModelPhoneSymbols([]byte{'M'})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(alias, []CMUPhone{{Label: "T"}}) {
		t.Fatalf("M alias = %+v", alias)
	}
	for _, symbol := range []byte{0, 0x7f, 0x80, 0xff} {
		if _, err := DecodePaul2013ModelPhoneSymbols([]byte{symbol}); err == nil {
			t.Fatalf("accepted unsupported symbol %#x", symbol)
		}
	}
}

func TestPreviousPhoneCategoryUsesNativeStressTable(t *testing.T) {
	record := make([]byte, paul2013MarkerRecordSize)
	for _, test := range []struct {
		symbol byte
		want   int16
	}{{1, 0}, {2, 1}, {3, 2}, {0x13, -1}, {'M', -1}} {
		record[0x2e8] = test.symbol
		if got := paul2013PreviousPhoneCategory(record, 0); got != test.want {
			t.Fatalf("symbol %#x category %d want %d", test.symbol, got, test.want)
		}
	}
}

func TestModelStructuralMAliasGroupsDurationAndPitch(t *testing.T) {
	record := make([]byte, paul2013MarkerRecordSize)
	record[0x95] = 2
	record[0x2e8], record[0x2e9] = 1, 'M'
	record[0x329], record[0x32a] = '0', '0'
	record[0x3bd] = 'Z'
	canonical := append([]byte(nil), record...)
	canonical[0x2e9] = 0x39 // T
	groups, err := BuildPaul2013ModelRecordPhoneGroupsFromSymbols(record)
	if err != nil {
		t.Fatal(err)
	}
	if groups.NativeGroupCount != 1 || !bytes.Equal(groups.PhoneLabels, []byte{2, 3}) {
		t.Fatalf("M groups = %+v", groups)
	}
	state := Paul2013ModelRecordDurationState{RecordIndex: 0}
	duration, err := BuildPaul2013ModelRecordDurationInputs(record, state)
	if err != nil {
		t.Fatal(err)
	}
	wantDuration, err := BuildPaul2013ModelRecordDurationInputs(canonical, state)
	if err != nil {
		t.Fatal(err)
	}
	pitch, err := BuildPaul2013ModelRecordPitchInputs(record, state)
	if err != nil {
		t.Fatal(err)
	}
	wantPitch, err := BuildPaul2013ModelRecordPitchInputs(canonical, state)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(duration, wantDuration) || !reflect.DeepEqual(pitch, wantPitch) {
		t.Fatal("M alias differs from native identity/stress equivalent")
	}
	input := finalizerModelStateFixture('U', 'Y', "\x01Mc\x00")
	finalized, err := FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(input, make([]byte, 0x4770b), 0, []int32{7}, func(rowIndex, rowOffset int) (uint32, error) { return uint32(0x1000 + rowOffset), nil })
	if err != nil {
		t.Fatal(err)
	}
	boundaries, err := finalized.RunTokenBoundaries(0x20000000, input.ParserStateOffset, nil, nil, []int32{-1}, []int32{4})
	if err != nil {
		t.Fatal(err)
	}
	output := boundaries.Final.Arena[0x64c:]
	if output[0x2e9] != 'M' || output[0x29f] != 1 || output[0x94] != 1 || boundaries.Final.State.Markers[0].Marker != '^' {
		t.Fatalf("M native fields not preserved: code %x flag %d groups %d marker %q", output[0x2e9], output[0x29f], output[0x94], boundaries.Final.State.Markers[0].Marker)
	}
}
