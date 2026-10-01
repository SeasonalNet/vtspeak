#!/usr/bin/env python3
"""Check the captured case-matching effects of [CI] in P dictionary rows."""

from __future__ import annotations

from pathlib import Path
import re
import wave


ROOT = Path(__file__).resolve().parent
MODEL = re.compile(
    r"^JAMES_DICT model_load_ax=(-?\d+) sample_speaker=(-?\d+) "
    r"gate=(\d+) capacity=(-?\d+)$",
    re.MULTILINE,
)
DICTIONARY = re.compile(
    r"^JAMES_DICT case=(p|a) load=(-?\d+) synth=(-?\d+) unload=(-?\d+)$",
    re.MULTILINE,
)
RUNS = {
    "upper_adjacent": "james-licensed-userdict-ci-source-case-v1-v3-api.log",
    "upper_spaced": "james-licensed-userdict-ci-source-case-spaced-v1-api.log",
    "mixed_adjacent": "james-licensed-userdict-ci-source-case-mixed-v1-api.log",
    "lower_tag_upper_source": "james-licensed-userdict-ci-source-case-lower-tag-v1-api.log",
    "lower_tag_lower_source": "james-licensed-userdict-ci-source-case-lower-tag-lower-source-v1-api.log",
    "mixed_tag_cI": "james-licensed-userdict-ci-tag-cI-case-v2-api.log",
    "mixed_tag_Ci": "james-licensed-userdict-ci-tag-ci-case-v3-api.log",
}


def read_pcm(stem: str) -> tuple[bytes, int]:
    path = ROOT / f"{stem}.wav"
    with wave.open(str(path), "rb") as source:
        props = (source.getnchannels(), source.getframerate(), source.getsampwidth())
        if props != (1, 16000, 2):
            raise ValueError(f"{path.name} has unexpected WAVE properties: {props}")
        frames = source.getnframes()
        pcm = source.readframes(frames)
    if len(pcm) != frames * 2:
        raise ValueError(f"{path.name} has an unexpected PCM length")
    return pcm, frames


def check_log(name: str, filename: str) -> None:
    log = (ROOT / filename).read_text(encoding="utf-8", errors="replace")
    if MODEL.findall(log) != [("0", "-1", "1", "6")]:
        raise ValueError(f"{name}: model/license state did not match the capture")
    rows = {
        key: tuple(map(int, values))
        for key, *values in DICTIONARY.findall(log)
    }
    if rows != {"p": (1, 1, 1), "a": (1, 1, 1)}:
        raise ValueError(f"{name}: dictionary call results were {rows}")


def check_run(
    name: str,
    prefix: str,
    *,
    plain_matches: bool = False,
) -> None:
    check_log(name, RUNS[name])
    control, control_frames = read_pcm(f"{prefix}control")
    plain, plain_frames = read_pcm(f"{prefix}p")
    marked, marked_frames = read_pcm(f"{prefix}a")
    restored, restored_frames = read_pcm(f"{prefix}restored")
    if control != restored or control_frames != restored_frames:
        raise ValueError(f"{name}: unloading the row did not restore control PCM")
    if (plain != control) != plain_matches or marked == control:
        raise ValueError(f"{name}: unexpected plain/marked match behavior")
    if marked_frames != 1150:
        raise ValueError(f"{name}: marked output had {marked_frames} frames")
    print(
        f"{name}: control={control_frames} plainP={plain_frames} "
        f"marked={marked_frames} plain_matches={plain_matches} PASS"
    )


def main() -> None:
    lowercase_plain, lowercase_frames = read_pcm("james-userdict-p")
    if lowercase_frames != 1150:
        raise ValueError(f"plain lowercase P output had {lowercase_frames} frames")
    check_run("upper_adjacent", "james-userdict-ci-case-")
    check_run("upper_spaced", "james-userdict-ci-case-spaced-")
    check_run("mixed_adjacent", "james-userdict-ci-case-mixed-")
    check_run("lower_tag_upper_source", "james-userdict-ci-case-lower-tag-")
    check_run(
        "lower_tag_lower_source",
        "james-userdict-ci-case-lower-tag-lower-source-",
        plain_matches=True,
    )
    check_run("mixed_tag_cI", "james-userdict-ci-tag-cI-case-v2-")
    check_run("mixed_tag_Ci", "james-userdict-ci-tag-ci-case-v3-")
    for prefix in (
        "james-userdict-ci-case-",
        "james-userdict-ci-case-spaced-",
        "james-userdict-ci-case-mixed-",
        "james-userdict-ci-case-lower-tag-lower-source-",
        "james-userdict-ci-tag-cI-case-v2-",
        "james-userdict-ci-tag-ci-case-v3-",
    ):
        marked, _ = read_pcm(f"{prefix}a")
        if marked != lowercase_plain:
            raise ValueError(f"{prefix}a did not match the plain lowercase P output")
    print("all four [CI] marker case masks perform tested ASCII case-insensitive source matching")


if __name__ == "__main__":
    main()
