#!/usr/bin/env python3
"""Verify VTML x-cmu to internal phone-byte observations from Wine traces."""

from __future__ import annotations

import re
import runpy
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / "stage10"
PROBES = runpy.run_path(
    str(Path(__file__).resolve().parents[2] / "scripts" / "generate_cmu_probes.py")
)
CONSONANTS = PROBES["CONSONANTS"]
VOWELS = PROBES["VOWELS"]

CONTEXT = re.compile(r"UPSTREAM_CONTEXT_BYTES row=(\d+) bytes:\n[^\n]*:\s*((?:0x[0-9a-f]{2}\s*)+)")


def capture(name: str, expected_count: int) -> list[bytes]:
    path = ROOT / f"{name}.log"
    log = path.read_text(errors="replace")
    if "[Inferior 1 (Remote target) exited normally]" not in log:
        raise ValueError(f"{name}: original process did not exit normally")
    match = re.search(r"CONTEXT_ENTRY count=(\d+)", log)
    if match is None or int(match.group(1)) != expected_count:
        raise ValueError(f"{name}: unexpected row count")
    rows: list[bytes] = []
    for match in CONTEXT.finditer(log):
        if int(match.group(1)) != len(rows):
            raise ValueError(f"{name}: out-of-order context row")
        data = bytes(int(value, 16) for value in re.findall(r"0x([0-9a-f]{2})", match.group(2)))
        if not data or 0 not in data:
            raise ValueError(f"{name}: missing NUL terminator in row {len(rows)}")
        rows.append(data[: data.index(0)])
    if len(rows) != expected_count:
        raise ValueError(f"{name}: expected {expected_count} context rows, got {len(rows)}")
    return rows


def main() -> None:
    mapping: dict[str, int] = {}
    for batch in range(3):
        symbols = CONSONANTS[batch * 8 : batch * 8 + 8]
        for symbol, row in zip(symbols, capture(f"cmu-consonants-{batch + 1}", len(symbols))):
            if len(row) != 2 or row[1] != 0x07:
                raise ValueError(f"{symbol}: expected consonant followed by AH0, got {row.hex()}")
            mapping[symbol] = row[0]
    stressed = tuple(f"{symbol}{stress}" for symbol in VOWELS for stress in range(3))
    for batch in range(6):
        symbols = stressed[batch * 8 : batch * 8 + 8]
        for symbol, row in zip(symbols, capture(f"cmu-vowel-stress-{batch + 1}", len(symbols))):
            if len(row) != 2 or row[0] != 0x39:
                raise ValueError(f"{symbol}: expected T followed by vowel, got {row.hex()}")
            mapping[symbol] = row[1]
    if mapping["T"] != 0x39 or mapping["AH0"] != 0x07:
        raise ValueError("anchor phoneme codes differ from initial control")
    if len(mapping) != 69 or set(mapping.values()) != set(range(1, 0x46)):
        raise ValueError("internal code mapping is not a complete one-to-one range 0x01..0x45")
    stress = capture("cmu-stress", 3)
    for index, row in enumerate(stress):
        if row != bytes((mapping["T"], mapping[f"AH{index}"])):
            raise ValueError(f"AH{index}: stress control disagrees with matrix")
    if capture("phoneme-tomato", 1)[0] != bytes(
        mapping[symbol] for symbol in ("T", "AH0", "M", "EY1", "T", "OW0")
    ):
        raise ValueError("vendor-guide tomato example disagrees with codebook")
    if capture("phoneme-trio", 3) != [
        bytes((mapping[symbol], mapping["AH0"])) for symbol in ("P", "T", "K")
    ]:
        raise ValueError("three-tag control disagrees with codebook")
    by_code = {code: symbol for symbol, code in mapping.items()}
    for code in range(1, 0x46):
        print(f"0x{code:02x} {by_code[code]}")
    print("verified 69 unique internal phone bytes, all 0x01..0x45")


if __name__ == "__main__":
    main()
