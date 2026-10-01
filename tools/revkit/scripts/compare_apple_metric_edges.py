#!/usr/bin/env python3
"""Compare Apple context-2 metric codes and predecessor distances offline."""

from __future__ import annotations

import argparse
import json
import re
import statistics
import struct
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
SCRIPTS = ROOT / "tools" / "revkit" / "scripts"
sys.path.insert(0, str(SCRIPTS))

from inspect_legacy_unit_idx import parse_index_data  # noqa: E402
from inspect_unit_idx import inspect as inspect_versioned  # noqa: E402

BANKS = (
    ("unit-gen.idx", 0),
    ("unit-gen2.idx", 179_995),
    ("unit-num.idx", 278_128),
    ("unit-etc.idx", 279_125),
    ("unit-alp.idx", 282_791),
)
OLD_NEIGHBOR = re.compile(
    r"APPLE_OLD_NEIGHBOR_LINK position=2 candidate=(\d+) "
    r"predecessor_index=\d+ predecessor_unit=(\d+)"
)
NEW_CANDIDATE = re.compile(
    r"APPLE_TRANSITION_CANDIDATE context=2 index=\d+ unit=(\d+) "
    r"cumulative_cost=([\d.]+) predecessor_index=\d+ predecessor_unit=(\d+)"
)


def read_codes(directory: Path, legacy: bool) -> dict[int, tuple[int, int, int]]:
    codes: dict[int, tuple[int, int, int]] = {}
    for filename, global_start in BANKS:
        path = directory / filename
        raw = path.read_bytes()
        if legacy:
            index = parse_index_data(path, raw)
            count = index.unit_count
            metrics_start = index.block_start + 27 * count
        else:
            index = inspect_versioned(path)
            count = int(index["unit_count"])
            metrics_start = int(index["opaque_block_start"]) + 28 * count

        for local_row in range(count):
            codes[global_start + local_row] = tuple(
                struct.unpack_from(
                    "<H", raw, metrics_start + count * group * 4 + local_row * 2
                )[0]
                & 0x3FFF
                for group in range(3)
            )
    return codes


def load_cepdist(path: Path) -> tuple[int, tuple[float, ...]]:
    raw = path.read_bytes()
    if len(raw) < 2 or (len(raw) - 2) % 4:
        raise ValueError(f"{path}: malformed 16-bit header and float table extent")
    count = struct.unpack_from("<H", raw)[0]
    values = struct.unpack_from(f"<{(len(raw) - 2) // 4}f", raw, 2)
    expected = count * (count + 1) // 2
    if len(values) != expected:
        raise ValueError(f"{path}: expected {expected} triangular values, found {len(values)}")
    return count, values


def table_distance(values: tuple[float, ...], left: int, right: int) -> float:
    high, low = max(left, right), min(left, right)
    return values[high * (high + 1) // 2 + low]


def summarize(values: list[float]) -> dict[str, float | int]:
    return {
        "count": len(values),
        "zero_count": sum(value == 0 for value in values),
        "min": min(values),
        "median": statistics.median(values),
        "mean": statistics.mean(values),
        "max": max(values),
    }


def analyze(args: argparse.Namespace) -> dict[str, object]:
    old_codes = read_codes(args.legacy_dir, legacy=True)
    new_codes = read_codes(args.adapted_dir, legacy=False)

    old_links = {
        int(unit): int(predecessor)
        for unit, predecessor in OLD_NEIGHBOR.findall(args.old_log.read_text())
    }
    new_links = {
        int(unit): int(predecessor)
        for unit, _cost, predecessor in NEW_CANDIDATE.findall(args.new_log.read_text())
    }
    old_ids, new_ids = set(old_links), set(new_links)
    shared, union = old_ids & new_ids, old_ids | new_ids
    if not old_ids or not new_ids:
        raise ValueError("candidate traces did not contain Apple context-2 rows")

    table_count, values = load_cepdist(args.cepdist)
    preserved_by_group = [
        sum(old_codes[unit][group] == new_codes[unit][group] for unit in union)
        for group in range(3)
    ]
    old_predecessors = {old_links[unit] for unit in shared}
    if len(old_predecessors) != 1:
        raise ValueError(f"expected one old predecessor in the shared set: {old_predecessors}")
    old_predecessor = next(iter(old_predecessors))

    edge_rows = []
    old_distances = []
    new_distances = []
    for unit in sorted(shared):
        old_distance = table_distance(
            values, old_codes[unit][0], old_codes[old_predecessor][2]
        )
        new_predecessor = new_links[unit]
        new_distance = table_distance(
            values, new_codes[unit][0], new_codes[new_predecessor][2]
        )
        old_distances.append(old_distance)
        new_distances.append(new_distance)
        edge_rows.append(
            {
                "unit": unit,
                "metric_codes_preserved": old_codes[unit] == new_codes[unit],
                "current_group0": new_codes[unit][0],
                "old_predecessor": old_predecessor,
                "old_predecessor_group2": old_codes[old_predecessor][2],
                "old_raw_distance": old_distance,
                "new_predecessor": new_predecessor,
                "new_predecessor_group2": new_codes[new_predecessor][2],
                "new_raw_distance": new_distance,
                "distance_delta_new_minus_old": new_distance - old_distance,
            }
        )

    return {
        "scope": "captured Apple context 2; raw table term only",
        "legacy_candidates": len(old_ids),
        "adapted_candidates": len(new_ids),
        "shared_candidates": len(shared),
        "old_only": sorted(old_ids - new_ids),
        "new_only": sorted(new_ids - old_ids),
        "union_metric_rows": len(union),
        "union_rows_with_all_three_codes_preserved": sum(
            old_codes[unit] == new_codes[unit] for unit in union
        ),
        "preserved_codes_by_group": preserved_by_group,
        "cepdist_dimension": table_count,
        "old_neighbor_predecessor": old_predecessor,
        "new_predecessor_count_on_shared_rows": len(
            {new_links[unit] for unit in shared}
        ),
        "shared_rows_raw_distance_using_old_links": summarize(old_distances),
        "shared_rows_raw_distance_using_adapted_trace_links": summarize(new_distances),
        "edge_rows": edge_rows,
    }


def main() -> int:
    default_evidence = (
        ROOT / "tools" / "revkit" / "work" / "corpus-parity" / "stage20"
        / "apple-continuity-bit7"
    )
    default_overlay = (
        ROOT / "tools" / "revkit" / "work" / "stage20"
        / "index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks"
    )
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--legacy-dir", type=Path, default=ROOT / "data-kate/M16/mc_idx_tbl")
    parser.add_argument("--adapted-dir", type=Path, default=default_overlay)
    parser.add_argument(
        "--old-log",
        type=Path,
        default=default_evidence / "old/apple-legacy-neighbor-links-gdb.log",
    )
    parser.add_argument(
        "--new-log",
        type=Path,
        default=default_evidence
        / "new-full-marker-metric-words-apple/adapted-forced-pah0-winedbg.log",
    )
    parser.add_argument(
        "--cepdist", type=Path, default=ROOT / "data-kate/M16/ttsdata/dist_tbl/cepdist.tbl"
    )
    args = parser.parse_args()
    print(json.dumps(analyze(args), indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
