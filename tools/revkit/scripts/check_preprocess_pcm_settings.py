#!/usr/bin/env python3
"""Verify the controlled flag-6 scalar-option capture."""

from __future__ import annotations

import hashlib
import re
import struct
from pathlib import Path


ROOT = Path(__file__).resolve().parents[3]
WORK = ROOT / "tools" / "revkit" / "work" / "stage21"
CAPTURE = WORK / "preprocess-pcm-settings-api.log"
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
ROW = re.compile(
    r"^PREPROCESS_PCM_SETTINGS variant=(\d+) "
    r"pitch=(-?\d+) speed=(-?\d+) volume=(-?\d+) pause=(-?\d+) "
    r"raw=(0x[0-9a-f]+)$",
    re.MULTILINE,
)


def main() -> None:
    log = CAPTURE.read_text(encoding="utf-8", errors="replace")
    rows = ROW.findall(log)
    if len(rows) != len(OPTIONS):
        raise SystemExit(f"expected 11 call records, found {len(rows)}")
    if "[Inferior 1 (Remote target) exited normally]" not in log:
        raise SystemExit("inferior did not exit normally")

    for fields in rows:
        variant, pitch, speed, volume, pause, raw = fields
        index = int(variant)
        if index not in OPTIONS:
            raise SystemExit(f"unexpected variant {index}")
        if tuple(map(int, (pitch, speed, volume, pause))) != OPTIONS[index]:
            raise SystemExit(f"wrong option tuple for variant {index}: {fields}")
        if raw != "0x1":
            raise SystemExit(f"unexpected raw EAX for variant {index}: {raw}")

    outputs = {
        index: (WORK / f"preprocess-pcm-settings-v{index:02d}.pcm").read_bytes()
        for index in OPTIONS
    }
    if len(outputs[0]) != 23606 or outputs[0] != outputs[10]:
        raise SystemExit("default PCM size or repeatability check failed")
    for index in (1, 2, 3, 4, 5, 6):
        if outputs[index] == outputs[0]:
            raise SystemExit(f"expected setting {index} to change PCM bytes")
    for index in (7, 8, 9):
        if outputs[index] != outputs[0]:
            raise SystemExit(f"pause setting {index} changed this fixture")
    if any(outputs[5]):
        raise SystemExit("volume 0 output is not silent")

    expected_samples = {
        0: 11803,
        1: 12850,
        2: 11998,
        3: 25377,
        4: 2438,
        5: 11803,
        6: 11803,
        7: 11803,
        8: 11803,
        9: 11803,
        10: 11803,
    }
    for index, data in outputs.items():
        if len(data) % 2 or len(data) // 2 != expected_samples[index]:
            raise SystemExit(f"unexpected PCM16 frame count in variant {index}")
        samples = struct.unpack(f"<{len(data) // 2}h", data)
        print(
            f"PREPROCESS_PCM_CHECK variant={index:02d} samples={len(samples)} "
            f"sha256={hashlib.sha256(data).hexdigest()}"
        )


if __name__ == "__main__":
    main()
