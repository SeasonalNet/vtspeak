package synthesis

import (
	"errors"
	"math"
)

// WindowAlignment selects which side receives zero padding when the requested
// window is longer than its source.
type WindowAlignment uint8

const (
	WindowAlignLeft WindowAlignment = iota
	WindowAlignRight
)

// ApplyPaul2013BlendWindow applies the rising coefficient window used by the
// two Stage 8 sample-window helpers. The DLL indexes a 4096-entry float table;
// its observed values follow sin(pi*x/2)^2. This implementation computes that
// curve analytically, so final floating-point rounding can differ by a sample
// from the table-driven DLL.
//
// Samples are cropped from the end when source is longer than outputLength.
// When the source is shorter, alignment chooses whether zeros follow or
// precede the window. The forward helper's short-source ramp is scaled across
// outputLength; the right-aligned helper scales it across the source length.
func ApplyPaul2013BlendWindow(
	source []int16,
	outputLength int,
	alignment WindowAlignment,
) ([]int16, error) {
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
	outputStart := 0
	phaseLength := outputLength
	if alignment == WindowAlignRight {
		if len(source) < outputLength {
			outputStart = outputLength - len(source)
			phaseLength = len(source)
		}
	}
	for index := 0; index < copyCount; index++ {
		phase := index * 4096 / phaseLength
		angle := math.Pi * float64(phase) / (2 * 4096)
		coefficient := math.Sin(angle)
		coefficient *= coefficient
		value := coefficient * float64(source[index])
		if value >= 32767 {
			value = 32766
		} else if value <= -32768 {
			value = -32767
		}
		output[outputStart+index] = int16(math.RoundToEven(value))
	}
	return output, nil
}
