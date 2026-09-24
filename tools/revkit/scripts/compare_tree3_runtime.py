#!/usr/bin/env python3
"""Cross-check captured runtime tree lookups against local tree3 resources."""

from __future__ import annotations

import argparse
import re
from collections import Counter, defaultdict
from pathlib import Path

from dict_resources import parse_atmt_trees
from tree3 import Tree, parse_tree


def values_after_colon(line: str) -> list[int]:
    return [int(value) for value in re.findall(r"-?\d+", line.split(":", 1)[1])]


def load_trees(root: Path) -> list[Tree]:
    common = root / "data-common/dict-eng"
    paul = root / "data-paul/M16/ttsdata/tree3"
    trees = [parse_tree(common / name) for name in ("engbi.tree3", "poly.tree3", "sbd.tree3")]
    trees.extend(parse_atmt_trees(common / "atmt.tree3"))
    trees.extend(parse_tree(path) for path in sorted(paul.rglob("*.tree3")))
    return trees


def parse_lookup(
    lines: list[str], index: int
) -> tuple[tuple[int, ...], tuple[int, ...], tuple[int, int]]:
    features = values_after_colon(lines[index + 1]) + values_after_colon(lines[index + 2])
    header = lines[index]
    scalar = header.startswith("SCALAR_ENTRY")
    result_prefix = "SCALAR_RETURN" if scalar else "VECTOR_RETURN values:"
    result_index = next(
        position
        for position in range(index + 3, min(index + 9, len(lines)))
        if lines[position].startswith(result_prefix)
    )
    if scalar:
        result = re.search(r"value=(-?\d+)", lines[result_index])
        if result is None:
            raise ValueError(f"malformed scalar result: {lines[result_index]}")
        outputs = (int(result.group(1)),)
    else:
        outputs = tuple(
            int(value) for value in re.findall(r"-?\d+", lines[result_index].rsplit(":", 1)[1])
        )
        outputs += tuple(values_after_colon(lines[result_index + 1]))
    shape = re.search(r"nodes=(\d+) output_width=(\d+)", header)
    if shape is None:
        raise ValueError(f"malformed tree lookup header: {header}")
    return features, outputs, (int(shape.group(1)), int(shape.group(2)))


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("logs", nargs="+", type=Path, help="GDB lookup logs")
    parser.add_argument("--root", type=Path, default=Path.cwd(), help="repository root")
    args = parser.parse_args()

    by_shape: dict[tuple[int, int], list[Tree]] = defaultdict(list)
    for tree in load_trees(args.root):
        by_shape[(len(tree.nodes), tree.output_width)].append(tree)

    compared: Counter[str] = Counter()
    skipped: Counter[str] = Counter()
    for log_path in args.logs:
        lines = log_path.read_text(errors="replace").splitlines()
        for index, line in enumerate(lines):
            if not (line.startswith("SCALAR_ENTRY") or line.startswith("VECTOR_ENTRY")):
                continue
            features, actual, shape = parse_lookup(lines, index)
            candidates = by_shape.get(shape, [])
            matches = []
            for tree in candidates:
                try:
                    _, expected = tree.evaluate(list(features))
                except ValueError:
                    continue
                if expected == actual:
                    matches.append(tree)
            if not matches:
                skipped["no_unique_resource_match"] += 1
                continue
            # A shape collision can still be distinguished when only one resource
            # predicts the returned value for the captured feature vector.
            if len(matches) != 1:
                skipped["ambiguous_resource_match"] += 1
                continue
            compared[matches[0].path.name] += 1

    print(f"exact resource matches: {sum(compared.values())}")
    print(f"unmatched or ambiguous: {sum(skipped.values())}")
    for name, count in sorted(compared.items()):
        print(f"{name}: {count}")
    for name, count in sorted(skipped.items()):
        print(f"{name}: {count}")
    return int(bool(skipped))


if __name__ == "__main__":
    raise SystemExit(main())
