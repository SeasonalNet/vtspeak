// Package voice reads the observed 2013 Paul unit index and its DAT/UPM banks.
package voice

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vtspeak/engine/dat"
)

type Unit struct {
	Index  uint32
	Record dat.UnitRecord
	UPM    []byte
	PCM    []byte
}

// Bank provides concurrent-safe reads from one matching index, DAT bank, and
// UPM bank. The model inputs are opened read-only.
type Bank struct {
	index     []byte
	datFile   *os.File
	upmFile   *os.File
	datSize   int64
	upmSize   int64
	unitCount uint32
}

func OpenBank(indexPath, datPath, upmPath string) (*Bank, error) {
	index, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read unit index: %w", err)
	}
	if len(index) < 45 {
		return nil, fmt.Errorf("unit index is shorter than its header")
	}
	name, err := dat.BankName(index)
	if err != nil {
		return nil, fmt.Errorf("validate unit index: %w", err)
	}
	if name != strings.TrimSuffix(filepath.Base(datPath), filepath.Ext(datPath)) ||
		name != strings.TrimSuffix(filepath.Base(upmPath), filepath.Ext(upmPath)) {
		return nil, fmt.Errorf("index, DAT, and UPM bank names do not match")
	}
	count := binary.LittleEndian.Uint32(index[39:43])
	if count == 0 {
		return nil, fmt.Errorf("unit index contains no units")
	}
	if _, err := dat.ReadUnit(index, 0); err != nil {
		return nil, fmt.Errorf("validate unit index: %w", err)
	}
	datFile, err := os.Open(datPath)
	if err != nil {
		return nil, fmt.Errorf("open DAT bank: %w", err)
	}
	upmFile, err := os.Open(upmPath)
	if err != nil {
		_ = datFile.Close()
		return nil, fmt.Errorf("open UPM bank: %w", err)
	}
	datInfo, err := datFile.Stat()
	if err != nil {
		_ = datFile.Close()
		_ = upmFile.Close()
		return nil, fmt.Errorf("stat DAT bank: %w", err)
	}
	upmInfo, err := upmFile.Stat()
	if err != nil {
		_ = datFile.Close()
		_ = upmFile.Close()
		return nil, fmt.Errorf("stat UPM bank: %w", err)
	}
	return &Bank{
		index:     index,
		datFile:   datFile,
		upmFile:   upmFile,
		datSize:   datInfo.Size(),
		upmSize:   upmInfo.Size(),
		unitCount: count,
	}, nil
}

func (b *Bank) UnitCount() uint32 { return b.unitCount }

// ReadRecord returns the indexed unit metadata and feature row without
// reading either waveform bank or decoding PCM.
func (b *Bank) ReadRecord(index uint32) (dat.UnitRecord, error) {
	if b == nil {
		return dat.UnitRecord{}, errors.New("unit bank is nil")
	}
	if index >= b.unitCount {
		return dat.UnitRecord{}, fmt.Errorf("unit %d outside bank range 0..%d", index, b.unitCount-1)
	}
	record, err := dat.ReadUnit(b.index, index)
	if err != nil {
		return dat.UnitRecord{}, err
	}
	return record, nil
}

func (b *Bank) ReadUnit(index uint32) (Unit, error) {
	record, err := b.ReadRecord(index)
	if err != nil {
		return Unit{}, err
	}
	payload, err := dat.UnitPayload(b.index, b.datFile, b.datSize, index)
	if err != nil {
		return Unit{}, err
	}
	upm, err := dat.UPMPayload(record, b.upmFile, b.upmSize)
	if err != nil {
		return Unit{}, err
	}
	pcm, err := dat.Decode(payload)
	if err != nil {
		return Unit{}, fmt.Errorf("decode unit %d: %w", index, err)
	}
	if err := validateUnitTiming(record, upm, pcm); err != nil {
		return Unit{}, fmt.Errorf("validate unit %d timing metadata: %w", index, err)
	}
	return Unit{Index: index, Record: record, UPM: upm, PCM: pcm}, nil
}

func validateUnitTiming(record dat.UnitRecord, upm, pcm []byte) error {
	if len(upm) == 0 {
		return errors.New("UPM vector is empty")
	}
	sharedIndex := int(record.UPMFirstCount) - 1
	if sharedIndex < 0 || sharedIndex >= len(upm) {
		return errors.New("shared UPM edge falls outside the combined vector")
	}
	if record.UPMEdges != [3]byte{upm[0], upm[sharedIndex], upm[len(upm)-1]} {
		return errors.New("cached UPM edge bytes do not match the combined vector")
	}
	firstSide, secondSide, err := record.UPMSides(upm)
	if err != nil {
		return err
	}
	var firstSideSamples, secondSideSamples uint64
	for _, period := range firstSide {
		firstSideSamples += uint64(period) * 2
	}
	for _, period := range secondSide {
		secondSideSamples += uint64(period) * 2
	}
	if firstSideSamples != uint64(record.FirstSideSamples) ||
		secondSideSamples != uint64(record.SecondSideSamples) {
		return fmt.Errorf(
			"UPM sides describe %d and %d samples; unit record declares %d and %d",
			firstSideSamples, secondSideSamples,
			record.FirstSideSamples, record.SecondSideSamples,
		)
	}
	var expectedSamples uint64
	for _, period := range upm {
		expectedSamples += uint64(period) * 2
	}
	if len(pcm)%2 != 0 || uint64(len(pcm)/2) != expectedSamples {
		return fmt.Errorf("decoded PCM has %d samples, UPM vector describes %d", len(pcm)/2, expectedSamples)
	}
	return nil
}

func (b *Bank) Close() error {
	datErr := b.datFile.Close()
	upmErr := b.upmFile.Close()
	return errors.Join(datErr, upmErr)
}
