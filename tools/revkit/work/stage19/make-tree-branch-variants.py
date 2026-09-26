#!/usr/bin/env python3
"""Make isolated converted-tree variants with selected branch pairs reversed."""

from __future__ import annotations

import hashlib
import shutil
import struct
import sys
from pathlib import Path

WORK = Path(__file__).resolve().parent
SOURCE = WORK / "data-kate-copy" / "M16" / "ttsdata" / "tree2"
OUTPUT = WORK / "tree-branch-variants"
sys.path.insert(0, str(Path(__file__).resolve().parents[3] / "revkit" / "scripts"))

from tree3 import parse_tree_data  # noqa: E402

VARIANTS: dict[str, tuple[set[int], str | None, str | None]] = {
    "invert-c": ({ord("C")}, None, None),
    "invert-d": ({ord("D")}, None, None),
    "invert-all": ({ord("C"), ord("D")}, None, None),
    "invert-c-duration": ({ord("C")}, "duration", None),
    "invert-c-pitch": ({ord("C")}, "pitch", None),
}
VARIANTS.update(
    {
        f"invert-c-duration-{name}": ({ord("C")}, "duration", f"{name}.tree2")
        for name in ("caff", "capp", "cfri", "cnas", "cstop", "vdi", "vlong", "vsch", "vshort")
    }
)


def reverse_branches(raw: bytes, operations: set[int]) -> tuple[bytes, int]:
    node_count = struct.unpack_from("<h", raw)[0]
    if node_count < 0:
        raise ValueError(f"negative converted-tree node count: {node_count}")
    changed = bytearray(raw)
    offset = 7
    swaps = 0
    for _ in range(node_count):
        if offset + 5 > len(raw):
            raise ValueError("truncated converted-tree node header")
        operation = raw[offset + 1]
        value_count = raw[offset + 4]
        child_offset = offset + 5 + value_count * 2
        if child_offset + 4 > len(raw):
            raise ValueError("truncated converted-tree child references")
        if operation in operations:
            first, second = struct.unpack_from("<hh", raw, child_offset)
            struct.pack_into("<hh", changed, child_offset, second, first)
            swaps += 1
        offset = child_offset + 4
    return bytes(changed), swaps


def main() -> None:
    sources = sorted(SOURCE.rglob("*.tree2"))
    if len(sources) != 17:
        raise SystemExit(f"expected 17 converted trees under {SOURCE}, found {len(sources)}")

    for label, (operations, tree_family, tree_name) in VARIANTS.items():
        destination_root = OUTPUT / label
        shutil.copytree(SOURCE, destination_root, dirs_exist_ok=True)
        selected_sources = [
            source
            for source in sources
            if (tree_family is None or source.relative_to(SOURCE).parts[0] == tree_family)
            and (tree_name is None or source.name == tree_name)
        ]
        rows = [f"variant={label}", f"reversed_trees={len(selected_sources)}"]
        for source in selected_sources:
            original = source.read_bytes()
            changed, swaps = reverse_branches(original, operations)
            relative = source.relative_to(SOURCE)
            destination = destination_root / relative
            destination.write_bytes(changed)
            tree, end = parse_tree_data(destination, changed)
            if end != len(changed) or len(tree.nodes) != struct.unpack_from("<h", changed)[0]:
                raise ValueError(f"{destination}: converted tree failed structural validation")
            rows.append(
                f"{relative} nodes={len(tree.nodes)} swapped={swaps} "
                f"source_sha256={hashlib.sha256(original).hexdigest()} "
                f"variant_sha256={hashlib.sha256(changed).hexdigest()}"
            )
        (destination_root / "variant-manifest.txt").write_text("\n".join(rows) + "\n")
        if tree_name is not None:
            (WORK / f"compose-tree-{label}.yaml").write_text(
                "services:\n"
                "  runtime:\n"
                "    volumes:\n"
                f"      - ../stage19/tree-branch-variants/{label}:/work/data-kate/M16/ttsdata/tree2:ro\n"
            )
        print("\n".join(rows))


if __name__ == "__main__":
    main()
