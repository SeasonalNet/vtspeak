package synthesis

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"math"
)

// paul2013WindowTableBytes contains the observed 8,192-entry little-endian
// float32 blend table at DLL RVA 0x6d1b0. The reproducible extractor is
// tools/revkit/scripts/extract_paul2013_window_table.py.
//
//go:embed window_table.bin
var paul2013WindowTableBytes []byte

// WindowAlignment selects which side receives zero padding when the requested
// window is longer than its source.
type WindowAlignment uint8

const (
	WindowAlignLeft WindowAlignment = iota
	WindowAlignRight
)

// ApplyPaul2013BlendWindow applies the coefficient window used by the two
// Stage 8 sample-window helpers. The DLL has an 8192-entry float table at
// 0x1006d1b0; the left helper indexes from zero and the right helper starts at
// index 4096. Coefficients come from the fixed table embedded with this
// package, preserving the DLL's float32 values exactly.
//
// This ports the aligned source copying in FUN_1002d010 (left) and
// FUN_1002cdf0 (right): a long source is cropped at the aligned edge; a short
// source is zero-padded on the opposite side. The native helper indexes the
// table using the smaller source/output length.
func ApplyPaul2013BlendWindow(
	source []int16,
	outputLength int,
	alignment WindowAlignment,
) ([]int16, error) {
	return applyPaul2013Window(source, outputLength, alignment, alignment == WindowAlignRight)
}

// BlendPaul2013UPMSegmentWindows combines the current output-context window
// and selected source-period window reconstructed by FUN_1002afb0. The output
// context uses the falling table half and is left aligned; the source period
// uses the rising half and is right aligned. Short windows are zero padded,
// long windows are cropped at the aligned edge, and the weighted int16 lanes
// are added with the native 16-bit wrap behavior.
func BlendPaul2013UPMSegmentWindows(first, second []int16, outputLength int) ([]int16, error) {
	if len(first) == 0 || len(second) == 0 {
		return nil, errors.New("UPM segment source windows must be nonempty")
	}
	if outputLength < 1 || outputLength > math.MaxInt16 ||
		len(first) > math.MaxInt16 || len(second) > math.MaxInt16 {
		return nil, errors.New("UPM segment window lengths must be in the range 1 through 32767")
	}
	falling, err := applyPaul2013Window(first, outputLength, WindowAlignLeft, true)
	if err != nil {
		return nil, err
	}
	rising, err := applyPaul2013Window(second, outputLength, WindowAlignRight, false)
	if err != nil {
		return nil, err
	}
	for index := range falling {
		falling[index] = int16(int32(falling[index]) + int32(rising[index]))
	}
	return falling, nil
}

func applyPaul2013Window(source []int16, outputLength int, alignment WindowAlignment, falling bool) ([]int16, error) {
	if len(source) == 0 {
		return nil, errors.New("blend window source is empty")
	}
	if outputLength < 1 || outputLength > math.MaxInt16 || len(source) > math.MaxInt16 {
		return nil, errors.New("blend window lengths must be in the range 1 through 32767")
	}
	if alignment != WindowAlignLeft && alignment != WindowAlignRight {
		return nil, errors.New("blend window alignment is unsupported")
	}
	output := make([]int16, outputLength)
	copyCount := len(source)
	if copyCount > outputLength {
		copyCount = outputLength
	}
	sourceStart := 0
	outputStart := 0
	if alignment == WindowAlignRight && len(source) < outputLength {
		outputStart = outputLength - len(source)
	} else if alignment == WindowAlignRight {
		sourceStart = len(source) - outputLength
	}
	for index := 0; index < copyCount; index++ {
		phase := index * paul2013WindowTableHalf / copyCount
		if falling {
			phase += paul2013WindowTableHalf
		}
		coefficient := paul2013WindowCoefficient(phase)
		output[outputStart+index] = paul2013WeightedLaneSample(float64(source[sourceStart+index]) * coefficient)
	}
	return output, nil
}

const paul2013WindowTableHalf = 4096

func paul2013WindowCoefficient(index int) float64 {
	return float64(math.Float32frombits(binary.LittleEndian.Uint32(paul2013WindowTableBytes[index*4:])))
}

func paul2013WeightedLaneSample(value float64) int16 {
	if value > 32766 {
		value = 32766
	} else if value < -32767 {
		value = -32767
	}
	// FUN_1006479c sets the x87 rounding-control bits to round toward zero
	// around FISTP, then restores the previous control word.
	return int16(math.Trunc(value))
}
