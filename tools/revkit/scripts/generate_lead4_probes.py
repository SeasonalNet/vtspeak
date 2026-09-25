#!/usr/bin/env python3
"""Generate controlled Lead 4 abbreviation-context fixtures."""

from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
STAGE = ROOT / "tools/revkit/work/stage14"
INPUTS = STAGE / "inputs"

# Cases are deliberately narrow: each compares one inventory key across
# contexts suggested by its spelling or the existing Stage 10 evidence.
PROBES: dict[str, list[tuple[str, str]]] = {
    "H-jr": [
        ("name-suffix", "John Smith Jr. spoke."),
        ("title-before-name", "Jr. Smith spoke."),
        ("unpunctuated-suffix", "John Smith Jr spoke."),
    ],
    "H-sr": [
        ("name-suffix", "John Smith Sr. spoke."),
        ("title-before-name", "Sr. Smith spoke."),
        ("unpunctuated-suffix", "John Smith Sr spoke."),
    ],
    "H-gen": [
        ("rank-title", "Gen. Lee spoke."),
        ("lowercase-title", "gen. Lee spoke."),
        ("ordinary-word", "The general spoke."),
    ],
    "H-rep": [
        ("title-before-name", "Rep. Jones spoke."),
        ("lowercase-title", "rep. Jones spoke."),
        ("ordinary-noun", "The representative spoke."),
    ],
    "T-co": [
        ("company-suffix", "Acme Co. hired Lee."),
        ("lowercase-suffix", "Acme co hired Lee."),
        ("standalone", "The co. was listed."),
    ],
    "T-inc": [
        ("company-suffix", "Acme Inc. hired Lee."),
        ("comma-company-suffix", "Acme, Inc. hired Lee."),
        ("standalone", "The inc. was listed."),
    ],
    "T-LLC": [
        ("inventory-case", "Acme LLC hired Lee."),
        ("lowercase", "Acme llc hired Lee."),
        ("title-case", "Acme Llc hired Lee."),
    ],
    "T-AG": [
        ("inventory-case", "Miller AG hired Lee."),
        ("punctuated", "Miller A.G. hired Lee."),
        ("standalone", "AG spoke today."),
    ],
    "C-am": [
        ("clock-lowercase", "It is 5 am."),
        ("ordinary-verb", "I am ready."),
        ("clock-uppercase", "It is 5 AM."),
    ],
    "C-in": [
        ("measurement", "It is 5 in. long."),
        ("preposition", "We are in town."),
        ("sentence-initial", "In May, we left."),
    ],
    "C-mi": [
        ("measurement", "It is 5 mi. away."),
        ("state-code", "We crossed MI today."),
        ("lowercase-state-spelling", "We crossed mi today."),
    ],
    "C-ct": [
        ("address-suffix", "42 Main Ct."),
        ("state-code", "We live in CT."),
        ("standalone", "The ct. was noted."),
    ],
    "C-CA": [
        ("state-code", "We live in CA."),
        ("mixed-case", "We live in Ca."),
        ("lowercase", "We live in ca."),
    ],
    "C-IN": [
        ("state-code", "We drove through IN."),
        ("preposition", "We are in town."),
        ("sentence-initial", "In May, we left."),
    ],
    "C-CO": [
        ("state-code", "We live in CO."),
        ("company-suffix", "Acme Co. hired Lee."),
        ("sentence-initial", "Co. was written on the page."),
    ],
    "C-Jan": [
        ("date", "We met on Jan. 5."),
        ("month-name", "It happened in Jan."),
        ("person-name", "Jan spoke today."),
    ],
}


def main() -> None:
    INPUTS.mkdir(parents=True, exist_ok=True)
    manifest: dict[str, dict[str, str]] = {}
    fixtures: list[str] = []
    for probe_index, (probe, contexts) in enumerate(PROBES.items(), start=1):
        for index, (context, text) in enumerate(contexts, start=1):
            fixture = f"lead4-{probe_index:02d}-{probe.lower()}-{index:02d}"
            (INPUTS / f"{fixture}.txt").write_text(text + "\n", encoding="utf-8")
            manifest[fixture] = {"probe": probe, "context": context, "input": text}
            fixtures.append(fixture)
    path = STAGE / "lead4-probe-manifest.json"
    path.write_text(json.dumps(manifest, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")
    (STAGE / "lead4-fixtures.txt").write_text("\n".join(fixtures) + "\n", encoding="utf-8")
    print(f"generated {len(manifest)} independent captures across {len(PROBES)} keys: {path}")


if __name__ == "__main__":
    main()
