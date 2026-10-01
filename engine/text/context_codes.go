package text

import "fmt"

const paul2013ContextCodeCapacity = 65

// Paul2013ContextCodes retains the opaque code bytes and per-phone delimiter
// markers parsed by FUN_10016c90. MFlags records the separate marker bit that
// the native routine sets for code 'M'; the code-byte meanings remain unknown.
type Paul2013ContextCodes struct {
	Codes                []byte
	PhoneMarkers         []byte
	MFlags               []bool
	StoppedAtUnsupported bool
}

// Paul2013ContextCodeSummary carries the other explicit outputs of
// FUN_10016c90. StateFlags preserves its per-token 0/12 mapping; ModeCode
// preserves the final 5/6/7 dispatch byte without assigning semantic names.
type Paul2013ContextCodeSummary struct {
	StateFlags []byte
	ModeCode   byte
}

// SummarizePaul2013ContextCodeState ports FUN_10016c90's state-to-flag pass
// and final mode dispatch. All inputs are caller-produced fields; their
// linguistic meaning is unresolved.
func SummarizePaul2013ContextCodeState(
	perTokenStateValues []int32,
	finalTokenStateCode int16,
	finalTokenAuxiliaryCode byte,
	terminalPitchValueCount int16,
) (Paul2013ContextCodeSummary, error) {
	if len(perTokenStateValues) == 0 {
		return Paul2013ContextCodeSummary{}, fmt.Errorf("Paul 2013 context code state requires at least one token")
	}
	flags := make([]byte, len(perTokenStateValues))
	for index, value := range perTokenStateValues {
		if value != -1 {
			flags[index] = 12
		}
	}
	mode := byte(7)
	switch {
	case finalTokenStateCode == 3 && terminalPitchValueCount == 0:
		mode = 6
	case finalTokenStateCode != 3 && finalTokenStateCode != 4:
		if finalTokenAuxiliaryCode != 12 {
			mode = 5
		}
	}
	return Paul2013ContextCodeSummary{StateFlags: flags, ModeCode: mode}, nil
}

// ParsePaul2013ContextCodes ports the string scan in FUN_10016c90. The native
// scan examines at most 65 source bytes. Codes 'd' and 'c' mark the preceding
// phone with '1' and '2'; other accepted bytes append a code with marker '0'.
// 'M' also sets its parallel flag. An unsupported byte stops parsing, and a
// sequence that produces no code bytes is rejected as the native caller does.
func ParsePaul2013ContextCodes(source []byte) (Paul2013ContextCodes, error) {
	limit := len(source)
	if limit > paul2013ContextCodeCapacity {
		limit = paul2013ContextCodeCapacity
	}
	result := Paul2013ContextCodes{
		Codes:        make([]byte, 0, limit),
		PhoneMarkers: make([]byte, 0, limit),
		MFlags:       make([]bool, 0, limit),
	}
	for _, value := range source[:limit] {
		if value == 0 {
			break
		}
		switch value {
		case 'd':
			if len(result.PhoneMarkers) != 0 {
				result.PhoneMarkers[len(result.PhoneMarkers)-1] = '1'
			}
			continue
		case 'c':
			if len(result.PhoneMarkers) != 0 {
				result.PhoneMarkers[len(result.PhoneMarkers)-1] = '2'
			}
			continue
		}
		if value == 'M' {
			result.Codes = append(result.Codes, value)
			result.PhoneMarkers = append(result.PhoneMarkers, '0')
			result.MFlags = append(result.MFlags, true)
			continue
		}
		if value > 'E' {
			result.StoppedAtUnsupported = true
			break
		}
		result.Codes = append(result.Codes, value)
		result.PhoneMarkers = append(result.PhoneMarkers, '0')
		result.MFlags = append(result.MFlags, false)
	}
	if len(result.Codes) == 0 {
		return Paul2013ContextCodes{}, fmt.Errorf("Paul 2013 context code sequence produced no phone codes")
	}
	return result, nil
}
