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
- `stage3/` through `stage8/` contain the controlled text inputs, trace
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
- The repository's `binary/`, `data-common/`, and `data-paul/` source inputs
  remain outside Git as before.

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

The helper in
`tools/revkit/scripts/compare_tree3_runtime.py` rechecks Stage 5 lookup logs
against locally available voice trees.

To rerun original-engine traces, use the scripts in the stage directories and
the isolated Docker/Wine setup described in [`tools/revkit/README.md`](../README.md).
That rerun requires locally available `binary/`, `data-common/`, and
`data-paul/` inputs; those full vendor assets and generated environment state
are intentionally not included here.
