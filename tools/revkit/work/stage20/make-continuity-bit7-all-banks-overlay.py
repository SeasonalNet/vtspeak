#!/usr/bin/env python3
"""Build an isolated full-bank Kate attr_48-to-signature-bit-7 overlay."""

from __future__ import annotations

import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
LEGACY = ROOT / "tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl"
BASE = WORK / "index-adapter-key-repacked-matched-context-attrb-transfer"
OUTPUT = WORK / "index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks"
BANKS = ("unit-alp", "unit-etc", "unit-gen", "unit-gen2", "unit-num")


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
    from inspect_legacy_unit_idx import parse_index_data
    from inspect_unit_idx import inspect as inspect_versioned

    if OUTPUT.exists():
        raise FileExistsError(f"refusing to overwrite existing overlay: {OUTPUT}")
    shutil.copytree(BASE, OUTPUT)

    report = [
        "variant=matched-context attr_40 transfer plus legacy attr_48 mapped to row-major signature byte 6 bit 7",
        f"base={BASE.relative_to(ROOT)}",
        f"legacy_source={LEGACY.relative_to(ROOT)}",
        "mapping=signature[row*7+6] |= 0x80 when legacy attr_48[row] is nonzero",
    ]
    total_units = 0
    total_marked = 0
    for bank in BANKS:
        name = f"{bank}.idx"
        source = LEGACY / name
        target = OUTPUT / name
        source_raw = source.read_bytes()
        target_raw = bytearray(target.read_bytes())
        legacy = parse_index_data(source, source_raw)
        converted = inspect_versioned(target)
        count = legacy.unit_count
        if converted["unit_count"] != count or legacy.record_stride != 19:
            raise ValueError(f"unit count or legacy stride mismatch for {name}")

        columns_start = legacy.block_start + count * legacy.record_stride
        columns = source_raw[columns_start:]
        if len(columns) != count * 20:
            raise ValueError(f"unexpected legacy trailing-column size for {name}")
        attr_48 = columns[count : 2 * count]
        if any(value not in (0, 1) for value in attr_48):
            raise ValueError(f"legacy attr_48 is not binary in {name}")

        signature_start = (
            int(converted["opaque_block_start"])
            + count * int(converted["opaque_block_stride_bytes"])
            + count
        )
        if len(target_raw) < signature_start + count * 7:
            raise ValueError(f"signature column is truncated in {name}")

        marked = 0
        for row, legacy_marker in enumerate(attr_48):
            if legacy_marker:
                byte_offset = signature_start + row * 7 + 6
                target_raw[byte_offset] |= 0x80
                marked += 1

        target.write_bytes(target_raw)
        checked = inspect_versioned(target)
        if checked["unit_count"] != count or not checked["size_matches"]:
            raise ValueError(f"patched versioned index failed validation: {name}")
        report.append(
            f"bank={bank} units={count} marked={marked} "
            f"base_sha256={sha256((BASE / name).read_bytes())} "
            f"overlay_sha256={sha256(target_raw)}"
        )
        total_units += count
        total_marked += marked

    report.extend((f"total_units={total_units}", f"total_marked={total_marked}"))
    (OUTPUT / "continuity-bit7-all-banks-manifest.txt").write_text(
        "\n".join(report) + "\n"
    )
    print("\n".join(report))


if __name__ == "__main__":
    main()
