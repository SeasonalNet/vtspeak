package selection

import (
	"reflect"
	"testing"
)

func TestBuildPaul2013FallbackRows(t *testing.T) {
	base := [6]byte{3, 4, 5, 0xfe, 8, 9}
	target := [7]byte{10, 11, 12, 13, 14, 15, 0xff}
	rows := BuildPaul2013FallbackRows(base, 2, 0, target)
	wantBytes := [2][6]byte{
		{3, 4, 5, 2, 1, 9},
		{3, 4, 5, 3, 2, 9},
	}
	for index, row := range rows {
		if row.Bytes != wantBytes[index] {
			t.Errorf("row %d bytes = %v, want %v", index, row.Bytes, wantBytes[index])
		}
		if len(row.TreeTargets) != 1 || row.TreeTargets[0] != (Paul2013FallbackTreeTarget{Signature: target, QueryMode: 1}) {
			t.Errorf("row %d tree targets = %#v, want original target in mode 1", index, row.TreeTargets)
		}
	}
}

func TestBuildPaul2013FallbackRowsSpecialModelClass(t *testing.T) {
	base := [6]byte{1, 2, 3, 4, 5, 6}
	target := [7]byte{10, 11, 12, 13, 14, 15, 0x65}
	rows := BuildPaul2013FallbackRows(base, -3, 12, target)
	if rows[0].Bytes[3] != 0xfd || rows[1].Bytes[3] != 0xfe {
		t.Fatalf("fallback byte +3 = %d,%d, want byte-wrapped signed delta and next ordinal", rows[0].Bytes[3], rows[1].Bytes[3])
	}
	if rows[0].Bytes[4] != 1 || rows[1].Bytes[4] != 2 {
		t.Fatalf("fallback byte +4 = %d,%d, want row classes 1,2", rows[0].Bytes[4], rows[1].Bytes[4])
	}
	for index, row := range rows {
		if len(row.TreeTargets) != 2 {
			t.Fatalf("row %d has %d tree targets, want 2", index, len(row.TreeTargets))
		}
		masked := target
		masked[6] &= 0x80
		if row.TreeTargets[0] != (Paul2013FallbackTreeTarget{Signature: masked, QueryMode: 1}) {
			t.Errorf("row %d first target = %#v, want masked signature in mode 1", index, row.TreeTargets[0])
		}
		if row.TreeTargets[1] != (Paul2013FallbackTreeTarget{Signature: target, QueryMode: 2}) {
			t.Errorf("row %d second target = %#v, want original signature in mode 2", index, row.TreeTargets[1])
		}
	}
	rows[0].TreeTargets[0].Signature[0] = 0
	if rows[1].TreeTargets[0].Signature[0] != target[0] {
		t.Fatal("fallback rows share mutable tree-target slices")
	}
}

func TestMergePaul2013FallbackCandidates(t *testing.T) {
	got, err := MergePaul2013FallbackCandidates(
		[]uint32{8, 2, 5},
		[]uint32{7, 5, 1, 7, 2},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013FallbackCandidateList{IDs: []uint32{8, 2, 5, 1, 7}, RankedCount: 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged fallback candidates = %#v, want %#v", got, want)
	}
}

func TestMergePaul2013FallbackCandidatesRejectsNativeBounds(t *testing.T) {
	if _, err := MergePaul2013FallbackCandidates([]uint32{1}, []uint32{1 << 16}); err == nil {
		t.Fatal("out-of-range class ID was accepted")
	}
	tooMany := make([]uint32, paul2013FallbackCandidateCapacity+1)
	for index := range tooMany {
		tooMany[index] = uint32(index)
	}
	if _, err := MergePaul2013FallbackCandidates(tooMany, nil); err == nil {
		t.Fatal("ranked list over native capacity was accepted")
	}
	if _, err := MergePaul2013FallbackCandidates(nil, tooMany); err == nil {
		t.Fatal("candidate union over native capacity was accepted")
	}
}
