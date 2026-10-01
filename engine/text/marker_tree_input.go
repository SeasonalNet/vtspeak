package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/internal/paul2013tables"
	"vtspeak/engine/tree3"
)

const paul2013MarkerTreeFeatureCount = 15

// Paul2013MarkerTreeInputState contains the running values observed in
// FUN_10012f00 immediately before it visits CurrentPhoneRow. The two string
// bytes are resolved by the caller because the native record stores pointers,
// not inline strings, at the corresponding locations.
type Paul2013MarkerTreeInputState struct {
	TokenPhoneCount        int16
	PhoneOrdinal           int16
	PhonesSinceLastSplit   int32
	CumulativePhoneValue   int32
	NextTokenCumulative    int16
	ValueSinceLastSplit    int32
	CurrentStringFirstByte byte
	PreviousSelectedByte   byte
}

// Paul2013MarkerTreeInput is one eligible 15-short vector for FUN_10001670.
// The updated counters include CurrentPhoneRow's +0x94 byte.
type Paul2013MarkerTreeInput struct {
	Features             [paul2013MarkerTreeFeatureCount]int16
	Eligible             bool
	PhoneOrdinal         int16
	PhonesSinceLastSplit int32
	CumulativePhoneValue int32
	ValueSinceLastSplit  int32
}

// Paul2013MarkerTreePhoneVisit supplies one native reverse-scan record and
// the two bytes reached through its pointer-backed strings.
type Paul2013MarkerTreePhoneVisit struct {
	PreviousPhoneRow       []byte
	CurrentPhoneRow        []byte
	CurrentStringFirstByte byte
	PreviousSelectedByte   byte
}

// Paul2013MarkerTreeTokenResult contains one token's native-order scan.
// Decision rows are copies of the preceding records and preserve their bytes
// except for the terminal marker write when a tree result exceeds 500.
type Paul2013MarkerTreeTokenResult struct {
	Inputs     []Paul2013MarkerTreeInput
	TreeValues []int16
	Leaves     []int
	Decisions  []Paul2013MarkerTreeDecision
	EndState   Paul2013MarkerTreeInputState
}

// BuildPaul2013MarkerTreeInput projects the vector assembled by
// FUN_10012f00. previousPhoneRow and currentPhoneRow are adjacent native
// 0x3c0-byte records. State counters are the values before the current row is
// added; token-boundary counts and native string-pointer targets remain
// explicit because their owning arena layout is not part of the portable row
// representation.
func BuildPaul2013MarkerTreeInput(
	previousPhoneRow []byte,
	currentPhoneRow []byte,
	state Paul2013MarkerTreeInputState,
) (Paul2013MarkerTreeInput, error) {
	if len(previousPhoneRow) < paul2013MarkerTreePhoneRowStride {
		return Paul2013MarkerTreeInput{}, fmt.Errorf("previous phone row has %d bytes, need %d", len(previousPhoneRow), paul2013MarkerTreePhoneRowStride)
	}
	if len(currentPhoneRow) < paul2013MarkerTreePhoneRowStride {
		return Paul2013MarkerTreeInput{}, fmt.Errorf("current phone row has %d bytes, need %d", len(currentPhoneRow), paul2013MarkerTreePhoneRowStride)
	}
	if state.TokenPhoneCount < 0 || state.PhoneOrdinal < 0 || state.PhonesSinceLastSplit < 0 {
		return Paul2013MarkerTreeInput{}, fmt.Errorf("negative marker-tree phone count: token=%d ordinal=%d since-split=%d", state.TokenPhoneCount, state.PhoneOrdinal, state.PhonesSinceLastSplit)
	}
	previousWord := int32(binary.LittleEndian.Uint32(previousPhoneRow[4:8]))
	currentWord := int32(binary.LittleEndian.Uint32(currentPhoneRow[0:4]))
	currentRowWord := int32(binary.LittleEndian.Uint32(currentPhoneRow[4:8]))
	phoneOrdinal := state.PhoneOrdinal + 1
	phonesSinceSplit := state.PhonesSinceLastSplit + 1
	phoneValue := state.CumulativePhoneValue + int32(currentPhoneRow[0x94])
	splitValue := state.ValueSinceLastSplit + int32(currentPhoneRow[0x94])
	input := Paul2013MarkerTreeInput{
		PhoneOrdinal:         phoneOrdinal,
		PhonesSinceLastSplit: phonesSinceSplit,
		CumulativePhoneValue: phoneValue,
		ValueSinceLastSplit:  splitValue,
	}
	if previousWord == currentRowWord || previousWord+1 == currentWord {
		return input, nil
	}
	currentStringClass, err := paul2013MarkerTreeClassByte(state.CurrentStringFirstByte)
	if err != nil {
		return Paul2013MarkerTreeInput{}, fmt.Errorf("current marker-tree string byte: %w", err)
	}
	previousStringClass, err := paul2013MarkerTreeClassByte(state.PreviousSelectedByte)
	if err != nil {
		return Paul2013MarkerTreeInput{}, fmt.Errorf("previous marker-tree string byte: %w", err)
	}
	input.Features = [paul2013MarkerTreeFeatureCount]int16{
		state.TokenPhoneCount - phoneOrdinal,
		phoneOrdinal,
		state.TokenPhoneCount,
		state.NextTokenCumulative - int16(phoneValue),
		int16(phoneValue),
		state.NextTokenCumulative,
		int16(phonesSinceSplit),
		int16(splitValue),
		int16(currentPhoneRow[0x3b3]),
		int16(previousPhoneRow[0x3b2]),
		int16(currentPhoneRow[0x3b2]),
		int16(binary.LittleEndian.Uint16(previousPhoneRow[0x3ae:0x3b0])),
		boolShort(previousWord == currentRowWord),
		int16(currentStringClass),
		int16(previousStringClass),
	}
	input.Eligible = true
	return input, nil
}

// EvaluatePaul2013MarkerTreeInput runs the supported scalar tree lookup and
// applies FUN_10012f00's result branch to the preceding row. Tree selection and
// the two pointer-backed class bytes are caller-provided.
func EvaluatePaul2013MarkerTreeInput(
	tree *tree3.Tree,
	previousPhoneRow []byte,
	currentPhoneRow []byte,
	state Paul2013MarkerTreeInputState,
) (Paul2013MarkerTreeInput, int16, int, Paul2013MarkerTreeDecision, error) {
	input, err := BuildPaul2013MarkerTreeInput(previousPhoneRow, currentPhoneRow, state)
	if err != nil {
		return Paul2013MarkerTreeInput{}, 0, 0, Paul2013MarkerTreeDecision{}, err
	}
	if !input.Eligible {
		decision, decisionErr := ApplyPaul2013MarkerTreeDecision(previousPhoneRow, 0)
		return input, 0, 0, decision, decisionErr
	}
	if tree == nil {
		return Paul2013MarkerTreeInput{}, 0, 0, Paul2013MarkerTreeDecision{}, fmt.Errorf("Paul 2013 marker tree is nil")
	}
	if tree.OutputWidth != 1 {
		return Paul2013MarkerTreeInput{}, 0, 0, Paul2013MarkerTreeDecision{}, fmt.Errorf("Paul 2013 marker tree has output width %d, want 1", tree.OutputWidth)
	}
	leaf, output, err := tree.Evaluate(input.Features[:])
	if err != nil {
		return Paul2013MarkerTreeInput{}, 0, 0, Paul2013MarkerTreeDecision{}, fmt.Errorf("evaluate Paul 2013 marker tree: %w", err)
	}
	if len(output) != 1 {
		return Paul2013MarkerTreeInput{}, 0, 0, Paul2013MarkerTreeDecision{}, fmt.Errorf("Paul 2013 marker tree returned %d outputs, want 1", len(output))
	}
	decision, err := ApplyPaul2013MarkerTreeDecision(previousPhoneRow, output[0])
	if err != nil {
		return Paul2013MarkerTreeInput{}, 0, 0, Paul2013MarkerTreeDecision{}, err
	}
	return input, output[0], leaf, decision, nil
}

// RunPaul2013MarkerTreeToken composes the reverse phone scan in
// FUN_10012f00 for one token. visits must be ordered as the native loop visits
// them; token grouping, previous-row associations, pointer resolution, and
// selected-tree choice remain explicit.
func RunPaul2013MarkerTreeToken(
	tree *tree3.Tree,
	tokenPhoneCount int16,
	nextTokenCumulative int16,
	visits []Paul2013MarkerTreePhoneVisit,
) (Paul2013MarkerTreeTokenResult, error) {
	if tokenPhoneCount < 0 {
		return Paul2013MarkerTreeTokenResult{}, fmt.Errorf("negative marker-tree token phone count %d", tokenPhoneCount)
	}
	result := Paul2013MarkerTreeTokenResult{
		Inputs:     make([]Paul2013MarkerTreeInput, 0, len(visits)),
		TreeValues: make([]int16, 0, len(visits)),
		Leaves:     make([]int, 0, len(visits)),
		Decisions:  make([]Paul2013MarkerTreeDecision, 0, len(visits)),
		EndState: Paul2013MarkerTreeInputState{
			TokenPhoneCount:     tokenPhoneCount,
			NextTokenCumulative: nextTokenCumulative,
		},
	}
	for index, visit := range visits {
		state := result.EndState
		state.CurrentStringFirstByte = visit.CurrentStringFirstByte
		state.PreviousSelectedByte = visit.PreviousSelectedByte
		input, treeValue, leaf, decision, err := EvaluatePaul2013MarkerTreeInput(
			tree, visit.PreviousPhoneRow, visit.CurrentPhoneRow, state,
		)
		if err != nil {
			return Paul2013MarkerTreeTokenResult{}, fmt.Errorf("marker-tree token phone visit %d: %w", index, err)
		}
		result.Inputs = append(result.Inputs, input)
		result.TreeValues = append(result.TreeValues, treeValue)
		result.Leaves = append(result.Leaves, leaf)
		result.Decisions = append(result.Decisions, decision)
		result.EndState.PhoneOrdinal = input.PhoneOrdinal
		result.EndState.PhonesSinceLastSplit = input.PhonesSinceLastSplit
		result.EndState.CumulativePhoneValue = input.CumulativePhoneValue
		result.EndState.ValueSinceLastSplit = input.ValueSinceLastSplit
		if decision.InsertedMarker {
			result.EndState.PhonesSinceLastSplit = 0
			result.EndState.ValueSinceLastSplit = 0
		}
	}
	return result, nil
}

func boolShort(value bool) int16 {
	if value {
		return 1
	}
	return 0
}

func paul2013MarkerTreeClassByte(value byte) (byte, error) {
	if value >= 0x80 {
		return 0, fmt.Errorf("high-bit byte 0x%02x follows the native signed-char index outside the recovered 256-byte table", value)
	}
	return paul2013tables.ByteTable1007B9E0(value), nil
}
