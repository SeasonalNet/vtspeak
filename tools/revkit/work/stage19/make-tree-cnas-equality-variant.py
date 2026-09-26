#!/usr/bin/env python3
"""Reverse only cnas comparison nodes hit at equality in the address trace."""

from __future__ import annotations

import hashlib
import shutil
import struct
import sys
from pathlib import Path

WORK = Path(__file__).resolve().parent
SOURCE = WORK / "data-kate-copy" / "M16" / "ttsdata" / "tree2"
OUTPUT = WORK / "tree-branch-variants" / "invert-cnas-equality-address"
TARGET_NODES = {39, 111, 221}
sys.path.insert(0, str(Path(__file__).resolve().parents[3] / "revkit" / "scripts"))

from tree3 import parse_tree_data  # noqa: E402


def main() -> None:
    shutil.copytree(SOURCE, OUTPUT, dirs_exist_ok=True)
    path = OUTPUT / "duration" / "cnas.tree2"
    raw = path.read_bytes()
    changed = bytearray(raw)
    node_count = struct.unpack_from("<h", raw)[0]
    offset = 7
    swapped: set[int] = set()
    for index in range(node_count):
        feature, operation, _threshold, value_count = struct.unpack_from("<BBhB", raw, offset)
        child_offset = offset + 5 + value_count * 2
        if index in TARGET_NODES:
            if operation != ord("C"):
                raise ValueError(f"node {index} is not a C comparison node")
            first, second = struct.unpack_from("<hh", raw, child_offset)
            struct.pack_into("<hh", changed, child_offset, second, first)
            swapped.add(index)
        offset = child_offset + 4
    if swapped != TARGET_NODES:
        raise ValueError(f"expected to swap {sorted(TARGET_NODES)}, swapped {sorted(swapped)}")
    path.write_bytes(changed)
    tree, end = parse_tree_data(path, changed)
    if end != len(changed) or len(tree.nodes) != node_count:
        raise ValueError("modified cnas tree failed structural validation")
    manifest = [
        "variant=invert-cnas-equality-address",
        "source=duration/cnas.tree2",
        "target_nodes=39,111,221",
        "selection_basis=address-fixture runtime trace showed feature value equal to threshold",
        f"node_count={node_count}",
        f"source_sha256={hashlib.sha256(raw).hexdigest()}",
        f"variant_sha256={hashlib.sha256(changed).hexdigest()}",
    ]
    (OUTPUT / "variant-manifest.txt").write_text("\n".join(manifest) + "\n")
    compose = WORK / "compose-tree-invert-cnas-equality-address.yaml"
    compose.write_text(
        "services:\n"
        "  runtime:\n"
        "    volumes:\n"
        "      - ../stage19/tree-branch-variants/invert-cnas-equality-address:/work/data-kate/M16/ttsdata/tree2:ro\n"
    )
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
