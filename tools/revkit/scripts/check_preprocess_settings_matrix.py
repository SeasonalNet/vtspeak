#!/usr/bin/env python3
"""Verify the controlled VT_TextToPreprocessInfo settings capture."""

from __future__ import annotations

import hashlib
import re
from pathlib import Path


ROOT = Path(__file__).resolve().parents[3]
WORK = ROOT / "tools" / "revkit" / "work" / "stage21"
CAPTURE = WORK / "preprocess-settings-matrix-api.log"
OPTIONS = {
    0: (-1, -1, -1, -1),
    1: (50, -1, -1, -1),
    2: (200, -1, -1, -1),
    3: (-1, 50, -1, -1),
    4: (-1, 400, -1, -1),
    5: (-1, -1, 0, -1),
    6: (-1, -1, 500, -1),
    7: (-1, -1, -1, 0),
    8: (-1, -1, -1, 250),
    9: (-1, -1, -1, 65535),
    10: (-1, -1, -1, -1),
}
PATH_CODE = {3: "3", 5: "5", 7: "7", 10: "a"}
ROW = re.compile(
    r"^PREPROCESS_SETTINGS flag=(\d+) variant=(\d+) "
    r"pitch=(-?\d+) speed=(-?\d+) volume=(-?\d+) pause=(-?\d+) "
    r"(?:state=(\d+) )?raw=(0x[0-9a-f]+) path=(\S+)$",
    re.MULTILINE,
)


def main() -> None:
    log = CAPTURE.read_text(encoding="utf-8", errors="replace")
    rows = ROW.findall(log)
    if len(rows) != 44:
        raise SystemExit(f"expected 44 call records, found {len(rows)}")
    if "[Inferior 1 (Remote target) exited normally]" not in log:
        raise SystemExit("inferior did not exit normally")

    observed: set[tuple[int, int]] = set()
    for fields in rows:
        flag, variant, pitch, speed, volume, pause, state, raw, path = fields
        key = (int(flag), int(variant))
        if key in observed:
            raise SystemExit(f"duplicate call record: {key}")
        observed.add(key)
        if key[0] not in PATH_CODE or key[1] not in OPTIONS:
            raise SystemExit(f"unexpected call record: {key}")
        got_options = tuple(map(int, (pitch, speed, volume, pause)))
        if got_options != OPTIONS[key[1]] or raw != "0x1":
            raise SystemExit(f"unexpected call result/options: {key} {fields}")
        expected_path = f"Z:/work/stage21/ptf{PATH_CODE[key[0]]}v{key[1]:02d}"
        if path != expected_path:
            raise SystemExit(f"unexpected output path for {key}: {path}")
        if key[0] == 7 and state != "134456":
            raise SystemExit(f"flag 7 state was not held constant for {key}: {state}")

    expected = {
        (flag, variant)
        for flag in PATH_CODE
        for variant in OPTIONS
    }
    if observed != expected:
        raise SystemExit(f"missing calls: {sorted(expected - observed)}")

    for flag, code in PATH_CODE.items():
        outputs = [WORK / f"ptf{code}v{variant:02d}" for variant in OPTIONS]
        missing = [path.name for path in outputs if not path.is_file()]
        if missing:
            raise SystemExit(f"missing outputs for flag {flag}: {missing}")
        data = [path.read_bytes() for path in outputs]
        if any(item != data[0] for item in data[1:]):
            changed = [path.name for path, item in zip(outputs, data) if item != data[0]]
            raise SystemExit(f"settings changed flag {flag} output: {changed}")
        digest = hashlib.sha256(data[0]).hexdigest()
        print(
            f"PREPROCESS_SETTINGS_CHECK flag={flag} variants={len(outputs)} "
            f"bytes={len(data[0])} sha256={digest} identical=yes"
        )


if __name__ == "__main__":
    main()
