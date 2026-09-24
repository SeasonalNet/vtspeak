#!/usr/bin/env python3
"""Compare local DAT decoder PCM hashes with the isolated DLL corpus capture."""

from __future__ import annotations

import argparse
import hashlib
import struct
import sys
from pathlib import Path
from typing import BinaryIO

from decode_dat import decode_payload

HEADER_SIZE = 45
RECORD_SIZE = 19
FEATURE_SIZE = 21
RESULT_MAGIC = b"VTDIGEST1"
RESULT_SIZE = 4 + 32
BANKS = ("gen", "num", "etc", "alp")


def compare_bank(root: Path, bank: str, capture: BinaryIO) -> int:
    index_path = root / "mc_idx_tbl" / f"unit-{bank}.idx"
    dat_path = root / "dat" / f"merged-{bank}.dat"
    index = index_path.read_bytes()
    dat = dat_path.read_bytes()
    if len(index) < HEADER_SIZE or not index[1:24].startswith(
        b"ver.2013\0VoiceText-Eng\0"
    ):
        raise ValueError(f"unsupported unit index: {index_path}")
    unit_count = struct.unpack_from("<I", index, 39)[0]
    stride = struct.unpack_from("<H", index, 43)[0]
    if stride != RECORD_SIZE or len(index) != HEADER_SIZE + unit_count * (stride + FEATURE_SIZE):
        raise ValueError(f"unsupported unit index dimensions: {index_path}")

    for unit in range(unit_count):
        record_start = HEADER_SIZE + unit * stride
        dat_offset = struct.unpack_from("<I", index, record_start)[0]
        dat_length = struct.unpack_from("<H", index, record_start + 8)[0]
        if dat_offset + dat_length > len(dat):
            raise ValueError(f"{bank} unit {unit}: DAT span exceeds bank")
        result = capture.read(RESULT_SIZE)
        if len(result) != RESULT_SIZE:
            raise ValueError(f"{bank} unit {unit}: truncated DLL result stream")
        dll_bytes = struct.unpack_from("<I", result)[0]
        dll_digest = result[4:]
        samples = decode_payload(dat[dat_offset : dat_offset + dat_length])
        pcm = struct.pack(f"<{len(samples)}h", *samples)
        local_digest = hashlib.sha256(pcm).digest()
        if len(pcm) != dll_bytes or local_digest != dll_digest:
            raise ValueError(
                f"{bank} unit {unit}: PCM mismatch; DLL bytes={dll_bytes}, "
                f"local bytes={len(pcm)}, DLL SHA-256={dll_digest.hex()}, "
                f"local SHA-256={local_digest.hex()}"
            )
        if unit and unit % 10_000 == 0:
            print(f"{bank}: compared {unit}/{unit_count} units", file=sys.stderr, flush=True)
    print(f"{bank}: matched {unit_count} units", flush=True)
    return unit_count


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=Path("data-paul/M16"))
    parser.add_argument(
        "--capture",
        type=Path,
        default=Path("tools/revkit/work/corpus-parity/dll-sha256.bin"),
    )
    args = parser.parse_args()

    with args.capture.open("rb") as capture:
        if capture.read(len(RESULT_MAGIC)) != RESULT_MAGIC:
            raise ValueError(f"unsupported DLL result stream: {args.capture}")
        total = sum(compare_bank(args.root, bank, capture) for bank in BANKS)
        if capture.read(1):
            raise ValueError("DLL result stream has trailing bytes")
    print(f"matched {total} DLL/local PCM SHA-256 pairs")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
