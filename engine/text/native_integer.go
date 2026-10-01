package text

// ParsePaul2013NativeInteger ports FUN_100645ba's character-attribute scan.
// It skips leading bytes marked 0x08, consumes an optional ASCII sign, then
// accumulates consecutive bytes marked 0x04 with native 32-bit wraparound.
// hasDigits is false when the scan finds no digit; trailing bytes after the
// first non-digit are left unconsumed, as in the native helper.
func ParsePaul2013NativeInteger(source []byte, attributes [256]byte) (value int32, hasDigits bool) {
	source = cString(source)
	index := 0
	for index < len(source) && attributes[source[index]]&0x08 != 0 {
		index++
	}

	negative := false
	if index < len(source) && (source[index] == '-' || source[index] == '+') {
		negative = source[index] == '-'
		index++
	}
	for index < len(source) && attributes[source[index]]&0x04 != 0 {
		value = value*10 + int32(source[index]-'0')
		hasDigits = true
		index++
	}
	if negative {
		value = -value
	}
	return value, hasDigits
}
