package text

import (
	"bytes"
	"context"
	"fmt"
)

// ExpandPaul2013CapturedVTMLSubstitutions replaces the captured
// <vtml_sub alias="...">surface</vtml_sub> form with its alias. It preserves
// ordinary text around each tag and rejects other markup because its parsing
// and source-span behavior has not been recovered.
func ExpandPaul2013CapturedVTMLSubstitutions(source string) (string, error) {
	expanded, _, err := expandPaul2013CapturedVTMLSubstitutions(source)
	return expanded, err
}

func expandPaul2013CapturedVTMLSubstitutions(source string) (string, []int, error) {
	const (
		openingPrefix = `<vtml_sub alias="`
		openingEnd    = `">`
		closingTag    = `</vtml_sub>`
	)

	input := []byte(source)
	output := make([]byte, 0, len(input))
	sourceOffsets := make([]int, 0, len(input))
	for cursor := 0; cursor < len(input); {
		relativeTag := bytes.IndexByte(input[cursor:], '<')
		if relativeTag < 0 {
			output = append(output, input[cursor:]...)
			sourceOffsets = append(sourceOffsets, sourceByteOffsets(cursor, len(input))...)
			break
		}
		tagStart := cursor + relativeTag
		output = append(output, input[cursor:tagStart]...)
		sourceOffsets = append(sourceOffsets, sourceByteOffsets(cursor, tagStart)...)
		if !bytes.HasPrefix(input[tagStart:], []byte(openingPrefix)) {
			return "", nil, fmt.Errorf("unsupported markup at source byte %d", tagStart)
		}
		aliasStart := tagStart + len(openingPrefix)
		openingRelativeEnd := bytes.Index(input[aliasStart:], []byte(openingEnd))
		if openingRelativeEnd < 0 {
			return "", nil, fmt.Errorf("unterminated substitution tag at source byte %d", tagStart)
		}
		aliasEnd := aliasStart + openingRelativeEnd
		alias := input[aliasStart:aliasEnd]
		if len(alias) == 0 || bytes.ContainsAny(alias, "<>\"&") || !isPrintableASCII(alias) {
			return "", nil, fmt.Errorf("unsupported substitution alias at source byte %d", tagStart)
		}
		bodyStart := aliasEnd + len(openingEnd)
		bodyRelativeEnd := bytes.Index(input[bodyStart:], []byte(closingTag))
		if bodyRelativeEnd < 0 {
			return "", nil, fmt.Errorf("unterminated substitution body at source byte %d", tagStart)
		}
		bodyEnd := bodyStart + bodyRelativeEnd
		if bytes.ContainsAny(input[bodyStart:bodyEnd], "<>") {
			return "", nil, fmt.Errorf("nested markup in substitution body at source byte %d", bodyStart)
		}
		output = append(output, alias...)
		sourceOffsets = append(sourceOffsets, sourceByteOffsets(aliasStart, aliasEnd)...)
		cursor = bodyEnd + len(closingTag)
	}
	return string(output), sourceOffsets, nil
}

func sourceByteOffsets(start, end int) []int {
	offsets := make([]int, end-start)
	for index := range offsets {
		offsets[index] = start + index
	}
	return offsets
}

func isPrintableASCII(value []byte) bool {
	for _, character := range value {
		if character < 0x20 || character > 0x7e {
			return false
		}
	}
	return true
}

// ResolveTextWithPaul2013CapturedVTMLSubstitutions resolves the captured alias
// form through the embedded dictionary. Token source spans point into the
// alias attribute in the original input; a token that crosses a rewrite
// boundary is rejected because its single source span would be ambiguous.
func (frontend LexiconFrontend) ResolveTextWithPaul2013CapturedVTMLSubstitutions(
	ctx context.Context,
	source string,
) ([]LexicalToken, error) {
	expanded, sourceOffsets, err := expandPaul2013CapturedVTMLSubstitutions(source)
	if err != nil {
		return nil, err
	}
	tokens, err := frontend.ResolveText(ctx, expanded)
	if err != nil {
		return nil, err
	}
	for tokenIndex := range tokens {
		token := &tokens[tokenIndex]
		if !token.HasSourceByteSpan || token.SourceByteStart < 0 ||
			token.SourceByteEnd < token.SourceByteStart || token.SourceByteEnd >= len(sourceOffsets) {
			return nil, fmt.Errorf("resolved token %d has invalid expanded-text span", tokenIndex)
		}
		start := token.SourceByteStart
		end := token.SourceByteEnd
		for offset := start + 1; offset <= end; offset++ {
			if sourceOffsets[offset] != sourceOffsets[offset-1]+1 {
				return nil, fmt.Errorf("resolved token %d crosses a rewritten source boundary", tokenIndex)
			}
		}
		token.SourceByteStart = sourceOffsets[start]
		token.SourceByteEnd = sourceOffsets[end]
	}
	return tokens, nil
}
