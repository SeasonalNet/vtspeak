package text

import "testing"

func TestPaul2013FeatureViews(t *testing.T) {
	key := [5]byte{0x21, 0x34, 0x56, 0x95, 0x20}
	mode1, err := Paul2013FeatureViewForKey(key, 1)
	if err != nil {
		t.Fatal(err)
	}
	flags := key[3]
	upperBoundary := boolByte(flags&0x80 != 0 && (flags>>3)&7 >= 2)
	lowerBoundary := boolByte(flags&0x40 != 0 && flags&7 >= 2)
	wantMode1 := [10]byte{
		0, key[4], key[1], paul2013ClassMap[key[0]], upperBoundary,
		paul2013PrimaryCategoryMap[key[0]], ((flags>>3)&7)*10 + lowerBoundary,
		key[0], (flags >> 3) & 7, paul2013ClassMap[key[2]],
	}
	if mode1 != wantMode1 {
		t.Fatalf("mode 1 view = % x, want % x", mode1, wantMode1)
	}

	mode2, err := Paul2013FeatureViewForKey(key, 2)
	if err != nil {
		t.Fatal(err)
	}
	wantMode2 := [10]byte{
		1, key[4], key[1], paul2013ClassMap[key[2]], lowerBoundary,
		paul2013SecondaryCategoryMap[key[2]], (flags&7)*10 + upperBoundary,
		key[2], flags & 7, paul2013ClassMap[key[0]],
	}
	if mode2 != wantMode2 {
		t.Fatalf("mode 2 view = % x, want % x", mode2, wantMode2)
	}
	if _, err := Paul2013FeatureViewForKey(key, 3); err == nil {
		t.Fatal("unsupported feature view accepted")
	}
}
