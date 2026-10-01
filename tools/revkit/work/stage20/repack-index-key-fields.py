#!/usr/bin/env python3
"""Build disposable 2013 index candidates from the native 2006 field layout.

The candidate mappings are hypotheses. They preserve the verified 12N metric
tail and vary only how the old five-byte key and adjacent byte fields populate
the seven-byte 2013 signature.
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
WORK = Path(__file__).resolve().parent
SOURCE = ROOT / "tools" / "revkit" / "work" / "stage19" / "data-kate-copy" / "M16" / "mc_idx_tbl"
OUTPUT = WORK
ENGINE = ROOT / "binary" / "vt_pau.dll"
IMAGE_BASE = 0x10000000
TABLE_ADDRESSES = (0x1007B7EC, 0x1007B788, 0x1007B850)


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--variant",
        help="build only the named disposable repack variant",
    )
    args = parser.parse_args()
    sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))
    from inspect_legacy_unit_idx import parse_index_data
    from inspect_unit_idx import inspect as inspect_versioned

    engine = ENGINE.read_bytes()
    class_tables = tuple(
        engine[address - IMAGE_BASE : address - IMAGE_BASE + 256]
        for address in TABLE_ADDRESSES
    )
    if any(len(table) != 256 for table in class_tables):
        raise ValueError("class-key lookup table extends beyond the local DLL")

    mappings = {
        "key5-b-c": (("key5", "attr48", "attr40"), None, False),
        "attr48-key5-attr40": (("attr48", "key5", "attr40"), None, False),
        "attr40-key5-attr48": (("attr40", "key5", "attr48"), None, False),
        "key5-b-c-copy-attr40": (("key5", "attr48", "attr40"), "attr40", False),
        "attr48-key3-attr40-key2": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-context-tail-map": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-context-prefix-tail-map": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-apple-context-prefix-tail-map": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-matched-context-tail-map": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-tail-map": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-0404-to-5d00": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-0404-to-5500": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-0404-to-4d00": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-0404-to-6500": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-exact-first-pool-filter": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-hi-first-candidate-prefix-filter": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-attrb40-hi-exact-first-pool-filter": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            "attr40",
            True,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slot8-tail-map": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slot8-0404-to-4500": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slot0-pool-to-a000": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slot0-pool-to-a000-v2": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slot5-pool-to-0000": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            None,
            False,
        ),
        "attr48-key3-attr40-key2-attrb40-transfer": (
            ("attr48", "key5[:3]", "attr40", "key5[3:]"),
            "attr40",
            True,
        ),
    }
    if args.variant is not None and args.variant not in mappings:
        parser.error(f"unknown variant {args.variant!r}")
    for name, (ordering, attr_b_source, transfer_attr40) in mappings.items():
        if args.variant is not None and name != args.variant:
            continue
        context_tail_map = (
            {
                bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                bytes((0x00, 0x1E)): bytes((0x20, 0x00)),
                bytes((0x00, 0x00)): bytes((0x00, 0x00)),
                bytes((0x04, 0x04)): bytes((0x45, 0x00)),
            }
            if name.endswith(("context-tail-map", "context-prefix-tail-map"))
            and not name.endswith(("matched-context-tail-map", "hi-tail-map"))
            else {}
        )
        context_prefix_tail_map = (
            {
                bytes((0x22, 0x17, 0x2B)): {
                    bytes((0x00, 0x1E)): bytes((0x20, 0x00)),
                    bytes((0x00, 0x00)): bytes((0x20, 0x00)),
                }
            }
            if name.endswith("context-prefix-tail-map")
            else {}
        )
        if name.endswith("apple-context-prefix-tail-map"):
            context_prefix_tail_map[bytes((0x47, 0x07, 0x2B))] = {
                bytes((0x00, 0x04)): bytes((0x05, 0x00)),
            }
        matched_context_tail_map = (
            {
                bytes((0x5A, 0x22, 0x17)): {
                    bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                },
                bytes((0x22, 0x17, 0x2B)): {
                    bytes((0x00, 0x1E)): bytes((0x20, 0x00)),
                    bytes((0x00, 0x00)): bytes((0x20, 0x00)),
                },
                bytes((0x2B, 0x30, 0x5A)): {
                    bytes((0x04, 0x04)): bytes((0x45, 0x00)),
                },
            }
            if name.endswith(("matched-context-tail-map", "hi-tail-map"))
            else {}
        )
        if name.endswith("repeat-hello-slot8-tail-map"):
            matched_context_tail_map = {
                bytes((0x2B, 0x30, 0x5A)): {
                    bytes((0x04, 0x00)): bytes((0x45, 0x00)),
                }
            }
        if name.endswith("repeat-hello-slot8-0404-to-4500"):
            matched_context_tail_map = {
                bytes((0x2B, 0x30, 0x5A)): {
                    bytes((0x04, 0x04)): bytes((0x45, 0x00)),
                }
            }
        candidate_prefix_tail_map = {}
        if name.endswith("repeat-hello-slot0-pool-to-a000"):
            candidate_prefix_tail_map = {
                (bytes((0x5A, 0x22, 0x04)), bytes((0x1E, 0x1E))): bytes(
                    (0x5A, 0x22, 0x17)
                )
            }
        if name.endswith("repeat-hello-slot0-pool-to-a000-v2"):
            candidate_prefix_tail_map = {
                (bytes((0x5A, 0x22, 0x05)), bytes((0x1E, 0x1E))): bytes(
                    (0x5A, 0x22, 0x17)
                ),
                (bytes((0x5B, 0x22, 0x05)), bytes((0x1E, 0x1E))): bytes(
                    (0x5A, 0x22, 0x17)
                ),
            }
        if name.endswith("repeat-hello-slot5-pool-to-0000"):
            target_prefix = bytes((0x2F, 0x22, 0x1E))
            candidate_prefix_tail_map = {
                (bytes.fromhex(prefix), bytes.fromhex(tail)): target_prefix
                for prefix, tail in (
                    ("112205", "0a0a"),
                    ("1a2217", "1414"),
                    ("1e2205", "0a0a"),
                    ("1e2217", "1414"),
                    ("272205", "0a0a"),
                    ("2f2205", "0a0a"),
                    ("362205", "0a0a"),
                    ("3e2205", "0a0a"),
                    ("3e2218", "0a0a"),
                    ("3f2205", "0a0a"),
                )
            }
        if name.endswith(
            ("repeat-hello-slots0-5-8-pools", "repeat-hello-slots0-3-5-8-pools")
        ):
            candidate_prefix_tail_map = {
                (bytes((0x5A, 0x22, 0x05)), bytes((0x1E, 0x1E))): bytes(
                    (0x5A, 0x22, 0x17)
                ),
                (bytes((0x5B, 0x22, 0x05)), bytes((0x1E, 0x1E))): bytes(
                    (0x5A, 0x22, 0x17)
                ),
                **{
                    (bytes.fromhex(prefix), bytes.fromhex(tail)): bytes(
                        (0x2F, 0x22, 0x1E)
                    )
                    for prefix, tail in (
                        ("112205", "0a0a"),
                        ("1a2217", "1414"),
                        ("1e2205", "0a0a"),
                        ("1e2217", "1414"),
                        ("272205", "0a0a"),
                        ("2f2205", "0a0a"),
                        ("362205", "0a0a"),
                        ("3e2205", "0a0a"),
                        ("3e2218", "0a0a"),
                        ("3f2205", "0a0a"),
                    )
                },
            }
        if candidate_prefix_tail_map:
            matched_context_tail_map = {
                bytes((0x5A, 0x22, 0x17)): {
                    bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                }
            }
        if name.endswith("repeat-hello-slot5-pool-to-0000"):
            matched_context_tail_map = {
                bytes((0x0D, 0x22, 0x17)): {
                    bytes((0x0A, 0x0A)): bytes((0x00, 0x00)),
                    bytes((0x14, 0x14)): bytes((0x00, 0x00)),
                }
            }
        if name.endswith(
            ("repeat-hello-slots0-5-8-pools", "repeat-hello-slots0-3-5-8-pools")
        ):
            matched_context_tail_map = {
                bytes((0x5A, 0x22, 0x17)): {
                    bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                },
                bytes((0x0D, 0x22, 0x17)): {
                    bytes((0x0A, 0x0A)): bytes((0x00, 0x00)),
                    bytes((0x14, 0x14)): bytes((0x00, 0x00)),
                },
                bytes((0x2B, 0x30, 0x5A)): {
                    bytes((0x04, 0x04)): bytes((0x45, 0x00)),
                },
            }
            if name.endswith("repeat-hello-slots0-3-5-8-pools"):
                matched_context_tail_map[bytes((0x2B, 0x30, 0x22))] = {
                    bytes((0x02, 0x02)): bytes((0x00, 0x00)),
                }
        if name.endswith("hi-tail-map"):
            matched_context_tail_map = {
                bytes((0x5A, 0x22, 0x01)): {
                    bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                },
                bytes((0x22, 0x11, 0x5A)): {
                    bytes((0x04, 0x22)): bytes((0x65, 0x00)),
                },
            }
        if name.endswith("hi-exact-first-pool-filter"):
            matched_context_tail_map = {
                bytes((0x5A, 0x22, 0x01)): {
                    bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                },
                bytes((0x22, 0x11, 0x5A)): {
                    bytes((0x04, 0x22)): bytes((0x65, 0x00)),
                    bytes((0x04, 0x04)): bytes((0x65, 0x00)),
                },
            }
            candidate_prefix_tail_map = {
                (bytes.fromhex(prefix), bytes((0x1E, 0x1E))): bytes.fromhex("5a2213")
                for prefix in (
                    "5a2202",
                    "5a2203",
                    "5a220e",
                    "5a220f",
                    "5b2202",
                    "5b220e",
                    "5b2212",
                )
            }
        if name.endswith("hi-first-candidate-prefix-filter"):
            candidate_prefix_tail_map = {
                (bytes.fromhex(prefix), bytes((0x1E, 0x1E))): bytes.fromhex("5a2213")
                for prefix in (
                    "5a2202",
                    "5a2203",
                    "5a220e",
                    "5a220f",
                    "5b2202",
                    "5b220e",
                    "5b2212",
                )
            }
        relaxation_targets = {
            "hi-0404-to-5d00": bytes((0x5D, 0x00)),
            "hi-0404-to-5500": bytes((0x55, 0x00)),
            "hi-0404-to-4d00": bytes((0x4D, 0x00)),
            "hi-0404-to-6500": bytes((0x65, 0x00)),
        }
        for suffix, target in relaxation_targets.items():
            if name.endswith(suffix):
                matched_context_tail_map = {
                    bytes((0x5A, 0x22, 0x01)): {
                        bytes((0x1E, 0x1E)): bytes((0xA0, 0x00)),
                    },
                    bytes((0x22, 0x11, 0x5A)): {
                        bytes((0x04, 0x22)): bytes((0x65, 0x00)),
                        bytes((0x04, 0x04)): target,
                    },
                }
        destination = OUTPUT / f"index-adapter-key-repacked-{name}"
        destination.mkdir(parents=True, exist_ok=True)
        manifest = [
            "source=Kate ver.2005 unit indexes",
            "legacy_read_order=attr_4c[N],attr_48[N],key[5N],attr_40[N],metrics[12N]",
            f"signature7_order={','.join(ordering)}",
            f"attr_b={attr_b_source or 'zero'}",
            f"transfer_attr_40={str(transfer_attr40).lower()}",
        ]
        if context_tail_map:
            manifest.append(
                "context_tail_map=1e1e:a000,001e:2000,0000:0000,0404:4500"
            )
        if context_prefix_tail_map:
            manifest.append("prefix_tail_map=22172b:001e:2000,0000:2000")
        if name.endswith("apple-context-prefix-tail-map"):
            manifest.append("apple_prefix_tail_map=47072b:0004:0500")
        if matched_context_tail_map:
            manifest.append(
                "matched_context_tail_map="
                + ",".join(
                    f"{prefix.hex()}:{source.hex()}:{target.hex()}"
                    for prefix, mappings_for_prefix in matched_context_tail_map.items()
                    for source, target in mappings_for_prefix.items()
                )
            )
        if candidate_prefix_tail_map:
            manifest.append(
                "candidate_prefix_tail_map="
                + ",".join(
                    f"{prefix.hex()}:{tail.hex()}:{target.hex()}"
                    for (prefix, tail), target in candidate_prefix_tail_map.items()
                )
            )
        for source in sorted(SOURCE.glob("unit-*.idx")):
            raw = source.read_bytes()
            parsed = parse_index_data(source, raw)
            units = parsed.unit_count
            block_end = parsed.block_start + units * parsed.record_stride
            tail = raw[block_end:]
            if len(tail) != units * 20:
                raise ValueError(f"{source.name}: expected 20*N legacy feature bytes")
            attr_4c = tail[:units]
            attr_48 = tail[units : 2 * units]
            key5 = tail[2 * units : 7 * units]
            attr_40 = tail[7 * units : 8 * units]
            metrics = tail[8 * units :]
            signature = bytearray()
            for unit in range(units):
                key = key5[unit * 5 : (unit + 1) * 5]
                mapped_prefix = candidate_prefix_tail_map.get(
                    (key[:3], key[3:5]), key[:3]
                )
                mapped_key = mapped_prefix + key[3:]
                fields = {
                    "key5": mapped_key,
                    "key5[:3]": mapped_prefix,
                    "key5[3:]": context_prefix_tail_map.get(mapped_prefix, {}).get(
                        mapped_key[3:5],
                        matched_context_tail_map.get(
                            bytes((
                                class_tables[0][mapped_prefix[0]],
                                class_tables[1][mapped_prefix[1]],
                                class_tables[2][mapped_prefix[2]],
                            )), {}
                        ).get(
                            mapped_key[3:5], context_tail_map.get(mapped_key[3:5], mapped_key[3:])
                        ),
                    ),
                    "attr48": bytes((attr_48[unit],)),
                    "attr40": bytes((0 if transfer_attr40 else attr_40[unit],)),
                }
                signature.extend(b"".join(fields[field] for field in ordering))

            converted_tail = (
                attr_4c
                + bytes(signature)
                + (attr_40 if attr_b_source == "attr40" else bytes(units))
                + metrics
            )
            header_end = 1 + raw[0]
            header = raw[1:header_end]
            if not header.startswith(b"ver.2005\0VoiceText-Eng\0"):
                raise ValueError(f"{source.name}: unexpected version header")
            converted_header = header.replace(b"ver.2005", b"ver.2013", 1)
            converted = (
                raw[:1]
                + converted_header
                + raw[header_end:block_end]
                + converted_tail
            )
            target = destination / source.name
            target.write_bytes(converted)
            inspected = inspect_versioned(target)
            if inspected["unit_count"] != units:
                raise ValueError(f"{source.name}: unit count changed")
            manifest.append(f"{source.name}: units={units} bytes={len(converted)}")
        (destination / "manifest.txt").write_text("\n".join(manifest) + "\n")
        (WORK / f"compose-key-repacked-{name}.yaml").write_text(
            "services:\n"
            "  runtime:\n"
            "    volumes:\n"
            f"      - ../stage20/index-adapter-key-repacked-{name}:/work/data-kate/M16/mc_idx_tbl:ro\n"
        )
        print("\n".join(manifest))


if __name__ == "__main__":
    main()
