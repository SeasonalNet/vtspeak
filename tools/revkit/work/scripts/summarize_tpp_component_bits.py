#!/usr/bin/env python3
"""Summarize component-bit values against standalone A0/A1 dictionary entries."""

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


def summarize(dictionary: Path, dll: Path) -> list[str]:
    rows = INSPECT_TPP.inspect(dictionary, dll, show_rows=True)["rows"]
    assert isinstance(rows, list)
    single_component = {
        row["key"]: row["payload"].split()[0]
        for row in rows
        if row["payload"].split()[0] in {"A0", "A1"}
    }
    occurrences: dict[tuple[str, str], list[str]] = collections.defaultdict(list)
    word_bits: dict[tuple[str, str], set[str]] = collections.defaultdict(set)

    for row in rows:
        primary = row["payload"].split()[0]
        if primary[0] not in "BCDE":
            continue
        words = row["key"].split("-")
        bits = primary[2:]
        if len(words) != len(bits):
            raise ValueError(f"component count mismatch: {row}")
        for word, bit in zip(words, bits):
            a_class = single_component.get(word, "NONE")
            occurrences[(a_class, bit)].append(row["key"])
            word_bits[(a_class, word)].add(bit)

    output = [
        "standalone_A_class\tcompound_bit\tcomponent_occurrences\t"
        "unique_components\tcomponents_seen_with_both_bits\tsample_compounds"
    ]
    for a_class in ("A0", "A1", "NONE"):
        for bit in ("0", "1"):
            compounds = occurrences.get((a_class, bit), [])
            words = {
                word
                for (candidate_class, word), values in word_bits.items()
                if candidate_class == a_class and bit in values
            }
            both = sum(
                values == {"0", "1"}
                for (candidate_class, _), values in word_bits.items()
                if candidate_class == a_class
            )
            output.append(
                f"{a_class}\t{bit}\t{len(compounds)}\t{len(words)}\t{both}\t"
                f"{','.join(compounds[:12])}"
            )
    return output


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--dictionary", type=Path, default=ROOT / "data-common/dict-eng/tppdict_eng"
    )
    parser.add_argument("--dll", type=Path, default=ROOT / "binary/vt_pau.dll")
    parser.add_argument(
        "--output",
        type=Path,
        default=ROOT / "tools/revkit/work/stage12/component-bit-corpus.tsv",
    )
    args = parser.parse_args()
    lines = summarize(args.dictionary, args.dll)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text("\n".join(lines) + "\n", encoding="ascii")
    print(f"wrote component-bit cross-reference to {args.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
