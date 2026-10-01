#!/usr/bin/env python3
"""Map legacy attr_48 to 2013 signature byte-6 bit 7 for repeat-Hello units."""

from __future__ import annotations

import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
LEGACY = ROOT / "tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl"
BASE = WORK / "index-adapter-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools"
OUTPUT = WORK / "index-adapter-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity"
GLOBAL_UNIT_IDS = (264071, 264072, 264073)
UNIT_GEN2_GLOBAL_START = 179995


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    sys.path.insert(0, str(ROOT / "tools" / "revkit/scripts"))
    from inspect_legacy_unit_idx import parse_index_data
    from inspect_unit_idx import inspect as inspect_versioned

    if OUTPUT.exists():
        raise FileExistsError(f"refusing to overwrite existing overlay: {OUTPUT}")
    shutil.copytree(BASE, OUTPUT)

    source = LEGACY / "unit-gen2.idx"
    target = OUTPUT / "unit-gen2.idx"
    source_raw = source.read_bytes()
    target_raw = bytearray(target.read_bytes())
    legacy = parse_index_data(source, source_raw)
    converted = inspect_versioned(target)
    count = legacy.unit_count
    if converted["unit_count"] != count or legacy.record_stride != 19:
        raise ValueError("unit-gen2 row count or legacy stride mismatch")

    columns_start = legacy.block_start + count * legacy.record_stride
    columns = source_raw[columns_start:]
    if len(columns) != count * 20:
        raise ValueError("unexpected legacy trailing-column size")
    attr_48 = columns[count : 2 * count]
    signature_start = (
        int(converted["opaque_block_start"])
        + count * int(converted["opaque_block_stride_bytes"])
        + count
    )

    report = [
        "mapping=legacy attr_48[row] != 0 => versioned signature[row*7+6] |= 0x80",
        f"base={BASE.relative_to(ROOT)}",
        f"legacy_source={source.relative_to(ROOT)}",
        f"base_sha256={sha256(target.read_bytes())}",
    ]
    for unit_id in GLOBAL_UNIT_IDS:
        row = unit_id - UNIT_GEN2_GLOBAL_START
        if row < 0 or row >= count or attr_48[row] != 1:
            raise ValueError(f"expected legacy attr_48=1 for global unit {unit_id}")
        offset = signature_start + row * 7 + 6
        before = target_raw[offset]
        target_raw[offset] |= 0x80
        report.append(
            f"unit_id={unit_id} local_row={row} legacy_attr_48={attr_48[row]} "
            f"signature_byte6_before={before:#04x} after={target_raw[offset]:#04x}"
        )

    target.write_bytes(target_raw)
    checked = inspect_versioned(target)
    if checked["unit_count"] != count or not checked["size_matches"]:
        raise ValueError("patched unit-gen2 failed versioned-index validation")
    report.append(f"overlay_sha256={sha256(target_raw)}")
    (OUTPUT / "continuity-bit7-repeat-hello-manifest.txt").write_text(
        "\n".join(report) + "\n"
    )
    print("\n".join(report))


if __name__ == "__main__":
    main()
