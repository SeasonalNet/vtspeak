# Stage 21: MakeInfo alphabet-bank selection

## Extended user-dictionary memory bounds

At the public extended-loader breakpoint, a ten-byte `hello,HH,P` row followed
by `0xa5` loads with advertised length 10 and fails with length 11; the latter
failed slot remains empty. A span containing `hello,HH,P\0,X\0` loads, so the
parser stops at the embedded NUL. A valid first row followed by malformed
`x,HH,X` also loads, showing that validity of every byte after the first row
is not required. A two-row `hello`/`world` memory dictionary changed both
texts' WAV hashes under the forced dictionary gate; the second row therefore
is retained and affects synthesis in this tested case. This does not cover
arbitrary row counts or malformed-row placement.

Run in the isolated Stage 21 Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-userdict-memory-boundaries.sh
```

The runner saves and restores the sandbox input and output, and changes into
that directory before launching Wine. The accepted GDB trace and capture are
`trace-userdict-memory-boundaries-v3.gdb` and
`userdict-memory-boundaries-v3-api.log`. Earlier unversioned/v2 captures were
made before the runner set its working directory and remain exploratory.

The paired effect check uses a 21-byte buffer with two valid rows, captures
same-process no-dictionary controls and loaded-dictionary outputs for `hello`
and `world`, unloads the dictionary, and restores the gate. Both loaded
outputs match the known `hello,HH,P` dictionary WAV hash; both controls differ.
Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-userdict-memory-multiline-effect.sh
```

The trace, capture, and four-WAVE hash manifest are
`trace-userdict-memory-multiline-effect-v2.gdb`,
`userdict-memory-multiline-effect-v2-api.log`, and
`userdict-memory-multiline-hashes-v2.txt`. The first unversioned capture was
made before the runner set its working directory; v2 repeats it in the
isolated sandbox.

## Version data export mutation cross-check

At the sample host's first API breakpoint, the four version exports read
`(3, 11, 7, 1)`. A controlled in-process write changed them to
`(203, 204, 205, 206)`, then called `VT_GetDefVersion_ENG`; the returned
string remained `Paul-M16-FileIO`. The trace restored all four original dwords
before resuming the host. The capture contains only the four values and the
returned string. This shows that the version data exports are writable in the
loaded image and that the tested default-version getter does not compose its
result from them. It does not establish whether external hosts normally read
or write the cells, or whether other APIs consume them.

Run from the standard Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-version-export-mutation.sh
```

The GDB trace and capture are `trace-version-export-mutation.gdb` and
`version-export-mutation-api.log`.

## Speaker query selector boundaries

With Paul loaded in slot 1, `VT_GetSpeakerName_ENG` was called for
`INT_MIN`, -1, slots 0–5, 6, and `INT_MAX`. The four out-of-range inputs all
returned the fixed `Paul` fallback; valid slots returned the six compiled
speaker names. `VT_GetDBSize_ENG` normalized those same invalid selectors to
slot 1 and returned `1 / 508121688`. For valid but unloaded slots 0 and 2–5,
it returned `-1` and preserved a four-byte sentinel output. Its static code
performs the normalization before indexing the size table. The trace uses
exported data cell `0x100ff11c` as a temporary output word and restores its
original value before detaching.

Run from the standard Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-speaker-query-edges.sh
```

The GDB trace and capture are `trace-speaker-query-edges.gdb` and
`speaker-query-edges-api.log`. The runner checks the full ten-selector result
matrix and restores and compares the Stage 5 input and output fixtures.

## Speaker metadata copy contract

The follow-up queries all six `VT_GetPathKey_ENG` keys. Each call returns the
same `0x100fe6e0` buffer; after the sweep it contains the Ashley key, while a
saved copy retains the Kate key. Static code uses a 32-bit byte offset
`(selector*24) mod 2^32`. Selector `s` aliases valid slot `j` when
`s ≡ j (mod 2^29)`, so each of the six slots has seven additional signed
32-bit aliases. The wrap-alias matrix exercised all 42 additional values for
both metadata exports; every call returned the same key/name/path as its
ordinary slot index. The prior `INT_MIN` case is one member of the slot-0
alias class. This matrix maps only selectors whose wrapped offset lands on a
valid slot start; separately tested non-alias offsets -2, -1, 6, 7, and
`INT_MAX` fault.

The capture is
[`speaker-metadata-wrap-alias-wrap-matrix2-api.log`](./speaker-metadata-wrap-alias-wrap-matrix2-api.log),
generated by [`run-speaker-metadata-wrap-alias-v1.sh`](./run-speaker-metadata-wrap-alias-v1.sh)
from [`trace-speaker-metadata-wrap-alias-v1.gdb`](./trace-speaker-metadata-wrap-alias-v1.gdb).
Reproduce in a fresh runtime using a unique run ID:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-speaker-metadata-wrap-alias-v1.sh wrap-matrix3
```

For each valid `VT_SpeakersInfo_ENG` slot, the capture measures the lowercase
name and database path lengths, checks the copied NUL, and verifies that all
remaining bytes in the 512-byte `0xa5`-filled destination buffers stay
unchanged. The measured name/path lengths excluding NUL are `4/20`, `4/20`,
`5/20`, `5/23`, `5/18`, and `6/20`; the function returns `6` for every slot.
This constrains the observed copies and minimum buffer sizes for these six
records; the API itself has no capacity arguments.

### Invalid selector dereferences

Separate fresh-process probes call both exports with selectors `-2`, `-1`,
`6`, `7`, and `INT_MAX`. `VT_GetPathKey_ENG` terminates with Windows status
`0xc0000005` before it returns a pointer or string for every tested selector.
`VT_SpeakersInfo_ENG` does the same with valid 512-byte destination buffers,
so the failures occur on selector handling/table access rather than because
the caller supplied null or undersized destinations. In contrast, the
previous `INT_MIN` probe confirms 32-bit `selector*24` wrap to slot 0 for both
exports. These are five tested out-of-range values, not a range-wide result.

Run each case in a fresh process, supplying a unique second argument when
repeating a capture:

```sh
for selector in minus2 minus1 6 7 intmax; do
  docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
    /bin/bash /work/stage21/run-pathkey-oob-v1.sh "$selector" replay
  docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
    /bin/bash /work/stage21/run-speakersinfo-oob-v1.sh "$selector" replay
done
```

The captures are named `pathkey-oob-<selector>-v1-api.log` and
`speakersinfo-oob-<selector>-v1-api.log`; the matching `trace-*-oob-*-v1.gdb`
files record the direct calls. The runners classify only successful return
markers or the observed access-violation exit status.

Three more fresh-process calls use valid slot 0 but pass a null name buffer,
null path buffer, or both null. All terminate with `0xc0000005` before API
return. Static pseudocode copies the name before the path; the null-path runtime
capture does not preserve the partially written name after the process fault.
The API has no capacity arguments. Guard-page calls confirm the longest
name (slot 5, seven bytes including NUL) and longest path (slot 3, 24 bytes
including NUL) succeed at exact capacity. Each one-byte-short destination
faults at the protected page. The full 18-call matrix now repeats that test
for both outputs on all six slots: exact capacities succeed, and each
independent one-byte-short name or path faults. Exact byte capacities include
the terminating NUL:

| Slot | Name | Path |
| ---: | ---: | ---: |
| 0 | 5 | 21 |
| 1 | 5 | 21 |
| 2 | 6 | 21 |
| 3 | 6 | 24 |
| 4 | 6 | 19 |
| 5 | 7 | 21 |

```sh
for pointer_case in null-name null-path null-both; do
  docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
    /bin/bash /work/stage21/run-speakersinfo-null-v1.sh "$pointer_case" replay
done
```

The logs are `speakersinfo-null-name-v1-api.log`,
`speakersinfo-null-path-v1-api.log`, and
`speakersinfo-null-both-v1-api.log`, with corresponding
`trace-speakersinfo-null-*-v1.gdb` files.

Run all six exact pairs and twelve one-byte-short cases in fresh Wine
processes, serially:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-speakersinfo-guard-slot-matrix-v1.sh replay
```

Use a new final argument on later runs because the runner refuses to overwrite
existing captures.

The runner generates each direct-call trace from
`trace-speakersinfo-guard-matrix-template-v1.gdb`. Each process allocates two
64 KiB regions and applies `VirtualProtect` to the page immediately after each
destination; all protection calls returned 1 with original protection
`PAGE_READWRITE` (`0x4`). Captures are named
`speakersinfo-guard-slot<slot>-<exact|name-short|path-short>-<run-id>-api.log`.
The runtime matrix reported six exact-capacity returns of 6 and twelve
one-byte-short access violations (`0xc0000005`). Per-slot log links and the
full result table are in the [behavior report](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md#vtspeakersinfo-destination-pointer-and-capacity-behavior).

Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-speaker-metadata-contract.sh
```

The trace and result are `trace-speaker-metadata-contract.gdb` and
`speaker-metadata-contract-api.log`. The runner restores and compares both
Stage 5 fixtures.

## Configuration getter pointer and slot contract

The pitch/speed/volume/sentence-pause getter was called under all 16
combinations of its four output pointers being null or non-null; every call
returned 1, and only non-null outputs were written. The comma-pause getter
also returns 1 with a null output pointer. For both exports, `INT_MIN`, -1,
6, and `INT_MAX` normalize to loaded slot 1; valid but unloaded slots 0 and
2–5 return -1 and preserve the output sentinels. The slot sweep covers
`INT_MIN`, -1, 0–6, and `INT_MAX`.

Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-config-getter-contract.sh
```

The trace and capture are `trace-config-getter-contract.gdb` and
`config-getter-contract-api.log`; the runner checks the pointer-mask matrix,
selector results, and Stage 5 fixture restoration.

## `VT_GetUserDictLimit_ENG` selector domain

The helper returns `30, 10, 50, 65, 65` for selectors 0–4. The static
five-case switch returns `-1` for every other signed 32-bit value; runtime
calls confirm `INT_MIN`, -1, 5, 6, and `INT_MAX` also return -1. Reproduce
with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-limit-domain.sh
```

The trace and capture are `trace-userdict-limit-domain.gdb` and
`userdict-limit-domain-api.log`. The runner checks all ten results and restores
and compares the Stage 5 fixtures.

## User-dictionary target-normalizer complete single-byte domain

The complementary sweep tested all 160 non-NUL values omitted from the
printable-ASCII matrix: bytes `0x01–0x1f` and `0x7f–0xff`. TAB, LF, and CR
returned `-1`; all other values returned `-2`. Combined with the printable
matrix and the empty/NUL-input observation, this maps every single-byte input
value from `0x00` through `0xff`. Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-nonprintable.sh
```

The trace is `trace-userdict-target-normalizer-nonprintable.gdb`; the capture
is `userdict-target-normalizer-nonprintable-api.log`.

## User-dictionary source-normalizer two-byte domain

`probe-source-normalizer-pairs.c` calls private helper RVA `0x5f2e0` for all
65,536 byte-pair combinations with a NUL terminator. For every call it checks
the helper return, normalized output and terminator, and input immutability
against the trim and byte-pair predicates recovered from static analysis.
All calls matched: 1,257 pair rejects (`-3`), 276 empty results (`-1`), 2,259
one-byte copies, and 61,744 two-byte copies. The overall FNV-1a digest is
`bc75caddeff3fd2b`; all 256 per-first-byte row digests are in the log.

Build from the repository root:

```sh
bash tools/revkit/work/stage21/build-source-normalizer-pairs.sh
```

Run under the Stage 2 Wine prefix:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-source-normalizer-pairs.sh
```

The source, executable, runner, and complete capture are
`probe-source-normalizer-pairs.c`, `probe-source-normalizer-pairs.exe`,
`run-source-normalizer-pairs.sh`, and `source-normalizer-pairs-api.log`.

The context follow-up tests every rejected pair with each possible one-byte
prefix and suffix (643,584 calls). It checks return value, output prefix, and
input immutability; all calls matched. A prefix is copied before a later
rejected pair, and the helper returns `-3` without terminating that partial
output. The public wrapper preinitializes only byte 0 of its scratch buffer,
so the remainder on this error path is not a valid NUL-terminated string.
The overall context digest is `74330184b571cac7`. Build and run:

```sh
bash tools/revkit/work/stage21/build-source-normalizer-pair-contexts.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-source-normalizer-pair-contexts.sh
```

The harness, executable, runner, and capture are
`probe-source-normalizer-pair-contexts.c`,
`probe-source-normalizer-pair-contexts.exe`,
`run-source-normalizer-pair-contexts.sh`, and
`source-normalizer-pair-contexts-api.log`.

## User-dictionary target-normalizer length boundary

`VT_CheckUserDict_TargetNorm_ENG` was called on fresh buffers containing 63
through 68 `A` bytes. Lengths 63–65 returned `1`; 66–68 returned `-5`. The
tested over-limit calls left the text and original NUL position unchanged.
This pins the boundary for repeated ASCII letters; it does not establish how
trimming, multibyte input, or marker-bearing input affects the length count.
Run from the standard Stage 21 sandbox with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-lengths.sh
```

The trace is `trace-userdict-target-normalizer-lengths.gdb`; the capture is
`userdict-target-normalizer-lengths-api.log`.

## User-dictionary target-normalizer separator sweep

The trace calls `VT_CheckUserDict_TargetNorm_ENG` on fresh strings containing
`A`, 1–35 ASCII spaces, then `A`. One through nine spaces return `1`; ten
through 35 return `-7`. This identifies the first failing separator count for
this exact two-letter pattern. The undocumented result code and its behavior
for other token lengths or whitespace remain unresolved. Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-spaces.sh
```

The trace is `trace-userdict-target-normalizer-spaces.gdb`; the capture is
`userdict-target-normalizer-spaces-api.log`.

## User-dictionary target-normalizer token-count boundary

The direct-export matrix uses fresh buffers with one through 12 `A` tokens,
separated by one ASCII space. One through ten tokens return `1`; 11 and 12
return `-7`. The prior repeated-space matrix reaches the same code at ten
spaces between two `A` bytes. Static control flow is consistent with a
ten-segment counter, while the count of empty segments created by repeated
spaces remains an interpretation. Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-words.sh
```

The trace is `trace-userdict-target-normalizer-words.gdb`; the capture is
`userdict-target-normalizer-words-api.log`.

## User-dictionary target-normalizer branch cases

Direct calls reached these additional paths: `<AB>` returned `1`;
unterminated `<AB` returned `-8`; `<A B>` returned `-4`; a 30-byte `A`
segment followed by a space returned `1`, while 31 bytes returned `-6`; and
`A` followed by each of `a1 a1`, `ae a1`, or `fd fe` returned `-3`. The tested
caller buffers remained unchanged. These inputs exercise the branches but do
not establish meanings for the codes or the complete tag grammar. Reproduce
with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-branches.sh
```

The trace is `trace-userdict-target-normalizer-branches.gdb`; the capture is
`userdict-target-normalizer-branches-api.log`.

## User-dictionary target-normalizer marker cases

The marker trace tests exact prefix/suffix handling and mutation. `[SKIP] ` is
trimmed and returns `1`; lowercase `[skip]` returns `1`, while `[SKIP]x`
returns `-11`. `A[CI]`, `A[ci]`, and `A [CI]` return `2` and write NUL at the
opening bracket, leaving the space in `A `. `A[OTHER]` and `A[CI]B` return
`-12` without changing the caller buffer. A follow-up tested all 16 ASCII
capitalization masks of `[SKIP]` and all four of `[CI]`; all prefix forms
returned `1`, and all suffix forms returned `2` with truncation at the
opening bracket. Marker placement beyond these tested cases remains open. Run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-markers.sh
```

The trace is `trace-userdict-target-normalizer-markers.gdb`; the capture is
`userdict-target-normalizer-markers-api.log`.

The capitalization follow-up tested all 16 case masks for the four letters
in `[SKIP]` and all four case masks for `[CI]` after an `A`. Every `[SKIP]`
variant returned `1`; every suffix variant returned `2` and truncated at the
opening bracket. This exhausts ASCII letter capitalization for these two
exact marker spellings, not their placement grammar. Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-target-normalizer-marker-case.sh
```

The trace is `trace-userdict-target-normalizer-marker-case.gdb`; the capture
is `userdict-target-normalizer-marker-case-api.log`.

## Highlight-setting byte matrix

The probe calls `VT_SetTextTypeForHighlight_ENG` for every byte value from 0
through 255 in the loaded sample process and checks slot 1's field at
`+0x20424` after each call. The complete mapping is `0 → 0` and `1–255 → 1`;
all 256 calls matched and the original field value was restored. This maps the
setter's normalization only; effects of this field on highlighting or
synthesis were not measured. Reproduce it with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-highlight-byte-matrix.sh
```

The trace is `trace-highlight-byte-matrix.gdb`; the capture is
`highlight-byte-matrix-api.log`.

### Two-byte position-map classifier

`FUN_1001c900` accepts a two-byte sequence only when both bytes are non-NUL,
and its positive domain is exactly: first byte `0xa1..0xad`, second byte
`0xa1..0xfe`; first byte `0xae`, second byte `0xa1..0xc2`; or `0xfd 0xfe`.
The direct runtime sweep called the helper for all 1,257 predicted-positive
pairs. Every call returned 1; row counts were 94 for each lead `0xa1..0xad`,
34 for `0xae`, and one for `0xfd 0xfe`. The capture is
`highlight-two-byte-accepted-api.log`; replay with
`run-highlight-two-byte-accepted.sh` and
`trace-highlight-two-byte-accepted.gdb`.

A second attempted sweep called the helper on all 65,536 pairs. It timed out
at the runner's 600-second limit after completing lead byte `0xcd`; its
partial log is `highlight-two-byte-classifier-api.log` and is not a complete
runtime-domain result. Static branch conditions establish rejection of the
remaining pairs. The runners restore Stage 5 input/output fixtures; after
these runs, their backups matched the restored files byte-for-byte.

The effect probe reruns the sample's exact text and settings with the
highlight field first 0 and then 1, writing `highlight-off.wav` and
`highlight-on.wav` in the sandbox. Both calls returned 1; the 23,650-byte WAV
files are byte-identical (SHA-256
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`). This
does not test auxiliary highlight/index output. Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-highlight-audio-effect.sh
```

The trace and capture are `trace-highlight-audio-effect.gdb` and
`highlight-audio-effect-api.log`.

### EX record effects of the highlight setting

`FUN_1001c990` builds per-context dword maps at `+0x47790/+0x47794/+0x47798`
and remaps each EX marker row's `+8` through `+0x47794` alone when the flag is
0, or through `+0x47794` followed by `+0x47790` when it is 1. The latter map
expands character classes detected by `FUN_1001c900` to two output positions.
Separately, `FUN_10022dc0` uses the setting to choose the endpoint transform
that feeds SyncInfo `+20/+24`. Ghidra outputs are pseudocode; see
`../reports/stage21-highlight-flag-map.c`,
`../reports/stage16-syncinfo-row-producers.txt`, and
`../reports/stage16-syncinfo-source-coordinate-consumers.txt`.

The flag-paired selector-0 capture of
`<vtml_mark name="start"/>Hello world.<vtml_mark name="end"/>` returns 23,606
bytes and two descriptor/SyncInfo rows at each setting. Marker coordinates,
names, kinds, and captured SyncInfo fields match (`start +8=0`, `end +8=37`,
SyncInfo spans `[25,29]` and `[31,35]`). The 512-byte inline fields differ
after their NUL-terminated names: flag 0's start/end rows contain dword runs
`0..21` and `27..60`; flag 1's start row contains `0..60`, while its end row
contains a different byte sequence. However, `FUN_1001ccc0` initializes only
row `+0/+4/+8`, and the mark writer stores only the name and NUL in this
field. The post-NUL bytes are thus uninitialized row storage in the recovered
path, possibly allocator residue; they are not established as flag-dependent
API data. The writer's defined marker coordinates, name, kind, and the
separate SyncInfo endpoint mapping remain the evidence for flag behavior.
Captures are
`highlight-ex-raw-mark-{0,1}-api.log` plus
`highlight-ex-raw-mark-flag{0,1}-descriptor.bin`, and
`highlight-ex-raw-0-{0,1}-api.log` plus
`highlight-ex-raw-flag{0,1}-descriptor.bin`; the runners are
`run-highlight-ex-raw-mark.sh` and `run-highlight-ex-raw-records.sh`.

The single-pair API probe uses raw bytes `0xa1 0xa1` between the two named
marks in `input-ex-raw-highlight-two-byte.txt`. Both flag settings returned a
descriptor with two kind-1 rows but no output bytes or SyncInfo rows; the end
marker `+8` was `0x1b` with flag 0 and `0x1a` with flag 1. Its parsed-field
logs are `highlight-ex-two-byte-0-{0,1}-api.log`. The six-pair successful
boundary probe below supersedes it for completed synthesis behavior.

The stronger boundary probe places six positive-domain edge pairs
(`a1 a1`, `a1 fe`, `ad a1`, `ae a1`, `ae c2`, `fd fe`) between seven marks.
Both flag settings return seven kind-1 rows and one SyncInfo row under
selectors 0–2. Encoded lengths are 13,936 bytes for selector 0 and 6,968 bytes
for selectors 1/2; all report 6,968 audio frames. Flag 0 marker coordinates
are `0,24,48,72,96,120,144`; flag 1 coordinates are
`0,23,46,69,92,115,138` for every selector. The cumulative difference is one
position per preceding accepted pair. SyncInfo `+12/+16/+20/+24` were
`0x2e/0x2f/0x2e/0x2f` with flag 0 and `0x2d` in all four fields with flag 1;
the frame total remained `0x1b38`. This is consistent with the flag-1 map
treating each pair as one source position. It establishes behavior for these
boundary cases, not an encoding identity or the host's coordinate convention.
Replay selector 0 with `run-highlight-ex-two-byte-classes-replay-pair.sh` and
selectors 1/2 with `run-highlight-ex-two-byte-classes-matrix.sh`; captures are
`highlight-ex-two-byte-classes-replay-{0,1,2}-{0,1}-api.log` and the matching
descriptor dumps. `run-highlight-ex-two-byte-classes-replay-all.sh` runs both
sets serially. The first-pass parsed logs remain under the
`highlight-ex-two-byte-classes-{0,1,2}-{0,1}-api.log` names.

## Unit-selection history byte matrix

At a pre-load breakpoint, the trace calls `VT_SetUnitSelectHistoryMode_ENG`
for bytes 0–255 while the load flag is zero. Only input 1 stores 1; the other
255 values store 0. All 256 comparisons matched, the previous state was
restored, and a later text-call breakpoint confirms the normal loader completed
with the gate at 1. Existing post-load calls with 0 and 1 remain ignored. This
maps the complete byte input domain at both sides of the gate for the tested
slot/process, but does not recover all synthesis effects of the mode. Run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-unit-history-byte-matrix.sh
```

The trace and capture are `trace-unit-history-byte-matrix.gdb` and
`unit-history-byte-matrix-api.log`.

## Unit-selection history load fallback

With the history flag set before model load, the DLL constructs one `.his`
path per bank: `unit-gen.his`, `unit-num.his`, `unit-etc.his`, and
`unit-alp.his`, under `mc_idx_tbl/`. None is present in the mounted Paul M16
package. Each lookup returns a null source; the helper still returns low AX 1
and its static null-source branch clears two dword arrays. Runtime checks
confirmed zero at the first and last entries. The array counts
440,124/24,508/115,723/119 match the corresponding `.idx` unit counts. The
probe also confirms the normal format-4 call still returns 1. This explains
the enabled-mode initialization for this package, but not present-file `.his`
semantics or whether later accounting changes selection or repeatability.

Run from the standard Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-unit-history-load.sh
```

The trace and capture are `trace-unit-history-load.gdb` and
`unit-history-load-api.log`.

## Unit-selection history file parser

The loader reads each `unit-*.his` as a little-endian 32-bit unit count,
followed by two little-endian dwords per index unit. It stores the first dword
in one per-unit array and the second in another, and requires the header count
to equal the corresponding `.idx` count. A runtime synthetic overlay supplied
nonzero sentinels at the first and last entries for all four banks; all values
were read into the expected arrays, all four helper calls returned low AX 1,
and synthesis completed. Four trailing bytes beyond the required payload were
ignored. A one-less header count and a file shortened by the final dword each
returned low AX 0 from the history helper.

Those two rejected files also expose an enclosing loader failure: the bank
loop returned -1, then `VT_LOADTTS_EXT_ENG` continued with a null state pointer
and zero error word. The host hit an access violation at DLL address
`0x10027d9c`, instruction `mov [eax+0x4d08], edx`, with `eax=0` and attempted
write address `0x4d08`. This records the malformed-file failure path in this
build; it does not establish the meaning of the two persisted values. No
vendor data was changed: the synthetic files and symlinked model overlay live
in the container's disposable `/tmp`.

The runner also compares the valid synthetic-file output with a fresh run using
the no-file fallback, with history mode enabled in both. For the current Stage
5 fixture the outputs are byte-identical: 23,650 bytes, SHA-256
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`. This
bounds the observed result to this utterance and sparse sentinel records; it
does not establish general effects or the values' meaning. The outputs are
`unit-history-present-output.wav` and `unit-history-fallback-output.wav`, and
the no-file fallback trace is `unit-history-fallback-comparison-api.log`. The
runner prints `UNIT_HISTORY_OUTPUT_COMPARE identical=1` after its bytewise
comparison succeeds.

Run from the Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-unit-history-present.sh
```

The four captures are `unit-history-present-api.log`,
`unit-history-header-mismatch-api.log`, `unit-history-short-read-api.log`,
and `unit-history-trailing-data-api.log`. The trace is
`trace-unit-history-present.gdb`.

## Unit-history value pattern matrix

The second history probe fills every record in all four banks with
index-alternating values. Even indexes contain two zeros; odd indexes contain
`0x7fffffff` in field A, field B, or both. Each pattern is loaded before
synthesis of `Hello world.`, `The quick brown fox jumps over the lazy dog.`,
and `I saw 123 birds at 10:30.`. Each output is compared bytewise with a fresh
process using the missing-file fallback on the same text. All nine patterned
outputs match their baselines, and every load and text call returns 1.

Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-unit-history-pattern-matrix.sh
```

The result matrix is `unit-history-pattern-matrix-results.txt`; each run has
its own `unit-history-pattern-*-api.log` and output WAVE. The runner creates
its model overlay in `/tmp` and restores and compares the Stage 5 fixtures.
The negative audio result bounds these nine fresh-process utterance cases; it
does not identify either field or exclude other synthesis inputs, repeat
effects within one process, or non-audio bookkeeping.

Pass `repeat` to the same runner to make three additional direct
`VT_TextToFile_ENG` calls in the same loaded process after the host call. The
runner reuses the captured text/settings arguments, writes each repeat to a
separate relative file in the overlay's working directory, and compares each
populated-history result with the matching fallback call. All 36 additional
calls returned 1; all 27 later-call comparisons match byte for byte. Their
outputs use the `unit-history-pattern-*-repeat-*-output.wav` names, and the
same result file records `UNIT_HISTORY_PATTERN_REPEAT_COMPARE` lines. Re-run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-unit-history-pattern-matrix.sh repeat
```

The GDB call sequence is in `trace-unit-history-repeat.gdb`. This only bounds
three consecutive replays of each sampled text/settings tuple; other call
sequences and non-audio effects remain open.

## `VT_VerifyTTS_ENG` short markup and delimiter matrices

The 15-case matrix calls the internal validator directly at the first
file-synthesis entry on loaded Paul, with dictionary index `-1` and text type
`0`. It captures both low AX and EAX. For example, `<`, `>`, `<>`, `<<`, and
`>>` return `0xfffb` in low AX, while `</>`, `<A>`, unterminated `<A`, `A>`, `<A></B>`,
`A<B>`, and `A</B>` return `1`. This does not establish the markup grammar.

The second trace also calls the internal validator directly and tests all 340
strings of length 1–4 over `<`, `>`, `/`, and `A`. It reports each length,
base-4 string code, and low AX. The results are
287 successes and 53 `-5` results; all strings containing `A` succeed, while
the failures are among the 120 delimiter-only strings. This exhausts that
small product only. Full captures are `verify-tts-markup-matrix-api.log` and
`verify-tts-delimiter-product-api.log`; reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-verify-tts-markup-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-verify-tts-delimiter-product.sh
```

The traces are `trace-verify-tts-markup-matrix.gdb` and
`trace-verify-tts-delimiter-product.gdb`.

## `VT_VerifyTTS_ENG` single-byte domain

The direct-export sweep tested all first-byte values 0–255, followed by NUL,
with slot 1 loaded, dictionary index `-1`, and text type `0`. NUL returned low
AX `-3` (upper EAX was stale); 160 non-NUL values returned `1`, and 95
returned `-5`. The accepted byte ranges are listed in the behavior report;
all non-listed non-NUL values returned `-5`. This does not establish
multi-byte encoding behavior. Capture: `verify-tts-byte-domain-api.log`.
Replay with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-verify-tts-byte-domain.sh
```

Trace: `trace-verify-tts-byte-domain.gdb`.

## `VT_VerifyTTS_ENG` with a populated dictionary

This probe loads the Stage 16 `hello,HH,P` fixture at dictionary index 27,
then directly calls the export with six texts using default index `-1`,
loaded index 27 while the gate is off, and loaded index 27 while the speaker
gate byte at `0x100a7489` is debugger-forced on. An empty-slot index 28 is a
control with the gate on. Load returned AX 1 and index 27 held a non-null
pointer. In all three six-text groups, five inputs returned low AX 1 and the
malformed `<` returned -5; both empty-slot controls returned 1. Unload
returned AX 1 and the prior gate byte was restored. These results show no
status change for these inputs; they do not prove that the internal parsed
representation is unaffected by the dictionary. The complete capture is
`verify-tts-dict-api.log`; reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-verify-tts-dict.sh
```

The trace is `trace-verify-tts-dict.gdb`.

The expanded product repeats all 340 strings of length 1–4 over `<`, `>`,
`/`, and `A` with dictionary index 27 and the gate first off, then on. Both
groups matched the no-dictionary capture at every input: 287 low-AX `1`
results and 53 `-5` results. The gate-on state is debugger-forced. Run it
with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-verify-tts-delimiter-dict-product.sh
```

The capture and trace are `verify-tts-delimiter-dict-product-api.log` and
`trace-verify-tts-delimiter-dict-product.gdb`. Static pseudocode in
`tools/revkit/work/reports/stage23-state-array-writers.c` shows the selected
dictionary pointer flowing into nested language state only when the gate is
set; this does not by itself establish a user-visible validation effect.

## `VT_CheckUserDict_TargetPhon_ENG` complete single-byte domain

The direct export at `0x1002a590` was tested for every byte `0x00..0xff`,
followed by NUL, in a writable engine-heap buffer. NUL, TAB, LF, CR, and
space returned `-1`; the 16 uppercase phone tokens and `#` returned `1`;
`[` returned `-2`; the other 233 values returned `-9`. This includes all
single-byte high-bit values but does not map multibyte encodings. Capture:
`targetphon-byte-domain-api.log`.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-userdict-targetphon-byte-domain.sh
```

Trace: `trace-userdict-targetphon-byte-domain.gdb`.

## `VT_CheckUserDict_TargetPhon_ENG` token case masks

This 244-call matrix tested uppercase/lowercase single-letter phones, all four
case masks of the eight accepted consonant pairs, and all four masks for the
15 accepted vowel stems with suffixes 0–2. Exactly the 69 canonical uppercase
spellings returned `1`; all 175 lowercase or mixed-case variants returned
`-9`. This maps case sensitivity for the known inventory; the `[CI]` marker
has separate case-folding behavior. Capture: `targetphon-case-matrix-api.log`.
Replay with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-userdict-targetphon-case-matrix.sh
```

Trace: `trace-userdict-targetphon-case-matrix.gdb`.

Static loader cross-reference: `VT_LOAD_UserDict_ENG` calls `FUN_1003bd50`
for file parsing. That parser selects row handling from third-field code `A`
or `P`; the `P` branch uses private `FUN_10056790` and `FUN_10056890` on its
target field. These helpers do not call the exported target-phone checker and
show no explicit `[CI]` comparison. This means the export's accepted marker
is implemented below these top-level helpers or in a callee. Runtime probes
show uppercase and lowercase marker forms enable case-insensitive source
matching for the tested ASCII `P` rows; spacing after the phone also worked.
Other marker capitalization masks and broader target-token interactions
remain open. Source:
`tools/revkit/work/reports/vt_eng-2006-all-functions-pseudocode.c`, functions
at `0x10021400`, `0x1003bd50`, `0x10056790`, and `0x10056890`.

## Window destruction export

The direct-flow probe redirects the sample's file call into
`VT_DestroyWindow_ENG`. The export's stored window handle was `0x10074`; the
trace reached the call instruction for `DestroyWindow` at `0x10027fd6`, then
the inferior exited normally before the post-call instruction at `0x10027fdc`.
The Win32 return value was therefore not captured. This is a host-specific
process outcome and does not establish behavior in a different host. Reproduce
with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-destroy-window-direct.sh
```

The trace and capture are `trace-destroy-window-direct.gdb` and
`destroy-window-direct-api.log`.

The invalid-handle follow-up redirected target flow through the wrapper with
the stored handle set to `NULL` and then `0xdeadbeef`. Both post-call EAX
values were zero. It restored the original `0x10074` handle before resuming;
the inferior then exited with code 1. Replay with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-destroy-window-invalid-handle.sh
```

The trace and capture are `trace-destroy-window-invalid-handle.gdb` and
`destroy-window-invalid-handle-api.log`. This does not recover the result for
a valid window handle.

## Scalar helper setter edges

The loaded-Paul probe calls `VT_SetEmphasisFactor_ENG` at `INT_MIN`, both
clamp-adjacent values on either side of ±95, the boundaries, zero, and `INT_MAX`.
It also calls `VT_SetParenthesisCharNumber_ENG` and
`VT_SetEnglishReadingRule_KOR` at `INT_MIN`, `-1`, `0`, `1`, and `INT_MAX`.
The emphasis field clamps to `[-95,95]`; the other two fields clamp negatives
to zero and retain the positive probes, including `INT_MAX`. It restores all
three initial values. Other speaker slots are not covered.

The same trace sweeps `VT_SetSoundCardID_ENG` through `INT_MIN`, `-1`, `0`,
`1`, and `INT_MAX`. Every input is stored verbatim, and each call returns the
previous value; the initial global `-1` is restored. Valid device IDs are not
established by these calls.

Reproduce both sweeps with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-scalar-helper-edges.sh
```

The trace and capture are `trace-scalar-helper-edges.gdb` and
`scalar-helper-edges-api.log`.

## Playback request 101 query

`run-play-state-info.sh` reuses the Paul's sample application from the isolated
Stage 21 sandbox, redirects its file-synthesis call into `VT_PLAYTTS_ENG`, and
queries `VT_GetTTSInfo_ENG(101, NULL, sentinel, 4)` after play start, pause, and
restart. In all three phases, the API returned `1`, left the output sentinel at
`0x5a5a5a5a`, and the backing global at `1`. Thus request 101 reports the
started-session value through the pause interval in this run. The trace's stop
call exits the sample during debugger evaluation, so it captures no
post-stop/completion query. Reproduce it with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-play-state-info.sh
```

The trace is `trace-play-state-info.gdb`; its capture is
`play-state-info-api.log`.

A stop follow-up redirected target execution into `VT_STOPTTS_ENG` with a
target-stack return address at the sample's successful-play continuation
(`0x4016f9`). That earlier trace exited before the continuation because its
return-boundary breakpoint was placed on padding after the `ret`. A corrected
lifecycle trace now reaches the return and samples request 101 after cleanup.
The new runner, trace, and capture are `run-play-stop-lifecycle.sh`,
`trace-play-stop-lifecycle.gdb`, and `play-stop-lifecycle-api.log`.

## Playback stop lifecycle

The corrected trace starts the sample's valid-text call through
`VT_PLAYTTS_ENG` with the ALSA null sink, then lets the target reach its
successful-play continuation. It calls pause and restart, redirects target
execution into `VT_STOPTTS_ENG`, and gives stop the continuation address
(`0x4016f9`) as its return address. The earlier zero-handle cleanup inside
play startup is visible separately from the active stop.

For the active stop, the captured MMRESULTs were zero for `waveOutReset`, both
`waveOutUnprepareHeader` calls, and `waveOutClose`. The output handle was
`0xff00` during cleanup and zero afterward. The export reached its `ret` with
raw EAX 0, then the sample continuation ran. Request 101 returned status 1,
left its sentinel output unchanged, and its backing global remained 1. The
earlier idle cleanup had handle/global zero and request-101 status 0. These are
single-run state observations; they do not establish other stop phases, error
paths, completion notifications, or physical-device playback.

The separate `play-pause-restart-results-api.log` and its `v2`, `v3`, and
`v4` variants are exploratory setup captures, not accepted result evidence.
The initial run's injected null-handle call failed in the DLL frame; `v2`
missed the disabled playback-return breakpoint; `v3` failed while GDB was
calling the wrappers from an inner breakpoint; and `v4` repeated the calls
from the sample return boundary but later failed its stop-state query. Since
these exports are statically `void` wrappers that discard WinMM's MMRESULT,
their raw EAX values do not reveal a public return value. Use
`play-state-info-api.log` for the supported single pause/restart observation;
it records request 101 as `1` after both calls but does not verify playback
position or audible pause/resume.

The refined `run-play-control-mmresult.sh` executes the controls through the
target's normal stack flow, so breakpoints can distinguish the zero-handle
branch from the imported WinMM call. Before playback, handle zero bypasses
both imports. With the active null-sink handle `0xff00`, two pauses and two
restarts each reach the import and return MMRESULT 0. The runner then stops
playback through target flow and exits normally. This measures Wine's null
sink behavior; it does not demonstrate audible pause/resume. Its accepted
trace and capture are `trace-play-control-mmresult-v4.gdb` and
`play-control-mmresult-v4-api.log`. The v4 capture repeats the earlier v3
trace from the isolated Stage 21 sandbox. The earlier
`play-control-mmresult-api.log` stopped inside GDB call evaluation; v2 treated
the shared return/join address as evidence of an imported-call result and did
not reach stop cleanup. Those logs are retained as failed probe attempts, not
runtime results.

## Playback callback messages and natural completion boundary

The initial GDB-driven pump received `MM_WOM_OPEN` (`0x3bb`) and
`MM_WOM_DONE` (`0x3bd`) messages, but the inferior exited inside its synthetic
`DispatchMessageA` call. A standalone PE32 host then called the exports and
pumped its normal thread queue. Request 101 changed from `1` to `0` after the
last done message for the seven original texts and two UTF-8 follow-ups.
Done-message counts ranged from one to eight. Caller-message `wParam/lParam`
values match inclusive source-byte spans: `A A A` produced `(0,0)`, `(2,2)`,
`(4,4)`; `Hello, world!`
produced `(0,4)` and `(7,11)`; CP1252 `café noir` produced `(0,3)` and
`(5,8)`. The final caller message was `(0,-1)` in every capture, and
`(31,38)` repeated for the long sentence. The callback pseudocode shows
`(0,-1)` is posted on the terminal branch after request 101 is cleared, then
`VT_STOPTTS_ENG` is called. With UTF-8 `café noir`, spans were `(0,3)`,
`(4,4)`, `(6,9)`; UTF-8 `é noir` gave `(0,0)`, `(1,1)`, `(3,6)`. The two
bytes of the tested UTF-8 `é` are reported at separate byte positions. Other
UTF-8 inputs, DBCS, and the cadence/repetition rules remain open. For the long
sentence, two `MM_WOM_DONE` messages were dispatched before the next two
caller updates, and `(31,38)` was sent twice. This shows queued notifications
can lag audio completions; it does not identify the exact block-to-span map.
Rebuild with
`bash tools/revkit/work/stage21/build-playback-natural-completion.sh` and run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-playback-natural-completion.sh
```

The direct-host source and captures are `probe-playback-natural-completion.c`,
`playback-natural-completion-host-api.log`,
`playback-notification-matrix-api.log`, and `playback-span-boundaries-api.log`.
The repeatable matrix runners are `run-playback-notification-matrix.sh` and
`run-playback-span-boundaries.sh`. The earlier debugger trace and capture are
`trace-play-state-natural-completion.gdb` and
`play-state-natural-completion-api.log`; pseudocode is in the Stage 26 reports.
Reproduce the CP1252/UTF-8 comparison with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-playback-encoding-matrix.sh
```

The capture is `playback-encoding-matrix-api.log`.

## Heap-start export access watch

The access-watch trace arms a hardware read/write watchpoint on the complete
four-byte cell exported as `VT_gHeapStartAddress_ENG` (RVA `0xff11c`) before
the sample calls `VT_LOADTTS_ENG`. It then allows the normal sample file
synthesis call to complete and explicitly unloads speaker 1. The cell is zero
before load, after synthesis (which returned EAX 1), and after unload. The
watchpoint reports zero accesses across the observed process path. This
narrows the ordinary-path behavior but does not establish how a host might
write or consume the export, or whether another engine path accesses it.
The four bundled voice-specific executables each import six functions from
their matching VoiceText DLL and none imports the five data exports. Although
they import `GetProcAddress`, static disassembly shows its only reference
resolves `___lc_codepage_func` or `__lc_codepage` from `msvcrt.dll` in
`_init_codepage_func`. There is no bundled-host dynamic VoiceText export
lookup in these executables. Other hosts remain outside this check. The
import/callsite survey and reproduction commands are in
`heap-export-host-lookup-static.txt`.

All four matching DLLs use the same separate Windows heap handle path: the
`HeapCreate` call at VA `0x10066924` stores its return value at `0x101004a0`,
which the allocation helpers pass to `HeapAlloc`, `HeapFree`, `HeapReAlloc`,
and `HeapDestroy`. This handle is distinct from the exported cell at
`0x100ff11c`; the exported cell's intended role remains unknown.

Reproduce from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-heap-export-access.sh
```

The trace and capture are `trace-heap-export-access.gdb` and
`heap-export-access-api.log`.

## No-op setter aliases

After the sample's file synthesis returned EAX 1, the trace directly called
the shared target at RVA `0x28420`, associated in the PE export table with
`VT_SetDecimal0Pron_ENG`, `VT_SetPhone0Pron_ENG`, and
`VT_SetVirtualTagMode_ENG`. The three calls supplied distinct EAX sentinels.
Each returned with EAX and ESP unchanged; the loaded speaker pointer,
play-state, and history-mode observations also stayed unchanged. This runtime
call result agrees with the one-instruction `ret` body. EAX values are
preservation sentinels, not API return values. Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-noop-setter-aliases.sh
```

The trace and capture are `trace-noop-setter-aliases.gdb` and
`noop-setter-aliases-api.log`.

## English-reading-rule consumer probe

The runner performs two fresh-process runs with the same format-4 text:
`Hello world. I read 123 books (three times) at 10:30 in the U.S.A.`. At
`VT_TextToFile_ENG` entry it sets the global `DAT_100a0460+0x20420` field to
0 or 1, then arms a hardware read/write watchpoint through API return. Both
calls return EAX 1, the field receives the requested value, and no later access
is observed. The two WAVs are byte-identical with SHA-256
`a54bcb6ea61bbb268dea4404a62bc6b388be82836c057d47a305e2fdabc89739`. This
shows no setting consumption on this utterance and format; other text, formats,
voices, and host-side uses remain open.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-english-reading-rule-effect.sh
```

The command captures `english-reading-rule-0-api.log`,
`english-reading-rule-1-api.log`, and both WAV outputs. Its GDB template is
`trace-english-reading-rule-effect.gdb.in`.

### Expanded text and stored-value matrix

The follow-up runner repeats the 0/1 comparison for the original sample,
homograph/heteronym sentences, date/currency/unit text, abbreviations, and an
acronym sample. It also tests `INT_MAX` on the date/currency sample. Each case
runs in a fresh process; the debugger sets the global field at the
`VT_TextToFile_ENG` entry and watches it until the API returns. All seven
paired WAVs are byte-identical, all thirteen calls return 1, and no watchpoint
reports a post-setter read or write.

Run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-english-reading-rule-matrix.sh
```

The captures, outputs, and comparison list are named
`english-reading-rule-matrix-*`. These cases show no consumer on the sampled
loaded-Paul file-synthesis paths; they do not cover other formats, voices,
APIs, or host-side reads.

This stage extends the Stage 16 `File Index` matrix to identify runtime cases
that select `merged-alp`. It uses a dedicated minimal working directory at
`sandbox/stage5/`, with its own `input1.txt` and empty `output.wav`. The
`sandbox/data-paul` and `sandbox/data-common` links preserve the sample's
relative model paths and resolve to the existing read-only vendor mounts. Run
each probe with Compose working directory `-w /work/stage21/sandbox/stage5`;
the scripts refuse to start from any other directory. Wine-gdb's remote target
detaches the sample when the trace ends, so the sample can continue and write
its normal `output.wav` in its current directory. The working-directory
override keeps that output inside the Stage 21 sandbox.

## Probe and result

The first matrix called `VTDTTS_MakeInfo_ENG` once for each standalone
uppercase letter A–Z and four controls: `ABC`, `A B C`, `A. B. C.`, and
`U.S.A.`. A second matrix repeated all 26 letters in lowercase. A third
22-case probe varied terminal punctuation for A, B, and C, then compared short
letter sequences with and without punctuation. Uppercase and lowercase
punctuation probes each applied `.`, `,`, `!`, and `?` to every letter A–Z
(104 calls apiece). Uppercase and lowercase pair matrices each covered all 676
pairs with a space, with and without a terminal question mark (1,352 calls
apiece). Two mixed-case matrices also covered all 676 pairs
in `Upper lower` and `lower Upper` forms (1,352 calls apiece). A further
8,112 calls covered all pairs ending in `.`, `,`, and `!` across those same
four case patterns. All 13,806 calls returned raw EAX `1` and emitted ASCII
and binary DTT files. The A-leading matrix first tested all 676 continuations
under `UUU`, `LLL`, `ULL`, and `LUU`, each plain and with terminal `?` (5,408
calls), then added `LLU`, `LUL`, `ULU`, and `UUL` (5,408 calls). The first
`LUL` and `ULU` batches used an incorrect middle-token case mapping; their
2,704 calls returned raw EAX `1` but are quarantined under `invalid-case-map/`
and excluded. Corrected replacement batches are included in the validated
46,254-capture set. Two eight-mask matrices tested all 676 surrounding-letter
pairs with `A` in the middle and at the end, plain and with terminal `?`
(21,632 calls). All 46,254 retained captures returned raw EAX `1` and emitted
ASCII and binary DTT files. A supplemental 64-call four-token probe is
documented below and is outside that aggregate.
Uppercase and lowercase calls
selected the same bank for every letter. The per-letter map was:

| Bank | Letters |
| --- | --- |
| `merged-alp` (index 3) | A, B, E, G, I, J, M, N, Q, R, U, V, W, Z |
| `merged-etc` (index 2) | C, D, F, H, K, O, P, S, T |
| `merged-gen` (index 0) | L, X, Y |

All four spelling/acronym controls used `merged-gen`. The punctuation probe
showed that a terminal period, comma, or exclamation mark preserved the
standalone bank for A, B, and C. A question mark changed A to `merged-etc`, B
to `merged-gen`, and C produced rows from both `merged-gen` and `merged-etc`.
Phrase contexts also varied by sequence: `A B`, `A B C`, and `B C D` used
`merged-gen`; `B C` used `merged-etc`; `A, B` mixed a gen row for A with alp
rows for B. Both `A.B.C.` and `A. B. C.` used gen throughout. The 78 captures
contain 231 TypeFlag=2 phone/unit rows: 73 gen, 60 etc, and 98 alp; none used
`merged-num`. One additional TypeFlag=1 silence row occurs in `A, B`.
Every TypeFlag=2 row's `PCM Pos` matches both a DAT offset and a unit ordinal
in the unit index corresponding to its `File Index`, including all index-3 rows. The selected
`merged-alp` unit ordinals follow the alphabetic record sequence for the
letters that choose that bank: A=0, B=1–2, E=7, G=10–11, I=14, J=15–16,
M=21–22, N=23–24, Q=28–30, R=31–32, U=37–38, V=39–40, W=41–47, and Z=53–54.
This directly demonstrates runtime selection of `merged-alp`, case invariance
for isolated letters, and punctuation/context-sensitive selection. The
alphabetic ordinal pattern suggests that at least part of `unit-alp.idx` stores
letter-name unit sequences; no unit labels or format specification confirm
that semantic name. These controlled contexts do not explain why the other
letters select `merged-etc` or `merged-gen`, or generalize phrase-context
behavior.

The complete terminal-punctuation matrices sharpen that boundary. In both
cases, for every letter A–Z, the TypeFlag=2 rows for a trailing period, comma,
and exclamation mark exactly match the corresponding standalone capture,
field-for-field (156/156 comparisons). Every trailing question-mark capture
differs from its standalone row list: nine use only `merged-gen`, two only
`merged-etc`, and fifteen contain rows from both banks. For all four marks,
the uppercase and lowercase row lists also match exactly (104/104 comparisons).
Neither 104-case matrix emitted a TypeFlag=1 silence row. All 208 calls
returned raw EAX `1`, and every phone-row offset cross-checks to the selected
bank's DAT/index record. This establishes the isolated-letter punctuation and
case behavior at the tested default options. It does not explain the
question-mark context rule or settle where/when silence rows are inserted in
multiword or longer text.

The ordered-pair matrices show that neighboring letter tokens materially
alter the result. In both cases, all 676 plain-versus-question-mark comparisons
change the TypeFlag=2 record list; 430 also change the bank set, while 246 keep
the same set but select different unit records. Pair-row counts are:

| Pair input | gen | num | etc | alp |
| --- | ---: | ---: | ---: | ---: |
| Uppercase, plain | 3,429 | 106 | 174 | 0 |
| Uppercase, `?` | 3,194 | 76 | 833 | 6 |
| Lowercase, plain | 3,426 | 110 | 184 | 0 |
| Lowercase, `?` | 3,198 | 80 | 837 | 6 |

The six alp rows occur only in five question-mark pairs: `E P?`, `Q J?`,
`T F?`, `T S?`, and `V P?`. No pair capture emits a TypeFlag=1 silence row.
Uppercase and lowercase pair outputs match exactly in 1,302 of 1,352
comparisons. All 50 differences begin with `A`/`a`; in the first phone row,
uppercase cases use `AH0` and lowercase cases use `EY1`. Seventeen of those
50 cases also change the selected bank set, and 33 change records within the
same bank set. The two mixed-case matrices add 2,704 calls: both `Upper lower`
and `lower Upper` outputs match the all-uppercase output field-for-field for
all 1,352 pair/form combinations in each pattern. Thus the observed
all-lowercase effect for leading `a` requires both tokens to be lowercase in
these tested two-token cases; lowering either token alone has no effect. This
is a bounded case interaction result, not an identified lexical/context rule.
Other four-token sequences and free-text contexts remain untested.

For exactness, the 17 of those 50 lowercase/uppercase comparisons that change
the selected bank set are: `aa plain` gen → gen+etc; `ad plain` gen → etc;
`ad?` gen+etc → etc; `ae?` gen+etc → gen; `af plain` gen → gen+num;
`al plain` gen → gen+num; `al?` gen+etc → gen+num+etc; `an plain` gen →
gen+num; `ao plain` gen+num → gen; `aq?` gen → gen+etc; `as plain` gen →
gen+num; `as?` gen → gen+etc; `at plain` gen → gen+num+etc; `at?`
gen+etc → gen+num+etc; `av plain` gen → num+etc; `av?` gen+etc →
gen+num+etc; and `ax plain` gen → gen+num. The other 33 differing outputs
retain the same bank set and select different records. The analyzer reports
both lists so the individual row differences can be reproduced from the
captures.

The pair punctuation matrices compare each ending against its same-case
unpunctuated pair. Uppercase `.` and `,` preserve every row list; `!` preserves
675/676, with only `AA!` changing rows (the bank set stays the same). Lowercase
`,` and `!` preserve all 676 row lists. Lowercase `.` preserves 670/676;
`a a.`, `a e.`, `a i.`, `a o.`, `a u.`, and `a y.` change rows, and `a a.` and
`a o.` also change the bank set. In each of those six cases, the lowercase
period result exactly matches the uppercase unpunctuated result, so the ending
removes the lowercase-only row variation seen without punctuation. Both
mixed-case patterns preserve all rows
for `.` and `,`; `!` preserves 675/676, with only the `aa` pair changing rows.
For all three marks, mixed-case outputs match uppercase field-for-field for
all 2,028 comparisons per pattern. All-lowercase versus uppercase differs in
69/2,028 terminal-punctuation comparisons: 19 period pairs, 25 comma pairs,
and 25 exclamation pairs, each difference beginning with `a`. The analyzer
prints the complete pair lists. No terminal-punctuated pair emitted a
TypeFlag=1 silence row.

The three-token matrix adds two letter positions of context after `A`/`a`.
Every uppercase, `A b c`, and `a B C` plain/question comparison changes the
row list; 415 change bank sets and 261 retain the same set. The all-lowercase
matrix has 394 bank-set changes and 282 same-bank record changes. Comparing
all-lowercase with uppercase, 649/676 suffixes differ in each form (1,298 of
1,352 outputs); the only exact suffix matches are `ga` through `gz` and `ne`,
for both forms. Both mixed-case patterns match uppercase field-for-field for
all 1,352 outputs. No three-token capture emits a TypeFlag=1 row. This extends
the all-lowercase sensitivity observed for two-token inputs into this tested
three-token set, but does not establish the general context or tokenization
rule.

Across all 46,254 captures, the analyzer finds 322,958 TypeFlag=2 detailed
rows: 289,681 gen, 10,689 num, 21,948 etc, and 640 alp, plus one TypeFlag=1
row in `A, B`. Every TypeFlag=2 `PCM Pos` resolves to a unit in the selected
bank.
The aggregate analyzer checks every row offset, compares isolated punctuation
records with standalone records in both cases, compares uppercase with
lowercase for all four marks, compares each pair's plain and question forms,
compares case patterns across all pairs, and checks period/comma/exclamation
against plain pairs in all case patterns. It compares the A-leading, A-middle,
and A-final three-token matrices under all eight per-token upper/lower masks,
including each plain/question result and each output's record rows against
uppercase.

For each three-letter form, the eight masks are `UUU`, `LLL`, `ULL`, `LUU`,
`LLU`, `LUL`, `ULU`, and `UUL` (token order; `U` is uppercase and `L` lowercase).
The question-mark result depends on the full case mask as well as A's position:

| A position | Case mask | Plain→`?` bank changes | Same-bank record changes | Outputs differing from `UUU` |
| --- | --- | ---: | ---: | ---: |
| Leading | `UUU` | 415 | 261 | 0/1,352 |
| Leading | `LLL` | 394 | 282 | 1,298/1,352 |
| Leading | `ULL` | 415 | 261 | 0/1,352 |
| Leading | `LUU` | 415 | 261 | 0/1,352 |
| Leading | `LLU` | 398 | 278 | 1,298/1,352 |
| Leading | `LUL` | 415 | 261 | 0/1,352 |
| Leading | `ULU` | 415 | 261 | 50/1,352 |
| Leading | `UUL` | 415 | 261 | 0/1,352 |
| Middle | `UUU` | 352 | 324 | 0/1,352 |
| Middle | `LLL` | 351 | 325 | 52/1,352 |
| Middle | `ULL` | 352 | 324 | 0/1,352 |
| Middle | `LUU` | 352 | 324 | 0/1,352 |
| Middle | `LLU` | 391 | 285 | 1,302/1,352 |
| Middle | `LUL` | 352 | 324 | 0/1,352 |
| Middle | `ULU` | 388 | 288 | 1,300/1,352 |
| Middle | `UUL` | 352 | 324 | 0/1,352 |
| Last | `UUU` | 141 | 535 | 0/1,352 |
| Last | `LLL` | 143 | 533 | 50/1,352 |
| Last | `ULL` | 141 | 535 | 0/1,352 |
| Last | `LUU` | 141 | 535 | 0/1,352 |
| Last | `LLU` | 140 | 536 | 100/1,352 |
| Last | `LUL` | 141 | 535 | 0/1,352 |
| Last | `ULU` | 137 | 539 | 52/1,352 |
| Last | `UUL` | 141 | 535 | 0/1,352 |

For leading A, `LLL` and `LLU` each differ from UUU in 1,298/1,352 outputs;
`ULU` differs in 50, while `ULL`, `LUU`, `LUL`, and `UUL` match UUU. Thus
these results do not reduce to a rule based only on A's case or on whether the
whole input is lowercase. The ULU differences are the same 25 continuations in
both forms: `aa`–`af` and `ah`–`az` (`ag` is unchanged). Of those, ULU changes
the bank set for 9 plain and 7 question outputs; the remaining 16 and 18,
respectively, change records within the same bank set. For `B A C`, the
partial-lowercase masks
`LLU` and `ULU` produce large changes from uppercase; for `B C A`, `LLU` and
`ULU` are exact matches while `UUL` differs in 52 outputs. Thus case effects
depend on the whole token mask and A's position, not on a blanket case-folding
rule. These are exhaustive results for three fixed three-token shapes, not a
recovered rule for arbitrary text.

A four-token follow-up first tested `A A A D` and `A A G D` under all 16
per-token case masks, each plain and with terminal `?` (64 calls). A further
400 calls swept the final token across A–Z for those same two letter contexts
under `UUUU`, `UUUL`, `ULUU`, and `ULUL`; the `D` results come from the first
set. These 464 supplemental captures are separate from the 46,254-capture
three-token aggregate. All returned raw EAX `1`; all 4,952 TypeFlag=2 rows
resolve to selected-bank units, and none emitted TypeFlag=1.

For `A A A D`, comparison with UUUU found seven masks with identical rows in
both forms, seven with same-bank record changes, and two with bank-set changes.
For `A A G D`, twelve masks matched in both forms and four changed records
within the same bank sets. Terminal `?` changed bank sets for all 16 masks in
both D contexts. Across the A–Z fourth-token sweep, `UUUU` and `UUUL` match
exactly for both contexts. For `A A A X`, ULUU and ULUL differ from UUUU for
all 26 final letters in both forms: plain outputs change banks for `X=a,b`,
while question outputs change banks for `X=b,e,h,j,m,n,s`; the other differing
outputs retain their bank set. For `A A G X`, all four tested masks match UUUU
for every final letter and form. The selected three-token contrast therefore
persists after any tested fourth letter in the `A A A X` context, but not in
the matched `A A G X` context. This is a bounded two-context result, not an
exhaustive four-token or free-text rule.

A full case-mask extension combines that 5,408-call set with 16,224 calls in
three additional batches. Together the captures cover all 676 ordered letter
pairs `X Y` in `A A X Y` and `A a X Y`, all 16 per-token upper/lower masks,
and both plain and terminal-`?` forms: 21,632 calls, 43,264 paired ASCII/BIN
files, separate from the 46,254-capture three-token aggregate. Every call
returned raw EAX `1`; all 220,076 TypeFlag=2 rows map to indexed units, and no
TypeFlag=1 row appeared. Across the 16 masks, every
plain-to-question comparison changes rows for all 676 pairs. The number of
bank-set changes / same-bank record changes is: UUUU, UUUL, UULL, LUUU, LUUL,
and LULL: 402/274; UULU and LULU: 403/273; ULUU and ULUL: 416/260; ULLU:
404/272; ULLL: 403/273; LLUU and LLUL: 410/266; LLLU: 410/266; LLLL:
410/266.

Field-for-field equivalence across all pairs and both forms groups the masks
as `UUUU UUUL UULL LUUU LUUL LULL`; `UULU LULU`; `ULUU ULUL`; `ULLU`;
`ULLL`; `LLUU LLUL`; `LLLU`; and `LLLL`. Thus lowercasing token 1 has no
independent effect when token 2 remains uppercase, but interacts with a
lowercase token 2 in other masks. Lowercasing token 4 has no effect in several
classes, though it is not generally inert. Against UUUU, UULU/LULU differ in
25/676 pairs per form (9 plain and 6 question bank-set changes); ULLU differs
in 26 per form (10 plain, 6 question bank-set changes); ULLL differs in one
plain pair (bank-set change) and one question pair (same-bank record change).
ULUU/ULUL differ in 650 pairs per form, with `ga`–`gz` as the only exact
matches. LLUU/LLUL and LLLU/LLLL differ for every pair. These results complete
the casing matrix for the two selected prefixes, not the selector rule for
other initial letters, token strings, punctuation, or free text.

An earlier 192-call fixed-`D` third-letter slice is also retained; the full
matrix above covers its same inputs and masks. It is redundant evidence, not a
separate expansion of the validated full matrix.

## First-two-letter identity with a fixed `A A` suffix

The next matrix varies the first two letters over all 676 ordered pairs
`X Y`, then supplies fixed `A A`. It covers all 16 per-token upper/lower case
masks and both plain and terminal-`?` forms. Four batches partition the masks
without overlap. This set is separate from both the 46,254-capture three-token
aggregate and the 21,632-call fixed-prefix four-token matrix above.

The four batches contain 5,408 API calls apiece (21,632 total) and 10,816
paired ASCII/BIN DTT files apiece (43,264 total). Every call returned raw EAX
`1`. All 204,466 TypeFlag=2 rows map to indexed units at their selected-bank
PCM offsets; no TypeFlag=1 row appeared. Row totals by batch are b1 51,096,
b2 51,100, b3 51,096, and b4 51,174. Every terminal-question comparison
changes records for all 676 pairs. The bank-set-change / same-bank-record-change
counts by mask are:

| Masks | Bank set changed | Records changed, same banks |
| --- | ---: | ---: |
| UUUU, UUUL, UULL, ULLL, LUUU, LUUL, LULL | 180 | 496 |
| UULU, ULLU, LULU | 30 | 646 |
| ULUU, ULUL | 191 | 485 |
| LLUU, LLUL | 186 | 490 |
| LLLU | 31 | 645 |
| LLLL | 175 | 501 |

Across all pairs and both forms, the field-for-field mask equivalence classes
are `UUUU UUUL UULL ULLL LUUU LUUL LULL`; `UULU ULLU LULU`; `ULUU ULUL`;
`LLUU LLUL`; `LLLU`; and `LLLL`. Against UUUU, UULU/ULLU/LULU differ on all
676 pairs in each form (254 plain and 129 question outputs change bank sets).
ULUU/ULUL differ on 26 pairs per form (16 plain and 4 question bank-set
changes); LLUU/LLUL differ on 50 pairs per form (24 plain and 12 question
bank-set changes); LLLU differs on all 676 pairs; LLLL differs on 25 pairs
per form. These are observed equivalences for this fixed suffix and input
shape, not a general capitalization or selector rule. Other suffix letters,
non-letter prefixes, punctuation beyond terminal `?`, free text, and broader
option/error behavior remain open.

Run the matrix validators and combined analysis from the repository root:

```sh
python3 tools/revkit/work/stage21/check_makeinfo_letter_vary_prefix_fixed_aa_batch.py b1
python3 tools/revkit/work/stage21/check_makeinfo_letter_vary_prefix_fixed_aa_batch.py b2
python3 tools/revkit/work/stage21/check_makeinfo_letter_vary_prefix_fixed_aa_batch.py b3
python3 tools/revkit/work/stage21/check_makeinfo_letter_vary_prefix_fixed_aa_batch.py b4
python3 tools/revkit/work/stage21/analyze_makeinfo_letter_vary_prefix_fixed_aa_matrix.py
```

Each batch runs with the Stage 8 Docker/Wine environment and the matching
`run-makeinfo-letter-vary-prefix-fixed-aa-matrix.sh bN` argument. Its API log,
GDB trace, file listing, and paired captures use the
`makeinfo-letter-vary-prefix-fixed-aa-bN-*` and
`mi-letter-vary-prefix-fixed-aa-bN-pair-*` prefixes.

For example, run b1 from the repository root with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-letter-vary-prefix-fixed-aa-matrix.sh b1
```

Use `b2`, `b3`, and `b4` for the remaining disjoint mask groups.

## Interword punctuation and TypeFlag=1 intervals

This follow-up tests all 676 ordered letter pairs as `X, Y`, `X. Y`, or
`X... Y`, with four two-token case masks (`UU`, `UL`, `LU`, `LL`). Each mark
has 2,704 calls and 5,408 paired ASCII/BIN files; the complete matrix has
8,112 calls and 16,224 captures. All calls returned raw EAX `1`, and every
TypeFlag=2 PCM position maps to an indexed unit in the selected bank. Every
capture containing TypeFlag=1 has exactly one such row, with TypeFlag=2 rows
before and after it.

| Interword text | UU | UL | LU | LL | TypeFlag=1 size |
| --- | ---: | ---: | ---: | ---: | ---: |
| `X, Y` | 676 | 676 | 676 | 676 | 3,200 samples |
| `X. Y` | 0 | 0 | 675 | 675 | 14,800 samples |
| `X... Y` | 52 | 0 | 676 | 676 | 14,800 samples |

For `X. Y`, the sole lower-leading exception is `n. e` in both `LU` and `LL`.
For uppercase `X... Y`, the 52 intervals are exactly `X... A` and `X... I`
for every first letter A–Z. At the established 16 kHz sample rate, 3,200
samples correspond to 200 ms and 14,800 samples to 925 ms. These are direct
results for one-letter tokens with the punctuation between them; they do not
establish behavior for arbitrary words, terminal punctuation, or pause-setting
combinations.

## Ordinary-word interword context grid

To extend the one-letter matrix to multi-character lexical tokens, this grid
tested all 144 ordered pairs from `a`, `i`, `hello`, `world`, `hi`, `kate`,
`paul`, `good`, `morning`, `weather`, `today`, and `voice`. Each pair was
called with an interword comma, period, or ellipsis under `UU`, `UL`, `LU`, and
`LL` case masks: 1,728 calls and 3,456 paired captures total. Every call
returned raw EAX `1`, every capture contained TypeFlag=2 phone rows, and every
TypeFlag=1 row was between TypeFlag=2 rows. Commas use `Size=3200`; periods and
ellipses use `Size=14800`.

| Interword mark | UU | UL | LU | LL |
| --- | ---: | ---: | ---: | ---: |
| comma | 143/144 | 144/144 | 144/144 | 144/144 |
| period | 51/144 | 80/144 | 144/144 | 126/144 |
| ellipsis | 55/144 | 80/144 | 144/144 | 144/144 |

The only comma omission is uppercase `PAUL, HI` (`paul/hi` under `UU`); the
other three case masks emit the 3,200-sample row. Period and ellipsis presence
depend on word identity as well as case: unlike the exhaustive one-letter
matrix, upper-leading cases can emit silence, and lower-leading periods can
omit it. To preserve the complete 144-pair cross-case result, case signatures
are listed in mask order `UU`, `UL`, `LU`, `LL`; `1` means the row is present:

| Mark | Signature counts |
| --- | --- |
| comma | `1111`: 143; `0111`: 1 |
| period | `0010`: 6; `0011`: 54; `0110`: 7; `0111`: 26; `1011`: 4; `1110`: 5; `1111`: 42 |
| ellipsis | `0011`: 56; `0111`: 33; `1011`: 8; `1111`: 47 |

These exact counts cover this twelve-word set and the two-token interword
shape. They do not establish a general word-level punctuation or case rule.
Replay and validation use `run-makeinfo-word-interword-grid.sh` and
`analyze_makeinfo_word_interword_grid.py`; the captures and value-bearing API
log are named `mi-word-interword-*` and `makeinfo-word-interword-grid-api.log`
in this Stage 21 workspace.

## Whitespace around interword punctuation

The same word pairs, marks, and case masks were then tested with four spacing
layouts (6,912 calls and 13,824 paired captures):

| Layout | Input shape |
| --- | --- |
| `after_space` | `left, right` (canonical baseline) |
| `both_space` | `left , right` |
| `before_only` | `left ,right` |
| `adjacent` | `left,right` |

For periods and ellipses, replace the comma with `.` or `...`. The table gives
TypeFlag=1 presence counts per 144 ordered word pairs, in `UU`, `UL`, `LU`,
`LL` order:

| Layout | Comma | Period | Ellipsis |
| --- | --- | --- | --- |
| `after_space` | 143, 144, 144, 144 | 51, 80, 144, 126 | 55, 80, 144, 144 |
| `both_space` | 143, 144, 144, 144 | 55, 80, 144, 126 | 144, 144, 144, 144 |
| `before_only` | 143, 144, 144, 144 | 0, 0, 0, 0 | 144, 144, 144, 144 |
| `adjacent` | 143, 144, 144, 144 | 0, 0, 144, 0 | 55, 80, 144, 144 |

Commas are invariant: all 576 outputs are byte-identical across layouts in
both ASCII and binary DTT. Adjacent ellipses also byte-match the canonical
layout in all 576 cases. For ellipses, `both_space` and `before_only` also
byte-match each other in all 576 cases, so the space after the mark is inert
for these word pairs; adding a space before the mark changes the result.
Moving the ellipsis to have a preceding space
(`both_space` or `before_only`) creates an interval in every case; the newly
present rows are 3,200 samples for 89 `UU` and 64 `UL` pairs, while the
remaining intervals are 14,800 samples. In these 153 newly-present cases,
phone rows also differ from the attached layout; the other 423 cases are
byte-identical in both DTT formats. Thus this spacing change affects both
silence presence and phone selection for the same cases.

Period spacing has a different effect. With a space on both sides, four `UU`
pairs (`a/a`, `a/i`, `i/a`, `i/i`) gain a 14,800-sample row, and 24/576
outputs differ from the canonical full DTT bytes. With a space before but not
after the period, no pair emits TypeFlag=1. With no spaces around the period,
only the `LU` mask emits it, for all 144 pairs. The non-silence DTT rows also
change for many period layouts, so these results capture both silence-row
and phone-selection changes. `both_space` versus `before_only` period outputs
match in only 171/576 cases; attached versus adjacent periods match in
255/576. The complete validator checks all calls,
TypeFlag=1 placement/sizes, and byte-equivalence counts:
`analyze_makeinfo_word_spacing_grid.py`. The runner is
`run-makeinfo-word-spacing-grid.sh`; captures and the API log use the
`mi-word-space-*` and `makeinfo-word-spacing-grid-*` prefixes.

Run the spacing matrix and validator from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-word-spacing-grid.sh
python3 tools/revkit/work/stage21/analyze_makeinfo_word_spacing_grid.py
```

## ASCII whitespace around interword punctuation

A follow-up varied the separator on each side independently over the same
three marks and case masks. It used eight ordered word pairs: `a/a`,
`a/hello`, `hello/world`, `paul/hi`, `good/morning`, `hi/kate`, `a/good`, and
`weather/today`. The six values on either side were no separator, one space,
two spaces, TAB, LF, and CRLF. This is 3,456 API calls and 6,912 paired
captures (32 pair/case combinations per mark/separator cell); all calls
returned raw EAX `1`, and every ASCII capture had phone rows with any silence
row placed between them. The LF and CRLF bytes were written directly into the
MakeInfo text argument in memory.

Comma output is fully invariant for this selected set: every one of the 36
before/after separator combinations byte-matches the no-prefix, one-space
control in both ASCII and binary DTT. `PAUL, HI` under `UU` remains the only
case without a TypeFlag=1 row; the other 31 combinations emit `Size=3200`.

For periods, the following table gives TypeFlag=1 counts per 32 pair/case
combinations, written as `3200 / 14800 / absent`:

| Before `.` | After none | After space | After two spaces | After TAB | After LF | After CRLF |
| --- | --- | --- | --- | --- | --- | --- |
| none | 0 / 8 / 24 | 0 / 22 / 10 | 0 / 32 / 0 | 0 / 22 / 10 | 10 / 22 / 0 | 0 / 32 / 0 |
| space, two spaces, or TAB | 0 / 0 / 32 | 0 / 23 / 9 | 0 / 32 / 0 | 0 / 23 / 9 | 9 / 23 / 0 | 0 / 32 / 0 |
| LF or CRLF | 32 / 0 / 0 | 9 / 23 / 0 | 0 / 32 / 0 | 9 / 23 / 0 | 9 / 23 / 0 | 0 / 32 / 0 |

The before-period `space`, `two spaces`, and TAB captures byte-match each other
when followed by one space; LF and CRLF before the period also match each
other in that condition. After a period, LF and CRLF are distinguishable: with
no preceding separator, LF yields ten 3,200-sample intervals and 22
14,800-sample intervals, while CRLF yields 32 intervals of 14,800 samples.
Period phone rows change along with interval presence in many cells. These
results rule out treating every ASCII separator as interchangeable for this
tested punctuation position; they do not identify the parser's general
whitespace rule.

For ellipses, the count table has the same `3200 / 14800 / absent` notation:

| Before `...` | After none, space, TAB, or LF | After two spaces or CRLF |
| --- | --- | --- |
| none | 0 / 25 / 7 | 0 / 32 / 0 |
| any tested nonempty separator | 7 / 25 / 0 | 0 / 32 / 0 |

With one space after `...`, any tested separator before the mark adds a
3,200-sample row and changes phone rows in the seven UU/UL cases that lack an
interval in the no-separator control. Two spaces or CRLF after `...` instead
makes all 32 intervals 14,800 samples, including those seven cases. The exact
byte and phone-row comparisons are produced by
`analyze_makeinfo_ascii_whitespace_grid.py` from the value logs and captures.
Reproduce from the Stage 21 sandbox with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-ascii-whitespace-grid.sh
python3 tools/revkit/work/stage21/analyze_makeinfo_ascii_whitespace_grid.py
```

The runner supports `START_CASE` for a contiguous resume after an interrupted
capture run; its default is zero and a fresh run refuses existing outputs.
These eight pairs deliberately sample earlier exceptions and controls. The
matrix does not establish whitespace behavior for arbitrary word pairs,
punctuation placement, or text longer than two words.

### Remaining ASCII whitespace controls

Another 576 calls tested vertical tab, form feed, and bare carriage return
independently before or after each punctuation mark. Every call returned raw
EAX `1` and produced both capture formats. Before punctuation, each of those
three bytes is byte-for-byte equivalent to each tested one-space, two-space,
or TAB layout when followed by one space. After punctuation, each byte is
byte-for-byte equivalent to the no-prefix/single-space control for all 32
pair/case combinations, including its phone rows and TypeFlag=1 presence and
size. This differs from the LF and CRLF results for period and ellipsis; the
observed behavior depends on the exact bytes and their side of punctuation.
The supplemental captures and API log use `mi-word-wsclass-*` and
`makeinfo-ascii-whitespace-class-grid-api.log`; validate them with
`analyze_makeinfo_ascii_whitespace_class_grid.py`.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-ascii-whitespace-class-grid.sh
python3 tools/revkit/work/stage21/analyze_makeinfo_ascii_whitespace_class_grid.py
```

#### Vertical tab, form feed, and carriage return across all word pairs

The previous control covered only eight selected pairs. A wider sweep tested
all 144 ordered pairs from the twelve-word inventory `a`, `i`, `hello`,
`world`, `hi`, `kate`, `paul`, `good`, `morning`, `weather`, `today`, and
`voice`. It crossed three punctuation marks, four upper/lowercase masks, both
separator positions, and vertical tab/form feed/bare carriage return: 10,368
calls and 20,736 paired captures. Every call returned raw EAX `1`; the
validator checked the complete coordinate set and each ASCII/BIN pair.
The zero-byte output from the timed-out attempt at the first resumed case is
retained as `mi-ws12-9883.bin.dtt.interrupted`; that case was rerun, and this
artifact is excluded from the paired-capture count.

For each position and punctuation mark, VT, FF, and CR produced byte-identical
ASCII and binary output in all 576 word/case combinations. Against the
canonical one-space-after-punctuation captures, all three controls matched in
every case after punctuation and before commas. Before periods, 552/576
matched; before ellipses, 423/576 matched. The observed TypeFlag=1 counts
(`3200 / 14800 / absent`, per 576 cases) were also identical across VT, FF,
and CR: before punctuation, comma `575 / 0 / 1`, period `0 / 405 / 171`, and
ellipsis `153 / 423 / 0`; after punctuation, comma `575 / 0 / 1`, period
`0 / 401 / 175`, and ellipsis `0 / 423 / 153`. Thus these three bytes behave
as one output-equivalence class in this tested two-word inventory, while
before-period and before-ellipsis outputs can differ from the canonical
one-space layout. This is a bounded output result, not a recovered general
tokenizer rule.

The reproducible runner and validator are
`run-makeinfo-ascii-whitespace-class-grid-all-pairs.sh` and
`analyze_makeinfo_ascii_whitespace_class_grid_all_pairs.py`. The runner accepts
`START_CASE` to resume a contiguous interrupted capture set; its default is
zero, and fresh runs refuse existing outputs.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-ascii-whitespace-class-grid-all-pairs.sh
python3 tools/revkit/work/stage21/analyze_makeinfo_ascii_whitespace_class_grid_all_pairs.py
```

### Six-class ASCII whitespace matrix across all word pairs

The VT/FF/CR sweep does not cover space, repeated space, TAB, LF, or CRLF in
both positions at once. This follow-up crossed the six layouts `none`, one
space, two spaces, TAB, LF, and CRLF before and after each punctuation mark,
for all 144 ordered pairs in the same twelve-word inventory and all four case
masks. It is 62,208 calls (144 pairs × 3 marks × 4 masks × 36 layout pairs)
and 124,416 paired ASCII/BIN captures. Every raw return was `1`; the analyzer
checked every coordinate, output pair, phone-row presence, optional internal
TypeFlag=1 row, and interval size. The two sandbox fixture files were restored
byte-for-byte after the run.

Counts below are `3200 / 14800 / absent` TypeFlag=1 rows per 576
word/case combinations. Commas are invariant: all 36 layout pairs are
byte-identical in both DTT formats to each other and the canonical
no-prefix/one-space output, with `575 / 0 / 1` rows in every layout. The sole
missing interval remains `PAUL, HI` under `UU`.

| Before `.` | After none | After space | After two spaces | After TAB | After LF | After CRLF |
| --- | --- | --- | --- | --- | --- | --- |
| none | 0 / 144 / 432 | 0 / 401 / 175 | 0 / 576 / 0 | 0 / 401 / 175 | 175 / 401 / 0 | 0 / 576 / 0 |
| space, two spaces, or TAB | 0 / 0 / 576 | 0 / 405 / 171 | 0 / 576 / 0 | 0 / 405 / 171 | 171 / 405 / 0 | 0 / 576 / 0 |
| LF or CRLF | 576 / 0 / 0 | 171 / 405 / 0 | 0 / 576 / 0 | 171 / 405 / 0 | 171 / 405 / 0 | 0 / 576 / 0 |

The period output depends on the exact separator bytes and position, and the
phone rows also vary. Against the canonical no-prefix/one-space layout,
before-separator `none` matches all 576 cases when followed by one space or
TAB; a preceding space/two spaces/TAB matches 552 complete outputs when
followed by either. A preceding LF/CRLF with one space or TAB after the mark
matches 401. Post-period two spaces or CRLF yields 576 TypeFlag=1 rows of
14,800 samples for every tested preceding layout. These are output
equivalences on this inventory, not a general whitespace-tokenization rule.

For ellipses, the count matrix reduces to three observed groups:

| Before `...` | After none, space, or TAB | After LF | After two spaces or CRLF |
| --- | --- | --- | --- |
| none | 0 / 423 / 153 | 153 / 423 / 0 | 0 / 576 / 0 |
| any nonempty layout | 153 / 423 / 0 | 153 / 423 / 0 | 0 / 576 / 0 |

With no prefix and no separator, one space, or TAB after the ellipsis, 153
cases have no interval; any nonempty preceding layout changes those cases to
3,200-sample rows. With LF after the mark, all preceding layouts yield
153/423 rows of 3,200/14,800 samples. Two spaces or CRLF after the mark puts
all cases in the 14,800-sample group regardless of the preceding layout. Full
ASCII+BIN comparisons show the no-prefix layouts with no separator, one
space, or TAB after `...` are identical; with two spaces or CRLF after it,
every preceding layout is identical. The remaining layouts form a third
byte-equivalence group. This captures the exact two-word output partition
without assigning a tokenizer meaning to it.

The capture-set analyzer also prints each of the 108 mark/layout silence
counts and exact/row-only matches against the canonical control. Reproduce a
fresh run in 4,096-case chunks (the runner refuses to overwrite captures and
the runtime has a 30-minute limit): set `START_CASE` to the next contiguous
case index and `CASE_LIMIT=4096`, advancing by the number of completed calls;
the final range is clipped automatically at 62,208. Validate the complete
capture set with:

```sh
START_CASE=0 CASE_LIMIT=4096 docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -e START_CASE -e CASE_LIMIT -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-ascii-whitespace-all-pairs.sh
python3 tools/revkit/work/stage21/analyze_makeinfo_ascii_whitespace_all_pairs.py
```

The runner and analyzer are `run-makeinfo-ascii-whitespace-all-pairs.sh` and
`analyze_makeinfo_ascii_whitespace_all_pairs.py` in this workspace.

## MakeInfo scalar pause argument

The seventh scalar argument was varied while speaker 1 and pitch, speed, and
volume stayed at `-1`. The first matrix ran 108 calls over nine texts and 12
pause values (`-1`, `0`, `1`, `119`, `120`, `121`, `249`, `250`, `251`,
`1000`, `65534`, and `65535`). A 27-call edge follow-up tested `-2`, `-100`,
`-2147483648`, `65536`, `65537`, `100000`, and `2147483647` on three
silence-producing contexts. Every call returned raw EAX `1` and emitted both
DTT files.

Where the tested punctuation/case context creates a TypeFlag=1 interval and
its per-unit pause is unset, positive scalar values set its size at 16
samples per unit: 1 gives 16 samples, 120 gives 1,920, 250 gives 4,000,
1,000 gives 16,000, and 65,534 gives 1,048,544. This matches milliseconds
converted at the established 16 kHz sample rate. Zero omits the TypeFlag=1
row. Negative values tested, including `INT_MIN`, preserve the default 14,800
samples (925 ms). Values from 65,535 through `INT_MAX` saturate at 65,535
units, producing 1,048,560 samples. In every variable-duration context, all
other ASCII rows remain byte-identical after removing the TypeFlag=1 row.

The fixed comma case `A, B` retains its 3,200-sample row for every tested
pause value, including zero and values above the cap. Conversely, the tested
period/ellipsis contexts that do not create a silence row remain without one
at every pause value. This shows the scalar pause controls an already-selected
unset silence interval; it does not itself select the punctuation/case rule
that inserts that interval. The capture paths and raw runtime logs are
`makeinfo-pause-silence-*` and `makeinfo-pause-edges-*` under this stage's
workspace. Runners are `run-makeinfo-pause-silence-probe.sh` and
`run-makeinfo-pause-edge-probe.sh`.

## Natural synthesis-context user-dictionary unload guard

`run-natural-userdict-inuse-unload.sh` pauses the sample at its normal
`VT_TextToFile_ENG` entry, loads `hello,HH,P` at dictionary index 0, and
temporarily sets Paul slot 1's dictionary gate to 1. The supplied Paul record
leaves that gate off because its database-size predicate fails. The runner
then rewrites the pending API call to synthesize `hello` as selector-4 WAVE
for slot 1 using dictionary index 0.

At breakpoint `0x100260C2`, reached after `FUN_10025FC0` assigns the selected
dictionary pointer to its allocated context, the actual synthesis context in
slot 1's first reference cell (`0x100A147C`) had a `+0x1312C0` dictionary
field matching the loaded dictionary pointer. `VT_UNLOAD_UserDict_ENG(0)`
returned `-3` there. The utterance then returned `1`; after the context was
released, idle unload returned `1`. The gate was restored to 0, and a further
synthesis returned `1`.

The same-process control and post-unload PCM match byte-for-byte: 8,984 mono
16 kHz PCM16 frames with WAVE SHA-256
`613ff17ff3d7b8c9a411ef771ea1d02a051f38fcb918874e95dde30db39c0f36`. The
active-dictionary output is 1,448 frames with SHA-256
`491b29d0c3deb1d95770f8bdfe6fbcdfce5ace03f8c15dbfb9c57cf95b601cf1`.
`analyze_natural_userdict_inuse_unload.py` verifies call statuses, pointer
identity, WAVE format, hashes, and control/restored equality. This covers one
real context phase on Paul slot 1 with a debugger-forced gate; other slots,
thread contention, and other context phases remain open.
The authoritative trace is `natural-userdict-inuse-unload-v3-api.log` with
`trace-natural-userdict-inuse-unload-v3.gdb`; the three `*-v3.wav` captures
are the control, active, and restored outputs. Earlier watchpoint experiments
in this filename family did not reach the guard and are superseded by v3.

Reproduce from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-natural-userdict-inuse-unload.sh
```

## Naturally license-enabled James user-dictionary effects

`run-james-licensed-userdict-effect.sh` attaches to the sample's normal Paul
process, loads James into speaker slot 4 with the supplied verification record
by file path, and probes the ordinary user-dictionary and file-synthesis
exports on that slot. The extended model loader returned AX `0`; direct slot
state reads showed license gate `1` and dictionary capacity `6`. A format-4
`hello` synthesis using dictionary index 0 produced the 7,798-frame control.
The runner then loaded `james-userdict-p.csv` (`hello,HH,P`) and
`james-userdict-a.csv` (`hello,world,A`) in separate cycles. Each load,
synthesis, and unload returned low AX `1`.

The P row produced 1,150 frames and the A row 9,912; both PCM payloads differ
from control. After unloading both rows, the same synthesis reproduced control
PCM byte-for-byte. This establishes a natural license-enabled dictionary
effect for these two rows and this James voice/record. It does not identify
the words heard in the output or generalize to other rows, voices, and license
states. The raw API log is `james-licensed-userdict-effect-api.log`; WAVE
captures are `james-userdict-{control,p,a,restored}.wav`.

Reproduce from the repository root with Stage 8 and Stage 17 Compose files,
mounting `data-james` read-only:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm -w /work/stage5 \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage21/run-james-licensed-userdict-effect.sh
```

Validate the capture with `analyze_james_licensed_userdict_effect.py`.

## `[CI]` suffix in `P` user-dictionary rows

The source text `HELLO` was synthesized against three cases: no dictionary,
plain `hello,HH,P`, and each of `hello,HH[CI],P` / `hello,HH [CI],P`. The
plain `P` row left output byte-identical to the no-dictionary control at 7,798
frames. Both marker rows loaded and unloaded with status `1`, and both changed
the output to the same 1,150-frame PCM produced for lowercase `hello` by the
plain `P` row. For these rows, `[CI]` makes source matching case-insensitive;
the tested marker spacing does not change the result or target-phone output.
This conclusion is bounded to the tested ASCII source and `P` rows. It does
not map other ASCII case masks, non-ASCII case folding, other word lengths, or
marker interaction with other target tokens. The decompiler cross-reference
and remaining scope are in the detailed API behavior report.

The adjacent form was also tested with mixed-case input `HeLLo`: the plain
`P` row remained equal to control, while the `[CI]` row produced the same
1,150-frame PCM as the lowercase match. Capture:
`james-licensed-userdict-ci-source-case-mixed-v1-api.log` and
`james-userdict-ci-case-mixed-{control,p,a,restored}.wav`.

Marker capitalization was checked with `hello,HH[ci],P`. It loaded and
unloaded with status `1`; lowercase `[ci]` produced the same 1,150-frame
uppercase `HELLO` result as `[CI]`, and lowercase `hello` matched the plain
`P` row's output. Thus uppercase and lowercase marker spellings enable the
tested case-insensitive match. Both mixed masks `[cI]` and `[Ci]` also loaded
and unloaded with status `1` and produced the same 1,150-frame match. The
direct checker and loader therefore agree across all four marker case masks.
Captures:
`james-licensed-userdict-ci-source-case-lower-tag-v1-api.log` and
`james-licensed-userdict-ci-source-case-lower-tag-lower-source-v1-api.log`.
Mixed-mask captures: `james-licensed-userdict-ci-tag-cI-case-v2-api.log` and
`james-licensed-userdict-ci-tag-ci-case-v3-api.log`.
Replay with the same Stage 8/Stage 17 Compose files and read-only James mount,
using `run-james-licensed-userdict-ci-source-case-lower-tag.sh` and
`run-james-licensed-userdict-ci-lower-tag-lower-source.sh`,
`run-james-userdict-ci-tag-cI-case-v2.sh`, and
`run-james-userdict-ci-tag-Ci-case-v3.sh`. Check the PCM comparisons and call
statuses with:

```sh
python3 tools/revkit/work/stage21/analyze_james_userdict_ci_case.py
```

Replay the adjacent and spaced cases with the documented read-only James data
mount:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm -w /work/stage5 \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage21/run-james-licensed-userdict-ci-source-case.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm -w /work/stage5 \
  -v "$PWD/data-james:/work/data-james:ro" runtime \
  /bin/bash /work/stage21/run-james-licensed-userdict-ci-source-case-spaced.sh
```

Inputs: `james-userdict-p-ci-adjacent.csv` and
`james-userdict-a-ci-spaced.csv`. Each capture includes the plain-P control,
marker result, and post-unload restoration. See
`james-licensed-userdict-ci-source-case-v1-v3-api.log` and
`james-licensed-userdict-ci-source-case-spaced-v1-api.log`, plus their
`james-userdict-ci-case-*.wav` outputs.

## File-API stored pause setters

The Stage 16 configuration probe changed the stored sentence and comma pause
values while synthesizing only `Hello world.`, which selects neither tested
boundary. Stage 21 therefore tested the setters against interword period and
comma contexts. Both stored settings add zero-valued PCM16 frames at 16
samples per millisecond when their punctuation context selects a pause. For
the fixed `Hello. World.` / `Hello, world.` fixture, the result is
`17,762 + 16 × pause` mono 16 kHz frames, with the inserted interval starting
at PCM byte offset 18,006. Tested ranges reached 65,535 ms (1,066,322 total
frames). The matching punctuation-free controls stay at 11,803 frames and
unchanged PCM. Case variants of the period form and comma spacing with a
space, no space, double space, or TAB also produce identical PCM in these
tested contexts.

Sentence pause: 12 settings × 4 texts = 48 calls. Comma pause: 16 settings ×
5 texts = 80 calls. Every synthesis returned 1 and every getter matched its
setting. The full details are in the [Lead 6 API behavior
report](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md#sentence-pause-synthesis-effect)
and its [comma-pause section](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md#comma-pause-synthesis-effect).
Runners and analyzers are `run-sentence-pause-context-grid.sh`,
`analyze_sentence_pause_context_grid.py`, `run-comma-pause-context-grid.sh`,
and `analyze_comma_pause_context_grid.py`. These matrices do not establish
other speaker behavior, every intermediate value,
VTML pause handling, or other synthesis formats.

The separate 36-call `run-pause-precedence-grid.sh` probe varies stored
sentence or comma pause (0/925), the per-call file pause argument (-1/0/250),
and period, comma, or no-punctuation text. `analyze_pause_precedence_grid.py`
checks all returns, WAVE formats, frame counts, and exact PCM insertions. A
nonnegative per-call value overrides stored sentence pause for the period
interval; -1 falls back to the stored sentence value. The comma interval uses
the stored comma setting for every tested per-call value. Captures and logs
use the `pause-precedence-*` prefix. See the [precedence results in the Lead
6 report](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md#stored-pause-versus-per-call-pause-argument).

`run-negative-pause-setter-grid.sh` sets baseline slot-1 configuration values,
then tests `-1`, `-2`, and `INT_MIN` independently in each pitch/speed/volume/
sentence-pause field and in the comma-pause setter. The 15 getter observations
all retain the full baseline tuple. `analyze_negative_pause_setter_grid.py`
validates the API statuses and state; the log uses the
`negative-pause-setter-grid` prefix.

The `run-text-file-path-error-grid.sh` probe calls `VT_TextToFile_ENG` with a
new absolute Z-drive filename, the existing Stage 21 sandbox directory, and
a filename under a deliberately absent parent. The valid call returned 1 and
produced a selector-4 WAVE byte-identical to the control; the two open failures
returned `-6`. The analyzer checks the return matrix, WAVE format, frame count,
control equality, and absence of the missing-parent path. Evidence uses the
`text-file-path-error-grid` prefix.

## `VT_SetCommaPause_ENG` synthesis effect

A single-process Stage 21 sweep set slot 1 to 16 values from 0 through 65,535
and synthesized four comma contexts (`Hello, world.`, attached comma, double
space, and TAB) plus `Hello world.` with the file API's other options at
`-1`. All 80 synthesis calls returned 1; each getter returned the requested
value. For comma contexts, output is mono 16 kHz PCM16 and frame count is
`17,762 + 16 × pause`. Against pause zero, the only changed bytes are
`16 × pause` zero-valued samples inserted at PCM byte offset 18,006. The four
comma whitespace layouts have identical PCM for each setting. The no-comma
control stays at 11,803 frames and identical PCM across all settings.

The runner backs up and restores the Stage 21 sandbox fixtures and refuses to
overwrite existing captures. Reproduce and validate from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-comma-pause-context-grid.sh
python3 tools/revkit/work/stage21/analyze_comma_pause_context_grid.py
```

Run a punctuation batch in the Stage 8 Docker/Wine environment:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-interword-punctuation-matrix.sh comma
```

Replace `comma` with `period` or `ellipsis`. The combined validator checks
all three batches:

```sh
python3 tools/revkit/work/stage21/analyze_makeinfo_interword_punctuation_matrix.py
```

Recheck the aggregate and all selected-bank DAT offsets from the repository
root:

```sh
python3 tools/revkit/scripts/analyze_makeinfo_alphabet_matrix.py
python3 tools/revkit/work/stage21/analyze_makeinfo_four_token_context.py
python3 tools/revkit/work/stage21/analyze_makeinfo_a_prefix_pair_matrix.py
python3 tools/revkit/work/stage21/check_makeinfo_letter_fullmask_batch.py b1
python3 tools/revkit/work/stage21/check_makeinfo_letter_fullmask_batch.py b2
python3 tools/revkit/work/stage21/check_makeinfo_letter_fullmask_batch.py b3
python3 tools/revkit/work/stage21/analyze_makeinfo_a_prefix_fullmask_pair_matrix.py
```

The runtime runners use the existing Stage 8 Docker/Wine environment:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-alphabet-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-lowercase-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-letter-context-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-alphabet-punctuation-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower /bin/bash /work/stage21/run-makeinfo-alphabet-punctuation-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-letter-pair-question-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower /bin/bash /work/stage21/run-makeinfo-letter-pair-question-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=upper-lower /bin/bash /work/stage21/run-makeinfo-letter-pair-question-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower-upper /bin/bash /work/stage21/run-makeinfo-letter-pair-question-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-letter-pair-punctuation-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower /bin/bash /work/stage21/run-makeinfo-letter-pair-punctuation-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=upper-lower /bin/bash /work/stage21/run-makeinfo-letter-pair-punctuation-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower-upper /bin/bash /work/stage21/run-makeinfo-letter-pair-punctuation-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-letter-triple-a-context-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-context-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=upper-lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-context-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower-upper /bin/bash /work/stage21/run-makeinfo-letter-triple-a-context-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=upper-lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower-upper /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower-lower-upper /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=lower-upper-lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=upper-lower-upper /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'LETTER_CASE=upper-upper-lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'POSITIONS=leading LETTER_CASE=lower-lower-upper /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'POSITIONS=leading LETTER_CASE=lower-upper-lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'POSITIONS=leading LETTER_CASE=upper-lower-upper /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash -c \
  'POSITIONS=leading LETTER_CASE=upper-upper-lower /bin/bash /work/stage21/run-makeinfo-letter-triple-a-position-matrix.sh'
```

The targeted four-token probe is reproducible with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-four-token-context.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-fourth-token-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-third-token-sweep.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-a-prefix-pair-matrix.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-a-prefix-fullmask-pair-matrix.sh b1
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-a-prefix-fullmask-pair-matrix.sh b2
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime \
  /bin/bash /work/stage21/run-makeinfo-letter-a-prefix-fullmask-pair-matrix.sh b3
```

The four-token follow-up logs and traces are `makeinfo-letter-four-token-context-*`,
`makeinfo-letter-fourth-token-sweep-*`, `makeinfo-letter-third-token-sweep-*`,
`makeinfo-letter-a-prefix-pair-matrix-*`, and
`makeinfo-letter-a-prefix-fullmask-b*-*`. Their paired outputs use the
`mi-letter-four-context-*`, `mi-letter-fourth-sweep-*`,
`mi-letter-third-sweep-*`, `mi-letter-a-prefix-pair-*`, and
`mi-letter-a-prefix-fullmask-b*-pair-*` prefixes. A separate oversized trace
was stopped after 9,463 successful calls and is retained under the distinct
`mi-letter-a-prefix-fullmask-pair-*` prefix; exclude it from the validated
21,632-call matrix. These four-token capture sets and analyzers are separate
from the exhaustive three-token aggregate.

Evidence files are `makeinfo-alphabet-matrix-v2-api.log`,
`trace-makeinfo-alphabet-matrix-v2.gdb`, the paired
`mi-alphabet-v2-*.asc.dtt`/`.bin.dtt` captures, and
`makeinfo-alphabet-matrix-crosscheck.txt`. Lowercase-call evidence is in
`makeinfo-lowercase-matrix-api.log`,
`trace-makeinfo-lowercase-matrix.gdb`, and the paired
`mi-alphabet-lower-*.asc.dtt`/`.bin.dtt` captures. Context-matrix evidence is
in `makeinfo-letter-context-matrix-api.log`,
`trace-makeinfo-letter-context-matrix.gdb`, and the paired
`mi-alphabet-context-*.asc.dtt`/`.bin.dtt` captures. The 104-case terminal
uppercase punctuation evidence is `makeinfo-alphabet-punctuation-matrix-api.log`,
`trace-makeinfo-alphabet-punctuation-matrix.gdb`, and the paired
`mi-alphabet-punct-*.asc.dtt`/`.bin.dtt` captures. The cross-check summary also
compares every punctuated TypeFlag=2 record field-for-field against its
standalone counterpart and across cases, and counts TypeFlag=1 rows.
Lowercase evidence is in `makeinfo-lowercase-punctuation-matrix-api.log`,
`trace-makeinfo-lowercase-punctuation-matrix.gdb`, and
`mi-alphabet-punct-lower-*.asc.dtt`/`.bin.dtt`. The ordered-pair evidence is
`makeinfo-letter-pair-question-matrix-api.log`,
`trace-makeinfo-letter-pair-question-matrix.gdb`, and
`mi-letter-pair-*.asc.dtt`/`.bin.dtt`. Lowercase captures use the
`mi-letter-pair-lower-` prefix, `makeinfo-lowercase-letter-pair-question-matrix-api.log`,
and `trace-makeinfo-lowercase-letter-pair-question-matrix.gdb`. Mixed-case
evidence is in `makeinfo-mixed-upper-lower-letter-pair-question-matrix-api.log`,
`trace-makeinfo-mixed-upper-lower-letter-pair-question-matrix.gdb`,
`makeinfo-mixed-lower-upper-letter-pair-question-matrix-api.log`, and
`trace-makeinfo-mixed-lower-upper-letter-pair-question-matrix.gdb`; captures use
the `mi-letter-pair-mixed-upper-lower-*` and
`mi-letter-pair-mixed-lower-upper-*` prefixes. The cross-check summary lists
every pair's selected bank/unit rows, the plain-to-question bank-set
transitions, and all case comparisons. Pair-punctuation evidence is in
`makeinfo-letter-punctuation-{upper,lower,mixed-upper-lower,mixed-lower-upper}-api.log`
and corresponding `trace-makeinfo-letter-punctuation-*.gdb` files; the ASCII
and binary captures use the `mi-letter-pair-punct-*` prefix. The analyzer
compares these rows with same-case plain pairs and with the uppercase
punctuation outputs. Three-token evidence is in
`makeinfo-letter-triple-a-context-{upper,lower,mixed-upper-lower,mixed-lower-upper}-api.log`
and corresponding `trace-makeinfo-letter-triple-a-context-*.gdb` files;
captures use `mi-letter-triple-a-*`. The A-position extension is in
`makeinfo-letter-triple-a-position-*-api.log` and
`trace-makeinfo-letter-triple-a-position-*.gdb`; it captures every
`B A C` and `B C A` letter continuation, plain and question-mark forms, for all
eight per-token case masks. The four added A-leading masks are captured by the
same runner with `POSITIONS=leading`; their logs and traces use the
`makeinfo-letter-triple-a-leading-{mask}-api.log` and
`trace-makeinfo-letter-triple-a-leading-{mask}.gdb` names, with output captures
using `mi-letter-triple-a-{mask}-*`. The analyzer compares all three A
positions' plain/question rows and capitalization variants. The initial
`LUL`/`ULU` A-leading attempts used an incorrect middle-token case mapping;
their captures and logs remain under `invalid-case-map/` and are excluded from
the analyzer's top-level capture glob. Corrected runs replace them in the
aggregate. The earlier runner attempt is retained as a setup failure: GDB rejected generated variable names containing
hyphens before any MakeInfo call, so its
`makeinfo-alphabet-matrix-api.log` contains no result evidence. The valid
capture runs in this session were initially launched without the working
directory override; after GDB detached, the sample wrote its WAV to Stage 5.
That file was restored from the zero-byte Stage 16 pre-probe backup and
verified against `HEAD` (both SHA-256
`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`). Future
runs are protected by both the documented `-w` override and each runner's
working-directory check.

## Progress-buffer long-text notifications

The target-flow probe constructs a 55-byte utterance by repeating the sample
text four times, then calls `VT_TextToPcmBuffer_ProgressBar_ENG` with a
60,000-byte guarded buffer, thread 0, speaker 1, scalar options `-1`, null
window, and message `0x8005`. It observes three 60,000-byte output chunks
(status 0), followed by 3,224 bytes (status 1), for 183,224 bytes total. Four
calls reach the DLL's `USER32.dll!PostMessageA` import at `0x1002151b`:

| Post | HWND | Message | wParam | lParam | API call |
| --- | ---: | ---: | ---: | ---: | --- |
| 1 | `0` | `0x8005` | `0xa7` | `0x18` | Initial call |
| 2 | `0` | `0x8005` | `0xa7` | `0x26` | Initial call |
| 3 | `0` | `0x8005` | `0xa7` | `0x34` | Poll 1 |
| 4 | `0` | `0x8005` | `0xa7` | `0x34` | Poll 2 |

Poll 3 completes without a post. Static operands show the supplied HWND and
message pass through, wParam comes from the first dword of the current
per-speaker progress state, and lParam comes from the selected SyncInfo row's
`+0x10` field. The earlier SyncInfo producer/consumer analysis establishes
that row `+0x10` is the inclusive source-buffer end byte offset. The progress
routine chooses row `cursor-1`, or row 599 when the cursor is zero. The host's
use of that source-span endpoint remains unknown.
The dedicated lParam cross-check read the current cursor and row fields at
each post. Cursors 5, 8, 11, and 11 selected rows 4, 7, 10, and 10; their
inclusive source-byte spans were `0x14..0x18`, `0x22..0x26`, `0x30..0x34`,
and `0x30..0x34`. All four posted lParams equaled the selected row's end field.
The repeated `0x34` therefore came from the same row remaining selected on
polls 1 and 2 in this run; why that cadence occurs more generally remains
open. Cursor-zero wrap to row 599 is statically established but was not
exercised by this trace.
At `0x10021513`, a corrected live trace saved the per-speaker context pointer
from EAX and read its first DWORD; at `0x10021519`, it compared that saved value
with the wParam already pushed for `PostMessageA`. All four posts matched in
each checked run. A watchpoint started immediately after context allocation
observed the first DWORD begin at zero and change once, from `0` to `0xc4`,
inside `FUN_10026ab0` at the store on `0x10026ad4`. The disassembly computes
that DWORD as the byte length of its third string argument (`strlen`), which
is the string returned by `FUN_1001c990` in this path. A separate capture
measured that helper-produced string as 167 bytes: 111 leading periods, a
space, and the repeated 55-byte utterance. Its `0xa7` length exactly matched
all four posted wParam values. Therefore wParam carries the byte length of the
working text string, including its leading filler; it is not a progress
percentage or a context identifier. Across completed runs, values `0x59`,
`0x72`, `0xa4`, `0xa7`, and `0xc4` each stayed constant within the utterance.
Subtracting the 55-byte source and one separating space leaves period-prefix
lengths 33, 58, 108, 111, and 140, all present in the independently measured
flag-7 filler-length set. This correspondence explains the observed variation
as a changing filler length, although whether both paths share the same
generator has not been proven. The lParam sequence repeated across these
runs, but duplicate `0x34` notification behavior has not been characterized
for other inputs. A separate trace captured BOOL `1` at
the first PostMessageA return, disabled that return breakpoint, and still saw
the inferior exit before notification two returned. That partial result does
not establish the remaining BOOL values. Its artifacts are
`trace-buffer-progress-notification-return-partial.gdb`,
`buffer-progress-notification-return-partial-api.log`, and
`run-buffer-progress-notification-return-partial.sh`. Since the captured
HWND is null, the stable trace does not verify a host-window handler or
message retrieval. The trace deliberately terminates
the inferior after collecting the final return, because execution was entered
through a synthetic target stack; GDB prints an internal shutdown diagnostic
after all the observations.

Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-buffer-progress-notifications-target.sh
```

The original trace and capture are `trace-buffer-progress-notifications-target.gdb`
and `buffer-progress-notifications-target-api.log`. The source cross-check is
`trace-buffer-progress-wparam-source.gdb` and
`buffer-progress-wparam-source-v3-api.log`, reproduced with
`run-buffer-progress-wparam-source.sh`. The initial source-probe attempt used
the wrong frame argument as a pointer and failed before producing notification
evidence (`buffer-progress-wparam-source-api.log`); v2 reached the first two
posts but had its polling breakpoint disabled after breakpoint insertion
(`buffer-progress-wparam-source-v2-api.log`). Both captures are retained as
setup failures and excluded from the observations above.

The context watchpoint and helper-string length captures use separate traces
and runners:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-buffer-progress-wparam-watch.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-buffer-progress-wparam-value.sh
```

The watchpoint is in `trace-buffer-progress-wparam-watch.gdb` with output
`buffer-progress-wparam-watch-api.log`. The helper-string capture is in
`trace-buffer-progress-wparam-value.gdb` with output
`buffer-progress-wparam-value-api.log`. The lParam cross-check is in
`trace-buffer-progress-syncinfo-lparam.gdb` with output
`buffer-progress-syncinfo-lparam-api.log`, reproduced by
`run-buffer-progress-syncinfo-lparam.sh`.

The live trace and static path establish wParam's meaning as the byte length
of the helper-produced working text. Context construction allocates a
`0x1312e0`-byte object; its first DWORD is then explicitly assigned by
`FUN_10026ab0` from the helper string length. The host handler's use of the
byte length and source-span endpoint remains unknown.

## User-dictionary source-normalizer output

The public `VT_CheckUserDict_SourceNorm_ENG` wrapper calls the private helper
with a 52-byte stack buffer but discards both its helper return and the
normalized bytes. The target-flow trace breaks immediately after that helper
call and records raw helper EAX plus all 52 scratch bytes for 17 cases: empty
and all-trimmed;
plain, punctuated, repeated-space, mixed-case, contraction, and UTF-8 text;
leading/trailing TAB/CR/LF; lengths 49, 50, and 51; and four byte pairs.

| Input | Helper EAX | Scratch output |
| --- | ---: | --- |
| Empty | `-1` | empty string |
| Space/TAB/CR/LF only | `-1` | empty string |
| `  Hello   world!  ` | `14` | `Hello   world!` |
| `\t\r\nHello\r\n\t` | `5` | `Hello` |
| `MiXeD_case-123` | `14` | unchanged |
| `can't stop` | `10` | unchanged |
| `caf` + bytes `c3 a9` | `5` | bytes preserved |
| 49 `A` bytes | `49` | all bytes copied |
| 50 or 51 `A` bytes | `-5` | first output byte cleared |
| `a1 a1`, `ae a1`, `fd fe` | `-3` | first byte remains NUL |
| `a1 a0` | `2` | both bytes copied |

Static code skips leading and trims trailing space, TAB, LF, and CR. The
successful byte-count limit is 49: at length 50 it returns `-5` and clears the
first output byte. It copies punctuation, case, internal spacing, and tested
UTF-8 bytes unchanged. At `0x1005f2e0`, the private helper checks whether two
bytes are present, then returns `-3` for `[a1-ad][a1-fe]`, `ae[a1-c2]`, or
`fd fe`; these byte exclusions' meaning is unknown. The public wrapper does
not expose the result or scratch; the trace is an instrumentation result from
the helper boundary. Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage21/sandbox/stage5 runtime /bin/bash \
  /work/stage21/run-userdict-source-normalizer-output.sh
```

The trace and capture are `trace-userdict-source-normalizer-output.gdb` and
`userdict-source-normalizer-output-api.log`. As with other target-flow traces,
the inferior is terminated after collection; GDB emits an internal shutdown
diagnostic after all matrix observations.

## Parenthesis-count setting consumer

`VT_SetParenthesisCharNumber_ENG` stores a nonnegative value in the loaded
Paul global configuration block at `DAT_100a0460+0x2041c`. A hardware access
watchpoint on the parser-context copy captured the write in `FUN_1003e470`,
then four reads in `FUN_100544f0` while parsing the parenthesized interiors
`tea`, `green tea`, `freshly brewed tea`, and `very fresh tea`. At the read,
the parser compared the setting against a helper result of 3, 9, 18, and 14,
respectively. The helper scans bytes to the first `)` or `]` and returns the
byte offset, so these ASCII spans happen to equal their character counts. The
static pseudocode takes the marked-record branch when the setting is nonzero
and no greater than that result. This supports a minimum byte-span threshold
for the tested parser path; it does not establish Unicode character counting,
the user-facing meaning of the record flag, or the square-bracket behavior.

Format-4 audio comparisons support the boundary: setting 0 leaves `(a)` at
the baseline output, setting 1 changes it; for `(tea)`, settings 1–3 produce
one identical changed output and settings 0 and 4–6 produce the baseline;
for `(book)`, settings 1 and 4 produce the same changed output; and for
`(green tea)`, settings 1 and 4 also produce the same changed output. These
fixtures show the threshold at byte spans of 1, 3, and 4. They do not
identify the marked record's exact synthesis consequence.

Run the selected-input sweep with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-parenthesis-number-effect.sh single-length3 \
  "0 1 2 3 4 5 6"
```

The watchpoint trace uses the broader `number` fixture:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-parenthesis-context-consumer.sh
```

The relevant artifacts are `trace-parenthesis-number-effect.gdb.in`, the
`parenthesis-single-length3-*-api.log` and `.wav` files,
`trace-parenthesis-context-consumer.gdb`,
`parenthesis-context-consumer-api.log`, and
`reports/stage25-parenthesis-parser-consumer.c` plus
`reports/stage25-parenthesis-length-helper.c`.

The record follow-up used the exact-length-3 fixture with setting 1. It
captured the DWORD slot written to `1`, then watched downstream access to that
slot. The runtime trace hit `FUN_1000d190` at `0x1000d208`, `FUN_1000ea20` at
`0x1000eb5c`, `FUN_10016c90` at `0x10016cee`, and context cleanup at
`FUN_1003e210` (`0x1003e238`). The static pseudocode for `FUN_1000d190` shows
the leading record DWORD entering a switch whose explicit cases are 2, 3, 4,
5, 11, and 12; value 1 follows the default arm. This proves the record is
consumed and copied downstream but does not identify its synthesis meaning.
Reproduce with `run-parenthesis-record-consumer.sh`; artifacts are
`trace-parenthesis-record-consumer.gdb`,
`parenthesis-record-consumer-api.log`, and
`reports/stage25-parenthesis-record-consumers.c`.

## CSV serializer single-byte domain

`VT_CsvParser_MakeCsv_ENG` was called once for each value `0x00`–`0xff` as a
single-field C string, with output capacity 8. Every nonzero byte returned
raw low AX 1 and appeared unchanged between CSV quote bytes, except `0x22`,
which became four quote bytes. The empty C string (`0x00`) returned raw low
AX -1 after writing only the opening quote; the `0x31` fill remained through
the final advertised NUL. Prefix and post-capacity guards stayed intact for
all cases. Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-makecsv-byte-domain-v1.sh matrix5
```

The capture is
[`csv-makecsv-byte-domain-matrix4-api.log`](./csv-makecsv-byte-domain-matrix4-api.log),
generated by [`run-csv-makecsv-byte-domain-v1.sh`](./run-csv-makecsv-byte-domain-v1.sh)
from [`trace-csv-makecsv-byte-domain-v1.gdb`](./trace-csv-makecsv-byte-domain-v1.gdb).

A selected 1,280-case pair matrix then tested every high-bit first byte with
NUL, `0x01`, quote, comma, `0x7f`, `0x80`, and `0xff`, plus all high-bit
second bytes after `A`, quote, and comma. A high-bit first byte consumed a
following non-NUL byte as a pair, leaving a quote/comma in the second position
literal. Reproduce with a unique run ID:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-makecsv-byte-pairs-v1.sh classes2
```

The capture is
[`csv-makecsv-byte-pairs-classes1-api.log`](./csv-makecsv-byte-pairs-classes1-api.log),
generated by [`run-csv-makecsv-byte-pairs-v1.sh`](./run-csv-makecsv-byte-pairs-v1.sh)
from [`trace-csv-makecsv-byte-pairs-v1.gdb`](./trace-csv-makecsv-byte-pairs-v1.gdb).

## CSV serializer exhaustive byte-pair domain

The follow-up executes the serializer 65,536 times inside the target process,
covering every ordered backing-buffer pair `[a,b]` followed by NUL. Capacity
8 is sufficient for all serialized outputs in this domain. Every call captures
raw low AX, all eight output bytes, and prefix/post guards. The capture has
65,536 records with intact guards and return values of -1 for the 256 pairs
whose first byte is NUL and 1 for the remaining 65,280 nonempty C strings.
The second byte is not part of the string when the first byte is NUL.

For every nonempty input, a first byte from `0x80` through `0xff` consumes any
following nonzero byte as a pair and copies both bytes unchanged, including a
quote or comma in second position. If the following byte is NUL, only the
high-bit byte is copied. With an ASCII first byte, each byte is processed
separately: quotes are doubled in either position, and every nonzero second
byte, including high-bit values, is copied as one byte. This is complete for
one- and two-byte NUL-terminated inputs. It does not prove locale or Windows
codepage semantics; runtime evidence and `FUN_10016a40` pseudocode show the
helper treating every high-bit first byte as a pair lead without checking a
locale lead-byte table.

Run with a unique ID:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-makecsv-byte-pairs-full-v1.sh matrix2
```

The native routine is `csv-makecsv-byte-pairs-native.S`, with its assembled
image in `csv-makecsv-byte-pairs-native.bin`. Rebuild it with host GNU binutils:

```sh
as --32 tools/revkit/work/stage21/csv-makecsv-byte-pairs-native.S \
  -o /tmp/csv-makecsv-byte-pairs-native.o
objcopy -O binary --only-section=.text \
  /tmp/csv-makecsv-byte-pairs-native.o \
  tools/revkit/work/stage21/csv-makecsv-byte-pairs-native.bin
```

The runner emits a bytewise GDB loader because remote `restore binary`
duplicated part of the image in this Wine/GDB setup. Each record is 14 bytes:
sign-extended low AX, eight output bytes, prefix guard, and post-capacity
guard; six trailing bytes record the last input pair and completed-call count.
The completed capture is
[`csv-makecsv-byte-pairs-full-native14.bin`](./csv-makecsv-byte-pairs-full-native14.bin),
with runtime trace [`csv-makecsv-byte-pairs-full-native14-api.log`](./csv-makecsv-byte-pairs-full-native14-api.log).
Verify every record, return, guard, and completion marker with:

```sh
python3 tools/revkit/work/stage21/verify_csv_makecsv_byte_pairs_full.py \
  tools/revkit/work/stage21/csv-makecsv-byte-pairs-full-native14.bin
```

## CSV serializer field/count/capacity matrix

The 2,404-call matrix covers 87 field arrays: all `A`/empty placements through
five fields, longer arrays of 8–128 fields with selected empty positions, and
four mixed quote/high-bit arrays. Small arrays sweep capacities 1–32; scaled
arrays include the exact output-size boundary. Every return and every byte
inside the advertised capacity matches the independent model, with intact
prefix and post-capacity guards. An empty field stops serialization after
its opening quote. For `N` fields containing `A`, capacity `4*N` is the first
success. Two-byte quote escapes and high-bit pairs are written together only
when both destination bytes precede the final NUL sentinel.

Generate the trace, run with a unique ID, and verify the capture:

```sh
python3 tools/revkit/work/stage21/csv_makecsv_field_matrix.py generate \
  tools/revkit/work/stage21/trace-csv-makecsv-fields-v1.gdb
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-makecsv-fields-v1.sh matrix2
python3 tools/revkit/work/stage21/csv_makecsv_field_matrix.py verify \
  tools/revkit/work/stage21/csv-makecsv-fields-matrix2-api.log
```

The accepted capture is
[`csv-makecsv-fields-matrix1-api.log`](./csv-makecsv-fields-matrix1-api.log).
The verifier was also checked against altered output, altered return value,
altered guard, missing record, and duplicate record; all five were rejected.
Regenerating the trace produced identical bytes. The runner restores and
compares the Stage 5 input and WAVE fixtures before exit.

## CSV serializer pointer access order

Fifteen isolated target-stack probes confirm that the output is initialized
before the array loop; the current field pointer is fetched at `0x100169a7`
before checking room for an opening quote; and the pointed-to string is scanned
at `0x10016a53` only after the opening quote is written. A null array therefore
faults at capacity 1, while an accessible array containing a null or
`0xffffffff` string pointer returns -1 at capacity 1 and faults at capacity 2.

With `A` followed by a null string pointer, capacities 4 and 5 return -1,
while capacity 6 reaches the string scan and faults. With only the first array
entry accessible and the second entry on a no-access page, capacity 4 returns
-1, but capacities 5 and 6 fault fetching the second array entry. An empty
first field prevents that later access even with count `INT_MAX`. A null array
with count `INT_MIN` returns 1 and an empty output at capacity 1.

The generator's `CASES` table lists all expected returns, fault EIPs, and exact
16-byte output snapshots. Guard setup is checked explicitly. Reproduce:

```sh
python3 tools/revkit/work/stage21/csv_makecsv_pointer_order.py generate \
  tools/revkit/work/stage21
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-pointer-order-v1.sh matrix2
python3 tools/revkit/work/stage21/csv_makecsv_pointer_order.py verify \
  tools/revkit/work/stage21 matrix2
```

Accepted captures are `csv-pointer-order-matrix1-<case>-api.log` for all 15
names in the generator. The guard cases were added in a separate invocation
using the runner's `guards` group; default replay runs all 15 serially. Each
process is killed after the expected return or fault checkpoint. The runner
restores and compares the Stage 5 fixtures. The verifier also rejects altered
EIPs, returns, bytes, failed/null guard setup, and contradictory outcomes;
regenerated traces are byte-identical.

## CSV serializer input/output overlap

The 95-call matrix initializes a single field to `AB` at `output+k`, covering
offsets 0–9 and every capacity `k+3`–16. The entire initial string lies inside
the region filled by the API. Every exact 18-byte physical-buffer snapshot,
return value, and field-array pointer check matches the independent model.

The initial fill overwrites `AB`. Offset 0 then propagates quotes through the
overlapping source while quote doubling proceeds; offset 1 copies fill bytes
in place and cannot write a closing quote; offsets 2–9 return 1 after
serializing the fill bytes. For capacity 8, offset 2 returns the string
`"11111"`, demonstrating that success does not mean the overlapping input
was preserved. Other overlap placements remain untested.

```sh
python3 tools/revkit/work/stage21/csv_makecsv_overlap.py generate \
  tools/revkit/work/stage21/trace-csv-makecsv-overlap-v1.gdb
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-makecsv-overlap-v1.sh matrix2
python3 tools/revkit/work/stage21/csv_makecsv_overlap.py verify \
  tools/revkit/work/stage21/csv-makecsv-overlap-matrix2-api.log
```

Accepted capture:
[`csv-makecsv-overlap-matrix1-api.log`](./csv-makecsv-overlap-matrix1-api.log).
The runner restores and compares the Stage 5 fixtures.

## CSV serializer 128 MiB capacity

`VT_CsvParser_MakeCsv_ENG` returned raw low AX 1 for a 134,217,728-byte
advertised capacity using a physically matched allocation. The serialized
bytes were unchanged (`"A","b,c"`), the sacrificial prefix `0x5a` and
post-capacity guard `0xa5` were preserved, and the final advertised byte was
NUL. Together with the Stage 16 64 MiB call, this extends large positive
capacity observations to 128 MiB for this small input. Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-capacity-128m.sh
```

The trace and capture are `trace-csv-capacity-128m.gdb` and
`csv-capacity-128m-api.log`. Static pseudocode for `FUN_10016960` and its
serializer helpers shows writes into caller-supplied field/output pointers
and no internal heap allocation; allocation failure is therefore outside
this function's call path. Huge unsigned capacities and 32-bit output-size
overflow remain untested.

## CSV parser private delimiter field

`VT_CsvParser_Init_ENG` initializes the opaque 0x18-byte parser object's
`+0x14` pointer to a duplicated comma string. The Stage 21 probe changed this
private pointer to a semicolon in GDB. Default parsing treated `a;b;c` as one
field; the forced semicolon set produced three fields. With the override,
`a;"b;c";d` produced three fields and retained `b;c` as the quoted middle
field. The exported parse function has no delimiter argument, and no setter
export was found, so this is internal-state behavior rather than a supported
configuration contract. The serializer still emits comma-separated fields.
Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-delimiter-field.sh
```

The trace and capture are `trace-csv-delimiter-field.gdb` and
`csv-delimiter-field-api.log`; parser setup and scanner pseudocode are in
`reports/stage25-csv-parser-delimiter.c`.

## CSV parser complete single-byte interior domain

The 256-call matrix places each possible byte value in `A<byte>B,X`. Every
call returned low AX 1. NUL terminated the source after `A`; comma produced
three fields; the remaining 254 nonzero values were preserved byte-for-byte
inside the first field. This maps individual-byte handling at one unquoted
interior position, without inferring codepage or Unicode character semantics.
Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-byte-domain.sh
```

The trace and full per-value capture are `trace-csv-byte-domain.gdb` and
`csv-byte-domain-api.log`; the runner requires the exact
`CSV_BYTE_DOMAIN passed=256 failed=0` aggregate.

## CSV parser quote/comma products

The original probe enumerates all strings of length 0–6 over `A`, double
quote, and comma: 1,093 calls on one reused parser object. The length-7
follow-up repeats that product and adds all 2,187 length-7 strings, for 3,280
calls total. Every call returned low AX 1; captures record input hex, field
count, and all returned fields. Field-count distributions are documented in
the Lead 6 behavior report. This is exhaustive for the stated finite alphabet
and maximum length, not for quote interactions with whitespace, arbitrary
payload bytes, or longer strings. Reproduce the original run with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-quote-domain.sh
```

The original trace and capture are `trace-csv-quote-domain.gdb` and
`csv-quote-domain-api.log`; that runner requires exactly 1,093 cases.
Reproduce the length-7 follow-up with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-quote-domain-len7.sh
```

Its trace and capture are `trace-csv-quote-domain-len7.gdb` and
`csv-quote-domain-len7-api.log`; the runner requires exactly 3,280 cases and
zero low-AX failures, and restores and compares the Stage 5 fixtures.

## CSV IsCsv bounded first-record predicate

The 33-call matrix tests `VT_CsvParser_IsCsv_ENG` with expected counts around
the parsed count, `INT_MIN`/`INT_MAX`, byte bounds around and below the row
length, CRLF/LF input, and empty input. For `a,b,c`, expected values
`INT_MIN`, -1, and 0–3 return raw low AX 1; 4 and `INT_MAX` return 0. With
expected 3, bounds -1/0/2/4 return 0, while 5/6 return 1. Empty input
returns 0 for expected counts 0–2. For both `a,b\r\nc,d` and `a,b\nc,d`,
expected 2 returns 1 and expected 3 returns 0; each call uses a fresh buffer,
and the export truncates it to `a,b` at the newline. This separates IsCsv's
first-record preprocessing from the standalone parser's direct CR/LF field
behavior. Paired direct parse/count calls show 98–100 fields parse
successfully, while 101 fields return `-4` and retain a count of 100. IsCsv
still returns 1 for expected 100 and 0 for expected 101/102, so it uses that
partial count despite the parse error.
Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-iscsv-matrix.sh
```

The trace and capture are `trace-csv-iscsv-matrix.gdb` and
`csv-iscsv-matrix-api.log`; the runner requires all 33 IsCsv result lines and
restores and compares the Stage 5 fixtures.

### IsCsv capacity beyond 101 fields

The v2 sweep tests every field count from 99 through 128, plus 255, 256, 512,
1,024, and 4,096. Each call's byte bound includes the complete generated row
and its NUL. The direct parser accepts 99 and 100 fields; every tested count
above 100 returns `-4` and retains exactly 100. `IsCsv` reports success for
expected count 100 even on those parse failures, then returns 0 when expected
count is the input count or one greater. This extends the partial-count result
well beyond the first overflow boundary.

Run the isolated probe with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-iscsv-capacity-v2.sh
```

The runner restores and compares the Stage 21 sandbox fixtures. The trace and
capture are `trace-csv-iscsv-capacity-v2.gdb` and
`csv-iscsv-capacity-v2-api.log`.

### IsCsv signed bounds and 32-bit end-pointer wrap

The v3 probe calls `VT_CsvParser_IsCsv_ENG` directly on a valid `a,b,c` row
and varies argument 3 across signed extremes, short limits, and values chosen
to make the 32-bit end pointer wrap to `0`, `1`, `4`, `5`, or `0xffffffff`.
The runtime row address was `0x003e0670`. Limits -1024, -1, 0, 1, and 4
returned raw low AX 0; 5 and 6 returned 1. `INT_MIN` and `INT_MAX` returned
1 because their computed end addresses were `0x803e0670` and `0x803e066f`,
above the buffer address, and the bounded scan stopped at the row NUL. Wrapped
ends 0, 1, 4, and 5 returned 0 because those unsigned addresses precede the
row; end `0xffffffff` returned 1. This confirms the bound is pointer addition
with 32-bit wrap followed by an unsigned address comparison. The observed
signed-limit result depends on the buffer address; these calls do not establish
safe behavior for unterminated inputs or inaccessible memory.

Reproduce in the isolated Stage 21 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-iscsv-bounds-v3.sh
```

The runner restores and compares the Stage 21 sandbox fixtures. The trace and
capture are `trace-csv-iscsv-bounds-v3.gdb` and
`csv-iscsv-bounds-v3-api.log`.

## CSV parser allocation lifecycle

Two concurrent `VT_CsvParser_Init_ENG` calls returned distinct 24-byte
objects. Their first five dwords were `0,0,0,100,0`; their delimiter pointers
were distinct allocations, both containing a comma. A parsed `first,second`
object held a two-entry field array and an owned text copy. `VT_CsvParser_Exit_ENG`
completed for NULL, unparsed objects, and a parsed object. Static
`FUN_10016860` conditionally frees the owned copy, field array, delimiter, and
object, while leaving the caller's original text pointer alone.

Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-csv-lifecycle.sh
```

The trace and capture are `trace-csv-lifecycle.gdb` and
`csv-lifecycle-api.log`; the runner verifies all five lifecycle observations
and restores and compares the Stage 5 fixtures.

## CSV getter null and index boundaries

`VT_CsvParser_GetNfields_ENG` returns 0 for a null object and a freshly
initialized, unparsed object. `VT_CsvParser_GetField_ENG` returns null for
those objects and for indexes 2, 3, and `INT_MAX` on a parsed two-field row.
Indexes 0 and 1 return the stored field pointers. Index -1 returns a nonnull
raw pointer read from before the field array; the probe records the pointer
without dereferencing it. The helper's static code checks only whether the
index is below the field count, with no lower-bound check.

Run in the isolated Stage 21 Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-csv-getter-boundaries.sh
```

The capture is `csv-getter-boundaries-api.log`, reproduced with
`trace-csv-getter-boundaries.gdb`. The runner preserves and compares the
sandbox input and output fixtures.

## Lip-sync report writer on a missing parent path

With `no-such-lipsync-dir/report.txt`, target-flow tracing reached the report
writer with the active state object's subobject at `state+0x2c` present and
its writer field at `+0x10` equal to null. The wrapper received null in EDX
and dereferenced it at `0x10025e4d` without checking it first. This captures
the immediate pointer and caller chain. A deeper trace follows the path through
`FUN_1001df10` and `FUN_10025dc0` to the stream helpers: mode `"wt"` reaches
the low-level open routine, `CreateFileA` (import slot `0x1006d054`) returns
`INVALID_HANDLE_VALUE` (`0xffffffff`), and `GetLastError` returns 3
(`ERROR_PATH_NOT_FOUND`). The open helper returns zero. `FUN_1001df10` stores
that zero at wrapper `+0x10` but still returns the allocated wrapper, which is
later passed into the unchecked writer dereference. The directory-target v4
case follows the same path with `ERROR_ACCESS_DENIED` (5). Static code funnels
all `CreateFileA` invalid-handle results through the same null-return path;
exact error codes from other path classes remain untested.

Run in the isolated Stage 21 Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-lipsync-missing-parent-null-writer-v3.sh
```

The v2 trace stops before executing the faulting dereference. Its GDB trace is
`trace-lipsync-missing-parent-null-writer-v2.gdb`, and its capture is
`lipsync-missing-parent-null-writer-v2-api.log`. The expanded v3 trace captures
the open return, mode, Win32 failure code, wrapper field, and later null
dereference; its files are `trace-lipsync-missing-parent-null-writer-v3.gdb`
and `lipsync-missing-parent-null-writer-v3-api.log`. Run it with
`run-lipsync-missing-parent-null-writer-v3.sh`. Both runners preserve and
compare the sandbox input and output fixtures.

The directory-target variant uses the existing `.` directory as the output
name. Its `CreateFileA` call returns `INVALID_HANDLE_VALUE` with
`GetLastError` 5 (`ERROR_ACCESS_DENIED`), and the same null writer field and
unchecked dereference follow. Reproduce it with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-lipsync-directory-target-null-writer-v4.sh
```

Its trace and capture are `trace-lipsync-directory-target-null-writer-v4.gdb`
and `lipsync-directory-target-null-writer-v4-api.log`.

## Preprocess flag-7 prefix recurrence

The DLL uses an eight-entry period-string table for the leading flag-7 filler.
Static arithmetic in `FUN_1001cd70` shows a one-time `GetTickCount` seed and a
stored-state recurrence equivalent to
`S' = 16807*(S mod 127773) - 2836*floor(S/127773)`, adding `2147483647` when
the result is nonpositive. The next-state residue modulo 8 selects the table
entry. A same-process sequence of three identical `Hello world.` calls
returned states `938191541`, `1360293313`, and `338805629`; the first two
recurrence calculations predict the next states exactly. Their residues
selected period counts 151, 111, and 151, matching output files `s2`, `s3`, and
`s4`. The runtime log is `preprocess-rng-sequence-v2-api.log`; the reproducible
trace and runner are `trace-preprocess-rng-sequence-v2.gdb` and
`run-preprocess-rng-sequence-v2.sh`. The runner preserves and compares the
Stage 5 fixtures. Verify the captured state transitions and prefix lengths
with `python3 tools/revkit/scripts/check_preprocess_rng_sequence.py`. This
recovers prefix selection, not the reason the filler is present.

## License token byte sensitivity

`license-token-byte-matrix.c` reads the archived verification record and
locates its single 96-character hexadecimal field without printing its
contents. It makes two in-memory mutations at each byte position: replace the
byte with the next hexadecimal digit, then replace it with `g`. The original
record returned 0; all 192 one-byte variants returned -2. This tests selected
positions and the intervening bytes across both 48-character halves. The
capture records only byte index, selected/gap class, mutation class, and
return code. The host-side analyzer verifies all 96 indexes in both mutation
modes, their selected/gap classifications, and the observed return codes.
Reproduce with:

```sh
bash tools/revkit/work/stage21/run-license-token-byte-matrix.sh
```

The outer runner builds with `vtspeak-pe32-builder:local`, then runs in the
read-only Stage 8 Wine runtime. The first capture is
`license-token-byte-matrix-api.log`. Version 2 adds case flips for each
alphabetic hex byte and writes `license-token-byte-matrix-v2-api.log`: 10 of
22 flips preserved return 0, all in second-half gap positions; 12 returned
-2. Version 3 tests all 1,024 combinations of the ten individually accepted
case flips; every combination returned 0. It writes
`license-token-byte-matrix-v3-api.log`. The analyzer checks the 192
value-changing variants, case-flip positions, and the complete combination
product. These captures establish byte-change and case behavior for one
supplied record, but do not generalize the token format to other records.

## License positional-field acceptance screen

`license-positional-field-matrix.c` makes same-length in-memory edits to the
sample expiry token, XML `expdate`, `VW_VTAPI`, and positional/XML channel
values. It prints only mutation case labels, byte count, and checker results.
The untouched record returns 0; each single edit and each tested paired edit
returns -2. The result establishes acceptance sensitivity with the original
derived token unchanged; it does not show that fields are directly compared
or establish their meanings. Reproduce with:

```sh
bash tools/revkit/work/stage21/run-license-positional-field-matrix.sh
```

The outer runner builds with the isolated PE32 builder and executes in the
Stage 8 Wine runtime. The capture is
`license-positional-field-matrix-api.log`.

## License token MD5 path

`probe-license-md5.c` calls the DLL's internal init/update/final helpers at
module offsets `0x142f0`, `0x14320`, and `0x14410`. It checks ten standard MD5
vectors: empty input, `a`, `abc`, and repeated-`a` inputs of lengths 55, 56,
63, 64, 65, 127, and 128 bytes. The results all match the reference digests;
the capture contains only lengths and match flags. Reproduce with:

```sh
bash tools/revkit/work/stage21/run-license-md5.sh
```

The capture is `license-md5-vectors-api.log`. Static disassembly shows
`FUN_10015110` concatenates its four string arguments in order 3, 4, 2, 1
without separators before hashing, and `FUN_10015220` encodes the resulting
16-byte digest as 32 lowercase hexadecimal characters. This verifies the MD5
primitive used in the license conversion path; it does not establish that the
license token is a signature or that MD5 provides an authentication guarantee.

## SyncInfo allocator transient-failure retry

The PE32 harness calls `VT_AllocSyncInfo_New_ENG` and
`VT_FreeSyncInfo_New_ENG` directly. GDB forces the first 520-byte allocation
for a nested entry to return null after the allocator helper's allocation call.
The helper waits 10 ms, retries, receives a nonnull result, and continues
through all 600 nested allocations. The object passes the runtime shape checks
for 600 rows, 65 entries per row, a row array, and distinct first/last nested
allocations; the direct Wine control reports success as well. This confirms
recovery from one transient allocation failure at that callsite. It does not
establish behavior under persistent memory exhaustion. Reproduce with:

```sh
bash tools/revkit/work/stage21/run-syncinfo-allocation-retry-v7.sh
```

The probe source is `probe-syncinfo-allocation-retry.c`; the GDB trace is
`trace-syncinfo-allocation-retry-v7.gdb`; the sanitized capture is
`syncinfo-allocation-retry-v7-api.log`. The injected-run checkpoint confirms
the returned object's 600×65 dimensions and distinct nested allocations.

## Preprocess flag 7 prefix table

The controlled runtime sweep calls `VT_TextToPreprocessInfoFile_ENG` eight
times with flag 7. Before each call it seeds the stored generator state so
the next state selects one desired residue modulo 8. The files `p0` through
`p7` verify every table entry's emitted period count: 108, 111, 140, 52, 33,
151, 106, and 58, respectively. The API log records the seed, resulting
state, raw EAX observation, and output path. This completes the table-index
to-prefix-length mapping; it does not explain why the report includes the
prefix or characterize natural seed distribution.

The GDB trace is `trace-preprocess-rng-all-prefixes.gdb`; the runner is
`run-preprocess-rng-all-prefixes.sh`; and the capture is
`preprocess-rng-all-prefixes-api.log`. Reproduce with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-preprocess-rng-all-prefixes.sh
```

## Preprocess report option matrix

The 44-call matrix uses `Hello world.`, speaker 1, dictionary 0, and text type 0.
For flags 3, 5, 7, and 10 it tests the all-`-1` baseline; pitch 50/200;
speed 50/400; volume 0/500; pause 0/250/65,535; and a repeated baseline.
All calls returned raw EAX `1`. Each flag's 11 files are byte-identical,
including flag 7 with its generator state reset before each call. This result
is bounded to this text, these flags, and these option values.

The trace is `trace-preprocess-settings-matrix.gdb`, the runner is
`run-preprocess-settings-matrix.sh`, and the capture is
`preprocess-settings-matrix-api.log`. From the repository root, run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-preprocess-settings-matrix.sh
python3 tools/revkit/scripts/check_preprocess_settings_matrix.py
```

The checker validates the 44 call records and byte-compares each option set.
The four common outputs are 53, 29, 202, and 15 bytes for flags 3, 5, 7, and
10, with SHA-256 values printed by the checker.

The first trial used a misordered argument list and is retained separately in
`quarantine-preprocess-option-argument-order/`; those files are not evidence
for the results above.

## Preprocess flag 6 PCM settings

The 11-call flag-6 matrix uses the same `Hello world.` fixture and setting
variants as the report matrix. It captures the API's fixed `test.pcm` output
after each direct call. The default and repeated-default files are both
23,606 bytes (11,803 PCM16 samples) and byte-identical. Pitch 50/200 produces
12,850/11,998 samples; speed 50/400 produces 25,377/2,438. Volume 0 produces
11,803 zero samples, while volume 500 changes the samples without changing
the count. Pause 0, 250, and 65,535 are byte-identical to the default on this
period-terminated utterance. All calls return raw EAX `1` through the void
wrapper. These results characterize only this fixture and these settings.

The trace is `trace-preprocess-pcm-settings.gdb`; the runner is
`run-preprocess-pcm-settings.sh`; the API log is
`preprocess-pcm-settings-api.log`; and the 11 captured files are named
`preprocess-pcm-settings-v00.pcm` through `v10.pcm`. Verify with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-preprocess-pcm-settings.sh
python3 tools/revkit/scripts/check_preprocess_pcm_settings.py
```

## `VT_GetTTSInfo_ENG` pointer and capacity precedence

The pointer matrix passed null output to valid requests 1, 2, 3, 23, and 101,
and invalid requests -1 and 27; all returned 3. With nonnull output, invalid
requests returned 2 and preserved the sentinel. `VT_BUILD_DATE` needs 12
bytes, `VT_DB_DIRECTORY` needs 4, and short capacities return 4 without
writing. A missing request-23 date file returns success without writing even
at capacity 1. Request 101 returns the playback-state value as status and
preserves its output word at tested size -1.

Run in the isolated Stage 21 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-info-pointer-capacity-matrix.sh
```

The trace and capture are `trace-info-pointer-capacity-matrix.gdb` and
`info-pointer-capacity-matrix-api.log`.

## `VT_GetTTSInfo_ENG` request-ID sweep

The trace calls every request integer from -128 through 256 with a
four-byte nonnull destination initialized to `0xa5a5a5a5`. All 357 IDs outside
the declared range 0–26 and special request 101 return 2 and preserve the
sentinel. ID 0 returns short-length 4 because four bytes cannot hold its
build-date string. ID 23 returns success without writing because the date
file is absent; ID 101 returns the pre-synthesis state as its status and leaves
the destination alone. Other declared requests return success. Separate
probes cover `INT_MIN` and `INT_MAX`.

Run in the isolated Stage 21 Stage 5 sandbox:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage21/run-info-request-id-sweep.sh
```

The runner saves and restores the sandbox input and output. The complete
request-by-request capture is `info-request-id-sweep-api.log`; reproduce it
with `trace-info-request-id-sweep.gdb`.

## `VT_GetTTSInfo_ENG` request 23 with a present date file

The probe creates a temporary `/work/db_build.date` containing `D23-OK` and
one newline, then calls request 23 with destination capacities `-1`, `0`, `1`,
`10`, `11`, `12`, `16`, and `64`. Each call receives a 64-byte buffer filled
with `0x5a`; the trace records the API result and the first 16 bytes. The
runner refuses to overwrite an existing date file or capture and removes only
the date file it creates. It stops the sample at the pre-synthesis breakpoint
after the queries.

From the repository root, run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-info23-date-file.sh
```

`run-info23-capacity-edges.sh` repeats against the same temporary payload for
capacities 2 through 9, resolving the minimum size boundary around the
seven-byte file content and terminating NUL. It uses distinct capture names:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-info23-capacity-edges.sh
```

`run-info23-long-file.sh` creates a temporary 1,100-byte file and compares
destination capacities 1,023, 1,024, and 1,100. It also calls request 23 with
a null destination after the file is present, recording the return without
reading through a null pointer.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-info23-long-file.sh
```

`run-info23-file-shapes.sh` calls request 23 against two temporary files:
`A\0B\nC\n` (embedded NUL) and `LINE1\nLINE2\n` (two lines). Each call
uses a 64-byte sentinel buffer and captures the first 24 bytes.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm \
  -w /work/stage5 runtime /bin/bash \
  /work/stage21/run-info23-file-shapes.sh
```
