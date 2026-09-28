#!/usr/bin/env python3
"""Swap only the metric fields of two same-key gen2 candidate rows."""

from __future__ import annotations

import sys
import argparse
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
SOURCE = (
    ROOT
    / "tools/revkit/work/stage20"
    / "index-adapter-key-repacked-matched-context-attrb-transfer"
)
WORK = ROOT / "tools/revkit/work/stage20"
ROWS = (92829, 93376)
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
from inspect_unit_idx import inspect


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "fields",
        choices=("all", "metric-codes", "features", "attr-b", "scoring-fields"),
        help="which metric fields to exchange between the two rows",
    )
    args = parser.parse_args()
    destination = WORK / f"index-adapter-matched-context-attrb-{args.fields}-swap"
    destination.mkdir(parents=True, exist_ok=True)
    for source in sorted(SOURCE.glob("unit-*.idx")):
        raw = bytearray(source.read_bytes())
        parsed = inspect(source)
        units = int(parsed["unit_count"])
        if source.name != "unit-gen2.idx":
            (destination / source.name).write_bytes(raw)
            continue
        if max(ROWS) >= units:
            raise ValueError(f"candidate row is outside {source.name}")
        block_start = int(parsed["opaque_block_start"])
        feature_start = block_start + units * 19
        metrics_start = feature_start + 9 * units
        if args.fields in ("attr-b", "scoring-fields"):
            attr_b_start = feature_start + 8 * units
            left = attr_b_start + ROWS[0]
            right = attr_b_start + ROWS[1]
            raw[left], raw[right] = raw[right], raw[left]
        for group in range(3):
            group_start = metrics_start + group * 4 * units
            arrays = []
            if args.fields in ("all", "metric-codes"):
                arrays.append((group_start, 2))
            if args.fields in ("all", "features", "scoring-fields"):
                arrays.extend(
                    ((group_start + 2 * units, 1), (group_start + 3 * units, 1))
                )
            for array_start, width in arrays:
                left = array_start + ROWS[0] * width
                right = array_start + ROWS[1] * width
                left_value = bytes(raw[left : left + width])
                right_value = bytes(raw[right : right + width])
                raw[left : left + width] = right_value
                raw[right : right + width] = left_value
        target = destination / source.name
        target.write_bytes(raw)
        if inspect(target)["unit_count"] != units:
            raise ValueError(f"unit count changed while converting {source.name}")
    (destination / "manifest.txt").write_text(
        "source=matched-context repack with attr_40 transferred to attr_b\n"
        "metric_swap=unit-gen2 rows 92829 and 93376\n"
        f"swapped_fields={args.fields}\n"
        "other_bytes=preserved\n"
    )


if __name__ == "__main__":
    main()
