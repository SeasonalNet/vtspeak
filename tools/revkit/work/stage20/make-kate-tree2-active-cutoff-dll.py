#!/usr/bin/env python3
"""Make a disposable tree2-compatible Kate DLL with the legacy split cutoff."""

from __future__ import annotations

from hashlib import sha256
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
SOURCE = ROOT / "tools" / "revkit" / "work" / "stage19" / "vt_kat-tree2-versioned.dll"
EXPECTED_SOURCE_SHA256 = "9f006d1474cc912535699bd94c04d5d7649a4489431059fa8b816ebcd10421b5"
OUTPUT = Path("/tmp/vtspeak-vt-kat-tree2-active-cutoff3.dll")
ACTIVE_CUTOFF_OFFSET = 0x2428F


def main() -> None:
    original = SOURCE.read_bytes()
    source_hash = sha256(original).hexdigest()
    if source_hash != EXPECTED_SOURCE_SHA256:
        raise SystemExit(f"unexpected source DLL hash: {source_hash}")

    patched = bytearray(original)
    before = bytes(patched[ACTIVE_CUTOFF_OFFSET : ACTIVE_CUTOFF_OFFSET + 1])
    if before != b"\x0a":
        raise SystemExit(
            f"unexpected cutoff byte at 0x{ACTIVE_CUTOFF_OFFSET:x}: {before.hex()}"
        )
    patched[ACTIVE_CUTOFF_OFFSET] = 0x03
    OUTPUT.write_bytes(patched)
    print(f"source_sha256={source_hash}")
    print(f"output={OUTPUT}")
    print(f"file_offset=0x{ACTIVE_CUTOFF_OFFSET:x} 0a->03 ordinary whole-position cutoff")
    print(f"patched_sha256={sha256(patched).hexdigest()}")


if __name__ == "__main__":
    main()
