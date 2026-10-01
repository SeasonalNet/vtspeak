#!/usr/bin/env python3
"""Build a disposable one-word metric-code intervention for Apple."""

from __future__ import annotations

import argparse
import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
BASE = (
    ROOT
    / "tools/revkit/work/stage20/"
    "index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks"
)
sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))

from inspect_unit_idx import inspect  # noqa: E402

UNIT_ID = 100_180
OLD_CODE = 744
MATCH_CODE = 884


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output-dir", type=Path, required=True)
    args = parser.parse_args()
    output = args.output_dir.resolve()
    if output == BASE.resolve() or BASE.resolve() in output.parents:
        raise ValueError("output directory must not overlap the source overlay")
    output.mkdir(parents=True, exist_ok=True)

    manifest = [
        "intervention=unit-gen row 100180 group-0 metric code 744 to 884",
        "rationale=match predecessor metric code 884 while preserving the class key",
        "scope=one little-endian u16 in one versioned index; all other files copied byte-for-byte",
    ]
    for source in sorted(BASE.glob("unit-*.idx")):
        destination = output / source.name
        shutil.copyfile(source, destination)
        base_hash = sha256(source)
        if source.name == "unit-gen.idx":
            metadata = inspect(source)
            count = int(metadata["unit_count"])
            if count <= UNIT_ID or int(metadata["opaque_block_stride_bytes"]) != 19:
                raise ValueError("unit-gen layout does not contain the expected row")
            block_start = int(metadata["opaque_block_start"])
            offset = block_start + count * (19 + 9) + UNIT_ID * 2
            raw = bytearray(destination.read_bytes())
            current = int.from_bytes(raw[offset : offset + 2], "little") & 0x3FFF
            if current != OLD_CODE:
                raise ValueError(f"expected metric code {OLD_CODE}, found {current}")
            raw[offset : offset + 2] = MATCH_CODE.to_bytes(2, "little")
            destination.write_bytes(raw)
            changed = [i for i, (a, b) in enumerate(zip(source.read_bytes(), raw)) if a != b]
            if changed != [offset, offset + 1]:
                raise ValueError(f"unexpected changed byte offsets: {changed}")
            manifest.append(
                f"{source.name}: sha256={base_hash}->{sha256(destination)} "
                f"offset={offset} bytes_changed={len(changed)} code={current}->{MATCH_CODE}"
            )
        else:
            manifest.append(f"{source.name}: sha256={base_hash}->{sha256(destination)} unchanged")
        inspect(destination)

    names = {path.name for path in output.glob("unit-*.idx")}
    expected = {path.name for path in BASE.glob("unit-*.idx")}
    if names != expected:
        raise ValueError(f"index set changed: expected {sorted(expected)}, found {sorted(names)}")
    (output / "apple-row100180-metric-match-manifest.txt").write_text(
        "\n".join(manifest) + "\n"
    )
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
