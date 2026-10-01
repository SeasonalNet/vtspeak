package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

type paul2013SuffixRewrite struct {
	suffix      string
	replacement string
}

var paul2013MappedSuffixRewrites = [...]paul2013SuffixRewrite{
	{suffix: "isations", replacement: "izations"},
	{suffix: "isation", replacement: "ization"},
	{suffix: "isingly", replacement: "izingly"},
	{suffix: "isably", replacement: "izably"},
	{suffix: "isable", replacement: "izable"},
	{suffix: "ising", replacement: "izing"},
	{suffix: "ised", replacement: "ized"},
	{suffix: "ises", replacement: "izes"},
	{suffix: "ise", replacement: "ize"},
}

// Paul2013SuffixNormalizationResult reports a directly compared suffix, its
// rewritten spelling, and the result of the following FUN_1000cb30 lookup.
// Output contains the rewritten pronunciation/context code bytes only when
// the lookup produced a nonempty result.
type Paul2013SuffixNormalizationResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Transformed   []byte
	Output        []byte
	RowFlags      byte
}

// NormalizePaul2013MappedSpellingSuffix ports the explicit literal suffix
// cases in FUN_1000b4d0. It compares the original suffix through the recovered
// character map, replaces the matching suffix with the directly observed
// lowercase spelling, then invokes FUN_1000cb30 through the loaded embedded
// dictionary. On success the native helper ORs row flag bits 0 and 3. The
// native minimum lengths require at least one prefix byte before each suffix.
// In the native cascade this helper runs only after FUN_10009dc0 and
// FUN_1000a140 report no match.
func (engine *Engine) NormalizePaul2013MappedSpellingSuffix(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013SuffixNormalizationResult, error) {
	result := Paul2013SuffixNormalizationResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013SuffixNormalizationResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013SuffixNormalizationResult{}, errors.New("suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013SuffixNormalizationResult{}, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) == 0 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	var rewrite *paul2013SuffixRewrite
	for index := range paul2013MappedSuffixRewrites {
		candidate := &paul2013MappedSuffixRewrites[index]
		if len(source) > len(candidate.suffix) && text.Paul2013MappedCStringEqual(
			source[len(source)-len(candidate.suffix):], []byte(candidate.suffix), characterMap,
		) {
			rewrite = candidate
			break
		}
	}
	if rewrite == nil {
		return result, nil
	}
	result.SuffixMatched = true
	if len(source) >= 32 {
		return Paul2013SuffixNormalizationResult{}, fmt.Errorf("matched suffix source has %d bytes, exceeding FUN_1000b4d0's 32-byte local buffer", len(source))
	}
	cut := len(source) - len(rewrite.suffix)
	transformed := make([]byte, 0, cut+len(rewrite.replacement))
	transformed = append(transformed, source[:cut]...)
	transformed = append(transformed, rewrite.replacement...)
	result.Transformed = transformed

	writes, err := engine.NormalizePaul2013ContextTokenToBuffer(
		ctx, transformed, engine.LookupPaul2013EmbeddedContextToken, make([]byte, 68),
	)
	if err != nil {
		return Paul2013SuffixNormalizationResult{}, fmt.Errorf("normalize rewritten suffix %q: %w", transformed, err)
	}
	if !writes.Result.Matched || len(writes.Output) == 0 || writes.Output[0] == 0 {
		return result, nil
	}
	end := bytes.IndexByte(writes.Output, 0)
	if end < 0 {
		return Paul2013SuffixNormalizationResult{}, errors.New("rewritten suffix output is not NUL-terminated")
	}
	result.Matched = true
	result.Output = append([]byte(nil), writes.Output[:end]...)
	result.RowFlags = currentFlags | 9
	return result, nil
}
