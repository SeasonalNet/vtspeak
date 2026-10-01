package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013ErsSuffixResult records FUN_1000a140's short `ers` branch.
type Paul2013ErsSuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013ErsSuffixContext ports the 6- to 8-byte `ers` branch in
// FUN_1000a140. It preserves the ordered generic stems, two short-table
// predicates, conditional i-to-y rewrite, and branch-specific output tails.
func (engine *Engine) NormalizePaul2013ErsSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013ErsSuffixResult, error) {
	result := Paul2013ErsSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013ErsSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013ErsSuffixResult{}, errors.New("ers-suffix normalizer has no context")
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
	if !text.Paul2013MappedCStringEqual(source[len(source)-3:], []byte("ers"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true

	length := len(source)
	previous := characterMap[source[length-4]]
	var stem []byte
	if previous != 'e' {
		stem = append([]byte(nil), source[:length-2]...)
		output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
		if err != nil {
			return Paul2013ErsSuffixResult{}, fmt.Errorf("normalize ers-suffix rs-truncated stem %q: %w", stem, err)
		}
		if matched {
			return finishPaul2013ErsSuffix(result, stem, output, []byte{'D'}, currentFlags)
		}

		stem = append([]byte(nil), source[:length-3]...)
		output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
		if err != nil {
			return Paul2013ErsSuffixResult{}, fmt.Errorf("normalize ers-suffix ers-truncated stem %q: %w", stem, err)
		}
		if matched {
			return finishPaul2013ErsSuffix(result, stem, output, []byte{0x1a, 'D'}, currentFlags)
		}
	}

	classAt := func(index int) (byte, bool) {
		return paul2013TailClassForSourceByte(source[index], characterMap), true
	}
	classN6, ok := classAt(length - 6)
	if !ok {
		result.TableClassUnsupported = true
		return result, nil
	}
	classN5, ok := classAt(length - 5)
	if !ok {
		result.TableClassUnsupported = true
		return result, nil
	}

	if classN6 == 0 || classN5 != 0 || source[length-5] != source[length-4] {
		if classN5 != 0 || characterMap[source[length-4]] != 'i' {
			return result, nil
		}
		stem = append(append([]byte(nil), source[:length-4]...), 'y')
		output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
		if err != nil {
			return Paul2013ErsSuffixResult{}, fmt.Errorf("normalize ers-suffix i-to-y stem %q: %w", stem, err)
		}
		if !matched {
			return result, nil
		}
		return finishPaul2013ErsSuffix(result, stem, output, []byte{0x07, '7', '9'}, currentFlags)
	}

	stem = append([]byte(nil), source[:length-4]...)
	output, matched, err := engine.normalizePaul2013EmbeddedStem(ctx, stem)
	if err != nil {
		return Paul2013ErsSuffixResult{}, fmt.Errorf("normalize ers-suffix table-truncated stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	return finishPaul2013ErsSuffix(result, stem, output, []byte{0x1a, 'D'}, currentFlags)
}

func finishPaul2013ErsSuffix(
	result Paul2013ErsSuffixResult,
	stem []byte,
	stemOutput []byte,
	tail []byte,
	currentFlags byte,
) (Paul2013ErsSuffixResult, error) {
	if len(stemOutput)+len(tail) >= 68 {
		return Paul2013ErsSuffixResult{}, fmt.Errorf("ers-suffix output needs %d bytes plus NUL, native local capacity is 68", len(stemOutput)+len(tail))
	}
	result.Stem = append([]byte(nil), stem...)
	result.Output = append(append([]byte(nil), stemOutput...), tail...)
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
