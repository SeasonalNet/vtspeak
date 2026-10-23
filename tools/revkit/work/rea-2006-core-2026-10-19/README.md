# REA: 2006 VoiceText core, 2026-10-19

Curated explanation:
[loading through PCM](../../../../docs/reverse-engineering/rea-2006-concatenative-core-2026-10-19.md).

Original input:
`corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/lib/vt_eng.dll`,
relative to the parent work directory. SHA-256:
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.
Provider: Ghidra 12.1.4 through REA. No vendor code was executed or modified.

## Portable artifacts

| File | Contents |
| --- | --- |
| `loading-prosody-evidence.json.gz` | Complete 51-record REA export: loaders, scalar/vector trees, prosody, core orchestration, PCM/WAVE output |
| `loading-prosody-ledger.md` | Function-to-Evidence IDs for that export |
| `selection-evidence.json.gz` | Complete 29-record export: class search, candidate scoring, DP, pruning, backtracking, constants |
| `selection-notes.md` | Selection call paths, fields, formulas, and limits |
| `audio-evidence.json.gz` | Complete 26-record export: DAT/UPM, decoder primitives/handler bytes, timeline, direct/grain synthesis, postfilter |
| `audio-notes.md` | Audio call paths, fields, formulas, and limits |
| `descriptor-evidence.json.gz` | Complete 16-record follow-up: global unit IDs, index-bank lookup, compact/interleaved record readers, original record bytes, and endianness |
| `pitch-smoothing-evidence.json.gz` | Assembly follow-up for contour smoothing, median, and square-root helpers |
| `pitch-smoothing-notes.md` | Recovered contour smoothing math and assembly Evidence IDs |
| `asset-byte-observations.json` | Original index/table hashes, widths, counts, exact sample offsets, and database list |
| `vector-tree-byte-observations.json` | Four original Kate vector-tree hashes, root fields, node/leaf counts, and EOF consumption |
| `evidence-manifest.json` | Gzip/JSON hashes, sizes, record counts, and target identity for all five exports |

Compressed exports are deterministic gzip (`gzip -n`) of the complete REA JSON;
no fields, native value-flow records, assembly, or limitations were removed.
Some provider records themselves report incomplete analysis; see the notes,
especially for the DAT decoder and pitch smoothing. The machine unknown
registry was not populated and must not be read as proof of full recovery.

## Reproduce or inspect

1. Keep the original extracted MSI inputs at their existing paths. Check the
   DLL hash above before using addresses in these notes.
2. Open that original DLL with REA `open_binary`, provider `ghidra`.
3. Use `analyze_function` for the addresses in the ledgers. Calls are static;
   they do not authorize or imply running the engine.
4. Export a new evidence bundle to a fresh path. Close the native session.
5. Compare claims against pseudocode **and** assembly when types or function
   bodies are incomplete. REA can omit computed-jump handlers from a recovered
   body; the audio pass retains the original handler/table bytes separately.

For portable bundle inspection, decompress one export to a new temporary path
with `gzip -dc ...json.gz > /tmp/<chosen-name>.json`, then use REA's
`inspect-evidence-bundle` with that path and `detail: summary`. Focus by exact
Evidence ID, canonical procedure address, or operation rather than dumping the
entire native value-flow inventory.

Asset observations are read-only. Index extents use the existing
`tools/revkit/scripts/inspect_legacy_unit_idx.py:parse_index_data`; feature
tail starts at `records_offset + 19*N` and spans `20*N` bytes. Vector-tree
observations recursively follow `0x10001a30` / `0x10001bf0`: read feature,
operation, LE signed threshold, output-width byte, that many LE signed words,
membership-count byte and words, then two trailer bytes; recurse into two
children when the second trailer is nonzero. Count nodes/leaves and compare
the final cursor with file size. This uses the vector format, not `tree2.py`'s
scalar format. All four observed files have output width 12 at every node.

No implementation or parser was changed. Existing runtime evidence referenced
in the notes was not rerun. Fresh engine parity, malformed-input coverage, and
all-voice format coverage are not claimed.
