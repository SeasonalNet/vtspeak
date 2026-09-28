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
