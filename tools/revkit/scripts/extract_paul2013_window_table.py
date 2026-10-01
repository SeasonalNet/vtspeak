#!/usr/bin/env python3
"""Extract the fixed 2013 blend coefficients from the read-only Paul DLL."""

from __future__ import annotations

import argparse
import math
import struct
from pathlib import Path


IMAGE_BASE = 0x10000000
TABLE_RVA = 0x6D1B0
TABLE_ENTRIES = 8192
TABLE_BYTES = TABLE_ENTRIES * 4


def section_file_offset(image: bytes, rva: int, size: int) -> int:
    if image[:2] != b"MZ":
        raise ValueError("input does not have a DOS MZ header")
    pe_offset = struct.unpack_from("<I", image, 0x3C)[0]
    if image[pe_offset : pe_offset + 4] != b"PE\0\0":
        raise ValueError("input does not have a PE signature")
    section_count = struct.unpack_from("<H", image, pe_offset + 6)[0]
    optional_size = struct.unpack_from("<H", image, pe_offset + 20)[0]
    section_offset = pe_offset + 24 + optional_size
    for index in range(section_count):
        offset = section_offset + index * 40
        virtual_size, virtual_address, raw_size, raw_offset = struct.unpack_from(
            "<IIII", image, offset + 8
        )
        section_size = max(virtual_size, raw_size)
        if virtual_address <= rva and rva + size <= virtual_address + section_size:
            file_offset = raw_offset + rva - virtual_address
            if file_offset + size > len(image) or rva - virtual_address + size > raw_size:
                raise ValueError("table RVA maps outside the section's raw data")
            return file_offset
    raise ValueError(f"RVA 0x{rva:x} is not contained in a PE section")


def extract(image_path: Path) -> bytes:
    image = image_path.read_bytes()
    pe_offset = struct.unpack_from("<I", image, 0x3C)[0]
    image_base = struct.unpack_from("<I", image, pe_offset + 24 + 28)[0]
    if image_base != IMAGE_BASE:
        raise ValueError(f"unexpected PE image base 0x{image_base:x}")
    file_offset = section_file_offset(image, TABLE_RVA, TABLE_BYTES)
    table = image[file_offset : file_offset + TABLE_BYTES]
    coefficients = struct.unpack(f"<{TABLE_ENTRIES}f", table)
    if any(not math.isfinite(value) or not 0.0 <= value <= 1.0 for value in coefficients):
        raise ValueError("blend table contains a non-finite or out-of-range coefficient")
    if coefficients[0] != 0.0 or coefficients[TABLE_ENTRIES // 2] != 1.0:
        raise ValueError("blend table endpoints do not match the observed half-cycle")
    midpoint = TABLE_ENTRIES // 2
    if any(left > right for left, right in zip(coefficients[:midpoint], coefficients[1:midpoint])):
        raise ValueError("rising half of the blend table is not monotonic")
    if any(left < right for left, right in zip(coefficients[midpoint:], coefficients[midpoint + 1 :])):
        raise ValueError("falling half of the blend table is not monotonic")
    return table


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("dll", type=Path, help="read-only 2013 Paul DLL")
    parser.add_argument("output", type=Path, help="output path for 8,192 little-endian float32 values")
    arguments = parser.parse_args()
    table = extract(arguments.dll)
    arguments.output.write_bytes(table)
    print(
        f"wrote {TABLE_ENTRIES} float32 coefficients from RVA 0x{TABLE_RVA:x} "
        f"to {arguments.output} ({len(table)} bytes)"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
