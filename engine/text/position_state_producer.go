package text

import "fmt"

// Paul2013PositionStateProducerSeries is one boundary/value pair consumed by
// FUN_1001d5d0 before FUN_10022970. The native matcher that produces these
// arrays remains an explicit input.
type Paul2013PositionStateProducerSeries struct {
	Boundaries []int32
	Values     []int32
}

// NormalizePaul2013PositionStateProducerSeries ports the mode 0-2 paired
// boundary compaction and negative-value fill in FUN_1001d5d0. It returns a
// copy and leaves the supplied arrays unchanged.
func NormalizePaul2013PositionStateProducerSeries(
	series Paul2013PositionStateProducerSeries,
) (Paul2013PositionStateProducerSeries, error) {
	if len(series.Boundaries) != len(series.Values) {
		return Paul2013PositionStateProducerSeries{}, fmt.Errorf(
			"received %d position boundaries for %d values",
			len(series.Boundaries),
			len(series.Values),
		)
	}

	boundaries := append([]int32(nil), series.Boundaries...)
	values := append([]int32(nil), series.Values...)
	paired := make([]bool, len(values))
	for index := len(values) - 1; index >= 0; index-- {
		if values[index] < 0 {
			continue
		}
		end := index
		for end < len(values) && (values[end] >= 0 || paired[end]) {
			end++
		}
		if end < len(values) {
			paired[index] = true
			paired[end] = true
		}
	}

	compactedBoundaries := make([]int32, 0, len(boundaries))
	compactedValues := make([]int32, 0, len(values))
	for index, keep := range paired {
		if keep {
			compactedBoundaries = append(compactedBoundaries, boundaries[index])
			compactedValues = append(compactedValues, values[index])
		}
	}

	paired = make([]bool, len(compactedValues))
	pairIndexes := make([]int, len(compactedValues))
	for index := len(compactedValues) - 1; index >= 0; index-- {
		if compactedValues[index] < 0 {
			continue
		}
		end := index
		for end < len(compactedValues) && (compactedValues[end] >= 0 || paired[end]) {
			end++
		}
		if end < len(compactedValues) {
			paired[index] = true
			paired[end] = true
			pairIndexes[index] = end
			pairIndexes[end] = index
		}
	}
	for index, value := range compactedValues {
		if value >= 0 {
			continue
		}
		previousBoundary := pairIndexes[index]
		if previousBoundary == 0 {
			compactedValues[index] = -1
			continue
		}
		compactedValues[index] = compactedValues[previousBoundary-1]
	}

	return Paul2013PositionStateProducerSeries{
		Boundaries: compactedBoundaries,
		Values:     compactedValues,
	}, nil
}

// NormalizePaul2013PositionStateProducerArrays applies the mode 0-2 producer
// step in native order. Missing modes remain absent, as in the caller's zero
// count arrays.
func NormalizePaul2013PositionStateProducerArrays(
	series [3]Paul2013PositionStateProducerSeries,
) ([3]Paul2013PositionStateProducerSeries, error) {
	var result [3]Paul2013PositionStateProducerSeries
	for mode, input := range series {
		var err error
		result[mode], err = NormalizePaul2013PositionStateProducerSeries(input)
		if err != nil {
			return [3]Paul2013PositionStateProducerSeries{}, fmt.Errorf(
				"normalize position-state producer mode %d: %w",
				mode,
				err,
			)
		}
	}
	return result, nil
}
