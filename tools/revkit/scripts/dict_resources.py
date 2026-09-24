#!/usr/bin/env python3
"""Inspect shared VoiceText English dictionary resources without modifying them."""

from __future__ import annotations

import argparse
import json
import struct
from collections import Counter
from pathlib import Path

from tree3 import Tree, describe, parse_tree, parse_tree_data


def parse_atmt_trees(path: Path) -> list[Tree]:
    raw = path.read_bytes()
    if len(raw) < 4:
        raise ValueError(f"{path}: missing 32-bit subtree count")
    subtree_count = struct.unpack_from("<I", raw, 0)[0]
    if subtree_count != 27:
        raise ValueError(f"{path}: expected 27 subtrees from FUN_10003a10, got {subtree_count}")

    offset = 4
    trees = []
    for index in range(subtree_count):
        label = path.with_name(f"{path.name}[{index:02d}]")
        tree, offset = parse_tree_data(label, raw, offset, require_eof=False)
        trees.append(tree)
    if offset != len(raw):
        raise ValueError(f"{path}: parsed subtrees end at {offset}, file ends at {len(raw)}")
    return trees


def inspect_atmt(path: Path) -> dict[str, object]:
    trees = parse_atmt_trees(path)
    descriptions = [describe(tree) for tree in trees]

    return {
        "file": str(path),
        "container_count": len(trees),
        "container_bytes": 4,
        "parsed_subtrees": len(descriptions),
        "total_nodes": sum(int(tree["nodes"]) for tree in descriptions),
        "subtrees": descriptions,
    }


def inspect_hash_family(root: Path, suffix: str) -> dict[str, object]:
    params_path = root / f"hashparams{suffix}"
    index_path = root / f"hashidx{suffix}"
    content_path = root / f"hashcont{suffix}"
    dictionary_name = "engttsdict_emb" if suffix == "_emb" else "tppdict_eng"
    dictionary_path = root / dictionary_name
    params = params_path.read_bytes()
    index = index_path.read_bytes()
    content = content_path.read_bytes()
    dictionary = dictionary_path.read_bytes()

    if not params.startswith(b"NLP"):
        raise ValueError(f"{params_path}: missing observed NLP marker")
    if len(params) < 36:
        raise ValueError(f"{params_path}: too short for the observed hash parameter header")
    version = params[3]
    fields = struct.unpack_from("<7I", params, 4)
    distribution_count = fields[6]
    distribution_offset = 32
    byte_count_offset = distribution_offset + distribution_count * 4
    if byte_count_offset + 4 > len(params):
        raise ValueError(f"{params_path}: distribution table extends beyond the file")
    byte_count = struct.unpack_from("<I", params, byte_count_offset)[0]
    end_offset = byte_count_offset + 4 + byte_count
    if end_offset != len(params):
        raise ValueError(
            f"{params_path}: parsed parameter sections end at {end_offset}, file ends at {len(params)}"
        )

    if len(index) < 4:
        raise ValueError(f"{index_path}: missing 32-bit header value")
    header_value = struct.unpack_from("<I", index, 0)[0]
    if header_value > 200:
        raise ValueError(f"{index_path}: header value {header_value} exceeds loader limit 200")
    index_bytes = len(index) - 4
    if index_bytes % 4:
        raise ValueError(f"{index_path}: index payload is not a sequence of 32-bit entries")
    if len(content) % 2:
        raise ValueError(f"{content_path}: content payload is not a sequence of 16-bit entries")
    index_count = index_bytes // 4
    content_count = len(content) // 2
    if index_count != content_count:
        raise ValueError(
            f"{suffix}: {index_count} index entries do not match {content_count} content entries"
        )

    content_values = struct.unpack(f"<{content_count}H", content) if content_count else ()
    index_values = struct.unpack_from(f"<{index_count}I", index, 4) if index_count else ()
    ordered_offsets = sorted(index_values)
    key_lengths: list[int] = []
    payload_lengths: list[int] = []
    for record_index, record_offset in enumerate(ordered_offsets):
        record_end = (
            ordered_offsets[record_index + 1]
            if record_index + 1 < len(ordered_offsets)
            else len(dictionary)
        )
        key_end = dictionary.find(b"\0", record_offset, record_end)
        if key_end < 0:
            raise ValueError(f"{dictionary_path}: record at {record_offset} has no key terminator")
        payload_start = key_end + 1
        payload_end = dictionary.find(b"\0", payload_start, record_end)
        if payload_end != record_end - 1:
            raise ValueError(
                f"{dictionary_path}: record at {record_offset} does not end with one "
                "NUL-terminated payload"
            )
        key_lengths.append(key_end - record_offset)
        payload_lengths.append(record_end - payload_start)
    if ordered_offsets and ordered_offsets[0] != 0:
        raise ValueError(f"{dictionary_path}: first indexed record does not begin at offset zero")
    return {
        "family": suffix,
        "dictionary_file": dictionary_path.name,
        "dictionary_bytes": len(dictionary),
        "dictionary_record_count": len(ordered_offsets),
        "dictionary_record_partition_bytes": len(dictionary),
        "dictionary_key_bytes_before_nul": [min(key_lengths), max(key_lengths)],
        "dictionary_payload_bytes_including_nul": [
            min(payload_lengths),
            max(payload_lengths),
        ],
        "parameter_file": params_path.name,
        "parameter_bytes": len(params),
        "marker": "NLP",
        "version": version,
        "parameter_words": [f"0x{value:08x}" for value in fields],
        "distribution_entries": distribution_count,
        "distribution_table_offset": distribution_offset,
        "byte_table_entries": byte_count,
        "parsed_parameter_end": end_offset,
        "index_file": index_path.name,
        "index_bytes": len(index),
        "index_header_value": header_value,
        "index_entries": index_count,
        "index_range": [min(index_values, default=0), max(index_values, default=0)],
        "content_file": content_path.name,
        "content_bytes": len(content),
        "content_entries_u16": content_count,
        "content_range": [min(content_values, default=0), max(content_values, default=0)],
    }


def inspect_exceptdict(path: Path) -> dict[str, object]:
    raw = path.read_bytes()
    if len(raw) < 8:
        raise ValueError(f"{path}: missing two 32-bit range fields")
    first_group, last_group = struct.unpack_from("<II", raw)
    if last_group < first_group or last_group - first_group > 64:
        raise ValueError(f"{path}: unreasonable group range {first_group}..{last_group}")
    offset = 8
    groups = []
    total_entries = 0
    for expected_group in range(first_group, last_group + 1):
        if offset + 8 > len(raw):
            raise ValueError(f"{path}: truncated group header at offset {offset}")
        group, count = struct.unpack_from("<II", raw, offset)
        offset += 8
        if group != expected_group:
            raise ValueError(f"{path}: expected group {expected_group}, found {group}")
        key_lengths = []
        value_lengths = []
        for _ in range(count):
            if offset + 8 > len(raw):
                raise ValueError(f"{path}: truncated key/value lengths at offset {offset}")
            key_length, value_length = struct.unpack_from("<II", raw, offset)
            offset += 8
            end = offset + key_length + value_length
            if end > len(raw):
                raise ValueError(f"{path}: key/value pair at offset {offset} exceeds file extent")
            key_lengths.append(key_length)
            value_lengths.append(value_length)
            offset = end
        total_entries += count
        groups.append(
            {
                "group": group,
                "entries": count,
                "key_bytes": [min(key_lengths), max(key_lengths)] if key_lengths else [0, 0],
                "value_bytes": [min(value_lengths), max(value_lengths)] if value_lengths else [0, 0],
            }
        )
    if offset != len(raw):
        raise ValueError(f"{path}: parsed through offset {offset}, file ends at {len(raw)}")
    return {
        "file": str(path),
        "bytes": len(raw),
        "first_group": first_group,
        "last_group": last_group,
        "groups": groups,
        "total_entries": total_entries,
        "parsed_end": offset,
    }


TXT2_MODES = {
    "abbrc_sort.txt2": "M",
    "abbrh_sort.txt2": "M",
    "abbrt_sort.txt2": "M",
    "chc_sort.txt2": "M",
    "citya_sort.txt2": "M",
    "sbdw_sort.txt2": "M",
    "streeta_sort.txt2": "M",
    "streetf_sort.txt2": "L",
    "wab.txt2": "L",
}


def inspect_txt2(path: Path, include_rows: bool = False) -> dict[str, object]:
    raw = path.read_bytes()
    if len(raw) < 17:
        raise ValueError(f"{path}: shorter than the 17-byte framed-resource overhead")
    marker = raw[:3]
    if marker != raw[-3:]:
        raise ValueError(f"{path}: leading and trailing three-byte markers differ")
    shift = raw[-14]
    body = bytes((value - shift) & 0xFF for value in raw[3:-14])
    count_field = bytes((value - shift) & 0xFF for value in raw[-13:-3])
    count_text, separator, padding = count_field.partition(b"\0")
    if not separator or any(padding) or not count_text.isdigit():
        raise ValueError(f"{path}: decoded 10-byte trailer field is not a padded decimal count")
    declared_rows = int(count_text)
    try:
        decoded_body = body.decode("ascii")
    except UnicodeDecodeError as exc:
        raise ValueError(f"{path}: decoded body is not ASCII after subtracting 0x{shift:02x}") from exc
    rows = decoded_body.splitlines()
    mode = TXT2_MODES.get(path.name)
    if mode is None:
        raise ValueError(f"{path}: no callsite mode recorded for this resource")
    expected_columns = 2 if mode == "M" else 1
    row_columns = [row.count("|") + 1 for row in rows]
    if len(rows) != declared_rows:
        raise ValueError(
            f"{path}: trailer declares {declared_rows} rows but body contains {len(rows)}"
        )
    if any(columns != expected_columns for columns in row_columns):
        raise ValueError(
            f"{path}: mode {mode} expects {expected_columns} column(s), "
            "but at least one row has a different shape"
        )
    trailer = raw[-14:-3]
    separators = Counter(body)
    report = {
        "file": str(path),
        "bytes": len(raw),
        "marker_hex": marker.hex(),
        "body_offset": 3,
        "body_transform": f"subtract 0x{shift:02x} modulo 256",
        "body_bytes": len(body),
        "trailer_bytes": len(trailer),
        "trailer_hex": trailer.hex(),
        "count_field_decoded": count_field.rstrip(b"\0").decode("ascii"),
        "declared_rows": declared_rows,
        "observed_rows": len(rows),
        "callsite_column_mode": mode,
        "columns_per_row": expected_columns,
        "sample_rows": rows[:3],
        "body_pipe_bytes": separators[ord("|")],
        "body_lf_bytes": separators[ord("\n")],
        "body_cr_bytes": separators[ord("\r")],
    }
    if include_rows:
        report["rows"] = rows
    return report


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--root",
        type=Path,
        default=Path("data-common/dict-eng"),
        help="shared dict-eng resource directory",
    )
    parser.add_argument(
        "--show-rows",
        action="store_true",
        help="include every decoded .txt2 row instead of only the first three samples",
    )
    args = parser.parse_args()
    root = args.root

    trees = [parse_tree(root / name) for name in ("engbi.tree3", "poly.tree3", "sbd.tree3")]
    text_paths = sorted(root.glob("*.txt2"))
    if {path.name for path in text_paths} != set(TXT2_MODES):
        raise ValueError(f"{root}: expected exactly these nine mapped .txt2 resources")
    text_tables = [inspect_txt2(path, include_rows=args.show_rows) for path in text_paths]
    report = {
        "tree3": [describe(tree) for tree in trees],
        "atmt": inspect_atmt(root / "atmt.tree3"),
        "hash_families": [
            inspect_hash_family(root, "_emb"),
            inspect_hash_family(root, "_eng_tpp"),
        ],
        "exceptdict": inspect_exceptdict(root / "exceptdict"),
        "txt2_tables": text_tables,
    }
    print(json.dumps(report, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
