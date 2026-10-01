package selection

import (
	"testing"

	"vtspeak/engine/text"
	"vtspeak/engine/voice"
)

type catalogFeatureQueryFixture struct {
	featureViewRangeFixture
}

func (*catalogFeatureQueryFixture) LookupPrefixLimited([5]byte, int, int) ([]voice.ClassRecord, error) {
	return nil, nil
}

func TestRunPaul2013CatalogQuerySequenceUsesFeatureViewIndexByDefault(t *testing.T) {
	catalog := &catalogFeatureQueryFixture{featureViewRangeFixture: featureViewRangeFixture{
		classes: []voice.ClassRecord{
			{ID: 1, Key: [5]byte{1}, Members: make([]voice.UnitLocation, 5)},
			{ID: 2, Key: [5]byte{2}, Members: make([]voice.UnitLocation, 5)},
		},
	}}
	contextRow := [6]byte{4: 1}
	result, err := RunPaul2013CatalogQuerySequence(
		catalog,
		[7]byte{},
		text.Context{Signature: [7]byte{1, 2, 3}},
		contextRow,
		nil,
		0,
		1,
		func() byte { return 0 },
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog.widths) == 0 {
		t.Fatal("feature-view catalog index was not queried")
	}
	if len(result.Candidates) == 0 {
		t.Fatalf("catalog query produced no feature-view candidates: %+v", result.Sequence)
	}
}

func TestQueryPaul2013FallbackRowsRereadsPriorMetricPerTarget(t *testing.T) {
	catalog := &catalogFeatureQueryFixture{}
	readCounts := [2]int{}
	runtime := [2]Paul2013FallbackRowRuntime{}
	for rowIndex := range runtime {
		index := rowIndex
		runtime[index] = Paul2013FallbackRowRuntime{
			ReadPriorMetric: func() uint32 {
				readCounts[index]++
				return uint32(readCounts[index])
			},
			ContextGate: func() byte { return 0 },
		}
	}
	_, err := QueryPaul2013FallbackRows(
		catalog,
		text.Context{Signature: [7]byte{1, 2, 3, 4, 5, 0x39, 6}},
		[6]byte{}, 0, 12, runtime,
	)
	if err != nil {
		t.Fatal(err)
	}
	if readCounts != [2]int{2, 2} {
		t.Fatalf("prior-metric reads per fallback row = %v, want one read per each of two tree targets", readCounts)
	}
}
