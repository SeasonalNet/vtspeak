#!/usr/bin/env python3
"""Crosswalk captured 2006 Hi candidate rows to 2013 repacked lookup keys."""

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

# IDs and feature subgroups come from the original 2006 runtime candidate
# trace: trace-legacy-hi-span-filter-bypass.gdb.
OLD_POOLS = {
    "slot0": {
        "feature_a": (
            272_820,
            22_336,
            65_331,
            75_811,
            135_441,
            164_116,
            205_125,
            219_724,
            255_985,
            277_527,
        ),
    },
    "slot1": {
        "feature_a": (272_821, 279_498, 279_499, 280_449),
        "feature_b": (4_087, 4_102, 84_094, 177_718, 243_540),
    },
}

HI_TAIL_MAP = {
    (bytes.fromhex("5a2201"), bytes.fromhex("1e1e")): bytes.fromhex("a000"),
    (bytes.fromhex("22115a"), bytes.fromhex("0422")): bytes.fromhex("6500"),
}
RELAXATION_TARGETS = tuple(
    bytes.fromhex(value) for value in ("5d00", "5500", "4d00")
)


def load_rows() -> dict[int, tuple[str, int, int, bytes]]:
    rows: dict[int, tuple[str, int, int, bytes]] = {}
    for bank, base in BANK_BASES.items():
        path = SOURCE / f"unit-{bank}.idx"
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        count = parsed.unit_count
        end = parsed.block_start + count * parsed.record_stride
        attr48 = raw[end + count : end + 2 * count]
        keys = raw[end + 2 * count : end + 7 * count]
        attr40 = raw[end + 7 * count : end + 8 * count]
        for local in range(count):
            rows[base + local] = (
                bank,
                local,
                attr48[local],
                bytes((attr40[local],)) + keys[local * 5 : local * 5 + 5],
            )
    return rows


def main() -> None:
    image = ENGINE.read_bytes()
    tables = tuple(
        image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    if any(len(table) != 256 for table in tables):
        raise ValueError("class-key lookup table extends beyond the local DLL")

    rows = load_rows()
    print(
        "pool\told_feature_group\tunit_id\tbank_row\tlegacy_key\t"
        "2013_key_hi_tail_map"
    )
    for pool, groups in OLD_POOLS.items():
        for group, unit_ids in groups.items():
            grouped_keys: dict[bytes, list[int]] = defaultdict(list)
            for unit_id in unit_ids:
                bank, local, attr48, fields = rows[unit_id]
                attr40, key = fields[:1], fields[1:]
                prefix = bytes(
                    (tables[0][key[0]], tables[1][key[1]], tables[2][key[2]])
                )
                suffix = HI_TAIL_MAP.get((prefix, key[3:5]), key[3:5])
                signature = bytes((attr48,)) + key[:3] + attr40 + suffix
                query_key = bytes(
                    (
                        tables[0][signature[1]],
                        tables[1][signature[2]],
                        tables[2][signature[3]],
                        signature[5],
                        signature[6] & 0x20,
                    )
                )
                grouped_keys[query_key].append(unit_id)
                print(
                    f"{pool}\t{group}\t{unit_id}\t{bank}:{local}\t"
                    f"{key.hex()}\t{query_key.hex()}"
                )

            expected = len(unit_ids)
            if sum(map(len, grouped_keys.values())) != expected:
                raise ValueError(f"{pool}/{group}: source row count changed")
            print(
                f"# {pool}/{group}: {expected} captured old rows map to "
                f"{len(grouped_keys)} distinct 2013 keys"
            )

    print(
        "\n# slot1 feature_b alternative target keys for the three observed relaxations"
    )
    for target in RELAXATION_TARGETS:
        print(f"22115a{target.hex()}")


if __name__ == "__main__":
    main()
