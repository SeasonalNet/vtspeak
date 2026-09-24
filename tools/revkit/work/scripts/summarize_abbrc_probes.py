#!/usr/bin/env python3
"""Validate Stage 10 C-key captures and summarize each returned token sequence."""

from __future__ import annotations

import json
import re
from pathlib import Path

STAGE = Path(__file__).resolve().parents[1] / "stage10"
MANIFEST = STAGE / "abbrc-probe-manifest.json"
OUTPUT = STAGE / "abbrc-probe-results.json"
ROW = re.compile(r"^ROW index=(\d+)\b.*surface=(.*)$")
SURFACE = re.compile(r'"([^"\\]*(?:\\.[^"\\]*)*)"')


def read_calls(path: Path) -> list[list[str]]:
    lines = path.read_text(errors="replace").splitlines()
    if "[Inferior 1 (Remote target) exited normally]" not in lines:
        raise ValueError(f"{path}: process did not exit normally")
    calls: list[list[str]] = []
    expected_count: int | None = None
    rows: list[str] = []
    index = 0
    while index < len(lines):
        returned = re.match(r"CONTEXT_RETURN count=(\d+)", lines[index])
        if returned is not None:
            if expected_count is not None:
                if len(rows) != expected_count:
                    raise ValueError(f"{path}: incomplete preceding context rows")
                calls.append(rows)
            expected_count = int(returned.group(1))
            rows = []
            index += 1
            continue
        row = ROW.match(lines[index])
        if row is not None:
            row_index = int(row.group(1))
            if row_index != len(rows):
                raise ValueError(f"{path}: expected row {len(rows)}, got {row_index}")
            surface = SURFACE.search(row.group(2))
            if surface is None and index + 1 < len(lines):
                surface = SURFACE.search(lines[index + 1])
                if surface is not None:
                    index += 1
            if surface is None:
                raise ValueError(f"{path}: row {row_index} lacks a surface")
            rows.append(surface.group(1))
        index += 1
    if expected_count is None or len(rows) != expected_count:
        raise ValueError(f"{path}: incomplete final context rows")
    calls.append(rows)
    return calls


def main() -> None:
    batches: dict[str, list[str]] = json.loads(MANIFEST.read_text(encoding="utf-8"))
    results: list[dict[str, object]] = []
    for name, keys in batches.items():
        calls = read_calls(STAGE / f"{name}.log")
        if len(calls) != len(keys):
            raise ValueError(f"{name}: {len(calls)} calls for {len(keys)} keys")
        results.extend(
            {"key": key, "surfaces": surfaces}
            for key, surfaces in zip(keys, calls, strict=True)
        )
    if len(results) != 330:
        raise ValueError(f"expected 330 C-key observations, got {len(results)}")
    OUTPUT.write_text(json.dumps(results, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    unchanged = sum(
        item["surfaces"] == ["We", "noted", item["key"], "today", "."]
        for item in results
    )
    changed = [item for item in results if item["surfaces"] != ["We", "noted", item["key"], "today", "."]]
    print(f"verified {len(results)} C-key context returns; {unchanged} preserve the template token sequence")
    print(f"{len(changed)} differ from the literal template token sequence; full rows: {OUTPUT}")
    for item in changed:
        print(f"{item['key']!r}: {' / '.join(item['surfaces'])}")


if __name__ == "__main__":
    main()
