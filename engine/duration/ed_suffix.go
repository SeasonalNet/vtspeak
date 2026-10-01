package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013EdSuffixResult reports FUN_1000a140's bounded long-input "ed"
// branch. TableClassUnsupported remains for API compatibility; the recovered
// signed table window covers all native int8 indexes.
type Paul2013EdSuffixResult struct {
	SuffixMatched         bool
	Matched               bool
	TableClassUnsupported bool
	Source                []byte
	Stem                  []byte
	Output                []byte
	RowFlags              byte
}

// NormalizePaul2013EdSuffixContext ports FUN_1000a140's direct "ed" branch
// for C-string lengths 9 through 31. It preserves the marker-aware '%' lookup,
// FUN_1000cb30 fallback, and the table-gated stem edits and output marker.
// The outer function's other suffix branches and non-ASCII table accesses
// remain separate.
func (engine *Engine) NormalizePaul2013EdSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013EdSuffixResult, error) {
	result := Paul2013EdSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013EdSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013EdSuffixResult{}, errors.New("ed-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013EdSuffixResult{}, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) < 9 || len(source) > 31 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(source[len(source)-2:], []byte("ed"), characterMap) {
		return result, nil
	}
	result.SuffixMatched = true

	// For a non-e byte immediately before "ed", the native first truncates
	// both suffix bytes and tries marker lookup. An e byte proceeds directly
	// to the table-gated stem edits.
	if characterMap[source[len(source)-3]] != 'e' {
		stem := append([]byte(nil), source[:len(source)-2]...)
		result.Stem = stem
		output, matched, err := engine.normalizePaul2013PercentStem(ctx, stem)
		if err != nil {
			return Paul2013EdSuffixResult{}, fmt.Errorf("normalize ed-suffix base form %q: %w", stem, err)
		}
		if matched {
			return finishPaul2013EdSuffix(result, output, currentFlags)
		}
	}

	stem, supported, hasStem, err := paul2013LongSuffixTableStem(source, characterMap)
	if err != nil {
		return Paul2013EdSuffixResult{}, err
	}
	if !supported {
		result.TableClassUnsupported = true
		return result, nil
	}
	if !hasStem {
		return result, nil
	}
	result.Stem = append([]byte(nil), stem...)
	output, matched, err := engine.normalizePaul2013PercentStem(ctx, stem)
	if err != nil {
		return Paul2013EdSuffixResult{}, fmt.Errorf("normalize ed-suffix table stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	return finishPaul2013EdSuffix(result, output, currentFlags)
}

func paul2013LongSuffixTableStem(source []byte, characterMap [256]byte) ([]byte, bool, bool, error) {
	if len(source) < 7 {
		return nil, true, false, nil
	}
	classAt := func(sourceIndex int) (byte, bool) {
		return paul2013TailClassForSourceByte(source[sourceIndex], characterMap), true
	}
	classN5, ok := classAt(len(source) - 5)
	if !ok {
		return nil, false, false, nil
	}
	classN4, ok := classAt(len(source) - 4)
	if !ok {
		return nil, false, false, nil
	}
	if classN5 == 0 || classN4 != 0 || source[len(source)-4] != source[len(source)-3] {
		mappedN3 := characterMap[source[len(source)-3]]
		if classN4 != 0 || mappedN3 != 'i' {
			return nil, true, false, nil
		}
		stem := append([]byte(nil), source[:len(source)-3]...)
		stem = append(stem, 'y')
		return stem, true, true, nil
	}
	return append([]byte(nil), source[:len(source)-3]...), true, true, nil
}

func (engine *Engine) normalizePaul2013PercentStem(
	ctx context.Context,
	stem []byte,
) ([]byte, bool, error) {
	selected, err := engine.SelectPaul2013ContextPronunciationFromEmbeddedDictionary(
		ctx, stem, '%',
	)
	if err != nil {
		return nil, false, err
	}
	if selected.ReturnCode == 1 {
		end := bytes.IndexByte(selected.Output, 0)
		if end < 0 {
			return nil, false, errors.New("path-selected stem output is not NUL-terminated")
		}
		return append([]byte(nil), selected.Output[:end]...), true, nil
	}
	return engine.normalizePaul2013EmbeddedStem(ctx, stem)
}

func (engine *Engine) normalizePaul2013EmbeddedStem(
	ctx context.Context,
	stem []byte,
) ([]byte, bool, error) {
	writes, err := engine.NormalizePaul2013ContextTokenToBuffer(
		ctx, stem, engine.LookupPaul2013EmbeddedContextToken, make([]byte, 68),
	)
	if err != nil {
		return nil, false, err
	}
	if !writes.Result.Matched || len(writes.Output) == 0 || writes.Output[0] == 0 {
		return nil, false, nil
	}
	end := bytes.IndexByte(writes.Output, 0)
	if end < 0 {
		return nil, false, errors.New("dictionary stem output is not NUL-terminated")
	}
	return append([]byte(nil), writes.Output[:end]...), true, nil
}

func finishPaul2013EdSuffix(
	result Paul2013EdSuffixResult,
	stemOutput []byte,
	currentFlags byte,
) (Paul2013EdSuffixResult, error) {
	if len(stemOutput) == 0 {
		return result, nil
	}
	last := stemOutput[len(stemOutput)-1]
	var marker []byte
	switch last {
	case '*', '5', '7', ' ', '"', '8', ':', 0x14:
		marker = []byte{'9'}
	case '9', 0x15:
		marker = []byte{'#', 0x15}
	default:
		marker = []byte{0x15}
	}
	if len(stemOutput)+len(marker) >= 68 {
		return Paul2013EdSuffixResult{}, fmt.Errorf("ed-suffix output needs %d bytes plus NUL, native local capacity is 68", len(stemOutput)+len(marker))
	}
	result.Output = append(append([]byte(nil), stemOutput...), marker...)
	result.Matched = true
	result.RowFlags = currentFlags | 0x08
	return result, nil
}
