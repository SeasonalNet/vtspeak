package text

import (
	"bytes"
	"context"
	"fmt"
)

// Paul2013CapturedVTMLText combines the three narrowly observed VTML forms
// that have standalone parsers. Mark and pause offsets refer to the original
// source; token spans are mapped back to original source bytes.
type Paul2013CapturedVTMLText struct {
	Tokens []LexicalToken
	Pauses []Paul2013InlinePause
	Marks  []Paul2013InlineMarkEvent
}

// SelectPaul2013CapturedVTMLPronunciations flattens the selected token paths
// and retains both supported inline event streams without assigning marks an
// output-frame location.
func SelectPaul2013CapturedVTMLPronunciations(
	input Paul2013CapturedVTMLText,
	alternativeIndexes []int,
) (LexicalPhoneSequence, error) {
	sequence, err := SelectLexicalPronunciationsWithPauses(
		Paul2013PausedLexicalText{Tokens: input.Tokens, Pauses: input.Pauses}, alternativeIndexes,
	)
	if err != nil {
		return LexicalPhoneSequence{}, err
	}
	sequence.InlineMarks = append([]Paul2013InlineMarkEvent(nil), input.Marks...)
	return sequence, nil
}

// ResolveTextWithPaul2013CapturedVTML resolves ordinary text containing the
// captured pause, mark, and substitution forms. It does not assign marks an
// output-frame location or interpret unsupported markup.
func (frontend LexiconFrontend) ResolveTextWithPaul2013CapturedVTML(
	ctx context.Context,
	source string,
) (Paul2013CapturedVTMLText, error) {
	plain, offsets, substitutions, pauses, marks, err := stripPaul2013CapturedVTML(source)
	if err != nil {
		return Paul2013CapturedVTMLText{}, err
	}
	tokens, err := frontend.ResolveText(ctx, plain)
	if err != nil {
		return Paul2013CapturedVTMLText{}, err
	}
	for tokenIndex := range tokens {
		token := &tokens[tokenIndex]
		if !token.HasSourceByteSpan || token.SourceByteStart < 0 || token.SourceByteEnd < token.SourceByteStart || token.SourceByteEnd >= len(offsets) {
			return Paul2013CapturedVTMLText{}, fmt.Errorf("resolved token %d has invalid transformed-text span", tokenIndex)
		}
		start, end := token.SourceByteStart, token.SourceByteEnd
		for position := start + 1; position <= end; position++ {
			if offsets[position] <= offsets[position-1] {
				return Paul2013CapturedVTMLText{}, fmt.Errorf("resolved token %d maps to a non-increasing source span", tokenIndex)
			}
			if offsets[position] != offsets[position-1]+1 && crossesSubstitutionBoundary(offsets[position-1], offsets[position], substitutions) {
				return Paul2013CapturedVTMLText{}, fmt.Errorf("resolved token %d crosses a rewritten substitution boundary", tokenIndex)
			}
		}
		token.SourceByteStart = offsets[start]
		token.SourceByteEnd = offsets[end]
	}
	resolvedPauses := make([]Paul2013InlinePause, len(pauses))
	for index, pause := range pauses {
		after := 0
		for _, token := range tokens {
			if token.SourceByteEnd >= pause.offset {
				break
			}
			after++
		}
		resolvedPauses[index] = Paul2013InlinePause{
			SourceByteOffset: pause.offset, AfterToken: after,
			DurationMilliseconds: pause.duration,
			OutputFramesAt16KHz:  uint64(pause.duration) * 16,
		}
	}
	return Paul2013CapturedVTMLText{Tokens: tokens, Pauses: resolvedPauses, Marks: marks}, nil
}

type capturedVTMLSubstitutionRange struct{ aliasStart, aliasEnd int }
type capturedVTMLPause struct {
	offset   int
	duration uint32
}

func crossesSubstitutionBoundary(left, right int, substitutions []capturedVTMLSubstitutionRange) bool {
	for _, span := range substitutions {
		if left < span.aliasStart && right >= span.aliasStart || left < span.aliasEnd && right > span.aliasEnd {
			return true
		}
	}
	return false
}

func stripPaul2013CapturedVTML(source string) (string, []int, []capturedVTMLSubstitutionRange, []capturedVTMLPause, []Paul2013InlineMarkEvent, error) {
	input := []byte(source)
	output := make([]byte, 0, len(input))
	offsets := make([]int, 0, len(input))
	var substitutions []capturedVTMLSubstitutionRange
	var pauses []capturedVTMLPause
	var marks []Paul2013InlineMarkEvent
	for cursor := 0; cursor < len(input); {
		relative := bytes.IndexByte(input[cursor:], '<')
		if relative < 0 {
			if unmatched := bytes.IndexByte(input[cursor:], '>'); unmatched >= 0 {
				return "", nil, nil, nil, nil, fmt.Errorf("unsupported markup delimiter at source byte %d", cursor+unmatched)
			}
			output = append(output, input[cursor:]...)
			offsets = append(offsets, sourceByteOffsets(cursor, len(input))...)
			break
		}
		start := cursor + relative
		if unmatched := bytes.IndexByte(input[cursor:start], '>'); unmatched >= 0 {
			return "", nil, nil, nil, nil, fmt.Errorf("unsupported markup delimiter at source byte %d", cursor+unmatched)
		}
		output = append(output, input[cursor:start]...)
		offsets = append(offsets, sourceByteOffsets(cursor, start)...)
		if bytes.HasPrefix(input[start:], []byte(`<vtml_sub alias="`)) {
			const prefix, opener, closer = `<vtml_sub alias="`, `">`, `</vtml_sub>`
			aliasStart := start + len(prefix)
			relEnd := bytes.Index(input[aliasStart:], []byte(opener))
			if relEnd < 0 {
				return "", nil, nil, nil, nil, fmt.Errorf("unterminated substitution tag at source byte %d", start)
			}
			aliasEnd := aliasStart + relEnd
			alias := input[aliasStart:aliasEnd]
			if len(alias) == 0 || bytes.ContainsAny(alias, "<>\"&") || !isPrintableASCII(alias) {
				return "", nil, nil, nil, nil, fmt.Errorf("unsupported substitution alias at source byte %d", start)
			}
			bodyStart := aliasEnd + len(opener)
			relBodyEnd := bytes.Index(input[bodyStart:], []byte(closer))
			if relBodyEnd < 0 {
				return "", nil, nil, nil, nil, fmt.Errorf("unterminated substitution body at source byte %d", start)
			}
			bodyEnd := bodyStart + relBodyEnd
			if bytes.ContainsAny(input[bodyStart:bodyEnd], "<>") {
				return "", nil, nil, nil, nil, fmt.Errorf("nested markup in substitution body at source byte %d", bodyStart)
			}
			output = append(output, alias...)
			offsets = append(offsets, sourceByteOffsets(aliasStart, aliasEnd)...)
			substitutions = append(substitutions, capturedVTMLSubstitutionRange{aliasStart: aliasStart, aliasEnd: aliasEnd})
			cursor = bodyEnd + len(closer)
			continue
		}
		relEnd := bytes.IndexByte(input[start:], '>')
		if relEnd < 0 {
			return "", nil, nil, nil, nil, fmt.Errorf("unterminated VTML tag at source byte %d", start)
		}
		end := start + relEnd + 1
		tag := input[start:end]
		if bytes.HasPrefix(tag, []byte("<vtml_pause")) {
			duration, err := parsePaul2013InlinePauseTag(string(tag), start)
			if err != nil {
				return "", nil, nil, nil, nil, err
			}
			pauses = append(pauses, capturedVTMLPause{offset: start, duration: duration})
			output = append(output, ' ')
			offsets = append(offsets, -1)
		} else if bytes.HasPrefix(tag, []byte("<vtml_mark")) {
			event, emit, err := parsePaul2013InlineMarkTag(tag, start)
			if err != nil {
				return "", nil, nil, nil, nil, err
			}
			if emit {
				marks = append(marks, event)
			}
		} else {
			return "", nil, nil, nil, nil, fmt.Errorf("unsupported markup at source byte %d", start)
		}
		cursor = end
	}
	return string(output), offsets, substitutions, pauses, marks, nil
}
