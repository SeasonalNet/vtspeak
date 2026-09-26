#!/usr/bin/env python3
"""Write reduced legacy unit indexes for a load-time boundary probe."""

from __future__ import annotations

import struct
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = WORK / "data-kate-copy" / "M16" / "mc_idx_tbl"
OUTPUT = WORK / "index-smoke" / "mc_idx_tbl"
LIMIT = 50

sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))
from inspect_legacy_unit_idx import parse_index_data  # noqa: E402


def main() -> None:
    OUTPUT.mkdir(parents=True, exist_ok=True)
    rows = [f"max_units={LIMIT}"]
    for source in sorted(SOURCE.glob("unit-*.idx")):
        raw = source.read_bytes()
        index = parse_index_data(source, raw)
        count = min(index.unit_count, LIMIT)
        prefix_end = index.block_start
        count_offset = prefix_end - 6
        prefix = bytearray(raw[:prefix_end])
        struct.pack_into("<I", prefix, count_offset, count)
        record_end = prefix_end + count * index.record_stride
        column_start = prefix_end + index.unit_count * index.record_stride
        column_end = column_start + count * index.feature_column_bytes
        encoded = bytes(prefix) + raw[prefix_end:record_end] + raw[
            column_start:column_end
        ]
        destination = OUTPUT / source.name
        destination.write_bytes(encoded)
        # Reparse so count, layout, and exact reduced extent are checked.
        checked = parse_index_data(destination, encoded)
        if checked.unit_count != count:
            raise ValueError(f"{source.name}: reduced count did not round-trip")
        rows.append(
            f"{source.name} original_units={index.unit_count} probe_units={count} "
            f"bytes={len(raw)}->{len(encoded)}"
        )
    (WORK / "index-smoke-manifest.txt").write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
