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
	Duration map[string]*Tree
	Pitch    map[string]*Tree
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

// Parse validates the complete header, node table, references, and output table.
func Parse(data []byte) (*Tree, error) {
	if len(data) < 7 {
		return nil, errors.New("fewer than 7 bytes remain for tree header")
	}
	nodeCount := int(int16(binary.LittleEndian.Uint16(data[0:2])))
	outputWidth := data[2]
	declaredLists := binary.LittleEndian.Uint32(data[3:7])
	if nodeCount < 0 {
		return nil, fmt.Errorf("negative node count %d", nodeCount)
	}
	if outputWidth == 0 {
		return nil, errors.New("output width must be positive")
	}
	tree := &Tree{
		OutputWidth:        outputWidth,
		Nodes:              make([]Node, 0, nodeCount),
		DeclaredListValues: declaredLists,
	}
	offset := 7
	var listCount uint64
	for index := 0; index < nodeCount; index++ {
		if len(data)-offset < 5 {
			return nil, fmt.Errorf("truncated node %d header at byte %d", index, offset)
		}
		feature := data[offset]
		operation := data[offset+1]
		threshold := int16(binary.LittleEndian.Uint16(data[offset+2 : offset+4]))
		valueCount := int(data[offset+4])
		offset += 5
		recordBytes := valueCount*2 + 4
		if len(data)-offset < recordBytes {
			return nil, fmt.Errorf("truncated node %d values or child references at byte %d", index, offset)
		}
		if operation != 'C' && operation != 'D' {
			return nil, fmt.Errorf("node %d has unknown operation byte 0x%02x", index, operation)
		}
		if operation == 'C' && valueCount != 0 {
			return nil, fmt.Errorf("node %d comparison operation has %d list values", index, valueCount)
		}
		values := make([]int16, valueCount)
		for valueIndex := range values {
			pos := offset + valueIndex*2
			values[valueIndex] = int16(binary.LittleEndian.Uint16(data[pos : pos+2]))
		}
		offset += valueCount * 2
		whenTrue := int16(binary.LittleEndian.Uint16(data[offset : offset+2]))
		whenFalse := int16(binary.LittleEndian.Uint16(data[offset+2 : offset+4]))
		offset += 4
		tree.Nodes = append(tree.Nodes, Node{
			Feature: feature, Operation: operation, Threshold: threshold,
			Values: values, WhenTrue: whenTrue, WhenFalse: whenFalse,
		})
		listCount += uint64(valueCount)
	}
	if listCount != uint64(declaredLists) {
		return nil, fmt.Errorf("header declares %d list values, nodes contain %d", declaredLists, listCount)
	}
	for index, node := range tree.Nodes {
		for _, child := range [...]int16{node.WhenTrue, node.WhenFalse} {
			if child >= 0 && int(child) >= nodeCount {
				return nil, fmt.Errorf("node %d has out-of-range child node %d", index, child)
			}
			if child < 0 && -int(child)-1 > nodeCount {
				return nil, fmt.Errorf("node %d has out-of-range leaf reference %d", index, child)
			}
		}
	}
	outputEntries := int64(nodeCount+1) * int64(outputWidth)
	outputBytes := outputEntries * 2
	if outputBytes != int64(len(data)-offset) {
		return nil, fmt.Errorf("expected %d output bytes at offset %d, found %d", outputBytes, offset, len(data)-offset)
	}
	tree.Outputs = make([][]int16, nodeCount+1)
	for row := range tree.Outputs {
		tree.Outputs[row] = make([]int16, outputWidth)
		for column := range tree.Outputs[row] {
			position := offset + (row*int(outputWidth)+column)*2
			tree.Outputs[row][column] = int16(binary.LittleEndian.Uint16(data[position : position+2]))
		}
	}
	return tree, nil
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
