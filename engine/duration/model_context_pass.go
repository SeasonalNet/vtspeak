package duration

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"vtspeak/engine/text"
)

const (
	paul2013ModelContextCountOffset = 0x429a2
	paul2013ModelContextRowBase     = 0x429a8
	paul2013ModelContextRowStride   = 0x70
)

func paul2013FUN100086C0ModelSource(model []byte, rowIndex int) ([]byte, error) {
	if len(model) < paul2013ModelContextCountOffset+2 {
		return nil, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return nil, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if rowIndex < 0 || rowIndex >= rowCount {
		return nil, fmt.Errorf("FUN_100086c0 model-source row %d is outside context count %d", rowIndex, rowCount)
	}
	rowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	// FUN_10007520 passes model + 0x429a2 + 0x0b + row*0x70; counted
	// context rows begin six bytes after that base, so this is row offset +5.
	const (
		sourceOffset = 0x05
		sourceEnd    = 0x23
	)
	sourceArea := model[rowStart+sourceOffset : rowStart+sourceEnd]
	nul := bytes.IndexByte(sourceArea, 0)
	if nul < 0 {
		return nil, fmt.Errorf("FUN_100086c0 source at row offset +%#x is not NUL-terminated", sourceOffset)
	}
	return append([]byte(nil), sourceArea[:nul]...), nil
}

func paul2013ModelContextParserClass(model, parserRows []byte, rowIndex int) (byte, error) {
	if rowIndex < 0 || len(model) < paul2013ModelContextRowBase+(rowIndex+1)*paul2013ModelContextRowStride {
		return 0, fmt.Errorf("model context row %d exceeds model size %d", rowIndex, len(model))
	}
	if len(parserRows)%0x94 != 0 {
		return 0, fmt.Errorf("parser rows have %d bytes, not a multiple of 0x94", len(parserRows))
	}
	rowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
	parserIndex := int(int16(binary.LittleEndian.Uint16(model[rowStart:])))
	parserCount := len(parserRows) / 0x94
	if parserIndex < 0 || parserIndex >= parserCount {
		return 0, fmt.Errorf("model context row %d references parser row %d outside %d rows", rowIndex, parserIndex, parserCount)
	}
	parserStart := parserIndex * 0x94
	return parserRows[parserStart+0x24], nil
}

// Paul2013ModelContextPassResult reports the known first-pass actions and
// final counted-row mutations. UnresolvedRows still require native handlers
// that are not represented by this pass.
type Paul2013ModelContextPassResult struct {
	Model                     []byte
	FormYRows                 []int
	FormSRows                 []int
	TailRows                  []int
	DuplicateRowsSkipped      []int
	ExceptionRows             []int
	ExceptionRowsSkipped      []int
	CompoundAttemptedRows     []int
	CompoundRows              []int
	CompoundUnmatchedRows     []int
	SuffixAttemptedRows       []int
	SuffixRows                []int
	ContractionRows           []int
	ContractionUnmatchedRows  []int
	B800Rows                  []int
	B800UnmatchedRows         []int
	C3A0AttemptedRows         []int
	UseMarkerRows             []int
	UpMarkerRows              []int
	MinuteRows                []int
	CloseMarkerRows           []int
	MouthMarkerRows           []int
	ArticleMarkerRows         []int
	TheMarkerRows             []int
	BowRows                   []int
	BowTableMatches           []int
	LeadRows                  []int
	ReadMarkerRows            []int
	C3A0Rows                  []int
	C3A0UnmatchedRows         []int
	TrailingApostropheRows    []int
	GenericFallbackRows       []int
	TerminalContextRows       []int
	PrecheckContextRows       []int
	GenericPrecheckRows       []int
	FallbackMarkerSkippedRows []int
	OuterGateSkippedRows      []int
	SourceClassRows           []int
	UnresolvedRows            []int
	SecondPassRows            []int
}

// Paul2013CompoundOuterMode derives FUN_10007520's short passed to
// FUN_10009dc0 for one counted context row. The native mode is one only when
// the preceding context row maps to a parser row whose dword at +0x2c is 0x0b.
func Paul2013CompoundOuterMode(model, parserRows []byte, contextRowIndex int) (bool, error) {
	const (
		modelContextCountOffset = 0x429a2
		modelContextRowBase     = 0x429a8
		modelContextRowStride   = 0x70
		modelContextParserIndex = 0x00
		parserRowStride         = 0x94
		parserRowTypeOffset     = 0x2c
		compoundParserRowType   = 0x0b
	)
	if len(model) < modelContextCountOffset+2 {
		return false, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), modelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[modelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-modelContextRowBase)/modelContextRowStride {
		return false, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	if contextRowIndex < 0 || contextRowIndex >= rowCount {
		return false, fmt.Errorf("compound outer-mode row %d is outside context count %d", contextRowIndex, rowCount)
	}
	if contextRowIndex == 0 {
		return false, nil
	}
	previousRow := modelContextRowBase + (contextRowIndex-1)*modelContextRowStride
	parserIndex := int(int16(binary.LittleEndian.Uint16(model[previousRow+modelContextParserIndex:])))
	if parserIndex < 0 {
		return false, fmt.Errorf("preceding context row %d has negative parser-row index %d", contextRowIndex-1, parserIndex)
	}
	parserRow := parserIndex * parserRowStride
	if parserRow > len(parserRows) || len(parserRows)-parserRow < parserRowStride {
		return false, fmt.Errorf("preceding context row %d references parser row %d outside %d bytes", contextRowIndex-1, parserIndex, len(parserRows))
	}
	return binary.LittleEndian.Uint32(parserRows[parserRow+parserRowTypeOffset:]) == compoundParserRowType, nil
}

// RunPaul2013KnownModelContextPassesWithSelectedCompoundRows derives explicit
// outer-mode overrides for the listed rows, then runs the recovered ordered
// passes. Compound and later supported handler eligibility is evaluated
// automatically; the historical row lists remain comparison controls.
func (engine *Engine) RunPaul2013KnownModelContextPassesWithSelectedCompoundRows(
	ctx context.Context,
	model, parserRows []byte,
	tailSurfaces map[int][]byte,
	compoundContextRowIndexes, c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013ModelContextPassResult, error) {
	compoundModes := make(map[int]bool, len(compoundContextRowIndexes))
	for _, rowIndex := range compoundContextRowIndexes {
		if _, exists := compoundModes[rowIndex]; exists {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("compound row selection repeats context row %d", rowIndex)
		}
		mode, err := Paul2013CompoundOuterMode(model, parserRows, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("derive compound outer mode for row %d: %w", rowIndex, err)
		}
		compoundModes[rowIndex] = mode
	}
	return engine.RunPaul2013KnownModelContextPassesWithCompoundAndC3A0Rows(
		ctx, model, parserRows, tailSurfaces, compoundModes,
		c3a0ContextRowIndexes, contractionPrefix,
	)
}

// RunPaul2013KnownModelContextPasses composes recovered FUN_10007520 stages
// in model-row order. It initializes per-row state, applies the directly
// ported Y/S context mutations, performs exception lookup and native
// matched-row skipping, runs the supported fallback chain, then applies
// FUN_100091b0's final counted-row pass.
//
// tailSurfaces carries optional explicit overrides for the literal suffix
// cases; when an override is absent, the form-Y/S helpers derive the suffix
// input from the parser-row arena. Unported suffix and TPP branches remain
// unresolved. The WithC3A0Rows wrapper retains its historical name and takes
// row indexes for conflict checks; C3A0's gate runs automatically.
func (engine *Engine) RunPaul2013KnownModelContextPasses(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	tailSurfaces map[int][]byte,
) (Paul2013ModelContextPassResult, error) {
	return engine.runPaul2013KnownModelContextPasses(ctx, model, parserRows, tailSurfaces, nil, nil, nil, nil, nil)
}

// RunPaul2013KnownModelContextPassesWithC3A0Rows is kept for API compatibility.
// C3A0 is attempted automatically after preceding native handlers; supplied
// row indexes are retained to detect rows consumed by higher-priority stages.
func (engine *Engine) RunPaul2013KnownModelContextPassesWithC3A0Rows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	tailSurfaces map[int][]byte,
	c3a0ContextRowIndexes []int,
) (Paul2013ModelContextPassResult, error) {
	selected := make(map[int]struct{}, len(c3a0ContextRowIndexes))
	for _, rowIndex := range c3a0ContextRowIndexes {
		if _, exists := selected[rowIndex]; exists {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("C3A0 row selection repeats context row %d", rowIndex)
		}
		selected[rowIndex] = struct{}{}
	}
	return engine.runPaul2013KnownModelContextPasses(ctx, model, parserRows, tailSurfaces, nil, nil, nil, selected, nil)
}

// RunPaul2013KnownModelContextPassesWithCompoundAndC3A0Rows preserves explicit
// FUN_10009dc0 mode overrides and C3A0 conflict checks. The handlers
// otherwise use their recovered native gates automatically.
func (engine *Engine) RunPaul2013KnownModelContextPassesWithCompoundAndC3A0Rows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	tailSurfaces map[int][]byte,
	compoundOuterModes map[int]bool,
	c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013ModelContextPassResult, error) {
	return engine.RunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0Rows(
		ctx, model, parserRows, tailSurfaces, compoundOuterModes, nil,
		c3a0ContextRowIndexes, contractionPrefix,
	)
}

// RunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0Rows keeps
// historical selections as mode overrides and conflict checks.
// Supported suffix, contraction, B800, and C3A0 handlers run automatically in
// native order after their predecessors miss.
func (engine *Engine) RunPaul2013KnownModelContextPassesWithCompoundContractionAndC3A0Rows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	tailSurfaces map[int][]byte,
	compoundOuterModes map[int]bool,
	contractionOuterModes map[int]bool,
	c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013ModelContextPassResult, error) {
	return engine.RunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0Rows(
		ctx, model, parserRows, tailSurfaces, compoundOuterModes, contractionOuterModes,
		nil, c3a0ContextRowIndexes, contractionPrefix,
	)
}

// RunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0Rows
// retains historical row lists for compatibility. B800 and C3A0 now run
// automatically after the preceding supported handlers miss.
func (engine *Engine) RunPaul2013KnownModelContextPassesWithCompoundContractionB800AndC3A0Rows(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	tailSurfaces map[int][]byte,
	compoundOuterModes map[int]bool,
	contractionOuterModes map[int]bool,
	b800ContextRowIndexes []int,
	c3a0ContextRowIndexes []int,
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013ModelContextPassResult, error) {
	selectedC3A0 := make(map[int]struct{}, len(c3a0ContextRowIndexes))
	for _, rowIndex := range c3a0ContextRowIndexes {
		if _, exists := selectedC3A0[rowIndex]; exists {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("C3A0 row selection repeats context row %d", rowIndex)
		}
		selectedC3A0[rowIndex] = struct{}{}
	}
	selectedCompound := make(map[int]bool, len(compoundOuterModes))
	for rowIndex, outerMode := range compoundOuterModes {
		selectedCompound[rowIndex] = outerMode
	}
	selectedContraction := make(map[int]bool, len(contractionOuterModes))
	for rowIndex, outerMode := range contractionOuterModes {
		selectedContraction[rowIndex] = outerMode
	}
	selectedB800 := make(map[int]struct{}, len(b800ContextRowIndexes))
	for _, rowIndex := range b800ContextRowIndexes {
		if _, exists := selectedB800[rowIndex]; exists {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("B800 row selection repeats context row %d", rowIndex)
		}
		selectedB800[rowIndex] = struct{}{}
	}
	return engine.runPaul2013KnownModelContextPasses(
		ctx, model, parserRows, tailSurfaces, selectedCompound, selectedContraction,
		selectedB800, selectedC3A0, contractionPrefix,
	)
}

func (engine *Engine) runPaul2013KnownModelContextPasses(
	ctx context.Context,
	model []byte,
	parserRows []byte,
	tailSurfaces map[int][]byte,
	selectedCompoundRows map[int]bool,
	selectedContractionRows map[int]bool,
	selectedB800Rows map[int]struct{},
	selectedC3A0Rows map[int]struct{},
	contractionPrefix Paul2013CompoundContractionPrefixNormalizer,
) (Paul2013ModelContextPassResult, error) {
	if engine == nil {
		return Paul2013ModelContextPassResult{}, errors.New("Paul 2013 duration engine is nil")
	}
	if ctx == nil {
		return Paul2013ModelContextPassResult{}, errors.New("model-context pass has no context")
	}
	if err := ctx.Err(); err != nil {
		return Paul2013ModelContextPassResult{}, err
	}
	if len(model) < paul2013ModelContextCountOffset+2 {
		return Paul2013ModelContextPassResult{}, fmt.Errorf("model has %d bytes, need context count at %#x", len(model), paul2013ModelContextCountOffset)
	}
	rowCount := int(int16(binary.LittleEndian.Uint16(model[paul2013ModelContextCountOffset:])))
	if rowCount < 0 || rowCount > (len(model)-paul2013ModelContextRowBase)/paul2013ModelContextRowStride {
		return Paul2013ModelContextPassResult{}, fmt.Errorf("model context count %d does not fit model size %d", rowCount, len(model))
	}
	for rowIndex := range selectedC3A0Rows {
		if rowIndex < 0 || rowIndex >= rowCount {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("selected C3A0 context row %d is outside model row count %d", rowIndex, rowCount)
		}
	}
	for rowIndex := range selectedB800Rows {
		if rowIndex < 0 || rowIndex >= rowCount {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("selected B800 context row %d is outside model row count %d", rowIndex, rowCount)
		}
	}
	for rowIndex := range selectedCompoundRows {
		if rowIndex < 0 || rowIndex >= rowCount {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("selected compound context row %d is outside model row count %d", rowIndex, rowCount)
		}
	}
	for rowIndex := range selectedContractionRows {
		if rowIndex < 0 || rowIndex >= rowCount {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("selected contraction context row %d is outside model row count %d", rowIndex, rowCount)
		}
	}
	initialized, err := text.InitializePaul2013ModelContextProcessingFlags(model)
	if err != nil {
		return Paul2013ModelContextPassResult{}, fmt.Errorf("initialize model-context processing flags: %w", err)
	}
	result := Paul2013ModelContextPassResult{Model: initialized.Model}
	postHandlerGate := func(rowIndex int) (bool, error) {
		gate, err := text.BuildPaul2013ModelContextExceptionGate(
			result.Model, parserRows, rowIndex, 0,
		)
		if err != nil {
			return false, err
		}
		return text.ShouldRunPaul2013ModelContextPostHandlers(gate), nil
	}
	applyUseMarker := func(rowIndex int) error {
		article, err := ApplyPaul2013FUN10007520ArticleMarker(result.Model, parserRows, rowIndex)
		if err != nil {
			return fmt.Errorf("apply article marker at row %d: %w", rowIndex, err)
		}
		result.Model = article.Model
		if article.Applied {
			result.ArticleMarkerRows = append(result.ArticleMarkerRows, rowIndex)
		}
		theMarker, err := ApplyPaul2013FUN10007520TheMarker(result.Model, parserRows, rowIndex)
		if err != nil {
			return fmt.Errorf("apply the marker at row %d: %w", rowIndex, err)
		}
		if theMarker.Applied {
			result.Model = theMarker.Model
			result.TheMarkerRows = append(result.TheMarkerRows, rowIndex)
		}
		surface, _, err := text.Paul2013ModelContextSurface(result.Model, rowIndex)
		if err != nil {
			return fmt.Errorf("read post-handler surface: %w", err)
		}
		rowOffset := paul2013ModelContextCountOffset + rowIndex*paul2013ModelContextRowStride
		if text.Paul2013ContextMappedCStringEqual(surface, []byte("use")) &&
			result.Model[rowOffset+0x2b] == '7' {
			marker, err := ApplyPaul2013FUN10007520UseMarker(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if marker.Applied {
				result.Model = marker.Model
				result.UseMarkerRows = append(result.UseMarkerRows, rowIndex)
			}
		}
		if rowIndex+1 < rowCount && bytes.Equal(cStringBytes(surface), []byte("UP")) {
			upMarker, err := ApplyPaul2013FUN10007520UpMarker(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if upMarker.Applied {
				result.Model = upMarker.Model
				result.UpMarkerRows = append(result.UpMarkerRows, rowIndex)
			}
		}
		if text.Paul2013ContextMappedCStringEqual(surface, []byte("read")) && result.Model[rowOffset+0x2a] == '\'' {
			readMarker, err := ApplyPaul2013FUN10007520ReadMarker(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if readMarker.Applied {
				result.Model = readMarker.Model
				result.ReadMarkerRows = append(result.ReadMarkerRows, rowIndex)
			}
		}
		if rowIndex > 0 && text.Paul2013ContextMappedCStringEqual(surface, []byte("minute")) {
			minute, err := ApplyPaul2013FUN10007520Minute(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if minute.Applied {
				result.Model = minute.Model
				result.MinuteRows = append(result.MinuteRows, rowIndex)
			}
		}
		if text.Paul2013ContextMappedCStringEqual(surface, []byte("close")) &&
			result.Model[rowOffset+0x2c] == '7' {
			closeMarker, err := ApplyPaul2013FUN10007520CloseMarker(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if closeMarker.Applied {
				result.Model = closeMarker.Model
				result.CloseMarkerRows = append(result.CloseMarkerRows, rowIndex)
			}
		}
		if text.Paul2013ContextMappedCStringEqual(surface, []byte("bow")) {
			bow, err := ApplyPaul2013FUN10007520Bow(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if bow.Applied {
				result.Model = bow.Model
				result.BowRows = append(result.BowRows, rowIndex)
				if bow.TableMatch {
					result.BowTableMatches = append(result.BowTableMatches, rowIndex)
				}
			}
		}
		if text.Paul2013ContextMappedCStringEqual(surface, []byte("lead")) && result.Model[rowOffset+0x2a] == '\'' {
			lead, err := ApplyPaul2013FUN10007520Lead(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if lead.Applied {
				result.Model = lead.Model
				result.LeadRows = append(result.LeadRows, rowIndex)
			}
		}
		if rowIndex > 1 && text.Paul2013ContextMappedCStringEqual(surface, []byte("mouth")) &&
			result.Model[rowOffset+0x2b] == 0x16 {
			mouthMarker, err := ApplyPaul2013FUN10007520MouthMarker(result.Model, rowIndex)
			if err != nil {
				return err
			}
			if mouthMarker.Applied {
				result.Model = mouthMarker.Model
				result.MouthMarkerRows = append(result.MouthMarkerRows, rowIndex)
			}
		}
		return nil
	}
	selectedByEarlierBranch := func(rowIndex int) bool {
		if _, selected := selectedCompoundRows[rowIndex]; selected {
			return true
		}
		if _, selected := selectedContractionRows[rowIndex]; selected {
			return true
		}
		if _, selected := selectedB800Rows[rowIndex]; selected {
			return true
		}
		_, selected := selectedC3A0Rows[rowIndex]
		return selected
	}

	for rowIndex := 0; rowIndex < rowCount; {
		if err := ctx.Err(); err != nil {
			return Paul2013ModelContextPassResult{}, err
		}
		tail := tailSurfaces[rowIndex]
		formY, err := text.ApplyPaul2013ModelContextFormYWithTail(
			result.Model, parserRows, rowIndex, tail,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d form-Y mutation: %w", rowIndex, err)
		}
		if formY.Applied {
			if selectedByEarlierBranch(rowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was handled by the higher-priority form-Y branch", rowIndex)
			}
			result.Model = formY.Model
			if !formY.TailHandled {
				tailResult, err := engine.ApplyPaul2013ModelContextTail(
					ctx, result.Model, formY.LastRowIndex, formY.TailSurface,
				)
				if err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_100091b0 tail cascade: %w", formY.LastRowIndex, err)
				}
				if tailResult.Applied {
					result.Model = tailResult.Model
					result.TailRows = append(result.TailRows, formY.LastRowIndex)
				}
			}
			postHandlers, err := postHandlerGate(formY.LastRowIndex)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("derive row %d FUN_10007520 post-handler gate after form-Y: %w", formY.LastRowIndex, err)
			}
			if postHandlers {
				if err := applyUseMarker(formY.LastRowIndex); err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", formY.LastRowIndex, err)
				}
			} else {
				result.OuterGateSkippedRows = append(result.OuterGateSkippedRows, formY.LastRowIndex)
			}
			result.FormYRows = append(result.FormYRows, rowIndex)
			for skipped := rowIndex + 1; skipped < formY.LastRowIndex; skipped++ {
				if selectedByEarlierBranch(skipped) {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was skipped by form-Y row %d", skipped, rowIndex)
				}
				result.DuplicateRowsSkipped = append(result.DuplicateRowsSkipped, skipped)
			}
			if selectedByEarlierBranch(formY.LastRowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was consumed by form-Y row %d", formY.LastRowIndex, rowIndex)
			}
			result.UnresolvedRows = append(result.UnresolvedRows, formY.LastRowIndex)
			rowIndex = formY.LastRowIndex + 1
			continue
		}
		formS, err := text.ApplyPaul2013ModelContextFormSWithTail(
			result.Model, parserRows, rowIndex, tail,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d form-S mutation: %w", rowIndex, err)
		}
		if formS.Applied {
			if selectedByEarlierBranch(rowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was handled by the higher-priority form-S branch", rowIndex)
			}
			result.Model = formS.Model
			if !formS.TailHandled {
				tailResult, err := engine.ApplyPaul2013ModelContextTail(
					ctx, result.Model, rowIndex, formS.TailSurface,
				)
				if err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_100091b0 tail cascade: %w", rowIndex, err)
				}
				if tailResult.Applied {
					result.Model = tailResult.Model
					result.TailRows = append(result.TailRows, rowIndex)
				}
			}
			postHandlers, err := postHandlerGate(rowIndex)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("derive row %d FUN_10007520 post-handler gate after form-S: %w", rowIndex, err)
			}
			if postHandlers {
				if err := applyUseMarker(rowIndex); err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
				}
			} else {
				result.OuterGateSkippedRows = append(result.OuterGateSkippedRows, rowIndex)
			}
			result.FormSRows = append(result.FormSRows, rowIndex)
			result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
			rowIndex++
			continue
		}

		exceptionGate, updated, match, _, found, err := engine.ResolvePaul2013PronunciationExceptionFromModelContextRows(
			parserRows, result.Model, rowIndex,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("resolve row %d pronunciation exception: %w", rowIndex, err)
		}
		if !text.ShouldRunPaul2013ModelContextPostHandlers(exceptionGate.Input) {
			if selectedByEarlierBranch(rowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d is blocked by the FUN_10007520 outer gate", rowIndex)
			}
			result.OuterGateSkippedRows = append(result.OuterGateSkippedRows, rowIndex)
			rowIndex++
			continue
		}
		if found {
			if match.TokenCount < 1 || rowIndex+match.TokenCount > rowCount {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("exception at row %d consumes invalid token count %d for %d rows", rowIndex, match.TokenCount, rowCount)
			}
			for skipped := rowIndex; skipped < rowIndex+match.TokenCount; skipped++ {
				if selectedByEarlierBranch(skipped) {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was consumed by pronunciation exception at row %d", skipped, rowIndex)
				}
			}
			result.Model = updated
			result.ExceptionRows = append(result.ExceptionRows, rowIndex)
			for skipped := rowIndex + 1; skipped < rowIndex+match.TokenCount; skipped++ {
				result.ExceptionRowsSkipped = append(result.ExceptionRowsSkipped, skipped)
			}
			rowIndex += match.TokenCount
			continue
		}
		if !exceptionGate.Dispatch {
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 post-handler rules after closed dispatch gate: %w", rowIndex, err)
			}
			result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
			rowIndex++
			continue
		}
		normalized, err := engine.NormalizePaul2013ModelContextRow(
			ctx, result.Model, parserRows, rowIndex,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d source class before fallback chain: %w", rowIndex, err)
		}
		if normalized.Handled {
			if selectedByEarlierBranch(rowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was consumed by higher-priority FUN_10009030", rowIndex)
			}
			result.Model = normalized.Model
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 post-handler rules: %w", rowIndex, err)
			}
			result.SourceClassRows = append(result.SourceClassRows, rowIndex)
			result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
			rowIndex++
			continue
		}
		precheckSource, err := paul2013FUN100086C0ModelSource(result.Model, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("read row %d FUN_100086c0 source: %w", rowIndex, err)
		}
		parserClass, err := paul2013ModelContextParserClass(result.Model, parserRows, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("read row %d parser class before fallback chain: %w", rowIndex, err)
		}
		phoneRowStart := paul2013ModelContextRowBase + rowIndex*paul2013ModelContextRowStride
		if parserClass == 'A' && binary.LittleEndian.Uint16(result.Model[phoneRowStart+0x6c:phoneRowStart+0x6e]) != 0 {
			if bytes.IndexByte(precheckSource, '-') >= 0 {
				currentFlags := result.Model[phoneRowStart-2]
				hyphenated, err := engine.NormalizePaul2013HyphenatedGenericToken(
					ctx, precheckSource, currentFlags,
				)
				if err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d FUN_10003110 precheck: %w", rowIndex, err)
				}
				result.Model, err = text.ApplyPaul2013ModelContextCodes(
					result.Model, rowIndex, hyphenated.Codes, hyphenated.Flags,
				)
				if err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("write row %d FUN_10003110 precheck: %w", rowIndex, err)
				}
				if hyphenated.NativeShort != 0 {
					if selectedByEarlierBranch(rowIndex) {
						return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was consumed by higher-priority FUN_10003110", rowIndex)
					}
					result.Model[phoneRowStart-2] |= 0x40
					if err := applyUseMarker(rowIndex); err != nil {
						return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 post-handler rules after hyphen precheck: %w", rowIndex, err)
					}
					result.GenericPrecheckRows = append(result.GenericPrecheckRows, rowIndex)
					result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
					rowIndex++
					continue
				}
			} else {
				generic, err := engine.NormalizePaul2013GenericToken(ctx, precheckSource)
				if err != nil {
					return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d FUN_10002f10 precheck: %w", rowIndex, err)
				}
				if generic.Eligible {
					if selectedByEarlierBranch(rowIndex) {
						return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was consumed by higher-priority FUN_10002f10", rowIndex)
					}
					stateOffset := phoneRowStart - 2
					phoneStart := phoneRowStart + 0x23
					phoneEnd := phoneRowStart + 0x64
					writes, err := ApplyPaul2013GenericNormalizerWrites(
						result.Model[phoneStart:phoneEnd], result.Model[stateOffset], generic,
					)
					if err != nil {
						return Paul2013ModelContextPassResult{}, fmt.Errorf("write row %d FUN_10002f10 precheck: %w", rowIndex, err)
					}
					copy(result.Model[phoneStart:phoneEnd], writes.Output)
					result.Model[stateOffset] = writes.Flags | 0x40
					if err := applyUseMarker(rowIndex); err != nil {
						return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 post-handler rules after generic precheck: %w", rowIndex, err)
					}
					result.GenericPrecheckRows = append(result.GenericPrecheckRows, rowIndex)
					result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
					rowIndex++
					continue
				}
			}
		}
		precheck, _, err := engine.EvaluatePaul2013FUN100086C0SupportedPath(
			ctx, precheckSource, result.Model, rowIndex, nil, false,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("evaluate row %d FUN_100086c0 gate: %w", rowIndex, err)
		}
		if precheck.Disposition == Paul2013FUN100086C0NeedsNeighborScan {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("row %d FUN_100086c0 neighbor decision remained unresolved", rowIndex)
		}
		if precheck.NativeShort != 0 {
			if selectedByEarlierBranch(rowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d was consumed by the higher-priority FUN_100086c0 branch", rowIndex)
			}
			codes, err := text.EncodePaul2013ContextString(precheckSource)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("encode row %d FUN_100086c0 context: %w", rowIndex, err)
			}
			result.Model, err = text.ApplyPaul2013ModelContextCodes(result.Model, rowIndex, codes, 0x08)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_100086c0 context: %w", rowIndex, err)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 post-handler rules: %w", rowIndex, err)
			}
			result.PrecheckContextRows = append(result.PrecheckContextRows, rowIndex)
			result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
			rowIndex++
			continue
		}
		if result.Model[phoneRowStart+0x29] != 0 {
			if selectedByEarlierBranch(rowIndex) {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected fallback row %d is blocked by the native nonzero phone-marker gate", rowIndex)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 post-handler rules after phone-marker gate: %w", rowIndex, err)
			}
			result.FallbackMarkerSkippedRows = append(result.FallbackMarkerSkippedRows, rowIndex)
			result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
			rowIndex++
			continue
		}
		outerMode, selectedCompound := selectedCompoundRows[rowIndex]
		if !selectedCompound {
			outerMode, err = Paul2013CompoundOuterMode(result.Model, parserRows, rowIndex)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("derive row %d FUN_10009dc0 outer mode: %w", rowIndex, err)
			}
		}
		compound, err := engine.NormalizePaul2013CompoundModelContextRow(
			ctx, result.Model, rowIndex, outerMode, contractionPrefix,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d FUN_10009dc0 compound branch: %w", rowIndex, err)
		}
		result.Model = compound.Model
		result.CompoundAttemptedRows = append(result.CompoundAttemptedRows, rowIndex)
		if compound.Applied {
			if _, selected := selectedContractionRows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected contraction row %d was consumed by FUN_10009dc0", rowIndex)
			}
			if _, selected := selectedB800Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected B800 row %d was consumed by FUN_10009dc0", rowIndex)
			}
			if _, selected := selectedC3A0Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected C3A0 row %d was consumed by FUN_10009dc0", rowIndex)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.CompoundRows = append(result.CompoundRows, rowIndex)
			rowIndex++
			continue
		}
		result.CompoundUnmatchedRows = append(result.CompoundUnmatchedRows, rowIndex)
		surface, currentFlags, err := text.Paul2013ModelContextSurface(result.Model, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("read suffix fallback row %d: %w", rowIndex, err)
		}
		suffix, err := engine.NormalizePaul2013SupportedSuffixFallbackAfter09DC0(
			ctx, surface, currentFlags,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d supported suffix fallback: %w", rowIndex, err)
		}
		result.SuffixAttemptedRows = append(result.SuffixAttemptedRows, rowIndex)
		if suffix.UnsupportedBranch || suffix.TableClassUnsupported {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("row %d reached unresolved %s suffix branch", rowIndex, suffix.Handler)
		}
		if suffix.Matched {
			if _, selected := selectedContractionRows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected contraction row %d was consumed by supported suffix handler %s", rowIndex, suffix.Handler)
			}
			if _, selected := selectedB800Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected B800 row %d was consumed by supported suffix handler %s", rowIndex, suffix.Handler)
			}
			if _, selected := selectedC3A0Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected C3A0 row %d was consumed by supported suffix handler %s", rowIndex, suffix.Handler)
			}
			result.Model, err = text.ApplyPaul2013ModelContextCodes(
				result.Model, rowIndex, suffix.Output, suffix.RowFlags,
			)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("write row %d supported suffix result: %w", rowIndex, err)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.SuffixRows = append(result.SuffixRows, rowIndex)
			rowIndex++
			continue
		}
		contractionMode, selectedContraction := selectedContractionRows[rowIndex]
		if !selectedContraction {
			contractionMode = outerMode
		}
		contraction, err := engine.NormalizePaul2013ContractionModelContextRow(
			ctx, result.Model, rowIndex, contractionMode,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d FUN_1000c040 contraction branch: %w", rowIndex, err)
		}
		result.Model = contraction.Model
		if contraction.Applied {
			if _, selected := selectedB800Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected B800 row %d was consumed by FUN_1000c040", rowIndex)
			}
			if _, selected := selectedC3A0Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected C3A0 row %d was consumed by FUN_1000c040", rowIndex)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.ContractionRows = append(result.ContractionRows, rowIndex)
			rowIndex++
			continue
		}
		result.ContractionUnmatchedRows = append(result.ContractionUnmatchedRows, rowIndex)
		b800, err := engine.ApplyPaul2013B800JoinedDictionaryModelContextRow(ctx, result.Model, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_1000b800: %w", rowIndex, err)
		}
		result.Model = b800.Model
		if b800.Applied {
			if _, selected := selectedC3A0Rows[rowIndex]; selected {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("selected C3A0 row %d was consumed by FUN_1000b800", rowIndex)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.B800Rows = append(result.B800Rows, rowIndex)
			rowIndex++
			continue
		}
		result.B800UnmatchedRows = append(result.B800UnmatchedRows, rowIndex)
		c3a0, err := engine.NormalizePaul2013C3A0ModelContextRow(ctx, result.Model, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d FUN_1000c3a0 branch: %w", rowIndex, err)
		}
		result.Model = c3a0.Model
		result.C3A0AttemptedRows = append(result.C3A0AttemptedRows, rowIndex)
		if c3a0.Applied {
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.C3A0Rows = append(result.C3A0Rows, rowIndex)
			rowIndex++
			continue
		}
		result.C3A0UnmatchedRows = append(result.C3A0UnmatchedRows, rowIndex)

		surface, currentFlags, err = text.Paul2013ModelContextSurface(result.Model, rowIndex)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("read trailing-apostrophe row %d: %w", rowIndex, err)
		}
		phoneStart := phoneRowStart + 0x23
		trailing, err := engine.NormalizePaul2013TrailingApostropheFromEmbeddedDictionary(
			ctx, surface, result.Model[phoneStart:phoneRowStart+0x64], currentFlags,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d FUN_1000c710 trailing-apostrophe branch: %w", rowIndex, err)
		}
		if trailing.Matched {
			result.Model, err = text.ApplyPaul2013ModelContextCodes(
				result.Model, rowIndex, cStringBytes(trailing.Output), trailing.RowFlags,
			)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("write row %d FUN_1000c710 result: %w", rowIndex, err)
			}
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.TrailingApostropheRows = append(result.TrailingApostropheRows, rowIndex)
			rowIndex++
			continue
		}

		generic, err := engine.NormalizePaul2013GenericToken(ctx, surface)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("normalize row %d terminal FUN_10002f10 branch: %w", rowIndex, err)
		}
		if generic.Eligible {
			writes, err := ApplyPaul2013GenericNormalizerWrites(
				result.Model[phoneStart:phoneRowStart+0x64], currentFlags, generic,
			)
			if err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("write row %d terminal FUN_10002f10 result: %w", rowIndex, err)
			}
			copy(result.Model[phoneStart:phoneRowStart+0x64], writes.Output)
			result.Model[phoneRowStart-2] = writes.Flags
			if err := applyUseMarker(rowIndex); err != nil {
				return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
			}
			result.GenericFallbackRows = append(result.GenericFallbackRows, rowIndex)
			rowIndex++
			continue
		}
		contextCodes, err := text.EncodePaul2013ContextString(surface)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("encode row %d terminal FUN_10009cd0 context: %w", rowIndex, err)
		}
		contextFlags := byte(0)
		if len(contextCodes) != 0 {
			contextFlags = 0x08
		}
		result.Model, err = text.ApplyPaul2013ModelContextCodes(
			result.Model, rowIndex, contextCodes, contextFlags,
		)
		if err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("write row %d terminal FUN_10009cd0 context: %w", rowIndex, err)
		}
		result.Model[phoneRowStart+4] = 0x45
		result.TerminalContextRows = append(result.TerminalContextRows, rowIndex)

		if err := applyUseMarker(rowIndex); err != nil {
			return Paul2013ModelContextPassResult{}, fmt.Errorf("apply row %d FUN_10007520 use-marker postpass: %w", rowIndex, err)
		}
		result.UnresolvedRows = append(result.UnresolvedRows, rowIndex)
		rowIndex++
	}

	secondPass, err := text.ApplyPaul2013ModelContextSecondPass(result.Model)
	if err != nil {
		return Paul2013ModelContextPassResult{}, fmt.Errorf("apply final counted-row pass: %w", err)
	}
	result.Model = secondPass.Model
	result.SecondPassRows = secondPass.ChangedRows
	return result, nil
}
