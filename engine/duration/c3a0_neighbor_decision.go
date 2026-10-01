package duration

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// EvaluatePaul2013FUN100086C0SupportedPath composes the recovered caller gate,
// precheck, previous/following row scans, compound-character gates, and
// post-scan decisions. callerGate is the native param_5 nonzero condition.
// The local model arena remains an explicit input.
func (engine *Engine) EvaluatePaul2013FUN100086C0SupportedPath(
	ctx context.Context,
	source []byte,
	model []byte,
	contextRowIndex int,
	alternateRecord []byte,
	callerGate bool,
) (Paul2013FUN100086C0PrecheckResult, Paul2013FUN100086C0Neighbors, error) {
	if callerGate {
		initial, shortCircuit, err := evaluatePaul2013FUN100086C0CallerGate(
			ctx, source, model, contextRowIndex, alternateRecord,
		)
		if err != nil {
			return Paul2013FUN100086C0PrecheckResult{}, Paul2013FUN100086C0Neighbors{}, err
		}
		if shortCircuit {
			return initial, Paul2013FUN100086C0Neighbors{}, nil
		}
	}
	precheck, err := engine.EvaluatePaul2013FUN100086C0Precheck(ctx, source, model, contextRowIndex, alternateRecord)
	if err != nil || precheck.Disposition != Paul2013FUN100086C0NeedsNeighborScan {
		return precheck, Paul2013FUN100086C0Neighbors{}, err
	}
	neighbors, err := ScanPaul2013FUN100086C0Neighbors(ctx, model, contextRowIndex)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, Paul2013FUN100086C0Neighbors{}, err
	}
	result, err := CompletePaul2013FUN100086C0NeighborDecision(
		ctx, engine, precheck, source, model, contextRowIndex, alternateRecord,
		neighbors,
	)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, neighbors, err
	}
	return result, neighbors, nil
}

func evaluatePaul2013FUN100086C0CallerGate(
	ctx context.Context,
	source []byte,
	model []byte,
	contextRowIndex int,
	alternateRecord []byte,
) (Paul2013FUN100086C0PrecheckResult, bool, error) {
	if ctx == nil {
		return Paul2013FUN100086C0PrecheckResult{}, false, errors.New("FUN_100086c0 caller gate has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, false, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if _, err := paul2013C3A0ValidateContextRows(model, contextRowIndex); err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, false, err
	}
	status, err := paul2013C3A0NeighborDecisionStatus(model, contextRowIndex, alternateRecord)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, false, err
	}
	if len(source) == 1 && status != 'E' {
		return Paul2013FUN100086C0PrecheckResult{}, false, nil
	}
	if text.Paul2013MappedCStringEqual(
		source,
		[]byte("Wi"),
		text.Paul2013EmbeddedKeyTables().CharacterMap,
	) {
		return paul2013C3A0PrecheckReturn(1, "caller-gated source matched mapped Wi"), true, nil
	}
	allowed, err := Paul2013CompoundCharacterGate(source)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, false, fmt.Errorf("run caller-gated FUN_1000ffd0: %w", err)
	}
	if allowed {
		return paul2013C3A0PrecheckReturn(1, "caller-gated FUN_1000ffd0 accepted the source"), true, nil
	}
	return Paul2013FUN100086C0PrecheckResult{}, false, nil
}

// CompletePaul2013FUN100086C0NeighborDecision ports the status-specific
// decisions after ScanPaul2013FUN100086C0Neighbors, including the static
// sorted-source table lookup. alternateRecord must retain the native byte
// field at offset +0x10.
func CompletePaul2013FUN100086C0NeighborDecision(
	ctx context.Context,
	engine *Engine,
	precheck Paul2013FUN100086C0PrecheckResult,
	source []byte,
	model []byte,
	contextRowIndex int,
	alternateRecord []byte,
	neighbors Paul2013FUN100086C0Neighbors,
) (Paul2013FUN100086C0PrecheckResult, error) {
	if precheck.Disposition != Paul2013FUN100086C0NeedsNeighborScan {
		return precheck, nil
	}
	if ctx == nil {
		return Paul2013FUN100086C0PrecheckResult{}, errors.New("FUN_100086c0 neighbor decision has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	rowCount, err := paul2013C3A0ValidateContextRows(model, contextRowIndex)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, err
	}
	status, err := paul2013C3A0NeighborDecisionStatus(model, contextRowIndex, alternateRecord)
	if err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, err
	}
	if err := paul2013C3A0ValidateNeighborIndexes(neighbors, contextRowIndex, rowCount); err != nil {
		return Paul2013FUN100086C0PrecheckResult{}, err
	}

	connected := (neighbors.Previous == nil || neighbors.Previous.CharacterGate) &&
		(neighbors.Following == nil || neighbors.Following.CharacterGate) &&
		(neighbors.Previous != nil || neighbors.Following != nil)
	if !connected {
		if status == 'A' && neighbors.Previous != nil &&
			text.Paul2013MappedCStringEqual(
				neighbors.Previous.Surface,
				[]byte("Lexus"),
				text.Paul2013EmbeddedKeyTables().CharacterMap,
			) && bytes.Equal(source, []byte("IS")) {
			return paul2013C3A0PrecheckReturn(1, "mapped Lexus previous neighbor and exact IS source matched the native special case"), nil
		}
		_, found, lookupErr := paul2013C3A0SortedSourceTable.Lookup(source)
		if lookupErr != nil {
			return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("run FUN_100560a0 sorted-source lookup: %w", lookupErr)
		}
		if !found {
			return paul2013C3A0PrecheckReturn(1, "FUN_100560a0 did not find the source in the recovered sorted table"), nil
		}
	}

	if status == 'A' {
		if len(alternateRecord) != 0 {
			if len(alternateRecord) < 0x11 {
				return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("alternate record has %d bytes, need phone byte at +0x10", len(alternateRecord))
			}
			if alternateRecord[0x10] == 0 && binary.LittleEndian.Uint16(alternateRecord[0x08:]) == 0 {
				return paul2013C3A0PrecheckReturn(0, "alternate A record has no phone and its +0x08 short is zero"), nil
			}
		} else {
			rowStart := paul2013C3A0ContextRowBase + contextRowIndex*paul2013C3A0ContextRowStride
			row := model[rowStart : rowStart+paul2013C3A0ContextRowStride]
			rowPrefix := model[paul2013C3A0ContextCountOffset+contextRowIndex*paul2013C3A0ContextRowStride:]
			if row[paul2013C3A0ContextPhone] == 0 && binary.LittleEndian.Uint16(rowPrefix[0x72:]) == 0 {
				return paul2013C3A0PrecheckReturn(0, "A context row has no phone and its +0x72 short is zero"), nil
			}
		}
		if len(source) < 3 {
			return paul2013C3A0PrecheckReturn(1, "A-row branch source length is below three"), nil
		}
		eligible, eligibilityErr := engine.IsPaul2013NormalizerEligible(source)
		if eligibilityErr != nil {
			return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("run final FUN_10002c70 eligibility check: %w", eligibilityErr)
		}
		if eligible {
			return paul2013C3A0PrecheckReturn(0, "final FUN_10002c70 eligibility check accepted the source"), nil
		}
		return paul2013C3A0PrecheckReturn(1, "A-row branch reached native low-short one"), nil
	}

	if len(alternateRecord) != 0 {
		if len(alternateRecord) < 8 {
			return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("alternate record has %d bytes, need fields through +0x07", len(alternateRecord))
		}
		if binary.LittleEndian.Uint16(alternateRecord[0x02:]) != 0 ||
			binary.LittleEndian.Uint16(alternateRecord[0x04:]) != 0 ||
			binary.LittleEndian.Uint16(alternateRecord[0x06:]) != 0 {
			return paul2013C3A0PrecheckReturn(0, "alternate-record post-scan state is nonzero"), nil
		}
	} else {
		rowPrefix := model[paul2013C3A0ContextCountOffset+contextRowIndex*paul2013C3A0ContextRowStride:]
		nextRowPrefix := model[paul2013C3A0ContextCountOffset+(contextRowIndex+1)*paul2013C3A0ContextRowStride:]
		if binary.LittleEndian.Uint16(rowPrefix[0x6c:]) != 0 ||
			binary.LittleEndian.Uint16(nextRowPrefix[:2]) != 0 {
			return paul2013C3A0PrecheckReturn(0, "context post-scan row state is nonzero"), nil
		}
		if binary.LittleEndian.Uint16(rowPrefix[0x6e:]) != 0 {
			return paul2013C3A0PrecheckReturn(0, "context post-scan +0x6e short is nonzero"), nil
		}
	}
	if len(source) < 3 {
		return paul2013C3A0PrecheckReturn(1, "non-A row branch source length is below three"), nil
	}
	eligible, eligibilityErr := engine.IsPaul2013NormalizerEligible(source)
	if eligibilityErr != nil {
		return Paul2013FUN100086C0PrecheckResult{}, fmt.Errorf("run final FUN_10002c70 eligibility check: %w", eligibilityErr)
	}
	if eligible {
		return paul2013C3A0PrecheckReturn(0, "final FUN_10002c70 eligibility check accepted the source"), nil
	}
	return paul2013C3A0PrecheckReturn(1, "native post-neighbor path returned low-short one"), nil
}

// The pointer array at DAT_100783d0 contains six direct string pointers, and
// DAT_100783e8 contains the signed count 6. Its mode-0x53 keys were read from
// the corresponding DLL .data strings in bytewise strcmp order.
var paul2013C3A0SortedSourceTable = &Paul2013NativeCStringTable{entries: [][]byte{
	[]byte("ALL"), []byte("BY"), []byte("DEAR"), []byte("DO"), []byte("IS"), []byte("NOT"),
}}

func paul2013C3A0ValidateContextRows(model []byte, rowIndex int) (int, error) {
	if len(model) < paul2013C3A0ContextCountOffset+2 {
		return 0, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013C3A0ContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013C3A0ContextCountOffset:])))
	if rowCount < 0 || rowIndex < 0 || rowIndex >= rowCount {
		return 0, fmt.Errorf("context row index %d is outside declared row count %d", rowIndex, rowCount)
	}
	if rowCount > (len(model)-paul2013C3A0ContextRowBase)/paul2013C3A0ContextRowStride {
		return 0, fmt.Errorf("model declares %d context rows that do not fit its %d bytes", rowCount, len(model))
	}
	return rowCount, nil
}

func paul2013C3A0NeighborDecisionStatus(model []byte, rowIndex int, alternateRecord []byte) (byte, error) {
	if len(alternateRecord) != 0 {
		return alternateRecord[0], nil
	}
	rowStart := paul2013C3A0ContextRowBase + rowIndex*paul2013C3A0ContextRowStride
	if rowStart+paul2013C3A0ContextStatus >= len(model) {
		return 0, errors.New("context row status is outside the model")
	}
	return model[rowStart+paul2013C3A0ContextStatus], nil
}

func paul2013C3A0ValidateNeighborIndexes(neighbors Paul2013FUN100086C0Neighbors, rowIndex, rowCount int) error {
	if neighbors.Previous != nil && (neighbors.Previous.RowIndex < 0 || neighbors.Previous.RowIndex >= rowIndex) {
		return fmt.Errorf("previous neighbor row %d is not before source row %d", neighbors.Previous.RowIndex, rowIndex)
	}
	if neighbors.Following != nil && (neighbors.Following.RowIndex <= rowIndex || neighbors.Following.RowIndex >= rowCount) {
		return fmt.Errorf("following neighbor row %d is outside the rows after source row %d of %d", neighbors.Following.RowIndex, rowIndex, rowCount)
	}
	return nil
}
