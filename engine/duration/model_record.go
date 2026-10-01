package duration

import (
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

// Paul2013ModelRecordPitchResult aligns a paired pitch-tree evaluation with
// its native group row and first physical phone.
type Paul2013ModelRecordPitchResult struct {
	GroupIndex int
	PhoneIndex int
	Pitch      PitchResult
}

// EvaluatePaul2013ModelRecordDuration evaluates the native per-phone duration
// vectors composed from one 0x3c0-byte model-state record. Record order and
// the preceding terminal marker remain explicit because they are held by the
// enclosing model arena, outside the record.
func EvaluatePaul2013ModelRecordDuration(
	record []byte,
	state text.Paul2013ModelRecordDurationState,
	catalog *tree3.Catalog,
) ([]PhoneResult, error) {
	if catalog == nil {
		return nil, errors.New("Paul 2013 duration catalog is nil")
	}
	inputs, err := text.BuildPaul2013ModelRecordDurationInputs(record, state)
	if err != nil {
		return nil, fmt.Errorf("build model-record duration inputs: %w", err)
	}
	phoneCount := int(record[0x95])
	phones, err := text.DecodeCMUPhones(record[0x2e8 : 0x2e8+phoneCount])
	if err != nil {
		return nil, fmt.Errorf("decode model-record phones: %w", err)
	}
	results := make([]PhoneResult, len(inputs))
	for inputIndex, input := range inputs {
		phoneIndex := input.PhoneIndex
		phone := phones[phoneIndex]
		treeName, err := Paul2013TreeName(phone)
		if err != nil {
			return nil, fmt.Errorf("model-record phone %d: %w", phoneIndex, err)
		}
		tree := catalog.Duration[treeName]
		if tree == nil {
			return nil, fmt.Errorf("Paul 2013 duration tree %q is not loaded", treeName)
		}
		leaf, output, err := tree.Evaluate(input.Values[:])
		if err != nil {
			return nil, fmt.Errorf("evaluate %s for model-record phone %d: %w", treeName, phoneIndex, err)
		}
		if len(output) != 1 {
			return nil, fmt.Errorf("duration tree %q returned %d values, want one", treeName, len(output))
		}
		results[inputIndex] = PhoneResult{
			Phone: phone, PhoneIndex: phoneIndex, TreeName: treeName, LeafOrdinal: leaf,
			Value: output[0], Input: input.Values,
		}
	}
	return results, nil
}

// EvaluatePaul2013ModelRecordStreamDuration evaluates all records in a packed
// model-record stream and derives each record's enclosing boundary state from
// stream order.
func EvaluatePaul2013ModelRecordStreamDuration(
	records []byte,
	catalog *tree3.Catalog,
) ([][]PhoneResult, error) {
	const recordStride = 0x3c0
	if catalog == nil {
		return nil, errors.New("Paul 2013 duration catalog is nil")
	}
	if len(records)%recordStride != 0 {
		return nil, fmt.Errorf("model-record stream has %d bytes, not a multiple of %d", len(records), recordStride)
	}
	recordCount := len(records) / recordStride
	results := make([][]PhoneResult, recordCount)
	for recordIndex := 0; recordIndex < recordCount; recordIndex++ {
		start := recordIndex * recordStride
		record := records[start : start+recordStride]
		state := text.Paul2013ModelRecordDurationState{RecordIndex: recordIndex}
		if recordIndex > 0 {
			state.PreviousTerminalMarker = records[start-recordStride+0x3bd]
		}
		phones, err := EvaluatePaul2013ModelRecordDuration(record, state, catalog)
		if err != nil {
			return nil, fmt.Errorf("evaluate model-record stream entry %d: %w", recordIndex, err)
		}
		results[recordIndex] = phones
	}
	return results, nil
}

// EvaluatePaul2013ModelStateArenaDuration reads the signed record count at +2
// and ordered model-state records at +0x64c, then evaluates their duration trees.
func EvaluatePaul2013ModelStateArenaDuration(
	modelStateArena []byte,
	catalog *tree3.Catalog,
) ([][]PhoneResult, error) {
	const (
		countOffset  = 2
		recordOffset = 0x64c
		recordStride = 0x3c0
	)
	if catalog == nil {
		return nil, errors.New("Paul 2013 duration catalog is nil")
	}
	if len(modelStateArena) < countOffset+2 {
		return nil, fmt.Errorf("model-state arena has %d bytes, need a signed record count at +%#x", len(modelStateArena), countOffset)
	}
	recordCount := int(int16(binary.LittleEndian.Uint16(modelStateArena[countOffset : countOffset+2])))
	if recordCount < 1 {
		return [][]PhoneResult{}, nil
	}
	if recordCount > 100 {
		return nil, fmt.Errorf("model-state arena record count %d exceeds native capacity 100", recordCount)
	}
	if recordOffset > len(modelStateArena) || recordCount > (len(modelStateArena)-recordOffset)/recordStride {
		return nil, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(modelStateArena), recordCount, recordOffset)
	}
	end := recordOffset + recordCount*recordStride
	return EvaluatePaul2013ModelRecordStreamDuration(modelStateArena[recordOffset:end], catalog)
}

// EvaluatePaul2013ModelRecordPitch evaluates the native per-group pitch
// vectors for one model-state record. It returns raw tree outputs before the
// separate FUN_100137c0 boundary smoothing pass.
func EvaluatePaul2013ModelRecordPitch(
	record []byte,
	state text.Paul2013ModelRecordDurationState,
	catalog *tree3.Catalog,
) ([]Paul2013ModelRecordPitchResult, error) {
	if catalog == nil {
		return nil, errors.New("Paul 2013 pitch catalog is nil")
	}
	inputs, err := text.BuildPaul2013ModelRecordPitchInputs(record, state)
	if err != nil {
		return nil, fmt.Errorf("build model-record pitch inputs: %w", err)
	}
	results := make([]Paul2013ModelRecordPitchResult, len(inputs))
	for inputIndex, input := range inputs {
		terminalGroup := input.GroupIndex == len(inputs)-1
		marker := record[0x3bd]
		scalarName, vectorName := paul2013PitchTreePair(marker, terminalGroup)
		scalarTree := catalog.Pitch[scalarName]
		if scalarTree == nil {
			return nil, fmt.Errorf("Paul 2013 scalar pitch tree %q is not loaded", scalarName)
		}
		vectorTree := catalog.Pitch[vectorName]
		if vectorTree == nil {
			return nil, fmt.Errorf("Paul 2013 vector pitch tree %q is not loaded", vectorName)
		}
		if scalarTree.OutputWidth != 1 || vectorTree.OutputWidth != 12 {
			return nil, fmt.Errorf("Paul 2013 pitch pair %q/%q has output widths %d/%d, want 1/12", scalarName, vectorName, scalarTree.OutputWidth, vectorTree.OutputWidth)
		}
		scalarLeaf, scalarOutput, err := scalarTree.Evaluate(input.Values[:])
		if err != nil {
			return nil, fmt.Errorf("evaluate %s for model-record group %d: %w", scalarName, input.GroupIndex, err)
		}
		if len(scalarOutput) != 1 {
			return nil, fmt.Errorf("scalar pitch tree %q returned %d values, want one", scalarName, len(scalarOutput))
		}
		storedScalar := int8(uint8(scalarOutput[0]))
		vectorInput := [12]int16{}
		copy(vectorInput[:], input.Values[:])
		vectorInput[11] = int16(storedScalar)
		vectorLeaf, vectorOutput, err := vectorTree.Evaluate(vectorInput[:])
		if err != nil {
			return nil, fmt.Errorf("evaluate %s for model-record group %d: %w", vectorName, input.GroupIndex, err)
		}
		if len(vectorOutput) != len(vectorInput) {
			return nil, fmt.Errorf("vector pitch tree %q returned %d values, want %d", vectorName, len(vectorOutput), len(vectorInput))
		}
		var vectorValues [12]int16
		copy(vectorValues[:], vectorOutput)
		results[inputIndex] = Paul2013ModelRecordPitchResult{
			GroupIndex: input.GroupIndex,
			PhoneIndex: input.PhoneIndex,
			Pitch: PitchResult{
				Input: vectorInput, ScalarTree: scalarName, ScalarLeaf: scalarLeaf,
				ScalarValue: scalarOutput[0], StoredScalarValue: storedScalar,
				VectorTree: vectorName, VectorLeaf: vectorLeaf,
				VectorTreeValues: vectorValues, VectorValues: vectorValues,
			},
		}
	}
	return results, nil
}

// EvaluatePaul2013ModelRecordStreamPitch evaluates per-group pitch trees over
// an ordered packed record stream, deriving each preceding terminal marker.
func EvaluatePaul2013ModelRecordStreamPitch(
	records []byte,
	catalog *tree3.Catalog,
) ([][]Paul2013ModelRecordPitchResult, error) {
	const recordStride = 0x3c0
	if catalog == nil {
		return nil, errors.New("Paul 2013 pitch catalog is nil")
	}
	if len(records)%recordStride != 0 {
		return nil, fmt.Errorf("model-record stream has %d bytes, not a multiple of %d", len(records), recordStride)
	}
	recordCount := len(records) / recordStride
	results := make([][]Paul2013ModelRecordPitchResult, recordCount)
	for recordIndex := 0; recordIndex < recordCount; recordIndex++ {
		start := recordIndex * recordStride
		state := text.Paul2013ModelRecordDurationState{RecordIndex: recordIndex}
		if recordIndex > 0 {
			state.PreviousTerminalMarker = records[start-recordStride+0x3bd]
		}
		pitch, err := EvaluatePaul2013ModelRecordPitch(records[start:start+recordStride], state, catalog)
		if err != nil {
			return nil, fmt.Errorf("evaluate model-record stream entry %d: %w", recordIndex, err)
		}
		results[recordIndex] = pitch
	}
	return results, nil
}

// EvaluatePaul2013ModelStateArenaPitch reads the native signed record count
// and record base from the full model-state arena, then evaluates each
// counted group's scalar and vector pitch trees.
func EvaluatePaul2013ModelStateArenaPitch(
	modelStateArena []byte,
	catalog *tree3.Catalog,
) ([][]Paul2013ModelRecordPitchResult, error) {
	const (
		countOffset  = 2
		recordOffset = 0x64c
		recordStride = 0x3c0
	)
	if catalog == nil {
		return nil, errors.New("Paul 2013 pitch catalog is nil")
	}
	if len(modelStateArena) < countOffset+2 {
		return nil, fmt.Errorf("model-state arena has %d bytes, need a signed record count at +%#x", len(modelStateArena), countOffset)
	}
	recordCount := int(int16(binary.LittleEndian.Uint16(modelStateArena[countOffset : countOffset+2])))
	if recordCount < 1 {
		return [][]Paul2013ModelRecordPitchResult{}, nil
	}
	if recordCount > 100 {
		return nil, fmt.Errorf("model-state arena record count %d exceeds native capacity 100", recordCount)
	}
	if recordOffset > len(modelStateArena) || recordCount > (len(modelStateArena)-recordOffset)/recordStride {
		return nil, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(modelStateArena), recordCount, recordOffset)
	}
	end := recordOffset + recordCount*recordStride
	return EvaluatePaul2013ModelRecordStreamPitch(modelStateArena[recordOffset:end], catalog)
}
