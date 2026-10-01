package text

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	paul2013FinalizerContextPrefix = 0x0b
	paul2013FinalizerContextStride = 0x70
	paul2013FinalizerGroupStride   = 0x94
)

// Paul2013ModelParserFinalizerInput describes the byte arenas consumed by
// FUN_1003e240 after its inner FUN_1000e2f0 parser has already run. The state
// arena must include the bytes preceding ParserStateOffset because the native
// routine writes into its caller's surrounding object. ContextTable begins at
// model+0x429a2, including its 11-byte prefix before the first 0x70-byte row.
// ModelResult is copied to parser-state offset +2 after row processing.
type Paul2013ModelParserFinalizerInput struct {
	StateArena        []byte
	ParserStateOffset int
	ContextTable      []byte
	ContextRowCount   int
	ModelResult       int16
}

// Paul2013ModelParserFinalizerResult returns the mutated caller arena and the
// native finalizer's low-short success value. ReturnValue is 1 when the
// expected parser-row count equals the number of output groups, and -1
// otherwise. The earlier inner-parser error path is not represented here.
type Paul2013ModelParserFinalizerResult struct {
	StateArena       []byte
	OutputGroupCount int
	ReturnValue      int16
}

// FinalizePaul2013ModelParserRows ports the post-FUN_1000e2f0 row pass in
// FUN_1003e240. It models the native surrounding memory as an explicit arena,
// preserving its signed offsets and C-string operations while checking every
// read and write. It does not run either source parser or construct its rows.
func FinalizePaul2013ModelParserRows(
	input Paul2013ModelParserFinalizerInput,
) (Paul2013ModelParserFinalizerResult, error) {
	if input.ContextRowCount < -0x8000 || input.ContextRowCount > 0x7fff {
		return Paul2013ModelParserFinalizerResult{}, fmt.Errorf(
			"model context row count %d is outside the signed-short range [-32768, 32767]",
			input.ContextRowCount,
		)
	}
	rowCount := input.ContextRowCount
	if rowCount < 0 {
		rowCount = 0
	}
	contextBytes := 0
	if rowCount > 0 {
		contextBytes = paul2013FinalizerContextPrefix + rowCount*paul2013FinalizerContextStride
	}
	if len(input.ContextTable) < contextBytes {
		return Paul2013ModelParserFinalizerResult{}, fmt.Errorf(
			"model context table has %d bytes, need %d for %d rows",
			len(input.ContextTable), contextBytes, rowCount,
		)
	}
	if input.ParserStateOffset < 0x78 || input.ParserStateOffset > len(input.StateArena) ||
		len(input.StateArena)-input.ParserStateOffset < 4 {
		return Paul2013ModelParserFinalizerResult{}, fmt.Errorf(
			"parser state offset %#x does not leave room for native surrounding fields",
			input.ParserStateOffset,
		)
	}
	state := append([]byte(nil), input.StateArena...)
	destination := input.ParserStateOffset - 0x4c
	expectedGroups := int(int16(binary.LittleEndian.Uint16(state[input.ParserStateOffset:])))
	groupCount := 0
	truncateRemainder := false
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		rowStart := paul2013FinalizerContextPrefix + rowIndex*paul2013FinalizerContextStride
		row := input.ContextTable[rowStart : rowStart+paul2013FinalizerContextStride]
		newGroup := rowIndex == 0 || binary.LittleEndian.Uint16(input.ContextTable[rowStart-0x75:rowStart-0x73]) !=
			binary.LittleEndian.Uint16(input.ContextTable[rowStart-5:rowStart-3])
		if newGroup {
			if err := finalizerWriteByte(state, destination+0x90, 0xff); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			surfaceMode, err := finalizerReadByte(state, destination+0x83)
			if err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			surfaceClass, err := finalizerReadByte(state, destination+0x84)
			if err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if surfaceMode == 'U' || surfaceClass != 'Y' {
				if err := finalizerCopyCString(state, destination+0xb2, row, 0x1e); err != nil {
					return Paul2013ModelParserFinalizerResult{}, fmt.Errorf("copy model surface for row %d: %w", rowIndex, err)
				}
			}
			priorCount, err := finalizerReadInt32(state, destination+0x70)
			if err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if err := finalizerWriteInt32(state, destination+0x70, priorCount+1); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if err := finalizerWriteByte(state, destination+int(priorCount)+0x91, 0xff); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if err := finalizerWriteByte(state, destination+0x82, input.ContextTable[rowStart-1]); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if err := finalizerWriteUint16(state, destination+0x7c, binary.LittleEndian.Uint16(input.ContextTable[rowStart-7:rowStart-5])); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			groupCount++
			truncateRemainder = false
			destination += paul2013FinalizerGroupStride
			continue
		}

		if row[0] == 0xa2 && row[1] == 0xfe {
			if !truncateRemainder {
				needle, err := finalizerCString(row, 0)
				if err != nil {
					return Paul2013ModelParserFinalizerResult{}, fmt.Errorf("read delimiter row %d: %w", rowIndex, err)
				}
				if err := finalizerTruncateAtCString(state, destination, needle); err != nil {
					return Paul2013ModelParserFinalizerResult{}, fmt.Errorf("truncate source at delimiter row %d: %w", rowIndex, err)
				}
				count, err := finalizerReadInt32(state, destination-0x2c)
				if err != nil {
					return Paul2013ModelParserFinalizerResult{}, err
				}
				if err := finalizerWriteInt32(state, destination-0x2c, count-2); err != nil {
					return Paul2013ModelParserFinalizerResult{}, err
				}
				truncateRemainder = true
			}
			continue
		}

		if truncateRemainder {
			if err := finalizerAppendCString(state, destination, row, 0); err != nil {
				return Paul2013ModelParserFinalizerResult{}, fmt.Errorf("append source row %d: %w", rowIndex, err)
			}
		}
		classCount, err := finalizerReadInt32(state, destination-0x24)
		if err != nil {
			return Paul2013ModelParserFinalizerResult{}, err
		}
		if classCount < 5 {
			if row[0x1e] != 0 {
				if err := finalizerAppendCString(state, destination+0x1e, row, 0x1e); err != nil {
					return Paul2013ModelParserFinalizerResult{}, fmt.Errorf("append phone text for row %d: %w", rowIndex, err)
				}
			}
			if err := finalizerWriteByte(state, destination+int(classCount)-4, 0xff); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if err := finalizerWriteInt32(state, destination-0x24, classCount+1); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			if err := finalizerWriteByte(state, destination+int(classCount)-3, 0xff); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			flags, err := finalizerReadUint16(state, destination-0x18)
			if err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
			phoneFlags := binary.LittleEndian.Uint16(input.ContextTable[rowStart-7 : rowStart-5])
			if err := finalizerWriteUint16(state, destination-0x18, flags|phoneFlags); err != nil {
				return Paul2013ModelParserFinalizerResult{}, err
			}
		}
	}
	if rowCount > 0 {
		if err := finalizerWriteUint16(state, input.ParserStateOffset+2, uint16(input.ModelResult)); err != nil {
			return Paul2013ModelParserFinalizerResult{}, err
		}
	}
	returnValue := int16(-1)
	if expectedGroups == groupCount {
		returnValue = 1
	}
	return Paul2013ModelParserFinalizerResult{
		StateArena:       state,
		OutputGroupCount: groupCount,
		ReturnValue:      returnValue,
	}, nil
}

func finalizerReadByte(arena []byte, offset int) (byte, error) {
	if offset < 0 || offset >= len(arena) {
		return 0, fmt.Errorf("native byte read at %#x exceeds arena size %d", offset, len(arena))
	}
	return arena[offset], nil
}

func finalizerWriteByte(arena []byte, offset int, value byte) error {
	if offset < 0 || offset >= len(arena) {
		return fmt.Errorf("native byte write at %#x exceeds arena size %d", offset, len(arena))
	}
	arena[offset] = value
	return nil
}

func finalizerReadInt32(arena []byte, offset int) (int32, error) {
	if offset < 0 || offset+4 > len(arena) {
		return 0, fmt.Errorf("native dword read at %#x exceeds arena size %d", offset, len(arena))
	}
	return int32(binary.LittleEndian.Uint32(arena[offset : offset+4])), nil
}

func finalizerWriteInt32(arena []byte, offset int, value int32) error {
	if offset < 0 || offset+4 > len(arena) {
		return fmt.Errorf("native dword write at %#x exceeds arena size %d", offset, len(arena))
	}
	binary.LittleEndian.PutUint32(arena[offset:offset+4], uint32(value))
	return nil
}

func finalizerReadUint16(arena []byte, offset int) (uint16, error) {
	if offset < 0 || offset+2 > len(arena) {
		return 0, fmt.Errorf("native word read at %#x exceeds arena size %d", offset, len(arena))
	}
	return binary.LittleEndian.Uint16(arena[offset : offset+2]), nil
}

func finalizerWriteUint16(arena []byte, offset int, value uint16) error {
	if offset < 0 || offset+2 > len(arena) {
		return fmt.Errorf("native word write at %#x exceeds arena size %d", offset, len(arena))
	}
	binary.LittleEndian.PutUint16(arena[offset:offset+2], value)
	return nil
}

func finalizerCString(arena []byte, offset int) ([]byte, error) {
	if offset < 0 || offset >= len(arena) {
		return nil, fmt.Errorf("native C-string read at %#x exceeds arena size %d", offset, len(arena))
	}
	end := bytes.IndexByte(arena[offset:], 0)
	if end < 0 {
		return nil, errors.New("native C string is not NUL terminated within its arena")
	}
	return arena[offset : offset+end], nil
}

func finalizerCopyCString(destination []byte, destinationOffset int, source []byte, sourceOffset int) error {
	value, err := finalizerCString(source, sourceOffset)
	if err != nil {
		return err
	}
	return finalizerWriteCString(destination, destinationOffset, value)
}

func finalizerAppendCString(destination []byte, destinationOffset int, source []byte, sourceOffset int) error {
	current, err := finalizerCString(destination, destinationOffset)
	if err != nil {
		return err
	}
	suffix, err := finalizerCString(source, sourceOffset)
	if err != nil {
		return err
	}
	combined := make([]byte, 0, len(current)+len(suffix))
	combined = append(combined, current...)
	combined = append(combined, suffix...)
	return finalizerWriteCString(destination, destinationOffset, combined)
}

func finalizerWriteCString(destination []byte, offset int, value []byte) error {
	if offset < 0 || offset+len(value)+1 > len(destination) {
		return fmt.Errorf("native C-string write at %#x with %d bytes exceeds arena size %d", offset, len(value)+1, len(destination))
	}
	copy(destination[offset:offset+len(value)], value)
	destination[offset+len(value)] = 0
	return nil
}

func finalizerTruncateAtCString(arena []byte, sourceOffset int, needle []byte) error {
	if len(needle) == 0 {
		return errors.New("native delimiter string is empty")
	}
	source, err := finalizerCString(arena, sourceOffset)
	if err != nil {
		return err
	}
	index := bytes.Index(source, needle)
	if index < 0 {
		return errors.New("native delimiter string was not found in source text")
	}
	return finalizerWriteByte(arena, sourceOffset+index, 0)
}
