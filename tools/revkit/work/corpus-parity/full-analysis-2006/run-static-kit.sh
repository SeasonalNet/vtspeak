#!/usr/bin/env bash
set -euo pipefail

dll='/work/corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll'
out='/work/corpus-parity/full-analysis-2006'

mkdir -p "$out/ghidra"
if [[ "${1:-all}" != ghidra ]]; then
  file "$dll" > "$out/file.txt"
  sha256sum "$dll" > "$out/input-sha256.txt"
  objdump -f -h -p "$dll" > "$out/objdump-pe.txt"
  objdump -d -Mintel --show-raw-insn "$dll" > "$out/objdump-disassembly.txt"
  strings -a -n 4 "$dll" > "$out/strings.txt"
  floss "$dll" > "$out/floss.txt"
  capa --rules /opt/capa-rules --signatures /opt/capa-sigs "$dll" > "$out/capa.txt"
  /opt/capa-venv/bin/python "$out/scripts/dump_pefile.py"
fi

/opt/ghidra/support/analyzeHeadless \
  "$out/ghidra" Kate2006Full \
  -import "$dll" \
  -scriptPath "$out/scripts" \
  -postScript DumpAllFunctions.java \
    "$out/all-functions-pseudocode.c" \
    "$out/function-inventory.tsv"
