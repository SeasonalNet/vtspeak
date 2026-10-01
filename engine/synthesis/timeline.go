package synthesis

import (
	"errors"
	"fmt"
	"math"

	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

const (
	Paul2013TimelineCombinedView byte = iota
	Paul2013TimelineFirstSideView
	Paul2013TimelineSecondSideView
)

// Paul2013TimelineUnitView contains the selected UPM/DAT view that
// FUN_1002c120 writes into a timeline unit row. Edge periods remain raw UPM
// bytes; consumers convert them to the 16 kHz sample grid where required.
type Paul2013TimelineUnitView struct {
	Unit           selection.UnitRef
	FileIndex      byte
	Mode           byte
	DATOffset      uint32
	SampleCount    uint32
	UPMOffset      uint32
	UPMCount       byte
	LeadingPeriod  byte
	TrailingPeriod byte
}

// Paul2013TimelineControlWords holds the three values copied into a normal
// FUN_1002c220 row. Their full semantics are not established, so they remain
// explicit 32-bit words supplied by the caller.
type Paul2013TimelineControlWords struct {
	Primary uint32
	Speed   uint32
	Gain    uint32
}

// Paul2013NormalTimelineRow is the portable field-level form of one 52-byte
// normal timeline row. RowType is the observed value 2. Descriptor addresses
// are represented by Unit rather than process-local pointers.
type Paul2013NormalTimelineRow struct {
	Controls     Paul2013TimelineControlWords
	SampleCount  uint32
	PrimaryIndex int16
	Unit         selection.UnitRef
	DATOffset    uint32
	UPMOffset    uint32
	UPMCount     uint16
	LeadingSpan  uint16
	TrailingSpan uint16
	Mode         byte
	RowKind      byte
	FileIndex    byte
}

// Paul2013SelectedTimelineRowState retains the native row values that are
// produced from model-state arrays rather than from a selected unit
// reference. Mode, primary index, and control words stay explicit.
type Paul2013SelectedTimelineRowState struct {
	Mode         byte
	PrimaryIndex int16
	Controls     Paul2013TimelineControlWords
}

// Paul2013SyntheticTimelineRow contains the portable values written for a
// row-kind-1 sentence boundary. OpaqueIndex is the word copied to +0x28; the
// sentinel words retain their observed offsets without assigning meaning.
type Paul2013SyntheticTimelineRow struct {
	Controls       Paul2013TimelineControlWords
	SampleCount    uint32
	OpaqueIndex    uint16
	PrimaryIndex   int16
	SentinelAt2A   uint16
	SentinelAt2C   uint16
	RowKind        byte
	TrailingMarker byte
}

// BuildPaul2013NormalTimelineRow ports the normal-row field assignments in
// FUN_1002c220 and the unit-view extraction in FUN_1002c120. Callers supply
// the otherwise-unresolved mode, controls, primary sequence index, and bank
// file index.
func BuildPaul2013NormalTimelineRow(
	unit voice.Unit,
	reference selection.UnitRef,
	fileIndex byte,
	mode byte,
	primaryIndex int16,
	controls Paul2013TimelineControlWords,
) (Paul2013NormalTimelineRow, error) {
	view, err := BuildPaul2013TimelineUnitView(unit, reference, fileIndex, mode)
	if err != nil {
		return Paul2013NormalTimelineRow{}, err
	}
	leading := uint16(view.LeadingPeriod) * 2
	trailing := uint16(view.TrailingPeriod) * 2
	return Paul2013NormalTimelineRow{
		Controls: controls, SampleCount: view.SampleCount,
		PrimaryIndex: primaryIndex, Unit: reference,
		DATOffset: view.DATOffset, UPMOffset: view.UPMOffset,
		UPMCount: uint16(view.UPMCount), LeadingSpan: leading,
		TrailingSpan: trailing, Mode: mode, RowKind: 2,
		FileIndex: fileIndex,
	}, nil
}

// Paul2013TimelineBankFileIndex maps a selected unit's bank to the file index
// order observed by VTDTTS_MakeInfo_ENG and the Paul timeline producer.
func Paul2013TimelineBankFileIndex(bank string) (byte, error) {
	switch bank {
	case "gen":
		return 0, nil
	case "num":
		return 1, nil
	case "etc":
		return 2, nil
	case "alp":
		return 3, nil
	default:
		return 0, fmt.Errorf("unknown Paul 2013 timeline bank %q", bank)
	}
}

// BuildPaul2013NormalTimelineRowsFromUnitRefs connects selected unit
// references to the recovered normal timeline row builder. Bank file indexes
// come from the reference; mode, primary sequence index, and carried controls
// remain explicit because their source-state producers are unresolved.
func BuildPaul2013NormalTimelineRowsFromUnitRefs(
	units []selection.UnitRef,
	states []Paul2013SelectedTimelineRowState,
	reader UnitReader,
) ([]Paul2013NormalTimelineRow, error) {
	if reader == nil {
		return nil, errors.New("normal timeline row builder has no unit reader")
	}
	if len(units) == 0 {
		return nil, errors.New("normal timeline row builder has no selected units")
	}
	if len(states) != len(units) {
		return nil, fmt.Errorf("normal timeline row builder has %d selected units and %d state records", len(units), len(states))
	}
	rows := make([]Paul2013NormalTimelineRow, len(units))
	for index, reference := range units {
		fileIndex, err := Paul2013TimelineBankFileIndex(reference.Bank)
		if err != nil {
			return nil, fmt.Errorf("timeline row %d: %w", index, err)
		}
		unit, err := reader.ReadUnit(reference.Bank, reference.Index)
		if err != nil {
			return nil, fmt.Errorf("read selected timeline unit %d (%s:%d): %w", index, reference.Bank, reference.Index, err)
		}
		rows[index], err = BuildPaul2013NormalTimelineRow(
			unit, reference, fileIndex, states[index].Mode,
			states[index].PrimaryIndex, states[index].Controls,
		)
		if err != nil {
			return nil, fmt.Errorf("build selected timeline row %d (%s:%d): %w", index, reference.Bank, reference.Index, err)
		}
	}
	return rows, nil
}

// BuildPaul2013SyntheticTimelineRow ports the sentence-boundary row writes in
// FUN_1002c220 when its lookup/state inputs are supplied. Its row kind,
// primary-index sentinel, trailing sentinels, and marker are fixed by the
// observed producer; source-text production of the inputs remains unresolved.
func BuildPaul2013SyntheticTimelineRow(
	row SyntheticTimelineRow,
	opaqueIndex uint16,
	controls Paul2013TimelineControlWords,
) (Paul2013SyntheticTimelineRow, error) {
	sampleCount, err := Paul2013SyntheticTimelineSamples(row)
	if err != nil {
		return Paul2013SyntheticTimelineRow{}, err
	}
	return Paul2013SyntheticTimelineRow{
		Controls: controls, SampleCount: sampleCount, OpaqueIndex: opaqueIndex,
		PrimaryIndex: -1, SentinelAt2A: ^uint16(0), SentinelAt2C: ^uint16(0),
		RowKind: 1, TrailingMarker: 0x5a,
	}, nil
}

// BuildPaul2013TimelineUnitView ports FUN_1002c120's mode-dependent unit
// view selection. Modes 1 and 2 select the first and second UPM sides; all
// other mode bytes select the combined view, matching the native fallback.
func BuildPaul2013TimelineUnitView(
	unit voice.Unit,
	reference selection.UnitRef,
	fileIndex byte,
	mode byte,
) (Paul2013TimelineUnitView, error) {
	record := unit.Record
	first, second, err := record.UPMSides(unit.UPM)
	if err != nil {
		return Paul2013TimelineUnitView{}, fmt.Errorf("validate unit UPM sides: %w", err)
	}
	if len(unit.PCM)%2 != 0 || len(unit.PCM)/2 == 0 {
		return Paul2013TimelineUnitView{}, errors.New("timeline unit has invalid decoded PCM length")
	}
	if record.UPMEdges != [3]byte{unit.UPM[0], unit.UPM[int(record.UPMFirstCount)-1], unit.UPM[len(unit.UPM)-1]} {
		return Paul2013TimelineUnitView{}, errors.New("timeline unit cached UPM edges do not match its combined vector")
	}
	firstSamples, err := sumPaul2013UPMSamples(first)
	if err != nil {
		return Paul2013TimelineUnitView{}, fmt.Errorf("sum first UPM side: %w", err)
	}
	secondSamples, err := sumPaul2013UPMSamples(second)
	if err != nil {
		return Paul2013TimelineUnitView{}, fmt.Errorf("sum second UPM side: %w", err)
	}
	if firstSamples != uint64(record.FirstSideSamples) || secondSamples != uint64(record.SecondSideSamples) {
		return Paul2013TimelineUnitView{}, errors.New("timeline unit UPM sides do not match declared side sample counts")
	}
	sharedSamples := uint64(record.UPMEdges[1]) * 2
	if sharedSamples > firstSamples || sharedSamples > secondSamples {
		return Paul2013TimelineUnitView{}, errors.New("shared UPM period exceeds a side sample count")
	}
	combinedSamples := firstSamples + secondSamples - sharedSamples
	if combinedSamples != uint64(len(unit.PCM)/2) {
		return Paul2013TimelineUnitView{}, errors.New("timeline unit UPM views do not match decoded PCM sample count")
	}

	view := Paul2013TimelineUnitView{
		Unit: reference, FileIndex: fileIndex, Mode: mode,
		DATOffset: record.DATOffset, UPMOffset: record.UPMOffset,
	}
	switch mode {
	case Paul2013TimelineFirstSideView:
		view.SampleCount = uint32(record.FirstSideSamples)
		view.UPMCount = record.UPMFirstCount
		view.LeadingPeriod = record.UPMEdges[0]
		view.TrailingPeriod = record.UPMEdges[1]
	case Paul2013TimelineSecondSideView:
		datOffset := uint64(record.DATOffset) + uint64(record.FirstSideSamples) - sharedSamples
		upmOffset := uint64(record.UPMOffset) + uint64(record.UPMFirstCount) - 1
		if datOffset > uint64(^uint32(0)) || upmOffset > uint64(^uint32(0)) {
			return Paul2013TimelineUnitView{}, errors.New("second-side timeline offset exceeds unsigned 32-bit range")
		}
		view.DATOffset = uint32(datOffset)
		view.SampleCount = uint32(record.SecondSideSamples)
		view.UPMOffset = uint32(upmOffset)
		view.UPMCount = record.UPMSecondCount
		view.LeadingPeriod = record.UPMEdges[1]
		view.TrailingPeriod = record.UPMEdges[2]
	default:
		count := uint16(record.UPMFirstCount) + uint16(record.UPMSecondCount) - 1
		if count > 255 {
			return Paul2013TimelineUnitView{}, errors.New("combined UPM view exceeds unsigned-byte period count")
		}
		view.SampleCount = uint32(combinedSamples)
		view.UPMCount = byte(count)
		view.LeadingPeriod = record.UPMEdges[0]
		view.TrailingPeriod = record.UPMEdges[2]
	}
	return view, nil
}

func sumPaul2013UPMSamples(periods []byte) (uint64, error) {
	if len(periods) == 0 {
		return 0, errors.New("UPM side is empty")
	}
	var samples uint64
	for _, period := range periods {
		samples += uint64(period) * 2
	}
	return samples, nil
}

// TimelineRow is the sample-window portion of one normal selected-unit row.
// Counts are 16 kHz samples. The legacy builder can also emit synthetic rows
// at boundaries; those are represented separately by SyntheticTimelineRow.
type TimelineRow struct {
	Unit         selection.UnitRef
	SampleCount  int
	LeadingSpan  int
	TrailingSpan int
}

// TimelineOutputRow contains exactly one normal selected-unit row or one
// synthetic boundary row. It lets output sizing preserve each row kind's
// different observed length source.
type TimelineOutputRow struct {
	Normal    *TimelineRow
	Synthetic *SyntheticTimelineRow
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
		view, err := BuildPaul2013TimelineUnitView(unit, reference, 0, Paul2013TimelineCombinedView)
		if err != nil {
			return nil, fmt.Errorf("build timeline unit %d (%s:%d): %w", index, reference.Bank, reference.Index, err)
		}
		rows[index] = TimelineRow{
			Unit: reference, SampleCount: int(view.SampleCount),
			LeadingSpan: int(view.LeadingPeriod) * 2, TrailingSpan: int(view.TrailingPeriod) * 2,
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
	outputRows := make([]TimelineOutputRow, len(rows))
	for index := range rows {
		outputRows[index].Normal = &rows[index]
	}
	return Paul2013TimelineOutputFramesWithSynthetic(outputRows)
}

// Paul2013TimelineOutputFramesWithSynthetic sizes an ordered mix of normal
// unit rows and synthetic sentence-boundary rows. Normal non-final rows
// advance by sample count minus trailing UPM span; final normal rows use their
// full sample count. Synthetic rows use FUN_1002c220's duration expression.
// Row ordering and this composition have not been compared against every
// boundary form in the native renderer.
func Paul2013TimelineOutputFramesWithSynthetic(rows []TimelineOutputRow) (int, error) {
	if len(rows) == 0 {
		return 0, errors.New("timeline has no rows")
	}
	total := int64(0)
	for index, row := range rows {
		if (row.Normal == nil) == (row.Synthetic == nil) {
			return 0, fmt.Errorf("timeline row %d must contain exactly one row kind", index)
		}
		advance := int64(0)
		if row.Normal != nil {
			normal := row.Normal
			if normal.SampleCount <= 0 || normal.LeadingSpan < 0 || normal.TrailingSpan < 0 ||
				normal.LeadingSpan > normal.SampleCount || normal.TrailingSpan > normal.SampleCount {
				return 0, fmt.Errorf("timeline row %d has invalid sample count or edge span", index)
			}
			advance = int64(normal.SampleCount)
			if index+1 < len(rows) {
				advance -= int64(normal.TrailingSpan)
			}
		} else {
			samples, err := Paul2013SyntheticTimelineSamples(*row.Synthetic)
			if err != nil {
				return 0, fmt.Errorf("timeline row %d synthetic duration: %w", index, err)
			}
			advance = int64(samples)
		}
		if advance > int64(math.MaxInt)-total {
			return 0, errors.New("timeline output length exceeds the platform integer range")
		}
		total += advance
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
