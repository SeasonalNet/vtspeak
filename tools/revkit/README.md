# VoiceText static analysis tools

This isolated Docker toolset contains Ghidra headless analysis, FLOSS, capa,
PEfile, `objdump`, `strings`, and `file`. The image does not copy the VoiceText
binary or model data. Compose mounts `binary/` read-only at `/samples`, the
shared and Paul model directories read-only under `/voice-data/`, and `include/`
read-only at `/voice-data/include`. It stores Ghidra projects and raw reports
under this directory's ignored `work/` folder.

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

Start a shell with the tools available on `PATH`:

```sh
docker compose run --rm revtools
```
