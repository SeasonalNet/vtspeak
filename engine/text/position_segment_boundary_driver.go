package text

import (
	"context"
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// Paul2013PositionSegmentBoundaryDriverResult retains the segment driver,
// populated phone/control records and the subsequent boundary writeback.
// HasBoundary is false for the null-source and end-of-source branches.
type Paul2013PositionSegmentBoundaryDriverResult struct {
	Segment     Paul2013PositionSegmentDriverResult
	Records     Paul2013ModelStateRecordArena
	Boundary    Paul2013TokenBoundaryWorkspaceResult
	HasBoundary bool
}

// PopulateModelStateRecords connects the successful native position driver
// to FUN_10016c90 using its produced slot-7 array and retained parser state.
// It preserves mapped intervals and does not invent a parser-finalizer result.
func (segment Paul2013PositionSegmentDriverResult) PopulateModelStateRecords(resolveAddress Paul2013ParserRowAddressResolver) (Paul2013ModelStateRecordArena, error) {
	count, err := readPaul2013PositionRecordCount(segment.ModelArena)
	if err != nil {
		return Paul2013ModelStateRecordArena{}, err
	}
	if segment.EndOfSource || segment.ParserCalls == 0 || count == 0 || len(segment.ParserState) < count*paul2013ParserStateRowStride || len(segment.Position.Program.Arrays[7]) != count {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("segment has no successful counted parser/state handoff")
	}
	if int(int16(binary.LittleEndian.Uint16(segment.ParserState))) != count {
		return Paul2013ModelStateRecordArena{}, fmt.Errorf("segment parser and model counts differ")
	}
	return PopulatePaul2013ModelStateRecordsFromParserRows(segment.ModelArena, segment.ParserState[:count*paul2013ParserStateRowStride], segment.Position.Program.Arrays[7], int16(binary.LittleEndian.Uint16(segment.ParserState[2:4])), resolveAddress)
}

// RunPaul2013PositionSegmentBoundaryDriver composes the recovered segment
// driver with phone/control projection and FUN_100130e0 boundary processing.
// The unresolved parser is a callback; native address ownership and pointer
// resolution remain explicit. Completion branches never evaluate the tree.
func RunPaul2013PositionSegmentBoundaryDriver(ctx context.Context, engineContext, workspace, sharedObject, modelArena []byte, parse Paul2013PositionSegmentParser, resolveStates Paul2013Int32PointerResolver, resolveAddress Paul2013ParserRowAddressResolver, arenaAddress uint32, catalog *tree3.Catalog, resolveStrings Paul2013CStringPointerResolver) (Paul2013PositionSegmentBoundaryDriverResult, error) {
	var result Paul2013PositionSegmentBoundaryDriverResult
	var err error
	result.Segment, err = RunPaul2013PositionSegmentDriver(ctx, engineContext, workspace, sharedObject, modelArena, parse, resolveStates)
	if err != nil {
		return result, err
	}
	if result.Segment.EndOfSource || result.Segment.ParserCalls == 0 {
		return result, nil
	}
	if catalog == nil || catalog.Pronunciation == nil {
		return result, fmt.Errorf("shared engbi.tree3 is not loaded")
	}
	result.Records, err = result.Segment.PopulateModelStateRecords(resolveAddress)
	if err != nil {
		return result, err
	}
	count := len(result.Segment.Position.Program.Arrays[7])
	values, states := make([]int32, count), make([]int32, count)
	outputs := make([]int32, count-1)
	for index := range values {
		values[index] = int32(binary.LittleEndian.Uint32(result.Segment.Workspace[paul2013BoundaryValuesWorkspaceOffset+index*4:]))
		states[index] = int32(binary.LittleEndian.Uint32(result.Segment.Workspace[paul2013BoundaryStatesWorkspaceOffset+index*4:]))
		if index < len(outputs) {
			outputs[index] = int32(binary.LittleEndian.Uint32(result.Segment.ParserState[0x20+index*paul2013ParserStateRowStride:]))
		}
	}
	result.Boundary.Pipeline, err = RunPaul2013TokenBoundaryPipeline(result.Records.Bytes, arenaAddress, catalog.Pronunciation, resolveStrings, values, states, outputs)
	if err != nil {
		return result, err
	}
	result.Boundary.Workspace = append([]byte(nil), result.Segment.Workspace...)
	for index := range values {
		binary.LittleEndian.PutUint32(result.Boundary.Workspace[paul2013BoundaryValuesWorkspaceOffset+index*4:], uint32(result.Boundary.Pipeline.Final.State.Values[index]))
		binary.LittleEndian.PutUint32(result.Boundary.Workspace[paul2013BoundaryStatesWorkspaceOffset+index*4:], uint32(result.Boundary.Pipeline.Final.State.States[index]))
	}
	result.Records.Bytes = result.Boundary.Pipeline.Final.Arena
	result.HasBoundary = true
	return result, nil
}
