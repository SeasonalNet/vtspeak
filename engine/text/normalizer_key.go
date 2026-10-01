package text

import (
	"bytes"
	"fmt"
)

// paul2013NormalizerClassRemap is DAT_10077f5a for category indexes 0 through
// 'Y'. The native routine reads the low byte of each little-endian short.
var paul2013NormalizerClassRemap = [90]uint16{
	0, 0xffff, 1, 2, 3, 4, 5, 6, 7, 0xffff, 0xffff, 0xffff, 0xffff, 8, 9, 10,
	11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 30, 31, 32, 33, 0xffff, 0xffff, 34, 0xffff, 35, 36, 37, 38,
	39, 40, 41, 42, 0xffff, 0xffff, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 58, 0xffff, 59, 60, 61, 62, 63, 64, 65, 66, 0xffff,
	0xffff, 67, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 68, 69,
}

// RemapPaul2013NormalizerClassCodes ports FUN_100024f0's explicit category
// remapping, then applies FUN_10002200's byte-level shaping. A false result
// corresponds to the native zero return. Sentinel entries in its table are
// handled by the explicit cases before fallback lookup.
func RemapPaul2013NormalizerClassCodes(source []byte) ([]byte, bool, error) {
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	mapped := make([]byte, 0, len(source)*2)
	for _, value := range source {
		switch value {
		case 1:
			continue
		case 9:
			mapped = append(mapped, 7, '+')
		case 10:
			mapped = append(mapped, 7, ',')
		case 11:
			mapped = append(mapped, 7, '-')
		case 12:
			mapped = append(mapped, 0x1a)
		case '\'':
			mapped = append(mapped, '!', 'D')
		case '(':
			mapped = append(mapped, '!', 'E')
		case '*':
			mapped = append(mapped, '"', 'B')
		case '3':
			mapped = append(mapped, '*', '7')
		case '4':
			mapped = append(mapped, '*', '8')
		case 'E':
			mapped = append(mapped, '9', '7')
		case 'N':
			mapped = append(mapped, 'B', 8)
		case 'O':
			mapped = append(mapped, 'B', 9)
		case 'Q':
			mapped = append(mapped, 'C', 7)
		case 'R':
			mapped = append(mapped, 'C', ';')
		case 'S':
			mapped = append(mapped, 'C', '<')
		case 'T':
			mapped = append(mapped, 'C', '=')
		case 'U':
			mapped = append(mapped, 'C', '>')
		case 'V':
			mapped = append(mapped, 'C', '?')
		case 'W':
			mapped = append(mapped, 'C', '@')
		default:
			if value < 1 || value > 'Y' {
				return nil, false, nil
			}
			mappedValue := paul2013NormalizerClassRemap[value]
			if mappedValue == 0xffff {
				return nil, false, fmt.Errorf("normalizer category 0x%02x maps through an uncharacterized table hole", value)
			}
			mapped = append(mapped, byte(mappedValue))
		}
	}
	shaped, ok := ShapePaul2013NormalizerClassCodes(mapped)
	return shaped, ok, nil
}

// ShapePaul2013NormalizerClassCodes ports FUN_10002200's in-place shaping of
// the opaque class-byte sequence produced by the ATMT normalizer. It returns
// false for a nil or empty C string. The result bytes retain their opaque
// representation; this function does not perform dictionary lookup.
func ShapePaul2013NormalizerClassCodes(source []byte) ([]byte, bool) {
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) == 0 {
		return nil, false
	}

	// FUN_10002200 first rewrites selected class bytes to the literal '6'
	// escape when preceded by one of its directly tested code ranges.
	firstPass := make([]byte, len(source))
	rewritten := false
	for index, value := range source {
		if index > 0 && value >= 0x1a && value <= 0x1c &&
			paul2013NormalizerEscapePrevious(source[index-1]) {
			firstPass[index] = '6'
			rewritten = true
		} else {
			firstPass[index] = value
		}
	}
	if rewritten {
		source = firstPass
	}

	hasContextCode := false
	for _, value := range source {
		if paul2013ContextCodeClass(value) {
			hasContextCode = true
			break
		}
	}

	// The native local buffer and loop guard allow output positions 0 through
	// 64; two-byte escapes may finish at position 65.
	output := make([]byte, 0, 66)
	for index := 0; index < len(source) && len(output) <= 64; index++ {
		value := source[index]
		last := index == len(source)-1
		nextIsContextCode := !last && paul2013ContextCodeClass(source[index+1])
		if paul2013NormalizerDirectPunctuation(value) && (last || nextIsContextCode) {
			switch value {
			case '#':
				output = append(output, '&')
			case '$':
				output = append(output, '\'')
			case '%':
				output = append(output, '(')
			case ';':
				output = append(output, '>')
			case '<':
				output = append(output, '?')
			case '=':
				output = append(output, '@')
			}
			continue
		}

		if last && (value == 'B' || value == 'C') {
			if value == 'B' {
				output = append(output, '>')
			} else {
				output = append(output, '&')
			}
			continue
		}

		if value >= 7 && value <= 9 && !last {
			next := source[index+1]
			switch {
			case next == '6':
				output = append(output, 0x1a+value-7)
				index++
			case next >= 0x1a && next <= 0x1c:
				output = append(output, next)
				index++
			default:
				output = append(output, value)
			}
			continue
		}

		priorIndex := index - 1
		needsSeparator := priorIndex >= 0 && value >= '+' && value <= '-' &&
			!paul2013ContextCodeClass(source[priorIndex])
		if needsSeparator {
			suppressSeparator := priorIndex > 0 &&
				(source[priorIndex] == '+' || source[priorIndex] == '6') &&
				paul2013ContextCodeClass(source[priorIndex-1])
			if !last && hasContextCode {
				suppressSeparator = true
			}
			if !suppressSeparator {
				output = append(output, 7)
			}
		}
		output = append(output, value)
	}
	return output, true
}

func paul2013NormalizerEscapePrevious(value byte) bool {
	return (value >= 1 && value <= 6) ||
		(value >= 10 && value <= 12) ||
		(value >= 0x17 && value <= 0x19) ||
		(value >= '#' && value <= '%') ||
		(value >= ';' && value <= '=')
}

func paul2013NormalizerDirectPunctuation(value byte) bool {
	return (value >= '#' && value <= '%') || (value >= ';' && value <= '=')
}
