package text

import (
	"bytes"
	"testing"
)

func TestShapePaul2013NormalizerClassCodesDecodesTreeEscapeBytes(t *testing.T) {
	got, ok := ShapePaul2013NormalizerClassCodes([]byte{7, '6', 8, '6', 9, '6'})
	if !ok {
		t.Fatal("nonempty class sequence was rejected")
	}
	if want := []byte{0x1a, 0x1b, 0x1c}; !bytes.Equal(got, want) {
		t.Fatalf("shaped class codes = % x, want % x", got, want)
	}
}

func TestShapePaul2013NormalizerClassCodesRewritesDirectPunctuation(t *testing.T) {
	got, ok := ShapePaul2013NormalizerClassCodes([]byte{'#', '$', '%', ';', '<', '='})
	if !ok {
		t.Fatal("nonempty class sequence was rejected")
	}
	if want := []byte{'&', '\'', '(', '>', '?', '@'}; !bytes.Equal(got, want) {
		t.Fatalf("shaped punctuation = % x, want % x", got, want)
	}
}

func TestShapePaul2013NormalizerClassCodesInsertsCommaEscape(t *testing.T) {
	got, ok := ShapePaul2013NormalizerClassCodes([]byte{'a', ','})
	if !ok {
		t.Fatal("nonempty class sequence was rejected")
	}
	if want := []byte{'a', 7, ','}; !bytes.Equal(got, want) {
		t.Fatalf("shaped comma sequence = % x, want % x", got, want)
	}
}

func TestShapePaul2013NormalizerClassCodesStopsAtCStringTerminator(t *testing.T) {
	got, ok := ShapePaul2013NormalizerClassCodes([]byte{'B', 0, 'C'})
	if !ok || !bytes.Equal(got, []byte{'>'}) {
		t.Fatalf("shaped C string = (% x, %t), want (>, true)", got, ok)
	}
	if got, ok := ShapePaul2013NormalizerClassCodes([]byte{0, 'B'}); ok || got != nil {
		t.Fatalf("empty C string = (% x, %t), want (nil, false)", got, ok)
	}
}

func TestRemapPaul2013NormalizerClassCodesAppliesExplicitAndTableMappings(t *testing.T) {
	got, ok, err := RemapPaul2013NormalizerClassCodes([]byte{13, '\''})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("mapped class sequence was rejected")
	}
	if want := []byte{8, '!', 'D'}; !bytes.Equal(got, want) {
		t.Fatalf("remapped class sequence = % x, want % x", got, want)
	}
}

func TestRemapPaul2013NormalizerClassCodesAppliesRawTableForPoundCode(t *testing.T) {
	got, ok, err := RemapPaul2013NormalizerClassCodes([]byte{'#'})
	if err != nil || !ok || !bytes.Equal(got, []byte{30}) {
		t.Fatalf("pound category mapping = (% x, %t, %v), want (1e, true, nil)", got, ok, err)
	}
	if got, ok, err := RemapPaul2013NormalizerClassCodes([]byte{0x5a}); err != nil || ok || got != nil {
		t.Fatalf("out-of-range category = (% x, %t, %v), want (nil, false, nil)", got, ok, err)
	}
}

func TestRemapPaul2013NormalizerClassCodesCoversEveryInRangeTableEntry(t *testing.T) {
	for value := byte(2); value <= 'Y'; value++ {
		if _, _, err := RemapPaul2013NormalizerClassCodes([]byte{value}); err != nil {
			t.Errorf("category 0x%02x hit an unhandled table sentinel: %v", value, err)
		}
	}
}
