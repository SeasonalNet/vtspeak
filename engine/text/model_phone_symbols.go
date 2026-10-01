package text

import (
	"fmt"

	"vtspeak/engine/internal/paul2013tables"
)

// DecodePaul2013ModelPhoneSymbols uses FUN_10012c70's native identity and
// stress tables rather than requiring a canonical pronunciation byte. Model
// symbols such as M alias a supported identity without changing their source
// record bytes. Entries without a supported identity remain explicit errors.
func DecodePaul2013ModelPhoneSymbols(symbols []byte) ([]CMUPhone, error) {
	phones := make([]CMUPhone, len(symbols))
	for index, symbol := range symbols {
		if symbol >= 0x80 {
			return nil, fmt.Errorf("model phone symbol %#x at %d uses unsupported negative native index", symbol, index)
		}
		ordinal := int(paul2013tables.ByteTable1007B6C0(symbol))
		if ordinal < 1 || ordinal > len(paul2013TreePhoneOrder) {
			return nil, fmt.Errorf("model phone symbol %#x at %d maps to unsupported identity %d", symbol, index, ordinal)
		}
		label := paul2013TreePhoneOrder[ordinal-1]
		_, vowel := paul2013TreeVowels[label]
		phone := CMUPhone{Label: label, Vowel: vowel, Stress: paul2013tables.ByteTable1007BAA8(symbol)}
		if _, err := phone.TreeFeatures(); err != nil {
			return nil, fmt.Errorf("model phone symbol %#x at %d: %w", symbol, index, err)
		}
		phones[index] = phone
	}
	return phones, nil
}
