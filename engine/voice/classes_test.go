package voice

import (
	"testing"

	"vtspeak/engine/text"
)

func TestClassCatalogFeatureViewPrefixUsesIndexedOrder(t *testing.T) {
	keys := [][5]byte{
		{0x21, 0x34, 0x56, 0x95, 0x20},
		{0x22, 0x34, 0x56, 0x95, 0x21},
		{0x21, 0x35, 0x56, 0x95, 0x22},
		{0x21, 0x34, 0x57, 0x95, 0x23},
	}
	classes := make([]ClassRecord, len(keys))
	for index, key := range keys {
		classes[index] = ClassRecord{ID: uint32(index), Key: key}
	}
	catalog := &ClassCatalog{classes: classes}
	if err := catalog.buildFeatureViewIndexes(); err != nil {
		t.Fatal(err)
	}

	for mode := byte(1); mode <= 2; mode++ {
		class := classes[2]
		query, err := text.Paul2013FeatureViewForKey(class.Key, mode)
		if err != nil {
			t.Fatal(err)
		}
		got, err := catalog.LookupFeatureViewPrefix(query, mode, len(query), 10)
		if err != nil {
			t.Fatalf("mode %d lookup: %v", mode, err)
		}
		if len(got) != 1 || got[0].ID != class.ID {
			t.Fatalf("mode %d exact view lookup = %#v, want class %d", mode, got, class.ID)
		}

		got, err = catalog.LookupFeatureViewPrefix(query, mode, 1, 2)
		if err != nil {
			t.Fatalf("mode %d short-prefix lookup: %v", mode, err)
		}
		if len(got) != 2 {
			t.Fatalf("mode %d short-prefix returned %d classes, want limited 2", mode, len(got))
		}
	}
}

func TestClassCatalogFeatureViewPrefixRejectsUnsupportedRanges(t *testing.T) {
	catalog := &ClassCatalog{}
	for _, test := range []struct {
		mode  byte
		width int
		limit int
	}{
		{mode: 0, width: 1, limit: 1},
		{mode: 1, width: 0, limit: 1},
		{mode: 2, width: 11, limit: 1},
		{mode: 1, width: 1, limit: -1},
	} {
		if _, err := catalog.LookupFeatureViewPrefix([10]byte{}, test.mode, test.width, test.limit); err == nil {
			t.Fatalf("mode=%d width=%d limit=%d unexpectedly succeeded", test.mode, test.width, test.limit)
		}
	}
}

func TestClassCatalogFeatureViewPrefixBreaksTiesByClassID(t *testing.T) {
	classes := []ClassRecord{
		{ID: 0, Key: [5]byte{0x21, 0x34, 0x56, 0x95, 0x20}},
		{ID: 1, Key: [5]byte{0x21, 0x34, 0x57, 0x95, 0x20}},
	}
	catalog := &ClassCatalog{classes: classes}
	if err := catalog.buildFeatureViewIndexes(); err != nil {
		t.Fatal(err)
	}
	query, err := text.Paul2013FeatureViewForKey(classes[0].Key, 1)
	if err != nil {
		t.Fatal(err)
	}
	got, err := catalog.LookupFeatureViewPrefix(query, 1, len(query), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != 0 || got[1].ID != 1 {
		t.Fatalf("equal feature-view results = %#v, want class IDs [0 1]", got)
	}
}
