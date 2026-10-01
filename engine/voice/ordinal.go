package voice

import (
	"errors"
	"fmt"
	"math"
)

// GlobalUnitOrdinal maps a bank-local unit reference into the concatenated
// unit order used by the model loader. The order follows dblist.idx.
func (p *Paul2013) GlobalUnitOrdinal(bank string, index uint32) (uint32, error) {
	if p == nil {
		return 0, errors.New("global unit ordinal has no Paul model")
	}
	var offset uint64
	for _, name := range paul2013BankNames {
		resource, ok := p.Banks[name]
		if !ok {
			return 0, fmt.Errorf("Paul %s bank is unavailable", name)
		}
		if resource == nil {
			return 0, fmt.Errorf("Paul %s bank is nil", name)
		}
		count := resource.UnitCount()
		if count == 0 {
			return 0, fmt.Errorf("Paul %s bank has no units", name)
		}
		if name == bank {
			if index >= count {
				return 0, fmt.Errorf("unit %s:%d outside bank range 0..%d", bank, index, count-1)
			}
			ordinal := offset + uint64(index)
			if ordinal > math.MaxUint32 {
				return 0, errors.New("global unit ordinal exceeds 32-bit range")
			}
			return uint32(ordinal), nil
		}
		offset += uint64(count)
	}
	return 0, fmt.Errorf("unknown Paul unit bank %q", bank)
}

// UnitAtGlobalOrdinal converts the concatenated dblist.idx unit order back to
// its bank-local reference. The returned reference can be used to read the
// indexed signature without decoding waveform data.
func (p *Paul2013) UnitAtGlobalOrdinal(ordinal uint32) (string, uint32, error) {
	if p == nil {
		return "", 0, errors.New("global unit lookup has no Paul model")
	}
	var offset uint64
	for _, name := range paul2013BankNames {
		resource, ok := p.Banks[name]
		if !ok {
			return "", 0, fmt.Errorf("Paul %s bank is unavailable", name)
		}
		if resource == nil {
			return "", 0, fmt.Errorf("Paul %s bank is nil", name)
		}
		count := uint64(resource.UnitCount())
		if count == 0 {
			return "", 0, fmt.Errorf("Paul %s bank has no units", name)
		}
		if uint64(ordinal) >= offset && uint64(ordinal) < offset+count {
			return name, uint32(uint64(ordinal) - offset), nil
		}
		offset += count
	}
	return "", 0, fmt.Errorf("global unit ordinal %d is outside the Paul model", ordinal)
}

func globalUnitOrdinal(order []string, counts map[string]uint32, bank string, index uint32) (uint32, error) {
	var offset uint64
	for _, name := range order {
		count, ok := counts[name]
		if !ok {
			return 0, fmt.Errorf("Paul %s bank count is unavailable", name)
		}
		if count == 0 {
			return 0, fmt.Errorf("Paul %s bank has no units", name)
		}
		if name == bank {
			if index >= count {
				return 0, fmt.Errorf("unit %s:%d outside bank range 0..%d", bank, index, count-1)
			}
			ordinal := offset + uint64(index)
			if ordinal > math.MaxUint32 {
				return 0, errors.New("global unit ordinal exceeds 32-bit range")
			}
			return uint32(ordinal), nil
		}
		offset += uint64(count)
	}
	return 0, fmt.Errorf("unknown Paul unit bank %q", bank)
}
