#!/usr/bin/env python3
"""Generate and verify MakeCsv input/output overlap observations."""

import argparse
import re
from pathlib import Path


def cases():
    return [(offset, cap) for offset in range(10) for cap in range(offset + 3, 17)]


def expected(offset, cap):
    if offset == 0:
        quotes = 1 + 2 * ((cap - 2) // 2)
        return -1, b'"' * quotes + b"1" * (cap - 1 - quotes) + b"\0"
    if offset == 1:
        return -1, b'"' + b"1" * (cap - 2) + b"\0"
    encoded = b'"' + b"1" * (cap - offset - 1) + b'"\0'
    if offset == 2:
        return 1, encoded
    return 1, encoded + b"1" * (offset - 3) + b"\0"


def generate(path):
    lines = ["set pagination off", "set confirm off", "set debuginfod enabled off",
             "handle SIGSEGV nostop noprint pass", "break *0x1001da50", "commands 1",
             "  silent", "  disable 1", "  set $fields = (char **)malloc(4)",
             "  set $raw = (unsigned char *)malloc(18)", "  set $out = $raw + 1"]
    args = ", ".join(f"$raw[{i}]" for i in range(18))
    for offset, cap in cases():
        lines.extend([
            "  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)",
            f"  set $fields[0] = (char *)($out + {offset})",
            f"  set $out[{offset}] = 65", f"  set $out[{offset + 1}] = 66",
            f"  set $out[{offset + 2}] = 0",
            f"  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, {cap})",
            f'  printf "CSV_ALIAS offset={offset} cap={cap} ret=%d raw={"%02x" * 18} pointer_same=%d\\n", $ret, {args}, $fields[0] == (char *)($out + {offset})',
        ])
    lines.extend([f'  printf "CSV_ALIAS_DONE calls={len(cases())}\\n"', "  continue", "end", "continue"])
    path.write_text("\n".join(lines) + "\n")
    print(f"generated_overlap_calls={len(cases())}")


def verify(path):
    required = set(cases())
    seen = set()
    text = path.read_text()
    if not re.search(r"\[Inferior .*exited normally\]", text):
        raise ValueError("normal host exit not confirmed")
    for match in re.finditer(r"CSV_ALIAS offset=(\d+) cap=(\d+) ret=(-?\d+) raw=([0-9a-f]{36}) pointer_same=(\d+)", text):
        offset, cap, ret = map(int, match.group(1, 2, 3))
        identity = (offset, cap)
        if identity not in required or identity in seen:
            raise ValueError(f"unexpected or duplicate capture: {identity}")
        expected_ret, output = expected(offset, cap)
        raw = b"\xa5" + output + b"\xa5" * (17 - cap)
        if ret != expected_ret or bytes.fromhex(match[4]) != raw or match[5] != "1":
            raise ValueError(f"mismatch: {match[0]}")
        seen.add(identity)
    if seen != required or f"CSV_ALIAS_DONE calls={len(required)}" not in text:
        raise ValueError(f"incomplete capture: matched {len(seen)} of {len(required)}")
    print(f"verified_overlap_calls={len(seen)} raw_bytes=exact array_pointer=unchanged")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("mode", choices=("generate", "verify"))
    parser.add_argument("path", type=Path)
    args = parser.parse_args()
    (generate if args.mode == "generate" else verify)(args.path)


if __name__ == "__main__":
    main()
