#!/usr/bin/env python3
"""Generate context, case, and punctuation probes for every H/T/C key."""

from __future__ import annotations

import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
INVENTORY = ROOT / "docs/reverse-engineering/abbreviation-table-inventory.md"
STAGE = ROOT / "tools/revkit/work/stage15"
INPUTS = STAGE / "inputs"
ROWS = re.compile(r"^\| `([^`]+)` \| ([12]) \|$", re.MULTILINE)
BATCH_SIZE = {"H": 8, "T": 7, "C": 8}

MONTHS = {
    "jan", "feb", "mar", "apr", "aug", "sep", "sept", "oct", "nov", "dec", "jun", "jul",
}
WEEKDAYS = {
    "mon", "tue", "tues", "wed", "thu", "thur", "thurs", "fri", "sat", "sun",
}
TIME_KEYS = {
    "am", "pm", "a.m", "p.m", "gmt", "g.m.t", "est", "e.s.t", "cst", "c.s.t",
    "mst", "pst", "p.s.t", "kst", "k.s.t",
}
STATE_KEYS = {
    "AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "ID",
    "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", "MA", "MI", "MN", "MS",
    "MO", "MT", "NE", "NV", "NH", "NJ", "NM", "NY", "NC", "ND", "OH", "OK",
    "OR", "PA", "RI", "SC", "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV",
    "WI", "WY", "DC",
}
ADDRESS_KEYS = {
    "apt", "bdrm", "bldg", "bsmt", "ct", "dr", "flr", "frnt", "hngr", "lowr",
    "rm", "st", "ste", "trlr", "uppr",
}


def read_set(markdown: str, name: str) -> list[tuple[str, str]]:
    heading = f"## `abbr{name.lower()}_sort.txt2`"
    start = markdown.index(heading) + len(heading)
    end = markdown.find("\n## `", start)
    section = markdown[start:] if end == -1 else markdown[start:end]
    return ROWS.findall(section)


def alternate_case(key: str) -> str:
    if key.islower():
        return key.upper()
    if key.isupper():
        return key.lower()
    if key[0].isupper() and key[1:].islower():
        return key.lower()
    return key.swapcase()


def c_context(key: str) -> tuple[str, str]:
    folded = key.replace(".", "").lower()
    if folded in MONTHS:
        return "calendar-date", f"We met on {key}. 5."
    if folded in WEEKDAYS:
        return "calendar-weekday", f"We met on {key}."
    if key in TIME_KEYS or folded in {"est", "cst", "mst", "pst", "kst", "gmt"}:
        return "clock-or-zone", f"The meeting is at 5 {key}."
    if key in STATE_KEYS:
        return "state-code", f"We live in {key}."
    if folded in ADDRESS_KEYS:
        return "address-or-building", f"42 Main {key}."
    return "numeric-or-ordinary", f"It is 5 {key}. long."


def contexts(family: str, key: str) -> list[tuple[str, str]]:
    variant = alternate_case(key)
    if family == "H":
        return [
            ("baseline-title", f"{key}. Smith spoke."),
            ("title-without-period", f"{key} Smith spoke."),
            ("ordinary-sentence", f"We noted {key} today."),
            ("case-variant-title", f"{variant}. Smith spoke."),
        ]
    if family == "T":
        return [
            ("baseline-company", f"Acme {key}. hired Lee."),
            ("comma-company", f"Acme, {key} hired Lee."),
            ("company-without-period", f"Acme {key} hired Lee."),
            ("ordinary-sentence", f"We noted {key} today."),
            ("case-variant-company", f"Acme {variant}. hired Lee."),
        ]
    category, likely_context = c_context(key)
    return [
        ("baseline-neutral", f"We noted {key} today."),
        (category, likely_context),
        ("comma-neighbor", f"We noted, {key}, today."),
        ("case-variant-neutral", f"We noted {variant} today."),
    ]


def main() -> None:
    markdown = INVENTORY.read_text(encoding="utf-8")
    inventories = {family: read_set(markdown, family) for family in "HTC"}
    all_key_sets = {family: {key for key, _tag in rows} for family, rows in inventories.items()}
    STAGE.mkdir(parents=True, exist_ok=True)
    INPUTS.mkdir(parents=True, exist_ok=True)

    manifest: dict[str, dict[str, object]] = {}
    fixtures: list[str] = []
    family_fixtures: dict[str, list[str]] = {family: [] for family in "HTC"}
    summary = {family: len(rows) for family, rows in inventories.items()}
    sequence = 0
    for family, rows in inventories.items():
        size = BATCH_SIZE[family]
        for batch_index, offset in enumerate(range(0, len(rows), size), start=1):
            batch = rows[offset : offset + size]
            sequence += 1
            fixture = f"lead4-full-{family.lower()}-{batch_index:02d}"
            entries: list[dict[str, str | bool]] = []
            sentences: list[str] = []
            for key, tag in batch:
                variant = alternate_case(key)
                variant_is_other_inventory_key = variant in all_key_sets[family] and variant != key
                for context, text in contexts(family, key):
                    sentences.append(f"{text} Next.")
                    entries.append(
                        {
                            "family": family,
                            "key": key,
                            "match_tag": tag,
                            "context": context,
                            "input": text,
                            "delimiter": "Next.",
                            "case_variant": variant,
                            "case_variant_is_other_inventory_key": variant_is_other_inventory_key,
                        }
                    )
            (INPUTS / f"{fixture}.txt").write_text(" ".join(sentences) + "\n", encoding="utf-8")
            manifest[fixture] = {"family": family, "entries": entries}
            fixtures.append(fixture)
            family_fixtures[family].append(fixture)

    (STAGE / "lead4-matrix-manifest.json").write_text(
        json.dumps(manifest, indent=2, ensure_ascii=False) + "\n", encoding="utf-8"
    )
    (STAGE / "lead4-matrix-fixtures.txt").write_text("\n".join(fixtures) + "\n", encoding="utf-8")
    for family, names in family_fixtures.items():
        (STAGE / f"lead4-matrix-fixtures-{family.lower()}.txt").write_text(
            "\n".join(names) + "\n", encoding="utf-8"
        )
    print(
        f"generated {sum(summary.values())} inventory keys across {len(fixtures)} batches "
        f"({summary['H']} H, {summary['T']} T, {summary['C']} C)"
    )
    print(f"each key has baseline, contextual, punctuation, and case coverage: {STAGE}")


if __name__ == "__main__":
    main()
