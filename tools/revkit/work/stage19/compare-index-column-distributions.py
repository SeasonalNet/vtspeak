#!/usr/bin/env python3
"""Compare Kate legacy byte-column distributions with Paul's attr-B column."""

from __future__ import annotations

import struct
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
PAUL = ROOT / "data-paul" / "M16" / "mc_idx_tbl"
KATE = ROOT / "data-kate" / "M16" / "mc_idx_tbl"
OUTPUT = Path(__file__).resolve().parent / "index-attribute-distribution.tsv"


def columns(path: Path, expected_width: int) -> tuple[int, list[bytes]]:
    raw = path.read_bytes()
    position = 1 + raw[0]
    bank_count = struct.unpack_from("<H", raw, position)[0]
    position += 2
    for _ in range(bank_count):
        name_length = struct.unpack_from("<H", raw, position)[0]
        position += 2 + name_length + 1
    unit_count = struct.unpack_from("<I", raw, position)[0]
    position += 4
    record_stride = struct.unpack_from("<H", raw, position)[0]
    position += 2 + unit_count * record_stride
    column_bytes = len(raw) - position
    if unit_count == 0 or column_bytes != unit_count * expected_width:
        raise ValueError(
            f"{path}: expected {expected_width} column bytes per unit, "
            f"found {column_bytes} bytes for {unit_count} units"
        )
    result = [
        raw[position + index * unit_count : position + (index + 1) * unit_count]
        for index in range(expected_width)
    ]
    return unit_count, result


def byte_emd(left: bytes, right: bytes) -> float:
    left_counts = Counter(left)
    right_counts = Counter(right)
    left_cdf = right_cdf = distance = 0.0
    for value in range(256):
        left_cdf += left_counts[value] / len(left)
        right_cdf += right_counts[value] / len(right)
        distance += abs(left_cdf - right_cdf)
    return distance


def main() -> None:
    rows = ["kate_index\tunits\tlegacy_column\tmin\tmax\tdistinct\tpaul_attr_b_emd"]
    for kate_path in sorted(KATE.glob("unit-*.idx")):
        bank = kate_path.stem.removeprefix("unit-")
        paul_bank = "gen" if bank == "gen2" else bank
        paul_path = PAUL / f"unit-{paul_bank}.idx"
        kate_units, kate_columns = columns(kate_path, 20)
        paul_units, paul_columns = columns(paul_path, 21)
        if len(paul_columns[8]) == 0:
            raise ValueError(f"{paul_path}: empty attr-B column")
        paul_attr_b = paul_columns[8]
        for index, values in enumerate(kate_columns):
            rows.append(
                f"{kate_path.name}\t{kate_units}\t{index:02d}\t{min(values)}\t"
                f"{max(values)}\t{len(set(values))}\t{byte_emd(values, paul_attr_b):.6f}"
            )
        if paul_units == 0:
            raise ValueError(f"{paul_path}: empty index")
    OUTPUT.write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
