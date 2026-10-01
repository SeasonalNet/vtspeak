package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013EstSuffixResult reports the direct short `est` path in
// FUN_1000a140. Output excludes its trailing C NUL.
type Paul2013EstSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// NormalizePaul2013EstSuffixContext ports the 6- to 8-byte `est` branch in
// FUN_1000a140. The native path removes the final four source bytes, performs
// generic stem normalization, replaces its terminator with 0x07, then writes
// `-7` and sets row-flag bit 3.
func (engine *Engine) NormalizePaul2013EstSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013EstSuffixResult, error) {
	result := Paul2013EstSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013EstSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013EstSuffixResult{}, errors.New("est-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 6 || len(source) > 8 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(source[len(source)-3:], []byte("est"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013EstSuffixResult{}, fmt.Errorf("normalize est-suffix stem %q: %w", stem, err)
	}
	if !matched || len(output) == 0 {
		return result, nil
	}
	tail := []byte{0x07, '-', '7'}
	if len(output)+len(tail) >= 68 {
		return Paul2013EstSuffixResult{}, fmt.Errorf("est-suffix output needs %d bytes plus NUL, native local capacity is 68", len(output)+len(tail))
	}
	result.Output = append(append([]byte(nil), output...), tail...)
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
