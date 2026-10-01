package text

import (
	"context"
	"strings"
	"testing"
)

func TestNativeASCIINumericScannerBranches(t *testing.T) {
	scan := NewPaul2013ModelScannerWithASCII(nil, nil)
	for _, test := range []struct {
		source, text string
		mode, status int32
	}{{"42!", "42", 0, 2}, {"3.14 next", "3.14", 1, 2}, {".42!", ".42", 0, 2}, {"1,234,567.8", "1,234,567.8", 0, 2}, {"0,123", "0", 1, 2}, {"1,23", "1", 1, 2}, {"1,2345", "1", 1, 2}, {"3.1.4", "3.1", 1, 2}, {"1st!", "1st", 1, 2}, {"2ND!", "2ND", 0x17, 2}, {"3rd!", "3rd", 0, 2}, {"3Rd!", "3", 0, 2}, {"3RD!", "3RD", 0, 2}, {"11st", "11", 1, 2}, {"11th", "11th", 1, 2}, {"12th", "12th", 1, 2}, {"13th", "13th", 1, 2}, {"21st", "21st", 1, 2}, {"31th", "31", 1, 2}, {"11st", "11", 1, 2}, {"3rdx", "3", 1, 2}, {"3.1st", "3.1", 1, 2}, {"401K!", "401K", 0, 1}, {"401(k)", "401(k)", 0, 1}, {"401Kx", "401", 0, 2}, {"401K!", "401", 1, 2}, {strings.Repeat("1", 35), strings.Repeat("1", 29), 0, 2}, {strings.Repeat("1", 29) + "th", strings.Repeat("1", 29) + "th", 0, 2}} {
		got, err := scan(context.Background(), []byte(test.source+"\x00"), 0, test.mode, 0)
		if err != nil || string(got.Text) != test.text || got.Status != test.status || got.Advance != int32(len(test.text)) || got.First != 0 || got.Second != uint32(len(test.text)) {
			t.Fatalf("%+v -> %+v, %v", test, got, err)
		}
	}
	if _, handled, err := ScanPaul2013ModelScannerASCII([]byte(".42\x00"), 100, 0); err != nil || handled {
		t.Fatal("unknown preceding byte was guessed", handled, err)
	}
	got, err := scan(context.Background(), []byte(" .42\x00"), 100, 0, 0)
	if err != nil || string(got.Text) != ".42" || got.Field0 != 1 || got.First != 101 || got.Second != 104 {
		t.Fatal(got, err)
	}
}

func TestNumericLookupBoundaryWithRecoveredScanner(t *testing.T) {
	scan := NewPaul2013ModelScannerWithASCII(nil, nil)
	for _, test := range []struct {
		source string
		length int32
		want   int16
	}{{"42", 2, 1}, {"42.1", 2, 0}, {"42st", 2, 1}, {"21st", 2, 0}, {"21st", 4, 1}, {"42!", 2, 1}} {
		got, err := CheckPaul2013ScannerLookupBoundary(context.Background(), []byte(test.source), test.length, scan)
		if err != nil || got != test.want {
			t.Fatal(test, got, err)
		}
	}
}
