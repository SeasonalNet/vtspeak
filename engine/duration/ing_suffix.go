package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013IngSuffixResult records FUN_1000a140's `ing` branch.
type Paul2013IngSuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013LongIngSuffixContext preserves the original entry point.
func (engine *Engine) NormalizePaul2013LongIngSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013IngSuffixResult, error) {
	return engine.NormalizePaul2013IngSuffixContext(ctx, source, currentFlags)
}

// NormalizePaul2013IngSuffixContext ports the 5- to 31-byte FUN_1000a140
// `ing` path for ASCII table indexes. It applies the marker-selected `%`
// lookup before the generic lookup, the recovered table-gated stem edits, and
// the native `#.` output tail.
func (engine *Engine) NormalizePaul2013IngSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013IngSuffixResult, error) {
	result := Paul2013IngSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013IngSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013IngSuffixResult{}, errors.New("ing-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) < 5 || len(source) > 31 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(source[len(source)-3:], []byte("ing"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true

	classAt := func(index int) (byte, bool) {
		return paul2013TailClassForSourceByte(source[index], characterMap), true
	}
	length := len(source)
	classN5, ok := classAt(length - 5)
	if !ok {
		result.TableClassUnsupported = true
		return result, nil
	}
	classN4, ok := classAt(length - 4)
	if !ok {
		result.TableClassUnsupported = true
		return result, nil
	}

	baseStem := append([]byte(nil), source[:length-3]...)
	result.Stem = append([]byte(nil), baseStem...)
	var output []byte
	var matched bool
	var err error
	if classN5 == 0 || classN4 != 0 {
		tryShorterStem := false
		// The native short-path guard skips the length-7 and length-6
		// table reads when the input has fewer than seven bytes.
		if length >= 7 {
			classN7, ok := classAt(length - 7)
			if !ok {
				result.TableClassUnsupported = true
				return result, nil
			}
			if classN7 == 0 {
				classN6, ok := classAt(length - 6)
				if !ok {
					result.TableClassUnsupported = true
					return result, nil
				}
				tryShorterStem = classN6 != 0 && classN5 == 0 && source[length-5] == source[length-4]
			}
		}
		output, matched, err = engine.normalizePaul2013PercentStem(ctx, baseStem)
		if err != nil {
			return Paul2013IngSuffixResult{}, fmt.Errorf("normalize ing base stem %q: %w", baseStem, err)
		}
		if !matched && tryShorterStem {
			stem := append([]byte(nil), source[:length-4]...)
			result.Stem = append([]byte(nil), stem...)
			output, matched, err = engine.normalizePaul2013PercentStem(ctx, stem)
			if err != nil {
				return Paul2013IngSuffixResult{}, fmt.Errorf("normalize ing shortened stem %q: %w", stem, err)
			}
		}
	} else {
		// When the short-table pair is 1/0, FUN_1000a140 first substitutes
		// `e` for the `i` beginning `ing`, then retries the direct base stem.
		stem := append(append([]byte(nil), baseStem...), 'e')
		result.Stem = append([]byte(nil), stem...)
		output, matched, err = engine.normalizePaul2013PercentStem(ctx, stem)
		if err != nil {
			return Paul2013IngSuffixResult{}, fmt.Errorf("normalize ing e-stem %q: %w", stem, err)
		}
		if !matched {
			result.Stem = append([]byte(nil), baseStem...)
			output, matched, err = engine.normalizePaul2013PercentStem(ctx, baseStem)
			if err != nil {
				return Paul2013IngSuffixResult{}, fmt.Errorf("normalize ing base retry %q: %w", baseStem, err)
			}
		}
	}
	if !matched || len(output) == 0 {
		return result, nil
	}
	if len(output)+2 >= 68 {
		return Paul2013IngSuffixResult{}, fmt.Errorf("ing-suffix output needs %d bytes plus NUL, native local capacity is 68", len(output)+2)
	}
	result.Output = append(append([]byte(nil), output...), '#', '.')
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
