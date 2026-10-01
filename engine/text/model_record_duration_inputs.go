package text

import (
	"encoding/binary"
	"fmt"
)

// Paul2013ModelRecordDurationState identifies the record's position in the
// ordered model stream and the marker immediately before its first record.
type Paul2013ModelRecordDurationState struct {
	RecordIndex            int
	PreviousTerminalMarker byte
}

// Paul2013ModelRecordDurationInput associates one native counted-group phone
// with the nine-short vector sent to FUN_100135d0's duration tree.
type Paul2013ModelRecordDurationInput struct {
	PhoneIndex int
	Values     [9]int16
}

// BuildPaul2013ModelRecordDurationInputs composes native symbol decoding,
// FUN_10012c70/FUN_10013c00 group production, and the nine-short input writer.
// The record supplies phone identities/stress, group row bytes, per-phone
// onset/nucleus/coda labels, its counted group total, and the right-edge
// marker. FUN_100135d0's left and right boundary identities and position
// states are derived from record order, group offsets, and marker bytes. The
// caller supplies the preceding marker because it lives before this record in
// the enclosing model arena. Results follow counted-group order and retain
// each physical phone index; uncounted zero-vowel fallback rows emit no input.
func BuildPaul2013ModelRecordDurationInputs(
	record []byte,
	state Paul2013ModelRecordDurationState,
) ([]Paul2013ModelRecordDurationInput, error) {
	if len(record) < paul2013MarkerRecordSize {
		return nil, fmt.Errorf("model-state record has %d bytes, need %d", len(record), paul2013MarkerRecordSize)
	}
	phoneCount := int(record[paul2013MarkerPhoneCountOffset])
	if phoneCount == 0 || phoneCount > paul2013MarkerPhoneCapacity {
		return nil, fmt.Errorf("model-state record has unsupported phone count %d", phoneCount)
	}
	phones, err := DecodePaul2013ModelPhoneSymbols(record[paul2013ModelRecordPhoneSymbolOffset : paul2013ModelRecordPhoneSymbolOffset+phoneCount])
	if err != nil {
		return nil, fmt.Errorf("decode model-record phones: %w", err)
	}
	groups, err := BuildPaul2013ModelRecordPhoneGroups(record, phones)
	if err != nil {
		return nil, fmt.Errorf("build model-record phone groups: %w", err)
	}
	if state.RecordIndex < 0 {
		return nil, fmt.Errorf("model-record index %d is negative", state.RecordIndex)
	}
	if groups.NativeGroupCount == 0 {
		return []Paul2013ModelRecordDurationInput{}, nil
	}
	coveredPhones := make([]bool, phoneCount)
	for groupIndex := 0; groupIndex < groups.NativeGroupCount; groupIndex++ {
		row := groups.NativeGroupRows[groupIndex]
		if len(row) < 0x1e {
			return nil, fmt.Errorf("native group row %d has %d bytes, need 30", groupIndex, len(row))
		}
		start := int(row[0x1c])
		count := int(row[0x1d])
		if count == 0 || start > phoneCount || count > phoneCount-start {
			return nil, fmt.Errorf("native group row %d covers invalid phone span [%d,%d) of %d", groupIndex, start, start+count, phoneCount)
		}
		for offset := 0; offset < count; offset++ {
			phoneIndex := start + offset
			if coveredPhones[phoneIndex] {
				return nil, fmt.Errorf("native group rows overlap at phone %d", phoneIndex)
			}
			coveredPhones[phoneIndex] = true
		}
	}
	inputs := make([]Paul2013ModelRecordDurationInput, 0, phoneCount)
	for groupIndex := 0; groupIndex < groups.NativeGroupCount; groupIndex++ {
		groupRow := groups.NativeGroupRows[groupIndex]
		start := int(groupRow[0x1c])
		count := int(groupRow[0x1d])
		for groupPhoneIndex := 0; groupPhoneIndex < count; groupPhoneIndex++ {
			phoneIndex := start + groupPhoneIndex
			phone := phones[phoneIndex]
			firstPhoneInFirstGroup := groupIndex == 0 && groupPhoneIndex == 0
			lastPhoneInLastGroup := groupIndex == groups.NativeGroupCount-1 && groupPhoneIndex == count-1
			positionState := uint8(2)
			if firstPhoneInFirstGroup && (state.RecordIndex == 0 || state.PreviousTerminalMarker == '[') {
				positionState = 1
			} else if lastPhoneInLastGroup &&
				Paul2013DurationBoundaryOrdinal(record[paul2013MarkerTerminalOffset]) == 40 {
				positionState = 3
			}
			previous := Paul2013DurationTreeNeighbor{BoundaryIdentityOrdinal: 42}
			if firstPhoneInFirstGroup {
				if state.RecordIndex == 0 || state.PreviousTerminalMarker == '[' {
					previous.BoundaryIdentityOrdinal = 40
				}
			} else if phoneIndex > 0 {
				previous = Paul2013DurationTreeNeighbor{Phone: &phones[phoneIndex-1]}
			}
			next := Paul2013DurationTreeNeighbor{BoundaryIdentityOrdinal: Paul2013DurationBoundaryOrdinal(record[paul2013MarkerTerminalOffset])}
			if !lastPhoneInLastGroup && phoneIndex+1 < phoneCount {
				next = Paul2013DurationTreeNeighbor{Phone: &phones[phoneIndex+1]}
			}
			metadata := Paul2013DurationTreeMetadata{
				RecordByte1:   int8(groupRow[1]),
				RecordByte2:   int8(groupRow[2]),
				PositionState: positionState,
				ContextCount:  uint8(groups.NativeGroupCount),
				AuxiliaryByte: int8(groups.PhoneLabels[phoneIndex]),
			}
			input, err := BuildObservedPaul2013DurationTreeInput(previous, phone, next, metadata)
			if err != nil {
				return nil, fmt.Errorf("build model-record duration input for phone %d: %w", phoneIndex, err)
			}
			inputs = append(inputs, Paul2013ModelRecordDurationInput{PhoneIndex: phoneIndex, Values: input})
		}
	}
	return inputs, nil
}

// BuildPaul2013ModelRecordStreamDurationInputs composes duration vectors for
// an ordered run of packed native model records. Record index and each
// preceding terminal marker are derived from stream order; the first record
// uses FUN_100135d0's record-zero boundary rule.
func BuildPaul2013ModelRecordStreamDurationInputs(records []byte) ([][]Paul2013ModelRecordDurationInput, error) {
	const recordStride = 0x3c0
	const terminalOffset = 0x3bd
	if len(records)%recordStride != 0 {
		return nil, fmt.Errorf("model-record stream has %d bytes, not a multiple of %d", len(records), recordStride)
	}
	recordCount := len(records) / recordStride
	result := make([][]Paul2013ModelRecordDurationInput, recordCount)
	for recordIndex := 0; recordIndex < recordCount; recordIndex++ {
		start := recordIndex * recordStride
		record := records[start : start+recordStride]
		state := Paul2013ModelRecordDurationState{RecordIndex: recordIndex}
		if recordIndex > 0 {
			previousStart := start - recordStride
			state.PreviousTerminalMarker = records[previousStart+terminalOffset]
		}
		inputs, err := BuildPaul2013ModelRecordDurationInputs(record, state)
		if err != nil {
			return nil, fmt.Errorf("model-record stream entry %d: %w", recordIndex, err)
		}
		result[recordIndex] = inputs
	}
	return result, nil
}

// BuildPaul2013ModelStateArenaDurationInputs reads the signed record count at
// model-state arena +2 and the 0x3c0-byte records beginning at +0x64c, then
// composes their native duration vectors in record order.
func BuildPaul2013ModelStateArenaDurationInputs(modelStateArena []byte) ([][]Paul2013ModelRecordDurationInput, error) {
	const (
		countOffset  = 2
		recordOffset = 0x64c
		recordStride = 0x3c0
	)
	if len(modelStateArena) < countOffset+2 {
		return nil, fmt.Errorf("model-state arena has %d bytes, need a signed record count at +%#x", len(modelStateArena), countOffset)
	}
	recordCount := int(int16(binary.LittleEndian.Uint16(modelStateArena[countOffset : countOffset+2])))
	if recordCount < 1 {
		return [][]Paul2013ModelRecordDurationInput{}, nil
	}
	if recordCount > paul2013RecordGroupCapacity {
		return nil, fmt.Errorf("model-state arena record count %d exceeds native capacity %d", recordCount, paul2013RecordGroupCapacity)
	}
	if recordOffset > len(modelStateArena) || recordCount > (len(modelStateArena)-recordOffset)/recordStride {
		return nil, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(modelStateArena), recordCount, recordOffset)
	}
	end := recordOffset + recordCount*recordStride
	return BuildPaul2013ModelRecordStreamDurationInputs(modelStateArena[recordOffset:end])
}
