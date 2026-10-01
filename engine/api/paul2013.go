// Package api contains evidence-backed behavior at the legacy VoiceText
// exported API boundary that is independent of text synthesis internals.
package api

var paul2013SpeakerNames = [...]string{"Kate", "Paul", "em001", "Julie", "James", "Ashley"}

var paul2013SpeakerMetadata = [...]struct {
	name string
	path string
}{
	{name: "kate", path: "d:/eng/db/susan/pcm/"},
	{name: "paul", path: "d:/eng/db/isaac/pcm/"},
	{name: "em001", path: "d:/eng/db/em001/pcm/"},
	{name: "julie", path: "d:/eng/db/jennifer/pcm/"},
	{name: "james", path: "d:/eng/db/lee/pcm/"},
	{name: "ashley", path: "d:/eng/db/casey/pcm/"},
}

// Paul2013SpeakerSettings contains the mutable per-speaker values exposed by
// the legacy configuration setters and getters. Loaded is supplied by the
// caller because the API's native speaker slots are process state.
type Paul2013SpeakerSettings struct {
	Loaded        bool
	Pitch         int32
	Speed         int32
	Volume        int32
	SentencePause int32
	CommaPause    int32
	Emphasis      int32
}

// Paul2013PitchSpeedVolumePauseUpdates mirrors the four setter arguments.
// Negative values mean that the corresponding native field is left alone.
type Paul2013PitchSpeedVolumePauseUpdates struct {
	Pitch         int32
	Speed         int32
	Volume        int32
	SentencePause int32
}

// Paul2013PitchSpeedVolumePauseOutputMask identifies the non-null output
// pointers passed to VT_GetPitchSpeedVolumePause_ENG.
type Paul2013PitchSpeedVolumePauseOutputMask struct {
	Pitch         bool
	Speed         bool
	Volume        bool
	SentencePause bool
}

// Paul2013PitchSpeedVolumePauseResult preserves the getter return code,
// values visible to the caller, and which output pointers were written.
type Paul2013PitchSpeedVolumePauseResult struct {
	ReturnCode int32
	Values     Paul2013PitchSpeedVolumePauseUpdates
	Wrote      Paul2013PitchSpeedVolumePauseOutputMask
}

// Paul2013CommaPauseResult preserves the comma-pause getter return code and
// its conditional output write.
type Paul2013CommaPauseResult struct {
	ReturnCode int32
	Value      int32
	WroteValue bool
}

// Paul2013GlobalSettings contains configuration fields stored on the
// initialized process-wide engine object rather than in a speaker slot.
type Paul2013GlobalSettings struct {
	Initialized          bool
	Highlight            byte
	ParenthesisCharCount int32
	EnglishReadingRule   int32
}

// Paul2013SpeakerName returns the name selected by VT_GetSpeakerName_ENG.
// Selectors outside 0 through 5 return the compiled Paul fallback.
func Paul2013SpeakerName(selector int32) string {
	if selector < 0 || selector >= int32(len(paul2013SpeakerNames)) {
		return "Paul"
	}
	return paul2013SpeakerNames[selector]
}

// Paul2013SpeakerInfo returns the fixed lowercase name and database path
// copied by VT_SpeakersInfo_ENG for a valid slot. The native API returns 6
// after copying both NUL-terminated strings and does not expose buffer sizes.
func Paul2013SpeakerInfo(selector int32) (name, path string, ok bool) {
	if selector < 0 || selector >= int32(len(paul2013SpeakerMetadata)) {
		return "", "", false
	}
	metadata := paul2013SpeakerMetadata[selector]
	return metadata.name, metadata.path, true
}

// Paul2013SetGlobalHighlight applies the process-object gate and boolean byte
// normalization used by VT_SetTextTypeForHighlight_ENG.
func Paul2013SetGlobalHighlight(current Paul2013GlobalSettings, value byte) Paul2013GlobalSettings {
	if current.Initialized {
		current.Highlight = Paul2013HighlightSetting(value)
	}
	return current
}

// Paul2013SetParenthesisCharCount applies the initialized-object gate and
// negative-to-zero clamp used by VT_SetParenthesisCharNumber_ENG.
func Paul2013SetParenthesisCharCount(current Paul2013GlobalSettings, value int32) Paul2013GlobalSettings {
	if current.Initialized {
		current.ParenthesisCharCount = max(value, 0)
	}
	return current
}

// Paul2013SetEnglishReadingRule applies the initialized-object gate and
// negative-to-zero clamp used by VT_SetEnglishReadingRule_KOR.
func Paul2013SetEnglishReadingRule(current Paul2013GlobalSettings, value int32) Paul2013GlobalSettings {
	if current.Initialized {
		current.EnglishReadingRule = max(value, 0)
	}
	return current
}

func paul2013SpeakerSlot(selector int32) int {
	if selector < 0 || selector >= 6 {
		return 1
	}
	return int(selector)
}

// Paul2013SetPitchSpeedVolumePause applies the native setter's loaded-slot
// gate, negative-value no-op, and per-field clamps. Invalid selectors use
// slot 1. The input array represents current process speaker state.
func Paul2013SetPitchSpeedVolumePause(
	selector int32,
	speakers [6]Paul2013SpeakerSettings,
	updates Paul2013PitchSpeedVolumePauseUpdates,
) [6]Paul2013SpeakerSettings {
	slot := paul2013SpeakerSlot(selector)
	if !speakers[slot].Loaded {
		return speakers
	}
	current := &speakers[slot]
	if updates.Pitch >= 0 {
		current.Pitch = clampPaul2013(updates.Pitch, 50, 200)
	}
	if updates.Speed >= 0 {
		current.Speed = clampPaul2013(updates.Speed, 50, 400)
	}
	if updates.Volume >= 0 {
		current.Volume = clampPaul2013(updates.Volume, 0, 500)
	}
	if updates.SentencePause >= 0 {
		current.SentencePause = clampPaul2013(updates.SentencePause, 0, 65535)
	}
	return speakers
}

// Paul2013GetPitchSpeedVolumePause applies the native selector and loaded-slot
// checks and writes only requested outputs. On failure, priorOutput is
// returned unchanged, matching the native output-pointer behavior.
func Paul2013GetPitchSpeedVolumePause(
	selector int32,
	speakers [6]Paul2013SpeakerSettings,
	priorOutput Paul2013PitchSpeedVolumePauseUpdates,
	outputs Paul2013PitchSpeedVolumePauseOutputMask,
) Paul2013PitchSpeedVolumePauseResult {
	result := Paul2013PitchSpeedVolumePauseResult{
		ReturnCode: -1,
		Values:     priorOutput,
	}
	slot := paul2013SpeakerSlot(selector)
	if !speakers[slot].Loaded {
		return result
	}
	current := speakers[slot]
	result.ReturnCode = 1
	if outputs.Pitch {
		result.Values.Pitch = current.Pitch
		result.Wrote.Pitch = true
	}
	if outputs.Speed {
		result.Values.Speed = current.Speed
		result.Wrote.Speed = true
	}
	if outputs.Volume {
		result.Values.Volume = current.Volume
		result.Wrote.Volume = true
	}
	if outputs.SentencePause {
		result.Values.SentencePause = current.SentencePause
		result.Wrote.SentencePause = true
	}
	return result
}

// Paul2013SetCommaPause applies VT_SetCommaPause_ENG's loaded-slot gate,
// negative-value no-op, and unsigned-short maximum.
func Paul2013SetCommaPause(selector int32, speakers [6]Paul2013SpeakerSettings, value int32) [6]Paul2013SpeakerSettings {
	slot := paul2013SpeakerSlot(selector)
	if !speakers[slot].Loaded || value < 0 {
		return speakers
	}
	speakers[slot].CommaPause = clampPaul2013(value, 0, 65535)
	return speakers
}

// Paul2013GetCommaPause applies VT_GetCommaPause_ENG's selector and loaded
// slot checks. A null output pointer is represented by wantOutput=false.
func Paul2013GetCommaPause(
	selector int32,
	speakers [6]Paul2013SpeakerSettings,
	priorOutput int32,
	wantOutput bool,
) Paul2013CommaPauseResult {
	result := Paul2013CommaPauseResult{ReturnCode: -1, Value: priorOutput}
	slot := paul2013SpeakerSlot(selector)
	if !speakers[slot].Loaded {
		return result
	}
	result.ReturnCode = 1
	if wantOutput {
		result.Value = speakers[slot].CommaPause
		result.WroteValue = true
	}
	return result
}

// Paul2013SetEmphasisFactor applies VT_SetEmphasisFactor_ENG's loaded-slot
// gate and inclusive -95 through 95 clamp.
func Paul2013SetEmphasisFactor(selector int32, speakers [6]Paul2013SpeakerSettings, value int32) [6]Paul2013SpeakerSettings {
	slot := paul2013SpeakerSlot(selector)
	if !speakers[slot].Loaded {
		return speakers
	}
	speakers[slot].Emphasis = clampPaul2013(value, -95, 95)
	return speakers
}

func clampPaul2013(value, minimum, maximum int32) int32 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

// Paul2013DBSizeResult retains VT_GetDBSize_ENG's return code and its
// conditional output write. SizeBytes remains the caller's prior value when
// the selected slot is not loaded.
type Paul2013DBSizeResult struct {
	ReturnCode int32
	SizeBytes  uint32
	WroteSize  bool
}

// Paul2013DBSizeQuery applies VT_GetDBSize_ENG's selector and loaded-slot
// contract to caller-supplied current database state. Invalid selectors map
// to slot 1. The native API returns -1 without writing the output for an
// unloaded slot; a loaded slot returns 1 and writes its size.
func Paul2013DBSizeQuery(
	selector int32,
	loaded [6]bool,
	sizesBytes [6]uint32,
	priorOutput uint32,
) Paul2013DBSizeResult {
	if selector < 0 || selector >= int32(len(loaded)) {
		selector = 1
	}
	if !loaded[selector] {
		return Paul2013DBSizeResult{ReturnCode: -1, SizeBytes: priorOutput}
	}
	return Paul2013DBSizeResult{
		ReturnCode: 1,
		SizeBytes:  sizesBytes[selector],
		WroteSize:  true,
	}
}

// Paul2013UserDictionaryLimit returns the limit selected by
// VT_GetUserDictLimit_ENG. The native switch is defined only for selectors
// zero through four; every other signed value returns -1.
func Paul2013UserDictionaryLimit(selector int32) int32 {
	switch selector {
	case 0:
		return 30
	case 1:
		return 10
	case 2:
		return 50
	case 3, 4:
		return 65
	default:
		return -1
	}
}

// Paul2013HighlightSetting normalizes the byte accepted by
// VT_SetTextTypeForHighlight_ENG. Stage 21 measured zero as disabled and all
// nonzero byte values as enabled. The setting's effect on auxiliary
// highlight/index output is not modeled here.
func Paul2013HighlightSetting(value byte) byte {
	if value == 0 {
		return 0
	}
	return 1
}

// Paul2013UnitHistorySetting applies VT_SetUnitSelectHistoryMode_ENG's
// observed pre-load gate. Before loading, only input byte 1 stores 1; all
// other byte values store 0. After loading, the setter leaves the prior value
// unchanged. The broader synthesis effect of this mode remains unresolved.
func Paul2013UnitHistorySetting(current byte, loaded bool, value byte) byte {
	if loaded {
		return current
	}
	if value == 1 {
		return 1
	}
	return 0
}
