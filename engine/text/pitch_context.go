package text

import (
	"fmt"
)

// Paul2013PitchInputCount is the number of values produced by FUN_10013a20.
// FUN_100138c0 appends the preceding scalar tree result as input 11 to the
// paired 12-value pitch tree.
const Paul2013PitchInputCount = 11

// Paul2013TokenPitchInputs contains FUN_10013a20's input row for each phone
// in one token.
type Paul2013TokenPitchInputs struct {
	Token  LexicalTokenSpan
	Inputs [][Paul2013PitchInputCount]int16
}

// BuildPaul2013TokenPitchInputsFromText ports FUN_10013a20 from the
// text-derived phone and group rows using the ordinary terminal-Z marker.
func BuildPaul2013TokenPitchInputsFromText(sequence LexicalPhoneSequence) ([]Paul2013TokenPitchInputs, error) {
	terminalMarkers := make([]byte, len(sequence.Tokens))
	if len(terminalMarkers) != 0 {
		terminalMarkers[len(terminalMarkers)-1] = 'Z'
	}
	return BuildPaul2013TokenPitchInputsWithMarkers(sequence, terminalMarkers)
}

// BuildPaul2013TokenPitchInputsWithMarkers evaluates the same context rows
// while applying caller-supplied final-token markers to the terminal position
// state consumed by FUN_10013a20. Marker-to-tree dispatch is handled by the
// duration package.
func BuildPaul2013TokenPitchInputsWithMarkers(
	sequence LexicalPhoneSequence,
	terminalMarkers []byte,
) ([]Paul2013TokenPitchInputs, error) {
	return BuildPaul2013TokenPitchInputsWithPhoneMarkers(sequence, nil, terminalMarkers)
}

// BuildPaul2013TokenPitchInputsWithPhoneMarkers composes explicit
// FUN_10012c70 block markers with pitch rows derived from the resulting phone
// groups. terminalMarkers independently select the last-token pitch tree.
func BuildPaul2013TokenPitchInputsWithPhoneMarkers(
	sequence LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
) ([]Paul2013TokenPitchInputs, error) {
	return buildPaul2013TokenPitchInputs(sequence, phoneMarkers, terminalMarkers, nil, false)
}

// BuildPaul2013TokenPitchInputsWithPositionStates composes explicit phone
// groups with caller-supplied per-phone position states. The same states are
// consumed by the duration-tree inputs and the pitch-tree position field.
func BuildPaul2013TokenPitchInputsWithPositionStates(
	sequence LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	positionStates []uint8,
) ([]Paul2013TokenPitchInputs, error) {
	return buildPaul2013TokenPitchInputs(sequence, phoneMarkers, terminalMarkers, positionStates, true)
}

func buildPaul2013TokenPitchInputs(
	sequence LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	positionStates []uint8,
	usePositionStates bool,
) ([]Paul2013TokenPitchInputs, error) {
	if err := validateLexicalPhoneSpans(sequence); err != nil {
		return nil, err
	}
	if len(terminalMarkers) != len(sequence.Tokens) {
		return nil, fmt.Errorf("received %d pitch terminal markers for %d tokens", len(terminalMarkers), len(sequence.Tokens))
	}
	var neighborhoods []LexicalTokenNeighborhoods
	var err error
	if usePositionStates {
		neighborhoods, err = BuildPaul2013TokenPhoneNeighborhoodsWithPositionStates(
			sequence, phoneMarkers, terminalMarkers, positionStates,
		)
	} else {
		if phoneMarkers == nil {
			neighborhoods, err = BuildPaul2013TokenPhoneNeighborhoods(sequence)
		} else {
			neighborhoods, err = BuildPaul2013TokenPhoneNeighborhoodsWithPhoneMarkers(sequence, phoneMarkers)
		}
		if err == nil {
			err = applyPaul2013TerminalMarkerStates(sequence, terminalMarkers, neighborhoods)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("build phone rows for pitch inputs: %w", err)
	}
	result := make([]Paul2013TokenPitchInputs, len(neighborhoods))
	for tokenIndex, neighborhood := range neighborhoods {
		phones := sequence.Phones[neighborhood.Token.PhoneStart:neighborhood.Token.PhoneEnd]
		if len(phones) == 0 {
			return nil, fmt.Errorf("token %d (%q) has no phones for pitch inputs", tokenIndex, neighborhood.Token.SourceSurface)
		}
		if len(phones) > 255 {
			return nil, fmt.Errorf("token %d (%q) has %d phones; pitch row stores an unsigned byte count", tokenIndex, neighborhood.Token.SourceSurface, len(phones))
		}
		if len(neighborhood.DurationContexts) != len(phones) {
			return nil, fmt.Errorf("token %d (%q) has %d pitch rows for %d phones", tokenIndex, neighborhood.Token.SourceSurface, len(neighborhood.DurationContexts), len(phones))
		}
		inputs := make([][Paul2013PitchInputCount]int16, len(phones))
		for phoneIndex, phone := range phones {
			if _, err := phone.TreeFeatures(); err != nil {
				return nil, fmt.Errorf("token %d phone %d: %w", tokenIndex, phoneIndex, err)
			}
			current := neighborhood.DurationContexts[phoneIndex]
			nucleusIndex := current.GroupEnd - 1
			for candidate := current.GroupStart; candidate < current.GroupEnd; candidate++ {
				if neighborhood.DurationContexts[candidate].AuxiliaryByte == 2 {
					nucleusIndex = candidate
					break
				}
			}
			nucleusFeatures, err := phones[nucleusIndex].TreeFeatures()
			if err != nil {
				return nil, fmt.Errorf("token %d group %d nucleus phone %d: %w", tokenIndex, current.GroupIndex, nucleusIndex, err)
			}

			row := &inputs[phoneIndex]
			pitchClass, ok := paul2013PitchPhoneClass[phone.Label]
			if !ok {
				return nil, fmt.Errorf("token %d phone %d (%q) has no recovered pitch class", tokenIndex, phoneIndex, phone.Label)
			}
			row[0] = pitchClass
			row[1] = nucleusFeatures.IdentityOrdinal
			row[2] = int16(current.NucleusStress)
			row[3] = 3
			row[4] = 3
			row[5] = int16(current.RecordByte2)
			row[6] = 0
			row[7] = 4
			row[8] = int16(current.TreePositionState)
			row[9] = int16(current.RecordByte1)
			row[10] = int16(len(phones))
			if phoneIndex > 0 {
				previous := neighborhood.DurationContexts[phoneIndex-1]
				row[3] = int16(previous.NucleusStress)
				row[6] = int16(previous.RecordByte2)
			}
			if phoneIndex+1 < len(phones) {
				next := neighborhood.DurationContexts[phoneIndex+1]
				row[4] = int16(next.NucleusStress)
				row[7] = int16(next.RecordByte2)
			}
			if row[0] < 0 || row[0] > 6 {
				return nil, fmt.Errorf("token %d phone %d (%q) has unsupported pitch class %d", tokenIndex, phoneIndex, phone.Label, row[0])
			}
		}
		result[tokenIndex] = Paul2013TokenPitchInputs{Token: neighborhood.Token, Inputs: inputs}
	}
	return result, nil
}

// paul2013PitchPhoneClass is the first byte at each 11-byte entry in
// DAT_10079880, indexed by the one-based phone identity from DAT_1007b6c0.
// It is zero for vowels and the observed consonant-class selector otherwise.
var paul2013PitchPhoneClass = map[string]int16{
	"AA": 0, "AE": 0, "AH": 0, "AO": 0, "AW": 0, "AY": 0,
	"B": 1, "CH": 3, "D": 1, "DH": 2, "EH": 0, "ER": 6,
	"EY": 0, "F": 2, "G": 1, "HH": 2, "IH": 0, "IY": 0,
	"JH": 3, "K": 1, "L": 5, "M": 4, "N": 4, "NG": 4,
	"OW": 0, "OY": 0, "P": 1, "R": 6, "S": 2, "SH": 2,
	"T": 1, "TH": 2, "UH": 0, "UW": 0, "V": 2, "W": 6,
	"Y": 6, "Z": 2, "ZH": 2,
}
