package text

import "testing"

func TestPaul2013NormalizerTreeIndexMapsASCIICharacters(t *testing.T) {
	for _, test := range []struct {
		value byte
		want  int
	}{
		{'A', 0}, {'Z', 25}, {'a', 0}, {'z', 25}, {'\'', 26},
	} {
		got, err := Paul2013NormalizerTreeIndex(test.value)
		if err != nil {
			t.Fatalf("tree index for %q: %v", test.value, err)
		}
		if got != test.want {
			t.Errorf("tree index for %q = %d, want %d", test.value, got, test.want)
		}
	}
	if _, err := Paul2013NormalizerTreeIndex(' '); err == nil {
		t.Fatal("unsupported punctuation selected a normalizer tree")
	}
}

func TestBuildPaul2013NormalizerFeatureWindowMapsCharactersAndCategories(t *testing.T) {
	got, err := BuildPaul2013NormalizerFeatureWindow(
		[]byte("abcdefg"), []byte{0x01, 0xfe, 0x03, 0x04, 0x05}, 3,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013NormalizerFeatureWindow{1, 2, 3, 4, 5, 6, 7, 3, -2, 1}
	if got != want {
		t.Fatalf("normalizer feature window = %v, want %v", got, want)
	}
}

func TestBuildPaul2013NormalizerFeatureWindowZeroFillsSourceEdges(t *testing.T) {
	got, err := BuildPaul2013NormalizerFeatureWindow([]byte("A'"), []byte{5}, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := Paul2013NormalizerFeatureWindow{0, 0, 0, 1, 0x1c, 0, 0, 5}
	if got != want {
		t.Fatalf("edge normalizer feature window = %v, want %v", got, want)
	}
}

func TestBuildPaul2013NormalizerFeatureWindowRejectsMissingCategories(t *testing.T) {
	_, err := BuildPaul2013NormalizerFeatureWindow([]byte("abcd"), []byte{1}, 0)
	if err == nil {
		t.Fatal("missing required trailing category byte was accepted")
	}
}

func TestBuildPaul2013NormalizerFeatureWindowRejectsInvalidIndex(t *testing.T) {
	for _, index := range []int{-1, 2} {
		if _, err := BuildPaul2013NormalizerFeatureWindow([]byte("a"), nil, index); err == nil {
			t.Errorf("invalid source index %d was accepted", index)
		}
	}
}
