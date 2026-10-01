package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
)

// Paul2013ApostropheSBaseResult is the prefix result consumed by the direct
// apostrophe-s branch in FUN_1000c3a0.
type Paul2013ApostropheSBaseResult struct {
	Matched bool
	Output  []byte
}

// Paul2013ApostropheSBaseNormalizer supplies FUN_1000c3a0's preceding
// FUN_1000cb30 stem lookup and normalization.
type Paul2013ApostropheSBaseNormalizer func(
	context.Context,
	[]byte,
) (Paul2013ApostropheSBaseResult, error)

// Paul2013ApostropheSSuffixResult reports the direct 's branch from
// FUN_1000c3a0. Output excludes the C NUL terminator.
type Paul2013ApostropheSSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
}

// NormalizePaul2013ApostropheSSuffixContext ports FUN_1000c3a0's exact 's
// suffix rewrite after its preceding handlers have missed. Stem normalization
// is a callback representing FUN_1000cb30. The native suffix comparison is
// case-sensitive.
func NormalizePaul2013ApostropheSSuffixContext(
	ctx context.Context,
	source []byte,
	normalize Paul2013ApostropheSBaseNormalizer,
) (Paul2013ApostropheSSuffixResult, error) {
	result := Paul2013ApostropheSSuffixResult{Source: append([]byte(nil), source...)}
	if ctx == nil {
		return Paul2013ApostropheSSuffixResult{}, errors.New("apostrophe-s suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) < 2 || !bytes.HasSuffix(source, []byte("'s")) {
		return result, nil
	}
	result.SuffixMatched = true
	result.Stem = append([]byte(nil), source[:len(source)-2]...)
	if normalize == nil {
		return Paul2013ApostropheSSuffixResult{}, errors.New("apostrophe-s suffix has no stem handler")
	}
	base, err := normalize(ctx, append([]byte(nil), result.Stem...))
	if err != nil {
		return Paul2013ApostropheSSuffixResult{}, fmt.Errorf("normalize apostrophe-s stem %q: %w", result.Stem, err)
	}
	if !base.Matched || len(base.Output) == 0 {
		return result, nil
	}
	if bytes.IndexByte(base.Output, 0) >= 0 {
		return Paul2013ApostropheSSuffixResult{}, errors.New("apostrophe-s stem output contains NUL")
	}
	output := append([]byte(nil), base.Output...)
	last := output[len(output)-1]
	switch {
	case last == ' ' || last == '*' || last == '5' || last == '9' || last == ':':
		output = append(output, '7')
	case last == 0x14 || last == ')' || last == '7' || last == '8' || last == 'D' || last == 'E':
		output = append(output, '#', 'D')
	default:
		output = append(output, 'D')
	}
	result.Output = output
	result.Matched = true
	return result, nil
}

// NormalizePaul2013ApostropheSSuffixFromEmbeddedDictionary connects the
// apostrophe-s rewrite to the loaded dictionary and generic normalizer for its
// truncated stem. Call it only after the earlier FUN_1000c3a0 handlers miss.
func (engine *Engine) NormalizePaul2013ApostropheSSuffixFromEmbeddedDictionary(
	ctx context.Context,
	source []byte,
) (Paul2013ApostropheSSuffixResult, error) {
	if engine == nil {
		return Paul2013ApostropheSSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	return NormalizePaul2013ApostropheSSuffixContext(ctx, source,
		func(ctx context.Context, stem []byte) (Paul2013ApostropheSBaseResult, error) {
			normalized, err := engine.NormalizePaul2013ContextTokenFromEmbeddedDictionary(ctx, stem)
			if err != nil {
				return Paul2013ApostropheSBaseResult{}, err
			}
			base := Paul2013ApostropheSBaseResult{Matched: normalized.Matched}
			switch {
			case normalized.UsedDictionary:
				base.Output = append([]byte(nil), normalized.DictionaryText...)
			case normalized.Generic.UsedContextFallback:
				base.Output = append([]byte(nil), normalized.Generic.ContextCodes...)
			default:
				base.Output = append([]byte(nil), normalized.Generic.ClassCodes...)
			}
			return base, nil
		},
	)
}
