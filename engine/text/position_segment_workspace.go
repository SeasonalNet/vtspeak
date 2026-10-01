package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// PreparePaul2013PositionSegmentWorkspace ports FUN_10022850. Each exact
// flag byte 1 selects a workspace override; other bytes select the engine
// context field. The three signed scalars are clamped, copied to active and
// default fields, and workspace +0x1312c4 is reset to -1. Inputs are preserved.
func PreparePaul2013PositionSegmentWorkspace(engineContext, workspace []byte) ([]byte, error) {
	if len(engineContext) < 0x4cfc || len(workspace) < 0x1312c8 {
		return nil, fmt.Errorf("segment default context/workspace fields are truncated")
	}
	result := append([]byte(nil), workspace...)
	for _, field := range []struct {
		flag, override, context, active, initial int
		minimum, maximum                         int32
	}{{0x1210d7, 0x1210dc, 0x4cf4, 0x1210e8, 0x1210fc, 50, 200}, {0x1210d8, 0x1210e0, 0x4cf0, 0x1210ec, 0x1210f4, 50, 400}, {0x1210d9, 0x1210e4, 0x4cf8, 0x1210f0, 0x1210f8, 0, 500}} {
		value := int32(binary.LittleEndian.Uint32(engineContext[field.context:]))
		if workspace[field.flag] == 1 {
			value = int32(binary.LittleEndian.Uint32(workspace[field.override:]))
		}
		value = min(max(value, field.minimum), field.maximum)
		binary.LittleEndian.PutUint32(result[field.active:], uint32(value))
		binary.LittleEndian.PutUint32(result[field.initial:], uint32(value))
	}
	binary.LittleEndian.PutUint32(result[0x1312c4:], 0xffffffff)
	return result, nil
}

// RunPreparedNativePositionWorkspaceBoundaries additionally prepares the
// native segment defaults. The already-finalized parser/model records remain
// inputs: this does not invoke the unrecovered source parser or segment loop.
func (state Paul2013FinalizedModelState) RunPreparedNativePositionWorkspaceBoundaries(engineContext, workspace, sharedObject []byte, resolveStates Paul2013Int32PointerResolver, arenaAddress uint32, parserStateOffset int, catalog *tree3.Catalog, resolveStrings Paul2013CStringPointerResolver) (Paul2013PositionBoundaryPipelineResult, error) {
	prepared, err := PreparePaul2013PositionSegmentWorkspace(engineContext, workspace)
	if err != nil {
		return Paul2013PositionBoundaryPipelineResult{State: state}, err
	}
	return state.RunNativePositionWorkspaceBoundaries(prepared, sharedObject, resolveStates, arenaAddress, parserStateOffset, catalog, resolveStrings)
}
