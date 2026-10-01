package text

import (
	"bytes"
	"context"
	"fmt"
)

// Paul2013MarkedLexicalText retains resolved tokens and source-positioned
// mark records. It does not assign marks an audio-frame coordinate.
type Paul2013MarkedLexicalText struct {
	Tokens []LexicalToken
	Marks  []Paul2013InlineMarkEvent
}

// ResolveTextWithPaul2013InlineMarks resolves the observed unnamed and named
// self-closing mark forms while removing them from lexical input. Token spans
// use the enclosing original-source interval, including any mark tags inside
// a token; mark positions remain exact source-byte offsets.
func (frontend LexiconFrontend) ResolveTextWithPaul2013InlineMarks(
	ctx context.Context,
	source string,
) (Paul2013MarkedLexicalText, error) {
	plain, sourceOffsets, marks, err := stripPaul2013InlineMarkTags(source)
	if err != nil {
		return Paul2013MarkedLexicalText{}, err
	}
	tokens, err := frontend.ResolveText(ctx, plain)
	if err != nil {
		return Paul2013MarkedLexicalText{}, err
	}
	for tokenIndex := range tokens {
		token := &tokens[tokenIndex]
		if !token.HasSourceByteSpan || token.SourceByteStart < 0 ||
			token.SourceByteEnd < token.SourceByteStart || token.SourceByteEnd >= len(sourceOffsets) {
			return Paul2013MarkedLexicalText{}, fmt.Errorf("resolved token %d has invalid source span", tokenIndex)
		}
		start := sourceOffsets[token.SourceByteStart]
		end := sourceOffsets[token.SourceByteEnd]
		if start < 0 || end < start {
			return Paul2013MarkedLexicalText{}, fmt.Errorf("resolved token %d has no ordered original-source interval", tokenIndex)
		}
		token.SourceByteStart = start
		token.SourceByteEnd = end
	}
	return Paul2013MarkedLexicalText{Tokens: tokens, Marks: marks}, nil
}

// SelectLexicalPronunciationsWithMarks carries mark records through explicit
// alternative selection without assigning them audio-frame positions.
func SelectLexicalPronunciationsWithMarks(
	input Paul2013MarkedLexicalText,
	alternativeIndexes []int,
) (LexicalPhoneSequence, error) {
	sequence, err := SelectLexicalPronunciations(input.Tokens, alternativeIndexes)
	if err != nil {
		return LexicalPhoneSequence{}, err
	}
	sequence.InlineMarks = append([]Paul2013InlineMarkEvent(nil), input.Marks...)
	return sequence, nil
}

func stripPaul2013InlineMarkTags(source string) (string, []int, []Paul2013InlineMarkEvent, error) {
	destination := make([]byte, len(source)+1)
	destinationCursor := 0
	sourceOffsets := make([]int, 0, len(source))
	var marks []Paul2013InlineMarkEvent
	rule, err := Paul2013InlineMarkModelSourceRule(func(event Paul2013InlineMarkEvent) error {
		marks = append(marks, event)
		return nil
	})
	if err != nil {
		return "", nil, nil, err
	}
	sourceBytes := []byte(source)
	for position := 0; position < len(source); {
		relativeDelimiter := bytes.IndexAny(sourceBytes[position:], "<>")
		if relativeDelimiter < 0 {
			copy(destination[destinationCursor:], sourceBytes[position:])
			sourceOffsets = append(sourceOffsets, sourceByteOffsets(position, len(source))...)
			destinationCursor += len(source) - position
			break
		}
		delimiter := position + relativeDelimiter
		if delimiter > position {
			copy(destination[destinationCursor:], sourceBytes[position:delimiter])
			sourceOffsets = append(sourceOffsets, sourceByteOffsets(position, delimiter)...)
			destinationCursor += delimiter - position
			position = delimiter
		}
		if source[position] == '>' {
			return "", nil, nil, fmt.Errorf("unsupported markup delimiter at source byte %d", position)
		}
		result, dispatchErr := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
			Source:            sourceBytes,
			SourceCursor:      position,
			Destination:       destination,
			DestinationCursor: destinationCursor,
			Rules:             []Paul2013ModelSourceRule{rule},
		})
		if dispatchErr != nil {
			return "", nil, nil, fmt.Errorf("dispatch mark at source byte %d: %w", position, dispatchErr)
		}
		if !result.Matched || !result.HandlerAccepted {
			tagEnd := bytes.IndexByte(sourceBytes[position:], '>')
			if tagEnd < 0 {
				return "", nil, nil, fmt.Errorf("unterminated markup at source byte %d", position)
			}
			return "", nil, nil, fmt.Errorf("unsupported markup %q at source byte %d", source[position:position+tagEnd+1], position)
		}
		position = result.SourceCursor
		destination = result.Destination
		destinationCursor = result.DestinationCursor
	}
	return string(destination[:destinationCursor]), sourceOffsets, marks, nil
}
