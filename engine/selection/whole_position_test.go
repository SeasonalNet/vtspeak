package selection

import (
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

func TestEvaluatePaul2013WholePositionDeduplicatesAndAppliesThreshold(t *testing.T) {
	leftMembers := make([]voice.UnitLocation, 6)
	for index := range leftMembers {
		leftMembers[index] = voice.UnitLocation{Bank: "gen", Index: uint32(index)}
	}
	rightMembers := make([]voice.UnitLocation, 3)
	for index := range rightMembers {
		rightMembers[index] = voice.UnitLocation{Bank: "num", Index: uint32(index)}
	}
	classes := []voice.ClassRecord{
		{ID: 2, Key: [5]byte{1}, Members: leftMembers},
		{ID: 1, Key: [5]byte{2}, Members: rightMembers},
		{ID: 2, Key: [5]byte{1}, Members: leftMembers},
	}
	result, err := EvaluatePaul2013WholePosition([5]byte{1}, classes, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Classes) != 2 || result.MetricSubtotal != 9 || result.Accepted {
		t.Fatalf("result = %#v, want two deduplicated classes, subtotal 9, rejected", result)
	}
	if result.Classes[0].ID != 1 || result.Classes[1].ID != 2 {
		t.Fatalf("ranked class order = [%d, %d], want ascending native class IDs", result.Classes[0].ID, result.Classes[1].ID)
	}

	classes[1].Members = append(classes[1].Members, voice.UnitLocation{Bank: "num", Index: 3})
	result, err = EvaluatePaul2013WholePosition([5]byte{1}, classes, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.MetricSubtotal != 10 || !result.Accepted {
		t.Fatalf("subtotal 10 result = %#v, want accepted", result)
	}
}

func TestEvaluatePaul2013WholePositionClassEightAcceptsAnyPositiveSubtotal(t *testing.T) {
	classes := []voice.ClassRecord{{ID: 1, Key: [5]byte{1}, Members: []voice.UnitLocation{{Bank: "gen", Index: 0}}}}
	result, err := EvaluatePaul2013WholePosition([5]byte{1}, classes, 8)
	if err != nil {
		t.Fatal(err)
	}
	if result.MetricSubtotal != 1 || !result.Accepted {
		t.Fatalf("result = %#v, want positive subtotal accepted for row class 8", result)
	}
	classes[0].Members = nil
	result, err = EvaluatePaul2013WholePosition([5]byte{1}, classes, 8)
	if err != nil {
		t.Fatal(err)
	}
	if result.Accepted {
		t.Fatal("row class 8 accepted a zero subtotal")
	}
}

func TestEvaluatePaul2013WholePositionRejectsConflictingDuplicate(t *testing.T) {
	base := voice.ClassRecord{ID: 1, Key: [5]byte{1}, Members: []voice.UnitLocation{{Bank: "gen", Index: 0}}}
	conflict := voice.ClassRecord{ID: 1, Key: [5]byte{2}, Members: base.Members}
	if _, err := EvaluatePaul2013WholePosition([5]byte{1}, []voice.ClassRecord{base, conflict}, 0); err == nil {
		t.Fatal("conflicting duplicate class was accepted")
	}
}

type wholePositionClassMap map[[5]byte]voice.ClassRecord

func (catalog wholePositionClassMap) LookupContext(context text.Context) (voice.ClassRecord, bool) {
	class, found := catalog[context.Paul2013SelectionKey()]
	return class, found
}

func TestLookupPaul2013WholePositionComposesOrderedVariants(t *testing.T) {
	target := text.Context{Signature: [7]byte{5: 1}}
	other := text.Context{Signature: [7]byte{5: 2}}
	missing := text.Context{Signature: [7]byte{5: 3}}
	targetKey := target.Paul2013SelectionKey()
	otherKey := other.Paul2013SelectionKey()
	otherMembers := make([]voice.UnitLocation, 9)
	for index := range otherMembers {
		otherMembers[index] = voice.UnitLocation{Bank: "gen", Index: uint32(index + 10)}
	}
	catalog := wholePositionClassMap{
		targetKey: {ID: 1, Key: targetKey, Members: []voice.UnitLocation{{Bank: "gen", Index: 0}}},
		otherKey:  {ID: 2, Key: otherKey, Members: otherMembers},
	}

	result, err := LookupPaul2013WholePosition(catalog, target, []text.Context{other, target, other, missing}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.QueryHitCount != 3 || len(result.Classes) != 2 || result.MetricSubtotal != 10 || !result.Accepted {
		t.Fatalf("result = %#v, want 3 hits, 2 unique classes, subtotal 10, accepted", result)
	}
	if result.Classes[0].ID != 1 || result.Classes[1].ID != 2 {
		t.Fatalf("class order = [%d, %d], want ascending IDs [1, 2] after native canonicalization", result.Classes[0].ID, result.Classes[1].ID)
	}
}

func TestLookupPaul2013WholePositionRequiresCatalogAndVariants(t *testing.T) {
	context := text.Context{Signature: [7]byte{5: 1}}
	if _, err := LookupPaul2013WholePosition(nil, context, []text.Context{context}, 0); err == nil {
		t.Fatal("lookup without a catalog succeeded")
	}
	if _, err := LookupPaul2013WholePosition(wholePositionClassMap{}, context, nil, 0); err == nil {
		t.Fatal("lookup without query variants succeeded")
	}
}

type featureViewRangeFixture struct {
	classes []voice.ClassRecord
	widths  []int
}

func (fixture *featureViewRangeFixture) LookupFeatureViewPrefix(_ [10]byte, _ byte, width, limit int) ([]voice.ClassRecord, error) {
	fixture.widths = append(fixture.widths, width)
	if width != 1 {
		return nil, nil
	}
	if limit < len(fixture.classes) {
		return fixture.classes[:limit], nil
	}
	return fixture.classes, nil
}

func TestLookupPaul2013WholePositionByFeatureViewCatalogRelaxesAndAccepts(t *testing.T) {
	fixture := &featureViewRangeFixture{classes: []voice.ClassRecord{
		{ID: 1, Key: [5]byte{1}, Members: make([]voice.UnitLocation, 5)},
		{ID: 2, Key: [5]byte{2}, Members: make([]voice.UnitLocation, 5)},
	}}
	query := text.Context{Signature: [7]byte{1, 2, 3}}
	target := text.Context{Signature: [7]byte{4, 5, 6}}
	result, err := LookupPaul2013WholePositionByFeatureViewCatalog(fixture, target, query, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Accepted || result.LookupPrefixLength != 1 || result.MetricSubtotal != 10 || result.LookupCandidateCount != 2 {
		t.Fatalf("catalog feature-view result = %#v, want accepted two-class range at width 1 with subtotal 10", result)
	}
	if len(fixture.widths) != 10 || fixture.widths[0] != 10 || fixture.widths[len(fixture.widths)-1] != 1 {
		t.Fatalf("queried widths = %v, want native relaxation from 10 through 1", fixture.widths)
	}
}
