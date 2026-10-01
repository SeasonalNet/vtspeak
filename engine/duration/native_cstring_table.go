package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
)

// Paul2013NativeCStringTable is the portable view of FUN_100560a0's sorted
// pointer array when comparison mode 0x53 selects native strcmp. The table
// entries are runtime data; this type implements the observed search over a
// caller-supplied snapshot without claiming to produce that data.
type Paul2013NativeCStringTable struct {
	entries [][]byte
}

// NewPaul2013NativeCStringTable copies and validates a nondecreasing array of
// NUL-free C-string values. Duplicate keys are retained because the native
// binary search may return any matching duplicate reached by its midpoint
// sequence.
func NewPaul2013NativeCStringTable(entries [][]byte) (*Paul2013NativeCStringTable, error) {
	table := &Paul2013NativeCStringTable{entries: make([][]byte, len(entries))}
	for index, entry := range entries {
		if bytes.IndexByte(entry, 0) >= 0 {
			return nil, fmt.Errorf("native C-string table entry %d contains an embedded NUL", index)
		}
		table.entries[index] = append([]byte(nil), entry...)
		if index > 0 && bytes.Compare(table.entries[index-1], table.entries[index]) > 0 {
			return nil, fmt.Errorf("native C-string table entries %d and %d are out of strcmp order", index-1, index)
		}
	}
	return table, nil
}

// Lookup returns the exact-match index produced by the native midpoint
// binary-search loop in FUN_100560a0 when param_4 is 0x53. A miss returns -1.
func (table *Paul2013NativeCStringTable) Lookup(key []byte) (int, bool, error) {
	if table == nil {
		return -1, false, errors.New("native C-string table is nil")
	}
	if nul := bytes.IndexByte(key, 0); nul >= 0 {
		key = key[:nul]
	}
	if len(key) == 0 {
		return -1, false, nil
	}
	low := 0
	high := len(table.entries) - 1
	for low <= high {
		middle := low + (high-low)/2
		comparison := bytes.Compare(table.entries[middle], key)
		if comparison == 0 {
			return middle, true, nil
		}
		if comparison < 0 {
			low = middle + 1
		} else {
			high = middle - 1
		}
	}
	return -1, false, nil
}

// SortedSourceLookup adapts this table to the callback consumed by the
// FUN_100086c0 post-neighbor decision. It preserves context cancellation and
// reports only membership because the parent function checks whether the
// native index is -1.
func (table *Paul2013NativeCStringTable) SortedSourceLookup(
	ctx context.Context,
	key []byte,
) (bool, error) {
	if ctx == nil {
		return false, errors.New("native sorted-source lookup has no context")
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	_, found, err := table.Lookup(key)
	if err != nil {
		return false, err
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return found, nil
}
