package synthesis

import (
	"reflect"
	"strings"
	"testing"

	"vtspeak/engine/dat"
	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

func TestBuildPaul2013TimelineUnitViewSelectsUPMSides(t *testing.T) {
	unit := voice.Unit{
		Index: 7,
		Record: dat.UnitRecord{
			DATOffset: 1000, FirstSideSamples: 24, SecondSideSamples: 36,
			DATLength: 64, UPMOffset: 2000, UPMFirstCount: 2,
			UPMSecondCount: 2, UPMEdges: [3]byte{5, 7, 11},
		},
		UPM: []byte{5, 7, 11},
		PCM: make([]byte, 46*2),
	}
	reference := selection.UnitRef{Bank: "gen", Index: 7}
	cases := []struct {
		name string
		mode byte
		want Paul2013TimelineUnitView
	}{
		{
			name: "combined", mode: Paul2013TimelineCombinedView,
			want: Paul2013TimelineUnitView{
				Unit: reference, FileIndex: 3, Mode: Paul2013TimelineCombinedView,
				DATOffset: 1000, SampleCount: 46, UPMOffset: 2000,
				UPMCount: 3, LeadingPeriod: 5, TrailingPeriod: 11,
			},
		},
		{
			name: "first side", mode: Paul2013TimelineFirstSideView,
			want: Paul2013TimelineUnitView{
				Unit: reference, FileIndex: 3, Mode: Paul2013TimelineFirstSideView,
				DATOffset: 1000, SampleCount: 24, UPMOffset: 2000,
				UPMCount: 2, LeadingPeriod: 5, TrailingPeriod: 7,
			},
		},
		{
			name: "second side", mode: Paul2013TimelineSecondSideView,
			want: Paul2013TimelineUnitView{
				Unit: reference, FileIndex: 3, Mode: Paul2013TimelineSecondSideView,
				DATOffset: 1010, SampleCount: 36, UPMOffset: 2001,
				UPMCount: 2, LeadingPeriod: 7, TrailingPeriod: 11,
			},
		},
		{
			name: "native combined fallback", mode: 9,
			want: Paul2013TimelineUnitView{
				Unit: reference, FileIndex: 3, Mode: 9,
				DATOffset: 1000, SampleCount: 46, UPMOffset: 2000,
				UPMCount: 3, LeadingPeriod: 5, TrailingPeriod: 11,
			},
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := BuildPaul2013TimelineUnitView(unit, reference, 3, test.mode)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("timeline unit view = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestBuildPaul2013TimelineUnitViewRejectsInconsistentMetadata(t *testing.T) {
	unit := voice.Unit{
		Record: dat.UnitRecord{
			DATOffset: 1000, FirstSideSamples: 1, SecondSideSamples: 36,
			UPMOffset: 2000, UPMFirstCount: 2, UPMSecondCount: 2,
			UPMEdges: [3]byte{5, 7, 11},
		},
		UPM: []byte{5, 7, 11}, PCM: make([]byte, 46*2),
	}
	_, err := BuildPaul2013TimelineUnitView(unit, selection.UnitRef{Bank: "gen"}, 0, Paul2013TimelineSecondSideView)
	if err == nil || !strings.Contains(err.Error(), "UPM sides") {
		t.Fatalf("inconsistent sample metadata error = %v", err)
	}
}

func TestBuildPaul2013NormalTimelineRowMapsCapturedFields(t *testing.T) {
	unit := voice.Unit{
		Record: dat.UnitRecord{
			DATOffset: 1000, FirstSideSamples: 24, SecondSideSamples: 36,
			UPMOffset: 2000, UPMFirstCount: 2, UPMSecondCount: 2,
			UPMEdges: [3]byte{5, 7, 11},
		},
		UPM: []byte{5, 7, 11}, PCM: make([]byte, 46*2),
	}
	reference := selection.UnitRef{Bank: "gen", Index: 7}
	controls := Paul2013TimelineControlWords{Primary: 100, Speed: 120, Gain: 200}
	row, err := BuildPaul2013NormalTimelineRow(
		unit, reference, 3, Paul2013TimelineSecondSideView, 4, controls,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013NormalTimelineRow{
		Controls: controls, SampleCount: 36, PrimaryIndex: 4, Unit: reference,
		DATOffset: 1010, UPMOffset: 2001, UPMCount: 2,
		LeadingSpan: 14, TrailingSpan: 22,
		Mode: Paul2013TimelineSecondSideView, RowKind: 2, FileIndex: 3,
	}
	if !reflect.DeepEqual(row, want) {
		t.Fatalf("normal timeline row = %+v, want %+v", row, want)
	}
}

func TestBuildPaul2013SyntheticTimelineRowMapsObservedFields(t *testing.T) {
	controls := Paul2013TimelineControlWords{Primary: 100, Speed: 100, Gain: 200}
	row, err := BuildPaul2013SyntheticTimelineRow(
		SyntheticTimelineRow{LookupValue: 1000, StateValue: 105}, 12, controls,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013SyntheticTimelineRow{
		Controls: controls, SampleCount: 11, OpaqueIndex: 12, PrimaryIndex: -1,
		SentinelAt2A: 0xffff, SentinelAt2C: 0xffff,
		RowKind: 1, TrailingMarker: 0x5a,
	}
	if !reflect.DeepEqual(row, want) {
		t.Fatalf("synthetic timeline row = %+v, want %+v", row, want)
	}
}
