package text

import "testing"

func TestPaul2013ModelSourceRuleDescriptorsMatchObservedTable(t *testing.T) {
	wantKeys := []string{
		"<vtml_break", "<vtml_emotion", "</vtml_emotion", "<vtml_mark",
		"<vtml_partofsp", "</vtml_partofsp", "<vtml_pause", "<vtml_phoneme",
		"</vtml_phoneme", "<vtml_pitch", "</vtml_pitch", "<vtml_sayas",
		"</vtml_sayas", "<vtml_skip", "</vtml_skip", "<vtml_speed",
		"</vtml_speed", "<vtml_sub", "</vtml_sub", "<vtml_volume",
		"</vtml_volume", "<vt_pitch", "</vt_pitch", "<vt_speed",
		"</vt_speed", "<vt_volume", "</vt_volume", "<vt_pause", "<vt_break",
		"<say-as", "</say-as",
	}
	wantHandlers := []uint32{
		0x1002ef70, 0x10030c60, 0x10030c60, 0x1002efd0, 0x1002f3a0,
		0x1002f3a0, 0x1002fb10, 0x1002fcf0, 0x1002fcf0, 0x10030b40,
		0x10030b40, 0x10030c90, 0x10030c90, 0x100313d0, 0x100313d0,
		0x10030b70, 0x10030b70, 0x10031ce0, 0x10031ce0, 0x10030ba0,
		0x10030ba0, 0x10030bd0, 0x10030bd0, 0x10030c00, 0x10030c00,
		0x10030c30, 0x10030c30, 0x1002fb40, 0x1002efa0, 0x10031d10,
		0x10031d10,
	}
	got := Paul2013ModelSourceRuleDescriptors()
	if len(got) != len(wantKeys) {
		t.Fatalf("descriptor count = %d, want %d", len(got), len(wantKeys))
	}
	for index, want := range wantKeys {
		if got[index].Key != want {
			t.Errorf("descriptor %d key = %q, want %q", index, got[index].Key, want)
		}
		if got[index].HandlerAddress != wantHandlers[index] {
			t.Errorf("descriptor %d handler = %#08x, want %#08x", index, got[index].HandlerAddress, wantHandlers[index])
		}
		if got[index].HandlerAddress == 0 || got[index].RuleDataAddress == 0 || got[index].HandlerDataAddress == 0 {
			t.Errorf("descriptor %d has a null observed table field: %+v", index, got[index])
		}
	}
	if got[0].RuleDataAddress != 0x1007d970 || got[0].Flag != 2 || got[0].HandlerDataAddress != 0x1007d964 ||
		got[2].HandlerDataAddress != 0x1007fd84 || got[27].RuleDataAddress != 0x1007d464 ||
		got[27].Flag != 1 || got[27].HandlerDataAddress != 0x1007fc74 ||
		got[13].HandlerDataAddress != 0x1009f948 {
		t.Fatalf("table metadata differs at representative records: first=%+v close=%+v vt-pause=%+v skip=%+v",
			got[0], got[2], got[27], got[13])
	}
	got[0].Key = "caller mutation"
	if Paul2013ModelSourceRuleDescriptors()[0].Key != wantKeys[0] {
		t.Fatal("descriptor result aliases the package table")
	}
}

func TestBindPaul2013ModelSourceRulesRequiresEveryHandler(t *testing.T) {
	handler := func([]byte, *int, []byte, *int, Paul2013ModelSourceRuleContext) (bool, error) {
		return false, nil
	}
	if _, err := BindPaul2013ModelSourceRules(nil); err == nil {
		t.Fatal("missing native handler implementation was accepted")
	}
	handlers := make(map[uint32]Paul2013ModelSourceRuleHandler)
	for _, descriptor := range Paul2013ModelSourceRuleDescriptors() {
		handlers[descriptor.HandlerAddress] = handler
	}
	rules, err := BindPaul2013ModelSourceRules(handlers)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 31 || string(rules[2].Key) != "</vtml_emotion" || rules[0].Handler == nil {
		t.Fatalf("bound rule table = %+v (count %d)", rules[:3], len(rules))
	}
	if len(rules[0].Data) != 12 || rules[0].Data[4] != 2 || rules[0].Data[5] != 0 {
		t.Fatalf("bound rule metadata bytes = % x", rules[0].Data)
	}
}
