# Follow-up investigation leads (2026-09-24)

This note records research directions and dated progress after the feasibility
pass for the local 2013 M16 Paul package. These leads do not change the status
of completed roadmap stages. The priority order favors questions that can
extend existing runtime evidence. Keep all VoiceText binaries and model data
read-only, and distinguish observed behavior from semantic interpretation.

## 1. Label more decision-tree inputs with controlled contrasts

**Question.** Which tree input fields track phone identity, stress, neighboring
phones, part of speech, and phrase or word boundaries? What do those fields
change in the selected tree outputs and final synthesis?

**Starting evidence.** Stage 10 already holds the visible token constant while
varying forced CMU phones and vowel stress. It identifies the one-based phone
identity values in some inputs, several consonant groups that select different
onset trees, and the stress value in selected tree positions. The remaining
tree selectors and the wider prosody consequences are still numeric. See the
[phone codebook](phone-symbol-codebook.md#controlled-tree-input-labels) and
the Stage 10 [tree tracing notes](../../tools/revkit/work/stage10/README.md).

**Bounded probe.** Extend the controlled fixtures one factor at a time: vary
one left or right phone, then a boundary, part-of-speech tag, or pronunciation
alternative while holding the rest fixed. Capture the tree identity, full
input vector, result, and downstream timeline rows. Repeat any apparent rule
with a second phone or word before labeling it.

**Useful result.** A verified mapping from one input field to a controlled
property and a measured downstream effect, with counterexamples or limits
recorded. The binary may expose operational roles without revealing the
vendor's original field names.

**Status: complete for the bounded local Paul result.** The
[lead 1 closure](lead1-tree-input-closure-2026-09-24.md) gathers the verified
input mappings, downstream selection path, counterexamples, and remaining
generalization limits. The detailed notes below preserve the investigation
sequence and raw-value references.

**Progress (2026-09-24).** The existing Stage 10 two-phone captures hold the
`probe` phone at `P AH0` and vary a forced adjacent phone. On the left-side
controls, the first two tree calls change and the next two remain identical;
on the right-side controls, the first two calls remain identical, the third
changes its input and return, and the fourth changes its input but keeps the
same return. The recorded values and limits are
in the [neighbor-context report](tree-context-neighbor-contrasts-2026-09-24.md).
Follow-up traces reach selection and timeline construction. With a fixed
`P AH0` span followed by P versus B, the two tree calls for the fixed span
retain identical vectors and results, but its first three primary selected
IDs and timeline sample counts change. With P versus B before the fixed span,
the changed first three primary rows track the changed left phone, while the
last three remain the same. The downstream matrix and raw captures are in the
neighbor-context report and `tools/revkit/work/stage10/`. This establishes a
selection effect beyond the unchanged fixed-span tree results for the tested
right-side contrast. Repeating the right-side P/B contrast with fixed `T AH0` produced
the same pattern: the fixed span's first two tree results remain equal, while
the first timeline call's first three selected IDs differ. For both T
boundary controls, a period leaves the preceding span's first timeline call
identical across the right-side P/B change; divergence appears in the
following timeline call. These results bound the tested behavior but do not
identify the underlying input field or establish a general linguistic rule.
Candidate traces now show that the right P/B pair changes fixed-span local
candidate costs and that the following context's transition chooses a
different predecessor. For the common candidate `3837`, its traced `+0x10`
field changes from 0 to 1 while the raw scorer target and feature view stay
identical; the scorer scale changes from 1 to 0.5 and the local cost follows
from `0.833333` to `0.416667`. Stage 7 pseudocode connects that scale to the
candidate's `+0x10` field and a normalized duration term, but the upstream
source of the field remains unknown. A trace at `FUN_100230a0` shows the
field changing during candidate metadata construction: the P/B runs enter
with candidate `3837`'s `+0x10` at 0, then leave it at 0 for P and 1 for B.
Its disassembly walks earlier and later context records. A focused trace
shows that for context-0 candidate `3837`, one lookup reaches following
context 1 in all four right-phone runs. Context-row byte `+0x04` at
`state + 0xec628 + context_index * 6` indexes a local 16-bit table: P/F/V
read byte 1 and contribute 1; B reads byte 0 and contributes 2. The forward
accumulator is halved before it is stored at `+0x10`, producing 0 versus 1.
The row byte is also used by the Stage 7 transition routine, but its semantic
meaning and the table's purpose remain unnamed. Its writer is now located
in the `FUN_10024680` fallback chain: `FUN_10024060` clears byte `+0x04` and
returns a status; when it returns false, `FUN_100242a0` builds the fallback
row and writes `+0x04 = 1` (or `2` on its copied second row). A hardware
watchpoint captured context 1's `0 → 1` write in right P at `0x10024305`;
right B made no context-1 change. The fallback-entry trace records indices
1 and 4 for right P/F/V, and index 3 for right B. The context-1 byte values
align with whether that context uses the fallback. A return-site trace at
`0x10024808` confirms context 1 returns `0` from `FUN_10024060` for P/F/V
and `1` for B; context 3 returns `1` for P/F/V and `0` for B. The caller
tests the result at `0x1002480b`, taking the fallback when it is zero. This
locates the immediate branch responsible for the row difference. The
candidate-count trace shows that all four fixtures use model flag `0`, so
the first pass accepts when the `FUN_10023060` sum is at least 10. At context
1, P has one candidate (`48721`, table value 1), B has two (`48655`/10 and
`48656`/4), and F/V generate no candidates. Their sums are respectively 1,
14, and none, matching false/true/false returns. At context 3, the sums are
113 for P, 3 for B, 141 for F, and 188 for V; only B takes the fallback.
`FUN_10023060` sums 16-bit table values addressed through the model pointer
at `+0x8c`. This explains the tested branch results, but not why candidate
generation yields these IDs and weights. A `FUN_10018770` entry trace shows
the context-1 input as `[90,57,7,X,7,97,0]`, with only `X` changing. Its
translation checks map `+2 = 7` to `0` and `+1 = 57` to `77` for P/B/F/V;
the changed `+3` maps P's `53` to `71`, while B/F/V values `19/32/65` map to
`0`. The P nonzero branch's three additional lookups return zero, as does
the B/F/V zero-route lookup using mapped `77`. Successful combined-signature
queries separately emit P ID `48721` and B IDs `48656`, `48655`; F/V emit no
candidate. Static disassembly places the zero-row-byte query through
`FUN_10023f90` → `FUN_10023dc0` → `FUN_10016ea0` and `FUN_10019450`, with a
table rooted at model pointer `+0x88`. `FUN_10016ea0` makes a five-byte key:
three lookup-table outputs from signature bytes `+1..+3`, a direct copy of
byte `+5`, and byte `+6` masked with `0x20`. `FUN_10019450` binary-searches
sorted five-byte records with a bytewise comparator and expands an equal key
to its contiguous record range. The class table has 61,566 records. Runtime
traces show P key `[57,7,53,16,0]` at class index `48721`, and B keys
`[57,7,19,32,0]` and `[57,7,19,16,0]` at `48656` and `48655`. Each
successful query returns one class, whose index matches the first-pass list
entry; F/V produce no exact class matches for either tested suffix variant.
The [phone-symbol codebook](phone-symbol-codebook.md) identifies the varied
byte values as P/B/F/V. Wag's separate `vt_eng.dll` 3.11.7.1 notes call this
signature position `R1`; that field name is a cross-build corroboration,
since our captures use `vt_pau.dll`. The suffix bytes `16` and `32` and
their full field semantics remain unresolved. A follow-up trace through
`FUN_10023350` reads class populations at model offset `+0x8c` and unit-list
pointers at `+0x94`: B classes `48655` and `48656` expand to 10 and 4 concrete
unit IDs respectively, matching the values previously summed by
`FUN_10023060`. In the broader P context-1 pool, class `48721` expands to its
single member, unit `246151`. This verifies class-to-unit expansion for the
captured classes; it does not label their acoustic or phonetic content. The
trace script and complete B expansion are in `trace-context-unit-expansion.gdb`
and `tree-context-t-right-b-unit-expansion.log` under
`tools/revkit/work/stage10/`. In the second right-side F/V
pair, context 0's 30 local candidate costs and `3837`'s
`+0x10` stay identical, while backtracking returns a different first two
units. Across the period, the preceding `0..2` context span has identical
candidate arrays, transition results, and backtracked units; the later span
diverges. The direct half-key lookup trace now captures numeric 10-byte keys,
query range lengths, and class populations for the split P/F/V cases; B remains
whole because it passes the whole-phone population threshold. This confirms
that the left-side first-pass pool is shared for P/F/V in this fixture while
right-side pools differ, but does not assign semantic names to the key fields.
The K/G pair adds a second threshold contrast: K falls back at context 1 and
uses half keys, while G passes the whole-phone first pass. A full scan of the
Paul M16 five-byte class table found no repeated keys among 61,566 records.
The left-neighbor P/B control reverses the context-1 whole-phone return seen
in the right-neighbor control; left B's half keys and candidate pools are now
captured. Three-phone B-left/T-probe controls with right P/B/F/V add a
context-4 split with identical six-byte rows: P/F/V return zero and use half
keys, while B returns one on the whole-phone path. Additional phone pairs
and context positions remain open, while class-to-unit mapping stays numeric
until the unit-record fields are independently identified.

## 2. Broaden full-synthesis parity coverage

**Question.** How far do the recovered text, selection, timing, and audio
paths reproduce original WAV output across different text classes and API
controls?

**Starting evidence.** DAT decode output matches the DLL for all 580,474 local
units by byte count and SHA-256. A smaller set of controlled synthesis runs
matches captured sample buffers to WAV data, including default and pitch
controls plus number and abbreviation fixtures. This is not whole-synthesis
corpus parity. Wag's latest status reports a separate 2,026-text byte-exact
engine-to-Rust WAV corpus, plus pitch/speed/volume waveform controls, and
independently describes the 580,474-unit decoder comparison as engine-decoder
output versus the C export. The new bundle adds a SHA-256/size manifest for
2,017 corpus texts and makes the test capture root configurable with
`VTPAUL_HARNESS`; the nine additional speak captures are older synthetic lines
that are not in that manifest. These remain peer results rather than local
reproductions: the bundle omits the DLLs, voice data, engine WAVs, and roughly
900 MB of ling/speak/stage captures. The manifest makes the expected whole-WAV
hashes reviewable but does not supply the missing inputs and captures needed
to rerun the complete parity suite here. The user also relays Wag's statement
that their comparison DLL has its license-check path patched to report a
licensed state, which avoids the engine's injected demo/license text. This is
peer-reported and cannot be independently checked from the bundle. It means
those WAV hashes do not establish original-license or demo-watermark behavior;
our local runtime evidence still uses the supplied `vt_pau.dll`. See the
[corpus parity and Stage 9 report](dat-corpus-parity-and-stage9-2026-09-24.md)
for our independent decoder-boundary result.

**Bounded probe.** Define a small, versioned golden input set covering ordinary
prose, names, numbers, dates/times, addresses, abbreviations, and short edge
cases. Capture deterministic full WAVs at default settings first, then vary
pitch, speed, and volume separately. Compare whole files and report the first
divergent boundary when bytes differ: token/context rows, selected unit IDs,
timeline, PCM blocks, or WAVE accounting.

**Useful result.** A reproducible pass/fail matrix with a clear coverage claim
and divergence localization. Do not extrapolate a small golden set to every
utterance or every engine/package version.

**Bounded result (2026-09-24).** The eight selected local Paul-package cases
all produced valid WAVs whose sample-buffer blocks exactly equal the local
WAV data chunks. Each local WAV has the peer manifest's byte length, but all
eight whole-file SHA-256 values differ. A prose case repeated identically
across two local runs. Stage 8's one-at-a-time pitch, speed, and volume runs
show the local frame-count and amplitude effects. The peer WAVs and runtime
inputs are absent, so the mismatch cannot be localized across engines; the
peer DLL's reported license patch also remains unverified. The bounded probe
is closed with negative whole-WAV parity; full-corpus and peer-boundary
parity remain open. See the [lead 2 WAV parity matrix](lead2-wav-parity-2026-09-24.md)
and its portable traces/manifests under `tools/revkit/work/stage11/`.

## 3. Investigate the unknown TPP code meanings

**Question.** What observable token or pronunciation behavior is associated
with TPP code families `A`–`G`, `AX`, and the numeric suffixes?

**Starting evidence.** All 31,550 TPP keys decode and re-encode; the typed
payload grammar, caller paths, and `F`/`G` suffix parsing are known. Static
analysis now confirms the B–E zero-bit rejection gate when the tokenizer's
component marker is `d` or `A`; a controlled marker/bit intervention confirms
the rule. A natural `N` → `North` replacement-table hit writes marker `d` in
the tokenizer; a `Main` control misses. The `A` producer appears in separate
address-form branches. Numeric suffix and component-bit linguistic meanings, `AX`'s
purpose, and the full context/acoustic roles remain unresolved. See the TPP discussion in the
[engine and model findings](voice-engine-and-model-formats.md#stage-6-text-and-pronunciation-resources)
and the [open-topics review](open-topics-review-2026-09-24.md).

**Bounded probe.** Select keys whose payloads differ in only one code or
suffix, where possible. Compare upstream token records, the returned
phone/context rows, selected pronunciation, and final surface/audio for
controlled neighboring text. Keep raw codes as the authoritative labels;
only assign a linguistic name when multiple independent observations support
it.

**Useful result.** A tested mapping from a code value to a particular
transformation or branch. If the available corpus has no discriminating
examples, record that limitation instead of guessing from key spelling.

**Bounded result (2026-09-24).** Runtime traces confirmed the `A`–`D`
component-count/bit returns, the `AX` return string, and the `F`/`G` numeric
suffix writes to token byte `+0x25`. An `F120` compound applies 120 to both
component tokens; `G95` and `G120` apply to individual tokens. A scoped
in-memory selector override retrieved both corpus `E` records and captured
their exact payloads (`Casa-de-Oro-Mount-Helix` → `500010`;
`Saint-Mary-of-the-Woods` → `500000`). Ordinary inputs still did not reach
selector E. Disassembly shows that `FUN_10034180` skips its selector switch
for proper-name counts above four, so its E case cannot execute; the other
two direct lookup callers use F/G. Natural E dispatch therefore appears dead
within the discovered direct caller graph. The corpus and caller-path
operations are now exhaustively inventoried for this local package, but
semantic closure is not supported. The full numeric suffix table leaves
linguistic meanings unnamed. Component-bit serialization and the runtime
zero-bit rejection gate are now verified, but the corpus shows that bit
assignment is not a function of component spelling alone. The reason for
`AX` remains unknown after all 25 keys were tested in one batch and its
classifier path was traced.

New controlled phrase contrasts show that zeroing `G83` changes selected
units and WAVs in `Does anybody know?`, `Does everybody know?`, and
`Can anybody help?`, while several matched pronoun and sentence controls are
byte-identical. In `Can anybody help?`, zeroing `G74` alone produces the same
change, while zeroing both codes restores the control exactly. Replacing
`G83` with `G95` changes `Does anybody see?` as well. Thus numeric codes can
affect audio in selected lexical/context combinations and can interact, though
the affected acoustic feature and predictive rule remain unknown. See the
[lead 3 TPP findings](lead3-tpp-typed-code-findings-2026-09-24.md) and
Stage 12 captures for the complete measured set and its limits.

## 4. Expand abbreviation probes by likely context

**Question.** Which of the 427 H/T/C abbreviation keys change with context,
case, punctuation, or neighboring token type?

**Starting evidence.** Every key has one runtime probe, and 40 targeted
contrasts cover selected honorific, company, date, measurement, address, and
repetition cases. The C sweep found 93 outputs different from its neutral
template, but a literal result in that template does not establish that a key
is inactive. See the [Stage 10 context report](abbreviation-context-contrasts-2026-09-24.md).

**Bounded probe.** Cover every inventory key with a baseline and controlled
context, punctuation, and case variants. Use family-specific contexts for H/T
and likely-reading contexts for C; keep the context-selection heuristic
visible where the inventory does not identify a more specific use. Capture
ordered token surfaces and phone/context rows.

**Useful result.** A per-key set of verified context contrasts, not a blanket
expansion table. Preserve unresolved keys and contexts explicitly.

**Bounded result (2026-09-25).** Original-engine captures now cover all 427
inventory keys: 55 H, 42 T, and 330 C. The 1,750 observations include each
key's baseline, a likely context, punctuation, and case contrasts, with full
returned rows and 24 phone/context bytes per row. Examples include `in.` →
`inches`, `mi.` → `miles`, `Ct.` → `Court`, title and ordinary `dr` readings
(`Doctor`/`Drive`), and the state/company contrast for `CO`. The C context
selector is an explicit heuristic; 21 C and 5 T alternate-case probes also
match a different key in their set. These controls are flagged in the
manifest. This closes the full inventory pass for the tested contexts; it
does not enumerate every possible neighbor or punctuation combination. The
[Lead 4 results](lead4-abbreviation-context-results-2026-09-25.md) link to the
per-key matrix, raw captures, and reproducible runner in Stage 15.

## 5. Compare other voices and package versions

**Question.** Which parts of the Paul findings generalize to the other
supplied voices or to other VoiceText package versions?

**Starting evidence.** The four supplied voice DLLs have matching executable
`.text` bytes in the recorded comparison, while small `.data` and voice/model
inputs differ. The deep runtime and format analysis covers Paul M16; the other
voice data and older index variants are not characterized to the same depth.
See the [asset map](../assets.md) and [open-topics review](open-topics-review-2026-09-24.md).

**Bounded probe.** If matching data packages are available, first inventory
file sets, hashes, headers, and tree/index layouts without modifying inputs.
Then run a short shared set of API and text probes against each voice and
compare load behavior, selected structural formats, and output properties.
Treat a matching code section as evidence of shared implementation code, not
proof that all voice behavior or data formats match.

**Useful result.** A compatibility matrix stating which layouts and behaviors
were checked for each package, and which findings remain Paul-specific.

## 6. Probe API behavior beyond the main text-to-WAV path

**Question.** What behavior does a native replacement need to preserve for
other output selectors, buffer/thread modes, and error cases?

**Starting evidence.** The public header and the 2013 Paul text-to-WAV path
are mapped, but runtime behavior across every selector, mode, and failure
path is not. See the API boundary in the [open-topics review](open-topics-review-2026-09-24.md)
and the recovered [engine findings](voice-engine-and-model-formats.md).

**Bounded probe.** First select the compatibility surface required by a
future implementation. For each in-scope API mode, record arguments,
return/error codes, output format, and whether state persists across calls.
Avoid probing unrelated modes without an implementation or research need.

**Useful result.** A tested API contract for the chosen scope, including
explicitly unsupported or untested cases.

**Bounded result (2026-09-25).** The selected scope is the local Paul M16
`VT_TextToFile_ENG` path. Runtime captures cover declared selectors 0–5 and
7–9, unsupported selector 6, out-of-range selector 10, and null-text,
empty-text, and null-path returns. Two identical selector 4 calls in one
loaded process both return success and produce byte-identical WAVE files. The
four buffer formats also return success at flag 0/thread ID 0 and match their
raw file-selector outputs; tested nonzero flag/thread combinations return the
header's create-thread error in Wine. Buffer null/empty/error cases and all
declared `VT_GetTTSInfo_ENG` requests were also probed. The output formats,
return codes, and limits are recorded in the [Lead 6 API report](lead6-file-api-behavior-2026-09-25.md),
with portable traces and outputs in Stage 16. Follow-up probes confirm that a
short synchronous buffer call treats `output_len` as the returned byte count,
not as a write limit; the configuration setters clamp upper values and alter
file output when synthesis arguments are `-1`; and valid playback returns the
initialization error in the local Wine runtime. This does not close the full
lead: successful threaded processing, polling, mixed-call state,
format-specific markup behavior, successful playback, per-field configuration
output effects, and other error paths remain untested. The declared text-format
values were tested only with a plain ASCII sample.

## Lower-priority leads and scope limits

- **DAT mode 8:** no mode-8 frame occurs in the local Paul corpus. Runtime
  validation needs a genuine payload from another package or an isolated
  decoder test that does not pretend to establish real-corpus compatibility.
- **Anonymous model feature labels:** more static inspection may refine field
  flow, but names are unlikely to emerge without discriminating controlled
  inputs, an independently labeled dataset, or a vendor schema.
- **Reuse rights:** technical feasibility does not establish rights to use or
  redistribute vendor code, voice data, or recordings. Any implementation
  plan needs its own licensing assumptions.

## Suggested first sequence

1. Continue the controlled tree-input contrasts, since fixtures and a tree
   trace already exist and the result can reduce the largest semantic gap.
2. In parallel at the research-planning level, define the small full-WAV
   golden corpus and its exact comparison boundaries.
3. Use the tree and golden-corpus observations to choose whether TPP semantics
   or abbreviation context breadth offers the next highest-value probe.
4. Broaden to other voices or API modes only when a concrete compatibility
   target calls for it.

This sequence is advisory. It does not start a new roadmap stage, authorize
changes to vendor inputs, or make a compatibility or licensing claim.
