# Stage 20: matched 2006 and 2013 tree-output trace

This read-only analysis compares the original Kate MSI engine and package with
the Stage 19 adapted 2013 engine and package. The runtime runners temporarily
replace Stage 5 synthesis fixtures, preserve the prior files, and verify
restoration. Vendor binaries and voice assets remain mounted read-only.

## Plain-text differential

The controlled input is `hello-plain.txt` (`Hello.`). Both runners accept an
`INPUT_FIXTURE` override so the same bytes reach each engine:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage20/compose-original-msi.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt runtime \
  /bin/bash /work/stage20/run-original-forced-pah0.sh

docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt runtime \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The old-engine GDB script traces parser entries, scalar recursive tree
evaluation at `FUN_10001570`, and vector evaluation at `FUN_100019b0`.
The adapted-engine script traces scalar/vector tree calls, selected unit IDs,
and timeline rows. Captures and WAVs are retained under the ignored
`tools/revkit/work/corpus-parity/stage20/` directory:

- `old-hello-gdb.log`, `old-hello-winedbg.log`, `old-hello.wav`
- `new-hello-gdb.log`, `new-hello-winedbg.log`, `new-hello.wav`

For this input, the first four old and new scalar tree returns match
(`1207, 555, 1740, 4832`). The first two old vector calls return 12 shorts,
while the paired new calls return 12 shorts with only the first value
populated. This is one observed input, not a whole-voice result. The old
runtime trace maps those calls to the engine's `nbf` and `sbf` slots; static
loader setup maps those slots to `pitch/nbf.tree2` and `pitch/sbf.tree2`.

The old loader reads byte 4 of each file as vector width 12 and uses the
special `FUN_10001c80` loader for these four pitch resources (`nbf`, `bf`,
`qbf`, `sbf`). The Stage 19 generic `tree2.py` parser reads bytes 4–5 as a
scalar output instead. For `nbf` and `sbf` this produces `0x500c` (20492) and
`0x600c` (24588), then treats the next bytes as a leaf marker and stops at
byte 9. The listed 19,931 and 3,691 “trailing bytes” are therefore almost the
entire special-format trees, not opaque tails after a complete parse. Stage 19
converts each mistaken scalar to a one-output leaf; the 2013 vector call
returns that value in its first position with the rest zero-filled. The old
trace script prints only the 12 shorts copied by the legacy evaluator.

This identifies a specific loader-format mismatch and explains the observed
vector outputs. It does not establish that correcting these four pitch
resources alone makes full sentences intelligible; the experimental indexed
conversion and its one-input effect are documented below.

## Experimental vector-pitch adapter

`convert_vector_pitch.py` implements the special 2006 node layout for the four
files, requires each parse to end exactly at EOF, preserves each terminal
12-short row, and emits indexed tree3 rows with width 12. It copies the other
four converted pitch trees into an isolated overlay so the engine still sees
the complete directory. The four files parse as follows:

| File | Parsed nodes | Indexed nodes | Leaves | Converted bytes |
| --- | ---: | ---: | ---: | ---: |
| `nbf.tree2` | 561 | 280 | 281 | 10,471 |
| `bf.tree2` | 21 | 10 | 11 | 425 |
| `qbf.tree2` | 7 | 3 | 4 | 140 |
| `sbf.tree2` | 99 | 49 | 50 | 1,948 |

Rebuild the overlay and run the same plain-text fixture with the adapted
engine:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 -B \
  tools/revkit/work/stage20/convert_vector_pitch.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt runtime \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The 2013 evaluator returned exactly the same two 12-short rows as the 2006
evaluator for the traced `nbf` and `sbf` contexts. The first four scalar
returns also remained `1207, 555, 1740, 4832`. Candidate unit IDs changed at
positions 3–5: baseline `265071, 278054, 71004`, converted `35672, 35672,
263578`; the other first-eight IDs stayed the same. The adapter WAV is 5,206
frames (10,456 bytes, SHA-256
`928cd898d658f422a40e8b91fc82ebe6e373a4e239c864072c5a81b36de5db36`), versus
4,798 frames for the unadapted 2013 package and 7,610 for the 2006 engine.
Those duration changes show the recovered vectors affect synthesis selection;
they do not establish that the resulting sentence is intelligible. The
captured trial files are `new-hello-vector-pitch-gdb.log`,
`new-hello-vector-pitch-winedbg.log`, and `new-hello-vector-pitch.wav`. Input
and output restoration hashes are in `adapted-forced-pah0-restore.txt`.

The requester listened to this vector-pitch WAV and described it as sounding
like “Shar” or similar, rather than intelligible “Hello.” The requester had
previously judged the selected-side waveform pieces to sound like the intended
phones. These are listening observations, not a phone-level transcription.
Together with the exact vector-row match, they suggest the remaining audible
failure is later in unit choice, ordering, or timing/coarticulation; that is an
inference, not a localized cause.

The adapter was also exercised on the existing `prose`, `numbers`, and
`address` fixtures. All three runs completed and produced different WAVs from
the baseline zero-filled-vector package:

| Fixture | Baseline bytes | Recovered-vector bytes | Recovered-vector SHA-256 |
| --- | ---: | ---: | --- |
| Prose | 100,144 | 87,100 | `bf6438c44103ec0261bae6fc151814da9b393d8f4d7032ca0600c27fc35c5901` |
| Numbers | 91,314 | 85,632 | `7ab68e0650574b2d23159712a9c67bc087c929a41ea5dc70b8e6dcd2e59f1a1c` |
| Address | 104,326 | 97,870 | `69a32ffd49e917870dc59d75e7b24f7d8dd851b2a047a36f13b1e97067153a9b` |

This broadens the measured effect of the vector-tree conversion beyond one
sentence, but duration and hashes do not show that speech became more
intelligible. The matrix and outputs are retained under
`tools/revkit/work/stage19/runtime-vector-pitch-matrix.tsv` and
`runtime-vector-pitch-{prose,numbers,address}.wav`.

An additional controlled `Hello.` run held the vector overlay fixed and used
the Stage 19 `signature-last-copy` index candidate for the inserted 2013
`attr_b`. It returned the same 12-value `nbf`/`sbf` rows and the same first
four scalar outputs as the zero-filled candidate, but changed the first eight
selected unit IDs from
`280411, 276348, 35672, 35672, 263578, 213073, 249391, 249391` to
`280411, 172786, 270092, 80917, 224948, 40384, 56352, 56352`. Its WAV is
`new-hello-vector-pitch-signature-last-copy.wav` in the ignored Stage 20
evidence directory (5,303 frames; SHA-256
`e86d005f2be73d6c93b9018d833cf38fb8bb3619763b17e60cf01de3b3adf1c8`). This
isolates a substantial selection effect from the index candidate while keeping
the corrected vector outputs constant; it does not establish that the
candidate units or the resulting speech are better. The candidate's field
meaning is still unknown.

After listening to that `signature-last-copy` output, the requester described
it as sounding like “Hello” with a major speech impediment, approximately
“helall.” This is a positive subjective difference from the zero-filled
vector-pitch output described as “Shar,” but it does not validate the mapping
or identify the remaining defect.

Two additional `Hello.` controls held the recovered vector trees fixed:

| Index candidate | First eight selected unit IDs | WAV frames | SHA-256 |
| --- | --- | ---: | --- |
| `signature-last-copy` | `280411, 172786, 270092, 80917, 224948, 40384, 56352, 56352` | 5,303 | `e86d005f2be73d6c93b9018d833cf38fb8bb3619763b17e60cf01de3b3adf1c8` |
| `signature-last-transfer` | `280529, 91559, 204071, 238753, 71004, 40384, 56352, 56352` | 4,959 | `dec81198b7a9672e73ac14810c948d31bb03d1867afdc09c946c22fb60b547da` |
| `constant-8` | `280411, 276348, 35672, 35672, 263578, 213073, 249391, 249391` | 5,206 | `928cd898d658f422a40e8b91fc82ebe6e373a4e239c864072c5a81b36de5db36` |

The transfer candidate copies a byte window beginning at legacy tail offset
`7 * unit_count` into `attr_b`, then zeroes that source window. It produces a
different selection from copying the window while preserving it. The
constant-8 output is byte-identical to the zero-fill baseline for this input,
despite filling the target value observed in one scoring trace. These outcomes
show that the inserted byte changes scoring and that zeroing part of the
source signature buffer also changes the earlier class query. The transfer WAV is
`hello-vector-signature-transfer.wav`; the constant-8 WAV is
`hello-vector-constant-8.wav` in the ignored Stage 20 evidence directory.
The requester describes the transfer WAV as colder than the copy version, with
an unclear initial fragment and an ending resembling “nall.” Those are
subjective comparisons, not recovered phone labels.

The `signature-last-copy` candidate was also run across the three existing
sentence fixtures with corrected vector trees. The run reproduced the
previously captured hashes:

| Fixture | Bytes | SHA-256 |
| --- | ---: | --- |
| Prose | 76,430 | `303e96eb3ecb58a9cc30b79f838a6a37fd355cdcf50bd468fd9a4e16f5115dcc` |
| Numbers | 69,542 | `818465e8ea5535326b69b990a0f67c00cfddcbb1c58ff2c37446022cc685365d` |
| Address | 83,924 | `bd6b6d97b50470ef17c70db1f0cd4655983ecc7fde5b1995eafe750b31714743` |

The reproduced matrix files are under
`tools/revkit/work/stage19/runtime-vector-signature-last-copy-*`; the earlier
captures are under `runtime-vector-pitch-signature-last-copy-*` and have the
same hashes. Reproduction confirms deterministic output for this setup; the
full-text outputs have not yet been assessed for intelligibility.

### What `attr_b` changes in the 2013 scorer

Matched GDB traces on `Hello.` captured the key-class shortlist, later unit
ranking, and `FUN_100182e0` pair costs for zero-fill and the copy-window
candidate
with the same recovered vector trees. Both variants use the identical target
signature `22 17 2b 30 5a 00 00`, return the same class ID (`46204`), and pass
the same first three final-ranking input lists with counts 10, 19, and 20.
Thus this copy candidate does not change the parser, tree traversal, or the
initial class shortlist on this input.

It does change the per-unit score inputs. With zero-fill, all candidates in a
given traced target context get the same `attr_b` pair cost: `1.33333` for
target value 12 and `6` for target value 6. Copying the legacy byte window produces
candidate-specific costs for those same targets. For example, in the first
target-6 group, candidate unit `35672` costs `59.6471`, `179163` costs `1.2`,
and `270092` costs `5.15789`. The first three final-ranking calls keep the
same candidate lists and counts in both variants, while returned score values
change. The eventual selected-unit sequence also changes, as shown above.
This directly locates the copy effect in later candidate scoring/selection,
after the initial class shortlist. It explains why zero-fill gives the scorer
no discrimination on this dimension; it does not establish that the copied
window has the intended semantics in `attr_b`.

The transfer control has the same inserted `attr_b` values as copy, but zeros
the source byte window. On this same `Hello.` input, its first key-class
query returns class ID `23148` instead of `46204`; its first three final-rank
candidate counts are 12, 30, and 14 instead of 11, 19, and 21 for copy. The
seven-byte target signature remains unchanged, while part of the adapted
candidate-signature buffer has been zeroed. This changes the 2013
class/selection path. The requester's colder
listening impression for transfer is consistent with this earlier shortlist
change, but does not prove the old and new class semantics are identical.

The captures are `hello-keyrank-gdb.log`, `hello-keyrank-zero-gdb.log`, and
`hello-keyrank-transfer-gdb.log` in the ignored Stage 20 evidence directory.
The paired `hello-vector-*-score-gdb.log` files include the per-unit pair-cost
rows and selected unit IDs.

## Causal assessment and highest-value follow-ups

### Index-layout correction

The native 2006 versioned-index reader `FUN_100123a0` reads the legacy feature
tail in this order: `attr_4c[N]`, `attr_48[N]`, `key[5N]`, `attr_40[N]`, then
three metric groups totaling `12N`. The 2013 reader consumes
`attr_a[N]`, a unit-major `signature[7N]`, `attr_b[N]`, then the same-width
metric area. The Stage 19 adapter only inserted `N` zero bytes at offset
`8N`; it did not transpose or rebuild the legacy fields. Consequently its
2013 signature is a bulk concatenation of `attr_48`, all five key columns,
and `attr_40`, which the 2013 engine then chunks every seven bytes as though
each chunk belonged to one unit. That is the concrete index-layout defect.

The 2013 key builder `FUN_10016ea0` reads signature positions 1, 2, 3, 5,
and 6 when making its five-byte class key. This gives a tightly constrained
candidate repack: `[attr_48, key[0:3], attr_40, key[3:5]]`. It keeps every
legacy key byte in a position consumed by the builder and places the two
remaining legacy byte fields at positions 0 and 4. For `unit-etc`, local row
1286, the on-disk fields are `attr_48=02`, `key=13 42 18 00 1e`, and
`attr_40=21`; the candidate signature is therefore
`02 13 42 18 21 00 1e`. This is structurally supported by the loader and
builder, but still needs confirmation against paired-generation data or the
2006 selector's live rows.

With that repack and `attr_b=0`, the traced 2013 builder produced class ID
9147 and key `17 2b 2f 00 00`. That key matches the output of
`FUN_10016ea0` for the row's reconstructed five-byte legacy key
`13 42 18 00 1e` (lookup-remapped first three bytes, unchanged fourth byte,
and the fifth byte masked with `0x20`). The same Hello fixture retained scalar returns
`1207, 555, 1740, 4832` and both recovered 12-value vector outputs. The
selected IDs changed to `273369, 273369, 273370, 273370, 273371, 282025,
282025`; the WAV is 8,896 frames (17,836 bytes, SHA-256
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`). This
isolates selection sensitivity to the signature repack while the tree outputs
remain fixed. It does not prove that the candidate sequence is linguistically
correct.

The old `[7N:8N]` comparison window is `attr_40`, a real legacy per-unit
field, not a signature byte column. Copying it into 2013 `attr_b` while
leaving it in the signature duplicates that field. The transfer variant
moves `attr_40` into `attr_b` and zeroes its old location. Under the corrected
signature mapping, this transfer produced an 8,132-frame WAV (16,308 bytes,
SHA-256 `6c07819596129485710b732b67b8eeb314a680eafd6849509a00c6db6b902646`)
and candidate-specific `attr_b` costs. Its first selected IDs were
`273369, 273369, 272152, 272152, 272153, 282025, 282025`. This shows the
legacy `attr_40` can drive the new scoring dimension when transferred; it
does not establish whether the 2013 contract intends duplication, transfer,
or a transformed value. The zero-fill and transfer WAVs have not been judged
side by side by the requester.

A watchpoint on the same row resolved a separate apparent mismatch.
Immediately after index loading, its first byte attribute is 0; during model
initialization, code at `0x1001a358` adds 3 to that bank's whole attribute
array. The scorer later reads 3. The disassembly shows a loop applying the
bank offset, so this is an intentional normalization rather than an adapter
or memory-layout side effect. The watchpoint stack is in
`hello-index-field-origin-gdb.log`; the runner restored both Stage 5 files
byte-for-byte and reproduced the 10,650-byte control WAV.

The earlier “signature-last” label and the interpretation of `[7N:8N]` as a
unit-major signature column were wrong. The distribution comparison remains a
descriptive comparison of legacy `attr_40` with Paul's 2013 `attr_b`; an EMD
similarity alone cannot establish that these fields are semantically paired.
The corrected experiments now separate the two hypotheses: zero-fill leaves
the added score dimension uninformative, while moving legacy `attr_40` into
it restores candidate-specific costs. Exact 2013 `attr_b` semantics remain
unresolved.

### Preserve the signature while duplicating `attr_40`

The earlier corrected transfer also zeroed reconstructed signature byte `+4`,
so it changed two input locations. A paired follow-up copies legacy `attr_40`
to `attr_b` while preserving the reconstructed seven-byte signature exactly.
The overlay builder is `duplicate-legacy-attr40-to-attrb.py`; its output is
mounted by `compose-matched-context-attrb-duplicate.yaml`. Compared with the
zero-`attr_b` control, the target keys and first three ranker candidate lists
are unchanged. The first eight `FUN_100182e0` returns match the transfer run
exactly (`16.032, 30.032, 0.0705882, 0.00851064, 15.0082, 16.0348, 16.96,
15.1455`). The duplicate, move, and zero controls all produce the same
18,364-byte WAV (`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`)
under this four-position cutoff intervention and select the same final-rank
lists. Thus the change in local costs comes from the value in `attr_b`; clearing
signature byte `+4` was not required for that effect, and these score changes
do not explain the chosen-chain difference in this capture.

The 2013 class-key producer consumes signature offsets `+1,+2,+3,+5,+6`, and
the traced local scorer reads candidate signature offsets `+1,+2,+3,+5` plus
flag bits at `+5`. Neither consumes candidate signature byte `+4` on this
path. The experiment therefore establishes that this byte can be preserved
without changing the measured key/rank/audio output; it does not establish
whether another 2013 path reads it or whether legacy `attr_40` is the intended
source for `attr_b`. Captures are under
`corpus-parity/stage20/matched-context-score-duplicate/`.

### Neutralize only the `attr_b` score contribution

The scorer loads the `attr_b` pair cost into the stack slot at `[ebp+0x1c]`.
For the flagged path, the isolated hook clears that value at `0x1001864a`
before its normalization. For the ordinary path, it clears the same slot at
`0x10018736`, after the other feature distances are combined and before the
`attr_b` value is added. The lookup table, candidate data, signature, keys,
shortlists, and all other score terms remain untouched.

On matched `Hello.`, the intervention changes the first eight local scores
from `25.6, 39.6, 9.6, 9.6, 24.6, 25.6, 25.6, 24.6` to `16, 30, 0, 0, 15,
16, 16, 15`. It preserves the same three ranked candidate lists and selected
chain `273369, 273370, 273371, 255282`; the WAV remains byte-identical at
18,364 bytes (`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`).
This rules out the `attr_b` term as the cause of the remaining chain difference
for this cutoff-controlled input. It may still matter on other contexts or
when score ties and cutoffs differ. The runtime trace is under
`corpus-parity/stage20/matched-context-score-zero-attrb-neutralized/`; the
GDB script is `trace-matched-context-score-zero-attrb-neutralized.gdb`.

The tests identify one concrete tree-adapter error and a second, independent
index-layout defect on the tested `Hello.` path:

1. The old and new engines reach matching first-four scalar values. The
   generic 2013 tree parser misread the four special vector pitch trees; with
   the corrected experimental adapter, the 2006 and 2013 `nbf`/`sbf`
   evaluations take the same branches and return the same 12-value leaves
   for this input. Tree evaluation is therefore not where this controlled
   phrase currently diverges.
2. The 20-byte Kate tail is column-major (`attr_4c`, `attr_48`, five key
   columns, `attr_40`, 12 metric columns), while the 2013 signature is read
   as seven bytes per unit. The adapter must rebuild the signature, not only
   insert an `attr_b` byte. The class-builder positions constrain a likely
   mapping to `[attr_48, key[0:3], attr_40, key[3:5]]`. Repacking changes the
   selected units while preserving tree return values, directly locating a
   second selection-path defect.
4. The new independent `attr_b` has no proven 2005 counterpart. Moving
   `attr_40` to it restores candidate-specific scoring costs, but whether to
   move, duplicate, or transform that field remains open.
5. The candidate sequence still differs substantially: the old path selects
   four consecutive `unit-gen2` rows, while the 2013 copy candidate selects
   eight rows across `etc`, `gen`, and `gen2`. DAT/UPM record spans and decoded
   sample counts agree with the legacy sources. Thus a bad payload decoder or
   simple span offset is not supported by the observations; candidate choice
   and segmentation are.

The native selector also makes this boundary worth checking before blaming
the renderer. In the 2006 pseudocode, `FUN_1001e470` builds each candidate's
score from three distance-table lookups, context weights, and categorical
penalties; signature mismatches can add 100 or 1,000 points. The 2013
`FUN_100182e0` likewise combines distance lookups with categorical costs, but
its parameter view includes two separate byte attributes and the seven-byte
signature. The broad scoring shape is similar, while the precise input fields
and penalty rules differ. This supports a data-contract mismatch at the
adapter boundary; it does not show that the new scoring algorithm itself is
incompatible with Kate's speech data.

The remaining high-value probes are:

1. **Resolve `attr_b`.** Compare a voice available in both `ver.2005` and
   `ver.2013`, or find a 2013 index writer/decompiler path that defines whether
   `attr_40` is moved, duplicated, transformed, or absent in 2005 data.
2. **Compare old and new phone decisions.** The four old selected rows now
   match the rebuilt index signatures byte-for-byte, and an end-to-end 2013
   control can render those same four IDs. Semantic phone labels and equivalent
   timing boundaries are still unavailable, so this does not yet say whether
   the old and new selectors made the same phone decisions.
3. **Signature/category score audit.** For the same target contexts, compare
   all seven signature bytes and the category-table results in both engines.
   The forced P control already shows a 500-point offset-1 penalty with a
   matching central phone byte. If old-engine candidate scores do not contain
   that penalty, the key/category contract remains incompatible even after
   the `attr_b` repair.

Lower-priority alternatives are the optional `class.idx`/`classhp.idx`
fallback and the crossfade arithmetic. The available traces show synthesis
completes when those optional files are absent, and the source peaks plus
nearly complementary weights make sample overflow less likely. Neither is a
good first explanation for the large change in selected unit sequences.

## Legacy 2006 phone-record handoff

### Follow-up tests: paired data, selected-row mapping, and forced sequence

The local inventory has no same-voice Kate pair for a generation comparison.
Every Kate index is `ver.2005`; the available James indices are `ver.2013`.
Those are different voice models, so comparing their attribute values cannot
establish a cross-generation mapping.

The four IDs selected by the 2006 runtime (`272822`–`272825`) map through the
bank ranges to local `unit-gen2` rows `92827`–`92830`. The reproducible
`compare-legacy-selected-index-signatures.py` probe reads their actual legacy
columns, compares their 19-byte payload records with the 2013 overlays, and
checks the rebuilt signatures. Every rebuilt signature matches the exact
legacy fields; every un-repacked signature differs:

| Global ID | Legacy `attr_48` | Legacy five-byte key | Legacy `attr_40` | Rebuilt signature | Un-repacked signature |
| ---: | ---: | --- | ---: | --- | --- |
| 272822 | `01` | `5a 22 17 1e 1e` | `17` | `01 5a 22 17 17 1e 1e` | `11 09 10 11 0c 0d 17` |
| 272823 | `01` | `22 17 2b 00 1e` | `0a` | `01 22 17 2b 0a 00 1e` | `10 2c 0e 1c 26 21 26` |
| 272824 | `01` | `17 2b 30 00 00` | `16` | `01 17 2b 30 16 00 00` | `33 2d 10 0d 1d 14 0f` |
| 272825 | `00` | `2b 30 5a 04 04` | `5f` | `00 2b 30 5a 5f 04 04` | `18 1c 30 41 1c 0d 10` |

This is direct row-level evidence for the adapter layout error on the exact
units the old engine chose. The table and trace files are under
`tools/revkit/work/stage20/legacy-selected-signature-comparison.tsv` and the
ignored `corpus-parity/stage20/hello-legacy-records-*` captures.

The forced-sequence GDB control replaced the first four 2013 selected IDs
with `272822, 272823, 272824, 272825`. The 2013 selection path had seven
slots, so the final old ID was repeated for the three unused slots; the
timeline-builder return then limited output to its first four rows. It
returned four timeline rows and wrote a valid WAV of 4,160 frames (8,364
bytes, SHA-256
`d67dfd9c840d703421b869dd0b9c55bd0809a78b568145b6222b45d71cea9eac`). The
same 2013 scalar and vector tree results were preserved. Fixture restoration
hashes match before and after. This confirms that the 2013 renderer can
complete with the four old-selected inventory IDs and makes a four-row output;
it does not establish intelligibility or match the 2006 timing. The old WAVE
has 7,610 frames, so row timing, trimming, and joining remain different.
The forced WAV and GDB trace are under the ignored
`corpus-parity/stage20/forced-legacy-units-trim4/` directory.

The requester heard this forced output as “Hello,” with a weak initial `h`.
To test whether the first transition was arriving too early, two follow-up
renders held the same four IDs, repacked index, recovered vector trees, and
four-row trim fixed. The first timeline row's sample count was changed from
664 to either the legacy append length (1,138 samples) or the full decoded
DAT length (1,210 samples). The 2006 append trace shows 1,210 raw samples
minus a 72-sample overlap, yielding 1,138 appended samples.

| First-row count | WAV frames | SHA-256 |
| ---: | ---: | --- |
| 664 (baseline) | 4,160 | `d67dfd9c840d703421b869dd0b9c55bd0809a78b568145b6222b45d71cea9eac` |
| 1,138 | 4,634 | `153fc0465f8916f76226e840c4cf9dae482f96b09488ead6f22ffb6e59ea42a8` |
| 1,210 | 4,706 | `109ea72ad38bf0eaba4e5642fc054b43304c87c1309bf6b07f80a35baf594e56` |

The 1,138-sample output shares its first 601 PCM samples with the baseline;
the changed timeline duration affects later samples. Ten-millisecond RMS
windows show the loud rise later in the extended renders:

| Capture | 30 ms | 40 ms | 50 ms | 60 ms | 70 ms |
| --- | ---: | ---: | ---: | ---: | ---: |
| 2006 plain `Hello.` | 180 | 284 | 477 | 1,731 | 4,237 |
| 2013 forced, 664 samples | 928 | 9,107 | 9,788 | 10,557 | 10,352 |
| 2013 forced, 1,138 samples | 360 | 568 | 955 | 3,566 | 9,063 |
| 2013 forced, 1,210 samples | 360 | 568 | 955 | 3,463 | 8,290 |

This is evidence that first-row duration alone can delay the transition that
the requester heard as a weak initial consonant; the 1,138-sample control
also moves the envelope closer in timing to the old capture. It supports a
segmentation/timing cause for this part of the impression, not a recovered
phonetic label or equivalent crossfade. The new extended outputs have not yet
been acoustically judged by the requester. Reproduce them with
`trace-forced-legacy-first-row-1138-repacked.gdb` or
`trace-forced-legacy-first-row-1210-repacked.gdb` and set `EVIDENCE_DIR` to a
new directory when invoking the runner. The successful captures are under
`corpus-parity/stage20/forced-legacy-first-row-{1138,1210}-repacked/`.

Two earlier exploratory attempts used the later signature-copy index overlay
and produced eight timeline rows, so their duration overrides did not run.
Their captures are retained under `forced-legacy-first-row-1138/` and
`forced-legacy-first-row-1138-controlled/`; they are not the results in the
table above.

To test the role of the four-row trim itself, another control kept all seven
timeline rows while forcing selections 1–4 to the four legacy IDs and
selections 5–7 to repeat `272825`. It returned seven rows with counts
`664, 360, 498, 2844, 4854, 2112, 2844` and produced 13,616 frames (27,276
bytes; SHA-256
`6dbe7b3d858e6263d2dc073c43880edaacfcd876f903a32299886866a0f19349`). Its
first 1,000 PCM samples are byte-identical to the four-row control, so the
trim affects later audio and the first four forced rows suffice for the
shared opening. Rows 5–7 deliberately repeat one unit and do not model a
recovered legacy phone sequence. The requester listened and reported
“hello-oh-oh,” matching the forced repeated tail.

A hybrid control forces only rows 1–4 to `272822, 272823, 272824, 272825`
and retains the original 2013 selections for rows 5–7 (`273371, 282025,
282025`). It produced 10,604 frames (21,252 bytes; SHA-256
`6099c0d9d5e6891ac27b1b1a50c15a62348b388b5cc59ca25e2e5af02be529f5`). The
first 4,035 PCM samples match the repeated-tail control exactly; the
waveforms diverge after the fourth forced row. This isolates the ending
difference to the changed later unit choices and their timeline rows. The
hybrid WAV and trace are in
`corpus-parity/stage20/forced-legacy-first-four-natural-tail/`; the requester
reports hearing “hello-glow,” with the latter sounding like a question-like
continuation after the conversational “hello.” This is a listening
description; the retained rows do not have recovered phone names.

To localize that continuation, three otherwise identical hybrid controls
retain five, six, or all seven timeline rows. They keep the first four forced
legacy IDs and the original 2013 IDs in the retained tail rows:

| Retained rows | 2013 tail IDs included | WAV frames | SHA-256 | Capture directory |
| ---: | --- | ---: | --- | --- |
| 5 | `273371` | 5,758 | `5f85dd17309309390305d83e64647f4114af55eb62fef5f65b10f89053946131` | `natural-tail-cut5/` |
| 6 | `273371, 282025` | 7,558 | `737b9c235714cc07ba4f2e5b030bf2d0257952e1c36f5b1be3898087dd110b4a` | `natural-tail-cut6/` |
| 7 | `273371, 282025, 282025` | 10,604 | `6099c0d9d5e6891ac27b1b1a50c15a62348b388b5cc59ca25e2e5af02be529f5` | `natural-tail-full7/` |

The requester reports “le—” in the 5-row cut, a “low—” sound in the 6-row
cut, and a complete “glow” with a question-like form in the full 7-row
render. This localizes the completion of that impression to the seventh row,
the second occurrence of unit `282025`. The two `282025` timeline rows are
not identical: the sixth row has 1,886 samples, type byte 1, and spans 98/102;
the seventh has 3,148 samples, type byte 2, and spans 102/98. The documented
type-2 copy path adjusts the starting sample from the decoded unit, while
type 1 starts at its base. Thus the last row adds a different window of the
same selected unit. This is a concrete timing/window distinction; it does
not establish the unit's phone or prosody label.

Two further renders isolate the retained tail after the timeline builder:
one compacts rows 5–7 to a three-row output (`natural-tail-only/`, 6,570
frames, SHA-256
`c59dae1d95452959c524f77efb46fac244d4578000702f6b6b5482692b871360`); the
other keeps only row 7 (`natural-last-row-only/`, 3,148 frames, SHA-256
`2bb2ccb90afa877619b70901aefb3480df0a6fa6d15468c27ac98608f86a4486`).
The requester hears “low?” in the tail-only render and “oh?” in row 7 alone.
This indicates the question-like ending is present in the isolated tail audio,
without the forced legacy prefix. Row 7 carries the “oh?” impression by
itself. These are listening descriptions, not recovered phone labels.

The selected units resolve to source inventory rows: row 5's ID `273371` is
`unit-gen2` local row 93,376, with 1,724 decoded DAT samples; rows 6 and 7's ID
`282025` is `unit-etc` local row 2,900, with 4,932 decoded DAT samples. Its
timeline counts (1,886 + 3,148) exceed that decoded length by 102 samples,
matching row 7's leading span. Row 5's 1,724-sample timeline count matches
its full decoded payload. These lengths support a shared-window/overlap
reading of the repeated ID; they do not identify a phone or explain the
question-like perception. These source values were read from the original
Kate M16 index, DAT, and UPM files, which remain unchanged.

Three more controls isolate row 5, row 6, and the pair:

| Isolated rows | WAV frames | SHA-256 | Capture directory |
| --- | ---: | --- | --- |
| 5 | 1,724 | `d77ae6062cd8c23fbdce375923a534a2ade74f5c1b92c210dec54b018fbf02a6` | `natural-row5-only/` |
| 6 | 1,886 | `f63dac8ff4571285bcf9766a305ef7c3824f609d56c9d3cfb21a97e50bc55b52` | `natural-row6-only-corrected/` |
| 5–6 | 3,524 | `5b04f9f32d6bc1ac3bcb82b412a8018e04ab14c14cbdb2f10dfb1a4e63d16150` | `natural-rows5-6-only/` |

These retain the corresponding row metadata and compact the selected row(s)
to the front of the timeline before assembly. The requester hears row 5 alone
as somewhat like “uh,” row 6 alone as somewhat like “ah,” and rows 5–6
together as “low.” These are listening descriptions, not recovered phone
labels. The first row-6-only attempt used an incorrect source offset and
yielded an empty WAV; it is retained under `natural-row6-only/` and excluded
from the results above. The corrected row-6-only capture is in
`natural-row6-only-corrected/`.

The pair output is 86 frames shorter than end-to-end concatenation of the
isolated clips. Direct PCM comparison shows that the first 1,638 frames are
identical to row 5 alone and the final 1,800 frames are identical to row 6
alone starting at frame 86. Only the 86-frame join region differs. This
matches row 5's 86-sample trailing span, so the pair's “low” impression is
associated with the combined sequence and its join rather than a matching
whole-row impression in either isolated clip.

A same-length no-crossfade control is at
`corpus-parity/stage20/natural-rows5-6-no-crossfade/adapted-pah0.wav` (3,524
frames; SHA-256
`864e2ad6f2832e885cc0e4ace7d39fb668336f2f8f6abe5893f2a84db85d39ed`). It is
the first 1,638 samples of row 5 followed by all 1,886 samples of row 6. It
therefore preserves the pair output's exact prefix and suffix and changes
only the 86-frame blended region. This control tests whether the engine's
boundary blend is needed for the reported “low” impression. On re-listening,
the requester reports that the hard-splice control and engine-blended pair
sound the same; “la—” is another plausible interpretation. The 86-frame blend
is therefore not necessary for the reported percept. This favors a role for
the larger ordered sequence, while leaving the phonetic interpretation open.

To test order, a reversed no-crossfade control is at
`corpus-parity/stage20/natural-rows6-5-reversed-no-crossfade/adapted-pah0.wav`
(3,524 frames; SHA-256
`a0fef9784064e899e1db4fab0b506ae130a914f3d04db81256dbfcec4e41ad34`). It
keeps the first 1,800 samples of row 6, then all 1,724 samples of row 5. This
uses the same 86-sample cut at the join as the forward hard-splice control,
but reverses the unit order. This is a derived PCM control, not a native
engine render. The requester hears it as “ah-uh—,” which supports an
order-dependent sequence percept. This remains a listening impression, not a
phone identification.

A forward-order gap control is at
`corpus-parity/stage20/natural-rows5-6-100ms-gap/adapted-pah0.wav` (5,124
frames; SHA-256
`e4bd5c71234173e65bd8c6348f55d4a7477385457a871ec3f3258976c0cf9f56`). It
uses the same 1,638-sample row-5 prefix and all 1,886 row-6 samples as the
forward hard splice, with 1,600 zero samples (100 ms at 16 kHz) inserted
between them. This derived PCM control tests whether immediate adjacency is
needed for the “low/la” impression. The requester hears the 100-ms-gap clip
as “uh ah” and is skeptical of assigning it a word impression. The 20-ms
and 50-ms versions sound the same in perceptual content; only the audible
distance changes. The requester also notes a possible word impression for
those clips, but remains skeptical of that reading. Across 20–100 ms, the
pause separates the two “uh ah” sounds without changing the basic impression.
This supports an adjacency-dependent percept for the no-gap sequence,
without establishing a word or phonetic identity.

Shorter forward-order gap controls are now available:

| Gap | Frames | SHA-256 | Capture directory |
| ---: | ---: | --- | --- |
| 5 ms | 3,604 | `1a669d5d28daa2ef8fa6f07fae7c1f24b17ce1c60cd07f4bbd8b5dc1ecb5764f` | `natural-rows5-6-5ms-gap/` |
| 10 ms | 3,684 | `f6e19470c9eadba22328d22865cda5b1641e5a11ff98f03d2acca97a85a57172` | `natural-rows5-6-10ms-gap/` |
| 15 ms | 3,764 | `7d70a3007efa393bd443a303ee174c698d1c60eb8e071547b4cb502ac2be36d4` | `natural-rows5-6-15ms-gap/` |
| 20 ms | 3,844 | `a14d307c2099e7ccfbee229b67af67f1c450d45c858fcd83395a8c40a60a9f95` | `natural-rows5-6-20ms-gap/` |
| 50 ms | 4,324 | `ed51cadb9bf7de004c0d6c20dbffa8bc163ab025ceeeb0c0d36d6879d675337c` | `natural-rows5-6-50ms-gap/` |

All preserve the same row audio and insert silence only between the two
units. The requester reports that the fused “low/la” impression occurs at
10 ms and below, while 20–100 ms versions remain the separated “uh ah”
impression; only the audible distance changes in that separated range. This
brackets the perceptual transition between 10 and 20 ms. The 15-ms control
is available inside that bracket; its result has not been stated separately.
These controls answer only the adjacency question and do not address the
engine-generation mismatch.

The repeated-tail WAV and trace are in
`corpus-parity/stage20/forced-legacy-units-no-trim/`.

`run-hello-legacy-records.sh` uses the same `Hello.` fixture with the original
2006 runtime and adds a breakpoint at `FUN_100232d0` (`0x100232d0`). The GDB
script records the per-token table pointer, its 16-bit entry count, and the
first four 20-byte rows. The trace sees one token (`start=0`, `count=1`) with
four populated rows; the same table is visited for each of the three Stage 17
output slots. Parser consumption and the first four scalar tree returns match
the earlier old-engine capture (`Hello.`, `1207, 555, 1740, 4832`).

The captured rows are implementation records, not recovered phone names:

```text
18 21 26 04 64 48 19 00 98 02 62 02 6c 02 0a 08 27 20 24 01
84 23 26 04 75 48 19 00 22 01 68 01 3c 01 04 05 24 24 24 01
c0 24 26 04 7d 48 19 00 f2 01 f0 02 64 02 07 0b 24 23 22 01
24 27 26 04 8e 48 19 00 40 08 1c 0b b0 09 1c 18 22 33 3f 01
```

Static pseudocode for `FUN_100232d0` reads these rows at a 20-byte stride,
copies their fields into synthesis arrays, and uses companion per-token
metadata. This gives us a concrete old-engine handoff to analyze next, but it
does not yet identify each field's meaning or establish a one-to-one mapping
to the 2013 candidate unit IDs. In particular, these four rows alone do not
show whether the engines diverge during candidate generation, ranking, or
later assembly. The complete captures are `hello-legacy-records-gdb.log` and
`hello-legacy-records-winedbg.log`; the runner's restoration hashes are in
`hello-legacy-records-restore.txt`. One initial trace attempt read the
16-bit count as 32-bit and overran the displayed count; the corrected trace
uses the width visible in the pseudocode and reports four entries.

The producer trace at `FUN_10013cf0` exposes the input model indices for those
records: `272822, 272823, 272824, 272825`. It returns the four rows above in
that order, and the same sequence repeats for each Stage 17 output slot. On
the adapted 2013 `signature-last-copy` path, the first eight selected IDs are
`280411, 172786, 270092, 80917, 224948, 40384, 56352, 56352`. The bank and
payload-row comparison below confirms these IDs address the same inventory
rows in the adapter as in the legacy source. Phone-position equivalence is
still open, so the ID difference alone does not prove the cause of the
distorted listening result.

### Unit-ID to bank and waveform-row check

The Kate `dblist.idx` contains five ordered index banks whose counts sum to
283,696 units. Using that order as the global-ID partition gives these ranges:

| Global ID range | Index | Count |
| --- | --- | ---: |
| `0–179994` | `unit-gen.idx` | 179,995 |
| `179995–278127` | `unit-gen2.idx` | 98,133 |
| `278128–279124` | `unit-num.idx` | 997 |
| `279125–282790` | `unit-etc.idx` | 3,666 |
| `282791–283695` | `unit-alp.idx` | 905 |

Applying those offsets maps the four legacy IDs to consecutive local rows in
`unit-gen2.idx` (local rows 92,827–92,830). The first eight 2013 copy-path
IDs span `unit-etc`, `unit-gen`, and `unit-gen2`. The legacy runtime manager
trace confirms the initial `unit-gen` count of 179,995 and the next ranges:
`[179995, 278128)` for `unit-gen2`, `[278128, 279125)` for `unit-num`, and
`[279125, 282791)` for `unit-etc`. Thus the four old IDs really resolve into
`unit-gen2`. For every listed ID, the corresponding 19-byte payload-span
record in the adapted 2013 index is byte-identical to the legacy source row.
For these IDs, the 2006 and adapted 2013 paths therefore point to the same
underlying inventory rows, though their selected sequences and bank mixtures
differ.

| Path and ID | Bank-local row | DAT offset / bytes | UPM offset / left,right | Decoded DAT samples |
| --- | ---: | ---: | ---: | ---: |
| 2006 `272822` | `gen2:92827` | `69607704 / 620` | `1656932 / 10,8` | 1,210 |
| 2006 `272823` | `gen2:92828` | `69608324 / 316` | `1656949 / 4,5` | 578 |
| 2006 `272824` | `gen2:92829` | `69608640 / 612` | `1656957 / 7,11` | 1,180 |
| 2006 `272825` | `gen2:92830` | `69609252 / 2480` | `1656974 / 28,24` | 4,854 |
| 2013 `280411` | `etc:1286` | `2063116 / 564` | `41554 / 5,9` | 1,112 |
| 2013 `172786` | `gen:172786` | `130743160 / 696` | `3066009 / 7,8` | 1,340 |
| 2013 `270092` | `gen2:90097` | `67621948 / 424` | `1611110 / 5,4` | 768 |
| 2013 `80917` | `gen:80917` | `65333696 / 556` | `1524107 / 4,8` | 1,078 |
| 2013 `224948` | `gen2:44953` | `31282776 / 472` | `757727 / 7,4` | 934 |
| 2013 `40384` | `gen:40384` | `33643860 / 612` | `782835 / 7,8` | 1,226 |
| 2013 `56352` | `gen:56352` | `46294388 / 1028` | `1077372 / 13,11` | 2,032 |

This is a measurable old/new selection difference at the inventory-row
boundary: the traced 2006 path produces four adjacent `gen2` rows, while the
2013 copy path begins with a mixture of rows from three banks. All listed DAT
payloads decode successfully. The four old payloads total 7,822 raw samples;
the old `Hello.` WAV has 7,610 frames. A trace at `FUN_10023b40` explains the
difference exactly: its four calls append 1,138, 506, 1,112, and 4,854
samples. For the first three rows, a record flag of 2 activates subtraction
of 72, 72, and 68 samples respectively; the final row is appended whole.
The resulting 7,610 samples are written as 15,220 bytes, matching the old
WAV's data length. Thus the old engine's boundary trimming is now directly
accounted for, not just inferred from the final length. The legacy
`FUN_10023a10` audio-preparation path receives these same four rows with
sample counts `1210, 578, 1180, 4854`.

### Split-helper threshold and candidate-lookup trace

The old/new threshold values (`>2` in `FUN_1001dfb0`, normally `>9` in
`FUN_10024060`) suggested that sums from 3 through 9 could cause extra 2013
splits. A paired runtime trace tested that explanation on plain `Hello.`
using the original 2006 MSI package and the adapted 2013 index/vector overlay.
The old capture is `split-threshold-old-gdb.log`. The 2013 capture was
repeated with only split-helper, metric-sum, fallback, and position-builder
breakpoints in `trace-split-threshold-new-stable.gdb`; its output and fixture
restoration hashes are under
`corpus-parity/stage20/split-threshold-new-stable/`.

The 2006 builder returns four positions. Each whole-position candidate set is
accepted: sums are 75 (`52164:45, 52165:2, 52166:24, 52167:4`), 9
(`32401:5, 32400:4`), 5 (`25383:5`), and 6 (`37302:6`). No old fallback is
taken. In the stable 2013 capture, the builder returns seven positions. Its
whole-position attempts at slots 0, 2, and 5 return false and enter the
two-row `FUN_100242a0` fallback. Those three attempts make no call to
`FUN_10023060`; `FUN_10024060` exits as soon as its `FUN_10018770` candidate
count is zero. Slot 4 is accepted with one class ID (`9147`), weight 25, and
sum 25. Thus the higher 2013 threshold cannot explain the three observed
splits: their score threshold is never evaluated.

The control-flow distinction explains why the empty lookup happens before
thresholding. For an ordinary whole-position row, `FUN_10018770` derives
candidate context variants and dispatches `FUN_10023f90` through
`FUN_10023dc0`. That helper derives a five-byte key using `FUN_10016ea0`;
`FUN_10019450` then binary-searches the sorted class-key records. A missing
key leaves the candidate count at zero, which causes `FUN_10024060` to return
false. The fallback labels its two rows with `+0x04` values 1 and 2, taking
the relaxed/expanded `FUN_10023e70` feature-query path. That separate helper
uses the population sum and normally retries until it reaches at least 10.
The class-byte-8 special branch in `FUN_10024060` instead accepts any
positive sum.

The low-intrusion whole-key trace records stable seven-byte target signatures
for slots 0, 2, 4, and 5. The corresponding 2013 keys, after the exact
`FUN_10016ea0` transform, are `[90,34,23,160,0]`, `[34,23,43,32,0]`,
`[23,43,47,0,0]`, and `[43,48,90,69,0]`. The 2006 query keys are
`[90,34,23,30,30]`, `[34,23,43,0,30]`, `[23,43,48,0,0]`, and
`[43,48,90,4,4]`. The prefixes align through the same three lookup tables;
the suffixes follow different encoding and search rules.

The 2006 path has its own candidate recovery before its `>2` final test.
`FUN_1001da80` emits key variants; `FUN_1001d8c0` and `FUN_10011fa0` retry
weaker feature queries when the initial lookup is empty or its accumulated
population is too small. The retries are limited by the current relaxation
level and key-specific conditions. This is a broader search route than the
2013 ordinary whole-position exact-key lookup. The measured cause of the
extra positions is therefore candidate lookup coverage, not the 3–9
threshold window for the baseline's empty lookups. For the four mapped
contexts, a source-row crosswalk now verifies which class memberships and
threshold effects match across generations.

For this `Hello.` trace, the old engine's subsequent feature-distance ranking
is bypassed. `FUN_1001dfb0` copies candidate IDs directly when the deduplicated
list has fewer than 31 entries; only larger lists enter `FUN_1001d5f0`, which
can call `FUN_1001d7c0` and `FUN_1001d530`. The four old lists have 4, 2, 1,
and 1 IDs, so their seven-byte feature-record comparisons do not determine
these accepts. A pointer-correct trace read records at
`*(model+0x98)+id*7`; their field semantics and source-unit identity remain
unverified. See `trace-old-class-features-v2.gdb` and
`corpus-parity/stage20/old-class-features-v5-gdb.log`; fixture restoration
hashes are in `corpus-parity/stage20/old-class-features-v5/`.

The actual old query producer is `FUN_1001b2d0`: it maps key bytes 0–2 using
the category byte at +1 and the decimal components of byte +3, copies bytes
3–4 to the corresponding output positions, and table-maps key bytes 2 and 0
into output bytes 5 and 6. `FUN_1001b3b0` then builds two directional,
eight-byte views for relaxed comparisons. Together with `FUN_1001da80`'s
context/category variants, this confirms that old candidate recovery is a
context-conditioned feature expansion, not a suffix-only conversion.

An additional trace placed breakpoints on every `FUN_10018770` entry to print
the query keys and candidate IDs. Those software breakpoints changed later
expansion and the builder returned 15 positions, so that capture is retained
only as exploratory evidence; it is not used for the seven-position result.

The offline scanner `scan-target-key-coverage.py` applies the three class-key
tables and `FUN_10016ea0`'s transform to target signatures, then compares the
keys against the experimental repacks. A minimal target-key trace and a
phone-ID-remapped control corrected the earlier higher-intrusion result: the
stable slot-5 prefix is `2b 30 5a`, not `2b 30 42`. Current and remapped
signatures and baseline audio are identical. The prefix/category mismatch is
therefore not supported for this fixture.

The baseline repack misses slots 0, 2, and 5 only at the suffix. Slot 4 still
matches runtime class 9147 with population 25, validating the key sort and
population calculation. Source rows exist for the other contexts: prefix
`5a 22 17` has 30 rows at source tail `1e 1e`; prefix `22 17 2b` has nine
rows split between `00 1e` (5) and `00 00` (4); prefix `2b 30 5a` has six
rows at source tail `04 04`.

An isolated index conversion maps `04 04` to the target suffix `45 00`.
The stable 2013 split-helper trace finds the resulting slot-5 class (ID 19786,
weight 6), then rejects it at the strict `>9` threshold and takes the
half-position fallback. The old path reports class 37302 with metric 6 for
the corresponding query and accepts it at `>2`. The six bank-local source
rows are identical on both sides, so this directly establishes the threshold
as the slot-5 fallback cause after lookup succeeds. In the matched-context
variant, the same six rows are class 18857. The prefix-specific conversion
also maps slot 2's two source tails under `22 17 2b` to `20 00`. Its
context-prefix-only capture finds class 14109, weight 9, rejects it at `>9`,
and takes the same fallback. The matched-context variant finds class 13478;
its nine bank-local rows are exactly the union of old classes 32401 (five
rows) and 32400 (four rows). The old sum 9 passes `>2`, while that same
nine-member set fails the 2013 `>9` gate. The mapped run returns six positions
and produces the same 18,364-byte output hash as the slot-5 tail mapping. The
capture and script are
`corpus-parity/stage20/split-threshold-context-prefix-tail-map/` and
`trace-split-threshold-context-prefix-tail-map.gdb`.

The tail conversion is context-dependent. Source tail `00 00` belongs to the
slot-2 aggregate mapped to `20 00` under prefix `22 17 2b`, while the same
source tail remains `00 00` under slot-4 prefix `17 2b 2f` to yield the
25-row class. A global pair-to-pair substitution would damage the latter
class. This matches the call paths: 2006 converts the key through phone and
category tables, unpacks and normalizes both decimal suffix fields in two
directional feature views, then searches relaxed variants; 2013 makes an
exact class key from the packed signature byte before expanding bitfields for
scoring. The old search policy cannot be reproduced by rewriting the suffix
pair alone.

`report-source-context-rows.py` reproduces the source member inventory by
bank-local row. For slot 2, tails `00 00` occur at `gen:13942,36540` and
`gen2:84077,92157`; tails `00 1e` occur at `gen:53487,96599` and
`gen2:88765,92828,93375`. Slot 4's unchanged `00 00` group has 25 rows;
slot 5's `04 04` group is `gen:1869,5170` and
`gen2:66682,75287,84079,92830`.

`compare-legacy-new-shortlist-members.py` applies the 2006
`FUN_1001b2d0` table transform to every source key and compares the result
with pointer-correct seven-byte rows captured from the original runtime at
`*(model+0x98)+id*7`. For each old class ID, its matching source-row count
equals its runtime `+0x8c` weight. It then compares bank-local row sets with
the exact 2013 mapped classes:

| Context | Matched 2013 class | Old member set | 2013 member set | Shared | Difference |
| --- | ---: | ---: | ---: | ---: | --- |
| Slot 0 | 32928 | 75 | 30 | 30 | 45 old-only |
| Slot 2 | 13478 | 9 | 9 | 9 | identical sets; old 5+4 classes merge |
| Slot 4 | 9147 | 5 | 25 | 5 | 20 new-only |
| Slot 5 | 18857 | 6 | 6 | 6 | identical sets |

This proves the 2013 population cutoff causes the slot-2 and slot-5 fallback
splits on the same source units; their old sums of 9 and 6 pass `>2`, while
the new values fail `>9`. Slots 0 and 4 pass in both generations, but their
different member sets remain a downstream ranking hypothesis. The script
prints the complete shared and one-sided bank-local row lists.

A controlled GDB intervention then changed only the split-helper return rule:
nonempty 2013 results with a measured sum above 2 were accepted. The slot-2
sum 9 and slot-5 sum 6 both passed, their two half-position fallbacks
disappeared, and the position builder returned four instead of six. The
30- and 25-row classes remained accepted without intervention. This confirms
that the `>9` cutoff alone causes the two extra positions for this matched
`Hello.` overlay. The WAV hash changed to
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`, so
matching the position count does not establish matching selected units or
sound. The override selected `273369, 273370, 273371, 255282`; the native
2006 sequence is `272822, 272823, 272824, 272825`. The 2013 selections map to
`gen2` rows `93374, 93375, 93376, 75287`; all four lie in the shared old/new
source-member set for their respective contexts. Thus the old choices were
available to the 2013 ranker, and the 2013 choices were available to the old
selector. The remaining difference for these four positions is downstream
ranking or timing. The trace is
`corpus-parity/stage20/old-cutoff-experiment-v3-gdb.log`; fixture hashes are
in `corpus-parity/stage20/old-cutoff-experiment-v3/`; the intervention script
is `trace-split-threshold-matched-context-old-cutoff.gdb`.

A paired scorer probe held the matched-context suffix map, vector trees, and
`>2` cutoff intervention fixed. One index overlay moved the legacy byte at
signature offset `+4` into the 2013 `attr_b` column and zeroed that signature
byte; the control left `attr_b` zero. The 2013 key builder does not consume
signature byte `+4`, so both runs retained class populations `30, 9, 25, 6`.
Both builders returned four positions and selected
`273369, 273370, 273371, 255282`. The WAVs are byte-identical, with SHA-256
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`.

The move does change per-unit `FUN_100182e0` costs. In the first ranking pass,
unit `273369` changes from `9.6` to `0.0705882`, while legacy-selected unit
`272822` changes from `9.6` to `0.00851064`. The new local score is slightly
lower for the old engine's choice, but the final sequence and WAV do not
change. This confirms that the transferred field reaches local scoring; its
mapping alone does not explain the selected-sequence divergence on this
fixture. The paired GDB scripts are
`trace-matched-context-score-{zero,transfer}.gdb`; logs, WAVs, and restoration
hashes are under `corpus-parity/stage20/matched-context-score-{zero,transfer}/`.
The transfer overlay is reproducible with
`transfer-signature-field-to-attrb.py`.

`attr48-key3-attr40-key2-matched-context-tail-map` conditions the three
suffix remappings on their transformed prefixes and leaves every other source
row unchanged; slot 4's `00 00` suffix stays intact. Its runtime classes are
`32928:30`, `13478:9`, `9147:25`, and
`18857:6`; the 30- and 25-weight classes are accepted, while the 9- and
6-weight classes fall back. The builder still returns six positions and the
WAV hash remains `f438243151b5657a37c4996ce242960ab3d0bde59e4ecb93491a8e2e2a9ede1e`.
This result shows that broad suffix rewriting is not needed to reproduce the
observed `Hello.` split decisions.

The four stable 2013 prefixes are `90 34 23`, `34 23 43`, `23 43 47`, and
`43 48 90`; the corresponding 2006 prefixes are `90 34 23`, `34 23 43`,
`23 43 48`, and `43 48 90`. The third class-key lookup table maps `48` to
`47` in slot 4. The stable slot-5 prefix therefore already agrees with the
legacy terminal context. The earlier `43 48 42` and `2b 30 42` account came
from a higher-intrusion trace that changed later expansion and is excluded
from the stable findings. A current-versus-phone-ID-remapped control produced
identical stable target signatures and identical output SHA-256
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`.
Numeric phone ID 66 is therefore not the source of the earlier apparent
prefix mismatch in this fixture.

The remaining compatibility problem is the suffix conversion and its
interaction with candidate thresholds. The 2006 query builder splits its
last two bytes as decimal-packed components and searches relaxed variants.
The 2013 class builder copies signature byte `+5` and masks byte `+6` with
`0x20`; the scorer then interprets the copied byte as two three-bit fields
plus flags. The baseline index adapter does not translate those distinct
representations, so slot 5's exact key `2b 30 5a 45 00` is absent even though
the source inventory has six rows at the matching legacy prefix and tail
`04 04`.

An isolated candidate index conversion maps source tail `04 04` to `45 00`
while leaving the target signature untouched. The stable split-helper trace
then finds class 19786 with population 6 for slot 5, but `FUN_10024060`
rejects the sum and takes the two-row fallback. This matches the 2006
slot-5 query's class population 6, which its `>2` acceptance permits. It
demonstrates both the suffix-format incompatibility and the independent 2013
`>9` sufficiency threshold on this matched source-row set. The fixture has
six source rows for this exact key, so no threshold change to `>9` alone can
make that candidate pass.

This is a controlled mapping hypothesis, not a general conversion rule: the
new overlay changes `Hello.` output to 18,364 bytes with SHA-256
`f438243151b5657a37c4996ce242960ab3d0bde59e4ecb93491a8e2e2a9ede1e`, and
the stable position builder returns six positions. Other suffixes need
equivalent source-row and candidate-set matching before the mapping can be
generalized. Current and remapped runtime traces, the suffix-mapped run, and
fixture-restoration hashes are under `corpus-parity/stage20/`; reproducible
GDB scripts are `trace-hello-whole-target-keys.gdb`,
`trace-hello-whole-target-keys-context-tail-map.gdb`, and
`trace-split-threshold-context-tail-map.gdb`.

The matched 2013 copy-path trace reports one timeline with eight rows for
`Hello.`. That is a concrete four-versus-eight segmentation difference,
alongside the changed bank mix. It locates a synthesis-path divergence after
the 2006 parser and matching scalar tree outputs; it does not establish that
the 2006 unitization is the only correct one, or that segmentation alone
causes the audible defect.

The 2013 timeline rows have fourth-word sample counts
`438, 760, 484, 798, 654, 716, 1166, 956`, totaling 5,972. In
`FUN_1002c8b0` pseudocode, the row field at offset `0x0c` is used as a 16-bit
sample count: the function copies `count * 2` bytes from the selected segment
into scratch audio. All eight rows have first and second scale fields 100/100,
so `FUN_1002afb0` takes the direct path through `FUN_1002aac0` rather than its
time-scaling branch. Their third scale field at row offset `0x08` is 200.

The captured row type byte at offset `0x26` alternates between 1 and 2.
`FUN_1002c8b0` uses the decoded segment buffer's base for type 1; for type 2,
it shifts the copy start by
`(short(source_record + 0x0c) - row_leading_span) * 2` bytes. It then copies
the row's sample count. Thus a timeline count is the size of a selected
window, not necessarily the whole decoded DAT payload. The same function
multiplies each copied signed-16-bit sample by the third scale field and
divides by 100, saturating values to the signed-16-bit range when necessary.

The `FUN_1002aac0` pseudocode reads a leading span at row offset `0x22` and a
trailing span at `0x24`. The 2006 `FUN_10023b40` and 2013 `FUN_1002aac0`
disassemblies both retain a trailing span for the next row and advance by
`count - trailing_span` for non-final rows. Both use two coefficients from
their byte-identical 8,192-entry tables to combine overlapping samples, but
their arithmetic differs: the 2006 path sums and clamps the combined weighted
value (`0x10023c42`–`0x10023c9d`), while the 2013 path clamps each weighted
lane before a later 16-bit addition (`0x1002add8`–`0x1002ade8`).
Their table-index schedules also differ. In the 2006 disassembly
(`0x10023be9`–`0x10023c09`), for overlap index `k` in a carry span of `m`,
the old path computes
`floor((k + 1) * 8191 / (2*m + 1))` and
`floor((m + k + 1) * 8191 / (2*m + 1))`, where `m` is the pending carry span.
In the 2013 direct path, `FUN_1002aac0` (`0x1002ab6f` onward) scales loop
counters in steps of 4096 over active spans and indexes the lower and upper
halves at
`floor(k * 4096 / m)` and `4096 + floor(k * 4096 / m)`. The 2013 routine has
separate branches when the pending carry and current leading span differ, so
`m` is the active span in each loop. These formulas are derived from x86
integer operations, not Ghidra's pseudocode. The 2013 routine also adds
paired scratch-buffer samples 1,000 shorts apart into output. It derives its
interior copy length as `max(count - trailing_span - leading_span, 0)` and
advances the PCM cursor by that length plus the leading span; the final row
emits its tail too. Applying that cursor arithmetic to the stable timeline
capture gives:

| Row | Samples | Leading span | Trailing span | Cursor contribution |
| ---: | ---: | ---: | ---: | ---: |
| 0 | 438 | 78 | 88 | 350 |
| 1 | 760 | 94 | 98 | 662 |
| 2 | 484 | 100 | 98 | 386 |
| 3 | 798 | 94 | 104 | 694 |
| 4 | 654 | 110 | 92 | 562 |
| 5 | 716 | 86 | 92 | 624 |
| 6 | 1,166 | 92 | 90 | 1,076 |
| 7 | 956 | 90 | 82 | 956 |
| **Total** | **5,972** |  |  | **5,310** |

Rows 6 and 7 both reference selected unit `56352`, whose decoded DAT payload
contains 2,032 samples. Their timeline counts sum to 2,122; subtracting the
90-sample overlap (row 6 trailing span and row 7 leading span) gives exactly
2,032. This is a useful consistency check for the repeated-unit case, but it
does not prove the timeline buffers are a direct partition of the DAT payload
or explain why this unit appears twice.

The comparison used the Kate MSI `vt_eng.dll` (SHA-256
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`) and
`binary/vt_pau.dll` (SHA-256
`200945bbb2853cc56b3eaa03770eedf26066b20bfc471ef5b0213387035b37ce`); the
Stage 1 copy of `vt_pau.dll` has the same hash. The old DLL stores the table
length `8192` at `0x10061198`, with coefficients beginning at `0x1006119c`.
The 2013 DLL stores the same length at `0x1006d1ac`, with coefficients at
`0x1006d1b0`. All 8,192 coefficient words are byte-identical between the two
DLLs (SHA-256
`6a71e6a37aaa3157d6edaae839a14986da4701112887c2c0b689d9affdbff6a6`). This
establishes identical coefficient data, while the disassembly comparison
above establishes different index schedules. The values fit
`sin(pi * i / 8192)^2` across all entries with maximum absolute error
`2.04e-7`; this identifies the curve's shape, while the formula used to
generate the original table remains an inference.

The first seven trailing spans total 662 samples, so this static accounting
predicts 5,310 output frames. A fresh selection/timeline control reproduced
all eight selected IDs, the eight-row timeline, and the existing 10,650-byte
`signature-last-copy` WAV (`e86d005f2be73d6c93b9018d833cf38fb8bb3619763b17e60cf01de3b3adf1c8`),
which contains 5,303 16-bit PCM frames. The seven-frame difference remains
unexplained. It is beyond the row-count/cursor accounting unless a captured
row input differs from the values used in that stable synthesis run.
For comparison, the four old DAT rows contain 7,822 samples and the legacy
append trace accounts for its 7,610-frame output with explicit 72/72/68-sample
trims. The old and new routines share trailing-span carry and coefficient
data, but use different integer schedules to choose the paired weights. The
observed old and new tail widths and row inputs also differ. The captured row
counts and spans predict 5,310 frames for the 2013 path, seven more than the
stable WAV; the residual remains open. This establishes a low-level join
difference, but does not yet measure its contribution to the speech distortion.

The gain field doubles each copied source sample before the join. Running
`decode_payload()` from `tools/revkit/scripts/decode_dat.py` on the seven
unique selected DAT payloads gives these peak magnitudes:

| Unit ID | Bank | Decoded peak | Peak after 2x |
| ---: | --- | ---: | ---: |
| 280411 | `etc` | 12,544 | 25,088 |
| 172786 | `gen` | 4,224 | 8,448 |
| 270092 | `gen2` | 8,640 | 17,280 |
| 80917 | `gen` | 10,752 | 21,504 |
| 224948 | `gen2` | 10,624 | 21,248 |
| 40384 | `gen` | 7,872 | 15,744 |
| 56352 | `gen` | 13,056 | 26,112 |

Even the largest becomes 26,112 after the gain, below the signed 16-bit
limit. Thus this gain step does not saturate these source samples before
joining. The captured final 2013 WAV has peak magnitude 25,088 and RMS 7,030;
the 2006 control has peak 14,336 and RMS 3,945. Neither final output contains
full-scale samples. The roughly 5 dB RMS difference is descriptive, not a
controlled gain comparison, because the engines selected different units and
waveforms. The 2006 join adds its two weighted values and clamps the combined
result before conversion. In the 2013 direct path, the x86 code converts the
weighted lanes separately, then adds two 16-bit samples 1,000 shorts apart
and stores the 16-bit sum without a visible saturation check. That sum can
wrap if its inputs exceed the signed-16-bit range. Source peaks after the 2x
gain are below that range. The table's observed sine-squared shape suggests
the paired lower/upper weights nearly sum to one, which may bound the result
by the lane peaks; exact paired sums and runtime lane values have not been
checked. The unchecked 16-bit addition is therefore a concrete trace target,
not evidence that this WAV overflowed.

An attempted `FUN_1002bd90` assembly-entry trace is not a valid output control:
the breakpoint changes the result even when its command only logs the entry
and continues. That minimal probe produced 5,518 bytes; a more detailed entry
and return probe produced 16,134 bytes. Both differ from the no-breakpoint
control. Stage 5 input and output hashes matched before and after each run, so
fixture restoration succeeded, but the breakpoint's effect on this runtime
path remains unexplained. Do not interpret its post-breakpoint row state as
normal join behavior. The seven-frame output residual and the relationship
between these row spans, source UPM spans, and perceived coarticulation remain
open. A follow-up should account for those values without stopping at
`FUN_1002bd90`. The legacy boundary and segment lines are in
`hello-legacy-records-gdb.log`; the read-only DAT decoding was an offline
analysis of the indexed payloads.

## Forced-phone parser control

The 2013 Stage 10 and Stage 19 forced-phone inputs use
`<vtml_phoneme alphabet="x-cmu" ...>`. This tag is not parsed as markup by the
2006 engine: its parser emits 11 records spanning the tag text and attributes.
Changing it to standard `<phoneme alphabet="x-cmu" ...>` still yields nine
records spanning tag text and attributes. Those old-engine WAVE files are not
valid forced-phone comparisons. Preserved evidence is:

- `old-vtml-pah0-gdb.log` and `old-vtml-pah0-winedbg.log`
- `old-phoneme-tag-gdb.log` and `old-phoneme-tag-winedbg.log`
- `new-vtml-pah0-gdb.log`, `new-vtml-pah0-winedbg.log`, `new-vtml-pah0.wav`

The runners preserve restoration hashes in `original-forced-pah0-restore.txt`
and `adapted-forced-pah0-restore.txt` for the last respective run. Their
default fixture remains the Stage 19 P AH0 input; `INPUT_FIXTURE` selects the
plain-text control for the documented differential.

## Local and transition score field isolation

The matched-context `>2` cutoff control makes a four-position 2013 path, but
its context-2 candidate class still contains 25 units. For same-key units
`272824` (the native 2006 choice) and `273371` (the 2013 choice), the direct
`FUN_100182e0` costs under the transferred-`attr_40` overlay are:

| Counterfactual | Unit 272824 | Unit 273371 |
| --- | ---: | ---: |
| No swap | 15.2893 | 1.06048 |
| Swap only three 16-bit distance codes | 15.2893 | 1.06048 |
| Swap only six scalar feature bytes | 6.87265 | 9.4771 |
| Swap only `attr_b` | 9.4771 | 6.87265 |
| Swap `attr_b` and six scalar feature bytes | 1.06048 | 15.2893 |

The six scalar feature bytes alone reverse the local preference in this pair;
`attr_b` changes the size and direction of the measured gap. The three
16-bit codes do not contribute to this local scorer result in this capture.
They are read by `FUN_10018c80` for adjacent-context costs. Swapping the codes
changes transition distances and the four selected IDs, which separates the
local scorer inputs from the transition inputs.

Swapping both `attr_b` and the scalar features reverses the two direct costs,
but does not make the ranker select either row at context 2. It selects
`273369, 264072, 264073, 264074`; the full 25-unit candidate set offers
lower-cost paths through other rows. In the paired edge trace, the same-row
old/new transitions use feature cost `46.2945`, while mixed-row edges use
`12.6564` plus categorical cost `10`. The Viterbi path therefore depends on
the transition fields and other candidates as well as the local score. This
is evidence for an interaction on one controlled sentence, not proof that a
single swapped field explains 2006 versus 2013 behavior across the voice.

The direct score captures are
`matched-context-unit-cost-fields-{metric-codes,features,attrb,scoring-fields}-swap-gdb.log`;
the selected path and transition costs are in
`matched-context-scoring-fields-swap-gdb.log`, all under
`corpus-parity/stage20/`. `swap-matched-context-metric-rows.py` builds each
disposable overlay under `stage20/`. The original 2006 runtime could not yet
be scored reliably with software breakpoints: stops inside and at the return
from `FUN_1001c860` caused an access violation after the first candidate.
Hardware breakpoints on the scorer entry and return avoid patching old code
and captured the native costs for all four selected rows:

| Target row | Legacy unit | `FUN_1001c860` cost |
| ---: | ---: | ---: |
| 0 | 272822 | 0.00141844 |
| 1 | 272823 | 1.64526 |
| 2 | 272824 | 6.19391 |
| 3 | 272825 | 3.9525 |

The context-2 old target features are `35, 101, 112, 2, 0, 0, 0, 0`; the
paired 2013 target view is `35, 94, 101, 1, 204, 118, 102, 0`. The 2006
candidate signatures and target records are logged in
`legacy-unit-score-hardware-gdb.log`. The original process exited normally;
its copied `old-pah0-prose.wav` hash equals the prior old-engine control
(`eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`), and
the fixture restoration hashes match. The reproducible trace is
`trace-legacy-unit-score-hardware.gdb` under `stage20/`.

These values establish that the 2006 scorer ran to completion for the
selected rows, but its raw costs cannot be ranked against the 2013 costs as
one common scale. The feature records differ between generations, and
`FUN_1001c860` and `FUN_100182e0` use different field and category rules.

The old neighbor-ranking return trace adds a structural comparison. On each
of the three identical-output runs, `FUN_1001cff0` returns one retained
candidate at each position: `272822`, then `272823` with predecessor
`272822`, `272824` with predecessor `272823`, and `272825` with predecessor
`272824`. The 2013 trace instead passes 9, 25, and 6 candidates into its
neighbor transitions for positions 1–3, then selects one complete path by
backtracking. The old and new search therefore diverge in effective path
width as well as score details. At the context corresponding to old class
`25383`, the member crosswalk is five old rows versus 25 new rows, with 20
new-only candidates; the 2013 path considers alternatives that do not exist
in the old set. This connects the measured class expansion to the third-chain
result under field swaps. The valid link-only capture is
`legacy-neighbor-links-hardware-gdb.log`, produced by
`trace-legacy-neighbor-links-hardware.gdb`. It reads candidate IDs and
predecessor indices from the per-position record and backpointer arrays; no
transition cost values are claimed from that trace.

A second trace locates the old pool reduction. `FUN_1001cd60` computes each
candidate's contiguous coverage across the four positions, stores the span at
candidate offset `+0x10`, and sorts the candidates. Its returned lists are:

| Position | Candidates before the old coverage gate | Full-span candidate(s) |
| ---: | ---: | --- |
| 0 | 75 | `272822` (coverage 4) |
| 1 | 9 | `272823` (coverage 4) |
| 2 | 5 | `272824` (coverage 4) |
| 3 | 6 | `272825` (coverage 4) |

For each position, the next `FUN_1001cff0` branch compares the first
candidate's `+0x10` span with the total position count (`4`). When they match,
it keeps only candidates whose span covers all four positions. The returned
lists therefore become the singletons `272822`, `272823`, `272824`, and
`272825`, with predecessor links in that order. This is the direct old-engine
mechanism behind its chain on this fixture: candidates with shorter coverage
are removed before `FUN_1001c860` scores the survivors. The 2013 path has a
parallel full-span gate in `FUN_100230a0`/`FUN_10023350`, but the matched
four-position run does not activate it: every candidate span is 1, so its 30,
9, 25, and 6 candidates continue to local scoring. The specific third-position
pool also grows from five old rows to 25 new rows, of which 20 are new-only.
The marker experiment below identifies why the 2013 gate sees no full-span
chain in the adapted model.

The continuity capture is `legacy-continuity-candidates-hardware-gdb.log`,
produced by `trace-legacy-continuity-candidates-hardware.gdb`. Hardware
breakpoints kept the original code unchanged; the old process exited normally,
reproduced the prior WAVE hash, and restored the shared fixture hashes. This
fully explains the 2006 path choice for the tested `Hello.` fixture. Whether
the same full-position coverage gate dominates other utterances, or takes a
different branch when no candidate covers all positions, remains untested.

A second old-only control, `Hi.`, confirms that the span threshold follows the
actual position count: the trace reports two positions, with candidate pools
of 10 and 9 and leading candidates `272820` and `272821` at coverage 2. The
remaining candidates in those lists have coverage 1. The old engine therefore
forms the same full-span distinction when the total is 2 rather than 4. This
probe did not run the adapted 2013 index/tree path for `Hi.`, so it extends the
legacy mechanism evidence without establishing cross-generation parity on a
second input. Its capture and fixture-restoration hashes are under
`corpus-parity/stage20/hi-continuity-probe/`.

### 2013 full-span marker comparison

The 2013 selector has a counterpart to the old full-span reduction. Static
disassembly of `FUN_100230a0` shows it builds per-candidate left, right, and
total span values, storing the total at node `+0x0e`. The caller
`FUN_10023350` compares the first node's span with the total position count;
when they match, scoring is limited to nodes that span every position.

With the matched-context index and controlled `>2` cutoff intervention, both
runs had four positions and the same target class sequence. The unmodified
adapted path reported span 1 throughout, including for IDs `272822`–`272825`.
The model `+0x30` per-unit lookups for those IDs aligned with the 2013 target
classes at the four positions, so the class sequence itself was present while
the continuity traversal marker was absent.

The old `FUN_1001cd60` reads a separate byte-per-unit array at model `+0x68`
and stops traversal when the byte is zero. For the four selected rows, the
captured bytes are `1, 1, 1, 0`, matching their legacy `attr_48` values. The
2013 `FUN_100230a0` instead tests bit `0x80` in byte `+6` of each seven-byte
candidate-signature row at model `+0x64`; for the same four rows those bytes
were `0, 0, 0, 0`.

A controlled in-memory intervention set only signature byte `+6` bit `0x80`
for IDs `272822`–`272824`, after class construction and before span building.
The spans then became 4 at all four positions, only one candidate per
position reached `FUN_100182e0`, and the selected IDs became
`272822, 272823, 272824, 272825`. The WAV was 15,264 bytes with SHA-256
`215ff65cf15e37f4dcce58527511aa31f143b58a1dca865ef082252a3a83e5f3`.
This directly establishes that the marker controls span traversal and
restores the native selected chain for this fixture.

The first disposable-index attempt set bit `0x80` at
`signature_offset + 6 * unit_count + row`; that column-major offset was wrong.
The corrected overlay sets unit-major signature byte 6 at
`signature_offset + row * 7 + 6` for local `unit-gen2` rows 92827–92829 and
passes structural validation. The first corrected runtime commands omitted
the Stage 19 versioned DLL layer and timed out before selector activity. With
`compose-versioned.yaml` included, the corrected overlay reached the selector
and the runtime bytes were `80,80,80,00` for IDs 272822–272825.

A matched no-marker/marker comparison used the same Stage 7, Stage 19
versioned-engine, and matched-context index layers. Both runs had candidate
pool sizes 211, 406, 78, and 11; both returned a 7,704-byte WAV with SHA-256
`89c9c78a1326e4b7f4a5dd93f5c2fbf64795d01e89117bd4351e32d53ae47ff0`. The
expanded GDB traces found identical span distributions, no span-3 or span-4
candidate, and no occurrence of any target unit 272822–272825 in the pools:

| Context | Candidate count | Span 1 | Span 2 | Span 3 | Span 4 | Target units present |
| ---: | ---: | ---: | ---: | ---: | ---: | --- |
| 0 | 211 | 178 | 33 | 0 | 0 | none |
| 1 | 406 | 376 | 30 | 0 | 0 | none |
| 2 | 78 | 76 | 2 | 0 | 0 | none |
| 3 | 11 | 9 | 2 | 0 | 0 | none |

The marker therefore reaches the selector but cannot alter this current
selection: the legacy target rows are absent before continuity traversal. The
previous recorded in-memory run that returned the chain
`272822,272823,272824,272825` was not reproduced under this complete
versioned-engine composition and is not the result of this disk-versus-memory
comparison.

The upstream mismatch is visible in the P AH0 target signatures captured at
`FUN_10018770`. Whole-position slots 0 and 2 return
`5a 5a 35 07 5a a0 00` and `5a 35 07 5a 5a 65 00`; `FUN_10016ea0` transforms
them to exact lookup keys `5a 35 07 a0 00` and `35 07 5a 65 00`. Neither key
exists in the active matched-context repack. The legacy-selected rows
272822–272825 instead map to `5a 22 17 a0 00`, `22 17 2b 20 00`,
`17 2b 2f 00 00`, and `2b 30 5a 45 00`. A scan across 13 available repacks
found the slot-0 query key only in two context-tail variants (33 rows each);
the slot-2 key was absent from all 13, and the slot-0 class in those variants
still excludes row 272822. The candidate-class miss precedes marker traversal
and scoring on this fixture.

`FUN_10024060` takes the target signature from a tree leaf. `FUN_10018770`
derives candidate lookup variants, then `FUN_10016ea0` builds the exact
five-byte key. Since `FUN_10016ea0` masks signature byte `+6` with `0x20`,
the continuity bit `0x80` cannot repair these missing query classes. The
working cause hypothesis has therefore shifted from marker propagation to a
mismatch between the tree-produced target signature and the repacked
candidate-class signatures. The full-pool captures and restore hashes are
under `corpus-parity/stage20/continuity-bit7-baseline/`,
`continuity-bit7-unit-major-overlay/`, and
`continuity-bit7-unit-major-coverage/`; target signatures are under
`corpus-parity/stage20/pah0-target-signatures/`. The earlier invalid overlay
builder `make-continuity-bit7-overlay.py` remains as historical evidence.

### Plain Hello: compatible key prefixes, different suffixes, stricter cutoff

The original 2006 package and matched 2013 setup were traced on the same
`Hello.` bytes. The 2006 class lookup emitted these four keys and metric
sums:

| 2006 position | Old query key | 2006 class metric sum |
| ---: | --- | ---: |
| 0 | `5a 22 17 1e 1e` | 75 |
| 1 | `22 17 2b 00 1e` | 9 |
| 2 | `17 2b 30 00 00` | 5 |
| 3 | `2b 30 5a 04 04` | 6 |

The paired 2013 tree leaves are transformed by `FUN_10016ea0` into the
repacked legacy-row keys below. These keys exist in the matched-context index;
the three prefix bytes align after the engine lookup tables, with the third
row showing the observed decimal `48` to `47` mapping. The final two fields
also use a different encoding across the engines:

| 2013 position / tree row | 2013 raw seven-byte signature | 2013 transformed key | Class and weight | 2013 decision |
| ---: | --- | --- | --- | --- |
| 0 / 0 | `5a 5a 22 17 2b a0 00` | `5a 22 17 a0 00` | `32928:30` | accept |
| 1 / 1 | `5a 22 17 2b 30 20 00` | `22 17 2b 20 00` | `13478:9` | fallback |
| 3 / 2 | `22 17 2b 30 5a 00 00` | `17 2b 2f 00 00` | `9147:25` | accept |
| 4 / 3 | `17 2b 30 5a 5a 45 00` | `2b 30 5a 45 00` | `18857:6` | fallback |

The 2013 whole-position helper therefore has candidates at all four matched
contexts. Its strict `>9` check accepts weights 30 and 25, but rejects 9 and
6, expanding the four-position selection to six positions. The old `>2`
check accepts its four shortlists, including sums 9, 5, and 6. These are
generation-specific metric tables, so the numerical sums are not claimed to
be a common scale; the live 2013 experiment nevertheless proves that its
strict boundary is the direct reason for the two fallbacks after successful
key lookup.

The old query keys are in `hello-old-query-classes-gdb.log`. The paired 2013
trace reports class IDs, weights, helper decisions, six final positions, and
a 15,264-byte WAV with SHA-256
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`.
The original 2006 WAV is also 15,264 bytes but has SHA-256
`eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`.
Matching file size therefore hides different selections/audio. Both fixture
runners restored their inputs and output files. The old trace is
`trace-hello-old-query-classes.gdb`; the 2013 trace is
`trace-hello-target-signatures.gdb`. Their captures are under
`corpus-parity/stage20/hello-old-query-classes/` and
`hello-new-query-classes/`.

### Hello phone-tree provenance

`trace-hello-phone-tree-provenance.gdb` uses hardware breakpoints to record
the position record passed to `FUN_10024680` and the query rows passed to
`FUN_10024060`. On the same plain `Hello.` input, the position record words
are `1, 1, 2, 0, 0, 0, 0, 1`. Four bank-0 rows are emitted: leaves 0 and 1
use feature index 0; leaves 2 and 3 use feature index 1. Their seven-byte
signatures are byte-identical to the previous target-signature trace. This
confirms the builder's row-to-leaf handoff for this fixture, while the record
word semantics and its earlier producer remain unknown.

This comparison constrains the P AH0 key-miss interpretation. P AH0 also
reaches bank-0 leaves 0 and 1, but their signatures contain `35` and `07` and
derive keys absent from the active adapted index. Plain `Hello.` follows the
same builder path and derives four keys present in that index. The P AH0
fixture is not a valid old-engine differential because the 2006 parser treats
its VTML tag as text. We therefore have a fixture-specific tree-key/index
mismatch, not evidence that `FUN_10024680` is generally incompatible. The
legacy index-column meaning needed to represent the P AH0 target fields is
still unresolved.

The capture is `corpus-parity/stage20/hello-phone-tree-provenance-gdb.log`.
The run exited normally, reproduced the prior 2013 WAV hash, and restored the
fixture input and output hashes.

### Natural `Apple.` selector comparison

The 2006 and matched 2013 engines were run on the same plain `Apple.` input.
Old query keys contain the phone-byte sequence `05 47 07 2b`; 2013 raw tree
signatures contain that same sequence across leaves 0–3. This supports an
input-level comparison but does not prove that both embedded dictionary
payloads are byte-identical.

| Old key | 2013 slot / raw signature | Transformed key | Runtime result |
| --- | --- | --- | --- |
| slot 0 `5a 05 47 30 30` | 0 / `5a 5a 05 47 07 a0 00` | `5a 05 35 a0 00` | empty; fallback |
| slot 2 `05 47 07 00 00` | 2 / `5a 05 47 07 2b 00 00` | `04 47 07 00 00` | class 1257, weight 50; accept |
| slot 3 `47 07 2b 00 04` | 3 / `05 47 07 2b 5a 05 00` | `47 07 2b 05 00` | empty; fallback |
| slot 4 `07 2b 5a 04 04` | 5 / `47 07 2b 5a 5a 45 00` | `07 2b 5a 45 00` | empty; fallback |

The `FUN_10016ea0` tables map `47` to `35` in one key field and `05` to
`04` in another. The later-context prefixes remain unchanged. The four suffix
pairs are context-specific; this trace does not support one global suffix
rewrite. In the matched-context index, the second transformed key has 50
members and the other three have none. Three fallbacks therefore happen
before the 2013 `>9` sufficiency test; the sole populated class passes it.
The old runtime accepts its final three contexts and splits its first. This
isolates a query-key/class-membership cause for the additional 2013 splits on
this plain-text control, without assigning semantics to the transformed
fields.

Isolated index-repack interventions tested whether alternate legacy column
and suffix mappings account for the misses. With
`attr_40/key[0:5]/attr_48`, Apple slot 3 acquired one member (class 36744,
weight 1), but the `>9` rule still rejected it; the selector remained at
seven positions and slot 0 still had no class. With the broader
`attr_48/key[0:3]/attr_40/key[3:5]` repack plus the generic context-tail map,
the dedicated Apple trace found class weights 1, 50, and 158 for slots 0, 2,
and 5; slot 3 remained empty. The weight-1 slot still failed the `>9` gate,
while slots 2 and 5 passed. The selected-position count fell from 7 to 6.
This no-pitch-overlay output is 16,216 bytes (SHA-256
`a242ff700b17a4f2511786cee2a5bbd439b4dd2f744ca5988ee769c44b0178a7`). The
alternative column-order output is 14,970 bytes (SHA-256
`91122aacbd6837b967d7ba8814ec35cf18b0f2f7f2acdd7599605dbb65e78ba2`).
These runtime interventions show that index
conversion changes class membership and downstream splits; they do not
validate either mapping as the intended legacy-to-2013 representation. The
single-member result isolates a second failure mode: once a mapping creates
a class, the 2013 sufficiency cutoff can still reject it. Verified traces and
fixture-restoration records are in
`corpus-parity/stage20/apple-alt-column-order-selector-no-vector/` and
`apple-generic-tail-map-selector-no-vector/`; read-only population counts
across generated variants are reproducible with `scan-target-key-coverage.py`.

A prefix-scoped candidate then mapped source prefix `47 07 2b` and legacy
tail `00 04` to target tail `05 00`. With the original no-pitch-overlay
composition, the 2013 selector found weights 1, 50, 20, and 158 at Apple
slots 0, 2, 3, and 4. Slot 3's 20-row class passed the `>9` gate; slot 0's
one-row class still failed it. The result was 5 positions, matching the old
engine's observed 5-position return on three repeated calls. Its WAV is
13,212 bytes (SHA-256
`167a388e30f74e83e6d955b39208fb9d62acd6063a3d833c12e0e12c6fef1e23`). This
is a controlled runtime confirmation that the missing slot-3 class came from
the legacy-to-2013 tail conversion in this context. The verified trace and
restore record are in
`corpus-parity/stage20/apple-prefix-tail-map-selector-no-vector/`; the old
position-count trace is in `apple-old-position-count/`.

A source-row crosswalk now compares the old runtime feature records with the
2013 class members for three Apple contexts. The old slot-2 class resolves to
48 source rows; the matched 2013 slot-2 class has 50, with 43 shared, five
old-only, and seven new-only rows. The old slot-3 class has 20 rows and its
2013 query key has no members. The old slot-4 class has 157 rows and the
corresponding 2013 slot-5 key also has no members. The crosswalk verifies its
old row counts against the observed 2006 lookup append counts and checks that
the Stage 19 index copies are byte-identical to the 2005 source indexes. This
is direct evidence that the current repack fails to preserve these three
legacy lookup populations, including one context that still passes the 2013
cutoff. Across all candidate repacks, the generic tail mapping recovers the
old slot-0 singleton row and 157 of the 157 old slot-4 rows, with one extra
new slot-5 row. It still misses the old slot-3 class and leaves the slot-2
class at 43 shared rows. This supports `1e 1e -> a0 00` and `04 04 -> 45 00`
for these observed contexts as index-key conversions, and the new prefix map
recovers all 20 old slot-3 rows with no extra members. These results do not
establish the conversions' scope beyond these contexts. The old trace prints ID 51804 twice for
slot 0, while the source crosswalk resolves one distinct row; that duplicate
append representation is still unexplained. Reproduce the all-repack
comparison with `crosswalk-apple-class-members.py`; the old feature trace is
`trace-apple-old-class-features.gdb` and the selector intervention uses
`trace-apple-selector-decisions.gdb`. The old feature capture is
`corpus-parity/stage20/apple-plain-features-v3-gdb.log`.

The old WAV hash is
`1ccabfc08b9ddbdf8cda2bd2a1f5d79b716649dbffe7c857572cb903bdbaa119`; the
2013 WAV hash is
`48f9610abc82debe3e7ce41000fc88e082ad1a2d47602438df02ce132db1c40e`.
Their equal 13,648-byte sizes and different hashes do not establish
intelligibility. Isolated old/new captures and restore records are in
`corpus-parity/stage20/apple-plain/`; the 2013 helper trace is in
`apple-plain-selector/`. Reproduce the captures with
`apple-plain.txt`, `trace-apple-old-query-classes.gdb`,
`trace-apple-target-signatures.gdb`, and the matched Stage 20 runner scripts.

### P AH0 key and cutoff interventions

The phone-tree builder trace confirms that `FUN_10024680` writes bank 0 leaf
indexes 0 and 1 into the query rows, and `FUN_10024060` reads those leaves as
the two P AH0 target signatures. This puts the missing keys at the tree/index
boundary before continuity or scoring.

An isolated GDB intervention changed the first target key fields in memory to
the adapted key for legacy row 272822. `FUN_10024060` then found class 32928
with weight 30 and accepted the whole position. After that acceptance, the
next query used position 1. Changing its key fields to legacy row 272823's
adapted key found class 13478 with weight 9, but the unmodified 2013 helper
rejected its score and entered the fallback. Overriding only that metric return
from 9 to 10 made the helper accept position 1; the resulting selected IDs
were exactly `272822, 272823`. This confirms both the upstream key miss and
the strict 2013 cutoff as separate decision points. The old `>2` cutoff would
accept 9 if the metrics were comparable, although the two engines' score
tables are not established as a common scale.

The measured-score key-substitution WAV contains 7,888 bytes (SHA-256
`61671203adaaa5d44667d751354467a1bf2afde16918c30d55d0e2afef029d7e`);
its selected stream is `272822, 282532, 282532, 272823`. With the 9-to-10
return override, the WAV contains 3,476 bytes (SHA-256
`3d85d8842976c02644861cd5e9b19bb919307bd4388585eb942906eb36090df9`) and
selects `272822, 272823`. These runtime-only edits are decision tests, not
model repairs or evidence of improved speech. Both runs restored the Stage 5
fixture hashes. Scripts are `trace-pah0-tree-leaf-provenance.gdb`,
`trace-pah0-legacy-key-substitution-metrics.gdb`, and
`trace-pah0-key-and-threshold-boundary.gdb`; captures are in the matching
`corpus-parity/stage20/pah0-legacy-key-substitution-metrics/` and
`pah0-key-and-threshold-boundary/` directories.

The immediate P AH0 selection cause is now pinned to two stages: current tree
leaf keys miss the adapted index, and the aligned second candidate class has
weight 9, one point below the 2013 whole-position acceptance boundary. The
origin of the key disagreement still needs to be followed into the old query
producer, its context expansion, and the 2013 phone-tree inputs.

`FUN_10016ea0` constrains the separate class-key path: it builds a five-byte
class key from signature offsets `+1`, `+2`, `+3`, and `+5`, and sets the final
key byte to `signature[+6] & 0x20`. The `0x80` bit used by span traversal is
discarded from this key representation, so continuity must reach the
selector's seven-byte table through a separate path. The corrected unit-major
overlay and manifest are under
`index-adapter-key-repacked-matched-context-attrb-transfer-bit7-unit-major/`;
the trace and runner are `trace-2013-continuity-gate-bit7-unit-major.gdb` and
`run-continuity-bit7-unit-major.sh`. Fixture restore hashes are under
`corpus-parity/stage20/continuity-bit7-unit-major-overlay/`.

An offline distribution cross-check supports the proposed old-to-new marker
mapping. Kate's legacy `attr_48` is strictly binary in all five indexes; the
`unit-gen2` bank has 84,286 nonzero values among 98,133 units. Native 2013
signatures from Paul and James also contain only `00` or `80` at byte `+6`:

| Voice / bank | Units | Legacy `attr_48 != 0` or 2013 byte `+6 == 0x80` |
| --- | ---: | ---: |
| Kate 2005 `unit-gen` | 179,995 | 153,991 (85.6%) |
| Kate 2005 `unit-gen2` | 98,133 | 84,286 (85.9%) |
| Kate 2005 `unit-etc` | 3,666 | 672 (18.3%) |
| Paul 2013 `unit-gen` | 440,124 | 417,223 (94.8%) |
| Paul 2013 `unit-etc` | 115,723 | 101,007 (87.3%) |
| James 2013 `unit-gen` | 277,100 | 246,465 (88.9%) |
| James 2013 `unit-etc` | 4,392 | 1,025 (23.3%) |

These are different voices, so the corpus frequencies do not establish a
general row-wise conversion. The controlled overlay did verify that the
specific legacy `attr_48` values for rows 92827–92829 can be written to the
expected 2013 selector bytes at unit-major offsets. The `0x80` bit does not
change the derived class key because `FUN_10016ea0` masks signature byte `+6`
with `0x20`. In the tested P AH0 path, the marker reaches runtime but has no
selection effect because the corresponding legacy unit IDs are absent from
the candidate pools. Reproduce the distribution scan with
`compare-continuity-marker-distributions.py` and the query/class scan with
`compare-pah0-target-classes.py` under `stage20/`.

## Natural Hello continuity marker and cutoff controls

A separate natural-text comparison used the same plain `Hello.` fixture with
the corrected vector-pitch trees and matched-context `attr_b`-transfer index.
The no-continuity-marker control was compared with an otherwise identical
overlay that transfers the
legacy `attr_48` values for `unit-gen2` rows 92827–92829 into bit `0x80` of
signature byte `+6`; row 92830 remains unmarked. These rows correspond to
global unit IDs 272822–272825. The comparison changes only those three marker
bits.

With the native 2013 cutoff, a hardware-breakpoint-only selection trace
reported six handoffs. The baseline selected
`273369, 273370, 273370, 273371, 264074, 264074`; the marker overlay selected
`272822, 272823, 272823, 272824, 272825, 272825`. Thus the marker changes the
chosen source rows to the legacy chain, but the stricter native cutoff still
splits positions and repeats IDs. The baseline WAVE is 15,876 bytes (SHA-256
`8d6b4f63906ad339f5b9b4dde977a495d2c7dfc7a7c9a0818e2a6237d81b3dfb`); the
marker WAVE is 15,264 bytes (SHA-256
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`).

A second A/B pair kept both the no-marker or marker index overlays, but applied
the same diagnostic cutoff override in both runs: the 2013 whole-position
helper accepted nonempty metrics above 2. The helper returned 9 at slot 1 and
6 at slot 3 in each run. With that gate held fixed, the
baseline selected `273369, 273370, 273371, 255282`, while the marker overlay
selected exactly `272822, 272823, 272824, 272825`. This is direct evidence
that the missing continuity markers cause the wrong source-row chain in this
natural `Hello.` comparison, independently of the split threshold. It also
shows the threshold difference has a separate effect: the marker overlay
still creates duplicate/split handoffs under the native `>9` rule.

The cutoff-override no-marker WAVE is 18,364 bytes (SHA-256
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`); the
marker WAVE is 15,264 bytes (SHA-256
`215ff65cf15e37f4dcce58527511aa31f143b58a1dca865ef082252a3a83e5f3`).
These are controlled selector experiments, not a repaired model or a claim
of audio quality. The native-cutoff comparison uses only hardware breakpoints
at the selected-unit handoff and process-ready point. More heavily
instrumented captures in the same evidence directory changed between runs and
are excluded from this comparison.

All four captures exited with status 0. The runner's restoration records show
the Stage 5 input hash remained
`6472bf692aaf270d5f9dc40c5ecab8f826ecc92425c8bac4d1ea69bcbbddaea4` and its
output hash remained
`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` before
and after each run. Logs, WAVE files, and restore records are under
`corpus-parity/stage20/hello-continuity-bit7/{hardware-base,hardware-marker,cutoff-base,cutoff-marker}/`.
Reproduce the low-intrusion traces with
`trace-hello-selected-units-hardware.gdb`; the paired threshold intervention
uses `trace-hello-cutoff-and-selected-hardware.gdb`. The marker overlay and
its file manifest are under
`index-adapter-key-repacked-matched-context-attrb-transfer-bit7-unit-major/`.
The manifest records base SHA-256
`17fed97418fd6eef6d3f3aa9a3cb41110051a45bb9a511ff8fc150f4a24a7b7b`, marker
overlay SHA-256
`9f8de9652cde5da5e3e51614c0d0484feeb357e2b06363ce0bfe8c6768bcebe4`, and
validation of 98,133 unit rows.

## Full-bank marker conversion on natural Apple

To check whether the Hello result generalizes beyond three rows, a second
disposable overlay applied the same conversion to every row in all five Kate
indexes: set signature byte `+6` bit `0x80` when the corresponding legacy
`attr_48` byte is nonzero. It starts from the same matched-context
`attr_b`-transfer index used by the no-marker control. All 283,696 rows were
validated; 240,417 received the marker. The per-bank source and overlay hashes
are recorded in
`index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks/continuity-bit7-all-banks-manifest.txt`.
The builder is `make-continuity-bit7-all-banks-overlay.py`.

The old-engine `Apple.` producer trace reports the same five unit IDs for each
of its three synthesis slots:
`177774, 177774, 100180, 59558, 59559`. The matched 2013 no-marker run handed
off nine units:
`103871, 142230, 128863, 96551, 96551, 271296, 271296, 142230, 103871`.
With the full-bank marker overlay, it handed off thirteen:
`103871, 142230, 128863, 266884, 266884, 266885, 266885, 142230, 103871, 128862, 142231, 266883, 128864`.
The first three handoffs are unchanged, but later choices differ; none of the
new handoff IDs equals an ID in the old five-unit sequence. This second-word
test confirms the marker mapping changes candidate continuation behavior
beyond the Hello rows, but it does not restore old selection for Apple.

The 2013 no-marker WAVE is 12,354 bytes (SHA-256
`d6055769ed787cfe74b68a0376b139014418dd43d2100e5b70bb072f73736d36`); the
full-marker WAVE is 13,178 bytes (SHA-256
`4faabfd82aa36ca97bfdc4454c185dd3cadd29d41997c814f4a66c1b4a1dcd37`). Both
2013 runs exited with status 0 and restored the Stage 5 input and output
hashes. The old-engine run also exited normally and restored all recorded
fixture hashes. These WAVE differences are not an intelligibility judgment.

Because `FUN_10016ea0` masks signature byte `+6` with `0x20`, the full-bank
conversion does not change the five-byte class keys. Existing Apple traces
show the no-marker setup has an empty candidate lookup at three whole-position
contexts and a 50-member class at another. The marker therefore cannot repair
those missing keys; its observed effect is downstream in continuity
traversal. The experiments jointly distinguish two compatibility defects:
legacy `attr_48` conversion changes traversal, while the Apple tree/index key
and class-membership mismatches remain. This result does not yet separate
later ranker differences from the remaining class-member differences. The
old-ID trace is `trace-apple-old-selected-records-hardware.gdb`; the 2013
selector-only trace is `trace-apple-selected-units-hardware.gdb`. Captures
and restoration records are under
`corpus-parity/stage20/apple-continuity-bit7/{old,new-no-marker,new-full-marker}/`.

An offline crosswalk of the old selected rows against their adapted signatures
sharpens the remaining key mismatch. The 19-byte span records for old unit IDs
177774, 100180, 59558, and 59559 are byte-identical in the 2005 source and
adapted 2013 `unit-gen` index. After `FUN_10016ea0`, the adapted row keys and
natural Apple target keys are:

| Old selected ID | Adapted row key | Apple target slot | Target key | Exact match |
| ---: | --- | ---: | --- | --- |
| 177774 | `5a 05 35 1e 00` | 0 | `5a 05 35 a0 00` | No |
| 100180 | `04 47 07 00 00` | 2 | `04 47 07 00 00` | Yes |
| 59558 | `47 07 2b 00 00` | 3 | `47 07 2b 05 00` | No |
| 59559 | `07 2b 5a 04 00` | 4 | `07 2b 5a 45 00` | No |

The row records therefore map to the same source units, but three of these
four old selected rows have a different five-byte key from the corresponding
2013 query. Exact lookup cannot place those rows in those target classes; the
continuity bit cannot change that because the key builder masks it out. The
fourth row, 100180, exactly matches the 50-member slot-2 class, yet it does not
appear in the 2013 selected handoff sequence. That makes downstream unit
ranking or transition cost a concrete next cause to test for this surviving
class. This row-key crosswalk does not yet explain which 2013 feature or cost
causes the alternate choice.

### Apple matched-row transition-cost trace

The surviving `100180` class member reaches `FUN_10018c80` alongside the
selected-side candidate `128863` in Apple context 2. A conditional GDB probe
captured both current rows against the same predecessor `142230`, then read
the scorer's temporary values and selected weight row. Both use weight row 2
and divisor 1, with zero categorical penalty. The row-2 coefficients are
`10` for the raw triangular `cepdist` lookup, `2` for derived feature distance
A, and `5` for derived feature distance B; the scorer then adds the fixed
`2.0` constant. The measured inputs reconstruct the transition subtotals:

| Current unit | Duration/prosody term | Raw distance | Derived A | Derived B | Weighted subtotal | Edge total over `142230` |
| ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| `100180` | `0` | `0.429843` | `0.115108` | `0.116129` | `7.10929` | `604.226` |
| `128863` | `0.0592593` | `0` | `0.26087` | `0.116129` | `3.10238` | `600.279` |

The subtotal difference is explained exactly by the observed inputs and
weights: the raw distance adds about `4.29843` for `100180`, while its lower
derived-A value saves about `0.29152`; derived B is equal. The fixed constant
cancels in the comparison. The separate duration/prosody term is `0` for
`100180` and `0.0592593` for `128863`, making the full edge to `100180` about
`3.94765` more expensive.

This is not a per-unit ranker reversal. The old `FUN_1001c860` scorer gives
`100180` score `0` and `128863` score `0.118519` for Apple row 2; the 2013
`FUN_100182e0` trace gives the same values for the same target and candidates.
Both local rankers prefer `100180`. The old neighbor trace retains both rows
at position 2, each linked to old selected predecessor `177774`, and its
selected-record trace chooses `100180`. The 2013 context-2 candidate list has
29 of 30 IDs in common with the old position-2 list, including both rows;
its lower-intrusion path trace instead links selected `128863` to predecessor
`142230`. The old and 2013 search structures and neighbor objectives differ,
but these captures isolate the immediate 2013 reversal to its transition
cost on a shared candidate rather than class-key membership or local score.

The relevant raw metric words are also concrete. In the adapted legacy index,
unit `100180` has group-0 code `744`, unit `128863` has group-0 code `884`, and
both old predecessor `177774` and 2013 predecessor `142230` have group-2 code
`884`. The traced 2013 table distance is therefore `cepdist[744,884] =
0.429843` for `100180`, versus the diagonal `cepdist[884,884] = 0` for
`128863`. This establishes the data values that trigger the penalty, but not
whether the 2005 code columns have equivalent meaning under the 2013 transition
model.

A controlled counterfactual changed only `unit-gen` row `100180`'s group-0
word from `744` to `884` in a disposable copy of the full-bank marker overlay.
The builder verifies that exactly two bytes differ and all five versioned
indexes remain structurally valid. The runtime trace confirms that the scorer
reads `884` for both sides of the edge. The raw distance falls to zero and the
subtotal falls from `7.10929` to `2.81086`; against `128863` on the same
predecessor, the edge to `100180` becomes about `0.3515` cheaper. The
intervention therefore reverses this pairwise edge comparison.

It does not change the global Apple result in the paired low-intrusion run:
control and intervention produce the same selected handoffs
`66216, 64727, 64728, 64729, 266884, 266885, 266885, 64727, 66216, 266884,
64729` and byte-identical 14,766-byte WAVE output (SHA-256
`75da51d1f964aa02bfc05bfb57a4775e3c091f02b452c5de6ab436d16e7223f3`). The
expanded intervention trace also backtracks through `64728` at context 2,
despite the improved `100180` edge. The raw metric mismatch is therefore a
real pairwise penalty, but correcting it alone is insufficient to restore the
legacy global path. Other candidate paths and the other transition features
remain active. Reproduce the overlay with
`make-apple-row100180-metric-match-overlay.py --output-dir
/tmp/vtspeak-apple-row100180-metric-match`; the Compose override is
`compose-apple-row100180-metric-match.yaml`. Both runs exited normally and
restored their Stage 5 fixtures. Captures are under
`corpus-parity/stage20/apple-continuity-bit7/new-metric-intervention-control/`,
`new-metric-intervention-884/`, and
`new-metric-intervention-884-costs/`.

The code contains three table terms here: one raw `cepdist` matrix lookup and
two derived 256-bin feature lookups. This corrects the earlier description of
three derived lookups in `voice-engine-and-model-formats.md`. The measured raw
term is the largest source of this edge difference, but the trace does not
establish the global cause by itself.

The lower-intrusion path capture with the same marker overlay had `142230` as
the selected predecessor for both rows and recorded cumulative costs
`526.309` (`100180`) and `522.361` (`128863`). The expanded component/weight
captures reproduce the edge calculation, but the most heavily instrumented
run produced a different overall backtrace and a different WAVE hash. Treat
the component arithmetic as observed scorer behavior; do not treat the
instrumented WAVE or its full selected path as a stable behavioral result.
The weight capture is
`corpus-parity/stage20/apple-continuity-bit7/new-full-marker-path-weights-apple/`;
its run exited successfully and its Stage 5 fixture before/after hashes match.
The preceding distance-only capture used the wrong input fixture and is not
part of this Apple comparison.

### Apple candidate metric-code crosswalk

An offline comparison reads the three 16-bit metric-code columns from the
2005 source indexes and the adapted all-bank indexes, then combines them with
the old neighbor-link and 2013 transition-candidate captures. Across the union
of the two context-2 candidate lists (31 distinct IDs), all three codes match
between source and adapted indexes for every row. The adapter therefore did
not change these candidates' metric codes. The raw `cepdist.tbl` used by this
comparison is the Kate package file, read as a 1,024-entry triangular matrix.

The predecessor links differ. In the old trace, all 29 IDs shared with the
2013 list link to `177774`, whose group-2 code is `884`. The 2013 trace assigns
those same 29 rows across five predecessors, with group-2 codes `342`, `500`,
`744`, and `822`. Recomputing only the raw matrix term over the shared set
gives four zero-distance edges under the old links and ten under the 2013
links; medians are `0.407533` and `0.374380`, respectively. Individual edges
move in both directions. For example, `100180` has code `744` and
`128863` has code `884` in both indexes. Against old predecessor `177774`
their raw distances are `0.429843` and `0`; in this 2013 trace they link to
`188825` (code `744`) and `64727` (code `822`), producing `0` and `0.450601`.
These are offline lookups using the 2013 triangular table and the captured
links; the old capture supplies predecessor IDs, not observed 2006 transition
costs.

This separates a repacking defect from a path-search difference: the shared
rows' codes are preserved, while the engines' captured predecessor relations
change which table entries the 2013 dynamic path scorer evaluates. The earlier
component probe used `142230` for both example rows, so predecessor IDs are
capture-specific and should not be combined across runs. The raw term alone
does not explain the final selected path; other transition features, candidate
sets, and backtracking remain involved. Reproduce the offline crosswalk with:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/compare_apple_metric_edges.py
```

The 2006 candidate builder has a concrete count gate at this position. In the
Ghidra pseudocode for `FUN_1001dfb0`, a deduplicated pool below 31 rows is
copied directly; a pool of 31 or more is reduced to 30 by
`FUN_1001d5f0`/`FUN_1001d530`, which compares weighted signature bytes. The
captured old Apple context-2 pool has 30 rows, so it takes the copy-all arm.
That old signature-distance helper is distinct from the 2013
`FUN_10018c80` transition sum over raw and derived `cepdist` tables. For this
position, the engines therefore feed different cost structures into their
neighbor/path decisions even though the candidate rows' metric codes match.

### Apple context-2 transition-weight factorial

A paired runtime intervention changed only the weight globals read by the
context-2 call to `FUN_10018c80`. The control read raw/A/B coefficients
`10/2/5`; each intervention restored the original values when that function
returned. Across the eight combinations of leaving or zeroing each weight,
all seven partial ablations produced the control's same 11 selected handoffs
and byte-identical 14,766-byte WAVE (SHA-256
`75da51d1f964aa02bfc05bfb57a4775e3c091f02b452c5de6ab436d16e7223f3`). This
includes raw-only, either derived term alone, and all three two-term
combinations.

Only the combined raw/A/B zeroing changed the backtrace. It selected 13
handoffs:
`66216, 64727, 128863, 266884, 266884, 266885, 266885, 64727, 66216, 128862,
64728, 266883, 128864`; the WAVE was 15,416 bytes with SHA-256
`159cf7612470a23845b96986e43c48b463af83f36a8a6cf7ec3b0a7f15ba7673`.
The trace reads `10/2/5` before the call, `0/0/0` during it, and confirms
`10/2/5` restored at return. All eight runs exited normally and their restore
records match the Stage 5 fixture hashes.

The raw term alone is therefore insufficient to change Apple selection in
this controlled run. The three weighted table terms jointly affect the path,
but no individual term or two-term subset changes its selected result. This
points to the combined transition subtotal competing with other path costs;
it does not identify which full-score margin changes sign or explain the
remaining class-key mismatches. The adjacent-row continuity gate is examined
separately below. The seven intervention
scripts are `trace-apple-context2-{raw-weight-zero,derived-a-zero,
derived-b-zero,raw-a-zero,raw-b-zero,a-b-zero,table-weights-zero}.gdb`; the
control is `trace-apple-context2-table-weights-control.gdb`. Captures are
under `corpus-parity/stage20/apple-continuity-bit7/new-transition-*`.

### Apple adjacent-row continuity gate

The `FUN_10018c80` disassembly and pseudocode place the `+6 & 0x80` test on
the previous unit's signature. When current ID is previous ID plus one and
that previous-row marker is set, the scorer assigns zero to the combined
raw/A/B table subtotal. The marker is thus read from the predecessor on the
outgoing edge. The incoming edge `64728 <- 64727` in context 2 confirms this
direction: row `64727` has signature byte `+6 = 0x8a`, and the combined table
subtotal is zero even though its three measured components are nonzero.

The next edge in the stable Apple trace is `64729 <- 64728` in context 3.
With row `64728`'s `+6` marker set, the captured raw/A/B distances are
`0.648320`, `1.30502`, and `0.318471`; the duration/prosody term is `75.9138`,
but the combined table subtotal is zero. A disposable overlay then cleared
only row `64728`'s `+6 & 0x80` bit (one byte changed across the index). The
selected path changed from 11 handoffs and a 14,766-byte WAVE
(`75da51d1f964aa02bfc05bfb57a4775e3c091f02b452c5de6ab436d16e7223f3`) to 10
handoffs and a 14,724-byte WAVE
(`5352979756a9f9fdd0923f5517864e9f054fa4dac67537428d015428de0e7fc4`). Row
`64729` was absent from the treatment's scorer calls across all contexts,
including a widened capture over all its possible predecessors. The control
proves the zero-subtotal shortcut on the outgoing edge. The one-byte
intervention proves that this marker participates in later path selection;
the treatment does not directly measure the nonzero counterfactual edge cost
because that candidate is no longer scored. Both runs restored the Stage 5
input/output hashes. The paired logs and WAVs are under
`corpus-parity/stage20/apple-continuity-bit7/new-context3-edge-control/`,
`new-context3-edge-bit-clear/`, `new-context3-edge-bit-clear-allpred/`, and
`new-row64729-all-contexts-bit-clear/`.
The capture is reproducible with
`trace-apple-path-selected-edge-components.gdb` and the row-clear overlay
builder `make-apple-row64728-clear-continuity-overlay.py`.

Dumping the candidate array at `FUN_10018c80` entry resolves where that
candidate disappears. Control and treatment both have 30 context-3 rows and
share 29 IDs. The control includes `64729` and excludes `28631`; the treatment
includes `28631` and excludes `64729`. The caller `FUN_100248b0` invokes
`FUN_10023350` before the transition scorer. The builder pseudocode calls
`FUN_100230a0` for continuity metadata, then scores candidates and caps lists
above 30. In `FUN_100230a0`, the disassembly tests signature byte `+6` bit
`0x80` while checking adjacent-context candidates (`0x1002316e` and
`0x10023216`). Together, the one-byte intervention, swapped candidate arrays,
and call order locate the marker's effect in continuity-aware metadata and
shortlist construction before `FUN_10018c80` evaluates transition edges.
Immediately after `FUN_100230a0`, row `64729` has metadata fields at `+0xc`,
`+0xe`, and `+0x10` of `1,4,0` in control and `1,2,0` in treatment; row `28631`
remains `2,3,0`. `FUN_100182e0` consumes the `+0xe` field in its score
normalization. The `1039.64` value is already present on `64729` immediately
after `FUN_100230a0` in both runs; it is the input score before the local
`FUN_100182e0` pass. The marker changes the row's position in that helper:
`64729` moves from index 68 to index 2 in control, and from index 68 to index
165 after clearing bit 7. Row `28631` remains at index 32 in both runs.
A Ghidra field-use scan of `FUN_100230a0` shows a calculated term
`-(+0xe + 100 * +0x10)`. Since `64729`'s `+0x10` stays zero, its `+0xe`
change alters that term by two points; this is consistent with the observed
reordering. The relative values of the other rows account for the resulting
list positions: the first control rows group at term `-4` (`+0xe=4`), then
`-3` (`+0xe=3`), while treatment row `64729` falls to `+0xe=2` and index 165.
This validates the ordering direction for the observed rows; full comparator
and tie behavior remain to be confirmed.

A direct mediation check cleared the original marker but forced row `64729`'s
computed `+0xe` value from 2 back to 4 at the helper's metadata-write
instruction (`0x1002328a`). The row returned to index 2, the score prefix
returned to 47 rows, and `FUN_100182e0` again computed `100 + 51.8277 =
151.828`. The resulting 14,766-byte WAVE is byte-identical to control
(SHA-256 `75da51d1f964aa02bfc05bfb57a4775e3c091f02b452c5de6ab436d16e7223f3`).
This establishes `+0xe` as sufficient to mediate the candidate-order and
rescoring change for this fixture, with the original marker still cleared.
It does not establish that no other marker-dependent paths affect other
contexts or voices. The experiment is under
`new-context3-force-e4-at-write-bit-clear/`, driven by
`trace-apple-path-force-64729-e4.gdb`.

At the next score-loop boundary, control processes a 47-row prefix and includes
`64729` at index 2; the bit-cleared run processes 46 rows and `64729` is at
index 165, outside that prefix. `28631` is at index 7 in control and index 32
in treatment, inside both prefixes. Thus the apparent `+887.812` change is
not a change to `FUN_100182e0`'s categorical penalty or weighted-distance
formula: the marker-dependent ordering determines whether the row is locally
rescored. In control, `FUN_100182e0` returns `100` categorical points plus
`51.8277` weighted distance, replacing `1039.64` with `151.828`. In treatment,
the row is skipped and retains `1039.64`. The local penalty for the control
call is the same `100` observed for `28631`; the hypothetical local score for
64729 in treatment was not computed. This explains the measured score delta
and shortlist substitution without claiming that the treatment's local
categorical or distance terms changed.

The paired score/metadata captures are under
`new-context3-pre-cutoff-ranks-control/` and
`new-context3-pre-cutoff-ranks-bit-clear/`; intermediate metadata captures
are under `new-context3-pre-rank-metadata-control/` and
`new-context3-pre-rank-metadata-bit-clear/`. The score-split, score-prefix, and
before/after `FUN_100230a0` traces are under `apple-continuity-bit7/` with
`new-context3-upstream-scores-*`, `new-context3-score-prefix-*`, and
`new-context3-pre-230a0-*` names; leading ordering-term dumps are under
`new-context3-order-key-*`; the +0xe mediation check is under
`new-context3-force-e4-at-write-bit-clear/`. Completed paired runs reproduce
the lower-intrusion WAV hashes and restore the Stage 5 fixture pair.

### Candidate coverage contract in 2006 and 2013

The decompiled producers and consumers show a related design across engine
generations, with materially different field meanings and ranking weights.
In 2006 `FUN_1001cd60` stores left-match count at candidate `+0x0c`,
right-match count at `+0x0e`, and total span at `+0x10`. Its ordering key is
`-(2 * total_span + bilateral_bonus + direction_bonus)`: one point is added
for matches on both sides, and two points when the current six-byte context
record's direction flag has coverage on its corresponding side. The caller
`FUN_1001cff0` uses the total and side counts to select a side-coverage value
for local-score normalization, preferring the available side when one side is
complete and the smaller count when neither is complete.

The 2013 `FUN_100230a0` stores right-match count at `+0x0c`, total span at
`+0x0e`, and a weighted side-coverage value at `+0x10`. For each matched
neighbor, context-record mode 0 contributes weight 2; modes 1 and 2 contribute
weight 1. The helper halves the left and right weighted sums, then stores the
complete opposite-side value when one side spans to the edge, their minimum
when neither side is complete, or their sum plus one for full-position span.
Its ordering key is `-(total_span + 100 * weighted_side_coverage)`, making
that weighted coverage term dominate ordinary span differences. The following
`FUN_10023350` uses `+0x0e` as total span and uses `+0x10` together with a
normalized total-span term to scale the local distance score in
`FUN_100182e0`.

This is not evidence of a caller reading a stale 2006 layout: each generation's
producer and consumers agree internally. It is evidence that the shared
candidate-node offsets were repurposed. The 2006 high-level behavior carries
into 2013 as adjacency coverage used for both ordering and score calibration,
but 2013 makes side coverage weighted and gives it a much stronger ranking
priority. The traversal input also changed representation: 2006 reads a
separate byte-per-unit flag at model `+0x68`, while 2013 checks bit `0x80` in
signature byte `+6`.

For Apple context 3, clearing the tested marker reduced row `64729`'s 2013
total span (`+0x0e`) from 4 to 2 while leaving its weighted side term at zero.
For this candidate the `100 * +0x10` component is therefore zero in both
runs; its two-point order-key change comes from total span alone.
The corresponding `unit-gen.idx` bytes identify the source field: legacy row
`64728` has `attr_48=1`; the matched-context adapter has signature byte `+6`
`0x00`, and the all-bank marker overlay has `0x80`. Rows `64727` and `64729`
also have legacy `attr_48=1` and change from `0x0a` to `0x8a` and `0x00` to
`0x80`, respectively. Since 2013 `+0x0c` is right-match count, row `64729`'s
captured tuples `(1,4,0)` and `(1,2,0)` imply left-match counts 2 and 0 via
`total - right - 1`. Clearing only row `64728`'s marker therefore removes two
left-side candidate matches while preserving the one right-side match.
The direct `+0x0e` intervention restored its rank, score-prefix inclusion, and
control WAV. This confirms that marker-dependent span is sufficient to mediate
the observed path change. Since `+0x0e` also feeds local-score normalization,
the intervention does not isolate ranking from normalization. The decompiled
2013 functions are preserved in
`work/reports/vt-pau-2013-span-builder-consumers.txt`; the report can be
regenerated with `work/scripts/DecompileNamedFunctionsToFile.java`. The 2006
producer and caller pseudocode are in
`work/reports/vt_eng-2006-all-functions-pseudocode.c`.

For this Apple edge, the proposed `attr_48` to bit-`0x80` conversion is verified
at the source and adapted row bytes, and the runtime intervention shows the
exact two-match span change. The remaining corpus-level question is whether
adjacent class membership and these coverage changes align across other
affected Kate contexts.

The first natural positive `+0x10` example occurs at Apple positions 5 and 6:
unit `266885` has total span 4 and weighted coverage 1 at both positions, with
signature byte `+6 = 0x00`. The producer trace shows the right side reaches the
utterance boundary, so the helper takes half the left weighted sum: position 5
has two left matches and a weight sum of 2; position 6 has three left matches
and a weight sum of 3. Integer halving yields 1 in both cases. The captured
rank key is `-(4 + 100*1) = -104`, ahead of the other observed candidates,
whose spans at these positions are at most 3.

A paired runtime intervention zeroed only this row's `+0x10` field before the
producer's node sort and recomputed its key from `-104` to `-4`. The row stayed
first at both positions because span 4 was unique there. This confirms that
the weighted field contributes 100 points to the rank key, while this Apple
case does not show a change in relative shortlist order after removing it.
The all-position capture, producer sums, and zero-field intervention are in
`corpus-parity/stage20/apple-continuity-bit7/weighted-coverage-full-marker-correct-stack/`,
`weighted-coverage-producer/`, and `weighted-coverage-zero-266885/`. Their
GDB scripts are `trace-apple-weighted-coverage.gdb`,
`trace-apple-weighted-coverage-producer.gdb`, and
`trace-apple-weighted-coverage-zero-266885.gdb`.

The `Hello.` capture provides a natural order-changing case. At position 0,
unit `272822` has `(right, span, weighted) = (5,6,4)` and rank key `-406`;
unit `273369` has `(4,5,2)` and key `-205`. The producer reports five right
matches with weighted sum 6, zero left matches, and a full-span result of
`0/2 + 6/2 + 1 = 4`. Zeroing only `272822`'s `+0x10` before the sort changes
its key to `-6` and moves it behind `273369` in the pre-score order. The
synthesized WAV hash remains unchanged. A follow-up order-only intervention
restored the original metric and rank after sorting while preserving the
altered node order. Its selected six-unit chain exactly matches the same-stack
control (`272822, 272823, 272823, 272824, 272825, 272825`), and both WAV hashes
are `936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`.
Thus the field changes the pre-score order, but later scoring recovers the
same selected chain for this fixture. The paired local-cost traces identify
the decisions. In the natural control, `272822` is first and spans the
complete six-phone position, so `FUN_10023350` scores only that leading
full-span candidate; its `FUN_100182e0` local cost is `0.00121581`. In the
order-only treatment, `273369` is first with span 5, disabling the full-span
shortcut, so all 30 candidates are scored. `273369` costs `0.0156863`;
`272822` still costs `0.00121581`. At context 1, the sole candidate `272823`
points back to `272822`, using predecessor index 0 in control and index 1 in
treatment, with the same cumulative cost `2.11883`. Contexts 2–5 each contain
one candidate and retain the same predecessor chain, with cumulative costs
`4.30964`, `12.1469`, `15.4763`, and `19.7343`. The reorder broadens the
first scored list but does not change the winning predecessor; subsequent
one-candidate transitions cannot diverge.

The 2006 engine reaches the same four underlying row IDs through its separate
full-span filter: it reduces candidate pools of 75, 9, 5, and 6 to singleton
rows `272822`, `272823`, `272824`, and `272825` before local scoring. The
2013 path records six positions with duplicate handoffs, so this establishes a
shared row chain for this input, not equal selector structure or general
parity. The captures are under
`corpus-parity/stage20/hello-weighted-coverage-full-marker/`,
`hello-weighted-coverage-producer/`, `hello-weighted-coverage-zero-272822/`,
`hello-weighted-order-only-272822/`, `hello-weighted-coverage-control-selected/`,
`hello-weighted-local-costs-control/`, and
`hello-weighted-order-only-local-costs/`; the GDB scripts are
`trace-hello-weighted-coverage.gdb`,
`trace-hello-weighted-coverage-producer.gdb`,
`trace-hello-weighted-coverage-zero-272822.gdb`,
`trace-hello-weighted-order-only-272822.gdb`, and
`trace-hello-weighted-local-costs-control.gdb`,
`trace-hello-weighted-order-only-local-costs.gdb`, and
`trace-hello-selected-units-hardware.gdb`.

## Repeated-Hello predecessor and transition comparison

The order-only intervention establishes that 2013's `+0x10` field can change
the pre-score order without changing the selected chain or PCM. A second
controlled input, `Hello hello.`, was used to follow the decision through
crowded later positions and compare the 2006 and 2013 transition calculations.

For 2013 with the all-bank bit-`0x80` overlay, the intervention moves `272822`
behind `273369` at position 0 (`272822` changes from key `-205` to `-5` before
the original metric is restored after sorting). The control and intervention
preserve all 20 selected unit IDs, all 286 keyed transition rows, their
cumulative costs, and predecessor IDs; only `272823`'s predecessor array
index changes from 0 to 1. The 2006 runtime was traced at its backpointer write
for candidate `272823`: it chooses predecessor index 1, which maps through the
live previous-position list to `272822`. This corrects the earlier repeated-
Hello link-only interpretation that read a post-backtracking state table as a
per-candidate link. On this edge both engines choose `272822`; the apparent
2006 choice of `273369` was an instrumentation mapping error.

The pair-cost traces identify the scores. In 2006, edge `273369 → 272823` has
transition feature `30.8357`, previous path cost `0.0201681`, and current
local cost `1.27656`, totaling `32.1324`. Edge `272822 → 272823` has zero
transition feature, previous cost `0.00243161`, and the same local cost,
totaling `1.27899`. The 2013 all-bank overlay run gives the corresponding
edges feature values `21.9468` and `0`, previous costs `0.00882353` and
`0.000945626`, and local cost `3.29502`; totals are `25.2507` and `3.29596`.
Both implementations select the same predecessor. The complete position-1
predecessor row was captured for this current candidate in each engine. Its
next-lowest total is `21.2669` in 2006 and `18.0691` in 2013, leaving margins
of `19.9879` and `14.7731` over the `272822` edge. The predecessor candidate
sets are not identical, so these margins establish a robust winner in each
observed list, not equal search spaces.

The raw feature lookups on the penalized edge are the same three values in
both builds, `{1.04473, 1.77778, 0.305296}`, but their slots and coefficients
differ. 2006 category 1 applies weights `{10, 10, 2}` and bias 2 to that
ordering. 2013 category 1 reads the latter two values in reverse slot order
and applies effective weights `{10, 2, 5}`, with the same bias 2. The subtotal
therefore changes from `30.8357` to `21.9468`. This is a measured transition
objective difference, but it does not flip the predecessor on this edge.
Relevant pseudocode functions are 2006 `FUN_1001e470` and 2013
`FUN_10018c80`. Captures are under `corpus-parity/stage20/` in
`hello-repeat-2006-paircosts/`, `hello-repeat-2006-pairmatrix/`,
`hello-repeat-2006-backpointer/`, `hello-repeat-2006-feature-components-v2/`,
`hello-repeat-2013-paircosts/`, `hello-repeat-2013-pairmatrix/`, and
`hello-repeat-2013-feature-components/`.

The plain-package controls keep engine and overlay effects distinct. With
identical `Hello hello.` text, the original 2006 row builder emits nine
rows (`273369, 273370, 273371, 232670, 266023, 264071, 264072, 264073,
264074`) and a 26,372-byte WAV. The standard adapted 2013 run reports 16
candidate positions, all with zero weighted-coverage fields, and produces a
17,334-byte WAV. The all-bank marker-overlay 2013 run reports a different
20-unit stream and a 26,660-byte WAV. These captures establish a segmentation
and synthesis-path difference for this fixture, but the overlay result cannot
be attributed to native engine changes alone.

The late 2013 `273370` decision changes under the full 2006 transition map:
baseline chooses predecessor `181914` at total `888.132` over `76688` at
`888.168`; the remap chooses `76688` at `888.396` over `181914` at `889.611`.
The first predecessor edge and position count remain unchanged. In 2006,
current `273370` is evaluated at context 1 against a different 30-row pool;
`273369` wins at `2.3531`, the next total is `16.809` through `182880`, and
neither 2013 alternative occurs in the old pool. This confirms a causal late
2013 scoring change while showing that the specific old/new edge is not
directly comparable.

The matched `Hi.` trace exposes a separate exact-key miss. The 2006 engine
accepts query keys `90 34 17 30 30` and `34 17 90 04 34`, selecting
`272820 → 272821`. Corrected-repack 2013 produces target-key prefixes that
match those rows after the three lookup tables, but suffixes differ:
`272820` is `5a 22 01 1e 00` versus target `5a 22 01 a0 00`; `272821` is
`22 11 5a 04 20` versus target `22 11 5a 65 00`. A two-entry, context-specific
suffix map (`1e 1e → a0 00` and `04 22 → 65 00`) makes the matching exact
lookups return one class each while the other four variants still return
zero. It restores both rows to candidate lists and changes the selected
handoffs to `272820, 272821, 272821`. Output is 14,560
bytes with SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`; fixture
restoration hashes match. The treatment proves sufficiency for this input,
not a general conversion law.

### Hi matched-row transition trace

An offline payload check confirms that the common IDs in the old and adapted
indexes are the same 19-byte source records: `272820` (SHA-256 prefix
`915fb2ae`), `65331` (`9280ba13`), and `272821` (`cd8c71d9`). This is direct
source-row identity, beyond matching numeric IDs. Recheck all four shared
rows with `crosscheck-hi-source-records.py`.

The old nine-row position-1 class contains two feature records: four rows
(`272821`, `279498`, `279499`, `280449`) share
`22 11 5a 04 22 00 22`, while five (`4087`, `4102`, `84094`, `177718`,
`243540`) share `22 11 5a 04 04 00 22`. Three isolated repacks mapped the
five-row class's final pair `04 04` to `5d 00`, `55 00`, and `4d 00` in
separate runs. Each made its matching 2013 exact lookup return one class, yet
all retained candidate counts `67, 19, 19`, selected
`272820, 272821, 272821`, and emitted the same 14,560-byte WAV hash
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.
This experiment cannot identify which of the three generated keys is the
intended conversion for the five-row class; only the field mapping is
established. Run outputs are under
`corpus-parity/stage20/hi-relax-0404-to-{5d00,5500,4d00}/`.

The old span-filter bypass admits `65331` with its own coverage-derived
weight. At position 0, 2006 scores `272820` at `0` and `65331` at `0.0285714`
and retains the old-selected `272820` row. In the matched 2013 suffix-map
capture, both rows reach the predecessor set with equal accumulated path cost
`5.4`. For current row `272821`, the transition term is `2` from `272820` and
`10.1671` from `65331`; with the same local cost `19.309`, totals are
`26.709` and `34.8761`. The measured adjacent-context term therefore favors
the same old row by `8.1671` in this pair. At the following context, repeating
`272821` adds zero transition term and zero penalty; its path total is
`54.709`, and the selected stream is `272820, 272821, 272821`.

This explains the matched first edge inside the suffix-map intervention:
neither the 2006 span gate nor an assumed equality of local scores is enough;
the 2013 transition minimum directly favors `272820 → 272821`. The 2006 path
has two positions, while this 2013 overlay has three. Earlier unadapted
four-position captures belong to a separate repack state and are not mixed
with these transition costs. The other four generated 2013 query variants
remain empty, so this is not a general repair. Reproduce with
`trace-2013-hi-pair-transition.gdb`; the control capture uses
`trace-2013-hi-pair-control.gdb`. Logs and the matched 14,560-byte WAV are in
`corpus-parity/stage20/hi-2013-pair-transition/` and
`hi-2013-pair-control/`; both runs report identical before/after fixture
hashes and the same WAV SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.

### Query producer trace and remaining work

The 2006 `Hi.` runtime trace at `FUN_1001da00` captures three unique key
operands from `FUN_1001da80`: `90 34 17 30 30`, `34 17 90 04 34`, and
`34 17 90 04 04`. The first two are the accepted split queries; all three
repeat across the runner's three render passes. See
`corpus-parity/stage20/legacy-expansion-hi/original-forced-pah0-winedbg.log`
and `trace-legacy-query-expansion.gdb`.

The paired 2013 capture at `FUN_10018770` records its source signature and
six-byte context row for each slot. Its key-class lookups emit suffixes
`a0 00`, `98 00`, `65 00`, `5d 00`, `55 00`, and `4d 00`. In the pseudocode,
these are composed into signature byte 5 from relaxation-vector entries and
flags; `FUN_10016ea0` copies that byte into key byte 3 and masks signature byte
6 with `0x20` into key byte 4. The legacy side uses its own decimal-pair key
expansion. This explains the exact `Hi.` misses: the converted candidate rows
retain legacy suffix values while the 2013 query producer emits packed
relaxation codes. The two-entry map recovers the two target rows, but is not a
general conversion rule. See
`corpus-parity/stage20/hi-2013-producer/adapted-forced-pah0-winedbg.log`,
`trace-2013-query-producer.gdb`, and
`corpus-parity/stage20/hi-2013-keyclass-lookup/adapted-keyclass-lookup-gdb.log`.

The same capture on `Hello.` broadens the comparison. The old engine accepts
four keys with counts 75, 9, 5, and 6. Corrected-repack 2013 emits ten exact
lookups as the contexts expand; nine return zero, and only
`17 2b 2f 00 00` returns one class. Its source signature groups carry byte-5
codes `a0`, `20`, `00`, and `69`, while the four old query suffix pairs are
`30 30`, `00 30`, `00 00`, and `04 04`. The 2013 trace includes both
source signatures/context rows and exact derived keys:
`corpus-parity/stage20/hello-2013-producer/adapted-forced-pah0-winedbg.log`,
`corpus-parity/stage20/hello-2013-producer/adapted-query-keyclasses-gdb.log`.
This confirms sparse 2013 exact lookups on a second matched input. Since the
transformed prefixes differ too, the observation does not isolate suffix
conversion as the only cause; the field-based repack must explain the prefix
and suffix fields together.

On `Hello hello.`, the same 2013 trace records 20 producer entries and 24
exact lookups. Five return a class and 19 return zero. The sequence includes
context-specific code families `a0/98`, `20/18/10/08`, and `45/4d`; the
repeated key `17 2b 2f 00 00` returns one class in each word instance. The
capture is `corpus-parity/stage20/hello-repeat-2013-producer/` and the
reproduction script is `trace-2013-repeat-producer-keyclasses.gdb`. The runner
exited normally and restored the fixture hashes. The existing 2006
selected-chain capture has nine output rows for this fixture. This links
sparse 2013 exact lookup coverage to the observed repeated-Hello segmentation
difference, but does not yet identify the general repack conversion.

The read-only class-member crosswalk explains why a suffix map alone cannot
repair the plain Hello candidate sets. Slot 0 has 75 old rows under legacy
feature prefix `5a 22 04` and tail `1e 1e`; the 2013 target asks for prefix
`5a 22 17` and tail `a0 00`. Thirty shared source rows map to the 2013 prefix
from raw prefixes `5a 22 18`, `5b 22 18`, `5a 22 1e`, or `5a 22 17`. The 45
old-only rows have raw prefixes `5a 22 05` or `5b 22 05`, which the 2013
third-byte table maps to `5a 22 04` instead. Slot 4 has five shared rows and
20 new-only rows: raw prefixes `18 2b 2f` and `19 2b 2f` collapse to the 2013
prefix `17 2b 2f`. Their old feature tails (`2f 18`/`2f 19`) differ from the
legacy class tail (`30 17`), but the five-byte 2013 lookup key does not retain
those two legacy feature bytes. Slots 2 and 5 have exact member parity. The
reproducible member analysis is
`analyze-key-class-prefix-crosswalk.py`.

The selected Hi rows now have a byte-level source identity check, and their
shared transition is traced below. A full class-member crosswalk across all
Hi query variants is still needed. Then test whether 2013's other signature
fields can preserve the two legacy distinctions lost in Hello slot 4, and
test field-based corrections for slot 0's third-byte mapping and the Hi
suffix misses without utterance-specific overrides. Once exact class
membership is recovered, rerun the `>2` versus `>9` comparison and trace
post-prune paths, then compare selected streams and waveform construction.
Finish with a held-out 2005 voice package. Completion requires producer,
member-set, and selector evidence that accounts for the extra positions and
wrong row choices across the matched inputs.

Earlier query and intervention evidence remains under
`corpus-parity/stage20/hi-2006-query-classes/`,
`hi-2013-keyclass-lookup/`, `hi-2013-suffix-map/`, and
`hello-repeat-2006-late-edge/`; paired 2013 late-edge evidence is under
`hello-repeat-2013-repacked-late-edge-control/` and
`hello-repeat-2013-repacked-late-edge-full-legacy-map/`. The disposable Hi
index treatment is reproduced with
`python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant
attr48-key3-attr40-key2-hi-tail-map`.

## Hello legacy-feature crosswalk and cutoff control

The diagnostic overlay `index-adapter-key-repacked-legacy-feature-class-crosswalk/`
maps measured old seven-byte feature classes onto the four primary 2013
`Hello.` keys. The four key populations are 75, 9, 5, and 6. Runtime traces
show these sums entering `FUN_10024060`: native 2013 accepts 75 and rejects
9, 5, and 6 at its `>9` check, taking three two-position fallbacks for seven
positions total. A paired GDB intervention that accepts nonempty sums above
2 takes no fallbacks and returns four positions. The native and override WAVs
are both 15,004 bytes, with hashes `d4878c9814ef2b9fe9186a2c1c6f1e7128519ed0f222bd5c3ef6ddeaa791fab0`
and `6be52f6b953a43d286ff8e23d0d707b36f60cf9e024686cf94e33f71108d4d35`.
Stage 5 restoration hashes match before and after both runs, but the source
text fixture was not hashed. These WAVs cannot be attributed to the explicit
`Hello.` control below.

These unprovenanced A/B captures show the strict cutoff changing the position
count after position-level populations are restored. The overlay combines
multiple old class groups under one 2013 key and leaves other generated
relaxation keys empty, so its four-position chain
(`52542, 264072, 264073, 264074`) cannot be attributed to the explicit
`Hello.` control. The next conversion must preserve class partitions across
the context-dependent 2013 key variants, then compare candidate rows and
scores per key. Reproduce with
`repack-index-legacy-feature-class-crosswalk.py`,
`trace-2013-hello-legacy-feature-crosswalk-native-cutoff.gdb`, and
`trace-2013-hello-legacy-feature-crosswalk-old-cutoff.gdb`; full findings are
in the Lead 5 comparison document.

An explicit replay uses the same `hello-plain.txt` bytes for both engines. The
2006 chain is `272822`–`272825`; with the crosswalk and `>2` intervention,
2013 selects `52542, 273370, 273371, 255282` and emits a reproducible 18,956-
byte WAV. In the 75-row first-position pool, old row `272822` and row `52542`
both score `9.6`, but only `52542` survives the ranker's 30-row reduction.
Positions 1 and 2 favor the selected 2013 rows locally; position 3 favors the
old row locally, while the cumulative path still favors the selected chain.
Swapping only the six scalar feature bytes in three old/new pairs exchanges
their direct 2013 local costs and changes the chain to
`179388, 272823, 272824, 272825`. This isolates the scalar fields' effect on
pairwise local ranking, but it does not recover the exact 2006 chain.

The explicit score-input, transition, and feature-swap traces are
`trace-2013-hello-crosswalk-score-fields.gdb`,
`trace-2013-hello-crosswalk-transition-path.gdb`, and
`trace-2013-hello-crosswalk-feature-swap.gdb`. Captures are in
`hello-plain-crosswalk-old-scores/`, `hello-plain-crosswalk-score-fields/`,
`hello-plain-crosswalk-transition-path-v3/`, and
`hello-plain-crosswalk-feature-swap/`. The field-by-field findings and
remaining class-partition work are documented in the Lead 5 comparison.

### Slot-2 class split control

`repack-index-legacy-feature-class-slot2-split.py` keeps the measured old
slot-2 weight-5 and weight-4 feature groups under separate 2013 keys:
`22 17 2b 20 00` and `22 17 2b 18 00`. The weight-4 assignment is
experimental. The matched `Hello.` trace shows the runtime querying both keys
and aggregating their two class IDs into one metric sum of 9. The native `>9`
check rejects it and expands that position; the paired `>2` intervention
accepts it. With the intervention, candidate counts, selected chain, and WAV
hash match the earlier merged-key capture exactly. This shows that the
key-level partition alone does not change downstream ranking for this
context, while confirming the cutoff effect on a populated nine-member set.

Run with `compose-key-repacked-legacy-feature-class-slot2-split.yaml` and the
exact fixture using
`trace-2013-hello-slot2-split-native-cutoff.gdb` or
`trace-2013-hello-slot2-split-candidate-path.gdb`. Native and override
captures are in `corpus-parity/stage20/hello-slot2-split-native/` and
`hello-slot2-split-override/`. The detailed comparison and remaining
multi-input score analysis are in the Lead 5 report.

### 2006 span-filter bypass and path-cost replay

The earlier unit-ID-only counterfactual is superseded because it retained the
old candidate's score weight after changing its ID. The corrected intervention
bypasses the full-span reduction immediately before scoring and leaves every
candidate's own coverage-derived weight intact. The old engine still selects
`272822, 272823, 272824, 272825` and produces its baseline WAV hash
`eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`.

Candidate-local 2006 costs (coverage / scorer weight / returned cost) are:

| Position | Old selected row | Partial alternative |
| ---: | --- | --- |
| 0 | `272822`: 4 / `0.166667` / `0.00141844` | `52542`: 1 / `1` / `14.5684`; `273369`: 3 / `0.285714` / `0.0201681` |
| 1 | `272823`: 4 / `0.166667` / `1.64526` | `273370`: 3 / `0.4` / `2.3784` |
| 2 | `272824`: 4 / `0.166667` / `6.19391` | `273371`: 3 / `0.666667` / `4.61288` |
| 3 | `272825`: 4 / `0.166667` / `3.9525` | `255282`: 1 / `1` / `6.04246` |

For the direct 2013-selected row chain, the old scorer favors its selected
rows at positions 0, 1, and 3, but locally favors `273371` at position 2.
The bypassed path's strongest partial-chain competitor starts with `273369`,
not `52542`; its local cost is also above old row `272822`. The path trace gives cumulative totals through
positions 1–3 of `1.64668`, `7.84059`, and `11.7931` for the old chain, versus
`2.39857`, `7.01144`, and `19.9869` for
`273369 -> 273370 -> 273371 -> 255282`; the latter's final edge adds `6.93303`.
The old complete path therefore remains cheaper after partial candidates are
admitted. The span gate explains their normal exclusion, but not the whole
row-choice difference. Trace scripts are
`trace-legacy-hello-bypass-full-span-filter-v3.gdb` and `-v4.gdb`; evidence is
in `corpus-parity/stage20/hello-plain-crosswalk-old-span-filter-bypass-v3/`
and `-v4/`. See the Lead 5 report for the correction history and the
within-engine limits on comparing cost scales.

### Per-field Hello score swaps

The six 2013 scalar feature arrays were swapped independently between the
same three old/new candidate pairs during `FUN_100182e0`. Only arrays at model
offsets `+0x54` and `+0x5c` changed the measured local scores in these target
contexts. The `+0x54` swap selected
`52542, 264072, 264073, 264074`; the `+0x5c` swap selected
`179388, 272823, 272824, 272825`, reproducing the earlier all-six swap
result. Swaps at `+0x48`, `+0x4c`, `+0x50`, and `+0x58` left the scores and
WAV unchanged. The field meanings are not recovered, and the result is limited
to these three Hello contexts. The six scripts are
`trace-2013-hello-crosswalk-feature-field-{0..5}-swap.gdb`; captures are in
`corpus-parity/stage20/hello-plain-crosswalk-feature-field-{0..5}/`. See the
Lead 5 comparison for score values, offsets, and limitations.

### Hi split threshold and decoded-sample gain

With the `5d 00` suffix overlay, the second Hi query contributes both old
classes to the 2013 metric: `13468:4` and `13467:5`, sum 9. Native `>9`
rejects the set and returns two fallback positions, producing three total.
The sum-9 override accepts it and returns two total positions; the selected
chain is `272820, 272821`, matching the 2006 two-row chain. The native overlay
selects `272820, 272821, 272821`. This isolates the additional Hi position to
the strict cutoff once the five-row class is present.

Native 2013 output (7,258 frames, SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`) and
the accepted sum-9 output (7,258 frames, SHA-256
`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`) differ
at 67 samples only, with maximum difference one PCM count. Repeating the
sum-9 capture produced the same hash.

For the accepted two-position chain, hardware-breakpoint inspection at
`FUN_1002c8b0` found `param2+8 = 200` on two decoded-row calls containing
1,420 and 5,910 PCM samples. The pseudocode multiplies each decoded sample by
this field and divides by 100. Changing both fields to 100 left the two
selected IDs and duration unchanged, and moved the output level from about
2.0x to about 1.0x the 2006 capture. The gain-adjusted output is not byte
identical: 283 samples differ from 2006, with maximum absolute difference
508. These differences cluster at samples 0–79, 1348–1419, and 7125–7257,
plus one sample at 7159. The 72-sample cluster at 1348–1419 is the end of the
first 1,420-sample decoded row, which meets a 72-sample overlap with the next
row. PCM gain explains the measured level mismatch. The paired PCM-window
capture documented below isolates the edge and join differences for this
Hi. control.

A paired hardware-breakpoint trace in the 2006 run records the corresponding
two rows at scale 100, pitch scale 100, and duration scale 100 for all three
render passes. Its row counts are also 1,420 and 5,910, exactly matching the
2013 calls. Ghidra pseudocode shows that 2006 `FUN_100232d0` copies its sample
scale from state `param5+0x3838`, while 2013 `FUN_1002c220` copies it from
`param5+0x4cf8`. Setter-entry traces identify the live state difference: the
2013 app calls `VT_SetPitchSpeedVolumePause_ENG` with
pitch/speed/volume/pause/speaker `100/100/200/925/1`; the paired 2006 probe
makes no corresponding setter call and keeps effective row scales at 100.
Changing the 2013 API volume argument to 100 sets every decoded row to scale
100. The three-position output retains 7,258 frames and has SHA-256
`f904e6c2718fc1be99c5a1b7ca5b2eace4df3223dab4800cae2791bb466693ef`. In the
two-position intervention, the API-level change is byte-identical to the
lower-level 200-to-100 row edit. The 2x level difference is therefore caused
by these callers supplying different volume state, independently of the extra
position from the 2013 selector cutoff. The instrumented 2006 WAVE hash
matches the prior control exactly.

Trace scripts are `trace-2013-hi-query-aggregation.gdb`,
`trace-2013-hi-query-cutoff9-override.gdb`, and
`trace-2013-hi-cutoff9-gain100.gdb`, `trace-legacy-hi-row-gain.gdb`,
`trace-legacy-hi-volume-setter.gdb`, `trace-2013-hi-volume-setter.gdb`, and
`trace-2013-hi-volume100-rows.gdb`, `trace-2013-hi-cutoff9-volume100.gdb`.
Captures are in
`corpus-parity/stage20/hi-query-aggregation-0404-to-5d00/`,
`hi-query-cutoff9-override/`, `hi-query-cutoff9-override-repeat/`, and
`hi-query-cutoff9-gain100-v2/`; the old capture is in
`hi-2006-row-gain/`. Setter and all-row captures are in
`hi-2006-volume-setter/`, `hi-2013-volume-setter/`, and `hi-volume100-rows/`.
The combined two-position API-level capture is in
`hi-query-cutoff9-volume100/` and matches the row-field override hash. All
runners restored their fixture hashes. Full measurements, including the
first invalid-offset probe, are in the Lead 5 comparison document.


### Matched four-row Hello PCM comparison

The same input-window capture was repeated for natural `Hello.`. The 2013
run used the measured continuity-marker overlay, accepted the whole-position
metric above 2, and changed the volume setter to 100. It selected
`272822, 272823, 272824, 272825`, which the 2006 engine selects natively.
All four decoded PCM buffers match byte for byte across DLL generations:

| Row | Samples | Leading span (2013) | Trailing span | PCM SHA-256 |
| ---: | ---: | ---: | ---: | --- |
| 0 | 1,210 | 78 | 72 | `8612a32fc5c91453afe09b9f38c5a1a2a996b4d53629d67de4304f927e6e359b` |
| 1 | 578 | 72 | 72 | `7d019679912d21cd2e847f130d6645ee187ea9199fa322b702ead957512ea7df` |
| 2 | 1,180 | 72 | 68 | `6cfc46ad8678768263c296316d8980fac9eafa9401ceb8896809eb84069d5d39` |
| 3 | 4,854 | 68 | 126 | `1bb1f78c80a60368045d258f73f5a604be5881d7fe9a4f6595265e9c4bee17e9` |

Both WAVs have 7,610 frames. The old output matches the existing baseline
hash `eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`;
the 2013 matched-chain, volume-100 output hashes to
`a8ec8d359acf9a2475c3698bc6b5e7ae4d6b41397cd22e842f9d2300c7564fba`. There
are 405 differing samples, confined to the 78-sample onset, joins of 72, 72,
and 68 samples, and the 126-sample final envelope. Maximum absolute
difference is 1,375 PCM counts; the intervening row interiors are identical
copies.

Sample-by-sample replay confirms that each new onset, overlap, and final-tail
equation matches the 2013 WAV, and each old overlap equation matches the 2006
WAV. This isolates three effects in the comparison: native 2013 selection can
choose a different chain and add positions; the paired applications pass
different volume state; and the generation-specific synthesis routines still
alter edge and overlap samples with the selected rows, PCM inputs, gain, and
duration held equal. The third effect touches 405 samples in this matched
four-row control; it does not create extra unit content by itself.

The old and new traces are `trace-legacy-hello-window-inputs.gdb` and
`trace-2013-hello-window-inputs.gdb`. Captured PCM, WAVs, runtime logs, and
fixture restoration records are in
`corpus-parity/stage20/hello-window-inputs-2006/` and
`corpus-parity/stage20/hello-window-inputs/`.

### Native cutoff split on the matched Hello rows

The same continuity-marker overlay and volume-100 setter were used without
the cutoff override. Native selection returned six IDs:
`272822, 272823, 272823, 272824, 272825, 272825`. The repeated IDs are
overlapping slices of the source rows:

- `272823`: a 290-sample slice `[0, 290)` and a 360-sample slice `[218, 578)`,
  overlapping by 72 identical samples.
- `272825`: a 2,112-sample slice `[0, 2112)` and a 2,844-sample slice
  `[2010, 4854)`, overlapping by 102 identical samples.

Native-cutoff and accepted-cutoff outputs both contain 7,610 frames. Their
hashes are
`21f8c6629bd7e42c7a86033ffe0d36d7f273f0f6d2a92afb42bb7e15b0e14ecf` and
`a8ec8d359acf9a2475c3698bc6b5e7ae4d6b41397cd22e842f9d2300c7564fba`.
Only 165 samples differ, at the two extra split joins, and every difference
is one PCM count. In this mapped Hello control, the native `>9` cutoff splits
existing source PCM into overlapping windows; it does not append a second
full-unit sound. Captures are in
`corpus-parity/stage20/hello-native-window-inputs/`; reproduce with
`trace-2013-hello-native-window-inputs.gdb`.

### Paired Hi PCM windows and join arithmetic

Read-only hardware-breakpoint captures at 2013 `FUN_1002aac0` and 2006
`FUN_10023b40` saved the two decoded PCM windows immediately before joining
them. The accepted two-position chain is `272820, 272821`, with scale 100 in
both engines. Both supply 1,420 samples for row 0 and 5,910 for row 1. The
captured row buffers match byte for byte across DLL generations (SHA-256
`47ba1f42a8bd102e87fad9d66ed42c5589c86ec2fd21a7dada99acc3c7a15823` and
`c9d05b1fef65ba6b96d2335d9129b2560b45dcf239e66177a0fa78f8592a7130`). The
instrumented WAVs match their controls: old
`331fe0de611f49ded95845376adf93c308d15ad3e98abce6968f1d2d2cbb9c2b`, new
`87310485389615cdfae8f8779372d9125eb3543d809551cc3b292542433bce53`, each
7,258 frames.

For the new 2013 output, onset samples `[0, 80)` are exactly
`trunc(row0[k] * T[floor(k*4096/80)])`; the final 134 samples are exactly
`trunc(row1[5776+k] * T[4096+floor(k*4096/134)])`. The 72-sample overlap at
output `[1348, 1420)` is exactly the sum of separately truncated lanes:
`trunc(row0[1348+k] * T[4096+floor(k*4096/72)]) + trunc(row1[k] * T[floor(k*4096/72)])`.
`T` is the shared 8,192-entry table. The old output copies row 0 `[0, 1348)`
and row 1 `[72, 5910)` directly. Its 72-sample join exactly matches
`trunc(row0[1348+k] * T[floor((72+k+1)*8191/145)] + row1[k] * T[floor((k+1)*8191/145)])`.
The old routine truncates and clamps the combined weighted sum; the new one
truncates and clamps each weighted lane before adding them.

The paired WAVs differ at 283 samples, all inside the 80-sample onset, the
72-sample overlap, or the 134-sample tail envelope. The largest difference is
508 PCM counts in the tail. The row data, selected chain, gain, and duration
are identical, so the remaining difference in this Hi control comes from
generation-specific edge and overlap handling. The result is scoped to this
two-row fixture; it does not yet explain the separate multi-row Hello defect.
Trace scripts are `trace-legacy-hi-window-inputs.gdb` and
`trace-2013-hi-window-inputs.gdb`; captures and restoration records are in
`corpus-parity/stage20/hi-window-inputs-2006/` and
`corpus-parity/stage20/hi-window-inputs/`.

### Equal-gain no-marker Hello source windows

A native-cutoff run on the no-marker matched-context index set the 2013
volume setter to 100, matching the marker-corrected native-cutoff capture.
The selected chain was `273369, 273370, 273370, 273371, 264074, 264074`.
Its six source windows contain 1,416, 648, 520, 1,724, 1,576, and 2,482
samples. Their first-five trailing spans are 84, 84, 90, 86, and 106, so the
result contains 7,916 frames. The marker-corrected six-window chain contains
7,610 frames; equalizing gain and cutoff does not equalize the chosen source
windows. The two native-cutoff WAVs differ at 7,504 of their 7,610 common
samples (maximum absolute difference 24,192), and no pair of captured source
rows shares an exact PCM substring of 16 or more samples.

This same-engine A/B supports the existing marker intervention finding:
setting bit `0x80` in signature byte `+6` for the three legacy-marked
`unit-gen2` rows redirects selection toward IDs 272822–272825. The broad
output change follows changed candidates and decoded windows, so this test
does not measure join arithmetic. It also does not establish the semantic
identity of every selected unit. Captures, WAV, logs, and fixture restoration
record are under `corpus-parity/stage20/hello-no-marker-window-inputs/`;
reproduce the capture with
`trace-2013-hello-no-marker-window-inputs.gdb`.

### Hi scalar-field swaps at both scored positions

The six per-field interventions were repeated on the matched Hi row pairs with
the `5d 00` suffix overlay and the sum-9 cutoff accepted, producing the
two-position chain `272820, 272821`. For position 0, rows `272820` and `65331`
have different values in every scalar array, but the six swaps all return
`10.8` for both rows and preserve the same 14,560-byte WAV
(`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`). The
context record is `00 00 00 00 00 01`. In the decompiled `FUN_100182e0`, a
nonzero context byte `+5` takes the branch that skips the scalar-array
lookups. The 2006 local costs for this pair are `0` and `0.0285714`; 2013
instead separates them on the next transition, favoring
`272820 → 272821` over `65331 → 272821` by `8.1671`.

At position 1, rows `272821` and `280449` share the same legacy seven-byte
feature record, and the context byte `+5` is zero. The controlled 2013 local
returns are `57.118` and `66.5484`. Swaps at model offsets `+0x54` and `+0x5c`
change those returns to `58.0484/65.618` and `55.618/68.0484`; the other four
scalar-array swaps leave them unchanged. All six runs retain the same selected
chain and exact WAV hash above. The old bypass trace gives `272821` cost
`10.3087` at weight `0.333333` and `280449` cost `64.2757` at weight `1`,
also favoring `272821`. These are generation-specific formulas and scales.
The two 2013 arrays alter local scores at this position but do not change the
matched path. This closes the Hi per-field consumer test for the two captured
pairs; the source-column mapping for these Kate rows is recorded below.

The 12 intervention captures are under
`corpus-parity/stage20/hi-2013-feature-field-{0..5}/` and
`hi-2013-current-field-{0..5}/`. The context-byte captures are under
`hi-2013-feature-context-gate/` and `hi-2013-score-contexts/`. Every runner
exited normally and restored the Stage 5 fixture hashes.

### Source columns for the six scalar arrays

The 2013 scorer-state trace now ties all six one-byte arrays to the Kate 2005
index tail. The legacy reader handles three metric groups, each arranged as a
two-byte value followed by two bytes. The 2013 index parser exposes those
groups as column-major `(word16, byte_a, byte_b)` columns; the repacker copies
the legacy `metrics[12N]` unchanged. For each row, scalar slots
`+0x48/+0x4c/+0x50` match metric columns `2/6/10`, and slots
`+0x54/+0x58/+0x5c` match columns `3/7/11`:

| Global row ID | Bank/local row | `+0x48/+0x4c/+0x50` | `+0x54/+0x58/+0x5c` |
| ---: | --- | --- | --- |
| 272820 | `gen2` / 92825 | `85, 126, 168` | `78, 69, 77` |
| 65331 | `gen` / 65331 | `71, 108, 156` | `64, 80, 81` |
| 272821 | `gen2` / 92826 | `168, 171, 127` | `77, 75, 117` |
| 280449 | `etc` / 1324 | `161, 169, 124` | `81, 86, 111` |

The original-index values match the live 2013 scorer arrays for all four rows.
The native 2006 trace also reads `85, 126, 168, 78, 69, 77` for row 272820.
This confirms the byte mapping on the matched rows, without assigning names to
the proprietary dimensions. In the Hi local-cost branch, `+0x54` and `+0x5c`
are the second byte lanes of metric groups 1 and 3. The score interventions
therefore exchanged original metric bytes.

The captures are produced by `trace-2013-hi-scalar-array-state.gdb` and
`trace-legacy-hi-scalar-array-state.gdb`; logs, WAVs, and restoration records
are under `corpus-parity/stage20/hi-scalar-array-state/` and
`hi-2006-scalar-array-state/`. Both runners exited normally and restored the
shared fixtures byte-for-byte. The 2013 state-only trace emitted the native
three-position 14,560-byte WAV with SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.
This clears the tested six-byte transfer as the source of the remaining
generation difference; query-key/candidate construction and the two engines'
distinct score and transition paths remain the active causes under study.

### 2026-09-29: Hi candidate-member crosswalk

`crosswalk-hi-class-members.py` applies the Hi suffix-map repack and the
2013 `FUN_10016ea0` key transform to every candidate ID in the two 2006 pools
captured by `trace-legacy-hi-span-filter-bypass.gdb`. The 10-row first pool
maps as one group to `5a 22 01 a0 00`. The second pool partitions into four
rows with legacy feature bytes `22 11 5a 04 22 00 22`, mapping to
`22 11 5a 65 00`, and five rows with legacy feature bytes
`22 11 5a 04 04 00 22`, left at `22 11 5a 04 00` by the primary overlay.
The source identities and complete row/key output are emitted by the script.

Three controlled repacks map that last group's `04 04` tail independently to
`5d 00`, `55 00`, and `4d 00`. Each yields one exact class lookup and the same
candidate-list counts `67, 19, 19`, selected stream
`272820, 272821, 272821`, and WAV hash
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.
Those equivalent outputs do not identify which target suffix is the intended
field conversion.

The old span-filter trace reduces its 10- and 9-row pools to full-span rows
`272820` and `272821`; native 2006 local costs are `0` and `10.3087`. The
2013 Hi suffix-map run forms lists of 67, 19, and 15 rows and selects
`272820, 272821, 272821`. Its first transition favors `272820` over `65331`
by `8.1671` (`26.709` versus `34.8761` total); repeating `272821` at the
third position adds no transition term. This closes member identity for the
captured Hi pools and distinguishes it from the remaining three-position
path behavior. The forced P AH0 runtime fixture remains unsuitable as a
matched old/new text comparison because the 2006 parser treats the VTML tag
as text. Four of the six generated exact keys still miss before scoring and
enter fallback expansion, so the extra position originates before final path
ranking.

### 2026-09-29: repeated-Hello query-hit membership

`crosswalk-repeat-hello-key-members.py` scans source rows for the four unique
successful keys in the 2013 `Hello hello.` producer trace. Calls 7 and 21
both find `17 2b 2f 00 00`, whose 25 rows include old-selected IDs `264073`
and `273371`. The existing plain-Hello crosswalk establishes that this is the
five-row legacy class plus 20 additional rows. Call 19 finds
`22 17 2b 00 00`; its nine rows include old-selected IDs `264072` and
`273370`, and match the previously crosswalked legacy nine-row class. Calls
10 and 16 find one-row and two-row classes (`2b 30 22 00 00` and
`0d 22 17 00 00`) with none of the nine IDs in the captured 2006 output
stream.

The live 2013 trace reports five successful lookups out of 24; this scan
measures their source-row membership as 25, 1, 2, and 9 rows. The existing
2006 repeated-Hello trace captures its nine selected rows, not every query
pool, so the two lookups containing old-selected rows are not a complete
old/new repeated-input member comparison. The 25-row class's 20 extra rows
remain a concrete source of additional path alternatives. The exact 2013
lookups and counts are in
`corpus-parity/stage20/hello-repeat-2013-producer/adapted-producer-keyclasses-gdb.log`.

### 2026-09-29: repeated-Hello 2006 pools and matched 2013 members

The isolated original-engine capture records 24 split invocations for
`Hello hello.` across three render passes. It prints each old query key,
shortlist ID, runtime weight, and seven-byte feature record. The runner
restores the Stage 5 input/output fixtures and verifies their hashes:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage20/compose-original-msi.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e EVIDENCE_DIR=/probe/stage20/repeat-hello-old-query-pools \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-hello-old-query-pools.gdb runtime \
  /bin/bash /work/stage20/run-original-forced-pah0.sh
```

`crosswalk-repeat-hello-old-pools.py` maps those feature records back to
source-index rows and checks each per-class row count against its captured
weight. For the 2013 side it applies the `attr48/key[:3]/attr40/key[3:]`
signature order, the three class-key tables, and the byte-6 `0x20` mask. The
9-row old slot-1 set exactly equals 2013 key `22172b0000` at call 19. The old
four-row slot-6 set is a subset of those same nine rows. Old slots 2 and 7
each contain the same five rows included in the 25-row 2013 key
`172b2f0000` at calls 7 and 21. Slot 3 has one old row but the successful
one-row 2013 key `2b30220000` identifies a different source row. Slot 0's
75 rows split into source keys `5a22041e00` (45) and `5a22171e00` (30), and
slot 5's 73 rows map across eight effective keys; neither old pool occurs in
the four successful exact-hit classes. Full intersections and row IDs are
printed by the script.

The old slot-8 pool contains six source rows whose stored tail is `04 04`.
The base 2013 signature makes that tail's effective key `04 00`, because it
masks signature byte 6 with `0x20`; the live repeated-Hello target asks for
`45 00` at call 23. A prefix-conditioned index candidate maps stored
`2b305a:0404` to `4500`:

```sh
python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant \
  attr48-key3-attr40-key2-repeat-hello-slot8-0404-to-4500
python3 tools/revkit/work/stage20/verify-repeat-hello-slot8-overlay.py
```

The verifier reports six rows under `2b305a4500`:
`1869`, `5170`, `246677`, `255282`, `264074`, and `272825`. A first diagnostic
variant mapped the effective `04 00` bytes instead of the stored `04 04`
tail and yielded no hit; its trace remains under
`corpus-parity/stage20/hello-repeat-2013-slot8-tail-map/`.

The corrected candidate overlay returns one class entry for key
`2b305a4500`; the following `2b305a4d00` query remains empty. To capture the
split-helper metric and fallback decisions for this candidate, run:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slot8-0404-to-4500.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-slot8-split-threshold.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-slot8-split-threshold \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

## Natural Hi query and membership cross-check

The ordinary `Hi.` fixture (`hi-plain.txt`, SHA-256
`c7e13072a15d4bec0f88822de02bc825d5c6055ac8cf928913d301e336a28c1e`) was
traced on both generations. The original 2006 engine accepts two queries with
10 and 9 candidates and selects `272820, 272821`. On the corrected-repack 2013
baseline, all six exact keys miss. The existing prefix-conditioned suffix map
fills two keys; 2013 then selects `272820, 272821, 272821`. This restores the
leading source rows while leaving one extra handoff and a different WAV. The
mapped second class has weight 4. The active 2013 cutoff is 10, so that class
fails and the builder returns a third position. A cutoff-only DLL copy lowers
the active threshold to 3 while leaving continuity logic native; the builder
then returns two positions and selects exactly `272820, 272821`. Forcing the
same sum-4 decision through GDB produces an identical WAV hash. Candidate
lists initially measure 67 and 4. A separate `04 04 → 65 00` suffix alias
adds the five old feature-B rows to the four existing feature-A rows; that
second 2013 pool then exactly matches the old nine IDs. The first pool's 57
extras came from seven raw key prefixes that the 2013 lookup tables collapse
into the target class. A candidate-prefix control routes those seven prefixes
to `5a2213`, leaving the exact old ten members. With both controls active, live
2013 candidate lists are the exact old 10 and 9 IDs; selection and WAV hash
remain unchanged. Across all scanned Kate index banks, the seven exact
prefix/tail patterns match only those same 57 rows (`gen`, `gen2`, and `etc`;
none in `num` or `alp`). Other query contexts could still use the redirected
signatures, so this remains a Hi-specific key projection experiment, not a
recovered general conversion rule.

With the same two IDs, rewriting the 2013 volume setter argument from 200 to
100 yields a waveform that differs from the old output at only 283 of 7,258
samples. Those differences fall in onset, overlap, and tail envelopes; the
decoded source buffers are byte-identical. The [Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-natural-hi-query-and-member-set-control)
records the hashes, class weights, and comparison boundary. Its [isolated Hi
PCM analysis](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#hi-split-threshold-and-pcm-gain-isolated)
details the generation-specific join arithmetic.

Build the cutoff-only DLL and reproduce the natural Hi selection capture:

```sh
python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant attr48-key3-attr40-key2-hi-exact-first-pool-filter
python3 tools/revkit/work/stage20/make-kate-tree2-active-cutoff-dll.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-hi-exact-first-pool-filter.yaml \
  -f tools/revkit/work/stage20/compose-kate-tree2-active-cutoff-dll.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hi-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-2013-hi-repacked-candidates.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hi-plain-2013-tailmap-cutoff3-exact-pools \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

Captures, including fixture restoration hashes, are under
`corpus-parity/stage20/hi-plain-2006-query-evidence-v2/`,
`hi-plain-2013-query-evidence/`, and
`hi-plain-2013-tailmap-query-evidence-v2/`,
`hi-plain-2013-tailmap-cutoff3-only/`,
`hi-plain-2013-tailmap-cutoff3-featureb-alias/`,
`hi-plain-2013-tailmap-cutoff3-exact-pools/`,
`hi-plain-2013-metric4-override/`, and
`hi-plain-2013-metric4-volume100/`.

As a held-out check, the Hi candidate-prefix map was applied to plain
`Hello.` without its Hi suffix aliases. Candidate counts changed at three
positions (465→408, 111→117, and 23→19), but the selected chain and WAV hash
stayed the same (`273369, 273369, 273370, 273370, 273371, 282025, 282025`; SHA-256
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`). This
shows a cross-context effect: all 57 redirected rows leave position 0, and
positions 1 and 6 also change, although their row-level differences do not
overlap those 57 IDs. Separate query traces show the same ten generated
signature/key/count lines, but those traces did not record returned class
identities. A paired trace at `FUN_10023350` captures its input class lists
and their unit membership, before unit expansion, continuity metadata, or
local scoring. A follow-up `FUN_10023f90` trace records the class keys and
member counts returned by each query. All three affected queries use the
relaxed `FUN_10023e70` → `FUN_10023c70` feature-scope path, called with a
10-class and 10,000-unit limit. The actual `FUN_10023c70` inputs have 84
classes for context 1 and 66 for context 6. Baseline and treatment contain the
same class-key/member-count sets and identical per-key distance scores, but
their traversal orders differ. The changed output classes tie at the ten-class
cutoff: context-1 keys `2922171400` and `3722171400` both score 620, while
context-6 keys `16305e0500` and `41305e0500` both score 1210. The 2013
`FUN_1001b5f0` partition sort can reorder equal scores; paired post-sort traces
confirm the context-1 rank-1 candidate changes from `2922171400` to
`3722171400` and the context-6 rank-9 candidate changes from `16305e0500` to
`41305e0500`, matching both returned shortlists. Context 0 is a separate direct
effect: query key `5a2217a000` returns class key `5a22011e00`, which shrinks
from 67 members to the intended 10 and produces 465→408 expanded units. The selected Hello chain
and WAV hash remain unchanged. Paired class captures are under
`hello-hi-builder-classes-baseline-repacked/` and
`hello-hi-builder-classes-treatment-repacked/`; both runs exited normally,
restored the Stage 5 fixtures, and produced the same 17,836-byte WAV hash
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`. The
full feature-range and per-class distance traces are under
`hello-hi-range-scope-baseline-v2/`, `hello-hi-range-scope-treatment-v2/`,
`hello-hi-scope-distances-baseline/`, and
`hello-hi-scope-distances-treatment/`. The direct post-sort traces are under
`hello-hi-scope-sorted-baseline-v6/` and `hello-hi-scope-sorted-treatment-v6/`.
Per-query class IDs, keys, and member counts are logged under
`hello-hi-query-membership-baseline-v3/` and
`hello-hi-query-membership-treatment-v3/` using the query-membership GDB scripts.
The earlier audio baseline and treatment captures remain under
`hello-hi-filter-unmapped-baseline/` and `hello-hi-prefix-filter-isolated/`;
the paired query traces are in the two
`hello-hi-prefix-filter-query-keys-*.log` files.

Reproduce the Hello prefix-only treatment with:

```sh
python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant attr48-key3-attr40-key2-hi-first-candidate-prefix-filter
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-hi-first-candidate-prefix-filter.yaml \
  -f tools/revkit/work/stage20/compose-selector-legacy-continuity-active-cutoff3-dll.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-2013-hi-repacked-candidates.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-hi-prefix-filter-isolated \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

At slot 12, 2013 records class ID `18857`, weight 6, rejects metric sum 6,
and enters its two-position fallback. The same trace records class ID
`13478`, sum 9, at slot 9; that exact nine-row class also falls back. The
original 2006 trace logs no fallback at the corresponding six-row and
nine-row pools, while its one-row slot-3 query does fall back. The corrected
2013 builder returns 14 positions. This directly confirms that the strict
`>9` cutoff splits both shared pools, while recovering the slot-8 lookup by
itself does not remove all repeated-Hello splits. Both runtime restore files
show identical before/after Stage 5 input and output hashes.

The curated member and threshold findings are in the
[Lead 5 comparison report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-repeated-hello-2006-pool-crosswalk-and-threshold-confirmation).

### 2026-09-29: repeated-Hello 75-row alias and first-selection tie

The 75-row old slot-0 pool transforms into two 2013 source keys: 45 rows at
`5a22041e00` and 30 at `5a22171e00`. The 45-row group is stored under raw
prefixes `5a2205` and `5b2205`; the 2013 lookup tables reduce those to
`5a2204`. The live first query is `5a2217a000`. A candidate-only variant maps
the two raw prefixes to `5a2217` and tail `1e1e` to `a000`, then combines both
groups at that target:

```sh
python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant \
  attr48-key3-attr40-key2-repeat-hello-slot0-pool-to-a000-v2
```

Capture the threshold and fallback decisions with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slot0-pool-to-a000-v2.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-slot0-pool-to-a000-final-split.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-final \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The original 2013 split helper then sees metric sum 75 and accepts slot 0.
Other slots still fall back at sums 1, 2, and 9; the position builder returns
13 positions. The corrected runtime evidence and restoration record are under
`corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-final/`.

With the full pool present, native 2013 ranking selects unit `149799` at the
first position while the old engine selects `273369`. Both units are in the
75-row pre-score list and both receive exactly `9.60000038147` from the
2013 local scorer. The final transition list has 30 candidates, retains
`149799`, and excludes `273369`. Their old query-pool classes differ:
`273369` is class 52165 with weight 2; `149799` is class 52164 with weight 45.
The exact score values are in
`hello-repeat-2013-slot0-pool-to-a000-v4/adapted-score-precision-gdb.log`.
The 75-item pre-score list and 2013 backtrack chain are in
`hello-repeat-2013-slot0-pool-to-a000-v2/`; the original 2006 selected chain
is in `hello-repeat-2006-selected-chain-v2/`.

A single-key intervention tests the first divergence. It changes unit
`273369`'s pre-sort node key from `-1` to `-2`, leaving the candidate index
and feature data unchanged. The selected first unit becomes `273369` at
index 0; the output hash changes to
`b572e5a7d45a9a245f0a76f116c9a980ff4638bfe5477f5eab26c431cfb207a7`.
Reproduce it with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slot0-pool-to-a000-v2.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-slot0-old-winner-order-override.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-slot0-old-winner-order-override \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

This establishes that rank ordering and the 30-candidate cap can cause the
first selected-row difference after lookup and split acceptance. The exact
equal-key permutation mechanism and the corresponding old-engine ordering
rule still need direct comparison. An initial alias mistakenly targeted the
effective class prefix as if it were stored bytes; its split log was
overwritten by the corrected run because the trace script kept the first
absolute output path. Its result artifacts remain, but it is not a complete
capture bundle. The corrected split log has a dedicated path as noted above.

The curated analysis is in the
[Lead 5 comparison report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-repeated-hello-75-row-pool-and-first-chain-divergence).

### 2026-09-29: repeated-Hello slot-5 pool alias

The old slot-5 pool has 73 rows across ten raw-prefix/tail combinations and
eight effective 2013 keys. The live call-16 key `0d22170000` has two native
rows (`4515`, `114916`) under raw prefixes `2f221e` and `2f2218`. Inspect the
partition with:

```sh
python3 tools/revkit/work/stage20/analyze-repeat-hello-slot5-keys.py
python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant \
  attr48-key3-attr40-key2-repeat-hello-slot5-pool-to-0000
```

Capture whether the remapped pool clears the split threshold with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slot5-pool-to-0000.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-slot5-pool-to-0000-split.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-slot5-pool-to-0000 \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The repack maps the 73 legacy source rows onto target raw prefix `2f221e` and
maps tails `0a0a` and `1414` to `0000`. The split trace reports weight 79 and
acceptance at slot 7; the old 73 rows have joined the two native rows, while
the preserved 2013 metrics sum to 79. Other contexts still split, and the
position builder returns 13 positions. Runtime output is 30,872 bytes with
SHA-256 `a3078cb23fcd7392dfa599e19fd0031bd67b84e312d365ff48977ea8968d3fe8`.
Stage 5 before/after hashes are in
`corpus-parity/stage20/hello-repeat-2013-slot5-pool-to-0000/`.

The backtrack trace selects `273369` at context 0, but this isolated overlay
still falls back at slot 0 and changes position segmentation. The contexts
cannot be aligned one-for-one with the 2006 chain; later selected rows remain
unresolved by this test. The selected chain and restoration record are under
`corpus-parity/stage20/hello-repeat-2013-slot5-pool-to-0000-backtrack/`.
The detailed comparison is in the
[Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-repeated-hello-73-row-slot-5-pool-alias).

### 2026-09-29: combined repeated-Hello pool overlay

The combined disposable variant applies the observed slot-0, slot-5, and
slot-8 candidate aliases together:

```sh
python3 tools/revkit/work/stage20/repack-index-key-fields.py --variant \
  attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools
```

The captured split trace accepts slot 0 at sum 75 and the mapped slot-5 class
at sum 79. It still rejects a six-weight class (slot 10), a nine-weight class
(slot 7), and a one-weight class (slot 4) under the native `>9` check. The
builder returns 12 positions. Backtracking begins with units
`149799, 273370, 273370, 273371`, then `208339, 34312, 49874`; the 2006 chain
begins with `273369`. This shows that restoring these lookup populations does
not by itself restore old unit choice. The fixture input/output hashes were
unchanged by both runs; output is 26,400 bytes with SHA-256
`602ce927b25b6e511513285148617feb76c08102fb9f7f797e1d30a90450cd07`.

Run the runtime captures with
`trace-repeat-slots0-5-8-split.gdb` and
`trace-repeat-slots0-5-8-backtrack.gdb`, using the generated
`compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools.yaml`.
The captures are in
`corpus-parity/stage20/hello-repeat-2013-slots0-5-8-pools/`. The rank-cap and
equal-key evidence, plus the remaining class aliases, is summarized in the
[Lead 5 comparison report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-combined-repeated-hello-pool-overlay-and-cutoff-result).

The slot-3 singleton has a separate suffix mismatch: the old candidate maps to
`2b30220200`, while the live 2013 query key is `2b30220000`. Variant
`attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools` maps the old
`0202` suffix to `0000`. This raises the observed metric sum from 1 to 2, but
the `>9` split check still rejects it. Builder count, selected 12-unit chain,
and output WAV hash remain identical to the prior combined capture. This
confirms that the one-row identity mismatch does not explain this fixture's
position-count divergence. Split and backtrack evidence is under
`corpus-parity/stage20/hello-repeat-2013-slots0-3-5-8-pools/`.

An order-only intervention on this combined overlay promotes `273369` to the
first selected row. The remaining sequence is
`273370, 273370, 273371, 208339, 34312, 49874, 264072, 264072, 264073, 264074,
264074`, versus the 2006 sequence `273370, 273371, 232670, 266023, 264071,
...` after its first row. Since the builders return 12 versus 9 positions,
this comparison does not assume index-for-index phone alignment. Output hash
is `184a7d1df5aa92b7822479b7a23612e633f382abf1e02a6e76902c2d14d19b7b`;
fixture hashes match. This shows the equal-rank/cap effect explains the first
choice in this run, while later candidate or edge differences remain.
Evidence is under
`corpus-parity/stage20/hello-repeat-2013-slots0-3-5-8-order-override/`.

Completed `FUN_10018c80` rows explain the later backtrack: context 4 stores
`232670` at 79.1433 and `208339` at 80.6948; context 5 stores `266023` at
98.4015 and `34312` at 99.3981. Although the old route is cheaper through
context 5, old-selected `264071` is present in the 79-row pool but absent from
the 30-row transition shortlist at context 6. The retained `49874` row costs
132.377 through `34312`, followed by `264072` at 158.995 through `49874`.

A correction to the earlier field interpretation is necessary: `81.4368` was
read before the final unit scorer. `FUN_100182e0` returns `59.8` for both
`264071` and `49874`. With native ordering, both enter the second sort with
key `59.8`; the 30-row output includes `49874` but excludes `264071`.

A ranking-key intervention changes `264071` from -1 to -2 at position 6. It
moves to the head of the initial sort and forms a one-row preserved prefix,
so that key survives the final cutoff. The path then uses
`232670 → 266023 → 264071` at contexts 4–6. Its cumulative cost is 140.651;
the next `264072` row is 158.802 through `264071`, 0.193 below the baseline
route through `49874`. Later expanded handoffs remain. Output is 26,374 bytes,
SHA-256 `41d793131f8e57596ad123718bcd17249a781c6e955b0ffce37c46f206a58b8f`;
the completed-row baseline is 27,448 bytes. Both runs restore the Stage 5
fixtures.

A field-preserving pointer-swap control moves `264071` from index 22 to index
0 immediately before local scoring. Its score/key/span/weight fields remain
unchanged, the preserved prefix remains zero, both target rows still score
59.8, and `264071` enters the final 30 rows at index 11. The WAV hash matches
the ranking-key intervention exactly. This confirms that equal-score ordering
and the 30-row cutoff determine whether this old-selected unit reaches path
scoring; the later transition then selects the restored route. The exact
native 2013 ordering cause versus the 2006 reduction remains under analysis.

Scripts are `trace-repeat-slots0-3-5-8-completed-path-rows.gdb`,
`trace-repeat-slots0-3-5-8-promote-264071.gdb`,
`trace-repeat-native-shortlist-sort-boundary.gdb`,
`trace-repeat-promoted-shortlist-sort-boundary.gdb`, and
`trace-repeat-order-only-264071-shortlist.gdb`. Captures are in
`corpus-parity/stage20/hello-repeat-2013-completed-path-rows/`,
`hello-repeat-2013-promote-264071/`,
`hello-repeat-2013-shortlist-sort-boundary/`,
`hello-repeat-2013-promoted-shortlist-sort-boundary/`, and
`hello-repeat-2013-order-only-264071-shortlist/`. See the [Lead 5 comparison
report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-completed-path-rows-and-the-late-shortlist-intervention) for the cost rows and evidence boundary.

### Matched-unit span and rank comparison

The paired shortlist traces show a distinct old/new ranking input for the
repeated-Hello candidate IDs. At old position 5, `264071` has left/right/total
coverage 0/3/4, initial rank key -8, and local score 7.304; it remains at
index 0 in the old 30-row shortlist. At 2013 context 6, the same unit has span
1, weighted coverage 0, and local score 59.8, tied with `49874`; the native
30-row sort removes it. The old engine also uses a 30-row cap, so the
difference is the candidate's coverage-derived rank before truncation.

`FUN_1001cd60` computes old continuity from neighboring per-position candidate
sets and boundary-control bytes. The old model's units `264071`–`264073` map
to phone IDs 3434, 32400, and 25383 and each carry continuation byte `1` in
the separate model `+0x68` array; `264074` carries zero. The 2013 signature
for `264071` is `1,47,34,30,13,0,0`, with continuation bit `0x80` clear in
its final byte. Setting only that bit in process memory changes the 2013 span
from 1 to 3, the local score from 59.8 to 54.8667, and the pre-sort key from
-1 to -3. The preserved prefix becomes one row and `264071` heads the 30-row
shortlist. The output is 26,374 bytes with SHA-256
`41d793131f8e57596ad123718bcd17249a781c6e955b0ffce37c46f206a58b8f`, matching
the earlier rank-key and pointer-order interventions; fixture hashes match.

This proves the missing 2013 marker causes the shortlist exclusion in this
repeated-Hello run. Setting bit `0x80` on all three consecutive 2013 records
`264071`–`264073` raises the target span to 6, lowers its score to 49.9333,
changes its key to -6, and keeps it first in the 30-row list. The output hash
is unchanged. The immediate failure is established: these 2013 model rows
have clear continuation markers where the corresponding old chain has them
set. Span 6 differs from old span 4 because the builders emit 12 versus 9
positions. The repacker source explains the clear marker in this adapted run:
it lays out the new signature as `attr48, key5[:3], attr40, key5[3:]`, placing
old `attr_48=1` in byte 0 while `FUN_100230a0` reads continuation bit `0x80`
from byte 6.

A disposable on-disk overlay writes `signature[row*7+6] |= 0x80` for the
three rows whose legacy `attr_48` is 1. The runtime reports signature byte 6
as 128, span 6, score 49.9333, and first shortlist position. The completed
path selects `232670 → 266023 → 264071` at contexts 4–6 and uses `264071` as
the predecessor of context-7 `264072`. The controlled WAV changes from 27,448
bytes (SHA-256 `184a7d1df5aa92b7822479b7a23612e633f382abf1e02a6e76902c2d14d19b7b`)
to 25,096 bytes (SHA-256
`7c5137c896e850f067d2f9686ef496b2a7ecf2bfd1406c6ee4cb3040e44d6b35`). Both
runs restore the fixtures. This confirms the repacker's field placement as
the immediate cause of the continuity shortlist miss; it does not establish
full audio parity, because the 12-versus-9 position layout and other query
key and split-threshold differences remain. The reproducible builder and
compose layer are `make-repeat-hello-bit7-overlay.py` and
`compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity.yaml`.
Rebuild the disposable overlay and reproduce its shortlist/path capture from
the repository root with:

```sh
python3 tools/revkit/work/stage20/make-repeat-hello-bit7-overlay.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-bit7-disk-completed-path-rows.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-bit7-disk-overlay \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

Evidence is in
`hello-repeat-2006-unit-boundary/`,
`hello-repeat-2013-span-boundary/`,
`hello-repeat-2013-shortlist-unit-signatures/`, and
`hello-repeat-2013-bit7-264071-span/`,
`hello-repeat-2013-bit7-chain-span/`, and
`hello-repeat-2013-bit7-disk-overlay/`; see the [Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-matched-unit-span-and-rank-comparison).

Setting the same continuation bit on consecutive units `264071`–`264073`
extends the target to span 6, lowers its score to 49.9333, changes its key to
-6, and retains it first in the 30-row shortlist. The WAV hash remains
`41d793131f8e57596ad123718bcd17249a781c6e955b0ffce37c46f206a58b8f`. Since
the 2013 builder emits 12 positions and the old builder emits 9, span 6 is not
a direct match for old span 4. The unit-by-unit position crosswalk remains
necessary. Capture and script: `hello-repeat-2013-bit7-chain-span/` and
`trace-repeat-2013-bit7-chain-span.gdb`.

## Native tail-sort permutations

`trace-native-tail-sort.gdb` captures the natural 75-entry local-score tail
from the repeated-Hello candidate overlay. Its exact IDs and float32 score
bits are in `native-tail-sort-order.tsv`. `trace-native-heap-sort.gdb` uses the
same callsite and ID list but overwrites only node scores in process memory:
the first, middle, and last scores are 0, 1, and 2, and every other score is 3.
This creates a 1-to-75 split at the first partition and forces the native
32:1 heap fallback. The captured permutation and score bits are in
`native-heap-tail-sort-order.tsv`; its output WAV and fixture restoration hashes
are in the ignored `corpus-parity/stage20/native-heap-tail-sort/` evidence
directory. The in-process score intervention does not modify model files or
shared input fixtures.

Reproduce the controlled heap-fallback vector from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-native-heap-sort.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/native-heap-tail-sort \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The runner verifies the Stage 5 input and output fixture hashes before and
after execution. The Go selection test replays both captured vectors through
the native sorter and the public tail-sort/cap function.

## Ordinary model-parser offset rows

`trace-model-parser-rows.gdb` records the parser state passed to
`FUN_1003e240` without changing process memory. The offset vectors are
preserved in `model-parser-offset-vectors.tsv`; the source fixtures are
`hello-repeat-plain.txt`, `parser-offsets-punctuation.txt`, and
`parser-offsets-boundaries.txt`. The runtime reports each ordinary `Hello`
word with an exclusive end offset. In `Hello hello.`, the rows are `(0,5)`
and `(6,11)`. In `Hello, world! Hello?`, the first segment has `(0,5)` and
`(7,12)`, then the next segment resets to `(0,5)`. In `Hello. World?
Hello!`, each of the three segments has one `(0,5)` row. Commas and spaces
remain part of the same segment's source coordinate system; `.`, `?`, and `!`
end the current segment, and the following word begins at zero.

Reproduce the three-segment boundary trace from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/parser-offsets-boundaries.txt \
  -e GDB_SCRIPT=/work/stage20/trace-model-parser-rows.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/model-parser-boundaries \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

`text.BuildPaul2013OrdinaryParserOffsetSegments` reproduces these raw offsets
for ASCII-letter words and unsigned cardinals. The cardinal captures show one
row for `1` and `12`, four identical source-span rows for `123`, and three for
`1234` and `2024`; each row retains the entire numeric token span. Other
unsigned-cardinal multiplicities follow the existing number expander as an
implementation inference.

`trace-model-parser-row-type.gdb` separately records row `+0x2c`, which
`FUN_1000ea20` maps into context-row type. The captured values are not a
constant or a function of the `+0x24` discriminator alone: repeated ordinary
words produced 0 and 2, punctuation-separated rows produced 5, 4, and 3, and
numeric rows used a mix of 0 and 4. The complete observed tuple set is in
`model-parser-row-types.tsv`. This closes the field observation; the full
producer remains incomplete, so callers still supply the dword outside the
directly ported punctuation branches.

The additional `trace-model-parser-row-type-writes.gdb` sets hardware
watchpoints on the first four row `+0x2c` fields before `FUN_1003d3d0`. On
`Hello hello.`, the second row changes from 0 to 2 at `FUN_10051a00`; on
`Hello, world! Hello?`, the writes include 0 to 5 at `FUN_100544f0`, 0 to 4,
and 0 to 3 at `FUN_10051a00`. The finalizer later clears the captured fields.
These traces establish the direct producers and bounded punctuation mappings;
the caller-side scanner gates remain open. The
read-only probe restored both temporary stage5 fixtures byte-for-byte. Its
captured log and WAVE are retained under the ignored
`corpus-parity/stage20/model-parser-row-type-writes/` and
`model-parser-row-type-punctuation-writes/` directories.
The isolated `1.25!` probe (`parser-row-type-decimal.txt`) also wrote 2 at
`FUN_10051a00`. It did not exercise the separate, statically recovered dot
branch in `FUN_100544f0` that writes 1. The `USA.` probe likewise wrote 2 at
`FUN_10051a00` and did not exercise that dot branch.
The comma predicate compares the following token with `too` at `0x10078ea0`
and `either` at `0x1007902c` using `FUN_1001c2c0`; the Go helper now applies
those mapped comparisons directly. The `Hello, either. Hello, too.` probe
did not hit the comma type writer, so the static code-12 outcome remains
unconfirmed at runtime. Its temporary fixtures were also restored byte-for-byte.

The word-and-year trace used `parser-offsets-number.txt`; the width matrix used
`parser-offsets-number-widths.txt`. Reproduce the width capture from the
repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/parser-offsets-number-widths.txt \
  -e GDB_SCRIPT=/work/stage20/trace-model-parser-rows.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/model-parser-number-widths \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The isolated numeric-class fixture covers `+12`, `1.25`, `$5.00`, `25%`,
`21st`, `01/02/2024`, `3:45 PM`, and `555-1234`. The trace directly records a
separate sign row for `+`, number rows for the remaining digits, a separate
`%` row, and distinct cardinal spans around the telephone hyphen. Decimal,
currency, ordinal, date, and time rows reuse their full token spans. The
telephone's captured 4 + 1 + 3 parser rows contain cardinalized groups and
the spoken connector `to`. Stage 5's independent phone-row capture also
records `five hundred fifty five to twelve thirty four`, confirming these are
normalized spoken words, not only offset bookkeeping. The Stage 5 record is
`tools/revkit/work/stage5/probes/stage6/records/numeric-exhaustive-capture.log`
(phone-row returns 224–232).

The trace script also reads the primary and auxiliary C strings consumed by
`FUN_1000d190`, at parser-row `+0x34` and `+0x52`. In the read-only numeric
capture under `corpus-parity/stage20/model-parser-row-strings-numeric/`,
separate rows contain `one`, `point`, `two`, and `five` while sharing `(0,4)`;
the currency rows contain `five` and `dollars`; the percentage suffix row
contains `percent`; and the time suffix row contains `PM`. Telephone rows
contain `five`, `hundred`, `fifty`, `five`, `to`, `twelve`, `thirty`, and
`four`. Every inspected auxiliary string was empty. The punctuation fixture
capture under `model-parser-row-strings-punctuation/` showed `Apple` at `+0x34`
and an empty `+0x52` string.

`BuildPaul2013OrdinaryParserOffsetRows` now pairs the bounded number
normalizers with their source-span rows and writes the normalized lexemes at
`+0x34`; mismatched row/text counts fail closed. The directly captured classes
match these strings. Broader number grammars remain an inference, and the
builder still does not produce the auxiliary string.

Reproduce the isolated numeric-class trace from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/parser-offsets-numeric-isolated.txt \
  -e GDB_SCRIPT=/work/stage20/trace-model-parser-rows.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/model-parser-numeric-isolated \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

This remains an offset producer, not raw native-row construction. Signed,
normalization outside the captured forms, abbreviation, exception, TPP,
callback, and source lookup paths remain unimplemented. Other supported
numeric row counts follow the existing normalizers as inference. The GDB logs
and fixture restoration hashes are retained under the ignored
`corpus-parity/stage20/model-parser-*` directories.

## Selector-side legacy continuation DLL experiment

`make-kate-selector-legacy-continuity-dll.py` makes a hash-guarded disposable
copy of `vt_kat.dll`. It preserves the tree2 path changes and changes both
`FUN_100230a0` span scans to read the adapted legacy continuation bit from
signature byte `+0`, bit `0x01`. The source DLL is not modified. The no-marker
matched-context index adapter then selects the same six Hello handoffs as the
existing index-side bit-7 conversion, with an identical WAV hash.

Build the selector-only DLL and run the hardware selection capture with the
unmarked index adapter:

```sh
python3 tools/revkit/work/stage20/make-kate-selector-legacy-continuity-dll.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-matched-context-attrb-transfer.yaml \
  -f tools/revkit/work/stage20/compose-selector-legacy-continuity-dll.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-hello-selected-units-hardware.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-selector-byte0-patch-native \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The ordinary whole-position cutoff is at `FUN_10024060` in Kate `vt_kat.dll`,
VA `0x1002428d`: `cmp eax,0xa` / `setge`, so sums of 10 or more pass. The
separate class-byte-8 branch accepts any positive sum. Empty class lookups
return before either metric decision and go directly to the two-position
fallback. The earlier edit to two `FUN_10023e70` comparisons targeted a later
relaxed-query helper and had no effect on this Hello fixture.

The active cutoff was then changed from 10 to 3 in a disposable DLL. Combined
with the byte-0 continuation patch, and without a GDB return hook, it selected
`272822, 272823, 272824, 272825`. Its 15,264-byte WAV hash exactly matches the
GDB intervention result. Stage 5 input/output fixture hashes match before and
after. The capture is under
`hello-selector-byte0-active-cutoff3-dll/`; the old ineffective candidate is
under `hello-selector-byte0-cutoff3-dll/`. See the [Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-active-cutoff-query-key-and-member-set-map)
for the exact key and candidate-member map and scope limits.

Reproduce the runtime-only cutoff intervention with the selector-patched
DLL:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-matched-context-attrb-transfer.yaml \
  -f tools/revkit/work/stage20/compose-selector-legacy-continuity-dll.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-hello-cutoff-and-selected-hardware.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-selector-byte0-patch-cutoff-override \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The earlier, ineffective relaxed-helper candidate is built with
`python3 tools/revkit/work/stage20/make-kate-selector-legacy-continuity-dll.py --legacy-cutoff`
and mounted with `compose-selector-legacy-continuity-cutoff3-dll.yaml`.

Reproduce the active-cutoff DLL run with:

```sh
python3 tools/revkit/work/stage20/make-kate-selector-legacy-continuity-dll.py --legacy-active-cutoff
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-matched-context-attrb-transfer.yaml \
  -f tools/revkit/work/stage20/compose-selector-legacy-continuity-active-cutoff3-dll.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-hello-selected-units-hardware.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-selector-byte0-active-cutoff3-dll \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

### `attr_b` term in repeated Hello

With the repeated-Hello pool aliases, corrected continuation marker, vector
trees, and first-row rank intervention held fixed, clearing only the cached
`attr_b` score changes the selected unit at output rows 4 and 15 from `232670`
to `205629`. Rows 5 and 6 stay `266023` and `264071`. The WAV changes from
25,096 bytes (`7c5137c896e850f067d2f9686ef496b2a7ecf2bfd1406c6ee4cb3040e44d6b35`)
to 24,470 bytes (`31d866f64d7fdac99252c34ae38dea0dac33b9761fe4199e3fd778d23770b749`).
The first completed-path shortlist divergence occurs at context 4: eight of
the 30 rows are replaced. This demonstrates that an all-zero `attr_b` column
still changes ranking because its looked-up cost is multiplied by a
candidate-dependent duration factor. Here, that default contribution happens
to retain the native-selected row; it is not proof that zero is the correct
legacy value.

Copying legacy `attr_40` into `attr_b` with the seven-byte signature intact
selects `53333` at rows 4 and 15 instead, with `266023` and `264071` unchanged
at rows 5 and 6. Its WAV is 24,154 bytes
(`b501f256fdad3be777011356a503e3bfb917913c4d55c0f2cd6065e47a16dd1f`). This
does not validate `attr_40` as the matching field. The old-format adapter
still needs the correct field mapping or a justified version-specific score
rule. Both runs restore the Stage 5 fixtures; details and evidence boundaries
are in the [Lead 5 comparison report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-attr_b-changes-the-repeated-hello-path).

Reproduce the isolated `attr_b` suppression:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-bit7-disk-attrb-neutralized.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-bit7-disk-overlay-attrb-neutralized \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

Rebuild and run the `attr_40` duplicate overlay:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 -B tools/revkit/work/stage20/duplicate-repeat-hello-attr40-to-attrb.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-repeat-hello-attrb40-bit7.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-bit7-disk-attrb40-duplicate.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-bit7-disk-attrb40-duplicate \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

### Shared predecessor intersections for repeated Hello

The native 2006 follow-up records transition costs for current units
`232670`, `266023`, and `264071`. The first 2013 matched-row trace used the
overlay before the disk-level continuation-bit correction: it recorded 15
candidates for `232670`, 25 for `266023`, and none for `264071`. That absence
was specific to the uncorrected candidate set.

The corrected 2013 trace uses the disk-level bit-7 continuity overlay while
retaining the first-position promotion of `273369`. It records 20, 28, and 19
predecessors for those rows, respectively; the 2006 counts are 5, 22, and 20.
The shared-predecessor winners are `273371` (2006) versus `264073` (2013) for
`232670`, `232670` in both for `266023`, and `266023` in both for `264071`.
The `264071` transition term is `12.3499` in both captures. Context ordinals
and accumulated costs differ, so these are matched unit and predecessor
comparisons, not identical target states. The corrected 2013 WAV matches the
bit-7 disk-overlay baseline at 25,096 bytes and SHA-256
`7c5137c896e850f067d2f9686ef496b2a7ecf2bfd1406c6ee4cb3040e44d6b35`; the
Stage 5 fixtures were restored byte-for-byte.

The separate completed-path backpointer trace confirms the selected 2013
route `264073 → 232670 → 266023 → 264071 → 264072` at contexts 3–7. The
shared-edge minima for `232670`, `266023`, and `264071` match its actual
predecessors in this run. The native 2006 selected chain is
`273371 → 232670 → 266023 → 264071 → 264072`, so both runs share the chain
from `232670` onward but enter that row from different predecessors.

One more edge was captured at current `264072`. Native 2006 evaluates 30
predecessors at context 6; corrected 2013 evaluates 25 at context 7, with 10
IDs shared. Both rank `264071` first in that intersection, with transition
term `0` in both. Their total path costs (`57.7644` and `127.809`) are not a
shared scale. This confirms the same `264071 → 264072` handoff in both
selected paths.

### 2013 cost decomposition before `232670`

Paired score traces hold the first-position promotion of `273369` and the
candidate-pool aliases fixed, changing only the disk-level continuation
overlay. For current `264073` at context 3, the edge from `264072` changes:

The preceding `6.0899` path-cost reduction is accounted for at two earlier
rows. At context 1, `264072` from `252497` has the same previous cost (`4.8`)
and feature (`8.50575`), while the local term falls from `11.1502` to
`8.07509`; total cost falls from `24.4559` to `21.3808`. At context 2,
`264072` from `264072` has feature `0` in both runs, while the local term
falls from `11.0296` to `8.01478`; cumulative cost falls from `35.4855` to
`29.3956`. The disk marker changes the accumulated path before the
`264073 → 232670` decision.

| Overlay | Previous path | Feature | Penalty | Local | Total |
| --- | ---: | ---: | ---: | ---: | ---: |
| Marker uncorrected | 35.4855 | 2 | 0 | 35.1225 | 72.608 |
| Disk bit-7 corrected | 29.3956 | 0 | 0 | 17.5612 | 46.9568 |

The total falls by 25.6512: 6.0899 from the prior path, 2 from the edge
feature, and 17.5613 from the local term, subject to displayed rounding. The
alternative `273370 → 264073` total also falls, from 89.8269 to 72.2656, as
the local term halves; its feature 15.7059, penalty 10, and previous path
28.9985 are unchanged. It still loses to the corrected `264072 → 264073`
route.

At current `232670`, the edge and local terms are nearly unchanged. Without
the marker correction, the 264073 and 273371 incoming totals are 85.3475 and
79.1433. Corrected, the 264073 path total falls to 59.6964, below 273371 at
79.1433. The selected paths therefore diverge at the predecessor entering
`232670`, then share `232670 → 266023 → 264071 → 264072`. This isolates how
the corrected marker changes path accumulation for this fixture; it does not
establish universal score parity or define every use of that marker.

The paired captures are under `hello-repeat-2013-unmarked-prefix-232670/`
and `hello-repeat-2013-bit7-prefix-232670/`. Reproduce the unmarked control:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-predecessors-of-232670-unmarked.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-unmarked-prefix-232670 \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

For the corrected run, replace the pool overlay with
`compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity.yaml`,
the GDB script with `trace-repeat-predecessors-of-232670-bit7.gdb`, and the
evidence directory with
`/work/corpus-parity/stage20/hello-repeat-2013-bit7-prefix-232670`.

The two earlier `264072` rows use the same paired overlay stacks and fixture;
substitute these per-run GDB and evidence paths:

| Overlay | GDB script | Evidence directory |
| --- | --- | --- |
| Unmarked | `/work/stage20/trace-repeat-264072-prefix-edges-unmarked.gdb` | `/work/corpus-parity/stage20/hello-repeat-2013-unmarked-prefix-264072` |
| Disk bit-7 | `/work/stage20/trace-repeat-264072-prefix-edges-bit7.gdb` | `/work/corpus-parity/stage20/hello-repeat-2013-bit7-prefix-264072` |

Reproduce the old-engine capture with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage20/compose-original-msi.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e EVIDENCE_DIR=/probe/stage20/hello-repeat-2006-matched-late-edges \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-2006-matched-late-edges.gdb \
  runtime /bin/bash /work/stage20/run-original-forced-pah0.sh
python3 tools/revkit/work/stage20/compare-common-transition-edges.py
```

Reproduce the corrected 2013 edge trace with the full disk-bit7 overlay stack:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools-bit7-continuity.yaml \
  run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-repeat-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-repeat-bit7-disk-matched-edges.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/hello-repeat-2013-bit7-matched-late-edges \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

This capture is under `hello-repeat-2013-bit7-matched-late-edges/`; the prior
uncorrected capture remains under `hello-repeat-2013-matched-old-row-edges/`.

The next-row captures are under `hello-repeat-2006-next-row-edge/` and
`hello-repeat-2013-bit7-next-row-edge/`. Rerun the Compose commands above with
the same overlays and fixture, substituting these GDB and evidence paths:

| Engine | GDB script | Evidence directory |
| --- | --- | --- |
| 2006 | `/work/stage20/trace-repeat-2006-next-row-edge.gdb` | `/probe/stage20/hello-repeat-2006-next-row-edge` |
| 2013 | `/work/stage20/trace-repeat-bit7-disk-next-row-edge.gdb` | `/work/corpus-parity/stage20/hello-repeat-2013-bit7-next-row-edge` |

### Selector-side continuation read on Hi and Apple

The selector-only DLL patch was checked against a matched Stage 19 tree2 DLL
on natural `Hi.` and `Apple.`. Their 688,128-byte DLL images differ at only
four offsets: `0x23155`, `0x23170`, `0x2320a`, and `0x23218`, the two span
scans' byte-offset and mask immediates. Tree2 path patches are identical.

On Hi, the patched and control DLLs produce the same selected sequence
`272820 → 272821 → 272821`, candidate counts `67, 19, 15`, and 14,560-byte WAV
(SHA-256 `3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`).
The first selected source row has legacy `attr_48=1`; the next has zero. The
flag's presence alone therefore does not imply a changed short output.

On Apple, an initial selector-patched run reported 16 callbacks and 25,090
bytes, but the same hardware-only trace later returned five callbacks and
14,248 bytes. The cause of that discrepancy is unknown; the 16-callback run
is retained as anomalous and should not be used to characterize the patch.

The matched hardware-breakpoint coverage captures both report five positions
and candidate counts `113, 34, 50, 20, 158`. The control selects
`177774 → 177774 → 64719 → 82036 → 82037` and produces 13,212 bytes (SHA-256
`167a388e30f74e83e6d955b39208fb9d62acd6063a3d833c12e0e12c6fef1e23`). The
patched DLL selects `177774 → 177774 → 64719 → 2554 → 2555` and produces
14,248 bytes (SHA-256
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`). The
candidate counts and position count are unchanged, but coverage spans increase
from context 1 onward. At context 3, all 20 candidates move from span 1 to
span 2; in context 4, 20 of the first 60 candidates have span 2 after patching
versus none in the control. The old 2006 neighbor-link trace contains links
for both late pairs (`82036` and `2554` from `189455`; `82037` and `2555` from
`2554`), but the selected `2554 → 2555` chain does not establish old selected
path parity.

The DP trace locates that change downstream: both runs enter the final
position with 30 candidates, while the preceding context retains 16 transition
rows in the control and 19 after patching. Backtracking selects
`82036 → 82037` in the control (`32.1952`, then `40.2424`) and
`2554 → 2555` after patching (`31.0598`, then `36.0851`). The patched context-4
30-row list contains `2555` at index 16 with predecessor `2554`; the control
list excludes it. Both `2555` and `82037` link from `2554` in the old
neighbor-link capture, so this identifies a 2013 coverage-to-retention change,
not complete old/new path parity. Captures are under
`apple-selector-byte0-path-costs-hw-control/` and
`apple-selector-byte0-path-costs-hw-patch/`; reproduce with
`trace-apple-path-costs-hardware.gdb`.

The Hi result remains unchanged at `272820 → 272821 → 272821`, counts
`67, 19, 15`, and 14,560 bytes. Captures include
`apple-selector-byte0-crosswalk-hw-control/`,
`apple-selector-byte0-crosswalk-hw-patch/`, and
`apple-selector-byte0-patch-repeat/`. Use
`trace-apple-continuity-crosswalk.gdb` for the matched metadata trace. The
interpretation and remaining scope are in the [Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#apple-marker-patch-coverage-crosswalk-and-reproducibility-correction).

The native 2006 selected-record capture is under
`apple-selector-byte0-old-selected-records-2026-09-29/`. It records
`177774 → 177774 → 100180 → 59558 → 59559` in each of three render passes.
The patched 2013 DP instead selects
`177774 → 177774 → 64719 → 2554 → 2555`; all three old-selected late IDs
remain present in its transition candidates, with higher cumulative costs
than the chosen alternatives in that run. Reproduce the 2006 capture with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage20/compose-original-msi.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/apple-plain.txt \
  -e EVIDENCE_DIR=/probe/stage20/apple-selector-byte0-old-selected-records-2026-09-29 \
  -e GDB_SCRIPT=/work/stage20/trace-apple-old-selected-records-hardware.gdb \
  runtime /bin/bash /work/stage20/run-original-forced-pah0.sh
```

This confirms a late selection/backpointer divergence for this input, but
does not expose the complete 2006 ranking rule. The companion
`apple-selector-byte0-old-candidate-scores-2026-09-29/` capture samples the
legacy local-score component: it favors old-selected `100180` and `59558` at
contexts 2 and 3, but favors `82037` over selected `59559` at context 4.
That reversal shows the component alone does not determine the legacy chain.
The Lead 5 report records the per-context costs and the paired path trace.

The legacy DP path trace now records the old chain costs and predecessor
writes. On this fixture, 2006 backtracks through
`100180 → 59558 → 59559` at cumulative costs `15.3267`, `27.5946`, and
`27.9018`; patched 2013 instead selects
`64719 → 2554 → 2555` at `18.8205`, `31.0598`, and `36.0851`. The paired
component traces show different context-2/3 feature weights on shared metric
inputs, plus a zero-cost adjacent-unit path in 2006 where 2013 charges a base
term. Substituting the measured old pair-feature subtotals for eight aligned
2013 edges restores the old five-ID chain, while the resulting WAVE hash
still differs from native 2006. Context-3-only and context-2/3 substitutions
instead select `64719 → 86327 → 86328`, so the complete tested change spans
contexts 2–4.

Reproduce the paired 2013 subtotal intervention with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-apple-context-prefix-tail-map.yaml \
  -f tools/revkit/work/stage20/compose-selector-legacy-continuity-dll.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/apple-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-apple-patch-2013-to-legacy-transition-weights.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/apple-selector-byte0-legacy-transition-weight-intervention-2026-09-29 \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

The matched old path-cost/component traces and the 2013 component, full
intervention, and partial intervention evidence are in the Lead 5 report.
The runtime-only edge writes establish selector-level sufficiency for this
Apple input; a production-compatible scoring adaptation and full audio parity
remain unverified.

The detailed evidence boundary is recorded in the [Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-29-later-transition-rows-on-shared-unit-ids).

The broader weight control applies the measured old context-2/3 coefficient
schedules to every 2013 transition edge, without candidate-ID exceptions.
It selects `177774 → 177774 → 188826 → 129559 → 129560`. Weight-only and
weight-plus-consecutive-zero runs produce the same 13,040-byte WAVE hash
`08561d76ee2924dba85f697075e204829ab4403e6b4a766bf36d74e96127f55f`; the
consecutive-zero-only run leaves the patched baseline chain and hash
unchanged.

The native 2006 context-2 boundary capture shows that `188826` is generated:
it is the 15th of 48 candidates after coverage metadata, with local score
`0.222`. The old ranker sorts all rows and drops it outside the top 30. The
2013 builder has 50 rows, preserves a five-row coverage-ranked prefix, and
keeps `188826` fourth in its final top 30. This locates the immediate
cross-generation difference at pre-transition candidate retention, rather
than a missing query hit. The [Lead 5 report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#generalized-weight-and-consecutive-unit-controls) records the paired boundary traces and later old-path costs.

The complete 2013 context-2 candidate list was then captured before scoring
and the shortlist cap. Its 50 IDs share 43 entries with the 2006 48-ID set;
the 2006-only IDs are `19870, 45873, 225775, 251672, 276016`, and the
2013-only IDs are `42412, 92769, 127405, 145945, 258032, 265201, 271727`.
This matches the 43 shared, five old-only, and seven new-only source rows from
`crosswalk-apple-class-members.py`. The capture is under
`apple-selector-byte0-2013-context2-candidate-members-2026-09-29/`; reproduce
it with `trace-apple-2013-context2-candidate-members.gdb` and the same adapted
Apple Compose stack above. The runner exited normally, reproduced baseline
WAVE hash `520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`,
and restored both Stage 5 fixture hashes. Candidate membership therefore
differs independently of the shortlist-retention rule; this trace does not
isolate how the unmatched IDs affect the selected Apple path.

A follow-up runtime edit removed the seven 2013-only rows from context 2,
leaving the 43 common IDs before scoring. The selected chain remained
`177774 → 177774 → 64719 → 2554 → 2555`, and the 14,248-byte WAVE hash stayed
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`. All five
2006-only IDs fall outside its top-30 cap; all seven 2013-only IDs fall
outside the adapted top-30 cap. The retained top-30 lists overlap at 22 IDs;
their eight unique entries per side are all from the 43 common candidates.
Thus, for this position, the candidate membership mismatch is present in the
full pool but does not explain the active shortlist or the selected path. The
capture and restore hashes are under
`apple-selector-byte0-2013-context2-prune-new-only-selected-2026-09-29/`;
reproduce it with `trace-apple-2013-context2-prune-new-only.gdb` and the same
adapted Apple Compose stack above.

Isolated exact-membership controls then restricted patched 2013 context 2 and
context 4 independently to the native 2006 top-30 IDs. The untouched context-4
list was verified as 158 entries with 158 distinct IDs before filtering; both
controls retained the patched chain and baseline WAVE. Context 3 already has
the same 20 source rows in both generations, with no truncation, but the two
rankers select different rows. Thus shortlist membership alone does not
account for the remaining Apple path divergence at these positions. The
captures are under `apple-selector-byte0-2013-context2-legacy-top30-2026-09-29/`
and `apple-selector-byte0-2013-context4-legacy-top30-2026-09-29/`, reproduced
with `trace-apple-2013-context2-legacy-top30.gdb` and
`trace-apple-2013-context4-legacy-top30.gdb`.

Factorial edge-subtotal controls further isolate the scoring boundary.
Replacing the sampled context-4 subtotals alone, context-2 subtotals alone, or
both context-2 and context-4 subtotals left the patched chain and WAVE
unchanged. Replacing the measured context-3 and context-4 subtotals together
selected the native late chain `100180 → 59558 → 59559`; its WAVE is 14,340
bytes with hash `4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8`.
This matches the selected IDs but not native audio bytes. Combined with the
earlier context-3-only and contexts-2+3 controls, the result shows that context
3's ranking terms and context 4's adjacent-unit zero branch are jointly
sufficient in this sampled path; context-2 subtotal restoration is not
necessary here. This remains an aligned-edge runtime intervention, not a
general scoring port. The captures are under
`apple-selector-byte0-context2-legacy-subtotals-2026-09-29/`,
`apple-selector-byte0-context4-legacy-subtotals-2026-09-29/`,
`apple-selector-byte0-context2-4-legacy-subtotals-2026-09-29/`, and
`apple-selector-byte0-context3-4-legacy-subtotals-2026-09-29/`; see the [Lead 5
report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#follow-up-context-factor-controls-for-the-apple-path) for measured edge values and limitations.

A broader control applied `raw×10 + derived-A×2 + derived-B×10 + 2` to every
2013 transition edge scored at context 3, kept context 2 on its patched
formula, and zeroed the three verified context-4 continuation pairs. It
selected `100180 → 59558 → 59559` and produced the same 14,340-byte WAVE hash
as the sampled-edge context-3+4 intervention. The capture is under
`apple-selector-byte0-context3-weights-context4-continuity-2026-09-29/`; use
`trace-apple-patch-2013-context3-weights-and-context4-continuity.gdb`. This
supports applying the old coefficient schedule by context. Its context-4
change still names the three observed pairs, so a general 2013 continuation
condition remains to be implemented and tested.

The old zero-cost branch is now verified directly for those context-4 rows.
At `0x1001e952`, the native 2006 runtime reports selector byte `0` and
previous-row continuation byte `1` for each consecutive pair
`59559←59558`, `82037←82036`, and `2555←2554`. The disassembly at
`FUN_1001e470` tests exactly this adjacent-index plus continuation condition
before writing a zero pair-feature term. All three records were observed in
each of the three render passes, and the native Stage 5 fixture hashes were
restored. The capture is under
`apple-selector-byte0-old-adjacent-zero-branch-2026-09-29/`; reproduce it with
`trace-apple-old-adjacent-zero-branch.gdb` and the native 2006 Apple Compose
stack.

A tightly scoped marker test changed only signature byte 6 bit 7 on rows
59558, 82036, and 2554 in a disposable `unit-gen.idx` overlay. The 2013
runtime read `128` for each predecessor and took its existing adjacent-row
zero-cost branch on `59559←59558`, `82037←82036`, and `2555←2554`. Applying
the old context-3 schedule on all context-3 edges selected
`177774 → 177774 → 100180 → 59558 → 59559` and produced the same 14,340-byte
WAVE hash as the earlier runtime edge-subtotal control:
`4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8`. Both
Stage 5 fixtures were restored byte-for-byte.

The marker-only factorial control took all three 2013 continuation branches
but retained `177774 → 177774 → 64719 → 2554 → 2555` and baseline WAVE hash
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`. The
matched schedule-only control on the unmarked index selected a different path
beginning `273369 → 273370 → 273370 → 273371` and produced a 24,720-byte WAVE
with SHA-256 `ff36d8aa0b5487a9940f4bb422181d463c8fda0a737f88d3500d4873f6960d04`.
Only the combined treatment recovers the old chain and the prior intervention
WAVE hash. The four cells support an interaction between index metadata and
transition scoring for this Apple input. Marker-only evidence is under
`apple-selector-byte0-selected-bit7-only-2026-09-29/`; reproduce it with
`trace-apple-patch-2013-selected-bit7-index-overlay-only.gdb`. The matched
schedule-only evidence is under
`apple-selector-byte0-context3-weights-only-2026-09-29/`; use
`trace-apple-patch-2013-context3-weights-only.gdb` with the unmarked adapted
Apple index.

The complete context-3 candidate pool contains 20 rows, and every matched
2005 source row has `attr_48=1`. An overlay setting bit 7 for exactly those 20
rows (16 in `unit-gen`, four in `unit-gen2`) plus the broad context-3 schedule
reproduced the same selected chain and 14,340-byte WAVE hash as the
three-predecessor marker overlay. The overlay changed exactly 20 bytes and
validated all five index structures. This checks the marker transfer across
the active context-3 pool; corpus-wide conversion remains unverified. The
pool trace is `trace-apple-2013-context3-candidate-pool.gdb`. Build the
overlay with `make-apple-context3-pool-bit7-overlay.py`, mount it with
`compose-apple-context3-pool-bit7.yaml`, and reproduce the combined run with
`trace-apple-patch-2013-context3-weights-context3-pool-bit7.gdb`; evidence is
under `apple-selector-byte0-context3-pool-bit7-weights-2026-09-29/`.

The overlay manifest confirms exactly three index bytes changed. This
resolves the full-index marker experiment's candidate-selection confound; it
does not establish a general mapping from legacy `attr_48` to signature byte
6 bit 7. Rebuild and reproduce with
`make-apple-selected-predecessors-bit7-overlay.py`,
`compose-apple-selected-predecessors-bit7.yaml`, and
`trace-apple-patch-2013-context3-weights-selected-bit7-index-overlay.gdb`.
The capture is under
`apple-selector-byte0-context3-weights-selected-bit7-2026-09-29/`.

The active 2013 table lookup was captured at `0x10018fbb`: Apple context 3
uses slot 1, whose live weights are `(10,2,5)`. The measured 2006 schedule is
`(10,2,10)`, so the exact code/data correction is slot 1's derived-B float
at VA `0x1007c298`, from `5.0` to `10.0`. A hash-guarded script patches only
that float in a disposable versioned DLL. Paired with the 20-row
source-derived context-3 marker overlay, this actual DLL patch selects
`177774 → 177774 → 100180 → 59558 → 59559` and produces the same 14,340-byte
WAVE hash as the GDB controls. GDB only observes branch inputs and selected
IDs in this run. The patch does not yet cover the context-2 table difference
or validate corpus-wide audio parity. Build the DLL with
`patch-apple-context3-weight-dll.py`; mount it with
`compose-apple-context3-bweight10-dll.yaml`. The runner uses
`trace-apple-context3-pool-markers-native-weight-patch.gdb`; evidence is under
`apple-selector-byte0-context3-native-weight-patch-2026-09-29/`.
The active context-2 edge uses slot 2 weights `(10,2,5)` versus measured old
weights `(10,1,5)`. A second DLL with both table corrections and the same
20-row marker overlay produced identical selected IDs and the same WAVE hash.
The context-2 difference therefore does not affect this recovered Apple
output after the context-3 correction. Use
`patch-apple-context3-weight-dll.py --include-context2`,
`compose-apple-context2-context3-weight-dll.yaml`, and
`trace-apple-context3-pool-markers-context2-context3-dll.gdb`; evidence is
under `apple-selector-byte0-context2-context3-native-weight-patch-2026-09-29/`.

### Plain-Hello coefficient transfer control

The slot-1 derived-B patch was checked on `Hello.` with the corrected
vector-pitch trees and matched-context `attr_b`-transfer index. Without
continuity markers, stock and patched DLL runs selected the same six IDs
`273369, 273370, 273370, 273371, 264074, 264074` and produced byte-identical
15,876-byte WAVEs (SHA-256
`8d6b4f63906ad339f5b9b4dde977a495d2c7dfc7a7c9a0818e2a6237d81b3dfb`). The
`0x10018fbb` trace recorded 1,776 coefficient lookups in each run. In 300
calls, slot 1 changed from `(10,2,5)` to `(10,2,10)`; the selected IDs and
WAVE stayed unchanged. This confirms that the patched coefficient is consumed
on the unmarked Hello path, but it does not affect this query's chosen chain.
A second stock/patched pair with identical path-cost tracing returned 111
cumulative candidate rows in each run. Twenty-nine rows changed cost and/or
best-predecessor link: 8 at context 1, 11 at context 2, 9 at context 3, and 1
at context 4; context 5 had no changes. Those later differences propagate
alternative predecessor costs. The selected rows retained the same
cumulative scores: `273370` at 2.98051 and 3.92356, `273371` at 6.45379,
then `264074` at 14.7554 and 15.7349. Candidate `13942` at context 1 had the
largest changed cost, rising from 16.6608 to 23.2013. The coefficient alters
losing alternatives here while leaving the selected route intact.

With the three-row continuity-marker overlay, the patched DLL again produced
the same result as the existing stock marker control: IDs
`272822, 272823, 272823, 272824, 272825, 272825` and 15,264 bytes, SHA-256
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`. The
trace breakpoint was installed and the process reached its ready breakpoint,
but `0x10018fbb` was not visited. That suggests the marker-driven path
bypasses this scorer entry for this input; it does not establish that all
scoring work is bypassed. These controls show that the Apple-derived weight
patch is not a general plain-Hello output fix. They do not invalidate the
Apple result or prove full schedule parity with the 2006 engine. Reproduce
with `trace-hello-active-weight-table.gdb` and
`compose-apple-context3-bweight10-dll.yaml`; captures are under
`hello-context3-weight-control-unmarked/`,
`hello-context3-weight-patch-unmarked/`, `hello-context3-weight-patch/`, and
`hello-context3-weight-patch-trace/` in `corpus-parity/stage20/`.
The paired path-cost captures use `trace-hello-weighted-local-costs-control.gdb`
and are under `hello-context3-weight-path-control/` and
`hello-context3-weight-path-patch/`.

### Cross-generation interpretation for plain Hello

The native 2006 `FUN_1001cd60`/`FUN_1001cff0` path reduces its four position
pools (75, 9, 5, 6 candidates) to their only full-span rows,
`272822 → 272823 → 272824 → 272825`, before local scoring. In the unmarked
2013 path trace, that legacy subsequence remains available and is cheaper
through the first two returned contexts: `272823` costs 2.32276 then 2.85916,
compared with `273370` at 2.98051 then 3.92356. The cumulative ordering
crosses at the next context: `272824` costs 12.5038 versus `273371` at
6.45379. Later, `272825` costs 17.4731 then 19.5923, compared with the
selected `264074` route at 14.7554 then 15.7349. The slot-1 DLL patch leaves
these legacy and selected-route costs unchanged; the two 2013 A/B runs choose
the same route and produce the same WAVE.

The 2006 span-filter-bypass diagnostic provides one direct edge comparison:
`272823` totals 1.64668 through `272822`, versus 32.5011 through `273369`.
That old-engine counterfactual strongly favors the legacy predecessor, but
the native 2006 run normally removes partial-span alternatives before this
scoring stage. This is not an active native candidate comparison.

The three-row continuity-marker overlay changes the 2013 span traversal and
selects `272822 → 272823 → 272823 → 272824 → 272825 → 272825` with the native
cutoff unchanged. This locates the plain-Hello divergence in continuity
metadata and candidate-pool width, upstream of the tested coefficient
correction. Old and new cost values are not a common numerical scale: the
legacy singleton filter and 2013 multi-candidate dynamic program use different
candidate sets and scoring paths. The within-engine A/B and marker controls
show these effects for this input; corpus-wide parity remains unverified.

The old full-span and counterfactual-cost capture is under
`hello-plain-crosswalk-old-span-filter-bypass-v4/`. The marker A/B captures
are under `hello-continuity-bit7/hardware-base/` and
`hello-continuity-bit7/hardware-marker/`.

### Held-out 2005 voice key-partition audit

The proposed signature layout was applied to every row in both local 2005
voices, then projected through the 2013 DLL's three class-key tables. On both
voices, the legacy side reuses the seven-byte feature transform captured from
the 2006 Kate DLL; Bridget has no paired 2006 runtime capture here. Its result
is a held-out structural projection, not an observed Bridget selector result.
The comparison shows how many derived legacy feature classes collapse into
each 2013 exact key:

The three 2013 lookup-table slices used here are byte-identical in
`vt_pau.dll` and runtime-tested `vt_kat.dll`.

| Voice | Rows | Legacy feature classes | 2013 exact keys | Keys merging multiple legacy classes | Rows outside each key's largest legacy class | Legacy classes split across 2013 keys |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Kate | 283,696 | 52,627 | 33,206 | 10,380 (31.3%) | 63,037 (22.2%) | 0 (0.0%) |
| Bridget | 780,882 | 76,893 | 35,198 | 14,119 (40.1%) | 282,904 (36.2%) | 3,764 (4.9%) |

The collision pattern is not confined to Kate rows under the reused 2006
transform. It does not establish Bridget runtime behavior or make each merged
candidate a false positive: the 2013 engine may intentionally rank broader
classes, and this offline audit does not replay queries through both engines.
Reproduce with
`audit-heldout-voice-class-projection.py`; the script reads both 2005 index
sets and DLL tables without writing model data.

The analyzer also compares the structural order
`[attr_48,key[0:5],attr_40]` with the runtime-tested signature repack
`[attr_48,key[0:3],attr_40,key[3:5]]`. A merged pair has different derived
2006 features but the same 2013 exact key; a split pair has the same old
feature but different 2013 keys:

| Voice | Layout | Merged row pairs | Split row pairs | Total partition disagreements |
| --- | --- | ---: | ---: | ---: |
| Kate | Structural | 2,426,245 | 910,383 | 3,336,628 |
| Kate | Signature repack | 4,465,016 | 0 | 4,465,016 |
| Bridget | Structural | 33,123,671 | 10,961,551 | 44,085,222 |
| Bridget | Signature repack | 87,388,861 | 5,193,169 | 92,582,030 |

The structural layout has fewer total pair disagreements in both packages,
but the repack produces fewer old-class splits and more merged classes. These
full-index counts rank the candidate partitions; they do not identify the
correct runtime conversion or prove that either output is more intelligible.

### Bridget 2005 runtime transfer control

The offline partition audit was followed by a same-input runtime A/B on the
second local 2005 voice. `build-bridget-2005-index-adapter.py` builds either
the candidate signature ordering
`[attr_48,key[0:3],attr_40,key[3:5]]` or a structural control retaining
`[attr_48,key[0:5],attr_40]`. Both preserve the 19-byte unit records and 12N
metric bytes, insert a zero `attr_b`, update the index version, and normalize
the producer tag from `VoiceText-Bre` to the 2013 reader's required
`VoiceText-Eng`. Only the temporary index overlay is changed; Bridget's source
data and `binary/vt_kat.dll` are mounted read-only.

For `Hello.`, the signature layout completed at 7,547 frames (15,138 bytes,
SHA-256 `bcb536ddeeff92730603d26b2d924bdd27e605118503a71d2246e5b5cfa0092e`).
The structural control completed at 7,991 frames (16,026 bytes, SHA-256
`9023035e14daa515c2a35499f6e5d4d73301884bf484b5351f69cb36e441ce10`). Both
were valid mono 16 kHz 16-bit PCM WAVs. Their selected-unit callback lists
differed:

| Layout | Selected-unit callbacks |
| --- | --- |
| Signature repack | `718258, 106109, 260196, 260196, 576411, 576411, 39073, 39073, 106109, 718258` |
| Structural control | `178150, 243779, 532206, 682208, 272930, 272930, 733259, 733259, 243779, 178150, 682208, 532206` |

The combined GDB trace shows identical tree query signatures and keys in both
layouts; all ten exact five-byte lookups return zero. Their post-builder
candidate sets differ:

| Position | Signature-repack candidates | Structural candidates | Shared IDs | Signature-only IDs | Structural-only IDs |
| ---: | ---: | ---: | ---: | ---: | ---: |
| 0 | 416 | 569 | 368 | 48 | 201 |
| 1 | 238 | 229 | 224 | 14 | 5 |
| 2 | 47 | 27 | 3 | 44 | 24 |
| 3 | 49 | 23 | 14 | 35 | 9 |
| 4 | 1,848 | 162 | 132 | 1,716 | 30 |
| 5 | 2,015 | 49 | 14 | 2,001 | 35 |
| 6 | 84 | 45 | 23 | 61 | 22 |
| 7 | 46 | 36 | 11 | 35 | 25 |

The observed divergence begins in index-driven fallback candidate
construction after exact lookup, before final backtracking. The 2013 tree
query producer is unchanged. Positions 4 and 5 have the strongest membership
change and the smallest overlap. These traces locate where the adapter alters
this Bridget run; they do not establish native 2006 selections.

`summarize-bridget-span-parity.py` compares the candidate-node span field
(`+0x0e`) across the complete captured pools. Positions 4 and 5 account for
the largest changes:

| Position | Signature total / span 2 | Structural total / span 2 | Shared IDs | Shared IDs with different span |
| ---: | ---: | ---: | ---: | ---: |
| 4 | 1,848 / 1,343 | 162 / 5 | 132 | 127 |
| 5 | 2,015 / 30 | 49 / 5 | 14 | 4 |

All 127 shared position-4 rows whose span rises from 1 to 2 in the signature
layout occur in its position-5 pool and are absent from the structural
position-5 pool. For example, row `16223` has signature bytes
`01 30 2b 32 43 00 0b` versus `01 30 2b 32 00 0b 43`; it is present at both
positions only in the signature layout. A focused breakpoint at
`FUN_100230a0+0x24a` (`0x100232ea`) records fields `(+0x0c,+0x0e,+0x10)` as
`(1,2,0)` versus `(0,1,0)`. This directly links the position-4 candidate
membership difference to 2013 right-match/span metadata before transition
scoring.

The four shared position-5 rows whose span changes the other way are present
in both raw position-4 pools, but only structural retains them in the scored
top-30 shortlist. `trace-bridget-position5-span-producer.gdb` records their
previous-list membership at `FUN_100230a0`; `trace-bridget-rank150652.gdb`
captures row `150652` before and after the position-4 ranking step. Its local
score is `163.89` in both layouts. Its post-coverage rank is 606 of 1,848 in
the signature run and 1 of 162 in the structural run, so it is dropped from
the signature top 30 and retained at structural rank 1. The four position-5
span changes follow that previous-shortlist membership difference.

A new read-only runtime trace reads each returned class key and its unit
members through the model tables at `+0x90`, `+0x8c`, and `+0x94`. For the
three keys containing these target rows, signature and structural layouts
return the same five-byte class keys but different populations:

| Class key | Target rows | Signature members | Structural members |
| --- | --- | ---: | ---: |
| `0d2b2f0000` | `365789`, `382197` | 412 | 2 |
| `102b2f0000` | `548226` | 444 | 1 |
| `232b2f0000` | `150652` | 38 | 1 |

The current key projection reads signature offsets `+1,+2,+3,+5,+6`. Two
disposable repacks isolate the final two consumed fields while retaining
`attr_40` at `+4`. Placing `key[3]` at `+5` and `attr_40` at `+6` reduces the
position-4 candidate count to 827; `150652` ranks 81 and remains outside the
top 30. Placing `key[4]` at `+5` and `attr_40` at `+6` reduces the count to
162; `150652` ranks 1 and is retained. That second treatment reproduces the
structural control's WAVE hash. Together, the paired captures establish that
both trailing selector fields change class population and shortlist order.
They do not establish which experimental placement matches a native Bridget
2006 selector.

These captures supersede the earlier conclusion that the four span flips
were unexplained by the raw adjacent-pool comparison. The earlier large-list
class-ID dump remains excluded because that instrumentation changed the
structural output; the new trace reads only the known model class-key/count/
member-list arrays and retains the established hashes.

The exact target-key lookups remain identical and empty in the paired runs.
The observed divergence occurs after those misses, as relaxed projected
class membership feeds candidate expansion and the top-30 shortlist, then
changes continuity metadata. It is Bridget/2013 adapter evidence; there is no
paired Bridget 2006 runtime here.

Recompute the stable full-pool summary with:

```sh
python3 tools/revkit/work/stage20/summarize-bridget-span-parity.py
```

This confirms that the index field placement changes the active 2013
selection path on Bridget, beyond the earlier Kate-only evidence. It does not
establish intelligibility or parity with a native Bridget 2006 runtime; no
paired Bridget engine is available here. The selected-row crosswalk is
reproducible with `compare-bridget-selected-index-rows.py`.

Rebuild and run the signature treatment:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage20/build-bridget-2005-index-adapter.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage20/compose-bridget-2005-adapter.yaml run --rm \
  -e EVIDENCE_PREFIX=bridget-hello-signature-trace -e TRACE_SELECTED=1 runtime \
  /bin/bash /probe/run-bridget-2005-adapter.sh
```

The structural control uses `--layout structural` and
`compose-bridget-2005-structural-adapter.yaml`. Both runs restored the
Stage 5 input/output hashes. Runtime logs, selected-unit traces, and WAV files
are retained under `tools/revkit/work/stage20/`.

### Kate key5 local-score and signature-field trace

The matched Kate `Hello.` control with `attr48-key5-attr40` records candidate
membership, local scores, and accumulated predecessor costs. A direct runtime
swap of signature byte `+5` between rows `272823` and `264072` changes their
local scores from `1003.89` / `22.3004` to `23.5481` / `1003.58`. The completed
path still differs from native 2006 because the 2013 context-3 candidate list
does not contain `272824`, and later transitions favor a different chain.
Interpretation and limits are in the [Lead 5 Kate report](../../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#2026-09-30-kate-key5-signature-penalty-and-path-replay).

Reproduce the baseline with:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key5-attr40.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/kate-key5-score-path-hello \
  -e GDB_SCRIPT=/work/stage20/trace-kate-key5-score-path.gdb runtime \
  /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

Use `trace-kate-key5-score-fields.gdb` for scorer inputs,
`trace-kate-key5-feature35-swap.gdb` for the two scalar-feature intervention,
and `trace-kate-key5-signature5-swap.gdb` for the one-byte scorer intervention
and resulting path rows. Evidence is retained under the ignored
`corpus-parity/stage20/kate-key5-score-path-hello/` directory. Each runner
capture restored the Stage 5 input and output fixture hashes.
