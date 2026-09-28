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
