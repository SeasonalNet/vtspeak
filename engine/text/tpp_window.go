package text

import (
	"bytes"
	"errors"
	"fmt"
)

// Paul2013TPPWindowRow exposes the three opaque signed words, state byte, and
// NUL-terminated text used by FUN_1000e160. Their linguistic meanings are not
// assigned here.
type Paul2013TPPWindowRow struct {
	Word0 int32
	Word1 int32
	Word2 int32
	State byte
	Text  []byte
}

// Paul2013TPPWindowInput supplies the recovered low-byte character map and
// the still-unresolved class-indexed TPP lookup around the bounded row scan.
type Paul2013TPPWindowInput struct {
	Rows                  []Paul2013TPPWindowRow
	StartRow              int
	CharacterMap          [256]byte
	LookupCompound        func(compound []byte, dictionaryIndex int) ([]byte, bool, error)
	DictionaryIndex       int
	InitialOutputRowIndex int
}

// Paul2013TPPWindowResult contains the last successful F-code lookup from the
// bounded scan. OutputRowIndex retains InitialOutputRowIndex if no lookup
// succeeds, matching the native pointer write behavior.
type Paul2013TPPWindowResult struct {
	Matched        bool
	Compound       []byte
	Replacement    []byte
	OutputRowIndex int
}

// MatchPaul2013TPPWindow ports the row traversal and control flow in
// FUN_1000e160. It examines following rows while the accumulated +8 values
// remain below 30, skips the observed case-mapped "dash" S-row, processes at
// most five eligible A-state rows, and retains the last successful F lookup.
// The class-indexed F lookup remains an explicit callback.
func MatchPaul2013TPPWindow(input Paul2013TPPWindowInput) (Paul2013TPPWindowResult, error) {
	result := Paul2013TPPWindowResult{OutputRowIndex: input.InitialOutputRowIndex}
	if input.StartRow < 0 || input.StartRow >= len(input.Rows) {
		return result, fmt.Errorf("TPP window start row %d outside %d rows", input.StartRow, len(input.Rows))
	}
	if len(input.Rows) > 100 {
		return result, fmt.Errorf("TPP window has %d rows, exceeds native 100-row state", len(input.Rows))
	}
	if input.LookupCompound == nil {
		return result, errors.New("TPP window requires a compound-lookup callback")
	}

	compound := cString(input.Rows[input.StartRow].Text)
	if len(compound) > 31 {
		return result, errors.New("TPP window initial text exceeds the native 32-byte buffer")
	}
	cumulative := int64(input.Rows[input.StartRow].Word2)
	processed := 1
	for rowIndex := input.StartRow + 1; rowIndex < len(input.Rows); rowIndex++ {
		row := input.Rows[rowIndex]
		cumulative += 1 + int64(row.Word2)
		if cumulative >= 0x1e {
			break
		}

		span := int64(row.Word0) - int64(row.Word1)
		priorEnd := int64(input.Rows[rowIndex-1].Word1)
		gap := int64(row.Word0) - priorEnd
		// The native fast-skip applies only to S-state rows with span 1,
		// mapped equality to "dash", and a gap less than 2.
		if row.State == 'S' && span == 1 && gap <= 1 &&
			Paul2013MappedCStringEqual(cString(row.Text), []byte("dash"), input.CharacterMap) {
			continue
		}
		if row.State != 'A' || span > 1 {
			break
		}

		candidateText := cString(row.Text)
		if len(compound) == 0 {
			compound = append(compound[:0], candidateText...)
		} else {
			compound = append(compound, '-')
			compound = append(compound, candidateText...)
		}
		if len(compound) > 31 {
			return result, errors.New("TPP window compound exceeds the native 32-byte buffer")
		}
		replacement, found, err := input.LookupCompound(append([]byte(nil), compound...), input.DictionaryIndex)
		if err != nil {
			return result, fmt.Errorf("lookup TPP compound at row %d: %w", rowIndex, err)
		}
		if found {
			result.Matched = true
			result.Compound = append(result.Compound[:0], compound...)
			result.Replacement = append(result.Replacement[:0], cString(replacement)...)
			result.OutputRowIndex = rowIndex
		}
		processed++
		if processed > 5 {
			break
		}
	}
	return result, nil
}

// MatchPaul2013TPPWindowWithDictionary connects the recovered F-code path to
// the loaded shared TPP dictionary. FUN_10011820 uses its nonindexed default
// resource when dictionary index is zero; other dictionary indexes are
// rejected because their backing lookup tables are not loaded by this engine.
func MatchPaul2013TPPWindowWithDictionary(
	input Paul2013TPPWindowInput,
	dictionary *TPPDictionary,
	tables EmbeddedKeyTables,
) (Paul2013TPPWindowResult, error) {
	if dictionary == nil {
		return Paul2013TPPWindowResult{OutputRowIndex: input.InitialOutputRowIndex}, errors.New("TPP window dictionary is nil")
	}
	if input.DictionaryIndex != 0 {
		return Paul2013TPPWindowResult{OutputRowIndex: input.InitialOutputRowIndex},
			fmt.Errorf("TPP dictionary index %d is not loaded", input.DictionaryIndex)
	}
	input.LookupCompound = func(compound []byte, dictionaryIndex int) ([]byte, bool, error) {
		if dictionaryIndex != 0 {
			return nil, false, fmt.Errorf("TPP dictionary index %d is not loaded", dictionaryIndex)
		}
		return dictionary.LookupNumericText(compound, 'F', tables)
	}
	return MatchPaul2013TPPWindow(input)
}

// Paul2013MappedCStringEqual ports the equality behavior of FUN_1001c2c0 for
// ASCII strings using the low-byte character-map projection provisioned with
// the embedded key tables. This projection is exact for the supported ASCII
// frontend; non-ASCII mapped comparisons are outside the engine input scope.
func Paul2013MappedCStringEqual(left, right []byte, characterMap [256]byte) bool {
	left = cString(left)
	right = cString(right)
	for index := 0; ; index++ {
		leftByte, rightByte := byte(0), byte(0)
		if index < len(left) {
			leftByte = left[index]
		}
		if index < len(right) {
			rightByte = right[index]
		}
		if characterMap[leftByte] != characterMap[rightByte] {
			return false
		}
		if leftByte == 0 {
			return true
		}
	}
}

func cString(value []byte) []byte {
	if end := bytes.IndexByte(value, 0); end >= 0 {
		return value[:end]
	}
	return value
}
