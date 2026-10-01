package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	paul2013FollowupComponentStride = 0x140
	paul2013FollowupComponentHeader = 0x0c
	paul2013FollowupSurfaceOffset   = 0x46
	paul2013SourceRowFollowupByte   = 0x1e
	paul2013SourceRowFollowupValue  = 0x1f
)

// AppendPaul2013MatchedComponentSourceRows ports FUN_10034110 after its
// caller has accepted a candidate with FUN_10034180. componentArena contains
// the counted 0x140-byte component rows; sourceArena is the FUN_10044fd0
// destination. The native append arguments remain raw field values.
func AppendPaul2013MatchedComponentSourceRows(
	sourceArena []byte,
	componentArena []byte,
) ([]byte, bool, error) {
	if len(componentArena) < paul2013FollowupComponentHeader {
		return nil, false, fmt.Errorf("matched-component arena has %d bytes, need a %#x-byte header", len(componentArena), paul2013FollowupComponentHeader)
	}
	componentCount := int32(binary.LittleEndian.Uint32(componentArena[:4]))
	if componentCount < 0 {
		return nil, false, fmt.Errorf("matched-component count is negative: %d", componentCount)
	}
	if int64(componentCount) > int64((len(componentArena)-paul2013FollowupComponentHeader)/paul2013FollowupComponentStride) {
		return nil, false, fmt.Errorf("matched-component arena has %d rows of %d bytes, count is %d", (len(componentArena)-paul2013FollowupComponentHeader)/paul2013FollowupComponentStride, paul2013FollowupComponentStride, componentCount)
	}
	if len(sourceArena) < paul2013SourceRowsOffset {
		return nil, false, fmt.Errorf("source-row arena has %d bytes, need at least %#x", len(sourceArena), paul2013SourceRowsOffset)
	}

	working := append([]byte(nil), sourceArena...)
	for index := int32(0); index < componentCount; index++ {
		componentStart := paul2013FollowupComponentHeader + int(index)*paul2013FollowupComponentStride
		component := componentArena[componentStart : componentStart+paul2013FollowupComponentStride]
		surfaceRegion := component[paul2013FollowupSurfaceOffset:]
		surfaceEnd := bytes.IndexByte(surfaceRegion, 0)
		if surfaceEnd < 0 {
			return nil, false, fmt.Errorf("matched component %d surface at +0x%x is not NUL-terminated", index, paul2013FollowupSurfaceOffset)
		}
		input := Paul2013ModelSourceRowInput{
			First:  binary.LittleEndian.Uint32(component[0x00:0x04]),
			Second: binary.LittleEndian.Uint32(component[0x04:0x08]),
			Type:   0x41,
			Class:  0x44,
			Flag:   0x12,
			Text:   surfaceRegion[:surfaceEnd],
		}
		var appended bool
		var err error
		working, appended, err = AppendPaul2013ModelSourceRowsSplit(working, input)
		if err != nil {
			return nil, false, fmt.Errorf("append matched component %d: %w", index, err)
		}
		if !appended {
			return working, false, nil
		}
		rowCount := int(binary.LittleEndian.Uint16(working[paul2013SourceRowCountOffset:]))
		if rowCount == 0 {
			return nil, false, fmt.Errorf("matched component %d follow-up write has no source row", index)
		}
		lastRowStart := paul2013SourceRowsOffset + (rowCount-1)*paul2013SourceRowStride
		if lastRowStart+paul2013SourceRowFollowupByte >= len(working) {
			return nil, false, fmt.Errorf("source-row arena has %d bytes, cannot write matched-component follow-up at row %d", len(working), rowCount-1)
		}
		working[lastRowStart+paul2013SourceRowFollowupByte] = paul2013SourceRowFollowupValue
	}
	return working, true, nil
}
