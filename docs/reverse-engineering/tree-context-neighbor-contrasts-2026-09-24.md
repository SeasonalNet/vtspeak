# Controlled decision-tree neighbor contrasts (2026-09-24)

This report advances follow-up lead 1 for the local 2013 M16 Paul package.
The Stage 10 `tree-context-*` fixtures force phones around a fixed `probe`
span, record decision-tree inputs and returns, and continue through selected
IDs and timeline rows. These are controlled synthetic inputs; no vendor field
names or linguistic meanings are inferred from numeric values.

## Fixed `P AH0` contrasts

The first controls hold `probe` at `P AH0` and vary a forced phone immediately
before or after it. Each capture uses one original process and records the
ordered tree calls, all 16 signed input halfwords, and scalar returns. The
four returns below follow the common initial scalar call, which returns `30`.

| Fixture contrast | Next four scalar tree returns, in call order |
| --- | --- |
| left P, control | `958, 755, 1159, 2341` |
| left B | `1040, 1470, 1159, 2341` |
| left T | `1456, 631, 1159, 2341` |
| left F | `1494, 1019, 1159, 2341` |
| left M | `1604, 1019, 1159, 2341` |
| right P, control | `958, 755, 1159, 2341` |
| right B | `958, 755, 1240, 2341` |
| right T | `958, 755, 1824, 2341` |

Across the left controls, the first two calls' printed vectors and returns
change while calls three and four remain identical. Across the right controls,
calls one and two remain identical; calls three and four both change their
input vectors, although only call three changes its return in the P/B/T
controls. In the P/B pair, the changed phone ordinal is at printed vector
positions 0 and 1 in those two calls respectively (`27` for P, `7` for B).
The corresponding left pair changes those same positions in calls one and
two. Thus the input slots track the controlled phone identity, while the
unchanged fourth return shows that an input change need not change a tree
output. These traces do not give the scalar values duration/pitch units. The
raw tree-only captures are
`tools/revkit/work/stage10/tree-context-{left,right}-*-trees.log`.

Four P/B controls were rerun with traces at `FUN_1001b200` and
`FUN_1002c220`. Each process exited normally and returned six primary rows.
The first six selected IDs precede timeline construction and match its six
rows in order. The table shows those IDs and the normal-row sample counts
(`+0x0c`) in the same order:

| Fixture | Six primary selected IDs | Six timeline sample counts |
| --- | --- | --- |
| left P | `134486, 134487, 81307, 568612, 568613, 281431` | `968, 580, 376, 1730, 610, 1662` |
| left B | `234256, 234257, 76329, 568612, 568613, 281431` | `1386, 470, 488, 1730, 610, 1662` |
| right P | `134486, 134487, 81307, 568612, 568613, 281431` | `968, 580, 376, 1730, 610, 1662` |
| right B | `14169, 14170, 3838, 94614, 94615, 248840` | `1302, 276, 1090, 1574, 368, 1876` |

Changing the left phone from P to B changes the first three primary rows;
the last three IDs and sample counts remain equal in this pair. Changing the
right phone from P to B changes all six primary rows. In the right pair, the
two tree calls associated with the unchanged `probe` span retain the same
printed inputs and scalar returns, while the first three selected IDs and
timeline rows differ. The selected-unit path therefore responds to the
right-side contrast without a change in those two tree results. The traces do
not identify which candidate feature causes this effect or establish a
universal left/right rule.

## Repeat with fixed `T AH0` and a sentence boundary

The P/B contrasts were repeated around fixed `T AH0`. After the common
initial call, the four tree returns are `958, 755, 1824, 2341` versus
`1040, 1470, 1824, 2341` when the left phone changes P to B. For the right
contrast they are `1456, 631, 1159, 2341` versus `1456, 631, 1240, 2341`.
As with `P AH0`, the first two calls for the fixed phone are unchanged when
the right neighbor changes, but the first timeline call's first three
selected IDs differ. It has six rows for right P and five for right B; this
comparison is limited to these observed ordered selections and counts.

For the boundary controls, fixed `T AH0` ends at a period before the P/B
right span. The first timeline call is identical in both runs: three rows
with IDs `3837, 157573, 157573` and sample counts `1104, 956, 1370`. The
following timeline call has four rows in both runs; its final three sample
counts are `728, 1324, 1390` for P and `1198, 1236, 1564` for B. In this
fixture, the right-side change does not alter the preceding span's first
timeline call; the changed rows appear in the following call. This does not
establish behavior at every punctuation or phrase boundary.

All ten downstream captures exited normally. The reproducing trace and
commands are in the [Stage 10 notes](../../tools/revkit/work/stage10/README.md);
the trace is `tools/revkit/work/stage10/trace-context-downstream.gdb`, and
the matching logs and input fixtures are beside it. Timeline field offsets
and their established uses are in the
[Stage 8 findings](voice-engine-and-model-formats.md#selected-records-and-timing).

## Candidate-scoring continuation

The four fixed-`T AH0` right/boundary controls were rerun with
`trace-context-scoring.gdb`. The trace records each context's candidate IDs,
per-unit local costs, transition cumulative costs and predecessor indexes,
pruning results, and backtracked units. A focused pass also records the raw
context row and scorer inputs for the path candidates. All four processes
exited normally. The detailed captures are
`tools/revkit/work/stage10/tree-context-t-{right,boundary}-{p,b}-scoring-row.log`;
the broader candidate captures are in the matching `*-scoring.log` files.

For the unbroken right-side P/B pair, context 0 enters candidate generation
with the same single class ID (`61497`) on both runs. Its resulting unit
candidate sets contain 28 common IDs and two different IDs on each side.
Candidate `3837`, present in both sets, has local cost `0.833333` and
context-penalty field `0` for right P; for right B it has local cost
`0.416667` and context-penalty field `1`. It is ranked at local candidate
index 27 in P and index 0 in B. The field labels here follow the traced
record offsets; their semantic meaning remains unknown. For candidate
`3837`, the logged six-byte context row is the first six bytes of the
scorer's seven-byte target at the same address; those bytes are identical
across runs (`00 00 00 00 00 01`), and the target's seventh byte is `00` in
both. The separate ten-byte feature view is identical too. Its scorer scale
changes from `1` in P to `0.5` in B, and its returned local cost is halved
accordingly. The recorded duration-penalty field stays `2`. The Stage 7
Ghidra pseudocode for
`FUN_10023350` forms this scale using the candidate node's `+0x10` field and
a normalized duration term; this controlled pair therefore ties the observed
scale change to the changed `+0x10` value, while the upstream source of that
field is still unknown. This explains the score change for this candidate in
this pair, not a general linguistic meaning for the field.

The transition at context 1 selects different predecessors in those same
runs:

| Right phone | Selected context 1 unit | Predecessor index | Context 0 unit at that index | Cumulative cost |
| --- | ---: | ---: | ---: | ---: |
| P | `152012` (candidate 12) | `12` | `152011` | `10.6287` |
| B | `3838` (candidate 0) | `0` | `3837` | `12.4167` |

The backtracker returns those context 0/1 pairs: `152011, 152012` for right
P and `3837, 3838` for right B. This traces the observed previous-unit change
to the candidate path selected after the neighboring context is scored. The
fixed span's decision-tree returns did not change, but candidate ranking and
the minimum-cost predecessor path did.

## Candidate metadata construction and second phone pair

The follow-up trace breaks at `FUN_100230a0` (`vt_pau.dll` address
`0x100230a0`), which runs after unit IDs are expanded and before local costs
are scored. For candidate `3837`, the P and B runs both enter this function
with `+0x0e = 0`, `+0x10 = 0`, and `+0x0c = 0`. On return, P has
`+0x0e = 2`, `+0x10 = 0`, `+0x0c = 1`; B has `+0x0e = 2`, `+0x10 = 1`,
`+0x0c = 1`. The right-side contrast therefore changes the value during
this metadata pass, rather than supplying a different initial value on the
candidate node.

For candidate `3837` at context 0, the focused trace reaches one forward
lookup hit at context index 1 in each of the four P/B/F/V runs. At
`0x1002326b`, byte `+0x04` of the six-byte context row (addressed as
`state + 0xec62c + index * 6`, with the row base at
`state + 0xec628 + index * 6`) indexes a 16-bit table local to
`FUN_100230a0`. The observed values are:

| Right phone | Context row at context 1 | Byte `+0x04` | Table contribution | Accumulator | Stored `+0x10` |
| --- | --- | ---: | ---: | ---: | ---: |
| P | `00 00 01 01 01 00` | `1` | `1` | `1` | `0` |
| B | `00 00 01 01 00 00` | `0` | `2` | `2` | `1` |
| F | `00 00 01 01 01 00` | `1` | `1` | `1` | `0` |
| V | `00 00 01 01 01 00` | `1` | `1` | `1` | `0` |

At `0x10023270` the contribution is added to the forward accumulator. For
this context-0 boundary case, `0x100232b5` halves that accumulator and
`0x100232ea` stores it at node `+0x10`. This directly accounts for the
P/B field difference and the unchanged F/V field in the tested pairs. The
byte `+0x04` is also read by the Stage 7 transition routine; this probe
establishes another operation that uses its raw value but does not assign
the byte a semantic name. The codes and local table's purpose remain
unnamed. The disassembly at `0x10023120–0x1002327b` also shows the routine
walking earlier and later context records; the endpoint branches select how
those side accumulators contribute to `+0x10` elsewhere. The matching
portable disassembly excerpt is
`tools/revkit/work/reports/stage10-candidate-metadata-disassembly.txt`; the
before/after records and focused row/table traces are in the
`*-scoring-row4.log` captures.

The source of row byte `+0x04` is now localized to a conditional fallback in
the `FUN_10024680` pseudocode. For each eligible context, it first calls
`FUN_10024060`; when that returns false, it calls `FUN_100242a0`. The latter
sets the current row's byte `+0x04` to `1` and, when constructing its second
row, copies the first row and sets that copied row's byte to `2`. The
`FUN_10024060` pseudocode clears the current row's byte to `0` before
candidate generation. A hardware watchpoint on context 1's byte confirmed
the P run writes `0 → 1` at instruction `0x10024305` in `FUN_100242a0`; the
right-B run made no change to that byte. A separate entry trace records the
fallback indices: right P, F, and V call `FUN_100242a0` for indices 1 and
4, while right B calls it for index 3. Thus the observed context-1 values
(`1` for P/F/V, `0` for B) align with whether that context uses the
fallback. A return-site trace at `0x10024808` confirms context 1 returns
`0` for P/F/V and `1` for B; the caller tests that result at `0x1002480b`
and branches around the fallback on a positive value. At context 3 the
returns are `1` for P/F/V and `0` for B, matching the reverse fallback
pattern. This identifies the immediate branch that distinguishes these
rows. The count trace follows candidate generation through `FUN_10018770`,
ranking at `FUN_10023af0`, and the final `FUN_10023060` sum. For these
fixtures the model flag tested by `FUN_10024060` is `0`, so the routine
returns true when that sum is at least 10; the special flag-8 path uses a
positive-sum check instead. `FUN_10023060` sums 16-bit values from the table
at model offset `+0x8c` for the candidate IDs it receives.

| Context | Right phone | Generated count | Candidate class IDs passed to `FUN_10023060` (table values) | Sum | Return |
| ---: | --- | ---: | --- | ---: | ---: |
| 1 | P | 1 | `48721` (`1`) | 1 | 0, fallback |
| 1 | B | 2 | `48655` (`10`), `48656` (`4`) | 14 | 1, no fallback |
| 1 | F/V | 0 | none; returns before ranking | — | 0, fallback |
| 3 | P | 3 | `3111` (`11`), `3113` (`65`), `3628` (`37`) | 113 | 1, no fallback |
| 3 | B | 1 | `9874` (`3`) | 3 | 0, fallback |
| 3 | F | 2 | `2369` (`64`), `2372` (`77`) | 141 | 1, no fallback |
| 3 | V | 2 | `3451` (`171`), `3454` (`17`) | 188 | 1, no fallback |

The measured candidate lists and sums account for the immediate fallback
choices, but not for why `FUN_10018770` yields those IDs or why their model
table values differ. A separate entry trace shows that the seven-byte
`FUN_10018770` input at context 1 is `[90, 57, 7, X, 7, 97, 0]`; only byte
`+3` varies across the controls (`X=53` for P, `19` for B, `32` for F, and
`65` for V). The six-byte context row at that entry is identical in all four
runs (`00 00 01 01 00 00`). The `FUN_10018770` translation-table checks map
bytes `+2 = 7` and `+1 = 57` to `0` and `77` in all four runs. Byte `+3`
maps as follows: P's `53` maps to `71`; B's `19`, F's `32`, and V's `65` map
to `0`. This routes P through the nonzero branch, whose three branch-specific
`FUN_10023f90` lookups return zero in these runs. B/F/V take the zero route
and make the corresponding lookup using mapped byte `77`, also returning
zero. Separately, successful combined-signature queries emit P candidate
`48721` from `[90,57,7,53,7,16,0]`, and B candidates `48656` and `48655`
from `[90,57,7,19,7,32,0]` and `[90,57,7,19,7,16,0]`, respectively. F/V
have no successful combined query and emit no candidates. These emissions
match the class IDs passed to the first-pass scoring sum. The local
phone-symbol codebook maps the changed raw values to P/B/F/V; the structural
slot mapping and the suffix bytes' field semantics remain qualified below.

Static disassembly shows `FUN_10023f90` dispatching to `FUN_10023dc0` when
context-row byte `+0x04` is zero. `FUN_10023dc0` converts the seven-byte
signature through `FUN_10016ea0`: the first three output bytes come from
three lookup tables indexed by signature bytes `+1`, `+2`, and `+3`; output
byte 3 copies signature byte `+5`; byte 4 is signature byte `+6` masked by
`0x20`. It passes this five-byte key to `FUN_10019450`, which binary-searches
sorted five-byte records rooted at model pointer `+0x88`; the bytewise
comparator is `FUN_10019350`, and an equal match expands to its contiguous
range. A focused runtime trace reports a table count of 61,566 and these
successful one-record matches:

| Seven-byte signature | Five-byte key | Range count | First class index |
| --- | --- | ---: | ---: |
| `[90,57,7,53,7,32,0]` (P, suffix 32) | `[57,7,53,32,0]` | 0 | — |
| `[90,57,7,53,7,16,0]` (P, suffix 16) | `[57,7,53,16,0]` | 1 | 48721 |
| `[90,57,7,19,7,32,0]` (B, suffix 32) | `[57,7,19,32,0]` | 1 | 48656 |
| `[90,57,7,19,7,16,0]` (B, suffix 16) | `[57,7,19,16,0]` | 1 | 48655 |
| `[90,57,7,32,7,32,0]` (F, suffix 32) | `[57,7,32,32,0]` | 0 | — |
| `[90,57,7,32,7,16,0]` (F, suffix 16) | `[57,7,32,16,0]` | 0 | — |
| `[90,57,7,65,7,32,0]` (V, suffix 32) | `[57,7,65,32,0]` | 0 | — |
| `[90,57,7,65,7,16,0]` (V, suffix 16) | `[57,7,65,16,0]` | 0 | — |

For each successful query, the one-record range index equals the generated
class ID recorded in the first-pass list; these are class indices at this
stage, not the final waveform-unit IDs. Thus the class lookup accounts for
the tested first-pass counts: P contributes one class, B two, and F/V none.
The local phone-symbol codebook identifies byte values `53/19/32/65` as
CMU-supported P/B/F/V codes. Wag's separate `vt_eng.dll` 3.11.7.1 notes label
the corresponding signature slots `L2/L1/C/R1/R2/posbits/flags`, placing
byte `+3` at `R1`; this is a cross-build corroboration, since these runtime
captures use our local `vt_pau.dll`. The key suffix bytes `16` and `32` and
their full field semantics remain unresolved. The trace scripts and logs are
`trace-context-row-writer.gdb`, `trace-context-row-build.gdb`,
`trace-context-firstpass.gdb`, `trace-context-firstpass-counts.gdb`,
`trace-context-firstpass-inputs.gdb`, `trace-context-firstpass-byte-branch.gdb`,
`trace-context-firstpass-class-ranges.gdb`, and the matching
`*-row-writer-verified.log`, `*-row-build.log`, `*-firstpass.log`,
`*-firstpass-audit.log`, `*-firstpass-inputs.log`, `*-byte-path-emitted.log`,
and `*-class-ranges.log` files under
`tools/revkit/work/stage10/`.

The follow-up trace at `FUN_10023350` reads each candidate class's 16-bit
population from the model count table and its member-unit list. In the right-B
control, the exact context-1 class set expands as follows:

| Class ID | Count-table value | Expanded waveform-unit IDs |
| ---: | ---: | --- |
| `48655` | 10 | `107607, 142293, 163010, 212458, 217023, 217816, 241367, 270718, 304501, 307897` |
| `48656` | 4 | `3838, 153466, 424604, 577216` |

The two values summed by `FUN_10023060` therefore equal these classes'
expanded member counts in this runtime. The P capture also shows class
`48721` has count 1 and expands to unit `246151` when present in the broader
context-1 pool. This identifies the direct class-to-unit mapping for these
captures; it does not identify the phonetic interpretation of the unit
records. The GDB trace and captured class expansions are in
`trace-context-unit-expansion.gdb` and `*-unit-expansion.log` under
`tools/revkit/work/stage10/`.

## Cross-check against Wag's independent port

Wag's response maps target-row byte `+0x04` to `half`: `0` is a whole-phone
target, `1` the left half, and `2` the right half. It calls
`FUN_10024060` the whole-phone lookup and `FUN_100242a0` the half-phone
fallback. This interpretation matches our local pseudocode and traces:
`FUN_10024060` initializes the target byte to zero and returns false when the
merged class population is at most nine; its caller then invokes
`FUN_100242a0`, which writes `1` and `2` to two adjacent target rows. The
runtime P/B contrast exercises both outcomes: the P whole-phone result sums
to one and splits, while B sums to fourteen and remains whole.

The local disassembly of `FUN_100230a0` independently supports Wag's
contiguity interpretation. It initializes the half-indexed weights to
`[2, 1, 1]`, walks neighboring unit IDs when the unit label's `0x80` bit
permits, and checks membership in neighboring target beams/class lists. It
stores the resulting `run` at candidate `+0x0e` and `wrun` at `+0x10`, then
sorts by `-(run + 100*wrun)`. In the traced context-0 case for unit `3837`,
there are no backward targets, so the branch sets `wrun` to the forward
weighted run divided by two. The measured `+0x10` change and halved local
scorer scale therefore fit this mechanism. More generally, `+0x10` is the
combined `wrun` selected from forward/backward runs; calling it simply
`wfwd / 2` applies to this context-0 case, not every target.

Wag's second reply accepts the combined-`wrun` interpretation and provides
the four endpoint branches. Register flow in the local disassembly matches
those branches. Let `back` and `fwd` be the contiguous target counts and `wb`
and `wf` their half-scaled weighted sums. The stored run is `back + 1 + fwd`;
`wrun` is `wf + 1 + wb` when the run covers all targets, `wf` when
`back == target_index`, `min(wb, wf)` when `fwd` does not reach the final
target, and `wb` otherwise. This records an independently supplied function
reading corroborated against the local instruction sequence.

The labels `whole phone`, `left half`, and `right half` are operational names
from Wag's port and are strongly supported by the local two-target split and
the target-index updates. They are not claimed as original vendor symbols.
Wag's second reply corrects its earlier `.text` identity claim. It reports
three changed bytes at `VT_CheckLicense_ENG` (`0x10029b80`), a 797-byte `.data`
run containing demo-notice strings, and identical `.rdata`, `.rsrc`, and
`.reloc` sections with matching section addresses and sizes. It also notes
that the license checker is called by `VT_LOADTTS_EXT_ENG` and
`VT_GetTTSInfo_ENG`. The response bundle omits both DLLs, so these precise
binary differences and call paths remain Wag-reported rather than locally
verified. They qualify the cross-build address comparison and warn that those
license-facing API paths can behave differently; our synthesis-path traces do
not exercise them.

### Direct half-key lookup trace

A runtime trace through `FUN_10019570` captured the transformed 10-byte keys
passed to `FUN_1002df50` for the split P/F/V targets. The trace follows raw
labels `[90,57,7,X,7,97,0]` with `X = 53` (P), `32` (F), or `65` (V). In
byte order, the left-half keys are P/F/V respectively
`[0,0,7,4,0,57,40,57,1,4]`,
`[0,0,7,4,0,57,40,57,1,4]`, and
`[0,0,7,4,0,57,40,57,1,3]`. Their first eight values are identical; the
last two are the mapped phone-dependent values. The right-half keys are
`[1,0,7,4,0,53,10,53,4,4]`,
`[1,0,7,4,0,19,10,32,4,4]`, and
`[1,0,7,3,0,65,10,65,4,4]`.

The search and ranking results differ by side and phone. For each P/F/V run,
the left-side query ultimately supplies eight exact classes to the first
candidate pass (40 concrete units by the captured class-population table).
The P right side first yields two exact classes (4 units), then a widened
query supplies ten classes (45 units); F and V each reach ten classes on the
right, with captured totals of 84 and 67 units. These are the candidate
arrays at the traced call sites, not claims that the 10-byte keys alone
determine final unit selection. The later ranker expands each side to its
30-entry candidate cap. The right-side B control emits no half-key lookup in
this path: its whole-phone lookup passes the population threshold and is not
split.

The local trace prints raw key bytes, query range counts, class IDs, and unit
populations without interpreting the key slots. In particular, the changed
byte values establish different operational keys for P/F/V; they do not
assign names to those fields. Reproduction instructions and complete captures
are in `trace-context-half-keys.gdb` and the four
`*-half-keys-verified.log` files under `tools/revkit/work/stage10/`.

A further K/G pair extends this comparison beyond P/B/F/V. K's left key is the
same `[0,0,7,4,0,57,40,57,1,4]`; its right key is
`[1,0,7,4,0,42,10,42,4,4]`. The K left first-pass pool again has eight
classes/40 units. Its right path first yields five classes/5 units and later
reaches ten classes/30 units. The adjacent K/G whole-phone trace returns zero
for context 1 in K and one in G; K enters the half path, while G has no
half-key calls. Thus, within these controls, a shared left key can coexist
with a changed right key and different split behavior. The traces do not
establish why K and G differ at the whole-phone lookup.

The existing left-neighbor P/B controls show the complementary context-1
result. With `P AH0` or `B AH0` before the fixed `T AH0` target, the whole-phone
pass returns `1` for P and `0` for B; for the right-neighbor pair after the
fixed target, the returns were `0` for P and `1` for B. Only left-B reaches
the traced half-key lookup in the left pair. Its label is
`[90,19,7,57,7,97,0]`; the left and right keys are
`[0,0,7,3,0,19,40,19,1,4]` and
`[1,0,7,4,0,57,10,57,4,3]`. The left candidate pool has six classes/13 units.
The right pool first has four classes/5 units, then five/6 after a broader
query, and ten/29 at the final pool. This orientation reversal is measured
only for the two controlled pairs; it does not establish a general P/B rule.
The full return, input, and half-key captures are the
`tree-context-t-left-*-firstpass-left.log`,
`tree-context-t-left-*-firstpass-inputs-left.log`, and
`tree-context-t-left-*-half-keys-left.log` files under
`tools/revkit/work/stage10/`.

A four-case three-phone control places `B AH0` before `T AH0` and varies the
phone after T across P/B/F/V. The context-1 input and row remain identical
across all four (`[90,19,7,57,7,98,0]`, row `00 00 01 01 00 00`) and all
return zero. At context 4, the row also remains identical
(`01 00 01 01 00 00`), while input byte `+3` changes among 53 (P), 19 (B),
32 (F), and 65 (V). P/F/V return zero; B returns one. For P, the half-key
trace records a left key `[0,0,7,4,0,57,20,57,1,4]` with an eight-record
final range and a right key `[1,0,7,4,0,53,10,53,2,4]` with a final 92-record
range. The ranked left pool contains eight classes/231 units; the right pool
has ten classes/17 units before the subsequent ranking pool. F and V also
take the half path; their full keys and query ranges are preserved in the raw
captures. B takes the whole-phone path and makes no matching context-4
half-key query. These are numeric lookup observations, not semantic names
for the key bytes. Other context indices differ in their input alignment
because the probe is embedded in a three-phone sequence, so the result is
limited to these controlled context-4 cases. Captures are the four
`*-middle-audit-all.log` and `*-half-keys-all.log` files under
`tools/revkit/work/stage10/`.

The duplicate-key question is now checked across the loaded class table, not
only the example query results. A read-only GDB Python scan bulk-read all
61,566 sorted five-byte records and found zero repeated-key runs. Therefore
the observed one-record equal-key results are representative of key
uniqueness in this Paul M16 table; this says nothing about other model
versions. The scan is reproducible with `trace-context-class-duplicates.gdb`,
and its capture is `tree-context-t-right-p-class-duplicates.log` in
`tools/revkit/work/stage10/`.

A second controlled right-phone pair holds the same `T AH0` span fixed and
changes the following phone from F to V. In context 0, all 30 candidate IDs
and local costs are identical between runs. Candidate `3837` has
`+0x10 = 0` before and after `FUN_100230a0` in both, and its local cost is
`0.833333` in both. The backtracker nevertheless returns different context
0/1 pairs: `256828, 256829` for F and `160452, 160453` for V. Thus this
second pair changes the selected path without changing this candidate's
metadata or the fixed context's local-cost array. The evidence bounds the
P/B field change to that contrast; it does not establish that all right-phone
changes modify `+0x10`.

For the period-separated P/B controls, the first sentence's contexts 0–2
have identical candidate local-cost arrays, transition results, and
backtracked units in both runs (`3837, 157573, 157573`). The engine calls the
backtracker twice per run, each time over contexts 0–2; the second span's
candidate lists and selected path differ across P/B. In these captures the
period therefore separates the dynamic-programming spans, and no changed
right-phone candidate is scored as a continuation of the preceding span.
This is a result for the tested period control, not every punctuation or
phrase boundary.

## Conclusion and limits

The measurements locate the tested selection effect in candidate
scoring/transition behavior. Local disassembly plus Wag's independent port
identify byte `+0x04` operationally as whole/left-half/right-half target mode,
and identify the 16-bit values read by `FUN_10023060` as class populations.
The local `FUN_100230a0` disassembly and traces tie context-0 candidate
`3837`'s `+0x10` to weighted forward contiguity in this case; other target
positions can combine forward and backward runs. At context 1, the first-pass
input differs only at byte `+3`, and the traces show its table mapping,
branch route, emitted IDs, transformed keys, matching model records, and
class-to-unit expansion. The tested successful equal-key ranges each contain
one record. The direct half-key trace now shows the split P/F/V key bytes and
first-pass candidate pools; B remains on the whole-phone path. K/G adds a
second contrast: K splits and G passes the whole-phone first pass. A full
scan of the 61,566-record Paul class table found no duplicate five-byte keys.
The left-side P/B control reverses the context-1 whole-phone return seen in
the right-side control, and left B's two half keys are now captured. A
three-phone B/T/P/B/F/V set adds a context-4 split: P/F/V take the half path
and B the whole-phone path with identical six-byte context rows. The
[lead 1 closure](lead1-tree-input-closure-2026-09-24.md) records the bounded
result and its counterexamples. Additional phone pairs and context positions
would extend that result to new contexts or package versions.
