package voice

import "testing"

func TestGlobalUnitOrdinalUsesDblistBankOrder(t *testing.T) {
	counts := map[string]uint32{"gen": 3, "num": 2, "etc": 4, "alp": 1}
	tests := []struct {
		bank  string
		index uint32
		want  uint32
	}{
		{bank: "gen", index: 0, want: 0},
		{bank: "gen", index: 2, want: 2},
		{bank: "num", index: 0, want: 3},
		{bank: "etc", index: 0, want: 5},
		{bank: "etc", index: 3, want: 8},
		{bank: "alp", index: 0, want: 9},
	}
	for _, test := range tests {
		got, err := globalUnitOrdinal(paul2013BankNames[:], counts, test.bank, test.index)
		if err != nil {
			t.Fatalf("ordinal for %s:%d: %v", test.bank, test.index, err)
		}
		if got != test.want {
			t.Errorf("ordinal for %s:%d = %d, want %d", test.bank, test.index, got, test.want)
		}
	}
}

func TestPaul2013GlobalUnitOrdinalUsesLoadedBankCounts(t *testing.T) {
	model := &Paul2013{Banks: map[string]*Bank{
		"gen": {unitCount: 3}, "num": {unitCount: 2},
		"etc": {unitCount: 4}, "alp": {unitCount: 1},
	}}
	got, err := model.GlobalUnitOrdinal("etc", 2)
	if err != nil {
		t.Fatal(err)
	}
	if got != 7 {
		t.Fatalf("global ordinal for etc:2 = %d, want 7", got)
	}
	bank, index, err := model.UnitAtGlobalOrdinal(got)
	if err != nil {
		t.Fatal(err)
	}
	if bank != "etc" || index != 2 {
		t.Fatalf("global ordinal %d resolved to %s:%d, want etc:2", got, bank, index)
	}
}

func TestGlobalUnitOrdinalRejectsMissingAndOutOfRangeBanks(t *testing.T) {
	counts := map[string]uint32{"gen": 3, "num": 2, "etc": 4, "alp": 1}
	for _, test := range []struct {
		name  string
		bank  string
		index uint32
	}{
		{name: "unknown bank", bank: "other", index: 0},
		{name: "out of range unit", bank: "num", index: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := globalUnitOrdinal(paul2013BankNames[:], counts, test.bank, test.index); err == nil {
				t.Fatalf("accepted %s:%d", test.bank, test.index)
			}
		})
	}
	delete(counts, "etc")
	if _, err := globalUnitOrdinal(paul2013BankNames[:], counts, "alp", 0); err == nil {
		t.Fatal("accepted missing preceding bank count")
	}
}
