# 2006 Kate MSI DLL full static reference

This is a local, ignored analysis bundle for the extracted MSI's
`lib/vt_eng.dll`. The source DLL was read without modification. Its SHA-256
before and after analysis is
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.

## Coverage

- Ghidra 12.1.4 imported and auto-analyzed the PE32 DLL in a separate project.
- Ghidra discovered 947 functions. The all-function script produced
  pseudocode for all 947; the inventory records zero decompilation failures.
  This is coverage of Ghidra's discovered functions, not proof that every
  executable byte was correctly identified as code or that pseudocode is
  original source.
- GNU `objdump` produced a complete PE metadata/import/export dump and
  raw-byte x86 disassembly of executable sections.
- The bundle includes `file`, `strings`, FLOSS 3.1.1, capa 9.4.0, and pefile
  2024.8.26 output. FLOSS reports 9,046 static strings, 6 stack strings,
  0 tight strings, and 2 decoded strings. Treat FLOSS and capa matches as
  triage clues; they do not establish engine behavior or maliciousness.
- The PE export table contains 44 named exports. `objdump-pe.txt` and
  `pefile-dump.txt` provide the full tables.

## Files

- `all-functions-pseudocode.c`: Ghidra pseudocode for all discovered functions.
- `function-inventory.tsv`: function address, size, thunk/external status, and
  decompilation result.
- `objdump-disassembly.txt`: complete `objdump -d -Mintel --show-raw-insn`
  output for disassembled executable sections.
- `objdump-pe.txt`, `pefile-dump.txt`: PE headers, sections, imports, exports,
  and data directories.
- `strings.txt`, `floss.txt`, `capa.txt`: string and behavior-scanner output.
- `ghidra/`: the isolated headless Ghidra project.
- `scripts/` and `run-static-kit.sh`: the all-function Ghidra script, pefile
  dumper, and suite runner.

## Reproduction

From `tools/revkit/`, with the local MSI extraction present, run:

```sh
docker compose run --rm revtools \
  /bin/bash /work/corpus-parity/full-analysis-2006/run-static-kit.sh
```

The `ghidra` argument reruns only the Ghidra import/decompilation stage. Use a
fresh output directory for a fresh Ghidra project.
