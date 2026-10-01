package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// Paul2013PositionBoundaryPipelineResult retains state production, updated
// model records, and the downstream native boundary/workspace result.
type Paul2013PositionBoundaryPipelineResult struct {
	State    Paul2013FinalizedModelState
	Position Paul2013PositionStateWorkspaceResult
	Boundary Paul2013TokenBoundaryWorkspaceResult
}

// ApplyPositionStatesAndRunTokenBoundaries connects FUN_10022970's recovered
// state program to optional FUN_10022dc0 interval mapping, FUN_10016c90's
// slot-7 flags/final mode, and the supported FUN_100130e0 boundary pipeline.
// Phone fields and pointers already come from the finalized record handoff;
// only fields affected by the newly produced position state are rewritten.
func (state Paul2013FinalizedModelState) ApplyPositionStatesAndRunTokenBoundaries(workspace []byte, program Paul2013PositionStateProgram, arenaAddress uint32, parserStateOffset int, catalog *tree3.Catalog, resolvePointer Paul2013CStringPointerResolver) (Paul2013PositionBoundaryPipelineResult, error) {
	result := Paul2013PositionBoundaryPipelineResult{State: state}
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(state.Records.Bytes)
	if err != nil {
		return result, err
	}
	count := len(inputs.InitialMarkers)
	if count < 1 || state.Finalized.ReturnValue != 1 || state.Finalized.OutputGroupCount != count {
		return result, fmt.Errorf("position/boundary record count does not match successful finalizer")
	}
	if parserStateOffset < 0 || parserStateOffset > len(state.Finalized.StateArena) || count*paul2013ParserSourceRowStride > len(state.Finalized.StateArena)-parserStateOffset {
		return result, fmt.Errorf("position/boundary parser rows outside finalized arena")
	}
	parser := state.Finalized.StateArena[parserStateOffset:]
	segmentBytes := int32(binary.LittleEndian.Uint32(parser[4:8]))
	result.Position, err = ApplyPaul2013PositionStateProgramFromModelArena(workspace, state.Records.Bytes, segmentBytes, program)
	if err != nil {
		return result, fmt.Errorf("produce model position states: %w", err)
	}
	values := result.Position.Program.Arrays[7]
	lastFlag := byte(0)
	if values[count-1] != -1 {
		lastFlag = 12
	}
	summary, err := SummarizePaul2013ContextCodeState(values, inputs.TokenCodes[count-1], lastFlag, int16(binary.LittleEndian.Uint16(parser[2:4])))
	if err != nil {
		return result, err
	}
	updated := append([]byte(nil), state.Records.Bytes...)
	if len(updated) <= paul2013ModelStateFinalModeOffset {
		return result, fmt.Errorf("model position/boundary arena lacks final mode")
	}
	for index, flag := range summary.StateFlags {
		row := updated[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		row[0x2df] = flag
		if result.Position.Program.HasMappedIntervals {
			interval := result.Position.Program.MappedIntervals[index]
			binary.LittleEndian.PutUint32(row[:4], uint32(interval.Minimum))
			binary.LittleEndian.PutUint32(row[4:8], uint32(interval.Maximum))
		}
	}
	updated[paul2013ModelStateFinalModeOffset] = summary.ModeCode
	result.State.Records.Bytes = updated
	result.Boundary, err = result.State.RunTokenBoundariesWithWorkspace(arenaAddress, parserStateOffset, catalog, resolvePointer, result.Position.Workspace)
	if err != nil {
		return result, fmt.Errorf("run produced-state boundaries: %w", err)
	}
	result.State.Records.Bytes = result.Boundary.Pipeline.Final.Arena
	return result, nil
}
