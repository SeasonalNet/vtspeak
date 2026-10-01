package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
)

// Paul2013ModelScannerResult is the recovered 0x44-byte scanner result used
// by FUN_1003e070. Coordinates are raw dwords +8/+0xc, text length is +0x10,
// status is +0x1c, and Text starts at +0x20. Advance is FUN_1005a350's return.
type Paul2013ModelScannerResult struct {
	Field0        int32  // Opaque scanner result +0, used by terminal handling.
	PrefixFlag    int32  // Scanner result +4.
	Field14       int32  // Opaque scanner result +0x14, read by FUN_1003d3d0.
	Field18       uint32 // Raw +0x18; post-scan lookup writes short +0x1a.
	First, Second uint32
	TextLength    int32
	Status        int32
	Text          []byte
	Advance       int32
}

// Paul2013ModelScanner supplies FUN_1005a350's still-incomplete scanner.
// source is the remaining NUL-terminated input, offset its native absolute
// offset, mode the scanner mode, and modePointer parser-state +0x39e8.
type Paul2013ModelScanner func(context.Context, []byte, int32, int32, uint32) (Paul2013ModelScannerResult, error)

type Paul2013AlternateModelParserResult struct {
	State         []byte
	ReturnValue   int16
	ScannerCalls  int
	ConsumedBytes int
}

// RunPaul2013AlternateModelParser ports FUN_1003e070's scanner/row loop.
// It receives an initialized state, as supplied by FUN_1003d350, and leaves
// inputs unchanged. Status 1 appends A/D rows; status 3 recognizes '<' and
// '[' single-byte openings. Angle contents append A/S rows until '>'; bracket
// contents are skipped until ']'. Appends use flag 0x15. Other tokens return
// -1; a full append returns 0; scanner termination returns 1 iff rows exist.
// Finalization is a separate native stage, and the scanner remains explicit.
func RunPaul2013AlternateModelParser(ctx context.Context, input, state []byte, scan Paul2013ModelScanner) (Paul2013AlternateModelParserResult, error) {
	result := Paul2013AlternateModelParserResult{State: append([]byte(nil), state...)}
	if len(state) < paul2013ModelParserModeOffset+4 || bytes.IndexByte(input, 0) < 0 {
		return result, fmt.Errorf("alternate parser state/input is truncated or unterminated")
	}
	if scan == nil {
		return result, fmt.Errorf("alternate model parser requires the native scanner")
	}
	nul := bytes.IndexByte(input, 0)
	modePointer := binary.LittleEndian.Uint32(state[paul2013ModelParserModeOffset:])
	position := 0
	read := func() (Paul2013ModelScannerResult, error) {
		if err := ctx.Err(); err != nil {
			return Paul2013ModelScannerResult{}, err
		}
		if position < 0 || position > nul {
			return Paul2013ModelScannerResult{}, fmt.Errorf("alternate parser cursor outside input")
		}
		result.ScannerCalls++
		token, err := scan(ctx, input[position:nul+1], int32(position), 0, modePointer)
		if err != nil {
			return token, err
		}
		if token.Advance < 0 || int64(token.Advance) > int64(nul-position) {
			return token, fmt.Errorf("scanner advance %d exceeds available source", token.Advance)
		}
		token.Text = append([]byte(nil), cString(token.Text)...)
		return token, nil
	}
	advance := func(token Paul2013ModelScannerResult) error {
		if token.Advance <= 0 {
			return fmt.Errorf("alternate parser continuation has no scanner progress")
		}
		position += int(token.Advance)
		result.ConsumedBytes = position
		return nil
	}
	appendRow := func(token Paul2013ModelScannerResult, class byte) (bool, error) {
		updated, appended, err := AppendPaul2013ModelSourceRow(result.State, Paul2013ModelSourceRowInput{First: token.First, Second: token.Second, Type: 'A', Class: class, Flag: 0x15, Text: token.Text})
		if err == nil {
			result.State = updated
		}
		return appended, err
	}
	for {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		for position < nul && (input[position] == ' ' || input[position] == '\t' || input[position] == '\n' || input[position] == '\r') {
			position++
		}
		result.ConsumedBytes = position
		token, err := read()
		if err != nil {
			return result, err
		}
		if token.Advance == 0 || token.TextLength == 0 {
			result.ReturnValue = -1
			if int16(binary.LittleEndian.Uint16(result.State[:2])) > 0 {
				result.ReturnValue = 1
			}
			return result, nil
		}
		if token.Status == 1 {
			appended, err := appendRow(token, 'D')
			if err != nil {
				return result, err
			}
			if !appended {
				return result, nil
			}
			if err := advance(token); err != nil {
				return result, err
			}
			continue
		}
		if token.Status != 3 || token.TextLength != 1 || len(token.Text) == 0 || (token.Text[0] != '<' && token.Text[0] != '[') {
			result.ReturnValue = -1
			return result, nil
		}
		opening := token.Text[0]
		for {
			if err := advance(token); err != nil {
				return result, err
			}
			token, err = read()
			if err != nil {
				return result, err
			}
			closing := byte('>')
			if opening == '[' {
				closing = ']'
			}
			if len(token.Text) > 0 && token.Text[0] == closing && token.TextLength == 1 {
				break
			}
			if opening == '<' {
				appended, err := appendRow(token, 'S')
				if err != nil {
					return result, err
				}
				if !appended {
					return result, nil
				}
			}
		}
		if err := advance(token); err != nil {
			return result, err
		}
	}
}

// DispatchPaul2013ModelParserWithAlternateScanner supplies the recovered
// alternate parser when the caller has not overridden it. Primary parsing
// and finalization retain the existing dispatcher contracts.
func DispatchPaul2013ModelParserWithAlternateScanner(ctx context.Context, input, state, modeProbe []byte, callbacks Paul2013ModelParserCallbacks, scan Paul2013ModelScanner) (Paul2013ModelParserDispatchResult, error) {
	if callbacks.Alternate == nil {
		callbacks.Alternate = func(input, state []byte) (int32, error) {
			result, err := RunPaul2013AlternateModelParser(ctx, input, state, scan)
			if err == nil {
				copy(state, result.State)
			}
			return int32(result.ReturnValue), err
		}
	}
	return DispatchPaul2013ModelParser(input, state, modeProbe, callbacks)
}
