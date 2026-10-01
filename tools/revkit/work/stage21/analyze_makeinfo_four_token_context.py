#!/usr/bin/env python3
"""Validate the bounded four-token MakeInfo context probe."""

from __future__ import annotations

import re
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
PROBE = ROOT / "tools/revkit/work/stage21"
SCRIPTS = ROOT / "tools/revkit/scripts"
sys.path.insert(0, str(SCRIPTS))
import analyze_makeinfo_alphabet_matrix as matrix  # noqa: E402

MASKS = (
    "UUUU", "UUUL", "UULU", "UULL", "ULUU", "ULUL", "ULLU", "ULLL",
    "LUUU", "LUUL", "LULU", "LULL", "LLUU", "LLUL", "LLLU", "LLLL",
)
CONTEXTS = ("aa", "ag")
FORMS = ("plain", "question")
SWEEP_MASKS = ("UUUU", "UUUL", "ULUU", "ULUL")


def read_rows(context: str, mask: str, form: str) -> tuple[list[dict[str, str]], int]:
    path = PROBE / f"mi-letter-four-context-{context}-{mask}-{form}.asc.dtt"
    blocks = path.read_text(encoding="ascii").split("\n\n")
    rows: list[dict[str, str]] = []
    type1 = 0
    for block in blocks:
        row = matrix.fields(block)
        type1 += row.get("TypeFlag") == "1"
        if row.get("TypeFlag") == "2":
            rows.append(row)
    return rows, type1


def read_sweep_rows(
    context: str, letter: str, mask: str, form: str
) -> tuple[list[dict[str, str]], int]:
    path = PROBE / f"mi-letter-fourth-sweep-{context}-{letter}-{mask}-{form}.asc.dtt"
    blocks = path.read_text(encoding="ascii").split("\n\n")
    rows: list[dict[str, str]] = []
    type1 = 0
    for block in blocks:
        row = matrix.fields(block)
        type1 += row.get("TypeFlag") == "1"
        if row.get("TypeFlag") == "2":
            rows.append(row)
    return rows, type1


def main() -> int:
    log = (PROBE / "makeinfo-letter-four-token-context-api.log").read_text(
        encoding="ascii"
    )
    observations = re.findall(
        r"MAKEINFO_FOUR_TOKEN case=([^ ]+) raw_eax=(0x[0-9a-f]+) text=(.*?) path=([^\r\n]+)",
        log,
    )
    if len(observations) != 64:
        raise ValueError(f"expected 64 logged calls, found {len(observations)}")
    if any(eax != "0x1" for _, eax, _, _ in observations):
        raise ValueError("one or more calls did not return raw EAX 1")

    sweep_log = (PROBE / "makeinfo-letter-fourth-token-sweep-api.log").read_text(
        encoding="ascii"
    )
    sweep_observations = re.findall(
        r"MAKEINFO_FOURTH_SWEEP case=([^ ]+) raw_eax=(0x[0-9a-f]+) text=(.*?) path=([^\r\n]+)",
        sweep_log,
    )
    if len(sweep_observations) != 400:
        raise ValueError(f"expected 400 logged sweep calls, found {len(sweep_observations)}")
    if any(eax != "0x1" for _, eax, _, _ in sweep_observations):
        raise ValueError("one or more sweep calls did not return raw EAX 1")

    offset_maps = [matrix.load_dat_offsets(bank) for bank in matrix.BANKS]
    type1_total = 0
    type2_total = 0
    mismatches: list[str] = []
    for context in CONTEXTS:
        for mask in MASKS:
            letters = ["A", "a", context[-1], "d"]
            expected = " ".join(
                char.upper() if case == "U" else char.lower()
                for char, case in zip(letters, mask, strict=True)
            )
            for form in FORMS:
                case = f"{context}-{mask}-{form}"
                observed = [item for item in observations if item[0] == case]
                if len(observed) != 1 or observed[0][2] != expected + ("?" if form == "question" else ""):
                    raise ValueError(f"unexpected or missing trace input for {case}")
                binary = PROBE / f"mi-letter-four-context-{case}.bin.dtt"
                if not binary.is_file():
                    raise ValueError(f"missing paired binary capture: {binary.name}")
                rows, type1 = read_rows(context, mask, form)
                type1_total += type1
                type2_total += len(rows)
                for row in rows:
                    bank = int(row["File Index"])
                    position = int(row["PCM Pos"])
                    if not 0 <= bank < len(offset_maps) or position not in offset_maps[bank]:
                        mismatches.append(f"{case}: bank={bank} PCM Pos={position}")

    sweep_by_case = {item[0]: item for item in sweep_observations}
    for context in CONTEXTS:
        for letter in "abcdefghijklmnopqrstuvwxyz":
            if letter == "d":
                continue
            for mask in SWEEP_MASKS:
                letters = ["A", "a", context[-1], letter]
                expected = " ".join(
                    char.upper() if case == "U" else char.lower()
                    for char, case in zip(letters, mask, strict=True)
                )
                for form in FORMS:
                    case = f"{context}-{letter}-{mask}-{form}"
                    observed = sweep_by_case.get(case)
                    if observed is None or observed[2] != expected + ("?" if form == "question" else ""):
                        raise ValueError(f"unexpected or missing sweep input for {case}")
                    binary = PROBE / f"mi-letter-fourth-sweep-{case}.bin.dtt"
                    if not binary.is_file():
                        raise ValueError(f"missing paired binary capture: {binary.name}")
                    rows, type1 = read_sweep_rows(context, letter, mask, form)
                    type1_total += type1
                    type2_total += len(rows)
                    for row in rows:
                        bank = int(row["File Index"])
                        position = int(row["PCM Pos"])
                        if not 0 <= bank < len(offset_maps) or position not in offset_maps[bank]:
                            mismatches.append(f"{case}: bank={bank} PCM Pos={position}")

    if mismatches:
        raise ValueError("selected-bank PCM position mismatch: " + "; ".join(mismatches))

    total_calls = len(observations) + len(sweep_observations)
    print(f"validated calls={total_calls} TypeFlag=2 rows={type2_total}")
    print(f"TypeFlag=1 rows={type1_total}; selected-bank PCM offsets=all matched")
    for context in CONTEXTS:
        by_form: dict[str, Counter[str]] = {form: Counter() for form in FORMS}
        question_effects: Counter[str] = Counter()
        for mask in MASKS:
            plain, _ = read_rows(context, mask, "plain")
            question, _ = read_rows(context, mask, "question")
            plain_key = tuple(tuple(sorted(row.items())) for row in plain)
            question_key = tuple(tuple(sorted(row.items())) for row in question)
            plain_banks = {row["File Index"] for row in plain}
            question_banks = {row["File Index"] for row in question}
            if plain_key == question_key:
                question_effects["identical"] += 1
            elif plain_banks == question_banks:
                question_effects["records changed, same banks"] += 1
            else:
                question_effects["bank set changed"] += 1
            for form in FORMS:
                baseline, _ = read_rows(context, "UUUU", form)
                candidate, _ = read_rows(context, mask, form)
                baseline_key = tuple(tuple(sorted(row.items())) for row in baseline)
                candidate_key = tuple(tuple(sorted(row.items())) for row in candidate)
                if baseline_key == candidate_key:
                    by_form[form]["identical to UUUU"] += 1
                elif {row["File Index"] for row in baseline} == {
                    row["File Index"] for row in candidate
                }:
                    by_form[form]["records changed, same banks"] += 1
                else:
                    by_form[form]["bank set changed"] += 1
        print(f"{context}: plain/question={dict(question_effects)}")
        for form in FORMS:
            print(f"{context} {form} vs UUUU={dict(by_form[form])}")
    for context in CONTEXTS:
        print(f"{context} A-Z fourth-token sweep vs UUUU")
        for mask in SWEEP_MASKS:
            for form in FORMS:
                changed: list[str] = []
                bank_changes: list[str] = []
                same_bank_changes: list[str] = []
                for letter in "abcdefghijklmnopqrstuvwxyz":
                    baseline = (
                        read_rows(context, "UUUU", form)[0]
                        if letter == "d"
                        else read_sweep_rows(context, letter, "UUUU", form)[0]
                    )
                    candidate = (
                        read_rows(context, mask, form)[0]
                        if letter == "d"
                        else read_sweep_rows(context, letter, mask, form)[0]
                    )
                    baseline_key = tuple(tuple(sorted(row.items())) for row in baseline)
                    candidate_key = tuple(tuple(sorted(row.items())) for row in candidate)
                    if baseline_key != candidate_key:
                        changed.append(letter)
                        baseline_banks = {row["File Index"] for row in baseline}
                        candidate_banks = {row["File Index"] for row in candidate}
                        (bank_changes if baseline_banks != candidate_banks else same_bank_changes).append(letter)
                print(
                    f"  {mask} {form}: changed={''.join(changed) or '-'} "
                    f"bank={''.join(bank_changes) or '-'} "
                    f"same-bank={''.join(same_bank_changes) or '-'}"
                )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
