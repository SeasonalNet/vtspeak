#!/usr/bin/env python3
"""Compare recorded transition costs over shared Hello predecessor IDs.

The matched later-row comparison uses the corrected disk-bit7 2013 capture.
Each engine's recorded predecessor-path cost is retained, so the output does
not claim a shared dynamic-programming history or whole-chain parity.
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parents[4]
CAPTURES = {
    "2006": ROOT
    / "tools/revkit/work/corpus-parity/stage20/hello-repeat-2006-pairmatrix/legacy-pairmatrix-gdb.log",
    "2013": ROOT
    / "tools/revkit/work/corpus-parity/stage20/hello-repeat-2013-pairmatrix/adapted-pairmatrix-gdb.log",
}
CURRENT_UNIT = 272823
POSITION = 1
LINE = re.compile(
    r"(?:LEGACY|NEW)_PAIR_COST\s+"
    r"(?:position|context)=(?P<position>\d+)\s+"
    r"current=(?P<current>\d+)\s+previous=(?P<previous>\d+)\s+"
    r"(?:transition|feature)=(?P<feature>[0-9.]+).*?"
    r"total=(?P<total>[0-9.]+)"
)
OLD_MATCHED_LINE = re.compile(
    r"LEGACY_MATCHED_EDGE context=(?P<context>\d+) "
    r"current=(?P<current>\d+) previous=(?P<previous>\d+) "
    r"transition=(?P<feature>[0-9.]+).*?total=(?P<total>[0-9.]+)"
)
NEW_MATCHED_LINE = re.compile(
    r"(?:MATCHED_EDGE|BIT7_MATCHED_EDGE) context=(?P<context>\d+) "
    r"current=(?P<current>\d+) previous=(?P<previous>\d+) "
    r"feature=(?P<feature>[0-9.]+).*?total=(?P<total>[0-9.]+)"
)
LEGACY_NEXT_LINE = re.compile(
    r"LEGACY_NEXT_EDGE context=(?P<context>\d+) "
    r"current=(?P<current>\d+) previous=(?P<previous>\d+) "
    r"transition=(?P<feature>[0-9.]+).*?total=(?P<total>[0-9.]+)"
)
NEW_NEXT_LINE = re.compile(
    r"BIT7_NEXT_EDGE context=(?P<context>\d+) "
    r"current=(?P<current>\d+) previous=(?P<previous>\d+) "
    r"feature=(?P<feature>[0-9.]+).*?total=(?P<total>[0-9.]+)"
)


@dataclass(frozen=True)
class EdgeCost:
    feature: float
    total: float


def load_edges(path: Path) -> dict[int, EdgeCost]:
    if not path.is_file():
        raise FileNotFoundError(path)

    edges: dict[int, EdgeCost] = {}
    for line in path.read_text(encoding="utf-8", errors="replace").splitlines():
        match = LINE.search(line)
        if match is None:
            continue
        if int(match["position"]) != POSITION:
            continue
        if int(match["current"]) != CURRENT_UNIT:
            continue

        previous = int(match["previous"])
        edges[previous] = EdgeCost(
            feature=float(match["feature"]),
            total=float(match["total"]),
        )
    if not edges:
        raise ValueError(f"no matching position/current row in {path}")
    return edges


def load_matched_rows(path: Path, pattern: re.Pattern[str]) -> dict[int, dict[int, EdgeCost]]:
    if not path.is_file():
        raise FileNotFoundError(path)

    rows: dict[int, dict[int, EdgeCost]] = {}
    for line in path.read_text(encoding="utf-8", errors="replace").splitlines():
        match = pattern.search(line)
        if match is None:
            continue
        current = int(match["current"])
        previous = int(match["previous"])
        rows.setdefault(current, {})[previous] = EdgeCost(
            feature=float(match["feature"]),
            total=float(match["total"]),
        )
    return rows


def load_one_row(
    path: Path,
    pattern: re.Pattern[str],
    *,
    context: int,
    current: int,
) -> dict[int, EdgeCost]:
    if not path.is_file():
        raise FileNotFoundError(path)

    edges: dict[int, EdgeCost] = {}
    for line in path.read_text(encoding="utf-8", errors="replace").splitlines():
        match = pattern.search(line)
        if match is None:
            continue
        if int(match["context"]) != context or int(match["current"]) != current:
            continue
        edges[int(match["previous"])] = EdgeCost(
            feature=float(match["feature"]),
            total=float(match["total"]),
        )
    if not edges:
        raise ValueError(f"no matching context/current row in {path}")
    return edges


def main() -> None:
    costs = {generation: load_edges(path) for generation, path in CAPTURES.items()}
    common = set(costs["2006"]) & set(costs["2013"])
    if not common:
        raise ValueError("the captures have no common predecessor IDs")

    print(
        f"position={POSITION} current={CURRENT_UNIT} "
        f"old={len(costs['2006'])} new={len(costs['2013'])} shared={len(common)}"
    )
    for generation in ("2006", "2013"):
        edges = costs[generation]
        best = min(common, key=lambda unit: edges[unit].total)
        print(
            f"{generation} best_shared_previous={best} "
            f"feature={edges[best].feature:g} total={edges[best].total:g}"
        )

    for previous in (272822, 273369):
        if previous not in common:
            raise ValueError(f"expected comparison edge {previous} is not shared")
        old = costs["2006"][previous]
        new = costs["2013"][previous]
        print(
            f"previous={previous} "
            f"old_feature={old.feature:g} old_total={old.total:g} "
            f"new_feature={new.feature:g} new_total={new.total:g}"
        )

    matched_paths = {
        "2006": ROOT
        / "tools/revkit/work/corpus-parity/stage20/hello-repeat-2006-matched-late-edges/legacy-matched-late-edges-gdb.log",
        "2013": ROOT
        / "tools/revkit/work/corpus-parity/stage20/hello-repeat-2013-bit7-matched-late-edges/adapted-matched-edges-gdb.log",
    }
    matched = {
        "2006": load_matched_rows(matched_paths["2006"], OLD_MATCHED_LINE),
        "2013": load_matched_rows(matched_paths["2013"], NEW_MATCHED_LINE),
    }
    print("matched current-unit rows (context ordinals differ between engines):")
    for current in (232670, 266023, 264071):
        old_edges = matched["2006"].get(current, {})
        new_edges = matched["2013"].get(current, {})
        shared_previous = set(old_edges) & set(new_edges)
        print(
            f"current={current} old={len(old_edges)} new={len(new_edges)} "
            f"shared_previous={len(shared_previous)}"
        )
        if shared_previous:
            for generation, edges in (("2006", old_edges), ("2013", new_edges)):
                best = min(shared_previous, key=lambda unit: edges[unit].total)
                print(
                    f" {generation} best_shared_previous={best} "
                    f"feature={edges[best].feature:g} total={edges[best].total:g}"
                )

    next_row_paths = {
        "2006": ROOT
        / "tools/revkit/work/corpus-parity/stage20/hello-repeat-2006-next-row-edge/legacy-next-row-edge-gdb.log",
        "2013": ROOT
        / "tools/revkit/work/corpus-parity/stage20/hello-repeat-2013-bit7-next-row-edge/adapted-next-row-edge-gdb.log",
    }
    next_edges = {
        "2006": load_one_row(
            next_row_paths["2006"], LEGACY_NEXT_LINE, context=6, current=264072
        ),
        "2013": load_one_row(
            next_row_paths["2013"], NEW_NEXT_LINE, context=7, current=264072
        ),
    }
    shared_next = set(next_edges["2006"]) & set(next_edges["2013"])
    print(
        f"current=264072 old={len(next_edges['2006'])} "
        f"new={len(next_edges['2013'])} shared_previous={len(shared_next)}"
    )
    for generation in ("2006", "2013"):
        edges = next_edges[generation]
        best = min(shared_next, key=lambda unit: edges[unit].total)
        print(
            f" {generation} best_shared_previous={best} "
            f"feature={edges[best].feature:g} total={edges[best].total:g}"
        )


if __name__ == "__main__":
    main()
