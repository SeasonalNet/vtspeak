package synthesis

import (
	"context"
	"encoding/binary"
	"reflect"
	"testing"

	"vtspeak/engine/dat"
	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

func TestBuildPaul2013ContextUPMPeriodsAlignsNeighborEdges(t *testing.T) {
	current := contextSource([]byte{2, 3, 4, 5})
	left := contextSource([]byte{6, 7, 8, 9, 10, 11})
	right := contextSource([]byte{1, 2, 3, 4, 5, 6})
	plan := Paul2013ContextEdgePlan{LeftContextCount: 2, RightContextCount: 2}
	got, err := buildPaul2013ContextUPMPeriods(&current, &left, &right, plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 {
		t.Fatalf("period count = %d, want 4", len(got))
	}
	want := []Paul2013UPMPeriod{
		{Length: 4, Windows: Paul2013UPMPeriodWindows{
			Current: current.samples[0:4], Left: left.samples[0:12],
		}},
		{Length: 6, Windows: Paul2013UPMPeriodWindows{
			Current: current.samples[4:10], Left: left.samples[12:26],
		}},
		{Length: 8, Windows: Paul2013UPMPeriodWindows{
			Current: current.samples[10:18], Right: right.samples[20:30],
		}},
		{Length: 10, Windows: Paul2013UPMPeriodWindows{
			Current: current.samples[18:28], Right: right.samples[30:42],
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("UPM period windows = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013ContextTimelineRowInputsPreservesNormalRowFields(t *testing.T) {
	current := Paul2013NormalTimelineRow{
		Unit: selection.UnitRef{Bank: "gen", Index: 0},
		Mode: Paul2013TimelineCombinedView, RowKind: 2,
		SampleCount: 10, LeadingSpan: 2, TrailingSpan: 4,
		Controls: Paul2013TimelineControlWords{Gain: 100},
	}
	left := current
	left.Unit = selection.UnitRef{Bank: "num", Index: 1}
	left.Mode = Paul2013TimelineFirstSideView
	right := current
	right.Unit = selection.UnitRef{Bank: "etc", Index: 2}
	rows, err := BuildPaul2013ContextTimelineRowInputs([]Paul2013NormalContextTimelineRowInput{{
		Current: current, Gate: ContextGateInput{TimelineSpeedControl: 100},
		Left: &left, Right: &right,
		CurrentScratchSamples: []int16{11}, LeftScratchSamples: []int16{12},
		RightScratchSamples: []int16{13},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Current.Unit != current.Unit || rows[0].Current.Mode != current.Mode ||
		rows[0].Current.SampleCount != 10 || rows[0].Current.LeadingSpan != 2 ||
		rows[0].Current.TrailingSpan != 4 || rows[0].Current.GainPercent != 100 {
		t.Fatalf("current row mapping = %+v", rows)
	}
	if rows[0].Left == nil || rows[0].Left.Unit != left.Unit || rows[0].Left.Mode != left.Mode ||
		!reflect.DeepEqual(rows[0].Left.SourceScratchSamples, []int16{12}) {
		t.Fatalf("left row mapping = %+v", rows[0].Left)
	}
	if rows[0].Right == nil || rows[0].Right.Unit != right.Unit ||
		!reflect.DeepEqual(rows[0].Right.SourceScratchSamples, []int16{13}) {
		t.Fatalf("right row mapping = %+v", rows[0].Right)
	}
}

func TestBuildPaul2013ContextTimelineRowInputsRejectsSyntheticRows(t *testing.T) {
	_, err := BuildPaul2013ContextTimelineRowInputs([]Paul2013NormalContextTimelineRowInput{{
		Current: Paul2013NormalTimelineRow{RowKind: 1},
	}})
	if err == nil {
		t.Fatal("synthetic boundary row was accepted as a normal context row")
	}
}

func TestRenderPaul2013ContextTimelineRowsComposesContextAndJoin(t *testing.T) {
	unitData := map[string]voice.Unit{
		"gen": contextTimelineUnit(0, 100),
		"num": contextTimelineUnit(1, 300),
		"etc": contextTimelineUnit(2, 500),
	}
	reader := fakeUnitReader{units: unitData}
	current := contextTimelineInput("gen", 0)
	left := contextTimelineInput("num", 1)
	right := contextTimelineInput("etc", 2)
	gateRows := make([][7]byte, 5)
	for index := range gateRows {
		gateRows[index][6] = 0x80
	}
	gateRows[2][6] = 0
	row := Paul2013ContextTimelineRowInput{
		Current: current,
		Gate: ContextGateInput{
			Rows: gateRows, Modes: []byte{0, 0, 0, 0, 0},
			CurrentIndex: 2, LeftIndex: 0, RightIndex: 4,
			TimelineSpeedControl: 100,
		},
		Left: &left, Right: &right,
	}
	got, err := RenderPaul2013ContextTimelineRows(context.Background(), reader, []Paul2013ContextTimelineRowInput{row})
	if err != nil {
		t.Fatal(err)
	}
	if got.SampleRate != 16000 || len(got.Samples) == 0 {
		t.Fatalf("context-rendered PCM rate/count = %d/%d", got.SampleRate, len(got.Samples))
	}

	currentSource, err := readPaul2013ContextTimelineSource(context.Background(), reader, current)
	if err != nil {
		t.Fatal(err)
	}
	leftSource, err := readPaul2013ContextTimelineSource(context.Background(), reader, left)
	if err != nil {
		t.Fatal(err)
	}
	rightSource, err := readPaul2013ContextTimelineSource(context.Background(), reader, right)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPaul2013ContextEdgePlan(row.Gate, len(currentSource.upm), len(leftSource.upm), len(rightSource.upm))
	if err != nil {
		t.Fatal(err)
	}
	periods, err := buildPaul2013ContextUPMPeriods(&currentSource, &leftSource, &rightSource, plan)
	if err != nil {
		t.Fatal(err)
	}
	wantRow := make([]int16, len(currentSource.samples))
	if err := MixPaul2013UPMTimeline(wantRow, 0, plan.Weights, periods); err != nil {
		t.Fatal(err)
	}
	wantSamples, err := MixPaul2013Timeline([]TimelinePCMRow{{
		Samples: wantRow, LeadingSpan: current.LeadingSpan, TrailingSpan: current.TrailingSpan,
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Samples, wantSamples) {
		t.Fatalf("context renderer output differs from composed row path: got %d samples, want %d", len(got.Samples), len(wantSamples))
	}
}

func TestRenderPaul2013ContextTimelineRowsPreservesNoContextRows(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{"gen": contextTimelineUnit(0, 100)}}
	current := contextTimelineInput("gen", 0)
	row := Paul2013ContextTimelineRowInput{
		Current: current,
		Gate: ContextGateInput{
			Rows: make([][7]byte, 1), Modes: []byte{0}, CurrentIndex: 0,
			TimelineSpeedControl: 100,
		},
	}
	got, err := RenderPaul2013ContextTimelineRows(context.Background(), reader, []Paul2013ContextTimelineRowInput{row})
	if err != nil {
		t.Fatal(err)
	}
	want, err := RenderPaul2013TimelineRows(context.Background(), reader, []Paul2013TimelinePCMInput{current})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("closed context gates changed the current row: got %v, want %v", got.Samples, want.Samples)
	}
}

func TestReadPaul2013ContextTimelineSourceRejectsRowViewMismatch(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{"gen": contextTimelineUnit(0, 100)}}
	row := contextTimelineInput("gen", 0)
	row.TrailingSpan++
	if _, err := readPaul2013ContextTimelineSource(context.Background(), reader, row); err == nil {
		t.Fatal("timeline edge span that disagrees with the decoded unit view was accepted")
	}
}

func contextSource(upm []byte) paul2013ContextTimelineSource {
	sampleCount := 0
	for _, period := range upm {
		sampleCount += int(period) * 2
	}
	samples := make([]int16, sampleCount)
	for index := range samples {
		samples[index] = int16(index + 1)
	}
	return paul2013ContextTimelineSource{samples: samples, upm: append([]byte(nil), upm...)}
}

func contextTimelineInput(bank string, index uint32) Paul2013TimelinePCMInput {
	return Paul2013TimelinePCMInput{
		Unit: selection.UnitRef{Bank: bank, Index: index}, RowType: 2,
		Mode: Paul2013TimelineCombinedView, SampleCount: 10,
		LeadingSpan: 2, TrailingSpan: 4, GainPercent: 100,
	}
}

func contextTimelineUnit(index uint32, firstSample int16) voice.Unit {
	upm := []byte{1, 2, 2}
	pcm := make([]byte, 20)
	for sampleIndex := 0; sampleIndex < len(pcm)/2; sampleIndex++ {
		binary.LittleEndian.PutUint16(pcm[sampleIndex*2:], uint16(firstSample+int16(sampleIndex)))
	}
	return voice.Unit{
		Index: index,
		Record: dat.UnitRecord{
			FirstSideSamples: 6, SecondSideSamples: 8, DATLength: 5,
			UPMFirstCount: 2, UPMSecondCount: 2,
			UPMEdges: [3]byte{1, 2, 2},
		},
		UPM: upm, PCM: pcm,
	}
}
