package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// GenericNormalizerResult contains the outputs of the recovered inner
// FUN_10002f10 dispatch. ClassCodes are opaque dictionary-dispatch bytes;
// ContextCodes are the direct FUN_10009cd0 encoding. Row flags are returned
// as a mask because this routine receives no native row to mutate.
type GenericNormalizerResult struct {
	Eligible            bool
	ClassCodes          []byte
	ContextCodes        []byte
	UsedContextFallback bool
	RowFlagMask         byte
}

// GenericNormalizerWrites contains the copied native output buffer and the
// separately updated byte passed as FUN_10002f10's third argument.
type GenericNormalizerWrites struct {
	Output []byte
	Flags  byte
}

// Paul2013HyphenNormalizerResult retains FUN_10003110's joined output,
// accumulated flags, native low-short result, and capacity-stop observation.
type Paul2013HyphenNormalizerResult struct {
	Codes           []byte
	Flags           byte
	NativeShort     int16
	CapacityStopped bool
}

// NormalizePaul2013HyphenatedGenericToken ports FUN_10003110's hyphen split,
// per-component FUN_10002f10 calls, d-joined output, and 64-byte local cap.
// Component normalization uses the already recovered generic-token path;
// input table behavior outside that supported path still fails through its
// existing errors.
func (engine *Engine) NormalizePaul2013HyphenatedGenericToken(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013HyphenNormalizerResult, error) {
	if engine == nil {
		return Paul2013HyphenNormalizerResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013HyphenNormalizerResult{}, errors.New("hyphenated generic normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013HyphenNormalizerResult{}, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}

	result := Paul2013HyphenNormalizerResult{Flags: currentFlags}
	nativeLength := 0
	for start := 0; start <= len(source); {
		end := start + bytes.IndexByte(source[start:], '-')
		if end < start {
			end = len(source)
		}
		component := source[start:end]
		generic, err := engine.NormalizePaul2013GenericToken(ctx, component)
		if err != nil {
			return Paul2013HyphenNormalizerResult{}, fmt.Errorf("normalize hyphen component at byte %d: %w", start, err)
		}
		if generic.Eligible {
			temporary := make([]byte, 0x44)
			writes, err := ApplyPaul2013GenericNormalizerWrites(temporary, result.Flags, generic)
			if err != nil {
				return Paul2013HyphenNormalizerResult{}, fmt.Errorf("write hyphen component at byte %d: %w", start, err)
			}
			result.Flags = writes.Flags
			componentCodes := cStringBytes(writes.Output)
			// FUN_10003110 tracks each component's terminating NUL in its
			// native count; a later separator overwrites that byte in output.
			if len(componentCodes)+nativeLength+1 >= 0x40 {
				result.CapacityStopped = true
				break
			}
			if nativeLength > 0 {
				result.Codes = append(result.Codes, 'd')
				nativeLength++
			}
			result.Codes = append(result.Codes, componentCodes...)
			nativeLength += len(componentCodes) + 1
		}
		if end == len(source) {
			break
		}
		start = end + 1
	}
	if nativeLength != 0 {
		result.NativeShort = 1
	}
	return result, nil
}

// ApplyPaul2013GenericNormalizerWrites ports FUN_10002f10's final writes.
// It copies the selected class/context code string into destination, writes a
// trailing NUL, preserves the rest of that buffer, and ORs the recovered flag
// mask into currentFlags. The caller-arena layout remains explicit.
func ApplyPaul2013GenericNormalizerWrites(
	destination []byte,
	currentFlags byte,
	result GenericNormalizerResult,
) (GenericNormalizerWrites, error) {
	if result.RowFlagMask & ^byte(0x30) != 0 {
		return GenericNormalizerWrites{}, fmt.Errorf("generic normalizer has unsupported row-flag mask %#02x", result.RowFlagMask)
	}
	if !result.Eligible {
		if len(result.ClassCodes) != 0 || len(result.ContextCodes) != 0 || result.RowFlagMask != 0 || result.UsedContextFallback {
			return GenericNormalizerWrites{}, errors.New("ineligible generic-normalizer result contains output writes")
		}
		return GenericNormalizerWrites{Output: append([]byte(nil), destination...), Flags: currentFlags}, nil
	}
	if result.UsedContextFallback {
		if len(result.ClassCodes) != 0 || result.RowFlagMask&0x10 != 0 {
			return GenericNormalizerWrites{}, errors.New("context-fallback result also contains class-dispatch writes")
		}
	} else if len(result.ClassCodes) == 0 || len(result.ContextCodes) != 0 || result.RowFlagMask != 0x10 {
		return GenericNormalizerWrites{}, errors.New("class-dispatch result has inconsistent output or row flags")
	}

	output := result.ClassCodes
	if result.UsedContextFallback {
		output = result.ContextCodes
		if len(output) == 0 && result.RowFlagMask != 0 {
			return GenericNormalizerWrites{}, errors.New("empty context fallback has nonzero row flags")
		}
		if len(output) != 0 && result.RowFlagMask != 0x20 {
			return GenericNormalizerWrites{}, errors.New("nonempty context fallback is missing its native 0x20 row flag")
		}
	}
	if bytes.IndexByte(output, 0) >= 0 {
		return GenericNormalizerWrites{}, errors.New("generic normalizer output contains an embedded NUL")
	}
	if len(output) >= len(destination) {
		return GenericNormalizerWrites{}, fmt.Errorf("generic normalizer output needs %d bytes including NUL, destination has %d", len(output)+1, len(destination))
	}
	written := append([]byte(nil), destination...)
	copy(written, output)
	written[len(output)] = 0
	return GenericNormalizerWrites{Output: written, Flags: currentFlags | result.RowFlagMask}, nil
}

// Paul2013ContextLookupResult contains the fields FUN_1000cb30 consumes from
// the FUN_10003a70/FUN_10003c50 lookup path. RecordCount is the native signed
// count at the recovered result offset; GenericFallbackGate is the nonzero
// short gate at its adjacent field. DictionaryText is copied only when the
// count is positive. The resource lookup that produces these fields remains
// caller supplied.
type Paul2013ContextLookupResult struct {
	RecordCount         int32
	GenericFallbackGate bool
	DictionaryText      []byte
}

// Paul2013ContextLookup performs the still-unported resource lookup before
// FUN_1000cb30's recovered outer dispatch.
type Paul2013ContextLookup func(context.Context, []byte) (Paul2013ContextLookupResult, error)

// Paul2013ContextNormalizationResult records the outer normalizer's selected
// branch. DictionaryText is set for a positive lookup count; Generic is set
// only when the lookup count is nonpositive and the native gate is nonzero.
type Paul2013ContextNormalizationResult struct {
	Matched        bool
	UsedDictionary bool
	DictionaryText []byte
	Generic        GenericNormalizerResult
	// NativeReturnCode preserves FUN_1000cb30's short return: 1 when its
	// dictionary or generic branch handles the token, and -1 on a miss.
	NativeReturnCode int16
}

// Paul2013ContextNormalizationWrites retains FUN_1000cb30's branch result,
// the copied output buffer, and the local generic-normalizer flag byte.
type Paul2013ContextNormalizationWrites struct {
	Result     Paul2013ContextNormalizationResult
	Output     []byte
	LocalFlags byte
}

// Paul2013TrailingApostropheResult records FUN_1000c710's supported branch.
// Path identifies which outer spelling branch returned first.
type Paul2013TrailingApostropheResult struct {
	Applied  bool
	Matched  bool
	Path     string
	Output   []byte
	RowFlags byte
}

// NormalizePaul2013TrailingApostrophe ports FUN_1000c710 for bounded ASCII
// strings ending in an apostrophe. It tries the mapped "in'" spelling with a
// case-preserving G/g substitution, then the spelling without its apostrophe,
// then the direct generic normalizer and context encoder. The caller supplies
// the context lookup used by FUN_1000cb30.
func (engine *Engine) NormalizePaul2013TrailingApostrophe(
	ctx context.Context,
	source []byte,
	lookup Paul2013ContextLookup,
	destination []byte,
	currentFlags byte,
) (Paul2013TrailingApostropheResult, error) {
	result := Paul2013TrailingApostropheResult{
		Output: append([]byte(nil), destination...), RowFlags: currentFlags,
	}
	if engine == nil {
		return result, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return result, errors.New("trailing-apostrophe normalizer has no context")
	}
	if lookup == nil {
		return result, errors.New("trailing-apostrophe normalizer has no context lookup")
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) == 0 || source[len(source)-1] != '\'' {
		return result, nil
	}
	result.Applied = true
	if len(source) >= 32 {
		return result, fmt.Errorf("trailing-apostrophe source has %d bytes, exceeding the native 32-byte local buffer", len(source))
	}
	if len(destination) == 0 {
		return result, errors.New("trailing-apostrophe destination has no byte for its NUL terminator")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}

	keyTables := text.Paul2013EmbeddedKeyTables()
	lookupDestination := destination
	if len(source) > 3 && text.Paul2013MappedCStringEqual(
		source[len(source)-3:], []byte("in'"), keyTables.CharacterMap,
	) {
		transformed := append([]byte(nil), source...)
		attributes := text.Paul2013UnsignedCharacterAttributeTable()
		replacement := byte('g')
		if attributes[transformed[len(transformed)-2]]&0x80 != 0 {
			replacement = 'G'
		}
		transformed[len(transformed)-1] = replacement
		writes, err := engine.NormalizePaul2013ContextTokenToBuffer(
			ctx, transformed, lookup, destination,
		)
		if err != nil {
			return Paul2013TrailingApostropheResult{}, fmt.Errorf("normalize mapped in-apostrophe form: %w", err)
		}
		if writes.Result.Matched {
			output := append([]byte(nil), writes.Output...)
			if end := bytes.IndexByte(output, 0); end >= 1 && output[end-1] == '.' {
				output[end-1] = '-'
			}
			result.Output = output
			result.RowFlags = currentFlags | 9
			result.Matched = true
			result.Path = "mapped-in-apostrophe-lookup"
			return result, nil
		}
		lookupDestination = writes.Output
	}

	base := append([]byte(nil), source[:len(source)-1]...)
	writes, err := engine.NormalizePaul2013ContextTokenToBuffer(ctx, base, lookup, lookupDestination)
	if err != nil {
		return Paul2013TrailingApostropheResult{}, fmt.Errorf("normalize apostrophe-stripped form: %w", err)
	}
	if writes.Result.Matched {
		result.Output = append([]byte(nil), writes.Output...)
		result.RowFlags = currentFlags | 9
		result.Matched = true
		result.Path = "apostrophe-stripped-lookup"
		return result, nil
	}

	generic, err := engine.NormalizePaul2013GenericToken(ctx, base)
	if err != nil {
		return Paul2013TrailingApostropheResult{}, fmt.Errorf("normalize apostrophe-stripped generic form: %w", err)
	}
	if generic.Eligible {
		genericWrites, err := ApplyPaul2013GenericNormalizerWrites(writes.Output, currentFlags, generic)
		if err != nil {
			return Paul2013TrailingApostropheResult{}, fmt.Errorf("write apostrophe-stripped generic form: %w", err)
		}
		result.Output = genericWrites.Output
		result.RowFlags = genericWrites.Flags | 8
		result.Matched = true
		result.Path = "generic-normalizer"
		return result, nil
	}

	codes, err := text.EncodePaul2013ContextString(base)
	if err != nil {
		return Paul2013TrailingApostropheResult{}, fmt.Errorf("encode apostrophe-stripped context form: %w", err)
	}
	if len(codes) >= len(writes.Output) {
		return Paul2013TrailingApostropheResult{}, fmt.Errorf("apostrophe-stripped context output needs %d bytes including NUL, destination has %d", len(codes)+1, len(writes.Output))
	}
	output := append([]byte(nil), writes.Output...)
	copy(output, codes)
	output[len(codes)] = 0
	result.Output = output
	if len(codes) != 0 {
		result.RowFlags = currentFlags | 0x28
		result.Matched = true
		result.Path = "context-encoder"
	} else {
		result.RowFlags = currentFlags
	}
	return result, nil
}

// NormalizePaul2013TrailingApostropheFromEmbeddedDictionary connects the
// supported FUN_1000c710 branch to the loaded Paul embedded lexicon.
func (engine *Engine) NormalizePaul2013TrailingApostropheFromEmbeddedDictionary(
	ctx context.Context,
	source []byte,
	destination []byte,
	currentFlags byte,
) (Paul2013TrailingApostropheResult, error) {
	return engine.NormalizePaul2013TrailingApostrophe(
		ctx, source, engine.LookupPaul2013EmbeddedContextToken, destination, currentFlags,
	)
}

// NormalizePaul2013ContextToken ports FUN_1000cb30's outer branch order while
// leaving its resource lookup callback explicit. A positive parsed-record
// count returns the lookup text directly. Otherwise the generic normalizer is
// entered only when the parsed short gate is nonzero; the existing inner
// implementation then applies its recovered eligibility and class/context
// dispatch. A zero count and zero gate is the native no-match result.
func (engine *Engine) NormalizePaul2013ContextToken(
	ctx context.Context,
	source []byte,
	lookup Paul2013ContextLookup,
) (Paul2013ContextNormalizationResult, error) {
	if engine == nil {
		return Paul2013ContextNormalizationResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013ContextNormalizationResult{}, errors.New("context normalization has no context")
	}
	if lookup == nil {
		return Paul2013ContextNormalizationResult{}, errors.New("context normalization has no dictionary lookup")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContextNormalizationResult{}, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	lookupResult, err := lookup(ctx, append([]byte(nil), source...))
	if err != nil {
		return Paul2013ContextNormalizationResult{}, fmt.Errorf("lookup normalized context for %q: %w", source, err)
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContextNormalizationResult{}, err
	}
	if lookupResult.RecordCount > 0 {
		text := lookupResult.DictionaryText
		if nul := bytes.IndexByte(text, 0); nul >= 0 {
			text = text[:nul]
		}
		return Paul2013ContextNormalizationResult{
			Matched: true, UsedDictionary: true, DictionaryText: append([]byte(nil), text...),
			NativeReturnCode: 1,
		}, nil
	}
	if !lookupResult.GenericFallbackGate {
		return Paul2013ContextNormalizationResult{NativeReturnCode: -1}, nil
	}
	generic, err := engine.NormalizePaul2013GenericToken(ctx, source)
	if err != nil {
		return Paul2013ContextNormalizationResult{}, fmt.Errorf("generic normalization for %q: %w", source, err)
	}
	if !generic.Eligible {
		return Paul2013ContextNormalizationResult{NativeReturnCode: -1}, nil
	}
	return Paul2013ContextNormalizationResult{
		Matched: true, Generic: generic, NativeReturnCode: 1,
	}, nil
}

// NormalizePaul2013ContextTokenToBuffer composes FUN_1000cb30's lookup and
// fallback order with its output-buffer behavior. Positive context-dictionary
// results are copied directly; generic results use
// ApplyPaul2013GenericNormalizerWrites; a no-match clears only the first
// destination byte, as the native entry point does before lookup. The lookup
// resource and caller-arena layout remain explicit.
func (engine *Engine) NormalizePaul2013ContextTokenToBuffer(
	ctx context.Context,
	source []byte,
	lookup Paul2013ContextLookup,
	destination []byte,
) (Paul2013ContextNormalizationWrites, error) {
	if len(destination) == 0 {
		return Paul2013ContextNormalizationWrites{}, errors.New("context normalization destination has no byte for its initial NUL")
	}
	result, err := engine.NormalizePaul2013ContextToken(ctx, source, lookup)
	if err != nil {
		return Paul2013ContextNormalizationWrites{}, err
	}
	if !result.Matched {
		output := append([]byte(nil), destination...)
		output[0] = 0
		return Paul2013ContextNormalizationWrites{Result: result, Output: output}, nil
	}
	if result.UsedDictionary {
		output, err := writePaul2013NormalizerCString(destination, result.DictionaryText)
		if err != nil {
			return Paul2013ContextNormalizationWrites{}, fmt.Errorf("write context-dictionary result: %w", err)
		}
		return Paul2013ContextNormalizationWrites{Result: result, Output: output}, nil
	}
	writes, err := ApplyPaul2013GenericNormalizerWrites(destination, 0, result.Generic)
	if err != nil {
		return Paul2013ContextNormalizationWrites{}, fmt.Errorf("write generic context result: %w", err)
	}
	return Paul2013ContextNormalizationWrites{
		Result: result, Output: writes.Output, LocalFlags: writes.Flags,
	}, nil
}

func writePaul2013NormalizerCString(destination, value []byte) ([]byte, error) {
	if bytes.IndexByte(value, 0) >= 0 {
		return nil, errors.New("context-dictionary result contains an embedded NUL")
	}
	if len(value) >= len(destination) {
		return nil, fmt.Errorf("context-dictionary result needs %d bytes including NUL, destination has %d", len(value)+1, len(destination))
	}
	output := append([]byte(nil), destination...)
	copy(output, value)
	output[len(value)] = 0
	return output, nil
}

// NormalizePaul2013GenericToken ports the supported ASCII branches of
// FUN_10002f10 after the native eligibility gate, including its vowel-gated
// MC-prefix expansion before character classification.
func (engine *Engine) NormalizePaul2013GenericToken(
	ctx context.Context,
	source []byte,
) (GenericNormalizerResult, error) {
	if engine == nil {
		return GenericNormalizerResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return GenericNormalizerResult{}, errors.New("generic normalizer has no context")
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	eligible, err := engine.IsPaul2013NormalizerEligible(source)
	if err != nil {
		return GenericNormalizerResult{}, err
	}
	result := GenericNormalizerResult{Eligible: eligible}
	if !eligible {
		return result, nil
	}
	if err := ctx.Err(); err != nil {
		return GenericNormalizerResult{}, err
	}

	normalized := normalizePaul2013GenericASCII(source)
	classifySurface := normalized
	if len(normalized) > 2 && normalized[0] == 'M' && normalized[1] == 'C' &&
		paul2013NormalizerHasVowel(normalized[2:]) {
		// FUN_1000fed0 uppercases the suffix and FUN_10063ed6 formats it
		// with the recovered literal "MAC%s" before the character scan.
		classifySurface = append([]byte("MAC"), normalized[2:]...)
	}
	if len(normalized) == 1 && paul2013NormalizerVowel(normalized[0]) {
		return paul2013GenericContextFallback(result, source)
	}

	classified, err := engine.ClassifyPaul2013NormalizerToken(ctx, classifySurface)
	if err != nil {
		return GenericNormalizerResult{}, fmt.Errorf("classify generic token %q: %w", source, err)
	}
	if !classified.Produced || !text.HasPaul2013NormalizerContextClass(classified.ClassCodes) {
		return paul2013GenericContextFallback(result, source)
	}
	result.ClassCodes = append([]byte(nil), classified.ClassCodes...)
	result.RowFlagMask = 0x10
	return result, nil
}

// NormalizePaul2013ModelContextRow composes FUN_10009030's model-row branch
// order with the loaded engine's FUN_10002f10 generic normalizer. The parser
// row association is read from the model row; parser-row eligibility and
// source bytes must already be present in the supplied arenas.
func (engine *Engine) NormalizePaul2013ModelContextRow(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	contextRowIndex int,
) (text.Paul2013ModelContextNormalizationResult, error) {
	if engine == nil {
		return text.Paul2013ModelContextNormalizationResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return text.Paul2013ModelContextNormalizationResult{}, errors.New("model context normalization has no context")
	}
	if err := ctx.Err(); err != nil {
		return text.Paul2013ModelContextNormalizationResult{}, err
	}
	return text.ApplyPaul2013ModelSourceClassNormalization(
		model, parserRows, contextRowIndex,
		func(surface []byte) (text.Paul2013GenericContextNormalization, error) {
			generic, err := engine.NormalizePaul2013GenericToken(ctx, surface)
			if err != nil {
				return text.Paul2013GenericContextNormalization{}, err
			}
			if !generic.Eligible {
				return text.Paul2013GenericContextNormalization{}, nil
			}
			writes, err := ApplyPaul2013GenericNormalizerWrites(
				make([]byte, 0x41), 0, generic,
			)
			if err != nil {
				return text.Paul2013GenericContextNormalization{}, err
			}
			end := bytes.IndexByte(writes.Output, 0)
			if end < 0 {
				return text.Paul2013GenericContextNormalization{}, errors.New("generic normalizer output is not NUL-terminated")
			}
			return text.Paul2013GenericContextNormalization{
				Handled: true, Codes: append([]byte(nil), writes.Output[:end]...), RowFlags: writes.Flags,
			}, nil
		},
	)
}

// NormalizePaul2013ModelContextRows applies FUN_10009030 in the explicit row
// order recovered by its caller. Row eligibility and order are not inferred
// here because FUN_10007520's preceding cascade determines both.
func (engine *Engine) NormalizePaul2013ModelContextRows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	contextRowIndexes []int,
) ([]byte, []int, error) {
	if engine == nil {
		return nil, nil, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return nil, nil, errors.New("model context normalization has no context")
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	working := append([]byte(nil), model...)
	handledRows := make([]int, 0, len(contextRowIndexes))
	for _, rowIndex := range contextRowIndexes {
		result, err := engine.NormalizePaul2013ModelContextRow(ctx, working, parserRows, rowIndex)
		if err != nil {
			return nil, nil, err
		}
		working = result.Model
		if result.Handled {
			handledRows = append(handledRows, rowIndex)
		}
	}
	return working, handledRows, nil
}

func paul2013GenericContextFallback(
	result GenericNormalizerResult,
	source []byte,
) (GenericNormalizerResult, error) {
	codes, err := text.EncodePaul2013ContextString(source)
	if err != nil {
		return GenericNormalizerResult{}, fmt.Errorf("encode generic-token context fallback: %w", err)
	}
	result.ContextCodes = codes
	result.UsedContextFallback = true
	if len(codes) != 0 {
		result.RowFlagMask |= 0x20
	}
	return result, nil
}

func normalizePaul2013GenericASCII(source []byte) []byte {
	result := append([]byte(nil), source...)
	for index, value := range result {
		if value >= 'a' && value <= 'z' {
			result[index] = value - ('a' - 'A')
		}
	}
	return result
}

func paul2013NormalizerHasVowel(source []byte) bool {
	for _, value := range source {
		if paul2013NormalizerVowel(value) {
			return true
		}
	}
	return false
}

func paul2013NormalizerVowel(value byte) bool {
	switch value {
	case 'A', 'E', 'I', 'O', 'U', 'Y', 'a', 'e', 'i', 'o', 'u', 'y':
		return true
	default:
		return false
	}
}
