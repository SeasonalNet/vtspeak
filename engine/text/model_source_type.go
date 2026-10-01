package text

import (
	"bytes"
	"fmt"
)

// ClassifyPaul2013ModelSourceType ports FUN_1000fd20. The source is the
// already transformed C string passed from FUN_1000d190; the caller supplies
// the native character-attribute table.
func ClassifyPaul2013ModelSourceType(source []byte, attributes [256]byte) byte {
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	result := byte(0)
	seenClass := false
	for _, character := range source {
		attribute := attributes[character]
		if attribute&0xc0 == 0 {
			continue
		}
		if attribute&0x80 == 0 {
			if !seenClass {
				result = 1
				seenClass = true
			}
			continue
		}
		if result != 0 {
			result = 3
		} else {
			result = 2
		}
		seenClass = true
	}
	return result
}

// ClassifyPaul2013ASCIIModelSourceType is a strict ASCII convenience wrapper.
// FUN_1000d190 passes this result to FUN_1000d450 for the token-result row's
// type byte. Use ClassifyPaul2013ModelSourceType with the effective signed-char
// attributes from Paul2013ExceptionCharacterAttributes for native byte input.
func ClassifyPaul2013ASCIIModelSourceType(source []byte) (byte, error) {
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	for index, character := range source {
		if character >= 0x80 {
			return 0, fmt.Errorf("model source type input has unsupported non-ASCII byte 0x%02x at offset %d", character, index)
		}
	}
	return ClassifyPaul2013ModelSourceType(source, Paul2013ExceptionCharacterAttributes()), nil
}

// WithPaul2013ModelSourceType returns a copy of the dictionary phone-row
// projection with FUN_1000fd20's output supplied by the transformed model
// source string. The transform itself remains an explicit upstream input.
func (rows Paul2013DictionaryPhoneRows) WithPaul2013ModelSourceType(source []byte) (Paul2013DictionaryPhoneRows, error) {
	resultType, err := ClassifyPaul2013ASCIIModelSourceType(source)
	if err != nil {
		return Paul2013DictionaryPhoneRows{}, err
	}
	rows.ResultType = resultType
	return rows, nil
}
