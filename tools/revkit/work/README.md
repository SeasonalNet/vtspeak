# Portable reverse-engineering evidence

This directory holds generated analysis evidence that the research documents
cite. Portable traces, small captured fixtures, reference outputs, scripts, and
Ghidra text exports are tracked so a checkout can inspect and recheck the
documented results.

## Tracked material

- `reports/` contains decompiler pseudocode, static-analysis reports, and
  disassembly excerpts. These are generated observations, not original source
  code.
- `scripts/` contains analysis helpers used by the research workflow.
- `stage2-copy/` contains small DAT payload excerpts, their captured PCM
  outputs, and the GDB capture scripts/logs used for decoder comparisons.
- `stage3/` through `stage16/` contain the controlled text inputs, trace
  scripts/logs, PCM buffers, WAVE outputs, and comparison evidence cited in the
  stage findings. `stage5/probes/stage6/` holds the captured tree lookups used
  by the checked-in comparison helper.
- `stage6-resource-report.json` records the checked-in TPP inventory result.

The outputs in this directory are evidence from the proprietary VoiceText
engine and voice package. Their inclusion makes these small research fixtures
available for local verification; it does not assert a license for the engine
or voice assets.

## Local-only exclusions

The following remain ignored because they are host-specific, large generated
state, or complete local copies of proprietary inputs:

- `home/`: Ghidra preferences, caches, and shell history.
- `projects/`: Ghidra project databases and caches.
- `data-common` and `data-paul`: container bind-mount links.
- `stage1-copy/` and `stage1-run/`: bulk runtime copies and very large raw
  syscall traces.
- `stage2-copy/wineprefix/` and its cache/config directories: Wine runtime
  state.
- The repository's `binary/`, `data-common/`, and per-voice `data-*` source
  inputs remain outside Git as before.

## Offline checks

Decode each captured Stage 2 payload and compare its PCM bytes with the
original-engine capture:

```sh
python3 -c '
import struct
import sys
from pathlib import Path

sys.path.insert(0, "tools/revkit/scripts")
from decode_dat import decode_payload

root = Path("tools/revkit/work/stage2-copy")
payloads = sorted(root.glob("candidate-*.dat"))
for path in payloads:
    samples = decode_payload(path.read_bytes())
    actual = struct.pack(f"<{len(samples)}h", *samples)
    expected = path.with_suffix(".pcm").read_bytes()
    if actual != expected:
        raise SystemExit(f"PCM mismatch: {path.name}")
print(f"matched {len(payloads)} captured DAT/PCM pairs")
'
```

For Stage 3 and Stage 8, concatenate the numbered `.pcm` blocks and compare
them with the corresponding WAVE data chunk; the stage reports record the
expected byte counts and hashes. Run the checked-in Stage 8 verifier from the
repository root to check the block/data matches, reference-WAVE identity, PCM
format, and controlled-output frame counts:

```sh
python3 tools/revkit/work/scripts/verify_stage8.py
```

With the full local 2013 Paul M16 assets present, inventory DAT frame controls
and compare each stream's decoded frame count with its combined UPM-derived
sample count:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/inventory_dat_modes.py
```

This checks stream framing and sample counts. The later full-corpus PCM value
comparison uses the PE32 probe and independent Python comparator documented in
the [corpus parity report](../../../docs/reverse-engineering/dat-corpus-parity-and-stage9-2026-09-24.md).

The Stage 9 verifier checks the two additional number/abbreviation WAVE files
against their captured PCM blocks and prior Stage 6 references:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_stage9.py
```

Stage 11 captures eight text-category WAVs with selected-unit IDs, timeline
row counts, and all PCM blocks. Its verifier checks each local WAV format,
compares the captured blocks with the WAV data chunk, and records whole-file
hash comparisons against the selected Wag peer manifest rows:

```sh
python3 tools/revkit/work/scripts/verify_stage11.py
```

Stage 12 contains controlled TPP family and suffix inputs, GDB captures for
typed lookup returns and numeric token-byte writes, and the resulting local
WAVs. Reproduce the runtime traces from the repository root with local vendor
inputs and the existing isolated Wine image:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run.sh \
  place-flags ax-special g-suffix f-suffix place-components
```

To capture the natural replacement-table path, use the `marker-producer`
mode with `marker-direction-address` and `marker-street-address-control`.
The trace logs lookup results, the `d` marker write, and the copied
replacement value; the `N` key hits and `Main` misses. The
`marker-apartment-address-span` and `marker-apartment-neutral` fixtures also
capture uppercase `A` and lowercase `a` writes in the separate apartment path.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-effects.sh marker-producer \
  marker-direction-address marker-street-address-control
```

To retrieve the two E records through the shared lookup helper, run the
scoped in-memory selector trace. It prints both returned strings, converts
each forced hit to a miss before resuming normal processing, and restores the
Stage 5 input/output files on exit:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-e-direct.sh \
  e-casa-hyphens e-saint-hyphens
```

### Lead 3 candidate-list trace

For the controlled `Can anybody help?` contrast, capture the class lists
returned by `FUN_10024680` and the unit-index lists returned by
`FUN_10023350`. The second run changes only the `G83` argument to zero at
`FUN_1000e0c0`. The runner restores the Stage 5 input and output on exit.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-candidate-lists.sh control \
  numeric-context-g83-can-anybody-help
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-candidate-lists.sh g83-zero \
  numeric-context-g83-can-anybody-help
```

The Stage 10 verifier checks 40 ordered abbreviation/context captures, 11
all-key H/T batch captures (55 H keys and 42 T keys across 87 context returns),
seven part-of-speech controls, five forced-pronunciation tree-input controls,
three address-guard captures, seven alias-branch traces, and 227
dictionary-metadata rows across ten observed Boolean patterns:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_stage10.py
```

The same stage's VTML controls map the original engine's 69 ordinary internal
phone bytes to documented `x-cmu` symbols:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_cmu_codebook.py
```

The controlled one-phone tree captures are summarized and checked for the
CMU alphabetical identity ordinals and five consonant onset-tree groups:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/analyze_stage10_tree_features.py
```

Stage 13 adds a host-side Intel PT collector for a single existing Stage 12
Wine probe. See [`stage13/README.md`](stage13/README.md) for prerequisites,
the gated attach sequence, and trace decoding. It leaves the runtime container's
non-root and capability-dropped settings unchanged.

Stage 14 contains the selected-key Lead 4 pilot; Stage 15 captures context,
punctuation, and case contrasts for all 427 H/T/C abbreviation keys. See
[`stage15/README.md`](stage15/README.md) and the curated
[Lead 4 results](../../../docs/reverse-engineering/lead4-abbreviation-context-results-2026-09-25.md).

Stage 16 probes file and buffer synthesis modes, selected errors,
`VT_GetTTSInfo_ENG`, configuration setters/getters, selected playback calls,
and the synchronous buffer length contract. See [`stage16/README.md`](stage16/README.md)
for the runners and [Lead 6 results](../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md)
for the captured behavior and limits.

The helper in
`tools/revkit/scripts/compare_tree3_runtime.py` rechecks Stage 5 lookup logs
against locally available voice trees.

To rerun original-engine traces, use the scripts in the stage directories and
the isolated Docker/Wine setup described in [`tools/revkit/README.md`](../README.md).
That rerun requires locally available `binary/`, `data-common/`, and
`data-paul/` inputs; those full vendor assets and generated environment state
are intentionally not included here.
