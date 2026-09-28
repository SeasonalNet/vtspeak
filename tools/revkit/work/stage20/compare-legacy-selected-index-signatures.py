#!/usr/bin/env python3
"""Compare old-engine Hello unit IDs with adapted and rebuilt index rows."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data
from inspect_unit_idx import inspect

GLOBAL_GEN2_START = 179_995
SELECTED_GLOBAL_IDS = (272_822, 272_823, 272_824, 272_825)
LEGACY_PATH = ROOT / "data-kate/M16/mc_idx_tbl/unit-gen2.idx"
UNREPACKED_PATH = (
    ROOT / "tools/revkit/work/stage19/index-adapter/unit-gen2.idx"
)
REPACKED_PATH = (
    ROOT
    / "tools/revkit/work/stage20/index-adapter-key-repacked-attr48-key3-attr40-key2"
    / "unit-gen2.idx"
)


def main() -> None:
    legacy = parse_index_data(LEGACY_PATH, LEGACY_PATH.read_bytes())
    legacy_raw = LEGACY_PATH.read_bytes()
    old_tail = legacy_raw[legacy.block_start + legacy.unit_count * legacy.record_stride :]

    versions = []
    for path in (UNREPACKED_PATH, REPACKED_PATH):
        parsed = inspect(path)
        versions.append((path, path.read_bytes(), int(parsed["opaque_block_start"])))

    print(
        "global_id\tlocal_row\tlegacy_attr4c\tlegacy_attr48\tlegacy_key5\t"
        "legacy_attr40\trebuilt_signature\trebuilt_matches\tunrepacked_signature"
    )
    for global_id in SELECTED_GLOBAL_IDS:
        local_row = global_id - GLOBAL_GEN2_START
        row_offset = local_row * legacy.record_stride
        record = legacy.records[row_offset : row_offset + legacy.record_stride]
        attr_4c = old_tail[local_row]
        attr_48 = old_tail[legacy.unit_count + local_row]
        key_start = 2 * legacy.unit_count + local_row * 5
        key_5 = old_tail[key_start : key_start + 5]
        attr_40 = old_tail[7 * legacy.unit_count + local_row]
        expected = bytes((attr_48,)) + key_5[:3] + bytes((attr_40,)) + key_5[3:]

        signatures: list[bytes] = []
        for path, raw, block_start in versions:
            record_start = block_start + row_offset
            if raw[record_start : record_start + legacy.record_stride] != record:
                raise ValueError(f"{path}: payload record differs for global ID {global_id}")
            tail_start = block_start + legacy.unit_count * legacy.record_stride
            signatures.append(
                raw[tail_start + legacy.unit_count + local_row * 7 :
                    tail_start + legacy.unit_count + (local_row + 1) * 7]
            )

        print(
            f"{global_id}\t{local_row}\t{attr_4c:02x}\t{attr_48:02x}\t"
            f"{key_5.hex(' ')}\t{attr_40:02x}\t{expected.hex(' ')}\t"
            f"{str(signatures[1] == expected).lower()}\t{signatures[0].hex(' ')}"
        )


if __name__ == "__main__":
    main()
