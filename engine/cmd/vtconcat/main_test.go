package main

import "testing"

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
