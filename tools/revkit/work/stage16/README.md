# Stage 16: Lead 6 API compatibility probe

This stage captures the local Paul's file, buffer, information,
configuration, and selected playback API behavior. It uses the existing Stage
8 Wine runtime container and calls APIs in the loaded process at the
`VT_TextToFile_ENG` entry breakpoint (`0x1001da50`). Vendor binary and model
mounts remain read-only.

The format runner accepts one selector from 0 through 10:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-file-format.sh 4
```

The error runner accepts `null-text`, `empty-text`, or `null-path`:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-file-error.sh null-text
```

The buffer runner accepts a format, flag, and thread ID:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer.sh 0 0 0
```

It calls `VT_TextToBuffer_ENG` in the already loaded process with a 1 MiB
buffer. It writes the captured buffer only when the call returns success, so
failed thread-creation cases do not expose uninitialized process memory.

`run-buffer-errors.sh` captures the buffer API's invalid-format, null-text,
empty-text, and null-buffer returns. `run-info.sh` queries all declared
`VT_GetTTSInfo_ENG` requests plus invalid-request, null-value, and string-size
errors in one loaded process:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-info.sh
```

`run-config.sh` queries the current speaker settings, applies upper-bound
values through the configuration setters, checks the resulting getters, then
captures file output with the synthesis arguments left at `-1`.
`run-config-fields.sh` resets the initial slot-1 values before each setter
boundary case, then captures getters and WAVE output for the pitch, speed,
volume, sentence-pause, and comma-pause minima and maxima.
`run-play-errors.sh` records null/empty/valid-text playback returns and calls
stop, pause, and restart controls. `run-play-waveout.sh` repeats valid
playback with `alsa-null.conf` configured as the ALSA default; this allows
WinMM initialization and active pause/restart calls in Wine while discarding
audio instead of producing audible sound.
`run-db-unloaded.sh` unloads speaker 1 inside the probe process and records the
file, synchronous buffer, and playback API returns for valid text.
`run-unload-ext-loaded.sh` calls the extended unload export directly after
checking loaded speaker 1's size. It records the size query, static speaker
name, and file/buffer results after unloading the loaded engine. The size query
changes from 508,121,688 bytes to `-1`; file and buffer calls return `-5` and
`-6` respectively.
`run-load-ext-wrapper.sh` captures the sample's call into
`VT_LOADTTS_EXT_ENG`: slot `-1`, null voice path, default extended arguments,
low AX 0, and the resulting 508,121,688-byte Paul database.
`run-load-ext-default-path.sh` records the selected default base (`../`).
`run-load-ext-explicit-runtime-root.sh` substitutes its absolute runtime
equivalent (`Z:\work\`) and confirms normalization to `Z:/work/` plus the
same successful database size. The two leaf candidates captured by
`run-load-ext-explicit-db-path.sh` and
`run-load-ext-explicit-parent-path.sh` return low AX 3 before slot state is
allocated. Together, these traces show that this parameter is a resource base
for the package's relative paths rather than the voice-model leaf directory.
The captures record license argument pointers and lengths only, not license
contents.
`run-load-ext-license-file-state.sh` passes the existing 468-byte
verification-file path as argument 6 with null in-memory text and length
`-1`. The checker returns `-11`, while the loader still returns low AX 0 and
loads the 508,121,688-byte model; slot 1 has license gate 0 and dictionary
capacity 1 afterward. Its GDB trace initializes the pathname with explicit
bytes and records statuses/state only. The earlier
`load-ext-license-file-api.log` has a corrupted injected path and is retained
as an invalid setup attempt; use the bytewise runner for reproduction.
`probe-load-ext-memory-license.c` reads that same file into a heap buffer,
passes it as extended-loader argument 7 with argument 6 null and the exact
byte length in argument 8, then calls the checker against the loaded speaker.
It prints only byte count, return statuses, database size, license gate, and
capacity. The result matches the file-path case: load AX 0, checker `-11`,
gate 0, capacity 1.

Build and run the memory-buffer trace with the PE32 builder and Wine runtime:

```sh
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage16:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage16/probe-load-ext-memory-license.c \
  -o /out/probe-load-ext-memory-license.exe
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-load-ext-memory-license-state-container.sh
```

`probe-load-ext-state.c` directly calls the extended loader and records return
codes, per-slot state occupancy, and database sizes without reading or printing
license data. Its `reload` mode checks an identical second load, unload/reload,
and queries database size for all six slots; `path-switch` supplies a different
nonexistent base while slot 1 is loaded; `multi` loads Paul and James together
in both orders, unloading each while querying the remaining slot; an integer
argument runs one first-call slot probe.
Build it with the same PE32 builder pattern, changing the source and output to
`probe-load-ext-state.c` and `probe-load-ext-state.exe`:

```sh
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage16:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage16/probe-load-ext-state.c \
  -o /out/probe-load-ext-state.exe
```

Then run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-load-ext-state-container.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-load-ext-path-switch-container.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage16/run-load-ext-slot-matrix-corrected-container.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage16/run-load-ext-multi-container.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage16/run-load-ext-multi-path-switch-container.sh
for slot in -2 -1 6 1000; do
  docker compose -f tools/revkit/work/stage8/compose.yaml \
    -f tools/revkit/work/stage17/compose.yaml run --rm \
    -v "$PWD/data-james:/work/data-james:ro" runtime \
    /bin/bash /work/stage16/run-load-ext-invalid-switch-container.sh "$slot"
done
for slot in 0 2 3 5; do
  docker compose -f tools/revkit/work/stage8/compose.yaml \
    -f tools/revkit/work/stage17/compose.yaml run --rm \
    -v "$PWD/data-james:/work/data-james:ro" runtime \
    /bin/bash /work/stage16/run-load-ext-prosody-filetrace-container.sh "$slot"
done
```

Stage 17 mounts James at `/work/data-jame` for its supplied host/DLL. This
Paul DLL requests `/work/data-james`, so the additional read-only alias is
required for slot 4. The final matrix loaded slots 1 (Paul) and 4 (James), with
sizes 508,121,688 and 253,797,273 bytes. Slots 0 and 3 fail on the missing
`data-kate/.../tree3/pitch/nbt.tree3` and
`data-julie/.../tree3/pitch/nbt.tree3`; their local packages provide `tree2`.
Slots 2 and 5 request equivalent first-tree paths below `data-em001` and
`data-ashley`, for which no local model roots exist. The first two slot-matrix
logs are retained as incomplete mount attempts; use
`load-ext-slot-matrix-alias-api.log`, `load-ext-multi-api.log`, and the
`load-ext-slot-*-prosody-filetrace-v2-api.log` captures for corrected results.
The multi-slot probe loaded both voices in either order, and unloading one
preserved the other slot and its database-size query.
The multi-path-switch probe tries a nonexistent base against each loaded slot,
then retries both slots with the original base and records occupancy and sizes.
The invalid-slot path-switch probe first normalizes an out-of-range slot to
Paul/slot 1, then repeats that same input with a different base and finally
loads James/slot 4 with the original base. It captures return codes, occupancy,
and sizes so invalid-index behavior can be correlated with whether both valid
voice states survive; each input runs in its own Wine process.

`make-load-ext-error-overlays.py OUTPUT` creates symlink-only views for
resource-failure probes. Each omitted file is left out while sibling files and
directories link to the read-only vendor root mounted under `/vendor-paul` or
`/vendor-common`. Mount the selected overlay read-only at the corresponding
`/work` path, then run
`run-load-ext-resource-error-container.sh CASE`. Cases map as follows:

| Cases | Overlay target | Vendor mount |
| --- | --- | --- |
| `omit-atmt-tree`, `omit-engbi-tree`, `omit-sbd-tree`, `omit-tppdict`, `omit-hashidx-tpp` | `/work/data-common/dict-eng` | `data-common:/vendor-common:ro` |
| `omit-cepdist-tbl` | `/work/data-paul/M16/ttsdata/dist_tbl` | `data-paul:/vendor-paul:ro` |
| `omit-pitch-nbt` | `/work/data-paul/M16/ttsdata/tree3/pitch` | `data-paul:/vendor-paul:ro` |
| `omit-unit-gen-idx` | `/work/data-paul/M16/mc_idx_tbl` | `data-paul:/vendor-paul:ro` |
| `omit-dblist-idx-preserve-tree` | `/work/data-paul/M16` | `data-paul:/vendor-paul:ro` |
| `omit-gen-dat`, `omit-gen-upm` | `/work/data-paul/M16/dat` | `data-paul:/vendor-paul:ro` |

For example, after generating overlays under `/tmp/vtspeak-load-ext-overlays`:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -v "$PWD/data-paul:/vendor-paul:ro" \
  -v /tmp/vtspeak-load-ext-overlays/omit-gen-dat:/work/data-paul/M16/dat:ro \
  runtime /bin/bash /work/stage16/run-load-ext-resource-error-container.sh omit-gen-dat
```

`make-load-ext-empty-files.py OUTPUT` creates zero-byte replacements for
`dblist.idx` and `merged-gen.dat`; mount either file over its original path to
reproduce the empty-but-present cases. Logs contain only loader status, slot
occupancy, and database-size results.

The error-2 probe statically imports `VT_LOADTTS_EXT_ENG` so the export's
address is stable before GDB installs breakpoints. Build its import library and
PE32 harness with the builder:

```sh
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage16:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-dlltool -d /src/tools/revkit/work/stage16/vt_pau_error2.def \
  -l /out/libvtpau_error2.a
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage16:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage16/probe-load-ext-error2-static.c \
  /out/libvtpau_error2.a -o /out/probe-load-ext-error2-static.exe
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-load-ext-error2-static-container.sh control
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-load-ext-error2-static-container.sh force-allocation
```

The `control` mode returns AX 0. `force-allocation` changes the shared
0x2042c-byte allocation return in GDB from its observed non-null pointer to
NULL; `force-speaker` similarly forces the separate 0x4d18-byte per-speaker
allocation return to NULL. Both return AX 2, confirming each error branch by
forcing its post-allocation null condition. They do not simulate actual system
memory exhaustion. The GDB transcripts are
`load-ext-error2-static-v3-force-allocation-api.log` and
`load-ext-error2-static-v5-force-speaker-api.log`.

`probe-load-ext-james-license.c` reads the same verification record but never
prints its contents. It loads James slot 4 once with argument 6 as a file path
and once with argument 7 as a memory buffer, then records loader/checker status,
database size, license gate, and dictionary capacity. Build and run it with:

```sh
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage16:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage16/probe-load-ext-james-license.c \
  -o /out/probe-load-ext-james-license.exe
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage16/run-load-ext-james-license-container.sh
```

Both forms returned loader AX `0`, checker `0`, gate `1`, and capacity `6`.
The record's `dbsize=300` permits James' 247,848-KiB database and rejects
Paul's 496,212-KiB database. Static consumers use the capacity as the
per-speaker user-dictionary index limit, so James' state accepts indexes 0–5.
The observed capacity `6` matches the record's `channel=6`; that field
relationship is inferred from one successful record.

`run-export-queries.sh` queries the six speaker slots through the exported
speaker-name, speaker-info, and database-size helpers, plus the DLL version
string, after Paul loads.
`run-helper-exports.sh` queries path-key and default dictionary-name helpers,
tests quoted-comma, unterminated-quote, empty-field, and embedded-LF CSV cases,
serialization/classification, and the sync-info allocation/initialization/
copy/free lifecycle. It also checks selected dictionary-validation and
configuration helpers and reads exported data values. Its license queries
retain only result lengths/codes, including
selectors 1 and 2, not license-field contents. `run-db-unloaded.sh`
also records the heap-start data export before and after unloading speaker 1.
`run-csv-capacity.sh` calls `VT_CsvParser_MakeCsv_ENG` with every capacity from
0 through 64 for fields `A` and `b,c`, then checks embedded-quote serialization
with `a"b` and `plain`, plus the null-field-array/zero-field case. It uses a
64-byte output region, captures bytes and guard positions, and restores the
Stage 5 fixtures from temporary backups. `run-csv-capacity-large.sh` repeats
the same short serialization at capacities 65, 128, 1,024, 4,096, and 65,536
with a correspondingly sized guarded allocation.
`run-csv-capacity-xlarge.sh` repeats the short serialization at capacities
131,072, 1,048,576, 4,194,304, 16,777,216, and 67,108,864 bytes with a
matched allocation and prefix/post guards.
`run-csv-negative-count.sh` tests field counts -1 and `INT_MIN` with a null
field-array and an eight-byte output region; it records raw AX, output bytes,
and the prefix guard, then restores the Stage 5 fixtures.
`run-csv-pointer-edges.sh` runs three separate processes for null output,
null field-array/count 1, and null field-element/count 1. It records each
process exit and restores the Stage 5 fixtures.
`run-csv-parser-edges.sh` probes doubled and embedded quotes, trailing
delimiters, CRLF and embedded LF, multiple input lines, and the same quoted
row with parse flags 0, 1, and 2. `run-csv-flag-matrix.sh` tests `INT_MIN`,
-1, 0, 1, 2, 3, 255, and `INT_MAX`, recording parsed fields and whether the
caller string changed. Only flag 1 parses in place. `run-csv-encoding-edges.sh`
checks CP1252 and UTF-8 high-byte sequences around ASCII delimiters and records
the returned field bytes in hex.
`run-csv-parser-deep.sh [RUN_ID]` generates bytewise GDB inputs for 41 quote,
empty-field, and CR/LF cases. It records every returned field as hex, so control
bytes do not distort the evidence. All 41 calls returned low AX 1. A run ID
keeps each batch's log separate; `csv-parser-deep-4-api.log` is the latest
capture.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-csv-flag-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-csv-encoding-edges.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-csv-parser-deep.sh next
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-csv-capacity-large.sh
```

`run-syncinfo-fields.sh` checks initialized header/row values and allocation
pointers, then calls copy at native dimensions (600 rows × 65 entries). Its
full comparison covers the 11 header dwords at `+0xc..+34`, 4,800 row values,
78,000 nested values, and all 600 destination nested pointers. The two
dimension words at `+4/+8` are checked by the shape probe. Static pseudocode
copies all 13 scalar header dwords at `+4..+34` and retains the destination's
row-array pointer at `+0`. It restores the allocator
dimensions before freeing both objects and snapshots/restores the Stage 5
fixtures.
`run-syncinfo-copy-edges.sh` calls copy with null arguments, performs a
self-copy using sentinels, and tests 2×3, 0×3, and 2×0 source shapes against
native-sized physical allocations. It records source-driven destination
dimensions, copied values, untouched tail sentinels, and pointer separation;
it restores native dimensions before freeing and also restores the Stage 5
fixtures.
`run-syncinfo-copy-shapes.sh` checks a 2×3 source against destination metadata
1×1, plus source dimensions -1×3 and 2×-1. It confirms the helper overwrites
destination dimension words from the source, uses signed-positive loop
conditions, and still copies fixed row fields when width is negative. Both
objects retain native physical allocations; dimensions are reset before free.
`run-syncinfo-copy-undersized.sh` runs two separate disposable processes with
guard pages: a 2×3 source against one physical destination row, and a 1×2
source against one physical nested entry. The copy faults on row 1 and nested
entry 1 after writing the earlier data. This confirms the source-driven loops
overrun physical destination storage without a capacity check. Replay with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-syncinfo-copy-undersized.sh
```

`run-syncinfo-runtime-fields.sh 1` snapshots the 14 header dwords, every
populated row field, and all nested entries through the active data polls of
`VT_TextToBufferEX_ENG`. Its safe terminal path does not dereference the
borrowed SyncInfo pointer after result 1 frees it, and it kills the inferior
before fixture restoration. The current capture is
`syncinfo-runtime-fields-safe2-1-api.log`; it confirms row +8 and nested first
dwords are frame counts, source intervals at row +12/+16, and the header's
complete start/end cursors across global frame, row, and nested-entry
coordinates. The pseudocode for `FUN_10021640`, `FUN_100216d0`, and
`FUN_100217e0` is retained in
`../reports/stage16-syncinfo-cursor-functions.txt`. Producer pseudocode for
the state and synthesis-record sources of row offsets 20–32 is retained in
`../reports/stage16-syncinfo-row-producers.txt`. It also records the full
`DAT_1007daa8` raw-class-byte to nested-selector mapping. Row `+28` is the
active parser-output record count; row `+32` copies synthesis-record `+0x28`,
not the separate loop ordinal at `+0x10`. The lookup maps raw bytes `0x5a–0x5e`
to the consumer's silence selector `0x28`. Semantic coordinate names for
row `+20/+24` are source-buffer byte offsets, as confirmed by the report
formatter's use of them to address original text and by repeated phrase
spans in the safe capture. The row-start predicate is mapped to adjacent
group/index (`+0x28`) and kind (`+0x27`) changes; ordinary selectors index the
phone mnemonic table (`AA` through `ZH`), and `0x28` is formatted as silence.
Still open are upstream schema names/domains for the group, kind, and class
values, and whether any separate downstream tool reads the generated report.
Physically undersized SyncInfo copies are confirmed to fault on the first
out-of-bounds write after partial updates. The row-slot tag at record `+0x2e`
is now mapped: the producer assigns the active circular
SyncInfo row index (0–599), with 599 as the wrapped previous slot when the
cursor is zero, and uses changes in that tag to avoid redundant source-span
refreshes. See the Lead 6 behavior report and
`../reports/stage16-syncinfo-source-coordinate-consumers.txt`.
The controlled normal-application trace records 41 kind-2 synthesis records,
their repeated `+0x28` indices, raw `+0x32` class bytes, and normalized
selectors in `syncinfo-index-classes-records-api.log`. A corrected sequence
trace reads parser state through `*(engine+0x4c)`: it observes active counts 9
and 6 for the 41- and 25-record batches. The earlier trace's zero count came
from reading `engine+2`, the wrong base.
The follow-up neighbor-field capture records `+0x2e == 0`, `+0x2a == 0`,
and leading dwords `100,100,200` across those same 41 records. Both sequence
and neighbor-field captures stop at producer entry, before it assigns current
`+0x2e` row-slot tags; entry-time zeros therefore do not measure producer
refresh decisions. Its trace and runner
are `trace-syncinfo-record-neighbor-fields.gdb` and
`run-syncinfo-record-neighbor-fields.sh`, with output in
`syncinfo-record-neighbor-fields-api.log`. The corrected sequence capture is
`syncinfo-synthesis-state-sequence-api.log`, replayed by
`run-syncinfo-synthesis-state-sequence.sh`.
`probe-license-api.c` directly calls the public checker, info, comment, and
TTS-info exports against the read-only verification record by file path and
memory buffer. It compares all info selectors across both paths, tests the
accepted frame transform and selected malformed/mutated forms, and calls the
internal attribute helper on value-suppressed baseline/changed/missing-key
cases. It also checks license-dependent TTS-info requests with default,
supplied, and invalid file paths. The harness passes all six built-in
voice-table entries, then loads and unloads the default Paul slot to exercise
loaded database-size checks and their threshold. Its output contains return
codes, lengths, slot numbers, capacity boundaries, guard status, and equality
flags only. The
capacity sweep shows selectors 3–15 need 322 destination bytes even when the
selected attribute output is shorter; the comment export has the same 322-byte
minimum. The framed file is created under container `/tmp` and deleted after
the calls.

Build and run it from the repository root with the local PE32 builder and Wine
images available:

```sh
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage16:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp \
  /src/tools/revkit/work/stage16/probe-license-api.c \
  -o /out/probe-license-api.exe
bash tools/revkit/work/stage16/run-license-probes.sh
```

The capture is [`license-paths-api.log`](license-paths-api.log). The probe
does not print license values or mutate vendor inputs.
`run-unit-history-mode.sh` sets the unit-selection history flag to 0 or 1 at
the extended-load entry, before model loading, then captures one format-4
output and log. Run the two modes serially to compare equivalent fresh
processes:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-unit-history-mode.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-unit-history-mode.sh 1
```

`run-verify-tts.sh` samples `VT_VerifyTTS_ENG` with the current synthesis text
at the first file-synthesis entry, after the model is loaded. It covers loaded,
unloaded, and invalid speaker slots plus a null-text control, records the raw
return register, and lets the executable finish normally.
`run-verify-tts-matrix.sh` crosses dictionary indexes -2, -1, 0, 1, 1,023,
and 1,024 with text types -1, 0, 1, 4, 6, 7, and 255.
`run-verify-tts-types.sh` sweeps every byte text type from 0 through 255
at dictionary index -1. `run-verify-tts-text-edges.sh` checks a non-null empty
string and a lone `<` byte. The latter reaches low-AX error `-5`; the empty
case returns low AX `-3` while upper EAX bits are stale. These probes report
result registers/counts and restore the Stage 5 fixtures; the dictionaries are
empty in this baseline runtime.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-verify-tts-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-verify-tts-types.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-verify-tts-text-edges.sh
```

`run-userdict-errors.sh` checks user-dictionary limits, invalid indexes,
loading an existing WAVE as a dictionary, and unloading an empty dictionary
slot. It does not change that file.

`run-userdict-heap-lifecycle.sh` uses engine-allocated strings and proves
successful three-field `P` and `A` load/synthesize/unload calls. Its
`vtspeakprobe` sample produces the same WAV with or without the dictionaries.
`run-userdict-hello-world-state.sh` records gate byte 0 plus a non-null
dictionary-slot pointer and captures `hello` with the same `P` and `A` rows
used by the forced-gate run; all three WAV hashes match.
`run-userdict-license-state3.sh` captures only the numeric license-check return,
dictionary capacity, and gate state: `-11`, `1`, and `0`. The existing
license probes map `-11` to the final `dbsize` predicate and show the supplied
value 300 falls below the loaded Paul's 452-unit minimum.
`run-userdict-hello-world.sh` is an explicitly debugger-forced experiment: it
sets the per-speaker dictionary gate byte to 1 in the disposable Wine process,
then compares a `hello,HH,P` row and a `hello,world,A` row against the control.
Both loaded rows produce different WAV hashes. No vendor binary, voice data,
license file, or persistent engine state is modified. The observed gate value
and the controlled outputs are recorded in `userdict-hello-state-api.log`,
`userdict-hello-world-api.log`, and their hash manifests.

`run-userdict-validation-matrix.sh` probes CSV row shapes with engine-heap
allocated path arguments. It captures accepted three/four-field rows, quote
and CRLF samples, malformed field counts, empty fields, and type errors. A
failed parse is followed by an empty-slot unload check and a valid reload at
the same index. Run the complete matrix with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-validation-matrix.sh
```

To replay the isolated unsupported-phoneme case, pass its fixture ID:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-validation-matrix.sh invalid-target-example
```

Rows with `example` or `NOTAPHONE` as the `P` target terminate this build with
an access violation during the loader call. The matrix traces and exact
acceptance/rejection boundary are documented in
`docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md`.
Optional case IDs `unmatched-type-quote`, `embedded-source-quote`, and
`escaped-source-quote` probe an unclosed final-type quote, an embedded
unquoted-source quote, and a doubled-quote source escape. All load/unload
successfully. `run-userdict-malformed-quote-effect.sh` compares the first two
forms against the plain mapping and tests both `hello` and `he"llo`. The
unmatched-type WAV equals the plain row; both source-quote forms change only
the quote-containing test text to the plain row's hash.

The `run-userdict-targetphon-*` probes sweep all printable single bytes,
uppercase pairs, every uppercase triple with suffixes 0–2, all stress digits
0–9 for the recognized vowel stems, sequence separators, and the repeated
phone-token boundary. Marker boundary probes vary case, separators, and
trailing content; the in-place probe prints each caller buffer after the API
call. `run-userdict-targetphon-length-fresh.sh` reconstructs 130- and
131-token strings in separate allocations, including a case with a trailing
space, so earlier API mutations cannot influence later cases. The
`run-userdict-targetphon-converter-map.sh` records one emitted byte for each
accepted spelling and `#`;
`run-userdict-targetphon-converter-sequence-check.sh` checks that grouped
inputs preserve token order; `run-userdict-targetphon-converter-limit.sh`
checks the 65-output-token maximum and 66-token failure. The
accepted inventory, mutation behavior, and post-trim length boundary are
described in the Lead 6 behavior report. Each runner writes its own API
capture and refuses to overwrite prior evidence.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-char-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-pair-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-stress-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-stress-digit-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-marker-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-marker-boundaries.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-marker-inplace.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-sequence-errors.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-token-count-boundary.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-length-fresh.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-converter-map.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-converter-sequence-check.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-converter-limit.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-targetphon-order-check.sh
```

`run-userdict-extended-buffer.sh` calls `VT_LOAD_UserDict_EXT_ENG` directly
with a null filename and heap buffer, repeats the call against its occupied
slot, unloads it, then forces the initialized-state flag off to capture the
loader and unloader guard returns before restoring the flag. The separate
`run-userdict-extended-argument-matrix.sh` varies filename/buffer presence and
buffer lengths `0`, `1`, `9`, `10`, `11`, and `-1`. It tests both file and
memory input in one process, with each call using a fresh index.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-extended-buffer.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-extended-argument-matrix.sh
```

`run-userdict-fourth-field-effect.sh` compares `hello,HH,P` with nonempty and
empty fourth columns, and with "hello","HH","P". It temporarily enables the
per-speaker user-dictionary gate, synthesizes `hello` for each, then restores
the gate. All four WAVs are byte-identical; the quoted row behaves like the
unquoted one.
`run-userdict-empty-source-effect.sh` loads `,HH,P` with the gate enabled and
compares `hello` and `.` against no-dictionary controls. Both pairs are
byte-identical; broader empty-source matching remains uncharacterized.
`run-userdict-inuse-unload.sh` injects a scratch context pointer into the
loaded speaker's user-dictionary reference slot and matches its
`+0x1312c0` dictionary field to a loaded dictionary. This drives the unload
in-use return, then restores the slot, frees the scratch context, and confirms
idle unload. It tests the guard with a controlled structure; it does not capture
a naturally active synthesis context.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-fourth-field-effect.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-empty-source-effect.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-malformed-quote-effect.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-inuse-unload.sh
```

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-heap-lifecycle.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-hello-world-state.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-license-state3.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-hello-world.sh
```

`run-buffer-capacity.sh` passes several `output_len` values with a physically
oversized output allocation and boundary sentinels; it captures only the
API-reported output bytes.
`run-texttypes.sh` sends the same ASCII text with the default and each declared
text-format value, capturing one selector-4 WAVE for each call.
`run-markup.sh` sends a VTML `<vtml_sub>` substitution plus its child-text
control through all declared text-format values.
`run-buffer-stream.sh` synthesizes a long repeated phrase to drive the buffer
API through flag-0 start, flag-1 chunk drains, and final completion, comparing
the concatenated bytes with file selector 0 before and after the buffer call.
`run-buffer-stream-state.sh` checks busy, cancel, and poll-after-cancel cases.
`run-buffer-ex.sh` probes one `VT_TextToBufferEX_ENG` selector (0, 1, or 2) in
an isolated process using a guarded 60,000-byte buffer. It records the initial
return/output length, performs flag-1 polls while the result is 0, and saves
each nonempty chunk. Its dedicated fixture backups avoid replacing the shared
Stage 16 baseline snapshots. Run selectors serially:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex.sh 1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex.sh 2
```
`run-buffer-ex-deep.sh` repeats the three selectors with a 78-byte utterance
to force multiple chunks. It reports each return/length pair, watches the
optional pointer outputs with guard bytes, then checks busy, flag-2
cancellation, poll-after-cancel, null/empty text, null buffer, nonzero
thread-ID, and database-unloaded cases in the same process. It records only
runtime results; captured pointer values are process-local engine addresses.
Run selectors serially:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-deep.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-deep.sh 1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-deep.sh 2
```
`run-buffer-ex-edges.sh` calls each selector first with a negative flag to
observe the fixed maximum-size query, then starts a long call with
`output_len = 1` and a 60,000-byte physical region. It records boundary
samples and snapshots the second returned descriptor's count, array pointer,
and selected fields from up to 64 live records. The captures show the input
length is not a capacity. `run-buffer-ex-guard-page.sh SELECTOR` goes further:
each selector receives exactly one writable byte followed by a no-access
page, with `output_len = 1`. All three fault writing into the protected page
at `0x10064000` (`rep movs`); Wine reports a write fault at the guard address.
This confirms a physically undersized target cannot be used safely and the
length argument does not cap the transfer. The runner restores the Stage 5
fixtures on exit.
The first selector-2 run returned 8,984 bytes and two repeats returned 30,000;
the shorter capture is preserved as an excerpt in
`buffer-ex-edges-2-first-capture.txt`, with full 30,000-byte reruns beside it.
Run selectors serially:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-edges.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-edges.sh 1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-edges.sh 2
```
`run-buffer-ex-guard-page.sh SELECTOR` places exactly one writable output
byte before a protected page. `run-buffer-ex-guard-size.sh SELECTOR BYTES`
repeats with a requested writable span from 1 to 4,096 bytes; the 16-byte
selector-0 case and 4,096-byte cases for all selectors fault on the next
protected page during the copy helper's `rep movs`, with incoming
`*output_len = 1`. Since the API has no capacity parameter, allocate the
fixed selector maximum (60,000 or 30,000 bytes) before calling it.
`run-buffer-ex-invalid-extra.sh a|b` separately
passes address 1 as
one optional output pointer while leaving the other null, and records the
resulting access violation. Run the guard-page selectors and invalid-pointer
slots serially:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-page.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-page.sh 1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-page.sh 2
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-size.sh 0 16
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-size.sh 0 4096
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-size.sh 1 4096
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-guard-size.sh 2 4096
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-invalid-extra.sh a
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-invalid-extra.sh b
```

`run-buffer-ex-records.sh` snapshots the descriptor plus up to 32 rows. The
`<vtml_sub>` fixture remains a useful negative control with count 0. A
`<vtml_mark name="start"/>...<vtml_mark name="end"/>` fixture produces two
kind-1 rows under selectors 0, 1, and 2. The boundary fixture confirms the
512-byte name takes the kind-3/truncation path. `boundary-sweep` varies lengths
1, 255, 256, 510, 511, and 512 in a single utterance and exposes unexplained
equal nonzero values in row offsets `+0` and `+4` for every row after the
first. A named-versus-unnamed contrast resolves the kind-2 trigger:
`<vtml_mark/>` yields a kind-2 row with an empty payload, and
`<vtml_mark name="named-control"/>` yields a kind-1 row carrying that name.
Selectors 0–2 return the same two-row structure and frame coordinates
(`0x90c`, `0x15e7`); their byte lengths differ by the expected PCM16 versus
one-byte encoding factor. An edge fixture shows `name=""` yields no row,
`name=" "` retains a one-space kind-1 payload, and uppercase `NAME="Upper"`
is accepted as kind 1. See `buffer-ex-tag-family-sweep-mark-forms-*.log` and
`buffer-ex-tag-family-sweep-mark-name-edges-0.log`.

A separate sweep put `<vtml_break>`, `<vtml_pause>`, `<vtml_partofsp>`,
`<vtml_phoneme>`, and `<vtml_sayas>` before a named-mark control. Only the
control produced a record. Legacy internal spellings produced no rows in the
same probe. These captures narrow the record writer to mark events for the
tested forms; they do not establish how an external caller consumes the rows.
Reproduce the family, legacy-spelling, unnamed-mark, and name-edge probes with
`run-buffer-ex-tag-family-sweep.sh SELECTOR [vtml|legacy|mark-forms|mark-name-edges]`.

A controlled position probe holds a 511-byte marker name constant and
moves its mark after 0, 1, and 2 ASCII characters: `+0/+4` are equal and
change from 0 to `0xefa` to `0x23a4`, while `+8` changes from 0 to 1 to 2.
This ties the fields to source position as well as synthesis options. Enhanced
GDB traces also snapshot the returned SyncInfo rows. At positions 1, 2, and 5,
descriptor `+0/+4` are `0xefa`, `0x23a4`, and `0x5d54`, exactly matching
SyncInfo row 0 `+8`; that row's source-text end is 0, 1, and 4, respectively.
The descriptor's own `+8` remains the marker-relative position 1, 2, or 5.
Static tracing of `FUN_100217e0` explains the short-position match: it selects
SyncInfo rows by source-text span, accumulates their `+8` values into
descriptor `+4`, and derives `+0` by adding a context-owned base. A long-mark
probe follows every poll with a marker after 51 source bytes. On poll 2,
descriptor `+0/+4` are `0x19fac/0x401c`, while SyncInfo header `+0x10/+0x18`
are `0x15f90/0x1065`; `+0` is exactly `+4 + +0x10`. The same frame base
advances by 30,000 on each full poll under all selectors despite their
different byte lengths. The poll-2 header gives start index 10, partial-row
offset `+0x18 = 0x1065`, and marker position `0x33`; SyncInfo rows 10–11 have
`+8 = 0x16b1` and `0x39d0`. The equation `0x16b1 + 0x39d0 - 0x1065 = 0x401c`
reproduces descriptor `+4`; adding base `+0x10 = 0x15f90` gives `+0 = 0x19fac`.
This directly establishes `+4` as chunk-local frames and `+0` as utterance-
global frames, including the nonzero partial-row subtraction.
The unit is one mono audio frame: at positions 1, 2, and 5, the two SyncInfo
row `+8` values sum to 9,191, 14,281, and 29,132. Selector-0 output lengths
are exactly twice those counts because selector 0 returns signed PCM16; at
position 2, selectors 1 and 2 each return 14,281 bytes of one-byte
A-law/μ-law. Descriptor `+4` subtracts header `+0x18` from the first selected
row; descriptor `+0` adds header `+0x10`. The external consumer's convention
remains open. Reproduce the nine-poll trace with
`run-buffer-ex-chunk-sync.sh 0|1|2`; the full-row selector-0 capture is
`buffer-ex-chunk-sync-full-0.log`. A ten-character prefix fills the 60,000-byte
output request and is not used to infer a value.

The 16-byte descriptor header's `+8/+0x0c` words are first/last selected row
indices from the timeline mapper, or `-1/-1` when no row matches: observed
examples are empty `(-1,-1)`, one record `(0,0)`, three records `(0,1)`, and
the six-record boundary sweep `(0,5)`. A two-marker sample also reports
`(0,0)`, so this is a selected subset rather than an automatic full-array
range. Each 528-byte row has three 32-bit leading fields, 512 inline bytes, a
kind byte at `+0x20c`, and three trailing bytes. The name branches of
`FUN_1002efd0` populate the inline bytes using `name` or `NAME`; kind 1 is a
nonempty named mark, kind 2 is an unnamed mark with an empty payload, and kind
3 is an overlength named mark whose payload is truncated to 511 bytes. An
explicitly empty `name=""` emitted no row in the edge probe.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-records.sh 0
```

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 0 mark
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 1 mark
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 2 mark
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 0 boundary-sweep
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 0 position-0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 0 position-1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 0 position-2
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 1 position-2
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 2 position-2
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark.sh 0 position-5
```
`run-buffer-ex-options.sh` varies the scalar arguments on plain text;
`run-buffer-ex-options-punctuation.sh` repeats that matrix on
`Hello, world.`. On both, pause 0/250/65535, text types 0–7, and dictionary
1023 match the baseline byte-for-byte, while pitch, speed, and volume alter
the captured output. These captures do not exhaust option cross-products or
the effects of alternate populated dictionaries and broader VTML/options
combinations.
`run-buffer-ex-mark-options.sh SELECTOR` tests representative options against
the three-row marker boundary fixture and snapshots each row's `+0/+4/+8`
fields and kind. On the 511-byte kind-1 row, `+0/+4` are equal and respond to
pitch/speed (default `0x59d`, pitch 50 `0x68d`, speed 50 `0xbda`), while
remaining equal across selectors 0–2 and unchanged by volume, pause, dictionary
1023, or text type 4. The simple first marker and 512-byte kind-3 row retain
zero `+0/+4`. The frame unit and internal bases are established for the
controlled positions; the external consumer's interpretation remains open.

`run-buffer-ex-pause-compare.sh SELECTOR CASE` drains identical text with an
inline `<vtml_pause time="N"/>` at `N=0`, `200`, and `1000`, then compares its
total bytes with the plain control. Across selectors 0–2, the added pause
frames are exactly 3,200 for 200 ms and 16,000 for 1,000 ms. Selectors 1/2
return 17,762/20,962/33,762 bytes for 0/200/1000 ms; selector 0 returns twice
those counts as PCM16. This confirms that VTML pause `time` is milliseconds at
16 kHz and that the chunked path reports the complete output after draining.
The separate scalar pause option at argument 13 is measured below. Logs are
`buffer-ex-pause-{control,pause0,pause,pause1000}-{0,1,2}.log`.

Scalar argument 13 was measured separately from the VTML attribute. On
`Hello. World.`, 120 and 250 add exactly 1,920 and 4,000 frames to the
17,762-frame zero setting, establishing milliseconds at 16 kHz in this
sentence-boundary case. The 65,535 setting completed a selector-1 drain at
1,066,322 bytes, which is the baseline plus 65,535 ms of mono frames. The
comma-separated `Hello, world.` and the inline-pause fixture are unchanged at
0, 120, 250, and 65,535. Thus the scalar duration applies in the tested
sentence boundary and does not stack onto `<vtml_pause time="200"/>` or the
tested comma. Placement and behavior for other text categories remain open.
The tag fixture matrix uses `run-buffer-ex-scalar-pause.sh SELECTOR`; sentence
and comma fixtures use `run-buffer-ex-scalar-pause-context.sh SELECTOR CASE`
with `CASE` set to `sentence` or `comma`. To reproduce the full maximum pause drain, use
`run-buffer-ex-scalar-pause-max.sh 1`. Logs are
`buffer-ex-scalar-pause-{0,1,2}.log`,
`buffer-ex-scalar-pause-{sentence,comma}-{0,1,2}.log`, and
`buffer-ex-scalar-pause-max-1.log`.

To test natural oversized-unit carry-over, `run-buffer-ex-long-token.sh 1`
drains unbroken `a` strings of 1,024 and 4,096 bytes. They completed in 4 and
13 calls at 84,789 and 343,768 bytes. `run-buffer-ex-long-token-v2.sh 1`
drains a 16,384-byte token in 47 calls at 1,377,407 bytes. All returned the
normal terminal result; none reached `-8`. These runs show that multi-chunk
output and long unbroken tokens are handled by the split/carry path at these
sizes, but do not establish a maximum token or carry-over size.

`run-buffer-ex-force-minus8-v3.sh SELECTOR` validates the defensive `-8`
return and recovery path. In the isolated process, its GDB trace saves the six
bytes of the cursor-limit conditional jump, replaces that jump with an
unconditional jump to the `-8` epilogue, calls the export, and restores the
original bytes immediately after the return. For selectors 0–2 the export
returns `-8`; flag-2 cancellation returns `1`, and the following flag-1 poll
returns `-2`, showing that cancellation cleared the active slot. This is an
injected branch test, not evidence that valid input naturally reaches the
comparison. The captured per-selector logs are
`buffer-ex-force-minus8-v3-{0,1,2}.log`.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-force-minus8-v3.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-force-minus8-v3.sh 1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-force-minus8-v3.sh 2
```

`run-buffer-ex-invalid-extra.sh a|b` passes address `1` as optional output
pointer 8 or 9, respectively. Both calls exit with Windows
`STATUS_ACCESS_VIOLATION` (`0xc0000005`) before returning from the injected API
call. Static disassembly locates the zero stores at `0x10021c45` and
`0x10021c52`, before synthesis begins. Null remains accepted; any invalid
non-null address faults. The captured logs include the actual pointer values.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-options.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-options-punctuation.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark-options.sh 0
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark-options.sh 1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-ex-mark-options.sh 2
```
`run-buffer-progress.sh` calls `VT_TextToPcmBuffer_ProgressBar_ENG` with a
null progress window/message, thread 0, speaker 1, and default option values.
It records the output length, raw EAX left by the wrapper, and guard bytes;
this short-text call does not exercise the long-text progress notification
path. The runner uses dedicated Stage 5 input/output snapshots.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-progress.sh
```
`run-preprocess-noop.sh` exercises the zero-byte-flag branch of
`VT_TextToPreprocessInfoFile_ENG` after model load, with every other argument
null/zero. Static disassembly shows this branch returns before inspecting
those pointers. The runner verifies that its dedicated Stage 5 fixture
snapshots were restored byte-for-byte.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-preprocess-noop.sh
```
`run-preprocess-nonzero.sh` is the original flag-1 probe with invalid `-1`
dictionary and text-type arguments; its low AX `-6` is superseded by the valid
argument matrix below. The wrapper is decompiled as `void`, so all captured
register values are observations rather than a declared return.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-preprocess-nonzero.sh
```

`run-preprocess-flags.sh` runs byte flags 1–10 and 255 in separate fresh
processes with loaded Stage 5 text, speaker 1, synthesis settings `-1`,
dictionary index 0, and text type 0. All calls reached common completion with
raw low AX 1. Flag 6 created fixed working-directory `test.pcm`; the 23,606-
byte result is preserved as `preprocess-followup-test-flag-6.pcm`. Every run
restores and compares dedicated Stage 5 input/output snapshots, and the
runner refuses to overwrite a preexisting `test.pcm`.

The corrected heap-backed runner
`run-preprocess-flag-heap-isolated.sh BYTE_FLAG` initializes the second
argument byte by byte in an engine-allocated buffer, verifies it before and
after the call, and captures new outputs without overwriting existing files.
For flags 1–5, 7–11, 254, and 255, each generated basename exactly matches
the supplied path. Flag 4 appends `.0` through `.3`; static code formats four
`%s.%d` names. Flag 6 ignores the supplied path and writes fixed `test.pcm`.
The earlier GDB-literal runners corrupted strings before API entry, so their
mangled filenames do not show an API transformation. Preserve those captures
as diagnostics; use the heap-backed logs and `heap-flagN-20260926` outputs for
path claims.

The heap-backed outputs classify the text-writing flags. Flag 3 writes one
period per source byte, followed by the count and decimal engine phone/control
bytes; flag 5 renders phone symbols as labels; flag 7 emits source text and
per-word pronunciation/boundary/span records; flag 10 re-emits text with a
space before tested terminal punctuation. The flag 7 `Hello world.` output is
127 bytes. Its integer is a mapped structural boundary-marker code, its
metadata tag letters are generated from a seven-bit mask, and its source
spans are inclusive for the tested ASCII text. The boundary-class names and
metadata-letter meanings remain unresolved. Five fresh runs of identical
`Hello world.` text produced 52, 140, 33, 151, and 33 leading periods. Static
code selects one of eight fixed period strings using a `GetTickCount`-seeded
generator; the string lengths are 108, 111, 140, 52, 33, 151, 106, and 58.
The reason for adding this time-selected filler remains unknown. Full samples
and interpretation are in the
[behavior report](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md#output-record-interpretation).

`run-preprocess-text-case.sh BYTE_FLAG CASE_ID TEXT_FILE` repeats selected
flags using one of the fixed `preprocess-case-*.txt` files. Flag 7 cases cover
one word, changed words, punctuation, word order, and homographs; flags 3, 5,
and 10 additionally compare `Hi.` and `Hello, world!`. The runner restores and
compares the Stage 5 snapshots and refuses to overwrite an existing capture.

`run-preprocess-flag-isolated.sh BYTE_FLAG` is the earlier convenience runner
and still demonstrates the call path, but its GDB literal path is not valid
evidence for filename behavior. Static dispatch sends values 12–253 through
the default mode.
`run-preprocess-input-errors.sh` calls the same nonzero flag with null and empty
text pointers; the captured raw low AX values are `-3` and `-4`, respectively.

`run-preprocess-default-flags-sweep.sh` calls every byte flag from 12 through
253 in sequence in one loaded process. It uses a unique output basename per
flag, records each raw EAX/low AX observation, and reports each file's length
and SHA-256. All 242 calls returned raw EAX `1` through the void wrapper, and
all 242 outputs were empty files. This closes the runtime inventory for those
flag values on the fixed Stage 5 text and settings, while not establishing
fresh-process or cross-text equivalence for every value.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-preprocess-default-flags-sweep.sh
```

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-preprocess-flags.sh
```
`run-lipsync-log.sh` calls `VT_TextToLipSyncLog_ENG` after load with the
current Stage 16 text, an empty second string, and remaining numeric arguments
set to -1. It repeats with speaker -1, null text, and empty text to capture
the fallback and text precondition branches. The wrapper is decompiled as
`void`, so the log records raw EAX/AX values rather than a declared return.
The runner restores and compares its dedicated fixture snapshots.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-log.sh
```
`run-lipsync-output.sh` used a GDB string literal as the path pointer; its
access-violation result is invalid because that setup did not provide a valid
target-process string. The bytewise heap-string probes below supersede it.
`run-lipsync-prefix.sh` repeats it with the simple string `vtspeak-lead6`.
That call returned raw EAX 1 and produced an ASCII report; its contents were
retained, but the runner did not preserve the original filename. Matched runs
with second strings `vtspeak-lipunits-100` and `vtspeak-lipunits-200` created
files at those exact basenames without extensions. The null-second-string call
returned raw EAX 1 and produced no report. On `Hello world.`, report totals
match PCM16 audio frames: speed 100 gives 11,803 frames and a 23,606-byte
buffer result; speed 200 gives 5,682 frames and an 11,364-byte result. The
report lists per-phone IDs/labels and integer values, word-index/text spans,
and aggregate word/phone lengths. Bytewise-initialized heap paths with an
extension, `./` prefix, relative subdirectory, and `Z:/work/...` absolute path
all returned raw EAX 1 and created a report at the supplied path. In
`FUN_1001df10`, the global at `0x1009f948` is an empty string; an exact match
returns a null logging context. The helper first formats a candidate name with
`length-sync-%s-%s.txt`, then compares the supplied path with that empty
global. An exact match returns null immediately; a nonempty path proceeds to
the caller-path branch. The later branch that selects the formatted candidate
also requires equality with the same empty global, so it is unreachable for a
stable argument through this public API call path. Runtime captures agree:
the empty-string sentinel call returns raw EAX 1 from the void wrapper and
creates no report, while nonempty paths produce reports at their supplied
paths. In the disassembly, the empty-path return is at `0x1001df82`–
`0x1001df95`; candidate selection is at `0x1001e005`, and the nonempty caller
path branch is at `0x1001e013`. The intended purpose of the candidate branch
remains unknown. The old absolute-path AV also came from invalid GDB string
setup. A bytewise heap-string call using
`no-such-lipsync-dir/report.txt` printed the intended path, then the inferior
exited with `0xc0000005` during the API call; GDB abandoned the injected call,
so no raw EAX was captured and no report was created. This is an observed
invalid-parent failure under Wine. A second run stopped at EIP `0x10025e4d`
inside the report writer's small formatting wrapper; the faulting instruction
is `mov eax, DWORD PTR [edx]`, before the call to `0x10064d92`. This locates
the access violation at the writer's first-argument dereference, but the
captured run did not preserve EDX, so whether that pointer was null and why it
was invalid remain unresolved. See
[`lipsync-filename-bytewise-missing-parent-stop.log`](lipsync-filename-bytewise-missing-parent-stop.log)
and the disassembly report.
A bytewise relative path containing `\` returned raw EAX 1 and produced a
report below the intended subdirectory, confirming that Windows-style
separators work for that case. An absolute `C:\windows\temp\...` path also
returned 1 and created its report under the Wine prefix's C: drive; the
absolute Z: host-drive case remains separately captured.

`run-lipsync-option.sh CASE` repeats the same fixture at speed 100. The
all-`-1` report totals 11,803 frames; pitch 50 and 200 produce 12,850 and
11,998. Types 0–7, volume 0, pause 250, and dictionary 0 are byte-identical to
the baseline report on this one plain-text fixture. Combining pitch 50 with
those three settings matches pitch 50 alone; adding type 4 or 6 makes the
report byte-identical to baseline, consistent with the static reset of earlier
options for those type values. The report captures timing, so unchanged report
bytes do not establish unchanged waveform audio. The complete matrix is in the
`lipsync-option-*.log` and `lipsync-option-*-report.txt` captures.

See `lipsync-filename-bytewise-*.log` and the matching `*-output.txt`
captures, `lipsync-sentinel.log`, and
`lipsync-speed-{100,200}.log`/`-report.txt`. Other path grammar, the intended
purpose of the formatted-candidate branch, write-error handling, and option
effects on other text or speakers remain open. An attempted breakpoint inside the helper during a
GDB-injected call caused GDB to abandon the call;
`lipsync-sentinel-open.log` is retained as an invalid probe and does not
establish file-open behavior.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-speed.sh 100
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-speed.sh 200
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-null-path.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-filename-case.sh subdir-relative-verified
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-filename-case.sh absolute-z
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-filename-case.sh missing-parent
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-filename-case.sh backslash-relative
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-filename-case.sh absolute-c-verified
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-sentinel.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-option.sh text-type-7
```

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-output.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-prefix.sh
```
`run-makeinfo-heap-path.sh` calls `VTDTTS_MakeInfo_ENG` with Stage 16 text,
speaker 1, remaining scalar arguments -1, and an engine-allocated second
argument initialized byte by byte. The DLL emits the exact prefix plus
`.bin.dtt` and `.asc.dtt`; both outputs byte-match the older pair. The earlier
GDB-literal path probe's mismatched basename and the no-output follow-up are
not valid path evidence because that argument setup corrupted the string.

The ASCII and binary outputs agree on ten records. Each record maps to a
selected Paul `unit-gen.idx` row: PCM position/coded size select a `.dat` span;
PM position/size select the combined UPM span in mode 0, the first side in
mode 1, or the second side in mode 2; PCM size matches the full decoded unit
or its indexed side span; and pitch endpoints match the selected doubled UPM
edge periods. The body is therefore a selected synthesis-timeline manifest
for this utterance. The binary header contains `03 04` after `VTDTTS BINARY`
and the four NUL-terminated bank names `merged-gen`, `merged-num`,
`merged-etc`, and `merged-alp`; the ASCII file shows `3`, `4`, and the same
names. Four matches the bank count; the meaning of three remains open.
Controlled text captures map `File Index` values 0, 1, and 2 to `merged-gen`,
`merged-num`, and `merged-etc`; index 3 has not been observed. Type 2 is the
detailed phone/unit record. Static synthesis handling shows that type 1 is
zero-filled into the PCM stream, while the DTT writer represents it as a
size-only row. A targeted matrix captured this type between OW1 and W for
`Hello, world.` (`Size=3200`) and `Hello. World.` / `Hello... World.`
(`Size=14800`); the plain-sentence control had no type-1 row. The synthesis
loop uses twice `Size` as the byte count, so this is a silent PCM16 sample
interval. Modes 0, 1, and 2 select complete, first-side, and
second-side unit spans. The exact physical meaning of mode-2 `Shift Size`
remains open, though the writer's subtraction is mapped in the behavior
report. Rate captures tie `Pitch_rate` to pitch, `Volume_rate` to volume, and
`Duration_rate` to integer inverse-normalized speed; pause did not affect the
three fields. Local host executables and binary searches revealed no DTT
reader, so the downstream consumer is unidentified. The heap-backed runner
restores its dedicated Stage 5 fixture snapshots and refuses to overwrite any
capture.

`run-makeinfo-rate-matrix.sh` varies pitch, speed, volume, and pause while
holding text fixed. `run-makeinfo-text-banks.sh` emits numeric,
spelled-letter, and mixed text captures used to map bank indices. Both use
heap-backed strings and preserve the Stage 5 fixtures. Their captures record
controlled examples; they do not establish all option/error behavior or use
of `merged-alp`. `run-makeinfo-silence-cases.sh` compares plain, comma,
sentence, ellipsis, and pause-120 cases to observe when type-1 rows appear.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-makeinfo.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-makeinfo-heap-path.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-makeinfo-rate-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-makeinfo-text-banks.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-makeinfo-silence-cases.sh
```

`run-heap-lifecycle.sh` reads `VT_gHeapStartAddress_ENG` at DLL load, after
model load, and after an unload attempt. The new pre-load reading is zero;
earlier helper and unload traces also read zero after load and unload. The
export's purpose remains unknown. The runner snapshots and restores Stage 5
fixtures.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-heap-lifecycle.sh
```

`run-userdict-valid.sh`, `run-userdict-alpha.sh`, and
`run-userdict-candidates.sh` used GDB literal strings and their `-3` results
are not accepted as row-rejection evidence. The replacement heap-backed matrix
establishes the tested row-shape returns and the unsupported-phoneme crash
boundary. The earlier heap-backed probes establish dictionary-dependent WAV
changes when the per-speaker gate is enabled. The ordinary Paul process records
gate byte 0; the debugger-only gate-1 experiment is separately labeled and
occurs only in the temporary Wine process. The runtime also opened the default
relative file `../data-common/userdict/userdict_eng.csv`; that checked-in file
is zero bytes and is not a known-good sample. `run-userdict-normalizers.sh`
still records the private source normalizer, target normalizer, phoneme
validator, and converter observations.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-normalizers.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-candidates.sh
```

Rerun these follow-up probes serially from the repository root with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-config-fields.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-db-unloaded.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer-stream.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-markup.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-play-waveout.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-userdict-errors.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-export-queries.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-helper-exports.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-csv-capacity.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-csv-parser-edges.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-syncinfo-fields.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-syncinfo-copy-shapes.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-syncinfo-runtime-fields.sh 1
```

Runners that call through the Stage 5 executable back up and restore its input
and output fixtures on exit; run them serially because concurrent runners can
overwrite the shared fixtures. Each selector result is saved as
`file-format-N.out`; GDB logs record
the API entry and signed-short return. The error runner mutates only the
in-process API arguments or first text byte. The captured matrix and its
limits are documented in
[`Lead 6 API behavior`](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md).

`run-file-repeat.sh` issues two identical selector 4 calls through the DLL in
one loaded process, saves each result, and logs both return values. This
single repeat control is not a general persistence guarantee.

The bounded local API behavior, including chunked buffer processing, one
markup substitution, file/buffer mixed calls, and playback through a
discarding null sink, is documented in the [Lead 6 report](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md).
Nonzero thread IDs, other markup tags, audible playback, untested call orders,
and other voice packages remain outside that evidence.
