package text

import (
	"strings"
	"testing"
)

func TestExpandPaul2013CapturedNumbers(t *testing.T) {
	tests := map[string]string{
		"0":                "zero",
		"7":                "seven",
		"42":               "forty two",
		"105":              "one hundred five",
		"1000":             "one thousand",
		"1001":             "one thousand one",
		"1009":             "one thousand nine",
		"1234":             "twelve thirty four",
		"2024":             "twenty twenty four",
		"1010":             "ten ten",
		"1099":             "ten ninety nine",
		"1100":             "eleven hundred",
		"1900":             "nineteen hundred",
		"1999":             "nineteen ninety nine",
		"2000":             "two thousand",
		"2001":             "two thousand one",
		"2005":             "two thousand five",
		"2009":             "two thousand nine",
		"2010":             "twenty ten",
		"2100":             "twenty one hundred",
		"2101":             "twenty one oh one",
		"9999":             "ninety nine ninety nine",
		"007":              "oh oh seven",
		"3.14":             "three point one four",
		"12.05":            "twelve point zero five",
		".5":               "point five",
		"1,000.00":         "one thousand point zero zero",
		"-12.5":            "minus twelve point five",
		"+7":               "plus seven",
		"1000000000000000": "one oh oh oh oh oh oh oh oh oh oh oh oh oh oh oh",
		"01/02/2024":       "January second twenty twenty four",
		"3:45":             "three forty five",
		"$5.00":            "five dollars",
		"25%":              "twenty five percent",
		"555-1234":         "five hundred fifty five to twelve thirty four",
	}
	for surface, want := range tests {
		t.Run(surface, func(t *testing.T) {
			words, err := expandNumericSurface(surface)
			if err != nil {
				t.Fatalf("expandNumericSurface(%q): %v", surface, err)
			}
			if got := strings.Join(words, " "); got != want {
				t.Fatalf("expandNumericSurface(%q) = %q, want %q", surface, got, want)
			}
		})
	}
}

func TestExpandPaul2013NumericFormsRejectMalformedInputs(t *testing.T) {
	tests := []string{
		"+", "1.2.3", "1,00", "13st", "13/01/2024", "13:00", "555-123",
	}
	for _, surface := range tests {
		t.Run(surface, func(t *testing.T) {
			if _, err := expandNumericSurface(surface); err == nil {
				t.Fatalf("expandNumericSurface(%q) accepted malformed input", surface)
			}
		})
	}
}
