#!/usr/bin/env python3
"""Inventory framed DAT controls and check decoded counts against UPM lengths."""

from __future__ import annotations

import argparse
import json
import struct
import sys
from collections import Counter
from pathlib import Path

from decode_dat import BitReader

INDEX_HEADER_SIZE = 45
UNIT_RECORD_SIZE = 19
FEATURE_COLUMNS_SIZE = 21
DEFAULT_FRAME_SAMPLES = 256


def inventory_bank(root: Path, bank: str) -> dict[str, object]:
    index_path = root / "mc_idx_tbl" / f"unit-{bank}.idx"
    dat_path = root / "dat" / f"merged-{bank}.dat"
    upm_path = root / "dat" / f"merged-{bank}.upm"
    index = index_path.read_bytes()
    if len(index) < INDEX_HEADER_SIZE or not index[1:24].startswith(
        b"ver.2013\0VoiceText-Eng\0"
    ):
        raise ValueError(f"unsupported unit index: {index_path}")
    unit_count = struct.unpack_from("<I", index, 39)[0]
    stride = struct.unpack_from("<H", index, 43)[0]
    if stride != UNIT_RECORD_SIZE:
        raise ValueError(
            f"{index_path}: expected {UNIT_RECORD_SIZE}-byte records, got {stride}"
        )
    expected_index_size = INDEX_HEADER_SIZE + unit_count * (stride + FEATURE_COLUMNS_SIZE)
    if len(index) != expected_index_size:
        raise ValueError(f"{index_path}: expected {expected_index_size} bytes, got {len(index)}")

    dat = dat_path.read_bytes()
    upm = upm_path.read_bytes()
    frame_modes: Counter[int] = Counter()
    control_modes: Counter[int] = Counter()
    mode8_units: set[int] = set()
    sample_count_mismatches: list[dict[str, int]] = []

    for unit in range(unit_count):
        start = INDEX_HEADER_SIZE + unit * stride
        record = index[start : start + stride]
        dat_offset = struct.unpack_from("<I", record, 0)[0]
        dat_length = struct.unpack_from("<H", record, 8)[0]
        upm_offset = struct.unpack_from("<I", record, 10)[0]
        left_count, right_count = record[14], record[15]
        upm_length = left_count + right_count - 1
        if (
            not upm_length
            or dat_offset + dat_length > len(dat)
            or upm_offset + upm_length > len(upm)
        ):
            raise ValueError(f"{bank} unit {unit}: payload span is out of bounds")

        reader = BitReader(dat[dat_offset : dat_offset + dat_length])
        frame_samples = DEFAULT_FRAME_SAMPLES
        residual_width = 0
        sample_count = 0
        saw_terminator = False
        while True:
            mode = reader.rice(2)
            if mode == 4:
                saw_terminator = True
                break
            if mode == 5:
                width = reader.rice(2)
                frame_samples = reader.rice(width)
                if not 1 <= frame_samples <= 1 << 16:
                    raise ValueError(f"{bank} unit {unit}: invalid frame size {frame_samples}")
                control_modes[mode] += 1
                continue
            if mode == 6:
                output_shift = reader.rice(2)
                if output_shift > 15:
                    raise ValueError(f"{bank} unit {unit}: invalid output shift {output_shift}")
                control_modes[mode] += 1
                continue
            if mode not in (0, 1, 2, 3, 8):
                raise ValueError(f"{bank} unit {unit}: unsupported control value {mode}")

            frame_modes[mode] += 1
            if mode == 8:
                mode8_units.add(unit)
            else:
                residual_width = reader.rice(3)
                for _ in range(frame_samples):
                    reader.rice(residual_width + 1)
            sample_count += frame_samples

        if not saw_terminator or sample_count == 0:
            raise ValueError(f"{bank} unit {unit}: missing terminator or decoded frames")
        expected = 2 * sum(upm[upm_offset : upm_offset + upm_length])
        if sample_count != expected:
            sample_count_mismatches.append(
                {"unit": unit, "stream_samples": sample_count, "upm_samples": expected}
            )
        if unit and unit % 100_000 == 0:
            print(f"{bank}: checked {unit}/{unit_count} units", file=sys.stderr, flush=True)

    if sample_count_mismatches:
        raise ValueError(
            f"{bank}: {len(sample_count_mismatches)} DAT/UPM count mismatches; "
            f"first: {sample_count_mismatches[0]}"
        )

    return {
        "bank": bank,
        "unit_count": unit_count,
        "dat_size_bytes": len(dat),
        "frame_modes": dict(sorted(frame_modes.items())),
        "control_modes": dict(sorted(control_modes.items())),
        "mode8_frame_count": frame_modes[8],
        "mode8_unit_count": len(mode8_units),
        "dat_upm_sample_count_mismatches": 0,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--root",
        type=Path,
        default=Path("data-paul/M16"),
        help="local 2013 Paul M16 data root (default: data-paul/M16)",
    )
    args = parser.parse_args()
    results = [inventory_bank(args.root, bank) for bank in ("gen", "num", "etc", "alp")]
    print(
        json.dumps(
            {"banks": results, "unit_count": sum(r["unit_count"] for r in results)},
            indent=2,
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
