package text

import (
	"bytes"
	"errors"
)

// NormalizePaul2013UserDictionaryTarget ports the observed input scan in
// FUN_1005f3b0, called by VT_CheckUserDict_TargetNorm_ENG. The argument must
// be a writable NUL-terminated C-string buffer. Leading/trailing ASCII
// whitespace is removed in place; a terminal [CI] marker is also replaced by
// NUL. Return values remain the native opaque signed codes.
func NormalizePaul2013UserDictionaryTarget(buffer []byte) (int16, error) {
	nul := bytes.IndexByte(buffer, 0)
	if nul < 0 {
		return 0, errors.New("Paul user-dictionary target buffer is not NUL-terminated")
	}
	source := buffer[:nul]
	start, end := 0, len(source)
	for start < end && paul2013UserDictionaryTrimByte(source[start]) {
		start++
	}
	for end > start && paul2013UserDictionaryTrimByte(source[end-1]) {
		end--
	}
	if start != 0 || end != len(source) {
		copy(buffer, source[start:end])
		buffer[end-start] = 0
		source = buffer[:end-start]
	}
	if len(source) == 0 {
		return -1, nil
	}

	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	if source[0] == '[' {
		if Paul2013MappedCStringEqual(source, []byte("[SKIP]"), characterMap) {
			return 1, nil
		}
		return -11, nil
	}

	attributes := Paul2013ExceptionCharacterAttributes()
	segmentLength := 0
	segmentCount := 0
	tagOpen := false
	result := int16(1)
	for index := 0; index < len(source); index++ {
		current := source[index]
		if index+1 < len(source) && paul2013UserDictionaryRejectedPair(current, source[index+1]) {
			return -3, nil
		}
		if (attributes[current] & 0xc0) != 0 {
			segmentLength++
		} else if current == '-' && index > 0 && index+1 < len(source) &&
			(attributes[source[index-1]]&0xc0) != 0 && (attributes[source[index+1]]&0xc0) != 0 {
			segmentLength++
		} else if !tagOpen && current == '<' && index+1 < len(source) && (attributes[source[index+1]]&0xc0) != 0 {
			tagOpen = true
		} else if tagOpen && current == '>' && index > 0 && (attributes[source[index-1]]&0xc0) != 0 {
			tagOpen = false
		} else if tagOpen && current == ' ' {
			return -4, nil
		} else if !tagOpen && current == '[' {
			if Paul2013MappedCStringEqual(source[index:], []byte("[CI]"), characterMap) {
				buffer[index] = 0
				result = 2
				break
			}
			return -12, nil
		} else if current == ' ' {
			if segmentLength > 30 {
				return -6, nil
			}
			segmentCount++
			if segmentCount > 10 {
				return -7, nil
			}
			segmentLength = 0
		} else {
			return -2, nil
		}
		if index+1 > 65 {
			return -5, nil
		}
	}

	if segmentLength > 0 {
		segmentCount++
	}
	if segmentCount > 10 {
		return -7, nil
	}
	if tagOpen {
		return -8, nil
	}
	return result, nil
}
