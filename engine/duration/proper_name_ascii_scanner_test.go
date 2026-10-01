package duration

import (
	"encoding/binary"
	"path/filepath"
	"testing"

	"vtspeak/engine/text"
)

func TestMatchAndAppendProperNameArenaWithASCIIScanner(t *testing.T) {
	dictionary, err := text.LoadTPPDictionary(filepath.Join("..", "..", "data-common", "dict-eng"))
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{tpp: dictionary}
	for _, test := range []struct {
		name             string
		contextScanInput uint32
		scannerInput     []byte
		wantScanner      bool
		wantAppend       bool
		wantPending      bool
	}{
		{name: "zero pointer bypasses unsupported scanner", scannerInput: []byte("?\x00"), wantAppend: true},
		{name: "supported scanner reaches append", scannerInput: []byte("AGUA-DULCE\x00"), wantScanner: true, wantAppend: true},
		{name: "letter token before punctuation reaches append", scannerInput: []byte("AGUA!DULCE\x00"), wantScanner: true, wantAppend: true},
		{name: "digit token before terminal whitespace reaches append", scannerInput: []byte("12345 \x00"), wantScanner: true, wantAppend: true},
		{name: "digit token before terminal punctuation reaches append", scannerInput: []byte("12345!\x00"), wantScanner: true, wantAppend: true},
		{name: "ordinal token before terminal punctuation reaches append", scannerInput: []byte("21st!\x00"), wantScanner: true, wantAppend: true},
		{name: "mapped uppercase ordinal reaches append", scannerInput: []byte("21ST\x00"), wantScanner: true, wantAppend: true},
		{name: "nonzero pointer keeps unsupported scanner pending", contextScanInput: 1, scannerInput: []byte("?\x00"), wantPending: true},
		{name: "longer numeric lookahead stays pending", contextScanInput: 1, scannerInput: []byte("123 next\x00"), wantPending: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			componentArena := properNameComponentArena(test.contextScanInput)
			sourceArena := make([]byte, 0x14+2*0x94)
			result, supported, reason, err := engine.MatchAndAppendProperNameTPPComponentArenaWithASCIIScannerInputAndContextTables(
				sourceArena, componentArena, text.Paul2013ProperNameTPPContextGateInput{},
				test.scannerInput, nil, nil,
			)
			if err != nil {
				t.Fatal(err)
			}
			if supported != test.wantScanner || (supported == (reason != "")) {
				t.Fatalf("scanner support = %t, reason = %q; want supported=%t with a reason only when unsupported", supported, reason, test.wantScanner)
			}
			if result.Accepted != true {
				t.Fatalf("proper-name match was not accepted: %+v", result)
			}
			if result.ContextChecksPending != test.wantPending {
				t.Fatalf("context pending = %t, want %t", result.ContextChecksPending, test.wantPending)
			}
			if result.FollowupSucceeded != test.wantAppend {
				t.Fatalf("follow-up succeeded = %t, want %t", result.FollowupSucceeded, test.wantAppend)
			}
			gotArena := result.SourceArena
			if !test.wantAppend {
				gotArena = sourceArena
			}
			gotRows := binary.LittleEndian.Uint16(gotArena)
			wantRows := uint16(0)
			if test.wantAppend {
				wantRows = 2
			}
			if gotRows != wantRows {
				t.Fatalf("source row count = %d, want %d", gotRows, wantRows)
			}
		})
	}
}

func properNameComponentArena(contextScanInput uint32) []byte {
	const (
		header = 0x0c
		stride = 0x140
	)
	arena := make([]byte, header+2*stride)
	binary.LittleEndian.PutUint32(arena[:4], 2)
	for index, component := range []string{"AGUA", "DULCE"} {
		row := arena[header+index*stride : header+(index+1)*stride]
		binary.LittleEndian.PutUint32(row[0x00:0x04], uint32(index*5))
		binary.LittleEndian.PutUint32(row[0x04:0x08], uint32(index*5+len(component)))
		binary.LittleEndian.PutUint32(row[0x08:0x0c], uint32(len(component)))
		row[0x0e] = '0'
		if index == 1 {
			binary.LittleEndian.PutUint32(row[0x18:0x1c], contextScanInput)
		}
		copy(row[0x3a:], component)
		copy(row[0x46:], component)
	}
	return arena
}
