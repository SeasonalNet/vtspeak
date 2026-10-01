package text

import (
	"context"
	"debug/pe"
	"encoding/binary"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestModelScannerASCIIFieldsAndNativeLimits(t *testing.T) {
	for _, test := range []struct {
		source, text            string
		mode                    int32
		prefix, advance, status int32
	}{{"  word next", "word", 0, 2, 6, 1}, {"\nWord.", "Word", 1, 1, 5, 1}, {"?next", "?", 1, 0, 1, 7}, {"??next", "??", 1, 0, 2, 8}, {"...next", "...", 1, 0, 3, 3}, {"?!", "?", 1, 0, 1, 7}, {";", ";", 1, 0, 1, 7}, {";;;", ";;;", 1, 0, 3, 8}, {"\"\"", "\"\"", 0, 0, 2, 3}, {"\\\\", "\\", 0, 0, 1, 3}, {"?", "?", 0x12, 0, 1, 3}, {".1", ".", 0x12, 0, 1, 3}, {"\x01\t word", "word", 0, 3, 7, 1}, {strings.Repeat("a", 35), strings.Repeat("a", 29), 0, 0, 29, 1}, {strings.Repeat("W", 12), strings.Repeat("W", 8), 0, 0, 8, 1}, {strings.Repeat("?", 35), strings.Repeat("?", 29), 0, 0, 29, 8}} {
		got, handled, err := ScanPaul2013ModelScannerASCII([]byte(test.source+"\x00"), 100, test.mode)
		want := Paul2013ModelScannerResult{Text: []byte(test.text), TextLength: int32(len(test.text)), Field0: test.prefix, First: uint32(100 + test.prefix), Second: uint32(100 + test.advance), Field14: test.advance, Advance: test.advance, Status: test.status}
		if strings.HasPrefix(test.source, "\n") {
			want.PrefixFlag = 1
		}
		if err != nil || !handled || !reflect.DeepEqual(got, want) {
			t.Fatalf("%q mode%d: %+v want %+v handled%v, %v", test.source, test.mode, got, want, handled, err)
		}
	}
	for _, source := range []string{"word\xff", "\xff"} {
		if _, handled, err := ScanPaul2013ModelScannerASCII([]byte(source+"\x00"), 0, 1); err != nil || handled {
			t.Fatalf("unported %q accepted: %v", source, err)
		}
	}
}

func TestASCIIScannerPostLookupAndFallback(t *testing.T) {
	calls := []uint32{}
	scan := NewPaul2013ModelScannerWithASCII(nil, func(_ context.Context, address, pointer uint32, source []byte) (int32, error) {
		calls = append(calls, address)
		if pointer != 123 || string(source) != "word" {
			t.Fatal(pointer, string(source))
		}
		if address == 0x1005e010 {
			return -1, nil
		}
		return 0, nil
	})
	got, err := scan(context.Background(), []byte("word\x00"), 0, 0, 123)
	if err != nil || got.Field18 != 0x10000 || !reflect.DeepEqual(calls, []uint32{0x1005e010, 0x1005e1a0}) {
		t.Fatal(got, calls, err)
	}
	missing := NewPaul2013ModelScannerWithASCII(nil, nil)
	if _, err := missing(context.Background(), []byte("word\x00"), 0, 0, 123); err == nil {
		t.Fatal("missing post-lookup accepted")
	}
	if _, err := missing(context.Background(), []byte("\xff\x00"), 0, 0x12, 0); err == nil {
		t.Fatal("missing non-ASCII scanner accepted")
	}
	if _, err := missing(context.Background(), []byte("\n\n\x00"), 0, 0, 123); err != nil {
		t.Fatal("early terminal path required post lookup", err)
	}
	fallback := NewPaul2013ModelScannerWithASCII(func(_ context.Context, source []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		if string(source) != " \xff\x00" || offset != 50 || mode != 9 || pointer != 123 {
			t.Fatal("fallback arguments changed")
		}
		return Paul2013ModelScannerResult{Advance: 99}, nil
	}, nil)
	if got, err := fallback(context.Background(), []byte(" \xff\x00"), 50, 9, 123); err != nil || got.Advance != 99 {
		t.Fatal(got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := scan(ctx, []byte("word\x00"), 0, 0, 0); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestRealASCIIScannerPunctuationModelAndTerminalHandler(t *testing.T) {
	if _, err := os.Stat("../../data-common/dict-eng/sbd.tree3"); os.IsNotExist(err) {
		t.Skip("local proprietary shared data absent")
	}
	model, err := LoadPaul2013TerminalModel("../../data-common/dict-eng")
	if err != nil {
		t.Fatal(err)
	}
	scan := NewPaul2013ModelScannerWithASCII(nil, nil)
	full := []byte("hello. next word\x00")
	state := terminalStateFixture("hello", 0)
	binary.LittleEndian.PutUint32(state[0x39e8:], 0)
	recognize := NewPaul2013TerminalContextRecognizer(full, scan, model.Lookup)
	current, err := scan(context.Background(), full[5:], 5, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	classification, err := recognize(context.Background(), state, &current, full[5:], 5)
	if err != nil {
		t.Fatal(err)
	}
	if classification != 'P' && classification != 'N' {
		t.Fatal(classification)
	}
	handler := NewPaul2013TerminalParserHandler(scan, recognize)
	advance, err := handler(context.Background(), state, &current, full[5:], 5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if classification == 'P' {
		if advance != 1 || binary.LittleEndian.Uint32(state[0x14+0x2c:]) != 2 {
			t.Fatal(advance)
		}
	} else if advance != -1 {
		t.Fatal(advance)
	}
}

func TestScannerASCIICostsMatchLocalDLL(t *testing.T) {
	file, err := pe.Open("../../binary/vt_pau.dll")
	if os.IsNotExist(err) {
		t.Skip("local proprietary DLL absent")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	header, ok := file.OptionalHeader.(*pe.OptionalHeader32)
	if !ok || header.ImageBase != 0x10000000 {
		t.Fatal("unexpected PE image base")
	}
	rva := uint32(0x10077e58) - header.ImageBase
	for _, section := range file.Sections {
		if rva < section.VirtualAddress || uint64(rva-section.VirtualAddress)+252 > uint64(section.Size) {
			continue
		}
		raw, err := section.Data()
		if err != nil {
			t.Fatal(err)
		}
		start := int(rva - section.VirtualAddress)
		for index, want := range paul2013ScannerASCIICosts {
			if got := int16(binary.LittleEndian.Uint16(raw[start+index*10:])); got != want {
				t.Fatalf("scanner ASCII cost[%d]=%d want %d", index, got, want)
			}
		}
		return
	}
	t.Fatal("unmapped scanner ASCII cost table")
}

func TestScannerJoinLiteralsMatchLocalDLL(t *testing.T) {
	file, err := pe.Open("../../binary/vt_pau.dll")
	if os.IsNotExist(err) {
		t.Skip("local proprietary DLL absent")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	header, ok := file.OptionalHeader.(*pe.OptionalHeader32)
	if !ok || header.ImageBase != 0x10000000 {
		t.Fatal("unexpected PE image base")
	}
	for address, want := range map[uint32]string{
		0x1009c060: "int\x00", 0x1009c428: "wi\x00", 0x1009c418: "United\x00", 0x1009c410: "States\x00",
		0x10081090: "A.M\x00", 0x10081080: "P.M\x00", 0x1009bdfc: "pa\x00", 0x1008000c: "Mrs\x00",
		0x1009c3fc: "S.&P\x00", 0x1009c3f8: "Sec\x00", 0x10077764: "ca\x00", 0x1009c3e4: "OK\x00",
		0x100813bc: ".?!;\x00", 0x1009c404: ",:{}[]()<>\"\x00", 0x1009c3e8: ",:{}[]()<>\";/\x00",
	} {
		rva := address - header.ImageBase
		found := false
		for _, section := range file.Sections {
			if rva < section.VirtualAddress || uint64(rva-section.VirtualAddress)+uint64(len(want)) > uint64(section.Size) {
				continue
			}
			raw, err := section.Data()
			if err != nil {
				t.Fatal(err)
			}
			start := int(rva - section.VirtualAddress)
			if string(raw[start:start+len(want)]) != want {
				t.Fatalf("native join literal %#x mismatch", address)
			}
			found = true
			break
		}
		if !found {
			t.Fatalf("unmapped join literal %#x", address)
		}
	}
}
