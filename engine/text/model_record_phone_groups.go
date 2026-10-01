package text

import "fmt"

const paul2013ModelRecordPhoneSymbolOffset = 0x2e8

// Paul2013ModelRecordPhoneGroups is the Go projection of FUN_10012c70 for one
// native model-state record. Groups are ordered by record marker span, then by
// vowel group within that span. NativeGroupRows contains the emitted 0x1e-byte
// rows; a final no-vowel fallback row may occupy the uncounted tail slot, as in
// FUN_10013c00. PhoneLabels is the parallel per-phone 1/2/3 stream.
type Paul2013ModelRecordPhoneGroups struct {
	Groups           []Paul2013PhoneGroup
	DurationContexts []Paul2013PhoneDurationContext
	NativeGroupCount int
	NativeGroupRows  [][]byte
	PhoneLabels      []byte
}

// BuildPaul2013ModelRecordPhoneGroups ports FUN_10012c70's subdivision and
// FUN_10013c00 row and label writes for caller-decoded phones. The native
// record supplies its phone count and +0x329 marker stream; use the FromSymbols
// entry point to decode the record's +0x2e8 bytes through the static codebook.
func BuildPaul2013ModelRecordPhoneGroups(
	record []byte,
	phones []CMUPhone,
) (Paul2013ModelRecordPhoneGroups, error) {
	if len(record) < paul2013MarkerRecordSize {
		return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("model-state record has %d bytes, need %d", len(record), paul2013MarkerRecordSize)
	}
	phoneCount := int(record[paul2013MarkerPhoneCountOffset])
	if phoneCount == 0 || phoneCount > paul2013MarkerPhoneCapacity {
		return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("model-state record has unsupported phone count %d", phoneCount)
	}
	if len(phones) != phoneCount {
		return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("model-state record has %d phones but received %d decoded phone features", phoneCount, len(phones))
	}

	var result Paul2013ModelRecordPhoneGroups
	result.PhoneLabels = make([]byte, phoneCount)
	blockStart := 0
	groupOffset := 0
	nativeRowOffset := 0
	for phoneIndex := 0; phoneIndex < phoneCount; phoneIndex++ {
		if record[paul2013MarkerFollowingOffset+phoneIndex] == '0' && phoneIndex != phoneCount-1 {
			continue
		}
		blockEnd := phoneIndex + 1
		blockPhones := phones[blockStart:blockEnd]
		groups, err := BuildPaul2013PhoneGroups(blockPhones)
		if err != nil {
			return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("build model-state phone groups for marker span [%d,%d): %w", blockStart, blockEnd, err)
		}
		contexts, err := buildPaul2013PhoneDurationContexts(blockPhones, groups)
		if err != nil {
			return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("build model-state duration contexts for marker span [%d,%d): %w", blockStart, blockEnd, err)
		}
		for index := range contexts {
			result.PhoneLabels[blockStart+index] = byte(contexts[index].AuxiliaryByte)
		}
		if len(groups) == 0 {
			fallback := make([]byte, 0x1e)
			fallback[1] = 4
			fallback[2] = 1
			fallback[0x1c] = byte(blockStart)
			fallback[0x1d] = byte(blockEnd - blockStart)
			if nativeRowOffset == len(result.NativeGroupRows) {
				result.NativeGroupRows = append(result.NativeGroupRows, fallback)
			} else {
				result.NativeGroupRows[nativeRowOffset] = fallback
			}
		} else {
			for groupIndex, group := range groups {
				row := make([]byte, 0x1e)
				onsetCount := group.Nucleus - group.OnsetStart
				codaCount := group.End - group.Nucleus - 1
				row[0] = phones[blockStart+group.Nucleus].Stress
				if onsetCount > 0 {
					row[1] |= 1
				}
				if codaCount > 0 {
					row[1] |= 2
				}
				row[2] = group.PositionState
				row[0x1c] = byte(blockStart + group.Start)
				row[0x1d] = byte(group.End - group.Start)
				rowIndex := nativeRowOffset + groupIndex
				if rowIndex == len(result.NativeGroupRows) {
					result.NativeGroupRows = append(result.NativeGroupRows, row)
				} else {
					result.NativeGroupRows[rowIndex] = row
				}
			}
		}
		for index := range groups {
			groups[index].Start += blockStart
			groups[index].End += blockStart
			groups[index].OnsetStart += blockStart
			groups[index].Nucleus += blockStart
			result.Groups = append(result.Groups, groups[index])
		}
		result.NativeGroupCount += len(groups)
		nativeRowOffset += len(groups)
		for index := range contexts {
			contexts[index].GroupStart += blockStart
			contexts[index].GroupEnd += blockStart
			contexts[index].GroupIndex += groupOffset
			contexts[index].RowStart += blockStart
			result.DurationContexts = append(result.DurationContexts, contexts[index])
		}
		groupOffset += max(1, len(groups))
		blockStart = blockEnd
	}
	for index := range result.DurationContexts {
		result.DurationContexts[index].GroupCount = groupOffset
	}
	return result, nil
}

// BuildPaul2013ModelRecordPhoneGroupsFromSymbols decodes the record's
// +0x2e8 internal phone-symbol bytes through the native identity/stress tables,
// then applies the native +0x329 marker subdivision and phone-group writer.
// Stage 10 runtime captures cover all 69 observed CMU symbol codes and all 39
// identity ordinals. Stage 162 also compares native table windows directly
// with the local DLL and supports model aliases such as structural M.
func BuildPaul2013ModelRecordPhoneGroupsFromSymbols(
	record []byte,
) (Paul2013ModelRecordPhoneGroups, error) {
	if len(record) < paul2013MarkerRecordSize {
		return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("model-state record has %d bytes, need %d", len(record), paul2013MarkerRecordSize)
	}
	phoneCount := int(record[paul2013MarkerPhoneCountOffset])
	if phoneCount == 0 || phoneCount > paul2013MarkerPhoneCapacity {
		return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("model-state record has unsupported phone count %d", phoneCount)
	}
	phones, err := DecodePaul2013ModelPhoneSymbols(record[paul2013ModelRecordPhoneSymbolOffset : paul2013ModelRecordPhoneSymbolOffset+phoneCount])
	if err != nil {
		return Paul2013ModelRecordPhoneGroups{}, fmt.Errorf("decode model-state phone symbols: %w", err)
	}
	return BuildPaul2013ModelRecordPhoneGroups(record, phones)
}
