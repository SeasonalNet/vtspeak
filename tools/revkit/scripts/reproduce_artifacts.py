#!/usr/bin/env python3
"""Rebuild disposable GDB command files and WAV copies from tracked sources.

This does not run VoiceText or reproduce a runtime observation. It reproduces
the exact command inputs and audio files identified by the checked-in manifest.
"""

import argparse
import hashlib
import io
import json
import struct
from pathlib import Path
import subprocess
import sys
import tempfile
import wave


REPOSITORY = Path(__file__).resolve().parents[3]
MANIFEST = Path("tools/revkit/work/trace-reproduction/manifest.json")


def read_input(root, relative, offset, size, digest):
    """Read a bounded, hash-checked vendor or research input without writing it."""
    path = Path(relative)
    if path.is_absolute() or ".." in path.parts or not path.parts:
        raise ValueError(f"unsafe input path: {relative}")
    if path.parts[0] not in {"data-kate", "data-paul"} and path.parts[:3] != ("tools", "revkit", "work"):
        raise ValueError(f"unsupported audio input: {relative}")
    target = root / path
    if not target.resolve().is_relative_to(root.resolve()) or offset < 0 or size < 0:
        raise ValueError(f"unsafe input span: {relative}")
    with target.open("rb") as stream:
        stream.seek(offset)
        data = stream.read(size)
    if len(data) != size or hashlib.sha256(data).hexdigest() != digest:
        raise ValueError(f"audio source differs: {relative}")
    return data


def decode_segments(root, segments):
    from decode_dat import decode_payload

    output = []
    for segment in segments:
        data = read_input(root, segment["dat"], segment["offset"], segment["bytes"], segment["sha256"])
        samples = decode_payload(data)
        start, count = segment["start"], segment["count"]
        if start < 0 or count < 0 or start + count > len(samples):
            raise ValueError("decoded audio slice is out of bounds")
        output.extend(samples[start:start + count])
    return struct.pack(f"<{len(output)}h", *output)


def wave_segments(root, manifest, segments, active, runtime_root=None):
    output = []
    for segment in segments:
        if "zero_samples" in segment:
            if segment["zero_samples"] < 0:
                raise ValueError("negative silence length")
            output.append(b"\0\0" * segment["zero_samples"])
            continue
        source = segment["wave"]
        data = build(root, manifest, "wav", source, active, runtime_root) if source in manifest["wav"] else safe_path(root, source).read_bytes()
        with wave.open(io.BytesIO(data), "rb") as audio:
            if (audio.getnchannels(), audio.getsampwidth(), audio.getframerate()) != (1, 2, 16000):
                raise ValueError(f"unsupported source WAV format: {source}")
            start, count = segment["start"], segment["count"]
            if start < 0 or count < 0 or start + count > audio.getnframes():
                raise ValueError("WAV sample slice is out of bounds")
            audio.setpos(start)
            output.append(audio.readframes(count))
    return b"".join(output)


def safe_path(root, relative):
    """Keep artifact operations within the analysis workspace, including links."""
    relative = Path(relative)
    if relative.is_absolute() or ".." in relative.parts:
        raise ValueError(f"unsafe relative path: {relative}")
    if relative.parts[:3] != ("tools", "revkit", "work"):
        raise ValueError(f"artifact is outside the analysis workspace: {relative}")
    target = root / relative
    if target.is_symlink() or not target.resolve().is_relative_to(root.resolve()):
        raise ValueError(f"unsafe artifact target: {relative}")
    return target


def expand_source(data):
    """Expand consecutive decimal byte writes in readable GDB source templates."""
    output = []
    for line in data.decode("utf-8").splitlines(keepends=True):
        if line.startswith("#@bytes "):
            prefix, start, suffix, values = json.loads(line[8:])
            output.extend(
                f"{prefix}{start + offset}{suffix}{value}\n"
                for offset, value in enumerate(bytes.fromhex(values))
            )
        elif line.startswith("#@literal "):
            output.append(json.loads(line[10:]))
        else:
            output.append(line)
    return "".join(output).encode("utf-8")


def checked(data, entry, path):
    if len(data) != entry["bytes"] or hashlib.sha256(data).hexdigest() != entry["sha256"]:
        raise ValueError(f"reproduction differs from recorded bytes: {path}")
    return data


def build(root, manifest, kind, path, active=(), runtime_root=None):
    entry = manifest[kind][path]
    if path in active:
        raise ValueError(f"cyclic artifact dependency: {path}")
    if kind == "gdb":
        if "generator" in entry:
            generator = safe_path(root, entry["generator"])
            result = subprocess.run(
                ["bash", str(generator), *entry["args"]],
                env={"PATH": "/usr/bin:/bin", "LC_ALL": "C", **entry["env"]},
                check=True,
                capture_output=True,
            )
            data = result.stdout
        elif "python_generator" in entry:
            generator = safe_path(root, entry["python_generator"])
            with tempfile.TemporaryDirectory(prefix="vtspeak-gdb-") as directory:
                output = Path(directory) / "trace.gdb"
                subprocess.run(
                    [sys.executable, str(generator), "generate", str(output)],
                    check=True,
                    capture_output=True,
                )
                data = output.read_bytes()
        else:
            data = expand_source(safe_path(root, entry["source"]).read_bytes())
    elif kind == "wav":
        if "runtime_batch" in entry:
            if runtime_root is None:
                source = safe_path(root, path)
            else:
                batch = Path(entry["runtime_batch"])
                if batch.is_absolute() or ".." in batch.parts or len(batch.parts) != 1:
                    raise ValueError("unsafe runtime batch name")
                source = safe_path(runtime_root / batch, entry.get("runtime_source", path))
            if not source.is_file():
                raise ValueError(f"run reproduce_wavs.py run --batch {entry['runtime_batch']} first, then pass --runtime-root to this tool: {path}")
            data = source.read_bytes()
        elif "copy" in entry:
            source = entry["copy"]
            if source in manifest["wav"]:
                data = build(root, manifest, kind, source, (*active, path), runtime_root)
            else:
                data = safe_path(root, source).read_bytes()
        else:
            if "wave_segments" in entry:
                payload = wave_segments(root, manifest, entry["wave_segments"], (*active, path), runtime_root)
            elif "dat_segments" in entry:
                payload = decode_segments(root, entry["dat_segments"])
            elif "pcm" in entry:
                payload = b"".join(safe_path(root, source).read_bytes() for source in entry["pcm"])
            else:
                payload = bytes([entry["repeat"]]) * entry["count"]
            data = bytes.fromhex(entry["prefix"]) + payload + bytes.fromhex(entry["suffix"])
    else:
        raise ValueError(f"unknown artifact kind: {kind}")
    return checked(data, entry, path)


def restore(root, path, data):
    target = safe_path(root, path)
    if target.exists():
        if target.read_bytes() != data:
            raise ValueError(f"refusing to overwrite changed artifact: {path}")
        return
    target.parent.mkdir(parents=True, exist_ok=True)
    # Exclusive creation protects concurrent work between the existence check and write.
    with target.open("xb") as output:
        output.write(data)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("action", choices=("verify", "restore", "clean"))
    parser.add_argument("--kind", choices=("gdb", "wav", "all"), default="all")
    parser.add_argument("--path", action="append", help="select a manifest path; repeat to select several")
    parser.add_argument("--output-root", type=Path, default=REPOSITORY)
    parser.add_argument("--runtime-root", type=Path, help="scratch root populated by reproduce_wavs.py run")
    args = parser.parse_args()
    manifest = json.loads((REPOSITORY / MANIFEST).read_text())
    if manifest["version"] != 1:
        raise ValueError("unsupported artifact manifest version")
    kinds = ("gdb", "wav") if args.kind == "all" else (args.kind,)
    selected = set(args.path or [])
    known = set().union(*(manifest[kind] for kind in kinds))
    if selected - known:
        raise ValueError(f"unknown selected paths: {sorted(selected - known)}")
    for kind in kinds:
        count = total = 0
        verb = {"verify": "verified", "restore": "restored", "clean": "cleaned"}[args.action]
        for path, entry in manifest[kind].items():
            if selected and path not in selected:
                continue
            target = safe_path(args.output_root, path)
            if args.action == "clean":
                if target.exists():
                    checked(target.read_bytes(), entry, path)
                    target.unlink()
            else:
                data = build(REPOSITORY, manifest, kind, path, runtime_root=args.runtime_root)
                if args.action == "restore":
                    restore(args.output_root, path, data)
                elif target.exists() and target.read_bytes() != data:
                    raise ValueError(f"existing artifact differs from reproduction: {path}")
            count += 1
            total += entry["bytes"]
            if count % 100 == 0:
                print(f"{kind}: {count} files {verb}", flush=True)
        print(f"{kind}: {count} files, {total} bytes; {args.action} completed", flush=True)


if __name__ == "__main__":
    main()
