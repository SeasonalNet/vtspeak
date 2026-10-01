package text

import (
	"fmt"
	"strings"
)

// Paul2013ModelSourceTransformInput contains the inputs used by the supported
// branches of FUN_1000d640. FinalRow, ParserType, AssociatedText, and
// PreviousCount correspond to native arguments 3, 4, 5, and 6; their
// higher-level meanings beyond the final-row gate are not inferred here.
type Paul2013ModelSourceTransformInput struct {
	Source         []byte
	ParserType     int8
	AssociatedText []byte
	PreviousCount  int16
	FinalRow       bool
}

// Paul2013ModelSourceTransformResult preserves both mutated native strings:
// Output corresponds to param_2 and Remainder to the post-call param_1.
type Paul2013ModelSourceTransformResult struct {
	Output            []byte
	Remainder         []byte
	SplitTrailingTwo  bool
	SplitAtTerminal   bool
	SplitAtNativePair bool
}

// TransformPaul2013ModelSourceString ports the directly recovered copy and
// split paths in FUN_1000d640. For parser types other than -1, it moves the
// final two bytes from Source to Remainder when Source has at least three bytes
// and its penultimate byte has no 0xc0 character attribute. The parser-type -1
// scanner trims leading spaces, copies source bytes, and handles contractions,
// terminal punctuation, hyphen splits, and the directly compared 0xa2 0xfe
// pair. A non-final terminal punctuation path that needs FUN_10003a70 fails
// closed when no lookup callback is available.
func TransformPaul2013ModelSourceString(
	input Paul2013ModelSourceTransformInput,
	attributes [256]byte,
) (Paul2013ModelSourceTransformResult, error) {
	return TransformPaul2013ModelSourceStringWithEmbeddedLookup(input, attributes, nil)
}

// TransformPaul2013ModelSourceStringWithEmbeddedLookup also ports the scanner's
// terminal punctuation split. lookup reports whether FUN_10003a70 finds the
// accumulated output token; it is called only for a non-final scanner row when
// the preceding copied byte was not a period. A nil lookup keeps this branch
// fail-closed.
func TransformPaul2013ModelSourceStringWithEmbeddedLookup(
	input Paul2013ModelSourceTransformInput,
	attributes [256]byte,
	lookup func([]byte) (bool, error),
) (Paul2013ModelSourceTransformResult, error) {
	source := cString(input.Source)
	associated := cString(input.AssociatedText)
	if input.ParserType == -1 && len(associated) == 0 && input.PreviousCount == 0 {
		return transformPaul2013ScannerSource(source, input.FinalRow, attributes, lookup)
	}
	if input.ParserType != -1 && len(source) > 2 && attributes[source[len(source)-2]]&0xc0 == 0 {
		return Paul2013ModelSourceTransformResult{
			Output:           append([]byte(nil), source[:len(source)-2]...),
			Remainder:        append([]byte(nil), source[len(source)-2:]...),
			SplitTrailingTwo: true,
		}, nil
	}
	return Paul2013ModelSourceTransformResult{
		Output:    append([]byte(nil), source...),
		Remainder: []byte{},
	}, nil
}

func transformPaul2013ScannerSource(
	source []byte,
	finalRow bool,
	attributes [256]byte,
	lookup func([]byte) (bool, error),
) (Paul2013ModelSourceTransformResult, error) {
	start := 0
	for start < len(source) && source[start] == ' ' {
		start++
	}
	source = source[start:]
	output := make([]byte, 0, len(source))
	precedingPeriod := false
	for offset, character := range source {
		if len(source) != 2 && len(output) != 0 && character == 0xa2 &&
			offset+1 < len(source) && source[offset+1] == 0xfe {
			if len(output) > 31 {
				return Paul2013ModelSourceTransformResult{}, fmt.Errorf("contraction scanner output has %d bytes, exceeds its 32-byte caller buffer", len(output))
			}
			return Paul2013ModelSourceTransformResult{
				Output:            append([]byte(nil), output...),
				Remainder:         append([]byte(nil), source[offset:]...),
				SplitAtNativePair: true,
			}, nil
		}
		if offset == len(source)-1 && strings.ContainsRune(";,.?!", rune(character)) {
			if finalRow {
				if len(output) > 31 {
					return Paul2013ModelSourceTransformResult{}, fmt.Errorf("contraction scanner output has %d bytes, exceeds its 32-byte caller buffer", len(output))
				}
				return Paul2013ModelSourceTransformResult{
					Output:          append([]byte(nil), output...),
					Remainder:       append([]byte(nil), source[offset:]...),
					SplitAtTerminal: true,
				}, nil
			}
			keepPunctuation := precedingPeriod
			if !keepPunctuation {
				if lookup == nil {
					return Paul2013ModelSourceTransformResult{}, fmt.Errorf("terminal punctuation split requires an embedded lookup callback")
				}
				found, err := lookup(append([]byte(nil), output...))
				if err != nil {
					return Paul2013ModelSourceTransformResult{}, fmt.Errorf("lookup terminal punctuation token: %w", err)
				}
				keepPunctuation = found
			}
			if keepPunctuation {
				if len(output)+1 > 31 {
					return Paul2013ModelSourceTransformResult{}, fmt.Errorf("contraction scanner output has %d bytes, exceeds its 32-byte caller buffer", len(output)+1)
				}
				output = append(output, character)
				return Paul2013ModelSourceTransformResult{
					Output:          output,
					Remainder:       []byte{},
					SplitAtTerminal: true,
				}, nil
			}
			if len(output) > 31 {
				return Paul2013ModelSourceTransformResult{}, fmt.Errorf("contraction scanner output has %d bytes, exceeds its 32-byte caller buffer", len(output))
			}
			return Paul2013ModelSourceTransformResult{
				Output:          append([]byte(nil), output...),
				Remainder:       append([]byte(nil), source[offset:]...),
				SplitAtTerminal: true,
			}, nil
		}
		if offset == len(source)-1 && character == '-' {
			if len(output) > 31 {
				return Paul2013ModelSourceTransformResult{}, fmt.Errorf("contraction scanner output has %d bytes, exceeds its 32-byte caller buffer", len(output))
			}
			return Paul2013ModelSourceTransformResult{
				Output:          append([]byte(nil), output...),
				Remainder:       append([]byte(nil), source[offset:]...),
				SplitAtTerminal: true,
			}, nil
		}
		if offset > 0 && offset+1 < len(source) && (character == '.' || character == ',') &&
			attributes[source[offset-1]]&0x10 != 0 && attributes[source[offset+1]]&0x10 != 0 {
			output = append(output, character)
			precedingPeriod = false
			continue
		}
		output = append(output, character)
		precedingPeriod = character == '.'
	}
	if len(output) > 31 {
		return Paul2013ModelSourceTransformResult{}, fmt.Errorf("contraction scanner output has %d bytes, exceeds its 32-byte caller buffer", len(output))
	}
	return Paul2013ModelSourceTransformResult{
		Output:    output,
		Remainder: []byte{},
	}, nil
}

// TransformPaul2013ASCIIModelSourceString is the strict ASCII convenience
// wrapper. Use the character-table entry point for native byte strings.
func TransformPaul2013ASCIIModelSourceString(
	input Paul2013ModelSourceTransformInput,
) (Paul2013ModelSourceTransformResult, error) {
	return TransformPaul2013ASCIIModelSourceStringWithEmbeddedLookup(input, nil)
}

// TransformPaul2013ASCIIModelSourceStringWithEmbeddedLookup supplies the
// recovered ASCII character attributes and an optional embedded-token lookup
// for the native terminal-punctuation decision.
func TransformPaul2013ASCIIModelSourceStringWithEmbeddedLookup(
	input Paul2013ModelSourceTransformInput,
	lookup func([]byte) (bool, error),
) (Paul2013ModelSourceTransformResult, error) {
	source := cString(input.Source)
	associated := cString(input.AssociatedText)
	if input.ParserType == -1 && len(associated) == 0 && input.PreviousCount == 0 {
		for offset, character := range source {
			if character >= 0x80 {
				return Paul2013ModelSourceTransformResult{}, fmt.Errorf("ASCII model source transform rejects high-bit byte 0x%02x at offset %d", character, offset)
			}
		}
		return transformPaul2013ScannerSource(source, input.FinalRow, Paul2013ScannerCharacterAttributes(), lookup)
	}
	if input.ParserType != -1 && len(source) > 2 && source[len(source)-2] >= 0x80 {
		return Paul2013ModelSourceTransformResult{}, fmt.Errorf("ASCII model source transform rejects high-bit penultimate byte 0x%02x at offset %d", source[len(source)-2], len(source)-2)
	}
	return TransformPaul2013ModelSourceStringWithEmbeddedLookup(
		input, Paul2013ExceptionCharacterAttributes(), lookup,
	)
}

// TransformPaul2013ModelSourceStringFromPaul2013CharacterTable applies the
// native signed-char lookup values. Scanner rows preserve the directly
// recovered high-byte copy and 0xa2 0xfe split behavior; non-scanner rows use
// the exact signed-index split predicate.
func TransformPaul2013ModelSourceStringFromPaul2013CharacterTable(
	input Paul2013ModelSourceTransformInput,
	lookup func([]byte) (bool, error),
) (Paul2013ModelSourceTransformResult, error) {
	return TransformPaul2013ModelSourceStringWithEmbeddedLookup(
		input, Paul2013ExceptionCharacterAttributes(), lookup,
	)
}

// Paul2013ScannerCharacterAttributes exposes the recovered ASCII table used
// by FUN_1000d640, including bit 0x10 on digits for decimal/comma grouping.
func Paul2013ScannerCharacterAttributes() [256]byte {
	return Paul2013ExceptionCharacterAttributes()
}

// TransformAndClassifyPaul2013ModelSource composes the supported
// FUN_1000d640 branches with FUN_1000fd20, applying the classification of the
// transformed Output to a copied dictionary phone-row projection. It returns
// the transformed strings so callers can continue the subsequent native
// dictionary-key conversion independently.
func (rows Paul2013DictionaryPhoneRows) TransformAndClassifyPaul2013ModelSource(
	input Paul2013ModelSourceTransformInput,
) (Paul2013DictionaryPhoneRows, Paul2013ModelSourceTransformResult, error) {
	return rows.TransformAndClassifyPaul2013ModelSourceWithEmbeddedLookup(input, nil)
}

// TransformAndClassifyPaul2013ModelSourceWithEmbeddedLookup composes the
// source transform, optional native embedded-token decision, and source-type
// classifier for one dictionary phone-row projection.
func (rows Paul2013DictionaryPhoneRows) TransformAndClassifyPaul2013ModelSourceWithEmbeddedLookup(
	input Paul2013ModelSourceTransformInput,
	lookup func([]byte) (bool, error),
) (Paul2013DictionaryPhoneRows, Paul2013ModelSourceTransformResult, error) {
	transformed, err := TransformPaul2013ModelSourceStringFromPaul2013CharacterTable(input, lookup)
	if err != nil {
		return Paul2013DictionaryPhoneRows{}, Paul2013ModelSourceTransformResult{}, err
	}
	resultType := ClassifyPaul2013ModelSourceType(transformed.Output, Paul2013ExceptionCharacterAttributes())
	rows.ResultType = resultType
	return rows, transformed, nil
}
