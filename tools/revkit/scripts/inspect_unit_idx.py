#!/usr/bin/env python3
"""Read-only structural inspector for the VoiceText 2013 unit-index layout."""

from __future__ import annotations

import argparse
import json
import struct
from pathlib import Path
from typing import BinaryIO

UNIT_RECORD_STRIDE = 19


def read_exact(stream: BinaryIO, size: int, label: str) -> bytes:
    value = stream.read(size)
    if len(value) != size:
        raise ValueError(f"truncated {label}: expected {size} bytes, got {len(value)}")
    return value


def read_u16(stream: BinaryIO, label: str) -> int:
    return struct.unpack("<H", read_exact(stream, 2, label))[0]


def read_u32(stream: BinaryIO, label: str) -> int:
    return struct.unpack("<I", read_exact(stream, 4, label))[0]


def inspect_payload_records(
    records: bytes, unit_count: int, bank_name: str, data_dir: Path
) -> dict[str, object]:
    if not bank_name.startswith("merged-"):
        raise ValueError(f"unsupported bank name for payload lookup: {bank_name!r}")

    bank = bank_name.removeprefix("merged-")
    dat_path = data_dir / f"merged-{bank}.dat"
    upm_path = data_dir / f"merged-{bank}.upm"
    dat_size = dat_path.stat().st_size
    upm_size = upm_path.stat().st_size
    dat_spans: list[tuple[int, int]] = []
    upm_spans: list[tuple[int, int]] = []

    for index in range(unit_count):
        record = records[index * UNIT_RECORD_STRIDE : (index + 1) * UNIT_RECORD_STRIDE]
        dat_spans.append(
            (struct.unpack_from("<I", record, 0)[0], struct.unpack_from("<H", record, 8)[0])
        )
        upm_spans.append((struct.unpack_from("<I", record, 10)[0], record[14] + record[15] - 1))

    dat_contiguous_boundaries = sum(
        start + length == next_start
        for (start, length), (next_start, _) in zip(dat_spans, dat_spans[1:])
    )
    dat_in_bounds = all(start + length <= dat_size for start, length in dat_spans)
    dat_starts_at_zero = bool(dat_spans) and dat_spans[0][0] == 0
    dat_last_matches_file = bool(dat_spans) and dat_spans[-1][0] + dat_spans[-1][1] == dat_size
    upm_contiguous_boundaries = sum(
        start + length == next_start
        for (start, length), (next_start, _) in zip(upm_spans, upm_spans[1:])
    )
    upm_monotonic = all(
        start <= next_start
        for (start, _), (next_start, _) in zip(upm_spans, upm_spans[1:])
    )
    upm_in_bounds = all(length > 0 and start + length <= upm_size for start, length in upm_spans)
    upm_starts_at_zero = bool(upm_spans) and upm_spans[0][0] == 0
    upm_last_matches_file = bool(upm_spans) and upm_spans[-1][0] + upm_spans[-1][1] == upm_size
    upm_overlaps = sum(
        start + length > next_start
        for (start, length), (next_start, _) in zip(upm_spans, upm_spans[1:])
    )

    def sample_record(index: int) -> dict[str, int]:
        record = records[index * UNIT_RECORD_STRIDE : (index + 1) * UNIT_RECORD_STRIDE]
        return {
            "unit_index": index,
            "dat_offset": struct.unpack_from("<I", record, 0)[0],
            "u16_at_4": struct.unpack_from("<H", record, 4)[0],
            "u16_at_6": struct.unpack_from("<H", record, 6)[0],
            "dat_length": struct.unpack_from("<H", record, 8)[0],
            "upm_offset": struct.unpack_from("<I", record, 10)[0],
            "upm_left_count": record[14],
            "upm_right_count": record[15],
            "upm_span_length": record[14] + record[15] - 1,
            "opaque_bytes_16_to_18": list(record[16:19]),
        }

    samples = [sample_record(0)]
    if unit_count > 1:
        samples.append(sample_record(unit_count - 1))

    return {
        "record_stride_bytes": UNIT_RECORD_STRIDE,
        "record_fields": {
            "dat_offset": {"offset": 0, "width": 4, "encoding": "little-endian u32"},
            "dat_length": {"offset": 8, "width": 2, "encoding": "little-endian u16"},
            "upm_offset": {"offset": 10, "width": 4, "encoding": "little-endian u32"},
            "upm_left_count": {"offset": 14, "width": 1, "encoding": "u8"},
            "upm_right_count": {"offset": 15, "width": 1, "encoding": "u8"},
            "upm_span_length": "upm_left_count + upm_right_count - 1",
        },
        "dat_file": str(dat_path),
        "dat_file_size_bytes": dat_size,
        "dat_spans_in_bounds": dat_in_bounds,
        "dat_first_span_starts_at_zero": dat_starts_at_zero,
        "dat_exactly_contiguous_boundaries": dat_contiguous_boundaries,
        "dat_expected_contiguous_boundaries": max(0, unit_count - 1),
        "dat_last_span_ends_at_eof": dat_last_matches_file,
        "upm_file": str(upm_path),
        "upm_file_size_bytes": upm_size,
        "upm_offsets_monotonic": upm_monotonic,
        "upm_spans_in_bounds": upm_in_bounds,
        "upm_first_span_starts_at_zero": upm_starts_at_zero,
        "upm_exactly_contiguous_boundaries": upm_contiguous_boundaries,
        "upm_expected_contiguous_boundaries": max(0, unit_count - 1),
        "upm_last_span_ends_at_eof": upm_last_matches_file,
        "upm_overlapping_adjacent_spans": upm_overlaps,
        "sample_records": samples,
    }


def inspect(path: Path, data_dir: Path | None = None) -> dict[str, object]:
    raw = path.read_bytes()
    with path.open("rb") as stream:
        header_length = read_exact(stream, 1, "header length")[0]
        header = read_exact(stream, header_length, "version header")
        if not header.startswith(b"ver.2013\0VoiceText-Eng\0"):
            raise ValueError(f"{path}: unsupported version header {header!r}")

        table_count = read_u16(stream, "bank table count")
        banks: list[dict[str, int | str]] = []
        for index in range(table_count):
            name_length = read_u16(stream, f"bank {index} name length")
            name = read_exact(stream, name_length, f"bank {index} name").decode("ascii")
            tag = read_exact(stream, 1, f"bank {index} tag")[0]
            banks.append({"name": name, "tag": tag})

        unit_count = read_u32(stream, "unit count")
        opaque_stride = read_u16(stream, "opaque per-unit block stride")
        block_start = stream.tell()
        opaque_block_bytes = unit_count * opaque_stride
        opaque_records = read_exact(stream, opaque_block_bytes, "opaque per-unit blocks")

        # The parser bulk-reads these columns after skipping the opaque block.
        columns: list[tuple[str, int]] = [
            ("column_a", 1),
            ("unit_signature", 7),
            ("column_b", 1),
        ]
        for group in range(1, 4):
            columns.extend(
                (
                    (f"group_{group}_word", 2),
                    (f"group_{group}_byte_a", 1),
                    (f"group_{group}_byte_b", 1),
                )
            )

        column_width = sum(width for _, width in columns)
        read_exact(stream, unit_count * column_width, "unit columns")
        if stream.read(1):
            raise ValueError(f"{path}: trailing bytes after parsed unit columns")

    expected_size = block_start + unit_count * (opaque_stride + column_width)
    if expected_size != len(raw):
        raise ValueError(f"{path}: expected {expected_size} bytes, found {len(raw)}")

    result: dict[str, object] = {
        "file": str(path),
        "version": "ver.2013",
        "producer": "VoiceText-Eng",
        "banks": banks,
        "unit_count": unit_count,
        "opaque_block_stride_bytes": opaque_stride,
        "opaque_block_start": block_start,
        "column_layout": columns,
        "column_bytes_per_unit": column_width,
        "expected_size_bytes": expected_size,
        "actual_size_bytes": len(raw),
        "size_matches": True,
    }
    if data_dir is not None:
        if opaque_stride != UNIT_RECORD_STRIDE:
            raise ValueError(
                f"{path}: payload span inspection expects {UNIT_RECORD_STRIDE}-byte "
                f"records, found stride {opaque_stride}"
            )
        if len(banks) != 1:
            raise ValueError(f"{path}: payload span inspection expects one bank entry")
        result["payload_spans"] = inspect_payload_records(
            opaque_records, unit_count, str(banks[0]["name"]), data_dir
        )
    return result


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("files", nargs="+", type=Path, help="2013 unit-index files to inspect")
    parser.add_argument(
        "--data-dir",
        type=Path,
        help="optional directory containing matching merged-*.dat and merged-*.upm files",
    )
    args = parser.parse_args()
    results = [inspect(path, args.data_dir) for path in args.files]
    print(json.dumps(results, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
