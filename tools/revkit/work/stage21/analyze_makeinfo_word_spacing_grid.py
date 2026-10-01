#!/usr/bin/env python3
"""Validate and compare Stage 21 punctuation-spacing captures."""

from __future__ import annotations

from collections import Counter
from pathlib import Path
import re


ROOT = Path(__file__).resolve().parent
WORDS = (
    "a", "i", "hello", "world", "hi", "kate", "paul", "good",
    "morning", "weather", "today", "voice",
)
MARKS = ("comma", "period", "ellipsis")
MASKS = ("UU", "UL", "LU", "LL")
LAYOUTS = ("after_space", "both_space", "before_only", "adjacent")
CALL_COUNT = len(WORDS) ** 2 * len(MARKS) * len(MASKS) * len(LAYOUTS)
CALL_PATTERN = re.compile(
    r"^MAKEINFO_WORD_SPACE layout=(\S+) mark=(\S+) mask=(\S+) "
    r"left=(\S+) right=(\S+) raw_eax=(0x[0-9a-f]+)$",
    re.MULTILINE,
)


def main() -> None:
    log = (ROOT / "makeinfo-word-spacing-grid-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    observed = CALL_PATTERN.findall(log)
    if len(observed) != CALL_COUNT:
        raise ValueError(f"expected {CALL_COUNT} calls, found {len(observed)}")

    calls: set[tuple[str, str, str, str, str]] = set()
    sizes: dict[tuple[str, str, str, str, str], str] = {}
    ascii_data: dict[tuple[str, str, str, str, str], bytes] = {}
    binary_data: dict[tuple[str, str, str, str, str], bytes] = {}
    phone_rows: dict[tuple[str, str, str, str, str], tuple[str, ...]] = {}
    size_counts: Counter[tuple[str, str, str, str]] = Counter()
    for layout, mark, mask, left, right, raw_eax in observed:
        key = (layout, mark, mask, left, right)
        if key in calls:
            raise ValueError(f"duplicate call: {key}")
        calls.add(key)
        if layout not in LAYOUTS or mark not in MARKS or mask not in MASKS:
            raise ValueError(f"unexpected call coordinates: {key}")
        if left not in WORDS or right not in WORDS:
            raise ValueError(f"unexpected word pair: {key}")
        if int(raw_eax, 16) != 1:
            raise ValueError(f"{key} returned raw EAX {raw_eax}")

        stem = f"mi-word-space-{layout}-{mark}-{mask}-{left}-{right}"
        ascii_path = ROOT / f"{stem}.asc.dtt"
        binary_path = ROOT / f"{stem}.bin.dtt"
        if not ascii_path.is_file() or not binary_path.is_file():
            raise FileNotFoundError(f"missing paired capture: {key}")
        data = ascii_path.read_bytes()
        binary = binary_path.read_bytes()
        blocks = data.decode("ascii").replace("\r", "").strip().split("\n\n")
        rows: list[tuple[str, str]] = []
        phones: list[str] = []
        for block in blocks:
            flag = re.search(r"^TypeFlag : (\d+)$", block, re.MULTILINE)
            if flag:
                size = re.search(r"^Size : (\d+)$", block, re.MULTILINE)
                rows.append((flag.group(1), size.group(1) if size else ""))
                if flag.group(1) == "2":
                    phones.append(block.strip())
        if not any(flag == "2" for flag, _ in rows):
            raise ValueError(f"no TypeFlag=2 phone rows: {key}")
        silence = [(index, row_size) for index, (flag, row_size) in enumerate(rows) if flag == "1"]
        if len(silence) > 1:
            raise ValueError(f"multiple silence rows: {key}")
        if silence:
            index, row_size = silence[0]
            if index == 0 or index + 1 == len(rows):
                raise ValueError(f"silence row is not between phone rows: {key}")
            if rows[index - 1][0] != "2" or rows[index + 1][0] != "2":
                raise ValueError(f"silence row is not between phone rows: {key}")
            if row_size not in {"3200", "14800"}:
                raise ValueError(f"unexpected TypeFlag=1 size {row_size}: {key}")
        else:
            row_size = "absent"
        sizes[key] = row_size
        ascii_data[key] = data
        binary_data[key] = binary
        phone_rows[key] = tuple(phones)
        size_counts[(layout, mark, mask, row_size)] += 1

    print(f"calls={len(calls)} paired_captures={len(calls) * 2} raw_eax=1")
    for layout in LAYOUTS:
        print(f"\n{layout}")
        for mark in MARKS:
            results = []
            for mask in MASKS:
                parts = ",".join(
                    f"{size}={size_counts[(layout, mark, mask, size)]}"
                    for size in ("3200", "14800", "absent")
                )
                results.append(f"{mask}:{parts}")
            print(f"  {mark}: " + " | ".join(results))

    for mark in MARKS:
        for layout in LAYOUTS[1:]:
            presence_same = 0
            silence_same = 0
            phone_rows_same = 0
            full_ascii_same = 0
            full_binary_same = 0
            for mask in MASKS:
                for left in WORDS:
                    for right in WORDS:
                        baseline = ("after_space", mark, mask, left, right)
                        changed = (layout, mark, mask, left, right)
                        presence_same += (sizes[baseline] != "absent") == (
                            sizes[changed] != "absent"
                        )
                        silence_same += sizes[baseline] == sizes[changed]
                        phone_rows_same += phone_rows[baseline] == phone_rows[changed]
                        full_ascii_same += ascii_data[baseline] == ascii_data[changed]
                        full_binary_same += binary_data[baseline] == binary_data[changed]
            print(
                f"{layout} vs attached {mark}: presence={presence_same}/576 "
                f"silence-size-and-presence={silence_same}/576 "
                f"phone-rows={phone_rows_same}/576 "
                f"full-ASCII={full_ascii_same}/576 full-BIN={full_binary_same}/576"
            )

        for left_layout, right_layout in (
            ("both_space", "before_only"),
            ("after_space", "adjacent"),
        ):
            ascii_same = 0
            binary_same = 0
            for mask in MASKS:
                for left in WORDS:
                    for right in WORDS:
                        left_key = (left_layout, mark, mask, left, right)
                        right_key = (right_layout, mark, mask, left, right)
                        ascii_same += ascii_data[left_key] == ascii_data[right_key]
                        binary_same += binary_data[left_key] == binary_data[right_key]
            print(
                f"{left_layout} vs {right_layout} {mark}: "
                f"ASCII={ascii_same}/576 BIN={binary_same}/576"
            )

    for mark in MARKS:
        for mask in MASKS:
            for left in WORDS:
                for right in WORDS:
                    current = ("after_space", mark, mask, left, right)
                    previous = ROOT / f"mi-word-interword-{mark}-{mask}-{left}-{right}"
                    if ascii_data[current] != (previous.with_suffix(".asc.dtt")).read_bytes():
                        raise ValueError(f"attached ASCII baseline mismatch: {current}")
                    if binary_data[current] != (previous.with_suffix(".bin.dtt")).read_bytes():
                        raise ValueError(f"attached binary baseline mismatch: {current}")
    print("after-space captures byte-match both prior canonical grids")


if __name__ == "__main__":
    main()
