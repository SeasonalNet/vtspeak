package synthesis

// NormalizePaul2013Controls resolves negative API sentinels to the observed
// Paul defaults and clamps explicit values to the ranges used by the engine.
// Speed zero follows the API's special case and becomes 50 before clamping.
// The returned values are effective control words, not the caller's inputs.
func NormalizePaul2013Controls(controls Controls) Controls {
	return Controls{
		Pitch:  normalizePaul2013Control(controls.Pitch, 100, 50, 200, 0),
		Speed:  normalizePaul2013Control(controls.Speed, 100, 50, 400, 50),
		Volume: normalizePaul2013Control(controls.Volume, 200, 0, 500, 0),
	}
}

func normalizePaul2013Control(value, defaultValue, minimum, maximum, zeroValue int32) int32 {
	if value < 0 {
		return defaultValue
	}
	if value == 0 && zeroValue != 0 {
		value = zeroValue
	}
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
