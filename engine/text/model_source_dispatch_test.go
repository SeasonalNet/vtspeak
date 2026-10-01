package text

import (
	"errors"
	"testing"
)

func TestDispatchPaul2013ModelSourceRuleChoosesLongestMappedKey(t *testing.T) {
	shortCalled, longCalled := false, false
	contextReceived := false
	originalDestination := make([]byte, 8)
	got, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source:      []byte("ALPHABET!"),
		Destination: originalDestination,
		Flag:        7,
		State:       []byte{0xaa},
		Rules: []Paul2013ModelSourceRule{
			{Key: []byte("ALPHA"), Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) {
				shortCalled = true
				return true, nil
			}},
			{Key: []byte("alphabet"), Data: []byte{0xbb}, Handler: func(source []byte, sourceCursor *int, destination []byte, destinationCursor *int, context Paul2013ModelSourceRuleContext) (bool, error) {
				longCalled = true
				contextReceived = context.Flag == 7 && len(context.RuleData) == 1 && context.RuleData[0] == 0xbb &&
					len(context.CallerState) == 1 && context.CallerState[0] == 0xaa
				copy(destination[*destinationCursor:], source[*sourceCursor:*sourceCursor+8])
				*sourceCursor += 8
				*destinationCursor += 8
				return true, nil
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Matched || !got.HandlerAccepted || got.RuleIndex != 1 || got.SourceCursor != 8 || got.DestinationCursor != 8 {
		t.Fatalf("dispatch result = %+v", got)
	}
	if shortCalled || !longCalled || !contextReceived || string(got.Destination) != "ALPHABET" {
		t.Fatalf("handler selection/context: short=%t long=%t context=%t destination=%q", shortCalled, longCalled, contextReceived, got.Destination)
	}
	if string(originalDestination) != "\x00\x00\x00\x00\x00\x00\x00\x00" {
		t.Fatalf("dispatcher mutated caller destination: %q", originalDestination)
	}
}

func TestDispatchPaul2013ModelSourceRuleTieUsesFirstAndMissRestoresCursors(t *testing.T) {
	firstCalled, secondCalled := false, false
	got, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source:            []byte("word!"),
		SourceCursor:      0,
		Destination:       []byte("....."),
		DestinationCursor: 2,
		Rules: []Paul2013ModelSourceRule{
			{Key: []byte("WORD"), Handler: func(_ []byte, sourceCursor *int, destination []byte, destinationCursor *int, _ Paul2013ModelSourceRuleContext) (bool, error) {
				firstCalled = true
				*sourceCursor = 4
				destination[*destinationCursor] = 'x'
				*destinationCursor++
				return false, nil
			}},
			{Key: []byte("word"), Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) {
				secondCalled = true
				return true, nil
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Matched || got.HandlerAccepted || got.RuleIndex != 0 || got.SourceCursor != 0 || got.DestinationCursor != 2 {
		t.Fatalf("declined longest handler result = %+v", got)
	}
	if !firstCalled || secondCalled || got.Destination[2] != 'x' {
		t.Fatalf("handler call or retained write: first=%t second=%t destination=%q", firstCalled, secondCalled, got.Destination)
	}
}

func TestDispatchPaul2013ModelSourceRuleSkipsNULTerminatedPrefixes(t *testing.T) {
	called := false
	got, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source:      []byte{'A', 0, 'B'},
		Destination: make([]byte, 2),
		Rules: []Paul2013ModelSourceRule{
			{Key: []byte("AB"), Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) {
				return true, nil
			}},
			{Key: []byte("A"), Handler: func(_ []byte, sourceCursor *int, _ []byte, _ *int, _ Paul2013ModelSourceRuleContext) (bool, error) {
				called = true
				*sourceCursor++
				return true, nil
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called || got.RuleIndex != 1 || got.SourceCursor != 1 || !got.HandlerAccepted {
		t.Fatalf("NUL-bounded dispatch result = %+v, handler called=%t", got, called)
	}
}

func TestDispatchPaul2013ModelSourceRuleNoMatchAndErrors(t *testing.T) {
	got, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source:            []byte("cat"),
		SourceCursor:      1,
		Destination:       []byte("out"),
		DestinationCursor: 1,
		Rules: []Paul2013ModelSourceRule{{
			Key: []byte("dog"), Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) { return true, nil },
		}},
	})
	if err != nil || got.Matched || got.SourceCursor != 1 || got.DestinationCursor != 1 {
		t.Fatalf("no-match result = %+v, err=%v", got, err)
	}

	handlerErr := errors.New("handler error")
	_, err = DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
		Source: []byte("cat"),
		Rules: []Paul2013ModelSourceRule{{
			Key: []byte("cat"), Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) {
				return false, handlerErr
			},
		}},
	})
	if !errors.Is(err, handlerErr) {
		t.Fatalf("handler error = %v, want wrapped handler error", err)
	}

	for name, input := range map[string]Paul2013ModelSourceDispatchInput{
		"source cursor":      {Source: []byte("x"), SourceCursor: 2},
		"destination cursor": {Destination: []byte("x"), DestinationCursor: 2},
		"empty key":          {Rules: []Paul2013ModelSourceRule{{Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) { return true, nil }}}},
		"missing handler":    {Rules: []Paul2013ModelSourceRule{{Key: []byte("x")}}},
		"NUL key":            {Rules: []Paul2013ModelSourceRule{{Key: []byte{'x', 0}, Handler: func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) { return true, nil }}}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := DispatchPaul2013ModelSourceRule(input); err == nil {
				t.Fatal("malformed dispatcher input was accepted")
			}
		})
	}
	for name, rule := range map[string]Paul2013ModelSourceRule{
		"source cursor": {Key: []byte("cat"), Handler: func(_ []byte, sourceCursor *int, _ []byte, _ *int, _ Paul2013ModelSourceRuleContext) (bool, error) {
			*sourceCursor = 4
			return true, nil
		}},
		"destination cursor": {Key: []byte("cat"), Handler: func(_ []byte, _ *int, _ []byte, destinationCursor *int, _ Paul2013ModelSourceRuleContext) (bool, error) {
			*destinationCursor = 1
			return true, nil
		}},
	} {
		t.Run("handler "+name, func(t *testing.T) {
			_, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
				Source: []byte("cat"), Rules: []Paul2013ModelSourceRule{rule},
			})
			if err == nil {
				t.Fatal("handler cursor outside its buffer was accepted")
			}
		})
	}
}
