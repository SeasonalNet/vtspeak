#!/usr/bin/env python3
"""Verify the Stage 21 exhaustive MakeCsv two-byte input capture."""

from __future__ import annotations

import argparse
import struct
from pathlib import Path


CASES = 256 * 256
ROW_SIZE = 14
ROWS_SIZE = CASES * ROW_SIZE
CAPTURE_SIZE = ROWS_SIZE + 6


def expected_row(first: int, second: int) -> tuple[int, bytes]:
    if first == 0:
        return -1, b'"' + b"1" * 6 + b"\0"

    if first & 0x80:
        payload = bytes((first,)) + (bytes((second,)) if second else b"")
    else:
        payload = bytes((first, first)) if first == 0x22 else bytes((first,))
        if second:
            payload += bytes((second, second)) if second == 0x22 else bytes((second,))

    encoded = b'"' + payload + b'"\0'
    if len(encoded) > 7:
        raise ValueError(f"expected output exceeds capacity: {first:02x} {second:02x}")
    return 1, encoded + b"1" * (7 - len(encoded)) + b"\0"


def verify(path: Path) -> None:
    data = path.read_bytes()
    if len(data) != CAPTURE_SIZE:
        raise ValueError(f"capture size is {len(data)}, expected {CAPTURE_SIZE}")

    empty = 0
    nonempty = 0
    for first in range(256):
        for second in range(256):
            offset = (first * 256 + second) * ROW_SIZE
            ret = struct.unpack_from("<i", data, offset)[0]
            output = data[offset + 4 : offset + 12]
            prefix, post = data[offset + 12 : offset + 14]
            expected_ret, expected_output = expected_row(first, second)
            if ret != expected_ret or output != expected_output:
                raise ValueError(
                    f"mismatch at {first:02x} {second:02x}: ret={ret}, "
                    f"output={output.hex()}, expected_ret={expected_ret}, "
                    f"expected_output={expected_output.hex()}"
                )
            if (prefix, post) != (0x5A, 0xA5):
                raise ValueError(
                    f"guard mismatch at {first:02x} {second:02x}: "
                    f"prefix={prefix:02x}, post={post:02x}"
                )
            if first == 0:
                empty += 1
            else:
                nonempty += 1

    if data[ROWS_SIZE : ROWS_SIZE + 2] != b"\xff\xff":
        raise ValueError("last input pair metadata is not ff ff")
    completed = struct.unpack_from("<I", data, ROWS_SIZE + 2)[0]
    if completed != CASES:
        raise ValueError(f"completed count is {completed}, expected {CASES}")

    print(
        f"verified_cases={CASES} effective_empty_rows={empty} "
        f"effective_nonempty_rows={nonempty} guards=intact "
        f"last_pair=ffff completed={completed}"
    )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("capture", type=Path)
    args = parser.parse_args()
    verify(args.capture)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
