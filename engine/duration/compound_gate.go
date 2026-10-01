package duration

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

// Paul2013CompoundDictionaryGateResult reports FUN_10003b10's direct payload
// flag gate without assigning meaning to the flag bits.
type Paul2013CompoundDictionaryGateResult struct {
	LookupFound bool
	PayloadFlag byte
	Eligible    bool
}

// LookupPaul2013CompoundDictionaryGate connects FUN_10003a70's mode-0
// embedded-key lookup to FUN_10003b10's payload-flag predicate. A present
// payload passes when bit 2 is set or when masking with 0xf8 yields 0x08.
func (engine *Engine) LookupPaul2013CompoundDictionaryGate(
	ctx context.Context,
	surface []byte,
) (Paul2013CompoundDictionaryGateResult, error) {
	if engine == nil || engine.dictionary == nil {
		return Paul2013CompoundDictionaryGateResult{}, errors.New("Paul 2013 embedded context dictionary is nil or unloaded")
	}
	if ctx == nil {
		return Paul2013CompoundDictionaryGateResult{}, errors.New("compound dictionary gate has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013CompoundDictionaryGateResult{}, err
	}
	if nul := bytes.IndexByte(surface, 0); nul >= 0 {
		surface = surface[:nul]
	}
	resolved, found, err := engine.dictionary.ResolvePaul2013Surface(surface)
	if err != nil {
		return Paul2013CompoundDictionaryGateResult{}, fmt.Errorf("lookup compound component %q: %w", surface, err)
	}
	result := Paul2013CompoundDictionaryGateResult{LookupFound: found}
	if found {
		result.PayloadFlag = resolved.Payload.Flags
		result.Eligible = result.PayloadFlag&0x04 != 0 || result.PayloadFlag&0xf8 == 0x08
	}
	if err := ctx.Err(); err != nil {
		return Paul2013CompoundDictionaryGateResult{}, err
	}
	return result, nil
}

// Paul2013CompoundCharacterGate ports FUN_1000ffd0 for the supported ASCII
// surface. It returns false for an empty surface and rejects a surface if any
// character has a nonzero 0xc0 class with bit 0x40 set.
func Paul2013CompoundCharacterGate(surface []byte) (bool, error) {
	if nul := bytes.IndexByte(surface, 0); nul >= 0 {
		surface = surface[:nul]
	}
	if len(surface) == 0 {
		return false, nil
	}
	attributes := text.Paul2013UnsignedCharacterAttributeTable()
	for index, value := range surface {
		if value >= 0x80 {
			return false, fmt.Errorf("compound character gate byte 0x%02x at offset %d is outside the supported ASCII input", value, index)
		}
		attribute := attributes[value]
		if attribute&0xc0 != 0 && attribute&0x40 != 0 {
			return false, nil
		}
	}
	return true, nil
}
