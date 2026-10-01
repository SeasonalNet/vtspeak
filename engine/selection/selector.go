// Package selection defines the model-context to waveform-unit boundary.
package selection

import (
	"context"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

// UnitRef identifies one unit in an observed model bank.
type UnitRef struct {
	Bank  string
	Index uint32
}

type Selector interface {
	Select(context.Context, []text.Context) ([]UnitRef, error)
}

// UniqueExactSelector connects the recovered exact-context class lookup to
// the pipeline for contexts whose indexed class contains exactly one unit.
// It does not rank multi-unit classes, broaden a missing exact match, or
// invent the query and path state used by the native selector.
type UniqueExactSelector struct {
	Catalog *voice.ClassCatalog
}

// Select resolves each context only when its exact class has one member.
// Empty or ambiguous classes fail closed because choosing among those cases
// requires the still-incomplete native candidate and path-scoring inputs.
func (selector UniqueExactSelector) Select(
	ctx context.Context,
	contexts []text.Context,
) ([]UnitRef, error) {
	if ctx == nil {
		return nil, fmt.Errorf("exact selector has no context")
	}
	if selector.Catalog == nil {
		return nil, fmt.Errorf("exact selector has no class catalog")
	}
	if len(contexts) == 0 {
		return nil, fmt.Errorf("exact selector received no phone contexts")
	}
	selected := make([]UnitRef, len(contexts))
	for index, phoneContext := range contexts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		candidates, found, err := ExactContextUnitCandidates(selector.Catalog, phoneContext, 2)
		if err != nil {
			return nil, fmt.Errorf("lookup exact class for phone context %d: %w", index, err)
		}
		if !found || len(candidates) == 0 {
			return nil, fmt.Errorf("phone context %d has no exact unit class", index)
		}
		if len(candidates) != 1 {
			return nil, fmt.Errorf("phone context %d has multiple exact unit candidates; native ranking is required", index)
		}
		selected[index] = candidates[0]
	}
	return selected, nil
}
