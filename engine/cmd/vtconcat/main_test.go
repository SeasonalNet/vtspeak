package main

import (
	"testing"

	"vtspeak/engine/synthesis"
)

func TestParseUnitRefs(t *testing.T) {
	units, err := parseUnitRefs("gen:12,num:7")
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 2 || units[0].Bank != "gen" || units[0].Index != 12 ||
		units[1].Bank != "num" || units[1].Index != 7 {
		t.Fatalf("unit refs = %+v", units)
	}
}

func TestParseUnitRefsRejectsMalformedInput(t *testing.T) {
	for _, value := range []string{"", "gen", ":2", "gen:", "gen:-1", "gen:4294967296", "gen:1,"} {
		t.Run(value, func(t *testing.T) {
			if _, err := parseUnitRefs(value); err == nil {
				t.Fatalf("parseUnitRefs(%q) accepted malformed input", value)
			}
		})
	}
}

func TestMakeRendererSelectsCursorTimeline(t *testing.T) {
	renderer, err := makeRenderer("cursor", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := renderer.(synthesis.CursorTimelineRenderer); !ok {
		t.Fatalf("cursor renderer type = %T, want synthesis.CursorTimelineRenderer", renderer)
	}
}

func TestMakeRendererSelectsEqualSpanJoin(t *testing.T) {
	renderer, err := makeRenderer("join", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := renderer.(synthesis.UPMTimelineJoinRenderer); !ok {
		t.Fatalf("join renderer type = %T, want synthesis.UPMTimelineJoinRenderer", renderer)
	}
}
