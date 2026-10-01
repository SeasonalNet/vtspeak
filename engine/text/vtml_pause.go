package text

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Paul2013InlinePause is one observed VTML pause boundary. AfterToken is the
// number of resolved lexical tokens preceding the tag, so zero denotes a
// leading pause and len(Tokens) denotes a trailing pause.
type Paul2013InlinePause struct {
	SourceByteOffset     int
	AfterToken           int
	DurationMilliseconds uint32
	OutputFramesAt16KHz  uint64
}

// Paul2013PausedLexicalText retains resolved tokens and the inline silence
// boundaries found in their source. These events are not yet wired into the
// unit-selection and timeline pipeline.
type Paul2013PausedLexicalText struct {
	Tokens []LexicalToken
	Pauses []Paul2013InlinePause
}

// SelectLexicalPronunciationsWithPauses carries the ordered pause events into
// the flattened phone sequence without assigning them acoustic effects.
func SelectLexicalPronunciationsWithPauses(
	input Paul2013PausedLexicalText,
	alternativeIndexes []int,
) (LexicalPhoneSequence, error) {
	sequence, err := SelectLexicalPronunciations(input.Tokens, alternativeIndexes)
	if err != nil {
		return LexicalPhoneSequence{}, err
	}
	for index, pause := range input.Pauses {
		if pause.AfterToken < 0 || pause.AfterToken > len(sequence.Tokens) {
			return LexicalPhoneSequence{}, fmt.Errorf("inline pause %d follows token boundary %d outside [0, %d]", index, pause.AfterToken, len(sequence.Tokens))
		}
		if pause.OutputFramesAt16KHz != uint64(pause.DurationMilliseconds)*16 {
			return LexicalPhoneSequence{}, fmt.Errorf("inline pause %d has inconsistent 16 kHz frame count", index)
		}
	}
	sequence.InlinePauses = append([]Paul2013InlinePause(nil), input.Pauses...)
	return sequence, nil
}

// ResolveTextWithPaul2013InlinePauses resolves ordinary text plus the exact
// self-closing <vtml_pause time="N"/> form observed in Stage 16. Each tag
// becomes a token boundary and an output-frame count at 16 kHz. Other markup
// fails closed because its source and audio behavior are not implemented.
func (frontend LexiconFrontend) ResolveTextWithPaul2013InlinePauses(
	ctx context.Context,
	source string,
) (Paul2013PausedLexicalText, error) {
	plain, sourceOffsets, rawPauses, err := stripPaul2013InlinePauseTags(source)
	if err != nil {
		return Paul2013PausedLexicalText{}, err
	}
	tokens, err := frontend.ResolveText(ctx, plain)
	if err != nil {
		return Paul2013PausedLexicalText{}, err
	}
	for tokenIndex := range tokens {
		token := &tokens[tokenIndex]
		if !token.HasSourceByteSpan || token.SourceByteStart < 0 || token.SourceByteEnd < token.SourceByteStart || token.SourceByteEnd >= len(sourceOffsets) {
			return Paul2013PausedLexicalText{}, fmt.Errorf("resolved token %d has invalid source span", tokenIndex)
		}
		originalStart := sourceOffsets[token.SourceByteStart]
		originalEnd := sourceOffsets[token.SourceByteEnd]
		if originalStart < 0 || originalEnd < originalStart {
			return Paul2013PausedLexicalText{}, fmt.Errorf("resolved token %d crosses a synthetic markup boundary", tokenIndex)
		}
		token.SourceByteStart = originalStart
		token.SourceByteEnd = originalEnd
	}
	pauses := make([]Paul2013InlinePause, len(rawPauses))
	for pauseIndex, pause := range rawPauses {
		afterToken := 0
		for _, token := range tokens {
			if token.SourceByteEnd >= pause.SourceByteOffset {
				break
			}
			afterToken++
		}
		pauses[pauseIndex] = Paul2013InlinePause{
			SourceByteOffset:     pause.SourceByteOffset,
			AfterToken:           afterToken,
			DurationMilliseconds: pause.DurationMilliseconds,
			OutputFramesAt16KHz:  uint64(pause.DurationMilliseconds) * 16,
		}
	}
	return Paul2013PausedLexicalText{Tokens: tokens, Pauses: pauses}, nil
}

type paul2013InlinePauseTag struct {
	SourceByteOffset     int
	DurationMilliseconds uint32
}

func stripPaul2013InlinePauseTags(source string) (string, []int, []paul2013InlinePauseTag, error) {
	destination := make([]byte, len(source)+1)
	destinationCursor := 0
	sourceOffsets := make([]int, 0, len(source))
	var pauses []paul2013InlinePauseTag
	rule, err := Paul2013InlinePauseModelSourceRule(func(event Paul2013InlinePauseTagEvent) error {
		pauses = append(pauses, paul2013InlinePauseTag{
			SourceByteOffset: event.SourceByteOffset, DurationMilliseconds: event.DurationMilliseconds,
		})
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
			for offset := position; offset < len(source); offset++ {
				sourceOffsets = append(sourceOffsets, offset)
			}
			destinationCursor += len(source) - position
			break
		}
		delimiter := position + relativeDelimiter
		if delimiter > position {
			copy(destination[destinationCursor:], sourceBytes[position:delimiter])
			for offset := position; offset < delimiter; offset++ {
				sourceOffsets = append(sourceOffsets, offset)
			}
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
			return "", nil, nil, fmt.Errorf("dispatch markup at source byte %d: %w", position, dispatchErr)
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
		sourceOffsets = append(sourceOffsets, -1)
	}
	return string(destination[:destinationCursor]), sourceOffsets, pauses, nil
}

func parsePaul2013InlinePauseTag(tag string, sourceOffset int) (uint32, error) {
	const prefix = `<vtml_pause time="`
	const suffix = `"/>`
	if !strings.HasPrefix(tag, prefix) || !strings.HasSuffix(tag, suffix) {
		return 0, fmt.Errorf("unsupported markup %q at source byte %d", tag, sourceOffset)
	}
	numeric := strings.TrimSuffix(strings.TrimPrefix(tag, prefix), suffix)
	if numeric == "" {
		return 0, fmt.Errorf("VTML pause at source byte %d has an empty time", sourceOffset)
	}
	for index := range numeric {
		if numeric[index] < '0' || numeric[index] > '9' {
			return 0, fmt.Errorf("VTML pause at source byte %d has a non-decimal time", sourceOffset)
		}
	}
	duration, err := strconv.ParseUint(numeric, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("VTML pause at source byte %d has an unsupported time: %w", sourceOffset, err)
	}
	return uint32(duration), nil
}
