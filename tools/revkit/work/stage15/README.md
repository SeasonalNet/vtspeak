# Stage 15: full Lead 4 abbreviation context matrix

Stage 15 captures every key in the H/T/C abbreviation inventory. It has 55 H,
42 T, and 330 C keys; the matrix contains 1,750 delimited observations and
records all returned phone/context rows plus 24 bytes from each row's code
field. The JSON groups the captures by key. The TSV has one row per key and
context, including the exact input and ordered returned surfaces.

## Context design

- H: baseline title before `Smith`, title without the trailing period, an
  ordinary sentence, and an alternate-case title.
- T: baseline company suffix, comma-separated suffix, suffix without a period,
  ordinary sentence, and an alternate-case company suffix.
- C: neutral baseline, a likely-reading phrase, a comma-neighbor control, and
  an alternate-case neutral phrase.

The C likely-reading phrase is selected by a checked-in spelling heuristic:
50 state-code, 24 calendar-date, 20 calendar-weekday, 14 clock/zone, 13
address/building, and 209 generic numeric phrase probes. These are input
selection categories, not VoiceText semantic labels. Alternate-case variants
that collide with another key in the same family are marked in the manifest
(21 C and 5 T cases).

Each test is followed by a separate `Next.` sentence. The Stage 10 context
trace can combine or split context returns at boundaries that depend on the
abbreviation. The parser uses the returned `Next` row to associate each
complete returned sequence with the corresponding manifest entry, while
retaining all captured rows, flags, source indexes, and phone/context bytes.

The full per-key outputs are in [`lead4-key-matrix.json`](lead4-key-matrix.json)
and [`lead4-key-summary.tsv`](lead4-key-summary.tsv). The fixture order and
context inputs are in [`lead4-matrix-manifest.json`](lead4-matrix-manifest.json).

## Reproduction

From the repository root:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/generate_lead4_matrix.py
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage15/run-lead4-matrix.sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_lead4_matrix.py
```

Pass `H`, `T`, or `C` to the runner to rerun one set. It uses a dedicated
scratch directory and restores the Stage 5 input and output on exit. The
container reuses the existing Wine prefix and mounts all vendor inputs
read-only.

The curated findings and limits are in
[`lead4-abbreviation-context-results-2026-09-25.md`](../../../../docs/reverse-engineering/lead4-abbreviation-context-results-2026-09-25.md).
