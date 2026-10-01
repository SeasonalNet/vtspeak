# Lead 6: file, buffer, configuration, playback, and information API behavior (2026-09-25)

**Status: bounded Lead 6 API lifecycle recovered for the local 2013 Paul M16
DLL.** The evidence covers file selectors, synchronous and chunked buffer
output, selected errors, configuration and information queries, repeated and
mixed file/buffer calls, one VTML substitution tag, and playback through an
ALSA null sink. It does not establish worker-thread behavior for nonzero
thread IDs, real audible playback, the full markup grammar, every error path,
or other packages. The [API coverage map](lead6-api-coverage-map-2026-09-25.md)
also inventories all 66 named PE exports and distinguishes header-declared
calls from export-only helpers.

## Scope and method

The runtime target was the local `vt_pau.dll` reached through
`voicetext_paul.exe` in the isolated Wine container. The executable loads the
local Paul M16 model, then calls `VT_TextToFile_ENG`. A GDB breakpoint at
`0x1001da50` changed only the selector argument for each run. The input was
`Hello world.`; all other API arguments came from the executable's usual
call. Selectors 0–10 were captured, covering the header's declared values
0–5 and 7–9, its unsupported value 6, and one out-of-range value 10.

The error probes used selector 4 and changed only the in-process text or path
argument: null text, an empty string made by clearing its first byte, and a
null output path. The caller declares a `short` return, so the logs print the
signed low 16 bits of EAX.

At the same loaded-process breakpoint, the buffer probe directly called
`VT_TextToBuffer_ENG` with a 1 MiB target-process output buffer and initialized
`*output_len` to 1 MiB. It checked formats 0–3 at `flag=0`, `nThreadID=0`,
then format 0 with `flag=1` and/or `nThreadID=1`. The GDB runner saves bytes
only when the API returns success.

The runner backs up `tools/revkit/work/stage5/input1.txt` and `output.wav`
before each run and restores both on exit. The saved and restored SHA-256
values match: `input1.txt` `023ba700365530190e180c5cb26a10288ad30596ca77a5a893d193c5e4b7baa6`;
`output.wav` `3bd8bbfde92f1d645715de40a03a6a68e3acedf0a3361158b869b8a07d803187`.
The Compose service mounts vendor inputs read-only. The trace template,
runners, logs, and selector output captures are under
[`stage16`](../../tools/revkit/work/stage16/README.md).

## Selector results

The `VT_TextToFile_ENG` return was `1` for every valid selector and `-1` for
6 and 10. The latter two left the pre-existing output file unchanged. The
header marks selector 6 (IMA WAVE) unsupported.

| Selector | Header name | Captured output | Bytes | SHA-256 |
| ---: | --- | --- | ---: | --- |
| 0 | S16PCM | Raw signed 16-bit PCM; little-endian samples | 23,606 | `e232e2454fc12176c7cc21b6573bb83b436d96dbf2ec2f33de2164fee0f61582` |
| 1 | ALAW | Raw A-law bytes | 11,803 | `9880164f3368865789c6aa51cb710063a2c50693bb422c7d1652a3ea479d5dfc` |
| 2 | MULAW | Raw μ-law bytes | 11,803 | `424afda0592dfe2f1b94ba8e29e072b3de93737a80bef3e14372087512ce53ab` |
| 3 | DADPCM | Raw DADPCM bytes | 5,901 | `abe150030a647354515e4589cf651f6a05748e19e37f2c313efb256bc94d06c8` |
| 4 | S16PCM_WAVE | RIFF/WAVE, PCM, mono, 16 kHz, 16-bit | 23,650 | `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69` |
| 5 | U08PCM_WAVE | RIFF/WAVE, PCM, mono, 16 kHz, 8-bit | 11,847 | `057cbcbec053734397bd608160c5f8f3905e8997320d33beb39f6bddb23550ea` |
| 6 | IMA_WAVE (unsupported) | Rejected; prior file unchanged | 89,478 | `3bd8bbfde92f1d645715de40a03a6a68e3acedf0a3361158b869b8a07d803187` |
| 7 | ALAW_WAVE | RIFF/WAVE, A-law, mono, 16 kHz | 11,861 | `2f6fc20155b5450f77693f1bc22cf7ce1d7244845498568d70dfe85787ac4128` |
| 8 | MULAW_WAVE | RIFF/WAVE, μ-law, mono, 16 kHz | 11,861 | `9f0d5d31726d61c53e535edb15665c423dac422ee761c1b187d9bf21119390c4` |
| 9 | MULAW_AU | Four zero bytes followed by exactly selector 2's 11,803 bytes | 11,807 | `9d60ebfb2e5a77c5084485ceb7951cb0b81178aeb90e5ac6df9e5777603f012a` |
| 10 | Out of range | Rejected; prior file unchanged | 89,478 | `3bd8bbfde92f1d645715de40a03a6a68e3acedf0a3361158b869b8a07d803187` |

Selector 9's observed prefix is four zero bytes; the remaining bytes compare
equal to selector 2's raw output. This capture does not contain the usual
`.snd` magic and is recorded as observed rather than treated as a validated
AU file. Selector 4's file has a 44-byte WAVE header. Selectors 5, 7, and 8
produce WAVE headers with the corresponding PCM or G.711 encoding tags.

## Buffer results

| Buffer format | Return | `output_len` | Captured bytes |
| ---: | ---: | ---: | --- |
| 0 (S16PCM) | `1` | 23,606 | Byte-identical to file selector 0 |
| 1 (ALAW) | `1` | 11,803 | Byte-identical to file selector 1 |
| 2 (MULAW) | `1` | 11,803 | Byte-identical to file selector 2 |
| 3 (DADPCM) | `1` | 5,901 | Byte-identical to file selector 3 |

Each success used `flag=0`, `nThreadID=0`, and the executable's default
speaker and synthesis arguments. The outputs establish the returned byte
count and encoding parity with the matching raw file selector for this text.
Short inputs complete in the initial call. For longer input, the same flag-0
call can return `VT_BUFFER_API_PROCESSING` and expose a chunk; the follow-up
state trace below documents that lifecycle. The separate format-0 capacity
probe below establishes behavior for a smaller incoming length value on the
short synchronous call.

| Format | `flag` | `nThreadID` | Return | `output_len` after call |
| ---: | ---: | ---: | ---: | ---: |
| 0 | 0 | 0 | `1` | 23,606 |
| 0 | 1 | 0 (no active stream) | `-2` (`VT_BUFFER_API_ERROR_CREATE_THREAD`) | unchanged at 1,048,576 |
| 0 | 0 | 1 | `-2` (`VT_BUFFER_API_ERROR_CREATE_THREAD`) | unchanged at 1,048,576 |
| 0 | 1 | 1 | `-2` (`VT_BUFFER_API_ERROR_CREATE_THREAD`) | unchanged at 1,048,576 |

The `flag=1` call without an active stream and the nonzero `nThreadID=1`
combinations return the declared create-thread error in this Wine runtime.
The `flag=1` result is state-dependent: with an active stream, it drains the
next chunk, as shown below. These results do not establish nonzero-thread
support on another host/runtime.

## Chunked buffer lifecycle and mixed file/buffer state

To exceed the internal output chunk threshold, the probe repeated `Hello
world. ` six times (78 input bytes) and used a 1 MiB physical output buffer
with `nThreadID=0`. File selector 0 returned 289,636 raw PCM bytes. The buffer
API then produced four 60,000-byte chunks and one 49,636-byte final chunk:

| Call | `flag` | Return | `output_len` | State interpretation |
| --- | ---: | ---: | ---: | --- |
| Initial buffer call | 0 | `0` | 60,000 | `VT_BUFFER_API_PROCESSING`; stream remains active |
| Poll 1 | 1 | `0` | 60,000 | another chunk; stream remains active |
| Poll 2 | 1 | `0` | 60,000 | another chunk; stream remains active |
| Poll 3 | 1 | `0` | 60,000 | another chunk; stream remains active |
| Poll 4 | 1 | `1` | 49,636 | `VT_BUFFER_API_DONE`; final chunk |

Concatenating the five buffer chunks exactly matches the raw file-selector-0
output (SHA-256
`dbddb69624804043ea2075c047c5ba3b00fdee5b3d60b2c6657f476740ad75ae`). A
selector-0 file call before the buffer stream and another after it both
returned `1` and produced that same hash. This is direct file → buffer → file
coverage within one loaded process and shows no output difference for this
sequence. It does not establish all interleavings, different parameters, or
recovery after arbitrary errors.

A second trace started a stream and called flag 0 again before draining: the
second start returned `-7` (`VT_BUFFER_API_ERROR_THREAD_BUSY`) and left the
output length unchanged. Calling flag 2 on that active context returned `1`
with `output_len=0`; a subsequent flag-1 call returned `-2` and left the
length unchanged. This is consistent with flag 2 cancelling/clearing an active
per-speaker stream and flag 1 requiring an existing one. The disassembly of
`FUN_100200c0` at `0x100200c0` also implements a 60,000-byte chunk threshold,
retains state on processing, returns done on the final chunk, and frees that
state. The GDB probes execute calls sequentially on `nThreadID=0`; they recover
the chunk/poll/cancel API contract but do not demonstrate a background worker
thread. Logs and chunks are `buffer-stream*.log` and
`buffer-stream-*.bin` under Stage 16.

## Error results

| Condition | Signed return | Evidence |
| --- | ---: | --- |
| Null text pointer | `-3` (`VT_FILE_API_ERROR_NULL_TEXT`) | [`file-error-null-text.log`](../../tools/revkit/work/stage16/file-error-null-text.log) |
| Empty text | `-4` (`VT_FILE_API_ERROR_EMPTY_TEXT`) | [`file-error-empty-text.log`](../../tools/revkit/work/stage16/file-error-empty-text.log) |
| Null output path | `-6` (`VT_FILE_API_ERROR_OUT_FILE_OPEN`) | [`file-error-null-path.log`](../../tools/revkit/work/stage16/file-error-null-path.log) |
| Existing directory passed as output filename | `-6` (`VT_FILE_API_ERROR_OUT_FILE_OPEN`) | [`text-file-path-error-grid-api.log`](../../tools/revkit/work/stage21/text-file-path-error-grid-api.log) |
| Output under a nonexistent parent directory | `-6` (`VT_FILE_API_ERROR_OUT_FILE_OPEN`) | [`text-file-path-error-grid-api.log`](../../tools/revkit/work/stage21/text-file-path-error-grid-api.log) |
| Unsupported selector 6 or out-of-range selector 10 | `-1` (`VT_FILE_API_ERROR_INVALID_FORMAT`) | [`file-format-6.log`](../../tools/revkit/work/stage16/file-format-6.log), [`file-format-10.log`](../../tools/revkit/work/stage16/file-format-10.log) |

The observed returns agree with the `VT_FILE_API_*` constants in
[`vt_eng.h`](../../include/vt_eng.h). The error probes establish these cases
only; database-unloaded behavior is separately covered. The two tested
non-null output-open failures return the same `-6` as a null output path.
Other inaccessible-path forms, thread creation failures, and unknown errors
remain untested. A valid absolute Z-drive path returned `1`, produced the
expected 11,803-frame mono 16 kHz PCM16 WAVE, and byte-matched the current
selector-4 control. Reproduce with
[`run-text-file-path-error-grid.sh`](../../tools/revkit/work/stage21/run-text-file-path-error-grid.sh)
and validate with
[`analyze_text_file_path_error_grid.py`](../../tools/revkit/work/stage21/analyze_text_file_path_error_grid.py).

Buffer format 4, null text, empty text, and null output buffer returned
`-1`, `-3`, `-4`, and `-5`, respectively, matching
`VT_BUFFER_API_ERROR_INVALID_FORMAT`, `NULL_TEXT`, `EMPTY_TEXT`, and
`NULL_BUFFER`. See [`buffer-errors.log`](../../tools/revkit/work/stage16/buffer-errors.log).
The database-unloaded (`-6`), busy-thread (`-7`), abnormal-condition (`-8`),
and unknown (`-9`) codes are declared in the header. The database-unloaded and
busy-thread results are exercised below; abnormal-condition and unknown
returns were not reached.

## User-dictionary API spot checks

`VT_GetUserDictLimit_ENG` returns `30`, `10`, `50`, `65`, and `65` for
selectors 0–4, and `-1` for every other signed 32-bit selector. Static
dispatch is a five-case switch with default `-1`; runtime probes confirm the
table and test `INT_MIN`, `-1`, 5, 6, and `INT_MAX`. Capture and reproduction
are in Stage 21. The public header declares short
returns for the load/unload wrappers (the Ghidra wrapper pseudocode renders
them as `void`). Heap-backed GDB strings avoid the pointer corruption in the
earlier candidate attempts. With them, `VT_LOAD_UserDict_ENG` returned low AX
`1` for the three-field rows `hello,HH,P` at index 27 and `hello,world,A` at
index 28. It returned low AX `-1` for indexes `-1` and `1024`. Earlier
candidate attempts exposed low AX `-3`, but used suspect literal-string
pointers, so a clean invalid-file result remains unestablished.
`VT_UNLOAD_UserDict_ENG` returned `1` for both loaded slots, `-2` for index
`-1`, and `-1` for empty slot 0. The matching
extended-loader statuses in the pseudocode are `-3` for an empty parse result
and `1` for a retained dictionary pointer; the simple wrapper returns those
low-AX values at the tested call boundary.

File synthesis with text `hello`, format 4, speaker 1, settings `-1`, and
text type 0 produced the same WAVE hash with no dictionary, with the
`hello,HH,P` dictionary loaded at 27, and with `hello,world,A` loaded at 28.
The ordinary process had per-speaker dictionary gate byte 0 at `0x100a7489`,
while global dictionary slot 27 held a non-null pointer. Static synthesis
setup copies the selected dictionary pointer into its language state only
when that per-speaker gate is nonzero; otherwise it clears the language state's
dictionary pointer.
The model-load pseudocode sets the gate from `VT_CheckLicense_ENG` and the
returned user-dictionary capacity. A status-only breakpoint captured checker
return `-11`, capacity `1`, and gate `0` for the loaded Paul slot. Static
control flow maps `-11` to the final `dbsize` predicate. The license record's
`dbsize=300` fails the already measured Paul boundary: its database is
496,212 KiB, and the checker allows 1,100 KiB per `dbsize` unit, so 300 gives
330,000 KiB and the smallest passing integer is 452. The model still
synthesizes, while its user-dictionary gate remains off and capacity is
restricted to 1. The return-only capture is
[`userdict-license-state3-api.log`](../../tools/revkit/work/stage16/userdict-license-state3-api.log);
the size-boundary evidence and derivation are in the license section below.

A debugger-only experiment set that byte to 1 in a disposable process, leaving
the PE, model, license files, and on-disk configuration unchanged. With input
`hello`, the loaded `hello,HH,P` row changed the WAVE hash from the control
`613ff17ff3d7b8c9a411ef771ea1d02a051f38fcb918874e95dde30db39c0f36` to
`491b29d0c3deb1d95770f8bdfe6fbcdfce5ace03f8c15dbfb9c57cf95b601cf1`; the
`hello,world,A` row changed it to
`2a105ccfa7df1a31aa335ff9b3bc9e62065c517487747eeb0b78bb80528857bb`.
This shows both row kinds can affect this sample when the gate is enabled.
The captures do not identify the resulting spoken words by listening or
independent transcription. This debugger-forced Paul experiment is distinct
from the naturally enabled James case below.

The exact gate-off WAVs for both rows all hash to
`613ff17ff3d7b8c9a411ef771ea1d02a051f38fcb918874e95dde30db39c0f36`. Under
the debugger-forced gate-1 run, the same control hash was retained while the
`P` and `A` outputs changed to the hashes above. This paired capture directly
isolates the gate byte for these two test rows. The regular gate's `dbsize`
cause is established above. The checked-in default file at
`data-common/userdict/userdict_eng.csv` is
zero bytes. Static extended-loader pseudocode bounds indexes to 0–1023,
checks initialization and occupancy, reserves an index during parsing, and
distinguishes invalid index, uninitialized subsystem, occupied index, parse
failure, and success. The extended unloader refuses an in-use dictionary,
frees an idle loaded dictionary, and reports an empty slot. Idle unload,
forced-uninitialized `-3`, and the injected in-use-reference `-3` are runtime
confirmed; the latter two exercise controlled state rather than natural engine
conditions. Runtime logs, WAV hashes, GDB traces, and reproducible runners are
in Stage 16, including `userdict-heap-lifecycle-api.log`,
`userdict-hello-world-state-api.log`, and `userdict-hello-world-api.log`.

### Naturally license-enabled James dictionary effects

A Stage 21 GDB probe ran in the sample's loaded Paul process and loaded James
in slot 4 using the supplied verification record by file path. The loader
returned AX `0`, and direct per-slot state reads showed license gate `1` and
dictionary capacity `6`. Before loading a dictionary, format-4 synthesis of
`hello` through slot 4 with dictionary index 0 returned `1` and produced a
7,798-frame mono 16 kHz PCM16 WAVE. `VT_LOAD_UserDict_ENG(0, ...)` then loaded
`hello,HH,P` and, in a separate load/unload cycle at the same index,
`hello,world,A`; each load, synthesis, and unload returned low AX `1`.

The `P` row produced 1,150 frames and the `A` row 9,912 frames. Both outputs
were byte-different from the no-dictionary control. After both dictionaries
were unloaded, the same call returned `1` and reproduced the original control
PCM byte-for-byte. This directly confirms that the naturally enabled James
license gate allows the ordinary dictionary API to affect synthesis and that
unloading restores the no-dictionary output for this text and these rows.
It does not identify the spoken results by independent transcription, nor
generalize to other rows, voices, or license records.

Captures are `james-licensed-userdict-effect-api.log` and
`james-userdict-{control,p,a,restored}.wav` under Stage 21. Reproduce with
`run-james-licensed-userdict-effect.sh` using the Stage 8 and Stage 17 Compose
files plus the read-only `data-james` mount; validate with
`analyze_james_licensed_userdict_effect.py`.

### CSV row validation matrix

The earlier candidate-file failures used literal strings as GDB call
arguments and are not valid parser evidence. The replacement matrix passes
engine-heap-allocated paths and invokes the simple exported loader at
`0x10027960`. In the Paul M16 runtime, `hello,HH,P` loads (`AX=1`) and unloads
(`AX=1`). The following row shapes also load and unload successfully:

| CSV contents | Load AX | What this establishes |
| --- | ---: | --- |
| `hello,HH,P,extra` | 1 | A nonempty 4th field is accepted. |
| `hello,world,A,extra` | 1 | The same 4-field acceptance for type `A`. |
| `"hello","HH","P"` | 1 | This quoted sample loads and produces the same tested synthesis WAV as the unquoted row. |
| `,HH,P` | 1 | An empty source field is accepted for this row. |
| `hello,HH,P,` | 1 | An empty 4th field is accepted. |
| `hello,HH,P` with CRLF terminator | 1 | A CRLF row terminator is accepted. |
| `hello,HH,"P` | 1 | An unmatched opening quote in the final type field is accepted. |
| `he"llo,HH,P` | 1 | A quote embedded in an otherwise unquoted source field is accepted. |
| `"he""llo",HH,P` | 1 | A doubled quote inside a quoted source field is accepted. |

Two fields, five fields, an empty file, an empty `P` or `A` target, type `X`,
and a two-character type `PP` each return `AX=-3`. For every such rejection,
unload at that index returns `-1`, showing that the failed parse leaves the
slot empty; loading the valid control row at the same index then returns `1`
and can be unloaded. This demonstrates cleanup and reuse for these parse
failures, not all parser errors. Static `FUN_1005e330` pseudocode reads fields
0–2 after recognizing counts 3 or 4. A separate controlled synthesis
comparison showed identical output for three fields, a nonempty fourth field,
an empty fourth field, and fully quoted three-field ASCII input. Combined with
the parser's independently observed quote-removal behavior, this establishes
that the fourth value is ignored by this loader path for the tested rows and
that ordinary enclosing quotes are removed for these fields; it does not
characterize all malformed quote cases.

Three quote cases were isolated through the heap-backed loader runner:
`hello,HH,"P`, `he"llo,HH,P`, and `"he""llo",HH,P` each returned load `1`
and unload `1`. The standalone CSV parser's captured behavior removes an
unmatched opening quote in a quoted field, preserves a quote in an unquoted
field, and decodes doubled quotes inside a quoted field. Under the forced
dictionary gate, the unmatched-type row's `hello` WAV equals the ordinary
`hello,HH,P` row (SHA-256
`491b29d0c3deb1d95770f8bdfe6fbcdfce5ace03f8c15dbfb9c57cf95b601cf1`). The
embedded-source row leaves `hello` equal to control (`613ff17f...`) but changes
`he"llo` from its control hash (`902d035a...`) to that same dictionary-row
hash. The well-formed escaped-source row `"he""llo",HH,P` produces the same
`he"llo` hash as the embedded-quote row. These observations show that the
final type quote behaves like the plain `P` field here, while both source
forms match the quote-containing text in this fixture. They do not establish
the output's spoken content or behavior for other placements and strings.
Replay load/unload acceptance with
`run-userdict-validation-matrix.sh unmatched-type-quote`,
`embedded-source-quote`, and `escaped-source-quote`; the synthesis comparison
is `run-userdict-malformed-quote-effect.sh`. Captures are
`userdict-malformed-quote-v2-effect-api.log` and
`userdict-malformed-quote-v2-hashes.txt`.

The accepted empty-source row `,HH,P` was then tested with the per-speaker
dictionary gate forced to 1. For `hello`, control and dictionary WAVs both
hash to
`613ff17ff3d7b8c9a411ef771ea1d02a051f38fcb918874e95dde30db39c0f36`; for
`.`, both hash to
`ba584a378b11d9e9c98736fd8c256fe1453a84ee4139416d24b07acff424f0fb`.
Each load, synthesis, and unload returned `1`, and the gate was restored to
zero. This shows no measurable effect for those two tested strings; it does
not establish whether an empty source can match another token boundary, and
empty utterances cannot be tested through this file-synthesis path because
they are rejected before dictionary matching. Reproduce with
[`run-userdict-empty-source-effect.sh`](../../tools/revkit/work/stage16/run-userdict-empty-source-effect.sh);
the API trace and hash manifest are
[`userdict-empty-source-effect-api.log`](../../tools/revkit/work/stage16/userdict-empty-source-effect-api.log)
and
[`userdict-empty-source-effect-hashes.txt`](../../tools/revkit/work/stage16/userdict-empty-source-effect-hashes.txt).

An unsupported phoneme target behaves differently. The exported
`VT_CheckUserDict_TargetPhon_ENG("example")` probe returns `-9`; passing the
same `example` target in `hello,example,P` terminates the Wine process with
access-violation status `0xc0000005` during `VT_LOAD_UserDict_ENG`. A second
unsupported target, `NOTAPHONE`, reproduces the same failure. Neither call
returns a loader status, and the crashed runs cannot establish whether the
reserved slot is cleared. This is a confirmed fault boundary for the tested
build and inputs; the exact internal faulting instruction is not captured.
See [`userdict-validation-matrix-v3-api.log`](../../tools/revkit/work/stage16/userdict-validation-matrix-v3-api.log),
the full capture ending in the `NOTAPHONE` crash
[`userdict-validation-matrix-v2-api.log`](../../tools/revkit/work/stage16/userdict-validation-matrix-v2-api.log),
and the isolated `example` reproduction
[`userdict-validation-case-invalid-target-example-api.log`](../../tools/revkit/work/stage16/userdict-validation-case-invalid-target-example-api.log).
The reproducible entry point is
[`run-userdict-validation-matrix.sh`](../../tools/revkit/work/stage16/run-userdict-validation-matrix.sh).

### Extended loader input selection and state returns

Direct calls to `VT_LOAD_UserDict_EXT_ENG` at `0x10027880` establish its tested
call shape as `(index, filename, buffer, length)`. Static disassembly of
`FUN_1005e330` selects the memory reader when `buffer != NULL && length > 0`;
otherwise it passes `filename` to the file reader. A 24-call runtime matrix
varied filename null/present, buffer null/present, and lengths `0`, `1`, `9`,
`10`, `11`, and `-1`. It used the valid disk row
`vtspeakprobe,HH,P` and a 10-byte heap buffer `hello,HH,P`.

| Filename | Buffer | Lengths | Result |
| --- | --- | --- | --- |
| null | null | all six | `-3`; subsequent unload `-1` |
| null | data | 0, 1, 9, -1 | `-3`; subsequent unload `-1` |
| null | data | 10, 11 | `1`; unload `1` |
| valid file | null | all six | `1`; unload `1` |
| valid file | data | 0, -1 | `1`; unload `1` |
| valid file | data | 1, 9 | `-3`; subsequent unload `-1` |
| valid file | data | 10, 11 | `1`; unload `1` |

The partial-buffer failures despite a valid filename show the positive-length
buffer path takes precedence over the filename. Lengths 10 and 11 both accept
the ten payload bytes, with 11 including the buffer's terminating NUL. For
null buffers or lengths at most zero, the file path is used; a null filename
then fails with `-3`. This establishes the branch and these tested lengths,
not general buffer safety for arbitrary sizes. All 24 cases use distinct
indexes, and every failed parse was confirmed to leave an empty slot.
Reproduce with
[`run-userdict-extended-argument-matrix.sh`](../../tools/revkit/work/stage16/run-userdict-extended-argument-matrix.sh)
and inspect
[`userdict-extended-argument-matrix-v2-api.log`](../../tools/revkit/work/stage16/userdict-extended-argument-matrix-v2-api.log).

A Stage 21 follow-up disambiguated the memory path's advertised byte count from
the NUL-terminated payload. With `hello,HH,P` followed physically by `0xa5`
and then NUL, passing length 10 loads (`AX=1`), while length 11 returns
`AX=-3`; the failed index remains empty (`unload AX=-1`). The exact ten-byte
record is accepted without the following byte, while extending the span to
include that byte causes rejection. A 14-byte span containing
`hello,HH,P\0,X\0` loads, showing that this parser stops at the embedded NUL
even when the supplied span continues beyond it. A 17-byte span containing a
valid first row followed by a malformed second row (`hello,HH,P\n` then
`x,HH,X`) also loads. Thus the API does not require the entire advertised span
to consist solely of valid records. A paired synthesis check then loaded
`hello,HH,P\nworld,HH,P` from a 21-byte memory span with the speaker's
dictionary gate enabled. Before loading, `hello` and `world` produced distinct
control WAVE hashes (`613ff17f…` and `2a105ccf…`). After loading, both texts
produced the same WAVE hash (`491b29d0…`), matching the already observed
single-row `hello,HH,P` dictionary result. Since `world` is not the first row's
source, this establishes that the second valid row is ingested and affects
synthesis in this tested two-row buffer. It does not establish behavior for
arbitrary row counts or malformed rows in other positions. The trace and
capture are
[`trace-userdict-memory-boundaries-v3.gdb`](../../tools/revkit/work/stage21/trace-userdict-memory-boundaries-v3.gdb)
and [`userdict-memory-boundaries-v3-api.log`](../../tools/revkit/work/stage21/userdict-memory-boundaries-v3-api.log); reproduce with
[`run-userdict-memory-boundaries.sh`](../../tools/revkit/work/stage21/run-userdict-memory-boundaries.sh).
The loaded-row effect check is captured in
[`userdict-memory-multiline-effect-v2-api.log`](../../tools/revkit/work/stage21/userdict-memory-multiline-effect-v2-api.log)
with WAVE hashes in
[`userdict-memory-multiline-hashes-v2.txt`](../../tools/revkit/work/stage21/userdict-memory-multiline-hashes-v2.txt); reproduce with
[`run-userdict-memory-multiline-effect.sh`](../../tools/revkit/work/stage21/run-userdict-memory-multiline-effect.sh).

A separate direct-memory call loaded index 200 with null filename, the
10-byte buffer, and length 10 (`AX=1`). Calling the loader again at index 200
returned `-2` (`VT_LOAD_USERDICT_ERROR_INDEX_BUSY`); idle unload returned `1`.
This confirms the occupied-slot guard, but does not simulate concurrent
threads racing the reservation sentinel. In the same disposable process, the
probe saved the initialized subsystem flag at `0x100a0458` (`1`), temporarily
set it to zero, and observed loader `-4` and unloader `-3`; it then restored
the flag to `1`. These forced-state calls confirm both uninitialized branches
at runtime without showing a natural application state in which they occur.
See [`userdict-extended-buffer-v2-api.log`](../../tools/revkit/work/stage16/userdict-extended-buffer-v2-api.log)
and [`run-userdict-extended-buffer.sh`](../../tools/revkit/work/stage16/run-userdict-extended-buffer.sh).

### Fourth CSV field and unload in-use comparison

The loader accepts a four-field row, but the fourth value has no consumed
meaning in the recovered parser path. `FUN_1005e330` first accepts a row count
of 3 or 4 at `0x1005e4a4`, then calls `VT_CsvParser_GetField_ENG` only for
indexes 0, 1, and 2 (`0x1005e4b4`, `0x1005e4bf`, `0x1005e4ca`). There is no
field-3 retrieval in that function. Runtime comparison used `hello,HH,P`,
`hello,HH,P,extra`, `hello,HH,P,`, and `"hello","HH","P"`; each loaded,
synthesized, and unloaded at index 200 while the per-speaker dictionary gate
was debugger-forced to 1. All four WAVs have identical SHA-256
`491b29d0c3deb1d95770f8bdfe6fbcdfce5ace03f8c15dbfb9c57cf95b601cf1`, and the
gate was restored to 0 before normal execution resumed. The quoted row's
identical output, combined with the standalone CSV parser's quoted-field
behavior, establishes ordinary enclosing-quote removal for these ASCII fields.
The fourth column is therefore an accepted but ignored field in this build's
loader path, for the tested rows and consumer. Reproduce with
[`run-userdict-fourth-field-effect.sh`](../../tools/revkit/work/stage16/run-userdict-fourth-field-effect.sh);
the latest capture and hash manifest are
[`userdict-fourth-field-effect-v2-api.log`](../../tools/revkit/work/stage16/userdict-fourth-field-effect-v2-api.log)
and [`userdict-fourth-field-hashes-v2.txt`](../../tools/revkit/work/stage16/userdict-fourth-field-hashes-v2.txt).

The unloader's in-use guard checks each nonnull synthesis-context pointer in
the per-speaker reference array, then compares that context's `+0x1312c0`
dictionary pointer with the global dictionary slot. To exercise this return
branch without claiming a naturally concurrent use, the runtime probe loads
dictionary 200, allocates a context-sized scratch record (`0x1312e0` bytes),
sets its `+0x1312c0` word to the actual dictionary pointer, and inserts that
context pointer into speaker 1's first reference slot at `0x100a147c` (the
loaded slot's capacity is 1). `VT_UNLOAD_UserDict_EXT_ENG(200)` then returns
`-3`. The probe restores the original null slot, frees the scratch record, and
the same unload returns `1`. This confirms the pointer-match guard and normal
idle cleanup after restoration. The repeatable runner and capture are
[`run-userdict-inuse-unload.sh`](../../tools/revkit/work/stage16/run-userdict-inuse-unload.sh)
and [`userdict-inuse-unload-v5-api.log`](../../tools/revkit/work/stage16/userdict-inuse-unload-v5-api.log).

The Stage 21 follow-up reaches the guard through the ordinary synthesis path.
It starts at the sample application's `VT_TextToFile_ENG` entry, loads
`hello,HH,P` at dictionary index 0, and forces Paul slot 1's dictionary gate
from 0 to 1. The supplied Paul record fails the `dbsize` predicate, so the
gate is naturally off. The probe rewrites the current call's arguments to
synthesize `hello` as WAVE for slot 1 with dictionary index 0. At
`0x100260C2`, immediately after `FUN_10025FC0` stores the selected dictionary
pointer in the allocated context, the reference slot at `0x100A147C` held
that context; its `+0x1312C0` field exactly matched the loaded dictionary
pointer. Calling the public `VT_UNLOAD_UserDict_ENG(0)` at this point exposed
low AX `-3`. Synthesis completed with status `1`; at the call's return, idle
unload returned low AX `1`. The probe restored the gate to 0, and post-unload
synthesis returned `1`.

The same-process no-dictionary control and post-unload WAVE are byte-identical
(8,984 mono 16 kHz PCM16 frames; SHA-256
`613ff17ff3d7b8c9a411ef771ea1d02a051f38fcb918874e95dde30db39c0f36`). The
active-dictionary WAVE is 1,448 frames with SHA-256
`491b29d0c3deb1d95770f8bdfe6fbcdfce5ace03f8c15dbfb9c57cf95b601cf1`, matching
the earlier forced-gate P-row output. This confirms that the simple unloader
rejects a dictionary while a real same-thread synthesis context references
it, then permits unload after that context is released. The test does not
cover a naturally licensed Paul gate, simultaneous calls on separate threads,
other speaker slots, or every point in the context's lifetime. The runtime
log, trace, and three WAVE files are under Stage 21; validate with
`analyze_natural_userdict_inuse_unload.py`.

## Information queries

The runtime accepted every declared `VT_GetTTSInfo_ENG` request 0–26 with a
null `licensefile` and returned success (`0`). String requests returned:

| Request | Result |
| --- | --- |
| `VT_BUILD_DATE` (0) | `Jan 14 2014` |
| `VT_DB_DIRECTORY` (3) | `../` |
| `VT_DB_BUILD_DATE` (23) | Empty string |

Integer request values from this DLL/runtime were:

| Requests | Values |
| --- | --- |
| `VT_VERIFY_CODE` (1), `VT_LOAD_SUCCESS_CODE` (4), `VT_DB_ACCESS_MODE` (8), `VT_FIXED_POINT_SUPPORT` (9), `VT_MIN_VOLUME` (19), `VT_MIN_SENT_PAUSE` (22), `VT_MIN_COMMA_PAUSE` (26) | `0` |
| `VT_MAX_CHANNEL` (2), `VT_MAX_SPEAKER` (5) | `6` |
| `VT_DEF_SPEAKER` (6) | `1` |
| `VT_CODEPAGE` (7) | `1252` in this Wine environment |
| `VT_SAMPLING_FREQUENCY` (10) | `16000` Hz |
| Pitch max/default/min (11–13) | `200` / `100` / `50` |
| Speed max/default/min (14–16) | `400` / `100` / `50` |
| Volume max/default (17–18) | `500` / `100` |
| Sentence-pause max/default/min (20–22) | `65535` / `687` / `0` ms |
| Comma-pause max/default/min (24–26) | `65535` / `200` / `0` ms |

These are observed query outputs with the null-license/default lookup path;
they do not establish values for a separately supplied license or another
package. `VT_DB_BUILD_DATE` returned success but no text in this file-I/O
build.

### Output-pointer precedence and capacity boundaries

A Stage 21 runtime matrix passed a null output pointer for valid requests 1, 2,
3, 23, and 101, and invalid requests -1 and 27. Every call returned
`VT_INFO_ERROR_NULL_VALUE` (3), showing that the common pointer check precedes
request dispatch. With a nonnull sentinel buffer, invalid requests -1 and 27
returned `VT_INFO_ERROR_UNKNOWN` (2) and preserved all four sentinel bytes.
For string requests, `VT_BUILD_DATE` needs capacity 12 for `Jan 14 2014` plus
NUL; capacity 11 returns `VT_INFO_ERROR_SHORT_LENGTH_VALUE` (4) without
writing. `VT_DB_DIRECTORY` needs capacity 4 for `../` plus NUL; capacities
-1, 0, 1, and 3 return 4 without writing, while 4 copies the string. In the
missing `db_build.date` case, request 23 returns 0 and leaves its buffer
untouched even with capacity 1. Request 101 with a nonnull output and capacity
-1 returns the playback-state value as the API status and preserves the output
word. This confirms the common null-pointer and these selected capacity paths;
it is not a full pointer/capacity cross-product for every request. The trace
and capture are
[`trace-info-pointer-capacity-matrix.gdb`](../../tools/revkit/work/stage21/trace-info-pointer-capacity-matrix.gdb)
and [`info-pointer-capacity-matrix-api.log`](../../tools/revkit/work/stage21/info-pointer-capacity-matrix-api.log); reproduce with
[`run-info-pointer-capacity-matrix.sh`](../../tools/revkit/work/stage21/run-info-pointer-capacity-matrix.sh).

### Information-query edge probes

The extended runtime trace called request IDs `-1`, `27`, `28`, `100`, `101`,
`102`, and `INT_MAX`. Every ID except 101 returned `2`; ID 101 returned `0`
before the sample's file-synthesis call. The destination word was initialized
to `-1` and remained `-1`, confirming that request 101 returns its value as the
API status rather than writing an output value. The portable reproduction is
[`trace-info-extended.gdb`](../../tools/revkit/work/stage16/trace-info-extended.gdb),
[`run-info-extended.sh`](../../tools/revkit/work/stage16/run-info-extended.sh),
and its
[`info-extended-api.log`](../../tools/revkit/work/stage16/info-extended-api.log)
capture.

A contiguous request-ID sweep then called every integer from `-128` through
`256` with a four-byte nonnull destination initialized to `0xa5a5a5a5`. The
357 IDs outside the declared range `0`–`26` and the special ID `101` all
returned `VT_INFO_ERROR_UNKNOWN` (`2`) and preserved the sentinel. ID `0`
returned short-length (`4`) because this sweep intentionally supplied only
four bytes for the build-date string. IDs `1`–`26` otherwise returned success;
ID `23` left the sentinel untouched because `db_build.date` was absent. ID
`101` returned the pre-synthesis playback state `0` as the status and also
left the output untouched. Combined with the separately tested `INT_MIN` and
`INT_MAX`, this closes the contiguous unlisted-ID band around the dispatch
values; it does not enumerate all 32-bit request values. The trace and capture
are [`trace-info-request-id-sweep.gdb`](../../tools/revkit/work/stage21/trace-info-request-id-sweep.gdb),
[`info-request-id-sweep-api.log`](../../tools/revkit/work/stage21/info-request-id-sweep-api.log),
and [`run-info-request-id-sweep.sh`](../../tools/revkit/work/stage21/run-info-request-id-sweep.sh).

Request 23 (`VT_DB_BUILD_DATE`) was separately called with destination sizes
`-1`, `0`, `1`, and `2`, each time starting with two dwords of `0x5a5a5a5a`.
Every call returned `0` and left both dwords unchanged. The local DLL therefore
does not populate a database build-date string on this path; success does not
imply a NUL byte was written. This also explains why the earlier zero-filled
buffer looked like an empty string. The trace stops the inferior at the API
breakpoint after querying, before allowing the sample to synthesize to
`output.wav`. See
[`trace-info-empty-string-edges.gdb`](../../tools/revkit/work/stage16/trace-info-empty-string-edges.gdb),
[`run-info-empty-string-edges.sh`](../../tools/revkit/work/stage16/run-info-empty-string-edges.sh),
and
[`info-empty-string-edges-api.log`](../../tools/revkit/work/stage16/info-empty-string-edges-api.log).

Static disassembly further clarifies request 23. The `VT_GetTTSInfo_ENG`
branch at `0x1002a690` selects a path from cached global `DAT_1009fa4c` when
nonempty, otherwise builds a default path through `FUN_10028620`, then calls
`FUN_1002a600`. That helper concatenates the selected path with the literal
`db_build.date`, opens the resulting filename in `rb` mode, reads at most
`0x3ff` bytes into shared buffer `DAT_1009fe58`, appends a NUL, and returns the
buffer; if the file open fails, it returns null. If the helper returns null,
the API returns success immediately without writing the caller's destination.
Only a non-null helper result reaches the ordinary string-size check and copy.
No local repository input named `db_build.date` was found. The observed
success-with-untouched-buffer is therefore consistent with the missing-file
branch, not evidence that the database build date is an empty string. A
controlled runtime probe then created the expected `../db_build.date` path in
the mounted work root with the seven bytes `D23-OK\n`. Request 23 returned 0
and copied those seven bytes plus a terminating NUL when destination capacity
was 8 or greater. Capacities -1 through 7 returned `VT_INFO_ERROR_SHORT_LENGTH_VALUE`
(4) and left the sentinel buffer unchanged. This confirms the short-length
check is reached only when the file exists; with a missing file the API returns
0 before checking or writing the destination. The capture and two
reproduction runners are
[`info23-date-file-api.log`](../../tools/revkit/work/stage21/info23-date-file-api.log),
[`info23-capacity-edges-api.log`](../../tools/revkit/work/stage21/info23-capacity-edges-api.log),
[`run-info23-date-file.sh`](../../tools/revkit/work/stage21/run-info23-date-file.sh),
and [`run-info23-capacity-edges.sh`](../../tools/revkit/work/stage21/run-info23-capacity-edges.sh).
The file content and minimum-capacity result are established for this fixture;
the actual package's date-file contents remain unknown. A separate 1,100-byte
`X` file establishes the helper's
truncation boundary: it returns the first `0x3ff` (1,023) bytes and a NUL,
then the API requires destination capacity 1,024; capacity 1,023 returns 4
without writing. With a present file and null destination, request 23 returns
`VT_INFO_ERROR_NULL_VALUE` (3). These latter observations are captured in
[`info23-long-file-api.log`](../../tools/revkit/work/stage21/info23-long-file-api.log)
and reproduced by
[`run-info23-long-file.sh`](../../tools/revkit/work/stage21/run-info23-long-file.sh).

Two additional file-shape probes show that the public copy follows C-string
semantics after the helper returns. A two-line `LINE1\nLINE2\n` file is copied
in full, including both line feeds, followed by NUL; the reader does not stop
at the first line. For `A\0B\nC\n`, the destination receives only `A\0` and
the remaining sentinel bytes stay untouched. Static code scans the returned
buffer to its first NUL to determine the destination length, so this output
does not by itself reveal how many post-NUL bytes the helper read into its
private shared buffer. Captures and replay are
[`info23-multiline-api.log`](../../tools/revkit/work/stage21/info23-multiline-api.log),
[`info23-embedded-nul-api.log`](../../tools/revkit/work/stage21/info23-embedded-nul-api.log),
and [`run-info23-file-shapes.sh`](../../tools/revkit/work/stage21/run-info23-file-shapes.sh).

Static inspection explains that exception: ID `0x65` (101) returns
`DAT_100a7498` directly. `VT_PLAYTTS_ENG` sets that global to 1 on its
successful playback-start path and 0 on its failure path; playback callback
code also writes the same global. Existing runtime evidence shows the global
at 1 when the null-sink playback call returns successfully. Together this
supports interpreting request 101 as a playback-state/result query. The raw successful-playback state capture
is [`play-waveout.log`](../../tools/revkit/work/stage16/play-waveout.log).
The header's `NOT_SUPPORTED_REQUEST` code
is 1, so this successful special request can numerically return that same
value when the global is 1; callers must interpret the special request ID
before treating 1 as an error. In the inspected handler, ordinary branches
return 0, 2, 3, or 4, while request 101 returns the global; no branch returning
the header's `UNKNOWN` value 5 was found.

A follow-up queried request 101 through the API itself at three points in a
single null-sink playback call: immediately after `VT_PLAYTTS_ENG` returned,
after `VT_PAUSETTS_ENG`, and after `VT_RESTARTTTS_ENG`. All three calls
returned `1`; the output dword remained `0x5a5a5a5a`, and the backing global
was `1` at each point. This directly shows that request 101 does not distinguish
the paused interval from the started playback session in this run. The trace
and capture are [`trace-play-state-info.gdb`](../../tools/revkit/work/stage21/trace-play-state-info.gdb),
[`run-play-state-info.sh`](../../tools/revkit/work/stage21/run-play-state-info.sh),
and [`play-state-info-api.log`](../../tools/revkit/work/stage21/play-state-info-api.log).
The initial stop trace did not reach its post-stop query because the continuation
breakpoint was on padding after the `ret`. A corrected lifecycle trace redirects
the target thread into `VT_STOPTTS_ENG` and reaches the sample continuation.
During this run, request 101 returned `1` at the active stop epilogue, left its
output sentinel unchanged, and the backing global remained `1` after the output
handle had been closed and after control returned to the sample. The same trace
observed an earlier idle cleanup at play startup, where request 101 returned
`0` and the global was `0`. This closes the previously missing post-stop value
for this bounded null-sink path, but does not establish the state after natural
completion. The trace and capture are
[`trace-play-stop-lifecycle.gdb`](../../tools/revkit/work/stage21/trace-play-stop-lifecycle.gdb),
[`run-play-stop-lifecycle.sh`](../../tools/revkit/work/stage21/run-play-stop-lifecycle.sh),
and [`play-stop-lifecycle-api.log`](../../tools/revkit/work/stage21/play-stop-lifecycle-api.log).

### Pause and restart wrapper dispatch

The machine bodies of `VT_PAUSETTS_ENG` and `VT_RESTARTTTS_ENG` load the
shared output handle at `DAT_100a7490`, return through the zero-handle branch,
or call `waveOutPause`/`waveOutRestart` with that handle. The `void` wrappers
discard the MMRESULT. A target-flow trace inserted pause and restart before
playback opened a device and observed handle zero at both entries; neither
branch reached its imported WinMM call site. After valid playback opened the
ALSA null sink as handle `0xff00`, two consecutive pauses and two consecutive
restarts each reached their WinMM call site and returned MMRESULT 0. The
trace then stopped playback through the normal target flow, closed the handle,
and exited normally. The zero observed at the null-handle join is a register
value from the wrapper branch, not an API return or WinMM result. These results
map wrapper dispatch and MMRESULT handling in this Wine null-sink run; they do
not establish audible pause/resume, physical-device results, or behavior on
other WinMM implementations. Reproduce with
[`run-play-control-mmresult.sh`](../../tools/revkit/work/stage21/run-play-control-mmresult.sh);
the trace and capture are
[`trace-play-control-mmresult-v4.gdb`](../../tools/revkit/work/stage21/trace-play-control-mmresult-v4.gdb)
and [`play-control-mmresult-v4-api.log`](../../tools/revkit/work/stage21/play-control-mmresult-v4-api.log).

A GDB-driven message-pump attempt received `MM_WOM_OPEN` (`0x3bb`) and
`MM_WOM_DONE` (`0x3bd`) messages, but the inferior exited inside GDB's
synthetic `DispatchMessageA` call. A standalone PE32 host then called the
exports directly and pumped its normal Win32 thread queue. Across the seven
original utterances and two UTF-8 follow-ups, request 101 changed from `1` at
play start to `0` after the last `MM_WOM_DONE`; the output sentinel remained
unchanged. The three new null-sink calls (CP1252 control, UTF-8 `café noir`,
and UTF-8 `é noir`) all returned play result `1` and reached natural state 0.
Observed done-message counts were one for `A` and `Hello.`, two for
`Hello there.` and `café noir`, three for `A A A` and `Hello, world!`, and
eight for the longer sentence. The CP1252 control in the encoding matrix had
two done messages; each UTF-8 case had three. The `Hello there.` control
transitioned after the second callback, about 101 ms after the first and 111
ms later after the second under the host's 100 ms polling loop.

The host passed a null caller HWND and message `WM_APP+0x51`, so the caller
notifications appeared as thread messages. The two message parameters are
source-text span coordinates for these cases: `A A A` produced inclusive byte
ranges `(0,0)`, `(2,2)`, `(4,4)`; `Hello, world!` produced `(0,4)` and
`(7,11)`; CP1252 bytes for `café noir` produced `(0,3)` and `(5,8)`. The
UTF-8 `café noir` bytes (`63 61 66 c3 a9 20 6e 6f 69 72`) produced regular
spans `(4,4)` and `(6,9)` after the initial `(0,3)` notification. UTF-8
`é noir` (`c3 a9 20 6e 6f 69 72`) produced `(0,0)`, `(1,1)`, and `(3,6)`.
Thus the two bytes of each tested UTF-8 `é` occupy separate reported source
positions; the first `café` span also includes the preceding ASCII `caf`.
These captures establish byte-coordinate spans for these UTF-8 inputs, not
general Unicode decoding or DBCS behavior. The new capture and replay script
are [`playback-encoding-matrix-api.log`](../../tools/revkit/work/stage21/playback-encoding-matrix-api.log)
and [`run-playback-encoding-matrix.sh`](../../tools/revkit/work/stage21/run-playback-encoding-matrix.sh).
Spaces and punctuation are outside the tested word spans. Notifications can
repeat a span: the long sentence produced `(31,38)` twice. In that capture,
`MM_WOM_DONE` messages 2 and 3 were dispatched before the queued caller
updates `(8,8)` and `(10,22)`; the two `(31,38)` updates were likewise
delivered after done message 6. This shows caller notifications can accumulate
behind already queued audio completion messages, so their observed dispatch
order/timing alone does not identify the audio block that produced each span.
The final caller notification was `(0,-1)` after the last done event in every
capture. Static pseudocode shows the done callback
clears `DAT_100a7498`, posts `(0,-1)`, and calls
`VT_STOPTTS_ENG` on its terminal branch. The same callback posts record fields
`+0x0c` and `+0x10` as the regular notification parameters; the record writer
`FUN_1002c530` fills those fields from parser source positions. This establishes
the parameter role and final sentinel for the tested path while leaving exact
notification scheduling/repetition rules open.

The GDB trace and its capture are
[`trace-play-state-natural-completion.gdb`](../../tools/revkit/work/stage21/trace-play-state-natural-completion.gdb),
[`run-play-state-natural-completion.sh`](../../tools/revkit/work/stage21/run-play-state-natural-completion.sh),
and [`play-state-natural-completion-api.log`](../../tools/revkit/work/stage21/play-state-natural-completion-api.log). The direct-host source, build/run scripts, and captures are
[`probe-playback-natural-completion.c`](../../tools/revkit/work/stage21/probe-playback-natural-completion.c),
[`build-playback-natural-completion.sh`](../../tools/revkit/work/stage21/build-playback-natural-completion.sh),
[`run-playback-natural-completion.sh`](../../tools/revkit/work/stage21/run-playback-natural-completion.sh),
[`playback-natural-completion-host-api.log`](../../tools/revkit/work/stage21/playback-natural-completion-host-api.log),
[`run-playback-notification-matrix.sh`](../../tools/revkit/work/stage21/run-playback-notification-matrix.sh),
[`playback-notification-matrix-api.log`](../../tools/revkit/work/stage21/playback-notification-matrix-api.log),
[`run-playback-span-boundaries.sh`](../../tools/revkit/work/stage21/run-playback-span-boundaries.sh),
and [`playback-span-boundaries-api.log`](../../tools/revkit/work/stage21/playback-span-boundaries-api.log). The callback and producer pseudocode are
[`stage26-playback-callback.c`](../../tools/revkit/work/reports/stage26-playback-callback.c)
and [`stage26-playback-callback-producers.c`](../../tools/revkit/work/reports/stage26-playback-callback-producers.c).

The license-dependent requests were then repeated with the supplied
`verification.txt` path and an invalid path. Requests 1 (`VT_VERIFY_CODE`) and
2 (`VT_MAX_CHANNEL`) both returned API status `0` with the default path and
the supplied file, writing values `0` and `6` respectively. With invalid path
`License`, both API calls still returned `0`; request 1 wrote `-1`, while
request 2 wrote `1`, following its fallback branch. Thus request 1 exposes the
checker result through the output value, while request 2 returns the license
channel only when the checker succeeds and otherwise falls through to the
default speaker value. This behavior is for this DLL and tested record.

Error probes returned `2` for request 27 and `3` for a null value pointer,
matching invalid-request and null-value. String sizes are checked including
the terminator: build date size 11 returned short-length (`4`) while size 12
succeeded; DB directory size 3 returned `4` while size 4 succeeded. Integer
queries succeeded with `valuesize` 1 and even 0, writing the full 32-bit value
into the supplied 256-byte buffer. This means the caller must provide an
integer-sized destination even though this implementation does not use
`valuesize` to reject short integer requests. Full results are in
[`info-api.log`](../../tools/revkit/work/stage16/info-api.log).

## Repeated-call control

At the API entry breakpoint, the harness called selector 4 twice in the same
loaded process with the same text and arguments. Both calls returned `1`.
Their captured WAVE files are byte-identical to each other and to the
single-call selector 4 capture (`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`).
See [`file-repeat.log`](../../tools/revkit/work/stage16/file-repeat.log) and
[`repeat-first.wav`](../../tools/revkit/work/stage16/repeat-first.wav) /
[`repeat-second.wav`](../../tools/revkit/work/stage16/repeat-second.wav).
This checks one immediate identical-call sequence; it does not establish
state behavior across different selectors, parameters, errors, or longer call
sequences. The longer file → buffer → file sequence is documented above.

## Buffer length contract

The earlier successful buffer calls supplied a 1 MiB output-length value, so
they did not establish whether the value was an input capacity. A follow-up
used a physically allocated 65,536-byte buffer, initialized the `output_len`
pointer with 0, 1, 23,605, 23,606, and 23,607 in turn, and placed a one-byte
sentinel at the supplied boundary. The fixed test text again generated 23,606
bytes. Every call returned `1` and wrote `23,606` back to `output_len`.

| Supplied value | Return | Returned length | Sentinel at offset `supplied` |
| ---: | ---: | ---: | --- |
| 0 | `1` | 23,606 | overwritten |
| 1 | `1` | 23,606 | overwritten |
| 23,605 | `1` | 23,606 | overwritten |
| 23,606 | `1` | 23,606 | unchanged |
| 23,607 | `1` | 23,606 | unchanged |

The sentinel result and byte-for-byte comparison of each captured 23,606-byte
result with the existing selector-0 output show that, on this successful
synchronous path, the incoming value does not bound writes. It behaves as an
output byte count. A caller must provide sufficient physical storage; setting
a smaller `*output_len` is not a safe capacity limit. This is one short-text,
format-0 result and does not establish behavior for errors, other formats, or
threaded operation. The trace, logs, and exact-length output captures are in
[`stage16`](../../tools/revkit/work/stage16/README.md).

## Configuration setters and getters

At the file API entry breakpoint, the probe called the exported configuration
getters and setters in the already loaded Wine process. Before changing state,
speaker slot 1 returned pitch 100, speed 100, volume 200, sentence pause 925,
and comma pause 200. The public information query in the earlier probe reports
nominal volume and sentence-pause defaults of 100 and 687; those information
values are distinct from the live per-speaker state observed here.

`VT_SetPitchSpeedVolumePause_ENG(999, 999, 999, 70000, 1)` followed by
`VT_SetCommaPause_ENG(70000, 1)` succeeded as void calls. The corresponding
getters then returned pitch 200, speed 400, volume 500, sentence pause 65,535,
and comma pause 65,535. These runtime values confirm the upper clamps for this
slot; exported-function pseudocode additionally shows lower clamps of 50 for
pitch and speed, 0 for volume, and nonnegative-only updates. The pseudocode
shows out-of-range speaker IDs redirected to slot 1. A getter request for slot
0 returned `-1` in this run, consistent with no loaded speaker-state pointer
for that slot; it does not mean all six declared IDs are usable in this
executable session.

A Stage 21 direct-export matrix closes the getter pointer and selector
contract for this loaded Paul process. `VT_GetPitchSpeedVolumePause_ENG` was
called for all 16 combinations of its four output pointers being null or
non-null. Every call returned `1`; each non-null destination received only
its corresponding value, while each null destination had no effect. The
values in argument order are pitch 100, speed 100, volume 200, and sentence
pause 925. `VT_GetCommaPause_ENG` also returned `1` with either a null or
non-null output pointer; the stored value was 200.

Both exports were then called for `INT_MIN`, `-1`, slots 0–5, 6, and
`INT_MAX`. Static code first normalizes any selector outside 0–5 to slot 1.
Runtime confirms that the four tested invalid values return slot 1's values.
Slots 0 and 2–5 are valid selectors but are unloaded in this process; both
getters return `-1` for those slots and leave every nonnull output sentinel
unchanged. The field order and state reads are explicit in their wrappers:
pitch `+0x4cf4`, speed `+0x4cf0`, volume `+0x4cf8`, sentence pause
`+0x4d00`, and comma pause `+0x4d04`. Capture and replay are
[`config-getter-contract-api.log`](../../tools/revkit/work/stage21/config-getter-contract-api.log)
and [`run-config-getter-contract.sh`](../../tools/revkit/work/stage21/run-config-getter-contract.sh).

A Stage 21 direct runtime check initialized the five stored values to
`(pitch, speed, volume, sentence, comma) = (123, 234, 345, 456, 567)`, then
passed `-1`, `-2`, and `INT_MIN` individually to each field of the two pause
setters. All 15 post-call getter rows returned success and retained the full
baseline tuple. This confirms at runtime that every tested negative input is
a no-op for these loaded slot-1 setters, including large-magnitude negatives;
it does not establish treatment of values below `INT_MIN` outside the C `int`
domain. Reproduce with `run-negative-pause-setter-grid.sh` and validate with
`analyze_negative_pause_setter_grid.py` under Stage 21.

With the file API's pitch, speed, volume, and sentence-pause arguments all
left at `-1`, the selector-4 result after the setter sequence differed from
the baseline captured immediately before the setters in the same process:
configured output SHA-256
`9544ee267ff6ccf81120535ec446095e1d747ed7400620bce2ef78b17c7e1463`, versus
baseline `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`.
This supports that the combined live configuration affects file synthesis
when those call arguments are `-1`. A follow-up single-process probe reset the
initial slot-1 values before each case, changed one field to its upper or
lower limit, queried the resulting state, and synthesized the same text with
all corresponding API arguments left at `-1`. Each setter-boundary case
returned success through the getters and file API. Pitch, speed, and volume
changed the WAVE bytes; sentence-pause and comma-pause changes did not change
this sample's WAVE bytes. The per-case WAVE hashes are:

| Case | Getter value | WAVE SHA-256 |
| --- | ---: | --- |
| Baseline | pitch 100, speed 100, volume 200, sentence pause 925, comma pause 200 | `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69` |
| Pitch minimum | 50 | `a7c37f950ad6dfd9e6c343c694362fb26c5c3b91839180c3db3fe39c7be2a67c` |
| Pitch maximum | 200 | `8521dec1c644d1b1065cad2674bc057eed8c3991afa25917e19c879218aa52c1` |
| Speed minimum | 50 | `ad2c336790d40bbf221562202dcb27acea85729c9f1a41848800bd44ab9c9e8d` |
| Speed maximum | 400 | `d9b5365291226d112cd0885ea77182cc6363b7f85c8aea5c63f7209bff12ea3d` |
| Volume minimum | 0 | `2e7044b22d6cf52a536fef0ea738e2f4443e4b1dec1b96613207ab352c1a8cf2` |
| Volume maximum | 500 | `5030eb1a48de20165ba406f2fee2991ea013227d57a5ef769302b43087167004` |
| Sentence pause minimum | 0 | `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69` |
| Sentence pause maximum | 65,535 | `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69` |
| Comma pause minimum | 0 | `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69` |
| Comma pause maximum | 65,535 | `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69` |

The original sentence-pause result on `Hello world.` did not expose a pause
boundary. The Stage 21 period-context probe below establishes its effect for
three tested interword-period case forms. The EX path separately confirms
that inline VTML pause `time` is milliseconds at 16 kHz. Other sentence
contexts and output formats remain untested. The original trace, log, and WAVE
captures are under Stage 16.

### Sentence-pause synthesis effect

The Stage 21 sweep set speaker 1's sentence pause to `0`, `1`, `199`, `200`,
`201`, `250`, `500`, `924`, `925`, `926`, `65534`, and `65535`. For each value
it queried all four settings and synthesized `Hello. World.`, `Hello. world.`,
`hello. world.`, and the no-period control `Hello world.` with the file API's
other options at `-1`. All 12 getter values matched the setting and all 48
synthesis calls returned `1`.

The three period forms produced mono 16 kHz PCM16. Each has
`17,762 + 16 × pause` frames: 17,762 at zero, 20,962 at 200, 32,562 at 925,
and 1,066,322 at 65,535. Relative to pause zero, only `16 × pause`
zero-valued samples are inserted at PCM byte offset 18,006; all PCM before and
after that interval is byte-identical. The three tested case forms produce
identical PCM at each value. `Hello world.` remains 11,803 frames and
byte-identical across the entire sweep. This directly establishes that the
stored sentence-pause value controls the duration of this period-selected
interval, at 16 samples per millisecond.

Captures are `sentence-pause-<value>-<context>.wav` under Stage 21;
`run-sentence-pause-context-grid.sh` and
`analyze_sentence_pause_context_grid.py` reproduce and validate the matrix.

### Comma-pause synthesis effect

The setter had previously only been varied on `Hello world.`, which contains
no comma. A follow-up set the speaker-1 comma pause to 16 values: `0`, `1`,
`199`–`201`, `249`–`251`, `499`–`501`, `924`–`926`, `65534`, and `65535`.
For each value it queried the getter and synthesized five inputs with the
other file-API options set to `-1`: `Hello, world.`, `Hello,world.`,
`Hello,  world.`, `Hello,<TAB>world.`, and the no-comma control
`Hello world.`. All 16 getter values matched the setter inputs and all 80
synthesis calls returned `1`.

For the four comma inputs, output was mono 16 kHz PCM16. Their frame count
obeys `17,762 + 16 × pause`: 17,762 frames at zero, 20,962 at 200, 32,562 at
925, and 1,066,322 at 65,535. Comparing each PCM capture with the same
context at pause zero shows the only change is insertion of `16 × pause`
zero-valued samples at byte offset 18,006; the audio before and after that
interval is byte-identical. This directly establishes the setter's duration
effect and its millisecond scale for these file-synthesis contexts. The four
comma-spacing forms produce identical PCM for each value, so these tested
spaces and TAB do not change that insertion or the selected-unit audio. The
no-comma control remains 11,803 frames and byte-identical across all 16
settings.

The captures are `comma-pause-<value>-<context>.wav` in Stage 21. Reproduce
and validate them with `run-comma-pause-context-grid.sh` and
`analyze_comma_pause_context_grid.py`. This matrix does not cover other
speakers, negative values other than `-1`, `-2`, and `INT_MIN`, every integer
value, VTML pauses, or other output formats.

### Stored pause versus per-call pause argument

A follow-up varied one stored pause setter at a time (0 or 925), the file
API's pause argument (`-1`, `0`, or `250`), and period, comma, and punctuation-
free controls. All 36 file calls returned `1`. Outputs remained mono 16 kHz
PCM16, and every pause interval was validated as an exact zero-sample
insertion at byte offset 18,006.

For `Hello. World.`, a nonnegative file-API pause argument controls the
interword period interval: 0 suppresses it and 250 inserts 4,000 samples,
regardless of whether stored sentence pause is 0 or 925. With file-API pause
`-1`, the stored sentence value is used (0 or 925). For `Hello, world.`, the
stored comma value controls the interval (0 or 925) regardless of the file-
API pause argument. Thus the per-call value falls back to stored sentence
pause when negative; it neither overrides the stored comma pause nor affects
the punctuation-free control. The control remains 11,803 frames in all 12
settings. This resolves precedence for these three ASCII fixtures and the
16 kHz WAVE path; VTML, other formats, and other utterances remain untested.

The captures are `pause-precedence-<axis>-<stored>-<call>-<context>.wav`.
Reproduce and validate them with `run-pause-precedence-grid.sh` and
`analyze_pause_precedence_grid.py` under Stage 21.

## Text-format and VTML substitution behavior

The header declares plain text (`0`), JEITA (`4`), JEITA Plus (`6`), and UMD
(`8`). With the same `Hello world.` input and all other arguments held at the
executable's defaults, each declared value and `-1` returned success from
`VT_TextToFile_ENG`; all five selector-4 WAVE outputs were byte-identical
(SHA-256 `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`).
This establishes no distinction for this plain ASCII sample. A follow-up
used the [VTML guide](https://static.carahsoft.com/concrete/files/1615/2520/8261/Voice-Text_Markup_Language.pdf)-defined
form `<vtml_sub alias="Hello world.">x</vtml_sub>` with each declared text format.
Values 0, 4, 6, and 8 all returned success and generated the same WAVE as the
plain expansion `Hello world.` (SHA-256
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`). The
control input `x` alone also succeeded but differed (SHA-256
`157b705d9bf51bda7a7027849613992a20978b84b0195ce0bbc6e7afb71360a8`). This
supports that `<vtml_sub>` uses its alias in all four declared values in this
sample; it does not identify broader differences among JEITA, JEITA Plus,
UMD, and plain text or recover other tags, malformed-tag handling, or encoding
rules. The DLL's static tag-dispatch table is described in
[`voice-engine-and-model-formats.md`](voice-engine-and-model-formats.md#static-text-path),
and WAVE captures plus the replayable trace are under Stage 16.

## Playback API

Direct calls to `VT_PLAYTTS_ENG` in the loaded process returned `-2` for a null
text pointer and `-3` for an empty string, matching the header's null/empty
text codes. With no audio device, valid text returned `-5`
(`VT_PLAY_API_ERROR_INITPLAY`). The internal initializer constructs mono
16 kHz, 16-bit PCM and calls WinMM `waveOutOpen`; the initial Wine run failed
that open.

For a second run, the isolated Compose overlay set `ALSA_CONFIG_PATH` to
`stage16/alsa-null.conf`, a local ALSA null PCM sink. In that run,
`waveOutOpen` returned `0` and a non-null handle (`0xff00`); playback returned
`1` with internal play state `1`. Pause and restart calls returned while this
handle was active. This demonstrates successful API initialization and active
control calls through Wine with a sink that discards audio. It does not verify
audible output, acoustic quality, completion notification behavior, or a
physical audio device. See `play-waveout.log`, `alsa-null.conf`, and
`trace-play-waveout.gdb` under Stage 16.

A Stage 21 lifecycle trace then let the sample reach its successful-play
continuation, paused and restarted playback, and redirected its target thread
into `VT_STOPTTS_ENG` with that continuation as the stop call's return address.
The stop export returned and the sample continuation was reached. At the WinMM
call boundaries, `waveOutReset`, both observed `waveOutUnprepareHeader` calls,
and `waveOutClose` each returned MMRESULT 0. The global output handle was
`0xff00` during cleanup and zero after close. The stop export reached its `ret`
with raw EAX 0; this register observation is not a declared C return value.

The trace also captured an earlier idle cleanup inside `VT_PLAYTTS_ENG`: the
handle and play-state global were both zero, and request 101 returned status 0
with its output sentinel unchanged. For the active stop, request 101 returned
status 1, also left the output sentinel unchanged, and the backing global stayed
1 after close and after returning to the sample. Thus this case closes the
handle and completes the cleanup path without clearing the queried play-state
value. It does not establish a universal stopped/completed-state meaning for
that global, nor behavior for repeated/concurrent stop, WinMM failures,
completion notification, or audible hardware. See
[`trace-play-stop-lifecycle.gdb`](../../tools/revkit/work/stage21/trace-play-stop-lifecycle.gdb),
[`run-play-stop-lifecycle.sh`](../../tools/revkit/work/stage21/run-play-stop-lifecycle.sh),
and [`play-stop-lifecycle-api.log`](../../tools/revkit/work/stage21/play-stop-lifecycle-api.log).

At the wrapper boundary, `VT_PAUSETTS_ENG` checks the global handle, calls
`waveOutPause` only when it is nonzero, and discards that call's MMRESULT;
`VT_RESTARTTTS_ENG` does the same with `waveOutRestart`. The Stage 21 null-sink
run directly observed both calls return to the sample with handle `0xff00`,
and request 101 stayed `1` after start, pause, and restart. That query does not
measure playback position or prove that audio was actually paused/resumed.
GDB-injected repeated wrapper calls were also tried, but their raw EAX was not
a declared result and a later stop-state query failed in that run; those
captures are excluded as evidence for repeated-call results or WinMM status.
Null-handle behavior is static-only: both wrappers skip the WinMM call when
the global handle is zero. Pause position, repeated/concurrent calls, and
WinMM failure behavior remain uncharacterized.

## Database-unloaded errors

After calling `VT_UNLOADTTS_ENG(1)` in the loaded process, the selected valid
text returned `-5` from `VT_TextToFile_ENG`, `-6` from synchronous
`VT_TextToBuffer_ENG` (with `output_len` set to 0), and `-4` from
`VT_PLAYTTS_ENG`. These match the header's database-not-loaded constants for
the three APIs. The call sequence and return values are captured in
`db-unloaded-api.log`; the runner restores the Stage 5 input and output
fixtures when the container exits. The debugger inferior exited with code 1
after the forced unload and probes, so this validates the reported API
returns only and is not a full post-unload process-lifecycle result.

A separate probe now calls `VT_UNLOADTTS_EXT_ENG(1)` after the Paul model is
loaded. `VT_GetDBSize_ENG(1)` returns `1` with 508,121,688 before unload and
`-1` afterward; the output word is left at zero. The speaker-name helper still
returns `Paul`, showing that this static slot name remains queryable when its
engine state is gone. Valid-text calls after the direct unload return `-5`
from the file API and `-6` from the synchronous buffer API, with buffer length
0. This documents the single-slot loaded-state effects directly. The later
two-slot Paul/James probe also confirms that unloading one loaded voice
preserves the other; invalid slot normalization remains static-only. Capture:
`unload-ext-loaded-api.log`, replayed by `run-unload-ext-loaded.sh`; see also
`load-ext-multi-api.log` for the multi-slot sequence.

## Additional export queries

The Stage 16 follow-up calls four export-only query helpers after the Paul M16
model has loaded. `VT_GetSpeakerName_ENG` returned `Kate`, `Paul`, `em001`,
`Julie`, `James`, and `Ashley` for slots 0–5. A Stage 21 signed-boundary sweep
confirmed that `INT_MIN`, `-1`, `6`, and `INT_MAX` all return the same fixed
`Paul` fallback. Static code checks the signed input against `[0,5]` before
indexing the six-entry name table; this fallback does not indicate that slot 1
is loaded.
`VT_SpeakersInfo_ENG` returned `6` for each slot and copied lowercase IDs and
the DLL's embedded `d:/eng/db/.../pcm/` path strings. These strings are
metadata returned by the DLL; this trace did not open those paths.
`VT_GetDefVersion_ENG` returned `Paul-M16-FileIO`. `VT_GetDBSize_ENG` returned
`1` and `508121688` for loaded slot 1; slots 0 and 2–5 returned `-1`, leaving
the output untouched. A Stage 21 four-byte sentinel probe confirms no write on
these unloaded-slot failures. Static code normalizes every selector outside
0–5 to slot 1 before checking and reading the size; runtime calls with
`INT_MIN`, `-1`, `6`, and `INT_MAX` all return `1` and `508121688` while Paul
is loaded. The byte count is recorded as the API result without inferring which
files or allocation it measures. The database-size trace uses a scratch word
at exported data cell `0x100ff11c` and restores its original value before
leaving the inferior. See
[`export-queries-api.log`](../../tools/revkit/work/stage16/export-queries-api.log)
and its replay script, plus Stage 21
[`speaker-query-edges-api.log`](../../tools/revkit/work/stage21/speaker-query-edges-api.log)
and [`run-speaker-query-edges.sh`](../../tools/revkit/work/stage21/run-speaker-query-edges.sh).

Static pseudocode and x86 instructions show `VT_GetPathKey_ENG` selects a
speaker record with a 32-bit `selector * 24` address calculation, then formats
`SOFTWARE\\VW\\VT\\` plus the selected display name, appends `\\M`, and
appends `16`. All six valid calls return the same global buffer at
`0x100fe6e0`; a saved copy of the first result remains
`SOFTWARE\\VW\\VT\\Kate\\M16` after the call sequence, while that shared
buffer contains the slot-5 key `SOFTWARE\\VW\\VT\\Ashley\\M16`. The 32-bit
byte offset is `(selector*24) mod 2^32`. For a valid table slot `j`, the
equation `selector*24 ≡ j*24 (mod 2^32)` reduces to
`selector ≡ j (mod 2^29)`, because `gcd(24,2^32)=8` and 3 is invertible
modulo `2^29`. Thus each slot has eight signed 32-bit selector representations
`j+k*2^29` for `k=0..7`; `k=0` is the ordinary selector and the other seven
values are out of range. A Stage 21 matrix called both this export and
`VT_SpeakersInfo_ENG` with all 42 out-of-range aliases. Every call resolved to
the matching table record and returned the expected path key or metadata.
There is no selector range check in the code. Separate one-call-per-process
probes for -2, -1, 6, 7, and `INT_MAX` terminate the Wine process with status
`0xc0000005`; `INT_MAX` wraps to the same byte offset as selector -1. Those
faults apply to the tested offsets only. Other offsets that do not land on a
valid slot start remain untested. The alias matrix capture is
[`speaker-metadata-wrap-alias-wrap-matrix2-api.log`](../../tools/revkit/work/stage21/speaker-metadata-wrap-alias-wrap-matrix2-api.log),
replayed by [`run-speaker-metadata-wrap-alias-v1.sh`](../../tools/revkit/work/stage21/run-speaker-metadata-wrap-alias-v1.sh)
and [`trace-speaker-metadata-wrap-alias-v1.gdb`](../../tools/revkit/work/stage21/trace-speaker-metadata-wrap-alias-v1.gdb).
Individual fault captures are the Stage 21 `pathkey-oob-*-v1-api.log`
files, reproduced by `run-pathkey-oob-v1.sh`.

For `VT_SpeakersInfo_ENG`, slots 0–5 returned `6` and copied lowercase IDs
and embedded database paths, including NUL bytes. Measured payload lengths
excluding NUL are name `4, 4, 5, 5, 5, 6` and path `20, 20, 20, 23, 18, 20`;
therefore those observed copies require destination capacities of at least
`5/5/6/6/6/7` and `21/21/21/24/19/21` bytes, respectively. Each output was
placed in a 512-byte buffer initialized to `0xa5`; the terminator was zero,
the following byte remained `0xa5`, and no tail byte changed. `INT_MIN` was
also runtime-confirmed to alias slot 0. The disassembly copies each full
NUL-terminated source through caller pointers and accepts no size values; no
selector or null-pointer checks are visible. The fixed return value `6` may
represent the six-entry speaker catalog, but that meaning is an inference.
The same 42 out-of-range signed aliases described above were runtime-tested
with valid destinations; every call returned `6` and copied the corresponding
slot's name/path. Separate probes with selectors -2, -1, 6, 7, and `INT_MAX`
terminate the Wine process with status `0xc0000005` before return. These
faults cover the tested offsets, not every selector whose wrapped offset does
not land on a valid slot start.
Null-pointer outcomes and the all-slot exact/short boundary matrix are
documented below. Captures are the Stage 21
`speakersinfo-oob-*-v1-api.log` files, reproduced by
`run-speakersinfo-oob-v1.sh`.
The exact captures are
[`speaker-metadata-contract-api.log`](../../tools/revkit/work/stage21/speaker-metadata-contract-api.log)
and [`trace-speaker-metadata-contract.gdb`](../../tools/revkit/work/stage21/trace-speaker-metadata-contract.gdb),
replayed by
[`run-speaker-metadata-contract.sh`](../../tools/revkit/work/stage21/run-speaker-metadata-contract.sh).

### VT_SpeakersInfo destination pointer and capacity behavior

Three isolated direct calls test null output pointers on valid slot 0: null
name with a valid path buffer, valid name with null path, and both null. Every
case terminates the Wine process with status `0xc0000005` before the export
returns. Static pseudocode copies the name first, then the path; this ordering
is visible in the implementation, while the null-path capture does not record
partial destination contents after the process fault. The API takes no
capacity arguments. The existing 512-byte guard sweep establishes exact
source lengths and bounded writes for valid pointers. Additional page-boundary
calls now test every slot with both destinations exactly sized and then each
destination individually one byte short. `VirtualProtect` protects the page
immediately following each destination; every setup call returns 1 and
reports its original protection as `PAGE_READWRITE` (`0x4`). All six exact
name/path pairs return 6. Every one-byte-short name and every one-byte-short
path faults with `0xc0000005` before return. The measured boundaries are:

| Slot | Exact name/path bytes, including NUL | Exact-capacity capture | Name one byte short | Path one byte short |
| --- | ---: | --- | --- | --- |
| 0 | 5 / 21 | [`slot0 exact`](../../tools/revkit/work/stage21/speakersinfo-guard-slot0-exact-matrix1-api.log) | [`slot0 name short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot0-name-short-matrix1-api.log) | [`slot0 path short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot0-path-short-matrix1-api.log) |
| 1 | 5 / 21 | [`slot1 exact`](../../tools/revkit/work/stage21/speakersinfo-guard-slot1-exact-matrix1-api.log) | [`slot1 name short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot1-name-short-matrix1-api.log) | [`slot1 path short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot1-path-short-matrix1-api.log) |
| 2 | 6 / 21 | [`slot2 exact`](../../tools/revkit/work/stage21/speakersinfo-guard-slot2-exact-matrix1-api.log) | [`slot2 name short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot2-name-short-matrix1-api.log) | [`slot2 path short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot2-path-short-matrix1-api.log) |
| 3 | 6 / 24 | [`slot3 exact`](../../tools/revkit/work/stage21/speakersinfo-guard-slot3-exact-matrix1-api.log) | [`slot3 name short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot3-name-short-matrix1-api.log) | [`slot3 path short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot3-path-short-matrix1-api.log) |
| 4 | 6 / 19 | [`slot4 exact`](../../tools/revkit/work/stage21/speakersinfo-guard-slot4-exact-matrix1-api.log) | [`slot4 name short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot4-name-short-matrix1-api.log) | [`slot4 path short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot4-path-short-matrix1-api.log) |
| 5 | 7 / 21 | [`slot5 exact`](../../tools/revkit/work/stage21/speakersinfo-guard-slot5-exact-matrix1-api.log) | [`slot5 name short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot5-name-short-matrix1-api.log) | [`slot5 path short`](../../tools/revkit/work/stage21/speakersinfo-guard-slot5-path-short-matrix1-api.log) |

The exact sizes match the observed source-string lengths plus NUL. This
establishes copy boundaries for every compiled slot under the tested Wine
process, not validity of arbitrary caller memory. The matrix runner generates
the per-case GDB calls from
[`trace-speakersinfo-guard-matrix-template-v1.gdb`](../../tools/revkit/work/stage21/trace-speakersinfo-guard-matrix-template-v1.gdb)
and is reproduced by
[`run-speakersinfo-guard-slot-matrix-v1.sh`](../../tools/revkit/work/stage21/run-speakersinfo-guard-slot-matrix-v1.sh).

The earlier null-pointer captures remain available:
[`speakersinfo-null-name-v1-api.log`](../../tools/revkit/work/stage21/speakersinfo-null-name-v1-api.log),
[`speakersinfo-null-path-v1-api.log`](../../tools/revkit/work/stage21/speakersinfo-null-path-v1-api.log),
[`speakersinfo-null-both-v1-api.log`](../../tools/revkit/work/stage21/speakersinfo-null-both-v1-api.log),
reproduced by [`run-speakersinfo-null-v1.sh`](../../tools/revkit/work/stage21/run-speakersinfo-null-v1.sh).

## Additional helper exports and data exports

The Stage 16 helper probe ran after the Paul model loaded. `VT_GetPathKey_ENG`
returned `SOFTWARE\VW\VT\<speaker>\M16` for all six speaker slots. The
absolute and relative default user-dictionary name helpers returned
`../data-common/userdict/userdict_eng.csv` and
`data-common/userdict/userdict_eng.csv`, respectively. These are the observed
strings in the local DLL; no dictionary was loaded by this query.

`VT_INIT_ENG` is named like an initializer but is a fixed stub in this DLL.
At `0x1002aa50`, machine code is `or ax,0xffff; ret`: every call forces the
low 16 bits of EAX to `0xffff`, preserves upper EAX, and makes no memory or
global-state accesses. The existing loaded-process call captured signed
short `-1`, consistent with that instruction sequence. It performs no engine
initialization. Capture: `helper-exports-api.log` under Stage 16.

The CSV parser initializer returned a non-null object. Parsing
`alpha,"beta,gamma",delta` with flag `0` returned `1`; the field-count helper
returned `3`, the field getters returned `alpha`, `beta,gamma`, and `delta`,
and requesting index 3 returned null. `VT_CsvParser_MakeCsv_ENG` serialized
`A` and `b,c` as `"A","b,c"` with return `1`. `VT_CsvParser_IsCsv_ENG`
returned `1` for expected field count 3 and `0` for 4. The matching exit call
completed. Further ASCII edge inputs also returned `1`: `a,"unterminated`
produced fields `a` and `unterminated` (the opening quote was removed);
`a,,c` preserved an empty middle field; and `a,<LF>b` returned fields `a` and
`b`, with the LF absent from the second field. These samples do not establish
general malformed-quote or multiline-record handling. Byte probes using
CP1252 `é/ï` (`e9`/`ef`) and UTF-8 `é/ï` (`c3 a9`/`c3 af`) split on the ASCII
comma and preserved every high byte in the returned fields. This establishes
byte preservation for those samples, not codepage detection or general
Unicode semantics. A full 256-call interior-byte sweep then used
`A<byte>B,X` for every byte value. All calls returned low AX 1: NUL ended the
input after `A` and yielded one field; comma split the string into `A`, `B`,
and `X`; all other 254 nonzero byte values were preserved exactly between
`A` and `B`, including all high-byte values. This exhausts individual byte
values in that unquoted interior position but does not test byte pairs,
encoding-specific character boundaries, or normalization. The complete
capture is
[`csv-byte-domain-api.log`](../../tools/revkit/work/stage21/csv-byte-domain-api.log),
reproduced by `run-csv-byte-domain.sh`.

A Stage 21 lifecycle trace called `VT_CsvParser_Init_ENG` twice while both
objects were live. Both returned distinct 24-byte objects, with dwords
`[0,0,0,100,0,delimiter_ptr]`; their `+0x14` delimiter pointers were distinct
and each pointed to `","`. For a parsed `first,second` row, the object held
two fields, an owned text copy at `+4`, and a field-pointer array at `+8`.
`VT_CsvParser_Exit_ENG` completed for a null pointer, an unparsed object, and
a parsed object. Static destructor `FUN_10016860` returns immediately for
null; otherwise it conditionally frees the owned copy at `+4`, the field
array at `+8`, the duplicated delimiter at `+0x14`, then the 24-byte object.
It does not free `+0`, the caller's original text pointer. This maps ordinary
constructor/destructor ownership, but not persistent allocation failure or
invalid nonnull object pointers. The capture is
[`csv-lifecycle-api.log`](../../tools/revkit/work/stage21/csv-lifecycle-api.log),
reproduced by `run-csv-lifecycle.sh` and `trace-csv-lifecycle.gdb`.

A Stage 21 getter-boundary call found that `VT_CsvParser_GetNfields_ENG`
returns 0 for both a null object and a freshly initialized, unparsed object.
`VT_CsvParser_GetField_ENG` returns null for a null object, an unparsed
object, and indexes 2, 3, and `INT_MAX` after parsing a two-field row. Indexes
0 and 1 return the two stored field pointers. Index -1 returns a nonnull raw
pointer (`0x750006` in this process); it is not a valid field result and was
not dereferenced. The helper pseudocode at `0x10016910` checks whether the
index is below the field count but has no lower-bound check, then indexes the
field-pointer array directly. This explains the negative-index result without
assigning meaning to the value read before the array. The capture is
[`csv-getter-boundaries-api.log`](../../tools/revkit/work/stage21/csv-getter-boundaries-api.log),
reproduced by `run-csv-getter-boundaries.sh`.

`VT_CsvParser_IsCsv_ENG` has a distinct bounded, first-record path. A direct
matrix called the export on `a,b,c` with expected field counts `INT_MIN`,
`-1`, `0`–`4`, and `INT_MAX`, with byte limit `64`: raw low AX was `1` for
the tested values `INT_MIN`, `-1`, and `0`–`3`, and `0` for `4` and
`INT_MAX`. This establishes an
at-least predicate (`expected <= parsed field count`) for nonempty input,
rather than exact equality. With expected count 3, limits `-1`, `0`, `2`, and
`4` returned 0, while limits `5` and `6` returned 1 for the same five-byte
row. Empty input returned 0 for expected counts 0, 1, and 2, even though the
standalone parser's empty-input path reports one field.

The third argument's bound is added to the input address with 32-bit pointer
arithmetic before the bounded scan. On the captured row address `0x003e0670`,
limits `INT_MIN` and `INT_MAX` produce end addresses `0x803e0670` and
`0x803e066f`; both calls return 1 because the scanner reaches the row's NUL
before that high bound. Limits -1024, -1, 0, 1, and 4 return 0. A negative
limit that wraps the computed end below the row address (end `0`, `1`, `4`, or
`5`) also returns 0; a limit that wraps to end `0xffffffff` returns 1. Together
with the positive limit 5/6 cases, this shows the bound is treated as an
unsigned end pointer after 32-bit addition, not rejected merely for being
negative. The exact result of a signed limit therefore depends on the input
address. This is a bounded observation on a valid NUL-terminated row, not a
safe general contract for arbitrary pointers or unterminated buffers. The
capture is [`csv-iscsv-bounds-v3-api.log`](../../tools/revkit/work/stage21/csv-iscsv-bounds-v3-api.log),
reproduced by [`run-csv-iscsv-bounds-v3.sh`](../../tools/revkit/work/stage21/run-csv-iscsv-bounds-v3.sh)
and [`trace-csv-iscsv-bounds-v3.gdb`](../../tools/revkit/work/stage21/trace-csv-iscsv-bounds-v3.gdb).

For both `a,b\r\nc,d` and `a,b\nc,d`, expected count 2 returned raw low AX 1
and expected count 3 returned 0; each call used a fresh input buffer. Both
calls changed the caller's source string to `a,b`, showing that this export
truncates the input at the first CR/LF and checks only the first record.
Static disassembly at
`0x10016b50` computes an end pointer from argument 3, then
`FUN_1002e1b0` searches within that bound, writes NUL at a CR/LF boundary,
and passes the resulting row to the parser with copy flag 2 before comparing
the expected count with the resulting field count. The direct capture is
[`csv-iscsv-matrix-api.log`](../../tools/revkit/work/stage21/csv-iscsv-matrix-api.log),
reproduced by `run-csv-iscsv-matrix.sh` and
`trace-csv-iscsv-matrix.gdb`. The export is decompiled as `void`; AX values
here are raw register observations. Other bounds, newline forms, and malformed
input interactions remain open. Pointer-wrap outcomes at other input addresses,
unterminated inputs, and inaccessible memory at or before the computed end
remain untested.

The same probe tested 98–101-field rows. Direct parser calls returned 1 and
stored the requested count for 98, 99, and 100 fields. At 101 fields, parsing
returned `-4` with `GetNfields` still reporting 100. For that same 101-field
row, `IsCsv` returned raw low AX 1 at expected count 100 and 0 at 101 and
102. The helper's disassembly calls `VT_CsvParser_Parsing` and then reads the
field count without branching on the parse result, so the runtime shows that
`IsCsv` evaluates the parser's retained partial count after this capacity
error. Its initializer sets capacity to 100. Static parsing increments the stored field count, then compares it with that capacity at `0x100167ee`; when more input remains at the capacity, it returns `-4` at `0x100167f5` without reducing the stored count. The expanded Stage 21 sweep covers every count 99–128 and 255, 256, 512, 1,024, and 4,096, using a byte bound that includes each complete row and its NUL. Counts 101 through 4,096 all return parser `-4` and retain exactly 100 fields. For each count above 100, `IsCsv` returns 1 for expected 100 and 0 for expected count and count+1. This confirms the partial-count behavior well beyond the first overflow and distinguishes parser status from the count predicate. The capture is [`csv-iscsv-capacity-v2-api.log`](../../tools/revkit/work/stage21/csv-iscsv-capacity-v2-api.log), reproduced by [`run-csv-iscsv-capacity-v2.sh`](../../tools/revkit/work/stage21/run-csv-iscsv-capacity-v2.sh) with [`trace-csv-iscsv-capacity-v2.gdb`](../../tools/revkit/work/stage21/trace-csv-iscsv-capacity-v2.gdb).

A separate edge probe observed that doubled quotes in
`"a""b",c` and a quote in unquoted `a"b,c` both returned `a"b` and `c`;
`a,` omitted its empty final field; and `a,\r\nb` removed the CRLF at the
beginning of the second field. In `a,b<LF>c,d`, the LF remained inside the
second field (`b<LF>c`) and the parser returned three fields. The exact row
`a,"b,c",d` returned the same three fields with parse flags 0, 1, and 2. A
new flag matrix covered `INT_MIN`, -1, 0, 1, 2, 3, 255, and `INT_MAX`; every
call returned 1 and produced the same three fields. Only flag 1 changed the
caller's buffer, which became `a` after the parser overwrote the first comma
with NUL. The machine code at `0x1001676a` branches on the exact comparison
`flag == 1`; flag 1 selects the caller buffer, while every other integer
selects the parser-owned copy. This closes the flag's control behavior across
the 32-bit integer extremes without claiming all CSV grammar behavior. A
41-case byte matrix then expanded quote and line-position coverage; every call
returned low AX 1. For details, see
[`csv-parser-deep-4-api.log`](../../tools/revkit/work/stage16/csv-parser-deep-4-api.log)
and the case-to-input hex mapping in
[`run-csv-parser-deep.sh`](../../tools/revkit/work/stage16/run-csv-parser-deep.sh).

Runtime and static observations from that matrix establish these bounded
rules. The parser skips ASCII tab, LF, CR, and space before a field. If the
first remaining byte is `"`, the helper at `0x10016810` removes the opening
quote, treats doubled quotes as one literal quote, and ends at the first
non-doubled closing quote; if no closing quote exists, it returns the bytes to
NUL without the opening quote. A quote after ordinary unquoted bytes is kept
literally, and does not protect a comma. After a quoted field, the parser skips
the same whitespace and consumes the next non-whitespace byte as the field
separator; the malformed input `"a"x,b` therefore yielded `a`, empty, and `b`.
That separator behavior is consistent with the pseudocode's one-byte advance,
not a claim that arbitrary separators are supported CSV.

Empty input and a lone comma each yield one empty field. `,a` yields an empty
field then `a`; `,,` yields two empty fields; `a,` omits a final empty field,
while `a,,` retains the empty field before the last delimiter. LF is not a
record separator in the tested inputs: leading LF/CRLF and LF/CRLF immediately
after a comma are skipped as field-leading whitespace, while embedded and
trailing LF/CRLF bytes remain inside the returned field. Bare CR is likewise
skipped at field start, but is preserved when embedded, quoted, or immediately
before a comma. The result is a permissive single-record field splitter, not
evidence of conventional multiline CSV handling. The exhaustive short quote
product below extends the malformed-quote evidence for its stated finite
alphabet and length.

### Exhaustive short quote/comma product

A direct-export matrix enumerated every byte string of length 0 through 7
over the alphabet `A` (`41`), double quote (`22`), and comma (`2c`): 3,280
inputs total. The original length-0–6 matrix had 1,093 cases; a separate
length-0–7 run repeated those inputs and added all 2,187 length-7 strings. One
parser object was reused for each run with parse flag 0; every call returned
low AX 1. The trace recorded each input as hex, the exported field count, and
every field returned by `VT_CsvParser_GetField_ENG`. Thus even malformed
quote arrangements in this bounded product produce a successful parse
result; the output shape, rather than the return value, distinguishes them.

| Input length | 1 field | 2 fields | 3 fields | 4 fields | 5 fields | 6 fields | 7 fields |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| 0 | 1 | 0 | 0 | 0 | 0 | 0 |
| 1 | 3 | 0 | 0 | 0 | 0 | 0 |
| 2 | 6 | 3 | 0 | 0 | 0 | 0 |
| 3 | 15 | 9 | 3 | 0 | 0 | 0 |
| 4 | 33 | 33 | 12 | 3 | 0 | 0 |
| 5 | 75 | 96 | 54 | 15 | 3 | 0 |
| 6 | 171 | 270 | 189 | 78 | 18 | 3 |
| 7 | 393 | 726 | 624 | 315 | 105 | 21 | 3 |

The counts exhaust only this alphabet through length 7. All 2,187 added
length-7 cases returned low AX 1, with field counts from one through seven;
only three produced seven fields: six leading commas followed by `A`, a
quote, or a comma. The earlier 41-case matrix covers selected
spaces, tabs, CR/LF, and other quote placements; quote behavior combined with
arbitrary whitespace, other payload bytes, longer malformed strings, and
alternate encodings remains outside the combined matrices. Reproduce the
length-0–6 product with `run-csv-quote-domain.sh` and the length-0–7 product
with `run-csv-quote-domain-len7.sh`; the added trace and capture are
[`trace-csv-quote-domain-len7.gdb`](../../tools/revkit/work/stage21/trace-csv-quote-domain-len7.gdb)
and
[`csv-quote-domain-len7-api.log`](../../tools/revkit/work/stage21/csv-quote-domain-len7-api.log).

The parser object is 0x18 bytes in the reviewed initializer. Its pointer field
at `+0x14` is initialized from a duplicated delimiter-set string whose bytes
are `2c 00` (comma, NUL). The scanner `FUN_10064a30` builds a byte-membership
table from that string and returns the first matching input-byte offset. A
runtime probe changed this private pointer to `";"`: `a;b;c` then returned
three fields (`a`, `b`, `c`), while default comma configuration returned one
field (`a;b;c`). With the private semicolon setting, `a;"b;c";d` returned
three fields and kept `b;c` intact as the quoted middle field. This is a
debugger-forced opaque-object mutation; the exported initializer has no
delimiter argument, the parse call has no delimiter parameter, and no
exported setter was found. It therefore characterizes internal capability,
not a supported caller configuration contract. The serializer separately
uses comma and quote literals in its implementation.

A direct single-byte sweep then passed each value `0x00`–`0xff` as the sole
field string to `VT_CsvParser_MakeCsv_ENG`, with capacity 8 and prefix/post
guards. For all 255 nonzero byte values, raw low AX was 1 and the exact output
was an opening quote, the byte, a closing quote, and NUL; this includes comma,
CR/LF, and every high-bit byte. The quote byte `0x22` was doubled, producing
four consecutive quote bytes followed by NUL. Input `0x00` is an empty C
string: the call returned raw low AX -1 after writing only the opening quote;
the rest of the initialized output remained `0x31` through the final
capacity NUL. Both guards stayed intact in all 256 cases. The `void` wrapper
means AX is an observed register value rather than a declared return. A full
ordered-pair matrix then tested every backing buffer `[a,b,0]` for all 256×256
byte pairs, at capacity 8. All 65,536 guards stayed intact. For the 256 cases
with `a=0`, the C string is empty regardless of `b`; each returned raw low AX
-1 after writing only the opening quote. The other 65,280 calls returned 1.
For every nonempty case, a first byte with its high bit set consumes any
following nonzero byte as a pair and copies both raw, including a quote or
comma in second position. A trailing high-bit byte (`b=0`) is copied alone.
When `a` is ASCII, bytes are processed separately: quotes are doubled in
either position, while every nonzero second byte, including high-bit bytes,
is copied as one byte. The data therefore fully maps this helper's one- and
two-byte NUL-terminated behavior. It does not identify a locale or Windows
codepage: the helper treats every `0x80`–`0xff` first byte as a lead byte,
without a runtime lead-byte table check. This is a byte-pair rule observed in
the helper, not a claim of valid DBCS decoding. The capture and independent
verifier are
[`csv-makecsv-byte-pairs-full-native14.bin`](../../tools/revkit/work/stage21/csv-makecsv-byte-pairs-full-native14.bin),
[`csv-makecsv-byte-pairs-full-native14-api.log`](../../tools/revkit/work/stage21/csv-makecsv-byte-pairs-full-native14-api.log),
and [`verify_csv_makecsv_byte_pairs_full.py`](../../tools/revkit/work/stage21/verify_csv_makecsv_byte_pairs_full.py).
Reproduce with
[`run-csv-makecsv-byte-pairs-full-v1.sh`](../../tools/revkit/work/stage21/run-csv-makecsv-byte-pairs-full-v1.sh),
whose native batch executes the API inside the target process.
The earlier selected-class capture remains
[`csv-makecsv-byte-domain-matrix4-api.log`](../../tools/revkit/work/stage21/csv-makecsv-byte-domain-matrix4-api.log),
[`trace-csv-makecsv-byte-domain-v1.gdb`](../../tools/revkit/work/stage21/trace-csv-makecsv-byte-domain-v1.gdb),
[`run-csv-makecsv-byte-domain-v1.sh`](../../tools/revkit/work/stage21/run-csv-makecsv-byte-domain-v1.sh),
[`csv-makecsv-byte-pairs-classes1-api.log`](../../tools/revkit/work/stage21/csv-makecsv-byte-pairs-classes1-api.log),
[`trace-csv-makecsv-byte-pairs-v1.gdb`](../../tools/revkit/work/stage21/trace-csv-makecsv-byte-pairs-v1.gdb),
and [`run-csv-makecsv-byte-pairs-v1.sh`](../../tools/revkit/work/stage21/run-csv-makecsv-byte-pairs-v1.sh).

The field/count/capacity follow-up tested 87 arrays in 2,404 calls. It covers
every placement of `A` and an empty field through counts 0–5, each with every
capacity 1–32; arrays of 8, 16, 32, 64, and 128 fields with an empty field at
the beginning, middle, or end or with no empty field; and four mixed arrays
containing ASCII quotes and the raw pair `80 22`. The scaled arrays include
capacities immediately below, at, and above their complete-output size.
Every captured advertised output byte and low-AX result matched the
independent model; all prefix/post-capacity guards were unchanged.

The empty-field failure is incremental. Given sufficient capacity, an empty
field at index `k` preserves the serialized preceding fields and separator,
writes its opening quote, then returns low AX -1. It does not serialize later
fields. For example, `A`, empty, `A` leaves `"A","` before the remaining
`0x31` fill and final capacity NUL. A capacity failure can occur before the
empty field is reached. Count zero returns 1 and writes an empty string for
all capacities 1–32. For `N` fields containing `A`, the complete string has
`4*N-1` bytes, so capacity `4*N` is the first success; this was cross-checked
for each `N=1`–5 and `N=8,16,32,64,128`. All smaller tested capacities return
-1 with partial initialized output.

The same matrix confirms the sentinel checks at each write boundary.
Ordinary bytes, opening/closing quotes, and separators reserve one output
byte. A doubled ASCII quote or high-bit pair reserves two bytes together;
when the advertised final NUL falls within that chunk, neither byte is
written. The runtime results match the checks in `FUN_10016a40` and
`FUN_1002ea80`. These results cover counts through 128 and capacities through
513 for the selected arrays; larger counts, near-32-bit cursor wrap, and other
pointer faults remain untested.

Capture:
[`csv-makecsv-fields-matrix1-api.log`](../../tools/revkit/work/stage21/csv-makecsv-fields-matrix1-api.log).
The generator and independent verifier are
[`csv_makecsv_field_matrix.py`](../../tools/revkit/work/stage21/csv_makecsv_field_matrix.py),
with replay runner
[`run-csv-makecsv-fields-v1.sh`](../../tools/revkit/work/stage21/run-csv-makecsv-fields-v1.sh).

The pointer-order follow-up ran 15 isolated processes and checked the terminal
return or fault instruction plus all 16 bytes of a sacrificial output region.
It establishes three distinct stages in `FUN_10016960`: output initialization;
the 32-bit field-array entry read at `0x100169a7`; and, after the opening quote
has been written, the input C-string scan at `0x10016a53` in `FUN_10016a40`.
The field-array entry read precedes the check for room for the opening quote.
The array cursor advances by four bytes per entry.

| Input arrangement | Capacity | Observed outcome |
| --- | ---: | --- |
| Null array, count 1 | 1 | Fault at `0x100169a7`; output byte 0 had already been initialized to NUL. |
| Accessible array containing a null or `0xffffffff` first string pointer | 1 | Returns low AX -1; opening-quote capacity check fails before the input string is scanned. |
| Same first string pointers | 2 | Writes the opening quote, then faults at the string scan `0x10016a53`. |
| `A`, null string pointer | 4, 5 | Returns -1 with `"A"` or `"A",` respectively, followed by the final capacity NUL. |
| `A`, null string pointer | 6 | Writes `"A","`, then faults at `0x10016a53`. |
| Only the first array entry accessible, containing `A`; second entry on a no-access page | 4 | Returns -1 at the separator boundary before fetching the second entry. |
| Same guarded array | 5, 6 | Faults at `0x100169a7` while fetching the second entry, even when capacity 5 has no room for its opening quote. |
| Empty first string followed by null, counts 2 and `INT_MAX`; or by a protected array entry, count `INT_MAX` | 8 | Returns -1 after the first opening quote; the next entry is not reached. |
| Null array, count `INT_MIN` | 1 | Returns 1 with an empty output; signed nonpositive count skips the array loop. |

The guarded array is a single four-byte pointer at allocation offset `0xffc`;
the next page, starting at `+0x1000`, is `PAGE_NOACCESS`. The runtime captures
confirm protection setup succeeded. The count-`INT_MAX` cases demonstrate
early failure before a later entry is read, not successful traversal of a
large array. Each process uses a manually prepared target stack so GDB catches
faults directly at the engine instruction; after capture, the probe kills that
process. These faulting calls did not produce an API return value.

The complete expected-case inventory, trace generator, and verifier are
[`csv_makecsv_pointer_order.py`](../../tools/revkit/work/stage21/csv_makecsv_pointer_order.py).
Representative captures are
[`csv-pointer-order-matrix1-null_array_cap1-api.log`](../../tools/revkit/work/stage21/csv-pointer-order-matrix1-null_array_cap1-api.log),
[`csv-pointer-order-matrix1-null_first_cap2-api.log`](../../tools/revkit/work/stage21/csv-pointer-order-matrix1-null_first_cap2-api.log),
and [`csv-pointer-order-matrix1-guard_array_cap5-api.log`](../../tools/revkit/work/stage21/csv-pointer-order-matrix1-guard_array_cap5-api.log).
Replay with
[`run-csv-pointer-order-v1.sh`](../../tools/revkit/work/stage21/run-csv-pointer-order-v1.sh).
Other malformed pointer graphs, persistent allocation failure, and large
successful traversals remain open.

A 95-call overlap matrix then supplied the initial field `AB` at `output+k`
for each offset `k=0`–9 and every capacity from `k+3` through 16. In every
case the initial field, including its NUL, lies inside the region that the
serializer initializes. Each capture includes all 18 bytes of the physical
output allocation, raw low AX, and a check that the array entry still holds
the same source pointer. All snapshots and returns matched the independent
model; the array pointer and surrounding bytes remained unchanged.

The initialization destroys the original `AB` bytes before input processing.
At offset 0, the opening quote also changes the first input byte; subsequent
quote doubling overwrites input bytes that have yet to be consumed, propagating
quotes until the sentinel prevents another pair write. At offset 1, the helper
copies the initialized bytes to the same addresses and fails at the closing
quote. At offsets 2–9, it copies the initialized bytes toward lower addresses
and returns 1 after serializing those bytes. At capacity 8, the exact results
for offsets 0, 1, and 2 are `22×7 00` / -1, `22 31×6 00` / -1, and
`22 31×5 22 00` / 1 respectively. Thus a return of 1 in these overlap layouts
does not demonstrate preservation of the original field content. These
observations cover only the specified offsets with the complete initial
field inside the initialized region; other overlap layouts remain untested.

Capture and reproducible tooling:
[`csv-makecsv-overlap-matrix1-api.log`](../../tools/revkit/work/stage21/csv-makecsv-overlap-matrix1-api.log),
[`csv_makecsv_overlap.py`](../../tools/revkit/work/stage21/csv_makecsv_overlap.py),
and [`run-csv-makecsv-overlap-v1.sh`](../../tools/revkit/work/stage21/run-csv-makecsv-overlap-v1.sh).

For `VT_CsvParser_MakeCsv_ENG`, a separate
call serialized `A` and `b,c` through every advertised capacity from 0 through
64 into a 64-byte output region. Capacity 10 is the first to produce the
complete nine-byte CSV text plus NUL and raw low AX 1; capacities 1–9 produce
raw low AX -1 with partial or empty output. Capacities 10–64 return raw low AX
1.
Capacity 0 also emits the full CSV text despite advertising no space. Static
disassembly at
`0x10016960` writes a terminator to `buffer + capacity - 1`, making that address
`buffer - 1` for capacity zero. A sacrificial prefix byte set to `0x5a`
changed to `0` at capacity 0 and stayed unchanged for positive capacities.
The function first fills the advertised buffer with byte `0x31` (`'1'`),
places a NUL at the final advertised byte, then serializes. Unused bytes after
the serialized terminator retain the fill byte. With fields `a"b` and `plain`,
the runtime output was `"a""b","plain"`, confirming embedded quote doubling.
The disassembly shows how the capacity boundary is enforced: before each
serialized byte, it checks whether the current destination byte is NUL; the
pre-planted NUL at `buffer + capacity - 1` therefore stops a too-large result.
There is no comparison between the running cursor and the capacity argument.
The capacity argument is consumed as an unsigned 32-bit value for fill and
final-byte placement. In contrast, the serialized cursor and the helper's
returned field length are 32-bit signed values; embedded quotes add two output
bytes, and the cursor is advanced with ordinary 32-bit `inc`/`add` operations.
This statically establishes that near-2-GiB serialized strings can wrap the
cursor before the advertised capacity boundary is reached. It does not
establish the resulting writes or return value: reaching that state requires
caller-owned field strings and an output allocation on that scale, and no
runtime overflow probe was run. The current successful 128-MiB probe remains
below this signed-cursor boundary. Capacity zero still places the sentinel at
`buffer - 1`, as directly observed above.
The underwrite stayed inside the probe allocation. The wrapper is decompiled
as `void`, so captured AX values are raw register state, not declared C
returns. A null field-array pointer with field count zero returned raw AX 1
and produced an empty string. With counts -1 and `INT_MIN` plus a null
field-array, the export likewise returned raw AX 1 and produced an empty
string; the output retained its `0x31` fill bytes except for the leading and
final NUL terminators. This matches the decompiled `if (0 < count)` traversal
guard. It does not establish behavior for every negative count or invalid
pointer combination. Null output, null field-array with positive count, and
null field elements remain unprobed. The bounded cases are recorded in
[`csv-negative-count-api.log`](../../tools/revkit/work/stage16/csv-negative-count-api.log)
and reproducible with
[`run-csv-negative-count.sh`](../../tools/revkit/work/stage16/run-csv-negative-count.sh).
Three more calls each ran in a fresh Wine process: a null output buffer with
capacity 8/count 0, a null field-array with count 1, and a one-entry array
whose field pointer was null. Each process exited with Windows status
`0xc0000005` (3221225477) while the export call was active, and none returned
to its post-call GDB print. Static pseudocode places the first invalid write
in the output fill, the second at `*param_4`, and the third in the field
string scan. This establishes an access violation for these exact inputs, not
a recoverable API error code. Captures and replay files are
[`csv-pointer-edges-0-api.log`](../../tools/revkit/work/stage16/csv-pointer-edges-0-api.log),
[`csv-pointer-edges-1-api.log`](../../tools/revkit/work/stage16/csv-pointer-edges-1-api.log),
[`csv-pointer-edges-2-api.log`](../../tools/revkit/work/stage16/csv-pointer-edges-2-api.log),
and [`run-csv-pointer-edges.sh`](../../tools/revkit/work/stage16/run-csv-pointer-edges.sh).
A new large-capacity probe extended the
guarded buffer to 65,537 bytes;
capacities 65, 128, 1,024, 4,096, and 65,536 all returned low AX 1, emitted
the same CSV bytes, placed the final capacity NUL at the advertised last byte,
preserved a prefix guard, and left the byte after the advertised region
unchanged. This confirms ordinary large capacities at those sample points;
overflow and allocation-failure behavior remain open. The exact trace is
[`csv-capacity-api.log`](../../tools/revkit/work/stage16/csv-capacity-api.log)
and runner is
[`run-csv-capacity.sh`](../../tools/revkit/work/stage16/run-csv-capacity.sh).
The extended capacity capture and trace are
[`csv-capacity-large-api.log`](../../tools/revkit/work/stage16/csv-capacity-large-api.log)
and [`trace-csv-capacity-large.gdb`](../../tools/revkit/work/stage16/trace-csv-capacity-large.gdb).
The matched-allocation sweep was extended to 131,072, 1,048,576, 4,194,304,
16,777,216, and 67,108,864 bytes. Each call returned raw low AX 1, retained the prefix and
post-buffer guards, wrote the same short serialization, and placed NUL at the
last advertised byte. A new 134,217,728-byte call also returned low AX 1,
emitted the same bytes, retained both guards, and placed NUL at the final
advertised byte. The implementation `FUN_10016960` and its helper chain
(`FUN_10016a40`, `FUN_1002ea80`, and `FUN_10063f30`) operate on the caller's
field pointers and output buffer; the reviewed call graph contains no heap
allocation. Thus allocation failure is a caller-side condition before this
API call, not an internal MakeCsv path. Near-32-bit output-size/count
overflow and huge unsigned capacities remain untested. The new capture and
runner are
[`csv-capacity-128m-api.log`](../../tools/revkit/work/stage21/csv-capacity-128m-api.log)
and [`run-csv-capacity-128m.sh`](../../tools/revkit/work/stage21/run-csv-capacity-128m.sh);
the decompilation is
[`stage25-makecsv-internals.c`](../../tools/revkit/work/reports/stage25-makecsv-internals.c).
The previous capture and runner are
[`csv-capacity-xlarge-api.log`](../../tools/revkit/work/stage16/csv-capacity-xlarge-api.log)
and [`run-csv-capacity-xlarge.sh`](../../tools/revkit/work/stage16/run-csv-capacity-xlarge.sh).
The parser-edge trace and runner are
[`csv-parser-edges-api.log`](../../tools/revkit/work/stage16/csv-parser-edges-api.log)
and [`run-csv-parser-edges.sh`](../../tools/revkit/work/stage16/run-csv-parser-edges.sh).
The expanded flag trace and capture are
[`trace-csv-flag-matrix.gdb`](../../tools/revkit/work/stage16/trace-csv-flag-matrix.gdb)
and [`csv-flag-matrix-api.log`](../../tools/revkit/work/stage16/csv-flag-matrix-api.log).
The CP1252/UTF-8 raw-byte probe is captured in
[`csv-encoding-edges-api.log`](../../tools/revkit/work/stage16/csv-encoding-edges-api.log)
with its replay trace at
[`trace-csv-encoding-edges.gdb`](../../tools/revkit/work/stage16/trace-csv-encoding-edges.gdb).

The sync-info helpers allocate and operate on a concrete nested structure.
Ghidra pseudocode shows a 56-byte header (14 dwords), a 600-row array with
36-byte rows, and one separately allocated 520-byte nested area per row,
holding 65 entries of 8 bytes. These are observed layout/control-flow facts
from pseudocode. The extended load/unload path
stores the object at state offset `+0x47774`, allocates and initializes it
during setup, then frees it during teardown. This establishes an internal
per-state lifecycle. The initializer sets header field 3 to zero, fields 1/2
to 600/65, and fields
4–13 to
`[-1,0,-1,0,-1,-1,0,-1,0,-1]`. It clears the 16-bit row value at offset 0,
the seven dwords at row offsets 8–32, and the 4-byte and 2-byte values in each
nested entry. It preserves the header row-array pointer and row nested
pointers. Runtime observations matched this on sampled rows 0, 1, and 599,
including unchanged nested pointers and cleared values.

Cross-references in internal pseudocode and a live `VT_TextToBufferEX_ENG`
capture establish more of the record contract. The producer `FUN_1002c530`
uses header field 3 as a row write cursor, advances it, and wraps it at 600.
The consumer `FUN_1001e0c0` walks pending rows relative to that cursor. Row
offset 0 is the nested-entry count; offset 4 is its separately allocated
pointer. For each nested entry, the consumer adds the first dword to a frame
accumulator. Its live values sum exactly to the row's offset-8 value for the
captured rows (for example, `0x586 + 0x38a + 0x51d + 0x93d = 0x176a`). The
row offset-8 total is an audio-frame count: the EX marker mapper uses it in
frame-coordinate sums, and earlier PCM16/A-law/μ-law output-length checks
confirmed the unit. The nested entry's following 16-bit word selects the
consumer's formatting branch; `0x28` selects the silence format, while the
other observed values select the ordinary entry format. No broader names for
those selector values have been recovered.

The producer fills row offsets 12/16 with a source-text interval. On the live
`Hello world.` repetitions, the pairs `[0,4]` and `[6,10]` cover the inclusive
source character positions of `Hello` and `world`; the consumer uses these
bounds when mapping a marker's source position to timeline rows. The other
four row dwords now have a traced source chain:

- Row `+20` copies synthesis-state `+0x64c`. `FUN_10022dc0` initially fills
  the per-record `+0x64c/+0x650` pair with parser-output byte coordinates;
  when a record spans multiple input bytes, the end is stored inclusively.
  The report formatter `FUN_100180b0` resolves the coordinate through an
  input-pointer-minus-one expression, while `FUN_10022970` uses it during
  per-record position remapping. In the safe runtime capture, row `+20`
  begins spans at values 0, 13, 26, and 39, matching the zero-based file
  offsets for this plain ASCII fixture.
- Row `+24` copies `state + 0x290 + 0x3c0 * state[+2]`, the `+0x290` word at
  the active-record-count index. `FUN_100180b0` adds one to this value before
  comparing it with `+0x64c` and reading intervening bytes from the original
  text buffer, establishing the same byte-coordinate domain and an inclusive
  endpoint for that consumer. Runtime row pairs `(0,10)`,
  `(13,23)`, `(26,36)`, and `(39,49)` match the inclusive spans of successive
  ASCII `Hello world` phrases. The producer's grouping rule and whether every
  output row must be interpreted as an entire source phrase remain open.
- Row `+28` is the sign-extended 16-bit `state+2` value. `FUN_10022dc0` sets
  it from the parser-output record count; `FUN_10012c70` and the remapping
  code also use it as the active record count. It is a count/index, not a
  phonetic label.
- Row `+32` is copied from synthesis-record `+0x28` (`psVar3[-3]` in the
  producer). This is not the loop ordinal: `FUN_1002c220` stores the ordinal
  separately at synthesis-record `+0x10`. For ordinary kind-2 records,
  `+0x28` is populated from the first byte of each 3-word record in the
  producer's interleaved array; the value is also used to index parallel
  synthesis-state arrays at `+0x47718`, `+0x47710`, and `+0x47714`. For
  inserted kind-1 records it is populated from caller `param_3`, which is
  used as an index into the `+0x47710` array. This establishes an array-index
  role in the producer, but not the upstream name or domain of that index. A
  separate normal-application capture observes indices 0–8, with each index
  reused by multiple kind-2 records; it also captures raw class bytes and the
  lookup results. A corrected call-sequence trace reads active record count
  from `*(short *)(*(int *)(engine+0x4c)+2)`: the 41-record batch has count 9
  and the following 25-record batch has count 6. An earlier trace mistakenly
  read `engine+2`; its reported zero count was a measurement error and is
  superseded.
  In the safe EX snapshot, row `+32` values 0 and 1 accompany source intervals
  `[0,4]` (`Hello`) and `[6,10]` (`world`); the following silence row retains
  value 1. `FUN_1002c530` also tests changes in the synthesis-record `+0x28`
  value as one condition for starting a new row. The full first-pass condition
  is now recovered from pointer offsets: start a row for the first record, a
  changed adjacent `+0x28` group/index value, or a changed adjacent `+0x27`
  record-kind byte. `psVar3[-0x1d]` and `psVar3[-3]` address previous/current
  `+0x28`; `psVar3-0x3b` and `psVar3-7` address previous/current `+0x27`.
  The condition does not compare the raw class byte at `+0x32`. Same-kind
  records with the same group/index reuse the current row. The row-field update
  within that reuse path computes and stores the most recently allocated row
  slot in record `+0x2e`: it uses `row_cursor - 1`, or `599` when the cursor is
  zero, matching the 600-slot circular row array. On a row-start path it stores
  the current row cursor before filling that slot and advancing/wrapping the
  cursor. These assignments are visible in `FUN_1002c530` at
  `0x1002c591–0x1002c59e` and `0x1002c645–0x1002c64f`; row creation advances the
  cursor at `0x1002c7b9–0x1002c7ce`. The same-kind reuse path compares the newly
  assigned row slot with the previous record's `+0x2e`, refreshing row `+12/+16`
  only when the slot changes. Thus `+0x2e` is a producer-maintained row-slot tag
  used to suppress redundant source-span refreshes, not an upstream synthesis
  field whose meaning must be inferred. The entry-time sequence capture shows
  the previous contents of this field before the producer rewrites it; its
  `0–4` progression is consistent with row-slot reuse but is not itself the
  post-call result. The active-record counts in that capture are 9 and 6 as
  corrected above. The remaining schema questions are the upstream
  names/domains for the `+0x28` group index, kind byte, and raw class byte, and
  finer phonetic meanings beyond the emitted phone mnemonics.

The nested selector is now traced through its complete static mapping.
`FUN_1002c220` writes a per-record class byte at synthesis-record `+0x32`
(kind-1 records use `0x5a`; kind-2 records read it from the per-state record
table). `FUN_1002c530` indexes `DAT_1007daa8` with that byte and stores the
mapped byte as the nested entry's 16-bit selector. The 0x60-byte table covers
indices `0x00..0x5f`; indices `0x5a..0x5e` all map to `0x28`, the silence
format already identified by the consumer. The full raw-index to output
mapping, expressed as inclusive ranges, is:

| Raw index | Nested selector | Raw index | Nested selector |
| --- | ---: | --- | ---: |
| `00–03` | `00` | `20` | `0d` |
| `04–06` | `01` | `21` | `0e` |
| `07–09` | `02` | `22` | `0f` |
| `0a–0c` | `03` | `23–25` | `10` |
| `0d–0f` | `04` | `26–28` | `11` |
| `10–12` | `05` | `29` | `12` |
| `13` | `06` | `2a` | `13` |
| `14` | `07` | `2b` | `14` |
| `15` | `08` | `2c` | `15` |
| `16` | `09` | `2d` | `16` |
| `17–19` | `0a` | `2e` | `17` |
| `1a–1c` | `0b` | `2f–31` | `18` |
| `1d–1f` | `0c` | `32–34` | `19` |
| `35` | `1a` | `36` | `1b` |
| `37` | `1c` | `38` | `1d` |
| `39` | `1e` | `3a` | `1f` |
| `3b–3d` | `20` | `3e–40` | `21` |
| `41` | `22` | `42` | `23` |
| `43` | `24` | `44` | `25` |
| `45` | `26` | `46` | `13` |
| `47` | `1a` | `48` | `1e` |
| `49` | `08` | `4a–4d` | `1e` |
| `4e–4f` | `08` | `50` | `1e` |
| `51–53` | `03` | `54–59` | `00` |
| `5a–5e` | `28` | `5f` | `00` |

The ordinary formatting branch also indexes a five-byte-stride NUL-terminated
label table at VA `0x1007bbd4` (file offset `0x7bbd4`). Its 39 entries are:

| Selector | Label | Selector | Label | Selector | Label |
| ---: | :--- | ---: | :--- | ---: | :--- |
| `00` | `AA` | `0d` | `F` | `1a` | `P` |
| `01` | `AE` | `0e` | `G` | `1b` | `R` |
| `02` | `AH` | `0f` | `HH` | `1c` | `S` |
| `03` | `AO` | `10` | `IH` | `1d` | `SH` |
| `04` | `AW` | `11` | `IY` | `1e` | `T` |
| `05` | `AY` | `12` | `JH` | `1f` | `TH` |
| `06` | `B` | `13` | `K` | `20` | `UH` |
| `07` | `CH` | `14` | `L` | `21` | `UW` |
| `08` | `D` | `15` | `M` | `22` | `V` |
| `09` | `DH` | `16` | `N` | `23` | `W` |
| `0a` | `EH` | `17` | `NG` | `24` | `Y` |
| `0b` | `ER` | `18` | `OW` | `25` | `Z` |
| `0c` | `EY` | `19` | `OY` | `26` | `ZH` |

Assembly passes the ordinary selector, label pointer (`0x1007bbd4 +
selector*5`), and nested duration to the literal format `%2d (%3s) : %d `.
These labels are the binary's phone mnemonics; that identifies the report
symbols, not a detailed phonetic interpretation of each upstream raw class.
Selector `0x28` bypasses the table and is passed with duration to the separate
literal format `%2d (sil) : %d `, establishing a silence entry. The producer
still has no visible bounds check on its raw class-table index; the nominal
0x00..0x5f input span follows from table size, not a producer guard.

After nested entries, `FUN_1001e0c0` copies the source substring using row
`+12/+16`, emits separator lines, and formats `%d - %d %s : %d ` with those
two offsets, the copied text, and row `+8` (total frames). The resulting
readable report line therefore carries source interval/text and total frame
count, preceded by phone/silence duration entries. Format strings and
argument order are direct static evidence. The public
`VT_TextToLipSyncLog_ENG` wrapper at `0x1001de50` calls `FUN_10022200`, whose
successful synthesis loop calls this formatter at `0x10022323`; runtime probes
confirm that API creates the report at the supplied path. This establishes
the public report producer and its output format, although no separate
downstream report reader was found. The offset interpretation and consumer
trace are in [`stage16-syncinfo-source-coordinate-consumers.txt`](../../tools/revkit/work/reports/stage16-syncinfo-source-coordinate-consumers.txt)
and by the row-source pseudocode in
[`stage16-syncinfo-row-producers.txt`](../../tools/revkit/work/reports/stage16-syncinfo-row-producers.txt);
the lookup bytes are at `vt_pau.dll` VA `0x1007daa8` (file offset `0x7daa8`).

The live capture snapshots all 14 header dwords and every populated row,
including nested values, on data polls 0–7. Header offsets `+0/+4/+8/+0xc`
are the row-array pointer, row capacity, per-row nested capacity, and row write
cursor. The remaining words form a pair of synchronized positions in the
audio timeline:

- `+0x10` is the global frame index at the beginning of the current output
  slice. It advances by `0x7530` (30,000 frames) per full poll in this capture.
- `+0x14/+0x18/+0x1c/+0x20` identify that slice's start within the SyncInfo
  rows: row index, frame offset in the row, nested-entry index, and frame
  offset in that nested entry.
- `+0x24/+0x28/+0x2c/+0x30/+0x34` identify the last frame in the slice: its
  inclusive global frame index, row index, frame offset in that row,
  nested-entry index, and frame offset in that nested entry.

This is a control-flow mapping, not a guess from field values. `FUN_100216d0`
converts the requested output byte count to frames (halving it for selector 0),
then walks each row's nested frame lengths and writes the endpoint fields.
`FUN_10021640` advances the start cursor to endpoint plus one; when it reaches
the nested entry's frame length it advances to the next entry, and then to the
next row. `FUN_100217e0` consumes the row range to associate EX marker
positions with source-text intervals. The live poll-0 values cross-check the
arithmetic: global end `59999` equals start `30000` plus 30,000 frames minus
one; row 7 / frame offset 799 / nested entry 0 / nested offset 799 also
identifies that same endpoint in the captured records. The initial snapshot
independently maps global frame 29,999 to row 3 / offset 3,396 / nested entry
2 / nested offset 1,076. Static decompilation and these captures are
[`stage16-syncinfo-cursor-functions.txt`](../../tools/revkit/work/reports/stage16-syncinfo-cursor-functions.txt)
and [`syncinfo-runtime-fields-safe2-1-api.log`](../../tools/revkit/work/stage16/syncinfo-runtime-fields-safe2-1-api.log).

Terminal result `1` frees the borrowed SyncInfo object; the safe trace skips
all pointer reads after that return. The runtime trace and replay are
[`syncinfo-runtime-fields-safe2-1-api.log`](../../tools/revkit/work/stage16/syncinfo-runtime-fields-safe2-1-api.log)
and [`run-syncinfo-runtime-fields.sh`](../../tools/revkit/work/stage16/run-syncinfo-runtime-fields.sh).
Evidence for the producer/consumer data sources is in the Ghidra pseudocode reports
[`vt_pau-core-pseudocode.c`](../../tools/revkit/work/reports/vt_pau-core-pseudocode.c)
and [`stage8-synthesis-functions.txt`](../../tools/revkit/work/reports/stage8-synthesis-functions.txt).

`VT_CopySyncInfo_New_ENG` copies all 13 scalar header dwords at offsets `+4`
through `+34`: the row and nested capacities, write cursor, and fields 4–13
(zero-based dword indexes). It does not replace the destination row-array
pointer at `+0`; the row copy also retains each destination nested pointer.
The runtime call used the native 600-row, 65-entry dimensions. The probe
compared all 11 header dwords at `+0xc..+34`, 4,800 row values (one 16-bit
value and seven dwords per row), and 78,000 nested values (one dword and one
16-bit value per entry); every comparison matched. The 600 destination nested
pointers and row-array pointer remained distinct. The 600/65 dimensions were
also observed being replaced by source dimensions in the separate shape probe.
The probe restored both objects' allocator dimensions before freeing them.
A separate edge trace called `(NULL,NULL)`, `(src,NULL)`, and `(NULL,src)`;
each returned to the caller without a fault. Self-copy retained distinctive
header, row, and nested-entry sentinels. For dimension mismatch, both objects
were allocated at native capacity, then source dimensions were set to 2 rows
and 3 nested entries. Copy set destination dimensions to 2×3, copied the
declared rows and nested entries, and left destination row 2 and nested entry
3 at their pre-copy sentinels. The source and destination nested pointers
remained distinct. This directly confirms that the source header controls the
copy loops and that smaller source dimensions leave destination tail storage
untouched. With source dimensions 0×3, the destination header changed to 0×3
while row 0 remained untouched. With dimensions 2×0, the destination header
changed to 2×0, its two fixed row records copied, and nested entry 0 remained
untouched. Thus row-copy iteration is controlled by the first source
dimension, nested-entry iteration by the second, and fixed row copying is
independent of nested width. Self-copy is an identity operation for a valid
object: the wrapper copies each header/data location back to itself, including
the pointed-to row and nested values; runtime sentinels agree. A follow-up set
destination metadata to 1×1 while retaining
native physical allocations, then copied a 2×3 source. The helper replaced
the destination dimensions with 2×3 and copied both rows and all three nested
entries, confirming it does not cap work to the destination's incoming
dimension words. With source dimensions -1×3, it copied those header values
and left row storage untouched. With 2×-1, it copied both fixed row records
but left nested storage untouched. Thus both loops use a signed-positive
condition; negative dimensions are not rejected or normalized. These
observations do not make undersized physical allocations safe. Static control
flow never compares source dimensions with destination physical capacity: it
copies source dimensions into the destination and loops using those source
values. Guard-page runtime probes now establish two exact fault boundaries.
With a source of 2 rows × 3 nested entries and a destination row array ending
immediately after one writable 36-byte row, the copy updates destination
dimensions to 2×3 and fully copies row 0, then faults writing row 1 at guard
address `0x1771000`, instruction `0x10026546`. With a 1×2 source and an
8-byte destination nested area ending at the guard page, the copy completes
entry 0, then faults writing entry 1 at `0x1771004`, instruction
`0x100265b9`. The Wine logs identify both as unhandled page faults on write;
the API does not return in either case. Thus the failure is a process-level
access violation at the first out-of-allocation write, after earlier header
and payload writes have already occurred. These probes ran in separate
disposable Wine processes. The shape and guard captures are
[`syncinfo-copy-shapes-api.log`](../../tools/revkit/work/stage16/syncinfo-copy-shapes-api.log),
[`syncinfo-copy-row-guard-api.log`](../../tools/revkit/work/stage16/syncinfo-copy-row-guard-api.log),
and [`syncinfo-copy-nested-guard-api.log`](../../tools/revkit/work/stage16/syncinfo-copy-nested-guard-api.log);
replay them with
[`run-syncinfo-copy-shapes.sh`](../../tools/revkit/work/stage16/run-syncinfo-copy-shapes.sh)
and [`run-syncinfo-copy-undersized.sh`](../../tools/revkit/work/stage16/run-syncinfo-copy-undersized.sh).

This copy contract is build-specific. The available 2006 `vt_eng.dll`
pseudocode in [`vt_eng-2006-all-functions-pseudocode.c`](../../tools/revkit/work/reports/vt_eng-2006-all-functions-pseudocode.c)
allocates 200 rows rather than 600. Its initializer also sets scalar header
index 4 (`+0x10`) to `-1`, but its copy helper does not copy that word; the
current Paul M16 `vt_pau.dll` helper does. This static difference means the
M16 copy contract should not be projected onto the 2006 build, or vice versa.
It does not by itself identify the omitted word's semantic meaning.

The header cursor fields `+0x1c..+0x34` are structurally mapped; any external
consumer use is absent from the supplied binaries. The remaining semantic
gaps are the grouping rule for row `+20/+24`, upstream raw-class/group-index
names for row `+32`, and the individual non-silence nested-selector meanings.
The physically undersized copy outcome is directly shown to fault at the first
write beyond row or nested storage; the copy is not capacity-safe. Full traces and
repeatable runners are [`syncinfo-fields-api.log`](../../tools/revkit/work/stage16/syncinfo-fields-api.log),
[`run-syncinfo-fields.sh`](../../tools/revkit/work/stage16/run-syncinfo-fields.sh),
[`syncinfo-copy-edges-api.log`](../../tools/revkit/work/stage16/syncinfo-copy-edges-api.log),
[`run-syncinfo-copy-edges.sh`](../../tools/revkit/work/stage16/run-syncinfo-copy-edges.sh),
[`run-syncinfo-copy-shapes.sh`](../../tools/revkit/work/stage16/run-syncinfo-copy-shapes.sh),
[`run-syncinfo-copy-undersized.sh`](../../tools/revkit/work/stage16/run-syncinfo-copy-undersized.sh),
and the live capture above.

### `VT_CheckUserDict_SourceNorm_ENG`

The public export at `0x1002a550` allocates a 52-byte local buffer at
`EBP-0x34`, initializes only byte 0, calls `FUN_1005f2e0(output,input)` at
`0x1005f2e0`, discards EAX, and returns `void`. A target-flow trace stopped at
`0x1002a567`, after the helper returns, to capture that raw helper EAX and the
scratch bytes. The input remained unchanged in all captured cases.

Runtime results for ordinary byte strings show that the helper copies bytes
without case folding or punctuation normalization: `Hello world` returns 11;
`Hello, world!` returns 13; `MiXeD_case-123` returns 14; `can't stop` returns
10; and UTF-8 bytes `c3 a9` in `caf\xc3\xa9` are retained, yielding length 5.
For `  Hello   world!  ` the output is `Hello   world!` (length 14), preserving
internal repeated spaces. The input `\t\r\nHello\r\n\t` yields `Hello` (5).
Both an empty string and a string consisting only of the trimmed bytes return
`-1`. Forty-nine `A` bytes return 49; 50 and 51 return `-5` and clear output
byte 0. Thus the observed maximum successful normalized length is 49 bytes.
The 52-byte scratch's remaining bytes after early error returns are not
initialized by the wrapper and must not be interpreted as output.

Static instructions show the helper skips leading and trims trailing ASCII
space, TAB, LF, and CR only. It rejects the tested two-byte patterns
`a1 a1`, `ae a1`, and `fd fe` with raw helper EAX `-3`, while `a1 a0` is copied
and returns 2. Its predicate at `0x10061f90` only checks that the next two
bytes are non-NUL; subsequent comparisons reject byte pairs in
`[a1-ad][a1-fe]`, `ae[a1-c2]`, and the single pair `fd fe`. The function
returns the copied byte count on success, `-1` when trimming leaves no bytes,
`-5` at the 50-byte limit, and `-3` for those byte checks in the tested
process. The names/intent behind `-3` and those excluded pairs remain unknown.
Replay with `run-userdict-source-normalizer-output.sh`; the complete capture is
`userdict-source-normalizer-output-api.log` and trace is
`trace-userdict-source-normalizer-output.gdb` in Stage 21.

A direct PE32 harness then called the same private helper at module RVA
`0x5f2e0` for all 65,536 two-byte buffer combinations, each followed by NUL.
It checked the raw return, normalized output bytes, terminator, and input
immutability against the statically observed trim and pair-rejection rules.
All 65,536 calls matched with zero mismatches: 1,257 returned `-3`, 276
normalized to empty and returned `-1`, 2,259 returned one copied byte, and
61,744 returned two copied bytes. The rejection count equals the complete
tested ranges `a1–ad × a1–fe` (1,222 pairs), `ae × a1–c2` (34), and `fd fe`
(1); trim bytes and leading NUL reduce which pairs reach those checks. Input
bytes remained unchanged. The log contains all 256 first-byte rows and their
per-row FNV-1a digests; the overall digest is `bc75caddeff3fd2b`. This closes
the two-byte buffer domain for the helper’s tested C-string interface, but
does not enumerate longer sequences, encoding semantics, or the undocumented
intent of `-3`. Build and run with
`build-source-normalizer-pairs.sh` and `run-source-normalizer-pairs.sh`.

A context extension placed each of the 1,257 rejected pairs after every
possible one-byte prefix and before every possible one-byte suffix (643,584
calls total). All returns, output prefixes, and input immutability checks
matched with zero mismatches. The prefix placement had 320,535 `-3` returns;
the 1,257 zero-byte prefixes terminate the C string before the pair and return
`-1`. All 321,792 suffix placements returned `-3`. When a nonempty prefix is
copied before a rejected pair, that partial prefix remains in the output
buffer; the helper does not append a NUL terminator before returning `-3`.
The wrapper initializes only output byte 0, so bytes after a partial prefix
remain uninitialized stack data. Treat the helper output as invalid on this
error even though its prefix bytes were written. The context-matrix digest is
`74330184b571cac7`; the capture has separate prefix/suffix digests. This tests
every one-byte context around each rejected pair, not all possible longer
strings. Build and run with `build-source-normalizer-pair-contexts.sh` and
`run-source-normalizer-pair-contexts.sh`.

### `VT_CheckUserDict_TargetNorm_ENG`

The export at `0x1002a570` returns the signed 16-bit result from
`FUN_1005f3b0` at `0x1005f3b0`. A target-flow matrix observed empty input
`-1`; `hello`, `HELLO`, `Hello World`, `hello-world`, and the exact string
`[SKIP]` each returned `1`. Leading/trailing spaces are trimmed in the caller's
buffer: `  hello   ` becomes `hello` and returns `1`. `hello[CI]` returns `2`
and changes the input to `hello`; `[CI]` alone and `[OTHER]` return `-11`.
`hello, world!`, `#`, `<`, and `<vtml_sub>` return `-2`. Fresh-buffer
repeated-`A` calls of lengths 63–65 return `1`; lengths 66–68 return `-5`.
The over-limit calls leave the leading bytes and the NUL at the original
length unchanged, so this tested path reports the limit without clearing or
truncating the caller's string.

A fresh-buffer sweep of every printable ASCII byte from `0x20` through `0x7e`
and the remaining non-NUL byte values from `0x01` through `0x1f` and `0x7f`
through `0xff` completes the single-byte map. Empty/NUL input, TAB, LF, CR,
and space return `-1`; all 52 ASCII letters return `1`; `[` returns `-11`;
every other value returns `-2` (198 one-byte values). The caller buffer
remains byte-identical in the printable sweep. The nonprintable sweep covered
all 160 values; only TAB/LF/CR return `-1`, while DEL and every high byte
return `-2`. This is exhaustive for single-byte inputs, not byte sequences or
multibyte encodings.

A second fresh-buffer sweep tested `A`, then 1–35 ASCII spaces, then `A`.
One through nine internal spaces returned `1`; ten through 35 returned `-7`.
The first rejection therefore occurs at ten consecutive separator bytes for
this two-letter input. This does not establish whether `-7` represents a
general token limit or how mixed whitespace, longer tokens, or parser markers
affect the same internal counter.

A third matrix tested one through 12 single-letter tokens separated by one
space. One through ten tokens returned `1`; 11 and 12 returned `-7`. Together
with the repeated-space result, this confirms the boundary for ordinary
single-space-separated tokens and shows repeated separators can reach the
same rejection earlier. The static counter path is consistent with a maximum
of ten counted segments, but whether empty segments from repeated spaces are
counted as words is still an interpretation of that implementation.

A direct branch matrix reached additional return paths. `<AB>` returned `1`,
while unterminated `<AB` returned `-8`; `<A B>` returned `-4` at the space
inside the open angle-bracket state. A 30-byte `A` segment followed by space
and `B` returned `1`, while 31 `A` bytes followed by the same separator
returned `-6`. Inserting each of `a1 a1`, `ae a1`, and `fd fe` after an ASCII
`A` returned `-3`. The caller buffers remained unchanged in these cases. The
observed sequences exercise the static branches, but do not establish the
intended markup language or meanings of the numeric results.

A marker matrix further tested exact placement and mutation. Leading
`[SKIP] ` is trimmed to `[SKIP]` and returns `1`; lowercase `[skip]` also
returns `1`, while `[SKIP]x` returns `-11`. `A[CI]`, `A[ci]`, and `A [CI]`
return `2`; the export writes NUL over the opening bracket, preserving the
space in the separated form (`A `). `A[OTHER]` and `A[CI]B` return `-12` and
remain unchanged. An exhaustive ASCII letter-case mask sweep found all 16
case variants of `[SKIP]` return `1` and all four case variants of the `[CI]`
suffix return `2`. This establishes case-insensitive matching over the marker
letters, exact leading `[SKIP]` matching after trim, and `[CI]` at the end of
the tested strings. Other placements, intervening text, and malformed
bracket combinations remain untested. Reproduce the capitalization sweep with
`run-userdict-target-normalizer-marker-case.sh`.

Static code at `0x1005f3b0` first trims through `FUN_10063230`, compares a
leading `[` string against the data constant `[SKIP]` at `0x1009c524`, and
handles `[CI]` at `0x1009c51c` in a later suffix branch that can terminate the
input at the bracket. The internal scan allows 65 bytes before its `>0x41`
length check; the fresh-buffer runtime sweep confirms 63–65 return `1` and
66–68 return `-5` for repeated `A`, with caller bytes unchanged. This pins the
boundary for that valid ASCII input class, not every multibyte, trimmed, or
marker-bearing form. The separator sweeps show `A` + 1–9 spaces + `A` returns
`1` and 10–35 spaces returns `-7`; one-space-separated sequences of 1–10
single-letter tokens return `1`, and 11–12 return `-7`. The static counter is
consistent with a ten-segment maximum, but whether repeated spaces count as
empty segments is not proven semantically. The behavior of `-7` outside these
ASCII patterns remains unknown. The branch matrix reaches `-3`, `-4`, `-6`,
and `-8` with the specific byte-pair, open-tag-space, long-segment, and
unterminated-tag inputs described above. Other parser-state domains remain
open. Replay with `run-userdict-target-normalizer-lengths.sh`,
`run-userdict-target-normalizer-spaces.sh`,
`run-userdict-target-normalizer-words.sh`, and
`run-userdict-target-normalizer-branches.sh`, and
`run-userdict-target-normalizer-marker-case.sh`; traces and captures are in
Stage 21. The initial matrix and printable-ASCII sweep use
`run-userdict-target-normalizer-matrix.sh` and
`run-userdict-target-normalizer-ascii.sh`.

Separate single-input observation: `VT_CheckUserDict_TargetPhon_ENG("example")`
returned `-9`; its exhaustive tested spelling inventory is documented in the
[target-phoneme validator section](#target-phoneme-validator-inventory). This
does not assign broader semantics to that return code.

`VT_SetEmphasisFactor_ENG` clamped slot 1 values `200` and `-200`
to `95` and `-95`. `VT_SetTextTypeForHighlight_ENG(7)` stored `1`;
`VT_SetParenthesisCharNumber_ENG(-1)` and
`VT_SetEnglishReadingRule_KOR(-1)` each left their respective fields at `0`.
`VT_SetSoundCardID_ENG` returned the prior `0xffffffff` value after a
temporary write. The probe restored every field it changed before normal
execution continued. `VT_SetUnitSelectHistoryMode_ENG` calls with `1` and `0`
left its global state at `0` while the loaded flag was `1`, matching the
decompiled guard that ignores calls after model loading.

A follow-up exercised setter boundaries in the same loaded slot. Emphasis
inputs `INT_MIN`, `-96`, `-95`, `-94`, `0`, `94`, `95`, `96`, and `INT_MAX`
stored `-95`, `-95`, `-95`, `-94`, `0`, `94`, `95`, `95`, and `95`. Parenthesis
count and reading-rule inputs `INT_MIN`, `-1`, `0`, `1`, and `INT_MAX` stored
`0`, `0`, `0`, `1`, and `INT_MAX` for each setter. The probe restored the
initial values after each sweep. These setters write the loaded Paul process's
global configuration block through `DAT_100a0460` (offsets `+0x2041c` and
`+0x20420`), not a caller-selected speaker field. The positive-value effects
remained open at that point. See the Stage 21
`scalar-helper-edges-api.log` capture and `run-scalar-helper-edges.sh` runner.
For a direct consumer check, the Stage 21 trace loaded Paul, reached
`VT_TextToFile_ENG` for the format-4 text `Hello world. I read 123 books
(three times) at 10:30 in the U.S.A.`, set the global field to 0 or 1, and then
armed a hardware read/write watchpoint on that four-byte field until the API
returned. Both calls returned 1, neither run accessed the field after the
setter, and the WAV outputs were byte-identical (SHA-256
`a54bcb6ea61bbb268dea4404a62bc6b388be82836c057d47a305e2fdabc89739`). The
reviewed `vt_pau` pseudocode and disassembly show the setter and its clamp
read, but no other reference to this offset. This shows no consumption for
this loaded-Paul format-4 utterance; it does not establish effects for other
text, formats, voices, or host-side uses. Captures, outputs, and runner are
`english-reading-rule-{0,1}-api.log`,
`english-reading-rule-{0,1}.wav`, and
[`run-english-reading-rule-effect.sh`](../../tools/revkit/work/stage21/run-english-reading-rule-effect.sh).

An expanded format-4 matrix reused the same setter/watchpoint sequence for six
texts: the original number/time sample, two homograph/heteronym samples, a
date/currency/unit sample, an abbreviation/time sample, and an acronym sample.
For all six, modes 0 and 1 returned 1, the field stayed unread/unwritten after
the setter through API return, and the WAVs were byte-identical. On the
date/currency sample, `INT_MAX` was also preserved by the setter; its output
matched mode 0 and the field had no post-setter accesses. This is seven paired
output comparisons and thirteen successful calls in fresh processes. Literal
references to offset `+0x20420` in the reviewed Paul disassembly/pseudocode
are confined to the setter's write/clamp path. Together these results give no
evidence of a consumer on these loaded-Paul `VT_TextToFile_ENG` paths; they do
not cover other formats, voices, APIs, or host-side access. The matrix runner,
result file, captures, and WAVs are under Stage 21 as
`run-english-reading-rule-matrix.sh`,
`english-reading-rule-matrix-results.txt`, and the
`english-reading-rule-matrix-*` artifacts.

Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-english-reading-rule-matrix.sh
```

The same capture swept `VT_SetSoundCardID_ENG` through `INT_MIN`, `-1`, `0`,
`1`, and `INT_MAX`: each value was stored verbatim, and each call returned the
immediately prior value. The initial global value `-1` was restored. These
boundary calls do not establish which device IDs are valid to the playback
subsystem.

At the pre-load breakpoint, `VT_SetUnitSelectHistoryMode_ENG` was then called
for every byte value. With the load gate at 0, input 1 alone stored 1; the
other 255 byte values stored 0. All 256 comparisons matched and the prior value
was restored. A later breakpoint at `VT_TextToFile_ENG` confirmed that the
ordinary loader completed and the gate had become 1. This exhausts the byte
domain for the pre-load setter path; the already observed post-load calls with
0 and 1 leave the value unchanged. The flag's detailed effect on synthesis
remains outside this input-state sweep. See
`unit-history-byte-matrix-api.log` and `run-unit-history-byte-matrix.sh` in
Stage 21.

An exhaustive byte sweep of `VT_SetTextTypeForHighlight_ENG` then called
values 0 through 255 in a loaded Paul process. Every call stored exactly 0 for
input 0 and 1 for inputs 1–255 in slot 1's field at `+0x20424`; all 256
comparisons matched, with no other stored values. The probe restored the
initial value 0. This recovers the setter's byte-to-state normalization. A
separate paired synthesis called format 4 with the same `Hello world.` text
and arguments at flag 0 and flag 1. Both calls returned 1; the two WAVE files
were byte-identical (23,650 bytes, SHA-256
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`). Thus
this tested flag change did not affect the produced audio. Static pseudocode
does show a downstream effect on source-position bookkeeping. In
`FUN_10022dc0`, the parser builds inclusive per-record endpoints at state
`+0x64c/+0x650` (the `+0x650` end is inclusive when the interval has multiple
positions). After `FUN_10022970`, `FUN_10022dc0` tests the normalized flag at
global-state `+0x20424` (disassembly address `0x10022ee4`). With the flag set,
each endpoint is remapped through per-state tables `+0x11de5/+0x11de6` and a
shared table `+0x11de4`; with it clear, endpoints are clamped to the
parsed-record count and then mapped through the endpoint-specific tables.
The transformed endpoints are subsequently copied by `FUN_1002c530` into
SyncInfo row `+20/+24`, which the source-coordinate consumers and EX marker
mapper use to select timeline rows.

The direct marker-record mapping is in `FUN_1001c990`. It builds three
per-context dword arrays at `+0x47790/+0x47794/+0x47798`; its `+0x47790`
construction expands a source character classified by `FUN_1001c900` to two
output positions, while the other two arrays begin as identity maps and are
updated through the text-rewrite stages. After those stages, flag 1 rewrites
each marker record's `+8` through `+0x47794` and then `+0x47790`; flag 0 uses
`+0x47794` alone. This directly maps the setting to the source-position path
for returned EX marker records. `FUN_1001c900` accepts exactly two non-NUL
bytes in these ranges: first byte `0xa1..0xad` with second byte `0xa1..0xfe`,
`0xae` with second byte `0xa1..0xc2`, or the one pair `0xfd 0xfe`. A runtime
direct-call sweep covered every one of those 1,257 accepted pairs; each
returned 1. Thus the classifier's positive domain is runtime-confirmed, while
the complement is established from the complete static branch conditions.
The external caller's interpretation of the final coordinates remains open.

A paired selector-0 EX capture now tests both settings on
`<vtml_mark name="start"/>Hello world.<vtml_mark name="end"/>`. Both runs
returned initial status 0 with 23,606 bytes, two descriptor rows, and two
SyncInfo rows. The descriptor's `+0/+4/+8` values, names, kind bytes, and all
captured SyncInfo row fields are identical: marker `start` has `+8=0`, marker
`end` has `+8=37`, and SyncInfo text spans are `[25,29]` and `[31,35]`. The
512-byte inline fields are not identical after their NUL-terminated names:
with flag 0 the first row contains a little-endian dword run `0..21` near the
field's end and the second contains `27..60` just after `end\0`; with flag 1
the first contains `0..60` beginning at row offset `+0x10c`, while the
second contains a different byte sequence. A second flag pair on the 1,091-byte
boundary fixture likewise changes post-NUL field bytes; its first row has
`0..21` at flag 0 versus `19..144` at flag 1. These bytes are inside the
returned descriptor, but `FUN_1001ccc0` initializes only row `+0/+4/+8`, and
`FUN_1002efd0` writes only the name and its terminator in this field. The
differing post-NUL bytes are therefore uninitialized row storage as far as the
recovered writer path establishes; they may reflect allocator reuse and are
not treated as meaningful flag-dependent payload. The visible row coordinates
and SyncInfo match for the ASCII fixtures, where the conversion map is
identity. A boundary-case input then placed six accepted pairs between seven
named marks: `a1 a1`, `a1 fe`, `ad a1`, `ae a1`, `ae c2`, and `fd fe`. Both
flag settings returned initial status 0, seven kind-1 descriptor rows, and one
SyncInfo row under every selector. Encoded output length was 13,936 bytes for
selector 0 and 6,968 bytes for selectors 1 and 2; each SyncInfo row reported
`+8=0x1b38` (6,968 frames). SyncInfo source/end fields `+12/+16` and mapped
endpoint fields `+20/+24` were `0x2e/0x2f/0x2e/0x2f` with flag 0 and
`0x2d/0x2d/0x2d/0x2d` with flag 1, consistently across selectors. Thus the
setting changes both returned marker positions and SyncInfo source-coordinate
fields for these classified pairs while leaving the frame total unchanged.
The marker coordinates matched across selectors:
flag 0 reported marker `+8` positions
`0,24,48,72,96,120,144`; flag 1 reported
`0,23,46,69,92,115,138`. The difference from the first marker grows by
exactly one for each preceding classified pair. This is consistent with the
flag-1 path accounting for each accepted two-byte source character as one
position in the marker map; the supplied caller does not establish the
consumer's coordinate convention or visible highlighting meaning. Replay
with `run-highlight-ex-two-byte-classes-replay-pair.sh` for selector 0 and
`run-highlight-ex-two-byte-classes-matrix.sh` for selectors 1 and 2. The
fixture is `input-ex-raw-highlight-two-byte-classes.txt`; flag-paired logs
and raw descriptor arrays are `highlight-ex-two-byte-classes-replay-{0,1,2}-{0,1}-api.log`
and matching `-descriptor.bin` files in Stage 21. The earlier single `a1 a1` fixture
returned no output bytes and is not needed for this completed audio result.

Replay the ASCII runs with `run-highlight-ex-raw-mark.sh` and
`run-highlight-ex-raw-records.sh`; captures and raw descriptor arrays are in
Stage 21. The original `Hello world.` format-4 audio equality covers only its
tested text/settings. Earlier v4 trace captures did reach the API breakpoint;
they are not negative evidence. Existing byte-matrix/audio probes are
`run-highlight-byte-matrix.sh` and `run-highlight-audio-effect.sh`.

`VT_DestroyWindow_ENG` is a one-call wrapper around `DestroyWindow` at
`0x10027fd6`. In target-flow execution, the sample had stored handle `0x10074`
and reached that import call; the inferior exited normally before the
post-import instruction at `0x10027fdc` could be observed. This captures the
call's input and the host-process outcome, but not the Win32 Boolean result or a
wrapper return value. It does not establish that other hosts will exit when
their window is destroyed. Replay with `run-destroy-window-direct.sh`; the
trace and log are `trace-destroy-window-direct.gdb` and
`destroy-window-direct-api.log` in Stage 21.

A controlled follow-up redirected target execution to the same wrapper after
replacing its process-local stored handle with `NULL`; it then repeated with
`0xdeadbeef`. At the instruction after the Win32 call (`0x10027fdc`), raw EAX
was 0 for both inputs. The original handle `0x10074` was restored before
continuing, and the capture confirms that restoration. The wrapper is
decompiled `void`, so this is the called Win32 function's raw result, not a
declared export return. The inferior exited with code 1 later in this modified
flow. Reproduce with `run-destroy-window-invalid-handle.sh`; the capture and
trace are `destroy-window-invalid-handle-api.log` and
`trace-destroy-window-invalid-handle.gdb` in Stage 21. The valid-handle call
still does not expose its post-call result.

### Target-phoneme validator inventory

The direct export at `0x1002a590` was called with engine-heap strings. The
single-byte sweep covered TAB, LF, VT, FF, CR, space, and printable ASCII
`0x21..0x7e` (100 values total): 16 uppercase letters and `#` returned `1`;
TAB/LF/CR/space returned `-1`; `[` returned `-2`; and the remaining 78
values returned `-9`. The accepted uppercase single-character strings were
`B D F G K L M N P R S T V W Y Z`. The `#` result follows a separate static
branch in `FUN_1005f5c0` and is recorded as a special accepted token, not
classified as a phone.

A Stage 21 sweep completed the single-byte domain, calling the same export
for every byte `0x00..0xff` followed by NUL. NUL, TAB, LF, CR, and space
returned `-1`; the 16 uppercase tokens plus `#` returned `1`; `[` returned
`-2`; every other byte returned `-9` (233 values total). No other one-byte
high-bit value is accepted. This closes the byte domain only; it does not
test multibyte encodings or change the uppercase-token finding. The complete
capture is
[`targetphon-byte-domain-api.log`](../../tools/revkit/work/stage21/targetphon-byte-domain-api.log);
reproduce with `run-userdict-targetphon-byte-domain.sh` in Stage 21.

The complete uppercase two-letter cross-product (676 calls) accepted exactly
`CH DH HH JH NG SH TH ZH`; the other 668 returned `-9`. The uppercase
two-letter cross-product with suffixes `0`, `1`, and `2` (2,028 calls)
accepted 45 strings: each of `AA AE AH AO AW AY EH ER EY IH IY OW OY UH UW`
with each suffix. A second 150-call sweep tested those 15 stems with digit
suffixes `0..9`; only suffixes `0..2` returned `1`, and suffixes `3..9`
returned `-9`. This maps the validator's tested spelling inventory; it does
not assign acoustic meanings to the codes.

A 244-call case-mask matrix tested both cases for all 16 accepted one-letter
tokens, all four letter-case masks for each of the eight accepted consonant
pairs, and all four masks for each accepted vowel stem plus suffix `0`–`2`.
Exactly the 69 canonical uppercase spellings returned `1`; each of the 175
lowercase or mixed-case variants returned `-9`. This establishes
case-sensitive acceptance for the known phone-token inventory. The direct
validator accepts case-folded `[CI]`; in loaded `P` rows, only the uppercase
marker forms tested so far enable case-insensitive source matching (see below).
Unknown phone tokens remain unmapped.
Capture: [`targetphon-case-matrix-api.log`](../../tools/revkit/work/stage21/targetphon-case-matrix-api.log);
reproduce with `run-userdict-targetphon-case-matrix.sh` in Stage 21.

Sequence probes returned `-1` for null, empty, and whitespace-only strings;
`1` for `HH B AA1 CH`, repeated spaces, and TAB-separated
`HH`/`B`; `-9` for `HH EXAMPLE`; and `-2` for isolated `[`. A dedicated
repeated-`B` sweep accepted 120–130 tokens (130 tokens, 259 input bytes,
returned `1`) and returned `-5` for 131–140 tokens (131 tokens, 261 input
bytes, returned `-5`). A fresh-buffer probe also checked 130 tokens plus a
trailing space (260 original bytes): trimming reduced it to 259 bytes and it
returned `1`; a fresh 261-byte/131-token string returned `-5`. Static code
increments a consumed-byte counter and returns `-5` when it reaches `0x104`
(260). The earlier mixed-call probe's nominal 261-byte case was not fresh:
its prior call trimmed a trailing space in the same writable buffer, and the
runner failed to restore the NUL-terminated byte before writing the last
token. It therefore retested the 130-token prefix and returned `1`. The
apparent inconsistency is a probe-construction artifact; the ASCII boundary is
cross-checked at 259 accepted versus 261 rejected bytes after trimming.

The bracket path is now bounded separately. `FUN_1001c2c0` at `0x1001c2c0`
is a case-folded string comparison (it indexes a 16-bit character-fold table);
it does not search for a substring. Before validation,
`FUN_10063230` calls trim helpers that remove leading and trailing space, TAB,
CR, and LF in place. The validator treats `[` as a token boundary, then
compares the remaining text to the literal `[CI]`. Runtime calls accepted
`HH[CI]`, `HH [CI]`, `HH [ci]`, `HH [Ci]`, repeated-space and TAB-separated
forms with return `2`; `[CI]` without a preceding valid phone returned `-2`.
Trailing spaces/TAB are trimmed before this comparison. A following phone or
second marker (`HH [CI] HH`, `HH [CI][CI]`) returned `-12`; `[SKIP]` after a
phone also returned `-12`, while isolated `[SKIP]` returned `-2`.

The export mutates the supplied writable string. Runtime inspection after the
call showed leading/trailing whitespace removed (`\t HH \r` became `HH`), and
the accepted marker's opening `[` replaced by NUL (`HH[CI]` became `HH`, while
`HH [CI] ` became `HH `). Invalid marker strings remained present after the
call, apart from the common trimming. This is a buffer-mutation contract, so a
read-only string literal is not a safe argument. The marker's intended semantic
meaning cannot be read from this export alone; the loaded-dictionary runtime
effect is established below.

The dictionary file loader provides a separate constraint on that meaning.
In the decompiler pseudocode, `VT_LOAD_UserDict_ENG` at `0x10021400` calls
`FUN_1003bd50` to parse the supplied file. That parser recognizes rows whose
third field is exactly the two-byte code `A` or `P`. Its `P` path validates the
target field with `FUN_10056790` and builds its byte sequence with
`FUN_10056890`; neither helper calls exported `VT_CheckUserDict_TargetPhon_ENG`
(`0x10022bc0`) or contains an explicit `[CI]` comparison. The private
converter splits on space/TAB (and its inner loop also stops at CR/LF), then
maps each token through the internal phoneme table, with a special `#` case.
The top-level pseudocode contains no literal `[CI]` comparison, so it does not
show which callee implements the marker. The successful runtime comparison
establishes its effect for two `P`-row spellings. With text `HELLO`, plain row
`hello,HH,P` produced PCM byte-identical to the no-dictionary control
(7,798 frames). `hello,HH[CI],P` and `hello,HH [CI],P` each loaded/unloaded
with status `1` and both changed synthesis to the same 1,150-frame PCM as the
plain `hello,HH,P` row produces for lowercase `hello`. This isolates the
marker as a case-insensitive source-match annotation for the tested `P` rows;
an additional mixed-case `HeLLo` run with the adjacent marker produced the
same changed PCM while its plain `P` row remained identical to control. The
marker does not alter the target-phone output in these samples. Captures are
`james-licensed-userdict-ci-source-case-v1-v3-api.log` and
`james-licensed-userdict-ci-source-case-spaced-v1-api.log`, with corresponding
WAVs and replay scripts in Stage 21; the mixed-case capture is
`james-licensed-userdict-ci-source-case-mixed-v1-api.log`. Other ASCII word shapes, non-ASCII case folding, and interaction with other
target tokens remain open. All four marker case masks were also tested on
otherwise equivalent `P` rows: each loaded/unloaded with status `1` and enabled the same
uppercase `HELLO` and mixed-case `HeLLo` matching effect as `[CI]`. On lowercase
`hello`, it produced the same 1,150-frame PCM as plain `hello,HH,P`. This
confirms all four letter-case masks (`[CI]`, `[ci]`, `[cI]`, `[Ci]`) as
case-insensitive source-match annotations in the tested `P` row. This agrees
with the direct export's case-folded marker recognition. Captures:
`james-licensed-userdict-ci-source-case-lower-tag-v1-api.log` and
`james-licensed-userdict-ci-source-case-lower-tag-lower-source-v1-api.log`,
`james-licensed-userdict-ci-tag-cI-case-v2-api.log`, and
`james-licensed-userdict-ci-tag-ci-case-v3-api.log`.

Reproduce with `run-userdict-targetphon-char-sweep.sh`,
`run-userdict-targetphon-pair-sweep.sh`,
`run-userdict-targetphon-stress-sweep.sh`,
`run-userdict-targetphon-stress-digit-sweep.sh`,
`run-userdict-targetphon-marker-sweep.sh`,
`run-userdict-targetphon-marker-boundaries.sh`,
`run-userdict-targetphon-marker-inplace.sh`,
`run-userdict-targetphon-sequence-errors.sh`,
`run-userdict-targetphon-token-count-boundary.sh`,
`run-userdict-targetphon-length-fresh.sh`, and
`run-userdict-targetphon-order-check.sh`. Converter mapping and bounds are
reproduced by `run-userdict-targetphon-converter-map.sh`,
`run-userdict-targetphon-converter-sequence-check.sh`, and
`run-userdict-targetphon-converter-limit.sh`. Captures in Stage 16 use the
corresponding `userdict-targetphon-*-api.log` names. The code path uses `-12`
for a nonmatching bracket suffix after a valid phone; it is not a substring
search or general markup parser.

The adjacent private converter `FUN_1005f710` takes an output byte buffer and
a target-phone string. Runtime calls for each of the 69 accepted spellings
above and `#` returned `1`, wrote one byte, then NUL. Multi-token input emitted
the same bytes in token order. This recovers the converter's spelling-to-byte
mapping without assigning meaning to the byte IDs:

| Input | Emitted byte (hex) |
| --- | --- |
| `B D F G K L M N P R S T V W Y Z` | `13 15 20 21 2a 2b 2c 2d 35 36 37 39 41 42 43 44` |
| `CH DH HH JH NG SH TH ZH` | `14 16 22 29 2e 38 3a 45` |
| `AA` (`0`, `1`, `2`) | `01 02 03` |
| `AE` (`0`, `1`, `2`) | `04 05 06` |
| `AH` (`0`, `1`, `2`) | `07 08 09` |
| `AO` (`0`, `1`, `2`) | `0a 0b 0c` |
| `AW` (`0`, `1`, `2`) | `0d 0e 0f` |
| `AY` (`0`, `1`, `2`) | `10 11 12` |
| `EH` (`0`, `1`, `2`) | `17 18 19` |
| `ER` (`0`, `1`, `2`) | `1a 1b 1c` |
| `EY` (`0`, `1`, `2`) | `1d 1e 1f` |
| `IH` (`0`, `1`, `2`) | `23 24 25` |
| `IY` (`0`, `1`, `2`) | `26 27 28` |
| `OW` (`0`, `1`, `2`) | `2f 30 31` |
| `OY` (`0`, `1`, `2`) | `32 33 34` |
| `UH` (`0`, `1`, `2`) | `3b 3c 3d` |
| `UW` (`0`, `1`, `2`) | `3e 3f 40` |
| `#` | `64` |

The numeric values are emitted bytes, not labels inferred from table order.
The special `#` mapping remains unknown. The `[CI]` marker makes source
matching case-insensitive for the tested uppercase marker spellings in `P`
rows; lowercase `[ci]` did not enable that behavior in the tested pair. The
per-token calls used a 66-byte output allocation; ordered
multi-token output is captured in
`userdict-targetphon-converter-sequence-check-api.log`.

The converter accepts at most 65 output tokens. Direct calls with 64 and 65
space-separated `B` tokens returned `1`; the latter wrote code `0x13` at output
offset 64 and NUL at offset 65. A 66-token call returned `0` and cleared
output byte 0, while the already written byte at offset 64 remained. This
matches the static `ebx >= 0x41` rejection branch. Callers need room for 66
bytes when requesting the maximum successful output.

`VT_TextToBufferEX_ENG` is absent from the public `vt_eng.h` declaration. Its
exported wrapper pseudocode at `0x1001ddf0` accepts selectors 0–2 and dispatches
them to format handlers at `0x100200c0`, `0x10020460`, and `0x10020930`; other
selectors return `-1`. The 15-argument call shape is recovered from the x86
wrapper and handler call flow. Ghidra's pointer-looking types for the final two
arguments are not reliable: the handler uses both as integer options.

| Argument positions | Recovered role | Evidence and limit |
| --- | --- | --- |
| 1–7 | selector, text, output buffer, output-length pointer, flag, thread ID, speaker ID | Handler branches on selector/flag, bounds-checks the thread ID, normalizes speaker IDs outside 0–5 to speaker 1, and uses speaker-indexed stream state. |
| 8–9 | Optional output pointers | The handler initializes them to zero, then writes borrowed pointers from context offsets `+0x47774` and `+0x122448`. The first is the active SyncInfo object. The second is a 16-byte per-text descriptor with count, row pointer, and selected row-index bounds. Its `0x210`-byte rows record mark source position, chunk-local/global audio-frame coordinates, inline name bytes, and kind: 1 named, 2 unnamed, 3 overlength named. Both pointers expire at terminal cleanup; copy needed data before the `DONE` return. |
| 10–13 | pitch, speed, volume, pause | Passed into synthesis-option setup in the same order as the ordinary buffer API. Negative pitch/speed/volume select engine defaults; speed 0 becomes 50; a later normalizer clamps effective pitch/speed/volume to 50–200, 50–400, and 0–500. Pause is stored directly and later consumed by the per-unit pause path. A text type of 4 or 6 resets all four to defaults. |
| 14 | dictionary index | Integer slot 0–1023 selects a loaded dictionary for the speaker when present; an empty slot, negative index, index >=1024, or unavailable slot falls back to the default dictionary. |
| 15 | text type | Negative values become 0; nonnegative values are truncated to a byte. Values 4 and 6 reset the four preceding options. Other byte values are stored without a range check in this setup helper. |

The compact probe calls passed speaker 1, thread 0, `-1` for positions 10–13,
dictionary index 0, text type 0, and non-null guarded output pointers at
positions 8–9. Those pointers received the two internal addresses above. The
initial short-text run still matches ordinary buffer formats byte-for-byte:
selector 0 matches format 0 (`e232e245…f61582`), selector 1 matches format 1
(`9880164f…d5dfc`), and selector 2 matches format 2
(`424afda0…e53ab`). That comparison establishes output equivalence only for
the tested input/options.

For a 78-byte repeated `Hello world. ` utterance and 60,000-byte guarded
storage, all three selectors enter the same polling contract. The result field
is `0` (`PROCESSING`) for every data-bearing call; a separate final flag-1 poll
returns `1` (`DONE`) with length zero:

| Selector | Data chunks in call order | Total bytes | Terminal call |
| ---: | --- | ---: | --- |
| 0 | 60,000; 60,000; 60,000; 60,000; 49,636 | 289,636 | fifth poll: return 1, length 0 |
| 1 | 30,000; 30,000; 30,000; 30,000; 24,818 | 144,818 | fifth poll: return 1, length 0 |
| 2 | 30,000; 30,000; 30,000; 30,000; 24,818 | 144,818 | fifth poll: return 1, length 0 |

The direct state probes also establish for each selector: a second flag-0
start while active returns `-7` and leaves the caller's length sentinel
unchanged; flag 2 cancels with return `1` and length 0; flag 1 after cancel
returns `-2` and leaves its length sentinel unchanged. Null text, empty text,
and null output buffer return `-3`, `-4`, and `-5`, respectively. After
unloading speaker 1, a valid call returns `-6` and writes length 0. Thread ID 1
returns `-2` in this Wine runtime without changing length or either guard
region. Earlier invalid-selector probes returned `-1`. For every long run,
the four leading and four trailing `0xa5` sentinels around the 60,000-byte
physical buffer remained intact, as did the guard bytes around both optional
pointer-output slots. A flag-1 call with no active state also returns `-2` in
the earlier short probe.

The runtime record does not establish nonzero-thread support on another host.
The negative-flag size query is directly observed for all three selectors: it
returns `60,000` for selector 0 and `30,000` for selectors 1 and 2, and writes
that value through `output_len` even with null text/buffer, thread 77, and
otherwise invalid scalar arguments. At `0x10021c2b–0x10021c3d`, disassembly
places this early return before database, thread, text, buffer, and
optional-output-pointer checks.
On normal calls, the handler never reads `*output_len` as a capacity. With
`*output_len` initialized to 1 and a physically writable 60,000-byte region,
selector 0 returned 60,000 bytes and selector 1 returned 30,000. Selector 2
returned 8,984 in its first capture, then 30,000 in two repeat runs with the
same setup; both exceed the incoming value. The 30,000-byte figure is its
observed maximum, while the cause of the shorter first result is unresolved.
The guard after 60,000 bytes stayed intact, and selector 1 left the region
after byte 30,000 intact. A separate guard-page probe passed selector 0 a
buffer with exactly one writable byte followed by a `PAGE_NOACCESS` page and
`*output_len = 1`. Selectors 0, 1, and 2 each faulted on a write to the first
protected address (`0x1771000`) at DLL instruction `0x10064000`, a `rep movs
DWORD PTR es:[edi],DWORD PTR ds:[esi]`. Wine reported an unhandled page fault
on write; the API could not return. This is direct runtime confirmation that
the output length does not bound writes even when the actual accessible
region is one byte, across all supported encodings. A parameterized follow-up
repeats the guard-page test with 16 writable bytes (selector 0) and 4,096
writable bytes (selectors 0–2); each run faults on a write to the immediately
following protected page, with `*output_len` still initialized to 1. The 16-byte
case uses `Hello world.`; the 4,096-byte cases use the 511-byte-mark fixture.
Thus tested undersized spans of 1, 16, and 4,096 bytes all fail before return,
and the API does not use the input length value as a capacity. Because the
interface accepts no capacity, the conservative call contract is to allocate
the fixed per-call limit: 60,000 bytes for selector 0 or 30,000 for selectors
1/2. The protected write occurs in `FUN_10063f30`'s `rep movs` copy path
(instruction `0x10063f63` or `0x10064000`, depending on copy alignment).
Reproduce the size sweep with
[`run-buffer-ex-guard-size.sh`](../../tools/revkit/work/stage16/run-buffer-ex-guard-size.sh)
and its [`trace template`](../../tools/revkit/work/stage16/trace-buffer-ex-guard-size.gdb.in);
logs are `buffer-ex-guard-size-0-16.log` and
`buffer-ex-guard-size-{0,1,2}-4096.log`. The one-byte runs remain in
[`buffer-ex-guard-page.log`](../../tools/revkit/work/stage16/buffer-ex-guard-page.log),
[`buffer-ex-guard-page-1.log`](../../tools/revkit/work/stage16/buffer-ex-guard-page-1.log),
and [`buffer-ex-guard-page-2.log`](../../tools/revkit/work/stage16/buffer-ex-guard-page-2.log).
Callers need writable storage for the selector maximum per call (60,000 or
30,000 bytes); the returned length is available only after the write. Reproduce
the one-byte test with
[`run-buffer-ex-guard-page.sh`](../../tools/revkit/work/stage16/run-buffer-ex-guard-page.sh).

The optional outputs at positions 8 and 9 are null-checked before their
initial zero stores, so null is accepted. Any non-null value is written
through directly; there is no pointer-validity check. Runtime calls with
position 8 equal to address `1` and position 9 equal to address `1` each
terminated the Wine process with Windows status `0xc0000005`
(`STATUS_ACCESS_VIOLATION`) before the injected API call returned. The logs
record GDB's octal exit code `030000000005` and the exact pointer arguments:
[`buffer-ex-invalid-extra-a.log`](../../tools/revkit/work/stage16/buffer-ex-invalid-extra-a.log)
and [`buffer-ex-invalid-extra-b.log`](../../tools/revkit/work/stage16/buffer-ex-invalid-extra-b.log).
Static disassembly identifies the initial stores at `0x10021c45` and
`0x10021c52`, immediately after each null test; synthesis starts later. Thus
invalid non-null addresses fault before synthesis and no API error is returned.
After context creation, output values are copied from context offsets
`+0x47774` and `+0x122448` at `0x10021e27–0x10021e44`.
These are borrowed pointers into the active synthesis context, not caller-owned
allocations. Terminal cleanup calls `VT_FreeSyncInfo_New_ENG` on the object at
context `+0x47774`, then destroys the context at `0x10027750`. The poll trace
keeps the pointer values but shows the nested SyncInfo row pointer changing
to invalid/reused data on the `DONE` return. Callers must copy data they need
while a data-bearing call is active; the outputs are not valid past terminal
cleanup. The second descriptor is also context-owned and is subject to that
same lifetime boundary.

The second output is a pointer to a 16-byte descriptor: `+0` is record count,
`+4` is the record-array pointer, and `+8`/`+0x0c` initialize to `-1`. The
timeline mapper uses the last two words as inclusive 0-based first/last row
indices selected for its active SyncInfo/text range; it restores both to
`-1` when no row matches. Runtime examples are count 0 with `(-1,-1)` for
plain text, `(0,0)` for the one-record position fixture, `(0,1)` for the
three-record marker-boundary fixture, and `(0,5)` for the six-record name
length sweep. The two-marker `start`/`end` fixture returns `(0,0)`, showing
these indices describe the subset selected by the mapper, not a blanket
`0..count-1` range. Each allocated row has a `0x210` stride and consists of
three 32-bit words at `+0/+4/+8`, a 512-byte inline field at `+0x0c..+0x20b`,
a kind byte at `+0x20c`, then three trailing bytes through the row end. The
allocator initializes row `+0` and `+4` to zero and row `+8` to `-1`. Runtime
snapshots for all three selectors show count 0 and a null array for the
repeated plain-ASCII input. Static writers establish additional row roles for
particular rewrite paths:

| Row offset | Directly established use | Limit |
| ---: | --- | --- |
| `+0`, `+4` | Initialized to zero per allocated row. `FUN_100217e0` scans SyncInfo rows selected by the marker source position, sums row `+8` counts from header start index `+0x14` through the matched row, and subtracts header `+0x18` from the first row; this is descriptor `+4`. It then adds header base `+0x10` for descriptor `+0` (`0x10021a9c–0x10021ac1`). | Both fields are mono audio-frame coordinates. `+4` is relative to the current synthesis chunk; `+0` is the utterance timeline coordinate. At poll 2 in the long fixture, marker position is `0x33`, header start index is 10, partial first-row offset is `0x1065`, and selected SyncInfo rows 10–11 have `+8 = 0x16b1` and `0x39d0`. Thus `+4 = 0x16b1 + 0x39d0 - 0x1065 = 0x401c`; adding header base `+0x10 = 0x15f90` gives `+0 = 0x19fac`. This directly resolves the nonzero `+0x18` contribution and the local/global distinction. The base advances 30,000 frames per full poll under every selector, independent of encoded byte length. Short captures have zero bases, so `+0 == +4`. At source positions 1/2/5, `+0/+4` are `0xefa`/`0x23a4`/`0x5d54`, matching SyncInfo row 0 `+8`; output byte checks across PCM16, A-law, and μ-law confirm the common frame unit. The external consumer's use of the pair is not present in the supplied caller binaries. |
| `+8` | Position `output_index + 1` while emitting an `A2 FE` two-byte marker into transformed text | Strongly indicates a marker-relative position; the downstream consumer’s coordinate convention is not recovered. |
| `+0x0c` | The first bytes hold inline mark-name text. `FUN_1002efd0` searches keys `name` then `NAME`; nonempty values shorter than 512 bytes are copied and NUL-terminated, and values at least 512 bytes are truncated to 511 bytes. An omitted `name` selects kind 2 and writes `DAT_1009f948` (zero in this DLL) at the first byte; `name=""` emits no row in the edge sweep. `FUN_1001ccc0` allocates `count * 0x210` bytes and initializes only each row's `+0`, `+4`, and `+8`; it does not clear the 512-byte field. The writer sets the name bytes and terminator (or just the first NUL for kind 2), leaving the remaining bytes untouched by these routines. | Captured post-NUL patterns, including their flag-paired differences, are returned heap contents but are not established payload fields. They may be allocator residue. No schema or reader for the remaining bytes is known. The tested 511-byte name does not establish any further semantics. |
| `+0x20c` | Byte values 1, 2, or 3 are assigned by distinct rewrite branches in `FUN_1002efd0`: kind 1 is a nonempty named mark, kind 2 is a mark with no `name` attribute, and kind 3 is the overlength named-mark branch. Runtime confirms all three under the tested selector set. | The kind values distinguish observed DLL record forms. The external caller's handling of these forms is not present in the supplied binaries. |

The row `+8` position is subsequently used as an index into the context's
per-output-character position arrays, and the mapped value is written back to
`+8` during the synthesis pipeline (`FUN_1001d5d0` and its caller
`FUN_1001c990`). One mode also adds a context-owned base before storing it.
This confirms that `+8` is a transformed-stream position used for an internal
position mapping; whether the final value denotes a unit, phone, or another
synthesis coordinate is not established. The descriptor allocator initializes
the 16-byte header; `FUN_1001ccc0` allocates the row array and initializes only
the three dwords at row offsets `+0/+4/+8`. The synthesis path reads and remaps
`+8`, but no in-engine reader of the inline name field or kind byte was found.
The API returns the descriptor address to its caller, so bytes left untouched
inside the allocated rows are observable but not evidence of a defined API
field.
Thus it is an internal rewrite/marker record list, not the
`VTDTTS_MakeInfo` selected-unit manifest. The supplied `voicetext_paul.exe`
import table lists six `vt_pau.dll` exports and omits
`VT_TextToBufferEX_ENG`; its ASCII strings also contain no EX export name.
This gives no local caller-side consumer to inspect (and does not rule out
computed dynamic lookup in that executable or callers not supplied here).
Accordingly, the public consumer contract remains unknown, but the DLL-side
meanings of kinds 1/2/3 and the inline name production rules are established
for the tested mark forms. No meaning is assigned to bytes after the first
NUL; the previous flag-dependent “payload” description was too strong because
the row allocator and mark writer leave those bytes uncleared.
The `<vtml_sub>` EX probe still returned count 0, so it did not dynamically
exercise these records. A tag-family sweep placed documented `<vtml_break>`,
`<vtml_pause>`, `<vtml_partofsp>`, `<vtml_phoneme>`, and `<vtml_sayas>` forms
before a named-mark control. Only the control emitted a row; the sweep returned
one kind-1 row for every selector. Repeating with the DLL's legacy internal
spellings (`<vt_break>`, `<vt_pause>`, `<partofsp>`, `<phoneme>`, and
`<say-as>`) emitted no rows, so those spellings are not treated as mark
records by this EX call path. The decisive contrast is
`<vtml_mark/>` followed by `<vtml_mark name="named-control"/>`: every selector
returns two rows, first kind 2 with an empty inline string, then kind 1 with
`named-control`. Both rows have equal `+0/+4` values; they are respectively
`0x90c` and `0x15e7`, matching the first SyncInfo row `+8` and the sum of the
first two row `+8` values. Their `+8` marker positions are 1 and 14. Thus kind
2 is directly identified as the unnamed VTML mark form. The descriptor header
reports count 2 and selected indices `(0,1)` in this fixture. An edge sweep
further found that an explicitly empty `name=""` produced no row, a one-space
name produced a kind-1 row carrying a space, and uppercase `NAME="Upper"`
produced a kind-1 row carrying `Upper`. Captures are
`buffer-ex-tag-family-sweep-mark-forms-{0,1,2}.log` and
`buffer-ex-tag-family-sweep-mark-name-edges-0.log`, reproduced by
`run-buffer-ex-tag-family-sweep.sh`. The positive `<vtml_mark name=.../>` probe returns
two kind-1 rows for names `start` and `end` under selectors 0, 1, and 2;
their final `+8` values were 0 and `0x25` for all three selectors. A boundary
sweep with name lengths 1, 255, 256, 510, 511, and 512 returns six rows. The
first five are kind 1; length 512 is kind 3 and confirms the static
`>= 0x200` branch at runtime. Names below 512 are stored inline, while the
512-byte value is truncated to 511 bytes by the handler. In that sweep, the
first row's `+0/+4` were zero and each later row had matching nonzero values;
the values did not equal the descriptor `+8` positions. A follow-up holds the
511-byte payload constant and moves its mark after 0, 1, 2, and 5 ASCII `x`
bytes. Descriptor fields follow the mark (`0`, `0xefa`, `0x23a4`, `0x5d54`)
while descriptor `+8` records marker-relative positions 0, 1, 2, and 5. The
enhanced GDB capture also snapshots the API's returned SyncInfo object. For
the 1-, 2-, and 5-byte prefixes, descriptor `+0/+4` exactly match SyncInfo
row 0 `+8` (`0xefa`, `0x23a4`, `0x5d54`), and that row's source-text end is
0, 1, and 4, respectively. The zero-byte case gives coordinate zero. This
rules out marker-name length as the cause and directly ties the returned
coordinate to a SyncInfo timeline entry selected at the mark's source
position. The second SyncInfo row completes the independent unit check: at
positions 1, 2, and 5, the two row `+8` values sum to 9,191, 14,281, and
29,132. Those totals equal selector-0 output lengths divided by two
(18,382/2, 28,562/2, and 58,264/2). At position 2, selector-1 and selector-2
each return 14,281 bytes, the same frame count in one-byte A-law and μ-law
output. This establishes SyncInfo row `+8` as an audio frame/sample count and
descriptor coordinates as cumulative sample boundaries independent of output
encoding. Static `FUN_100217e0` gives the exact rebasing: it adds selected row
`+8` values into descriptor `+4`, subtracting SyncInfo header `+0x18` from
the first selected row, then sets `+0` to `+4` plus header `+0x10`. Both
header bases are zero in the short captures, explaining `+0 == +4`. A later
chunk trace resolves both rebasing fields dynamically. At poll 2, the marker
position `0x33` maps through SyncInfo rows 10–11 (`+8 = 0x16b1, 0x39d0`);
header `+0x18 = 0x1065` is subtracted to produce `+4 = 0x401c`. Header
`+0x10 = 0x15f90` is then added to produce `+0 = 0x19fac`. This identifies
`+4` as the chunk-local frame coordinate and `+0` as the utterance-global
frame coordinate in the later poll. The supplied caller binaries do not show
how the external consumer uses the pair. A
ten-byte prefix fills the 60,000-byte requested output region; its zero fields
are treated as inconclusive. The option
comparison used the same
short/511/512-byte marker fixture and produced the values above for selectors
0–2; captures are in `buffer-ex-mark-0.log`,
`buffer-ex-mark-mark-1.log`, `buffer-ex-mark-mark-2.log`,
`buffer-ex-mark-boundary-0.log`, and `buffer-ex-mark-boundary-sweep-0.log`.
The option snapshots are `buffer-ex-mark-options-0.log`,
`buffer-ex-mark-options-1.log`, and `buffer-ex-mark-options-2.log`.

The controlled-position captures are `buffer-ex-mark-position-0-0.log`,
`buffer-ex-mark-position-1-0.log`, `buffer-ex-mark-position-2-0.log`,
`buffer-ex-mark-position-5-0.log`, and `buffer-ex-mark-position-10-0.log`.
The 1-, 2-, and 5-byte captures include both descriptor and SyncInfo rows.
The position-2 captures for selectors 0, 1, and 2 confirm identical
coordinates and a stable frame-count timeline across PCM16, A-law, and μ-law.
The long-mark poll captures are `buffer-ex-chunk-sync-0.log`,
`buffer-ex-chunk-sync-1.log`, and `buffer-ex-chunk-sync-2.log`; the full
selector-0 row dump is `buffer-ex-chunk-sync-full-0.log`. They are reproduced
by `run-buffer-ex-chunk-sync.sh` and show the frame base advancing across
polls before the marker row becomes populated.

The shared handler returns `-7` at `0x10021d44` when the active-object table
entry indexed by normalized speaker and thread is nonzero
(`[0x100f86e0 + speaker*0x1000 + thread_id*4]`). This is the occupied-slot
case and matches the captured duplicate-start result. The `-8` return at
`0x10022045` occurs only if the handler enters its copy loop with the
per-call copied-position cursor (`[ebp-4]`) already at or past the fixed
output limit (`[ebp-8]`). That limit is 60,000 bytes for selector 0 and 30,000
for selectors 1/2; a resume call can restore the cursor from a per-thread
carry-over slot. The normal path for a unit that crosses the limit is separate
(`0x100220ca–0x100221d1`): it writes the current chunk and saves the remainder
for a later poll. The loop's back-edge is taken only while the cursor is
below the limit, so `-8` is not the normal full-buffer result. A 78-byte VTML
control completed all nine data/terminal polls under each selector without
`-8`. Larger single-token inputs made of 1,024, 4,096, and 16,384 `a` bytes
also drained normally under selector 1 in 4, 13, and 47 calls, producing
84,789, 343,768, and 1,377,407 bytes; none returned `-8`. A separate
debugger-only probe replaced the conditional jump at
`0x10021f1a` with an unconditional jump to the `-8` epilogue, restored all six
original code bytes after the call, and repeated this for selectors 0–2. All
three calls returned `-8` without changing the incoming length value (60,000);
flag-2 cancellation then returned `1`, and a subsequent flag-1 poll returned
`-2`, confirming the active slot had been cleared. This proves the error
return and recovery path, not that ordinary input or corrupted carry-over
naturally satisfies the cursor comparison. No valid caller input has
triggered it, and a static upper bound on an individual carry-over unit has
not been proved. Reproduce with
`tools/revkit/work/stage16/run-buffer-ex-force-minus8-v3.sh SELECTOR`; the
trace forces the branch and restores the code bytes in the isolated process.
The `-9` value (`0xfffffff7`) has not been found among the returns on this
handler path or observed through `VT_TextToBufferEX_ENG`. The same signed
constant does occur in a separate export: `VT_CheckLicense_ENG` assigns `-9`
when its `dbaccess` field check fails (`FUN_10029cbe`–`FUN_10029cd6`). That
license-parser result is unrelated to the EX handler and does not establish an
EX `-9` condition. The static path shows flags other than 0 and 1 take the
active-state cancellation branch; flag 2 is the runtime-verified cancellation
value.

`FUN_100286c0` stores the text type as a byte, maps negative values to 0,
and resets the four preceding synthesis options to defaults when the stored
type is 4 or 6. Thus 260 also stores as 4 and triggers that reset. For
pitch/speed/volume it records whether each input is nonnegative; negative
values select engine defaults and speed zero first becomes 50. The later
normalizer at `0x10022850` selects explicit or engine-default values and
clamps the effective pitch to 50–200, speed to 50–400, and volume to 0–500.
Argument 13 (the scalar pause option) is stored at context `+0x1312c8`; the
downstream pause helper uses it only when its existing per-unit pause value is unset, consults
context flag `+0x122850`, and caps the resulting unsigned value at 65,535.
The exact meanings of that flag and the pause timing units remain unresolved.
The dictionary selector uses a per-speaker table for indexes 0–1023 and falls
back to the default dictionary. The EX runtime matrix used selector 0 and
plain `Hello world.`: pitch values 49/50 produced identical bytes, as did
200/201; speed 49/50 and 400/401 also paired identically. This confirms the
static clamp boundaries affect generated output as expected for this fixture.
Volume 0 and 500 produced different bytes from the default and from each
other. Pause 0, 250, and 65,535 produced the same bytes as default on both
`Hello world.` and `Hello, world.`. Text types 0–7 produced identical bytes
on both fixtures; values 4, 6, and 260 also matched when explicit
pitch/speed/volume/pause values were supplied, consistent with (but not by
themselves proving) the static reset behavior. Dictionary index 1023 matched
default on this fixture, which does not show whether a populated alternate
dictionary changes output. Record-bearing VTML captures add interaction
coverage for all selectors: explicit pitch 50 and speed 50 change the marker
coordinate and output length; volume 0, pause 250, dictionary index 1023,
and a combined extreme pitch/speed/volume/pause tuple with text type 4 leave
the observed rows unchanged. The combined type-4 case is consistent with the
static reset of those four options. This fixture includes a 512-byte marker
name, which takes the truncation path. These are controlled observations,
not a full cross-product: a populated alternate dictionary and combinations
that move the record writer through distinct branches remain open. Captures are under
`tools/revkit/work/stage16/options/` and the repeatable matrix runner is
`tools/revkit/work/stage16/run-buffer-ex-options.sh`. The punctuation run is
captured under `tools/revkit/work/stage16/options-punctuation/` and reproduced
by `tools/revkit/work/stage16/run-buffer-ex-options-punctuation.sh`.

The inline `<vtml_pause time="N"/>` command was tested by fully draining the
same `Hello<vtml_pause time="N"/>world.` text for `N=0`, `200`, and `1000`.
Selectors 1 and 2 returned 17,762, 20,962, and 33,762 bytes; selector 0
returned exactly twice those lengths as PCM16. Relative to the zero-time tag,
the 200 ms value adds 3,200 frames and the 1,000 ms value adds 16,000 frames:
both are exactly 16,000 mono frames per second times the declared milliseconds.
Each run completed on the next flag-1 poll with result 1. The untagged
`Hello world.` control is 11,803 frames, so the zero-time tag itself also
changes segmentation; comparisons of pause duration therefore use the same
tagged text. This establishes the VTML `time` unit and duration effect across
all output selectors. Scalar argument 13 is a separate setting. With
`Hello. World.`, values 0, 120, and 250 produced 17,762, 19,682, and 21,762
frames under selectors 1/2; selector-0 byte lengths are exactly twice those
values. The increments are 1,920 and 4,000 frames, or exactly 16,000 frames
per second times 120 and 250 ms. The maximum 65,535 value completed a full
selector-1 drain at 1,066,322 bytes, exactly 17,762 plus `65,535 * 16` frames.
The unchanged comma fixture (`Hello, world.`) and the unchanged VTML-pause
fixture at all four scalar values show that this setting does not add duration
in those tested contexts. These results establish argument 13's millisecond
unit and its effect on the tested sentence-boundary case; exact placement and
behavior on other punctuation/text types remain open. Reproduce with
`run-buffer-ex-scalar-pause-context.sh SELECTOR sentence` or `comma`; the
extended maximum drain is `run-buffer-ex-scalar-pause-max.sh 1`. The traces and
logs are `buffer-ex-scalar-pause-{0,1,2}.log`,
`buffer-ex-scalar-pause-{sentence,comma}-{0,1,2}.log`, and
`buffer-ex-scalar-pause-max-1.log`. Captures for the VTML attribute remain
`buffer-ex-pause-{control,pause0,pause,pause1000}-{0,1,2}.log`; reproduce them
with `run-buffer-ex-pause-compare.sh SELECTOR CASE`.

The three run logs and captured streams are in [Stage 16](../../tools/revkit/work/stage16/README.md):
[`buffer-ex-0.log`](../../tools/revkit/work/stage16/buffer-ex-0.log),
[`buffer-ex-1.log`](../../tools/revkit/work/stage16/buffer-ex-1.log), and
[`buffer-ex-2.log`](../../tools/revkit/work/stage16/buffer-ex-2.log). The
probe is reproduced by
[`run-buffer-ex.sh`](../../tools/revkit/work/stage16/run-buffer-ex.sh). The
expanded traces are [`buffer-ex-deep-0.log`](../../tools/revkit/work/stage16/buffer-ex-deep-0.log),
[`buffer-ex-deep-1.log`](../../tools/revkit/work/stage16/buffer-ex-deep-1.log),
and [`buffer-ex-deep-2.log`](../../tools/revkit/work/stage16/buffer-ex-deep-2.log),
reproduced with [`run-buffer-ex-deep.sh`](../../tools/revkit/work/stage16/run-buffer-ex-deep.sh).
Negative-flag size queries, `output_len = 1`, boundary bytes, and the live
second-object descriptor snapshot for each selector are captured in
[`buffer-ex-edges-0.log`](../../tools/revkit/work/stage16/buffer-ex-edges-0.log),
[`buffer-ex-edges-1.log`](../../tools/revkit/work/stage16/buffer-ex-edges-1.log),
and [`buffer-ex-edges-2.log`](../../tools/revkit/work/stage16/buffer-ex-edges-2.log),
reproduced by [`run-buffer-ex-edges.sh`](../../tools/revkit/work/stage16/run-buffer-ex-edges.sh).
The first selector-2 8,984-byte observation and its later 30,000-byte repeats
are preserved as excerpts in
[`buffer-ex-edges-2-first-capture.txt`](../../tools/revkit/work/stage16/buffer-ex-edges-2-first-capture.txt)
and [`buffer-ex-edges-2-30000.log`](../../tools/revkit/work/stage16/buffer-ex-edges-2-30000.log).
The `<vtml_sub>` descriptor snapshot and the unsuccessful writer/pass trace
are in [`buffer-ex-records-0.log`](../../tools/revkit/work/stage16/buffer-ex-records-0.log),
reproduced with [`run-buffer-ex-records.sh`](../../tools/revkit/work/stage16/run-buffer-ex-records.sh).

`VT_GetLicenseInfo_ENG` was called for every selector, 0–15, using both the
mounted `verification.txt` path and its bytes as an in-memory record. Every
call returned success through both paths, and each output matched byte for
byte. The log records only output lengths and equality flags. Selector 0 writes
a four-byte integer matching the positional channel value. Selector 1 returns
the positional `0` string, matching the XML `expdate` value in this record.
Selector 2 returns the host-ID field, matching XML `hostid`. Selector 3 returns
exactly the attribute text between `<vw_verify ` and `/>`, excluding both delimiters and
the adjacent `[V00]`/` Per agreement` text. Selectors 4–11 and 13–15 match
`os`, `lang`, `speaker`, `version`, `dbaccess`, `sampling`, `app`, `wavsave`,
`bgaudio`, `dbsize`, and `realtime`; selector 12 parses `savetime` into a
four-byte integer. All selectors matched the same source values in file-backed
and memory-backed calls.

Capacity probes confirm selector 0 requires four output bytes: capacity 3
returns `-5`; a four-byte buffer succeeds. Selector 1 needs room for its string
plus NUL: capacities 0 and 1 return `-5`, while capacity 2 succeeds. A null
output pointer returns `-4`; selector 99 returns `-1`. The null-path query
lengths remain 4, 1, 12, and 4 for selectors 0, 1, 2, and 12.

An expanded capacity sweep tested selectors 0–15 through both source modes
with a guarded destination and capacities from zero through 512. Selectors 0
and 12 first succeed at 4 bytes, selector 1 at 2, and selector 2 at 13. Every
selector from 3 through 15 first succeeds at 322 bytes, although the returned
value ranges from the one-byte empty `lang` string to the 290-byte attribute
text. For selectors 3–15, capacity 321 returns `-5`; selectors 0, 1, 2, and
12 fail one byte below their own minimum. A negative-capacity size query
returns 321 for selectors 3–15, but returns the selected string length for
selectors 1 and 2 (1 and 12), and 4 for numeric selectors 0 and 12. A null
destination with adequate capacity returns `-4` for every selector. No call
changed bytes beyond the advertised capacity. The identical 321-byte size
query and 322-byte minimum across selectors 3–15 match the full comment
payload length and remain fixed even for short attribute values. This runtime
behavior strongly indicates that these selectors size-check the shared
comment/tag source before returning the selected attribute; the precise local
variable assignment comes from decompiler/disassembly analysis.

`VT_GetLicenseComment_ENG` was separately swept over capacities 0–2048 through
both the default and explicit-file paths. Both first succeed at 322 bytes;
capacity 321 returns `-4`, and guard bytes remain intact. Successful calls
return 322 (including NUL), with 321 payload bytes.

### License record parsing and runtime trace

The decompiled `VT_CheckLicense_ENG` path passes its source/key/length
arguments to `FUN_10028db0`, then to `FUN_10028f20`. The reader substitutes
built-in defaults for absent source/key inputs. A non-null in-memory input
with positive length selects `FUN_1002e3a0`; otherwise `FUN_1002e210` opens
the selected source through the engine's file abstraction. Both helpers copy
ordinary input unchanged. For input at least seven bytes long, they compare
the first three bytes with the final three; when equal, they allocate an
output seven bytes shorter, copy the bytes between the three-byte prefix and
the one-byte seed plus three-byte suffix, then subtract that seed byte from
each copied byte. The path helper performs the same first-three/last-three
comparison and reads the seed from immediately before the suffix. This is
a simple framing/byte-offset transform observed in the instructions; this
framing step itself is not the MD5 conversion described below and does not
establish a signature check.
The resulting length is used to NUL-terminate the buffer before parsing.

`FUN_10029680` copies the resulting text and splits it at semicolons. For each
entry, its key scanner uses the delimiter set space, tab, colon, and equals;
the value scanner stops at a tab. The scanner walks over quoted spans while
searching for delimiters, and entries whose first key character is `#` or `/`
are skipped. Key/value pairs are stored in an internal linked list.
`FUN_10028f20` looks up the exact key `License`, then splits that value at its
first colon in place. It retains the part before the colon and looks up the
entire remaining substring as a second exact key. A successful second lookup
sets a separate presence flag. This is the parser's observed linkage rule; it
does not by itself identify the human meaning of either value or prove an
authentication scheme.

After conversion succeeds, `FUN_10028990` splits the linked value at colons
into seven stored fields at context offsets `+0x20` through `+0x38`. Runtime
file and memory calls establish these sample mappings:

| Stored field | Source-record position | Runtime observation |
| ---: | --- | --- |
| 1 (`+0x20`) | `hostid` token | Selector 2 returned 13 bytes including NUL (12 characters) and matched XML `hostid`. |
| 2 (`+0x24`) | `VW_VTAPI` token | Present in the record; no public selector returns this field. |
| 3 (`+0x28`) | `0` token | Selector 1 returned the string and it matched XML `expdate`. |
| 4 (`+0x2c`) | `6` token | Selector 0 returned a four-byte integer matching XML `channel`. |
| 5 (`+0x30`) | `Kurzweil` token | Matched XML `user`. |
| 6 (`+0x34`) | `win32` token | Matched XML `os`. |
| 7 (`+0x38`) | `[V00]<vw_verify .../> Per agreement` | `VT_GetLicenseComment_ENG` returned this complete colon field byte-for-byte; selector 3 returned only its tag attribute text. |

The sample comparisons establish the positional associations above for this
record, not a universal schema guarantee. The first 96-character token before
these seven fields is not exposed by these selectors. Static disassembly
requires exactly 96 characters, splits it into two 48-character halves, and
passes each through a fixed conversion routine. That routine selects repeated
two-character slices at six-character intervals, parses hexadecimal values,
performs further arithmetic, and formats comparison values. The checker also
passes the preceding `License` value and the value reached by the linked-key
lookup into this routine. Its consistency path contains the verified standard
MD5 primitive described below; this does not establish an authentication
guarantee or the business meaning of its outputs. The `VW_VTAPI` token remains
an opaque product/record discriminator.
The local `data-*` inventory contains only this one file with the
`vw_verify`/`VW_VTAPI` record markers, so there is no second accepted local
record with which to validate whether these positions and associations
generalize.

The `hostid` label above names the field in the archived record. It does not
show that the value was generated for, or issued to, the workstation used for
this probe; the record came with the archived VoiceText package.

A value-suppressed same-length mutation matrix called `VT_CheckLicense_ENG`
against the in-memory sample with one positional/XML edit per call. The
unmodified control and its repeat returned 0. Changing only stored field 3
(`+0x28`, the sampled `0` expiry token), only XML `expdate`, or both together
to `1` returned -2. Changing only stored field 2 (`+0x24`, `VW_VTAPI`) to an
eight-byte alternate, only stored field 4 (`+0x2c`, channel) from `6` to `7`,
only XML `channel` from `6` to `7`, or both channel values together also
returned -2.
The capture prints case labels and statuses, not source values:
[`license-positional-field-matrix-api.log`](../../tools/revkit/work/stage21/license-positional-field-matrix-api.log),
reproduced with
[`run-license-positional-field-matrix.sh`](../../tools/revkit/work/stage21/run-license-positional-field-matrix.sh).
This establishes rejection under edits to each of these fields in this
record's existing derived-token state. Since even paired field edits still
reject, this does not demonstrate that the XML `expdate` mirrors the
positional token, that `channel` is cross-validated, or what `VW_VTAPI` means;
the unchanged 96-character derived field may itself depend on the edited
record content. The field meanings and validation relation remain open.

`FUN_10029170` passes the retained `License` value, the second lookup's value,
and the presence flag (as either a null pointer or a built-in string) to
`FUN_10014dd0`. The 96-character value must be exactly 96 bytes in the
observed C-string path, then is split into two 48-character halves. Each half
is handled by the same conversion routine. For each half, the disassembly
selects two-character slices at offsets `0, 6, 12, 18, 24, 30, 36, 42`;
`%.2s` formatting joins the first four slices and the last four slices into
two eight-character strings, and base-16 integer conversion produces two
32-bit values. The four intervening characters in each six-character group
do not enter those two parsed values through this formatting path. The helper
then combines the parsed values with fixed arithmetic and formatting. A
second helper concatenates four strings without separators, in argument order
3, 4, 2, 1, and hashes the byte sequence with standard MD5. Its initialization
state is the standard MD5 IV, its block routine uses the MD5 round constants
and 64 steps, and it emits the 16-byte digest as 32 lowercase hexadecimal
characters. The conversion code incorporates that digest into a generated
48-character comparison string and compares it against the supplied half.
The checker processes the first half, then passes that first-half string as an
input to the second-half conversion; consequently, first-half byte changes
also change the second digest input. The caller bounds one output at
99,999,999 and another at 1,024 (zero selects the default 1,024). The date
validator at `FUN_10029510` repeats that conversion path, formats the local
system date as `YYYYMMDD`, and rejects when today is later than its converted
value. This establishes a linked two-value validation pipeline and its
numeric bounds. The positional `0` matching XML `expdate` is a sample
correlation; it does not prove that the date conversion consumes that field.
The MD5 primitive itself was called at module offsets `0x142f0` (initialize),
`0x14320` (update), and `0x14410` (finalize) for ten standard vectors covering
empty input, `a`, `abc`, and repeated-`a` lengths 55, 56, 63, 64, 65, 127, and
128 bytes; every digest matched its reference. This confirms the hash
primitive and its padding/block behavior. It does not establish that the
license token is a signature, that the use of MD5 provides an authentication
guarantee, or what the derived numeric values mean. See the cited PE
disassembly at `0x10014dd0`–`0x100152b3` and MD5 helpers
`0x100142f0`–`0x10014510` for the conversion and hash paths. The sanitized
vector capture is
[`license-md5-vectors-api.log`](../../tools/revkit/work/stage21/license-md5-vectors-api.log);
the PE32 probe and runner are in Stage 21.

A value-suppressed byte-mutation matrix strengthens the acceptance boundary.
Starting from the accepted archived record, each of the 96 token bytes was
changed individually in two ways: to the next hexadecimal digit, and to `g`.
All 192 variants returned `-2`; the unchanged record returned `0`. This
includes all 32 byte positions selected into the four parsed hexadecimal
values and all 64 intervening positions. Thus every position rejected these
two value-changing substitutions in this record. The intervening bytes do not
feed those four parsed values through the observed `%.2s` extraction, so their
rejection is consistent with the later comparison of each generated value
against its supplied 48-character half. The probe does not identify the
individual rejecting instruction.

A follow-up changed the case of each alphabetic hex byte individually. Of 22
letters, 10 case flips preserved the checker result `0` and 12 returned `-2`.
All 10 accepted flips were in gap positions of the second half; the 10
first-half gap flips and both selected-byte flips failed. The ten individually
accepted second-half flips were also tested as a complete Boolean product:
all 1,024 combinations returned `0`. Static flow explains the asymmetry for
this record: the first half is passed as raw text into the second-half MD5
input, whereas second-half gap bytes are outside the parsed values/hash inputs
and compare through the DLL's ASCII case-folding routine. The second-half
selected bytes did not contain alphabetic hex characters in this record, so
their case behavior is untested. The three captures contain only byte indexes,
selected/gap class, mutation class, subset masks, and return codes, never token
contents. Reproduce with
`tools/revkit/work/stage21/run-license-token-byte-matrix.sh`; see
[`license-token-byte-matrix-api.log`](../../tools/revkit/work/stage21/license-token-byte-matrix-api.log),
[`license-token-byte-matrix-v2-api.log`](../../tools/revkit/work/stage21/license-token-byte-matrix-v2-api.log),
[`license-token-byte-matrix-v3-api.log`](../../tools/revkit/work/stage21/license-token-byte-matrix-v3-api.log),
plus the runner/source in Stage 21.

The mounted 468-byte `data-common/verify/verification.txt` record was accepted
directly through `VT_CheckLicense_ENG` in both modes: as the explicit file
path and as the in-memory buffer (`return=0` in each case). Passing the literal
filename `License` returns `-1`, confirming that the first argument is a source
path rather than the record key. The null-path default call also returned `0`,
but is separate from the supplied-file result. `VT_GetLicenseComment_ENG` on
the supplied path returned 322 bytes including the terminating NUL; the 321
payload bytes matched the record's final nonempty colon field exactly.

The file and memory readers both accepted an equivalent framed record built
with a three-byte `XYZ` prefix/suffix and a one-byte additive seed. Changing
the seed or the final suffix byte made the check return `-2`. Plain, valid
framed, and supplied-file forms all returned `0`. Mutating each of the first
seven visible fields, selected XML values, removing queried attributes present
in the record, changing tag boundaries, or passing a short buffer returned
`-2`.
Those mutations therefore fail before the individual attribute return codes
can be isolated; they do not show which changed field caused rejection.

After successful record setup, `VT_CheckLicense_ENG` searches the resulting
text for the literal `<vw_verify ` and then the next `/>`; it copies the
attribute text between those delimiters and queries the names in order: `os`,
`lang`, `speaker`, `version`, `dbaccess`, `sampling`, and `dbsize`. Failed
checks take return branches `-5` through `-11` in that order. The attribute
helper applies local value checks (including comparisons against built-in
allowed-value strings for some fields). This is custom string/token processing
in the recovered code path; it does not establish general XML parsing or
validation. See the [Ghidra pseudocode](../../tools/revkit/work/reports/vt_pau-exported-api.c)
and [PE disassembly](../../tools/revkit/work/reports/vt_pau-objdump-disassembly.txt)
for the code paths at `0x10029b80`, `0x10029d30`, `0x10028db0`,
`0x10028f20`, `0x10029680`, `0x100297a0`, `0x10029920`, `0x10029970`,
`0x10028990`, `0x10029170`, `0x10029510`, and `0x1002a110`.

The `FUN_10029d30` attribute helper was also called directly at its module
address with the record's extracted attribute text, to avoid the earlier
96-character record-conversion gate. Baseline values passed for `os`, `version`,
`dbaccess`, and `sampling`; `linux`, `0`, `none`, and `8` returned false for
those respective keys. Missing-key lookups returned `-1`, which the enclosing
checker treats as true because it rejects only a zero result; this includes
absent `lang` and the tested removals.

For the per-voice checks, the probe passed each of the DLL's six actual voice
table pointers. The supplied `speaker` list passed slots 0/1/3/4 (Kate, Paul,
Julie, James) and failed slots 2/5 (em001, Ashley), matching membership in the
tag's list. This establishes the `speaker` value as an allow-list for the
current DLL's voice IDs. With no loaded model, `dbsize` values 0 and
999,999,999 passed. After `VT_LOADTTS_ENG(NULL, -1, NULL, NULL)` loaded the
default Paul slot, the same `dbsize=300` check failed, and the full file-backed
checker returned `-11`. The export-reported Paul database size is
508,121,688 bytes (496,212 KiB). Direct values 300, 450, and 451 failed;
452 and 500 passed. The disassembly parses `dbsize`, multiplies by 1,000 and
then 11/10, and compares that threshold in KiB with the loaded database size:
451 yields 496,100 KiB, below the observed size, while 452 yields 497,200 KiB.
Thus `dbsize` is a per-voice loaded-database size allowance in 1,100-KiB
steps. The unload call completed after the probes.

The same 468-byte record was then passed through `VT_LOADTTS_EXT_ENG` for
James slot 4, using both argument-6 file path and argument-7 memory forms.
Both calls returned AX `0`; the follow-up checker returned `0`, the license
gate became `1`, and dictionary capacity became `6`. The loaded James database
is 253,797,273 bytes (247,848 KiB by integer division), below the record's
330,000-KiB allowance (`dbsize=300`). This contrasts with Paul's 496,212-KiB
database, which exceeds that allowance and fails with `-11`. The per-voice
database-size check therefore explains the different loader license states
for this record.

On a successful check, the loader calls `FUN_10028c50` for the capacity value;
its disassembly returns the parsed record's dword at structure offset `+0x10`
(or 1 when the parser status is nonzero). The loader stores it at per-speaker
state offset `+0x4d14`. Six text-synthesis entrypoints compare a selected
dictionary index with this limit; `VT_UNLOAD_UserDict_EXT_ENG` also scans only
that many per-speaker dictionary references when checking whether a global
dictionary can be unloaded. Capacity 6 therefore permits per-speaker indexes
0 through 5 in the James state. It matches this record's `channel=6`, which
supports but does not prove that `channel` supplies the capacity; only one
successful record is available. Captures: `load-ext-james-license-api.log`
and `userdict-license-state3-api.log`.

The loader call was traced into `VT_LOADTTS_EXT_ENG` at `0x10027af0`. The
sample's legacy wrapper passes slot `-1`, a null database path, and extended
arguments `(0, 0xffffffff, NULL, 0, 0xffffffff)`. The extended export returned
low AX `0`; immediately afterward the database-size query returned `1` with
508,121,688 bytes. At `0x10025e70`, after default-path selection and before
shared initialization, the path buffer contains `../`. A separate entry
substitution of `Z:\work\` reaches the same checkpoint as `Z:/work/` and
also returns low AX `0` with the same database size. In this container,
`Z:\work\` is the absolute equivalent of the sample's relative parent
directory. The path parameter therefore supplies the base used for the
package's relative resources; it is not the `M16` model directory itself.
Substituting `Z:\voice-data\paul\M16` or its parent
`Z:\voice-data\paul` instead returns low AX `3` before the
`0x10025e70` checkpoint, matching the header's tagger-error constant. These
failures characterize only those path candidates, not every malformed-path
case.

Static pseudocode normalizes slot values outside 0–5 to 1; null and empty
database paths use the default-path helper. A nonempty path is copied into a
514-byte local buffer using the measured input length, passed through the path
normalizer, compared with the shared current path when initialized, then
copied into that shared path before initialization continues. The decompiled
copy has no visible bound check against the 514-byte destination; overlong-path
behavior has not been probed. Extended arguments 4 and 5 are not read in the
decompiled body. Arguments 6–8 flow into `VT_CheckLicense_ENG` and the
license-derived database-size helper. Their roles are file path, in-memory
license text, and text length, consistent with the shared license parser's
call sites. The original wrapper call uses the null-file/null-text/-1-length
tuple. A separate run passed the 468-byte
`data-common/verify/verification.txt` file as argument 6 while arguments 7
and 8 remained null and `-1`. The checker returned `-11`, but the loader returned
low AX `0` and the database-size query still returned
`1 / 508,121,688`. The slot state afterward directly reads license gate 0 and
dictionary capacity 1. This confirms that this nonzero license status does
not prevent model loading, while the license-dependent features remain gated.
A PE32 harness then read the same record into a heap buffer and passed it as
argument 7 with argument 6 null and argument 8 set to the 468-byte length.
That call also returned low AX `0` and loaded the same database. Calling the
checker against the loaded speaker and memory buffer returned `-11`; direct
slot reads again showed gate 0 and capacity 1. For this record, file-path and
memory-buffer loader forms therefore converge on the same status and state.
The same argument forms were then tested against loaded James slot 4, whose
model root required the explicit `/work/data-james` alias. Both returned AX
`0`; their direct checker calls returned `0`, and the loaded slot had license
gate 1 and dictionary capacity 6. Its 247,848-KiB size is below this record's
330,000-KiB `dbsize` allowance, unlike Paul's 496,212-KiB model. This
cross-slot result confirms the size check is per loaded voice and that a
passing license enables the loader's per-speaker gate. The capacity value's
match with `channel=6` remains an inference as described above.
Additional direct PE32 calls characterize repeated calls and slot selection.
With slot 1 loaded from `Z:\work\`, a second call using a separate buffer
containing the identical path returned AX `0`; the state occupancy mask stayed
`0x02`, and the size query stayed `1 / 508,121,688`. After unloading slot 1,
the mask became zero and a fresh call with the same base returned AX `0` and
restored the same size. Thus same-slot, same-base calls are accepted both as
an already-loaded no-op and as a reload after teardown.

While slot 1 remained loaded, a call with the deliberately nonexistent base
`Z:\work\not-a-real-directory\` returned AX `1`, with the occupancy mask and
database size unchanged. Calling again with `Z:\work\` returned AX `0`. The
static branch explains this as a shared-base mismatch combined with an already
loaded selected slot: it returns before copying the new base or entering
resource initialization. AX `1` in this case does not mean the supplied path
was loaded. This result is observed for a single different-base call against
loaded slot 1. A separate two-slot call loaded Paul and James together, then
tried the same different base against each loaded slot. Both calls returned
AX `1`, occupancy stayed `0x12`, and both database-size queries stayed at
their original values. Repeating each call with the original base returned
AX `0` and preserved both voices. Thus the changed-base short-circuit was
confirmed for both selected slots in this loaded pair. See
`load-ext-multi-path-switch-api.log` and
`run-load-ext-multi-path-switch-container.sh`.

Separate fresh-process first calls with speaker slots `-2`, `-1`, and `6`
each returned AX `0` and allocated only slot 1 (`occupancy=0x02`). This
directly confirms normalization on the first call. A further isolated sequence
repeated each of `-2`, `-1`, `6`, and `1000` with a different base path after
slot 1 had loaded. Every call returned AX `0`; slot 1 remained queryable at
508,121,688 bytes. A subsequent valid slot-4 call using the original base
loaded James successfully, yielding occupancy `0x12` and both original
database sizes. No process exception occurred in those four cases. This closes
the observed outcomes for those inputs, not the invalid-index boundary
generally: pseudocode reads the per-slot flag at
`(&DAT_100a7480)[param_2]` using the original argument before the normalized
slot is used for the state lookup. For these invalid values, that is an
out-of-range read from the six-entry array, so behavior for other integers or
memory layouts remains unspecified. See
`load-ext-invalid-switch-recovery-*-api.log` and
`run-load-ext-invalid-switch-container.sh`.

A six-slot first-call matrix used `vt_pau.dll`, the Paul M16 engine build, and
the local Kate, Paul, Julie, and James model roots mounted read-only at the
paths requested by the DLL. Slot 1 returned AX `0` and allocated slot 1 with
database size 508,121,688 bytes. Slot 4 returned AX `0` and allocated slot 4
with database size 253,797,273 bytes. Slots 0, 2, 3, and 5 returned AX `8`
without allocating state; header constant 8 is `VT_LOADTTS_ERROR_PROSODY_DB`.
File traces show these failures occur on each slot's first tree open:
`data-kate/M16/ttsdata/tree3/pitch/nbt.tree3` for slot 0,
`data-em001/M16/ttsdata/tree3/pitch/nbt.tree3` for slot 2,
`data-julie/M16/ttsdata/tree3/pitch/nbt.tree3` for slot 3, and
`data-ashley/M16/ttsdata/tree3/pitch/nbt.tree3` for slot 5; each open reports
not found (`0xc000003a`). Kate and Julie provide older `tree2` layouts, so the
requested `tree3` file is absent. The local tree inventory has no model roots
for `em001` or `Ashley`. These results describe this DLL and local assets, not
all possible builds for those speaker IDs.

The file trace for slot 4 identifies why the intermediate matrix still failed:
`vt_pau.dll` requested `Z:\work\data-james\M16\ttsdata\tree3\pitch\nbt.tree3`.
The Stage 8 runtime alone had no `/work/data-james` mount; Stage 17 maps the
James package to `/work/data-jame` for the supplied James host/DLL. Adding a
separate read-only `/work/data-james` alias allowed slot 4 to load, along with
its `data-james` DAT and UPM resources. The earlier
`load-ext-slot-matrix-api.log` and `load-ext-slot-matrix-corrected-api.log`
are retained as mount-setup attempts; `load-ext-slot-matrix-alias-api.log` is
the six-slot capture with the additional alias. `load-ext-slot4-filetrace-v3-api.log`
records the James opens, and the per-slot file traces record the four missing
tree requests.

A further two-slot probe loaded Paul/slot 1 and James/slot 4 into the same DLL
process in both orders. Both calls returned AX `0`; occupancy was `0x12`, and
each size query continued to report its own value (508,121,688 and
253,797,273 bytes). Unloading either voice left the other slot loaded and
queryable; unloading the remaining voice cleared both slots. This directly
confirms simultaneous state for these two package paths and last-slot teardown
in both load orders. It does not extend to concurrent calls from multiple
threads. See `load-ext-multi-api.log`.

Read-only missing-resource overlays now map several loader failures in a fresh
Paul process. Omitting `data-common/dict-eng/atmt.tree3` returns 3;
`engbi.tree3` returns 4; `sbd.tree3`, `tppdict_eng`, or `hashidx_eng_tpp`
returns 5; `data-paul/M16/ttsdata/dist_tbl/cepdist.tbl` returns 6;
`data-paul/M16/mc_idx_tbl/unit-gen.idx` and `data-paul/M16/dblist.idx`
return 7; a missing Paul prosody tree returns 8; missing
`data-paul/M16/dat/merged-gen.dat` returns 9; and missing
`data-paul/M16/dat/merged-gen.upm` returns 10. Each failure leaves occupancy
zero. These associate the exercised missing resources with the public header
codes for break index, TPP dictionary, table, unit index, prosody, PCM, and PM
database errors. `load-ext-resource-error-*-api.log` captures the results.
The matching overlay generator creates symlink-only views; all vendor files
remain untouched. The initial attempt to omit `dblist.idx` also masked the
whole M16 tree and is not evidence for that file; the corrected
`omit-dblist-idx-preserve-tree` overlay preserved the other directories.

An empty `dblist.idx` returns 7 as well. An empty `merged-gen.dat` is accepted:
the load returns 0, and database size drops from 508,121,688 to 161,530,436
bytes. That 346,591,252-byte reduction exactly equals the original file's
length. For this resource, the loader accepts a zero-length payload and its
reported size tracks the file's physical length. This does not establish
validation behavior for arbitrary nonempty corruption.

Static pseudocode sets error 2 when either shared initialization's 0x2042c-byte
allocation (`FUN_10025e70`) or per-speaker state's 0x4d18-byte allocation
(`FUN_10025f10`) fails. A GDB probe replaced the shared allocation's valid
return pointer (`0x01430020`) with NULL; the extended-loader export then
returned low AX 2, while an unmodified control returned 0. This confirms the
shared-allocation error branch by forcing its post-allocation condition. A
second probe replaced the per-speaker allocation's valid return pointer
(`0x0147f198`) with NULL at `0x10025f2a`; this independently returned low AX 2.
Both probes force the null condition after a successful allocator return; they
do not reproduce natural allocator failure. Error 11 is declared in the header
but no assignment was found in the recovered extended-loader error-propagation
chain: the reviewed helpers map allocation failure to 2 and resource failures
to 3–10. This finding is scoped to this DLL's recovered loader path. It is
distinct from the license checker returning `-11` for Paul's insufficient
`dbsize`: the extended loader still returns AX 0, with license gate 0 and
capacity 1. Concurrent calls from multiple threads, malformed/overlong paths,
and natural allocation failures remain open. The successful-license
behavior is established only for this one record and the Paul/James slot pair. The valid
license captures retain only lengths, statuses, and internal gate/capacity
values, never record contents.
The valid capture logs
only the path, pointer, length, and status, not record contents. The first
string-injection attempt produced a corrupted path and is invalid evidence;
`load-ext-license-file-api.log` is retained for that failed setup. Traces:
`load-ext-wrapper-api.log`, `load-ext-default-path-api.log`,
`load-ext-explicit-runtime-root-api.log`, `load-ext-explicit-db-path-api.log`,
`load-ext-explicit-parent-path-api.log`,
`load-ext-license-file-bytewise-api.log`,
`load-ext-license-file-state-api.log`, and
`load-ext-memory-license-state-api.log`.

The export's attribute sequence queries `os`, `lang`, `speaker`, `version`,
`dbaccess`, `sampling`, and `dbsize`; it does not query `app` or `expdate`
there. After loading Paul, the direct helper rejected the sample's
`dbsize=300`; `VT_GetDBSize_ENG` reported 508,121,688 bytes for that voice in
[`export-queries-api.log`](../../tools/revkit/work/stage16/export-queries-api.log).
Although the sample's positional `0` matches XML `expdate`, this
equality alone does not connect that XML attribute to the separate positional
date check.

The parser is custom token processing rather than general XML validation. The
supplied record succeeds because it has valid linked conversion data and a
tag envelope; changing its bytes causes the setup path to return `-2` before
we can exercise the checker’s individual `-5` through `-11` failures through
the export. A separate framed temporary file was created and deleted under
container `/tmp`; vendor inputs remained mounted read-only. The runtime output
omits record contents and the reproducible harness is
[`license-paths-api.log`](../../tools/revkit/work/stage16/license-paths-api.log),
[`probe-license-api.c`](../../tools/revkit/work/stage16/probe-license-api.c),
and [`run-license-probes.sh`](../../tools/revkit/work/stage16/run-license-probes.sh).

The five data exports read as `VT_gHeapStartAddress_ENG = 0` and version
components `(first, second, third, fourth) = (3, 11, 7, 1)` after model load.
The tuple matches both the DLL's PE FileVersion and ProductVersion strings
`3.11.7.1`; this supports interpreting the four exports as version
components. Separate traces read the heap-start export as zero before model
load, after load, and after unloading speaker 1. PE section headers place its
RVA `0xFF11C` inside `.data` (RVA `0x77000`, VirtualSize `0x894B0`, raw size
`0x29000`), beyond the section's file-backed bytes; the image loader therefore
zero-fills this virtual-tail location. A Stage 21 hardware read/write watchpoint
was armed on the complete four-byte cell before the sample's model load. It
observed zero accesses while the sample loaded Paul, returned EAX 1 from its
successful file-synthesis call, and explicitly unloaded speaker 1; the cell
read as zero before load, after synthesis, and after unload. The capture is
[`heap-export-access-api.log`](../../tools/revkit/work/stage21/heap-export-access-api.log),
reproduced by [`run-heap-export-access.sh`](../../tools/revkit/work/stage21/run-heap-export-access.sh).
In a separate disposable process, the four version cells accepted writes of
203, 204, 205, and 206, then were restored to 3, 11, 7, and 1 before the host
continued. While the changed values were live,
`VT_GetDefVersion_ENG` still returned `Paul-M16-FileIO`, consistent with its
static implementation building a string from separate fixed text. This
establishes writable data cells and independence of that getter in the tested
call path; it does not establish normal external-client writes or other
consumers. See Stage 21 `version-export-mutation-api.log` and
`trace-version-export-mutation.gdb`.
Each of the four bundled voice-specific executables imports six named
functions from its matching VoiceText DLL and none of the five DLL data
exports. Each also imports `GetProcAddress`; its sole code reference is in
`_init_codepage_func`, which gets `msvcrt.dll` and looks up
`___lc_codepage_func`, falling back to `__lc_codepage`. It does not dynamically
resolve a VoiceText export. This closes the named-import and direct
GetProcAddress paths for the four bundled hosts, but not for external hosts or
other dynamic-resolution mechanisms.
Separately, all four matching DLLs have the same internal Windows-heap setup:
the `HeapCreate` call at `0x10066924` stores its returned handle in global
`0x101004a0`; allocator helpers use that handle for `HeapAlloc`, `HeapFree`,
`HeapReAlloc`, and `HeapDestroy`. That handle is distinct from the exported
`VT_gHeapStartAddress_ENG` cell at `0x100ff11c`. This maps the handle used by
the observed internal allocation path, not the intended use of the exported
cell. Disassembly scans of all four DLLs found no direct absolute-address
reference to it, and no reference appears in the reviewed pseudocode. Together
these observations show no access on this ordinary tested lifecycle or
evidence of consumption by the bundled hosts; other branches, external host
writes, derived addressing, and the intended contract remain unknown. The
import and callsite survey with reproduction commands is
[`heap-export-host-lookup-static.txt`](../../tools/revkit/work/stage21/heap-export-host-lookup-static.txt).
The earlier value-only capture is
[`helper-exports-api.log`](../../tools/revkit/work/stage16/helper-exports-api.log),
reproduced by `run-helper-exports.sh`.
The three export names `VT_SetDecimal0Pron_ENG`, `VT_SetPhone0Pron_ENG`, and
`VT_SetVirtualTagMode_ENG` all map to RVA `0x28420`, a single `ret` in the
x86 listing. After the sample's file-synthesis call returned 1, a Stage 21
trace called that shared target three times, using one EAT-associated export
name per call. Each preserved its distinct EAX sentinel and the pre-call ESP;
the sampled speaker pointer, play-state, and history-mode values also stayed
unchanged. The direct calls and reproduction are in
[`noop-setter-aliases-api.log`](../../tools/revkit/work/stage21/noop-setter-aliases-api.log)
and [`run-noop-setter-aliases.sh`](../../tools/revkit/work/stage21/run-noop-setter-aliases.sh).

The NULL-path `VT_GetLicenseComment_ENG` query returned
322, but its output is deliberately omitted from the capture because it
contains embedded license-record fields. The direct null-path
`VT_CheckLicense_ENG` call returned `0`; `VT_INIT_ENG` returned `-1`. Exact
trace lines are in
[`helper-exports-api.log`](../../tools/revkit/work/stage16/helper-exports-api.log);
the probe is reproducible with `run-helper-exports.sh`.

### Remaining export-only text helpers: static trace

These wrappers are not declared by `vt_eng.h`. The export pseudocode routes
them to internal routines; the following are disassembly observations, not
runtime contracts:

| Export | Internal path observations | Still open |
| --- | --- | --- |
| `VT_TextToBufferEX_ENG` | Selectors 0–2 match ordinary buffer formats on the short fixture; chunking, polling, cancel/busy/errors, size query, length-not-capacity, optional outputs, and row descriptor are documented above. `<vtml_mark/>` returns kind 2 with an empty payload, nonempty named marks return kind 1, and a 512-byte name returns kind 3 truncated to 511 bytes; all three kinds are runtime-confirmed across selectors 0–2. `FUN_100217e0` maps descriptor `+0/+4` through SyncInfo row `+8`; their unit is one audio frame, confirmed by exact sums across both timeline rows and output byte lengths for PCM16, A-law, and μ-law captures. Static code maps `+4` through header `+0x18` and `+0` through header `+0x10`; the long poll trace confirms nonzero `+0x10` and the `+0` equation directly. The supplied caller binaries do not reveal how an external consumer uses the decoded chunk-local/global frame pair. Inline `<vtml_pause time="N"/>` was fully drained at N=0/200/1000; relative to N=0, 200 and 1,000 ms add 3,200 and 16,000 frames under all selectors, confirming milliseconds at 16 kHz. Scalar argument 13 also uses milliseconds in the tested sentence-boundary case: values 120/250 add 1,920/4,000 frames across selectors, and selector 1 fully drained value 65,535 at 1,066,322 bytes. Values 0/120/250/65,535 leave the tested comma and inline-VTML-pause fixtures unchanged; exact placement and other text categories remain open. Pause 0/250/65535, dictionary 1023, and text types 0–7 are byte-identical to default on the single-sentence plain and comma-punctuation fixtures; pitch, speed, and volume change output. Static wrapper rejects selectors outside 0–2 with `-1` and dispatches to `0x100200c0`, `0x10020460`, and `0x10020930`. A debugger-forced branch probe on all selectors returned `-8`; flag-2 cancellation returned `1` and the next poll returned `-2`. | Whether the cursor comparison can be reached naturally from valid or corrupted carry-over state; external handling of descriptor rows; populated alternate-dictionary effects; scalar pause-argument placement on other text categories; and the full option cross-product. |
| `VT_TextToPreprocessInfoFile_ENG` | The wrapper at `0x1001dd60` delegates to `FUN_1001ef20` (`0x1001ef20`) and is decompiled `void`. Flag 0 returns before pointer checks. A fresh-process matrix used speaker 1, pitch/speed/volume/pause `-1`, dictionary 0, and text type 0; byte flags 1–10, 11, 254, and 255 completed with raw low AX 1 (not a declared return). Correct heap-backed calls establish that flags 1–5, 7–11, 254, and 255 use `param2` literally; flag 4 appends `.0`–`.3`, while flag 6 ignores `param2` and writes fixed `test.pcm`. The earlier malformed names resulted from GDB setup corrupting the path before API entry. Flag 3 emits one dot per source byte, a count, and decimal engine phone/control bytes that match the recovered phone codebook; flag 5 renders the same stream as phone labels. Flag 7 writes source text and per-word surfaces, structural boundary codes, phone labels, metadata-bit letters, and inclusive ASCII spans. Its leading period string is selected among eight fixed lengths by a `GetTickCount`-seeded generator; identical-input runs vary, and the reason for including the filler is unknown. Flag 10 spaces tested terminal punctuation; flag 6 wrote 23,606 bytes of PCM. The boundary-class names, bit-label meanings, broader text/punctuation behavior, malformed nonempty input, and invalid scalar combinations remain open. A corrected 44-call matrix crossed flags 3/5/7/10 with baseline, pitch 50/200, speed 50/400, volume 0/500, pause 0/250/65,535, and a repeat baseline; all calls returned raw EAX 1 and each mode’s output was byte-identical across variants. A separate 11-call flag-6 capture shows pitch 50/200 and speed 50/400 change PCM length, volume 0 produces silence, volume 500 changes samples, and pause 0/250/65,535 leaves this `Hello world.` output byte-identical to baseline. These settings results are limited to this fixture. Flags 12–253 were swept in one process on the fixed fixture; cross-text behavior remains open. See the isolated observations and output-record interpretation below, the heap-backed captures and runners in Stage 16/21, and the coverage-map entry. |
| `VT_TextToPcmBuffer_ProgressBar_ENG` | Runtime after model load with the Stage 16 text, thread 0, speaker 1, option values -1, and null window/message: 23,606 bytes were written, matching ordinary buffer format 0; guard bytes remained intact. A target-flow long call used a 55-byte repeated sample utterance, a 60,000-byte buffer, `hwnd=0`, and message `0x8005`. It returned status 0 with three full 60,000-byte output chunks, then status 1 with 3,224 bytes (183,224 bytes total). Four calls reached the `PostMessageA` import. Corrected source traces compare the context's first DWORD with pushed wParam on every post. Static instructions show `FUN_10026ab0` computes the byte length of its third string argument into context `+0`; in this path that string is returned by `FUN_1001c990`. A hardware watchpoint observed context `+0` change from zero to `0xc4` at the store in `FUN_10026ab0`. A separate capture measured the working string at 167 bytes, with 111 leading periods followed by a space and the repeated 55-byte utterance; its wParam was `0xa7` on all four posts. Five completed runs observed `0x59`, `0x72`, `0xa4`, `0xa7`, and `0xc4`, stable within each run. After subtracting the 55-byte source and one space, their prefix lengths are 33, 58, 108, 111, and 140, all members of the independently measured flag-7 filler-length set. This explains wParam as the working string's byte length; the relationship between the two paths' filler generators remains an inference. Two posts occurred during the initial call, then one each during polls 1 and 2; poll 3 completed without another post. Static instructions pass arguments `+8/+0xc` as hwnd/message and select lParam from SyncInfo row `cursor-1` at `+0x10` (cursor zero selects row 599). Prior SyncInfo producer/consumer analysis establishes row `+0x10` as the inclusive source-buffer end byte offset. A direct callsite cross-check matched all lParams to the selected row field: cursors 5/8/11/11 selected rows 4/7/10/10 with inclusive spans `0x14..0x18`, `0x22..0x26`, `0x30..0x34`, and `0x30..0x34`; the repeated `0x34` came from row 10 remaining selected on polls 1 and 2 in this run. Cursor-zero wrap to row 599 remains static-only. The PE import table maps IAT `0x1006d110` to `USER32.dll!PostMessageA`; the callsite is `0x1002151b`. The host's use of the endpoint remains unknown. A separate instrumented run captured BOOL `1` at the first PostMessageA return; after disabling that return breakpoint, the inferior exited before the second notification returned. This establishes one call's API return, not complete delivery/consumption. The stable trace captures attempted arguments but does not establish host retrieval or handler interpretation. The void wrapper's raw EAX values are run observations, not declared return values. Short-call evidence is in [`buffer-progress.log`](../../tools/revkit/work/stage16/buffer-progress.log); long-call trace/capture are [`trace-buffer-progress-notifications-target.gdb`](../../tools/revkit/work/stage21/trace-buffer-progress-notifications-target.gdb) and [`buffer-progress-notifications-target-api.log`](../../tools/revkit/work/stage21/buffer-progress-notifications-target-api.log). The source cross-check and field watchpoint are in [`buffer-progress-wparam-source-v3-api.log`](../../tools/revkit/work/stage21/buffer-progress-wparam-source-v3-api.log) and [`buffer-progress-wparam-watch-api.log`](../../tools/revkit/work/stage21/buffer-progress-wparam-watch-api.log); the string-length capture is [`buffer-progress-wparam-value-api.log`](../../tools/revkit/work/stage21/buffer-progress-wparam-value-api.log). The lParam cross-check is [`buffer-progress-syncinfo-lparam-api.log`](../../tools/revkit/work/stage21/buffer-progress-syncinfo-lparam-api.log), reproduced by `run-buffer-progress-syncinfo-lparam.sh`. The source, watchpoint, and value runners are `run-buffer-progress-wparam-source.sh`, `run-buffer-progress-wparam-watch.sh`, and `run-buffer-progress-wparam-value.sh`. The failed initial/v2 attempts are documented in the Stage 21 README. The partial return observation is in [`buffer-progress-notification-return-partial-api.log`](../../tools/revkit/work/stage21/buffer-progress-notification-return-partial-api.log) with its trace and runner. | Host message retrieval/handler behavior for a real HWND, message cadence for other text lengths/voices/options, return values beyond the first observed post, failures, and broader arguments. |
| `VT_TextToLipSyncLog_ENG` | Runtime after model load on `Hello world.`: bytewise heap-backed calls return raw EAX 1 and create reports at exact tested paths: `vtspeak-lipname.txt`, `./vtspeak-lipdot`, `lipsync-probe/rel.txt`, `lipsync-probe/report.txt`, `lipsync-probe/vtspeak-liprelative.txt`, and `Z:/work/stage16/lipsync-filename-bytewise-absolute-z-output.txt`. The supplied path is used as the output path for these extension, relative, subdirectory, and absolute Z-drive cases. At speed 100, eight phone lengths and two word totals sum to 11,803 frames; `VT_TextToBufferEX_ENG` format 0 returns 23,606 PCM bytes = 11,803 16-bit frames. At speed 200, report totals are 5,682 frames and the paired PCM16 call returns 11,364 bytes. Thus units are audio frames and values vary with speed in these matched cases. Null and empty second strings each returned raw EAX 1 and produced no report in the Stage 5 working directory. Static `FUN_1001df10` builds a candidate using `length-sync-%s-%s.txt`; its path argument exactly matching the empty global at `0x1009f948` returns a null context, confirmed by the sentinel runtime call. The later formatted-candidate branch also requires that equality, so it is unreachable for a stable argument through the ordinary public call path; a nonempty path selects the caller-supplied path. The intended use of that candidate branch remains unknown. An earlier absolute-path fault came from the invalid GDB string setup and is not API evidence. The older `vtspeak-lead6` capture retained report content but not its original pathname. Speaker -1 falls back and returns raw EAX 1; null text produced AX -3; empty text produced AX -4 with raw EAX `0x003efffc`. The decompiled wrapper is `void`; register values are not declared C returns. Two tested file-open failures reach the stream helper in mode `wt`: a missing parent produces `CreateFileA` (import slot `0x1006d054`) `INVALID_HANDLE_VALUE` with `GetLastError` 3 (`ERROR_PATH_NOT_FOUND`), and using the existing directory `.` as the target produces the same invalid handle with `GetLastError` 5 (`ERROR_ACCESS_DENIED`). In both cases the null writer field at wrapper `+0x10` is dereferenced unchecked. | UNC paths and other pathname forms; exact `GetLastError` values for path classes beyond the observed missing-parent and directory-target cases; intended use of the formatted-candidate branch; and report/audio effects on other text, speakers, and options beyond the measured matrix. Evidence: [`lipsync-filename-bytewise-extension.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-extension.log), [`lipsync-filename-bytewise-dot-relative.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-dot-relative.log), [`lipsync-filename-bytewise-subdir-relative-verified.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-subdir-relative-verified.log), [`lipsync-filename-bytewise-absolute-z.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-absolute-z.log), [`lipsync-speed-100.log`](../../tools/revkit/work/stage16/lipsync-speed-100.log), [`lipsync-speed-200.log`](../../tools/revkit/work/stage16/lipsync-speed-200.log), [`lipsync-null-path.log`](../../tools/revkit/work/stage16/lipsync-null-path.log), [`lipsync-sentinel.log`](../../tools/revkit/work/stage16/lipsync-sentinel.log), [`lipsync-filename-bytewise-missing-parent-stop.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-missing-parent-stop.log), [`lipsync-missing-parent-null-writer-v3-api.log`](../../tools/revkit/work/stage21/lipsync-missing-parent-null-writer-v3-api.log), [`lipsync-directory-target-null-writer-v4-api.log`](../../tools/revkit/work/stage21/lipsync-directory-target-null-writer-v4-api.log), and [`buffer-ex-0.log`](../../tools/revkit/work/stage16/buffer-ex-0.log). |
| `VTDTTS_MakeInfo_ENG` | Runtime with Stage 16 text, speaker 1, remaining scalar args -1, and a heap-backed second argument completed with raw EAX `1`. It emitted the literal prefix plus `.bin.dtt` and `.asc.dtt`; the heap-path pair byte-matches the original captures. The binary magic is `VTDTTS BINARY\0`; bytes `03 04` and four NUL-terminated bank names (`merged-gen`, `merged-num`, `merged-etc`, `merged-alp`) follow. The ASCII header spells these as `3`, `4`, and the same names; `4` matches the named-bank count, while `3` remains unlabeled. Captured records decode identically in ASCII and binary. Phone rows cross-check to `unit-*.idx` plus their `.dat`/`.upm` spans. Mode 0 selects a whole decoded unit and combined UPM span; mode 1 selects first-side spans; mode 2 selects second-side spans and adds `Shift Size`. The writer computes that field as the 16-bit value at the selected-unit descriptor's `+0x0c` minus the timeline row's `-0x05` word. The subtraction is directly visible in disassembly; calling its result a splice/crop location would go beyond the evidence. `TypeFlag=2` is the detailed phone/unit record. `TypeFlag=1` is observed between OW1 and W for `Hello, world.` (`Size=3200`) and `Hello. World.` / `Hello... World.` (`Size=14800`), but absent for `Hello world.`. The synthesis loop zero-fills `Size*2` bytes, so Size is a count of PCM16 samples and type 1 is an inserted silence interval in these cases. Repeating the comma case with pause 120 retained the same size. Stage 21 confirms zero-based `File Index` mappings 0=`merged-gen`, 1=`merged-num`, 2=`merged-etc`, and 3=`merged-alp`; all four are runtime-observed. Rate captures establish `Pitch_rate` = supplied pitch, `Volume_rate` = supplied volume, and `Duration_rate` = `((speed >> 1)+10000)/speed` with integer division (speed 120 yields 83); pause 120 left the fields unchanged. The Stage 21 upper/lower isolated-letter matrices add 208 calls: `.`, `,`, and `!` yield TypeFlag=2 records identical field-for-field to standalone for all 26 letters; `?` changes the record list in every case (nine gen-only, two etc-only, fifteen mixed gen/etc). All 598 TypeFlag=2 rows cross-check to indexed units and neither 104-case matrix emits TypeFlag=1; upper/lower rows are field-for-field equal for each letter and mark. The 22 phrase/punctuation cases include one TypeFlag=1 row in `A, B`. The format is therefore a selected synthesis-timeline manifest of model unit spans, timing/pitch spans, rate metadata, and observed silent intervals. The four checked-in host executables have no import for this export, and the local binary/report search found no DTT consumer; its intended external client remains unidentified. The record layout is: `u8` type flag; NUL-terminated phone string; `u8` file index; little-endian `u32` PCM position; `u16` PCM size and coded size; `u8` mode; two `u16` pitch endpoints; `u32` PM position; `u16` PM size; then three `u16` rates, with an extra `u16` shift field for mode 2. Static code also checks loaded/nonempty text, applies an invalid-speaker fallback to slot 1, writes both files, and frees the context. The export wrapper is decompiled `void`, so EAX is a raw register observation. See the Stage 16 rate/text/punctuation captures and disassembly around `0x1002ca10`. | Header byte 3; physical meaning of mode-2 shift arithmetic; silence-row placement/duration outside the tested phrase cases; `merged-alp` row selection; broader option/text/error behavior; identity of any external DTT consumer. |
| `VT_VerifyTTS_ENG` | Export `0x1001ded0` is a `void` wrapper over `FUN_100226f0`; static flow applies speaker-slot fallback and loaded/text checks before internal parsing. Runtime matrices cover selected short delimiter strings and compare status with no dictionary, a populated index-27 dictionary at gate 0/1, and empty index 28. | Full markup grammar, dictionary-dependent internal state beyond returned low AX, additional rows/texts, and other voice slots. |

### `VT_TextToLipSyncLog_ENG` option and path probes

The option matrix uses the same loaded Paul engine, speaker 1, `Hello world.`
text, speed 100, and a bytewise-initialized heap filename. Every call that
completed returned raw EAX 1 after the decompiled `void` wrapper. With all
remaining numeric options at `-1`, the two word totals and eight phone lengths
sum to 11,803 frames. Setting pitch to 50 changes the report total to 12,850;
pitch 200 gives 11,998. The phone IDs, labels, and word spans remain the same;
individual phone durations and both word durations change. This is consistent
with the independently paired PCM16 count for the baseline. The pitch-variant
captures measure report output only; no buffer call was paired to those runs.
The report therefore records selected utterance timing, not a static
dictionary of phone durations.

On that fixture, reports for text types 0–7 are byte-identical to the all-`-1`
baseline. Setting volume 0, pause 250, or dictionary 0 individually also leaves
the report unchanged. Pitch 50 combined with volume 0, pause 250, and dictionary
0 produces the same report as pitch 50 alone. Repeating that combination with
text type 4 or 6 returns the report to the baseline bytes, consistent with the
separately observed static behavior that those types reset earlier option
fields. This does not establish that volume, pause, dictionary, or text type
are globally ignored: the report captures timing/labels for one plain-text
fixture, and the output samples were not compared acoustically in this matrix.
The raw return was 1 in every completed case. Captures are
`lipsync-option-*.log` and `lipsync-option-*-report.txt` in Stage 16; compare
the [baseline](../../tools/revkit/work/stage16/lipsync-option-baseline-report.txt),
[pitch 50](../../tools/revkit/work/stage16/lipsync-option-pitch-low-report.txt),
[pitch 200](../../tools/revkit/work/stage16/lipsync-option-pitch-high-report.txt),
and [type-4 extreme control](../../tools/revkit/work/stage16/lipsync-option-type4-extremes-report.txt).

The bytewise absolute path
`C:\windows\temp\vtspeak-lipsync-c-drive-output-verified.txt` returned raw
EAX 1 and created a report under the Wine prefix's `drive_c/windows/temp`.
Together with the prior `Z:/work/...` result, this confirms the tested paths on
both the mapped C: prefix drive and Z: host drive. See
[`lipsync-filename-bytewise-absolute-c-verified.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-absolute-c-verified.log)
and [`lipsync-filename-bytewise-absolute-c-verified-output.txt`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-absolute-c-verified-output.txt).

The explicit relative path `no-such-lipsync-dir/report.txt` with a nonexistent
parent did not create a report. The target process exited during the injected
API call with exception status `0xc0000005`; GDB abandoned the call, so no raw
EAX value was captured. The heap string itself was printed intact before the
call, and working relative subdirectory paths succeed in the matched filename
probes. This is direct evidence of a failure on this invalid-path case under
Wine. A follow-up stopped at EIP `0x10025e4d` in the report writer wrapper;
the instruction is `mov eax, DWORD PTR [edx]`, before a call to `0x10064d92`.
The capture did not retain EDX, so it locates the invalid dereference without
establishing whether the pointer was null or why it was invalid. See
[`lipsync-filename-bytewise-missing-parent.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-missing-parent.log).
The follow-up trace is
[`lipsync-filename-bytewise-missing-parent-stop.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-missing-parent-stop.log).

A Stage 21 target-flow replay of the same missing-parent path stopped
immediately before that instruction. At the `FUN_1001e0c0` caller, the active
state object had a report subobject at `state+0x2c`; its writer field at
subobject `+0x10` was zero. The writer wrapper received that same null pointer
in EDX and dereferenced it without a null check at `0x10025e4d`. This identifies
the immediate invalid value and call chain. The expanded trace shows that
`FUN_1001df10` calls `FUN_10025dc0` with the supplied pathname and mode `"wt"`;
the stream layer parses this into flags `0x4301` and calls `CreateFileA`
through import slot `0x1006d054`. For `no-such-lipsync-dir/report.txt`, that call returns
`INVALID_HANDLE_VALUE`, and the immediately queried `GetLastError` is 3
(`ERROR_PATH_NOT_FOUND`). `FUN_10025dc0` consequently returns zero. Its caller
stores zero at the allocated wrapper's `+0x10`, clears the other four fields,
and nevertheless returns the nonnull wrapper. The later report writer passes
the zero field to `FUN_10025e40`, whose first dereference faults. Thus the same unchecked failure path occurs for both `ERROR_PATH_NOT_FOUND`
and `ERROR_ACCESS_DENIED`. Static disassembly shows that every
`CreateFileA` invalid-handle result follows the same branch: it queries
`GetLastError`, returns a negative low-level result, and the stream opener
converts that to zero without code-specific recovery. Runtime directly confirms
that propagation for error codes 3 and 5; exact codes from other path classes
remain untested. The v2 capture is
[`lipsync-missing-parent-null-writer-v2-api.log`](../../tools/revkit/work/stage21/lipsync-missing-parent-null-writer-v2-api.log),
reproduced by
[`run-lipsync-missing-parent-null-writer-v2.sh`](../../tools/revkit/work/stage21/run-lipsync-missing-parent-null-writer-v2.sh)
with
[`trace-lipsync-missing-parent-null-writer-v2.gdb`](../../tools/revkit/work/stage21/trace-lipsync-missing-parent-null-writer-v2.gdb).
The deeper capture is
[`lipsync-missing-parent-null-writer-v3-api.log`](../../tools/revkit/work/stage21/lipsync-missing-parent-null-writer-v3-api.log),
reproduced by
[`run-lipsync-missing-parent-null-writer-v3.sh`](../../tools/revkit/work/stage21/run-lipsync-missing-parent-null-writer-v3.sh)
and
[`trace-lipsync-missing-parent-null-writer-v3.gdb`](../../tools/revkit/work/stage21/trace-lipsync-missing-parent-null-writer-v3.gdb).
The existing-directory variant is captured in
[`lipsync-directory-target-null-writer-v4-api.log`](../../tools/revkit/work/stage21/lipsync-directory-target-null-writer-v4-api.log),
reproduced by
[`run-lipsync-directory-target-null-writer-v4.sh`](../../tools/revkit/work/stage21/run-lipsync-directory-target-null-writer-v4.sh)
and
[`trace-lipsync-directory-target-null-writer-v4.gdb`](../../tools/revkit/work/stage21/trace-lipsync-directory-target-null-writer-v4.gdb).

A bytewise path containing the Windows separator, `lipsync-probe\backslash.txt`,
returned raw EAX 1 and created the report in the `lipsync-probe` directory as
`backslash.txt`; no POSIX filename containing a literal backslash was left
behind. Thus backslash is treated as a separator for this relative-path case
under Wine. See
[`lipsync-filename-bytewise-backslash-relative.log`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-backslash-relative.log)
and [`lipsync-filename-bytewise-backslash-relative-output.txt`](../../tools/revkit/work/stage16/lipsync-filename-bytewise-backslash-relative-output.txt).

### `VT_TextToPreprocessInfoFile_ENG` isolated flag observations

The baseline flag probes used Stage 5 text `Hello world.`, speaker 1,
synthesis settings `-1`, dictionary index 0, and text type 0. The wrapper is
decompiled as `void`; raw EAX/AX observations do not establish an API return
value. Flag 7 was then repeated with five controlled text fixtures; flags 3,
5, and 10 were checked with `Hi.` and `Hello, world!`. With nonzero flag 5
and otherwise valid arguments, a null text pointer left raw low AX `-3`; an
empty text string left raw low AX `-4`. These error-path register observations
are in [`preprocess-input-errors-api.log`](../../tools/revkit/work/stage16/preprocess-input-errors-api.log).

| Byte flag | Internal mode byte at context `+0x20` (static) | Isolated file effect (runtime) |
| ---: | ---: | --- |
| 0 | Not initialized; early return before pointer validation | No isolated output probe; separate zero-flag probe confirms early return with null/zero arguments. |
| 1 | 1 | One empty file, `preprocess-isolapreprocess-isolated`. |
| 2 | 0 | One empty file, `flag2-preprocessflag2-preprocess-is`. |
| 3 | 2 | CRLF report: one dot per input text byte (excluding its line ending), then a count and decimal internal symbol/control bytes. `Hello world.` gives `12` dots and `11 : 90 34 23 43 48 92 66 27 43 21 90`; `Hi.` gives `3` dots and `4 : 90 34 17 90`; `Hello, world!` gives `13` dots and `11 : 90 34 23 43 48 91 66 27 43 21 90`. The ordinary phone bytes agree with the internal phone-symbol codebook (for example 34=`HH`, 23=`EH0`, 43=`L`, 48=`OW1`, 66=`W`, 27=`ER1`, and 21=`D`). In these captures, `91` occurs at the comma boundary, `92` between word phone runs, and `90` at the sentence start and end positions; the exact general roles of those control bytes remain open. |
| 4 | 2 | Four empty text files. Correct heap-backed path probe created exactly `<param2>.0` through `<param2>.3`; static code formats four `%s.%d` names. |
| 5 | 2 | CRLF phone-symbol rendering. `Hello world.` gives `HH EH0 L OW1 - W ER1 L D .`; `Hi.` gives `HH AY1 .`; `Hello, world!` gives `HH EH0 L OW1 - W ER1 L D !`. Ordinary phone labels use the recovered internal phone-symbol map; the writer converts engine phone bytes to those labels. The numeric flag-3 sequence includes an additional leading boundary value not printed by this renderer. |
| 6 | 4 | Fixed working-directory `test.pcm`, 23,606 bytes (11,803 16-bit samples); the static path opens it as binary output. |
| 7 | 2 | Variable-size CRLF parse/pronunciation report; see the field-level description below. The corrected `Hello world.` capture is 127 bytes, not the earlier 245-byte GDB-string capture. |
| 8 | 4 | One empty file. |
| 9 | 3 | One empty file. |
| 10 | 4 | CRLF text re-emission with punctuation separated: `Hello world .`; `Hi.` becomes `Hi .`, and `Hello, world!` becomes `Hello, world !`. |
| 11 | 4 | One empty file. |
| 12–253 | 4 (static default branch) | One in-process sweep called each byte value with the fixed `Hello world.` fixture and a unique output pathname. All 242 calls reached common completion with raw EAX 1; all created empty files (SHA-256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`). The wrapper is `void`, so raw EAX is not a declared return. This is one call sequence and fixture, not fresh-process or cross-text validation. |
| 254 | 4 | One empty file. |
| 255 | 4 | One empty file. |

The sweep's per-call register observations and output hashes are in
[`preprocess-default-flags-sweep-api.log`](../../tools/revkit/work/stage16/preprocess-default-flags-sweep-api.log);
the generated output inventory is
[`preprocess-default-flags-sweep-files.txt`](../../tools/revkit/work/stage16/preprocess-default-flags-sweep-files.txt),
and the runner is documented in the [Stage 16 README](../../tools/revkit/work/stage16/README.md).

The mode mapping is from the dispatch at `0x1001eff2`–`0x1001f024`: the
original byte is stored at context `+0x21`; flags 1–9 index the table at
`0x1001f1cc`, while the default sets mode 4. The bounds check routes byte
values 10 and above to that default; 10, 11, every value 12–253, 254, and 255
were checked at runtime. The fixed-text sweep results for 12–253 match the
empty-file effect previously observed for 11, 254, and 255. Thus the internal
dispatch and the tested runtime result are mapped, while public names for
these modes and behavior on other texts remain open.

#### Output record interpretation

Static writers and controlled text changes establish these output roles:

* Flag 3 writes one period per input text byte (excluding the file line
  ending), then a count
  and a sequence of decimal internal phone/control bytes. Its ordinary phone
  values match the internal phone-symbol codebook. In the baseline sequence,
  `92` separates the two word phone runs; the leading and trailing `90` values
  bracket the sentence. A comma changes the internal separator to `91` while
  a terminal exclamation mark still ends in `90`. This is a compact engine
  representation, not an independent standard phonetic notation.
* Flag 5 writes readable CMU-style labels for the phone sequence, with `-`
  between word runs and the input terminal punctuation at the end. Ordinary
  phone symbols agree with flag 3 and the recovered byte-to-phone table; flag
  3 also emits a leading boundary value that flag 5 omits.
* Flag 7 writes the engine's per-word parse/pronunciation view. The line
  begins with an engine-selected run of periods, a space, and the source text.
  Each following row contains a normalized word surface, a numeric rendering
  of an internal boundary marker, its phone labels, a letter rendering of a metadata bitmask,
  and inclusive source-byte start/end offsets. For `Hello world.`, for
  example, the ranges `(0 4)` and `(6 10)` select the five letters of each
  word; for `Hello, world!`, the second range shifts to `(7 11)` and excludes
  punctuation. The integer in `word(code)` is generated from the marker byte
  at record offset `+0xa09`: `[` prints as `0`, `` as `1`, `]` as `2`, and
  `^` as `4`; other markers, including the observed `Z`, print as `3`. The
  upstream writer at `0x100130e0` assigns these marker characters across the
  parsed word sequence before `FUN_10017e80` serializes them. This is a
  structural boundary-marker code, not a part-of-speech label; exact names
  for each boundary class remain open. The metadata formatter at `0x10018020`
  emits letters for set bits: bit 0 `G`, bit 1 `E`, bit 2 `Y`, bit 3 `R`, bit
  4 `A`, bit 5 `S`, and bit 6 `T`. Those letters are observed bit labels;
  their expansions/linguistic interpretation are not established. The
  period-prefix length is not tied to the text: five fresh-process repeats of
  the same `Hello world.` input selected 52, 140, 33, 151, and 33 periods. In
  static code, `FUN_1001ce10` prepends a selected string from the eight-entry
  pointer table at `0x1007cf88` plus one space; it obtains an index in 0–7
  through `FUN_1001cd70`, whose first seed comes from `GetTickCount` (IAT
  `0x1006d020`). The table's period counts by index are `108, 111, 140, 52,
  33, 151, 106, 58`. A controlled runtime sweep then forced all eight table
  indexes in isolated API calls and verified each emitted period count:
  index 0 → 108, 1 → 111, 2 → 140, 3 → 52, 4 → 33, 5 → 151, 6 → 106, and
  7 → 58. Each forced call used a state seed chosen so its next Park–Miller
  state had the desired residue modulo 8; the trace records the resulting
  state and raw EAX observation, while the files confirm the prefix bytes.
  The captures are `preprocess-rng-all-prefixes-api.log` and `p0`–`p7` in
  Stage 21. Reproduce with `run-preprocess-rng-all-prefixes.sh`. Together,
  static arithmetic, the natural multi-call sequence, and the controlled
  table-index sweep establish how the prefix is selected and the complete
  table-to-length mapping. They do not establish why the report includes
  this filler or the distribution across independent natural process seeds;
  those questions remain open.
  Disassembly of `FUN_1001cd70` gives the state transition explicitly. A
  one-time flag at `0x1007cfb0` selects `GetTickCount` as the first input;
  subsequent calls use the saved state at `0x1009fc4c`. Each call advances it
  with `S' = 16807*(S mod 127773) - 2836*floor(S/127773)`, adding
  `2147483647` when `S' <= 0`, then selects `S' mod 8` from the table. The
  constants make this the Park–Miller form of a multiplicative congruential
  generator (algorithm identification from the recovered arithmetic). A
  single-process runtime sequence on three identical `Hello world.` inputs
  returned states `938191541`, `1360293313`, and `338805629`; applying the
  formula to the first two states predicts the next state exactly. Their
  residues 5, 1, and 5 select prefixes of 151, 111, and 151 periods, matching
  the three output files. The call outputs and state values are in
  [`preprocess-rng-sequence-v2-api.log`](../../tools/revkit/work/stage21/preprocess-rng-sequence-v2-api.log);
  the files are `s2`, `s3`, and `s4` in Stage 21. The probe does not capture
  the initial `GetTickCount` value itself, but directly confirms the stored
  recurrence and table lookup after seeding. This explains how the filler is
  chosen, not why the report includes it. Recheck the captured transitions
  and output prefixes with
  [`check_preprocess_rng_sequence.py`](../../tools/revkit/scripts/check_preprocess_rng_sequence.py).
* Flag 10 re-emits source text with spaces before terminal punctuation in the
  tested ASCII cases. This behavior is directly shown for `.`, `!`, and a
  comma followed by a terminal mark, but broader punctuation and markup rules
  are not covered.

#### Preprocess report sensitivity to scalar options

A controlled 44-call matrix used `Hello world.`, speaker 1, and
dictionary 0. For each of flags 3, 5, 7, and 10, it tested the all-`-1`
baseline, pitch 50/200, speed 50/400, volume 0/500, pause 0/250/65,535, and a
second baseline call. Every call completed with raw EAX `1` (the wrapper is
decompiled `void`). For each flag, all 11 output files are byte-identical,
including flag 7 after its random prefix was reset to the same state before
each call. This shows that these scalar settings do not change these four
reports for this utterance; it does not establish invariance for other text,
speaker/dictionary states, invalid scalar values, or other text categories.

The separate flag-6 matrix shows the PCM writer does consume some of those
settings. Against the 11,803-sample default, pitch 50/200 produced 12,850 and
11,998 samples; speed 50/400 produced 25,377 and 2,438 samples. Volume 0
produced 11,803 zero samples; volume 500 changed sample values but retained
the 11,803-sample count. Pause 0, 250, and 65,535 preserved the default PCM
byte-for-byte for `Hello world.`. The repeated default also matched exactly.
This matrix uses one period-terminated phrase, so the pause result does not
cover punctuation or other sentence-boundary contexts. Its capture is
[`preprocess-pcm-settings-api.log`](../../tools/revkit/work/stage21/preprocess-pcm-settings-api.log),
with outputs `preprocess-pcm-settings-v00.pcm` through `v10.pcm`; the
reproduction trace and runner are
[`trace-preprocess-pcm-settings.gdb`](../../tools/revkit/work/stage21/trace-preprocess-pcm-settings.gdb)
and [`run-preprocess-pcm-settings.sh`](../../tools/revkit/work/stage21/run-preprocess-pcm-settings.sh).
Verify with
[`check_preprocess_pcm_settings.py`](../../tools/revkit/scripts/check_preprocess_pcm_settings.py).

The call log is
[`preprocess-settings-matrix-api.log`](../../tools/revkit/work/stage21/preprocess-settings-matrix-api.log).
The GDB trace and runner are
[`trace-preprocess-settings-matrix.gdb`](../../tools/revkit/work/stage21/trace-preprocess-settings-matrix.gdb)
and [`run-preprocess-settings-matrix.sh`](../../tools/revkit/work/stage21/run-preprocess-settings-matrix.sh).
The host-side checker is
[`check_preprocess_settings_matrix.py`](../../tools/revkit/scripts/check_preprocess_settings_matrix.py).
For flags 3, 5, 7, and 10 respectively, each 11-file set is 53, 29, 202, and
15 bytes; the checker prints the common SHA-256 for each set.

Flag 6 writes the generated 16-bit PCM stream to fixed `test.pcm`; it ignores
the supplied path. Flag 4 appends `.0` through `.3` to the supplied path.
Heap-backed GDB probes initialize every path byte separately, confirm the
path string before and after the API call, and create files at the exact
expected names for flags 1–5, 7–11, 254, and 255. Thus those flags pass
`param2` through literally; no filename transformation was found in the
wrapper. Earlier malformed names and the earlier 245-byte flag-7 body came
from GDB literal-string setup that altered memory before the call, and are
superseded for filename and content claims. The traces and outputs are listed
in [`Stage 16`](../../tools/revkit/work/stage16/README.md). The repeatable
heap-backed runner is
[`run-preprocess-flag-heap-isolated.sh`](../../tools/revkit/work/stage16/run-preprocess-flag-heap-isolated.sh);
controlled text cases use
[`run-preprocess-text-case.sh`](../../tools/revkit/work/stage16/run-preprocess-text-case.sh).
Five identical-input captures preserve the observed prefix variation:
[`repeat-a`](../../tools/revkit/work/stage16/heap-flag7-case-repeat-a-20260926),
[`repeat-b`](../../tools/revkit/work/stage16/heap-flag7-case-repeat-b-20260926),
[`repeat-c`](../../tools/revkit/work/stage16/heap-flag7-case-repeat-c-20260926),
[`repeat-d`](../../tools/revkit/work/stage16/heap-flag7-case-repeat-d-20260926),
and [`repeat-e`](../../tools/revkit/work/stage16/heap-flag7-case-repeat-e-20260926).
The null/empty text check is reproducible with
[`run-preprocess-input-errors.sh`](../../tools/revkit/work/stage16/run-preprocess-input-errors.sh).

### `VTDTTS_MakeInfo_ENG` binary record evidence

The checked-in disassembly identifies the exported implementation at
`0x10022350`. The observed binary file is
[`heap-makeinfo-path-20260926.bin.dtt`](../../tools/revkit/work/stage16/heap-makeinfo-path-20260926.bin.dtt).
In that file, record 1 starts at byte 60 and its fields align with the ASCII
row for phone `HH` as follows. Offsets are decimal from the beginning of the
file; integer byte order is little-endian, established by matching decoded
values against the ASCII output.

| Offset | Width | ASCII label / interpretation supported by output | `HH` bytes/value |
| ---: | ---: | --- | --- |
| 60 | 1 | TypeFlag | `02` / 2 |
| 61 | variable + NUL | PhoneString | `HH\0` |
| 64 | 1 | File Index | `00` / 0 |
| 65 | 4 | PCM Pos | `00000000` / 0 |
| 69 | 2 | PCM Size | `ee05` / 1518 |
| 71 | 2 | PCM Coded Size | `2003` / 800 |
| 73 | 1 | Phone Mode | `00` / 0 |
| 74 | 2 | Pitch_first | `6c00` / 108 |
| 76 | 2 | Pitch_last | `6800` / 104 |
| 78 | 4 | PM Pos | `00000000` / 0 |
| 82 | 2 | PM Size | `0d00` / 13 |
| 84, 86, 88 | 2 each | Pitch_rate, Duration_rate, Volume_rate | 100, 100, 200 |

Mode 2 inserts a two-byte `Shift Size` after the mode byte. For example, the
`OW1` mode-2 row starts at byte 181, with its shift size at 196–197; the
following pitch endpoints begin at 198. Decoding all ten records using this
conditional width consumes the complete 365-byte binary file, and each
decoded value matches its corresponding ASCII field. This validates the
serialization for this capture. Cross-checking all ten rows against
`unit-gen.idx` shows that `PCM Pos` and `PCM Coded Size` are the selected
`.dat` offset and encoded span length. `PM Pos` and `PM Size` select the
matching combined UPM span for mode 0, the first-side span for mode 1, and the
second-side span for mode 2. `PCM Size` matches the complete decoded sample
count for mode 0 and the index's corresponding first/second sample span for
modes 1/2. In all ten rows, the pitch endpoints equal the matching UPM edge
periods doubled to the synthesis sample grid. These checks identify the
records as unit selections and the associated waveform/timing spans for this
utterance; they do not establish what downstream application consumes the
DTT files. The one-based DTT row mappings are:

| Row | Phone | Mode | `unit-gen.idx` unit | PCM size selection | UPM selection |
| ---: | --- | ---: | ---: | --- | --- |
| 1 | HH | 0 | 0 | Full unit | Combined span |
| 2 | EH0 | 0 | 1 | Full unit | Combined span |
| 3 | L | 0 | 2 | Full unit | Combined span |
| 4 | OW1 | 1 | 165222 | First-side span | First-side span |
| 5 | OW1 | 2 | 369960 | Second-side span | Second-side span |
| 6 | W | 0 | 369961 | Full unit | Combined span |
| 7 | ER1 | 1 | 420365 | First-side span | First-side span |
| 8 | ER1 | 2 | 420365 | Second-side span | Second-side span |
| 9 | L | 0 | 420366 | Full unit | Combined span |
| 10 | D | 0 | 420367 | Full unit | Combined span |

The original `Hello world.` output's `File Index` is zero in every row, and
those records come from `merged-gen`. Controlled text captures extend the
mapping: `42.` selected rows from `merged-num` with index 1; in `Hello 42.`,
the initial HH/EH0/L selections came from `merged-etc` with index 2, subsequent
general phones from `merged-gen` with index 0, and number phones from
`merged-num` with index 1. `A.B.C.` remained on `merged-gen`; all six rows
use `File Index` 0. Across all 15 paired Stage 16 ASCII captures (143 rows),
indices 0, 1, and 2 occur, while 3 does not. The static path narrows the open
question: per-phone candidate records carry a bank index; the generic
candidate routines use it to index the bank descriptor table at
`selector * 0x3c0`, try `FUN_10024060`, and then use `FUN_100242a0` as a fallback.
`FUN_1002c220` serializes the candidate record's bank byte into `File Index`;
the loader has a descriptor for each of the four named banks. This shows that
3 is representable by the serialized path and covered by generic selector
logic. At this point in the Stage 16 work, the input/context that sets
selector 3 was unknown; the Stage 21 runtime probes below resolve its
reachability.
The Stage 21 isolated runtime matrix resolves selection of index 3. It called
MakeInfo on each standalone uppercase letter A–Z, repeated A–Z in lowercase,
and used four controls (`ABC`, `A B C`, `A. B. C.`, and `U.S.A.`). Uppercase
and lowercase calls selected the same bank for each letter. Fourteen letters
selected index 3 (`merged-alp`): A, B, E, G, I, J, M, N, Q, R, U, V, W, and Z.
Nine selected index 2 (`merged-etc`): C, D, F, H, K, O, P, S, and T. L, X,
and Y selected index 0 (`merged-gen`). All four controls selected index 0.
Across those 56 letter/control calls, 165 TypeFlag=2 rows occurred: 41 gen, 44 etc, and
80 alp, with no num rows. Every row's `PCM Pos` matches a DAT offset and a
unit ordinal in the unit index named by its `File Index`. The alp rows' unit
ordinals follow alphabetic record order for the selected letters (A=0,
B=1–2, E=7, G=10–11,
I=14, J=15–16, M=21–22, N=23–24, Q=28–30, R=31–32, U=37–38, V=39–40,
W=41–47, Z=53–54). The reproducible analysis is
[`analyze_makeinfo_alphabet_matrix.py`](../../tools/revkit/scripts/analyze_makeinfo_alphabet_matrix.py);
the isolated runtime traces and raw captures are recorded in Stage 21. This
proves `merged-alp` selection is reachable and letter-dependent, and that
case does not change the tested standalone mappings. The separate 22-case
punctuation/context probe found that `.`, `,`, and `!` preserve the standalone
bank for A/B/C, while `?` maps A to `etc`, B to `gen`, and C to a mix of gen
and etc rows. `A B`, `A B C`, and `B C D` use gen; `B C` uses etc; `A, B`
mixes gen for A and alp for B. The full 78 captures contain 231 TypeFlag=2
rows: 73 gen, 60 etc, and 98 alp; one additional TypeFlag=1 silence row
occurs in `A, B`. Every TypeFlag=2 row maps to the selected bank's unit-index
record.
Uppercase and lowercase A–Z terminal-punctuation matrices (208 calls total)
found that period, comma, and exclamation produce TypeFlag=2 records
field-for-field equal to each corresponding standalone letter (52/52
comparisons per mark). Question mark changes every letter's record list in
both cases: nine captures contain only gen rows, two only etc rows, and
fifteen mix gen and etc. For each punctuation mark, uppercase and lowercase
records are also field-for-field equal across all 26 letters. Neither
104-capture matrix contains a TypeFlag=1 row. All 208 calls returned raw EAX
1 and their 598 TypeFlag=2 rows map to indexed units. These results establish
the tested single-letter punctuation and case behavior at default options;
they do not explain why standalone letters or question-mark cases choose
their banks, or how silence intervals are placed in broader phrase contexts.
The Stage 21 ordered-pair probes called MakeInfo for all 676 two-letter
space-separated sequences in uppercase and lowercase, each without terminal
punctuation and with `?` (2,704 successful calls). In both cases, every
plain/question pair changes its TypeFlag=2 row list; 430 pairs also change the
set of selected banks, while 246 retain the set but change unit records.
Uppercase pair rows total 3,429 gen, 106 num, and 174 etc without `?`, and
3,194 gen, 76 num, 833 etc, and 6 alp with `?`. Lowercase totals are 3,426
gen, 110 num, and 184 etc without `?`, and 3,198 gen, 80 num, 837 etc, and
6 alp with `?`. The six alp rows in each case appear only in five
question-mark pairs: `E P?`, `Q J?`, `T F?`, `T S?`, and `V P?`. `merged-num`
occurs in 54 uppercase plain and 43 uppercase question-mark outputs, versus
60 and 46 lowercase outputs. No pair capture emitted TypeFlag=1. Across
1,352 uppercase/lowercase comparisons, 1,302 record lists are field-for-field
equal; all 50 differences begin with A/a. The first phone row is `AH0` for
uppercase A and `EY1` for lowercase a; 17 cases also change bank sets, while
33 change records within the same set. Every selected phone row's PCM offset
maps to an index record in its named bank. Two additional 676-pair matrices
tested `Upper lower` and `lower Upper` spellings, again with and without
terminal `?` (2,704 calls). All 1,352 outputs in each mixed-case pattern match
the all-uppercase output field-for-field. Therefore, changing only the first
or only the second token to lowercase has no effect in this matrix; the
all-lower differences require both tokens lowercase and are confined to pairs
beginning with `a`. This establishes an interaction between case positions
for the tested inputs without identifying a parser or lexical rule. Phrases
longer than two letter tokens remain unprobed.
The exact 17 all-lowercase pairs/forms whose bank sets differ from uppercase
are `aa-plain` (gen → gen+etc), `ad-plain` (gen → etc), `ad-question`
(gen+etc → etc), `ae-question` (gen+etc → gen), `af-plain` (gen → gen+num),
`al-plain` (gen → gen+num), `al-question` (gen+etc → gen+num+etc),
`an-plain` (gen → gen+num), `ao-plain` (gen+num → gen), `aq-question`
(gen → gen+etc), `as-plain` (gen → gen+num), `as-question` (gen → gen+etc),
`at-plain` (gen → gen+num+etc), `at-question` (gen+etc → gen+num+etc),
`av-plain` (gen → num+etc), `av-question` (gen+etc → gen+num+etc), and
`ax-plain` (gen → gen+num). The other 33 differing pair/forms preserve the
bank set but select different rows; the Stage 21 analyzer enumerates them.
Bank indices are 0=gen, 1=num, 2=etc, and 3=alp.
The next 8,112 calls tested terminal period, comma, and exclamation mark on
every ordered pair under all four case patterns. Uppercase period/comma and
lowercase comma/exclamation preserve every unpunctuated row list. Uppercase
exclamation changes only `AA!` (without changing its bank set). Lowercase
period changes the row lists for `aa`, `ae`, `ai`, `ao`, `au`, and `ay`; only
`aa` and `ao` change bank sets. Each mixed-case pattern preserves every pair
row list for period and comma, and changes only `aa!` for exclamation. For
each lowercase-period exception, the resulting rows exactly match the
corresponding uppercase unpunctuated pair, so the ending removes the
lowercase-only variation observed for the plain pair. These are record-list
comparisons against the same-case plain pair.

Against uppercase results, all mixed-case terminal-punctuation captures match
field-for-field (2,028/2,028 per pattern). All-lowercase differs from uppercase
in 69/2,028 comparisons: 19 pairs with period, 25 with comma, and 25 with
exclamation. Every differing pair begins with `a`. For period, the second
letters are `b c d f h j k l m n p q r s t v w x z`; for comma and
exclamation they are every letter except `g`. Thus lowercasing either token
alone still reproduces uppercase output for these punctuation endings. None
of the 8,112 captures emitted TypeFlag=1. Across all 13,806 Stage 21 captures,
76,664 TypeFlag=2 rows cross-check to a DAT offset and unit ordinal in the
selected bank (67,874 gen, 2,022 num, 6,406 etc, 362 alp); the only TypeFlag=1
row remains the one in `A, B`.
The final Stage 21 extension tested 5,408 three-token calls: all 676 possible
two-letter continuations after initial `A`/`a`, in uppercase, lowercase,
`A b c`, and `a B C` patterns, each plain and with terminal `?`. All returned
raw EAX 1, and every detailed row offset maps to the selected bank's unit
index. For uppercase and both mixed patterns, every plain/question comparison
changes rows: 415 also change the bank set and 261 keep the same set while
changing records. In all-lowercase, 394 comparisons change bank sets and 282
change records within the same sets. Against uppercase, 649/676 all-lowercase
suffixes differ in each form (1,298/1,352 outputs); the only exact matches in
both forms are suffixes `ga` through `gz` and `ne`. Both mixed-case patterns
match uppercase for every suffix and form. No three-token capture emits
TypeFlag=1. These results extend the observed case interaction through one
three-token shape beginning with A; they do not determine arbitrary
three-token contexts or a general selector/tokenization rule.

The A-leading extension added the four masks not in the first run: `LLU`
(`a b C`), `LUL` (`a B c`), `ULU` (`A b C`), and `UUL` (`A B c`), each plain
and with terminal `?` for all 676 continuations. The initial LUL and ULU runs
used an incorrect middle-token case mapping; those outputs are retained under
Stage 21 `invalid-case-map/` and excluded from the aggregate. Corrected runs
replace both masks in the 46,254-capture analysis. Every retained call returned
raw EAX 1; every detailed row offset maps to its selected-bank unit, and none
emitted TypeFlag=1. Against UUU, LLL and
LLU each differ in 1,298/1,352 outputs; ULU differs in 50/1,352, while ULL,
LUU, LUL, and UUL match UUU exactly. ULU differs for the same 25
continuations in both forms (`aa`–`af` and `ah`–`az`; `ag` is unchanged),
with 9 plain and 7 question outputs changing bank sets and the rest changing
records within the same bank set. The ULU difference shows a case interaction
even with uppercase A; the exact rule remains bounded to this fixed
three-token shape.

The next Stage 21 extension held the full `B A C` or `B C A` shape constant
while testing all 676 ordered letter pairs around A under all eight per-token
uppercase/lowercase masks, both plain and with terminal `?` (21,632 calls).
Every call returned raw EAX 1; the analyzer validated each TypeFlag=2 PCM
position against the unit index for its selected bank. No call emitted
TypeFlag=1. The plain-to-question split depends on position and case mask. For
`B A C`, UUU changes banks in 352/676 comparisons and changes same-bank rows
in 324; LLL is 351/325, LLU is 391/285, and ULU is 388/288. The other four
masks match UUU's split. For `B C A`, UUU is 141 bank changes and 535 same-bank
row changes; LLL is 143/533, LLU is 140/536, and ULU is 137/539. The other
four masks match UUU's split.

Comparing output records with UUU isolates larger token-case interactions.
For `B A C`, LLL differs in 52/1,352 outputs, LLU in 1,302, and ULU in 1,300;
the other five masks match all 1,352 outputs. For `B C A`, LLL differs in
50/1,352, LLU in 100, and ULU in 52; the other five masks match exactly.
Thus the observed casing interaction depends jointly on which tokens are
lowercase and where A occurs, rather than on all-lowercase text or A's case
alone. This closes capitalization coverage for these two fixed three-token
shapes. It does not establish behavior for arbitrary three-token contexts or
free-text and other longer strings.

A separate four-token follow-up contains 464 captures. The first 64 tested
`A A A D` and `A A G D` across all 16 per-token case masks, plain and with
terminal `?`; the other 400 swept the final token across A–Z for both contexts
under `UUUU`, `UUUL`, `ULUU`, and `ULUL`, using the D results from the first
set. Every call returned raw EAX 1; all 4,952 TypeFlag=2 positions map to units
in their selected banks, and no TypeFlag=1 row appeared. Relative to UUUU,
`A A A D` has seven masks with identical rows in both forms, seven with
same-bank record changes, and two with bank-set changes. `A A G D` has twelve
exact masks and four same-bank record changes. In both D contexts, terminal
`?` changes the bank set for all 16 masks. Across every tested final letter X,
ULUU/ULUL (`A a A X`/`A a A x`) differs from UUUU in both forms, while the
same masks for `A a G X`/`A a G x` match UUUU. For the first context, plain
outputs change banks only at X=a,b; question outputs change banks at
X=b,e,h,j,m,n,s, with the other differences retaining their bank sets. This
carries the selected three-token contrast across each tested fourth letter,
but does not establish a general four-token rule.

The initial 5,408-call matrix tested all 676 ordered `X Y` pairs in `A A X Y` and `A a X Y` under UUUU, UUUL, ULUU, and ULUL, plain and with terminal `?`. Three more batches filled the other 12 case masks, yielding 21,632 calls and 43,264 paired ASCII/BIN captures across all 16 masks, separate from the 46,254-capture aggregate. Every call returned raw EAX 1; all 220,076 TypeFlag=2 rows resolve to indexed units in their selected banks; no TypeFlag=1 row appeared. Under every mask, all 676 plain/question comparisons changed records. Bank-set / same-bank record change counts: UUUU, UUUL, UULL, LUUU, LUUL, LULL: 402/274; UULU and LULU: 403/273; ULUU and ULUL: 416/260; ULLU: 404/272; ULLL: 403/273; LLUU, LLUL, LLLU, and LLLL: 410/266.

Across all pairs and both forms, masks partition into exact field-for-field equivalence classes: `UUUU, UUUL, UULL, LUUU, LUUL, LULL`; `UULU, LULU`; `ULUU, ULUL`; `ULLU`; `ULLL`; `LLUU, LLUL`; `LLLU`; and `LLLL`. Thus lowering token 1 alone is inert when token 2 remains uppercase, while its effect depends on token 2's case in other masks. Lowering token 4 is inert for some classes and changes outcomes in others. Compared with UUUU, UULU/LULU differ for 25/676 pairs per form (plain: 9 bank-set and 16 same-bank changes; question: 6 and 19); ULLU differs for 26/676 per form (plain: 10 and 16; question: 6 and 20); ULLL differs in one plain pair by bank set and one question pair within the same set. ULUU/ULUL differ in 650/676 pairs per form; the only exact pairs are `ga` through `gz`. LLUU/LLUL and LLLU/LLLL differ on all pairs. These results exhaust casing only for the selected `A A` and `A a` prefixes; they do not establish other initial-token or free-text rules.

A complementary four-token matrix varies both initial letters over all 676
ordered pairs `X Y`, followed by fixed `A A`, under all 16 case masks and both
plain and terminal-question forms. This is separate from the preceding
fixed-prefix matrix and the 46,254-capture aggregate. All 21,632 calls returned
raw EAX 1 and produced 43,264 paired ASCII/BIN files. The 204,466 TypeFlag=2
rows all map to indexed units at their selected-bank PCM offsets; no TypeFlag=1
row appeared. The exact case-mask equivalence classes are `UUUU, UUUL, UULL,
ULLL, LUUU, LUUL, LULL`; `UULU, ULLU, LULU`; `ULUU, ULUL`; `LLUU, LLUL`;
`LLLU`; and `LLLL`. Every terminal-question comparison changes records for all
676 pairs; the bank-set/same-bank change counts by group are respectively
180/496, 30/646, 191/485, 186/490, 31/645, and 175/501. These equivalences
describe the tested first-two-letter identity with this fixed suffix only;
other suffixes, non-letter prefixes, punctuation beyond terminal `?`, free
text, and general option/error behavior remain untested.

The aggregate Stage 21 corpus now has 46,254 captures and 322,958 TypeFlag=2
rows: 289,681 gen, 10,689 num, 21,948 etc, and 640 alp. The sole TypeFlag=1
row remains `A, B`.
Thus the observed index is zero-based bank order, with all four bank mappings
now demonstrated at runtime. `03 04` is the observed
binary header after the magic; `04` matches the four names that follow. The
writer at `0x10022513` explicitly stores the immediate byte `3` in the binary
stream and formats the same literal into the ASCII header. It then emits the
bank count and iterates that many bank names. All 15 paired Stage 16 captures
retain ASCII header `3`, `4` and binary bytes `03 04`; the read-only check is
[`check_makeinfo_header.py`](../../tools/revkit/scripts/check_makeinfo_header.py).
Because `3` is a fixed field immediately following the `VTDTTS BINARY` magic
and preceding the bank count, a DTT format revision/variant marker is the
leading interpretation. No field label, DTT reader, or format specification
in the inspected package confirms that meaning, so it remains an inference.

The rate matrix held text and speaker constant while changing one API setting
at a time. Default rows contained `Pitch_rate=100`, `Duration_rate=100`, and
`Volume_rate=200`; pitch 120 changed only the first value to 120, volume 120
changed only the last to 120, and pause 120 changed none of them. Speed 120
set `Duration_rate=83`. The serializer at `0x1002cd4d` computes this field as
`((speed >> 1) + 10000) / speed` using signed integer division. This makes it
an inverse-normalized duration/speed control, rather than the raw speed value;
the matrix establishes its relationship to the API setting, not physical time
units. Record selections and modes stayed the same across these rate calls.

Rows in the default/rate captures and the controlled text captures (`42.`,
`A.B.C.`, and `Hello 42.`) had `TypeFlag=2`, the detailed phone/unit record.
A targeted punctuation matrix then produced type 1: `Hello world.` had none;
`Hello, world.` inserted one between the OW1 row and W with `Size=3200`; and
`Hello. World.` and `Hello... World.` each inserted one at the same boundary
with `Size=14800`. Repeating the comma case with the pause argument set to 120
kept the same 3200 size. These captures show that punctuation can make the
timeline include a type-1 interval. The scalar pause behavior is now measured
for the tested contexts below; the broader punctuation/case rule remains open.

A Stage 21 follow-up exhaustively tested all 676 ordered pairs of one-letter
tokens around interword comma, period, and ellipsis punctuation, under `UU`,
`UL`, `LU`, and `LL` case masks (8,112 calls total). Each call returned raw
EAX 1; all 44,185 TypeFlag=2 rows map to selected-bank unit-index offsets.
There are 5,458 TypeFlag=1 rows across the matrix, always one per affected
capture, with at least one TypeFlag=2 row on each side. `X, Y` emits one
size-3,200 row for all pairs and masks. `X. Y` emits none for uppercase-leading
forms and one size-14,800 row for 675/676 lowercase-leading pairs; `n. e` is
the exception in both `LU` and `LL`. `X... Y` emits one size-14,800 row for all
lowercase-leading pairs, none in `UL`, and 52/676 in `UU`: exactly the cases
whose second token is `A` or `I`. At the established 16 kHz rate, these sizes
correspond to 200 ms and 925 ms. These observations establish punctuation and
case effects for the single-letter interword shape. The runner and combined
index validator are in Stage 21.

A 12-word follow-up tested all 144 ordered pairs from `a`, `i`, `hello`,
`world`, `hi`, `kate`, `paul`, `good`, `morning`, `weather`, `today`, and
`voice`, with comma, period, and ellipsis between tokens under `UU`, `UL`,
`LU`, and `LL` casing (1,728 calls). All returned raw EAX 1; all 3,456 paired
captures contained phone rows, and each silence row was between phone rows.
Comma emitted a 3,200-sample row in 575/576 cases; the sole omission was
uppercase `PAUL, HI`. Period emitted 14,800-sample rows in 401/576 cases,
while ellipsis did so in 423/576. Per-mark, per-mask counts and all 144-pair
case-signature distributions are in the Stage 21 README and reproducible with
`analyze_makeinfo_word_interword_grid.py`. Word identity changes the case
pattern: upper-leading forms sometimes emit silence, and some lower-leading
period forms omit it. The measured matrix therefore does not support carrying
the one-letter trigger rule over to ordinary words.

A further 6,912-call grid moved whitespace around the same punctuation while
holding the word pairs and case masks fixed. All comma layouts were byte
identical in both DTT formats (576/576 per layout), and attached versus
adjacent ellipses also matched in all 576 cases. Giving an ellipsis a
preceding space made a TypeFlag=1 row appear in every pair: 89 `UU` and 64
`UL` cases changed to 3,200 samples; the other rows remained at 14,800. The
phone rows also differ in those same 153 newly-present cases; the other 423
outputs are byte-identical in both DTT formats.
For ellipses, spacing after the mark is inert in these inputs: `both_space`
and `before_only` match in all 576 cases, as do `after_space` and `adjacent`.
Periods are more spacing-sensitive: `left .right` emitted no silence row in
any mask, while `left.right` emitted a 14,800-sample row for every `LU` pair
and none for the other masks. The complete per-layout counts show that
spacing changes some phone rows as well as silence placement. Full ASCII and
binary equivalence counts and the replay validator are in the Stage 21 README
and `analyze_makeinfo_word_spacing_grid.py`.

A further 3,456-call matrix independently varied the bytes before and after
interword punctuation: no separator, one or two spaces, TAB, LF, and CRLF, on
eight ordered word pairs and all four case masks. All calls returned raw EAX
1 and produced 6,912 paired ASCII/BIN captures. Every comma layout
byte-matches its no-prefix/single-space control; the `PAUL, HI`/`UU` case
remains the single comma interval omission. For periods, separator type and
position both affect silence and phone records. With no separator before the
period, a following LF yields ten 3,200-sample and 22 14,800-sample intervals,
while a following CRLF yields 32 intervals of 14,800 samples. A preceding LF
or CRLF with no following separator yields 32 intervals of 3,200 samples.
The Stage 21 README gives the complete 3×6 period count matrix. For ellipses,
with one space after the mark, adding any tested separator before it changes
the seven UU/UL cases that previously lacked silence: each gains a
3,200-sample interval and changes its phone rows. Two spaces or CRLF after the
ellipsis instead makes all 32 intervals 14,800 samples. The cases sample
earlier exceptions and controls;
they establish neither a general tokenizer nor behavior for arbitrary text.
`analyze_makeinfo_ascii_whitespace_grid.py` checks all coordinates, return
values, paired files, row placement, sizes, and byte/phone equivalence classes.
A further 576 calls tested vertical tab, form feed, and bare carriage return
one side at a time. Before punctuation, each exactly matched the tested
space/two-space/TAB-before layouts when followed by one space. After
punctuation, each exactly matched the no-prefix/single-space control in both
DTT formats for all selected cases. This confirms those byte-specific
equivalences on the same eight pairs; it does not extend the result to other
text shapes. The supplemental replay and validator are
`run-makeinfo-ascii-whitespace-class-grid.sh` and
`analyze_makeinfo_ascii_whitespace_class_grid.py`.

The eight-pair boundary was extended across all 144 ordered pairs of the
twelve-word Stage 21 inventory (`a`, `i`, `hello`, `world`, `hi`, `kate`,
`paul`, `good`, `morning`, `weather`, `today`, `voice`), retaining all three
marks, four case masks, and both separator positions. This 10,368-call sweep
captured VT, FF, and bare CR independently (20,736 paired captures); every
call returned raw EAX `1`. Within each position/mark, all three bytes produce
byte-identical ASCII and binary DTT output for all 576 word/case combinations.
Against the canonical one-space-after-mark output, all 576 cases match after
punctuation and before commas; before periods 552 match, and before ellipses
423 match. Per-byte silence counts also agree: before punctuation, comma
`575/0/1`, period `0/405/171`, and ellipsis `153/423/0`; after punctuation,
comma `575/0/1`, period `0/401/175`, and ellipsis `0/423/153` (3200-frame /
14800-frame / absent TypeFlag=1 rows). The larger inventory confirms an
output-equivalence class for VT/FF/CR over these two-word inputs, but does not
establish behavior for arbitrary text, punctuation placement, or longer
sequences. The exhaustive coordinates and byte comparisons are checked by
`analyze_makeinfo_ascii_whitespace_class_grid_all_pairs.py`; the replay is
`run-makeinfo-ascii-whitespace-class-grid-all-pairs.sh` in the Stage 21
evidence directory.

A second all-pairs sweep tested six separator layouts independently before and
after the mark: none, one space, two spaces, TAB, LF, and CRLF. It crossed all
144 ordered pairs from the twelve-word inventory above, the three punctuation
marks, and four case masks: 62,208 calls and 124,416 paired ASCII/BIN captures.
Every call returned raw EAX `1`. The analyzer checked the complete coordinate
set, paired captures, phone rows, and the optional internal TypeFlag=1 row;
both sandbox input/output fixtures were restored byte-for-byte. Full
per-layout counts are in the [Stage 21 matrix report](../../tools/revkit/work/stage21/README.md#six-class-ascii-whitespace-matrix-across-all-word-pairs).

For commas, all 36 layout pairs produce byte-identical ASCII and binary DTT
outputs across all 576 pair/case combinations. Every layout contains a
3,200-sample TypeFlag=1 row in 575 cases; `PAUL, HI`/`UU` is the single
omission. The 14,800-sample size never occurs for this comma matrix.

For periods, whitespace bytes and their side affect both silence and phone
rows. Counts are reported as `3200 / 14800 / absent` per 576 pair/case cases.
No prefix and no following separator yields `0 / 144 / 432`; a preceding
space, two spaces, or TAB with no following separator yields no interval in
all cases. A preceding LF or CRLF with no following separator yields 576
3,200-sample rows. Two spaces or CRLF after the period yields 576
14,800-sample rows regardless of the preceding layout. With one space or TAB
after it, `none` before matches the canonical output in 576 cases; a preceding
space/two spaces/TAB matches in 552, and preceding LF/CRLF matches in 401.
The count matrix distinguishes LF from CRLF and shows other combinations in
the linked report; these are measured output rules for this finite inventory.

For ellipses, all preceding layouts followed by two spaces or CRLF produce 576
14,800-sample rows. With no preceding separator and no separator/one space/TAB
after the mark, the counts are `0 / 423 / 153`; with no prefix but LF after,
the counts are `153 / 423 / 0`. Any nonempty preceding layout followed by
none, one space, or TAB gives `153 / 423 / 0`. In the 153 cases that omit the
row in the no-prefix/no-separator, one-space, or TAB controls, adding the
leading separator changes phone rows as well as inserting a 3,200-sample
interval. With LF after the mark, all preceding layouts yield the same
`153 / 423 / 0` counts. Full ASCII+BIN comparisons partition the 36 layouts
into three global equivalence groups: no prefix with none/space/TAB after; any
preceding layout with two spaces/CRLF after; and the remaining layouts. This
exhausts the chosen two-word matrix; it does not establish a general text
parser rule for arbitrary token strings or longer utterances.

The reproducible runner supports contiguous `START_CASE`/`CASE_LIMIT` chunks;
the capture-set validator is
`analyze_makeinfo_ascii_whitespace_all_pairs.py` in the Stage 21 workspace.

The Stage 21 scalar-pause matrix varied MakeInfo's seventh argument across 108
calls: `-1`, `0`, `1`, `119`, `120`, `121`, `249`, `250`, `251`, `1000`,
`65534`, and `65535` on nine punctuation/control texts. A 27-call follow-up
tested `-2`, `-100`, `-2147483648`, `65536`, `65537`, `100000`, and
`2147483647` on three contexts that produce a silence row. Every call returned
raw EAX 1. In lowercase-period `a. b`, lowercase-ellipsis `a... b`, and
uppercase-ellipsis `A... A`, positive values set the interval size to 16
sample frames per unit (120→1,920; 250→4,000; 1,000→16,000). This is
milliseconds converted at 16 kHz. Zero omits the TypeFlag=1 row. The tested
negative values, including `INT_MIN`, retain the default 14,800-frame
(925 ms) interval. Values 65,535 through `INT_MAX` saturate at 65,535 units,
or 1,048,560 frames. Removing the TypeFlag=1 row makes all other ASCII
capture bytes identical across pause values in each triggering context.

The comma context `A, B` always retains its 3,200-frame (200 ms) interval,
including at zero and above the cap. Period/ellipsis controls that do not
produce a TypeFlag=1 row at default also remain without one for every tested
pause value. Therefore the argument sets the duration only after the text
context has selected an interval whose per-unit pause is unset; it does not
create the punctuation/case trigger. Captures and replay scripts are recorded
in Stage 21 as `makeinfo-pause-silence-*` and `makeinfo-pause-edges-*`.

The writer at `0x1002ca10` serializes type 1 using only its type and size. The
synthesis loop at `0x1002bd90` checks the same row discriminator and fills the
corresponding PCM region with zeros. It computes the byte count as stored
`Size * 2`, identifying `Size` as a PCM16 sample count for this row. Thus
type 1 is a silent timeline interval, not a phone/model record. Type 2 remains
the selected phone/unit record.

The punctuation captures and replay script are in Stage 16:
[`makeinfo-silence-cases-api.log`](../../tools/revkit/work/stage16/makeinfo-silence-cases-api.log),
[`mi-silence-comma-20260926.asc.dtt`](../../tools/revkit/work/stage16/mi-silence-comma-20260926.asc.dtt),
[`mi-silence-sentence-20260926.asc.dtt`](../../tools/revkit/work/stage16/mi-silence-sentence-20260926.asc.dtt),
and [`run-makeinfo-silence-cases.sh`](../../tools/revkit/work/stage16/run-makeinfo-silence-cases.sh).

The phone modes describe record selection, not yet named high-level phonetic
operations. Mode 0 selects the complete decoded unit and combined UPM span;
mode 1 selects the index's first-side PCM/UPM spans; mode 2 selects its
second-side spans and adds the `Shift Size` word. The writer computes that
word as the 16-bit value at the selected-unit descriptor offset `+0x0c` minus
the timeline row word at offset `-0x05` relative to its row pointer. This is
the exact arithmetic in the disassembly. Cross-referencing the selected-unit
descriptor with the 19-byte `unit-*.idx` record identifies descriptor `+0x0c`
as its first-side PCM sample span (index-record bytes 4–5). The subtracted row
word is `Pitch_first`, a doubled UPM edge period already mapped to the 16 kHz
sample grid. `Shift Size` is therefore an exact sample count: first-side PCM
sample span minus `Pitch_first`. A read-only checker matched this relation for
all 24 mode-2 row appearances across 13 ASCII captures, representing 7 unique
bank/unit/value tuples. OW1 gives `1360 - 94 = 1266`; ER1 gives
`1042 - 140 = 902`. This identifies the operands and units but not the
downstream operation that consumes the count; the field name alone does not
prove a crop or splice action. Reproduce the cross-check with
[`check_makeinfo_shift_size.py`](../../tools/revkit/scripts/check_makeinfo_shift_size.py).

The package's four host executables (`voicetext_paul`, `voicetext_julie`,
`voicetext_james`, and `voicetext_kate`) do not import `VTDTTS_MakeInfo_ENG`.
The local binary and decompiler-report search found writer strings and the
export, but no DTT reader or code opening these suffixes. Thus the output's
producer-side meaning is mapped as a selected synthesis-timeline manifest;
its downstream consumer and user-facing purpose remain unverified.

A separate heap-backed path probe preserved and printed the complete second
argument before and after the API call. The exact output names were
`heap-makeinfo-path-20260926.bin.dtt` and
`heap-makeinfo-path-20260926.asc.dtt`, both byte-identical to the original
pair. The wrapper formats the supplied prefix directly with `%s.bin.dtt` and
`%s.asc.dtt`; it does not apply an additional basename transformation. The
earlier `...plead6` name and no-new-output follow-up used GDB literal strings
and are invalid for path conclusions because that setup corrupted the path
argument. See the new heap-path trace and runner in Stage 16.

The slot fallback observed in these paths is not proof that callers should pass
an invalid slot; it records the implementation's fallback branch. Source-level
signatures for these exports are unavailable in the public header. Address
references point to the local Paul M16 PE disassembly; see
[`vt_pau-objdump-disassembly.txt`](../../tools/revkit/work/reports/vt_pau-objdump-disassembly.txt)
and the wrapper pseudocode in
[`vt_pau-exported-api.c`](../../tools/revkit/work/reports/vt_pau-exported-api.c).

The two extended unload routines also have useful static contracts. For
`VT_UNLOADTTS_EXT_ENG`, an out-of-range speaker slot is replaced with slot 1;
it stops TTS and unloads that slot when loaded. If no speaker slots remain
loaded, it additionally attempts to unload all 1,024 user-dictionary slots
and tears down shared window and critical-section state. Runtime calls with
Paul slot 1 and James slot 4 loaded together confirm that unloading either
one preserves the other slot and its database-size query; unloading the last
slot clears both states. This was observed in both load orders. Invalid-slot
normalization and other voice-pair combinations remain static-only. For
`VT_UNLOAD_UserDict_EXT_ENG`, indexes outside 0–1,023 return `-2`; an
uninitialized subsystem or a dictionary referenced by a loaded speaker returns
`-3`; an idle dictionary is freed and cleared with result `1`; an empty slot
returns `-1`. The header declares a short return for the two-argument wrapper,
and the runtime probe captured low AX `1` for two successful loaded-idle
unloads. A forced-uninitialized call and a controlled reference-table match
both returned `-3`; these confirm the error branches without demonstrating
naturally uninitialized or active-synthesis states. See the Stage 16 in-use
comparison above.

The history-mode setter also has a clearer static boundary than the original
post-load runtime probe suggested. It writes its byte only while the global
load flag is clear, retains `1`, and maps all other inputs to `0`. The `1`
branch calls the shared no-op routine exported under the decimal-pronunciation
name. Other internal routines consult the flag to allocate/copy/free three
additional per-item arrays and to run extra per-item counter/update loops
(notably at `0x10019a70`, `0x10019b10`, `0x10019bb0`, `0x10012b5a`, and
`0x10024a96`). This supports describing it as a gate for internal
unit-selection history accounting; it does not establish a user-visible
quality or repeatability effect. The existing runtime calls happened after
load and therefore only demonstrated that this setter ignores calls in that
state.

A paired runtime probe set the flag to `0` or `1` at the entry to
`VT_LOADTTS_EXT_ENG`, before the model-loaded flag became set, then allowed the
same executable and Stage 16 input to synthesize file format 4. Both calls
returned `1`. Both WAVE files have SHA-256
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`, matching
the existing format-4 baseline. The GDB captures are
[`unit-history-0-api.log`](../../tools/revkit/work/stage16/unit-history-0-api.log)
and
[`unit-history-1-api.log`](../../tools/revkit/work/stage16/unit-history-1-api.log);
the captured WAVE files and runner are under Stage 16. This establishes no
output difference for this one text/settings case only; it does not show
whether the extra history accounting changes selection on other input or
affects repeatability across longer runs.

A Stage 21 trace followed the enabled loader path through `FUN_100125d0` and
`FUN_1001b240`. It requested
`../data-paul/M16/mc_idx_tbl/unit-gen.his`, `unit-num.his`, `unit-etc.his`,
and `unit-alp.his`; the mounted package contains none of those files. For all
four, the source handle was null, the fallback helper returned low AX 1, and
the count field plus two destination arrays were visible at the helper
boundary. The counts were 440,124, 24,508, 115,723, and 119, respectively,
matching the independently parsed `.idx` unit counts. Static instructions
clear each count-sized dword array when the `.his` source is absent; runtime
checks confirmed zero at the first and last entries in each array. This maps
the missing-file fallback for the local Paul package. It does not reveal the
format or purpose of a present `.his` file, nor whether the accumulated
counters affect selection or repeatability. See Stage 21
[`unit-history-load-api.log`](../../tools/revkit/work/stage21/unit-history-load-api.log),
[`trace-unit-history-load.gdb`](../../tools/revkit/work/stage21/trace-unit-history-load.gdb),
and [`run-unit-history-load.sh`](../../tools/revkit/work/stage21/run-unit-history-load.sh).

A synthetic read-only model overlay then exercised the present-file parser
without touching the vendor tree. Static reads and runtime sentinels establish
the accepted layout: a little-endian 32-bit count, then two little-endian
dwords per unit, first filling the `+0x14` array and then the `+0x18` array;
the count must equal the `.idx` unit count. Correct-sized files loaded all
four banks, and first/last sentinels landed in the expected array positions.
Four extra trailing bytes were ignored. A one-less header count and a file
missing its final dword each made `FUN_1001b240` return low AX 0, which made
the bank loop return -1. On both failures, `VT_LOADTTS_EXT_ENG` then continued
with a null state pointer and zero error word and faulted at
`0x10027d9c` (`mov [eax+0x4d08],edx`): SEH records a write access violation
with `eax=0` and target `0x4d08`. This maps the failure behavior for these two
malformed cases in this build. It does not identify what the two stored
per-unit values mean or prove behavior for every malformed size/count. Stage
21 captures are `unit-history-present-api.log`,
`unit-history-header-mismatch-api.log`, `unit-history-short-read-api.log`,
and `unit-history-trailing-data-api.log`; reproduce them with
[`run-unit-history-present.sh`](../../tools/revkit/work/stage21/run-unit-history-present.sh).

To test whether populated history values change synthesis, the valid synthetic
overlay's output for the existing Stage 5 utterance was compared with a fresh
load using the package's no-file fallback, with history mode enabled in both
processes. The two 23,650-byte WAVs are byte-identical (SHA-256
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`). This
shows no effect from these sparse first/last sentinels on this utterance and
settings. It does not show that history data is ignored generally: the test
does not exercise different utterances, repeated synthesis, or realistic
records, and the values' semantics remain unknown. The WAV artifacts and
fallback trace are `unit-history-present-output.wav`,
`unit-history-fallback-output.wav`, and
`unit-history-fallback-comparison-api.log` under Stage 21.

A wider synthetic-value matrix filled every unit record in every bank with
alternating pairs: even unit indexes held `(0,0)` and odd indexes held
`(0x7fffffff,0)`, `(0,0x7fffffff)`, or `(0x7fffffff,0x7fffffff)`. Each
pattern was compared against a fresh no-file fallback load on `Hello world.`,
`The quick brown fox jumps over the lazy dog.`, and `I saw 123 birds at
10:30.`. All nine patterned WAVs were byte-identical to their matching
fallback WAVs. Every load and text call returned success. Since each run
started in a fresh process, this tests whether these loaded values change the
first synthesis of the sampled utterances; it found no such audio effect.
That weighs against the fields being direct per-unit selection penalties on
these paths, but does not establish that they are ignored by other inputs,
later utterances in one process, or non-audio bookkeeping. The reproducible
runner, result matrix, twelve GDB logs, and paired WAVs are under Stage 21 as
`run-unit-history-pattern-matrix.sh`,
`unit-history-pattern-matrix-results.txt`, and the
`unit-history-pattern-*` artifacts.

The matrix runner's `repeat` mode extends each case within the same loaded
process. After the host's initial file synthesis, GDB calls the public
`VT_TextToFile_ENG` export three more times, reusing the captured text pointer,
format, speaker, pitch, speed, volume, pause, dictionary, and text-type
arguments while changing only the output filename. The fallback process and
each populated-pattern process therefore each perform four successive calls.
All 36 repeated API calls returned 1, and all 27 populated-pattern repeat WAVs
(three later calls × three patterns × three utterances) matched the WAV from
the same call position in the no-file fallback process. In this bounded
sequence, neither the loaded values nor same-process accumulation changed the
audio. The result still does not name the fields or exclude effects on other
text sequences or internal state that is not reflected in these WAVs. The
repeat trace is `trace-unit-history-repeat.gdb`; reproduce with
`run-unit-history-pattern-matrix.sh repeat`.

`VT_VerifyTTS_ENG` was called after slot 1 loaded, at the first file-synthesis
entry. Its wrapper is decompiled as `void`, but it returns immediately after
the internal validator, preserving that routine's status in EAX. These are
machine-register observations rather than a declared C result. Valid text
returns 1 on slot 1; negative slots and slots >=6 normalize to 1 before the
loaded-state check, while unloaded in-range slot 0 returns -4. Null text
returned low `AX=-2`; non-null empty text returned low `AX=-3`. The empty
case's upper EAX bits retained an unrelated prior value, so error results are
normalized from low AX rather than treated as full signed EAX values.

The third argument is a user-dictionary index. The helper at `0x10025fc0`
looks up the speaker's selected dictionary slot; values outside 0–1,023 use
the default slot, and a missing in-range pointer falls back to slot 0. Runtime
indexes -2, -1, 0, 1, 1,023, and 1,024 all returned 1 on the Stage 16 text.
Those ordinary-state calls did not test an active nonempty dictionary: the
supplied runtime's user-dictionary gate is closed, and the tested slots were
empty. A controlled comparison below loads a dictionary and debugger-forces
the gate for a separate export-call matrix.

The fourth argument is the text-type byte stored by the same setup helper
(`0x100286c0`) used by `VT_TextToBufferEX_ENG`. Its low byte is stored for
nonnegative values, negative values become zero, and types 4 or 6 reset the
four accompanying pitch/speed/volume/pause fields to -1. A complete runtime
sweep of types 0–255 with dictionary index -1 returned 1 for all 256 calls on
the Stage 16 text; a separate boundary matrix crossed dictionary indexes
-2/-1/0/1/1,023/1,024 with types -1/0/1/4/6/7/255 and all 42 calls returned 1.
Thus the arguments select dictionary/text-type processing state, but this
successful plain-text fixture does not show that they are semantically
equivalent or expose the resulting internal representation.

The one-character malformed input `<` returned low `AX=-5`, confirming runtime
reachability of the final-helper error branch; its upper EAX was `65531` in
this call. A 15-case follow-up found the same `-5` for `>`, `<>`, `<<`, and
`>>`, while `</>`, `<A>`, unterminated `<A`, bare `A>`, `<A></A>`, the
mismatched `<A></B>`, `A<B>`, and `A</B>` returned 1. These cases already
show that this result is not a balanced-tag check.

A further matrix directly called export `VT_VerifyTTS_ENG` at `0x1001ded0`
(a thin `void` wrapper over `FUN_100226f0`) for every string of length 1–4
over the four-byte alphabet `<`, `>`, `/`, `A` (340 strings) with slot 1 loaded,
dictionary index -1, and text type 0. It found 287 low-AX successes (`1`) and
53 `0xfffb` results (`-5`). Every tested string containing `A` succeeded;
53 of the 120 delimiter-only strings failed and the other 67 succeeded. This
exhausts that bounded alphabet and length, not the markup grammar: other
letters, whitespace, attributes, nesting depth, dictionary effects, and
other voice slots remain untested. The test records low AX, because the
export is decompiled `void` and high EAX bits are not a stable status source.
The 15-case capture is
[`verify-tts-markup-matrix-api.log`](../../tools/revkit/work/stage21/verify-tts-markup-matrix-api.log);
the full product capture is
[`verify-tts-delimiter-product-api.log`](../../tools/revkit/work/stage21/verify-tts-delimiter-product-api.log).
Reproduce them with `run-verify-tts-markup-matrix.sh` and
`run-verify-tts-delimiter-product.sh` in Stage 21.

A direct-export byte sweep then called `VT_VerifyTTS_ENG` with every first
byte 0–255 followed by NUL (slot 1, dictionary index -1, text type 0). NUL
returned low AX `-3`; its full EAX was stale (`21889021`). Of the 255
non-NUL byte strings, 160 returned `1` and 95 returned `-5` (low AX
`0xfffb`). The complete set returning `1` is:

```text
0x24-0x26, 0x2b, 0x2f-0x39, 0x3d, 0x40-0x5a, 0x5c,
0x61-0x7a, 0x80, 0x83, 0x89-0x8a, 0x8c, 0x99-0x9a, 0x9c,
0x9f, 0xa2-0xa5, 0xa7, 0xa9, 0xae, 0xb0-0xb3, 0xb5-0xb6,
0xb9, 0xbc-0xbe, 0xc0-0xff
```

Every other non-NUL value returned `-5`. These are one-byte NUL-terminated
strings, not a character-encoding result: the high-byte outcomes do not map
multi-byte encodings or establish the meaning of those byte values. The full
capture is
[`verify-tts-byte-domain-api.log`](../../tools/revkit/work/stage21/verify-tts-byte-domain-api.log);
reproduce with `run-verify-tts-byte-domain.sh` in Stage 21.

The active-dictionary matrix loaded the existing `hello,HH,P` row at index 27
(load AX 1 and non-null slot pointer). For each of six texts (`hello`,
`world`, `HELLO`, `hello hello`, `hello<`, and `<`), it called the export with
default index `-1`, index 27 while the gate was 0, and index 27 with the
speaker gate debugger-forced to 1. The first three groups returned low AX 1
for the first five inputs and -5 for `<`. With the gate on, index 28's empty
slot returned 1 for both `hello` and `world`. The dictionary unloaded with AX
1, and the gate was restored to 0. Thus, the populated dictionary did not
change the exported status on these inputs, whether the gate was on or off.
This status-only call does not establish whether the verifier uses dictionary
entries internally or changes an unobserved parse representation. The
capture is
[`verify-tts-dict-api.log`](../../tools/revkit/work/stage21/verify-tts-dict-api.log);
reproduce with `run-verify-tts-dict.sh` in Stage 21. Other dictionary rows,
inputs, voices, and dictionary-dependent internal effects remain open.

A larger 680-call repeat ran the full 340-string delimiter product first with
the populated dictionary gate off and then with it on. Each group returned
287 low-AX successes (`1`) and 53 `-5` results, matching every one of the
no-dictionary baseline's 340 individual return values. The gate-on group was
debugger-forced; the capture is
[`verify-tts-delimiter-dict-product-api.log`](../../tools/revkit/work/stage21/verify-tts-delimiter-dict-product-api.log).
This establishes status invariance only for this bounded text product.

Static pseudocode explains the setup difference: `VT_VerifyTTS_ENG` calls
`FUN_10025fc0`, which selects the dictionary pointer and stores it in the
new context at `+0x1312c0`; after creating the nested language state, the
code writes zero at language-state `+0x39e8` when the speaker gate is zero,
and otherwise copies the selected dictionary pointer there. The text parser
then receives the context. This proves that the two gate states prepare
different dictionary state before parsing, but does not show whether these
particular texts trigger a dictionary match or how such a match changes an
unreturned representation. This is based on Ghidra pseudocode in
[`stage23-state-array-writers.c`](../../tools/revkit/work/reports/stage23-state-array-writers.c),
not a runtime read of the temporary context fields.
Captures and reproduction scripts are [`verify-tts-api.log`](../../tools/revkit/work/stage16/verify-tts-api.log),
[`verify-tts-matrix-api.log`](../../tools/revkit/work/stage16/verify-tts-matrix-api.log),
[`verify-tts-types-api.log`](../../tools/revkit/work/stage16/verify-tts-types-api.log),
[`verify-tts-text-edges-api.log`](../../tools/revkit/work/stage16/verify-tts-text-edges-api.log),
[`run-verify-tts.sh`](../../tools/revkit/work/stage16/run-verify-tts.sh),
[`run-verify-tts-matrix.sh`](../../tools/revkit/work/stage16/run-verify-tts-matrix.sh),
[`run-verify-tts-types.sh`](../../tools/revkit/work/stage16/run-verify-tts-types.sh),
and [`run-verify-tts-text-edges.sh`](../../tools/revkit/work/stage16/run-verify-tts-text-edges.sh).

## SyncInfo allocator transient-failure retry

`VT_AllocSyncInfo_New_ENG` has three allocation classes: a 56-byte header,
the 600-row array of 36-byte rows, and one 520-byte block for each row's 65
eight-byte nested entries. Ghidra pseudocode for `FUN_1001d9c0` calls the
internal allocation routine and, on a null result, sleeps for 10 ms and
repeats the same request until it succeeds. This makes the allocator's direct
null-return checks in its caller unreachable for ordinary transient
allocation failures, though it can still loop indefinitely if memory never
becomes available.

A PE32 client calls the exported allocator and freer. The GDB probe forces the
first nested 520-byte allocator result to null once at return address
`0x100263be`, the callsite following the `0x100263b9` allocation call in the
allocator export. The retry returns nonnull; the export then completes all
600 nested allocation calls. The injected run's post-call checkpoint confirms
a valid 600×65 object with nonnull row/first/last nested pointers and distinct
first/last nested allocations. A separate direct Wine control reports the
same shape. The probe covers one
transient failure at one nested-allocation site. It does not test persistent
out-of-memory behavior, allocator-wide fault injection, or cleanup after a
forced API-level null return. The capture and replay files are
[`syncinfo-allocation-retry-v7-api.log`](../../tools/revkit/work/stage21/syncinfo-allocation-retry-v7-api.log),
[`probe-syncinfo-allocation-retry.c`](../../tools/revkit/work/stage21/probe-syncinfo-allocation-retry.c),
and [`run-syncinfo-allocation-retry-v7.sh`](../../tools/revkit/work/stage21/run-syncinfo-allocation-retry-v7.sh).

## Compatibility boundary

The incoming output length was tested only for short format-0 output and for
the 1 MiB long-stream buffer; other formats and error paths may differ. The
state transitions, cancel, one mixed file/buffer sequence, one VTML
substitution, and null-sink playback are established only for these cases.
Nonzero thread IDs, other markup handlers, malformed markup, actual audible
playback, request-101 state on other natural-completion paths, DBCS and other
UTF-8 caller-span cases, exact notification scheduling/repetition rules,
abnormal/unknown errors, cross-voice or cross-version compatibility, and
whole-synthesis parity remain unestablished. The direct-host Stage 21 probes
establish request 101's natural-completion value for the seven original texts
and two UTF-8 follow-ups. Callback positions are byte offsets for the ASCII,
CP1252, and two UTF-8 inputs captured so far; DBCS and other UTF-8 shapes
remain untested.

## `VT_SetParenthesisCharNumber_ENG` consumer and threshold

The prior boundary sweep established storage only: on the loaded Paul's
process-global configuration, negative inputs clamp to zero and the tested
nonnegative values `0`, `1`, and `INT_MAX` are retained. The consumer trace
now shows how the field enters the parser. At the call to `FUN_1003e470`, the
global value is copied to a temporary parser context at `context+8`. A
hardware access watchpoint then observed the write and four subsequent reads
during `VT_TextToFile_ENG` on
`I like (tea). I like (green tea). I like (freshly brewed tea). I enjoy (very
fresh tea) today.`. The read instruction is in `FUN_100544f0`, whose Ghidra
pseudocode handles parser items associated with `(` and `[` markers. For the
four parenthesis spans, the compared helper results were 3, 9, 18, and 14.
The helper `FUN_10062ea0` scans bytes until the first `)` or `]` and returns
that byte offset; the values match the interior byte lengths of these ASCII
examples. This is a byte distance, not a demonstrated Unicode character
count.

The pseudocode branches around the special record path when the setting is
zero or greater than the helper result. This makes a nonzero setting a
minimum byte-span threshold in the observed parser route: the marked-record
branch is eligible when the nonzero byte offset is at least the setting. At
the eligible branch, the machine code writes `1` to a DWORD in a parser
record slot indexed from local parser state; the field's schema name and
downstream interpretation are not recovered. A second watchpoint run with
`We saw (tea) today.` followed that exact DWORD after the write: it was read
at `0x1000d208` in `FUN_1000d190`, read again in the later `FUN_1000ea20`
path, copied as part of a record in `FUN_10016c90`, and eventually zeroed
during context cleanup at `FUN_1003e210`. In `FUN_1000d190`, the record's
leading DWORD feeds a switch with explicit cases 2, 3, 4, 5, 11, and 12; the
observed value 1 takes its default arm. This establishes that the marked slot
is consumed and carried into later processing, but not what the marker means
to synthesis. Direct format-4 audio
comparisons corroborate boundaries at byte spans 1, 3, and 4:
for `(a)`, 0 matches the baseline and 1 changes the WAV; for `(tea)`, 1, 2,
and 3 are byte-identical while 0 and 4–6 match the baseline; for `(book)`, 1
and 4 produce byte-identical changed WAVs; and `(green tea)` also produces
identical output at 1 and 4. These audio results show an effect on this
selected voice and format, but do not identify what the parser's marked
record means to the synthesis stages.

The watchpoint capture is
[`parenthesis-context-consumer-api.log`](../../tools/revkit/work/stage21/parenthesis-context-consumer-api.log),
reproduced by `run-parenthesis-context-consumer.sh`. The controlled WAV
comparisons and per-mode runtime logs are under Stage 21 and are reproduced
with `run-parenthesis-number-effect.sh`. The consumer pseudocode is
[`stage25-parenthesis-parser-consumer.c`](../../tools/revkit/work/reports/stage25-parenthesis-parser-consumer.c).
The byte-scanning helper pseudocode is
[`stage25-parenthesis-length-helper.c`](../../tools/revkit/work/reports/stage25-parenthesis-length-helper.c).
The reader pseudocode is
[`stage25-parenthesis-record-consumers.c`](../../tools/revkit/work/reports/stage25-parenthesis-record-consumers.c),
and the record follow-up trace is
[`parenthesis-record-consumer-api.log`](../../tools/revkit/work/stage21/parenthesis-record-consumer-api.log).
This establishes the byte-threshold mechanic only for the sampled parser
route; square-bracket effects, other voices/formats, multibyte-span behavior,
and downstream record interpretation remain open.
