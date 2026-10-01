package duration

import (
	"context"
	"testing"

	"vtspeak/engine/text"
)

func TestModelTokenBoundaryEngineRequiresSharedCatalog(t *testing.T) {
	for _, engine := range []*Engine{nil, {}} {
		if _, err := engine.RunPaul2013PositionSegmentBoundaryDriver(context.Background(), nil, nil, nil, nil, nil, nil, nil, 0, nil); err == nil {
			t.Fatal("accepted unloaded segment boundary engine")
		}
		if _, err := engine.RunPaul2013ModelTokenBoundaries(text.Paul2013FinalizedModelState{}, 0, 0, nil, nil); err == nil {
			t.Fatal("accepted unloaded boundary engine")
		}
		if _, err := engine.RunPaul2013PositionStateBoundaries(text.Paul2013FinalizedModelState{}, nil, text.Paul2013PositionStateProgram{}, 0, 0, nil); err == nil {
			t.Fatal("accepted unloaded position/boundary engine")
		}
		if _, err := engine.RunPaul2013PositionDescriptorBoundaries(text.Paul2013FinalizedModelState{}, nil, nil, nil, 0, 0, nil); err == nil {
			t.Fatal("accepted unloaded descriptor/boundary engine")
		}
		if _, err := engine.RunPaul2013NativePositionWorkspaceBoundaries(text.Paul2013FinalizedModelState{}, nil, nil, nil, 0, 0, nil); err == nil {
			t.Fatal("accepted unloaded native workspace engine")
		}
		if _, err := engine.RunPaul2013PreparedPositionWorkspaceBoundaries(text.Paul2013FinalizedModelState{}, nil, nil, nil, nil, 0, 0, nil); err == nil {
			t.Fatal("accepted unloaded prepared workspace engine")
		}
	}
}
