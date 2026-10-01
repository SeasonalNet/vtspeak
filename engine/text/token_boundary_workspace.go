package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

const (
	paul2013BoundaryValuesWorkspaceOffset = 0x121a64
	paul2013BoundaryStatesWorkspaceOffset = 0x121d80
)

var paul2013PositionStateWorkspaceOffsets = [8]int{0x121100, 0x121420, 0x121740, 0x121a60, 0x121d80, 0, 0, 0x1220a0}

// Paul2013PositionStateWorkspaceResult retains the recovered state passes
// and native destination bytes. Workspace memory beyond each initialized
// row prefix is preserved, including the boundary-value tail cell.
type Paul2013PositionStateWorkspaceResult struct {
	Workspace []byte
	Program   Paul2013PositionStateProgramResult
}

// ApplyPaul2013PositionStateProgramToWorkspace reads defaults from the native
// +0x1210fc/+0x1210f4/+0x1210f8 fields, runs the supplied range/event program,
// and writes its six recovered arrays to their native workspace offsets.
// The supplied program's Initial and nonempty event-pass starting cursors are
// replaced by native workspace fields. Event cursors are written back at
// +0x1223f4/+0x122404; terminal processing continues selector 3's same cursor.
// Row keys, event tables, clamps, and remapping inputs remain explicit.
func ApplyPaul2013PositionStateProgramToWorkspace(workspace []byte, input Paul2013PositionStateProgram) (Paul2013PositionStateWorkspaceResult, error) {
	if input.RowCount < 0 || input.RowCount > paul2013ParserStateRowLimit {
		return Paul2013PositionStateWorkspaceResult{}, fmt.Errorf("workspace state row count %d outside native capacity", input.RowCount)
	}
	needed := 0x1220a0 + input.RowCount*4
	if (input.Terminal != nil || len(input.Events[3].Boundaries) != 0) && needed < 0x122444 {
		needed = 0x122444
	}
	if len(input.Events[4].Boundaries) != 0 && needed < 0x122408 {
		needed = 0x122408
	}
	if len(workspace) < needed {
		return Paul2013PositionStateWorkspaceResult{}, fmt.Errorf("position-state workspace has %d bytes, need %d", len(workspace), needed)
	}
	input.Initial = [3]int32{int32(binary.LittleEndian.Uint32(workspace[0x1210fc:])), int32(binary.LittleEndian.Uint32(workspace[0x1210f4:])), int32(binary.LittleEndian.Uint32(workspace[0x1210f8:]))}
	if input.Events != nil {
		passes := make(map[byte]Paul2013PositionEventValues, len(input.Events))
		for mode, pass := range input.Events {
			if (mode == 3 || mode == 4) && len(pass.Boundaries) != 0 {
				pass.StartBoundary = int(int32(binary.LittleEndian.Uint32(workspace[0x1223c4+int(mode)*0x10:])))
				// The legacy overflow field is a diagnostic, not a separate
				// native destination. Start both views at the native cursor.
				pass.OverflowStart = pass.StartBoundary
			}
			passes[mode] = pass
		}
		input.Events = passes
	}
	program, err := ApplyPaul2013PositionStateProgram(input)
	if err != nil {
		return Paul2013PositionStateWorkspaceResult{}, err
	}
	result := Paul2013PositionStateWorkspaceResult{Workspace: append([]byte(nil), workspace...), Program: program}
	for _, mode := range []int{0, 1, 2, 3, 4, 7} {
		for index, value := range program.Arrays[mode] {
			binary.LittleEndian.PutUint32(result.Workspace[paul2013PositionStateWorkspaceOffsets[mode]+index*4:], uint32(value))
		}
	}
	for mode, event := range program.Events {
		if len(input.Events[mode].Boundaries) != 0 {
			binary.LittleEndian.PutUint32(result.Workspace[0x1223c4+int(mode)*0x10:], uint32(event.NextBoundary))
			if mode == 3 {
				binary.LittleEndian.PutUint32(result.Workspace[0x122440:], 0xffffffff)
			}
		}
	}
	if program.HasTerminalPass {
		binary.LittleEndian.PutUint32(result.Workspace[0x122440:], uint32(program.Terminal.Value))
		binary.LittleEndian.PutUint32(result.Workspace[0x1223f4:], uint32(program.Terminal.NextBoundary))
	}
	return result, nil
}

// Paul2013TokenBoundaryWorkspaceResult retains the complete supported
// boundary pipeline and its native workspace value/state writes.
type Paul2013TokenBoundaryWorkspaceResult struct {
	Pipeline  Paul2013TokenBoundaryPipelineResult
	Workspace []byte
}

// RunTokenBoundariesWithWorkspace uses the shared engbi.tree3 already loaded
// as Catalog.Pronunciation. FUN_10012280 formats that filename and passes it
// to FUN_10001000, which loads the +0x200 tree used by FUN_10012f00.
// Values begin one cell into state slot 3, so the final initial value comes
// from the retained workspace tail, not a fabricated default.
func (state Paul2013FinalizedModelState) RunTokenBoundariesWithWorkspace(arenaAddress uint32, parserStateOffset int, catalog *tree3.Catalog, resolvePointer Paul2013CStringPointerResolver, workspace []byte) (Paul2013TokenBoundaryWorkspaceResult, error) {
	if catalog == nil || catalog.Pronunciation == nil {
		return Paul2013TokenBoundaryWorkspaceResult{}, fmt.Errorf("shared engbi.tree3 is not loaded")
	}
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(state.Records.Bytes)
	if err != nil {
		return Paul2013TokenBoundaryWorkspaceResult{}, err
	}
	count := len(inputs.InitialMarkers)
	if len(workspace) < paul2013BoundaryStatesWorkspaceOffset+count*4 {
		return Paul2013TokenBoundaryWorkspaceResult{}, fmt.Errorf("boundary workspace cannot hold %d native states", count)
	}
	values, states := make([]int32, count), make([]int32, count)
	for index := range values {
		values[index] = int32(binary.LittleEndian.Uint32(workspace[paul2013BoundaryValuesWorkspaceOffset+index*4:]))
		states[index] = int32(binary.LittleEndian.Uint32(workspace[paul2013BoundaryStatesWorkspaceOffset+index*4:]))
	}
	pipeline, err := state.RunTokenBoundaries(arenaAddress, parserStateOffset, catalog.Pronunciation, resolvePointer, values, states)
	if err != nil {
		return Paul2013TokenBoundaryWorkspaceResult{Pipeline: pipeline}, err
	}
	result := Paul2013TokenBoundaryWorkspaceResult{Pipeline: pipeline, Workspace: append([]byte(nil), workspace...)}
	for index := range values {
		binary.LittleEndian.PutUint32(result.Workspace[paul2013BoundaryValuesWorkspaceOffset+index*4:], uint32(pipeline.Final.State.Values[index]))
		binary.LittleEndian.PutUint32(result.Workspace[paul2013BoundaryStatesWorkspaceOffset+index*4:], uint32(pipeline.Final.State.States[index]))
	}
	return result, nil
}
