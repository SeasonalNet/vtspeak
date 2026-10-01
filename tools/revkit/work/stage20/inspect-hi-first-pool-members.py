#!/usr/bin/env python3
"""Summarize legacy signatures in the captured 2013 natural-Hi first pool."""

from __future__ import annotations

import re
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))
from inspect_legacy_unit_idx import parse_index_data

TRACE = (
    ROOT
    / "tools"
    / "revkit"
    / "work"
    / "corpus-parity"
    / "stage20"
    / "hi-plain-2013-tailmap-cutoff3-featureb-alias"
    / "adapted-forced-pah0-winedbg.log"
)
ENGINE = ROOT / "binary" / "vt_pau.dll"
IMAGE_BASE = 0x10000000
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)
BANK_BASES = {
    "gen": 0,
    "gen2": 179_995,
    "num": 278_128,
    "etc": 279_125,
    "alp": 282_791,
}
OLD_SLOT0 = {
    272_820,
    22_336,
    65_331,
    75_811,
    135_441,
    164_116,
    205_125,
    219_724,
    255_985,
    277_527,
}
SOURCE = ROOT / "tools" / "revkit" / "work" / "stage19" / "data-kate-copy" / "M16" / "mc_idx_tbl"


def load_rows() -> dict[int, tuple[int, bytes, int]]:
    rows: dict[int, tuple[int, bytes, int]] = {}
    for bank, base in BANK_BASES.items():
        path = SOURCE / f"unit-{bank}.idx"
        raw = path.read_bytes()
        parsed = parse_index_data(path, raw)
        count = parsed.unit_count
        end = parsed.block_start + count * parsed.record_stride
        attr48 = raw[end + count : end + 2 * count]
        keys = raw[end + 2 * count : end + 7 * count]
        attr40 = raw[end + 7 * count : end + 8 * count]
        for local in range(count):
            rows[base + local] = (
                attr48[local],
                bytes(keys[local * 5 : local * 5 + 5]),
                attr40[local],
            )
    return rows


def main() -> None:
    trace = TRACE.read_text(encoding="utf-8", errors="replace")
    match = re.search(
        r"HI_2013_CANDIDATE_LIST position=0 count=\d+(.*?)(?=HI_2013_CANDIDATE_LIST|HI_2013_SELECTED_UNIT)",
        trace,
        re.DOTALL,
    )
    if match is None:
        raise SystemExit(f"no position-0 candidate list in {TRACE}")
    ids = [int(value) for value in re.findall(r"\[\d+\]=(\d+)\(", match.group(1))]
    engine = ENGINE.read_bytes()
    tables = tuple(
        engine[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    rows = load_rows()
    if any(unit_id not in rows for unit_id in ids):
        raise SystemExit("candidate list contains an ID absent from the parsed Kate indexes")

    print(f"candidate_count={len(ids)} old_slot0_members={sum(i in OLD_SLOT0 for i in ids)}")
    print(f"old_members={','.join(str(i) for i in ids if i in OLD_SLOT0)}")
    print(f"extra_members={','.join(str(i) for i in ids if i not in OLD_SLOT0)}")
    for label, key_fn in (
        ("attr48", lambda row: f"{row[0]:02x}"),
        ("attr40", lambda row: f"{row[2]:02x}"),
        ("legacy_key5", lambda row: row[1].hex()),
    ):
        groups = Counter((key_fn(rows[i]), i in OLD_SLOT0) for i in ids)
        print(f"[{label}]")
        for (value, old_member), count in sorted(groups.items()):
            print(f"{value}\t{'old' if old_member else 'extra'}\t{count}")

    query_keys = Counter()
    for unit_id in ids:
        _attr48, key, attr40 = rows[unit_id]
        prefix = bytes((tables[0][key[0]], tables[1][key[1]], tables[2][key[2]]))
        suffix = bytes((0xA0, 0x00)) if prefix == bytes.fromhex("5a2201") and key[3:] == bytes.fromhex("1e1e") else key[3:]
        signature = bytes((rows[unit_id][0],)) + key[:3] + bytes((attr40,)) + suffix
        query_key = bytes((tables[0][signature[1]], tables[1][signature[2]], tables[2][signature[3]], signature[5], signature[6] & 0x20))
        query_keys[query_key.hex()] += 1
    print(f"mapped_query_keys={dict(query_keys)}")

    prefixes = (
        "5a2202",
        "5a2203",
        "5a220e",
        "5a220f",
        "5b2202",
        "5b220e",
        "5b2212",
    )
    print("[candidate_filter_scope]")
    for prefix in prefixes:
        matching = [
            (unit_id, rows[unit_id])
            for unit_id in rows
            if rows[unit_id][1][:3].hex() == prefix
            and rows[unit_id][1][3:] == bytes.fromhex("1e1e")
        ]
        by_bank: Counter[str] = Counter()
        for unit_id, _row in matching:
            for bank, start in reversed(tuple(BANK_BASES.items())):
                if unit_id >= start:
                    by_bank[bank] += 1
                    break
        in_pool = sum(unit_id in ids for unit_id, _row in matching)
        print(
            f"{prefix}\ttotal={len(matching)}\tstage20_pool={in_pool}"
            f"\tbanks={dict(by_bank)}"
        )


if __name__ == "__main__":
    main()
