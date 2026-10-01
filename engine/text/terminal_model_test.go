package text

import (
	"context"
	"debug/pe"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"vtspeak/engine/tree3"
)

func TestTerminalFollowingWordClassifierBranches(t *testing.T) {
	classifier, err := NewPaul2013TerminalWordClassifier(&Table{Name: "sbdw_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "dictionary", Value: "19"}}})
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		text string
		want int16
	}{{"DICTIONARY", 1}, {"42!", 9}, {"aB-C", 3}, {".next", 17}, {",next", 11}, {"'next", 15}, {"`", 15}, {"\"", 15}, {":", 12}, {"-", 12}, {"(", 13}, {"[", 13}, {"{", 13}, {"<", 13}, {")", 14}, {"]", 14}, {"}", 14}, {">", 14}, {"kindliness", 1}, {"allowance", 1}, {"action", 1}, {"kindness", 1}, {"shipment", 1}, {"revision", 1}, {"friendship", 1}, {"acidity", 1}, {"artist", 1}, {"harmless", 5}, {"working", 5}, {"useful", 5}, {"slowly", 5}, {"ly", 0}, {"word", 0}, {"Word", 0}, {"can't", 0}, {"kindness's", 1}, {"kindness'd", 1}, {"kindness'm", 1}, {"kindness'em", 1}, {"kindness've", 1}, {"kindness're", 1}, {"kindness'll", 1}, {"kindness'nonsense", 0}, {"@", 16}, {"", 16}, {"\xff", 16}, {"word@", 0}} {
		got, err := classifier.Classify([]byte(test.text))
		if err != nil || got != test.want {
			t.Fatalf("%q -> %d want %d: %v", test.text, got, test.want, err)
		}
	}
	if _, err := classifier.Classify([]byte(strings.Repeat("a", 32) + "'s")); err == nil {
		t.Fatal("native contraction buffer overrun accepted")
	}
	if _, err := NewPaul2013TerminalWordClassifier(&Table{Name: "sbdw_sort.txt2", Mode: TwoColumns, Rows: []Row{{Key: "x", Value: "48"}}}); err == nil {
		t.Fatal("class table overread accepted")
	}
}

func terminalModelFixture(t *testing.T) *Paul2013TerminalModel {
	t.Helper()
	tables := Tables{}
	for _, name := range []string{"abbrh_sort.txt2", "abbrt_sort.txt2", "abbrc_sort.txt2", "sbdw_sort.txt2"} {
		tables[name] = &Table{Name: name, Mode: TwoColumns}
	}
	tables["abbrh_sort.txt2"].Rows = []Row{{Key: "word", Value: "2"}}
	tables["abbrt_sort.txt2"].Rows = []Row{{Key: "kindness", Value: "2"}}
	tree := &tree3.Tree{OutputWidth: 1, Nodes: []tree3.Node{{Feature: 14, Operation: 'D', Values: []int16{1}, WhenTrue: -1, WhenFalse: -2}}, Outputs: [][]int16{{1}, {0}}}
	model, err := NewPaul2013TerminalModel(tables, tree)
	if err != nil {
		t.Fatal(err)
	}
	return model
}

func TestTerminalModelCompleteKeyAndRecognizerIntegration(t *testing.T) {
	model := terminalModelFixture(t)
	window := terminalContextFixture("word", ".", "kindness", 1)
	window[2].Status = 1
	got, err := model.BuildKey(window)
	want := [16]int16{4, 2, 1, 0, 1, 0, -1, 8, 2, 0, 1, 1, 0, -1, 1, 1}
	if err != nil || got != want {
		t.Fatalf("key %v want %v: %v", got, want, err)
	}
	empty, err := model.BuildKey(Paul2013TerminalContextWindow{})
	want = [16]int16{0, -1, -1, -1, -1, -1, -1, 0, -1, -1, -1, -1, -1, -1, -1, -1}
	if err != nil || empty != want {
		t.Fatal(empty, err)
	}
	for _, test := range []struct {
		word string
		want int32
	}{{"kindness", 'P'}, {"later", 'N'}} {
		full := []byte("word. " + test.word + "\x00")
		recognize := NewPaul2013TerminalContextRecognizer(full, fixtureTerminalWordScanner, model.Lookup)
		current := Paul2013ModelScannerResult{Text: []byte("."), TextLength: 1, Field14: 1, Status: 3}
		value, err := recognize(context.Background(), terminalStateFixture("word", 0), &current, full[4:], 4)
		if err != nil || value != test.want {
			t.Fatalf("recognize %s -> %c, %v", test.word, value, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := model.Lookup(ctx, window); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestTerminalWordClassifierAndTreeLocalResources(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "sbd.tree3")); os.IsNotExist(err) {
		t.Skip("local proprietary shared data absent")
	}
	model, err := LoadPaul2013TerminalModel(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("local sbd.tree3: %d nodes, %d leaves, scalar width %d", len(model.tree.Nodes), len(model.tree.Outputs), model.tree.OutputWidth)
	raw, err := os.ReadFile(filepath.Join(root, "sbdw_sort.txt2"))
	if err != nil {
		t.Fatal(err)
	}
	table, err := ParseTXT2("sbdw_sort.txt2", raw)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range table.Rows {
		code, err := strconv.Atoi(row.Value)
		if err != nil {
			t.Fatal(err)
		}
		for _, word := range []string{row.Key, strings.ToUpper(row.Key), strings.ToLower(row.Key)} {
			got, err := model.words.Classify([]byte(word))
			if err != nil || got != paul2013TerminalWordClasses[code] {
				t.Fatalf("%s -> %d, %v", word, got, err)
			}
		}
	}
	for _, word := range []string{"next", "Kindness", "42", "Mr", "word's", "?"} {
		window := terminalContextFixture("word", ".", word, 1)
		window[2].Status = 1
		key, err := model.BuildKey(window)
		if err != nil {
			t.Fatal(err)
		}
		_, direct, err := model.tree.Evaluate(key[:])
		if err != nil {
			t.Fatal(err)
		}
		result, err := model.Lookup(context.Background(), window)
		if err != nil || result != uint32(uint16(direct[0])) {
			t.Fatal(result, direct, err)
		}
	}
}

func TestTerminalWordClassMapMatchesLocalDLL(t *testing.T) {
	file, err := pe.Open("../../binary/vt_pau.dll")
	if os.IsNotExist(err) {
		t.Skip("local proprietary DLL absent")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	header, ok := file.OptionalHeader.(*pe.OptionalHeader32)
	if !ok || header.ImageBase != 0x10000000 {
		t.Fatal("unexpected PE image base")
	}
	rva := uint32(0x10081358) - header.ImageBase
	for _, section := range file.Sections {
		if rva < section.VirtualAddress || uint64(rva-section.VirtualAddress)+uint64(len(paul2013TerminalWordClasses)*2) > uint64(section.Size) {
			continue
		}
		raw, err := section.Data()
		if err != nil {
			t.Fatal(err)
		}
		start := int(rva - section.VirtualAddress)
		for index, want := range paul2013TerminalWordClasses {
			if got := int16(binary.LittleEndian.Uint16(raw[start+index*2:])); got != want {
				t.Fatalf("word class map [%d]=%d want %d", index, got, want)
			}
		}
		return
	}
	t.Fatal("word class translation table is unmapped")
}
