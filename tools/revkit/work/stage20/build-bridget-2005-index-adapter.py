#!/usr/bin/env python3
"""Build a disposable 2013-reader index overlay from Bridget's 2005 indexes."""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
SOURCE = ROOT / "data-bridget" / "M16" / "mc_idx_tbl"
WORK = Path(__file__).resolve().parent
sys.path.insert(0, str(ROOT / "tools" / "revkit" / "scripts"))

from inspect_legacy_unit_idx import parse_index_data
from inspect_unit_idx import inspect as inspect_versioned


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--layout",
        choices=(
            "structural",
            "signature",
            "selector-key3-attr40",
            "selector-key4-attr40",
        ),
        default="signature",
        help=(
            "choose the experimental seven-byte feature layout; selector layouts "
            "control which legacy key byte appears at +5 before attr_40 at +6"
        ),
    )
    layout = parser.parse_args().layout
    output_root = WORK / f"bridget-2005-{layout}-adapter"
    output = output_root / "M16" / "mc_idx_tbl"
    output.mkdir(parents=True, exist_ok=True)
    manifest = [
        "source=Bridget 2005 indexes; source files remain read-only",
        "adapter=ver.2005 to ver.2013",
        "attr_b=zero; metric columns preserved; producer tag normalized to VoiceText-Eng",
        f"layout={layout}",
        "semantic_status=experimental mapping, not a recovered native Bridget 2013 format",
    ]
    if layout == "signature":
        manifest.insert(2, "signature=[attr_48,key[0:3],attr_40,key[3:5]]")
    elif layout == "selector-key3-attr40":
        manifest.insert(2, "signature=[attr_48,key[0:3],attr_40,key[3],attr_40]")
    elif layout == "selector-key4-attr40":
        manifest.insert(2, "signature=[attr_48,key[0:3],attr_40,key[4],attr_40]")
    else:
        manifest.insert(2, "signature=[attr_48,key[0:5],attr_40]")

    for source in sorted(SOURCE.glob("unit-*.idx")):
        raw = source.read_bytes()
        parsed = parse_index_data(source, raw)
        if parsed.version != "ver.2005" or parsed.producer != "VoiceText-Bre":
            raise ValueError(
                f"{source.name}: expected ver.2005/VoiceText-Bre, found "
                f"{parsed.version}/{parsed.producer}"
            )
        count = parsed.unit_count
        block_end = parsed.block_start + count * parsed.record_stride
        tail = raw[block_end:]
        if len(tail) != count * 20:
            raise ValueError(f"{source.name}: expected 20 feature bytes per unit")

        attr_4c = tail[:count]
        attr_48 = tail[count : 2 * count]
        key5 = tail[2 * count : 7 * count]
        attr_40 = tail[7 * count : 8 * count]
        metrics = tail[8 * count :]
        signature = bytearray()
        for unit in range(count):
            key = key5[unit * 5 : unit * 5 + 5]
            if layout == "signature":
                signature.extend(
                    (attr_48[unit], *key[:3], attr_40[unit], *key[3:])
                )
            elif layout == "selector-key3-attr40":
                signature.extend(
                    (attr_48[unit], *key[:3], attr_40[unit], key[3], attr_40[unit])
                )
            elif layout == "selector-key4-attr40":
                signature.extend(
                    (attr_48[unit], *key[:3], attr_40[unit], key[4], attr_40[unit])
                )
            else:
                signature.extend((attr_48[unit], *key, attr_40[unit]))

        header_start = 1
        header_end = header_start + raw[0]
        header = raw[header_start:header_end]
        converted_header_bytes = header.replace(b"ver.2005", b"ver.2013", 1).replace(
            b"VoiceText-Bre", b"VoiceText-Eng", 1
        )
        if len(converted_header_bytes) != len(header):
            raise ValueError(f"{source.name}: header replacement changed its size")

        converted_tail = attr_4c + bytes(signature) + bytes(count) + metrics
        converted = (
            raw[:1]
            + converted_header_bytes
            + raw[header_end:block_end]
            + converted_tail
        )
        target = output / source.name
        target.write_bytes(converted)
        report = inspect_versioned(target)
        if report["unit_count"] != count:
            raise ValueError(f"{source.name}: unit count changed")
        manifest.append(
            f"{source.name}: units={count} bytes={len(raw)}->{len(converted)} "
            f"producer={parsed.producer}->VoiceText-Eng"
        )

    (output_root / "manifest.txt").write_text("\n".join(manifest) + "\n")
    print("\n".join(manifest))


if __name__ == "__main__":
    main()
