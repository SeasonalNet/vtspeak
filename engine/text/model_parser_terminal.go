package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
)

// Paul2013TerminalPunctuationRecognizer is FUN_10051cc0's remaining context
// recognizer. It can update the shared scanner result; the caller tests only
// its low short against 0x50, preserving the native contract.
type Paul2013TerminalPunctuationRecognizer func(context.Context, []byte, *Paul2013ModelScannerResult, []byte, int32) (int32, error)

// NewPaul2013TerminalParserHandler ports FUN_10051a00 and adapts it to the
// primary handler map at address 0x10051a00. It rescans in mode 1, replacing
// the shared scanner result, and derives prior-row eligibility from native
// row text/length. The punctuation context recognizer remains explicit and
// is required only on the single .?!; and exact "..." branches.
func NewPaul2013TerminalParserHandler(scan Paul2013ModelScanner, recognize Paul2013TerminalPunctuationRecognizer) Paul2013PrimaryParserHandler {
	return func(ctx context.Context, state []byte, token *Paul2013ModelScannerResult, source []byte, offset int32, _ int16) (int32, error) {
		if token == nil || scan == nil || len(state) < 0x39ec {
			return 0, fmt.Errorf("terminal parser scanner/state is unavailable")
		}
		count, err := paul2013SourceRowCount(state)
		if err != nil {
			return 0, err
		}
		rescanned, err := scan(ctx, source, offset, 1, binary.LittleEndian.Uint32(state[0x39e8:]))
		if err != nil {
			return 0, err
		}
		rescanned.Text = append([]byte(nil), cString(rescanned.Text)...)
		*token = rescanned
		var previous []byte
		if count > 0 {
			start := 0x14 + (count-1)*0x94
			if start+0x94 > len(state) {
				return 0, fmt.Errorf("terminal parser prior row is truncated")
			}
			previous = state[start : start+0x94]
		}
		previousAccepted := func() (bool, error) {
			if count == 0 {
				return false, nil
			}
			length := int32(binary.LittleEndian.Uint32(previous[8:12]))
			if length < 0 || int64(length) > int64(len(previous)-0x34) {
				return false, fmt.Errorf("terminal prior-row length %d outside row", length)
			}
			return !Paul2013ParserPreviousRowTerminalPunctuation(previous[0x34:], int(length)), nil
		}
		writeType := func() {
			if count == 0 || binary.LittleEndian.Uint32(previous[0x2c:]) != 0 {
				return
			}
			value := uint32(2)
			if len(token.Text) > 0 {
				if token.Text[0] == '?' {
					value = 3
				} else if token.Text[0] == '!' {
					value = 4
				}
			}
			binary.LittleEndian.PutUint32(previous[0x2c:], value)
		}
		if token.Status == 8 {
			if count > 0 {
				accepted, err := previousAccepted()
				if err != nil {
					return 0, err
				}
				if accepted {
					writeType()
				}
			}
			if count == 0 {
				return -token.Field14, nil
			}
			return token.Field14, nil
		}
		punctuation := len(token.Text) == 1 && bytes.ContainsAny(token.Text, ".?!;") || bytes.Equal(token.Text, []byte("..."))
		if punctuation {
			if recognize == nil {
				return 0, fmt.Errorf("terminal punctuation recognizer 0x10051cc0 is unavailable")
			}
			classification, err := recognize(ctx, state, token, source, offset)
			if err != nil {
				return 0, err
			}
			if uint16(classification) == 0x50 {
				if count > 0 {
					writeType()
					return token.Field14, nil
				}
				return -token.Field14, nil
			}
			if count > 0 && binary.LittleEndian.Uint32(previous[4:8]) == uint32(token.Field0+offset) && binary.LittleEndian.Uint32(previous[0x28:]) != 0 {
				return -token.Field14, nil
			}
			if token.Field0 == 0 && token.TextLength < 3 && len(token.Text) > 0 && token.Text[0] == '.' {
				return -token.Field14, nil
			}
			return 0, nil
		}
		if token.Status == 9 {
			// Native's conditional zero-to-zero row-type store has no change.
			if count > 0 {
				if _, err := previousAccepted(); err != nil {
					return 0, err
				}
			}
			if token.Field14 == 0 {
				return 1, nil
			}
			return token.Field14, nil
		}
		if token.Status == 3 && token.TextLength == 2 && len(token.Text) >= 2 && token.Text[0] == 0xa2 && token.Text[1] == 0xfd {
			if count > 0 {
				if _, err := previousAccepted(); err != nil {
					return 0, err
				}
			}
			return token.Field14, nil
		}
		return 0, nil
	}
}

// BindPaul2013TerminalParserHandler copies the handler map and provisions
// the recovered terminal body only when no explicit override exists.
func BindPaul2013TerminalParserHandler(handlers map[uint32]Paul2013PrimaryParserHandler, scan Paul2013ModelScanner, recognize Paul2013TerminalPunctuationRecognizer) map[uint32]Paul2013PrimaryParserHandler {
	result := make(map[uint32]Paul2013PrimaryParserHandler, len(handlers)+1)
	for address, handler := range handlers {
		result[address] = handler
	}
	if result[0x10051a00] == nil {
		result[0x10051a00] = NewPaul2013TerminalParserHandler(scan, recognize)
	}
	return result
}
