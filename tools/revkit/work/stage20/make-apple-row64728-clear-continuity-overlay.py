#!/usr/bin/env python3
"""Clear one Apple row's continuity marker in a disposable index overlay."""

from __future__ import annotations

import argparse
import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_unit_idx import inspect as inspect_versioned  # noqa: E402

SOURCE = (
    ROOT
    / "tools"
    / "revkit"
    / "work"
    / "stage20"
    / "index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks"
)
GLOBAL_UNIT = 64_728
SIGNATURE_BYTE = 6
CONTINUITY_BIT = 0x80


def sha256(raw: bytes) -> str:
    return hashlib.sha256(raw).hexdigest()


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--output-dir",
        type=Path,
        default=Path("/tmp/vtspeak-apple-row64728-clear-bit7"),
    )
    args = parser.parse_args()
    output = args.output_dir.resolve()
    if output.exists():
        raise FileExistsError(f"refusing to overwrite existing overlay: {output}")
    shutil.copytree(SOURCE, output)

    index_path = output / "unit-gen.idx"
    raw = index_path.read_bytes()
    info = inspect_versioned(index_path)
    unit_count = int(info["unit_count"])
    if GLOBAL_UNIT >= unit_count:
        raise ValueError(f"unit {GLOBAL_UNIT} is outside unit-gen ({unit_count} rows)")
    signature_start = int(info["opaque_block_start"]) + 20 * unit_count
    byte_offset = signature_start + GLOBAL_UNIT * 7 + SIGNATURE_BYTE
    before_hash = sha256(raw)
    old_byte = raw[byte_offset]
    if not old_byte & CONTINUITY_BIT:
        raise ValueError(
            f"unit {GLOBAL_UNIT} signature byte {SIGNATURE_BYTE} lacks bit 0x80: "
            f"0x{old_byte:02x}"
        )
    changed = bytearray(raw)
    changed[byte_offset] &= ~CONTINUITY_BIT
    index_path.write_bytes(changed)
    after = index_path.read_bytes()
    if sum(left != right for left, right in zip(raw, after)) != 1:
        raise ValueError("expected exactly one changed byte")
    checked = inspect_versioned(index_path)
    if int(checked["unit_count"]) != unit_count:
        raise ValueError("unit count changed")

    (output / "apple-row64728-clear-bit7-manifest.txt").write_text(
        "base=matched-context attr_b transfer plus full-bank continuity overlay\n"
        f"unit=unit-gen row {GLOBAL_UNIT}\n"
        f"signature_byte={SIGNATURE_BYTE}\n"
        f"byte_offset={byte_offset}\n"
        f"old_byte=0x{old_byte:02x}\n"
        f"new_byte=0x{after[byte_offset]:02x}\n"
        "changed_bytes=1\n"
        f"unit_gen_sha256_before={before_hash}\n"
        f"unit_gen_sha256_after={sha256(after)}\n"
        "versioned_index_validation=passed\n"
    )
    print((output / "apple-row64728-clear-bit7-manifest.txt").read_text(), end="")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
