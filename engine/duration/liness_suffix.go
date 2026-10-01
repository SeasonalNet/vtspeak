package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013LinessSuffixResult reports the directly gated long-input
// FUN_1000a140 path. TableClassUnsupported remains for API compatibility;
// the recovered signed table window covers all native int8 indexes.
type Paul2013LinessSuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013LinessSuffixContext ports FUN_1000a140's 9- through
// 31-byte path gated by the mapped literal suffix "liness". It first
// normalizes the source with "ness" removed; after a miss, it applies the
// recovered table-gated final-i-to-y retry. Success appends the observed
// -#7 marker and sets row-flag bit 3.
func (engine *Engine) NormalizePaul2013LinessSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013LinessSuffixResult, error) {
	result := Paul2013LinessSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013LinessSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013LinessSuffixResult{}, errors.New("liness-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 9 || len(source) > 31 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(
		source[len(source)-6:], []byte("liness"), characterMap,
	) {
		return result, nil
	}
	result.SuffixMatched = true

	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013LinessSuffixResult{}, fmt.Errorf("normalize liness stem %q: %w", stem, err)
	}
	if matched {
		return finishPaul2013LinessSuffix(result, output, currentFlags)
	}

	if paul2013TailClassForSourceByte(source[len(source)-6], characterMap) != 0 ||
		characterMap[source[len(source)-5]] != 'i' {
		return result, nil
	}
	stem = append(append([]byte(nil), source[:len(source)-5]...), 'y')
	result.Stem = append([]byte(nil), stem...)
	output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013LinessSuffixResult{}, fmt.Errorf("normalize liness y-retry stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	return finishPaul2013LinessSuffix(result, output, currentFlags)
}

func finishPaul2013LinessSuffix(
	result Paul2013LinessSuffixResult,
	stemOutput []byte,
	currentFlags byte,
) (Paul2013LinessSuffixResult, error) {
	if len(stemOutput) == 0 {
		return result, nil
	}
	if len(stemOutput)+3 >= 68 {
		return Paul2013LinessSuffixResult{}, fmt.Errorf("liness-suffix output needs %d bytes plus NUL, native local capacity is 68", len(stemOutput)+3)
	}
	result.Output = append(append([]byte(nil), stemOutput...), '-', '#', '7')
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
