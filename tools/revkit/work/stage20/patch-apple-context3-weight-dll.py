#!/usr/bin/env python3
"""Create a hash-guarded DLL with the measured context-3 B weight."""

from __future__ import annotations

import argparse
import hashlib
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = ROOT / "tools/revkit/work/stage19/vt_kat-tree2-versioned.dll"
OUTPUT = WORK / "vt_kat-context3-bweight10.dll"
OUTPUT_CONTEXT2_CONTEXT3 = WORK / "vt_kat-context2-aweight1-context3-bweight10.dll"
EXPECTED_SOURCE_SHA256 = "9f006d1474cc912535699bd94c04d5d7649a4489431059fa8b816ebcd10421b5"
WEIGHT_PATCHES = (
    (0x7C298, bytes.fromhex("0000a040"), bytes.fromhex("00002041"), "context-3 slot-1 derived-B 5.0->10.0"),
)
CONTEXT2_PATCH = (
    0x7C2A8,
    bytes.fromhex("00000040"),
    bytes.fromhex("0000803f"),
    "context-2 slot-2 derived-A 2.0->1.0",
)


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--include-context2",
        action="store_true",
        help="also patch the measured context-2 slot-2 derived-A coefficient",
    )
    args = parser.parse_args()
    output = OUTPUT_CONTEXT2_CONTEXT3 if args.include_context2 else OUTPUT

    source = SOURCE.read_bytes()
    source_hash = sha256(source)
    if source_hash != EXPECTED_SOURCE_SHA256:
        raise SystemExit(f"unexpected versioned DLL hash: {source_hash}")
    patches = list(WEIGHT_PATCHES)
    if args.include_context2:
        patches.append(CONTEXT2_PATCH)
    for file_offset, before, _, description in patches:
        actual = source[file_offset : file_offset + 4]
        if actual != before:
            raise SystemExit(
                f"unexpected coefficient at file offset 0x{file_offset:x} "
                f"({description}): expected {before.hex()}, found {actual.hex()}"
            )
    if output.exists():
        raise FileExistsError(f"refusing to overwrite existing patched DLL: {output}")

    patched = bytearray(source)
    manifest_lines = [
        f"source={SOURCE.relative_to(ROOT)}",
        f"source_sha256={source_hash}",
    ]
    for file_offset, before, after, description in patches:
        patched[file_offset : file_offset + 4] = after
        manifest_lines.append(
            f"file_offset=0x{file_offset:x} virtual_address=0x{0x10000000 + file_offset:x} "
            f"{description} bytes={before.hex()}->{after.hex()}"
        )
    output.write_bytes(patched)
    manifest_lines.append(f"patched_sha256={sha256(patched)}")
    manifest = "\n".join(manifest_lines) + "\n"
    output.with_suffix(".manifest.txt").write_text(manifest)
    print(manifest, end="")


if __name__ == "__main__":
    main()
