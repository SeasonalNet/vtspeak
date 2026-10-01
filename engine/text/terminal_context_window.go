package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
)

// CopyPaul2013BackwardContext ports FUN_10063300. Count limits whitespace
// encounters rather than lexical tokens; consecutive spaces count separately.
// Consumed includes stripped trailing whitespace, while Text does not.
func CopyPaul2013BackwardContext(source []byte, position, count int) (text []byte, consumed int16, err error) {
	if position < 0 || position > len(source) || count <= 0 || count > 71582788 {
		return nil, 0, fmt.Errorf("backward context position/count outside supported bounds")
	}
	if position == 0 {
		return []byte{0}, 0, nil
	}
	space := func(value byte) bool { return value == ' ' || value == '\t' || value == '\n' || value == '\r' }
	index := position - 1
	trailing := 0
	for space(source[index]) {
		if index == 0 {
			return []byte{0}, 0, nil
		}
		index--
		trailing++
	}
	end := index + 1
	length, encounters := 1, 1
	limit := 30*count - 1
	for index > 0 {
		if length >= limit || space(source[index]) && encounters >= count {
			index++
			length--
			break
		}
		if space(source[index]) {
			encounters++
		}
		length++
		index--
	}
	// The native copy excludes trailing whitespace and appends a zero byte.
	text = append(append([]byte(nil), source[index:end]...), 0)
	return text, int16(length + trailing), nil
}

// ScanPaul2013BackwardContext ports FUN_10063660's three-record rolling
// window. Scanner return values are ignored: progress uses result +0x14.
// Nonprogressing or out-of-bounds scanner results fail explicitly.
func ScanPaul2013BackwardContext(ctx context.Context, source []byte, offset, mode int32, scan Paul2013ModelScanner) (records [3]Paul2013ModelScannerResult, accepted bool, err error) {
	end := bytes.IndexByte(source, 0)
	if end < 0 {
		return records, false, fmt.Errorf("backward scanner source is not NUL-terminated")
	}
	if end == 0 {
		return records, false, nil
	}
	if scan == nil {
		return records, false, fmt.Errorf("backward context scanner is unavailable")
	}
	for position := 0; position < end; {
		if err := ctx.Err(); err != nil {
			return records, false, err
		}
		token, err := scan(ctx, source[position:end+1], offset+int32(position), mode, 0)
		if err != nil {
			return records, false, err
		}
		if token.Field14 <= 0 || int64(token.Field14) > int64(end-position) {
			return records, false, fmt.Errorf("backward scanner advance %d outside remaining source", token.Field14)
		}
		token.Text = append([]byte(nil), cString(token.Text)...)
		records[2] = token
		position += int(token.Field14)
		if position < end {
			if records[1].TextLength != 0 {
				records[0] = records[1]
			}
			if records[2].TextLength != 0 {
				records[1] = records[2]
			}
		}
	}
	return records, true, nil
}

// Paul2013TerminalModelLookup represents FUN_10052520 plus FUN_10001670.
// Only the low byte of its result is tested. Key production and table lookup
// remain explicit until their resource binding has been recovered.
type Paul2013TerminalModelLookup func(context.Context, Paul2013TerminalContextWindow) (uint32, error)

// NewPaul2013TerminalContextRecognizer ports FUN_10051cc0's window assembly
// and classifier dispatch. FullSource is the NUL-terminated utterance whose
// first byte has source coordinate zero. It is copied for callback ownership.
// Ordinary token scanning and the model lookup remain explicit dependencies.
func NewPaul2013TerminalContextRecognizer(fullSource []byte, scan Paul2013ModelScanner, lookup Paul2013TerminalModelLookup) Paul2013TerminalPunctuationRecognizer {
	fullSource = append([]byte(nil), fullSource...)
	return func(ctx context.Context, state []byte, current *Paul2013ModelScannerResult, source []byte, offset int32) (int32, error) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		end := bytes.IndexByte(fullSource, 0)
		if end < 0 || offset < 0 || int64(offset) > int64(end) || current == nil || len(state) < 0x39ec {
			return 0, fmt.Errorf("terminal context source, coordinate, token or state is unavailable")
		}
		position := int(offset)
		if !bytes.Equal(cString(source), fullSource[position:end]) {
			return 0, fmt.Errorf("terminal source does not match the bound utterance at %d", offset)
		}
		backward, consumed, err := CopyPaul2013BackwardContext(fullSource, position, 3)
		if err != nil {
			return 0, err
		}
		var window Paul2013TerminalContextWindow
		// Native window projection copies just +0, +0x10, +0x1c and the C string.
		project := func(token Paul2013ModelScannerResult) Paul2013ModelScannerResult {
			return Paul2013ModelScannerResult{Field0: token.Field0, TextLength: token.TextLength, Status: token.Status, Text: append([]byte(nil), cString(token.Text)...)}
		}
		if consumed != 0 {
			previous, accepted, err := ScanPaul2013BackwardContext(ctx, backward, offset-int32(consumed), 1, scan)
			if err != nil {
				return 0, err
			}
			if !accepted {
				return 'N', nil
			}
			for index, token := range previous {
				if token.TextLength > 0 {
					window[index] = project(token)
				}
			}
		}
		window[3] = project(*current)
		if current.Field14 < 0 || int64(current.Field14) > int64(end-position) || scan == nil {
			return 0, fmt.Errorf("terminal forward scanner/advance is unavailable")
		}
		following := position + int(current.Field14)
		modePointer := binary.LittleEndian.Uint32(state[0x39e8:])
		next, err := scan(ctx, fullSource[following:end+1], int32(following), 1, modePointer)
		if err != nil {
			return 0, err
		}
		if next.TextLength > 0 {
			window[4] = project(next)
			if next.Field14 < 0 || int64(next.Field14) > int64(end-following) {
				return 0, fmt.Errorf("terminal second forward scanner advance outside source")
			}
			second := following + int(next.Field14)
			after, err := scan(ctx, fullSource[second:end+1], int32(second), 1, modePointer)
			if err != nil {
				return 0, err
			}
			if after.TextLength > 0 {
				window[5] = project(after)
			}
		}
		value, err := ClassifyPaul2013TerminalContext(window, fullSource[following:end+1])
		if err != nil || value != 'A' {
			return int32(value), err
		}
		if lookup == nil {
			return 0, fmt.Errorf("terminal context model key/table lookup is unavailable")
		}
		result, err := lookup(ctx, window)
		if err != nil {
			return 0, err
		}
		if byte(result) != 0 {
			return 'P', nil
		}
		return 'N', nil
	}
}
