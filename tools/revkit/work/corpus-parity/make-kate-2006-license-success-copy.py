#!/usr/bin/env python3
"""Create a guarded, disposable 2006 Kate DLL with the check-success branch."""

from __future__ import annotations

import hashlib
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll"
OUTPUT = ROOT / "vt_eng-license-success-2006.dll"
EXPECTED_SOURCE_SHA256 = "00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d"
FILE_OFFSET = 0x22542
EXPECTED_BYTES = bytes.fromhex("7d 0e")
PATCHED_BYTES = bytes.fromhex("eb 0e")
LIMIT_FILE_OFFSET = 0x22010
LIMIT_EXPECTED_BYTES = bytes.fromhex(
    "83 ec 0c 53 55 56 bb 01 00 00 00 57 53 e8 be 03"
)
LIMIT_PATCHED_BYTES = bytes.fromhex(
    "c7 05 f0 62 09 10 01 00 00 00 b8 01 00 00 00 c3"
)


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
    limit_actual = source[
        LIMIT_FILE_OFFSET : LIMIT_FILE_OFFSET + len(LIMIT_EXPECTED_BYTES)
    ]
    if limit_actual != LIMIT_EXPECTED_BYTES:
        raise SystemExit(
            "unexpected license-limit helper bytes at "
            f"file offset 0x{LIMIT_FILE_OFFSET:x}: {limit_actual.hex(' ')}"
        )

    patched = bytearray(source)
    patched[FILE_OFFSET : FILE_OFFSET + len(PATCHED_BYTES)] = PATCHED_BYTES
    patched[
        LIMIT_FILE_OFFSET : LIMIT_FILE_OFFSET + len(LIMIT_PATCHED_BYTES)
    ] = LIMIT_PATCHED_BYTES
    OUTPUT.write_bytes(patched)
    print(f"source={SOURCE}")
    print(f"source_sha256={source_hash}")
    print(f"output={OUTPUT}")
    print(f"output_sha256={sha256(patched)}")
    print(
        f"patch=VA 0x10022542 / file offset 0x{FILE_OFFSET:x}: "
        f"{EXPECTED_BYTES.hex(' ')} -> {PATCHED_BYTES.hex(' ')}"
    )
    print(
        f"patch=VA 0x10022010 / file offset 0x{LIMIT_FILE_OFFSET:x}: "
        f"license-limit helper returns 1 and stores limit 1"
    )


if __name__ == "__main__":
    main()
