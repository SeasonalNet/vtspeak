package text

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const paul2013ExceptionCategoryCount = 4

// Paul2013ExceptionDispatchGateInput contains the fields tested immediately
// before FUN_10008dc0 in FUN_10007520. SourceTokenStatus is the signed byte at
// source-record +0x30; PhoneRowMarker is phone/context-row +0x29; SourceKind
// and ContextFormIsEmpty read source-record +0x23 and the first byte at +0x52;
// ContextPreparationCode is the low short returned by FUN_100091b0. This is a
// post-preprocessing snapshot, not a producer of those fields from text.
type Paul2013ExceptionDispatchGateInput struct {
	SourceTokenStatus      int8
	PhoneRowMarker         byte
	SourceKind             byte
	ContextFormIsEmpty     bool
	ContextPreparationCode int16
}

// BuildPaul2013ExceptionDispatchGateFromRows extracts the fields read by the
// FUN_10007520 dispatch gate from one 0x94-byte source/parser row and one
// 0x70-byte phone/context row. The low-short result is explicit in this
// helper; DerivePaul2013PhoneContextPreparationFromRows can supply it from the
// recovered return branches without reproducing the remaining row mutations.
func BuildPaul2013ExceptionDispatchGateFromRows(
	sourceParserRows []byte,
	sourceTokenIndex uint16,
	phoneContextRow []byte,
	contextPreparationCode int16,
) (Paul2013ExceptionDispatchGateInput, error) {
	const (
		sourceParserRowSize     = 0x94
		sourceKindOffset        = 0x23
		sourceTokenStatusOffset = 0x30
		sourceContextFormOffset = 0x52
	)
	rowStart := int(sourceTokenIndex) * sourceParserRowSize
	if rowStart > len(sourceParserRows) || len(sourceParserRows)-rowStart < sourceParserRowSize {
		return Paul2013ExceptionDispatchGateInput{}, fmt.Errorf(
			"source/parser row %d is outside %d available bytes", sourceTokenIndex, len(sourceParserRows),
		)
	}
	if len(phoneContextRow) < Paul2013PhoneContextRowSize {
		return Paul2013ExceptionDispatchGateInput{}, fmt.Errorf(
			"phone/context row has %d bytes, need %d", len(phoneContextRow), Paul2013PhoneContextRowSize,
		)
	}
	row := sourceParserRows[rowStart : rowStart+sourceParserRowSize]
	return Paul2013ExceptionDispatchGateInput{
		SourceTokenStatus:      int8(row[sourceTokenStatusOffset]),
		PhoneRowMarker:         phoneContextRow[0x29],
		SourceKind:             row[sourceKindOffset],
		ContextFormIsEmpty:     row[sourceContextFormOffset] == 0,
		ContextPreparationCode: contextPreparationCode,
	}, nil
}

// ShouldDispatchPaul2013PronunciationException ports the caller's complete
// gate for entering FUN_10008dc0, using the low-short status it checks from
// FUN_100091b0.
func ShouldDispatchPaul2013PronunciationException(input Paul2013ExceptionDispatchGateInput) bool {
	return (input.SourceTokenStatus == -1 || input.PhoneRowMarker == 0) &&
		(input.SourceKind == 'U' || input.ContextFormIsEmpty) &&
		input.ContextPreparationCode == 0
}

// ShouldRunPaul2013ModelContextPostHandlers ports the FUN_10007520 outer
// counted-row gate. Unlike exception dispatch, this gate does not include the
// FUN_100091b0 preparation short: a nonzero short skips normalization and
// exception lookup but still reaches the post-handler block when these fields
// pass.
func ShouldRunPaul2013ModelContextPostHandlers(input Paul2013ExceptionDispatchGateInput) bool {
	return (input.SourceTokenStatus == -1 || input.PhoneRowMarker == 0) &&
		(input.SourceKind == 'U' || input.ContextFormIsEmpty)
}

// Paul2013UnsignedCharacterAttributeTable returns the raw 256-byte
// DAT_1007e188 table from vt_pau.dll (VA 0x1007e188, file offset 0x7e188).
// Use only for consumers that index it with an unsigned byte.
func Paul2013UnsignedCharacterAttributeTable() [256]byte {
	var attributes [256]byte
	for value := byte(0); value <= 0x20; value++ {
		attributes[value] = 0x01
	}
	for value := byte(0x08); value <= 0x0d; value++ {
		attributes[value] = 0x02
	}
	attributes[' '] = 0x04
	for value := byte('!'); value <= '@'; value++ {
		attributes[value] = 0x08
	}
	for value := byte('0'); value <= '9'; value++ {
		attributes[value] = 0x30
	}
	for value := byte('A'); value <= 'F'; value++ {
		attributes[value] = 0xa0
	}
	for value := byte('G'); value <= 'Z'; value++ {
		attributes[value] = 0x80
	}
	for value := byte('['); value <= '`'; value++ {
		attributes[value] = 0x08
	}
	for value := byte('a'); value <= 'f'; value++ {
		attributes[value] = 0x60
	}
	for value := byte('g'); value <= 'z'; value++ {
		attributes[value] = 0x40
	}
	for value := byte('{'); value <= '~'; value++ {
		attributes[value] = 0x08
	}
	attributes[0x7f] = 0x01
	for value := byte(0x80); value <= 0x82; value++ {
		attributes[value] = 0x08
	}
	attributes[0x83] = 0x40
	for value := byte(0x84); value <= 0x89; value++ {
		attributes[value] = 0x08
	}
	attributes[0x8a] = 0x80
	attributes[0x8b] = 0x08
	attributes[0x8c] = 0x80
	for value := byte(0x8d); value <= 0x99; value++ {
		attributes[value] = 0x08
	}
	attributes[0x9a] = 0x40
	attributes[0x9b] = 0x08
	attributes[0x9c] = 0x40
	attributes[0x9d] = 0x08
	attributes[0x9e] = 0x08
	attributes[0x9f] = 0x80
	for value := byte(0xa0); value <= 0xbf; value++ {
		attributes[value] = 0x08
	}
	for value := byte(0xc0); value <= 0xd6; value++ {
		attributes[value] = 0x80
	}
	attributes[0xd7] = 0x08
	for value := byte(0xd8); value <= 0xde; value++ {
		attributes[value] = 0x80
	}
	attributes[0xdf] = 0x40
	for value := byte(0xe0); value <= 0xe6; value++ {
		attributes[value] = 0x40
	}
	attributes[0xe7] = 0x08
	for value := byte(0xe8); value <= 0xf6; value++ {
		attributes[value] = 0x40
	}
	attributes[0xf7] = 0x08
	for value := 0xf8; value <= 0xff; value++ {
		attributes[byte(value)] = 0x40
	}
	return attributes
}

// Paul2013ExceptionCharacterAttributes returns the effective lookup values
// for native signed-char string consumers. Bytes 0x80-0xff index the 128-byte
// prefix before DAT_1007e188; that prefix is zero in vt_pau.dll. Unsigned-byte
// consumers should use Paul2013UnsignedCharacterAttributeTable instead.
func Paul2013ExceptionCharacterAttributes() [256]byte {
	attributes := Paul2013UnsignedCharacterAttributeTable()
	for value := byte(0x80); ; value++ {
		attributes[value] = 0
		if value == 0xff {
			break
		}
	}
	return attributes
}

// ExceptionRow stores one compressed surface key and its compact phone-code
// value. The bytes are retained verbatim because the resource uses the
// embedded dictionary's private key transform and code alphabet.
type ExceptionRow struct {
	Key   []byte
	Value []byte
}

// ExceptionMatch is one exact normalized lookup result over a consecutive
// caller-selected token prefix.
type ExceptionMatch struct {
	PhoneCodes []byte
	Category   uint32
	TokenCount int
}

// DecodePhoneRows applies the exception phone-code writer using one explicit
// surface-component count for each matched destination token.
func (match ExceptionMatch) DecodePhoneRows(delimiterCounts []int) ([][]CMUPhone, error) {
	if match.TokenCount <= 0 {
		return nil, fmt.Errorf("exception match has invalid token count %d", match.TokenCount)
	}
	if len(delimiterCounts) != match.TokenCount {
		return nil, fmt.Errorf("exception match spans %d tokens but received %d delimiter counts", match.TokenCount, len(delimiterCounts))
	}
	return DecodePaul2013ExceptionPhoneRows(match.PhoneCodes, delimiterCounts)
}

// DecodePhoneRowsForSurfaces derives FUN_1000ca30's destination-row counts
// from the matched source surfaces, then applies FUN_1000ca50's code split.
func (match ExceptionMatch) DecodePhoneRowsForSurfaces(surfaces []string) ([][]CMUPhone, error) {
	if match.TokenCount <= 0 {
		return nil, fmt.Errorf("exception match has invalid token count %d", match.TokenCount)
	}
	if len(surfaces) < match.TokenCount {
		return nil, fmt.Errorf("exception match spans %d tokens but received %d source surfaces", match.TokenCount, len(surfaces))
	}
	delimiterCounts := make([]int, match.TokenCount)
	for index, surface := range surfaces[:match.TokenCount] {
		delimiterCounts[index] = Paul2013ExceptionSurfaceDelimiterCount(surface)
	}
	return match.DecodePhoneRows(delimiterCounts)
}

// Paul2013ExceptionSurfaceDelimiterCount ports FUN_1000ca30: count literal
// hyphens before NUL and add one. The input should match the surface buffer at
// +7 of one native phone/context row.
func Paul2013ExceptionSurfaceDelimiterCount(surface string) int {
	hyphens := 0
	for index := 0; index < len(surface); index++ {
		if surface[index] == 0 {
			break
		}
		if surface[index] == '-' {
			hyphens++
		}
	}
	return hyphens + 1
}

// ExceptionDictionary keeps the four component-count groups from
// dict-eng/exceptdict in source order. Keys within each group are sorted, as
// required by the native binary-search lookup.
type ExceptionDictionary struct {
	minComponents uint32
	maxComponents uint32
	groups        [paul2013ExceptionCategoryCount][]ExceptionRow
}

// NormalizePaul2013ExceptionSurface ports FUN_1000c9c0 for explicit
// character-attribute flags. It maps accepted bytes through the observed
// character table, preserves apostrophes, optionally preserves hyphens, and
// returns one plus the number of accepted hyphens as the category selector.
// A NUL byte terminates the supplied surface, matching the native C-string
// input. The source of the per-byte attribute flags remains caller supplied.
func NormalizePaul2013ExceptionSurface(
	surface []byte,
	allowHyphen bool,
	attributes [256]byte,
	characterMap [256]byte,
) ([]byte, uint32, error) {
	componentCount := uint32(1)
	normalized := make([]byte, 0, len(surface))
	for offset, value := range surface {
		if value == 0 {
			return normalized, componentCount, nil
		}
		if attributes[value]&0xc0 == 0 {
			switch value {
			case '\'':
				normalized = append(normalized, '\'')
			case '-':
				if !allowHyphen {
					return nil, 0, fmt.Errorf("hyphen at byte %d is not allowed", offset)
				}
				normalized = append(normalized, '-')
				componentCount++
			default:
				return nil, 0, fmt.Errorf("byte 0x%02x at offset %d is not accepted by the exception normalizer", value, offset)
			}
			continue
		}
		normalized = append(normalized, characterMap[value])
	}
	return normalized, componentCount, nil
}

// ParseExceptionDictionary parses the observed little-endian exceptdict
// framing and rejects incomplete, reordered, duplicate, or trailing data.
func ParseExceptionDictionary(data []byte) (*ExceptionDictionary, error) {
	reader := exceptionReader{data: data}
	minimum, err := reader.u32("minimum component category")
	if err != nil {
		return nil, err
	}
	maximum, err := reader.u32("maximum component category")
	if err != nil {
		return nil, err
	}
	if minimum != 1 || maximum != paul2013ExceptionCategoryCount {
		return nil, fmt.Errorf("exception dictionary declares categories %d..%d, want 1..%d", minimum, maximum, paul2013ExceptionCategoryCount)
	}
	dictionary := &ExceptionDictionary{minComponents: minimum, maxComponents: maximum}
	for expectedCategory := uint32(1); expectedCategory <= maximum; expectedCategory++ {
		category, err := reader.u32("category id")
		if err != nil {
			return nil, fmt.Errorf("exception group %d: %w", expectedCategory, err)
		}
		if category != expectedCategory {
			return nil, fmt.Errorf("exception group %d declares category %d", expectedCategory, category)
		}
		rowCount, err := reader.u32("row count")
		if err != nil {
			return nil, fmt.Errorf("exception category %d: %w", category, err)
		}
		if uint64(rowCount) > uint64(reader.remaining()/8) {
			return nil, fmt.Errorf("exception category %d declares %d rows beyond remaining input", category, rowCount)
		}
		rows := make([]ExceptionRow, 0, int(rowCount))
		for rowIndex := uint32(0); rowIndex < rowCount; rowIndex++ {
			keyLength, err := reader.u32("key length")
			if err != nil {
				return nil, fmt.Errorf("exception category %d row %d: %w", category, rowIndex, err)
			}
			valueLength, err := reader.u32("value length")
			if err != nil {
				return nil, fmt.Errorf("exception category %d row %d: %w", category, rowIndex, err)
			}
			key, err := reader.bytes(keyLength, "key")
			if err != nil {
				return nil, fmt.Errorf("exception category %d row %d: %w", category, rowIndex, err)
			}
			value, err := reader.bytes(valueLength, "value")
			if err != nil {
				return nil, fmt.Errorf("exception category %d row %d: %w", category, rowIndex, err)
			}
			if len(key) == 0 || len(value) == 0 {
				return nil, fmt.Errorf("exception category %d row %d has an empty key or value", category, rowIndex)
			}
			if rowIndex > 0 && bytes.Compare(rows[rowIndex-1].Key, key) >= 0 {
				return nil, fmt.Errorf("exception category %d keys are not strictly sorted at row %d", category, rowIndex)
			}
			rows = append(rows, ExceptionRow{Key: key, Value: value})
		}
		dictionary.groups[category-1] = rows
	}
	if reader.remaining() != 0 {
		return nil, fmt.Errorf("exception dictionary has %d trailing bytes", reader.remaining())
	}
	return dictionary, nil
}

// LoadExceptionDictionary reads the shared English exceptdict resource.
func LoadExceptionDictionary(dictionaryRoot string) (*ExceptionDictionary, error) {
	path := filepath.Join(dictionaryRoot, "exceptdict")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	dictionary, err := ParseExceptionDictionary(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return dictionary, nil
}

// Lookup returns the exact key's compact value in the category selected by
// the normalized token component count. The input key is expected to have
// already passed through the native embedded-key transform.
func (dictionary *ExceptionDictionary) Lookup(category uint32, encodedKey []byte) ([]byte, bool) {
	if dictionary == nil || category < dictionary.minComponents || category > dictionary.maxComponents {
		return nil, false
	}
	rows := dictionary.groups[category-1]
	low, high := 0, len(rows)
	for low < high {
		middle := low + (high-low)/2
		if bytes.Compare(rows[middle].Key, encodedKey) < 0 {
			low = middle + 1
		} else {
			high = middle
		}
	}
	if low == len(rows) || !bytes.Equal(rows[low].Key, encodedKey) {
		return nil, false
	}
	return append([]byte(nil), rows[low].Value...), true
}

// LookupSurface applies the caller-supplied native character-attribute table,
// observed character map, and embedded-key compression before selecting the
// exception group by the normalized hyphen count.
func (dictionary *ExceptionDictionary) LookupSurface(
	surface []byte,
	allowHyphen bool,
	attributes [256]byte,
	tables EmbeddedKeyTables,
) ([]byte, uint32, bool, error) {
	normalized, category, err := NormalizePaul2013ExceptionSurface(
		surface, allowHyphen, attributes, tables.CharacterMap,
	)
	if err != nil {
		return nil, 0, false, err
	}
	key, err := EncodeEmbeddedKey(normalized, tables)
	if err != nil {
		return nil, category, false, err
	}
	value, found := dictionary.Lookup(category, key)
	return value, category, found, nil
}

// LookupPaul2013Surface applies the local Paul 2013 character-attribute mask
// and key tables before exact exception lookup.
func (dictionary *ExceptionDictionary) LookupPaul2013Surface(
	surface []byte,
	allowHyphen bool,
) ([]byte, uint32, bool, error) {
	return dictionary.LookupSurface(
		surface,
		allowHyphen,
		Paul2013ExceptionCharacterAttributes(),
		Paul2013EmbeddedKeyTables(),
	)
}

// LookupPaul2013SurfaceSequence ports the exact-prefix composition in
// FUN_10008dc0. It normalizes each supplied surface, joins successive surfaces
// with hyphens, accumulates the one-to-four-component category, and retains
// the longest prefix that has an exact dictionary entry. This compatibility
// form supplies no X-row retry markers; callers with those markers can use
// LookupPaul2013SurfaceSequenceWithHPrefixRetry. Token eligibility remains
// caller supplied.
func (dictionary *ExceptionDictionary) LookupPaul2013SurfaceSequence(
	surfaces []string,
) (ExceptionMatch, bool, error) {
	return dictionary.LookupPaul2013SurfaceSequenceWithHPrefixRetry(surfaces, nil)
}

// LookupPaul2013SurfaceSequenceWithHPrefixRetry ports the X-marker retry in
// FUN_10008dc0. The flags correspond to its source marker at row +0x66; when
// a joined-prefix lookup misses on a flagged row, the native caller retries
// the current normalized surface as "h'" plus that surface. The producer of
// those marker bytes remains outside this text helper.
func (dictionary *ExceptionDictionary) LookupPaul2013SurfaceSequenceWithHPrefixRetry(
	surfaces []string,
	retryHPrefix []bool,
) (ExceptionMatch, bool, error) {
	if dictionary == nil {
		return ExceptionMatch{}, false, errors.New("exception dictionary is nil")
	}
	if len(surfaces) == 0 {
		return ExceptionMatch{}, false, errors.New("exception surface sequence is empty")
	}
	if retryHPrefix != nil && len(retryHPrefix) != len(surfaces) {
		return ExceptionMatch{}, false, fmt.Errorf(
			"received %d exception retry markers for %d surfaces",
			len(retryHPrefix), len(surfaces),
		)
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	tables := Paul2013EmbeddedKeyTables()
	var joined strings.Builder
	var accumulatedCategory uint32
	var match ExceptionMatch
	found := false
	for tokenIndex, surface := range surfaces {
		normalized, category, err := NormalizePaul2013ExceptionSurface(
			[]byte(surface), true, attributes, tables.CharacterMap,
		)
		if err != nil {
			return ExceptionMatch{}, false, fmt.Errorf("normalize exception surface %d %q: %w", tokenIndex, surface, err)
		}
		if tokenIndex != 0 {
			joined.WriteByte('-')
		}
		joined.Write(normalized)
		accumulatedCategory += category
		if accumulatedCategory > dictionary.maxComponents {
			break
		}
		key, err := EncodeEmbeddedKey([]byte(joined.String()), tables)
		if err != nil {
			return ExceptionMatch{}, false, fmt.Errorf("encode exception surface prefix through token %d: %w", tokenIndex, err)
		}
		phoneCodes, ok := dictionary.Lookup(accumulatedCategory, key)
		if !ok && retryHPrefix != nil && retryHPrefix[tokenIndex] {
			retrySurface := make([]byte, 0, len(normalized)+2)
			retrySurface = append(retrySurface, 'h', '\'')
			retrySurface = append(retrySurface, normalized...)
			retryKey, err := EncodeEmbeddedKey(retrySurface, tables)
			if err != nil {
				return ExceptionMatch{}, false, fmt.Errorf("encode h-prefix exception retry for token %d: %w", tokenIndex, err)
			}
			phoneCodes, ok = dictionary.Lookup(accumulatedCategory, retryKey)
		}
		if ok {
			match = ExceptionMatch{
				PhoneCodes: phoneCodes,
				Category:   accumulatedCategory,
				TokenCount: tokenIndex + 1,
			}
			found = true
		}
	}
	return match, found, nil
}

// SplitPaul2013ExceptionPhoneCodes ports the delimiter decision in
// FUN_1000ca50. Each supplied count is FUN_1000ca30's one-plus-hyphens result
// for a destination row. A 'd' is retained while the row's count is below its
// boundary count; the boundary 'd' is consumed and starts the next row. The
// caller supplies the reachable destination-row counts because the upstream
// exception matcher owns their layout.
func SplitPaul2013ExceptionPhoneCodes(code []byte, delimiterCounts []int) ([][]byte, error) {
	if len(delimiterCounts) == 0 {
		return nil, errors.New("exception phone-code split has no destination rows")
	}
	for rowIndex, count := range delimiterCounts {
		if count < 1 {
			return nil, fmt.Errorf("destination row %d has invalid delimiter count %d", rowIndex, count)
		}
	}
	rows := make([][]byte, 1, len(delimiterCounts))
	rows[0] = make([]byte, 0)
	rowIndex, seenDelimiters := 0, 0
	for _, value := range code {
		if value == 0 {
			break
		}
		if value != 'd' {
			rows[rowIndex] = append(rows[rowIndex], value)
			continue
		}
		seenDelimiters++
		if seenDelimiters < delimiterCounts[rowIndex] {
			rows[rowIndex] = append(rows[rowIndex], value)
			continue
		}
		rowIndex++
		if rowIndex >= len(delimiterCounts) {
			return nil, fmt.Errorf("exception code crosses beyond %d supplied destination rows", len(delimiterCounts))
		}
		rows = append(rows, make([]byte, 0))
		seenDelimiters = 0
	}
	return rows, nil
}

// DecodePaul2013ExceptionPhoneRows applies the recovered row split and labels
// each resulting internal phone-symbol byte with the runtime-backed CMU map.
// This low-level form accepts explicit counts; ExceptionMatch can derive them
// from source surfaces when those row buffers are available.
func DecodePaul2013ExceptionPhoneRows(code []byte, delimiterCounts []int) ([][]CMUPhone, error) {
	symbolRows, err := SplitPaul2013ExceptionPhoneCodes(code, delimiterCounts)
	if err != nil {
		return nil, err
	}
	phoneRows := make([][]CMUPhone, len(symbolRows))
	for rowIndex, symbols := range symbolRows {
		phones, err := DecodeCMUPhones(symbols)
		if err != nil {
			return nil, fmt.Errorf("decode exception phone row %d: %w", rowIndex, err)
		}
		phoneRows[rowIndex] = phones
	}
	return phoneRows, nil
}

// ExceptionGroupSizes reports the observed row count for each 1..4 component
// category without exposing mutable internal row slices.
func (dictionary *ExceptionDictionary) ExceptionGroupSizes() [paul2013ExceptionCategoryCount]int {
	var sizes [paul2013ExceptionCategoryCount]int
	if dictionary == nil {
		return sizes
	}
	for i := range dictionary.groups {
		sizes[i] = len(dictionary.groups[i])
	}
	return sizes
}

type exceptionReader struct {
	data []byte
	off  int
}

func (reader *exceptionReader) remaining() int { return len(reader.data) - reader.off }

func (reader *exceptionReader) u32(name string) (uint32, error) {
	if reader.remaining() < 4 {
		return 0, fmt.Errorf("truncated %s at byte %d", name, reader.off)
	}
	value := binary.LittleEndian.Uint32(reader.data[reader.off : reader.off+4])
	reader.off += 4
	return value, nil
}

func (reader *exceptionReader) bytes(length uint32, name string) ([]byte, error) {
	if uint64(length) > uint64(reader.remaining()) {
		return nil, fmt.Errorf("truncated %s of %d bytes at byte %d", name, length, reader.off)
	}
	if length == 0 {
		return nil, errors.New("empty byte field")
	}
	start := reader.off
	reader.off += int(length)
	return append([]byte(nil), reader.data[start:reader.off]...), nil
}
