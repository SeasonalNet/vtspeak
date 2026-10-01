#!/usr/bin/env python3
"""Validate VT, FF, and CR spacing over the complete 12-word pair set."""

from __future__ import annotations

from collections import Counter
from pathlib import Path
import re


ROOT = Path(__file__).resolve().parent
WORDS = ("a", "i", "hello", "world", "hi", "kate", "paul", "good", "morning", "weather", "today", "voice")
PAIRS = tuple((left, right) for left in WORDS for right in WORDS)
MARKS = ("comma", "period", "ellipsis")
MASKS = ("UU", "UL", "LU", "LL")
SIDES = ("before", "after")
SEPARATORS = ("vt", "ff", "cr")
MARK_INDEX = {mark: index for index, mark in enumerate(MARKS)}
MASK_INDEX = {mask: index for index, mask in enumerate(MASKS)}
SIDE_INDEX = {side: index for index, side in enumerate(SIDES)}
SEPARATOR_INDEX = {sep: index for index, sep in enumerate(SEPARATORS)}
EXPECTED_CALLS = len(PAIRS) * len(MARKS) * len(MASKS) * len(SIDES) * len(SEPARATORS)
CALL_PATTERN = re.compile(
    r"^MAKEINFO_ASCII_WSCLASS pair=(\d+) side=(\S+) sep=(\S+) "
    r"mark=(\S+) mask=(\S+) raw_eax=(0x[0-9a-f]+)$",
    re.MULTILINE,
)


def parse_rows(data: bytes, key: tuple[str, ...]) -> tuple[list[tuple[str, str]], tuple[bytes, ...]]:
    blocks = data.decode("ascii").replace("\r", "").strip().split("\n\n")
    rows: list[tuple[str, str]] = []
    phones: list[bytes] = []
    for block in blocks:
        match = re.search(r"^TypeFlag : (\d+)$", block, re.MULTILINE)
        if not match:
            continue
        size = re.search(r"^Size : (\d+)$", block, re.MULTILINE)
        rows.append((match.group(1), size.group(1) if size else ""))
        if match.group(1) == "2":
            phones.append(block.encode("ascii"))
    if not any(flag == "2" for flag, _ in rows):
        raise ValueError(f"no TypeFlag=2 rows: {key}")
    silence = [(index, size) for index, (flag, size) in enumerate(rows) if flag == "1"]
    if len(silence) > 1:
        raise ValueError(f"multiple TypeFlag=1 rows: {key}")
    for index, size in silence:
        if index == 0 or index + 1 >= len(rows):
            raise ValueError(f"silence row is not internal: {key}")
        if rows[index - 1][0] != "2" or rows[index + 1][0] != "2":
            raise ValueError(f"silence row is not between TypeFlag=2 rows: {key}")
        if size not in {"3200", "14800"}:
            raise ValueError(f"unexpected silence size {size}: {key}")
    return rows, tuple(phones)


def main() -> None:
    log_paths = (
        ROOT / "makeinfo-ascii-whitespace-class-grid-all-pairs-api.log",
        ROOT / "makeinfo-ascii-whitespace-class-grid-all-pairs-api-resume-1083.log",
        ROOT / "makeinfo-ascii-whitespace-class-grid-all-pairs-api-resume-9883.log",
    )
    log = "\n".join(path.read_text(encoding="utf-8", errors="replace") for path in log_paths)
    observed = CALL_PATTERN.findall(log)
    if len(observed) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} calls, found {len(observed)}")

    calls: set[tuple[str, str, str, str, str, str]] = set()
    counts: Counter[tuple[str, str, str, str]] = Counter()
    results: dict[tuple[str, str, str, str, str, str], tuple[bytes, bytes, list[tuple[str, str]], tuple[bytes, ...]]] = {}
    for pair_text, side, sep, mark, mask, raw_eax in observed:
        pair_index = int(pair_text)
        if pair_index >= len(PAIRS):
            raise ValueError(f"invalid pair index {pair_index}")
        left, right = PAIRS[pair_index]
        key = (side, sep, mark, mask, left, right)
        if key in calls:
            raise ValueError(f"duplicate call: {key}")
        calls.add(key)
        if side not in SIDES or sep not in SEPARATORS or mark not in MARKS or mask not in MASKS:
            raise ValueError(f"unexpected call coordinates: {key}")
        if int(raw_eax, 16) != 1:
            raise ValueError(f"{key} returned raw EAX {raw_eax}")

        stem = f"mi-word-wsclass12-{side}-{sep}-{mark}-{mask}-{left}-{right}"
        ascii_path = ROOT / f"{stem}.asc.dtt"
        binary_path = ROOT / f"{stem}.bin.dtt"
        if not ascii_path.is_file() or not binary_path.is_file():
            pair_order = pair_index
            case_index = (
                (((pair_order * len(MARKS) + MARK_INDEX[mark]) * len(MASKS)
                   + MASK_INDEX[mask]) * len(SIDES) + SIDE_INDEX[side])
                * len(SEPARATORS) + SEPARATOR_INDEX[sep]
            )
            compact_stem = ROOT / f"mi-ws12-{case_index}"
            ascii_path = compact_stem.with_suffix(".asc.dtt")
            binary_path = compact_stem.with_suffix(".bin.dtt")
        if not ascii_path.is_file() or not binary_path.is_file():
            raise FileNotFoundError(f"missing paired capture: {key}")
        ascii_data = ascii_path.read_bytes()
        binary_data = binary_path.read_bytes()
        rows, phones = parse_rows(ascii_data, key)
        size = next((size for flag, size in rows if flag == "1"), "absent")
        counts[(side, sep, mark, size)] += 1
        results[key] = (ascii_data, binary_data, rows, phones)

    if len(calls) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} distinct calls, found {len(calls)}")

    for side in SIDES:
        for mark in MARKS:
            for mask in MASKS:
                for left, right in PAIRS:
                    keys = [(side, sep, mark, mask, left, right) for sep in SEPARATORS]
                    outputs = [(results[key][0], results[key][1]) for key in keys]
                    if outputs[0] != outputs[1] or outputs[0] != outputs[2]:
                        raise ValueError(f"VT/FF/CR captures differ: {side} {mark} {mask} {left}/{right}")

    print(f"calls={len(calls)} paired_captures={2 * len(calls)} raw_eax=1")
    for side in SIDES:
        for sep in SEPARATORS:
            print(f"\n{side} {sep}: TypeFlag=1 sizes over 576 word/case combinations")
            for mark in MARKS:
                print(
                    f"  {mark}: 3200={counts[(side, sep, mark, '3200')]} "
                    f"14800={counts[(side, sep, mark, '14800')]} "
                    f"absent={counts[(side, sep, mark, 'absent')]}"
                )

            for mark in MARKS:
                same_ascii = same_binary = same_rows = same_phones = 0
                changed = []
                for mask in MASKS:
                    for left, right in PAIRS:
                        key = (side, sep, mark, mask, left, right)
                        base_key = (mark, mask, left, right, "after_space")
                        base_stem = f"mi-word-space-after_space-{mark}-{mask}-{left}-{right}"
                        base_ascii = (ROOT / f"{base_stem}.asc.dtt").read_bytes()
                        base_binary = (ROOT / f"{base_stem}.bin.dtt").read_bytes()
                        base_rows, base_phones = parse_rows(base_ascii, base_key)
                        ascii_data, binary_data, rows, phones = results[key]
                        same_ascii += ascii_data == base_ascii
                        same_binary += binary_data == base_binary
                        same_rows += rows == base_rows
                        same_phones += phones == base_phones
                        if ascii_data != base_ascii:
                            size = next((s for f, s in rows if f == "1"), "absent")
                            base_size = next((s for f, s in base_rows if f == "1"), "absent")
                            changed.append(f"{mask}:{left}/{right}({base_size}->{size})")
                print(
                    f"  {mark} vs canonical layout: rows={same_rows}/576 phones={same_phones}/576 "
                    f"ASCII={same_ascii}/576 BIN={same_binary}/576"
                )
                if changed:
                    print(f"    changed: {', '.join(changed)}")
        for mark in MARKS:
            print(f"{side} {mark}: VT/FF/CR byte-identical=576/576")


if __name__ == "__main__":
    main()
