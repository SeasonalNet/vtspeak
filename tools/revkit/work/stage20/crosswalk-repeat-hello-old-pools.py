#!/usr/bin/env python3
"""Crosswalk repeated-Hello 2006 shortlist features to 2013 exact-key hits."""

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
SOURCE = ROOT / "tools" / "revkit" / "work" / "stage19" / "data-kate-copy" / "M16" / "mc_idx_tbl"
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
NEW_TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)
BANK_BASES = {
    "gen": 0,
    "gen2": 179_995,
    "num": 278_128,
    "etc": 279_125,
    "alp": 282_791,
}

# Query bytes, shortlisted IDs, +0x8c weights, and seven-byte feature records
# come from trace-repeat-hello-old-query-pools.gdb on the original 2006 DLL.
OLD_POOLS = {
    "slot0": {
        "key": bytes.fromhex("5a22171e1e"),
        "records": {
            52164: (45, bytes((90, 34, 4, 30, 30, 5, 0))),
            52165: (2, bytes((90, 34, 4, 30, 30, 23, 0))),
            52166: (24, bytes((90, 34, 4, 30, 30, 24, 0))),
            52167: (4, bytes((90, 34, 4, 30, 30, 30, 0))),
        },
    },
    "slot1": {
        "key": bytes.fromhex("22172b001e"),
        "records": {
            32401: (5, bytes((34, 23, 43, 0, 30, 43, 34))),
            32400: (4, bytes((34, 23, 43, 0, 0, 43, 34))),
        },
    },
    "slot2": {
        "key": bytes.fromhex("172b300000"),
        "records": {25383: (5, bytes((23, 43, 47, 0, 0, 48, 23)))},
    },
    "slot3": {
        "key": bytes.fromhex("2b30220202"),
        "records": {37258: (1, bytes((43, 48, 34, 2, 2, 34, 43)))},
    },
    "slot5": {
        "key": bytes.fromhex("3022171414"),
        "records": {
            3408: (13, bytes((1, 34, 4, 10, 10, 5, 17))),
            3409: (14, bytes((1, 34, 4, 10, 10, 5, 30))),
            3410: (19, bytes((1, 34, 4, 10, 10, 5, 39))),
            3411: (1, bytes((1, 34, 4, 10, 10, 5, 47))),
            3412: (1, bytes((1, 34, 4, 10, 10, 5, 54))),
            3413: (6, bytes((1, 34, 4, 10, 10, 5, 62))),
            3414: (15, bytes((1, 34, 4, 10, 10, 5, 63))),
            3415: (1, bytes((1, 34, 4, 10, 10, 24, 62))),
            3433: (1, bytes((1, 34, 4, 20, 20, 23, 26))),
            3434: (2, bytes((1, 34, 4, 20, 20, 23, 30))),
        },
    },
    "slot6": {
        "key": bytes.fromhex("22172b0000"),
        "records": {32400: (4, bytes((34, 23, 43, 0, 0, 43, 34)))},
    },
    "slot7": {
        "key": bytes.fromhex("172b300000"),
        "records": {25383: (5, bytes((23, 43, 47, 0, 0, 48, 23)))},
    },
    "slot8": {
        "key": bytes.fromhex("2b305a0404"),
        "records": {37302: (6, bytes((43, 48, 90, 4, 4, 0, 43)))},
    },
}

# Successful exact lookups from the 2013 repeated-Hello producer capture.
NEW_HITS = {
    bytes.fromhex("172b2f0000"): (7, 21),
    bytes.fromhex("2b30220000"): (10,),
    bytes.fromhex("0d22170000"): (16,),
    bytes.fromhex("22172b0000"): (19,),
}
NEW_EXPECTED_COUNTS = {
    bytes.fromhex("172b2f0000"): 25,
    bytes.fromhex("2b30220000"): 1,
    bytes.fromhex("0d22170000"): 2,
    bytes.fromhex("22172b0000"): 9,
}


def load_tables(path: Path, addresses: dict[str, int] | tuple[int, ...]):
    image = path.read_bytes()
    if isinstance(addresses, dict):
        result = {
            name: image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
            for name, address in addresses.items()
        }
    else:
        result = tuple(
            image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
            for address in addresses
        )
    tables = result.values() if isinstance(result, dict) else result
    if any(len(table) != 256 for table in tables):
        raise ValueError(f"lookup table extends beyond {path}")
    return result


def old_feature(key: bytes, tables: dict[str, bytes]) -> bytes:
    category_zero = tables["category"][key[1]] == 0
    if category_zero:
        first = tables["phone_zero"][key[0]]
    elif key[3] // 10 == 2:
        first = tables["phone_tens_two"][key[0]]
    else:
        first = tables["phone_other"][key[0]]
    if key[3] % 10 == 2:
        third = tables["context_remainder_two"][key[2]]
    elif category_zero:
        third = tables["context_category_zero"][key[2]]
    else:
        third = tables["context_category_nonzero"][key[2]]
    return bytes(
        (
            first,
            tables["phone_tail"][key[1]],
            third,
            key[3],
            key[4],
            tables["phone_tail"][key[2]],
            tables["phone_tail"][key[0]],
        )
    )


def source_rows():
    rows = {}
    for bank, base in BANK_BASES.items():
        path = SOURCE / f"unit-{bank}.idx"
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        end = parsed.block_start + parsed.unit_count * parsed.record_stride
        attr48 = raw[end + parsed.unit_count : end + 2 * parsed.unit_count]
        keys = raw[end + 2 * parsed.unit_count : end + 7 * parsed.unit_count]
        attr40 = raw[end + 7 * parsed.unit_count : end + 8 * parsed.unit_count]
        for local in range(parsed.unit_count):
            key = keys[local * 5 : local * 5 + 5]
            signature = (
                bytes((attr48[local],))
                + key[:3]
                + bytes((attr40[local],))
                + key[3:]
            )
            rows[base + local] = (bank, local, key, signature)
    return rows


def main() -> None:
    old_tables = load_tables(OLD_DLL, OLD_TABLES)
    new_tables = load_tables(NEW_DLL, NEW_TABLE_ADDRESSES)
    rows = source_rows()
    by_feature: dict[bytes, set[int]] = defaultdict(set)
    by_new_key: dict[bytes, set[int]] = defaultdict(set)
    new_key_for_unit = {}
    for unit_id, (_, _, key, signature) in rows.items():
        by_feature[old_feature(key, old_tables)].add(unit_id)
        new_key = bytes(
            (
                new_tables[0][signature[1]],
                new_tables[1][signature[2]],
                new_tables[2][signature[3]],
                signature[5],
                signature[6] & 0x20,
            )
        )
        new_key_for_unit[unit_id] = new_key
        by_new_key[new_key].add(unit_id)

    print("old_slot\told_key\told_rows\t2013_hit_key\tshared\told_only\tnew_only")
    for slot, info in OLD_POOLS.items():
        members = set()
        weight_sum = 0
        for class_id, (weight, feature) in info["records"].items():
            matched = by_feature[feature]
            if len(matched) != weight:
                raise ValueError(
                    f"{slot} class {class_id}: source rows {len(matched)} != weight {weight}"
                )
            members |= matched
            weight_sum += weight
        if len(members) != weight_sum:
            raise ValueError(f"{slot}: overlapping class features in a shortlist")
        print(f"{slot}\t{info['key'].hex()}\t{len(members)}\t-\t-\t-\t-")
        print("  members:", " ".join(map(str, sorted(members))))
        key_counts: dict[bytes, int] = defaultdict(int)
        for unit_id in members:
            key_counts[new_key_for_unit[unit_id]] += 1
        print(
            "  source_keys:",
            " ".join(
                f"{key.hex()}={count}"
                for key, count in sorted(key_counts.items())
            ),
        )
        projected_candidates = set().union(
            *(by_new_key[key] for key in key_counts)
        )
        print(
            "  projected_key_buckets:",
            " ".join(
                f"{key.hex()}={len(by_new_key[key])}/{count}"
                for key, count in sorted(key_counts.items())
            ),
        )
        extras = projected_candidates - members
        print(f"  union_extra_members: {len(extras)}")
        for new_key, calls in NEW_HITS.items():
            new_members = by_new_key[new_key]
            if len(new_members) != NEW_EXPECTED_COUNTS[new_key]:
                raise ValueError(
                    f"2013 key {new_key.hex()}: expected "
                    f"{NEW_EXPECTED_COUNTS[new_key]} rows, found {len(new_members)}"
                )
            shared = members & new_members
            print(
                f"{slot}\t{info['key'].hex()}\t{len(members)}\t{new_key.hex()}"
                f"@{','.join(map(str, calls))}\t{len(shared)}\t"
                f"{len(members - new_members)}\t{len(new_members - members)}"
            )
            if shared:
                print("  shared:", " ".join(map(str, sorted(shared))))


if __name__ == "__main__":
    main()
