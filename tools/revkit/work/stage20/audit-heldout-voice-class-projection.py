#!/usr/bin/env python3
"""Measure projected 2006-feature to 2013-key collisions on 2005 voices."""

from __future__ import annotations

import importlib.util
import sys
from collections import Counter, defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
VOICES = {
    "Kate": ROOT / "data-kate" / "M16" / "mc_idx_tbl",
    "Bridget": ROOT / "data-bridget" / "M16" / "mc_idx_tbl",
}
LEGACY_RESEARCH = ROOT / "tools" / "revkit" / "work" / "stage20"
LEGACY_HELPER = LEGACY_RESEARCH / "compare-legacy-new-shortlist-members.py"
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data


def load_legacy_helper():
    spec = importlib.util.spec_from_file_location("legacy_crosswalk", LEGACY_HELPER)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load crosswalk helper: {LEGACY_HELPER}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def indexes(directory: Path) -> list[tuple[str, bytes, int]]:
    result = []
    for path in sorted(directory.glob("unit-*.idx")):
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        if parsed.version != "ver.2005":
            raise ValueError(f"{path}: expected ver.2005, found {parsed.version}")
        count = parsed.unit_count
        records_end = parsed.block_start + count * parsed.record_stride
        tail = raw[records_end:]
        if len(tail) != count * 20:
            raise ValueError(f"{path}: expected 20 feature bytes per row")
        result.append((path.stem.removeprefix("unit-"), tail, count))
    if not result:
        raise ValueError(f"no legacy unit indexes found under {directory}")
    return result


def main() -> None:
    helper = load_legacy_helper()
    old_tables = helper.tables(helper.OLD_DLL, helper.OLD_TABLES)
    new_tables = helper.tables(helper.NEW_DLL, helper.NEW_TABLES)
    print(
        "voice\tlayout\trows\tlegacy_features\t2013_keys\t"
        "keys_merging_features\trows_outside_largest_feature\tfeatures_split_across_keys"
        "\tfalse_merge_pairs\tfalse_split_pairs\ttotal_mismatched_pairs"
    )
    for voice, directory in VOICES.items():
        voice_indexes = indexes(directory)
        legacy_rows: list[tuple[str, bytes, bytes, bytes]] = []
        for bank, tail, count in voice_indexes:
            attr48 = tail[count : 2 * count]
            key_column = tail[2 * count : 7 * count]
            attr40 = tail[7 * count : 8 * count]
            for index in range(count):
                key = key_column[index * 5 : index * 5 + 5]
                old_class = helper.old_feature(key, old_tables)
                legacy_rows.append(
                    (f"{bank}:{index}", old_class, key, bytes((attr48[index], attr40[index])))
                )

        for layout in ("structural", "signature"):
            old_members: dict[bytes, set[str]] = defaultdict(set)
            new_members: dict[bytes, set[str]] = defaultdict(set)
            old_classes_by_new: dict[bytes, set[bytes]] = defaultdict(set)
            new_keys_by_old: dict[bytes, set[bytes]] = defaultdict(set)
            old_class_counts_by_new: dict[bytes, Counter[bytes]] = defaultdict(Counter)
            for row_id, old_class, key, attributes in legacy_rows:
                attr48, attr40 = attributes
                if layout == "structural":
                    signature = bytes((attr48, *key, attr40))
                else:
                    signature = bytes((attr48, *key[:3], attr40, *key[3:]))
                new_key = bytes(
                    (
                        new_tables["first"][signature[1]],
                        new_tables["second"][signature[2]],
                        new_tables["third"][signature[3]],
                        signature[5],
                        signature[6] & 0x20,
                    )
                )
                old_members[old_class].add(row_id)
                new_members[new_key].add(row_id)
                old_classes_by_new[new_key].add(old_class)
                new_keys_by_old[old_class].add(new_key)
                old_class_counts_by_new[new_key][old_class] += 1

            contaminated = [
                key for key, classes in old_classes_by_new.items() if len(classes) > 1
            ]
            split = [feature for feature, keys in new_keys_by_old.items() if len(keys) > 1]
            extra_rows = sum(
                sum(counts.values()) - max(counts.values())
                for key, counts in old_class_counts_by_new.items()
                if len(old_classes_by_new[key]) > 1
            )
            false_merge_pairs = sum(
                sum(counts.values()) * (sum(counts.values()) - 1) // 2
                - sum(count * (count - 1) // 2 for count in counts.values())
                for counts in old_class_counts_by_new.values()
            )
            false_split_pairs = 0
            old_class_key_counts: dict[bytes, Counter[bytes]] = defaultdict(Counter)
            for new_key, counts in old_class_counts_by_new.items():
                for old_class, count in counts.items():
                    old_class_key_counts[old_class][new_key] += count
            for counts in old_class_key_counts.values():
                total = sum(counts.values())
                false_split_pairs += total * (total - 1) // 2 - sum(
                    count * (count - 1) // 2 for count in counts.values()
                )
            print(
                f"{voice}\t{layout}\t{len(legacy_rows)}\t{len(old_members)}\t"
                f"{len(new_members)}\t{len(contaminated)} "
                f"({len(contaminated) / len(new_members):.1%})\t"
                f"{extra_rows} ({extra_rows / len(legacy_rows):.1%})\t"
                f"{len(split)} ({len(split) / len(old_members):.1%})\t"
                f"{false_merge_pairs}\t{false_split_pairs}\t"
                f"{false_merge_pairs + false_split_pairs}"
            )


if __name__ == "__main__":
    main()
