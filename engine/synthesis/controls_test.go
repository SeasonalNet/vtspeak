package synthesis

import "testing"

func TestResolvePaul2013ControlOverridesSelectsIndependentlyAndClamps(t *testing.T) {
	got := ResolvePaul2013ControlOverrides(Paul2013ControlOverrides{
		Defaults: Controls{Pitch: 180, Speed: 240, Volume: 300},
		Values:   Controls{Pitch: 25, Speed: 450, Volume: -1},
		Pitch:    true,
		Speed:    true,
	})
	want := Controls{Pitch: 50, Speed: 400, Volume: 300}
	if got != want {
		t.Fatalf("resolved controls = %+v, want %+v", got, want)
	}
}

func TestResolvePaul2013ControlOverridesClampsSelectedValues(t *testing.T) {
	got := ResolvePaul2013ControlOverrides(Paul2013ControlOverrides{
		Defaults: Controls{Pitch: -1, Speed: -1, Volume: -1},
		Values:   Controls{Pitch: -1, Speed: 0, Volume: 700},
		Pitch:    true,
		Speed:    true,
		Volume:   true,
	})
	want := Controls{Pitch: 50, Speed: 50, Volume: 500}
	if got != want {
		t.Fatalf("resolved controls = %+v, want %+v", got, want)
	}
}
