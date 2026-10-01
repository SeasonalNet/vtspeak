package synthesis

// Paul2013ControlOverrides records the per-control override switches used by
// FUN_10022850. An enabled switch selects the matching value from Values;
// otherwise the corresponding value from Defaults is used.
type Paul2013ControlOverrides struct {
	Defaults Controls
	Values   Controls
	Pitch    bool
	Speed    bool
	Volume   bool
}

// ResolvePaul2013ControlOverrides ports FUN_10022850's per-control selection
// and clamp behavior. Overrides remain
// independent: enabling one does not change the source of either other
// control. Defaults are the already-selected engine-level words. This helper
// mirrors FUN_10022850's clamps; API sentinel handling remains in
// NormalizePaul2013Controls.
func ResolvePaul2013ControlOverrides(input Paul2013ControlOverrides) Controls {
	selected := input.Defaults
	if input.Pitch {
		selected.Pitch = input.Values.Pitch
	}
	if input.Speed {
		selected.Speed = input.Values.Speed
	}
	if input.Volume {
		selected.Volume = input.Values.Volume
	}
	return clampPaul2013Controls(selected)
}

// NormalizePaul2013Controls resolves negative API sentinels to the observed
// Paul defaults and clamps explicit values to the ranges used by the engine.
// Speed zero follows the API's special case and becomes 50 before clamping.
// The returned values are effective control words, not the caller's inputs.
func NormalizePaul2013Controls(controls Controls) Controls {
	resolved := Controls{
		Pitch:  normalizePaul2013Control(controls.Pitch, 100, 50, 200, 0),
		Speed:  normalizePaul2013Control(controls.Speed, 100, 50, 400, 50),
		Volume: normalizePaul2013Control(controls.Volume, 200, 0, 500, 0),
	}
	return clampPaul2013Controls(resolved)
}

func clampPaul2013Controls(controls Controls) Controls {
	return Controls{
		Pitch:  clampPaul2013Control(controls.Pitch, 50, 200),
		Speed:  clampPaul2013Control(controls.Speed, 50, 400),
		Volume: clampPaul2013Control(controls.Volume, 0, 500),
	}
}

func normalizePaul2013Control(value, defaultValue, minimum, maximum, zeroValue int32) int32 {
	if value < 0 {
		return defaultValue
	}
	if value == 0 && zeroValue != 0 {
		value = zeroValue
	}
	return clampPaul2013Control(value, minimum, maximum)
}

func clampPaul2013Control(value, minimum, maximum int32) int32 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
