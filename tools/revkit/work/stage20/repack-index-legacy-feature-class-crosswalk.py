#!/usr/bin/env python3
"""Repack candidate rows using measured 2006 feature classes as an intervention.

This is a diagnostic overlay, not a proposed production conversion. It maps
the measured Hello. 2006 feature classes to their observed 2013 target keys
and moves accidental exact-key collisions to a sentinel suffix.
"""

from __future__ import annotations

import importlib.util
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = ROOT / "tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl"
OUTPUT = WORK / "index-adapter-key-repacked-legacy-feature-class-crosswalk"
COMPOSE = WORK / "compose-key-repacked-legacy-feature-class-crosswalk.yaml"
CROSSWALK = WORK / "compare-legacy-new-shortlist-members.py"
TARGET_KEYS = {
    "slot0": bytes.fromhex("5a 22 17 a0 00"),
    "slot2": bytes.fromhex("22 17 2b 20 00"),
    "slot4": bytes.fromhex("17 2b 2f 00 00"),
    "slot5": bytes.fromhex("2b 30 5a 45 00"),
}
SENTINEL_TAIL = bytes((0xFF, 0x20))
IMAGE_BASE = 0x10000000


def load_crosswalk():
    spec = importlib.util.spec_from_file_location("legacy_new_crosswalk", CROSSWALK)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load crosswalk helper: {CROSSWALK}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def inverse_tables(tables: tuple[bytes, bytes, bytes]) -> tuple[dict[int, int], ...]:
    inverses = []
    for table in tables:
        values: dict[int, int] = {}
        for source, target in enumerate(table):
            values.setdefault(target, source)
        inverses.append(values)
    return tuple(inverses)


def main() -> None:
    crosswalk = load_crosswalk()
    from inspect_legacy_unit_idx import parse_index_data
    from inspect_unit_idx import inspect as inspect_versioned

    old_tables = crosswalk.tables(crosswalk.OLD_DLL, crosswalk.OLD_TABLES)
    new_tables = crosswalk.tables(crosswalk.NEW_DLL, crosswalk.NEW_TABLES)
    new_table_tuple = tuple(new_tables[name] for name in ("first", "second", "third"))
    inverses = inverse_tables(new_table_tuple)
    feature_to_target: dict[bytes, bytes] = {}
    for slot, records in crosswalk.OLD_CLASS_RECORDS.items():
        target = TARGET_KEYS[slot]
        for _, feature in records.values():
            prior = feature_to_target.setdefault(feature, target)
            if prior != target:
                raise ValueError(f"legacy feature maps to two 2013 classes: {feature.hex()}")
    target_keys = set(TARGET_KEYS.values())
    counts: Counter[str] = Counter()
    OUTPUT.mkdir(parents=True, exist_ok=True)

    manifest = [
        "source=Kate ver.2005 unit indexes",
        "method=measured 2006 seven-byte feature to observed 2013 five-byte class key",
        "target_keys=" + ",".join(f"{slot}:{key.hex()}" for slot, key in TARGET_KEYS.items()),
        f"collision_sentinel_tail={SENTINEL_TAIL.hex()}",
    ]
    for source in sorted(SOURCE.glob("unit-*.idx")):
        raw = source.read_bytes()
        parsed = parse_index_data(source, raw)
        units = parsed.unit_count
        block_end = parsed.block_start + units * parsed.record_stride
        tail = raw[block_end:]
        if len(tail) != units * 20:
            raise ValueError(f"{source.name}: expected 20*N legacy feature bytes")
        attr_4c = tail[:units]
        attr_48 = tail[units : 2 * units]
        key5 = tail[2 * units : 7 * units]
        attr_40 = tail[7 * units : 8 * units]
        metrics = tail[8 * units :]
        signatures = bytearray()
        for unit in range(units):
            source_key = key5[unit * 5 : (unit + 1) * 5]
            signature = bytearray(
                (attr_48[unit], *source_key[:3], attr_40[unit], *source_key[3:5])
            )
            current_key = bytes(
                (
                    new_tables["first"][signature[1]],
                    new_tables["second"][signature[2]],
                    new_tables["third"][signature[3]],
                    signature[5],
                    signature[6] & 0x20,
                )
            )
            legacy_feature = crosswalk.old_feature(source_key, old_tables)
            desired_key = feature_to_target.get(legacy_feature)
            if desired_key is None and current_key in target_keys:
                desired_key = current_key[:3] + SENTINEL_TAIL
                counts["redirected_collisions"] += 1
            elif desired_key is not None:
                counts["mapped_legacy_feature_rows"] += 1
            if desired_key is not None and desired_key != current_key:
                for index, value in enumerate(desired_key[:3], start=0):
                    if value not in inverses[index]:
                        raise ValueError(f"no 2013 table preimage for key byte {index}: {value:02x}")
                    signature[index + 1] = inverses[index][value]
                signature[5] = desired_key[3]
                signature[6] = desired_key[4]
            signatures.extend(signature)

        converted_tail = attr_4c + bytes(signatures) + bytes(units) + metrics
        header_end = 1 + raw[0]
        header = raw[1:header_end]
        if not header.startswith(b"ver.2005\0VoiceText-Eng\0"):
            raise ValueError(f"{source.name}: unexpected version header")
        converted = (
            raw[:1]
            + header.replace(b"ver.2005", b"ver.2013", 1)
            + raw[header_end:block_end]
            + converted_tail
        )
        target = OUTPUT / source.name
        target.write_bytes(converted)
        inspected = inspect_versioned(target)
        if inspected["unit_count"] != units:
            raise ValueError(f"{source.name}: unit count changed")
        manifest.append(f"{source.name}: units={units} bytes={len(converted)}")

    (OUTPUT / "manifest.txt").write_text("\n".join(manifest) + "\n")
    COMPOSE.write_text(
        "services:\n"
        "  runtime:\n"
        "    volumes:\n"
        "      - ../stage20/index-adapter-key-repacked-legacy-feature-class-crosswalk:/work/data-kate/M16/mc_idx_tbl:ro\n"
    )
    print("\n".join(manifest))
    print(f"mapped_legacy_feature_rows={counts['mapped_legacy_feature_rows']}")
    print(f"redirected_collisions={counts['redirected_collisions']}")


if __name__ == "__main__":
    main()
