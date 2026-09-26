#!/usr/bin/env python3
"""Build a candidate adapter matching the legacy seven-byte key ordering."""

from __future__ import annotations

import sys
from pathlib import Path

sys.dont_write_bytecode = True

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = WORK / "data-kate-copy" / "M16" / "mc_idx_tbl"
OUTPUT = WORK / "index-adapter-reordered"
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data  # noqa: E402
from inspect_unit_idx import inspect as inspect_versioned  # noqa: E402


def main() -> None:
    OUTPUT.mkdir(parents=True, exist_ok=True)
    rows = [
        "candidate=zero new leading column; preserve legacy 7-byte key then 1-byte attribute",
        "mapping=[key7][attribute1][groups12] -> [zero1][key7][attribute1][groups12]",
    ]
    for source in sorted(SOURCE.glob("unit-*.idx")):
        raw = source.read_bytes()
        parsed = parse_index_data(source, raw)
        header_end = 1 + raw[0]
        header = raw[1:header_end]
        if not header.startswith(b"ver.2005\0VoiceText-Eng\0"):
            raise ValueError(f"{source.name}: unexpected legacy header {header!r}")
        converted_header = header.replace(b"ver.2005", b"ver.2013", 1)
        block_end = parsed.block_start + parsed.unit_count * parsed.record_stride
        old_columns = raw[block_end:]
        count = parsed.unit_count
        if len(old_columns) != count * 20:
            raise ValueError(f"{source.name}: unexpected legacy feature-column extent")

        key_end = count * 7
        attribute_end = count * 8
        new_columns = (
            bytes(count)
            + old_columns[:key_end]
            + old_columns[key_end:attribute_end]
            + old_columns[attribute_end:]
        )
        if len(new_columns) != count * 21:
            raise ValueError(f"{source.name}: candidate layout has wrong extent")

        cursor = count * 9
        for group in range(3):
            for unit in range(count):
                metric = int.from_bytes(
                    new_columns[cursor + unit * 2 : cursor + unit * 2 + 2], "little"
                ) & 0x3FFF
                if metric >= 1024:
                    raise ValueError(
                        f"{source.name}: group {group + 1}, unit {unit} has metric {metric}"
                    )
            cursor += count * 4

        destination = OUTPUT / source.name
        destination.write_bytes(
            raw[:1] + converted_header + raw[header_end:block_end] + new_columns
        )
        checked = inspect_versioned(destination)
        if checked["unit_count"] != count:
            raise ValueError(f"{source.name}: unit count changed during conversion")
        rows.append(f"{source.name}: units={count} current_reader_structure=valid")

    (WORK / "index-reordered-variant-manifest.txt").write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
