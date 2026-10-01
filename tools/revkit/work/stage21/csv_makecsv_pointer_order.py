#!/usr/bin/env python3
"""Generate and verify isolated MakeCsv pointer-access-order captures."""

import argparse
import re
from pathlib import Path


# name, array kind, count, capacity, outcome, return/fault EIP, output bytes
CASES = (
    ("null_array_cap1", "null_array", 1, 1, "fault", 0x100169A7, b"\0"),
    ("null_first_cap1", "null_first", 1, 1, "return", -1, b"\0"),
    ("null_first_cap2", "null_first", 1, 2, "fault", 0x10016A53, b'"\0'),
    ("bad_first_cap1", "bad_first", 1, 1, "return", -1, b"\0"),
    ("bad_first_cap2", "bad_first", 1, 2, "fault", 0x10016A53, b'"\0'),
    ("late_null_cap4", "late_null", 2, 4, "return", -1, b'"A"\0'),
    ("late_null_cap5", "late_null", 2, 5, "return", -1, b'"A",\0'),
    ("late_null_cap6", "late_null", 2, 6, "fault", 0x10016A53, b'"A","\0'),
    ("empty_then_null", "empty_then_null", 2, 8, "return", -1, b'"111111\0'),
    ("empty_huge_count", "empty_then_null", 2147483647, 8, "return", -1, b'"111111\0'),
    ("negative_null_array", "null_array", -2147483648, 1, "return", 1, b"\0"),
    ("guard_array_cap4", "guard_array", 2, 4, "return", -1, b'"A"\0'),
    ("guard_array_cap5", "guard_array", 2, 5, "fault", 0x100169A7, b'"A",\0'),
    ("guard_array_cap6", "guard_array", 2, 6, "fault", 0x100169A7, b'"A",1\0'),
    ("guard_empty_huge", "guard_empty", 2147483647, 8, "return", -1, b'"111111\0'),
)


def generate(directory):
    for name, kind, count, cap, _, _, _ in CASES:
        lines = ["set pagination off", "set confirm off", "set debuginfod enabled off",
                 "handle SIGSEGV stop print nopass", "break *0x1001da50", "commands 1",
                 "  silent", "  disable 1", "  set $fields = (char **)malloc(8)",
                 "  set $a = (unsigned char *)malloc(2)", "  set $a[0] = 65", "  set $a[1] = 0",
                 "  set $empty = (unsigned char *)malloc(1)", "  set $empty[0] = 0",
                 "  set $fields[0] = (char *)$a", "  set $fields[1] = 0",
                 "  set $raw = (unsigned char *)malloc(16)",
                 "  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 16)",
                 "  set $out = $raw + 1", "  set $stack = (unsigned char *)malloc(16384)"]
        if kind == "null_array":
            lines.append("  set $fields = (char **)0")
        elif kind == "null_first":
            lines.append("  set $fields[0] = 0")
        elif kind == "bad_first":
            lines.append("  set $fields[0] = (char *)0xffffffff")
        elif kind == "empty_then_null":
            lines.append("  set $fields[0] = (char *)$empty")
        elif kind in ("guard_array", "guard_empty"):
            lines.extend([
                "  set $allocation = ((char *(*)(void *, unsigned int, unsigned int, unsigned int))0x7b68d8a0)(0, 12288, 12288, 4)",
                "  set $old_protect = (unsigned int *)malloc(4)",
                "  set $protect_ok = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))0x7b68df20)($allocation + 4096, 4096, 1, $old_protect)",
                "  set $fields = (char **)($allocation + 4092)",
                "  set $fields[0] = (char *)$" + ("empty" if kind == "guard_empty" else "a"),
                '  printf "CSV_PTR_GUARD allocation=%08x protect_ok=%d\\n", $allocation, $protect_ok',
            ])
        args = ", ".join(f"$raw[{i}]" for i in range(16))
        lines.extend([
            '  printf "CSV_PTR_BEGIN name=' + name + f' count={count} cap={cap}\\n"',
            "  break *0x10016c6d", "  commands 2", "    silent",
            f'    printf "CSV_PTR_RETURN name={name} low_ax=%d raw={"%02x" * 16}\\n", (short)$eax, {args}',
            "    kill", "    quit", "  end", "  catch signal SIGSEGV", "  commands 3", "    silent",
            f'    printf "CSV_PTR_FAULT name={name} eip=%08x raw={"%02x" * 16}\\n", $eip, {args}',
            "    kill", "    quit", "  end", "  set $esp = (unsigned int)$stack + 16128",
            "  set {unsigned int}$esp = 0x10016c6d",
            "  set {unsigned int}($esp + 4) = (unsigned int)$fields",
            f"  set {{int}}($esp + 8) = {count}",
            "  set {unsigned int}($esp + 12) = (unsigned int)$out",
            f"  set {{unsigned int}}($esp + 16) = {cap}",
            "  set $eip = 0x10016c50", "  continue", "end", "continue",
        ])
        (directory / f"trace-csv-pointer-order-{name}.gdb").write_text("\n".join(lines) + "\n")
    print(f"generated_pointer_order_cases={len(CASES)}")


def verify(directory, run_id):
    for name, kind, _, _, outcome, value, output in CASES:
        path = directory / f"csv-pointer-order-{run_id}-{name}-api.log"
        text = path.read_text()
        if not re.search(r"\[Inferior .* killed\]", text):
            raise ValueError(f"process termination not confirmed: {name}")
        exits = re.findall(r"DEBUGGER_EXIT code=(\d+)", text)
        if exits and exits != ["0"]:
            raise ValueError(f"debugger exit failed: {name}: {exits}")
        if len(re.findall(r"CSV_PTR_(?:RETURN|FAULT) name=", text)) != 1:
            raise ValueError(f"missing or repeated terminal outcome: {name}")
        if kind.startswith("guard"):
            setup = re.search(r"CSV_PTR_GUARD allocation=([0-9a-f]{8}) protect_ok=1", text)
            if setup is None or int(setup[1], 16) == 0:
                raise ValueError(f"guard-page setup not confirmed: {name}")
        if outcome == "return":
            expression = rf"CSV_PTR_RETURN name={name} low_ax=(-?\d+) raw=([0-9a-f]{{32}})"
        else:
            expression = rf"CSV_PTR_FAULT name={name} eip=([0-9a-f]{{8}}) raw=([0-9a-f]{{32}})"
        matches = list(re.finditer(expression, text))
        if len(matches) != 1:
            raise ValueError(f"missing or repeated outcome: {name}")
        match = matches[0]
        actual_value = int(match[1], 10 if outcome == "return" else 16)
        expected_raw = b"\xa5" + output + b"\xa5" * (15 - len(output))
        if actual_value != value or bytes.fromhex(match[2]) != expected_raw:
            raise ValueError(f"mismatch: {name}: {match[0]}")
        print(f"verified={name} outcome={outcome}")
    print(f"verified_pointer_order_cases={len(CASES)} raw_bytes=exact")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("mode", choices=("generate", "verify"))
    parser.add_argument("directory", type=Path)
    parser.add_argument("run_id", nargs="?", default="matrix1")
    args = parser.parse_args()
    if args.mode == "generate":
        generate(args.directory)
    else:
        verify(args.directory, args.run_id)


if __name__ == "__main__":
    main()
