# Stage 10: controlled linguistic and abbreviation contrasts

`inputs/` contains 46 targeted abbreviation/context fixtures, 11 generated
all-key H/T batches, and the controlled VTML phoneme fixtures.
`run.sh` executes each selected fixture in a
fresh original Paul process under Wine/GDB. The GDB script breaks at
`FUN_10007520` (`0x10007520`), records the upstream `0x94`-byte token surface
at `+0x34`, and records the returned `0x70`-byte phone/context row's flags,
source index, class, surface, and first 24 code bytes. The row's surface is at
`+0x07` and code at `+0x25`. The resulting logs are named after the input
fixtures. `work/scripts/verify_stage10.py` checks the complete ordered surface
sequence and normal process exit for every capture.

The [VoiceText VTML guide](https://static.carahsoft.com/concrete/files/1615/2520/8261/Voice-Text_Markup_Language.pdf)
defines the `x-cmu` input used in `phoneme-tomato.txt` and the generated
`cmu-*` fixtures. `tools/revkit/scripts/generate_cmu_probes.py` creates those
small input files from its checked-in symbol list. The upstream token record
at `+0x52` contains the forced internal phone bytes. The verifier below
checks all 24 consonants, 15 vowels at three stress levels, the guide's
tomato example, a three-tag control, and a separate stress contrast:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_cmu_probes.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage10/run.sh cmu-consonants-1 cmu-consonants-2 cmu-consonants-3 \
  cmu-vowel-stress-1 cmu-vowel-stress-2 cmu-vowel-stress-3 \
  cmu-vowel-stress-4 cmu-vowel-stress-5 cmu-vowel-stress-6 \
  cmu-stress phoneme-tomato phoneme-trio
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_cmu_codebook.py
```

`pos-record-single-*` fixtures apply the guide's six part-of-speech values
and a no-tag control to `record.`. The same context GDB script captures the
upstream `+0x30` status, returned phone codes, and four trailing dictionary
metadata words. `trace-pos-trees.gdb` captures decision-tree input vectors
for the noun/modifier and verb/interjection pairs, whose phone sequences are
the same within each pair. The Stage 10 verifier checks these controls and
their matching tree input sequences. Five `feature-*` fixtures hold the
visible token constant while changing its forced first consonant or AH stress;
the verifier checks the exact tree-input differences. Reproduce the tree
traces with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-pos-trees.gdb \
  TRACE_SUFFIX=-trees /bin/bash /work/stage10/run.sh \
  pos-record-single-noun pos-record-single-modifier \
  pos-record-single-verb pos-record-single-interjection \
  feature-p-ah0 feature-b-ah0 feature-t-ah0 feature-p-ah1 feature-p-ah2'
```

`trace-lookups.gdb` and the two `*-lookups.log` files provide an additional
comparison of the address-alias lookup helpers for `apartment-mixed` and
`apartment-second`. `trace-apt-branch.gdb` and its two `*-branch.log` files
identify which occurrence reaches the address replacement branch. Run either
focused trace with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-lookups.gdb \
  TRACE_SUFFIX=-lookups /bin/bash /work/stage10/run.sh apartment-mixed apartment-second'
```

Use `TRACE_SCRIPT=/work/stage10/trace-apt-branch.gdb` and
`TRACE_SUFFIX=-branch` for the branch trace.
`trace-address-guard.gdb` and the three `*-guard.log` files capture the
address-rule guard on the repeated `Apt` controls; only the first `Apt` of
the repeated expression reaches that guard in these runs.
The phrase-control `*-branch.log` files show one alias branch for the
uninterrupted `Apt number and Apt number` form, and two visits when a sentence
boundary, comma, or intervening word separates the occurrences.
`trace-address-call.gdb` records entry text and return lengths from
`FUN_10043810`. The repeated form consumes the 14-byte prefix ` Apt 4 and Apt`
in one match; the word-break form consumes ` Apt 4 and Hello` (16 bytes) and
then matches the later ` Apt 4.` independently. The comma form returns once
for each six-byte ` Apt 4` span. These captures locate the repeated-key
suppression in address-phrase span matching before the phone/context pass.
Reproduce these three call/return captures from the repository root with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-address-call.gdb \
  TRACE_SUFFIX=-address-call /bin/bash /work/stage10/run.sh \
  apartment-repeated apartment-word-break-repeat apartment-comma-repeat'
```

The generator also creates one-token `tree-phone-consonant-*` and
`tree-phone-vowel-*` fixtures. Each isolates a CMU symbol in its own original
process. Capture and summarize their decision-tree input fields with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-pos-trees.gdb \
  TRACE_SUFFIX=-trees /bin/bash /work/stage10/run.sh \
  "tree-phone-consonant-*" "tree-phone-vowel-*"'
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/analyze_stage10_tree_features.py
```

`trace-context-downstream.gdb` extends the adjacent-phone controls through
the selected-ID handoff (`FUN_1001b200`) and the primary timeline rows
(`FUN_1002c220`). It records the full scalar tree input vectors, tree
returns, selected IDs, and all 13 words of each returned 52-byte row. Re-run
the four P/B controls from the repository root with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-downstream.gdb \
  TRACE_SUFFIX=-downstream /bin/bash /work/stage10/run.sh \
  tree-context-left-p tree-context-left-b \
  tree-context-right-p tree-context-right-b'
```

The matching `*-downstream.log` files preserve the four successful captures.
The runner restores the Stage 5 input and output on exit. Interpret the first
six selected IDs in order as the six primary timeline rows, as established
for the Stage 8 selection-to-timeline path; subsequent selected IDs are
additional context-neighbor selections.

The repeated fixed-`T AH0` contrast and period-separated control use the same
trace:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-downstream.gdb \
  TRACE_SUFFIX=-downstream /bin/bash /work/stage10/run.sh \
  tree-context-t-left-p tree-context-t-left-b \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-boundary-p tree-context-t-boundary-b'
```

The sentence-boundary controls produce two timeline calls each; compare their
first call separately from the next sentence's rows.

`trace-context-scoring.gdb` follows per-context candidate generation and
local costs at `FUN_10023350`, transition costs and predecessor links at
`FUN_10018c80`, pruning at `FUN_10024510`, and the final backtrack at
`FUN_10024900`. It also captures candidate metadata before and after
`FUN_100230a0`, local and cumulative costs, and selected predecessor indexes.
For focused scorer calls, it records the scale and detailed inputs for the
candidate IDs whose backtracked path changes. A targeted trace follows
candidate `3837` through the forward scan at `0x1002326b` and accumulator at
`0x10023270`, including context-row byte `+0x04` and the selected table value. The
static excerpt for
`FUN_100230a0` is in `tools/revkit/work/reports/`. Run the fixed-T right-side
and period-separated pairs with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-scoring.gdb \
  TRACE_SUFFIX=-scoring-row /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-boundary-p tree-context-t-boundary-b'
```

The runner restores the Stage 5 input and output after each fixture.

The additional F/V contrast and before/after P/B metadata captures use the
same trace:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-scoring.gdb \
  TRACE_SUFFIX=-scoring-row4 /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

The corresponding `*-scoring-row4.log` files preserve the before/after
records, row-byte/table lookup, and complete candidate traces.

The context-row byte writer and its fallback caller can be reproduced with
the focused traces below. The watchpoint trace records writes to context 1's
byte `+0x04`; the row-build trace records which context indices enter
`FUN_100242a0`:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-row-writer.gdb \
  TRACE_SUFFIX=-row-writer-verified /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-row-build.gdb \
  TRACE_SUFFIX=-row-build /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

The writer is in `FUN_100242a0` at `0x10024305` for value 1 and in its
second-row copy path for value 2. `FUN_10024680` calls it only when
`FUN_10024060` returns false. These traces locate the writer and conditional
path, but do not explain the row-byte values semantically.

To capture the preceding candidate-pass return values directly at the
caller's branch, use:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass.gdb \
  TRACE_SUFFIX=-firstpass /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

This logs `FUN_10024060`'s signed return, context index, and row byte at
return. At `0x1002480b`, the caller takes the fallback for a zero result and
skips it for a positive result. The candidate-generation cause and semantic
interpretation remain open.

To trace candidate counts and table weights feeding that return, use the
count trace:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass-counts.gdb \
  TRACE_SUFFIX=-firstpass-audit /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

This records generated and ranked counts, the final candidate IDs and their
16-bit model-table values, and the `FUN_10023060` sum before the threshold
branch. For the captured flag-0 cases, a sum of 10 or more takes the first
pass; the byte and table values still have no assigned linguistic meaning.

To compare the seven-byte input records delivered to `FUN_10018770` for
those same contexts, run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass-inputs.gdb \
  TRACE_SUFFIX=-firstpass-inputs /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

This records the raw seven bytes and row bytes at the `FUN_10018770` entry
called by `FUN_10024060`; the byte positions remain unnamed.

To record the translation-table checks for the changed context-1 byte, its
branch route, combined-signature lookup results, and emitted candidate IDs,
use:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass-byte-branch.gdb \
  TRACE_SUFFIX=-byte-path-emitted /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

The breakpoints follow the `FUN_10018770` branches at `0x10018ae4`,
`0x10018af1`, and `0x10018b03`; they report raw byte/table values, branch
routes, successful combined-signature lookups, and emitted candidate IDs.
Static disassembly reduces the seven-byte signature to a five-byte key and
searches sorted five-byte records. To capture those transformed keys and the
matching record ranges for the controlled P/B/F/V queries, use:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass-class-ranges.gdb \
  TRACE_SUFFIX=-class-ranges /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

The `*-class-ranges.log` captures report the transformed key, its equal-key
range, and the five-byte model records in that range. The probes keep the
records numeric and do not assign linguistic meanings. In the captured
context-1 queries, P produced one exact class match, B two, and F/V none;
the returned class indices equal the class-list entries for the matches.

To follow class candidates through `FUN_10023350` into model unit IDs, use:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-unit-expansion.gdb \
  TRACE_SUFFIX=-unit-expansion /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b'
```

The trace reads the 16-bit class populations through the model pointer at
`+0x8c` and each class's unit-ID list through `+0x94`. It prints up to 100
members per class. In the B fixture both matching classes are fully captured;
the complete expansion totals 14 unit IDs.

To capture the whole/half target lookups, transformed 10-byte keys, and
candidate class populations for the controlled right-phone P/B/F/V cases, run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-half-keys.gdb \
  TRACE_SUFFIX=-half-keys-verified /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p tree-context-t-right-b \
  tree-context-t-right-f tree-context-t-right-v'
```

The trace captures the arguments to `FUN_10019570`, key bytes passed to
`FUN_1002df50`, returned range lengths, and class IDs/populations returned by
`FUN_10023c70`. It masks that function's packed return to its low 16-bit
candidate count. The logs are `tree-context-t-right-*-half-keys-verified.log`.
The B fixture has no half-key query because its whole-phone class population
passes the split threshold. See the neighbor-context report for the numeric
key and candidate comparison; key slots remain semantically unnamed.

An additional controlled K/G pair uses the same traces and input fixtures:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-half-keys.gdb \
  TRACE_SUFFIX=-half-keys-kg /bin/bash /work/stage10/run.sh \
  tree-context-t-right-k tree-context-t-right-g'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass.gdb \
  TRACE_SUFFIX=-firstpass-kg /bin/bash /work/stage10/run.sh \
  tree-context-t-right-k tree-context-t-right-g'
```

The left-neighbor P/B threshold reversal and half-key inputs use the existing
left controls:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass-inputs.gdb \
  TRACE_SUFFIX=-firstpass-inputs-left /bin/bash /work/stage10/run.sh \
  tree-context-t-left-p tree-context-t-left-b'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-firstpass.gdb \
  TRACE_SUFFIX=-firstpass-left /bin/bash /work/stage10/run.sh \
  tree-context-t-left-p tree-context-t-left-b'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-half-keys.gdb \
  TRACE_SUFFIX=-half-keys-left /bin/bash /work/stage10/run.sh \
  tree-context-t-left-p tree-context-t-left-b'
```

The three-phone B/T/P/B/F/V controls compare the same context-4 row while
varying input byte `+3`. Capture all four first-pass rows and returns, then
the resulting half-key searches, with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-middle-audit.gdb \
  TRACE_SUFFIX=-middle-audit-all /bin/bash /work/stage10/run.sh \
  tree-context-t-middle-b-p tree-context-t-middle-b-b \
  tree-context-t-middle-b-f tree-context-t-middle-b-v'
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-half-keys.gdb \
  TRACE_SUFFIX=-half-keys-all /bin/bash /work/stage10/run.sh \
  tree-context-t-middle-b-p tree-context-t-middle-b-b \
  tree-context-t-middle-b-f tree-context-t-middle-b-v'
```

At context 4, P/F/V return zero and emit left/right half-key lookups; B
returns one and stays on the whole-phone path. The matching
`*-middle-audit-all.log` and `*-half-keys-all.log` files preserve the raw
inputs, six-byte rows, returns, key bytes, ranges, and candidate populations.

The GDB class-duplicate audit reads the loaded five-byte key array in one bulk
memory transfer and reports all repeated-key ranges:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash -c 'TRACE_SCRIPT=/work/stage10/trace-context-class-duplicates.gdb \
  TRACE_SUFFIX=-class-duplicates /bin/bash /work/stage10/run.sh \
  tree-context-t-right-p'
```

Its log reports the loaded record count, repeated-key groups, records in those
groups, and maximum run size. The current Paul M16 capture reports zero
repeated keys across 61,566 records.

The runner backs up and restores the Stage 5 input and output in the ignored
`work/corpus-parity/` scratch directory. It accepts optional fixture basenames
for focused reruns. From the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage10/run.sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_stage10.py
```

The curated interpretation and its limits are in
`docs/reverse-engineering/abbreviation-context-contrasts-2026-09-24.md`.

To generate and run one controlled context for all 55 H-set keys and 42 T-set
keys, then verify the 87 resulting context-return calls:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_abbreviation_probes.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage10/run.sh "abbr-h-*" "abbr-t-*"
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_stage10.py
```

The C set has 330 keys spanning mixed categories. The separate probe uses one
neutral `We noted KEY today.` context per key, in batches of ten. It preserves
the key order in `abbrc-probe-manifest.json`; the summarizer requires all 330
context calls and normal process exits, then records the full returned
surface sequences in `abbrc-probe-results.json`:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_abbrc_probes.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage10/run.sh 'abbr-c-*'
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_abbrc_probes.py
```

This gives one baseline observation per key. It is not a per-key context
matrix, and the neutral sentence is not suitable for resolving date, time,
measurement, location, or case-sensitive distinctions.
