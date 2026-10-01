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
	classes        []ClassRecord
	featureViewIDs [2][]uint32
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
	catalog := &ClassCatalog{classes: classes}
	if err := catalog.buildFeatureViewIndexes(); err != nil {
		return nil, err
	}
	return catalog, nil
}

// LookupFeatureViewPrefix returns classes in reconstructed feature-view index
// order, matching the range-search shape used by FUN_1002df50. The two native
// view vectors are loaded from sclass.idx; this catalog reconstructs them from
// the local unit indexes and orders equal views by ascending class ID. That tie
// order is deterministic but has not been shown to match the native vectors.
func (catalog *ClassCatalog) LookupFeatureViewPrefix(query [10]byte, mode byte, prefixLength, limit int) ([]ClassRecord, error) {
	if catalog == nil {
		return nil, errors.New("feature-view lookup has no class catalog")
	}
	if mode != 1 && mode != 2 {
		return nil, fmt.Errorf("unsupported Paul 2013 feature-view mode %d", mode)
	}
	if prefixLength < 1 || prefixLength > len(query) {
		return nil, fmt.Errorf("feature-view prefix length %d outside native range [1, %d]", prefixLength, len(query))
	}
	if limit < 0 {
		return nil, fmt.Errorf("feature-view result limit %d is negative", limit)
	}
	if limit == 0 {
		return []ClassRecord{}, nil
	}
	ids := catalog.featureViewIDs[mode-1]
	start := sort.Search(len(ids), func(index int) bool {
		return compareFeatureViewPrefix(catalog.classes[ids[index]].Key, query, mode, prefixLength) >= 0
	})
	end := start
	for end < len(ids) && end-start < limit && compareFeatureViewPrefix(catalog.classes[ids[end]].Key, query, mode, prefixLength) == 0 {
		end++
	}
	result := make([]ClassRecord, end-start)
	for index, id := range ids[start:end] {
		class := catalog.classes[id]
		result[index] = class
		result[index].Members = append([]UnitLocation(nil), class.Members...)
	}
	return result, nil
}

func (catalog *ClassCatalog) buildFeatureViewIndexes() error {
	for mode := byte(1); mode <= 2; mode++ {
		ids := make([]uint32, len(catalog.classes))
		views := make([][10]byte, len(catalog.classes))
		for index, class := range catalog.classes {
			if class.ID != uint32(index) {
				return fmt.Errorf("class catalog record %d has non-contiguous ID %d", index, class.ID)
			}
			view, err := text.Paul2013FeatureViewForKey(class.Key, mode)
			if err != nil {
				return fmt.Errorf("build class %d feature view %d: %w", class.ID, mode, err)
			}
			ids[index] = class.ID
			views[index] = view
		}
		sort.Slice(ids, func(i, j int) bool {
			left, right := views[ids[i]], views[ids[j]]
			if comparison := compareFeatureViews(left, right); comparison != 0 {
				return comparison < 0
			}
			return ids[i] < ids[j]
		})
		catalog.featureViewIDs[mode-1] = ids
	}
	return nil
}

func compareFeatureViewPrefix(key [5]byte, query [10]byte, mode byte, prefixLength int) int {
	view, err := text.Paul2013FeatureViewForKey(key, mode)
	if err != nil {
		return 1
	}
	for position := 0; position < prefixLength; position++ {
		if view[position] < query[position] {
			return -1
		}
		if view[position] > query[position] {
			return 1
		}
	}
	return 0
}

func compareFeatureViews(left, right [10]byte) int {
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

// LookupPrefix returns every sorted class whose key begins with the requested
// prefix. prefixLength follows FUN_10019450's caller-supplied comparison
// length and must be in the native key width [1, 5].
func (catalog *ClassCatalog) LookupPrefix(key [5]byte, prefixLength int) ([]ClassRecord, error) {
	return catalog.LookupPrefixLimited(key, prefixLength, int(^uint(0)>>1))
}

// LookupPrefixLimited returns at most limit classes from the beginning of the
// sorted prefix range. It preserves the native class-table order and avoids
// copying unneeded memberships from broad prefixes.
func (catalog *ClassCatalog) LookupPrefixLimited(key [5]byte, prefixLength, limit int) ([]ClassRecord, error) {
	if catalog == nil {
		return nil, errors.New("class prefix lookup has no catalog")
	}
	if prefixLength < 1 || prefixLength > len(key) {
		return nil, fmt.Errorf("class prefix length %d outside native range [1, %d]", prefixLength, len(key))
	}
	if limit < 0 {
		return nil, fmt.Errorf("class prefix result limit %d is negative", limit)
	}
	if limit == 0 {
		return []ClassRecord{}, nil
	}
	start := sort.Search(len(catalog.classes), func(index int) bool {
		return compareClassKeyPrefix(catalog.classes[index].Key, key, prefixLength) >= 0
	})
	end := start
	for end < len(catalog.classes) && end-start < limit && compareClassKeyPrefix(catalog.classes[end].Key, key, prefixLength) == 0 {
		end++
	}
	result := make([]ClassRecord, end-start)
	for index, class := range catalog.classes[start:end] {
		result[index] = class
		result[index].Members = append([]UnitLocation(nil), class.Members...)
	}
	return result, nil
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

func compareClassKeyPrefix(left, right [5]byte, prefixLength int) int {
	for position := 0; position < prefixLength; position++ {
		if left[position] < right[position] {
			return -1
		}
		if left[position] > right[position] {
			return 1
		}
	}
	return 0
}
