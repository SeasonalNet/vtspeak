package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// NormalizePaul2013C3A0EmbeddedRecordComponent ports the direct
// FUN_1000c3a0 parsed-record branch for components of at least two bytes with
// pronunciation rows. Payload bit 7 gates a separate generic call and does
// not select the parsed string.
func (engine *Engine) NormalizePaul2013C3A0EmbeddedRecordComponent(
	ctx context.Context,
	input Paul2013C3A0ComponentInput,
) (Paul2013C3A0ComponentResult, error) {
	if engine == nil {
		return Paul2013C3A0ComponentResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013C3A0ComponentResult{}, errors.New("FUN_1000c3a0 component lookup has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013C3A0ComponentResult{}, err
	}
	lookup, err := engine.LookupPaul2013EmbeddedContextToken(ctx, input.SourceSpan.Output)
	if err != nil {
		return Paul2013C3A0ComponentResult{}, fmt.Errorf("look up component record: %w", err)
	}
	if len(input.SourceSpan.Output) < 2 || lookup.RecordCount < 1 {
		return Paul2013C3A0ComponentResult{}, nil
	}
	if bytes.IndexByte(lookup.DictionaryText, 0) >= 0 {
		return Paul2013C3A0ComponentResult{}, errors.New("component record text contains NUL")
	}
	return Paul2013C3A0ComponentResult{
		Output:      append([]byte(nil), lookup.DictionaryText...),
		RowFlagMask: 0x01,
		Handled:     true,
	}, nil
}

// NormalizePaul2013C3A0ComponentSequenceWithEmbeddedDictionary runs the
// direct embedded-record branch before the caller's remaining native handler
// chain, then connects both to the recovered source gate and sequence shell.
func (engine *Engine) NormalizePaul2013C3A0ComponentSequenceWithEmbeddedDictionary(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	fallback Paul2013C3A0ComponentHandler,
) (Paul2013C3A0SequenceResult, error) {
	if engine == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	handler := func(
		ctx context.Context,
		input Paul2013C3A0ComponentInput,
	) (Paul2013C3A0ComponentResult, error) {
		result, err := engine.NormalizePaul2013C3A0EmbeddedRecordComponent(ctx, input)
		if err != nil || result.Handled {
			return result, err
		}
		if fallback == nil {
			return Paul2013C3A0ComponentResult{}, nil
		}
		return fallback(ctx, input)
	}
	return NormalizePaul2013C3A0ComponentSequenceFromSource(ctx, source, currentFlags, handler)
}

// Paul2013C3A0OrderedHandlers divides the still-caller-supplied component
// paths around the directly ported apostrophe-s handler. BeforeApostropheS
// represents the earlier exception and FUN_1000a140 branches;
// AfterApostropheS represents the later generic/context-encoder fallback.
type Paul2013C3A0OrderedHandlers struct {
	BeforeApostropheS Paul2013C3A0ComponentHandler
	AfterApostropheS  Paul2013C3A0ComponentHandler
}

// Paul2013C3A0AlternateRecordProvider supplies the raw local_25c record passed
// as param_4 to FUN_100086c0 for one scanned component. Its native offsets must
// be preserved; the dictionary-record-to-local_25c projection remains a
// caller-owned producer.
type Paul2013C3A0AlternateRecordProvider func(
	context.Context,
	Paul2013C3A0ComponentInput,
) ([]byte, error)

// NormalizePaul2013C3A0ComponentSequenceWithSupportedHandlers composes the
// parsed-record branches, the caller-supplied earlier handlers, the recovered
// apostrophe-s rewrite, and the supported generic/context fallback in native
// order. Callers can supply a later handler to replace that final fallback.
func (engine *Engine) NormalizePaul2013C3A0ComponentSequenceWithSupportedHandlers(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	handlers Paul2013C3A0OrderedHandlers,
) (Paul2013C3A0SequenceResult, error) {
	if engine == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	handle := func(
		ctx context.Context,
		input Paul2013C3A0ComponentInput,
	) (Paul2013C3A0ComponentResult, error) {
		recordResult, err := engine.NormalizePaul2013C3A0EmbeddedRecordComponent(ctx, input)
		if err != nil || recordResult.Handled {
			return recordResult, err
		}
		rowFlagMask := recordResult.RowFlagMask
		lookup, err := engine.LookupPaul2013EmbeddedContextToken(ctx, input.SourceSpan.Output)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, fmt.Errorf("look up component fallback gate: %w", err)
		}
		runEarlierHandlers := true
		if lookup.GenericFallbackGate {
			precheck, err := engine.NormalizePaul2013GenericToken(ctx, input.SourceSpan.Output)
			if err != nil {
				return Paul2013C3A0ComponentResult{}, fmt.Errorf("run bit-7-gated generic precheck: %w", err)
			}
			rowFlagMask |= precheck.RowFlagMask
			if precheck.Eligible {
				runEarlierHandlers = false
			}
		}
		if runEarlierHandlers && handlers.BeforeApostropheS != nil {
			earlier, err := handlers.BeforeApostropheS(ctx, input)
			if err != nil {
				return earlier, err
			}
			rowFlagMask |= earlier.RowFlagMask
			if earlier.Handled {
				earlier.RowFlagMask |= rowFlagMask
				return earlier, nil
			}
		}
		suffix, err := engine.NormalizePaul2013ApostropheSSuffixFromEmbeddedDictionary(ctx, input.SourceSpan.Output)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, fmt.Errorf("apply apostrophe-s branch: %w", err)
		}
		if suffix.Matched {
			return Paul2013C3A0ComponentResult{
				Output: suffix.Output, RowFlagMask: rowFlagMask, Handled: true,
			}, nil
		}
		laterHandler := handlers.AfterApostropheS
		if laterHandler == nil {
			laterHandler = engine.NormalizePaul2013C3A0GenericFallbackComponent
		}
		later, err := laterHandler(ctx, input)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, err
		}
		later.RowFlagMask |= rowFlagMask
		return later, nil
	}
	return NormalizePaul2013C3A0ComponentSequenceFromSource(ctx, source, currentFlags, handle)
}

// NormalizePaul2013C3A0ComponentSequenceWithSupportedA140 composes the
// caller-owned handlers before FUN_1000a140 with the recovered A140 branches,
// apostrophe-s rewrite, and generic/context fallback. earlierHandlers must
// preserve the native FUN_100086c0 decision; the function invokes A140 only
// when that callback misses. Use NormalizePaul2013C3A0ComponentSequenceWithFUN100086C0
// when the model row and alternate-record producer are available.
func (engine *Engine) NormalizePaul2013C3A0ComponentSequenceWithSupportedA140(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	earlierHandlers Paul2013C3A0ComponentHandler,
	afterApostropheS Paul2013C3A0ComponentHandler,
) (Paul2013C3A0SequenceResult, error) {
	if engine == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if earlierHandlers == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("C3A0 supported-A140 composition requires the earlier FUN_100086c0/exception handler")
	}
	beforeApostropheS := func(
		ctx context.Context,
		input Paul2013C3A0ComponentInput,
	) (Paul2013C3A0ComponentResult, error) {
		earlier, err := earlierHandlers(ctx, input)
		if err != nil || earlier.Handled {
			return earlier, err
		}
		a140, err := engine.NormalizePaul2013C3A0SupportedA140Component(ctx, input)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, err
		}
		a140.RowFlagMask |= earlier.RowFlagMask
		return a140, nil
	}
	return engine.NormalizePaul2013C3A0ComponentSequenceWithSupportedHandlers(
		ctx,
		source,
		currentFlags,
		Paul2013C3A0OrderedHandlers{
			BeforeApostropheS: beforeApostropheS,
			AfterApostropheS:  afterApostropheS,
		},
	)
}

// NormalizePaul2013C3A0ComponentSequenceWithFUN100086C0 connects the
// recovered full FUN_100086c0 path to the ordered C3A0 component handlers. It
// uses the caller's model row and raw alternate record, then applies the
// native direct context-string output when FUN_100086c0 returns nonzero. A
// zero return continues through supported A140 branches. The outer TPP/local
// record producer and runtime row parity remain caller-owned.
func (engine *Engine) NormalizePaul2013C3A0ComponentSequenceWithFUN100086C0(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	model []byte,
	contextRowIndex int,
	alternateRecord Paul2013C3A0AlternateRecordProvider,
	afterApostropheS Paul2013C3A0ComponentHandler,
) (Paul2013C3A0SequenceResult, error) {
	if engine == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if alternateRecord == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("C3A0 FUN_100086c0 composition has no alternate-record producer")
	}
	earlier := func(
		ctx context.Context,
		input Paul2013C3A0ComponentInput,
	) (Paul2013C3A0ComponentResult, error) {
		record, err := alternateRecord(ctx, input)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, fmt.Errorf("produce FUN_100086c0 alternate record: %w", err)
		}
		if len(record) == 0 {
			return Paul2013C3A0ComponentResult{}, errors.New("FUN_100086c0 alternate-record producer returned an empty record")
		}
		decision, _, err := engine.EvaluatePaul2013FUN100086C0SupportedPath(
			ctx, input.SourceSpan.Output, model, contextRowIndex, record, true,
		)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, fmt.Errorf("evaluate FUN_100086c0: %w", err)
		}
		if decision.Disposition == Paul2013FUN100086C0NeedsNeighborScan {
			return Paul2013C3A0ComponentResult{}, errors.New("FUN_100086c0 composition stopped before its low-short return")
		}
		if decision.NativeShort == 0 {
			return Paul2013C3A0ComponentResult{}, nil
		}
		output, err := text.EncodePaul2013ContextString(input.SourceSpan.Output)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, fmt.Errorf("run FUN_10009cd0 after FUN_100086c0: %w", err)
		}
		rowFlagMask := byte(0)
		if len(output) != 0 {
			rowFlagMask = 0x20
		}
		return Paul2013C3A0ComponentResult{
			Output: output, RowFlagMask: rowFlagMask, Handled: true,
		}, nil
	}
	return engine.NormalizePaul2013C3A0ComponentSequenceWithSupportedA140(
		ctx, source, currentFlags, earlier, afterApostropheS,
	)
}

// NormalizePaul2013C3A0ComponentSequenceWithEmbeddedFUN100086C0 composes the
// embedded dictionary lookup, FUN_10003c50 local-record projection, supported
// FUN_100086c0 path, and later C3A0 handlers. The model arena and context-row
// index remain caller-owned because they identify mutable per-utterance state.
func (engine *Engine) NormalizePaul2013C3A0ComponentSequenceWithEmbeddedFUN100086C0(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	model []byte,
	contextRowIndex int,
	afterApostropheS Paul2013C3A0ComponentHandler,
) (Paul2013C3A0SequenceResult, error) {
	if engine == nil {
		return Paul2013C3A0SequenceResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	producer := func(
		ctx context.Context,
		input Paul2013C3A0ComponentInput,
	) ([]byte, error) {
		return engine.BuildPaul2013C3A0AlternateRecord(ctx, input.SourceSpan.Output)
	}
	return engine.NormalizePaul2013C3A0ComponentSequenceWithFUN100086C0(
		ctx, source, currentFlags, model, contextRowIndex, producer, afterApostropheS,
	)
}

// NormalizePaul2013C3A0GenericFallbackComponent ports the final
// FUN_10002f10/FUN_10009cd0 output choice for a scanned component. The generic
// normalizer supplies class or context codes when eligible; components that
// do not pass its gate use the direct context-string encoder. Exception and
// FUN_1000a140 rewrites still precede this helper in the ordered composition.
func (engine *Engine) NormalizePaul2013C3A0GenericFallbackComponent(
	ctx context.Context,
	input Paul2013C3A0ComponentInput,
) (Paul2013C3A0ComponentResult, error) {
	if engine == nil {
		return Paul2013C3A0ComponentResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013C3A0ComponentResult{}, errors.New("C3A0 generic fallback has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013C3A0ComponentResult{}, err
	}
	source := input.SourceSpan.Output
	generic, err := engine.NormalizePaul2013GenericToken(ctx, source)
	if err != nil {
		return Paul2013C3A0ComponentResult{}, fmt.Errorf("normalize generic component: %w", err)
	}
	var output []byte
	rowFlagMask := byte(0)
	if generic.Eligible {
		if generic.UsedContextFallback {
			output = generic.ContextCodes
		} else {
			output = generic.ClassCodes
		}
		rowFlagMask = generic.RowFlagMask
	} else {
		output, err = text.EncodePaul2013ContextString(source)
		if err != nil {
			return Paul2013C3A0ComponentResult{}, fmt.Errorf("encode component context: %w", err)
		}
		if len(output) != 0 {
			rowFlagMask = 0x20
		}
	}
	if bytes.IndexByte(output, 0) >= 0 {
		return Paul2013C3A0ComponentResult{}, errors.New("C3A0 generic fallback output contains NUL")
	}
	return Paul2013C3A0ComponentResult{
		Output: append([]byte(nil), output...), RowFlagMask: rowFlagMask, Handled: true,
	}, nil
}

// NormalizePaul2013C3A0SupportedA140Component adapts the bounded supported
// FUN_1000a140 suffix cases to the C3A0 component handler contract. Invoke it
// only after the earlier exception handler misses. Inputs that select a known
// any A140 branch explicitly marked unsupported fail closed rather than
// continuing to later handlers.
func (engine *Engine) NormalizePaul2013C3A0SupportedA140Component(
	ctx context.Context,
	input Paul2013C3A0ComponentInput,
) (Paul2013C3A0ComponentResult, error) {
	if engine == nil {
		return Paul2013C3A0ComponentResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	cascade, err := engine.NormalizePaul2013SupportedA140Suffixes(
		ctx, input.SourceSpan.Output, 0,
	)
	if err != nil {
		return Paul2013C3A0ComponentResult{}, fmt.Errorf("normalize supported FUN_1000a140 case: %w", err)
	}
	if cascade.TableClassUnsupported {
		return Paul2013C3A0ComponentResult{}, fmt.Errorf("FUN_1000a140 %s suffix uses an unrecovered short-table index", cascade.Handler)
	}
	if cascade.UnsupportedBranch {
		return Paul2013C3A0ComponentResult{}, errors.New("FUN_1000a140 selected an unported branch")
	}
	if !cascade.Matched {
		return Paul2013C3A0ComponentResult{}, nil
	}
	return Paul2013C3A0ComponentResult{
		Output:      append([]byte(nil), cascade.Output...),
		RowFlagMask: cascade.RowFlags,
		Handled:     true,
	}, nil
}
