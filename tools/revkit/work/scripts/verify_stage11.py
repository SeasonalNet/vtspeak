#!/usr/bin/env python3
"""Verify local Stage 11 WAVs and compare the selected Wag peer manifest rows."""

from __future__ import annotations

import csv
import hashlib
import wave
from pathlib import Path


WORK = Path(__file__).resolve().parents[1]
STAGE11 = WORK / "stage11"
PEER_MANIFEST = STAGE11 / "peer-golden.tsv"
LOCAL_MANIFEST = STAGE11 / "local-golden.tsv"


def read_wave(path: Path) -> tuple[int, bytes]:
    with wave.open(str(path), "rb") as wav:
        params = wav.getparams()
        expected = (1, 2, 16000, "NONE")
        actual = (params.nchannels, params.sampwidth, params.framerate, params.comptype)
        if actual != expected:
            raise ValueError(f"unexpected PCM format in {path}: {params}")
        return params.nframes, wav.readframes(params.nframes)


def main() -> None:
    with PEER_MANIFEST.open(encoding="utf-8", newline="") as stream:
        peer_rows = {row["fixture"]: row for row in csv.DictReader(stream, delimiter="\t")}

    output_rows: list[dict[str, str]] = []
    for fixture, peer in peer_rows.items():
        input_path = STAGE11 / "inputs" / f"{fixture}.txt"
        if input_path.read_bytes() != bytes.fromhex(peer["input_cp1252_hex"]):
            raise SystemExit(f"input bytes differ from peer manifest: {fixture}")
        wav_path = STAGE11 / f"{fixture}.wav"
        blocks_path = STAGE11 / f"{fixture}-blocks.pcm"
        if not wav_path.is_file() or not blocks_path.is_file():
            continue
        frames, pcm = read_wave(wav_path)
        blocks = blocks_path.read_bytes()
        if blocks != pcm:
            raise SystemExit(f"PCM block mismatch: {fixture}: {len(blocks)} != {len(pcm)} bytes")
        wav_bytes = wav_path.read_bytes()
        digest = hashlib.sha256(wav_bytes).hexdigest()
        peer_match = str(
            len(wav_bytes) == int(peer["expected_wav_bytes"])
            and digest == peer["expected_wav_sha256"]
        ).lower()
        print(
            f"{fixture}: frames={frames} wav_bytes={len(wav_bytes)} "
            f"sha256={digest} peer_match={peer_match} pcm_blocks_match=true"
        )
        output_rows.append(
            {
                "fixture": fixture,
                "frames": str(frames),
                "wav_bytes": str(len(wav_bytes)),
                "wav_sha256": digest,
                "peer_expected_wav_bytes": peer["expected_wav_bytes"],
                "peer_expected_wav_sha256": peer["expected_wav_sha256"],
                "peer_hash_match": peer_match,
                "pcm_blocks_match_wav_data": "true",
            }
        )

    if len(output_rows) != len(peer_rows):
        raise SystemExit(f"captured {len(output_rows)} of {len(peer_rows)} peer golden fixtures")
    with LOCAL_MANIFEST.open("w", encoding="utf-8", newline="") as stream:
        writer = csv.DictWriter(stream, fieldnames=list(output_rows[0]), delimiter="\t")
        writer.writeheader()
        writer.writerows(output_rows)


if __name__ == "__main__":
    main()
