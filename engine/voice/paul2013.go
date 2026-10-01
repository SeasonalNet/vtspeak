package voice

import (
	"errors"
	"fmt"
	"path/filepath"

	"vtspeak/engine/dat"
	"vtspeak/engine/distance"
	"vtspeak/engine/tree3"
)

var paul2013BankNames = [...]string{"gen", "num", "etc", "alp"}

type Paul2013 struct {
	Banks            map[string]*Bank
	Trees            *tree3.Catalog
	Distances        *distance.Table
	FeatureDistances *distance.FeatureTable
}

// OpenPaul2013 loads the four documented unit banks and all duration/pitch
// trees beneath a local 2013 M16 Paul data root. Inputs are opened read-only.
func OpenPaul2013(dataRoot string) (*Paul2013, error) {
	model := &Paul2013{Banks: make(map[string]*Bank, len(paul2013BankNames))}
	for _, name := range paul2013BankNames {
		bank, err := OpenBank(
			filepath.Join(dataRoot, "mc_idx_tbl", "unit-"+name+".idx"),
			filepath.Join(dataRoot, "dat", "merged-"+name+".dat"),
			filepath.Join(dataRoot, "dat", "merged-"+name+".upm"),
		)
		if err != nil {
			_ = model.Close()
			return nil, fmt.Errorf("open Paul %s bank: %w", name, err)
		}
		model.Banks[name] = bank
	}
	trees, err := tree3.LoadPaul2013(dataRoot)
	if err != nil {
		_ = model.Close()
		return nil, fmt.Errorf("load Paul decision trees: %w", err)
	}
	model.Trees = trees
	distances, err := distance.LoadFile(filepath.Join(dataRoot, "ttsdata", "dist_tbl", "cepdist.tbl"))
	if err != nil {
		_ = model.Close()
		return nil, fmt.Errorf("load Paul metric-distance table: %w", err)
	}
	if distances.Size() != 1024 {
		_ = model.Close()
		return nil, fmt.Errorf("Paul metric-distance table has %d entries, want 1024", distances.Size())
	}
	model.Distances = distances
	model.FeatureDistances = distance.NewFeatureTable()
	return model, nil
}

func (p *Paul2013) ReadUnit(bank string, index uint32) (Unit, error) {
	resource, ok := p.Banks[bank]
	if !ok {
		return Unit{}, fmt.Errorf("unknown Paul unit bank %q", bank)
	}
	return resource.ReadUnit(index)
}

// ReadRecord reads indexed metadata without accessing DAT or UPM payloads.
func (p *Paul2013) ReadRecord(bank string, index uint32) (dat.UnitRecord, error) {
	if p == nil {
		return dat.UnitRecord{}, errors.New("unit record lookup has no Paul model")
	}
	resource, ok := p.Banks[bank]
	if !ok || resource == nil {
		return dat.UnitRecord{}, fmt.Errorf("unknown or unavailable Paul unit bank %q", bank)
	}
	return resource.ReadRecord(index)
}

func (p *Paul2013) Close() error {
	var failures []error
	for _, name := range paul2013BankNames {
		if bank, ok := p.Banks[name]; ok {
			if err := bank.Close(); err != nil {
				failures = append(failures, fmt.Errorf("close %s bank: %w", name, err))
			}
		}
	}
	return errors.Join(failures...)
}
