package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013C3A0ComponentLimit = 31
	paul2013C3A0OutputLimit    = 63
)

// Paul2013C3A0ComponentInput is one source component emitted by
// FUN_1000c860 and passed to FUN_1000c3a0's handler chain.
type Paul2013C3A0ComponentInput struct {
	Index      int
	SourceSpan text.Paul2013C3A0Component
}

// Paul2013C3A0ComponentResult contains one component handler's native output,
// row-flag changes, and local handled-state update.
type Paul2013C3A0ComponentResult struct {
	Output      []byte
	RowFlagMask byte
	Handled     bool
}

// Paul2013C3A0ComponentHandler supplies the per-component dictionary,
// exception, TPP, suffix, and generic fallback chain. The outer scanner and
// result assembly are ported independently of unresolved handler branches.
type Paul2013C3A0ComponentHandler func(
	context.Context,
	Paul2013C3A0ComponentInput,
) (Paul2013C3A0ComponentResult, error)

// Paul2013C3A0SequenceResult reports the supported outer portion of
// FUN_1000c3a0. Output excludes its terminating C NUL.
type Paul2013C3A0SequenceResult struct {
	Processed      bool
	Matched        bool
	Source         []byte
	Components     []Paul2013C3A0ComponentInput
	Output         []byte
	RowFlags       byte
	OutputOverflow bool
}

// NormalizePaul2013C3A0ComponentSequence ports FUN_1000c3a0's repeated
// FUN_1000c860 scan, one-component 32-byte-buffer bound, `d` output joining,
// 63-byte length check, and final bit-3 return gate. enabled is the explicit
// result of FUN_10008cc0; component-specific normalization remains a callback.
func NormalizePaul2013C3A0ComponentSequence(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	enabled bool,
	handle Paul2013C3A0ComponentHandler,
) (Paul2013C3A0SequenceResult, error) {
	result := Paul2013C3A0SequenceResult{Source: append([]byte(nil), source...), RowFlags: currentFlags}
	if ctx == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("FUN_1000c3a0 component sequence has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	for offset, value := range source {
		if value >= 0x80 {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("FUN_1000c3a0 source byte 0x%02x at offset %d is outside the supported ASCII input", value, offset)
		}
	}
	if !enabled || len(source) == 0 {
		return result, nil
	}
	if handle == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("FUN_1000c3a0 component sequence has no component handler")
	}
	result.Processed = true
	sourceLength := len(source)
	offset := 0
	handledAny := false
	for componentIndex := 0; offset < sourceLength; componentIndex++ {
		if err := ctx.Err(); err != nil {
			return Paul2013C3A0SequenceResult{}, err
		}
		component, err := text.ScanPaul2013C3A0Component(source, offset, sourceLength)
		if err != nil {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("scan FUN_1000c3a0 component %d: %w", componentIndex, err)
		}
		if componentIndex == 0 && component.EndOffset == sourceLength {
			// Native FUN_1000c3a0 bypasses a source that scans as one component.
			return result, nil
		}
		if component.EndOffset <= offset {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("FUN_1000c860 made no source progress at byte %d", offset)
		}
		if len(component.Output) > paul2013C3A0ComponentLimit {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("FUN_1000c3a0 component %d has %d bytes; its native local buffer holds %d plus NUL", componentIndex, len(component.Output), paul2013C3A0ComponentLimit)
		}
		input := Paul2013C3A0ComponentInput{Index: componentIndex, SourceSpan: component}
		result.Components = append(result.Components, input)
		componentResult, err := handle(ctx, input)
		if err != nil {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("handle FUN_1000c3a0 component %d %q: %w", componentIndex, component.Output, err)
		}
		if componentResult.RowFlagMask & ^byte(0x39) != 0 {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("FUN_1000c3a0 component %d returned unsupported row-flag mask %#02x", componentIndex, componentResult.RowFlagMask)
		}
		result.RowFlags |= componentResult.RowFlagMask
		handledAny = handledAny || componentResult.Handled
		if bytes.IndexByte(componentResult.Output, 0) >= 0 {
			return Paul2013C3A0SequenceResult{}, fmt.Errorf("FUN_1000c3a0 component %d output contains NUL", componentIndex)
		}
		if len(result.Output)+len(componentResult.Output)+1 > paul2013C3A0OutputLimit {
			result.OutputOverflow = true
			break
		}
		if len(result.Output) > 0 {
			result.Output = append(result.Output, 'd')
		}
		result.Output = append(result.Output, componentResult.Output...)
		offset = component.EndOffset
	}
	if len(result.Output) > 0 && handledAny {
		result.RowFlags |= 0x08
		result.Matched = true
	}
	return result, nil
}

// NormalizePaul2013C3A0ComponentSequenceFromSource derives the native
// FUN_10008cc0 gate and then runs the supported sequence shell.
func NormalizePaul2013C3A0ComponentSequenceFromSource(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	handle Paul2013C3A0ComponentHandler,
) (Paul2013C3A0SequenceResult, error) {
	if ctx == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("FUN_1000c3a0 component sequence has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013C3A0SequenceResult{}, err
	}
	enabled, err := text.Paul2013C3A0Gate(source)
	if err != nil {
		return Paul2013C3A0SequenceResult{}, err
	}
	return NormalizePaul2013C3A0ComponentSequence(ctx, source, currentFlags, enabled, handle)
}
