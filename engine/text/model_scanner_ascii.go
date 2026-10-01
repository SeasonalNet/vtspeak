package text

import (
	"bytes"
	"context"
	"fmt"
)

var paul2013ScannerASCIICosts = [...]int16{1, 2, 2, 2, 1, 2, 2, 2, 1, 2, 2, 2, 2, 2, 1, 2, 3, 2, 2, 2, 2, 2, 7, 3, 2, 2}

// Paul2013ScannerTokenLookup supplies the two post-scan model lookups at
// FUN_1005e010 and FUN_1005e1a0, including internal hyphen join gates.
// A nonnegative post-token result sets raw short +0x1a.
type Paul2013ScannerTokenLookup func(context.Context, uint32, uint32, []byte) (int32, error)

// ScanPaul2013ModelScannerASCII ports ASCII letters, numbers and punctuation
// in modes 0/1/0x17, mode-0x12 punctuation and recovered ASCII word
// branches across the native mode dispatch, including its default branch.
// The standalone entry point uses no model pointer; the composed scanner
// supplies mode-8 join lookup gates. Non-ASCII and other unrecovered paths
// return handled=false for the original source to be scanned elsewhere.
func ScanPaul2013ModelScannerASCII(source []byte, offset, mode int32) (Paul2013ModelScannerResult, bool, error) {
	return scanPaul2013ModelScannerASCII(context.Background(), source, offset, mode, 0, nil)
}

func scanPaul2013ModelScannerASCII(ctx context.Context, source []byte, offset, mode int32, pointer uint32, lookup Paul2013ScannerTokenLookup) (result Paul2013ModelScannerResult, handled bool, err error) {
	if err := ctx.Err(); err != nil {
		return result, false, err
	}

	end := bytes.IndexByte(source, 0)
	if end < 0 {
		return result, false, fmt.Errorf("ASCII scanner source is not NUL-terminated")
	}
	prefix := ScanPaul2013ModelParserWhitespacePrefix(source[:end+1])
	position := prefix.ConsumedBytes
	if prefix.SingleLineFeedMarker && prefix.LineFeeds < 2 {
		result.PrefixFlag = 1
	}
	space := func(value byte) bool { return value == ' ' || value == '\t' || value == '\n' || value == '\r' }
	letter := func(value byte) bool { return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z' }
	structural := []byte("~`@#$%&^_|,: '+-*/={}[]()<>\\\"")
	terminal := []byte(".?!;")
	for position < end {
		value := source[position]
		if space(value) {
			extra := ScanPaul2013ModelParserWhitespacePrefix(source[position : end+1])
			position += extra.ConsumedBytes
			prefix.LineFeeds += extra.LineFeeds
			if prefix.LineFeeds >= 2 {
				break
			}
		} else if letter(value) || value >= '0' && value <= '9' || bytes.IndexByte(structural, value) >= 0 || bytes.IndexByte(terminal, value) >= 0 || value >= 0x80 {
			break
		} else {
			position++
		}
	}
	result.Field0 = int32(position)
	result.First = uint32(offset) + uint32(position)
	skipPostLookup := false
	finish := func(text []byte, status int32, consumed int) (Paul2013ModelScannerResult, bool, error) {
		result.Text = append([]byte(nil), text...)
		result.TextLength = int32(len(text))
		result.Status = status
		result.Advance = int32(consumed)
		result.Field14 = int32(consumed)
		result.Second = uint32(offset) + uint32(consumed)
		if pointer != 0 && !skipPostLookup {
			found, err := lookupPaul2013ScannerToken(ctx, pointer, result.Text, lookup)
			if err != nil {
				return result, true, err
			}
			if found {
				result.Field18 |= 1 << 16
			}
		}
		return result, true, nil
	}
	if prefix.LineFeeds >= 2 {
		skipPostLookup = true
		return finish(nil, 8, position)
	}
	if position == end {
		skipPostLookup = true
		return finish(nil, 9, position)
	}
	start := position
	value := source[position]
	if mode != 0 && mode != 1 && mode != 0x12 && mode != 0x17 && !letter(value) && !(value >= '0' && value <= '9') {
		return Paul2013ModelScannerResult{}, false, nil
	}
	if value >= 0x80 {
		return Paul2013ModelScannerResult{}, false, nil
	}
	if letter(value) || value >= '0' && value <= '9' {
		cost := 0
		numericLetterUsed := false
		status := int32(1)
		for {
			if position < end && source[position] >= '0' && source[position] <= '9' {
				var resume bool
				position, status, resume = scanPaul2013ASCIINumber(source, start, position, end, mode, &numericLetterUsed)
				if resume {
					continue
				}
				break
			}
			status = 1
			letterStart := position
			for position < end && letter(source[position]) {
				index := source[position]
				if index >= 'a' {
					index -= 'a'
				} else {
					index -= 'A'
				}
				nextCost := cost + 1 + int(paul2013ScannerASCIICosts[index])
				if nextCost > 64 || position-start > 28 {
					if boundary := bytes.LastIndexAny(source[start:position], "-."); boundary >= 0 {
						position = start + boundary
					}
					return finish(source[start:position], 1, position)
				}
				cost = nextCost
				position++
			}
			if position-start > 28 || position == letterStart {
				break
			}
			if consumed := scanPaul2013ASCIIWordJoin(source, start, position, end, mode); consumed != 0 {
				position += consumed
				continue
			}
			if mode != 8 && mode != 0x17 {
				position += scanPaul2013ASCIIWordEnd(source, start, position, end, mode)
				break
			}
			join := position+1 < end && source[position] == '-' && letter(source[position+1])
			if mode == 8 && position+1 < end && source[position] == '.' && letter(source[position+1]) {
				join = bytes.IndexByte(source[start:position], '.') >= 0 || source[start] < 'a' || source[start] > 'z' || source[position+1] < 'A' || source[position+1] > 'Z'
			}
			if join {
				if mode == 8 && pointer != 0 {
					found, err := lookupPaul2013ScannerToken(ctx, pointer, source[position+1:end+1], lookup)
					if err != nil {
						return result, false, err
					}
					if found {
						skipPostLookup = true
						return finish(source[start:position], 1, position)
					}
				}
				position++
				continue
			}
			if position < end && source[position] == '\'' && bytes.IndexByte(source[start:position], '\'') < 0 {
				if suffix, accepted := ScanPaul2013ScannerApostropheSuffix(source[position-1], source[position:end+1]); accepted {
					position += len(suffix)
					if mode == 8 && position+1 < end && source[position] == '-' && letter(source[position+1]) {
						position++
						continue
					}
				} else if mode == 8 && bytes.EqualFold(source[start:position], []byte("wi")) {
					position++
				}
			} else if position+1 < end && source[position] == '.' && source[position+1] == '\'' {
				if suffix, accepted := ScanPaul2013ScannerApostropheSuffix('.', source[position+1:end+1]); accepted {
					position += 1 + len(suffix)
				}
			}
			break
		}
		if position-start >= 32 {
			return Paul2013ModelScannerResult{}, false, fmt.Errorf("joined scanner token exceeds native local buffer")
		}
		// A non-ASCII continuation can enter native expansion or join paths.
		if position < end && (source[position] >= 0x80 || position+1 < end && source[position+1] >= 0x80 && (source[position] == '.' || source[position] == '-' || source[position] == '\'' || mode == 0x1b && (source[position] == '&' || source[position] == '/'))) {
			return Paul2013ModelScannerResult{}, false, nil
		}
		return finish(source[start:position], status, position)
	}
	if (mode == 0 || mode == 1) && value == '.' && position+1 < end && source[position+1] >= '0' && source[position+1] <= '9' {
		if offset != 0 && start == 0 {
			return Paul2013ModelScannerResult{}, false, nil
		}
		used := false
		position, status, _ := scanPaul2013ASCIINumber(source, start, start, end, mode, &used)
		return finish(source[start:position], status, position)
	}
	if value == '<' && position+1 < end && source[position+1] >= 0x80 {
		return Paul2013ModelScannerResult{}, false, nil
	}
	if value == '\\' {
		return finish(source[position:position+1], 3, position+1)
	}
	position++
	for position < end && position-start < 29 && source[position] == value {
		position++
	}
	status := int32(3)
	if bytes.IndexByte(terminal, value) >= 0 {
		status = 7
		if position-start > 1 {
			status = 8
			if value == '.' {
				status = 3
			}
		} else if mode == 0x12 {
			status = 3
		}
	}
	return finish(source[start:position], status, position)
}

// NewPaul2013ModelScannerWithASCII composes recovered terminal prefixes and
// ASCII branches, preserving explicit fallback for all remaining token paths.
func NewPaul2013ModelScannerWithASCII(fallback Paul2013ModelScanner, lookup Paul2013ScannerTokenLookup) Paul2013ModelScanner {
	return func(ctx context.Context, source []byte, offset, mode int32, pointer uint32) (Paul2013ModelScannerResult, error) {
		if err := ctx.Err(); err != nil {
			return Paul2013ModelScannerResult{}, err
		}
		if result, handled, err := ScanPaul2013ModelScannerTerminalPrefix(source, offset); handled || err != nil {
			return result, err
		}
		result, handled, err := scanPaul2013ModelScannerASCII(ctx, source, offset, mode, pointer, lookup)
		if err != nil {
			return result, err
		}
		if !handled {
			if fallback == nil {
				return result, fmt.Errorf("general scanner path is unavailable for mode %d", mode)
			}
			return fallback(ctx, source, offset, mode, pointer)
		}
		return result, nil
	}
}

// lookupPaul2013ScannerToken preserves native exact-before-mapped short circuiting.
func lookupPaul2013ScannerToken(ctx context.Context, pointer uint32, source []byte, lookup Paul2013ScannerTokenLookup) (bool, error) {
	if lookup == nil {
		return false, fmt.Errorf("scanner model lookup is unavailable")
	}
	for _, address := range []uint32{0x1005e010, 0x1005e1a0} {
		if err := ctx.Err(); err != nil {
			return false, err
		}
		value, err := lookup(ctx, address, pointer, source)
		if err != nil {
			return false, err
		}
		if value >= 0 {
			return true, nil
		}
	}
	return false, nil
}
