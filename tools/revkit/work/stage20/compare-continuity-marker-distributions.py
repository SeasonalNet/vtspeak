#!/usr/bin/env python3
"""Compare Kate's legacy continuation byte with native 2013 signature flags."""

from __future__ import annotations

import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))

from inspect_legacy_unit_idx import parse_index_data
from inspect_unit_idx import inspect as inspect_versioned


def summarize_legacy(path: Path) -> tuple[int, int, str]:
    raw = path.read_bytes()
    index = parse_index_data(path, raw)
    count = index.unit_count
    columns_start = index.block_start + count * index.record_stride
    attr_48 = raw[columns_start + count : columns_start + 2 * count]
    marked = sum(value != 0 for value in attr_48)
    values = ",".join(f"{value:02x}:{frequency}" for value, frequency in sorted(Counter(attr_48).items()))
    return count, marked, values


def summarize_versioned(path: Path) -> tuple[int, int, str]:
    metadata = inspect_versioned(path)
    count = int(metadata["unit_count"])
    start = (
        int(metadata["opaque_block_start"])
        + count * int(metadata["opaque_block_stride_bytes"])
        + count
    )
    raw = path.read_bytes()
    signatures = raw[start : start + count * 7]
    byte_6 = bytes(signatures[row * 7 + 6] for row in range(count))
    marked = sum(bool(value & 0x80) for value in byte_6)
    values = ",".join(f"{value:02x}:{frequency}" for value, frequency in sorted(Counter(byte_6).items()))
    return count, marked, values


def main() -> None:
    sources = (
        ("Kate", "2005", ROOT / "data-kate/M16/mc_idx_tbl", summarize_legacy),
        ("Paul", "2013", ROOT / "data-paul/M16/mc_idx_tbl", summarize_versioned),
        ("James", "2013", ROOT / "data-james/M16/mc_idx_tbl", summarize_versioned),
    )
    print("voice\tgeneration\tbank\tunits\tmarked\tpercent\tbyte_values")
    for voice, generation, directory, summarize in sources:
        for path in sorted(directory.glob("unit-*.idx")):
            count, marked, values = summarize(path)
            print(
                f"{voice}\t{generation}\t{path.stem.removeprefix('unit-')}\t"
                f"{count}\t{marked}\t{100 * marked / count:.1f}\t{values}"
            )


if __name__ == "__main__":
    main()
