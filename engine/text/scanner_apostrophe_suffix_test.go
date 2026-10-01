package text

import (
	"context"
	"testing"
)

func TestScannerApostropheSuffixNativeRules(t *testing.T) {
	for _, test := range []struct {
		previous     byte
		source, want string
		accepted     bool
	}{{'s', "' next", "'", true}, {'S', "'\t", "'", true}, {'s', "'", "", false}, {'n', "'T!", "'T", true}, {'a', "'t!", "", false}, {'a', "'s1", "'s", true}, {'a', "'VE.", "'VE", true}, {'a', "'em ", "'em", true}, {'a', "'re!", "'re", true}, {'a', "'ll ", "'ll", true}, {'a', "'m ", "'m", true}, {'a', "'d ", "'d", true}, {'n', "'thing", "", false}, {'a', "'something", "", false}, {0, "'s", "", false}} {
		got, accepted := ScanPaul2013ScannerApostropheSuffix(test.previous, []byte(test.source+"\x00"))
		if accepted != test.accepted || string(got) != test.want {
			t.Fatal(test, string(got), accepted)
		}
	}
}

func TestMode17ASCIIJoinsAndBoundaryWithoutFixtures(t *testing.T) {
	scan := NewPaul2013ModelScannerWithASCII(nil, nil)
	for _, test := range []struct {
		source, text string
		advance      int32
	}{{"word", "word", 4}, {"word-next", "word-next", 9}, {"one-two-three", "one-two-three", 13}, {"can't!", "can't", 5}, {"word'S9", "word'S", 6}, {"words' next", "words'", 6}, {"words'", "words", 5}, {"word.'s!", "word.'s", 7}, {"word.'ve!", "word.'ve", 8}, {"word.next", "word.next", 9}, {"word'bad", "word", 4}, {"can't-next", "can't", 5}, {"WWWWWWWW-next", "WWWWWWWW", 8}} {
		got, err := scan(context.Background(), []byte(test.source+"\x00"), 0, 0x17, 0)
		if err != nil || string(got.Text) != test.text || got.Advance != test.advance || got.Status != 1 {
			t.Fatalf("%q -> %+v, %v", test.source, got, err)
		}
	}
	for _, test := range []struct {
		source string
		length int32
		want   int16
	}{{"word-next", 4, 1}, {"word-next", 5, 0}, {"can't", 3, 0}, {"word's", 4, 1}, {"word.'s", 4, 0}, {"word-more", 9, 1}} {
		got, err := CheckPaul2013ScannerLookupBoundary(context.Background(), []byte(test.source), test.length, scan)
		if err != nil || got != test.want {
			t.Fatal(test, got, err)
		}
	}
}
