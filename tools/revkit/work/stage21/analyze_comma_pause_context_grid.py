#!/usr/bin/env python3
"""Validate comma-pause setter effects on Stage 21 WAVE captures."""

from __future__ import annotations

import hashlib
from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
VALUES = (0, 1, 199, 200, 201, 249, 250, 251, 499, 500, 501, 924, 925, 926, 65534, 65535)
COMMA_CONTEXTS = ("space", "adjacent", "double", "tab")
CONTEXTS = (*COMMA_CONTEXTS, "no_comma")
GETTER = re.compile(
    r"^COMMA_PAUSE value=(\d+) getter_ret=(\d+) getter_value=(\d+)$",
    re.MULTILINE,
)
SYNTH = re.compile(
    r"^COMMA_PAUSE value=(\d+) context=(\S+) synth_ret=(\d+)$",
    re.MULTILINE,
)


def read_wave(value: int, context: str) -> tuple[bytes, int, str]:
    path = ROOT / f"comma-pause-{value}-{context}.wav"
    with wave.open(str(path), "rb") as source:
        if (source.getnchannels(), source.getframerate(), source.getsampwidth()) != (1, 16000, 2):
            raise ValueError(f"unexpected WAVE format: {path}")
        frames = source.getnframes()
        pcm = source.readframes(frames)
    if len(pcm) != frames * 2:
        raise ValueError(f"PCM byte count differs from frame count: {path}")
    return pcm, frames, hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    log_path = ROOT / "comma-pause-context-grid-api.log"
    log = log_path.read_text(encoding="utf-8", errors="replace")
    getter_rows = GETTER.findall(log)
    synth_rows = SYNTH.findall(log)
    if len(getter_rows) != len(VALUES):
        raise ValueError(f"expected {len(VALUES)} getter rows, got {len(getter_rows)}")
    if len(synth_rows) != len(VALUES) * len(CONTEXTS):
        raise ValueError(f"expected {len(VALUES) * len(CONTEXTS)} synthesis rows, got {len(synth_rows)}")
    for value_text, return_text, observed_text in getter_rows:
        value = int(value_text)
        if value not in VALUES or int(return_text) != 1 or int(observed_text) != value:
            raise ValueError(f"getter mismatch: {value_text} {return_text} {observed_text}")
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
    for value in VALUES:
        control_pcm, control_frames, _ = captures[(value, "no_comma")]
        if control_frames != 11803:
            raise ValueError(f"no-comma control has {control_frames} frames at pause {value}")
        if control_pcm != captures[(0, "no_comma")][0]:
            raise ValueError(f"no-comma PCM changed at pause {value}")

        reference_pcm, reference_frames, _ = captures[(value, COMMA_CONTEXTS[0])]
        if reference_frames != 17762 + 16 * value:
            raise ValueError(f"unexpected comma duration at pause {value}: {reference_frames} frames")
        for context in COMMA_CONTEXTS[1:]:
            pcm, frames, _ = captures[(value, context)]
            if frames != reference_frames or pcm != reference_pcm:
                raise ValueError(f"comma whitespace changes output at pause {value}: {context}")

    base_pcm = captures[(0, "space")][0]
    split: int | None = None
    for value in VALUES[1:]:
        pcm = captures[(value, "space")][0]
        delta = 32 * value
        prefix = 0
        while prefix < min(len(base_pcm), len(pcm)) and base_pcm[prefix] == pcm[prefix]:
            prefix += 1
        if split is None:
            split = prefix
        if prefix != split or pcm[prefix : prefix + delta] != bytes(delta):
            raise ValueError(f"pause {value} does not insert {delta} zero bytes at a stable position")
        if pcm[prefix + delta :] != base_pcm[prefix:]:
            raise ValueError(f"pause {value} changes PCM outside the inserted silence interval")
        for context in COMMA_CONTEXTS[1:]:
            if captures[(value, context)][0] != pcm:
                raise ValueError(f"context {context} differs at pause {value}")

    print(
        f"values={len(VALUES)} contexts={len(CONTEXTS)} calls={len(synth_rows)} "
        f"getter_rows={len(getter_rows)} comma_split_byte={split}"
    )
    for value in (0, 1, 200, 925, 65535):
        print(
            f"pause={value} comma_frames={captures[(value, 'space')][1]} "
            f"pcm_sha256={hashlib.sha256(captures[(value, 'space')][0]).hexdigest()} "
            f"no_comma_frames={captures[(value, 'no_comma')][1]}"
        )


if __name__ == "__main__":
    main()
