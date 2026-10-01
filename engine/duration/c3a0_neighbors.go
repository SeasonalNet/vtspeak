package duration

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013C3A0ContextSurfaceOffset = 0x05
	paul2013C3A0ContextSurfaceEnd    = 0x23
)

// Paul2013FUN100086C0NeighborCandidate describes the nearest row in one
// direction whose first or second surface byte carries a native class bit.
// CharacterGate is the low-short result of FUN_1000ffd0 for that surface.
type Paul2013FUN100086C0NeighborCandidate struct {
	RowIndex      int
	Surface       []byte
	CharacterGate bool
}

// Paul2013FUN100086C0Neighbors contains the independently scanned previous
// and following candidates. A nil candidate means that the native scan found
// no qualifying row in that direction.
type Paul2013FUN100086C0Neighbors struct {
	Previous  *Paul2013FUN100086C0NeighborCandidate
	Following *Paul2013FUN100086C0NeighborCandidate
}

// ScanPaul2013FUN100086C0Neighbors ports FUN_100086c0's outward row scans for
// an already-reached neighbor-scan branch. It scans to the nearest row whose
// first or second surface byte has either 0xc0 class bit, then evaluates that
// row with FUN_1000ffd0. The later runtime sorted-table lookup and status-
// specific decisions are applied by CompletePaul2013FUN100086C0NeighborDecision.
func ScanPaul2013FUN100086C0Neighbors(
	ctx context.Context,
	model []byte,
	contextRowIndex int,
) (Paul2013FUN100086C0Neighbors, error) {
	if ctx == nil {
		return Paul2013FUN100086C0Neighbors{}, errors.New("FUN_100086c0 neighbor scan has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013FUN100086C0Neighbors{}, err
	}
	if len(model) < paul2013C3A0ContextCountOffset+2 {
		return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013C3A0ContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013C3A0ContextCountOffset:])))
	if rowCount < 0 || contextRowIndex < 0 || contextRowIndex >= rowCount {
		return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("context row index %d is outside declared row count %d", contextRowIndex, rowCount)
	}
	if rowCount > (len(model)-paul2013C3A0ContextRowBase)/paul2013C3A0ContextRowStride {
		return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("model declares %d context rows that do not fit its %d bytes", rowCount, len(model))
	}

	var result Paul2013FUN100086C0Neighbors
	for rowIndex := contextRowIndex - 1; rowIndex >= 0; rowIndex-- {
		candidate, qualifies, err := paul2013C3A0NeighborCandidate(model, rowIndex)
		if err != nil {
			return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("inspect previous context row %d: %w", rowIndex, err)
		}
		if !qualifies {
			continue
		}
		candidate.CharacterGate, err = Paul2013CompoundCharacterGate(candidate.Surface)
		if err != nil {
			return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("run FUN_1000ffd0 for previous context row %d: %w", rowIndex, err)
		}
		result.Previous = &candidate
		break
	}
	for rowIndex := contextRowIndex + 1; rowIndex < rowCount; rowIndex++ {
		candidate, qualifies, err := paul2013C3A0NeighborCandidate(model, rowIndex)
		if err != nil {
			return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("inspect following context row %d: %w", rowIndex, err)
		}
		if !qualifies {
			continue
		}
		candidate.CharacterGate, err = Paul2013CompoundCharacterGate(candidate.Surface)
		if err != nil {
			return Paul2013FUN100086C0Neighbors{}, fmt.Errorf("run FUN_1000ffd0 for following context row %d: %w", rowIndex, err)
		}
		result.Following = &candidate
		break
	}
	if err := ctx.Err(); err != nil {
		return Paul2013FUN100086C0Neighbors{}, err
	}
	return result, nil
}

func paul2013C3A0NeighborCandidate(
	model []byte,
	rowIndex int,
) (Paul2013FUN100086C0NeighborCandidate, bool, error) {
	rowStart := paul2013C3A0ContextRowBase + rowIndex*paul2013C3A0ContextRowStride
	surfaceArea := model[rowStart+paul2013C3A0ContextSurfaceOffset : rowStart+paul2013C3A0ContextSurfaceEnd]
	end := bytes.IndexByte(surfaceArea, 0)
	if end < 0 {
		return Paul2013FUN100086C0NeighborCandidate{}, false, errors.New("surface has no NUL terminator before the phone string")
	}
	surface := append([]byte(nil), surfaceArea[:end]...)
	attributes := text.Paul2013UnsignedCharacterAttributeTable()
	qualifies := false
	classWindow := surface
	if len(classWindow) > 2 {
		classWindow = classWindow[:2]
	}
	for _, value := range classWindow {
		if value >= 0x80 {
			return Paul2013FUN100086C0NeighborCandidate{}, false, fmt.Errorf("surface class byte 0x%02x is outside the supported ASCII table view", value)
		}
		if attributes[value]&0xc0 != 0 {
			qualifies = true
			break
		}
	}
	return Paul2013FUN100086C0NeighborCandidate{RowIndex: rowIndex, Surface: surface}, qualifies, nil
}
