package voice

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"vtspeak/engine/text"
)

// UnitLocation identifies a unit in one of the Paul model's merged banks.
type UnitLocation struct {
	Bank  string
	Index uint32
}

// ClassRecord is a group of units with an identical transformed selection
// key. ID is its position in the lexicographically sorted class table.
type ClassRecord struct {
	ID      uint32
	Key     [5]byte
	Members []UnitLocation
}

// ClassCatalog is the sorted class index reconstructed from unit signatures.
// It is generated in memory from the read-only unit indexes; vendor resources
// are never changed.
type ClassCatalog struct {
	classes []ClassRecord
}

// BuildPaul2013ClassCatalog transforms each indexed unit signature through
// the observed seven-to-five-byte key mapping, groups equal keys, and assigns
// IDs in the sorted order used by the class-table search. Bank order is the
// documented Paul descriptor order: gen, num, etc, alp.
func BuildPaul2013ClassCatalog(ctx context.Context, model *Paul2013) (*ClassCatalog, error) {
	if model == nil {
		return nil, errors.New("class catalog has no Paul model")
	}
	type classKey = [5]byte
	groups := make(map[classKey][]UnitLocation)
	for _, bankName := range paul2013BankNames {
		bank, ok := model.Banks[bankName]
		if !ok || bank == nil {
			return nil, fmt.Errorf("class catalog is missing Paul %s bank", bankName)
		}
		for unitIndex := uint32(0); unitIndex < bank.UnitCount(); unitIndex++ {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			record, err := bank.ReadRecord(unitIndex)
			if err != nil {
				return nil, fmt.Errorf("read Paul %s unit %d signature: %w", bankName, unitIndex, err)
			}
			key := (text.Context{Signature: record.Signature}).Paul2013SelectionKey()
			groups[key] = append(groups[key], UnitLocation{Bank: bankName, Index: unitIndex})
		}
	}
	keys := make([]classKey, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return compareClassKeys(keys[i], keys[j]) < 0
	})
	classes := make([]ClassRecord, len(keys))
	for index, key := range keys {
		classes[index] = ClassRecord{ID: uint32(index), Key: key, Members: groups[key]}
	}
	return &ClassCatalog{classes: classes}, nil
}

// Len reports the number of distinct transformed signature keys.
func (catalog *ClassCatalog) Len() int {
	if catalog == nil {
		return 0
	}
	return len(catalog.classes)
}

// Lookup returns the exact sorted class for key. The returned member slice is
// copied so callers cannot mutate the catalog.
func (catalog *ClassCatalog) Lookup(key [5]byte) (ClassRecord, bool) {
	if catalog == nil {
		return ClassRecord{}, false
	}
	index := sort.Search(len(catalog.classes), func(index int) bool {
		return compareClassKeys(catalog.classes[index].Key, key) >= 0
	})
	if index >= len(catalog.classes) || catalog.classes[index].Key != key {
		return ClassRecord{}, false
	}
	result := catalog.classes[index]
	result.Members = append([]UnitLocation(nil), result.Members...)
	return result, true
}

// LookupContext transforms one existing seven-byte model context to its
// five-byte key, then returns the exact matching class when present.
func (catalog *ClassCatalog) LookupContext(context text.Context) (ClassRecord, bool) {
	return catalog.Lookup(context.Paul2013SelectionKey())
}

// Records returns the complete sorted class table with independent member
// slices.
func (catalog *ClassCatalog) Records() []ClassRecord {
	if catalog == nil {
		return nil
	}
	result := make([]ClassRecord, len(catalog.classes))
	for index, class := range catalog.classes {
		result[index] = class
		result[index].Members = append([]UnitLocation(nil), class.Members...)
	}
	return result
}

func compareClassKeys(left, right [5]byte) int {
	for position := range left {
		if left[position] < right[position] {
			return -1
		}
		if left[position] > right[position] {
			return 1
		}
	}
	return 0
}
