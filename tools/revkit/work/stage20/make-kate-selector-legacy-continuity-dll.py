#!/usr/bin/env python3
"""Create a disposable Kate DLL with legacy attr_48 span reads."""

from __future__ import annotations

import argparse
from hashlib import sha256
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
SOURCE = ROOT / "binary" / "vt_kat.dll"
EXPECTED_SOURCE_SHA256 = "f93c60a7ab0baa3f4ec2271cdc46997f48982143f49088349c32ff1755388461"

# FUN_100230a0 computes each signature row's byte +6 at these two sites. The
# compatibility experiment instead reads legacy attr_48 from signature byte 0.
PATCHES = (
    (0x191F, b"\x33", b"\x32", "tree filename suffix 3 to 2"),
    (0x7C8F0, b"\x33", b"\x32", "voice tree directory tree3 to tree2"),
    (0x23155, b"\x06", b"\x00", "left span scan signature byte +6 to +0"),
    (0x23170, b"\x80", b"\x01", "left span scan marker mask 0x80 to 0x01"),
    (0x2320A, b"\x06", b"\x00", "right span scan signature byte +6 to +0"),
    (0x23218, b"\x80", b"\x01", "right span scan marker mask 0x80 to 0x01"),
)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--output",
        type=Path,
    )
    parser.add_argument(
        "--legacy-cutoff",
        action="store_true",
        help="also change FUN_10023e70's two sum cutoffs from >=10 to >=3",
    )
    parser.add_argument(
        "--legacy-active-cutoff",
        action="store_true",
        help="change FUN_10024060's ordinary whole-position cutoff from >=10 to >=3",
    )
    args = parser.parse_args()

    output = args.output or Path(
        "/tmp/vtspeak-vt-kat-legacy-continuity-byte0-cutoff3.dll"
        if args.legacy_cutoff
        else "/tmp/vtspeak-vt-kat-legacy-continuity-byte0-active-cutoff3.dll"
        if args.legacy_active_cutoff
        else "/tmp/vtspeak-vt-kat-legacy-continuity-byte0.dll"
    )

    original = SOURCE.read_bytes()
    source_hash = sha256(original).hexdigest()
    if source_hash != EXPECTED_SOURCE_SHA256:
        raise SystemExit(f"unexpected source DLL hash: {source_hash}")

    patched = bytearray(original)
    manifest = [f"source_sha256={source_hash}", f"output={output}"]
    selected_patches = list(PATCHES)
    if args.legacy_cutoff:
        selected_patches.extend(
            (
                (0x23F0F, b"\x0a", b"\x03", "primary accumulated-metric cutoff 10 to 3"),
                (0x23F6E, b"\x0a", b"\x03", "retry accumulated-metric cutoff 10 to 3"),
            )
        )
    if args.legacy_active_cutoff:
        selected_patches.append(
            (0x2428F, b"\x0a", b"\x03", "FUN_10024060 ordinary whole-position cutoff 10 to 3")
        )
    for offset, before, after, description in selected_patches:
        actual = bytes(patched[offset : offset + len(before)])
        if actual != before:
            raise SystemExit(
                f"unexpected bytes at file offset 0x{offset:x}: "
                f"expected {before.hex()}, found {actual.hex()}"
            )
        patched[offset : offset + len(before)] = after
        manifest.append(
            f"file_offset=0x{offset:x} {before.hex()}->{after.hex()} {description}"
        )

    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_bytes(patched)
    manifest.append(f"patched_sha256={sha256(patched).hexdigest()}")
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
