package synthesis

import (
	"fmt"
	"reflect"
	"testing"

	"vtspeak/engine/text"
)

func TestInsertPaul2013InlinePausesAtTokenBoundaries(t *testing.T) {
	input := PCM{SampleRate: 16000, Samples: []int16{11, 22}}
	pauses := []text.Paul2013InlinePause{
		{AfterToken: 0, DurationMilliseconds: 1, OutputFramesAt16KHz: 16},
		{AfterToken: 1, DurationMilliseconds: 0, OutputFramesAt16KHz: 0},
		{AfterToken: 1, DurationMilliseconds: 2, OutputFramesAt16KHz: 32},
		{AfterToken: 2, DurationMilliseconds: 1, OutputFramesAt16KHz: 16},
	}
	got, err := InsertPaul2013InlinePauses(input, []uint64{0, 1, 2}, pauses)
	if err != nil {
		t.Fatal(err)
	}
	if got.SampleRate != 16000 || len(got.Samples) != 2+16+32+16 {
		t.Fatalf("paused PCM profile = %d Hz, %d frames", got.SampleRate, len(got.Samples))
	}
	if got.Samples[16] != 11 || got.Samples[17] != 0 || got.Samples[49] != 22 {
		t.Fatalf("PCM samples at token boundaries = %v, want source frames preserved around inserted silence", got.Samples[16:50])
	}
	for index, sample := range got.Samples {
		if index != 16 && index != 49 && sample != 0 {
			t.Fatalf("sample %d = %d, want zero-valued pause frame", index, sample)
		}
	}
	if !reflect.DeepEqual(input.Samples, []int16{11, 22}) {
		t.Fatalf("input PCM was mutated: %v", input.Samples)
	}
}

func TestInsertPaul2013InlinePausesMatchesCapturedDurations(t *testing.T) {
	for _, duration := range []uint32{200, 1000} {
		t.Run(fmt.Sprintf("%dms", duration), func(t *testing.T) {
			frames := uint64(duration) * 16
			got, err := InsertPaul2013InlinePauses(
				PCM{SampleRate: 16000, Samples: []int16{7}},
				[]uint64{0, 1},
				[]text.Paul2013InlinePause{{
					AfterToken: 0, DurationMilliseconds: duration, OutputFramesAt16KHz: frames,
				}},
			)
			if err != nil {
				t.Fatal(err)
			}
			if uint64(len(got.Samples)) != frames+1 || got.Samples[len(got.Samples)-1] != 7 {
				t.Fatalf("%d ms pause produced %d frames, want %d including the source frame", duration, len(got.Samples), frames+1)
			}
		})
	}
}

func TestInsertPaul2013InlinePausesValidatesInputs(t *testing.T) {
	pcm := PCM{SampleRate: 16000, Samples: []int16{1, 2}}
	validPause := text.Paul2013InlinePause{AfterToken: 1, DurationMilliseconds: 1, OutputFramesAt16KHz: 16}
	for name, input := range map[string]struct {
		pcm        PCM
		boundaries []uint64
		pauses     []text.Paul2013InlinePause
		wantError  bool
	}{
		"sample rate":                {PCM{SampleRate: 8000, Samples: pcm.Samples}, []uint64{0, 2}, nil, true},
		"missing boundary table":     {pcm, nil, nil, true},
		"out-of-range frame":         {pcm, []uint64{0, 3}, nil, true},
		"unordered frame offsets":    {pcm, []uint64{2, 1}, nil, true},
		"out-of-range token":         {pcm, []uint64{0, 2}, []text.Paul2013InlinePause{{AfterToken: 2}}, true},
		"unordered token boundaries": {pcm, []uint64{0, 1, 2}, []text.Paul2013InlinePause{{AfterToken: 2}, {AfterToken: 1}}, true},
		"inconsistent duration":      {pcm, []uint64{0, 2}, []text.Paul2013InlinePause{{AfterToken: 1, DurationMilliseconds: 1, OutputFramesAt16KHz: 15}}, true},
		"valid control":              {pcm, []uint64{0, 2}, []text.Paul2013InlinePause{validPause}, false},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := InsertPaul2013InlinePauses(input.pcm, input.boundaries, input.pauses)
			if (err != nil) != input.wantError {
				t.Fatalf("error = %v, want error %v", err, input.wantError)
			}
		})
	}
}
