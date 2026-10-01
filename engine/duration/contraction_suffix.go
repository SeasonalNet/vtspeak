package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013ContractionBaseResult contains the normalized prefix supplied to
// FUN_1000c040's contraction branch.
type Paul2013ContractionBaseResult struct {
	Matched  bool
	Output   []byte
	RowFlags byte
}

// Paul2013ContractionBaseNormalizer supplies FUN_1000c040's preceding prefix
// chain, which can include dictionary, exception, compound, and generic
// normalization paths.
type Paul2013ContractionBaseNormalizer func(
	context.Context,
	[]byte,
) (Paul2013ContractionBaseResult, error)

// Paul2013ContractionSuffixResult reports the bounded direct suffix cases in
// FUN_1000c040. Output excludes its trailing C NUL.
type Paul2013ContractionSuffixResult struct {
	SuffixMatched bool
	Matched       bool
	Source        []byte
	Prefix        []byte
	Suffix        string
	Output        []byte
	RowFlags      byte
}

// NormalizePaul2013ContractionSuffixContext ports FUN_1000c040's directly
// compared "'n", "'ee", "'er", and "'ers" endings and output markers. The
// source is bounded to the native 32-byte local buffer. Prefix normalization
// remains an explicit callback because that chain depends on row/model state
// and other handlers.
func NormalizePaul2013ContractionSuffixContext(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	normalize Paul2013ContractionBaseNormalizer,
) (Paul2013ContractionSuffixResult, error) {
	result := Paul2013ContractionSuffixResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if ctx == nil {
		return Paul2013ContractionSuffixResult{}, errors.New("contraction suffix normalizer has no context")
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) >= 32 {
		return result, fmt.Errorf("contraction source has %d bytes, exceeding FUN_1000c040's 32-byte local buffer", len(source))
	}
	if len(source) < 3 {
		return result, nil
	}

	suffix, suffixLength, ok := matchPaul2013ContractionSuffix(source)
	if !ok {
		return result, nil
	}
	result.SuffixMatched = true
	result.Suffix = suffix
	result.Prefix = append([]byte(nil), source[:len(source)-suffixLength]...)
	if normalize == nil {
		return Paul2013ContractionSuffixResult{}, errors.New("contraction suffix normalizer has no prefix handler")
	}

	base, err := normalize(ctx, append([]byte(nil), result.Prefix...))
	if err != nil {
		return Paul2013ContractionSuffixResult{}, fmt.Errorf("normalize contraction prefix %q: %w", result.Prefix, err)
	}
	if unsupported := base.RowFlags & ^(currentFlags | byte(0x39)); unsupported != 0 {
		return Paul2013ContractionSuffixResult{}, fmt.Errorf("contraction prefix introduced unsupported row flags %#02x", unsupported)
	}
	result.RowFlags |= base.RowFlags
	if !base.Matched || len(base.Output) == 0 {
		return result, nil
	}
	if bytes.IndexByte(base.Output, 0) >= 0 {
		return Paul2013ContractionSuffixResult{}, errors.New("contraction prefix output contains NUL")
	}

	output := append([]byte(nil), base.Output...)
	switch suffix {
	case "'n", "'N":
		last := output[len(output)-1]
		previous := byte(0)
		if len(output) > 1 {
			previous = output[len(output)-2]
		}
		lastClass := last > 0 && last <= 'E' && text.IsPaul2013ContextCodeClass(last)
		previousClass := previous > 0 && previous <= 'E' && text.IsPaul2013ContextCodeClass(previous)
		if !lastClass && (!previousClass || last != '6') {
			output = append(output, 0x07, '-')
		} else {
			output = append(output, '-')
		}
	case "'ee", "'EE":
		output = append(output, '&')
	case "'er", "'ER":
		output = append(output, 0x1a)
	case "'ers", "'ERS":
		output = append(output, 0x1a, 'D')
	}
	if len(output) >= 68 {
		return Paul2013ContractionSuffixResult{}, fmt.Errorf("contraction output needs %d bytes plus NUL, native local capacity is 68", len(output))
	}
	result.Output = output
	result.Matched = true
	result.RowFlags |= 0x08
	return result, nil
}

// NormalizePaul2013ContractionSuffixWithSupportedPrefix composes the direct
// FUN_1000c040 suffix handler with its recovered prefix chain. The model arena
// and row index are required because the E/A dictionary-record branch calls
// FUN_100086c0 with per-row state. It then tries the compound handler, the
// supported A140 subset, the generic normalizer, and finally FUN_10009cd0.
// Branches in those handlers that remain unresolved return an error instead
// of being treated as a miss.
func (engine *Engine) NormalizePaul2013ContractionSuffixWithSupportedPrefix(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	outerMode bool,
	model []byte,
	contextRowIndex int,
) (Paul2013ContractionSuffixResult, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013ContractionSuffixResult{}, errors.New("Paul 2013 contraction resources are nil or unloaded")
	}
	normalize := func(ctx context.Context, prefix []byte) (Paul2013ContractionBaseResult, error) {
		lookup, resultType, err := engine.lookupPaul2013EmbeddedContextTokenWithType(ctx, prefix)
		if err != nil {
			return Paul2013ContractionBaseResult{}, fmt.Errorf("look up contraction prefix: %w", err)
		}
		if resultType == 'E' || resultType == 'A' {
			record, err := engine.BuildPaul2013C3A0AlternateRecord(ctx, prefix)
			if err != nil {
				return Paul2013ContractionBaseResult{}, fmt.Errorf("build contraction prefix alternate record: %w", err)
			}
			decision, _, err := engine.EvaluatePaul2013FUN100086C0SupportedPath(
				ctx, prefix, model, contextRowIndex, record, true,
			)
			if err != nil {
				return Paul2013ContractionBaseResult{}, fmt.Errorf("evaluate contraction prefix FUN_100086c0: %w", err)
			}
			if decision.Disposition == Paul2013FUN100086C0NeedsNeighborScan {
				return Paul2013ContractionBaseResult{}, errors.New("contraction prefix FUN_100086c0 did not reach a native low-short return")
			}
			if decision.NativeShort != 0 {
				output, err := text.EncodePaul2013ContextString(prefix)
				if err != nil {
					return Paul2013ContractionBaseResult{}, fmt.Errorf("encode contraction prefix context: %w", err)
				}
				flags := byte(0)
				if len(output) != 0 {
					flags = 0x20
				}
				return Paul2013ContractionBaseResult{Matched: len(output) != 0, Output: output, RowFlags: flags}, nil
			}
		}
		if lookup.RecordCount > 0 {
			if bytes.IndexByte(lookup.DictionaryText, 0) >= 0 {
				return Paul2013ContractionBaseResult{}, errors.New("contraction prefix dictionary text contains NUL")
			}
			return Paul2013ContractionBaseResult{
				Matched:  len(lookup.DictionaryText) != 0,
				Output:   append([]byte(nil), lookup.DictionaryText...),
				RowFlags: currentFlags | 0x01,
			}, nil
		}

		compound, err := engine.NormalizePaul2013CompoundContextWithSupportedHandlers(
			ctx, prefix, currentFlags, outerMode, model, contextRowIndex,
			func(ctx context.Context, nested []byte, nestedOuterMode bool) (Paul2013ContractionBaseResult, error) {
				result, err := engine.NormalizePaul2013ContractionSuffixWithSupportedPrefix(
					ctx, nested, currentFlags, nestedOuterMode, model, contextRowIndex,
				)
				return Paul2013ContractionBaseResult{
					Matched: result.Matched, Output: result.Output, RowFlags: result.RowFlags,
				}, err
			},
		)
		if err != nil {
			return Paul2013ContractionBaseResult{}, fmt.Errorf("normalize contraction prefix compound chain: %w", err)
		}
		if compound.Matched {
			return Paul2013ContractionBaseResult{
				Matched: true, Output: append([]byte(nil), compound.Output...), RowFlags: compound.RowFlags,
			}, nil
		}

		a140, err := engine.NormalizePaul2013SupportedA140Suffixes(ctx, prefix, currentFlags)
		if err != nil {
			return Paul2013ContractionBaseResult{}, fmt.Errorf("normalize contraction prefix A140 chain: %w", err)
		}
		if a140.UnsupportedBranch || a140.TableClassUnsupported {
			return Paul2013ContractionBaseResult{}, fmt.Errorf("contraction prefix selected unresolved A140 %q branch", a140.Handler)
		}
		if a140.Matched {
			return Paul2013ContractionBaseResult{
				Matched: true, Output: append([]byte(nil), a140.Output...), RowFlags: a140.RowFlags,
			}, nil
		}

		generic, err := engine.NormalizePaul2013GenericToken(ctx, prefix)
		if err != nil {
			return Paul2013ContractionBaseResult{}, fmt.Errorf("normalize contraction prefix generic chain: %w", err)
		}
		if generic.Eligible {
			output := generic.ClassCodes
			if generic.UsedContextFallback {
				output = generic.ContextCodes
			}
			return Paul2013ContractionBaseResult{
				Matched: len(output) != 0, Output: append([]byte(nil), output...),
				RowFlags: currentFlags | generic.RowFlagMask,
			}, nil
		}
		output, err := text.EncodePaul2013ContextString(prefix)
		if err != nil {
			return Paul2013ContractionBaseResult{}, fmt.Errorf("encode contraction prefix fallback: %w", err)
		}
		flags := currentFlags
		if len(output) != 0 {
			flags |= 0x20
		}
		return Paul2013ContractionBaseResult{Matched: len(output) != 0, Output: output, RowFlags: flags}, nil
	}
	return NormalizePaul2013ContractionSuffixContext(ctx, source, currentFlags, normalize)
}

func matchPaul2013ContractionSuffix(source []byte) (string, int, bool) {
	for _, suffix := range []string{"'n", "'N"} {
		if len(source) >= 3 && bytes.HasSuffix(source, []byte(suffix)) {
			return suffix, len(suffix), true
		}
	}
	if len(source) >= 4 {
		for _, suffix := range []string{"'ee", "'EE", "'er", "'ER"} {
			if bytes.HasSuffix(source, []byte(suffix)) {
				return suffix, len(suffix), true
			}
		}
	}
	if len(source) >= 5 {
		for _, suffix := range []string{"'ers", "'ERS"} {
			if bytes.HasSuffix(source, []byte(suffix)) {
				return suffix, len(suffix), true
			}
		}
	}
	return "", 0, false
}
