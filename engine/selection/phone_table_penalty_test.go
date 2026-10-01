package selection

import "testing"

func TestCrossBytePenaltyUsesNativeVowelAndStressWindows(t *testing.T) {
	for _, test := range []struct {
		name               string
		candidate, context [7]byte
		want               float32
	}{
		{"unstressed preceding vowel skips special branch", [7]byte{0, 2, 0x36}, [7]byte{0, 1, 0x36}, 0},
		{"stressed preceding vowel adds signature difference", [7]byte{0, 1, 0x36}, [7]byte{0, 2, 0x36}, 100},
		{"central vowel stress is same category", [7]byte{0, 0, 2}, [7]byte{0, 0, 1}, 0},
		{"central vowel consonant category difference", [7]byte{0, 0, 0x13}, [7]byte{0, 0, 1}, 100},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := UnitCrossBytePenalty(test.candidate, test.context)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("penalty %v want %v", got, test.want)
			}
		})
	}
}
