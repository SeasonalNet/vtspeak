#!/usr/bin/env python3
"""Read-only parser/evaluator for the observed recursive VoiceText tree2 nodes."""

from __future__ import annotations

import argparse
import json
import struct
from collections import Counter
from dataclasses import dataclass
from pathlib import Path

MAX_NODES = 100_000
MAX_DEPTH = 900


@dataclass(frozen=True)
class Node:
    feature: int
    operation: str
    threshold: int
    output: int
    values: tuple[int, ...]
    trailer_bytes: tuple[int, int]
    children: tuple[Node, Node] | None


@dataclass(frozen=True)
class Tree:
    path: Path
    root: Node
    nodes: int
    leaves: int
    end_offset: int
    file_size: int

    @property
    def trailing_bytes(self) -> int:
        return self.file_size - self.end_offset

    def evaluate(self, features: list[int]) -> int:
        node = self.root
        while node.children is not None:
            if node.feature >= len(features):
                raise ValueError(
                    f"{self.path}: expected feature {node.feature}, "
                    f"got {len(features)} inputs"
                )
            value = features[node.feature]
            if node.operation == "D":
                take_first = value in node.values
            else:
                take_first = value <= node.threshold
            node = node.children[0 if take_first else 1]
        return node.output


def parse_tree2_data(path: Path, raw: bytes) -> Tree:
    position = 0
    node_count = 0
    leaf_count = 0

    def require(size: int, label: str) -> None:
        if position + size > len(raw):
            raise ValueError(
                f"{path}: truncated {label} at byte {position}; "
                f"need {size}, have {len(raw) - position}"
            )

    def parse_node(depth: int) -> Node:
        nonlocal position, node_count, leaf_count
        if depth > MAX_DEPTH:
            raise ValueError(f"{path}: tree depth exceeds safety limit {MAX_DEPTH}")
        if node_count >= MAX_NODES:
            raise ValueError(f"{path}: node count exceeds safety limit {MAX_NODES}")
        require(7, "node header")
        feature, operation_byte, threshold, output, value_count = struct.unpack_from(
            "<BBhhB", raw, position
        )
        node_offset = position
        position += 7
        if operation_byte not in (ord("C"), ord("D")):
            raise ValueError(
                f"{path}: node at byte {node_offset} has operation byte "
                f"0x{operation_byte:02x}"
            )
        require(value_count * 2 + 2, "node list and trailer")
        values = (
            struct.unpack_from(f"<{value_count}h", raw, position)
            if value_count
            else ()
        )
        position += value_count * 2
        trailer_bytes = (raw[position], raw[position + 1])
        position += 2
        children = None
        if trailer_bytes[1] != 0:
            left = parse_node(depth + 1)
            right = parse_node(depth + 1)
            children = (left, right)
        else:
            leaf_count += 1
        node_count += 1
        return Node(
            feature=feature,
            operation=chr(operation_byte),
            threshold=threshold,
            output=output,
            values=tuple(values),
            trailer_bytes=trailer_bytes,
            children=children,
        )

    root = parse_node(0)
    return Tree(path, root, node_count, leaf_count, position, len(raw))


def parse_tree2(path: Path, *, require_eof: bool = False) -> Tree:
    tree = parse_tree2_data(path, path.read_bytes())
    if require_eof and tree.trailing_bytes:
        raise ValueError(
            f"{path}: parsed tree ends at {tree.end_offset}, "
            f"file has {tree.trailing_bytes} trailing bytes"
        )
    return tree


def describe(tree: Tree) -> dict[str, object]:
    operations: Counter[str] = Counter()
    trailers: Counter[str] = Counter()
    feature_selectors: set[int] = set()
    outputs: list[int] = []
    list_values = 0
    max_depth = 0

    def visit(node: Node, depth: int) -> None:
        nonlocal list_values, max_depth
        operations[node.operation] += 1
        trailers[f"{node.trailer_bytes[0]:02x} {node.trailer_bytes[1]:02x}"] += 1
        feature_selectors.add(node.feature)
        outputs.append(node.output)
        list_values += len(node.values)
        max_depth = max(max_depth, depth)
        if node.children is not None:
            visit(node.children[0], depth + 1)
            visit(node.children[1], depth + 1)

    visit(tree.root, 0)
    return {
        "file": str(tree.path),
        "file_size_bytes": tree.file_size,
        "parsed_tree_bytes": tree.end_offset,
        "trailing_bytes": tree.trailing_bytes,
        "nodes": tree.nodes,
        "leaves": tree.leaves,
        "max_depth": max_depth,
        "operations": dict(sorted(operations.items())),
        "feature_selectors": sorted(feature_selectors),
        "list_value_count": list_values,
        "leaf_output_range": [min(outputs), max(outputs)],
        "trailer_bytes": dict(sorted(trailers.items())),
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("files", nargs="+", type=Path, help="tree2 files to parse")
    parser.add_argument(
        "--strict-eof",
        action="store_true",
        help="fail if a parsed tree leaves any file suffix unconsumed",
    )
    parser.add_argument(
        "--features",
        help="optional comma-separated signed feature values to evaluate",
    )
    args = parser.parse_args()
    features = None
    if args.features is not None:
        features = [int(value, 0) for value in args.features.split(",") if value]
    results = []
    for path in args.files:
        tree = parse_tree2(path, require_eof=args.strict_eof)
        result = describe(tree)
        if features is not None:
            result["evaluation"] = {"output": tree.evaluate(features)}
        results.append(result)
    print(json.dumps(results, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
