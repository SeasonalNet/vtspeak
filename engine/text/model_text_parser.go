package text

import (
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/tree3"
)

const paul2013ModelContextCountOffset = 0x429a2
const paul2013ModelFinalizerTPPOutputOffset = 0x39ec

// Paul2013ModelParserFinalizerRunner binds the model parser callbacks and the
// caller arena needed by FUN_1003e240. Model and StateArena are updated after
// a successful callback; the original vendor model files are never modified.
type Paul2013ModelParserFinalizerRunner struct {
	StateArena        []byte
	ParserStateOffset int
	Model             []byte
	Callbacks         Paul2013ModelTextParserCallbacks
}

// Paul2013ModelParserFinalizerRun records both inner-wrapper and finalizer
// results. The finalizer is skipped when FUN_1000e2f0 returns a negative value.
type Paul2013ModelParserFinalizerRun struct {
	StateArena       []byte
	Model            []byte
	TextParser       Paul2013ModelTextParserResult
	OutputGroupCount int
	ReturnValue      int16
}

// Paul2013ModelTextParserCallbacks supplies the routines called by
// FUN_1000e2f0. They receive writable model/output buffers because the native
// helpers mutate their caller-owned structures. The classifier has a Go
// implementation and may be overridden for controlled comparisons.
type Paul2013ModelTextParserCallbacks struct {
	ParsePhoneRows func(model, source []byte, sourceLength int, mode int16) (int32, error)
	// NormalizeRows overrides the built-in FUN_1000ea20 row projection.
	NormalizeRows func(model, source []byte, sourceLength int) error
	// ParserOffsetRows supplies the caller-arena 0x94-byte rows used by the
	// built-in projection when NormalizeRows is nil. When omitted, the parser
	// derives the bounded ordinary-word/numeric rows for one sentence segment.
	ParserOffsetRows []byte
	// PronunciationTree supplies engbi.tree3 for the built-in ordinary
	// multi-alternative path scorer. SelectMultiAlternative overrides it.
	PronunciationTree *tree3.Tree
	// SelectMultiAlternative overrides the built-in FUN_100068b0 path scorer.
	SelectMultiAlternative Paul2013ModelMultiAlternativePhoneSelector
	ProcessRows            func(model, source []byte, mode int16) error
	ClassifyContextRow     func(contextTable []byte, selector int) (uint16, error)
	ParseTPPRows           func(source []byte, sourceLength int, output []byte) (int32, error)
}

// Paul2013ModelTextParserResult preserves the mutated buffers, the phone-row
// parser's raw result, and FUN_1000e2f0's own status. A negative TPP result is
// normalized to -1 by the native wrapper, regardless of its original value.
type Paul2013ModelTextParserResult struct {
	Model        []byte
	TPPOutput    []byte
	ParserResult int32
	ReturnValue  int32
}

// RunPaul2013ModelTextParser ports FUN_1000e2f0's counter initialization,
// helper ordering, conditional calls, and error returns. For empty input it
// clears only the model context-row count and returns zero. For nonempty input
// it resets the leading model word and context count, calls the phone-row
// parser, conditionally normalizes when model word zero becomes positive,
// conditionally processes rows when mode is nonzero, always updates the
// context-table flag, and finally parses TPP rows. Normalization uses the
// built-in FUN_1000ea20 row projection when ParserOffsetRows is supplied or
// can be derived for a single ordinary sentence segment;
// multi-alternative rows use the built-in tree scorer when PronunciationTree
// is supplied, or the SelectMultiAlternative override. The phone-row parser,
// row processor, and TPP parser remain callback inputs.
func RunPaul2013ModelTextParser(
	model []byte,
	source []byte,
	sourceLength int,
	mode int16,
	tppOutput []byte,
	callbacks Paul2013ModelTextParserCallbacks,
) (Paul2013ModelTextParserResult, error) {
	if sourceLength < 0 || sourceLength > len(source) {
		return Paul2013ModelTextParserResult{}, fmt.Errorf(
			"source length %d is outside available byte range [0, %d]", sourceLength, len(source),
		)
	}
	minimumModelSize := paul2013ModelContextCountOffset + 2
	if sourceLength > 0 {
		minimumModelSize = paul2013ModelContextCountOffset + 4
	}
	if len(model) < minimumModelSize {
		return Paul2013ModelTextParserResult{}, fmt.Errorf(
			"model buffer has %d bytes, need at least %d for context row count", len(model), minimumModelSize,
		)
	}
	result := Paul2013ModelTextParserResult{
		Model:     append([]byte(nil), model...),
		TPPOutput: append([]byte(nil), tppOutput...),
	}
	if sourceLength == 0 {
		binary.LittleEndian.PutUint16(result.Model[paul2013ModelContextCountOffset:], 0)
		return result, nil
	}
	if len(result.Model) < 2 {
		return Paul2013ModelTextParserResult{}, errors.New("model buffer has no leading parser counter")
	}
	binary.LittleEndian.PutUint16(result.Model[:2], 0)
	binary.LittleEndian.PutUint16(result.Model[paul2013ModelContextCountOffset:], 0)
	if callbacks.ParsePhoneRows == nil {
		return result, errors.New("FUN_1000e2f0 requires a phone-row parser callback")
	}
	parserResult, err := callbacks.ParsePhoneRows(result.Model, source[:sourceLength], sourceLength, mode)
	if err != nil {
		return result, fmt.Errorf("parse model phone rows: %w", err)
	}
	result.ParserResult = parserResult
	if parserResult < 0 {
		result.ReturnValue = -1
		return result, nil
	}
	if int16(binary.LittleEndian.Uint16(result.Model[:2])) > 0 {
		if callbacks.NormalizeRows != nil {
			if err := callbacks.NormalizeRows(result.Model, source[:sourceLength], sourceLength); err != nil {
				return result, fmt.Errorf("normalize model phone rows: %w", err)
			}
		} else {
			parserOffsetRows := callbacks.ParserOffsetRows
			if parserOffsetRows == nil {
				parserOffsetRows, err = BuildPaul2013OrdinaryParserOffsetRowArena(source[:sourceLength])
				if err != nil {
					return result, fmt.Errorf("derive parser offset rows: %w", err)
				}
			}
			selectMultiAlternative := callbacks.SelectMultiAlternative
			if selectMultiAlternative == nil && callbacks.PronunciationTree != nil {
				selectMultiAlternative = SelectPaul2013ModelPronunciationAlternative(callbacks.PronunciationTree)
			}
			normalized, err := NormalizePaul2013ModelPhoneRows(
				result.Model,
				parserOffsetRows,
				sourceLength,
				selectMultiAlternative,
			)
			if err != nil {
				return result, fmt.Errorf("normalize model phone rows: %w", err)
			}
			result.Model = normalized
		}
	}
	if mode != 0 {
		if callbacks.ProcessRows == nil {
			return result, errors.New("FUN_1000e2f0 requires row processing when mode is nonzero")
		}
		if err := callbacks.ProcessRows(result.Model, source[:sourceLength], mode); err != nil {
			return result, fmt.Errorf("process model phone rows: %w", err)
		}
	}
	if err := UpdatePaul2013ModelContextTableFlag(result.Model, callbacks.ClassifyContextRow); err != nil {
		return result, fmt.Errorf("update model context-table flag: %w", err)
	}
	if callbacks.ParseTPPRows == nil {
		return result, errors.New("FUN_1000e2f0 requires TPP row parsing")
	}
	tppResult, err := callbacks.ParseTPPRows(source[:sourceLength], sourceLength, result.TPPOutput)
	if err != nil {
		return result, fmt.Errorf("parse model TPP rows: %w", err)
	}
	if tppResult < 0 {
		result.ReturnValue = -1
		return result, nil
	}
	result.ReturnValue = parserResult
	return result, nil
}

// UpdatePaul2013ModelContextTableFlag ports FUN_1000e990. It checks the final
// row's signed-short field at table offset count*0x70-0x68, then scans all but
// the final row through FUN_1000e570 in native order. selector 0 represents
// its null second argument; positive selectors are the native i+1 pointer
// values. A non-nil classifier overrides the built-in FUN_1000e570 port for
// controlled comparisons. The result short is written at model+0x429a4.
func UpdatePaul2013ModelContextTableFlag(
	model []byte,
	classify func(contextTable []byte, selector int) (uint16, error),
) error {
	if len(model) < paul2013ModelContextCountOffset+4 {
		return fmt.Errorf("model buffer has %d bytes, need %d for context count and result words", len(model), paul2013ModelContextCountOffset+4)
	}
	table := model[paul2013ModelContextCountOffset:]
	count := int(int16(binary.LittleEndian.Uint16(table[:2])))
	flag := uint16(0)
	if count > 0 {
		tableBytes := paul2013FinalizerContextPrefix + count*paul2013FinalizerContextStride
		if len(table) < tableBytes {
			return fmt.Errorf("context table has %d bytes, need %d for %d rows", len(table), tableBytes, count)
		}
		tailField := count*paul2013FinalizerContextStride - 0x68
		if binary.LittleEndian.Uint16(table[tailField:tailField+2]) == 3 && count > 1 {
			value := uint16(0x44)
			for index := 0; index < count-1; index++ {
				callClassifier := index == 0 || binary.LittleEndian.Uint16(table[index*paul2013FinalizerContextStride+8:]) != 0
				if callClassifier {
					selector := 0
					if index != 0 {
						selector = index + 1
					}
					var classified uint16
					var err error
					if classify == nil {
						classified, err = ClassifyPaul2013ModelContextRow(table, selector)
					} else {
						classified, err = classify(table, selector)
					}
					if err != nil {
						return fmt.Errorf("classify context row selector %d: %w", selector, err)
					}
					value = classified
				}
				if value == 'P' {
					flag = 1
					break
				}
				if value == 'N' {
					break
				}
			}
		}
	}
	binary.LittleEndian.PutUint16(model[paul2013ModelContextCountOffset+2:], flag)
	return nil
}

// RunPaul2013ModelParserFinalizer composes the call sequence in FUN_1003e240:
// it reads the signed source length at parser-state +0, source bytes at +0x14,
// the model context table at model +0x429a2, and the TPP output buffer at
// parser-state +0x39ec. FUN_1000e2f0 runs in mode 1. A negative inner result
// returns -1 without applying the context-row finalizer; otherwise the
// FUN_1003e240 row pass runs and its low-short result is returned.
func RunPaul2013ModelParserFinalizer(
	stateArena []byte,
	parserStateOffset int,
	model []byte,
	callbacks Paul2013ModelTextParserCallbacks,
) (Paul2013ModelParserFinalizerRun, error) {
	if parserStateOffset < 0 || parserStateOffset > len(stateArena) ||
		paul2013ModelFinalizerTPPOutputOffset > len(stateArena)-parserStateOffset {
		return Paul2013ModelParserFinalizerRun{}, fmt.Errorf(
			"parser state offset %#x does not contain the TPP output buffer at +%#x",
			parserStateOffset, paul2013ModelFinalizerTPPOutputOffset,
		)
	}
	if len(stateArena)-parserStateOffset < 2 {
		return Paul2013ModelParserFinalizerRun{}, errors.New("parser state has no signed source-length field")
	}
	workingState := append([]byte(nil), stateArena...)
	workingModel := append([]byte(nil), model...)
	sourceLength := int(int16(binary.LittleEndian.Uint16(workingState[parserStateOffset:])))
	if sourceLength < 0 {
		return Paul2013ModelParserFinalizerRun{}, fmt.Errorf("negative native source length %d is unsupported", sourceLength)
	}
	sourceStart := parserStateOffset + 0x14
	if sourceStart > len(workingState) || sourceLength > len(workingState)-sourceStart {
		return Paul2013ModelParserFinalizerRun{}, fmt.Errorf(
			"source range [%#x, %#x) exceeds parser arena size %d",
			sourceStart, sourceStart+sourceLength, len(workingState),
		)
	}
	minimumModelSize := paul2013ModelContextCountOffset + 2
	if sourceLength > 0 {
		minimumModelSize = paul2013ModelContextCountOffset + 4
	}
	if len(workingModel) < minimumModelSize {
		return Paul2013ModelParserFinalizerRun{}, fmt.Errorf(
			"model buffer has %d bytes, need at least %d for this source length",
			len(workingModel), minimumModelSize,
		)
	}
	tppStart := parserStateOffset + paul2013ModelFinalizerTPPOutputOffset
	if sourceLength > 0 && tppStart >= len(workingState) {
		return Paul2013ModelParserFinalizerRun{}, errors.New("parser arena has no writable TPP output byte")
	}
	textResult, err := RunPaul2013ModelTextParser(
		workingModel,
		workingState[sourceStart:sourceStart+sourceLength],
		sourceLength,
		1,
		workingState[tppStart:],
		callbacks,
	)
	if err != nil {
		return Paul2013ModelParserFinalizerRun{}, err
	}
	copy(workingState[tppStart:], textResult.TPPOutput)
	workingModel = textResult.Model
	if textResult.ReturnValue < 0 {
		return Paul2013ModelParserFinalizerRun{
			StateArena:  workingState,
			Model:       workingModel,
			TextParser:  textResult,
			ReturnValue: -1,
		}, nil
	}
	contextCount := int(int16(binary.LittleEndian.Uint16(workingModel[paul2013ModelContextCountOffset:])))
	modelResult := int16(0)
	if contextCount > 0 {
		modelResult = int16(binary.LittleEndian.Uint16(workingModel[paul2013ModelContextCountOffset+2:]))
	}
	finalized, err := FinalizePaul2013ModelParserRows(Paul2013ModelParserFinalizerInput{
		StateArena:        workingState,
		ParserStateOffset: parserStateOffset,
		ContextTable:      workingModel[paul2013ModelContextCountOffset:],
		ContextRowCount:   contextCount,
		ModelResult:       modelResult,
	})
	if err != nil {
		return Paul2013ModelParserFinalizerRun{}, err
	}
	return Paul2013ModelParserFinalizerRun{
		StateArena:       finalized.StateArena,
		Model:            workingModel,
		TextParser:       textResult,
		OutputGroupCount: finalized.OutputGroupCount,
		ReturnValue:      finalized.ReturnValue,
	}, nil
}

// Callback adapts RunPaul2013ModelParserFinalizer to the dispatcher finalizer
// signature. parserResult is supplied by FUN_1003d3d0/FUN_1003e070 but the
// native FUN_1003e240 entry reads its own source count and calls
// FUN_1000e2f0, so that outer callback result is not consumed here.
func (runner *Paul2013ModelParserFinalizerRunner) Callback(
	state []byte,
	parserResult int32,
) (int16, error) {
	_ = parserResult
	if runner == nil {
		return 0, errors.New("nil Paul 2013 model parser finalizer runner")
	}
	if runner.ParserStateOffset < 0 || runner.ParserStateOffset > len(runner.StateArena) ||
		len(state) > len(runner.StateArena)-runner.ParserStateOffset {
		return 0, fmt.Errorf("parser state range exceeds caller arena size %d", len(runner.StateArena))
	}
	copy(runner.StateArena[runner.ParserStateOffset:], state)
	run, err := RunPaul2013ModelParserFinalizer(
		runner.StateArena,
		runner.ParserStateOffset,
		runner.Model,
		runner.Callbacks,
	)
	if err != nil {
		return 0, err
	}
	runner.StateArena = run.StateArena
	runner.Model = run.Model
	copy(state, runner.StateArena[runner.ParserStateOffset:runner.ParserStateOffset+len(state)])
	return run.ReturnValue, nil
}
