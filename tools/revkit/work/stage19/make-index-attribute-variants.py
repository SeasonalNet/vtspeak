#!/usr/bin/env python3
"""Build controlled fills for the 2013-only unit-index byte column."""

from __future__ import annotations

import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = WORK / "index-adapter"
KATE_INDEXES = WORK / "data-kate-copy" / "M16" / "mc_idx_tbl"
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data  # noqa: E402
from inspect_unit_idx import inspect as inspect_versioned  # noqa: E402

VARIANTS = {
    "attr-a-copy": "copy the preserved one-byte source column into the inserted column",
    "signature-last-copy": "copy the last preserved signature byte into the inserted column",
    "signature-last-transfer": "move the last preserved signature byte into the inserted column and zero-pad its signature slot",
    "constant-8": "fill the inserted column with 8, the observed target value in the trace",
}


def main() -> None:
    for variant, description in VARIANTS.items():
        destination = WORK / f"index-adapter-{variant}"
        shutil.copytree(SOURCE, destination, dirs_exist_ok=True)
        rows = [
            f"variant={variant}",
            f"hypothesis={description}",
            "all source bytes, including the 7-byte signature and three feature groups, otherwise preserved",
        ]
        for original in sorted(KATE_INDEXES.glob("unit-*.idx")):
            parsed = parse_index_data(original, original.read_bytes())
            adapted = destination / original.name
            raw = bytearray(adapted.read_bytes())
            units = parsed.unit_count
            block_end = parsed.block_start + units * parsed.record_stride
            inserted_column = block_end + units * 8
            source_attr_a = block_end
            source_signature_last = block_end + units * 7
            if variant == "attr-a-copy":
                values = bytes(raw[source_attr_a : source_attr_a + units])
            elif variant == "signature-last-copy":
                values = bytes(raw[source_signature_last : source_signature_last + units])
            elif variant == "signature-last-transfer":
                values = bytes(raw[source_signature_last : source_signature_last + units])
                raw[source_signature_last : source_signature_last + units] = bytes(units)
            else:
                values = bytes([8]) * units
            if len(values) != units:
                raise ValueError(f"{original.name}: inserted column has wrong extent")
            raw[inserted_column : inserted_column + units] = values
            adapted.write_bytes(raw)
            checked = inspect_versioned(adapted)
            if checked["unit_count"] != units:
                raise ValueError(f"{adapted}: unit count changed")
            rows.append(
                f"{original.name}: units={units} attr_b_sha256={hashlib.sha256(values).hexdigest()}"
            )
        (destination / "variant-manifest.txt").write_text("\n".join(rows) + "\n")
        (WORK / f"compose-index-{variant}.yaml").write_text(
            "services:\n"
            "  runtime:\n"
            "    volumes:\n"
            f"      - ../stage19/index-adapter-{variant}:/work/data-kate/M16/mc_idx_tbl:ro\n"
        )
        print("\n".join(rows))


if __name__ == "__main__":
    main()
