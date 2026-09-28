#!/usr/bin/env python3
"""Create symlink-only resource overlays that omit one VoiceText resource."""

from __future__ import annotations

import argparse
from pathlib import Path


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("output", type=Path)
    args = parser.parse_args()

    repository = Path(__file__).resolve().parents[4]
    cases = {
        "omit-gen-dat": ("data-paul/M16/dat", "merged-gen.dat", "/vendor-paul/M16/dat"),
        "omit-gen-upm": ("data-paul/M16/dat", "merged-gen.upm", "/vendor-paul/M16/dat"),
        "omit-unit-gen-idx": (
            "data-paul/M16/mc_idx_tbl",
            "unit-gen.idx",
            "/vendor-paul/M16/mc_idx_tbl",
        ),
        "omit-dblist-idx-preserve-tree": (
            "data-paul/M16",
            "dblist.idx",
            "/vendor-paul/M16",
        ),
        "omit-cepdist-tbl": (
            "data-paul/M16/ttsdata/dist_tbl",
            "cepdist.tbl",
            "/vendor-paul/M16/ttsdata/dist_tbl",
        ),
        "omit-pitch-nbt": (
            "data-paul/M16/ttsdata/tree3/pitch",
            "nbt.tree3",
            "/vendor-paul/M16/ttsdata/tree3/pitch",
        ),
        "omit-atmt-tree": (
            "data-common/dict-eng",
            "atmt.tree3",
            "/vendor-common/dict-eng",
        ),
        "omit-engbi-tree": (
            "data-common/dict-eng",
            "engbi.tree3",
            "/vendor-common/dict-eng",
        ),
        "omit-sbd-tree": (
            "data-common/dict-eng",
            "sbd.tree3",
            "/vendor-common/dict-eng",
        ),
        "omit-tppdict": (
            "data-common/dict-eng",
            "tppdict_eng",
            "/vendor-common/dict-eng",
        ),
        "omit-hashidx-tpp": (
            "data-common/dict-eng",
            "hashidx_eng_tpp",
            "/vendor-common/dict-eng",
        ),
    }
    for case, (source_relative, omitted_name, vendor_path) in cases.items():
        target = args.output / case
        target.mkdir(parents=True, exist_ok=False)
        source = repository / source_relative
        for entry in sorted(source.iterdir()):
            if entry.name != omitted_name and (entry.is_file() or entry.is_dir()):
                (target / entry.name).symlink_to(f"{vendor_path}/{entry.name}")


if __name__ == "__main__":
    main()
