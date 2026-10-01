package text

import (
	"bytes"
	"fmt"
)

// Paul2013ParserRowTypeMarker returns the fixed marker selected by
// FUN_1000d190 for a source-row +0x2c value. The numeric codes retain their
// native values; this mapping does not infer the parser branch that produced
// one.
func Paul2013ParserRowTypeMarker(rowType uint32) ([]byte, bool) {
	switch rowType {
	case 2:
		return []byte{'.'}, true
	case 3:
		return []byte{'?'}, true
	case 4:
		return []byte{'!'}, true
	case 5, 12:
		return []byte{','}, true
	case 11:
		return []byte{'-'}, true
	default:
		return nil, false
	}
}

// Paul2013ParserRowMarkerInput contains the values consumed by the marker
// append branch in FUN_1000d190. The native header value is at the pointer
// expression param_1[-9]; its higher-level meaning is not assigned here.
type Paul2013ParserRowMarkerInput struct {
	HeaderValue int32
	RowType     uint32
	Text        []byte
}

// AppendPaul2013ParserRowTypeMarker ports FUN_1000d190's fixed-marker append.
// It acts only when HeaderValue < 0x1d, the row type is nonzero and recognized,
// and the row text is nonempty. C-string termination and the native 32-byte
// local buffer are preserved; overlong inputs fail closed instead of
// reproducing a stack overwrite.
func AppendPaul2013ParserRowTypeMarker(
	input Paul2013ParserRowMarkerInput,
) ([]byte, bool, error) {
	text := input.Text
	if nul := bytes.IndexByte(text, 0); nul >= 0 {
		text = text[:nul]
	}
	marker, recognized := Paul2013ParserRowTypeMarker(input.RowType)
	if input.HeaderValue >= 0x1d || input.RowType == 0 || !recognized || len(text) == 0 {
		return append([]byte(nil), text...), false, nil
	}
	if len(text)+len(marker) > 31 {
		return nil, false, fmt.Errorf("parser-row marker append would use %d bytes in the native 32-byte local buffer", len(text)+len(marker)+1)
	}
	result := make([]byte, len(text)+len(marker))
	copy(result, text)
	copy(result[len(text):], marker)
	return result, true, nil
}
