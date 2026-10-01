#!/usr/bin/env python3
"""Map the Bridget A/B selected IDs to their source and converted index fields."""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data
from inspect_unit_idx import inspect

WORK = Path(__file__).resolve().parent
SOURCE_ROOT = ROOT / "data-bridget" / "M16" / "mc_idx_tbl"
STRUCTURAL_ROOT = WORK / "bridget-2005-structural-adapter" / "M16" / "mc_idx_tbl"
SIGNATURE_ROOT = WORK / "bridget-2005-signature-adapter" / "M16" / "mc_idx_tbl"
BANKS = (
    ("gen", 0, 723_614),
    ("alp", 723_614, 1_130),
    ("etc", 724_744, 47_249),
    ("exp", 771_993, 8_889),
)
SELECTED = {
    "signature": (718258, 106109, 260196, 576411, 39073),
    "structural": (178150, 243779, 532206, 682208, 272930, 733259),
}
IMAGE_BASE = 0x10000000
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)


def columns(path: Path) -> tuple[int, bytes]:
    raw = path.read_bytes()
    parsed = parse_index_data(path, raw)
    count = parsed.unit_count
    start = parsed.block_start + count * parsed.record_stride
    return count, raw[start:]


def signature_column(path: Path) -> bytes:
    parsed = inspect(path)
    raw = path.read_bytes()
    count = int(parsed["unit_count"])
    start = int(parsed["opaque_block_start"]) + count * int(
        parsed["opaque_block_stride_bytes"]
    )
    start += count
    return raw[start : start + count * 7]


def main() -> None:
    image = (ROOT / "binary" / "vt_kat.dll").read_bytes()
    tables = tuple(
        image[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    if any(len(table) != 256 for table in tables):
        raise ValueError("2013 key table extends beyond vt_kat.dll")

    print(
        "layout\tglobal_id\tbank_row\tattr48\tlegacy_key5\tattr40\t"
        "signature7\tprojected_2013_key5"
    )
    legacy_cache: dict[str, tuple[int, bytes]] = {}
    signature_cache: dict[tuple[str, str], bytes] = {}
    for layout, ids in SELECTED.items():
        for unit in ids:
            bank, base, count_expected = next(
                (name, base, count)
                for name, base, count in BANKS
                if base <= unit < base + count
            )
            local = unit - base
            if bank not in legacy_cache:
                legacy_cache[bank] = columns(SOURCE_ROOT / f"unit-{bank}.idx")
            actual_count, legacy_tail = legacy_cache[bank]
            count = count_expected
            if actual_count != count:
                raise ValueError(f"unexpected {bank} row count: {actual_count}")
            if len(legacy_tail) != count * 20:
                raise ValueError(f"unexpected {bank} legacy tail length")
            attr48 = legacy_tail[count : 2 * count]
            keys = legacy_tail[2 * count : 7 * count]
            attr40 = legacy_tail[7 * count : 8 * count]
            adapter_root = STRUCTURAL_ROOT if layout == "structural" else SIGNATURE_ROOT
            cache_key = (layout, bank)
            if cache_key not in signature_cache:
                signature_cache[cache_key] = signature_column(
                    adapter_root / f"unit-{bank}.idx"
                )
            signatures = signature_cache[cache_key]
            key = keys[local * 5 : local * 5 + 5]
            sig = signatures[local * 7 : local * 7 + 7]
            projected = bytes(
                (
                    tables[0][sig[1]],
                    tables[1][sig[2]],
                    tables[2][sig[3]],
                    sig[5],
                    sig[6] & 0x20,
                )
            )
            print(
                f"{layout}\t{unit}\t{bank}:{local}\t{attr48[local]:02x}\t{key.hex()}\t"
                f"{attr40[local]:02x}\t{sig.hex()}\t{projected.hex()}"
            )


if __name__ == "__main__":
    main()
