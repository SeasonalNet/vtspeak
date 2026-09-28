package text

import "fmt"

// CMUPhone is the runtime-backed label for one internal phone-symbol byte.
// Stress is meaningful only when Vowel is true.
type CMUPhone struct {
	Label  string
	Stress uint8
	Vowel  bool
}

// TreePhoneFeatures contains the directly observed internal code, identity,
// and stress values derived from one CMU phone. It is not a complete
// decision-tree input vector.
type TreePhoneFeatures struct {
	// SymbolCode is the original internal phone byte from the embedded
	// pronunciation payload. It is distinct from IdentityOrdinal.
	SymbolCode      byte
	IdentityOrdinal int16
	StressValue     int16
}

var paul2013TreePhoneOrder = [...]string{
	"AA", "AE", "AH", "AO", "AW", "AY", "B", "CH", "D", "DH", "EH", "ER", "EY", "F", "G", "HH",
	"IH", "IY", "JH", "K", "L", "M", "N", "NG", "OW", "OY", "P", "R", "S", "SH", "T", "TH",
	"UH", "UW", "V", "W", "Y", "Z", "ZH",
}

var paul2013TreeVowels = map[string]struct{}{
	"AA": {}, "AE": {}, "AH": {}, "AO": {}, "AW": {}, "AY": {}, "EH": {}, "ER": {}, "EY": {},
	"IH": {}, "IY": {}, "OW": {}, "OY": {}, "UH": {}, "UW": {},
}

// TreeFeatures maps a labeled phone to the observed one-based CMU identity
// ordinal and vowel stress value used in selected Paul 2013 tree inputs.
// Neighbor features and call-specific vector positions are not inferred.
func (phone CMUPhone) TreeFeatures() (TreePhoneFeatures, error) {
	ordinal := int16(0)
	for index, label := range paul2013TreePhoneOrder {
		if phone.Label == label {
			ordinal = int16(index + 1)
			break
		}
	}
	if ordinal == 0 {
		return TreePhoneFeatures{}, fmt.Errorf("CMU phone label %q has no Paul 2013 tree ordinal", phone.Label)
	}
	_, isVowel := paul2013TreeVowels[phone.Label]
	if phone.Vowel != isVowel {
		return TreePhoneFeatures{}, fmt.Errorf("CMU phone %q has inconsistent vowel marker", phone.Label)
	}
	if phone.Vowel {
		if phone.Stress > 2 {
			return TreePhoneFeatures{}, fmt.Errorf("CMU vowel %q has unsupported stress %d", phone.Label, phone.Stress)
		}
	} else if phone.Stress != 0 {
		return TreePhoneFeatures{}, fmt.Errorf("CMU consonant %q has nonzero stress %d", phone.Label, phone.Stress)
	}
	code, ok := internalPhoneSymbolCode(phone)
	if !ok {
		return TreePhoneFeatures{}, fmt.Errorf("CMU phone %q stress %d has no internal symbol code", phone.Label, phone.Stress)
	}
	return TreePhoneFeatures{
		SymbolCode:      code,
		IdentityOrdinal: ordinal,
		StressValue:     int16(phone.Stress),
	}, nil
}

func internalPhoneSymbolCode(wanted CMUPhone) (byte, bool) {
	for code, phone := range internalPhoneLabels {
		if phone == wanted {
			return code, true
		}
	}
	return 0, false
}

var internalPhoneLabels = map[byte]CMUPhone{
	0x01: {Label: "AA", Stress: 0, Vowel: true},
	0x02: {Label: "AA", Stress: 1, Vowel: true},
	0x03: {Label: "AA", Stress: 2, Vowel: true},
	0x04: {Label: "AE", Stress: 0, Vowel: true},
	0x05: {Label: "AE", Stress: 1, Vowel: true},
	0x06: {Label: "AE", Stress: 2, Vowel: true},
	0x07: {Label: "AH", Stress: 0, Vowel: true},
	0x08: {Label: "AH", Stress: 1, Vowel: true},
	0x09: {Label: "AH", Stress: 2, Vowel: true},
	0x0a: {Label: "AO", Stress: 0, Vowel: true},
	0x0b: {Label: "AO", Stress: 1, Vowel: true},
	0x0c: {Label: "AO", Stress: 2, Vowel: true},
	0x0d: {Label: "AW", Stress: 0, Vowel: true},
	0x0e: {Label: "AW", Stress: 1, Vowel: true},
	0x0f: {Label: "AW", Stress: 2, Vowel: true},
	0x10: {Label: "AY", Stress: 0, Vowel: true},
	0x11: {Label: "AY", Stress: 1, Vowel: true},
	0x12: {Label: "AY", Stress: 2, Vowel: true},
	0x13: {Label: "B"},
	0x14: {Label: "CH"},
	0x15: {Label: "D"},
	0x16: {Label: "DH"},
	0x17: {Label: "EH", Stress: 0, Vowel: true},
	0x18: {Label: "EH", Stress: 1, Vowel: true},
	0x19: {Label: "EH", Stress: 2, Vowel: true},
	0x1a: {Label: "ER", Stress: 0, Vowel: true},
	0x1b: {Label: "ER", Stress: 1, Vowel: true},
	0x1c: {Label: "ER", Stress: 2, Vowel: true},
	0x1d: {Label: "EY", Stress: 0, Vowel: true},
	0x1e: {Label: "EY", Stress: 1, Vowel: true},
	0x1f: {Label: "EY", Stress: 2, Vowel: true},
	0x20: {Label: "F"},
	0x21: {Label: "G"},
	0x22: {Label: "HH"},
	0x23: {Label: "IH", Stress: 0, Vowel: true},
	0x24: {Label: "IH", Stress: 1, Vowel: true},
	0x25: {Label: "IH", Stress: 2, Vowel: true},
	0x26: {Label: "IY", Stress: 0, Vowel: true},
	0x27: {Label: "IY", Stress: 1, Vowel: true},
	0x28: {Label: "IY", Stress: 2, Vowel: true},
	0x29: {Label: "JH"},
	0x2a: {Label: "K"},
	0x2b: {Label: "L"},
	0x2c: {Label: "M"},
	0x2d: {Label: "N"},
	0x2e: {Label: "NG"},
	0x2f: {Label: "OW", Stress: 0, Vowel: true},
	0x30: {Label: "OW", Stress: 1, Vowel: true},
	0x31: {Label: "OW", Stress: 2, Vowel: true},
	0x32: {Label: "OY", Stress: 0, Vowel: true},
	0x33: {Label: "OY", Stress: 1, Vowel: true},
	0x34: {Label: "OY", Stress: 2, Vowel: true},
	0x35: {Label: "P"},
	0x36: {Label: "R"},
	0x37: {Label: "S"},
	0x38: {Label: "SH"},
	0x39: {Label: "T"},
	0x3a: {Label: "TH"},
	0x3b: {Label: "UH", Stress: 0, Vowel: true},
	0x3c: {Label: "UH", Stress: 1, Vowel: true},
	0x3d: {Label: "UH", Stress: 2, Vowel: true},
	0x3e: {Label: "UW", Stress: 0, Vowel: true},
	0x3f: {Label: "UW", Stress: 1, Vowel: true},
	0x40: {Label: "UW", Stress: 2, Vowel: true},
	0x41: {Label: "V"},
	0x42: {Label: "W"},
	0x43: {Label: "Y"},
	0x44: {Label: "Z"},
	0x45: {Label: "ZH"},
}

// DecodeCMUPhones names the observed internal symbol bytes 0x01 through
// 0x45. Structural and unknown bytes are rejected rather than labeled as
// phonemes.
func DecodeCMUPhones(symbols []byte) ([]CMUPhone, error) {
	phones := make([]CMUPhone, len(symbols))
	for i, symbol := range symbols {
		phone, ok := internalPhoneLabels[symbol]
		if !ok {
			return nil, fmt.Errorf("internal symbol 0x%02x at position %d has no CMU phone label", symbol, i)
		}
		phones[i] = phone
	}
	return phones, nil
}
