import hashlib
import importlib.util
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch


SPEC = importlib.util.spec_from_file_location(
    "reproduce_artifacts", Path(__file__).resolve().parents[1] / "reproduce_artifacts.py"
)
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)


def entry(data, **recipe):
    return {"bytes": len(data), "sha256": hashlib.sha256(data).hexdigest(), **recipe}


class ReproductionTests(unittest.TestCase):
    def test_decoded_segments_check_source_hash_and_sample_bounds(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            source = root / "data-kate/example.dat"
            source.parent.mkdir()
            source.write_bytes(b"prefixDATAtrailer")
            segment = {"dat": "data-kate/example.dat", "offset": 6, "bytes": 4,
                       "sha256": hashlib.sha256(b"DATA").hexdigest(), "start": 1, "count": 2}
            with patch("decode_dat.decode_payload", return_value=[-1, 123, -234]):
                self.assertEqual(MODULE.decode_segments(root, [segment]), b"{\x00\x16\xff")
                with self.assertRaisesRegex(ValueError, "out of bounds"):
                    MODULE.decode_segments(root, [{**segment, "count": 3}])
            source.write_bytes(b"prefixEDITtrailer")
            with self.assertRaisesRegex(ValueError, "source differs"):
                MODULE.decode_segments(root, [segment])

    def test_vendor_reads_reject_traversal_and_external_links(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            for source in ("../data-kate/sample.dat", "binary/sample.dat", "/etc/passwd"):
                with self.assertRaises(ValueError):
                    MODULE.read_input(root, source, 0, 0, "")
            (root / "data-kate").symlink_to("/tmp")
            with self.assertRaises(ValueError):
                MODULE.read_input(root, "data-kate/sample.dat", 0, 0, "")

    def test_byte_write_expansion_preserves_offsets_and_plain_lines(self):
        macro = '#@bytes ' + json.dumps(["  set {char}($buffer + ", 7, ") = ", "410aff00"])
        actual = MODULE.expand_source(("set pagination off\n" + macro + "\ncontinue\n").encode())
        self.assertEqual(actual, b"set pagination off\n  set {char}($buffer + 7) = 65\n  set {char}($buffer + 8) = 10\n  set {char}($buffer + 9) = 255\n  set {char}($buffer + 10) = 0\ncontinue\n")
        literal = '#@literal ' + json.dumps("  \n\n") + "\n"
        self.assertEqual(MODULE.expand_source(literal.encode()), b"  \n\n")

    def test_audio_copy_rebuilds_dependency_from_pcm_with_exact_header(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            folder = root / "tools/revkit/work/example"
            folder.mkdir(parents=True)
            (folder / "samples.pcm").write_bytes(b"\x01\x02\x03\x04")
            original = b"header" + b"\x01\x02\x03\x04" + b"trailer"
            path = "tools/revkit/work/example/output.wav"
            copy = "tools/revkit/work/example/copy.wav"
            manifest = {"wav": {
                path: entry(original, pcm=["tools/revkit/work/example/samples.pcm"], prefix=b"header".hex(), suffix=b"trailer".hex()),
                copy: entry(original, copy=path),
            }}
            self.assertEqual(MODULE.build(root, manifest, "wav", copy), original)
            (folder / "samples.pcm").write_bytes(b"bad")
            with self.assertRaisesRegex(ValueError, "differs"):
                MODULE.build(root, manifest, "wav", copy)

    def test_python_and_shell_generators_return_verified_command_bytes(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            folder = root / "tools/revkit/work/example"
            folder.mkdir(parents=True)
            (folder / "generate.py").write_text(
                "import sys\nfrom pathlib import Path\nPath(sys.argv[2]).write_bytes(b'continue\\n')\n"
            )
            (folder / "generate.sh").write_text("printf 'continue\\n'\n")
            path = "tools/revkit/work/example/trace.gdb"
            for recipe in (
                {"python_generator": "tools/revkit/work/example/generate.py"},
                {"generator": "tools/revkit/work/example/generate.sh", "args": [], "env": {}},
            ):
                manifest = {"gdb": {path: entry(b"continue\n", **recipe)}}
                self.assertEqual(MODULE.build(root, manifest, "gdb", path), b"continue\n")

    def test_restore_is_idempotent_and_preserves_changed_files(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            path = "tools/revkit/work/example/trace.gdb"
            MODULE.restore(root, path, b"continue\n")
            MODULE.restore(root, path, b"continue\n")
            target = root / path
            target.write_bytes(b"user edit\n")
            with self.assertRaisesRegex(ValueError, "overwrite"):
                MODULE.restore(root, path, b"continue\n")
            self.assertEqual(target.read_bytes(), b"user edit\n")

    def test_paths_reject_vendor_inputs_traversal_and_links(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            for path in ("binary/sample.wav", "tools/revkit/work/../../../binary/sample.wav", "/tmp/output.gdb"):
                with self.assertRaises(ValueError):
                    MODULE.safe_path(root, path)
            folder = root / "tools/revkit/work"
            folder.mkdir(parents=True)
            (folder / "link").symlink_to("/tmp")
            with self.assertRaises(ValueError):
                MODULE.safe_path(root, "tools/revkit/work/link/trace.gdb")

    def test_hash_mismatch_and_dependency_cycle_fail(self):
        with self.assertRaisesRegex(ValueError, "differs"):
            MODULE.checked(b"changed", entry(b"original"), "trace.gdb")
        path = "tools/revkit/work/example/loop.wav"
        manifest = {"wav": {path: entry(b"", copy=path)}}
        with self.assertRaisesRegex(ValueError, "cyclic"):
            MODULE.build(Path("."), manifest, "wav", path)


if __name__ == "__main__":
    unittest.main()
