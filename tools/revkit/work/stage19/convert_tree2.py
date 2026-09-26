#!/usr/bin/env python3
"""Convert observed legacy tree2 trees to the indexed tree3 layout."""

from __future__ import annotations

import random
import shutil
import struct
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))
from tree2 import Node as LegacyNode, parse_tree2  # noqa: E402
from tree3 import parse_tree_data  # noqa: E402


def encode(root: LegacyNode) -> bytes:
    indexed: list[LegacyNode] = []
    outputs: list[int] = []
    edges: dict[int, tuple[int, int]] = {}

    def emit(node: LegacyNode) -> int:
        if node.children is None:
            leaf = len(outputs)
            outputs.append(node.output)
            return -leaf - 1
        index = len(indexed)
        indexed.append(node)
        first, second = node.children
        edges[index] = (emit(first), emit(second))
        return index

    emit(root)

    list_values = sum(len(node.values) for node in indexed)
    if len(indexed) > 32767 or list_values > 0xFFFFFFFF:
        raise ValueError("tree exceeds tree3 header limits")
    raw = bytearray(struct.pack("<hBI", len(indexed), 1, list_values))
    for index, node in enumerate(indexed):
        if node.operation == "C" and node.values:
            raise ValueError("comparison node unexpectedly has list values")
        if len(node.values) > 255:
            raise ValueError("node list exceeds byte count")
        raw.extend(struct.pack("<BBhB", node.feature, ord(node.operation), node.threshold, len(node.values)))
        if node.values:
            raw.extend(struct.pack(f"<{len(node.values)}h", *node.values))
        raw.extend(struct.pack("<hh", *edges[index]))
    raw.extend(struct.pack(f"<{len(outputs)}h", *outputs))
    return bytes(raw)


def verify(source: Path, encoded: bytes, root: LegacyNode) -> None:
    tree, end = parse_tree_data(source, encoded)
    if end != len(encoded):
        raise ValueError(f"{source}: converted tree has trailing bytes")
    max_feature = 0

    def visit(node: LegacyNode) -> None:
        nonlocal max_feature
        if node.children is not None:
            max_feature = max(max_feature, node.feature)
            visit(node.children[0])
            visit(node.children[1])

    visit(root)
    if not tree.nodes:
        if tree.outputs != ((root.output,),):
            raise ValueError(f"{source}: converted leaf output differs")
        return
    rng = random.Random(0x54524545)
    for _ in range(128):
        features = [rng.randint(-32768, 32767) for _ in range(max_feature + 1)]
        old_value = root
        while old_value.children is not None:
            value = features[old_value.feature]
            first = value in old_value.values if old_value.operation == "D" else value <= old_value.threshold
            old_value = old_value.children[0 if first else 1]
        _leaf, row = tree.evaluate(features)
        if row != (old_value.output,):
            raise ValueError(f"{source}: converted evaluation differs for {features}")


def main() -> None:
    source_root = ROOT / "data-kate" / "M16" / "ttsdata" / "tree"
    package_copy = WORK / "data-kate-copy"
    if not package_copy.exists():
        shutil.copytree(ROOT / "data-kate", package_copy)
    output_root = package_copy / "M16" / "ttsdata" / "tree2"
    files = sorted(source_root.rglob("*.tree2"))
    if len(files) != 17:
        raise SystemExit(f"expected 17 Kate tree2 files, found {len(files)}")
    rows = ["format=tree2 recursive to tree3 indexed bridge", f"source_files={len(files)}"]
    for source in files:
        legacy = parse_tree2(source)
        converted = encode(legacy.root)
        verify(source, converted, legacy.root)
        relative = source.relative_to(source_root)
        destination = output_root / relative
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_bytes(converted)
        rows.append(
            f"{relative} source={source.stat().st_size} converted={len(converted)} "
            f"parsed={legacy.end_offset} trailing={legacy.trailing_bytes} "
            f"nodes={legacy.nodes}"
        )
    report = WORK / "conversion-manifest.txt"
    report.write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
