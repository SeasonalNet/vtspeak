// Package duration builds and evaluates the observed Paul 2013 duration inputs.
package duration

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

// PhoneResult is one phone's duration-tree lookup plus its paired pitch-tree
// results when produced by Engine.Evaluate. Value and Pitch values are raw tree
// outputs whose physical units and timeline conversions remain unresolved.
type PhoneResult struct {
	Phone       text.CMUPhone
	PhoneIndex  int
	TreeName    string
	LeafOrdinal int
	Value       int16
	Input       [9]int16
	Pitch       *PitchResult
}

// TokenResult keeps duration outputs aligned with the source token span.
type TokenResult struct {
	Token  text.LexicalTokenSpan
	Phones []PhoneResult
	// PitchBoundaryAverages has one six-value entry per adjacent within-token
	// phone pair, split into the native left-record and right-record groups.
	PitchBoundaryAverages [][2][3]int16
}

// NormalizerCharacterResult retains one character-class ATMT tree lookup from
// the generic normalization path. The class value remains an opaque signed
// short exactly as returned by the shared tree.
type NormalizerCharacterResult struct {
	Character    byte
	TreeIndex    int
	LeafOrdinal  int
	Features     text.Paul2013NormalizerFeatureWindow
	Value        int16
	CategoryByte byte
}

// NormalizerTokenResult retains each per-character trace and the shaped
// opaque class-code string produced by the following remapping stages.
type NormalizerTokenResult struct {
	Characters []NormalizerCharacterResult
	ClassCodes []byte
	Produced   bool
}

// Engine owns the read-only shared text dictionaries, pronunciation
// classifier, ATMT catalog, and Paul duration/pitch resources.
type Engine struct {
	dictionary        *text.EmbeddedDictionary
	exceptions        *text.ExceptionDictionary
	tpp               *text.TPPDictionary
	txt2              text.Tables
	catalog           *tree3.Catalog
	pronunciationTree *tree3.Tree
	atmtTrees         []*tree3.Tree
}

// OpenPaul2013 loads the shared English dictionaries and text trees plus the
// Paul M16 duration/pitch trees. The voiceDataRoot contains the voice's
// ttsdata/tree3 tree; dictionaryRoot contains the shared dict-eng resources.
func OpenPaul2013(voiceDataRoot, dictionaryRoot string) (*Engine, error) {
	dictionary, err := text.LoadEmbeddedDictionary(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 embedded dictionary: %w", err)
	}
	exceptions, err := text.LoadExceptionDictionary(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 pronunciation exceptions: %w", err)
	}
	tpp, err := text.LoadTPPDictionary(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 typed text dictionary: %w", err)
	}
	txt2, err := text.LoadTXT2Tables(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 shared text tables: %w", err)
	}
	catalog, err := tree3.LoadPaul2013(voiceDataRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 decision trees: %w", err)
	}
	pronunciationTree, err := tree3.LoadPaul2013PronunciationTree(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 pronunciation tree: %w", err)
	}
	atmtTrees, err := tree3.LoadPaul2013ATMTTrees(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 shared ATMT trees: %w", err)
	}
	catalog.Pronunciation = pronunciationTree
	return &Engine{
		dictionary:        dictionary,
		exceptions:        exceptions,
		tpp:               tpp,
		txt2:              txt2,
		catalog:           catalog,
		pronunciationTree: pronunciationTree,
		atmtTrees:         atmtTrees,
	}, nil
}

// ApplyPaul2013TPPNumericTextRows routes normalized lexical rows through the
// recovered direct F/G TPP operations and WAB class fallback. Dictionary
// index zero uses the loaded shared resources; alternate indexes remain
// unsupported until their lookup tables are available.
func (engine *Engine) ApplyPaul2013TPPNumericTextRows(
	rows []text.Paul2013TPPWindowRow,
	dictionaryIndex int,
) (text.Paul2013TPPNumericTextResult, error) {
	if engine == nil || engine.tpp == nil {
		return text.Paul2013TPPNumericTextResult{}, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	wab := engine.txt2["wab.txt2"]
	if wab == nil {
		return text.Paul2013TPPNumericTextResult{}, errors.New("Paul 2013 WAB table is nil or unloaded")
	}
	return text.ApplyPaul2013TPPNumericTextRows(text.Paul2013TPPNumericTextInput{
		Rows: rows, Dictionary: engine.tpp, WABTable: wab,
		KeyTables: text.Paul2013EmbeddedKeyTables(), DictionaryIndex: dictionaryIndex,
	})
}

// ApplyPaul2013TPPNumericParserArena projects the fields read by the native
// TPP row pass from parserRows, applies the supported F/G and WAB routing, and
// writes the recovered fields into a copy of that arena. Unknown TPP routes
// remain unmodified by the underlying fail-closed F/G dispatcher.
func (engine *Engine) ApplyPaul2013TPPNumericParserArena(
	parserRows []byte,
	dictionaryIndex int,
) ([]byte, text.Paul2013TPPNumericTextResult, error) {
	if engine == nil || engine.tpp == nil {
		return nil, text.Paul2013TPPNumericTextResult{}, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	wab := engine.txt2["wab.txt2"]
	if wab == nil {
		return nil, text.Paul2013TPPNumericTextResult{}, errors.New("Paul 2013 WAB table is nil or unloaded")
	}
	rows, err := text.Paul2013TPPWindowRowsFromParserArena(parserRows)
	if err != nil {
		return nil, text.Paul2013TPPNumericTextResult{}, fmt.Errorf("project TPP parser rows: %w", err)
	}
	return text.ApplyPaul2013TPPNumericTextRowsToParserArena(text.Paul2013TPPNumericTextInput{
		Rows: rows, Dictionary: engine.tpp, WABTable: wab,
		KeyTables: text.Paul2013EmbeddedKeyTables(), DictionaryIndex: dictionaryIndex,
	}, parserRows)
}

// LookupTypedText performs exact TPP lookup after the caller chooses a
// normalization form. Returned atom tags and suffixes retain their opaque
// data representation; they are not interpreted as phonetic or linguistic
// labels here.
func (engine *Engine) LookupTypedText(surface []byte) (text.TPPRecord, bool, error) {
	if engine == nil || engine.tpp == nil {
		return text.TPPRecord{}, false, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	return engine.tpp.Lookup(surface, text.Paul2013EmbeddedKeyTables())
}

// LookupTypedTextCode returns the raw numeric suffix for the selected F/G
// code in an exact TPP key. It leaves byte narrowing and token-row updates to
// their separately ported consumers.
func (engine *Engine) LookupTypedTextCode(surface []byte, tag byte) ([]byte, bool, error) {
	if engine == nil || engine.tpp == nil {
		return nil, false, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	return engine.tpp.LookupNumericText(surface, tag, text.Paul2013EmbeddedKeyTables())
}

// LookupSelectedTypedText returns the opaque suffix selected by one A-G TPP
// discriminator. This includes secondary atoms in two-code payloads; AX is
// returned as the literal X suffix. It does not infer how callers use the
// selected bytes.
func (engine *Engine) LookupSelectedTypedText(surface []byte, selector byte) ([]byte, bool, error) {
	if engine == nil || engine.tpp == nil {
		return nil, false, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	return engine.tpp.LookupSelectedText(surface, selector, text.Paul2013EmbeddedKeyTables())
}

// LookupSelectedTypedAtom returns the selected discriminator and copied raw
// suffix for a TPP key. The tag remains a format field; its linguistic meaning
// is not assigned.
func (engine *Engine) LookupSelectedTypedAtom(surface []byte, selector byte) (text.TPPAtom, bool, error) {
	if engine == nil || engine.tpp == nil {
		return text.TPPAtom{}, false, errors.New("Paul 2013 typed text dictionary is nil or unloaded")
	}
	return engine.tpp.LookupSelectedAtom(surface, selector, text.Paul2013EmbeddedKeyTables())
}

// LookupPaul2013CHCFlags applies the loaded chc_sort.txt2 exact-key lookup
// used by FUN_10002680. The four returned class bits and their linguistic
// roles remain opaque.
func (engine *Engine) LookupPaul2013CHCFlags(key []byte, requiredMask uint16) (int, bool, error) {
	if engine == nil {
		return 0, false, errors.New("Paul 2013 duration engine is nil")
	}
	table := engine.txt2["chc_sort.txt2"]
	if table == nil {
		return 0, false, errors.New("Paul 2013 CHC table is nil or unloaded")
	}
	return table.LookupPaul2013CHCFlags(key, requiredMask)
}

// MatchPaul2013CHCSubstring applies FUN_100026f0's substring lookup strategy
// using the loaded CHC table and mapped-character table.
func (engine *Engine) MatchPaul2013CHCSubstring(
	word []byte,
	tokenLength int,
	requiredMask uint16,
) (bool, error) {
	if engine == nil {
		return false, errors.New("Paul 2013 duration engine is nil")
	}
	table := engine.txt2["chc_sort.txt2"]
	if table == nil {
		return false, errors.New("Paul 2013 CHC table is nil or unloaded")
	}
	return table.MatchPaul2013CHCSubstring(
		word,
		tokenLength,
		requiredMask,
		text.Paul2013EmbeddedKeyTables().CharacterMap,
	)
}

// IsPaul2013NormalizerEligible ports FUN_10002c70 for ASCII letter and
// apostrophe tokens. Consonant runs use the loaded CHC matcher; the upstream
// byte buffer is 32 bytes, so longer tokens fail closed.
func (engine *Engine) IsPaul2013NormalizerEligible(source []byte) (bool, error) {
	if engine == nil {
		return false, errors.New("Paul 2013 duration engine is nil")
	}
	if nul := bytes.IndexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	if len(source) == 0 {
		return false, nil
	}
	if len(source) > 31 {
		return false, fmt.Errorf("normalizer eligibility token has %d bytes; native local run buffer supports at most 31", len(source))
	}
	table := engine.txt2["chc_sort.txt2"]
	if table == nil {
		return false, errors.New("Paul 2013 CHC table is nil or unloaded")
	}
	attributes := text.Paul2013ExceptionCharacterAttributes()
	characterMap := text.Paul2013EmbeddedKeyTables().CharacterMap
	consonants := make([]byte, 0, len(source))
	seenVowel := false
	for index, value := range source {
		if attributes[value]&0xc0 == 0 {
			if value == '\'' {
				continue
			}
			return false, nil
		}
		if paul2013NormalizerGateVowel(value) {
			seenVowel = true
			if len(consonants) >= 2 {
				mask := uint16(4)
				if index != len(consonants) {
					mask = 2
				}
				eligible, err := table.MatchPaul2013CHCSubstring(
					consonants, len(source), mask, characterMap,
				)
				if err != nil || !eligible {
					return false, err
				}
			}
			consonants = consonants[:0]
			continue
		}

		priorRunLength := len(consonants)
		consonants = append(consonants, value)
		if index == len(source)-1 && len(consonants) > 1 {
			mask := uint16(8)
			if index != priorRunLength {
				mask = 1
			}
			eligible, err := table.MatchPaul2013CHCSubstring(
				consonants, len(source), mask, characterMap,
			)
			if err != nil || !eligible {
				return false, err
			}
			consonants = consonants[:0]
		}
	}
	return seenVowel, nil
}

func paul2013NormalizerGateVowel(value byte) bool {
	if value >= 'a' && value <= 'z' {
		value -= 'a' - 'A'
	}
	return paul2013NormalizerSpecialEAdjacentClass(value) == 1
}

// LookupPronunciationException performs the recovered exact lookup for one
// already-normalized, already-encoded key. Category is the observed one-to-
// four-component group selector. Surface normalization and the private key
// transform remain explicit caller stages.
func (engine *Engine) LookupPronunciationException(category uint32, encodedKey []byte) ([]byte, bool, error) {
	if engine == nil || engine.exceptions == nil {
		return nil, false, errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	value, ok := engine.exceptions.Lookup(category, encodedKey)
	return value, ok, nil
}

// LookupPronunciationExceptionSurface ports the recovered exception-surface
// normalization and lookup when supplied the DLL's character-attribute table.
// The caller decides whether hyphens are valid in its current text context.
func (engine *Engine) LookupPronunciationExceptionSurface(
	surface []byte,
	allowHyphen bool,
	attributes [256]byte,
) ([]byte, uint32, bool, error) {
	if engine == nil || engine.exceptions == nil {
		return nil, 0, false, errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	return engine.exceptions.LookupSurface(
		surface,
		allowHyphen,
		attributes,
		text.Paul2013EmbeddedKeyTables(),
	)
}

// LookupPaul2013PronunciationExceptionSurface performs the same lookup with
// the local DLL's signed-char attribute map and key tables. High-bit source
// bytes use the verified zero prefix and are rejected by the native gate.
func (engine *Engine) LookupPaul2013PronunciationExceptionSurface(
	surface []byte,
	allowHyphen bool,
) ([]byte, uint32, bool, error) {
	if engine == nil || engine.exceptions == nil {
		return nil, 0, false, errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	return engine.exceptions.LookupPaul2013Surface(surface, allowHyphen)
}

// LookupPaul2013PronunciationExceptionSequence checks exact exception keys
// across the caller-selected consecutive surface candidates.
func (engine *Engine) LookupPaul2013PronunciationExceptionSequence(
	surfaces []string,
) (text.ExceptionMatch, bool, error) {
	return engine.LookupPaul2013PronunciationExceptionSequenceWithHPrefixRetry(surfaces, nil)
}

// LookupPaul2013PronunciationExceptionSequenceWithHPrefixRetry also accepts
// the explicit per-row source marker used by FUN_10008dc0's h'-prefixed retry.
// Marker production remains a separate frontend stage.
func (engine *Engine) LookupPaul2013PronunciationExceptionSequenceWithHPrefixRetry(
	surfaces []string,
	retryHPrefix []bool,
) (text.ExceptionMatch, bool, error) {
	if engine == nil || engine.exceptions == nil {
		return text.ExceptionMatch{}, false, errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	return engine.exceptions.LookupPaul2013SurfaceSequenceWithHPrefixRetry(surfaces, retryHPrefix)
}

// LookupPaul2013PronunciationExceptionWithDispatchGate applies the recovered
// FUN_10007520 entry gate before the exception sequence lookup. The gate's
// row fields and retry markers are explicit because their source-text
// producers remain upstream of this engine boundary.
func (engine *Engine) LookupPaul2013PronunciationExceptionWithDispatchGate(
	gate text.Paul2013ExceptionDispatchGateInput,
	surfaces []string,
	retryHPrefix []bool,
) (text.ExceptionMatch, bool, error) {
	if engine == nil || engine.exceptions == nil {
		return text.ExceptionMatch{}, false, errors.New("Paul 2013 pronunciation exceptions are nil or unloaded")
	}
	if !text.ShouldDispatchPaul2013PronunciationException(gate) {
		return text.ExceptionMatch{}, false, nil
	}
	return engine.LookupPaul2013PronunciationExceptionSequenceWithHPrefixRetry(surfaces, retryHPrefix)
}

// ResolvePaul2013PronunciationExceptionSequence looks up an adjacent surface
// sequence and decodes a matching value into destination phone rows. The
// caller supplies one component count per matched token because the native
// phone-row buffer layout is produced earlier in text processing.
func (engine *Engine) ResolvePaul2013PronunciationExceptionSequence(
	surfaces []string,
	delimiterCounts []int,
) (text.ExceptionMatch, [][]text.CMUPhone, bool, error) {
	return engine.ResolvePaul2013PronunciationExceptionSequenceWithHPrefixRetry(surfaces, nil, delimiterCounts)
}

// ResolvePaul2013PronunciationExceptionSequenceWithHPrefixRetry composes the
// explicit native retry markers with caller-supplied destination-row counts.
func (engine *Engine) ResolvePaul2013PronunciationExceptionSequenceWithHPrefixRetry(
	surfaces []string,
	retryHPrefix []bool,
	delimiterCounts []int,
) (text.ExceptionMatch, [][]text.CMUPhone, bool, error) {
	match, found, err := engine.LookupPaul2013PronunciationExceptionSequenceWithHPrefixRetry(surfaces, retryHPrefix)
	if err != nil || !found {
		return match, nil, found, err
	}
	phones, err := match.DecodePhoneRows(delimiterCounts)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, err
	}
	return match, phones, true, nil
}

// ResolvePaul2013PronunciationExceptionSurfaces composes lookup, explicit
// h'-prefix retry markers, native hyphen-count derivation, and phone-row
// decoding from the source surfaces.
func (engine *Engine) ResolvePaul2013PronunciationExceptionSurfaces(
	surfaces []string,
	retryHPrefix []bool,
) (text.ExceptionMatch, [][]text.CMUPhone, bool, error) {
	match, found, err := engine.LookupPaul2013PronunciationExceptionSequenceWithHPrefixRetry(surfaces, retryHPrefix)
	if err != nil || !found {
		return match, nil, found, err
	}
	phones, err := match.DecodePhoneRowsForSurfaces(surfaces)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, err
	}
	return match, phones, true, nil
}

// ResolvePaul2013PronunciationExceptionWithDispatchGate composes the recovered
// caller gate, prefix lookup, optional h'-prefix retry, surface-derived row
// counts, and phone-code decoding. The gate and retry-marker producers remain
// explicit inputs.
func (engine *Engine) ResolvePaul2013PronunciationExceptionWithDispatchGate(
	gate text.Paul2013ExceptionDispatchGateInput,
	surfaces []string,
	retryHPrefix []bool,
) (text.ExceptionMatch, [][]text.CMUPhone, bool, error) {
	match, found, err := engine.LookupPaul2013PronunciationExceptionWithDispatchGate(gate, surfaces, retryHPrefix)
	if err != nil || !found {
		return match, nil, found, err
	}
	phones, err := match.DecodePhoneRowsForSurfaces(surfaces)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, err
	}
	return match, phones, true, nil
}

// ResolvePaul2013PronunciationExceptionFromRows extracts the native dispatch
// gate from one source/parser row and a phone/context row, then composes the
// existing exception lookup and decoding path. The FUN_100091b0 result,
// normalized surfaces, and retry markers remain explicit inputs.
func (engine *Engine) ResolvePaul2013PronunciationExceptionFromRows(
	sourceParserRows []byte,
	sourceTokenIndex uint16,
	phoneContextRow []byte,
	contextPreparationCode int16,
	surfaces []string,
	retryHPrefix []bool,
) (text.ExceptionMatch, [][]text.CMUPhone, bool, error) {
	gate, err := text.BuildPaul2013ExceptionDispatchGateFromRows(
		sourceParserRows, sourceTokenIndex, phoneContextRow, contextPreparationCode,
	)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, err
	}
	return engine.ResolvePaul2013PronunciationExceptionWithDispatchGate(gate, surfaces, retryHPrefix)
}

// ResolvePaul2013PronunciationExceptionFromCountedRows derives the
// low-short FUN_100091b0 result from its complete return control flow, then
// resolves an exception through the counted phone/context table. The row
// mutation cascade can remain unported because this caller's gate consumes
// only that low short. preparationKnown is false only when a future supported
// input form cannot be classified by the recovered branch structure.
func (engine *Engine) ResolvePaul2013PronunciationExceptionFromCountedRows(
	sourceParserRows []byte,
	phoneContextTable []byte,
	phoneContextRowIndex int,
	surfaces []string,
	retryHPrefix []bool,
) (match text.ExceptionMatch, phones [][]text.CMUPhone, found, preparationKnown bool, err error) {
	preparation, err := text.DerivePaul2013PhoneContextPreparationFromRows(
		sourceParserRows, phoneContextTable, phoneContextRowIndex,
	)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, false, err
	}
	if !preparation.Known {
		return text.ExceptionMatch{}, nil, false, false, nil
	}
	rowStart := text.Paul2013PhoneContextTableHeaderSize + phoneContextRowIndex*text.Paul2013PhoneContextRowSize
	phoneRow := phoneContextTable[rowStart : rowStart+text.Paul2013PhoneContextRowSize]
	rowIndex := int16(binary.LittleEndian.Uint16(phoneRow[text.Paul2013PhoneContextRowTokenIndex:]))
	gate, err := text.BuildPaul2013ExceptionDispatchGateFromRows(
		sourceParserRows, uint16(rowIndex), phoneRow, preparation.Code,
	)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, true, err
	}
	match, phones, found, err = engine.ResolvePaul2013PronunciationExceptionWithDispatchGate(
		gate, surfaces, retryHPrefix,
	)
	return match, phones, found, true, err
}

// ResolvePaul2013PronunciationExceptionFromPhoneContextRows reads the
// normalized surfaces and X-row retry markers from caller-selected ordered
// phone/context rows, then applies the derived dispatch gate and exception
// lookup. candidateRowIndexes remain explicit because native row eligibility
// and sequence construction from source text are not yet recovered.
func (engine *Engine) ResolvePaul2013PronunciationExceptionFromPhoneContextRows(
	sourceParserRows []byte,
	phoneContextTable []byte,
	gateRowIndex int,
	candidateRowIndexes []int,
) (match text.ExceptionMatch, phones [][]text.CMUPhone, found, preparationKnown bool, err error) {
	inputs, err := text.ExtractPaul2013ExceptionRowInputs(phoneContextTable, candidateRowIndexes)
	if err != nil {
		return text.ExceptionMatch{}, nil, false, false, err
	}
	surfaces := make([]string, len(inputs))
	retryHPrefix := make([]bool, len(inputs))
	for index, input := range inputs {
		surfaces[index] = input.Surface
		retryHPrefix[index] = input.RetryHPrefix
	}
	return engine.ResolvePaul2013PronunciationExceptionFromCountedRows(
		sourceParserRows, phoneContextTable, gateRowIndex, surfaces, retryHPrefix,
	)
}

// EvaluateATMTTree evaluates one of the 27 shared English text trees with
// caller-supplied numeric inputs. Tree roles and feature meanings remain
// caller-defined where the native path has not established them.
func (engine *Engine) EvaluateATMTTree(
	ctx context.Context,
	treeIndex int,
	features []int16,
) (int, []int16, error) {
	if engine == nil || len(engine.atmtTrees) == 0 {
		return 0, nil, errors.New("Paul 2013 shared ATMT trees are nil or unloaded")
	}
	if ctx == nil {
		return 0, nil, errors.New("ATMT tree evaluation has no context")
	}
	if err := ctx.Err(); err != nil {
		return 0, nil, err
	}
	if treeIndex < 0 || treeIndex >= len(engine.atmtTrees) {
		return 0, nil, fmt.Errorf("ATMT tree index %d is outside 0..%d", treeIndex, len(engine.atmtTrees)-1)
	}
	leaf, output, err := engine.atmtTrees[treeIndex].Evaluate(features)
	if err != nil {
		return 0, nil, fmt.Errorf("evaluate ATMT tree %d: %w", treeIndex, err)
	}
	return leaf, output, nil
}

// EvaluatePaul2013NormalizerCharacter builds FUN_10002db0's ten-short input
// for one source byte, dispatches to the corresponding letter/apostrophe
// tree in the loaded 27-entry ATMT catalog, and returns its scalar class.
// Category bytes are supplied in their native reverse-scan order. For
// whole-token processing, EvaluatePaul2013NormalizerCharacterSequence
// produces them from right to left.
func (engine *Engine) EvaluatePaul2013NormalizerCharacter(
	ctx context.Context,
	source []byte,
	categoryBytes []byte,
	index int,
) (NormalizerCharacterResult, error) {
	if engine == nil || len(engine.atmtTrees) != 27 {
		return NormalizerCharacterResult{}, errors.New("Paul 2013 normalizer ATMT catalog is nil or incomplete")
	}
	if index < 0 || index >= len(source) {
		return NormalizerCharacterResult{}, fmt.Errorf("normalizer character index %d is outside source length %d", index, len(source))
	}
	treeIndex, err := text.Paul2013NormalizerTreeIndex(source[index])
	if err != nil {
		return NormalizerCharacterResult{}, err
	}
	features, err := text.BuildPaul2013NormalizerFeatureWindow(source, categoryBytes, index)
	if err != nil {
		return NormalizerCharacterResult{}, err
	}
	leaf, values, err := engine.EvaluateATMTTree(ctx, treeIndex, features[:])
	if err != nil {
		return NormalizerCharacterResult{}, fmt.Errorf("evaluate normalizer character tree %d: %w", treeIndex, err)
	}
	if len(values) != 1 {
		return NormalizerCharacterResult{}, fmt.Errorf("normalizer character tree %d returned %d values, want one", treeIndex, len(values))
	}
	return NormalizerCharacterResult{
		Character: source[index], TreeIndex: treeIndex,
		LeafOrdinal: leaf, Features: features, Value: values[0], CategoryByte: byte(values[0]),
	}, nil
}

// EvaluatePaul2013NormalizerCharacterSequence ports the reverse character
// classifier pass in FUN_10002f10 for ASCII letter/apostrophe tokens. It
// folds lowercase letters, carries each tree result forward as a signed-byte
// category, applies the observed short-word E override, then restores result
// order to match the native output. The native FUN_10002c70 eligibility gate
// and the later dictionary/TPP lookup are not part of this operation.
func (engine *Engine) EvaluatePaul2013NormalizerCharacterSequence(
	ctx context.Context,
	source []byte,
) ([]NormalizerCharacterResult, error) {
	if len(source) == 0 {
		return nil, errors.New("normalizer token is empty")
	}
	normalized := append([]byte(nil), source...)
	for index, value := range normalized {
		switch {
		case value >= 'a' && value <= 'z':
			normalized[index] = value - ('a' - 'A')
		case value >= 'A' && value <= 'Z', value == '\'':
		default:
			return nil, fmt.Errorf("normalizer token byte 0x%02x at position %d is outside the recovered letter/apostrophe path", value, index)
		}
	}
	results := make([]NormalizerCharacterResult, 0, len(normalized))
	categoryBytes := make([]byte, 0, len(normalized))
	for index := len(normalized) - 1; index >= 0; index-- {
		result, err := engine.EvaluatePaul2013NormalizerCharacter(
			ctx, normalized, categoryBytes, index,
		)
		if err != nil {
			return nil, fmt.Errorf("evaluate normalizer token byte %d (%q): %w", index, normalized[index], err)
		}
		if normalized[index] == 'E' && result.Value == 1 &&
			len(normalized) > 1 && len(normalized) < 4 && index > 0 && index+1 < len(normalized) &&
			paul2013NormalizerSpecialEAdjacentClass(normalized[index-1]) == 0 &&
			paul2013NormalizerSpecialEAdjacentClass(normalized[index+1]) == 0 {
			result.CategoryByte = 0x1c
		}
		results = append(results, result)
		categoryBytes = append(categoryBytes, result.CategoryByte)
	}
	for left, right := 0, len(results)-1; left < right; left, right = left+1, right-1 {
		results[left], results[right] = results[right], results[left]
	}
	return results, nil
}

// paul2013NormalizerSpecialEAdjacentClass preserves the ASCII entries used
// by the native 0x10077d9e short lookup for the letter/apostrophe path.
// Captured uppercase vowels map to one; consonants and apostrophe map to zero.
func paul2013NormalizerSpecialEAdjacentClass(value byte) int16 {
	switch value {
	case 'A', 'E', 'I', 'O', 'U', 'Y':
		return 1
	default:
		return 0
	}
}

// ClassifyPaul2013NormalizerToken composes the reverse ATMT scan with the
// observed FUN_100024f0/FUN_10002200 class-byte mapping and shaping. The native
// eligibility gate and later dispatch of the returned opaque bytes remain
// caller responsibilities.
func (engine *Engine) ClassifyPaul2013NormalizerToken(
	ctx context.Context,
	source []byte,
) (NormalizerTokenResult, error) {
	characters, err := engine.EvaluatePaul2013NormalizerCharacterSequence(ctx, source)
	if err != nil {
		return NormalizerTokenResult{}, err
	}
	categories := make([]byte, len(characters))
	for index, character := range characters {
		categories[index] = character.CategoryByte
	}
	classCodes, produced, err := text.RemapPaul2013NormalizerClassCodes(categories)
	if err != nil {
		return NormalizerTokenResult{}, fmt.Errorf("remap Paul 2013 normalizer categories: %w", err)
	}
	return NormalizerTokenResult{
		Characters: characters,
		ClassCodes: classCodes,
		Produced:   produced,
	}, nil
}

// Evaluate resolves source text, selects pronunciation alternatives through
// the loaded shared classifier, and returns per-phone duration and pitch tree
// results for the ordinary start/terminal-Z marker path.
func (engine *Engine) Evaluate(ctx context.Context, source string) ([]TokenResult, error) {
	if engine == nil {
		return nil, errors.New("Paul 2013 duration engine is nil")
	}
	sequence, err := engine.ResolvePronunciationSequence(ctx, source)
	if err != nil {
		return nil, err
	}
	terminalMarkers := make([]byte, len(sequence.Tokens))
	if len(terminalMarkers) != 0 {
		terminalMarkers[len(terminalMarkers)-1] = 'Z'
	}
	return engine.EvaluateSequenceWithMarkers(sequence, terminalMarkers)
}

// EvaluateWithMarkers resolves source text and evaluates its duration and
// pitch trees with one caller-supplied terminal marker per token. It accepts
// marker values explicitly; source-text marker production is a separate,
// unresolved part of the frontend. Pitch dispatch uses each token's marker;
// the duration boundary uses the final nonempty token's marker.
func (engine *Engine) EvaluateWithMarkers(
	ctx context.Context,
	source string,
	terminalMarkers []byte,
) ([]TokenResult, error) {
	if engine == nil {
		return nil, errors.New("Paul 2013 duration engine is nil")
	}
	sequence, err := engine.ResolvePronunciationSequence(ctx, source)
	if err != nil {
		return nil, err
	}
	return engine.EvaluateSequenceWithMarkers(sequence, terminalMarkers)
}

// EvaluateWithPhoneMarkers resolves source text and evaluates its duration
// and pitch trees with caller-supplied per-phone block bytes and per-token
// terminal markers. Per-phone source-marker production remains unresolved.
func (engine *Engine) EvaluateWithPhoneMarkers(
	ctx context.Context,
	source string,
	phoneMarkers []byte,
	terminalMarkers []byte,
) ([]TokenResult, error) {
	if engine == nil {
		return nil, errors.New("Paul 2013 duration engine is nil")
	}
	sequence, err := engine.ResolvePronunciationSequence(ctx, source)
	if err != nil {
		return nil, err
	}
	return engine.EvaluateSequenceWithPhoneMarkers(sequence, phoneMarkers, terminalMarkers)
}

// EvaluateSequenceWithMarkers evaluates already-resolved phones without
// discarding their token alignment. Terminal markers are shared by the
// duration boundary-state and pitch-tree dispatchers.
func (engine *Engine) EvaluateSequenceWithMarkers(
	sequence text.LexicalPhoneSequence,
	terminalMarkers []byte,
) ([]TokenResult, error) {
	return engine.EvaluateSequenceWithPhoneMarkers(sequence, nil, terminalMarkers)
}

// EvaluateSequenceWithPhoneMarkers evaluates an already-resolved sequence
// after rebuilding its per-token phone groups from explicit per-phone block
// markers.
func (engine *Engine) EvaluateSequenceWithPhoneMarkers(
	sequence text.LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
) ([]TokenResult, error) {
	if engine == nil || engine.catalog == nil {
		return nil, errors.New("Paul 2013 duration engine is nil or unloaded")
	}
	results, err := EvaluatePaul2013TextWithPhoneMarkers(sequence, phoneMarkers, terminalMarkers, engine.catalog)
	if err != nil {
		return nil, err
	}
	pitch, err := EvaluatePaul2013PitchTextWithPhoneMarkers(sequence, phoneMarkers, terminalMarkers, engine.catalog)
	if err != nil {
		return nil, err
	}
	return attachPaul2013PitchResults(results, pitch)
}

// EvaluateSequenceWithPositionStates evaluates duration and pitch trees for
// an already-resolved sequence using caller-supplied per-phone states. It is
// intended for captured paths whose state producer is not recovered.
func (engine *Engine) EvaluateSequenceWithPositionStates(
	sequence text.LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	positionStates []uint8,
) ([]TokenResult, error) {
	if engine == nil || engine.catalog == nil {
		return nil, errors.New("Paul 2013 duration engine is nil or unloaded")
	}
	results, err := EvaluatePaul2013TextWithPositionStates(
		sequence, phoneMarkers, terminalMarkers, positionStates, engine.catalog,
	)
	if err != nil {
		return nil, err
	}
	pitch, err := EvaluatePaul2013PitchTextWithPositionStates(
		sequence, phoneMarkers, terminalMarkers, positionStates, engine.catalog,
	)
	if err != nil {
		return nil, err
	}
	return attachPaul2013PitchResults(results, pitch)
}

// EvaluateVTMLCMUPhoneme evaluates the directly captured forced x-cmu VTML
// form using caller-supplied phone markers, terminal markers, and position
// states. It does not derive the latter state arrays or produce selected units.
func (engine *Engine) EvaluateVTMLCMUPhoneme(
	source string,
	phoneMarkers []byte,
	terminalMarkers []byte,
	positionStates []uint8,
) ([]TokenResult, error) {
	if engine == nil {
		return nil, errors.New("Paul 2013 duration engine is nil")
	}
	sequence, err := text.ParsePaul2013VTMLCMUPhoneme(source)
	if err != nil {
		return nil, fmt.Errorf("parse forced Paul 2013 phonemes: %w", err)
	}
	return engine.EvaluateSequenceWithPositionStates(
		sequence, phoneMarkers, terminalMarkers, positionStates,
	)
}

func attachPaul2013PitchResults(results []TokenResult, pitch []TokenPitchResult) ([]TokenResult, error) {
	if len(results) != len(pitch) {
		return nil, fmt.Errorf("duration produced %d tokens but pitch produced %d", len(results), len(pitch))
	}
	for tokenIndex := range results {
		if len(results[tokenIndex].Phones) != len(pitch[tokenIndex].Phones) {
			return nil, fmt.Errorf("token %d has %d duration results but %d pitch results", tokenIndex, len(results[tokenIndex].Phones), len(pitch[tokenIndex].Phones))
		}
		for phoneIndex := range results[tokenIndex].Phones {
			pitchResult := pitch[tokenIndex].Phones[phoneIndex]
			results[tokenIndex].Phones[phoneIndex].Pitch = &pitchResult
		}
		results[tokenIndex].PitchBoundaryAverages = pitch[tokenIndex].BoundaryAverages
	}
	return results, nil
}

// EvaluatePronunciationFeatures evaluates the shared engbi classifier for one
// complete legacy feature row. Engine text resolution builds these rows from
// token context and each code in a candidate path group.
func (engine *Engine) EvaluatePronunciationFeatures(features [15]int16) (int16, error) {
	if engine == nil || engine.pronunciationTree == nil {
		return 0, errors.New("Paul 2013 pronunciation tree is not loaded")
	}
	_, output, err := engine.pronunciationTree.Evaluate(features[:])
	if err != nil {
		return 0, fmt.Errorf("evaluate Paul 2013 pronunciation tree: %w", err)
	}
	if len(output) != 1 {
		return 0, fmt.Errorf("Paul 2013 pronunciation tree returned %d values, want one", len(output))
	}
	return output[0], nil
}

// EvaluatePaul2013SourceText resolves text through the local Paul dictionary,
// automatically selects only unambiguous pronunciations, derives duration
// contexts, and evaluates the selected trees. The caller supplies source text
// and already-open read-only resources, not per-phone metadata.
func EvaluatePaul2013SourceText(
	ctx context.Context,
	source string,
	dictionary *text.EmbeddedDictionary,
	catalog *tree3.Catalog,
) ([]TokenResult, error) {
	frontend := text.LexiconFrontend{Dictionary: dictionary}
	sequence, err := frontend.ResolveUniquePhoneSequence(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("resolve Paul 2013 source text: %w", err)
	}
	return EvaluatePaul2013Text(sequence, catalog)
}

// EvaluatePaul2013Text derives the duration contexts from the phone sequence,
// selects the observed tree family for each phone, and evaluates its tree.
// Inputs use the currently ported ordinary single-utterance marker path.
func EvaluatePaul2013Text(sequence text.LexicalPhoneSequence, catalog *tree3.Catalog) ([]TokenResult, error) {
	terminalMarkers := make([]byte, len(sequence.Tokens))
	if len(terminalMarkers) != 0 {
		terminalMarkers[len(terminalMarkers)-1] = 'Z'
	}
	return EvaluatePaul2013TextWithMarkers(sequence, terminalMarkers, catalog)
}

// EvaluatePaul2013TextWithMarkers evaluates duration trees using explicit
// final-token markers while retaining the ordinary phone-family dispatch.
func EvaluatePaul2013TextWithMarkers(
	sequence text.LexicalPhoneSequence,
	terminalMarkers []byte,
	catalog *tree3.Catalog,
) ([]TokenResult, error) {
	return EvaluatePaul2013TextWithPhoneMarkers(sequence, nil, terminalMarkers, catalog)
}

// EvaluatePaul2013TextWithPhoneMarkers evaluates duration trees after
// splitting each token at explicit FUN_10012c70 per-phone markers. Source
// text still does not produce those marker bytes.
func EvaluatePaul2013TextWithPhoneMarkers(
	sequence text.LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	catalog *tree3.Catalog,
) ([]TokenResult, error) {
	return evaluatePaul2013TextInputs(sequence, terminalMarkers, catalog, func() ([]text.Paul2013TokenDurationInputs, error) {
		return text.BuildPaul2013TokenDurationInputsWithPhoneMarkers(sequence, phoneMarkers, terminalMarkers)
	})
}

// EvaluatePaul2013TextWithPositionStates evaluates duration trees from the
// phone-group rows and caller-supplied per-phone position states. This keeps
// captured non-lexical state paths usable without inferring their producer.
func EvaluatePaul2013TextWithPositionStates(
	sequence text.LexicalPhoneSequence,
	phoneMarkers []byte,
	terminalMarkers []byte,
	positionStates []uint8,
	catalog *tree3.Catalog,
) ([]TokenResult, error) {
	return evaluatePaul2013TextInputs(sequence, terminalMarkers, catalog, func() ([]text.Paul2013TokenDurationInputs, error) {
		return text.BuildPaul2013TokenDurationInputsWithPositionStates(
			sequence, phoneMarkers, terminalMarkers, positionStates,
		)
	})
}

func evaluatePaul2013TextInputs(
	sequence text.LexicalPhoneSequence,
	terminalMarkers []byte,
	catalog *tree3.Catalog,
	buildInputs func() ([]text.Paul2013TokenDurationInputs, error),
) ([]TokenResult, error) {
	if catalog == nil {
		return nil, errors.New("Paul 2013 duration catalog is nil")
	}
	if len(terminalMarkers) != len(sequence.Tokens) {
		return nil, fmt.Errorf("received %d duration terminal markers for %d tokens", len(terminalMarkers), len(sequence.Tokens))
	}
	inputs, err := buildInputs()
	if err != nil {
		return nil, fmt.Errorf("build Paul 2013 duration inputs: %w", err)
	}
	results := make([]TokenResult, len(inputs))
	for tokenIndex, tokenInputs := range inputs {
		phoneResults := make([]PhoneResult, len(tokenInputs.Inputs))
		for phoneIndex, input := range tokenInputs.Inputs {
			phone := sequence.Phones[tokenInputs.Token.PhoneStart+phoneIndex]
			treeName, err := Paul2013TreeName(phone)
			if err != nil {
				return nil, fmt.Errorf("token %d phone %d: %w", tokenIndex, phoneIndex, err)
			}
			tree := catalog.Duration[treeName]
			if tree == nil {
				return nil, fmt.Errorf("Paul 2013 duration tree %q is not loaded", treeName)
			}
			leaf, output, err := tree.Evaluate(input[:])
			if err != nil {
				return nil, fmt.Errorf("evaluate %s for token %d phone %d: %w", treeName, tokenIndex, phoneIndex, err)
			}
			if len(output) != 1 {
				return nil, fmt.Errorf("duration tree %q returned %d values, want one", treeName, len(output))
			}
			phoneResults[phoneIndex] = PhoneResult{
				Phone: phone, PhoneIndex: phoneIndex, TreeName: treeName, LeafOrdinal: leaf,
				Value: output[0], Input: input,
			}
		}
		results[tokenIndex] = TokenResult{Token: tokenInputs.Token, Phones: phoneResults}
	}
	return results, nil
}

// Paul2013TreeName ports FUN_10013380's phone-family dispatch. The table
// selectors are indexed by the observed one-based phone identity in the DLL's
// DAT_1007987b/7c/80 arrays. The tree filenames are established by the model
// loader; family abbreviations remain opaque.
func Paul2013TreeName(phone text.CMUPhone) (string, error) {
	if _, err := phone.TreeFeatures(); err != nil {
		return "", err
	}
	if phone.Vowel {
		selector, ok := paul2013VowelSelectors[phone.Label]
		if !ok {
			return "", fmt.Errorf("no Paul 2013 vowel duration selector for %q", phone.Label)
		}
		name, ok := paul2013VowelTrees[selector]
		if !ok {
			return "", fmt.Errorf("Paul 2013 vowel duration selector %d for %q is unsupported", selector, phone.Label)
		}
		return name, nil
	}
	selector, ok := paul2013ConsonantSelectors[phone.Label]
	if !ok {
		return "", fmt.Errorf("no Paul 2013 consonant duration selector for %q", phone.Label)
	}
	name, ok := paul2013ConsonantTrees[selector]
	if !ok {
		return "", fmt.Errorf("Paul 2013 consonant duration selector %d for %q is unsupported", selector, phone.Label)
	}
	return name, nil
}

var paul2013VowelTrees = map[byte]string{
	1: "vshort.tree3", 2: "vlong.tree3", 3: "vdi.tree3", 4: "vsch.tree3",
}

var paul2013ConsonantTrees = map[byte]string{
	1: "cstop.tree3", 2: "cfri.tree3", 3: "caff.tree3",
	4: "cnas.tree3", 5: "capp.tree3", 6: "capp.tree3",
}

var paul2013VowelSelectors = map[string]byte{
	"AA": 2, "AE": 1, "AH": 1, "AO": 2, "AW": 3, "AY": 3,
	"EH": 1, "ER": 4, "EY": 3, "IH": 1, "IY": 2,
	"OW": 3, "OY": 3, "UH": 1, "UW": 2,
}

var paul2013ConsonantSelectors = map[string]byte{
	"B": 1, "CH": 3, "D": 1, "DH": 2, "F": 2, "G": 1,
	"HH": 2, "JH": 3, "K": 1, "L": 5, "M": 4, "N": 4,
	"NG": 4, "P": 1, "R": 6, "S": 2, "SH": 2, "T": 1,
	"TH": 2, "V": 2, "W": 6, "Y": 6, "Z": 2, "ZH": 2,
}
