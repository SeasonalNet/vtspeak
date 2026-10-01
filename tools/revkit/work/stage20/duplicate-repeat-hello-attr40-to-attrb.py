#!/usr/bin/env python3
"""Duplicate legacy attr_40 into attr_b on the repeated-Hello bit-7 overlay."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
SOURCE = (
    ROOT
    / "tools/revkit/work/stage20"
    / "index-adapter-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity"
)
DESTINATION = (
    ROOT
    / "tools/revkit/work/stage20"
    / "index-adapter-key-repacked-repeat-hello-attrb40-bit7-continuity"
)
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
from inspect_unit_idx import inspect


def main() -> None:
    DESTINATION.mkdir(parents=True, exist_ok=True)
    for source in sorted(SOURCE.glob("unit-*.idx")):
        original = source.read_bytes()
        raw = bytearray(original)
        parsed = inspect(source)
        count = int(parsed["unit_count"])
        block_start = int(parsed["opaque_block_start"])
        feature_start = block_start + count * 19
        signature_start = feature_start + count
        attrb_start = feature_start + count * 8
        for unit in range(count):
            if raw[attrb_start + unit] != 0:
                raise ValueError(f"{source.name}: expected zero-filled attr_b")
            raw[attrb_start + unit] = raw[signature_start + 7 * unit + 4]
        if raw[:attrb_start] != original[:attrb_start]:
            raise ValueError(f"{source.name}: bytes before attr_b changed")
        if raw[attrb_start + count :] != original[attrb_start + count :]:
            raise ValueError(f"{source.name}: bytes after attr_b changed")
        (DESTINATION / source.name).write_bytes(raw)
    (DESTINATION / "manifest.txt").write_text(
        "source=repeated-Hello bit-7-continuity overlay\n"
        "attr_b=duplicate legacy attr_40\n"
        "signature=preserved byte-for-byte\n"
    )


if __name__ == "__main__":
    main()
