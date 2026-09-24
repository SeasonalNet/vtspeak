#!/usr/bin/env python3
"""Verify the checked-in Stage 8 PCM and WAVE evidence without vendor inputs."""

from __future__ import annotations

import hashlib
import wave
from pathlib import Path


WORK = Path(__file__).resolve().parents[1]
STAGE8 = WORK / "stage8"


def wave_data(path: Path) -> tuple[int, bytes]:
    with wave.open(str(path), "rb") as wav:
        params = wav.getparams()
        if (params.nchannels, params.sampwidth, params.framerate, params.comptype) != (
            1,
            2,
            16000,
            "NONE",
        ):
            raise ValueError(f"unexpected PCM format in {path}: {params}")
        return params.nframes, wav.readframes(params.nframes)


def require_equal(label: str, actual: bytes, expected: bytes) -> None:
    if actual != expected:
        raise SystemExit(f"mismatch: {label}")
    print(f"matched {label} ({len(actual)} bytes)")


def main() -> None:
    cases = {
        "default": ("default.wav", ("default-block-000.pcm", "default-block-001.pcm")),
        "pitch120": ("pitch120.wav", ("pitch120-block-000.pcm", "pitch120-block-001.pcm")),
        "Hi": ("map-short.wav", ("map-block.pcm",)),
    }
    for name, (wav_name, block_names) in cases.items():
        frames, data = wave_data(STAGE8 / wav_name)
        blocks = b"".join((STAGE8 / block_name).read_bytes() for block_name in block_names)
        require_equal(f"{name} PCM blocks == {wav_name} data chunk", blocks, data)
        print(f"{wav_name}: {frames} frames")

    default_frames, default_data = wave_data(STAGE8 / "default.wav")
    reference_frames, reference_data = wave_data(WORK / "stage3" / "reference.wav")
    if default_frames != reference_frames:
        raise SystemExit("frame-count mismatch: default.wav vs Stage 3 reference")
    require_equal("default.wav == Stage 3 reference.wav", default_data, reference_data)
    print(f"default SHA-256: {hashlib.sha256((STAGE8 / 'default.wav').read_bytes()).hexdigest()}")

    pitch_frames, pitch_data = wave_data(STAGE8 / "pitch120.wav")
    pitch_reference_frames, pitch_reference_data = wave_data(
        WORK / "stage4" / "pitch120-reference.wav"
    )
    if pitch_frames != pitch_reference_frames:
        raise SystemExit("frame-count mismatch: pitch120.wav vs Stage 4 reference")
    require_equal("pitch120.wav == Stage 4 reference.wav", pitch_data, pitch_reference_data)
    print(f"pitch120 SHA-256: {hashlib.sha256((STAGE8 / 'pitch120.wav').read_bytes()).hexdigest()}")

    expected_frames = {
        "map.wav": 44717,
        "map-pitch120.wav": 47156,
        "map-speed120.wav": 38326,
        "map-volume120.wav": 44717,
    }
    for name, expected in expected_frames.items():
        frames, _ = wave_data(STAGE8 / name)
        if frames != expected:
            raise SystemExit(f"frame-count mismatch: {name}: got {frames}, expected {expected}")
        print(f"{name}: {frames} frames")


if __name__ == "__main__":
    main()
