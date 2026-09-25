# Lead 3: TPP typed-code behavior (2026-09-24)

## Bounded result

The 2013 M16 Paul TPP corpus and runtime traces establish the tested
operational behavior of the `A`–`G` atoms and the exceptional `AX` spelling.
Place-name atoms `A`–`E` return a component count followed by one binary
value per hyphen-separated component. `F` and `G` return decimal suffixes;
in the tested runtime paths, each number is written to token-record byte
`+0x25`. For a matched `F` compound, the same value is written across the
component tokens. The `AX` payload is the literal `X`; two tested keys
returned `X` through the `A` selector.

These are mechanics and corpus associations. Later corpus analysis below
connects numeric IDs to recurring lexical forms and component bits to the
`A0`/`A1` distinction, but it does not recover authoritative vendor labels,
explain the purpose of `AX`, or establish a universal acoustic role.

## Corpus findings

[`inspect_tpp_dictionary.py`](../../tools/revkit/scripts/inspect_tpp_dictionary.py)
decoded and re-encoded all 31,550 keys against the local `vt_pau.dll` tables.
Each payload is one or two atoms. The observed two-atom sequences are `A G`
(66 records), `A A` (5), `B B` (1), and `C C` (7).

| Primary form | Records | Directly observed association |
| --- | ---: | --- |
| `A0`/`A1` | 7,899 | Standard single-component place-name atom; two binary values |
| `B` | 3,980 | Two hyphen-separated components; then two binary values |
| `C` | 303 | Three components; then three binary values |
| `D` | 19 | Four components; then four binary values |
| `E` | 2 | Five components; then five binary values |
| `F` | 829 | Hyphenated lexical keys with one to four hyphens and a decimal suffix |
| `G` | 18,493 | Non-hyphenated lexical keys and a decimal suffix |
| `AX` | 25 | Special payload; key contents do not identify a stable semantic class |

There are 7,924 records with an `A`-prefixed primary atom when the 25 `AX`
records are included; the standard place-name `A0`/`A1` count is 7,899.

`F` and `G` suffixes range from 1 through 124 in this corpus. The raw
decimal values are preserved in the decoded manifest; no numeric class name is
assigned. The component-bit correlation with standalone `A0`/`A1` entries is
analyzed below.

The exact serialization rule for `B`–`E` is recovered: the payload carries
the family/count, followed by one literal `0` or `1` for each hyphen-delimited
component in key order. The parser validates that the encoded count, bit
count, and number of hyphens agree; the runtime returns those bits. This
defines which bit belongs to which component. It does not define what the two
bit values mean or a rule for deriving them from spelling.

The complete numeric association inventory is in
[`numeric-class-corpus.tsv`](../../tools/revkit/work/stage12/numeric-class-corpus.tsv):
20 distinct `F` suffixes and 121 distinct `G` suffixes, covering 123 of the
124 values in the range; value `43` is absent. Each row records the full corpus
count, suffix-ending counts (`-ING`, `-LY`, `-ER`, `-ED`, final
`-S`, or other), and up to eight examples. This is a reproducible inventory of
spelling associations, not 124 complete linguistic definitions. Several
classes have strong spelling clusters (`G120`/`F120` with `-ING`, `G95`/`F95`
with `-LY`, and `G112` with `-ER`); many others are heterogeneous by these
coarse tests. No vendor class-name table or consumer source labels were found.

The summary can be regenerated with
`PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_tpp_numeric_classes.py`.

## Runtime contrasts

The isolated inputs and traces are under
[`tools/revkit/work/stage12/`](../../tools/revkit/work/stage12/). They use the
supplied `voicetext_paul.exe` and unchanged local `vt_pau.dll` under the
existing isolated Wine setup. At `0x1003a7b0`, GDB records selector, lookup
key, return value, and output string. At `0x1000e0c0`, it records the numeric
code and token byte `+0x25` after the call.

| Input/key | TPP result | Follow-on observation |
| --- | --- | --- |
| `Accord` | selector `A`: `A0` returns `0`; selector `G`: `G95` returns `95` | `FUN_1000e0c0` writes 95 to token byte `+0x25` |
| `Center` | selector `A`: `A1` returns `1`; selector `G`: `G95` returns `95` | Same suffix write; the place-name flag differs while the suffix is held constant |
| `Abbot Village` | selector `B`: `B201` returns `201` | Two components and two returned bits; neither bit is semantically named |
| `Agua Dulce` | selector `B`: `B200` returns `200` | Same component count with a different second bit |
| `Amelia Court House` | selector `C`: `C3010` returns `3010` | Three component positions are present in the runtime result |
| `Cape May Court House` | selector `D`: `D41010` returns `41010` | The separate `May` token also matches `G74`, which writes 74 to its byte `+0x25` |
| `Opelika`, `Talladega` | selector `A`: `AX` returns `X` for each key | Their `G` lookups miss; the probe finds no numeric suffix write for `AX` |
| `About shipping` | selector `F`: `F120` returns `120` for `About-shipping` | Value 120 is written to byte `+0x25` for both component tokens |
| `Abed`, `Abandoning` | selector `G`: suffixes `95`, `120` | Each value is written to the corresponding token byte `+0x25` |

The `E` family is represented by two corpus records, but the ordinary-input
probe for `Casa de Oro Mount Helix` did not issue an `E` selector call. The
follow-up phrase controls (`Casa-de-Oro-Mount-Helix` and
`Saint-Mary-of-the-Woods`) also reached the general `F` selector with no
match. The proper-name grouping traces show shorter spans around `de`, `of`,
and `the`.

A separate GDB probe then changed the selector argument at `FUN_1003a7b0`
from the naturally observed `F` (`0x46`) to `E` (`0x45`) only for those exact
23-byte keys. Both calls returned the lookup-hit status `eax=0x660001` and
these strings:

| Key | E result | Five component bits |
| --- | --- | --- |
| `Casa-de-Oro-Mount-Helix` | `500010` | `00010` |
| `Saint-Mary-of-the-Woods` | `500000` | `00000` |

The trace cleared the output byte and changed the result to a miss immediately
after capture, then let the process exit normally. This verifies that the
original DLL can retrieve and return both E records through the shared lookup
helper. It does not show a natural caller selecting E or establish how the
returned bits affect pronunciation.

The x86 instructions resolve that apparent pseudocode contradiction for this
caller. At `0x1003419e`, the function executes `cmp eax,0x4`, followed at
`0x100341a1` by `jg 0x1003424e`, which skips the complete selector-construction
block for counts of 5 or more. The E case stub at `0x10034210` therefore
cannot run through this path. This was checked with
`objdump -d -Mintel --start-address=0x10034180 --stop-address=0x10034255
binary/vt_pau.dll`. The recorded direct references to the lookup helper are
`FUN_10034180`,
`FUN_1000dfc0`, and `FUN_1000e160`; the latter two pass selectors G and F,
respectively ([xref report](../../tools/revkit/work/reports/stage6-tpp-xrefs.txt),
[caller pseudocode](../../tools/revkit/work/reports/stage6-tpp-callers.txt)).
Thus, in the discovered direct caller graph, neither natural phrase crafting
nor a five-part proper-name span can dispatch E. This is a call-graph-bound
finding; it does not rule out an unobserved indirect call elsewhere.

The initial `F`/`G` traces demonstrated a field write but did not isolate its
effect. The subsequent same-input overrides and acoustic comparisons are
reported below.

## Extended corpus associations and controlled effects

The numeric payloads have strong spelling associations in the full corpus.
Of 9,488 `G120` keys, 9,386 end in `-ING` (98.9%); of 7,098 `G95` keys,
6,299 end in `-LY` (88.7%). The corresponding `F` sets show the same
association: 357/362 `F120` keys reduce to an `-ING` ending after removing
hyphens, and 119/256 `F95` keys reduce to `-LY`. Other frequent examples also
cluster: 48/49 `G112` keys end in `-ER` and are comparative forms; the `G83`
sample is rich in pronouns, and `G20` in number words. This supports an
inference that at least some numeric IDs encode lexical or morphological
families. It does not give a complete definition: `G95` includes 799 keys not
ending in `-LY`, `G120` includes 102 keys not ending in `-ING`, and the class
IDs still lack authoritative linguistic names.

The place-name bits closely track the single-component `A0`/`A1` split. For
multi-component records, among components whose spelling also has a standalone
`A` record, all 2,114 components mapped to `A0` have bit `0`; 1,133/1,135
components mapped to `A1` have bit `1` (the exceptions are `MISSION` and
`PIKE`). The `A1` inventory is concentrated in geographic feature terms such
as `CENTER`, `VALLEY`, `PARK`, and `COVE`. This strongly supports the
operational reading that the component bit mirrors the standalone `A0`/`A1`
distinction when that component has an `A` record. The vendor's intended name
for the distinction, and the rule for components without standalone entries,
remain unknown.

The complete cross-reference is in
[`component-bit-corpus.tsv`](../../tools/revkit/work/stage12/component-bit-corpus.tsv).
It separates component occurrences from unique spellings: the 2,114 `A0`
occurrences all carry bit `0`; 1,133 of 1,135 `A1` occurrences carry bit `1`.
Both exception spellings also appear with bit `1` elsewhere, so even those
exceptions do not define a stable spelling rule. Of 5,706 compound-component
occurrences without a standalone `A0`/`A1` key, 3,387 carry `0` and 2,319
carry `1`. Seven such spellings occur with both values (`EAST`, `SAINT`,
`TOWNSHIP`, `FERRY`, `TOWN`, `GATE`, `SAINTE`). Thus the bits are stored on
the compound entry and cannot be reconstructed as a single-valued mapping
from component spelling. The corpus establishes their values and correlations,
but does not reveal the vendor's exact linguistic rule or names for the two
states. Regenerate the cross-reference with
`PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/summarize_tpp_component_bits.py`.

Static caller analysis does expose one exact operational effect for a bit.
`FUN_10034180` rejects a multi-component candidate when a returned bit is
literal `0` and that component's internal marker at `+0x1a` is `d` or `A`.
For a one-component lookup, the payload bit is not read in that gate; the
function checks only the marker. A controlled runtime probe verified the
multi-component branch against `Abbot Village` (`B201`, bits `01`): the
ordinary `F` markers return low-word success; forcing component 0's marker to
`d` returns `eax=0`; forcing component 1's marker to `d` preserves
`eax=0xffff0001`; and changing the first bit to `1` (`B211`) while that
marker is `d` restores `eax=0xffff0001`. This confirms the bit-to-component
order and the zero-bit rejection rule. It does not identify what the vendor
intended the `d`/`A` markers or bit values to mean, or how unmatched component
bits were assigned.

The producer narrows the `d` marker's immediate source. In pseudocode,
`FUN_10043810` assigns `d` after a successful lookup through the table at
`0x10099078`, then copies a 0x30-byte record from `0x10099088 + index * 0x30`
into its working record. The read-only DLL table contains short-form keys and
replacement text: `N` → `North`, `APT` → `Apartment`, `acctg` → `accounting`,
`aka` → `also known as`, and `dept` → `department`. A runtime trace on
`42 N St.` confirms lookup index 1, marker byte `0x64` (`d`), and replacement
text `North`. In matched `42 Main St.`, lookup attempts for `Main` and `St`
miss; `42 North St.` also misses because `North` is the replacement value,
not this entry's key. Thus a natural replacement-table hit does produce the
same `d` marker used by the controlled bit-gate test. This ties `d` to the
replacement path, but does not prove the vendor intended the marker's name to
mean “abbreviation” or characterize every downstream use. The separate
Apartment/Basement branch writes different marker bytes: in `Unit 4 and Apt
4.`, the `Unit` path writes uppercase `A` with working text `Unit`, while the
`Apt` table hit writes lowercase `a` with working text `Apartment`. `The apt
is ready.` also writes lowercase `a`, although the later text path can retain
the surface `apt`. Since `FUN_10034180` checks uppercase `A` and lowercase `d`,
the distinction is operationally relevant; these traces do not assign vendor
names to the markers. The `42 N St.` run did not call `FUN_10034180`, so the
natural producer trace does not establish that this `d` value reaches the TPP
bit gate on that input. In `Unit 4 and Apt 4.`, the tokenizer emits `A` for
`Unit` and `a` for `Apt`, but the later TPP component trace still sees marker
`F` for the `Unit` and `Apt` candidates. Thus these tokenizer-local markers
are not shown to propagate into the classifier's `+0x1a` field in the tested
address path.

Same-input overrides then isolated downstream behavior. The probes are
`trace-numeric-zero.gdb`, `trace-numeric-swap.gdb`,
`trace-category-override.gdb`, and `trace-effects-control.gdb`; the runner is
`run-effects.sh`. For the numeric tests, the return-site trace confirms the
overwritten byte `+0x25` is actually zero or the swapped value after
`FUN_1000e0c0` returns. For the place-name tests, the lookup return buffer is
changed in memory before processing continues. Controls and interventions
used the same text inputs and runtime setup.

| Intervention | Inputs | Measured result versus matched control |
| --- | --- | --- |
| Numeric code set to zero | `Abed`, `Abandoning`, `About shipping` | Phone-row bytes, ordered selected-unit IDs, and whole WAV were identical |
| Numeric code `95` swapped with `120` (both directions) | `Abed`, `Abandoning`; `About shipping` (`120` to `95`) | Phone-row bytes, ordered selected-unit IDs, and whole WAV were identical |
| `A0`/`A1` bit flipped | `Accord`, `Center` | Phone-row bytes, ordered selected-unit IDs, and whole WAV were identical |
| One component bit cleared in B/C/D payloads | `Abbot Village`, `Amelia Court House`, `Cape May Court House` | Phone-row bytes, ordered selected-unit IDs, and whole WAV were identical |
| `AX` changed to `0`, `1`, or a lookup miss | All 25 `AX` keys in one phrase batch | Phone-row bytes, 376 ordered selected-unit records, and whole WAV were identical |
| Numeric code zeroed in additional phrase contexts | Seven additional `G` codes, `F95`, and eight initial sentence contexts | Most were identical; focused `G83` probes found effects in selected question contexts below |
| `G83` changed to `G95` | Four matched `anybody` contexts | Two question cases changed units and WAV; two other contexts were byte-identical |

The selected-unit sequences were equal in all five A–D/AX or numeric-ablation
fixtures (11, 149, 32, 23, and 17 selected-unit records respectively), and in
both numeric-swap fixtures and the AX-to-1 fixture. WAVs were byte-identical
in each comparison. This
shows that the tested values and bits do not affect the measured pronunciation
path for these specific words and contexts. It does not prove that they are
unused in all text contexts, engine builds, or processing paths.

### Numeric-code effects in additional contexts (follow-up through 2026-09-25)

The broader context sweep found a real, context-sensitive `G83` effect. In
these controls, the GDB intervention changes the numeric argument passed to
`FUN_1000e0c0`; the matched control keeps the original code. The phone-row
bytes remain identical in every comparison, so the measured change occurs
after those rows are formed.

| Fixed phrase | Intervention | Selected units | WAV frames | Result |
| --- | --- | ---: | ---: | --- |
| `Does anybody know?` | `G83` to zero | 22 → 26; 14 same-position IDs differ | 15,060 → 15,100 | Changed |
| `Does everybody know?` | `G83` to zero | 27 → 27; 6 IDs differ | 17,200 → 16,932 | Changed |
| `Does anybody see?` | `G83` to `G95` | 24 → 24; 8 IDs differ | 16,725 → 16,022 | Changed |
| `Does anyone know?` | `G83` to zero | Same sequence | 13,628 → 13,628 | Identical |
| `Does somebody know?` | `G83` to zero | Same sequence | 14,357 → 14,357 | Identical |
| `Does she know?` | `G83` to zero | Same sequence | 10,581 → 10,581 | Identical |
| `Can anybody help?` | `G83` to zero | 32 → 32; 12 IDs differ | 16,245 → 16,876 | Changed |
| `Anybody knows.` | `G83` to zero or `G95` | Same sequence | 16,087 → 16,087 | Identical |
| `Nobody answered.` / `Somebody came.` | `G83` to zero | Same sequence in each | Per-file WAVs identical | Identical |
| `Do they know?` | `G83` alone to zero | Same sequence | 10,428 → 10,428 | Identical |

The code does not act as a uniform global prosody switch: the effect depends
on both the dictionary key and its phrase context. The positive interventions
establish that a numeric code can influence unit selection and synthesized
audio after the phone rows are fixed. They do not identify the affected
acoustic feature (for example, stress, duration, pitch, or coarticulation), nor
do they provide a general rule for predicting which contexts use it. The raw
controls, interventions, selected-unit traces, and WAVs are preserved under
`tools/revkit/work/stage12/`.

An additional `Can anybody help?` contrast exposed an interaction: zeroing
either `G83` alone or `G74` alone produces the same 32-unit sequence and the
same 16,876-frame WAV (12 selected IDs differ from control). Zeroing both
codes together restores the exact control sequence and 16,245-frame WAV;
replacing `G83` with `G95` also matches control. This is a non-additive
interaction in this phrase, not evidence for a general role for either code.

### `AX` caller behavior and purpose boundary

The complete corpus contains 25 `AX` keys: `OPELIKA`, `OPPORTUNITY`,
`ORCHARDS`, `PINEWOOD`, `PONCE`, `PRATTVILLE`, `PRICHARD`, `REDLAND`,
`ROSSMOOR`, `SHERRELWOOD`, `SOUTHGLENN`, `TALLADEGA`, `TAMIAMI`, `TROTWOOD`,
`WAIMALU`, `FAJARDO`, `OATFIELD`, `OILDALE`, `OJUS`, `SMITHS`, `TUSCALOOSA`,
`VALINDA`, `VINCENT`, `WOODLAWN`, and `YAUCO`. Each is returned as literal
`X` by the `A`-selector lookup. The set has no separate `G` record for the same
key. In one phrase batch containing every key, changing each returned `X` to
`0`, to `1`, or to a lookup miss left phone rows, all 376 selected-unit
records, and the WAV unchanged. This bounds the tested A-selector path, but
the caller's purpose for preserving these dictionary records is still
unknown; neither the key set nor the intervention identifies it.

The static caller narrows the role. `FUN_10034180` checks the lookup hit, then
for a one-component key checks the token marker at component-record `+0x1a`
against literal `d` and `A`; it does not inspect the returned `A0`, `A1`, or
`AX` payload in this branch. For a multi-component key, it rejects the
candidate only when a returned component byte is literal `0` and that
component marker is `d` or `A`. The enclosing caller tests the low 16-bit
return from `FUN_10034180` for success. This is a directly observed branch
condition, not a semantic name for `d` or `A`.

A focused runtime trace on `Opelika` and `Talladega` supports that path:
`AX` returns caller status `eax=1`; replacing only the payload byte `X` with
`0` still returns `eax=1`; forcing the lookup to miss returns `eax=0`. The
`X`→`0` and miss runs produced identical WAVs. Their token marker was `F`, so
neither literal `d` nor `A` branch was active. Together, these findings
support the limited inference that an `AX` entry acts as a recognized-key
marker in this one-component classifier path, while its literal `X` is not
used there. In the multi-component guard, `X` would also differ from literal
`0` and therefore avoid that specific rejection; however, the corpus has no
multi-component `AX` record to test this case. Why these 25 entries use `AX`,
and whether another consumer interprets `X`, remain unresolved.

## E selector caller audit

The local Ghidra report lists three direct call references to
`FUN_1003a7b0`; each supplies G, F, or the A-through-D proper-name selector
path. The five-component branch in `FUN_10034180` is skipped by the observed
`count > 4` guard. Wag's independent review in `vtpaul_codex_status_2026-09-24-new-3`
reports the same three references and no data references to the helper, so
the reviewed image contains no discovered function-pointer-table route to E.
This strengthens the negative call-graph result, while remaining a static
analysis bound: it cannot rule out a runtime-generated target or code outside
the analyzed DLL. The forced E lookups still establish that both E records are
stored and retrievable through the helper.

## Adjacent record paths and offset separation (2026-09-25)

The DLL has another path that uses byte offset `+0x25`, but in a different
record layout. The TPP tokenizer writes the Lead 3 code through
`FUN_1000e0c0` into a `0x94`-byte token record. Separately,
`FUN_10008dc0` calls `FUN_1000c9c0` to normalize exception-dictionary phone
text and, on a match, passes its encoded result to `FUN_1000ca50`. That helper
writes phone-symbol bytes into records spaced `0x70` bytes apart; its `+0x25`
is therefore not the TPP token's `+0x25`. In this path, `d` is a delimiter in
the encoded phone string: `FUN_1000ca50` either retains it as a byte or ends
one `0x70` record and advances to the next, depending on the count returned by
`FUN_1000ca30`. The pseudocode for these functions is preserved in
`stage6-exception-dictionary-expanded.txt` and
`stage6-dictionary-transform-expanded.txt`.

The wider phone-row path calls `FUN_1000d190` from `FUN_1000e2f0`, builds
phone rows with `FUN_1000d450`, and later passes the rows through
`FUN_10007520`. The runtime capture in `stage12/trace.gdb` prints both the
upstream `0x94` token fields and the resulting `0x70` context rows. This gives
us a useful tracing boundary for future work, but it does not establish that
the Lead 3 numeric byte itself is consumed by this phone-row path: the tested
`G83` interventions leave the phone-row bytes unchanged while later selected
unit IDs and WAVs can change.

This audit rules out one tempting offset-based conflation and sharpens the
numeric-code trace. The static path in `stage7-backtrack.txt` shows
`FUN_10024980` calling `FUN_10024680` to assemble candidate units. For each
candidate, `FUN_10024680` tries `FUN_10024060` and falls back to
`FUN_100242a0`; both call `FUN_10018770` to construct candidate context, then
use different downstream filters (`FUN_10023060` and `FUN_10023c70`). After
candidate assembly, `FUN_10024980` scores candidate rows and calls
`FUN_1001b200`, which extracts unit-feature records through `FUN_1001af50` or
`FUN_1001b0d0`.

The new runtime pair uses `numeric-context-g83-can-anybody-help` with
`candidate-path-control` and `candidate-path-83-zero`. Zeroing `G83` leaves
the phone rows unchanged, produces the same 18 six-byte candidate descriptors,
and gives the same 18 filter outcomes across `FUN_10024060` and
`FUN_100242a0`. The selected-unit stream is still 32 entries, but 12 IDs
differ. The first observed score-vector difference is at the return from
`FUN_10018c80` for candidate index 1: slot 28 is `0.0176470596` in control
and `0.0157894734` with `G83` zeroed; the other 29 values match. After
`FUN_10024510`, the corresponding 30-value candidate vector differs more
broadly. At the final `FUN_10024900` ranking vector captured for this phrase,
all 30 candidate costs are higher by approximately `5.013245` when `G83` is
zeroed; the selected index remains 0. This traces a path from the changed code
through score state after candidate generation. It does not explain why
zeroing this code changes the `FUN_10018c80` vector, which other candidate
groups it affects, or which acoustic dimension changes. The next step is
tracing score-vector writes and backtracking inputs in `FUN_10018c80`,
`FUN_10024510`, `FUN_10023350`, `FUN_10024900`, and their callers; this still
does not identify a linguistic name for `G83` or the other suffix classes.

The reproducing GDB scripts and full traces are
`stage12/trace-candidate-path-control.gdb`,
`stage12/trace-candidate-path-83-zero.gdb`, and the corresponding
`numeric-context-g83-can-anybody-help-effect-candidate-path-*.log` files.
Stage 5 input/output hashes were unchanged after both runs.

## Suffix-class and acoustic follow-up (2026-09-25)

The full-corpus associations continue to support lexical or morphological
families, without supplying vendor terminology. The corpus has 17/18 `G96`
keys ending in `-ER`, 48/49 `G112` keys ending in `-ER`, and 7/7 `G113` keys
ending in `-EST`; the `G96` and `G112` examples are comparative forms and the
`G113` examples are superlatives. `G83` has 57 keys and is rich in pronouns,
but also contains other forms. These are spelling and lexical correlations,
not a proof that the numbers name parts of speech, inflection rules, or sound
features. The corpus summary is `stage12/numeric-class-corpus.tsv`.

A second paired intervention traces selected IDs into the unit descriptor
constructed by `FUN_1002d1e0`/`FUN_1002c120`. For the first six of 14 mapped
unit calls, zeroing `G83` changes the selected ID, DAT sample count, UPM
period count, or UPM edge values. The later eight mapped descriptors are
identical between runs. Both runs use the same input and otherwise follow the
same process path:

| Call | Control ID; samples, periods, edges | `G83=0` ID; samples, periods, edges |
| ---: | --- | --- |
| 1 | `121166`; 840, 6, 130/152 | `94212`; 798, 5, 152/164 |
| 2 | `130020`; 1,402, 10, 146/142 | `255932`; 958, 6, 172/156 |
| 3 | `130760`; 1,800, 10, 188/170 | `42414`; 1,690, 14, 120/116 |
| 4 | `121168`; 2,428, 19, 138/130 | `94215`; 758, 7, 114/98 |
| 5 | `393743`; 1,566, 13, 158/126 | `393745`; 1,434, 13, 102/122 |
| 6 | `130762`; 602, 5, 124/116 | `42417`; 1,224, 10, 120/118 |

At the complete-WAV level, the control is 1.0153 s and the `G83=0` output is
1.0548 s. A simple 16 kHz PCM summary measured RMS 5755 versus 6017 and peak
29,696 versus 30,665. An exploratory 25 ms autocorrelation pitch estimate
gave medians of 134.5 and 135.6 Hz respectively; this estimator and one phrase
are not sufficient to identify a pitch or other acoustic feature as the
numeric code's intended control. These output differences could arise from
the units and timing selected downstream, and do not isolate a particular
source feature.

This closes a useful trace gap: the effect is no longer only a score-vector
change or a WAV hash difference; runtime evidence connects the intervention to
different selected unit descriptors and their timing metadata. The suffixes'
linguistic labels remain unproven. A later candidate-list trace places an
earlier divergence before `FUN_10018c80`, while the exact context input that
changes those class candidates and the intended linguistic feature remain
unresolved. The paired traces are
`stage12/numeric-context-g83-can-anybody-help-effect-acoustic-control.log`
and `stage12/numeric-context-g83-can-anybody-help-effect-acoustic-83-zero.log`;
their WAVs share those basenames. Reproduce them with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-effects.sh acoustic-control \
  numeric-context-g83-can-anybody-help
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-effects.sh acoustic-83-zero \
  numeric-context-g83-can-anybody-help
```

Reproduce the pair with the existing isolated runtime container:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-effects.sh candidate-path-control \
  numeric-context-g83-can-anybody-help
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-effects.sh candidate-path-83-zero \
  numeric-context-g83-can-anybody-help
```

### Candidate-pool follow-up (2026-09-25)

A focused paired trace now observes the class lists after
`FUN_10024680` returns and the unit-index lists after `FUN_10023350` returns.
The runtime override is applied at entry to `FUN_1000e0c0`; the trace prints
`TPP_NUMERIC_ABLATE83` for the changed call. Both runs complete normally.

`FUN_10024680` produces 18 contexts. The first class-list difference is at
context 1: control returns six IDs (`30504, 30497, 30498, 30500, 30501,
30503`), while `G83=0` returns eight (`30503, 30497, 30498, 30499, 30500,
30501, 30502, 30504`). Contexts 2 and 3 also differ (2 versus 3 class IDs
each); contexts 0 and 4–17 have identical ordered class lists.

After `FUN_10023350`, unit-index candidate lists differ at contexts 0–3.
Context 3 has 25 candidates in control and 30 with `G83=0`; contexts 0–2
have 30 in both runs but different candidate multisets. Contexts 4–17 have
identical ordered lists. This places a measured code effect in candidate
generation/local ranking before the later adjacent-context transition score.
The previously reported `FUN_10018c80` score-vector difference is therefore
not the earliest observed divergence. The control and ablated output WAVs
have SHA-256 values `29052da3a11213f4616c30b2da31ccd9c4400d657afa1c163a51135653ddad09`
and `f8aa3133ef5701e78f46db0bd55cdcf3674430c988d7b6f5834bd804a0683305`,
respectively.

These captures establish a local, context-dependent change in the class and
unit candidate sets, not a vendor label for code 83 or for the numeric class
IDs. They also do not show a universal rule for which contexts are affected.
The logs preserve all 18 class lists and all post-`FUN_10023350` candidate
lists in `stage12/numeric-context-g83-can-anybody-help-candidate-lists-{control,g83-zero}.log`.
Reproduce them with the script documented in
[`tools/revkit/work/README.md`](../../tools/revkit/work/README.md#lead-3-candidate-list-trace).

## Reproduction

With the local vendor inputs and existing runtime image:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run.sh \
  place-flags ax-special g-suffix f-suffix place-components
```

The two E records can be retrieved through the shared lookup helper with a
scoped in-memory selector override. The trace discards the returned value
before normal processing resumes:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-e-direct.sh \
  e-casa-hyphens e-saint-hyphens
```

See [`trace-e-direct.gdb`](../../tools/revkit/work/stage12/trace-e-direct.gdb)
and the two `*-e-direct.log` captures for the exact key filters, selector
mutation, returned bytes, and successful process exits.

The isolated `G83` and `G74` context probes use the same runner with
`numeric-83-zero`, `numeric-74-zero`, `numeric-83-to95`, `numeric-zero`, or
`effects-control` modes. For example:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage12/run-effects.sh numeric-83-zero \
  numeric-context-g83-can-anybody-help numeric-context-g83-do-they
```

The paired control and intervention `.log`/`.wav` files use the same fixture
name under `tools/revkit/work/stage12/`. The focused runs on 2026-09-25 also
restored Stage 5's protected input and output hashes shown below.

The caller status for `AX` can be traced with `place-calls`,
`place-calls-ax-zero`, and `place-calls-ax-miss` modes using the `ax-special`
fixture. `trace-place-call.gdb` captures the component marker and return from
`FUN_10034180`; the two variants alter only the matching A-selector result.

The component-bit gate is reproducible with `place-gate-first-d`,
`place-gate-second-d`, and `place-gate-first-d-bit1` using
`place-bit-gate-ab-village`. These alter the parsed component marker and, in
the final mode, only the first returned B bit; they do not patch the DLL or
vendor dictionaries.

The natural `d` producer trace uses `marker-producer` with
`marker-direction-address` and `marker-street-address-control`. Its GDB script
records the lookup result at `0x10043a8f`, the marker write at `0x10043b10`,
and the copied replacement value. The runtime also completed normally for
these controls, and the runner restored Stage 5's protected files. The
`marker-apartment-address-span` and `marker-apartment-neutral` fixtures capture
the distinct uppercase `A` and lowercase `a` writes described above.

The runner restores Stage 5's input and output files on exit. The final run's
restoration was checked byte-for-byte against its pre-run copies. The
corpus-wide result is reproducible with:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/inspect_tpp_dictionary.py
```

Lead 3 is not fully closed. The empirical inventory is now exhaustive for the
local 31,550-key corpus, the exact B–E serialization and zero-bit gate are
established, and context-controlled probes establish that `G83`
can affect unit selection and output audio in some phrases. However, the
corpus and tested binary do not expose authoritative meanings for all numeric
classes, a universal linguistic assignment rule for compound-component bits,
or the reason the 25 `AX` records exist. The `d` marker's producer is now tied
to a replacement table with abbreviation/expansion entries, but that does not
establish the vendor's intended marker name or fully characterize its runtime
role. The acoustic effect is measured but its specific feature and full
context condition remain unresolved. Independent static review found
no indirect E reference in the analyzed image, while a runtime-generated or
out-of-image caller remains outside the evidence bound. These are evidence
limits, rather than gaps that can be closed by assigning labels from spelling
alone.
