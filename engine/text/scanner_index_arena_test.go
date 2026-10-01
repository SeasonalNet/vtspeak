package text

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestScannerLookupArenaProjectionAndOwnership(t *testing.T) {
	memory := make([]byte, 0x900)
	put := func(offset int, value uint32) { binary.LittleEndian.PutUint32(memory[offset:], value) }
	put(0x10c, 0x200)
	put(0x110, 0x300)
	for _, address := range []int{0x200, 0x300} {
		put(address, 4)
		put(address+4, 0x400)
		put(address+8, 0x440)
		put(address+12, 0x500)
	}
	for index := 0; index < 4; index++ {
		put(0x400+index*4, 0xffffffff)
		put(0x440+index*4, 0xffffffff)
	}
	put(0x40c, 5)
	put(0x44c, 5)
	put(0x500+5*20+12, 0x800)
	copy(memory[0x800:], "Word\x00")
	reads := 0
	read := func(address uint32, size int) ([]byte, error) {
		reads++
		if int64(address)+int64(size) > int64(len(memory)) {
			return nil, fmt.Errorf("unmapped")
		}
		return memory[int(address) : int(address)+size], nil
	}
	stringsRead := 0
	resolve := func(address uint32) ([]byte, error) {
		stringsRead++
		if address != 0x800 {
			t.Fatal(address)
		}
		return memory[address:], nil
	}
	exact, mapped, err := ReadPaul2013ScannerLookupIndexes(0x100, read, resolve)
	if err != nil || reads != 9 || stringsRead != 2 {
		t.Fatal(reads, stringsRead, err)
	}
	memory[0x800] = 'X'
	put(0x40c, 0)
	gate := func(context.Context, []byte, int32) (int16, error) { return 1, nil }
	if got, err := exact.Lookup(context.Background(), []byte("Word"), false, gate); err != nil || got != 5 {
		t.Fatal(got, err)
	}
	if got, err := mapped.Lookup(context.Background(), []byte("word"), true, gate); err != nil || got != 5 {
		t.Fatal(got, err)
	}
}

func TestScannerLookupArenaInactiveAndErrors(t *testing.T) {
	read := func(address uint32, size int) ([]byte, error) {
		if address == 0x10c && size == 8 {
			raw := make([]byte, 8)
			binary.LittleEndian.PutUint32(raw, 0x200)
			return raw, nil
		}
		if address == 0x200 && size == 16 {
			return make([]byte, 16), nil
		}
		return nil, fmt.Errorf("unmapped")
	}
	exact, mapped, err := ReadPaul2013ScannerLookupIndexes(0x100, read, nil)
	if err != nil || exact == nil || mapped != nil || len(exact.Lower) != 0 {
		t.Fatal(exact, mapped, err)
	}
	if _, _, err := ReadPaul2013ScannerLookupIndexes(0xfffffff8, read, nil); err == nil {
		t.Fatal("model address overflow accepted")
	}
	if _, _, err := ReadPaul2013ScannerLookupIndexes(0x100, func(uint32, int) ([]byte, error) { return nil, nil }, nil); err == nil {
		t.Fatal("truncated memory span accepted")
	}
}
