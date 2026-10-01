package text

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseEmbeddedDictionaryUsesIndexedRecordPartitions(t *testing.T) {
	index := make([]byte, 12)
	binary.LittleEndian.PutUint32(index[0:], 0)
	binary.LittleEndian.PutUint32(index[4:], 6)
	binary.LittleEndian.PutUint32(index[8:], 0)
	data := []byte("go\x00\x01\x02\x00hi\x00\x03\x04\x00")
	dictionary, err := ParseEmbeddedDictionary(index, data)
	if err != nil {
		t.Fatal(err)
	}
	if dictionary.Len() != 2 {
		t.Fatalf("record count = %d, want 2", dictionary.Len())
	}
	record, ok := dictionary.Lookup("go")
	if !ok || string(record.Payload) != "\x01\x02\x00" {
		t.Fatalf("lookup go = (%+v, %t)", record, ok)
	}
	if _, ok := dictionary.Lookup("GO"); ok {
		t.Fatal("lookup unexpectedly changed key case")
	}
}

func TestParseEmbeddedDictionaryRejectsBrokenPartitions(t *testing.T) {
	tests := map[string]struct {
		index []byte
		data  []byte
	}{
		"short index":      {index: []byte{1, 2}, data: []byte("a\x00b\x00")},
		"misaligned index": {index: []byte{0, 0, 0, 0, 0}, data: []byte("a\x00b\x00")},
		"first offset": {
			index: offsets(4), data: []byte("a\x00b\x00"),
		},
		"duplicate offset": {
			index: offsets(0, 0), data: []byte("a\x00b\x00"),
		},
		"missing payload terminator": {
			index: offsets(0), data: []byte("a\x00b"),
		},
		"multiple payloads": {
			index: offsets(0), data: []byte("a\x00b\x00c\x00"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseEmbeddedDictionary(test.index, test.data); err == nil {
				t.Fatal("invalid dictionary accepted")
			}
		})
	}
}

func TestEncodeEmbeddedKeyMapsAndGreedilyCompressesPairs(t *testing.T) {
	tables := EmbeddedKeyTables{Pairs: [][2]byte{{'A', 'B'}, {'C', 'D'}}}
	for i := range tables.CharacterMap {
		tables.CharacterMap[i] = byte(i)
	}
	encoded, err := EncodeEmbeddedKey([]byte("ABCDEFG"), tables)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x81, 0x82, 'E', 'F', 'G'}
	if string(encoded) != string(want) {
		t.Fatalf("encoded key = % x, want % x", encoded, want)
	}

	if _, err := EncodeEmbeddedKey(make([]byte, 151), tables); err == nil {
		t.Fatal("oversized key accepted")
	}
	if _, err := EncodeEmbeddedKey([]byte("AB"), EmbeddedKeyTables{Pairs: [][2]byte{{'B', 'C'}, {'A', 'Z'}}}); err == nil {
		t.Fatal("unsorted key pair table accepted")
	}
}

func TestLookupSurfaceUsesEncodedKey(t *testing.T) {
	dictionary, err := ParseEmbeddedDictionary(offsets(0), []byte{0x81, 0, 0x01, 0})
	if err != nil {
		t.Fatal(err)
	}
	tables := EmbeddedKeyTables{Pairs: [][2]byte{{'A', 'B'}}}
	for i := range tables.CharacterMap {
		tables.CharacterMap[i] = byte(i)
	}
	record, ok, err := dictionary.LookupSurface([]byte("AB"), tables)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || record.Key != string([]byte{0x81}) || string(record.Payload) != "\x01\x00" {
		t.Fatalf("surface lookup = (%+v, %t)", record, ok)
	}
}

func offsets(values ...uint32) []byte {
	data := make([]byte, 4+len(values)*4)
	for i, value := range values {
		binary.LittleEndian.PutUint32(data[4+i*4:], value)
	}
	return data
}

func TestParsePhonePayloadFormsAndMetadata(t *testing.T) {
	direct, err := ParsePhonePayload([]byte{0x4d, 0x01, 0x02, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if direct.ResultType != 'A' || direct.Metadata != [4]bool{true, false, false, false} {
		t.Fatalf("direct metadata = %+v", direct)
	}
	if len(direct.Pronunciations) != 1 || string(direct.Pronunciations[0].Phone) != "\x01\x02" {
		t.Fatalf("direct pronunciation = %+v", direct.Pronunciations)
	}

	paths, err := ParsePhonePayload([]byte{0x02, 0x02, 0x03, '|', 0x01, 0x04, 0xff, '|', 0x06, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths.Pronunciations) != 2 {
		t.Fatalf("alternative count = %d, want 2", len(paths.Pronunciations))
	}
	if string(paths.Pronunciations[0].Path) != "\x01\x02" || string(paths.Pronunciations[0].Phone) != "\x01\x04" {
		t.Fatalf("first alternative = %+v", paths.Pronunciations[0])
	}
	if len(paths.Pronunciations[1].Path) != 0 || string(paths.Pronunciations[1].Phone) != "\x06" {
		t.Fatalf("second alternative = %+v", paths.Pronunciations[1])
	}

	marker, err := ParsePhonePayload([]byte{0})
	if err != nil || len(marker.Pronunciations) != 0 {
		t.Fatalf("marker payload = %+v, error = %v", marker, err)
	}
}

func TestExpandPronunciationsThroughInjectedCodebook(t *testing.T) {
	payload, err := ParsePhonePayload([]byte{0x01, 0x10, 0x11, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	var codebook PhoneIDCodebook
	codebook[0x10] = [5]byte{0x22, 0x17, 0}
	codebook[0x11] = [5]byte{0x2b, 0}
	expanded, err := payload.ExpandPronunciations(codebook)
	if err != nil {
		t.Fatal(err)
	}
	if len(expanded) != 1 || string(expanded[0].Symbols) != "\x22\x17\x2b" {
		t.Fatalf("expanded pronunciation = %+v", expanded)
	}
	emptyCodebook := PhoneIDCodebook{}
	if _, err := payload.ExpandPronunciations(emptyCodebook); err == nil {
		t.Fatal("unmapped phone ID accepted")
	}
}

func TestPaul2013PhoneIDCodebookCoversDocumentedSlots(t *testing.T) {
	codebook := Paul2013PhoneIDCodebook()
	if got, want := codebook[0x01], [5]byte{0x07, 0x13, 0x07, 0x2b, 0}; got != want {
		t.Errorf("ID 0x01 = % x, want % x", got, want)
	}
	if got, want := codebook[0xd9], [5]byte{0x22, 0, 0, 0, 0}; got != want {
		t.Errorf("ID 0xd9 = % x, want % x", got, want)
	}
	if got, want := codebook[0xce], [5]byte{0x17, 0, 0, 0, 0}; got != want {
		t.Errorf("ID 0xce = % x, want % x", got, want)
	}
	if got, want := codebook[0x26], [5]byte{0x2b, 0x30, 0, 0, 0}; got != want {
		t.Errorf("ID 0x26 = % x, want % x", got, want)
	}
	if got, want := codebook[0xfd], [5]byte{0x64, 0, 0, 0, 0}; got != want {
		t.Errorf("ID 0xfd = % x, want % x", got, want)
	}
	if got, want := codebook[0xfe], [5]byte{0x63, 0, 0, 0, 0}; got != want {
		t.Errorf("ID 0xfe = % x, want % x", got, want)
	}
	if codebook[0] != [5]byte{} || codebook[0xff] != [5]byte{} {
		t.Fatalf("reserved slots 0x00/0xff are not empty: % x / % x", codebook[0], codebook[0xff])
	}

	assigned := 0
	for id, slot := range codebook {
		terminator := indexByte(slot[:], 0)
		if terminator == 0 {
			if id != 0 && id != 0xff {
				t.Errorf("unexpected empty slot 0x%02x", id)
			}
			continue
		}
		assigned++
		if terminator < 0 || terminator > 4 {
			t.Errorf("slot 0x%02x is not NUL-terminated within five bytes: % x", id, slot)
			continue
		}
		for _, symbol := range slot[:terminator] {
			if symbol == 0 || (symbol > 0x45 && symbol != 0x63 && symbol != 0x64) {
				t.Errorf("slot 0x%02x contains undocumented symbol 0x%02x", id, symbol)
			}
		}
	}
	if assigned != 254 {
		t.Errorf("assigned slot count = %d, want 254", assigned)
	}
}

func TestDecodeCMUPhonesCoversObservedSymbolsAndRejectsControls(t *testing.T) {
	if len(internalPhoneLabels) != 0x45 {
		t.Fatalf("CMU label count = %d, want %d", len(internalPhoneLabels), 0x45)
	}
	for symbol := byte(1); symbol <= 0x45; symbol++ {
		if _, ok := internalPhoneLabels[symbol]; !ok {
			t.Errorf("internal phone byte 0x%02x has no label", symbol)
		}
	}
	for _, symbol := range []byte{0, 0x63, 0x64, 0xff} {
		if _, err := DecodeCMUPhones([]byte{symbol}); err == nil {
			t.Errorf("structural/unknown byte 0x%02x accepted as a phone", symbol)
		}
	}
}

func TestParsePhonePayloadRejectsMalformedPronunciation(t *testing.T) {
	for name, payload := range map[string][]byte{
		"direct missing terminator":            {0x01, 0x02},
		"alternative missing separator":        {0x02, 0x03, 0x00},
		"alternative missing final terminator": {0x02, '|', 0x01},
		"too many alternatives":                {0x02, '|', 0x01, 0xff, '|', 0x01, 0xff, '|', 0x01, 0xff, '|', 0x01, 0xff, '|', 0x01, 0xff, '|', 0x01, 0x00},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParsePhonePayload(payload); err == nil {
				t.Fatal("invalid phone payload accepted")
			}
		})
	}
}

func TestLoadEmbeddedDictionaryAndParseEveryPhonePayload(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	if dictionary.Len() != 228591 {
		t.Fatalf("embedded record count = %d, want 228591", dictionary.Len())
	}
	pronunciationRecords := 0
	codebook := Paul2013PhoneIDCodebook()
	row := make([]byte, Paul2013TokenResultRowSize)
	for key, record := range dictionary.records {
		parsed, err := ParsePhonePayload(record.Payload)
		if err != nil {
			t.Fatalf("parse %q: %v", key, err)
		}
		if _, err := parsed.ExpandPronunciations(codebook); err != nil {
			t.Fatalf("expand %q: %v", key, err)
		}
		rows, err := parsed.BuildPaul2013DictionaryPhoneRows(codebook)
		if err != nil {
			t.Fatalf("build parser phone rows for %q: %v", key, err)
		}
		if err := WritePaul2013DictionaryPhoneRow(row, 0, []byte("x"), rows); err != nil {
			t.Fatalf("write parser token-result row for %q: %v", key, err)
		}
		if len(parsed.Pronunciations) > 0 {
			pronunciationRecords++
		}
	}
	if pronunciationRecords == 0 {
		t.Fatal("embedded dictionary has no pronunciation records")
	}
}

func TestLookupSurfaceWithLocalKeyTablesWhenAssetsArePresent(t *testing.T) {
	root := filepath.Join("..", "..")
	dictionaryRoot := filepath.Join(root, "data-common", "dict-eng")
	dllPath := filepath.Join(root, "binary", "vt_pau.dll")
	dll, err := os.ReadFile(dllPath)
	if err != nil {
		t.Skip("local VoiceText DLL is unavailable")
	}
	const characterMapOffset = 0x7e388
	const pairTableOffset = 0x81568
	const pairTableStride = 3
	const pairCount = 127
	if len(dll) < characterMapOffset+256*2 || len(dll) < pairTableOffset+pairCount*pairTableStride {
		t.Fatal("local VoiceText DLL is shorter than the observed key-table ranges")
	}
	tables := EmbeddedKeyTables{}
	for i := range tables.CharacterMap {
		tables.CharacterMap[i] = dll[characterMapOffset+i*2]
	}
	tables.Pairs = make([][2]byte, pairCount)
	for i := range tables.Pairs {
		start := pairTableOffset + i*pairTableStride
		copy(tables.Pairs[i][:], dll[start:start+2])
	}
	provisionedTables := Paul2013EmbeddedKeyTables()
	if provisionedTables.CharacterMap != tables.CharacterMap || len(provisionedTables.Pairs) != len(tables.Pairs) {
		t.Fatal("provisioned Paul 2013 key tables do not match the documented DLL ranges")
	}
	for i := range tables.Pairs {
		if provisionedTables.Pairs[i] != tables.Pairs[i] {
			t.Fatalf("provisioned key pair %d = % x, want DLL value % x", i, provisionedTables.Pairs[i], tables.Pairs[i])
		}
	}

	dictionary, err := LoadEmbeddedDictionary(dictionaryRoot)
	if err != nil {
		t.Fatal(err)
	}
	pronunciation, ok, err := dictionary.ResolvePaul2013Surface([]byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("encoded Hello key was not found in the local embedded dictionary")
	}
	if len(pronunciation.Payload.Pronunciations) != 1 ||
		string(pronunciation.Payload.Pronunciations[0].Phone) != "\xd9\xce\x26" {
		t.Fatalf("Hello compact phone IDs = % x, want d9 ce 26", pronunciation.Payload.Pronunciations)
	}
	if len(pronunciation.Alternatives) != 1 || string(pronunciation.Alternatives[0].Symbols) != "\x22\x17\x2b\x30" {
		t.Fatalf("Hello internal symbols = %+v, want 22 17 2b 30", pronunciation.Alternatives)
	}
	phones, err := DecodeCMUPhones(pronunciation.Alternatives[0].Symbols)
	if err != nil {
		t.Fatalf("label Hello internal symbols: %v", err)
	}
	wantPhones := []CMUPhone{
		{Label: "HH"},
		{Label: "EH", Stress: 0, Vowel: true},
		{Label: "L"},
		{Label: "OW", Stress: 1, Vowel: true},
	}
	if len(phones) != len(wantPhones) {
		t.Fatalf("Hello CMU phone count = %d, want %d", len(phones), len(wantPhones))
	}
	for i := range wantPhones {
		if phones[i] != wantPhones[i] {
			t.Errorf("Hello CMU phone %d = %+v, want %+v", i, phones[i], wantPhones[i])
		}
	}
}

func TestTokenizeASCIISurfacesPreservesSeparators(t *testing.T) {
	tokens, err := tokenizeASCIISurfaces("  Hello, world!")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 || tokens[0].Surface != "Hello" || tokens[0].SeparatorBefore != "  " ||
		tokens[0].SeparatorAfter != ", " || tokens[1].Surface != "world" ||
		tokens[1].SeparatorBefore != ", " || tokens[1].SeparatorAfter != "!" {
		t.Fatalf("tokens = %+v", tokens)
	}
	if _, err := tokenizeASCIISurfaces("café"); err == nil {
		t.Fatal("non-ASCII token accepted")
	}
	if _, err := tokenizeASCIISurfaces("... "); err == nil {
		t.Fatal("punctuation-only input accepted")
	}
}

func TestLexiconFrontendResolvesKnownLocalText(t *testing.T) {
	root := filepath.Join("..", "..", "data-common", "dict-eng")
	if _, err := os.Stat(filepath.Join(root, "engttsdict_emb")); err != nil {
		t.Skip("local shared dictionary files are unavailable")
	}
	dictionary, err := LoadEmbeddedDictionary(root)
	if err != nil {
		t.Fatal(err)
	}
	analysis, err := (LexiconFrontend{Dictionary: dictionary}).AnalyzeText(context.Background(), "Hello.")
	if err != nil {
		t.Fatal(err)
	}
	tokens := analysis.Tokens
	if len(tokens) != 1 || tokens[0].Surface != "Hello" || tokens[0].SeparatorAfter != "." ||
		len(tokens[0].Alternatives) != 1 {
		t.Fatalf("resolved tokens = %+v", tokens)
	}
	if len(analysis.PronunciationFeatures) != len(tokens) {
		t.Fatalf("pronunciation feature rows = %d, want %d", len(analysis.PronunciationFeatures), len(tokens))
	}
	features := analysis.PronunciationFeatures[0]
	if features.Values[14] != 1 || features.Available != 0x5fff {
		t.Fatalf("Hello pronunciation features = %+v, want available mask 0x5fff and uppercase flag 1", features)
	}
	if !reflect.DeepEqual(features.MissingPositions(), []int{13}) {
		t.Fatalf("Hello missing pronunciation features = %v", features.MissingPositions())
	}
	wantCandidateCount := 0
	for _, alternative := range tokens[0].Alternatives {
		for _, pathGroup := range alternative.PathGroups {
			wantCandidateCount += len(pathGroup)
		}
	}
	candidates := analysis.PronunciationCandidates[0]
	if len(candidates) != wantCandidateCount {
		t.Fatalf("pronunciation candidates = %d, want one row for each of %d nonempty path groups", len(candidates), wantCandidateCount)
	}
	for _, candidate := range candidates {
		if candidate.Features.Available&(1<<13) == 0 {
			t.Errorf("candidate feature row omits path position 13: %+v", candidate.Features)
		}
		if len(candidate.Features.MissingPositions()) != 0 {
			t.Errorf("candidate feature row has unexpected missing positions: %v", candidate.Features.MissingPositions())
		}
		pathCode := tokens[0].Alternatives[candidate.AlternativeIndex].PathGroups[candidate.PathGroupIndex][candidate.PathCodeIndex]
		wantClass, _, err := Paul2013PronunciationPathClass([]byte{pathCode})
		if err != nil || candidate.Features.Values[13] != wantClass {
			t.Errorf("candidate path code %#x feature = %d, want %d (error %v)", pathCode, candidate.Features.Values[13], wantClass, err)
		}
	}
	want := []CMUPhone{{Label: "HH"}, {Label: "EH", Stress: 0, Vowel: true}, {Label: "L"}, {Label: "OW", Stress: 1, Vowel: true}}
	got := tokens[0].Alternatives[0].Phones
	if len(got) != len(want) {
		t.Fatalf("Hello phone count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Hello phone %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}
