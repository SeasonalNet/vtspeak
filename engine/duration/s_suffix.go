package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013SSuffixResult reports the bounded terminal-s branch from
// FUN_1000a140. TableClassUnsupported remains for API compatibility; the
// recovered signed table window covers all native int8 indexes.
type Paul2013SSuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013SSuffixContext ports the 5- to 31-byte terminal-s branch
// in FUN_1000a140. It preserves the native stem lookup order, mapped
// fixed-width comparisons, fallback rejection, marker tail, and row flag.
func (engine *Engine) NormalizePaul2013SSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013SSuffixResult, error) {
	result := Paul2013SSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013SSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013SSuffixResult{}, errors.New("s-suffix normalizer has no context")
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
	if len(source) >= 9 && text.Paul2013MappedCStringEqual(
		source[len(source)-6:], []byte("liness"), characterMap,
	) {
		return result, nil
	}
	if characterMap[source[len(source)-1]] != 's' {
		return result, nil
	}
	result.SuffixMatched = true

	classAt := func(index int) (byte, bool) {
		return paul2013TailClassForSourceByte(source[index], characterMap), true
	}
	mappedPairEqual := func(start int, expected string) bool {
		if len(expected) != 2 || start < 0 || start+2 > len(source) {
			return false
		}
		return characterMap[source[start]] == expected[0] && characterMap[source[start+1]] == expected[1]
	}
	isOneOf := func(value byte, options string) bool {
		for index := 0; index < len(options); index++ {
			if value == options[index] {
				return true
			}
		}
		return false
	}
	length := len(source)
	previous := characterMap[source[length-2]]
	var stem []byte
	var output []byte
	var matched bool
	rejectProtectedOutput := false
	var err error

	// The first branch strips only s, then retries without es on a miss.
	firstPluralCase := false
	if previous == 'e' {
		beforeE := characterMap[source[length-3]]
		firstPluralCase = isOneOf(beforeE, "szx") ||
			mappedPairEqual(length-4, "ch") || mappedPairEqual(length-4, "sh")
		if !firstPluralCase {
			class, ok := classAt(length - 4)
			if !ok {
				result.TableClassUnsupported = true
				return result, nil
			}
			firstPluralCase = class == 0 && beforeE == 'o'
		}
	}
	if firstPluralCase {
		stem = append([]byte(nil), source[:length-1]...)
		output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
		if err != nil {
			return Paul2013SSuffixResult{}, fmt.Errorf("normalize s-suffix first plural stem %q: %w", stem, err)
		}
		if !matched {
			stem = append([]byte(nil), source[:length-2]...)
			output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
			if err != nil {
				return Paul2013SSuffixResult{}, fmt.Errorf("normalize s-suffix plural fallback stem %q: %w", stem, err)
			}
		}
	} else {
		class, ok := classAt(length - 4)
		if !ok {
			result.TableClassUnsupported = true
			return result, nil
		}
		if class == 0 && mappedPairEqual(length-3, "ie") {
			stem = append(append([]byte(nil), source[:length-3]...), 'y')
			output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
			if err != nil {
				return Paul2013SSuffixResult{}, fmt.Errorf("normalize s-suffix ie-to-y stem %q: %w", stem, err)
			}
		} else {
			endsChSh := mappedPairEqual(length-3, "ch") || mappedPairEqual(length-3, "sh")
			class, ok = classAt(length - 3)
			if !ok {
				result.TableClassUnsupported = true
				return result, nil
			}
			ordinarySReduction := !isOneOf(previous, "szx") && !endsChSh &&
				(class != 0 || previous != 'o')
			if ordinarySReduction {
				stem = append([]byte(nil), source[:length-1]...)
				output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
				if err != nil {
					return Paul2013SSuffixResult{}, fmt.Errorf("normalize s-suffix base stem %q: %w", stem, err)
				}
				if !matched && previous == 'e' {
					stem = append([]byte(nil), source[:length-2]...)
					output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
					if err != nil {
						return Paul2013SSuffixResult{}, fmt.Errorf("normalize s-suffix es fallback stem %q: %w", stem, err)
					}
				}
			} else {
				stem = append([]byte(nil), source[:length-2]...)
				output, matched, err = engine.normalizePaul2013EmbeddedStem(ctx, stem)
				if err != nil {
					return Paul2013SSuffixResult{}, fmt.Errorf("normalize s-suffix plural stem %q: %w", stem, err)
				}
				rejectProtectedOutput = class != 0 || previous != 'o'
			}
		}
	}
	result.Stem = append([]byte(nil), stem...)
	if !matched || len(output) == 0 {
		return result, nil
	}

	// The final general remove-es branch rejects an empty or protected-ending result.
	if rejectProtectedOutput {
		last := output[len(output)-1]
		if isOneOf(last, string([]byte{0x14, ')', '7', '8', 'D', 'E'})) {
			return result, nil
		}
	}

	last := output[len(output)-1]
	var tail []byte
	switch {
	case isOneOf(last, " *59:"):
		tail = []byte{'7'}
	case isOneOf(last, string([]byte{0x14, ')', '7', '8', 'D', 'E'})):
		tail = []byte{'#', 'D'}
	default:
		tail = []byte{'D'}
	}
	if len(output)+len(tail) >= 68 {
		return Paul2013SSuffixResult{}, fmt.Errorf("s-suffix output needs %d bytes plus NUL, native local capacity is 68", len(output)+len(tail))
	}
	result.Output = append(append([]byte(nil), output...), tail...)
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
