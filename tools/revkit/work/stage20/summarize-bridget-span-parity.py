#!/usr/bin/env python3
"""Compare candidate membership and continuity spans in the Bridget A/B traces."""

from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent
QUERY_RE = re.compile(
    r"BRIDGET_2013_QUERY call=(\d+) signature=([^ ]+) key=([^ ]+) count=(\d+)"
)
LIST_RE = re.compile(
    r"BRIDGET_CANDIDATE_LIST position=(\d+) count=(\d+)(.*?)(?="
    r"BRIDGET_CANDIDATE_LIST|BRIDGET_TRANSITION|BRIDGET_SELECTED|$)",
    re.DOTALL,
)
ROW_RE = re.compile(
    r"\[(\d+)\]=(\d+)\(span=(-?\d+),weighted=(-?\d+),key=([^)]+)\)"
)
TRANSITION_RE = re.compile(
    r"BRIDGET_TRANSITION context=(\d+) count=(\d+) units:([^\n]+)"
)
TRACES = {
    "signature": ROOT / "bridget-hello-signature-deeptrace-runtime.log",
    "structural": ROOT / "bridget-hello-structural-deeptrace-runtime.log",
}


def read_trace(
    path: Path,
) -> tuple[
    list[tuple[str, str, int]],
    dict[int, dict[int, int]],
    dict[int, set[int]],
]:
    text = path.read_text(errors="replace")
    queries = [
        (match.group(2), match.group(3), int(match.group(4)))
        for match in QUERY_RE.finditer(text)
    ]
    positions: dict[int, dict[int, int]] = {}
    for match in LIST_RE.finditer(text):
        position = int(match.group(1))
        expected = int(match.group(2))
        rows = {
            int(row.group(2)): int(row.group(3))
            for row in ROW_RE.finditer(match.group(3))
        }
        if len(rows) != expected:
            raise ValueError(
                f"{path.name}: position {position} parsed {len(rows)} of {expected} rows"
            )
        positions[position] = rows
    transitions = {
        int(match.group(1)): {int(unit) for unit in match.group(3).split()}
        for match in TRANSITION_RE.finditer(text)
    }
    return queries, positions, transitions


def main() -> None:
    traces = {name: read_trace(path) for name, path in TRACES.items()}
    query_sets = [trace[0] for trace in traces.values()]
    if len(query_sets[0]) != 10 or query_sets[0] != query_sets[1]:
        raise ValueError("A/B query signatures or keys differ")
    if any(count != 0 for _, _, count in query_sets[0]):
        raise ValueError("the A/B trace contains a nonempty exact-key lookup")
    print("exact queries: identical=10 all_zero=yes")
    print(
        "position\tsignature_total\tsignature_span2\tstructural_total\t"
        "structural_span2\tshared_ids\tshared_span_changes"
    )

    signature = traces["signature"][1]
    structural = traces["structural"][1]
    for position in sorted(signature):
        left = signature[position]
        right = structural[position]
        common = left.keys() & right.keys()
        changes = sum(left[unit] != right[unit] for unit in common)
        print(
            f"{position}\t{len(left)}\t{sum(span == 2 for span in left.values())}"
            f"\t{len(right)}\t{sum(span == 2 for span in right.values())}"
            f"\t{len(common)}\t{changes}"
        )

    left = signature[4]
    right = structural[4]
    differing = {
        unit for unit in left.keys() & right.keys() if left[unit] != right[unit]
    }
    left_next = signature[5]
    right_next = structural[5]
    changed_next_membership = {
        unit for unit in differing if (unit in left_next) != (unit in right_next)
    }
    print(
        "position4_span_flips_with_position5_membership_change="
        f"{len(changed_next_membership)}/{len(differing)}"
    )
    print(
        "position4_span_flips_signature_next_only="
        f"{sum(unit in left_next and unit not in right_next for unit in differing)}"
    )

    position5_flips = {
        unit
        for unit in signature[5].keys() & structural[5].keys()
        if signature[5][unit] == 1 and structural[5][unit] == 2
    }
    signature_predecessors = traces["signature"][2][4]
    structural_predecessors = traces["structural"][2][4]
    print(
        "position5_span_flips_in_context4_shortlist_signature/structural="
        f"{sum(unit in signature_predecessors for unit in position5_flips)}/"
        f"{sum(unit in structural_predecessors for unit in position5_flips)}"
    )


if __name__ == "__main__":
    main()
