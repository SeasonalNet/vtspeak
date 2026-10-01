package text

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	Paul2013TokenResultRowSize      = 0x554
	Paul2013TokenResultRowCount     = 0x00
	Paul2013TokenResultRowIndex     = 0x02
	Paul2013TokenResultRowType      = 0x04
	Paul2013TokenResultRowSurface   = 0x05
	Paul2013TokenResultPathControls = 0x23
	Paul2013TokenResultPhoneStrings = 0x37
	Paul2013TokenResultPhoneStride  = 0x41
	Paul2013TokenResultMetadata     = 0x54c

	paul2013TokenResultSurfaceCapacity     = Paul2013TokenResultPathControls - Paul2013TokenResultRowSurface
	paul2013TokenResultPathCapacity        = Paul2013TokenResultPhoneStrings - Paul2013TokenResultPathControls
	paul2013TokenResultPhoneCapacity       = Paul2013TokenResultPhoneStride - 1
	paul2013TokenResultAlternativeCapacity = 5
)

// WritePaul2013DictionaryPhoneRow applies the fields written by
// FUN_1000d450 to one 0x554-byte token-result row. destination is the row
// pointer passed to that function. Bytes the native function leaves unwritten
// are preserved, so callers can retain upstream or preinitialized state.
func WritePaul2013DictionaryPhoneRow(
	destination []byte,
	tokenIndex uint16,
	surface []byte,
	rows Paul2013DictionaryPhoneRows,
) error {
	if len(destination) < Paul2013TokenResultRowSize {
		return fmt.Errorf("token-result row has %d bytes, need %d", len(destination), Paul2013TokenResultRowSize)
	}
	if len(surface) >= paul2013TokenResultSurfaceCapacity {
		return fmt.Errorf("token surface has %d bytes, maximum is %d", len(surface), paul2013TokenResultSurfaceCapacity-1)
	}
	if containsNUL(surface) {
		return errors.New("token surface contains an embedded NUL")
	}
	if rows.AlternativeCount < 0 || rows.AlternativeCount > paul2013TokenResultAlternativeCapacity {
		return fmt.Errorf("token-result row has %d alternatives, supported maximum is %d", rows.AlternativeCount, paul2013TokenResultAlternativeCapacity)
	}
	if rows.HasContextMarker {
		if rows.AlternativeCount != 0 || len(rows.PhoneStrings) != 0 || len(rows.PathControlBytes) != 0 ||
			!rows.HasMarkerSentinel || rows.MarkerSentinel != 0xff {
			return errors.New("conditional marker row has inconsistent fields")
		}
		if len(rows.SelectedPhone) > paul2013TokenResultPhoneCapacity || containsNUL(rows.SelectedPhone) {
			return fmt.Errorf("selected phone string has invalid length or embedded NUL")
		}
	} else {
		if rows.HasMarkerSentinel || rows.AlternativeCount != len(rows.PhoneStrings) || len(rows.SelectedPhone) != 0 {
			return errors.New("ordinary token-result row has inconsistent alternative fields")
		}
		if rows.AlternativeCount == 0 {
			if len(rows.PhoneStrings) != 0 || len(rows.PathControlBytes) != 1 || rows.PathControlBytes[0] != 0xff {
				return errors.New("empty token-result row must contain only the 0xff path terminator")
			}
		} else if len(rows.PathControlBytes) == 0 || len(rows.PathControlBytes) > paul2013TokenResultPathCapacity ||
			rows.PathControlBytes[len(rows.PathControlBytes)-1] != 0xff {
			return fmt.Errorf("path-control stream must fit %d bytes and end in 0xff", paul2013TokenResultPathCapacity)
		}
		for index, value := range rows.PathControlBytes[:len(rows.PathControlBytes)-1] {
			if value == 0xff {
				return fmt.Errorf("path-control stream has an early terminator at byte %d", index)
			}
		}
		for index, phone := range rows.PhoneStrings {
			if len(phone) > paul2013TokenResultPhoneCapacity || containsNUL(phone) {
				return fmt.Errorf("phone string %d has invalid length or embedded NUL", index)
			}
		}
	}

	// Validate first, then write, so malformed projections leave the row intact.
	row := destination[:Paul2013TokenResultRowSize]
	binary.LittleEndian.PutUint16(row[Paul2013TokenResultRowCount:], uint16(rows.AlternativeCount))
	binary.LittleEndian.PutUint16(row[Paul2013TokenResultRowIndex:], tokenIndex)
	row[Paul2013TokenResultRowType] = rows.ResultType
	copy(row[Paul2013TokenResultRowSurface:], surface)
	row[Paul2013TokenResultRowSurface+len(surface)] = 0
	if rows.HasContextMarker {
		row[0x12] = rows.MarkerSentinel
		row[Paul2013TokenResultPathControls] = rows.ContextMarker
		copyNULTerminated(row[Paul2013TokenResultPhoneStrings:], rows.SelectedPhone)
	} else {
		copy(row[Paul2013TokenResultPathControls:], rows.PathControlBytes)
		for index, phone := range rows.PhoneStrings {
			offset := Paul2013TokenResultPhoneStrings + index*Paul2013TokenResultPhoneStride
			copyNULTerminated(row[offset:], phone)
		}
		terminatorOffset := Paul2013TokenResultPhoneStrings + rows.AlternativeCount*Paul2013TokenResultPhoneStride
		row[terminatorOffset] = 0
	}
	for index, metadata := range rows.Metadata {
		if metadata {
			binary.LittleEndian.PutUint16(row[Paul2013TokenResultMetadata+index*2:], 1)
		} else {
			binary.LittleEndian.PutUint16(row[Paul2013TokenResultMetadata+index*2:], 0)
		}
	}
	return nil
}

// WritePaul2013DictionaryPhoneRow writes this resolved token's native row
// using the caller's zero-based parser token index.
func (token LexicalToken) WritePaul2013DictionaryPhoneRow(destination []byte, tokenIndex uint16) error {
	return WritePaul2013DictionaryPhoneRow(destination, tokenIndex, []byte(token.Surface), token.ModelPhoneRows)
}

// WritePaul2013DictionaryPhoneRow writes this selected token span's native row
// using the caller's zero-based parser token index.
func (token LexicalTokenSpan) WritePaul2013DictionaryPhoneRow(destination []byte, tokenIndex uint16) error {
	return WritePaul2013DictionaryPhoneRow(destination, tokenIndex, []byte(token.Surface), token.ModelPhoneRows)
}

func copyNULTerminated(destination, value []byte) {
	copy(destination, value)
	destination[len(value)] = 0
}

func containsNUL(value []byte) bool {
	for _, character := range value {
		if character == 0 {
			return true
		}
	}
	return false
}
