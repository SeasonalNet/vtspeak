package synthesis

import (
	"errors"
	"fmt"
	"math"

	"vtspeak/engine/selection"
)

// TimelineRow is the sample-window portion of one normal selected-unit row.
// Counts are 16 kHz samples. The legacy builder can also emit synthetic rows
// at boundaries; those are represented separately by SyntheticTimelineRow.
type TimelineRow struct {
	Unit         selection.UnitRef
	SampleCount  int
	LeadingSpan  int
	TrailingSpan int
}

// BuildPaul2013TimelineRows derives normal row lengths and edge spans from
// selected units. The source trace maps normal row length to the decoded unit
// sample span and maps both edge spans to the first/last UPM periods, doubled
// onto the 16 kHz sample grid. It does not determine row kind or insert
// synthetic boundary rows.
func BuildPaul2013TimelineRows(units []selection.UnitRef, reader UnitReader) ([]TimelineRow, error) {
	if reader == nil {
		return nil, errors.New("timeline builder has no unit reader")
	}
	if len(units) == 0 {
		return nil, errors.New("timeline requires at least one selected unit")
	}
	rows := make([]TimelineRow, len(units))
	for index, reference := range units {
		unit, err := reader.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return nil, fmt.Errorf("read timeline unit %d (%s:%d): %w", index, reference.Bank, reference.Index, err)
		}
		if len(unit.PCM) == 0 || len(unit.PCM)%2 != 0 {
			return nil, fmt.Errorf("timeline unit %d (%s:%d) has invalid PCM byte length %d", index, reference.Bank, reference.Index, len(unit.PCM))
		}
		if len(unit.UPM) == 0 {
			return nil, fmt.Errorf("timeline unit %d (%s:%d) has no UPM periods", index, reference.Bank, reference.Index)
		}
		leading := int(unit.UPM[0]) * 2
		trailing := int(unit.UPM[len(unit.UPM)-1]) * 2
		sampleCount := len(unit.PCM) / 2
		if leading > sampleCount || trailing > sampleCount {
			return nil, fmt.Errorf("timeline unit %d (%s:%d) edge span exceeds its decoded sample count", index, reference.Bank, reference.Index)
		}
		rows[index] = TimelineRow{
			Unit: reference, SampleCount: sampleCount,
			LeadingSpan: leading, TrailingSpan: trailing,
		}
	}
	return rows, nil
}

// Paul2013TimelineOutputFrames applies the observed cursor accounting for a
// sequence of normal timeline rows: each non-final row advances by its sample
// count minus its trailing span; the final row emits its full sample count.
// The 2013 join routine has a seven-frame residual against this accounting in
// one stable capture, so this is a validated plan length, not a claim of exact
// renderer output length for every path.
func Paul2013TimelineOutputFrames(rows []TimelineRow) (int, error) {
	if len(rows) == 0 {
		return 0, errors.New("timeline has no rows")
	}
	total := int64(0)
	for index, row := range rows {
		if row.SampleCount <= 0 || row.LeadingSpan < 0 || row.TrailingSpan < 0 ||
			row.LeadingSpan > row.SampleCount || row.TrailingSpan > row.SampleCount {
			return 0, fmt.Errorf("timeline row %d has invalid sample count or edge span", index)
		}
		advance := row.SampleCount
		if index+1 < len(rows) {
			advance -= row.TrailingSpan
		}
		total += int64(advance)
		if total > math.MaxInt {
			return 0, errors.New("timeline output length exceeds the platform integer range")
		}
	}
	return int(total), nil
}

// SyntheticTimelineRow carries the two source values used by the observed
// sentence-boundary duration expression. Their semantic names are not
// established, so they remain numeric rather than being called pause rate or
// duration.
type SyntheticTimelineRow struct {
	LookupValue uint32
	StateValue  uint32
}

// Paul2013SyntheticTimelineSamples ports the captured integer expression
// ((((v >> 1) + 10000) / v) * q + 50) / 100. It rejects v=0 and overflow
// instead of emulating unchecked 32-bit wraparound.
func Paul2013SyntheticTimelineSamples(row SyntheticTimelineRow) (uint32, error) {
	if row.LookupValue == 0 {
		return 0, errors.New("synthetic timeline lookup value is zero")
	}
	ratio := (uint64(row.LookupValue>>1) + 10000) / uint64(row.LookupValue)
	value := ratio * uint64(row.StateValue)
	if value > math.MaxUint64-50 {
		return 0, errors.New("synthetic timeline duration arithmetic overflows")
	}
	value = (value + 50) / 100
	if value > math.MaxUint32 {
		return 0, errors.New("synthetic timeline sample count exceeds 32-bit range")
	}
	return uint32(value), nil
}
