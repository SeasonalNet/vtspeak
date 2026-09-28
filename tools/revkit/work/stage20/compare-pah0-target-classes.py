#!/usr/bin/env python3
"""Compare P AH0 runtime target keys with rows in the adapted Kate indexes."""

from __future__ import annotations

import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))

from inspect_unit_idx import inspect

ENGINE = ROOT / "binary/vt_pau.dll"
ADAPTERS = sorted((ROOT / "tools/revkit/work/stage20").glob("index-adapter-key-repacked-*"))
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)
IMAGE_BASE = 0x10000000
TARGETS = {
    "slot0": bytes((90, 90, 53, 7, 90, 160, 0)),
    "slot2": bytes((90, 53, 7, 90, 90, 101, 0)),
}
LEGACY_ROWS = {272822: ("unit-gen2.idx", 92827), 272823: ("unit-gen2.idx", 92828),
               272824: ("unit-gen2.idx", 92829), 272825: ("unit-gen2.idx", 92830)}


def main() -> None:
    binary = ENGINE.read_bytes()
    tables = tuple(binary[a - IMAGE_BASE : a - IMAGE_BASE + 256] for a in TABLE_ADDRESSES)
    targets = {
        name: bytes((tables[0][sig[1]], tables[1][sig[2]], tables[2][sig[3]], sig[5], sig[6] & 0x20))
        for name, sig in TARGETS.items()
    }
    print("adapter\ttarget\truntime_signature\tquery_key\tclass_id\tpopulation")
    print("adapter\tlegacy_unit\tadapted_row_key\tclass_id")
    for adapter in ADAPTERS:
        populations: Counter[bytes] = Counter()
        for path in sorted(adapter.glob("unit-*.idx")):
            parsed = inspect(path)
            count = int(parsed["unit_count"])
            signature_start = int(parsed["opaque_block_start"]) + count * 20
            raw = path.read_bytes()
            for row in range(count):
                candidate = raw[signature_start + row * 7 : signature_start + (row + 1) * 7]
                candidate_key = bytes((tables[0][candidate[1]], tables[1][candidate[2]],
                                       tables[2][candidate[3]], candidate[5], candidate[6] & 0x20))
                populations[candidate_key] += 1
        class_ids = {candidate_key: index for index, candidate_key in enumerate(sorted(populations))}
        for name, signature in TARGETS.items():
            key = targets[name]
            print(f"{adapter.name}\t{name}\t{signature.hex(' ')}\t{key.hex(' ')}\t"
                  f"{class_ids.get(key, 'MISS')}\t{populations.get(key, 0)}")
        for unit, (file_name, row) in LEGACY_ROWS.items():
            path = adapter / file_name
            parsed = inspect(path)
            count = int(parsed["unit_count"])
            signature_start = int(parsed["opaque_block_start"]) + count * 20
            signature = path.read_bytes()[signature_start + row * 7 : signature_start + (row + 1) * 7]
            key = bytes((tables[0][signature[1]], tables[1][signature[2]], tables[2][signature[3]],
                         signature[5], signature[6] & 0x20))
            print(f"{adapter.name}\t{unit}\t{key.hex(' ')}\t{class_ids.get(key, 'MISS')}")


if __name__ == "__main__":
    main()
