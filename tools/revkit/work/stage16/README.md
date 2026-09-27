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
Stage 5 fixtures from temporary backups.
`run-csv-parser-edges.sh` probes doubled and embedded quotes, trailing
delimiters, CRLF and embedded LF, multiple input lines, and the same quoted
row with parse flags 0, 1, and 2.
`run-syncinfo-fields.sh` checks initialized header/row values and allocation
pointers, then calls copy at native dimensions (600 rows × 65 entries). Its
full comparison covers 11 header fields, 4,800 row values, 78,000 nested
values, and all 600 destination nested pointers. It restores the allocator
dimensions before freeing both objects and snapshots/restores the Stage 5
fixtures.
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
return register, and lets the executable finish normally; empty-text and other
argument cases are not covered.

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
evidence for filename behavior. Values 12–253 were not individually run, but
static dispatch sends them through the same default mode as 10, 11, 254, and
255.
`run-preprocess-input-errors.sh` calls the same nonzero flag with null and empty
text pointers; the captured raw low AX values are `-3` and `-4`, respectively.

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
`run-lipsync-output.sh` repeats the call with a fresh nonempty path-like
second argument. In this GDB-injected probe the process terminated with an
access violation before the call result could be captured. The runner restores
its dedicated Stage 5 snapshots and refuses to overwrite the candidate output
path; this single failure does not establish the API's normal behavior.
`run-lipsync-prefix.sh` repeats it with the simple string `vtspeak-lead6`.
That call returned raw EAX 1 and produced an ASCII length report under the
`length-sync-*` naming pattern. The runner moves the generated file into this
stage's capture area and refuses to overwrite a preexisting matching output.
The report lists per-phone IDs/labels and integer values, word-index/text
spans, and aggregate word/phone lengths. The values' units and the exact
filename components remain open.

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
are not accepted as row-rejection evidence. Heap-backed probes now establish
successful three-field `P` and `A` rows, load and unload statuses of `1`, and
dictionary-dependent WAV changes when the per-speaker gate is enabled. The
ordinary Paul process records gate byte 0; the debugger-only gate-1 experiment
is separately labeled and occurs only in the temporary Wine process. The
runtime also opened the default relative file
`../data-common/userdict/userdict_eng.csv`; that checked-in file is zero bytes
and is not a known-good sample. `run-userdict-normalizers.sh` still records
the private source normalizer, target normalizer, phoneme validator, and
converter observations.

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
