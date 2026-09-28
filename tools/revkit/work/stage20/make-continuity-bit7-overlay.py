#!/usr/bin/env python3
"""Reproduce the first, invalid column-offset marker experiment.

This historical experiment writes at 6 * unit_count + row. That is not the
selector's row-major signature byte offset and must not be used as a valid
marker mapping. Use make-continuity-bit7-unit-major-overlay.py instead.
Vendor inputs and the source repack remain read-only.
"""

from __future__ import annotations

import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
LEGACY = ROOT / "tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl"
BASE = WORK / "index-adapter-key-repacked-matched-context-attrb-transfer"
OUTPUT = WORK / "index-adapter-key-repacked-matched-context-attrb-transfer-bit7"
ROWS = (92827, 92828, 92829)


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
    from inspect_legacy_unit_idx import parse_index_data
    from inspect_unit_idx import inspect as inspect_versioned

    if OUTPUT.exists():
        raise FileExistsError(f"refusing to overwrite existing overlay: {OUTPUT}")
    shutil.copytree(BASE, OUTPUT)

    source = LEGACY / "unit-gen2.idx"
    target = OUTPUT / source.name
    legacy_raw = source.read_bytes()
    target_raw = bytearray(target.read_bytes())
    legacy = parse_index_data(source, legacy_raw)
    converted = inspect_versioned(target)
    units = legacy.unit_count
    if converted["unit_count"] != units:
        raise ValueError("source and adapted unit counts differ")

    legacy_block_end = legacy.block_start + units * legacy.record_stride
    legacy_columns = legacy_raw[legacy_block_end:]
    if legacy.record_stride != 19 or len(legacy_columns) != units * 20:
        raise ValueError("unexpected legacy index record or column layout")
    attr_48 = legacy_columns[units : 2 * units]

    signature_offset = int(converted["opaque_block_start"]) + units * 19 + units
    rows: list[str] = []
    for row in ROWS:
        if row >= units or attr_48[row] == 0:
            raise ValueError(f"legacy row {row} does not carry the expected attr_48 marker")
        byte_offset = signature_offset + 6 * units + row
        before = target_raw[byte_offset]
        target_raw[byte_offset] = before | 0x80
        rows.append(
            f"row={row} legacy_attr_48={attr_48[row]:02x} "
            f"signature_byte6={before:02x}->{target_raw[byte_offset]:02x}"
        )

    target.write_bytes(target_raw)
    checked = inspect_versioned(target)
    if checked["unit_count"] != units or not checked["size_matches"]:
        raise ValueError("patched versioned index failed structural validation")

    manifest = [
        "variant=matched-context attr_40 transfer plus legacy attr_48 to signature byte 6 bit 7",
        f"base={BASE.relative_to(ROOT)}",
        f"legacy_source={source.relative_to(ROOT)}",
        "modified_global_units=272822,272823,272824",
        *rows,
        f"base_sha256={sha256(BASE / source.name)}",
        f"variant_sha256={sha256(target)}",
        f"validated_units={units}",
    ]
    (OUTPUT / "continuity-bit7-manifest.txt").write_text("\n".join(manifest) + "\n")
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
