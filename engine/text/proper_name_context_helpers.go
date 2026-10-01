package text

import "bytes"

// Paul2013ProperNameContextLateLookup ports the nonnegative-result cases in
// FUN_10041df0 as used by FUN_10034180. It recognizes uppercase Roman
// numerals I through XX and the observed two-byte A5/AF..B9 encodings.
func Paul2013ProperNameContextLateLookup(surface []byte) bool {
	surface = cString(surface)
	for _, numeral := range [...]string{
		"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X",
		"XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX",
	} {
		if bytes.Equal(surface, []byte(numeral)) {
			return true
		}
	}
	if len(surface) != 2 {
		return false
	}
	first, second := surface[0], surface[1]
	validTwoByteString := (first > 0xa0 && first < 0xae && second > 0xa0 && second != 0xff) ||
		(first == 0xfd && second == 0xfe) ||
		(first == 0xae && second > 0xa0 && second < 0xc3)
	return validTwoByteString && first == 0xa5 && second > 0xaf && second < 0xba
}

// Paul2013ProperNameScannerWordPredicate ports FUN_10010120's short-result
// predicate for the supplied scanner output. It preserves the leading MC/DE
// and apostrophe-prefix recursive cases, then checks the native class and
// punctuation sequence using the embedded 2013 character tables.
func Paul2013ProperNameScannerWordPredicate(surface []byte) bool {
	surface = cString(surface)
	attributes := Paul2013ExceptionCharacterAttributes()
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	return paul2013ProperNameScannerWordPredicate(surface, attributes, characterMap)
}

func paul2013ProperNameScannerWordPredicate(
	surface []byte,
	attributes [256]byte,
	characterMap [256]byte,
) bool {
	surface = cString(surface)
	if len(surface) >= 2 {
		first := surface[0]
		second := surface[1]
		mappedSecond := characterMap[second]
		if (first == 'M' && mappedSecond == 'C') || (first == 'D' && mappedSecond == 'E') {
			if paul2013ProperNameScannerWordPredicate(surface[2:], attributes, characterMap) {
				return true
			}
		}
		if len(surface) >= 3 && attributes[first]&0x80 != 0 && second == '\'' &&
			attributes[surface[2]]&0x80 != 0 {
			return paul2013ProperNameScannerWordPredicate(surface[2:], attributes, characterMap)
		}
	}
	if len(surface) == 0 || attributes[surface[0]]&0x80 == 0 {
		return false
	}
	if len(surface) == 1 {
		// FUN_10010120 returns success when the first character passes its
		// class gate and the next byte is already NUL.
		return true
	}

	index := 1
	for index < len(surface) && paul2013ScannerPredicatePunctuation(surface[index]) {
		index++
	}
	if index == len(surface) {
		return false
	}
	class := attributes[surface[index]]
	classMask := byte(0x80)
	if class&0x80 == 0 {
		if class&0x40 == 0 {
			return false
		}
		classMask = 0x40
	}
	for index++; index < len(surface); index++ {
		if paul2013ScannerPredicatePunctuation(surface[index]) {
			continue
		}
		if attributes[surface[index]]&classMask == 0 {
			return false
		}
	}
	return true
}

func paul2013ScannerPredicatePunctuation(value byte) bool {
	switch value {
	case '.', '-', '_', '&', '\'':
		return true
	default:
		return false
	}
}
