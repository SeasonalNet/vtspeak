package text

import (
	"reflect"
	"testing"
)

func TestSplitPaul2013TokenPhoneBlocks(t *testing.T) {
	tests := map[string]struct {
		markers []byte
		want    []Paul2013PhoneBlock
	}{
		"empty token": {},
		"default markers close only at token end": {
			markers: []byte{'0', '0', '0'},
			want:    []Paul2013PhoneBlock{{Start: 0, End: 3, TerminalMarker: '0'}},
		},
		"nondefault markers close blocks": {
			markers: []byte{'0', '0', '^', '0', '[', '0'},
			want: []Paul2013PhoneBlock{
				{Start: 0, End: 3, TerminalMarker: '^'},
				{Start: 3, End: 5, TerminalMarker: '['},
				{Start: 5, End: 6, TerminalMarker: '0'},
			},
		},
		"consecutive markers close single phone blocks": {
			markers: []byte{'Z', '^'},
			want: []Paul2013PhoneBlock{
				{Start: 0, End: 1, TerminalMarker: 'Z'},
				{Start: 1, End: 2, TerminalMarker: '^'},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := SplitPaul2013TokenPhoneBlocks(test.markers)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(test.want) {
				t.Fatalf("got %d blocks %+v, want %+v", len(got), got, test.want)
			}
			for index := range test.want {
				if got[index] != test.want[index] {
					t.Errorf("block %d = %+v, want %+v", index, got[index], test.want[index])
				}
			}
		})
	}
}

func TestBuildPaul2013InitialTokenMarkers(t *testing.T) {
	got, err := BuildPaul2013InitialTokenMarkers([]int16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{']', '[', 'Z', '^', 'Z', '[', ']', '[', '[', ']', '[', ']', '\\'}; !reflect.DeepEqual(got, want) {
		t.Fatalf("initial token markers = %q, want %q", got, want)
	}
	for _, stateCode := range []int16{-1, 13} {
		if _, err := LookupPaul2013InitialTokenMarker(stateCode); err == nil {
			t.Errorf("accepted initial marker state %d", stateCode)
		}
	}
}

func TestSplitPaul2013TokenPhoneBlocksRejectsOversizedToken(t *testing.T) {
	if _, err := SplitPaul2013TokenPhoneBlocks(make([]byte, 66)); err == nil {
		t.Fatal("66-phone token accepted by the 65-phone block buffer")
	}
}

func TestApplyPaul2013PhoneBlockDurationLimitInsertsBoundary(t *testing.T) {
	markers := []byte{'0', '0', '0', '0'}
	got, err := ApplyPaul2013PhoneBlockDurationLimit(markers, []byte{250, 250, 1, 10})
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{'0', '[', '0', '0'}; !reflect.DeepEqual(got, want) {
		t.Fatalf("limited phone markers = %q, want %q", got, want)
	}
	if want := []byte{'0', '0', '0', '0'}; !reflect.DeepEqual(markers, want) {
		t.Fatalf("input markers mutated to %q", markers)
	}
	blocks, err := SplitPaul2013TokenPhoneBlocks(got)
	if err != nil {
		t.Fatal(err)
	}
	wantBlocks := []Paul2013PhoneBlock{
		{Start: 0, End: 2, TerminalMarker: '['},
		{Start: 2, End: 4, TerminalMarker: '0'},
	}
	if !reflect.DeepEqual(blocks, wantBlocks) {
		t.Fatalf("duration-limited blocks = %+v, want %+v", blocks, wantBlocks)
	}
}

func TestApplyPaul2013PhoneBlockDurationLimitHonorsExistingBoundaries(t *testing.T) {
	got, err := ApplyPaul2013PhoneBlockDurationLimit(
		[]byte{'0', '^', '0'}, []byte{250, 250, 250},
	)
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{'0', '^', '0'}; !reflect.DeepEqual(got, want) {
		t.Fatalf("existing boundary markers = %q, want %q", got, want)
	}
}

func TestApplyPaul2013PhoneBlockDurationLimitRejectsMismatchedVectors(t *testing.T) {
	if _, err := ApplyPaul2013PhoneBlockDurationLimit([]byte{'0'}, nil); err == nil {
		t.Fatal("mismatched marker and duration vectors were accepted")
	}
}

func TestBuildPaul2013MarkerTreeGroupsCombinesSlashBoundaries(t *testing.T) {
	got, err := BuildPaul2013MarkerTreeGroups(
		[]byte{'\\', '0', ']', '0'}, []byte{2, 3, 5, 7},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013MarkerTreeGroup{
		{PhoneStart: 0, PhoneEnd: 2, DurationByteSum: 5},
		{PhoneStart: 2, PhoneEnd: 4, DurationByteSum: 12},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("marker-tree groups = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013MarkerTreeGroupsRejectsMismatchedOrOversizedInput(t *testing.T) {
	if _, err := BuildPaul2013MarkerTreeGroups([]byte{'0'}, nil); err == nil {
		t.Fatal("mismatched marker and duration vectors were accepted")
	}
	if _, err := BuildPaul2013MarkerTreeGroups(make([]byte, 66), make([]byte, 66)); err == nil {
		t.Fatal("66-phone marker group input was accepted")
	}
}

func TestBuildPaul2013TokenBoundaryMarkerTransitions(t *testing.T) {
	initial := []Paul2013TokenBoundaryMarker{
		{Marker: ']', PreviousStateFlag: false},
		{Marker: '\\', PreviousStateFlag: false},
		{Marker: '[', PreviousStateFlag: false},
	}
	got, err := BuildPaul2013TokenBoundaryMarkerTransitions(
		initial,
		[]byte{7, 8, 8},
		[]int16{0, 12, 12},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013TokenBoundaryMarker{
		{Marker: '\\', PreviousStateFlag: true},
		{Marker: '\\', PreviousStateFlag: true},
		{Marker: '[', PreviousStateFlag: false},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("marker transitions = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013TokenBoundaryMarkerTransitionsHonorsTokenCodeGate(t *testing.T) {
	initial := []Paul2013TokenBoundaryMarker{
		{Marker: ']', PreviousStateFlag: false},
		{Marker: 'Z', PreviousStateFlag: false},
	}
	got, err := BuildPaul2013TokenBoundaryMarkerTransitions(
		initial,
		[]byte{8, 7},
		[]int16{12, 0},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []Paul2013TokenBoundaryMarker{
		{Marker: ']', PreviousStateFlag: true},
		{Marker: 'Z', PreviousStateFlag: false},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("marker transitions = %+v, want %+v", got, want)
	}
}

func TestBuildPaul2013TokenBoundaryStateResultWithTransitions(t *testing.T) {
	got, err := BuildPaul2013TokenBoundaryStateResultWithTransitions(
		[]byte{'Z', 'Z'},
		[]int32{-1, 0},
		[]int32{0, 0},
		[]int32{-2},
		[]byte{7, 8},
		[]int16{0, 0},
	)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int32{100, 0}; !reflect.DeepEqual(got.Values, want) {
		t.Fatalf("state values = %v, want %v", got.Values, want)
	}
	if want := []int32{2, 0}; !reflect.DeepEqual(got.States, want) {
		t.Fatalf("token states = %v, want %v", got.States, want)
	}
	wantMarkers := []Paul2013TokenBoundaryMarker{
		{Marker: '[', PreviousStateFlag: false},
		{Marker: '[', PreviousStateFlag: false},
	}
	if !reflect.DeepEqual(got.Markers, wantMarkers) {
		t.Fatalf("token markers = %+v, want %+v", got.Markers, wantMarkers)
	}
}
