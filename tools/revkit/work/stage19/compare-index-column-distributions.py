#!/usr/bin/env python3
"""Compare Kate's legacy attr_40 field with Paul's attr-B distribution.

This is only a byte-distribution comparison. It does not establish semantic
equivalence between the window and the 2013 field.
"""

from __future__ import annotations

import struct
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
PAUL = ROOT / "data-paul" / "M16" / "mc_idx_tbl"
KATE = ROOT / "data-kate" / "M16" / "mc_idx_tbl"
OUTPUT = Path(__file__).resolve().parent / "index-attribute-distribution.tsv"


def index_tail(path: Path, expected_width: int) -> tuple[int, bytes]:
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
    return unit_count, raw[position:]


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
    rows = [
        "kate_index\tunits\tkate_field\tmin\tmax\tdistinct\t"
        "paul_attr_b_emd\tsemantic_mapping_established"
    ]
    for kate_path in sorted(KATE.glob("unit-*.idx")):
        bank = kate_path.stem.removeprefix("unit-")
        paul_bank = "gen" if bank == "gen2" else bank
        paul_path = PAUL / f"unit-{paul_bank}.idx"
        kate_units, kate_tail = index_tail(kate_path, 20)
        paul_units, paul_tail = index_tail(paul_path, 21)
        kate_attr_40 = kate_tail[7 * kate_units : 8 * kate_units]
        paul_attr_b = paul_tail[8 * paul_units : 9 * paul_units]
        if not paul_attr_b:
            raise ValueError(f"{paul_path}: empty attr-B column")
        rows.append(
            f"{kate_path.name}\t{kate_units}\tlegacy_attr_40\t{min(kate_attr_40)}\t"
            f"{max(kate_attr_40)}\t{len(set(kate_attr_40))}\t"
            f"{byte_emd(kate_attr_40, paul_attr_b):.6f}\tno"
        )
        if paul_units == 0:
            raise ValueError(f"{paul_path}: empty index")
    OUTPUT.write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
