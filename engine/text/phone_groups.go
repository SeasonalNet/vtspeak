package text

import (
	"errors"
	"fmt"
	"strings"
)

// Paul2013PhoneGroup is one vowel-centered phone group produced by the
// observed Paul 2013 context-row splitter. The indexes are half-open; Nucleus
// identifies the vowel that anchors the group.
type Paul2013PhoneGroup struct {
	Start      int
	End        int
	OnsetStart int
	Nucleus    int
	// PositionState is the observed context-row value: 1 for the first
	// group, 2 for interior groups, and 3 for the final group in a
	// multi-group span. A single group has value 1.
	PositionState uint8
}

// BuildPaul2013PhoneGroups ports the onset-boundary rule in FUN_10013f30.
// The DLL starts a group at each vowel, then moves the boundary left across a
// permitted singleton or consonant cluster. Remaining phones belong to the
// preceding vowel's group. Inputs with no vowel produce no groups, matching
// the zero-group branch consumed by FUN_10013c00.
func BuildPaul2013PhoneGroups(phones []CMUPhone) ([]Paul2013PhoneGroup, error) {
	if len(phones) == 0 {
		return nil, errors.New("phone group input is empty")
	}
	if len(phones) > 65 {
		return nil, fmt.Errorf("phone group input has %d phones; Paul 2013 splitter caps input at 65 phones", len(phones))
	}

	vowels := make([]int, 0, len(phones)/2)
	for index, phone := range phones {
		if _, err := phone.TreeFeatures(); err != nil {
			return nil, fmt.Errorf("phone %d: %w", index, err)
		}
		if phone.Vowel {
			vowels = append(vowels, index)
		}
	}
	if len(vowels) == 0 {
		return nil, nil
	}
	starts := make([]int, len(vowels))
	for groupIndex, nucleus := range vowels {
		start := nucleus
		lowerBound := 0
		if groupIndex > 0 {
			lowerBound = vowels[groupIndex-1] + 1
		}
		for candidate := nucleus - 1; candidate >= lowerBound; candidate-- {
			if !paul2013AllowedOnset(phones[candidate:nucleus]) {
				break
			}
			start = candidate
		}
		starts[groupIndex] = start
	}
	// FUN_10013f30 forces the first one-based onset position to 1 so initial
	// consonants remain in the first group even when that cluster is not in the
	// ordinary onset table.
	starts[0] = 0

	groups := make([]Paul2013PhoneGroup, len(vowels))
	for index, nucleus := range vowels {
		start := starts[index]
		if index == 0 {
			start = 0
		}
		end := len(phones)
		if index+1 < len(starts) {
			end = starts[index+1]
		}
		if start >= end || starts[index] < start || nucleus < starts[index] || nucleus >= end {
			return nil, fmt.Errorf("derived invalid Paul 2013 phone group %d [%d,%d) around onset %d and nucleus %d", index, start, end, starts[index], nucleus)
		}
		groups[index] = Paul2013PhoneGroup{
			Start:         start,
			End:           end,
			OnsetStart:    starts[index],
			Nucleus:       nucleus,
			PositionState: paul2013GroupPositionState(index, len(vowels)),
		}
	}
	return groups, nil
}

func paul2013GroupPositionState(index, count int) uint8 {
	if index == 0 {
		return 1
	}
	if index == count-1 {
		return 3
	}
	return 2
}

func paul2013AllowedOnset(phones []CMUPhone) bool {
	if len(phones) == 1 {
		return !phones[0].Vowel && phones[0].Label != "NG"
	}
	if len(phones) < 2 || len(phones) > 3 {
		return false
	}
	labels := make([]string, len(phones))
	for index, phone := range phones {
		if phone.Vowel {
			return false
		}
		labels[index] = phone.Label
	}
	_, ok := paul2013OnsetClusters[strings.Join(labels, " ")]
	return ok
}

// Distinct strings referenced by FUN_10013da0's onset-cluster table at
// PTR_DAT_10079a4c. Duplicate pointers in the DLL are collapsed here.
var paul2013OnsetClusters = map[string]struct{}{
	"B L": {}, "B R": {}, "B Y": {},
	"D R": {}, "D W": {}, "D Y": {},
	"F L": {}, "F R": {}, "F Y": {},
	"G L": {}, "G R": {}, "G W": {}, "G Y": {},
	"K L": {}, "K R": {}, "K W": {}, "K Y": {},
	"M Y": {}, "N Y": {},
	"P L": {}, "P R": {}, "P Y": {}, "P W": {},
	"S F": {}, "S FR": {}, "S K": {}, "S KL": {}, "S KR": {}, "S KW": {}, "S KY": {},
	"S L": {}, "S M": {}, "S MY": {}, "S N": {}, "S P": {}, "S PL": {}, "S PR": {},
	"S PY": {}, "S T": {}, "S TR": {}, "S TY": {}, "S V": {}, "S W": {}, "S Y": {},
	"T R": {}, "T W": {}, "T Y": {}, "V Y": {}, "W Y": {},
}
