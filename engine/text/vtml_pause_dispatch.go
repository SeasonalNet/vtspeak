package text

import (
	"bytes"
	"errors"
	"fmt"
)

// Paul2013InlinePauseTagEvent contains the timing fields parsed by the
// model-source rule. Token-boundary placement must be derived by the caller
// from the source byte offset.
type Paul2013InlinePauseTagEvent struct {
	SourceByteOffset     int
	DurationMilliseconds uint32
	OutputFramesAt16KHz  uint64
}

// Paul2013InlinePauseModelSourceRule returns the observed `<vtml_pause` rule
// bound to the exact supported `<vtml_pause time="N"/>` form. It writes a
// separating space and reports timing through onPause. Other attributes, tag
// forms, and the separate <vt_pause> handler remain unsupported.
func Paul2013InlinePauseModelSourceRule(
	onPause func(Paul2013InlinePauseTagEvent) error,
) (Paul2013ModelSourceRule, error) {
	for _, descriptor := range Paul2013ModelSourceRuleDescriptors() {
		if descriptor.Key != "<vtml_pause" || descriptor.HandlerAddress != 0x1002fb10 {
			continue
		}
		data := make([]byte, 12)
		putPaul2013ModelSourceUint32(data[0:4], descriptor.RuleDataAddress)
		putPaul2013ModelSourceUint32(data[4:8], descriptor.Flag)
		putPaul2013ModelSourceUint32(data[8:12], descriptor.HandlerDataAddress)
		return Paul2013ModelSourceRule{
			Key:     []byte(descriptor.Key),
			Data:    data,
			Handler: newPaul2013InlinePauseModelSourceHandler(onPause),
		}, nil
	}
	return Paul2013ModelSourceRule{}, errors.New("observed Paul 2013 <vtml_pause> source rule is missing")
}

func newPaul2013InlinePauseModelSourceHandler(
	onPause func(Paul2013InlinePauseTagEvent) error,
) Paul2013ModelSourceRuleHandler {
	return func(source []byte, sourceCursor *int, destination []byte, destinationCursor *int, _ Paul2013ModelSourceRuleContext) (bool, error) {
		if onPause == nil {
			return false, errors.New("inline pause source rule has no event callback")
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
		duration, err := parsePaul2013InlinePauseTag(string(source[start:tagEnd]), start)
		if err != nil {
			return false, nil
		}
		if *destinationCursor < 0 || *destinationCursor >= len(destination) {
			return false, fmt.Errorf("inline pause destination cursor %d has no room for a separator", *destinationCursor)
		}
		event := Paul2013InlinePauseTagEvent{
			SourceByteOffset:     start,
			DurationMilliseconds: duration,
			OutputFramesAt16KHz:  uint64(duration) * 16,
		}
		if err := onPause(event); err != nil {
			return false, fmt.Errorf("record inline pause at source byte %d: %w", start, err)
		}
		destination[*destinationCursor] = ' '
		*destinationCursor = *destinationCursor + 1
		*sourceCursor = tagEnd
		return true, nil
	}
}
