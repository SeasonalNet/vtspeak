#!/usr/bin/env python3
"""Validate all-pair MakeInfo punctuation whitespace captures."""

from __future__ import annotations

from collections import Counter
from pathlib import Path
import re


ROOT = Path(__file__).resolve().parent
WORDS = ("a", "i", "hello", "world", "hi", "kate", "paul", "good", "morning", "weather", "today", "voice")
PAIRS = tuple((left, right) for left in WORDS for right in WORDS)
MARKS = ("comma", "period", "ellipsis")
MASKS = ("UU", "UL", "LU", "LL")
SEPARATORS = ("none", "space", "double", "tab", "lf", "crlf")
EXPECTED_CALLS = len(PAIRS) * len(MARKS) * len(MASKS) * len(SEPARATORS) ** 2
MARK_INDEX = {value: index for index, value in enumerate(MARKS)}
MASK_INDEX = {value: index for index, value in enumerate(MASKS)}
SEPARATOR_INDEX = {value: index for index, value in enumerate(SEPARATORS)}
CALL_PATTERN = re.compile(
    r"^MAKEINFO_ASCII_WS_FULL pair=(\d+) mark=(\S+) mask=(\S+) "
    r"pre=(\S+) post=(\S+) raw_eax=(0x[0-9a-f]+)$",
    re.MULTILINE,
)


def parse_rows(data: bytes, key: tuple[str, ...]) -> list[tuple[str, str]]:
    blocks = data.decode("ascii").replace("\r", "").strip().split("\n\n")
    rows: list[tuple[str, str]] = []
    for block in blocks:
        match = re.search(r"^TypeFlag : (\d+)$", block, re.MULTILINE)
        if not match:
            continue
        size = re.search(r"^Size : (\d+)$", block, re.MULTILINE)
        rows.append((match.group(1), size.group(1) if size else ""))
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
    return rows


def case_index(pair: int, mark: str, mask: str, pre: str, post: str) -> int:
    return (
        (((pair * len(MARKS) + MARK_INDEX[mark]) * len(MASKS) + MASK_INDEX[mask])
         * len(SEPARATORS) + SEPARATOR_INDEX[pre])
        * len(SEPARATORS) + SEPARATOR_INDEX[post]
    )


def main() -> None:
    log_paths = sorted(ROOT.glob("makeinfo-ascii-whitespace-all-pairs-api*.log"))
    if not log_paths:
        raise FileNotFoundError("no all-pairs whitespace API logs")
    log = "\n".join(path.read_text(encoding="utf-8", errors="replace") for path in log_paths)
    observed = CALL_PATTERN.findall(log)
    if len(observed) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} calls, found {len(observed)}")

    calls: set[tuple[str, str, str, str, str, str]] = set()
    results: dict[tuple[str, str, str, str, str, str], tuple[bytes, bytes, list[tuple[str, str]]]] = {}
    silence_counts: Counter[tuple[str, str, str, str]] = Counter()
    for pair_text, mark, mask, pre, post, raw_eax in observed:
        pair = int(pair_text)
        if pair >= len(PAIRS):
            raise ValueError(f"invalid pair index {pair}")
        if mark not in MARKS or mask not in MASKS or pre not in SEPARATORS or post not in SEPARATORS:
            raise ValueError(f"invalid call coordinate: {pair_text} {mark} {mask} {pre} {post}")
        if int(raw_eax, 16) != 1:
            raise ValueError(f"raw EAX {raw_eax}: {pair_text} {mark} {mask} {pre} {post}")
        left, right = PAIRS[pair]
        key = (mark, mask, left, right, pre, post)
        if key in calls:
            raise ValueError(f"duplicate call: {key}")
        calls.add(key)
        index = case_index(pair, mark, mask, pre, post)
        stem = ROOT / f"mi-wsfull12-{index}"
        ascii_path = stem.with_suffix(".asc.dtt")
        binary_path = stem.with_suffix(".bin.dtt")
        if not ascii_path.is_file() or not binary_path.is_file():
            raise FileNotFoundError(f"missing paired capture: {key}")
        ascii_data = ascii_path.read_bytes()
        binary_data = binary_path.read_bytes()
        rows = parse_rows(ascii_data, key)
        size = next((value for flag, value in rows if flag == "1"), "absent")
        silence_counts[(mark, pre, post, size)] += 1
        results[key] = (ascii_data, binary_data, rows)

    if len(calls) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} distinct calls, found {len(calls)}")

    # Ensure the four pre-existing all-pair layouts are represented exactly.
    controls = {
        ("none", "space"): "after_space",
        ("space", "space"): "both_space",
        ("space", "none"): "before_only",
        ("none", "none"): "adjacent",
    }
    control_counts: Counter[tuple[str, str]] = Counter()
    for (pre, post), layout in controls.items():
        for mark in MARKS:
            for mask in MASKS:
                for left, right in PAIRS:
                    key = (mark, mask, left, right, pre, post)
                    stem = ROOT / f"mi-word-space-{layout}-{mark}-{mask}-{left}-{right}"
                    expected = (stem.with_suffix(".asc.dtt").read_bytes(), stem.with_suffix(".bin.dtt").read_bytes())
                    actual = results[key][:2]
                    control_counts[(layout, "same")] += actual == expected

    print(f"calls={len(calls)} paired_captures={2 * len(calls)} raw_eax=1")
    for layout in controls.values():
        count = control_counts[(layout, "same")]
        expected_count = len(MARKS) * len(MASKS) * len(PAIRS)
        if count != expected_count:
            raise ValueError(
                f"layout {layout} matches its existing control in {count}/{expected_count} cases"
            )
        print(f"existing layout {layout}: ASCII+BIN identical={count}/{expected_count}")

    layouts = tuple((pre, post) for pre in SEPARATORS for post in SEPARATORS)
    for mark in MARKS:
        groups: list[list[tuple[str, str]]] = []
        for layout in layouts:
            for group in groups:
                representative = group[0]
                same = all(
                    results[(mark, mask, left, right, *layout)][:2]
                    == results[(mark, mask, left, right, *representative)][:2]
                    for mask in MASKS
                    for left, right in PAIRS
                )
                if same:
                    group.append(layout)
                    break
            else:
                groups.append([layout])
        rendered = [
            "=".join(f"{pre}/{post}" for pre, post in group)
            for group in groups
        ]
        print(f"{mark} globally identical layout groups: " + " | ".join(rendered))

    for mark in MARKS:
        for pre in SEPARATORS:
            for post in SEPARATORS:
                counts = {
                    size: silence_counts[(mark, pre, post, size)]
                    for size in ("3200", "14800", "absent")
                }
                exact = 0
                rows_same = 0
                for mask in MASKS:
                    for left, right in PAIRS:
                        key = (mark, mask, left, right, pre, post)
                        base = (mark, mask, left, right, "none", "space")
                        exact += results[key][:2] == results[base][:2]
                        rows_same += results[key][2] == results[base][2]
                print(
                    f"{mark} pre={pre} post={post}: silence="
                    f"3200/{counts['3200']} 14800/{counts['14800']} absent/{counts['absent']}; "
                    f"vs after_space rows={rows_same}/576 ASCII+BIN={exact}/576"
                )


if __name__ == "__main__":
    main()
