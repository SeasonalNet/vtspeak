package synthesis

import (
	"reflect"
	"testing"
)

func TestSplitPaul2013UPMSamples(t *testing.T) {
	samples := []int16{10, 11, 12, 13, 14, 15, 16, 17}
	periods, err := SplitPaul2013UPMSamples([]byte{2, 2}, samples)
	if err != nil {
		t.Fatal(err)
	}
	want := []UPMSamplePeriod{
		{Index: 0, Offset: 0, Samples: []int16{10, 11, 12, 13}},
		{Index: 1, Offset: 4, Samples: []int16{14, 15, 16, 17}},
	}
	if !reflect.DeepEqual(periods, want) {
		t.Fatalf("periods = %#v, want %#v", periods, want)
	}
	periods[0].Samples[0] = -1
	if samples[0] != 10 {
		t.Fatal("period windows alias the caller's PCM")
	}
}

func TestSplitPaul2013UPMSamplesRejectsInvalidSpans(t *testing.T) {
	tests := []struct {
		name    string
		upm     []byte
		samples []int16
	}{
		{name: "empty UPM", samples: []int16{1}},
		{name: "empty PCM", upm: []byte{1}},
		{name: "zero period", upm: []byte{0}, samples: []int16{1}},
		{name: "span too long", upm: []byte{2}, samples: []int16{1, 2}},
		{name: "uncovered tail", upm: []byte{1}, samples: []int16{1, 2, 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if periods, err := SplitPaul2013UPMSamples(test.upm, test.samples); err == nil {
				t.Fatalf("SplitPaul2013UPMSamples() = %#v, want error", periods)
			}
		})
	}
}
