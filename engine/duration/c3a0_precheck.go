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
	paul2013C3A0ContextCountOffset = 0x429a2
	paul2013C3A0ContextRowBase     = 0x429a8
	paul2013C3A0ContextRowStride   = 0x70
	paul2013C3A0ContextStatus      = 0x04
	paul2013C3A0ContextPhone       = 0x23
	paul2013C3A0ContextState       = 0x6c
)

// Paul2013FUN100086C0Disposition is the outcome of the directly recovered
// early control-flow slice inside FUN_100086c0.
type Paul2013FUN100086C0Disposition uint8

const (
	// Paul2013FUN100086C0ReturnsZero is a native low-short zero return. The
	// caller may continue to its next normalization handler.
	Paul2013FUN100086C0ReturnsZero Paul2013FUN100086C0Disposition = iota
	// Paul2013FUN100086C0ReturnsOne is a native low-short one return. The
	// caller treats this helper as having handled or blocked the current path.
	Paul2013FUN100086C0ReturnsOne
	// Paul2013FUN100086C0NeedsNeighborScan marks the point where the native
	// function enters its unresolved previous/next context-row search.
	Paul2013FUN100086C0NeedsNeighborScan
)

// Paul2013FUN100086C0PrecheckResult retains the branch selected before the
// neighbor scan. NativeShort is set only for paths that return immediately.
type Paul2013FUN100086C0PrecheckResult struct {
	Disposition Paul2013FUN100086C0Disposition
	NativeShort int16
	Reason      string
}

// EvaluatePaul2013FUN100086C0Precheck ports FUN_100086c0 from the label after
// its caller-dependent initial table/character gate through the vowel and
// generic-normalizer short-circuits. alternateRecord is the optional native
// param_4 record and must retain its raw offsets. Paths reaching the native
// previous/next context-row scan return NeedsNeighborScan; that scan and its
// following sorted-table checks are not inferred here.
func (engine *Engine) EvaluatePaul2013FUN100086C0Precheck(
	ctx context.Context,
	source []byte,
	model []byte,
	contextRowIndex int,
	alternateRecord []byte,
) (Paul2013FUN100086C0PrecheckResult, error) {
	if engine == nil {
		return Paul2013FUN100086C0PrecheckResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013FUN100086C0PrecheckResult{}, errors.New("FUN_100086c0 precheck has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(model) < paul2013C3A0ContextCountOffset+2 {
		return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013C3A0ContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013C3A0ContextCountOffset:])))
	if rowCount < 0 || contextRowIndex < 0 || contextRowIndex >= rowCount {
		return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("context row index %d is outside declared row count %d", contextRowIndex, rowCount)
	}
	if rowCount > (len(model)-paul2013C3A0ContextRowBase)/paul2013C3A0ContextRowStride {
		return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("model declares %d context rows that do not fit its %d bytes", rowCount, len(model))
	}
	allFlagged, err := text.Paul2013FUN10010010(source)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("run FUN_10010010 precheck: %w", err)
	}
	if !allFlagged {
		return paul2013C3A0PrecheckReturn(0, "FUN_10010010 returned zero"), nil
	}

	rowStart := paul2013C3A0ContextRowBase + contextRowIndex*paul2013C3A0ContextRowStride
	var row []byte
	if len(alternateRecord) != 0 {
		status := alternateRecord[0]
		if status == 'E' {
			return paul2013C3A0PrecheckReturn(1, "alternate record status is E"), nil
		}
		if status == 'A' {
			return Paul2013FUN100086C0PrecheckResult{
				Disposition: Paul2013FUN100086C0NeedsNeighborScan,
				Reason:      "alternate record status is A",
			}, nil
		}
		if len(alternateRecord) < 0x10 {
			return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("alternate record has %d bytes, need fields through +0x0f", len(alternateRecord))
		}
		if binary.LittleEndian.Uint32(alternateRecord[0x0c:]) != 0 ||
			binary.LittleEndian.Uint16(alternateRecord[0x08:]) != 0 {
			return paul2013C3A0PrecheckReturn(0, "non-A alternate record carries native state"), nil
		}
	} else {
		row = model[rowStart : rowStart+paul2013C3A0ContextRowStride]
		status := row[paul2013C3A0ContextStatus]
		if status == 'E' {
			return paul2013C3A0PrecheckReturn(1, "context row status is E"), nil
		}
		if status == 'A' {
			return Paul2013FUN100086C0PrecheckResult{
				Disposition: Paul2013FUN100086C0NeedsNeighborScan,
				Reason:      "context row status is A",
			}, nil
		}
		if row[paul2013C3A0ContextPhone] != 0 {
			return paul2013C3A0PrecheckReturn(0, "context row already has a phone string"), nil
		}
		if binary.LittleEndian.Uint16(row[paul2013C3A0ContextState:]) != 0 {
			return paul2013C3A0PrecheckReturn(0, "context row state at +0x6c is nonzero"), nil
		}
	}

	vowelCount, err := text.Paul2013FUN100100D0(source)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("run FUN_100100d0 vowel count: %w", err)
	}
	if vowelCount == 0 {
		return paul2013C3A0PrecheckReturn(1, "FUN_100100d0 found no mapped vowels"), nil
	}
	if len(source) > 2 {
		eligible, err := engine.IsPaul2013NormalizerEligible(source)
		if err != nil {
			return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("run FUN_10002c70 eligibility check: %w", err)
		}
		if eligible {
			return paul2013C3A0PrecheckReturn(0, "FUN_10002c70 accepted the source"), nil
		}
	}
	return Paul2013FUN100086C0PrecheckResult{
		Disposition: Paul2013FUN100086C0NeedsNeighborScan,
		Reason:      "source passed local prechecks and requires the native neighbor scan",
	}, nil
}

func paul2013C3A0PrecheckReturn(nativeShort int16, reason string) Paul2013FUN100086C0PrecheckResult {
	disposition := Paul2013FUN100086C0ReturnsZero
	if nativeShort != 0 {
		disposition = Paul2013FUN100086C0ReturnsOne
	}
	return Paul2013FUN100086C0PrecheckResult{
		Disposition: disposition,
		NativeShort: nativeShort,
		Reason:      reason,
	}
}
