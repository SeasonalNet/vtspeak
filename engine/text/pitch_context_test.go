package text

import "testing"

func TestBuildPaul2013TokenPitchInputsFromText(t *testing.T) {
	sequence := LexicalPhoneSequence{
		Phones: []CMUPhone{{Label: "B"}, {Label: "AH", Stress: 1, Vowel: true}, {Label: "T"}},
		Tokens: []LexicalTokenSpan{{SourceSurface: "bat", Surface: "bat", PhoneStart: 0, PhoneEnd: 3}},
	}
	got, err := BuildPaul2013TokenPitchInputsFromText(sequence)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].Inputs) != 3 {
		t.Fatalf("pitch rows = %+v, want one token with three phone rows", got)
	}
	want := [][Paul2013PitchInputCount]int16{
		{1, 3, 1, 3, 1, 1, 0, 1, 1, 3, 3},
		{0, 3, 1, 1, 1, 1, 1, 1, 2, 3, 3},
		{1, 3, 1, 1, 3, 1, 1, 4, 3, 3, 3},
	}
	for phoneIndex, row := range got[0].Inputs {
		if row != want[phoneIndex] {
			t.Errorf("pitch row %d = %v, want %v", phoneIndex, row, want[phoneIndex])
		}
	}
}

func TestBuildPaul2013TokenPitchInputsWithPositionStates(t *testing.T) {
	sequence := LexicalPhoneSequence{
		Phones: []CMUPhone{{Label: "P"}, {Label: "AH", Stress: 0, Vowel: true}},
		Tokens: []LexicalTokenSpan{{SourceSurface: "pa", Surface: "pa", PhoneStart: 0, PhoneEnd: 2}},
	}
	got, err := BuildPaul2013TokenPitchInputsWithPositionStates(
		sequence, []byte{'0', '0'}, []byte{'Z'}, []uint8{1, 1},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].Inputs) != 2 {
		t.Fatalf("pitch inputs = %+v, want one token with two phone rows", got)
	}
	for phoneIndex, input := range got[0].Inputs {
		if input[8] != 1 {
			t.Errorf("phone %d position state = %d, want caller-supplied state 1", phoneIndex, input[8])
		}
	}
}

func TestPaul2013PitchPhoneClassMatchesRecoveredOrdinalTable(t *testing.T) {
	want := []int16{
		0, 0, 0, 0, 0, 0, 1, 3, 1, 2, 0, 6, 0, 2, 1, 2, 0, 0, 3, 1,
		5, 4, 4, 4, 0, 0, 1, 6, 2, 2, 1, 2, 0, 0, 2, 6, 6, 2, 2,
	}
	if len(paul2013PitchPhoneClass) != len(paul2013TreePhoneOrder) {
		t.Fatalf("pitch classes = %d, phone ordinals = %d", len(paul2013PitchPhoneClass), len(paul2013TreePhoneOrder))
	}
	for ordinal, label := range paul2013TreePhoneOrder {
		if got := paul2013PitchPhoneClass[label]; got != want[ordinal] {
			t.Errorf("pitch class for ordinal %d (%s) = %d, want %d", ordinal+1, label, got, want[ordinal])
		}
	}
}
