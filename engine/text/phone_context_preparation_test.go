package text

import (
	"encoding/binary"
	"testing"
)

func TestDerivePaul2013PhoneContextPreparationFromRows(t *testing.T) {
	const parserRowSize = 0x94
	sourceRows := make([]byte, 2*parserRowSize)
	contextRows := make([]byte, Paul2013PhoneContextTableHeaderSize+2*Paul2013PhoneContextRowSize)
	binary.LittleEndian.PutUint16(contextRows, 2)
	firstRow := contextRows[Paul2013PhoneContextTableHeaderSize : Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize]
	secondRow := contextRows[Paul2013PhoneContextTableHeaderSize+Paul2013PhoneContextRowSize:]
	binary.LittleEndian.PutUint16(firstRow[Paul2013PhoneContextRowTokenIndex:], 0)
	binary.LittleEndian.PutUint16(secondRow[Paul2013PhoneContextRowTokenIndex:], 1)

	cases := []struct {
		name       string
		sourceKind byte
		sourceForm byte
		rowIndex   int
		repeated   bool
		wantCode   int16
		wantEarly  bool
	}{
		{name: "Y subtype returns one", sourceKind: 'U', sourceForm: 'Y', rowIndex: 0, wantCode: 1},
		{name: "S subtype returns one", sourceKind: 'U', sourceForm: 'S', rowIndex: 0, wantCode: 1},
		{name: "other subtype reaches zero tail", sourceKind: 'U', sourceForm: 'N', rowIndex: 0, wantCode: 0},
		{name: "non-U early return", sourceKind: 'X', sourceForm: 'S', rowIndex: 0, wantCode: 0, wantEarly: true},
		{name: "repeated source row early return", sourceKind: 'U', sourceForm: 'S', rowIndex: 1, repeated: true, wantCode: 0, wantEarly: true},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			clear(sourceRows)
			selectedIndex := 0
			if test.repeated {
				selectedIndex = 1
				binary.LittleEndian.PutUint16(firstRow[Paul2013PhoneContextRowTokenIndex:], 1)
			}
			selected := sourceRows[selectedIndex*parserRowSize : (selectedIndex+1)*parserRowSize]
			selected[0x23] = test.sourceKind
			selected[0x24] = test.sourceForm
			got, err := DerivePaul2013PhoneContextPreparationFromRows(sourceRows, contextRows, test.rowIndex)
			if err != nil {
				t.Fatal(err)
			}
			if !got.Known || got.Code != test.wantCode || got.EarlyReturn != test.wantEarly {
				t.Fatalf("preparation = %+v, want known=%t code=%d early=%t", got, true, test.wantCode, test.wantEarly)
			}
		})
	}
}

func TestObservedPaul2013DurationInputsMatchCapturedCmuPhoneControl(t *testing.T) {
	// The two nine-short inputs are the first values in the captured duration
	// tree calls in tools/revkit/work/stage10/feature-p-ah0-trees.log. Position
	// state is supplied explicitly because the VTML phoneme fixture's producer
	// differs from the ordinary lexical-token state path.
	phones := []CMUPhone{
		{Label: "P"},
		{Label: "AH", Stress: 0, Vowel: true},
	}
	firstMetadata := Paul2013DurationTreeMetadata{
		RecordByte1: 1, RecordByte2: 1, PositionState: 1, ContextCount: 1, AuxiliaryByte: 1,
	}
	secondMetadata := firstMetadata
	secondMetadata.AuxiliaryByte = 2
	previousBoundary := Paul2013DurationTreeNeighbor{BoundaryIdentityOrdinal: 40}
	first, err := BuildObservedPaul2013DurationTreeInput(
		previousBoundary, phones[0], Paul2013DurationTreeNeighbor{Phone: &phones[1]}, firstMetadata,
	)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildObservedPaul2013DurationTreeInput(
		Paul2013DurationTreeNeighbor{Phone: &phones[0]}, phones[1], previousBoundary, secondMetadata,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := [][9]int16{
		{27, 40, 3, 0, 1, 1, 1, 1, 1},
		{3, 27, 40, 0, 1, 1, 1, 1, 2},
	}
	got := [][9]int16{first, second}
	for phoneIndex := range want {
		if got[phoneIndex] != want[phoneIndex] {
			t.Errorf("phone %d duration input = %v, want captured vector %v", phoneIndex, got[phoneIndex], want[phoneIndex])
		}
	}
}

func TestTokenDurationInputsWithPositionStatesMatchCapturedCmuPhoneControl(t *testing.T) {
	sequence := LexicalPhoneSequence{
		Phones: []CMUPhone{
			{Label: "P"},
			{Label: "AH", Stress: 0, Vowel: true},
		},
		Tokens: []LexicalTokenSpan{{
			SourceSurface: "pa", Surface: "pa", PhoneStart: 0, PhoneEnd: 2,
		}},
	}
	got, err := BuildPaul2013TokenDurationInputsWithPositionStates(
		sequence, []byte{'0', '0'}, []byte{'Z'}, []uint8{1, 1},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := [][9]int16{
		{27, 40, 3, 0, 1, 1, 1, 1, 1},
		{3, 27, 40, 0, 1, 1, 1, 1, 2},
	}
	if len(got) != 1 || len(got[0].Inputs) != len(want) {
		t.Fatalf("duration inputs = %+v, want one token with %d phone inputs", got, len(want))
	}
	for phoneIndex := range want {
		if got[0].Inputs[phoneIndex] != want[phoneIndex] {
			t.Errorf("phone %d duration input = %v, want captured vector %v", phoneIndex, got[0].Inputs[phoneIndex], want[phoneIndex])
		}
	}
}

func TestTokenDurationInputsWithPositionStatesRejectsInvalidStates(t *testing.T) {
	sequence := LexicalPhoneSequence{
		Phones: []CMUPhone{{Label: "P"}},
		Tokens: []LexicalTokenSpan{{SourceSurface: "p", Surface: "p", PhoneStart: 0, PhoneEnd: 1}},
	}
	for _, test := range []struct {
		name   string
		states []uint8
	}{
		{name: "missing state"},
		{name: "too many states", states: []uint8{1, 2}},
		{name: "zero state", states: []uint8{0}},
		{name: "state above range", states: []uint8{4}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := BuildPaul2013TokenDurationInputsWithPositionStates(
				sequence, []byte{'0'}, []byte{'Z'}, test.states,
			); err == nil {
				t.Fatal("invalid position states were accepted")
			}
		})
	}
}
