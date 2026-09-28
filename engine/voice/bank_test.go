package voice

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLocalPaulUnitWhenAssetsArePresent(t *testing.T) {
	root := filepath.Join("..", "..", "data-paul", "M16")
	indexPath := filepath.Join(root, "mc_idx_tbl", "unit-gen.idx")
	datPath := filepath.Join(root, "dat", "merged-gen.dat")
	upmPath := filepath.Join(root, "dat", "merged-gen.upm")
	for _, path := range []string{indexPath, datPath, upmPath} {
		if _, err := os.Stat(path); err != nil {
			t.Skip("local Paul model files are unavailable")
		}
	}
	bank, err := OpenBank(indexPath, datPath, upmPath)
	if err != nil {
		t.Fatal(err)
	}
	defer bank.Close()
	unit, err := bank.ReadUnit(0)
	if err != nil {
		t.Fatal(err)
	}
	wantSamples := 0
	for _, period := range unit.UPM {
		wantSamples += int(period) * 2
	}
	if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 || len(unit.UPM) == 0 || len(unit.PCM)/2 != wantSamples {
		t.Fatalf("invalid decoded unit: signature=%d pcm=%d upm=%d", len(unit.Record.Signature), len(unit.PCM), len(unit.UPM))
	}
}

func TestOpenPaul2013ModelWhenAssetsArePresent(t *testing.T) {
	dataRoot := filepath.Join("..", "..", "data-paul", "M16")
	if _, err := os.Stat(filepath.Join(dataRoot, "dblist.idx")); err != nil {
		t.Skip("local Paul model files are unavailable")
	}
	model, err := OpenPaul2013(dataRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer model.Close()
	wantCounts := map[string]uint32{"gen": 440124, "num": 24508, "etc": 115723, "alp": 119}
	for bankName, wantCount := range wantCounts {
		bank := model.Banks[bankName]
		var gotCount uint32
		if bank != nil {
			gotCount = bank.UnitCount()
		}
		if gotCount != wantCount {
			t.Errorf("%s bank count = %d, want %d", bankName, gotCount, wantCount)
		}
	}
	if len(model.Trees.Duration) != 9 || len(model.Trees.Pitch) != 8 {
		t.Fatalf("loaded %d duration and %d pitch trees", len(model.Trees.Duration), len(model.Trees.Pitch))
	}
}
