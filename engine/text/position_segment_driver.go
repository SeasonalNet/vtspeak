package text

import (
	"context"
	"encoding/binary"
	"fmt"
)

// Paul2013PositionSegmentParser supplies the unresolved FUN_1003d350 parser
// boundary. It receives the native source address and an owned prepared
// workspace. Return parser state starts at the native parser-state base;
// Workspace is optional and carries parser-side workspace writes when present.
type Paul2013PositionSegmentParser func(context.Context, uint32, []byte) (Paul2013PositionSegmentParserOutput, error)

type Paul2013PositionSegmentParserOutput struct {
	Workspace []byte
	State     []byte
}

// Paul2013PositionSegmentDriverResult retains FUN_10022dc0's next counted
// segment or end-of-source branch. Parser output and all returned arenas are
// copied. A null source returns zero without calling the parser or preparing
// defaults, and does not set EndOfSource.
type Paul2013PositionSegmentDriverResult struct {
	Workspace   []byte
	ModelArena  []byte
	ParserState []byte
	Position    Paul2013PositionStateWorkspaceResult
	ReturnValue int32
	ParserCalls int
	EndOfSource bool
}

// RunPaul2013PositionSegmentDriver ports FUN_10022dc0's preparation, parser
// call, interval/default writes, empty-record retries, zero-byte completion,
// native state/mapping passes and successful cursor advancement. The parser
// callback and native pointer resolver remain explicit. Cancellation is checked
// before every retry; no invented native retry limit is imposed.
func RunPaul2013PositionSegmentDriver(ctx context.Context, engineContext, workspace, sharedObject, modelArena []byte, parse Paul2013PositionSegmentParser, resolve Paul2013Int32PointerResolver) (Paul2013PositionSegmentDriverResult, error) {
	result := Paul2013PositionSegmentDriverResult{Workspace: append([]byte(nil), workspace...), ModelArena: append([]byte(nil), modelArena...)}
	if len(workspace) < 12 {
		return result, fmt.Errorf("segment workspace lacks source fields")
	}
	if binary.LittleEndian.Uint32(workspace[8:12]) == 0 {
		return result, nil
	}
	if parse == nil {
		return result, fmt.Errorf("native segment parser is unavailable")
	}
	for {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		prepared, err := PreparePaul2013PositionSegmentWorkspace(engineContext, result.Workspace)
		if err != nil {
			return result, err
		}
		result.Workspace = prepared
		source := binary.LittleEndian.Uint32(prepared[8:12]) + binary.LittleEndian.Uint32(prepared[4:8])
		parsed, err := parse(ctx, source, append([]byte(nil), prepared...))
		result.ParserCalls++
		if err != nil {
			return result, fmt.Errorf("parse native segment: %w", err)
		}
		if parsed.Workspace != nil {
			if len(parsed.Workspace) != len(prepared) {
				return result, fmt.Errorf("segment parser changed workspace extent")
			}
			result.Workspace = append([]byte(nil), parsed.Workspace...)
		}
		if len(parsed.State) < 8 || len(result.ModelArena) < 4 {
			return result, fmt.Errorf("segment parser/model header is truncated")
		}
		result.ParserState = append([]byte(nil), parsed.State...)
		count := int(int16(binary.LittleEndian.Uint16(parsed.State[:2])))
		if count > paul2013ParserStateRowLimit {
			return result, fmt.Errorf("native segment count %d exceeds record capacity", count)
		}
		binary.LittleEndian.PutUint16(result.ModelArena[2:4], uint16(count))
		activeCount := max(0, count)
		if activeCount > 0 {
			if len(result.ModelArena) < paul2013TokenMarkerArenaRecords+activeCount*paul2013TokenMarkerArenaStride || len(parsed.State) < 0x1c+(activeCount-1)*paul2013ParserStateRowStride {
				return result, fmt.Errorf("segment parser/model records are truncated")
			}
			base := binary.LittleEndian.Uint32(result.Workspace[4:8])
			for index := 0; index < activeCount; index++ {
				row := parsed.State[0x14+index*paul2013ParserStateRowStride:]
				start, end := binary.LittleEndian.Uint32(row[:4]), binary.LittleEndian.Uint32(row[4:8])
				absoluteEnd := base + start
				if int32(start) < int32(end) {
					absoluteEnd = base + end - 1
				}
				record := result.ModelArena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
				binary.LittleEndian.PutUint32(record[:4], base+start)
				binary.LittleEndian.PutUint32(record[4:8], absoluteEnd)
			}
		}
		// Native defaults are initialized even when consumed bytes are zero.
		initialized, err := ApplyPaul2013PositionStateProgramToWorkspace(result.Workspace, Paul2013PositionStateProgram{RowCount: activeCount})
		if err != nil {
			return result, err
		}
		result.Workspace, result.Position = initialized.Workspace, initialized
		consumed := binary.LittleEndian.Uint32(parsed.State[4:8])
		if consumed == 0 {
			binary.LittleEndian.PutUint32(result.Workspace[0x44:0x48], 1)
			result.Position.Workspace = result.Workspace
			result.ReturnValue = int32(binary.LittleEndian.Uint32(result.Workspace[4:8]))
			result.EndOfSource = true
			return result, nil
		}
		if activeCount > 0 {
			tables, err := ReadPaul2013PositionIndexWorkspaceTables(result.Workspace, sharedObject, result.ModelArena, resolve)
			if err != nil {
				return result, err
			}
			program, err := ReadPaul2013PositionStateDescriptorProgram(result.Workspace, resolve, &tables)
			if err != nil {
				return result, err
			}
			result.Position, err = ApplyPaul2013PositionStateProgramFromModelArena(result.Workspace, result.ModelArena, int32(consumed), program)
			if err != nil {
				return result, err
			}
			result.Workspace = result.Position.Workspace
			for index, interval := range result.Position.Program.MappedIntervals {
				row := result.ModelArena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
				binary.LittleEndian.PutUint32(row[:4], uint32(interval.Minimum))
				binary.LittleEndian.PutUint32(row[4:8], uint32(interval.Maximum))
			}
		}
		next := binary.LittleEndian.Uint32(result.Workspace[4:8]) + consumed
		binary.LittleEndian.PutUint32(result.Workspace[4:8], next)
		result.Position.Workspace = result.Workspace
		result.ReturnValue = int32(next)
		if activeCount > 0 {
			return result, nil
		}
	}
}
