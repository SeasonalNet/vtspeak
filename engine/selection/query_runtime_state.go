package selection

import "fmt"

const (
	paul2013QueryRuntimeRecordStride = 0x3c0
	paul2013QueryPriorMetricOffset   = 0x8a9
	paul2013QueryContextGateOffset   = 0x8ea
)

// Paul2013QueryRuntimeState contains the per-position values read by
// FUN_10018770 through the pointer stored at its state argument +0x4c.
type Paul2013QueryRuntimeState struct {
	PriorMetric uint32
	ContextGate byte
}

// ReadPaul2013QueryRuntimeState projects FUN_10018770's query metric and live
// gate from the byte arena addressed by the native pointer at state +0x4c.
// The six-byte state row selects a 0x3c0-byte slot with byte 0 and a phone
// offset with byte 2. This function accepts the pointer target as a bounded
// slice because the native process pointer itself is not portable.
func ReadPaul2013QueryRuntimeState(
	stateArena []byte,
	state Paul2013TransitionContextState,
) (Paul2013QueryRuntimeState, error) {
	metricOffset, gateOffset := paul2013QueryRuntimeStateOffsets(state)
	if metricOffset >= len(stateArena) {
		return Paul2013QueryRuntimeState{}, fmt.Errorf("query prior metric at +%#x exceeds state arena size %d", metricOffset, len(stateArena))
	}
	if gateOffset >= len(stateArena) {
		return Paul2013QueryRuntimeState{}, fmt.Errorf("query context gate at +%#x exceeds state arena size %d", gateOffset, len(stateArena))
	}
	return Paul2013QueryRuntimeState{
		PriorMetric: uint32(stateArena[metricOffset]),
		ContextGate: stateArena[gateOffset],
	}, nil
}

func paul2013QueryRuntimeStateOffsets(state Paul2013TransitionContextState) (int, int) {
	blockOffset := int(state[0]) * paul2013QueryRuntimeRecordStride
	phoneOffset := int(state[2])
	return blockOffset + phoneOffset + paul2013QueryPriorMetricOffset,
		blockOffset + phoneOffset + paul2013QueryContextGateOffset
}

func readPaul2013QueryPriorMetric(stateArena []byte, state Paul2013TransitionContextState) uint32 {
	metricOffset, _ := paul2013QueryRuntimeStateOffsets(state)
	return uint32(stateArena[metricOffset])
}
