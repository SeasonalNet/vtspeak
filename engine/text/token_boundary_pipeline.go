package text

import (
	"fmt"

	"vtspeak/engine/tree3"
)

// Paul2013TokenBoundaryPipelineResult retains the native pre-tree groups,
// reverse marker-tree scan, and final boundary arena/state outputs.
type Paul2013TokenBoundaryPipelineResult struct {
	Prepared   Paul2013PreparedTokenBoundaryArena
	MarkerTree Paul2013MarkerTreeArenaResult
	Final      Paul2013TokenBoundaryArenaResult
}

// RunPaul2013TokenBoundaryPipeline composes the recovered FUN_100130e0
// sequence. Pointer mapping, the scalar tree selected at model +0x200, and
// workspace array initialization remain explicit inputs. Unsupported phone
// symbols fail at group preparation rather than bypassing that native stage.
func RunPaul2013TokenBoundaryPipeline(arena []byte, arenaAddress uint32, tree *tree3.Tree, resolvePointer Paul2013CStringPointerResolver, initialValues, initialStates, perTokenOutputs []int32) (Paul2013TokenBoundaryPipelineResult, error) {
	return runPaul2013TokenBoundaryPipeline(arena, arenaAddress, tree, resolvePointer, func(prepared []byte) (Paul2013TokenBoundaryArenaResult, error) {
		return FinishPaul2013TokenBoundaryArena(prepared, initialValues, initialStates, perTokenOutputs)
	})
}

// RunTokenBoundaries connects the finalized context-to-record handoff with
// grouping, marker-tree evaluation, parser-output extraction, and dispatch.
func (state Paul2013FinalizedModelState) RunTokenBoundaries(arenaAddress uint32, parserStateOffset int, tree *tree3.Tree, resolvePointer Paul2013CStringPointerResolver, initialValues, initialStates []int32) (Paul2013TokenBoundaryPipelineResult, error) {
	return runPaul2013TokenBoundaryPipeline(state.Records.Bytes, arenaAddress, tree, resolvePointer, func(prepared []byte) (Paul2013TokenBoundaryArenaResult, error) {
		return state.FinishTokenBoundaries(prepared, parserStateOffset, initialValues, initialStates)
	})
}

func runPaul2013TokenBoundaryPipeline(arena []byte, arenaAddress uint32, tree *tree3.Tree, resolvePointer Paul2013CStringPointerResolver, finish func([]byte) (Paul2013TokenBoundaryArenaResult, error)) (Paul2013TokenBoundaryPipelineResult, error) {
	result := Paul2013TokenBoundaryPipelineResult{}
	var err error
	result.Prepared, err = PreparePaul2013TokenBoundaryArena(arena, arenaAddress)
	if err != nil {
		return result, fmt.Errorf("prepare token boundaries: %w", err)
	}
	result.MarkerTree, err = RunPaul2013MarkerTreeArena(result.Prepared.PhoneGroups.Arena, tree, resolvePointer)
	if err != nil {
		return result, fmt.Errorf("scan token marker tree: %w", err)
	}
	result.Final, err = finish(result.MarkerTree.Arena)
	if err != nil {
		return result, fmt.Errorf("finish token boundaries: %w", err)
	}
	result.Final.Arena, result.Final.Groups, err = PopulatePaul2013RecordGroupDescriptorArena(result.Final.Arena, arenaAddress)
	if err != nil {
		return result, fmt.Errorf("write final boundary descriptors: %w", err)
	}
	return result, nil
}
