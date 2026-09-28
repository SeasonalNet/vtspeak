#!/usr/bin/env python3
"""Crosswalk old-engine Apple class rows against the matched 2013 repack."""

from __future__ import annotations

import hashlib
import sys
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data

OLD_DLL = (
    ROOT
    / "tools/revkit/work/corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll"
)
NEW_DLL = ROOT / "binary/vt_pau.dll"
OLD_INDEX = ROOT / "data-kate/M16/mc_idx_tbl"
STAGE19_INDEX = ROOT / "tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl"
NEW_INDEX_ROOT = ROOT / "tools/revkit/work/stage20"
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
NEW_TABLES = (0x1007B7EC, 0x1007B788, 0x1007B850)

# Feature records captured at the old engine's class-lookup return for Apple.
OLD_FEATURES = {
    "slot0": {bytes((90, 5, 53, 30, 30, 71, 0))},
    "slot2": {
        bytes((1, 71, 7, 0, 0, 7, 5)),
        bytes((1, 53, 7, 0, 0, 7, 2)),
        bytes((1, 53, 7, 0, 0, 7, 7)),
    },
    "slot3": {bytes((71, 7, 43, 0, 4, 43, 71))},
    "slot4": {bytes((7, 43, 90, 4, 4, 0, 7))},
}
OLD_APPEND_COUNTS = {"slot2": 48, "slot3": 20, "slot4": 157}
NEW_TARGETS = {
    "slot0": bytes((90, 90, 5, 71, 7, 160, 0)),
    "slot2": bytes((90, 5, 71, 7, 43, 0, 0)),
    "slot3": bytes((5, 71, 7, 43, 90, 5, 0)),
    "slot4": bytes((71, 7, 43, 90, 90, 69, 0)),
}


def tables(path: Path, addresses: tuple[int, ...] | list[int]) -> list[bytes]:
    image = path.read_bytes()
    result = [image[address - IMAGE_BASE : address - IMAGE_BASE + 256] for address in addresses]
    if any(len(table) != 256 for table in result):
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


def new_key(signature: bytes, lookup: list[bytes]) -> bytes:
    return bytes(
        (
            lookup[0][signature[1]],
            lookup[1][signature[2]],
            lookup[2][signature[3]],
            signature[5],
            signature[6] & 0x20,
        )
    )


def main() -> None:
    old_binary = OLD_DLL.read_bytes()
    old_lookup = {
        name: old_binary[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for name, address in OLD_TABLES.items()
    }
    if any(len(table) != 256 for table in old_lookup.values()):
        raise ValueError("old lookup table extends beyond the local DLL")
    new_lookup = tables(NEW_DLL, NEW_TABLES)

    old_members: dict[str, set[str]] = defaultdict(set)
    for path in sorted(OLD_INDEX.glob("unit-*.idx")):
        source = path.read_bytes()
        staged = (STAGE19_INDEX / path.name).read_bytes()
        if hashlib.sha256(source).digest() != hashlib.sha256(staged).digest():
            raise ValueError(f"{path.name}: Stage 19 source copy differs from the 2005 index")

        parsed = parse_index_data(path, source)
        count = parsed.unit_count
        records_end = parsed.block_start + count * parsed.record_stride
        keys = source[
            records_end + 2 * count : records_end + 7 * count
        ]
        for row in range(count):
            row_id = f"{path.stem.removeprefix('unit-')}:{row}"
            key = keys[row * 5 : row * 5 + 5]
            transformed = old_feature(key, old_lookup)
            for slot, features in OLD_FEATURES.items():
                if transformed in features:
                    old_members[slot].add(row_id)

    print("adapter\tslot\told_rows\t2013_rows\tshared\told_only\tnew_only")
    adapters = sorted(NEW_INDEX_ROOT.glob("index-adapter-key-repacked-*"))
    if not adapters:
        raise ValueError("no generated index repacks found")
    for slot in OLD_FEATURES:
        old_rows = old_members[slot]
        if slot in OLD_APPEND_COUNTS and len(old_rows) != OLD_APPEND_COUNTS[slot]:
            raise ValueError(
                f"{slot}: recovered source rows {len(old_rows)} "
                f"!= old runtime append count {OLD_APPEND_COUNTS[slot]}"
            )
    for adapter in adapters:
        new_members: dict[str, set[str]] = defaultdict(set)
        for path in sorted(OLD_INDEX.glob("unit-*.idx")):
            raw = path.read_bytes()
            parsed = parse_index_data(path, raw)
            count = parsed.unit_count
            records_end = parsed.block_start + count * parsed.record_stride
            repacked = (adapter / path.name).read_bytes()
            signature_start = records_end + count
            for row in range(count):
                row_id = f"{path.stem.removeprefix('unit-')}:{row}"
                signature = repacked[
                    signature_start + row * 7 : signature_start + (row + 1) * 7
                ]
                key_new = new_key(signature, new_lookup)
                for slot, target in NEW_TARGETS.items():
                    if key_new == new_key(target, new_lookup):
                        new_members[slot].add(row_id)
        for slot, old_rows in old_members.items():
            new_rows = new_members[slot]
            print(
                f"{adapter.name}\t{slot}\t{len(old_rows)}\t{len(new_rows)}"
                f"\t{len(old_rows & new_rows)}\t{len(old_rows - new_rows)}"
                f"\t{len(new_rows - old_rows)}"
            )


if __name__ == "__main__":
    main()
