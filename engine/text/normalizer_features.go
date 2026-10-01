package text

import (
	"errors"
	"fmt"
)

// Paul2013NormalizerFeatureWindow is the ten-short input assembled by
// FUN_10002db0 before its caller evaluates a character-class tree. The values
// remain numeric; their model semantics are not assigned here.
type Paul2013NormalizerFeatureWindow [10]int16

// Paul2013NormalizerTreeIndex maps one supported classifier character to its
// observed position in the shared 27-tree ATMT container. ASCII case folds to
// the same letter tree; apostrophe selects the final tree.
func Paul2013NormalizerTreeIndex(value byte) (int, error) {
	switch {
	case value >= 'A' && value <= 'Z':
		return int(value - 'A'), nil
	case value >= 'a' && value <= 'z':
		return int(value - 'a'), nil
	case value == '\'':
		return 26, nil
	default:
		return 0, fmt.Errorf("byte 0x%02x has no Paul 2013 normalizer ATMT tree", value)
	}
}

// BuildPaul2013NormalizerFeatureWindow ports FUN_10002db0's feature layout.
// It maps the source bytes at index-3 through index+3, zero-fills missing
// source neighbors, and appends up to three signed category bytes selected
// by the source-relative suffix positions. Category values remain caller
// inputs; character mapping uses the recovered DLL table behavior.
func BuildPaul2013NormalizerFeatureWindow(
	source []byte,
	categoryBytes []byte,
	index int,
) (Paul2013NormalizerFeatureWindow, error) {
	var result Paul2013NormalizerFeatureWindow
	if len(source) == 0 {
		return result, errors.New("normalizer feature source is empty")
	}
	if index < 0 || index >= len(source) {
		return result, fmt.Errorf("normalizer feature index %d is outside source length %d", index, len(source))
	}
	categoryIndexes := [3]int{len(source) - index - 2, len(source) - index - 3, len(source) - index - 4}
	for offset := -3; offset <= 3; offset++ {
		sourceIndex := index + offset
		if sourceIndex >= 0 && sourceIndex < len(source) {
			result[offset+3] = paul2013NormalizerCharacterClass(source[sourceIndex])
		}
	}
	for categoryOffset, categoryIndex := range categoryIndexes {
		if index+categoryOffset+1 >= len(source) {
			continue
		}
		if categoryIndex < 0 || categoryIndex >= len(categoryBytes) {
			return Paul2013NormalizerFeatureWindow{}, fmt.Errorf(
				"normalizer category index %d for feature %d is outside category length %d",
				categoryIndex, categoryOffset, len(categoryBytes),
			)
		}
		result[7+categoryOffset] = int16(int8(categoryBytes[categoryIndex]))
	}
	return result, nil
}

func paul2013NormalizerCharacterClass(value byte) int16 {
	switch {
	case value >= 'A' && value <= 'Z':
		return int16(value - '@')
	case value >= 'a' && value <= 'z':
		// FUN_10002d70 reads DAT_1007e688 before subtracting 0x40. The
		// recovered table maps ASCII lowercase to uppercase.
		return int16(value-'a') + 1
	case value == '\'':
		return 0x1c
	default:
		return 0
	}
}
