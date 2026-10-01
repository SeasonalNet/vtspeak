#!/usr/bin/env python3
"""Count 2013 class rows produced by the repeated-Hello slot-8 tail map."""

from __future__ import annotations

import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_unit_idx import inspect

INDEX_DIR = (
    ROOT
    / "tools"
    / "revkit"
    / "work"
    / "stage20"
    / "index-adapter-key-repacked-attr48-key3-attr40-key2-repeat-hello-slot8-0404-to-4500"
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
TARGET_KEY = bytes.fromhex("2b305a4500")


def main() -> None:
    image = ENGINE.read_bytes()
    tables = tuple(
        image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    if any(len(table) != 256 for table in tables):
        raise ValueError("class-key lookup table extends beyond the local DLL")

    members = []
    prefix_keys = Counter()
    for bank, base in BANK_BASES.items():
        path = INDEX_DIR / f"unit-{bank}.idx"
        metadata = inspect(path)
        count = int(metadata["unit_count"])
        block_start = int(metadata["opaque_block_start"])
        record_stride = int(metadata["opaque_block_stride_bytes"])
        if record_stride != 19:
            raise ValueError(f"{path.name}: expected 19-byte records")
        signature_start = block_start + count * record_stride + count
        raw = path.read_bytes()
        signatures = raw[signature_start : signature_start + count * 7]
        if len(signatures) != count * 7:
            raise ValueError(f"{path.name}: truncated signature column")
        for local in range(count):
            signature = signatures[local * 7 : local * 7 + 7]
            key = bytes(
                (
                    tables[0][signature[1]],
                    tables[1][signature[2]],
                    tables[2][signature[3]],
                    signature[5],
                    signature[6] & 0x20,
                )
            )
            if key[:3] == TARGET_KEY[:3]:
                prefix_keys[key] += 1
            if key == TARGET_KEY:
                members.append((base + local, bank, local))

    print(f"key={TARGET_KEY.hex()} source_rows={len(members)}")
    for key, count in sorted(prefix_keys.items()):
        print(f"prefix_key={key.hex()} source_rows={count}")
    for unit_id, bank, local in members:
        print(f"{unit_id}\t{bank}:{local}")


if __name__ == "__main__":
    main()
