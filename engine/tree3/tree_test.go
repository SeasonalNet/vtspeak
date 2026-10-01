package tree3

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestParseAndEvaluateMembershipTree(t *testing.T) {
	data := make([]byte, 7+5+4+4+8)
	binary.LittleEndian.PutUint16(data[0:2], 1)
	data[2] = 2
	binary.LittleEndian.PutUint32(data[3:7], 2)
	data[7] = 0
	data[8] = 'D'
	data[11] = 2
	putI16(data[12:14], -1)
	binary.LittleEndian.PutUint16(data[14:16], 7)
	putI16(data[16:18], -1)
	putI16(data[18:20], -2)
	binary.LittleEndian.PutUint16(data[20:22], 100)
	binary.LittleEndian.PutUint16(data[22:24], 101)
	putI16(data[24:26], -2)
	binary.LittleEndian.PutUint16(data[26:28], 300)

	tree, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	leaf, output, err := tree.Evaluate([]int16{7})
	if err != nil || leaf != 0 || !equal(output, []int16{100, 101}) {
		t.Fatalf("membership match = (%d, %v, %v)", leaf, output, err)
	}
	leaf, output, err = tree.Evaluate([]int16{5})
	if err != nil || leaf != 1 || !equal(output, []int16{-2, 300}) {
		t.Fatalf("membership miss = (%d, %v, %v)", leaf, output, err)
	}
}

func TestParseAndEvaluateSignedComparison(t *testing.T) {
	data := make([]byte, 7+9+4)
	binary.LittleEndian.PutUint16(data[0:2], 1)
	data[2] = 1
	data[8] = 'C'
	putI16(data[9:11], -3)
	putI16(data[12:14], -1)
	putI16(data[14:16], -2)
	binary.LittleEndian.PutUint16(data[16:18], 14)
	binary.LittleEndian.PutUint16(data[18:20], 25)
	tree, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	leaf, output, err := tree.Evaluate([]int16{-4})
	if err != nil || leaf != 0 || !equal(output, []int16{14}) {
		t.Fatalf("comparison match = (%d, %v, %v)", leaf, output, err)
	}
}

func TestRejectMalformedTrees(t *testing.T) {
	valid := []byte{
		1, 0, 1, 0, 0, 0, 0,
		0, 'C', 0, 0, 0, 0xff, 0xff, 0xfe, 0xff,
		1, 0, 2, 0,
	}
	malformed := map[string][]byte{
		"short header":        valid[:6],
		"unknown operation":   append(append([]byte(nil), valid[:7]...), []byte{0, 'X', 0, 0, 0, 0xff, 0xff, 0xfe, 0xff, 1, 0, 2, 0}...),
		"invalid leaf":        append(append([]byte(nil), valid[:15]...), []byte{0xfd, 0xff, 1, 0, 2, 0}...),
		"output extent":       valid[:len(valid)-1],
		"trailing bytes":      append(append([]byte(nil), valid...), 0),
		"comparison has list": []byte{1, 0, 1, 1, 0, 0, 0, 0, 'C', 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 2, 0},
	}
	for name, data := range malformed {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse(data); err == nil {
				t.Fatal("malformed tree accepted")
			}
		})
	}
}

func TestParseAtReadsOneTreeFromConcatenatedBuffer(t *testing.T) {
	first := makeEmptyTree(1)
	second := makeEmptyTree(2)
	data := append(append([]byte(nil), first...), second...)
	tree, end, err := ParseAt(data, 0)
	if err != nil {
		t.Fatal(err)
	}
	if end != len(first) || tree.OutputWidth != 1 {
		t.Fatalf("first tree end/width = %d/%d, want %d/1", end, tree.OutputWidth, len(first))
	}
	tree, end, err = ParseAt(data, end)
	if err != nil {
		t.Fatal(err)
	}
	if end != len(data) || tree.OutputWidth != 2 {
		t.Fatalf("second tree end/width = %d/%d, want %d/2", end, tree.OutputWidth, len(data))
	}
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse accepted a concatenated tree buffer")
	}
}

func TestParseATMTValidatesCountAndEOF(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, paul2013ATMTTreeCount)
	for range paul2013ATMTTreeCount {
		data = append(data, makeEmptyTree(1)...)
	}
	trees, err := ParseATMT(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(trees) != paul2013ATMTTreeCount {
		t.Fatalf("parsed %d ATMT trees, want %d", len(trees), paul2013ATMTTreeCount)
	}
	withTrailingByte := append(append([]byte(nil), data...), 0)
	if _, err := ParseATMT(withTrailingByte); err == nil {
		t.Fatal("ATMT parser accepted trailing bytes")
	}
	wrongCount := append([]byte(nil), data...)
	binary.LittleEndian.PutUint32(wrongCount[:4], 26)
	if _, err := ParseATMT(wrongCount); err == nil {
		t.Fatal("ATMT parser accepted an unexpected tree count")
	}
}

func TestLoadPaul2013ATMTTreesFromLocalCommonData(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "atmt.tree3")); err != nil {
		t.Skip("local shared English ATMT tree container is unavailable")
	}
	trees, err := LoadPaul2013ATMTTrees(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(trees) != paul2013ATMTTreeCount {
		t.Fatalf("loaded %d ATMT trees, want %d", len(trees), paul2013ATMTTreeCount)
	}
	nodes := 0
	for _, tree := range trees {
		nodes += len(tree.Nodes)
	}
	if nodes != 13879 {
		t.Fatalf("ATMT container has %d nodes, want observed total 13879", nodes)
	}
}

func makeEmptyTree(outputWidth byte) []byte {
	data := make([]byte, 7+2*int(outputWidth))
	data[2] = outputWidth
	return data
}

func TestLoadPaul2013PronunciationTreeRejectsUnexpectedShape(t *testing.T) {
	root := t.TempDir()
	data := make([]byte, 9)
	binary.LittleEndian.PutUint16(data[0:2], 0)
	data[2] = 1
	binary.LittleEndian.PutUint16(data[7:9], 7)
	if err := os.WriteFile(filepath.Join(root, "engbi.tree3"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPaul2013PronunciationTree(root); err == nil {
		t.Fatal("unexpected pronunciation-tree shape accepted")
	}
}

func TestLoadPaul2013PronunciationTreeFromLocalCommonData(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engbi.tree3")); err != nil {
		t.Skip("local shared English tree is unavailable")
	}
	tree, err := LoadPaul2013PronunciationTree(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(tree.Nodes) != 531 || tree.OutputWidth != 1 {
		t.Fatalf("pronunciation tree shape = %d nodes by %d outputs", len(tree.Nodes), tree.OutputWidth)
	}
	_, output, err := tree.Evaluate(make([]int16, 15))
	if err != nil {
		t.Fatal(err)
	}
	if len(output) != 1 {
		t.Fatalf("pronunciation output width = %d, want 1", len(output))
	}
}

func TestAllPaulTreesAndRuntimeLookups(t *testing.T) {
	root := filepath.Join("..", "..")
	treeDir := filepath.Join(root, "data-paul", "M16", "ttsdata", "tree3")
	if _, err := os.Stat(treeDir); err != nil {
		t.Skip("local Paul tree files are unavailable")
	}
	catalog, err := LoadPaul2013(filepath.Join(root, "data-paul", "M16"))
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.Duration) != 9 || len(catalog.Pitch) != 8 {
		t.Fatalf("loaded %d duration and %d pitch trees", len(catalog.Duration), len(catalog.Pitch))
	}
	byShape := make(map[[2]int]*Tree, len(paul2013Trees))
	for _, group := range []map[string]*Tree{catalog.Duration, catalog.Pitch} {
		for _, tree := range group {
			shape := [2]int{len(tree.Nodes), int(tree.OutputWidth)}
			if _, exists := byShape[shape]; exists {
				t.Fatalf("tree shape %v is ambiguous", shape)
			}
			byShape[shape] = tree
		}
	}

	tracePath := filepath.Join(root, "tools", "revkit", "work", "stage5", "gdb-tree-run.log")
	trace, err := os.ReadFile(tracePath)
	if err != nil {
		t.Skip("tracked original-engine tree trace is unavailable")
	}
	count, err := compareRuntimeLookups(string(trace), byShape)
	if err != nil {
		t.Fatal(err)
	}
	if count != 103 {
		t.Fatalf("compared %d lookups from the Stage 5 trace, expected 103", count)
	}
}

var (
	entryPattern  = regexp.MustCompile(`(?:SCALAR|VECTOR)_ENTRY .* nodes=(\d+) output_width=(\d+) features=`)
	valuePattern  = regexp.MustCompile(`-?\d+`)
	scalarPattern = regexp.MustCompile(`SCALAR_RETURN value=(-?\d+)`)
)

func compareRuntimeLookups(log string, trees map[[2]int]*Tree) (int, error) {
	lines := strings.Split(log, "\n")
	compared := 0
	for index, line := range lines {
		entry := entryPattern.FindStringSubmatch(line)
		if entry == nil || index+2 >= len(lines) {
			continue
		}
		nodes, _ := strconv.Atoi(entry[1])
		width, _ := strconv.Atoi(entry[2])
		tree := trees[[2]int{nodes, width}]
		if tree == nil {
			continue
		}
		features := append(numbersAfterColon(lines[index+1]), numbersAfterColon(lines[index+2])...)
		leaf, expected, err := tree.Evaluate(features)
		if err != nil {
			return compared, fmt.Errorf("runtime entry %d: %w", index, err)
		}
		if strings.HasPrefix(line, "SCALAR_ENTRY") {
			var actual []int16
			for next := index + 3; next < len(lines) && next <= index+8; next++ {
				match := scalarPattern.FindStringSubmatch(lines[next])
				if match == nil {
					continue
				}
				value, parseErr := strconv.ParseInt(match[1], 10, 16)
				if parseErr != nil {
					return compared, parseErr
				}
				actual = []int16{int16(value)}
				break
			}
			if !equal(actual, expected) {
				return compared, fmt.Errorf("scalar runtime output at line %d: leaf %d expected %v, got %v", index+1, leaf, expected, actual)
			}
		} else {
			var actual []int16
			resultIndex := -1
			for next := index + 3; next < len(lines) && next <= index+9; next++ {
				if strings.HasPrefix(lines[next], "VECTOR_RETURN values:") {
					resultIndex = next
					actual = append(actual, numbersAfterColon(lines[next])...)
					break
				}
			}
			if resultIndex >= 0 && resultIndex+1 < len(lines) {
				actual = append(actual, numbersAfterColon(lines[resultIndex+1])...)
			}
			if !equal(actual, expected) {
				return compared, fmt.Errorf("vector runtime output at line %d: leaf %d expected %v, got %v", index+1, leaf, expected, actual)
			}
		}
		compared++
	}
	return compared, nil
}

func numbersAfterColon(line string) []int16 {
	if separator := strings.LastIndexByte(line, ':'); separator >= 0 {
		line = line[separator+1:]
	}
	items := valuePattern.FindAllString(line, -1)
	values := make([]int16, 0, len(items))
	for _, item := range items {
		value, err := strconv.ParseInt(item, 10, 16)
		if err != nil {
			return nil
		}
		values = append(values, int16(value))
	}
	return values
}

func equal(left, right []int16) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func putI16(dst []byte, value int16) {
	binary.LittleEndian.PutUint16(dst, uint16(value))
}
