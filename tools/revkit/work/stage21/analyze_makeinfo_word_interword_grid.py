#!/usr/bin/env python3
"""Validate and summarize the Stage 21 word-level punctuation grid."""

from __future__ import annotations

from collections import Counter, defaultdict
from pathlib import Path
import re


ROOT = Path(__file__).resolve().parent
WORDS = (
    "a", "i", "hello", "world", "hi", "kate", "paul", "good",
    "morning", "weather", "today", "voice",
)
MARKS = ("comma", "period", "ellipsis")
MASKS = ("UU", "UL", "LU", "LL")
EXPECTED_CALLS = len(WORDS) ** 2 * len(MARKS) * len(MASKS)


def main() -> None:
    log = (ROOT / "makeinfo-word-interword-grid-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    observed = re.findall(
        r"^MAKEINFO_WORD_INTERWORD mark=(\S+) mask=(\S+) left=(\S+) "
        r"right=(\S+) raw_eax=(0x[0-9a-f]+)$",
        log,
        re.MULTILINE,
    )
    if len(observed) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} calls, found {len(observed)}")

    calls: set[tuple[str, str, str, str]] = set()
    presence: dict[tuple[str, str, str], bool] = {}
    size_counts: Counter[tuple[str, str, str]] = Counter()
    for mark, mask, left, right, raw_eax in observed:
        key = (mark, mask, left, right)
        if key in calls:
            raise ValueError(f"duplicate API call record: {key}")
        calls.add(key)
        if mark not in MARKS or mask not in MASKS or left not in WORDS or right not in WORDS:
            raise ValueError(f"unexpected call coordinates: {key}")
        if int(raw_eax, 16) != 1:
            raise ValueError(f"{key} returned raw EAX {raw_eax}")

        stem = f"mi-word-interword-{mark}-{mask}-{left}-{right}"
        ascii_path = ROOT / f"{stem}.asc.dtt"
        binary_path = ROOT / f"{stem}.bin.dtt"
        if not ascii_path.is_file() or not binary_path.is_file():
            raise FileNotFoundError(f"missing paired DTT capture for {key}")
        blocks = ascii_path.read_text(encoding="ascii").strip().split("\n\n")
        rows = []
        for block in blocks:
            match = re.search(r"^TypeFlag : (\d+)$", block, re.MULTILINE)
            if match:
                rows.append((block, match.group(1)))
        row_flags = [flag for _, flag in rows]
        if "2" not in row_flags:
            raise ValueError(f"{key} contains no TypeFlag=2 row")
        silence_rows = [
            (index, re.search(r"^Size : (\d+)$", block, re.MULTILINE))
            for index, (block, flag) in enumerate(rows)
            if flag == "1"
        ]
        if len(silence_rows) > 1:
            raise ValueError(f"{key} contains multiple TypeFlag=1 rows")
        if silence_rows:
            index, size_match = silence_rows[0]
            if index == 0 or index + 1 >= len(rows):
                raise ValueError(f"{key} silence row is not between phone rows")
            if row_flags[index - 1] != "2" or row_flags[index + 1] != "2":
                raise ValueError(f"{key} silence row is not between TypeFlag=2 rows")
            if size_match is None or size_match.group(1) not in {"3200", "14800"}:
                raise ValueError(f"{key} has an unexpected silence size")
            size = size_match.group(1)
        else:
            size = "absent"
        presence[key] = size != "absent"
        size_counts[(mark, mask, size)] += 1

    print(f"calls={len(calls)} paired_captures={len(calls) * 2} raw_eax=1")
    for mark in MARKS:
        for mask in MASKS:
            present = sum(
                size_counts[(mark, mask, size)] for size in ("3200", "14800")
            )
            print(
                f"{mark} {mask}: silence={present}/{len(WORDS) ** 2} "
                f"sizes=3200:{size_counts[(mark, mask, '3200')]} "
                f"14800:{size_counts[(mark, mask, '14800')]} "
                f"absent:{size_counts[(mark, mask, 'absent')]}"
            )

    for mark in MARKS:
        signatures: Counter[str] = Counter()
        examples: dict[str, list[str]] = defaultdict(list)
        for left in WORDS:
            for right in WORDS:
                signature = "".join(
                    "1" if presence[(mark, mask, left, right)] else "0"
                    for mask in MASKS
                )
                signatures[signature] += 1
                if signature != "1111":
                    examples[signature].append(f"{left}/{right}")
        print(f"{mark} case signatures (UU,UL,LU,LL): {dict(sorted(signatures.items()))}")
        for signature, pairs in sorted(examples.items()):
            print(f"  {signature}: {', '.join(pairs)}")


if __name__ == "__main__":
    main()
