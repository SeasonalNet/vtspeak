# VoiceText reverse-engineering roadmap

This is the living work list for understanding the VoiceText engine and its
Paul voice data well enough to judge a native POSIX reimplementation. Update
the status and evidence links here as work progresses. Effort figures are
rough focused analyst-days for the remaining investigation, not calendar
deadlines. Runtime setup, missing reference outputs, or undocumented format
details may change them.

## Working method: static analysis checked against runtime behavior

Use a hybrid loop throughout the investigation:

1. **Form a bounded hypothesis statically.** Record the binary or data file,
   function/address or byte offset, observed reads/writes and calls, and what
   remains inferred. Use Ghidra pseudocode as a guide, not original source.
2. **Choose a runtime question that can distinguish explanations.** Examples
   include which index fields are read for one unit, what byte range is fetched
   from a bank, or how a known input changes a sample buffer.
3. **Observe the original program or DLL in a controlled run.** Capture only
   the necessary call, file, memory, or output evidence. Use a known input and
   preserve the input assets. Do not alter the license check or verification
   data as part of behavior tracing.
4. **Compare the observation with the static model.** Record agreement,
   mismatch, and the next question. If behavior disagrees with the current
   interpretation, revise the static model before building on it.
5. **Repeat at the next boundary.** Keep static mapping and runtime checks
   interleaved; do not postpone all execution validation until the end.

For every material finding, label it as **observed statically**,
**observed at runtime**, **cross-checked**, or **inferred**. A parser or codec
should not be called recovered until it has been checked against multiple
files or known engine output, consistent with [the repository evidence
rules](../../AGENTS.md).

## Progress and sequence

| Status | Stage | Main work and evidence checkpoint | Estimate |
| --- | --- | --- | ---: |
| Done | Baseline static map | API/load/output path, `tree3` reader shape, four versioned Paul `.idx` layouts, `.dat` predictor outline, and `.upm` processing path. See [current findings](voice-engine-and-model-formats.md). | Complete |
| Done | 1. Map unit records to payload spans | The read-only inspector finds contiguous DAT and UPM spans through EOF in all four banks. Runtime reads confirm DAT fields and both single-side and combined UPM lengths; see [Stage 1 findings](voice-engine-and-model-formats.md#stage-1-unit-to-payload-span-mapping). | Complete |
| Done | 2. Recover `.dat` decoding | Standalone prototype models the four-block mean history and predictor modes. Runtime checkpoint B: 27 unique `gen` payloads match captured DLL PCM byte-for-byte across modes 0–3; 16 records across all four banks match UPM-derived sample counts. Mode 8 is statically mapped but absent from the captured run. See [Stage 2 findings](voice-engine-and-model-formats.md#stage-2-dat-decoding). | Complete; full-corpus parity remains follow-up |
| Done | 3. Validate decode and synthesis observations | Controlled known-text run produced a stable 16 kHz mono 16-bit PCM WAV. Two captured sample-buffer returns concatenate byte-for-byte to its full WAVE data chunk; Stage 2 decoder matches remain consistent for this input. See [Stage 3 findings](voice-engine-and-model-formats.md#stage-3-decode-to-wave-boundary-check). | Complete; broader synthesis coverage remains follow-up |
| Done | 4. Explain `.upm` and remaining index fields | All 580,474 local records cross-check bytes 4/6 as doubled side-period sums and bytes 16–18 as first/shared/last UPM values. Runtime checkpoint C traces two selected vectors through segment generation and reconstruction under default and non-default pitch. See [Stage 4 findings](voice-engine-and-model-formats.md#stage-4-upm-and-index-field-semantics). Feature columns and full-corpus PCM parity remain open. | Complete |
| Done | 5. Validate `tree3` family | Parser validates node/list/output extents and references for all 17 Paul duration/pitch trees. Caller mapping identifies nine scalar duration trees and four scalar/12-value pitch pairs. Runtime checkpoint D compares 307 lookups across all 17 trees from three controlled phrases with no mismatches. See [Stage 5 findings](voice-engine-and-model-formats.md#stage-5-tree3-parser-and-caller-behavior). | Complete |
| Done | 6. Map text and pronunciation resources | All nine `.txt2` row schemas/callsites and all 427 H/T/C abbreviation keys/tags are mapped. `exceptdict` is grouped by one through four hyphen-separated components. Both hash families now have exact record framing; all 31,550 TPP keys decode and round-trip, with typed payload shapes and their single/multi-token callers documented. The embedded phone payload and the complete `0x554` token-result / `0x70` phone-context layouts are mapped. Cardinal, decimal, ordinal, four-digit year, case-match, address-alias, and ordered abbreviation/context paths are documented with runtime examples. **Runtime checkpoint E:** controlled Wine/GDB captures cover ordinary text, numbers, dates/times, addresses, and abbreviations, including token surfaces and phone/context rows. TPP code meanings and the full per-key/context output matrix remain interpretation work, not blockers to the resource-format map. See [Stage 6 findings](voice-engine-and-model-formats.md#stage-6-text-and-pronunciation-resources), [abbreviation inventory](abbreviation-table-inventory.md), [phone-symbol codebook](phone-symbol-codebook.md), and the [TPP key decoder](../../tools/revkit/scripts/inspect_tpp_dictionary.py). | Complete |
| Done | 7. Recover candidate selection and scoring | Class keys and view fields are named by data flow; weighted mismatch scoring, unit and adjacent-context costs, the 30/10,000 and 10/10,000 limits, divisor-4 pruning, flag retention, and backtracking are mapped. Three varied Wine/GDB runs exercise 2–77 contexts, different candidate sizes, and flag bypass. A controlled GDB budget override exercises ranker scoring and class-boundary truncation. A separate 41-context trace matches every selected ID to the `FUN_1001b200` handoff; destination record offsets are statically mapped. See [Stage 7 findings](voice-engine-and-model-formats.md#stage-7-candidate-selection-and-scoring). | Complete |
| Done | 8. Trace synthesis end to end | Maps 52-byte timeline rows, speed/pitch/volume effects, UPM-period reconstruction, left/right context eligibility and neighbor selections, PCM tail/block assembly, and WAVE writes. Four controlled settings plus a short no-neighbor case cross-check output counts; default and pitch-120 WAVs match prior references byte-for-byte. Remaining bounded limits are listed in the [Stage 8 findings](voice-engine-and-model-formats.md#stage-8-synthesis-end-to-end-complete-for-feasibility-scope). | Complete for feasibility scope |

Stages 1 through 8 are complete for the documented 2013 M16 Paul package
scope and support an implementation feasibility decision. Full-corpus parity,
other voices or package versions, and the upstream meanings of packed
conversion fields remain separate follow-up work; see the Stage 8 findings for
their consequences and evidence boundaries.

## Runtime analysis entry criteria

Runtime-backed analysis begins as soon as Stage 1 has identified a small,
traceable set of unit records and a concrete question that static analysis
cannot settle. It does not wait for every static pass to finish. Before each
run, record:

- the exact executable/DLL and voice/model input set;
- the specific function, file access, or buffer behavior being observed;
- the controlled input and expected evidence;
- how the observation will distinguish the competing explanations; and
- any environmental constraints that affect interpretation.

Use the least intrusive observation that answers the question. Start with
process/file/API traces; use debugger or memory instrumentation when those
cannot expose the relevant boundary. Preserve portable trace logs, fixtures,
and reference WAVs under tracked `tools/revkit/work/`, then summarize curated
results in the relevant documentation. Keep vendor binaries and voice/model files
read-only. Do not inspect, modify, or bypass license verification material as
part of this work.

## Completion criteria for the analysis phase

The analysis phase is ready for an implementation feasibility decision when:

- the index-to-payload mapping is validated against runtime reads;
- a `.dat` decoder reproduces selected original-engine unit samples;
- `.upm` and `tree3` behavior is described with cross-checked examples;
- text/pronunciation, candidate selection, and synthesis stages have a
  supported end-to-end call/data path;
- remaining unknowns and their consequences for a POSIX port are listed; and
- findings clearly separate observations from inferences and vendor assets
  remain treated as proprietary.

At that point, create a separate implementation plan with its own scope,
licensing assumptions, and acceptance criteria. Completing this roadmap does
not itself establish reuse or redistribution rights for the vendor voices.
