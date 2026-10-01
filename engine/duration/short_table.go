package duration

// paul2013TailClassWindow stores the zero/nonzero result of the native
// 16-bit read at DAT_10077d5e + signedIndex*2 for signedIndex -128 through
// 127. The first 128 entries are the exact preceding memory window; the
// final 128 cover the table and following data. FUN_1000a140 only compares
// these reads with zero. Values are from binary/vt_pau.dll file offsets
// 0x77c5e through 0x77e5d.
var paul2013TailClassWindow = [16][16]byte{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1},
	{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0},
	{0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0},
}

func paul2013TailClassForSourceByte(sourceByte byte, characterMap [256]byte) byte {
	index := int(characterMap[sourceByte])
	if sourceByte >= 0x80 {
		// The u16 map's signed-index prefix contains 0xff80..0xffff, which
		// sign-extend to the original source byte's -128..-1 value.
		index = int(int8(sourceByte))
	}
	windowIndex := index + 128
	return paul2013TailClassWindow[windowIndex/16][windowIndex%16]
}
