package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013ErSuffixResult reports the bounded long-input "er" branch of
// FUN_1000a140. TableClassUnsupported remains for API compatibility; the
// recovered signed table window covers all native int8 indexes.
type Paul2013ErSuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013ErSuffixContext ports FUN_1000a140's long-input "er"
// branch for C-string lengths 9 through 31. The preceding-letter truncation
// lookups, raw table predicates, output marker, and row flag are preserved.
// The outer suffix cascade and non-letter table entries remain separate.
func (engine *Engine) NormalizePaul2013ErSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013ErSuffixResult, error) {
	result := Paul2013ErSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013ErSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013ErSuffixResult{}, errors.New("er-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ErSuffixResult{}, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) < 9 || len(source) > 31 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(source[len(source)-2:], []byte("er"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true

	if characterMap[source[len(source)-3]] != 'e' {
		for _, trim := range []int{1, 2} {
			stem := append([]byte(nil), source[:len(source)-trim]...)
			result.Stem = stem
			output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
			if err != nil {
				return Paul2013ErSuffixResult{}, fmt.Errorf("normalize er-suffix base form %q: %w", stem, err)
			}
			if matched {
				return finishPaul2013ErSuffix(result, output, currentFlags)
			}
		}
	}

	stem, supported, hasStem, err := paul2013LongSuffixTableStem(source, characterMap)
	if err != nil {
		return Paul2013ErSuffixResult{}, err
	}
	if !supported {
		result.TableClassUnsupported = true
		return result, nil
	}
	if !hasStem {
		return result, nil
	}
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013ErSuffixResult{}, fmt.Errorf("normalize er-suffix table stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	return finishPaul2013ErSuffix(result, output, currentFlags)
}

func finishPaul2013ErSuffix(
	result Paul2013ErSuffixResult,
	stemOutput []byte,
	currentFlags byte,
) (Paul2013ErSuffixResult, error) {
	if len(stemOutput) == 0 {
		return result, nil
	}
	if len(stemOutput)+1 >= 68 {
		return Paul2013ErSuffixResult{}, fmt.Errorf("er-suffix output needs %d bytes plus NUL, native local capacity is 68", len(stemOutput)+1)
	}
	result.Output = append(append([]byte(nil), stemOutput...), 0x1a)
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
