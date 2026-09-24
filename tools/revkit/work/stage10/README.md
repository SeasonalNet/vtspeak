# Stage 10: controlled linguistic and abbreviation contrasts

`inputs/` contains 40 targeted abbreviation/context fixtures, 11 generated
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
