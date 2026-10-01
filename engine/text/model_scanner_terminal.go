package text

import (
	"bytes"
	"context"
	"fmt"
)

// ScanPaul2013ModelScannerTerminalPrefix ports FUN_1005a350's mode-independent
// early returns after leading whitespace. It supplies the complete recovered
// result for two-or-more LFs (status 8) and NUL after whitespace (status 9).
// handled=false means the remaining token state machine must run.
func ScanPaul2013ModelScannerTerminalPrefix(source []byte, offset int32) (Paul2013ModelScannerResult, bool, error) {
	if bytes.IndexByte(source, 0) < 0 {
		return Paul2013ModelScannerResult{}, false, fmt.Errorf("model scanner input is not NUL-terminated")
	}
	prefix := ScanPaul2013ModelParserWhitespacePrefix(source)
	if prefix.LineFeeds < 2 && !prefix.NULTerminated {
		return Paul2013ModelScannerResult{}, false, nil
	}
	consumed := int32(prefix.ConsumedBytes)
	result := Paul2013ModelScannerResult{Field0: consumed, First: uint32(offset + consumed), Second: uint32(offset + consumed), Field14: consumed, Advance: consumed, Status: 9}
	if prefix.LineFeeds >= 2 {
		result.Status = 8
	} else if prefix.LineFeeds == 1 {
		result.PrefixFlag = 1
	}
	return result, true, nil
}

// NewPaul2013ModelScannerWithTerminalPrefixes provisions the recovered early
// paths before the supplied general scanner. Missing ordinary-token support
// returns an error; terminal input does not need the unresolved scanner.
func NewPaul2013ModelScannerWithTerminalPrefixes(fallback Paul2013ModelScanner) Paul2013ModelScanner {
	return func(ctx context.Context, source []byte, offset, mode int32, modePointer uint32) (Paul2013ModelScannerResult, error) {
		if err := ctx.Err(); err != nil {
			return Paul2013ModelScannerResult{}, err
		}
		result, handled, err := ScanPaul2013ModelScannerTerminalPrefix(source, offset)
		if err != nil || handled {
			return result, err
		}
		if fallback == nil {
			return result, fmt.Errorf("general model token scanner is unavailable")
		}
		return fallback(ctx, source, offset, mode, modePointer)
	}
}
