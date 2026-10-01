#!/usr/bin/env python3
"""Cross-check MakeInfo mode-2 Shift Size against Paul index sample spans."""

from __future__ import annotations

import struct
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
STAGE16 = ROOT / "tools/revkit/work/stage16"
INDEX_DIR = ROOT / "data-paul/M16/mc_idx_tbl"
BANKS = ("gen", "num", "etc", "alp")
INDEX_HEADER_BYTES = 45
INDEX_RECORD_BYTES = 19


def load_dat_offsets(bank: str) -> dict[int, tuple[int, int]]:
    path = INDEX_DIR / f"unit-{bank}.idx"
    raw = path.read_bytes()
    if len(raw) < INDEX_HEADER_BYTES:
        raise ValueError(f"{path}: truncated unit-index header")
    unit_count = struct.unpack_from("<I", raw, 39)[0]
    record_bytes = struct.unpack_from("<H", raw, 43)[0]
    if record_bytes != INDEX_RECORD_BYTES:
        raise ValueError(f"{path}: expected {INDEX_RECORD_BYTES}-byte records, got {record_bytes}")
    records_end = INDEX_HEADER_BYTES + unit_count * record_bytes
    if records_end > len(raw):
        raise ValueError(f"{path}: unit records extend past EOF")
    records = raw[INDEX_HEADER_BYTES:records_end]

    offsets: dict[int, tuple[int, int]] = {}
    for unit_index in range(len(records) // INDEX_RECORD_BYTES):
        start = unit_index * INDEX_RECORD_BYTES
        record = records[start : start + INDEX_RECORD_BYTES]
        dat_offset = struct.unpack_from("<I", record, 0)[0]
        first_side_samples = struct.unpack_from("<H", record, 4)[0]
        if dat_offset in offsets:
            raise ValueError(f"{path}: duplicate DAT offset {dat_offset}")
        offsets[dat_offset] = (unit_index, first_side_samples)
    return offsets


def parse_fields(block: str) -> dict[str, str]:
    fields: dict[str, str] = {}
    for line in block.splitlines():
        if " : " in line:
            key, value = line.split(" : ", 1)
            fields[key] = value
    return fields


def main() -> int:
    indexes = [load_dat_offsets(bank) for bank in BANKS]
    checked = 0
    matched = 0
    failures: list[str] = []
    captures: set[str] = set()
    unique_rows: set[tuple[int, int, int, int, int]] = set()

    for path in sorted(STAGE16.glob("*.asc.dtt")):
        for block in path.read_text(encoding="ascii").split("\n\n"):
            fields = parse_fields(block)
            if fields.get("TypeFlag") != "2" or fields.get("Phone Mode") != "2":
                continue

            checked += 1
            captures.add(path.name)
            file_index = int(fields["File Index"])
            if not 0 <= file_index < len(BANKS):
                failures.append(f"{path.name}: invalid file index {file_index}")
                continue

            dat_offset = int(fields["PCM Pos"])
            source = indexes[file_index].get(dat_offset)
            if source is None:
                failures.append(
                    f"{path.name}: no {BANKS[file_index]} unit at DAT offset {dat_offset}"
                )
                continue

            unit_index, first_side_samples = source
            pitch_first = int(fields["Pitch_first"])
            shift_size = int(fields["Shift Size"])
            unique_rows.add(
                (file_index, unit_index, first_side_samples, pitch_first, shift_size)
            )
            expected = first_side_samples - pitch_first
            if shift_size == expected:
                matched += 1
            else:
                failures.append(
                    f"{path.name}: unit {unit_index} Shift Size={shift_size}, "
                    f"first-side samples {first_side_samples} - Pitch_first "
                    f"{pitch_first} = {expected}"
                )

    print(f"mode2 rows checked: {checked}")
    print(f"captures checked: {len(captures)}")
    print(f"unique bank/unit/value tuples: {len(unique_rows)}")
    print(f"first-side samples minus Pitch_first matches: {matched}")
    for failure in failures:
        print(f"mismatch: {failure}", file=sys.stderr)
    return 1 if failures or not checked else 0


if __name__ == "__main__":
    raise SystemExit(main())
