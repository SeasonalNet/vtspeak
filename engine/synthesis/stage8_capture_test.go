package synthesis

import (
	"bytes"
	"context"
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
	want, err := os.ReadFile(filepath.Join("..", "..", "tools", "revkit", "work", "stage8", "map-short.wav"))
	if err != nil {
		t.Fatalf("read Stage 8 runtime reference: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("default no-context join differs from Stage 8 runtime WAV: got %d bytes, want %d", len(got), len(want))
	}
}
