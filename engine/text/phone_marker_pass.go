package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/internal/paul2013tables"
)

const paul2013MarkerRecordSize = 0x3c0

const paul2013MarkerPhoneCapacity = 65

const paul2013MarkerSpecialKeyPointerOffset = 0x3b8

const (
	paul2013MarkerPhoneCountOffset   = 0x95
	paul2013MarkerPhoneCodeOffset    = 0x2e8
	paul2013MarkerPreviousByteOffset = 0x328
	paul2013MarkerFollowingOffset    = 0x329
	paul2013MarkerTerminalOffset     = 0x3bd
)

// paul2013SpecialTerminalKeys is the ten-entry pointer table rooted at
// DAT_1007beec and searched by FUN_100560a0. Their linguistic roles remain
// opaque.
var paul2013SpecialTerminalKeys = [...]string{
	"ya", "you", "you'd", "you'll", "you're", "you've", "your", "yours", "yourself", "yourselves",
}

// Paul2013SpecialTerminalKeyMatch ports the membership result consumed from
// FUN_100560a0. The caller supplies the bytes reached through the native
// pointer field and the engine's character map for FUN_1001c2c0 comparisons.
func Paul2013SpecialTerminalKeyMatch(key []byte, characterMap [256]byte) bool {
	key = cString(key)
	if len(key) == 0 {
		return false
	}
	for _, candidate := range paul2013SpecialTerminalKeys {
		if Paul2013MappedCStringEqual(key, []byte(candidate), characterMap) {
			return true
		}
	}
	return false
}

// Paul2013PreviousPhoneCategory ports FUN_10017030's category read for an
// already validated phone index. The category bytes remain opaque.
func paul2013PreviousPhoneCategory(record []byte, phoneIndex int) int16 {
	if phoneIndex < 0 {
		return -1
	}
	phone := record[paul2013MarkerPhoneCodeOffset+phoneIndex]
	if paul2013tables.ByteTable1007B9E0(phone) == 1 {
		return int16(paul2013tables.ByteTable1007BAA8(phone))
	}
	if phone == '6' && phoneIndex > 0 {
		previous := record[paul2013MarkerPhoneCodeOffset+phoneIndex-1]
		if paul2013tables.ByteTable1007B9E0(previous) == 1 {
			return int16(paul2013tables.ByteTable1007BAA8(previous))
		}
	}
	return -1
}

// Paul2013PhoneBoundaryCategory ports FUN_10017090's next-phone category
// lookup, including its special 'C' skip to the following phone.
func paul2013PhoneBoundaryCategory(record []byte, phoneIndex int) int16 {
	phoneCount := int(record[paul2013MarkerPhoneCountOffset])
	if phoneIndex < 0 || phoneIndex >= phoneCount {
		return -1
	}
	phone := record[paul2013MarkerPhoneCodeOffset+phoneIndex]
	if paul2013tables.ByteTable1007B9E0(phone) == 1 {
		return int16(int8(paul2013tables.ByteTable1007BB70(phone)))
	}
	if phone == 'C' && phoneIndex < phoneCount-1 {
		next := record[paul2013MarkerPhoneCodeOffset+phoneIndex+1]
		if paul2013tables.ByteTable1007B9E0(next) == 1 {
			return int16(int8(paul2013tables.ByteTable1007BB70(next)))
		}
	}
	return -1
}

// ApplyPaul2013PhoneMarkerPrepass ports the record-byte transformations in
// FUN_10017100 over caller-ordered 0x3c0-byte records. Its outputs are copied
// before mutation. specialKeys corresponds to adjacent record boundaries and
// carries bytes reached through the native pointer fields; its comparison map
// is used only if a ']' terminal row reaches the DAT_1007beec lookup.
func ApplyPaul2013PhoneMarkerPrepass(
	inputRecords [][]byte,
	specialKeys [][]byte,
	characterMap [256]byte,
) ([][]byte, error) {
	if len(specialKeys) != 0 && len(specialKeys) != max(0, len(inputRecords)-1) {
		return nil, fmt.Errorf("received %d special-key values for %d adjacent boundaries", len(specialKeys), max(0, len(inputRecords)-1))
	}
	var resolve func(int, []byte) ([]byte, error)
	if len(specialKeys) != 0 {
		resolve = func(index int, _ []byte) ([]byte, error) { return specialKeys[index], nil }
	}
	return applyPaul2013PhoneMarkerPrepass(inputRecords, characterMap, resolve)
}

// Paul2013CStringPointerResolver reads a NUL-terminated string at a 32-bit
// address stored in a native record. The callback owns address-space bounds
// and lifetime; the engine never dereferences process pointers itself.
type Paul2013CStringPointerResolver func(address uint32) ([]byte, error)

// ApplyPaul2013PhoneMarkerPrepassWithPointerResolver resolves the native
// pointer field at record +0x3b8 only if the terminal marker pass reaches its
// DAT_1007beec search. The resolver must map captured or live native addresses
// into readable C-string bytes.
func ApplyPaul2013PhoneMarkerPrepassWithPointerResolver(
	inputRecords [][]byte,
	resolvePointer Paul2013CStringPointerResolver,
	characterMap [256]byte,
) ([][]byte, error) {
	var resolve func(int, []byte) ([]byte, error)
	if resolvePointer != nil {
		resolve = func(_ int, record []byte) ([]byte, error) {
			address := binary.LittleEndian.Uint32(record[paul2013MarkerSpecialKeyPointerOffset:])
			if address == 0 {
				return nil, nil
			}
			return resolvePointer(address)
		}
	}
	return applyPaul2013PhoneMarkerPrepass(inputRecords, characterMap, resolve)
}

func applyPaul2013PhoneMarkerPrepass(
	inputRecords [][]byte,
	characterMap [256]byte,
	resolveSpecialKey func(int, []byte) ([]byte, error),
) ([][]byte, error) {
	if len(inputRecords) > paul2013PackedContextRecordCapacity {
		return nil, fmt.Errorf("received %d marker records, exceeding native capacity %d", len(inputRecords), paul2013PackedContextRecordCapacity)
	}
	records := make([][]byte, len(inputRecords))
	for recordIndex, record := range inputRecords {
		if len(record) < paul2013MarkerRecordSize {
			return nil, fmt.Errorf("marker record %d has %d bytes, need %d", recordIndex, len(record), paul2013MarkerRecordSize)
		}
		phoneCount := int(record[paul2013MarkerPhoneCountOffset])
		if phoneCount == 0 || phoneCount > paul2013MarkerPhoneCapacity {
			return nil, fmt.Errorf("marker record %d has unsupported phone count %d", recordIndex, phoneCount)
		}
		records[recordIndex] = append([]byte(nil), record...)
	}

	for _, record := range records {
		phoneCount := int(record[paul2013MarkerPhoneCountOffset])
		phones := record[paul2013MarkerPhoneCodeOffset:]
		previousMarkers := record[paul2013MarkerPreviousByteOffset:]
		followingMarkers := record[paul2013MarkerFollowingOffset:]
		for phoneIndex := 0; phoneIndex < phoneCount-1; phoneIndex++ {
			phone := phones[phoneIndex]
			mapped := paul2013tables.ByteTable1007BE98(phone)
			if followingMarkers[phoneIndex] == '0' && mapped != 0 && phones[phoneIndex+1] == '6' {
				phones[phoneIndex] = mapped
			}
		}
		for phoneIndex := 1; phoneIndex < phoneCount-1; phoneIndex++ {
			if previousMarkers[phoneIndex] != '0' {
				continue
			}
			phone := phones[phoneIndex]
			previousPhone := phones[phoneIndex-1]
			followingPhone := phones[phoneIndex+1]
			previousCategory := paul2013PreviousPhoneCategory(record, phoneIndex-1)
			followingCategory := paul2013PhoneBoundaryCategory(record, phoneIndex)
			if followingMarkers[phoneIndex] == '0' {
				switch phone {
				case '9':
					if (previousPhone == '7' || previousPhone == '8') && followingCategory >= 0 {
						phones[phoneIndex] = 'H'
					} else if paul2013tables.ByteTable1007BE3C(previousPhone) == 1 && followingCategory == 0 {
						phones[phoneIndex] = 'H'
					} else if previousCategory >= 1 && phoneIndex < phoneCount-2 &&
						followingMarkers[phoneIndex+1] == '0' && followingPhone == '7' && phones[phoneIndex+2] == '-' {
						phones[phoneIndex] = 'K'
					} else if previousCategory >= 0 {
						if followingCategory == 0 {
							phones[phoneIndex] = 'M'
						} else if followingPhone == '+' || followingPhone == ',' || followingPhone == '-' {
							phones[phoneIndex] = 'K'
						}
					}
				case '*', '5':
					if followingCategory >= 0 && (previousPhone == '7' || followingCategory == 0) {
						if phone == '*' {
							phones[phoneIndex] = 'F'
						} else {
							phones[phoneIndex] = 'G'
						}
					}
				case 0x15:
					if previousCategory >= 0 && followingCategory == 0 {
						phones[phoneIndex] = 'N'
					}
				}
			} else if previousCategory >= 0 && followingCategory >= 0 {
				if followingPhone != 'C' {
					switch phone {
					case '9':
						phones[phoneIndex] = 'M'
					case 0x15:
						phones[phoneIndex] = 'N'
					}
				}
				if phone == '*' {
					phones[phoneIndex] = 'F'
				} else if phone == '5' {
					phones[phoneIndex] = 'G'
				}
			}
			if phone == '9' && paul2013tables.ByteTable1007BE3C(previousPhone) == 1 && followingCategory >= 0 {
				phones[phoneIndex] = 'H'
			}
		}
		for phoneIndex := 0; phoneIndex < phoneCount-1; phoneIndex++ {
			if followingMarkers[phoneIndex] == '0' && phones[phoneIndex] == '9' {
				category := paul2013PhoneBoundaryCategory(record, phoneIndex)
				if category == 1 || category == 2 {
					phones[phoneIndex] = 'L'
				}
			}
		}
	}

	for recordIndex := 0; recordIndex+1 < len(records); recordIndex++ {
		current, next := records[recordIndex], records[recordIndex+1]
		if current[paul2013MarkerTerminalOffset] != ']' {
			continue
		}
		phoneCount := int(current[paul2013MarkerPhoneCountOffset])
		lastPhoneIndex := phoneCount - 1
		previousPhoneIndex := phoneCount - 2
		followingCategory := paul2013PhoneBoundaryCategory(next, 0)
		previousCategory := paul2013PreviousPhoneCategory(current, previousPhoneIndex)
		phones := current[paul2013MarkerPhoneCodeOffset:]
		if followingCategory >= 0 {
			if previousCategory >= 0 && next[paul2013MarkerPhoneCodeOffset] != 'C' {
				switch phones[lastPhoneIndex] {
				case '9':
					phones[lastPhoneIndex] = 'M'
				case 0x15:
					phones[lastPhoneIndex] = 'N'
				}
			}
			if previousCategory >= 0 {
				if phones[lastPhoneIndex] == '*' {
					phones[lastPhoneIndex] = 'F'
				} else if phones[lastPhoneIndex] == '5' {
					phones[lastPhoneIndex] = 'G'
				}
			}
			if previousPhoneIndex > 0 && paul2013tables.ByteTable1007BE3C(phones[previousPhoneIndex-1]) == 1 {
				switch phones[lastPhoneIndex] {
				case '9':
					phones[lastPhoneIndex] = 'H'
				case '*':
					phones[lastPhoneIndex] = 'F'
				case '5':
					phones[lastPhoneIndex] = 'G'
				}
			}
			if phones[lastPhoneIndex] == 0x15 || phones[lastPhoneIndex] == '9' {
				if resolveSpecialKey == nil {
					return nil, fmt.Errorf("record boundary %d requires the pointer-backed DAT_1007beec lookup result", recordIndex)
				}
				key, err := resolveSpecialKey(recordIndex, current)
				if err != nil {
					return nil, fmt.Errorf("resolve special terminal key for record boundary %d: %w", recordIndex, err)
				}
				if Paul2013SpecialTerminalKeyMatch(key, characterMap) {
					if phones[lastPhoneIndex] == 0x15 {
						phones[lastPhoneIndex] = 'O'
					} else {
						phones[lastPhoneIndex] = 'P'
					}
				}
			}
		}
	}
	return records, nil
}

// Paul2013PreparedContextGroup retains the preprocessed source arena and the
// packed-context result for one caller-selected ordered record group.
type Paul2013PreparedContextGroup struct {
	PreprocessedRecords [][]byte
	Group               Paul2013PackedContextPhoneGroupResult
}

// PreparePaul2013PhoneMarkerContextGroup composes the all-record marker
// prepass with packed-context construction for explicit native record indexes.
// Callers that need FUN_10012df0's terminal-based groups can use
// PreparePaul2013PhoneMarkerRecordGroups instead.
func PreparePaul2013PhoneMarkerContextGroup(
	records [][]byte,
	specialKeys [][]byte,
	characterMap [256]byte,
	groupRecordIndexes []int,
) (Paul2013PreparedContextGroup, error) {
	if len(groupRecordIndexes) == 0 {
		return Paul2013PreparedContextGroup{}, fmt.Errorf("context group has no selected records")
	}
	preprocessed, err := ApplyPaul2013PhoneMarkerPrepass(records, specialKeys, characterMap)
	if err != nil {
		return Paul2013PreparedContextGroup{}, fmt.Errorf("preprocess phone markers: %w", err)
	}
	selected := make([][]byte, len(groupRecordIndexes))
	for groupIndex, recordIndex := range groupRecordIndexes {
		if recordIndex < 0 || recordIndex >= len(preprocessed) {
			return Paul2013PreparedContextGroup{}, fmt.Errorf("group record index %d at position %d outside %d records", recordIndex, groupIndex, len(preprocessed))
		}
		selected[groupIndex] = preprocessed[recordIndex]
	}
	group, err := PreparePaul2013PackedContextPhoneGroup(selected)
	if err != nil {
		return Paul2013PreparedContextGroup{}, fmt.Errorf("prepare packed phone context: %w", err)
	}
	return Paul2013PreparedContextGroup{
		PreprocessedRecords: preprocessed,
		Group:               group,
	}, nil
}
