package text

import (
	"bytes"
	"testing"
)

func TestEncodePaul2013ContextString(t *testing.T) {
	got, err := EncodePaul2013ContextString([]byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x1e, 0x14, 'd', 0x27, 'd', 0x18, 0x2b, 'd', 0x18, 0x2b, 'd', 0x30}
	if !bytes.Equal(got, want) {
		t.Fatalf("encoded Hello = % x, want % x", got, want)
	}
}

func TestBuildPaul2013ContextCodesFromSourceString(t *testing.T) {
	got, err := BuildPaul2013ContextCodesFromSourceString([]byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}
	wantCodes := []byte{0x1e, 0x14, 0x27, 0x18, 0x2b, 0x18, 0x2b, 0x30}
	wantMarkers := []byte{'0', '1', '1', '0', '1', '0', '1', '0'}
	if !bytes.Equal(got.Codes, wantCodes) || !bytes.Equal(got.PhoneMarkers, wantMarkers) {
		t.Fatalf("normalized context codes/markers = % x/%q, want % x/%q", got.Codes, got.PhoneMarkers, wantCodes, wantMarkers)
	}
}

func TestEncodePaul2013ContextStringDelimiters(t *testing.T) {
	tests := []struct {
		input string
		want  []byte
	}{
		{input: "A.B", want: []byte{0x1e, 'd', 'd', 0x13, 0x27}},
		{input: "A.", want: []byte{0x1e}},
		{input: "A B", want: []byte{0x1e}},
		{input: "a-b", want: []byte{0x1e, 'd', 'd', 0x13, 0x27}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := EncodePaul2013ContextString([]byte(test.input))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, test.want) {
				t.Fatalf("encoded %q = % x, want % x", test.input, got, test.want)
			}
		})
	}
}

func TestApplyPaul2013ContextStringToPhoneRow(t *testing.T) {
	row := make([]byte, Paul2013PhoneContextRowSize)
	row[Paul2013PhoneContextRowFlags] = 0x04
	for index := Paul2013PhoneContextRowPhoneCodes; index < Paul2013PhoneContextRowSourceMarker; index++ {
		row[index] = 0xaa
	}
	if err := ApplyPaul2013ContextStringToPhoneRow(row, []byte("A")); err != nil {
		t.Fatal(err)
	}
	if row[Paul2013PhoneContextRowFlags] != 0x24 || row[Paul2013PhoneContextRowPhoneCodes] != 0x1e || row[Paul2013PhoneContextRowPhoneCodes+1] != 0 || row[Paul2013PhoneContextRowPhoneCodes+2] != 0xaa {
		t.Fatalf("normalized phone row flags/codes = % x", row[:Paul2013PhoneContextRowSourceMarker])
	}
	if err := ApplyPaul2013ContextStringToPhoneRow(row, []byte("!")); err != nil {
		t.Fatal(err)
	}
	if row[Paul2013PhoneContextRowFlags] != 0x24 || row[Paul2013PhoneContextRowPhoneCodes] != 0 {
		t.Fatalf("empty normalization did not clear codes while preserving row flags: % x", row[:Paul2013PhoneContextRowSourceMarker])
	}
}

func TestEncodePaul2013ContextStringBounds(t *testing.T) {
	if _, err := EncodePaul2013ContextString([]byte("é")); err == nil {
		t.Fatal("non-ASCII context string was accepted")
	}
	long := bytes.Repeat([]byte{'W'}, 80)
	got, err := EncodePaul2013ContextString(long)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) > paul2013ContextCodeCapacity {
		t.Fatalf("encoded context string has %d bytes, exceeds native cap", len(got))
	}
}
