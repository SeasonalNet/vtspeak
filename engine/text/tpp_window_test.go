package text

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestMatchPaul2013TPPWindowSkipsAndRetainsLastMatch(t *testing.T) {
	rows := []Paul2013TPPWindowRow{
		{Word0: 10, Word1: 10, Word2: 1, State: 'S', Text: []byte("START")},
		{Word0: 11, Word1: 10, Word2: 1, State: 'S', Text: []byte("DASH")},
		{Word0: 11, Word1: 10, Word2: 1, State: 'A', Text: []byte("FIRST")},
		{Word0: 12, Word1: 11, Word2: 1, State: 'A', Text: []byte("SECOND")},
		{Word0: 13, Word1: 12, Word2: 1, State: 'B', Text: []byte("STOP")},
		{Word0: 14, Word1: 13, Word2: 1, State: 'A', Text: []byte("AFTER")},
	}
	var lookupCalls int
	got, err := MatchPaul2013TPPWindow(Paul2013TPPWindowInput{
		Rows: rows, StartRow: 0, CharacterMap: Paul2013EmbeddedKeyTables().CharacterMap,
		DictionaryIndex: 4, InitialOutputRowIndex: 77,
		LookupCompound: func(compound []byte, dictionaryIndex int) ([]byte, bool, error) {
			lookupCalls++
			if dictionaryIndex != 4 {
				t.Fatalf("dictionary index = %d, want 4", dictionaryIndex)
			}
			want := [][]byte{[]byte("START-FIRST"), []byte("START-FIRST-SECOND")}
			if !reflect.DeepEqual(compound, want[lookupCalls-1]) {
				t.Fatalf("lookup compound %d = %q, want %q", lookupCalls, compound, want[lookupCalls-1])
			}
			return []byte{byte('0' + lookupCalls)}, true, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if lookupCalls != 2 {
		t.Fatalf("lookup calls = %d, want 2", lookupCalls)
	}
	if !got.Matched || got.OutputRowIndex != 3 || string(got.Compound) != "START-FIRST-SECOND" || string(got.Replacement) != "2" {
		t.Fatalf("window result = %+v, want last successful match at row 3", got)
	}
}

func TestMatchPaul2013TPPWindowNonDashStateRowBreaks(t *testing.T) {
	rows := []Paul2013TPPWindowRow{
		{Word0: 10, Word1: 10, Word2: 1, State: 'S', Text: []byte("START")},
		{Word0: 11, Word1: 10, Word2: 1, State: 'S', Text: []byte("SPECIAL")},
	}
	lookedUp := false
	got, err := MatchPaul2013TPPWindow(Paul2013TPPWindowInput{
		Rows: rows, CharacterMap: Paul2013EmbeddedKeyTables().CharacterMap,
		LookupCompound: func([]byte, int) ([]byte, bool, error) { lookedUp = true; return nil, false, nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Matched || lookedUp {
		t.Fatalf("non-dash S row must break after failing fast-skip: result=%+v lookedUp=%t", got, lookedUp)
	}
}

func TestMatchPaul2013TPPWindowHonorsCumulativeCutoffAndErrors(t *testing.T) {
	rows := []Paul2013TPPWindowRow{
		{Word2: 1, State: 'S', Text: []byte("START")},
		{Word2: 28, State: 'A', Text: []byte("TOO-LATE")},
	}
	called := false
	input := Paul2013TPPWindowInput{
		Rows:           rows,
		LookupCompound: func([]byte, int) ([]byte, bool, error) { called = true; return nil, false, nil },
	}
	if _, err := MatchPaul2013TPPWindow(input); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("lookup ran after the native cumulative cutoff")
	}

	input.Rows = input.Rows[:1]
	input.LookupCompound = func([]byte, int) ([]byte, bool, error) { return nil, false, errors.New("lookup failed") }
	// A second row is required to reach the lookup callback.
	input.Rows = append(input.Rows, Paul2013TPPWindowRow{Word0: 1, Word1: 0, State: 'A', Text: []byte("NEXT")})
	if _, err := MatchPaul2013TPPWindow(input); err == nil {
		t.Fatal("lookup callback error was lost")
	}
}

func TestMatchPaul2013TPPWindowProcessesAtMostFiveEligibleRows(t *testing.T) {
	rows := make([]Paul2013TPPWindowRow, 8)
	rows[0] = Paul2013TPPWindowRow{State: 'S', Text: []byte("START")}
	for index := 1; index < len(rows); index++ {
		rows[index] = Paul2013TPPWindowRow{State: 'A', Text: []byte("NEXT")}
	}
	lookups := 0
	_, err := MatchPaul2013TPPWindow(Paul2013TPPWindowInput{
		Rows: rows,
		LookupCompound: func([]byte, int) ([]byte, bool, error) {
			lookups++
			return nil, false, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if lookups != 5 {
		t.Fatalf("eligible-row lookups = %d, want native maximum 5", lookups)
	}
}

func TestMatchPaul2013TPPWindowRejectsUnsafeNativeBuffer(t *testing.T) {
	input := Paul2013TPPWindowInput{
		Rows:           []Paul2013TPPWindowRow{{Text: bytes.Repeat([]byte{'X'}, 32)}},
		LookupCompound: func([]byte, int) ([]byte, bool, error) { return nil, false, nil },
	}
	if _, err := MatchPaul2013TPPWindow(input); err == nil {
		t.Fatal("oversized initial text accepted")
	}
}
