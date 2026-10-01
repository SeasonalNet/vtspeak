package duration

import (
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/text"
)

func TestApplyPaul2013TPPNumericTextRowsUsesLoadedSharedResources(t *testing.T) {
	dictionaryRoot := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(dictionaryRoot, "tppdict_eng")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := text.LoadTPPDictionary(dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	tables, err := text.LoadTXT2Tables(dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{tpp: dictionary, txt2: tables}
	rows := []text.Paul2013TPPWindowRow{{Text: []byte("ACCORD")}}
	result, err := engine.ApplyPaul2013TPPNumericTextRows(rows, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Operations) != 1 || result.Operations[0] != (text.Paul2013TPPNumericTextOperation{
		Tag: 'G', StartRow: 0, EndRow: 0,
	}) {
		t.Fatalf("operations = %+v, want one G operation on row 0", result.Operations)
	}
	if len(result.Updates) != 1 || !result.Updates[0].TypedCodeWritten || result.Updates[0].TypedCode != 95 {
		t.Fatalf("updates = %+v, want typed code 95", result.Updates)
	}
	classIndex, found, err := tables["wab.txt2"].LookupPaul2013WABClass(
		[]byte("ACCORD"), text.Paul2013EmbeddedKeyTables().CharacterMap,
	)
	if err != nil {
		t.Fatal(err)
	}
	update := result.Updates[0]
	if update.ClassCodeWritten != found || (found && update.ClassCode != byte(classIndex+1)) {
		t.Fatalf("WAB update = %+v; lookup = (%d, %t)", update, classIndex, found)
	}

	compound, err := engine.ApplyPaul2013TPPNumericTextRows([]text.Paul2013TPPWindowRow{
		{Word2: 1, State: 'S', Text: []byte("ABOUT")},
		{Word0: 1, Word1: 0, Word2: 1, State: 'A', Text: []byte("SHIPPING")},
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(compound.Operations) != 1 || compound.Operations[0] != (text.Paul2013TPPNumericTextOperation{
		Tag: 'F', StartRow: 0, EndRow: 1,
	}) || len(compound.Updates) != 2 || compound.Updates[0].TypedCode != 120 || compound.Updates[1].TypedCode != 120 {
		t.Fatalf("F compound result = %+v, want one F operation with code 120 on both rows", compound)
	}

	var wabOnlySurface []byte
	for _, row := range tables["wab.txt2"].Rows {
		if _, exists, lookupErr := dictionary.Lookup([]byte(row.Key), text.Paul2013EmbeddedKeyTables()); lookupErr != nil {
			t.Fatal(lookupErr)
		} else if !exists {
			wabOnlySurface = []byte(row.Key)
			break
		}
	}
	if wabOnlySurface == nil {
		t.Fatal("could not find a WAB word without a direct TPP record")
	}
	fallback, err := engine.ApplyPaul2013TPPNumericTextRows(
		[]text.Paul2013TPPWindowRow{{Text: wabOnlySurface}}, 0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(fallback.Operations) != 0 || len(fallback.Updates) != 1 || !fallback.Updates[0].ClassCodeWritten {
		t.Fatalf("WAB fallback result = %+v, want class-only update", fallback)
	}
	if _, err := engine.ApplyPaul2013TPPNumericTextRows(rows, 1); err == nil {
		t.Fatal("unsupported nonzero TPP index was accepted")
	}
}
