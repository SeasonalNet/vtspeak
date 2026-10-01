package text

import "testing"

func TestPaul2013ProperNameScannerWordPredicateSingleCharacter(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{name: "uppercase letter", input: []byte("A"), want: true},
		{name: "lowercase first-byte class", input: []byte("x"), want: false},
		{name: "digit", input: []byte("1"), want: false},
		{name: "punctuation", input: []byte("-"), want: false},
		{name: "empty C string", input: []byte{0}, want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Paul2013ProperNameScannerWordPredicate(test.input); got != test.want {
				t.Fatalf("Paul2013ProperNameScannerWordPredicate(%q) = %t, want %t", test.input, got, test.want)
			}
		})
	}
}
