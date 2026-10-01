package text

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
)

// Paul2013PrimaryParserHandler is one native row handler. Address keys retain
// their native identity; the scanner and state may be changed by the handler.
// markedRows is the native local low short also passed to FUN_100544f0.
type Paul2013PrimaryParserHandler func(context.Context, []byte, *Paul2013ModelScannerResult, []byte, int32, int16) (int32, error)

type Paul2013PrimaryModelParserResult struct {
	State         []byte
	ReturnValue   int32
	ConsumedBytes int
	MarkedRows    int16
	HandlerCalls  int
	ScannerCalls  int
}

// RunPaul2013PrimaryModelParser ports FUN_1003d3d0's caller control flow.
// Native handler bodies and the shared scanner remain explicit dependencies.
// Missing reached handlers return errors, rather than being treated as misses.
// State is already initialized by FUN_1003d350. The ordered fallback cascade,
// candidate rollback, marked-row boundary and 0x10033760 merge are handled here.
func RunPaul2013PrimaryModelParser(ctx context.Context, input, state []byte, scan Paul2013ModelScanner, handlers map[uint32]Paul2013PrimaryParserHandler) (Paul2013PrimaryModelParserResult, error) {
	result := Paul2013PrimaryModelParserResult{State: append([]byte(nil), state...)}
	nul := bytes.IndexByte(input, 0)
	if nul < 0 || len(state) < 0x39ec {
		return result, fmt.Errorf("primary parser input/state is truncated")
	}
	position := 0
	count := func() (int, error) {
		n := int(int16(binary.LittleEndian.Uint16(result.State[:2])))
		if n < 0 || n > 100 {
			return 0, fmt.Errorf("primary parser count %d outside capacity", n)
		}
		return n, nil
	}
	finish := func() (Paul2013PrimaryModelParserResult, error) {
		n, err := count()
		result.ReturnValue = -1
		if n > 0 {
			result.ReturnValue = 1
		}
		result.ConsumedBytes = position
		return result, err
	}
	rollback := func(n int) error {
		updated, err := RollbackPaul2013ModelParserCandidateRow(result.State, n)
		if err == nil {
			result.State = updated
		}
		return err
	}
	move := func(amount int32) error {
		next := int64(position) + int64(amount)
		if amount <= 0 || next > int64(nul) {
			return fmt.Errorf("primary parser advance %d outside source", amount)
		}
		position = int(next)
		result.ConsumedBytes = position
		return nil
	}
	call := func(address uint32, token *Paul2013ModelScannerResult) (int32, error) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		handler := handlers[address]
		if handler == nil {
			return 0, fmt.Errorf("primary parser handler %#x is unavailable", address)
		}
		result.HandlerCalls++
		value, err := handler(ctx, result.State, token, input[position:nul+1], int32(position), result.MarkedRows)
		if err != nil {
			return value, fmt.Errorf("primary parser handler %#x: %w", address, err)
		}
		_, err = count()
		return value, err
	}
	boundary := func(committed int, setMode bool) (Paul2013PrimaryModelParserResult, error) {
		row := 0x14 + int(result.MarkedRows)*0x94
		if row+4 > len(result.State) {
			return result, fmt.Errorf("marked parser boundary lacks next row")
		}
		start := int32(binary.LittleEndian.Uint32(result.State[row:]))
		delta := start - int32(position)
		if delta < 1 {
			delta = 0
		}
		if err := rollback(committed); err != nil {
			return result, err
		}
		binary.LittleEndian.PutUint32(result.State[4:8], uint32(int32(position)+delta))
		if setMode && binary.LittleEndian.Uint16(result.State[0xc:]) == 0 {
			binary.LittleEndian.PutUint16(result.State[0xc:], 1)
		}
		return finish()
	}
	// Native calls between FUN_1004d990 and the terminal FUN_1003a880.
	fallbacks := []uint32{0x1003d290, 0x1005c890, 0x100344e0, 0x10036940, 0x10058490, 0x100570c0, 0x10052e80, 0x10053440, 0x1003e4a0, 0x10033760, 0x1003c1a0, 0x1003aad0, 0x100420f0, 0x10031d40, 0x10040ef0, 0x100544f0, 0x10063780, 0x1003a880}
	for {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		committed, err := count()
		if err != nil {
			return result, err
		}
		binary.LittleEndian.PutUint32(result.State[4:8], uint32(position))
		var token Paul2013ModelScannerResult
		accept := func(value int32, setMode bool) (bool, error) {
			if value == 0 {
				if err := rollback(committed); err != nil {
					return false, err
				}
				return true, nil
			}
			n, err := count()
			if err != nil {
				return false, err
			}
			if result.MarkedRows != 0 && n > int(result.MarkedRows) {
				return true, nil
			}
			if err := move(value); err != nil {
				return false, err
			}
			if setMode && binary.LittleEndian.Uint16(result.State[0xc:]) != 0 {
				if result.MarkedRows != 0 {
					return true, nil
				}
				result.MarkedRows = int16(n)
				if n < 1 {
					return true, nil
				}
			}
			return false, nil
		}
		if binary.LittleEndian.Uint16(result.State[0xc:]) != 0 {
			value, err := call(0x1005dd20, &token)
			if err != nil {
				return result, err
			}
			if value >= 0 {
				n, _ := count()
				if value > 0 && result.MarkedRows != 0 && n > int(result.MarkedRows) {
					return boundary(committed, true)
				}
				stop, err := accept(value, true)
				if err != nil {
					return result, err
				}
				if stop {
					return finish()
				}
				continue
			}
			if err := rollback(committed); err != nil {
				return result, err
			}
		}
		if scan == nil {
			return result, fmt.Errorf("primary model parser scanner is unavailable")
		}
		result.ScannerCalls++
		token, err = scan(ctx, input[position:nul+1], int32(position), 0, binary.LittleEndian.Uint32(result.State[0x39e8:]))
		if err != nil {
			return result, err
		}
		if token.PrefixFlag == 1 {
			n, _ := count()
			if n > 0 {
				field := 0x14 + (n-1)*0x94 + 0x2c
				if field+4 > len(result.State) {
					return result, fmt.Errorf("scanner row type field is truncated")
				}
				if binary.LittleEndian.Uint32(result.State[field:]) == 0 {
					binary.LittleEndian.PutUint32(result.State[field:], 1)
				}
			}
		}
		if token.Status >= 1 && token.Status <= 3 {
			value, err := call(0x1005deb0, &token)
			if err != nil {
				return result, err
			}
			if value >= 0 {
				n, _ := count()
				if value > 0 && result.MarkedRows != 0 && n > int(result.MarkedRows) {
					return boundary(committed, false)
				}
				stop, err := accept(value, false)
				if err != nil {
					return result, err
				}
				if stop {
					return finish()
				}
				continue
			}
			if err := rollback(committed); err != nil {
				return result, err
			}
		}
		value, err := call(0x10051a00, &token)
		if err != nil {
			return result, err
		}
		if value < 0 {
			magnitude := -int64(value)
			if magnitude > 0x7fffffff {
				return result, fmt.Errorf("terminal handler advance overflows")
			}
			if err := move(int32(magnitude)); err != nil {
				return result, err
			}
			continue
		}
		if value > 0 {
			if result.MarkedRows == 0 {
				binary.LittleEndian.PutUint32(result.State[4:8], uint32(position))
				if token.Field14 != 0 {
					if err := move(value); err != nil {
						return result, err
					}
					binary.LittleEndian.PutUint32(result.State[4:8], uint32(position))
				}
				n, _ := count()
				result.MarkedRows = int16(n)
				if n < 1 {
					return finish()
				}
				continue
			}
			if token.Field14 == 0 {
				if err := rollback(committed); err != nil {
					return result, err
				}
				binary.LittleEndian.PutUint32(result.State[4:8], uint32(position))
				return finish()
			}
			if err := move(value); err != nil {
				return result, err
			}
			continue
		}
		if err := rollback(committed); err != nil {
			return result, err
		}
		value, err = call(0x1004d990, &token)
		if err != nil {
			return result, err
		}
		if value >= 0 {
			n, _ := count()
			if value == 0 || result.MarkedRows != 0 && n > int(result.MarkedRows) {
				if err := rollback(committed); err != nil {
					return result, err
				}
				if value > 0 {
					binary.LittleEndian.PutUint32(result.State[4:8], uint32(position))
				}
				return finish()
			}
			if err := move(value); err != nil {
				return result, err
			}
			continue
		}
		for _, address := range fallbacks {
			if err := rollback(committed); err != nil {
				return result, err
			}
			value, err = call(address, &token)
			if err != nil {
				return result, err
			}
			if value < 0 {
				continue
			}
			n, _ := count()
			if value > 0 && result.MarkedRows != 0 && n > int(result.MarkedRows) {
				if address == 0x10033760 {
					marked := 0x14 + (int(result.MarkedRows)-1)*0x94
					last := 0x14 + (n-1)*0x94
					if last+0x30 > len(result.State) {
						return result, fmt.Errorf("parser merge rows are truncated")
					}
					if binary.LittleEndian.Uint32(result.State[marked+0x28:]) == 0x12 && binary.LittleEndian.Uint32(result.State[marked+0x2c:]) == 2 && binary.LittleEndian.Uint32(result.State[last+0x28:]) == 0x12 && !(binary.LittleEndian.Uint16(result.State[marked+0x1e:]) == 0x1f && binary.LittleEndian.Uint16(result.State[last+0x1e:]) == 0x1f) {
						binary.LittleEndian.PutUint32(result.State[marked+0x2c:], 0)
						result.MarkedRows = 0
					} else {
						return boundary(committed, false)
					}
				} else {
					return boundary(committed, false)
				}
			}
			stop, err := accept(value, address == 0x1005c890)
			if err != nil {
				return result, err
			}
			if stop {
				return finish()
			}
			break
		}
		if value < 0 {
			if err := rollback(committed); err != nil {
				return result, err
			}
		}
	}
}

// DispatchPaul2013ModelParserWithHandlers supplies both recovered inner
// parser loops while preserving explicit parser overrides and finalization.
func DispatchPaul2013ModelParserWithHandlers(ctx context.Context, input, state, modeProbe []byte, callbacks Paul2013ModelParserCallbacks, scan Paul2013ModelScanner, handlers map[uint32]Paul2013PrimaryParserHandler) (Paul2013ModelParserDispatchResult, error) {
	if callbacks.Primary == nil {
		callbacks.Primary = func(input, state []byte) (int32, error) {
			result, err := RunPaul2013PrimaryModelParser(ctx, input, state, scan, handlers)
			if err == nil {
				copy(state, result.State)
			}
			return result.ReturnValue, err
		}
	}
	return DispatchPaul2013ModelParserWithAlternateScanner(ctx, input, state, modeProbe, callbacks, scan)
}
