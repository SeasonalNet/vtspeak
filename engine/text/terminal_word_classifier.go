package text

import (
	"bytes"
	"fmt"
)

// Paul2013TerminalWordClassifier ports FUN_100528b0 using the copied
// sbdw_sort.txt2 table and the DLL's short translation table at 0x10081358.
type Paul2013TerminalWordClassifier struct {
	search  *Paul2013NativeStringSearchTable
	classes []int16
}

var paul2013TerminalWordClasses = [...]int16{16, 16, 15, 13, 14, 11, 17, 12, 7, 9, 6, 2, 0, 8, 5, 5, 5, 9, 4, 1, 3, 3, 1, -1, 6, 0, 2, 2, 5, 5, 5, 8, -1, -1, 16, 8, 0, 4, 4, 4, 4, 4, 4, 10, 10, 10, 10, 15}

func NewPaul2013TerminalWordClassifier(table *Table) (*Paul2013TerminalWordClassifier, error) {
	if table == nil || table.Name != "sbdw_sort.txt2" || table.Mode != TwoColumns {
		return nil, fmt.Errorf("terminal word classifier requires two-column sbdw_sort.txt2")
	}
	keys := make([][]byte, len(table.Rows))
	classes := make([]int16, len(table.Rows))
	var attributes [256]byte
	for value := byte('0'); value <= '9'; value++ {
		attributes[value] = 4
	}
	for _, value := range []byte(" \t\r\n\v\f") {
		attributes[value] = 8
	}
	for index, row := range table.Rows {
		keys[index] = []byte(row.Key)
		value, _ := ParsePaul2013NativeInteger([]byte(row.Value), attributes)
		if value < 0 || int(value) >= len(paul2013TerminalWordClasses) {
			return nil, fmt.Errorf("terminal word row %d class index %d outside recovered table", index, value)
		}
		classes[index] = paul2013TerminalWordClasses[value]
	}
	search, err := NewPaul2013NativeStringSearchTable(keys, 'I', Paul2013ContextCharacterWeights())
	if err != nil {
		return nil, err
	}
	return &Paul2013TerminalWordClassifier{search: search, classes: classes}, nil
}

func (classifier *Paul2013TerminalWordClassifier) Classify(source []byte) (int16, error) {
	if classifier == nil {
		return 0, fmt.Errorf("terminal word classifier is unavailable")
	}
	source = cString(source)
	index, found, err := classifier.search.Lookup(source)
	if err != nil {
		return 0, err
	}
	if found {
		return classifier.classes[index], nil
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	digits, upper := 0, 0
	for _, value := range source {
		if attributes[value]&0x10 != 0 {
			digits++
		}
		if attributes[value]&0x80 != 0 {
			upper++
		}
	}
	if digits > 0 {
		return 9, nil
	}
	if upper >= 2 {
		return 3, nil
	}
	first := byte(0)
	if len(source) > 0 {
		first = source[0]
	}
	for _, gate := range []struct {
		characters string
		value      int16
	}{{".?!;", 17}, {",", 11}, {"'`\"", 15}, {":-", 12}, {"([{<", 13}, {")]}>", 14}} {
		if first != 0 && bytes.IndexByte([]byte(gate.characters), first) >= 0 {
			return gate.value, nil
		}
	}
	letters, _ := Paul2013CStringHasCharacterMask(source, 0xc0)
	weights := Paul2013ContextCharacterWeights()
	equal := func(value []byte, want string) bool {
		return ComparePaul2013MappedCString(value, []byte(want), weights) == 0
	}
	if letters {
		for _, suffix := range []string{"liness", "ance", "tion", "ness", "ment", "sion", "ship", "dity", "ist"} {
			if len(source) > len(suffix) && equal(source[len(source)-len(suffix):], suffix) {
				return 1, nil
			}
		}
		for _, suffix := range []string{"less", "ing", "ful", "ly"} {
			if len(source) > len(suffix) && equal(source[len(source)-len(suffix):], suffix) {
				return 5, nil
			}
		}
		return 0, nil
	}
	apostrophe := bytes.IndexByte(source, '\'')
	if apostrophe >= 0 {
		prefixLength := -1
		if apostrophe > 0 && equal(source[apostrophe-1:], "n't") {
			prefixLength = apostrophe - 1
		} else {
			for _, suffix := range []string{"s", "d", "m", "em", "ve", "re", "ll"} {
				if equal(source[apostrophe+1:], suffix) {
					prefixLength = apostrophe
					break
				}
			}
		}
		if prefixLength <= 0 {
			return 0, nil
		}
		// The native recursive prefix is copied into a 32-byte local buffer.
		if prefixLength >= 32 {
			return 0, fmt.Errorf("terminal contraction prefix exceeds native local buffer")
		}
		return classifier.Classify(source[:prefixLength])
	}
	if attributes[first]&0xc0 != 0 {
		return 0, nil
	}
	return 16, nil
}
