package text

import (
	"errors"
	"fmt"
)

// Paul2013ModelParserLoopState is the portable state surface read and written
// by FUN_100267d0. Field names retain native offsets where their meanings are
// not established. Advance and Process callbacks may update these values to
// model the native shared structures.
type Paul2013ModelParserLoopState struct {
	StopHandle             uint32
	ParserStatus           int16
	DispatchStatus         int32
	ModelStateWord0        int16
	ModelStateWord1        int16
	GeneratedRowCount      int16
	ModelStateCountAt47704 int16
	ModelStateCountAt47706 uint16
	ProgressDwordAt3ef5b   uint32
}

// Paul2013ModelParserLoopCallbacks supplies the two routines called by
// FUN_100267d0. Advance represents FUN_1002c9b0 and Process represents
// FUN_10026630. Each callback returns its raw 32-bit value and can mutate the
// shared loop state through the pointer.
type Paul2013ModelParserLoopCallbacks struct {
	Advance func(state *Paul2013ModelParserLoopState) (uint32, error)
	Process func(state *Paul2013ModelParserLoopState) (uint32, error)
}

// Paul2013ModelParserLoopResult records the native return value. The ordinary
// exit masks to the high word; the stop-handle exit inside the process retry
// returns that callback's full value.
type Paul2013ModelParserLoopResult struct {
	ReturnHighWord uint16
	ReturnValue    uint32
	LastRawResult  uint32
	ReturnedEarly  bool
	State          Paul2013ModelParserLoopState
}

// RunPaul2013ModelParserLoop ports FUN_100267d0's outer advance/process loop
// and nested retry condition. It initializes DispatchStatus to zero, increments
// the signed-short counter at native offset +0x47704 on status 0, clears the
// short at +0x47706 and dword at +0x3ef5b, and processes rows when either
// model-state word +2 is zero or the +0x47704 count equals word +0. Native
// short comparisons and the final high-word return are preserved. The parser
// and row-producing routines remain callback inputs. The early stop-handle
// return inside the retry loop preserves the full process result; ordinary
// exits return only its high word, as in the native function.
func RunPaul2013ModelParserLoop(
	state Paul2013ModelParserLoopState,
	callbacks Paul2013ModelParserLoopCallbacks,
) (Paul2013ModelParserLoopResult, error) {
	if callbacks.Advance == nil {
		return Paul2013ModelParserLoopResult{}, errors.New("Paul 2013 parser loop requires an advance callback")
	}
	if callbacks.Process == nil {
		return Paul2013ModelParserLoopResult{}, errors.New("Paul 2013 parser loop requires a process callback")
	}
	state.DispatchStatus = 0
	lastRawResult := state.StopHandle
	if state.StopHandle == 0 {
		for {
			advanceResult, err := callbacks.Advance(&state)
			if err != nil {
				return Paul2013ModelParserLoopResult{LastRawResult: lastRawResult, State: state}, fmt.Errorf("advance Paul 2013 parser loop: %w", err)
			}
			lastRawResult = uint32(uint16(state.ParserStatus)) | (advanceResult & 0xffff0000)
			if state.ParserStatus == 1 {
				break
			}
			if state.ParserStatus == 0 {
				state.ModelStateCountAt47704 = int16(uint16(state.ModelStateCountAt47704) + 1)
				state.ModelStateCountAt47706 = 0
				state.ProgressDwordAt3ef5b = 0
				if state.ModelStateWord1 == 0 || state.ModelStateCountAt47704 == state.ModelStateWord0 {
					for {
						processResult, processErr := callbacks.Process(&state)
						if processErr != nil {
							return Paul2013ModelParserLoopResult{LastRawResult: lastRawResult, State: state}, fmt.Errorf("process Paul 2013 model rows: %w", processErr)
						}
						lastRawResult = processResult
						processStatus := int16(uint16(processResult))
						if processStatus >= 0 && state.GeneratedRowCount != 0 {
							break
						}
						if state.StopHandle != 0 {
							return paul2013ModelParserLoopResult(lastRawResult, state, true), nil
						}
					}
				}
			}
			if state.StopHandle != 0 {
				break
			}
		}
	}
	return paul2013ModelParserLoopResult(lastRawResult, state, false), nil
}

func paul2013ModelParserLoopResult(
	raw uint32,
	state Paul2013ModelParserLoopState,
	returnedEarly bool,
) Paul2013ModelParserLoopResult {
	returnValue := raw & 0xffff0000
	if returnedEarly {
		returnValue = raw
	}
	return Paul2013ModelParserLoopResult{
		ReturnHighWord: uint16(returnValue >> 16),
		ReturnValue:    returnValue,
		LastRawResult:  raw,
		ReturnedEarly:  returnedEarly,
		State:          state,
	}
}
