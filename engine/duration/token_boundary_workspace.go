package duration

import (
	"context"
	"fmt"

	"vtspeak/engine/text"
)

// RunPaul2013ModelTokenBoundaries uses the shared engbi.tree3 provisioned by
// OpenPaul2013 and the native workspace value/state fields. The model record
// arena, parser-state base, and pointer resolver remain explicit.
func (engine *Engine) RunPaul2013ModelTokenBoundaries(state text.Paul2013FinalizedModelState, arenaAddress uint32, parserStateOffset int, resolvePointer text.Paul2013CStringPointerResolver, workspace []byte) (text.Paul2013TokenBoundaryWorkspaceResult, error) {
	if engine == nil {
		return text.Paul2013TokenBoundaryWorkspaceResult{}, fmt.Errorf("Paul 2013 engine is nil")
	}
	return state.RunTokenBoundariesWithWorkspace(arenaAddress, parserStateOffset, engine.catalog, resolvePointer, workspace)
}

// RunPaul2013PositionSegmentBoundaryDriver supplies the loaded shared tree
// to the native segment/state/record/boundary driver. The source parser and
// native pointer ownership remain explicit dependencies.
func (engine *Engine) RunPaul2013PositionSegmentBoundaryDriver(ctx context.Context, engineContext, workspace, sharedObject, modelArena []byte, parse text.Paul2013PositionSegmentParser, resolveStates text.Paul2013Int32PointerResolver, resolveAddress text.Paul2013ParserRowAddressResolver, arenaAddress uint32, resolveStrings text.Paul2013CStringPointerResolver) (text.Paul2013PositionSegmentBoundaryDriverResult, error) {
	if engine == nil || engine.catalog == nil {
		return text.Paul2013PositionSegmentBoundaryDriverResult{}, fmt.Errorf("Paul 2013 shared boundary resources are not loaded")
	}
	return text.RunPaul2013PositionSegmentBoundaryDriver(ctx, engineContext, workspace, sharedObject, modelArena, parse, resolveStates, resolveAddress, arenaAddress, engine.catalog, resolveStrings)
}

// RunPaul2013PositionStateBoundaries uses the loaded shared tree while
// composing position-state production, interval mapping, record flags/mode,
// and boundary/workspace writeback.
func (engine *Engine) RunPaul2013PositionStateBoundaries(state text.Paul2013FinalizedModelState, workspace []byte, program text.Paul2013PositionStateProgram, arenaAddress uint32, parserStateOffset int, resolvePointer text.Paul2013CStringPointerResolver) (text.Paul2013PositionBoundaryPipelineResult, error) {
	if engine == nil || engine.catalog == nil {
		return text.Paul2013PositionBoundaryPipelineResult{}, fmt.Errorf("Paul 2013 shared boundary resources are not loaded")
	}
	return state.ApplyPositionStatesAndRunTokenBoundaries(workspace, program, arenaAddress, parserStateOffset, engine.catalog, resolvePointer)
}

// RunPaul2013PositionDescriptorBoundaries reads native producer descriptors
// and clamp tables before composing the loaded model boundary pipeline.
func (engine *Engine) RunPaul2013PositionDescriptorBoundaries(state text.Paul2013FinalizedModelState, workspace []byte, resolveStates text.Paul2013Int32PointerResolver, indexTables *text.Paul2013PositionIndexTables, arenaAddress uint32, parserStateOffset int, resolveStrings text.Paul2013CStringPointerResolver) (text.Paul2013PositionBoundaryPipelineResult, error) {
	if engine == nil || engine.catalog == nil {
		return text.Paul2013PositionBoundaryPipelineResult{}, fmt.Errorf("Paul 2013 shared boundary resources are not loaded")
	}
	return state.RunPositionDescriptorBoundaries(workspace, resolveStates, indexTables, arenaAddress, parserStateOffset, engine.catalog, resolveStrings)
}

// RunPaul2013NativePositionWorkspaceBoundaries also derives index-table
// pointers and the remap gate from native workspace/shared-object snapshots.
func (engine *Engine) RunPaul2013NativePositionWorkspaceBoundaries(state text.Paul2013FinalizedModelState, workspace, sharedObject []byte, resolveStates text.Paul2013Int32PointerResolver, arenaAddress uint32, parserStateOffset int, resolveStrings text.Paul2013CStringPointerResolver) (text.Paul2013PositionBoundaryPipelineResult, error) {
	if engine == nil || engine.catalog == nil {
		return text.Paul2013PositionBoundaryPipelineResult{}, fmt.Errorf("Paul 2013 shared boundary resources are not loaded")
	}
	return state.RunNativePositionWorkspaceBoundaries(workspace, sharedObject, resolveStates, arenaAddress, parserStateOffset, engine.catalog, resolveStrings)
}

// RunPaul2013PreparedPositionWorkspaceBoundaries also selects and clamps
// native engine-context/workspace segment defaults before state processing.
func (engine *Engine) RunPaul2013PreparedPositionWorkspaceBoundaries(state text.Paul2013FinalizedModelState, engineContext, workspace, sharedObject []byte, resolveStates text.Paul2013Int32PointerResolver, arenaAddress uint32, parserStateOffset int, resolveStrings text.Paul2013CStringPointerResolver) (text.Paul2013PositionBoundaryPipelineResult, error) {
	if engine == nil || engine.catalog == nil {
		return text.Paul2013PositionBoundaryPipelineResult{}, fmt.Errorf("Paul 2013 shared boundary resources are not loaded")
	}
	return state.RunPreparedNativePositionWorkspaceBoundaries(engineContext, workspace, sharedObject, resolveStates, arenaAddress, parserStateOffset, engine.catalog, resolveStrings)
}
