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
tests a CSV parse/field/serialization/classification case and sync-info
allocation/initialization/copy/free lifecycle, checks selected
dictionary-validation and configuration helpers, and reads exported data
values. Its license queries retain only result lengths/codes, not
host/license-specific output. `run-db-unloaded.sh` also records the heap-start
data export before and after unloading speaker 1.
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
`run-preprocess-nonzero.sh` calls the same export with byte flag 1, loaded
Stage 16 text, speaker 1, and a fresh path-like second argument. The observed
low AX was `-6`, and no file appeared at that argument; because the wrapper is
decompiled as `void`, this register value is not a declared C return. The
successful nonzero-flag behavior and second-argument role remain unresolved.
It restores and compares dedicated Stage 5 fixture snapshots.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-preprocess-nonzero.sh
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

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-lipsync-output.sh
```
`run-makeinfo.sh` calls `VTDTTS_MakeInfo_ENG` with Stage 16 text, speaker 1,
remaining scalar arguments -1, and a unique second-argument prefix. It records
the raw EAX left by the decompiled `void` wrapper and preserves the emitted
`.bin.dtt`/`.asc.dtt` files. The observed pair has `VTDTTS BINARY` and
`VTDTTS ASCII` headers; the ASCII file lists per-phone metadata. The observed
basename did not match the supplied prefix directly, so filename mapping and
field interpretation remain open. The runner restores its dedicated Stage 5
fixture snapshots and refuses to overwrite the observed output pair.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-makeinfo.sh
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

`run-userdict-valid.sh` tries a plain `P` row after checking the isolated
target validators; `run-userdict-alpha.sh` tries a normalized `A` row. Both
expose low AX `-3` from the void wrapper. `run-userdict-normalizers.sh` calls
the private source normalizer, target normalizer, phoneme validator, and
converter directly for the same sample values. `run-userdict-candidates.sh`
compares plain, quoted, and four-field `P` rows for `hello`; all three expose
low AX `-3`. An earlier quoted `vtspeakprobe` load attempt in
`run-userdict-debug.sh` was interrupted by nested parser breakpoints and has
no result. The successful row format, loaded-dictionary lifecycle, and
synthesis effect remain open.

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
