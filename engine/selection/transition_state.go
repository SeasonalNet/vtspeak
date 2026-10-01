package selection

import (
	"encoding/binary"
	"fmt"

	"vtspeak/engine/internal/paul2013tables"
)

const (
	paul2013TransitionEligibilityArenaBase   = 0x65c
	paul2013TransitionEligibilityGroupStride = 0x1e0
)

// ReadPaul2013TransitionEligibilityShorts reads FUN_10024680's signed-short
// per-phone gates from the workspace arena. Native indexing is
// base+0x65c+(groupOrdinal*0x1e0+phoneIndex)*2. groupOrdinal is explicit
// because the descriptor's +0x0c value, not compact output position, selects
// the arena group.
func ReadPaul2013TransitionEligibilityShorts(
	arena []byte,
	groupOrdinal int,
	phoneCount int,
) ([]int16, error) {
	if groupOrdinal < 0 {
		return nil, fmt.Errorf("transition eligibility group ordinal %d is negative", groupOrdinal)
	}
	if phoneCount < 0 || phoneCount > paul2013TransitionEligibilityGroupStride {
		return nil, fmt.Errorf("transition eligibility phone count %d is outside 0..%d", phoneCount, paul2013TransitionEligibilityGroupStride)
	}
	groupStrideBytes := paul2013TransitionEligibilityGroupStride * 2
	maxInt := int(^uint(0) >> 1)
	if groupOrdinal > (maxInt-paul2013TransitionEligibilityArenaBase)/groupStrideBytes {
		return nil, fmt.Errorf("transition eligibility group ordinal %d overflows arena indexing", groupOrdinal)
	}
	groupOffset := paul2013TransitionEligibilityArenaBase + groupOrdinal*groupStrideBytes
	if groupOffset > len(arena) || phoneCount*2 > len(arena)-groupOffset {
		return nil, fmt.Errorf("transition eligibility group %d needs bytes [%#x,%#x), arena has %d", groupOrdinal, groupOffset, groupOffset+phoneCount*2, len(arena))
	}
	result := make([]int16, phoneCount)
	for phoneIndex := range result {
		result[phoneIndex] = int16(binary.LittleEndian.Uint16(arena[groupOffset+phoneIndex*2:]))
	}
	return result, nil
}

// BuildPaul2013TransitionPhoneInputsFromEligibilityArena applies the native
// positive-short gate to one already-mapped group of phone inputs. It checks
// that phone order matches the arena order and preserves all record and query
// inputs supplied by the caller.
func BuildPaul2013TransitionPhoneInputsFromEligibilityArena(
	phones []Paul2013TransitionPhoneInput,
	arena []byte,
	groupOrdinal int,
) ([]Paul2013TransitionPhoneInput, error) {
	shorts, err := ReadPaul2013TransitionEligibilityShorts(arena, groupOrdinal, len(phones))
	if err != nil {
		return nil, err
	}
	result := append([]Paul2013TransitionPhoneInput(nil), phones...)
	for index := range result {
		if int(result[index].PhoneIndex) != index {
			return nil, fmt.Errorf("transition phone input %d has phone index %d; eligibility arena order requires %d", index, result[index].PhoneIndex, index)
		}
		result[index].QueryEligible = shorts[index] > 0
	}
	return result, nil
}

// BuildPaul2013TransitionContextStatesFromEligibilityArena composes the
// native per-phone short reads with FUN_10024680's phone-field writes and
// compact one-row/two-row state advancement. The caller supplies the already
// mapped phone records, pointer targets, and position-query results.
func BuildPaul2013TransitionContextStatesFromEligibilityArena(
	phones []Paul2013TransitionPhoneInput,
	arena []byte,
	groupOrdinal int,
) ([]Paul2013TransitionContextState, error) {
	prepared, err := BuildPaul2013TransitionPhoneInputsFromEligibilityArena(phones, arena, groupOrdinal)
	if err != nil {
		return nil, fmt.Errorf("read native transition eligibility: %w", err)
	}
	states, err := BuildPaul2013TransitionContextStates(prepared)
	if err != nil {
		return nil, fmt.Errorf("build native transition context states: %w", err)
	}
	return states, nil
}

// Paul2013TransitionContextState is the six-byte state record read by
// FUN_10018c80. Its weight-row choice uses byte 0 and signed byte 5. The
// recovered phone-state writer also fills bytes 1 and 2; bytes 3 and 4 are
// preserved for later query/fallback writers.
type Paul2013TransitionContextState [6]byte

// Paul2013ContextModesFromTransitionStates extracts the byte +4 mode written
// by FUN_10024060/FUN_100242a0: accepted whole-position rows use 0, and the
// two fallback rows use 1 and 2. The caller must keep these rows aligned with
// the positions passed to continuity ranking.
func Paul2013ContextModesFromTransitionStates(states []Paul2013TransitionContextState) ([]byte, error) {
	modes := make([]byte, len(states))
	for index, state := range states {
		if state[4] > 2 {
			return nil, fmt.Errorf("transition context state %d has unsupported mode %d at +4", index, state[4])
		}
		modes[index] = state[4]
	}
	return modes, nil
}

// Paul2013TransitionPhoneInput supplies one phone visited by FUN_10024680.
// QueryEligible reflects the positive per-phone short at workspace offset
// +0x65c; BuildPaul2013TransitionPhoneInputsFromEligibilityArena derives it
// from the native short array. QueryResult is the result of
// QueryPaul2013Position for an eligible phone; its accepted/fallback state
// controls the native one-row/two-row transition-state write.
type Paul2013TransitionPhoneInput struct {
	RecordOrdinal byte
	PhoneIndex    byte
	Record        []byte
	PointedRows   []byte
	QueryEligible bool
	QueryResult   *Paul2013PositionQueryResult
}

// Paul2013TransitionPeriodRowsResolver maps the native pointer stored at
// record +0x08 to its 0x1e-byte period-row target. The pointer value is a
// vendor-process address and must be resolved by the caller's address space.
type Paul2013TransitionPeriodRowsResolver func(pointer uint32) ([]byte, error)

// BuildPaul2013TransitionContextStatesFromRecordGroup maps prepared native
// records to FUN_10024680's phone iteration order, resolves each record's
// period-row pointer, applies the workspace eligibility shorts, and writes
// compact transition states from the supplied per-phone position-query
// results. recordOrdinalBase is the native byte conversion of the descriptor's
// +0x0c value; each record's loop ordinal is added with byte conversion.
// queryResults is flattened by record order and then phone index, with nil
// entries for ineligible phones.
func BuildPaul2013TransitionContextStatesFromRecordGroup(
	records [][]byte,
	recordOrdinalBase byte,
	arena []byte,
	groupOrdinal int,
	queryResults []*Paul2013PositionQueryResult,
	resolvePeriodRows Paul2013TransitionPeriodRowsResolver,
) ([]Paul2013TransitionContextState, error) {
	if resolvePeriodRows == nil {
		return nil, fmt.Errorf("transition period-row pointer resolver is nil")
	}
	if groupOrdinal < 0 {
		return nil, fmt.Errorf("transition eligibility group ordinal %d is negative", groupOrdinal)
	}
	maxInt := int(^uint(0) >> 1)
	phoneCount := 0
	for recordIndex, record := range records {
		if recordIndex > maxInt-groupOrdinal {
			return nil, fmt.Errorf("transition eligibility group ordinal %d overflows at record %d", groupOrdinal, recordIndex)
		}
		if len(record) < 0x96 {
			return nil, fmt.Errorf("transition record %d has %d bytes, need 150 for phone count", recordIndex, len(record))
		}
		phoneCount += int(record[0x95])
		if phoneCount > paul2013TransitionEligibilityGroupStride {
			return nil, fmt.Errorf("transition record group has %d phones, exceeding native eligibility stride %d", phoneCount, paul2013TransitionEligibilityGroupStride)
		}
	}
	if len(queryResults) != phoneCount {
		return nil, fmt.Errorf("received %d position-query results for %d record-group phones", len(queryResults), phoneCount)
	}
	phones := make([]Paul2013TransitionPhoneInput, 0, phoneCount)
	queryIndex := 0
	for recordIndex, record := range records {
		shorts, readErr := ReadPaul2013TransitionEligibilityShorts(arena, groupOrdinal+recordIndex, int(record[0x95]))
		if readErr != nil {
			return nil, fmt.Errorf("read native transition eligibility for record %d: %w", recordIndex, readErr)
		}
		pointedRows := []byte(nil)
		if record[0x95] != 0 {
			pointer := binary.LittleEndian.Uint32(record[0x08:0x0c])
			resolvedRows, resolveErr := resolvePeriodRows(pointer)
			if resolveErr != nil {
				return nil, fmt.Errorf("resolve period rows for transition record %d pointer %#x: %w", recordIndex, pointer, resolveErr)
			}
			pointedRows = resolvedRows
		}
		for phoneIndex := 0; phoneIndex < int(record[0x95]); phoneIndex++ {
			if shorts[phoneIndex] > 0 && queryResults[queryIndex] == nil {
				return nil, fmt.Errorf("eligible transition phone %d has no position-query result", queryIndex)
			}
			if shorts[phoneIndex] <= 0 && queryResults[queryIndex] != nil {
				return nil, fmt.Errorf("ineligible transition phone %d has a position-query result", queryIndex)
			}
			phones = append(phones, Paul2013TransitionPhoneInput{
				RecordOrdinal: byte((int(recordOrdinalBase) + recordIndex) & 0xff),
				PhoneIndex:    byte(phoneIndex),
				Record:        record,
				PointedRows:   pointedRows,
				QueryEligible: shorts[phoneIndex] > 0,
				QueryResult:   queryResults[queryIndex],
			})
			queryIndex++
		}
	}
	return BuildPaul2013TransitionContextStates(phones)
}

// BuildPaul2013TransitionContextStates ports the six-byte row writes and
// compact output advancement in FUN_10024680. Noneligible phones do not
// contribute to the returned compact prefix; their scratch-slot write is
// overwritten if a later eligible phone uses that slot. Eligible accepted
// phones emit one row; eligible rejected phones emit the two rows produced by
// FUN_100242a0. The query result is validated and supplies the acceptance
// decision directly. Callers may use the arena adapter to derive eligibility
// from the native short array; its upstream writer remains separate.
func BuildPaul2013TransitionContextStates(
	phones []Paul2013TransitionPhoneInput,
) ([]Paul2013TransitionContextState, error) {
	states := make([]Paul2013TransitionContextState, 0, len(phones)*2)
	for inputIndex, phone := range phones {
		base, err := ApplyPaul2013TransitionStatePhoneFields(
			Paul2013TransitionContextState{},
			phone.RecordOrdinal,
			phone.PhoneIndex,
			phone.Record,
			phone.PointedRows,
		)
		if err != nil {
			return nil, fmt.Errorf("transition phone input %d: %w", inputIndex, err)
		}
		if !phone.QueryEligible {
			if phone.QueryResult != nil {
				return nil, fmt.Errorf("transition phone input %d has a query result although its native eligibility gate is false", inputIndex)
			}
			continue
		}
		if phone.QueryResult == nil {
			return nil, fmt.Errorf("transition phone input %d is query-eligible but has no position-query result", inputIndex)
		}
		accepted := phone.QueryResult.WholePosition.Selection.Accepted
		if accepted && (phone.QueryResult.ReturnCount != 1 || phone.QueryResult.UsedFallback) {
			return nil, fmt.Errorf("transition phone input %d has inconsistent accepted position-query state", inputIndex)
		}
		if !accepted && (phone.QueryResult.ReturnCount != 2 || !phone.QueryResult.UsedFallback) {
			return nil, fmt.Errorf("transition phone input %d has inconsistent fallback position-query state", inputIndex)
		}
		if accepted {
			if len(phone.Record) <= 0x92 {
				return nil, fmt.Errorf("transition phone input %d has %d record bytes, need query state at 0x92", inputIndex, len(phone.Record))
			}
			base[3] = phone.Record[0x92]
			base[4] = 0
			states = append(states, base)
			continue
		}
		if len(phone.Record) <= 0x92 {
			return nil, fmt.Errorf("transition phone input %d has %d record bytes, need fallback state at 0x92", inputIndex, len(phone.Record))
		}
		base[3] = phone.Record[0x92]
		base[4] = 1
		states = append(states, base)
		second := base
		second[3]++
		second[4] = 2
		states = append(states, second)
	}
	return states, nil
}

// Paul2013PhoneRowHasTransitionFlag ports FUN_10017010: it reads byte +2 of
// the seven-byte phone row and returns true when the DLL table at RVA 0x7ba44
// contains zero for that byte. The byte's meaning remains opaque.
func Paul2013PhoneRowHasTransitionFlag(phoneRow []byte) (bool, error) {
	if len(phoneRow) < 3 {
		return false, fmt.Errorf("phone row has %d bytes, need 3 for FUN_10017010", len(phoneRow))
	}
	return paul2013tables.ByteTable1007BA44(phoneRow[2]) == 0, nil
}

// ReadPaul2013TransitionPeriodPhoneSpans extracts the period-span byte at
// +0x1d from each 0x1e-byte row reached through a native record's pointer at
// +0x08. The pointer target is caller supplied because native addresses do not
// survive the independent record representation.
func ReadPaul2013TransitionPeriodPhoneSpans(record, pointedRows []byte) ([]byte, error) {
	if len(record) <= 0x94 {
		return nil, fmt.Errorf("native record has %d bytes, need count byte at 0x94", len(record))
	}
	periodCount := int(record[0x94])
	const rowStride = 0x1e
	if periodCount > len(pointedRows)/rowStride {
		return nil, fmt.Errorf("record declares %d period rows but pointer data has %d complete rows", periodCount, len(pointedRows)/rowStride)
	}
	spans := make([]byte, periodCount)
	for index := range spans {
		spans[index] = pointedRows[index*rowStride+0x1d]
	}
	return spans, nil
}

// ApplyPaul2013TransitionStatePhoneFields ports the per-phone state writes in
// FUN_10024680. recordOrdinal is the caller's native byte conversion of
// local14+sVar4. pointedRows is the resolved target of the record pointer at
// +0x08. Fields 3 and 4 are preserved because later query/fallback stages
// write them.
func ApplyPaul2013TransitionStatePhoneFields(
	state Paul2013TransitionContextState,
	recordOrdinal byte,
	phoneIndex byte,
	record []byte,
	pointedRows []byte,
) (Paul2013TransitionContextState, error) {
	if len(record) < 0x96 {
		return Paul2013TransitionContextState{}, fmt.Errorf("native record has %d bytes, need 150 for phone state fields", len(record))
	}
	if phoneIndex >= record[0x95] {
		return Paul2013TransitionContextState{}, fmt.Errorf("phone index %d is outside native phone count %d", phoneIndex, record[0x95])
	}
	periodPhoneSpans, err := ReadPaul2013TransitionPeriodPhoneSpans(record, pointedRows)
	if err != nil {
		return Paul2013TransitionContextState{}, err
	}
	periodCount := len(periodPhoneSpans)
	phoneRowStart := 0x96 + int(phoneIndex)*7
	if phoneRowStart+7 > len(record) {
		return Paul2013TransitionContextState{}, fmt.Errorf("phone row %d requires bytes through %#x, record has %d", phoneIndex, phoneRowStart+7, len(record))
	}
	flag, err := Paul2013PhoneRowHasTransitionFlag(record[phoneRowStart : phoneRowStart+7])
	if err != nil {
		return Paul2013TransitionContextState{}, err
	}
	state[0] = recordOrdinal
	periodIndex := 0
	periodStart := int16(0)
	for ; periodIndex < periodCount; periodIndex++ {
		span := int(periodPhoneSpans[periodIndex])
		if periodStart <= int16(phoneIndex) && int(phoneIndex) < int(periodStart)+span {
			break
		}
		periodStart = int16(int32(periodStart) + int32(span))
	}
	state[1] = byte(periodIndex)
	state[2] = phoneIndex
	state[5] = 0
	if flag {
		state[5] = 1
	}
	return state, nil
}

// Paul2013TransitionWeightRow derives the transition coefficient row from
// the current and previous six-byte state records. The returned row is in
// the range 0 through 4 and follows the signed-byte branches in
// FUN_10018c80.
func Paul2013TransitionWeightRow(
	current Paul2013TransitionContextState,
	previous Paul2013TransitionContextState,
) uint8 {
	currentState := int8(current[5])
	previousState := int8(previous[5])
	if currentState == 0 {
		if previousState == 0 {
			if current[0] != previous[0] {
				return 4
			}
			return 0
		}
		if previousState > 0 {
			return 1
		}
	}
	if currentState < 1 || previousState != 0 {
		return 3
	}
	return 2
}
