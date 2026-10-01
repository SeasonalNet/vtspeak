package dat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	indexHeaderSize = 45
	unitStride      = 19
	featureStride   = 21
)

var paul2013Header = []byte("ver.2013\x00VoiceText-Eng\x00")

// BankName returns the single bank identifier stored in a 2013 Paul index.
func BankName(index []byte) (string, error) {
	if len(index) < indexHeaderSize || int(index[0]) != len(paul2013Header) || !bytes.Equal(index[1:24], paul2013Header) {
		return "", errors.New("unsupported 2013 Paul unit index")
	}
	if binary.LittleEndian.Uint16(index[24:26]) != 1 {
		return "", errors.New("2013 Paul index must identify exactly one bank")
	}
	nameLength := int(binary.LittleEndian.Uint16(index[26:28]))
	if nameLength == 0 || 28+nameLength+1 != 39 || index[38] != 0 {
		return "", errors.New("invalid 2013 Paul bank table")
	}
	name := string(index[28 : 28+nameLength])
	if !strings.HasPrefix(name, "merged-") {
		return "", errors.New("unsupported 2013 Paul bank name")
	}
	return name, nil
}

// UnitRecord contains the observed payload span fields and one logical
// 21-byte feature row reassembled from the index's column-major arrays.
// MetricCodes and FeatureCodes expose the three opaque four-byte groups in
// their original order without assigning semantic names to their values.
type UnitRecord struct {
	DATOffset         uint32
	FirstSideSamples  uint16
	SecondSideSamples uint16
	DATLength         uint16
	UPMOffset         uint32
	UPMFirstCount     byte
	UPMSecondCount    byte
	UPMEdges          [3]byte
	Signature         [7]byte
	Features          [featureStride]byte
	MetricCodes       [3]uint16
	FeatureCodes      [3][2]byte
}

// ReadUnit reads the 2013 Paul record and its columnar 21-byte feature row.
func ReadUnit(index []byte, unit uint32) (UnitRecord, error) {
	var result UnitRecord
	if _, err := BankName(index); err != nil {
		return result, err
	}
	count := binary.LittleEndian.Uint32(index[39:43])
	stride := binary.LittleEndian.Uint16(index[43:45])
	if stride != unitStride || unit >= count {
		return result, fmt.Errorf("unit %d outside supported index range", unit)
	}
	expectedSize := uint64(indexHeaderSize) + uint64(count)*uint64(unitStride+featureStride)
	if expectedSize != uint64(len(index)) {
		return result, errors.New("unit index size does not match its declared layout")
	}
	recordStart := uint64(indexHeaderSize) + uint64(unit)*unitStride
	record := index[recordStart : recordStart+unitStride]
	result.DATOffset = binary.LittleEndian.Uint32(record[:4])
	result.FirstSideSamples = binary.LittleEndian.Uint16(record[4:6])
	result.SecondSideSamples = binary.LittleEndian.Uint16(record[6:8])
	result.DATLength = binary.LittleEndian.Uint16(record[8:10])
	result.UPMOffset = binary.LittleEndian.Uint32(record[10:14])
	result.UPMFirstCount = record[14]
	result.UPMSecondCount = record[15]
	copy(result.UPMEdges[:], record[16:19])
	// FUN_10019940 reads the feature section as a sequence of columns, not
	// 21-byte per-unit rows. Reassemble one logical row from those arrays.
	columnsStart := uint64(indexHeaderSize) + uint64(count)*unitStride
	column := func(offset, width uint64) []byte {
		start := columnsStart + offset*uint64(count) + uint64(unit)*width
		return index[start : start+width]
	}
	copy(result.Features[0:1], column(0, 1))
	copy(result.Signature[:], column(1, 7))
	copy(result.Features[1:8], result.Signature[:])
	copy(result.Features[8:9], column(8, 1))
	for group := uint64(0); group < 3; group++ {
		columnBase := uint64(9) + group*4
		rowBase := 9 + group*4
		metric := column(columnBase, 2)
		result.MetricCodes[group] = binary.LittleEndian.Uint16(metric)
		copy(result.Features[rowBase:rowBase+2], metric)
		copy(result.Features[rowBase+2:rowBase+3], column(columnBase+2, 1))
		copy(result.Features[rowBase+3:rowBase+4], column(columnBase+3, 1))
		result.FeatureCodes[group] = [2]byte{
			result.Features[rowBase+2],
			result.Features[rowBase+3],
		}
	}
	return result, nil
}

// UnitPayload reads a single 2013 Paul DAT span identified by its unit index.
// The caller supplies the complete index and a matching DAT bank reader.
func UnitPayload(index []byte, bank io.ReaderAt, bankSize int64, unit uint32) ([]byte, error) {
	record, err := ReadUnit(index, unit)
	if err != nil {
		return nil, err
	}
	offset := int64(record.DATOffset)
	length := int64(record.DATLength)
	if length == 0 || offset+length > bankSize {
		return nil, errors.New("unit DAT span outside bank")
	}
	payload := make([]byte, length)
	if _, err := bank.ReadAt(payload, offset); err != nil {
		return nil, fmt.Errorf("read unit DAT span: %w", err)
	}
	return payload, nil
}

// UPMPayload reads the combined first/second UPM span. The sides share one
// boundary byte, so their combined length is first + second - 1.
func UPMPayload(record UnitRecord, bank io.ReaderAt, bankSize int64) ([]byte, error) {
	if record.UPMFirstCount == 0 || record.UPMSecondCount == 0 {
		return nil, errors.New("unit UPM side count is zero")
	}
	length := int(record.UPMFirstCount) + int(record.UPMSecondCount) - 1
	start := int64(record.UPMOffset)
	if start+int64(length) > bankSize {
		return nil, errors.New("unit UPM span outside bank")
	}
	value := make([]byte, length)
	if _, err := bank.ReadAt(value, start); err != nil {
		return nil, fmt.Errorf("read unit UPM span: %w", err)
	}
	return value, nil
}

// UPMSides returns the first and second UPM side vectors. The sides overlap
// at the final period of the first side, as FUN_1002c120's second-side offset
// is base + firstCount - 1. The supplied combined vector must match the two
// counts stored in the unit record.
func (record UnitRecord) UPMSides(upm []byte) ([]byte, []byte, error) {
	firstCount := int(record.UPMFirstCount)
	secondCount := int(record.UPMSecondCount)
	if firstCount == 0 || secondCount == 0 {
		return nil, nil, errors.New("unit UPM side count is zero")
	}
	combinedCount := firstCount + secondCount - 1
	if len(upm) != combinedCount {
		return nil, nil, fmt.Errorf("combined UPM vector has %d periods; record requires %d", len(upm), combinedCount)
	}
	sharedBoundary := firstCount - 1
	return upm[:firstCount], upm[sharedBoundary:], nil
}
