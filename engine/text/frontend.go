// Package text defines the boundary between source text and model contexts.
package text

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Context carries the observed seven-byte signature. Its individual feature
// meanings remain opaque until their semantics are established.
type Context struct {
	Signature [7]byte
}

type Frontend interface {
	Process(context.Context, string) ([]Context, error)
}

// LexicalToken is one normalized dictionary surface with its original source
// token. Separators are preserved verbatim; their duration or prosodic meaning
// is not inferred.
type LexicalToken struct {
	Surface       string
	SourceSurface string
	// SourceByteStart and SourceByteEnd are inclusive offsets in the original
	// ASCII input. Numeric expansions inherit their source token's span.
	SourceByteStart   int
	SourceByteEnd     int
	HasSourceByteSpan bool
	SeparatorBefore   string
	SeparatorAfter    string
	// DictionaryMetadata retains the four embedded-payload flags copied into
	// the legacy token/context records. Their individual semantics are unknown.
	DictionaryMetadata [4]bool
	// ModelPhoneRows preserves the recovered FUN_1000d450 projection for this
	// resolved surface. The row index and surface copy are represented by the
	// token fields above; marker-branch inputs remain caller supplied.
	ModelPhoneRows Paul2013DictionaryPhoneRows
	Alternatives   []LabeledAlternative
}

// LabeledAlternative retains opaque path controls, original model-coded
// symbols, and their readable labels for one dictionary alternative.
type LabeledAlternative struct {
	Path       []byte
	PathGroups [][]byte
	Symbols    []byte
	Phones     []CMUPhone
}

// LexiconFrontend resolves ASCII surface tokens through the Paul 2013
// embedded dictionary. It expands supported numeric forms but does not
// normalize word spelling or build model contexts.
type LexiconFrontend struct {
	Dictionary *EmbeddedDictionary
}

// ResolveText expands supported numeric tokens, tokenizes the ASCII word
// stream, and resolves every normalized token to its dictionary alternatives.
// Unsupported number grammars, symbols, missing words, and private control
// symbols fail closed.
func (frontend LexiconFrontend) ResolveText(ctx context.Context, source string) ([]LexicalToken, error) {
	if frontend.Dictionary == nil {
		return nil, errors.New("lexicon frontend has no embedded dictionary")
	}
	words, err := tokenizeASCIISurfaces(source)
	if err != nil {
		return nil, err
	}
	resolved := make([]LexicalToken, 0, len(words))
	for _, word := range words {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		normalized, err := expandNumericSurface(word.Surface)
		if err != nil {
			return nil, fmt.Errorf("normalize %q: %w", word.Surface, err)
		}
		for normalizedIndex, surface := range normalized {
			pronunciation, found, err := frontend.Dictionary.ResolvePaul2013Surface([]byte(surface))
			if err != nil {
				return nil, fmt.Errorf("resolve %q as %q: %w", word.Surface, surface, err)
			}
			if !found {
				return nil, fmt.Errorf("resolve %q as %q: no embedded dictionary entry", word.Surface, surface)
			}
			if len(pronunciation.Alternatives) == 0 {
				return nil, fmt.Errorf("resolve %q as %q: dictionary entry has no pronunciation", word.Surface, surface)
			}
			modelPhoneRows, err := pronunciation.Payload.BuildPaul2013DictionaryPhoneRows(Paul2013PhoneIDCodebook())
			if err != nil {
				return nil, fmt.Errorf("build model phone rows for %q: %w", surface, err)
			}
			token := LexicalToken{
				Surface:            surface,
				SourceSurface:      word.Surface,
				SourceByteStart:    word.SourceByteStart,
				SourceByteEnd:      word.SourceByteEnd,
				HasSourceByteSpan:  true,
				DictionaryMetadata: pronunciation.Payload.Metadata,
				ModelPhoneRows:     modelPhoneRows,
			}
			if normalizedIndex == 0 {
				token.SeparatorBefore = word.SeparatorBefore
			} else {
				token.SeparatorBefore = " "
			}
			if normalizedIndex == len(normalized)-1 {
				token.SeparatorAfter = word.SeparatorAfter
			} else {
				token.SeparatorAfter = " "
			}
			for _, alternative := range pronunciation.Alternatives {
				pathGroups, err := ParsePaul2013PronunciationPath(alternative.Path)
				if err != nil {
					return nil, fmt.Errorf("resolve %q pronunciation path: %w", surface, err)
				}
				phones, err := DecodeCMUPhones(alternative.Symbols)
				if err != nil {
					return nil, fmt.Errorf("resolve %q pronunciation: %w", surface, err)
				}
				token.Alternatives = append(token.Alternatives, LabeledAlternative{
					Path:       append([]byte(nil), alternative.Path...),
					PathGroups: pathGroups,
					Symbols:    append([]byte(nil), alternative.Symbols...),
					Phones:     phones,
				})
			}
			resolved = append(resolved, token)
		}
	}
	return resolved, nil
}

// LexicalAnalysis combines resolved token alternatives with the classifier
// features that can currently be derived from token surfaces. Missing
// classifier positions remain explicit in each row's availability mask.
type LexicalAnalysis struct {
	Tokens                  []LexicalToken
	PronunciationFeatures   []Paul2013PronunciationFeatures
	PronunciationCandidates [][]Paul2013PronunciationCandidate
}

// Paul2013PronunciationCandidate is one classifier row produced from a
// token's dictionary alternative and one nonempty pronunciation path group.
type Paul2013PronunciationCandidate struct {
	AlternativeIndex int
	PathGroupIndex   int
	PathCodeIndex    int
	Features         Paul2013PronunciationFeatures
}

// AnalyzeText resolves source text and produces all currently known
// pronunciation classifier features without caller-supplied context values.
func (frontend LexiconFrontend) AnalyzeText(ctx context.Context, source string) (LexicalAnalysis, error) {
	tokens, err := frontend.ResolveText(ctx, source)
	if err != nil {
		return LexicalAnalysis{}, err
	}
	return frontend.AnalyzeResolvedTokens(ctx, tokens)
}

// AnalyzeResolvedTokens builds the currently recovered pronunciation
// classifier rows for tokens produced by another source parser, such as the
// supported inline-VTML pause parser. It does not change token spans or
// choose pronunciation alternatives.
func (frontend LexiconFrontend) AnalyzeResolvedTokens(ctx context.Context, tokens []LexicalToken) (LexicalAnalysis, error) {
	if frontend.Dictionary == nil {
		return LexicalAnalysis{}, errors.New("lexicon frontend has no embedded dictionary")
	}
	if ctx == nil {
		return LexicalAnalysis{}, errors.New("lexical analysis has no context")
	}
	if len(tokens) == 0 {
		return LexicalAnalysis{}, errors.New("lexical analysis has no resolved tokens")
	}
	analysis := LexicalAnalysis{
		Tokens:                  tokens,
		PronunciationFeatures:   make([]Paul2013PronunciationFeatures, len(tokens)),
		PronunciationCandidates: make([][]Paul2013PronunciationCandidate, len(tokens)),
	}
	for tokenIndex := range tokens {
		if err := ctx.Err(); err != nil {
			return LexicalAnalysis{}, err
		}
		features, err := BuildPaul2013KnownPronunciationFeatures(tokens, tokenIndex)
		if err != nil {
			return LexicalAnalysis{}, fmt.Errorf("build pronunciation features for %q: %w", tokens[tokenIndex].Surface, err)
		}
		analysis.PronunciationFeatures[tokenIndex] = features
		for alternativeIndex, alternative := range tokens[tokenIndex].Alternatives {
			for pathGroupIndex, pathGroup := range alternative.PathGroups {
				for pathCodeIndex, pathCode := range pathGroup {
					candidateFeatures, present, err := BuildPaul2013PronunciationPathCodeFeatures(tokens, tokenIndex, pathCode)
					if err != nil {
						return LexicalAnalysis{}, fmt.Errorf(
							"build pronunciation features for %q alternative %d path group %d code %d: %w",
							tokens[tokenIndex].Surface, alternativeIndex, pathGroupIndex, pathCodeIndex, err,
						)
					}
					if !present {
						continue
					}
					analysis.PronunciationCandidates[tokenIndex] = append(
						analysis.PronunciationCandidates[tokenIndex],
						Paul2013PronunciationCandidate{
							AlternativeIndex: alternativeIndex,
							PathGroupIndex:   pathGroupIndex,
							PathCodeIndex:    pathCodeIndex,
							Features:         candidateFeatures,
						},
					)
				}
			}
		}
	}
	return analysis, nil
}

// ResolveUniquePhoneSequence resolves text and selects the only dictionary
// pronunciation for each token. It is fully automatic when every token is
// unambiguous and fails closed when legacy context-based pronunciation
// selection is required.
func (frontend LexiconFrontend) ResolveUniquePhoneSequence(ctx context.Context, source string) (LexicalPhoneSequence, error) {
	tokens, err := frontend.ResolveText(ctx, source)
	if err != nil {
		return LexicalPhoneSequence{}, err
	}
	choices := make([]int, len(tokens))
	for _, token := range tokens {
		if len(token.Alternatives) != 1 {
			return LexicalPhoneSequence{}, fmt.Errorf(
				"choose pronunciation for %q: automatic selection among %d alternatives is not implemented",
				token.SourceSurface,
				len(token.Alternatives),
			)
		}
	}
	return SelectLexicalPronunciations(tokens, choices)
}

func expandNumericSurface(surface string) ([]string, error) {
	if surface == "" {
		return nil, errors.New("surface token is empty")
	}
	if strings.HasPrefix(surface, "$") {
		return ExpandPaul2013Currency(surface)
	}
	if strings.HasSuffix(surface, "%") {
		return ExpandPaul2013Percentage(surface)
	}
	if strings.Contains(surface, "/") {
		return ExpandPaul2013SlashDate(surface)
	}
	if strings.Contains(surface, ":") {
		return ExpandPaul2013ClockTime(surface)
	}
	if strings.Contains(surface[1:], "-") {
		return ExpandPaul2013Telephone(surface)
	}
	if len(surface) > 2 {
		suffix := surface[len(surface)-2:]
		if (suffix == "st" || suffix == "nd" || suffix == "rd" || suffix == "th") &&
			allASCIIDigits(surface[:len(surface)-2]) {
			return ExpandPaul2013Ordinal(surface)
		}
	}
	for position := 0; position < len(surface); position++ {
		if !isNumericTokenByte(surface[position]) {
			return []string{surface}, nil
		}
	}
	return ExpandPaul2013Number(surface)
}

func allASCIIDigits(value string) bool {
	if value == "" {
		return false
	}
	for position := 0; position < len(value); position++ {
		if value[position] < '0' || value[position] > '9' {
			return false
		}
	}
	return true
}

type surfaceToken struct {
	Surface         string
	SourceByteStart int
	SourceByteEnd   int
	SeparatorBefore string
	SeparatorAfter  string
}

func tokenizeASCIISurfaces(source string) ([]surfaceToken, error) {
	var tokens []surfaceToken
	var surface strings.Builder
	var separator strings.Builder
	separatorBefore := ""
	sourceByteStart := 0
	sourceByteEnd := 0
	appendSurfaceByte := func(value byte, position int) {
		if surface.Len() == 0 {
			sourceByteStart = position
		}
		surface.WriteByte(value)
		sourceByteEnd = position
	}
	flush := func() {
		if surface.Len() == 0 {
			return
		}
		tokens = append(tokens, surfaceToken{
			Surface:         surface.String(),
			SourceByteStart: sourceByteStart,
			SourceByteEnd:   sourceByteEnd,
			SeparatorBefore: separatorBefore,
		})
		surface.Reset()
		separatorBefore = ""
	}
	for position, r := range source {
		if r > 0x7f {
			return nil, fmt.Errorf("text contains unsupported non-ASCII character %q", r)
		}
		if isSurfaceByte(byte(r)) {
			if surface.Len() == 0 {
				separatorBefore = separator.String()
				separator.Reset()
			}
			appendSurfaceByte(byte(r), position)
			continue
		}
		value := byte(r)
		nextIsDigit := position+1 < len(source) && source[position+1] >= '0' && source[position+1] <= '9'
		if value == '$' && nextIsDigit {
			if surface.Len() != 0 {
				return nil, errors.New("currency sign must begin a numeric token")
			}
			separatorBefore = separator.String()
			separator.Reset()
			appendSurfaceByte(value, position)
			continue
		}
		if value == '/' && isNumericPrefix(surface.String()) && nextIsDigit {
			if strings.Count(surface.String(), "/") >= 2 {
				return nil, errors.New("slash date contains more than two separators")
			}
			appendSurfaceByte(value, position)
			continue
		}
		if value == ':' && isNumericPrefix(surface.String()) && nextIsDigit {
			if strings.Contains(surface.String(), ":") {
				return nil, errors.New("clock time contains more than one colon")
			}
			appendSurfaceByte(value, position)
			continue
		}
		if value == '-' && allASCIIDigits(surface.String()) && len(surface.String()) == 3 && nextIsDigit {
			appendSurfaceByte(value, position)
			continue
		}
		if value == '-' && nextIsDigit &&
			(allASCIIDigits(surface.String()) || strings.Contains(surface.String(), "-")) {
			return nil, errors.New("numeric separator '-' requires an unsupported telephone grammar")
		}
		if (value == '%' || value == '-' || value == '+') && isNumericPrefix(surface.String()) && nextIsDigit {
			return nil, fmt.Errorf("numeric separator %q requires an unsupported number grammar", value)
		}
		if value == '%' && isNumericPrefix(surface.String()) {
			appendSurfaceByte(value, position)
			continue
		}
		if (value == '+' || value == '-') && surface.Len() == 0 && nextIsDigit {
			separatorBefore = separator.String()
			separator.Reset()
			appendSurfaceByte(value, position)
			continue
		}
		if value == '.' && nextIsDigit && (surface.Len() == 0 || isNumericPrefix(surface.String())) {
			if surface.Len() == 0 {
				separatorBefore = separator.String()
				separator.Reset()
			}
			appendSurfaceByte(value, position)
			continue
		}
		if value == ',' && isNumericPrefix(surface.String()) &&
			(nextIsDigit || position+1 < len(source) && source[position+1] == ',') {
			appendSurfaceByte(value, position)
			continue
		}
		flush()
		separator.WriteByte(byte(r))
	}
	flush()
	if len(tokens) == 0 {
		return nil, errors.New("text contains no ASCII surface tokens")
	}
	for i := range tokens {
		if i+1 < len(tokens) {
			tokens[i].SeparatorAfter = tokens[i+1].SeparatorBefore
		} else {
			tokens[i].SeparatorAfter = separator.String()
		}
	}
	return tokens, nil
}

func isNumericPrefix(value string) bool {
	if value == "" {
		return false
	}
	if strings.HasPrefix(value, "$") {
		value = value[1:]
		if value == "" {
			return false
		}
	}
	for position := 0; position < len(value); position++ {
		if !isNumericTokenByte(value[position]) {
			return false
		}
		if (value[position] == '+' || value[position] == '-') && position != 0 {
			return false
		}
	}
	return true
}

func isNumericTokenByte(value byte) bool {
	return value >= '0' && value <= '9' || value == '+' || value == '-' || value == ',' || value == '.' || value == ':' || value == '/'
}

func isSurfaceByte(value byte) bool {
	return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z' ||
		value >= '0' && value <= '9' || value == '\''
}
