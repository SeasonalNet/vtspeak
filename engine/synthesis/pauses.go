package synthesis

import (
	"fmt"

	"vtspeak/engine/text"
)

const paul2013OutputSampleRate = 16000

// InsertPaul2013InlinePauses inserts the observed zero-valued PCM frames for
// inline VTML pauses. boundaryFrameOffsets contains the frame offset in the
// original PCM after each token boundary, including the leading boundary at
// index 0 and trailing boundary at index len(boundaryFrameOffsets)-1. Mapping
// lexical token boundaries to rendered frame offsets remains caller work.
func InsertPaul2013InlinePauses(
	pcm PCM,
	boundaryFrameOffsets []uint64,
	pauses []text.Paul2013InlinePause,
) (PCM, error) {
	if pcm.SampleRate != paul2013OutputSampleRate {
		return PCM{}, fmt.Errorf("Paul 2013 inline pauses require %d Hz PCM, got %d Hz", paul2013OutputSampleRate, pcm.SampleRate)
	}
	if len(boundaryFrameOffsets) == 0 {
		return PCM{}, fmt.Errorf("inline pause insertion requires at least the leading token boundary")
	}
	if len(pauses) > 0 && len(boundaryFrameOffsets) < 2 {
		return PCM{}, fmt.Errorf("inline pauses require a boundary table for at least one token")
	}
	previousOffset := uint64(0)
	for index, offset := range boundaryFrameOffsets {
		if offset < previousOffset || offset > uint64(len(pcm.Samples)) {
			return PCM{}, fmt.Errorf("token boundary %d frame offset %d is not nondecreasing within %d frames", index, offset, len(pcm.Samples))
		}
		previousOffset = offset
	}
	maxInt := uint64(^uint(0) >> 1)
	outputLength := uint64(len(pcm.Samples))
	previousTokenBoundary := 0
	for index, pause := range pauses {
		if pause.AfterToken < previousTokenBoundary || pause.AfterToken >= len(boundaryFrameOffsets) {
			return PCM{}, fmt.Errorf("inline pause %d token boundary %d is outside the ordered boundary table", index, pause.AfterToken)
		}
		if pause.OutputFramesAt16KHz != uint64(pause.DurationMilliseconds)*16 {
			return PCM{}, fmt.Errorf("inline pause %d frame count does not match its millisecond duration", index)
		}
		if pause.OutputFramesAt16KHz > maxInt-outputLength {
			return PCM{}, fmt.Errorf("inline pause %d would exceed the platform slice capacity", index)
		}
		outputLength += pause.OutputFramesAt16KHz
		previousTokenBoundary = pause.AfterToken
	}
	if outputLength > maxInt {
		return PCM{}, fmt.Errorf("inline-pause PCM length %d exceeds the platform slice capacity", outputLength)
	}

	result := make([]int16, int(outputLength))
	inputOffset := uint64(0)
	outputOffset := 0
	for _, pause := range pauses {
		boundary := boundaryFrameOffsets[pause.AfterToken]
		copyCount := int(boundary - inputOffset)
		copy(result[outputOffset:outputOffset+copyCount], pcm.Samples[int(inputOffset):int(boundary)])
		outputOffset += copyCount + int(pause.OutputFramesAt16KHz)
		inputOffset = boundary
	}
	copy(result[outputOffset:], pcm.Samples[int(inputOffset):])
	return PCM{SampleRate: pcm.SampleRate, Samples: result}, nil
}

// InsertPaul2013InlinePausesAfterTokenChunks derives frame offsets from
// rendered token chunks, then inserts the captured inline pause events. Each
// chunk must correspond to one resolved lexical token; phoneme or unit chunks
// need a separate token-to-frame mapping before calling this helper.
func InsertPaul2013InlinePausesAfterTokenChunks(
	tokenChunks []PCM,
	pauses []text.Paul2013InlinePause,
) (PCM, error) {
	if len(tokenChunks) == 0 {
		return PCM{}, fmt.Errorf("inline pause insertion requires at least one rendered token chunk")
	}
	sampleRate := tokenChunks[0].SampleRate
	if sampleRate != paul2013OutputSampleRate {
		return PCM{}, fmt.Errorf("Paul 2013 inline pauses require %d Hz PCM, got %d Hz", paul2013OutputSampleRate, sampleRate)
	}
	maxInt := int(^uint(0) >> 1)
	totalFrames := 0
	for index, chunk := range tokenChunks {
		if chunk.SampleRate != sampleRate {
			return PCM{}, fmt.Errorf("token chunk %d sample rate %d differs from first chunk rate %d", index, chunk.SampleRate, sampleRate)
		}
		if len(chunk.Samples) > maxInt-totalFrames {
			return PCM{}, fmt.Errorf("token PCM length exceeds the platform slice capacity")
		}
		totalFrames += len(chunk.Samples)
	}
	combined := make([]int16, 0, totalFrames)
	boundaryOffsets := make([]uint64, len(tokenChunks)+1)
	for index, chunk := range tokenChunks {
		combined = append(combined, chunk.Samples...)
		boundaryOffsets[index+1] = uint64(len(combined))
	}
	return InsertPaul2013InlinePauses(
		PCM{SampleRate: sampleRate, Samples: combined}, boundaryOffsets, pauses,
	)
}
