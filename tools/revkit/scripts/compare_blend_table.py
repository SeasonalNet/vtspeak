#!/usr/bin/env python3
"""Compare the DLL's blend coefficients with the documented analytic curve."""

from __future__ import annotations

import argparse
import json
import math
import struct
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
TABLE_RVA = 0x6D1B0
TABLE_ENTRIES = 8192


def table_file_offset(data: bytes, rva: int) -> int:
    if data[:2] != b"MZ":
        raise ValueError("input is not a PE image")
    pe_offset = struct.unpack_from("<I", data, 0x3C)[0]
    if data[pe_offset : pe_offset + 4] != b"PE\0\0":
        raise ValueError("PE signature is missing")
    section_count = struct.unpack_from("<H", data, pe_offset + 6)[0]
    optional_size = struct.unpack_from("<H", data, pe_offset + 20)[0]
    section_offset = pe_offset + 24 + optional_size
    for index in range(section_count):
        offset = section_offset + index * 40
        name = data[offset : offset + 8].rstrip(b"\0").decode("ascii", errors="replace")
        virtual_size, virtual_address, raw_size, raw_offset = struct.unpack_from(
            "<IIII", data, offset + 8
        )
        extent = max(virtual_size, raw_size)
        if virtual_address <= rva < virtual_address + extent:
            file_offset = raw_offset + rva - virtual_address
            if file_offset + TABLE_ENTRIES * 4 > raw_offset + raw_size:
                raise ValueError(f"coefficient table extends past section {name!r}")
            return file_offset
    raise ValueError(f"RVA 0x{rva:x} is outside the PE sections")


def compare(dll_path: Path) -> dict[str, int | str]:
    data = dll_path.read_bytes()
    file_offset = table_file_offset(data, TABLE_RVA)
    table = data[file_offset : file_offset + TABLE_ENTRIES * 4]
    if len(table) != TABLE_ENTRIES * 4:
        raise ValueError("coefficient table is truncated")

    mismatches = 0
    max_ulp_difference = 0
    for index in range(TABLE_ENTRIES):
        native_bits = struct.unpack_from("<I", table, index * 4)[0]
        formula = math.sin(math.pi * index / TABLE_ENTRIES) ** 2
        formula_bits = struct.unpack("<I", struct.pack("<f", formula))[0]
        if native_bits != formula_bits:
            mismatches += 1
            max_ulp_difference = max(max_ulp_difference, abs(native_bits - formula_bits))

    return {
        "dll": str(dll_path),
        "table_rva": f"0x{TABLE_RVA:x}",
        "file_offset": f"0x{file_offset:x}",
        "entries": TABLE_ENTRIES,
        "bitwise_mismatches": mismatches,
        "maximum_ulp_difference": max_ulp_difference,
        "formula": "float32(sin(pi * index / 8192)^2)",
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "dll",
        nargs="?",
        type=Path,
        default=ROOT / "binary" / "vt_pau.dll",
        help="read-only 2013 Paul DLL (default: binary/vt_pau.dll)",
    )
    arguments = parser.parse_args()
    print(json.dumps(compare(arguments.dll), indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
