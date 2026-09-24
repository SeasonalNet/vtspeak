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
| Full-corpus DAT sample parity | The local stream inventory parsed all 580,474 records with zero decoded-sample-count mismatches. A normal-load Wine probe then decoded every indexed DAT payload through the DLL. Every returned PCM byte count and SHA-256 matched the independent Python decoder in bank/record order. | **Resolved for this local 2013 Paul package at the decoder output boundary.** The earlier 27 captured payloads also match byte-for-byte. The corpus-wide comparison uses cryptographic digests; it does not prove whole-synthesis or other-package parity. See the [corpus parity report](dat-corpus-parity-and-stage9-2026-09-24.md). |
| DAT mode 8 | A full framed-stream scan found no mode-8 frames in any of the 580,474 local records across `gen`, `num`, `etc`, and `alp`. | **Resolved as absent from this local package.** There is no real mode-8 payload to compare against the DLL; the static zero-frame branch remains unvalidated at runtime. The reproducible inventory is `PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/inventory_dat_modes.py`. |
| `.idx` span and cached UPM fields | For all 580,474 local records, bytes 4–7 match the two doubled side-period sums and bytes 16–18 match first/shared/last UPM bytes. | **Resolved structurally for local versioned indexes.** It does not label every feature column or establish older index variants. |
| UPM timing and joins | UPM bytes are widened and doubled onto the observed 16 kHz frame grid. Non-default pitch creates adjacent-period segment records; interpolation, edge windows, overlap weights, context eligibility, and output counts are mapped at the code/data-flow level. | **Closed for the declared target.** The mechanics and numeric timebase are established. Acoustic names for packed fields are absent from the local package; “8 kHz pitch period” remains a numeric interpretation rather than a vendor label. |
| `tree3` traversal | Parser structure and signed branch behavior are checked; 307 captured scalar/vector lookups across 17 Paul trees match runtime values. Single-phone contrasts establish that the identity values across all 39 CMU base symbols are their one-based alphabetical ranks; consonant manner groups select five different onset trees. | **Closed for traversal and observed inputs.** Phone identity and stress fields are labeled in these paths. Other selectors/outputs have numeric behavior, but no local schema provides context-slot names, physical units, family-name expansions, or full downstream prosody semantics. See the [controlled tree input labels](phone-symbol-codebook.md#controlled-tree-input-labels). |
| Text/pronunciation code labels | Resource framing, phone-ID expansion bytes, token/context layouts, typed TPP grammar, and representative ordered rules are documented. Controlled VTML `x-cmu` inputs map every internal phone byte `0x01`–`0x45` one-to-one to a documented CMU label. Single-phone contrasts map all 39 CMU base symbols to their one-based alphabetical identity values in tree inputs and distinguish five first-tree groups. Part-of-speech tags identify an upstream lexical-class constraint at token offset `+0x30` for documented values. All 427 H/T/C keys have one original-runtime context probe each; the C sweep found 93 returned surface sequences differing from the neutral input template, including lexical expansions and case changes. The `F`/`G` numeric suffix path parses signed decimal values and writes the narrowed byte to token field `+0x25`. A call/return trace resolves repeated `Apt` suppression to address-phrase span consumption: ` Apt 4 and Apt` is consumed as one 14-byte match; separated and comma controls are processed as independent spans. | **Observed mechanics closed for tested paths.** Numeric TPP class names, `AX`, dictionary metadata names, and outputs for untested per-key contexts are not encoded in the local corpus. The evidence does not support assigning those semantic labels. See the [phone codebook](phone-symbol-codebook.md) and [Stage 10 contrasts](abbreviation-context-contrasts-2026-09-24.md). |
| Token-result stride gap | `0x17c`–`0x54b` is not written by `FUN_1000d450`. A whole-DLL scalar scan found no direct offset reference in this interval; the one `0x17c` occurrence is stack allocation. The direct `0x554` references occur in `FUN_1000cf00` and `FUN_1000ea20`; decompilation shows they compare token surfaces and convert record fields. | **Direct-reference check closed for this DLL.** No direct reader or writer for the gap was found. Opaque indirect pointers and other modules are outside this binary-level statement. |
| Unit/model feature columns and class keys | Field widths, transformations, projections, scoring weights, branches, and runtime paths are mapped. The values serve as categorical selection/context features in the observed algorithms. | **Operational behavior closed.** The local data does not label the categories phonetically; assigning names would require an independent labeled dataset or vendor schema. |
| Whole-synthesis coverage | Default output, pitch/speed/volume controls, a short no-neighbor case, and two additional number/abbreviation inputs have sample-buffer-to-WAVE byte/count cross-checks. | **Feasibility validation closed at the documented boundaries.** This evidence does not claim every context-gate combination or corpus-wide whole-synthesis parity. Those are separate parity targets, not unanswered mechanics in this review. See the [Stage 9 report](dat-corpus-parity-and-stage9-2026-09-24.md#additional-synthesis-boundary-checks). |
| API and compatibility boundary | Public exported signatures and argument names are documented in `include/vt_eng.h`; the 2013 Paul `VT_TextToFile_ENG` path is mapped. The executable is PE32 and uses Windows ABI/runtime dependencies. | **Interface inventory closed for the selected target.** Full behavior for every selector, buffer/thread mode, error path, and other voice/package is not claimed; no compatibility target for those cases was selected in this review. |

## Closure boundary

Every item from the prior review now has a disposition for the documented
2013 M16 Paul feasibility target. Direct binary behavior, available runtime
effects, and checked-in format structures have been analyzed to the boundaries
listed above. A small set of semantic names remains unavailable from local
evidence: numeric TPP class names, `AX`, dictionary metadata labels, and most
tree/model feature labels. Assigning them would need an independently labeled
corpus or vendor schema. Compatibility outside this
voice/package, untested API modes, and whole-synthesis corpus parity are
separate targets and are not part of the current feasibility claim.

The current evidence supports a feasibility decision for the documented 2013
M16 Paul package. It does not by itself justify naming undocumented fields,
claiming whole-synthesis corpus parity, or declaring compatibility with another
VoiceText package. Vendor binaries and voice assets remain proprietary inputs.

The full local scan used
`PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/inventory_dat_modes.py`.
It parsed all 580,474 records, found only frame modes 0–3, and reported zero
DAT/UPM sample-count mismatches. Per-bank mode counts are emitted as JSON by
the script. The subsequent DLL-to-Python comparison established all 580,474
PCM byte counts and SHA-256 values at the decoder boundary; see the linked
corpus parity report. The 27 captured `gen` payloads remain the direct
byte-for-byte fixtures from Stage 2.
