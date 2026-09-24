#!/usr/bin/env python3
"""Summarize controlled CMU-to-tree input values from Stage 10 captures."""

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
CALL = re.compile(
    r"^TREE_(?:SCALAR|VECTOR) tree=(0x[0-9a-f]+) return=.*? features:\n"
    r"([^\n]+)\n([^\n]+)",
    re.MULTILINE,
)
TARGETS = ("0x14bf434", "0x14bf3c4", "0x14bf4f8", "0x14bf568")


def feature_values(line: str) -> tuple[int, ...]:
    _, values = line.split(":", 1)
    return tuple(int(value) for value in values.split())


def capture(name: str) -> list[tuple[str, tuple[int, ...]]]:
    path = ROOT / f"tree-phone-{name}-trees.log"
    log = path.read_text(errors="replace")
    if "[Inferior 1 (Remote target) exited normally]" not in log:
        raise ValueError(f"{path}: process did not exit normally")
    calls = [
        (match.group(1), feature_values(match.group(2)) + feature_values(match.group(3)))
        for match in CALL.finditer(log)
    ]
    addresses = tuple(address for address, _ in calls)
    if len(calls) != 4 or addresses[2:] != TARGETS[2:] or not addresses[1].startswith("0x14bf"):
        raise ValueError(f"{path}: unexpected ordered duration/pitch calls {addresses!r}")
    return calls


def main() -> None:
    symbols = tuple(sorted((*CONSONANTS, *VOWELS)))
    phone_ordinals = {symbol: index for index, symbol in enumerate(symbols, start=1)}
    manner_trees = {
        **dict.fromkeys(("P", "B", "T", "D", "K", "G"), "0x14bf434"),
        **dict.fromkeys(("CH", "JH"), "0x14bf46c"),
        **dict.fromkeys(("F", "V", "TH", "DH", "S", "Z", "SH", "ZH", "HH"), "0x14bf450"),
        **dict.fromkeys(("L", "M", "N", "NG"), "0x14bf488"),
        **dict.fromkeys(("R", "W", "Y"), "0x14bf4a4"),
    }
    print("consonant, phone identity ordinal, selected onset tree")
    consonant_values: dict[str, int] = {}
    for symbol in CONSONANTS:
        calls = capture(f"consonant-{symbol.lower()}")
        first = calls[0][1][0]
        second = calls[1][1][1]
        if first != second or first != phone_ordinals[symbol]:
            raise ValueError(f"{symbol}: first and second tree categories disagree")
        if calls[0][0] != manner_trees[symbol]:
            raise ValueError(f"{symbol}: unexpected selected onset tree {calls[0][0]}")
        consonant_values[symbol] = first
        print(f"{symbol},{first},{calls[0][0]}")
    print("vowel, phone identity ordinal")
    vowel_values: dict[str, int] = {}
    for symbol in VOWELS:
        calls = capture(f"vowel-{symbol.lower()}")
        vowel_values[symbol] = calls[1][1][0]
        if (
            vowel_values[symbol] != phone_ordinals[symbol]
            or calls[0][1][2] != phone_ordinals[symbol]
        ):
            raise ValueError(f"{symbol}: unexpected phone identity ordinal {vowel_values[symbol]}")
        print(f"{symbol},{vowel_values[symbol]}")
    if set(consonant_values.values()) | set(vowel_values.values()) != set(range(1, 40)):
        raise ValueError("phone identity ordinals are not a one-to-one CMU alphabet ranking")
    print(
        f"verified alphabetical CMU identity ordinals for {len(consonant_values)} consonants "
        f"and {len(vowel_values)} vowels, plus five observed onset tree families"
    )


if __name__ == "__main__":
    main()
