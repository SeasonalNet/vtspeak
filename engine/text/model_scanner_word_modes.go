package text

import "bytes"

// scanPaul2013ASCIIWordJoin supplies common and mode-specific letter joins
// at FUN_1005a350:1005ac34 through 1005b503. The returned byte count
// is consumed before resuming the common letter loop; separators add no cost.
func scanPaul2013ASCIIWordJoin(source []byte, start, position, end int, mode int32) int {
	if position >= end {
		return 0
	}
	letter := func(value byte) bool { return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z' }
	upper := func(value byte) bool { return value >= 'A' && value <= 'Z' }
	value := source[position]
	followingLetter := position+1 < end && letter(source[position+1])
	switch mode {
	case 8, 9, 0xc, 0x12, 0x18, 0x1a, 0x1c, 0x1d:
	default:
		if followingLetter && (value == '.' || value == '\'' && position-start == 3 && bytes.EqualFold(source[start:position], []byte("int"))) {
			return 1
		}
	}
	switch mode {
	case 5:
		if value == '\'' && bytes.EqualFold(source[start:position], []byte("pa")) && position+5 <= end && bytes.EqualFold(source[position+1:position+5], []byte("anga")) {
			return 1
		}
	case 0x1d:
		if value >= '0' && value <= '9' {
			return 1
		}
	case 0x1b:
		if position-start == 1 {
			if upper(source[start]) && value == '.' && position+2 < end && source[position+1] == '&' && upper(source[position+2]) {
				return 2
			}
			if upper(source[start]) && value == '&' && followingLetter && upper(source[position+1]) || value == '/' && followingLetter {
				return 1
			}
		}
	case 0xc, 0x1a:
		if value == '\'' && followingLetter {
			return 1
		}
	case 0x16:
		if (value == '\'' || value == '-') && followingLetter {
			return 1
		}
	case 0x11, 0x15:
		if value == '-' && followingLetter {
			return 1
		}
		if position-start == 1 && upper(source[start]) && value >= '0' && value <= '9' && followingLetter && upper(source[position+1]) {
			return 1
		}
		// FUN_10062370 checks that all six following bytes precede NUL.
		// The mapped prefix match consumes only the space; States then runs
		// through the ordinary letter loop and its accumulated cost bound.
		if value == ' ' && bytes.EqualFold(source[start:position], []byte("United")) && position+7 <= end && bytes.EqualFold(source[position+1:position+7], []byte("States")) && (position+7 == end || !letter(source[position+7])) {
			return 1
		}
	case 0x19:
		if position-start == 1 && upper(source[start]) && value >= '0' && value <= '9' {
			return 1
		}
	}
	return 0
}

// scanPaul2013ASCIIWordEnd consumes special trailing punctuation whose native
// branches finish the token instead of reentering its letter loop.
func scanPaul2013ASCIIWordEnd(source []byte, start, position, end int, mode int32) int {
	if position >= end {
		return 0
	}
	word := source[start:position]
	if mode == 0xf && source[position] == '.' && (bytes.EqualFold(word, []byte("A.M")) || bytes.EqualFold(word, []byte("P.M"))) {
		return 1
	}
	if mode != 0x1b {
		return 0
	}
	following := source[position+1] // the NUL terminator is in the supplied slice
	space := following == ' ' || following == '\t' || following == '\n' || following == '\r'
	terminal := bytes.IndexByte([]byte(".?!;"), following) >= 0
	if len(word) < 3 && source[position] == '/' && (following == 0 || space || terminal || bytes.IndexByte([]byte(",:{}[]()<>\""), following) >= 0) {
		return 1
	}
	if source[position] == '.' && (string(word) == "Mrs" || string(word) == "S.&P" || string(word) == "Sec" || string(word) == "ca") && (following == 0 || following >= '0' && following <= '9' || space || terminal || bytes.IndexByte([]byte(",:{}[]()<>\";/"), following) >= 0) {
		return 1
	}
	if source[position] == '!' && string(word) == "OK" {
		return 1
	}
	return 0
}
