# 2013 Paul front-end research batch

Date: 2026-11-14

## Objective

Close the remaining path from source text to the model and selection inputs
already consumed by the independent Go engine for the local 2013 M16 Paul
package. The batch should recover and validate producers, ordering, and state;
it should not repeat the established DAT, UPM, tree, unit-selection, or PCM
mechanics.

The completion target is a bounded native path that accepts the selected text
fixtures and produces the same meaningful intermediate rows, context
signatures, query modes, and selected unit IDs as the original DLL. Arbitrary
text and whole-engine compatibility remain later acceptance targets.

## Current baseline

The Go implementation already contains more than isolated front-end helpers:

- `text.LexiconFrontend` resolves the embedded dictionary, supported ASCII and
  numeric forms, pronunciation alternatives, source spans, and the captured
  VTML pause, mark, and substitution forms.
- The model parser code implements the outer dispatcher, primary and alternate
  parser loops, rollback, selected scanner branches, terminal-context
  classification, parser finalization, `0x554` token-row construction,
  `0x70` context-row projection, and parts of the normalization, exception,
  TPP, compound, suffix, and special-pronunciation cascades.
- The model-state code projects finalized `0x94` parser rows into counted
  `0x3c0` records at arena offset `+0x64c`, derives phone groups, runs the
  marker/boundary tree, and models the position-state workspace and its
  handoff toward selection.
- The selection packages implement the known whole-position and fallback
  queries, candidate expansion, scoring, pruning, transition costs, and
  backtracking when their model and workspace inputs are supplied.

The remaining problem is that several producers and live predicates are still
callbacks or caller inputs, and most of the composed parser/model path has
only static or synthetic validation. The provisional `pipeline.Pipeline`
therefore cannot yet supply a concrete Paul frontend, selector, and renderer
for arbitrary source text.

## Evidence boundary

Use the original read-only `binary/vt_pau.dll`, SHA-256
`200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce`,
with the matching local Paul M16 and shared English resources. Ghidra output is
pseudocode. Runtime captures establish behavior only for their exact DLL,
resources, options, and text fixtures.

Existing static work through Stages 60-186 should be treated as the map for
instrumentation. Do not decompile another broad function inventory. Pull a
new body or table only after a live fixture reaches it or an existing Go
callback cannot otherwise be resolved.

## What must be pulled

### 1. One reusable original-engine capture

Extend the existing isolated Wine/GDB harness instead of creating per-question
scripts. Each run must record:

- DLL and resource hashes, source bytes, text type, speaker, pitch, speed,
  volume, pause, dictionary selector, and parser mode;
- every checkpoint's function address, call ordinal, row count, source cursor,
  return value, and exact captured byte ranges;
- SHA-256 values for larger snapshots plus a decoded, reviewable row summary;
  and
- the final selected unit sequence and reference WAVE hash as canaries.

Capture only live counts and bounded used spans. Do not dump unrelated process
memory or proprietary banks into Git.

### 2. Source scanner and parser-row production

Instrument `FUN_1003d350`, the primary parser `FUN_1003d3d0`, the alternate
parser at `0x1003e070`, and `FUN_1005a350`. For each scanner call, retain the
input cursor, absolute source offset, mode and mode-pointer state, return
advance, and complete `0x44` result.

For each primary-handler attempt, record the handler address, result, parser
row count and cursor before and after, and the changed bytes in the current
`0x94` row. The known handler order is already recorded in Stage 172; use the
trace to identify which bodies actually matter for the fixture corpus. Pull
only those reached bodies and their directly required tables/callees.

This packet must resolve or preserve as explicitly unsupported:

- the remaining ordinary, numeric, punctuation, structural-marker, and
  non-ASCII branches of `FUN_1005a350`;
- the native index builders used by the scanner's post-token lookup;
- production of parser-row fields `+0x2c`, `+0x30`, and `+0x52` rather than
  supplying them as test inputs; and
- the primary cascade bodies still represented by injected handlers.

### 3. Token, pronunciation, context, and TPP rows

Capture entry and return state around `FUN_1000e2f0` and its live callees. The
required snapshots are:

- initialized parser rows and the rows after source parsing;
- used `0x554` token rows written through `FUN_1000d450`;
- the context-table count at model `+0x429a2` and used `0x70` rows beginning
  at model `+0x429ad`;
- pre/post row bytes for `FUN_10007520`, including the handler that accepted
  each row;
- `FUN_1000dfc0` TPP window inputs, selected records, and row updates; and
- the used TPP output at parser-state `+0x39ec`.

The trace must distinguish dictionary lookup, pronunciation selection,
normalization, exception, TPP, compound/contraction, suffix, and special-word
effects. A final row diff alone is insufficient because it cannot establish
native ordering when more than one handler could have matched.

Use these captures to replace the remaining `ParsePhoneRows`, `ProcessRows`,
and `ParseTPPRows` callback boundaries in `RunPaul2013ModelTextParser` for the
supported fixtures. Keep opaque numeric classes opaque; behavioral parity does
not require speculative linguistic names.

### 4. Finalizer, model-state, and selector handoff

Capture the state immediately before and after `FUN_1003e240`, then the inputs
and used outputs of `FUN_10016c90`, `FUN_10022dc0`, `FUN_10022970`, and
`FUN_100130e0`. Retain:

- finalized `0x94` rows, including source intervals, phone strings, per-token
  values, flags, and pointer fields;
- model-state count at `+2`, used `0x3c0` records from `+0x64c`, and final
  mode at `+0x4770a`;
- live workspace descriptors, cursors, boundary/value arrays, index tables,
  clamp inputs, and only the pointed prefixes actually consumed; and
- the resulting seven-byte selection signatures, query modes, whole/fallback
  class queries, candidate IDs, and chosen path.

This packet is intended to resolve current caller inputs such as per-token
state values, pointer ownership, event/range producers, descriptor and index
table producers, and workspace initialization. The final selector capture is
a boundary proof; candidate scoring and rendering do not need to be
re-researched.

## Fixture matrix

Start with existing reproducible inputs so new captures can be compared with
tracked evidence:

| Fixture family | Required behavior |
| --- | --- |
| Stage 1 short prose | Ordinary words, spaces, terminal period, unique pronunciations |
| Stage 6/9 number and time | Integer, decimal, unit, date/time, and multi-chunk rows |
| Stage 6/9 address and telephone | Honorific, address abbreviation, telephone, punctuation |
| Ambiguous words | Existing `record`, `dogged`, `lied`, `resume`, `separate`, and `live` branches |
| Proper-name/TPP | Existing component-pattern and G83-sensitive controls |
| Scanner boundaries | Newline terminal prefix, quotes, apostrophe, dot/hyphen joins, ordinal and grouped-number forms |
| Captured VTML | Existing pause, mark, and substitution spellings only |

Add a fixture only when it reaches a distinct unresolved producer or branch.
Do not begin with all abbreviation keys or an arbitrary text corpus.

## Packet order and deliverables

### F1: capture harness and ordinary baseline

- Produce one manifest-driven runner and one ordinary-prose capture.
- Define compact binary snapshots plus a deterministic decoded summary.
- Add a verifier that checks counts, extents, hashes, and expected checkpoint
  order.

### F2: scanner and primary-handler closure

- Run the scanner-boundary and numeric fixtures.
- Produce a live handler coverage table for the primary cascade.
- Extract and port only reached missing bodies and index producers.
- Add Go differential fixtures for scanner results and complete `0x94` rows.

### F3: model-context and TPP closure

- Run ordinary, number/address, ambiguous-word, and proper-name fixtures.
- Record handler-attributed diffs across token and context rows.
- Remove the three `FUN_1000e2f0` callback requirements for supported paths.
- Add row-for-row Go comparisons with explicit masks only for bytes proven
  untouched or pointer-valued.

### F4: state/workspace and selector closure

- Capture the finalizer, model-state records, workspace producers, and
  selection boundary for the same fixtures.
- Replace synthetic caller inputs where the native producer is now known.
- Compare context signatures, query modes, candidate lists, and selected unit
  IDs before exercising the existing renderer.

### F5: bounded end-to-end proof

- Run the Go path without injected linguistic callbacks for every fixture in
  the matrix.
- Require byte parity at the recovered intermediate fields and exact selected
  unit IDs; compare final PCM/WAVE where the existing renderer supports the
  path.
- Update the implementation plan and coverage notes with resolved and still
  unsupported branches.

Portable scripts, compact snapshots, manifests, and verifiers belong under a
new tracked `tools/revkit/work/` batch directory. Curated conclusions belong
in this document or a dated result note. Full memory dumps, Wine prefixes,
Ghidra projects, and vendor resources remain ignored.

## Acceptance criteria

The batch is complete when:

- every fixture has a reproducible original-engine capture with exact target
  identity and bounded snapshot extents;
- the Go implementation constructs matching meaningful parser, token,
  context, finalized-model, and selection rows without fixture-provided row
  controls or linguistic callbacks;
- native and Go row counts, source intervals, phone bytes, state/marker
  fields, context signatures, query modes, and selected unit IDs agree at the
  documented checkpoints;
- every compared byte is either equal or has a documented pointer/address or
  native-unwritten exclusion;
- unsupported source forms fail explicitly rather than silently entering an
  inferred path; and
- the existing Go unit suite, build, vet, and repository diff checks pass.

## Non-goals

- Re-recovering DAT/UPM decoding, tree formats, candidate scoring, waveform
  reconstruction, or WAVE serialization.
- Supporting the 2006 engine, another voice/package, or every public API mode.
- Assigning names to opaque TPP, feature, or state values without independent
  evidence.
- General Unicode, arbitrary VTML/SSML, malformed-input compatibility, or all
  abbreviation contexts in this batch.
- Establishing redistribution rights for vendor binaries or voice data.
