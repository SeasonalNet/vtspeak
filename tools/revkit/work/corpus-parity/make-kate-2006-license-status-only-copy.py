#!/usr/bin/env python3
"""Create a disposable 2006 Kate DLL that changes only checker status."""

from __future__ import annotations

import hashlib
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SOURCE = ROOT / "kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll"
OUTPUT = ROOT / "vt_eng-license-status-only-2006.dll"
EXPECTED_SOURCE_SHA256 = "00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d"
FILE_OFFSET = 0x22542
EXPECTED_BYTES = bytes.fromhex("7d 0e")
PATCHED_BYTES = bytes.fromhex("eb 0e")

source = SOURCE.read_bytes()
source_hash = hashlib.sha256(source).hexdigest()
if source_hash != EXPECTED_SOURCE_SHA256:
    raise SystemExit(f"unexpected source SHA-256: {source_hash}")
if source[FILE_OFFSET : FILE_OFFSET + 2] != EXPECTED_BYTES:
    raise SystemExit("unexpected checker-status branch bytes")
patched = bytearray(source)
patched[FILE_OFFSET : FILE_OFFSET + 2] = PATCHED_BYTES
OUTPUT.write_bytes(patched)
print(f"output={OUTPUT}")
print(f"output_sha256={hashlib.sha256(patched).hexdigest()}")
