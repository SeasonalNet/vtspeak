package text

import (
	"encoding/binary"
	"fmt"
)

const paul2013ModelPhoneGroupArenaBase = 0x17d4c

// Paul2013ModelPhoneGroupArena owns FUN_10012c70's arena writes and retains
// the decoded per-record groups for duration and pitch input construction.
type Paul2013ModelPhoneGroupArena struct {
	Arena   []byte
	Records []Paul2013ModelRecordPhoneGroups
	Count   int
}

// Paul2013PreparedTokenBoundaryArena is the pre-tree FUN_100130e0 handoff:
// phone-group arena writes, initial token markers, and record descriptors.
type Paul2013PreparedTokenBoundaryArena struct {
	PhoneGroups Paul2013ModelPhoneGroupArena
	Groups      []Paul2013RecordGroupDescriptor
}

// PreparePaul2013TokenBoundaryArena composes FUN_10012c70, initial marker
// writes, and FUN_10012df0 in caller order. The returned arena and descriptors
// are the explicit input to marker-tree evaluation, then Finish.
func PreparePaul2013TokenBoundaryArena(arena []byte, arenaAddress uint32) (Paul2013PreparedTokenBoundaryArena, error) {
	phoneGroups, err := PopulatePaul2013ModelPhoneGroupArena(arena, arenaAddress)
	if err != nil {
		return Paul2013PreparedTokenBoundaryArena{}, err
	}
	phoneGroups.Arena, err = InitializePaul2013TokenBoundaryArena(phoneGroups.Arena)
	if err != nil {
		return Paul2013PreparedTokenBoundaryArena{}, err
	}
	var groups []Paul2013RecordGroupDescriptor
	phoneGroups.Arena, groups, err = PopulatePaul2013RecordGroupDescriptorArena(phoneGroups.Arena, arenaAddress)
	if err != nil {
		return Paul2013PreparedTokenBoundaryArena{}, err
	}
	return Paul2013PreparedTokenBoundaryArena{PhoneGroups: phoneGroups, Groups: groups}, nil
}

// PopulatePaul2013ModelPhoneGroupArena writes group count +4, record group
// pointers +8, group counts +0x94, labels +0x36a, and the recovered fields
// of 0x1e-byte group rows at arena +0x17d4c. arenaAddress is the caller's
// explicit native address mapping. Unwritten row fields retain their bytes.
// A no-vowel fallback occupies the current tail without advancing the count;
// the next record may overwrite it, matching FUN_10012c70/FUN_10013c00.
func PopulatePaul2013ModelPhoneGroupArena(arena []byte, arenaAddress uint32) (Paul2013ModelPhoneGroupArena, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(arena)
	if err != nil {
		return Paul2013ModelPhoneGroupArena{}, err
	}
	if len(arena) < paul2013ModelPhoneGroupArenaBase {
		return Paul2013ModelPhoneGroupArena{}, fmt.Errorf("model arena lacks group row base +%#x", paul2013ModelPhoneGroupArenaBase)
	}
	result := Paul2013ModelPhoneGroupArena{Arena: append([]byte(nil), arena...), Records: make([]Paul2013ModelRecordPhoneGroups, len(inputs.InitialMarkers))}
	binary.LittleEndian.PutUint16(result.Arena[4:6], 0)
	for index := range result.Records {
		start := paul2013TokenMarkerArenaRecords + index*paul2013TokenMarkerArenaStride
		record := result.Arena[start : start+paul2013TokenMarkerArenaStride]
		groupStart := paul2013ModelPhoneGroupArenaBase + result.Count*0x1e
		address := uint64(arenaAddress) + uint64(groupStart)
		if address > 0xffffffff {
			return Paul2013ModelPhoneGroupArena{}, fmt.Errorf("record %d group pointer overflows native address", index)
		}
		binary.LittleEndian.PutUint32(record[8:12], uint32(address))
		record[0x94] = 0
		if record[0x95] == 0 {
			continue
		}
		groups, err := BuildPaul2013ModelRecordPhoneGroupsFromSymbols(record)
		if err != nil {
			return Paul2013ModelPhoneGroupArena{}, fmt.Errorf("record %d phone groups: %w", index, err)
		}
		if len(groups.NativeGroupRows)*0x1e > len(result.Arena)-groupStart {
			return Paul2013ModelPhoneGroupArena{}, fmt.Errorf("record %d group rows exceed arena", index)
		}
		for rowIndex, row := range groups.NativeGroupRows {
			output := result.Arena[groupStart+rowIndex*0x1e:]
			for _, field := range []int{0, 1, 2, 0x1c, 0x1d} {
				output[field] = row[field]
			}
		}
		copy(record[0x36a:], groups.PhoneLabels)
		record[0x94] = byte(groups.NativeGroupCount)
		result.Records[index] = groups
		result.Count += groups.NativeGroupCount
		binary.LittleEndian.PutUint16(result.Arena[4:6], uint16(result.Count))
	}
	return result, nil
}
