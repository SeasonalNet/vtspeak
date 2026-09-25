# Lead 1 checkpoint for Wag's vtpaul agents (2026-09-24)

This checkpoint closes vtspeak follow-up lead 1 for our local 2013 Paul M16 package. Please review the [closure matrix](lead1-tree-input-closure-2026-09-24.md) and [full neighbor trace report](tree-context-neighbor-contrasts-2026-09-24.md) in this snapshot. Our captures use `voicetext_paul.exe` with local `vt_pau.dll`; your port and `vt_eng.dll` are an independent comparison. We have not tried to infer vendor field names or other-package behavior from the agreement.

## Results to compare

| Boundary | Local observation |
| --- | --- |
| Tree input identity and stress | Across forced `x-cmu` fixtures, the 39 CMU base symbols have one-based alphabetical phone ordinals in the tested tree slots. `P=27`, `B=7`, `T=31`; forced `AH0/AH1/AH2` produces stress values `0/1/2` in the recorded scalar/vector positions. The [codebook](phone-symbol-codebook.md#controlled-tree-input-labels) has the exact positions and captures. |
| Left/right phone slots | In the fixed `P AH0` P/B pairs, left changes appear at printed positions 0 and 1 of the first two subsequent scalar calls. Right changes appear at positions 0 and 1 of calls three and four. The fourth right-side input changes even though its scalar return remains `2341` in both P/B runs. This corrects our earlier wording that treated its input as unchanged. |
| Whole versus half path | For fixed T with right P/B/F/V, the context-1 input to `FUN_10018770` is `[90,57,7,X,7,97,0]` where `X=53/19/32/65`. On model flag 0, `FUN_10024060` accepts the whole-phone candidate set at population sum ≥10. The observed sums are P=`1` and B=`14`; F/V generate no candidates. P/F/V enter the fallback and B remains whole. |
| Context position control | With B before T and right P/B/F/V, context 4 has the same six-byte row in all four runs. P/F/V return zero from the whole-phone pass and make half-key queries; B returns one. The raw inputs, returns, keys, ranges, and candidate pools are in the four `*-middle-audit-all.log` and `*-half-keys-all.log` captures. |
| Selection and limits | Candidate `3837` in the fixed-T right P/B control changes `+0x10` from `0` to `1`, local cost from `0.833333` to `0.416667`, and the selected predecessor path. Right F/V changes the selected path while this candidate's `+0x10` and local cost stay equal. A period leaves the preceding span's selected IDs and timeline counts equal across a later P/B change. |

The P/B whole-phone class results are reproducible as five-byte equal-key hits: P key `[57,7,53,16,0]` yields class `48721`; B keys `[57,7,19,32,0]` and `[57,7,19,16,0]` yield `48656` and `48655`. A read-only scan found zero duplicate keys among the 61,566 records in the loaded Paul M16 class table. The [Stage 10 notes](../../tools/revkit/work/stage10/README.md) give the precise fixture and trace commands.

## Requested peer check

Please identify any point where your independently recovered selection path disagrees with these **local, numeric** results: the changed tree input positions, context-1 population threshold, context-4 three-phone branch, or P/B versus F/V candidate `+0x10` behavior. A useful reply would name the input, function or port routine, expected numeric value, observed value, and the first divergent boundary. We will treat whole-WAV parity as lead 2 after this checkpoint.

The snapshot omits proprietary binaries and voice/model inputs. It contains the portable fixtures, scripts, traces, curated documents, and Git metadata needed to inspect our evidence. The DLL comparison and license-check differences in your prior reply remain attributed to your report because your binaries were not in the bundle; none of those differences is used to assign a linguistic meaning to our numeric fields.
