#!/usr/bin/env python3
"""Map repeated-Hello 2013 query hits to source-index unit rows."""

from __future__ import annotations

import sys
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data

SOURCE = (
    ROOT
    / "tools"
    / "revkit"
    / "work"
    / "stage19"
    / "data-kate-copy"
    / "M16"
    / "mc_idx_tbl"
)
ENGINE = ROOT / "binary" / "vt_pau.dll"
IMAGE_BASE = 0x10000000
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)
BANK_BASES = {
    "gen": 0,
    "gen2": 179_995,
    "num": 278_128,
    "etc": 279_125,
    "alp": 282_791,
}

# Successful exact lookups in
# corpus-parity/stage20/hello-repeat-2013-producer/adapted-producer-keyclasses-gdb.log.
QUERY_HITS = {
    bytes.fromhex("172b2f0000"): (7, 21),
    bytes.fromhex("2b30220000"): (10,),
    bytes.fromhex("0d22170000"): (16,),
    bytes.fromhex("22172b0000"): (19,),
}
EXPECTED_MEMBER_COUNTS = {
    bytes.fromhex("172b2f0000"): 25,
    bytes.fromhex("2b30220000"): 1,
    bytes.fromhex("0d22170000"): 2,
    bytes.fromhex("22172b0000"): 9,
}
OLD_SELECTED = (
    273_369,
    273_370,
    273_371,
    232_670,
    266_023,
    264_071,
    264_072,
    264_073,
    264_074,
)


def lookup_tables() -> tuple[bytes, bytes, bytes]:
    image = ENGINE.read_bytes()
    first, second, third = (
        image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    if any(len(table) != 256 for table in (first, second, third)):
        raise ValueError("class-key lookup table extends beyond the local DLL")
    return first, second, third


def source_keys() -> dict[int, tuple[str, int, bytes]]:
    rows: dict[int, tuple[str, int, bytes]] = {}
    for bank, base in BANK_BASES.items():
        path = SOURCE / f"unit-{bank}.idx"
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        count = parsed.unit_count
        end = parsed.block_start + count * parsed.record_stride
        attr48 = raw[end + count : end + 2 * count]
        key_column = raw[end + 2 * count : end + 7 * count]
        attr40 = raw[end + 7 * count : end + 8 * count]
        for local in range(count):
            key = key_column[local * 5 : local * 5 + 5]
            signature = (
                bytes((attr48[local],))
                + key[:3]
                + bytes((attr40[local],))
                + key[3:]
            )
            rows[base + local] = (bank, local, signature)
    return rows


def main() -> None:
    tables = lookup_tables()
    rows = source_keys()
    members: dict[bytes, list[tuple[int, str, int]]] = defaultdict(list)
    for unit_id, (bank, local, signature) in rows.items():
        key = bytes(
            (
                tables[0][signature[1]],
                tables[1][signature[2]],
                tables[2][signature[3]],
                signature[5],
                signature[6] & 0x20,
            )
        )
        if key in QUERY_HITS:
            members[key].append((unit_id, bank, local))

    old_selected = set(OLD_SELECTED)
    print("query_calls\t2013_key\tsource_rows\told_selected_rows")
    for key, calls in QUERY_HITS.items():
        rows_for_key = members[key]
        expected_count = EXPECTED_MEMBER_COUNTS[key]
        if len(rows_for_key) != expected_count:
            raise ValueError(
                f"{key.hex()}: expected {expected_count} source rows, "
                f"found {len(rows_for_key)}"
            )
        selected = [unit_id for unit_id, _, _ in rows_for_key if unit_id in old_selected]
        print(
            f"{','.join(map(str, calls))}\t{key.hex()}\t{len(rows_for_key)}\t"
            f"{','.join(map(str, selected)) or '-'}"
        )
        print(
            "  members: "
            + " ".join(
                f"{unit_id}({bank}:{local})"
                for unit_id, bank, local in rows_for_key
            )
        )


if __name__ == "__main__":
    main()
