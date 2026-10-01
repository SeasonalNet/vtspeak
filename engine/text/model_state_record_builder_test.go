package text

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestPopulatePaul2013ModelStateRecordsFromParserRows(t *testing.T) {
	const rows = 2
	parserRows := make([]byte, rows*paul2013ParserSourceRowStride)
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		row := parserRows[rowIndex*paul2013ParserSourceRowStride : (rowIndex+1)*paul2013ParserSourceRowStride]
		binary.LittleEndian.PutUint16(row[0x24:], uint16(0x1200+rowIndex))
		copy(row[0x34:], []byte{0x21, 0x22})
		row[0x37] = byte(0x30 + rowIndex)
		row[0x39] = byte(0x40 + rowIndex)
		row[0x3a] = byte(0x50 + rowIndex)
		row[0x3c] = byte(0x60 + rowIndex)
		binary.LittleEndian.PutUint16(row[0x40:], uint16(0x7000+rowIndex))
		copy(row[0x66:], "ABdMc\x00")
	}
	arena := make([]byte, paul2013ModelStateFinalModeOffset+1)
	for index := range arena {
		arena[index] = 0xa5
	}
	binary.LittleEndian.PutUint16(arena[2:], rows)
	got, err := PopulatePaul2013ModelStateRecordsFromParserRows(
		arena,
		parserRows,
		[]int32{-1, 7},
		0,
		func(rowIndex, rowOffset int) (uint32, error) {
			return uint32(0x10000000 + rowIndex*paul2013ParserSourceRowStride + rowOffset), nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ContextCodes) != rows || string(got.ContextCodes[0].Codes) != "ABM" || string(got.ContextCodes[0].PhoneMarkers) != "012" || !got.ContextCodes[0].MFlags[2] {
		t.Fatalf("context codes = %+v", got.ContextCodes)
	}
	for rowIndex := 0; rowIndex < rows; rowIndex++ {
		recordStart := paul2013ModelStateRecordBaseOffset + rowIndex*paul2013ModelStateRecordStride
		record := got.Bytes[recordStart : recordStart+paul2013ModelStateRecordStride]
		for _, field := range [...]struct{ recordOffset, rowOffset int }{{0x2e0, 0x48}, {0x2e4, 0x66}, {0x3b8, 0x44}} {
			want := uint32(0x10000000 + rowIndex*paul2013ParserSourceRowStride + field.rowOffset)
			if value := binary.LittleEndian.Uint32(record[field.recordOffset:]); value != want {
				t.Errorf("row %d pointer at +%#x = %#x, want %#x", rowIndex, field.recordOffset, value, want)
			}
		}
		if record[0x95] != 3 || string(record[0x2e8:0x2eb]) != "ABM" ||
			!bytes.Equal(record[0x328:0x32b], []byte{0xa5, '1', '2'}) ||
			string(record[0x329:0x32c]) != "120" {
			t.Errorf("row %d phone fields: count=%d codes=%q previous=%x following=%q", rowIndex, record[0x95], record[0x2e8:0x2eb], record[0x328:0x32b], record[0x329:0x32c])
		}
		wantStateFlag := byte(0)
		if rowIndex == 1 {
			wantStateFlag = 0x0c
		}
		if record[0x2df] != wantStateFlag {
			t.Errorf("row %d +0x2df = %#x, want %#x", rowIndex, record[0x2df], wantStateFlag)
		}
		if record[0x100] != 0xa5 {
			t.Errorf("row %d untouched record byte +0x100 = %#x, want sentinel", rowIndex, record[0x100])
		}
		if record[0x3b1] != byte(0x30+rowIndex) || record[0x3b2] != byte(0x40+rowIndex) || record[0x3b3] != byte(0x50+rowIndex) {
			t.Errorf("row %d trailing controls = %x %x %x", rowIndex, record[0x3b1], record[0x3b2], record[0x3b3])
		}
	}
	if got.Bytes[0x64c+0x2df] != 0 {
		t.Fatalf("first row per-token state byte = %#x, want zero for input -1", got.Bytes[0x64c+0x2df])
	}
	if got.Bytes[paul2013ModelStateFinalModeOffset] != 7 {
		t.Fatalf("final mode byte = %d, want 7", got.Bytes[paul2013ModelStateFinalModeOffset])
	}
}

func TestPopulatePaul2013ModelStateRecordsDerivesFinalMode(t *testing.T) {
	for _, test := range []struct {
		name       string
		stateCode  int16
		stateValue int32
		pitchCount int16
		wantMode   byte
	}{
		{name: "ordinary", stateCode: 2, stateValue: -1, wantMode: 5},
		{name: "empty terminal pitch", stateCode: 3, stateValue: -1, wantMode: 6},
		{name: "state four", stateCode: 4, stateValue: 7, wantMode: 7},
	} {
		t.Run(test.name, func(t *testing.T) {
			rows := make([]byte, paul2013ParserSourceRowStride)
			rows[0x66] = 'A'
			rows[0x67] = 0
			binary.LittleEndian.PutUint16(rows[0x40:], uint16(test.stateCode))
			arena := make([]byte, paul2013ModelStateFinalModeOffset+1)
			binary.LittleEndian.PutUint16(arena[2:], 1)
			got, err := PopulatePaul2013ModelStateRecordsFromParserRows(
				arena,
				rows,
				[]int32{test.stateValue},
				test.pitchCount,
				func(int, int) (uint32, error) { return 1, nil },
			)
			if err != nil {
				t.Fatal(err)
			}
			if mode := got.Bytes[paul2013ModelStateFinalModeOffset]; mode != test.wantMode {
				t.Fatalf("final mode byte = %d, want %d", mode, test.wantMode)
			}
		})
	}
}

func TestPopulatePaul2013ModelStateRecordsRejectsInvalidInputs(t *testing.T) {
	validRows := make([]byte, paul2013ParserSourceRowStride)
	validRows[0x66] = 'A'
	validRows[0x67] = 0
	validArena := make([]byte, paul2013ModelStateFinalModeOffset+1)
	binary.LittleEndian.PutUint16(validArena[2:], 1)
	for name, input := range map[string]struct {
		arena   []byte
		rows    []byte
		states  []int32
		resolve Paul2013ParserRowAddressResolver
	}{
		"record count mismatch":    {validArena, validRows, []int32{-1, -1}, func(int, int) (uint32, error) { return 1, nil }},
		"missing pointer resolver": {validArena, validRows, []int32{-1}, nil},
		"missing code terminator": func() struct {
			arena   []byte
			rows    []byte
			states  []int32
			resolve Paul2013ParserRowAddressResolver
		} {
			rows := append([]byte(nil), validRows...)
			for index := 0x66; index < len(rows); index++ {
				rows[index] = 'A'
			}
			return struct {
				arena   []byte
				rows    []byte
				states  []int32
				resolve Paul2013ParserRowAddressResolver
			}{validArena, rows, []int32{-1}, func(int, int) (uint32, error) { return 1, nil }}
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := PopulatePaul2013ModelStateRecordsFromParserRows(input.arena, input.rows, input.states, 0, input.resolve); err == nil {
				t.Fatal("PopulatePaul2013ModelStateRecordsFromParserRows() succeeded")
			}
		})
	}
}
