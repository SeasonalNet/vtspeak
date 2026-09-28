#!/usr/bin/env python3
"""Experimental adapter for the 2006 special vector pitch-tree records."""

from __future__ import annotations

import struct
import sys
from dataclasses import dataclass
from pathlib import Path
from shutil import copytree

ROOT = Path(__file__).resolve().parents[4]
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))
from tree3 import parse_tree_data  # noqa: E402

SOURCE = ROOT / "data-kate" / "M16" / "ttsdata" / "tree" / "pitch"
OUTPUT = Path(__file__).resolve().parent / "vector-pitch-overlay"
CURRENT_PITCH = (
    ROOT / "tools" / "revkit" / "work" / "stage19" / "data-kate-copy"
    / "M16" / "ttsdata" / "tree2" / "pitch"
)
NAMES = ("nbf.tree2", "bf.tree2", "qbf.tree2", "sbf.tree2")
KNOWN_ROWS = {
    "nbf.tree2": (
        [2, 11, 0, 3, 1, 1, 0, 3, 1, 1, 2, 0, 1073, 4098, 6960, 398,
         32, 623, 6960, 398, 32, 623, 4304, 416, 1918, 4098, 6960, 398,
         32, 623, 1, 0],
        (85, 84, 83, 83, 84, 84, 85, 85, 86, 86, 86, 86),
    ),
    "sbf.tree2": (
        [5, 25, 1, 0, 3, 3, 1, 4, 3, 1, 2, 0, 1073, 4098, 6960, 398,
         32, 623, 6960, 398, 32, 623, 4304, 416, 1918, 4098, 6960, 398,
         32, 623, 1, 0],
        (96, 98, 100, 104, 108, 112, 116, 120, 122, 124, 119, 113),
    ),
}


@dataclass(frozen=True)
class Node:
    feature: int
    operation: str
    threshold: int
    values: tuple[int, ...]
    output: tuple[int, ...]
    children: tuple[Node, Node] | None


def parse(path: Path) -> tuple[Node, int, int]:
    raw = path.read_bytes()
    offset = 0
    output_widths: set[int] = set()
    node_count = 0

    def require(size: int, label: str) -> None:
        if offset + size > len(raw):
            raise ValueError(f"{path}: truncated {label} at byte {offset}")

    def read_node(depth: int = 0) -> Node:
        nonlocal offset, node_count
        if depth > 900:
            raise ValueError(f"{path}: tree depth exceeds 900")
        require(5, "special node header")
        feature, operation_byte, threshold, width = struct.unpack_from(
            "<BBhB", raw, offset
        )
        offset += 5
        operation = chr(operation_byte)
        if operation not in ("C", "D"):
            raise ValueError(f"{path}: unknown operation {operation!r}")
        output_widths.add(width)
        require(width * 2 + 1, "output vector and value count")
        output = struct.unpack_from(f"<{width}h", raw, offset) if width else ()
        offset += width * 2
        value_count = raw[offset]
        offset += 1
        require(value_count * 2 + 2, "decision values and child marker")
        values = (
            struct.unpack_from(f"<{value_count}h", raw, offset)
            if value_count
            else ()
        )
        offset += value_count * 2
        _reserved, has_children = raw[offset : offset + 2]
        offset += 2
        if operation == "C" and values:
            raise ValueError(f"{path}: comparison node has decision values")
        children = (read_node(depth + 1), read_node(depth + 1)) if has_children else None
        node_count += 1
        return Node(feature, operation, threshold, tuple(values), tuple(output), children)

    root = read_node()
    if offset != len(raw):
        raise ValueError(f"{path}: parsed through byte {offset} of {len(raw)}")
    if len(output_widths) != 1 or 0 in output_widths:
        raise ValueError(f"{path}: inconsistent or empty vector widths {output_widths}")
    return root, node_count, output_widths.pop()


def encode(root: Node, node_count: int, output_width: int) -> bytes:
    indexed: list[Node] = []
    outputs: list[tuple[int, ...]] = []
    edges: dict[int, tuple[int, int]] = {}

    def emit(node: Node) -> int:
        if node.children is None:
            if len(node.output) != output_width:
                raise ValueError("leaf output width changed during parse")
            leaf = len(outputs)
            outputs.append(node.output)
            return -leaf - 1
        index = len(indexed)
        indexed.append(node)
        first, second = node.children
        edges[index] = (emit(first), emit(second))
        return index

    emit(root)
    if len(indexed) + len(outputs) != node_count or len(outputs) != len(indexed) + 1:
        raise ValueError("tree is not a complete binary tree")
    list_values = sum(len(node.values) for node in indexed)
    if len(indexed) > 32767 or list_values > 0xFFFFFFFF:
        raise ValueError("tree exceeds tree3 header limits")

    raw = bytearray(struct.pack("<hBI", len(indexed), output_width, list_values))
    for index, node in enumerate(indexed):
        if len(node.values) > 255:
            raise ValueError("decision list exceeds byte count")
        raw.extend(
            struct.pack(
                "<BBhB", node.feature, ord(node.operation), node.threshold, len(node.values)
            )
        )
        if node.values:
            raw.extend(struct.pack(f"<{len(node.values)}h", *node.values))
        raw.extend(struct.pack("<hh", *edges[index]))
    for row in outputs:
        raw.extend(struct.pack(f"<{output_width}h", *row))

    tree, end = parse_tree_data(Path("<experimental-vector-tree3>"), bytes(raw))
    if end != len(raw) or tree.output_width != output_width:
        raise ValueError("tree3 round-trip failed")
    return bytes(raw)


def evaluate(root: Node, features: list[int]) -> tuple[int, ...]:
    node = root
    while node.children is not None:
        value = features[node.feature]
        take_first = (
            value in node.values
            if node.operation == "D"
            else value <= node.threshold
        )
        node = node.children[0 if take_first else 1]
    return node.output


def main() -> None:
    destination = OUTPUT / "pitch"
    copytree(CURRENT_PITCH, destination, dirs_exist_ok=True)
    for name in NAMES:
        source = SOURCE / name
        root, nodes, width = parse(source)
        converted = encode(root, nodes, width)
        tree, _ = parse_tree_data(Path(name), converted)
        if name in KNOWN_ROWS:
            features, expected = KNOWN_ROWS[name]
            actual_source = evaluate(root, features)
            _leaf, actual_indexed = tree.evaluate(features)
            if actual_source != expected or actual_indexed != expected:
                raise ValueError(f"{name}: known Stage 20 vector row mismatch")
        (destination / name).write_bytes(converted)
        print(
            f"{name} source_bytes={source.stat().st_size} parsed_nodes={nodes} "
            f"output_width={width} tree3_bytes={len(converted)} "
            f"tree3_nodes={len(tree.nodes)} leaves={len(tree.outputs)}"
        )


if __name__ == "__main__":
    main()
