package text

import (
	"fmt"
)

const paul2013ContextCodeOutputCapacity = 0x40

// FUN_10009cd0's 10-byte records at DLL VA 0x10077e58 contain a little-endian
// length followed by up to eight output bytes. These are retained as opaque
// model codes; the first 26 records correspond to A through Z.
var paul2013ContextCodeExpansions = [...][]byte{
	{0x1e},
	{0x13, 0x27},
	{0x37, 0x27},
	{0x15, 0x27},
	{0x27},
	{0x18, 0x20},
	{0x29, 0x27},
	{0x1e, 0x14},
	{0x11},
	{0x29, 0x1e},
	{0x2a, 0x1e},
	{0x18, 0x2b},
	{0x18, 0x2c},
	{0x18, 0x2d},
	{0x30},
	{0x35, 0x27},
	{0x2a, 0x43, 0x3f},
	{0x02, 0x36},
	{0x18, 0x37},
	{0x39, 0x27},
	{0x43, 0x3f},
	{0x41, 0x27},
	{0x15, 0x08, 0x13, 0x07, 0x2b, 0x43, 0x3e},
	{0x18, 0x2a, 0x37},
	{0x42, 0x11},
	{0x44, 0x27},
}

// EncodePaul2013ContextString ports FUN_10009cd0 for ASCII input. Letters
// expand through the DLL table above; adjacent letter groups and the observed
// . _ - ' separators use the native 'd' delimiter byte. The 64-byte native
// output cap is applied before the terminating NUL, which is omitted here.
// High-bit bytes are rejected because their ANSI-table behavior is outside
// the current ASCII frontend contract.
func EncodePaul2013ContextString(source []byte) ([]byte, error) {
	attributes := Paul2013ExceptionCharacterAttributes()
	output := make([]byte, 0, paul2013ContextCodeOutputCapacity)
	for index := 0; index < len(source); index++ {
		value := source[index]
		if value == 0 {
			break
		}
		if value >= 0x80 {
			return nil, fmt.Errorf("context string byte 0x%02x at offset %d is outside the recovered ASCII mapping", value, index)
		}
		if attributes[value]&0xc0 == 0 {
			if paul2013ContextCodeDelimiter(value) && index+1 < len(source) && attributes[source[index+1]]&0xc0 != 0 {
				if len(output) < paul2013ContextCodeOutputCapacity {
					output = append(output, 'd')
				} else {
					break
				}
			}
			continue
		}

		letter := value
		if letter >= 'a' && letter <= 'z' {
			letter -= 'a' - 'A'
		}
		if letter < 'A' || letter > 'Z' {
			return nil, fmt.Errorf("context string byte 0x%02x at offset %d has no recovered ASCII expansion", value, index)
		}
		expansion := paul2013ContextCodeExpansions[letter-'A']
		for _, code := range expansion {
			if len(output) > 0x3f {
				break
			}
			output = append(output, code)
		}
		if len(output) > 0x3f {
			break
		}
		if index+1 >= len(source) || source[index+1] == 0 || len(output) > 0x3e {
			break
		}
		next := source[index+1]
		if next >= 0x80 {
			return nil, fmt.Errorf("context string byte 0x%02x at offset %d is outside the recovered ASCII mapping", next, index+1)
		}
		nextIsLetter := attributes[next]&0xc0 != 0
		if !nextIsLetter && (!paul2013ContextCodeDelimiter(next) || index+2 >= len(source) || source[index+2] == 0 || source[index+2] >= 0x80 || attributes[source[index+2]]&0xc0 == 0) {
			break
		}
		output = append(output, 'd')
	}
	return output, nil
}

// BuildPaul2013ContextCodesFromSourceString composes FUN_10009cd0's ASCII
// expansion with FUN_10016c90's delimiter scan. It produces per-phone marker
// bytes but does not infer the separate token-state arrays or final mode code.
func BuildPaul2013ContextCodesFromSourceString(source []byte) (Paul2013ContextCodes, error) {
	encoded, err := EncodePaul2013ContextString(source)
	if err != nil {
		return Paul2013ContextCodes{}, err
	}
	codes, err := ParsePaul2013ContextCodes(encoded)
	if err != nil {
		return Paul2013ContextCodes{}, fmt.Errorf("parse normalized Paul 2013 context string: %w", err)
	}
	return codes, nil
}

// ApplyPaul2013ContextStringToPhoneRow writes FUN_10009cd0's output into the
// phone-code area and applies its +0x20 row-flag update when output is nonempty.
func ApplyPaul2013ContextStringToPhoneRow(row []byte, source []byte) error {
	if len(row) < Paul2013PhoneContextRowSize {
		return fmt.Errorf("phone/context row has %d bytes, need %d", len(row), Paul2013PhoneContextRowSize)
	}
	codes, err := EncodePaul2013ContextString(source)
	if err != nil {
		return err
	}
	if len(codes) >= paul2013PhoneContextPhoneCapacity {
		return fmt.Errorf("encoded context string has %d bytes, exceeds phone-row capacity %d", len(codes), paul2013PhoneContextPhoneCapacity-1)
	}
	codeStart := Paul2013PhoneContextRowPhoneCodes
	copy(row[codeStart:codeStart+len(codes)], codes)
	row[codeStart+len(codes)] = 0
	if len(codes) != 0 {
		row[Paul2013PhoneContextRowFlags] |= 0x20
	}
	return nil
}

func paul2013ContextCodeDelimiter(value byte) bool {
	switch value {
	case '.', '_', '-', '\'':
		return true
	default:
		return false
	}
}
