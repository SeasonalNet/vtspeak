# Deferred follow-up: normalized voice-data export

**Status:** Deferred until the initial native reimplementation can load a
matched voice package and produce speech through an end-to-end path. This is a
follow-up for review and planning, not part of the initial engine scope.

## Purpose

After that milestone, assess a reproducible exporter that converts validated
vendor package structures into a more convenient working representation for
the concatenative engine. Keep the vendor inputs read-only and preserve the
source package identity and unit IDs throughout the export.

Potential outputs include:

- Packed PCM16 unit banks with an index mapping each source bank/unit ID to
  its waveform offset and length. Individual WAV files could be generated on
  demand for listening to selected units; writing one file per unit is likely
  inconvenient for the full corpus.
- Structured per-unit timing data derived from the paired UPM vectors, with
  the original bytes retained and any physical interpretation clearly marked
  as observed or inferred.
- Normalized index and tree records that preserve source ordering, raw feature
  bytes, package generation, and opaque fields or suffixes whose meanings are
  unresolved.
- Queryable tables for pronunciation and other shared resources, again
  retaining original IDs and unclassified values.

The runtime representation and the analysis representation need not be the
same. A packed waveform bank may suit random access in the engine, while
WAVs, JSON, or a small database may be more convenient for inspection. Choose
formats only after measuring corpus size, load time, and the access patterns
of the working engine. Decoding compressed DAT streams to PCM is expected to
increase storage size; the tradeoff is simpler direct waveform access.

## What this could help with

- Avoid repeatedly parsing package-specific index layouts or decoding a unit
  every time an analysis tool needs to inspect it.
- Give the Go engine a consistent in-memory or on-disk interface while keeping
  generation-specific importers for the actual 2005/2006-era and 2013 package
  layouts.
- Make units easy to browse, audition, plot, and compare by their original
  IDs, banks, timing vectors, and feature bytes.
- Isolate and test unit selection, waveform decoding, timing adjustment, and
  joins against smaller, directly inspectable records.

## Boundaries

Export does not supply the synthesis algorithm. Text analysis, candidate
selection, prosody, duration and pitch changes, waveform reconstruction, and
context joins still need to be implemented and validated. A normalized
container also cannot infer unknown feature meanings or resolve a mismatch
between a newer reader and older package data; preserve such differences and
unknown bytes rather than smoothing them over.

An exported corpus is derived vendor data. Conversion does not change its
licensing status. Treat both the source assets and generated exports as
proprietary, and keep generated corpus copies out of version control unless a
separately reviewed policy explicitly permits a small fixture.

## Revisit and validation criteria

When the initial engine has a working end-to-end path for an explicitly
matched package, make a bounded export plan. Before calling any generation's
export validated:

1. Check its reader against multiple packages/files and preserve all source
   fields, unit IDs, ordering, and byte offsets needed to trace results back.
2. Compare decoded waveforms with the corresponding original engine at the
   decoded-unit boundary, over the intended corpus scope. Keep decoder
   validation separate from whole-synthesis parity.
3. Verify UPM/timing exports against source vectors and observed synthesis
   behavior; retain raw vectors when a derived interpretation is uncertain.
4. Demonstrate that importing the exported representation gives the
   reimplementation the same inputs and byte-level results as reading the
   source package directly for the checked fixtures.
5. Record storage and load-time costs, package/version scope, unresolved
   fields, and the fact that the exports remain proprietary.

Current evidence is uneven across generations. The local 2013 M16 Paul DAT
decoder matches the original DLL's decoded output across all 580,474 indexed
payloads by PCM byte count and SHA-256. Older package indexes and DAT/UPM spans
have structural and sampled consistency checks, while opaque tree suffixes
and other generation-specific behavior remain. Use the [format findings](reverse-engineering/voice-engine-and-model-formats.md)
and [package generation comparison](reverse-engineering/voice-package-generation-comparison.md)
as the starting evidence; do not treat one generation's exporter as proof
that the others are understood.
