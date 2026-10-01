#!/usr/bin/env python3
"""Explain the 2006/2013 member-set differences in selected Kate classes."""

from __future__ import annotations

import importlib.util
from collections import Counter
from pathlib import Path

SOURCE = Path(__file__).with_name("compare-legacy-new-shortlist-members.py")


def load_crosswalk():
    spec = importlib.util.spec_from_file_location("legacy_new_crosswalk", SOURCE)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load crosswalk helper: {SOURCE}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def main() -> None:
    crosswalk = load_crosswalk()
    old_tables = crosswalk.tables(crosswalk.OLD_DLL, crosswalk.OLD_TABLES)
    new_tables = crosswalk.tables(crosswalk.NEW_DLL, crosswalk.NEW_TABLES)
    rows = crosswalk.source_rows()
    key_by_row = dict(rows)

    for slot in ("slot0", "slot2", "slot4", "slot5"):
        records = crosswalk.OLD_CLASS_RECORDS[slot]
        old_members = {
            row_id
            for row_id, key in rows
            if any(
                crosswalk.old_feature(key, old_tables) == feature
                for _, feature in records.values()
            )
        }
        target_prefix, target_tails = crosswalk.NEW_TARGETS[slot]
        new_members = set()
        for row_id, key in rows:
            transformed_prefix = bytes(
                (
                    new_tables["first"][key[0]],
                    new_tables["second"][key[1]],
                    new_tables["third"][key[2]],
                )
            )
            if transformed_prefix == target_prefix and key[3:5] in target_tails:
                new_members.add(row_id)

        print(
            f"{slot}: old={len(old_members)} new={len(new_members)} "
            f"shared={len(old_members & new_members)} "
            f"old_only={len(old_members - new_members)} "
            f"new_only={len(new_members - old_members)}"
        )
        for label, members in (
            ("old_only", old_members - new_members),
            ("shared", old_members & new_members),
            ("new_only", new_members - old_members),
        ):
            if not members:
                continue
            raw_prefixes: Counter[str] = Counter()
            mapped_prefixes: Counter[str] = Counter()
            old_features: Counter[str] = Counter()
            for row_id in members:
                key = key_by_row[row_id]
                raw_prefixes[key[:3].hex()] += 1
                mapped = bytes(
                    (
                        new_tables["first"][key[0]],
                        new_tables["second"][key[1]],
                        new_tables["third"][key[2]],
                    )
                )
                mapped_prefixes[mapped.hex()] += 1
                old_features[crosswalk.old_feature(key, old_tables).hex()] += 1
            print(
                f"  {label}: raw_prefixes={dict(raw_prefixes.most_common())} "
                f"new_prefixes={dict(mapped_prefixes.most_common())} "
                f"old_features={dict(old_features.most_common())}"
            )


if __name__ == "__main__":
    main()
