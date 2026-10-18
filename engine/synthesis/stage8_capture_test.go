package synthesis

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/dat"
	"vtspeak/engine/selection"
	"vtspeak/engine/voice"
)

// This opt-in local-model check exercises the no-context boundary case from
// the Stage 8 Wine capture. It skips when proprietary model inputs are absent.
func TestStage8ShortNoContextJoinMatchesRuntimeCapture(t *testing.T) {
	dataRoot := os.Getenv("VTSPEAK_PAUL2013_DATA_ROOT")
	if dataRoot == "" {
		dataRoot = filepath.Join("..", "..", "data-paul", "M16")
		if _, err := os.Stat(dataRoot); err != nil {
			t.Skip("set VTSPEAK_PAUL2013_DATA_ROOT to run the local-model parity check")
		}
	}

	model, err := voice.OpenPaul2013(dataRoot)
	if err != nil {
		t.Fatalf("open local Paul 2013 model: %v", err)
	}
	defer model.Close()

	refs := []selection.UnitRef{
		{Bank: "gen", Index: 210603},
		{Bank: "gen", Index: 210604},
	}
	controls := Paul2013TimelineControlWords{Primary: 100, Speed: 100, Gain: 200}
	rows := make([]Paul2013NormalTimelineRow, len(refs))
	for index, ref := range refs {
		unit, err := model.ReadUnit(ref.Bank, ref.Index)
		if err != nil {
			t.Fatalf("read captured unit %d: %v", index, err)
		}
		rows[index], err = BuildPaul2013NormalTimelineRow(
			unit, ref, 0, Paul2013TimelineCombinedView, int16(index), controls,
		)
		if err != nil {
			t.Fatalf("build captured timeline row %d: %v", index, err)
		}
	}
	pcm, err := RenderPaul2013NormalTimelineRows(context.Background(), model, rows)
	if err != nil {
		t.Fatalf("render captured normal rows: %v", err)
	}
	got, err := dat.WAVFromSamples(pcm.Samples)
	if err != nil {
		t.Fatalf("encode rendered PCM: %v", err)
	}
	// Pin the complete original WAV, including its header, without retaining a
	// generated audio file. The reproduction manifest records the same digest.
	const want = "473d40d76100c8ec48f125407b9d96d2a896b3728fc6b74aea504ea784f76360"
	if digest := fmt.Sprintf("%x", sha256.Sum256(got)); digest != want {
		t.Fatalf("default no-context join differs from Stage 8 runtime WAV: SHA-256 %s, want %s", digest, want)
	}
}
