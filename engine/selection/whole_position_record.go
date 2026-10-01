package selection

import (
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013QueryRecordStride      = 0x3c0
	paul2013QueryRecordBlockStride = 3 * paul2013QueryRecordStride
	paul2013QueryRowClassOffset    = 0x6de
	paul2013QuerySignatureOffset   = 0x6e2
	paul2013QueryModelClassOffset  = 0x92b
)

// Paul2013WholePositionRecordInput contains the fields FUN_10024060 reads
// from the three-record workspace block selected by state[0]. State[2]
// indexes seven-byte signatures within that block. The row class and model
// class are read from their distinct observed offsets; they are not assigned
// broader semantic labels here.
type Paul2013WholePositionRecordInput struct {
	Target     text.Context
	ContextRow [6]byte
	ModelClass byte
}

// BuildPaul2013WholePositionRecordInput projects the target signature,
// row-class byte, and model-class byte consumed by FUN_10024060 from a
// contiguous record arena and its already-produced six-byte workspace
// row. It preserves row bytes 0..2 and 5, writes the observed row class to
// byte 3, and forces byte 4 to zero for a whole-position query.
func BuildPaul2013WholePositionRecordInput(
	recordArena []byte,
	state [6]byte,
) (Paul2013WholePositionRecordInput, error) {
	blockOrdinal := int(state[0])
	if blockOrdinal > (int(^uint(0)>>1)-paul2013QueryModelClassOffset)/paul2013QueryRecordBlockStride {
		return Paul2013WholePositionRecordInput{}, fmt.Errorf("query record block ordinal %d overflows arena indexing", blockOrdinal)
	}
	blockOffset := blockOrdinal * paul2013QueryRecordBlockStride
	signatureOffset := blockOffset + paul2013QuerySignatureOffset + int(state[2])*7
	modelClassOffset := blockOffset + paul2013QueryModelClassOffset
	if signatureOffset > len(recordArena) || len(recordArena)-signatureOffset < 7 {
		return Paul2013WholePositionRecordInput{}, fmt.Errorf("query signature at +%#x needs 7 bytes, arena has %d", signatureOffset, len(recordArena))
	}
	if modelClassOffset >= len(recordArena) {
		return Paul2013WholePositionRecordInput{}, fmt.Errorf("query model class at +%#x exceeds arena size %d", modelClassOffset, len(recordArena))
	}
	rowClassOffset := blockOffset + paul2013QueryRowClassOffset
	if rowClassOffset >= len(recordArena) {
		return Paul2013WholePositionRecordInput{}, fmt.Errorf("query row class at +%#x exceeds arena size %d", rowClassOffset, len(recordArena))
	}
	input := Paul2013WholePositionRecordInput{
		ContextRow: state,
		ModelClass: recordArena[modelClassOffset],
	}
	copy(input.Target.Signature[:], recordArena[signatureOffset:signatureOffset+7])
	input.ContextRow[3] = recordArena[rowClassOffset]
	input.ContextRow[4] = 0
	return input, nil
}

// QueryPaul2013WholePositionFromRecordState connects the native record-arena
// projection above to FUN_10024060's whole-position query. The six-byte
// workspace row, prior metric, and context-gate producer remain caller inputs;
// conditional fallback-row generation is handled by QueryPaul2013Position.
func QueryPaul2013WholePositionFromRecordState(
	catalog Paul2013PrefixClassLookup,
	recordArena []byte,
	state [6]byte,
	priorMetric uint32,
	contextGate func() byte,
) (Paul2013WholePositionQueryResult, error) {
	input, err := BuildPaul2013WholePositionRecordInput(recordArena, state)
	if err != nil {
		return Paul2013WholePositionQueryResult{}, fmt.Errorf("read whole-position query state: %w", err)
	}
	return QueryPaul2013WholePosition(
		catalog, input.Target, input.ContextRow, priorMetric, input.ModelClass, contextGate,
	)
}

// QueryPaul2013WholePositionFromTransitionState joins the compact state row
// written by FUN_10024680 to FUN_10024060's record-backed query. Both native
// functions address the same six-byte workspace row at +0xec628: byte 0 is
// the record-block ordinal and byte 2 is the phone/signature ordinal. The
// prior metric and context gate are still produced outside this projection.
func QueryPaul2013WholePositionFromTransitionState(
	catalog Paul2013PrefixClassLookup,
	recordArena []byte,
	state Paul2013TransitionContextState,
	priorMetric uint32,
	contextGate func() byte,
) (Paul2013WholePositionQueryResult, error) {
	return QueryPaul2013WholePositionFromRecordState(
		catalog, recordArena, [6]byte(state), priorMetric, contextGate,
	)
}

// QueryPaul2013PositionFromTransitionState composes the record-backed
// whole-position query with FUN_100242a0's conditional two-row fallback.
// The fallback base row, model class, and signed row-class delta come from the
// same record projection. Per-row fallback runtime values remain explicit.
func QueryPaul2013PositionFromTransitionState(
	catalog Paul2013PrefixClassLookup,
	recordArena []byte,
	state Paul2013TransitionContextState,
	priorMetric uint32,
	contextGate func() byte,
	fallbackRuntime [2]Paul2013FallbackRowRuntime,
) (Paul2013PositionQueryResult, error) {
	input, err := BuildPaul2013WholePositionRecordInput(recordArena, [6]byte(state))
	if err != nil {
		return Paul2013PositionQueryResult{}, fmt.Errorf("read transition-state position query: %w", err)
	}
	return QueryPaul2013Position(
		catalog,
		input.Target,
		input.ContextRow,
		priorMetric,
		input.ModelClass,
		contextGate,
		input.ContextRow,
		int8(input.ContextRow[3]),
		input.ModelClass,
		fallbackRuntime,
	)
}

// QueryPaul2013PositionFromStateArenas reads the record-backed query row and
// the per-position metric/gate values addressed through FUN_10018770's
// state-pointer field, then executes whole-position lookup and conditional
// fallback. The query-state arena is the target of the native pointer at
// state +0x4c; the returned prior metric is a byte value widened to uint32.
func QueryPaul2013PositionFromStateArenas(
	catalog Paul2013PrefixClassLookup,
	recordArena []byte,
	queryStateArena []byte,
	state Paul2013TransitionContextState,
) (Paul2013PositionQueryResult, error) {
	input, err := BuildPaul2013WholePositionRecordInput(recordArena, [6]byte(state))
	if err != nil {
		return Paul2013PositionQueryResult{}, fmt.Errorf("read transition-state position query: %w", err)
	}
	runtimeState, err := ReadPaul2013QueryRuntimeState(queryStateArena, state)
	if err != nil {
		return Paul2013PositionQueryResult{}, fmt.Errorf("read transition-state query runtime: %w", err)
	}
	_, gateOffset := paul2013QueryRuntimeStateOffsets(state)
	contextGate := func() byte { return queryStateArena[gateOffset] }
	fallbackRuntime := [2]Paul2013FallbackRowRuntime{
		{
			PriorMetric:     runtimeState.PriorMetric,
			ReadPriorMetric: func() uint32 { return readPaul2013QueryPriorMetric(queryStateArena, state) },
			ContextGate:     contextGate,
		},
		{
			PriorMetric:     runtimeState.PriorMetric,
			ReadPriorMetric: func() uint32 { return readPaul2013QueryPriorMetric(queryStateArena, state) },
			ContextGate:     contextGate,
		},
	}
	return QueryPaul2013Position(
		catalog,
		input.Target,
		input.ContextRow,
		runtimeState.PriorMetric,
		input.ModelClass,
		contextGate,
		input.ContextRow,
		int8(input.ContextRow[3]),
		input.ModelClass,
		fallbackRuntime,
	)
}
