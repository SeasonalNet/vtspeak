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
	Surface         string
	SourceSurface   string
	SeparatorBefore string
	SeparatorAfter  string
	// DictionaryMetadata retains the four embedded-payload flags copied into
	// the legacy token/context records. Their individual semantics are unknown.
	DictionaryMetadata [4]bool
	Alternatives       []LabeledAlternative
}

// LabeledAlternative retains opaque path controls, original model-coded
// symbols, and their readable labels for one dictionary alternative.
type LabeledAlternative struct {
	Path    []byte
	Symbols []byte
	Phones  []CMUPhone
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
			token := LexicalToken{
				Surface:            surface,
				SourceSurface:      word.Surface,
				DictionaryMetadata: pronunciation.Payload.Metadata,
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
				phones, err := DecodeCMUPhones(alternative.Symbols)
				if err != nil {
					return nil, fmt.Errorf("resolve %q pronunciation: %w", surface, err)
				}
				token.Alternatives = append(token.Alternatives, LabeledAlternative{
					Path:    append([]byte(nil), alternative.Path...),
					Symbols: append([]byte(nil), alternative.Symbols...),
					Phones:  phones,
				})
			}
			resolved = append(resolved, token)
		}
	}
	return resolved, nil
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
	SeparatorBefore string
	SeparatorAfter  string
}

func tokenizeASCIISurfaces(source string) ([]surfaceToken, error) {
	var tokens []surfaceToken
	var surface strings.Builder
	var separator strings.Builder
	separatorBefore := ""
	flush := func() {
		if surface.Len() == 0 {
			return
		}
		tokens = append(tokens, surfaceToken{Surface: surface.String(), SeparatorBefore: separatorBefore})
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
			surface.WriteByte(byte(r))
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
			surface.WriteByte(value)
			continue
		}
		if value == '/' && isNumericPrefix(surface.String()) && nextIsDigit {
			if strings.Count(surface.String(), "/") >= 2 {
				return nil, errors.New("slash date contains more than two separators")
			}
			surface.WriteByte(value)
			continue
		}
		if value == ':' && isNumericPrefix(surface.String()) && nextIsDigit {
			if strings.Contains(surface.String(), ":") {
				return nil, errors.New("clock time contains more than one colon")
			}
			surface.WriteByte(value)
			continue
		}
		if (value == '%' || value == '-' || value == '+') && isNumericPrefix(surface.String()) && nextIsDigit {
			return nil, fmt.Errorf("numeric separator %q requires an unsupported number grammar", value)
		}
		if value == '%' && isNumericPrefix(surface.String()) {
			surface.WriteByte(value)
			continue
		}
		if (value == '+' || value == '-') && surface.Len() == 0 && nextIsDigit {
			separatorBefore = separator.String()
			separator.Reset()
			surface.WriteByte(value)
			continue
		}
		if value == '.' && nextIsDigit && (surface.Len() == 0 || isNumericPrefix(surface.String())) {
			if surface.Len() == 0 {
				separatorBefore = separator.String()
				separator.Reset()
			}
			surface.WriteByte(value)
			continue
		}
		if value == ',' && isNumericPrefix(surface.String()) &&
			(nextIsDigit || position+1 < len(source) && source[position+1] == ',') {
			surface.WriteByte(value)
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
