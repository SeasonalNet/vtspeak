package text

import (
	"encoding/binary"
	"fmt"
)

// PopulatePaul2013RecordGroupDescriptorArena ports FUN_10012df0's header
// and descriptor writes. Descriptors stride by 0x10: start index +0x0c,
// count +0x0e, row span +0x10, record pointer +0x14, group-row pointer +0x18.
// The first group count is at arena +0. OutputShortOffset in the value view
// uses short units; native group-row pointers use twice that byte offset.
func PopulatePaul2013RecordGroupDescriptorArena(arena []byte, arenaAddress uint32) ([]byte, []Paul2013RecordGroupDescriptor, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(arena)
	if err != nil {
		return nil, nil, err
	}
	records := make([][]byte, len(inputs.InitialMarkers))
	for index := range records {
		start := paul2013TokenMarkerArenaRecords + index*paul2013TokenMarkerArenaStride
		records[index] = arena[start : start+paul2013TokenMarkerArenaStride]
	}
	groups, err := BuildPaul2013RecordGroupDescriptors(records)
	if err != nil {
		return nil, nil, err
	}
	result := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint16(result[0:2], uint16(len(groups)))
	for index, group := range groups {
		base := index * 0x10
		recordPointer := uint64(arenaAddress) + uint64(paul2013TokenMarkerArenaRecords+group.StartRecordIndex*paul2013TokenMarkerArenaStride)
		groupPointer := uint64(arenaAddress) + uint64(group.OutputShortOffset)*2
		if recordPointer > 0xffffffff || groupPointer > 0xffffffff {
			return nil, nil, fmt.Errorf("group %d native descriptor pointer overflows", index)
		}
		binary.LittleEndian.PutUint16(result[base+0x0c:base+0x0e], uint16(group.StartRecordIndex))
		binary.LittleEndian.PutUint16(result[base+0x0e:base+0x10], uint16(group.RecordCount))
		binary.LittleEndian.PutUint16(result[base+0x10:base+0x12], group.OpaqueRowSpan)
		binary.LittleEndian.PutUint32(result[base+0x14:base+0x18], uint32(recordPointer))
		binary.LittleEndian.PutUint32(result[base+0x18:base+0x1c], uint32(groupPointer))
	}
	return result, groups, nil
}
