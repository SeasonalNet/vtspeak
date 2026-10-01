package text

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	paul2013ModelParserModeOffset  = 0x39e8
	paul2013ModelParserRowsOffset  = 0x20
	paul2013ModelParserRowStride   = 0x94
	paul2013ModelParserRowCount    = 100
	paul2013ModelParserRowMode     = 0x30
	paul2013ModelParserRowText     = 0x34
	paul2013ModelParserRowAux      = 0x52
	paul2013ModelParserRowSentinel = uint32(0xffffffff)
	paul2013ModelParserInitOffset  = 0x14
	paul2013ModelParserInitDwords  = 0xe74
)

// Paul2013ModelParserKind identifies the inner parser selected by
// FUN_1003d350. The parser implementations remain caller supplied.
type Paul2013ModelParserKind uint8

const (
	Paul2013PrimaryModelParser Paul2013ModelParserKind = iota + 1
	Paul2013AlternateModelParser
)

// Paul2013ModelParserCallbacks supplies the still-unported parser routines
// and finalizer. Each callback receives the raw parser state because its
// layout is only partially recovered.
type Paul2013ModelParserCallbacks struct {
	Primary   func(input, state []byte) (int32, error)
	Alternate func(input, state []byte) (int32, error)
	Finalize  func(state []byte, parserResult int32) (int16, error)
}

// Paul2013ModelParserDispatchResult records the selected inner parser and
// both results from FUN_1003d350's parser/finalizer sequence.
type Paul2013ModelParserDispatchResult struct {
	Kind         Paul2013ModelParserKind
	ParserResult int32
	FinalResult  int16
	// Success is the native low-short-to-BOOL conversion (FinalResult != 0).
	// The finalizer can return -1 on an error path, so true does not establish
	// that parsing was valid.
	Success bool
}

// Paul2013ModeProbeHasWord ports FUN_1005f2c0: a null pointer or a zero
// little-endian word at pointed offset +8 returns false; any nonzero word
// returns true. probe represents the bytes at the pointed-to address.
func Paul2013ModeProbeHasWord(probe []byte) (bool, error) {
	if len(probe) == 0 {
		return false, nil
	}
	if len(probe) < 10 {
		return false, fmt.Errorf("Paul 2013 parser mode probe has %d bytes, need at least 10", len(probe))
	}
	return binary.LittleEndian.Uint16(probe[8:10]) != 0, nil
}

// InitializePaul2013ModelParserState ports FUN_1003e210. It clears the
// leading counters and 0xe74 dwords beginning at +0x14. The control word at
// +0x0c is preserved only when it is exactly 1; otherwise it is cleared.
func InitializePaul2013ModelParserState(state []byte) error {
	const required = paul2013ModelParserInitOffset + paul2013ModelParserInitDwords*4
	if len(state) < required {
		return fmt.Errorf("Paul 2013 parser state has %d bytes, need at least %d", len(state), required)
	}
	preserveControl := binary.LittleEndian.Uint16(state[0x0c:0x0e]) == 1
	binary.LittleEndian.PutUint32(state[4:8], 0)
	binary.LittleEndian.PutUint16(state[0:2], 0)
	binary.LittleEndian.PutUint16(state[2:4], 0)
	if !preserveControl {
		binary.LittleEndian.PutUint16(state[0x0c:0x0e], 0)
	}
	clear(state[paul2013ModelParserInitOffset:required])
	return nil
}

// WritePaul2013ParserOffsetRows writes observed source offsets and the type
// byte into 0x94-byte parser-state rows. When Text is present it writes its
// byte length at +0x08 and the NUL-terminated string at +0x34, matching the
// source-row append layout and FUN_1000d190 reads. AuxiliaryText is written at
// +0x52. Other bytes, including initialized sentinels and parser controls,
// remain unchanged. Offsets are written as supplied; callers that combine
// sentence segments preserve each row's segment-relative origin.
func WritePaul2013ParserOffsetRows(
	state []byte,
	rows []Paul2013OrdinaryParserOffsetRow,
) error {
	if len(rows) > paul2013ModelParserRowCount {
		return fmt.Errorf("Paul 2013 parser has %d rows, maximum is %d", len(rows), paul2013ModelParserRowCount)
	}
	if len(state) < 2 {
		return fmt.Errorf("Paul 2013 parser state has %d bytes, need 2 for row count", len(state))
	}
	if len(rows) > 0 {
		lastRow := paul2013ModelParserRowsOffset + (len(rows)-1)*paul2013ModelParserRowStride
		needed := lastRow + 0x25
		for index, row := range rows {
			rowBase := paul2013ModelParserRowsOffset + index*paul2013ModelParserRowStride
			if row.Text != nil {
				if containsNUL(row.Text) {
					return fmt.Errorf("parser offset row %d text contains an embedded NUL", index)
				}
				maxTextLength := paul2013ModelParserRowAux - paul2013ModelParserRowText - 1
				if len(row.Text) > maxTextLength {
					return fmt.Errorf("parser offset row %d text has %d bytes, maximum before auxiliary field is %d", index, len(row.Text), maxTextLength)
				}
				rowEnd := rowBase + paul2013ModelParserRowText + len(row.Text) + 1
				if rowEnd > needed {
					needed = rowEnd
				}
			}
			if row.AuxiliaryText != nil {
				if containsNUL(row.AuxiliaryText) {
					return fmt.Errorf("parser offset row %d auxiliary text contains an embedded NUL", index)
				}
				maxAuxiliaryLength := paul2013ModelParserRowStride - paul2013ModelParserRowAux - 1
				if len(row.AuxiliaryText) > maxAuxiliaryLength {
					return fmt.Errorf("parser offset row %d auxiliary text has %d bytes, maximum is %d", index, len(row.AuxiliaryText), maxAuxiliaryLength)
				}
				rowEnd := rowBase + paul2013ModelParserRowAux + len(row.AuxiliaryText) + 1
				if rowEnd > needed {
					needed = rowEnd
				}
			}
		}
		if len(state) < needed {
			return fmt.Errorf("Paul 2013 parser state has %d bytes, need %d for %d projected rows", len(state), needed, len(rows))
		}
	}
	binary.LittleEndian.PutUint16(state[0:2], uint16(int16(len(rows))))
	for index, row := range rows {
		offset := paul2013ModelParserRowsOffset + index*paul2013ModelParserRowStride
		binary.LittleEndian.PutUint32(state[offset:offset+4], uint32(row.Start))
		binary.LittleEndian.PutUint32(state[offset+4:offset+8], uint32(row.End))
		state[offset+0x24] = row.RawTypeByte
		if row.Text != nil {
			binary.LittleEndian.PutUint32(state[offset+0x08:offset+0x0c], uint32(len(row.Text)))
			copy(state[offset+paul2013ModelParserRowText:], row.Text)
			state[offset+paul2013ModelParserRowText+len(row.Text)] = 0
		}
		if row.AuxiliaryText != nil {
			copy(state[offset+paul2013ModelParserRowAux:], row.AuxiliaryText)
			state[offset+paul2013ModelParserRowAux+len(row.AuxiliaryText)] = 0
		}
	}
	return nil
}

// Paul2013ParserRowTypeWrite is one explicitly observed write to parser-row
// dword +0x2c. RowIndex is zero-based in the row slice passed to the writer.
type Paul2013ParserRowTypeWrite struct {
	RowIndex int
	Value    uint32
}

// Paul2013ParserRowTypeApplication applies one existing parser-row type
// helper to an explicitly indexed ordinary parser-offset row. Apply receives
// a temporary source-row arena whose final row corresponds to RowIndex. Its
// row count and prior +0x2c value are initialized; the callback must provide
// any other source fields its helper reads, preserve gate inputs, and return
// its normal (arena, applied, error) result.
type Paul2013ParserRowTypeApplication struct {
	RowIndex int
	Apply    func([]byte) ([]byte, bool, error)
}

// Paul2013ParserRowTypeGate applies one already-resolved native scanner gate
// to an ordinary parser-offset row. Exactly one of Terminal, Comma, Dot, or
// ScanFallback must be non-nil. RowIndex addresses the row in the same slice
// passed to the offset writer; it does not infer which source token a native
// scanner accepts.
type Paul2013ParserRowTypeGate struct {
	RowIndex     int
	Terminal     *Paul2013ParserTerminalTypeGate
	Comma        *Paul2013ParserCommaTypeGate
	Dot          *Paul2013ParserDotTypeGate
	ScanFallback *Paul2013ParserScanFallbackRowTypeGate
}

// BuildPaul2013ParserRowTypeWritesFromGates composes the terminal, comma,
// statically recovered dot, and scan-fallback +0x2c writers with an ordered
// ordinary-row sequence. The scanner predicates remain explicit in each
// gate. A temporary source-row arena is used so each helper retains its native
// row-index behavior, while writes from earlier gates to the same row are
// preserved in order. The dot and scan-fallback branches lack direct runtime
// write observations.
func BuildPaul2013ParserRowTypeWritesFromGates(
	rowCount int,
	gates []Paul2013ParserRowTypeGate,
) ([]Paul2013ParserRowTypeWrite, error) {
	applications := make([]Paul2013ParserRowTypeApplication, 0, len(gates))
	for gateIndex, gate := range gates {
		gateCount := 0
		if gate.Terminal != nil {
			gateCount++
		}
		if gate.Comma != nil {
			gateCount++
		}
		if gate.Dot != nil {
			gateCount++
		}
		if gate.ScanFallback != nil {
			gateCount++
		}
		if gateCount != 1 {
			return nil, fmt.Errorf("parser row type gate %d must provide exactly one gate", gateIndex)
		}
		application := Paul2013ParserRowTypeApplication{RowIndex: gate.RowIndex}
		if gate.Terminal != nil {
			terminalGate := *gate.Terminal
			application.Apply = func(arena []byte) ([]byte, bool, error) {
				return ApplyPaul2013ParserTerminalRowType(arena, terminalGate)
			}
		} else if gate.Comma != nil {
			commaGate := *gate.Comma
			application.Apply = func(arena []byte) ([]byte, bool, error) {
				return ApplyPaul2013ParserCommaRowType(arena, commaGate)
			}
		} else if gate.Dot != nil {
			dotGate := *gate.Dot
			application.Apply = func(arena []byte) ([]byte, bool, error) {
				return ApplyPaul2013ParserDotRowType(arena, dotGate)
			}
		} else {
			scanFallbackGate := *gate.ScanFallback
			scanFallbackGate.RowIndex = gate.RowIndex + 1
			application.Apply = func(arena []byte) ([]byte, bool, error) {
				return ApplyPaul2013ParserScanFallbackRowType(arena, scanFallbackGate)
			}
		}
		applications = append(applications, application)
	}
	return BuildPaul2013ParserRowTypeWritesFromApplications(rowCount, applications)
}

// BuildPaul2013ParserRowTypeWritesFromApplications composes any existing
// final-row +0x2c helper with ordinary parser-offset rows. It retains type
// values between applications to the same row and returns successful writes
// in callback order. Each application still supplies the helper's recovered
// gate fields; this function does not infer scanner or model-parser state.
func BuildPaul2013ParserRowTypeWritesFromApplications(
	rowCount int,
	applications []Paul2013ParserRowTypeApplication,
) ([]Paul2013ParserRowTypeWrite, error) {
	if rowCount < 0 || rowCount > paul2013ModelParserRowCount {
		return nil, fmt.Errorf("Paul 2013 parser has %d rows, maximum is %d", rowCount, paul2013ModelParserRowCount)
	}
	return BuildPaul2013ParserRowTypeWritesFromApplicationsWithInitialTypes(
		make([]uint32, rowCount), applications,
	)
}

// BuildPaul2013ParserRowTypeWritesFromApplicationsWithInitialTypes composes
// the helper callbacks over caller-supplied initial +0x2c values. This models
// source rows whose earlier native producer has already assigned a type; each
// callback then observes the prior value when applying its own zero checks or
// overwrites.
func BuildPaul2013ParserRowTypeWritesFromApplicationsWithInitialTypes(
	initialTypes []uint32,
	applications []Paul2013ParserRowTypeApplication,
) ([]Paul2013ParserRowTypeWrite, error) {
	rowCount := len(initialTypes)
	if rowCount > paul2013ModelParserRowCount {
		return nil, fmt.Errorf("Paul 2013 parser has %d rows, maximum is %d", rowCount, paul2013ModelParserRowCount)
	}
	types := append([]uint32(nil), initialTypes...)
	writes := make([]Paul2013ParserRowTypeWrite, 0, len(applications))
	for applicationIndex, application := range applications {
		if application.RowIndex < 0 || application.RowIndex >= rowCount {
			return nil, fmt.Errorf("parser row type application %d references row %d outside %d rows", applicationIndex, application.RowIndex, rowCount)
		}
		if application.Apply == nil {
			return nil, fmt.Errorf("parser row type application %d has no callback", applicationIndex)
		}
		arenaSize := paul2013SourceRowsOffset + (application.RowIndex+1)*paul2013SourceRowStride
		arena := make([]byte, arenaSize)
		binary.LittleEndian.PutUint16(arena[paul2013SourceRowCountOffset:], uint16(application.RowIndex+1))
		field := paul2013SourceRowsOffset + application.RowIndex*paul2013SourceRowStride + 0x2c
		binary.LittleEndian.PutUint32(arena[field:], types[application.RowIndex])
		updated, applied, err := application.Apply(arena)
		if err != nil {
			return nil, fmt.Errorf("apply parser row type application %d: %w", applicationIndex, err)
		}
		if !applied {
			continue
		}
		if len(updated) < field+4 {
			return nil, fmt.Errorf("parser row type application %d returned %d bytes, need %d for row type", applicationIndex, len(updated), field+4)
		}
		value := binary.LittleEndian.Uint32(updated[field:])
		types[application.RowIndex] = value
		writes = append(writes, Paul2013ParserRowTypeWrite{RowIndex: application.RowIndex, Value: value})
	}
	return writes, nil
}

// WritePaul2013ParserOffsetRowsWithTypeGates composes the ordinary offset
// projection with the captured terminal and comma writers plus the statically
// recovered dot writer. It requires explicit, row-indexed scanner gate values
// and applies successful writes in input order. Other +0x2c producers remain
// caller supplied through WritePaul2013ParserOffsetRowsWithTypeWrites.
func WritePaul2013ParserOffsetRowsWithTypeGates(
	state []byte,
	rows []Paul2013OrdinaryParserOffsetRow,
	gates []Paul2013ParserRowTypeGate,
) error {
	writes, err := BuildPaul2013ParserRowTypeWritesFromGates(len(rows), gates)
	if err != nil {
		return err
	}
	return WritePaul2013ParserOffsetRowsWithTypeWrites(state, rows, writes)
}

// WritePaul2013ParserOffsetRowsWithTypeApplications composes ordinary offset
// rows with any existing final-row +0x2c helper. Applications remain
// row-indexed, ordered, and responsible for supplying their native gate
// values. It provides composition for recovered branches beyond the typed
// terminal/comma/dot convenience API.
func WritePaul2013ParserOffsetRowsWithTypeApplications(
	state []byte,
	rows []Paul2013OrdinaryParserOffsetRow,
	applications []Paul2013ParserRowTypeApplication,
) error {
	return WritePaul2013ParserOffsetRowsWithInitialTypesAndTypeApplications(
		state, rows, make([]uint32, len(rows)), applications,
	)
}

// WritePaul2013ParserOffsetRowsWithInitialTypesAndTypeApplications composes
// ordinary offsets, caller-supplied initial row types, and ordered helper
// applications. Initial types are written before applications, so each
// helper receives the current +0x2c value for its target row.
func WritePaul2013ParserOffsetRowsWithInitialTypesAndTypeApplications(
	state []byte,
	rows []Paul2013OrdinaryParserOffsetRow,
	initialTypes []uint32,
	applications []Paul2013ParserRowTypeApplication,
) error {
	if len(initialTypes) != len(rows) {
		return fmt.Errorf("ordinary parser has %d rows but %d initial row types", len(rows), len(initialTypes))
	}
	applicationWrites, err := BuildPaul2013ParserRowTypeWritesFromApplicationsWithInitialTypes(initialTypes, applications)
	if err != nil {
		return err
	}
	writes := make([]Paul2013ParserRowTypeWrite, 0, len(rows)+len(applicationWrites))
	for rowIndex, value := range initialTypes {
		writes = append(writes, Paul2013ParserRowTypeWrite{RowIndex: rowIndex, Value: value})
	}
	writes = append(writes, applicationWrites...)
	return WritePaul2013ParserOffsetRowsWithTypeWrites(state, rows, writes)
}

// WritePaul2013ParserOffsetRowsWithTypeWrites composes the ordinary offset
// projection with sparse, caller-supplied +0x2c writes. It does not infer row
// types from offsets or the +0x24 discriminator. Writes are applied in slice
// order, matching the source order when multiple native branches target one
// row.
func WritePaul2013ParserOffsetRowsWithTypeWrites(
	state []byte,
	rows []Paul2013OrdinaryParserOffsetRow,
	writes []Paul2013ParserRowTypeWrite,
) error {
	for writeIndex, write := range writes {
		if write.RowIndex < 0 || write.RowIndex >= len(rows) {
			return fmt.Errorf("parser row type write %d references row %d outside %d rows", writeIndex, write.RowIndex, len(rows))
		}
	}
	working := append([]byte(nil), state...)
	if err := WritePaul2013ParserOffsetRows(working, rows); err != nil {
		return err
	}
	for writeIndex, write := range writes {
		field := paul2013ModelParserRowsOffset + write.RowIndex*paul2013ModelParserRowStride + 0x2c
		if field+4 > len(working) {
			return fmt.Errorf("parser row type write %d needs bytes through %#x, state has %d", writeIndex, field+4, len(working))
		}
		binary.LittleEndian.PutUint32(working[field:field+4], write.Value)
	}
	copy(state, working)
	return nil
}

// BuildPaul2013OrdinaryParserOffsetRowArenaWithTypeGates builds one captured
// ordinary sentence segment and composes the explicit terminal, comma, and
// dot gates into its parser rows. Multi-segment input remains rejected because
// each segment has its own offset origin.
func BuildPaul2013OrdinaryParserOffsetRowArenaWithTypeGates(
	source []byte,
	gates []Paul2013ParserRowTypeGate,
) ([]byte, error) {
	return buildPaul2013OrdinaryParserOffsetRowArena(source, func(state []byte, rows []Paul2013OrdinaryParserOffsetRow) error {
		return WritePaul2013ParserOffsetRowsWithTypeGates(state, rows, gates)
	})
}

// BuildPaul2013OrdinaryParserOffsetRowArenaWithTypeApplications builds one
// captured ordinary sentence segment and invokes the supplied existing row
// type helpers for their indexed rows. Applications retain their explicit
// native gates and any additional source-field setup.
func BuildPaul2013OrdinaryParserOffsetRowArenaWithTypeApplications(
	source []byte,
	applications []Paul2013ParserRowTypeApplication,
) ([]byte, error) {
	return buildPaul2013OrdinaryParserOffsetRowArena(source, func(state []byte, rows []Paul2013OrdinaryParserOffsetRow) error {
		return WritePaul2013ParserOffsetRowsWithTypeApplications(state, rows, applications)
	})
}

// BuildPaul2013OrdinaryParserOffsetRowArena composes the bounded ordinary
// offset producer with the raw 0x94-byte parser-row representation consumed by
// FUN_1000ea20. The captured producer resets offsets at terminal sentence
// separators, while this arena represents one parser segment, so inputs that
// produce more than one segment are rejected instead of flattening offsets
// from distinct coordinate origins.
func BuildPaul2013OrdinaryParserOffsetRowArena(source []byte) ([]byte, error) {
	return buildPaul2013OrdinaryParserOffsetRowArena(source, func(state []byte, rows []Paul2013OrdinaryParserOffsetRow) error {
		return WritePaul2013ParserOffsetRows(state, rows)
	})
}

func buildPaul2013OrdinaryParserOffsetRowArena(
	source []byte,
	writeRows func([]byte, []Paul2013OrdinaryParserOffsetRow) error,
) ([]byte, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(string(source))
	if err != nil {
		return nil, fmt.Errorf("build ordinary parser offset rows: %w", err)
	}
	if len(segments) > 1 {
		return nil, fmt.Errorf("ordinary parser source produced %d sentence segments; one parser arena represents one segment", len(segments))
	}
	if len(segments) == 0 || len(segments[0]) == 0 {
		return nil, errors.New("ordinary parser source produced no rows")
	}
	rows := segments[0]
	state := make([]byte, paul2013ModelParserRowsOffset+len(rows)*paul2013ModelParserRowStride)
	for rowIndex := range rows {
		rowModeOffset := paul2013ModelParserRowsOffset + rowIndex*paul2013ModelParserRowStride + paul2013ModelParserRowMode
		state[rowModeOffset] = 0xff
	}
	if err := writeRows(state, rows); err != nil {
		return nil, fmt.Errorf("write ordinary parser rows: %w", err)
	}
	return append([]byte(nil), state[paul2013ModelParserRowsOffset:]...), nil
}

// DispatchPaul2013ModelParser ports the initialization and parser selection
// in FUN_1003d350. modeProbe must contain the pointed-to bytes when the state
// dword at +0x39e8 is nonzero; it is ignored when that dword is zero. The
// function initializes the first dword of each of 100 records to -1, selects
// the alternate parser only for a nonnull mode pointer whose +8 word is zero,
// and converts any nonzero finalizer low short, including -1, to true. This
// mirrors the native BOOL conversion and is not a parse-validity guarantee.
func DispatchPaul2013ModelParser(
	input []byte,
	state []byte,
	modeProbe []byte,
	callbacks Paul2013ModelParserCallbacks,
) (Paul2013ModelParserDispatchResult, error) {
	if len(state) < paul2013ModelParserRequiredSize() {
		return Paul2013ModelParserDispatchResult{}, fmt.Errorf(
			"Paul 2013 parser state has %d bytes, need at least %d",
			len(state), paul2013ModelParserRequiredSize(),
		)
	}
	if callbacks.Finalize == nil {
		return Paul2013ModelParserDispatchResult{}, errors.New("Paul 2013 parser dispatch requires a finalizer")
	}
	if err := InitializePaul2013ModelParserState(state); err != nil {
		return Paul2013ModelParserDispatchResult{}, fmt.Errorf("initialize Paul 2013 parser state: %w", err)
	}
	for row := 0; row < paul2013ModelParserRowCount; row++ {
		offset := paul2013ModelParserRowsOffset + row*paul2013ModelParserRowStride
		binary.LittleEndian.PutUint32(state[offset:offset+4], paul2013ModelParserRowSentinel)
	}

	kind := Paul2013PrimaryModelParser
	if binary.LittleEndian.Uint32(state[paul2013ModelParserModeOffset:paul2013ModelParserModeOffset+4]) != 0 {
		if len(modeProbe) == 0 {
			return Paul2013ModelParserDispatchResult{}, errors.New("Paul 2013 parser mode pointer is nonzero but its pointed-to bytes were not supplied")
		}
		hasWord, err := Paul2013ModeProbeHasWord(modeProbe)
		if err != nil {
			return Paul2013ModelParserDispatchResult{}, err
		}
		if !hasWord {
			kind = Paul2013AlternateModelParser
		}
	}

	parser := callbacks.Primary
	if kind == Paul2013AlternateModelParser {
		parser = callbacks.Alternate
	}
	if parser == nil {
		return Paul2013ModelParserDispatchResult{}, fmt.Errorf("Paul 2013 parser callback for selected kind %d is nil", kind)
	}
	parserResult, err := parser(input, state)
	if err != nil {
		return Paul2013ModelParserDispatchResult{Kind: kind}, fmt.Errorf("run Paul 2013 parser kind %d: %w", kind, err)
	}
	finalResult, err := callbacks.Finalize(state, parserResult)
	if err != nil {
		return Paul2013ModelParserDispatchResult{Kind: kind, ParserResult: parserResult}, fmt.Errorf("finalize Paul 2013 parser result: %w", err)
	}
	return Paul2013ModelParserDispatchResult{
		Kind: kind, ParserResult: parserResult, FinalResult: finalResult, Success: finalResult != 0,
	}, nil
}

func paul2013ModelParserRequiredSize() int {
	rowsEnd := paul2013ModelParserRowsOffset +
		(paul2013ModelParserRowCount-1)*paul2013ModelParserRowStride + 4
	modeEnd := paul2013ModelParserModeOffset + 4
	if rowsEnd > modeEnd {
		return rowsEnd
	}
	return modeEnd
}
