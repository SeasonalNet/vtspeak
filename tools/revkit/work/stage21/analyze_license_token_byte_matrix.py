#!/usr/bin/env python3
"""Validate the value-suppressed license-token mutation capture."""

from __future__ import annotations

import re
import sys
from pathlib import Path


CASE = re.compile(
    r"^TOKEN_BYTE index=(\d+) class=(selected|gap) "
    r"mutation=(next_hex|nonhex|case_flip) result=(-?\d+)$"
)
BASELINE = re.compile(r"^TOKEN_MATRIX baseline=(-?\d+) record_bytes=(\d+) token_length=96$")
SUMMARY = re.compile(r"^TOKEN_MATRIX_SUMMARY (.+)$")
ACCEPTED_FLIP = re.compile(r"^TOKEN_CASEFLIP_ACCEPTED index=(\d+)$")
COMBO = re.compile(r"^TOKEN_CASE_COMBO mask=(\d+) result=(-?\d+)$")


def main() -> int:
    if len(sys.argv) != 2:
        print(f"usage: {Path(sys.argv[0]).name} CAPTURE", file=sys.stderr)
        return 2
    capture = Path(sys.argv[1]).read_text(encoding="utf-8", errors="replace")
    baseline_lines = [line for line in capture.splitlines() if BASELINE.fullmatch(line)]
    if len(baseline_lines) != 1:
        print("expected exactly one sanitized baseline line", file=sys.stderr)
        return 1
    baseline_match = BASELINE.fullmatch(baseline_lines[0])
    assert baseline_match is not None
    baseline = int(baseline_match.group(1))
    if baseline != 0:
        print(f"expected accepted baseline return 0, got {baseline}", file=sys.stderr)
        return 1
    summary_lines = [line for line in capture.splitlines() if SUMMARY.fullmatch(line)]
    if len(summary_lines) != 1:
        print("expected exactly one mutation summary line", file=sys.stderr)
        return 1
    summary_match = SUMMARY.fullmatch(summary_lines[0])
    assert summary_match is not None
    summary = {
        key: int(value)
        for key, value in re.findall(r"([a-z_]+)=(\d+)", summary_match.group(1))
    }
    for key in (
        "caseflip_count",
        "caseflip_pass",
        "caseflip_fail",
        "combo_count",
        "combo_pass",
        "combo_fail",
    ):
        summary.setdefault(key, 0)

    cases: dict[tuple[int, str], tuple[str, int]] = {}
    accepted_flip_positions: list[int] = []
    combos: dict[int, int] = {}
    for line in capture.splitlines():
        accepted_match = ACCEPTED_FLIP.fullmatch(line)
        if accepted_match is not None:
            accepted_flip_positions.append(int(accepted_match.group(1)))
            continue
        combo_match = COMBO.fullmatch(line)
        if combo_match is not None:
            mask, result = combo_match.groups()
            mask_value = int(mask)
            if mask_value in combos:
                print(f"duplicate case-combination mask={mask_value}", file=sys.stderr)
                return 1
            combos[mask_value] = int(result)
            continue
        match = CASE.fullmatch(line)
        if match is None:
            continue
        index, byte_class, mutation, result = match.groups()
        key = (int(index), mutation)
        if key in cases:
            print(f"duplicate case index={index} mutation={mutation}", file=sys.stderr)
            return 1
        cases[key] = (byte_class, int(result))

    expected_keys = {
        (index, mode)
        for index in range(96)
        for mode in ("next_hex", "nonhex")
    }
    base_keys = {key for key in cases if key[1] != "case_flip"}
    if base_keys != expected_keys:
        print(f"expected 192 unique base cases, captured {len(base_keys)}", file=sys.stderr)
        return 1

    caseflip_count = 0
    caseflip_pass = 0
    actual = {
        "selected_hex_pass": 0,
        "selected_hex_fail": 0,
        "gap_hex_pass": 0,
        "gap_hex_fail": 0,
        "selected_nonhex_pass": 0,
        "selected_nonhex_fail": 0,
        "gap_nonhex_pass": 0,
        "gap_nonhex_fail": 0,
        "combo_count": 0,
        "combo_pass": 0,
        "combo_fail": 0,
    }
    for (index, mutation), (byte_class, result) in cases.items():
        expected_class = "selected" if index % 6 < 2 else "gap"
        if byte_class != expected_class:
            print(f"wrong classification at byte index {index}", file=sys.stderr)
            return 1
        if mutation == "case_flip":
            caseflip_count += 1
            caseflip_pass += result == baseline
            continue
        family = "hex" if mutation == "next_hex" else "nonhex"
        outcome = "pass" if result == baseline else "fail"
        actual[f"{byte_class}_{family}_{outcome}"] += 1
        if result != -2:
            print(f"unexpected checker result at byte index {index} mutation={mutation}: {result}",
                  file=sys.stderr)
            return 1

    actual["caseflip_count"] = caseflip_count
    actual["caseflip_pass"] = caseflip_pass
    actual["caseflip_fail"] = caseflip_count - caseflip_pass
    if accepted_flip_positions and (
        len(accepted_flip_positions) != caseflip_pass
        or len(set(accepted_flip_positions)) != len(accepted_flip_positions)
    ):
        print("accepted case-flip position list does not match individual results", file=sys.stderr)
        return 1
    if combos:
        if len(accepted_flip_positions) != caseflip_pass:
            print("case-combination candidates do not match accepted individual flips", file=sys.stderr)
            return 1
        if len(accepted_flip_positions) > 16:
            print("case-combination product exceeds analyzer bound", file=sys.stderr)
            return 1
        expected_masks = set(range(1 << len(accepted_flip_positions)))
        if set(combos) != expected_masks:
            print("case-combination masks are incomplete", file=sys.stderr)
            return 1
        actual["combo_count"] = len(combos)
        actual["combo_pass"] = sum(result == baseline for result in combos.values())
        actual["combo_fail"] = len(combos) - actual["combo_pass"]
        if combos.get(0) != baseline:
            print("empty case-combination should preserve baseline", file=sys.stderr)
            return 1
    if summary != actual:
        print("summary counters do not match the per-byte capture", file=sys.stderr)
        return 1

    print(
        "PASS license_token_byte_matrix baseline=0 cases=192 "
        "selected_positions=32 gap_positions=64 all_digit_mutations_return=-2 "
        f"case_flip_calls={caseflip_count} case_flip_accepts={caseflip_pass} "
        f"case_combos={len(combos)} case_combo_accepts={actual['combo_pass']}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
