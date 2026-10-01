package selection

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestBuildPaul2013TransitionContextStatesFromRecordGroup(t *testing.T) {
	makeRecord := func(pointer uint32, state byte) []byte {
		record := make([]byte, 0x96+7)
		binary.LittleEndian.PutUint32(record[0x08:0x0c], pointer)
		record[0x92] = state
		record[0x94] = 1
		record[0x95] = 1
		return record
	}
	records := [][]byte{makeRecord(0x1001, 3), makeRecord(0x1002, 9)}
	arena := make([]byte, paul2013TransitionEligibilityArenaBase+4*paul2013TransitionEligibilityGroupStride*2)
	groupOffset := paul2013TransitionEligibilityArenaBase + 2*paul2013TransitionEligibilityGroupStride*2
	binary.LittleEndian.PutUint16(arena[groupOffset:], 0)
	secondRecordOffset := groupOffset + paul2013TransitionEligibilityGroupStride*2
	binary.LittleEndian.PutUint16(arena[secondRecordOffset:], 1)
	accepted := &Paul2013PositionQueryResult{
		WholePosition: Paul2013WholePositionQueryResult{Selection: Paul2013WholePositionResult{Accepted: true}},
		ReturnCount:   1,
	}
	var resolved []uint32
	states, err := BuildPaul2013TransitionContextStatesFromRecordGroup(
		records, 255, arena, 2, []*Paul2013PositionQueryResult{nil, accepted},
		func(pointer uint32) ([]byte, error) {
			resolved = append(resolved, pointer)
			rows := make([]byte, 0x1e)
			rows[0x1d] = 1
			return rows, nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(resolved, []uint32{0x1001, 0x1002}) {
		t.Fatalf("resolved pointers = %#v, want [0x1001 0x1002]", resolved)
	}
	want := []Paul2013TransitionContextState{{0, 0, 0, 9, 0, 1}}
	if !reflect.DeepEqual(states, want) {
		t.Fatalf("transition states = %#v, want %#v", states, want)
	}
}

func TestBuildPaul2013TransitionContextStatesFromRecordGroupRejectsMisalignedQueries(t *testing.T) {
	record := make([]byte, 0x96+7)
	record[0x95] = 1
	_, err := BuildPaul2013TransitionContextStatesFromRecordGroup(
		[][]byte{record}, 0, make([]byte, paul2013TransitionEligibilityArenaBase+2), 0,
		nil, func(uint32) ([]byte, error) { return nil, nil },
	)
	if err == nil {
		t.Fatal("misaligned query-result count unexpectedly succeeded")
	}
}

func TestPaul2013TransitionWeightRow(t *testing.T) {
	tests := []struct {
		name     string
		current  Paul2013TransitionContextState
		previous Paul2013TransitionContextState
		want     uint8
	}{
		{
			name:     "both state bytes zero and first bytes equal selects row zero",
			current:  Paul2013TransitionContextState{0, 1, 2, 3, 4, 0},
			previous: Paul2013TransitionContextState{0, 5, 6, 7, 8, 0},
			want:     0,
		},
		{
			name:     "both state bytes zero and first bytes differ selects row four",
			current:  Paul2013TransitionContextState{1, 0, 0, 0, 0, 0},
			previous: Paul2013TransitionContextState{2, 0, 0, 0, 0, 0},
			want:     4,
		},
		{
			name:     "zero current and positive previous selects row one",
			current:  Paul2013TransitionContextState{0, 0, 0, 0, 0, 0},
			previous: Paul2013TransitionContextState{0, 0, 0, 0, 0, 1},
			want:     1,
		},
		{
			name:     "positive current and zero previous selects row two",
			current:  Paul2013TransitionContextState{0, 0, 0, 0, 0, 1},
			previous: Paul2013TransitionContextState{0, 0, 0, 0, 0, 0},
			want:     2,
		},
		{
			name:     "nonzero state pair selects row three",
			current:  Paul2013TransitionContextState{0, 0, 0, 0, 0, 2},
			previous: Paul2013TransitionContextState{0, 0, 0, 0, 0, 1},
			want:     3,
		},
		{
			name:     "negative signed current selects row three",
			current:  Paul2013TransitionContextState{0, 0, 0, 0, 0, 0xff},
			previous: Paul2013TransitionContextState{0, 0, 0, 0, 0, 0},
			want:     3,
		},
		{
			name:     "negative signed previous with zero current selects row three",
			current:  Paul2013TransitionContextState{0, 0, 0, 0, 0, 0},
			previous: Paul2013TransitionContextState{0, 0, 0, 0, 0, 0xff},
			want:     3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Paul2013TransitionWeightRow(test.current, test.previous)
			if got != test.want {
				t.Fatalf("weight row = %d, want %d", got, test.want)
			}
		})
	}
}

func TestPaul2013TransitionWeightRowExhaustiveSignedBytePairs(t *testing.T) {
	for currentValue := 0; currentValue <= 255; currentValue++ {
		for previousValue := 0; previousValue <= 255; previousValue++ {
			for _, firstBytesDiffer := range []bool{false, true} {
				current := Paul2013TransitionContextState{0, 0x11, 0x22, 0x33, 0x44, byte(currentValue)}
				previous := Paul2013TransitionContextState{0, 0xaa, 0xbb, 0xcc, 0xdd, byte(previousValue)}
				if firstBytesDiffer {
					current[0] = 1
				}
				got := Paul2013TransitionWeightRow(current, previous)
				want := referencePaul2013TransitionWeightRow(current, previous)
				if got != want {
					t.Fatalf("current state %02x previous state %02x first-byte-diff=%v: row %d, want %d", currentValue, previousValue, firstBytesDiffer, got, want)
				}
			}
		}
	}
}

func referencePaul2013TransitionWeightRow(
	current Paul2013TransitionContextState,
	previous Paul2013TransitionContextState,
) uint8 {
	currentState := int8(current[5])
	previousState := int8(previous[5])
	switch {
	case currentState == 0 && previousState == 0 && current[0] == previous[0]:
		return 0
	case currentState == 0 && previousState > 0:
		return 1
	case currentState > 0 && previousState == 0:
		return 2
	case currentState == 0 && previousState == 0 && current[0] != previous[0]:
		return 4
	default:
		return 3
	}
}
