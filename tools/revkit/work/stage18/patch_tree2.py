#!/usr/bin/env python3
"""Make a guarded, two-byte tree2 path probe copy of vt_kat.dll."""

from hashlib import sha256
from pathlib import Path


ROOT = Path(__file__).resolve().parents[4]
SOURCE = ROOT / "binary" / "vt_kat.dll"
OUTPUT = Path(__file__).with_name("vt_kat-tree2-probe.dll")
REPORT = Path(__file__).with_name("patch-manifest.txt")
EXPECTED_SOURCE_SHA256 = "f93c60a7ab0baa3f4ec2271cdc46997f48982143f49088349c32ff1755388461"
PATCHES = (
    (0x191F, b"\x33", b"\x32", "FUN_10001900 tree suffix digit 3 to 2"),
    (0x7C8F0, b"\x33", b"\x32", "tree-directory template digit 3 to 2"),
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
    output_hash = sha256(patched).hexdigest()
    rows = [
        f"source={SOURCE.relative_to(ROOT)}",
        f"source_sha256={source_hash}",
        f"probe={OUTPUT.name}",
        f"probe_sha256={output_hash}",
        "changes=",
    ]
    rows.extend(
        f"  file_offset=0x{offset:x} {before.hex()}->{after.hex()} {description}"
        for offset, before, after, description in PATCHES
    )
    REPORT.write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
