package text

import "fmt"

// Paul2013C3A0ComponentStop reports why FUN_1000c860 stopped scanning one
// component from the source buffer.
type Paul2013C3A0ComponentStop uint8

const (
	Paul2013C3A0ComponentAtLimit Paul2013C3A0ComponentStop = iota
	Paul2013C3A0ComponentAtNUL
	Paul2013C3A0ComponentAtCharacterBoundary
	Paul2013C3A0ComponentAfterApostropheS
	Paul2013C3A0ComponentAfterApostrophe
)

// Paul2013C3A0Component is one component emitted by FUN_1000c860. Offsets
// are byte offsets into the original source; EndOffset is exclusive.
type Paul2013C3A0Component struct {
	Output      []byte
	StartOffset int
	EndOffset   int
	Stop        Paul2013C3A0ComponentStop
}

// ScanPaul2013C3A0Component ports the bounded ASCII byte scan in
// FUN_1000c860, called by FUN_1000c3a0. It preserves source bytes and the
// native uppercase/lowercase boundary checks, including its special handling
// for apostrophe-s before a class-0x80 byte. sourceLimit is the exclusive
// source boundary passed as the native fourth argument.
func ScanPaul2013C3A0Component(source []byte, startOffset, sourceLimit int) (Paul2013C3A0Component, error) {
	if sourceLimit < 0 || sourceLimit > len(source) {
		return Paul2013C3A0Component{}, fmt.Errorf("component source limit %d is outside 0..%d", sourceLimit, len(source))
	}
	if startOffset < 0 || startOffset > sourceLimit {
		return Paul2013C3A0Component{}, fmt.Errorf("component start offset %d is outside 0..%d", startOffset, sourceLimit)
	}
	for offset, value := range source[:sourceLimit] {
		if value == 0 {
			break
		}
		if value >= 0x80 {
			return Paul2013C3A0Component{}, fmt.Errorf("component source byte 0x%02x at offset %d is outside the supported ASCII input", value, offset)
		}
	}

	result := Paul2013C3A0Component{StartOffset: startOffset, EndOffset: startOffset, Stop: Paul2013C3A0ComponentAtLimit}
	attributes := Paul2013ExceptionCharacterAttributes()
	for offset := startOffset; offset < sourceLimit; {
		current := source[offset]
		if current == 0 {
			result.Stop = Paul2013C3A0ComponentAtNUL
			return result, nil
		}
		if len(result.Output) > 0 {
			previous := result.Output[len(result.Output)-1]
			currentAttributes := attributes[current]
			if offset+1 < sourceLimit && attributes[previous]&0x80 != 0 && currentAttributes&0x80 != 0 &&
				attributes[source[offset+1]]&0x40 != 0 {
				result.Stop = Paul2013C3A0ComponentAtCharacterBoundary
				result.EndOffset = offset
				return result, nil
			}
			if attributes[previous]&0x40 != 0 && currentAttributes&0x80 != 0 {
				result.Stop = Paul2013C3A0ComponentAtCharacterBoundary
				result.EndOffset = offset
				return result, nil
			}
		}

		if current == '\'' {
			if offset+2 < sourceLimit && source[offset+1] == 's' && attributes[source[offset+2]]&0x80 != 0 {
				result.Output = append(result.Output, '\'', 's')
				result.EndOffset = offset + 2
				result.Stop = Paul2013C3A0ComponentAfterApostropheS
				return result, nil
			}
			if len(result.Output) > 0 && result.Output[len(result.Output)-1] == 's' &&
				offset+1 < sourceLimit && attributes[source[offset+1]]&0x80 != 0 {
				result.EndOffset = offset + 1
				result.Stop = Paul2013C3A0ComponentAfterApostrophe
				return result, nil
			}
		}

		result.Output = append(result.Output, current)
		offset++
		result.EndOffset = offset
	}
	return result, nil
}
