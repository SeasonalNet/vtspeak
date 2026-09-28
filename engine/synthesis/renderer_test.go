package synthesis

import (
	"context"
	"errors"
	"testing"

	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

type fakeUnitReader struct {
	units map[string]voice.Unit
	err   error
}

func (reader fakeUnitReader) ReadUnit(bank string, index uint32) (voice.Unit, error) {
	if reader.err != nil {
		return voice.Unit{}, reader.err
	}
	unit, ok := reader.units[bank]
	if !ok || unit.Index != index {
		return voice.Unit{}, errors.New("unit not found")
	}
	return unit, nil
}

func TestConcatenatingRendererCombinesDecodedUnitSamples(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 4, PCM: []byte{0x01, 0x00, 0xfe, 0xff}},
		"num": {Index: 9, PCM: []byte{0x34, 0x12}},
	}}
	renderer := ConcatenatingRenderer{Units: reader}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "gen", Index: 4}, {Bank: "num", Index: 9},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1})
	if err != nil {
		t.Fatal(err)
	}
	if pcm.SampleRate != 16000 {
		t.Fatalf("sample rate = %d, want 16000", pcm.SampleRate)
	}
	want := []int16{1, -2, 0x1234}
	if len(pcm.Samples) != len(want) {
		t.Fatalf("sample count = %d, want %d", len(pcm.Samples), len(want))
	}
	for i := range want {
		if pcm.Samples[i] != want[i] {
			t.Errorf("sample %d = %d, want %d", i, pcm.Samples[i], want[i])
		}
	}
}

func TestConcatenatingRendererRejectsUnsupportedInputs(t *testing.T) {
	renderer := ConcatenatingRenderer{Units: fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 0, PCM: []byte{0x01, 0x00}},
	}}}
	tests := map[string]struct {
		ctx      context.Context
		units    []selection.UnitRef
		controls Controls
	}{
		"empty sequence":       {context.Background(), nil, Controls{Pitch: -1, Speed: -1, Volume: -1}},
		"non-default controls": {context.Background(), []selection.UnitRef{{Bank: "gen"}}, Controls{Pitch: 120, Speed: -1, Volume: -1}},
		"canceled context":     {canceledContext(t), []selection.UnitRef{{Bank: "gen"}}, Controls{Pitch: -1, Speed: -1, Volume: -1}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := renderer.Render(test.ctx, test.units, test.controls); err == nil {
				t.Fatal("unsupported input accepted")
			}
		})
	}
}

func TestConcatenatingRendererRejectsMalformedUnitPCM(t *testing.T) {
	renderer := ConcatenatingRenderer{Units: fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 0, PCM: []byte{0x01}},
	}}}
	if _, err := renderer.Render(context.Background(), []selection.UnitRef{{Bank: "gen"}}, Controls{Pitch: -1, Speed: -1, Volume: -1}); err == nil {
		t.Fatal("odd byte count accepted as PCM")
	}
}

func canceledContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
