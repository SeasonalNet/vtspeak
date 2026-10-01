package duration

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"vtspeak/engine/text"
)

const paul2013FUN10008580ModelBase = 0x429a2

type Paul2013FUN10007520UseMarkerResult struct {
	Model   []byte
	Applied bool
}

type Paul2013FUN10007520UpMarkerResult struct {
	Model   []byte
	Applied bool
}

type Paul2013FUN10007520ArticleMarkerResult struct {
	Model   []byte
	Applied bool
	Path    string
}

type Paul2013FUN10007520TheMarkerResult struct {
	Model   []byte
	Applied bool
	Path    string
}

type Paul2013FUN10007520MinuteResult struct {
	Model   []byte
	Applied bool
}

type Paul2013FUN10007520CloseMarkerResult struct {
	Model   []byte
	Applied bool
}

type Paul2013FUN10007520MouthMarkerResult struct {
	Model   []byte
	Applied bool
}

type Paul2013FUN10007520BowResult struct {
	Model      []byte
	Applied    bool
	TableMatch bool
}

type Paul2013FUN10007520LeadResult struct {
	Model   []byte
	Applied bool
}

type Paul2013FUN10007520ReadMarkerResult struct {
	Model   []byte
	Applied bool
}

var paul2013FUN1000CC60BowWords = [][]byte{
	[]byte("arrow"), []byte("compass"), []byte("fiddle"), []byte("tie"), []byte("violin"), []byte("window"),
}

var paul2013FUN1000CC60LeadWords = [][]byte{
	[]byte("atomic"), []byte("dust"), []byte("element"), []byte("hazard"), []byte("metal"),
	[]byte("paint"), []byte("pb"), []byte("pencil"), []byte("pencils"), []byte("toxic"),
}

var paul2013FUN10007520ReadSpecialWords = [][]byte{
	[]byte("'ve"), []byte("sh"), []byte("had not"), []byte("hadn't"), []byte("ch"),
	[]byte("has not"), []byte("hasn't"), []byte("have"), []byte("have not"), []byte("haven't"),
}

var paul2013FUN10007520ReadTemporalPhrases = [][]byte{
	[]byte("ago"), []byte("already"), []byte("before"), []byte("just now"), []byte("last"),
	[]byte("this afternoon"), []byte("this evening"), []byte("this morning"), []byte("yesterday"),
}

// EvaluatePaul2013FUN10008580 ports the row-index predicate used by the
// counted model-context pass. It keeps the native relative offsets and
// mapped-string conditions; the field semantics remain opaque.
func EvaluatePaul2013FUN10008580(model []byte, rowIndex int) (bool, error) {
	if rowIndex < 0 {
		return false, fmt.Errorf("FUN_10008580 row index %d is negative", rowIndex)
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return false, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return false, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex >= rowCount {
		return false, fmt.Errorf("FUN_10008580 row index %d is outside model row count %d", rowIndex, rowCount)
	}
	if rowIndex == 0 {
		return true, nil
	}
	readCString := func(offset int) ([]byte, error) {
		if offset < 0 || offset >= len(model) {
			return nil, fmt.Errorf("FUN_10008580 string starts outside model at %#x", offset)
		}
		area := model[offset:]
		for index, value := range area {
			if value == 0 {
				return area[:index], nil
			}
		}
		return nil, fmt.Errorf("FUN_10008580 string at %#x is not NUL-terminated", offset)
	}
	rowOffset := paul2013FUN10008580ModelBase + rowIndex*paul2013ModelContextRowStride
	if rowIndex > 1 {
		first, err := readCString(rowOffset - 0xd5)
		if err != nil {
			return false, err
		}
		second, err := readCString(rowOffset - 0x65)
		if err != nil {
			return false, err
		}
		if text.Paul2013ContextMappedCStringEqual(first, []byte("let")) ||
			text.Paul2013ContextWordPairSpecial(first, second) {
			return true, nil
		}
		row0, err := readCString(paul2013FUN10008580ModelBase + 0x0b)
		if err != nil {
			return false, err
		}
		row1, err := readCString(paul2013FUN10008580ModelBase + 0x7b)
		if err != nil {
			return false, err
		}
		if text.Paul2013ContextWordPairSpecial(row0, row1) {
			return true, nil
		}
		if text.Paul2013ContextWordPairSpecial(first, nil) &&
			!text.Paul2013ContextMappedCStringEqual(second, []byte("not")) &&
			!text.Paul2013ContextMappedCStringEqual(second, []byte("can't")) {
			return true, nil
		}
	}
	first, err := readCString(rowOffset - 0x65)
	if err != nil {
		return false, err
	}
	if text.Paul2013ContextMappedCStringEqual(first, []byte("do")) ||
		text.Paul2013ContextMappedCStringEqual(first, []byte("to")) ||
		text.Paul2013ContextWordPairSpecial(first, nil) {
		return true, nil
	}
	return false, nil
}

// ApplyPaul2013FUN10007520UseMarker ports the post-handler "use" row rule in
// FUN_10007520. It changes phone-code byte +2 from '7' to 'D' only when
// FUN_10008580 accepts the same counted row, and sets the separate row-state
// bit 3 observed by the caller.
func ApplyPaul2013FUN10007520UseMarker(model []byte, rowIndex int) (Paul2013FUN10007520UseMarkerResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520UseMarkerResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520UseMarkerResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520UseMarkerResult{}, fmt.Errorf("FUN_10007520 use-marker row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520UseMarkerResult{Model: append([]byte(nil), model...)}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520UseMarkerResult{}, fmt.Errorf("read FUN_10007520 use-marker row %d: %w", rowIndex, err)
	}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	codeOffset := rowOffset + 0x2b
	stateOffset := rowOffset + 4
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("use")) || model[codeOffset] != '7' {
		return result, nil
	}
	accepted, err := EvaluatePaul2013FUN10008580(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520UseMarkerResult{}, fmt.Errorf("evaluate FUN_10008580 for use-marker row %d: %w", rowIndex, err)
	}
	if !accepted {
		return result, nil
	}
	result.Model[codeOffset] = 'D'
	result.Model[stateOffset] |= 0x08
	result.Applied = true
	return result, nil
}

// ApplyPaul2013FUN10007520UpMarker ports the adjacent-row "UP" rule from
// FUN_10007520. It checks the next counted surface and writes the three
// observed phone-code bytes plus row-state bit 3 when either condition matches.
func ApplyPaul2013FUN10007520UpMarker(model []byte, rowIndex int) (Paul2013FUN10007520UpMarkerResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520UpMarkerResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520UpMarkerResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520UpMarkerResult{}, fmt.Errorf("FUN_10007520 UP-marker row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520UpMarkerResult{Model: append([]byte(nil), model...)}
	if rowIndex+1 >= rowCount {
		return result, nil
	}
	currentSurface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520UpMarkerResult{}, fmt.Errorf("read FUN_10007520 UP-marker row %d: %w", rowIndex, err)
	}
	if !bytes.Equal(cStringBytes(currentSurface), []byte("UP")) {
		return result, nil
	}
	nextSurface, _, err := text.Paul2013ModelContextSurface(model, rowIndex+1)
	if err != nil {
		return Paul2013FUN10007520UpMarkerResult{}, fmt.Errorf("read FUN_10007520 UP-marker next row %d: %w", rowIndex+1, err)
	}
	if !text.Paul2013ContextMappedCStringEqual(nextSurface, []byte("a")) &&
		!text.Paul2013ContextMappedCStringEqual(nextSurface, []byte("the")) {
		return result, nil
	}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	codeOffset := rowOffset + 0x29
	stateOffset := rowOffset + 4
	if codeOffset+2 >= len(model) || stateOffset >= len(model) {
		return Paul2013FUN10007520UpMarkerResult{}, fmt.Errorf("FUN_10007520 UP-marker output does not fit row %d", rowIndex)
	}
	result.Model[codeOffset] = 0x08
	result.Model[codeOffset+1] = 0x35
	result.Model[codeOffset+2] = 0
	result.Model[stateOffset] |= 0x08
	result.Applied = true
	return result, nil
}

// ApplyPaul2013FUN10007520ArticleMarker ports the `a` post-handler class and
// apostrophe-s cases, including the initial-row gate over recovered class-3
// and class-4 tables.
func ApplyPaul2013FUN10007520ArticleMarker(
	model, parserRows []byte,
	rowIndex int,
) (Paul2013FUN10007520ArticleMarkerResult, error) {
	result := Paul2013FUN10007520ArticleMarkerResult{Model: append([]byte(nil), model...)}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520ArticleMarkerResult{}, err
	}
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("a")) {
		return result, nil
	}
	parserClass, err := paul2013ModelContextParserClass(model, parserRows, rowIndex)
	if err != nil {
		return Paul2013FUN10007520ArticleMarkerResult{}, fmt.Errorf("read article marker parser class at row %d: %w", rowIndex, err)
	}
	if parserClass == 'S' {
		return result, nil
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	rowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	if rowIndex+1 >= rowCount {
		return result, nil
	}
	nextSurface, _, err := text.Paul2013ModelContextSurface(model, rowIndex+1)
	if err != nil {
		return Paul2013FUN10007520ArticleMarkerResult{}, fmt.Errorf("read article marker following surface at row %d: %w", rowIndex, err)
	}
	attributes := text.Paul2013UnsignedCharacterAttributeTable()
	initialContextGate := rowIndex == 0 && len(surface) != 0 &&
		attributes[surface[0]]&0x80 != 0 && model[rowStart+0x23] > 0x1c && model[rowStart+0x23] < 0x20
	if initialContextGate {
		nextClassGate := text.Paul2013TripletTableContains(4, nextSurface) ||
			text.Paul2013TripletTableContains(3, nextSurface) ||
			(len(nextSurface) == 0 || attributes[nextSurface[0]]&0xc0 == 0)
		if !nextClassGate {
			result.Model[rowStart+0x23] = 7
			result.Model[rowStart-2] |= 0x08
			result.Applied = true
			result.Path = "initial-following-class"
			return result, nil
		}
	} else if len(surface) != 0 && len(nextSurface) != 0 &&
		attributes[surface[0]]&0x40 != 0 &&
		model[rowStart+0x23] >= 0x1d && model[rowStart+0x23] <= 0x1f &&
		attributes[nextSurface[0]]&0xc0 != 0 && attributes[nextSurface[0]]&0x80 != 0 {
		result.Model[rowStart+0x23] = 7
		result.Model[rowStart-2] |= 0x08
		result.Applied = true
		result.Path = "following-class"
		return result, nil
	}
	nextModelRowStart := rowStart + paul2013ModelContextRowStride
	if binary.LittleEndian.Uint16(model[nextModelRowStart:]) == binary.LittleEndian.Uint16(model[rowStart:]) &&
		text.Paul2013ContextMappedCStringEqual(nextSurface, []byte("'s")) {
		result.Model[rowStart+0x23] = 0x1e
		result.Model[rowStart-2] |= 0x08
		result.Applied = true
		result.Path = "following-apostrophe-s"
	}
	return result, nil
}

// ApplyPaul2013FUN10007520TheMarker ports the `the` post-handler predicates
// in FUN_10007520. The native branch may write 0x07 or 0x1e to the phone
// marker at row offset +0x23; these values are preserved without assigning
// them phonetic meaning.
func ApplyPaul2013FUN10007520TheMarker(
	model, parserRows []byte,
	rowIndex int,
) (Paul2013FUN10007520TheMarkerResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520TheMarkerResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520TheMarkerResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520TheMarkerResult{}, fmt.Errorf("FUN_10007520 the-marker row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520TheMarkerResult{Model: append([]byte(nil), model...)}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520TheMarkerResult{}, fmt.Errorf("read FUN_10007520 the-marker row %d: %w", rowIndex, err)
	}
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("the")) {
		return result, nil
	}
	parserClass, err := paul2013ModelContextParserClass(model, parserRows, rowIndex)
	if err != nil {
		return Paul2013FUN10007520TheMarkerResult{}, fmt.Errorf("read the-marker parser class at row %d: %w", rowIndex, err)
	}
	if parserClass == 'S' {
		return result, nil
	}
	rowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	phoneMarkerOffset := rowStart + 0x23
	attributes := text.Paul2013UnsignedCharacterAttributeTable()
	apply := func(value byte, path string) {
		if result.Model[phoneMarkerOffset] != value {
			result.Model[phoneMarkerOffset] = value
			result.Applied = true
			result.Path = path
		}
	}
	if rowIndex+1 < rowCount {
		nextSurface, _, err := text.Paul2013ModelContextSurface(model, rowIndex+1)
		if err != nil {
			return Paul2013FUN10007520TheMarkerResult{}, fmt.Errorf("read the-marker following surface at row %d: %w", rowIndex, err)
		}
		marker := model[phoneMarkerOffset]
		if rowIndex == 0 && len(surface) != 0 && attributes[surface[0]]&0x80 != 0 && marker > 0x1c && marker < 0x20 {
			nextClassGate := text.Paul2013TripletTableContains(4, nextSurface) ||
				text.Paul2013TripletTableContains(3, nextSurface) ||
				len(nextSurface) == 0 || attributes[nextSurface[0]]&0xc0 == 0
			if !nextClassGate {
				apply(0x1e, "initial-following-class")
			}
		} else if len(surface) != 0 && attributes[surface[0]]&0x40 != 0 &&
			marker >= 0x1d && marker <= 0x1f &&
			len(nextSurface) != 0 && attributes[nextSurface[0]]&0xc0 != 0 &&
			attributes[nextSurface[0]]&0x80 != 0 {
			apply(0x07, "following-class")
		}
		if binary.LittleEndian.Uint16(model[rowStart:rowStart+2]) ==
			binary.LittleEndian.Uint16(model[rowStart+paul2013ModelContextRowStride:rowStart+paul2013ModelContextRowStride+2]) &&
			!text.Paul2013ContextMappedCStringEqual(nextSurface, []byte("'s")) {
			apply(0x1e, "following-same-parser-row")
		}
	}
	return result, nil
}

// ApplyPaul2013FUN10007520Minute ports the direct "minute" rewrite in
// FUN_10007520. It writes the observed six-byte phone-code region when the
// previous counted surface is accepted by FUN_10008550.
func ApplyPaul2013FUN10007520Minute(model []byte, rowIndex int) (Paul2013FUN10007520MinuteResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520MinuteResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520MinuteResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520MinuteResult{}, fmt.Errorf("FUN_10007520 minute row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520MinuteResult{Model: append([]byte(nil), model...)}
	if rowIndex == 0 {
		return result, nil
	}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520MinuteResult{}, fmt.Errorf("read FUN_10007520 minute row %d: %w", rowIndex, err)
	}
	previousSurface, _, err := text.Paul2013ModelContextSurface(model, rowIndex-1)
	if err != nil {
		return Paul2013FUN10007520MinuteResult{}, fmt.Errorf("read FUN_10007520 minute previous row %d: %w", rowIndex-1, err)
	}
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("minute")) || !text.Paul2013IsNumberWord(previousSurface) {
		return result, nil
	}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	copy(result.Model[rowOffset+0x29:rowOffset+0x2f], []byte{0x2c, 0x24, 0x2d, 0x23, 0x39, 0})
	result.Model[rowOffset+4] |= 0x08
	result.Applied = true
	return result, nil
}

// ApplyPaul2013FUN10007520CloseMarker ports the direct "close" rewrite in
// FUN_10007520, gated by phone-code byte +2 and FUN_10008580.
func ApplyPaul2013FUN10007520CloseMarker(model []byte, rowIndex int) (Paul2013FUN10007520CloseMarkerResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520CloseMarkerResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520CloseMarkerResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520CloseMarkerResult{}, fmt.Errorf("FUN_10007520 close-marker row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520CloseMarkerResult{Model: append([]byte(nil), model...)}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520CloseMarkerResult{}, fmt.Errorf("read FUN_10007520 close-marker row %d: %w", rowIndex, err)
	}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	codeOffset := rowOffset + 0x2c
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("close")) || model[codeOffset] != '7' {
		return result, nil
	}
	accepted, err := EvaluatePaul2013FUN10008580(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520CloseMarkerResult{}, fmt.Errorf("evaluate FUN_10008580 for close-marker row %d: %w", rowIndex, err)
	}
	if !accepted {
		return result, nil
	}
	result.Model[codeOffset] = 'D'
	result.Model[rowOffset+4] |= 0x08
	result.Applied = true
	return result, nil
}

// ApplyPaul2013FUN10007520MouthMarker ports the mode-R, two-row class-12
// phrase lookup and marker rewrite for "mouth" in FUN_10007520.
func ApplyPaul2013FUN10007520MouthMarker(model []byte, rowIndex int) (Paul2013FUN10007520MouthMarkerResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520MouthMarkerResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520MouthMarkerResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520MouthMarkerResult{}, fmt.Errorf("FUN_10007520 mouth-marker row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520MouthMarkerResult{Model: append([]byte(nil), model...)}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520MouthMarkerResult{}, fmt.Errorf("read FUN_10007520 mouth-marker row %d: %w", rowIndex, err)
	}
	if rowIndex <= 1 || !text.Paul2013ContextMappedCStringEqual(surface, []byte("mouth")) || model[rowOffset+0x2b] != 0x16 {
		return result, nil
	}
	for nextRow := rowIndex + 1; nextRow < rowCount && nextRow <= rowIndex+2; nextRow++ {
		nextSurface, _, err := text.Paul2013ModelContextSurface(model, nextRow)
		if err != nil {
			return Paul2013FUN10007520MouthMarkerResult{}, fmt.Errorf("read FUN_10007520 mouth-marker next row %d: %w", nextRow, err)
		}
		if text.Paul2013TripletClass12Contains(nextSurface) {
			result.Model[rowOffset+0x2b] = 0x3a
			result.Model[rowOffset+4] |= 0x08
			result.Applied = true
			return result, nil
		}
	}
	return result, nil
}

// ApplyPaul2013FUN10007520Bow ports the bounded mode-B table lookup for
// "bow" in FUN_10007520. The native branch writes one of two observed phone
// codes and sets row-state bit 3 regardless of lookup outcome.
func ApplyPaul2013FUN10007520Bow(model []byte, rowIndex int) (Paul2013FUN10007520BowResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520BowResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520BowResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520BowResult{}, fmt.Errorf("FUN_10007520 bow row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520BowResult{Model: append([]byte(nil), model...)}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520BowResult{}, fmt.Errorf("read FUN_10007520 bow row %d: %w", rowIndex, err)
	}
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("bow")) {
		return result, nil
	}
	_, matched, err := FindPaul2013ContextPhraseWindow(model, rowIndex, 'B', 5, paul2013FUN1000CC60BowWords, bytes.Compare)
	if err != nil {
		return Paul2013FUN10007520BowResult{}, fmt.Errorf("search FUN_10007520 bow context row %d: %w", rowIndex, err)
	}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	result.Model[rowOffset+0x2a] = 0x0e
	if matched {
		result.Model[rowOffset+0x2a] = 0x30
	}
	result.Model[rowOffset+4] |= 0x08
	result.Applied = true
	result.TableMatch = matched
	return result, nil
}

// ApplyPaul2013FUN10007520Lead ports the mode-B, ten-row context-table
// rewrite for "lead" when its first phone-code byte is an apostrophe.
func ApplyPaul2013FUN10007520Lead(model []byte, rowIndex int) (Paul2013FUN10007520LeadResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520LeadResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520LeadResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520LeadResult{}, fmt.Errorf("FUN_10007520 lead row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520LeadResult{Model: append([]byte(nil), model...)}
	surface, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520LeadResult{}, fmt.Errorf("read FUN_10007520 lead row %d: %w", rowIndex, err)
	}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	if !text.Paul2013ContextMappedCStringEqual(surface, []byte("lead")) || model[rowOffset+0x2a] != '\'' {
		return result, nil
	}
	_, matched, err := FindPaul2013ContextPhraseWindow(model, rowIndex, 'B', 10, paul2013FUN1000CC60LeadWords, bytes.Compare)
	if err != nil {
		return Paul2013FUN10007520LeadResult{}, fmt.Errorf("search FUN_10007520 lead context row %d: %w", rowIndex, err)
	}
	if !matched {
		return result, nil
	}
	result.Model[rowOffset+0x2a] = 0x18
	result.Model[rowOffset+4] |= 0x08
	result.Applied = true
	return result, nil
}

// ApplyPaul2013FUN10007520ReadMarker ports the direct row predicates and
// bounded temporal phrase lookup used by the ambiguous "read" branch.
func ApplyPaul2013FUN10007520ReadMarker(model []byte, rowIndex int) (Paul2013FUN10007520ReadMarkerResult, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013FUN10007520ReadMarkerResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013FUN10007520ReadMarkerResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return Paul2013FUN10007520ReadMarkerResult{}, fmt.Errorf("FUN_10007520 read row %d is outside model row count %d", rowIndex, rowCount)
	}
	result := Paul2013FUN10007520ReadMarkerResult{Model: append([]byte(nil), model...)}
	rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
	current, _, err := text.Paul2013ModelContextSurface(model, rowIndex)
	if err != nil {
		return Paul2013FUN10007520ReadMarkerResult{}, fmt.Errorf("read FUN_10007520 current row %d: %w", rowIndex, err)
	}
	if !text.Paul2013ContextMappedCStringEqual(current, []byte("read")) || model[rowOffset+0x2a] != '\'' || rowIndex == 0 {
		return result, nil
	}
	previous, err := paul2013ContextSurfaceBytes(model, rowIndex-1)
	if err != nil {
		return Paul2013FUN10007520ReadMarkerResult{}, err
	}
	previousPrevious := []byte(nil)
	if rowIndex > 1 {
		previousPrevious, err = paul2013ContextSurfaceBytes(model, rowIndex-2)
		if err != nil {
			return Paul2013FUN10007520ReadMarkerResult{}, err
		}
	}
	previousThird := []byte(nil)
	if rowIndex > 2 {
		previousThird, err = paul2013ContextSurfaceBytes(model, rowIndex-3)
		if err != nil {
			return Paul2013FUN10007520ReadMarkerResult{}, err
		}
	}
	readSpecialTableContains := func(surface []byte) bool {
		for _, entry := range paul2013FUN10007520ReadSpecialWords {
			if text.Paul2013ContextMappedCStringEqual(surface, entry) {
				return true
			}
		}
		return false
	}
	shouldMark := text.Paul2013TripletTableContains(3, previous) ||
		(rowIndex > 1 && text.Paul2013TripletTableContains(3, previousPrevious) &&
			text.Paul2013ContextMappedCStringEqual(previous, []byte("not"))) ||
		readSpecialTableContains(previous) ||
		(rowIndex > 1 && readSpecialTableContains(previousPrevious) &&
			text.Paul2013ContextMappedCStringEqual(previous, []byte("not"))) ||
		text.Paul2013ContextMappedCStringEqual(previous, []byte("he")) ||
		text.Paul2013ContextMappedCStringEqual(previous, []byte("she"))
	if !shouldMark {
		_, matched, err := FindPaul2013ContextPhraseWindow(
			model, rowIndex, 'B', 20, paul2013FUN10007520ReadTemporalPhrases, bytes.Compare,
		)
		if err != nil {
			return Paul2013FUN10007520ReadMarkerResult{}, fmt.Errorf("search FUN_10007520 read temporal context at row %d: %w", rowIndex, err)
		}
		previousIsI := text.Paul2013ContextMappedCStringEqual(previous, []byte("i"))
		previousIsClass8 := text.Paul2013TripletTableContains(8, previous)
		if (previousIsI || previousIsClass8) && matched {
			shouldMark = true
		} else if text.Paul2013TripletTableContains(11, previous) {
			shouldMark = true
		} else if rowIndex > 2 &&
			paul2013OneOfMapped(previousThird, []byte("get"), []byte("got"), []byte("gets")) &&
			paul2013OneOfMapped(previousPrevious, []byte("get"), []byte("got"), []byte("gets")) &&
			!text.Paul2013ContextMappedCStringEqual(previous, []byte("to")) {
			shouldMark = true
		}
	}
	if !shouldMark {
		return result, nil
	}
	result.Model[rowOffset+0x2a] = 0x18
	result.Model[rowOffset+4] |= 0x08
	result.Applied = true
	return result, nil
}

func paul2013OneOfMapped(surface []byte, choices ...[]byte) bool {
	for _, choice := range choices {
		if text.Paul2013ContextMappedCStringEqual(surface, choice) {
			return true
		}
	}
	return false
}
