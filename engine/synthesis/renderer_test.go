package synthesis

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

type fakeUnitReader struct {
	units map[string]voice.Unit
	err   error
}

func (reader fakeUnitReader) ReadUnit(bank string, index uint32) (voice.Unit, error) {
	if reader.err != nil {
		return voice.Unit{}, reader.err
	}
	unit, ok := reader.units[bank]
	if !ok || unit.Index != index {
		return voice.Unit{}, errors.New("unit not found")
	}
	return unit, nil
}

func TestConcatenatingRendererCombinesDecodedUnitSamples(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 4, PCM: []byte{0x01, 0x00, 0xfe, 0xff}},
		"num": {Index: 9, PCM: []byte{0x34, 0x12}},
	}}
	renderer := ConcatenatingRenderer{Units: reader}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "gen", Index: 4}, {Bank: "num", Index: 9},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1})
	if err != nil {
		t.Fatal(err)
	}
	if pcm.SampleRate != 16000 {
		t.Fatalf("sample rate = %d, want 16000", pcm.SampleRate)
	}
	want := []int16{1, -2, 0x1234}
	if len(pcm.Samples) != len(want) {
		t.Fatalf("sample count = %d, want %d", len(pcm.Samples), len(want))
	}
	for i := range want {
		if pcm.Samples[i] != want[i] {
			t.Errorf("sample %d = %d, want %d", i, pcm.Samples[i], want[i])
		}
	}
}

func TestConcatenatingRendererRejectsUnsupportedInputs(t *testing.T) {
	renderer := ConcatenatingRenderer{Units: fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 0, PCM: []byte{0x01, 0x00}},
	}}}
	tests := map[string]struct {
		ctx      context.Context
		units    []selection.UnitRef
		controls Controls
	}{
		"empty sequence":       {context.Background(), nil, Controls{Pitch: -1, Speed: -1, Volume: -1}},
		"non-default controls": {context.Background(), []selection.UnitRef{{Bank: "gen"}}, Controls{Pitch: 120, Speed: -1, Volume: -1}},
		"canceled context":     {canceledContext(t), []selection.UnitRef{{Bank: "gen"}}, Controls{Pitch: -1, Speed: -1, Volume: -1}},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := renderer.Render(test.ctx, test.units, test.controls); err == nil {
				t.Fatal("unsupported input accepted")
			}
		})
	}
}

func TestConcatenatingRendererRejectsMalformedUnitPCM(t *testing.T) {
	renderer := ConcatenatingRenderer{Units: fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 0, PCM: []byte{0x01}},
	}}}
	if _, err := renderer.Render(context.Background(), []selection.UnitRef{{Bank: "gen"}}, Controls{Pitch: -1, Speed: -1, Volume: -1}); err == nil {
		t.Fatal("odd byte count accepted as PCM")
	}
}

func TestCursorTimelineRendererOmitsNonFinalTrailingUPMSpan(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 4, PCM: []byte{1, 0, 2, 0, 3, 0, 4, 0}, UPM: []byte{1, 1}},
		"num": {Index: 9, PCM: []byte{5, 0, 6, 0, 7, 0, 8, 0}, UPM: []byte{1, 1}},
	}}
	renderer := CursorTimelineRenderer{Units: reader}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "gen", Index: 4}, {Bank: "num", Index: 9},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1})
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{1, 2, 5, 6, 7, 8}
	if pcm.SampleRate != 16000 || len(pcm.Samples) != len(want) {
		t.Fatalf("rendered PCM rate/count = %d/%d, want 16000/%d", pcm.SampleRate, len(pcm.Samples), len(want))
	}
	for index, sample := range want {
		if pcm.Samples[index] != sample {
			t.Errorf("sample %d = %d, want %d", index, pcm.Samples[index], sample)
		}
	}
}

func TestCursorTimelineRendererRejectsInvalidUnitsAndControls(t *testing.T) {
	valid := selection.UnitRef{Bank: "gen", Index: 0}
	tests := map[string]struct {
		reader   fakeUnitReader
		units    []selection.UnitRef
		controls Controls
	}{
		"non-default pitch": {
			reader: fakeUnitReader{units: map[string]voice.Unit{"gen": {Index: 0, PCM: []byte{1, 0, 2, 0}, UPM: []byte{1}}}},
			units:  []selection.UnitRef{valid}, controls: Controls{Pitch: 120, Speed: -1, Volume: -1},
		},
		"upm pcm mismatch": {
			reader: fakeUnitReader{units: map[string]voice.Unit{"gen": {Index: 0, PCM: []byte{1, 0, 2, 0}, UPM: []byte{2}}}},
			units:  []selection.UnitRef{valid}, controls: Controls{Pitch: -1, Speed: -1, Volume: -1},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			renderer := CursorTimelineRenderer{Units: test.reader}
			if _, err := renderer.Render(context.Background(), test.units, test.controls); err == nil {
				t.Fatal("invalid render request succeeded")
			}
		})
	}
}

func TestUPMEqualSpanJoinRendererJoinsSelectedUnits(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"one": {Index: 1, PCM: []byte{100, 0, 100, 0, 100, 0, 100, 0}, UPM: []byte{1, 1}},
		"two": {Index: 2, PCM: []byte{200, 0, 200, 0, 200, 0, 200, 0}, UPM: []byte{1, 1}},
	}}
	renderer := UPMEqualSpanJoinRenderer{Units: reader}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "one", Index: 1}, {Bank: "two", Index: 2},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1})
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{0, 50, 100, 150, 200, 100}
	if pcm.SampleRate != 16000 || len(pcm.Samples) != len(want) {
		t.Fatalf("rendered PCM rate/count = %d/%d, want 16000/%d", pcm.SampleRate, len(pcm.Samples), len(want))
	}
	for index, sample := range want {
		if pcm.Samples[index] != sample {
			t.Errorf("sample %d = %d, want %d", index, pcm.Samples[index], sample)
		}
	}
}

func TestUPMEqualSpanJoinRendererRejectsUnequalEdgeSpans(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"one": {Index: 1, PCM: []byte{100, 0, 100, 0, 100, 0, 100, 0}, UPM: []byte{1, 1}},
		"two": {Index: 2, PCM: []byte{200, 0, 200, 0, 200, 0, 200, 0, 200, 0, 200, 0, 200, 0, 200, 0}, UPM: []byte{2, 2}},
	}}
	renderer := UPMEqualSpanJoinRenderer{Units: reader}
	if _, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "one", Index: 1}, {Bank: "two", Index: 2},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1}); err == nil {
		t.Fatal("renderer accepted unequal adjacent edge spans")
	}
}

func TestUPMSegmentPlanRendererRendersPlannedPeriods(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {
			Index: 4,
			PCM:   []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10, 0, 11, 0, 12, 0},
			UPM:   []byte{2, 2, 2},
		},
	}}
	renderer := UPMSegmentPlanRenderer{Units: reader}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{{Bank: "gen", Index: 4}}, Controls{
		Pitch: 100, Speed: 100, Volume: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	if pcm.SampleRate != 16000 || len(pcm.Samples) != len(want) {
		t.Fatalf("rendered PCM rate/count = %d/%d, want 16000/%d", pcm.SampleRate, len(pcm.Samples), len(want))
	}
	for index, sample := range want {
		if pcm.Samples[index] != sample {
			t.Errorf("sample %d = %d, want %d", index, pcm.Samples[index], sample)
		}
	}
}

func TestUPMSegmentPlanRendererRejectsUPMThatDoesNotCoverPCM(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 4, PCM: []byte{1, 0, 2, 0, 3, 0, 4, 0}, UPM: []byte{1}},
	}}
	renderer := UPMSegmentPlanRenderer{Units: reader}
	if _, err := renderer.Render(context.Background(), []selection.UnitRef{{Bank: "gen", Index: 4}}, Controls{
		Pitch: -1, Speed: -1, Volume: -1,
	}); err == nil {
		t.Fatal("renderer accepted UPM metadata that covers only half the PCM")
	}
}

func TestUPMSegmentPlanRendererDefaultControlsBypassSegmentConstruction(t *testing.T) {
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 4, PCM: []byte{1, 0, 2, 0, 3, 0, 4, 0}, UPM: []byte{2}},
	}}
	renderer := UPMSegmentPlanRenderer{Units: reader}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{{Bank: "gen", Index: 4}}, Controls{
		Pitch: 100, Speed: 100, Volume: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{1, 2, 3, 4}
	if len(pcm.Samples) != len(want) {
		t.Fatalf("default-path sample count = %d, want %d", len(pcm.Samples), len(want))
	}
	for index := range want {
		if pcm.Samples[index] != want[index] {
			t.Errorf("default-path sample %d = %d, want %d", index, pcm.Samples[index], want[index])
		}
	}
}

func TestUPMSegmentPlanRendererCarriesOutputContextBetweenSteps(t *testing.T) {
	upm := []byte{54, 53, 54, 53, 53, 53, 53, 55, 52, 57, 57, 59, 59, 64}
	var sampleCount int
	for _, period := range upm {
		sampleCount += int(period) * 2
	}
	pcmBytes := make([]byte, sampleCount*2)
	samples := make([]int16, sampleCount)
	for index := range samples {
		sample := uint16((index + 1) * 100)
		samples[index] = int16(sample)
		pcmBytes[index*2] = byte(sample)
		pcmBytes[index*2+1] = byte(sample >> 8)
	}
	reader := fakeUnitReader{units: map[string]voice.Unit{
		"gen": {Index: 4, PCM: pcmBytes, UPM: upm},
	}}
	renderer := UPMSegmentPlanRenderer{Units: reader}
	got, err := renderer.Render(context.Background(), []selection.UnitRef{{Bank: "gen", Index: 4}}, Controls{
		Pitch: 120, Speed: 100, Volume: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	segments, err := BuildPaul2013UPMSegments(upm, 120, 100)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := PlanPaul2013UPMSegmentResampling(segments)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) < 2 {
		t.Fatalf("segment plan has %d steps, need multiple steps to exercise context carry", len(plan))
	}
	firstSegment := segments[plan[0].SegmentIndex]
	firstEnd := int(firstSegment.StartSample + firstSegment.FirstPeriod)
	firstBlend, err := BlendPaul2013UPMSegmentWindows(
		make([]int16, firstSegment.FirstPeriod),
		samples[firstEnd-int(plan[0].ResampledPeriodSize):firstEnd],
		int(plan[0].ResampledPeriodSize),
	)
	if err != nil {
		t.Fatal(err)
	}
	for index := range firstBlend {
		if got.Samples[index] != firstBlend[index] {
			t.Errorf("first reconstructed sample %d = %d, want %d", index, got.Samples[index], firstBlend[index])
		}
	}
	secondSegment := segments[plan[1].SegmentIndex]
	secondStart := int(secondSegment.StartSample)
	secondEnd := secondStart + int(secondSegment.FirstPeriod)
	previousSegment := segments[plan[0].SegmentIndex]
	previousSecondStart := int(previousSegment.StartSample + previousSegment.FirstPeriod)
	previousContext := samples[previousSecondStart : previousSecondStart+int(previousSegment.SecondPeriod)]
	secondBlend, err := BlendPaul2013UPMSegmentWindows(
		previousContext,
		samples[secondEnd-int(plan[1].ResampledPeriodSize):secondEnd],
		int(plan[1].ResampledPeriodSize),
	)
	if err != nil {
		t.Fatal(err)
	}
	secondOutputOffset := int(plan[0].ResampledPeriodSize)
	if len(got.Samples) < secondOutputOffset+len(secondBlend) {
		t.Fatalf("rendered PCM has %d samples, too short for the second segment at %d", len(got.Samples), secondOutputOffset)
	}
	for index := range secondBlend {
		if got.Samples[secondOutputOffset+index] != secondBlend[index] {
			t.Errorf("second reconstructed sample %d = %d, want %d", index, got.Samples[secondOutputOffset+index], secondBlend[index])
		}
	}
}

func TestPlanUPMSegmentResamplingStopsAtCumulativeBudget(t *testing.T) {
	upm := []byte{54, 53, 54, 53, 53, 53, 53, 55, 52, 57, 57, 59, 59, 64}
	segments, err := BuildPaul2013UPMSegments(upm, 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := PlanPaul2013UPMSegmentResampling(segments)
	if err != nil {
		t.Fatalf("plan segments %+v: %v", segments, err)
	}
	if len(plan) != 3 {
		t.Fatalf("plan has %d steps, want 3 before the observed cumulative budget is exhausted", len(plan))
	}
	if plan[len(plan)-1].ResampledPeriodSize != 2 {
		t.Fatalf("last planned period size = %d, want the 2-sample remainder of the budget", plan[len(plan)-1].ResampledPeriodSize)
	}
}

func TestPlanUPMSegmentResamplingCanReuseSegmentUntilBudgetEnds(t *testing.T) {
	segments, err := BuildPaul2013UPMSegments([]byte{2, 2, 2}, 120, 100)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := PlanPaul2013UPMSegmentResampling(segments)
	if err != nil {
		t.Fatal(err)
	}
	wantSegmentIndexes := []int{0, 1, 1}
	wantTargetLengths := []int32{3, 3, 2}
	if len(plan) != len(wantSegmentIndexes) {
		t.Fatalf("plan has %d steps, want %d: %+v", len(plan), len(wantSegmentIndexes), plan)
	}
	for index := range plan {
		if plan[index].SegmentIndex != wantSegmentIndexes[index] ||
			plan[index].ResampledPeriodSize != wantTargetLengths[index] {
			t.Errorf("plan step %d = %+v, want segment %d and target %d",
				index, plan[index], wantSegmentIndexes[index], wantTargetLengths[index])
		}
	}
}

func TestUPMSegmentPlanRendererWithLocalPaulUnits(t *testing.T) {
	root := filepath.Join("..", "..", "data-paul", "M16")
	paths := []string{
		filepath.Join(root, "mc_idx_tbl", "unit-gen.idx"),
		filepath.Join(root, "dat", "merged-gen.dat"),
		filepath.Join(root, "dat", "merged-gen.upm"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Skip("local Paul model files are unavailable")
		}
	}
	bank, err := voice.OpenBank(paths[0], paths[1], paths[2])
	if err != nil {
		t.Fatal(err)
	}
	defer bank.Close()
	renderer := UPMSegmentPlanRenderer{Units: oneBankUnitReader{bank: bank}}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "gen", Index: 12}, {Bank: "gen", Index: 13},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1})
	if err != nil {
		t.Fatal(err)
	}
	if pcm.SampleRate != 16000 || len(pcm.Samples) == 0 {
		t.Fatalf("local renderer PCM = %d samples at %d Hz", len(pcm.Samples), pcm.SampleRate)
	}
}

func TestUPMEqualSpanJoinRendererWithLocalPaulUnits(t *testing.T) {
	root := filepath.Join("..", "..", "data-paul", "M16")
	paths := []string{
		filepath.Join(root, "mc_idx_tbl", "unit-gen.idx"),
		filepath.Join(root, "dat", "merged-gen.dat"),
		filepath.Join(root, "dat", "merged-gen.upm"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Skip("local Paul model files are unavailable")
		}
	}
	bank, err := voice.OpenBank(paths[0], paths[1], paths[2])
	if err != nil {
		t.Fatal(err)
	}
	defer bank.Close()
	renderer := UPMEqualSpanJoinRenderer{Units: oneBankUnitReader{bank: bank}}
	pcm, err := renderer.Render(context.Background(), []selection.UnitRef{
		{Bank: "gen", Index: 0}, {Bank: "gen", Index: 1},
	}, Controls{Pitch: -1, Speed: -1, Volume: -1})
	if err != nil {
		t.Fatal(err)
	}
	if pcm.SampleRate != 16000 || len(pcm.Samples) == 0 {
		t.Fatalf("local equal-span renderer PCM = %d samples at %d Hz", len(pcm.Samples), pcm.SampleRate)
	}
}

type oneBankUnitReader struct {
	bank *voice.Bank
}

func (reader oneBankUnitReader) ReadUnit(bank string, index uint32) (voice.Unit, error) {
	if bank != "gen" {
		return voice.Unit{}, errors.New("test reader supports only gen")
	}
	return reader.bank.ReadUnit(index)
}

func canceledContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
