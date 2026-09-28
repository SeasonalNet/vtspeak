#!/usr/bin/env python3
"""Crosswalk observed 2006 shortlist feature records to 2013 class members."""

from __future__ import annotations

import sys
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data

OLD_DLL = (
    ROOT
    / "tools"
    / "revkit"
    / "work"
    / "corpus-parity"
    / "kate-msi"
    / "Program Files"
    / "NeoSpeech"
    / "Kate16"
    / "lib"
    / "vt_eng.dll"
)
NEW_DLL = ROOT / "binary" / "vt_pau.dll"
INDEX_DIR = ROOT / "data-kate" / "M16" / "mc_idx_tbl"
IMAGE_BASE = 0x10000000
OLD_TABLES = {
    "category": 0x10070410,
    "phone_zero": 0x10070170,
    "phone_tens_two": 0x10070230,
    "phone_other": 0x100700B0,
    "phone_tail": 0x10070050,
    "context_remainder_two": 0x100702F0,
    "context_category_zero": 0x100701D0,
    "context_category_nonzero": 0x10070110,
}
NEW_TABLES = {
    "first": 0x1007B7EC,
    "second": 0x1007B788,
    "third": 0x1007B850,
}

# Runtime bytes are decimal values printed by GDB's x/7ub at
# *(model + 0x98) + candidate_id * 7.
OLD_CLASS_RECORDS = {
    "slot0": {
        52164: (45, bytes((90, 34, 4, 30, 30, 5, 0))),
        52165: (2, bytes((90, 34, 4, 30, 30, 23, 0))),
        52166: (24, bytes((90, 34, 4, 30, 30, 24, 0))),
        52167: (4, bytes((90, 34, 4, 30, 30, 30, 0))),
    },
    "slot2": {
        32401: (5, bytes((34, 23, 43, 0, 30, 43, 34))),
        32400: (4, bytes((34, 23, 43, 0, 0, 43, 34))),
    },
    "slot4": {25383: (5, bytes((23, 43, 47, 0, 0, 48, 23)))},
    "slot5": {37302: (6, bytes((43, 48, 90, 4, 4, 0, 43)))},
}
NEW_TARGETS = {
    "slot0": (bytes((0x5A, 0x22, 0x17)), {bytes((0x1E, 0x1E))}),
    "slot2": (
        bytes((0x22, 0x17, 0x2B)),
        {bytes((0x00, 0x1E)), bytes((0x00, 0x00))},
    ),
    "slot4": (bytes((0x17, 0x2B, 0x2F)), {bytes((0x00, 0x00))}),
    "slot5": (bytes((0x2B, 0x30, 0x5A)), {bytes((0x04, 0x04))}),
}


def tables(path: Path, addresses: dict[str, int]) -> dict[str, bytes]:
    image = path.read_bytes()
    result = {
        name: image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for name, address in addresses.items()
    }
    if any(len(table) != 256 for table in result.values()):
        raise ValueError(f"lookup table extends beyond {path}")
    return result


def old_feature(key: bytes, lookup: dict[str, bytes]) -> bytes:
    category_zero = lookup["category"][key[1]] == 0
    if category_zero:
        first = lookup["phone_zero"][key[0]]
    elif key[3] // 10 == 2:
        first = lookup["phone_tens_two"][key[0]]
    else:
        first = lookup["phone_other"][key[0]]

    if key[3] % 10 == 2:
        third = lookup["context_remainder_two"][key[2]]
    elif category_zero:
        third = lookup["context_category_zero"][key[2]]
    else:
        third = lookup["context_category_nonzero"][key[2]]

    return bytes(
        (
            first,
            lookup["phone_tail"][key[1]],
            third,
            key[3],
            key[4],
            lookup["phone_tail"][key[2]],
            lookup["phone_tail"][key[0]],
        )
    )


def source_rows() -> list[tuple[str, bytes]]:
    rows: list[tuple[str, bytes]] = []
    for path in sorted(INDEX_DIR.glob("unit-*.idx")):
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        records_end = parsed.block_start + parsed.unit_count * parsed.record_stride
        key_column = raw[
            records_end + 2 * parsed.unit_count : records_end + 7 * parsed.unit_count
        ]
        for row in range(parsed.unit_count):
            key = key_column[row * 5 : row * 5 + 5]
            rows.append((f"{path.stem.removeprefix('unit-')}:{row}", key))
    return rows


def main() -> None:
    old_lookup = tables(OLD_DLL, OLD_TABLES)
    new_lookup = tables(NEW_DLL, NEW_TABLES)
    rows = source_rows()
    old_members: dict[str, dict[int, set[str]]] = defaultdict(dict)
    new_members: dict[str, set[str]] = defaultdict(set)

    for slot, records in OLD_CLASS_RECORDS.items():
        for class_id, (_, expected_feature) in records.items():
            old_members[slot][class_id] = {
                row_id
                for row_id, key in rows
                if old_feature(key, old_lookup) == expected_feature
            }

    for row_id, key in rows:
        for slot, (target_prefix, target_tails) in NEW_TARGETS.items():
            prefix = bytes(
                (
                    new_lookup["first"][key[0]],
                    new_lookup["second"][key[1]],
                    new_lookup["third"][key[2]],
                )
            )
            if prefix == target_prefix and key[3:5] in target_tails:
                new_members[slot].add(row_id)

    print("slot\told_class_id\truntime_weight\told_rows")
    for slot, records in OLD_CLASS_RECORDS.items():
        for class_id, (weight, _) in records.items():
            count = len(old_members[slot][class_id])
            print(f"{slot}\t{class_id}\t{weight}\t{count}")
            if count != weight:
                raise ValueError(
                    f"{slot} class {class_id}: source rows {count} != runtime weight {weight}"
                )

    print("\nslot\told_union\tnew_class\tintersection\told_only\tnew_only")
    for slot, records in OLD_CLASS_RECORDS.items():
        old_union = set().union(*old_members[slot].values())
        new_set = new_members[slot]
        intersection = old_union & new_set
        print(
            f"{slot}\t{len(old_union)}\t{len(new_set)}\t{len(intersection)}"
            f"\t{len(old_union - new_set)}\t{len(new_set - old_union)}"
        )
        print(f"  shared: {', '.join(sorted(intersection))}")
        print(f"  old only: {', '.join(sorted(old_union - new_set))}")
        print(f"  new only: {', '.join(sorted(new_set - old_union))}")


if __name__ == "__main__":
    main()
