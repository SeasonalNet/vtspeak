package text

import (
	"bytes"
	"context"
	"fmt"
)

// Paul2013ScannerLookupIndex projects the native index object: signed count
// at +0, lower/upper row arrays at +4/+8, and 20-byte rows whose key pointer
// is at +0xc. Each array cell selects the inclusive range for length index+1.
type Paul2013ScannerLookupIndex struct {
	Lower, Upper []int32
	Keys         [][]byte
	nativeKeys   map[int32][]byte // Sparse copied native ranges; inactive gaps stay unread.
}

// Lookup ports FUN_1005e010 (exact prefixes) and FUN_1005e1a0 (mapped
// prefixes). The boundary gate remains explicit and is invoked in the native
// endpoint/midpoint order. Null native indexes are represented by nil.
func (table *Paul2013ScannerLookupIndex) Lookup(ctx context.Context, source []byte, mapped bool, gate func(context.Context, []byte, int32) (int16, error)) (int32, error) {
	if table == nil {
		return -1, nil
	}
	if len(table.Lower) != len(table.Upper) {
		return -1, fmt.Errorf("scanner lookup bound array counts differ")
	}
	source = cString(source)
	weights := Paul2013ContextCharacterWeights()
	compare := func(index int32, length int) (int, error) {
		var key []byte
		if table.nativeKeys != nil {
			var exists bool
			key, exists = table.nativeKeys[index]
			if !exists {
				return 0, fmt.Errorf("scanner lookup row %d outside resolved ranges", index)
			}
		} else {
			if index < 0 || int64(index) >= int64(len(table.Keys)) {
				return 0, fmt.Errorf("scanner lookup row %d outside index", index)
			}
			key = table.Keys[index]
		}
		key = cString(key)
		for position := 0; position < length; position++ {
			left, right := byte(0), byte(0)
			if position < len(source) {
				left = source[position]
			}
			if position < len(key) {
				right = key[position]
			}
			difference := int(left) - int(right)
			if mapped {
				difference = int(weights[left]) - int(weights[right])
			}
			if difference != 0 || left == 0 {
				return difference, nil
			}
		}
		return 0, nil
	}
	accept := func(index int32, length int) (bool, error) {
		comparison, err := compare(index, length)
		if err != nil || comparison != 0 {
			return false, err
		}
		if gate == nil {
			return false, fmt.Errorf("scanner lookup boundary gate is unavailable")
		}
		value, err := gate(ctx, source, int32(length))
		return value != 0, err
	}
	for length := len(table.Lower); length > 0; length-- {
		if err := ctx.Err(); err != nil {
			return -1, err
		}
		low, high := table.Lower[length-1], table.Upper[length-1]
		if low < 0 || high < 0 || len(source) < length {
			continue
		}
		if high < low {
			return -1, fmt.Errorf("scanner lookup range is reversed")
		}
		for {
			if err := ctx.Err(); err != nil {
				return -1, err
			}
			if accepted, err := accept(low, length); err != nil || accepted {
				if err != nil {
					return -1, err
				}
				return low, nil
			}
			if accepted, err := accept(high, length); err != nil || accepted {
				if err != nil {
					return -1, err
				}
				return high, nil
			}
			if high-low <= 1 {
				break
			}
			middle := low + (high-low)/2
			comparison, err := compare(middle, length)
			if err != nil {
				return -1, err
			}
			if comparison == 0 {
				if gate == nil {
					return -1, fmt.Errorf("scanner lookup boundary gate is unavailable")
				}
				value, err := gate(ctx, source, int32(length))
				if err != nil {
					return -1, err
				}
				if value != 0 {
					return middle, nil
				}
			}
			if comparison < 0 {
				high = middle
			} else {
				low = middle
			}
		}
	}
	return -1, nil
}

// Paul2013ScannerApostropheBoundary ports FUN_1005f9b0. Its exact suffix
// set differs from the general contraction normalizer: s/d/m/ve/em/re/ll.
func Paul2013ScannerApostropheBoundary(source []byte) bool {
	source = cString(source)
	if len(source) == 0 || source[0] != '\'' {
		return false
	}
	weights := Paul2013ContextCharacterWeights()
	for _, suffix := range []string{"s", "d", "m", "ve", "em", "re", "ll"} {
		if ComparePaul2013MappedCString(source[1:], []byte(suffix), weights) == 0 {
			return true
		}
	}
	return false
}

// CheckPaul2013ScannerLookupBoundary ports FUN_1005f840's complete return
// decisions. Its model argument is unused; the third argument's special dot
// check is repeated by the ordinary dot check, so both have identical output.
// Internal scanning requires mode 0x17, pointer zero, coordinate zero-origin.
func CheckPaul2013ScannerLookupBoundary(ctx context.Context, source []byte, length int32, scan Paul2013ModelScanner) (int16, error) {
	source = append(append([]byte(nil), cString(source)...), 0)
	if length < 0 || int64(length) >= int64(len(source)) {
		return 0, fmt.Errorf("scanner lookup boundary length outside source")
	}
	value := source[length]
	if value == ' ' || value == '\t' || value == '\n' || value == '\r' {
		return 1, nil
	}
	if scan == nil {
		return 0, fmt.Errorf("scanner lookup boundary mode-0x17 scanner is unavailable")
	}
	position := int32(0)
	var token Paul2013ModelScannerResult
	for {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		var err error
		token, err = scan(ctx, source[position:], position, 0x17, 0)
		if err != nil {
			return 0, err
		}
		if token.Advance == 0 || token.Status == 9 || token.TextLength == 0 {
			return 1, nil
		}
		if token.Advance < 0 || int64(token.Advance) > int64(len(source)-1)-int64(position) {
			return 0, fmt.Errorf("boundary scanner advance outside source")
		}
		position += token.Advance
		if position == length {
			return 1, nil
		}
		if position > length {
			break
		}
	}
	text := cString(token.Text)
	attributes := Paul2013ExceptionCharacterAttributes()
	for _, punctuation := range []byte{'.', '\'', '-'} {
		index := bytes.IndexByte(text, punctuation)
		if index <= 0 {
			continue
		}
		accepted := false
		if punctuation == '\'' {
			accepted = Paul2013ScannerApostropheBoundary(text[index:])
		} else if index+1 < len(text) {
			accepted = attributes[text[index+1]]&0xc0 != 0
		}
		if accepted && position-int32(len(text)-index) == length {
			return 1, nil
		}
	}
	return 0, nil
}

// NewPaul2013ScannerIndexedTokenLookup binds native model addresses to their
// exact/mapped indexes and supplies the recovered boundary body.
func NewPaul2013ScannerIndexedTokenLookup(resolve func(uint32) (*Paul2013ScannerLookupIndex, *Paul2013ScannerLookupIndex, error), scan Paul2013ModelScanner) Paul2013ScannerTokenLookup {
	return func(ctx context.Context, address, pointer uint32, source []byte) (int32, error) {
		if resolve == nil {
			return -1, fmt.Errorf("scanner model index resolver is unavailable")
		}
		exact, mapped, err := resolve(pointer)
		if err != nil {
			return -1, err
		}
		gate := func(ctx context.Context, source []byte, length int32) (int16, error) {
			return CheckPaul2013ScannerLookupBoundary(ctx, source, length, scan)
		}
		switch address {
		case 0x1005e010:
			return exact.Lookup(ctx, source, false, gate)
		case 0x1005e1a0:
			return mapped.Lookup(ctx, source, true, gate)
		default:
			return -1, fmt.Errorf("unsupported scanner lookup address %#x", address)
		}
	}
}
