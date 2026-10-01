package text

import (
	"bytes"
)

// Paul2013FUN10010010 ports FUN_10010010. It returns
// true only when the C string is nonempty and every byte's entry in
// DAT_1007e188 has both bits 0xc0 and 0x80 set. The native helper indexes the
// table with signed chars; high-bit inputs use the verified zero-valued
// prefix before DAT_1007e188.
func Paul2013FUN10010010(source []byte) (bool, error) {
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	if len(source) == 0 {
		return false, nil
	}
	attributes := Paul2013ExceptionCharacterAttributes()
	for _, value := range source {
		attribute := attributes[value]
		if attribute&0xc0 == 0 || attribute&0x80 == 0 {
			return false, nil
		}
	}
	return true, nil
}

// Paul2013FUN100100D0 ports FUN_100100d0's vowel count for C strings.
// Each source byte is mapped through the low byte of the 256-entry u16 table
// at 0x1007e388; mapped values equal to lowercase a, e, i, o, or u increment
// the count. For high-bit input, the signed index selects words 0xff80 through
// 0xffff immediately before the table; none can equal an ASCII vowel.
func Paul2013FUN100100D0(source []byte) (int, error) {
	if end := bytes.IndexByte(source, 0); end >= 0 {
		source = source[:end]
	}
	characterMap := Paul2013EmbeddedKeyTables().CharacterMap
	vowelCount := 0
	for _, value := range source {
		switch characterMap[value] {
		case 'a', 'e', 'i', 'o', 'u':
			vowelCount++
		}
	}
	return vowelCount, nil
}
