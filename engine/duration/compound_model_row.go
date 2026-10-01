package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013CompoundModelContextResult reports one explicit FUN_10009dc0 model
// row branch and its counted-row write.
type Paul2013CompoundModelContextResult struct {
	Model    []byte
	Compound Paul2013CompoundContextResult
	Applied  bool
}

// NormalizePaul2013CompoundModelContextRow runs the supported FUN_10009dc0
// component chain for one row selected by the caller's outer dispatch. It
// writes a matched result to that row's phone string and state byte.
func (engine *Engine) NormalizePaul2013CompoundModelContextRow(
	ctx context.Context,
	model []byte,
	contextRowIndex int,
	outerMode bool,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013CompoundModelContextResult, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013CompoundModelContextResult{}, errors.New("Paul 2013 compound model-row resources are nil or unloaded")
	}
	if ctx == nil {
		return Paul2013CompoundModelContextResult{}, errors.New("compound model-row normalization has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013CompoundModelContextResult{}, err
	}
	surface, currentFlags, err := text.Paul2013ModelContextSurface(model, contextRowIndex)
	if err != nil {
		return Paul2013CompoundModelContextResult{}, fmt.Errorf("read compound model context row %d: %w", contextRowIndex, err)
	}
	compound, err := engine.NormalizePaul2013CompoundContextWithSupportedHandlers(
		ctx, surface, currentFlags, outerMode, model, contextRowIndex, contractionPrefix,
	)
	if err != nil {
		return Paul2013CompoundModelContextResult{}, fmt.Errorf("normalize compound model context row %d: %w", contextRowIndex, err)
	}
	result := Paul2013CompoundModelContextResult{
		Model: append([]byte(nil), model...), Compound: compound,
	}
	if !compound.Matched {
		return result, nil
	}
	result.Model, err = text.ApplyPaul2013ModelContextCodes(
		model, contextRowIndex, compound.Output, compound.RowFlags,
	)
	if err != nil {
		return Paul2013CompoundModelContextResult{}, fmt.Errorf("write compound model context row %d: %w", contextRowIndex, err)
	}
	result.Applied = true
	return result, nil
}
