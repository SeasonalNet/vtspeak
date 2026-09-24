#!/usr/bin/env python3
"""Check Stage 9 sample-buffer captures against WAVE and Stage 6 outputs."""

from __future__ import annotations

import hashlib
import wave
from pathlib import Path

WORK = Path(__file__).resolve().parents[1]


def main() -> None:
    for name, expected_frames in (("numbers", 108_601), ("abbreviations", 116_689)):
        wav_path = WORK / "stage9" / f"{name}.wav"
        with wave.open(str(wav_path), "rb") as wav:
            params = (wav.getnchannels(), wav.getsampwidth(), wav.getframerate(), wav.getcomptype())
            if params != (1, 2, 16_000, "NONE") or wav.getnframes() != expected_frames:
                raise ValueError(f"unexpected WAVE format or frame count: {wav_path}")
            pcm = wav.readframes(expected_frames)
        blocks = (WORK / "stage9" / f"{name}-blocks.pcm").read_bytes()
        if blocks != pcm:
            raise ValueError(f"captured PCM blocks differ from WAVE data: {name}")
        prior = WORK / "stage5" / "probes" / "stage6" / name / "output.wav"
        if wav_path.read_bytes() != prior.read_bytes():
            raise ValueError(f"Stage 9 WAVE differs from Stage 6 reference: {name}")
        print(
            f"{name}: {expected_frames} frames, {len(pcm)} matching PCM bytes, "
            f"SHA-256={hashlib.sha256(pcm).hexdigest()}"
        )


if __name__ == "__main__":
    main()
