#!/usr/bin/env python3
"""Generate and independently verify a field/count/capacity MakeCsv probe."""

import argparse
import itertools
import re
from pathlib import Path


TOKENS = {"a": b"A", "e": b"", "q": b'""', "h": b'\x80"'}


def cases():
    result = []
    for count in range(6):
        for fields in itertools.product("ae", repeat=count):
            result.append((fields, list(range(1, 33))))
    for count in (8, 16, 32, 64, 128):
        caps = sorted(set((1, 2, 3, 4, 8, 16, 32, 64, 128, 256, 512,
                           4 * count - 1, 4 * count, 4 * count + 1)))
        result.append((tuple("a" * count), caps))
        for position in (0, count // 2, count - 1):
            fields = list("a" * count)
            fields[position] = "e"
            result.append((tuple(fields), caps))
    for fields in ("qh", "hq", "aqha", "hqeh"):
        result.append((tuple(fields), list(range(1, 33))))
    return result


def expected(fields, capacity):
    output = bytearray(b"1" * capacity)
    output[-1] = 0
    cursor = 0

    def emit(chunk):
        nonlocal cursor
        if len(chunk) > capacity - cursor or 0 in output[cursor:cursor + len(chunk)]:
            return False
        output[cursor:cursor + len(chunk)] = chunk
        cursor += len(chunk)
        return True

    for field_index, key in enumerate(fields):
        if not emit(b'"'):
            return -1, bytes(output)
        field = TOKENS[key]
        if not field:
            return -1, bytes(output)
        offset = 0
        while offset < len(field):
            byte = field[offset]
            if byte & 0x80 and offset + 1 < len(field):
                chunk = field[offset:offset + 2]
                offset += 2
            else:
                chunk = b'""' if byte == 0x22 else bytes((byte,))
                offset += 1
            if not emit(chunk):
                return -1, bytes(output)
        if not emit(b'"'):
            return -1, bytes(output)
        if field_index + 1 < len(fields) and not emit(b","):
            return -1, bytes(output)
    output[cursor] = 0
    return 1, bytes(output)


def generate(path):
    lines = ["set pagination off", "set confirm off", "set debuginfod enabled off",
             "handle SIGSEGV nostop noprint pass", "break *0x1001da50", "commands 1",
             "  silent", "  disable 1", "  set $fields = (char **)malloc(512)",
             "  set $raw = (unsigned char *)malloc(1026)", "  set $out = $raw + 1"]
    for key, value in TOKENS.items():
        lines.append(f"  set ${key} = (unsigned char *)malloc({len(value) + 1})")
        for index, byte in enumerate(value + b"\0"):
            lines.append(f"  set ${key}[{index}] = {byte}")
    calls = 0
    for case_id, (fields, caps) in enumerate(cases()):
        for index, key in enumerate(fields):
            lines.append(f"  set $fields[{index}] = (char *)${key}")
        for cap in caps:
            lines.extend([
                "  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 1026)",
                f"  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, {len(fields)}, $out, {cap})",
                f'  printf "CSV_FIELD case={case_id} count={len(fields)} cap={cap} ret=%d output=", $ret',
            ])
            for start in range(0, cap, 32):
                size = min(32, cap - start)
                args = ", ".join(f"$out[{index}]" for index in range(start, start + size))
                lines.append(f'  printf "{"%02x" * size}", {args}')
            lines.append(f'  printf " prefix=%02x post=%02x\\n", $raw[0], $out[{cap}]')
            calls += 1
    lines.extend([f'  printf "CSV_FIELD_DONE calls={calls}\\n"', "  continue", "end", "continue"])
    path.write_text("\n".join(lines) + "\n")
    print(f"generated_cases={len(cases())} calls={calls} trace={path}")


def verify(path):
    matrix = cases()
    required = {(index, cap) for index, (_, caps) in enumerate(matrix) for cap in caps}
    seen = set()
    for match in re.finditer(r"CSV_FIELD case=(\d+) count=(\d+) cap=(\d+) ret=(-?\d+) output=([0-9a-f\s]+) prefix=([0-9a-f]{2}) post=([0-9a-f]{2})", path.read_text()):
        case_id, count, cap, ret = map(int, match.group(1, 2, 3, 4))
        identity = (case_id, cap)
        if identity not in required or identity in seen:
            raise ValueError(f"unexpected or duplicate record {identity}")
        fields, _ = matrix[case_id]
        expected_ret, expected_output = expected(fields, cap)
        if count != len(fields) or ret != expected_ret or bytes.fromhex(match[5]) != expected_output:
            raise ValueError(f"output mismatch {identity}: {match[0]}")
        if (match[6], match[7]) != ("a5", "a5"):
            raise ValueError(f"guard mismatch {identity}")
        seen.add(identity)
    if seen != required:
        raise ValueError(f"missing {len(required - seen)} records")
    if f"CSV_FIELD_DONE calls={len(required)}" not in path.read_text():
        raise ValueError("missing completion marker")
    print(f"verified_cases={len(matrix)} calls={len(seen)} outputs=exact guards=intact")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("mode", choices=("generate", "verify"))
    parser.add_argument("path", type=Path)
    args = parser.parse_args()
    (generate if args.mode == "generate" else verify)(args.path)


if __name__ == "__main__":
    main()
