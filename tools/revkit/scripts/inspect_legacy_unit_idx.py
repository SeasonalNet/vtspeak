#!/usr/bin/env python3
"""Inspect the observed VoiceText 2005/2009 M16 unit-index layouts."""

from __future__ import annotations

import argparse
import json
import struct
from dataclasses import dataclass
from pathlib import Path

from decode_dat import decode_payload

UNIT_RECORD_STRIDE = 19
FEATURE_COLUMN_BYTES = 20
SUPPORTED_VERSIONS = {"ver.2005", "ver.2009"}


@dataclass(frozen=True)
class Bank:
    name: str
    tag: int


@dataclass(frozen=True)
class LegacyIndex:
    path: Path
    version: str
    producer: str
    banks: tuple[Bank, ...]
    unit_count: int
    record_stride: int
    feature_column_bytes: int
    block_start: int
    records: bytes
    actual_size: int


def require_bytes(raw: bytes, offset: int, size: int, label: str) -> bytes:
    if size < 0 or offset < 0 or offset + size > len(raw):
        raise ValueError(
            f"truncated {label} at byte {offset}: need {size} bytes, "
            f"file has {max(0, len(raw) - offset)} remaining"
        )
    return raw[offset : offset + size]


def parse_index_data(path: Path, raw: bytes) -> LegacyIndex:
    if not raw:
        raise ValueError(f"{path}: empty unit index")
    header_size = raw[0]
    header = require_bytes(raw, 1, header_size, "version header")
    parts = header.rstrip(b"\0").split(b"\0")
    if len(parts) != 2:
        raise ValueError(f"{path}: expected NUL-separated version and producer fields")
    try:
        version, producer = (part.decode("ascii") for part in parts)
    except UnicodeDecodeError as exc:
        raise ValueError(f"{path}: version header is not ASCII") from exc
    if version not in SUPPORTED_VERSIONS:
        raise ValueError(f"{path}: unsupported legacy version {version!r}")

    offset = 1 + header_size
    bank_count = struct.unpack("<H", require_bytes(raw, offset, 2, "bank count"))[0]
    offset += 2
    if bank_count != 1:
        raise ValueError(f"{path}: expected one bank entry, found {bank_count}")
    name_size = struct.unpack("<H", require_bytes(raw, offset, 2, "bank-name size"))[0]
    offset += 2
    name_bytes = require_bytes(raw, offset, name_size, "bank name")
    offset += name_size
    try:
        bank_name = name_bytes.decode("ascii")
    except UnicodeDecodeError as exc:
        raise ValueError(f"{path}: bank name is not ASCII") from exc
    tag = require_bytes(raw, offset, 1, "bank tag")[0]
    offset += 1
    banks = (Bank(bank_name, tag),)

    unit_count = struct.unpack("<I", require_bytes(raw, offset, 4, "unit count"))[0]
    offset += 4
    record_stride = struct.unpack("<H", require_bytes(raw, offset, 2, "record stride"))[0]
    offset += 2
    if record_stride != UNIT_RECORD_STRIDE:
        raise ValueError(
            f"{path}: expected {UNIT_RECORD_STRIDE}-byte records, got {record_stride}"
        )
    if unit_count == 0:
        raise ValueError(f"{path}: unit count must be positive")

    records_size = unit_count * record_stride
    records = require_bytes(raw, offset, records_size, "unit records")
    expected_size = offset + records_size + unit_count * FEATURE_COLUMN_BYTES
    if len(raw) != expected_size:
        raise ValueError(
            f"{path}: expected {expected_size} bytes for {record_stride}-byte records "
            f"and {FEATURE_COLUMN_BYTES} feature-column bytes per unit, "
            f"found {len(raw)}"
        )

    return LegacyIndex(
        path=path,
        version=version,
        producer=producer,
        banks=banks,
        unit_count=unit_count,
        record_stride=record_stride,
        feature_column_bytes=FEATURE_COLUMN_BYTES,
        block_start=offset,
        records=records,
        actual_size=len(raw),
    )


def inspect_payloads(index: LegacyIndex, data_dir: Path) -> dict[str, object]:
    bank_name = index.banks[0].name
    if not bank_name.startswith("merged-"):
        raise ValueError(f"{index.path}: unsupported bank name {bank_name!r}")
    bank = bank_name.removeprefix("merged-")
    dat_path = data_dir / f"merged-{bank}.dat"
    upm_path = data_dir / f"merged-{bank}.upm"
    dat = dat_path.read_bytes()
    upm = upm_path.read_bytes()
    dat_spans: list[tuple[int, int]] = []
    upm_spans: list[tuple[int, int]] = []

    for unit in range(index.unit_count):
        record_start = unit * index.record_stride
        record = index.records[record_start : record_start + index.record_stride]
        dat_offset = struct.unpack_from("<I", record, 0)[0]
        dat_length = struct.unpack_from("<H", record, 8)[0]
        upm_offset = struct.unpack_from("<I", record, 10)[0]
        upm_length = record[14] + record[15] - 1
        if dat_length == 0 or upm_length <= 0:
            raise ValueError(f"{bank} unit {unit}: empty DAT or UPM span")
        if dat_offset + dat_length > len(dat) or upm_offset + upm_length > len(upm):
            raise ValueError(f"{bank} unit {unit}: DAT/UPM span is out of bounds")
        dat_spans.append((dat_offset, dat_length))
        upm_spans.append((upm_offset, upm_length))

    if dat_spans[0][0] != 0 or dat_spans[-1][0] + dat_spans[-1][1] != len(dat):
        raise ValueError(f"{bank}: DAT spans do not begin at zero and end at EOF")
    if upm_spans[0][0] != 0 or upm_spans[-1][0] + upm_spans[-1][1] != len(upm):
        raise ValueError(f"{bank}: UPM spans do not begin at zero and end at EOF")

    dat_contiguous = sum(
        start + size == next_start
        for (start, size), (next_start, _) in zip(dat_spans, dat_spans[1:])
    )
    upm_contiguous = sum(
        start + size == next_start
        for (start, size), (next_start, _) in zip(upm_spans, upm_spans[1:])
    )
    if dat_contiguous != index.unit_count - 1:
        raise ValueError(f"{bank}: DAT spans are not contiguous in unit order")
    if upm_contiguous != index.unit_count - 1:
        raise ValueError(f"{bank}: UPM spans are not contiguous in unit order")

    sample_units = sorted({0, index.unit_count // 2, index.unit_count - 1})
    decoded_samples = []
    for unit in sample_units:
        record_start = unit * index.record_stride
        record = index.records[record_start : record_start + index.record_stride]
        dat_offset = struct.unpack_from("<I", record, 0)[0]
        dat_length = struct.unpack_from("<H", record, 8)[0]
        upm_offset = struct.unpack_from("<I", record, 10)[0]
        upm_length = record[14] + record[15] - 1
        payload = dat[dat_offset : dat_offset + dat_length]
        decoded_count = len(decode_payload(payload))
        upm_sample_count = 2 * sum(upm[upm_offset : upm_offset + upm_length])
        if decoded_count != upm_sample_count:
            raise ValueError(
                f"{bank} unit {unit}: decoded {decoded_count} samples, "
                f"UPM implies {upm_sample_count}"
            )
        decoded_samples.append(
            {
                "unit": unit,
                "dat_bytes": dat_length,
                "decoded_samples": decoded_count,
                "upm_samples": upm_sample_count,
            }
        )

    return {
        "bank": bank,
        "dat_file_size": len(dat),
        "upm_file_size": len(upm),
        "all_dat_spans_in_bounds": True,
        "all_upm_spans_in_bounds": True,
        "dat_spans_contiguous_through_eof": True,
        "upm_spans_contiguous_through_eof": True,
        "sample_units": decoded_samples,
    }


def inspect_index(path: Path, data_root: Path | None = None) -> dict[str, object]:
    index = parse_index_data(path, path.read_bytes())
    result: dict[str, object] = {
        "file": str(path),
        "version": index.version,
        "producer": index.producer,
        "banks": [{"name": bank.name, "tag": bank.tag} for bank in index.banks],
        "unit_count": index.unit_count,
        "record_stride_bytes": index.record_stride,
        "feature_column_bytes_per_unit": index.feature_column_bytes,
        "records_offset": index.block_start,
        "expected_size_bytes": index.actual_size,
        "size_matches": True,
    }
    if data_root is not None:
        result["payload_checks"] = inspect_payloads(index, data_root / "dat")
    return result


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, required=True, help="voice M16 directory")
    parser.add_argument(
        "--skip-payload-checks",
        action="store_true",
        help="check index layout only; do not read matching DAT/UPM files",
    )
    args = parser.parse_args()
    indexes = sorted((args.root / "mc_idx_tbl").glob("unit-*.idx"))
    if not indexes:
        parser.error(f"no unit-*.idx files found under {args.root / 'mc_idx_tbl'}")
    results = [
        inspect_index(path, None if args.skip_payload_checks else args.root)
        for path in indexes
    ]
    print(json.dumps({"root": str(args.root), "indexes": results}, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
