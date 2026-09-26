#!/usr/bin/env python3
"""Adapt legacy unit indexes for the current DLL's versioned index reader."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = WORK / "data-kate-copy" / "M16" / "mc_idx_tbl"
OUTPUT = WORK / "index-adapter"
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data  # noqa: E402
from inspect_unit_idx import inspect as inspect_versioned  # noqa: E402


def main() -> None:
    OUTPUT.mkdir(parents=True, exist_ok=True)
    manifest = [
        "source=tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl",
        "adapter=ver.2005 header to ver.2013; insert one zero-filled column",
        "insertion=after the legacy 1-byte and 7-byte column groups",
        "layout=column-major; check all three 10-bit cepdist row codes per unit",
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
        unit_count = parsed.unit_count
        if len(old_columns) != unit_count * 20:
            raise ValueError(f"{source.name}: unexpected legacy column extent")

        # The legacy 20-byte tail is column-major: a 1-byte column, a 7-byte
        # signature, then three 4-byte feature groups. The 2013 reader adds a
        # one-byte column between the signature and feature groups.
        split = unit_count * 8
        new_columns = old_columns[:split] + bytes(unit_count) + old_columns[split:]

        # The current DLL masks each feature-group word to a 10-bit row in
        # cepdist.tbl, whose declared dimension is 1,024. Reject an adapted
        # layout that violates that runtime invariant.
        cursor = unit_count * 9
        for group in range(3):
            for unit in range(unit_count):
                metric = int.from_bytes(
                    new_columns[cursor + unit * 2 : cursor + unit * 2 + 2], "little"
                ) & 0x3FFF
                if metric >= 1024:
                    raise ValueError(
                        f"{source.name}: adapted group {group + 1} unit {unit} "
                        f"has out-of-table metric code {metric}"
                    )
            cursor += unit_count * 4

        converted = raw[:1] + converted_header + raw[header_end:block_end] + new_columns
        destination = OUTPUT / source.name
        destination.write_bytes(converted)
        checked = inspect_versioned(destination)
        if checked["unit_count"] != parsed.unit_count:
            raise ValueError(f"{source.name}: unit count changed during conversion")
        if checked["opaque_block_stride_bytes"] != parsed.record_stride:
            raise ValueError(f"{source.name}: fixed-block stride changed during conversion")
        manifest.append(
            f"{source.name}: units={parsed.unit_count} bytes={len(raw)}->{len(converted)} "
            "current_reader_structure=valid all_metric_codes_in_0_1023"
        )

    (WORK / "index-adapter-manifest.txt").write_text("\n".join(manifest) + "\n")
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
