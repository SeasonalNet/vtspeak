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
| [ ] | 4. Explain `.upm` and remaining index fields | Correlate per-unit `.upm` vectors and unknown index columns with the interpolation and timing paths. **Runtime checkpoint C:** inspect selected records at the reader and reconstruction boundaries to establish value ranges and effects. | 1–3 days |
| [ ] | 5. Validate `tree3` family | Check the parser against all 17 trees and identify each tree's inputs and output meaning from its callers. **Runtime checkpoint D:** trace representative lookups and compare selected leaves/values with runtime results. | 1–2 days |
| [ ] | 6. Map text and pronunciation resources | Inspect `data-common/dict-eng` files, including `engbi.tree3`, `atmt.tree3`, hash tables, and `txt2`; follow normalization, pronunciation lookup, and number handling. **Runtime checkpoint E:** compare a small set of controlled text inputs with the generated phone/context records and output. | 2–5 days |
| [ ] | 7. Recover candidate selection and scoring | Connect the five-byte class keys, feature views, `cepdist.tbl`, and scorer to candidate ranking and duration/prosody choices. **Runtime checkpoint F:** observe candidate inputs/ranks for selected contexts and reconcile them with the static scoring path. | 2–4 days |
| [ ] | 8. Trace synthesis end to end | Follow selected units through timing, pitch or other transforms, joins/blending, sample assembly, and WAV writing. Use a short reproducible corpus and cross-check the whole path against original-engine execution. | 3–6 days |

The remaining static and runtime investigation is estimated at about
**9–20 focused analyst-days**. Stages can overlap where their dependencies
allow. Stages 1, 2, and 3 are complete for the documented 2013 Paul package
scope. The remaining estimates cover analysis and decoder prototypes, not
the later clean-room engine implementation, compatibility work, or rights
review for voice assets.

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
cannot expose the relevant boundary. Preserve trace logs and reference WAVs
under the ignored `tools/revkit/work/` area, then summarize curated results in
the relevant documentation. Keep vendor binaries and voice/model files
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
