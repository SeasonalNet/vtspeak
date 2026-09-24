#!/usr/bin/env python3
"""Verify the controlled Stage 10 phone/context-row captures."""

from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / "stage10"
ROW = re.compile(r"^ROW index=(\d+) flags=(?:0x[0-9a-f]+|0) source=(\d+) class=(?:0x[0-9a-f]+|0) marker=(?:0x[0-9a-f]+|0) surface=(.*)$")
SURFACE = re.compile(r'"([^"]*)"')
EXPECTED = {
    "doctor-title": ("doctor", "Smith", "spoke", "."),
    "doctor-word": ("The", "doctor", "spoke", "."),
    "drive-address": ("forty", "two", "Main", "Drive", "."),
    "saint-name": ("Saint", "Paul", "spoke", "."),
    "street-address": ("forty", "two", "Main", "Street", "."),
    "apartment-bare": ("Apartment", "four", "is", "ready", "."),
    "apartment-period": ("Apartment", "four", "is", "ready", "."),
    "apartment-short-bare": ("Apartment", "four", "."),
    "apartment-short-period": ("Apartment", "four", "."),
    "apartment-mixed": ("Apartment", "four", "and", "Apt", "four", "."),
    "apartment-repeated": ("Apartment", "four", "and", "Apt", "four", "."),
    "apartment-second": ("Unit", "four", "and", "Apartment", "four", "."),
    "apartment-after-and": ("Hello", "and", "Apartment", "four", "."),
    "apartment-repeat-number": ("Apartment", "four", "and", "Apt", "five", "."),
    "apartment-sentence-repeat": ("Apartment", "four", "Apartment", "four", "."),
    "apartment-repeat-reversed": ("Apartment", "four", "and", "Unit", "four", "."),
    "apartment-comma": ("Unit", "four", ",", "Apartment", "four", "."),
    "apartment-comma-repeat": ("Apartment", "four", ",", "Apartment", "four", "."),
    "apartment-word-break-repeat": ("Apartment", "four", "and", "Hello", "and", "Apartment", "four", "."),
    "pm-bare": ("It", "is", "five", "PM", "today", "."),
    "pm-period": ("It", "is", "five", "PM", "today", "."),
    "company-inc-period": ("Acme", "incorporated", "hired", "Lee", "."),
    "company-inc-comma": ("Acme", ",", "incorporated", "hired", "Lee", "."),
    "company-llc": ("Acme", "LLC", "hired", "Lee", "."),
    "company-co-period": ("Acme", "Company", "hired", "Lee", "."),
    "number-no-period": ("Number", "four", "arrived", "."),
    "number-no-comma": ("No", ",", "thanks", "."),
    "mr-title": ("mister", "Smith", "spoke", "."),
    "mrs-title": ("misses", "Smith", "spoke", "."),
    "prof-title": ("professor", "Smith", "spoke", "."),
    "capt-title": ("captain", "Smith", "spoke", "."),
    "rev-title": ("reverend", "Smith", "spoke", "."),
    "company-corp": ("Acme", "corporation", "hired", "Lee", "."),
    "company-ltd": ("Acme", "limited", "hired", "Lee", "."),
    "company-gmbh": ("Acme", "GmbH", "hired", "Lee", "."),
    "month-date": ("We", "met", "on", "January", "fifth", "."),
    "month-standalone": ("It", "happened", "in", "January", "."),
    "measurement-ft": ("It", "is", "ten", "feet", "long", "."),
    "measurement-mi": ("The", "road", "is", "two", "miles", "long", "."),
    "apartment-context": ("The", "apt", "is", "ready", "."),
}
POS_EXPECTED = {
    "default": (0xFF, (0x36, 0x23, 0x2A, 0x0B, 0x36, 0x15)),
    "unknown": (0xFF, (0x36, 0x23, 0x2A, 0x0B, 0x36, 0x15)),
    "function": (0xFF, (0x36, 0x23, 0x2A, 0x0B, 0x36, 0x15)),
    "noun": (0x13, (0x36, 0x18, 0x2A, 0x1A, 0x15)),
    "modifier": (0x0E, (0x36, 0x18, 0x2A, 0x1A, 0x15)),
    "verb": (0x25, (0x36, 0x23, 0x2A, 0x0B, 0x36, 0x15)),
    "interjection": (0x24, (0x36, 0x23, 0x2A, 0x0B, 0x36, 0x15)),
}
TREE = re.compile(r"TREE_(SCALAR|VECTOR) tree=(0x[0-9a-f]+) return=(0x[0-9a-f]+) features:\n([^\n]+)\n([^\n]+)")
BRANCH_EXPECTED = {
    "apartment-repeated": ("Apt",),
    "apartment-repeat-number": ("Apt",),
    "apartment-sentence-repeat": ("Apt", "Apt"),
    "apartment-comma-repeat": ("Apt", "Apt"),
    "apartment-word-break-repeat": ("Apt", "Apt"),
    "apartment-mixed": ("Apt",),
    "apartment-second": ("Unit", "Apt"),
}


def capture_surfaces(path: Path) -> tuple[str, ...]:
    lines = path.read_text(errors="replace").splitlines()
    if not any("[Inferior 1 (Remote target) exited normally]" in line for line in lines):
        raise ValueError(f"{path}: process did not exit normally")
    counts = [int(match.group(1)) for line in lines if (match := re.match(r"CONTEXT_RETURN count=(\d+)", line))]
    if len(counts) != 1:
        raise ValueError(f"{path}: expected one context return, got {counts}")
    rows: list[str] = []
    for position, line in enumerate(lines):
        match = ROW.match(line)
        if match is None:
            continue
        if int(match.group(1)) != len(rows):
            raise ValueError(f"{path}: out-of-order row {match.group(1)}")
        surface = SURFACE.search(match.group(3))
        if surface is None and position + 1 < len(lines):
            surface = SURFACE.search(lines[position + 1])
        if surface is None:
            raise ValueError(f"{path}: row {len(rows)} lacks a surface")
        rows.append(surface.group(1))
    if len(rows) != counts[0]:
        raise ValueError(f"{path}: {len(rows)} rows versus declared count {counts[0]}")
    return tuple(rows)


def tree_inputs(name: str) -> list[tuple[str, str, str, tuple[int, ...]]]:
    path = ROOT / f"{name}-trees.log"
    log = path.read_text(errors="replace")
    if "[Inferior 1 (Remote target) exited normally]" not in log:
        raise ValueError(f"{path}: process did not exit normally")
    calls = []
    for match in TREE.finditer(log):
        data = " ".join((match.group(4), match.group(5)))
        data = re.sub(r"0x[0-9a-f]+:\s*", "", data)
        features = tuple(int(value) for value in re.findall(r"-?\d+", data))
        calls.append((match.group(1), match.group(2), match.group(3), features))
    return calls


def row_metadata(name: str, row: int) -> tuple[int, int, int, int]:
    path = ROOT / f"{name}.log"
    log = path.read_text(errors="replace")
    match = re.search(
        rf"ROW_META index={row} words:[^\n]*:\s*(\d+)\s+(\d+)\s+(\d+)\s+(\d+)",
        log,
    )
    if match is None:
        raise ValueError(f"{path}: row {row} lacks trailing metadata")
    values = tuple(int(match.group(index)) for index in range(1, 5))
    return (values[0], values[1], values[2], values[3])


def verify_inventory_batches() -> tuple[int, int]:
    captures = 0
    context_calls = 0
    for path in sorted(ROOT.glob("abbr-[ht]-*.log")):
        lines = path.read_text(errors="replace").splitlines()
        if "[Inferior 1 (Remote target) exited normally]" not in lines:
            raise ValueError(f"{path}: process did not exit normally")
        expected_count: int | None = None
        row_indices: list[int] = []
        for line in lines:
            returned = re.match(r"CONTEXT_RETURN count=(\d+)", line)
            if returned is not None:
                if expected_count is not None and row_indices != list(range(expected_count)):
                    raise ValueError(f"{path}: incomplete or unordered context rows")
                expected_count = int(returned.group(1))
                row_indices = []
                context_calls += 1
                continue
            row = ROW.match(line)
            if row is not None:
                row_indices.append(int(row.group(1)))
        if expected_count is None or row_indices != list(range(expected_count)):
            raise ValueError(f"{path}: missing or incomplete final context return")
        captures += 1
    if captures != 11:
        raise ValueError(f"expected 11 all-key H/T captures, got {captures}")
    return captures, context_calls


def main() -> None:
    for name, expected in EXPECTED.items():
        actual = capture_surfaces(ROOT / f"{name}.log")
        if actual != expected:
            raise ValueError(f"{name}: {actual!r} != {expected!r}")
        print(f"{name}: {' | '.join(actual)}")
    print(f"verified {len(EXPECTED)} controlled context captures")
    for part, (expected_status, expected_code) in POS_EXPECTED.items():
        path = ROOT / f"pos-record-single-{part}.log"
        log = path.read_text(errors="replace")
        if "[Inferior 1 (Remote target) exited normally]" not in log:
            raise ValueError(f"{path}: process did not exit normally")
        status = re.search(r"UPSTREAM row=0[^\n]*status=(0x[0-9a-f]+)", log)
        code = re.search(r"ROW_CODE index=0 bytes:\n[^\n]*:\s*((?:0x[0-9a-f]{2}\s*)+)", log)
        metadata = re.search(r"ROW_META index=0 words:[^\n]*:\s*(\d+)\s+(\d+)\s+(\d+)\s+(\d+)", log)
        if status is None or code is None or metadata is None:
            raise ValueError(f"{path}: missing status, code, or metadata")
        actual_code = tuple(int(value, 16) for value in re.findall(r"0x([0-9a-f]{2})", code.group(1)))
        actual_code = actual_code[: actual_code.index(0)]
        if int(status.group(1), 16) != expected_status or actual_code != expected_code:
            raise ValueError(f"{path}: unexpected status or phone sequence")
        if tuple(int(metadata.group(i)) for i in range(1, 5)) != (0, 1, 0, 0):
            raise ValueError(f"{path}: trailing dictionary metadata changed")
        print(f"record/{part}: status=0x{expected_status:02x}, code={actual_code}")
    print(f"verified {len(POS_EXPECTED)} part-of-speech controls")
    for first, second, expected_calls in (("noun", "modifier", 9), ("verb", "interjection", 10)):
        left = tree_inputs(f"pos-record-single-{first}")
        right = tree_inputs(f"pos-record-single-{second}")
        if len(left) != expected_calls or left != right:
            raise ValueError(f"{first}/{second}: tree input sequences differ")
        print(f"{first}/{second}: {expected_calls} matching tree inputs")
    feature_cases = {name: tree_inputs(f"feature-{name}") for name in ("p-ah0", "b-ah0", "t-ah0", "p-ah1", "p-ah2")}
    if any(len(calls) != 4 for calls in feature_cases.values()):
        raise ValueError("forced-phone controls did not each reach four tree calls")
    baseline = feature_cases["p-ah0"]
    expected_differences = {
        "b-ah0": {(0, 0): (27, 7), (1, 1): (27, 7)},
        "t-ah0": {(0, 0): (27, 31), (1, 1): (27, 31)},
        "p-ah1": {(1, 3): (0, 1), (2, 2): (0, 1), (3, 2): (0, 1), (3, 11): (0, 4)},
        "p-ah2": {(1, 3): (0, 2), (2, 2): (0, 2), (3, 2): (0, 2)},
    }
    for name, expected in expected_differences.items():
        actual = {}
        for call_index, (left, right) in enumerate(zip(baseline, feature_cases[name], strict=True)):
            if left[:3] != right[:3]:
                raise ValueError(f"{name}: tree call identity changed at {call_index}")
            actual.update({
                (call_index, feature_index): (old, new)
                for feature_index, (old, new) in enumerate(zip(left[3], right[3], strict=True))
                if old != new
            })
        if actual != expected:
            raise ValueError(f"{name}: tree input differences {actual!r} != {expected!r}")
        print(f"{name}: {len(actual)} controlled tree input changes")
    for name in ("apartment-mixed", "apartment-repeated", "apartment-second"):
        path = ROOT / f"{name}-guard.log"
        log = path.read_text(errors="replace")
        if "[Inferior 1 (Remote target) exited normally]" not in log:
            raise ValueError(f"{path}: process did not exit normally")
        apt = re.findall(r'^ADDRESS_GUARD token=(\d+) pending=(\d+) kind=(\d+) key=.*"Apt"$', log, re.MULTILINE)
        if apt != [("0", "0", "1")]:
            raise ValueError(f"{path}: unexpected Apt guard visits {apt!r}")
        print(f"{name}: one Apt address guard visit")
    for name, expected in BRANCH_EXPECTED.items():
        path = ROOT / f"{name}-branch.log"
        lines = path.read_text(errors="replace").splitlines()
        if "[Inferior 1 (Remote target) exited normally]" not in lines:
            raise ValueError(f"{path}: process did not exit normally")
        keys = []
        for line in lines:
            if not line.startswith("ADDRESS_BRANCH "):
                continue
            surface = SURFACE.search(line)
            if surface is not None and surface.group(1) in ("Apt", "Unit"):
                keys.append(surface.group(1))
        if tuple(keys) != expected:
            raise ValueError(f"{path}: address alias branches {tuple(keys)!r} != {expected!r}")
        print(f"{name}: address aliases reaching branch = {', '.join(keys)}")
    metadata_cases = {
        ("apartment-repeated", 0): (0, 1, 1, 1),
        ("apartment-repeated", 3): (0, 0, 0, 1),
        ("apartment-word-break-repeat", 3): (0, 1, 0, 0),
        ("apartment-word-break-repeat", 5): (0, 1, 1, 1),
    }
    for (name, row), expected in metadata_cases.items():
        actual = row_metadata(name, row)
        if actual != expected:
            raise ValueError(f"{name} row {row}: metadata {actual!r} != {expected!r}")
    print(f"verified {len(metadata_cases)} controlled dictionary-metadata rows")
    metadata_logs = (
        tuple(EXPECTED)
        + tuple(f"pos-record-single-{part}" for part in POS_EXPECTED)
        + ("pos-record-default", "pos-record-marked")
    )
    metadata_values = []
    for name in metadata_logs:
        path = ROOT / f"{name}.log"
        values = re.findall(
            r"ROW_META index=\d+ words:[^\n]*:\s*(\d+)\s+(\d+)\s+(\d+)\s+(\d+)",
            path.read_text(errors="replace"),
        )
        metadata_values.extend(tuple(map(int, row)) for row in values)
    if len(metadata_values) != 227 or len(set(metadata_values)) != 10:
        raise ValueError("full Stage 10 metadata row count or observed pattern count changed")
    print(f"verified {len(metadata_values)} metadata rows across {len(set(metadata_values))} observed patterns")
    captures, context_calls = verify_inventory_batches()
    print(f"verified {captures} all-key H/T inventory batch captures across {context_calls} context returns")


if __name__ == "__main__":
    main()
