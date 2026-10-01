package text

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	paul2013VTMLMarkNamedKind     byte = 1
	paul2013VTMLMarkUnnamedKind   byte = 2
	paul2013VTMLMarkTruncatedKind byte = 3
	paul2013VTMLMarkNameCapacity       = 0x200
)

// Paul2013InlineMarkEvent contains the fields established for one inline
// VTML mark record. Output-frame placement is intentionally not included.
type Paul2013InlineMarkEvent struct {
	SourceByteOffset int
	Kind             byte
	Name             string
}

// Paul2013InlineMarkModelSourceRule binds the observed <vtml_mark rule to
// the directly captured unnamed and named self-closing forms. Marks consume
// source bytes without adding lexical text; the callback receives the
// captured record kind and inline name.
func Paul2013InlineMarkModelSourceRule(
	onMark func(Paul2013InlineMarkEvent) error,
) (Paul2013ModelSourceRule, error) {
	for _, descriptor := range Paul2013ModelSourceRuleDescriptors() {
		if descriptor.Key != "<vtml_mark" || descriptor.HandlerAddress != 0x1002efd0 {
			continue
		}
		data := make([]byte, 12)
		putPaul2013ModelSourceUint32(data[0:4], descriptor.RuleDataAddress)
		putPaul2013ModelSourceUint32(data[4:8], descriptor.Flag)
		putPaul2013ModelSourceUint32(data[8:12], descriptor.HandlerDataAddress)
		return Paul2013ModelSourceRule{
			Key:     []byte(descriptor.Key),
			Data:    data,
			Handler: newPaul2013InlineMarkModelSourceHandler(onMark),
		}, nil
	}
	return Paul2013ModelSourceRule{}, errors.New("observed Paul 2013 <vtml_mark> source rule is missing")
}

func newPaul2013InlineMarkModelSourceHandler(
	onMark func(Paul2013InlineMarkEvent) error,
) Paul2013ModelSourceRuleHandler {
	return func(source []byte, sourceCursor *int, _ []byte, _ *int, _ Paul2013ModelSourceRuleContext) (bool, error) {
		if onMark == nil {
			return false, errors.New("inline mark source rule has no event callback")
		}
		start := *sourceCursor
		if start < 0 || start >= len(source) {
			return false, nil
		}
		relativeEnd := bytes.IndexByte(source[start:], '>')
		if relativeEnd < 0 {
			return false, nil
		}
		tagEnd := start + relativeEnd + 1
		event, emit, err := parsePaul2013InlineMarkTag(source[start:tagEnd], start)
		if err != nil {
			return false, nil
		}
		if emit {
			if err := onMark(event); err != nil {
				return false, fmt.Errorf("record inline mark at source byte %d: %w", start, err)
			}
		}
		*sourceCursor = tagEnd
		return true, nil
	}
}

func parsePaul2013InlineMarkTag(tag []byte, sourceOffset int) (Paul2013InlineMarkEvent, bool, error) {
	if !bytes.HasPrefix(tag, []byte("<vtml_mark")) || !bytes.HasSuffix(tag, []byte("/>")) {
		return Paul2013InlineMarkEvent{}, false, fmt.Errorf("unsupported mark tag %q", tag)
	}
	decoder := xml.NewDecoder(strings.NewReader(string(tag)))
	token, err := decoder.Token()
	if err != nil {
		return Paul2013InlineMarkEvent{}, false, fmt.Errorf("parse mark start tag: %w", err)
	}
	start, ok := token.(xml.StartElement)
	if !ok || start.Name.Local != "vtml_mark" || start.Name.Space != "" {
		return Paul2013InlineMarkEvent{}, false, errors.New("expected an unnamespaced vtml_mark element")
	}
	name := ""
	hasName := false
	for _, attribute := range start.Attr {
		if attribute.Name.Space != "" || !strings.EqualFold(attribute.Name.Local, "name") || hasName {
			return Paul2013InlineMarkEvent{}, false, fmt.Errorf("unsupported or duplicate mark attribute %q", attribute.Name.Local)
		}
		hasName = true
		name = attribute.Value
	}
	for {
		var ending xml.Token
		ending, err = decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Paul2013InlineMarkEvent{}, false, fmt.Errorf("finish mark tag: %w", err)
		}
		end, ok := ending.(xml.EndElement)
		if !ok || end.Name != start.Name {
			return Paul2013InlineMarkEvent{}, false, errors.New("mark tag is not empty")
		}
	}
	if !hasName {
		return Paul2013InlineMarkEvent{SourceByteOffset: sourceOffset, Kind: paul2013VTMLMarkUnnamedKind}, true, nil
	}
	if name == "" {
		return Paul2013InlineMarkEvent{}, false, nil
	}
	if !isASCIIBytes([]byte(name)) {
		return Paul2013InlineMarkEvent{}, false, errors.New("non-ASCII VTML mark names are unsupported")
	}
	kind := paul2013VTMLMarkNamedKind
	if len(name) >= paul2013VTMLMarkNameCapacity {
		kind = paul2013VTMLMarkTruncatedKind
		name = name[:paul2013VTMLMarkNameCapacity-1]
	}
	return Paul2013InlineMarkEvent{SourceByteOffset: sourceOffset, Kind: kind, Name: name}, true, nil
}

func isASCIIBytes(value []byte) bool {
	for _, current := range value {
		if current >= 0x80 {
			return false
		}
	}
	return true
}
