package text

import (
	"encoding/binary"
	"fmt"
)

// Paul2013TokenBoundaryArenaResult owns the arena after the recovered
// FUN_100130e0 boundary writes and its rebuilt FUN_10012df0 descriptors.
type Paul2013TokenBoundaryArenaResult struct {
	Arena  []byte
	State  Paul2013TokenBoundaryStateResult
	Groups []Paul2013RecordGroupDescriptor
}

// FinishTokenBoundaries reads the nonfinal parser outputs retained by the
// finalized handoff. treePreparedArena must be the result of marker
// initialization followed by the caller's explicit tree stage. The parser
// state offset is the same base supplied to the finalizer, not the source
// rows after its header.
func (state Paul2013FinalizedModelState) FinishTokenBoundaries(treePreparedArena []byte, parserStateOffset int, initialValues, initialStates []int32) (Paul2013TokenBoundaryArenaResult, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(treePreparedArena)
	if err != nil {
		return Paul2013TokenBoundaryArenaResult{}, err
	}
	count := len(inputs.InitialMarkers)
	if state.Finalized.ReturnValue != 1 || count != state.Finalized.OutputGroupCount || count < 1 {
		return Paul2013TokenBoundaryArenaResult{}, fmt.Errorf("boundary count %d does not match successful finalized group count %d", count, state.Finalized.OutputGroupCount)
	}
	if parserStateOffset < 0 || parserStateOffset > len(state.Finalized.StateArena) || count*paul2013ParserSourceRowStride > len(state.Finalized.StateArena)-parserStateOffset {
		return Paul2013TokenBoundaryArenaResult{}, fmt.Errorf("finalized parser arena cannot hold %d rows from +%#x", count, parserStateOffset)
	}
	outputs := make([]int32, count-1)
	for index := range outputs {
		start := parserStateOffset + index*paul2013ParserSourceRowStride + 0x20
		outputs[index] = int32(binary.LittleEndian.Uint32(state.Finalized.StateArena[start : start+4]))
	}
	return FinishPaul2013TokenBoundaryArena(treePreparedArena, initialValues, initialStates, outputs)
}

// InitializePaul2013TokenBoundaryArena writes the native marker table and
// final mode override before FUN_10012df0/FUN_10012f00. Phone grouping and
// marker-tree evaluation belong between this function and Finish below.
func InitializePaul2013TokenBoundaryArena(arena []byte) ([]byte, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(arena)
	if err != nil {
		return nil, err
	}
	if len(inputs.InitialMarkers) == 0 {
		return nil, fmt.Errorf("token boundary pass requires at least one record")
	}
	if len(arena) <= paul2013ModelStateFinalModeOffset {
		return nil, fmt.Errorf("model-state arena lacks mode byte at +%#x", paul2013ModelStateFinalModeOffset)
	}
	result := append([]byte(nil), arena...)
	for index, marker := range inputs.InitialMarkers {
		result[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride+0x3bd] = marker
	}
	terminal := byte('Z')
	if arena[paul2013ModelStateFinalModeOffset] == 6 {
		terminal = '^'
	}
	result[paul2013TokenMarkerArenaRecords+(len(inputs.InitialMarkers)-1)*paul2013TokenMarkerArenaStride+0x3bd] = terminal
	return result, nil
}

// FinishPaul2013TokenBoundaryArena accepts the arena after the marker-tree
// stage, resets +0x3bc flags, applies adjacent-record transitions, updates
// caller-owned value/state arrays, dispatches markers, limits spans, and
// rebuilds descriptors in native order. It does not assume an absent tree
// stage is a no-op. perTokenOutputs are parser-state +0x20 at 0x94 strides.
func FinishPaul2013TokenBoundaryArena(arena []byte, initialValues, initialStates, perTokenOutputs []int32) (Paul2013TokenBoundaryArenaResult, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(arena)
	if err != nil {
		return Paul2013TokenBoundaryArenaResult{}, err
	}
	count := len(inputs.InitialMarkers)
	if count == 0 {
		return Paul2013TokenBoundaryArenaResult{}, fmt.Errorf("token boundary pass requires at least one record")
	}
	markers := make([]Paul2013TokenBoundaryMarker, count)
	raw := make([]byte, count)
	for index := range markers {
		raw[index] = arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride+0x3bd]
		markers[index].Marker = raw[index]
	}
	markers, err = BuildPaul2013TokenBoundaryMarkerTransitions(markers, inputs.RowStateBytes, inputs.TokenCodes)
	if err != nil {
		return Paul2013TokenBoundaryArenaResult{}, err
	}
	for index := range raw {
		raw[index] = markers[index].Marker
	}
	state, err := BuildPaul2013TokenBoundaryStateResult(raw, initialValues, initialStates, perTokenOutputs)
	if err != nil {
		return Paul2013TokenBoundaryArenaResult{}, err
	}
	for index := range state.Markers {
		state.Markers[index].PreviousStateFlag = state.Markers[index].PreviousStateFlag || markers[index].PreviousStateFlag
		raw[index] = state.Markers[index].Marker
	}
	raw, err = applyPaul2013RecordBoundaryDurationLimit(raw, inputs.DurationValues)
	if err != nil {
		return Paul2013TokenBoundaryArenaResult{}, err
	}
	result := Paul2013TokenBoundaryArenaResult{Arena: append([]byte(nil), arena...), State: state}
	records := make([][]byte, count)
	for index, marker := range raw {
		start := paul2013TokenMarkerArenaRecords + index*paul2013TokenMarkerArenaStride
		records[index] = result.Arena[start : start+paul2013TokenMarkerArenaStride]
		records[index][0x3bc] = 0
		if state.Markers[index].PreviousStateFlag {
			records[index][0x3bc] = 1
		}
		records[index][0x3bd] = marker
		result.State.Markers[index].Marker = marker
	}
	result.Groups, err = BuildPaul2013RecordGroupDescriptors(records)
	if err != nil {
		return Paul2013TokenBoundaryArenaResult{}, err
	}
	return result, nil
}

// The native overflow branch backs up one record, inserts '[', resets the
// span, and reprocesses the overflowing record. Record capacity is 100,
// distinct from the 65-phone capacity of a single token.
func applyPaul2013RecordBoundaryDurationLimit(markers, values []byte) ([]byte, error) {
	if len(markers) != len(values) || len(markers) > paul2013RecordGroupCapacity {
		return nil, fmt.Errorf("invalid record boundary lengths %d/%d", len(markers), len(values))
	}
	result := append([]byte(nil), markers...)
	accumulated, spanStart := 0, 0
	for index := 0; index < len(result); index++ {
		accumulated += int(values[index]) * 2
		if accumulated > paul2013PhoneBlockDurationLimit {
			if index == spanStart {
				return nil, fmt.Errorf("record %d cannot split overflowing boundary span", index)
			}
			index--
			result[index] = '['
		}
		switch result[index] {
		case '[', 'Z', '^', '`':
			spanStart, accumulated = index+1, 0
		}
	}
	return result, nil
}
