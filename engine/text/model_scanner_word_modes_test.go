package text

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestScannerAdditionalWordModes(t *testing.T) {
	for _, test := range []struct {
		mode         int32
		source, want string
	}{
		{0x17, "word.next", "word.next"}, {0x17, "int'l", "int'l"}, {9, "word.next", "word"}, {9, "can't", "can"},
		{0x12, "word.next", "word"}, {0x18, "word-next", "word"}, {0x1c, "word'next", "word"},
		{0xc, "word'foo'bar", "word'foo'bar"}, {0xc, "word.next", "word"},
		{0x1a, "word'foo'bar", "word'foo'bar"}, {0x1a, "word-next", "word"},
		{0x16, "one-two'three.four", "one-two'three.four"}, {0x16, "word'-next", "word"},
		{0x11, "one-two.three", "one-two.three"}, {0x15, "one-two.three", "one-two.three"},
		{0x11, "A1B", "A1B"}, {0x15, "Z0ABC", "Z0ABC"},
		{0x11, "a1B", "a"}, {0x15, "A1b", "A"}, {0x11, "A12B", "A"},
		{0x15, "AB1C", "AB"}, {0x11, "A1", "A"}, {0x15, "can't", "can"},
		{0x11, "int'l", "int'l"}, {0x16, "int'l", "int'l"},
		{0x11, "United States", "United States"}, {0x15, "UNITED states9", "UNITED states"},
		{0x11, "United Statesman", "United"}, {0x15, "United State", "United"},
		{0x11, "United  States", "United"}, {0x15, "United\tStates", "United"},
		{0x19, "A1b", "A1b"}, {0x19, "A1", "A1"},
		{0x19, "a1B", "a"}, {0x19, "A1B2C", "A1B"}, {0x19, "word.Next", "word.Next"},
		{0x16, "one.two-" + strings.Repeat("W", 10), "one.two"},
		{5, "pa'anga", "pa'anga"}, {5, "PA'ANGAX", "PA'ANGAX"},
		{5, "pa'ang", "pa"}, {5, "pa'ango", "pa"}, {5, "word'anga", "word"},
		{5, "pa.next", "pa.next"}, {5, "int'l", "int'l"},
		{0xf, "A.M.", "A.M."}, {0xf, "p.m.!", "p.m."},
		{0xf, "a.M..", "a.M."}, {0xf, "A.M.next", "A.M.next"},
		{0xf, "X.M.", "X.M"}, {0xf, "A.M-foo", "A.M"},
		{0x1d, "abc1def2ghi", "abc1def2ghi"}, {0x1d, "abc12def", "abc12def"},
		{0x1d, "abc1", "abc1"}, {0x1d, "abc.def", "abc"}, {0x1d, "int'l", "int"},
		{0x1b, "Mrs.", "Mrs."}, {0x1b, "Mrs.9", "Mrs."}, {0x1b, "mrs.", "mrs"},
		{0x1b, "Sec./", "Sec."}, {0x1b, "SEC.", "SEC"}, {0x1b, "ca.:", "ca."},
		{0x1b, "CA.", "CA"}, {0x1b, "Mrs.next", "Mrs.next"},
		{0x1b, "ab/ next", "ab/"}, {0x1b, "abc/ next", "abc"},
		{0x1b, "ab/9", "ab"}, {0x1b, "ab/", "ab/"}, {0x1b, "ab//", "ab"},
		{0x1b, "ab/;", "ab/"}, {0x1b, "ab/:next", "ab/"},
		{0x1b, "a/b", "a/b"}, {0x1b, "ab/c", "ab"},
		{0x1b, "S.&P.", "S.&P."}, {0x1b, "S&P", "S&P"},
		{0x1b, "s&P", "s"}, {0x1b, "S&p", "S"}, {0x1b, "s.&P", "s"},
		{0x1b, "OK!next", "OK!"}, {0x1b, "ok!", "ok"}, {0x1b, "OK!!", "OK!"},
	} {
		got, handled, err := ScanPaul2013ModelScannerASCII([]byte(" \n"+test.source+"\x00"), 100, test.mode)
		advance := int32(2 + len(test.want))
		want := Paul2013ModelScannerResult{Field0: 2, PrefixFlag: 1, First: 102, Second: uint32(100 + advance), Field14: advance, Advance: advance, Text: []byte(test.want), TextLength: int32(len(test.want)), Status: 1}
		if err != nil || !handled || !reflect.DeepEqual(got, want) {
			t.Fatalf("mode%d %q: %+v handled%v err%v want %+v", test.mode, test.source, got, handled, err, want)
		}
	}
}

func TestScannerAdditionalWordModesLookupAndLimits(t *testing.T) {
	for _, mode := range []int32{5, 9, 0xc, 0xf, 0x11, 0x12, 0x15, 0x16, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d} {
		var calls []uint32
		scan := NewPaul2013ModelScannerWithASCII(nil, func(_ context.Context, address, pointer uint32, source []byte) (int32, error) {
			if pointer != 123 || string(source) != strings.Repeat("a", 29) {
				t.Fatal(pointer, string(source))
			}
			calls = append(calls, address)
			return 0, nil
		})
		got, err := scan(context.Background(), []byte(strings.Repeat("a", 35)+"\x00"), 0, mode, 123)
		if err != nil || got.Advance != 29 || got.Field18 != 0x10000 || !reflect.DeepEqual(calls, []uint32{0x1005e010}) {
			t.Fatal(mode, got, calls, err)
		}
		got, handled, err := ScanPaul2013ModelScannerASCII([]byte(strings.Repeat("W", 10)+"\x00"), 0, mode)
		if err != nil || !handled || got.Advance != 8 {
			t.Fatal(mode, got, handled, err)
		}

	}
}

func TestScannerDefaultWordModes(t *testing.T) {
	for _, mode := range []int32{-1, 2, 3, 4, 6, 7, 0xa, 0xb, 0xd, 0xe, 0x10, 0x13, 0x14, 0x1e, 100, 0x7fffffff} {
		for _, test := range []struct{ source, want string }{
			{"word.Next", "word.Next"}, {"int'l", "int'l"},
			{"word-next", "word"}, {"can't", "can"}, {"word", "word"},
		} {
			got, handled, err := ScanPaul2013ModelScannerASCII([]byte(test.source+"\x00"), 0, mode)
			if err != nil || !handled || string(got.Text) != test.want || got.Advance != int32(len(test.want)) || got.Status != 1 {
				t.Fatal(mode, test, got, handled, err)
			}
		}

	}
}
