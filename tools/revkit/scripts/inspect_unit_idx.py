#!/usr/bin/env python3
"""Read-only structural inspector for the VoiceText 2013 unit-index layout."""

from __future__ import annotations

import argparse
import json
import struct
from pathlib import Path
from typing import BinaryIO


def read_exact(stream: BinaryIO, size: int, label: str) -> bytes:
    value = stream.read(size)
    if len(value) != size:
        raise ValueError(f"truncated {label}: expected {size} bytes, got {len(value)}")
    return value


def read_u16(stream: BinaryIO, label: str) -> int:
    return struct.unpack("<H", read_exact(stream, 2, label))[0]


def read_u32(stream: BinaryIO, label: str) -> int:
    return struct.unpack("<I", read_exact(stream, 4, label))[0]


def inspect(path: Path) -> dict[str, object]:
    raw = path.read_bytes()
    with path.open("rb") as stream:
        header_length = read_exact(stream, 1, "header length")[0]
        header = read_exact(stream, header_length, "version header")
        if not header.startswith(b"ver.2013\0VoiceText-Eng\0"):
            raise ValueError(f"{path}: unsupported version header {header!r}")

        table_count = read_u16(stream, "bank table count")
        banks: list[dict[str, int | str]] = []
        for index in range(table_count):
            name_length = read_u16(stream, f"bank {index} name length")
            name = read_exact(stream, name_length, f"bank {index} name").decode("ascii")
            tag = read_exact(stream, 1, f"bank {index} tag")[0]
            banks.append({"name": name, "tag": tag})

        unit_count = read_u32(stream, "unit count")
        opaque_stride = read_u16(stream, "opaque per-unit block stride")
        block_start = stream.tell()
        opaque_block_bytes = unit_count * opaque_stride
        read_exact(stream, opaque_block_bytes, "opaque per-unit blocks")

        # The parser bulk-reads these columns after skipping the opaque block.
        columns: list[tuple[str, int]] = [
            ("column_a", 1),
            ("unit_signature", 7),
            ("column_b", 1),
        ]
        for group in range(1, 4):
            columns.extend(
                (
                    (f"group_{group}_word", 2),
                    (f"group_{group}_byte_a", 1),
                    (f"group_{group}_byte_b", 1),
                )
            )

        column_width = sum(width for _, width in columns)
        read_exact(stream, unit_count * column_width, "unit columns")
        if stream.read(1):
            raise ValueError(f"{path}: trailing bytes after parsed unit columns")

    expected_size = block_start + unit_count * (opaque_stride + column_width)
    if expected_size != len(raw):
        raise ValueError(f"{path}: expected {expected_size} bytes, found {len(raw)}")

    return {
        "file": str(path),
        "version": "ver.2013",
        "producer": "VoiceText-Eng",
        "banks": banks,
        "unit_count": unit_count,
        "opaque_block_stride_bytes": opaque_stride,
        "opaque_block_start": block_start,
        "column_layout": columns,
        "column_bytes_per_unit": column_width,
        "expected_size_bytes": expected_size,
        "actual_size_bytes": len(raw),
        "size_matches": True,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("files", nargs="+", type=Path, help="2013 unit-index files to inspect")
    args = parser.parse_args()
    results = [inspect(path) for path in args.files]
    print(json.dumps(results, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
