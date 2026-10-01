#!/usr/bin/env python3
"""Validate sentence-pause setter effects on Stage 21 WAVE captures."""

from __future__ import annotations

import hashlib
from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
VALUES = (0, 1, 199, 200, 201, 250, 500, 924, 925, 926, 65534, 65535)
PERIOD_CONTEXTS = ("upper", "upper_lower", "lower")
CONTEXTS = (*PERIOD_CONTEXTS, "no_period")
GETTER = re.compile(
    r"^SENT_PAUSE value=(\d+) getter_ret=(\d+) pitch=(\d+) speed=(\d+) "
    r"volume=(\d+) getter_value=(\d+)$",
    re.MULTILINE,
)
SYNTH = re.compile(
    r"^SENT_PAUSE value=(\d+) context=(\S+) synth_ret=(\d+)$",
    re.MULTILINE,
)


def read_wave(value: int, context: str) -> tuple[bytes, int, str]:
    path = ROOT / f"sentence-pause-{value}-{context}.wav"
    with wave.open(str(path), "rb") as source:
        if (source.getnchannels(), source.getframerate(), source.getsampwidth()) != (1, 16000, 2):
            raise ValueError(f"unexpected WAVE format: {path}")
        frames = source.getnframes()
        pcm = source.readframes(frames)
    if len(pcm) != frames * 2:
        raise ValueError(f"PCM byte count differs from frame count: {path}")
    return pcm, frames, hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    log = (ROOT / "sentence-pause-context-grid-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    getter_rows = GETTER.findall(log)
    synth_rows = SYNTH.findall(log)
    if len(getter_rows) != len(VALUES):
        raise ValueError(f"expected {len(VALUES)} getter rows, got {len(getter_rows)}")
    if len(synth_rows) != len(VALUES) * len(CONTEXTS):
        raise ValueError(f"expected {len(VALUES) * len(CONTEXTS)} synthesis rows, got {len(synth_rows)}")
    for value_text, result_text, pitch_text, speed_text, volume_text, observed_text in getter_rows:
        value = int(value_text)
        if (
            value not in VALUES
            or int(result_text) != 1
            or (int(pitch_text), int(speed_text), int(volume_text)) != (100, 100, 200)
            or int(observed_text) != value
        ):
            raise ValueError(f"getter mismatch for sentence pause {value_text}")
    observed_synth = {
        (int(value), context): int(result)
        for value, context, result in synth_rows
    }
    expected_synth = {(value, context) for value in VALUES for context in CONTEXTS}
    if set(observed_synth) != expected_synth or any(result != 1 for result in observed_synth.values()):
        raise ValueError("synthesis coordinates or return values differ from the complete grid")

    captures = {
        (value, context): read_wave(value, context)
        for value in VALUES
        for context in CONTEXTS
    }
    baseline_period = captures[(0, PERIOD_CONTEXTS[0])][0]
    baseline_control = captures[(0, "no_period")][0]
    if captures[(0, "no_period")][1] != 11803:
        raise ValueError("no-period control has unexpected frame count")
    split: int | None = None
    for value in VALUES:
        no_period, frames, _ = captures[(value, "no_period")]
        if frames != 11803 or no_period != baseline_control:
            raise ValueError(f"sentence pause changed no-period output at value {value}")
        reference_pcm, reference_frames, _ = captures[(value, PERIOD_CONTEXTS[0])]
        if reference_frames != 17762 + 16 * value:
            raise ValueError(f"unexpected period-context frame count at value {value}: {reference_frames}")
        for context in PERIOD_CONTEXTS[1:]:
            pcm, context_frames, _ = captures[(value, context)]
            if context_frames != reference_frames or pcm != reference_pcm:
                raise ValueError(f"period-context outputs differ at value {value}: {context}")
        if value == 0:
            if reference_pcm != baseline_period:
                raise ValueError("pause-zero reference mismatch")
            continue

        delta = 32 * value
        prefix = 0
        while prefix < min(len(baseline_period), len(reference_pcm)) and baseline_period[prefix] == reference_pcm[prefix]:
            prefix += 1
        if split is None:
            split = prefix
        if prefix != split or reference_pcm[prefix : prefix + delta] != bytes(delta):
            raise ValueError(f"pause {value} does not insert {delta} zero bytes at a stable position")
        if reference_pcm[prefix + delta :] != baseline_period[prefix:]:
            raise ValueError(f"pause {value} changes PCM outside the inserted silence interval")

    print(
        f"values={len(VALUES)} contexts={len(CONTEXTS)} calls={len(synth_rows)} "
        f"getter_rows={len(getter_rows)} period_split_byte={split}"
    )
    for value in (0, 1, 200, 925, 65535):
        print(
            f"pause={value} period_frames={captures[(value, 'upper')][1]} "
            f"pcm_sha256={hashlib.sha256(captures[(value, 'upper')][0]).hexdigest()} "
            f"no_period_frames={captures[(value, 'no_period')][1]}"
        )


if __name__ == "__main__":
    main()
