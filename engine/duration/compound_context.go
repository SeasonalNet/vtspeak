package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013CompoundOutputLimit    = 63
	paul2013CompoundComponentLimit = 31
)

// Paul2013CompoundComponentInput describes one delimiter-separated component
// passed by FUN_10009dc0 to its component-normalization chain.
type Paul2013CompoundComponentInput struct {
	Index               int
	Surface             []byte
	CurrentFlags        byte
	SeparatorBefore     byte
	SeparatorAfter      byte
	UseHyphenContext    bool
	DictionaryGate      Paul2013CompoundDictionaryGateResult
	CharacterGatePassed bool
}

// Paul2013CompoundComponentResult is the explicit result of the component
// normalizer callback. ContextFallback requests FUN_10009cd0 encoding when
// the supplied native-handler chain misses; a nil/zero result without that
// bit leaves the component empty.
type Paul2013CompoundComponentResult struct {
	Output          []byte
	RowFlagMask     byte
	ContextFallback bool
}

// Paul2013CompoundComponentNormalizer supplies the per-component path through
// dictionary, suffix, and special handlers. The scanner and aggregate output
// behavior are implemented here; unported component handlers remain the
// callback's responsibility.
type Paul2013CompoundComponentNormalizer func(
	context.Context,
	Paul2013CompoundComponentInput,
) (Paul2013CompoundComponentResult, error)

// Paul2013CompoundContractionPrefixNormalizer receives FUN_1000c040's prefix
// and native hyphen-context argument for the current compound component.
type Paul2013CompoundContractionPrefixNormalizer func(
	context.Context,
	[]byte,
	bool,
) (Paul2013ContractionBaseResult, error)

// Paul2013CompoundContextResult reports FUN_10009dc0's supported compound
// scanning and output assembly. Output omits the trailing NUL.
type Paul2013CompoundContextResult struct {
	Processed                    bool
	Matched                      bool
	Source                       []byte
	Components                   []Paul2013CompoundComponentInput
	Output                       []byte
	RowFlags                     byte
	OutputOverflow               bool
	PriorOutputMarkerRows        []int
	MarkerWithoutPriorOutputRows []int
}

// NormalizePaul2013CompoundContext ports FUN_10009dc0's compound scanner and
// result assembly for nonempty, bounded components separated by hyphens or
// periods. It preserves component order, the native special one-letter `a`
// marker, `d` separators between emitted component outputs, the 63-byte
// output bound, and final row flag bit 3. Component normalization is supplied
// because the dictionary/TPP and special-handler chain is not complete.
func (engine *Engine) NormalizePaul2013CompoundContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	outerMode bool,
	normalize Paul2013CompoundComponentNormalizer,
) (Paul2013CompoundContextResult, error) {
	result := Paul2013CompoundContextResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if ctx == nil {
		return Paul2013CompoundContextResult{}, errors.New("compound context normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) == 0 {
		return result, nil
	}
	firstSeparator := bytes.IndexByte(source, '-')
	if firstSeparator < 0 {
		firstSeparator = bytes.IndexByte(source, '.')
	}
	if firstSeparator <= 0 {
		return result, nil
	}
	if firstSeparator+1 >= len(source) {
		return result, nil
	}
	result.Processed = true

	components, err := splitPaul2013CompoundComponents(source)
	if err != nil {
		return Paul2013CompoundContextResult{}, err
	}
	result.Components = components
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	markerPending := false
	for _, component := range components {
		if err := ctx.Err(); err != nil {
			return Paul2013CompoundContextResult{}, err
		}
		if len(component.Surface) > paul2013CompoundComponentLimit {
			return Paul2013CompoundContextResult{}, fmt.Errorf(
				"compound component %d has %d bytes; native local capacity is %d plus NUL",
				component.Index, len(component.Surface), paul2013CompoundComponentLimit,
			)
		}
		component.UseHyphenContext = (component.Index == 0 && outerMode) ||
			(component.Index > 0 && component.SeparatorBefore == '-')
		result.Components[component.Index].UseHyphenContext = component.UseHyphenContext

		output := []byte(nil)
		if len(component.Surface) == 1 {
			markerPending = false
			if characterMap[component.Surface[0]] == 'a' &&
				((component.Index == 0 && outerMode && component.SeparatorAfter == '-') ||
					(component.Index > 0 && component.SeparatorBefore == '-' && component.SeparatorAfter == '-')) {
				output = []byte{0x07}
			} else {
				output, err = text.EncodePaul2013ContextString(component.Surface)
				if err != nil {
					return Paul2013CompoundContextResult{}, fmt.Errorf(
						"encode one-byte compound component %d: %w", component.Index, err,
					)
				}
				if len(output) != 0 {
					result.RowFlags |= 0x20
				}
			}
		} else {
			if normalize == nil {
				return Paul2013CompoundContextResult{}, errors.New("compound context normalizer has no multi-byte component handler")
			}
			component.DictionaryGate, err = engine.LookupPaul2013CompoundDictionaryGate(ctx, component.Surface)
			if err != nil {
				return Paul2013CompoundContextResult{}, fmt.Errorf("read compound component %d dictionary gate: %w", component.Index, err)
			}
			component.CharacterGatePassed, err = Paul2013CompoundCharacterGate(component.Surface)
			if err != nil {
				return Paul2013CompoundContextResult{}, fmt.Errorf("read compound component %d character gate: %w", component.Index, err)
			}
			result.Components[component.Index] = component
			component.CurrentFlags = result.RowFlags
			result.Components[component.Index].CurrentFlags = component.CurrentFlags
			componentResult, err := normalize(ctx, component)
			if err != nil {
				return Paul2013CompoundContextResult{}, fmt.Errorf(
					"normalize compound component %d %q: %w", component.Index, component.Surface, err,
				)
			}
			if componentResult.RowFlagMask & ^byte(0x31) != 0 {
				return Paul2013CompoundContextResult{}, fmt.Errorf(
					"compound component %d returned unsupported row flags %#02x",
					component.Index, componentResult.RowFlagMask,
				)
			}
			result.RowFlags |= componentResult.RowFlagMask
			output = append([]byte(nil), componentResult.Output...)
			if componentResult.ContextFallback && len(output) == 0 {
				output, err = text.EncodePaul2013ContextString(component.Surface)
				if err != nil {
					return Paul2013CompoundContextResult{}, fmt.Errorf(
						"encode compound component %d context fallback: %w", component.Index, err,
					)
				}
				if len(output) != 0 {
					result.RowFlags |= 0x20
				}
			}
		}
		if bytes.IndexByte(output, 0) >= 0 {
			return Paul2013CompoundContextResult{}, fmt.Errorf("compound component %d output contains NUL", component.Index)
		}
		if len(component.Surface) > 1 {
			if markerPending && compoundCodeClassNeedsPriorMarker(output) {
				if len(result.Output) == 0 {
					result.MarkerWithoutPriorOutputRows = append(result.MarkerWithoutPriorOutputRows, component.Index)
				} else {
					result.Output[len(result.Output)-1] = '&'
					result.PriorOutputMarkerRows = append(result.PriorOutputMarkerRows, component.Index)
				}
			}
			markerPending = text.Paul2013MappedCStringEqual(component.Surface, []byte("the"), characterMap) ||
				(component.Index > 0 && text.Paul2013MappedCStringEqual(component.Surface, []byte("de"), characterMap))
		}
		if len(output) == 0 {
			continue
		}
		if len(result.Output)+len(output)+1 > paul2013CompoundOutputLimit {
			hadOutput := len(result.Output) != 0
			result.OutputOverflow = true
			if hadOutput {
				result.RowFlags |= 0x08
				result.Matched = true
			}
			return result, nil
		}
		if len(result.Output) > 0 {
			result.Output = append(result.Output, 'd')
		}
		result.Output = append(result.Output, output...)
	}
	if len(result.Output) != 0 {
		result.RowFlags |= 0x08
		result.Matched = true
	}
	return result, nil
}

func compoundCodeClassNeedsPriorMarker(output []byte) bool {
	if len(output) == 0 {
		return false
	}
	first := output[0]
	return first > 0 && first < 'F' && first != 'C' && text.IsPaul2013ContextCodeClass(first)
}

// NormalizePaul2013CompoundContextWithEmbeddedDictionary composes the
// recovered FUN_10009dc0 gates and FUN_1000cb30 branch with a caller-supplied
// remainder of the component fallback chain. Components that pass both
// native gates go directly to FUN_10009cd0; otherwise a native return of 1
// from FUN_1000cb30 terminates the chain and sets row flag bit 0. The
// remaining A140, apostrophe, contraction, C3A0, and generic order must be
// supplied by fallback and is never skipped implicitly.
func (engine *Engine) NormalizePaul2013CompoundContextWithEmbeddedDictionary(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	outerMode bool,
	fallback Paul2013CompoundComponentNormalizer,
) (Paul2013CompoundContextResult, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013CompoundContextResult{}, errors.New("Paul 2013 compound context dictionary is nil or unloaded")
	}
	if fallback == nil {
		return Paul2013CompoundContextResult{}, errors.New("compound context composition has no remaining component fallback chain")
	}
	normalize := func(
		ctx context.Context,
		component Paul2013CompoundComponentInput,
	) (Paul2013CompoundComponentResult, error) {
		if component.DictionaryGate.Eligible && component.CharacterGatePassed {
			return Paul2013CompoundComponentResult{ContextFallback: true}, nil
		}
		writes, err := engine.NormalizePaul2013ContextTokenToBuffer(
			ctx,
			component.Surface,
			engine.LookupPaul2013EmbeddedContextToken,
			make([]byte, 68),
		)
		if err != nil {
			return Paul2013CompoundComponentResult{}, fmt.Errorf(
				"normalize compound component %d through FUN_1000cb30: %w",
				component.Index, err,
			)
		}
		if writes.Result.NativeReturnCode == 1 {
			end := bytes.IndexByte(writes.Output, 0)
			if end < 0 {
				return Paul2013CompoundComponentResult{}, fmt.Errorf(
					"compound component %d FUN_1000cb30 output is not NUL-terminated",
					component.Index,
				)
			}
			return Paul2013CompoundComponentResult{
				Output: append([]byte(nil), writes.Output[:end]...), RowFlagMask: 0x01,
			}, nil
		}
		return fallback(ctx, component)
	}
	return engine.NormalizePaul2013CompoundContext(
		ctx, source, currentFlags, outerMode, normalize,
	)
}

// NormalizePaul2013CompoundContextWithSupportedHandlers composes the
// dictionary/character gates and FUN_1000cb30 lookup with the recovered A140,
// trailing-apostrophe, contraction-suffix, C3A0, and generic handlers in
// their FUN_10009dc0 order. Prefix normalization for FUN_1000c040 remains an
// explicit callback. Prior-output `&` writes are applied when an accumulated
// output byte exists; cases that would write before an empty output buffer are
// listed in MarkerWithoutPriorOutputRows.
func (engine *Engine) NormalizePaul2013CompoundContextWithSupportedHandlers(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	outerMode bool,
	model []byte,
	contextRowIndex int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013CompoundContextResult, error) {
	fallback := func(
		ctx context.Context,
		component Paul2013CompoundComponentInput,
	) (Paul2013CompoundComponentResult, error) {
		cascade, err := engine.NormalizePaul2013SupportedA140Suffixes(
			ctx, component.Surface, component.CurrentFlags,
		)
		if err != nil {
			return Paul2013CompoundComponentResult{}, fmt.Errorf("run compound component %d A140 handlers: %w", component.Index, err)
		}
		if cascade.UnsupportedBranch || cascade.TableClassUnsupported {
			return Paul2013CompoundComponentResult{}, fmt.Errorf("compound component %d selected an unresolved A140 %q branch", component.Index, cascade.Handler)
		}
		if cascade.Matched {
			return compoundComponentResultFromFlags(cascade.Output, cascade.RowFlags, component.CurrentFlags), nil
		}

		apostrophe, err := engine.NormalizePaul2013TrailingApostropheFromEmbeddedDictionary(
			ctx, component.Surface, make([]byte, 68), component.CurrentFlags,
		)
		if err != nil {
			return Paul2013CompoundComponentResult{}, fmt.Errorf("run compound component %d trailing-apostrophe handler: %w", component.Index, err)
		}
		if apostrophe.Matched {
			end := bytes.IndexByte(apostrophe.Output, 0)
			if end < 0 {
				return Paul2013CompoundComponentResult{}, fmt.Errorf("compound component %d trailing-apostrophe output is not NUL-terminated", component.Index)
			}
			return compoundComponentResultFromFlags(apostrophe.Output[:end], apostrophe.RowFlags, component.CurrentFlags), nil
		}

		var contraction Paul2013ContractionSuffixResult
		if contractionPrefix == nil {
			contraction, err = engine.NormalizePaul2013ContractionSuffixWithSupportedPrefix(
				ctx, component.Surface, component.CurrentFlags,
				component.UseHyphenContext, model, contextRowIndex,
			)
		} else {
			prefixNormalizer := func(ctx context.Context, prefix []byte) (Paul2013ContractionBaseResult, error) {
				return contractionPrefix(ctx, prefix, component.UseHyphenContext)
			}
			contraction, err = NormalizePaul2013ContractionSuffixContext(
				ctx, component.Surface, component.CurrentFlags, prefixNormalizer,
			)
		}
		if err != nil {
			return Paul2013CompoundComponentResult{}, fmt.Errorf("run compound component %d contraction handler: %w", component.Index, err)
		}
		if contraction.Matched {
			return compoundComponentResultFromFlags(contraction.Output, contraction.RowFlags, component.CurrentFlags), nil
		}

		c3a0, err := engine.NormalizePaul2013C3A0ComponentSequenceWithEmbeddedFUN100086C0(
			ctx, component.Surface, component.CurrentFlags, model, contextRowIndex, nil,
		)
		if err != nil {
			return Paul2013CompoundComponentResult{}, fmt.Errorf("run compound component %d C3A0 handler: %w", component.Index, err)
		}
		if c3a0.Matched {
			return compoundComponentResultFromFlags(c3a0.Output, c3a0.RowFlags, component.CurrentFlags), nil
		}

		generic, err := engine.NormalizePaul2013GenericToken(ctx, component.Surface)
		if err != nil {
			return Paul2013CompoundComponentResult{}, fmt.Errorf("run compound component %d generic handler: %w", component.Index, err)
		}
		if generic.Eligible {
			output := generic.ClassCodes
			if generic.UsedContextFallback {
				output = generic.ContextCodes
			}
			return Paul2013CompoundComponentResult{
				Output:      append([]byte(nil), output...),
				RowFlagMask: generic.RowFlagMask &^ 0x08,
			}, nil
		}
		return Paul2013CompoundComponentResult{ContextFallback: true}, nil
	}
	return engine.NormalizePaul2013CompoundContextWithEmbeddedDictionary(
		ctx, source, currentFlags, outerMode, fallback,
	)
}

func compoundComponentResultFromFlags(output []byte, flags, currentFlags byte) Paul2013CompoundComponentResult {
	return Paul2013CompoundComponentResult{
		Output: append([]byte(nil), output...), RowFlagMask: (flags &^ currentFlags) &^ 0x08,
	}
}

func splitPaul2013CompoundComponents(source []byte) ([]Paul2013CompoundComponentInput, error) {
	components := make([]Paul2013CompoundComponentInput, 0, bytes.Count(source, []byte{'-'})+bytes.Count(source, []byte{'.'})+1)
	start := 0
	separatorBefore := byte(0)
	for index := 0; index <= len(source); index++ {
		if index < len(source) && source[index] != '-' && source[index] != '.' {
			continue
		}
		end := index
		if end-start > paul2013CompoundComponentLimit {
			return nil, fmt.Errorf("compound component at byte %d has %d bytes; native local capacity is %d plus NUL", start, end-start, paul2013CompoundComponentLimit)
		}
		component := Paul2013CompoundComponentInput{
			Index: len(components), Surface: append([]byte(nil), source[start:end]...),
			SeparatorBefore: separatorBefore,
		}
		if index < len(source) {
			component.SeparatorAfter = source[index]
		}
		components = append(components, component)
		if index == len(source) {
			break
		}
		separatorBefore = source[index]
		start = index + 1
	}
	return components, nil
}
