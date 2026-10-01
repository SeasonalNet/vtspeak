package text

import (
	"bytes"
	"fmt"
)

// Paul2013ModelParserWhitespacePrefix is the portion of FUN_1005a350's
// scanner result produced while it skips leading ASCII whitespace.
type Paul2013ModelParserWhitespacePrefix struct {
	ConsumedBytes        int
	LineFeeds            int
	SingleLineFeedMarker bool
	ScannerStatus        int32
	NULTerminated        bool
}

// ScanPaul2013ModelParserWhitespacePrefix ports the leading-whitespace loop
// at FUN_1005a350 (0x1005a384 onward). It consumes spaces, tabs, CR, and LF
// from a NUL-terminated byte string. Exactly one LF sets the scanner's +4
// marker; two or more LFs set status 8. The subsequent token-scanning state
// machine is not represented here.
func ScanPaul2013ModelParserWhitespacePrefix(source []byte) Paul2013ModelParserWhitespacePrefix {
	result := Paul2013ModelParserWhitespacePrefix{}
	for result.ConsumedBytes < len(source) {
		value := source[result.ConsumedBytes]
		if value == 0 {
			result.NULTerminated = true
			break
		}
		switch value {
		case ' ', '\t', '\r':
			result.ConsumedBytes++
		case '\n':
			result.ConsumedBytes++
			result.LineFeeds++
		default:
			if result.LineFeeds == 1 {
				result.SingleLineFeedMarker = true
			}
			if result.LineFeeds >= 2 {
				result.ScannerStatus = 8
			}
			return result
		}
	}
	if result.LineFeeds == 1 {
		result.SingleLineFeedMarker = true
	}
	if result.LineFeeds >= 2 {
		result.ScannerStatus = 8
	}
	return result
}

// Paul2013ProperNameScannerObservation records the scanner counters and
// status relevant to FUN_10034180. Output is supplied for ordinary token
// paths; terminal multiline/empty paths are derived directly from Input.
type Paul2013ProperNameScannerObservation struct {
	Output               []byte
	ResultCount          int32
	PrefixCount          int32
	Status               int32
	SingleLineFeedMarker bool
	Terminal             bool
}

// Paul2013ProperNameASCIITokenObservation records bounded ordinary token
// paths recovered for FUN_1005a350 mode 0x15. Supported results contain an
// ASCII letter word or a digit-based number after the leading whitespace
// prefix. The token cursor stops before the supported native ASCII boundary;
// ConsumedBytes excludes that delimiter.
type Paul2013ProperNameASCIITokenObservation struct {
	Output               []byte
	ResultCount          int32
	PrefixCount          int32
	ConsumedBytes        int32
	Status               int32
	SingleLineFeedMarker bool
	Supported            bool
	Reason               string
}

// ObservePaul2013ProperNameASCIIToken ports bounded mode-0x15 ordinary paths:
// ASCII digit runs, decimal and comma-grouped numbers, numeric ordinal
// suffixes, the native one-digit letter-number-letter forms, and ASCII letter words with internal
// dot/hyphen continuation, plus the exact two-component United States branch.
// Letter-token paths stop at NUL or an ASCII nonletter boundary. A numeric
// token may also stop before one ASCII byte followed by NUL. ConsumedBytes
// excludes either delimiter. Other mixed alphanumeric forms, longer numeric
// lookahead, multibyte input, oversized tokens, and other scanner modes remain
// unsupported.
func ObservePaul2013ProperNameASCIIToken(
	input []byte,
) Paul2013ProperNameASCIITokenObservation {
	nul := bytes.IndexByte(input, 0)
	if nul < 0 {
		return Paul2013ProperNameASCIITokenObservation{Reason: "input is not NUL-terminated"}
	}
	prefix := ScanPaul2013ModelParserWhitespacePrefix(input[:nul+1])
	result := Paul2013ProperNameASCIITokenObservation{
		PrefixCount:          int32(prefix.ConsumedBytes),
		ConsumedBytes:        int32(prefix.ConsumedBytes),
		SingleLineFeedMarker: prefix.SingleLineFeedMarker,
	}
	if prefix.LineFeeds >= 2 {
		result.Status = 8
		result.Supported = true
		result.Reason = "native multiline early return"
		return result
	}
	start := prefix.ConsumedBytes
	if start == nul {
		result.Status = 9
		result.Supported = true
		result.Reason = "native empty-input early return"
		return result
	}
	remaining := input[start:nul]
	if numeric, handled := observePaul2013ASCIIProperNameNumericToken(remaining); handled {
		numeric.PrefixCount = int32(start)
		numeric.ConsumedBytes += int32(start)
		numeric.SingleLineFeedMarker = prefix.SingleLineFeedMarker
		return numeric
	}
	tokenEnd := len(remaining)
	const unitedStatesLength = len("United States")
	if len(remaining) >= unitedStatesLength &&
		(len(remaining) == unitedStatesLength || isPaul2013ASCIIWordBoundary(remaining[unitedStatesLength])) &&
		paul2013ProperNameUnitedStates(remaining[:unitedStatesLength]) {
		tokenEnd = unitedStatesLength
	} else {
		for index, value := range remaining {
			if (value == '.' || value == '-') && index > 0 && index+1 < len(remaining) &&
				isPaul2013ASCIIWordLetter(remaining[index-1]) && isPaul2013ASCIIWordLetter(remaining[index+1]) {
				continue
			}
			if isPaul2013ASCIIDigit(value) || isPaul2013ASCIIWordLetter(value) || value >= 0x80 {
				continue
			}
			if value < 0x80 {
				tokenEnd = index
				break
			}
		}
	}
	for tokenEnd > 0 && (remaining[tokenEnd-1] == '.' || remaining[tokenEnd-1] == '-') {
		tokenEnd--
	}
	token := remaining[:tokenEnd]
	if isPaul2013ASCIIDigit(remaining[0]) && tokenEnd < len(remaining) &&
		(len(remaining)-tokenEnd != 1 || !isPaul2013ASCIIDigitRun(token)) {
		result.Reason = "numeric token delimiter lookahead is outside the recovered one-byte terminal boundary"
		return result
	}
	if len(token) == 0 || len(token) > 28 {
		result.Reason = "token length is outside the directly supported 1..28 byte range"
		return result
	}
	if paul2013ProperNameUnitedStates(token) {
		weightTotal := 0
		for _, value := range token {
			if value == ' ' {
				continue
			}
			letter := value
			if letter >= 'a' && letter <= 'z' {
				letter -= 'a' - 'A'
			}
			weightTotal += 1 + len(paul2013ContextCodeExpansions[letter-'A'])
		}
		if weightTotal > 0x40 {
			result.Reason = "native FUN_10062f50 accumulation leaves the United States scanner path"
			return result
		}
		result.Output = append([]byte(nil), token...)
		result.ResultCount = int32(len(token))
		result.ConsumedBytes = int32(start + len(token))
		result.Status = 1
		result.Supported = true
		return result
	}
	if len(token) >= 3 && isPaul2013ASCIIDigit(token[0]) &&
		isPaul2013ASCIIUpper(token[1]) && isPaul2013ASCIIDigit(token[2]) &&
		isPaul2013ASCIIDigitRun(token[2:]) {
		result.Output = append([]byte(nil), token...)
		result.ResultCount = int32(len(token))
		result.ConsumedBytes = int32(start + len(token))
		result.Status = 2
		result.Supported = true
		return result
	}
	if isPaul2013ASCIIDigitRun(token) {
		result.Output = append([]byte(nil), token...)
		result.ResultCount = int32(len(token))
		result.ConsumedBytes = int32(start + len(token))
		result.Status = 2
		result.Supported = true
		return result
	}
	wordStart := 0
	if len(token) >= 3 && isPaul2013ASCIIUpper(token[0]) &&
		isPaul2013ASCIIDigit(token[1]) && isPaul2013ASCIIUpper(token[2]) {
		wordStart = 2
	}
	weightTotal := 0
	for index := wordStart; index < len(token); index++ {
		value := token[index]
		if (value >= 'A' && value <= 'Z') || (value >= 'a' && value <= 'z') {
			letter := value
			if letter >= 'a' && letter <= 'z' {
				letter -= 'a' - 'A'
			}
			weightTotal += 1 + len(paul2013ContextCodeExpansions[letter-'A'])
			if weightTotal > 0x40 {
				result.Reason = "native FUN_10062f50 accumulation leaves the simple token path"
				return result
			}
			continue
		}
		if (value == '.' || value == '-') && index > 0 && index+1 < len(token) &&
			isPaul2013ASCIIWordLetter(token[index-1]) && isPaul2013ASCIIWordLetter(token[index+1]) {
			continue
		}
		result.Reason = "token is outside the recovered ASCII mode-0x15 token paths"
		return result
	}
	if wordStart == 2 {
		letter := token[0] - 'A'
		weightTotal += 1 + len(paul2013ContextCodeExpansions[letter])
		if weightTotal > 0x40 {
			result.Reason = "native FUN_10062f50 accumulation leaves the simple token path"
			return result
		}
	}
	result.Output = append([]byte(nil), token...)
	result.ResultCount = int32(len(token))
	result.ConsumedBytes = int32(start + len(token))
	result.Status = 1
	result.Supported = true
	return result
}

func isPaul2013ASCIIDigitRun(token []byte) bool {
	for _, value := range token {
		if value < '0' || value > '9' {
			return false
		}
	}
	return true
}

func observePaul2013ASCIIProperNameNumericToken(
	remaining []byte,
) (Paul2013ProperNameASCIITokenObservation, bool) {
	if len(remaining) == 0 || !isPaul2013ASCIIDigit(remaining[0]) {
		return Paul2013ProperNameASCIITokenObservation{}, false
	}
	// Keep the directly supported digit-uppercase-digit word form on its
	// existing branch; this helper covers numeric punctuation only.
	if len(remaining) >= 3 && isPaul2013ASCIIUpper(remaining[1]) && isPaul2013ASCIIDigit(remaining[2]) {
		return Paul2013ProperNameASCIITokenObservation{}, false
	}

	index := 0
	for index < len(remaining) && isPaul2013ASCIIDigit(remaining[index]) {
		index++
	}
	if index == 0 {
		return Paul2013ProperNameASCIITokenObservation{}, false
	}
	if suffixLength := paul2013ASCIIOrdinalSuffixLength(remaining[:index], remaining[index:]); suffixLength != 0 {
		tokenEnd := index + suffixLength
		if tokenEnd < len(remaining) && (len(remaining)-tokenEnd != 1 || remaining[tokenEnd] >= 0x80) {
			return Paul2013ProperNameASCIITokenObservation{
				Reason: "ordinal token delimiter lookahead is outside the recovered one-byte ASCII boundary",
			}, true
		}
		token := remaining[:tokenEnd]
		if len(token) > 28 {
			return Paul2013ProperNameASCIITokenObservation{
				Reason: "numeric token length is outside the directly supported 1..28 byte range",
			}, true
		}
		return Paul2013ProperNameASCIITokenObservation{
			Output:        append([]byte(nil), token...),
			ResultCount:   int32(len(token)),
			ConsumedBytes: int32(tokenEnd),
			Status:        2,
			Supported:     true,
		}, true
	}
	decimalSeen := false
	for index < len(remaining) {
		switch remaining[index] {
		case ',':
			if decimalSeen || index+1 >= len(remaining) || !isPaul2013ASCIIDigit(remaining[index+1]) {
				goto boundary
			}
			if remaining[0] == '0' {
				return Paul2013ProperNameASCIITokenObservation{
					Reason: "native thousands separator path rejects a leading zero",
				}, true
			}
			groupStart := index + 1
			groupEnd := groupStart
			for groupEnd < len(remaining) && isPaul2013ASCIIDigit(remaining[groupEnd]) {
				groupEnd++
			}
			if groupEnd-groupStart != 3 {
				return Paul2013ProperNameASCIITokenObservation{
					Reason: "native comma grouping requires exactly three following digits",
				}, true
			}
			index = groupEnd
		case '.':
			if index+1 >= len(remaining) || !isPaul2013ASCIIDigit(remaining[index+1]) {
				goto boundary
			}
			if decimalSeen {
				return Paul2013ProperNameASCIITokenObservation{
					Reason: "native numeric scanner supports one decimal point on this recovered path",
				}, true
			}
			decimalSeen = true
			index++
			for index < len(remaining) && isPaul2013ASCIIDigit(remaining[index]) {
				index++
			}
		default:
			goto boundary
		}
	}

boundary:
	if index < len(remaining) && (len(remaining)-index != 1 || remaining[index] >= 0x80) {
		return Paul2013ProperNameASCIITokenObservation{
			Reason: "numeric token delimiter lookahead is outside the recovered one-byte ASCII boundary",
		}, true
	}
	token := remaining[:index]
	if len(token) == 0 || len(token) > 28 {
		return Paul2013ProperNameASCIITokenObservation{
			Reason: "numeric token length is outside the directly supported 1..28 byte range",
		}, true
	}
	return Paul2013ProperNameASCIITokenObservation{
		Output:        append([]byte(nil), token...),
		ResultCount:   int32(len(token)),
		ConsumedBytes: int32(index),
		Status:        2,
		Supported:     true,
	}, true
}

func paul2013ASCIIOrdinalSuffixLength(digits, following []byte) int {
	if len(digits) == 0 {
		return 0
	}
	last := digits[len(digits)-1]
	want := "th"
	if len(digits) < 2 || digits[len(digits)-2] != '1' {
		switch last {
		case '1':
			want = "st"
		case '2':
			want = "nd"
		case '3':
			want = "rd"
		}
	}
	if len(digits)+len(want) > 28 {
		return 0
	}
	if len(following) < len(want) || !equalPaul2013MappedASCIISuffix(following[:len(want)], want) {
		return 0
	}
	// The native rd branch also performs a byte-exact comparison against "Rd".
	if want == "rd" && following[0] == 'R' && following[1] == 'd' {
		return 0
	}
	return len(want)
}

func equalPaul2013MappedASCIISuffix(value []byte, expected string) bool {
	if len(value) != len(expected) {
		return false
	}
	for index := range value {
		actual := value[index]
		if actual >= 'A' && actual <= 'Z' {
			actual += 'a' - 'A'
		}
		if actual != expected[index] {
			return false
		}
	}
	return true
}

func isPaul2013ASCIIDigit(value byte) bool {
	return value >= '0' && value <= '9'
}

func isPaul2013ASCIIWordBoundary(value byte) bool {
	return value < 0x80 && !isPaul2013ASCIIWordLetter(value)
}

func isPaul2013ASCIIUpper(value byte) bool {
	return value >= 'A' && value <= 'Z'
}

func paul2013ProperNameUnitedStates(token []byte) bool {
	if len(token) != len("United States") || token[6] != ' ' {
		return false
	}
	weights := Paul2013ContextCharacterWeights()
	return ComparePaul2013MappedCString(token[:6], []byte("United"), weights) == 0 &&
		ComparePaul2013MappedCString(token[7:], []byte("States"), weights) == 0
}

func isPaul2013ASCIIWordLetter(value byte) bool {
	return (value >= 'A' && value <= 'Z') || (value >= 'a' && value <= 'z')
}

// ObservePaul2013ProperNameScannerInput ports FUN_1005a350's leading
// whitespace counters and its status 8/9 early returns. For a nonterminal
// input outside the bounded letter-token path, status and transformed output
// still come from the caller's observed token-scanning result.
func ObservePaul2013ProperNameScannerInput(
	input []byte,
	ordinaryOutput []byte,
	ordinaryStatus int32,
) (Paul2013ProperNameScannerObservation, error) {
	if bytes.IndexByte(input, 0) < 0 {
		return Paul2013ProperNameScannerObservation{}, fmt.Errorf("proper-name scanner input is not NUL-terminated")
	}
	prefix := ScanPaul2013ModelParserWhitespacePrefix(input)
	result := Paul2013ProperNameScannerObservation{
		Output:               append([]byte(nil), cString(ordinaryOutput)...),
		PrefixCount:          int32(prefix.ConsumedBytes),
		Status:               ordinaryStatus,
		SingleLineFeedMarker: prefix.SingleLineFeedMarker,
	}
	if prefix.LineFeeds >= 2 {
		result.Output = nil
		result.Status = 8
		result.Terminal = true
	} else if prefix.ConsumedBytes < len(input) && input[prefix.ConsumedBytes] == 0 {
		result.Output = nil
		result.Status = 9
		result.Terminal = true
	}
	result.ResultCount = int32(len(result.Output))
	return result, nil
}
