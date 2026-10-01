package text

import "fmt"

const paul2013RecordGroupCapacity = 100

const paul2013RecordGroupRowSpanOffset = 0x94

const paul2013RecordGroupTerminalOffset = 0x3bd

const paul2013RecordGroupOutputShortBase = 0xbea6

const paul2013RecordGroupOutputShortStride = 0x0f

// Paul2013RecordGroupDescriptor is the value view of one FUN_10012df0
// descriptor. BoundaryRecordIndex is the native +0xe field, which stores the
// group's starting record index. OpaqueRowSpan is the sum of record byte
// +0x94; OutputShortOffset is the pointer offset from the caller's short-array
// base.
type Paul2013RecordGroupDescriptor struct {
	StartRecordIndex    int
	RecordCount         int
	BoundaryRecordIndex int
	OpaqueRowSpan       uint16
	OutputShortOffset   int
}

// BuildPaul2013RecordGroupDescriptors ports FUN_10012df0's record grouping
// and metadata-span construction. A '\\' or ']' terminal keeps a record in
// the current group; any other terminal closes a group after that record,
// except that the final record always belongs to the final group.
func BuildPaul2013RecordGroupDescriptors(records [][]byte) ([]Paul2013RecordGroupDescriptor, error) {
	if len(records) > paul2013RecordGroupCapacity {
		return nil, fmt.Errorf("received %d native records, exceeding group capacity %d", len(records), paul2013RecordGroupCapacity)
	}
	for recordIndex, record := range records {
		if len(record) < paul2013MarkerRecordSize {
			return nil, fmt.Errorf("group record %d has %d bytes, need %d", recordIndex, len(record), paul2013MarkerRecordSize)
		}
	}
	if len(records) == 0 {
		return []Paul2013RecordGroupDescriptor{}, nil
	}

	groups := make([]Paul2013RecordGroupDescriptor, 0, len(records))
	start := 0
	for recordIndex := 0; recordIndex < len(records)-1; recordIndex++ {
		terminal := records[recordIndex][paul2013RecordGroupTerminalOffset]
		if terminal == '\\' || terminal == ']' {
			continue
		}
		groups = append(groups, paul2013RecordGroupDescriptor(records, start, recordIndex+1, start))
		start = recordIndex + 1
	}
	groups = append(groups, paul2013RecordGroupDescriptor(records, start, len(records), start))

	rowPrefix := 0
	for index := range groups {
		groups[index].OutputShortOffset = paul2013RecordGroupOutputShortBase + rowPrefix*paul2013RecordGroupOutputShortStride
		rowPrefix += int(groups[index].OpaqueRowSpan)
	}
	return groups, nil
}

func paul2013RecordGroupDescriptor(records [][]byte, start, end, boundaryRecordIndex int) Paul2013RecordGroupDescriptor {
	var rowSpan uint16
	for _, record := range records[start:end] {
		rowSpan += uint16(record[paul2013RecordGroupRowSpanOffset])
	}
	return Paul2013RecordGroupDescriptor{
		StartRecordIndex:    start,
		RecordCount:         end - start,
		BoundaryRecordIndex: boundaryRecordIndex,
		OpaqueRowSpan:       rowSpan,
	}
}

// Paul2013PreparedRecordGroups retains the all-record marker results, the
// recovered group descriptors, and each group's packed phone-context output.
type Paul2013PreparedRecordGroups struct {
	Records           [][]byte
	Descriptors       []Paul2013RecordGroupDescriptor
	ContextGroups     []Paul2013PackedContextPhoneGroupResult
	SelectionContexts [][]Context
	ResetArena        []byte
}

// PreparePaul2013PhoneMarkerRecordGroups composes the native record-group
// descriptors, FUN_10017100 marker prepass, and per-group FUN_10017510 packed
// context construction. specialKeys are the caller-resolved bytes reached by
// the native pointer table and are needed only for eligible terminal joins.
func PreparePaul2013PhoneMarkerRecordGroups(
	records [][]byte,
	specialKeys [][]byte,
	characterMap [256]byte,
) (Paul2013PreparedRecordGroups, error) {
	return preparePaul2013PhoneMarkerRecordGroups(records, func() ([][]byte, error) {
		return ApplyPaul2013PhoneMarkerPrepass(records, specialKeys, characterMap)
	})
}

// PreparePaul2013PhoneMarkerRecordGroupsWithPointerResolver composes the
// marker pass with the native record-group and packed-context stages while
// resolving +0x3b8 pointers through a caller-owned address space.
func PreparePaul2013PhoneMarkerRecordGroupsWithPointerResolver(
	records [][]byte,
	resolvePointer Paul2013CStringPointerResolver,
	characterMap [256]byte,
) (Paul2013PreparedRecordGroups, error) {
	return preparePaul2013PhoneMarkerRecordGroups(records, func() ([][]byte, error) {
		return ApplyPaul2013PhoneMarkerPrepassWithPointerResolver(records, resolvePointer, characterMap)
	})
}

func preparePaul2013PhoneMarkerRecordGroups(
	records [][]byte,
	preprocess func() ([][]byte, error),
) (Paul2013PreparedRecordGroups, error) {
	descriptors, err := BuildPaul2013RecordGroupDescriptors(records)
	if err != nil {
		return Paul2013PreparedRecordGroups{}, err
	}
	preprocessed, err := preprocess()
	if err != nil {
		return Paul2013PreparedRecordGroups{}, fmt.Errorf("preprocess record phone markers: %w", err)
	}
	groups := make([]Paul2013PackedContextPhoneGroupResult, len(descriptors))
	selectionContexts := make([][]Context, len(descriptors))
	for groupIndex, descriptor := range descriptors {
		start := descriptor.StartRecordIndex
		end := start + descriptor.RecordCount
		groups[groupIndex], err = PreparePaul2013PackedContextPhoneGroup(preprocessed[start:end])
		if err != nil {
			return Paul2013PreparedRecordGroups{}, fmt.Errorf("prepare record group %d: %w", groupIndex, err)
		}
		selectionContexts[groupIndex], err = ExtractPaul2013SelectionContexts(groups[groupIndex].Records)
		if err != nil {
			return Paul2013PreparedRecordGroups{}, fmt.Errorf("extract selection contexts for record group %d: %w", groupIndex, err)
		}
	}
	return Paul2013PreparedRecordGroups{
		Records: preprocessed, Descriptors: descriptors, ContextGroups: groups,
		SelectionContexts: selectionContexts,
	}, nil
}

// PreparePaul2013PhoneMarkerRecordGroupsFromArena also applies FUN_10017510's
// initial shared-state reset when the caller supplies the contiguous native
// record arena. The reset copy is retained for inspection; the later
// per-group phone-row outputs remain in ContextGroups, matching the native
// stage boundary.
func PreparePaul2013PhoneMarkerRecordGroupsFromArena(
	arena []byte,
	recordCount int,
	specialKeys [][]byte,
	characterMap [256]byte,
) (Paul2013PreparedRecordGroups, error) {
	return preparePaul2013PhoneMarkerRecordGroupsFromArena(arena, recordCount, func(records [][]byte) (Paul2013PreparedRecordGroups, error) {
		return PreparePaul2013PhoneMarkerRecordGroups(records, specialKeys, characterMap)
	})
}

// PreparePaul2013PhoneMarkerRecordGroupsFromArenaWithPointerResolver is the
// contiguous-arena variant for native +0x3b8 key pointers.
func PreparePaul2013PhoneMarkerRecordGroupsFromArenaWithPointerResolver(
	arena []byte,
	recordCount int,
	resolvePointer Paul2013CStringPointerResolver,
	characterMap [256]byte,
) (Paul2013PreparedRecordGroups, error) {
	return preparePaul2013PhoneMarkerRecordGroupsFromArena(arena, recordCount, func(records [][]byte) (Paul2013PreparedRecordGroups, error) {
		return PreparePaul2013PhoneMarkerRecordGroupsWithPointerResolver(records, resolvePointer, characterMap)
	})
}

func preparePaul2013PhoneMarkerRecordGroupsFromArena(
	arena []byte,
	recordCount int,
	prepare func([][]byte) (Paul2013PreparedRecordGroups, error),
) (Paul2013PreparedRecordGroups, error) {
	if recordCount < 0 || recordCount > paul2013RecordGroupCapacity {
		return Paul2013PreparedRecordGroups{}, fmt.Errorf("native record count %d outside [0, %d]", recordCount, paul2013RecordGroupCapacity)
	}
	needed := recordCount * paul2013MarkerRecordSize
	if len(arena) < needed {
		return Paul2013PreparedRecordGroups{}, fmt.Errorf("record arena has %d bytes, need at least %d for %d records", len(arena), needed, recordCount)
	}
	resetArena, err := ResetPaul2013PackedContextSharedPhoneState(arena, recordCount)
	if err != nil {
		return Paul2013PreparedRecordGroups{}, fmt.Errorf("reset shared phone state: %w", err)
	}
	records := make([][]byte, recordCount)
	for recordIndex := range records {
		start := recordIndex * paul2013MarkerRecordSize
		records[recordIndex] = resetArena[start : start+paul2013MarkerRecordSize]
	}
	prepared, err := prepare(records)
	if err != nil {
		return Paul2013PreparedRecordGroups{}, err
	}
	prepared.ResetArena = resetArena
	return prepared, nil
}
