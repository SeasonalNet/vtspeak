package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013C3A0ModelContextResult reports the explicit C3A0 row branch and its
// model write. The caller owns the outer FUN_10007520 dispatch gate; this
// helper only runs the recovered C3A0 source gate and component sequence.
type Paul2013C3A0ModelContextResult struct {
	Model    []byte
	Sequence Paul2013C3A0SequenceResult
	Applied  bool
}

// Paul2013C3A0ModelContextRowsResult records ordered branch outcomes for a
// caller-selected set of counted model rows.
type Paul2013C3A0ModelContextRowsResult struct {
	Model         []byte
	AppliedRows   []int
	UnmatchedRows []int
}

// NormalizePaul2013C3A0ModelContextRow composes the recovered embedded-record,
// FUN_100086C0, supported A140, apostrophe-s, and generic C3A0 handlers for a
// counted model context row. Call it only after the caller's native dispatch
// has selected the C3A0 branch. A FUN_100086C0 path that still needs neighbor
// state fails closed through the existing supported-path boundary.
func (engine *Engine) NormalizePaul2013C3A0ModelContextRow(
	ctx context.Context,
	model []byte,
	contextRowIndex int,
) (Paul2013C3A0ModelContextResult, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013C3A0ModelContextResult{}, errors.New("Paul 2013 C3A0 model-row resources are nil or unloaded")
	}
	if ctx == nil {
		return Paul2013C3A0ModelContextResult{}, errors.New("C3A0 model-row normalization has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013C3A0ModelContextResult{}, err
	}
	surface, currentFlags, err := text.Paul2013ModelContextSurface(model, contextRowIndex)
	if err != nil {
		return Paul2013C3A0ModelContextResult{}, fmt.Errorf("read C3A0 model context row %d: %w", contextRowIndex, err)
	}
	sequence, err := engine.NormalizePaul2013C3A0ComponentSequenceWithEmbeddedFUN100086C0(
		ctx, surface, currentFlags, model, contextRowIndex, nil,
	)
	if err != nil {
		return Paul2013C3A0ModelContextResult{}, fmt.Errorf("normalize C3A0 model context row %d: %w", contextRowIndex, err)
	}
	result := Paul2013C3A0ModelContextResult{
		Model: append([]byte(nil), model...), Sequence: sequence,
	}
	if !sequence.Matched {
		return result, nil
	}
	result.Model, err = text.ApplyPaul2013ModelContextCodes(
		model, contextRowIndex, sequence.Output, sequence.RowFlags,
	)
	if err != nil {
		return Paul2013C3A0ModelContextResult{}, fmt.Errorf("write C3A0 model context row %d: %w", contextRowIndex, err)
	}
	result.Applied = true
	return result, nil
}

// NormalizePaul2013C3A0ModelContextRows applies the explicit C3A0 branch to
// each caller-selected row in order, carrying every successful phone/state
// write into the next row's FUN_100086C0 evaluation. The caller supplies the
// outer-dispatch decisions and native row order.
func (engine *Engine) NormalizePaul2013C3A0ModelContextRows(
	ctx context.Context,
	model []byte,
	contextRowIndexes []int,
) (Paul2013C3A0ModelContextRowsResult, error) {
	if ctx == nil {
		return Paul2013C3A0ModelContextRowsResult{}, errors.New("C3A0 model-row pass has no context")
	}
	result := Paul2013C3A0ModelContextRowsResult{Model: append([]byte(nil), model...)}
	seen := make(map[int]struct{}, len(contextRowIndexes))
	for _, rowIndex := range contextRowIndexes {
		if err := ctx.Err(); err != nil {
			return Paul2013C3A0ModelContextRowsResult{}, err
		}
		if _, exists := seen[rowIndex]; exists {
			return Paul2013C3A0ModelContextRowsResult{}, fmt.Errorf("C3A0 model-row pass repeats row %d", rowIndex)
		}
		seen[rowIndex] = struct{}{}
		row, err := engine.NormalizePaul2013C3A0ModelContextRow(ctx, result.Model, rowIndex)
		if err != nil {
			return Paul2013C3A0ModelContextRowsResult{}, fmt.Errorf("normalize selected C3A0 model row %d: %w", rowIndex, err)
		}
		result.Model = row.Model
		if row.Applied {
			result.AppliedRows = append(result.AppliedRows, rowIndex)
		} else {
			result.UnmatchedRows = append(result.UnmatchedRows, rowIndex)
		}
	}
	return result, nil
}
