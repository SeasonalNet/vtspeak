package duration

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013ContextPhraseWindowMatch identifies the table entry and the
// neighboring context rows that matched it in FUN_1000cc60.
type Paul2013ContextPhraseWindowMatch struct {
	TableIndex      int
	FirstContextRow int
	ContextRowCount int
}

// FindPaul2013ContextPhraseWindow ports FUN_1000cc60 over caller-supplied
// sorted table entries and a caller-supplied native string comparator. Mode
// values are the observed 'L', 'R', and 'B'; the window is capped at 20 rows
// and split across both sides for 'B'. A miss returns ok=false.
func FindPaul2013ContextPhraseWindow(
	model []byte,
	rowIndex int,
	mode byte,
	window uint,
	entries [][]byte,
	compare func(left, right []byte) int,
) (Paul2013ContextPhraseWindowMatch, bool, error) {
	if compare == nil {
		return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("FUN_1000cc60 has no native string comparator")
	}
	if mode != 'L' && mode != 'R' && mode != 'B' {
		return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("FUN_1000cc60 mode %#x is unsupported", mode)
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("FUN_1000cc60 row %d is outside model row count %d", rowIndex, rowCount)
	}
	if window > 20 {
		window = 20
	}
	if mode == 'B' {
		if window&1 != 0 {
			window++
		}
		window /= 2
	}
	leftCount, rightCount := 0, 0
	if mode == 'L' || mode == 'B' {
		leftCount = rowIndex
		if leftCount > int(window) {
			leftCount = int(window)
		}
	}
	if mode == 'R' || mode == 'B' {
		rightCount = rowCount - rowIndex - 1
		if rightCount > int(window) {
			rightCount = int(window)
		}
	}
	contextRows := make([]int, 0, leftCount+rightCount)
	for index := rowIndex - leftCount; index < rowIndex; index++ {
		contextRows = append(contextRows, index)
	}
	for index := rowIndex + 1; index <= rowIndex+rightCount; index++ {
		contextRows = append(contextRows, index)
	}
	if len(contextRows) == 0 || len(entries) == 0 {
		return Paul2013ContextPhraseWindowMatch{}, false, nil
	}
	for index, entry := range entries {
		if bytes.IndexByte(entry, 0) >= 0 {
			return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("FUN_1000cc60 table entry %d contains an embedded NUL", index)
		}
		if index > 0 && compare(entries[index-1], entry) > 0 {
			return Paul2013ContextPhraseWindowMatch{}, false, fmt.Errorf("FUN_1000cc60 table entries %d and %d are not sorted under the native comparator", index-1, index)
		}
	}
	for start := range contextRows {
		query, err := paul2013ContextSurfaceBytes(model, contextRows[start])
		if err != nil {
			return Paul2013ContextPhraseWindowMatch{}, false, err
		}
		if len(query) == 0 {
			continue
		}
		low, high := 0, len(entries)-1
		for low <= high {
			middle := low + (high-low)/2
			entry := entries[middle]
			prefixLength := len(query)
			if prefixLength > len(entry) {
				prefixLength = len(entry)
			}
			comparison := compare(entry[:prefixLength], query)
			if comparison == 0 && len(entry) < len(query) {
				comparison = -1
			}
			if comparison == 0 {
				suffix := entry[len(query):]
				if len(suffix) == 0 {
					return Paul2013ContextPhraseWindowMatch{
						TableIndex: middle, FirstContextRow: contextRows[start], ContextRowCount: 1,
					}, true, nil
				}
				if suffix[0] == ' ' {
					matchedRows := 1
					for next := start + 1; next < len(contextRows) && len(suffix) > 0; next++ {
						if suffix[0] != ' ' {
							break
						}
						suffix = suffix[1:]
						space := bytes.IndexByte(suffix, ' ')
						part := suffix
						if space >= 0 {
							part = suffix[:space]
						}
						nextSurface, err := paul2013ContextSurfaceBytes(model, contextRows[next])
						if err != nil {
							return Paul2013ContextPhraseWindowMatch{}, false, err
						}
						if compare(part, nextSurface) != 0 {
							break
						}
						matchedRows++
						if space < 0 {
							return Paul2013ContextPhraseWindowMatch{
								TableIndex: middle, FirstContextRow: contextRows[start], ContextRowCount: matchedRows,
							}, true, nil
						}
						suffix = suffix[space:]
					}
				}
				high = middle - 1
				continue
			}
			if comparison < 0 {
				low = middle + 1
			} else {
				high = middle - 1
			}
		}
	}
	return Paul2013ContextPhraseWindowMatch{}, false, nil
}

func paul2013ContextSurfaceBytes(model []byte, rowIndex int) ([]byte, error) {
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return nil, fmt.Errorf("read FUN_1000cc60 context row %d: %w", rowIndex, err)
	}
	if end := bytes.IndexByte(surface, 0); end >= 0 {
		return surface[:end], nil
	}
	return nil, fmt.Errorf("FUN_1000cc60 context row %d has no NUL-terminated surface", rowIndex)
}
