package text

import (
	"testing"
)

func terminalContextFixture(previous, punctuation, following string, prefix int32) Paul2013TerminalContextWindow {
	var window Paul2013TerminalContextWindow
	window[2] = Paul2013ModelScannerResult{Text: []byte(previous), TextLength: int32(len(previous))}
	window[3] = Paul2013ModelScannerResult{Text: []byte(punctuation), TextLength: int32(len(punctuation))}
	window[4] = Paul2013ModelScannerResult{Text: []byte(following), TextLength: int32(len(following)), Field0: prefix, Status: 1}
	return window
}

func TestTerminalContextClassifierOrderedBranches(t *testing.T) {
	for _, test := range []struct {
		name   string
		window Paul2013TerminalContextWindow
		want   int16
	}{
		{"question", terminalContextFixture("word", "?", "next", 1), 'P'},
		{"exclamation", terminalContextFixture("word", "!", "next", 1), 'P'},
		{"semicolon", terminalContextFixture("word", ";", "next", 1), 'P'},
		{"empty following", terminalContextFixture("word", ".", "", 0), 'P'},
		{"two byte prefix", terminalContextFixture("word", ".", "Next", 2), 'P'},
		{"upper upper adjacency", terminalContextFixture("Word", ".", "Next", 0), 'N'},
		{"upper lower adjacency", terminalContextFixture("Word", ".", "next", 0), 'N'},
		{"lower upper adjacency", terminalContextFixture("word", ".", "Next", 0), 'P'},
		{"lower lower adjacency", terminalContextFixture("word", ".", "next", 0), 'N'},
		{"word numeric adjacency", terminalContextFixture("word", ".", "12", 0), 'N'},
		{"single letters spaced", terminalContextFixture("A", ".", "B", 1), 'N'},
		{"and exact", terminalContextFixture("word", ".", "and", 1), 'N'},
		{"and case sensitive", terminalContextFixture("word", ".", "And", 1), 'A'},
		{"generic", terminalContextFixture("word", ".", "next", 1), 'A'},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := ClassifyPaul2013TerminalContext(test.window, []byte(" following\x00"))
			if err != nil || got != test.want {
				t.Fatalf("classification %c want %c, %v", got, test.want, err)
			}
		})
	}
	window := terminalContextFixture("word", ".", "next", 2)
	window[4].Status = 8
	if got, err := ClassifyPaul2013TerminalContext(window, nil); err != nil || got != 'N' {
		t.Fatal("multiline status gate lost native precedence")
	}
}

func TestTerminalContextClassifierMultiTokenAndQuotedPaths(t *testing.T) {
	window := terminalContextFixture("word", ".", "Next", 1)
	window[1] = Paul2013ModelScannerResult{Text: []byte("before"), TextLength: 6}
	if got, err := ClassifyPaul2013TerminalContext(window, nil); err != nil || got != 'P' {
		t.Fatal("unseparated prior tokens did not select P")
	}
	window = terminalContextFixture("word", ".", "'", 1)
	window[4].Status = 3
	window[5] = Paul2013ModelScannerResult{Text: []byte("Next"), TextLength: 4, Status: 1}
	if got, err := ClassifyPaul2013TerminalContext(window, nil); err != nil || got != 'P' {
		t.Fatal("quoted following token did not select P")
	}
	for _, left := range []string{"neither", "NOR", "So"} {
		for _, middle := range []string{"do", "AM"} {
			window = terminalContextFixture("I", ".", "next", 1)
			window[0] = Paul2013ModelScannerResult{Text: []byte(left), TextLength: int32(len(left))}
			window[1] = Paul2013ModelScannerResult{Text: []byte(middle), TextLength: int32(len(middle))}
			if got, err := ClassifyPaul2013TerminalContext(window, nil); err != nil || got != 'P' {
				t.Fatalf("native mapped literal branch %s/%s -> %c, %v", left, middle, got, err)
			}
		}
	}
}

func TestNativeTerminalCharacterPredicates(t *testing.T) {
	for _, test := range []struct {
		text []byte
		mask byte
		want bool
	}{{[]byte("Hello"), 0xc0, true}, {[]byte("Hello9"), 0xc0, false}, {[]byte("0123"), 0x10, true}, {[]byte("12.3"), 0x10, false}, {nil, 0xc0, false}, {[]byte{0x80}, 0xc0, false}, {[]byte{0xff}, 0x10, false}, {[]byte{'A', 0, 0xff}, 0xc0, true}} {
		if got, err := Paul2013CStringHasCharacterMask(test.text, test.mask); err != nil || got != test.want {
			t.Fatalf("mask %#x text %v = %v, %v", test.mask, test.text, got, err)
		}
	}
}
