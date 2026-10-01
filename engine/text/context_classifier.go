package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const paul2013ContextClassifierRowStride = 0x70
const paul2013ContextClassifierRowTextOffset = 0x0b

// The six sorted pointer tables and literal strings below were read from the
// local 2013 Paul DLL by extract_paul2013_context_classifier_tables.py. Their
// table roles are kept as function names because their linguistic meanings
// have not been established.
var paul2013ContextWordsE540 = byteWords(
	"about", "above", "across", "after", "again", "against", "along", "among", "amongst", "around", "as", "aside", "at", "away", "before", "behind", "below", "beside", "besides", "beyond", "both", "but", "by", "down", "every", "except", "far", "for", "from", "general", "in", "into", "just", "nearby", "next", "no", "now", "of", "off", "on", "once", "onto", "otherwise", "out", "over", "so", "such", "then", "through", "throughout", "till", "to", "too", "toward", "towards", "under", "unlike", "up", "upon", "via", "well", "with", "without", "yet",
)

var paul2013ContextWordsE4E0 = byteWords("n't", "never", "not")
var paul2013ContextWordsE4B0 = byteWords("and", "both", "but", "either", "neither", "nor", "or")
var paul2013ContextWordsE510 = byteWords("how", "what", "whatever", "when", "whence", "whenever", "where", "whereby", "whereever", "which", "whichever", "who", "whoever", "whom", "whose", "whoseever", "why")
var paul2013ContextWordsE430 = byteWords("'s", "am", "are", "is", "was", "were")
var paul2013ContextWordsE3A0 = byteWords("'d", "'ll", "can", "can't", "cannot", "canst", "could", "couldest", "couldn't", "dare", "did", "do", "does", "may", "mayest", "mayst", "might", "must", "mustn't", "shall", "shalt", "shan't", "should", "shouldest", "shouldn't", "shouldst", "will", "won't", "would", "wouldn't")

// ClassifyPaul2013ModelContextRow ports FUN_1000e570 for one selector into
// the count-prefixed 0x70-byte context table. Results are the native opaque
// bytes 'P', 'N', and 'D'. The selector addresses a row; callers must retain
// native row ordering.
func ClassifyPaul2013ModelContextRow(contextTable []byte, selector int) (uint16, error) {
	if len(contextTable) < 2 {
		return 0, fmt.Errorf("context classifier table has %d bytes, need its count word", len(contextTable))
	}
	count := int(int16(binary.LittleEndian.Uint16(contextTable[:2])))
	if count <= 0 {
		return 0, fmt.Errorf("context classifier has invalid signed row count %d", count)
	}
	if selector < 0 || selector >= count {
		return 0, fmt.Errorf("context classifier selector %d outside %d rows", selector, count)
	}
	rowsEnd := paul2013ContextClassifierRowTextOffset + count*paul2013ContextClassifierRowStride
	if len(contextTable) < rowsEnd {
		return 0, fmt.Errorf("context classifier table has %d bytes, need %d for %d rows", len(contextTable), rowsEnd, count)
	}
	for index := 0; index < count; index++ {
		start := paul2013ContextClassifierRowTextOffset + index*paul2013ContextClassifierRowStride
		if bytes.IndexByte(contextTable[start:start+paul2013ContextClassifierRowStride], 0) < 0 {
			return 0, fmt.Errorf("context classifier row %d has no NUL-terminated text within its 0x70-byte record", index)
		}
	}
	weights := paul2013ContextCharacterWeights()
	var classify func(index int) (uint16, error)
	classify = func(index int) (uint16, error) {
		if index < 0 || index >= count {
			return 0, fmt.Errorf("recursive context classifier selector %d outside %d rows", index, count)
		}
		rowText := func(row int) []byte {
			start := paul2013ContextClassifierRowTextOffset + row*paul2013ContextClassifierRowStride
			return cString(contextTable[start : start+paul2013ContextClassifierRowStride])
		}
		equal := func(row int, literal string) bool {
			return ComparePaul2013MappedCString(rowText(row), []byte(literal), weights) == 0
		}
		contains := func(table [][]byte, row int) bool {
			return FindPaul2013SortedCString(table, rowText(row), weights) >= 0
		}
		e430 := func(row, position int) bool {
			if position > 0 {
				return contains(paul2013ContextWordsE430, row)
			}
			return position == 0 && !equal(row, "'s") && contains(paul2013ContextWordsE430, row)
		}
		e3a0 := func(row int, other []byte) bool {
			if contains(paul2013ContextWordsE3A0, row) {
				return true
			}
			if (equal(row, "ca") || equal(row, "sha") || equal(row, "wo")) &&
				ComparePaul2013MappedCString(other, []byte("n't"), weights) == 0 {
				return true
			}
			return false
		}
		otherText := func(row int) []byte {
			if row < 0 || row >= count {
				return nil
			}
			return rowText(row)
		}
		last := count - 1
		inFirstGroup := contains(paul2013ContextWordsE540, index) ||
			contains(paul2013ContextWordsE4E0, index) ||
			contains(paul2013ContextWordsE4B0, index)
		if inFirstGroup && index+1 < last {
			nextClass, err := classify(index + 1)
			if err != nil {
				return 0, err
			}
			if nextClass == 'P' || nextClass == 'N' {
				return nextClass, nil
			}
		}

		if equal(index, "why") {
			return 'P', nil
		}
		if equal(index, "how") && index+1 < last &&
			(equal(index+1, "about") || equal(index+1, "come") || equal(index+1, "much") || equal(index+1, "say")) {
			return 'P', nil
		}
		if equal(index, "what") {
			if index+1 < last &&
				(equal(index+1, "about") || equal(index+1, "if") || equal(index+1, "in") || equal(index+1, "thouth") || equal(index+1, "then") || equal(index+1, "time") || equal(index+1, "else")) {
				return 'P', nil
			}
			if index+2 < last && equal(index+1, "kind") && equal(index+2, "of") {
				return 'P', nil
			}
		}

		if contains(paul2013ContextWordsE510, index) {
			if index+1 < last && (e430(index+1, index+1) || e3a0(index+1, []byte("")) || equal(index+1, ",")) {
				return 'P', nil
			}
			if index+2 < last && e3a0(index+1, otherText(index+2)) {
				return 'P', nil
			}
			if index+1 >= last {
				return 'P', nil
			}
			for scan := index + 1; scan < last; scan++ {
				if !equal(scan, ",") {
					continue
				}
				if index+1 >= last {
					return 'P', nil
				}
				for cursor := index + 1; cursor < last; cursor++ {
					if equal(cursor, ",") && cursor+1 < last {
						if equal(cursor+1, "if") {
							search := cursor + 1
							for search < last && !equal(search, ",") {
								search++
							}
							if search >= last {
								return 'P', nil
							}
						} else {
							nested, err := classify(cursor + 1)
							if err != nil {
								return 0, err
							}
							if nested == 'N' {
								return 'N', nil
							}
						}
					}
				}
				return 'P', nil
			}
			return 'P', nil
		}

		if index < last && (e430(index, index) || e3a0(index, []byte(""))) ||
			(index+1 < last && e3a0(index, otherText(index+1))) {
			return 'N', nil
		}
		return 'D', nil
	}
	return classify(selector)
}

func paul2013ContextCharacterWeights() [256]int16 {
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	var weights [256]int16
	for index, value := range characterMap {
		weights[index] = int16(value)
	}
	return weights
}

// Paul2013ContextCharacterWeights exposes the local DLL's mapped C-string
// comparison order used by FUN_1001c2c0 in context-table lookups.
func Paul2013ContextCharacterWeights() [256]int16 {
	return paul2013ContextCharacterWeights()
}

func byteWords(words ...string) [][]byte {
	result := make([][]byte, len(words))
	for index, word := range words {
		result[index] = []byte(word)
	}
	return result
}
