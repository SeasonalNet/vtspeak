package duration

import (
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

// PitchResult retains both paired tree lookups for one phone. ScalarValue is
// the raw tree result; StoredScalarValue is the signed byte written to the
// phone record and reused as feature 11 of the vector tree.
type PitchResult struct {
	Input             [12]int16
	ScalarTree        string
	ScalarLeaf        int
	ScalarValue       int16
	StoredScalarValue int8
	VectorTree        string
	VectorLeaf        int
	VectorTreeValues  [12]int16
	VectorValues      [12]int16
}

// TokenPitchResult aligns pitch-tree results with one selected token span.
type TokenPitchResult struct {
	Token  text.LexicalTokenSpan
	Phones []PitchResult
	// BoundaryAverages has one entry per adjacent phone pair. Each entry holds
	// the three windows written to the left record and the three written to the
	// right record by FUN_100137c0.
	BoundaryAverages [][2][3]int16
}

// BuildPaul2013PitchBoundaryAverages ports the moving-average consumer in
// FUN_100137c0. It joins the final six values of the left vector to the first
// six values of the right vector, then evaluates six overlapping windows of
// seven signed shorts. The outputs are split at the phone boundary exactly as
// the native routine writes them: the first three replace the left vector's
// final three values and the last three replace the right vector's first three.
func BuildPaul2013PitchBoundaryAverages(left, right [12]int16) [2][3]int16 {
	var joined [24]int16
	copy(joined[:12], left[:])
	copy(joined[12:], right[:])
	var averages [2][3]int16
	for window := 0; window < 6; window++ {
		average := roundedPaul2013PitchAverage(joined[6+window : 13+window])
		averages[window/3][window%3] = average
	}
	return averages
}

// SmoothPaul2013PitchVectors applies FUN_100137c0's adjacent-vector pass in
// place and returns its six averages for each boundary. The left vector's
// last three values and the right vector's first three values are overwritten
// as each pair is processed. This helper applies the recovered transform to
// an explicit vector sequence; it does not derive FUN_100137c0's descriptor
// spans or map them to model-record groups.
func SmoothPaul2013PitchVectors(vectors [][12]int16) [][2][3]int16 {
	if len(vectors) < 2 {
		return nil
	}
	averages := make([][2][3]int16, len(vectors)-1)
	for phoneIndex := range averages {
		pair := BuildPaul2013PitchBoundaryAverages(vectors[phoneIndex], vectors[phoneIndex+1])
		averages[phoneIndex] = pair
		copy(vectors[phoneIndex][9:], pair[0][:])
		copy(vectors[phoneIndex+1][:3], pair[1][:])
	}
	return averages
}

// Paul2013PitchSmoothingDescriptor is the pointer-free projection of one
// 16-byte FUN_100137c0 descriptor: Count is its leading signed short and
// Vectors supplies the rows reached through its +0x18 pointer.
type Paul2013PitchSmoothingDescriptor struct {
	Count   int16
	Vectors [][12]int16
}

// Paul2013PitchSmoothingRun contains one descriptor's copied vector rows and
// the boundary averages written by FUN_100137c0.
type Paul2013PitchSmoothingRun struct {
	Vectors  [][12]int16
	Averages [][2][3]int16
}

// Paul2013PitchSmoothingTable carries the table-level loop count and its
// pointer-free descriptor projections. FUN_100138c0 reads IterationCount
// from table offset +2; each FUN_100137c0 call reads the following descriptor
// count from that descriptor's leading signed short.
type Paul2013PitchSmoothingTable struct {
	IterationCount int16
	Descriptors    []Paul2013PitchSmoothingDescriptor
}

// Paul2013PitchSmoothingVectorResolver maps a native row-array address to
// signed 12-short vectors. Native process pointers require a caller-provided
// address-space view; the descriptor reader supplies the requested row count.
type Paul2013PitchSmoothingVectorResolver func(address uint32, rowCount int) ([][12]int16, error)

// ReadPaul2013PitchSmoothingTable projects the table fields consumed by
// FUN_100138c0 and FUN_100137c0 from a contiguous snapshot beginning at the
// table base. It reads the signed iteration count at +0x2, the descriptor
// counts at 16-byte strides, and each current row pointer at +0x18 plus that
// stride. The final descriptor supplies the last call's count but has no row
// pointer of its own.
func ReadPaul2013PitchSmoothingTable(
	tableBytes []byte,
	resolve Paul2013PitchSmoothingVectorResolver,
) (Paul2013PitchSmoothingTable, error) {
	const (
		iterationCountOffset = 0x02
		descriptorStride     = 0x10
		rowPointerBase       = 0x18
	)
	if len(tableBytes) < iterationCountOffset+2 {
		return Paul2013PitchSmoothingTable{}, fmt.Errorf("pitch smoothing table has %d bytes, need count through %#x", len(tableBytes), iterationCountOffset+2)
	}
	iterationCount := int16(binary.LittleEndian.Uint16(tableBytes[iterationCountOffset : iterationCountOffset+2]))
	table := Paul2013PitchSmoothingTable{IterationCount: iterationCount}
	if iterationCount <= 0 {
		return table, nil
	}
	iterations := int(iterationCount)
	countEnd := iterations*descriptorStride + 2
	pointerEnd := rowPointerBase + (iterations-1)*descriptorStride + 4
	required := countEnd
	if pointerEnd > required {
		required = pointerEnd
	}
	if len(tableBytes) < required {
		return Paul2013PitchSmoothingTable{}, fmt.Errorf("pitch smoothing table requests %d iterations but snapshot has %d bytes, need %d", iterations, len(tableBytes), required)
	}
	descriptors := make([]Paul2013PitchSmoothingDescriptor, iterations+1)
	for descriptorIndex := range descriptors {
		countOffset := descriptorIndex * descriptorStride
		descriptors[descriptorIndex].Count = int16(binary.LittleEndian.Uint16(tableBytes[countOffset : countOffset+2]))
	}
	for descriptorIndex := 0; descriptorIndex < iterations; descriptorIndex++ {
		rowCount := int(descriptors[descriptorIndex+1].Count)
		if rowCount <= 1 {
			continue
		}
		if resolve == nil {
			return Paul2013PitchSmoothingTable{}, errors.New("pitch smoothing table has vector rows but no native pointer resolver")
		}
		pointerOffset := rowPointerBase + descriptorIndex*descriptorStride
		address := binary.LittleEndian.Uint32(tableBytes[pointerOffset : pointerOffset+4])
		if address == 0 {
			return Paul2013PitchSmoothingTable{}, fmt.Errorf("pitch smoothing descriptor %d has a null row-array pointer for %d rows", descriptorIndex, rowCount)
		}
		vectors, err := resolve(address, rowCount)
		if err != nil {
			return Paul2013PitchSmoothingTable{}, fmt.Errorf("resolve pitch smoothing descriptor %d row pointer %#x: %w", descriptorIndex, address, err)
		}
		if len(vectors) < rowCount {
			return Paul2013PitchSmoothingTable{}, fmt.Errorf("pitch smoothing descriptor %d resolves %d rows, need %d", descriptorIndex, len(vectors), rowCount)
		}
		descriptors[descriptorIndex].Vectors = vectors
	}
	table.Descriptors = descriptors
	return table, nil
}

// SmoothPaul2013PitchDescriptorTableRuns ports FUN_100138c0's table-level
// iteration and FUN_100137c0's per-descriptor count traversal over projected
// row data. Each smoother call reads the next descriptor's count, processes
// count-1 neighboring vector pairs, and uses the current descriptor's row
// pointer. The caller supplies those pointer targets as Vectors; extraction
// and mapping of native addresses remain outside this helper.
func SmoothPaul2013PitchDescriptorTableRuns(
	table Paul2013PitchSmoothingTable,
) ([]Paul2013PitchSmoothingRun, error) {
	if table.IterationCount <= 0 {
		return []Paul2013PitchSmoothingRun{}, nil
	}
	iterationCount := int(table.IterationCount)
	if iterationCount >= len(table.Descriptors) {
		return nil, fmt.Errorf("pitch smoothing table requests %d iterations but has only %d descriptors; each iteration reads a following descriptor", iterationCount, len(table.Descriptors))
	}
	runs := make([]Paul2013PitchSmoothingRun, iterationCount)
	for descriptorIndex := range runs {
		count := int(table.Descriptors[descriptorIndex+1].Count)
		vectors := append([][12]int16(nil), table.Descriptors[descriptorIndex].Vectors...)
		if count <= 1 {
			runs[descriptorIndex] = Paul2013PitchSmoothingRun{Vectors: vectors}
			continue
		}
		if count > len(vectors) {
			return nil, fmt.Errorf("pitch smoothing descriptor %d requests %d rows but supplies %d", descriptorIndex, count, len(vectors))
		}
		averages := SmoothPaul2013PitchVectors(vectors[:count])
		runs[descriptorIndex] = Paul2013PitchSmoothingRun{Vectors: vectors, Averages: averages}
	}
	return runs, nil
}

// SmoothPaul2013PitchDescriptorRuns applies the descriptor traversal to an
// explicit pointer-free sequence, using all adjacent descriptor pairs. Use
// SmoothPaul2013PitchDescriptorTableRuns when the native table-level +2 loop
// count is available. The caller supplies pointer targets as Vectors.
func SmoothPaul2013PitchDescriptorRuns(
	descriptors []Paul2013PitchSmoothingDescriptor,
) ([]Paul2013PitchSmoothingRun, error) {
	if len(descriptors) == 0 {
		return []Paul2013PitchSmoothingRun{}, nil
	}
	if len(descriptors)-1 > int(^uint16(0)>>1) {
		return nil, fmt.Errorf("pitch smoothing descriptor sequence has %d entries, exceeding the signed table-count range", len(descriptors))
	}
	return SmoothPaul2013PitchDescriptorTableRuns(Paul2013PitchSmoothingTable{
		IterationCount: int16(len(descriptors) - 1),
		Descriptors:    descriptors,
	})
}

// roundedPaul2013PitchAverage ports FUN_10013790's signed rounding expression
// for the fixed seven-short windows used by FUN_100137c0.
func roundedPaul2013PitchAverage(values []int16) int16 {
	var sum int32
	for _, value := range values {
		sum += int32(value)
	}
	length := int32(len(values))
	return int16((length + 2*sum) / (2 * length))
}

// EvaluatePaul2013PitchText builds FUN_10013a20's eleven context values,
// evaluates the ordinary nbt/nbf pair, and uses the sbt/sbf pair for the
// final phone at the standard terminal-Z boundary. The scalar result is
// truncated to the signed byte stored by FUN_100138c0 before it is appended
// as input 11 to the vector tree.
func EvaluatePaul2013PitchText(sequence text.LexicalPhoneSequence, catalog *tree3.Catalog) ([]TokenPitchResult, error) {
	markers := make([]byte, len(sequence.Tokens))
	if len(markers) > 0 {
		markers[len(markers)-1] = 'Z'
	}
	return EvaluatePaul2013PitchTextWithMarkers(sequence, markers, catalog)
}

// EvaluatePaul2013PitchTextWithMarkers evaluates the paired pitch trees using
// one terminal marker per lexical token. FUN_100138c0 selects the alternate
// tree pair only for the final phone in that token. The same markers also
// select the terminal position state in the eleven-value input row.
func EvaluatePaul2013PitchTextWithMarkers(
	sequence text.LexicalPhoneSequence,
	terminalMarkers []byte,
	catalog *tree3.Catalog,
) ([]TokenPitchResult, error) {
	return EvaluatePaul2013PitchTextWithPhoneMarkers(sequence, nil, terminalMarkers, catalog)
}

// EvaluatePaul2013PitchTextWithPhoneMarkers evaluates pitch trees using phone
// groups rebuilt from explicit FUN_10012c70 block markers. Terminal markers
// still control per-token pitch-tree dispatch.
func EvaluatePaul2013PitchTextWithPhoneMarkers(
	sequence text.LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	catalog *tree3.Catalog,
) ([]TokenPitchResult, error) {
	if len(terminalMarkers) != len(sequence.Tokens) {
		return nil, fmt.Errorf("received %d pitch terminal markers for %d tokens", len(terminalMarkers), len(sequence.Tokens))
	}
	inputs, err := text.BuildPaul2013TokenPitchInputsWithPhoneMarkers(sequence, phoneMarkers, terminalMarkers)
	if err != nil {
		return nil, fmt.Errorf("build Paul 2013 pitch inputs: %w", err)
	}
	return evaluatePaul2013PitchInputs(inputs, terminalMarkers, catalog)
}

// EvaluatePaul2013PitchTextWithPositionStates evaluates pitch trees using
// caller-supplied per-phone position states shared with duration evaluation.
func EvaluatePaul2013PitchTextWithPositionStates(
	sequence text.LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	positionStates []uint8,
	catalog *tree3.Catalog,
) ([]TokenPitchResult, error) {
	if len(terminalMarkers) != len(sequence.Tokens) {
		return nil, fmt.Errorf("received %d pitch terminal markers for %d tokens", len(terminalMarkers), len(sequence.Tokens))
	}
	inputs, err := text.BuildPaul2013TokenPitchInputsWithPositionStates(
		sequence, phoneMarkers, terminalMarkers, positionStates,
	)
	if err != nil {
		return nil, fmt.Errorf("build Paul 2013 pitch inputs: %w", err)
	}
	return evaluatePaul2013PitchInputs(inputs, terminalMarkers, catalog)
}

func evaluatePaul2013PitchInputs(
	inputs []text.Paul2013TokenPitchInputs,
	terminalMarkers []byte,
	catalog *tree3.Catalog,
) ([]TokenPitchResult, error) {
	if catalog == nil {
		return nil, errors.New("Paul 2013 pitch catalog is nil")
	}
	results := make([]TokenPitchResult, len(inputs))
	for tokenIndex, tokenInputs := range inputs {
		phoneResults := make([]PitchResult, len(tokenInputs.Inputs))
		for phoneIndex, input := range tokenInputs.Inputs {
			terminalPhone := phoneIndex == len(tokenInputs.Inputs)-1
			scalarName, vectorName := paul2013PitchTreePair(terminalMarkers[tokenIndex], terminalPhone)
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
			scalarLeaf, scalarOutput, err := scalarTree.Evaluate(input[:])
			if err != nil {
				return nil, fmt.Errorf("evaluate %s for token %d phone %d: %w", scalarName, tokenIndex, phoneIndex, err)
			}
			if len(scalarOutput) != 1 {
				return nil, fmt.Errorf("scalar pitch tree %q returned %d values, want one", scalarName, len(scalarOutput))
			}
			storedScalar := int8(uint8(scalarOutput[0]))
			vectorInput := [12]int16{}
			copy(vectorInput[:], input[:])
			vectorInput[11] = int16(storedScalar)
			vectorLeaf, vectorOutput, err := vectorTree.Evaluate(vectorInput[:])
			if err != nil {
				return nil, fmt.Errorf("evaluate %s for token %d phone %d: %w", vectorName, tokenIndex, phoneIndex, err)
			}
			if len(vectorOutput) != len(vectorInput) {
				return nil, fmt.Errorf("vector pitch tree %q returned %d values, want %d", vectorName, len(vectorOutput), len(vectorInput))
			}
			var vectorValues [12]int16
			copy(vectorValues[:], vectorOutput)
			phoneResults[phoneIndex] = PitchResult{
				Input:             vectorInput,
				ScalarTree:        scalarName,
				ScalarLeaf:        scalarLeaf,
				ScalarValue:       scalarOutput[0],
				StoredScalarValue: storedScalar,
				VectorTree:        vectorName,
				VectorLeaf:        vectorLeaf,
				VectorTreeValues:  vectorValues,
				VectorValues:      vectorValues,
			}
		}
		vectorValues := make([][12]int16, len(phoneResults))
		for phoneIndex := range phoneResults {
			vectorValues[phoneIndex] = phoneResults[phoneIndex].VectorValues
		}
		boundaryAverages := SmoothPaul2013PitchVectors(vectorValues)
		for phoneIndex := range phoneResults {
			phoneResults[phoneIndex].VectorValues = vectorValues[phoneIndex]
		}
		results[tokenIndex] = TokenPitchResult{Token: tokenInputs.Token, Phones: phoneResults, BoundaryAverages: boundaryAverages}
	}
	return results, nil
}

func paul2013PitchTreePair(marker byte, terminalPhone bool) (string, string) {
	if !terminalPhone {
		return "nbt.tree3", "nbf.tree3"
	}
	switch marker {
	case '^':
		return "bt.tree3", "bf.tree3"
	case 'Z':
		return "sbt.tree3", "sbf.tree3"
	case '[':
		return "qbt.tree3", "qbf.tree3"
	default:
		return "nbt.tree3", "nbf.tree3"
	}
}
