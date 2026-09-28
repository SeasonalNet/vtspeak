#!/usr/bin/env python3
"""Check whether captured 2013 target keys occur in experimental index repacks."""

from __future__ import annotations

import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))

from inspect_legacy_unit_idx import parse_index_data
from inspect_unit_idx import inspect

ENGINE = ROOT / "binary/vt_pau.dll"
ADAPTERS = sorted((ROOT / "tools/revkit/work/stage20").glob("index-adapter-key-repacked-*"))
TARGETS = {
    "slot0": bytes((90, 90, 34, 23, 43, 160, 0)),
    "slot2": bytes((90, 34, 23, 43, 48, 32, 0)),
    "slot4": bytes((34, 23, 43, 48, 66, 0, 0)),
    "slot5": bytes((23, 43, 48, 90, 90, 69, 0)),
}
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)
IMAGE_BASE = 0x10000000


def class_key(signature: bytes, tables: tuple[bytes, bytes, bytes]) -> bytes:
    return bytes(
        (
            tables[0][signature[1]],
            tables[1][signature[2]],
            tables[2][signature[3]],
            signature[5],
            signature[6] & 0x20,
        )
    )


def main() -> None:
    binary = ENGINE.read_bytes()
    tables = tuple(binary[address - IMAGE_BASE : address - IMAGE_BASE + 256] for address in TABLE_ADDRESSES)
    if any(len(table) != 256 for table in tables):
        raise ValueError("class-key lookup table extends beyond the local DLL")

    target_keys = {name: class_key(signature, tables) for name, signature in TARGETS.items()}
    print("target\traw_signature\tderived_class_key")
    for name, signature in TARGETS.items():
        print(f"{name}\t{signature.hex(' ')}\t{target_keys[name].hex(' ')}")

    old_tails: Counter[bytes] = Counter()
    for path in sorted((ROOT / "data-kate/M16/mc_idx_tbl").glob("unit-*.idx")):
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        tail_start = parsed.block_start + parsed.unit_count * parsed.record_stride
        tail = raw[tail_start:]
        key_start = 2 * parsed.unit_count
        for index in range(parsed.unit_count):
            row_key = tail[key_start + index * 5 : key_start + (index + 1) * 5]
            old_tails[row_key[3:5]] += 1
    print("\nlegacy_key_tail\tunit_count")
    for tail, count in sorted(old_tails.items()):
        print(f"{tail.hex(' ')}\t{count}")

    print("\nadapter\ttarget\tclass_index\tunit_population")
    for adapter in ADAPTERS:
        populations: Counter[bytes] = Counter()
        prefix_tails: dict[bytes, Counter[bytes]] = {}
        prefix_legacy_tails: dict[bytes, Counter[bytes]] = {}
        for path in sorted(adapter.glob("unit-*.idx")):
            parsed = inspect(path)
            raw = path.read_bytes()
            count = int(parsed["unit_count"])
            tail = int(parsed["opaque_block_start"]) + count * 19
            signatures_start = tail + count
            source_path = ROOT / "data-kate/M16/mc_idx_tbl" / path.name
            source_raw = source_path.read_bytes()
            source = parse_index_data(source_path, source_raw)
            if source.unit_count != count:
                raise ValueError(f"{path.name}: source and repack counts differ")
            source_key_start = (
                source.block_start + count * source.record_stride + 2 * count
            )
            for index in range(count):
                signature = raw[
                    signatures_start + index * 7 : signatures_start + (index + 1) * 7
                ]
                key = class_key(signature, tables)
                populations[key] += 1
                prefix_tails.setdefault(key[:3], Counter())[key[3:]] += 1
                if adapter.name.endswith("attr48-key3-attr40-key2"):
                    source_key = source_raw[
                        source_key_start + index * 5 : source_key_start + (index + 1) * 5
                    ]
                    prefix_legacy_tails.setdefault(key[:3], Counter())[source_key[3:5]] += 1
        if not populations:
            continue
        sorted_keys = sorted(populations)
        class_indexes = {key: index for index, key in enumerate(sorted_keys)}
        for name, key in target_keys.items():
            if key in populations:
                index = class_indexes[key]
                population = populations[key]
                print(f"{adapter.name}\t{name}\t{index}\t{population}")
            else:
                print(f"{adapter.name}\t{name}\tMISS\t0")
        if adapter.name.endswith("attr48-key3-attr40-key2"):
            print("\ncurrent_repack_prefix\ttarget\tquery_tail\tmatching_units\ttop_inventory_tails")
            for name, key in target_keys.items():
                tails = prefix_tails.get(key[:3], Counter())
                legacy_tails = prefix_legacy_tails.get(key[:3], Counter())
                common = ", ".join(
                    f"{tail.hex()}:{count}" for tail, count in tails.most_common(8)
                )
                legacy_common = ", ".join(
                    f"{tail.hex()}:{count}" for tail, count in legacy_tails.most_common(8)
                )
                print(
                    f"{key[:3].hex()}\t{name}\t{key[3:].hex()}\t"
                    f"{sum(tails.values())}\t{common or 'none'}\t"
                    f"legacy_source_tails={legacy_common or 'none'}"
                )
            print("\ncurrent_repack_first2\ttarget\tmatching_units\tthird_byte_populations")
            for name, key in target_keys.items():
                neighbors = sorted(
                    (
                        (prefix[2], sum(tails.values()))
                        for prefix, tails in prefix_tails.items()
                        if prefix[:2] == key[:2]
                    ),
                    key=lambda pair: (-pair[1], pair[0]),
                )
                third_bytes = ", ".join(
                    f"{value:02x}:{count}" for value, count in neighbors[:12]
                )
                print(
                    f"{key[:2].hex()}\t{name}\t{sum(count for _, count in neighbors)}\t"
                    f"{third_bytes or 'none'}"
                )
            slot5_legacy_suffix = bytearray(TARGETS["slot5"])
            slot5_legacy_suffix[5] = 0x04
            alternate_key = class_key(bytes(slot5_legacy_suffix), tables)
            alternate_tails = prefix_tails.get(alternate_key[:3], Counter())
            alternate_source_tails = prefix_legacy_tails.get(alternate_key[:3], Counter())
            print(
                "\nslot5_suffix_counterfactual\t"
                f"signature={bytes(slot5_legacy_suffix).hex(' ')}\t"
                f"key={alternate_key.hex(' ')}\t"
                f"class_index={class_indexes.get(alternate_key, 'MISS')}\t"
                f"population={populations.get(alternate_key, 0)}\t"
                f"prefix_tails={dict((tail.hex(), count) for tail, count in alternate_tails.items())}\t"
                "legacy_source_tails="
                f"{dict((tail.hex(), count) for tail, count in alternate_source_tails.items())}"
            )


if __name__ == "__main__":
    main()
