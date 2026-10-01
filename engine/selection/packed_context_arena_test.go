package selection

import (
	"encoding/binary"
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

func TestPhoneMarkerModelStateArenaCandidatesReadsCountAndRecordBase(t *testing.T) {
	const stateOffset = 0x80
	state := make([]byte, stateOffset+0x94)
	binary.LittleEndian.PutUint16(state[stateOffset:], 1)
	state[stateOffset+0x37] = 'U'
	state[stateOffset+0x38] = 'Y'
	binary.LittleEndian.PutUint16(state[stateOffset+0x40:], 2)
	table := make([]byte, 0x0b+0x70)
	copy(table[0x0b+0x1e:], "ABdMc\x00")
	arena := make([]byte, 0x4770b)
	built, err := text.FinalizePaul2013ModelParserRowsAndBuildModelStateRecords(
		text.Paul2013ModelParserFinalizerInput{
			StateArena: state, ParserStateOffset: stateOffset,
			ContextTable: table, ContextRowCount: 1,
		}, arena, 0, []int32{-1},
		func(int, int) (uint32, error) { return 0x10203040, nil },
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := LookupPaul2013PhoneMarkerModelStateArenaCandidatesWithPointerResolver(
		&voice.ClassCatalog{}, built.Records.Bytes,
		func(address uint32) ([]byte, error) {
			if address != 0x10203040 {
				t.Fatalf("resolved special-key address = %#x, want %#x", address, 0x10203040)
			}
			return []byte{0}, nil
		}, [256]byte{}, 100,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Prepared.Records) != 1 || len(got.Groups) != 1 {
		t.Fatalf("model-state candidate result has %d records and %d groups; want 1 and 1", len(got.Prepared.Records), len(got.Groups))
	}
	if len(got.Groups[0].Contexts) != 3 || len(got.Groups[0].CandidatePools) != 3 {
		t.Fatalf("model-state candidate group has %d contexts and %d pools; want 3 and 3", len(got.Groups[0].Contexts), len(got.Groups[0].CandidatePools))
	}
	if len(built.Records.ContextCodes) != 1 || string(built.Records.ContextCodes[0].PhoneMarkers) != "012" {
		t.Fatalf("model-state phone markers = %+v, want parsed preceding markers 012", built.Records.ContextCodes)
	}
}

func TestPhoneMarkerModelStateArenaCandidatesRejectsInvalidCountAndTruncation(t *testing.T) {
	catalog := &voice.ClassCatalog{}
	for _, test := range []struct {
		name  string
		count int16
		size  int
	}{
		{name: "negative count", count: -1, size: 4},
		{name: "record range truncated", count: 1, size: 0x64c},
	} {
		t.Run(test.name, func(t *testing.T) {
			arena := make([]byte, test.size)
			binary.LittleEndian.PutUint16(arena[2:4], uint16(test.count))
			if _, err := LookupPaul2013PhoneMarkerModelStateArenaCandidatesWithPointerResolver(
				catalog, arena, nil, [256]byte{}, 100,
			); err == nil {
				t.Fatal("model-state candidate lookup accepted malformed arena")
			}
		})
	}
}

func TestPhoneMarkerArenaPathFromTransitionStatesRequiresCatalog(t *testing.T) {
	_, err := SelectPaul2013PhoneMarkerArenaPathFromTransitionStates(
		nil, nil, nil, nil, nil, 0, nil, [256]byte{}, 0, nil, 0,
		nil, nil, nil, nil, nil, nil,
	)
	if err == nil {
		t.Fatal("arena path accepted a nil class catalog")
	}
}

func TestPhoneMarkerArenaPathWithPointerResolverFromTransitionStatesRequiresCatalog(t *testing.T) {
	_, err := SelectPaul2013PhoneMarkerArenaPathWithPointerResolverFromTransitionStates(
		nil, nil, nil, nil, nil, 0, nil, [256]byte{}, 0, nil, 0,
		nil, nil, nil, nil, nil, nil,
	)
	if err == nil {
		t.Fatal("pointer-resolved arena path accepted a nil class catalog")
	}
}
