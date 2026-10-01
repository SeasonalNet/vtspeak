#!/usr/bin/env python3
"""Check paired MakeInfo captures for stable ASCII and binary headers."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
CAPTURE_DIR = ROOT / "tools/revkit/work/stage16"
ASCII_HEADER = ("VTDTTS ASCII", "3", "4")
BINARY_PREFIX = b"VTDTTS BINARY\0\x03\x04"


def main() -> int:
    captures = sorted(CAPTURE_DIR.glob("*.asc.dtt"))
    failures: list[str] = []
    binary_headers: set[bytes] = set()
    ascii_headers: set[tuple[str, ...]] = set()
    paired = 0

    for ascii_path in captures:
        ascii_header = tuple(ascii_path.read_text(encoding="ascii").splitlines()[:3])
        ascii_headers.add(ascii_header)
        binary_path = ascii_path.with_name(ascii_path.name.replace(".asc.dtt", ".bin.dtt"))
        if not binary_path.is_file():
            failures.append(f"{ascii_path.name}: missing binary partner")
            continue

        paired += 1
        binary_header = binary_path.read_bytes()[: len(BINARY_PREFIX)]
        binary_headers.add(binary_header)
        if ascii_header != ASCII_HEADER:
            failures.append(f"{ascii_path.name}: unexpected ASCII header {ascii_header!r}")
        if binary_header != BINARY_PREFIX:
            failures.append(f"{binary_path.name}: unexpected binary prefix {binary_header!r}")

    print(f"ASCII captures checked: {len(captures)}")
    print(f"binary pairs checked: {paired}")
    print(f"distinct ASCII headers: {sorted(ascii_headers)!r}")
    print(f"distinct binary signature/field prefixes: {sorted(h.hex() for h in binary_headers)!r}")
    for failure in failures:
        print(f"mismatch: {failure}", file=sys.stderr)
    return 1 if failures or not captures else 0


if __name__ == "__main__":
    raise SystemExit(main())
