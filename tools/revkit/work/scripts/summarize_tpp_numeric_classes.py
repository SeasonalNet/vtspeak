#!/usr/bin/env python3
"""Summarize corpus associations for primary TPP F/G numeric codes."""

from __future__ import annotations

import argparse
import collections
import importlib.util
import sys
from pathlib import Path

sys.dont_write_bytecode = True

ROOT = Path(__file__).resolve().parents[4]
INSPECTOR = ROOT / "tools/revkit/scripts/inspect_tpp_dictionary.py"
SPEC = importlib.util.spec_from_file_location("inspect_tpp_dictionary", INSPECTOR)
if SPEC is None or SPEC.loader is None:
    raise RuntimeError(f"cannot load {INSPECTOR}")
INSPECT_TPP = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(INSPECT_TPP)


def ending_bucket(key: str) -> str:
    normalized = key.replace("-", "")
    for ending in ("ING", "LY", "ER", "ED", "S"):
        if normalized.endswith(ending):
            return ending
    return "OTHER"


def summarize(dictionary: Path, dll: Path) -> list[str]:
    rows = INSPECT_TPP.inspect(dictionary, dll, show_rows=True)["rows"]
    groups: dict[tuple[str, int], list[str]] = collections.defaultdict(list)
    for row in rows:
        atoms = row["payload"].split()
        primary = atoms[0]
        if primary[0] in "FG":
            groups[(primary[0], int(primary[1:]))].append(row["key"])

    endings = ("ING", "LY", "ER", "ED", "S", "OTHER")
    output = [
        "family\tcode\tcount\t" + "\t".join(endings) + "\tsamples"
    ]
    for (family, code), keys in sorted(groups.items()):
        counts = collections.Counter(ending_bucket(key) for key in keys)
        samples = ",".join(keys[:8])
        output.append(
            f"{family}\t{code}\t{len(keys)}\t"
            + "\t".join(str(counts[ending]) for ending in endings)
            + f"\t{samples}"
        )
    return output


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--dictionary",
        type=Path,
        default=ROOT / "data-common/dict-eng/tppdict_eng",
    )
    parser.add_argument("--dll", type=Path, default=ROOT / "binary/vt_pau.dll")
    parser.add_argument(
        "--output",
        type=Path,
        default=ROOT / "tools/revkit/work/stage12/numeric-class-corpus.tsv",
    )
    args = parser.parse_args()
    lines = summarize(args.dictionary, args.dll)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text("\n".join(lines) + "\n", encoding="ascii")
    codes = sum(line.startswith(("F\t", "G\t")) for line in lines[1:])
    print(f"wrote {codes} F/G numeric-class rows to {args.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
