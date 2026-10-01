#!/usr/bin/env python3
"""Verify source-record identity for the shared Hi candidate IDs."""

from __future__ import annotations

import hashlib
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data
from inspect_unit_idx import inspect as inspect_versioned

SOURCE = ROOT / "data-kate" / "M16" / "mc_idx_tbl"
ADAPTED = (
    Path(__file__).resolve().parent
    / "index-adapter-key-repacked-attr48-key3-attr40-key2-hi-tail-map"
)
BANK_BASES = {
    "gen": 0,
    "gen2": 179_995,
    "num": 278_128,
    "etc": 279_125,
    "alp": 282_791,
}
UNIT_IDS = (65_331, 272_820, 272_821, 281_932)
RECORD_SIZE = 19


def record_at(directory: Path, unit_id: int, *, versioned: bool) -> tuple[str, bytes]:
    bank, local = next(
        (name, unit_id - base)
        for name, base in reversed(tuple(BANK_BASES.items()))
        if unit_id >= base
    )
    path = directory / f"unit-{bank}.idx"
    raw = path.read_bytes()
    if versioned:
        info = inspect_versioned(path)
        start = int(info["opaque_block_start"])
        stride = int(info["opaque_block_stride_bytes"])
        count = int(info["unit_count"])
    else:
        info = parse_index_data(path, raw)
        start = info.block_start
        stride = info.record_stride
        count = info.unit_count
    if stride != RECORD_SIZE or local >= count:
        raise ValueError(f"{path}: invalid record location for unit {unit_id}")
    offset = start + local * stride
    return f"{bank}:{local}", raw[offset : offset + stride]


def main() -> None:
    for unit_id in UNIT_IDS:
        source_location, source_record = record_at(SOURCE, unit_id, versioned=False)
        adapted_location, adapted_record = record_at(ADAPTED, unit_id, versioned=True)
        if source_record != adapted_record:
            raise SystemExit(f"unit {unit_id}: payload record differs")
        digest = hashlib.sha256(source_record).hexdigest()
        print(
            f"unit={unit_id} source={source_location} adapted={adapted_location} "
            f"record_sha256={digest} identical=true"
        )


if __name__ == "__main__":
    main()
