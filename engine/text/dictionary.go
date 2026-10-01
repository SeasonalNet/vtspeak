package text

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// EmbeddedRecord is one exact key and compact pronunciation payload from the
// shared embedded English lexicon. Payload retains the vendor bytes; callers
// can inspect the parsed form with ParsePhonePayload.
type EmbeddedRecord struct {
	Key     string
	Payload []byte
}

// EmbeddedDictionary provides exact, case-sensitive lookup. The legacy hash
// table is not needed here: its indexed record offsets are sufficient to
// validate and load the complete lexicon independently.
type EmbeddedDictionary struct {
	records map[string]EmbeddedRecord
}

// EmbeddedKeyTables contains the table data consumed by the observed
// FUN_1000fd70 key transform. CharacterMap is the low byte of each u16 entry;
// Pairs must be in the original lexicographic order.
type EmbeddedKeyTables struct {
	CharacterMap [256]byte
	Pairs        [][2]byte
}

// EncodeEmbeddedKey applies the recovered character-map and greedy pair
// compression used before embedded dictionary lookup. It accepts one
// canonical token as bytes; token normalization and extraction of the small
// lookup tables from their source remain separate stages.
func EncodeEmbeddedKey(surface []byte, tables EmbeddedKeyTables) ([]byte, error) {
	if len(surface) > 150 {
		return nil, errors.New("embedded dictionary key exceeds the observed 150-byte transform limit")
	}
	if len(tables.Pairs) > 127 {
		return nil, errors.New("embedded key pair table exceeds the observed 127-entry search range")
	}
	for i := 1; i < len(tables.Pairs); i++ {
		if tables.Pairs[i-1][0] > tables.Pairs[i][0] ||
			(tables.Pairs[i-1][0] == tables.Pairs[i][0] && tables.Pairs[i-1][1] >= tables.Pairs[i][1]) {
			return nil, fmt.Errorf("embedded key pair table is not strictly sorted at entry %d", i)
		}
	}

	mapped := make([]byte, len(surface))
	for i, value := range surface {
		mapped[i] = tables.CharacterMap[value]
	}
	encoded := make([]byte, 0, len(mapped))
	for i := 0; i < len(mapped); {
		if i+1 < len(mapped) {
			pair := [2]byte{mapped[i], mapped[i+1]}
			if pairIndex, ok := findKeyPair(tables.Pairs, pair); ok {
				encoded = append(encoded, byte(pairIndex+1)|0x80)
				i += 2
				continue
			}
		}
		encoded = append(encoded, mapped[i])
		i++
	}
	return encoded, nil
}

func findKeyPair(pairs [][2]byte, wanted [2]byte) (int, bool) {
	index := sort.Search(len(pairs), func(i int) bool {
		return pairs[i][0] > wanted[0] ||
			(pairs[i][0] == wanted[0] && pairs[i][1] >= wanted[1])
	})
	return index, index < len(pairs) && pairs[index] == wanted
}

// ParseEmbeddedDictionary joins the observed little-endian offset table with
// the lexicon. The first four index bytes are the loader header; remaining
// entries are u32 offsets into dictionaryData. Sorted offsets define the
// record partitions, as observed in hashidx_emb and engttsdict_emb.
func ParseEmbeddedDictionary(indexData, dictionaryData []byte) (*EmbeddedDictionary, error) {
	if len(indexData) < 4 {
		return nil, errors.New("embedded dictionary index is shorter than its header")
	}
	if (len(indexData)-4)%4 != 0 {
		return nil, errors.New("embedded dictionary index entries are not 32-bit words")
	}
	entryCount := (len(indexData) - 4) / 4
	if entryCount == 0 {
		return nil, errors.New("embedded dictionary index has no records")
	}
	offsets := make([]uint32, entryCount)
	for i := range offsets {
		offsets[i] = binary.LittleEndian.Uint32(indexData[4+i*4:])
	}
	sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
	if offsets[0] != 0 {
		return nil, errors.New("first embedded dictionary record does not start at offset zero")
	}

	records := make(map[string]EmbeddedRecord, entryCount)
	for i, start32 := range offsets {
		start := int(start32)
		end := len(dictionaryData)
		if i+1 < len(offsets) {
			end = int(offsets[i+1])
		}
		if start < 0 || start >= end || end > len(dictionaryData) {
			return nil, fmt.Errorf("invalid embedded dictionary record span [%d,%d)", start, end)
		}
		record := dictionaryData[start:end]
		keyEnd := indexByte(record, 0)
		if keyEnd <= 0 || keyEnd == len(record)-1 {
			return nil, fmt.Errorf("record at offset %d has an empty key or payload", start)
		}
		payloadTerminator := indexByte(record[keyEnd+1:], 0)
		if payloadTerminator != len(record)-keyEnd-2 {
			return nil, fmt.Errorf("record at offset %d does not end at one payload terminator", start)
		}
		key := string(record[:keyEnd])
		if _, exists := records[key]; exists {
			return nil, fmt.Errorf("duplicate embedded dictionary key %q", key)
		}
		payload := append([]byte(nil), record[keyEnd+1:]...)
		records[key] = EmbeddedRecord{Key: key, Payload: payload}
	}
	return &EmbeddedDictionary{records: records}, nil
}

func indexByte(data []byte, value byte) int {
	for i, item := range data {
		if item == value {
			return i
		}
	}
	return -1
}

// LoadEmbeddedDictionary loads the shared embedded English dictionary files
// without modifying them.
func LoadEmbeddedDictionary(root string) (*EmbeddedDictionary, error) {
	indexPath := filepath.Join(root, "hashidx_emb")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", indexPath, err)
	}
	dictionaryPath := filepath.Join(root, "engttsdict_emb")
	dictionaryData, err := os.ReadFile(dictionaryPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dictionaryPath, err)
	}
	dictionary, err := ParseEmbeddedDictionary(indexData, dictionaryData)
	if err != nil {
		return nil, fmt.Errorf("parse embedded dictionary: %w", err)
	}
	return dictionary, nil
}

// Lookup returns a copy of the exact matching record. Keys are not normalized;
// normalization and lexical fallback belong to the frontend stage.
func (dictionary *EmbeddedDictionary) Lookup(key string) (EmbeddedRecord, bool) {
	record, ok := dictionary.records[key]
	if !ok {
		return EmbeddedRecord{}, false
	}
	record.Payload = append([]byte(nil), record.Payload...)
	return record, true
}

// LookupSurface transforms a canonical surface token and looks up its
// encoded embedded-dictionary key. The caller remains responsible for the
// upstream token-normalization step and supplying the recovered key tables.
func (dictionary *EmbeddedDictionary) LookupSurface(surface []byte, tables EmbeddedKeyTables) (EmbeddedRecord, bool, error) {
	key, err := EncodeEmbeddedKey(surface, tables)
	if err != nil {
		return EmbeddedRecord{}, false, err
	}
	record, ok := dictionary.Lookup(string(key))
	return record, ok, nil
}

// ContainsPaul2013Surface performs the canonical key transform and exact
// embedded-record lookup used by the parser's punctuation-retention probe.
// It checks record presence without parsing or expanding the pronunciation.
func (dictionary *EmbeddedDictionary) ContainsPaul2013Surface(surface []byte) (bool, error) {
	if dictionary == nil {
		return false, errors.New("embedded dictionary is nil")
	}
	_, found, err := dictionary.LookupSurface(surface, Paul2013EmbeddedKeyTables())
	return found, err
}

// ResolveSurfacePronunciation performs key encoding, exact embedded lookup,
// payload parsing, and phone-ID expansion for one canonical surface token.
func (dictionary *EmbeddedDictionary) ResolveSurfacePronunciation(
	surface []byte,
	keyTables EmbeddedKeyTables,
	codebook PhoneIDCodebook,
) (SurfacePronunciation, bool, error) {
	record, found, err := dictionary.LookupSurface(surface, keyTables)
	if err != nil || !found {
		return SurfacePronunciation{}, found, err
	}
	payload, err := ParsePhonePayload(record.Payload)
	if err != nil {
		return SurfacePronunciation{}, false, fmt.Errorf("parse embedded pronunciation payload: %w", err)
	}
	alternatives, err := payload.ExpandPronunciations(codebook)
	if err != nil {
		return SurfacePronunciation{}, false, fmt.Errorf("expand embedded pronunciation IDs: %w", err)
	}
	return SurfacePronunciation{Record: record, Payload: payload, Alternatives: alternatives}, true, nil
}

// ResolvePaul2013SurfacePronunciation resolves and expands one canonical Paul
// 2013 token with the documented compact phone-ID table.
func (dictionary *EmbeddedDictionary) ResolvePaul2013SurfacePronunciation(
	surface []byte,
	keyTables EmbeddedKeyTables,
) (SurfacePronunciation, bool, error) {
	return dictionary.ResolveSurfacePronunciation(surface, keyTables, Paul2013PhoneIDCodebook())
}

// ResolvePaul2013Surface resolves one canonical token with both documented
// Paul 2013 lookup resources. Full text normalization is a separate stage.
func (dictionary *EmbeddedDictionary) ResolvePaul2013Surface(
	surface []byte,
) (SurfacePronunciation, bool, error) {
	return dictionary.ResolvePaul2013SurfacePronunciation(surface, Paul2013EmbeddedKeyTables())
}

// Len returns the number of validated records in the embedded dictionary.
func (dictionary *EmbeddedDictionary) Len() int { return len(dictionary.records) }

// Pronunciation contains compact phone-symbol IDs and any preceding path
// controls. Path meanings are intentionally left opaque.
type Pronunciation struct {
	Path  []byte
	Phone []byte
}

// PhoneIDCodebook maps each compact dictionary ID to its observed five-byte
// NUL-terminated internal-symbol slot.
type PhoneIDCodebook [256][5]byte

// ExpandedPronunciation retains opaque path controls and expands the compact
// phone IDs into internal symbol bytes.
type ExpandedPronunciation struct {
	Path    []byte
	Symbols []byte
}

// SurfacePronunciation is the result of resolving one canonical token through
// the embedded lexicon and expanding its compact IDs.
type SurfacePronunciation struct {
	Record       EmbeddedRecord
	Payload      PhonePayload
	Alternatives []ExpandedPronunciation
}

// PhonePayload retains the observed flag fields and decoded pronunciation
// alternatives from FUN_10003c50. The metadata field meanings remain unknown.
type PhonePayload struct {
	Flags          byte
	ResultType     byte
	Metadata       [4]bool
	Pronunciations []Pronunciation
}

const (
	// Paul2013ParsedDictionaryRecordSize spans the fixed fields, five 65-byte
	// phone slots, and five 20-byte path rows written by FUN_10003c50.
	Paul2013ParsedDictionaryRecordSize = 0x1b9
	paul2013ParsedDictionaryPhoneBase  = 0x10
	paul2013ParsedDictionaryPathBase   = 0x155
)

// BuildPaul2013ParsedDictionaryRecord projects a parsed payload into the
// caller-local record written by FUN_10003c50. It retains the native offsets
// consumed by FUN_100086c0; metadata meanings remain opaque.
func (payload PhonePayload) BuildPaul2013ParsedDictionaryRecord(
	codebook PhoneIDCodebook,
) ([]byte, error) {
	result := make([]byte, Paul2013ParsedDictionaryRecordSize)
	result[0] = payload.ResultType
	for index, set := range payload.Metadata {
		if set {
			binary.LittleEndian.PutUint16(result[2+index*2:], 1)
		}
	}

	expanded, err := payload.ExpandPronunciations(codebook)
	if err != nil {
		return nil, err
	}
	if len(expanded) > 5 {
		return nil, fmt.Errorf("parsed dictionary record has %d alternatives, maximum is 5", len(expanded))
	}
	binary.LittleEndian.PutUint32(result[0x0c:], uint32(len(expanded)))
	for index, alternative := range expanded {
		if len(alternative.Symbols) >= 65 {
			return nil, fmt.Errorf("parsed dictionary phone string %d has %d bytes, maximum is 64", index, len(alternative.Symbols))
		}
		phoneStart := paul2013ParsedDictionaryPhoneBase + index*65
		copy(result[phoneStart:phoneStart+65], alternative.Symbols)
		if len(expanded) == 1 {
			continue
		}
		if len(alternative.Path) >= 20 {
			return nil, fmt.Errorf("parsed dictionary path %d has %d bytes, maximum is 19", index, len(alternative.Path))
		}
		pathStart := paul2013ParsedDictionaryPathBase + index*20
		copy(result[pathStart:pathStart+20], alternative.Path)
		result[pathStart+len(alternative.Path)] = 0xff
	}
	if len(expanded) == 1 {
		result[paul2013ParsedDictionaryPathBase] = 0xff
	}
	return result, nil
}

// Paul2013DictionaryPhoneRows is the field-level output of the ordinary
// branch in FUN_1000d450. PathControlBytes holds its 0xff-terminated path
// stream; PhoneStrings contains the NUL-free internal-symbol string for each
// alternative. Their meanings remain opaque.
type Paul2013DictionaryPhoneRows struct {
	ResultType        byte
	Metadata          [4]bool
	AlternativeCount  int
	PathControlBytes  []byte
	PhoneStrings      [][]byte
	ContextMarker     byte
	HasContextMarker  bool
	MarkerSentinel    byte
	HasMarkerSentinel bool
	SelectedPhone     []byte
}

// Clone returns an independently owned copy of the row projection.
func (rows Paul2013DictionaryPhoneRows) Clone() Paul2013DictionaryPhoneRows {
	clone := rows
	clone.PathControlBytes = append([]byte(nil), rows.PathControlBytes...)
	clone.PhoneStrings = make([][]byte, len(rows.PhoneStrings))
	for index := range rows.PhoneStrings {
		clone.PhoneStrings[index] = append([]byte(nil), rows.PhoneStrings[index]...)
	}
	clone.SelectedPhone = append([]byte(nil), rows.SelectedPhone...)
	return clone
}

// SelectAlternativeForPathCode applies FUN_10010640 to this row projection.
// It returns the first aligned phone string whose path contains a code in the
// requested DLL class. Rows produced by the conditional marker branch have no
// alternatives and return found=false.
func (rows Paul2013DictionaryPhoneRows) SelectAlternativeForPathCode(
	requestedCode byte,
) (alternativeIndex int, phoneString []byte, found bool, err error) {
	if rows.AlternativeCount < 0 {
		return 0, nil, false, fmt.Errorf("pronunciation row has negative alternative count %d", rows.AlternativeCount)
	}
	if rows.AlternativeCount != len(rows.PhoneStrings) {
		return 0, nil, false, fmt.Errorf("pronunciation row declares %d alternatives but has %d phone strings", rows.AlternativeCount, len(rows.PhoneStrings))
	}
	if rows.AlternativeCount == 0 {
		return 0, nil, false, nil
	}
	return FindPaul2013PronunciationForPathCode(
		rows.PathControlBytes, rows.PhoneStrings, requestedCode,
	)
}

// BuildPaul2013DictionaryPhoneRows expands the parsed phone IDs and applies
// FUN_1000d450's default/final-row path-stream writes. The caller still owns
// the native row index, surface copy, and conditional marker branch.
func (payload PhonePayload) BuildPaul2013DictionaryPhoneRows(
	codebook PhoneIDCodebook,
) (Paul2013DictionaryPhoneRows, error) {
	return payload.BuildPaul2013DictionaryPhoneRowsForState(
		codebook,
		Paul2013DictionaryPhoneRowState{FinalRow: true},
	)
}

// Paul2013DictionaryPhoneRowState contains the row-index/marker inputs that
// select FUN_1000d450's default/final or conditional non-final branch.
// SelectPathPhone means the caller's two native pointer/flag gates both held.
type Paul2013DictionaryPhoneRowState struct {
	FinalRow        bool
	HasMarker       bool
	Marker          byte
	SelectPathPhone bool
}

// BuildPaul2013DictionaryPhoneRowsForState ports both branches of
// FUN_1000d450 over one parsed embedded payload. The selector and the gate
// that enables marker-based phone extraction remain explicit inputs.
func (payload PhonePayload) BuildPaul2013DictionaryPhoneRowsForState(
	codebook PhoneIDCodebook,
	state Paul2013DictionaryPhoneRowState,
) (Paul2013DictionaryPhoneRows, error) {
	result := Paul2013DictionaryPhoneRows{
		ResultType: payload.ResultType,
		Metadata:   payload.Metadata,
	}
	if state.HasMarker && !state.FinalRow {
		result.ContextMarker = state.Marker
		result.HasContextMarker = true
		result.MarkerSentinel = 0xff
		result.HasMarkerSentinel = true
		if state.SelectPathPhone {
			phone, found, err := payload.SelectPaul2013PronunciationByPathMarker(state.Marker, codebook)
			if err != nil {
				return Paul2013DictionaryPhoneRows{}, err
			}
			if found {
				result.SelectedPhone = phone
			}
		}
		return result, nil
	}

	expanded, err := payload.ExpandPronunciations(codebook)
	if err != nil {
		return Paul2013DictionaryPhoneRows{}, err
	}
	result.AlternativeCount = len(expanded)
	result.PathControlBytes = []byte{0xff}
	result.PhoneStrings = make([][]byte, len(expanded))
	if len(expanded) != 1 {
		result.PathControlBytes = make([]byte, 0, len(expanded)*2+1)
		for index, alternative := range expanded {
			result.PathControlBytes = append(result.PathControlBytes, paul2013PathRowBytes(alternative.Path)...)
			if index+1 < len(expanded) {
				result.PathControlBytes = append(result.PathControlBytes, 'd')
			}
		}
		result.PathControlBytes = append(result.PathControlBytes, 0xff)
	}
	for index := range expanded {
		result.PhoneStrings[index] = append([]byte(nil), expanded[index].Symbols...)
	}
	return result, nil
}

// ParsePhonePayload applies the recovered embedded payload grammar. A leading
// NUL denotes a marker record. Direct-ID form (bit 0) takes precedence over
// alternative form (bit 1), matching the legacy parser branch order.
func ParsePhonePayload(payload []byte) (PhonePayload, error) {
	var parsed PhonePayload
	if len(payload) == 0 {
		return parsed, errors.New("phone payload is empty")
	}
	parsed.Flags = payload[0]
	if parsed.Flags == 0 {
		return parsed, nil
	}
	if parsed.Flags&0x04 != 0 {
		parsed.ResultType = 'E'
	}
	if parsed.Flags&0x08 != 0 {
		parsed.ResultType = 'A'
	}
	parsed.Metadata = [4]bool{
		parsed.Flags&0x40 != 0,
		parsed.Flags&0x10 != 0,
		parsed.Flags&0x20 != 0,
		parsed.Flags&0x80 != 0,
	}

	switch {
	case parsed.Flags&0x01 != 0:
		phones, err := parseDirectPhoneIDs(payload[1:])
		if err != nil {
			return PhonePayload{}, err
		}
		parsed.Pronunciations = []Pronunciation{{Phone: phones}}
	case parsed.Flags&0x02 != 0:
		pronunciations, err := parsePhoneAlternatives(payload[1:])
		if err != nil {
			return PhonePayload{}, err
		}
		parsed.Pronunciations = pronunciations
	}
	return parsed, nil
}

// ExpandPronunciations expands every parsed ID sequence through codebook. The
// resulting symbols are internal engine bytes, not IPA or audio phone labels.
func (payload PhonePayload) ExpandPronunciations(codebook PhoneIDCodebook) ([]ExpandedPronunciation, error) {
	expanded := make([]ExpandedPronunciation, 0, len(payload.Pronunciations))
	for alternativeIndex, pronunciation := range payload.Pronunciations {
		symbols := make([]byte, 0, len(pronunciation.Phone))
		for phoneIndex, id := range pronunciation.Phone {
			slot := codebook[id]
			terminator := indexByte(slot[:], 0)
			if terminator <= 0 {
				return nil, fmt.Errorf("pronunciation %d phone ID 0x%02x at position %d has no symbol mapping", alternativeIndex, id, phoneIndex)
			}
			symbols = append(symbols, slot[:terminator]...)
		}
		expanded = append(expanded, ExpandedPronunciation{
			Path:    append([]byte(nil), pronunciation.Path...),
			Symbols: symbols,
		})
	}
	return expanded, nil
}

// SelectPaul2013PronunciationByPathMarker composes compact phone-ID expansion
// with FUN_10003f10's marker lookup. It returns the internal phone-symbol
// string aligned with the selected path row; the source of selector remains
// outside this operation.
func (payload PhonePayload) SelectPaul2013PronunciationByPathMarker(
	selector byte,
	codebook PhoneIDCodebook,
) ([]byte, bool, error) {
	expanded, err := payload.ExpandPronunciations(codebook)
	if err != nil {
		return nil, false, err
	}
	paths := make([][]byte, len(expanded))
	phones := make([][]byte, len(expanded))
	for index := range expanded {
		paths[index] = expanded[index].Path
		phones[index] = expanded[index].Symbols
	}
	return SelectPaul2013PronunciationByPathMarker(paths, phones, selector)
}

func parseDirectPhoneIDs(data []byte) ([]byte, error) {
	terminator := indexByte(data, 0)
	if terminator != len(data)-1 {
		return nil, errors.New("direct phone IDs must end at the payload terminator")
	}
	return append([]byte(nil), data[:terminator]...), nil
}

func parsePhoneAlternatives(data []byte) ([]Pronunciation, error) {
	if len(data) == 0 || data[len(data)-1] != 0 {
		return nil, errors.New("alternative payload is not NUL-terminated")
	}
	body := data[:len(data)-1]
	var alternatives []Pronunciation
	for len(body) > 0 {
		end := indexByte(body, 0xff)
		if end < 0 {
			end = len(body)
		}
		part := body[:end]
		separator := indexByte(part, '|')
		if separator < 0 {
			return nil, errors.New("pronunciation alternative has no path/phone separator")
		}
		path := append([]byte(nil), part[:separator]...)
		for i := range path {
			if path[i] == 0 || path[i] == 0xff {
				return nil, errors.New("path contains a reserved byte")
			}
			path[i]--
		}
		phones := append([]byte(nil), part[separator+1:]...)
		for _, phoneID := range phones {
			if phoneID == 0 || phoneID == 0xff {
				return nil, errors.New("phone ID sequence contains a reserved byte")
			}
		}
		alternatives = append(alternatives, Pronunciation{Path: path, Phone: phones})
		if end == len(body) {
			break
		}
		body = body[end+1:]
	}
	if len(alternatives) == 0 {
		return nil, errors.New("alternative payload has no pronunciation records")
	}
	if len(alternatives) > 5 {
		return nil, errors.New("alternative payload exceeds the observed five-record limit")
	}
	return alternatives, nil
}
