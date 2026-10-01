package duration

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

const paul2013ModelContextTailOutputCapacity = 68

// Paul2013ModelContextTailResult records the ordered Y/S tail path in
// FUN_100091b0. Path is one of literal, embedded-dictionary,
// embedded-generic, direct-generic, or context-encoder; an empty path means
// the native tail produced no output.
type Paul2013ModelContextTailResult struct {
	Model   []byte
	Applied bool
	Path    string
}

// ApplyPaul2013ModelContextTail completes the tail work following the direct
// contraction literals in FUN_100091b0. On a literal miss it removes the
// leading tail byte, then tries FUN_1000cb30, FUN_10002f10, and FUN_10009cd0
// in native order. The engine's existing resource and ASCII boundaries still
// apply to the embedded lookup and normalizers.
func (engine *Engine) ApplyPaul2013ModelContextTail(
	ctx context.Context,
	model []byte,
	rowIndex int,
	tailSurface []byte,
) (Paul2013ModelContextTailResult, error) {
	if engine == nil {
		return Paul2013ModelContextTailResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013ModelContextTailResult{}, errors.New("model-context tail has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ModelContextTailResult{}, err
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model-context tail row %d is outside row count %d", rowIndex, rowCount)
	}
	if nul := bytes.IndexByte(tailSurface, 0); nul >= 0 {
		tailSurface = tailSurface[:nul]
	}
	result := Paul2013ModelContextTailResult{Model: append([]byte(nil), model...)}
	if len(tailSurface) == 0 {
		return result, nil
	}
	if len(tailSurface) >= 32 {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model-context tail has %d bytes, exceeding FUN_100091b0's 32-byte local buffer", len(tailSurface))
	}

	rowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	flagOffset := rowStart - 2
	rowEnd := rowStart + paul2013ModelContextRowStride
	if flagOffset < 0 || rowEnd > len(result.Model) {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model-context tail row %d exceeds model arena", rowIndex)
	}
	_, handled, err := text.ApplyPaul2013LiteralContextHandler(
		result.Model[flagOffset:rowEnd], tailSurface,
	)
	if err != nil {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("apply literal model-context tail at row %d: %w", rowIndex, err)
	}
	if handled {
		result.Applied = true
		result.Path = "literal"
		return result, nil
	}

	// FUN_100091b0 copies local_2c + 1 into local_4c before the shared
	// dictionary/generic/context fallback chain.
	source := append([]byte(nil), tailSurface[1:]...)
	destination := make([]byte, paul2013ModelContextTailOutputCapacity)
	lookupWrites, err := engine.NormalizePaul2013ContextTokenToBuffer(
		ctx, source, engine.LookupPaul2013EmbeddedContextToken, destination,
	)
	if err != nil {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("run FUN_1000cb30 for model-context tail row %d: %w", rowIndex, err)
	}
	var output []byte
	stateFlags := byte(0)
	produced := false
	if lookupWrites.Result.NativeReturnCode == 1 {
		output = cStringBytes(lookupWrites.Output)
		produced = true
		result.Path = "embedded-dictionary"
		if !lookupWrites.Result.UsedDictionary {
			result.Path = "embedded-generic"
		}
	} else {
		generic, err := engine.NormalizePaul2013GenericToken(ctx, source)
		if err != nil {
			return Paul2013ModelContextTailResult{}, fmt.Errorf("run direct FUN_10002f10 for model-context tail row %d: %w", rowIndex, err)
		}
		if generic.Eligible {
			genericWrites, err := ApplyPaul2013GenericNormalizerWrites(
				destination, result.Model[flagOffset], generic,
			)
			if err != nil {
				return Paul2013ModelContextTailResult{}, fmt.Errorf("write direct generic tail at row %d: %w", rowIndex, err)
			}
			output = cStringBytes(genericWrites.Output)
			stateFlags = genericWrites.Flags
			produced = true
			result.Path = "direct-generic"
		} else {
			output, err = text.EncodePaul2013ContextString(source)
			if err != nil {
				return Paul2013ModelContextTailResult{}, fmt.Errorf("encode FUN_10009cd0 model-context tail at row %d: %w", rowIndex, err)
			}
			if len(output) != 0 {
				stateFlags |= 0x20
				produced = true
			}
			result.Path = "context-encoder"
		}
	}
	if !produced {
		return result, nil
	}
	if bytes.IndexByte(output, 0) >= 0 {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model-context tail row %d output contains an embedded NUL", rowIndex)
	}
	phoneStart := rowStart + 0x23
	phoneArea := result.Model[phoneStart : rowStart+0x64]
	terminator := bytes.IndexByte(phoneArea, 0)
	if terminator < 0 {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model-context tail row %d phone codes are not NUL-terminated", rowIndex)
	}
	appendLength := 1 + len(output) // Native writes 'd' before local_90.
	if terminator+appendLength >= len(phoneArea) {
		return Paul2013ModelContextTailResult{}, fmt.Errorf("model-context tail row %d needs %d bytes plus NUL, phone field has %d", rowIndex, appendLength, len(phoneArea)-terminator)
	}
	phoneArea[terminator] = 'd'
	copy(phoneArea[terminator+1:], output)
	phoneArea[terminator+appendLength] = 0
	result.Model[flagOffset] |= stateFlags | 0x08
	result.Applied = true
	return result, nil
}
