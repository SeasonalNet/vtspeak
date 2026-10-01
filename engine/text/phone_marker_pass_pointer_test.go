package text

import (
	"encoding/binary"
	"testing"

	"vtspeak/engine/internal/paul2013tables"
)

func TestPhoneMarkerPrepassResolvesNativeTerminalPointerLazily(t *testing.T) {
	const pointer = uint32(0x12345678)
	current, following := phoneMarkerPointerFixture(pointer)
	followingPhone := following[paul2013MarkerPhoneCodeOffset]
	var characterMap [256]byte
	for value := range characterMap {
		characterMap[value] = byte(value)
	}
	called := 0
	got, err := ApplyPaul2013PhoneMarkerPrepassWithPointerResolver(
		[][]byte{current, following},
		func(address uint32) ([]byte, error) {
			called++
			if address != pointer {
				t.Fatalf("resolved address = %#x, want %#x", address, pointer)
			}
			return []byte("ya\x00"), nil
		},
		characterMap,
	)
	if err != nil {
		t.Fatal(err)
	}
	if called != 1 {
		t.Fatalf("pointer resolver called %d times, want once", called)
	}
	if got[0][paul2013MarkerPhoneCodeOffset] != 'P' {
		t.Fatalf("terminal phone = %#x, want 'P' after matching key", got[0][paul2013MarkerPhoneCodeOffset])
	}
	if got[1][paul2013MarkerPhoneCodeOffset] != followingPhone {
		t.Fatalf("following phone changed to %#x, want %#x", got[1][paul2013MarkerPhoneCodeOffset], followingPhone)
	}
}

func TestPhoneMarkerPrepassDoesNotResolveUnusedNativeTerminalPointer(t *testing.T) {
	current, following := phoneMarkerPointerFixture(0x87654321)
	current[paul2013MarkerTerminalOffset] = '.'
	got, err := ApplyPaul2013PhoneMarkerPrepassWithPointerResolver(
		[][]byte{current, following}, nil, [256]byte{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got[0][paul2013MarkerPhoneCodeOffset] != '9' {
		t.Fatalf("terminal phone = %#x, want unchanged '9'", got[0][paul2013MarkerPhoneCodeOffset])
	}
}

func TestPhoneMarkerRecordGroupsCarryPointerResolvedPrepass(t *testing.T) {
	const pointer = uint32(0x10203040)
	current, following := phoneMarkerPointerFixture(pointer)
	arena := append(append([]byte(nil), current...), following...)
	var characterMap [256]byte
	for value := range characterMap {
		characterMap[value] = byte(value)
	}
	prepared, err := PreparePaul2013PhoneMarkerRecordGroupsFromArenaWithPointerResolver(
		arena, 2, func(address uint32) ([]byte, error) {
			if address != pointer {
				t.Fatalf("resolved address = %#x, want %#x", address, pointer)
			}
			return []byte("ya\x00"), nil
		}, characterMap,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(prepared.Descriptors) != 1 || prepared.Descriptors[0].RecordCount != 2 {
		t.Fatalf("prepared groups = %+v, want one two-record group", prepared.Descriptors)
	}
	if len(prepared.SelectionContexts) != 1 || len(prepared.SelectionContexts[0]) != 2 {
		t.Fatalf("selection contexts = %+v, want one group with two contexts", prepared.SelectionContexts)
	}
	if got := prepared.Records[0][paul2013MarkerPhoneCodeOffset]; got != 'P' {
		t.Fatalf("preprocessed terminal phone = %#x, want 'P'", got)
	}
	if len(prepared.ResetArena) != len(arena) {
		t.Fatalf("reset arena length = %d, want %d", len(prepared.ResetArena), len(arena))
	}
}

func phoneMarkerPointerFixture(pointer uint32) ([]byte, []byte) {
	current := make([]byte, paul2013MarkerRecordSize)
	following := make([]byte, paul2013MarkerRecordSize)
	current[paul2013MarkerPhoneCountOffset] = 1
	following[paul2013MarkerPhoneCountOffset] = 1
	current[paul2013MarkerPhoneCodeOffset] = '9'
	following[paul2013MarkerPhoneCodeOffset] = paul2013TestBoundaryPhoneCode()
	current[paul2013MarkerTerminalOffset] = ']'
	binary.LittleEndian.PutUint32(current[paul2013MarkerSpecialKeyPointerOffset:], pointer)
	return current, following
}

func paul2013TestBoundaryPhoneCode() byte {
	for value := 1; value < 256; value++ {
		phone := byte(value)
		if phone != 'C' && paul2013tables.ByteTable1007B9E0(phone) == 1 &&
			int8(paul2013tables.ByteTable1007BB70(phone)) >= 0 {
			return phone
		}
	}
	return 0
}
