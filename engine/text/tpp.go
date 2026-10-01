package text

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const paul2013TPPRecordCount = 31550

// TPPAtom retains one typed code from a TPP payload. Tag and suffix are raw
// schema fields; the linguistic meaning of these codes is not established.
type TPPAtom struct {
	Tag    byte
	Suffix []byte
}

// TPPRecord stores a compressed key and its typed payload without assigning
// semantic names to the vendor's code letters or numeric suffixes.
type TPPRecord struct {
	Key     []byte
	Payload []byte
	Atoms   []TPPAtom
}

// TPPDictionary indexes the sequentially framed shared TPP records by their
// original compressed key bytes.
type TPPDictionary struct {
	records map[string]TPPRecord
}

// NumericByte returns the observed one-byte numeric result for the F and G
// suffix families. Other tags, including AX, have no established numeric
// operation. The DLL parses a signed integer and narrows it to one byte.
func (atom TPPAtom) NumericByte() (byte, bool, error) {
	if atom.Tag != 'F' && atom.Tag != 'G' {
		return 0, false, nil
	}
	if len(atom.Suffix) == 0 {
		return 0, false, errors.New("TPP numeric suffix is empty")
	}
	attributes := [256]byte{}
	for digit := byte('0'); digit <= '9'; digit++ {
		attributes[digit] = 0x04
	}
	for index, digit := range atom.Suffix {
		if digit < '0' || digit > '9' {
			return 0, false, fmt.Errorf("parse TPP %c numeric suffix %q: byte 0x%02x at offset %d is not decimal", atom.Tag, atom.Suffix, digit, index)
		}
	}
	value, hasDigits := ParsePaul2013NativeInteger(atom.Suffix, attributes)
	if !hasDigits {
		return 0, false, fmt.Errorf("parse TPP %c numeric suffix %q: no digits", atom.Tag, atom.Suffix)
	}
	return byte(value), true, nil
}

// ParseTPPDictionary parses NUL-terminated key/payload pairs to exact EOF and
// validates the observed one- and two-atom code grammar.
func ParseTPPDictionary(data []byte) (*TPPDictionary, error) {
	dictionary := &TPPDictionary{records: make(map[string]TPPRecord)}
	offset := 0
	for offset < len(data) {
		keyEnd := bytes.IndexByte(data[offset:], 0)
		if keyEnd < 0 || keyEnd == 0 {
			return nil, fmt.Errorf("TPP record at byte %d has a missing or empty key", offset)
		}
		keyEnd += offset
		payloadStart := keyEnd + 1
		payloadEndRelative := bytes.IndexByte(data[payloadStart:], 0)
		if payloadEndRelative < 0 || payloadEndRelative == 0 {
			return nil, fmt.Errorf("TPP record at byte %d has a missing or empty payload", offset)
		}
		payloadEnd := payloadStart + payloadEndRelative
		key := append([]byte(nil), data[offset:keyEnd]...)
		payload := append([]byte(nil), data[payloadStart:payloadEnd]...)
		atoms, err := parseTPPAtoms(payload)
		if err != nil {
			return nil, fmt.Errorf("TPP record at byte %d: %w", offset, err)
		}
		if _, exists := dictionary.records[string(key)]; exists {
			return nil, fmt.Errorf("duplicate TPP key % x", key)
		}
		dictionary.records[string(key)] = TPPRecord{Key: key, Payload: payload, Atoms: atoms}
		offset = payloadEnd + 1
	}
	if len(dictionary.records) == 0 {
		return nil, errors.New("TPP dictionary contains no records")
	}
	if len(dictionary.records) != paul2013TPPRecordCount {
		return nil, fmt.Errorf("TPP dictionary has %d records, want observed count %d", len(dictionary.records), paul2013TPPRecordCount)
	}
	return dictionary, nil
}

func parseTPPAtoms(payload []byte) ([]TPPAtom, error) {
	parts := bytes.Split(payload, []byte{' '})
	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("payload has %d atoms, want one or two", len(parts))
	}
	atoms := make([]TPPAtom, len(parts))
	for index, part := range parts {
		if len(part) < 2 || part[0] < 'A' || part[0] > 'G' {
			return nil, fmt.Errorf("atom %d has invalid tag or suffix %q", index, part)
		}
		if part[0] == 'A' && len(part) == 2 && part[1] == 'X' {
			atoms[index] = TPPAtom{Tag: 'A', Suffix: []byte{'X'}}
			continue
		}
		for _, digit := range part[1:] {
			if digit < '0' || digit > '9' {
				return nil, fmt.Errorf("atom %d has invalid suffix %q", index, part[1:])
			}
		}
		atoms[index] = TPPAtom{Tag: part[0], Suffix: append([]byte(nil), part[1:]...)}
	}
	if len(atoms) == 2 {
		if bytes.Equal(atoms[0].Suffix, []byte{'X'}) || bytes.Equal(atoms[1].Suffix, []byte{'X'}) {
			return nil, errors.New("AX records are observed only as single-atom payloads")
		}
		allowed := (atoms[0].Tag == 'A' &&
			(atoms[1].Tag == 'A' || atoms[1].Tag == 'G')) ||
			(atoms[0].Tag == 'B' && atoms[1].Tag == 'B') ||
			(atoms[0].Tag == 'C' && atoms[1].Tag == 'C')
		if !allowed {
			return nil, fmt.Errorf("two-atom sequence %c %c is not observed", atoms[0].Tag, atoms[1].Tag)
		}
	}
	return atoms, nil
}

// LoadTPPDictionary reads the shared sequential TPP text dictionary.
func LoadTPPDictionary(dictionaryRoot string) (*TPPDictionary, error) {
	path := filepath.Join(dictionaryRoot, "tppdict_eng")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	dictionary, err := ParseTPPDictionary(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return dictionary, nil
}

// Lookup finds an exact surface after applying the observed embedded-key
// transform. The caller supplies the already-selected normalization form.
func (dictionary *TPPDictionary) Lookup(surface []byte, tables EmbeddedKeyTables) (TPPRecord, bool, error) {
	if dictionary == nil {
		return TPPRecord{}, false, errors.New("TPP dictionary is nil")
	}
	key, err := EncodeEmbeddedKey(surface, tables)
	if err != nil {
		return TPPRecord{}, false, err
	}
	record, ok := dictionary.records[string(key)]
	if !ok {
		return TPPRecord{}, false, nil
	}
	record.Key = append([]byte(nil), record.Key...)
	record.Payload = append([]byte(nil), record.Payload...)
	record.Atoms = append([]TPPAtom(nil), record.Atoms...)
	for index := range record.Atoms {
		record.Atoms[index].Suffix = append([]byte(nil), record.Atoms[index].Suffix...)
	}
	return record, true, nil
}

// LookupNumericText returns the raw decimal suffix for an exact F or G TPP
// atom. The native F/G callers consume these suffixes as numeric text.
func (dictionary *TPPDictionary) LookupNumericText(
	surface []byte,
	tag byte,
	tables EmbeddedKeyTables,
) ([]byte, bool, error) {
	if tag != 'F' && tag != 'G' {
		return nil, false, fmt.Errorf("unsupported Paul 2013 TPP numeric selector %q", tag)
	}
	record, found, err := dictionary.Lookup(surface, tables)
	if err != nil || !found {
		return nil, false, err
	}
	for _, atom := range record.Atoms {
		if atom.Tag == tag {
			return append([]byte(nil), atom.Suffix...), true, nil
		}
	}
	return nil, false, nil
}

// LookupSelectedText returns the suffix selected by the native TPP selector
// path. If a record contains multiple atoms, it returns the first atom whose
// discriminator matches selector, including a secondary atom such as G95 in
// "A0 G95". It preserves AX as its literal X suffix and does not assign
// linguistic meaning to any returned bytes.
func (dictionary *TPPDictionary) LookupSelectedText(
	surface []byte,
	selector byte,
	tables EmbeddedKeyTables,
) ([]byte, bool, error) {
	atom, found, err := dictionary.LookupSelectedAtom(surface, selector, tables)
	if err != nil || !found {
		return nil, found, err
	}
	return append([]byte(nil), atom.Suffix...), true, nil
}

// LookupSelectedAtom returns the first parsed atom whose discriminator
// matches the native caller's selector. The suffix is copied so callers can
// retain or transform it without changing dictionary storage.
func (dictionary *TPPDictionary) LookupSelectedAtom(
	surface []byte,
	selector byte,
	tables EmbeddedKeyTables,
) (TPPAtom, bool, error) {
	if selector < 'A' || selector > 'G' {
		return TPPAtom{}, false, fmt.Errorf("unsupported Paul 2013 TPP selector %q", selector)
	}
	record, found, err := dictionary.Lookup(surface, tables)
	if err != nil || !found {
		return TPPAtom{}, found, err
	}
	for _, atom := range record.Atoms {
		if atom.Tag == selector {
			return TPPAtom{Tag: atom.Tag, Suffix: append([]byte(nil), atom.Suffix...)}, true, nil
		}
	}
	return TPPAtom{}, false, nil
}
