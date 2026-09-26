#!/usr/bin/env python3
"""Build a disposable variant that preserves opaque tree2 pitch suffixes."""

from __future__ import annotations

import sys
from pathlib import Path

sys.dont_write_bytecode = True

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = ROOT / "data-kate" / "M16" / "ttsdata" / "tree"
CONVERTED = WORK / "data-kate-copy" / "M16" / "ttsdata" / "tree2"
OUTPUT = WORK / "tree2-suffix"
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))
from tree2 import parse_tree2  # noqa: E402


def main() -> None:
    rows = ["variant=append opaque source bytes after converted tree payload"]
    for source in sorted(SOURCE.rglob("*.tree2")):
        relative = source.relative_to(SOURCE)
        original = source.read_bytes()
        tree = parse_tree2(source)
        converted = (CONVERTED / relative).read_bytes()
        suffix = original[tree.end_offset :]
        destination = OUTPUT / relative
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_bytes(converted + suffix)
        rows.append(f"{relative}: parsed={tree.end_offset} suffix={len(suffix)}")
    (WORK / "tree-suffix-variant-manifest.txt").write_text("\n".join(rows) + "\n")
    print("\n".join(rows))


if __name__ == "__main__":
    main()
