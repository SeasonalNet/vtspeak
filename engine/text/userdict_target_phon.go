package text

import (
	"bytes"
	"fmt"
)

const (
	paul2013UserDictTargetPhonMarker    = "[CI]"
	paul2013UserDictTargetPhonMaxBytes  = 260
	paul2013UserDictTargetPhonMaxPhones = 65
)

// CheckPaul2013UserDictionaryTargetPhon ports FUN_1005f5c0, called by
// VT_CheckUserDict_TargetPhon_ENG. buffer must be a writable NUL-terminated
// C-string. The native helper trims space, TAB, LF, and CR in place, accepts
// the recovered uppercase phone inventory, and truncates an accepted terminal
// [CI] marker at its opening bracket. Return values remain native status
// codes; their undocumented meanings are not inferred here.
func CheckPaul2013UserDictionaryTargetPhon(buffer []byte) (int16, error) {
	nul := bytes.IndexByte(buffer, 0)
	if nul < 0 {
		return 0, fmt.Errorf("Paul target-phone buffer is not NUL-terminated")
	}
	start, end := 0, nul
	for start < end && paul2013UserDictionaryTrimByte(buffer[start]) {
		start++
	}
	for end > start && paul2013UserDictionaryTrimByte(buffer[end-1]) {
		end--
	}
	if start != 0 {
		copy(buffer, buffer[start:end])
	}
	trimmedLength := end - start
	buffer[trimmedLength] = 0
	if trimmedLength == 0 {
		return -1, nil
	}
	input := buffer[:trimmedLength]

	phoneCount := 0
	marker := false
	for offset := 0; offset < len(input); {
		for offset < len(input) && (input[offset] == ' ' || input[offset] == '\t') {
			offset++
		}
		if offset == len(input) {
			break
		}
		if input[offset] == '[' {
			if phoneCount == 0 {
				return -2, nil
			}
			if !Paul2013MappedCStringEqual(
				input[offset:], []byte(paul2013UserDictTargetPhonMarker),
				Paul2013EmbeddedKeyTables().CharacterMap,
			) {
				return -12, nil
			}
			input[offset] = 0
			marker = true
			break
		}

		tokenStart := offset
		for offset < len(input) && input[offset] != ' ' && input[offset] != '\t' && input[offset] != '[' {
			offset++
		}
		token := string(input[tokenStart:offset])
		if token == "#" {
			if phoneCount != 0 || marker || offset != len(input) {
				return -9, nil
			}
		} else if _, err := ParsePaul2013CMUPronunciation(token); err != nil {
			return -9, nil
		}
		phoneCount++
	}

	if phoneCount == 0 {
		return -2, nil
	}
	if trimmedLength >= paul2013UserDictTargetPhonMaxBytes {
		return -5, nil
	}
	if marker {
		return 2, nil
	}
	return 1, nil
}

// ConvertPaul2013UserDictionaryTargetPhon ports FUN_1005f710 over the
// recovered spelling-to-byte map. It returns a 66-byte NUL-terminated output
// area and the native 1/0 result. On invalid input or more than 65 tokens,
// byte zero is cleared while any already-written tail bytes are retained.
func ConvertPaul2013UserDictionaryTargetPhon(source []byte) ([]byte, int16) {
	output := make([]byte, paul2013UserDictTargetPhonMaxPhones+1)
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	phoneCount := 0
	for offset := 0; offset < len(source); {
		for offset < len(source) && (source[offset] == ' ' || source[offset] == '\t' || source[offset] == '\n' || source[offset] == '\r') {
			offset++
		}
		if offset == len(source) || source[offset] == 0 {
			return output, 1
		}
		tokenStart := offset
		for offset < len(source) && source[offset] != 0 && source[offset] != ' ' && source[offset] != '\t' && source[offset] != '\n' && source[offset] != '\r' {
			offset++
		}
		token := string(source[tokenStart:offset])
		if token == "#" {
			for offset < len(source) && (source[offset] == ' ' || source[offset] == '\t' || source[offset] == '\n' || source[offset] == '\r') {
				offset++
			}
			if offset != len(source) || phoneCount != 0 {
				output[0] = 0
				return output, 0
			}
			output[0] = 0x64
			output[1] = 0
			return output, 1
		}
		phones, err := ParsePaul2013CMUPronunciation(token)
		if err != nil || len(phones) != 1 {
			output[0] = 0
			return output, 0
		}
		if phoneCount >= paul2013UserDictTargetPhonMaxPhones {
			output[0] = 0
			return output, 0
		}
		features, err := phones[0].TreeFeatures()
		if err != nil {
			output[0] = 0
			return output, 0
		}
		output[phoneCount] = features.SymbolCode
		phoneCount++
		output[phoneCount] = 0
	}
	return output, 1
}
