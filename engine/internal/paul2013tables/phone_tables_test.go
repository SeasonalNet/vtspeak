package paul2013tables

import (
	"debug/pe"
	"os"
	"testing"
)

func TestNativePhoneTableDistinctEntries(t *testing.T) {
	for _, entry := range []struct {
		code                    byte
		identity, vowel, stress byte
	}{{1, 1, 1, 0}, {2, 1, 1, 1}, {3, 1, 1, 2}, {0x13, 7, 0, 0}, {0x4d, 31, 0, 0}} {
		if ByteTable1007B6C0(entry.code) != entry.identity || ByteTable1007B9E0(entry.code) != entry.vowel || ByteTable1007BAA8(entry.code) != entry.stress {
			t.Fatalf("symbol %#x: identity/vowel/stress %d/%d/%d", entry.code, ByteTable1007B6C0(entry.code), ByteTable1007B9E0(entry.code), ByteTable1007BAA8(entry.code))
		}
	}
}

// This local-input check is skipped in vendor-free checkouts. All windows
// are compared through PE section mapping, not a presumed file/RVA equality.
func TestPhoneTablesMatchLocalReadOnlyDLL(t *testing.T) {
	file, err := pe.Open("../../../binary/vt_pau.dll")
	if os.IsNotExist(err) {
		t.Skip("local proprietary DLL is absent")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	header, ok := file.OptionalHeader.(*pe.OptionalHeader32)
	if !ok || header.ImageBase != 0x10000000 {
		t.Fatal("unexpected DLL PE image base")
	}
	for _, table := range []struct {
		va     uint32
		size   int
		lookup func(byte) byte
	}{{0x1007b6c0, 128, ByteTable1007B6C0}, {0x1007b9e0, 256, ByteTable1007B9E0}, {0x1007baa8, 256, ByteTable1007BAA8}, {0x1007bb70, 256, ByteTable1007BB70}} {
		rva := table.va - header.ImageBase
		var data []byte
		for _, section := range file.Sections {
			if rva >= section.VirtualAddress && uint64(rva-section.VirtualAddress)+uint64(table.size) <= uint64(section.Size) {
				raw, err := section.Data()
				if err != nil {
					t.Fatal(err)
				}
				start := int(rva - section.VirtualAddress)
				data = raw[start : start+table.size]
				break
			}
		}
		if len(data) != table.size {
			t.Fatalf("table VA %#x not mapped", table.va)
		}
		for index, want := range data {
			if got := table.lookup(byte(index)); got != want {
				t.Fatalf("VA %#x index %#x got %d want %d", table.va, index, got, want)
			}
		}
	}
}
