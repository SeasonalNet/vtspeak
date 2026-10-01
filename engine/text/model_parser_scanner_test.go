package text

import (
	"bytes"
	"strings"
	"testing"
)

func TestObservePaul2013ProperNameASCIITokenSupportedPaths(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantOutput string
		wantCount  int32
		wantStatus int32
	}{
		{name: "digit run", input: []byte("12345\x00"), wantOutput: "12345", wantCount: 5, wantStatus: 2},
		{name: "digit run before terminal space", input: []byte("12345 \x00"), wantOutput: "12345", wantCount: 5, wantStatus: 2},
		{name: "digit run before terminal punctuation", input: []byte("12345!\x00"), wantOutput: "12345", wantCount: 5, wantStatus: 2},
		{name: "decimal number", input: []byte("1.25\x00"), wantOutput: "1.25", wantCount: 4, wantStatus: 2},
		{name: "decimal before terminal punctuation", input: []byte("1.25!\x00"), wantOutput: "1.25", wantCount: 4, wantStatus: 2},
		{name: "thousands group", input: []byte("1,234\x00"), wantOutput: "1,234", wantCount: 5, wantStatus: 2},
		{name: "repeated groups and decimal", input: []byte("12,345,678.90\x00"), wantOutput: "12,345,678.90", wantCount: 13, wantStatus: 2},
		{name: "first ordinal", input: []byte("1st\x00"), wantOutput: "1st", wantCount: 3, wantStatus: 2},
		{name: "second ordinal", input: []byte("22nd\x00"), wantOutput: "22nd", wantCount: 4, wantStatus: 2},
		{name: "third ordinal", input: []byte("103rd\x00"), wantOutput: "103rd", wantCount: 5, wantStatus: 2},
		{name: "teen ordinal uses th", input: []byte("113th\x00"), wantOutput: "113th", wantCount: 5, wantStatus: 2},
		{name: "ordinal before terminal punctuation", input: []byte("21st!\x00"), wantOutput: "21st", wantCount: 4, wantStatus: 2},
		{name: "mapped uppercase ordinal suffix", input: []byte("21ST\x00"), wantOutput: "21ST", wantCount: 4, wantStatus: 2},
		{name: "mapped mixed-case ordinal suffix", input: []byte("22Nd\x00"), wantOutput: "22Nd", wantCount: 4, wantStatus: 2},
		{name: "maximum digit run", input: []byte(strings.Repeat("1", 28) + "\x00"), wantOutput: strings.Repeat("1", 28), wantCount: 28, wantStatus: 2},
		{name: "uppercase digit prefix and word", input: []byte("A1BETA\x00"), wantOutput: "A1BETA", wantCount: 6, wantStatus: 1},
		{name: "digit uppercase digit run", input: []byte("1A234\x00"), wantOutput: "1A234", wantCount: 5, wantStatus: 2},
		{name: "internal hyphen", input: []byte("AGUA-DULCE\x00"), wantOutput: "AGUA-DULCE", wantCount: 10, wantStatus: 1},
		{name: "internal dot", input: []byte("ST.LOUIS\x00"), wantOutput: "ST.LOUIS", wantCount: 8, wantStatus: 1},
		{name: "mapped phrase", input: []byte("United States\x00"), wantOutput: "United States", wantCount: 13, wantStatus: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ObservePaul2013ProperNameASCIIToken(test.input)
			if !got.Supported || got.Reason != "" {
				t.Fatalf("observation = %+v, want supported path", got)
			}
			if !bytes.Equal(got.Output, []byte(test.wantOutput)) || got.ResultCount != test.wantCount || got.Status != test.wantStatus || got.PrefixCount != 0 {
				t.Fatalf("observation = %+v, want output=%q count=%d status=%d prefix=0", got, test.wantOutput, test.wantCount, test.wantStatus)
			}
		})
	}
}

func TestObservePaul2013ProperNameASCIITokenStopsAtNativeBoundaries(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		wantOutput   string
		wantCount    int32
		wantConsumed int32
		wantStatus   int32
		wantPrefix   int32
	}{
		{name: "letters before following token", input: []byte("AGUA-DULCE next\x00"), wantOutput: "AGUA-DULCE", wantCount: 10, wantConsumed: 10, wantStatus: 1},
		{name: "digit run before terminal whitespace", input: []byte("12345 \x00"), wantOutput: "12345", wantCount: 5, wantConsumed: 5, wantStatus: 2},
		{name: "digit run before terminal punctuation", input: []byte("12345!\x00"), wantOutput: "12345", wantCount: 5, wantConsumed: 5, wantStatus: 2},
		{name: "decimal before terminal punctuation", input: []byte("1.25!\x00"), wantOutput: "1.25", wantCount: 4, wantConsumed: 4, wantStatus: 2},
		{name: "ordinal before terminal punctuation", input: []byte("21st!\x00"), wantOutput: "21st", wantCount: 4, wantConsumed: 4, wantStatus: 2},
		{name: "comma group before terminal space", input: []byte("1,234 \x00"), wantOutput: "1,234", wantCount: 5, wantConsumed: 5, wantStatus: 2},
		{name: "punctuation ends letter token", input: []byte("AGUA!DULCE\x00"), wantOutput: "AGUA", wantCount: 4, wantConsumed: 4, wantStatus: 1},
		{name: "United States before following token", input: []byte("United States next\x00"), wantOutput: "United States", wantCount: 13, wantConsumed: 13, wantStatus: 1},
		{name: "United States before punctuation", input: []byte("United States!next\x00"), wantOutput: "United States", wantCount: 13, wantConsumed: 13, wantStatus: 1},
		{name: "United States before numeric token", input: []byte("United States123\x00"), wantOutput: "United States", wantCount: 13, wantConsumed: 13, wantStatus: 1},
		{name: "trailing hyphen is a delimiter", input: []byte("AGUA-\x00"), wantOutput: "AGUA", wantCount: 4, wantConsumed: 4, wantStatus: 1},
		{name: "trailing dot is a delimiter", input: []byte("ST.\x00"), wantOutput: "ST", wantCount: 2, wantConsumed: 2, wantStatus: 1},
		{name: "leading prefix and following token", input: []byte(" \tAGUA next\x00"), wantOutput: "AGUA", wantCount: 4, wantConsumed: 6, wantStatus: 1, wantPrefix: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ObservePaul2013ProperNameASCIIToken(test.input)
			if !got.Supported || string(got.Output) != test.wantOutput || got.ResultCount != test.wantCount ||
				got.ConsumedBytes != test.wantConsumed || got.Status != test.wantStatus || got.PrefixCount != test.wantPrefix {
				t.Fatalf("observation = %+v; want output=%q count=%d consumed=%d status=%d prefix=%d", got, test.wantOutput, test.wantCount, test.wantConsumed, test.wantStatus, test.wantPrefix)
			}
		})
	}
}

func TestObservePaul2013ProperNameASCIITokenPrefixesAndTerminalPaths(t *testing.T) {
	leading := ObservePaul2013ProperNameASCIIToken([]byte(" \tAGUA-DULCE\x00"))
	if !leading.Supported || leading.PrefixCount != 2 || leading.SingleLineFeedMarker || string(leading.Output) != "AGUA-DULCE" {
		t.Fatalf("leading whitespace observation = %+v", leading)
	}
	oneLineFeed := ObservePaul2013ProperNameASCIIToken([]byte("\nAGUA-DULCE\x00"))
	if !oneLineFeed.Supported || oneLineFeed.PrefixCount != 1 || !oneLineFeed.SingleLineFeedMarker || string(oneLineFeed.Output) != "AGUA-DULCE" {
		t.Fatalf("single-LF observation = %+v", oneLineFeed)
	}
	multiline := ObservePaul2013ProperNameASCIIToken([]byte("\n\n\x00"))
	if !multiline.Supported || multiline.Status != 8 || multiline.ResultCount != 0 {
		t.Fatalf("multiline observation = %+v", multiline)
	}
	empty := ObservePaul2013ProperNameASCIIToken([]byte(" \t\x00"))
	if !empty.Supported || empty.Status != 9 || empty.PrefixCount != 2 || empty.ResultCount != 0 {
		t.Fatalf("empty observation = %+v", empty)
	}
}

func TestObservePaul2013ProperNameASCIITokenRejectsUnsupportedShapes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "missing NUL", input: []byte("AGUA-DULCE")},
		{name: "unsupported non-ASCII token", input: []byte("ÉGUA\x00")},
		{name: "unsupported mixed form", input: []byte("A12BETA\x00")},
		{name: "numeric token with longer lookahead", input: []byte("123 next\x00")},
		{name: "decimal token with longer lookahead", input: []byte("1.25 next\x00")},
		{name: "comma group with short digits", input: []byte("1,23\x00")},
		{name: "comma group with excess digits", input: []byte("1,2345\x00")},
		{name: "comma group with leading zero", input: []byte("0,123\x00")},
		{name: "multiple decimal points", input: []byte("1.2.3\x00")},
		{name: "ordinal suffix does not match last digit", input: []byte("21th\x00")},
		{name: "teen ordinal requires th", input: []byte("11st\x00")},
		{name: "native exact Rd exception", input: []byte("33Rd\x00")},
		{name: "ordinal suffix followed by longer lookahead", input: []byte("21st next\x00")},
		{name: "oversized digit run", input: []byte(strings.Repeat("1", 29) + "\x00")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ObservePaul2013ProperNameASCIIToken(test.input)
			if got.Supported || got.Reason == "" {
				t.Fatalf("observation = %+v, want unsupported with a reason", got)
			}
		})
	}
}
