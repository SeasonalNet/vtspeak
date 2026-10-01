#!/usr/bin/env python3
"""Verify the captured flag-7 period-prefix state sequence."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
STAGE21 = ROOT / "tools/revkit/work/stage21"
PERIOD_COUNTS = (108, 111, 140, 52, 33, 151, 106, 58)
CALL_RE = re.compile(r"PREPROCESS_RNG_CALL n=(\d+) raw=0x([0-9a-fA-F]+) state=(\d+)")


def main() -> int:
    log_path = STAGE21 / "preprocess-rng-sequence-v2-api.log"
    log = log_path.read_text(encoding="utf-8", errors="replace")
    calls = [
        (int(match.group(1)), int(match.group(2), 16), int(match.group(3)))
        for line in log.splitlines()
        if (match := CALL_RE.search(line)) is not None
    ]
    if [(n, raw) for n, raw, _ in calls] != [(1, 1), (2, 1), (3, 1)]:
        print(f"unexpected call sequence in {log_path}: {calls}", file=sys.stderr)
        return 1

    states = [state for _, _, state in calls]
    for previous, observed in zip(states, states[1:]):
        expected = 16807 * (previous % 127773) - 2836 * (previous // 127773)
        if expected <= 0:
            expected += 2147483647
        if observed != expected:
            print(f"state transition mismatch: {previous} -> {observed}, expected {expected}", file=sys.stderr)
            return 1

    for index, state in enumerate(states, start=2):
        path = STAGE21 / f"s{index}"
        first_line = path.read_bytes().split(b"\n", maxsplit=1)[0]
        period_count = len(first_line) - len(first_line.lstrip(b"."))
        expected_periods = PERIOD_COUNTS[state % len(PERIOD_COUNTS)]
        if period_count != expected_periods or first_line[period_count : period_count + 1] != b" ":
            print(
                f"{path}: observed {period_count} periods, expected {expected_periods}",
                file=sys.stderr,
            )
            return 1
        print(f"call={index - 1} state={state} table_index={state % 8} periods={period_count}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
