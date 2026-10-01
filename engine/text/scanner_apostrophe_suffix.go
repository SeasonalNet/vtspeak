package text

import "bytes"

// ScanPaul2013ScannerApostropheSuffix ports FUN_100621f0's returned suffix.
// Unlike FUN_1005f9b0 it permits trailing nonletters, recognizes n't using
// the preceding n, and permits a plural-s apostrophe only before whitespace.
func ScanPaul2013ScannerApostropheSuffix(previous byte, source []byte) ([]byte, bool) {
	source = cString(source)
	if previous == 0 || len(source) == 0 || source[0] != '\'' {
		return nil, false
	}
	weights := Paul2013ContextCharacterWeights()
	attributes := Paul2013ExceptionCharacterAttributes()
	value := func(index int) byte {
		if index >= len(source) {
			return 0
		}
		return source[index]
	}
	if weights[previous] == 's' && bytes.IndexByte([]byte(" \t\n\r"), value(1)) >= 0 {
		return []byte{'\''}, true
	}
	if weights[previous] == 'n' && weights[value(1)] == 't' && attributes[value(2)]&0xc0 == 0 {
		return append([]byte(nil), source[:2]...), true
	}
	for _, suffix := range []string{"s", "d", "m", "ve", "em", "re", "ll"} {
		if len(source) < len(suffix)+1 {
			continue
		}
		if ComparePaul2013MappedCString(source[1:len(suffix)+1], []byte(suffix), weights) == 0 && attributes[value(len(suffix)+1)]&0xc0 == 0 {
			return append([]byte(nil), source[:len(suffix)+1]...), true
		}
	}
	return nil, false
}
