package text

import (
	"encoding/binary"
	"fmt"
)

const paul2013PhoneBlockDurationLimit = 1000

const (
	paul2013MarkerTreePhoneRowStride = 0x3c0
	paul2013MarkerTreeRowMarker      = 0x3bd
	paul2013MarkerTreeSplitCutoff    = 500
	paul2013TokenMarkerArenaCount    = 0x2
	paul2013TokenMarkerArenaRecords  = 0x64c
	paul2013TokenMarkerArenaStride   = 0x3c0
	paul2013TokenMarkerCodeOffset    = 0x3ac
	paul2013TokenMarkerStateOffset   = 0x3b0
	paul2013TokenMarkerDurationByte  = 0x95
)

var paul2013InitialMarkerByState = [...]byte{
	']', '[', 'Z', '^', 'Z', '[', ']', '[', '[', ']', '[', ']', '\\',
}

// Paul2013PhoneBlock is one half-open phone span delimited by the marker bytes
// consumed by FUN_10012c70. TerminalMarker is copied from the span's final
// phone record and is preserved for the pitch-tree dispatcher.
type Paul2013PhoneBlock struct {
	Start          int
	End            int
	TerminalMarker byte
}

// Paul2013MarkerTreeGroup is one marker-derived span used by the frontend
// marker-tree preparation. DurationByteSum is the sum of the source row's
// opaque +0x94 byte for each phone in the span.
type Paul2013MarkerTreeGroup struct {
	PhoneStart      int
	PhoneEnd        int
	DurationByteSum uint16
}

// Paul2013MarkerTreeDecision contains the supported output branch from
// FUN_10012f00. The input vector and tree selection remain separate stages.
type Paul2013MarkerTreeDecision struct {
	PhoneRow       []byte
	InsertedMarker bool
	ResetLocal8    bool
	ResetParam2    bool
}

// ApplyPaul2013MarkerTreeDecision ports FUN_10012f00's tree-result write.
// Results above 500 set the preceding record's marker at +0x3bd to '\\' and
// reset both running scan accumulators; all other rows and accumulator states
// are preserved. phoneRow is that preceding record, not the current row used
// to assemble the feature vector (x86 0x10013092 writes current base minus 3).
func ApplyPaul2013MarkerTreeDecision(phoneRow []byte, treeValue int16) (Paul2013MarkerTreeDecision, error) {
	if len(phoneRow) < paul2013MarkerTreePhoneRowStride {
		return Paul2013MarkerTreeDecision{}, fmt.Errorf(
			"Paul 2013 marker-tree phone row has %d bytes, need %d",
			len(phoneRow), paul2013MarkerTreePhoneRowStride,
		)
	}
	result := Paul2013MarkerTreeDecision{
		PhoneRow: append([]byte(nil), phoneRow...),
	}
	if treeValue > paul2013MarkerTreeSplitCutoff {
		result.PhoneRow[paul2013MarkerTreeRowMarker] = '\\'
		result.InsertedMarker = true
		result.ResetLocal8 = true
		result.ResetParam2 = true
	}
	return result, nil
}

// Paul2013TokenBoundaryMarker is one token's output from the state-to-marker
// pass in FUN_100130e0. PreviousStateFlag is set for state 0 and 1, matching
// the adjacent byte write performed alongside markers ']' and '\\'.
type Paul2013TokenBoundaryMarker struct {
	Marker            byte
	PreviousStateFlag bool
}

// Paul2013TokenBoundaryStateResult retains the intermediate arrays and
// markers produced by the model-state portion of FUN_100130e0.
type Paul2013TokenBoundaryStateResult struct {
	Values  []int32
	States  []int32
	Markers []Paul2013TokenBoundaryMarker
}

// Paul2013TokenBoundaryArenaInputs contains fields read directly from the
// model-state record stream by FUN_100130e0. InitialValues and InitialStates
// are separate caller-owned arrays and remain outside this projection.
type Paul2013TokenBoundaryArenaInputs struct {
	InitialMarkers []byte
	RowStateBytes  []byte
	TokenCodes     []int16
	DurationValues []byte
}

// ReadPaul2013TokenBoundaryArenaInputs reads the record-backed inputs used by
// FUN_100130e0. The signed code at record +0x3ac indexes the 13-byte marker
// table; the transition pass also reads that code and the byte at +0x3b0.
// Duration-limit bytes are read from +0x95. The arena uses a signed count at
// +2 and 0x3c0-byte records beginning at +0x64c.
func ReadPaul2013TokenBoundaryArenaInputs(modelStateArena []byte) (Paul2013TokenBoundaryArenaInputs, error) {
	if len(modelStateArena) < paul2013TokenMarkerArenaCount+2 {
		return Paul2013TokenBoundaryArenaInputs{}, fmt.Errorf("model-state arena has %d bytes, need count at +%#x", len(modelStateArena), paul2013TokenMarkerArenaCount)
	}
	count := int(int16(binary.LittleEndian.Uint16(modelStateArena[paul2013TokenMarkerArenaCount : paul2013TokenMarkerArenaCount+2])))
	if count < 0 {
		return Paul2013TokenBoundaryArenaInputs{}, fmt.Errorf("model-state arena has negative record count %d", count)
	}
	if count > paul2013RecordGroupCapacity {
		return Paul2013TokenBoundaryArenaInputs{}, fmt.Errorf("model-state arena record count %d exceeds native capacity %d", count, paul2013RecordGroupCapacity)
	}
	if paul2013TokenMarkerArenaRecords > len(modelStateArena) || count > (len(modelStateArena)-paul2013TokenMarkerArenaRecords)/paul2013TokenMarkerArenaStride {
		return Paul2013TokenBoundaryArenaInputs{}, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(modelStateArena), count, paul2013TokenMarkerArenaRecords)
	}
	inputs := Paul2013TokenBoundaryArenaInputs{
		InitialMarkers: make([]byte, count),
		RowStateBytes:  make([]byte, count),
		TokenCodes:     make([]int16, count),
		DurationValues: make([]byte, count),
	}
	for index := 0; index < count; index++ {
		start := paul2013TokenMarkerArenaRecords + index*paul2013TokenMarkerArenaStride
		record := modelStateArena[start : start+paul2013TokenMarkerArenaStride]
		code := int16(binary.LittleEndian.Uint16(record[paul2013TokenMarkerCodeOffset : paul2013TokenMarkerCodeOffset+2]))
		marker, err := LookupPaul2013InitialTokenMarker(code)
		if err != nil {
			return Paul2013TokenBoundaryArenaInputs{}, fmt.Errorf("model-state record %d: %w", index, err)
		}
		inputs.InitialMarkers[index] = marker
		inputs.RowStateBytes[index] = record[paul2013TokenMarkerStateOffset]
		inputs.TokenCodes[index] = code
		inputs.DurationValues[index] = record[paul2013TokenMarkerDurationByte]
	}
	return inputs, nil
}

// BuildPaul2013TokenBoundaryMarkersFromModelStateArena composes the record
// field projection with the native initial-marker table and adjacent-record
// transition, in the order observed in FUN_100130e0. State/value array
// updates and their later marker dispatch remain separate because their
// producers are not established by the current static evidence.
func BuildPaul2013TokenBoundaryMarkersFromModelStateArena(modelStateArena []byte) ([]Paul2013TokenBoundaryMarker, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(modelStateArena)
	if err != nil {
		return nil, err
	}
	markers, err := BuildPaul2013InitialTokenMarkers(inputs.TokenCodes)
	if err != nil {
		return nil, err
	}
	result := make([]Paul2013TokenBoundaryMarker, len(markers))
	for index, marker := range markers {
		result[index].Marker = marker
	}
	return BuildPaul2013TokenBoundaryMarkerTransitions(result, inputs.RowStateBytes, inputs.TokenCodes)
}

// LookupPaul2013InitialTokenMarker applies the byte table at DLL address
// 0x10079a30, indexed by the signed token-state code used in FUN_100130e0.
func LookupPaul2013InitialTokenMarker(stateCode int16) (byte, error) {
	if stateCode < 0 || int(stateCode) >= len(paul2013InitialMarkerByState) {
		return 0, fmt.Errorf("Paul 2013 initial marker state %d is outside the observed table", stateCode)
	}
	return paul2013InitialMarkerByState[stateCode], nil
}

// BuildPaul2013InitialTokenMarkers resolves explicit token-state codes through
// the native marker table. Production of the signed state codes remains
// caller-provided.
func BuildPaul2013InitialTokenMarkers(stateCodes []int16) ([]byte, error) {
	markers := make([]byte, len(stateCodes))
	for index, stateCode := range stateCodes {
		marker, err := LookupPaul2013InitialTokenMarker(stateCode)
		if err != nil {
			return nil, fmt.Errorf("token %d initial marker: %w", index, err)
		}
		markers[index] = marker
	}
	return markers, nil
}

// BuildPaul2013TokenBoundaryStateResultWithTransitions composes explicit
// state updates and marker dispatch followed by adjacent-token transitions.
// This post-dispatch convenience order is retained for existing callers;
// BuildPaul2013TokenBoundaryMarkersFromModelStateArena exposes the earlier
// record-backed transition stage in FUN_100130e0's observed order.
func BuildPaul2013TokenBoundaryStateResultWithTransitions(
	initialMarkers []byte,
	initialValues []int32,
	initialStates []int32,
	perTokenOutputs []int32,
	rowStateBytes []byte,
	tokenCodes []int16,
) (Paul2013TokenBoundaryStateResult, error) {
	result, err := BuildPaul2013TokenBoundaryStateResult(
		initialMarkers, initialValues, initialStates, perTokenOutputs,
	)
	if err != nil {
		return Paul2013TokenBoundaryStateResult{}, err
	}
	result.Markers, err = BuildPaul2013TokenBoundaryMarkerTransitions(result.Markers, rowStateBytes, tokenCodes)
	if err != nil {
		return Paul2013TokenBoundaryStateResult{}, err
	}
	return result, nil
}

// BuildPaul2013TokenBoundaryStateResult ports the state-array updates in
// FUN_100130e0 before its marker dispatch. initialValues and initialStates
// retain values initialized elsewhere; perTokenOutputs supplies the values
// read at 0x94-byte strides. A result of -2 writes 100, a nonnegative
// result is copied, and other negative results leave the initialized value in
// place. For token indexes 1..N-1, a nonnegative resulting value sets the
// preceding token's state to 2 (0x1001324e reads state-pointer minus 0x31c).
// The producer of all three input arrays remains outside this helper.
func BuildPaul2013TokenBoundaryStateResult(
	initialMarkers []byte,
	initialValues []int32,
	initialStates []int32,
	perTokenOutputs []int32,
) (Paul2013TokenBoundaryStateResult, error) {
	tokenCount := len(initialMarkers)
	if len(initialValues) != tokenCount || len(initialStates) != tokenCount {
		return Paul2013TokenBoundaryStateResult{}, fmt.Errorf(
			"received %d marker values and %d states for %d tokens",
			len(initialValues), len(initialStates), tokenCount,
		)
	}
	wantOutputs := tokenCount - 1
	if tokenCount == 0 {
		wantOutputs = 0
	}
	if len(perTokenOutputs) != wantOutputs {
		return Paul2013TokenBoundaryStateResult{}, fmt.Errorf(
			"received %d per-token state outputs for %d tokens, want %d",
			len(perTokenOutputs), tokenCount, wantOutputs,
		)
	}
	values := append([]int32(nil), initialValues...)
	states := append([]int32(nil), initialStates...)
	for index, output := range perTokenOutputs {
		switch {
		case output == -2:
			values[index] = 100
		case output >= 0:
			values[index] = output
		}
	}
	for index := 1; index < tokenCount; index++ {
		if values[index] >= 0 {
			states[index-1] = 2
		}
	}
	markers, err := BuildPaul2013TokenBoundaryMarkers(initialMarkers, states)
	if err != nil {
		return Paul2013TokenBoundaryStateResult{}, err
	}
	return Paul2013TokenBoundaryStateResult{Values: values, States: states, Markers: markers}, nil
}

// BuildPaul2013TokenBoundaryMarkers ports FUN_100130e0's state dispatch after
// FUN_10012f00 has produced per-token states. initialMarkers are the already
// initialized marker bytes; states are the model-state values copied into the
// dispatch array. State production and initial marker lookup remain caller
// inputs. For each non-final token, states 0..3 select ']', '\\', '[', and
// 'Z'. The final token is forced to '[' for states 0..2 and otherwise keeps
// its initial marker.
func BuildPaul2013TokenBoundaryMarkers(
	initialMarkers []byte,
	states []int32,
) ([]Paul2013TokenBoundaryMarker, error) {
	if len(initialMarkers) != len(states) {
		return nil, fmt.Errorf("received %d initial markers for %d token states", len(initialMarkers), len(states))
	}
	result := make([]Paul2013TokenBoundaryMarker, len(states))
	for index, marker := range initialMarkers {
		result[index].Marker = marker
	}
	if len(states) == 0 {
		return result, nil
	}
	for index, state := range states[:len(states)-1] {
		switch state {
		case 0:
			result[index] = Paul2013TokenBoundaryMarker{Marker: ']', PreviousStateFlag: true}
		case 1:
			result[index] = Paul2013TokenBoundaryMarker{Marker: '\\', PreviousStateFlag: true}
		case 2:
			result[index].Marker = '['
		case 3:
			result[index].Marker = 'Z'
		}
	}
	finalIndex := len(states) - 1
	if states[finalIndex] >= 0 && states[finalIndex] < 3 {
		result[finalIndex].Marker = '['
	}
	return result, nil
}

// BuildPaul2013TokenBoundaryMarkerTransitions ports the adjacent-record
// transition in FUN_100130e0. stateBytes and tokenCodes are the opaque values
// read from current and next rows; prior markers and flags are copied before
// the transition is applied.
func BuildPaul2013TokenBoundaryMarkerTransitions(
	markers []Paul2013TokenBoundaryMarker,
	stateBytes []byte,
	tokenCodes []int16,
) ([]Paul2013TokenBoundaryMarker, error) {
	if len(stateBytes) != len(markers) || len(tokenCodes) != len(markers) {
		return nil, fmt.Errorf("received %d state bytes and %d token codes for %d markers", len(stateBytes), len(tokenCodes), len(markers))
	}
	result := append([]Paul2013TokenBoundaryMarker(nil), markers...)
	for index := 0; index+1 < len(result); index++ {
		marker := result[index].Marker
		ordinaryDelimiter := marker == ']' || marker == '\\'
		if stateBytes[index] == 8 || stateBytes[index+1] != 8 || !ordinaryDelimiter {
			if tokenCodes[index] == 12 {
				result[index].PreviousStateFlag = true
			}
			continue
		}
		result[index].Marker = '\\'
		result[index].PreviousStateFlag = true
	}
	return result, nil
}

// SplitPaul2013TokenPhoneBlocks ports FUN_10012c70's marker scan for one
// already-normalized token. Every marker other than ASCII '0' closes a block;
// the final phone closes a block regardless of its marker. The source-text
// producer of per-phone markers is not part of this helper.
func SplitPaul2013TokenPhoneBlocks(markers []byte) ([]Paul2013PhoneBlock, error) {
	if len(markers) > 65 {
		return nil, fmt.Errorf("Paul 2013 token has %d phones; block splitter supports at most 65", len(markers))
	}
	blocks := make([]Paul2013PhoneBlock, 0, len(markers))
	start := 0
	for phoneIndex, marker := range markers {
		if marker == '0' && phoneIndex != len(markers)-1 {
			continue
		}
		blocks = append(blocks, Paul2013PhoneBlock{
			Start:          start,
			End:            phoneIndex + 1,
			TerminalMarker: marker,
		})
		start = phoneIndex + 1
	}
	return blocks, nil
}

// ApplyPaul2013PhoneBlockDurationLimit ports FUN_100130e0's cumulative
// boundary insertion over the explicit per-phone byte values. Each value is
// doubled before accumulation; a total strictly greater than 1000 inserts
// '[' after the preceding phone. Existing '[', 'Z', '^', and '`' markers
// reset the running total. The values' text/model producer remains external.
func ApplyPaul2013PhoneBlockDurationLimit(
	markers []byte,
	values []byte,
) ([]byte, error) {
	if len(markers) != len(values) {
		return nil, fmt.Errorf("received %d phone markers for %d duration values", len(markers), len(values))
	}
	if len(markers) > 65 {
		return nil, fmt.Errorf("Paul 2013 token has %d phones; duration limiter supports at most 65", len(markers))
	}
	result := append([]byte(nil), markers...)
	accumulated := 0
	for index, value := range values {
		accumulated += int(value) * 2
		if accumulated > paul2013PhoneBlockDurationLimit {
			if index == 0 {
				return nil, fmt.Errorf("first phone value %d exceeds the duration limit without a preceding boundary", value)
			}
			result[index-1] = '['
			accumulated = int(value) * 2
		}
		switch result[index] {
		case '[', 'Z', '^', '`':
			accumulated = 0
		}
	}
	return result, nil
}

// ApplyPaul2013TokenBoundaryDurationLimitFromModelStateArena composes the
// arena's per-record +0x95 bytes with the cumulative 1,000-unit marker pass.
// The caller supplies the markers after the preceding marker stages; this
// helper does not infer the still-unresolved state/value arrays.
func ApplyPaul2013TokenBoundaryDurationLimitFromModelStateArena(
	modelStateArena []byte,
	markers []byte,
) ([]byte, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(modelStateArena)
	if err != nil {
		return nil, err
	}
	if len(markers) != len(inputs.DurationValues) {
		return nil, fmt.Errorf("received %d token markers for %d model-state records", len(markers), len(inputs.DurationValues))
	}
	return applyPaul2013RecordBoundaryDurationLimit(markers, inputs.DurationValues)
}

// BuildPaul2013MarkerTreeGroups ports FUN_10012df0's group construction and
// per-group byte accumulation. Markers '\\' and ']' keep the adjacent phone
// in the same group; other markers end the current group. Values are the
// caller-supplied bytes read from each phone row at +0x94.
func BuildPaul2013MarkerTreeGroups(
	markers []byte,
	durationByteValues []byte,
) ([]Paul2013MarkerTreeGroup, error) {
	if len(markers) != len(durationByteValues) {
		return nil, fmt.Errorf("received %d phone markers for %d duration byte values", len(markers), len(durationByteValues))
	}
	if len(markers) > 65 {
		return nil, fmt.Errorf("Paul 2013 token has %d phones; marker grouping supports at most 65", len(markers))
	}
	if len(markers) == 0 {
		return nil, nil
	}
	groups := make([]Paul2013MarkerTreeGroup, 0, len(markers))
	start := 0
	var sum uint16
	for index, value := range durationByteValues {
		sum += uint16(value)
		if index == len(markers)-1 || markers[index] != '\\' && markers[index] != ']' {
			groups = append(groups, Paul2013MarkerTreeGroup{
				PhoneStart: start, PhoneEnd: index + 1, DurationByteSum: sum,
			})
			start = index + 1
			sum = 0
		}
	}
	return groups, nil
}
