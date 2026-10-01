#!/usr/bin/env python3
"""Analyze user-dictionary effects with a naturally licensed James slot."""

from __future__ import annotations

import hashlib
from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
MODEL = re.compile(
    r"^JAMES_DICT model_load_ax=(-?\d+) sample_speaker=(-?\d+) "
    r"gate=(\d+) capacity=(-?\d+)$",
    re.MULTILINE,
)
SINGLE = re.compile(
    r"^JAMES_DICT case=(control|restored) synth=(-?\d+)$", re.MULTILINE
)
DICTIONARY = re.compile(
    r"^JAMES_DICT case=(p|a) load=(-?\d+) synth=(-?\d+) unload=(-?\d+)$",
    re.MULTILINE,
)
FILES = {
    "control": "james-userdict-control.wav",
    "p": "james-userdict-p.wav",
    "a": "james-userdict-a.wav",
    "restored": "james-userdict-restored.wav",
}


def read_wave(name: str) -> tuple[bytes, int, str]:
    path = ROOT / FILES[name]
    with wave.open(str(path), "rb") as source:
        props = (source.getnchannels(), source.getframerate(), source.getsampwidth())
        if props != (1, 16000, 2):
            raise ValueError(f"{name} has unexpected WAVE properties: {props}")
        frames = source.getnframes()
        pcm = source.readframes(frames)
    if len(pcm) != frames * 2:
        raise ValueError(f"{name} has unexpected frame or PCM byte count: {frames}")
    digest = hashlib.sha256(path.read_bytes()).hexdigest()
    return pcm, frames, digest


def main() -> None:
    log = (ROOT / "james-licensed-userdict-effect-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    model_rows = MODEL.findall(log)
    if model_rows != [("0", "-1", "1", "6")]:
        raise ValueError(f"unexpected licensed James state: {model_rows}")

    singles = {name: int(result) for name, result in SINGLE.findall(log)}
    if singles != {"control": 1, "restored": 1}:
        raise ValueError(f"unexpected control call results: {singles}")
    dictionaries = {
        name: tuple(map(int, (load, synth, unload)))
        for name, load, synth, unload in DICTIONARY.findall(log)
    }
    if dictionaries != {"p": (1, 1, 1), "a": (1, 1, 1)}:
        raise ValueError(f"unexpected dictionary call results: {dictionaries}")

    captures = {name: read_wave(name) for name in FILES}
    if captures["control"][0] != captures["restored"][0]:
        raise ValueError("unloading both dictionaries did not restore control PCM")
    for name in ("control", "p", "a", "restored"):
        print(
            f"{name} frames={captures[name][1]} sha256={captures[name][2]} "
            f"pcm_matches_control={captures[name][0] == captures['control'][0]}"
        )
    print("james_license_gate=1 dictionary_capacity=6 p_and_a_load_synth_unload=1")


if __name__ == "__main__":
    main()
