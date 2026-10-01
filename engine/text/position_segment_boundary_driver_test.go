package text

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"vtspeak/engine/tree3"
)

func TestPositionSegmentDriverFeedsPhoneRecordsAndBoundaries(t *testing.T) {
	engineContext, workspace, shared, model := segmentDriverFixture()
	parser := segmentParserFixture(2, 8)
	for index := 0; index < 2; index++ {
		copy(parser[index*0x94+0x66:], []byte{1, 0x13, 0})
	}
	binary.LittleEndian.PutUint32(parser[0x20:], 0xfffffffe)
	parse := func(context.Context, uint32, []byte) (Paul2013PositionSegmentParserOutput, error) {
		return Paul2013PositionSegmentParserOutput{State: parser}, nil
	}
	resolveState := func(_ uint32, count int) ([]int32, error) {
		cells := make([]int32, count)
		for index := range cells {
			cells[index] = int32(index)
		}
		return cells, nil
	}
	resolveAddress := func(index, offset int) (uint32, error) { return uint32(0x200000 + index*0x94 + offset), nil }
	resolveStrings := func(pointer uint32) ([]byte, error) {
		offset := int(pointer) - 0x200000
		if offset < 0 || offset >= len(parser) {
			return nil, fmt.Errorf("invalid parser pointer")
		}
		return parser[offset:], nil
	}
	got, err := RunPaul2013PositionSegmentBoundaryDriver(context.Background(), engineContext, workspace, shared, model, parse, resolveState, resolveAddress, 0x10000000, &tree3.Catalog{Pronunciation: markerConstantTree(0)}, resolveStrings)
	if err != nil {
		t.Fatal(err)
	}
	if !got.HasBoundary || got.Segment.ReturnValue != 8 || len(got.Records.ContextCodes) != 2 || binary.LittleEndian.Uint32(got.Boundary.Workspace[4:8]) != 8 {
		t.Fatal("segment handoff missing")
	}
	if got.Boundary.Pipeline.Final.State.Values[0] != 100 || binary.LittleEndian.Uint32(got.Boundary.Workspace[0x121a64:]) != 100 || got.Records.Bytes[0x4770a] != 5 {
		t.Fatalf("boundary writeback = %+v", got.Boundary.Pipeline.Final.State)
	}
	for index := 0; index < 2; index++ {
		row := got.Records.Bytes[0x64c+index*0x3c0:]
		if row[0x95] != 2 || row[0x2df] != 0 || binary.LittleEndian.Uint32(row[0x2e4:]) != uint32(0x200000+index*0x94+0x66) || binary.LittleEndian.Uint32(row[:4]) != uint32(index*4) {
			t.Fatal("projected record fields or mapped intervals differ")
		}
	}
	segment := got.Segment
	segment.ParserState = segment.ParserState[:7]
	if _, err := segment.PopulateModelStateRecords(resolveAddress); err == nil {
		t.Fatal("accepted short parser rows")
	}
	segment = got.Segment
	segment.ParserState = append([]byte(nil), segment.ParserState...)
	binary.LittleEndian.PutUint16(segment.ParserState[:2], 1)
	if _, err := segment.PopulateModelStateRecords(resolveAddress); err == nil {
		t.Fatal("accepted mismatched segment count")
	}
}

func TestPositionSegmentBoundaryDriverCompletionSkipsDependencies(t *testing.T) {
	engineContext, workspace, shared, model := segmentDriverFixture()
	got, err := RunPaul2013PositionSegmentBoundaryDriver(context.Background(), engineContext, workspace, shared, model, func(context.Context, uint32, []byte) (Paul2013PositionSegmentParserOutput, error) {
		return Paul2013PositionSegmentParserOutput{State: segmentParserFixture(0, 0)}, nil
	}, nil, nil, 0, nil, nil)
	if err != nil || got.HasBoundary || !got.Segment.EndOfSource {
		t.Fatal("completion invoked record/tree dependencies")
	}
	if _, err := got.Segment.PopulateModelStateRecords(nil); err == nil {
		t.Fatal("projected completed segment")
	}
}
