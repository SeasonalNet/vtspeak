#!/usr/bin/env python3
"""Extract the Paul DLL's model-context classifier tables for review."""

from __future__ import annotations

import argparse
import hashlib
import json
import struct
from pathlib import Path


IMAGE_BASE = 0x10000000
WEIGHT_TABLE_RVA = 0x7E388
STRING_TABLES = {
    "FUN_1000e540": (0x10078140, 0x1007824C),
    "FUN_1000e4e0": (0x10078038, 0x100780C0),
    "FUN_1000e4b0": (0x100780E0, 0x10078244),
    "FUN_1000e510": (0x100780FC, 0x10078248),
    "FUN_1000e430": (0x100780C8, 0x10078240),
    "FUN_1000e3a0": (0x10078044, 0x100780C4),
}
LITERAL_ADDRESSES = (
    0x10078F90,
    0x10077374,
    0x10077538,
    0x10079270,
    0x1007760C,
    0x10078AE8,
    0x10079014,
    0x1007926C,
    0x10077688,
    0x10079264,
    0x10078EB8,
    0x1007925C,
    0x10079254,
    0x1007924C,
    0x10077368,
    0x100776FC,
    0x1009F948,
    0x100773A0,
    0x10077764,
    0x10077768,
    0x1007776C,
    0x10077568,
)


def parse_pe(image: bytes) -> tuple[int, list[tuple[int, int, int, int]]]:
    if image[:2] != b"MZ":
        raise ValueError("input does not have a DOS MZ header")
    pe_offset = struct.unpack_from("<I", image, 0x3C)[0]
    if image[pe_offset : pe_offset + 4] != b"PE\0\0":
        raise ValueError("input does not have a PE signature")
    section_count = struct.unpack_from("<H", image, pe_offset + 6)[0]
    optional_size = struct.unpack_from("<H", image, pe_offset + 20)[0]
    image_base = struct.unpack_from("<I", image, pe_offset + 24 + 28)[0]
    section_offset = pe_offset + 24 + optional_size
    sections = []
    for index in range(section_count):
        offset = section_offset + index * 40
        virtual_size, virtual_address, raw_size, raw_offset = struct.unpack_from(
            "<IIII", image, offset + 8
        )
        sections.append((virtual_size, virtual_address, raw_size, raw_offset))
    return image_base, sections


def rva_offset(
    sections: list[tuple[int, int, int, int]], rva: int, size: int, image_size: int
) -> int:
    for virtual_size, virtual_address, raw_size, raw_offset in sections:
        if virtual_address <= rva and rva + size <= virtual_address + max(virtual_size, raw_size):
            offset = raw_offset + rva - virtual_address
            if offset + size > image_size or rva - virtual_address + size > raw_size:
                raise ValueError(f"RVA 0x{rva:x} is not backed by section data")
            return offset
    raise ValueError(f"RVA 0x{rva:x} is not contained in a PE section")


def va_bytes(
    image: bytes,
    sections: list[tuple[int, int, int, int]],
    image_base: int,
    address: int,
    size: int,
) -> bytes:
    if address < image_base:
        raise ValueError(f"VA 0x{address:x} is below image base")
    offset = rva_offset(sections, address - image_base, size, len(image))
    return image[offset : offset + size]


def c_string(image: bytes, sections: list[tuple[int, int, int, int]], image_base: int, address: int) -> str:
    for length in range(1, 4097):
        raw = va_bytes(image, sections, image_base, address, length)
        if raw[-1] == 0:
            try:
                return raw[:-1].decode("ascii")
            except UnicodeDecodeError as error:
                raise ValueError(f"string at VA 0x{address:x} is not ASCII") from error
    raise ValueError(f"string at VA 0x{address:x} has no NUL within 4096 bytes")


def extract(dll: Path) -> dict[str, object]:
    image = dll.read_bytes()
    image_base, sections = parse_pe(image)
    if image_base != IMAGE_BASE:
        raise ValueError(f"unexpected image base 0x{image_base:x}")

    raw_weights = va_bytes(image, sections, image_base, image_base + WEIGHT_TABLE_RVA, 512)
    weights = list(struct.unpack("<256h", raw_weights))
    result: dict[str, object] = {
        "source": dll.name,
        "sha256": hashlib.sha256(image).hexdigest(),
        "image_base": f"0x{image_base:08x}",
        "weight_table_rva": f"0x{WEIGHT_TABLE_RVA:x}",
        "weights_i16": weights,
        "sorted_tables": {},
        "literals": {},
    }
    sorted_tables: dict[str, object] = result["sorted_tables"]  # type: ignore[assignment]
    for function, (table_address, count_address) in STRING_TABLES.items():
        count = struct.unpack(
            "<i", va_bytes(image, sections, image_base, count_address, 4)
        )[0]
        if count < 0 or count > 4096:
            raise ValueError(f"{function} has unsupported string count {count}")
        pointers = struct.unpack(
            f"<{count}I",
            va_bytes(image, sections, image_base, table_address, count * 4),
        )
        values = [c_string(image, sections, image_base, address) for address in pointers]
        mapped = [tuple(weights[value] for value in value.encode("ascii")) + (weights[0],) for value in values]
        if mapped != sorted(mapped):
            raise ValueError(f"{function} table is not sorted by the recovered mapped comparator")
        sorted_tables[function] = {
            "pointer_table_va": f"0x{table_address:08x}",
            "count_va": f"0x{count_address:08x}",
            "count": count,
            "strings": values,
        }
    literals: dict[str, str] = result["literals"]  # type: ignore[assignment]
    for address in LITERAL_ADDRESSES:
        literals[f"0x{address:08x}"] = c_string(image, sections, image_base, address)
    return result


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("dll", type=Path, help="read-only Paul 2013 DLL")
    parser.add_argument("output", type=Path, help="JSON report destination")
    args = parser.parse_args()
    report = extract(args.dll)
    args.output.write_text(json.dumps(report, indent=2, ensure_ascii=True) + "\n")
    table_count = len(report["sorted_tables"])
    print(f"wrote {table_count} sorted string tables and classifier literals to {args.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
