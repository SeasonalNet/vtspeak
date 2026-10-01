package duration

import (
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/text"
)

func TestLookupTypedComponentPatternUsesLoadedTPPDictionary(t *testing.T) {
	dictionaryRoot := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(dictionaryRoot, "tppdict_eng")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := text.LoadTPPDictionary(dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{tpp: dictionary}

	got, found, err := engine.LookupTypedComponentPattern([]byte("AGUA-DULCE"))
	if err != nil || !found {
		t.Fatalf("LookupTypedComponentPattern() = (%+v, %t, %v)", got, found, err)
	}
	if got.Family != 'B' || got.ComponentCount != 2 || len(got.ComponentValues) != 2 ||
		got.ComponentValues[0] != 0 || got.ComponentValues[1] != 0 {
		t.Fatalf("component pattern = %+v, want B with two opaque zero values", got)
	}
	properName, err := engine.LookupProperNameTPPComponentPattern([]byte("AGUA-DULCE"), 2)
	if err != nil || !properName.LookupPerformed || !properName.Found || properName.Selector != 'B' ||
		properName.Pattern.ComponentCount != 2 || len(properName.Pattern.ComponentValues) != 2 {
		t.Fatalf("proper-name two-component lookup = (%+v, %v)", properName, err)
	}
	fivePart, err := engine.LookupProperNameTPPComponentPattern([]byte("CASA-DE-ORO-MOUNT-HELIX"), 5)
	if err != nil || fivePart.LookupPerformed || fivePart.Found {
		t.Fatalf("proper-name five-component gate = (%+v, %v), want native lookup bypass", fivePart, err)
	}
	forcedE, found, err := engine.LookupTypedComponentPattern([]byte("CASA-DE-ORO-MOUNT-HELIX"))
	if err != nil || !found || forcedE.Family != 'E' || forcedE.ComponentCount != 5 ||
		len(forcedE.ComponentValues) != 5 || forcedE.ComponentValues[0] != 0 ||
		forcedE.ComponentValues[1] != 0 || forcedE.ComponentValues[2] != 0 ||
		forcedE.ComponentValues[3] != 1 || forcedE.ComponentValues[4] != 0 {
		t.Fatalf("direct five-component E lookup = (%+v, %t, %v)", forcedE, found, err)
	}
}

func TestLookupTypedComponentPatternRequiresLoadedDictionary(t *testing.T) {
	if _, found, err := (&Engine{}).LookupTypedComponentPattern([]byte("AGUA-DULCE")); err == nil || found {
		t.Fatalf("lookup without dictionary = found %t, err %v; want unloaded-resource error", found, err)
	}
}
