#!/usr/bin/env python3
"""Create a disposable 2006 Kate comparison DLL that skips demo-text insertion."""

from __future__ import annotations

import hashlib
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll"
OUTPUT = ROOT / "vt_eng-banner-free-2006.dll"
EXPECTED_SOURCE_SHA256 = "00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d"
FILE_OFFSET = 0x15957
EXPECTED_BYTES = bytes.fromhex("0f 85 e7 01 00 00")
PATCHED_BYTES = bytes.fromhex("e9 e8 01 00 00 90")


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    source = SOURCE.read_bytes()
    source_hash = sha256(source)
    if source_hash != EXPECTED_SOURCE_SHA256:
        raise SystemExit(f"unexpected source SHA-256: {source_hash}")
    actual = source[FILE_OFFSET : FILE_OFFSET + len(EXPECTED_BYTES)]
    if actual != EXPECTED_BYTES:
        raise SystemExit(
            f"unexpected bytes at file offset 0x{FILE_OFFSET:x}: {actual.hex(' ')}"
        )

    patched = bytearray(source)
    patched[FILE_OFFSET : FILE_OFFSET + len(PATCHED_BYTES)] = PATCHED_BYTES
    OUTPUT.write_bytes(patched)
    print(f"source_sha256={source_hash}")
    print(f"output={OUTPUT}")
    print(f"output_sha256={sha256(patched)}")
    print(
        f"patch=VA 0x10015957 / file offset 0x{FILE_OFFSET:x}: "
        f"{EXPECTED_BYTES.hex(' ')} -> {PATCHED_BYTES.hex(' ')}"
    )


if __name__ == "__main__":
    main()
