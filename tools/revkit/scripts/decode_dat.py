#!/usr/bin/env python3
"""Decode one VoiceText 2013 Paul DAT unit into a standalone PCM WAV."""

from __future__ import annotations

import argparse
import struct
import wave
from pathlib import Path

INDEX_HEADER_SIZE = 45
UNIT_RECORD_SIZE = 19
DEFAULT_FRAME_SAMPLES = 256
MAX_FRAME_SAMPLES = 1 << 16


class BitReader:
    """MSB-first reader for the DLL's unary quotient plus fixed-width remainder code."""

    def __init__(self, data: bytes):
        self.data = data
        self.bit_offset = 0

    def bit(self) -> int:
        if self.bit_offset >= len(self.data) * 8:
            raise ValueError("DAT bitstream ended before its terminator")
        byte = self.data[self.bit_offset // 8]
        shift = 7 - self.bit_offset % 8
        self.bit_offset += 1
        return (byte >> shift) & 1

    def rice(self, remainder_width: int) -> int:
        if not 0 <= remainder_width <= 31:
            raise ValueError(f"unsupported Rice remainder width: {remainder_width}")
        quotient = 0
        while self.bit() == 0:
            quotient += 1
            if quotient > 1 << 20:
                raise ValueError("DAT unary quotient exceeds safety limit")
        remainder = 0
        for _ in range(remainder_width):
            remainder = (remainder << 1) | self.bit()
        return (quotient << remainder_width) | remainder


def signed_residual(value: int) -> int:
    """Apply FUN_100020a0's even/odd signed fold."""
    half = value >> 1
    return half if value & 1 == 0 else -half - 1


def _signed_i32(value: int) -> int:
    value &= 0xFFFFFFFF
    return value - 0x100000000 if value & 0x80000000 else value


def _previous_sample(frame: list[int], history: list[int], distance: int) -> int:
    """Read a predictor sample across the current-frame boundary."""
    if len(frame) >= distance:
        return frame[-distance]
    return history[-(distance - len(frame))]


def decode_payload(data: bytes) -> list[int]:
    """Decode the original engine's framed predictive DAT bitstream.

    This models the static path in FUN_10001b30 and its mode handlers. It has
    been compared byte-for-byte with 27 unique original-DLL payloads, covering
    predictor modes 0-3. The zero-frame mode has not appeared in our runtime
    sample.
    """
    reader = BitReader(data)
    frame_samples = DEFAULT_FRAME_SAMPLES
    residual_width = 0
    output_shift = 0
    mean_history = [0, 0, 0, 0]
    history = [0, 0, 0]
    samples: list[int] = []
    max_output_samples = max(DEFAULT_FRAME_SAMPLES, len(data) * 128)

    while True:
        mode = reader.rice(2)
        if mode == 4:
            break
        if mode == 5:
            width = reader.rice(2)
            frame_samples = reader.rice(width)
            if not 1 <= frame_samples <= MAX_FRAME_SAMPLES:
                raise ValueError(f"invalid DAT frame size: {frame_samples}")
            continue
        if mode == 6:
            output_shift = reader.rice(2)
            if output_shift > 15:
                raise ValueError(f"invalid DAT output shift: {output_shift}")
            continue
        if mode not in (0, 1, 2, 3, 8):
            raise ValueError(f"unsupported DAT control value: {mode}")

        if mode != 8:
            residual_width = reader.rice(3)

        # FUN_10001b30 averages the four stored block means before dispatching
        # a frame. The shift reverses the bitshift applied when each mean was
        # stored. Keep the operations in the same order as the x86 arithmetic.
        mean_sum = _signed_i32(sum(mean_history) + 2)
        if mean_sum < 0:
            mean_sum = _signed_i32(mean_sum + 3)
        mean_offset = mean_sum >> 2
        if output_shift:
            mean_offset = (mean_offset >> (output_shift - 1)) >> 1

        decoded_frame: list[int] = []
        for _ in range(frame_samples):
            if mode == 8:
                value = 0
            else:
                residual = signed_residual(reader.rice(residual_width + 1))
                if mode == 0:
                    prediction = mean_offset
                elif mode == 1:
                    prediction = _previous_sample(decoded_frame, history, 1)
                elif mode == 2:
                    previous = _previous_sample(decoded_frame, history, 1)
                    previous2 = _previous_sample(decoded_frame, history, 2)
                    prediction = 2 * previous - previous2
                else:
                    previous = _previous_sample(decoded_frame, history, 1)
                    previous2 = _previous_sample(decoded_frame, history, 2)
                    previous3 = _previous_sample(decoded_frame, history, 3)
                    prediction = 3 * (previous - previous2) + previous3
                value = _signed_i32(prediction + residual)
            decoded_frame.append(value)

        # The shared mode tail updates the mean predictor and three-sample
        # history before shifting and clipping the frame written to PCM.
        mean_total = sum(decoded_frame) + frame_samples // 2
        rounded_mean = (abs(mean_total) // frame_samples) * (-1 if mean_total < 0 else 1)
        mean_history = mean_history[1:] + [
            _signed_i32(rounded_mean << output_shift)
        ]
        if len(decoded_frame) >= 3:
            history = decoded_frame[-3:]
        else:
            history = (history + decoded_frame)[-3:]

        for value in decoded_frame:
            scaled = _signed_i32(value << output_shift)
            if scaled > 0x7FFF:
                scaled = 0x7FFF
            # FUN_10002100 clips the positive side; the DLL then stores a
            # 16-bit word. Preserve that observed conversion for now.
            samples.append(((scaled + 0x8000) & 0xFFFF) - 0x8000)
        if len(samples) > max_output_samples:
            raise ValueError("decoded DAT output exceeded its safety limit")

    if not samples:
        raise ValueError("DAT payload terminated without producing samples")
    return samples


def unit_payload(index_path: Path, dat_path: Path, unit: int) -> bytes:
    raw_index = index_path.read_bytes()
    if len(raw_index) < INDEX_HEADER_SIZE or not raw_index[1:24].startswith(
        b"ver.2013\0VoiceText-Eng\0"
    ):
        raise ValueError(f"unsupported unit index: {index_path}")
    unit_count = struct.unpack_from("<I", raw_index, 39)[0]
    stride = struct.unpack_from("<H", raw_index, 43)[0]
    if stride != UNIT_RECORD_SIZE or not 0 <= unit < unit_count:
        raise ValueError(f"unit {unit} is outside the supported index range 0..{unit_count - 1}")
    offset_in_index = INDEX_HEADER_SIZE + unit * stride
    record = raw_index[offset_in_index : offset_in_index + stride]
    dat_offset = struct.unpack_from("<I", record, 0)[0]
    dat_length = struct.unpack_from("<H", record, 8)[0]
    end = dat_offset + dat_length
    if end > dat_path.stat().st_size:
        raise ValueError("unit DAT span exceeds the bank file")
    with dat_path.open("rb") as stream:
        stream.seek(dat_offset)
        payload = stream.read(dat_length)
    if len(payload) != dat_length:
        raise ValueError("truncated unit DAT span")
    return payload


def write_wav(path: Path, samples: list[int]) -> None:
    with wave.open(str(path), "wb") as output:
        output.setnchannels(1)
        output.setsampwidth(2)
        output.setframerate(16_000)
        output.writeframes(struct.pack(f"<{len(samples)}h", *samples))


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--index", type=Path, required=True, help="versioned VoiceText unit index")
    parser.add_argument("--dat", type=Path, required=True, help="matching merged-*.dat bank")
    parser.add_argument("--unit", type=int, required=True, help="zero-based unit record number")
    parser.add_argument("--output", type=Path, required=True, help="output 16 kHz mono PCM WAV")
    args = parser.parse_args()
    payload = unit_payload(args.index, args.dat, args.unit)
    samples = decode_payload(payload)
    write_wav(args.output, samples)
    print(f"decoded {len(samples)} samples from {len(payload)} DAT bytes to {args.output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
