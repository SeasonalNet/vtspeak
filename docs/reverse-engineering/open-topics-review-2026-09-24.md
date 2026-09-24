# Open reverse-engineering topics review (2026-09-24)

This review reconciles the follow-ups left by Stages 1–8 with the checked-in
runtime captures, parsers, and Ghidra pseudocode. It covers the local 2013 M16
Paul package. A question is marked resolved only where the current evidence
answers it; descriptive confidence is not extended to other voices or package
versions.

## Dispositions

| Topic | Current answer | Status and evidence boundary |
| --- | --- | --- |
| DAT `nmean` history | The decoder keeps four signed 32-bit rounded frame means. The DLL sums all four before mode dispatch, uses signed rounding, compensates for the output shift, and then shifts in the newly rounded, shifted frame mean. The history starts as four zeros in the decoder invocation; the prototype resets it for each payload. | **Resolved for the analyzed DLL.** `FUN_10001b30` pseudocode and Stage 2 describe the four-entry state and 32 captured calls / 27 unique payloads with exact PCM matches, including repeated mode-0 frames. Those captures establish compatibility for the captured cases. They do not document every possible state sequence or a distinct-mean adversarial vector. See [Stage 2](voice-engine-and-model-formats.md#stage-2-dat-decoding) and the [peer-review reconciliation](wag-peer-review-followup-2026-09-23.md#resolution-nmean4). |
| Full-corpus DAT sample parity | Local index/UPM scans and 16 cross-bank length checks do not compare every decoded sample. Wag's separate 580,474-payload result is external corroboration. | **Still open locally.** Need a reproducible comparison boundary and results tied to this checkout's decoder and DLL; do not describe external parity as independently reproduced here. |
| DAT mode 8 | The static decoder path emits zero-valued frames for this control. No captured local runtime payload has exercised it. | **Still open at runtime.** First inventory all local streams for mode 8; if present, capture one original-DLL decode. If absent, retain as a static-only branch for this package. |
| `.idx` span and cached UPM fields | For all 580,474 local records, bytes 4–7 match the two doubled side-period sums and bytes 16–18 match first/shared/last UPM bytes. | **Resolved structurally for local versioned indexes.** It does not label every feature column or establish older index variants. |
| UPM timing and joins | UPM bytes are widened and doubled onto the observed 16 kHz frame grid. Non-default pitch creates adjacent-period segment records; interpolation, edge windows, overlap weights, context eligibility, and output counts are mapped at the code/data-flow level. | **Mechanics mapped; model naming remains open.** “8 kHz pitch period” is a numeric interpretation. The exact acoustic meaning of all packed fields and complete behavior across all flag combinations are not established. |
| `tree3` traversal | Parser structure and signed branch behavior are checked; 307 captured scalar/vector lookups across 17 Paul trees match runtime values. | **Traversal resolved; linguistic meaning open.** Duration/pitch roles are supported by callers, while exact feature names, physical units, family-name expansions, and downstream prosody interpretation are not recovered. |
| Text/pronunciation code labels | Resource framing, phone-ID expansion bytes, token/context layouts, typed TPP grammar, and representative ordered rules are documented. | **Formats resolved; interpretation open.** Compact phone labels, dictionary metadata and context flags, TPP atoms `A`–`G`/`AX`, and some abbreviation outputs require additional data/caller or controlled context evidence. Do not infer these names from corpus frequency alone. |
| Token-result stride gap | `0x17c`–`0x54b` is not written by the traced builder or read by the traced converter/context path. | **Bounded answer only.** The gap is unreferenced in the traced path; readers elsewhere in the DLL have not been ruled out. A full xref/reader search is still needed before calling it unused globally. |
| Unit/model feature columns and class keys | Field widths, transformations, projections, scoring weights, and branches are mapped, with runtime paths exercised. | **Data flow resolved; semantic labels open.** The values act as categorical selection/context features. Their phonetic names need a controlled mapping to known phone/context records or vendor documentation. |
| Whole-synthesis coverage | Default output, pitch/speed/volume controls, and a short no-neighbor case have byte/count cross-checks. | **Feasibility scope complete; coverage open.** Broader text/settings, every context-gate combination, and corpus-wide parity remain unvalidated. |
| API and compatibility boundary | The main exported calls and 2013 Paul call path are mapped, with remaining `param_N` names. The executable is PE32 and uses Windows ABI/runtime dependencies. | **Implementation planning input, not a completed compatibility contract.** Resolve parameter semantics and other package/voice differences only for a separately scoped port target. |

## What to revisit next

1. **Close local decoder coverage:** inventory all four DAT banks for mode 8,
   then define a practical reproducible full-corpus PCM check. Keep the DLL
   decoder output boundary distinct from UPM/prosody and final WAV comparisons.
2. **Resolve semantic labels only with discriminating evidence:** select
   representative phone/context records and trace how candidate labels affect
   tree inputs, TPP transformations, and final phone rows. Preserve raw values
   alongside any proposed names.
3. **Bound remaining text-rule questions:** decide whether a per-key/context
   abbreviation matrix is needed; if so, generate it as a controlled probe
   artifact with input, token/context rows, and output WAV references.
4. **Expand synthesis validation only for a stated target:** choose additional
   texts/settings or a second package, and compare timeline rows and PCM at
   the same boundaries already used in Stages 3 and 8.

The current evidence supports a feasibility decision for the documented 2013
M16 Paul package. It does not by itself justify naming undocumented fields,
claiming complete corpus parity, or declaring compatibility with another
VoiceText package. Vendor binaries and voice assets remain proprietary inputs.
