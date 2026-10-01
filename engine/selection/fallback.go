package selection

import (
	"fmt"
	"sort"
)

const (
	paul2013FallbackCandidateCapacity = 63
	paul2013FallbackRawCandidateLimit = 30
)

// Paul2013FallbackTreeTarget is one seven-byte tree target passed to
// FUN_10018770 while constructing a fallback position.
type Paul2013FallbackTreeTarget struct {
	Signature [7]byte
	QueryMode byte
}

// Paul2013FallbackQueryRow is one six-byte position row assembled by
// FUN_100242a0. The byte meanings remain opaque; fallback assigns bytes +3
// and +4 according to the observed row-construction path.
type Paul2013FallbackQueryRow struct {
	Bytes       [6]byte
	TreeTargets []Paul2013FallbackTreeTarget
}

// BuildPaul2013FallbackRows ports the row-copy and byte-update portion of
// FUN_100242a0 after whole-position lookup fails. contextDelta is the signed
// increment read from the model table at +0x6de; its production remains
// caller-owned. The native routine emits two rows, assigns byte +3 from that
// delta plus the row ordinal with byte-width wraparound, and writes row
// classes 1 and 2 to byte +4. Normally the original seven-byte target is
// queried in mode 1. For model
// row class 12, mode 1 uses a copy with signature byte +6 masked to bit 0x80,
// followed by the original target in mode 2.
func BuildPaul2013FallbackRows(
	base [6]byte,
	contextDelta int8,
	modelRowClass byte,
	target [7]byte,
) [2]Paul2013FallbackQueryRow {
	rows := [2]Paul2013FallbackQueryRow{}
	targets := []Paul2013FallbackTreeTarget{{Signature: target, QueryMode: 1}}
	if modelRowClass == 12 {
		maskedTarget := target
		maskedTarget[6] &= 0x80
		targets = []Paul2013FallbackTreeTarget{
			{Signature: maskedTarget, QueryMode: 1},
			{Signature: target, QueryMode: 2},
		}
	}
	for index := range rows {
		row := base
		row[3] = byte(int16(contextDelta) + int16(index))
		row[4] = byte(index + 1)
		rows[index] = Paul2013FallbackQueryRow{
			Bytes:       row,
			TreeTargets: append([]Paul2013FallbackTreeTarget(nil), targets...),
		}
	}
	return rows
}

// Paul2013FallbackCandidateList is the per-row candidate pool written by
// FUN_100242a0. RankedCount identifies the leading ranked entries; remaining
// IDs are raw generated candidates not present in that ranked list.
type Paul2013FallbackCandidateList struct {
	IDs         []uint32
	RankedCount int
}

// MergePaul2013FallbackCandidates ports FUN_100242a0's final row assembly.
// The ranked input retains FUN_10023c70's score order; raw IDs are
// canonicalized through FUN_10024010's ascending numeric sort and duplicate
// collapse. Ranked IDs are followed by ascending raw IDs absent from the
// ranked list. Native rows have 63 ID slots and class IDs are stored in the
// unsigned 16-bit range.
func MergePaul2013FallbackCandidates(
	rankedClassIDs, rawClassIDs []uint32,
) (Paul2013FallbackCandidateList, error) {
	ranked := append([]uint32(nil), rankedClassIDs...)
	seenRanked := make(map[uint32]struct{}, len(ranked))
	for index, classID := range ranked {
		if classID > uint32(^uint16(0)) {
			return Paul2013FallbackCandidateList{}, fmt.Errorf("ranked class ID at index %d (%d) exceeds native 16-bit range", index, classID)
		}
		if _, exists := seenRanked[classID]; exists {
			return Paul2013FallbackCandidateList{}, fmt.Errorf("ranked fallback list contains duplicate class ID %d", classID)
		}
		seenRanked[classID] = struct{}{}
	}
	if len(ranked) > paul2013FallbackCandidateCapacity {
		return Paul2013FallbackCandidateList{}, fmt.Errorf("ranked fallback list has %d IDs, exceeding native capacity %d", len(ranked), paul2013FallbackCandidateCapacity)
	}
	raw, err := canonicalPaul2013ClassIDs(rawClassIDs)
	if err != nil {
		return Paul2013FallbackCandidateList{}, fmt.Errorf("canonicalize raw fallback IDs: %w", err)
	}
	combined := append([]uint32(nil), ranked...)
	for _, rawID := range raw {
		if _, exists := seenRanked[rawID]; exists {
			continue
		}
		if len(combined) == paul2013FallbackCandidateCapacity {
			return Paul2013FallbackCandidateList{}, fmt.Errorf("fallback candidate union exceeds native capacity %d", paul2013FallbackCandidateCapacity)
		}
		combined = append(combined, rawID)
	}
	return Paul2013FallbackCandidateList{IDs: combined, RankedCount: len(ranked)}, nil
}

func canonicalPaul2013ClassIDs(classIDs []uint32) ([]uint32, error) {
	ordered := append([]uint32(nil), classIDs...)
	for index, classID := range ordered {
		if classID > uint32(^uint16(0)) {
			return nil, fmt.Errorf("class ID at index %d (%d) exceeds native 16-bit range", index, classID)
		}
	}
	sort.Slice(ordered, func(left, right int) bool { return ordered[left] < ordered[right] })
	unique := ordered[:0]
	for _, classID := range ordered {
		if len(unique) == 0 || unique[len(unique)-1] != classID {
			unique = append(unique, classID)
		}
	}
	return unique, nil
}
