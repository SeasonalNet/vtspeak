// Package tree3 parses and evaluates the observed version-2013 decision trees.
package tree3

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Node struct {
	Feature   uint8
	Operation byte
	Threshold int16
	Values    []int16
	WhenTrue  int16
	WhenFalse int16
}

type Tree struct {
	OutputWidth        uint8
	Nodes              []Node
	Outputs            [][]int16
	DeclaredListValues uint32
}

type Catalog struct {
	Duration      map[string]*Tree
	Pitch         map[string]*Tree
	Pronunciation *Tree
}

var paul2013Trees = map[string][2]int{
	"duration/caff.tree3":   {48, 1},
	"duration/capp.tree3":   {315, 1},
	"duration/cfri.tree3":   {248, 1},
	"duration/cnas.tree3":   {352, 1},
	"duration/cstop.tree3":  {717, 1},
	"duration/vdi.tree3":    {433, 1},
	"duration/vlong.tree3":  {369, 1},
	"duration/vsch.tree3":   {101, 1},
	"duration/vshort.tree3": {812, 1},
	"pitch/bf.tree3":        {18, 12},
	"pitch/bt.tree3":        {22, 1},
	"pitch/nbf.tree3":       {263, 12},
	"pitch/nbt.tree3":       {167, 1},
	"pitch/qbf.tree3":       {9, 12},
	"pitch/qbt.tree3":       {7, 1},
	"pitch/sbf.tree3":       {42, 12},
	"pitch/sbt.tree3":       {26, 1},
}

const paul2013ATMTTreeCount = 27

// LoadPaul2013 loads the nine duration and eight pitch trees from a Paul M16
// data root and checks each against its runtime-observed node/output shape.
func LoadPaul2013(dataRoot string) (*Catalog, error) {
	catalog := &Catalog{
		Duration: make(map[string]*Tree, 9),
		Pitch:    make(map[string]*Tree, 8),
	}
	paths := make([]string, 0, len(paul2013Trees))
	for relativePath := range paul2013Trees {
		paths = append(paths, relativePath)
	}
	sort.Strings(paths)
	for _, relativePath := range paths {
		expected := paul2013Trees[relativePath]
		path := filepath.Join(dataRoot, "ttsdata", "tree3", filepath.FromSlash(relativePath))
		tree, err := ParseFile(path)
		if err != nil {
			return nil, err
		}
		shape := [2]int{len(tree.Nodes), int(tree.OutputWidth)}
		if shape != expected {
			return nil, fmt.Errorf("%s: tree shape %v does not match observed shape %v", path, shape, expected)
		}
		parts := strings.SplitN(relativePath, "/", 2)
		if parts[0] == "duration" {
			catalog.Duration[filepath.Base(path)] = tree
		} else {
			catalog.Pitch[filepath.Base(path)] = tree
		}
	}
	return catalog, nil
}

// LoadPaul2013PronunciationTree loads the shared English ambiguity classifier
// used by the Paul 2013 engine. It lives in data-common/dict-eng beside the
// embedded dictionary, rather than in the voice-specific tree3 directory.
// The observed resource has 531 nodes and emits one byte per evaluation.
func LoadPaul2013PronunciationTree(dictionaryRoot string) (*Tree, error) {
	path := filepath.Join(dictionaryRoot, "engbi.tree3")
	tree, err := ParseFile(path)
	if err != nil {
		return nil, err
	}
	if len(tree.Nodes) != 531 || tree.OutputWidth != 1 {
		return nil, fmt.Errorf("%s: pronunciation tree shape is %d nodes by %d outputs, want 531 by 1",
			path, len(tree.Nodes), tree.OutputWidth)
	}
	return tree, nil
}

// LoadPaul2013ATMTTrees loads the 27 concatenated shared English trees from
// dict-eng/atmt.tree3. The container count and exact EOF are validated.
func LoadPaul2013ATMTTrees(dictionaryRoot string) ([]*Tree, error) {
	path := filepath.Join(dictionaryRoot, "atmt.tree3")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trees, err := ParseATMT(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return trees, nil
}

// ParseATMT parses the shared English tree container used by the 2013 text
// pipeline. Its 32-bit count is observed as 27 in the local package and in the
// native loader; all declared trees must parse and consume the file exactly.
func ParseATMT(data []byte) ([]*Tree, error) {
	if len(data) < 4 {
		return nil, errors.New("ATMT container is missing its 32-bit tree count")
	}
	count := binary.LittleEndian.Uint32(data[:4])
	if count != paul2013ATMTTreeCount {
		return nil, fmt.Errorf("ATMT container declares %d trees, want %d", count, paul2013ATMTTreeCount)
	}
	trees := make([]*Tree, 0, count)
	offset := 4
	for index := uint32(0); index < count; index++ {
		tree, end, err := ParseAt(data, offset)
		if err != nil {
			return nil, fmt.Errorf("ATMT tree %d at byte %d: %w", index, offset, err)
		}
		trees = append(trees, tree)
		offset = end
	}
	if offset != len(data) {
		return nil, fmt.Errorf("ATMT trees end at byte %d, found %d trailing bytes", offset, len(data)-offset)
	}
	return trees, nil
}

// ParseFile reads one complete tree3 file.
func ParseFile(path string) (*Tree, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tree, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return tree, nil
}

// Parse validates one complete tree and rejects trailing bytes.
func Parse(data []byte) (*Tree, error) {
	tree, end, err := ParseAt(data, 0)
	if err != nil {
		return nil, err
	}
	if end != len(data) {
		return nil, fmt.Errorf("tree ends at byte %d, found %d trailing bytes", end, len(data)-end)
	}
	return tree, nil
}

// ParseAt validates one tree beginning at offset and returns the absolute byte
// offset immediately after its output table. Bytes after that offset are left
// for a containing resource parser.
func ParseAt(data []byte, offset int) (*Tree, int, error) {
	if offset < 0 || offset > len(data) {
		return nil, 0, fmt.Errorf("tree offset %d is outside the %d-byte buffer", offset, len(data))
	}
	if len(data)-offset < 7 {
		return nil, 0, fmt.Errorf("fewer than 7 bytes remain for tree header at byte %d", offset)
	}
	nodeCount := int(int16(binary.LittleEndian.Uint16(data[offset : offset+2])))
	outputWidth := data[offset+2]
	declaredLists := binary.LittleEndian.Uint32(data[offset+3 : offset+7])
	if nodeCount < 0 {
		return nil, 0, fmt.Errorf("negative node count %d", nodeCount)
	}
	if outputWidth == 0 {
		return nil, 0, errors.New("output width must be positive")
	}
	tree := &Tree{
		OutputWidth:        outputWidth,
		Nodes:              make([]Node, 0, nodeCount),
		DeclaredListValues: declaredLists,
	}
	cursor := offset + 7
	var listCount uint64
	for index := 0; index < nodeCount; index++ {
		if len(data)-cursor < 5 {
			return nil, 0, fmt.Errorf("truncated node %d header at byte %d", index, cursor)
		}
		feature := data[cursor]
		operation := data[cursor+1]
		threshold := int16(binary.LittleEndian.Uint16(data[cursor+2 : cursor+4]))
		valueCount := int(data[cursor+4])
		cursor += 5
		recordBytes := valueCount*2 + 4
		if len(data)-cursor < recordBytes {
			return nil, 0, fmt.Errorf("truncated node %d values or child references at byte %d", index, cursor)
		}
		if operation != 'C' && operation != 'D' {
			return nil, 0, fmt.Errorf("node %d has unknown operation byte 0x%02x", index, operation)
		}
		if operation == 'C' && valueCount != 0 {
			return nil, 0, fmt.Errorf("node %d comparison operation has %d list values", index, valueCount)
		}
		values := make([]int16, valueCount)
		for valueIndex := range values {
			pos := cursor + valueIndex*2
			values[valueIndex] = int16(binary.LittleEndian.Uint16(data[pos : pos+2]))
		}
		cursor += valueCount * 2
		whenTrue := int16(binary.LittleEndian.Uint16(data[cursor : cursor+2]))
		whenFalse := int16(binary.LittleEndian.Uint16(data[cursor+2 : cursor+4]))
		cursor += 4
		tree.Nodes = append(tree.Nodes, Node{
			Feature: feature, Operation: operation, Threshold: threshold,
			Values: values, WhenTrue: whenTrue, WhenFalse: whenFalse,
		})
		listCount += uint64(valueCount)
	}
	if listCount != uint64(declaredLists) {
		return nil, 0, fmt.Errorf("header declares %d list values, nodes contain %d", declaredLists, listCount)
	}
	for index, node := range tree.Nodes {
		for _, child := range [...]int16{node.WhenTrue, node.WhenFalse} {
			if child >= 0 && int(child) >= nodeCount {
				return nil, 0, fmt.Errorf("node %d has out-of-range child node %d", index, child)
			}
			if child < 0 && -int(child)-1 > nodeCount {
				return nil, 0, fmt.Errorf("node %d has out-of-range leaf reference %d", index, child)
			}
		}
	}
	outputEntries := int64(nodeCount+1) * int64(outputWidth)
	outputBytes := outputEntries * 2
	if outputBytes > int64(len(data)-cursor) {
		return nil, 0, fmt.Errorf("expected %d output bytes at offset %d, found %d", outputBytes, cursor, len(data)-cursor)
	}
	tree.Outputs = make([][]int16, nodeCount+1)
	for row := range tree.Outputs {
		tree.Outputs[row] = make([]int16, outputWidth)
		for column := range tree.Outputs[row] {
			position := cursor + (row*int(outputWidth)+column)*2
			tree.Outputs[row][column] = int16(binary.LittleEndian.Uint16(data[position : position+2]))
		}
	}
	return tree, cursor + int(outputBytes), nil
}

// Evaluate returns the reached leaf ordinal and a copy of its output row.
func (tree *Tree) Evaluate(features []int16) (int, []int16, error) {
	required := 0
	for _, node := range tree.Nodes {
		if int(node.Feature)+1 > required {
			required = int(node.Feature) + 1
		}
	}
	if len(features) < required {
		return 0, nil, fmt.Errorf("expected at least %d input features, got %d", required, len(features))
	}
	index := int16(0)
	for steps := 0; ; steps++ {
		if index < 0 {
			leaf := -int(index) - 1
			if leaf >= len(tree.Outputs) {
				return 0, nil, fmt.Errorf("leaf ordinal %d is out of range", leaf)
			}
			return leaf, append([]int16(nil), tree.Outputs[leaf]...), nil
		}
		if int(index) >= len(tree.Nodes) {
			return 0, nil, fmt.Errorf("child node index %d is out of range", index)
		}
		if steps > len(tree.Nodes) {
			return 0, nil, errors.New("cycle encountered during tree evaluation")
		}
		node := tree.Nodes[index]
		value := features[node.Feature]
		if node.Operation == 'D' {
			index = node.WhenFalse
			for _, listed := range node.Values {
				if value == listed {
					index = node.WhenTrue
					break
				}
			}
		} else if value <= node.Threshold {
			index = node.WhenTrue
		} else {
			index = node.WhenFalse
		}
	}
}
