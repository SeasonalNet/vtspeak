#!/usr/bin/env python3
"""Validate non-null TextToFile path outcomes from Stage 21."""

from __future__ import annotations

from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
ROWS = re.compile(
    r"^TEXT_FILE_PATH name=(\w+) result=(-?\d+)$", re.MULTILINE
)
EXPECTED = {"valid": 1, "directory": -6, "missing_parent": -6}


def main() -> None:
    log = (ROOT / "text-file-path-error-grid-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    rows = ROWS.findall(log)
    observed = {name: int(result) for name, result in rows}
    if len(rows) != len(EXPECTED) or observed != EXPECTED:
        raise ValueError(f"unexpected path results: {observed}")

    valid = ROOT / "text-file-path-probe-valid.wav"
    control = ROOT / "sandbox/stage5/output.wav"
    with wave.open(str(valid), "rb") as source:
        if (source.getnchannels(), source.getframerate(), source.getsampwidth()) != (
            1,
            16000,
            2,
        ):
            raise ValueError("valid-path output has unexpected WAVE properties")
        if source.getnframes() != 11803:
            raise ValueError("valid-path output has unexpected frame count")
    if valid.read_bytes() != control.read_bytes():
        raise ValueError("valid-path output differs from the sandbox control")
    if (ROOT / "text-file-path-missing-parent-v1").exists():
        raise ValueError("missing-parent path unexpectedly appeared")
    if not (ROOT / "sandbox/stage5").is_dir():
        raise ValueError("directory control is not present")

    print("paths=3 valid=1 directory=-6 missing_parent=-6")
    print("valid_output=11803_frames_pcm16_mono_16000Hz exact_control_match=true")
    print("failed_paths_created_no_missing_parent=true")


if __name__ == "__main__":
    main()
