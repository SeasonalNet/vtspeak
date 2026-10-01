#!/usr/bin/env python3
"""Set continuation bit 7 on only the three measured Apple predecessor rows."""

from __future__ import annotations

import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
BASE = WORK / "index-adapter-key-repacked-attr48-key3-attr40-key2-apple-context-prefix-tail-map"
OUTPUT = WORK / "index-adapter-key-repacked-attr48-key3-attr40-key2-apple-selected-predecessors-bit7"
UNIT_IDS = (59558, 82036, 2554)


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
    from inspect_unit_idx import inspect

    if OUTPUT.exists():
        raise FileExistsError(f"refusing to overwrite existing overlay: {OUTPUT}")
    shutil.copytree(BASE, OUTPUT)

    path = OUTPUT / "unit-gen.idx"
    raw = bytearray(path.read_bytes())
    metadata = inspect(path)
    count = int(metadata["unit_count"])
    stride = int(metadata["opaque_block_stride_bytes"])
    signature_start = (
        int(metadata["opaque_block_start"]) + count * stride + count
    )
    if max(UNIT_IDS) >= count:
        raise ValueError(f"selected ID exceeds unit-gen count {count}")

    before_hash = sha256(raw)
    changes: list[str] = []
    for unit_id in UNIT_IDS:
        offset = signature_start + unit_id * 7 + 6
        before = raw[offset]
        if before & 0x80:
            raise ValueError(f"unit {unit_id} already has bit 7 set: {before:#04x}")
        raw[offset] = before | 0x80
        changes.append(
            f"unit={unit_id} row={unit_id} offset=0x{offset:x} "
            f"before={before:#04x} after={raw[offset]:#04x}"
        )

    path.write_bytes(raw)
    checked = inspect(path)
    if checked["unit_count"] != count or not checked["size_matches"]:
        raise ValueError("targeted overlay failed versioned-index structure validation")
    changed = [
        offset
        for offset, (old, new) in enumerate(zip((BASE / "unit-gen.idx").read_bytes(), raw))
        if old != new
    ]
    expected = sorted(signature_start + unit_id * 7 + 6 for unit_id in UNIT_IDS)
    if changed != expected:
        raise ValueError(f"unexpected byte changes: {changed!r}")

    manifest = [
        "variant=Apple prefix/tail-map index with bit 7 set on three measured predecessor signatures",
        f"base={BASE.relative_to(ROOT)}",
        "bank=unit-gen",
        f"unit_count={count}",
        f"signature_start=0x{signature_start:x}",
        f"base_sha256={before_hash}",
        f"overlay_sha256={sha256(raw)}",
        f"changed_byte_count={len(changed)}",
        *changes,
    ]
    (OUTPUT / "selected-predecessors-bit7-manifest.txt").write_text(
        "\n".join(manifest) + "\n"
    )
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
