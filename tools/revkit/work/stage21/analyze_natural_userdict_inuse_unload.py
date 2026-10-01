#!/usr/bin/env python3
"""Validate the natural synthesis-context user-dictionary unload capture."""

from __future__ import annotations

import hashlib
from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
LOG = ROOT / "natural-userdict-inuse-unload-v3-api.log"
FILES = {
    "control": ROOT / "natural-userdict-inuse-control-v3.wav",
    "active": ROOT / "natural-userdict-inuse-active-v3.wav",
    "restored": ROOT / "natural-userdict-inuse-restored-v3.wav",
}
EXPECTED_HASHES = {
    "control": "613ff17ff3d7b8c9a411ef771ea1d02a051f38fcb918874e95dde30db39c0f36",
    "active": "491b29d0c3deb1d95770f8bdfe6fbcdfce5ace03f8c15dbfb9c57cf95b601cf1",
}


def require_once(pattern: str, log: str) -> tuple[str, ...]:
    rows = [match.groups() for match in re.finditer(pattern, log, re.MULTILINE)]
    if len(rows) != 1:
        raise ValueError(f"expected one log row for {pattern!r}, found {rows}")
    return rows[0]


def read_wave(path: Path) -> tuple[bytes, int, str]:
    with wave.open(str(path), "rb") as source:
        props = (source.getnchannels(), source.getframerate(), source.getsampwidth())
        if props != (1, 16000, 2):
            raise ValueError(f"{path.name} has unexpected WAVE properties: {props}")
        frames = source.getnframes()
        pcm = source.readframes(frames)
    if len(pcm) != frames * 2:
        raise ValueError(f"{path.name} frame and PCM lengths disagree")
    return pcm, frames, hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    log = LOG.read_text(encoding="utf-8", errors="replace")
    if require_once(
        r"^NATURAL_INUSE control_synth=(-?\d+)$", log
    ) != ("1",):
        raise ValueError("same-process no-dictionary synthesis failed")
    setup = require_once(
        r"^NATURAL_INUSE setup selector=(\d+) dict_load=(-?\d+) "
        r"gate_before=(\d+) gate_forced=(\d+) return=(0x[0-9a-f]+)$",
        log,
    )
    if setup[0:4] != ("4", "1", "0", "1"):
        raise ValueError(f"unexpected setup state: {setup}")
    context = require_once(
        r"^NATURAL_INUSE context=(0x[0-9a-f]+) "
        r"context_dictionary=(0x[0-9a-f]+) dictionary=(0x[0-9a-f]+)$",
        log,
    )
    if context[0] == "0x0" or context[1] == "0x0" or context[1] != context[2]:
        raise ValueError(f"context did not hold the loaded dictionary: {context}")
    expected_statuses = {
        r"^NATURAL_INUSE active_unload=(-?\d+)$": "-3",
        r"^NATURAL_INUSE synthesis_return=(-?\d+)$": "1",
        r"^NATURAL_INUSE idle_unload=(-?\d+)$": "1",
        r"^NATURAL_INUSE gate_restored=(\d+)$": "0",
        r"^NATURAL_INUSE restored_synth=(-?\d+)$": "1",
    }
    for pattern, expected in expected_statuses.items():
        if require_once(pattern, log) != (expected,):
            raise ValueError(f"unexpected status for {pattern!r}")

    captures = {name: read_wave(path) for name, path in FILES.items()}
    for name, expected in EXPECTED_HASHES.items():
        if captures[name][2] != expected:
            raise ValueError(f"{name} WAVE hash changed: {captures[name][2]}")
    if captures["control"][0] == captures["active"][0]:
        raise ValueError("active dictionary output unexpectedly matches control")
    if captures["control"][0] != captures["restored"][0]:
        raise ValueError("idle unload did not restore control PCM")
    for name, (_, frames, digest) in captures.items():
        print(f"{name} frames={frames} sha256={digest}")
    print(
        "context_dictionary_matches_loaded_dictionary=True "
        "active_unload=-3 synthesis=1 idle_unload=1 gate_restored=0 "
        "restored_pcm_matches_control=True"
    )


if __name__ == "__main__":
    main()
