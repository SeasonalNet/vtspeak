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

The user-dictionary exports were called in the loaded process using the
header-declared wrappers. `VT_GetUserDictLimit_ENG` returned `30`, `10`, `50`,
`65`, and `65` for selectors 0–4, and `-1` for selector 5. `VT_LOAD_UserDict_ENG`
exposed low AX `-1` for indexes `-1` and `1024`; loading the current WAVE
output path as a dictionary exposed low AX `-3`. `VT_UNLOAD_UserDict_ENG`
exposed low AX `-2` for index `-1` and `-1` for empty slot 0. These wrappers
are declared `void`, so the values are raw register observations that happen
to match selected header error values, not declared C returns. Additional
probes tried plain `source,target,P` and `source,target,A` rows using targets
accepted by the isolated target helpers. Both load attempts exposed low AX
`-3`. Additional runtime calls with `hello,HH,P`, its fully quoted form, and
the four-field form `hello,HH,P,` also exposed low AX `-3`. The earlier quoted
`vtspeakprobe` candidate call was interrupted by nested parser breakpoints and
captured no result. A control synthesis and synthesis with the
rejected dictionary index both produced the same WAVE
SHA-256 (`c748068e46c635381303fa4beccbc76661e17adb7b97c928f75024550f94588f`);
this does not measure dictionary behavior because no candidate loaded. The
probe did not load a valid user dictionary, exercise an index-busy or in-use
case, or establish the effects of a loaded dictionary on synthesis. Exact
calls are recorded in
[`userdict-errors-api.log`](../../tools/revkit/work/stage16/userdict-errors-api.log)
and [`trace-userdict-errors.gdb`](../../tools/revkit/work/stage16/trace-userdict-errors.gdb).
Candidate traces and CSV fixtures are `userdict-valid-api.log`,
`userdict-alpha-api.log`, `userdict-debug-api.log`,
`trace-userdict-valid.gdb`, `trace-userdict-alpha.gdb`,
`trace-userdict-debug.gdb`, `userdict-candidates-api.log`,
`trace-userdict-candidates.gdb`, and the userdict CSV files under Stage 16.

Static disassembly of `FUN_1005e330` (`0x1005e330`) establishes part of the
loader's row gate: its CSV check permits 3 or 4 fields; it reads fields 0, 1,
and 2; field 2 must be exactly one character and map to `A` or `P`; and the
source normalizer must return nonzero. The `A` path applies target normalization
and rejects a negative result. The `P` path applies target-phon validation,
then calls the phoneme converter and rejects a zero result. A direct private
helper probe returned source-normalized length 12 for `vtspeakprobe`, target
normalization result 1 for `hello`, target-phon result 1 for `HH`, and
conversion result 1 with a one-character quote-string output for `HH`. Those
accepted helper results plus the failing candidate loads show these checks are
not the complete load acceptance contract. The public source-normalizer
wrapper discards its scratch output. See
[`userdict-normalizers-api.log`](../../tools/revkit/work/stage16/userdict-normalizers-api.log),
[`userdict-candidates-api.log`](../../tools/revkit/work/stage16/userdict-candidates-api.log),
and their Stage 16 runners.

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
general malformed-quote or multiline-record handling. Encoding and
output-capacity limits remain open.

Two synchronization-info objects were allocated and freed. Calling
`VT_InitSyncInfo_New_ENG` on one API-allocated object set observed fields 1 and
2 to `600` and `65`, and fields 4, 6, and 13 to `-1`. Before copying, scalar
fields 4 and 6 in the source were set to `0x13579bdf` and `0x2468ace0`; both
values appeared in the target after `VT_CopySyncInfo_New_ENG`. This confirms
selected initialization and copy fields only, not the full structure
contract.

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

For `VT_GetLicenseInfo_ENG`, null-path size queries with selectors 0 and 12,
null output, and capacity `-1` each returned `4`. For selector 0, capacities
0 and 3 returned `-5`; capacity 4 with null output returned `-4`. This matches
the decompiled required-size and output-pointer checks. No license contents
were requested or captured. Static selector dispatch shows selector 3 copies
the complete `<vw_verify .../>` attribute text; selectors 4–11 and 13–15
look up `os`, `lang`, `speaker`, `version`, `dbaccess`, `sampling`, `app`,
`wavsave`, `bgaudio`, `dbsize`, and `realtime`, respectively; selector 12
looks up `savetime` and returns a four-byte value (or `-1` when absent).
Selectors 1 and 2 return two other parser-produced strings whose field names
are not identified in the recovered code. Selector 0 writes a four-byte
parser-produced value. This maps selector routing, not the meaning or actual
values of every license field; runtime output for those selectors and
supplied-license inputs remains untested.

### License record parsing: static trace

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

`FUN_10029680` copies the resulting text, splits it at semicolons, and parses
each nonempty entry into an internal linked list using token extraction around
whitespace, colon, and equals delimiters. The validators then access entries
by position. The recovered numeric checks reject parsed values above
`99,999,999`; another field is defaulted to `1024` when zero and rejected
above `1024`. These are observed bounds, but their field names and licensing
meanings are not established. A separate validator calls `GetSystemTime`,
formats the current date as `YYYYMMDD`, and compares it with a parsed numeric
field. That is consistent with an expiration-date check, but the exact field
mapping and comparison rule are still unresolved. Thus we know the framing
transform and part of the post-read parser, not the syntax or semantics of
the complete outer license envelope.

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
`0x10028f20`, `0x10029680`, `0x10029170`, and `0x10029510`.

In the supplied sample, `app` is not queried in that
`VT_CheckLicense_ENG` sequence, while the code queries `lang`. The attribute
helper initializes its result to true and keeps that value when its extractor
returns zero, so the decompiled path does not establish that a missing `lang`
attribute is rejected. This does not establish whether `app` or other fields
are consumed by a different code path. The branch initializes its result to
zero after record setup and leaves it zero if either tag-boundary search
fails. Static control flow therefore shows what it does after the record
helper returns success; it does not prove malformed or absent envelopes are
accepted by the entire API. No runtime call against the supplied file was
made.

The five data exports read as `VT_gHeapStartAddress_ENG = 0` and version
components `(first, second, third, fourth) = (3, 11, 7, 1)` after model load.
The tuple matches both the DLL's PE FileVersion and ProductVersion strings
`3.11.7.1`; this supports interpreting the four exports as version
components. Separate traces read the heap-start export as zero before model
load, after load, and after unloading speaker 1. Its purpose remains unknown.
The NULL-path `VT_GetLicenseComment_ENG` query returned
322, but its output is deliberately redacted because it contains
host/license-specific data. The direct null-path `VT_CheckLicense_ENG` call
returned `0`; `VT_INIT_ENG` returned `-1`. Exact trace lines are in
[`helper-exports-api.log`](../../tools/revkit/work/stage16/helper-exports-api.log);
the probe is reproducible with `run-helper-exports.sh`.

### Remaining export-only text helpers: static trace

These wrappers are not declared by `vt_eng.h`. The export pseudocode routes
them to internal routines; the following are disassembly observations, not
runtime contracts:

| Export | Internal path observations | Still open |
| --- | --- | --- |
| `VT_TextToBufferEX_ENG` | Runtime selectors 0–2 match the corresponding ordinary buffer formats on the Stage 16 text and argument tuple; each finishes with one flag-1 poll. Static wrapper rejects selectors outside 0–2 with `-1` and dispatches to handlers at `0x100200c0`, `0x10020460`, and `0x10020930`. | Other pointer/scalar arguments, invalid inputs, multi-poll/chunk/cancel behavior, and extra-argument meanings. |
| `VT_TextToPreprocessInfoFile_ENG` | Runtime with byte flag 0 and every other argument null/zero completed normally; raw EAX after its decompiled `void` wrapper was 1. Static disassembly confirms this branch returns before pointer checks. With flag 1, Stage 16 text, speaker 1, and second argument `Z:/work/stage16/preprocess-flag1-probe.out`, low AX was `-6`; no file appeared at the supplied path. The wrapper is `void`, so this is a raw register observation. Static branches also show a four-resource loop for flag 4 and a `test.pcm` reference for flag 6; these branches were not run. See [`preprocess-noop.log`](../../tools/revkit/work/stage16/preprocess-noop.log), [`preprocess-nonzero.log`](../../tools/revkit/work/stage16/preprocess-nonzero.log), and runners `run-preprocess-noop.sh` / `run-preprocess-nonzero.sh`. | Meaning of the second argument, successful nonzero-flag behavior, file layout/output, errors, and meanings of other byte values. |
| `VT_TextToPcmBuffer_ProgressBar_ENG` | Runtime after model load with the Stage 16 text, thread 0, speaker 1, option values -1, and null window/message: 23,606 bytes were written, matching ordinary buffer format 0; guard bytes remained intact. Raw EAX after the decompiled `void` wrapper was 1, which is not treated as a declared return value. Static path uses speaker-indexed state and calls `PostMessageW` in the chunk-progress path. See [`buffer-progress.log`](../../tools/revkit/work/stage16/buffer-progress.log) and [`buffer-progress-start.bin`](../../tools/revkit/work/stage16/buffer-progress-start.bin), reproduced by [`run-buffer-progress.sh`](../../tools/revkit/work/stage16/run-buffer-progress.sh). | Long-text progress notifications, callback delivery, error cases, and other arguments. |
| `VT_TextToLipSyncLog_ENG` | Runtime after model load with Stage 16 text, empty second string, and remaining numeric args -1: raw EAX was 1 for speaker 1 and -1 (fallback). Null text produced AX -3; empty text produced AX -4 with raw EAX `0x003efffc`. A separate GDB-injected call with a nonempty path-like second argument terminated with an access violation before a result was captured; this single probe does not establish normal API behavior. The decompiled export wrapper is `void`; these register values are not declared C returns. Static routine checks loaded state and non-null/nonempty text, then creates and frees engine state. | Meaning of the second string, log format/destination, and output semantics. See [`lipsync-log.log`](../../tools/revkit/work/stage16/lipsync-log.log) and [`lipsync-output.log`](../../tools/revkit/work/stage16/lipsync-output.log); runners are `run-lipsync-log.sh` and `run-lipsync-output.sh`. |
| `VTDTTS_MakeInfo_ENG` | Runtime with Stage 16 text, speaker 1, remaining scalar arguments -1, and second argument `Z:/work/stage16/lead6-makeinfo-probe` completed with raw EAX `1`. It emitted `lead6-makeinfo-plead6.bin.dtt` (365 bytes) and `.asc.dtt` (2,596 bytes). Their headers read `VTDTTS BINARY` and `VTDTTS ASCII`; the ASCII output contains phone records with fields including file index, PCM/PM offsets and sizes, phone mode, pitch endpoints, and pitch/duration/volume rates. Static disassembly normalizes invalid speaker slots to slot 1, checks loaded/nonempty text, builds output records, then frees its internal context. The export wrapper is decompiled `void`, so EAX is a raw register observation. See [`makeinfo-api.log`](../../tools/revkit/work/stage16/makeinfo-api.log), both generated files, and [`run-makeinfo.sh`](../../tools/revkit/work/stage16/run-makeinfo.sh). | Why the generated basename does not directly match the supplied prefix; binary field widths/semantics, option effects, error cases, and whether output varies with text/settings. |
| `VT_VerifyTTS_ENG` | The internal validator applies the speaker-slot fallback and loaded/text checks before calling internal parsing and verification helpers. | Other argument combinations and result meanings. |

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
returns `-1`. The `void` `VT_UNLOAD_UserDict_ENG` wrapper discards that result.
These are static paths; only selected empty/invalid cases were called through
the header-declared wrapper in the runtime probe.

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
