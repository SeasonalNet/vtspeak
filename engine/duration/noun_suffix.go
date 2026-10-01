package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013NessSuffixResult reports the bounded "ness" branch in
// FUN_1000a140. Output contains the normalized stem plus the native suffix
// marker bytes only when the embedded context normalizer succeeds.
type Paul2013NessSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// Paul2013AnceSuffixResult reports the bounded "ance" branch in
// FUN_1000a140.
type Paul2013AnceSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// Paul2013MentSuffixResult reports the bounded "ment" branch in
// FUN_1000a140.
type Paul2013MentSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// Paul2013ShipSuffixResult reports the bounded "ship" branch in
// FUN_1000a140.
type Paul2013ShipSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// Paul2013LessSuffixResult reports the bounded "less" branch in
// FUN_1000a140.
type Paul2013LessSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Stem          []byte
	Output        []byte
	RowFlags      byte
}

// NormalizePaul2013NessSuffixContext ports FUN_1000a140's 7- and 8-byte
// source branch for the directly compared "ness" suffix. It removes the
// suffix, sends the stem through FUN_1000cb30 using the loaded embedded
// dictionary, appends the observed comma/control/hyphen/9 marker bytes, and
// ORs bit 3 into the caller's row flags when the output is nonempty.
//
// Other source lengths and the remaining table-driven cases in FUN_1000a140
// are outside this helper. In the native order this branch follows the
// function's "ance" check; the whole helper is called after FUN_10009dc0 and
// before FUN_1000b4d0.
func (engine *Engine) NormalizePaul2013NessSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013NessSuffixResult, error) {
	result := Paul2013NessSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013NessSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013NessSuffixResult{}, errors.New("ness-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013NessSuffixResult{}, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 7 || len(source) > 8 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(
		source[len(source)-4:], []byte("ness"), characterMap,
	) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = stem

	output, matched, err := engine.normalizePaul2013SuffixStemWithMarker(
		ctx, stem, []byte{',', 0x07, '-', '9'},
	)
	if err != nil {
		return Paul2013NessSuffixResult{}, fmt.Errorf("normalize ness-suffix stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	result.Matched = true
	result.Output = output
	result.RowFlags = currentFlags | 0x08
	return result, nil
}

// NormalizePaul2013AnceSuffixContext ports the bounded 7- and 8-byte "ance"
// branch in FUN_1000a140. It removes the suffix, normalizes the stem, appends
// the observed control/hyphen/hash/7 marker, and sets row flag bit 3 only when
// the resulting string is nonempty.
//
// Other source lengths and table-driven cases in FUN_1000a140 remain outside
// this helper. It is the first suffix check in the function's short-input
// branch, after the preceding "liness" condition.
func (engine *Engine) NormalizePaul2013AnceSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013AnceSuffixResult, error) {
	result := Paul2013AnceSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013AnceSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013AnceSuffixResult{}, errors.New("ance-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013AnceSuffixResult{}, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 7 || len(source) > 8 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(
		source[len(source)-4:], []byte("ance"), characterMap,
	) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = stem
	output, matched, err := engine.normalizePaul2013SuffixStemWithMarker(
		ctx, stem, []byte{0x07, '-', '#', '7'},
	)
	if err != nil {
		return Paul2013AnceSuffixResult{}, fmt.Errorf("normalize ance-suffix stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	result.Matched = true
	result.Output = output
	result.RowFlags = currentFlags | 0x08
	return result, nil
}

// NormalizePaul2013MentSuffixContext ports the bounded 7- and 8-byte "ment"
// branch in FUN_1000a140. It removes the suffix, normalizes the stem, appends
// the same observed marker bytes as the adjacent "ness" branch, and sets row
// flag bit 3 only when the resulting string is nonempty.
//
// Other source lengths and table-driven cases in FUN_1000a140 remain outside
// this helper. It occupies the same fallback stage as the "ness" branch:
// after FUN_10009dc0 and before FUN_1000b4d0.
func (engine *Engine) NormalizePaul2013MentSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013MentSuffixResult, error) {
	result := Paul2013MentSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013MentSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013MentSuffixResult{}, errors.New("ment-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013MentSuffixResult{}, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 7 || len(source) > 8 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(
		source[len(source)-4:], []byte("ment"), characterMap,
	) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = stem
	output, matched, err := engine.normalizePaul2013SuffixStemWithMarker(
		ctx, stem, []byte{',', 0x07, '-', '9'},
	)
	if err != nil {
		return Paul2013MentSuffixResult{}, fmt.Errorf("normalize ment-suffix stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	result.Matched = true
	result.Output = output
	result.RowFlags = currentFlags | 0x08
	return result, nil
}

// NormalizePaul2013ShipSuffixContext ports the bounded 7- and 8-byte "ship"
// branch in FUN_1000a140. It removes the suffix, normalizes the stem, appends
// the observed control/hyphen/7 marker, and sets row flag bit 3 only when the
// resulting string is nonempty.
func (engine *Engine) NormalizePaul2013ShipSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013ShipSuffixResult, error) {
	result := Paul2013ShipSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013ShipSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013ShipSuffixResult{}, errors.New("ship-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ShipSuffixResult{}, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 7 || len(source) > 8 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(
		source[len(source)-4:], []byte("ship"), characterMap,
	) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = stem
	output, matched, err := engine.normalizePaul2013SuffixStemWithMarker(
		ctx, stem, []byte{0x07, '-', '7'},
	)
	if err != nil {
		return Paul2013ShipSuffixResult{}, fmt.Errorf("normalize ship-suffix stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	result.Matched = true
	result.Output = output
	result.RowFlags = currentFlags | 0x08
	return result, nil
}

// NormalizePaul2013LessSuffixContext ports the bounded 7- and 8-byte "less"
// branch in FUN_1000a140. It removes the suffix, normalizes the stem, appends
// the observed 8/%/5 marker, and sets row flag bit 3 only when the resulting
// string is nonempty.
func (engine *Engine) NormalizePaul2013LessSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013LessSuffixResult, error) {
	result := Paul2013LessSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if engine == nil {
		return Paul2013LessSuffixResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013LessSuffixResult{}, errors.New("less-suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013LessSuffixResult{}, err
	}
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) < 7 || len(source) > 8 {
		return result, nil
	}
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	if !text.Paul2013MappedCStringEqual(
		source[len(source)-4:], []byte("less"), characterMap,
	) {
		return result, nil
	}
	result.SuffixMatched = true
	stem := append([]byte(nil), source[:len(source)-4]...)
	result.Stem = stem
	output, matched, err := engine.normalizePaul2013SuffixStemWithMarker(
		ctx, stem, []byte{'8', '%', '5'},
	)
	if err != nil {
		return Paul2013LessSuffixResult{}, fmt.Errorf("normalize less-suffix stem %q: %w", stem, err)
	}
	if !matched {
		return result, nil
	}
	result.Matched = true
	result.Output = output
	result.RowFlags = currentFlags | 0x08
	return result, nil
}

func (engine *Engine) normalizePaul2013SuffixStemWithMarker(
	ctx context.Context,
	stem []byte,
	marker []byte,
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
		return nil, false, errors.New("suffix-stem output is not NUL-terminated")
	}
	if len(marker) > len(writes.Output)-end-1 {
		return nil, false, fmt.Errorf("suffix-marker output needs %d bytes including NUL, capacity is %d", end+len(marker)+1, len(writes.Output))
	}
	output := append([]byte(nil), writes.Output[:end]...)
	output = append(output, marker...)
	return output, true, nil
}
