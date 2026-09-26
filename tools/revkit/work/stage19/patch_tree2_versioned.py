#!/usr/bin/env python3
"""Build a disposable tree2-path DLL while retaining the versioned index path."""

from hashlib import sha256
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = ROOT / "binary" / "vt_kat.dll"
OUTPUT = WORK / "vt_kat-tree2-versioned.dll"
EXPECTED_SOURCE_SHA256 = "f93c60a7ab0baa3f4ec2271cdc46997f48982143f49088349c32ff1755388461"
PATCHES = (
    (0x191F, b"\x33", b"\x32", "generated tree suffix 3 to 2"),
    (0x7C8F0, b"\x33", b"\x32", "voice tree directory tree3 to tree2"),
)


def main() -> None:
    original = SOURCE.read_bytes()
    source_hash = sha256(original).hexdigest()
    if source_hash != EXPECTED_SOURCE_SHA256:
        raise SystemExit(f"unexpected source DLL hash: {source_hash}")
    patched = bytearray(original)
    for offset, before, after, _description in PATCHES:
        if patched[offset : offset + len(before)] != before:
            raise SystemExit(f"unexpected source bytes at file offset 0x{offset:x}")
        patched[offset : offset + len(before)] = after
    OUTPUT.write_bytes(patched)

    rows = [
        f"source_sha256={source_hash}",
        f"patched_copy={OUTPUT.name}",
        f"patched_sha256={sha256(patched).hexdigest()}",
        "index_guard=unchanged; selects versioned reader for converted ver.2013 indexes",
    ]
    rows.extend(
        f"file_offset=0x{offset:x} {before.hex()}->{after.hex()} {description}"
        for offset, before, after, description in PATCHES
    )
    (WORK / "versioned-patch-manifest.txt").write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
