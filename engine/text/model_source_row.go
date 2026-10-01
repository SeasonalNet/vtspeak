package text

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	paul2013SourceRowCountOffset = 0
	paul2013SourceRowTailOffset  = 4
	paul2013SourceRowsOffset     = 0x14
	paul2013SourceRowStride      = 0x94
	paul2013SourceRowLimit       = 100
	paul2013SourceRowTextOffset  = 0x34
	paul2013SourceRowAuxOffset   = 0x52
)

// Paul2013ModelSourceRowInput contains the raw arguments written by
// FUN_10044fd0. The two coordinate values, discriminator bytes, and low-byte
// flag retain their native roles without assigning higher-level meanings.
type Paul2013ModelSourceRowInput struct {
	First  uint32
	Second uint32
	Type   byte
	Class  byte
	Flag   uint32
	Text   []byte
	// AuxiliaryText selects the FUN_10045070 variant and is copied at row +0x52.
	// Nil keeps the FUN_10044fd0 single-string behavior.
	AuxiliaryText []byte
}

// AppendPaul2013ModelSourceRow ports the successful append path in
// FUN_10044fd0. The arena has a signed-short row count at +0, the most recent
// second coordinate at +4, and 0x94-byte rows beginning at +0x14. A full row
// returns appended=false, matching the native zero result. Empty text succeeds
// without appending a row. The function preserves all bytes not written by the
// native helper, including row field +0x2c, whose producer is unresolved.
func AppendPaul2013ModelSourceRow(
	arena []byte,
	input Paul2013ModelSourceRowInput,
) ([]byte, bool, error) {
	if len(arena) < paul2013SourceRowsOffset {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need at least %#x", len(arena), paul2013SourceRowsOffset)
	}
	if int16(binary.LittleEndian.Uint16(arena[paul2013SourceRowCountOffset:])) < 0 {
		return nil, false, errors.New("Paul 2013 source-row count is negative")
	}
	count := int(binary.LittleEndian.Uint16(arena[paul2013SourceRowCountOffset:]))
	text := input.Text
	if nul := bytes.IndexByte(text, 0); nul >= 0 {
		text = text[:nul]
	}
	if len(text) == 0 {
		return append([]byte(nil), arena...), true, nil
	}
	if count >= paul2013SourceRowLimit {
		return append([]byte(nil), arena...), false, nil
	}
	maxTextLength := paul2013SourceRowStride - paul2013SourceRowTextOffset - 1
	if len(text) > maxTextLength {
		return nil, false, fmt.Errorf("source-row text has %d bytes, maximum is %d", len(text), maxTextLength)
	}
	rowStart := paul2013SourceRowsOffset + count*paul2013SourceRowStride
	rowEnd := rowStart + paul2013SourceRowStride
	if rowEnd > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row %d", len(arena), rowEnd, count)
	}
	working := append([]byte(nil), arena...)
	row := working[rowStart:rowEnd]
	binary.LittleEndian.PutUint32(row[0:4], input.First)
	binary.LittleEndian.PutUint32(row[4:8], input.Second)
	binary.LittleEndian.PutUint32(row[8:12], uint32(len(text)))
	row[0x23] = input.Type
	row[0x24] = input.Class
	binary.LittleEndian.PutUint32(row[0x28:0x2c], input.Flag&0xff)
	row[0x30] = 0xff
	copy(row[paul2013SourceRowTextOffset:], text)
	row[paul2013SourceRowTextOffset+len(text)] = 0
	if input.AuxiliaryText != nil {
		auxiliaryText := input.AuxiliaryText
		if nul := bytes.IndexByte(auxiliaryText, 0); nul >= 0 {
			auxiliaryText = auxiliaryText[:nul]
		}
		if len(text) >= paul2013SourceRowAuxOffset-paul2013SourceRowTextOffset {
			return nil, false, fmt.Errorf("source-row text has %d bytes and overlaps auxiliary text at +%#x", len(text), paul2013SourceRowAuxOffset)
		}
		maxAuxiliaryLength := paul2013SourceRowStride - paul2013SourceRowAuxOffset - 1
		if len(auxiliaryText) > maxAuxiliaryLength {
			return nil, false, fmt.Errorf("source-row auxiliary text has %d bytes, maximum is %d", len(auxiliaryText), maxAuxiliaryLength)
		}
		copy(row[paul2013SourceRowAuxOffset:], auxiliaryText)
		row[paul2013SourceRowAuxOffset+len(auxiliaryText)] = 0
	}
	binary.LittleEndian.PutUint16(working[paul2013SourceRowCountOffset:], uint16(count+1))
	binary.LittleEndian.PutUint32(working[paul2013SourceRowTailOffset:], input.Second)
	return working, true, nil
}

// AppendPaul2013ModelSourceRowsSplit ports FUN_100451e0's optional split on
// spaces and tabs. A string without either delimiter is passed directly to
// the row writer. Split components retain the same coordinate/type inputs;
// empty components are ignored by the native row writer. The native helper
// uses a 32-byte local component buffer, so this function rejects components
// longer than 31 bytes instead of reproducing its possible stack overwrite.
func AppendPaul2013ModelSourceRowsSplit(
	arena []byte,
	input Paul2013ModelSourceRowInput,
) ([]byte, bool, error) {
	input.AuxiliaryText = nil
	text := input.Text
	if nul := bytes.IndexByte(text, 0); nul >= 0 {
		text = text[:nul]
	}
	if !bytes.ContainsAny(text, " \t") {
		input.Text = text
		return AppendPaul2013ModelSourceRow(arena, input)
	}
	working := append([]byte(nil), arena...)
	for start := 0; start < len(text); {
		for start < len(text) && (text[start] == ' ' || text[start] == '\t') {
			start++
		}
		end := start
		for end < len(text) && text[end] != ' ' && text[end] != '\t' {
			end++
		}
		if end == start {
			continue
		}
		if end-start > 31 {
			return nil, false, fmt.Errorf("split source-row component has %d bytes, native local capacity is 31", end-start)
		}
		part := input
		part.Text = text[start:end]
		var appended bool
		var err error
		working, appended, err = AppendPaul2013ModelSourceRow(working, part)
		if err != nil || !appended {
			return working, appended, err
		}
		start = end
	}
	return working, true, nil
}

// AppendPaul2013ModelSourceRowsBackslash ports FUN_100452c0. It splits only
// AuxiliaryText on backslashes and appends one row for every component,
// including empty components, through the FUN_10045070 row variant. Its
// 68-byte native local component buffer limits each component to 67 bytes.
func AppendPaul2013ModelSourceRowsBackslash(
	arena []byte,
	input Paul2013ModelSourceRowInput,
) ([]byte, bool, error) {
	auxiliary := input.AuxiliaryText
	if auxiliary == nil {
		auxiliary = []byte{}
	}
	if nul := bytes.IndexByte(auxiliary, 0); nul >= 0 {
		auxiliary = auxiliary[:nul]
	}
	if !bytes.Contains(auxiliary, []byte{'\\'}) {
		input.AuxiliaryText = auxiliary
		return AppendPaul2013ModelSourceRow(arena, input)
	}
	working := append([]byte(nil), arena...)
	start := 0
	for index := 0; index <= len(auxiliary); index++ {
		if index != len(auxiliary) && auxiliary[index] != '\\' {
			continue
		}
		component := auxiliary[start:index]
		if len(component) > 67 {
			return nil, false, fmt.Errorf("backslash source-row component has %d bytes, native local capacity is 67", len(component))
		}
		part := input
		part.AuxiliaryText = component
		var appended bool
		var err error
		working, appended, err = AppendPaul2013ModelSourceRow(working, part)
		if err != nil || !appended {
			return working, appended, err
		}
		start = index + 1
	}
	return working, true, nil
}

// Paul2013ParserTerminalTypeGate carries the predicates checked by
// FUN_10051a00 before it writes the preceding source row's +0x2c dword.
// PreviousRowAccepted represents the zero result from FUN_10062f10 or the
// accepted 0x50 result from FUN_10051cc0. Use
// Paul2013ParserTerminalGateFromPreviousRow to derive the former from a source
// row; scanner status and the latter producer remain explicit.
type Paul2013ParserTerminalTypeGate struct {
	ScannerStatus             int16
	RecognizedPunctuationPath bool
	PreviousRowAccepted       bool
	Punctuation               byte
}

// Paul2013ParserPreviousRowTerminalPunctuation ports FUN_10062f10's observed
// predicate. For a nonempty byte sequence, it returns true when the byte at
// row[index-1] is in the native NUL-terminated set ".?!;". The native helper
// does not check the index against the string length; this Go version returns
// false for an index outside the supplied slice.
func Paul2013ParserPreviousRowTerminalPunctuation(row []byte, index int) bool {
	if len(row) == 0 || row[0] == 0 || index < 1 || index > len(row) {
		return false
	}
	value := row[index-1]
	if value == 0 {
		return false
	}
	return value == '.' || value == '?' || value == '!' || value == ';'
}

// Paul2013ParserTerminalGateFromPreviousRow derives FUN_10051a00's
// PreviousRowAccepted input from the FUN_10062f10 predicate. That native
// branch proceeds when the predicate returns zero, so a row ending in one of
// ".?!;" leaves PreviousRowAccepted false. Scanner status and punctuation
// path recognition remain caller-supplied because their producers are not
// derived by this helper.
func Paul2013ParserTerminalGateFromPreviousRow(
	scannerStatus int16,
	recognizedPunctuationPath bool,
	previousRow []byte,
	previousRowIndex int,
	punctuation byte,
) Paul2013ParserTerminalTypeGate {
	return Paul2013ParserTerminalTypeGate{
		ScannerStatus:             scannerStatus,
		RecognizedPunctuationPath: recognizedPunctuationPath,
		PreviousRowAccepted:       !Paul2013ParserPreviousRowTerminalPunctuation(previousRow, previousRowIndex),
		Punctuation:               punctuation,
	}
}

// Paul2013ParserTerminalGateFromWhitespacePrefix composes the recovered
// two-line-feed scanner exit with FUN_10062f10's previous-row predicate. The
// ordinary status-zero path and the alternate FUN_10051cc0 path are not
// inferred by this constructor.
func Paul2013ParserTerminalGateFromWhitespacePrefix(
	prefix Paul2013ModelParserWhitespacePrefix,
	previousRow []byte,
	previousRowIndex int,
	punctuation byte,
) Paul2013ParserTerminalTypeGate {
	status := int16(0)
	if prefix.ScannerStatus == 8 {
		status = 8
	}
	return Paul2013ParserTerminalGateFromPreviousRow(
		status, false, previousRow, previousRowIndex, punctuation,
	)
}

// ApplyPaul2013ParserTerminalRowType ports FUN_10051a00's observed terminal
// type writes. When a qualifying scanner path is selected and the final row's
// +0x2c dword is zero, question mark writes 3, exclamation mark writes 4, and
// other characters write 2. The native scanner and previous-row eligibility
// producers are represented by the explicit gate.
func ApplyPaul2013ParserTerminalRowType(
	arena []byte,
	gate Paul2013ParserTerminalTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || (gate.ScannerStatus != 8 && !gate.RecognizedPunctuationPath) || !gate.PreviousRowAccepted {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	if binary.LittleEndian.Uint32(arena[field:]) != 0 {
		return append([]byte(nil), arena...), false, nil
	}
	value := uint32(2)
	if gate.Punctuation == '?' {
		value = 3
	} else if gate.Punctuation == '!' {
		value = 4
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], value)
	return working, true, nil
}

// Paul2013ParserCommaTypeGate carries the values used by the comma branch in
// FUN_100544f0. FollowingToken is compared under the DLL's mapped character
// table against "too" and "either". Scanner status, advance, and row count
// remain explicit outputs of the source scanner.
type Paul2013ParserCommaTypeGate struct {
	ScannerStatus         int16
	Character             byte
	CharacterCount        int16
	FollowingTokenAdvance int32
	FollowingTokenStatus  int32
	FollowingToken        []byte
	FollowingRowCount     int32
}

// ApplyPaul2013ParserCommaRowType ports FUN_100544f0's direct comma writes.
// It sets the last row's +0x2c dword to 5 when the following token differs
// from both "too" and "either" under FUN_1001c2c0, or when the following
// parser row count exceeds one; otherwise the
// qualifying path writes 12. The helper returns false when its scanner gates
// do not match. Other FUN_100544f0 type writes remain separate branches.
func ApplyPaul2013ParserCommaRowType(
	arena []byte,
	gate Paul2013ParserCommaTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || gate.ScannerStatus != 3 || gate.Character != ',' || gate.CharacterCount != 1 ||
		gate.FollowingTokenAdvance == 0 || gate.FollowingTokenStatus == 9 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	weights := paul2013ContextCharacterWeights()
	followingDiffersFromToo := ComparePaul2013MappedCString(gate.FollowingToken, []byte("too"), weights) != 0
	followingDiffersFromEither := ComparePaul2013MappedCString(gate.FollowingToken, []byte("either"), weights) != 0
	value := uint32(12)
	if (followingDiffersFromToo && followingDiffersFromEither) || gate.FollowingRowCount > 1 {
		value = 5
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], value)
	return working, true, nil
}

// Paul2013ParserDotTypeGate carries the scanner fields read by the dotted
// branch in FUN_100544f0. Scanner status is restricted to the two observed
// branches; a count of at least three represents the helper's native length
// gate.
type Paul2013ParserDotTypeGate struct {
	ScannerStatus  int16
	Character      byte
	CharacterCount int16
}

// ApplyPaul2013ParserDotRowType ports FUN_100544f0's statically recovered
// dot branch. When scanner status is 3 or 6, the character is '.', its count
// is at least three, and a source row exists, the last row's +0x2c dword is
// set to 1. This branch has not been directly observed in the runtime traces.
func ApplyPaul2013ParserDotRowType(
	arena []byte,
	gate Paul2013ParserDotTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || (gate.ScannerStatus != 3 && gate.ScannerStatus != 6) || gate.Character != '.' || gate.CharacterCount < 3 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserQuoteTypeGate carries the decompiler-visible inputs to the
// quote branch in FUN_100544f0. ScannerPathAccepted represents preceding
// scanner/table checks that are not reconstructed here; LookupIndex is the
// native sorted-table result, and RowCount is the local count checked < 2.
type Paul2013ParserQuoteTypeGate struct {
	ScannerStatus       int16
	Character           byte
	ScannerPathAccepted bool
	LookupIndex         int32
	RowCount            int32
}

// ApplyPaul2013ParserQuoteRowType ports FUN_100544f0's quote-path assignment.
// For accepted status 3/6 paths, a double quote, a nonnegative lookup result,
// a local row count below two, and at least one source row write 1 to the last
// source row's +0x2c dword. The preceding scanner path remains explicit.
func ApplyPaul2013ParserQuoteRowType(
	arena []byte,
	gate Paul2013ParserQuoteTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || (gate.ScannerStatus != 3 && gate.ScannerStatus != 6) ||
		gate.Character != '"' || !gate.ScannerPathAccepted || gate.LookupIndex < 0 || gate.RowCount >= 2 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserMultiSourceRowTypeGate carries the recovered post-append
// predicates in FUN_10054050. The selected generic row append must have
// succeeded, and the original source-row count must exceed one.
type Paul2013ParserMultiSourceRowTypeGate struct {
	RowAppendSucceeded bool
	InputRowCount      int32
}

// ApplyPaul2013ParserMultiSourceRowType ports FUN_10054050's direct write of
// 10 to the newly appended source row's +0x2c field. The native helper checks
// the input source-row count after a successful FUN_100451e0 append; append
// path selection remains caller supplied.
func ApplyPaul2013ParserMultiSourceRowType(
	arena []byte,
	gate Paul2013ParserMultiSourceRowTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.RowAppendSucceeded || gate.InputRowCount <= 1 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 10)
	return working, true, nil
}

// Paul2013ParserUnmatchedTypeGate carries the final predicates in the
// lookup-miss branch of FUN_100544f0. The scanner path and local row count
// remain explicit; LookupIndex is the result for the table at 0x100780fc.
type Paul2013ParserUnmatchedTypeGate struct {
	ScannerPathAccepted bool
	LookupIndex         int32
	RowCount            int32
}

// ApplyPaul2013ParserUnmatchedRowType ports FUN_100544f0's static lookup-miss
// assignment of 1. A missing table key, local count below two, accepted
// scanner path, and at least one source row write the last row's +0x2c dword.
func ApplyPaul2013ParserUnmatchedRowType(
	arena []byte,
	gate Paul2013ParserUnmatchedTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted || gate.LookupIndex > -1 || gate.RowCount >= 2 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserCloseDelimiterTypeGate carries the final checks for the
// closing-delimiter path in FUN_100544f0. Earlier token and scanner decisions
// remain explicit through ScannerPathAccepted.
type Paul2013ParserCloseDelimiterTypeGate struct {
	ScannerPathAccepted bool
	Character           byte
	ScannerStatus       int32
}

// ApplyPaul2013ParserCloseDelimiterRowType ports the static type-1 write for
// a closing parenthesis or bracket. The preceding scanner path must be
// accepted, its final status must be 1, 2, or 3, and a source row must exist.
func ApplyPaul2013ParserCloseDelimiterRowType(
	arena []byte,
	gate Paul2013ParserCloseDelimiterTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted ||
		(gate.Character != ')' && gate.Character != ']') ||
		gate.ScannerStatus < 1 || gate.ScannerStatus > 3 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserModeOneTypeGate carries the final count and mode checks in
// one FUN_100544f0 type-1 branch. The preceding scanner path is explicit.
type Paul2013ParserModeOneTypeGate struct {
	ScannerPathAccepted bool
	PriorScanCount      int32
	Mode                int32
}

// ApplyPaul2013ParserModeOneRowType ports the static write at 0x10055a8e.
// It writes 1 when the scanner path is accepted, its prior scan count is
// positive, mode is one, and a source row exists.
func ApplyPaul2013ParserModeOneRowType(
	arena []byte,
	gate Paul2013ParserModeOneTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted || gate.PriorScanCount <= 0 || gate.Mode != 1 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserAcceptedPathTypeGate represents a scanner branch whose
// internal predicates have not been recovered, but whose final source-row
// count check and type write are direct in FUN_100544f0.
type Paul2013ParserAcceptedPathTypeGate struct {
	ScannerPathAccepted bool
}

// ApplyPaul2013ParserAcceptedPathRowType ports the final count check and
// type-1 assignment at 0x10055a4a. The branch's earlier scanner/table
// predicates remain caller supplied.
func ApplyPaul2013ParserAcceptedPathRowType(
	arena []byte,
	gate Paul2013ParserAcceptedPathTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserBoundedScanTypeGate preserves the raw comparison values read
// by the remaining type-1 branch in FUN_100544f0. Their semantic roles are
// unresolved; the scanner path is supplied explicitly.
type Paul2013ParserBoundedScanTypeGate struct {
	ScannerPathAccepted bool
	ComparedValue       int32
	LimitValue          int32
}

// ApplyPaul2013ParserBoundedScanRowType ports the final predicates and type-1
// write at 0x10055677. It requires a nonzero signed value no greater than the
// supplied limit, an accepted scanner path, and at least one source row.
func ApplyPaul2013ParserBoundedScanRowType(
	arena []byte,
	gate Paul2013ParserBoundedScanTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted || gate.ComparedValue == 0 || gate.ComparedValue > gate.LimitValue {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserPositiveScanDotTypeGate carries the final status and
// character checks from a separate FUN_100544f0 dot branch.
type Paul2013ParserPositiveScanDotTypeGate struct {
	ScannerStatus int32
	Character     byte
}

// ApplyPaul2013ParserPositiveScanDotRowType ports the static type-1 write at
// 0x10056077. A positive scanner status, dot character, and nonempty source
// row list assign 1 to the last source row's +0x2c dword.
func ApplyPaul2013ParserPositiveScanDotRowType(
	arena []byte,
	gate Paul2013ParserPositiveScanDotTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || gate.ScannerStatus <= 0 || gate.Character != '.' {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserScanFallbackRowTypeGate carries the final predicates before
// FUN_100544f0's byte-indexed type-1 store at 0x1005551d. RowIndex is the
// native one-based index: the target is row RowIndex-1. The scanner and
// preceding lookup results remain explicit because their producers are not
// reconstructed here.
type Paul2013ParserScanFallbackRowTypeGate struct {
	RowIndex             int
	ScannerState         int32
	ScannerStateTwoMatch bool
	ParameterZero        int32
	PositiveScanLength   int32
	ParserMode           int32
	CandidateCount       int32
	SourceRowMatched     bool
	ZeroFlag             bool
	PositiveScanAccepted bool
}

// ApplyPaul2013ParserScanFallbackRowType ports the final type-1 assignment at
// 0x1005551d. The native branch requires a source row and matching row key,
// scanner state 1 or accepted state 2, parameter zero below two, positive scan
// length, parser mode 1 through 3, candidate count below two, a zero local
// flag, and a positive preceding scan result. The native write is
// unconditional once these predicates pass; the helper bounds the one-based
// row index to the supplied arena.
func ApplyPaul2013ParserScanFallbackRowType(
	arena []byte,
	gate Paul2013ParserScanFallbackRowTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || gate.RowIndex < 1 || gate.RowIndex > count ||
		(gate.ScannerState != 1 && gate.ScannerState != 2) ||
		(gate.ScannerState == 2 && !gate.ScannerStateTwoMatch) ||
		gate.ParameterZero >= 2 || gate.PositiveScanLength <= 0 ||
		gate.ParserMode < 1 || gate.ParserMode > 3 || gate.CandidateCount >= 2 ||
		!gate.SourceRowMatched || gate.ZeroFlag || !gate.PositiveScanAccepted {
		return append([]byte(nil), arena...), false, nil
	}
	field := paul2013SourceRowsOffset + (gate.RowIndex-1)*paul2013SourceRowStride + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserHyphenFallbackTypeGate carries the terminal predicates for a
// type-7 write in the separate sign-scanner branch. Scanner classification
// before the hyphen decision remains caller-owned.
type Paul2013ParserHyphenFallbackTypeGate struct {
	ScannerPathAccepted bool
	Character           byte
	CharacterCount      int32
	Coordinate          int32
}

// ApplyPaul2013ParserHyphenFallbackRowType ports the static type-7 assignment
// at 0x1004eb4d. An accepted one-character hyphen path with a negative
// coordinate sentinel writes 7 to the last source row's +0x2c dword.
func ApplyPaul2013ParserHyphenFallbackRowType(
	arena []byte,
	gate Paul2013ParserHyphenFallbackTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted || gate.Character != '-' ||
		gate.CharacterCount != 1 || gate.Coordinate > -1 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 7)
	return working, true, nil
}

// Paul2013ParserQuoteCoordinateTypeGate carries the repeated quote fallback
// conditions found in several source-scanner routines. Earlier scanner path
// selection remains explicit.
type Paul2013ParserQuoteCoordinateTypeGate struct {
	ScannerPathAccepted bool
	Character           byte
	CharacterCount      int32
	Coordinate          int32
}

// ApplyPaul2013ParserQuoteCoordinateRowType ports the repeated static type-1
// fallback: a one-character double quote with a coordinate sentinel <= -1
// marks the last source row. The gate is shared by the disassembly paths at
// 0x1004f426, 0x1004f867, 0x1004fbb4, 0x1004f6c5, 0x1004fa17, 0x10050073,
// 0x10050229, 0x10050278, 0x1005084b, 0x10050b4e, and 0x100510df where the
// surrounding scanner state establishes the same input predicates.
func ApplyPaul2013ParserQuoteCoordinateRowType(
	arena []byte,
	gate Paul2013ParserQuoteCoordinateTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.ScannerPathAccepted || gate.Character != '"' ||
		gate.CharacterCount != 1 || gate.Coordinate > -1 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserInputFlagTypeGate carries the append result and short input
// flag read by the direct type-5 producer at 0x100445a1.
type Paul2013ParserInputFlagTypeGate struct {
	RowAppendSucceeded bool
	InputFlag          int16
}

// ApplyPaul2013ParserInputFlagRowType ports the static type-5 write at
// 0x100445b1. It marks the newly appended row when the input record's short
// at +2 is nonzero; row selection and append success remain explicit.
func ApplyPaul2013ParserInputFlagRowType(
	arena []byte,
	gate Paul2013ParserInputFlagTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.RowAppendSucceeded || gate.InputFlag == 0 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 5)
	return working, true, nil
}

// Paul2013ParserAppendType7Gate carries the successful row-append result and
// whether the native caller still has an input row to process.
type Paul2013ParserAppendType7Gate struct {
	RowAppendSucceeded bool
	MoreInputRows      bool
}

// ApplyPaul2013ParserAcceptedAppendRowType7 ports the direct type-7 write at
// 0x10057722 immediately after a successful A-row append.
func ApplyPaul2013ParserAcceptedAppendRowType7(
	arena []byte,
	gate Paul2013ParserAppendType7Gate,
) ([]byte, bool, error) {
	return applyPaul2013ParserAppendType7(arena, gate, false)
}

// ApplyPaul2013ParserIntermediateAppendRowType7 ports the type-7 write at
// 0x100577db, which occurs after a successful append only when another input
// row remains in the native loop.
func ApplyPaul2013ParserIntermediateAppendRowType7(
	arena []byte,
	gate Paul2013ParserAppendType7Gate,
) ([]byte, bool, error) {
	return applyPaul2013ParserAppendType7(arena, gate, true)
}

func applyPaul2013ParserAppendType7(
	arena []byte,
	gate Paul2013ParserAppendType7Gate,
	requireMoreRows bool,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.RowAppendSucceeded || (requireMoreRows && !gate.MoreInputRows) {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 7)
	return working, true, nil
}

// Paul2013ParserZeroInputFlagTypeGate carries the shared post-append
// predicates for two static type-8 producers.
type Paul2013ParserZeroInputFlagTypeGate struct {
	RowAppendSucceeded bool
	InputFlag          int16
}

// ApplyPaul2013ParserZeroInputFlagRowType ports the type-8 stores at
// 0x10053229 and 0x100538c0. A successful append and zero input short at +2
// assign 8 to the last source row; earlier recognition gates remain explicit.
func ApplyPaul2013ParserZeroInputFlagRowType(
	arena []byte,
	gate Paul2013ParserZeroInputFlagTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.RowAppendSucceeded || gate.InputFlag != 0 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 8)
	return working, true, nil
}

// Paul2013ParserHyphenModeTypeGate carries the direct post-append checks at
// 0x10045a7c. Mode retains its opaque native value.
type Paul2013ParserHyphenModeTypeGate struct {
	RowAppendSucceeded bool
	Mode               int16
	FollowingByte      byte
}

// ApplyPaul2013ParserHyphenModeRowType ports the static type-11 assignment
// at 0x10045a90. After a successful append, mode 2 and a following hyphen
// write 11 to the last source row's +0x2c field.
func ApplyPaul2013ParserHyphenModeRowType(
	arena []byte,
	gate Paul2013ParserHyphenModeTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.RowAppendSucceeded || gate.Mode != 2 || gate.FollowingByte != '-' {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 11)
	return working, true, nil
}

// Paul2013ParserExactInputFlagTypeGate carries the exact flag comparison and
// append result for the second post-append type-5 producer.
type Paul2013ParserExactInputFlagTypeGate struct {
	RowAppendSucceeded bool
	InputFlag          int16
}

// ApplyPaul2013ParserExactInputFlagRowType ports the type-5 store at
// 0x10059f8b, which requires the input short at +2 to equal exactly one.
func ApplyPaul2013ParserExactInputFlagRowType(
	arena []byte,
	gate Paul2013ParserExactInputFlagTypeGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.RowAppendSucceeded || gate.InputFlag != 1 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 5)
	return working, true, nil
}

// Paul2013ParserContextualTypeOneGate preserves the raw final conditions at
// 0x10042e8f. The integer fields are deliberately unlabeled native values.
type Paul2013ParserContextualTypeOneGate struct {
	ScannerPathAccepted     bool
	PriorRowCount           int16
	PriorRowKeyMatches      bool
	PredicateValue          int32
	ScanValue               int32
	ModeValue               int32
	SecondaryPredicateValue int32
}

// Paul2013ParserModeOneTypeOneGate retains the raw predicates for the
// assignment at 0x100430ce. Its native state values are supplied as booleans
// without inferred meanings.
type Paul2013ParserModeOneTypeOneGate struct {
	PriorRowKeyMatches bool
	ModeValue          int32
	PositiveStateA     bool
	PositiveStateB     bool
}

// ApplyPaul2013ParserModeOnePriorKeyRowTypeOne ports the static assignment at
// 0x100430ce. It writes the row selected by the prior-row count only when the
// prior key matches and the native mode/state predicates pass.
func ApplyPaul2013ParserModeOnePriorKeyRowTypeOne(
	arena []byte,
	gate Paul2013ParserModeOneTypeOneGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !gate.PriorRowKeyMatches || gate.ModeValue != 1 ||
		!gate.PositiveStateA || !gate.PositiveStateB {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserTypeSevenPreAppendGate contains the two positive caller-state
// checks observed before the 0x1005782b type-7 write.
type Paul2013ParserTypeSevenPreAppendGate struct {
	StateA int32
	StateB int32
}

// ApplyPaul2013ParserAcceptedScannerRowTypeOne ports the direct type-1 write
// at 0x1003d4b6 after FUN_1005a350 returns 1. The earlier scan inputs remain
// represented by the raw scanner result.
func ApplyPaul2013ParserAcceptedScannerRowTypeOne(
	arena []byte,
	scannerResult int32,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || scannerResult != 1 {
		return append([]byte(nil), arena...), false, nil
	}
	field := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride + 0x2c
	if binary.LittleEndian.Uint32(arena[field:]) != 0 {
		return append([]byte(nil), arena...), false, nil
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

// Paul2013ParserTypeTwoResetGate retains the raw peer-row comparisons before
// the static type-2-to-zero write at 0x1003da7a. The two +0x28 words and
// +0x22 sentinels have no semantic labels here.
type Paul2013ParserTypeTwoResetGate struct {
	RowIndex               int
	FirstPeerWordAtPlus28  uint32
	SecondPeerWordAtPlus28 uint32
	CurrentWordAtPlus22    uint16
	PeerWordAtPlus22       uint16
}

// ApplyPaul2013ParserTypeTwoReset ports the write of zero at 0x1003da7a when
// the selected row is type 2, both compared peer words equal 0x12, and at
// least one compared +0x22 word differs from 0x1f.
func ApplyPaul2013ParserTypeTwoReset(
	arena []byte,
	gate Paul2013ParserTypeTwoResetGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if gate.RowIndex < 0 || gate.RowIndex >= count {
		return nil, false, fmt.Errorf("source-row type reset index %d outside %d rows", gate.RowIndex, count)
	}
	if gate.FirstPeerWordAtPlus28 != 0x12 || gate.SecondPeerWordAtPlus28 != 0x12 ||
		(gate.CurrentWordAtPlus22 == 0x1f && gate.PeerWordAtPlus22 == 0x1f) {
		return append([]byte(nil), arena...), false, nil
	}
	field := paul2013SourceRowsOffset + gate.RowIndex*paul2013SourceRowStride + 0x2c
	if binary.LittleEndian.Uint32(arena[field:]) != 2 {
		return append([]byte(nil), arena...), false, nil
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 0)
	return working, true, nil
}

// ApplyPaul2013ParserPreAppendRowType7 ports the static write at 0x1005782b,
// which marks the next source row before FUN_10044fd0 appends it. It requires
// capacity for the pending row and fails closed when the source-row limit is
// already reached.
func ApplyPaul2013ParserPreAppendRowType7(
	arena []byte,
	gate Paul2013ParserTypeSevenPreAppendGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if gate.StateA <= 0 || gate.StateB <= 0 || count >= paul2013SourceRowLimit {
		return append([]byte(nil), arena...), false, nil
	}
	field := paul2013SourceRowsOffset + count*paul2013SourceRowStride + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for pending row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 7)
	return working, true, nil
}

// ApplyPaul2013ParserSplitAppendRowType7 ports the static final-row write at
// 0x1005788b after a successful space/tab split append.
func ApplyPaul2013ParserSplitAppendRowType7(
	arena []byte,
	appendSucceeded bool,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	if count == 0 || !appendSucceeded {
		return append([]byte(nil), arena...), false, nil
	}
	field := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride + 0x2c
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 7)
	return working, true, nil
}

// ApplyPaul2013ParserContextualRowTypeOne ports the static type-1 assignment
// at 0x10042e8f. It checks the preceding row/key relation and the native
// bounds/equalities explicitly supplied in the gate.
func ApplyPaul2013ParserContextualRowTypeOne(
	arena []byte,
	gate Paul2013ParserContextualTypeOneGate,
) ([]byte, bool, error) {
	count, err := paul2013SourceRowCount(arena)
	if err != nil {
		return nil, false, err
	}
	modeAccepted := gate.ModeValue == 1 || gate.ModeValue == 2
	if count == 0 || !gate.ScannerPathAccepted || gate.PriorRowCount <= 0 ||
		!gate.PriorRowKeyMatches || gate.PredicateValue != 1 ||
		gate.ScanValue < 1 || gate.ScanValue > 2 || !modeAccepted ||
		gate.SecondaryPredicateValue < 1 || gate.SecondaryPredicateValue > 2 {
		return append([]byte(nil), arena...), false, nil
	}
	rowStart := paul2013SourceRowsOffset + (count-1)*paul2013SourceRowStride
	field := rowStart + 0x2c
	if field+4 > len(arena) {
		return nil, false, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for row type", len(arena), field+4)
	}
	working := append([]byte(nil), arena...)
	binary.LittleEndian.PutUint32(working[field:], 1)
	return working, true, nil
}

func paul2013SourceRowCount(arena []byte) (int, error) {
	if len(arena) < paul2013SourceRowsOffset {
		return 0, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need at least %#x", len(arena), paul2013SourceRowsOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(arena[paul2013SourceRowCountOffset:])))
	if count < 0 || count > paul2013SourceRowLimit {
		return 0, fmt.Errorf("Paul 2013 source-row count %d is outside [0, %d]", count, paul2013SourceRowLimit)
	}
	if count > 0 {
		end := paul2013SourceRowsOffset + count*paul2013SourceRowStride
		if end > len(arena) {
			return 0, fmt.Errorf("Paul 2013 source-row arena has %d bytes, need %d for %d rows", len(arena), end, count)
		}
	}
	return count, nil
}
