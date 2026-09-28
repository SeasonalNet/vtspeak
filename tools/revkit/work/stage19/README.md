# Stage 19: Kate tree and index compatibility probe

This experiment uses an ignored, disposable copy of `vt_kat.dll` and a copied
Kate package. It keeps the supplied DLL and vendor model data read-only.

The DLL copy has three guarded byte changes:

- `0x191f`: the common tree loader changes generated `.tree3` names to
  `.tree2`;
- `0x7c8f0`: the voice-tree directory changes from `tree3/` to `tree2/`;
- `0x7c59c`: the `ver.` prefix guard is changed so `FUN_10019e80` takes the
  existing legacy unit-index reader (`FUN_100197d0`).

The supplied 2013 DLL does not contain the old recursive tree reader. The
converter therefore parses Kate's recursive tree2 data and serializes
candidate indexed tree3 payloads into the copied package's `tree2/` directory.
Its 128 deterministic signed feature checks per non-leaf tree compare one
scalar result per lookup. They do not validate the legacy multi-output tree
contract. Four common dictionary tree files are copied under the
requested `.tree2` suffix in a local overlay. The generic converter also
misreads four special vector pitch files (`nbf`, `bf`, `qbf`, and `sbf`) as
one-leaf scalar trees: the old engine uses a separate loader for these files,
and their byte 4 is a vector width. The converted Stage 19 package therefore
does not preserve their old decision trees or vector outputs.

Build the disposable artifacts and run the Wine probe from the repository
root:

```sh
python3 tools/revkit/work/stage19/convert_tree2.py
python3 tools/revkit/work/stage19/patch_compat.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml run --rm runtime \
  /bin/bash /probe/run-probe.sh
```

The probe backs up and restores Stage 5 input/output fixtures. The patched DLL,
full Kate package copy, common overlay, and converted trees are ignored local
artifacts. `conversion-manifest.txt`, `patch-manifest.txt`, runtime result,
and the Wine trace record the experiment.

## Observed result

The legacy-reader variant opened the four common dictionary aliases, all 17
converted Kate trees, Kate's distance table, and its database list. It then
timed out while reading `unit-gen.idx`. Its unchanged Paul control WAVE was not
evidence of Kate synthesis. That variant remains an unsuccessful experiment;
the versioned-index adapter below is the working path.

## Versioned-index adapter follow-up

`convert_legacy_indexes.py` writes a disposable index overlay for the current
versioned reader. It preserves each legacy 19-byte unit record, changes the
version string from `ver.2005` to `ver.2013`, and inserts one zero-filled
unit-count-sized block after the first `N`-byte attribute and `7N`-byte
signature blocks. Inserting one byte into each 20-byte slice instead split
those bulk-read regions and caused the scorer to dereference an invalid
distance row. The corrected converter inserts a full block at the boundary.
It checks every masked word in all three feature
groups against the 1,024-row `cepdist.tbl` dimension. All five indexes retain
their unit counts and 19-byte record stride and pass the current structural
inspector. This is a structurally valid mapping, but its byte-to-column
semantics remain unverified.

`patch_tree2_versioned.py` makes a standard-engine DLL copy with only the two
tree path changes. The converted `ver.2013` indexes therefore use the existing
versioned reader. A GDB breakpoint at `0x10018f17` confirmed that the corrected
adapter reaches candidate scoring with metric codes in range and valid
distance-row pointers. The runtime attempted optional local `class.idx` and
`classhp.idx` files, then opened all five index files, every Kate `.dat` and
`.upm` bank, and the common verification file. It completed synthesis and
wrote a new 88,956-byte, 16 kHz mono 16-bit PCM WAVE for the adapter fixture
(SHA-256 `b82b97019cd422c37e2568a3947c028f3078d7d33bda35f40ae92a6fd2930d68`).
The result differs from the preserved Paul control WAVE.

The same DLL and model adapter also completed the existing Stage 17 prose,
number/time, and address runs. All three returned zero and produced
non-silent PCM WAVE files. The requester listened to the Kate output and
reported that all three are gibberish, with only slight sentence-like
fragments in the prose sample. These are valid audio containers and completed
runtime calls, not evidence of intelligible speech or a compatible voice.

| Fixture | Bytes | Frames | SHA-256 |
| --- | ---: | ---: | --- |
| Prose | 100,144 | 50,050 | `c8dbaf820a15d4d2d198b574a3149306f72e1dff7669d8ab599f86d28eb06083` |
| Number/time | 91,314 | 45,635 | `f27e6879bf25f5b3a02fce8e4efd981287497ffe701e07b1ddf202125414d71a` |
| Address/abbreviation | 104,326 | 52,141 | `a6114ff8c22965241d77845d7ba98e9e9d859d43694124d84a4dd07bf67163c8` |

All files are mono, 16 kHz, 16-bit PCM. Their sample peaks are nonzero. The
captured WAVs, per-fixture runtime logs, `runtime-matrix.tsv`, and fault trace
are retained beside this README. The run scripts back up and restore Stage 5
input/output fixtures.

### Index-column mapping and suffix append control

Two controls now isolate parts of the remaining mismatch. First, a disposable
variant appends the unparsed source bytes to the four converted one-leaf pitch
files. The appended byte counts are 803, 19,931, 245, and 3,691. The adapted
2013 reader produced byte-identical WAVs for all three fixtures with and
without these appended bytes. This only shows that the 2013 reader ignores
bytes after the converted tree3 payload; it says nothing about how the 2006
special loader interprets those bytes as part of its tree format. The captures
and hashes are `runtime-suffix-{prose,numbers,address}.wav` and
`runtime-suffix-matrix.tsv`. Recreate that comparison with:

```sh
python3 -B tools/revkit/work/stage19/make-tree-suffix-variant.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-tree-suffix.yaml run --rm runtime \
  env OUTPUT_PREFIX=runtime-suffix /bin/bash /probe/run-versioned-matrix.sh
```

Second, the legacy reader's call sequence copies a seven-byte key and a
separate one-byte column before its three feature groups. To test whether the
current adapter placed those bytes incorrectly, `make-index-reordered-variant.py`
builds a candidate mapping: `[key7][attribute1][groups12]` becomes
`[zero1][key7][attribute1][groups12]`. This mapping is inferred from the
reader's copy widths; the correspondence has not been confirmed by tracing the
legacy reader on a Kate index. The alternate adapter passed the current
structural inspector, kept all three metric codes in range, and completed all
three synthesis runs. Its WAVs differ from the original adapter's outputs:

| Fixture | Bytes | Frames | SHA-256 |
| --- | ---: | ---: | --- |
| Prose | 68,340 | 34,148 | `be9b43423d56e87030cc973090769b5fbd720e145a25e7ad9f0308237a7e8b6a` |
| Number/time | 75,336 | 37,644 | `3a8ae67b1d26c97cbaceefee94749a1e20569b40358df9a6fd54920cfad5c676` |
| Address/abbreviation | 87,746 | 43,701 | `01e6f81bd952171330141c227be3cac1ab43ef9d60ac041a53206dc7fe424129` |

These files are also valid mono 16 kHz 16-bit PCM, but runtime success and
changed durations do not show whether either mapping selects intelligible
speech. The requester compared the prose samples and judged the original
mapping somewhat more speech-like; the reordered mapping sounded as though it
struggled to form syllables and word parts. Both remain gibberish.
`selection-index-columns.log` and
`selection-index-columns-reordered.log` snapshot the current reader buffers
after loading all five indexes. In the original adapter, the `+0x44` buffer
starts with the source column-1 `0/1` values and the `+0x40` buffer is zero.
In the reordered candidate, `+0x44` is zero for the captured prefix and
`+0x40` starts with varying source column-7 values; `+0x48` is zero in both.
This matches the candidate moving source column 0 into the seven-byte key and
moving the key's last byte into the following one-byte slot. Along with the
legacy reader's 1-byte/7-byte copy widths and the listening comparison, this
rules against the reordered candidate as a useful adapter. The remaining
column semantics and the source-to-engine feature equivalence remain open.
Each trace's matching `*-result.txt` records a successful process exit and the
Stage 5 output hash; the fixture hashes were checked after both runs. Recreate
the alternate with:

```sh
python3 -B tools/revkit/work/stage19/make-index-reordered-variant.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-index-reordered.yaml run --rm runtime \
  env OUTPUT_PREFIX=runtime-reordered /bin/bash /probe/run-versioned-matrix.sh
```

The post-load buffer trace can be reproduced with
`trace-versioned-index-columns.gdb` and
`run-versioned-selection-trace.sh`, setting `TRACE_PREFIX` to
`index-columns` or `index-columns-reordered`, `TRACE_GDB` to the trace file,
and `TRACE_LOG` to `/probe/versioned-index-columns-gdb.log`. Mount the matching
original or reordered Compose overlay for each run.

A focused GDB comparison on the short Stage 19 input confirms that the two
mappings feed different unit-selection data. In the first class-ranking call,
the original mapping presents signature bytes `22 17 2b 30 42 00 00` and one
class candidate (`46204`); the reordered mapping presents
`5a 16 07 2b 18 61 00` and a different candidate (`22588`). For the first
three final-ranking contexts, the original mapping returns 11, 19, and 21 unit
candidates; the reordered mapping returns 30 in each. This is direct runtime
evidence that the index remapping changes the candidate pools, not just the
WAV header or output length. It does not show which candidates are linguistically
correct. The captured traces are `selection-original.log` and the partial
`trace-selection-reordered.log`; the latter ended during further GDB tracing,
so it is not a complete run record.

### Tree branch-order experiment

The converter writes the recursive tree's first child as the indexed reader's
true edge and its second child as the false edge. To test the alternative
child-order interpretation, `make-tree-branch-variants.py` generated three
separate copies of all 17 converted trees: reverse only `C` comparison edges,
only `D` membership edges, or both. Each output still passes the current
`tree3` structural parser. The index adapter, DLL patch, input fixtures, and
other model files were held fixed.

| Variant | Fixture | Bytes | Frames | SHA-256 |
| --- | --- | ---: | ---: | --- |
| `C` edges reversed | Prose | 100,144 | 50,050 | `c8dbaf820a15d4d2d198b574a3149306f72e1dff7669d8ab599f86d28eb06083` |
| `C` edges reversed | Number/time | 91,314 | 45,635 | `f27e6879bf25f5b3a02fce8e4efd981287497ffe701e07b1ddf202125414d71a` |
| `C` edges reversed | Address | 103,840 | 51,898 | `5f0caa43e601cc2fc3680055a314525404b785917cf55121091a42509e56bf25` |
| `D` edges reversed | Prose | 98,042 | 48,999 | `52e2ed9861c6b9dacf9352a956a86fbdabd4d0e72480e30ec219bf95d7c9ecb0` |
| `D` edges reversed | Number/time | 86,216 | 43,086 | `1e72eff14d5dc3d6b24879daf6b2db6c609e28222ace20bd5d7bd658364b8682` |
| `D` edges reversed | Address | 103,050 | 51,503 | `372f12c766d17cd41079ab371fbec7901134b7fcc5c60badeea8898e650188f3` |
| Both edge types reversed | Prose | 98,042 | 48,999 | `52e2ed9861c6b9dacf9352a956a86fbdabd4d0e72480e30ec219bf95d7c9ecb0` |
| Both edge types reversed | Number/time | 87,782 | 43,869 | `49dd080e79a66c12e1c27e942ac4dbb81ecece070700e756e81b2551228a02db` |
| Both edge types reversed | Address | 103,050 | 51,503 | `372f12c766d17cd41079ab371fbec7901134b7fcc5c60badeea8898e650188f3` |

The `C`-only prose and number/time files are byte-identical to the original
adapter's outputs; the address file differs. `D`-edge reversal changes all
three files. Reversing both types matches the `D`-only output for prose and
address, while number/time differs. This says these tree paths can affect
output and that the `D` edges dominate the first and third tested paths. It
does not establish which branch interpretation the legacy evaluator uses or
whether any candidate is more intelligible. Captures are
`runtime-invert-{c,d,all}-{prose,numbers,address}.wav`; per-variant runtime
logs and matrices have matching `runtime-invert-*-matrix.tsv` names. The
variant manifest records each tree's node count, swap count, and hashes.

Recreate the structural variants and runs from the repository root:

```sh
python3 tools/revkit/work/stage19/convert_tree2.py
python3 tools/revkit/work/stage19/convert_legacy_indexes.py
python3 tools/revkit/work/stage19/patch_tree2_versioned.py
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage19/make-tree-branch-variants.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-tree-invert-c.yaml run --rm runtime \
  env OUTPUT_PREFIX=runtime-invert-c /bin/bash /probe/run-versioned-matrix.sh
```

Repeat the Docker command with `compose-tree-invert-d.yaml` and
`compose-tree-invert-all.yaml`, changing `OUTPUT_PREFIX` to
`runtime-invert-d` and `runtime-invert-all`, respectively.

Recreate and rerun that variant from the repository root with:

```sh
python3 tools/revkit/work/stage19/convert_tree2.py
python3 tools/revkit/work/stage19/convert_legacy_indexes.py
python3 tools/revkit/work/stage19/patch_tree2_versioned.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml run --rm runtime \
  /bin/bash /probe/run-versioned-matrix.sh
```

This establishes that the patched standard engine can load the adapted
resources, traverse candidate scoring, and write non-silent WAVs. The
requester's listening review rejects the output as intelligible Kate speech,
so functional voice compatibility remains unresolved. The tree converter
still serializes legacy recursive nodes into the standard reader's indexed
layout under a `tree2/` directory; the standard DLL's native recursive tree2
parser has not been restored. A same-input Stage 20 trace shows identical
results for four initial scalar lookups, but different outputs for paired
vector lookups. The old special loader reads 12-value vector trees from the
`nbf`, `bf`, `qbf`, and `sbf` resources; the generic Stage 19 parser instead
misreads these files as one-scalar leaves and stops after 9 bytes. This is the
specific gap in the default Stage 19 conversion. The inserted index byte's
meaning and the remaining index feature semantics are still open. Stage 20
now contains an experimental parser and converted overlay for the four vector
pitch files; the default Stage 19 package remains the misparsed baseline. See
the
[Lead 5 comparison](../../../docs/reverse-engineering/lead5-kate-common-resource-comparison-2026-09-26.md#matched-ordinary-text-tree-trace).
The reversible differential and preserved traces are documented in
[Stage 20](../stage20/README.md).

### C-edge localization and runtime path trace

After the requester reported that the C-reversed output sounded slightly
better but remained gibberish, the C-edge effect was split by tree family and
then by individual duration tree. Reversing C edges in duration trees alone
reproduced the changed address WAV; reversing pitch-tree C edges left all
three fixtures byte-identical to baseline. Among the nine duration trees,
only `duration/cnas.tree2` changed the address WAV. Reversing its 33 C nodes
produced 103,840 bytes with SHA-256
`5f0caa43e601cc2fc3680055a314525404b785917cf55121091a42509e56bf25`; each
of the other eight duration-tree variants produced the baseline 104,326-byte
address WAV (`a6114f...`). Those eight were `caff`, `capp`, `cfri`, `cstop`,
`vdi`, `vlong`, `vsch`, and `vshort`; the per-tree WAVs and runtime logs are
named `address-tree-<stem>.wav` and `address-tree-<stem>.log`. The C-only prose
and number/time outputs remain byte-identical to baseline.

The new `trace-cnas-path.gdb` captures each scalar lookup of the 281-node
`cnas` tree in the address fixture: its caller, feature vector, traversed
nodes, C comparison value/threshold, and returned scalar tree value. The calls
come from `FUN_10013380` at return address `0x1001350b` (the context slot at
`+0x178`). For example, the first call compares feature 7 value 1 with
threshold 2 and returns 1518 with source edge order, versus 1338 when every C
edge is reversed. Equality cases also change, including node 221 (1 == 1),
node 111 (2 == 2), and node 39 (1 == 1).

To isolate those equal-value decisions, `make-tree-cnas-equality-variant.py`
reverses only nodes 39, 111, and 221, selected from the observed address trace.
This is a fixture-specific experiment, not a recovered global branch rule.
Its full three-fixture matrix matches baseline exactly for prose and
number/time, and matches the full C-reversed address WAV byte-for-byte. The
traced values at these equality nodes change (for example, node 221 returns
1263 instead of baseline 885; node 111 returns 649 instead of 1068). Thus
these three choices are sufficient to reproduce the address WAV change in
this fixture; neither this equality-only variant nor the broader C-reversed
variant establishes intelligible speech or the legacy evaluator's rule.

| Equality-only fixture | Bytes | SHA-256 | Comparison |
| --- | ---: | --- | --- |
| Prose | 100,144 | `c8dbaf820a15d4d2d198b574a3149306f72e1dff7669d8ab599f86d28eb06083` | baseline-identical |
| Number/time | 91,314 | `f27e6879bf25f5b3a02fce8e4efd981287497ffe701e07b1ddf202125414d71a` | baseline-identical |
| Address | 103,840 | `5f0caa43e601cc2fc3680055a314525404b785917cf55121091a42509e56bf25` | identical to full C reversal |

The requester has listened to the C variant and reports it sounds a bit
better, but remains garbled/gibberish. The fixture they auditioned was not
identified, so the subjective comparison is recorded without assigning it to
the address-only change. Runtime evidence shows the address fixture is the
only tested sample whose WAV changes under C reversal.

Recreate the runtime trace and equality-only matrix from the repository root:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage19/make-tree-branch-variants.py
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage19/make-tree-cnas-equality-variant.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml run --rm runtime \
  env FIXTURE=address TRACE_PREFIX=cnas-baseline \
  /bin/bash /probe/run-versioned-fixture-trace.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-tree-invert-c-duration-cnas.yaml \
  run --rm runtime env FIXTURE=address TRACE_PREFIX=cnas-invert-cnas \
  /bin/bash /probe/run-versioned-fixture-trace.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-tree-invert-cnas-equality-address.yaml \
  run --rm runtime env OUTPUT_PREFIX=runtime-cnas-equality-only \
  /bin/bash /probe/run-versioned-matrix.sh
```

The trace WAVs, logs, results, equality-only matrix, and generated variant
manifest are stored in this directory. The run scripts restore the Stage 5
input and output fixtures after each probe.

### Legacy index layout and attribute-B probes

The native 2006 versioned-index reader `FUN_100123a0` reads the legacy tail as
`attr_4c[N]`, `attr_48[N]`, `key[5N]`, `attr_40[N]`, then three metric groups
totaling `12N`. The 2013 reader expects `attr_a[N]`, a unit-major
`signature[7N]`, `attr_b[N]`, and metrics. The adapter that produced the
original experimental outputs only inserted `N` zero bytes at offset `8N`;
it did not transpose the old per-field arrays into per-unit signatures. The
bytes presented as the 2013 signature therefore mix `attr_48`, all five key
columns, and `attr_40`, then get chunked every seven bytes as if each chunk
belonged to one unit.

The 2013 key builder `FUN_10016ea0` reads signature positions 1, 2, 3, 5,
and 6 to make its five-byte key. A structurally constrained candidate is
`[attr_48, key[0:3], attr_40, key[3:5]]`; this keeps the five old key bytes in
all five consumed slots and the two remaining legacy bytes in slots 0 and 4.
For `unit-etc` row 1286, this reconstructs
`02 13 42 18 21 00 1e` from `attr_48=02`, key `13 42 18 00 1e`, and
`attr_40=21`. This fit is strong structural evidence, not paired-generation
confirmation.

The original `signature-last-*` variants remain useful as historical
experiments, but their names were misleading. Their `[7N:8N]` source window
is the genuine legacy `attr_40` field. The copy variant duplicates that field
into 2013 `attr_b` while leaving it in the un-repacked signature. The transfer
variant copies it into `attr_b` and zeros its original location in that
malformed signature. Neither variant repaired the seven-byte signature
layout.

The four fixture runs all exited 0 and wrote valid WAVs. `attr-a-copy` is
byte-identical to zero-fill baseline for all three fixtures. The other
variants produce:

| Variant | Inserted `attr_b` | Signature handling |
| --- | --- | --- |
| `attr-a-copy` | Copy legacy column 00 | Preserve all seven bytes |
| `signature-last-copy` | Copy legacy `attr_40` | Preserve its source bytes in the malformed signature |
| `signature-last-transfer` | Copy legacy `attr_40` | Zero its source bytes in the malformed signature |
| `constant-8` | Fill every unit with 8 | Preserve all seven bytes |

The `signature-last-*` names reflect an earlier layout interpretation and are
misleading. Their `[7N:8N]` source window is the old `attr_40` field. Because
the adapter has not rebuilt a unit-major signature, that field sits at the
end of the bulk `7N` region and is misread as part of the signature. Copying
it to `attr_b` duplicates the field in the malformed view; transferring it
zeros its original bytes. These runs do not test the corrected per-unit
mapping described in Stage 20.

The three-fixture runs all exited 0 and wrote valid WAVs. `attr-a-copy` is
byte-identical to zero-fill baseline for all three fixtures. The other
variants produce:

| Variant | Fixture | Bytes | SHA-256 |
| --- | --- | ---: | --- |
| `signature-last-copy` | Prose | 75,882 | `f74b95f834cc23f2fb390fd6e4b882386d6db60ecff87d35841e463705f3f69a` |
| `signature-last-copy` | Number/time | 72,584 | `7ecbcc62d39b36c781bb9146fbf8c467f623ea6908cba6911b78cba937f06222` |
| `signature-last-copy` | Address | 88,774 | `1a638f315b83ffda9827b78a1feb783d1915ad8cd2b6d0f4b444be987de25f2e` |
| `constant-8` | Prose | 101,220 | `6d262a3431e875b2c8f1740b7348100cfe2e6611c1cbd9a03b6ceec86d055c10` |
| `constant-8` | Number/time | 91,314 | `f27e6879bf25f5b3a02fce8e4efd981287497ffe701e07b1ddf202125414d71a` |
| `constant-8` | Address | 109,082 | `11078d7d5b0e00aaff3653e52740bb671bcbeba4c387cf762661a33e196852af` |
| `signature-last-transfer` | Prose | 73,108 | `c4daeac316906032a4e92c6a082f7a571ad8b1069538b50a3cfc65b477250adb` |
| `signature-last-transfer` | Number/time | 77,602 | `27947c0263ccc718c81448f639f730f1a391d63d3277630e93465d4e079b16d4` |
| `signature-last-transfer` | Address | 79,314 | `beb8601a56803fe2733b03318e8fdcc6461a857324072ca1f6425c32e243afa5` |

The attribute-cost trace confirms that copying legacy `attr_40` into
`attr_b` restores candidate-specific pair costs where zero-fill gave the
same cost to every candidate in the traced context. This establishes a real
scoring effect, not the intended mapping. The old distribution comparison
remains valid only as a descriptive comparison of legacy `attr_40` against
Paul's 2013 `attr_b`; byte-distribution similarity does not establish semantic
equivalence.

A separate Stage 20 experiment rebuilt the 7-byte signature using the
constrained ordering above. It preserved the first four scalar tree results
and both recovered vector rows but changed the class ID, candidate sets, and
selected units. Transferring `attr_40` into `attr_b` after this rebuild also
restored variable score costs while changing the selected sequence. The new
`attr_b` field has no proven `ver.2005` counterpart; move, duplicate, and
transformation remain testable alternatives.

No speech recognizer is installed in the local environment, and WAV hashes or
durations do not measure intelligibility. The generated files are retained
as `runtime-index-<variant>-<fixture>.wav` for listening. The old candidates
are not accepted mappings; the rebuilt signature experiment narrows the
layout but does not yet establish intelligible synthesis.

Recreate the variants and their matrices from the repository root:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage19/make-index-attribute-variants.py
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage19/compare-index-column-distributions.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-index-signature-last-copy.yaml \
  run --rm runtime env OUTPUT_PREFIX=runtime-index-signature-last-copy \
  /bin/bash /probe/run-versioned-matrix.sh
```

Repeat with `compose-index-constant-8.yaml` and
`compose-index-signature-last-transfer.yaml`, using matching output prefixes.
The candidate-cost and candidate-selection traces are
`selection-attributes-signature-last.log`,
`selection-attributes-constant-8.log`, and
`selection-selection-signature-last-transfer.log`.

### Exhaustive Kate DAT decoder parity

The legacy package had previously been checked for indexed DAT/UPM bounds and
sampled decoded lengths, but its decoded waveform bytes had not been compared
with the standard engine. `probe-kate-dat.c` calls the patched standard DLL's
`FUN_10001b30` at module offset `0x1b30` for every indexed Kate record in the
`gen`, `gen2`, `num`, `etc`, and `alp` banks. The probe checks each return
length against four times the sum of that record's UPM periods and captures a
SHA-256 for each DLL PCM buffer. `compare-kate-dat-probe.py` decodes the same
DAT spans with the independent Python implementation and compares every PCM
length and digest with the DLL capture.

All 283,696 records passed both checks: 179,995 `gen`, 98,133 `gen2`, 997
`num`, 3,666 `etc`, and 905 `alp`. There were zero length, UPM, or PCM hash
mismatches. This establishes local decoder-output parity across Kate's
observed DAT corpus with the patched standard DLL. It does not validate the
legacy index feature semantics, converted-tree decisions, candidate-to-unit
selection, prosody, or whole-synthesis intelligibility.

The full 23 MB per-unit text capture is retained locally under
`tools/revkit/work/corpus-parity/kate/kate-dat-dll-parity.tsv`; its SHA-256 is
recorded in `kate-dat-parity-manifest.txt`. The PE32 source and Python
comparison helper are portable and tracked here. Run from the repository root
with the local vendor inputs and Docker images present:

```sh
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$PWD:/src:ro" \
  -v "$PWD/tools/revkit/work/stage19:/out" \
  vtspeak-pe32-builder:local \
  i686-w64-mingw32-gcc -O2 -Wall -Wextra -Werror -std=c11 \
  -Wl,--no-insert-timestamp /src/tools/revkit/work/stage19/probe-kate-dat.c \
  -o /out/probe-kate-dat.exe -ladvapi32
mkdir -p tools/revkit/work/corpus-parity/kate
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage19/compose-index-signature-last-copy.yaml \
  run --rm runtime /bin/bash -lc \
  'Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >/tmp/kate-dat-xvfb.log 2>&1 & xpid=$!; trap "kill $xpid 2>/dev/null || true" EXIT; export DISPLAY=:99; sleep 1; WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all wine /probe/probe-kate-dat.exe --all > /work/corpus-parity/kate/kate-dat-dll-parity.tsv 2> /work/corpus-parity/kate/kate-dat-dll-parity.log'
PYTHONDONTWRITEBYTECODE=1 PYTHONPATH=tools/revkit/scripts \
  python3 tools/revkit/work/stage19/compare-kate-dat-probe.py
```

These results narrow the gibberish away from compressed waveform decoding:
the indexed audio payloads and their UPM timing counts agree with the standard
decoder over the complete Kate corpus. The zero-fill baseline produces
constant pair costs for the traced target `attr_b`; the experimental copied
byte window changes those costs and the output WAVs. A later layout audit
showed that this window is not a per-unit signature byte or a verified field
mapping. The 2006 and 2013 tree formats also differ for the special vector
pitch resources, and the ordinary tree2-to-indexed conversion remains an
unproven semantic translation. The available traces do not yet identify which
adapted unit sequence corresponds to the intended phonetic stream.
