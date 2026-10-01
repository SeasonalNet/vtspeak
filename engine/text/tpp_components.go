package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// TPPComponentPattern retains the numeric family, component count, and raw
// binary values decoded from an A-E TPP atom. The bit meanings remain opaque.
// AX is represented separately because its X payload is not a component bit.
type TPPComponentPattern struct {
	Family             byte
	ComponentCount     uint8
	ComponentValues    []byte
	IsAXRecognitionTag bool
}

// Paul2013ProperNameTPPComponent contains the source row fields consumed by
// FUN_10034180 while it assembles the candidate TPP surface. Length is the
// stored dword at component-row +0x14; Marker is the byte at +0x1a. The raw
// +0x18 dword is passed as FUN_1005a350's input pointer; zero bypasses that
// scanner and its later context predicates.
type Paul2013ProperNameTPPComponent struct {
	Surface             []byte
	Length              int32
	Marker              byte
	ContextScanInput    uint32
	HasContextScanInput bool
}

// Paul2013ProperNameTPPComponentsFromArena projects the counted component
// rows read by FUN_10034180. The arena has a 12-byte header, then 0x140-byte
// rows whose stored length, marker, and C-string surface are at +0x08,
// +0x0e, and +0x3a respectively. It also projects the final-row +0x18
// scanner-input dword.
func Paul2013ProperNameTPPComponentsFromArena(arena []byte) ([]Paul2013ProperNameTPPComponent, error) {
	if len(arena) < paul2013FollowupComponentHeader {
		return nil, fmt.Errorf("proper-name component arena has %d bytes, need %#x-byte header", len(arena), paul2013FollowupComponentHeader)
	}
	count := int32(binary.LittleEndian.Uint32(arena[:4]))
	if count < 0 {
		return nil, fmt.Errorf("proper-name component count is negative: %d", count)
	}
	available := (len(arena) - paul2013FollowupComponentHeader) / paul2013FollowupComponentStride
	if int64(count) > int64(available) {
		return nil, fmt.Errorf("proper-name component arena has %d rows of %d bytes, count is %d", available, paul2013FollowupComponentStride, count)
	}
	components := make([]Paul2013ProperNameTPPComponent, count)
	for index := range components {
		start := paul2013FollowupComponentHeader + index*paul2013FollowupComponentStride
		row := arena[start : start+paul2013FollowupComponentStride]
		surfaceRegion := row[0x3a:]
		end := bytes.IndexByte(surfaceRegion, 0)
		if end < 0 {
			return nil, fmt.Errorf("proper-name component %d surface at row +0x3a is not NUL-terminated", index)
		}
		components[index] = Paul2013ProperNameTPPComponent{
			Surface:             append([]byte(nil), surfaceRegion[:end]...),
			Length:              int32(binary.LittleEndian.Uint32(row[0x08:0x0c])),
			Marker:              row[0x0e],
			ContextScanInput:    binary.LittleEndian.Uint32(row[0x18:0x1c]),
			HasContextScanInput: true,
		}
	}
	return components, nil
}

// BuildPaul2013ProperNameTPPSurface ports FUN_10034180's one-through-four
// component surface assembly. It joins C-string surfaces with '-', and
// applies the native sum(length+1) <= 0x1d gate before dictionary lookup.
// Counts outside 1..4 and over-limit keys are ineligible without an error.
func BuildPaul2013ProperNameTPPSurface(
	components []Paul2013ProperNameTPPComponent,
) ([]byte, bool, error) {
	if len(components) < 1 || len(components) > 4 {
		return nil, false, nil
	}
	var nativeLength int32
	surface := make([]byte, 0, 0x1c)
	for index, component := range components {
		if component.Length < 0 {
			return nil, false, fmt.Errorf("proper-name component %d has negative stored length %d", index, component.Length)
		}
		text := component.Surface
		if nul := bytes.IndexByte(text, 0); nul >= 0 {
			text = text[:nul]
		}
		if int32(len(text)) != component.Length {
			return nil, false, fmt.Errorf("proper-name component %d stores length %d for %d C-string bytes", index, component.Length, len(text))
		}
		if component.Length >= 0x1d || component.Length+1 > 0x1d-nativeLength {
			return nil, false, nil
		}
		nativeLength += component.Length + 1
		if index > 0 {
			surface = append(surface, '-')
		}
		surface = append(surface, text...)
	}
	return surface, true, nil
}

// Paul2013ProperNameTPPSelector ports the family byte selected by
// FUN_10034180 for short component sequences. The caller selects A-D for one
// through four components; its x86 count gate skips selector construction at
// five or more components, even though E records exist in the dictionary.
func Paul2013ProperNameTPPSelector(componentCount int) (byte, bool) {
	if componentCount < 1 || componentCount > 4 {
		return 0, false
	}
	return byte('A' + componentCount - 1), true
}

// AcceptPaul2013ProperNameTPPComponentPattern ports the post-lookup gate in
// FUN_10034180. A miss is rejected. For one component the returned A payload
// is ignored and only marker d/A rejects the candidate. For multiple
// components, a literal zero bit rejects only when its corresponding marker
// is d/A. Marker production remains an upstream input; the marker bytes are
// intentionally not assigned linguistic names.
func AcceptPaul2013ProperNameTPPComponentPattern(
	pattern TPPComponentPattern,
	componentMarkers []byte,
	lookupFound bool,
) (bool, error) {
	if !lookupFound {
		return false, nil
	}
	count := int(pattern.ComponentCount)
	if count < 1 || count > 4 {
		return false, fmt.Errorf("proper-name TPP pattern has unsupported component count %d", count)
	}
	if len(componentMarkers) != count {
		return false, fmt.Errorf("proper-name TPP pattern has %d components but %d markers were supplied", count, len(componentMarkers))
	}
	if count == 1 {
		return !paul2013ProperNameMarkerRejects(componentMarkers[0]), nil
	}
	if len(pattern.ComponentValues) != count {
		return false, fmt.Errorf("proper-name TPP pattern has %d component values, want %d", len(pattern.ComponentValues), count)
	}
	for index, value := range pattern.ComponentValues {
		if value == 0 && paul2013ProperNameMarkerRejects(componentMarkers[index]) {
			return false, nil
		}
	}
	return true, nil
}

func paul2013ProperNameMarkerRejects(marker byte) bool {
	return marker == 'd' || marker == 'A'
}

// ComponentPattern decodes the runtime-observed A-E component payload shape.
// A0/A1 carry one bit for one component; B-E carry a count digit followed by
// one bit for each hyphen-separated component. The AX spelling is retained as
// a single-component marker with no inferred bit value.
func (atom TPPAtom) ComponentPattern() (TPPComponentPattern, bool, error) {
	if atom.Tag < 'A' || atom.Tag > 'E' {
		return TPPComponentPattern{}, false, nil
	}
	if atom.Tag == 'A' && string(atom.Suffix) == "X" {
		return TPPComponentPattern{
			Family: 'A', ComponentCount: 1, IsAXRecognitionTag: true,
		}, true, nil
	}
	if len(atom.Suffix) == 0 {
		return TPPComponentPattern{}, false, fmt.Errorf("TPP %c component suffix is empty", atom.Tag)
	}

	expectedCount := int(atom.Tag - 'A' + 1)
	count := 1
	bitStart := 0
	if atom.Tag != 'A' {
		if atom.Suffix[0] < '0' || atom.Suffix[0] > '9' {
			return TPPComponentPattern{}, false, fmt.Errorf("TPP %c component count %q is not decimal", atom.Tag, atom.Suffix[:1])
		}
		count = int(atom.Suffix[0] - '0')
		bitStart = 1
		if count != expectedCount {
			return TPPComponentPattern{}, false, fmt.Errorf("TPP %c component count %d, want %d", atom.Tag, count, expectedCount)
		}
	}
	if len(atom.Suffix)-bitStart != count {
		return TPPComponentPattern{}, false, fmt.Errorf("TPP %c payload has %d component values, want %d", atom.Tag, len(atom.Suffix)-bitStart, count)
	}

	values := make([]byte, count)
	for index, value := range atom.Suffix[bitStart:] {
		if value != '0' && value != '1' {
			return TPPComponentPattern{}, false, fmt.Errorf("TPP %c component value %d is %q, want 0 or 1", atom.Tag, index, value)
		}
		values[index] = value - '0'
	}
	return TPPComponentPattern{
		Family: atom.Tag, ComponentCount: uint8(count), ComponentValues: values,
	}, true, nil
}

// LookupTPPComponentPattern returns the first A-E atom from an exact TPP key.
// Other atoms in a two-atom payload remain available from LookupTypedText.
func (dictionary *TPPDictionary) LookupTPPComponentPattern(
	surface []byte,
	tables EmbeddedKeyTables,
) (TPPComponentPattern, bool, error) {
	record, found, err := dictionary.Lookup(surface, tables)
	if err != nil || !found {
		return TPPComponentPattern{}, found, err
	}
	for _, atom := range record.Atoms {
		pattern, matched, decodeErr := atom.ComponentPattern()
		if decodeErr != nil {
			return TPPComponentPattern{}, false, fmt.Errorf("decode TPP component atom for %q: %w", surface, decodeErr)
		}
		if matched {
			return pattern, true, nil
		}
	}
	return TPPComponentPattern{}, false, nil
}

// LookupTPPComponentPatternForSelector returns only the A-E atom selected by
// one native caller family. Other atoms in the same record are left intact in
// the record returned by LookupTypedText and are not substituted for a miss.
func (dictionary *TPPDictionary) LookupTPPComponentPatternForSelector(
	surface []byte,
	selector byte,
	tables EmbeddedKeyTables,
) (TPPComponentPattern, bool, error) {
	if selector < 'A' || selector > 'E' {
		return TPPComponentPattern{}, false, fmt.Errorf("TPP component selector %q is outside A-E", selector)
	}
	record, found, err := dictionary.Lookup(surface, tables)
	if err != nil || !found {
		return TPPComponentPattern{}, found, err
	}
	for _, atom := range record.Atoms {
		if atom.Tag != selector {
			continue
		}
		pattern, matched, decodeErr := atom.ComponentPattern()
		if decodeErr != nil {
			return TPPComponentPattern{}, false, fmt.Errorf("decode TPP %c component atom for %q: %w", selector, surface, decodeErr)
		}
		return pattern, matched, nil
	}
	return TPPComponentPattern{}, false, nil
}
