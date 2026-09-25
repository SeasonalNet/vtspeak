# Lead 1 closure: controlled tree and context inputs (2026-09-24)

**Status:** Complete for the bounded result in the [lead tracker](follow-up-investigation-leads-2026-09-24.md#1-label-more-decision-tree-inputs-with-controlled-contrasts). The target is the local 2013 Paul M16 package, using `voicetext_paul.exe` and `vt_pau.dll` in Wine. Their SHA-256 values are `b591add774bf7a50d7f5f36f637e1f6d762be6d9a076d98379c7ae6310707ee9` and `200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce`, respectively. The vendor binary and model inputs were read only.

Lead 1 asked for controlled mappings from phone, stress, neighboring-phone, tag, or boundary changes into tree inputs and downstream results. Its stated useful result is a verified mapping for at least one input field, a measured downstream effect, and explicit counterexamples or limits. The evidence now supplies several such mappings and follows one neighbor effect through class lookup, candidate scoring, backtracking, and timeline construction. It does not provide a full semantic schema for every packed field.

## Directly observed mappings

| Control and evidence | Input field or operation | Observed result | Boundary |
| --- | --- | --- | --- |
| 24 consonants and 15 vowels forced through one-token `x-cmu` fixtures; [codebook](phone-symbol-codebook.md#controlled-tree-input-labels) | Printed tree input phone-identity values are the one-based alphabetical rank of the 39 CMU base symbols. Onset-scalar position 0 and paired-scalar positions 0/1 carry tested phone ordinals. | P=`27`, B=`7`, T=`31`; consonants select five observed onset-tree groups. | This is a runtime mapping for these forced symbols, not a vendor name for every tree feature. |
| `P AH0/AH1/AH2` controls; [codebook](phone-symbol-codebook.md#controlled-tree-input-labels) | Paired-scalar position 3 and the `0x14bf4f8` scalar / `0x14bf568` vector position 2 track the forced AH stress value. | Values change `0/1/2` in order. | The later vector's position 11 follows a preceding scalar result; it is not an independent stress code. |
| Fixed `P AH0` with left or right P/B; [neighbor report](tree-context-neighbor-contrasts-2026-09-24.md#fixed-p-ah0-contrasts) | Left change appears as the phone ordinal in printed positions 0 and 1 of the first two subsequent scalar calls; right change appears in positions 0 and 1 of calls three and four. | Left P→B changes the first two returns `958,755`→`1040,1470`; right P→B changes call three `1159`→`1240`, while call four stays `2341` despite its changed input. | A changed input does not guarantee a changed tree return. The same orientation pattern was repeated with fixed `T AH0`; the numeric tree outputs have no established physical units. |
| Isolated `record` with documented part-of-speech tags; [codebook](phone-symbol-codebook.md#controlled-part-of-speech-field-contrast) | Upstream token byte `+0x30` responds to the tag. | Noun/modifier and verb/interjection pairs have different status bytes but identical phone sequences and identical ordered duration/pitch tree inputs within each pair. | This confirms an upstream lexical constraint. It does not show a separate tree feature encoding those tag distinctions after pronunciation is fixed. |
| Fixed `T AH0`, right P/B, with and without an intervening period; [neighbor report](tree-context-neighbor-contrasts-2026-09-24.md#repeat-with-fixed-t-ah0-and-a-sentence-boundary) | The period separates the observed dynamic-programming spans. | The preceding span has identical selected IDs and timeline sample counts (`3837,157573,157573`; `1104,956,1370`); the later span differs. | This establishes the tested period boundary, not all punctuation, prosodic, or word boundaries. |

The phone identity mapping is independently checked at two representations: the upstream forced internal phone bytes in the returned phone/context rows, and the duration/pitch tree input ordinals. The internal byte value for P is `53`, whereas its tree ordinal is `27`; those are distinct encodings.

## Measured downstream effect of a right-neighbor change

With a fixed `T AH0` probe and a following P or B phone, the fixed probe's first two tree vectors and returns remain equal. Its first timeline call nevertheless selects different unit IDs and sample counts. The [neighbor report](tree-context-neighbor-contrasts-2026-09-24.md#candidate-scoring-continuation) follows the intervening selection path in the local DLL:

1. At `FUN_10018770`, the context-1 seven-byte candidate signature is `[90,57,7,X,7,97,0]`; byte `+3` is the varied internal phone code (`53` for P, `19` for B, `32` for F, `65` for V). The six-byte context row at this entry is the same in the four right-phone controls.
2. `FUN_10016ea0` maps tested signatures to five-byte class keys. `FUN_10019450` searches the loaded sorted class table; the successful tested whole-phone keys yield class `48721` for P and classes `48656`/`48655` for B. A read-only scan found no duplicate five-byte keys in the 61,566-record local table.
3. At `FUN_10024060`, model flag `0` makes the tested whole-phone pass accept a `FUN_10023060` class-population sum of at least 10. The context-1 sums are 1 for P and 14 for B; F/V emit no exact whole-phone class. A zero return enters `FUN_100242a0`, which sets context-row byte `+0x04` to `1` for the fallback row (or `2` for its copied second row). P/F/V take that path; B does not in this right-neighbor control.
4. For context-0 candidate `3837`, `FUN_100230a0` reads the following row's byte `+0x04` and stores a weighted run value at candidate `+0x10`. The P/B value is `0/1`, the local scorer scale is `1/0.5`, and its cost is `0.833333/0.416667`. The transition and backtracker choose different predecessor paths and selected units, which reach different timeline rows.

This mechanism is supported by the raw GDB captures and local disassembly cited in the [neighbor report](tree-context-neighbor-contrasts-2026-09-24.md). Wag's independent port calls row byte `+0x04` the whole/left-half/right-half mode and candidate `+0x10` the combined weighted run. Those names are useful operational interpretations; the local trace directly verifies the tested branches and values.

The contrasts bound the rule. Changing F to V leaves candidate `3837`'s context-0 `+0x10` and local cost unchanged, yet changes the backtracked path. Moving P/B to the left of T reverses the tested context-1 whole-phone return. In three-phone B/T/P/B/F/V controls, context 4 keeps the same six-byte row while P/F/V return zero and enter half-key searches and B returns one. Thus neither P/B identity nor a changed `+0x10` alone is a general rule for every position.

## Peer cross-check (Wag's Claude agent, 2026-09-24)

The peer's third reply reports no disagreement with the local numeric Lead 1
results. For the 186 Stage 10 fixture texts, it reports byte-exact WAV parity
between its captured engine output and its Rust port. It also reports 65/65
unit-selection capture comparisons exact across eight fixtures, including the
eight new right/middle context probes. The `T AH0` right-neighbor table agrees
with the local observations for B's 14 whole-phone candidates, P/F/V half-phone
fallbacks, candidate `3837`'s P/B `wrun` and target-cost values, and the
different F/V selected paths. The peer explicitly did not capture the raw
population sum for the P fallback or the tree input vectors for these
fixtures; the final WAV parity is supporting evidence but does not directly
verify a changed tree input with an unchanged tree return.

These are peer-reported results from `vt_eng.dll` captures and are not locally
reproduced here. The bundle provides hashes for the 186 WAVs and the eight
unit-selection JSON captures, but omits the engine WAVs and DLL/data needed to
rerun the capture set. The peer's prior DLL comparison also reports a
license-check code difference and demo-notice string differences between
`vt_eng.dll` and local `vt_pau.dll`; those reported differences are kept
separate from the field-level unit-selection agreement above.

## Completion boundary and reproduction

Lead 1 is closed on its stated useful result: phone identity and stress have direct input mappings; a neighboring phone has both a direct tree-input mapping and a traced selection/timeline effect; the POS and F/V controls provide counterexamples; and the period control bounds one phrase-break case. Numeric model features, untested boundaries, alternate dictionary pronunciations, other voice packages, and byte-exact full-WAV parity remain unclaimed. Whole-WAV parity belongs to lead 2.

Fixtures, GDB scripts, and reproduction commands are in the [Stage 10 notes](../../tools/revkit/work/stage10/README.md). The detailed [neighbor report](tree-context-neighbor-contrasts-2026-09-24.md) links each operation to its capture and address. All four recent three-phone first-pass and half-key runs exited normally; the runner restored the Stage 5 input and output byte-for-byte. No vendor input was edited.
