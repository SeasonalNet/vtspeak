package text

import (
	"bytes"
	"encoding/binary"
	"testing"

	"vtspeak/engine/tree3"
)

func TestPositionSegmentDefaultsNativeSelectionClampsAndPreservation(t *testing.T) {
	for _, flag := range []byte{0, 1, 2, 0xff} {
		for _, value := range []int32{-10, 100, 999} {
			context := bytes.Repeat([]byte{0xa5}, 0x4cfc)
			workspace := bytes.Repeat([]byte{0xb6}, 0x1312c8)
			for _, field := range []struct{ flag, override, context int }{{0x1210d7, 0x1210dc, 0x4cf4}, {0x1210d8, 0x1210e0, 0x4cf0}, {0x1210d9, 0x1210e4, 0x4cf8}} {
				workspace[field.flag] = flag
				binary.LittleEndian.PutUint32(workspace[field.override:], uint32(value))
				binary.LittleEndian.PutUint32(context[field.context:], 175)
			}
			before := append([]byte(nil), workspace...)
			expected := append([]byte(nil), workspace...)
			selected := int32(175)
			if flag == 1 {
				selected = value
			}
			for _, field := range []struct {
				active, initial  int
				minimum, maximum int32
			}{{0x1210e8, 0x1210fc, 50, 200}, {0x1210ec, 0x1210f4, 50, 400}, {0x1210f0, 0x1210f8, 0, 500}} {
				clamped := min(max(selected, field.minimum), field.maximum)
				binary.LittleEndian.PutUint32(expected[field.active:], uint32(clamped))
				binary.LittleEndian.PutUint32(expected[field.initial:], uint32(clamped))
			}
			binary.LittleEndian.PutUint32(expected[0x1312c4:], 0xffffffff)
			got, err := PreparePaul2013PositionSegmentWorkspace(context, workspace)
			if err != nil || !bytes.Equal(got, expected) || !bytes.Equal(workspace, before) {
				t.Fatalf("flag %d value %d: invalid defaults or memory writes: %v", flag, value, err)
			}
		}
	}
	for _, lengths := range [][2]int{{0x4cfb, 0x1312c8}, {0x4cfc, 0x1312c7}} {
		if _, err := PreparePaul2013PositionSegmentWorkspace(make([]byte, lengths[0]), make([]byte, lengths[1])); err == nil {
			t.Fatal("accepted truncated segment defaults")
		}
	}
}

func TestPreparedPositionSegmentPipelineAdvancesCursor(t *testing.T) {
	for _, start := range []uint32{0, 100, 0xfffffff8} {
		workspace, shared, state := positionIndexWorkspaceFixture()
		workspace = append(workspace, make([]byte, 0x1312c8-len(workspace))...)
		binary.LittleEndian.PutUint32(workspace[4:8], start)
		context := make([]byte, 0x4cfc)
		for _, offset := range []int{0x4cf4, 0x4cf0, 0x4cf8} {
			binary.LittleEndian.PutUint32(context[offset:], 150)
		}
		before := append([]byte(nil), workspace...)
		got, err := state.RunPreparedNativePositionWorkspaceBoundaries(context, workspace, shared, func(_ uint32, count int) ([]int32, error) {
			cells := make([]int32, count)
			for index := range cells {
				cells[index] = int32(index)
			}
			return cells, nil
		}, 0x10000000, 0, &tree3.Catalog{Pronunciation: markerConstantTree(0)}, func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
		if err != nil {
			t.Fatal(err)
		}
		for _, memory := range [][]byte{got.Position.Workspace, got.Boundary.Workspace} {
			if binary.LittleEndian.Uint32(memory[4:8]) != start+15 || binary.LittleEndian.Uint32(memory[0x1312c4:]) != 0xffffffff {
				t.Fatal("missing cursor advancement or reset")
			}
		}
		for _, mode := range []int{0, 1, 2} {
			if got.Position.Program.Arrays[mode][0] != 150 {
				t.Fatal("prepared defaults did not reach state arrays")
			}
		}
		if !bytes.Equal(workspace, before) {
			t.Fatal("mutated caller segment workspace")
		}
	}
}

func TestNativeSegmentHeaderAndEndOfSourceErrors(t *testing.T) {
	for _, offset := range []int{-1, 0, 1000000} {
		workspace, shared, state := positionIndexWorkspaceFixture()
		binary.LittleEndian.PutUint32(state.Finalized.StateArena[4:8], 0)
		if _, err := state.RunNativePositionWorkspaceBoundaries(workspace, shared, nil, 0, offset, nil, nil); err == nil {
			t.Fatal("accepted unsupported/truncated segment")
		}
	}
}
