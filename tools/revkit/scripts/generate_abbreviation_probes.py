#!/usr/bin/env python3
"""Generate bounded H/T abbreviation context batches from the checked-in inventory."""

from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
INVENTORY = ROOT / "docs/reverse-engineering/abbreviation-table-inventory.md"
OUTPUT = ROOT / "tools/revkit/work/stage10/inputs"
SETS = {"abbrh": "H", "abbrt": "T"}
ROWS = re.compile(r"^\| `([^`]+)` \| ([12]) \|$", re.MULTILINE)


def read_keys(markdown: str, filename: str) -> list[tuple[str, str]]:
    heading = f"## `{filename}_sort.txt2`"
    start = markdown.index(heading) + len(heading)
    end = markdown.find("\n## `", start)
    section = markdown[start:] if end == -1 else markdown[start:end]
    return ROWS.findall(section)


def main() -> None:
    markdown = INVENTORY.read_text(encoding="utf-8")
    OUTPUT.mkdir(parents=True, exist_ok=True)
    generated = 0
    for filename, label in SETS.items():
        rows = read_keys(markdown, filename)
        for batch_index, offset in enumerate(range(0, len(rows), 10), start=1):
            batch = rows[offset : offset + 10]
            if label == "H":
                sentences = [f"{key}. Smith spoke." for key, _tag in batch]
            else:
                sentences = [f"Acme {key}. hired Lee." for key, _tag in batch]
            path = OUTPUT / f"abbr-{label.lower()}-{batch_index:02d}.txt"
            path.write_text(" ".join(sentences) + "\n", encoding="utf-8")
            generated += 1
        print(f"{label}: {len(rows)} keys in {(len(rows) + 9) // 10} batches")
    print(f"generated {generated} fixtures in {OUTPUT}")


if __name__ == "__main__":
    main()
