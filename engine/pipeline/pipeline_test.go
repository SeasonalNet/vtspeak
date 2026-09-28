package pipeline

import (
	"context"
	"errors"
	"testing"

	"vtspeak/engine/selection"
	"vtspeak/engine/synthesis"
	"vtspeak/engine/text"
)

type testFrontend struct{}

func (testFrontend) Process(context.Context, string) ([]text.Context, error) {
	return []text.Context{{Signature: [7]byte{1, 2, 3, 4, 5, 6, 7}}}, nil
}

type testSelector struct{}

func (testSelector) Select(context.Context, []text.Context) ([]selection.UnitRef, error) {
	return []selection.UnitRef{{Bank: "gen", Index: 42}}, nil
}

type testRenderer struct{}

func (testRenderer) Render(context.Context, []selection.UnitRef, synthesis.Controls) (synthesis.PCM, error) {
	return synthesis.PCM{SampleRate: 16000, Samples: []int16{123}}, nil
}

func TestPipelineRequiresUnimplementedStage(t *testing.T) {
	_, err := (Pipeline{}).Synthesize(context.Background(), Request{Text: "hello"})
	if !errors.Is(err, ErrStageUnavailable) {
		t.Fatalf("error = %v, want ErrStageUnavailable", err)
	}
}

func TestPipelineRunsConfiguredStages(t *testing.T) {
	pipeline := Pipeline{Frontend: testFrontend{}, Selector: testSelector{}, Renderer: testRenderer{}}
	got, err := pipeline.Synthesize(context.Background(), Request{Text: "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if got.SampleRate != 16000 || len(got.Samples) != 1 || got.Samples[0] != 123 {
		t.Fatalf("unexpected PCM: %+v", got)
	}
}

func TestPipelineRejectsEmptyText(t *testing.T) {
	_, err := (Pipeline{}).Synthesize(context.Background(), Request{})
	if err == nil {
		t.Fatal("empty text accepted")
	}
}
