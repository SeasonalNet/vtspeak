#!/usr/bin/env python3
"""Move the matched-context signature byte into the 2013 attr_b column."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
SOURCE = (
    ROOT
    / "tools/revkit/work/stage20"
    / "index-adapter-key-repacked-attr48-key3-attr40-key2-matched-context-tail-map"
)
DESTINATION = (
    ROOT
    / "tools/revkit/work/stage20"
    / "index-adapter-key-repacked-matched-context-attrb-transfer"
)
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
from inspect_unit_idx import inspect


def main() -> None:
    DESTINATION.mkdir(parents=True, exist_ok=True)
    for source in sorted(SOURCE.glob("unit-*.idx")):
        raw = bytearray(source.read_bytes())
        parsed = inspect(source)
        units = int(parsed["unit_count"])
        block_start = int(parsed["opaque_block_start"])
        tail_start = block_start + units * 19
        signature_start = tail_start + units
        attr_b_start = tail_start + 8 * units
        for unit in range(units):
            signature_byte = signature_start + 7 * unit + 4
            if raw[attr_b_start + unit] != 0:
                raise ValueError(f"{source.name}: expected zero-filled attr_b")
            raw[attr_b_start + unit] = raw[signature_byte]
            raw[signature_byte] = 0
        target = DESTINATION / source.name
        target.write_bytes(raw)
        updated = inspect(target)
        if updated["unit_count"] != units:
            raise ValueError(f"unit count changed while converting {source.name}")
    (DESTINATION / "manifest.txt").write_text(
        "source=matched-context 2013 repack\n"
        "signature_byte_4=moved_to_attr_b_then_zeroed\n"
        "other_bytes=preserved\n"
    )


if __name__ == "__main__":
    main()
