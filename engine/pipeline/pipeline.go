// Package pipeline defines the provisional text-to-speech stage contracts.
// Concrete linguistic and synthesis stages will be added as evidence is
// translated into the independent engine.
package pipeline

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/dat"
	"vtspeak/engine/selection"
	"vtspeak/engine/synthesis"
	"vtspeak/engine/text"
)

var ErrStageUnavailable = errors.New("pipeline stage is not implemented")

// Request is the engine-level text synthesis request.
type Request struct {
	Text     string
	Controls synthesis.Controls
}

// Pipeline joins the text, selection, and selected-unit rendering stages.
// Each dependency can be developed and compared with its own runtime evidence.
type Pipeline struct {
	Frontend text.Frontend
	Selector selection.Selector
	Renderer synthesis.Renderer
}

func (p Pipeline) Synthesize(ctx context.Context, request Request) (synthesis.PCM, error) {
	if request.Text == "" {
		return synthesis.PCM{}, errors.New("text is empty")
	}
	if p.Frontend == nil {
		return synthesis.PCM{}, fmt.Errorf("frontend: %w", ErrStageUnavailable)
	}
	phones, err := p.Frontend.Process(ctx, request.Text)
	if err != nil {
		return synthesis.PCM{}, fmt.Errorf("frontend: %w", err)
	}
	if len(phones) == 0 {
		return synthesis.PCM{}, errors.New("frontend returned no phone contexts")
	}
	if p.Selector == nil {
		return synthesis.PCM{}, fmt.Errorf("unit selection: %w", ErrStageUnavailable)
	}
	units, err := p.Selector.Select(ctx, phones)
	if err != nil {
		return synthesis.PCM{}, fmt.Errorf("unit selection: %w", err)
	}
	if len(units) == 0 {
		return synthesis.PCM{}, errors.New("selector returned no units")
	}
	if p.Renderer == nil {
		return synthesis.PCM{}, fmt.Errorf("synthesis: %w", ErrStageUnavailable)
	}
	pcm, err := p.Renderer.Render(ctx, units, request.Controls)
	if err != nil {
		return synthesis.PCM{}, fmt.Errorf("synthesis: %w", err)
	}
	if pcm.SampleRate != 16000 || len(pcm.Samples) == 0 {
		return synthesis.PCM{}, errors.New("renderer returned invalid PCM profile")
	}
	return pcm, nil
}

// SynthesizeWAV returns the standard mono 16 kHz PCM WAVE representation.
func (p Pipeline) SynthesizeWAV(ctx context.Context, request Request) ([]byte, error) {
	pcm, err := p.Synthesize(ctx, request)
	if err != nil {
		return nil, err
	}
	return dat.WAVFromSamples(pcm.Samples)
}
