# Lead 6: file, buffer, configuration, playback, and information API behavior (2026-09-25)

**Status: bounded API result; Lead 6 remains open.** This probe characterizes
the local 2013 Paul M16 file-output selectors, four buffer formats, selected
file and buffer errors, tested buffer flag/thread combinations and capacity
boundary, selected information requests, configuration setter/getter behavior,
selected text-format values and playback returns, and one repeated-call control. It does not
characterize the full buffer/thread lifecycle or general state across mixed
calls.

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
The separate format-0 capacity probe below establishes behavior for a smaller
incoming length value on this short synchronous call.

| Format | `flag` | `nThreadID` | Return | `output_len` after call |
| ---: | ---: | ---: | ---: | ---: |
| 0 | 0 | 0 | `1` | 23,606 |
| 0 | 1 | 0 | `-2` (`VT_BUFFER_API_ERROR_CREATE_THREAD`) | unchanged at 1,048,576 |
| 0 | 0 | 1 | `-2` (`VT_BUFFER_API_ERROR_CREATE_THREAD`) | unchanged at 1,048,576 |
| 0 | 1 | 1 | `-2` (`VT_BUFFER_API_ERROR_CREATE_THREAD`) | unchanged at 1,048,576 |

The nonzero `flag`/thread-ID cases failed at thread creation in this Wine
runtime. This records the observed failure for those arguments; it does not
show whether those modes work under another host/runtime or what their
successful processing/done sequence would be. The logs and successful byte
captures are `buffer-FORMAT-FLAG-THREAD.log` and
`buffer-FORMAT-FLAG-THREAD.bin` under Stage 16.

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
and unknown (`-9`) paths were not reached.

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
sequences.

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
when those call arguments are `-1`; it does not isolate the effect of each
field. The trace, log, baseline, and configured output are under Stage 16.

## Text-format argument

The header declares plain text (`0`), JEITA (`4`), JEITA Plus (`6`), and UMD
(`8`). With the same `Hello world.` input and all other arguments held at the
executable's defaults, each declared value and `-1` returned success from
`VT_TextToFile_ENG`; all five selector-4 WAVE outputs were byte-identical
(SHA-256 `a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`).
This establishes no distinction for this plain ASCII sample. It does not
establish format-specific markup, encoding, or malformed-input behavior. The
results are in `texttypes-api.log` and `texttype-*.wav` under Stage 16.

## Playback API

Direct calls to `VT_PLAYTTS_ENG` in the loaded process returned `-2` for a null
text pointer and `-3` for an empty string, matching the header's null/empty
text codes. A valid pointer to the executable's `Hello world.` text returned
`-5` (`VT_PLAY_API_ERROR_INITPLAY`) in this Wine run. The pseudocode reaches
the internal playback initializer for valid text and returns this error when
that helper reports failure; the probe does not identify the host-level cause
or establish successful playback.

The void `VT_STOPTTS_ENG`, `VT_PAUSETTS_ENG`, and `VT_RESTARTTTS_ENG` calls
returned to the harness without a crash after the failed play attempt. Their
pseudocode calls `waveOutPause`/`waveOutRestart` only when a global wave-output
handle exists, while stop resets and releases playback state. Since the valid
play call failed initialization, these calls do not verify active playback,
pause/resume behavior, completion notifications, or the thread lifecycle.
The captured calls are in `play-api.log` under Stage 16.

## Compatibility boundary

The public header declares buffer return values for processing/done and
additional errors. The incoming length was tested only for short synchronous
format-0 output; other formats, larger outputs, and error paths may behave
differently. Successful nonzero-thread processing, polling/completion
transitions, other error paths, and buffer-call persistence remain untested.
Playback success and active pause/resume behavior are also untested because
valid playback returns the initialization error in this Wine environment.
Apart from the selector 4 repeat control above, each format capture uses a
fresh process and one synthesis call. The setter probe changes several fields
together; per-field output effects and lower-clamp boundary outputs remain
untested. Text-format values were tested only with plain ASCII; format-specific
markup and malformed input remain untested.

The file-output selector map is now runtime-confirmed for the local Paul M16
engine at these arguments and output boundaries. This is not a cross-voice or
cross-version compatibility claim, nor whole-synthesis parity.
