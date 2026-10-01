#!/usr/bin/env python3
"""Validate MakeInfo records for exhaustive two-letter interword punctuation probes."""

from __future__ import annotations

import re
import struct
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
PROBE = ROOT / "tools/revkit/work/stage21"
INDEX_DIR = ROOT / "data-paul/M16/mc_idx_tbl"
LETTERS = "abcdefghijklmnopqrstuvwxyz"
MASKS = ("UU", "UL", "LU", "LL")
MARKS = ("comma", "period", "ellipsis")
BANKS = ("merged-gen", "merged-num", "merged-etc", "merged-alp")


def read_fields(block: str) -> dict[str, str]:
    result: dict[str, str] = {}
    for line in block.splitlines():
        if " : " in line:
            key, value = line.split(" : ", 1)
            result[key] = value
    return result


def load_dat_offsets(bank: str) -> set[int]:
    path = INDEX_DIR / f"unit-{bank.removeprefix('merged-')}.idx"
    raw = path.read_bytes()
    if len(raw) < 45:
        raise ValueError(f"{path}: truncated index header")
    count = struct.unpack_from("<I", raw, 39)[0]
    stride = struct.unpack_from("<H", raw, 43)[0]
    if stride != 19:
        raise ValueError(f"{path}: unexpected unit stride {stride}")
    end = 45 + count * stride
    if end > len(raw):
        raise ValueError(f"{path}: index rows extend beyond EOF")
    return {struct.unpack_from("<I", raw, offset)[0] for offset in range(45, end, stride)}


def expected_text(first: str, second: str, mask: str, mark: str) -> str:
    left = first.upper() if mask[0] == "U" else first
    right = second.upper() if mask[1] == "U" else second
    punctuation = {"comma": ",", "period": ".", "ellipsis": "..."}[mark]
    return f"{left}{punctuation} {right}"


def main() -> int:
    selected_marks = tuple(sys.argv[1:]) or MARKS
    if any(mark not in MARKS for mark in selected_marks):
        raise SystemExit("usage: analyze_makeinfo_interword_punctuation_matrix.py [comma|period|ellipsis ...]")
    offsets = [load_dat_offsets(bank) for bank in BANKS]
    totals = {mark: {"calls": 0, "type1": 0, "type2": 0, "captures_with_type1": 0} for mark in selected_marks}
    mismatches: list[str] = []
    rows_by_case: dict[tuple[str, str, str], list[dict[str, str]]] = {}

    for mark in selected_marks:
        log_path = PROBE / f"makeinfo-interword-punctuation-{mark}-api.log"
        log = log_path.read_text(encoding="ascii")
        matches = re.findall(
            rf"MAKEINFO_INTERWORD_PUNCT mark={mark} case=([^ ]+) "
            r"raw_eax=(0x[0-9a-f]+) text=(.*?) path=([^\r\n]+)",
            log,
        )
        expected_calls = 26 * 26 * len(MASKS)
        if len(matches) != expected_calls:
            raise ValueError(f"{mark}: expected {expected_calls} calls, found {len(matches)}")
        by_case = {case: (eax, text, path) for case, eax, text, path in matches}
        if len(by_case) != expected_calls:
            raise ValueError(f"{mark}: duplicate API case labels")
        if any(eax != "0x1" for eax, _, _ in by_case.values()):
            raise ValueError(f"{mark}: one or more calls did not return raw EAX 1")

        mark_rows = totals[mark]
        mark_rows["calls"] = len(matches)
        for first in LETTERS:
            for second in LETTERS:
                pair = first + second
                for mask in MASKS:
                    case = f"{pair}-{mask}"
                    _, text, path = by_case[case]
                    stem = f"mi-interword-punct-{mark}-{mask}-{pair}"
                    expected_path = f"Z:/work/stage21/{stem}"
                    if text != expected_text(first, second, mask, mark):
                        raise ValueError(f"{mark}/{case}: unexpected input text {text!r}")
                    if path != expected_path:
                        raise ValueError(f"{mark}/{case}: unexpected output prefix {path!r}")
                    ascii_path = PROBE / f"{stem}.asc.dtt"
                    binary_path = PROBE / f"{stem}.bin.dtt"
                    if not ascii_path.is_file() or not binary_path.is_file():
                        raise ValueError(f"{mark}/{case}: missing paired DTT output")
                    blocks = ascii_path.read_text(encoding="ascii").split("\n\n")
                    rows = [row for block in blocks if (row := read_fields(block))]
                    rows_by_case[(mark, mask, pair)] = rows
                    type1_rows = [row for row in rows if row.get("TypeFlag") == "1"]
                    if len(type1_rows) > 1:
                        raise ValueError(f"{mark}/{case}: more than one TypeFlag=1 row")
                    if type1_rows:
                        mark_rows["captures_with_type1"] += 1
                    mark_rows["type1"] += len(type1_rows)
                    for row in rows:
                        if row.get("TypeFlag") != "2":
                            continue
                        mark_rows["type2"] += 1
                        bank = int(row["File Index"])
                        position = int(row["PCM Pos"])
                        if not 0 <= bank < len(offsets):
                            mismatches.append(f"{mark}/{case}: invalid bank {bank}")
                        elif position not in offsets[bank]:
                            mismatches.append(f"{mark}/{case}: unmapped PCM offset {bank}/{position}")

    for mark in selected_marks:
        files = list(PROBE.glob(f"mi-interword-punct-{mark}-*.*.dtt"))
        if len(files) != 2 * totals[mark]["calls"]:
            raise ValueError(
                f"{mark}: expected {2 * totals[mark]['calls']} paired files, found {len(files)}"
            )
    if mismatches:
        raise ValueError("selected-bank offset mismatch: " + "; ".join(mismatches[:20]))

    call_total = sum(int(values["calls"]) for values in totals.values())
    print(
        f"validated calls={call_total} paired captures={2 * call_total}; "
        "every call returned raw EAX 1"
    )
    for mark in selected_marks:
        values = totals[mark]
        print(
            f"{mark}: calls={values['calls']} TypeFlag=2 rows={values['type2']} "
            f"TypeFlag=1 rows={values['type1']} captures-with-TypeFlag=1="
            f"{values['captures_with_type1']}"
        )
        for mask in MASKS:
            affected: list[str] = []
            unaffected: list[str] = []
            size_counts: Counter[str] = Counter()
            before_counts: Counter[int] = Counter()
            after_counts: Counter[int] = Counter()
            by_first: Counter[str] = Counter()
            for first in LETTERS:
                for second in LETTERS:
                    pair = first + second
                    rows = rows_by_case[(mark, mask, pair)]
                    type1_indexes = [i for i, row in enumerate(rows) if row.get("TypeFlag") == "1"]
                    if type1_indexes:
                        affected.append(pair)
                        by_first[first] += 1
                    else:
                        unaffected.append(pair)
                    for index in type1_indexes:
                        row = rows[index]
                        size_counts[row.get("Size", "?")] += 1
                        before = sum(item.get("TypeFlag") == "2" for item in rows[:index])
                        after = sum(item.get("TypeFlag") == "2" for item in rows[index + 1 :])
                        if not before or not after:
                            raise ValueError(f"{mark}/{case}: silence row is not between phone rows")
                        before_counts[before] += 1
                        after_counts[after] += 1
            if len(affected) == 676:
                first_counts = "all 26 initial letters"
            elif len(affected) >= 650:
                first_counts = "all initial letters except " + ",".join(
                    f"{letter}:{by_first[letter]}" for letter in LETTERS if by_first[letter] != 26
                )
            else:
                first_counts = ",".join(
                    f"{letter}:{by_first[letter]}" for letter in LETTERS if by_first[letter]
                ) or "-"
            pair_summary = ",".join(affected) if len(affected) <= 100 else f"all {len(affected)} pairs"
            unaffected_summary = (
                ",".join(unaffected) if len(unaffected) <= 30 else f"{len(unaffected)} pairs"
            )
            print(
                f"  {mask}: affected_pairs={len(affected)} by_first={first_counts} "
                f"pairs={pair_summary} unaffected={unaffected_summary or '-'} sizes={dict(size_counts)} "
                f"TypeFlag2_rows_before={dict(sorted(before_counts.items()))} "
                f"TypeFlag2_rows_after={dict(sorted(after_counts.items()))}"
            )
    print("selected-bank PCM offsets=all matched")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
