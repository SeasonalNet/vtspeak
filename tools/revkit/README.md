# VoiceText static analysis tools

This isolated Docker toolset contains Ghidra headless analysis, FLOSS, capa,
PEfile, `objdump`, `strings`, and `file`. The image does not copy the VoiceText
binary or model data. Compose mounts `binary/` read-only at `/samples`, the
shared and Paul model directories read-only under `/voice-data/`, and `include/`
read-only at `/voice-data/include`. Portable runtime captures, small fixtures,
trace scripts, and decompiler reports cited by the research notes are tracked
under this directory's `work/` folder. Host-local Ghidra projects, caches,
Wine prefixes, and full vendor-input copies remain ignored; see
[`work/README.md`](work/README.md).

## Build

From this directory:

```sh
docker compose build
```

The image pins Ghidra 12.1.4 and checks its SHA-256 against the value published
with the [official Ghidra release](https://github.com/NationalSecurityAgency/ghidra/releases/tag/Ghidra_12.1.4_build).
It also pins capa 9.4.0, the matching [capa-rules v9.4.0 release](https://github.com/mandiant/capa-rules/releases/tag/v9.4.0),
and the matching FLIRT signatures from the [capa v9.4.0 source](https://github.com/mandiant/capa/releases/tag/v9.4.0).
The rules and signatures archive SHA-256 values are verified during the build.
The FLOSS and PEfile versions are pinned in the Dockerfile too.

## Analyze the executable and DLL with Ghidra

From this directory, import both PE files into a persistent Ghidra project:

```sh
docker compose run --rm \
  --entrypoint /opt/ghidra/support/analyzeHeadless \
  revtools /work/projects VoiceText \
  -import /samples/voicetext_paul.exe /samples/vt_pau.dll
```

The project is saved as `work/projects/VoiceText.gpr` with its companion
`VoiceText.rep` directory. To inspect it with Ghidra's GUI, open a Ghidra
installation on the host and open that project, or run a GUI-enabled Ghidra
setup separately.

## Other analyses

```sh
docker compose run --rm --entrypoint floss revtools /samples/vt_pau.dll
docker compose run --rm --entrypoint capa revtools \
  --rules /opt/capa-rules --signatures /opt/capa-sigs /samples/vt_pau.dll
docker compose run --rm --entrypoint objdump revtools -p /samples/vt_pau.dll
docker compose run --rm --entrypoint strings revtools -a /samples/vt_pau.dll
```

## Inspect the versioned unit-index layout

The read-only inspector checks the `ver.2013` header, bank names, unit count,
19-byte record stride, and packed-column extents. With `--data-dir`, it also
checks that DAT and UPM record spans partition their matching files through
EOF. The UPM span length is the sum of the two stored counts minus one; the
two runtime views share a byte. Fields whose use has not been traced remain
unnamed. From the repository root, run:

```sh
python3 tools/revkit/scripts/inspect_unit_idx.py \
  --data-dir data-paul/M16/dat \
  data-paul/M16/mc_idx_tbl/unit-gen.idx \
  data-paul/M16/mc_idx_tbl/unit-num.idx \
  data-paul/M16/mc_idx_tbl/unit-etc.idx \
  data-paul/M16/mc_idx_tbl/unit-alp.idx
```

## Inspect Paul `tree3` files

The read-only parser checks the observed header, node/list records, output
rows, and child/leaf references. It can also evaluate a caller-supplied signed
feature vector; use a vector whose entries correspond to the selected tree's
feature selectors. For example, list the shape of every Paul duration and
pitch tree from the repository root:

```sh
python3 tools/revkit/scripts/tree3.py \
  data-paul/M16/ttsdata/tree3/duration/*.tree3 \
  data-paul/M16/ttsdata/tree3/pitch/*.tree3
```

The evaluator is an analysis aid for the observed 2013 package. It does not
assign phonetic names or physical units to feature slots or output values.

## Inspect shared English dictionary resources

The read-only inspector validates the three standalone shared `tree3` files,
the 27-tree `atmt.tree3` container, both hash-indexed dictionary record
partitions, the structured `exceptdict`, and the shared `.txt2` files. It
subtracts each `.txt2` shift byte modulo 256, validates the padded decimal row
count and loader-selected column mode, and prints decoded sample rows. From
the repository root, run:

```sh
python3 tools/revkit/scripts/dict_resources.py
```

To print every decoded `.txt2` row in the JSON output, add `--show-rows`.

Decode the TPP dictionary's compressed keys and validate every decoded-key
round trip and typed payload shape with the matching Paul DLL:

```sh
python3 tools/revkit/scripts/inspect_tpp_dictionary.py
```

Use `--show-rows` to include all 31,550 decoded key/payload pairs in the JSON
report. The script validates structural codes but does not assign
pronunciation or linguistic labels to every TPP code.

Compare captured GDB scalar/vector tree returns against those shared resources
and the Paul voice trees with:

```sh
python3 tools/revkit/scripts/compare_tree3_runtime.py \
  tools/revkit/work/stage5/probes/stage6/ordinary/tree-lookups-gdb.log \
  tools/revkit/work/stage5/probes/stage6/numbers/tree-lookups-gdb.log \
  tools/revkit/work/stage5/probes/stage6/abbreviations/tree-lookups-gdb.log
```

The comparison reports only a unique tree whose evaluation matches the
captured return values for the observed input feature vector. It does not
assign phonetic labels to those features.

Start a shell with the tools available on `PATH`:

```sh
docker compose run --rm revtools
```
