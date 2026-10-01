package synthesis

import (
	"errors"
	"fmt"

	"vtspeak/engine/dat"
)

// UPMSamplePeriod identifies one UPM-bounded source interval in decoded PCM.
// Offset and Samples are measured in mono 16 kHz samples.
type UPMSamplePeriod struct {
	Index   int
	Offset  int
	Samples []int16
}

// UPMSampleSides contains the first and second model sides in decoded PCM.
// The shared boundary period is present in both slices, with its original
// combined-vector Index and PCM Offset retained in each view.
type UPMSampleSides struct {
	First  []UPMSamplePeriod
	Second []UPMSamplePeriod
}

// SplitPaul2013UnitUPMSides joins a unit record's side counts to its decoded
// PCM and combined UPM vector. It preserves the shared period in both sides
// and checks each byte-4/byte-6 sample span against the mapped PCM windows.
func SplitPaul2013UnitUPMSides(
	record dat.UnitRecord,
	upm []byte,
	samples []int16,
) (UPMSampleSides, error) {
	firstUPM, secondUPM, err := record.UPMSides(upm)
	if err != nil {
		return UPMSampleSides{}, err
	}
	periods, err := SplitPaul2013UPMSamples(upm, samples)
	if err != nil {
		return UPMSampleSides{}, err
	}
	firstEnd := len(firstUPM)
	secondStart := len(periods) - len(secondUPM)
	if secondStart != firstEnd-1 {
		return UPMSampleSides{}, errors.New("unit UPM sides do not share exactly one boundary period")
	}
	first := append([]UPMSamplePeriod(nil), periods[:firstEnd]...)
	second := append([]UPMSamplePeriod(nil), periods[secondStart:]...)
	second[0].Samples = append([]int16(nil), second[0].Samples...)
	if sampleWindowLength(first) != int(record.FirstSideSamples) {
		return UPMSampleSides{}, fmt.Errorf(
			"first UPM side covers %d PCM samples; unit record declares %d",
			sampleWindowLength(first), record.FirstSideSamples,
		)
	}
	if sampleWindowLength(second) != int(record.SecondSideSamples) {
		return UPMSampleSides{}, fmt.Errorf(
			"second UPM side covers %d PCM samples; unit record declares %d",
			sampleWindowLength(second), record.SecondSideSamples,
		)
	}
	return UPMSampleSides{First: first, Second: second}, nil
}

func sampleWindowLength(periods []UPMSamplePeriod) int {
	length := 0
	for _, period := range periods {
		length += len(period.Samples)
	}
	return length
}

// SplitPaul2013UPMSamples maps the observed UPM vector onto decoded PCM.
// Each UPM byte describes a period of twice that many samples. The returned
// windows are copies so callers can prepare or modify them independently.
func SplitPaul2013UPMSamples(upm []byte, samples []int16) ([]UPMSamplePeriod, error) {
	if len(upm) == 0 {
		return nil, errors.New("UPM vector is empty")
	}
	if len(samples) == 0 {
		return nil, errors.New("UPM source PCM is empty")
	}
	periods := make([]UPMSamplePeriod, len(upm))
	offset := 0
	for index, period := range upm {
		length := int(period) * 2
		if length == 0 {
			return nil, errors.New("UPM period has zero length")
		}
		if length > len(samples)-offset {
			return nil, errors.New("UPM periods extend past decoded PCM")
		}
		window := append([]int16(nil), samples[offset:offset+length]...)
		periods[index] = UPMSamplePeriod{Index: index, Offset: offset, Samples: window}
		offset += length
	}
	if offset != len(samples) {
		return nil, errors.New("UPM periods do not cover decoded PCM")
	}
	return periods, nil
}
