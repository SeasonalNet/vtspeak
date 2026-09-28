#!/usr/bin/env python3
"""Create zero-byte files for isolated malformed-resource loader probes."""

from __future__ import annotations

import argparse
from pathlib import Path


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("output", type=Path)
    args = parser.parse_args()
    args.output.mkdir(parents=True, exist_ok=False)
    for name in ("dblist.idx", "merged-gen.dat"):
        (args.output / name).touch()


if __name__ == "__main__":
    main()
