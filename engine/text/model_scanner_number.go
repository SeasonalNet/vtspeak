package text

import "bytes"

// scanPaul2013ASCIINumber ports the native ASCII digit dispatch. Resume
// reenters the shared token loop after appending a mode-specific letter.
// Source is bounded at its NUL byte; start includes preceding joined bytes.
func scanPaul2013ASCIINumber(source []byte, start, position, end int, mode int32, letterUsed *bool) (int, int32, bool) {
	digit := func(value byte) bool { return value >= '0' && value <= '9' }
	letter := func(value byte) bool { return value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' }
	if source[position] == '.' {
		position++
	}
	for {
		for position < end && digit(source[position]) && position-start < 29 {
			position++
		}
		if position < end && digit(source[position]) && position-start >= 29 {
			if boundary := bytes.LastIndexByte(source[start:position], ','); boundary >= 0 {
				position = start + boundary
			}
			return position, 2, false
		}
		if position-start >= 29 {
			break
		}
		if mode == 0x13 && !*letterUsed && position < end && letter(source[position]) && source[position] != 'X' && source[position] != 'x' && (position+1 == end || !digit(source[position+1])) {
			*letterUsed = true
			return position + 1, 1, true
		}
		if mode == 3 || mode == 6 || mode == 0xe || mode == 0x10 || mode == 0x1c {
			break
		}
		if position+3 < end && source[position] == ',' && bytes.IndexByte(source[start:position], '.') < 0 && source[start] != '0' && digit(source[position+1]) && digit(source[position+2]) && digit(source[position+3]) && (position+4 == end || !digit(source[position+4])) {
			position++
			continue
		}
		if position+1 < end && source[position] == '.' && (mode == 2 || bytes.IndexByte(source[start:position], '.') < 0) && digit(source[position+1]) {
			position++
			continue
		}
		break
	}
	if (mode == 0 || mode == 0x1b) && position-start < 29 && bytes.Equal(source[start:position], []byte("401")) {
		if position < end && (source[position] == 'K' || source[position] == 'k') && (position+1 == end || !letter(source[position+1])) {
			return position + 1, 1, false
		}
		if position+2 < end && source[position] == '(' && (source[position+1] == 'K' || source[position+1] == 'k') && source[position+2] == ')' {
			return position + 3, 1, false
		}
	}
	if position-start < 29 && position < end {
		if (mode == 0x11 || mode == 0x15) && position-start == 1 && source[position] >= 'A' && source[position] <= 'Z' && position+1 < end && digit(source[position+1]) {
			return position + 1, 1, true
		}
		if mode == 0x1d && letter(source[position]) {
			return position + 1, 1, true
		}
	}
	if mode == 0x1c {
		return position, 2, false
	}
	if bytes.IndexByte(source[start:position], '.') >= 0 || position+1 >= end || position+2 < end && letter(source[position+2]) {
		return position, 2, false
	}
	text := source[start:position]
	if len(text) == 0 {
		return position, 2, false
	}
	last := text[len(text)-1]
	want := ""
	if len(text) > 1 && text[len(text)-2] == '1' && last >= '1' && last <= '3' {
		want = "th"
	} else {
		switch last {
		case '1':
			want = "st"
		case '2':
			want = "nd"
		case '3':
			want = "rd"
		default:
			if digit(last) {
				want = "th"
			}
		}
	}
	if want != "" && ComparePaul2013MappedCString(source[position:position+2], []byte(want), Paul2013ContextCharacterWeights()) == 0 && !(want == "rd" && bytes.Equal(source[position:position+2], []byte("Rd"))) {
		position += 2
	}
	return position, 2, false
}
