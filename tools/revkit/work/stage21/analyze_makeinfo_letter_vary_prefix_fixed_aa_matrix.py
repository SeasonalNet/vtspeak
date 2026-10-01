#!/usr/bin/env python3
"""Validate and compare the all 16 case masks for fixed-AA-suffix prefix-pair contexts."""

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

MASKS = ("UUUU", "UUUL", "UULU", "UULL", "ULUU", "ULUL", "ULLU", "ULLL", "LUUU", "LUUL", "LULU", "LULL", "LLUU", "LLUL", "LLLU", "LLLL")
FORMS = ("plain", "question")
LETTERS = "abcdefghijklmnopqrstuvwxyz"
BATCH_BY_MASK = {
    "UUUU": "b1", "UUUL": "b1", "UULU": "b1", "UULL": "b1",
    "ULUU": "b2", "ULUL": "b2", "ULLU": "b2", "ULLL": "b2",
    "LUUU": "b3", "LUUL": "b3", "LULU": "b3", "LULL": "b3",
    "LLUU": "b4", "LLUL": "b4", "LLLU": "b4", "LLLL": "b4",
}


def capture_stem(pair: str, mask: str, form: str) -> str:
    batch = BATCH_BY_MASK[mask]
    return f"mi-letter-vary-prefix-fixed-aa-{batch}-pair-{pair}-{mask}-{form}"


def read_rows(
    pair: str, mask: str, form: str
) -> tuple[tuple[tuple[tuple[str, str], ...], ...], int]:
    stem = capture_stem(pair, mask, form)
    path = PROBE / f"{stem}.asc.dtt"
    blocks = path.read_text(encoding="ascii").split("\n\n")
    rows: list[tuple[tuple[str, str], ...]] = []
    type1 = 0
    for block in blocks:
        row = matrix.fields(block)
        type1 += row.get("TypeFlag") == "1"
        if row.get("TypeFlag") == "2":
            rows.append(tuple(sorted(row.items())))
    return tuple(rows), type1


def expected_text(pair: str, mask: str, form: str) -> str:
    letters = (pair[0], pair[1], "A", "A")
    text = " ".join(
        letter.upper() if case == "U" else letter.lower()
        for letter, case in zip(letters, mask, strict=True)
    )
    return text + ("?" if form == "question" else "")


def main() -> int:
    observations: list[tuple[str, str, str, str]] = []
    for batch in ("b1", "b2", "b3", "b4"):
        log = (PROBE / f"makeinfo-letter-vary-prefix-fixed-aa-{batch}-api.log").read_text(
            encoding="ascii"
        )
        batch_observations = re.findall(
            rf"MAKEINFO_PREFIX_PAIR_FIXED_AA batch={batch} case=([^ ]+) "
            r"raw_eax=(0x[0-9a-f]+) text=(.*?) path=([^\r\n]+)",
            log,
        )
        observations.extend(batch_observations)
    expected_calls = 26 * 26 * len(MASKS) * len(FORMS)
    if len(observations) != expected_calls:
        raise ValueError(f"expected {expected_calls} calls, found {len(observations)}")
    if any(eax != "0x1" for _, eax, _, _ in observations):
        raise ValueError("one or more calls did not return raw EAX 1")
    by_case = {case: (eax, text, path) for case, eax, text, path in observations}
    if len(by_case) != expected_calls:
        raise ValueError("the API log contains duplicate case labels")

    offset_maps = [matrix.load_dat_offsets(bank) for bank in matrix.BANKS]
    rows_by_case: dict[tuple[str, str, str], tuple[tuple[tuple[str, str], ...], ...]] = {}
    type1_total = 0
    type2_total = 0
    mismatches: list[str] = []
    for first in LETTERS:
        for second in LETTERS:
            pair = first + second
            for mask in MASKS:
                for form in FORMS:
                    case = f"{pair}-{mask}-{form}"
                    eax, text, path = by_case[case]
                    if eax != "0x1" or text != expected_text(pair, mask, form):
                        raise ValueError(f"unexpected API input for {case}")
                    expected_path = f"Z:/work/stage21/{capture_stem(pair, mask, form)}"
                    if path != expected_path:
                        raise ValueError(f"unexpected output prefix for {case}")
                    binary = PROBE / f"{capture_stem(pair, mask, form)}.bin.dtt"
                    if not binary.is_file():
                        raise ValueError(f"missing paired binary capture: {binary.name}")
                    rows, type1 = read_rows(pair, mask, form)
                    rows_by_case[(pair, mask, form)] = rows
                    type1_total += type1
                    type2_total += len(rows)
                    for packed in rows:
                        row = dict(packed)
                        bank = int(row["File Index"])
                        position = int(row["PCM Pos"])
                        if not 0 <= bank < len(offset_maps):
                            mismatches.append(f"{case}: invalid bank index {bank}")
                        elif position not in offset_maps[bank]:
                            mismatches.append(f"{case}: no indexed unit at PCM Pos {position}")

    if mismatches:
        raise ValueError("selected-bank PCM position mismatch: " + "; ".join(mismatches))

    print(
        f"validated calls={len(observations)} paired captures={2 * len(observations)} "
        f"TypeFlag=2 rows={type2_total}"
    )
    print(f"TypeFlag=1 rows={type1_total}; selected-bank PCM offsets=all matched")
    for mask in MASKS:
        question_effects: Counter[str] = Counter()
        for first in LETTERS:
            for second in LETTERS:
                pair = first + second
                plain = rows_by_case[(pair, mask, "plain")]
                question = rows_by_case[(pair, mask, "question")]
                plain_banks = {dict(row)["File Index"] for row in plain}
                question_banks = {dict(row)["File Index"] for row in question}
                if plain == question:
                    question_effects["identical"] += 1
                elif plain_banks == question_banks:
                    question_effects["records changed, same banks"] += 1
                else:
                    question_effects["bank set changed"] += 1
        print(f"{mask} plain/question={dict(question_effects)}")

    all_differences = 0
    for mask in MASKS:
        for form in FORMS:
            differences: Counter[str] = Counter()
            changed: list[str] = []
            for first in LETTERS:
                for second in LETTERS:
                    pair = first + second
                    baseline = rows_by_case[(pair, "UUUU", form)]
                    candidate = rows_by_case[(pair, mask, form)]
                    if baseline == candidate:
                        differences["identical"] += 1
                        continue
                    changed.append(pair)
                    baseline_banks = {dict(row)["File Index"] for row in baseline}
                    candidate_banks = {dict(row)["File Index"] for row in candidate}
                    if baseline_banks == candidate_banks:
                        differences["records changed, same banks"] += 1
                    else:
                        differences["bank set changed"] += 1
            if mask == "UUUU":
                all_differences += sum(differences.values()) - differences["identical"]
            print(f"{mask} {form} vs UUUU={dict(differences)}")
            if mask in ("ULUU", "ULUL"):
                exact = [pair for pair in (x + y for x in LETTERS for y in LETTERS) if pair not in changed]
                print(f"{mask} {form} exact pairs vs UUUU={','.join(exact)}")

    if all_differences != 0:
        raise ValueError("UUUU must be identical to itself")

    equivalent_masks: dict[
        tuple[tuple[tuple[tuple[str, str], ...], ...], ...], list[str]
    ] = {}
    for mask in MASKS:
        signature = tuple(
            rows_by_case[(first + second, mask, form)]
            for first in LETTERS
            for second in LETTERS
            for form in FORMS
        )
        equivalent_masks.setdefault(signature, []).append(mask)
    print(
        "field-for-field equivalence classes across all pairs/forms="
        + "; ".join(",".join(group) for group in equivalent_masks.values())
    )

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
