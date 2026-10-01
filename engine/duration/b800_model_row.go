package duration

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013B800ModelContextFlagBackstep = 2
	paul2013B800ModelContextPhoneOffset  = text.Paul2013PhoneContextRowPhoneCodes - 2
)

// Paul2013B800ModelContextResult reports FUN_1000b800 for one caller-selected
// context row. NeedsUnportedSuffixDispatch is retained for API compatibility;
// the recovered literal suffix dispatch is now included in this handler.
type Paul2013B800ModelContextResult struct {
	Model                       []byte
	FirstRow                    int
	LastRow                     int
	JoinedSurface               []byte
	SuffixRule                  string
	Eligible                    bool
	NativeReturnCode            int16
	Applied                     bool
	NeedsUnportedSuffixDispatch bool
}

// ApplyPaul2013B800JoinedDictionaryModelContextRow ports the supported
// FUN_1000b800 control flow. It joins the current row with the maximal
// preceding run whose native comparison words match, calls the embedded
// dictionary normalizer, then applies the recovered literal and table-backed
// suffix rules on a dictionary miss. A miss from every branch returns native
// code zero without changing the model.
func (engine *Engine) ApplyPaul2013B800JoinedDictionaryModelContextRow(
	ctx context.Context,
	model []byte,
	contextRowIndex int,
) (Paul2013B800ModelContextResult, error) {
	result := Paul2013B800ModelContextResult{
		Model: append([]byte(nil), model...), FirstRow: contextRowIndex,
		LastRow: contextRowIndex,
	}
	if engine == nil || engine.dictionary == nil {
		return Paul2013B800ModelContextResult{}, errors.New("Paul 2013 B800 dictionary resources are nil or unloaded")
	}
	if ctx == nil {
		return Paul2013B800ModelContextResult{}, errors.New("FUN_1000b800 row handler has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013B800ModelContextResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013B800ModelContextResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if contextRowIndex < 1 || contextRowIndex >= rowCount {
		return result, nil
	}
	rowStart := paul2013ModelContextRowBase + contextRowIndex*paul2013ModelContextRowStride
	comparisonOffset := rowStart - paul2013ModelContextRowStride
	if binary.LittleEndian.Uint16(model[comparisonOffset:comparisonOffset+2]) !=
		binary.LittleEndian.Uint16(model[rowStart:rowStart+2]) {
		return result, nil
	}
	currentSurface, _, err := text.Paul2013ModelContextSurface(model, contextRowIndex)
	if err != nil {
		return Paul2013B800ModelContextResult{}, fmt.Errorf("read B800 model context row %d: %w", contextRowIndex, err)
	}
	if len(currentSurface) == 1 {
		switch currentSurface[0] {
		case ',', '.', '!', '?':
			return result, nil
		}
	}

	firstRow := contextRowIndex - 1
	for firstRow > 0 {
		if err := ctx.Err(); err != nil {
			return Paul2013B800ModelContextResult{}, err
		}
		candidateStart := paul2013ModelContextRowBase + firstRow*paul2013ModelContextRowStride
		previousStart := candidateStart - paul2013ModelContextRowStride
		if binary.LittleEndian.Uint16(model[previousStart:previousStart+2]) !=
			binary.LittleEndian.Uint16(model[candidateStart:candidateStart+2]) {
			break
		}
		firstRow--
	}
	joined := make([]byte, 0, 32)
	for rowIndex := firstRow; rowIndex <= contextRowIndex; rowIndex++ {
		surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
		if err != nil {
			return Paul2013B800ModelContextResult{}, fmt.Errorf("read B800 joined row %d: %w", rowIndex, err)
		}
		if len(joined)+len(surface) >= 32 {
			return Paul2013B800ModelContextResult{}, fmt.Errorf("B800 joined surface requires %d bytes; native local buffer holds 31 bytes plus NUL", len(joined)+len(surface))
		}
		joined = append(joined, surface...)
	}
	result.FirstRow = firstRow
	result.LastRow = contextRowIndex
	result.JoinedSurface = append([]byte(nil), joined...)
	result.Eligible = true

	writes, err := engine.NormalizePaul2013ContextTokenFromEmbeddedDictionaryToBuffer(
		ctx, joined, make([]byte, 68),
	)
	if err != nil {
		return Paul2013B800ModelContextResult{}, fmt.Errorf("normalize B800 joined surface: %w", err)
	}
	result.NativeReturnCode = writes.Result.NativeReturnCode
	if result.NativeReturnCode != 1 {
		result.NativeReturnCode = 0
		suffix, matched, err := applyPaul2013B800LiteralSuffixModelContextRow(result.Model, contextRowIndex)
		if err != nil {
			return Paul2013B800ModelContextResult{}, fmt.Errorf("apply B800 literal suffix: %w", err)
		}
		if matched {
			suffix.FirstRow = firstRow
			suffix.LastRow = contextRowIndex
			suffix.JoinedSurface = append([]byte(nil), joined...)
			return suffix, nil
		}
		return result, nil
	}
	end := bytes.IndexByte(writes.Output, 0)
	if end < 0 {
		return Paul2013B800ModelContextResult{}, errors.New("B800 dictionary output is not NUL-terminated")
	}
	updated, err := text.ApplyPaul2013ModelContextCodes(
		model, firstRow, writes.Output[:end], 0x08,
	)
	if err != nil {
		return Paul2013B800ModelContextResult{}, fmt.Errorf("write B800 joined output on row %d: %w", firstRow, err)
	}
	for rowIndex := firstRow + 1; rowIndex <= contextRowIndex; rowIndex++ {
		rowFlagsOffset := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride - paul2013B800ModelContextFlagBackstep
		rowPhoneOffset := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride + paul2013B800ModelContextPhoneOffset
		updated[rowFlagsOffset] |= 0x08
		updated[rowPhoneOffset] = 0
	}
	result.Model = updated
	result.Applied = true
	result.NeedsUnportedSuffixDispatch = false
	return result, nil
}
