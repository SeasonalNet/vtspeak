package text

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestScannerMode8WordJoins(t *testing.T) {
	for _, mode := range []int32{8} {
		for _, test := range []struct{ source, want string }{
			{"word.next", "word.next"}, {"word.Next.more", "word"}, {"Word.Next.more", "Word.Next.more"}, {"word.next.More", "word.next.More"},
			{"one-two-three", "one-two-three"}, {"can't-next", "can't-next"},
			{"word.'ve-next", "word.'ve"}, {"word's-next", "word's-next"},
			{"int'l", "int"}, {"INT'ernational", "INT"},
			{"int'foo", "int"}, {"wi'abc", "wi'"}, {"WI'", "WI'"},
			{"word'foo", "word"}, {"words' next", "words'"}, {"words'", "words"},
			{"word-9", "word"}, {"word..next", "word"},
			{strings.Repeat("W", 8) + "-next", strings.Repeat("W", 8)},
			{strings.Repeat("W", 8) + ".next", strings.Repeat("W", 8)},
			{"one-two-" + strings.Repeat("W", 10), "one-two"},
			{"one.two." + strings.Repeat("W", 10), "one.two"},
		} {
			got, handled, err := ScanPaul2013ModelScannerASCII([]byte(" \n"+test.source+"\x00"), 50, mode)
			advance := int32(2 + len(test.want))
			want := Paul2013ModelScannerResult{Field0: 2, PrefixFlag: 1, First: 52, Second: uint32(50 + advance), Text: []byte(test.want), TextLength: int32(len(test.want)), Field14: advance, Advance: advance, Status: 1}
			if err != nil || !handled || !reflect.DeepEqual(got, want) {
				t.Fatalf("mode%d %q: %+v handled%v err%v want %+v", mode, test.source, got, handled, err, want)
			}
		}
	}
}

func TestScannerJoinLookupOrdering(t *testing.T) {
	for _, separator := range []string{"-", "."} {
		for _, match := range []uint32{0, 0x1005e010, 0x1005e1a0} {
			var calls []uint32
			var sources []string
			scan := NewPaul2013ModelScannerWithASCII(nil, func(_ context.Context, address, pointer uint32, source []byte) (int32, error) {
				if pointer != 123 {
					t.Fatal(pointer)
				}
				calls = append(calls, address)
				sources = append(sources, string(source))
				if string(source) == "next\x00" && address == match {
					return 0, nil
				}
				return -1, nil
			})
			got, err := scan(context.Background(), []byte("word"+separator+"next\x00"), 10, 8, 123)
			wantCalls := []uint32{0x1005e010, 0x1005e1a0}
			wantSources := []string{"next\x00", "next\x00"}
			wantText := "word"
			if match == 0x1005e010 {
				wantCalls = wantCalls[:1]
				wantSources = wantSources[:1]
			} else if match == 0 {
				wantText = "word" + separator + "next"
				wantCalls = append(wantCalls, 0x1005e010, 0x1005e1a0)
				wantSources = append(wantSources, wantText, wantText)
			}
			if err != nil || string(got.Text) != wantText || got.Advance != int32(len(wantText)) || got.Field18 != 0 || !reflect.DeepEqual(calls, wantCalls) || !reflect.DeepEqual(sources, wantSources) {
				t.Fatal(match, got, calls, sources, err)
			}
		}
	}
}

func TestScannerJoinsWithoutInternalLookup(t *testing.T) {
	for _, test := range []struct {
		source string
		mode   int32
	}{{"word.Next", 1}, {"can't-next", 8}, {"word-next", 0x17}} {
		var calls []uint32
		scan := NewPaul2013ModelScannerWithASCII(nil, func(_ context.Context, address, _ uint32, source []byte) (int32, error) {
			if string(source) != test.source {
				t.Fatalf("unexpected internal lookup %q", source)
			}
			calls = append(calls, address)
			return 4, nil
		})
		got, err := scan(context.Background(), []byte(test.source+"\x00"), 0, test.mode, 123)
		if err != nil || string(got.Text) != test.source || got.Field18 != 0x10000 || !reflect.DeepEqual(calls, []uint32{0x1005e010}) {
			t.Fatal(test, got, calls, err)
		}
	}
	marker := errors.New("lookup error")
	scan := NewPaul2013ModelScannerWithASCII(nil, func(context.Context, uint32, uint32, []byte) (int32, error) { return 0, marker })
	if _, err := scan(context.Background(), []byte("word-next\x00"), 0, 8, 123); !errors.Is(err, marker) {
		t.Fatal(err)
	}
}

func TestScannerMode01JoinsRemainDistinct(t *testing.T) {
	for _, mode := range []int32{0, 1} {
		for _, test := range []struct{ source, want string }{
			{"word.Next.more", "word.Next.more"}, {"word-next", "word"},
			{"can't-next", "can"}, {"word.'ve", "word"}, {"wi'abc", "wi"},
			{"int'l", "int'l"}, {"INT'ernational", "INT'ernational"},
			{"int'foo", "int'foo"}, {"word's", "word"},
			{strings.Repeat("W", 8) + ".next", strings.Repeat("W", 8)},
		} {
			got, handled, err := ScanPaul2013ModelScannerASCII([]byte(test.source+"\x00"), 0, mode)
			if err != nil || !handled || string(got.Text) != test.want || got.Advance != int32(len(test.want)) {
				t.Fatal(mode, test, got, handled, err)
			}
		}
	}
}

func TestScannerMode8FallbackPreservesArguments(t *testing.T) {
	for _, source := range []string{"?", "word-\xff", "word.\xff", "int'\xff"} {
		scan := NewPaul2013ModelScannerWithASCII(func(_ context.Context, got []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
			if string(got) != source+"\x00" || offset != 50 || mode != 8 || pointer != 123 {
				t.Fatal(got, offset, mode, pointer)
			}
			return Paul2013ModelScannerResult{Advance: 99}, nil
		}, nil)
		got, err := scan(context.Background(), []byte(source+"\x00"), 50, 8, 123)
		if err != nil || got.Advance != 99 {
			t.Fatal(source, got, err)
		}
	}
}
