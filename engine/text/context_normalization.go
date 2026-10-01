package text

import (
	"bytes"
	"fmt"
)

// ApplyPaul2013LiteralContextHandler ports the table-backed contraction and
// suffix cases in FUN_100091b0. It appends their observed opaque code bytes to
// the phone/context row and sets row flag bit 3. The caller supplies the
// surface extracted by the native parser; non-matching strings are unchanged.
// Other words continue through the unported generic dictionary/TPP path.
func ApplyPaul2013LiteralContextHandler(row []byte, surface []byte) (int16, bool, error) {
	if len(row) < Paul2013PhoneContextRowSize {
		return 0, false, fmt.Errorf("phone/context row has %d bytes, need %d", len(row), Paul2013PhoneContextRowSize)
	}
	if nul := bytes.IndexByte(surface, 0); nul >= 0 {
		surface = surface[:nul]
	}
	var appendCodes []byte
	switch string(surface) {
	case "'s":
		last, _, _, err := paul2013PhoneRowCodeTail(row)
		if err != nil {
			return 0, false, err
		}
		switch {
		case int8(last) < 1:
			appendCodes = []byte{0x44}
		case last == ' ' || last == '*' || last == '5' || last == '9' || last == ':':
			appendCodes = []byte{0x37}
		case last == 0x14 || last == ')' || last == '7' || last == '8' || last == 'D' || last == 'E':
			appendCodes = []byte{0x23, 0x44}
		default:
			appendCodes = []byte{0x44}
		}
	case "'ve", "'VE":
		last, _, _, err := paul2013PhoneRowCodeTail(row)
		if err != nil {
			return 0, false, err
		}
		if !paul2013ContextCodeClass(last) {
			appendCodes = []byte{0x07, 0x41}
		} else {
			appendCodes = []byte{0x41}
		}
	case "'re", "'RE":
		last, _, _, err := paul2013PhoneRowCodeTail(row)
		if err != nil {
			return 0, false, err
		}
		if paul2013ContextCodeClass(last) {
			appendCodes = []byte{0x36}
		} else {
			appendCodes = []byte{0x1a}
		}
	case "'em", "'ll", "'LL":
		last, previous, _, err := paul2013PhoneRowCodeTail(row)
		if err != nil {
			return 0, false, err
		}
		appendCodes = paul2013ConditionalSuffixCodes(0x2b, last, previous)
	case "'EM", "'m":
		last, previous, _, err := paul2013PhoneRowCodeTail(row)
		if err != nil {
			return 0, false, err
		}
		appendCodes = paul2013ConditionalSuffixCodes(0x15, last, previous)
	case "'n":
		last, _, _, err := paul2013PhoneRowCodeTail(row)
		if err != nil {
			return 0, false, err
		}
		if paul2013ContextCodeClass(last) {
			appendCodes = []byte{0x2c}
		} else {
			appendCodes = []byte{0x07, 0x2c}
		}
	case "'d":
		appendCodes = []byte{0x26}
	case "'er":
		appendCodes = []byte{0x1a}
	case "'ee":
		appendCodes = []byte{0x1a, 0x44}
	default:
		return 0, false, nil
	}

	_, _, length, err := paul2013PhoneRowCodeTail(row)
	if err != nil {
		return 0, false, err
	}
	codeEnd := Paul2013PhoneContextRowPhoneCodes + length
	if codeEnd+len(appendCodes) >= Paul2013PhoneContextRowSourceMarker {
		return 0, false, fmt.Errorf("literal context handler needs %d code bytes, row has %d bytes available", len(appendCodes), Paul2013PhoneContextRowSourceMarker-codeEnd-1)
	}
	copy(row[codeEnd:codeEnd+len(appendCodes)], appendCodes)
	row[codeEnd+len(appendCodes)] = 0
	row[Paul2013PhoneContextRowFlags] |= 0x08
	return 1, true, nil
}

func paul2013PhoneRowCodeTail(row []byte) (last, previous byte, length int, err error) {
	codeArea := row[Paul2013PhoneContextRowPhoneCodes:Paul2013PhoneContextRowSourceMarker]
	nul := bytes.IndexByte(codeArea, 0)
	if nul < 0 {
		return 0, 0, 0, fmt.Errorf("phone/context row code string is not NUL-terminated")
	}
	length = nul
	if length > 0 {
		last = codeArea[length-1]
	}
	if length > 1 {
		previous = codeArea[length-2]
	}
	return last, previous, length, nil
}

func paul2013ContextCodeClass(value byte) bool {
	if value == 0 || value > 'E' {
		return false
	}
	// DAT_10077d94 contains the 16-bit membership flags indexed by code byte.
	// These ranges are the nonzero entries through 'E', which bounds every
	// native lookup in FUN_100091b0.
	return value <= 0x12 || (value >= 0x17 && value <= 0x1f) ||
		(value >= 0x23 && value <= 0x28) || (value >= 0x2f && value <= 0x34) ||
		(value >= 0x3b && value <= 0x40) || value == 0x43
}

// IsPaul2013ContextCodeClass reports membership in the raw 16-bit table at
// DLL address 0x10077d94 for code bytes through 'E'. The table values remain
// opaque; this helper exposes only the nonzero membership predicate.
func IsPaul2013ContextCodeClass(value byte) bool {
	return paul2013ContextCodeClass(value)
}

// HasPaul2013NormalizerContextClass reports whether a shaped generic
// normalizer result contains a code that FUN_1000ff60 accepts for dictionary
// dispatch. The byte meanings remain opaque.
func HasPaul2013NormalizerContextClass(codes []byte) bool {
	if nul := bytes.IndexByte(codes, 0); nul >= 0 {
		codes = codes[:nul]
	}
	for _, value := range codes {
		if value > 0 && value < 'F' && paul2013ContextCodeClass(value) {
			return true
		}
	}
	return false
}

func paul2013ConditionalSuffixCodes(code, last, previous byte) []byte {
	if !paul2013ContextCodeClass(last) && (!paul2013ContextCodeClass(previous) || last != '6') {
		return []byte{0x07, code}
	}
	return []byte{code}
}
