package text

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ParsePaul2013VTMLCMUPhoneme parses the directly captured
// <vtml_phoneme alphabet="x-cmu" ph="...">surface</vtml_phoneme> form,
// optionally followed by one sentence-final period. It returns the forced
// phone sequence but does not infer the position-state or unit-selection rows.
func ParsePaul2013VTMLCMUPhoneme(source string) (LexicalPhoneSequence, error) {
	decoder := xml.NewDecoder(strings.NewReader(source))
	token, err := decoder.Token()
	start, ok := token.(xml.StartElement)
	if err != nil || !ok || start.Name.Local != "vtml_phoneme" {
		return LexicalPhoneSequence{}, fmt.Errorf("expected one vtml_phoneme start tag")
	}
	if start.Name.Space != "" {
		return LexicalPhoneSequence{}, fmt.Errorf("namespaced VTML phoneme tags are unsupported")
	}
	attributes := make(map[string]string, len(start.Attr))
	for _, attribute := range start.Attr {
		if attribute.Name.Space != "" {
			return LexicalPhoneSequence{}, fmt.Errorf("namespaced VTML phoneme attributes are unsupported")
		}
		if _, exists := attributes[attribute.Name.Local]; exists {
			return LexicalPhoneSequence{}, fmt.Errorf("duplicate VTML phoneme attribute %q", attribute.Name.Local)
		}
		attributes[attribute.Name.Local] = attribute.Value
	}
	if len(attributes) != 2 || attributes["alphabet"] != "x-cmu" {
		return LexicalPhoneSequence{}, fmt.Errorf("VTML phoneme requires exactly alphabet=x-cmu and ph attributes")
	}
	pronunciation, found := attributes["ph"]
	if !found {
		return LexicalPhoneSequence{}, fmt.Errorf("VTML phoneme has no ph attribute")
	}

	var surface strings.Builder
	contentStart := -1
	contentEnd := -1
	markupEnd := -1
	for {
		token, err = decoder.Token()
		if err != nil {
			return LexicalPhoneSequence{}, fmt.Errorf("read VTML phoneme content: %w", err)
		}
		switch value := token.(type) {
		case xml.CharData:
			if contentStart < 0 {
				contentStart = int(decoder.InputOffset()) - len(value)
			}
			contentEnd = int(decoder.InputOffset())
			surface.Write([]byte(value))
		case xml.EndElement:
			if value.Name != start.Name {
				return LexicalPhoneSequence{}, fmt.Errorf("unexpected VTML phoneme end tag %q", value.Name.Local)
			}
			markupEnd = int(decoder.InputOffset())
			break
		default:
			return LexicalPhoneSequence{}, fmt.Errorf("nested or non-text VTML phoneme content is unsupported")
		}
		if _, ended := token.(xml.EndElement); ended {
			break
		}
	}
	if surface.Len() == 0 || strings.TrimSpace(surface.String()) == "" || contentStart < 0 ||
		contentEnd < contentStart || contentEnd > len(source) || source[contentStart:contentEnd] != surface.String() {
		return LexicalPhoneSequence{}, fmt.Errorf("VTML phoneme surface is empty, escaped, or has invalid source offsets")
	}
	if strings.ContainsAny(surface.String(), "&<>") {
		return LexicalPhoneSequence{}, fmt.Errorf("escaped or markup VTML phoneme surfaces are unsupported")
	}
	for {
		token, err = decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return LexicalPhoneSequence{}, fmt.Errorf("read content after VTML phoneme: %w", err)
		}
		text, ok := token.(xml.CharData)
		if !ok || strings.TrimSpace(string(text)) != "" && strings.TrimSpace(string(text)) != "." {
			return LexicalPhoneSequence{}, fmt.Errorf("content after the VTML phoneme is unsupported")
		}
	}
	trailing := strings.TrimSpace(source[markupEnd:])
	if trailing != "" && trailing != "." {
		return LexicalPhoneSequence{}, fmt.Errorf("content after the VTML phoneme is unsupported")
	}
	sequence, err := BuildPaul2013CMUPhoneSequence(surface.String(), pronunciation)
	if err != nil {
		return LexicalPhoneSequence{}, err
	}
	tokenSpan := &sequence.Tokens[0]
	tokenSpan.SourceByteStart = contentStart
	tokenSpan.SourceByteEnd = contentEnd - 1
	if trailing == "." {
		tokenSpan.SeparatorAfter = "."
	}
	return sequence, nil
}
