"""Synthetic checks for the observed legacy VoiceText parser layouts."""

from __future__ import annotations

import struct
import sys
import tempfile
import unittest
from pathlib import Path

SCRIPT_DIR = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(SCRIPT_DIR))

from scripts.inspect_legacy_unit_idx import parse_index_data
from scripts.tree2 import parse_tree2, parse_tree2_data


def make_node(
    feature: int,
    operation: str,
    threshold: int,
    output: int,
    values: tuple[int, ...] = (),
    trailer: bytes = b"\0\0",
) -> bytes:
    return (
        struct.pack("<BBhhB", feature, ord(operation), threshold, output, len(values))
        + struct.pack(f"<{len(values)}h", *values)
        + trailer
    )


def make_index(version: str = "ver.2005") -> bytes:
    header = f"{version}\0VoiceText-Eng\0".encode("ascii")
    bank_name = b"merged-gen"
    record = bytearray(19)
    record[14:16] = b"\1\1"
    return (
        bytes([len(header)])
        + header
        + struct.pack("<H", 1)
        + struct.pack("<H", len(bank_name))
        + bank_name
        + b"\0"
        + struct.pack("<I", 1)
        + struct.pack("<H", 19)
        + record
        + bytes(20)
    )


class LegacyIndexTests(unittest.TestCase):
    def test_parses_supported_index_headers_and_extents(self) -> None:
        for version in ("ver.2005", "ver.2009"):
            with self.subTest(version=version):
                raw = make_index(version)
                index = parse_index_data(Path("synthetic.idx"), raw)
                self.assertEqual(index.version, version)
                self.assertEqual(index.producer, "VoiceText-Eng")
                self.assertEqual(index.banks[0].name, "merged-gen")
                self.assertEqual(index.record_stride, 19)
                self.assertEqual(index.feature_column_bytes, 20)
                self.assertEqual(index.actual_size, len(raw))

    def test_rejects_unexpected_record_stride(self) -> None:
        raw = bytearray(make_index())
        stride_offset = raw.index(b"merged-gen") + len(b"merged-gen") + 1 + 4
        struct.pack_into("<H", raw, stride_offset, 20)
        with self.assertRaisesRegex(ValueError, "expected 19-byte records"):
            parse_index_data(Path("synthetic.idx"), bytes(raw))

    def test_rejects_truncated_feature_columns(self) -> None:
        with self.assertRaisesRegex(ValueError, "expected .* feature-column bytes"):
            parse_index_data(Path("synthetic.idx"), make_index()[:-1])


class Tree2Tests(unittest.TestCase):
    def setUp(self) -> None:
        self.raw = (
            make_node(0, "D", 0, 0, (1,), b"\1\1")
            + make_node(1, "C", -1, 11)
            + make_node(2, "D", 0, 22)
        )

    def test_parses_recursive_nodes_and_evaluates_both_branches(self) -> None:
        tree = parse_tree2_data(Path("synthetic.tree2"), self.raw)
        self.assertEqual((tree.nodes, tree.leaves, tree.end_offset), (3, 2, len(self.raw)))
        self.assertEqual(tree.evaluate([1, -2, 0]), 11)
        self.assertEqual(tree.evaluate([2, 0, 0]), 22)

    def test_reports_opaque_suffix_and_strict_mode_rejects_it(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "synthetic.tree2"
            path.write_bytes(self.raw + b"opaque")
            tree = parse_tree2(path)
            self.assertEqual(tree.trailing_bytes, 6)
            with self.assertRaisesRegex(ValueError, "trailing bytes"):
                parse_tree2(path, require_eof=True)

    def test_rejects_truncated_child_node(self) -> None:
        with self.assertRaisesRegex(ValueError, "truncated node"):
            parse_tree2_data(Path("synthetic.tree2"), self.raw[:-1])

    def test_rejects_unknown_operation(self) -> None:
        raw = make_node(0, "X", 0, 0)
        with self.assertRaisesRegex(ValueError, "operation byte"):
            parse_tree2_data(Path("synthetic.tree2"), raw)


if __name__ == "__main__":
    unittest.main()
