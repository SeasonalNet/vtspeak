#!/usr/bin/env python3
"""Write a full pefile structure dump for the local 2006 Kate MSI DLL."""

from __future__ import annotations

from pathlib import Path

import pefile


ROOT = Path("/work/corpus-parity/full-analysis-2006")
DLL = Path(
    "/work/corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll"
)


def main() -> None:
    pe = pefile.PE(str(DLL), fast_load=False)
    (ROOT / "pefile-dump.txt").write_text(pe.dump_info(), encoding="utf-8")


if __name__ == "__main__":
    main()
