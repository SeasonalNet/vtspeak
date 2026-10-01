package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// ReadPaul2013PositionIndexWorkspaceTables derives FUN_10022dc0's lookup
// pointers from workspace +0x47790/+0x47794/+0x47798, its ordinary clamp
// bound from workspace +0, and its remap branch from shared-object +0x20424.
// It resolves only the table prefixes needed by the supplied model records.
// Native table allocation and source-index producers remain caller-owned.
func ReadPaul2013PositionIndexWorkspaceTables(workspace, sharedObject, modelArena []byte, resolve Paul2013Int32PointerResolver) (Paul2013PositionIndexTables, error) {
	var tables Paul2013PositionIndexTables
	if len(workspace) < 0x4779c || len(sharedObject) <= 0x20424 {
		return tables, fmt.Errorf("native position mapping workspace/shared fields are truncated")
	}
	count, err := readPaul2013PositionRecordCount(modelArena)
	if err != nil {
		return tables, err
	}
	tables.MaximumIndex = int32(binary.LittleEndian.Uint32(workspace[:4]))
	tables.RemappedMode = sharedObject[0x20424] == 1
	if count == 0 {
		return tables, nil
	}
	if !tables.RemappedMode && tables.MaximumIndex <= 0 {
		return tables, fmt.Errorf("ordinary position mapping requires positive source length")
	}
	var startCount, endCount int64
	for index := 0; index < count; index++ {
		row := modelArena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		start, end := int32(binary.LittleEndian.Uint32(row)), int32(binary.LittleEndian.Uint32(row[4:]))
		if !tables.RemappedMode {
			start = clampPaul2013PositionIndex(start, tables.MaximumIndex)
			end = clampPaul2013PositionIndex(end, tables.MaximumIndex)
		}
		if start < 0 || end < 0 {
			return tables, fmt.Errorf("remapped record %d has negative source index", index)
		}
		startCount = max(startCount, int64(start)+1)
		endCount = max(endCount, int64(end)+1)
	}
	// Prefix byte sizes exceeding native 32-bit addressing are invalid even
	// on a host whose int can hold their cell counts.
	if startCount > 0x3fffffff || endCount > 0x3fffffff {
		return tables, fmt.Errorf("position table prefix exceeds native address space")
	}
	tables.Start, err = resolvePaul2013Int32Cells(resolve, binary.LittleEndian.Uint32(workspace[0x47794:]), int(startCount))
	if err != nil {
		return tables, fmt.Errorf("position start table: %w", err)
	}
	tables.End, err = resolvePaul2013Int32Cells(resolve, binary.LittleEndian.Uint32(workspace[0x47798:]), int(endCount))
	if err != nil {
		return tables, fmt.Errorf("position end table: %w", err)
	}
	if tables.RemappedMode {
		// Only values actually indexed by records determine final-table size;
		// unused entries in the resolved prefixes may contain arbitrary data.
		var finalCount int64
		for index := 0; index < count; index++ {
			row := modelArena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
			start := tables.Start[int32(binary.LittleEndian.Uint32(row))]
			end := tables.End[int32(binary.LittleEndian.Uint32(row[4:]))]
			if start < 0 || end < 0 {
				return tables, fmt.Errorf("record %d has negative final-table index", index)
			}
			finalCount = max(finalCount, int64(start)+1, int64(end)+1)
		}
		if finalCount > 0x3fffffff {
			return tables, fmt.Errorf("final position table prefix exceeds native address space")
		}
		tables.Final, err = resolvePaul2013Int32Cells(resolve, binary.LittleEndian.Uint32(workspace[0x47790:]), int(finalCount))
		if err != nil {
			return tables, fmt.Errorf("final position table: %w", err)
		}
	}
	return tables, nil
}

// RunNativePositionWorkspaceBoundaries derives descriptors, clamp limits,
// mapping tables and the shared remap gate from native memory snapshots,
// then runs position-state and loaded boundary processing in caller order.
// On success it advances workspace +4 by the parser's consumed-byte count,
// using native 32-bit wrap semantics. Zero consumed bytes indicate the native
// end-of-source branch, which requires the separate source-parser driver.
func (state Paul2013FinalizedModelState) RunNativePositionWorkspaceBoundaries(workspace, sharedObject []byte, resolveStates Paul2013Int32PointerResolver, arenaAddress uint32, parserStateOffset int, catalog *tree3.Catalog, resolveStrings Paul2013CStringPointerResolver) (Paul2013PositionBoundaryPipelineResult, error) {
	if parserStateOffset < 0 || parserStateOffset > len(state.Finalized.StateArena) || len(state.Finalized.StateArena)-parserStateOffset < 8 {
		return Paul2013PositionBoundaryPipelineResult{State: state}, fmt.Errorf("native segment parser header is truncated")
	}
	consumed := binary.LittleEndian.Uint32(state.Finalized.StateArena[parserStateOffset+4:])
	if consumed == 0 {
		return Paul2013PositionBoundaryPipelineResult{State: state}, fmt.Errorf("zero-byte parser segment requires the native end-of-source driver")
	}
	tables, err := ReadPaul2013PositionIndexWorkspaceTables(workspace, sharedObject, state.Records.Bytes, resolveStates)
	if err != nil {
		return Paul2013PositionBoundaryPipelineResult{State: state}, err
	}
	result, err := state.RunPositionDescriptorBoundaries(workspace, resolveStates, &tables, arenaAddress, parserStateOffset, catalog, resolveStrings)
	if err != nil {
		return result, err
	}
	next := binary.LittleEndian.Uint32(workspace[4:8]) + consumed
	binary.LittleEndian.PutUint32(result.Position.Workspace[4:8], next)
	binary.LittleEndian.PutUint32(result.Boundary.Workspace[4:8], next)
	return result, nil
}
