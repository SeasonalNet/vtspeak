#!/usr/bin/env python3
"""Decode and validate the compressed keys in VoiceText's TPP dictionary."""

from __future__ import annotations

import argparse
import json
import re
from collections import Counter
from pathlib import Path


CHAR_MAP_OFFSET = 0x7E388
CHAR_MAP_STRIDE = 2
PAIR_TABLE_OFFSET = 0x81568
PAIR_TABLE_STRIDE = 3
PAIR_COUNT = 127
ATOM = re.compile(rb"[A-G](?:X|[0-9]+)\Z")


def parse_records(data: bytes) -> list[tuple[bytes, bytes]]:
    records: list[tuple[bytes, bytes]] = []
    offset = 0
    while offset < len(data):
        key_end = data.find(b"\0", offset)
        if key_end < 0:
            raise ValueError(f"missing key terminator at byte {offset}")
        value_start = key_end + 1
        value_end = data.find(b"\0", value_start)
        if value_end < 0:
            raise ValueError(f"missing value terminator at byte {value_start}")
        records.append((data[offset:key_end], data[value_start:value_end]))
        offset = value_end + 1
    return records


def decode_key(
    key: bytes, char_map: bytes, pairs: list[bytes], inverse: dict[int, int]
) -> bytes:
    decoded = bytearray()
    for value in key:
        if value & 0x80:
            pair_index = (value & 0x7F) - 1
            if not 0 <= pair_index < len(pairs):
                raise ValueError(f"invalid pair token 0x{value:02x}")
            transformed = pairs[pair_index]
        else:
            transformed = bytes((value,))
        for character in transformed:
            try:
                decoded.append(inverse[character])
            except KeyError as exc:
                raise ValueError(f"no printable source byte for 0x{character:02x}") from exc
    return bytes(decoded)


def encode_key(
    key: bytes, char_map: bytes, pair_indexes: dict[bytes, int]
) -> bytes:
    transformed = bytes(char_map[value] for value in key)
    encoded = bytearray()
    offset = 0
    while offset < len(transformed):
        pair = transformed[offset : offset + 2]
        pair_index = pair_indexes.get(pair) if len(pair) == 2 else None
        if pair_index is None:
            encoded.append(transformed[offset])
            offset += 1
        else:
            encoded.append((pair_index + 1) | 0x80)
            offset += 2
    return bytes(encoded)


def inspect(dictionary_path: Path, dll_path: Path, show_rows: bool) -> dict[str, object]:
    dll = dll_path.read_bytes()
    if dll[:2] != b"MZ":
        raise ValueError(f"{dll_path}: missing MZ executable header")
    char_map = bytes(
        dll[CHAR_MAP_OFFSET + CHAR_MAP_STRIDE * index]
        for index in range(256)
    )
    pairs = [
        dll[PAIR_TABLE_OFFSET + PAIR_TABLE_STRIDE * index :
            PAIR_TABLE_OFFSET + PAIR_TABLE_STRIDE * index + 2]
        for index in range(PAIR_COUNT)
    ]
    if len(set(pairs)) != PAIR_COUNT or pairs != sorted(pairs):
        raise ValueError("the recovered two-byte pair table is not unique and sorted")
    inverse: dict[int, int] = {}
    for source in range(0x20, 0x7F):
        inverse.setdefault(char_map[source], source)
    pair_indexes = {pair: index for index, pair in enumerate(pairs)}

    records = parse_records(dictionary_path.read_bytes())
    sequence_counts: Counter[str] = Counter()
    atom_counts: Counter[str] = Counter()
    component_counts: Counter[str] = Counter()
    compound_hyphen_counts: Counter[int] = Counter()
    tag_values: dict[str, set[str]] = {}
    decoded_rows: list[dict[str, str]] = []
    max_atoms = 0
    for key, value in records:
        decoded_key = decode_key(key, char_map, pairs, inverse)
        if encode_key(decoded_key, char_map, pair_indexes) != key:
            raise ValueError(f"key did not round-trip: {key.hex()} -> {decoded_key!r}")
        atoms = value.split(b" ")
        if not atoms or any(not ATOM.fullmatch(atom) for atom in atoms):
            raise ValueError(f"unexpected TPP payload grammar for key {decoded_key!r}: {value!r}")
        key_text = decoded_key.decode("ascii")
        primary = atoms[0]
        primary_tag = primary[:1].decode("ascii")
        hyphen_count = key_text.count("-")
        if primary_tag == "A":
            if hyphen_count != 0:
                raise ValueError(f"single-component A key contains a hyphen: {key_text!r}")
            component_counts[primary_tag] += 1
        if primary_tag in "BCDE":
            encoded_count = int(primary[1:2])
            flags = primary[2:]
            if (
                encoded_count != ord(primary_tag) - ord("A") + 1
                or len(flags) != encoded_count
                or any(flag not in b"01" for flag in flags)
                or hyphen_count != encoded_count - 1
            ):
                raise ValueError(f"inconsistent place-name component code: {key_text!r} -> {value!r}")
            component_counts[primary_tag] += 1
        if primary_tag == "F":
            if not 1 <= hyphen_count <= 4:
                raise ValueError(f"F compound must contain one to four hyphens: {key_text!r}")
            compound_hyphen_counts[hyphen_count] += 1
        if primary_tag == "G" and hyphen_count != 0:
            raise ValueError(f"G general lexical key unexpectedly contains a hyphen: {key_text!r}")
        atom_tags = " ".join(atom[:1].decode("ascii") for atom in atoms)
        sequence_counts[atom_tags] += 1
        max_atoms = max(max_atoms, len(atoms))
        for atom in atoms:
            tag = atom[:1].decode("ascii")
            atom_counts[tag] += 1
            tag_values.setdefault(tag, set()).add(atom[1:].decode("ascii"))
        decoded_rows.append(
            {
                "key": key_text,
                "payload": value.decode("ascii"),
            }
        )

    report: dict[str, object] = {
        "dictionary": str(dictionary_path),
        "records": len(records),
        "key_transform": "low byte of a 256-entry u16 character map, then greedy sorted two-byte replacement",
        "character_map_offset": f"0x{CHAR_MAP_OFFSET:x}",
        "pair_table_offset": f"0x{PAIR_TABLE_OFFSET:x}",
        "pair_table_entries": len(pairs),
        "decoded_key_round_trips": len(records),
        "max_atoms_per_payload": max_atoms,
        "payload_sequence_counts": dict(sorted(sequence_counts.items())),
        "payload_atom_counts": dict(sorted(atom_counts.items())),
        "place_name_records_by_component_count": dict(sorted(component_counts.items())),
        "F_compound_records_by_hyphen_count": {
            str(count): records for count, records in sorted(compound_hyphen_counts.items())
        },
        "payload_suffix_ranges": {
            tag: [min(map(int, values)), max(map(int, values))]
            for tag, values in sorted(tag_values.items())
            if all(value.isdigit() for value in values)
        },
        "sample_rows": decoded_rows[:12],
    }
    if show_rows:
        report["rows"] = decoded_rows
    return report


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--dictionary",
        type=Path,
        default=Path("data-common/dict-eng/tppdict_eng"),
        help="TPP dictionary data file",
    )
    parser.add_argument(
        "--dll",
        type=Path,
        default=Path("binary/vt_pau.dll"),
        help="matching VoiceText DLL containing the character and pair tables",
    )
    parser.add_argument("--show-rows", action="store_true", help="include every decoded row")
    args = parser.parse_args()
    print(json.dumps(inspect(args.dictionary, args.dll, args.show_rows), indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
