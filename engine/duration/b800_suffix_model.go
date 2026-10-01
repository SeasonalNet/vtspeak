package duration

import (
	"bytes"
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013B800AuxiliaryPhoneBase = 0x4295b
	paul2013B800AuxiliaryTokenBase = 0x4293d
	paul2013B800ContextCodesOffset = 0x23
	paul2013B800ContextFlagsOffset = -2
)

type paul2013B800SuffixResult struct {
	Matched bool
	Codes   []byte
	Rule    string
}

// applyPaul2013B800LiteralSuffixModelContextRow ports FUN_1000b800's
// post-dictionary marker dispatch. Its output is the exact opaque code byte
// sequence written at row+0x23; the source and neighboring bytes retain their
// native numeric identities.
func applyPaul2013B800LiteralSuffixModelContextRow(model []byte, rowIndex int) (Paul2013B800ModelContextResult, bool, error) {
	rowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	currentSurface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013B800ModelContextResult{}, false, fmt.Errorf("read B800 suffix row %d: %w", rowIndex, err)
	}
	currentSurface = cStringBytes(currentSurface)
	auxiliaryOffset := paul2013B800AuxiliaryPhoneBase + rowIndex*paul2013ModelContextRowStride
	if auxiliaryOffset < 0 || auxiliaryOffset >= rowStart || rowStart > len(model) {
		return Paul2013B800ModelContextResult{}, false, fmt.Errorf("B800 auxiliary code string starts outside model at %#x", auxiliaryOffset)
	}
	auxiliary := model[auxiliaryOffset:rowStart]
	codeEnd := bytes.IndexByte(auxiliary, 0)
	if codeEnd < 0 {
		return Paul2013B800ModelContextResult{}, false, fmt.Errorf("B800 auxiliary code string at %#x is not NUL-terminated", auxiliaryOffset)
	}
	var last, previous byte
	if codeEnd > 0 {
		last = auxiliary[codeEnd-1]
	}
	if codeEnd > 1 {
		previous = auxiliary[codeEnd-2]
	}
	lastClass := b800ContextClass(last)
	previousClass := b800ContextClass(previous)
	result := paul2013B800SuffixResult{}
	set := func(rule string, codes ...byte) {
		result = paul2013B800SuffixResult{Matched: true, Codes: codes, Rule: rule}
	}

	switch {
	case text.Paul2013ContextMappedCStringEqual(currentSurface, []byte("'s")):
		switch {
		case int8(last) < 1:
			set("mapped-'s", 0x44)
		case last == ' ' || last == '*' || last == '5' || last == '9' || last == ':':
			set("mapped-'s", 0x37)
		case last == 0x14 || last == ')' || last == '7' || last == '8' || last == 'D' || last == 'E':
			set("mapped-'s", 0x23, 0x44)
		default:
			set("mapped-'s", 0x44)
		}
	case bytes.Equal(currentSurface, []byte("n't")), bytes.Equal(currentSurface, []byte("N'T")):
		if lastClass || (previousClass && last == '6') {
			set("n't/N'T", 0x2d, 0x39)
		} else {
			set("n't/N'T", 0x07, 0x2d, 0x39)
		}
	case bytes.Equal(currentSurface, []byte("'ve")), bytes.Equal(currentSurface, []byte("'VE")), bytes.Equal(currentSurface, []byte("'em")):
		if lastClass || (previousClass && last == '6') {
			set("'ve/'VE/'em", 0x41)
		} else {
			set("'ve/'VE/'em", 0x07, 0x41)
		}
	case bytes.Equal(currentSurface, []byte("'EM")), bytes.Equal(currentSurface, []byte("'re")):
		if lastClass || (previousClass && last == '6') {
			set("'EM/'re", 0x2c)
		} else {
			set("'EM/'re", 0x07, 0x2c)
		}
	case bytes.Equal(currentSurface, []byte("'RE")), bytes.Equal(currentSurface, []byte("'ll")):
		if lastClass {
			set("'RE/'ll", 0x36)
		} else {
			set("'RE/'ll", 0x1a)
		}
	case bytes.Equal(currentSurface, []byte("'LL")):
		if lastClass || (previousClass && last == '6') {
			set("'LL", 0x2b)
		} else {
			set("'LL", 0x07, 0x2b)
		}
	case text.Paul2013ContextMappedCStringEqual(currentSurface, []byte("'d")):
		leftClass3, class11, err := paul2013B800SuffixTableMatches(model, rowIndex)
		if err != nil {
			return Paul2013B800ModelContextResult{}, false, err
		}
		if !leftClass3 && !class11 {
			if last == '*' || last == '5' || last == '7' || last == ' ' || last == '"' || last == '8' || last == ':' || last == 0x14 {
				set("mapped-'d-unlisted", 0x39)
			} else if last == '9' || last == 0x15 {
				set("mapped-'d-unlisted", 0x23, 0x15)
			} else {
				set("mapped-'d-unlisted", 0x15)
			}
		} else if lastClass || (previousClass && last == '6') {
			set("mapped-'d-listed", 0x15)
		} else {
			set("mapped-'d-listed", 0x07, 0x15)
		}
	case text.Paul2013ContextMappedCStringEqual(currentSurface, []byte("'m")):
		if lastClass {
			set("mapped-'m", 0x2c)
		} else {
			set("mapped-'m", 0x07, 0x2c)
		}
	default:
		return Paul2013B800ModelContextResult{}, false, nil
	}

	codeOffset := rowStart + paul2013B800ContextCodesOffset
	flagOffset := rowStart + paul2013B800ContextFlagsOffset
	if codeOffset < 0 || codeOffset+len(result.Codes) >= len(model) || flagOffset < 0 || flagOffset >= len(model) {
		return Paul2013B800ModelContextResult{}, false, fmt.Errorf("B800 suffix %q output does not fit model row %d", result.Rule, rowIndex)
	}
	updated := append([]byte(nil), model...)
	copy(updated[codeOffset:], result.Codes)
	updated[codeOffset+len(result.Codes)] = 0
	updated[flagOffset] |= 0x08
	return Paul2013B800ModelContextResult{
		Model: updated, FirstRow: rowIndex, LastRow: rowIndex, JoinedSurface: append([]byte(nil), currentSurface...),
		SuffixRule: result.Rule, Eligible: true, NativeReturnCode: 1, Applied: true,
	}, true, nil
}

func b800ContextClass(value byte) bool {
	return value > 0 && value <= 'E' && text.IsPaul2013ContextCodeClass(value)
}

func paul2013B800SuffixTableMatches(model []byte, rowIndex int) (bool, bool, error) {
	start := paul2013B800AuxiliaryTokenBase + rowIndex*paul2013ModelContextRowStride
	end := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	if start < 0 || start >= end || end > len(model) {
		return false, false, fmt.Errorf("B800 auxiliary token starts outside model at %#x", start)
	}
	token := model[start:end]
	if bytes.IndexByte(token, 0) < 0 {
		return false, false, fmt.Errorf("B800 auxiliary token at %#x is not NUL-terminated", start)
	}
	leftClass3, class11 := text.Paul2013B800SuffixWordTableMatches(cStringBytes(token))
	return leftClass3, class11, nil
}

func cStringBytes(value []byte) []byte {
	if end := bytes.IndexByte(value, 0); end >= 0 {
		return value[:end]
	}
	return value
}
