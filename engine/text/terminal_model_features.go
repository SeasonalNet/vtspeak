package text

import (
	"bytes"
	"fmt"
)

// Paul2013TerminalAbbreviationTable is the copied two-column input to
// FUN_10052d00. Loader mode I searches by mapped weights; column 2 starts
// with '1' for an exact-case match or '2' for a mapped-case match.
type Paul2013TerminalAbbreviationTable struct {
	search *Paul2013NativeStringSearchTable
	rows   []Row
}

func NewPaul2013TerminalAbbreviationTable(table *Table) (*Paul2013TerminalAbbreviationTable, error) {
	if table == nil || table.Mode != TwoColumns || (table.Name != "abbrh_sort.txt2" && table.Name != "abbrt_sort.txt2" && table.Name != "abbrc_sort.txt2") {
		return nil, fmt.Errorf("terminal abbreviation lookup requires a recovered two-column abbreviation resource")
	}
	rows := append([]Row(nil), table.Rows...)
	keys := make([][]byte, len(rows))
	for index, row := range rows {
		keys[index] = []byte(row.Key)
	}
	search, err := NewPaul2013NativeStringSearchTable(keys, 'I', Paul2013ContextCharacterWeights())
	if err != nil {
		return nil, err
	}
	return &Paul2013TerminalAbbreviationTable{search: search, rows: rows}, nil
}

// Lookup preserves the native return index: a qualifying adjacent variant
// returns the initial binary-search index, rather than the variant's index.
func (table *Paul2013TerminalAbbreviationTable) Lookup(source []byte) (int, error) {
	if table == nil {
		return -1, fmt.Errorf("terminal abbreviation table is unavailable")
	}
	index, found, err := table.search.Lookup(source)
	if err != nil || !found {
		return -1, err
	}
	source = cString(source)
	accepted := func(row Row) bool {
		value := cString([]byte(row.Value))
		return len(value) > 0 && (value[0] == '2' || value[0] == '1' && bytes.Equal(cString([]byte(row.Key)), source))
	}
	if accepted(table.rows[index]) {
		return index, nil
	}
	weights := Paul2013ContextCharacterWeights()
	for previous := index - 1; previous >= 0 && ComparePaul2013MappedCString([]byte(table.rows[previous].Key), source, weights) == 0; previous-- {
		if accepted(table.rows[previous]) {
			return index, nil
		}
	}
	for following := index + 1; following < len(table.rows) && ComparePaul2013MappedCString([]byte(table.rows[following].Key), source, weights) == 0; following++ {
		if accepted(table.rows[following]) {
			return index, nil
		}
	}
	return -1, nil
}

// ClassifyPaul2013TerminalTokenCase ports FUN_10052710's raw case/status
// feature. The native window has 36 zero-initialized bytes at +0x20;
// a counted length beyond that storage is rejected rather than overread.
func ClassifyPaul2013TerminalTokenCase(token Paul2013ModelScannerResult) (int16, error) {
	if token.TextLength == 0 {
		return -1, nil
	}
	if token.Status == 2 {
		return 4, nil
	}
	if token.Status != 1 {
		return 5, nil
	}
	if token.TextLength > 36 {
		return 0, fmt.Errorf("terminal case feature length exceeds scanner record")
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	text := cString(token.Text)
	upper, lower, firstUpper := false, false, false
	for index := int32(0); index < token.TextLength; index++ {
		value := byte(0)
		if int(index) < len(text) {
			value = text[index]
		}
		if attributes[value]&0x80 != 0 {
			upper = true
			if index == 0 {
				firstUpper = true
			}
		} else if attributes[value]&0x40 != 0 {
			lower = true
		}
	}
	if upper && !lower {
		return 1, nil
	}
	if lower && !upper {
		return 2, nil
	}
	if firstUpper && lower {
		return 3, nil
	}
	return 5, nil
}

// ClassifyPaul2013TerminalTokenPunctuation ports FUN_100527c0's literal gate
// and jump table at 0x1005284c/0x1005286c. The single quote has a table entry
// but fails the earlier literal gate, so it still returns -1.
func ClassifyPaul2013TerminalTokenPunctuation(token Paul2013ModelScannerResult) int16 {
	if token.TextLength != 1 || len(token.Text) == 0 {
		return -1
	}
	switch token.Text[0] {
	case '.':
		return 1
	case '?':
		return 2
	case '!':
		return 3
	case ';':
		return 4
	case ',', ':':
		return 5
	case '"', '`':
		return 6
	case '-':
		return 7
	default:
		return -1
	}
}

// Paul2013TerminalTokenHasInternalDot ports FUN_10052c90. It tests only
// the first dot returned by strchr, which must have bytes on both sides.
func Paul2013TerminalTokenHasInternalDot(source []byte) int16 {
	source = cString(source)
	index := bytes.IndexByte(source, '.')
	if index > 0 && index+1 < len(source) {
		return 1
	}
	return 0
}

// Paul2013TerminalAbbreviationMask ports FUN_100526b0's three lookups.
// Bit 4 is independent here; the key writer later combines all three bits.
func Paul2013TerminalAbbreviationMask(tables [3]*Paul2013TerminalAbbreviationTable, source []byte) (uint16, error) {
	var mask uint16
	for index, table := range tables {
		result, err := table.Lookup(source)
		if err != nil {
			return 0, err
		}
		if result >= 0 {
			mask |= 1 << index
		}
	}
	return mask, nil
}
