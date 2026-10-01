package text

import (
	"encoding/binary"
	"fmt"
)

// Paul2013ModelRecordPitchInput associates one native group row with the
// eleven-short input produced by FUN_10013a20.
type Paul2013ModelRecordPitchInput struct {
	GroupIndex int
	PhoneIndex int
	Values     [Paul2013PitchInputCount]int16
}

// BuildPaul2013ModelRecordPitchInputs composes FUN_10012c70/FUN_10013c00
// group production with FUN_10013a20's one-input-per-counted-group layout.
// Record order and its preceding terminal marker are supplied because they
// belong to the enclosing model-state arena.
func BuildPaul2013ModelRecordPitchInputs(
	record []byte,
	state Paul2013ModelRecordDurationState,
) ([]Paul2013ModelRecordPitchInput, error) {
	if len(record) < paul2013MarkerRecordSize {
		return nil, fmt.Errorf("model-state record has %d bytes, need %d", len(record), paul2013MarkerRecordSize)
	}
	if state.RecordIndex < 0 {
		return nil, fmt.Errorf("model-record index %d is negative", state.RecordIndex)
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
	inputs := make([]Paul2013ModelRecordPitchInput, 0, groups.NativeGroupCount)
	for groupIndex := 0; groupIndex < groups.NativeGroupCount; groupIndex++ {
		row := groups.NativeGroupRows[groupIndex]
		if len(row) < 0x1e {
			return nil, fmt.Errorf("native group row %d has %d bytes, need 30", groupIndex, len(row))
		}
		phoneIndex := int(row[0x1c])
		groupPhones := int(row[0x1d])
		if groupPhones == 0 || phoneIndex > phoneCount || groupPhones > phoneCount-phoneIndex {
			return nil, fmt.Errorf("native group row %d covers invalid phone span [%d,%d) of %d", groupIndex, phoneIndex, phoneIndex+groupPhones, phoneCount)
		}
		class, ok := paul2013PitchPhoneClass[phones[phoneIndex].Label]
		if !ok {
			return nil, fmt.Errorf("model-record group %d first phone %d (%q) has no recovered pitch class", groupIndex, phoneIndex, phones[phoneIndex].Label)
		}
		nucleusOffset := groupPhones - 1
		for offset := 0; offset < groupPhones; offset++ {
			if groups.PhoneLabels[phoneIndex+offset] == 2 {
				nucleusOffset = offset
				break
			}
		}
		nucleus, err := phones[phoneIndex+nucleusOffset].TreeFeatures()
		if err != nil {
			return nil, fmt.Errorf("model-record group %d nucleus phone: %w", groupIndex, err)
		}
		rowInput := [Paul2013PitchInputCount]int16{
			class,
			nucleus.IdentityOrdinal,
			int16(row[0]),
			3,
			3,
			int16(row[2]),
			0,
			4,
			2,
			int16(row[1]),
			int16(groups.NativeGroupCount),
		}
		if groupIndex > 0 {
			previous := groups.NativeGroupRows[groupIndex-1]
			rowInput[3] = int16(previous[0])
			rowInput[6] = int16(previous[2])
		}
		if groupIndex+1 < groups.NativeGroupCount {
			next := groups.NativeGroupRows[groupIndex+1]
			rowInput[4] = int16(next[0])
			rowInput[7] = int16(next[2])
		}
		firstGroup := groupIndex == 0
		if firstGroup && (state.RecordIndex == 0 || state.PreviousTerminalMarker == '[') {
			rowInput[8] = 1
		} else if groupIndex == groups.NativeGroupCount-1 &&
			Paul2013DurationBoundaryOrdinal(record[paul2013MarkerTerminalOffset]) == 40 {
			rowInput[8] = 3
		} else {
			rowInput[8] = 2
		}
		inputs = append(inputs, Paul2013ModelRecordPitchInput{
			GroupIndex: groupIndex, PhoneIndex: phoneIndex, Values: rowInput,
		})
	}
	return inputs, nil
}

// BuildPaul2013ModelRecordStreamPitchInputs derives per-record pitch inputs
// and preceding terminal markers from an ordered packed record stream.
func BuildPaul2013ModelRecordStreamPitchInputs(records []byte) ([][]Paul2013ModelRecordPitchInput, error) {
	const recordStride = 0x3c0
	const terminalOffset = 0x3bd
	if len(records)%recordStride != 0 {
		return nil, fmt.Errorf("model-record stream has %d bytes, not a multiple of %d", len(records), recordStride)
	}
	recordCount := len(records) / recordStride
	result := make([][]Paul2013ModelRecordPitchInput, recordCount)
	for recordIndex := 0; recordIndex < recordCount; recordIndex++ {
		start := recordIndex * recordStride
		state := Paul2013ModelRecordDurationState{RecordIndex: recordIndex}
		if recordIndex > 0 {
			state.PreviousTerminalMarker = records[start-recordStride+terminalOffset]
		}
		inputs, err := BuildPaul2013ModelRecordPitchInputs(records[start:start+recordStride], state)
		if err != nil {
			return nil, fmt.Errorf("model-record stream entry %d: %w", recordIndex, err)
		}
		result[recordIndex] = inputs
	}
	return result, nil
}

// BuildPaul2013ModelStateArenaPitchInputs reads the signed record count at
// model-state arena +2 and records at +0x64c, then composes per-group pitch
// inputs in native record order.
func BuildPaul2013ModelStateArenaPitchInputs(modelStateArena []byte) ([][]Paul2013ModelRecordPitchInput, error) {
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
		return [][]Paul2013ModelRecordPitchInput{}, nil
	}
	if recordCount > paul2013RecordGroupCapacity {
		return nil, fmt.Errorf("model-state arena record count %d exceeds native capacity %d", recordCount, paul2013RecordGroupCapacity)
	}
	if recordOffset > len(modelStateArena) || recordCount > (len(modelStateArena)-recordOffset)/recordStride {
		return nil, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(modelStateArena), recordCount, recordOffset)
	}
	end := recordOffset + recordCount*recordStride
	return BuildPaul2013ModelRecordStreamPitchInputs(modelStateArena[recordOffset:end])
}
