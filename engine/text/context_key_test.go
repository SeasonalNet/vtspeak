package text

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestPaul2013SelectionKeyMatchesDocumentedTables(t *testing.T) {
	dll, err := os.ReadFile(filepath.Join("..", "..", "binary", "vt_pau.dll"))
	if err != nil {
		t.Skip("local VoiceText DLL is unavailable")
	}
	const (
		firstMap  = 0x7b7ec
		secondMap = 0x7b788
		thirdMap  = 0x7b850
	)
	for name, source := range map[string]struct {
		offset int
		values [256]byte
	}{
		"first":  {firstMap, paul2013ContextKeyMap1},
		"second": {secondMap, paul2013ContextKeyMap2},
		"third":  {thirdMap, paul2013ContextKeyMap3},
	} {
		t.Run(name, func(t *testing.T) {
			if len(dll) < source.offset+len(source.values) {
				t.Fatal("DLL is shorter than the observed context-key table")
			}
			for i, value := range source.values {
				if value != dll[source.offset+i] {
					t.Fatalf("table byte %d = 0x%02x, want DLL byte 0x%02x", i, value, dll[source.offset+i])
				}
			}
		})
	}

	signature := [7]byte{0xa1, 0x01, 0x52, 0x33, 0xb2, 0x45, 0x7f}
	want := [5]byte{
		dll[firstMap+int(signature[1])],
		dll[secondMap+int(signature[2])],
		dll[thirdMap+int(signature[3])],
		signature[5],
		signature[6] & 0x20,
	}
	if got := (Context{Signature: signature}).Paul2013SelectionKey(); got != want {
		t.Fatalf("selection key = % x, want % x", got, want)
	}

	for name, source := range map[string]struct {
		offset int
		values [256]byte
	}{
		"class map":     {0x7b97c, paul2013ClassMap},
		"primary map":   {0x7b8b4, paul2013PrimaryCategoryMap},
		"secondary map": {0x7b918, paul2013SecondaryCategoryMap},
	} {
		t.Run(name, func(t *testing.T) {
			for i, value := range source.values {
				if value != dll[source.offset+i] {
					t.Fatalf("table byte %d = 0x%02x, want DLL byte 0x%02x", i, value, dll[source.offset+i])
				}
			}
		})
	}
	for name, source := range map[string]struct {
		offset int
		count  int
		weight func(int) (int32, bool)
	}{
		"key mismatch weights":  {0x7bf14, 5, Paul2013KeyMismatchWeight},
		"view mismatch weights": {0x7bf28, 10, Paul2013ViewMismatchWeight},
	} {
		t.Run(name, func(t *testing.T) {
			for i := 0; i < source.count; i++ {
				got, ok := source.weight(i)
				want := int32(binary.LittleEndian.Uint32(dll[source.offset+i*4:]))
				if !ok || got != want {
					t.Fatalf("weight %d = %d (valid=%t), want %d", i, got, ok, want)
				}
			}
		})
	}
}
