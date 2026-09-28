package dat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCapturedDLLPayloads(t *testing.T) {
	for i := range 32 {
		name := fmt.Sprintf("candidate-%02d", i)
		t.Run(name, func(t *testing.T) {
			path := filepath.Join("..", "..", "tools", "revkit", "work", "stage2-copy")
			payload, err := os.ReadFile(filepath.Join(path, name+".dat"))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(path, name+".pcm"))
			if err != nil {
				t.Fatal(err)
			}
			got, err := Decode(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("decoded %d bytes, expected %d; PCM differs", len(got), len(want))
			}
		})
	}
}

func TestRejectsTruncatedPayload(t *testing.T) {
	if _, err := Decode(nil); err == nil {
		t.Fatal("empty payload accepted")
	}
}

func TestUnitPayloadAndWAV(t *testing.T) {
	index := make([]byte, indexHeaderSize+unitStride+featureStride)
	index[0] = byte(len(paul2013Header))
	copy(index[1:], paul2013Header)
	binary.LittleEndian.PutUint16(index[24:26], 1)
	binary.LittleEndian.PutUint16(index[26:28], 10)
	copy(index[28:38], "merged-gen")
	binary.LittleEndian.PutUint32(index[39:43], 1)
	binary.LittleEndian.PutUint16(index[43:45], unitStride)
	binary.LittleEndian.PutUint32(index[indexHeaderSize:], 2)
	binary.LittleEndian.PutUint16(index[indexHeaderSize+8:], 3)
	index[indexHeaderSize+unitStride+1] = 9
	index[indexHeaderSize+unitStride+2] = 8
	index[indexHeaderSize+unitStride+3] = 7
	index[indexHeaderSize+unitStride+4] = 6
	index[indexHeaderSize+unitStride+5] = 5
	index[indexHeaderSize+unitStride+6] = 4
	index[indexHeaderSize+unitStride+7] = 3
	index[indexHeaderSize+unitStride+8] = 2
	record, err := ReadUnit(index, 0)
	if err != nil || record.Signature != [7]byte{9, 8, 7, 6, 5, 4, 3} {
		t.Fatalf("record signature = %v, %v", record.Signature, err)
	}
	got, err := UnitPayload(index, bytes.NewReader([]byte{0, 0, 1, 2, 3}), 5, 0)
	if err != nil || !bytes.Equal(got, []byte{1, 2, 3}) {
		t.Fatalf("payload = %v, %v", got, err)
	}
	if _, err := UnitPayload(index, bytes.NewReader(nil), 0, 1); err == nil {
		t.Fatal("out-of-range unit accepted")
	}
	if _, err := UnitPayload(index[:indexHeaderSize], bytes.NewReader(nil), 0, 0); err == nil {
		t.Fatal("truncated record accepted")
	}
	upm, err := UPMPayload(UnitRecord{UPMOffset: 1, UPMFirstCount: 2, UPMSecondCount: 2}, bytes.NewReader([]byte{0, 10, 11, 12}), 4)
	if err != nil || !bytes.Equal(upm, []byte{10, 11, 12}) {
		t.Fatalf("UPM span = %v, %v", upm, err)
	}
	wav, err := WAV([]byte{0x34, 0x12})
	if err != nil {
		t.Fatal(err)
	}
	if len(wav) != 46 || string(wav[:4]) != "RIFF" || binary.LittleEndian.Uint32(wav[40:44]) != 2 || !bytes.Equal(wav[44:], []byte{0x34, 0x12}) {
		t.Fatal("invalid WAV framing")
	}
}
