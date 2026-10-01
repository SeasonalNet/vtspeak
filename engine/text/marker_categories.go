package text

import (
	"fmt"

	"vtspeak/engine/internal/paul2013tables"
)

// FUN_10017510's category copy targets abStackY_8bd + 1 with a 999-byte
// local array, leaving 998 indexed bytes after its leading slot.
const paul2013PackedContextCategoryCapacity = 998

const paul2013PackedContextRecordCapacity = 100

const paul2013PackedContextPhoneRowStride = 7

const paul2013PackedContextNativeRecordSize = 0x3c0

// Paul2013PackedContextPhoneRecord supplies the bytes read by FUN_10017510
// from one 0x3c0-byte phone record. The fields retain their native offsets;
// the phone-symbol and marker values remain opaque.
type Paul2013PackedContextPhoneRecord struct {
	PhoneSymbols     []byte // +0x2e8
	PreviousMarkers  []byte // +0x328
	FollowingMarkers []byte // +0x329
	TerminalMarker   byte   // +0x3bd
}

// Paul2013PackedContextPhonePosition maps one flattened packed byte back to
// its caller-selected record and phone index.
type Paul2013PackedContextPhonePosition struct {
	RecordIndex int
	PhoneIndex  int
}

// Paul2013PackedContextCategories contains the initial packed bytes and their
// parallel symbol stream constructed inside FUN_10017510.
type Paul2013PackedContextCategories struct {
	CategoryBytes  []byte
	BoundaryBytes  []byte
	Positions      []Paul2013PackedContextPhonePosition
	TerminalMarker byte
}

// BuildPaul2013PackedContextCategories ports the record-to-workspace pass in
// FUN_10017510 for one already-selected group. It constructs one packed byte
// per phone, applies the two per-phone marker increments, and adds each
// record's converted terminal marker to that record's last phone. Record
// selection and ordering, produced by FUN_10017100 and its caller, remain
// explicit inputs.
func BuildPaul2013PackedContextCategories(
	records []Paul2013PackedContextPhoneRecord,
) (Paul2013PackedContextCategories, error) {
	if len(records) > paul2013PackedContextRecordCapacity {
		return Paul2013PackedContextCategories{}, fmt.Errorf("received %d context records, exceeding native capacity %d", len(records), paul2013PackedContextRecordCapacity)
	}
	result := Paul2013PackedContextCategories{}
	for recordIndex, record := range records {
		phoneCount := len(record.PhoneSymbols)
		if phoneCount == 0 {
			return Paul2013PackedContextCategories{}, fmt.Errorf("context record %d has no phone symbols", recordIndex)
		}
		if phoneCount > 65 {
			return Paul2013PackedContextCategories{}, fmt.Errorf("context record %d has %d phone symbols, exceeding native byte count", recordIndex, phoneCount)
		}
		if len(record.PreviousMarkers) != phoneCount || len(record.FollowingMarkers) != phoneCount {
			return Paul2013PackedContextCategories{}, fmt.Errorf("context record %d has %d phone symbols, %d previous markers, and %d following markers", recordIndex, phoneCount, len(record.PreviousMarkers), len(record.FollowingMarkers))
		}
		if len(result.CategoryBytes)+phoneCount > paul2013PackedContextCategoryCapacity {
			return Paul2013PackedContextCategories{}, fmt.Errorf("flattened context group exceeds native category capacity %d", paul2013PackedContextCategoryCapacity)
		}
		start := len(result.CategoryBytes)
		startClass := byte(4)
		if recordIndex > 0 {
			startClass = paul2013ContextTerminalMarkerClass(records[recordIndex-1].TerminalMarker)
		}
		for phoneIndex, symbol := range record.PhoneSymbols {
			category := startClass << 3
			if phoneIndex > 0 && record.PreviousMarkers[phoneIndex] != '0' {
				category += 8
			}
			if phoneIndex < phoneCount-1 && record.FollowingMarkers[phoneIndex] != '0' {
				category++
			}
			result.CategoryBytes = append(result.CategoryBytes, category)
			result.BoundaryBytes = append(result.BoundaryBytes, symbol)
			result.Positions = append(result.Positions, Paul2013PackedContextPhonePosition{
				RecordIndex: recordIndex,
				PhoneIndex:  phoneIndex,
			})
		}
		result.CategoryBytes[start+phoneCount-1] += paul2013ContextTerminalMarkerClass(record.TerminalMarker)
	}
	if len(records) > 0 {
		result.TerminalMarker = records[len(records)-1].TerminalMarker
	}
	return result, nil
}

func paul2013ContextTerminalMarkerClass(marker byte) byte {
	switch marker {
	case 'Z':
		return 5
	case '^':
		return 6
	case '[':
		return 4
	case 'a':
		return 3
	case ']':
		return 1
	default:
		return 2
	}
}

// PropagatePaul2013PackedContextCategories ports the bidirectional packed-byte
// propagation passes in FUN_10017510. categoryBytes correspond to the local
// packed workspace, and boundaryBytes are the parallel bytes tested through
// DAT_1007b9e0. The table values and packed three-bit fields remain opaque.
// This is a helper from FUN_10017510, not a port of that whole function.
func PropagatePaul2013PackedContextCategories(
	categoryBytes []byte,
	boundaryBytes []byte,
) ([]byte, error) {
	if len(categoryBytes) != len(boundaryBytes) {
		return nil, fmt.Errorf("received %d category bytes and %d boundary bytes", len(categoryBytes), len(boundaryBytes))
	}
	if len(categoryBytes) > paul2013PackedContextCategoryCapacity {
		return nil, fmt.Errorf("received %d packed context categories, exceeding native capacity %d", len(categoryBytes), paul2013PackedContextCategoryCapacity)
	}
	result := append([]byte(nil), categoryBytes...)
	for index := 1; index < len(result); {
		upperCategory := (result[index-1] >> 3) & 7
		if upperCategory == 0 || paul2013tables.ByteTable1007B9E0(boundaryBytes[index-1]) != 0 {
			index++
			continue
		}
		for index < len(result) {
			result[index] = result[index]&7 | upperCategory<<3
			if paul2013tables.ByteTable1007B9E0(boundaryBytes[index]) != 0 {
				index++
				break
			}
			index++
		}
	}
	for cursor := len(result) - 1; cursor > 0; cursor-- {
		lowerCategory := result[cursor] & 7
		if lowerCategory == 0 || paul2013tables.ByteTable1007B9E0(boundaryBytes[cursor]) != 0 {
			continue
		}
		for index := cursor - 1; index >= 0; index-- {
			if result[index]&7 != 0 {
				break
			}
			result[index] = result[index]&0x38 | lowerCategory
			if paul2013tables.ByteTable1007B9E0(boundaryBytes[index]) != 0 {
				break
			}
		}
	}
	return result, nil
}

// ApplyPaul2013PackedContextPhoneStateAdjustments ports the three intermediate
// per-phone mutation passes in FUN_10017510. It works on copied native records
// and uses the propagated packed bytes plus DAT_1007b9e0 boundary lookups.
// These destinations are opaque byte fields; their semantic meanings remain
// unresolved.
func ApplyPaul2013PackedContextPhoneStateAdjustments(
	inputRecords [][]byte,
	categories Paul2013PackedContextCategories,
	propagatedCategories []byte,
) ([][]byte, error) {
	count := len(categories.CategoryBytes)
	if len(categories.BoundaryBytes) != count || len(categories.Positions) != count || len(propagatedCategories) != count {
		return nil, fmt.Errorf("packed context has %d initial bytes, %d boundary bytes, %d positions, and %d propagated bytes", count, len(categories.BoundaryBytes), len(categories.Positions), len(propagatedCategories))
	}
	if len(inputRecords) > paul2013PackedContextRecordCapacity {
		return nil, fmt.Errorf("received %d input records, exceeding native capacity %d", len(inputRecords), paul2013PackedContextRecordCapacity)
	}
	records := make([][]byte, len(inputRecords))
	starts := make([]int, len(inputRecords))
	positionIndex := 0
	for recordIndex, record := range inputRecords {
		if len(record) < paul2013PackedContextNativeRecordSize {
			return nil, fmt.Errorf("input record %d has %d bytes, need %d", recordIndex, len(record), paul2013PackedContextNativeRecordSize)
		}
		phoneCount := int(record[0x95])
		if phoneCount == 0 || phoneCount > 65 {
			return nil, fmt.Errorf("input record %d has unsupported native phone count %d", recordIndex, phoneCount)
		}
		if positionIndex+phoneCount > count {
			return nil, fmt.Errorf("record %d requires %d phone positions after offset %d, but packed context has %d", recordIndex, phoneCount, positionIndex, count)
		}
		starts[recordIndex] = positionIndex
		for phoneIndex := 0; phoneIndex < phoneCount; phoneIndex++ {
			position := categories.Positions[positionIndex]
			if position.RecordIndex != recordIndex || position.PhoneIndex != phoneIndex {
				return nil, fmt.Errorf("packed phone %d maps to record/phone %d/%d, want %d/%d", positionIndex, position.RecordIndex, position.PhoneIndex, recordIndex, phoneIndex)
			}
			positionIndex++
		}
		records[recordIndex] = append([]byte(nil), record...)
	}
	if positionIndex != count {
		return nil, fmt.Errorf("native records account for %d phones, packed context contains %d", positionIndex, count)
	}

	for recordIndex := 0; recordIndex+1 < len(records); recordIndex++ {
		current := records[recordIndex]
		phoneCount := int(current[0x95])
		if current[0x3bc] == 1 {
			for phoneIndex := phoneCount - 1; phoneIndex >= 0; phoneIndex-- {
				flatIndex := starts[recordIndex] + phoneIndex
				lowCategory := propagatedCategories[flatIndex] & 7
				if lowCategory == 1 || lowCategory == 2 {
					current[0x25d+phoneIndex]++
				}
				if paul2013tables.ByteTable1007B9E0(categories.BoundaryBytes[flatIndex]) != 0 {
					break
				}
			}
		}
		next := records[recordIndex+1]
		nextPhoneCount := int(next[0x95])
		if current[0x3bc] == 1 && nextPhoneCount != 0 {
			for phoneIndex := 0; phoneIndex < nextPhoneCount; phoneIndex++ {
				flatIndex := starts[recordIndex+1] + phoneIndex
				upperCategory := (propagatedCategories[flatIndex] >> 3) & 7
				if upperCategory == 1 || upperCategory == 2 {
					next[0x25d+phoneIndex] += 10
				}
				if paul2013tables.ByteTable1007B9E0(categories.BoundaryBytes[flatIndex]) != 0 {
					break
				}
			}
		}
	}
	for _, record := range records {
		phoneCount := int(record[0x95])
		for phoneIndex := 0; phoneIndex < phoneCount; phoneIndex++ {
			if record[0x329+phoneIndex] != '2' {
				continue
			}
			if phoneIndex < phoneCount-1 && record[0x25e+phoneIndex]/10 == 0 {
				record[0x25e+phoneIndex] += 10
			}
			if record[0x25d+phoneIndex]%10 == 0 {
				record[0x25d+phoneIndex]++
			}
		}
	}
	return records, nil
}

// WritePaul2013PackedContextPhoneRows ports FUN_10017510's final per-phone
// context-byte writes. outputRecords are copied before mutation and must be
// full native 0x3c0-byte records; categories must be the initial result from
// BuildPaul2013PackedContextCategories and propagatedCategories its output
// after PropagatePaul2013PackedContextCategories. Opaque bytes at +0x2df,
// +0x96..+0x9c, and in the parallel phone stream keep their numeric meanings.
func WritePaul2013PackedContextPhoneRows(
	outputRecords [][]byte,
	categories Paul2013PackedContextCategories,
	propagatedCategories []byte,
) ([][]byte, error) {
	count := len(categories.CategoryBytes)
	if len(categories.BoundaryBytes) != count || len(categories.Positions) != count || len(propagatedCategories) != count {
		return nil, fmt.Errorf("packed context has %d initial bytes, %d boundary bytes, %d positions, and %d propagated bytes", count, len(categories.BoundaryBytes), len(categories.Positions), len(propagatedCategories))
	}
	if count > paul2013PackedContextCategoryCapacity {
		return nil, fmt.Errorf("packed context has %d phones, exceeding native capacity %d", count, paul2013PackedContextCategoryCapacity)
	}
	result := make([][]byte, len(outputRecords))
	for recordIndex, record := range outputRecords {
		if len(record) < paul2013PackedContextNativeRecordSize {
			return nil, fmt.Errorf("output record %d has %d bytes, need %d", recordIndex, len(record), paul2013PackedContextNativeRecordSize)
		}
		result[recordIndex] = append([]byte(nil), record...)
	}
	for phoneIndex, position := range categories.Positions {
		if position.RecordIndex < 0 || position.RecordIndex >= len(result) {
			return nil, fmt.Errorf("phone %d references output record %d outside %d records", phoneIndex, position.RecordIndex, len(result))
		}
		if position.PhoneIndex < 0 || position.PhoneIndex >= 65 {
			return nil, fmt.Errorf("phone %d has native record phone index %d outside [0, 65)", phoneIndex, position.PhoneIndex)
		}
	}
	for phoneIndex, position := range categories.Positions {
		row := result[position.RecordIndex]
		rowOffset := 0x96 + position.PhoneIndex*paul2013PackedContextPhoneRowStride
		leftFirst, leftSecond := byte('Z'), byte('Z')
		if phoneIndex == 1 {
			leftSecond = categories.BoundaryBytes[0]
		} else if phoneIndex >= 2 {
			leftFirst = categories.BoundaryBytes[phoneIndex-2]
			leftSecond = categories.BoundaryBytes[phoneIndex-1]
		}
		rightFirst := categories.TerminalMarker
		rightSecond := categories.TerminalMarker
		if phoneIndex+1 < count {
			rightFirst = categories.BoundaryBytes[phoneIndex+1]
		}
		if phoneIndex+2 < count {
			rightSecond = categories.BoundaryBytes[phoneIndex+2]
		}
		row[rowOffset] = leftFirst
		row[rowOffset+1] = leftSecond
		row[rowOffset+2] = categories.BoundaryBytes[phoneIndex]
		row[rowOffset+3] = rightFirst
		row[rowOffset+4] = rightSecond
		row[rowOffset+5] = propagatedCategories[phoneIndex] + categories.CategoryBytes[phoneIndex]
		row[rowOffset+6] = 0
		if row[0x2df] == 0x0c {
			row[rowOffset+6] = 0x20
		}
	}
	return result, nil
}

// Paul2013PackedContextPhoneGroupResult retains the copied output records and
// both packed-byte stages so callers can inspect the intermediate values.
type Paul2013PackedContextPhoneGroupResult struct {
	Records              [][]byte
	Categories           Paul2013PackedContextCategories
	PropagatedCategories []byte
}

// ResetPaul2013PackedContextSharedPhoneState ports FUN_10017510's initial
// reset over a contiguous 0x3c0-stride record arena. In the caller's enclosing
// group structure, record zero starts at +0x64c, so native offsets +0x6e1 and
// +0x8a9 correspond to record offsets +0x95 and +0x25d here.
func ResetPaul2013PackedContextSharedPhoneState(arena []byte, recordCount int) ([]byte, error) {
	if recordCount < 0 || recordCount > paul2013PackedContextRecordCapacity {
		return nil, fmt.Errorf("native record count %d outside [0, %d]", recordCount, paul2013PackedContextRecordCapacity)
	}
	result := append([]byte(nil), arena...)
	for recordIndex := 0; recordIndex < recordCount; recordIndex++ {
		base := recordIndex * paul2013PackedContextNativeRecordSize
		countOffset := base + 0x95
		if countOffset >= len(result) {
			return nil, fmt.Errorf("record %d reset-count byte at %#x exceeds arena size %d", recordIndex, countOffset, len(result))
		}
		if result[countOffset] == 0 {
			continue
		}
		start := base + 0x25d
		end := start + int(result[countOffset])
		if end > len(result) {
			return nil, fmt.Errorf("record %d reset span [%#x, %#x) exceeds arena size %d", recordIndex, start, end, len(result))
		}
		clear(result[start:end])
	}
	return result, nil
}

// PreparePaul2013PackedContextPhoneGroup composes the recovered record
// projection, propagation, and final phone-row writes for one caller-selected
// group of native 0x3c0-byte records. Callers with the native record stream can
// use BuildPaul2013RecordGroupDescriptors to recover the grouping boundaries.
func PreparePaul2013PackedContextPhoneGroup(
	records [][]byte,
) (Paul2013PackedContextPhoneGroupResult, error) {
	phoneRecords := make([]Paul2013PackedContextPhoneRecord, len(records))
	for recordIndex, record := range records {
		if len(record) < paul2013PackedContextNativeRecordSize {
			return Paul2013PackedContextPhoneGroupResult{}, fmt.Errorf("input record %d has %d bytes, need %d", recordIndex, len(record), paul2013PackedContextNativeRecordSize)
		}
		phoneCount := int(record[0x95])
		if phoneCount == 0 || phoneCount > 65 {
			return Paul2013PackedContextPhoneGroupResult{}, fmt.Errorf("input record %d has unsupported native phone count %d", recordIndex, phoneCount)
		}
		phoneRecords[recordIndex] = Paul2013PackedContextPhoneRecord{
			PhoneSymbols:     record[0x2e8 : 0x2e8+phoneCount],
			PreviousMarkers:  record[0x328 : 0x328+phoneCount],
			FollowingMarkers: record[0x329 : 0x329+phoneCount],
			TerminalMarker:   record[0x3bd],
		}
	}
	categories, err := BuildPaul2013PackedContextCategories(phoneRecords)
	if err != nil {
		return Paul2013PackedContextPhoneGroupResult{}, err
	}
	propagated, err := PropagatePaul2013PackedContextCategories(categories.CategoryBytes, categories.BoundaryBytes)
	if err != nil {
		return Paul2013PackedContextPhoneGroupResult{}, fmt.Errorf("propagate packed context categories: %w", err)
	}
	adjustedRecords, err := ApplyPaul2013PackedContextPhoneStateAdjustments(records, categories, propagated)
	if err != nil {
		return Paul2013PackedContextPhoneGroupResult{}, fmt.Errorf("adjust packed context phone state: %w", err)
	}
	outputRecords, err := WritePaul2013PackedContextPhoneRows(adjustedRecords, categories, propagated)
	if err != nil {
		return Paul2013PackedContextPhoneGroupResult{}, fmt.Errorf("write packed context phone rows: %w", err)
	}
	return Paul2013PackedContextPhoneGroupResult{
		Records: outputRecords, Categories: categories,
		PropagatedCategories: propagated,
	}, nil
}

// ExtractPaul2013SelectionContexts reads the seven-byte per-phone signatures
// written by FUN_10017510 from already-prepared native records. Signatures
// retain record order and then phone order, matching the packed record stream.
func ExtractPaul2013SelectionContexts(records [][]byte) ([]Context, error) {
	contexts := make([]Context, 0)
	for recordIndex, record := range records {
		if len(record) < paul2013PackedContextNativeRecordSize {
			return nil, fmt.Errorf("prepared context record %d has %d bytes, need %d", recordIndex, len(record), paul2013PackedContextNativeRecordSize)
		}
		phoneCount := int(record[0x95])
		if phoneCount == 0 || phoneCount > 65 {
			return nil, fmt.Errorf("prepared context record %d has unsupported native phone count %d", recordIndex, phoneCount)
		}
		for phoneIndex := 0; phoneIndex < phoneCount; phoneIndex++ {
			start := 0x96 + phoneIndex*paul2013PackedContextPhoneRowStride
			end := start + paul2013PackedContextPhoneRowStride
			if end > len(record) {
				return nil, fmt.Errorf("prepared context record %d phone %d requires bytes through %#x, record has %d", recordIndex, phoneIndex, end, len(record))
			}
			var context Context
			copy(context.Signature[:], record[start:end])
			contexts = append(contexts, context)
		}
	}
	return contexts, nil
}

// PreparePaul2013SelectionContextsFromPackedContextPhoneGroup connects the
// recovered marker/category pass to the seven-byte selection signatures.
// Input records must already belong to one native group; group dispatch and
// upstream record production remain outside this function.
func PreparePaul2013SelectionContextsFromPackedContextPhoneGroup(
	records [][]byte,
) (Paul2013PackedContextPhoneGroupResult, []Context, error) {
	prepared, err := PreparePaul2013PackedContextPhoneGroup(records)
	if err != nil {
		return Paul2013PackedContextPhoneGroupResult{}, nil, err
	}
	contexts, err := ExtractPaul2013SelectionContexts(prepared.Records)
	if err != nil {
		return Paul2013PackedContextPhoneGroupResult{}, nil, fmt.Errorf("extract prepared Paul 2013 selection contexts: %w", err)
	}
	return prepared, contexts, nil
}
