#!/usr/bin/env python3
"""Validate file-API pause-argument precedence captures from Stage 21."""

from __future__ import annotations

from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
AXES = ("sentence", "comma")
STORED = (0, 925)
CALLS = (-1, 0, 250)
CONTEXTS = ("period", "comma", "control")
ROWS = re.compile(
    r"^PAUSE_PRECEDENCE axis=(\w+) stored=(\d+) call=(-?\d+) "
    r"context=(\w+) ret=(-?\d+)$",
    re.MULTILINE,
)


def pcm(axis: str, stored: int, call: int, context: str) -> tuple[bytes, int]:
    path = ROOT / f"pause-precedence-{axis}-{stored}-{call}-{context}.wav"
    with wave.open(str(path), "rb") as source:
        if (source.getnchannels(), source.getframerate(), source.getsampwidth()) != (
            1,
            16000,
            2,
        ):
            raise ValueError(f"unexpected WAVE format: {path}")
        frames = source.getnframes()
        samples = source.readframes(frames)
    if len(samples) != frames * 2:
        raise ValueError(f"invalid PCM size: {path}")
    return samples, frames


def assert_inserted_pause(actual: bytes, baseline: bytes, duration_ms: int) -> None:
    count = duration_ms * 16
    split = 18006
    inserted = bytes(count * 2)
    if actual != baseline[:split] + inserted + baseline[split:]:
        raise ValueError(f"pause {duration_ms} ms is not an exact zero insertion")


def main() -> None:
    log = (ROOT / "pause-precedence-grid-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    rows = ROWS.findall(log)
    expected = {
        (axis, stored, call, context)
        for axis in AXES
        for stored in STORED
        for call in CALLS
        for context in CONTEXTS
    }
    observed = {
        (axis, int(stored), int(call), context): int(result)
        for axis, stored, call, context, result in rows
    }
    if set(observed) != expected or len(rows) != len(expected):
        raise ValueError(f"expected {len(expected)} unique runtime rows, got {len(rows)}")
    if any(result != 1 for result in observed.values()):
        raise ValueError("one or more file API calls failed")

    baseline, base_frames = pcm("sentence", 0, 0, "period")
    if base_frames != 17762:
        raise ValueError(f"unexpected period baseline: {base_frames}")
    for axis in AXES:
        for stored in STORED:
            for call in CALLS:
                period_pcm, period_frames = pcm(axis, stored, call, "period")
                comma_pcm, comma_frames = pcm(axis, stored, call, "comma")
                control_pcm, control_frames = pcm(axis, stored, call, "control")
                sentence_stored = stored if axis == "sentence" else 200
                comma_stored = 200 if axis == "sentence" else stored
                period_ms = sentence_stored if call == -1 else call
                comma_ms = comma_stored
                expected_frames = (
                    17762 + 16 * period_ms,
                    17762 + 16 * comma_ms,
                    11803,
                )
                if (period_frames, comma_frames, control_frames) != expected_frames:
                    raise ValueError(
                        f"frame mismatch for {axis}/{stored}/{call}: "
                        f"{period_frames, comma_frames, control_frames} != {expected_frames}"
                    )
                if control_pcm != pcm("sentence", 0, 0, "control")[0]:
                    raise ValueError("control PCM changed across pause settings")
                assert_inserted_pause(period_pcm, baseline, period_ms)
                comma_baseline, _ = pcm("comma", 0, 0, "comma")
                assert_inserted_pause(comma_pcm, comma_baseline, comma_ms)

    print(f"axes={len(AXES)} stored_values={len(STORED)} call_values={len(CALLS)}")
    print(f"contexts={len(CONTEXTS)} calls={len(expected)} returns=1 exact_pcm_relations=validated")
    for axis in AXES:
        for stored in STORED:
            rows_out = []
            for call in CALLS:
                frames = tuple(
                    pcm(axis, stored, call, context)[1] for context in CONTEXTS
                )
                rows_out.append(f"call={call}:{frames}")
            print(f"{axis} stored={stored} " + " ".join(rows_out))


if __name__ == "__main__":
    main()
