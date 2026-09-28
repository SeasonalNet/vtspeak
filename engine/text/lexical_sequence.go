package text

import (
	"errors"
	"fmt"
)

// LexicalTokenSpan maps one explicitly selected token pronunciation into a
// flattened phone sequence. The offsets are half-open phone indexes.
type LexicalTokenSpan struct {
	SourceSurface      string
	Surface            string
	SeparatorBefore    string
	SeparatorAfter     string
	DictionaryMetadata [4]bool
	AlternativeIndex   int
	PhoneStart         int
	PhoneEnd           int
}

// LexicalPhoneSequence contains pronunciations in source order while
// retaining each token boundary and its original separators.
type LexicalPhoneSequence struct {
	Phones  []CMUPhone
	Symbols []byte
	Tokens  []LexicalTokenSpan
}

// LexicalTokenNeighborhoods pairs one token's source metadata with its slice
// of the utterance-wide phone features.
type LexicalTokenNeighborhoods struct {
	Token            LexicalTokenSpan
	PhoneFeatures    []Paul2013PhoneNeighborhood
	PhoneGroups      []Paul2013PhoneGroup
	DurationContexts []Paul2013PhoneDurationContext
}

// Paul2013PhoneDurationContext contains one phone's projection of its legacy
// context row. GroupCount is the number of rows in the current token block,
// while RecordByte2 is the per-token group position stored in the row.
type Paul2013PhoneDurationContext struct {
	GroupStart        int
	GroupEnd          int
	GroupIndex        int
	GroupCount        int
	RowStart          int
	RowPhoneCount     int
	NucleusStress     uint8
	RecordByte1       int8
	RecordByte2       uint8
	TreePositionState uint8
	AuxiliaryByte     int8
}

// Paul2013TokenDurationBoundaries supplies the explicit previous/next
// identity ordinals used at one token span's edges.
type Paul2013TokenDurationBoundaries struct {
	PreviousIdentityOrdinal int16
	NextIdentityOrdinal     int16
}

const paul2013OrdinaryBoundaryOrdinal int16 = 40

// Paul2013TokenDurationInputs contains the observed nine-short duration input
// for each phone in one token span.
type Paul2013TokenDurationInputs struct {
	Token  LexicalTokenSpan
	Inputs [][9]int16
}

// SelectLexicalPronunciations flattens one explicitly selected alternative
// per lexical token. It does not rank alternatives or merge token boundaries.
func SelectLexicalPronunciations(tokens []LexicalToken, alternativeIndexes []int) (LexicalPhoneSequence, error) {
	if len(tokens) == 0 {
		return LexicalPhoneSequence{}, errors.New("lexical phone sequence has no tokens")
	}
	if len(alternativeIndexes) != len(tokens) {
		return LexicalPhoneSequence{}, fmt.Errorf("received %d pronunciation choices for %d tokens", len(alternativeIndexes), len(tokens))
	}

	sequence := LexicalPhoneSequence{Tokens: make([]LexicalTokenSpan, 0, len(tokens))}
	for tokenIndex, token := range tokens {
		alternativeIndex := alternativeIndexes[tokenIndex]
		if alternativeIndex < 0 || alternativeIndex >= len(token.Alternatives) {
			return LexicalPhoneSequence{}, fmt.Errorf("token %d (%q) has no pronunciation alternative %d", tokenIndex, token.SourceSurface, alternativeIndex)
		}
		alternative := token.Alternatives[alternativeIndex]
		start := len(sequence.Phones)
		sequence.Phones = append(sequence.Phones, alternative.Phones...)
		sequence.Symbols = append(sequence.Symbols, alternative.Symbols...)
		sequence.Tokens = append(sequence.Tokens, LexicalTokenSpan{
			SourceSurface:      token.SourceSurface,
			Surface:            token.Surface,
			SeparatorBefore:    token.SeparatorBefore,
			SeparatorAfter:     token.SeparatorAfter,
			DictionaryMetadata: token.DictionaryMetadata,
			AlternativeIndex:   alternativeIndex,
			PhoneStart:         start,
			PhoneEnd:           len(sequence.Phones),
		})
	}
	return sequence, nil
}

// BuildPaul2013TokenPhoneNeighborhoods derives identity/stress neighborhoods
// over the complete utterance, then groups the results by token. Phones on
// either side of a word boundary are adjacent in the phone sequence; preserved
// separators remain available to later prosody and boundary processing.
func BuildPaul2013TokenPhoneNeighborhoods(sequence LexicalPhoneSequence) ([]LexicalTokenNeighborhoods, error) {
	if err := validateLexicalPhoneSpans(sequence); err != nil {
		return nil, err
	}
	utteranceFeatures, err := BuildPaul2013PhoneNeighborhoods(sequence.Phones)
	if err != nil {
		return nil, fmt.Errorf("build utterance phone neighborhoods: %w", err)
	}
	result := make([]LexicalTokenNeighborhoods, len(sequence.Tokens))
	for tokenIndex, token := range sequence.Tokens {
		groups, err := BuildPaul2013PhoneGroups(sequence.Phones[token.PhoneStart:token.PhoneEnd])
		if err != nil {
			return nil, fmt.Errorf("build phone groups for token %d (%q): %w", tokenIndex, token.SourceSurface, err)
		}
		durationContexts, err := buildPaul2013PhoneDurationContexts(sequence.Phones[token.PhoneStart:token.PhoneEnd], groups)
		if err != nil {
			return nil, fmt.Errorf("build duration contexts for token %d (%q): %w", tokenIndex, token.SourceSurface, err)
		}
		result[tokenIndex] = LexicalTokenNeighborhoods{
			Token:            token,
			PhoneFeatures:    append([]Paul2013PhoneNeighborhood(nil), utteranceFeatures[token.PhoneStart:token.PhoneEnd]...),
			PhoneGroups:      groups,
			DurationContexts: durationContexts,
		}
	}
	if len(sequence.Phones) > 0 {
		result[0].DurationContexts[0].TreePositionState = 1
		if len(sequence.Phones) > 1 {
			lastTokenIndex := len(result) - 1
			lastPhoneIndex := len(result[lastTokenIndex].DurationContexts) - 1
			// FUN_100135d0 uses marker-dependent terminal states. The independent
			// text path currently emits the standard 'Z' utterance terminator.
			result[lastTokenIndex].DurationContexts[lastPhoneIndex].TreePositionState = 3
		}
	}
	for tokenIndex := range result {
		for phoneIndex := range result[tokenIndex].DurationContexts {
			if result[tokenIndex].DurationContexts[phoneIndex].TreePositionState == 0 {
				result[tokenIndex].DurationContexts[phoneIndex].TreePositionState = 2
			}
		}
	}
	return result, nil
}

func buildPaul2013PhoneDurationContexts(phones []CMUPhone, groups []Paul2013PhoneGroup) ([]Paul2013PhoneDurationContext, error) {
	contexts := make([]Paul2013PhoneDurationContext, len(phones))
	if len(groups) == 0 {
		// FUN_10013c00's no-vowel branch emits one fallback row with row bytes
		// (0, 4, 1) and auxiliary label 1 for every phone.
		for phoneIndex := range contexts {
			contexts[phoneIndex] = Paul2013PhoneDurationContext{
				GroupStart: 0, GroupEnd: len(phones), GroupIndex: 0,
				GroupCount: 1, RowStart: 0, RowPhoneCount: len(phones),
				RecordByte1: 4, RecordByte2: 1, AuxiliaryByte: 1,
			}
		}
		return contexts, nil
	}
	rowOffset := 0
	for groupIndex, group := range groups {
		onsetCount := group.Nucleus - group.OnsetStart
		codaCount := group.End - group.Nucleus - 1
		rowPhoneCount := onsetCount + 1 + codaCount
		if onsetCount < 0 || codaCount < 0 || rowPhoneCount < 1 || rowOffset+rowPhoneCount > len(contexts) {
			return nil, fmt.Errorf("group %d has invalid duration row range", groupIndex)
		}
		rowByte1 := int8(0)
		if onsetCount > 0 {
			rowByte1 |= 1
		}
		if codaCount > 0 {
			rowByte1 |= 2
		}
		for rowPhoneIndex := 0; rowPhoneIndex < rowPhoneCount; rowPhoneIndex++ {
			phoneIndex := rowOffset + rowPhoneIndex
			auxiliary := int8(2)
			if rowPhoneIndex < onsetCount {
				auxiliary = 1
			} else if rowPhoneIndex > onsetCount {
				auxiliary = 3
			}
			contexts[phoneIndex] = Paul2013PhoneDurationContext{
				GroupStart: group.Start, GroupEnd: group.End,
				GroupIndex: groupIndex, GroupCount: len(groups),
				RowStart: rowOffset, RowPhoneCount: rowPhoneCount,
				NucleusStress: phones[group.Nucleus].Stress,
				RecordByte1:   rowByte1, RecordByte2: group.PositionState,
				AuxiliaryByte: auxiliary,
			}
		}
		rowOffset += rowPhoneCount
	}
	if rowOffset != len(phones) {
		return nil, fmt.Errorf("duration groups cover %d of %d phones", rowOffset, len(phones))
	}
	return contexts, nil
}

// BuildPaul2013TokenDurationInputs is a low-level trace helper. Callers supply
// raw row metadata and one pair of boundary ordinals per token; use
// BuildPaul2013TokenDurationInputsFromText for engine-derived row state.
func BuildPaul2013TokenDurationInputs(
	sequence LexicalPhoneSequence,
	metadata []Paul2013DurationTreeMetadata,
	boundaries []Paul2013TokenDurationBoundaries,
) ([]Paul2013TokenDurationInputs, error) {
	if err := validateLexicalPhoneSpans(sequence); err != nil {
		return nil, err
	}
	if len(metadata) != len(sequence.Phones) {
		return nil, fmt.Errorf("received %d duration metadata records for %d phones", len(metadata), len(sequence.Phones))
	}
	if len(boundaries) != len(sequence.Tokens) {
		return nil, fmt.Errorf("received %d duration boundary pairs for %d tokens", len(boundaries), len(sequence.Tokens))
	}
	result := make([]Paul2013TokenDurationInputs, len(sequence.Tokens))
	for tokenIndex, token := range sequence.Tokens {
		inputs := make([][9]int16, token.PhoneEnd-token.PhoneStart)
		for phoneIndex := token.PhoneStart; phoneIndex < token.PhoneEnd; phoneIndex++ {
			previous := Paul2013DurationTreeNeighbor{
				BoundaryIdentityOrdinal: boundaries[tokenIndex].PreviousIdentityOrdinal,
			}
			if phoneIndex > token.PhoneStart {
				previous = Paul2013DurationTreeNeighbor{Phone: &sequence.Phones[phoneIndex-1]}
			}
			next := Paul2013DurationTreeNeighbor{
				BoundaryIdentityOrdinal: boundaries[tokenIndex].NextIdentityOrdinal,
			}
			if phoneIndex+1 < token.PhoneEnd {
				next = Paul2013DurationTreeNeighbor{Phone: &sequence.Phones[phoneIndex+1]}
			}
			input, err := BuildObservedPaul2013DurationTreeInput(
				previous, sequence.Phones[phoneIndex], next, metadata[phoneIndex],
			)
			if err != nil {
				return nil, fmt.Errorf("build duration input for token %d phone %d: %w", tokenIndex, phoneIndex-token.PhoneStart, err)
			}
			inputs[phoneIndex-token.PhoneStart] = input
		}
		result[tokenIndex] = Paul2013TokenDurationInputs{Token: token, Inputs: inputs}
	}
	return result, nil
}

// BuildPaul2013TokenDurationInputsFromText builds the nine-short duration
// vectors from a resolved phone sequence. It ports the FUN_10013c00 row bytes
// and auxiliary labels, uses adjacent phones across token boundaries, and
// applies the ordinary single-utterance edge sentinels from FUN_100135d0.
// This text path assumes the standard terminal 'Z' marker; alternate marker
// transitions used for special pauses and phrase boundaries remain unsupported.
func BuildPaul2013TokenDurationInputsFromText(sequence LexicalPhoneSequence) ([]Paul2013TokenDurationInputs, error) {
	if err := validateLexicalPhoneSpans(sequence); err != nil {
		return nil, err
	}
	tokenNeighborhoods, err := BuildPaul2013TokenPhoneNeighborhoods(sequence)
	if err != nil {
		return nil, err
	}
	result := make([]Paul2013TokenDurationInputs, len(sequence.Tokens))
	for tokenIndex, neighborhood := range tokenNeighborhoods {
		token := neighborhood.Token
		inputs := make([][9]int16, token.PhoneEnd-token.PhoneStart)
		for phoneIndex := token.PhoneStart; phoneIndex < token.PhoneEnd; phoneIndex++ {
			previous := Paul2013DurationTreeNeighbor{BoundaryIdentityOrdinal: paul2013OrdinaryBoundaryOrdinal}
			if phoneIndex > 0 {
				previous = Paul2013DurationTreeNeighbor{Phone: &sequence.Phones[phoneIndex-1]}
			}
			next := Paul2013DurationTreeNeighbor{BoundaryIdentityOrdinal: paul2013OrdinaryBoundaryOrdinal}
			if phoneIndex+1 < len(sequence.Phones) {
				next = Paul2013DurationTreeNeighbor{Phone: &sequence.Phones[phoneIndex+1]}
			}
			durationContext := neighborhood.DurationContexts[phoneIndex-token.PhoneStart]
			metadata := Paul2013DurationTreeMetadata{
				RecordByte1:   durationContext.RecordByte1,
				RecordByte2:   int8(durationContext.RecordByte2),
				PositionState: durationContext.TreePositionState,
				ContextCount:  uint8(durationContext.GroupCount),
				AuxiliaryByte: durationContext.AuxiliaryByte,
			}
			input, err := BuildObservedPaul2013DurationTreeInput(previous, sequence.Phones[phoneIndex], next, metadata)
			if err != nil {
				return nil, fmt.Errorf("build duration input for token %d phone %d: %w", tokenIndex, phoneIndex-token.PhoneStart, err)
			}
			inputs[phoneIndex-token.PhoneStart] = input
		}
		result[tokenIndex] = Paul2013TokenDurationInputs{Token: token, Inputs: inputs}
	}
	return result, nil
}

func validateLexicalPhoneSpans(sequence LexicalPhoneSequence) error {
	if len(sequence.Tokens) == 0 {
		return errors.New("lexical phone sequence has no token spans")
	}
	expectedStart := 0
	for tokenIndex, token := range sequence.Tokens {
		if token.PhoneStart != expectedStart || token.PhoneEnd < token.PhoneStart || token.PhoneEnd > len(sequence.Phones) {
			return fmt.Errorf("token span %d is not a valid contiguous range of the phone sequence", tokenIndex)
		}
		expectedStart = token.PhoneEnd
	}
	if expectedStart != len(sequence.Phones) {
		return errors.New("token spans do not cover the complete phone sequence")
	}
	return nil
}
