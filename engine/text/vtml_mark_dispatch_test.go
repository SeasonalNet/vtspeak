package text

import (
	"reflect"
	"strings"
	"testing"
)

func TestPaul2013InlineMarkModelSourceRule(t *testing.T) {
	name512 := strings.Repeat("n", paul2013VTMLMarkNameCapacity)
	cases := []struct {
		name        string
		tag         string
		wantKind    byte
		wantName    string
		wantEmitted bool
	}{
		{name: "unnamed", tag: `<vtml_mark/>`, wantKind: paul2013VTMLMarkUnnamedKind, wantEmitted: true},
		{name: "named", tag: `<vtml_mark name="start"/>`, wantKind: paul2013VTMLMarkNamedKind, wantName: "start", wantEmitted: true},
		{name: "uppercase attribute", tag: `<vtml_mark NAME="Upper"/>`, wantKind: paul2013VTMLMarkNamedKind, wantName: "Upper", wantEmitted: true},
		{name: "space name", tag: `<vtml_mark name=" "/>`, wantKind: paul2013VTMLMarkNamedKind, wantName: " ", wantEmitted: true},
		{name: "empty name", tag: `<vtml_mark name=""/>`},
		{name: "511 bytes", tag: `<vtml_mark name="` + strings.Repeat("x", 511) + `"/>`, wantKind: paul2013VTMLMarkNamedKind, wantName: strings.Repeat("x", 511), wantEmitted: true},
		{name: "512 bytes truncate", tag: `<vtml_mark name="` + name512 + `"/>`, wantKind: paul2013VTMLMarkTruncatedKind, wantName: strings.Repeat("n", 511), wantEmitted: true},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var events []Paul2013InlineMarkEvent
			rule, err := Paul2013InlineMarkModelSourceRule(func(event Paul2013InlineMarkEvent) error {
				events = append(events, event)
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			source := []byte(test.tag + "tail")
			result, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
				Source: source, Destination: []byte("existing"), DestinationCursor: 3,
				Rules: []Paul2013ModelSourceRule{rule},
			})
			if err != nil {
				t.Fatal(err)
			}
			if !result.Matched || !result.HandlerAccepted || result.SourceCursor != len(test.tag) {
				t.Fatalf("mark dispatch result = %+v, want consumed tag and preserved following source", result)
			}
			if result.DestinationCursor != 3 || string(result.Destination) != "existing" {
				t.Fatalf("mark dispatch changed lexical destination: cursor=%d data=%q", result.DestinationCursor, result.Destination)
			}
			if len(events) != boolInt(test.wantEmitted) {
				t.Fatalf("mark events = %+v, want emitted=%t", events, test.wantEmitted)
			}
			if test.wantEmitted {
				want := Paul2013InlineMarkEvent{SourceByteOffset: 0, Kind: test.wantKind, Name: test.wantName}
				if !reflect.DeepEqual(events[0], want) {
					t.Fatalf("mark event = %+v, want %+v", events[0], want)
				}
			}
		})
	}
}

func TestPaul2013InlineMarkModelSourceRuleDeclinesUnsupportedForms(t *testing.T) {
	for _, source := range []string{
		`<vtml_mark name="one" extra="two"/>`,
		`<vtml_mark name="é"/>`,
		`<vtml_mark name="nested">x</vtml_mark>`,
	} {
		rule, err := Paul2013InlineMarkModelSourceRule(func(Paul2013InlineMarkEvent) error { return nil })
		if err != nil {
			t.Fatal(err)
		}
		result, err := DispatchPaul2013ModelSourceRule(Paul2013ModelSourceDispatchInput{
			Source: []byte(source), Destination: make([]byte, 1), Rules: []Paul2013ModelSourceRule{rule},
		})
		if err != nil {
			t.Fatal(err)
		}
		if result.HandlerAccepted || result.SourceCursor != 0 || result.DestinationCursor != 0 {
			t.Errorf("unsupported mark %q dispatch = %+v, want declined with restored cursors", source, result)
		}
	}
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
