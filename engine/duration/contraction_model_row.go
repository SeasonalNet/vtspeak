package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013ContractionModelContextResult reports the c040 suffix result and
// counted-row write for one caller-selected model context row.
type Paul2013ContractionModelContextResult struct {
	Model       []byte
	Contraction Paul2013ContractionSuffixResult
	Applied     bool
}

// NormalizePaul2013ContractionModelContextRow runs the supported FUN_1000c040
// path for one selected row and writes a match into a copied model arena.
// Outer branch selection remains the caller's responsibility.
func (engine *Engine) NormalizePaul2013ContractionModelContextRow(
	ctx context.Context,
	model []byte,
	contextRowIndex int,
	outerMode bool,
) (Paul2013ContractionModelContextResult, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013ContractionModelContextResult{}, errors.New("Paul 2013 contraction model-row resources are nil or unloaded")
	}
	if ctx == nil {
		return Paul2013ContractionModelContextResult{}, errors.New("contraction model-row normalization has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContractionModelContextResult{}, err
	}
	surface, currentFlags, err := text.Paul2013ModelContextSurface(model, contextRowIndex)
	if err != nil {
		return Paul2013ContractionModelContextResult{}, fmt.Errorf("read contraction model context row %d: %w", contextRowIndex, err)
	}
	contraction, err := engine.NormalizePaul2013ContractionSuffixWithSupportedPrefix(
		ctx, surface, currentFlags, outerMode, model, contextRowIndex,
	)
	if err != nil {
		return Paul2013ContractionModelContextResult{}, fmt.Errorf("normalize contraction model context row %d: %w", contextRowIndex, err)
	}
	result := Paul2013ContractionModelContextResult{
		Model: append([]byte(nil), model...), Contraction: contraction,
	}
	if !contraction.Matched {
		return result, nil
	}
	result.Model, err = text.ApplyPaul2013ModelContextCodes(
		model, contextRowIndex, contraction.Output, contraction.RowFlags,
	)
	if err != nil {
		return Paul2013ContractionModelContextResult{}, fmt.Errorf("write contraction model context row %d: %w", contextRowIndex, err)
	}
	result.Applied = true
	return result, nil
}
