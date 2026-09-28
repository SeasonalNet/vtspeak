// Package distance reads the raw Paul 2013 metric-distance table.
package distance

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
)

// Table stores the lower triangle of a symmetric distance matrix in row-major
// order. The model file is retained as read-only input and decoded in memory.
type Table struct {
	size   uint16
	values []float32
}

// FeatureTable stores the lower triangle of the generated 256-bin feature
// difference matrix built by FUN_1001ae70.
type FeatureTable struct {
	values []float32
}

// SymmetricFeatureDifference computes the feature-value cost used by the
// observed unit and transition scorers. The DLL stores the averaged
// denominator as float32 before the final arithmetic; a zero denominator
// produces zero cost.
func SymmetricFeatureDifference(left, right float32) float64 {
	denominator := (left + right) / 2
	if denominator == 0 {
		return 0
	}
	difference := float64(left) - float64(right)
	return difference * difference / float64(denominator)
}

// NewFeatureTable constructs the deterministic table built by the Paul 2013
// DLL. Entries are the float32-stored results of FUN_1001ae20 for all pairs
// of integer feature-bin values from 0 through 255.
func NewFeatureTable() *FeatureTable {
	const size = 256
	values := make([]float32, size*(size+1)/2)
	for row := 0; row < size; row++ {
		rowStart := row * (row + 1) / 2
		for column := 0; column <= row; column++ {
			values[rowStart+column] = float32(SymmetricFeatureDifference(float32(row), float32(column)))
		}
	}
	return &FeatureTable{values: values}
}

// Lookup returns the symmetric cost for two observed 8-bit feature codes.
func (table *FeatureTable) Lookup(left, right byte) (float32, error) {
	if table == nil || len(table.values) != 256*257/2 {
		return 0, errors.New("feature-distance table is not initialized")
	}
	row, column := int(left), int(right)
	if row < column {
		row, column = column, row
	}
	return table.values[row*(row+1)/2+column], nil
}

// LoadFile reads one complete little-endian distance table.
func LoadFile(path string) (*Table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read distance table: %w", err)
	}
	table, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse distance table: %w", err)
	}
	return table, nil
}

// Parse validates the 16-bit entry count and complete triangular float32
// payload observed in Paul 2013's cepdist.tbl.
func Parse(data []byte) (*Table, error) {
	if len(data) < 2 {
		return nil, errors.New("distance table is shorter than its count header")
	}
	size := binary.LittleEndian.Uint16(data[:2])
	if size == 0 {
		return nil, errors.New("distance table count must be positive")
	}
	count := uint64(size) * (uint64(size) + 1) / 2
	expectedSize := uint64(2) + count*4
	if expectedSize != uint64(len(data)) {
		return nil, fmt.Errorf("distance table has %d bytes, want %d for %d entries", len(data), expectedSize, count)
	}
	values := make([]float32, int(count))
	for index := range values {
		bits := binary.LittleEndian.Uint32(data[2+index*4 : 6+index*4])
		value := math.Float32frombits(bits)
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) || value < 0 {
			return nil, fmt.Errorf("distance table value %d is not a finite nonnegative float", index)
		}
		values[index] = value
	}
	return &Table{size: size, values: values}, nil
}

// Size returns the number of metric codes represented by the matrix.
func (table *Table) Size() uint16 {
	if table == nil {
		return 0
	}
	return table.size
}

// Lookup returns the symmetric distance for two metric codes. Callers must
// apply any model-specific packed-code masks before lookup.
func (table *Table) Lookup(left, right uint16) (float32, error) {
	if table == nil || table.size == 0 {
		return 0, errors.New("distance table is not loaded")
	}
	if left >= table.size || right >= table.size {
		return 0, fmt.Errorf("metric codes %d and %d exceed table size %d", left, right, table.size)
	}
	row, column := left, right
	if row < column {
		row, column = column, row
	}
	index := uint64(row)*(uint64(row)+1)/2 + uint64(column)
	return table.values[index], nil
}

// LookupPacked applies the observed 0x3fff metric-code mask before lookup.
func (table *Table) LookupPacked(left, right uint16) (float32, error) {
	return table.Lookup(left&0x3fff, right&0x3fff)
}
