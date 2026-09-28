// Package duration builds and evaluates the observed Paul 2013 duration inputs.
package duration

import (
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
	"vtspeak/engine/tree3"
)

// PhoneResult is one phone's duration-tree lookup. Value is the raw 16-bit
// tree output; its physical unit and conversion into the final timeline remain
// unresolved.
type PhoneResult struct {
	Phone       text.CMUPhone
	TreeName    string
	LeafOrdinal int
	Value       int16
	Input       [9]int16
}

// TokenResult keeps duration outputs aligned with the source token span.
type TokenResult struct {
	Token  text.LexicalTokenSpan
	Phones []PhoneResult
}

// Engine owns the read-only lexical and duration resources needed by the
// supported Paul 2013 text-to-duration path.
type Engine struct {
	dictionary *text.EmbeddedDictionary
	catalog    *tree3.Catalog
}

// OpenPaul2013 loads the embedded common English dictionary and Paul M16
// duration trees. The voiceDataRoot contains the voice's ttsdata/tree3 tree;
// dictionaryRoot contains the shared embedded English dictionary files.
func OpenPaul2013(voiceDataRoot, dictionaryRoot string) (*Engine, error) {
	dictionary, err := text.LoadEmbeddedDictionary(dictionaryRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 embedded dictionary: %w", err)
	}
	catalog, err := tree3.LoadPaul2013(voiceDataRoot)
	if err != nil {
		return nil, fmt.Errorf("load Paul 2013 decision trees: %w", err)
	}
	return &Engine{dictionary: dictionary, catalog: catalog}, nil
}

// Evaluate resolves source text and returns per-phone raw duration-tree
// results. It automatically handles uniquely pronounced words; ambiguous
// pronunciations fail closed until the legacy context ranker is ported.
func (engine *Engine) Evaluate(ctx context.Context, source string) ([]TokenResult, error) {
	if engine == nil {
		return nil, errors.New("Paul 2013 duration engine is nil")
	}
	return EvaluatePaul2013SourceText(ctx, source, engine.dictionary, engine.catalog)
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
	if catalog == nil {
		return nil, errors.New("Paul 2013 duration catalog is nil")
	}
	inputs, err := text.BuildPaul2013TokenDurationInputsFromText(sequence)
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
				Phone: phone, TreeName: treeName, LeafOrdinal: leaf,
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
