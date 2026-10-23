# Stage 22: G83 query-signature bridge

This paired trace follows the existing `numeric-context-g83-can-anybody-help`
fixture through a narrower boundary than the Stage 12 score/candidate probes.
It records writes of numeric token code 83 at `FUN_1000e0c0`, then records
seven-byte query signatures at the class lookup `FUN_10023dc0`. REA's
decompilation of the immediately called `FUN_10016ea0` shows how those bytes
become the five-byte key. The intervention run zeros only
the `FUN_1000e0c0` code-83 argument, as in the prior Stage 12 controls.

The paired capture is reproducible with the repository's isolated, no-network
Wine image. Vendor DLLs and voice data remain read-only; the runner backs up
and restores the Stage 5 input and preserves its original output-file state.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage22/run-g83-query.sh control \
  numeric-context-g83-can-anybody-help
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage22/run-g83-query.sh g83-zero \
  numeric-context-g83-can-anybody-help
```

Each run completed normally and recorded 48 seven-byte signatures. In the
paired comparison, 41 of 48 signatures differ; the earliest difference is
signature byte 5 (zero-based) at hit 5. The changed positions across the full
capture are bytes 0 through 5. Byte 6 is zero in all 96 records. REA's
decompilation shows signature byte 5 is copied directly into key byte 3,
while signature byte 6 is masked with `0x20` into key byte 4. Thus this probe
ties the G83 intervention to a class-lookup key component. It does not provide
a phonetic or acoustic name for that component.

The control and intervention WAVs match the previously recorded Stage 12
outputs (`29052da3...` and `f8aa3133...`, respectively), so this query-boundary
instrumentation did not introduce a new output result. Detailed counts,
representative bytes, source-function mapping, and capture hashes are in the
curated [runtime G83 query-signature report](../../../docs/reverse-engineering/rea-runtime-g83-query-signature-2026-10-19.md).

## Rebuild the ignored assets

The two GDB command files and three WAV files are hash-checked entries in the
shared artifact manifest. From the repository root, verify that their sources
still reproduce the recorded bytes:

```sh
python3 tools/revkit/scripts/reproduce_artifacts.py verify --kind gdb \
  --path tools/revkit/work/stage22/g83-query-control.gdb \
  --path tools/revkit/work/stage22/g83-query-g83-zero.gdb
python3 tools/revkit/scripts/reproduce_artifacts.py verify --kind wav \
  --runtime-root /tmp/vtspeak-stage22-wavs \
  --path tools/revkit/work/stage22/numeric-context-g83-can-anybody-help-g83-query-control.wav \
  --path tools/revkit/work/stage22/numeric-context-g83-can-anybody-help-g83-query-g83-zero.wav \
  --path tools/revkit/work/stage22/diagnostic/numeric-context-g83-can-anybody-help-effect-candidate-path-control.wav
```

Use the matching `restore --kind ... --path ...` command to materialize a
missing asset. Reproduce just the Stage 22 capture batch into scratch storage:

```sh
python3 tools/revkit/scripts/reproduce_wavs.py check-inputs --batch stage22
python3 tools/revkit/scripts/reproduce_wavs.py run --batch stage22 \
  --scratch /tmp/vtspeak-stage22-wavs
```

The batch checks only the Paul voice, common data, executable, and headers it
uses; it does not require unrelated Stage 19/20 overlays. The hash-checked
WAV recipes require the local vendor inputs and isolated Wine image. To repeat
the query traces directly, use the runtime commands above; the runner writes
fresh logs and WAVs under Stage 22.
