#!/usr/bin/env python3
"""Extract the documented Kate installer without executing vendor setup code."""

import argparse
import hashlib
import json
from pathlib import Path
import struct
import subprocess
import tempfile

INSTALLER_SHA256 = 'a1b1b0d189f8d9751e2d34f73fc9701dac695293e70c06c6ba9fdd1243b42377'
MSI_SHA256 = 'b3f5b3022a02037d83352878a279695e5b095b37031b1df79959c37563185351'
FILE_MAP = Path(__file__).resolve().parents[1] / 'work/trace-reproduction/kate-installer-files.json'
ROTATE = bytes(((value << 4) | (value >> 4)) & 255 for value in range(256))


def decode_block(raw, filename, offset):
    key = bytes(value ^ (0x13, 0x35, 0x86, 7)[index % 4] for index, value in enumerate(filename))
    if not key:
        raise ValueError('empty embedded filename')
    sequence = bytes(key[index % len(key)] ^ 255 for index in range(1024))
    start = offset % 1024
    stream = (sequence * ((len(raw) + start + 1023) // 1024))[start:start + len(raw)]
    rotated = raw.translate(ROTATE)
    return (int.from_bytes(rotated, 'little') ^ int.from_bytes(stream, 'little')).to_bytes(len(raw), 'little')


def extract_msi(installer, output):
    digest = hashlib.sha256()
    with installer.open('rb') as source:
        while block := source.read(1 << 20):
            digest.update(block)
    if digest.hexdigest() != INSTALLER_SHA256:
        raise ValueError('installer SHA-256 differs from the documented source')
    with installer.open('rb') as source:
        source.seek(0x3e000)
        header = source.read(46)
        if header[:13] != b'InstallShield' or struct.unpack_from('<H', header, 14)[0] != 3:
            raise ValueError('unexpected installer overlay header')
        found = False
        for _ in range(3):
            attributes = source.read(312)
            if len(attributes) != 312:
                raise ValueError('truncated embedded-file attributes')
            name = attributes[:260].split(b'\0', 1)[0]
            flags, size = struct.unpack_from('<I', attributes, 260)[0], struct.unpack_from('<I', attributes, 268)[0]
            if name != b'NextUp.com-NeoSpeech Kate16 Voice.msi':
                source.seek(size, 1)
                continue
            found = True
            with output.open('xb') as target:
                offset = 0
                while offset < size:
                    raw = source.read(min(1 << 20, size - offset))
                    if not raw:
                        raise ValueError('truncated embedded MSI')
                    target.write(decode_block(raw, name, offset) if flags & 6 == 6 else raw)
                    offset += len(raw)
        if not found or hashlib.sha256(output.read_bytes()).hexdigest() != MSI_SHA256:
            raise ValueError('embedded MSI differs from the documented SHA-256')


def extract_package(installer, destination):
    repository = Path(__file__).resolve().parents[3]
    if destination.resolve().is_relative_to(repository):
        relative = destination.resolve().relative_to(repository)
        if not relative.parts or relative.parts[0] == 'binary' or relative.parts[0].startswith('data-'):
            raise ValueError('vendor input directories are read-only')
    if destination.exists():
        raise ValueError(f'refusing to overwrite existing package: {destination}')
    mapping = json.loads(FILE_MAP.read_text())
    destination.mkdir(parents=True)
    with tempfile.TemporaryDirectory(prefix='vtspeak-kate-install-', dir=destination.parent) as directory:
        scratch = Path(directory)
        msi = scratch / 'Kate.msi'
        extract_msi(installer, msi)
        streams = scratch / 'streams'
        subprocess.run(['7z', 'x', '-y', f'-o{streams}', str(msi)], check=True, capture_output=True)
        cabinets = []
        for path in streams.iterdir():
            if path.is_file():
                with path.open('rb') as stream:
                    if stream.read(4) == b'MSCF':
                        cabinets.append(path)
        if len(cabinets) != 1:
            raise ValueError('expected one embedded MSI cabinet')
        files = scratch / 'files'
        subprocess.run(['7z', 'x', '-y', f'-o{files}', str(cabinets[0])], check=True, capture_output=True)
        content = {hashlib.sha256(p.read_bytes()).hexdigest(): p for p in files.rglob('*') if p.is_file()}
        for relative, digest in mapping.items():
            path = Path(relative)
            if path.is_absolute() or '..' in path.parts:
                raise ValueError('unsafe mapped package path')
            data = content[digest].read_bytes()
            target = destination / path
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(data)
    return len(mapping)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('installer', type=Path)
    parser.add_argument('destination', type=Path)
    args = parser.parse_args()
    print(f'{extract_package(args.installer, args.destination)} package files extracted and verified')


if __name__ == '__main__':
    main()
