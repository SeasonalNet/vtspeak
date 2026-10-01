#!/usr/bin/env python3
"""Cross-check Stage 21 MakeInfo alphabet rows against the Paul 2013 unit indexes."""

from __future__ import annotations

import struct
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
PROBE = ROOT / "tools/revkit/work/stage21"
INDEX_DIR = ROOT / "data-paul/M16/mc_idx_tbl"
BANKS = ("merged-gen", "merged-num", "merged-etc", "merged-alp")
HEADER_BYTES = 45
RECORD_BYTES = 19


def load_dat_offsets(bank: str) -> dict[int, int]:
    path = INDEX_DIR / f"unit-{bank.removeprefix('merged-')}.idx"
    raw = path.read_bytes()
    if len(raw) < HEADER_BYTES:
        raise ValueError(f"{path}: truncated header")
    count = struct.unpack_from("<I", raw, 39)[0]
    stride = struct.unpack_from("<H", raw, 43)[0]
    if stride != RECORD_BYTES:
        raise ValueError(f"{path}: unexpected unit record stride {stride}")
    end = HEADER_BYTES + count * stride
    if end > len(raw):
        raise ValueError(f"{path}: unit records extend past EOF")
    return {
        struct.unpack_from("<I", raw, offset)[0]: (offset - HEADER_BYTES) // stride
        for offset in range(HEADER_BYTES, end, stride)
    }


def fields(block: str) -> dict[str, str]:
    result: dict[str, str] = {}
    for line in block.splitlines():
        if " : " in line:
            key, value = line.split(" : ", 1)
            result[key] = value
    return result


def main() -> int:
    paths = sorted(
        (
            *PROBE.glob("mi-alphabet-v2-*.asc.dtt"),
            *PROBE.glob("mi-alphabet-lower-*.asc.dtt"),
            *PROBE.glob("mi-alphabet-context-*.asc.dtt"),
            *PROBE.glob("mi-alphabet-punct-[a-z]-*.asc.dtt"),
            *PROBE.glob("mi-alphabet-punct-lower-*.asc.dtt"),
            *PROBE.glob("mi-letter-pair-[a-z][a-z]-*.asc.dtt"),
            *PROBE.glob("mi-letter-pair-lower-*.asc.dtt"),
            *PROBE.glob("mi-letter-pair-mixed-*.asc.dtt"),
            *PROBE.glob("mi-letter-pair-punct-*.asc.dtt"),
            *PROBE.glob("mi-letter-triple-a-*.asc.dtt"),
        )
    )
    if len(paths) != 46254:
        raise ValueError(f"expected 46254 ASCII captures, found {len(paths)}")
    offsets = [load_dat_offsets(bank) for bank in BANKS]
    totals = [0] * len(BANKS)
    type1_total = 0
    type2_total = 0
    mismatches: list[str] = []
    phone_rows_by_capture: dict[str, tuple[tuple[tuple[str, str], ...], ...]] = {}

    for path in paths:
        body = path.read_text(encoding="ascii").split("\n\n")
        row_fields = [fields(block) for block in body]
        type1_total += sum(row.get("TypeFlag") == "1" for row in row_fields)
        row_fields = [row for row in row_fields if row.get("TypeFlag") == "2"]
        type2_total += len(row_fields)
        phone_rows_by_capture[path.name] = tuple(
            tuple(sorted(row.items())) for row in row_fields
        )
        used: set[int] = set()
        labels: list[str] = []
        for row in row_fields:
            bank_index = int(row["File Index"])
            dat_offset = int(row["PCM Pos"])
            if not 0 <= bank_index < len(BANKS):
                mismatches.append(f"{path.name}: invalid bank index {bank_index}")
                continue
            totals[bank_index] += 1
            used.add(bank_index)
            unit_index = offsets[bank_index].get(dat_offset)
            labels.append(f"{bank_index}@{unit_index}:{row.get('PhoneString', '?')}")
            if unit_index is None:
                mismatches.append(
                    f"{path.name}: {BANKS[bank_index]} has no unit at DAT offset {dat_offset}"
                )
        name = path.name.removeprefix("mi-alphabet-v2-").removeprefix(
            "mi-alphabet-lower-"
        ).removeprefix("mi-alphabet-context-").removeprefix(
            "mi-alphabet-punct-"
        ).removeprefix("mi-alphabet-punct-lower-").removesuffix(".asc.dtt")
        name = name.removeprefix("mi-letter-pair-lower-").removeprefix(
            "mi-letter-pair-"
        )
        print(
            f"{name}: rows={len(row_fields)} "
            f"banks={','.join(str(index) for index in sorted(used)) or '-'} "
            f"phones={','.join(labels)}"
        )

    print(f"captures={len(paths)} rows_by_bank={dict(zip(BANKS, totals))}")
    print(f"TypeFlag=2 detailed unit rows={type2_total}")
    print(f"TypeFlag=1 silence rows across captures={type1_total}")
    for case, baseline_prefix, punctuation_prefix in (
        ("upper", "mi-alphabet-v2-letter", "mi-alphabet-punct"),
        ("lower", "mi-alphabet-lower", "mi-alphabet-punct-lower"),
    ):
        punctuation_matches: dict[str, int] = {}
        punctuation_differences: dict[str, list[str]] = {}
        for punctuation in ("period", "comma", "bang", "question"):
            matches = 0
            differences: list[str] = []
            for letter in "abcdefghijklmnopqrstuvwxyz":
                base_name = f"{baseline_prefix}-{letter}.asc.dtt"
                punct_name = f"{punctuation_prefix}-{letter}-{punctuation}.asc.dtt"
                if base_name not in phone_rows_by_capture or punct_name not in phone_rows_by_capture:
                    continue
                if phone_rows_by_capture[base_name] == phone_rows_by_capture[punct_name]:
                    matches += 1
                else:
                    differences.append(letter)
            punctuation_matches[punctuation] = matches
            punctuation_differences[punctuation] = differences
        print(
            f"{case} terminal punctuation exact TypeFlag=2 row matches to standalone: "
            + ", ".join(f"{key}={value}/26" for key, value in punctuation_matches.items())
        )
        for punctuation, letters in punctuation_differences.items():
            print(f"{case} terminal {punctuation} differing letters: {','.join(letters) or '-'}")
    for punctuation in ("period", "comma", "bang", "question"):
        equal = sum(
            phone_rows_by_capture.get(f"mi-alphabet-punct-{letter}-{punctuation}.asc.dtt")
            == phone_rows_by_capture.get(
                f"mi-alphabet-punct-lower-{letter}-{punctuation}.asc.dtt"
            )
            for letter in "abcdefghijklmnopqrstuvwxyz"
        )
        print(f"upper/lower {punctuation} exact TypeFlag=2 row equality: {equal}/26")
    for case, prefix in (
        ("upper", "mi-letter-pair"),
        ("lower", "mi-letter-pair-lower"),
        ("upper-lower", "mi-letter-pair-mixed-upper-lower"),
        ("lower-upper", "mi-letter-pair-mixed-lower-upper"),
    ):
        pair_deltas = {"identical": 0, "records changed, same banks": 0, "banks changed": 0}
        pair_bank_rows = {"plain": Counter(), "question": Counter()}
        pair_num_counts = {"plain": 0, "question": 0}
        bank_transitions: Counter[tuple[str, str]] = Counter()
        for left in "abcdefghijklmnopqrstuvwxyz":
            for right in "abcdefghijklmnopqrstuvwxyz":
                plain_name = f"{prefix}-{left}{right}-plain.asc.dtt"
                question_name = f"{prefix}-{left}{right}-question.asc.dtt"
                plain = phone_rows_by_capture.get(plain_name)
                question = phone_rows_by_capture.get(question_name)
                if plain is None or question is None:
                    continue
                plain_bank_set = {dict(row).get("File Index") for row in plain}
                question_bank_set = {dict(row).get("File Index") for row in question}
                bank_transitions[
                    (
                        ",".join(sorted(plain_bank_set)),
                        ",".join(sorted(question_bank_set)),
                    )
                ] += 1
                for row in plain:
                    pair_bank_rows["plain"][dict(row).get("File Index", "?")] += 1
                for row in question:
                    pair_bank_rows["question"][dict(row).get("File Index", "?")] += 1
                if "1" in plain_bank_set:
                    pair_num_counts["plain"] += 1
                if "1" in question_bank_set:
                    pair_num_counts["question"] += 1
                if plain == question:
                    pair_deltas["identical"] += 1
                elif plain_bank_set == question_bank_set:
                    pair_deltas["records changed, same banks"] += 1
                else:
                    pair_deltas["banks changed"] += 1
        print(f"{case} two-letter plain/question record comparisons={pair_deltas}")
        print(f"{case} two-letter pair TypeFlag=2 rows by bank={dict(pair_bank_rows)}")
        print(f"{case} two-letter pairs using merged-num={pair_num_counts}")
        print(f"{case} two-letter bank-set transitions={dict(sorted(bank_transitions.items()))}")
    case_prefixes = {
        "upper": "mi-letter-pair",
        "lower": "mi-letter-pair-lower",
        "upper-lower": "mi-letter-pair-mixed-upper-lower",
        "lower-upper": "mi-letter-pair-mixed-lower-upper",
    }
    punctuation_case_labels = {
        "upper": "upper",
        "lower": "lower",
        "upper-lower": "mixed-upper-lower",
        "lower-upper": "mixed-lower-upper",
    }
    for case, case_label in punctuation_case_labels.items():
        prefix = f"mi-letter-pair-punct-{case_label}"
        for punctuation in ("period", "comma", "bang"):
            same_rows = 0
            bank_set_changes = 0
            changed_pairs: list[str] = []
            bank_changed_pairs: list[str] = []
            for left in "abcdefghijklmnopqrstuvwxyz":
                for right in "abcdefghijklmnopqrstuvwxyz":
                    punct_name = f"{prefix}-{left}{right}-{punctuation}.asc.dtt"
                    plain_prefix = case_prefixes[case]
                    plain_name = f"{plain_prefix}-{left}{right}-plain.asc.dtt"
                    punct_rows = phone_rows_by_capture.get(punct_name)
                    plain_rows = phone_rows_by_capture.get(plain_name)
                    if punct_rows is None or plain_rows is None:
                        continue
                    if punct_rows == plain_rows:
                        same_rows += 1
                    else:
                        changed_pairs.append(left + right)
                    punct_banks = {dict(row).get("File Index") for row in punct_rows}
                    plain_banks = {dict(row).get("File Index") for row in plain_rows}
                    if punct_banks != plain_banks:
                        bank_set_changes += 1
                        bank_changed_pairs.append(left + right)
            print(
                f"{case} terminal {punctuation} vs plain pair rows exact="
                f"{same_rows}/676 bank-set changes={bank_set_changes}/676"
            )
            if changed_pairs:
                print(f"{case} terminal {punctuation} changed pairs={changed_pairs}")
            if bank_changed_pairs:
                print(
                    f"{case} terminal {punctuation} bank-set changed pairs="
                    f"{bank_changed_pairs}"
                )
    for case, case_label in punctuation_case_labels.items():
        if case == "upper":
            continue
        differences = 0
        by_mark = Counter()
        differences_by_mark: dict[str, list[str]] = {
            mark: [] for mark in ("period", "comma", "bang")
        }
        lower_period_matches_upper_plain: list[str] = []
        for left in "abcdefghijklmnopqrstuvwxyz":
            for right in "abcdefghijklmnopqrstuvwxyz":
                for punctuation in ("period", "comma", "bang"):
                    baseline = phone_rows_by_capture.get(
                        f"mi-letter-pair-punct-upper-{left}{right}-{punctuation}.asc.dtt"
                    )
                    candidate = phone_rows_by_capture.get(
                        f"mi-letter-pair-punct-{case_label}-{left}{right}-{punctuation}.asc.dtt"
                    )
                    if baseline != candidate:
                        differences += 1
                        by_mark[punctuation] += 1
                        differences_by_mark[punctuation].append(left + right)
                    if case == "lower" and punctuation == "period":
                        upper_plain = phone_rows_by_capture.get(
                            f"{case_prefixes['upper']}-{left}{right}-plain.asc.dtt"
                        )
                        lower_plain = phone_rows_by_capture.get(
                            f"{case_prefixes['lower']}-{left}{right}-plain.asc.dtt"
                        )
                        if candidate == upper_plain and lower_plain != upper_plain:
                            lower_period_matches_upper_plain.append(left + right)
        print(
            f"uppercase vs {case} terminal punctuation differences="
            f"{differences}/2028 by mark={dict(by_mark)}"
        )
        if differences:
            print(
                f"uppercase vs {case} terminal punctuation differing pairs="
                f"{differences_by_mark}"
            )
        if case == "lower":
            print(
                "lowercase terminal period restores uppercase plain rows for="
                f"{lower_period_matches_upper_plain}"
            )
    triple_case_labels = {
        "upper": "upper",
        "lower": "lower",
        "upper-lower": "mixed-upper-lower",
        "lower-upper": "mixed-lower-upper",
        "lower-lower-upper": "lower-lower-upper",
        "lower-upper-lower": "lower-upper-lower",
        "upper-lower-upper": "upper-lower-upper",
        "upper-upper-lower": "upper-upper-lower",
    }
    for case, case_label in triple_case_labels.items():
        prefix = f"mi-letter-triple-a-{case_label}"
        deltas = Counter()
        bank_rows = {"plain": Counter(), "question": Counter()}
        for middle in "abcdefghijklmnopqrstuvwxyz":
            for right in "abcdefghijklmnopqrstuvwxyz":
                plain = phone_rows_by_capture.get(
                    f"{prefix}-{middle}{right}-plain.asc.dtt"
                )
                question = phone_rows_by_capture.get(
                    f"{prefix}-{middle}{right}-question.asc.dtt"
                )
                if plain is None or question is None:
                    continue
                plain_banks = {dict(row).get("File Index") for row in plain}
                question_banks = {dict(row).get("File Index") for row in question}
                if plain == question:
                    deltas["identical"] += 1
                elif plain_banks == question_banks:
                    deltas["records changed, same banks"] += 1
                else:
                    deltas["banks changed"] += 1
                for row in plain:
                    bank_rows["plain"][dict(row).get("File Index", "?")] += 1
                for row in question:
                    bank_rows["question"][dict(row).get("File Index", "?")] += 1
        print(f"A-leading triple {case} plain/question comparisons={deltas}")
        print(f"A-leading triple {case} rows by bank={dict(bank_rows)}")
    for case, case_label in triple_case_labels.items():
        if case == "upper":
            continue
        differences = 0
        by_form = Counter()
        bank_changes_by_form = Counter()
        same_bank_changes_by_form = Counter()
        identical_suffixes: dict[str, list[str]] = {"plain": [], "question": []}
        bank_transitions_by_form: dict[str, Counter[tuple[str, str]]] = {
            "plain": Counter(),
            "question": Counter(),
        }
        for middle in "abcdefghijklmnopqrstuvwxyz":
            for right in "abcdefghijklmnopqrstuvwxyz":
                for form in ("plain", "question"):
                    baseline = phone_rows_by_capture.get(
                        f"mi-letter-triple-a-upper-{middle}{right}-{form}.asc.dtt"
                    )
                    candidate = phone_rows_by_capture.get(
                        f"mi-letter-triple-a-{case_label}-{middle}{right}-{form}.asc.dtt"
                    )
                    if baseline != candidate:
                        differences += 1
                        by_form[form] += 1
                        baseline_banks = sorted(
                            {dict(row).get("File Index", "?") for row in baseline or ()}
                        )
                        candidate_banks = sorted(
                            {dict(row).get("File Index", "?") for row in candidate or ()}
                        )
                        bank_transitions_by_form[form][
                            (",".join(baseline_banks), ",".join(candidate_banks))
                        ] += 1
                        if baseline_banks == candidate_banks:
                            same_bank_changes_by_form[form] += 1
                        else:
                            bank_changes_by_form[form] += 1
                    else:
                        identical_suffixes[form].append(middle + right)
        print(
            f"uppercase vs A-leading triple {case} differences="
            f"{differences}/1352 by form={dict(by_form)}"
        )
        print(
            f"uppercase vs A-leading triple {case} bank-set changes by form="
            f"{dict(bank_changes_by_form)} same-bank row changes="
            f"{dict(same_bank_changes_by_form)}"
        )
        if case == "lower":
            print(
                "uppercase vs all-lower A-leading triple exact suffixes by form="
                f"{identical_suffixes}"
            )
            print(
                "uppercase vs all-lower A-leading triple bank transitions="
                f"{dict(bank_transitions_by_form)}"
            )
    for position in ("middle", "last"):
        for case, case_label in triple_case_labels.items():
            prefix = f"mi-letter-triple-a-{position}-{case_label}"
            deltas = Counter()
            bank_rows = {"plain": Counter(), "question": Counter()}
            for left in "abcdefghijklmnopqrstuvwxyz":
                for right in "abcdefghijklmnopqrstuvwxyz":
                    plain = phone_rows_by_capture.get(
                        f"{prefix}-{left}{right}-plain.asc.dtt"
                    )
                    question = phone_rows_by_capture.get(
                        f"{prefix}-{left}{right}-question.asc.dtt"
                    )
                    if plain is None or question is None:
                        continue
                    plain_banks = {dict(row).get("File Index") for row in plain}
                    question_banks = {dict(row).get("File Index") for row in question}
                    if plain == question:
                        deltas["identical"] += 1
                    elif plain_banks == question_banks:
                        deltas["records changed, same banks"] += 1
                    else:
                        deltas["banks changed"] += 1
                    for row in plain:
                        bank_rows["plain"][dict(row).get("File Index", "?")] += 1
                    for row in question:
                        bank_rows["question"][dict(row).get("File Index", "?")] += 1
            print(f"A-in-{position} triple {case} plain/question comparisons={deltas}")
            print(f"A-in-{position} triple {case} rows by bank={dict(bank_rows)}")
    for position in ("middle", "last"):
        for case, case_label in triple_case_labels.items():
            if case == "upper":
                continue
            differences = 0
            by_form = Counter()
            bank_changes_by_form = Counter()
            identical_suffixes: dict[str, list[str]] = {"plain": [], "question": []}
            for left in "abcdefghijklmnopqrstuvwxyz":
                for right in "abcdefghijklmnopqrstuvwxyz":
                    for form in ("plain", "question"):
                        baseline = phone_rows_by_capture.get(
                            f"mi-letter-triple-a-{position}-upper-{left}{right}-{form}.asc.dtt"
                        )
                        candidate = phone_rows_by_capture.get(
                            f"mi-letter-triple-a-{position}-{case_label}-{left}{right}-{form}.asc.dtt"
                        )
                        if baseline == candidate:
                            identical_suffixes[form].append(left + right)
                            continue
                        differences += 1
                        by_form[form] += 1
                        baseline_banks = {
                            dict(row).get("File Index") for row in baseline or ()
                        }
                        candidate_banks = {
                            dict(row).get("File Index") for row in candidate or ()
                        }
                        bank_changes_by_form[form] += baseline_banks != candidate_banks
            print(
                f"uppercase vs A-in-{position} triple {case} differences="
                f"{differences}/1352 by form={dict(by_form)} bank-set changes="
                f"{dict(bank_changes_by_form)}"
            )
            if case == "lower":
                print(
                    f"uppercase vs all-lower A-in-{position} exact suffixes="
                    f"{identical_suffixes}"
                )
    for case, prefix in case_prefixes.items():
        if case == "upper":
            continue
        differences = 0
        by_first = Counter()
        by_second = Counter()
        bank_set_changes: list[str] = []
        record_only_changes: list[str] = []
        for left in "abcdefghijklmnopqrstuvwxyz":
            for right in "abcdefghijklmnopqrstuvwxyz":
                for form in ("plain", "question"):
                    baseline = phone_rows_by_capture.get(
                        f"{case_prefixes['upper']}-{left}{right}-{form}.asc.dtt"
                    )
                    candidate = phone_rows_by_capture.get(
                        f"{prefix}-{left}{right}-{form}.asc.dtt"
                    )
                    if baseline != candidate:
                        differences += 1
                        by_first[left] += 1
                        by_second[right] += 1
                        baseline_banks = sorted(
                            {dict(row).get("File Index", "?") for row in baseline or ()}
                        )
                        candidate_banks = sorted(
                            {dict(row).get("File Index", "?") for row in candidate or ()}
                        )
                        change = (
                            f"{left}{right}-{form} "
                            f"banks={','.join(baseline_banks) or '-'}"
                            f"->{','.join(candidate_banks) or '-'}"
                        )
                        target = (
                            bank_set_changes
                            if baseline_banks != candidate_banks
                            else record_only_changes
                        )
                        target.append(change)
        print(f"uppercase vs {case} pair record differences={differences}/1352")
        print(f"uppercase vs {case} differences by first letter={dict(by_first)}")
        print(f"uppercase vs {case} differences by second letter={dict(by_second)}")
        if case == "lower":
            print(f"uppercase vs lower bank-set differences={bank_set_changes}")
            print(f"uppercase vs lower same-bank record differences={record_only_changes}")
    for mismatch in mismatches:
        print(f"mismatch: {mismatch}")
    if mismatches:
        return 1
    print("all PCM Pos values match a unit in the DTT-selected bank")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
