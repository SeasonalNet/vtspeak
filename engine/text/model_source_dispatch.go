package text

import "fmt"

// Paul2013ModelSourceRuleHandler runs the selected FUN_1002e990 handler. It
// may advance sourceCursor and destinationCursor and write destination bytes.
// A false handled result restores both cursors, matching the native wrapper;
// bytes written by a handler are not rolled back.
type Paul2013ModelSourceRuleHandler func(
	source []byte,
	sourceCursor *int,
	destination []byte,
	destinationCursor *int,
	context Paul2013ModelSourceRuleContext,
) (handled bool, err error)

// Paul2013ModelSourceRuleContext retains the non-cursor arguments forwarded
// by FUN_1002e990 to the selected native handler.
type Paul2013ModelSourceRuleContext struct {
	Flag        uint32
	RuleData    []byte
	CallerState []byte
}

// Paul2013ModelSourceRule is one nonempty literal key and its native handler.
type Paul2013ModelSourceRule struct {
	Key     []byte
	Data    []byte
	Handler Paul2013ModelSourceRuleHandler
}

// Paul2013ModelSourceDispatchInput supplies the source/destination cursors and
// the handler table consumed by FUN_1002e990. State represents the opaque
// caller object forwarded to a selected handler.
type Paul2013ModelSourceDispatchInput struct {
	Source            []byte
	SourceCursor      int
	Destination       []byte
	DestinationCursor int
	Flag              uint32
	State             []byte
	Rules             []Paul2013ModelSourceRule
}

// Paul2013ModelSourceDispatchResult records the longest selected rule and
// the cursor values after the native handler's accept/miss behavior.
type Paul2013ModelSourceDispatchResult struct {
	SourceCursor      int
	DestinationCursor int
	RuleIndex         int
	Matched           bool
	HandlerAccepted   bool
	Destination       []byte
}

// DispatchPaul2013ModelSourceRule ports FUN_1002e990's longest-key scan,
// mapped-byte comparison, and selected-handler call. Ties retain the first
// rule. The native wrapper does not try a shorter rule when the longest
// handler declines; it restores both cursors and reports a miss.
func DispatchPaul2013ModelSourceRule(
	input Paul2013ModelSourceDispatchInput,
) (Paul2013ModelSourceDispatchResult, error) {
	if input.SourceCursor < 0 || input.SourceCursor > len(input.Source) {
		return Paul2013ModelSourceDispatchResult{}, fmt.Errorf(
			"source cursor %d is outside [0, %d]",
			input.SourceCursor,
			len(input.Source),
		)
	}
	if input.DestinationCursor < 0 || input.DestinationCursor > len(input.Destination) {
		return Paul2013ModelSourceDispatchResult{}, fmt.Errorf(
			"destination cursor %d is outside [0, %d]",
			input.DestinationCursor,
			len(input.Destination),
		)
	}
	result := Paul2013ModelSourceDispatchResult{
		SourceCursor:      input.SourceCursor,
		DestinationCursor: input.DestinationCursor,
		RuleIndex:         -1,
		Destination:       append([]byte(nil), input.Destination...),
	}
	weights := paul2013ContextCharacterWeights()
	longestLength := 0
	for ruleIndex, rule := range input.Rules {
		if len(rule.Key) == 0 {
			return Paul2013ModelSourceDispatchResult{}, fmt.Errorf("source rule %d has an empty key", ruleIndex)
		}
		if containsNUL(rule.Key) {
			return Paul2013ModelSourceDispatchResult{}, fmt.Errorf("source rule %d key contains NUL", ruleIndex)
		}
		if rule.Handler == nil {
			return Paul2013ModelSourceDispatchResult{}, fmt.Errorf("source rule %d has no handler", ruleIndex)
		}
		if len(rule.Key) <= longestLength || len(rule.Key) > len(input.Source)-input.SourceCursor {
			continue
		}
		sourcePrefix := input.Source[input.SourceCursor : input.SourceCursor+len(rule.Key)]
		if containsNUL(sourcePrefix) || ComparePaul2013MappedCString(sourcePrefix, rule.Key, weights) != 0 {
			continue
		}
		longestLength = len(rule.Key)
		result.RuleIndex = ruleIndex
	}
	if result.RuleIndex < 0 {
		return result, nil
	}
	result.Matched = true
	originalSourceCursor := input.SourceCursor
	originalDestinationCursor := input.DestinationCursor
	sourceCursor := originalSourceCursor
	destinationCursor := originalDestinationCursor
	handled, err := input.Rules[result.RuleIndex].Handler(
		input.Source,
		&sourceCursor,
		result.Destination,
		&destinationCursor,
		Paul2013ModelSourceRuleContext{
			Flag:        input.Flag,
			RuleData:    input.Rules[result.RuleIndex].Data,
			CallerState: input.State,
		},
	)
	if err != nil || !handled {
		result.SourceCursor = originalSourceCursor
		result.DestinationCursor = originalDestinationCursor
		if err != nil {
			return result, fmt.Errorf("source rule %d handler: %w", result.RuleIndex, err)
		}
		return result, nil
	}
	if sourceCursor < originalSourceCursor || sourceCursor > len(input.Source) {
		return Paul2013ModelSourceDispatchResult{}, fmt.Errorf(
			"source rule %d returned cursor %d outside [%d, %d]",
			result.RuleIndex,
			sourceCursor,
			originalSourceCursor,
			len(input.Source),
		)
	}
	if destinationCursor < originalDestinationCursor || destinationCursor > len(result.Destination) {
		return Paul2013ModelSourceDispatchResult{}, fmt.Errorf(
			"source rule %d returned destination cursor %d outside [%d, %d]",
			result.RuleIndex,
			destinationCursor,
			originalDestinationCursor,
			len(result.Destination),
		)
	}
	result.SourceCursor = sourceCursor
	result.DestinationCursor = destinationCursor
	result.HandlerAccepted = true
	return result, nil
}
