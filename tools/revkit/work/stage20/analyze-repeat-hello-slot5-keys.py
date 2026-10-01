#!/usr/bin/env python3
"""Show raw and effective 2013 keys for the repeated-Hello 2006 slot-5 pool."""

from __future__ import annotations

import runpy
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

research = runpy.run_path(
    str(Path(__file__).resolve().with_name("crosswalk-repeat-hello-old-pools.py"))
)
old_pools = research["OLD_POOLS"]
source_rows = research["source_rows"]()
old_tables = research["load_tables"](
    research["OLD_DLL"], research["OLD_TABLES"]
)
new_tables = research["load_tables"](
    research["NEW_DLL"], research["NEW_TABLE_ADDRESSES"]
)
feature_records = old_pools["slot5"]["records"]
features = {feature for _, feature in feature_records.values()}
groups: Counter[tuple[bytes, bytes, bytes]] = Counter()
target = bytes.fromhex("0d22170000")
target_rows: list[tuple[int, bytes, bytes]] = []

for unit_id, (_, _, key, signature) in source_rows.items():
    effective = bytes(
        (
            new_tables[0][signature[1]],
            new_tables[1][signature[2]],
            new_tables[2][signature[3]],
            signature[5],
            signature[6] & 0x20,
        )
    )
    if effective == target:
        target_rows.append((unit_id, key[:3], key[3:5]))
    if research["old_feature"](key, old_tables) not in features:
        continue
    groups[(key[:3], key[3:5], effective)] += 1

print("old_rows", sum(groups.values()))
print("raw_prefix:stored_tail -> effective_key count")
for (prefix, tail, effective), count in sorted(groups.items()):
    print(f"{prefix.hex()}:{tail.hex()} -> {effective.hex()} {count}")
print("live_call16_target_rows", len(target_rows))
for unit_id, prefix, tail in target_rows:
    print(f"{unit_id} raw={prefix.hex()}:{tail.hex()}")
