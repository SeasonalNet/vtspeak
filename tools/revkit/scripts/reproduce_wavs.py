#!/usr/bin/env python3
"""Reproduce catalogued WAVs in scratch storage and verify their exact bytes."""

import argparse
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import re
import shlex
import shutil
import subprocess

try:
    from . import reproduce_artifacts as artifacts
except ImportError:
    import reproduce_artifacts as artifacts

ROOT = artifacts.REPOSITORY
CATALOG = ROOT / 'tools/revkit/work/trace-reproduction/wav-runtime.json'


def check_inputs(index):
    for relative, digest in index.items():
        source = ROOT / relative_path(relative)
        with source.open('rb') as stream:
            if hashlib.file_digest(stream, 'sha256').hexdigest() != digest:
                raise ValueError(f'runtime prerequisite differs: {relative}')
    return len(index)


def relative_path(value):
    path = Path(value)
    if path.is_absolute() or '..' in path.parts or not path.parts:
        raise ValueError(f'unsafe scratch path: {value}')
    return path


def work_relative(value):
    path = relative_path(value)
    if path.parts[:3] != ('tools', 'revkit', 'work'):
        raise ValueError(f'not a research workspace path: {value}')
    return Path(*path.parts[3:])


def prepare(scratch, batch, manifest):
    work = scratch / 'work'
    (scratch / 'tools/revkit').mkdir(parents=True, exist_ok=True)
    (scratch / 'tools/revkit/work').symlink_to('../../work', target_is_directory=True)
    for mount in batch.get('mounts', []):
        if mount.get('prepare'):
            subprocess.run([arg.replace('{scratch}', str(scratch)) for arg in mount['prepare']['argv']], cwd=ROOT, check=True)
    for directory in ('home', 'stage5', 'stage2-copy/wineprefix', 'stage21/sandbox/stage5'):
        (work / directory).mkdir(parents=True, exist_ok=True)
    for folder in (work / 'stage5', work / 'stage21/sandbox/stage5'):
        for filename, data in (('input1.txt', b'Hello.\n'), ('output.wav', b'')):
            target = folder / filename
            if not target.exists():
                target.write_bytes(data)
    for voice in ('data-paul', 'data-common'):
        target = work / 'stage21/sandbox' / voice
        if not target.exists() and not target.is_symlink():
            target.symlink_to('../../' + voice)
    if batch.get('kate_adapter'):
        package = scratch / 'kate-package'
        shutil.copytree(ROOT / relative_path(batch.get('model_source', 'data-kate')), package)
        (package / 'M16/ttsdata/tree2').mkdir(parents=True, exist_ok=True)
        trees = scratch / 'kate-trees'
        module_path = ROOT / 'tools/revkit/work/stage19/convert_tree2.py'
        spec = importlib.util.spec_from_file_location('wav_tree_converter', module_path)
        converter = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(converter)
        source_root = ROOT / 'data-kate/M16/ttsdata/tree'
        for source in sorted(source_root.rglob('*.tree2')):
            legacy = converter.parse_tree2(source)
            encoded = converter.encode(legacy.root)
            converter.verify(source, encoded, legacy.root)
            destination = trees / source.relative_to(source_root)
            destination.parent.mkdir(parents=True, exist_ok=True)
            destination.write_bytes(encoded)
        common = scratch / 'common-dict'
        shutil.copytree(ROOT / 'data-common/dict-eng', common)
        for source in sorted(common.glob('*.tree3')):
            shutil.copyfile(source, source.with_suffix('.tree2'))
    for source in batch.get('sources', []):
        destination = work / work_relative(source)
        destination.parent.mkdir(parents=True, exist_ok=True)
        data = artifacts.safe_path(ROOT, source).read_bytes()
        if destination.exists() and destination.read_bytes() != data:
            raise ValueError(f'refusing to replace changed scratch source: {source}')
        if not destination.exists():
            destination.write_bytes(data)
    for path in batch.get('placeholders', []):
        target = work / work_relative(path)
        target.parent.mkdir(parents=True, exist_ok=True)
        target.touch(exist_ok=False)
    for source in batch.get('traces', []):
        trace_manifest = {'gdb': {**manifest['gdb'], **batch.get('trace_overrides', {})}}
        artifacts.restore(scratch, source, artifacts.build(ROOT, trace_manifest, 'gdb', source))
        # The runtime layout mounts tools/revkit/work directly as /work.
        generated = scratch / source
        destination = work / work_relative(source)
        destination.parent.mkdir(parents=True, exist_ok=True)
        if destination.exists() and destination.read_bytes() != generated.read_bytes():
            raise ValueError(f'refusing to replace changed scratch trace: {source}')
        if not destination.exists():
            destination.write_bytes(generated.read_bytes())
        for output in re.findall(rb'/work/([A-Za-z0-9._/-]+\.(?:log|pcm|dat|wav|txt))', generated.read_bytes()):
            (work / relative_path(output.decode('ascii'))).parent.mkdir(parents=True, exist_ok=True)
    for output, entry in batch['outputs'].items():
        destination = work / work_relative(output)
        destination.parent.mkdir(parents=True, exist_ok=True)
        generated = work / work_relative(entry.get('generated_path', output))
        generated.parent.mkdir(parents=True, exist_ok=True)
        if destination.exists():
            raise ValueError(f'output already exists; use a fresh scratch directory: {output}')
    lines = ['#!/bin/bash', 'set -euo pipefail']
    if batch.get('probes'):
        lines += ['Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >/work/xvfb-reproduction.log 2>&1 &',
                  'xpid=$!', "trap 'kill \"$xpid\" 2>/dev/null || true' EXIT", 'export DISPLAY=:99', 'sleep 1']
        for index, probe in enumerate(batch['probes']):
            source = '/work/' + str(work_relative(probe['input']))
            lines += [f'cp {shlex.quote(source)} input1.txt', ': > output.wav']
            trace = probe.get('trace')
            executable = shlex.quote(probe.get('executable', '/samples/voicetext_paul.exe'))
            command = f'timeout 300s wine {executable}'
            if trace:
                command = f'timeout 300s winedbg --gdb {executable} < ' + shlex.quote('/work/' + str(work_relative(trace)))
            lines.append(command + f' >/work/reproduction-probe-{index}.log 2>&1')
            if probe.get('copy_output'):
                lines.append('cp output.wav ' + shlex.quote('/work/' + str(work_relative(probe['copy_output']))))
    for command in batch.get('commands', []):
        env = command.get('env', {})
        lines.append(shlex.join(['env', *(f'{key}={value}' for key, value in env.items()), *command['argv']]))
    runner = work / 'reproduce-batch.sh'
    runner.write_text('\n'.join(lines) + '\n')
    return work


def docker_command(work, batch):
    overrides = {mount['target'] for mount in batch.get('mounts', [])}
    command = ['docker', 'run', '--rm', '--network', 'none', '--read-only', '--cap-drop', 'ALL',
               '--security-opt', 'no-new-privileges', '--user', f'{os.getuid()}:{os.getgid()}',
               '--tmpfs', '/tmp:rw,nosuid,nodev,size=256m', '-e', 'HOME=/work/home',
               '-e', 'WINEARCH=win32', '-e', 'WINEPREFIX=/work/stage2-copy/wineprefix',
               '-e', 'WINEDEBUG=-all', '-v', f'{work}:/work']
    for source, target in (('binary', '/samples'), ('include', '/voice-data/include')):
        if target not in overrides:
            command += ['-v', f'{ROOT}/{source}:{target}:ro']
    for name, alias in (('common', 'common'), ('paul', 'paul'), ('james', 'jame'),
                        ('james', 'james'), ('julie', 'julie'), ('julie', 'juli'),
                        ('kate', 'kate'), ('bridget', 'bridget')):
        if f'/work/data-{alias}' in overrides:
            continue
        source = work.parent / 'kate-package' if name == 'kate' and batch.get('kate_adapter') else ROOT / f'data-{name}'
        command += ['-v', f'{source}:/work/data-{alias}:ro']
    # Mount parents before children so a base tree cannot hide its pitch overlay.
    for mount in sorted(batch.get('mounts', []), key=lambda item: item['target'].count('/')):
        if 'work' in mount:
            source = work / relative_path(mount['work'])
            source.mkdir(parents=True, exist_ok=True)
            command += ['-v', f'{source}:{mount["target"]}']
            continue
        elif 'generated' in mount:
            source = work.parent / relative_path(mount['generated'])
        else:
            source = ROOT / relative_path(mount['source'])
        if not source.exists():
            raise ValueError(f'missing runtime prerequisite: {source}')
        command += ['-v', f'{source}:{mount["target"]}:ro']
    command += ['-w', '/work/' + str(relative_path(batch.get('cwd', 'stage5'))),
                batch.get('image', 'vtspeak-stage2-wine:latest'), '/bin/bash', '/work/reproduce-batch.sh']
    return command


def verify_batch(scratch, batch):
    failures = []
    for path, entry in batch['outputs'].items():
        try:
            data = (scratch / 'work' / work_relative(entry.get('generated_path', path))).read_bytes()
            artifacts.checked(data, entry, path)
            original = ROOT / path
            if original.exists() and original.read_bytes() != data:
                raise ValueError(f'existing WAV differs from fresh reproduction: {path}')
        except (ValueError, FileNotFoundError) as error:
            failures.append(str(error))
    if failures:
        raise ValueError('\n'.join(failures))
    return len(batch['outputs'])


def batch_input_index(catalog, names):
    """Combine the prerequisite indexes for the selected runtime batches."""
    index = {}
    for name in names:
        batch = catalog['batches'][name]
        source = batch.get('input_index', catalog['input_index'])
        index.update(json.loads((ROOT / relative_path(source)).read_text()))
    return index


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('action', choices=('list', 'run', 'verify', 'check-inputs'))
    parser.add_argument('--batch', action='append')
    parser.add_argument('--scratch', type=Path, default=Path('/tmp/vtspeak-wav-reproduction'))
    parser.add_argument('--image', help='override the recorded local Wine image; output hashes must still match')
    args = parser.parse_args()
    catalog = json.loads(CATALOG.read_text())
    selected = args.batch or list(catalog['batches'])
    unknown = set(selected) - catalog['batches'].keys()
    if unknown:
        raise ValueError(f'unknown batches: {sorted(unknown)}')
    if args.action in ('run', 'check-inputs'):
        index = batch_input_index(catalog, selected)
        print(f'{check_inputs(index)} runtime prerequisites verified', flush=True)
        if args.action == 'check-inputs':
            return
    manifest = json.loads((ROOT / artifacts.MANIFEST).read_text())
    for name in selected:
        batch = {**catalog['batches'][name], 'image': args.image or catalog['image']}
        if args.action == 'list':
            print(f'{name}: {len(batch["outputs"])} WAVs')
            continue
        scratch = args.scratch.resolve() / relative_path(name)
        if args.action == 'run':
            if scratch.exists():
                raise ValueError(f'use a fresh scratch directory: {scratch}')
            work = prepare(scratch, batch, manifest)
            try:
                subprocess.run(docker_command(work, batch), check=True)
                for path in batch.get('derived', []):
                    artifacts.restore(scratch, path, artifacts.build(scratch, manifest, 'wav', path))
                print(f'{name}: {verify_batch(scratch, batch)} WAVs verified byte for byte', flush=True)
            finally:
                # These machine inputs belong to this invocation's previously absent directory.
                shutil.rmtree(scratch / 'work/stage2-copy/wineprefix')
                if batch.get('kate_adapter'):
                    shutil.rmtree(scratch / 'kate-package')
        else:
            print(f'{name}: {verify_batch(scratch, batch)} WAVs verified byte for byte', flush=True)


if __name__ == '__main__':
    main()
