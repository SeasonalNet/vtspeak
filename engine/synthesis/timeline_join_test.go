package synthesis

import (
	"reflect"
	"testing"
)

func TestBuildPaul2013TimelinePCMRowUsesCapturedModeOffsets(t *testing.T) {
	source := []int16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	tests := []struct {
		name string
		mode uint8
		want []int16
	}{
		{name: "combined mode starts at unit start", mode: 0, want: []int16{0, 1, 2, 3}},
		{name: "first-side mode starts at unit start", mode: 1, want: []int16{0, 1, 2, 3}},
		{name: "second-side mode selects trailing row window", mode: 2, want: []int16{10, 11}},
		{name: "unknown mode follows native non-second-side branch", mode: 3, want: []int16{0, 1, 2, 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			row, err := BuildPaul2013TimelinePCMRow(source, 12, test.mode, len(test.want), 2, 1)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(row.Samples, test.want) {
				t.Fatalf("row samples = %v, want %v", row.Samples, test.want)
			}
			if row.LeadingSpan != 2 || row.TrailingSpan != 1 {
				t.Fatalf("row spans = %d/%d, want 2/1", row.LeadingSpan, row.TrailingSpan)
			}
			row.Samples[0] = -1
			if source[0] == -1 || (test.mode == 2 && source[10] == -1) {
				t.Fatal("row extraction aliases the decoded source buffer")
			}
		})
	}
}

func TestBuildPaul2013TimelinePCMRowAllowsObservedScratchTail(t *testing.T) {
	source := []int16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	row, err := BuildPaul2013TimelinePCMRow(source, 12, 2, 4, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int16{10, 11, 12, 13}; !reflect.DeepEqual(row.Samples, want) {
		t.Fatalf("type-two scratch window = %v, want %v", row.Samples, want)
	}
}

func TestBuildPaul2013TimelinePCMRowRejectsOutOfBoundsRows(t *testing.T) {
	source := make([]int16, 12)
	for name, args := range map[string]struct {
		mode         uint8
		decodedCount int
		count        int
		leading      int
		trailing     int
	}{
		"zero count":                              {mode: 0, decodedCount: 12, count: 0, leading: 0, trailing: 0},
		"count does not fit combined view":        {mode: 0, decodedCount: 12, count: 13, leading: 1, trailing: 1},
		"second-side offset and count do not fit": {mode: 2, decodedCount: 12, count: 3, leading: 1, trailing: 1},
		"leading exceeds decoded samples":         {mode: 2, decodedCount: 12, count: 1, leading: 13, trailing: 0},
		"decoded count exceeds buffer":            {mode: 1, decodedCount: 13, count: 1, leading: 0, trailing: 0},
		"edge exceeds row":                        {mode: 1, decodedCount: 12, count: 2, leading: 3, trailing: 0},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := BuildPaul2013TimelinePCMRow(source, args.decodedCount, args.mode, args.count, args.leading, args.trailing); err == nil {
				t.Fatalf("accepted invalid timeline row: %+v", args)
			}
		})
	}
}

func TestScalePaul2013TimelinePCMRowUsesObservedIntegerGain(t *testing.T) {
	row := TimelinePCMRow{
		Samples:     []int16{1001, -1001, 20000, -20000},
		LeadingSpan: 1, TrailingSpan: 1,
	}
	scaled, err := ScalePaul2013TimelinePCMRow(row, 200)
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{2002, -2002, 32766, -32767}
	if !reflect.DeepEqual(scaled.Samples, want) {
		t.Fatalf("scaled samples = %v, want %v", scaled.Samples, want)
	}
	if !reflect.DeepEqual(row.Samples, []int16{1001, -1001, 20000, -20000}) {
		t.Fatalf("scaling changed the source row: %v", row.Samples)
	}
	if scaled.LeadingSpan != row.LeadingSpan || scaled.TrailingSpan != row.TrailingSpan {
		t.Fatalf("scaled row spans = %d/%d, want %d/%d", scaled.LeadingSpan, scaled.TrailingSpan, row.LeadingSpan, row.TrailingSpan)
	}

	truncated, err := ScalePaul2013TimelinePCMRow(TimelinePCMRow{Samples: []int16{1001, -1001}}, 50)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(truncated.Samples, []int16{500, -500}) {
		t.Fatalf("fractional gain samples = %v, want [500 -500]", truncated.Samples)
	}
}

func TestScalePaul2013TimelinePCMRowRejectsInvalidInput(t *testing.T) {
	for name, row := range map[string]TimelinePCMRow{
		"empty":     {},
		"bad spans": {Samples: []int16{1, 2}, LeadingSpan: 3},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ScalePaul2013TimelinePCMRow(row, 200); err == nil {
				t.Fatalf("accepted invalid row: %+v", row)
			}
		})
	}
	for _, gain := range []int32{-1, 501} {
		if _, err := ScalePaul2013TimelinePCMRow(TimelinePCMRow{Samples: []int16{1}}, gain); err == nil {
			t.Fatalf("accepted gain %d", gain)
		}
	}
}

func TestMixPaul2013EqualSpanTimelineBlendsRowsAndPreservesCursorLength(t *testing.T) {
	rows := []TimelinePCMRow{
		{Samples: []int16{100, 100, 11, 12, 100, 100}, LeadingSpan: 2, TrailingSpan: 2},
		{Samples: []int16{200, 200, 23, 24, 200, 200}, LeadingSpan: 2, TrailingSpan: 2},
	}
	got, err := MixPaul2013EqualSpanTimeline(rows)
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{0, 50, 11, 12, 100, 150, 23, 24, 200, 100}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("joined samples = %v, want %v", got, want)
	}
	frames, err := Paul2013TimelineOutputFrames([]TimelineRow{
		{SampleCount: 6, LeadingSpan: 2, TrailingSpan: 2},
		{SampleCount: 6, LeadingSpan: 2, TrailingSpan: 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != frames {
		t.Fatalf("joined frame count = %d, cursor plan = %d", len(got), frames)
	}
}

func TestMixPaul2013EqualSpanTimelineRejectsInvalidRows(t *testing.T) {
	for name, rows := range map[string][]TimelinePCMRow{
		"empty":              nil,
		"empty row":          {{LeadingSpan: 0, TrailingSpan: 0}},
		"invalid edge spans": {{Samples: []int16{1, 2}, LeadingSpan: 2, TrailingSpan: 1}},
		"unequal boundary spans": {
			{Samples: []int16{1, 2, 3}, LeadingSpan: 0, TrailingSpan: 1},
			{Samples: []int16{4, 5, 6}, LeadingSpan: 2, TrailingSpan: 0},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := MixPaul2013EqualSpanTimeline(rows); err == nil {
				t.Fatalf("accepted invalid timeline rows: %+v", rows)
			}
		})
	}
}

func TestPaul2013WeightedLaneSaturatesAndTruncates(t *testing.T) {
	if got := timelineJoinLane(32767, 0.999999); got != 32766 {
		t.Fatalf("positive weighted lane = %d, want 32766", got)
	}
	if got := timelineJoinLane(-32768, 0.999999); got != -32767 {
		t.Fatalf("negative weighted lane = %d, want -32767", got)
	}
	if got := timelineJoinLane(1234, 0.5); got != 617 {
		t.Fatalf("ordinary weighted lane = %d, want 617", got)
	}
	if got := timelineJoinLane(1235, 0.5); got != 617 {
		t.Fatalf("positive fractional weighted lane = %d, want truncation to 617", got)
	}
	if got := timelineJoinLane(-1235, 0.5); got != -617 {
		t.Fatalf("negative fractional weighted lane = %d, want truncation to -617", got)
	}
}

func TestApplyPaul2013BlendWindowSaturatesFractionalEndpoints(t *testing.T) {
	positive := make([]int16, 4097)
	positive[len(positive)-1] = 32767
	gotPositive, err := ApplyPaul2013BlendWindow(positive, len(positive), WindowAlignLeft)
	if err != nil {
		t.Fatal(err)
	}
	if got := gotPositive[len(gotPositive)-1]; got != 32766 {
		t.Fatalf("positive blend endpoint = %d, want 32766", got)
	}

	negative := make([]int16, 4097)
	negative[len(negative)-1] = -32768
	gotNegative, err := ApplyPaul2013BlendWindow(negative, len(negative), WindowAlignLeft)
	if err != nil {
		t.Fatal(err)
	}
	if got := gotNegative[len(gotNegative)-1]; got != -32767 {
		t.Fatalf("negative blend endpoint = %d, want -32767", got)
	}
}

func TestBlendPaul2013UPMSegmentWindowsUsesAlignedTableHalves(t *testing.T) {
	got, err := BlendPaul2013UPMSegmentWindows(
		[]int16{100, 100},
		[]int16{40, 40},
		4,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{100, 50, 0, 20}
	if len(got) != len(want) {
		t.Fatalf("blended window length = %d, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Errorf("blended sample %d = %d, want %d", index, got[index], want[index])
		}
	}
}

func TestBlendPaul2013UPMSegmentWindowsCropsAtNativeEdges(t *testing.T) {
	got, err := BlendPaul2013UPMSegmentWindows(
		[]int16{100, 100, 100, 100},
		[]int16{10, 20, 30, 40},
		2,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{100, 70}
	for index := range want {
		if got[index] != want[index] {
			t.Errorf("cropped blended sample %d = %d, want %d", index, got[index], want[index])
		}
	}
}
