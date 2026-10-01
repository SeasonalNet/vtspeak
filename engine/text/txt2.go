package text

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TableMode uint8

const (
	OneColumn TableMode = iota + 1
	TwoColumns
)

type Row struct {
	Key   string
	Value string
}

type Table struct {
	Name string
	Mode TableMode
	Rows []Row
}

type Tables map[string]*Table

var txt2Modes = map[string]TableMode{
	"abbrc_sort.txt2":   TwoColumns,
	"abbrh_sort.txt2":   TwoColumns,
	"abbrt_sort.txt2":   TwoColumns,
	"chc_sort.txt2":     TwoColumns,
	"citya_sort.txt2":   TwoColumns,
	"sbdw_sort.txt2":    TwoColumns,
	"streeta_sort.txt2": TwoColumns,
	"streetf_sort.txt2": OneColumn,
	"wab.txt2":          OneColumn,
}

var txt2Names = [...]string{
	"abbrc_sort.txt2", "abbrh_sort.txt2", "abbrt_sort.txt2", "chc_sort.txt2",
	"citya_sort.txt2", "sbdw_sort.txt2", "streeta_sort.txt2", "streetf_sort.txt2", "wab.txt2",
}

// ParseTXT2 decodes one callsite-mapped shared resource and validates its
// framing, row count, ASCII body, and documented column count.
func ParseTXT2(name string, raw []byte) (*Table, error) {
	mode, ok := txt2Modes[filepath.Base(name)]
	if !ok {
		return nil, fmt.Errorf("no callsite column mode is recorded for %q", name)
	}
	if len(raw) < 17 {
		return nil, errors.New("TXT2 resource is shorter than its 17-byte framing")
	}
	if string(raw[:3]) != string(raw[len(raw)-3:]) {
		return nil, errors.New("TXT2 start and end markers differ")
	}
	shift := raw[len(raw)-14]
	bodyBytes := make([]byte, len(raw)-17)
	for i, value := range raw[3 : len(raw)-14] {
		bodyBytes[i] = value - shift
	}
	countBytes := make([]byte, 10)
	for i, value := range raw[len(raw)-13 : len(raw)-3] {
		countBytes[i] = value - shift
	}
	countEnd := 0
	for countEnd < len(countBytes) && countBytes[countEnd] != 0 {
		countEnd++
	}
	if countEnd == 0 || countEnd == len(countBytes) {
		return nil, errors.New("TXT2 row count is not NUL-terminated decimal text")
	}
	for _, value := range countBytes[countEnd+1:] {
		if value != 0 {
			return nil, errors.New("TXT2 row count has nonzero padding")
		}
	}
	for _, value := range countBytes[:countEnd] {
		if value < '0' || value > '9' {
			return nil, errors.New("TXT2 row count is not decimal text")
		}
	}
	declared, err := strconv.Atoi(string(countBytes[:countEnd]))
	if err != nil {
		return nil, fmt.Errorf("parse TXT2 row count: %w", err)
	}
	for _, value := range bodyBytes {
		if value > 0x7f {
			return nil, errors.New("TXT2 decoded body is not ASCII")
		}
	}
	body := string(bodyBytes)
	lines := strings.FieldsFunc(body, func(r rune) bool { return r == '\r' || r == '\n' })
	if len(lines) != declared {
		return nil, fmt.Errorf("TXT2 trailer declares %d rows but body contains %d", declared, len(lines))
	}
	table := &Table{Name: filepath.Base(name), Mode: mode, Rows: make([]Row, 0, len(lines))}
	for rowIndex, line := range lines {
		line = strings.TrimRight(line, " \t")
		columns := strings.Split(line, "|")
		if mode == OneColumn && len(columns) != 1 {
			return nil, fmt.Errorf("TXT2 row %d has %d columns, expected one", rowIndex, len(columns))
		}
		if mode == TwoColumns && len(columns) != 2 {
			return nil, fmt.Errorf("TXT2 row %d has %d columns, expected two", rowIndex, len(columns))
		}
		row := Row{Key: columns[0]}
		if mode == TwoColumns {
			row.Value = columns[1]
		}
		table.Rows = append(table.Rows, row)
	}
	return table, nil
}

// FindExact returns every row whose first column exactly matches key,
// preserving source order and case variants.
func (table *Table) FindExact(key string) []Row {
	var matches []Row
	for _, row := range table.Rows {
		if row.Key == key {
			matches = append(matches, row)
		}
	}
	return matches
}

// LookupPaul2013WABClass returns the zero-based row index used by the Paul
// 2013 WAB classifier. FUN_100035d0 binary-searches this one-column table
// using FUN_1001c2c0 when the loader's mode byte is 'I'.
func (table *Table) LookupPaul2013WABClass(surface []byte, characterMap [256]byte) (int, bool, error) {
	if table == nil {
		return 0, false, errors.New("WAB table is nil")
	}
	if table.Name != "wab.txt2" || table.Mode != OneColumn {
		return 0, false, fmt.Errorf("WAB lookup requires the one-column wab.txt2 table, got %q mode %d", table.Name, table.Mode)
	}
	for rowIndex := 1; rowIndex < len(table.Rows); rowIndex++ {
		if paul2013MappedCStringCompare([]byte(table.Rows[rowIndex-1].Key), []byte(table.Rows[rowIndex].Key), characterMap) > 0 {
			return 0, false, fmt.Errorf("WAB table rows %d and %d are not sorted by the native character map", rowIndex-1, rowIndex)
		}
	}
	key := cString(surface)
	if len(key) == 0 {
		return 0, false, nil
	}
	low, high := 0, len(table.Rows)-1
	for low <= high {
		middle := low + (high-low)/2
		comparison := paul2013MappedCStringCompare([]byte(table.Rows[middle].Key), key, characterMap)
		if comparison == 0 {
			return middle, true, nil
		}
		if comparison < 0 {
			low = middle + 1
		} else {
			high = middle - 1
		}
	}
	return 0, false, nil
}

// LookupPaul2013CHCFlags ports FUN_10002680's exact-key search in the
// two-column chc_sort.txt2 table. The second column's first four bytes map
// '1' values to mask bits 8, 4, 2, and 1; all requested bits must be present.
// It returns the native zero-based row index when the mask matches.
func (table *Table) LookupPaul2013CHCFlags(key []byte, requiredMask uint16) (int, bool, error) {
	if table == nil {
		return 0, false, errors.New("CHC table is nil")
	}
	if table.Name != "chc_sort.txt2" || table.Mode != TwoColumns {
		return 0, false, fmt.Errorf("CHC lookup requires the two-column chc_sort.txt2 table, got %q mode %d", table.Name, table.Mode)
	}
	for rowIndex := 1; rowIndex < len(table.Rows); rowIndex++ {
		if table.Rows[rowIndex-1].Key >= table.Rows[rowIndex].Key {
			return 0, false, fmt.Errorf("CHC table rows %d and %d are not strictly sorted", rowIndex-1, rowIndex)
		}
	}
	key = cString(key)
	if len(key) == 0 {
		return 0, false, nil
	}
	wanted := string(key)
	low, high := 0, len(table.Rows)-1
	for low <= high {
		middle := low + (high-low)/2
		row := table.Rows[middle]
		comparison := strings.Compare(row.Key, wanted)
		if comparison == 0 {
			if len(row.Value) < 4 {
				return 0, false, fmt.Errorf("CHC table row %d has a %d-byte mask, need at least four", middle, len(row.Value))
			}
			var flags byte
			for index, value := range []byte(row.Value[:4]) {
				if value == '1' {
					flags |= 1 << (3 - index)
				} else if value != '0' {
					return 0, false, fmt.Errorf("CHC table row %d mask byte %d is %q, want '0' or '1'", middle, index, value)
				}
			}
			return middle, uint16(flags)&requiredMask == requiredMask, nil
		}
		if comparison < 0 {
			low = middle + 1
		} else {
			high = middle - 1
		}
	}
	return 0, false, nil
}

// MatchPaul2013CHCSubstring ports FUN_100026f0's exact CHC lookup, compound
// prefix/suffix checks, trailing-S retry, and long-token fallback. word is
// the candidate substring; tokenLength is the full source token length.
func (table *Table) MatchPaul2013CHCSubstring(
	word []byte,
	tokenLength int,
	requiredMask uint16,
	characterMap [256]byte,
) (bool, error) {
	word = cString(word)
	length := len(word)
	if length == 0 || length > 8 || tokenLength < length {
		return false, nil
	}
	if _, found, err := table.LookupPaul2013CHCFlags(word, requiredMask); err != nil || found {
		return found, err
	}

	trySplit := func(prefixLength int) (bool, error) {
		prefix := word[:prefixLength]
		suffix := word[prefixLength:]
		if _, found, err := table.LookupPaul2013CHCFlags(prefix, 1); err != nil || !found {
			return false, err
		}
		_, found, err := table.LookupPaul2013CHCFlags(suffix, 4)
		return found, err
	}

	if length < 3 || requiredMask != 1 {
		var splits []int
		if requiredMask == 2 {
			switch length {
			case 4:
				splits = []int{2}
			case 5:
				splits = []int{2, 3}
			case 6:
				splits = []int{3, 2, 4}
			case 7:
				splits = []int{2, 3, 4, 5}
			case 8:
				splits = []int{2, 3, 4, 5, 6}
			}
		} else if length == 8 {
			splits = []int{2, 3, 4, 5, 6}
		}
		for _, prefixLength := range splits {
			found, err := trySplit(prefixLength)
			if err != nil || found {
				return found, err
			}
		}
	} else if characterMap[word[length-1]] == 's' {
		if _, found, err := table.LookupPaul2013CHCFlags(word[:length-1], 1); err != nil || found {
			return found, err
		}
	}

	if tokenLength > 6 && ((length < 5 && requiredMask != 8) || (length >= 5 && requiredMask == 2)) {
		return true, nil
	}
	return false, nil
}

func paul2013MappedCStringCompare(left, right []byte, characterMap [256]byte) int {
	left = cString(left)
	right = cString(right)
	for index := 0; ; index++ {
		leftByte, rightByte := byte(0), byte(0)
		if index < len(left) {
			leftByte = left[index]
		}
		if index < len(right) {
			rightByte = right[index]
		}
		leftMapped, rightMapped := characterMap[leftByte], characterMap[rightByte]
		if leftMapped < rightMapped {
			return -1
		}
		if leftMapped > rightMapped {
			return 1
		}
		if leftByte == 0 {
			return 0
		}
	}
}

// LoadTXT2Tables loads the nine shared Paul text resources whose callsite
// column modes are documented in the reverse-engineering findings.
func LoadTXT2Tables(root string) (Tables, error) {
	tables := make(Tables, len(txt2Modes))
	for _, name := range txt2Names {
		path := filepath.Join(root, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		table, err := ParseTXT2(name, raw)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		tables[name] = table
	}
	return tables, nil
}
