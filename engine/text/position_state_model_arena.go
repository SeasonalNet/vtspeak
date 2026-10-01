package text

import (
	"encoding/binary"
	"fmt"
)

// ApplyPaul2013PositionStateProgramFromModelArena derives FUN_10022970's
// row keys, selector-specific gaps, and final-segment gate from model records
// and workspace fields, then applies the program with native cursor writeback.
// segmentByteCount is parser-state +4. Event values, clamp tables, and range
// boundary producers remain inputs. IndexTables are applied after the state
// passes, as in the native caller; intervals used by this scan precede mapping.
func ApplyPaul2013PositionStateProgramFromModelArena(workspace, modelArena []byte, segmentByteCount int32, input Paul2013PositionStateProgram) (Paul2013PositionStateWorkspaceResult, error) {
	count, err := readPaul2013PositionRecordCount(modelArena)
	if err != nil {
		return Paul2013PositionStateWorkspaceResult{}, err
	}
	if count == 0 {
		return Paul2013PositionStateWorkspaceResult{}, fmt.Errorf("model position-state pass requires at least one record")
	}
	if len(workspace) < 8 {
		return Paul2013PositionStateWorkspaceResult{}, fmt.Errorf("workspace lacks source length and segment position")
	}
	starts, ends := make([]int32, count), make([]int32, count)
	for index := range starts {
		row := modelArena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		starts[index] = int32(binary.LittleEndian.Uint32(row[:4]))
		ends[index] = int32(binary.LittleEndian.Uint32(row[4:8]))
	}
	segmentEnd := int32(binary.LittleEndian.Uint32(workspace[4:8])) + segmentByteCount
	input.RowCount = count
	input.ModelStateRows = modelArena
	ranges := make(map[byte]Paul2013PositionStateRanges, len(input.Ranges))
	for mode, pass := range input.Ranges {
		pass.RowKeys = append([]int32(nil), starts...)
		ranges[mode] = pass
	}
	input.Ranges = ranges
	events := make(map[byte]Paul2013PositionEventValues, len(input.Events))
	for mode, pass := range input.Events {
		if mode != 3 && mode != 4 {
			return Paul2013PositionStateWorkspaceResult{}, fmt.Errorf("unsupported model event mode %d", mode)
		}
		pass.Ranges = make([]Paul2013PositionEventRange, count)
		for index := range starts {
			if mode == 3 {
				minimum := int32(0)
				if index > 0 {
					minimum = ends[index-1]
				}
				pass.Ranges[index] = Paul2013PositionEventRange{Minimum: minimum, Maximum: starts[index]}
			} else {
				maximum := segmentEnd
				if index+1 < count {
					maximum = starts[index+1]
				}
				pass.Ranges[index] = Paul2013PositionEventRange{Minimum: ends[index], Maximum: maximum}
			}
		}
		events[mode] = pass
	}
	input.Events = events
	if pass, ok := events[3]; ok && len(pass.Boundaries) > 0 {
		input.Terminal = &Paul2013PositionTerminalAccumulator{Enabled: segmentEnd == int32(binary.LittleEndian.Uint32(workspace[:4])), MinimumBoundary: ends[count-1], MaximumBoundary: segmentEnd, Boundaries: pass.Boundaries, Values: pass.Values, Minimum: pass.Minimum, Maximum: pass.Maximum}
	} else {
		input.Terminal = nil
	}
	return ApplyPaul2013PositionStateProgramToWorkspace(workspace, input)
}

// Position processing precedes phone/control record population in the native
// caller, so only count and record storage are required here.
func readPaul2013PositionRecordCount(arena []byte) (int, error) {
	if len(arena) < 4 {
		return 0, fmt.Errorf("position arena lacks record count")
	}
	count := int(int16(binary.LittleEndian.Uint16(arena[2:4])))
	if count < 0 || count > paul2013ParserStateRowLimit {
		return 0, fmt.Errorf("position record count %d outside native capacity", count)
	}
	if count > 0 && len(arena) < paul2013TokenMarkerArenaRecords+count*paul2013TokenMarkerArenaStride {
		return 0, fmt.Errorf("position arena lacks %d record slots", count)
	}
	return count, nil
}
