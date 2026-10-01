package text

import "bytes"

// ComparePaul2013MappedCString ports FUN_1001c2c0. Each input byte indexes
// the DLL's 256-entry signed-short weight table, and the first differing
// weight determines the result. The table is explicit because the caller
// must provide the exact table for the DLL generation being modeled.
func ComparePaul2013MappedCString(left, right []byte, weights [256]int16) int {
	if len(left) > 0 && len(right) > 0 && &left[0] == &right[0] {
		return 0
	}
	left = cString(left)
	right = cString(right)
	for index := 0; ; index++ {
		leftByte, rightByte := byte(0), byte(0)
		if index < len(left) {
			leftByte = left[index]
		}
		if index < len(right) {
			rightByte = right[index]
		}
		if difference := int(weights[leftByte]) - int(weights[rightByte]); difference != 0 {
			return difference
		}
		if leftByte == 0 {
			return 0
		}
	}
}

// FindPaul2013SortedCString ports the binary-search branch of FUN_100560a0
// for a pointer table already sorted under FUN_1001c2c0. It returns the
// native matching row index; a miss returns -1. Empty C strings are rejected
// by the native lookup routine.
func FindPaul2013SortedCString(rows [][]byte, key []byte, weights [256]int16) int {
	key = cString(key)
	if len(key) == 0 || len(rows) == 0 {
		return -1
	}
	low, high := 0, len(rows)-1
	for low <= high {
		middle := low + (high-low)/2
		comparison := ComparePaul2013MappedCString(rows[middle], key, weights)
		switch {
		case comparison == 0:
			return middle
		case comparison < 0:
			low = middle + 1
		default:
			high = middle - 1
		}
	}
	return -1
}

// FindPaul2013CString ports the linear strcmp branch of FUN_100560a0. The
// mode selector 0x49 used by the context-classifier helpers chooses mapped
// comparison instead; this helper is for callers that observed strcmp mode.
func FindPaul2013CString(rows [][]byte, key []byte) int {
	key = cString(key)
	if len(key) == 0 {
		return -1
	}
	for index, row := range rows {
		if bytes.Equal(cString(row), key) {
			return index
		}
	}
	return -1
}
