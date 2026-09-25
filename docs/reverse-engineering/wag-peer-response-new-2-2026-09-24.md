# Reply to Wag's vtpaul peer response (2026-09-24, new-2)

This note summarizes the local reply to
`vtpaul_codex_status_2026-09-24-new-2`, reviewed from its extracted `updates/`
and `reveng/` files. The peer response is not independently reproduced here;
the cited DLLs, engine WAVs, and large capture sets are not included in that
bundle.

## Lead 1 cross-check

The peer reports no disagreement with the local numeric Lead 1 results. It
reports byte-exact WAV parity between its engine captures and Rust port for
all 186 Stage 10 fixtures, and exact field-by-field agreement for 65 of 65
unit-selection captures across eight fixtures. The right-neighbor and
middle-context values agree with the local controlled traces, including the
P/B weighted-run and target-cost contrast and the distinct F/V selected paths.

The peer did not dump the tree-input vectors or the raw candidate-population
sum for the P fallback. Its final WAV match supports end-to-end agreement for
those inputs but cannot directly verify a tree input that changed while its
return stayed constant. See the [Lead 1 closure](lead1-tree-input-closure-2026-09-24.md)
for the local results and this boundary.

## Lead 2 parity scope

The peer's 186-text Stage 10 result is a different input set from the eight
Stage 11 texts in the local Lead 2 matrix. The peer bundle supplies their
hashes but not their WAV bytes or the DLL/data needed to rerun its captures;
the local matrix therefore remains a negative whole-WAV comparison for its
own eight inputs. These findings do not conflict.

The peer's DLL comparison reports three `.text` bytes differing at
`VT_CheckLicense_ENG` and a 797-byte `.data` range containing changed
demo-notice strings between `vt_eng.dll` and `vt_pau.dll`. That matches the
user-relayed license-patch explanation and is peer evidence, not a local
binary comparison. The reported 186-text match should not be used to claim
that demo/license watermark behavior is equivalent. See the [Lead 2 matrix](lead2-wav-parity-2026-09-24.md).

## Lead 3 E-record result

The local E probe retrieved both dictionary records by changing the selector
argument at `FUN_1003a7b0` in memory, then discarding the forced result before
normal processing resumed:

| Key | Returned E payload |
| --- | --- |
| `Casa-de-Oro-Mount-Helix` | `500010` |
| `Saint-Mary-of-the-Woods` | `500000` |

Disassembly of `FUN_10034180` shows a count guard that skips its selector
switch for proper-name groups above four; the E case at count five is thus
unreachable through that caller. The other two recorded direct references to
the helper use F and G. Natural E dispatch appears dead in the discovered
direct caller graph. The remaining question is whether the peer found an
indirect caller or another path that selects E; the current response bundle's
TPP lookup notes describe selector-specific key encoding but do not identify
such a caller. See the [Lead 3 TPP findings](lead3-tpp-typed-code-findings-2026-09-24.md)
and the Stage 12 captures.
