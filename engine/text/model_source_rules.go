package text

import "fmt"

// Paul2013ModelSourceRuleDescriptor is one entry in the 2013 Paul DLL's
// FUN_1002e990 matcher table. Addresses are DLL virtual addresses observed in
// the read-only image; the descriptor does not assign semantics to opaque
// auxiliary pointers.
type Paul2013ModelSourceRuleDescriptor struct {
	Key                string
	HandlerAddress     uint32
	RuleDataAddress    uint32
	Flag               uint32
	HandlerDataAddress uint32
}

// Paul2013ModelSourceRuleDescriptors returns the 31 nonterminal entries
// observed at DLL VA 0x1007e888. The original table's following entry has a
// zero key length and terminates the scan.
func Paul2013ModelSourceRuleDescriptors() []Paul2013ModelSourceRuleDescriptor {
	return append([]Paul2013ModelSourceRuleDescriptor(nil), paul2013ModelSourceRules[:]...)
}

// BindPaul2013ModelSourceRules combines the observed key/order/metadata table
// with caller-supplied handler implementations. It fails closed if any native
// handler has no implementation, rather than treating an unknown tag as
// handled. RuleData contains the little-endian bytes of the two observed
// auxiliary pointers and flag; HandlerDataAddress remains available in the
// descriptor for callers that need its distinct pointer.
func BindPaul2013ModelSourceRules(
	handlers map[uint32]Paul2013ModelSourceRuleHandler,
) ([]Paul2013ModelSourceRule, error) {
	descriptors := Paul2013ModelSourceRuleDescriptors()
	rules := make([]Paul2013ModelSourceRule, len(descriptors))
	for index, descriptor := range descriptors {
		handler, ok := handlers[descriptor.HandlerAddress]
		if !ok || handler == nil {
			return nil, fmt.Errorf("Paul 2013 model-source handler at %#08x is not implemented", descriptor.HandlerAddress)
		}
		data := make([]byte, 12)
		putPaul2013ModelSourceUint32(data[0:4], descriptor.RuleDataAddress)
		putPaul2013ModelSourceUint32(data[4:8], descriptor.Flag)
		putPaul2013ModelSourceUint32(data[8:12], descriptor.HandlerDataAddress)
		rules[index] = Paul2013ModelSourceRule{Key: []byte(descriptor.Key), Data: data, Handler: handler}
	}
	return rules, nil
}

func putPaul2013ModelSourceUint32(destination []byte, value uint32) {
	destination[0] = byte(value)
	destination[1] = byte(value >> 8)
	destination[2] = byte(value >> 16)
	destination[3] = byte(value >> 24)
}

var paul2013ModelSourceRules = [...]Paul2013ModelSourceRuleDescriptor{
	{Key: "<vtml_break", HandlerAddress: 0x1002ef70, RuleDataAddress: 0x1007d970, Flag: 2, HandlerDataAddress: 0x1007d964},
	{Key: "<vtml_emotion", HandlerAddress: 0x10030c60, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_emotion", HandlerAddress: 0x10030c60, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_mark", HandlerAddress: 0x1002efd0, RuleDataAddress: 0x1007d970, Flag: 2, HandlerDataAddress: 0x1007d964},
	{Key: "<vtml_partofsp", HandlerAddress: 0x1002f3a0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_partofsp", HandlerAddress: 0x1002f3a0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_pause", HandlerAddress: 0x1002fb10, RuleDataAddress: 0x1007d970, Flag: 2, HandlerDataAddress: 0x1007d964},
	{Key: "<vtml_phoneme", HandlerAddress: 0x1002fcf0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_phoneme", HandlerAddress: 0x1002fcf0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_pitch", HandlerAddress: 0x10030b40, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_pitch", HandlerAddress: 0x10030b40, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_sayas", HandlerAddress: 0x10030c90, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_sayas", HandlerAddress: 0x10030c90, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_skip", HandlerAddress: 0x100313d0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1009f948},
	{Key: "</vtml_skip", HandlerAddress: 0x100313d0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1009f948},
	{Key: "<vtml_speed", HandlerAddress: 0x10030b70, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_speed", HandlerAddress: 0x10030b70, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_sub", HandlerAddress: 0x10031ce0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_sub", HandlerAddress: 0x10031ce0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vtml_volume", HandlerAddress: 0x10030ba0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</vtml_volume", HandlerAddress: 0x10030ba0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vt_pitch", HandlerAddress: 0x10030bd0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fc74},
	{Key: "</vt_pitch", HandlerAddress: 0x10030bd0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vt_speed", HandlerAddress: 0x10030c00, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fc74},
	{Key: "</vt_speed", HandlerAddress: 0x10030c00, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vt_volume", HandlerAddress: 0x10030c30, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fc74},
	{Key: "</vt_volume", HandlerAddress: 0x10030c30, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
	{Key: "<vt_pause", HandlerAddress: 0x1002fb40, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fc74},
	{Key: "<vt_break", HandlerAddress: 0x1002efa0, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fc74},
	{Key: "<say-as", HandlerAddress: 0x10031d10, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007d964},
	{Key: "</say-as", HandlerAddress: 0x10031d10, RuleDataAddress: 0x1007d464, Flag: 1, HandlerDataAddress: 0x1007fd84},
}
