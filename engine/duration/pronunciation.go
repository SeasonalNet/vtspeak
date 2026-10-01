package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// ResolvePronunciationSequence selects one dictionary alternative for each
// token by evaluating the loaded shared classifier and applying the recovered
// path-group ranker. The caller supplies only text; dictionary metadata,
// classifier rows, outputs, and group scores are produced by the engine.
func (engine *Engine) ResolvePronunciationSequence(ctx context.Context, source string) (text.LexicalPhoneSequence, error) {
	if engine == nil {
		return text.LexicalPhoneSequence{}, errors.New("Paul 2013 duration engine is nil")
	}
	if engine.pronunciationTree == nil {
		return text.LexicalPhoneSequence{}, errors.New("Paul 2013 pronunciation tree is not loaded")
	}
	analysis, err := (text.LexiconFrontend{Dictionary: engine.dictionary}).AnalyzeText(ctx, source)
	if err != nil {
		return text.LexicalPhoneSequence{}, fmt.Errorf("analyze Paul 2013 pronunciations: %w", err)
	}
	choices, err := choosePaul2013PronunciationAlternatives(analysis, engine.EvaluatePronunciationFeatures)
	if err != nil {
		return text.LexicalPhoneSequence{}, err
	}
	sequence, err := text.SelectLexicalPronunciations(analysis.Tokens, choices)
	if err != nil {
		return text.LexicalPhoneSequence{}, fmt.Errorf("assemble selected Paul 2013 pronunciations: %w", err)
	}
	return sequence, nil
}

type paul2013PronunciationEvaluator func([text.Paul2013PronunciationFeatureCount]int16) (int16, error)

func choosePaul2013PronunciationAlternatives(
	analysis text.LexicalAnalysis,
	evaluate paul2013PronunciationEvaluator,
) ([]int, error) {
	if evaluate == nil {
		return nil, errors.New("Paul 2013 pronunciation classifier evaluator is nil")
	}
	if len(analysis.Tokens) == 0 {
		return nil, errors.New("Paul 2013 pronunciation analysis has no tokens")
	}
	if len(analysis.PronunciationCandidates) != len(analysis.Tokens) {
		return nil, fmt.Errorf("Paul 2013 analysis has %d candidate lists for %d tokens", len(analysis.PronunciationCandidates), len(analysis.Tokens))
	}

	choices := make([]int, len(analysis.Tokens))
	for tokenIndex, token := range analysis.Tokens {
		if len(token.Alternatives) == 0 {
			return nil, fmt.Errorf("token %d (%q) has no pronunciation alternatives", tokenIndex, token.SourceSurface)
		}
		if len(token.Alternatives) == 1 {
			continue
		}

		groups := make([][]byte, len(token.Alternatives))
		for alternativeIndex, alternative := range token.Alternatives {
			if len(alternative.PathGroups) != 1 {
				return nil, fmt.Errorf(
					"token %d (%q) alternative %d has %d path groups; native pronunciation slots require one group per alternative",
					tokenIndex, token.SourceSurface, alternativeIndex, len(alternative.PathGroups),
				)
			}
			groups[alternativeIndex] = alternative.PathGroups[0]
		}

		outputs := make([]int16, 0, len(token.Alternatives))
		for _, candidate := range analysis.PronunciationCandidates[tokenIndex] {
			if candidate.AlternativeIndex < 0 || candidate.AlternativeIndex >= len(groups) || candidate.PathGroupIndex != 0 {
				return nil, fmt.Errorf("token %d (%q) has a candidate outside the native alternative/group layout", tokenIndex, token.SourceSurface)
			}
			if candidate.PathCodeIndex < 0 || candidate.PathCodeIndex >= len(groups[candidate.AlternativeIndex]) {
				return nil, fmt.Errorf("token %d (%q) candidate code index %d is outside alternative %d's path group", tokenIndex, token.SourceSurface, candidate.PathCodeIndex, candidate.AlternativeIndex)
			}
			if !candidate.Features.Complete() {
				return nil, fmt.Errorf("token %d (%q) candidate features are missing positions %v", tokenIndex, token.SourceSurface, candidate.Features.MissingPositions())
			}
			output, err := evaluate(candidate.Features.Values)
			if err != nil {
				return nil, fmt.Errorf("evaluate token %d (%q) alternative %d path code %d: %w", tokenIndex, token.SourceSurface, candidate.AlternativeIndex, candidate.PathCodeIndex, err)
			}
			outputs = append(outputs, output)
		}
		selected, _, err := text.RankPaul2013PronunciationPathGroups(groups, outputs)
		if err != nil {
			return nil, fmt.Errorf("rank token %d (%q) pronunciations: %w", tokenIndex, token.SourceSurface, err)
		}
		choices[tokenIndex] = selected
	}
	return choices, nil
}
