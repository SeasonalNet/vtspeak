#!/usr/bin/env python3
"""Check that negative pause setter inputs preserve stored runtime state."""

from __future__ import annotations

from pathlib import Path
import re


ROOT = Path(__file__).resolve().parent
FIELDS = ("pitch", "speed", "volume", "sentence", "comma")
NEGATIVE = (-2147483648, -2, -1)
EXPECTED_STATE = (123, 234, 345, 456, 567)
BASELINE = re.compile(
    r"^NEG_PAUSE stage=baseline getter_ret=(\d+) pitch=(\d+) speed=(\d+) "
    r"volume=(\d+) sentence=(\d+) comma_ret=(\d+) comma=(\d+)(?:\\n)?$",
    re.MULTILINE,
)
ROW = re.compile(
    r"^NEG_PAUSE field=(\w+) input=(-?\d+) getter_ret=(\d+) "
    r"pitch=(\d+) speed=(\d+) volume=(\d+) sentence=(\d+) "
    r"comma_ret=(\d+) comma=(\d+)$",
    re.MULTILINE,
)


def main() -> None:
    log = (ROOT / "negative-pause-setter-grid-api.log").read_text(
        encoding="utf-8", errors="replace"
    )
    baseline = BASELINE.findall(log)
    if len(baseline) != 1:
        raise ValueError(f"expected one baseline row, got {len(baseline)}")
    base = tuple(map(int, baseline[0]))
    if base != (1, *EXPECTED_STATE[:4], 1, EXPECTED_STATE[4]):
        raise ValueError(f"unexpected initial getter state: {base}")

    rows = ROW.findall(log)
    expected = {(field, value) for field in FIELDS for value in NEGATIVE}
    observed = {(field, int(value)) for field, value, *_ in rows}
    if len(rows) != len(expected) or observed != expected:
        raise ValueError(f"expected {len(expected)} unique input rows, got {len(rows)}")
    for field, value, getter_ret, pitch, speed, volume, sentence, comma_ret, comma in rows:
        state = (int(pitch), int(speed), int(volume), int(sentence), int(comma))
        if int(getter_ret) != 1 or int(comma_ret) != 1 or state != EXPECTED_STATE:
            raise ValueError(
                f"negative {field}={value} changed state or getter status: "
                f"{getter_ret, comma_ret, state}"
            )

    print(f"fields={len(FIELDS)} negative_values={len(NEGATIVE)} calls={len(rows)}")
    print(f"baseline={EXPECTED_STATE} all_getters=1 all_negative_inputs_preserved_state=true")


if __name__ == "__main__":
    main()
