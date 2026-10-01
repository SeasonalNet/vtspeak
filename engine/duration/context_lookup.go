package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// LookupPaul2013EmbeddedContextToken supplies the embedded-dictionary branch
// consumed by FUN_1000cb30. It follows FUN_10003a70's embedded-key transform
// and FUN_10003c50's payload parse: RecordCount is the pronunciation count,
// DictionaryText is the first expanded internal-symbol string, and the
// generic-normalizer gate is payload metadata bit 7. The caller still supplies
// an already-normalized single-token surface.
func (engine *Engine) LookupPaul2013EmbeddedContextToken(
	ctx context.Context,
	surface []byte,
) (Paul2013ContextLookupResult, error) {
	result, _, err := engine.lookupPaul2013EmbeddedContextTokenWithType(ctx, surface)
	return result, err
}

// BuildPaul2013C3A0AlternateRecord connects the embedded lookup and phone-ID
// expansion to the caller-local FUN_10003c50 record consumed by FUN_100086c0.
// A missing token produces the zeroed native record, matching the parser's
// empty-payload branch.
func (engine *Engine) BuildPaul2013C3A0AlternateRecord(
	ctx context.Context,
	surface []byte,
) ([]byte, error) {
	if engine == nil || engine.dictionary == nil {
		return nil, errors.New("Paul 2013 C3A0 dictionary is nil or unloaded")
	}
	if ctx == nil {
		return nil, errors.New("Paul 2013 C3A0 record lookup has no context")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if nul := bytes.IndexByte(surface, 0); nul >= 0 {
		surface = surface[:nul]
	}
	resolved, found, err := engine.dictionary.ResolvePaul2013Surface(surface)
	if err != nil {
		return nil, fmt.Errorf("resolve C3A0 alternate record for %q: %w", surface, err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !found {
		return (text.PhonePayload{}).BuildPaul2013ParsedDictionaryRecord(text.Paul2013PhoneIDCodebook())
	}
	return resolved.Payload.BuildPaul2013ParsedDictionaryRecord(text.Paul2013PhoneIDCodebook())
}

func (engine *Engine) lookupPaul2013EmbeddedContextTokenWithType(
	ctx context.Context,
	surface []byte,
) (Paul2013ContextLookupResult, byte, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013ContextLookupResult{}, 0, errors.New("Paul 2013 embedded context dictionary is nil or unloaded")
	}
	if ctx == nil {
		return Paul2013ContextLookupResult{}, 0, errors.New("embedded context lookup has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContextLookupResult{}, 0, err
	}
	if nul := bytes.IndexByte(surface, 0); nul >= 0 {
		surface = surface[:nul]
	}
	resolved, found, err := engine.dictionary.ResolvePaul2013Surface(surface)
	if err != nil {
		return Paul2013ContextLookupResult{}, 0, fmt.Errorf("resolve embedded context token %q: %w", surface, err)
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContextLookupResult{}, 0, err
	}
	if !found {
		return Paul2013ContextLookupResult{}, 0, nil
	}
	result := Paul2013ContextLookupResult{
		RecordCount:         int32(len(resolved.Alternatives)),
		GenericFallbackGate: resolved.Payload.Metadata[3],
	}
	if len(resolved.Alternatives) > 0 {
		result.DictionaryText = append([]byte(nil), resolved.Alternatives[0].Symbols...)
	}
	return result, resolved.Payload.ResultType, nil
}

// TransformPaul2013ModelSourceFromEmbeddedDictionary connects
// FUN_1000d640's terminal-punctuation decision to the loaded embedded
// dictionary. The native branch tests FUN_10003c50's result type byte, so a
// record counts as a positive decision only when that decoded byte is nonzero.
func (engine *Engine) TransformPaul2013ModelSourceFromEmbeddedDictionary(
	ctx context.Context,
	input text.Paul2013ModelSourceTransformInput,
) (text.Paul2013ModelSourceTransformResult, error) {
	lookup := func(surface []byte) (bool, error) {
		_, resultType, err := engine.lookupPaul2013EmbeddedContextTokenWithType(ctx, surface)
		if err != nil {
			return false, err
		}
		return resultType != 0, nil
	}
	transformed, err := text.TransformPaul2013ModelSourceStringFromPaul2013CharacterTable(input, lookup)
	if err != nil {
		return text.Paul2013ModelSourceTransformResult{}, err
	}
	return transformed, nil
}

// NormalizePaul2013ContextTokenFromEmbeddedDictionary connects the loaded
// embedded lexicon to FUN_1000cb30's recovered outer branch order. The TPP
// resource family and any caller-specific token transformation remain
// separate inputs to other native paths.
func (engine *Engine) NormalizePaul2013ContextTokenFromEmbeddedDictionary(
	ctx context.Context,
	surface []byte,
) (Paul2013ContextNormalizationResult, error) {
	return engine.NormalizePaul2013ContextToken(ctx, surface, engine.LookupPaul2013EmbeddedContextToken)
}

// NormalizePaul2013ContextTokenFromEmbeddedDictionaryToBuffer connects the
// loaded embedded lexicon, supported generic fallback, and native output
// buffer writes for one already-normalized token.
func (engine *Engine) NormalizePaul2013ContextTokenFromEmbeddedDictionaryToBuffer(
	ctx context.Context,
	surface []byte,
	destination []byte,
) (Paul2013ContextNormalizationWrites, error) {
	return engine.NormalizePaul2013ContextTokenToBuffer(
		ctx, surface, engine.LookupPaul2013EmbeddedContextToken, destination,
	)
}
