package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013FulSuffixResult records FUN_1000a140's mapped `ful` branch.
type Paul2013FulSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// NormalizePaul2013FulSuffixContext ports the 5- to 31-byte mapped `ful`
// branch. It normalizes the suffix-stripped stem and appends space, 0x07, and
// `+` to a nonempty result.
func (engine *Engine) NormalizePaul2013FulSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013FulSuffixResult, error) {
	result := Paul2013FulSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013FulSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013FulSuffixResult{}, errors.New("ful-suffix normalizer has no context")
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
	if !text.Paul2013MappedCStringEqual(source[len(source)-3:], []byte("ful"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-3]...)
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013FulSuffixResult{}, fmt.Errorf("normalize ful stem %q: %w", stem, err)
	}
	if !matched || len(output) == 0 {
		return result, nil
	}
	if len(output)+3 >= 68 {
		return Paul2013FulSuffixResult{}, fmt.Errorf("ful-suffix output needs %d bytes plus NUL, native local capacity is 68", len(output)+3)
	}
	result.Output = append(append([]byte(nil), output...), ' ', 0x07, '+')
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
