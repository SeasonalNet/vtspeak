package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// PausedTextEvaluation joins the supported inline-pause parser, dictionary
// pronunciation selection, and duration/pitch tree evaluation. Pause frame
// boundaries remain attached to the sequence for the later timeline stage;
// this result does not insert audio silence by itself.
type PausedTextEvaluation struct {
	Sequence text.LexicalPhoneSequence
	Tokens   []TokenResult
	Pauses   []text.Paul2013InlinePause
}

// CapturedVTMLTextEvaluation carries duration and pitch results for the
// supported pause, mark, and substitution forms. Marks remain source-positioned
// because their output-frame mapping is not yet established.
type CapturedVTMLTextEvaluation struct {
	Sequence text.LexicalPhoneSequence
	Tokens   []TokenResult
	Pauses   []text.Paul2013InlinePause
	Marks    []text.Paul2013InlineMarkEvent
}

// EvaluateWithCapturedVTML composes the bounded VTML frontend with automatic
// pronunciation choice and the loaded Paul duration/pitch trees.
func (engine *Engine) EvaluateWithCapturedVTML(
	ctx context.Context,
	source string,
) (CapturedVTMLTextEvaluation, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil || engine.catalog == nil {
		return CapturedVTMLTextEvaluation{}, errors.New("Paul 2013 text evaluation engine is nil or incompletely loaded")
	}
	if ctx == nil {
		return CapturedVTMLTextEvaluation{}, errors.New("Paul 2013 VTML evaluation has no context")
	}
	frontend := text.LexiconFrontend{Dictionary: engine.dictionary}
	parsed, err := frontend.ResolveTextWithPaul2013CapturedVTML(ctx, source)
	if err != nil {
		return CapturedVTMLTextEvaluation{}, fmt.Errorf("resolve Paul 2013 captured VTML text: %w", err)
	}
	analysis, err := frontend.AnalyzeResolvedTokens(ctx, parsed.Tokens)
	if err != nil {
		return CapturedVTMLTextEvaluation{}, fmt.Errorf("analyze Paul 2013 VTML pronunciations: %w", err)
	}
	choices, err := choosePaul2013PronunciationAlternatives(analysis, engine.EvaluatePronunciationFeatures)
	if err != nil {
		return CapturedVTMLTextEvaluation{}, err
	}
	sequence, err := text.SelectPaul2013CapturedVTMLPronunciations(parsed, choices)
	if err != nil {
		return CapturedVTMLTextEvaluation{}, fmt.Errorf("assemble Paul 2013 VTML phone sequence: %w", err)
	}
	terminalMarkers := make([]byte, len(sequence.Tokens))
	if len(terminalMarkers) != 0 {
		terminalMarkers[len(terminalMarkers)-1] = 'Z'
	}
	results, err := engine.EvaluateSequenceWithMarkers(sequence, terminalMarkers)
	if err != nil {
		return CapturedVTMLTextEvaluation{}, fmt.Errorf("evaluate Paul 2013 VTML duration and pitch: %w", err)
	}
	return CapturedVTMLTextEvaluation{
		Sequence: sequence, Tokens: results,
		Pauses: append([]text.Paul2013InlinePause(nil), sequence.InlinePauses...),
		Marks:  append([]text.Paul2013InlineMarkEvent(nil), sequence.InlineMarks...),
	}, nil
}

// EvaluateWithInlinePauses resolves the observed self-closing VTML pause form
// alongside ordinary text, selects pronunciations with the loaded shared
// classifier, and evaluates both Paul duration and pitch trees. Pause events
// are retained in source order for a renderer that can map token boundaries
// to output frames.
func (engine *Engine) EvaluateWithInlinePauses(
	ctx context.Context,
	source string,
) (PausedTextEvaluation, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil || engine.catalog == nil {
		return PausedTextEvaluation{}, errors.New("Paul 2013 text evaluation engine is nil or incompletely loaded")
	}
	if ctx == nil {
		return PausedTextEvaluation{}, errors.New("Paul 2013 paused text evaluation has no context")
	}
	frontend := text.LexiconFrontend{Dictionary: engine.dictionary}
	paused, err := frontend.ResolveTextWithPaul2013InlinePauses(ctx, source)
	if err != nil {
		return PausedTextEvaluation{}, fmt.Errorf("resolve Paul 2013 inline-pause text: %w", err)
	}
	analysis, err := frontend.AnalyzeResolvedTokens(ctx, paused.Tokens)
	if err != nil {
		return PausedTextEvaluation{}, fmt.Errorf("analyze Paul 2013 inline-pause pronunciations: %w", err)
	}
	choices, err := choosePaul2013PronunciationAlternatives(analysis, engine.EvaluatePronunciationFeatures)
	if err != nil {
		return PausedTextEvaluation{}, err
	}
	sequence, err := text.SelectLexicalPronunciationsWithPauses(paused, choices)
	if err != nil {
		return PausedTextEvaluation{}, fmt.Errorf("assemble Paul 2013 inline-pause phone sequence: %w", err)
	}
	terminalMarkers := make([]byte, len(sequence.Tokens))
	if len(terminalMarkers) != 0 {
		terminalMarkers[len(terminalMarkers)-1] = 'Z'
	}
	results, err := engine.EvaluateSequenceWithMarkers(sequence, terminalMarkers)
	if err != nil {
		return PausedTextEvaluation{}, fmt.Errorf("evaluate Paul 2013 inline-pause duration and pitch: %w", err)
	}
	return PausedTextEvaluation{
		Sequence: sequence,
		Tokens:   results,
		Pauses:   append([]text.Paul2013InlinePause(nil), sequence.InlinePauses...),
	}, nil
}
