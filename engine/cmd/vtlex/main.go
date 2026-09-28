package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"vtspeak/engine/text"
)

func main() {
	dictionaryRoot := flag.String("dictionary-root", "", "local shared dict-eng resource directory")
	source := flag.String("text", "", "ASCII text to normalize and resolve")
	choices := flag.String("choices", "", "comma-separated zero-based pronunciation indexes, one per normalized token")
	flag.Parse()
	if err := run(*dictionaryRoot, *source, *choices); err != nil {
		fmt.Fprintln(os.Stderr, "vtlex:", err)
		os.Exit(1)
	}
}

func run(dictionaryRoot, source, choiceSpec string) error {
	if dictionaryRoot == "" {
		return errors.New("-dictionary-root is required")
	}
	if source == "" {
		return errors.New("-text is required")
	}
	dictionary, err := text.LoadEmbeddedDictionary(dictionaryRoot)
	if err != nil {
		return fmt.Errorf("load embedded dictionary: %w", err)
	}
	tokens, err := (text.LexiconFrontend{Dictionary: dictionary}).ResolveText(context.Background(), source)
	if err != nil {
		return fmt.Errorf("resolve text: %w", err)
	}
	choices, err := parseAlternativeIndexes(choiceSpec)
	if err != nil {
		return err
	}
	var output strings.Builder
	for _, token := range tokens {
		fmt.Fprintf(&output, "token %q surface=%q before=%q after=%q\n",
			token.SourceSurface, token.Surface, token.SeparatorBefore, token.SeparatorAfter)
		for alternativeIndex, alternative := range token.Alternatives {
			fmt.Fprintf(&output, "  alt[%d] path=%s symbols=%s phones=%s\n",
				alternativeIndex,
				hex.EncodeToString(alternative.Path),
				hex.EncodeToString(alternative.Symbols),
				formatPhones(alternative.Phones),
			)
		}
	}
	if choiceSpec != "" {
		sequence, err := text.SelectLexicalPronunciations(tokens, choices)
		if err != nil {
			return fmt.Errorf("assemble selected pronunciation sequence: %w", err)
		}
		fmt.Fprintf(&output, "selected symbols=%s phones=%s\n",
			hex.EncodeToString(sequence.Symbols), formatPhones(sequence.Phones))
		for tokenIndex, span := range sequence.Tokens {
			fmt.Fprintf(&output, "  selected-token[%d] alt=%d phones=[%d,%d)\n",
				tokenIndex, span.AlternativeIndex, span.PhoneStart, span.PhoneEnd)
		}
	}
	if _, err := os.Stdout.WriteString(output.String()); err != nil {
		return fmt.Errorf("write pronunciation output: %w", err)
	}
	return nil
}

func parseAlternativeIndexes(value string) ([]int, error) {
	if value == "" {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	indexes := make([]int, len(parts))
	for index, part := range parts {
		choice, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || choice < 0 {
			return nil, fmt.Errorf("pronunciation choice %d %q must be a nonnegative integer", index, part)
		}
		indexes[index] = choice
	}
	return indexes, nil
}

func formatPhones(phones []text.CMUPhone) string {
	labels := make([]string, len(phones))
	for index, phone := range phones {
		labels[index] = phone.Label
		if phone.Vowel {
			labels[index] += fmt.Sprint(phone.Stress)
		}
	}
	return strings.Join(labels, " ")
}
