import hashlib
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import wave

from scripts import extract_kate_installer as installer
from scripts import reproduce_artifacts as artifacts
from scripts import reproduce_wavs as runtime


def audio(samples):
    output = io.BytesIO()
    with wave.open(output, 'wb') as stream:
        stream.setnchannels(1)
        stream.setsampwidth(2)
        stream.setframerate(16000)
        stream.writeframes(samples)
    return output.getvalue()


class WavReproductionTests(unittest.TestCase):
    def test_runtime_seed_restores_a_copy_and_rejects_changed_audio(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            scratch = root / 'scratch'
            path = 'tools/revkit/work/example/output.wav'
            source = 'tools/revkit/work/example/generated.wav'
            data = audio(b'\x01\x00\x02\x00')
            entry = {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest(),
                     'runtime_batch': 'batch', 'runtime_source': source}
            manifest = {'wav': {path: entry, path + '.copy': {
                'bytes': len(data), 'sha256': entry['sha256'], 'copy': path}}}
            artifacts.restore(scratch / 'batch', source, data)
            self.assertEqual(artifacts.build(root, manifest, 'wav', path, runtime_root=scratch), data)
            self.assertEqual(artifacts.build(root, manifest, 'wav', path + '.copy', runtime_root=scratch), data)
            (scratch / 'batch' / source).write_bytes(b'changed')
            with self.assertRaisesRegex(ValueError, 'differs'):
                artifacts.build(root, manifest, 'wav', path, runtime_root=scratch)
            entry['runtime_batch'] = '../outside'
            with self.assertRaisesRegex(ValueError, 'unsafe runtime'):
                artifacts.build(root, manifest, 'wav', path, runtime_root=scratch)

    def test_input_hash_check_rejects_drift_and_path_traversal(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / 'fixture.txt').write_bytes(b'fixture')
            index = {'fixture.txt': hashlib.sha256(b'fixture').hexdigest()}
            with patch.object(runtime, 'ROOT', root):
                self.assertEqual(runtime.check_inputs(index), 1)
                (root / 'fixture.txt').write_bytes(b'changed')
                with self.assertRaisesRegex(ValueError, 'prerequisite differs'):
                    runtime.check_inputs(index)
                with self.assertRaisesRegex(ValueError, 'unsafe scratch path'):
                    runtime.check_inputs({'../outside': 'unused'})

    def test_embedded_file_decoding_handles_the_1024_byte_key_reset(self):
        name = b'example.msi'
        plain = bytes(range(256)) * 9
        encoded = []
        for index, value in enumerate(plain):
            key_index = index % 1024 % len(name)
            key = name[key_index] ^ (0x13, 0x35, 0x86, 7)[key_index % 4]
            rotated = value ^ key ^ 255
            encoded.append(((rotated << 4) | (rotated >> 4)) & 255)
        encoded = bytes(encoded)
        self.assertEqual(installer.decode_block(encoded, name, 0), plain)
        self.assertEqual(installer.decode_block(encoded[1019:], name, 1019), plain[1019:])
        with self.assertRaises(ValueError):
            installer.decode_block(b'data', b'', 0)

    def test_installer_rejects_wrong_input_hash_and_vendor_destination(self):
        with tempfile.TemporaryDirectory() as directory:
            source = Path(directory) / 'setup.exe'
            source.write_bytes(b'not the known installer')
            with self.assertRaisesRegex(ValueError, 'SHA-256'):
                installer.extract_msi(source, Path(directory) / 'output.msi')
            with self.assertRaisesRegex(ValueError, 'read-only'):
                installer.extract_package(source, artifacts.REPOSITORY / 'data-kate/forbidden')

    def test_audio_splice_and_silence_use_sample_counts_and_reject_bounds(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            folder = root / 'tools/revkit/work/example'
            folder.mkdir(parents=True)
            (folder / 'input.wav').write_bytes(audio(b'\x01\x00\x02\x00\x03\x00'))
            segments = [{'wave': 'tools/revkit/work/example/input.wav', 'start': 1, 'count': 2},
                        {'zero_samples': 2}]
            self.assertEqual(artifacts.wave_segments(root, {'wav': {}}, segments, ()), b'\x02\x00\x03\x00\0\0\0\0')
            with self.assertRaisesRegex(ValueError, 'bounds'):
                artifacts.wave_segments(root, {'wav': {}}, [{**segments[0], 'count': 3}], ())

    def test_splice_preserves_header_and_audio_beyond_declared_data_length(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            short = audio(b'\x01\x00')
            long = audio(b'\x02\x00\x03\x00')
            first = 'tools/revkit/work/example/first.wav'
            second = 'tools/revkit/work/example/second.wav'
            artifacts.restore(root, first, long)
            artifacts.restore(root, second, short)
            expected = short[:44] + long[44:] + short[44:]
            target = 'tools/revkit/work/example/mixed.wav'
            entry = {'bytes': len(expected), 'sha256': hashlib.sha256(expected).hexdigest(),
                     'prefix': short[:44].hex(), 'suffix': '', 'wave_segments': [
                         {'wave': first, 'start': 0, 'count': 2},
                         {'wave': second, 'start': 0, 'count': 1}]}
            actual = artifacts.build(root, {'wav': {target: entry}}, 'wav', target)
            self.assertEqual(actual, expected)
            with wave.open(io.BytesIO(actual), 'rb') as stream:
                self.assertEqual(stream.getnframes(), 1)
            self.assertEqual(len(actual), 50)

    def test_runtime_mounts_are_unique_and_parents_precede_children(self):
        batch = {'mounts': [{'source': 'data-bridget', 'target': '/work/data-kate'},
                            {'source': 'tools/revkit/work/stage20/vector-pitch-overlay/pitch',
                             'target': '/work/data-kate/M16/ttsdata/tree2/pitch'},
                            {'source': 'tools/revkit/work/stage19/tree2-suffix',
                             'target': '/work/data-kate/M16/ttsdata/tree2'}]}
        with patch.object(Path, 'exists', return_value=True):
            command = runtime.docker_command(Path('/tmp/scratch/work'), batch)
        mounts = [command[index + 1] for index, value in enumerate(command) if value == '-v']
        targets = [value.split(':')[1] for value in mounts]
        self.assertEqual(len(targets), len(set(targets)))
        self.assertLess(targets.index('/work/data-kate/M16/ttsdata/tree2'), targets.index('/work/data-kate/M16/ttsdata/tree2/pitch'))
        self.assertIn('--read-only', command)
        self.assertIn('no-new-privileges', command)
        self.assertEqual(command[command.index('--network') + 1], 'none')

    def test_scratch_preparation_does_not_copy_wavs_and_preserves_edits(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory) / 'repository'
            source = 'tools/revkit/work/example/input.txt'
            (root / source).parent.mkdir(parents=True)
            (root / source).write_bytes(b'fixture\n')
            scratch = Path(directory) / 'scratch'
            batch = {'sources': [source], 'outputs': {'tools/revkit/work/example/output.wav': {}}}
            with patch.object(runtime, 'ROOT', root):
                work = runtime.prepare(scratch, batch, {'gdb': {}})
            self.assertEqual((work / 'example/input.txt').read_bytes(), b'fixture\n')
            self.assertFalse((work / 'example/output.wav').exists())
            self.assertEqual((work / 'stage5/output.wav').read_bytes(), b'')
            with self.assertRaises(ValueError):
                runtime.relative_path('../outside')
            with self.assertRaises(ValueError):
                runtime.work_relative('data-kate/sample.wav')

    def test_preparation_creates_parent_for_aliased_generated_output(self):
        with tempfile.TemporaryDirectory() as directory:
            scratch = Path(directory) / 'scratch'
            generated = 'tools/revkit/work/different/folder/native.wav'
            batch = {'outputs': {'tools/revkit/work/catalogued.wav': {
                'generated_path': generated}}}
            work = runtime.prepare(scratch, batch, {'gdb': {}})
            self.assertTrue((work / runtime.work_relative(generated)).parent.is_dir())
            self.assertFalse((work / runtime.work_relative(generated)).exists())


if __name__ == '__main__':
    unittest.main()
