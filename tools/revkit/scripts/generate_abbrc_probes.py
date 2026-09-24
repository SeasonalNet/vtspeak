#!/usr/bin/env python3
"""Generate neutral one-context batches for the C abbreviation inventory."""

from __future__ import annotations

import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
INVENTORY = ROOT / "docs/reverse-engineering/abbreviation-table-inventory.md"
OUTPUT = ROOT / "tools/revkit/work/stage10/inputs"
ROWS = re.compile(r"^\| `([^`]+)` \| ([12]) \|$", re.MULTILINE)


def read_keys(markdown: str) -> list[tuple[str, str]]:
    heading = "## `abbrc_sort.txt2`"
    start = markdown.index(heading) + len(heading)
    end = markdown.find("\n## `", start)
    section = markdown[start:] if end == -1 else markdown[start:end]
    return ROWS.findall(section)


def main() -> None:
    rows = read_keys(INVENTORY.read_text(encoding="utf-8"))
    OUTPUT.mkdir(parents=True, exist_ok=True)
    batches: dict[str, list[str]] = {}
    for batch_index, offset in enumerate(range(0, len(rows), 10), start=1):
        batch = rows[offset : offset + 10]
        name = f"abbr-c-{batch_index:02d}"
        sentences = [f"We noted {key} today." for key, _tag in batch]
        (OUTPUT / f"{name}.txt").write_text(" ".join(sentences) + "\n", encoding="utf-8")
        batches[name] = [key for key, _tag in batch]
    manifest = OUTPUT.parent / "abbrc-probe-manifest.json"
    manifest.write_text(json.dumps(batches, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    print(f"generated {len(rows)} keys in {len(batches)} neutral-context batches: {manifest}")


if __name__ == "__main__":
    main()
