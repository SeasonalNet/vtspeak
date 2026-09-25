#!/usr/bin/env python3
"""Validate and summarize Lead 4 original-engine context captures."""

from __future__ import annotations

import json
import re
from pathlib import Path

STAGE = Path(__file__).resolve().parents[1] / "stage14"
MANIFEST = STAGE / "lead4-probe-manifest.json"
OUTPUT = STAGE / "lead4-results.json"
ROW = re.compile(
    r"^ROW index=(\d+) flags=(0x[0-9a-f]+|0) source=(\d+) "
    r"class=(0x[0-9a-f]+|0) marker=(0x[0-9a-f]+|0) surface=(.*)$"
)
SURFACE = re.compile(r'"([^"\\]*(?:\\.[^"\\]*)*)"')
BYTE = re.compile(r"0x([0-9a-f]{2})")


def parse_capture(path: Path) -> list[dict[str, object]]:
    lines = path.read_text(errors="replace").splitlines()
    if "[Inferior 1 (Remote target) exited normally]" not in lines:
        raise ValueError(f"{path}: process did not exit normally")
    calls: list[dict[str, object]] = []
    index = 0
    while index < len(lines):
        returned = re.match(r"CONTEXT_RETURN count=(\d+)", lines[index])
        if returned is None:
            index += 1
            continue
        expected = int(returned.group(1))
        index += 1
        rows: list[dict[str, object]] = []
        while index < len(lines) and len(rows) < expected:
            match = ROW.match(lines[index])
            if match is None:
                index += 1
                continue
            row_index = int(match.group(1))
            if row_index != len(rows):
                raise ValueError(f"{path}: expected row {len(rows)}, found {row_index}")
            surface = SURFACE.search(match.group(6))
            if surface is None and index + 1 < len(lines):
                surface = SURFACE.search(lines[index + 1])
            if surface is None:
                raise ValueError(f"{path}: row {row_index} lacks a surface")
            code_line = next(
                (pos for pos in range(index + 1, len(lines)) if lines[pos].startswith(f"ROW_CODE index={row_index} bytes:")),
                None,
            )
            if code_line is None:
                raise ValueError(f"{path}: row {row_index} lacks phone/context bytes")
            code_bytes: list[int] = []
            for code_pos in range(code_line + 1, min(code_line + 4, len(lines))):
                byte_values = lines[code_pos].partition(":")[2]
                code_bytes.extend(int(value, 16) for value in BYTE.findall(byte_values))
            if len(code_bytes) != 24:
                raise ValueError(f"{path}: row {row_index} has {len(code_bytes)} code bytes")
            rows.append(
                {
                    "index": row_index,
                    "flags": match.group(2),
                    "source": int(match.group(3)),
                    "class": match.group(4),
                    "marker": match.group(5),
                    "surface": surface.group(1),
                    "phone_context_bytes": code_bytes,
                }
            )
            index = code_line + 4
        if len(rows) != expected:
            raise ValueError(f"{path}: expected {expected} returned rows, got {len(rows)}")
        calls.append({"count": expected, "rows": rows})
    if len(calls) != 1:
        raise ValueError(f"{path}: expected one context return, got {len(calls)}")
    return calls


def main() -> None:
    manifest: dict[str, dict[str, str]] = json.loads(MANIFEST.read_text(encoding="utf-8"))
    results: list[dict[str, object]] = []
    for fixture, metadata in manifest.items():
        capture = parse_capture(STAGE / f"{fixture}.log")[0]
        results.append({"fixture": fixture, **metadata, **capture})
    OUTPUT.write_text(json.dumps(results, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    print(f"verified {len(results)} independent context captures across {len({item['probe'] for item in results})} keys")
    for result in results:
        surfaces = " / ".join(row["surface"] for row in result["rows"])
        print(f"{result['probe']} [{result['context']}]: {surfaces}")
    print(f"full returned rows and phone/context bytes: {OUTPUT}")


if __name__ == "__main__":
    main()
