package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// BuildPaul2013ModelPhoneRowsFromParserRows connects the loaded dictionary
// and pronunciation classifier to the supported parser-row-to-model path.
// Parser construction and the FUN_10007520 cascade remain separate stages.
func (engine *Engine) BuildPaul2013ModelPhoneRowsFromParserRows(
	model []byte,
	parserRows []byte,
	sourceLength int,
	pathPhoneEnabled bool,
) (text.Paul2013ModelParserOrchestrationResult, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, errors.New("Paul 2013 model-text resources are nil or unloaded")
	}
	return text.BuildPaul2013ModelPhoneRowsFromParserRows(
		model,
		parserRows,
		engine.dictionary,
		engine.pronunciationTree,
		sourceLength,
		pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySource connects the bounded ordinary
// source-row producer to the loaded dictionary and pronunciation classifier.
// The native parser-row controls that remain unresolved must be supplied for
// every produced source row.
func (engine *Engine) BuildPaul2013ModelPhoneRowsFromOrdinarySource(
	model []byte,
	source []byte,
	rowControls []text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
) (text.Paul2013ModelParserOrchestrationResult, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, errors.New("Paul 2013 model-text resources are nil or unloaded")
	}
	return text.BuildPaul2013ModelPhoneRowsFromOrdinarySource(
		model,
		source,
		rowControls,
		engine.dictionary,
		engine.pronunciationTree,
		pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications connects
// explicit parser row-type applications to the loaded dictionary and
// pronunciation classifier. Initial row types, mode bytes, and auxiliary
// strings remain caller supplied where their native producers are unresolved.
func (engine *Engine) BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications(
	model []byte,
	source []byte,
	rowControls []text.Paul2013ParserSourceRowControls,
	applications []text.Paul2013ParserRowTypeApplication,
	pathPhoneEnabled bool,
) (text.Paul2013ModelParserOrchestrationResult, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, errors.New("Paul 2013 model-text resources are nil or unloaded")
	}
	return text.BuildPaul2013ModelPhoneRowsFromOrdinarySourceWithTypeApplications(
		model, source, rowControls, applications, engine.dictionary,
		engine.pronunciationTree, pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena connects the
// captured sentence-segment offset rows to one loaded-dictionary model pass.
// Each segment keeps its relative offset origin; parser controls are required
// for every produced source row.
func (engine *Engine) BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
) (text.Paul2013ModelParserOrchestrationResult, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, errors.New("Paul 2013 model-text resources are nil or unloaded")
	}
	return text.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
		model,
		source,
		rowControls,
		engine.dictionary,
		engine.pronunciationTree,
		pathPhoneEnabled,
	)
}

// BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications
// composes per-segment row-type helpers with the loaded dictionary and
// pronunciation classifier. Application row indexes are local to each
// captured sentence segment; mode bytes and auxiliary strings remain explicit
// in rowControls.
func (engine *Engine) BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications(
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	applications [][]text.Paul2013ParserRowTypeApplication,
	pathPhoneEnabled bool,
) (text.Paul2013ModelParserOrchestrationResult, error) {
	if engine == nil || engine.dictionary == nil || engine.pronunciationTree == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, errors.New("Paul 2013 model-text resources are nil or unloaded")
	}
	return text.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		model, source, rowControls, applications, engine.dictionary,
		engine.pronunciationTree, pathPhoneEnabled,
	)
}

// BuildAndNormalizePaul2013ModelPhoneRowsFromParserRows runs the supported
// parser-row-to-model projection, then applies FUN_10009030 in the explicit
// context-row order. Earlier FUN_10007520 gates remain caller-owned.
func (engine *Engine) BuildAndNormalizePaul2013ModelPhoneRowsFromParserRows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	sourceLength int,
	pathPhoneEnabled bool,
	contextRowIndexes []int,
) (text.Paul2013ModelParserOrchestrationResult, []int, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil,
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil, err
	}
	result, err := engine.BuildPaul2013ModelPhoneRowsFromParserRows(
		model, parserRows, sourceLength, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil, err
	}
	processing, err := text.InitializePaul2013ModelContextProcessingFlags(result.Model)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil,
			fmt.Errorf("initialize projected model context states: %w", err)
	}
	normalizedModel, handledRows, err := engine.NormalizePaul2013ModelContextRows(
		ctx, processing.Model, result.ParserRows, contextRowIndexes,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil,
			fmt.Errorf("normalize projected model context rows: %w", err)
	}
	result.Model = normalizedModel
	return result, handledRows, nil
}

// BuildAndNormalizePaul2013ModelPhoneRowsFromOrdinarySource composes the
// bounded ordinary source builder with the supported model-row normalizer.
// It requires caller-supplied parser controls and context-row order; it does
// not run the earlier exception/TPP cascade.
func (engine *Engine) BuildAndNormalizePaul2013ModelPhoneRowsFromOrdinarySource(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls []text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	contextRowIndexes []int,
) (text.Paul2013ModelParserOrchestrationResult, []int, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil,
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil, err
	}
	result, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySource(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil, err
	}
	processing, err := text.InitializePaul2013ModelContextProcessingFlags(result.Model)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil,
			fmt.Errorf("initialize projected model context states: %w", err)
	}
	normalizedModel, handledRows, err := engine.NormalizePaul2013ModelContextRows(
		ctx, processing.Model, result.ParserRows, contextRowIndexes,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, nil,
			fmt.Errorf("normalize projected model context rows: %w", err)
	}
	result.Model = normalizedModel
	return result, handledRows, nil
}

// BuildAndRunPaul2013KnownModelContextPassesFromParserRows connects the
// parser-row projection to the supported ordered model-context passes. Rows
// still requiring unported dictionary/TPP or special handlers are available
// in the returned pass result.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesFromParserRows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	sourceLength int,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromParserRows(
		model, parserRows, sourceLength, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPasses(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithTPPFromParserRows runs the
// supported parser-to-model and model-context stages before applying the
// recoverable TPP row pass, matching their relative order in FUN_1000e2f0.
// It also applies FUN_1000e990 between the model-context and TPP stages. The
// returned TPP result records F/G operations and WAB fallback writes; other
// TPP routes remain unresolved.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithTPPFromParserRows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	sourceLength int,
	pathPhoneEnabled bool,
	dictionaryIndex int,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, text.Paul2013TPPNumericTextResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	projection, passes, err := engine.BuildAndRunPaul2013KnownModelContextPassesFromParserRows(
		ctx, model, parserRows, sourceLength, pathPhoneEnabled, tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	projection, passes, tppResult, err := engine.finishPaul2013KnownModelContextPassesWithTPP(
		projection, passes, dictionaryIndex,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	return projection, passes, tppResult, nil
}

func (engine *Engine) finishPaul2013KnownModelContextPassesWithTPP(
	projection text.Paul2013ModelParserOrchestrationResult,
	passes Paul2013ModelContextPassResult,
	dictionaryIndex int,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, text.Paul2013TPPNumericTextResult, error) {
	if err := text.UpdatePaul2013ModelContextTableFlag(projection.Model, nil); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, fmt.Errorf("update model context-table flag before TPP: %w", err)
	}
	passes.Model = projection.Model
	updatedRows, tppResult, err := engine.ApplyPaul2013TPPNumericParserArena(projection.ParserRows, dictionaryIndex)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, fmt.Errorf("apply parser-row TPP pass after model-context stages: %w", err)
	}
	projection.ParserRows = updatedRows
	return projection, passes, tppResult, nil
}

// BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySource connects the
// bounded ordinary-source row producer, context projection, and supported
// ordered model-context passes. Parser controls and literal-tail strings
// remain explicit inputs; unresolved rows remain listed in the pass result.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySource(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls []text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySource(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPasses(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySource
// composes the bounded one-segment source-row producer, supported model
// context stages, FUN_1000e990, and the later F/G plus WAB TPP pass. Row
// controls remain explicit because the ordinary scanner does not produce
// them.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySource(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls []text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	dictionaryIndex int,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, text.Paul2013TPPNumericTextResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	projection, passes, err := engine.BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySource(
		ctx, model, source, rowControls, pathPhoneEnabled, tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	return engine.finishPaul2013KnownModelContextPassesWithTPP(projection, passes, dictionaryIndex)
}

// BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArena
// applies the same ordered passes after flattening the captured sentence
// segments through the shared parser/model arenas.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	dictionaryIndex int,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, text.Paul2013TPPNumericTextResult, error) {
	return engine.BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		ctx, model, source, rowControls, nil, pathPhoneEnabled, dictionaryIndex, tailSurfaces,
	)
}

// BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArenaWithTypeApplications
// runs the supported ordered context passes and numeric TPP row updates after
// applying per-segment parser row type helpers. Marker, mode, and auxiliary
// fields remain explicit inputs; nonnumeric TPP component effects remain
// unresolved.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArenaWithTypeApplications(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	applications [][]text.Paul2013ParserRowTypeApplication,
	pathPhoneEnabled bool,
	dictionaryIndex int,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, text.Paul2013TPPNumericTextResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	projection, passes, err := engine.BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		ctx, model, source, rowControls, applications, pathPhoneEnabled, tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	return engine.finishPaul2013KnownModelContextPassesWithTPP(projection, passes, dictionaryIndex)
}

// BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArenaWithDefaultParserControls
// composes the captured ordinary parser defaults, per-segment row-type
// applications, shared-arena model projection, supported ordered context
// passes, FUN_1000e990, and the direct F/G plus WAB TPP pass. Initial +0x2c
// row types and type-application inputs remain explicit because their native
// producers are not fully recovered.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithTPPFromOrdinarySegmentsInSharedArenaWithDefaultParserControls(
	ctx context.Context,
	model []byte,
	source []byte,
	initialRowTypes [][]uint32,
	applications [][]text.Paul2013ParserRowTypeApplication,
	pathPhoneEnabled bool,
	dictionaryIndex int,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, text.Paul2013TPPNumericTextResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	projection, err := text.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithDefaultParserControls(
		model,
		source,
		initialRowTypes,
		applications,
		engine.dictionary,
		engine.pronunciationTree,
		pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPasses(
		ctx,
		projection.Model,
		projection.ParserRows,
		tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, text.Paul2013TPPNumericTextResult{}, err
	}
	return engine.finishPaul2013KnownModelContextPassesWithTPP(projection, passes, dictionaryIndex)
}

// BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArena
// connects ordered multi-segment source rows to the shared model projection
// and supported counted-row passes. Parser controls and literal-tail strings
// remain explicit; unported handlers are listed in UnresolvedRows.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	return engine.BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		ctx, model, source, rowControls, nil, pathPhoneEnabled, tailSurfaces,
	)
}

// BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArenaWithTypeApplications
// runs the supported counted-row passes after applying recovered parser-row
// type helpers to each captured source segment. Applications keep local row
// indexes; mode bytes, auxiliary strings, and literal-tail values remain
// caller supplied.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesFromOrdinarySegmentsInSharedArenaWithTypeApplications(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	applications [][]text.Paul2013ParserRowTypeApplication,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArenaWithTypeApplications(
		model, source, rowControls, applications, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPasses(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithSelectedCompoundRowsFromOrdinarySegmentsInSharedArena
// connects shared-arena source projection to the known model-context passes
// and derives each selected compound row's native outer mode from its
// preceding parser row. Compound/C3A0 row eligibility and parser controls
// remain explicit inputs.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithSelectedCompoundRowsFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
	compoundContextRowIndexes, c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPassesWithSelectedCompoundRows(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
		compoundContextRowIndexes, c3a0ContextRowIndexes, contractionPrefix,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithC3A0RowsFromOrdinarySegmentsInSharedArena
// connects shared-arena segment projection to the known counted-row passes
// and caller-selected C3A0 branches. Parser controls, literal-tail strings,
// and outer C3A0 dispatch row indexes remain explicit inputs.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithC3A0RowsFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
	c3a0ContextRowIndexes []int,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPassesWithC3A0Rows(
		ctx, projection.Model, projection.ParserRows, tailSurfaces, c3a0ContextRowIndexes,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithCompoundAndC3A0RowsFromOrdinarySegmentsInSharedArena
// carries explicit outer compound modes and C3A0 row selections through the
// shared-arena segment projection and ordered counted-row pass.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithCompoundAndC3A0RowsFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
	compoundOuterModes map[int]bool,
	c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPassesWithCompoundAndC3A0Rows(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
		compoundOuterModes, c3a0ContextRowIndexes, contractionPrefix,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0RowsFromOrdinarySegmentsInSharedArena
// carries explicit compound and contraction outer modes plus C3A0 row
// selections through shared-arena segment projection and the supported
// counted-row fallback chain.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0RowsFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
	compoundOuterModes map[int]bool,
	contractionOuterModes map[int]bool,
	c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0Rows(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
		compoundOuterModes, contractionOuterModes, c3a0ContextRowIndexes,
		contractionPrefix,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}

// BuildAndRunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0RowsFromOrdinarySegmentsInSharedArena
// connects the supported B800 handler to ordinary-segment model projection.
// Compound/contraction modes and B800/C3A0 row selections remain explicit
// because their outer eligibility predicates are not all recovered.
func (engine *Engine) BuildAndRunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0RowsFromOrdinarySegmentsInSharedArena(
	ctx context.Context,
	model []byte,
	source []byte,
	rowControls [][]text.Paul2013ParserSourceRowControls,
	pathPhoneEnabled bool,
	tailSurfaces map[int][]byte,
	compoundOuterModes map[int]bool,
	contractionOuterModes map[int]bool,
	b800ContextRowIndexes []int,
	c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (text.Paul2013ModelParserOrchestrationResult, Paul2013ModelContextPassResult, error) {
	if ctx == nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{},
			errors.New("model-text orchestration has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection, err := engine.BuildPaul2013ModelPhoneRowsFromOrdinarySegmentsInSharedArena(
		model, source, rowControls, pathPhoneEnabled,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	passes, err := engine.RunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0Rows(
		ctx, projection.Model, projection.ParserRows, tailSurfaces,
		compoundOuterModes, contractionOuterModes, b800ContextRowIndexes,
		c3a0ContextRowIndexes, contractionPrefix,
	)
	if err != nil {
		return text.Paul2013ModelParserOrchestrationResult{}, Paul2013ModelContextPassResult{}, err
	}
	projection.Model = passes.Model
	return projection, passes, nil
}
