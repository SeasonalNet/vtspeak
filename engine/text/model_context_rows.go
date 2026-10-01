package text

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

const (
	paul2013ModelParserTokenRowsOffset = 2
	paul2013ModelParserOffsetRowStride = 0x94
	paul2013ModelParserOffsetRowClass  = 0x24
	paul2013ModelParserOffsetRowType   = 0x2c
	paul2013ModelContextRowBaseOffset  = 0x429a8
	paul2013ModelContextTokenIndex     = 0x00
	paul2013ModelContextType           = 0x02
	paul2013ModelContextStatus         = 0x04
	paul2013ModelContextSurface        = 0x05
	paul2013ModelContextPhoneString    = 0x23
	paul2013ModelContextClass          = 0x64
	paul2013ModelContextMetadata       = 0x68
)

// Paul2013ModelMultiAlternativePhoneSelector implements the FUN_100068b0 call
// made by FUN_1000ea20 for token rows with at least two pronunciation alternatives. It
// receives the contiguous 0x554-byte token rows, the current zero-based row,
// and the original source length; its returned bytes become the output
// context row's phone string.
type Paul2013ModelMultiAlternativePhoneSelector func(
	tokenRows []byte,
	rowIndex int,
	sourceLength int,
) ([]byte, error)

// SelectPaul2013ModelPronunciationAlternative ports the recovered scoring
// path in FUN_100068b0 for multi-alternative token rows. It rebuilds the
// currently recovered lexical features from model token surfaces, evaluates each
// path code with engbi.tree3, and returns the phone string at the earliest
// maximum-scoring path group. It also ports FUN_100068b0's short-input fast
// path for surfaces other than the mapped string "a". Callers should run
// SelectPaul2013ModelNameContextPronunciation first; it includes the known
// `August` after `an`/`the` shortcut. Other contexts handled by
// FUN_10003ff0 remain unresolved, and this scorer does not detect them.
func SelectPaul2013ModelPronunciationAlternative(
	tree *tree3.Tree,
) Paul2013ModelMultiAlternativePhoneSelector {
	return func(tokenRows []byte, rowIndex int, sourceLength int) ([]byte, error) {
		pronunciationRows, err := ReadPaul2013ModelPronunciationRows(tokenRows)
		if err != nil {
			return nil, fmt.Errorf("read pronunciation context rows: %w", err)
		}
		tokenCount := len(pronunciationRows)
		if rowIndex < 0 || rowIndex >= tokenCount {
			return nil, fmt.Errorf("pronunciation token row %d is outside %d rows", rowIndex, tokenCount)
		}
		lexicalTokens := make([]LexicalToken, tokenCount)
		for index := range lexicalTokens {
			lexicalTokens[index].Surface = string(pronunciationRows[index].Surface)
		}
		row := tokenRows[rowIndex*Paul2013TokenResultRowSize : (rowIndex+1)*Paul2013TokenResultRowSize]
		pronunciationRow := pronunciationRows[rowIndex]
		if sourceLength < 2 && ComparePaul2013MappedCString(
			[]byte(lexicalTokens[rowIndex].Surface), []byte("a"), paul2013ContextCharacterWeights(),
		) != 0 {
			if len(pronunciationRow.PhoneAlternatives) == 1 {
				return append([]byte(nil), pronunciationRow.PhoneAlternatives[0]...), nil
			}
			phoneArea := row[Paul2013TokenResultPhoneStrings : Paul2013TokenResultPhoneStrings+Paul2013TokenResultPhoneStride]
			phoneEnd := bytes.IndexByte(phoneArea, 0)
			if phoneEnd < 0 {
				return nil, fmt.Errorf("model token row %d first alternative has no NUL-terminated phone string", rowIndex)
			}
			return append([]byte(nil), phoneArea[:phoneEnd]...), nil
		}
		nameContext, err := SelectPaul2013ModelNameContextPronunciation(tokenRows, rowIndex)
		if err != nil {
			return nil, fmt.Errorf("select model token row %d name/context pronunciation: %w", rowIndex, err)
		}
		if nameContext.Applicable {
			return append([]byte(nil), nameContext.PhoneString...), nil
		}
		if tree == nil {
			return nil, fmt.Errorf("pronunciation alternative selection has no engbi.tree3")
		}
		alternativeCount := pronunciationRow.AlternativeCount
		if alternativeCount < 2 || alternativeCount > paul2013TokenResultAlternativeCapacity {
			return nil, fmt.Errorf("model token row %d has unsupported alternative count %d", rowIndex, alternativeCount)
		}
		groups, err := ParsePaul2013PronunciationPath(pronunciationRow.PathControlBytes)
		if err != nil {
			return nil, fmt.Errorf("parse model token row %d path controls: %w", rowIndex, err)
		}
		if len(groups) != alternativeCount {
			return nil, fmt.Errorf("model token row %d has %d path groups for %d alternatives", rowIndex, len(groups), alternativeCount)
		}
		classifierOutputs := make([]int16, 0, len(groups))
		for groupIndex, group := range groups {
			for _, pathCode := range group {
				features, present, err := BuildPaul2013PronunciationPathCodeFeatures(lexicalTokens, rowIndex, pathCode)
				if err != nil {
					return nil, fmt.Errorf("build path features for token row %d group %d: %w", rowIndex, groupIndex, err)
				}
				if !present {
					continue
				}
				if features.Available != (1<<Paul2013PronunciationFeatureCount)-1 {
					return nil, fmt.Errorf("token row %d path group %d has incomplete pronunciation features", rowIndex, groupIndex)
				}
				_, output, err := tree.Evaluate(features.Values[:])
				if err != nil {
					return nil, fmt.Errorf("evaluate engbi.tree3 for token row %d group %d: %w", rowIndex, groupIndex, err)
				}
				if len(output) != 1 {
					return nil, fmt.Errorf("engbi.tree3 returned %d outputs, want one", len(output))
				}
				classifierOutputs = append(classifierOutputs, int16(int8(uint8(output[0]))))
			}
		}
		selected, scores, err := RankPaul2013PronunciationPathGroups(groups, classifierOutputs)
		if err != nil {
			return nil, fmt.Errorf("score model token row %d pronunciation paths: %w", rowIndex, err)
		}
		if scores[selected] == 0 && ComparePaul2013MappedCString(
			[]byte(lexicalTokens[rowIndex].Surface), []byte("august"), paul2013ContextCharacterWeights(),
		) == 0 {
			attribute := Paul2013ExceptionCharacterAttributes()[lexicalTokens[rowIndex].Surface[0]]
			fallbackCode := byte(0x0e)
			if attribute&0x80 != 0 {
				fallbackCode = 0x14
			}
			fallbackClass, _, err := Paul2013PronunciationPathClass([]byte{fallbackCode})
			if err != nil {
				return nil, fmt.Errorf("map August pronunciation fallback code: %w", err)
			}
			selected = -1
			for groupIndex, group := range groups {
				if groupContainsClass(group, fallbackClass) {
					selected = groupIndex
					break
				}
			}
			if selected < 0 {
				return nil, fmt.Errorf("model token row %d has no path for the August fallback class", rowIndex)
			}
		}
		if selected < 0 || selected >= len(pronunciationRow.PhoneAlternatives) {
			return nil, fmt.Errorf("model token row %d selected alternative %d outside %d phone strings", rowIndex, selected, len(pronunciationRow.PhoneAlternatives))
		}
		return append([]byte(nil), pronunciationRow.PhoneAlternatives[selected]...), nil
	}
}

// NormalizePaul2013ModelPhoneRows ports FUN_1000ea20's row projection from
// the 0x554-byte model token rows into the model's count-prefixed 0x70-byte
// context rows. parserRows is the caller-arena array consumed by the token
// index, with 0x94-byte records. Multi-alternative phone strings use the
// supplied FUN_100068b0 selector.
func NormalizePaul2013ModelPhoneRows(
	model []byte,
	parserRows []byte,
	sourceLength int,
	selectMultiAlternative Paul2013ModelMultiAlternativePhoneSelector,
) ([]byte, error) {
	if sourceLength < 0 {
		return nil, fmt.Errorf("source length %d is negative", sourceLength)
	}
	if len(model) < paul2013ModelParserTokenRowsOffset {
		return nil, fmt.Errorf("model has %d bytes, need its parser token count", len(model))
	}
	tokenCount := int(int16(binary.LittleEndian.Uint16(model[:2])))
	if tokenCount < 0 {
		return nil, fmt.Errorf("model parser token count %d is negative", tokenCount)
	}
	tokenBytes := tokenCount * Paul2013TokenResultRowSize
	if tokenBytes > len(model)-paul2013ModelParserTokenRowsOffset {
		return nil, fmt.Errorf("model has %d bytes after token count, need %d for %d token rows", len(model)-paul2013ModelParserTokenRowsOffset, tokenBytes, tokenCount)
	}
	contextBytes := 2
	if tokenCount > 0 {
		contextBytes += 4 + tokenCount*Paul2013PhoneContextRowSize
	}
	if len(model) < paul2013ModelContextCountOffset+contextBytes {
		return nil, fmt.Errorf("model has %d bytes, need %d for %d context rows", len(model), paul2013ModelContextCountOffset+contextBytes, tokenCount)
	}
	if len(parserRows)%paul2013ModelParserOffsetRowStride != 0 {
		return nil, fmt.Errorf("parser offset rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013ModelParserOffsetRowStride)
	}
	parserRowCount := len(parserRows) / paul2013ModelParserOffsetRowStride
	working := append([]byte(nil), model...)
	contextTable := working[paul2013ModelContextCountOffset:]
	binary.LittleEndian.PutUint16(contextTable[:2], 0)
	tokenRows := working[paul2013ModelParserTokenRowsOffset : paul2013ModelParserTokenRowsOffset+tokenBytes]
	for rowIndex := 0; rowIndex < tokenCount; rowIndex++ {
		tokenStart := rowIndex * Paul2013TokenResultRowSize
		tokenRow := tokenRows[tokenStart : tokenStart+Paul2013TokenResultRowSize]
		tokenIndex := int(int16(binary.LittleEndian.Uint16(tokenRow[Paul2013TokenResultRowIndex:])))
		if tokenIndex < 0 || tokenIndex >= parserRowCount {
			return nil, fmt.Errorf("model token row %d references parser row %d outside %d rows", rowIndex, tokenIndex, parserRowCount)
		}
		parserStart := tokenIndex * paul2013ModelParserOffsetRowStride
		parserRow := parserRows[parserStart : parserStart+paul2013ModelParserOffsetRowStride]
		contextStart := paul2013ModelContextRowBaseOffset + rowIndex*Paul2013PhoneContextRowSize
		contextRow := working[contextStart : contextStart+Paul2013PhoneContextRowSize]

		surfaceArea := tokenRow[Paul2013TokenResultRowSurface:Paul2013TokenResultPathControls]
		surfaceEnd := bytes.IndexByte(surfaceArea, 0)
		if surfaceEnd < 0 {
			return nil, fmt.Errorf("model token row %d has no NUL-terminated surface in its 0x1e-byte field", rowIndex)
		}
		binary.LittleEndian.PutUint16(contextRow[paul2013ModelContextTokenIndex:], uint16(tokenIndex))
		sourceType := int32(binary.LittleEndian.Uint32(parserRow[paul2013ModelParserOffsetRowType:]))
		contextType := uint16(1)
		if sourceType == 0 {
			contextType = 0
		} else if sourceType == 3 {
			contextType = 3
		}
		binary.LittleEndian.PutUint16(contextRow[paul2013ModelContextType:], contextType)
		contextRow[paul2013ModelContextStatus] = tokenRow[Paul2013TokenResultRowType]
		copyNULTerminated(contextRow[paul2013ModelContextSurface:], surfaceArea[:surfaceEnd])
		if parserRow[paul2013ModelParserOffsetRowClass] == 'X' {
			contextRow[paul2013ModelContextClass] = 'X'
		} else {
			contextRow[paul2013ModelContextClass] = '0'
		}
		for metadataIndex := 0; metadataIndex < 4; metadataIndex++ {
			valueOffset := Paul2013TokenResultMetadata + metadataIndex*2
			binary.LittleEndian.PutUint16(
				contextRow[paul2013ModelContextMetadata+metadataIndex*2:],
				binary.LittleEndian.Uint16(tokenRow[valueOffset:]),
			)
		}

		alternativeCount := int(int16(binary.LittleEndian.Uint16(tokenRow[Paul2013TokenResultRowCount:])))
		var phoneString []byte
		if alternativeCount < 2 {
			phoneArea := tokenRow[Paul2013TokenResultPhoneStrings:]
			phoneEnd := bytes.IndexByte(phoneArea, 0)
			if phoneEnd < 0 {
				return nil, fmt.Errorf("model token row %d has no NUL-terminated phone string", rowIndex)
			}
			phoneString = phoneArea[:phoneEnd]
		} else {
			if selectMultiAlternative == nil {
				return nil, fmt.Errorf("model token row %d has %d pronunciation alternatives and requires FUN_100068b0", rowIndex, alternativeCount)
			}
			var err error
			phoneString, err = selectMultiAlternative(tokenRows, rowIndex, sourceLength)
			if err != nil {
				return nil, fmt.Errorf("normalize model token row %d: %w", rowIndex, err)
			}
		}
		phoneCapacity := len(contextRow) - paul2013ModelContextPhoneString
		if len(phoneString) >= phoneCapacity || bytes.IndexByte(phoneString, 0) >= 0 {
			return nil, fmt.Errorf("model token row %d produced a phone string of invalid length or with an embedded NUL", rowIndex)
		}
		copyNULTerminated(contextRow[paul2013ModelContextPhoneString:], phoneString)
		binary.LittleEndian.PutUint16(contextTable[:2], uint16(rowIndex+1))
	}
	return working, nil
}
