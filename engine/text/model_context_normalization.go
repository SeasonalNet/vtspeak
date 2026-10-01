package text

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	paul2013SourceRowKindOffset       = 0x23
	paul2013SourceRowClassOffset      = 0x24
	paul2013ModelContextSurfaceOffset = 0x05
	paul2013ModelContextPhoneOffset   = 0x23
	paul2013ModelContextPhoneEnd      = 0x64
	paul2013ModelContextFlagBackstep  = 2
	paul2013ModelContextGateOffset    = 0x6c
)

// Paul2013GenericContextNormalization is FUN_10002f10's result as consumed by
// FUN_10009030. RowFlags is the native mask to OR into the model row's flag
// byte; the code bytes remain opaque.
type Paul2013GenericContextNormalization struct {
	Handled  bool
	Codes    []byte
	RowFlags byte
}

// Paul2013GenericContextNormalizer supplies the recovered generic-normalizer
// stage when FUN_10009030's row metadata gate allows or requires it.
type Paul2013GenericContextNormalizer func(surface []byte) (Paul2013GenericContextNormalization, error)

// Paul2013ModelContextNormalizationResult records the copied model arena and
// whether FUN_10009030 handled this model context row.
type Paul2013ModelContextNormalizationResult struct {
	Model   []byte
	Handled bool
}

// Paul2013ModelContextSurface reads the bounded C-string surface and current
// low-byte state for one counted model row.
func Paul2013ModelContextSurface(model []byte, contextRowIndex int) ([]byte, byte, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return nil, 0, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count < 0 || contextRowIndex < 0 || contextRowIndex >= count {
		return nil, 0, fmt.Errorf("model context row index %d is outside signed count %d", contextRowIndex, count)
	}
	rowStart := paul2013ModelContextRowBaseOffset + contextRowIndex*Paul2013PhoneContextRowSize
	rowEnd := rowStart + Paul2013PhoneContextRowSize
	if rowStart < paul2013ModelContextCountOffset+4 || rowEnd > len(model) {
		return nil, 0, fmt.Errorf("model context row %d range [%#x,%#x) exceeds model size %d", contextRowIndex, rowStart, rowEnd, len(model))
	}
	surfaceArea := model[rowStart+paul2013ModelContextSurfaceOffset : rowStart+paul2013ModelContextPhoneOffset]
	surfaceEnd := bytes.IndexByte(surfaceArea, 0)
	if surfaceEnd < 0 {
		return nil, 0, fmt.Errorf("model context row %d surface is not NUL-terminated", contextRowIndex)
	}
	return append([]byte(nil), surfaceArea[:surfaceEnd]...), model[rowStart-paul2013ModelContextFlagBackstep], nil
}

// ApplyPaul2013ModelContextCodes writes one already-selected native handler
// result into a counted model context row. It preserves the row and arena
// bytes outside the phone string and the low byte of the preceding state
// short, matching the write shape shared by the FUN_10009030 and C3A0 paths.
func ApplyPaul2013ModelContextCodes(
	model []byte,
	contextRowIndex int,
	codes []byte,
	rowFlags byte,
) ([]byte, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return nil, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count < 0 || contextRowIndex < 0 || contextRowIndex >= count {
		return nil, fmt.Errorf("model context row index %d is outside signed count %d", contextRowIndex, count)
	}
	rowStart := paul2013ModelContextRowBaseOffset + contextRowIndex*Paul2013PhoneContextRowSize
	rowEnd := rowStart + Paul2013PhoneContextRowSize
	if rowStart < paul2013ModelContextCountOffset+4 || rowEnd > len(model) {
		return nil, fmt.Errorf("model context row %d range [%#x,%#x) exceeds model size %d", contextRowIndex, rowStart, rowEnd, len(model))
	}
	if rowFlags & ^byte(0x39) != 0 {
		return nil, fmt.Errorf("model context row %d handler returned unsupported row flags %#02x", contextRowIndex, rowFlags)
	}
	if bytes.IndexByte(codes, 0) >= 0 {
		return nil, errors.New("model context handler codes contain an embedded NUL")
	}
	phoneCapacity := paul2013ModelContextPhoneEnd - paul2013ModelContextPhoneOffset
	if len(codes) >= phoneCapacity {
		return nil, fmt.Errorf("model context handler returned %d codes, phone field capacity is %d", len(codes), phoneCapacity-1)
	}
	working := append([]byte(nil), model...)
	copy(working[rowStart+paul2013ModelContextPhoneOffset:], codes)
	working[rowStart+paul2013ModelContextPhoneOffset+len(codes)] = 0
	working[rowStart-paul2013ModelContextFlagBackstep] |= rowFlags
	return working, nil
}

// Paul2013ModelContextExceptionGate records FUN_10007520's exception lookup
// decision for one full model context row, including FUN_100091b0's recovered
// low-short result.
type Paul2013ModelContextExceptionGate struct {
	RowIndex       int
	ParserRowIndex int
	Preparation    Paul2013PhoneContextPreparation
	Input          Paul2013ExceptionDispatchGateInput
	Dispatch       bool
}

// Paul2013ModelContextFormYResult reports the known parser-form Y writes and
// the final context-row index after native duplicate-row skipping.
type Paul2013ModelContextFormYResult struct {
	Model         []byte
	Applied       bool
	TailHandled   bool
	TailSurface   []byte
	FirstRowIndex int
	LastRowIndex  int
}

// Paul2013ModelContextFormSResult reports the initial context-string
// normalization writes made by parser form S.
type Paul2013ModelContextFormSResult struct {
	Model       []byte
	Applied     bool
	TailHandled bool
	TailSurface []byte
	RowIndex    int
}

// Paul2013ModelContextSecondPassResult reports the preceding context rows
// rewritten by the final pass in FUN_100091b0.
type Paul2013ModelContextSecondPassResult struct {
	Model       []byte
	ChangedRows []int
}

// Paul2013ModelContextProcessingFlagsResult reports the rows whose initial
// native processing state was set by FUN_10007520 before its first pass.
type Paul2013ModelContextProcessingFlagsResult struct {
	Model       []byte
	SetRows     []int
	ClearedRows []int
}

// InitializePaul2013ModelContextProcessingFlags ports FUN_10007520's initial
// counted-row loop. It sets the 16-bit state field immediately before each
// 0x70-byte model context row to 1 when the first phone byte is nonzero and to
// 0 otherwise. Later native handlers narrow updates to the low byte of this
// field, so initialization writes the full short as observed.
func InitializePaul2013ModelContextProcessingFlags(model []byte) (Paul2013ModelContextProcessingFlagsResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ModelContextProcessingFlagsResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count < 0 {
		return Paul2013ModelContextProcessingFlagsResult{}, fmt.Errorf("model context count %d is negative", count)
	}
	rowCapacity := (len(model) - paul2013ModelContextRowBaseOffset) / Paul2013PhoneContextRowSize
	if count > rowCapacity {
		return Paul2013ModelContextProcessingFlagsResult{}, fmt.Errorf("model declares %d context rows but only %d fit", count, rowCapacity)
	}
	result := Paul2013ModelContextProcessingFlagsResult{Model: append([]byte(nil), model...)}
	for rowIndex := 0; rowIndex < count; rowIndex++ {
		rowStart := paul2013ModelContextRowBaseOffset + rowIndex*Paul2013PhoneContextRowSize
		stateOffset := rowStart - paul2013ModelContextFlagBackstep
		state := uint16(0)
		if model[rowStart+paul2013ModelContextPhoneOffset] != 0 {
			state = 1
			result.SetRows = append(result.SetRows, rowIndex)
		} else {
			result.ClearedRows = append(result.ClearedRows, rowIndex)
		}
		binary.LittleEndian.PutUint16(result.Model[stateOffset:stateOffset+2], state)
	}
	return result, nil
}

// ApplyPaul2013ModelContextSecondPass ports FUN_100091b0's final counted-row
// pass. It recognizes the directly observed previous-surface "the" and
// following-surface "de" cases, validates the phone-code class, then writes
// the two phone bytes and state bit into the preceding context row. The
// comparison for the following surface follows the native pointer even for
// the final counted row, provided that its full surface field is present in
// the arena.
func ApplyPaul2013ModelContextSecondPass(model []byte) (Paul2013ModelContextSecondPassResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ModelContextSecondPassResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count < 0 {
		return Paul2013ModelContextSecondPassResult{}, fmt.Errorf("model context count %d is negative", count)
	}
	rowBytes := len(model) - paul2013ModelContextRowBaseOffset
	if count > rowBytes/Paul2013PhoneContextRowSize {
		return Paul2013ModelContextSecondPassResult{}, fmt.Errorf("model declares %d context rows but only %d fit", count, rowBytes/Paul2013PhoneContextRowSize)
	}
	result := Paul2013ModelContextSecondPassResult{Model: append([]byte(nil), model...)}
	const (
		followingSurfaceOffset = 0x7b
		previousSurfaceOffset  = -0x65
	)
	shortBase := paul2013ModelContextCountOffset
	for rowIndex := 1; rowIndex < count; rowIndex++ {
		byteBase := shortBase + rowIndex*Paul2013PhoneContextRowSize
		phoneOffset := byteBase + 0x29
		phoneCode := model[phoneOffset]
		if phoneCode == 0 || phoneCode >= 'F' || !paul2013ContextCodeClass(phoneCode) || phoneCode == 'C' {
			continue
		}
		previousMatch, err := paul2013ModelCStringEquals(model, byteBase+previousSurfaceOffset, paul2013PhoneContextSurfaceCapacity, "the")
		if err != nil {
			return Paul2013ModelContextSecondPassResult{}, fmt.Errorf("read previous surface for context row %d: %w", rowIndex, err)
		}
		followingMatch := false
		if !previousMatch {
			followingMatch, err = paul2013ModelCStringEquals(model, byteBase+followingSurfaceOffset, paul2013PhoneContextSurfaceCapacity, "de")
			if err != nil {
				return Paul2013ModelContextSecondPassResult{}, fmt.Errorf("read following surface for context row %d: %w", rowIndex, err)
			}
		}
		if !previousMatch && !followingMatch {
			continue
		}
		previousRowStart := paul2013ModelContextRowBaseOffset + (rowIndex-1)*Paul2013PhoneContextRowSize
		phoneWriteOffset := previousRowStart + paul2013ModelContextPhoneOffset + 1
		flagOffset := previousRowStart - paul2013ModelContextFlagBackstep
		if phoneWriteOffset < 0 || phoneWriteOffset+2 > len(result.Model) || flagOffset < 0 || flagOffset >= len(result.Model) {
			return Paul2013ModelContextSecondPassResult{}, fmt.Errorf("context row %d second-pass writes exceed model size %d", rowIndex, len(model))
		}
		result.Model[phoneWriteOffset] = 0x26
		result.Model[phoneWriteOffset+1] = 0
		result.Model[flagOffset] |= 8
		result.ChangedRows = append(result.ChangedRows, rowIndex-1)
	}
	return result, nil
}

func paul2013ModelCStringEquals(model []byte, offset, capacity int, expected string) (bool, error) {
	if offset < 0 || offset > len(model)-capacity {
		return false, fmt.Errorf("surface at offset %#x with capacity %#x exceeds model size %d", offset, capacity, len(model))
	}
	field := model[offset : offset+capacity]
	nul := bytes.IndexByte(field, 0)
	if nul < 0 {
		return false, fmt.Errorf("surface at offset %#x has no NUL within %#x bytes", offset, capacity)
	}
	return string(field[:nul]) == expected, nil
}

// DerivePaul2013ModelContextPreparation ports FUN_100091b0's known low-short
// return branches against the full FUN_1000ea20 model arena. It does not apply
// the function's still-unrecovered context-row mutations.
func DerivePaul2013ModelContextPreparation(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
) (Paul2013PhoneContextPreparation, error) {
	row, parserRow, err := paul2013ModelContextAndParserRows(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013PhoneContextPreparation{}, err
	}
	parserIndex := int(int16(binary.LittleEndian.Uint16(row[paul2013ModelContextTokenIndex:])))
	earlyReturn := parserRow[paul2013SourceRowKindOffset] != 'U'
	if contextRowIndex > 0 {
		previousStart := paul2013ModelContextRowBaseOffset + (contextRowIndex-1)*Paul2013PhoneContextRowSize
		previousIndex := int16(binary.LittleEndian.Uint16(model[previousStart+paul2013ModelContextTokenIndex:]))
		earlyReturn = earlyReturn || previousIndex == int16(parserIndex)
	}
	if earlyReturn {
		return Paul2013PhoneContextPreparation{Code: 0, Known: true, EarlyReturn: true}, nil
	}
	form := parserRow[paul2013SourceRowClassOffset]
	if form == 'Y' || form == 'S' {
		return Paul2013PhoneContextPreparation{Code: 1, Known: true}, nil
	}
	return Paul2013PhoneContextPreparation{Code: 0, Known: true}, nil
}

// ApplyPaul2013ModelContextFormY ports FUN_100091b0's parser-form Y copy,
// flag-bit update, and adjacent duplicate-token row advance. The later
// auxiliary-string and suffix-handler cascade remains a separate stage.
func ApplyPaul2013ModelContextFormY(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
) (Paul2013ModelContextFormYResult, error) {
	return ApplyPaul2013ModelContextFormYWithTail(model, parserRows, contextRowIndex, nil)
}

// ApplyPaul2013ModelContextFormYWithTail uses an explicit parser tail override
// when non-nil; otherwise it derives the tail from the parser row selected
// after native duplicate-row advancement.
func ApplyPaul2013ModelContextFormYWithTail(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
	tailSurface []byte,
) (Paul2013ModelContextFormYResult, error) {
	row, parserRow, err := paul2013ModelContextAndParserRows(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ModelContextFormYResult{}, err
	}
	result := Paul2013ModelContextFormYResult{
		Model: append([]byte(nil), model...), FirstRowIndex: contextRowIndex,
		LastRowIndex: contextRowIndex,
	}
	preparation, err := DerivePaul2013ModelContextPreparation(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ModelContextFormYResult{}, err
	}
	if parserRow[paul2013SourceRowClassOffset] != 'Y' || preparation.Code != 1 {
		return result, nil
	}
	auxiliaryArea := parserRow[0x52:]
	auxiliaryEnd := bytes.IndexByte(auxiliaryArea, 0)
	if auxiliaryEnd < 0 {
		return Paul2013ModelContextFormYResult{}, fmt.Errorf("parser row %d form-Y text is not NUL-terminated", int(int16(binary.LittleEndian.Uint16(row[paul2013ModelContextTokenIndex:]))))
	}
	if auxiliaryEnd >= paul2013ModelContextPhoneEnd-paul2013ModelContextPhoneOffset {
		return Paul2013ModelContextFormYResult{}, fmt.Errorf("parser form-Y text has %d bytes, model phone field capacity is %d", auxiliaryEnd, paul2013ModelContextPhoneEnd-paul2013ModelContextPhoneOffset-1)
	}
	rowStart := paul2013ModelContextRowBaseOffset + contextRowIndex*Paul2013PhoneContextRowSize
	copy(result.Model[rowStart+paul2013ModelContextPhoneOffset:], auxiliaryArea[:auxiliaryEnd])
	result.Model[rowStart+paul2013ModelContextPhoneOffset+auxiliaryEnd] = 0
	result.Model[rowStart-paul2013ModelContextFlagBackstep] |= 4
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count > (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize {
		return Paul2013ModelContextFormYResult{}, fmt.Errorf("model declares %d context rows but only %d fit", count, (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize)
	}
	for result.LastRowIndex+1 < count {
		currentStart := paul2013ModelContextRowBaseOffset + result.LastRowIndex*Paul2013PhoneContextRowSize
		nextStart := currentStart + Paul2013PhoneContextRowSize
		currentIndex := binary.LittleEndian.Uint16(model[currentStart+paul2013ModelContextTokenIndex:])
		nextIndex := binary.LittleEndian.Uint16(model[nextStart+paul2013ModelContextTokenIndex:])
		if currentIndex != nextIndex {
			break
		}
		result.LastRowIndex++
	}
	if tailSurface == nil {
		lastRow, _, err := paul2013ModelContextAndParserRows(
			result.Model, parserRows, result.LastRowIndex,
		)
		if err != nil {
			return Paul2013ModelContextFormYResult{}, err
		}
		parserIndex := int(int16(binary.LittleEndian.Uint16(lastRow[paul2013ModelContextTokenIndex:])))
		tailSurface, err = paul2013ParserTailSurface(parserRows, parserIndex)
		if err != nil {
			return Paul2013ModelContextFormYResult{}, fmt.Errorf("derive form-Y parser tail at context row %d: %w", result.LastRowIndex, err)
		}
	}
	result.TailSurface = append([]byte(nil), tailSurface...)
	result.TailHandled, err = applyPaul2013ModelContextLiteralTail(
		result.Model, result.LastRowIndex, tailSurface,
	)
	if err != nil {
		return Paul2013ModelContextFormYResult{}, err
	}
	result.Applied = true
	return result, nil
}

// ApplyPaul2013ModelContextFormS ports FUN_100091b0's parser-form S context
// encoding into the model phone field. The later auxiliary-string and suffix
// handler cascade remains a separate stage.
func ApplyPaul2013ModelContextFormS(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
) (Paul2013ModelContextFormSResult, error) {
	return ApplyPaul2013ModelContextFormSWithTail(model, parserRows, contextRowIndex, nil)
}

// ApplyPaul2013ModelContextFormSWithTail uses an explicit parser tail override
// when non-nil; otherwise it derives the tail from the current parser row.
func ApplyPaul2013ModelContextFormSWithTail(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
	tailSurface []byte,
) (Paul2013ModelContextFormSResult, error) {
	row, parserRow, err := paul2013ModelContextAndParserRows(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ModelContextFormSResult{}, err
	}
	result := Paul2013ModelContextFormSResult{
		Model: append([]byte(nil), model...), RowIndex: contextRowIndex,
	}
	preparation, err := DerivePaul2013ModelContextPreparation(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ModelContextFormSResult{}, err
	}
	if parserRow[paul2013SourceRowClassOffset] != 'S' || preparation.Code != 1 {
		return result, nil
	}
	if tailSurface == nil {
		parserIndex := int(int16(binary.LittleEndian.Uint16(row[paul2013ModelContextTokenIndex:])))
		tailSurface, err = paul2013ParserTailSurface(parserRows, parserIndex)
		if err != nil {
			return Paul2013ModelContextFormSResult{}, fmt.Errorf("derive form-S parser tail at context row %d: %w", contextRowIndex, err)
		}
	}
	result.TailSurface = append([]byte(nil), tailSurface...)
	surfaceArea := row[paul2013ModelContextSurfaceOffset:paul2013ModelContextPhoneOffset]
	surfaceEnd := bytes.IndexByte(surfaceArea, 0)
	if surfaceEnd < 0 {
		return Paul2013ModelContextFormSResult{}, fmt.Errorf("model context row %d surface is not NUL-terminated", contextRowIndex)
	}
	codes, err := EncodePaul2013ContextString(surfaceArea[:surfaceEnd])
	if err != nil {
		return Paul2013ModelContextFormSResult{}, fmt.Errorf("encode form-S model context row %d: %w", contextRowIndex, err)
	}
	if len(codes) >= paul2013ModelContextPhoneEnd-paul2013ModelContextPhoneOffset {
		return Paul2013ModelContextFormSResult{}, fmt.Errorf("form-S context encoding has %d bytes, phone field capacity is %d", len(codes), paul2013ModelContextPhoneEnd-paul2013ModelContextPhoneOffset-1)
	}
	rowStart := paul2013ModelContextRowBaseOffset + contextRowIndex*Paul2013PhoneContextRowSize
	copy(result.Model[rowStart+paul2013ModelContextPhoneOffset:], codes)
	result.Model[rowStart+paul2013ModelContextPhoneOffset+len(codes)] = 0
	if len(codes) != 0 {
		result.Model[rowStart-paul2013ModelContextFlagBackstep] |= 0x20
	}
	result.TailHandled, err = applyPaul2013ModelContextLiteralTail(
		result.Model, contextRowIndex, tailSurface,
	)
	if err != nil {
		return Paul2013ModelContextFormSResult{}, err
	}
	result.Applied = true
	return result, nil
}

// paul2013ParserTailSurface ports FUN_100091b0's read of the selected
// 0x94-byte parser row: a positive dword at +0x14 is an offset from +0x34,
// and the dword at +0x18 gives the bounded copy length. The native local is
// 32 bytes; larger lengths fail closed instead of reproducing a stack write.
func paul2013ParserTailSurface(parserRows []byte, parserRowIndex int) ([]byte, error) {
	if len(parserRows)%paul2013ParserSourceRowStride != 0 {
		return nil, fmt.Errorf("parser rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013ParserSourceRowStride)
	}
	rowCount := len(parserRows) / paul2013ParserSourceRowStride
	if parserRowIndex < 0 || parserRowIndex >= rowCount {
		return nil, fmt.Errorf("parser tail row index %d is outside %d rows", parserRowIndex, rowCount)
	}
	rowStart := parserRowIndex * paul2013ParserSourceRowStride
	row := parserRows[rowStart : rowStart+paul2013ParserSourceRowStride]
	sourceOffset := int32(binary.LittleEndian.Uint32(row[0x14:0x18]))
	copyLength := int32(binary.LittleEndian.Uint32(row[0x18:0x1c]))
	if sourceOffset < 1 || copyLength <= 0 {
		return nil, nil
	}
	if copyLength >= 32 {
		return nil, fmt.Errorf("parser tail length %d exceeds native 32-byte local capacity", copyLength)
	}
	sourceStart := int64(rowStart) + int64(paul2013ParserSourceRowTextOffset) + int64(sourceOffset)
	sourceEnd := sourceStart + int64(copyLength)
	if sourceStart < 0 || sourceEnd > int64(len(parserRows)) {
		return nil, fmt.Errorf("parser tail range [%#x,%#x) exceeds parser arena size %d", sourceStart, sourceEnd, len(parserRows))
	}
	surface := parserRows[int(sourceStart):int(sourceEnd)]
	if nul := bytes.IndexByte(surface, 0); nul >= 0 {
		surface = surface[:nul]
	}
	return append([]byte(nil), surface...), nil
}

func applyPaul2013ModelContextLiteralTail(model []byte, rowIndex int, surface []byte) (bool, error) {
	rowStart := paul2013ModelContextRowBaseOffset + rowIndex*Paul2013PhoneContextRowSize
	flagOffset := rowStart - paul2013ModelContextFlagBackstep
	rowEnd := flagOffset + Paul2013PhoneContextRowSize
	if flagOffset < 0 || rowEnd > len(model) {
		return false, fmt.Errorf("model context row %d literal tail exceeds model size %d", rowIndex, len(model))
	}
	_, handled, err := ApplyPaul2013LiteralContextHandler(model[flagOffset:rowEnd], surface)
	return handled, err
}

// BuildPaul2013ModelContextExceptionGate joins the recovered preparation
// result to FUN_10007520's full-model/source-row dispatch gate.
func BuildPaul2013ModelContextExceptionGate(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
	preparationCode int16,
) (Paul2013ExceptionDispatchGateInput, error) {
	row, parserRow, err := paul2013ModelContextAndParserRows(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ExceptionDispatchGateInput{}, err
	}
	return Paul2013ExceptionDispatchGateInput{
		SourceTokenStatus:      int8(parserRow[0x30]),
		PhoneRowMarker:         row[paul2013ModelContextPhoneOffset],
		SourceKind:             parserRow[paul2013SourceRowKindOffset],
		ContextFormIsEmpty:     parserRow[0x52] == 0,
		ContextPreparationCode: preparationCode,
	}, nil
}

// EvaluatePaul2013ModelContextExceptionGate derives FUN_100091b0's known
// return and applies the caller gate before FUN_10008dc0. The result describes
// eligibility only; an eligible row still needs exception lookup and, when
// unmatched, the later normalization cascade.
func EvaluatePaul2013ModelContextExceptionGate(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
) (Paul2013ModelContextExceptionGate, error) {
	row, _, err := paul2013ModelContextAndParserRows(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ModelContextExceptionGate{}, err
	}
	preparation, err := DerivePaul2013ModelContextPreparation(model, parserRows, contextRowIndex)
	if err != nil {
		return Paul2013ModelContextExceptionGate{}, err
	}
	parserIndex := int(int16(binary.LittleEndian.Uint16(row[paul2013ModelContextTokenIndex:])))
	input, err := BuildPaul2013ModelContextExceptionGate(model, parserRows, contextRowIndex, preparation.Code)
	if err != nil {
		return Paul2013ModelContextExceptionGate{}, err
	}
	return Paul2013ModelContextExceptionGate{
		RowIndex: contextRowIndex, ParserRowIndex: parserIndex,
		Preparation: preparation, Input: input,
		Dispatch: ShouldDispatchPaul2013PronunciationException(input),
	}, nil
}

// ExtractPaul2013ModelExceptionRowInputs reads normalized source surfaces
// and X retry markers from full model context rows in caller-selected order.
func ExtractPaul2013ModelExceptionRowInputs(
	model []byte,
	rowIndexes []int,
) ([]Paul2013ExceptionRowInput, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return nil, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	if len(model) < paul2013ModelContextRowBaseOffset {
		return nil, fmt.Errorf("model has %d bytes, need context rows to start at %#x", len(model), paul2013ModelContextRowBaseOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count < 0 {
		return nil, fmt.Errorf("model has negative context row count %d", count)
	}
	if count > (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize {
		return nil, fmt.Errorf("model declares %d context rows but only %d fit", count, (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize)
	}
	result := make([]Paul2013ExceptionRowInput, len(rowIndexes))
	for index, rowIndex := range rowIndexes {
		if rowIndex < 0 || rowIndex >= count {
			return nil, fmt.Errorf("exception row index %d at position %d is outside 0..%d", rowIndex, index, count-1)
		}
		rowStart := paul2013ModelContextRowBaseOffset + rowIndex*Paul2013PhoneContextRowSize
		row := model[rowStart : rowStart+Paul2013PhoneContextRowSize]
		surfaceArea := row[paul2013ModelContextSurfaceOffset:paul2013ModelContextPhoneOffset]
		nul := bytes.IndexByte(surfaceArea, 0)
		if nul < 0 {
			return nil, fmt.Errorf("model context row %d surface is not NUL-terminated", rowIndex)
		}
		result[index] = Paul2013ExceptionRowInput{
			RowIndex: rowIndex, Surface: string(surfaceArea[:nul]),
			RetryHPrefix: row[paul2013ModelContextClass] == 'X',
		}
	}
	return result, nil
}

// BuildPaul2013NativeExceptionCandidateRows ports FUN_10008dc0's contiguous
// candidate scan from one context row. It retains rows while their normalized
// component total remains within the exception dictionary's observed 1..4
// range. An invalid first surface yields no candidates; a later invalid
// surface or an over-limit component total ends the scan after the valid
// prefix, matching the native early-stop behavior. The per-row X retry marker
// is retained for the lookup stage.
func BuildPaul2013NativeExceptionCandidateRows(
	model []byte,
	startRowIndex int,
) ([]Paul2013ExceptionRowInput, error) {
	if len(model) < paul2013ModelContextCountOffset+2 || len(model) < paul2013ModelContextRowBaseOffset {
		return nil, fmt.Errorf("model has %d bytes, context rows require offsets through %#x", len(model), paul2013ModelContextRowBaseOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 {
		return nil, fmt.Errorf("model context row count %d is negative", rowCount)
	}
	if rowCount > (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize {
		return nil, fmt.Errorf("model declares %d context rows but only %d fit", rowCount, (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize)
	}
	if startRowIndex < 0 || startRowIndex >= rowCount {
		return nil, fmt.Errorf("exception start row %d is outside 0..%d", startRowIndex, rowCount-1)
	}

	attributes := Paul2013ExceptionCharacterAttributes()
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	result := make([]Paul2013ExceptionRowInput, 0, min(rowCount-startRowIndex, paul2013ExceptionCategoryCount))
	componentTotal := uint32(0)
	for rowIndex := startRowIndex; rowIndex < rowCount; rowIndex++ {
		inputs, err := ExtractPaul2013ModelExceptionRowInputs(model, []int{rowIndex})
		if err != nil {
			if len(result) == 0 {
				return nil, nil
			}
			break
		}
		_, category, err := NormalizePaul2013ExceptionSurface(
			[]byte(inputs[0].Surface), true, attributes, characterMap,
		)
		if err != nil {
			if len(result) == 0 {
				return nil, nil
			}
			break
		}
		if componentTotal+category > paul2013ExceptionCategoryCount {
			break
		}
		componentTotal += category
		result = append(result, inputs[0])
		if componentTotal == paul2013ExceptionCategoryCount {
			break
		}
	}
	return result, nil
}

// ApplyPaul2013ExceptionMatchToModelRows ports FUN_1000ca50's exception
// phone-code splitting and model-row writes. A d delimiter is retained until
// the current source surface's one-plus-hyphen count is reached; that boundary
// delimiter advances to the next contiguous model context row.
func ApplyPaul2013ExceptionMatchToModelRows(
	model []byte,
	startRow int,
	match ExceptionMatch,
) ([]byte, []int, error) {
	if match.TokenCount <= 0 {
		return nil, nil, fmt.Errorf("exception match has invalid token count %d", match.TokenCount)
	}
	if len(model) < paul2013ModelContextCountOffset+2 || len(model) < paul2013ModelContextRowBaseOffset {
		return nil, nil, fmt.Errorf("model has %d bytes, context rows require offsets through %#x", len(model), paul2013ModelContextRowBaseOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBaseOffset)/Paul2013PhoneContextRowSize {
		return nil, nil, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if startRow < 0 || startRow >= rowCount {
		return nil, nil, fmt.Errorf("exception destination row %d is outside 0..%d", startRow, rowCount-1)
	}
	working := append([]byte(nil), model...)
	destinationRows := [][]byte{make([]byte, 0)}
	rowIndexes := []int{startRow}
	rowIndex := startRow
	delimiterCount, err := paul2013ModelContextDelimiterCount(model, rowIndex)
	if err != nil {
		return nil, nil, err
	}
	seenDelimiters := 0
	code := match.PhoneCodes
	if nul := bytes.IndexByte(code, 0); nul >= 0 {
		code = code[:nul]
	}
	for _, value := range code {
		if value != 'd' {
			destinationRows[len(destinationRows)-1] = append(destinationRows[len(destinationRows)-1], value)
			continue
		}
		seenDelimiters++
		if seenDelimiters < delimiterCount {
			destinationRows[len(destinationRows)-1] = append(destinationRows[len(destinationRows)-1], value)
			continue
		}
		rowIndex++
		if rowIndex >= rowCount {
			return nil, nil, fmt.Errorf("exception phone codes cross beyond %d model context rows", rowCount)
		}
		rowIndexes = append(rowIndexes, rowIndex)
		destinationRows = append(destinationRows, make([]byte, 0))
		delimiterCount, err = paul2013ModelContextDelimiterCount(model, rowIndex)
		if err != nil {
			return nil, nil, err
		}
		seenDelimiters = 0
	}
	for index, codes := range destinationRows {
		if len(codes) >= paul2013ModelContextPhoneEnd-paul2013ModelContextPhoneOffset {
			return nil, nil, fmt.Errorf("exception output row %d has %d phone codes, capacity is %d", rowIndexes[index], len(codes), paul2013ModelContextPhoneEnd-paul2013ModelContextPhoneOffset-1)
		}
		rowStart := paul2013ModelContextRowBaseOffset + rowIndexes[index]*Paul2013PhoneContextRowSize
		copy(working[rowStart+paul2013ModelContextPhoneOffset:], codes)
		working[rowStart+paul2013ModelContextPhoneOffset+len(codes)] = 0
		binary.LittleEndian.PutUint16(working[rowStart-paul2013ModelContextFlagBackstep:], 2)
	}
	return working, rowIndexes, nil
}

func paul2013ModelContextDelimiterCount(model []byte, rowIndex int) (int, error) {
	rowStart := paul2013ModelContextRowBaseOffset + rowIndex*Paul2013PhoneContextRowSize
	rowEnd := rowStart + Paul2013PhoneContextRowSize
	if rowStart < 0 || rowEnd > len(model) {
		return 0, fmt.Errorf("model context row %d exceeds model size %d", rowIndex, len(model))
	}
	row := model[rowStart:rowEnd]
	surface := row[paul2013ModelContextSurfaceOffset:paul2013ModelContextPhoneOffset]
	nul := bytes.IndexByte(surface, 0)
	if nul < 0 {
		return 0, fmt.Errorf("model context row %d surface is not NUL-terminated", rowIndex)
	}
	return Paul2013ExceptionSurfaceDelimiterCount(string(surface[:nul])), nil
}

func paul2013ModelContextAndParserRows(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
) ([]byte, []byte, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return nil, nil, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	count := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if count < 0 || contextRowIndex < 0 || contextRowIndex >= count {
		return nil, nil, fmt.Errorf("model context row index %d is outside signed count %d", contextRowIndex, count)
	}
	rowStart := paul2013ModelContextRowBaseOffset + contextRowIndex*Paul2013PhoneContextRowSize
	rowEnd := rowStart + Paul2013PhoneContextRowSize
	if rowEnd > len(model) {
		return nil, nil, fmt.Errorf("model context row %d range [%#x,%#x) exceeds model size %d", contextRowIndex, rowStart, rowEnd, len(model))
	}
	if len(parserRows)%paul2013ModelParserOffsetRowStride != 0 {
		return nil, nil, fmt.Errorf("parser rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013ModelParserOffsetRowStride)
	}
	row := model[rowStart:rowEnd]
	parserIndex := int(int16(binary.LittleEndian.Uint16(row[paul2013ModelContextTokenIndex:])))
	if parserIndex < 0 || parserIndex >= len(parserRows)/paul2013ModelParserOffsetRowStride {
		return nil, nil, fmt.Errorf("model context row %d references parser row %d outside %d bytes", contextRowIndex, parserIndex, len(parserRows))
	}
	parserStart := parserIndex * paul2013ModelParserOffsetRowStride
	return row, parserRows[parserStart : parserStart+paul2013ModelParserOffsetRowStride], nil
}

// ApplyPaul2013ModelSourceClassNormalization ports FUN_10009030's ordered
// source-class, phone-type, existing-code, metadata-gate, generic-normalizer,
// and context-encoding branches for one FUN_1000ea20 model row. It reads the
// parser-row index from the model row and preserves all unrelated arena bytes.
func ApplyPaul2013ModelSourceClassNormalization(
	model []byte,
	parserRows []byte,
	contextRowIndex int,
	normalize Paul2013GenericContextNormalizer,
) (Paul2013ModelContextNormalizationResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ModelContextNormalizationResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	contextCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if contextCount < 0 || contextRowIndex < 0 || contextRowIndex >= contextCount {
		return Paul2013ModelContextNormalizationResult{}, fmt.Errorf("model context row index %d is outside signed count %d", contextRowIndex, contextCount)
	}
	rowStart := paul2013ModelContextRowBaseOffset + contextRowIndex*Paul2013PhoneContextRowSize
	rowEnd := rowStart + Paul2013PhoneContextRowSize
	if rowStart < paul2013ModelContextCountOffset+4 || rowEnd > len(model) {
		return Paul2013ModelContextNormalizationResult{}, fmt.Errorf("model context row %d range [%#x,%#x) exceeds model size %d", contextRowIndex, rowStart, rowEnd, len(model))
	}
	row := model[rowStart:rowEnd]
	parserIndex := int(int16(binary.LittleEndian.Uint16(row[paul2013ModelContextTokenIndex:])))
	if parserIndex < 0 || len(parserRows)%paul2013ModelParserOffsetRowStride != 0 || parserIndex >= len(parserRows)/paul2013ModelParserOffsetRowStride {
		return Paul2013ModelContextNormalizationResult{}, fmt.Errorf("model context row %d references parser row %d outside %d bytes", contextRowIndex, parserIndex, len(parserRows))
	}
	parserStart := parserIndex * paul2013ModelParserOffsetRowStride
	parserRow := parserRows[parserStart : parserStart+paul2013ModelParserOffsetRowStride]
	sourceKind := parserRow[paul2013SourceRowKindOffset]
	sourceClass := parserRow[paul2013SourceRowClassOffset]
	if sourceClass != 'S' && sourceClass != 'C' {
		return Paul2013ModelContextNormalizationResult{Model: append([]byte(nil), model...)}, nil
	}
	phoneType := row[paul2013ModelContextStatus]
	surfaceArea := row[paul2013ModelContextSurfaceOffset:paul2013ModelContextPhoneOffset]
	surfaceEnd := bytes.IndexByte(surfaceArea, 0)
	if surfaceEnd < 0 {
		return Paul2013ModelContextNormalizationResult{}, fmt.Errorf("model context row %d surface is not NUL-terminated", contextRowIndex)
	}
	surface := append([]byte(nil), surfaceArea[:surfaceEnd]...)
	phoneArea := row[paul2013ModelContextPhoneOffset:paul2013ModelContextPhoneEnd]
	flagOffset := rowStart - paul2013ModelContextFlagBackstep
	metadataGate := int16(binary.LittleEndian.Uint16(row[paul2013ModelContextGateOffset:]))
	working := append([]byte(nil), model...)
	writeContextCodes := func(codes []byte, flags byte) error {
		if bytes.IndexByte(codes, 0) >= 0 {
			return errors.New("normalizer returned codes with an embedded NUL")
		}
		if len(codes) >= len(phoneArea) {
			return fmt.Errorf("normalizer returned %d codes, phone field capacity is %d", len(codes), len(phoneArea)-1)
		}
		copy(working[rowStart+paul2013ModelContextPhoneOffset:], codes)
		working[rowStart+paul2013ModelContextPhoneOffset+len(codes)] = 0
		working[flagOffset] |= flags
		return nil
	}
	encodeContext := func() error {
		codes, err := EncodePaul2013ContextString(surface)
		if err != nil {
			return fmt.Errorf("encode model context row %d surface: %w", contextRowIndex, err)
		}
		flags := byte(0)
		if len(codes) != 0 {
			flags = 0x20
		}
		return writeContextCodes(codes, flags)
	}
	applyGeneric := func() (bool, error) {
		if normalize == nil {
			return false, errors.New("FUN_10009030 requires the generic normalizer for this model row")
		}
		result, err := normalize(append([]byte(nil), surface...))
		if err != nil {
			return false, fmt.Errorf("normalize model context row %d: %w", contextRowIndex, err)
		}
		if !result.Handled {
			if len(result.Codes) != 0 || result.RowFlags != 0 {
				return false, errors.New("unhandled generic-normalizer result contains writes")
			}
			return false, nil
		}
		if result.RowFlags & ^byte(0x30) != 0 {
			return false, fmt.Errorf("generic normalizer returned unsupported row flags %#02x", result.RowFlags)
		}
		if err := writeContextCodes(result.Codes, result.RowFlags); err != nil {
			return false, err
		}
		return true, nil
	}
	setExistingCodeFlag := func() {
		working[flagOffset] |= 1
	}
	if sourceClass == 'S' {
		if err := encodeContext(); err != nil {
			return Paul2013ModelContextNormalizationResult{}, err
		}
		return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
	}
	contextType := phoneType == 'E' || phoneType == 'A'
	if sourceKind != 'I' && sourceKind != 'M' && sourceKind != 'E' {
		return Paul2013ModelContextNormalizationResult{Model: append([]byte(nil), model...)}, nil
	}
	if contextType {
		if err := encodeContext(); err != nil {
			return Paul2013ModelContextNormalizationResult{}, err
		}
		return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
	}
	if bytes.IndexByte(phoneArea, 0) < 0 {
		return Paul2013ModelContextNormalizationResult{}, fmt.Errorf("model context row %d phone-code field is not NUL-terminated", contextRowIndex)
	}
	if phoneArea[0] != 0 {
		setExistingCodeFlag()
		return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
	}
	if sourceKind == 'I' {
		if metadataGate != 0 {
			handled, err := applyGeneric()
			if err != nil {
				return Paul2013ModelContextNormalizationResult{}, err
			}
			if handled {
				return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
			}
		}
		handled, err := applyGeneric()
		if err != nil {
			return Paul2013ModelContextNormalizationResult{}, err
		}
		if handled {
			return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
		}
	} else if metadataGate != 0 {
		handled, err := applyGeneric()
		if err != nil {
			return Paul2013ModelContextNormalizationResult{}, err
		}
		if handled {
			return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
		}
	}
	if err := encodeContext(); err != nil {
		return Paul2013ModelContextNormalizationResult{}, err
	}
	return Paul2013ModelContextNormalizationResult{Model: working, Handled: true}, nil
}
