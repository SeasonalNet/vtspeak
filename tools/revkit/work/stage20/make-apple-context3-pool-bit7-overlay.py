#!/usr/bin/env python3
"""Apply legacy attr_48 markers to the captured Apple context-3 candidate pool."""

from __future__ import annotations

import hashlib
import shutil
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
LEGACY = ROOT / "tools/revkit/work/stage19/data-kate-copy/M16/mc_idx_tbl"
BASE = WORK / "index-adapter-key-repacked-attr48-key3-attr40-key2-apple-context-prefix-tail-map"
OUTPUT = WORK / "index-adapter-key-repacked-attr48-key3-attr40-key2-apple-context3-pool-bit7"
CONTEXT3_POOL_IDS = (
    2554, 256581, 232462, 218670, 158830, 131118, 129559, 103238, 97884, 86327,
    82036, 59558, 41057, 40214, 28631, 27966, 27794, 27440, 11135, 258822,
)
BANKS = ("unit-gen", "unit-gen2", "unit-etc", "unit-alp", "unit-num")


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def main() -> None:
    sys.path.insert(0, str(ROOT / "tools/revkit/scripts"))
    from inspect_legacy_unit_idx import parse_index_data
    from inspect_unit_idx import inspect

    if OUTPUT.exists():
        raise FileExistsError(f"refusing to overwrite existing overlay: {OUTPUT}")
    shutil.copytree(BASE, OUTPUT)

    ranges: list[tuple[int, int, str, bytes]] = []
    offset = 0
    for bank in BANKS:
        source = LEGACY / f"{bank}.idx"
        source_raw = source.read_bytes()
        legacy = parse_index_data(source, source_raw)
        count = legacy.unit_count
        columns_start = legacy.block_start + count * legacy.record_stride
        attr_48 = source_raw[columns_start + count : columns_start + 2 * count]
        converted = inspect(OUTPUT / f"{bank}.idx")
        if converted["unit_count"] != count:
            raise ValueError(f"legacy/versioned row count mismatch for {bank}")
        ranges.append((offset, offset + count, bank, attr_48))
        offset += count

    changes: list[str] = []
    changed_offsets: dict[str, list[int]] = {bank: [] for bank in BANKS}
    for global_id in CONTEXT3_POOL_IDS:
        for low, high, bank, attr_48 in ranges:
            if low <= global_id < high:
                local_row = global_id - low
                if attr_48[local_row] != 1:
                    raise ValueError(
                        f"context-3 candidate {global_id} has legacy attr_48="
                        f"{attr_48[local_row]}"
                    )
                target = OUTPUT / f"{bank}.idx"
                raw = bytearray(target.read_bytes())
                metadata = inspect(target)
                count = int(metadata["unit_count"])
                signature_start = (
                    int(metadata["opaque_block_start"])
                    + count * int(metadata["opaque_block_stride_bytes"])
                    + count
                )
                byte_offset = signature_start + local_row * 7 + 6
                before = raw[byte_offset]
                if before & 0x80:
                    raise ValueError(
                        f"context-3 candidate {global_id} already has bit 7 set"
                    )
                raw[byte_offset] = before | 0x80
                target.write_bytes(raw)
                changed_offsets[bank].append(byte_offset)
                changes.append(
                    f"unit={global_id} bank={bank} row={local_row} attr_48=1 "
                    f"offset=0x{byte_offset:x} before={before:#04x} "
                    f"after={raw[byte_offset]:#04x}"
                )
                break
        else:
            raise ValueError(f"context-3 candidate {global_id} is outside all banks")

    manifest = [
        "variant=Apple adapter with legacy attr_48 markers on all 20 captured context-3 candidates",
        f"base={BASE.relative_to(ROOT)}",
        f"legacy_source={LEGACY.relative_to(ROOT)}",
        f"candidate_count={len(CONTEXT3_POOL_IDS)}",
        f"changed_byte_count={len(changes)}",
    ]
    for bank in BANKS:
        path = OUTPUT / f"{bank}.idx"
        metadata = inspect(path)
        if not metadata["size_matches"]:
            raise ValueError(f"modified index failed structural validation: {bank}")
        base = BASE / f"{bank}.idx"
        base_raw = base.read_bytes()
        overlay_raw = path.read_bytes()
        observed = [i for i, (old, new) in enumerate(zip(base_raw, overlay_raw)) if old != new]
        if observed != sorted(changed_offsets[bank]):
            raise ValueError(f"unexpected changed bytes in {bank}: {observed!r}")
        manifest.append(
            f"bank={bank} changed_bytes={len(observed)} "
            f"base_sha256={sha256(base_raw)} overlay_sha256={sha256(overlay_raw)}"
        )
    manifest.extend(changes)
    (OUTPUT / "context3-pool-bit7-manifest.txt").write_text(
        "\n".join(manifest) + "\n"
    )
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
