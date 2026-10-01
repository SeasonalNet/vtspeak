package duration

import (
	"context"
	"errors"
)

// Paul2013SuffixCascadeResult records the first recovered suffix handler used
// after FUN_10009dc0 has already missed. The function does not substitute for
// that preceding handler or for unrecovered FUN_1000a140 cases.
type Paul2013SuffixCascadeResult struct {
	Handler               string
	SuffixMatched         bool
	Matched               bool
	UnsupportedBranch     bool
	Source                []byte
	Stem                  []byte
	Transformed           []byte
	Output                []byte
	RowFlags              byte
	TableClassUnsupported bool
}

// NormalizePaul2013SupportedSuffixFallbackAfter09DC0 composes the directly
// ported suffix handlers in their recovered order: the bounded
// FUN_1000a140 branches first, then FUN_1000b4d0's mapped spelling rewrites.
// Callers must invoke it only after FUN_10009dc0 misses. Unported
// FUN_1000a140 branches remain outside this composition.
func (engine *Engine) NormalizePaul2013SupportedSuffixFallbackAfter09DC0(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013SuffixCascadeResult, error) {
	return engine.normalizePaul2013SupportedSuffixFallback(ctx, source, currentFlags, true)
}

// NormalizePaul2013SupportedA140Suffixes composes the directly ported bounded
// FUN_1000a140 suffix branches, including the short `ers`, `est`, `ing`,
// `ist`, `ful`, `ly`, and terminal-`s` cases. Call it after FUN_100086c0 and
// before the later C3A0 suffix and generic handlers. It excludes
// FUN_1000b4d0's mapped spelling rewrites.
func (engine *Engine) NormalizePaul2013SupportedA140Suffixes(
	ctx context.Context,
	source []byte,
	currentFlags byte,
) (Paul2013SuffixCascadeResult, error) {
	return engine.normalizePaul2013SupportedSuffixFallback(ctx, source, currentFlags, false)
}

func (engine *Engine) normalizePaul2013SupportedSuffixFallback(
	ctx context.Context,
	source []byte,
	currentFlags byte,
	includeMappedSpelling bool,
) (Paul2013SuffixCascadeResult, error) {
	if engine == nil {
		return Paul2013SuffixCascadeResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013SuffixCascadeResult{}, errors.New("suffix cascade has no context")
	}
	result := Paul2013SuffixCascadeResult{
		Source: append([]byte(nil), source...), RowFlags: currentFlags,
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}

	length := cStringLength(source)
	var attempted *Paul2013SuffixCascadeResult
	normalizeS := func() (*Paul2013SuffixCascadeResult, error) {
		value, err := engine.NormalizePaul2013SSuffixContext(ctx, source, currentFlags)
		if err != nil {
			return nil, err
		}
		if !value.SuffixMatched {
			return nil, nil
		}
		cascade := suffixCascadeResult("s", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags)
		cascade.TableClassUnsupported = value.TableClassUnsupported
		return &cascade, nil
	}
	normalizeLy := func() (*Paul2013SuffixCascadeResult, error) {
		value, err := engine.NormalizePaul2013LySuffixContext(ctx, source, currentFlags)
		if err != nil {
			return nil, err
		}
		if !value.SuffixMatched {
			return nil, nil
		}
		cascade := suffixCascadeResult("ly", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags)
		cascade.TableClassUnsupported = value.TableClassUnsupported
		return &cascade, nil
	}
	normalizeIng := func() (*Paul2013SuffixCascadeResult, error) {
		value, err := engine.NormalizePaul2013IngSuffixContext(ctx, source, currentFlags)
		if err != nil {
			return nil, err
		}
		if !value.SuffixMatched {
			return nil, nil
		}
		cascade := suffixCascadeResult("ing", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags)
		cascade.TableClassUnsupported = value.TableClassUnsupported
		return &cascade, nil
	}
	normalizeIst := func() (*Paul2013SuffixCascadeResult, error) {
		value, err := engine.NormalizePaul2013IstSuffixContext(ctx, source, currentFlags)
		if err != nil {
			return nil, err
		}
		if !value.SuffixMatched {
			return nil, nil
		}
		cascade := suffixCascadeResult("ist", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags)
		return &cascade, nil
	}
	normalizeFul := func() (*Paul2013SuffixCascadeResult, error) {
		value, err := engine.NormalizePaul2013FulSuffixContext(ctx, source, currentFlags)
		if err != nil {
			return nil, err
		}
		if !value.SuffixMatched {
			return nil, nil
		}
		cascade := suffixCascadeResult("ful", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags)
		return &cascade, nil
	}
	if length >= 7 && length <= 8 {
		if value, err := engine.NormalizePaul2013AnceSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ance", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			if value.Matched {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013NessSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ness", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			if value.Matched {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013MentSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ment", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			if value.Matched {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013ShipSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ship", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			if value.Matched {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013LessSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("less", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			if value.Matched {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013ErsSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ers", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			return *attempted, nil
		}
		if value, err := engine.NormalizePaul2013EstSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("est", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			return *attempted, nil
		}
		if value, err := normalizeIng(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeIst(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeFul(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if attempted == nil {
			if value, err := normalizeLy(); err != nil {
				return Paul2013SuffixCascadeResult{}, err
			} else if value != nil {
				return *value, nil
			}
			if value, err := normalizeS(); err != nil {
				return Paul2013SuffixCascadeResult{}, err
			} else if value != nil {
				return *value, nil
			}
		}
	} else if length == 6 {
		if value, err := engine.NormalizePaul2013ErsSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ers", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			return *attempted, nil
		}
		if value, err := engine.NormalizePaul2013EstSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("est", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			return *attempted, nil
		}
		if value, err := normalizeIng(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeIst(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeFul(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if attempted == nil {
			if value, err := normalizeLy(); err != nil {
				return Paul2013SuffixCascadeResult{}, err
			} else if value != nil {
				return *value, nil
			}
			if value, err := normalizeS(); err != nil {
				return Paul2013SuffixCascadeResult{}, err
			} else if value != nil {
				return *value, nil
			}
		}
	} else if length == 5 {
		if value, err := normalizeIng(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeIst(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeFul(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeLy(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeS(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
	} else if length >= 9 && length <= 31 {
		if value, err := engine.NormalizePaul2013LinessSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("liness", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			if value.Matched || value.TableClassUnsupported {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013EdSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ed", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			if value.Matched || value.TableClassUnsupported {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013ErSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("er", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			if value.Matched || value.TableClassUnsupported {
				return *attempted, nil
			}
		}
		if value, err := engine.NormalizePaul2013LongIngSuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ing", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			return *attempted, nil
		}
		if value, err := normalizeIst(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := normalizeFul(); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value != nil {
			return *value, nil
		}
		if value, err := engine.NormalizePaul2013LongLySuffixContext(ctx, source, currentFlags); err != nil {
			return Paul2013SuffixCascadeResult{}, err
		} else if value.SuffixMatched {
			attempted = ptrSuffixCascadeResult(suffixCascadeResult("ly", value.SuffixMatched, value.Matched, value.Source, value.Stem, nil, value.Output, value.RowFlags))
			attempted.TableClassUnsupported = value.TableClassUnsupported
			if value.Matched || value.TableClassUnsupported {
				return *attempted, nil
			}
		}
		if attempted == nil {
			if value, err := normalizeS(); err != nil {
				return Paul2013SuffixCascadeResult{}, err
			} else if value != nil {
				return *value, nil
			}
		}
	}
	if !includeMappedSpelling {
		if attempted != nil {
			return *attempted, nil
		}
		return result, nil
	}

	value, err := engine.NormalizePaul2013MappedSpellingSuffix(ctx, source, currentFlags)
	if err != nil {
		return Paul2013SuffixCascadeResult{}, err
	}
	if value.SuffixMatched {
		result := suffixCascadeResult("mapped-spelling", value.SuffixMatched, value.Matched, value.Source, nil, value.Transformed, value.Output, value.RowFlags)
		if value.Matched {
			return result, nil
		}
		attempted = &result
	}
	if attempted != nil {
		return *attempted, nil
	}
	return result, nil
}

func suffixCascadeResult(
	handler string,
	suffixMatched bool,
	matched bool,
	source []byte,
	stem []byte,
	transformed []byte,
	output []byte,
	flags byte,
) Paul2013SuffixCascadeResult {
	return Paul2013SuffixCascadeResult{
		Handler: handler, SuffixMatched: suffixMatched, Matched: matched,
		Source: append([]byte(nil), source...), Stem: append([]byte(nil), stem...),
		Transformed: append([]byte(nil), transformed...), Output: append([]byte(nil), output...),
		RowFlags: flags,
	}
}

func cStringLength(source []byte) int {
	for index, value := range source {
		if value == 0 {
			return index
		}
	}
	return len(source)
}

func ptrSuffixCascadeResult(result Paul2013SuffixCascadeResult) *Paul2013SuffixCascadeResult {
	return &result
}
