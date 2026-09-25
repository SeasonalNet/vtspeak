#!/usr/bin/env python3
"""Validate and map the full H/T/C Lead 4 context matrix."""

from __future__ import annotations

import json
import re
from pathlib import Path

STAGE = Path(__file__).resolve().parents[1] / "stage15"
MANIFEST = STAGE / "lead4-matrix-manifest.json"
OUTPUT = STAGE / "lead4-key-matrix.json"
SUMMARY = STAGE / "lead4-key-summary.tsv"
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
                (
                    pos
                    for pos in range(index + 1, len(lines))
                    if lines[pos].startswith(f"ROW_CODE index={row_index} bytes:")
                ),
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
    return calls


def target_rows(family: str, context: str, rows: list[dict[str, object]]) -> list[dict[str, object]]:
    output = [str(row["surface"]) for row in rows]

    def locate(token: str, start: int = 0) -> int:
        needle = token.casefold().rstrip(".,!?;:")
        for position in range(start, len(output)):
            if output[position].casefold().rstrip(".,!?;:") == needle:
                return position
        raise ValueError(f"{family} {context}: output anchor {token!r} was not found in {output!r}")

    def between(left: str, right: str, *, skip_comma: bool = False) -> list[dict[str, object]]:
        start = locate(left) + 1
        if skip_comma and start < len(rows) and output[start] == ",":
            start += 1
        end = locate(right, start)
        selected = rows[start:end]
        while selected and selected[0]["surface"] == ",":
            selected = selected[1:]
        while selected and selected[-1]["surface"] == ",":
            selected = selected[:-1]
        while selected and selected[-1]["surface"] == ".":
            selected = selected[:-1]
        if not selected:
            raise ValueError(f"{family} {context}: target span is empty in {output!r}")
        return selected

    if family == "H":
        if context == "ordinary-sentence":
            return between("noted", "today")
        selected = rows[: locate("Smith")]
        while selected and selected[-1]["surface"] in {".", ",", "?", "!"}:
            selected = selected[:-1]
        if not selected:
            raise ValueError(f"H {context}: target span is empty in {output!r}")
        return selected
    if family == "T":
        if context == "ordinary-sentence":
            return between("noted", "today")
        return between("Acme", "hired", skip_comma=context == "comma-company")
    if context in {"baseline-neutral", "case-variant-neutral"}:
        return between("noted", "today")
    if context == "comma-neighbor":
        return between("noted", "today", skip_comma=True)
    if context == "calendar-date":
        return between("on", "fifth")
    if context == "calendar-weekday":
        return between("on", "Next")
    if context == "clock-or-zone":
        return between("five", "Next")
    if context == "state-code":
        return between("in", "Next")
    if context == "address-or-building":
        return between("Main", "Next")
    if context == "numeric-or-ordinary":
        return between("is", "long")
    raise ValueError(f"unknown context type {family} {context}")


def main() -> None:
    manifest: dict[str, dict[str, object]] = json.loads(MANIFEST.read_text(encoding="utf-8"))
    mapped: dict[tuple[str, str], dict[str, object]] = {}
    total_calls = 0
    for fixture, metadata in manifest.items():
        entries = metadata["entries"]
        calls = parse_capture(STAGE / f"{fixture}.log")
        total_calls += len(calls)
        flat_rows = [row for call in calls for row in call["rows"]]
        segments: list[list[dict[str, object]]] = []
        current: list[dict[str, object]] = []
        index = 0
        while index < len(flat_rows):
            row = flat_rows[index]
            current.append(row)
            if str(row["surface"]).casefold().rstrip(".") == "next":
                if index + 1 < len(flat_rows) and flat_rows[index + 1]["surface"] in {".", "?", "!"}:
                    current.append(flat_rows[index + 1])
                    index += 1
                segments.append(current)
                current = []
            index += 1
        if current:
            raise ValueError(f"{fixture}: rows remain after the last Next delimiter")
        if len(segments) != len(entries):
            raise ValueError(f"{fixture}: expected {len(entries)} delimited probes, got {len(segments)}")
        for entry, rows in zip(entries, segments, strict=True):
            identity = (str(entry["family"]), str(entry["key"]))
            if identity not in mapped:
                mapped[identity] = {
                    "family": identity[0],
                    "key": identity[1],
                    "match_tag": entry["match_tag"],
                    "case_variant": entry["case_variant"],
                    "case_variant_is_other_inventory_key": entry["case_variant_is_other_inventory_key"],
                    "contexts": [],
                }
            elif mapped[identity]["match_tag"] != entry["match_tag"]:
                raise ValueError(f"{identity}: conflicting match tags")
            mapped[identity]["contexts"].append(
                {
                    "name": entry["context"],
                    "input": entry["input"],
                    "delimiter": entry["delimiter"],
                    "rows": rows,
                    "target_rows": target_rows(str(entry["family"]), str(entry["context"]), rows),
                }
            )

    expected_counts = {"H": 55, "T": 42, "C": 330}
    family_counts = {family: 0 for family in expected_counts}
    summary_lines = [
        "family\tkey\tmatch_tag\tcase_variant\tcase_variant_key_collision\tcontext\tinput\ttarget_surfaces\tordered_returned_surfaces"
    ]
    results: list[dict[str, object]] = []
    for identity, item in sorted(mapped.items()):
        family, key = identity
        contexts = item["contexts"]
        names = [str(context["name"]) for context in contexts]
        baseline_name = {"H": "baseline-title", "T": "baseline-company", "C": "baseline-neutral"}[family]
        if names.count(baseline_name) != 1:
            raise ValueError(f"{family} {key}: baseline capture is missing or duplicated")
        family_counts[family] += 1
        for context in contexts:
            summary_lines.append(
                "\t".join(
                    (
                        family,
                        key,
                        str(item["match_tag"]),
                        str(item["case_variant"]),
                        str(item["case_variant_is_other_inventory_key"]),
                        str(context["name"]),
                        str(context["input"]),
                        " / ".join(str(row["surface"]) for row in context["target_rows"]),
                        " / ".join(str(row["surface"]) for row in context["rows"]),
                    )
                )
            )
        results.append(item)

    if family_counts != expected_counts:
        raise ValueError(f"inventory coverage mismatch: {family_counts!r} != {expected_counts!r}")
    expected_contexts = {"H": 4, "T": 5, "C": 4}
    for item in results:
        if len(item["contexts"]) != expected_contexts[item["family"]]:
            raise ValueError(f"{item['family']} {item['key']}: incomplete context matrix")

    OUTPUT.write_text(json.dumps(results, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    SUMMARY.write_text("\n".join(summary_lines) + "\n", encoding="utf-8")
    print(f"mapped {len(results)} keys in {sum(len(item['contexts']) for item in results)} delimited observations and {total_calls} original-engine context returns")
    print("per-family keys: " + ", ".join(f"{family}={count}" for family, count in family_counts.items()))
    print(f"full rows: {OUTPUT}")
    print(f"per-key summary: {SUMMARY}")


if __name__ == "__main__":
    main()
