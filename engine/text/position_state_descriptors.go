package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// Paul2013Int32PointerResolver returns count signed little-endian int32 cells
// from a native pointer. Address ownership remains with the caller; returned
// cells are copied before use. No native addresses are dereferenced by Go.
type Paul2013Int32PointerResolver func(address uint32, count int) ([]int32, error)

// FUN_10022970 reads these tables at VA 0x1007d66c and 0x1007d690.
var paul2013PositionMaximum = [8]int32{200, 400, 500, 65535, 3, 0, 9, 9}
var paul2013PositionMinimum = [8]int32{50, 50, 0, 0, 0, 0, 0, 0}

// ReadPaul2013PositionStateDescriptorProgram reads the six descriptors at
// workspace +0x1223c0 + selector*0x10: count, cursor, boundary pointer, value
// pointer. It supplies native clamp values, retains event cursors, and skips
// inactive descriptors using FUN_10022970's pointer/count gates. Row keys,
// event gaps and terminal processing are filled by the model-arena adapter.
// The producer that fills these descriptors remains outside this reader.
func ReadPaul2013PositionStateDescriptorProgram(workspace []byte, resolve Paul2013Int32PointerResolver, indexTables *Paul2013PositionIndexTables) (Paul2013PositionStateProgram, error) {
	program := Paul2013PositionStateProgram{Ranges: make(map[byte]Paul2013PositionStateRanges), Events: make(map[byte]Paul2013PositionEventValues), IndexTables: indexTables}
	if len(workspace) < 0x122440 {
		return program, fmt.Errorf("workspace lacks native position-state descriptors")
	}
	for _, mode := range []byte{0, 1, 2, 7, 3, 4} {
		descriptor := workspace[0x1223c0+int(mode)*0x10:]
		count := int(int32(binary.LittleEndian.Uint32(descriptor)))
		boundaryPointer := binary.LittleEndian.Uint32(descriptor[8:])
		valuePointer := binary.LittleEndian.Uint32(descriptor[12:])
		isRange := mode != 3 && mode != 4
		if count <= 0 || boundaryPointer == 0 || (isRange && valuePointer == 0) {
			continue
		}
		boundaries, err := resolvePaul2013Int32Cells(resolve, boundaryPointer, count)
		if err != nil {
			return program, fmt.Errorf("selector %d boundaries: %w", mode, err)
		}
		valueCount := count
		if isRange {
			valueCount-- // Native range loop reads values[boundaryIndex-1].
		}
		values, err := resolvePaul2013Int32Cells(resolve, valuePointer, valueCount)
		if err != nil {
			return program, fmt.Errorf("selector %d values: %w", mode, err)
		}
		if isRange {
			program.Ranges[mode] = Paul2013PositionStateRanges{Mode: mode, Boundaries: boundaries, Values: values, Minimum: paul2013PositionMinimum[mode], Maximum: paul2013PositionMaximum[mode]}
		} else {
			cursor := int(int32(binary.LittleEndian.Uint32(descriptor[4:])))
			program.Events[mode] = Paul2013PositionEventValues{Mode: mode, Boundaries: boundaries, Values: values, StartBoundary: cursor, OverflowStart: cursor, Minimum: paul2013PositionMinimum[mode], Maximum: paul2013PositionMaximum[mode]}
		}
	}
	return program, nil
}

func resolvePaul2013Int32Cells(resolve Paul2013Int32PointerResolver, address uint32, count int) ([]int32, error) {
	if count == 0 {
		return nil, nil
	}
	if count < 0 || address == 0 || resolve == nil || uint64(address)+uint64(count)*4 > 1<<32 {
		return nil, fmt.Errorf("cannot resolve %d int32 cells at %#x", count, address)
	}
	values, err := resolve(address, count)
	if err != nil {
		return nil, err
	}
	if len(values) < count {
		return nil, fmt.Errorf("pointer %#x returned %d cells, need %d", address, len(values), count)
	}
	return append([]int32(nil), values[:count]...), nil
}

// RunPositionDescriptorBoundaries connects native descriptors to record
// state, optional interval mapping, and the loaded engbi boundary pipeline.
func (state Paul2013FinalizedModelState) RunPositionDescriptorBoundaries(workspace []byte, resolveStates Paul2013Int32PointerResolver, indexTables *Paul2013PositionIndexTables, arenaAddress uint32, parserStateOffset int, catalog *tree3.Catalog, resolveStrings Paul2013CStringPointerResolver) (Paul2013PositionBoundaryPipelineResult, error) {
	program, err := ReadPaul2013PositionStateDescriptorProgram(workspace, resolveStates, indexTables)
	if err != nil {
		return Paul2013PositionBoundaryPipelineResult{State: state}, err
	}
	return state.ApplyPositionStatesAndRunTokenBoundaries(workspace, program, arenaAddress, parserStateOffset, catalog, resolveStrings)
}
