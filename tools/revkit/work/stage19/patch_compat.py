#!/usr/bin/env python3
"""Build an isolated current-DLL copy and common-tree suffix aliases."""

from hashlib import sha256
from pathlib import Path
import shutil

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = ROOT / "binary" / "vt_kat.dll"
OUTPUT = WORK / "vt_kat-tree2-compat.dll"
COMMON_SOURCE = ROOT / "data-common" / "dict-eng"
COMMON_OUTPUT = WORK / "data-common" / "dict-eng"
EXPECTED_SOURCE_SHA256 = "f93c60a7ab0baa3f4ec2271cdc46997f48982143f49088349c32ff1755388461"
PATCHES = (
    (0x191F, b"\x33", b"\x32", "FUN_10001900 generated suffix 3 to 2"),
    (0x7C8F0, b"\x33", b"\x32", "voice tree directory tree3 to tree2"),
    (0x7C59C, b"v", b"x", "version prefix guard selects legacy unit-index reader"),
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

    if COMMON_OUTPUT.exists():
        shutil.rmtree(COMMON_OUTPUT)
    shutil.copytree(COMMON_SOURCE, COMMON_OUTPUT)
    aliases = []
    for source in sorted(COMMON_OUTPUT.glob("*.tree3")):
        alias = source.with_suffix(".tree2")
        shutil.copyfile(source, alias)
        aliases.append(alias.name)

    output_hash = sha256(patched).hexdigest()
    rows = [
        f"source={SOURCE.relative_to(ROOT)}",
        f"source_sha256={source_hash}",
        f"patched_copy={OUTPUT.name}",
        f"patched_sha256={output_hash}",
        "patches=",
    ]
    rows.extend(
        f"  file_offset=0x{offset:x} {before.hex()}->{after.hex()} {description}"
        for offset, before, after, description in PATCHES
    )
    rows.append(f"common_tree3_as_tree2_aliases={','.join(aliases)}")
    (WORK / "patch-manifest.txt").write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
