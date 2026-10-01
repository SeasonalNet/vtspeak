package text

import (
	"fmt"
	"strings"
)

// ParsePaul2013CMUPronunciation parses the recovered x-cmu phone spelling
// form. Vowels require one CMU stress digit (0 through 2); consonants omit it.
// The result carries the established Paul phone identities and symbol codes.
func ParsePaul2013CMUPronunciation(source string) ([]CMUPhone, error) {
	if source == "" {
		return nil, fmt.Errorf("CMU pronunciation is empty")
	}
	for offset := 0; offset < len(source); offset++ {
		if source[offset] >= 0x80 {
			return nil, fmt.Errorf("CMU pronunciation byte 0x%02x at offset %d is non-ASCII", source[offset], offset)
		}
		if source[offset] != ' ' && (source[offset] < 'A' || source[offset] > 'Z') && (source[offset] < '0' || source[offset] > '9') {
			return nil, fmt.Errorf("CMU pronunciation byte 0x%02x at offset %d is unsupported", source[offset], offset)
		}
	}
	labels := strings.Fields(source)
	if len(labels) == 0 {
		return nil, fmt.Errorf("CMU pronunciation has no phones")
	}
	phones := make([]CMUPhone, 0, len(labels))
	for index, label := range labels {
		phone := CMUPhone{Label: label}
		if last := label[len(label)-1]; last >= '0' && last <= '9' {
			phone.Label = label[:len(label)-1]
			phone.Stress = uint8(last - '0')
			phone.Vowel = true
		}
		if _, err := phone.TreeFeatures(); err != nil {
			return nil, fmt.Errorf("CMU pronunciation phone %d %q: %w", index, label, err)
		}
		phones = append(phones, phone)
	}
	return phones, nil
}

// BuildPaul2013CMUPhoneSequence creates one explicit token sequence from a
// caller-selected x-cmu pronunciation. It does not resolve a lexical surface,
// produce the native position-state arrays, or select model units.
func BuildPaul2013CMUPhoneSequence(surface, pronunciation string) (LexicalPhoneSequence, error) {
	if surface == "" {
		return LexicalPhoneSequence{}, fmt.Errorf("CMU pronunciation surface is empty")
	}
	for offset := 0; offset < len(surface); offset++ {
		if surface[offset] == 0 || surface[offset] >= 0x80 {
			return LexicalPhoneSequence{}, fmt.Errorf("CMU pronunciation surface byte 0x%02x at offset %d is unsupported", surface[offset], offset)
		}
	}
	phones, err := ParsePaul2013CMUPronunciation(pronunciation)
	if err != nil {
		return LexicalPhoneSequence{}, err
	}
	sequence := LexicalPhoneSequence{
		Phones: phones,
		Tokens: []LexicalTokenSpan{{
			SourceSurface: surface, Surface: surface,
			SourceByteStart: 0, SourceByteEnd: len(surface) - 1, HasSourceByteSpan: true,
			PhoneStart: 0, PhoneEnd: len(phones),
		}},
	}
	sequence.Symbols = make([]byte, len(phones))
	for index, phone := range phones {
		features, err := phone.TreeFeatures()
		if err != nil {
			return LexicalPhoneSequence{}, fmt.Errorf("CMU pronunciation phone %d: %w", index, err)
		}
		sequence.Symbols[index] = features.SymbolCode
	}
	return sequence, nil
}
