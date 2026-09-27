#!/usr/bin/env python3
"""Compare sampled patched-DLL Kate PCM with the independent DAT decoder."""

from __future__ import annotations

import hashlib
import struct
from pathlib import Path

from decode_dat import decode_payload

ROOT = Path("tools/revkit/work/stage19/data-kate-copy/M16")
CAPTURE = Path("tools/revkit/work/corpus-parity/kate/kate-dat-dll-parity.tsv")
BANKS = ("gen", "gen2", "num", "etc", "alp")


def main() -> int:
    rows = [line.split("\t") for line in CAPTURE.read_text().splitlines() if line]
    row_iter = iter(rows)
    checked = 0
    for bank in BANKS:
        idx = (ROOT / "mc_idx_tbl" / f"unit-{bank}.idx").read_bytes()
        dat = (ROOT / "dat" / f"merged-{bank}.dat").read_bytes()
        upm = (ROOT / "dat" / f"merged-{bank}.upm").read_bytes()
        header_len = idx[0]
        pos = 1 + header_len + 2
        name_len = struct.unpack_from("<H", idx, pos)[0]
        pos += 2 + name_len + 1
        count = struct.unpack_from("<I", idx, pos)[0]
        stride = struct.unpack_from("<H", idx, pos + 4)[0]
        pos += 6
        if stride != 19 or len(idx) != pos + count * 39:
            raise ValueError(f"{bank}: unexpected legacy index dimensions")
        units = range(count)
        for unit in units:
            row = next(row_iter)
            if len(row) != 4 or row[:2] != [bank, str(unit)]:
                raise ValueError(f"capture order mismatch at {bank}/{unit}: {row}")
            record = idx[pos + unit * stride : pos + (unit + 1) * stride]
            dat_offset = struct.unpack_from("<I", record, 0)[0]
            dat_length = struct.unpack_from("<H", record, 8)[0]
            upm_offset = struct.unpack_from("<I", record, 10)[0]
            upm_length = record[14] + record[15] - 1
            samples = decode_payload(dat[dat_offset : dat_offset + dat_length])
            pcm = struct.pack(f"<{len(samples)}h", *samples)
            dll_bytes = int(row[2])
            dll_hash = row[3]
            local_hash = hashlib.sha256(pcm).hexdigest()
            upm_bytes = 4 * sum(upm[upm_offset : upm_offset + upm_length])
            if len(pcm) != dll_bytes or local_hash != dll_hash or len(pcm) != upm_bytes:
                raise ValueError(
                    f"{bank}/{unit}: DLL={dll_bytes}/{dll_hash}, "
                    f"local={len(pcm)}/{local_hash}, UPM={upm_bytes}"
                )
            checked += 1
        print(f"{bank}: {count} units matched DLL/local PCM SHA-256 and UPM bytes")
    try:
        next(row_iter)
    except StopIteration:
        print(f"matched {checked} units")
        return 0
    raise ValueError("capture has trailing rows")


if __name__ == "__main__":
    raise SystemExit(main())
