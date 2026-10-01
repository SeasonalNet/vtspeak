#!/usr/bin/env python3
"""Validate one complete Stage 21 four-token letter-pair mask batch."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
PROBE = ROOT / "tools/revkit/work/stage21"
SCRIPTS = ROOT / "tools/revkit/scripts"
sys.path.insert(0, str(SCRIPTS))
import analyze_makeinfo_alphabet_matrix as matrix  # noqa: E402

BATCH_MASKS = {
    "b1": ("UULU", "UULL", "LUUU", "LUUL"),
    "b2": ("LULU", "LULL", "LLUU", "LLUL"),
    "b3": ("ULLU", "ULLL", "LLLU", "LLLL"),
}
LETTERS = "abcdefghijklmnopqrstuvwxyz"


def main() -> int:
    if len(sys.argv) != 2 or sys.argv[1] not in BATCH_MASKS:
        raise SystemExit("usage: check_makeinfo_letter_fullmask_batch.py b1|b2|b3")
    batch = sys.argv[1]
    masks = BATCH_MASKS[batch]
    log = (PROBE / f"makeinfo-letter-a-prefix-fullmask-{batch}-api.log").read_text(
        encoding="ascii"
    )
    observations = re.findall(
        rf"MAKEINFO_A_PREFIX_FULLMASK_PAIR batch={batch} case=([^ ]+) "
        r"raw_eax=(0x[0-9a-f]+) text=(.*?) path=([^\r\n]+)",
        log,
    )
    expected_calls = 26 * 26 * len(masks) * 2
    if len(observations) != expected_calls:
        raise ValueError(f"expected {expected_calls} calls, found {len(observations)}")
    by_case = {case: (eax, text, path) for case, eax, text, path in observations}
    if len(by_case) != expected_calls:
        raise ValueError("duplicate case labels in API log")
    if any(eax != "0x1" for eax, _, _ in by_case.values()):
        raise ValueError("one or more calls did not return raw EAX 1")

    offset_maps = [matrix.load_dat_offsets(bank) for bank in matrix.BANKS]
    type1_total = 0
    type2_total = 0
    for first in LETTERS:
        for second in LETTERS:
            pair = first + second
            for mask in masks:
                for form in ("plain", "question"):
                    case = f"{pair}-{mask}-{form}"
                    _, text, path = by_case[case]
                    source = ("A", "a", first, second)
                    expected_text = " ".join(
                        letter.upper() if case_flag == "U" else letter.lower()
                        for letter, case_flag in zip(source, mask, strict=True)
                    )
                    if form == "question":
                        expected_text += "?"
                    stem = f"mi-letter-a-prefix-fullmask-{batch}-pair-{case}"
                    if text != expected_text:
                        raise ValueError(f"unexpected text for {case}: {text!r}")
                    if path != f"Z:/work/stage21/{stem}":
                        raise ValueError(f"unexpected output path for {case}: {path!r}")
                    if not (PROBE / f"{stem}.bin.dtt").is_file():
                        raise ValueError(f"missing paired binary file for {case}")
                    blocks = (PROBE / f"{stem}.asc.dtt").read_text(
                        encoding="ascii"
                    ).split("\n\n")
                    for block in blocks:
                        row = matrix.fields(block)
                        if row.get("TypeFlag") == "1":
                            type1_total += 1
                        elif row.get("TypeFlag") == "2":
                            type2_total += 1
                            bank = int(row["File Index"])
                            position = int(row["PCM Pos"])
                            if not 0 <= bank < len(offset_maps):
                                raise ValueError(f"invalid bank in {case}: {bank}")
                            if position not in offset_maps[bank]:
                                raise ValueError(
                                    f"unmapped PCM position in {case}: {bank}/{position}"
                                )

    captures = list(PROBE.glob(f"mi-letter-a-prefix-fullmask-{batch}-pair-*.*.dtt"))
    if len(captures) != 2 * expected_calls:
        raise ValueError(f"expected {2 * expected_calls} paired files, found {len(captures)}")
    print(
        f"validated batch={batch} calls={expected_calls} paired captures={len(captures)} "
        f"TypeFlag=2 rows={type2_total} TypeFlag=1 rows={type1_total}; offsets matched"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
