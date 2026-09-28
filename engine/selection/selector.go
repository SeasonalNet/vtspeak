// Package selection defines the model-context to waveform-unit boundary.
package selection

import (
	"context"

	"vtspeak/engine/text"
)

// UnitRef identifies one unit in an observed model bank.
type UnitRef struct {
	Bank  string
	Index uint32
}

type Selector interface {
	Select(context.Context, []text.Context) ([]UnitRef, error)
}
