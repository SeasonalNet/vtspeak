#!/usr/bin/env python3
"""Build an ignored, comparison-only copy of the 2006 Kate engine DLL."""

from __future__ import annotations

import hashlib
from pathlib import Path

HERE = Path(__file__).resolve().parent
WORK = HERE.parents[1]
CORPUS = WORK / "corpus-parity"
SOURCE = (
    CORPUS
    / "kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll"
)
OUTPUT = CORPUS / "vt_eng-license-success-2006.dll"
SOURCE_SHA256 = "00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d"

# FUN_10022430: force its final checker-status branch to return zero.
STATUS_OFFSET = 0x22542
STATUS_BEFORE = bytes.fromhex("7d 0e")
STATUS_AFTER = bytes.fromhex("eb 0e")


def digest(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    original = SOURCE.read_bytes()
    if digest(original) != SOURCE_SHA256:
        raise SystemExit("the extracted 2006 source DLL hash did not match")
    if original[STATUS_OFFSET : STATUS_OFFSET + len(STATUS_BEFORE)] != STATUS_BEFORE:
        raise SystemExit("unexpected checker-status branch bytes")
    patched = bytearray(original)
    patched[STATUS_OFFSET : STATUS_OFFSET + len(STATUS_AFTER)] = STATUS_AFTER
    OUTPUT.write_bytes(patched)
    print(f"source_sha256={digest(original)}")
    print(f"output={OUTPUT}")
    print(f"output_sha256={digest(patched)}")
    print("status_patch_va=0x10022542 bytes=7d0e->eb0e")


if __name__ == "__main__":
    main()
