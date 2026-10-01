package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// Paul2013MarkerTreeArenaResult owns the arena after FUN_10012f00 and
// retains each group's reverse-scan inputs and scalar outputs.
type Paul2013MarkerTreeArenaResult struct {
	Arena  []byte
	Groups []Paul2013RecordGroupDescriptor
	Scans  []Paul2013MarkerTreeTokenResult
}

// RunPaul2013MarkerTreeArena derives record groups from initialized terminal
// markers, then visits each group's records from last to second. It resolves
// +0x2e4 string pointers only on the eligible branch and writes the preceding
// record's +0x3bd terminal when the supplied native scalar tree exceeds 500.
// OpaqueRowSpan is this group's total +0x94 sum, the short read at native
// arena +(group+1)*0x10; it is not the following group's cumulative value.
func RunPaul2013MarkerTreeArena(arena []byte, tree *tree3.Tree, resolvePointer Paul2013CStringPointerResolver) (Paul2013MarkerTreeArenaResult, error) {
	inputs, err := ReadPaul2013TokenBoundaryArenaInputs(arena)
	if err != nil {
		return Paul2013MarkerTreeArenaResult{}, err
	}
	result := Paul2013MarkerTreeArenaResult{Arena: append([]byte(nil), arena...)}
	records := make([][]byte, len(inputs.InitialMarkers))
	for index := range records {
		start := paul2013TokenMarkerArenaRecords + index*paul2013TokenMarkerArenaStride
		records[index] = result.Arena[start : start+paul2013TokenMarkerArenaStride]
	}
	result.Groups, err = BuildPaul2013RecordGroupDescriptors(records)
	if err != nil {
		return Paul2013MarkerTreeArenaResult{}, err
	}
	result.Scans = make([]Paul2013MarkerTreeTokenResult, len(result.Groups))
	for groupIndex, group := range result.Groups {
		scan := Paul2013MarkerTreeTokenResult{EndState: Paul2013MarkerTreeInputState{TokenPhoneCount: int16(group.RecordCount), NextTokenCumulative: int16(group.OpaqueRowSpan)}}
		for index := group.StartRecordIndex + group.RecordCount - 1; index > group.StartRecordIndex; index-- {
			previous, current := records[index-1], records[index]
			state := scan.EndState
			// The branch test precedes native string dereferences. Rejected
			// visits still increment all four counters.
			probe, err := BuildPaul2013MarkerTreeInput(previous, current, state)
			if err != nil {
				return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("marker group %d record %d: %w", groupIndex, index, err)
			}
			if probe.Eligible {
				if resolvePointer == nil {
					return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("marker group %d record %d needs string pointer resolver", groupIndex, index)
				}
				currentString, err := resolvePointer(binary.LittleEndian.Uint32(current[0x2e4:0x2e8]))
				if err != nil {
					return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("current record %d string: %w", index, err)
				}
				if len(currentString) == 0 {
					return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("current record %d string has no first byte", index)
				}
				previousString, err := resolvePointer(binary.LittleEndian.Uint32(previous[0x2e4:0x2e8]))
				if err != nil {
					return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("previous record %d string: %w", index-1, err)
				}
				selected := int(previous[0x95]) - 1
				if selected < 0 || selected >= len(previousString) {
					return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("previous record %d string lacks selected index %d", index-1, selected)
				}
				state.CurrentStringFirstByte, state.PreviousSelectedByte = currentString[0], previousString[selected]
			}
			input, value, leaf, decision, err := EvaluatePaul2013MarkerTreeInput(tree, previous, current, state)
			if err != nil {
				return Paul2013MarkerTreeArenaResult{}, fmt.Errorf("marker group %d record %d: %w", groupIndex, index, err)
			}
			copy(previous, decision.PhoneRow)
			scan.Inputs = append(scan.Inputs, input)
			scan.TreeValues = append(scan.TreeValues, value)
			scan.Leaves = append(scan.Leaves, leaf)
			scan.Decisions = append(scan.Decisions, decision)
			scan.EndState.PhoneOrdinal = input.PhoneOrdinal
			scan.EndState.PhonesSinceLastSplit = input.PhonesSinceLastSplit
			scan.EndState.CumulativePhoneValue = input.CumulativePhoneValue
			scan.EndState.ValueSinceLastSplit = input.ValueSinceLastSplit
			if decision.InsertedMarker {
				scan.EndState.PhonesSinceLastSplit, scan.EndState.ValueSinceLastSplit = 0, 0
			}
		}
		result.Scans[groupIndex] = scan
	}
	return result, nil
}
