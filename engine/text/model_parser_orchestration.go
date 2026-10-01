package text

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/tree3"
)

// Paul2013ModelParserOrchestrationResult reports the parser and model arenas
// produced by the supported portion of FUN_1000d190. A source row can yield
// multiple token rows when FUN_1000d640 leaves a nonempty remainder.
type Paul2013ModelParserOrchestrationResult struct {
	Model             []byte
	ParserRows        []byte
	SourceSegments    []Paul2013ModelParserSegmentRange
	TokenRowsWritten  int
	SourceRowsVisited int
	StoppedAtLimit    bool
}

// Paul2013ModelParserSegmentRange identifies one sentence segment's
// half-open row range in a shared parser-row arena.
type Paul2013ModelParserSegmentRange struct {
	FirstRow int
	RowCount int
}

// Paul2013ParserSourceRowControls contains the parser-row fields that the
// bounded ordinary offset builder does not recover. Values map to +0x2c,
// +0x30, and +0x52 respectively; their producers remain caller-owned.
type Paul2013ParserSourceRowControls struct {
	RowType       uint32
	Mode          int8
	AuxiliaryText []byte
}

// BuildPaul2013ModelTokenRowsFromParserRows ports the observed
// FUN_1000d190 traversal, dictionary lookup, and FUN_1000d450 row writes for
// initialized 0x94-byte parser rows. It preserves bytes outside the model
// fields written by the native helpers.
//
// pathPhoneEnabled is the caller's FUN_1000d190 param_4 gate used by the
// conditional marker branch in FUN_1000d450. Parser-row construction and the
// earlier normalization/TPP cascade remain upstream stages.
func BuildPaul2013ModelTokenRowsFromParserRows(
	model []byte,
	parserRows []byte,
	dictionary *EmbeddedDictionary,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	if len(parserRows)%paul2013ParserSourceRowStride != 0 {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("parser source rows have %d bytes, not a multiple of 0x%x", len(parserRows), paul2013ParserSourceRowStride)
	}
	rowCount := len(parserRows) / paul2013ParserSourceRowStride
	if rowCount > paul2013SourceRowLimit {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("parser source rows have %d entries, native limit is %d", rowCount, paul2013SourceRowLimit)
	}
	if rowCount == 0 {
		return Paul2013ModelParserOrchestrationResult{
			Model: append([]byte(nil), model...), ParserRows: append([]byte(nil), parserRows...),
		}, nil
	}
	if dictionary == nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("model parser orchestration has no embedded dictionary")
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("model buffer has %d bytes, need at least %d", len(model), paul2013ModelContextCountOffset+2)
	}
	modelTokenCount := int(binary.LittleEndian.Uint16(model[paul2013ModelTokenRowCountOffset:]))
	if modelTokenCount > paul2013ModelTokenRowLimit {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("model already has %d token rows, native limit is %d", modelTokenCount, paul2013ModelTokenRowLimit)
	}
	if paul2013ModelTokenRowsOffset+modelTokenCount*Paul2013TokenResultRowSize > len(model) {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("model buffer does not contain its %d existing token rows", modelTokenCount)
	}
	working := append([]byte(nil), model...)
	result := Paul2013ModelParserOrchestrationResult{
		Model: working, ParserRows: append([]byte(nil), parserRows...),
	}
	for parserIndex := 0; parserIndex < rowCount; parserIndex++ {
		if modelTokenCount >= paul2013ModelTokenRowLimit {
			result.StoppedAtLimit = true
			break
		}
		row := parserRows[parserIndex*paul2013ParserSourceRowStride : (parserIndex+1)*paul2013ParserSourceRowStride]
		transformed, err := TransformPaul2013ModelParserSourceRowWithEmbeddedDictionary(
			row, parserIndex, rowCount, dictionary,
		)
		if err != nil {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("transform parser source row %d: %w", parserIndex, err)
		}
		result.SourceRowsVisited++
		mode := int8(row[paul2013ParserSourceRowModeOffset])
		auxiliary, err := paul2013ParserSourceCString(row, paul2013ParserSourceRowAuxOffset, paul2013ParserSourceRowStride)
		if err != nil {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("read parser source row %d auxiliary text: %w", parserIndex, err)
		}
		priorParserIndex := -1
		if modelTokenCount > 0 {
			priorOffset := paul2013ModelTokenRowsOffset + (modelTokenCount-1)*Paul2013TokenResultRowSize
			priorParserIndex = int(int16(binary.LittleEndian.Uint16(working[priorOffset+Paul2013TokenResultRowIndex:])))
		}
		current := transformed.Transform
		for splitIndex := 0; ; splitIndex++ {
			if splitIndex >= 32 {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("parser source row %d produced more than 32 native source fragments", parserIndex)
			}
			if modelTokenCount >= paul2013ModelTokenRowLimit {
				result.StoppedAtLimit = true
				break
			}
			pronunciation, found, err := dictionary.ResolvePaul2013Surface(current.Output)
			if err != nil {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("resolve parser source row %d fragment %d: %w", parserIndex, splitIndex, err)
			}
			phoneRows := Paul2013DictionaryPhoneRows{
				PathControlBytes: []byte{0xff},
			}
			if found {
				phoneRows, err = pronunciation.Payload.BuildPaul2013DictionaryPhoneRowsForState(
					Paul2013PhoneIDCodebook(),
					Paul2013DictionaryPhoneRowState{
						FinalRow:        parserIndex == priorParserIndex,
						HasMarker:       mode != -1,
						Marker:          byte(mode),
						SelectPathPhone: pathPhoneEnabled && len(auxiliary) == 0,
					},
				)
				if err != nil {
					return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("build parser source row %d fragment %d phone rows: %w", parserIndex, splitIndex, err)
				}
			}
			tokenRowOffset := paul2013ModelTokenRowsOffset + modelTokenCount*Paul2013TokenResultRowSize
			tokenRowEnd := tokenRowOffset + Paul2013TokenResultRowSize
			if tokenRowEnd > paul2013ModelContextCountOffset || tokenRowEnd > len(working) {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("model buffer cannot hold token row %d before the context table", modelTokenCount)
			}
			if err := WritePaul2013DictionaryPhoneRow(
				working[tokenRowOffset:tokenRowEnd],
				uint16(parserIndex), current.Output, phoneRows,
			); err != nil {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("write parser source row %d fragment %d token row: %w", parserIndex, splitIndex, err)
			}
			working[Paul2013TokenResultRowType+tokenRowOffset] = ClassifyPaul2013ModelSourceType(
				current.Output, Paul2013ExceptionCharacterAttributes(),
			)
			modelTokenCount++
			binary.LittleEndian.PutUint16(working[paul2013ModelTokenRowCountOffset:], uint16(modelTokenCount))
			result.TokenRowsWritten++
			priorParserIndex = parserIndex
			if len(current.Remainder) == 0 {
				break
			}
			previousCount := int16(0)
			if row[paul2013ParserSourceRowRawClassOffset] == 'S' {
				previousCount = 1
			}
			current, err = TransformPaul2013ModelSourceStringFromPaul2013CharacterTable(
				Paul2013ModelSourceTransformInput{
					Source:         current.Remainder,
					ParserType:     mode,
					AssociatedText: auxiliary,
					PreviousCount:  previousCount,
					FinalRow:       parserIndex == rowCount-1,
				},
				func(surface []byte) (bool, error) { return dictionary.ContainsPaul2013Surface(surface) },
			)
			if err != nil {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("transform parser source row %d fragment %d: %w", parserIndex, splitIndex+1, err)
			}
		}
		if result.StoppedAtLimit {
			break
		}
	}
	result.Model = working
	return result, nil
}

// BuildPaul2013ModelPhoneRowsFromParserRows composes the supported
// FUN_1000d190 token-row traversal with FUN_1000ea20's fixed context-row
// projection and the recovered ordinary multi-pronunciation selector.
// FUN_100049b0's special name/context shortcut and the upstream normalization
// and TPP cascade remain outside this composition.
func BuildPaul2013ModelPhoneRowsFromParserRows(
	model []byte,
	parserRows []byte,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	sourceLength int,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	result, err := BuildPaul2013ModelTokenRowsFromParserRows(
		model, parserRows, dictionary, pathPhoneEnabled,
	)
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, err
	}
	selector := SelectPaul2013ModelPronunciationAlternative(pronunciationTree)
	result.Model, err = NormalizePaul2013ModelPhoneRows(
		result.Model, parserRows, sourceLength, selector,
	)
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("project model phone-context rows: %w", err)
	}
	return result, nil
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySource composes the bounded ordinary
// word/number offset producer with the model token/context path. The caller
// must provide one native control record for every generated parser row;
// this function does not infer row types, mode bytes, or the auxiliary string
// from source spelling. It derives +0x08 from the generated +0x34 C-string
// length, as observed in the native source-row append helper.
func BuildPaul2013ModelPhoneRowsFromOrdinarySource(
	model []byte,
	source []byte,
	rowControls []Paul2013ParserSourceRowControls,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	return BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications(
		model, source, rowControls, nil, dictionary, pronunciationTree, pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications composes
// the bounded ordinary offset rows, caller-provided initial row controls,
// ordered recovered +0x2c helper applications, and the supported token/context
// projection. Mode and auxiliary text remain explicit because their producers
// are not recovered by the ordinary offset builder.
func BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications(
	model []byte,
	source []byte,
	rowControls []Paul2013ParserSourceRowControls,
	applications []Paul2013ParserRowTypeApplication,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(string(source))
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("derive ordinary parser offset rows: %w", err)
	}
	if len(segments) != 1 {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("ordinary source produced %d parser segments; one model parser call accepts one segment", len(segments))
	}
	rows := segments[0]
	if len(rowControls) != len(rows) {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("ordinary source produced %d parser rows but %d row controls were supplied", len(rows), len(rowControls))
	}
	parserState := make([]byte, paul2013ModelParserRowsOffset+len(rows)*paul2013ModelParserRowStride)
	initialTypes := make([]uint32, len(rowControls))
	for rowIndex, controls := range rowControls {
		initialTypes[rowIndex] = controls.RowType
	}
	if err := WritePaul2013ParserOffsetRowsWithInitialTypesAndTypeApplications(
		parserState, rows, initialTypes, applications,
	); err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("write ordinary parser rows and type applications: %w", err)
	}
	parserRows := parserState[paul2013ModelParserRowsOffset:]
	for rowIndex, controls := range rowControls {
		row := parserRows[rowIndex*paul2013ModelParserRowStride : (rowIndex+1)*paul2013ModelParserRowStride]
		if containsNUL(controls.AuxiliaryText) {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("parser row %d auxiliary text contains an embedded NUL", rowIndex)
		}
		if len(controls.AuxiliaryText) >= paul2013ModelParserRowStride-paul2013ParserSourceRowAuxOffset {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("parser row %d auxiliary text has %d bytes, maximum is %d", rowIndex, len(controls.AuxiliaryText), paul2013ModelParserRowStride-paul2013ParserSourceRowAuxOffset-1)
		}
		row[paul2013ParserSourceRowModeOffset] = byte(controls.Mode)
		copy(row[paul2013ParserSourceRowAuxOffset:], controls.AuxiliaryText)
		row[paul2013ParserSourceRowAuxOffset+len(controls.AuxiliaryText)] = 0
	}
	return BuildPaul2013ModelPhoneRowsFromParserRows(
		model, parserRows, dictionary, pronunciationTree, len(source), pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithDefaultParserControls
// composes the bounded ordinary source path with caller-supplied initial
// parser-row +0x2c values and ordered type applications. It initializes
// parser-row +0x30 to the captured untagged mode 0xff and leaves +0x52 empty,
// matching the ordinary word, punctuation, and bounded numeric captures.
// The +0x2c row-type producer remains an explicit input.
func BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithDefaultParserControls(
	model []byte,
	source []byte,
	initialRowTypes []uint32,
	applications []Paul2013ParserRowTypeApplication,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(string(source))
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("derive ordinary parser offset rows: %w", err)
	}
	if len(segments) != 1 {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("ordinary source produced %d parser segments; one model parser call accepts one segment", len(segments))
	}
	if len(initialRowTypes) != len(segments[0]) {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("ordinary source produced %d parser rows but %d initial row types were supplied", len(segments[0]), len(initialRowTypes))
	}
	controls := make([]Paul2013ParserSourceRowControls, len(initialRowTypes))
	for rowIndex, rowType := range initialRowTypes {
		controls[rowIndex] = Paul2013ParserSourceRowControls{
			RowType: rowType,
			Mode:    -1,
		}
	}
	return BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications(
		model, source, controls, applications, dictionary, pronunciationTree, pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySegments composes every sentence
// segment emitted by the bounded ordinary offset producer with the model
// token/context projection. rowControls has one entry per produced segment,
// then one control record per row in that segment. Each result starts from
// the same input model, so results are segment-local. Use
// BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena to preserve the
// native shared model counters across those segments.
//
// This covers the observed offset reset at `.`, `?`, and `!` while keeping
// parser control fields explicit. It does not merge segment context tables or
// synthesize sentence-boundary rows.
func BuildPaul2013ModelPhoneRowsFromOrdinarySegments(
	model []byte,
	source []byte,
	rowControls [][]Paul2013ParserSourceRowControls,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) ([]Paul2013ModelParserOrchestrationResult, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(string(source))
	if err != nil {
		return nil, fmt.Errorf("derive ordinary parser offset rows: %w", err)
	}
	if len(rowControls) != len(segments) {
		return nil, fmt.Errorf("ordinary source produced %d parser segments but %d segment control lists were supplied", len(segments), len(rowControls))
	}
	if len(segments) != 0 {
		if len(model) < paul2013ModelContextCountOffset+2 || len(model) < 2 {
			return nil, fmt.Errorf("segment-local model buffer has %d bytes; need parser and context row counters", len(model))
		}
		if binary.LittleEndian.Uint16(model[:2]) != 0 ||
			binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:]) != 0 {
			return nil, fmt.Errorf("segment-local ordinary-source projection requires zero initial token and context row counts")
		}
	}
	results := make([]Paul2013ModelParserOrchestrationResult, 0, len(segments))
	for segmentIndex, rows := range segments {
		controls := rowControls[segmentIndex]
		if len(controls) != len(rows) {
			return nil, fmt.Errorf("ordinary parser segment %d produced %d rows but %d row controls were supplied", segmentIndex, len(rows), len(controls))
		}
		parserState := make([]byte, paul2013ModelParserRowsOffset+len(rows)*paul2013ModelParserRowStride)
		if err := WritePaul2013ParserOffsetRows(parserState, rows); err != nil {
			return nil, fmt.Errorf("write ordinary parser segment %d rows: %w", segmentIndex, err)
		}
		parserRows := parserState[paul2013ModelParserRowsOffset:]
		for rowIndex, control := range controls {
			row := parserRows[rowIndex*paul2013ModelParserRowStride : (rowIndex+1)*paul2013ModelParserRowStride]
			if containsNUL(control.AuxiliaryText) {
				return nil, fmt.Errorf("parser segment %d row %d auxiliary text contains an embedded NUL", segmentIndex, rowIndex)
			}
			if len(control.AuxiliaryText) >= paul2013ModelParserRowStride-paul2013ParserSourceRowAuxOffset {
				return nil, fmt.Errorf("parser segment %d row %d auxiliary text has %d bytes, maximum is %d", segmentIndex, rowIndex, len(control.AuxiliaryText), paul2013ModelParserRowStride-paul2013ParserSourceRowAuxOffset-1)
			}
			binary.LittleEndian.PutUint32(row[paul2013ParserSourceRowTypeOffset:], control.RowType)
			row[paul2013ParserSourceRowModeOffset] = byte(control.Mode)
			copy(row[paul2013ParserSourceRowAuxOffset:], control.AuxiliaryText)
			row[paul2013ParserSourceRowAuxOffset+len(control.AuxiliaryText)] = 0
		}
		result, err := BuildPaul2013ModelPhoneRowsFromParserRows(
			model, parserRows, dictionary, pronunciationTree, len(source), pathPhoneEnabled,
		)
		if err != nil {
			return nil, fmt.Errorf("build model phone rows for ordinary parser segment %d: %w", segmentIndex, err)
		}
		results = append(results, result)
	}
	return results, nil
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena flattens the
// bounded sentence segments in source order, preserves each segment's
// segment-relative offsets, then runs the native parser-row traversal and
// context projection once over the shared model arena. FUN_1000d190 advances
// its source-row and model-token counters across the complete row array, while
// FUN_1000e2f0 initializes that arena once per parser call. Row controls remain
// explicit because the ordinary offset producer does not recover them.
func BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
	model []byte,
	source []byte,
	rowControls [][]Paul2013ParserSourceRowControls,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	return BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		model, source, rowControls, nil, dictionary, pronunciationTree, pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications
// extends the shared-arena projection with per-segment ordered parser-row
// type applications. Application RowIndex values are local to their segment;
// successful writes are translated to the flattened parser-row indexes after
// each segment's native row-count gates have been evaluated.
func BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications(
	model []byte,
	source []byte,
	rowControls [][]Paul2013ParserSourceRowControls,
	applications [][]Paul2013ParserRowTypeApplication,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(string(source))
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("derive ordinary parser offset rows: %w", err)
	}
	if len(rowControls) != len(segments) {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
			"ordinary source produced %d parser segments but %d segment control lists were supplied",
			len(segments), len(rowControls),
		)
	}
	if len(applications) == 0 {
		applications = make([][]Paul2013ParserRowTypeApplication, len(segments))
	}
	if len(applications) != len(segments) {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
			"ordinary source produced %d parser segments but %d segment type-application lists were supplied",
			len(segments), len(applications),
		)
	}
	if len(model) < paul2013ModelContextCountOffset+2 || len(model) < 2 {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
			"shared model buffer has %d bytes; need parser and context row counters", len(model),
		)
	}
	if binary.LittleEndian.Uint16(model[:2]) != 0 ||
		binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:]) != 0 {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
			"shared-arena ordinary-source projection requires zero initial token and context row counts",
		)
	}

	var rows []Paul2013OrdinaryParserOffsetRow
	ranges := make([]Paul2013ModelParserSegmentRange, 0, len(segments))
	var typeWrites []Paul2013ParserRowTypeWrite
	for segmentIndex, segmentRows := range segments {
		controls := rowControls[segmentIndex]
		if len(controls) != len(segmentRows) {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
				"ordinary parser segment %d produced %d rows but %d row controls were supplied",
				segmentIndex, len(segmentRows), len(controls),
			)
		}
		firstRow := len(rows)
		initialTypes := make([]uint32, len(controls))
		for localRowIndex, control := range controls {
			initialTypes[localRowIndex] = control.RowType
			typeWrites = append(typeWrites, Paul2013ParserRowTypeWrite{
				RowIndex: firstRow + localRowIndex,
				Value:    control.RowType,
			})
		}
		localWrites, err := BuildPaul2013ParserRowTypeWritesFromApplicationsWithInitialTypes(
			initialTypes, applications[segmentIndex],
		)
		if err != nil {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
				"build parser row type writes for segment %d: %w", segmentIndex, err,
			)
		}
		for _, write := range localWrites {
			write.RowIndex += firstRow
			typeWrites = append(typeWrites, write)
		}
		rows = append(rows, segmentRows...)
		ranges = append(ranges, Paul2013ModelParserSegmentRange{
			FirstRow: firstRow,
			RowCount: len(segmentRows),
		})
	}
	if len(rows) == 0 {
		return Paul2013ModelParserOrchestrationResult{
			Model: append([]byte(nil), model...), SourceSegments: ranges,
		}, nil
	}

	parserState := make([]byte, paul2013ModelParserRowsOffset+len(rows)*paul2013ModelParserRowStride)
	for rowIndex := range rows {
		parserState[paul2013ModelParserRowsOffset+rowIndex*paul2013ModelParserRowStride+paul2013ModelParserRowMode] = 0xff
	}
	if err := WritePaul2013ParserOffsetRowsWithTypeWrites(parserState, rows, typeWrites); err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("write shared ordinary parser rows and type applications: %w", err)
	}
	parserRows := parserState[paul2013ModelParserRowsOffset:]
	for segmentIndex, segmentRange := range ranges {
		for localRowIndex, control := range rowControls[segmentIndex] {
			rowIndex := segmentRange.FirstRow + localRowIndex
			row := parserRows[rowIndex*paul2013ModelParserRowStride : (rowIndex+1)*paul2013ModelParserRowStride]
			if containsNUL(control.AuxiliaryText) {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
					"parser segment %d row %d auxiliary text contains an embedded NUL", segmentIndex, localRowIndex,
				)
			}
			if len(control.AuxiliaryText) >= paul2013ModelParserRowStride-paul2013ParserSourceRowAuxOffset {
				return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
					"parser segment %d row %d auxiliary text has %d bytes, maximum is %d",
					segmentIndex, localRowIndex, len(control.AuxiliaryText),
					paul2013ModelParserRowStride-paul2013ParserSourceRowAuxOffset-1,
				)
			}
			row[paul2013ParserSourceRowModeOffset] = byte(control.Mode)
			copy(row[paul2013ParserSourceRowAuxOffset:], control.AuxiliaryText)
			row[paul2013ParserSourceRowAuxOffset+len(control.AuxiliaryText)] = 0
		}
	}

	result, err := BuildPaul2013ModelPhoneRowsFromParserRows(
		model, parserRows, dictionary, pronunciationTree, len(source), pathPhoneEnabled,
	)
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("project shared ordinary parser rows: %w", err)
	}
	result.SourceSegments = ranges
	return result, nil
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithDefaultParserControls
// applies the captured untagged mode 0xff and empty auxiliary string to every
// ordinary parser row while retaining caller-supplied initial +0x2c row types
// and per-segment type applications. It preserves the shared-arena counters
// and segment-relative offsets of the underlying composition.
func BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithDefaultParserControls(
	model []byte,
	source []byte,
	initialRowTypes [][]uint32,
	applications [][]Paul2013ParserRowTypeApplication,
	dictionary *EmbeddedDictionary,
	pronunciationTree *tree3.Tree,
	pathPhoneEnabled bool,
) (Paul2013ModelParserOrchestrationResult, error) {
	segments, err := BuildPaul2013OrdinaryParserOffsetRows(string(source))
	if err != nil {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf("derive ordinary parser offset rows: %w", err)
	}
	if len(initialRowTypes) != len(segments) {
		return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
			"ordinary source produced %d parser segments but %d segment row-type lists were supplied",
			len(segments), len(initialRowTypes),
		)
	}
	controls := make([][]Paul2013ParserSourceRowControls, len(segments))
	for segmentIndex, segment := range segments {
		if len(initialRowTypes[segmentIndex]) != len(segment) {
			return Paul2013ModelParserOrchestrationResult{}, fmt.Errorf(
				"ordinary parser segment %d produced %d rows but %d initial row types were supplied",
				segmentIndex, len(segment), len(initialRowTypes[segmentIndex]),
			)
		}
		controls[segmentIndex] = make([]Paul2013ParserSourceRowControls, len(segment))
		for rowIndex, rowType := range initialRowTypes[segmentIndex] {
			controls[segmentIndex][rowIndex] = Paul2013ParserSourceRowControls{
				RowType: rowType,
				Mode:    -1,
			}
		}
	}
	return BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		model, source, controls, applications, dictionary, pronunciationTree, pathPhoneEnabled,
	)
}
