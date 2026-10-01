package text

// Paul2013UserDictionarySourceNormalization retains the private source
// normalizer's signed return code and the bytes it wrote. On a pair rejection,
// Output can contain an unterminated prefix; bytes after that prefix are not
// represented because the native wrapper leaves them uninitialized.
type Paul2013UserDictionarySourceNormalization struct {
	ReturnCode    int32
	Output        []byte
	NULTerminated bool
}

// NormalizePaul2013UserDictionarySource ports the behavior of
// FUN_1005f2e0, called by VT_CheckUserDict_SourceNorm_ENG. It trims only
// leading/trailing space, TAB, LF, and CR; copies all other bytes without
// rewriting; enforces the observed 49-byte output limit; and rejects the
// recovered byte-pair ranges. The public export discards this result, but
// retaining the helper result makes its evidence-backed mechanics callable.
func NormalizePaul2013UserDictionarySource(source []byte) Paul2013UserDictionarySourceNormalization {
	if nul := indexByte(source, 0); nul >= 0 {
		source = source[:nul]
	}
	start, end := 0, len(source)
	for start < end && paul2013UserDictionaryTrimByte(source[start]) {
		start++
	}
	for end > start && paul2013UserDictionaryTrimByte(source[end-1]) {
		end--
	}
	source = source[start:end]
	if len(source) == 0 {
		return Paul2013UserDictionarySourceNormalization{ReturnCode: -1, NULTerminated: true}
	}
	if len(source) > 49 {
		return Paul2013UserDictionarySourceNormalization{ReturnCode: -5, NULTerminated: true}
	}

	output := make([]byte, 0, len(source))
	for index := 0; index+1 < len(source); index++ {
		if paul2013UserDictionaryRejectedPair(source[index], source[index+1]) {
			return Paul2013UserDictionarySourceNormalization{
				ReturnCode: -3, Output: output, NULTerminated: len(output) == 0,
			}
		}
		output = append(output, source[index])
	}
	output = append(output, source[len(source)-1])
	return Paul2013UserDictionarySourceNormalization{
		ReturnCode: int32(len(output)), Output: output, NULTerminated: true,
	}
}

func paul2013UserDictionaryTrimByte(value byte) bool {
	return value == ' ' || value == '\t' || value == '\n' || value == '\r'
}

func paul2013UserDictionaryRejectedPair(first, second byte) bool {
	return (first >= 0xa1 && first <= 0xad && second >= 0xa1 && second <= 0xfe) ||
		(first == 0xae && second >= 0xa1 && second <= 0xc2) ||
		(first == 0xfd && second == 0xfe)
}
