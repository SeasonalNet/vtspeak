package text

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"vtspeak/engine/tree3"
)

// Paul2013TerminalModel joins the FUN_10052520 key producer with the shared
// English sbd.tree3 scalar evaluator used by FUN_10051cc0's A branch.
type Paul2013TerminalModel struct {
	abbreviations [3]*Paul2013TerminalAbbreviationTable
	words         *Paul2013TerminalWordClassifier
	tree          *tree3.Tree
}

func NewPaul2013TerminalModel(tables Tables, tree *tree3.Tree) (*Paul2013TerminalModel, error) {
	if tree == nil || tree.OutputWidth != 1 {
		return nil, fmt.Errorf("terminal model requires a scalar sbd.tree3")
	}
	model := &Paul2013TerminalModel{tree: tree}
	for index, name := range []string{"abbrh_sort.txt2", "abbrt_sort.txt2", "abbrc_sort.txt2"} {
		var err error
		model.abbreviations[index], err = NewPaul2013TerminalAbbreviationTable(tables[name])
		if err != nil {
			return nil, err
		}
	}
	var err error
	model.words, err = NewPaul2013TerminalWordClassifier(tables["sbdw_sort.txt2"])
	if err != nil {
		return nil, err
	}
	return model, nil
}

func LoadPaul2013TerminalModel(dictionaryRoot string) (*Paul2013TerminalModel, error) {
	tables := make(Tables, 4)
	for _, name := range []string{"abbrh_sort.txt2", "abbrt_sort.txt2", "abbrc_sort.txt2", "sbdw_sort.txt2"} {
		raw, err := os.ReadFile(filepath.Join(dictionaryRoot, name))
		if err != nil {
			return nil, err
		}
		table, err := ParseTXT2(name, raw)
		if err != nil {
			return nil, err
		}
		tables[name] = table
	}
	tree, err := tree3.ParseFile(filepath.Join(dictionaryRoot, "sbd.tree3"))
	if err != nil {
		return nil, err
	}
	return NewPaul2013TerminalModel(tables, tree)
}

// BuildKey returns the complete 16-short key in native order. Empty neighbors
// write zero only to length; their six other fields retain the initialized -1.
func (model *Paul2013TerminalModel) BuildKey(window Paul2013TerminalContextWindow) (key [16]int16, err error) {
	for index := range key {
		key[index] = -1
	}
	if model == nil {
		return key, fmt.Errorf("terminal model is unavailable")
	}
	for side, index := range []int{2, 4} {
		token := window[index]
		base := side * 7
		key[base] = int16(token.TextLength)
		if token.TextLength == 0 {
			continue
		}
		key[base+1], err = ClassifyPaul2013TerminalTokenCase(token)
		if err != nil {
			return key, err
		}
		mask, maskErr := Paul2013TerminalAbbreviationMask(model.abbreviations, token.Text)
		if maskErr != nil {
			return key, maskErr
		}
		key[base+2] = 0
		if mask&1 != 0 {
			key[base+2] = 1
		}
		key[base+3] = 0
		if mask&2 != 0 {
			key[base+3] = 1
		}
		key[base+4] = 0
		if mask != 0 {
			key[base+4] = 1
		}
		key[base+5] = Paul2013TerminalTokenHasInternalDot(token.Text)
		key[base+6] = ClassifyPaul2013TerminalTokenPunctuation(token)
	}
	if window[4].TextLength != 0 {
		key[14], err = model.words.Classify(window[4].Text)
		if err != nil {
			return key, err
		}
	}
	key[15] = ClassifyPaul2013TerminalTokenPunctuation(window[3])
	return key, nil
}

// Lookup implements Paul2013TerminalModelLookup. The recognizer tests only
// its result's low byte, as the original scalar-short caller does.
func (model *Paul2013TerminalModel) Lookup(ctx context.Context, window Paul2013TerminalContextWindow) (uint32, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	key, err := model.BuildKey(window)
	if err != nil {
		return 0, err
	}
	if model.tree == nil {
		return 0, fmt.Errorf("terminal model tree is unavailable")
	}
	_, output, err := model.tree.Evaluate(key[:])
	if err != nil {
		return 0, err
	}
	if len(output) != 1 {
		return 0, fmt.Errorf("terminal model output is not scalar")
	}
	return uint32(uint16(output[0])), nil
}
