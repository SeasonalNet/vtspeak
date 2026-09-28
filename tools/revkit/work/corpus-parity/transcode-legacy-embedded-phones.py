#!/usr/bin/env python3
"""Build a disposable MSI dictionary variant with phone IDs remapped to the current codebook."""

from __future__ import annotations

import json
import re
import shutil
import struct
from pathlib import Path


ROOT = Path("tools/revkit/work/corpus-parity")
SOURCE = ROOT / "variant-old-embedded"
TARGET = ROOT / "variant-old-embedded-idremap-v4"
RUNTIME_LOG = ROOT / "runtime-original-msi-2006-dict.log"
CODEBOOK = Path("docs/reverse-engineering/phone-symbol-codebook.md")


def phone_tables() -> tuple[list[bytes], list[bytes | None]]:
    log = RUNTIME_LOG.read_text(encoding="utf-8")
    segment = log.split("LEGACY_PHONE_ID_TABLE_256X5_BYTES\n", 1)[1].split(
        "EMBEDDED_PAYLOAD_DECODE_ENTRY", 1
    )[0]
    legacy_bytes: list[int] = []
    for line in segment.splitlines():
        if ":" in line:
            legacy_bytes.extend(
                int(value, 16)
                for value in re.findall(r"0x([0-9a-fA-F]{2})", line.split(":", 1)[1])
            )
    if len(legacy_bytes) != 256 * 5:
        raise ValueError(f"expected 1280 bytes from legacy table, got {len(legacy_bytes)}")
    legacy = [
        bytes(legacy_bytes[index * 5 : index * 5 + 5]).split(b"\0", 1)[0]
        for index in range(256)
    ]

    text = CODEBOOK.read_text(encoding="utf-8")
    section = text.split("| ID | Bytes | ID | Bytes | ID | Bytes | ID | Bytes |", 1)[1]
    current: list[bytes | None] = [None] * 256
    for line in section.splitlines():
        for match in re.finditer(r"`0x([0-9a-fA-F]{2})`\s*\|\s*`([^`]*)`", line):
            raw = match.group(2).strip()
            current[int(match.group(1), 16)] = (
                b"" if raw in {"—", ""} else bytes.fromhex(raw)
            )
    if any(value is None for value in current):
        raise ValueError("current codebook does not describe all 256 IDs")
    return legacy, current


def indexed_records(dictionary: bytes, index: bytes) -> list[tuple[int, int, bytes, bytes]]:
    if len(index) < 4 or (len(index) - 4) % 4:
        raise ValueError("unexpected hashidx_emb size")
    count = (len(index) - 4) // 4
    offsets = sorted(struct.unpack_from(f"<{count}I", index, 4))
    records = []
    for record_number, start in enumerate(offsets):
        end = offsets[record_number + 1] if record_number + 1 < count else len(dictionary)
        key_end = dictionary.find(b"\0", start, end)
        if key_end < 0:
            raise ValueError(f"missing key terminator at dictionary offset {start}")
        payload_start = key_end + 1
        payload_end = dictionary.find(b"\0", payload_start, end)
        if payload_end != end - 1:
            raise ValueError(f"unexpected payload extent at dictionary offset {start}")
        records.append((start, payload_start, dictionary[start:key_end], dictionary[payload_start:payload_end]))
    return records


def main() -> None:
    if TARGET.exists():
        raise FileExistsError(f"refusing to overwrite existing experiment directory: {TARGET}")
    shutil.copytree(SOURCE, TARGET)
    legacy, current = phone_tables()
    reverse: dict[bytes, list[int]] = {}
    for new_id, symbols in enumerate(current):
        assert symbols is not None
        reverse.setdefault(symbols, []).append(new_id)
    mapping = {}
    crosswalk = ["old_id\told_symbols_hex\tcurrent_id\tcurrent_symbols_hex\tstatus"]
    for old_id, symbols in enumerate(legacy):
        ids = reverse.get(symbols, [])
        if len(ids) == 1:
            mapping[old_id] = ids[0]
            status = "same-id" if ids[0] == old_id else "remapped"
            current_id = str(ids[0])
            current_hex = current[ids[0]].hex()
        elif ids:
            status = "ambiguous-current-equivalent"
            current_id = ",".join(str(value) for value in ids)
            current_hex = "|".join(current[value].hex() for value in ids)
        else:
            status = "no-current-equivalent"
            current_id = ""
            current_hex = ""
        crosswalk.append(
            f"{old_id}\t{symbols.hex()}\t{current_id}\t{current_hex}\t{status}"
        )

    dictionary_path = TARGET / "engttsdict_emb"
    index_path = TARGET / "hashidx_emb"
    raw = bytearray(dictionary_path.read_bytes())
    records = indexed_records(bytes(raw), index_path.read_bytes())
    transformed_direct = 0
    transformed_variant = 0
    skipped_other_form = 0
    skipped_malformed_variant = 0
    skipped_unmapped = 0
    changed_payloads = 0
    exact_current = 0
    current_dict = indexed_records(
        Path("data-common/dict-eng/engttsdict_emb").read_bytes(),
        Path("data-common/dict-eng/hashidx_emb").read_bytes(),
    )
    current_payloads = {key: payload for _, _, key, payload in current_dict}
    manifest = ["key_hex\tmsi_payload_hex\tremapped_payload_hex\tcurrent_payload_hex\tresult"]

    for start, payload_start, key, payload in records:
        if not payload:
            skipped_other_form += 1
            continue
        if payload[0] & 1:
            try:
                remapped = bytes(
                    [payload[0], *(mapping[value] for value in payload[1:])]
                )
            except KeyError:
                skipped_unmapped += 1
                continue
            transformed_direct += 1
        elif (payload[0] & 0xFE) == 2:
            chunks = payload[1:].split(b"\xff")
            if any(chunk.count(b"|") != 1 for chunk in chunks):
                skipped_malformed_variant += 1
                continue
            try:
                remapped_chunks = []
                for chunk in chunks:
                    context, phone_ids = chunk.split(b"|", 1)
                    remapped_chunks.append(
                        context
                        + b"|"
                        + bytes(mapping[value] for value in phone_ids)
                    )
            except KeyError:
                skipped_unmapped += 1
                continue
            remapped = bytes([payload[0]]) + b"\xff".join(remapped_chunks)
            transformed_variant += 1
        else:
            skipped_other_form += 1
            continue
        raw[payload_start : payload_start + len(payload)] = remapped
        current_payload = current_payloads.get(key)
        if current_payload != payload:
            changed_payloads += 1
        if current_payload == remapped:
            exact_current += 1
            result = "matches-current"
        elif current_payload is None:
            result = "current-key-missing"
        else:
            result = "differs-after-remap"
        manifest.append(
            f"{key.hex()}\t{payload.hex()}\t{remapped.hex()}\t"
            f"{current_payload.hex() if current_payload is not None else ''}\t{result}"
        )

    dictionary_path.write_bytes(raw)
    (TARGET / "phone-id-remap.tsv").write_text("\n".join(manifest) + "\n", encoding="utf-8")
    (TARGET / "phone-id-crosswalk.tsv").write_text(
        "\n".join(crosswalk) + "\n", encoding="utf-8"
    )
    report = {
        "indexed_msi_records": len(records),
        "old_to_current_id_mappings": len(mapping),
        "same_numeric_id_symbol_vectors": sum(
            legacy[index] == current[index] for index in range(256)
        ),
        "different_numeric_id_symbol_vectors": sum(
            legacy[index] != current[index] for index in range(256)
        ),
        "ambiguous_old_ids": sum(
            len(reverse.get(symbols, [])) > 1 for symbols in legacy
        ),
        "old_ids_without_current_equivalent": sum(
            not reverse.get(symbols) for symbols in legacy
        ),
        "transformed_direct_payloads": transformed_direct,
        "transformed_alternative_payloads": transformed_variant,
        "transformed_payloads_different_from_current": changed_payloads,
        "transformed_payloads_matching_current": exact_current,
        "skipped_other_payload_forms": skipped_other_form,
        "skipped_malformed_alternative_payloads": skipped_malformed_variant,
        "skipped_unmapped_payloads": skipped_unmapped,
        "dictionary_bytes_preserved": len(raw),
        "hash_index_and_parameters_unchanged": True,
        "note": "Remapped direct phone IDs and alternative-branch phone IDs; context bytes were preserved.",
    }
    (TARGET / "phone-id-remap.json").write_text(
        json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    print(json.dumps(report, indent=2, sort_keys=True))


if __name__ == "__main__":
    main()
