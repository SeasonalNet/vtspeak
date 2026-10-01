package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013ContextPathLookupResult records FUN_1000cbe0's dictionary gate and
// selected C-string output. Output always includes its terminating NUL.
type Paul2013ContextPathLookupResult struct {
	Source              []byte
	PathMarker          byte
	RecordMatched       bool
	PronunciationChosen bool
	Output              []byte
	ReturnCode          int32
}

// SelectPaul2013ContextPronunciationFromEmbeddedDictionary ports
// FUN_1000cbe0's dictionary-only path: resolve and parse one embedded record,
// then select the aligned phone string through FUN_10003f10 using pathMarker.
// A record with no alternatives produces an empty C string. When alternatives
// exist, the native selector's marker match and 0x17 fallback order apply.
func (engine *Engine) SelectPaul2013ContextPronunciationFromEmbeddedDictionary(
	ctx context.Context,
	source []byte,
	pathMarker byte,
) (Paul2013ContextPathLookupResult, error) {
	result := Paul2013ContextPathLookupResult{
		Source:     append([]byte(nil), source...),
		PathMarker: pathMarker,
		Output:     []byte{0},
		ReturnCode: -1,
	}
	if engine == nil || engine.dictionary == nil {
		return Paul2013ContextPathLookupResult{}, errors.New("Paul 2013 embedded context dictionary is nil or unloaded")
	}
	if ctx == nil {
		return Paul2013ContextPathLookupResult{}, errors.New("context path lookup has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContextPathLookupResult{}, err
	}
	lookupSource := source
	if nul := bytes.IndexByte(lookupSource, 0); nul >= 0 {
		lookupSource = lookupSource[:nul]
	}
	resolved, found, err := engine.dictionary.ResolvePaul2013Surface(lookupSource)
	if err != nil {
		return Paul2013ContextPathLookupResult{}, fmt.Errorf("resolve embedded context path token %q: %w", source, err)
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ContextPathLookupResult{}, err
	}
	if !found {
		return result, nil
	}
	result.RecordMatched = true
	if len(resolved.Alternatives) == 0 {
		return result, nil
	}
	paths := make([][]byte, len(resolved.Alternatives))
	phones := make([][]byte, len(resolved.Alternatives))
	for index, alternative := range resolved.Alternatives {
		paths[index] = alternative.Path
		phones[index] = alternative.Symbols
	}
	selected, chosen, err := text.SelectPaul2013PronunciationByPathMarker(
		paths, phones, pathMarker,
	)
	if err != nil {
		return Paul2013ContextPathLookupResult{}, fmt.Errorf("select embedded context pronunciation: %w", err)
	}
	if !chosen {
		return result, nil
	}
	result.PronunciationChosen = true
	result.Output = append(append([]byte(nil), selected...), 0)
	if len(selected) != 0 {
		result.ReturnCode = 1
	}
	return result, nil
}
