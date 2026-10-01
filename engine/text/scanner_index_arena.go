package text

import (
	"encoding/binary"
	"fmt"
)

// Paul2013ScannerMemoryReader resolves an exact span of native address space.
// The result is borrowed only during the call; index snapshots copy all cells.
type Paul2013ScannerMemoryReader func(address uint32, size int) ([]byte, error)

// ReadPaul2013ScannerLookupIndexes projects FUN_1005e010/FUN_1005e1a0's
// index objects through model +0xc/+0x10. It resolves only active key ranges;
// the native lookup never touches a bucket with either endpoint negative.
func ReadPaul2013ScannerLookupIndexes(modelAddress uint32, read Paul2013ScannerMemoryReader, resolveCString Paul2013CStringPointerResolver) (exact, mapped *Paul2013ScannerLookupIndex, err error) {
	span := func(address uint64, size int) ([]byte, error) {
		if address > 0xffffffff || size < 0 || address+uint64(size) > 0x100000000 {
			return nil, fmt.Errorf("scanner index native span exceeds address space")
		}
		if read == nil {
			return nil, fmt.Errorf("scanner index memory reader is unavailable")
		}
		raw, err := read(uint32(address), size)
		if err != nil {
			return nil, err
		}
		if len(raw) < size {
			return nil, fmt.Errorf("scanner index span at %#x is truncated", address)
		}
		return append([]byte(nil), raw[:size]...), nil
	}
	model, err := span(uint64(modelAddress)+0xc, 8)
	if err != nil {
		return nil, nil, err
	}
	project := func(address uint32) (*Paul2013ScannerLookupIndex, error) {
		if address == 0 {
			return nil, nil
		}
		header, err := span(uint64(address), 16)
		if err != nil {
			return nil, err
		}
		count := int32(binary.LittleEndian.Uint32(header))
		index := &Paul2013ScannerLookupIndex{}
		if count <= 0 {
			return index, nil
		}
		lowerAddress, upperAddress, rowAddress := binary.LittleEndian.Uint32(header[4:]), binary.LittleEndian.Uint32(header[8:]), binary.LittleEndian.Uint32(header[12:])
		if lowerAddress == 0 || upperAddress == 0 {
			return nil, fmt.Errorf("active scanner index has a null bound pointer")
		}
		lower, err := span(uint64(lowerAddress), int(count)*4)
		if err != nil {
			return nil, err
		}
		upper, err := span(uint64(upperAddress), int(count)*4)
		if err != nil {
			return nil, err
		}
		index.Lower = make([]int32, count)
		index.Upper = make([]int32, count)
		maximum := int32(-1)
		for bucket := range index.Lower {
			low, high := int32(binary.LittleEndian.Uint32(lower[bucket*4:])), int32(binary.LittleEndian.Uint32(upper[bucket*4:]))
			index.Lower[bucket], index.Upper[bucket] = low, high
			if low < 0 || high < 0 {
				continue
			}
			if high < low {
				return nil, fmt.Errorf("scanner index bucket %d is reversed", bucket)
			}
			if uint64(rowAddress)+uint64(high)*20+16 > 0x100000000 {
				return nil, fmt.Errorf("scanner index row address exceeds native address space")
			}
			if high > maximum {
				maximum = high
			}
		}
		if maximum < 0 {
			return index, nil
		}
		if rowAddress == 0 || resolveCString == nil {
			return nil, fmt.Errorf("active scanner index row/string resolver is unavailable")
		}
		index.nativeKeys = make(map[int32][]byte)
		for bucket, low := range index.Lower {
			high := index.Upper[bucket]
			if low < 0 || high < 0 {
				continue
			}
			for row := int64(low); row <= int64(high); row++ {
				if _, exists := index.nativeKeys[int32(row)]; exists {
					continue
				}
				cell, err := span(uint64(rowAddress)+uint64(row)*20+12, 4)
				if err != nil {
					return nil, err
				}
				pointer := binary.LittleEndian.Uint32(cell)
				if pointer == 0 {
					return nil, fmt.Errorf("scanner index key row %d has a null pointer", row)
				}
				text, err := resolveCString(pointer)
				if err != nil {
					return nil, err
				}
				index.nativeKeys[int32(row)] = append([]byte(nil), cString(text)...)
			}
		}
		return index, nil
	}
	exact, err = project(binary.LittleEndian.Uint32(model))
	if err != nil {
		return nil, nil, err
	}
	mapped, err = project(binary.LittleEndian.Uint32(model[4:]))
	if err != nil {
		return nil, nil, err
	}
	return exact, mapped, nil
}

// NewPaul2013ScannerArenaTokenLookup composes the native pointer reader,
// index searches and boundary body. Arena reads remain caller-owned so that
// changes to a native model object are observed on each lookup invocation.
func NewPaul2013ScannerArenaTokenLookup(read Paul2013ScannerMemoryReader, resolveCString Paul2013CStringPointerResolver, scan Paul2013ModelScanner) Paul2013ScannerTokenLookup {
	return NewPaul2013ScannerIndexedTokenLookup(func(address uint32) (*Paul2013ScannerLookupIndex, *Paul2013ScannerLookupIndex, error) {
		return ReadPaul2013ScannerLookupIndexes(address, read, resolveCString)
	}, scan)
}
