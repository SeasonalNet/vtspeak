package text

import "testing"

func TestParsePaul2013PronunciationPathGroupsAndMapsMarkers(t *testing.T) {
	got, err := ParsePaul2013PronunciationPath([]byte{'A', 'd', ')', 'x', ')', '(', '(', 'd'})
	if err != nil {
		t.Fatal(err)
	}
	want := [][]byte{{'A'}, {'%', 'x', '&'}, {}}
	if len(got) != len(want) {
		t.Fatalf("group count = %d, want %d: %q", len(got), len(want), got)
	}
	for groupIndex := range want {
		if string(got[groupIndex]) != string(want[groupIndex]) {
			t.Errorf("group %d = %q, want %q", groupIndex, got[groupIndex], want[groupIndex])
		}
	}
}

func TestSelectPaul2013PronunciationByPathMarker(t *testing.T) {
	tests := []struct {
		name      string
		paths     [][]byte
		phoneRows [][]byte
		selector  byte
		want      []byte
		present   bool
	}{
		{name: "direct selector beats earlier fallback", paths: [][]byte{{0x17, 'a'}, {'b'}, {'x', 'b'}}, phoneRows: [][]byte{{'p', 'a'}, {'p', 'b'}, {'p', 'c'}}, selector: 'b', want: []byte{'p', 'b'}, present: true},
		{name: "first direct selector wins", paths: [][]byte{{'x', 'b'}, {'b'}}, phoneRows: [][]byte{{'p', 'a'}, {'p', 'b'}}, selector: 'b', want: []byte{'p', 'a'}, present: true},
		{name: "close marker matches percent row", paths: [][]byte{{'x', '%'}, {'x', ')'}}, phoneRows: [][]byte{{'p', 'a'}, {'p', 'b'}}, selector: ')', want: []byte{'p', 'a'}, present: true},
		{name: "first special byte fallback", paths: [][]byte{{'a'}, {'b', 0x17}, {'c', 0x17}}, phoneRows: [][]byte{{'p', 'a'}, {'p', 'b'}, {'p', 'c'}}, selector: 'z', want: []byte{'p', 'b'}, present: true},
		{name: "terminator ends selector scan", paths: [][]byte{{'a', 0xff, 'z'}, {'b'}}, phoneRows: [][]byte{{'p', 'a'}, {'p', 'b'}}, selector: 'z', want: []byte{'p', 'a'}, present: true},
		{name: "phone row is nul terminated", paths: [][]byte{{'a'}}, phoneRows: [][]byte{{'p', 'a', 0, 'x'}}, selector: 'a', want: []byte{'p', 'a'}, present: true},
		{name: "default first row", paths: [][]byte{{'a'}, {'b'}}, phoneRows: [][]byte{{'p', 'a'}, {'p', 'b'}}, selector: 'z', want: []byte{'p', 'a'}, present: true},
		{name: "empty rows", selector: 'z'},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := make([][]byte, len(test.phoneRows))
			for index := range test.paths {
				original[index] = append([]byte(nil), test.phoneRows[index]...)
			}
			got, present, err := SelectPaul2013PronunciationByPathMarker(test.paths, test.phoneRows, test.selector)
			if err != nil {
				t.Fatal(err)
			}
			if present != test.present || string(got) != string(test.want) {
				t.Fatalf("selected path = (%q, %t), want (%q, %t)", got, present, test.want, test.present)
			}
			if present && len(test.paths) != 0 {
				got[0] = 'X'
				for index := range test.phoneRows {
					if string(test.phoneRows[index]) != string(original[index]) {
						t.Fatal("selected phone row aliases the caller's phone rows")
					}
				}
			}
		})
	}
	if _, _, err := SelectPaul2013PronunciationByPathMarker([][]byte{{'a'}}, nil, 'a'); err == nil {
		t.Fatal("mismatched path and phone rows were accepted")
	}
}

func TestPhonePayloadSelectPaul2013PronunciationByPathMarker(t *testing.T) {
	var codebook PhoneIDCodebook
	codebook[1] = [5]byte{'A', 'H', '0', 0}
	codebook[2] = [5]byte{'B', 'I', 'Y', '1'}
	payload := PhonePayload{Pronunciations: []Pronunciation{
		{Path: []byte{'a'}, Phone: []byte{1}},
		{Path: []byte{'b'}, Phone: []byte{2}},
	}}
	got, found, err := payload.SelectPaul2013PronunciationByPathMarker('b', codebook)
	if err != nil || !found || string(got) != "BIY1" {
		t.Fatalf("selected pronunciation = (%q, %t, %v), want (BIY1, true, nil)", got, found, err)
	}
}

func TestParsePaul2013PronunciationPathTerminatorAndEmptyPath(t *testing.T) {
	for name, input := range map[string][]byte{
		"terminated": {'A', 'B', 0xff},
		"empty":      {},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := ParsePaul2013PronunciationPath(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("groups = %q, want one group", got)
			}
			want := input
			if name == "terminated" {
				want = input[:len(input)-1]
			}
			if string(got[0]) != string(want) {
				t.Fatalf("group = %q, want %q", got[0], want)
			}
		})
	}
}

func TestParsePaul2013PronunciationPathRejectsUnsafeInputs(t *testing.T) {
	if _, err := ParsePaul2013PronunciationPath([]byte{'a', 0xff, 'b'}); err == nil {
		t.Fatal("bytes following path terminator accepted")
	}
	tooLong := make([]byte, paul2013PathGroupCapacity)
	for index := range tooLong {
		tooLong[index] = 'a'
	}
	if _, err := ParsePaul2013PronunciationPath(tooLong); err == nil {
		t.Fatal("oversized path group accepted")
	}
}

func TestRankPaul2013PronunciationPathGroups(t *testing.T) {
	selected, scores, err := RankPaul2013PronunciationPathGroups(
		[][]byte{{13, 19}, {14}}, []int16{258, 1, 2, 99},
	)
	if err != nil {
		t.Fatal(err)
	}
	if selected != 1 || len(scores) != 2 || scores[0] != 1 || scores[1] != 2 {
		t.Fatalf("selection = %d with scores %v, want 1 with [1 2]", selected, scores)
	}

	selected, scores, err = RankPaul2013PronunciationPathGroups(
		[][]byte{{13}, {14}}, nil,
	)
	if err != nil || selected != 0 || scores[0] != 0 || scores[1] != 0 {
		t.Fatalf("zero-score tie = (%d, %v, %v), want (0, [0 0], nil)", selected, scores, err)
	}
}

func TestRankPaul2013PronunciationPathGroupsRequiresCandidates(t *testing.T) {
	if _, _, err := RankPaul2013PronunciationPathGroups(nil, nil); err == nil {
		t.Fatal("empty candidate list accepted")
	}
	if _, _, err := RankPaul2013PronunciationPathGroups([][]byte{{46}}, nil); err == nil {
		t.Fatal("path code beyond the recovered map accepted")
	}
}

func TestPaul2013PronunciationPathClassUsesFirstCode(t *testing.T) {
	if len(paul2013PathCodeClasses) != 46 {
		t.Fatalf("pronunciation path class table has %d entries, want 46", len(paul2013PathCodeClasses))
	}
	for _, test := range []struct {
		pathCode byte
		want     int16
	}{
		{pathCode: 13, want: 1},
		{pathCode: 28, want: 8},
		{pathCode: 38, want: 10},
		{pathCode: 42, want: 12},
		{pathCode: 45, want: -1},
	} {
		got, present, err := Paul2013PronunciationPathClass([]byte{test.pathCode, 12})
		if err != nil || !present || got != test.want {
			t.Errorf("path class(%d) = (%d, %t, %v), want (%d, true, nil)", test.pathCode, got, present, err, test.want)
		}
	}
	if _, present, err := Paul2013PronunciationPathClass(nil); err != nil || present {
		t.Fatalf("empty path class = (present %t, error %v), want false, nil", present, err)
	}
	if _, _, err := Paul2013PronunciationPathClass([]byte{0xff}); err == nil {
		t.Fatal("out-of-range path class accepted")
	}
}
