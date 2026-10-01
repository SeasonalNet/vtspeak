package text

import (
	"bytes"
	"errors"
	"fmt"
)

// Paul2013NativeStringSearchTable is a copied snapshot of FUN_100035d0's
// count-plus-pointer-array input. ComparisonMode 'I' selects FUN_1001c2c0;
// every other mode selects strcmp. The entries are required to be sorted by
// the selected comparator, as they are in the native resource table.
type Paul2013NativeStringSearchTable struct {
	entries          [][]byte
	comparisonMode   byte
	characterWeights [256]int16
}

// NewPaul2013NativeStringSearchTable copies a sorted array of NUL-free C
// strings. Empty entries are allowed because FUN_100035d0 itself does not
// reject an empty key or table row.
func NewPaul2013NativeStringSearchTable(
	entries [][]byte,
	comparisonMode byte,
	characterWeights [256]int16,
) (*Paul2013NativeStringSearchTable, error) {
	table := &Paul2013NativeStringSearchTable{
		entries:          make([][]byte, len(entries)),
		comparisonMode:   comparisonMode,
		characterWeights: characterWeights,
	}
	for index, entry := range entries {
		if bytes.IndexByte(entry, 0) >= 0 {
			return nil, fmt.Errorf("native string-search entry %d contains an embedded NUL", index)
		}
		table.entries[index] = append([]byte(nil), entry...)
		if index == 0 {
			continue
		}
		comparison := table.compare(table.entries[index-1], table.entries[index])
		if comparison > 0 {
			return nil, fmt.Errorf("native string-search entries %d and %d are not sorted for mode %q", index-1, index, comparisonMode)
		}
	}
	return table, nil
}

// Lookup ports FUN_100035d0's inclusive midpoint binary search. It returns
// the native zero-based matching index; a miss returns -1. Unlike
// FUN_100560a0, this helper allows an empty C-string key.
func (table *Paul2013NativeStringSearchTable) Lookup(key []byte) (int, bool, error) {
	if table == nil {
		return -1, false, errors.New("native string-search table is nil")
	}
	key = cString(key)
	low, high := 0, len(table.entries)-1
	for low <= high {
		middle := low + (high-low)/2
		comparison := table.compare(table.entries[middle], key)
		switch {
		case comparison == 0:
			return middle, true, nil
		case comparison < 0:
			low = middle + 1
		default:
			high = middle - 1
		}
	}
	return -1, false, nil
}

func (table *Paul2013NativeStringSearchTable) compare(left, right []byte) int {
	if table.comparisonMode == 'I' {
		return ComparePaul2013MappedCString(left, right, table.characterWeights)
	}
	return bytes.Compare(cString(left), cString(right))
}
