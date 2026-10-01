package text

import (
	"reflect"
	"testing"
)

func TestTPPComponentPatternDecodesRuntimePayloadFamilies(t *testing.T) {
	tests := []struct {
		atom     TPPAtom
		family   byte
		count    uint8
		values   []byte
		axMarker bool
	}{
		{atom: TPPAtom{Tag: 'A', Suffix: []byte("0")}, family: 'A', count: 1, values: []byte{0}},
		{atom: TPPAtom{Tag: 'A', Suffix: []byte("1")}, family: 'A', count: 1, values: []byte{1}},
		{atom: TPPAtom{Tag: 'A', Suffix: []byte("X")}, family: 'A', count: 1, axMarker: true},
		{atom: TPPAtom{Tag: 'B', Suffix: []byte("201")}, family: 'B', count: 2, values: []byte{0, 1}},
		{atom: TPPAtom{Tag: 'C', Suffix: []byte("3010")}, family: 'C', count: 3, values: []byte{0, 1, 0}},
		{atom: TPPAtom{Tag: 'D', Suffix: []byte("41010")}, family: 'D', count: 4, values: []byte{1, 0, 1, 0}},
		{atom: TPPAtom{Tag: 'E', Suffix: []byte("500010")}, family: 'E', count: 5, values: []byte{0, 0, 0, 1, 0}},
	}
	for _, test := range tests {
		got, matched, err := test.atom.ComponentPattern()
		if err != nil || !matched {
			t.Fatalf("ComponentPattern(%c%s) = (%+v, %t, %v)", test.atom.Tag, test.atom.Suffix, got, matched, err)
		}
		if got.Family != test.family || got.ComponentCount != test.count ||
			got.IsAXRecognitionTag != test.axMarker || !reflect.DeepEqual(got.ComponentValues, test.values) {
			t.Errorf("ComponentPattern(%c%s) = %+v", test.atom.Tag, test.atom.Suffix, got)
		}
	}
}

func TestTPPComponentPatternRejectsMalformedFamilyPayloads(t *testing.T) {
	for _, atom := range []TPPAtom{
		{Tag: 'A', Suffix: []byte("2")},
		{Tag: 'B', Suffix: []byte("301")},
		{Tag: 'C', Suffix: []byte("3012")},
		{Tag: 'D', Suffix: []byte("4101")},
		{Tag: 'E', Suffix: []byte("50001x")},
	} {
		if _, matched, err := atom.ComponentPattern(); err == nil || matched {
			t.Errorf("ComponentPattern(%c%s) = matched %t, err %v; want malformed-payload error", atom.Tag, atom.Suffix, matched, err)
		}
	}
	if _, matched, err := (TPPAtom{Tag: 'G', Suffix: []byte("95")}).ComponentPattern(); err != nil || matched {
		t.Fatalf("numeric G atom component dispatch = matched %t, err %v; want no component match", matched, err)
	}
}

func TestPaul2013ProperNameTPPSelectorUsesNativeShortCountGate(t *testing.T) {
	for count := 1; count <= 4; count++ {
		got, selected := Paul2013ProperNameTPPSelector(count)
		want := byte('A' + count - 1)
		if !selected || got != want {
			t.Errorf("selector for %d components = (%q, %t), want (%q, true)", count, got, selected, want)
		}
	}
	for _, count := range []int{0, 5, 6} {
		if got, selected := Paul2013ProperNameTPPSelector(count); selected {
			t.Errorf("selector for %d components = %q, true; native caller skips selector construction", count, got)
		}
	}
}

func TestLookupTPPComponentPatternFindsAAtomBeforeNumericSuffix(t *testing.T) {
	tables := Paul2013EmbeddedKeyTables()
	key, err := EncodeEmbeddedKey([]byte("ACCORD"), tables)
	if err != nil {
		t.Fatal(err)
	}
	dictionary := &TPPDictionary{records: map[string]TPPRecord{
		string(key): {Atoms: []TPPAtom{
			{Tag: 'A', Suffix: []byte("0")},
			{Tag: 'G', Suffix: []byte("95")},
		}},
	}}

	got, found, err := dictionary.LookupTPPComponentPattern([]byte("ACCORD"), tables)
	if err != nil || !found {
		t.Fatalf("LookupTPPComponentPattern(ACCORD) = (%+v, %t, %v)", got, found, err)
	}
	if got.Family != 'A' || got.ComponentCount != 1 || !reflect.DeepEqual(got.ComponentValues, []byte{0}) {
		t.Fatalf("ACCORD component pattern = %+v, want one A-family zero value", got)
	}
}

func TestLookupTPPComponentPatternForSelectorDoesNotSubstituteOtherAtoms(t *testing.T) {
	tables := Paul2013EmbeddedKeyTables()
	key, err := EncodeEmbeddedKey([]byte("ACCORD"), tables)
	if err != nil {
		t.Fatal(err)
	}
	dictionary := &TPPDictionary{records: map[string]TPPRecord{
		string(key): {Atoms: []TPPAtom{{Tag: 'A', Suffix: []byte("0")}, {Tag: 'G', Suffix: []byte("95")}}},
	}}

	if _, found, err := dictionary.LookupTPPComponentPatternForSelector([]byte("ACCORD"), 'B', tables); err != nil || found {
		t.Fatalf("B selector against A G record = found %t, err %v; want miss", found, err)
	}
	if _, _, err := dictionary.LookupTPPComponentPatternForSelector([]byte("ACCORD"), 'G', tables); err == nil {
		t.Fatal("non-component selector G was accepted")
	}
}
