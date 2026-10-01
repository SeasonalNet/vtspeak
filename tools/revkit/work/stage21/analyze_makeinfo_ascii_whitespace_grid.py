#!/usr/bin/env python3
"""Validate and summarize the Stage 21 ASCII-whitespace punctuation grid."""

from __future__ import annotations

from collections import Counter
from pathlib import Path
import re


ROOT = Path(__file__).resolve().parent
PAIRS = (
    ("a", "a"), ("a", "hello"), ("hello", "world"), ("paul", "hi"),
    ("good", "morning"), ("hi", "kate"), ("a", "good"),
    ("weather", "today"),
)
MARKS = ("comma", "period", "ellipsis")
MASKS = ("UU", "UL", "LU", "LL")
SEPARATORS = ("none", "space", "double", "tab", "lf", "crlf")
EXPECTED_CALLS = len(PAIRS) * len(MARKS) * len(MASKS) * len(SEPARATORS) ** 2
CALL_PATTERN = re.compile(
    r"^MAKEINFO_ASCII_WS pair=(\d+) mark=(\S+) mask=(\S+) "
    r"pre=(\S+) post=(\S+) raw_eax=(0x[0-9a-f]+)$",
    re.MULTILINE,
)


def read_rows(data: bytes, key: tuple[str, ...]) -> tuple[list[tuple[str, str]], tuple[bytes, ...]]:
    blocks = data.decode("ascii").replace("\r", "").strip().split("\n\n")
    rows: list[tuple[str, str]] = []
    phones: list[bytes] = []
    for block in blocks:
        flag = re.search(r"^TypeFlag : (\d+)$", block, re.MULTILINE)
        if not flag:
            continue
        size = re.search(r"^Size : (\d+)$", block, re.MULTILINE)
        rows.append((flag.group(1), size.group(1) if size else ""))
        if flag.group(1) == "2":
            phones.append(block.encode("ascii"))
    if not any(flag == "2" for flag, _ in rows):
        raise ValueError(f"no TypeFlag=2 phone rows: {key}")
    silences = [
        (index, size)
        for index, (flag, size) in enumerate(rows)
        if flag == "1"
    ]
    if len(silences) > 1:
        raise ValueError(f"multiple TypeFlag=1 rows: {key}")
    for index, size in silences:
        if index == 0 or index + 1 >= len(rows):
            raise ValueError(f"silence row is not between phone rows: {key}")
        if rows[index - 1][0] != "2" or rows[index + 1][0] != "2":
            raise ValueError(f"silence row is not between phone rows: {key}")
        if size not in {"3200", "14800"}:
            raise ValueError(f"unexpected TypeFlag=1 size {size}: {key}")
    return rows, tuple(phones)


def main() -> None:
    logs = (
        ROOT / "makeinfo-ascii-whitespace-grid-api.log",
        ROOT / "makeinfo-ascii-whitespace-grid-api-resume-1975.log",
    )
    text = "\n".join(path.read_text(encoding="utf-8", errors="replace") for path in logs)
    observed = CALL_PATTERN.findall(text)
    if len(observed) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} calls, found {len(observed)}")

    rows_by_key: dict[tuple[str, str, str, str, str, str], list[tuple[str, str]]] = {}
    phones_by_key: dict[tuple[str, str, str, str, str, str], tuple[bytes, ...]] = {}
    ascii_by_key: dict[tuple[str, str, str, str, str, str], bytes] = {}
    binary_by_key: dict[tuple[str, str, str, str, str, str], bytes] = {}
    silence_counts: Counter[tuple[str, str, str, str]] = Counter()
    calls: set[tuple[str, str, str, str, str, str]] = set()
    for pair_text, mark, mask, pre, post, raw_eax in observed:
        pair_index = int(pair_text)
        if pair_index >= len(PAIRS):
            raise ValueError(f"unexpected pair index: {pair_index}")
        left, right = PAIRS[pair_index]
        key = (mark, mask, left, right, pre, post)
        if key in calls:
            raise ValueError(f"duplicate call: {key}")
        calls.add(key)
        if mark not in MARKS or mask not in MASKS or pre not in SEPARATORS or post not in SEPARATORS:
            raise ValueError(f"unexpected call coordinates: {key}")
        if int(raw_eax, 16) != 1:
            raise ValueError(f"{key} returned raw EAX {raw_eax}")

        stem = f"mi-word-wsgrid-{mark}-{mask}-{left}-{right}-pre-{pre}-post-{post}"
        ascii_path = ROOT / f"{stem}.asc.dtt"
        binary_path = ROOT / f"{stem}.bin.dtt"
        if not ascii_path.is_file() or not binary_path.is_file():
            raise FileNotFoundError(f"missing paired capture: {key}")
        ascii_data = ascii_path.read_bytes()
        binary_data = binary_path.read_bytes()
        rows, phone_rows = read_rows(ascii_data, key)
        ascii_by_key[key] = ascii_data
        binary_by_key[key] = binary_data
        rows_by_key[key] = rows
        phones_by_key[key] = phone_rows
        size = next((size for flag, size in rows if flag == "1"), "absent")
        silence_counts[(mark, pre, post, size)] += 1

    if len(calls) != EXPECTED_CALLS:
        raise ValueError(f"expected {EXPECTED_CALLS} distinct calls, found {len(calls)}")
    print(f"calls={len(calls)} paired_captures={len(calls) * 2} raw_eax=1")
    for mark in MARKS:
        print(f"\n{mark}: TypeFlag=1 rows by pre/post separator; each cell is count / 32")
        for pre in SEPARATORS:
            print(
                "  " + pre + ": " + " | ".join(
                    f"{post} 3200={silence_counts[(mark, pre, post, '3200')]} "
                    f"14800={silence_counts[(mark, pre, post, '14800')]} "
                    f"absent={silence_counts[(mark, pre, post, 'absent')]}"
                    for post in SEPARATORS
                )
            )

        baseline = ("none", "space")
        for pre in SEPARATORS:
            for post in SEPARATORS:
                same_full_ascii = same_full_binary = same_phones = same_rows = 0
                for mask in MASKS:
                    for left, right in PAIRS:
                        current = (mark, mask, left, right, pre, post)
                        control = (mark, mask, left, right, *baseline)
                        same_full_ascii += ascii_by_key[current] == ascii_by_key[control]
                        same_full_binary += binary_by_key[current] == binary_by_key[control]
                        same_phones += phones_by_key[current] == phones_by_key[control]
                        same_rows += rows_by_key[current] == rows_by_key[control]
                if (pre, post) != baseline:
                    print(
                        f"  {pre}/{post} vs no-pre/single-space: rows={same_rows}/32 "
                        f"phones={same_phones}/32 ASCII={same_full_ascii}/32 "
                        f"BIN={same_full_binary}/32"
                    )

    print("\nPer-separator changes from no separator before punctuation and one space after:")
    for mark in MARKS:
        for pre in SEPARATORS:
            same = sum(
                rows_by_key[(mark, mask, left, right, pre, "space")]
                == rows_by_key[(mark, mask, left, right, "none", "space")]
                for mask in MASKS for left, right in PAIRS
            )
            print(f"  {mark} pre={pre}: identical record rows={same}/32")
        for post in SEPARATORS:
            same = sum(
                rows_by_key[(mark, mask, left, right, "none", post)]
                == rows_by_key[(mark, mask, left, right, "none", "space")]
                for mask in MASKS for left, right in PAIRS
            )
            print(f"  {mark} post={post}: identical record rows={same}/32")

        print(f"  exact equality with post=space while varying pre separator:")
        for left_pre in SEPARATORS:
            equal_to = []
            for right_pre in SEPARATORS:
                if all(
                    ascii_by_key[(mark, mask, left, right, left_pre, "space")]
                    == ascii_by_key[(mark, mask, left, right, right_pre, "space")]
                    and binary_by_key[(mark, mask, left, right, left_pre, "space")]
                    == binary_by_key[(mark, mask, left, right, right_pre, "space")]
                    for mask in MASKS for left, right in PAIRS
                ):
                    equal_to.append(right_pre)
            print(f"    {left_pre}: {','.join(equal_to)}")
        print(f"  exact equality with pre=none while varying post separator:")
        for left_post in SEPARATORS:
            equal_to = []
            for right_post in SEPARATORS:
                if all(
                    ascii_by_key[(mark, mask, left, right, "none", left_post)]
                    == ascii_by_key[(mark, mask, left, right, "none", right_post)]
                    and binary_by_key[(mark, mask, left, right, "none", left_post)]
                    == binary_by_key[(mark, mask, left, right, "none", right_post)]
                    for mask in MASKS for left, right in PAIRS
                ):
                    equal_to.append(right_post)
            print(f"    {left_post}: {','.join(equal_to)}")

    print("\nFocused case differences (only cases with changed ASCII output):")
    comparisons = (
        ("period", "none", "space", "space", "space", "ordinary pre-space effect"),
        ("period", "none", "lf", "space", "space", "LF before mark"),
        ("period", "none", "crlf", "space", "space", "CRLF before mark"),
        ("period", "none", "none", "none", "space", "remove post-space"),
        ("period", "none", "none", "none", "double", "double post-space"),
        ("period", "none", "none", "none", "lf", "LF after mark"),
        ("period", "none", "none", "none", "crlf", "CRLF after mark"),
        ("ellipsis", "none", "space", "space", "space", "space before ellipsis"),
        ("ellipsis", "none", "none", "none", "double", "double after ellipsis"),
        ("ellipsis", "none", "none", "none", "crlf", "CRLF after ellipsis"),
    )
    for mark, left_pre, right_pre, left_post, right_post, label in comparisons:
        changed = []
        for mask in MASKS:
            for left, right in PAIRS:
                left_key = (mark, mask, left, right, left_pre, left_post)
                right_key = (mark, mask, left, right, right_pre, right_post)
                if ascii_by_key[left_key] != ascii_by_key[right_key]:
                    left_size = next((s for f, s in rows_by_key[left_key] if f == "1"), "absent")
                    right_size = next((s for f, s in rows_by_key[right_key] if f == "1"), "absent")
                    changed.append(f"{mask}:{left}/{right}({left_size}->{right_size})")
        print(f"  {label}: {', '.join(changed) if changed else 'none'}")


if __name__ == "__main__":
    main()
