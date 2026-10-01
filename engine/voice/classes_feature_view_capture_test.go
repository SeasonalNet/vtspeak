package voice

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/text"
)

func TestPaul2013FeatureViewPrefixMatchesCapturedPHalfKeys(t *testing.T) {
	dataRoot := filepath.Join("..", "..", "data-paul", "M16")
	if _, err := os.Stat(filepath.Join(dataRoot, "dblist.idx")); err != nil {
		t.Skip("local Paul model files are unavailable")
	}
	model, err := OpenPaul2013(dataRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer model.Close()
	catalog, err := BuildPaul2013ClassCatalog(context.Background(), model)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name       string
		signature  [7]byte
		mode1View  [10]byte
		mode2View  [10]byte
		mode1Count int
		mode2Count int
	}{
		{
			name:       "P with low class-field value 1",
			signature:  [7]byte{90, 19, 7, 57, 7, 97, 0},
			mode1View:  [10]byte{0, 0, 7, 3, 0, 19, 40, 19, 1, 4},
			mode2View:  [10]byte{1, 0, 7, 4, 0, 57, 10, 57, 4, 3},
			mode1Count: 6,
			mode2Count: 4,
		},
		{
			name:       "P with low class-field value 1 and upper field 2",
			signature:  [7]byte{90, 19, 7, 57, 7, 81, 0},
			mode1View:  [10]byte{0, 0, 7, 3, 0, 19, 20, 19, 1, 4},
			mode2View:  [10]byte{1, 0, 7, 4, 0, 57, 10, 57, 2, 3},
			mode1Count: 8,
			mode2Count: 4,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			query := text.Context{Signature: test.signature}
			for _, mode := range []byte{1, 2} {
				view, err := query.Paul2013FeatureView(mode)
				if err != nil {
					t.Fatal(err)
				}
				wantView, wantCount := test.mode1View, test.mode1Count
				if mode == 2 {
					wantView, wantCount = test.mode2View, test.mode2Count
				}
				if view != wantView {
					t.Fatalf("mode %d feature view = %v, want captured %v", mode, view, wantView)
				}
				for _, width := range []int{10, 9, 8} {
					classes, err := catalog.LookupFeatureViewPrefix(view, mode, width, 100)
					if err != nil {
						t.Fatalf("mode %d lookup width %d: %v", mode, width, err)
					}
					want := wantCount
					if width > 8 {
						want = 0
					}
					if len(classes) != want {
						t.Errorf("mode %d width %d range count = %d, want captured %d", mode, width, len(classes), want)
					}
				}
			}
		})
	}
}
