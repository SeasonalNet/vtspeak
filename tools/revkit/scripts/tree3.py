#!/usr/bin/env python3
"""Read-only parser and evaluator for the observed VoiceText tree3 format."""

from __future__ import annotations

import argparse
import json
import struct
from collections import Counter
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class Node:
    feature: int
    operation: str
    threshold: int
    values: tuple[int, ...]
    when_true: int
    when_false: int


@dataclass(frozen=True)
class Tree:
    path: Path
    output_width: int
    nodes: tuple[Node, ...]
    outputs: tuple[tuple[int, ...], ...]
    declared_list_values: int

    def evaluate(self, features: list[int]) -> tuple[int, tuple[int, ...]]:
        required_features = max((node.feature + 1 for node in self.nodes), default=0)
        if len(features) < required_features:
            raise ValueError(
                f"{self.path}: expected at least {required_features} input features, "
                f"got {len(features)}"
            )
        node_index = 0
        steps = 0
        while node_index >= 0:
            if node_index >= len(self.nodes):
                raise ValueError(f"{self.path}: child node index {node_index} is out of range")
            node = self.nodes[node_index]
            feature_value = features[node.feature]
            if node.operation == "D":
                next_index = node.when_true if feature_value in node.values else node.when_false
            else:
                next_index = node.when_true if feature_value <= node.threshold else node.when_false
            node_index = next_index
            steps += 1
            if steps > len(self.nodes) + 1:
                raise ValueError(f"{self.path}: cycle encountered during tree evaluation")

        leaf = -node_index - 1
        if leaf >= len(self.outputs):
            raise ValueError(f"{self.path}: leaf ordinal {leaf} is out of range")
        return leaf, self.outputs[leaf]


def parse_tree(path: Path) -> Tree:
    raw = path.read_bytes()
    if len(raw) < 7:
        raise ValueError(f"{path}: file is shorter than the 7-byte header")

    node_count, output_width, declared_list_values = struct.unpack_from("<hBI", raw, 0)
    if node_count < 0:
        raise ValueError(f"{path}: negative node count {node_count}")
    if output_width == 0:
        raise ValueError(f"{path}: output width must be positive")
    offset = 7
    nodes: list[Node] = []
    list_value_count = 0

    def require(size: int, label: str) -> None:
        if offset + size > len(raw):
            raise ValueError(f"{path}: truncated {label} at byte {offset}")

    for index in range(node_count):
        require(5, f"node {index} header")
        feature, operation_byte, threshold, value_count = struct.unpack_from("<BBhB", raw, offset)
        offset += 5
        require(value_count * 2 + 4, f"node {index} list and children")
        values = struct.unpack_from(f"<{value_count}h", raw, offset) if value_count else ()
        offset += value_count * 2
        when_true, when_false = struct.unpack_from("<hh", raw, offset)
        offset += 4
        try:
            operation = chr(operation_byte)
        except ValueError as exc:
            raise ValueError(f"{path}: invalid operation byte {operation_byte}") from exc
        if operation not in ("C", "D"):
            raise ValueError(f"{path}: node {index} has unknown operation {operation!r}")
        if operation == "C" and value_count:
            raise ValueError(
                f"{path}: node {index} operation {operation!r} has an unexpected list"
            )
        nodes.append(Node(feature, operation, threshold, tuple(values), when_true, when_false))
        list_value_count += value_count

    if list_value_count != declared_list_values:
        raise ValueError(
            f"{path}: header declares {declared_list_values} list values, "
            f"nodes contain {list_value_count}"
        )

    for index, node in enumerate(nodes):
        for child in (node.when_true, node.when_false):
            if child >= node_count or (child < 0 and -child - 1 > node_count):
                raise ValueError(f"{path}: node {index} has invalid child/leaf reference {child}")

    output_count = (node_count + 1) * output_width
    output_bytes = output_count * 2
    if offset + output_bytes != len(raw):
        raise ValueError(
            f"{path}: expected {output_bytes} output bytes at offset {offset}, "
            f"file has {len(raw) - offset} remaining"
        )
    flat_outputs = struct.unpack_from(f"<{output_count}h", raw, offset) if output_count else ()
    outputs = tuple(
        tuple(flat_outputs[index * output_width : (index + 1) * output_width])
        for index in range(node_count + 1)
    )
    return Tree(path, output_width, tuple(nodes), outputs, declared_list_values)


def describe(tree: Tree) -> dict[str, object]:
    values = [value for row in tree.outputs for value in row]
    return {
        "file": str(tree.path),
        "nodes": len(tree.nodes),
        "output_width": tree.output_width,
        "list_values": tree.declared_list_values,
        "output_rows": len(tree.outputs),
        "operations": dict(sorted(Counter(node.operation for node in tree.nodes).items())),
        "feature_selectors": sorted({node.feature for node in tree.nodes}),
        "threshold_range": [
            min((node.threshold for node in tree.nodes), default=0),
            max((node.threshold for node in tree.nodes), default=0),
        ],
        "output_range": [min(values, default=0), max(values, default=0)],
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("files", nargs="+", type=Path, help="tree3 files to parse")
    parser.add_argument(
        "--features",
        help="comma-separated signed feature values; evaluate each tree with this vector",
    )
    args = parser.parse_args()
    feature_values = None
    if args.features is not None:
        feature_values = [int(value, 0) for value in args.features.split(",") if value]

    results = []
    for path in args.files:
        tree = parse_tree(path)
        result = describe(tree)
        if feature_values is not None:
            leaf, output = tree.evaluate(feature_values)
            result["evaluation"] = {"leaf": leaf, "output": list(output)}
        results.append(result)
    print(json.dumps(results, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
