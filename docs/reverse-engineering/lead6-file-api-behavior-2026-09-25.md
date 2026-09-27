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
| Unsupported selector 6 or out-of-range selector 10 | `-1` (`VT_FILE_API_ERROR_INVALID_FORMAT`) | [`file-format-6.log`](../../tools/revkit/work/stage16/file-format-6.log), [`file-format-10.log`](../../tools/revkit/work/stage16/file-format-10.log) |

The observed returns agree with the `VT_FILE_API_*` constants in
[`vt_eng.h`](../../include/vt_eng.h). The error probes establish these cases
only; they do not cover database-unloaded behavior, inaccessible non-null
paths, thread creation failures, or unknown errors.

Buffer format 4, null text, empty text, and null output buffer returned
`-1`, `-3`, `-4`, and `-5`, respectively, matching
`VT_BUFFER_API_ERROR_INVALID_FORMAT`, `NULL_TEXT`, `EMPTY_TEXT`, and
`NULL_BUFFER`. See [`buffer-errors.log`](../../tools/revkit/work/stage16/buffer-errors.log).
The database-unloaded (`-6`), busy-thread (`-7`), abnormal-condition (`-8`),
and unknown (`-9`) codes are declared in the header. The database-unloaded and
busy-thread results are exercised below; abnormal-condition and unknown
returns were not reached.

## User-dictionary API spot checks

`VT_GetUserDictLimit_ENG` returned `30`, `10`, `50`, `65`, and `65` for
selectors 0–4, and `-1` for selector 5. The public header declares short
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
independent transcription, and the ordinary run's gate remains unexplained.

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
frees an idle loaded dictionary, and reports an empty slot. The loaded-idle
unload path is runtime-confirmed; in-use and uninitialized errors remain
static-only. Runtime logs, WAV hashes, GDB traces, and reproducible runners are
in Stage 16, including `userdict-heap-lifecycle-api.log`,
`userdict-hello-world-state-api.log`, and `userdict-hello-world-api.log`.

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

This only shows per-field byte changes for `Hello world.` and the listed
limits. It does not establish pause effects for text with internal pauses,
their audible duration, or output differences at intermediate values. The
trace, log, and WAVE captures are under Stage 16.

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
audible output, acoustic quality, completion notification behavior, a physical
audio device, or the thread lifecycle. The stop call was attempted, but the
inferior exited normally while GDB was evaluating it, so no stop-call return
is claimed. See `play-waveout.log`, `alsa-null.conf`, and
`trace-play-waveout.gdb` under Stage 16.

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

## Additional export queries

The Stage 16 follow-up calls four export-only query helpers after the Paul M16
model has loaded. `VT_GetSpeakerName_ENG` returned `Kate`, `Paul`, `em001`,
`Julie`, `James`, and `Ashley` for slots 0–5; slot `-1` fell back to `Paul`.
`VT_SpeakersInfo_ENG` returned `6` for each slot and copied lowercase IDs and
the DLL's embedded `d:/eng/db/.../pcm/` path strings. These strings are
metadata returned by the DLL; this trace did not open those paths.
`VT_GetDefVersion_ENG` returned `Paul-M16-FileIO`. `VT_GetDBSize_ENG` returned
`1` and `508121688` for loaded slot 1; slots 0 and 2–5 returned `-1`, leaving
the zero-initialized output untouched. The byte count is recorded as the API
result without inferring which files or allocation it measures. See
[`export-queries-api.log`](../../tools/revkit/work/stage16/export-queries-api.log)
and its replay script.

## Additional helper exports and data exports

The Stage 16 helper probe ran after the Paul model loaded. `VT_GetPathKey_ENG`
returned `SOFTWARE\VW\VT\<speaker>\M16` for all six speaker slots. The
absolute and relative default user-dictionary name helpers returned
`../data-common/userdict/userdict_eng.csv` and
`data-common/userdict/userdict_eng.csv`, respectively. These are the observed
strings in the local DLL; no dictionary was loaded by this query.

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
general malformed-quote or multiline-record handling. Non-ASCII encoding
behavior remains open. A separate edge probe observed that doubled quotes in
`"a""b",c` and a quote in unquoted `a"b,c` both returned `a"b` and `c`;
`a,` omitted its empty final field; and `a,\r\nb` removed the CRLF at the
beginning of the second field. In `a,b<LF>c,d`, the LF remained inside the
second field (`b<LF>c`) and the parser returned three fields. The exact row
`a,"b,c",d` returned the same three fields with parse flags 0, 1, and 2.
These samples do not establish other quote, line-position, or flag behavior.
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
The underwrite stayed inside the probe allocation. The wrapper is decompiled
as `void`, so captured AX values are raw register state, not declared C
returns. A null field-array pointer with field count zero returned raw AX 1
and produced an empty string. Null pointers in other argument combinations
and capacities above 64 remain open. The exact trace is
[`csv-capacity-api.log`](../../tools/revkit/work/stage16/csv-capacity-api.log)
and runner is
[`run-csv-capacity.sh`](../../tools/revkit/work/stage16/run-csv-capacity.sh).
The parser-edge trace and runner are
[`csv-parser-edges-api.log`](../../tools/revkit/work/stage16/csv-parser-edges-api.log)
and [`run-csv-parser-edges.sh`](../../tools/revkit/work/stage16/run-csv-parser-edges.sh).

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

Cross-references in internal pseudocode narrow several structural roles. The
producer `FUN_1002c530` uses header field 3 as a row write cursor, advances it,
and wraps it at 600. The consumer `FUN_1001e0c0` walks pending rows relative to
that cursor; row offset 0 controls its nested-entry loop. For each nested
entry, it adds the first dword to an accumulator, while the following 16-bit
value selects a different formatting branch when it equals `0x28`. Row offset
8 contributes to another accumulator, and offsets 12/16 bound a text slice
copied by the consumer. These observations identify control/data roles, not
the semantic units of the accumulated values or meanings of the remaining row
fields. Evidence is in the Ghidra pseudocode reports
[`vt_pau-core-pseudocode.c`](../../tools/revkit/work/reports/vt_pau-core-pseudocode.c)
and [`stage8-synthesis-functions.txt`](../../tools/revkit/work/reports/stage8-synthesis-functions.txt).

`VT_CopySyncInfo_New_ENG` copies header fields 3–13, the row 16-bit value and
seven dwords, and both nested-entry values, while retaining the destination's
row-array and per-row nested pointers. The runtime call used the native 600-row,
65-entry dimensions. The probe compared all 11 header fields, 4,800 row values
(one 16-bit value and seven dwords per row), and 78,000 nested values (one
dword and one 16-bit value per entry); every comparison matched. All 600
destination nested pointers and the row-array pointer remained distinct.
The probe restored both objects' allocator dimensions before freeing them.
Null/self-copy behavior, mismatched dimensions, and semantic names for these
values remain open. Full trace and repeatable runner are
[`syncinfo-fields-api.log`](../../tools/revkit/work/stage16/syncinfo-fields-api.log)
and [`run-syncinfo-fields.sh`](../../tools/revkit/work/stage16/run-syncinfo-fields.sh).

On the literal input `example`, `VT_CheckUserDict_SourceNorm_ENG` completed
and left the input unchanged; its wrapper writes the normalized value to local
scratch that is not returned. `VT_CheckUserDict_TargetNorm_ENG` returned `1`
and `VT_CheckUserDict_TargetPhon_ENG` returned `-9`. These are single-input
observations; the labels and codes are not assigned broader validation
semantics. `VT_SetEmphasisFactor_ENG` clamped slot 1 values `200` and `-200`
to `95` and `-95`. `VT_SetTextTypeForHighlight_ENG(7)` stored `1`;
`VT_SetParenthesisCharNumber_ENG(-1)` and
`VT_SetEnglishReadingRule_KOR(-1)` each left their state field at `0`.
`VT_SetSoundCardID_ENG` returned the prior `0xffffffff` value after a
temporary write. The probe restored every field it changed before normal
execution continued. `VT_SetUnitSelectHistoryMode_ENG` calls with `1` and `0`
left its global state at `0` while the loaded flag was `1`, matching the
decompiled guard that ignores calls after model loading.

`VT_TextToBufferEX_ENG` returned `-1` for selectors `-1` and `3` when all
remaining arguments were null or zero. For selectors 0–2, the Stage 16 probe
passed the same sample text used by the ordinary buffer calls, a 60,000-byte
capacity within a 60,000-byte guarded allocation, flag 0, thread 0, speaker 1,
null optional pointer arguments, and `-1` for the four intervening scalar
arguments. The first call returned `0` and wrote 23,606 bytes for selector 0;
selectors 1 and 2 returned `0` and wrote 11,803 bytes each. A follow-up call
with flag 1 returned `1` and length 0 for each selector. All four leading and
trailing guard bytes retained their `0xa5` sentinel values. The process exited
normally.

Each captured byte stream has the same SHA-256 as the matching ordinary
`VT_TextToBuffer_ENG` format call: selector 0 matches format 0
(`e232e245…f61582`), selector 1 matches format 1 (`9880164f…d5dfc`), and
selector 2 matches format 2 (`424afda0…e53ab`). This establishes output
equivalence for that input and argument tuple. It does not identify the
meanings of the extra pointer/scalar arguments or establish other input,
capacity, error, repeated-poll, chunking, or cancellation behavior. Static
wrapper pseudocode rejects selectors outside 0–2 before dispatching valid
selectors to the three format handlers at `0x100200c0`, `0x10020460`, and
`0x10020930`.

The three run logs and captured streams are in [Stage 16](../../tools/revkit/work/stage16/README.md):
[`buffer-ex-0.log`](../../tools/revkit/work/stage16/buffer-ex-0.log),
[`buffer-ex-1.log`](../../tools/revkit/work/stage16/buffer-ex-1.log), and
[`buffer-ex-2.log`](../../tools/revkit/work/stage16/buffer-ex-2.log). The
probe is reproduced by
[`run-buffer-ex.sh`](../../tools/revkit/work/stage16/run-buffer-ex.sh).

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
a simple framing/byte-offset transform observed in the instructions; the
code path does not establish a cryptographic checksum or signature check.
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
lookup into this routine. This shows that the token participates in a custom
linked conversion/consistency check; it does not establish a cryptographic
algorithm, an authentication guarantee, or the business meaning of its
outputs. The `VW_VTAPI` token remains an opaque product/record discriminator.
The local `data-*` inventory contains only this one file with the
`vw_verify`/`VW_VTAPI` record markers, so there is no second accepted local
record with which to validate whether these positions and associations
generalize.

The `hostid` label above names the field in the archived record. It does not
show that the value was generated for, or issued to, the workstation used for
this probe; the record came with the archived VoiceText package.

`FUN_10029170` passes the retained `License` value, the second lookup's value,
and the presence flag (as either a null pointer or a built-in string) to
`FUN_10014dd0`. The 96-character value is split into two 48-character halves;
the disassembly passes them through separate fixed-format conversion helpers,
compares the converted strings against the other supplied values, and returns
two numeric outputs on success. The caller bounds one output at
99,999,999 and another at 1,024 (zero selects the default 1,024). The date
validator at `FUN_10029510` repeats that conversion path, formats the local
system date as `YYYYMMDD`, and rejects when today is later than its converted
value. This establishes a linked two-value validation pipeline and its
numeric bounds, but the numeric outputs are not mapped back to specific
colon-separated fields. The positional `0` matching XML `expdate` is a sample
correlation; it does not prove that the date conversion consumes that field.

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
load, after load, and after unloading speaker 1. Its purpose remains unknown.
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
| `VT_TextToBufferEX_ENG` | Runtime selectors 0–2 match the corresponding ordinary buffer formats on the Stage 16 text and argument tuple; each finishes with one flag-1 poll. Static wrapper rejects selectors outside 0–2 with `-1` and dispatches to handlers at `0x100200c0`, `0x10020460`, and `0x10020930`. | Other pointer/scalar arguments, invalid inputs, multi-poll/chunk/cancel behavior, and extra-argument meanings. |
| `VT_TextToPreprocessInfoFile_ENG` | The wrapper at `0x1001dd60` delegates to `FUN_1001ef20` (`0x1001ef20`) and is decompiled `void`. Flag 0 returns before pointer checks. A fresh-process matrix used speaker 1, pitch/speed/volume/pause `-1`, dictionary 0, and text type 0; byte flags 1–10, 11, 254, and 255 completed with raw low AX 1 (not a declared return). Correct heap-backed calls establish that flags 1–5, 7–11, 254, and 255 use `param2` literally; flag 4 appends `.0`–`.3`, while flag 6 ignores `param2` and writes fixed `test.pcm`. The earlier malformed names resulted from GDB setup corrupting the path before API entry. Flag 3 emits one dot per source byte, a count, and decimal engine phone/control bytes that match the recovered phone codebook; flag 5 renders the same stream as phone labels. Flag 7 writes source text and per-word surfaces, structural boundary codes, phone labels, metadata-bit letters, and inclusive ASCII spans. Its leading period string is selected among eight fixed lengths by a `GetTickCount`-seeded generator; identical-input runs vary, and the reason for including the filler is unknown. Flag 10 spaces tested terminal punctuation; flag 6 wrote 23,606 bytes of PCM. The boundary-class names, bit-label meanings, broader text/punctuation behavior, malformed nonempty input, invalid scalar combinations, and individual runtime calls for flags 12–253 remain open. See the isolated observations and output-record interpretation below, the heap-backed captures and runners in Stage 16, and the coverage-map entry. |
| `VT_TextToPcmBuffer_ProgressBar_ENG` | Runtime after model load with the Stage 16 text, thread 0, speaker 1, option values -1, and null window/message: 23,606 bytes were written, matching ordinary buffer format 0; guard bytes remained intact. Raw EAX after the decompiled `void` wrapper was 1, which is not treated as a declared return value. Static path uses speaker-indexed state and calls `PostMessageW` in the chunk-progress path. See [`buffer-progress.log`](../../tools/revkit/work/stage16/buffer-progress.log) and [`buffer-progress-start.bin`](../../tools/revkit/work/stage16/buffer-progress-start.bin), reproduced by [`run-buffer-progress.sh`](../../tools/revkit/work/stage16/run-buffer-progress.sh). | Long-text progress notifications, callback delivery, error cases, and other arguments. |
| `VT_TextToLipSyncLog_ENG` | Runtime after model load with Stage 16 text and second string `vtspeak-lead6`: raw EAX was 1 and an ASCII report was produced. The captured report lists phone IDs/labels and integer length values, source word-index/text spans, and total word/phone lengths of 11,803 for this sample. Empty second string also returns raw EAX 1; speaker -1 falls back and returns 1; null text produced AX -3; empty text produced AX -4 with raw EAX `0x003efffc`. Static code routes the second string through `FUN_1001df10`; static data at `0x1007d538` contains `length-sync-%s-%s.txt`, and the runtime produced a matching `length-sync-*` file. This supports, but does not prove, using the string as a naming/tag input. A path-like second string caused an access violation. The decompiled wrapper is `void`; register values are not declared C returns. | Accepted second-string grammar, exact filename components, units for numeric lengths, path-like failure, and option/text variation. Evidence: [`lipsync-prefix.log`](../../tools/revkit/work/stage16/lipsync-prefix.log), [`lipsync-prefix-output.txt`](../../tools/revkit/work/stage16/lipsync-prefix-output.txt), and prior [`lipsync-output.log`](../../tools/revkit/work/stage16/lipsync-output.log). Reproduce with [`run-lipsync-prefix.sh`](../../tools/revkit/work/stage16/run-lipsync-prefix.sh). |
| `VTDTTS_MakeInfo_ENG` | Runtime with Stage 16 text, speaker 1, remaining scalar args -1, and a heap-backed second argument completed with raw EAX `1`. It emitted the literal prefix plus `.bin.dtt` and `.asc.dtt`; the heap-path pair byte-matches the original captures. The binary magic is `VTDTTS BINARY\0`; bytes `03 04` and four NUL-terminated bank names (`merged-gen`, `merged-num`, `merged-etc`, `merged-alp`) follow. The ASCII header spells these as `3`, `4`, and the same names; `4` matches the named-bank count, while `3` remains unlabeled. Captured records decode identically in ASCII and binary. Phone rows cross-check to `unit-*.idx` plus their `.dat`/`.upm` spans. Mode 0 selects a whole decoded unit and combined UPM span; mode 1 selects first-side spans; mode 2 selects second-side spans and adds `Shift Size`. The writer computes that field as the 16-bit value at the selected-unit descriptor's `+0x0c` minus the timeline row's `-0x05` word. The subtraction is directly visible in disassembly; calling its result a splice/crop location would go beyond the evidence. `TypeFlag=2` is the detailed phone/unit record. `TypeFlag=1` is observed between OW1 and W for `Hello, world.` (`Size=3200`) and `Hello. World.` / `Hello... World.` (`Size=14800`), but absent for `Hello world.`. The synthesis loop zero-fills `Size*2` bytes, so Size is a count of PCM16 samples and type 1 is an inserted silence interval in these cases. Repeating the comma case with pause 120 retained the same size. Controlled text captures establish `File Index` values 0=`merged-gen`, 1=`merged-num`, and 2=`merged-etc`; index 3 has not been observed. Rate captures establish `Pitch_rate` = supplied pitch, `Volume_rate` = supplied volume, and `Duration_rate` = `((speed >> 1)+10000)/speed` with integer division (speed 120 yields 83); pause 120 left the fields unchanged. The format is therefore a selected synthesis-timeline manifest of model unit spans, timing/pitch spans, rate metadata, and punctuation-related silent intervals. The four checked-in host executables have no import for this export, and the local binary/report search found no DTT consumer; its intended external client remains unidentified. The record layout is: `u8` type flag; NUL-terminated phone string; `u8` file index; little-endian `u32` PCM position; `u16` PCM size and coded size; `u8` mode; two `u16` pitch endpoints; `u32` PM position; `u16` PM size; then three `u16` rates, with an extra `u16` shift field for mode 2. Static code also checks loaded/nonempty text, applies an invalid-speaker fallback to slot 1, writes both files, and frees the context. The export wrapper is decompiled `void`, so EAX is a raw register observation. See the Stage 16 rate/text/punctuation captures and disassembly around `0x1002ca10`. | Header byte 3; physical meaning of mode-2 shift arithmetic; punctuation-to-silence duration rules beyond the three cases; `merged-alp` row selection; broader option/text/error behavior; identity of any external DTT consumer. |
| `VT_VerifyTTS_ENG` | The internal validator applies the speaker-slot fallback and loaded/text checks before calling internal parsing and verification helpers. | Other argument combinations and result meanings. |

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
| 12–253 | 4 (static default branch; not individually runtime-probed) | Not individually probed. |
| 254 | 4 | One empty file. |
| 255 | 4 | One empty file. |

The mode mapping is from the dispatch at `0x1001eff2`–`0x1001f024`: the
original byte is stored at context `+0x21`; flags 1–9 index the table at
`0x1001f1cc`, while the default sets mode 4. The bounds check routes byte
values 10 and above to that default; 11, 254, and 255 were checked at runtime.
Thus the internal dispatch is mapped, while public names for these modes are
not established.

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
  at record offset `+0xa09`: `[` prints as `0`, `\` as `1`, `]` as `2`, and
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
  33, 151, 106, 58`. This explains the prefix's variable construction; why
  the report deliberately gets this tick-count-selected filler remains open.
* Flag 10 re-emits source text with spaces before terminal punctuation in the
  tested ASCII cases. This behavior is directly shown for `.`, `!`, and a
  comma followed by a terminal mark, but broader punctuation and markup rules
  are not covered.

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
`merged-num` with index 1. `A.B.C.` remained on `merged-gen`; no tested row
used index 3 (`merged-alp`). Thus the observed index is zero-based bank order,
with only the first three mappings demonstrated. `03 04` is the observed
binary header after the magic; `04` matches the four names that follow, but
the meaning of `03` is still not established.

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
timeline include a type-1 interval, but do not establish the full punctuation
or pause-setting rule.

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
the exact arithmetic in the disassembly, but the two fields' physical units
and the resulting shift's role are not yet identified.

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
and tears down shared window and critical-section state. For
`VT_UNLOAD_UserDict_EXT_ENG`, indexes outside 0–1,023 return `-2`; an
uninitialized subsystem or a dictionary referenced by a loaded speaker returns
`-3`; an idle dictionary is freed and cleared with result `1`; an empty slot
returns `-1`. The header declares a short return for the two-argument wrapper,
and the runtime probe captured low AX `1` for two successful loaded-idle
unloads. In-use and uninitialized errors remain static-only.

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

The same runtime setup sampled `VT_VerifyTTS_ENG` at the first file-synthesis
entry, after speaker 1 was loaded. With the current synthesis text and
trailing arguments `0, 0`, slot 1 returned raw `EAX=1`; slot `-1` also
returned `EAX=1`, matching the static fallback to slot 1; unloaded slot 0
returned `AX=-4` (unsigned `65532`). A null text pointer on slot 1 returned
`AX=-2` (unsigned `65534`). The exported wrapper's decompiler signature is
`void`, but its machine code calls the internal validator and returns without
changing EAX, so these are observed register results rather than a documented
C return contract. Empty text and other argument values remain unprobed. See
[`verify-tts-api.log`](../../tools/revkit/work/stage16/verify-tts-api.log)
and [`run-verify-tts.sh`](../../tools/revkit/work/stage16/run-verify-tts.sh).

## Compatibility boundary

The incoming output length was tested only for short format-0 output and for
the 1 MiB long-stream buffer; other formats and error paths may differ. The
state transitions, cancel, one mixed file/buffer sequence, one VTML
substitution, and null-sink playback are established only for these cases.
Nonzero thread IDs, other markup handlers, malformed markup, actual audible
playback, completion notifications, abnormal/unknown errors, cross-voice or
cross-version compatibility, and whole-synthesis parity remain unestablished.
