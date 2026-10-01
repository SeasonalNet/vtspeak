package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013IstSuffixResult records FUN_1000a140's mapped `ist` branch.
type Paul2013IstSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// Paul2013LongIstSuffixResult preserves the original result type name.
type Paul2013LongIstSuffixResult = Paul2013IstSuffixResult

// NormalizePaul2013LongIstSuffixContext preserves the original entry point.
func (engine *Engine) NormalizePaul2013LongIstSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013LongIstSuffixResult, error) {
	return engine.NormalizePaul2013IstSuffixContext(ctx, source, currentFlags)
}

// NormalizePaul2013IstSuffixContext ports the 5- to 31-byte mapped `ist`
// path in FUN_1000a140. It removes the suffix, runs FUN_1000cb30 through the
// loaded embedded dictionary, and appends `#79` to a nonempty result.
func (engine *Engine) NormalizePaul2013IstSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013IstSuffixResult, error) {
	result := Paul2013LongIstSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013IstSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013IstSuffixResult{}, errors.New("ist-suffix normalizer has no context")
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
	if !text.Paul2013MappedCStringEqual(source[len(source)-3:], []byte("ist"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-3]...)
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013IstSuffixResult{}, fmt.Errorf("normalize ist stem %q: %w", stem, err)
	}
	if !matched || len(output) == 0 {
		return result, nil
	}
	if len(output)+3 >= 68 {
		return Paul2013IstSuffixResult{}, fmt.Errorf("ist-suffix output needs %d bytes plus NUL, native local capacity is 68", len(output)+3)
	}
	result.Output = append(append([]byte(nil), output...), '#', '7', '9')
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
