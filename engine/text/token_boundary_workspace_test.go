package text

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"vtspeak/engine/tree3"
)

func stateWorkspaceFixture() []byte {
	workspace := bytes.Repeat([]byte{0xa5}, 0x122500)
	binary.LittleEndian.PutUint32(workspace[0x1223f4:], 0)
	binary.LittleEndian.PutUint32(workspace[0x122404:], 0)
	for _, field := range []struct {
		offset int
		value  uint32
	}{{0x1210fc, 17}, {0x1210f4, 19}, {0x1210f8, 23}} {
		binary.LittleEndian.PutUint32(workspace[field.offset:], field.value)
	}
	return workspace
}

func TestPositionStateWorkspacePreservesShiftedValueTail(t *testing.T) {
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[0x121a68:], 0)
	before := append([]byte(nil), workspace...)
	got, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace, Paul2013PositionStateProgram{RowCount: 2, Initial: [3]int32{999, 999, 999}})
	if err != nil {
		t.Fatal(err)
	}
	for mode, offset := range paul2013PositionStateWorkspaceOffsets {
		if offset == 0 {
			continue
		}
		want := int32(-1)
		if mode < 3 {
			want = []int32{17, 19, 23}[mode]
		}
		for index := 0; index < 2; index++ {
			if value := int32(binary.LittleEndian.Uint32(got.Workspace[offset+index*4:])); value != want {
				t.Fatalf("mode %d row %d = %d want %d", mode, index, value, want)
			}
		}
	}
	if binary.LittleEndian.Uint32(got.Workspace[0x121a68:]) != 0 || !bytes.Equal(workspace, before) || !bytes.Equal(got.Workspace[0x1220a8:], before[0x1220a8:]) {
		t.Fatal("overwrote retained tail or caller workspace")
	}
}

func TestPositionStateWorkspaceWritesEventPrograms(t *testing.T) {
	workspace := stateWorkspaceFixture()
	got, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace, Paul2013PositionStateProgram{RowCount: 2, Events: map[byte]Paul2013PositionEventValues{
		3: {Mode: 3, Ranges: []Paul2013PositionEventRange{{Minimum: 0, Maximum: 9}, {Minimum: 10, Maximum: 19}}, Boundaries: []int32{5, 15}, Values: []int32{7, 12}, Minimum: 0, Maximum: 100},
		4: {Mode: 4, Ranges: []Paul2013PositionEventRange{{Minimum: 0, Maximum: 9}, {Minimum: 10, Maximum: 19}}, Boundaries: []int32{5, 15}, Values: []int32{0, 1}, Minimum: 0, Maximum: 3},
	}})
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []struct {
		offset int
		want   int32
	}{{0x121a60, 7}, {0x121a64, 12}, {0x121d80, 0}, {0x121d84, 1}} {
		if value := int32(binary.LittleEndian.Uint32(got.Workspace[field.offset:])); value != field.want {
			t.Fatalf("event output +%#x = %d want %d", field.offset, value, field.want)
		}
	}
	if binary.LittleEndian.Uint32(got.Workspace[0x1223f4:]) != 2 || binary.LittleEndian.Uint32(got.Workspace[0x122404:]) != 2 || int32(binary.LittleEndian.Uint32(got.Workspace[0x122440:])) != -1 {
		t.Fatal("event cursor/reset writeback missing")
	}
}

func TestPositionWorkspaceResumesNativeCursorAndWritesTerminalCursor(t *testing.T) {
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[0x1223f4:], 1)
	pass := Paul2013PositionEventValues{Mode: 3, Ranges: []Paul2013PositionEventRange{{Minimum: 0, Maximum: 3}}, Boundaries: []int32{0, 1, 2, 3, 4}, Values: []int32{100, 8, -2, 1, 7}, StartBoundary: 0, Minimum: 0, Maximum: 10}
	input := Paul2013PositionStateProgram{RowCount: 1, Events: map[byte]Paul2013PositionEventValues{3: pass}, Terminal: &Paul2013PositionTerminalAccumulator{Enabled: true, MinimumBoundary: 0, MaximumBoundary: 5, Boundaries: pass.Boundaries, Values: pass.Values, Minimum: 0, Maximum: 10}}
	got, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace, input)
	if err != nil {
		t.Fatal(err)
	}
	if got.Program.Arrays[3][0] != 7 || got.Program.Terminal.Value != 7 || binary.LittleEndian.Uint32(got.Workspace[0x1223f4:]) != 5 || binary.LittleEndian.Uint32(got.Workspace[0x122440:]) != 7 {
		t.Fatalf("resumed native cursor = %+v", got.Program)
	}
	if input.Events[3].StartBoundary != 0 || binary.LittleEndian.Uint32(workspace[0x1223f4:]) != 1 {
		t.Fatal("mutated supplied program/workspace")
	}
	binary.LittleEndian.PutUint32(workspace[0x1223f4:], 0xffffffff)
	if _, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace, input); err == nil {
		t.Fatal("accepted negative native cursor")
	}
}

func finalizedBoundaryWorkspaceFixture() Paul2013FinalizedModelState {
	parser := make([]byte, 2*0x94)
	binary.LittleEndian.PutUint32(parser[0x20:], 0xfffffffe)
	return Paul2013FinalizedModelState{Finalized: Paul2013ModelParserFinalizerResult{ReturnValue: 1, OutputGroupCount: 2, StateArena: parser}, Records: Paul2013ModelStateRecordArena{Bytes: markerArenaFixture(2)}}
}

func TestBoundaryWorkspaceUsesShiftedValuesAndWritesNativeArrays(t *testing.T) {
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[0x121a68:], 0)
	initialized, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace, Paul2013PositionStateProgram{RowCount: 2})
	if err != nil {
		t.Fatal(err)
	}
	before := append([]byte(nil), initialized.Workspace...)
	got, err := finalizedBoundaryWorkspaceFixture().RunTokenBoundariesWithWorkspace(0x10000000, 0, &tree3.Catalog{Pronunciation: markerConstantTree(501)}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil }, initialized.Workspace)
	if err != nil {
		t.Fatal(err)
	}
	if got.Pipeline.Final.State.Values[0] != 100 || got.Pipeline.Final.State.Values[1] != 0 || got.Pipeline.Final.State.States[0] != 2 || got.Pipeline.Final.State.States[1] != -1 {
		t.Fatalf("shifted workspace state = %+v", got.Pipeline.Final.State)
	}
	if int32(binary.LittleEndian.Uint32(got.Workspace[0x121a60:])) != -1 || binary.LittleEndian.Uint32(got.Workspace[0x121a64:]) != 100 || binary.LittleEndian.Uint32(got.Workspace[0x121d80:]) != 2 || !bytes.Equal(initialized.Workspace, before) {
		t.Fatal("wrong workspace writeback or input mutation")
	}
	if _, err := finalizedBoundaryWorkspaceFixture().RunTokenBoundariesWithWorkspace(0, 0, nil, nil, workspace); err == nil {
		t.Fatal("accepted absent shared tree")
	}
	if _, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace[:100], Paul2013PositionStateProgram{RowCount: 2}); err == nil {
		t.Fatal("accepted short workspace")
	}
}

func TestBoundaryWorkspaceRunsLoadedEngbiTree(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engbi.tree3")); os.IsNotExist(err) {
		t.Skip("local shared proprietary tree is absent")
	}
	tree, err := tree3.LoadPaul2013PronunciationTree(root)
	if err != nil {
		t.Fatal(err)
	}
	workspace := stateWorkspaceFixture()
	binary.LittleEndian.PutUint32(workspace[0x121a68:], 0xffffffff)
	initialized, err := ApplyPaul2013PositionStateProgramToWorkspace(workspace, Paul2013PositionStateProgram{RowCount: 2})
	if err != nil {
		t.Fatal(err)
	}
	got, err := finalizedBoundaryWorkspaceFixture().RunTokenBoundariesWithWorkspace(0x10000000, 0, &tree3.Catalog{Pronunciation: tree}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil }, initialized.Workspace)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Pipeline.MarkerTree.Scans) != 1 || len(got.Pipeline.MarkerTree.Scans[0].Inputs) != 1 || !got.Pipeline.MarkerTree.Scans[0].Inputs[0].Eligible {
		t.Fatal("loaded engbi tree did not evaluate eligible native scan")
	}
}
