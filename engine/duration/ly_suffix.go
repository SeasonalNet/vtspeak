package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013LySuffixResult reports FUN_1000a140's mapped "ly" branch and its
// bounded short-table lookup.
type Paul2013LySuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013LongLySuffixContext preserves the original entry point and
// supports the recovered 5- to 31-byte FUN_1000a140 `ly` branch. It
// normalizes the source without the suffix, then applies the observed
// table-gated final-i-to-y retry on a miss. Success applies the direct
// ampersand/control/plus tail and row-flag bit 3.
func (engine *Engine) NormalizePaul2013LongLySuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013LySuffixResult, error) {
	return engine.NormalizePaul2013LySuffixContext(ctx, source, currentFlags)
}

// NormalizePaul2013LySuffixContext ports the 5- to 31-byte FUN_1000a140
// branch for C strings ending in mapped "ly". It normalizes the source
// without the suffix, then applies the observed
// table-gated final-i-to-y retry on a miss. Success applies the direct
// ampersand/control/plus tail and row-flag bit 3.
func (engine *Engine) NormalizePaul2013LySuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013LySuffixResult, error) {
	result := Paul2013LySuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013LySuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013LySuffixResult{}, errors.New("ly-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 5 || len(source) > 31 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(source[len(source)-2:], []byte("ly"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true

	stem := append([]byte(nil), source[:len(source)-2]...)
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013LySuffixResult{}, fmt.Errorf("normalize ly stem %q: %w", stem, err)
	}
	if matched {
		return finishPaul2013LySuffix(result, output, currentFlags)
	}

	if paul2013TailClassForSourceByte(source[len(source)-4], characterMap) != 0 ||
		characterMap[source[len(source)-3]] != 'i' {
		return result, nil
	}
	stem = append(append([]byte(nil), source[:len(source)-3]...), 'y')
	result.Stem = append([]byte(nil), stem...)
	output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013LySuffixResult{}, fmt.Errorf("normalize ly y-retry stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	return finishPaul2013LySuffix(result, output, currentFlags)
}

func finishPaul2013LySuffix(
	result Paul2013LySuffixResult,
	stemOutput []byte,
	currentFlags byte,
) (Paul2013LySuffixResult, error) {
	if len(stemOutput) == 0 {
		return result, nil
	}
	if stemOutput[len(stemOutput)-1] == '&' {
		stemOutput = append([]byte(nil), stemOutput...)
		stemOutput[len(stemOutput)-1] = 0x07
	}
	marker := []byte{'&'}
	if stemOutput[len(stemOutput)-1] != '+' {
		marker = append([]byte{'+'}, marker...)
	}
	if len(stemOutput)+len(marker) >= 68 {
		return Paul2013LySuffixResult{}, fmt.Errorf("ly-suffix output needs %d bytes plus NUL, native local capacity is 68", len(stemOutput)+len(marker))
	}
	result.Output = append(append([]byte(nil), stemOutput...), marker...)
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
