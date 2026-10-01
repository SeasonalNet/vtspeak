package text

import (
	"reflect"
	"testing"
)

func TestParsePaul2013ContextCodesAppliesTrailingControls(t *testing.T) {
	got, err := ParsePaul2013ContextCodes([]byte("ABdMcE"))
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte("ABME"); !reflect.DeepEqual(got.Codes, want) {
		t.Fatalf("codes = %q, want %q", got.Codes, want)
	}
	if want := []byte{'0', '1', '2', '0'}; !reflect.DeepEqual(got.PhoneMarkers, want) {
		t.Fatalf("phone markers = %q, want %q", got.PhoneMarkers, want)
	}
	if want := []bool{false, false, true, false}; !reflect.DeepEqual(got.MFlags, want) {
		t.Fatalf("M flags = %v, want %v", got.MFlags, want)
	}
	if got.StoppedAtUnsupported {
		t.Fatal("supported codes unexpectedly stopped parsing")
	}
}

func TestParsePaul2013ContextCodesStopsAtUnsupportedAndNUL(t *testing.T) {
	got, err := ParsePaul2013ContextCodes([]byte{'A', 'c', 'F', 'B', 0, 'C'})
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{'A'}; !reflect.DeepEqual(got.Codes, want) {
		t.Fatalf("codes = %q, want %q", got.Codes, want)
	}
	if want := []byte{'2'}; !reflect.DeepEqual(got.PhoneMarkers, want) {
		t.Fatalf("phone markers = %q, want %q", got.PhoneMarkers, want)
	}
	if !got.StoppedAtUnsupported {
		t.Fatal("unsupported byte did not stop parsing")
	}

	got, err = ParsePaul2013ContextCodes([]byte{'A', 0, '!'})
	if err != nil {
		t.Fatal(err)
	}
	if got.StoppedAtUnsupported || !reflect.DeepEqual(got.Codes, []byte{'A'}) {
		t.Fatalf("NUL termination result = %+v", got)
	}
}

func TestParsePaul2013ContextCodesEnforcesNativeCapacityAndNonemptyOutput(t *testing.T) {
	input := make([]byte, paul2013ContextCodeCapacity+1)
	for index := range input {
		input[index] = 'A'
	}
	got, err := ParsePaul2013ContextCodes(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Codes) != paul2013ContextCodeCapacity {
		t.Fatalf("parsed %d codes, want %d", len(got.Codes), paul2013ContextCodeCapacity)
	}
	got, err = ParsePaul2013ContextCodes([]byte{0x20})
	if err != nil || !reflect.DeepEqual(got.Codes, []byte{0x20}) {
		t.Fatalf("native low-byte acceptance = %+v, %v", got, err)
	}
	if _, err := ParsePaul2013ContextCodes([]byte("dcd")); err == nil {
		t.Fatal("marker-only source produced a successful result")
	}
}

func TestSummarizePaul2013ContextCodeStateMapsFlagsAndMode(t *testing.T) {
	cases := []struct {
		name       string
		state      int16
		auxiliary  byte
		pitchCount int16
		wantMode   byte
	}{
		{name: "empty pitch terminal state", state: 3, wantMode: 6},
		{name: "nonempty pitch terminal state", state: 3, pitchCount: 1, wantMode: 7},
		{name: "special auxiliary state", state: 2, auxiliary: 12, wantMode: 7},
		{name: "ordinary state", state: 2, auxiliary: 11, wantMode: 5},
		{name: "state four", state: 4, wantMode: 7},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := SummarizePaul2013ContextCodeState(
				[]int32{-1, 0, 12}, test.state, test.auxiliary, test.pitchCount,
			)
			if err != nil {
				t.Fatal(err)
			}
			if want := []byte{0, 12, 12}; !reflect.DeepEqual(got.StateFlags, want) {
				t.Fatalf("state flags = %v, want %v", got.StateFlags, want)
			}
			if got.ModeCode != test.wantMode {
				t.Fatalf("mode code = %d, want %d", got.ModeCode, test.wantMode)
			}
		})
	}
}
