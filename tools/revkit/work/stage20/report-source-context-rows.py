#!/usr/bin/env python3
"""List legacy Kate rows matching selected transformed context keys."""

from __future__ import annotations

import sys
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data

ENGINE = ROOT / "binary" / "vt_pau.dll"
SOURCE = ROOT / "data-kate" / "M16" / "mc_idx_tbl"
IMAGE_BASE = 0x10000000
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)
CONTEXTS = {
    "slot0": (bytes((0x5A, 0x22, 0x17)), {bytes((0x1E, 0x1E))}),
    "slot2": (
        bytes((0x22, 0x17, 0x2B)),
        {bytes((0x00, 0x1E)), bytes((0x00, 0x00))},
    ),
    "slot4": (bytes((0x17, 0x2B, 0x2F)), {bytes((0x00, 0x00))}),
    "slot5": (bytes((0x2B, 0x30, 0x5A)), {bytes((0x04, 0x04))}),
}


def main() -> None:
    binary = ENGINE.read_bytes()
    tables = tuple(
        binary[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    if any(len(table) != 256 for table in tables):
        raise ValueError("class-key lookup table extends beyond the local DLL")

    matches: dict[tuple[str, bytes], list[str]] = defaultdict(list)
    for path in sorted(SOURCE.glob("unit-*.idx")):
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        count = parsed.unit_count
        block_end = parsed.block_start + count * parsed.record_stride
        source_keys = raw[block_end + 2 * count : block_end + 7 * count]
        for row in range(count):
            key = source_keys[row * 5 : row * 5 + 5]
            prefix = bytes(
                (tables[0][key[0]], tables[1][key[1]], tables[2][key[2]])
            )
            for context, (wanted_prefix, wanted_tails) in CONTEXTS.items():
                if prefix == wanted_prefix and key[3:5] in wanted_tails:
                    matches[(context, key[3:5])].append(f"{path.stem.removeprefix('unit-')}:{row}")

    print("context\tlegacy_tail\trow_count\tbank_local_rows")
    for context, (_, wanted_tails) in CONTEXTS.items():
        for tail in sorted(wanted_tails):
            rows = matches.get((context, tail), [])
            print(f"{context}\t{tail.hex(' ')}\t{len(rows)}\t{', '.join(rows)}")


if __name__ == "__main__":
    main()
