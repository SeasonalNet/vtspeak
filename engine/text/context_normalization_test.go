package text

import (
	"bytes"
	"testing"
)

func TestApplyPaul2013LiteralContextHandler(t *testing.T) {
	tests := []struct {
		name       string
		surface    string
		preceding  []byte
		wantAppend []byte
	}{
		{name: "possessive fallback", surface: "'s", wantAppend: []byte{0x44}},
		{name: "possessive punctuation", surface: "'s", preceding: []byte{' '}, wantAppend: []byte{0x37}},
		{name: "possessive voiced context", surface: "'s", preceding: []byte{0x14}, wantAppend: []byte{0x23, 0x44}},
		{name: "ve context", surface: "'ve", preceding: []byte{0x30}, wantAppend: []byte{0x41}},
		{name: "ve fallback", surface: "'VE", wantAppend: []byte{0x07, 0x41}},
		{name: "re context", surface: "'re", preceding: []byte{0x30}, wantAppend: []byte{0x36}},
		{name: "re fallback", surface: "'RE", wantAppend: []byte{0x1a}},
		{name: "ll after six", surface: "'ll", preceding: []byte{0x30, '6'}, wantAppend: []byte{0x2b}},
		{name: "ll fallback", surface: "'LL", wantAppend: []byte{0x07, 0x2b}},
		{name: "em uppercase", surface: "'EM", preceding: []byte{0x30, '6'}, wantAppend: []byte{0x15}},
		{name: "m fallback", surface: "'m", preceding: []byte{0x30, 0x20}, wantAppend: []byte{0x07, 0x15}},
		{name: "n context", surface: "'n", preceding: []byte{0x30}, wantAppend: []byte{0x2c}},
		{name: "d", surface: "'d", wantAppend: []byte{0x26}},
		{name: "er", surface: "'er", wantAppend: []byte{0x1a}},
		{name: "ee", surface: "'ee", wantAppend: []byte{0x1a, 0x44}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			row := make([]byte, Paul2013PhoneContextRowSize)
			row[Paul2013PhoneContextRowFlags] = 0x04
			copy(row[Paul2013PhoneContextRowPhoneCodes:], test.preceding)
			row[Paul2013PhoneContextRowPhoneCodes+len(test.preceding)] = 0
			code, handled, err := ApplyPaul2013LiteralContextHandler(row, []byte(test.surface))
			if err != nil {
				t.Fatal(err)
			}
			if !handled || code != 1 {
				t.Fatalf("handler result = (%d,%t), want (1,true)", code, handled)
			}
			start := Paul2013PhoneContextRowPhoneCodes + len(test.preceding)
			got := row[start : start+len(test.wantAppend)]
			if !bytes.Equal(got, test.wantAppend) || row[start+len(test.wantAppend)] != 0 {
				t.Fatalf("appended codes = % x plus %02x, want % x plus NUL", got, row[start+len(test.wantAppend)], test.wantAppend)
			}
			if row[Paul2013PhoneContextRowFlags] != 0x0c {
				t.Fatalf("row flags = 0x%02x, want 0x0c", row[Paul2013PhoneContextRowFlags])
			}
		})
	}
}

func TestApplyPaul2013LiteralContextHandlerLeavesUnknownAndRejectsBadRows(t *testing.T) {
	row := make([]byte, Paul2013PhoneContextRowSize)
	copy(row[Paul2013PhoneContextRowPhoneCodes:], []byte{0x30, 0})
	before := append([]byte(nil), row...)
	if _, handled, err := ApplyPaul2013LiteralContextHandler(row, []byte("unmapped")); err != nil || handled {
		t.Fatalf("unknown handler = (%t,%v), want (false,nil)", handled, err)
	}
	if !bytes.Equal(row, before) {
		t.Fatal("unknown handler changed the phone/context row")
	}
	if _, _, err := ApplyPaul2013LiteralContextHandler(row[:Paul2013PhoneContextRowSize-1], []byte("'s")); err == nil {
		t.Fatal("short phone/context row was accepted")
	}
	unterminated := make([]byte, Paul2013PhoneContextRowSize)
	for index := Paul2013PhoneContextRowPhoneCodes; index < Paul2013PhoneContextRowSourceMarker; index++ {
		unterminated[index] = 1
	}
	if _, _, err := ApplyPaul2013LiteralContextHandler(unterminated, []byte("'s")); err == nil {
		t.Fatal("unterminated phone-code row was accepted")
	}
}
