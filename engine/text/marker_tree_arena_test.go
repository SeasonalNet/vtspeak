package text

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"vtspeak/engine/tree3"
)

func markerArenaFixture(count int) []byte {
	arena := boundaryTestArena(count)
	for index := 0; index < count; index++ {
		row := arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		binary.LittleEndian.PutUint32(row[0:4], uint32(index*10))
		binary.LittleEndian.PutUint32(row[4:8], uint32(index*10+2))
		binary.LittleEndian.PutUint32(row[0x2e4:0x2e8], uint32(0x1000+index*0x100))
		row[0x94], row[0x95] = byte(index+1), 2
		row[0x2e8], row[0x2e9] = 0x01, 0x13
		row[0x329], row[0x32a] = '0', '0'
		row[0x3bd] = ']'
		if index == count-1 {
			row[0x3bd] = 'Z'
		}
	}
	return arena
}

func markerConstantTree(value int16) *tree3.Tree {
	return &tree3.Tree{OutputWidth: 1, Nodes: []tree3.Node{{Feature: 0, Threshold: 32767, WhenTrue: -1, WhenFalse: -1}}, Outputs: [][]int16{{value}}}
}

func TestMarkerTreeArenaReverseScanAndPreviousTerminal(t *testing.T) {
	arena := markerArenaFixture(4)
	before := append([]byte(nil), arena...)
	var addresses []uint32
	got, err := RunPaul2013MarkerTreeArena(arena, markerConstantTree(501), func(address uint32) ([]byte, error) {
		addresses = append(addresses, address)
		return []byte{0x01, 0x13, 0}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(addresses, []uint32{0x1300, 0x1200, 0x1200, 0x1100, 0x1100, 0x1000}) {
		t.Fatalf("pointer order = %x", addresses)
	}
	if len(got.Scans) != 1 || len(got.Scans[0].Inputs) != 3 {
		t.Fatalf("scans = %+v", got.Scans)
	}
	scan := got.Scans[0]
	for index, input := range scan.Inputs {
		if !input.Eligible || input.Features[0] != int16(3-index) || input.Features[1] != int16(index+1) || input.Features[2] != 4 || input.Features[5] != 10 || input.Features[6] != 1 || input.Features[7] != int16(4-index) {
			t.Fatalf("visit %d = %+v", index, input)
		}
		if input.Features[4] != []int16{4, 7, 9}[index] {
			t.Fatalf("visit %d cumulative = %d", index, input.Features[4])
		}
	}
	for index := 0; index < 4; index++ {
		row := got.Arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		want := byte('\\')
		if index == 3 {
			want = 'Z'
		}
		if row[0x3bd] != want || row[0x91] != 0xa5 || row[0x3bc] != 0xa5 {
			t.Fatalf("record %d wrong write: terminal %q, +91 %x, flag %x", index, row[0x3bd], row[0x91], row[0x3bc])
		}
	}
	if scan.EndState.PhonesSinceLastSplit != 0 || scan.EndState.ValueSinceLastSplit != 0 || scan.EndState.CumulativePhoneValue != 9 {
		t.Fatalf("end state = %+v", scan.EndState)
	}
	if !bytes.Equal(arena, before) {
		t.Fatal("mutated caller arena")
	}
}

func TestMarkerTreeArenaCutoffAndGroupCounters(t *testing.T) {
	arena := markerArenaFixture(4)
	arena[paul2013TokenMarkerArenaRecords+paul2013TokenMarkerArenaStride+0x3bd] = '['
	got, err := RunPaul2013MarkerTreeArena(arena, markerConstantTree(500), func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Arena, arena) {
		t.Fatal("500 inserted marker")
	}
	if len(got.Scans) != 2 || got.Scans[0].Inputs[0].Features[5] != 3 || got.Scans[1].Inputs[0].Features[5] != 7 || got.Scans[1].EndState.PhoneOrdinal != 1 {
		t.Fatalf("group-local spans/counters = %+v", got.Scans)
	}
}

func TestMarkerTreeArenaSkippedVisitCountsWithoutDereferences(t *testing.T) {
	arena := markerArenaFixture(3)
	for index := 1; index < 3; index++ {
		previous := arena[paul2013TokenMarkerArenaRecords+(index-1)*paul2013TokenMarkerArenaStride:]
		current := arena[paul2013TokenMarkerArenaRecords+index*paul2013TokenMarkerArenaStride:]
		binary.LittleEndian.PutUint32(current[0:4], binary.LittleEndian.Uint32(previous[4:8])+1)
	}
	got, err := RunPaul2013MarkerTreeArena(arena, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Scans[0].EndState.PhoneOrdinal != 2 || got.Scans[0].EndState.CumulativePhoneValue != 5 || got.Scans[0].Inputs[0].Eligible || !bytes.Equal(got.Arena, arena) {
		t.Fatalf("skipped scan = %+v", got.Scans)
	}
}

func TestMarkerTreeArenaRejectsMissingTreeAndPointers(t *testing.T) {
	arena := markerArenaFixture(2)
	if _, err := RunPaul2013MarkerTreeArena(arena, markerConstantTree(0), nil); err == nil {
		t.Fatal("accepted missing resolver")
	}
	resolver := func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil }
	if _, err := RunPaul2013MarkerTreeArena(arena, nil, resolver); err == nil {
		t.Fatal("accepted missing eligible tree")
	}
	if _, err := RunPaul2013MarkerTreeArena(arena, markerConstantTree(0), func(uint32) ([]byte, error) { return []byte{1}, nil }); err == nil {
		t.Fatal("accepted short selected string")
	}
	if _, err := RunPaul2013MarkerTreeArena(arena, markerConstantTree(0), func(uint32) ([]byte, error) { return []byte{0x80, 0x13, 0}, nil }); err == nil {
		t.Fatal("accepted signed-char table underflow")
	}
}

func TestPreparedMarkerTreeBoundaryComposition(t *testing.T) {
	prepared, err := PreparePaul2013TokenBoundaryArena(markerArenaFixture(3), 0x20000000)
	if err != nil {
		t.Fatal(err)
	}
	marked, err := RunPaul2013MarkerTreeArena(prepared.PhoneGroups.Arena, markerConstantTree(501), func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil })
	if err != nil {
		t.Fatal(err)
	}
	finished, err := FinishPaul2013TokenBoundaryArena(marked.Arena, []int32{-1, -1, -1}, []int32{4, 4, 4}, []int32{-1, -1})
	if err != nil {
		t.Fatal(err)
	}
	if finished.State.Markers[0].Marker != '\\' || finished.State.Markers[1].Marker != '\\' || finished.State.Markers[2].Marker != 'Z' || finished.Groups[0].OpaqueRowSpan != 3 {
		t.Fatalf("prepared-tree-finish = %+v/%+v", finished.State.Markers, finished.Groups)
	}
	combined, err := RunPaul2013TokenBoundaryPipeline(markerArenaFixture(3), 0x20000000, markerConstantTree(501), func(uint32) ([]byte, error) { return []byte{1, 0x13, 0}, nil }, []int32{-1, -1, -1}, []int32{4, 4, 4}, []int32{-1, -1})
	if err != nil {
		t.Fatal(err)
	}
	finished.Arena, finished.Groups, err = PopulatePaul2013RecordGroupDescriptorArena(finished.Arena, 0x20000000)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(combined.Final, finished) {
		t.Fatal("composed pipeline differs from native stage sequence")
	}
}
