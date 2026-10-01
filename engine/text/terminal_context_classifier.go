package text

import (
	"bytes"
)

// Paul2013TerminalContextWindow is FUN_10052000's six 0x44-byte scanner
// records: three preceding tokens, the punctuation, and two following tokens.
// Scanner fields remain raw; their coordinate/length predicates are retained.
type Paul2013TerminalContextWindow [6]Paul2013ModelScannerResult

// Paul2013CStringHasCharacterMask ports the nonempty all-character predicate
// in FUN_1000ffa0 (mask C0) and FUN_10061940 (mask 10). The effective table
// includes the recovered zero prefix reached by negative signed-byte indexes.
func Paul2013CStringHasCharacterMask(text []byte, mask byte) (bool, error) {
	text = cString(text)
	if len(text) == 0 {
		return false, nil
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	for _, character := range text {
		if attributes[character]&mask == 0 {
			return false, nil
		}
	}
	return true, nil
}

// ClassifyPaul2013TerminalContext ports FUN_10052000's ordered N/P/A
// classifier over a recovered window. FollowingSource begins immediately
// after the punctuation's scanner +0x14 advance. No semantic interpretation
// of the native codes is assigned; A selects the later model-table fallback.
func ClassifyPaul2013TerminalContext(window Paul2013TerminalContextWindow, followingSource []byte) (value int16, err error) {
	previous, current, next, after := window[2], window[3], window[4], window[5]
	first := func(text []byte) byte {
		if len(text) == 0 {
			return 0
		}
		return text[0]
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	mask := func(text []byte, bits byte) bool {
		character := first(text)
		return attributes[character]&bits != 0
	}
	allMask := func(text []byte, bits byte) bool {
		accepted, _ := Paul2013CStringHasCharacterMask(text, bits)
		return accepted
	}
	punct := first(current.Text)
	if punct != 0 && bytes.ContainsAny([]byte{punct}, "?!;") && current.Field0 == 0 && current.TextLength == 1 {
		return 'P', nil
	}
	if punct == '.' && next.TextLength == 0 && len(followingSource) > 0 && (followingSource[0] == '\t' || followingSource[0] == '\n' || followingSource[0] == '\r') {
		return 'P', nil
	}
	if next.TextLength == 0 {
		return 'P', nil
	}
	if next.Field0 < 3 && next.Status == 8 {
		return 'N', nil
	}
	if next.Field0 >= 2 {
		return 'P', nil
	}
	wordDot := previous.TextLength > 0 && current.Field0 == 0 && current.TextLength == 1 && punct == '.'
	if wordDot && mask(previous.Text, 0x80) && mask(next.Text, 0x80) && next.Field0 == 0 {
		return 'N', nil
	}
	if wordDot && mask(previous.Text, 0x80) && mask(next.Text, 0x40) && next.Field0 == 0 {
		return 'N', nil
	}
	if wordDot && mask(previous.Text, 0x40) && mask(next.Text, 0x80) && next.Field0 == 0 {
		return 'P', nil
	}
	if current.TextLength == 1 && punct == '.' && mask(next.Text, 0xc0) && next.Field0 == 0 {
		return 'N', nil
	}
	if wordDot && allMask(previous.Text, 0xc0) && allMask(next.Text, 0x10) && next.Field0 == 0 {
		return 'N', nil
	}
	if wordDot && previous.TextLength == 1 && mask(previous.Text, 0x80) && next.Field0 < 2 && next.TextLength == 1 && mask(next.Text, 0x80) {
		return 'N', nil
	}
	// Native repeats the zero-following-length gate in this single-letter path.
	if wordDot && previous.TextLength == 1 && mask(previous.Text, 0x80) && next.TextLength == 0 {
		return 'P', nil
	}
	if previous.TextLength > 0 && current.Field0 == 0 && current.TextLength < 3 && punct == '.' && next.Field0 < 2 && bytes.Equal(cString(next.Text), []byte("and")) {
		return 'N', nil
	}
	if window[1].TextLength > 0 && window[1].Field0 == 0 && previous.Field0 == 0 && current.Field0 == 0 && current.TextLength == 1 && punct == '.' && next.Field0 > 0 && next.TextLength > 0 && mask(next.Text, 0x80) {
		return 'P', nil
	}
	if current.Field0 == 0 && current.TextLength == 1 && punct == '.' && next.TextLength == 1 && next.Field0 > 0 && next.Status == 3 && first(next.Text) != 0 && bytes.ContainsAny([]byte{first(next.Text)}, "`'\"") && after.TextLength > 0 && after.Field0 == 0 && after.Status == 1 && mask(after.Text, 0x80) {
		return 'P', nil
	}
	weights := Paul2013ContextCharacterWeights()
	equal := func(text []byte, want string) bool {
		return ComparePaul2013MappedCString(cString(text), []byte(want), weights) == 0
	}
	if window[0].TextLength > 0 && (equal(window[0].Text, "neither") || equal(window[0].Text, "nor") || equal(window[0].Text, "so")) && window[1].TextLength > 0 && (equal(window[1].Text, "do") || equal(window[1].Text, "am")) && previous.TextLength > 0 && equal(previous.Text, "I") && current.Field0 == 0 && current.TextLength < 3 && punct == '.' {
		return 'P', nil
	}
	return 'A', nil
}
