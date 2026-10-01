package text

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

const (
	paul2013ParserStateRowsOffset = 0x14
	paul2013ParserStateRowStride  = 0x94
	paul2013ParserStateRowLimit   = 100
	paul2013ParserStateStartField = 0
	paul2013ParserStateEndField   = 4
)

// Paul2013PositionIntervalOffsets are the two signed offsets emitted by the
// parser rows consumed in FUN_10022dc0.
type Paul2013PositionIntervalOffsets struct {
	Start int32
	End   int32
}

// Paul2013OrdinaryParserOffsetRow retains the row-type byte observed beside
// the parser's exclusive source offsets. Start and End are relative to the
// current sentence segment; RawTypeByte is an opaque parser discriminator.
type Paul2013OrdinaryParserOffsetRow struct {
	Start       int32
	End         int32
	RawTypeByte byte
	// Text is the primary C string at parser-row +0x34 consumed by FUN_1000d190.
	// Numeric rows receive expanded lexemes from the existing bounded number
	// normalizers. Nil preserves the existing row field when callers project
	// offsets only.
	Text []byte
	// AuxiliaryText represents an explicitly recovered row +0x52 C string.
	// The ordinary offset builder does not currently produce it.
	AuxiliaryText []byte
}

// BuildPaul2013PositionIntervals ports FUN_10022dc0's conversion from parser
// offsets to inclusive absolute spans. When End is greater than Start, the
// native path subtracts one from the absolute end; otherwise the interval is
// a single point at Start. Parser offsets remain caller-provided.
func BuildPaul2013PositionIntervals(
	basePosition int32,
	offsets []Paul2013PositionIntervalOffsets,
) ([]Paul2013PositionEventRange, error) {
	result := make([]Paul2013PositionEventRange, len(offsets))
	for index, row := range offsets {
		start := int64(basePosition) + int64(row.Start)
		end := start
		if row.Start < row.End {
			end = int64(basePosition) + int64(row.End) - 1
		}
		if start < math.MinInt32 || start > math.MaxInt32 || end < math.MinInt32 || end > math.MaxInt32 {
			return nil, fmt.Errorf("position interval %d exceeds signed 32-bit range", index)
		}
		result[index] = Paul2013PositionEventRange{Minimum: int32(start), Maximum: int32(end)}
	}
	return result, nil
}

// BuildPaul2013OrdinaryWordParserOffsetSegments derives the raw start/end
// offsets emitted for ordinary ASCII letter words by the traced model-parser
// path. It rejects numeric and other parser classes.
func BuildPaul2013OrdinaryWordParserOffsetSegments(
	source string,
) ([][]Paul2013PositionIntervalOffsets, error) {
	return buildPaul2013OrdinaryParserOffsetSegments(source, false)
}

// BuildPaul2013OrdinaryParserOffsetSegments derives ordinary ASCII word and
// normalized-number parser offsets. Runtime rows use exclusive end offsets. A
// single period, question mark, or exclamation mark closes a segment and
// resets the next segment's offsets to zero. Commas and ASCII whitespace stay
// within a segment. Runtime captures cover words, unsigned cardinals, +12,
// 1.25, $5.00, 25%, 21st, 01/02/2024, 3:45 PM, and 555-1234.
// The existing normalizers supply row counts for related numeric forms as
// inference. Telephone parser rows use cardinal group spans plus a hyphen row;
// the frontend expansion uses the directly observed group wording.
func BuildPaul2013OrdinaryParserOffsetSegments(
	source string,
) ([][]Paul2013PositionIntervalOffsets, error) {
	return buildPaul2013OrdinaryParserOffsetSegments(source, true)
}

// BuildPaul2013OrdinaryParserOffsetRows adds the runtime-observed row-type
// byte and +0x34 text to the bounded ordinary-word and numeric parser rows.
// Numeric text is supplied by the existing bounded number normalizers and
// matches expanded lexemes observed in Stage 20. ASCII words use
// 'A'; normalized numeric rows use 'D', with the captured leading sign,
// final currency, and telephone hyphen rows using 'S'. Percent suffix rows
// use 'A'. Numeric row multiplicities beyond the directly captured forms
// remain inference.
// It does not construct the remaining 0x94-byte parser-row fields or
// finalizer state.
func BuildPaul2013OrdinaryParserOffsetRows(
	source string,
) ([][]Paul2013OrdinaryParserOffsetRow, error) {
	if source == "" {
		return nil, nil
	}
	tokens, err := tokenizeASCIISurfaces(source)
	if err != nil {
		return nil, fmt.Errorf("tokenize ordinary parser source: %w", err)
	}
	if len(tokens) == 0 {
		return nil, nil
	}
	segments := make([][]Paul2013OrdinaryParserOffsetRow, 0, 1)
	segment := make([]Paul2013OrdinaryParserOffsetRow, 0, len(tokens))
	segmentBase := tokens[0].SourceByteStart
	previousEnd := -1
	for tokenIndex, token := range tokens {
		if token.SourceByteStart < 0 || token.SourceByteEnd < token.SourceByteStart ||
			previousEnd >= token.SourceByteStart {
			return nil, fmt.Errorf("parser token %d has invalid or overlapping source span %d..%d", tokenIndex, token.SourceByteStart, token.SourceByteEnd)
		}
		start := token.SourceByteStart - segmentBase
		end := token.SourceByteEnd + 1 - segmentBase
		if start < 0 || end < start || int64(end) > math.MaxInt32 {
			return nil, fmt.Errorf("parser token %d source offsets exceed the signed 32-bit range", tokenIndex)
		}
		var tokenRows []Paul2013PositionIntervalOffsets
		var tokenTexts []string
		if asciiLettersOnly(token.Surface) {
			tokenRows = []Paul2013PositionIntervalOffsets{{Start: 0, End: int32(len(token.Surface))}}
			tokenTexts = []string{token.Surface}
		} else {
			tokenRows, err = paul2013ParserTokenOffsets(token.Surface)
			if err != nil {
				return nil, fmt.Errorf("derive parser offsets for token %d %q: %w", tokenIndex, token.Surface, err)
			}
			tokenTexts, err = paul2013ParserTokenTexts(token.Surface)
			if err != nil {
				return nil, fmt.Errorf("derive parser text for token %d %q: %w", tokenIndex, token.Surface, err)
			}
			if len(tokenTexts) != len(tokenRows) {
				return nil, fmt.Errorf("parser token %d %q produced %d offsets and %d text rows", tokenIndex, token.Surface, len(tokenRows), len(tokenTexts))
			}
		}
		for rowIndex, row := range tokenRows {
			typeByte := paul2013ParserOffsetRowType(token.Surface, row, rowIndex, len(tokenRows))
			segment = append(segment, Paul2013OrdinaryParserOffsetRow{
				Start:       int32(start) + row.Start,
				End:         int32(start) + row.End,
				RawTypeByte: typeByte,
				Text:        []byte(tokenTexts[rowIndex]),
			})
		}
		previousEnd = token.SourceByteEnd
		terminal, err := ordinaryParserTerminalSeparator(token.SeparatorAfter)
		if err != nil {
			return nil, fmt.Errorf("parser token %d separator: %w", tokenIndex, err)
		}
		if terminal {
			segments = append(segments, segment)
			segment = make([]Paul2013OrdinaryParserOffsetRow, 0, len(tokens)-tokenIndex-1)
			if tokenIndex+1 < len(tokens) {
				segmentBase = tokens[tokenIndex+1].SourceByteStart
			}
		}
	}
	if len(segment) > 0 {
		segments = append(segments, segment)
	}
	return segments, nil
}

// BuildPaul2013OrdinaryParserPositionIntervals composes the bounded ordinary
// text offset rows with FUN_10022dc0's inclusive-coordinate conversion. Each
// returned inner slice is one sentence segment. basePosition is the native
// caller's absolute coordinate at source byte zero; each segment's first
// source byte is added because the native caller advances its coordinate by
// consumed source length. The source parser still emits only the recovered
// offsets and row-type byte, not complete 0x94-byte rows or the event/table
// values consumed by later state passes.
func BuildPaul2013OrdinaryParserPositionIntervals(
	source string,
	basePosition int32,
) ([][]Paul2013PositionEventRange, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(source)
	if err != nil {
		return nil, err
	}
	tokens, err := tokenizeASCIISurfaces(source)
	if err != nil {
		return nil, fmt.Errorf("tokenize ordinary parser source origins: %w", err)
	}
	if len(segments) == 0 {
		return nil, nil
	}
	segmentBases := make([]int, 0, len(segments))
	if len(tokens) > 0 {
		segmentBases = append(segmentBases, tokens[0].SourceByteStart)
	}
	for tokenIndex, token := range tokens {
		terminal, separatorErr := ordinaryParserTerminalSeparator(token.SeparatorAfter)
		if separatorErr != nil {
			return nil, fmt.Errorf("parser token %d separator: %w", tokenIndex, separatorErr)
		}
		if terminal && tokenIndex+1 < len(tokens) {
			segmentBases = append(segmentBases, tokens[tokenIndex+1].SourceByteStart)
		}
	}
	if len(segmentBases) != len(segments) {
		return nil, fmt.Errorf("derived %d sentence origins for %d parser segments", len(segmentBases), len(segments))
	}
	result := make([][]Paul2013PositionEventRange, len(segments))
	for segmentIndex, segment := range segments {
		offsets := make([]Paul2013PositionIntervalOffsets, len(segment))
		for rowIndex, row := range segment {
			offsets[rowIndex] = Paul2013PositionIntervalOffsets{Start: row.Start, End: row.End}
		}
		absoluteBase := int64(basePosition) + int64(segmentBases[segmentIndex])
		if absoluteBase < math.MinInt32 || absoluteBase > math.MaxInt32 {
			return nil, fmt.Errorf("ordinary parser segment %d base exceeds signed 32-bit range", segmentIndex)
		}
		intervals, intervalErr := BuildPaul2013PositionIntervals(int32(absoluteBase), offsets)
		if intervalErr != nil {
			return nil, fmt.Errorf("convert ordinary parser segment %d: %w", segmentIndex, intervalErr)
		}
		result[segmentIndex] = intervals
	}
	return result, nil
}

func paul2013ParserOffsetRowType(
	surface string,
	row Paul2013PositionIntervalOffsets,
	rowIndex int,
	rowCount int,
) byte {
	if asciiLettersOnly(surface) {
		return 'A'
	}
	if (surface[0] == '+' || surface[0] == '-') && row.Start == 0 && row.End == 1 {
		return 'S'
	}
	if strings.HasPrefix(surface, "$") && rowIndex == rowCount-1 {
		return 'S'
	}
	if strings.HasSuffix(surface, "%") && row.Start == int32(len(surface)-1) && row.End == int32(len(surface)) {
		return 'A'
	}
	if strings.Contains(surface, "-") {
		separatorOffset := int32(strings.IndexByte(surface, '-'))
		if row.Start == separatorOffset && row.End == separatorOffset+1 {
			return 'S'
		}
	}
	return 'D'
}

func paul2013ParserTokenTexts(surface string) ([]string, error) {
	if allASCIIDigits(surface) {
		return ExpandPaul2013UnsignedInteger(surface)
	}
	if len(surface) > 1 && (surface[0] == '+' || surface[0] == '-') && allASCIIDigits(surface[1:]) {
		words, err := ExpandPaul2013UnsignedInteger(surface[1:])
		if err != nil {
			return nil, err
		}
		sign := "minus"
		if surface[0] == '+' {
			sign = "plus"
		}
		return append([]string{sign}, words...), nil
	}
	if strings.HasPrefix(surface, "$") {
		return ExpandPaul2013Currency(surface)
	}
	if strings.HasSuffix(surface, "%") {
		return ExpandPaul2013Percentage(surface)
	}
	if strings.Contains(surface, "/") {
		return ExpandPaul2013SlashDate(surface)
	}
	if strings.Contains(surface, ":") {
		return ExpandPaul2013ClockTime(surface)
	}
	if strings.Contains(surface, "-") {
		return ExpandPaul2013Telephone(surface)
	}
	if strings.HasSuffix(surface, "st") || strings.HasSuffix(surface, "nd") ||
		strings.HasSuffix(surface, "rd") || strings.HasSuffix(surface, "th") {
		return ExpandPaul2013Ordinal(surface)
	}
	return ExpandPaul2013Number(surface)
}

func buildPaul2013OrdinaryParserOffsetSegments(
	source string,
	allowNumericClasses bool,
) ([][]Paul2013PositionIntervalOffsets, error) {
	if source == "" {
		return nil, nil
	}
	tokens, err := tokenizeASCIISurfaces(source)
	if err != nil {
		return nil, fmt.Errorf("tokenize ordinary parser source: %w", err)
	}
	segments := make([][]Paul2013PositionIntervalOffsets, 0, 1)
	segment := make([]Paul2013PositionIntervalOffsets, 0, len(tokens))
	segmentBase := tokens[0].SourceByteStart
	previousEnd := -1
	for index, token := range tokens {
		if token.Surface == "" {
			return nil, fmt.Errorf("parser token %d %q is outside the ordinary ASCII-letter path", index, token.Surface)
		}
		tokenRows := []Paul2013PositionIntervalOffsets{{Start: 0, End: int32(len(token.Surface))}}
		if !asciiLettersOnly(token.Surface) {
			if !allowNumericClasses {
				return nil, fmt.Errorf("parser token %d %q is outside the ordinary ASCII-letter path", index, token.Surface)
			}
			var expansionErr error
			tokenRows, expansionErr = paul2013ParserTokenOffsets(token.Surface)
			if expansionErr != nil {
				return nil, fmt.Errorf("derive parser offsets for token %d %q: %w", index, token.Surface, expansionErr)
			}
		}
		if token.SourceByteStart < 0 || token.SourceByteEnd < token.SourceByteStart ||
			previousEnd >= token.SourceByteStart {
			return nil, fmt.Errorf("parser token %d has invalid or overlapping source span %d..%d", index, token.SourceByteStart, token.SourceByteEnd)
		}
		start := token.SourceByteStart - segmentBase
		end := token.SourceByteEnd + 1 - segmentBase
		if start < 0 || end < start || int64(end) > math.MaxInt32 {
			return nil, fmt.Errorf("parser token %d source offsets exceed the signed 32-bit range", index)
		}
		for _, tokenRow := range tokenRows {
			segment = append(segment, Paul2013PositionIntervalOffsets{
				Start: int32(start) + tokenRow.Start,
				End:   int32(start) + tokenRow.End,
			})
		}
		previousEnd = token.SourceByteEnd

		terminal, separatorErr := ordinaryParserTerminalSeparator(token.SeparatorAfter)
		if separatorErr != nil {
			return nil, fmt.Errorf("parser token %d separator: %w", index, separatorErr)
		}
		if terminal {
			segments = append(segments, segment)
			segment = make([]Paul2013PositionIntervalOffsets, 0, len(tokens)-index-1)
			if index+1 < len(tokens) {
				segmentBase = tokens[index+1].SourceByteStart
			}
		}
	}
	if len(segment) > 0 {
		segments = append(segments, segment)
	}
	return segments, nil
}

func paul2013ParserTokenOffsets(surface string) ([]Paul2013PositionIntervalOffsets, error) {
	if allASCIIDigits(surface) {
		words, err := ExpandPaul2013UnsignedInteger(surface)
		return repeatedParserTokenSpan(len(surface), len(words), err)
	}
	if len(surface) > 1 && (surface[0] == '+' || surface[0] == '-') && allASCIIDigits(surface[1:]) {
		words, err := ExpandPaul2013UnsignedInteger(surface[1:])
		rows, err := repeatedParserTokenSpan(len(surface)-1, len(words), err)
		if err != nil {
			return nil, err
		}
		result := make([]Paul2013PositionIntervalOffsets, 0, len(rows)+1)
		result = append(result, Paul2013PositionIntervalOffsets{Start: 0, End: 1})
		for _, row := range rows {
			row.Start++
			row.End++
			result = append(result, row)
		}
		return result, nil
	}
	if strings.HasPrefix(surface, "$") {
		words, err := ExpandPaul2013Currency(surface)
		return repeatedParserTokenSpan(len(surface), len(words), err)
	}
	if strings.HasSuffix(surface, "%") {
		number := surface[:len(surface)-1]
		rows, err := paul2013ParserTokenOffsets(number)
		if err != nil {
			return nil, fmt.Errorf("percentage number: %w", err)
		}
		return append(rows, Paul2013PositionIntervalOffsets{Start: int32(len(number)), End: int32(len(surface))}), nil
	}
	if strings.Contains(surface, "/") {
		words, err := ExpandPaul2013SlashDate(surface)
		return repeatedParserTokenSpan(len(surface), len(words), err)
	}
	if strings.Contains(surface, ":") {
		words, err := ExpandPaul2013ClockTime(surface)
		return repeatedParserTokenSpan(len(surface), len(words), err)
	}
	if strings.Contains(surface, "-") {
		parts := strings.Split(surface, "-")
		if len(parts) != 2 || len(parts[0]) != 3 || len(parts[1]) != 4 ||
			!allASCIIDigits(parts[0]) || !allASCIIDigits(parts[1]) {
			return nil, fmt.Errorf("telephone token must use the captured NNN-NNNN form")
		}
		leftWords, err := ExpandPaul2013UnsignedInteger(parts[0])
		if err != nil {
			return nil, fmt.Errorf("telephone first group: %w", err)
		}
		rightWords, err := ExpandPaul2013UnsignedInteger(parts[1])
		if err != nil {
			return nil, fmt.Errorf("telephone second group: %w", err)
		}
		left, err := repeatedParserTokenSpan(len(parts[0]), len(leftWords), nil)
		if err != nil {
			return nil, err
		}
		left = append(left, Paul2013PositionIntervalOffsets{Start: int32(len(parts[0])), End: int32(len(parts[0]) + 1)})
		right, err := repeatedParserTokenSpan(len(parts[1]), len(rightWords), nil)
		if err != nil {
			return nil, err
		}
		for index := range right {
			right[index].Start += int32(len(parts[0]) + 1)
			right[index].End += int32(len(parts[0]) + 1)
		}
		return append(left, right...), nil
	}
	if strings.HasSuffix(surface, "st") || strings.HasSuffix(surface, "nd") ||
		strings.HasSuffix(surface, "rd") || strings.HasSuffix(surface, "th") {
		words, err := ExpandPaul2013Ordinal(surface)
		return repeatedParserTokenSpan(len(surface), len(words), err)
	}
	words, err := ExpandPaul2013Number(surface)
	if err != nil {
		return nil, fmt.Errorf("unsupported numeric parser token: %w", err)
	}
	return repeatedParserTokenSpan(len(surface), len(words), nil)
}

func repeatedParserTokenSpan(
	spanLength int,
	count int,
	expansionErr error,
) ([]Paul2013PositionIntervalOffsets, error) {
	if expansionErr != nil {
		return nil, expansionErr
	}
	if spanLength <= 0 || count <= 0 || int64(spanLength) > math.MaxInt32 {
		return nil, fmt.Errorf("invalid expanded token span length %d or row count %d", spanLength, count)
	}
	rows := make([]Paul2013PositionIntervalOffsets, count)
	for index := range rows {
		rows[index].End = int32(spanLength)
	}
	return rows, nil
}

func asciiLettersOnly(value string) bool {
	for index := 0; index < len(value); index++ {
		if (value[index] < 'A' || value[index] > 'Z') && (value[index] < 'a' || value[index] > 'z') {
			return false
		}
	}
	return len(value) > 0
}

func ordinaryParserTerminalSeparator(value string) (bool, error) {
	terminal := false
	for index := 0; index < len(value); index++ {
		char := value[index]
		if char == '.' || char == '!' || char == '?' {
			if terminal {
				return false, fmt.Errorf("multiple sentence terminators in %q are unsupported", value)
			}
			terminal = true
			continue
		}
		if char == ' ' || char == '\t' || char == '\r' || char == '\n' || (!terminal && char == ',') {
			continue
		}
		return false, fmt.Errorf("separator byte %#x in %q is unsupported", char, value)
	}
	return terminal, nil
}

// BuildPaul2013PositionIntervalsFromParserState reads the parser result rows
// consumed by FUN_10022dc0 and converts their signed start/end offsets to
// inclusive absolute intervals. The parser state uses a signed 16-bit row
// count at +0 and 0x94-byte records beginning at +0x14; each row's signed
// 32-bit start and end offsets are at row offsets +0 and +4. It does not produce
// those rows or claim that either inner parser has been ported.
func BuildPaul2013PositionIntervalsFromParserState(
	basePosition int32,
	state []byte,
) ([]Paul2013PositionEventRange, error) {
	if len(state) < 2 {
		return nil, fmt.Errorf("Paul 2013 parser state has %d bytes, need at least 2 for row count", len(state))
	}
	count := int16(binary.LittleEndian.Uint16(state[:2]))
	if count < 0 || count > paul2013ParserStateRowLimit {
		return nil, fmt.Errorf("Paul 2013 parser state row count %d is outside [0, %d]", count, paul2013ParserStateRowLimit)
	}
	if count == 0 {
		return []Paul2013PositionEventRange{}, nil
	}
	lastRow := paul2013ParserStateRowsOffset + (int(count)-1)*paul2013ParserStateRowStride
	needed := lastRow + paul2013ParserStateEndField + 4
	if len(state) < needed {
		return nil, fmt.Errorf("Paul 2013 parser state has %d bytes, need at least %d for %d rows", len(state), needed, count)
	}
	offsets := make([]Paul2013PositionIntervalOffsets, int(count))
	for row := range offsets {
		rowOffset := paul2013ParserStateRowsOffset + row*paul2013ParserStateRowStride
		startOffset := rowOffset + paul2013ParserStateStartField
		endOffset := rowOffset + paul2013ParserStateEndField
		offsets[row] = Paul2013PositionIntervalOffsets{
			Start: int32(binary.LittleEndian.Uint32(state[startOffset : startOffset+4])),
			End:   int32(binary.LittleEndian.Uint32(state[endOffset : endOffset+4])),
		}
	}
	return BuildPaul2013PositionIntervals(basePosition, offsets)
}

// MapPaul2013PositionStateIndexes ports FUN_10022dc0's post-state-array map
// pass. In ordinary mode it clamps each index to [0, maximumIndex) before
// separate start/end table lookups. In remapped mode it skips that clamp and
// applies a final lookup to both table results. All model tables are explicit.
func MapPaul2013PositionStateIndexes(
	startIndexes []int32,
	endIndexes []int32,
	startTable []int32,
	endTable []int32,
	finalTable []int32,
	maximumIndex int32,
	remappedMode bool,
) ([]Paul2013PositionEventRange, error) {
	if len(startIndexes) != len(endIndexes) {
		return nil, fmt.Errorf("received %d start indexes and %d end indexes", len(startIndexes), len(endIndexes))
	}
	if !remappedMode && maximumIndex <= 0 {
		return nil, fmt.Errorf("ordinary position mapping requires a positive maximum index")
	}
	if len(startTable) == 0 || len(endTable) == 0 {
		return nil, fmt.Errorf("position mapping requires nonempty start and end tables")
	}
	if remappedMode && len(finalTable) == 0 {
		return nil, fmt.Errorf("remapped position mapping requires a nonempty final table")
	}
	result := make([]Paul2013PositionEventRange, len(startIndexes))
	for index := range startIndexes {
		startIndex, endIndex := startIndexes[index], endIndexes[index]
		if !remappedMode {
			startIndex = clampPaul2013PositionIndex(startIndex, maximumIndex)
			endIndex = clampPaul2013PositionIndex(endIndex, maximumIndex)
		}
		startValue, err := paul2013PositionTableValue(startTable, startIndex)
		if err != nil {
			return nil, fmt.Errorf("position %d start: %w", index, err)
		}
		endValue, err := paul2013PositionTableValue(endTable, endIndex)
		if err != nil {
			return nil, fmt.Errorf("position %d end: %w", index, err)
		}
		if remappedMode {
			startValue, err = paul2013PositionTableValue(finalTable, startValue)
			if err != nil {
				return nil, fmt.Errorf("position %d remapped start: %w", index, err)
			}
			endValue, err = paul2013PositionTableValue(finalTable, endValue)
			if err != nil {
				return nil, fmt.Errorf("position %d remapped end: %w", index, err)
			}
		}
		result[index] = Paul2013PositionEventRange{Minimum: startValue, Maximum: endValue}
	}
	return result, nil
}

// Paul2013PositionIndexTables contains the explicit lookup resources consumed
// by FUN_10022dc0 after its state passes. The table names retain their
// observed start/end/final roles without assigning linguistic meaning.
type Paul2013PositionIndexTables struct {
	Start        []int32
	End          []int32
	Final        []int32
	MaximumIndex int32
	RemappedMode bool
}

// MapPaul2013PositionStateRecordIndexes reads the signed start/end indexes
// at +0x64c/+0x650 in each 0x3c0-byte model-state record, then applies the
// separate table lookups and optional final remap from FUN_10022dc0. records
// is the complete enclosing arena; its first record begins at +0x64c.
func MapPaul2013PositionStateRecordIndexes(
	records []byte,
	rowCount int,
	tables Paul2013PositionIndexTables,
) ([]Paul2013PositionEventRange, error) {
	const (
		recordStride = 0x3c0
		startOffset  = 0x64c
		endOffset    = 0x650
	)
	if rowCount < 0 || rowCount > paul2013ParserStateRowLimit {
		return nil, fmt.Errorf("Paul 2013 mapped row count %d is outside [0, %d]", rowCount, paul2013ParserStateRowLimit)
	}
	if rowCount == 0 {
		return []Paul2013PositionEventRange{}, nil
	}
	if len(records) < startOffset || rowCount > (len(records)-startOffset)/recordStride {
		return nil, fmt.Errorf("model-state arena has %d bytes, cannot hold %d records from +%#x", len(records), rowCount, startOffset)
	}
	starts := make([]int32, rowCount)
	ends := make([]int32, rowCount)
	for row := 0; row < rowCount; row++ {
		recordStart := row * recordStride
		starts[row] = int32(binary.LittleEndian.Uint32(records[recordStart+startOffset : recordStart+startOffset+4]))
		ends[row] = int32(binary.LittleEndian.Uint32(records[recordStart+endOffset : recordStart+endOffset+4]))
	}
	return MapPaul2013PositionStateIndexes(
		starts, ends, tables.Start, tables.End, tables.Final,
		tables.MaximumIndex, tables.RemappedMode,
	)
}

func clampPaul2013PositionIndex(value, maximum int32) int32 {
	if value < 0 {
		return 0
	}
	if value >= maximum {
		return maximum - 1
	}
	return value
}

func paul2013PositionTableValue(table []int32, index int32) (int32, error) {
	if index < 0 || int64(index) >= int64(len(table)) {
		return 0, fmt.Errorf("table index %d is outside %d entries", index, len(table))
	}
	return table[index], nil
}
