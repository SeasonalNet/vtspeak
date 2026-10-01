package synthesis

import (
	"context"
	"math"
	"testing"

	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

func TestPaul2013TimelineOutputFramesWithSyntheticRow(t *testing.T) {
	normalBefore := TimelineRow{SampleCount: 100, LeadingSpan: 8, TrailingSpan: 10}
	normalFinal := TimelineRow{SampleCount: 50, LeadingSpan: 6, TrailingSpan: 7}
	synthetic := SyntheticTimelineRow{LookupValue: 1000, StateValue: 105}
	rows := []TimelineOutputRow{
		{Normal: &normalBefore},
		{Synthetic: &synthetic},
		{Normal: &normalFinal},
	}
	got, err := Paul2013TimelineOutputFramesWithSynthetic(rows)
	if err != nil {
		t.Fatal(err)
	}
	if got != 151 {
		t.Fatalf("timeline output frames = %d, want 151 (90 + 11 + 50)", got)
	}

	plain, err := Paul2013TimelineOutputFrames([]TimelineRow{normalBefore, normalFinal})
	if err != nil {
		t.Fatal(err)
	}
	withNormalKinds, err := Paul2013TimelineOutputFramesWithSynthetic([]TimelineOutputRow{
		{Normal: &normalBefore}, {Normal: &normalFinal},
	})
	if err != nil {
		t.Fatal(err)
	}
	if withNormalKinds != plain {
		t.Fatalf("normal-only output frames = %d, wrapper returned %d", withNormalKinds, plain)
	}
}

func TestBuildPaul2013TimelinePCMInputsFromNormalRows(t *testing.T) {
	rows := []Paul2013NormalTimelineRow{{
		Controls:    Paul2013TimelineControlWords{Primary: 100, Speed: 100, Gain: 120},
		SampleCount: 240, PrimaryIndex: 3, Unit: selection.UnitRef{Bank: "gen", Index: 12},
		LeadingSpan: 8, TrailingSpan: 10, RowKind: 2,
		Mode: Paul2013TimelineSecondSideView,
	}}
	got, err := BuildPaul2013TimelinePCMInputs(rows)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013TimelinePCMInput{
		Unit: selection.UnitRef{Bank: "gen", Index: 12}, RowType: 2,
		Mode:        Paul2013TimelineSecondSideView,
		SampleCount: 240, LeadingSpan: 8, TrailingSpan: 10, GainPercent: 120,
	}
	if len(got) != 1 || got[0].Unit != want.Unit || got[0].RowType != want.RowType || got[0].Mode != want.Mode ||
		got[0].SampleCount != want.SampleCount || got[0].LeadingSpan != want.LeadingSpan ||
		got[0].TrailingSpan != want.TrailingSpan || got[0].GainPercent != want.GainPercent ||
		len(got[0].SourceScratchSamples) != 0 {
		t.Fatalf("timeline PCM inputs = %+v, want [%+v]", got, want)
	}
	if _, err := BuildPaul2013TimelinePCMInputs(nil); err == nil {
		t.Fatal("empty normal rows accepted")
	}
	rows[0].RowKind = 1
	if _, err := BuildPaul2013TimelinePCMInputs(rows); err == nil {
		t.Fatal("synthetic row accepted as a normal unit row")
	}
}

func TestRenderPaul2013TimelineRowsUsesViewModeAndCallerScratch(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 12, PCM: []byte{10, 0, 20, 0, 30, 0, 40, 0}},
	}}
	for _, test := range []struct {
		name string
		mode uint8
		want []int16
	}{
		{name: "combined starts at unit start", mode: Paul2013TimelineCombinedView, want: []int16{0, 20}},
		{name: "first side starts at unit start", mode: Paul2013TimelineFirstSideView, want: []int16{0, 20}},
		{name: "second side uses trailing offset", mode: Paul2013TimelineSecondSideView, want: []int16{0, 50}},
	} {
		t.Run(test.name, func(t *testing.T) {
			rows := []Paul2013TimelinePCMInput{{
				Unit: selection.UnitRef{Bank: "gen", Index: 12}, RowType: 2, Mode: test.mode,
				SampleCount: 2, LeadingSpan: 1, TrailingSpan: 1, GainPercent: 100,
				SourceScratchSamples: []int16{50, 60},
			}}
			got, err := RenderPaul2013TimelineRows(context.Background(), reader, rows)
			if err != nil {
				t.Fatal(err)
			}
			if got.SampleRate != 16000 || len(got.Samples) != len(test.want) {
				t.Fatalf("rendered PCM rate/count = %d/%d, want 16000/%d", got.SampleRate, len(got.Samples), len(test.want))
			}
			for index := range test.want {
				if got.Samples[index] != test.want[index] {
					t.Errorf("sample %d = %d, want %d", index, got.Samples[index], test.want[index])
				}
			}
		})
	}
}

func TestPaul2013TimelineOutputFramesWithSyntheticRejectsInvalidRows(t *testing.T) {
	normal := TimelineRow{SampleCount: 10}
	maxLength := TimelineRow{SampleCount: math.MaxInt}
	synthetic := SyntheticTimelineRow{LookupValue: 0, StateValue: 10}
	for name, rows := range map[string][]TimelineOutputRow{
		"empty":       nil,
		"no kind":     {{}},
		"two kinds":   {{Normal: &normal, Synthetic: &synthetic}},
		"zero lookup": {{Synthetic: &synthetic}},
		"bad normal":  {{Normal: &TimelineRow{SampleCount: 0}}},
		"overflow":    {{Normal: &maxLength}, {Synthetic: &SyntheticTimelineRow{LookupValue: 20000, StateValue: 50}}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Paul2013TimelineOutputFramesWithSynthetic(rows); err == nil {
				t.Fatalf("accepted invalid timeline rows: %+v", rows)
			}
		})
	}
}
