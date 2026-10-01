# Lead 5: Kate MSI common-resource comparison (2026-09-26)

## Result

The Kate MSI contains the original 2006 shared English pronunciation resources
alongside its voice model. The Stage 19 experiment adapted Kate's per-voice
`tree2` and index data for the supplied standard DLL, but used the checkout's
current `data-common` resources. It also copied four current `.tree3` common
trees under `.tree2` names. That makes the experiment a mixed-generation
package. The old and current embedded dictionaries also use different
phone-ID tables: a controlled ID remap makes the address fixture's output
byte-identical to the current-data control. This establishes one concrete
cross-generation incompatibility that affects synthesis. Prose and number
outputs still differ, and the remapped output's intelligibility has not yet
been assessed.

All 34 Kate M16 files in the MSI are byte-identical to the local `data-kate/M16`
copies. This rules out an extraction mismatch in those voice-specific files.

The listening comparison now separates the license symptom from the earlier
garbling. The requester reports clean, intelligible output from the original
2006 MSI engine both with its matching MSI resources and with the mismatched
workspace verification record; the latter only prepends a cleanly synthesized
demo/nag sentence. The Stage 19 standard-engine adapter output remains
gibberish. Since those runs change the engine generation and the tree/index
adaptation together, this comparison rules against licensing as the cause of
the garbling but does not yet identify which compatibility mapping is wrong.

## Extraction and coverage

Source: `/home/seasonal/Downloads/Neospeech/Neospeech Mike9012/Kate3(1).exe`,
234,313,945 bytes, SHA-256
`a1b1b0d189f8d9751e2d34f73fc9701dac695293e70c06c6ba9fdd1243b42377`. Its
InstallShield overlay starts at `0x3e000`, declares three embedded files, and
ends at EOF after their declared lengths. The extracted
`NextUp.com-NeoSpeech Kate16 Voice.msi` is a valid Windows Installer database,
234,051,584 bytes, SHA-256
`b3f5b3022a02037d83352878a279695e5b095b37031b1df79959c37563185351`.
The MSI metadata identifies InstallShield 11 Express and a save date of
2006-05-31. The setup and MSI were not executed.

The MSI File table lists 70 installed files: 34 under `data-kate/M16`, 21
under `data-common/dict-eng`, two user-dictionary files, one verification
file, ten program/library/font files, and two package-level documents. The
entire MSI extraction is retained locally, outside Git, at
`tools/revkit/work/corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/`.
`msiinfo` and `msiextract` from msitools 0.106+repack-2 were unpacked under
`/tmp`; 7-Zip 26.00 extracted the MSI's embedded cabinet. The outer stream
layout was checked against the public [ISx InstallShield extractor
notes](https://github.com/Coldblackice/InstallShield-installer-extractor-ISx/blob/master/ISx.c).

The MSI also supplies a matching-era `lib/vt_eng.dll` and `lib/TTSApp.exe`.
They are available in that local extraction and were used in the matched
runtime and parser experiments below.

## Common dictionary comparison

The MSI's embedded pronunciation family is structurally valid under the
existing resource inspector, but it is a different corpus from the current
checkout:

| Resource | Kate MSI (2006) | Current `data-common` | Direct comparison |
| --- | ---: | ---: | --- |
| `engttsdict_emb` | 2,055,146 bytes; 227,310 records | 2,157,937 bytes; 228,591 records | Different |
| `hashidx_emb` | 909,244 bytes; 227,310 entries | 914,368 bytes; 228,591 entries | Different |
| `hashcont_emb` | 454,620 bytes; 227,310 entries | 457,182 bytes; 228,591 entries | Different |
| `hashparams_emb` | 132,132 bytes | 132,132 bytes | Different bytes and hash parameters |
| `exceptdict` | 1,345 bytes; 55 rows in groups 2–4 | 2,891 bytes; 123 rows in groups 1–4 | Different |

The two indexed embedded dictionaries have 227,236 raw keys in common. Of
those, 170,776 key/payload pairs are byte-identical and 56,460 have different
payload bytes. The MSI has 74 keys absent from the current corpus; the current
corpus has 1,355 keys absent from the MSI. These are parsed record keys and
payloads, not decoded word-to-phone labels. Runtime tracing below shows that
many payload differences are systematic phone-ID recodings between DLL
generations, rather than changed pronunciations.

The MSI does not contain the current four-file TPP family
(`tppdict_eng`, `hashidx_eng_tpp`, `hashcont_eng_tpp`,
`hashparams_eng_tpp`). It contains four older `city1.txt2` through
`city4.txt2` tables instead. After uppercasing and replacing spaces with the
hyphens used by decoded TPP keys, every unique old city key is present in the
current TPP corpus. Its primary class agrees for 7,924/7,942 city1 keys,
3,971/3,972 city2 keys, 296/296 city3 keys, and 17/18 city4 keys. Twenty keys
have a changed primary class. This shows a close successor relationship for
these place-name tables, with some reclassification; it does not establish
identical behavior for the old and new engines.

## Trees and text tables

The Kate MSI carries common dictionary trees as `tree2`; the current checkout
has `tree3` files. The existing `tree2` to `tree3` converter produces bytes
identical to current `poly.tree3` and `sbd.tree3` from the MSI's `poly.tree2`
and `sbd.tree2`. Those two common trees are therefore not a byte-level content
mismatch after that conversion. `ENGBI_comp.tree2` parses to 41 recursive
records, but its converted bytes do not match current `engbi.tree3`. The
448,653-byte `atmt.tree2` does not parse with the recovered recursive reader
(the first operation byte is `0xeb`); its layout remains unresolved. The
current `atmt.tree3` is a 27-tree container of 292,420 bytes.

The Stage 19 common-resource directory differs from current `data-common`
only by four additional `.tree2`-named aliases: `engbi`, `poly`, `sbd`, and
`atmt`. Their bytes match the corresponding current `.tree3` files, so Stage
19 did not load the original MSI's common tree payloads. This preserves the
old format's request names for the new engine while leaving the new tree
content in place.

The older `.txt2` files lack the current format's 11-byte row-count trailer.
After applying the observed `0x0a` byte transform and removing that legacy
terminator, decoded rows compare as follows:

| Kate MSI table | Rows | Current table | Rows | Difference |
| --- | ---: | --- | ---: | --- |
| `abbrc.txt2` | 330 | `abbrc_sort.txt2` | 330 | Identical rows and order |
| `abbrh.txt2` | 55 | `abbrh_sort.txt2` | 55 | Identical rows and order |
| `abbrt.txt2` | 42 | `abbrt_sort.txt2` | 42 | Identical rows and order |
| `sbdw.txt2` | 1,067 | `sbdw_sort.txt2` | 1,067 | Identical rows and order |
| `chc.txt2` | 2,582 | `chc_sort.txt2` | 2,581 | `gm|0011` and `hp|1111` replaced by `gm|0111` |
| `citya.txt2` | 116 | `citya_sort.txt2` | 117 | Current adds `CRT|Court` |
| `streeta.txt2` | 357 | `streeta_sort.txt2` | 347 | Eleven MSI-only aliases; current adds one |
| `streetf.txt2` | 203 | `streetf_sort.txt2` | 211 | All MSI rows remain; current adds eight rows (listed below) |

The eleven MSI-only `streeta` rows are `DVD|Divide`, `HOLLOWS|Hollow`,
`ISLES|Isle`, `LANES|Lane`, `LOOPS|Loop`, `PATHS|Path`, `PIKES|Pike`,
`RANCHES|Ranch`, `TRACES|Trace`, `TRACKS|Track`, and `VIA|Viaduct`. The
current-only `streeta` row is `EXPWY|Expressway`. The eight current-only
`streetf` rows are `Hollows`, `Isles`, `Loops`, `Paths`, `Pikes`, `Ranches`,
`Traces`, and `Tracks`.

The current `data-common/dict-eng` also has `wab.txt2`, which is absent from
the MSI. The MSI bundles `data-common/userdict/user_eng.dict` (2,067 bytes),
which has no same-name current file, and
`data-common/userdict/userdict_eng.csv` (280 bytes); the current CSV at that
path is empty. The MSI's `data-common/verify/verification.txt` is 315 bytes,
while the current file is 468 bytes. These files were not part of the
pronunciation-resource swaps. The verification-file difference did affect the
original-engine control, as documented below.

## Controlled Stage 19 resource swaps

We reran the same Stage 17 `prose`, `numbers`, and `address` fixtures with the
same supplied standard DLL, patched tree2 adapter, Kate M16 data, TPP/text
tables, and converted common-tree aliases. A current-data control reproduced
the original Stage 19 WAVE hashes exactly. The isolated MSI-embedded variant
replaced only `engttsdict_emb`, `hashidx_emb`, `hashcont_emb`, and
`hashparams_emb`; it retained the current `exceptdict` and all other common
resources. All three runs exited zero and produced mono, 16 kHz, 16-bit PCM
WAVs:

| Fixture | Current control | MSI embedded family | Frame change | MSI variant SHA-256 |
| --- | ---: | ---: | ---: | --- |
| Prose | 50,050 frames / 3.128 s | 41,689 / 2.606 s | −16.7% | `ac35dfe3470afdc48b34d9f75d522f6d144cc39b0b812cae3c45925ae99f001f` |
| Numbers | 45,635 / 2.852 s | 44,552 / 2.784 s | −2.4% | `0fa838c11b400112af9579f8204736c24592eafb8cacbb50ccc66c81b590140d` |
| Address | 52,141 / 3.259 s | 49,716 / 3.107 s | −4.7% | `949695c331193e07fad718980d7647bdadf8925eda9372744bc1581dfe5cd72c` |

The result is consistent with the embedded family changing the generated phone
or unit sequence: replacing that family changes every WAVE while holding the
engine and voice model fixed. The prose change is substantial. Duration and
hash differences alone do not establish which pronunciations changed or
whether the MSI variant is more intelligible; listen to the paired files
before drawing that conclusion. Outputs are retained locally under
`tools/revkit/work/corpus-parity/runtime-current-control-*.wav` and
`runtime-old-embedded-*.wav`.

We then replaced only `exceptdict` with the MSI copy, retaining the current
embedded family. Its three WAVE hashes matched the current control byte for
byte. This exception-table swap has no audible or measurable effect on these
fixtures; it does not rule out an effect for input that reaches one of its
entries.

A GDB lookup probe at `0x10011820` did not fire during the run, so this
experiment did not capture individual query keys or returned pronunciation
payloads. The WAVE changes establish resource sensitivity, not the exact
lookup-to-phone path. The run script restored the shared Stage 5 input and
output fixtures on exit.

The first matched-generation attempt used the supplied 2013 wrapper and failed
because it imports `VT_SetCommaPause_ENG`, absent from the MSI's 2006 DLL. We
then replaced that wrapper dependency with a small C caller that resolves only
the four exports present in the old DLL: `VT_GetTTSInfo_ENG`,
`VT_LOADTTS_ENG`, `VT_TextToFile_ENG`, and `VT_UNLOADTTS_ENG`. This caller
loaded the original DLL, obtained its default speaker and 16 kHz output rate,
loaded Kate, and successfully wrote all three requested WAVs. The earlier
wrapper failure and fixture restore are retained as historical evidence; the
successful direct-caller run is `runtime-original-msi-2006-api.log` and the
parser trace is `runtime-original-msi-2006-parser.log` in the ignored local
evidence directory.

### Matched 2006 runtime: files opened and text path

The earlier successful run mounted the MSI's `lib/vt_eng.dll`, MSI Kate M16
directory, and MSI `data-common/dict-eng` subtree. Wine's file trace shows the
engine opening the MSI's `hashparams_emb`, `engttsdict_emb`, `hashidx_emb`, and
`hashcont_emb` from `../data-common/dict-eng/`, as well as `exceptdict`, the
common `.txt2` tables, and `atmt.tree2`, `poly.tree2`, and `sbd.tree2`. It also
opened Kate's `.tree2` resources, all five `.dat` acoustic banks and matching
`.upm` maps, unit indexes, distance table, and `dblist.idx`. The verification
path resolved to the current 468-byte workspace record rather than the MSI's
315-byte record. The engine attempted to open `class.idx` and `classhp.idx`,
which are absent from the MSI package. It still returned a successful load
and synthesized all three outputs. This does not establish the purpose or
consequence of the two missing class indexes.

The old DLL's pseudocode gives this high-level path (function names are
addresses because their original symbols are unavailable):

1. `VT_TextToFile_ENG` selects the WAVE output path and calls
   `FUN_10017420`.
2. `FUN_10017420` normalizes/copies the input through `FUN_10015370`, runs
   text callbacks including `FUN_100157b0`, then initializes the segment
   state with `FUN_10020750`.
3. `FUN_10020750` repeatedly calls `FUN_100202f0`. For each remaining input
   span, `FUN_100202f0` calls `FUN_1001c180`, which reaches the parser at
   `FUN_100546d0`. That parser calls `FUN_10054700`, `FUN_100547d0`, and
   `FUN_10055050`; the handler chain includes special text, number, and
   punctuation routines.
4. The parser returns records with a `0x88`-byte row stride. The caller reads
   each row's source start/end offsets at `+0x14` and `+0x18`, builds context
   arrays in `FUN_1001c440`, then proceeds through `FUN_1000f3e0` and
   `FUN_1001b800`. The later path calls `FUN_1000f5a0` for duration-tree work
   and `FUN_1000fd30` for phone/unit-tree work; recursive traversals include
   `FUN_10001570` and `FUN_100019b0`.

This establishes a staged text-record/context/tree path into acoustic-unit
selection. The dictionary lookup and its old-engine phone decoding are now
identified below. Selected leaf/unit IDs in the supplied standard-engine run
remain untraced.

The GDB parser trace for the run using the newer 468-byte verification record
caught 16 parser calls across the three fixtures. It shows the parser
receiving additional text before the requested sentences:

| Fixture | Text before requested sentence | Requested sentence parsed as |
| --- | --- | --- |
| Prose | `Thank you for using VoiceText. You need to contact our sales support team to obtain a valid verification file.` | `Hello from the VoiceText Stage 1 runtime check.` |
| Numbers | `This is a demo of VoiceText TTS.` | `The temperature is 72 degrees at 5 PM.` |
| Address | `This is a demo of the VoiceText English TTS system.` | `Meet me at 123 Main St. in Albany, NY.` |

The extra spans are injected on the legacy verification/demo path; the
supplied text files contain only the requested sentences. The output WAVs
therefore include that extra speech and cannot be compared directly by
duration or hash with the Stage 19 runs. Different injected banners also
confound cross-fixture or cross-run listening comparisons. The trace proves
that input reaches the parser and that it emits populated source-span records:
for example, the requested prose is consumed in 49 characters as eight
records. Numeric and address spans contain repeated overlapping records for
`72`, `123`, and `NY`. This is direct evidence of multiple parser records for
those spans, but their semantic role (alternatives, context variants, or
another representation) is not yet established.

The newer verification file was successfully opened, but the engine followed
the demo/verification-text path. At that point it was unclear whether this
reflected the record contents or an environment-dependent check.

With the newer record, the banner text passes through the same parser and
synthesis path before the requested sentence. It changes duration, hash, and
listening context; the requested prose is still parsed afterward. A matched
control using the MSI's own 315-byte verification record later showed that
the unpatched DLL accepts that record under the same Wine environment and
does not inject a banner. The initial banner was therefore caused by the
verification-file mismatch, not an established registry, machine, or time
dependency. The Stage 19 standard-engine `T AH0` control also has no banner.

The earlier banner-bearing WAVs are retained locally as
`runtime-original-msi-2006-{prose,numbers,address}.wav`. A fresh matched
unpatched control with the MSI verification file produced 51,512, 46,566, and
54,510 frames in `runtime-original-control-2006-{prose,numbers,address}.wav`.
The unpatched and one-byte status-patched controls have byte-identical output
for all three fixtures. The requester listened to the `runtime-banner-free`,
`runtime-license-success`, and `runtime-license-status-only` outputs and
reported that Kate is fully intelligible in all three: no garbling, audible
nags, audio errors, or artifacting. Each of those three output sets has the
same per-fixture SHA-256 as the matched unpatched MSI-verification control.
This is direct requester listening evidence for the tested 2006 MSI engine,
Kate M16 voice files, and matching common resources.

The requester also listened to the unpatched
`runtime-original-msi-2006-{prose,numbers,address}.wav` set, generated with
the 468-byte workspace verification record. Those outputs are also clean and
intelligible, with the verification/demo sentence audibly prepended. The
verification mismatch changes the utterance's content; the listening report
does not indicate that it degrades the requested speech or introduces
garbling, artifacting, or synthesis errors.

The comparison now isolates the licensing effect on these fixtures:

| Run | Verification input | Behavior and output |
| --- | --- | --- |
| Original DLL, no patch | Newer 468-byte workspace record | Parser trace shows demo/nag text prepended to each requested sentence; requester reports the audio is clean and intelligible, including the prepended nag. |
| Banner-free branch patch | Same 468-byte workspace record | Skips demo-text insertion; WAVs match the clean set byte for byte. |
| Original DLL, matched MSI control | MSI's 315-byte record | Accepts the record and parses only the requested sentences. |
| License-status or status-only patch | MSI's 315-byte record | Byte-identical WAVs to the original matched control; the status-patched parser trace also contains only the requested sentences. |

The invalid or mismatched verification path therefore changes the synthesis
input itself. The inserted phrase goes through the same parser and downstream
text-to-speech path as the requested text, so its phones and units become part
of the generated audio. Skipping that insertion alone reproduces the
valid-verification WAVs byte for byte, isolating the injected demo text as the
cause of the waveform difference for these fixtures. The requester reports
both versions as clean and intelligible; these observations show a content
difference, not a demonstrated audio-quality improvement from licensing
success. In the matched MSI case, forcing the status to success has no
additional waveform effect because the original DLL already accepts the MSI's
own verification record.

This listening result applies to the matched 2006 MSI runtime and resources.
It does not establish that the Stage 19 standard-engine `.tree2` adapter now
produces intelligible Kate audio; that remains a separate compatibility test.

### Embedded dictionary lookup and phone-ID decoding

The runtime lookup trace now reaches the dictionary-to-phone link. The old
DLL normalizes each token with `FUN_1000cd50`, hashes the encoded key with
`FUN_1000e730`, selects a slot from `hashidx_emb`, compares the candidate key
byte-for-byte, and copies its payload from `engttsdict_emb` (`FUN_1000e440`).
`FUN_100040a0` interprets payload flags and expands numeric IDs using the
256-by-5-byte table at runtime address `0x1009e9a0`. The GDB trace records the
surface token, normalized key, hit/miss, payload, and decoded symbol bytes in
`tools/revkit/work/corpus-parity/runtime-original-msi-2006-dict.log`. These
function names are Ghidra pseudocode names, not vendor symbols.

For example, `Hello` hits MSI payload `11 da cf 26`; the old runtime expands
its IDs to `22 17 2b 30`, or `HH EH0 L OW1`. The current corpus stores
`11 d9 ce 26`, which the 2013 phone table expands to the same symbols. If the
2013 decoder reads the old IDs without conversion, IDs `da` and `cf` instead
mean `IH0` and `EH1`; the resulting sequence is `IH0 EH1 L OW1`. `This` and
`from` show the same systematic pattern: MSI `01 ce dc ef` and `01 d8 39 e4`
become current `01 cd db ee` and `01 d7 39 e3` while preserving their decoded
phone-symbol sequences. This directly demonstrates a generation-specific
phone-ID incompatibility.

The old and current ID tables agree at the same numeric ID for 185 of 256
slots and differ at 71. A full 256-row crosswalk, including empty, ambiguous,
and unrepresented slots, is in
[the Kate phone-ID crosswalk](lead5-kate-phone-id-crosswalk-2026-09-26.tsv).
There are 253 old IDs with a unique current equivalent, two empty old IDs
whose expansion is ambiguous between current IDs 0 and 255, and one old ID
with no current equivalent: old ID 184 expands to `02 36`. Of 38,320 changed
common records whose old payload has the direct-phone flag, 671 contain that
unrepresented ID. For the other 37,649 records, bytewise remapping translates
30,659 payloads exactly to the current corpus payload; 6,990 still differ.
The old decoder also has a second payload form: first byte `02`, then one or
more `context-bytes | phone-IDs` branches separated by `ff`. It preserves the
context bytes (subtracting one into a separate feature buffer) and expands
each branch's phone IDs through the same table. Among changed shared keys,
238 old payloads use this form: 121 become byte-identical to current after
phone-ID remapping; 106 remain different. Of those 106, 82 have the same
branch/context bytes but different phone sequences, three have matching phone
sequences but changed context bytes, 13 differ in both, and eight changed
from this branch form to another payload form. Four additional old branch
records have a layout the recovered parser does not accept, and seven contain
unrepresented ID 184. The remaining 17,902 changed old payloads have other
flag values and are not decoded as phone sequences by this function. These
counts cover all shared raw keys, not only the words in the three fixtures.

The complete indexed-corpus inventory by the old decoder's visible flag
branches is:

| Record set | Direct phone flag | Alternative-branch flag | Other flags |
| --- | ---: | ---: | ---: |
| MSI corpus (all keys) | 39,709 | 246 | 187,355 |
| Current corpus (all keys) | 53,171 | 252 | 175,168 |
| Common byte-identical key/payload pairs | 1,369 | 8 | 169,399 |
| Changed common records, MSI payload class | 38,320 | 238 | 17,902 |
| Changed common records, current payload class | 50,565 | 240 | 5,655 |

“Direct” means payload flag bit 0 is set. “Alternative” means the old
decoder's mask test selects its `02` branch parser; those records carry `|`
and `ff` delimiters. “Other” is the remainder and is not assigned a semantic
label. These counts account for every indexed record in both corpora and all
56,460 changed shared payloads.

### Standard-engine phone-ID remap experiment

We made a disposable copy of the MSI embedded family and changed only payload
phone IDs with a unique same-symbol mapping. The dictionary file length and
all resource offsets were preserved; hash parameters, index files, key bytes,
other common resources, Kate voice files, and the patched standard engine were
unchanged. The experiment skipped all 671 direct payloads containing old ID
184; it also remapped phone IDs after `|` in recognized alternative branches
while preserving their context bytes. Other payload flags were left untouched.
The script and local manifest are
`tools/revkit/work/corpus-parity/transcode-legacy-embedded-phones.py` and
`variant-old-embedded-idremap-v4/phone-id-remap.json`.

| Fixture | Unconverted MSI family | Phone-ID remapped MSI family | Current control | Remap vs current |
| --- | ---: | ---: | ---: | --- |
| Prose | 41,689 frames / 2.606 s | 48,930 / 3.058 s | 50,050 / 3.128 s | Different WAVE |
| Numbers | 44,552 / 2.784 s | 43,525 / 2.720 s | 45,635 / 2.852 s | Different WAVE |
| Address | 49,716 / 3.107 s | 52,141 / 3.259 s | 52,141 / 3.259 s | Byte-identical WAVE |

The address WAVE after remapping has the same SHA-256 as the current-data
control (`a6114ff8c22965241d77845d7ba98e9e9d859d43694124d84a4dd07bf67163c8`).
For that fixture, correcting old phone IDs restores the current engine's
exact output while retaining the MSI dictionary family. Prose and numbers
move toward, but do not reach, current-control durations and retain different
hashes. Those differences may involve genuinely changed entries, parser
behavior, other payload forms, or tree/model interpretation; duration alone
cannot distinguish them. This establishes a phone-ID mismatch as one
cross-generation defect, but does not prove that it is the reason the audio is
unintelligible or that fixing it makes the voice legible. The new WAVs have
not yet been heard. The extended alternative-branch conversion generated
byte-identical WAVE hashes to the direct-only remap for all three fixtures.
Adding this conversion did not affect these utterances, although it accounts
for 238 changed shared dictionary records.

The extracted MSI `lib/vt_eng.dll` (598,016 bytes) and `lib/TTSApp.exe` are
also distinct from the supplied `binary/vt_kat.dll` and
`binary/voicetext_kate.exe` by whole-file hash. The original pair is therefore
a useful future runtime control, rather than evidence that Stage 19 used it.

## Controlled explicit-phone selection trace

To remove dictionary lookup from the experiment, we sent single visible words
through the standard engine with `<vtml_phoneme alphabet="x-cmu">` overrides
and traced the final unit ranking. The disposable inputs and GDB captures are
under `tools/revkit/work/corpus-parity/phone-boundary-probe/`.

First, `Hello.` was synthesized under the current common dictionary and the
phone-ID-remapped MSI dictionary with the same Stage 19 Kate resources. Both
runs reached the selector with the same seven-byte key signature
`22 17 2b 30 5a 00 00` (the controlled sequence HH EH0 L OW1 followed by two
zero fields), returned the same ranked units, and wrote the same WAVE
(`f0abad28d3a8b423e3d59dc03f9bfe2cca095e3ccc9f394be5cd6a8ae3a869db`). For
this word, the dictionary-family difference is already normalized before unit
selection and does not explain the bad sound.

We then forced `P AH0`, `T AH0`, `B AH0`, and `F AH0` on the same visible word. The
forced-phone path did not enter the common-dictionary key-ranking breakpoint;
it proceeded into final unit ranking. We mapped the returned flattened unit IDs
back to the original Kate index banks in `dblist.idx` order and read their
seven-byte signatures. The phone byte in signature position 2 is consistent
with the requested consonant for every one of the first ten returned units for
P (`35`), B (`13`), and F (`20`); for T (`39`), seven of the first ten match,
while the other three contain `4a` or `4b` in that position. The latter values
have no established phone interpretation in this index-signature context.
This is evidence that the index adapter preserves the central phone field in
these controlled cases, not proof that the rest of each signature or the tree
and feature tables are compatible.

The first-call T exceptions do not persist through the next ranking call: its
first ten rows all carry `39`. In all four probes, the first two final-ranking
calls return the requested consonant byte in signature position 2, and the
third returns AH0 (`07`) for all ten inspected rows. Across those 120
candidate rows, 117 carry the expected consonant or vowel byte for their
corresponding ranking call; the three exceptions are the `4a`/`4b` T rows in
the first call. This narrows the concern: the direct-phone path and the
adapted unit signatures agree on these segment bytes, while the reason for
large local costs and unintelligible audio remains unresolved.

The requester initially described the full forced P AH0 output as “pah” and
T AH0 as “tih.” After listening to the isolated decoded rows, they report
P AH0 sounds like a clear “I,” T IH0 has hiss like /s/, T AH0 is very short
and sounds like “su,” and P IH0 is too short to judge. These are perceptual
reports about the generated files, not recovered phone meanings. A matched
P/T × AH0/IH0 probe traced the actual selected unit handoff at
`FUN_1001b200`. After collapsing repeated forward/backward
handoffs, the unique selected units grouped by their central phone are:

| Forced phones | Onset units | Vowel units | Signature phone bytes |
| --- | --- | --- | --- |
| P AH0 | 274307 (GEN) | 248899, 251906 (GEN2) | P `35`, AH0 `07` |
| T AH0 | 154608 (GEN) | 53417, 123517 (GEN) | T `39`, AH0 `07` |
| P IH0 | 282432 (ETC), 102934 (GEN) | 271868, 269729 (GEN2) | P `35`, IH0 `23` |
| T IH0 | 42440 (GEN), 196073 (GEN2) | 210306, 269729 (GEN2) | T `39`, IH0 `23` |

The pre-synthesis trace for T IH0 also records internal phone bytes `39 23`
and the selected index rows carry T then IH0. Thus the initial T AH0 “tih”
impression is not explained by the front end simply substituting IH0 for AH0.
The later isolated-unit impression is different, and the side-level runtime
comparison below narrows which waveform spans need to be assessed. Whether the
selected spans themselves are poor examples or later timing and
coarticulation alter their percept remains unresolved.

We decoded representative selected vowel units through the already
parity-checked DAT decoder and saved standalone WAVs. T AH0 unit 53417 decodes
to 1,218 samples; P AH0 unit 248899 to 2,274; T IH0 unit 210306 to 2,356; and
P IH0 unit 271868 to 516. These files are complete DAT records, but the
runtime does not play each complete record as one phone. Each record contains
two contextual sample spans with a shared overlap, and timeline mode 1 selects
the first span while mode 2 selects the second. `FUN_1002c120` computes the
second-span DAT start as `base + first_span - shared_overlap`; the timeline
row's `+0x0c` sample count is the selected side's span.

A downstream trace at `FUN_1002c220` resolves the actual side choices for all
four controls. Its DAT offsets and durations match the original legacy index
record fields and the corresponding Kate DAT banks:

| Phones | Mode 1: first-side unit and samples | Mode 2: second-side unit and samples | Selected total |
| --- | --- | --- | ---: |
| P AH0 | 248899, 1,102 samples | 251906, 544 samples | 1,646 |
| T AH0 | 53417, 496 samples | 123517, 482 samples | 978 |
| P IH0 | 271868, 298 samples | 269729, 542 samples | 840 |
| T IH0 | 210306, 1,226 samples | 269729, 542 samples | 1,768 |

For example, T AH0's mode-1 row reads 496 samples at DAT offset
`0x02a0c074`, exactly the base offset for unit 53417. Its mode-2 row reads
482 samples at `0x05a8fe2c`, equal to unit 123517's DAT base plus its 506
sample first span less the 102-sample shared period. P AH0's mode-2 row
similarly starts 516 samples after unit 251906's base and reads its 544-sample
second span. This directly verifies the legacy index-to-DAT mapping for these
selected units; it does not prove that the selector chose linguistically or
acoustically appropriate records.

We also saved unadjusted concatenations of the exact mode-1 and mode-2 sample
spans for each control as `timeline-selected-sides-{pah0,tah0,pih0,tih0}.wav`
under the ignored `tools/revkit/work/corpus-parity/phone-boundary-probe/`
directory. `timeline-selected-sides.tsv` records the selected global ID,
physical bank/local row, raw DAT hash and offset, span lengths, and matching
runtime WAVE hash. These concatenations omit engine timing adjustment, gain,
and boundary reconstruction, so they isolate which portions of the selected
records contribute without reproducing the final phone sound. The user's
whole-record impressions cannot distinguish a bad selected span from an
unused half or from later timing/coarticulation.

The requester listened to all four side-selected composites and reports that
they sound like the intended phones. For these forced-phone controls, this
reduces the likelihood that the selected first/second DAT spans themselves or
their index-to-DAT offsets explain the garbling. The composites omit engine
timing, gain, and boundary reconstruction, so this does not establish that the
full synthesis path preserves those phones in context or explain the garbled
Stage 19 sentences.

The fresh downstream runs for P AH0 and T AH0 reproduced their prior WAVE
hashes. The IH0 rows above use the earlier valid captures
`trace-selection-valid-pih0.log` and `trace-selection-valid-tih0.log`, whose
WAVEs are 6,366 and 7,690 bytes. A separate rerun of both IH0 inputs with the
new generalized runner produced the same 22,084-byte output for both and
selected unrelated IDs; we excluded those two reruns from the table. Every
fresh run's Stage 5 input and output fixture hashes matched its saved
pre-run hashes after cleanup.

The earlier probe note treated `0x13` printed in the `FUN_100182e0`
prosody-feature view as P. That was an incorrect reading: the runtime-backed
CMU mapping identifies P as `0x35`, B as `0x13`, and T as `0x39`; the feature
view is not a direct dump of the requested phone symbol. The actual selected
index signatures, rather than those prosody-view bytes, support the alignment
finding above.

The selector still assigns substantial scores to these candidates (roughly
500–1,500 in the captured first context). The current analysis shows that
`FUN_100182e0` can add 500- and 1,000-point categorical penalties for signature
field mismatches, in addition to table-derived and duration terms. Thus,
correct central-phone alignment coexists with large candidate costs. That
keeps neighboring signature fields, class/tree mapping, and adapted feature
columns in scope; the captures do not yet identify which one is wrong. The
output hashes are P `01f4bdcaf34542e6a27fd65c352e3eaecf88085696268c5c0ba0e2216afa6fdd`,
T `e4bf1518d3c3d21420875a6a0c067b094834b161fc46511935708475d2d4a445`, and B
`61762aaf050695eef8779ac5f24ba0b4c88bb597162d07fb70e3b1221b417fed`.
The F output is `6f411ddb62590b4f7d2ab9bc73c75d6b73b61bd3d7bfcb5afcf04e0d629df6ae`.

A second GDB capture evaluated the scorer's two signature pointers directly
for forced P AH0. Its target row is `5a 5a 35 07 5a a0 00`; the first selected
Kate unit is `24 10 35 09 17 20 1d`. The central phone byte (`35`) matches,
but byte offset 1 is target `5a` versus unit `10`. At the runtime tables these
map to different category values (`0` versus `6`) and different lookup classes
(`5` versus `1`), so `FUN_100182e0` takes its 500-point fallback for that field.
That accounts for the 500-point portion of this unit's 512.6 local score. The
exact phone-symbol mapping is therefore not the main failure in this control;
the candidate's context/category code is being treated as incompatible by the
current scorer.

We also scanned all 283,696 adapted Kate signatures. Of 2,089 rows whose third
byte is P (`35`), only 25 have byte offset 1 equal to the target's `5a`; 32
have byte offset 3 equal to AH0 (`07`), and just three have both. None of
those three rows is among the ten units returned for this forced P control.
This suggests the current tree/class shortlist is not surfacing the rare rows
with that pair of expected context codes, or that the copied legacy bytes do
not mean the same thing to the current scorer. It does not yet distinguish
those explanations, and the remaining feature-column contribution to the
score has not been fully decomposed.

## Matched ordinary-text tree trace

The first old-versus-new forced-phone run was not comparable: the 2006 parser
receives `<vtml_phoneme ...>` as text and emits records for pieces of the tag
and its attributes. Replacing it with standard `<phoneme ...>` does not fix
that old-engine input path. Those 2006 WAVE files are therefore not phone
controls. The corrected differential uses the identical plain input
`Hello.` on both matched Kate engine/package pairs.

The 2006 parser consumes six bytes and emits one text record spanning
`Hello.`. For its first four scalar tree calls, the first eight printed feature
values match the 2013 calls, and the returned values also match exactly:
`1207`, `555`, `1740`, and `4832`. This is direct evidence that the current
tree2-to-tree3 bridge can preserve these scalar lookups for this input; it
does not establish broad scalar-tree equivalence.

The paired vector lookups diverge. At the first two matching contexts, the
old evaluator returns 12 shorts:

| Legacy lookup | Returned shorts |
| --- | --- |
| 1 | `85, 84, 83, 83, 84, 84, 85, 85, 86, 86, 86, 86` |
| 2 | `96, 98, 100, 104, 108, 112, 116, 120, 122, 124, 119, 113` |

The corresponding 2013 calls return 12-short arrays with only the first
value populated: `20492, 0, ...` and `24588, 0, ...`. The old GDB trace shows
the called tree pointers equal the engine's loaded `nbf` and `sbf` slots. Static
loader setup in `FUN_10001000` maps those slots to `pitch/nbf.tree2` and
`pitch/sbf.tree2`; both old calls therefore map to those exact source files.
The old tree header's byte at offset 4 controls the copied vector width; the
runtime-loaded `nbf` and `sbf` trees both report 12. `FUN_100019b0` copies
exactly that many shorts from the selected terminal node. The trace script now
captures those 12 shorts only.

This also resolves the apparent Stage 19 root outputs. The source files begin
`0b 44 00 00 0c 50 ...` (`nbf`) and `02 44 00 00 0c 60 ...` (`sbf`). The
2006 special loader `FUN_10001c80` parses these as tree records; the Stage 19
`tree2.py` parser instead reads bytes 4–5 as a scalar output, yielding
`0x500c` = 20492 and `0x600c` = 24588, then mistakes the next two bytes for a
leaf trailer and stops after 9 bytes. The “19,931/3,691 trailing bytes” are
therefore not evidence of an opaque suffix after a correctly parsed tree:
they are almost the entire old tree that this parser does not understand.
The converter then writes the bogus scalar as a one-output leaf, which
explains the populated first value and zero-filled remainder in the 2013
vector call.

This is a specific loader-format mismatch for the four files handled by the
old special loader (`nbf`, `bf`, `qbf`, and `sbf`), distinct from the other
recursive tree files and their scalar path. The evidence establishes a strong
compatibility cause for the tested vector lookups, but not that this alone
accounts for all sentence garbling. The 12 returned values' semantic labels
remain unassigned. An experimental parser now consumes each of the four files
to EOF and converts their leaf rows to width-12 indexed trees. On this same
plain-text input, the 2013 evaluator returned exactly the two 2006 vector rows;
candidate unit IDs changed at positions 3–5, confirming these outputs affect
the selection path. The generated WAV changed from 4,798 to 5,206 frames, but
the trial has not established intelligibility or compatibility across a
sentence set.

The plain `Hello.` captures are in the ignored local evidence directory at
`tools/revkit/work/corpus-parity/stage20/`. The old WAV has 7,610 frames
(0.476 s); the adapted 2013 WAV has 4,798 frames (0.300 s). These durations
record different synthesis output but do not measure intelligibility. The
parser, tree, and timeline logs preserve the input and returned values for
both engines. The [Stage 20 README](../../tools/revkit/work/stage20/README.md)
documents the reversible capture procedure and evidence filenames.

## What this changes for Lead 5

The leading unresolved cause remains the index/scoring bridge from Kate's
2005 unit data to the 2013 reader's 21-byte view, but the earlier `attr_b`
fill experiment was misinterpreted. A matched scorer trace shows zero-fill
removes candidate discrimination on that score dimension. The experimental
“signature-last-copy” variant changes costs and unit selection and the
requester described its `Hello.` output as closer to the intended word.
However, it copies the raw `[7N:8N]` window, which the native 2006 reader
identifies as `attr_40`. It does not copy each unit's last signature byte.
This is evidence that an injected legacy field affects selection and
listening, not evidence of the intended `attr_b` mapping.

The runtime reads each signature as seven bytes per unit. For Kate
`unit-etc`, local row 1286 has `02 5a 22 22 43 30 5a` at
`tail + N + 7 * 1286`, matching the adapted file and live scorer. A watchpoint
also shows initialization changes this row's first byte attribute from 0 to 3
at `0x1001a358`, consistent with adding the bank offset. The prior
byte-distribution comparison is a descriptive comparison of legacy
`attr_40` with Paul's `attr_b`; its EMD result does not establish that those
fields are semantically paired. The old transfer variant zeros the bulk
`attr_40` bytes in the malformed signature view; it does not clear one final
byte per unit.

This is a compatibility-layer diagnosis, not evidence that the 2013 synthesis
engine is inherently unable to use 2005 voice data. The 2006 engine loads the
native tree2 and ver.2005 formats; the tested 2013 path uses converted indexed
trees and a versioned index adapter. In the controlled `Hello.` path, the
first four scalar tree returns match, and an offline evaluation confirms that
the corrected vector adapter follows the same branches and returns the same
`nbf` and `sbf` leaf rows. This removes the known vector-tree defect as the
remaining explanation for this input. The 2006 engine selects four consecutive
`unit-gen2` rows; the corrected 2013 copy candidate selects eight units across
three banks. The source DAT and UPM spans decode and match the legacy records,
so the clearest remaining divergence is which units the adapted scorer selects
and how it segments them.

Several lower-level causes have been narrowed. All Kate DAT records match the
independent decoder at the PCM boundary; selected timeline offsets and sample
counts match the original legacy index records; and the requester reports
that the four selected-side composites sound like the intended phones. The
MSI embedded dictionary also has a confirmed phone-ID mismatch with a
measurable effect, but remapping it only makes the Stage 19 address output
match the current-data control, which is itself unintelligible. The old
vector-tree parser defect is now a corrected compatibility bug rather than
the remaining explanation. Its experimental 2013 conversion changes
WAVs across the existing prose, numbers, and address fixtures, but the
controlled `Hello.` vector calls now return the same 12-value rows in both
engines. The requester listened to the corrected-vector `Hello.` output with
zero-filled `attr_b` and reported that it sounded like “Shar” or similar,
rather than intelligible “Hello.” The requester reports that the selected-side
composites sound like intended phones, which points the remaining gap toward unit choice,
sequence, or timing/coarticulation.

A controlled `Hello.` comparison then held the corrected vector trees fixed
and changed only the experimental fill for the inserted 2013 `attr_b` from
zero to a copy of the legacy `[7N:8N]` byte window. The 12-value `nbf`/`sbf` rows and
first four scalar outputs stayed identical, while most of the first eight
selected unit IDs changed. Thus the vector-tree and index-field experiments
affect separate stages of the observed path. `signature-last-copy` remains a
candidate only: its field meaning and linguistic quality are unproven. The
next tests should recover the legacy reader's actual field mapping and compare old
and new unit selection for the same target contexts. See the
[Stage 20 record](../../tools/revkit/work/stage20/README.md) for exact
matrices, captures, and hashes.

The requester then listened to the same `Hello.` with the copy-window variant
and described it as sounding like “Hello” with a major speech impediment,
approximately “helall.” This is a subjective improvement over “Shar,” but it
does not validate a field mapping. Two controls reinforce that distinction:
`signature-last-transfer` (copy the same `[7N:8N]` window to `attr_b`, then
zero that source window)
selected a different unit sequence, while `constant-8` produced a WAV
byte-identical to the zero-fill case. With corrected vector trees held
constant, the `signature-last-copy` run across prose, numbers, and address
also reproduced the prior candidate hashes exactly. The broader outputs have
not yet been listened to or judged for intelligibility.

The requester reports that the transfer variant sounds colder than the
copy candidate, with an unclear initial fragment and a final sound described
as “nall.” We treat those as listening descriptions, not phone labels. A
matched `Hello.` scorer trace now explains a mechanical difference: zero-fill
and copy preserve the same malformed seven-byte signature view, return the
same first class ID (`46204`), and provide identical first three
unit-candidate lists to final ranking; copy injects the `[7N:8N]` window and restores candidate-specific
`attr_b` pair costs where zero-fill gave every unit the same cost for a target.
These traces used the un-repacked adapter. Its transfer variant zeros old
`attr_40` bytes inside the malformed signature block and returns a different
first class (`23148`) with different final-ranking candidate counts. The listening comparison therefore aligns
with two separate observed effects: copy changes later unit scoring, while
transfer also alters the earlier class shortlist. This supports preserving
the original byte if the copy hypothesis is pursued, but still does not
recover the intended `attr_b` definition or explain the remaining distortion.

### Generation-specific class-record consumers

The decompiled class-record readers show a second compatibility boundary
after tree evaluation. In the 2006 DLL, `FUN_10013cf0` dispatches to
`FUN_10013a40` or `FUN_10013bc0` according to a format flag. Those readers
select an interval record using a `0x68`-byte stride and a row width taken
from the legacy context at `+0x76`. The matched runtime trace confirms the
downstream `FUN_100232d0` consumes four legacy rows at a `0x14` (20-byte)
stride and resolves them to IDs `272822`–`272825`.

The 2013 DLL has the parallel dispatcher `FUN_1001b200`, branching to
`FUN_1001af50` or `FUN_1001b0d0`. Its decompiled readers use a `0x64`-byte
interval stride and a row width from context offset `+0x72`; the alternate
reader also extracts a different final field. The 2013 handoff
`FUN_1002d1e0` builds a 24-byte local record, passes it through
`FUN_1002c120`, then calls `FUN_1002c8b0` and `FUN_1002bbd0`. The stable
2013 runtime trace records eight timeline rows for the same `Hello.` input.
These are distinct class-record and handoff contracts, not a shared byte
layout. The names and semantics of every copied field remain unresolved.

### Candidate-sequence construction and where the row counts can diverge

The caller paths show that segmentation is decided before the class-record
handoff. In the 2006 DLL, `FUN_10022f60` calls `FUN_1001e470` for the ordinary
selection path. That routine builds a sequence of candidate positions,
deduplicates candidate IDs, expands difficult positions through
`FUN_1001dfb0` and fallback `FUN_1001e110`, scores neighboring choices in
`FUN_1001cff0`, prunes with `FUN_1001e2f0`, and finally reads one selected
class record per output position through `FUN_10013cf0`. Its return value is
the constructed position count. The matched runtime capture independently
shows four populated 20-byte selected rows for `Hello.`.

The 2013 path has a corresponding but split pipeline. `FUN_1002c9b0` calls
`FUN_10024980`; that calls `FUN_10024680` to construct positions, then builds
costs and performs ranking/backtracking through `FUN_10023350`,
`FUN_10018c80`, `FUN_10024510`, and `FUN_10024900`. `FUN_10024680` invokes
`FUN_10024060` and fallback `FUN_100242a0` where the old path invokes its two
split helpers. After selecting positions, `FUN_10024980` calls `FUN_1001b200`
for each chosen ID. That record consumer is then called again by synthesis
handoff `FUN_1002d1e0`, which produces the 2013 timeline rows.

This identifies the next causal test precisely: compare the old and new
position lists and split-helper returns for the same initial record and same
scalar/vector tree outputs, then compare candidate IDs and scores before and
after pruning. The current captures do not yet expose those intermediate
lists side by side. Therefore the observed four old rows and seven/eight new
timeline rows are not yet proof that the split helpers themselves cause the
cardinality difference. A divergence in class-row interpretation, candidate
generation, score weights/pruning, or the later timeline handoff could each
produce it. Static code already establishes distinct class-row layouts and
selection pipelines; only the 2005 index/signature layout defect has been
demonstrated to corrupt concrete scoring inputs. The crossfade and gap
controls do not resolve this selector question.

### Split-helper trace and cause status

A paired runtime trace now rules out the simple threshold-window explanation
for the extra 2013 positions in this `Hello.` run. The 2006 selector returns
four positions. Its four `FUN_1001dfb0` shortlists contain class IDs and
16-bit table values `52164:(45), 52165:(2), 52166:(24), 52167:(4)` (sum 75),
`32401:(5), 32400:(4)` (sum 9), `25383:(5)` (sum 5), and `37302:(6)` (sum 6).
All four sums exceed the old `>2` test, and no old two-position fallback is
taken. The capture repeats the same four-position result for all three
output slots. See
`tools/revkit/work/corpus-parity/stage20/split-threshold-old-gdb.log`.

The matched 2013 run uses the recovered vector-tree overlay, the
`attr_48/key[0:3]/attr_40/key[3:5]` index repack, and `attr_b=0`. The
lower-intrusion trace returns seven positions. At split-helper slots 0, 2,
and 5, `FUN_10024060` returns false without calling `FUN_10023060`; the
helper first checks whether `FUN_10018770` produced any class IDs and exits
immediately when the count is zero. Each false result enters
`FUN_100242a0`, which constructs two half-position rows. Slot 4 is the
observed ordinary whole-position acceptance: one class ID (`9147`), table
value 25, sum 25, then true. Therefore, for the three recorded splits, the
`>9` metric test is not reached. The extra positions arise before score
threshold evaluation because the whole-position class lookup is empty.
Runtime evidence is in
`tools/revkit/work/corpus-parity/stage20/split-threshold-new-stable-gdb.log`
and the matching output/hash record in
`tools/revkit/work/corpus-parity/stage20/split-threshold-new-stable/`.

The detailed lookup trace shows the mechanism at the next layer. For ordinary
whole-position rows (`row[+0x04] == 0`), `FUN_10018770` generates context
variants and calls `FUN_10023f90`; that dispatcher uses `FUN_10023dc0` to
derive a five-byte key through `FUN_10016ea0`, then `FUN_10019450` performs an
exact binary search in the sorted five-byte class table. If no matching key
exists, the candidate count remains zero and `FUN_10024060` returns before
scoring. The fallback marks its two rows with `+0x04` values 1 and 2. Those
half-position lookups take the `FUN_10023e70` path, which can relax the
feature query and broaden class candidates; that helper repeatedly sums
class-table values and normally stops once the sum is at least 10. A
class-byte-8 special case in `FUN_10024060` accepts any positive sum. The
2013 threshold therefore governs candidate sufficiency after lookup, but it
does not explain a zero-result whole-position lookup.

### Exact class-key inventory and suffix conversion

A minimal `FUN_10018770` trace with the same input and phone-ID-remapped
control gives stable target signatures for slots 0, 2, 4, and 5. Their
derived class keys under `FUN_10016ea0` are:

| Slot | 2013 target signature | 2013 exact key | 2006 query key |
| ---: | --- | --- | --- |
| 0 | `5a 5a 22 17 2b a0 00` | `5a 22 17 a0 00` | `90 34 23 30 30` |
| 2 | `5a 22 17 2b 30 20 00` | `22 17 2b 20 00` | `34 23 43 00 30` |
| 4 | `22 17 2b 30 42 00 00` | `17 2b 2f 00 00` | `23 43 48 00 00` |
| 5 | `17 2b 30 5a 5a 45 00` | `2b 30 5a 45 00` | `43 48 90 04 04` |

The key suffixes use different construction rules. The 2006
`FUN_1001da80` path divides its final two key bytes by 10 and takes their
remainders, then issues table-directed variants. The 2013 `FUN_10016ea0`
path copies signature byte `+5`, masks signature byte `+6` with `0x20`, and
the scoring code interprets byte `+5` as two three-bit fields with flags.
This is a demonstrated producer/consumer difference, rather than an
unexplained missing class in the source inventory.

The stable prefixes also correct the earlier phone/category hypothesis. The
2013 prefixes are `90 34 23`, `34 23 43`, `23 43 47`, and `43 48 90`; the old
prefixes are `90 34 23`, `34 23 43`, `23 43 48`, and `43 48 90`. The `48` to
`47` change is the established third lookup-table transform in slot 4. Slot
5's stable prefix agrees through the terminal `90`. The earlier `43 48 42`
and `2b 30 42` values came from a higher-intrusion GDB capture that changed
later expansion and are excluded. Repeating the exact target trace with the
embedded phone-ID remap produced identical signatures and byte-identical
baseline audio (`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`).
The remap is not causal for these Hello target keys.

The source-member crosswalk now resolves the threshold question for slots 2
and 5. Pointer-correct runtime bytes from `*(model+0x98)+id*7`, combined with
the 2006 `FUN_1001b2d0` table transform, identify the source rows represented
by each old shortlist ID. For every traced old ID, the number of matching
source rows equals its observed `+0x8c` metric weight. The reproducible
crosswalk is `compare-legacy-new-shortlist-members.py` under `stage20/`.

For slot 5, old class 37302 has exactly six source rows, the same six rows as
the mapped 2013 class 18857. The old engine accepts their metric 6 under
`>2`; the 2013 engine rejects that same set at `>9` and takes the two-row
fallback. For slot 2, old classes 32401 and 32400 contain five and four rows;
their nine-row union is exactly the mapped 2013 class 13478. The old engine
accepts the total 9, while the 2013 `>9` gate rejects the identical nine units.
This proves that the strict population cutoff causes the observed slot-2 and
slot-5 fallback splits after lookup succeeds; equal totals are not a
coincidence or merely a count-level resemblance.

A runtime intervention then changed only those failed split-helper returns:
for a nonempty class, a measured sum above 2 was accepted, matching the old
`FUN_1001dfb0` sufficiency rule. The sums 9 and 6 passed, both half-position
fallbacks disappeared, and the 2013 builder returned four positions instead
of six. The 30- and 25-row controls stayed accepted without intervention.
This confirms that the cutoff by itself causes the two extra positions on
this matched `Hello.` index overlay. The controlled WAV hash is
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`,
different from the unmodified matched-context run; four positions therefore
do not mean that its chosen units or sound match the 2006 output. The trace
and fixture-restoration hashes are
`old-cutoff-experiment-v3-gdb.log` and
`old-cutoff-experiment-v3/adapted-forced-pah0-restore.txt` under
`tools/revkit/work/corpus-parity/stage20/`; the override script is
`trace-split-threshold-matched-context-old-cutoff.gdb` under `stage20/`.

The four-position override selects 2013 IDs `273369, 273370, 273371, 255282`,
while the native 2006 path selects `272822, 272823, 272824, 272825`. They map
to `gen2` rows `93374, 93375, 93376, 75287` and `92827, 92828, 92829, 92830`,
respectively. Every selected row belongs to the shared old/new member set for
its context. The 2013 engine therefore had the native engine's chosen rows
available in all four contexts, and the native engine had the 2013 choices
available. The restored position count with different selected IDs localizes
the remaining `Hello.` difference to later unit ranking or timing, rather
than missing source candidates for these picks.

A paired control moved the legacy signature byte at offset `+4` into the 2013
`attr_b` column and zeroed the original signature byte, while retaining the
matched-context suffix map, vector trees, and `>2` cutoff intervention. The
2013 key builder does not consume signature byte `+4`, and both runs retained
class populations `30, 9, 25, 6`. Both runs returned four positions and chose
`273369, 273370, 273371, 255282`; their WAVs are byte-identical
(`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`).

The transfer changes the first per-unit `FUN_100182e0` scores: candidate
`273369` changes from `9.6` to `0.0705882`, and legacy-selected candidate
`272822` changes from `9.6` to `0.00851064`. The legacy choice gets the
slightly lower local cost, yet the final sequence and waveform remain
unchanged. Thus this field reaches local scoring, but its placement by itself
does not explain the sequence difference. Paired captures and the reproducible
overlay builder are documented in the Stage 20 README and retained under
`tools/revkit/work/corpus-parity/stage20/matched-context-score-{zero,transfer}/`.

The next counterfactual isolates the remaining fields for the two competing
same-key rows in context 2: 2006-selected `272824` and 2013-selected `273371`.
With `attr_b` populated from legacy `attr_40`, their direct 2013 local costs
are `15.2893` and `1.06048`. Swapping only the three 16-bit distance codes
leaves those local costs unchanged; swapping only the six scalar feature
bytes changes them to `6.87265` and `9.4771`; swapping only `attr_b` changes
them to `9.4771` and `6.87265`; swapping both `attr_b` and the six scalar
features exchanges the original costs exactly (`1.06048` and `15.2893`).
For this target and pair, the local-score preference is therefore driven by
the six scalar features, with `attr_b` modifying their effect. The three
distance codes are consumed by the adjacent-context pass instead: swapping
them changes transition distances and changes the selected four-unit chain.

The combined `attr_b` plus scalar-feature swap reverses the two local costs,
but does not make the ranker choose either candidate for context 2. The full
25-unit context-2 set instead backtracks through `264072, 264073, 264074`
after `273369`. Under that intervention the direct old/new candidate edges
also carry large feature or categorical penalties, while mixed-generation
edges are cheaper. This establishes an interaction between local feature
costs, distance-code transitions, and the broader candidate set; exchanging
two rows' score fields is not sufficient to recover the native 2006 path.
Field-swap input captures are
`matched-context-unit-cost-fields-{metric-codes,features,attrb,scoring-fields}-swap-gdb.log`;
the combined-swap path trace is
`matched-context-scoring-fields-swap-gdb.log`, all under
`tools/revkit/work/corpus-parity/stage20/`. The disposable row-swap builder
is `swap-matched-context-metric-rows.py` under `stage20/`.

The other two contexts expose candidate-set differences. Slot 0's four old
classes contain 75 source rows total; the 30-row 2013 class 32928 is an exact
subset, leaving 45 rows in the old set only. Slot 4's old class 25383 contains
five rows, all included in the 25-row 2013 class 9147, which adds 20 rows.
Both generations accept those contexts, but the different candidate sets can
change later unit ranking. The exact shared and one-sided bank-local row lists
are printed by the crosswalk script.

The conversion must depend on the complete class context, not only on the
two legacy tail bytes. Under slot 2's transformed prefix `22 17 2b`, source
tails `00 1e` and `00 00` combine as target suffix `20 00`. Under slot 4's
prefix `17 2b 2f`, source tail `00 00` remains target suffix `00 00` for its
25-row class. A global rewrite of `00 00` would therefore break the slot-4
match while repairing slot 2.

A stricter experimental adapter conditions the three suffix remappings on
their matched transformed prefixes and leaves all other source rows
unchanged; slot 4's `00 00` suffix stays intact.
Its split-helper trace reports classes `32928:30`, `13478:9`, `9147:25`, and
`18857:6` in context order. The first and third are accepted; the second and
fourth are rejected and fall back. It returns six positions and produces the
same WAV hash as the broader suffix-map overlays. The member crosswalk shows
that this controlled repack gives exact old/new membership for slots 2 and 5,
while slots 0 and 4 retain the measured subset/superset differences. The
reproducible repack variant and trace are
`attr48-key3-attr40-key2-matched-context-tail-map` and
`trace-split-threshold-matched-context-tail-map.gdb` under `stage20/`.

The static call paths explain why. In the 2006 engine, `FUN_1001b2d0`
converts the five-byte key into a seven-byte feature record using phone and
category tables. `FUN_1001b3b0` splits the last two fields by 10, applies
different normalization to their decimal components, and emits one of two
directional feature views. `FUN_10011fa0` applies a relaxation level starting
at 8; `FUN_1001d8c0` lowers it while the query has no candidates or its
accumulated population is below 3, subject to special guards. `FUN_1001da80`
also tries context and category variants. This is a structured candidate
expansion policy, not a direct conversion of two suffix bytes.

The exact `FUN_1001b2d0` producer is more specific than that summary: output
bytes 0–2 are lookup-table results selected from key bytes 0–2, the category
byte at +1, and the decimal tens/remainder of byte +3. Output bytes 3–4 copy
the final two query bytes; bytes 5–6 are table lookups of query bytes 2 and 0.
`FUN_1001b3b0` then builds the two eight-byte relaxed comparison views. The
2006 suffix path therefore depends on both the first-three-byte context and
the table-selected category mode. It cannot be represented by a global
two-byte suffix rewrite.

For the measured `Hello.` input, the seven-byte feature-distance stage is
skipped. `FUN_1001dfb0` copies candidate IDs directly when the deduplicated
list has fewer than 31 entries; only lists of at least 31 enter
`FUN_1001d5f0`, which can call `FUN_1001d7c0` and its `FUN_1001d530`
distance-ranker. The four old-engine shortlists contain 4, 2, 1, and 1 IDs,
so the measured accepts use lookup expansion and accumulated 16-bit class
weights (`75`, `9`, `5`, and `6`) passing the old `>2` rule. A pointer-correct
trace read the stored seven-byte rows at `*(model+0x98)+id*7`; the rows are not
ranked on this short-list path, and their field semantics and source-unit
identity remain unverified. Evidence is in
`old-class-features-v5-gdb.log` and
`old-class-features-v5/original-forced-pah0-winedbg.log` under
`tools/revkit/work/corpus-parity/stage20/`.

In 2013, `FUN_10016ea0` instead forms an exact five-byte key by copying
signature byte `+5` and masking byte `+6`; `FUN_10016ef0` later interprets
byte `+5` as two three-bit feature fields and conditional flags. Ordinary
whole-position lookups therefore do not reproduce the old decimal relaxation
sequence. The measured conversion must preserve both context-conditioned
class membership and the generation-specific candidate-expansion policy.

The matched-context repack returns six positions, versus seven in the stable
baseline. The exact old cutoff intervention reduces it to four positions,
confirming that the two added positions come from the stricter 2013 threshold
on the same slot-2 and slot-5 source members. It still does not establish the
old unit sequence or intelligible output: slots 0 and 4 have different
candidate sets, and the experimental four-position WAV differs from the
unmodified matched-context WAV. General suffix conversion across the full
voice corpus and the downstream effects of the slot-0/slot-4 candidate-set
differences remain open.
Runtime captures and reproducible scripts are under
`tools/revkit/work/corpus-parity/stage20/` and `tools/revkit/work/stage20/`.

### Cause status and discriminating hypotheses

The compatibility failure is established at multiple interfaces: the 2005
index is not in the 2013 unit-major signature layout; four legacy pitch trees
are vector-valued but the generic adapter read them as scalar leaves; and the
2013 class-record consumer uses a different interval and row contract. The
index and vector errors have both been shown to change 2013 tree/scoring or
unit-selection results when corrected in controlled runs. Together these
explain why the 2006 engine can use its native package while an unadapted 2013
engine receives the wrong model inputs. They are not yet a complete
explanation of the remaining unnatural 2013 output: no adapted run has yet
matched the full 2006 candidate sequence and timing, and the 2013 `attr_b`
meaning is still unverified.

The remaining questions have distinct signatures in the decompiled call
paths:

#### Plain Hello: lookup succeeds, then the 2013 cutoff splits

A fresh pair used identical plain `Hello.` input with the original 2006
package and the matched-context 2013 repack. The old query keys were
`5a 22 17 1e 1e`, `22 17 2b 00 1e`, `17 2b 30 00 00`, and
`2b 30 5a 04 04`; their metric sums were 75, 9, 5, and 6. The 2013 tree
leaves transform to keys present in the repack:

| Context | 2013 raw tree signature | Transformed key | Class weight | Result |
| ---: | --- | --- | ---: | --- |
| 0 | `5a 5a 22 17 2b a0 00` | `5a 22 17 a0 00` | 30 | accept |
| 1 | `5a 22 17 2b 30 20 00` | `22 17 2b 20 00` | 9 | split |
| 2 | `22 17 2b 30 5a 00 00` | `17 2b 2f 00 00` | 25 | accept |
| 3 | `17 2b 30 5a 5a 45 00` | `2b 30 5a 45 00` | 6 | split |

The first three key components align after the 2013 tables, including the
decimal `48` to `47` conversion in context 2. The last two key bytes are
encoded differently: in these four pairs the old suffixes `1e 1e`, `00 1e`,
`00 00`, and `04 04` correspond to 2013 suffixes `a0 00`, `20 00`, `00 00`,
and `45 00`. These are observed paired contexts, not a general suffix
conversion rule.

On this 2013 run the candidate lookup succeeded for every context. Its strict
`>9` cutoff accepted weights 30 and 25 and rejected weights 9 and 6, producing
six positions instead of the old four. This directly identifies the cutoff
as the cause of the two additional splits for this matched `Hello.` context.
The old and new metric tables are generation-specific, so their sums are not
treated as a shared numeric scale. Both WAVs contain 15,264 bytes but differ:
the old SHA-256 is
`eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`, and the
2013 SHA-256 is
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`. The
old and new query traces and restoration records are under
`corpus-parity/stage20/hello-old-query-classes/` and
`hello-new-query-classes/`; their GDB scripts are
`trace-hello-old-query-classes.gdb` and
`trace-hello-target-signatures.gdb` under `stage20/`.

An additional hardware-breakpoint trace records the 2013 tree-builder input
and its emitted query rows for the same `Hello.` bytes. The builder receives
position record words `1, 1, 2, 0, 0, 0, 0, 1`. It emits four bank-0 rows:
leaves 0 and 1 use feature index 0; leaves 2 and 3 use feature index 1. The
leaf signatures are respectively `5a 5a 22 17 2b a0 00`,
`5a 22 17 2b 30 20 00`, `22 17 2b 30 5a 00 00`, and
`17 2b 30 5a 5a 45 00`, identical to the prior target-signature capture.
This verifies the `FUN_10024680` row/leaf handoff for this natural-text case;
the record-word meanings and the upstream producer of that record remain
unlabeled.

This narrows the P AH0 mismatch without equating unlike inputs. The P AH0
trace starts from bank-0 leaves 0 and 1 whose signatures contain phone bytes
`35` and `07`; its derived keys miss the current adapted index. The `Hello.`
trace shows that the same builder path can emit rows whose keys are present in
the repack, after which the measured cutoff causes two fallbacks. Thus the
builder-to-class lookup is not universally broken. The forced VTML P AH0
fixture is not a valid 2006 comparison because the old parser treats that tag
as text; the trace does not yet identify which 2005 index columns should
represent the P AH0 tree fields.

The capture is `hello-phone-tree-provenance-gdb.log`, produced by
`trace-hello-phone-tree-provenance.gdb`. The hardware-breakpoint run exited
normally, reproduced the baseline 2013 WAV hash above, and restored the
fixture input/output hashes.

#### Natural `Apple.` query comparison

To compare the forced P AH0 case without VTML, both engines were run on the
same plain `Apple.` bytes. The 2006 query keys expose the phone sequence
`05 47 07 2b`; the 2013 raw tree signatures contain that same sequence across
leaves 0–3. This makes the input comparable at the observed phone-byte level,
though it does not prove the dictionary payloads are identical.

| Old query slot / key | 2013 query slot / raw signature | `FUN_10016ea0` key | 2013 lookup result |
| --- | --- | --- | --- |
| 0 / `5a 05 47 30 30` | 0 / `5a 5a 05 47 07 a0 00` | `5a 05 35 a0 00` | no candidates; split |
| 2 / `05 47 07 00 00` | 2 / `5a 05 47 07 2b 00 00` | `04 47 07 00 00` | class 1257, weight 50; accept |
| 3 / `47 07 2b 00 04` | 3 / `05 47 07 2b 5a 05 00` | `47 07 2b 05 00` | no candidates; split |
| 4 / `07 2b 5a 04 04` | 5 / `47 07 2b 5a 5a 45 00` | `07 2b 5a 45 00` | no candidates; split |

The 2013 table transform changes `47` to `35` in one key field and `05` to
`04` in another; it leaves the corresponding later-context prefix bytes
unchanged. The suffix pairs also vary by context (`30 30` to `a0 00`,
`00 00` to `00 00`, `00 04` to `05 00`, and `04 04` to `45 00`). The active
matched-context repack contains 50 members for the second transformed key and
none for the other three. Thus three of the four 2013 fallbacks occur before
the `>9` sufficiency check; only the populated class is accepted. The 2006
runtime accepts the last three old contexts and splits the first. This is
direct evidence that query-key conversion/class membership explains these
extra 2013 splits for this ordinary-text control. It does not yet say whether
each transformed field represents a phone, a context category, or another
selector feature, and the suffix conversion remains context-dependent.

The original WAV is 14,340 bytes and the adapted WAV is 13,648 bytes, with
SHA-256
`1ccabfc08b9ddbdf8cda2bd2a1f5d79b716649dbffe7c857572cb903bdbaa119` and
`48f9610abc82debe3e7ce41000fc88e082ad1a2d47602438df02ce132db1c40e`,
respectively. These hashes confirm distinct output, not intelligibility. The
old and new hardware-breakpoint traces and isolated fixture-restoration
records are under `corpus-parity/stage20/apple-plain/`; the 2013 split-helper
trace and restore proof are under `apple-plain-selector/`. The capture scripts
and one-line input fixture are under `stage20/`.

Several runtime-only index-repack interventions tested candidate explanations
for the missing classes. Reordering the legacy columns to
`attr_40/key[0:5]/attr_48` made the third Apple key resolve to a one-member
class with weight 1. The 2013 `>9` gate still rejected it, so the selector
remained at seven positions. With the broader generic suffix map, the
dedicated Apple trace found class weights 1, 50, and 158 at slots 0, 2, and
5; slot 3 remained empty. The weight-1 class still failed the 2013 `>9`
gate, while slots 2 and 5 passed, reducing the selected-position count from
7 to 6. This no-pitch-overlay output is 16,216 bytes with SHA-256
`a242ff700b17a4f2511786cee2a5bbd439b4dd2f744ca5988ee769c44b0178a7`.
These interventions prove that repack mapping changes lookup and split
behavior, and show that key membership and the 2013 sufficiency cutoff can
fail independently. They do not establish either candidate mapping as the
correct legacy-to-2013 conversion. Verified traces, WAV hashes, and
fixture-restoration records are in
`corpus-parity/stage20/apple-alt-column-order-selector-no-vector/` and
`apple-generic-tail-map-selector-no-vector/`.

A prefix-scoped candidate mapped source prefix `47 07 2b` with legacy tail
`00 04` to 2013 tail `05 00`. On the original no-pitch-overlay composition,
the live selector found weights 1, 50, 20, and 158 at Apple slots 0, 2, 3,
and 4. The 20-row slot-3 class passed the strict 2013 gate; the singleton at
slot 0 remained below it. The selector returned five positions, matching the
old engine's five-position result on three repeated calls. This output is
13,212 bytes with SHA-256
`167a388e30f74e83e6d955b39208fb9d62acd6063a3d833c12e0e12c6fef1e23`.
The new candidate recovers all 20 old slot-3 rows with no extras in the
source crosswalk. This pins the missing class for this context to a specific
legacy-to-2013 suffix conversion. The runtime trace and restoration record
are in `corpus-parity/stage20/apple-prefix-tail-map-selector-no-vector/`; the
old selector count is in `corpus-parity/stage20/apple-old-position-count/`.

The old runtime feature records were also crosswalked back through the 2005
index columns and compared with members of the matched 2013 repack. For
Apple's second context, the old class resolves to 48 source rows and the 2013
class contains 50: 43 rows are shared, five occur only in the old set, and
seven only in the 2013 set. The next two old classes resolve to 20 and 157
rows, while their corresponding 2013 query keys have no class members. These
counts match the old runtime's append counts for those three contexts, and
the Stage 19 copies of all five source indexes match the 2005 files byte for
byte. The comparison confirms that the current 2013 repack does not preserve
three ordinary-text lookup populations; one of those 2013 lookups already
passes its cutoff, so this discrepancy cannot be explained by threshold
policy alone. Among the tested repacks, the generic tail map recovers the old
slot-0 singleton row and 157 of 157 old slot-4 rows, with one extra new row;
the prefix-scoped candidate recovers all 20 old slot-3 rows. The generic tail
map still preserves only 43 of 48 rows for slot 2. These results support
`1e 1e -> a0 00`, `00 04 -> 05 00`, and `04 04 -> 45 00` as conversions for
these observed contexts, while leaving their scope unresolved. The old
trace prints ID 51804 twice for slot 0, while the source crosswalk resolves
one distinct row; that duplicate append representation remains unexplained.
`crosswalk-apple-class-members.py` reproduces the row-level comparison across
all candidate repacks; `trace-apple-old-class-features.gdb` captures the old
feature records, and `trace-apple-selector-decisions.gdb` captures the
verified 2013 intervention.

1. **How does suffix conversion generalize beyond these matched contexts?**
   Slot 5's source tail `04 04` maps to target `45 00` for the 2013 six-row
   class, while slot 0's `1e 1e` maps to `a0 00` for a 30-row class. These
   are controlled context matches, not yet a general conversion algorithm.
   Compare transformed contexts and member sets across more phones, stress
   marks, categories, and suffix combinations. The four `Hello.` contexts
   establish context-conditioned mappings, but not a general conversion
   algorithm for the whole Kate corpus.
2. **Why does the 2013 ranker choose a different path from 2006?**
   For the tested context-2 pair, six scalar feature bytes reverse the
   `FUN_100182e0` local-score preference when exchanged; `attr_b` changes the
   margin, and three 16-bit distance codes affect transition costs rather
   than that local score. Reversing the two local costs still leaves the
   ranker on a third chain because the 25-unit set has lower-cost mixed
   transitions. The 2006 path calls `FUN_1001c860` for local scoring and
   `FUN_1001cff0` for neighboring choices, with different feature tables and
   category rules. A GDB stop inside or immediately after the 2006 local
   scorer caused an access violation with software breakpoints. Hardware
   breakpoints now capture native 2006 `FUN_1001c860` returns for the selected
   IDs: `272822` = `0.00141844`, `272823` = `1.64526`, `272824` = `6.19391`,
   and `272825` = `3.9525`. The same-key 2013 pair scores `272824` at
   `15.2893` and `273371` at `1.06048` under the transferred-`attr_40`
   overlay. These raw values use different target feature records and scorer
   contracts, so their magnitudes are not a common scale; they show that the
   native and adapted scorers evaluate the same candidate row through
   different inputs and implementations. Hardware tracing shows
   `FUN_1001cff0` leaves one candidate at each old position, linked as
   `272822 → 272823 → 272824 → 272825`; the 2013 transition pass evaluates
   9, 25, and 6 candidates at positions 1–3 before final backtracking. The
   context corresponding to old class `25383` has five old source rows versus
   25 new rows, including 20 new-only rows. A hardware trace directly locates
   the old reduction: `FUN_1001cd60` returns candidate sets of 75, 9, 5, and
   6, with a coverage field at candidate `+0x10`; the first candidate in
   every set covers all four positions. `FUN_1001cff0` then detects that
   full-position span and retains only candidates with coverage 4. The
   survivors are exactly `272822`, `272823`, `272824`, and `272825`; shorter
   spans are discarded before `FUN_1001c860` local scoring. This is the
   concrete old-engine path on `Hello.`. The old target feature bytes for
   context 2 are
   `35, 101, 112, 2, 0, 0, 0, 0`; the paired 2013 target view is
   `35, 94, 101, 1, 204, 118, 102, 0`. Hardware-breakpoint captures are
   `legacy-unit-score-hardware-gdb.log` and
   `legacy-neighbor-links-hardware-gdb.log` under the Stage 20 evidence
   directory; the coverage rows are in
   `legacy-continuity-candidates-hardware-gdb.log`. Corresponding scripts are
   under `stage20/`. The instrumented 2006 process exited normally, its output hash
   `eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`
   matches the prior old-engine control, and fixture restoration hashes
   match. This fully explains the old chain for this fixture. A second
   old-only `Hi.` trace reports two positions, candidate pools of 10 and 9,
   and leading candidates `272820` and `272821` with coverage 2; the other
   listed candidates have coverage 1. This supports a gate relative to the
   runtime position count, rather than a fixed span of four. A paired 2013
   `Hi.` run and a phrase with no full-span candidate are still needed to
   measure how broadly this explains cross-generation selection. The second
   input capture and restoration proof are in
   `corpus-parity/stage20/hi-continuity-probe/`.
   The old selector bypasses
   `FUN_1001d530` for these four shortlists because each is below 31 entries;
   larger shortlists use that additional distance path.

   **Correction to the earlier comparison:** the 2013 path also has a
   full-span gate. `FUN_100230a0` computes the span at candidate `+0x0e`, and
   `FUN_10023350` filters to full-span candidates when the first candidate's
   span equals the position count. A renewed matched-context trace found that
   the four legacy-selected unit IDs were absent from all four 2013 candidate
   pools in that run. The earlier in-memory run that reported them as
   full-span survivors remains a separate recorded result, but was not
   reproduced under the complete versioned-engine composition used below.

   Static disassembly shows `FUN_100230a0` reads seven-byte candidate
   signatures and requires bit `0x80` in byte `+6` to extend the span. The
   2006 `FUN_1001cd60` uses a separate one-byte-per-unit array at model `+0x68`
   and stops when that byte is zero. For the four old rows, legacy `attr_48`
   and the captured old array bytes are `1, 1, 1, 0`; the adapted 2013
   signature byte `+6` is `0, 0, 0, 0`. Both generations' model `+0x30`
   lookups align the consecutive old IDs with their corresponding target
   class sequence, so those class IDs alone do not explain the missing span.

   The first disk-overlay attempt was invalid: it wrote byte 6 at
   `6 * unit_count + row`, but the adapted 2013 signature rows and selector
   use `row * 7 + 6`. A corrected disposable overlay sets the unit-major
   offsets for rows 92827–92829. A matched no-marker/marker run then verified
   that the loaded model bytes change from `00,00,00,00` to `80,80,80,00` at
   IDs 272822–272825. Both runs produced the same four candidate-pool sizes
   (211, 406, 78, 11), identical span distributions, and byte-identical
   7,704-byte WAVs (SHA-256
   `89c9c78a1326e4b7f4a5dd93f5c2fbf64795d01e89117bd4351e32d53ae47ff0`).
   No target unit appeared in any pool. Thus the corrected overlay reaches
   the selector, but this marker alone does not restore the legacy chain in
   the current matched-context path. The earlier recorded in-memory run
   reporting the chain `272822 → 272823 → 272824 → 272825` was not reproduced
   under this composition and should not be treated as the current
   disk-versus-memory A/B result.

   The span histograms were unchanged in every context: `178/33/0/0`,
   `376/30/0/0`, `76/2/0/0`, and `9/2/0/0` for spans 1–4. The target-unit
   mask and counts for their expected class IDs were zero throughout. This
   puts the current miss before continuity traversal: those unit rows never
   reach the candidate lists for the marker to extend.

   A runtime trace of `FUN_10018770` captured the P AH0 fixture's target
   signatures for whole-position slots 0 and 2 as
   `5a 5a 35 07 5a a0 00` and `5a 35 07 5a 5a 65 00`. Applying the
   `FUN_10016ea0` transform gives query keys `5a 35 07 a0 00` and
   `35 07 5a 65 00`. Neither key exists in the active matched-context
   repack. The four legacy-selected rows instead map to keys
   `5a 22 17 a0 00`, `22 17 2b 20 00`, `17 2b 2f 00 00`, and
   `2b 30 5a 45 00`. The query and row keys therefore diverge before
   continuity or scoring. A scan of 13 available repacks found the slot-0
   query key only in two context-tail variants (33 units each); slot 2 was
   absent from all 13. Even in those variants, the slot-0 class does not
   contain legacy unit 272822. This is direct evidence for an upstream
   target-key/class-membership mismatch on this fixture, not proof of a
   universal phone-ID mapping error.

   The target signature is selected from a tree leaf in `FUN_10024060`, then
   `FUN_10018770` derives lookup variants and `FUN_10016ea0` forms the exact
   five-byte key. The `0x80` marker is discarded by the key transform because
   byte `+6` is masked with `0x20`; it cannot repair a missing class key.
   This separates the tested causes: marker propagation now works, while the
   current blocker is the tree-produced target signature versus the repacked
   candidate-class signatures. Captures, full-pool GDB traces, and fixture
   restoration proofs are under
   `corpus-parity/stage20/continuity-bit7-{baseline,unit-major-overlay,unit-major-coverage}/`
   and `pah0-target-signatures/`; the comparison tool and trace scripts are
   under `stage20/`.

   `FUN_10016ea0` shows that the marker is outside the five-byte class key: it
   stores only `signature[+6] & 0x20` as the final key byte, discarding bit
   `0x80`. A corpus scan adds support for the mapping. All five Kate 2005
   `attr_48` columns contain only `00` and `01`; all Paul and James 2013
   signature byte `+6` values are `00` or `80`. Kate `unit-gen2` has 84,286
   marked rows out of 98,133 (85.9%). Native 2013 general-unit banks also
   carry high-bit markers on most rows, though percentages vary by voice and
   bank. These are unpaired voices, so frequency agreement does not prove a
   row mapping. Combined with the exact `1,1,1,0` legacy sequence and the
   successful in-memory `80,80,80,00` intervention, `attr_48 != 0` to
   signature byte `+6` bit `0x80` remains a plausible field conversion, but
   its effect is limited to candidate continuity and is not sufficient for
   this current selection miss. The full scan is
   reproducible with `compare-continuity-marker-distributions.py` under Stage
   20.

   **Controlled target-key substitution separates lookup failure from the
   cutoff.** A new runtime capture at `FUN_10024680` shows how the P AH0
   query is formed: the phone-tree walk writes bank 0 and leaf indexes 0 and 1
   into the per-position rows, and `FUN_10024060` reads those leaves as the
   target signatures already reported above. This confirms the mismatching
   keys are selected from the active tree output before index lookup; they
   are not caused by the continuity overlay.

   In an isolated GDB run, the first leaf's key fields were temporarily
   changed in memory to match the adapted class key of legacy row 272822. The
   exact lookup then produced class 32928 with weight 30 and `FUN_10024060`
   accepted the whole position. The selected stream began with unit 272822.
   This proves the original empty lookup is an immediate cause of the first
   fallback in this path.

   After that first position was accepted, the next query moved to position
   1. Temporarily changing its key fields to match legacy row 272823 produced
   class 13478 with weight 9. The unmodified 2013 metric returned 9 and the
   whole-position helper rejected it, then took the two-position fallback.
   A separate run held both key substitutions fixed and overrode only this
   metric return from 9 to 10. The helper then accepted the second whole
   position; the selected stream was exactly `272822, 272823`. This one-point
   intervention confirms the strict `>9` 2013 gate is a second, independent
   split cause once the key mismatch is repaired. The old `>2` gate would
   accept a score of 9 if both generations received comparable candidate
   metrics, though the score tables are not generally on a proven shared
   scale.

   The key-substitution run with the measured score 9 selected
   `272822, 282532, 282532, 272823` and produced a 7,888-byte WAV
   (SHA-256
   `61671203adaaa5d44667d751354467a1bf2afde16918c30d55d0e2afef029d7e`).
   The forced score-10 boundary run selected only `272822, 272823` and
   produced a 3,476-byte WAV (SHA-256
   `3d85d8842976c02644861cd5e9b19bb919307bd4388585eb942906eb36090df9`).
   These in-memory edits demonstrate selector mechanics; they are not valid
   model fixes and do not establish audio quality. Each runner restored the
   Stage 5 input and output files byte-for-byte. The logs and WAVs are under
   `corpus-parity/stage20/pah0-tree-leaf-provenance/`,
   `pah0-legacy-key-substitution-metrics/`, and
   `pah0-key-and-threshold-boundary/`; their GDB scripts are
   `trace-pah0-tree-leaf-provenance.gdb`,
   `trace-pah0-legacy-key-substitution-metrics.gdb`, and
   `trace-pah0-key-and-threshold-boundary.gdb` under `stage20/`.

   The concrete P AH0 cause is now a two-stage interaction: the unmodified
   tree leaves produce keys absent from the adapted candidate index, and an
   otherwise matching second class has population weight 9, exactly below the
   2013 acceptance boundary. The original key miss is upstream of scoring;
   the stricter cutoff becomes active only after the target key is aligned.
   The deeper source of the tree/index key disagreement remains to be traced
   through the old query producer, its context expansion, and the 2013 phone
   tree inputs.
3. **How much do the different class-member sets add beyond the coverage gate?**
   Slot 0 has 45 legacy-only rows, while slot 4 has 20 2013-only rows. The
   2006 full-position gate already excludes the observed third-chain candidates
   `273369`, `264072`, `264073`, and `264074` because their measured coverage
   is 3 rather than 4. In the renewed P AH0 trace, however, the legacy-selected
   IDs were absent from all 2013 candidate pools before the span gate. The
   measured candidate pools contained only span-1 and span-2 nodes, with no
   span-3 or span-4 node. Candidate membership therefore remains the earlier
   failure point for this path. The slot-4 superset gives 2013 still more
   alternatives, but the exact contribution of its 20 new-only rows remains
   unmeasured. A matched-set intervention can isolate that contribution
   after the target-class mismatch is addressed.
4. **The selected sequence is right but the handoff window/timing differs.**
   The old selected rows are 20 bytes and the new handoff has a separate
   layout; row spans and overlap arithmetic also differ. Force a matched
   sequence only after candidate selection is understood, then compare each
   selected row's source ID, DAT offset, sample count, leading/trailing span,
   and assembled PCM. A mismatch only at this boundary would localize the
   residual defect to synthesis handoff rather than candidate choice.

These tests are ordered by the earliest point where the current traces can
diverge. The reported “low/la” threshold from the derived gap controls is
useful for understanding adjacency in two already selected rows, but it is
not evidence for any of these engine-generation hypotheses.

The decompilation strengthens the format-compatibility diagnosis: the 2013
consumer has its own class-record contract and sequence-selection pipeline,
so copying legacy table bytes is not enough to demonstrate compatibility.
The Ghidra pseudocode references are
[`vt_eng-2006-all-functions-pseudocode.c`](../../tools/revkit/work/reports/vt_eng-2006-all-functions-pseudocode.c),
[`stage7-metric-builders.txt`](../../tools/revkit/work/reports/stage7-metric-builders.txt),
[`stage7-metric-offset-scan.txt`](../../tools/revkit/work/reports/stage7-metric-offset-scan.txt),
[`stage7-backtrack.txt`](../../tools/revkit/work/reports/stage7-backtrack.txt),
[`stage7-outer-ranking.txt`](../../tools/revkit/work/reports/stage7-outer-ranking.txt),
and [`vt_pau-synthesis-pseudocode.c`](../../tools/revkit/work/reports/vt_pau-synthesis-pseudocode.c).

### Index-layout reconstruction: primary compatibility cause

The native 2006 reader pseudocode at `FUN_100123a0` resolves the old tail
layout: `attr_4c[N]`, `attr_48[N]`, `key[5N]`, `attr_40[N]`, then three metric
groups totaling `12N`. The Stage 19 adapter inserted one `N`-byte column at
offset `8N` but left the old columns in place. The 2013 engine reads the
resulting `7N` region as unit-major seven-byte signatures, so it groups bytes
from different legacy columns and units together. This is a concrete adapter
bug that corrupts the signature used to classify and rank candidate units.

The 2013 key builder `FUN_10016ea0` takes signature positions 1, 2, 3, 5, and
6 for its five-byte key. A constrained reconstruction is therefore
`[attr_48, key[0:3], attr_40, key[3:5]]`: it places all five legacy key bytes
in the consumed positions, with the two adjacent old fields at positions 0
and 4. For `unit-etc` row 1286, the old fields are `attr_48=02`, key
`13 42 18 00 1e`, and `attr_40=21`, giving candidate signature
`02 13 42 18 21 00 1e`. This exact placement is supported by the reader and
builder pseudocode, but remains a candidate until checked against a paired
2005/2013 index or the legacy selector's live unit rows.

The targeted 2013 run with this repack and zero-filled `attr_b` preserved the
four traced scalar results and the recovered 12-value vector outputs, while
changing the class ID, unit shortlist, selected sequence, and WAV. The class
builder returned ID `9147` with key `17 2b 2f 00 00`; the first selected IDs
were `273369, 273369, 273370, 273370, 273371, 282025, 282025`, and the WAV
was 8,896 PCM frames (SHA-256
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`). The
2013 tree outputs stayed fixed, so this demonstrates that signature layout
alone changes selection. It does not establish that this is the final correct
mapping or that the new waveform is intelligible.

The earlier `[7N:8N]` fill was legacy `attr_40`, a genuine per-unit field;
the “signature-last” name was inaccurate. Copying it to `attr_b` while
preserving it in the malformed signature duplicates it. A corrected transfer
test put `attr_40` in `attr_b` and zeroed its reconstructed signature slot 4.
That kept the old key bytes in positions 1, 2, 3, 5, and 6, restored
candidate-specific `attr_b` costs, and changed the first selected IDs to
`273369, 273369, 272152, 272152, 272153, 282025, 282025` (8,132 frames;
SHA-256 `6c07819596129485710b732b67b8eeb314a680eafd6849509a00c6db6b902646`).
It does not prove whether `attr_40` should be moved, duplicated, or transformed
into the 2013-only independent `attr_b` field. The pair search found no
same-voice Kate 2013 model in the local inventory: Kate's indexes are all
`ver.2005`, and the available `ver.2013` indexes belong to James. A direct
row comparison then checked the four legacy-selected IDs (`272822`–`272825`)
against their raw `unit-gen2` fields. All four reconstructed signatures match
the candidate overlay byte-for-byte, while all four signatures in the
un-repacked adapter differ.

The forced-sequence control made the 2013 renderer use those four old IDs and
trimmed its seven-row timeline to the first four rows. It produced a valid
4,160-frame WAV (SHA-256
`d67dfd9c840d703421b869dd0b9c55bd0809a78b568145b6222b45d71cea9eac`). This
shows the 2013 path can render those inventory rows, while leaving 2006 versus
2013 timing and joining different. The requester heard it as “Hello,” though
the initial `h` was not well defined. This confirms a meaningful perceptual
improvement for the forced legacy sequence, while exact phone labels and
`attr_b` semantics remain open.

We then held the same IDs, repacked index, vector trees, and four-row timeline
fixed while changing only the first timeline row's sample count from 664 to
1,138 or 1,210. The legacy append trace shows the first 1,210-sample DAT unit
contributes 1,138 samples after a 72-sample overlap. These 2013 renders have
4,634 and 4,706 frames, respectively; the 1,138-sample WAV is
`153fc0465f8916f76226e840c4cf9dae482f96b09488ead6f22ffb6e59ea42a8`.
Its first 601 PCM samples match the 664-sample baseline exactly. The first
major energy rise shifts from the 30–40 ms windows in the baseline toward
60–70 ms, nearer the 2006 control's rise at 60–70 ms. This makes the initial
row-duration mismatch a concrete candidate for the weak initial consonant:
the 664-sample timeline reaches the next row substantially earlier than the
legacy 1,138-sample append. The requester reports the 1,138- and 1,210-sample
outputs sound identical to each other and both sound like “hello.” This
supports the forced sequence's intelligibility while suggesting the extra 72
samples between these two settings have little audible effect. It does not
prove a phone identity or equivalent crossfade.

A separate control retained all seven 2013 timeline rows while forcing the
first four legacy IDs and repeating `272825` in the last three slots. It
produced 13,616 frames (SHA-256
`6dbe7b3d858e6263d2dc073c43880edaacfcd876f903a32299886866a0f19349`). Its
first 1,000 PCM samples are byte-identical to the four-row control. Therefore
the four-row trim changes later audio, while the shared first four forced
rows contain the opening recognized as “Hello.” The repeated tail rows are an
artificial control and do not identify a natural 2013-to-2006 mapping. The
requester hears the result as “hello-oh-oh,” consistent with the same forced
ID appearing in each of the final three timeline rows.

A hybrid then retained the original 2013 IDs for rows 5–7 (`273371, 282025,
282025`) while forcing only the first four legacy IDs. It produced 10,604
frames (SHA-256
`6099c0d9d5e6891ac27b1b1a50c15a62348b388b5cc59ca25e2e5af02be529f5`); its
first 4,035 PCM samples match the repeated-tail capture exactly, and it
diverges after row 4. The requester describes it as “hello-glow,” with a
question-like continuation following conversational “hello.” This identifies
an audible difference after the four legacy rows, but the sound is not yet a
phonetic label for any selected unit.

We generated five-, six-, and seven-row cuts of this hybrid. These preserve
the first four forced legacy IDs and add, in order, `273371`, a `282025` row,
then its second `282025` occurrence. The requester hears “le—” in the
five-row cut, “low—” in the six-row cut, and a complete question-like “glow”
in the seven-row render. The seventh row is therefore needed to complete the
reported sound. Both of the last two rows select global ID `282025`, but
their timeline records differ: 1,886 versus 3,148 samples, type byte 1 versus
2, and leading/trailing spans 98/102 versus 102/98. The type-2 row starts at
an adjusted offset into the decoded unit. This shows the two selected rows
use different windows of that unit; it does not give the sound a recovered
phone or intonation label.

We also generated tail-only and last-row-only controls by compacting the
timeline rows after construction. They produced 6,570 and 3,148 frames,
respectively. The requester hears “low?” in rows 5–7 alone and “oh?” in row 7
alone. These are listening descriptions, not recovered phone labels. The
question-like impression is therefore present in the tail without the forced
legacy prefix; row 7 alone carries the reported “oh?” sound.

The source rows now explain the repeated-ID structure: row 5's `273371` maps
to `unit-gen2` local row 93,376 and decodes to 1,724 samples, matching its
timeline count. Rows 6 and 7's
`282025` map to `unit-etc` local row 2,900 and decode to 4,932 samples. Their
timeline counts total 5,034 samples, 102 more than that DAT payload, matching
row 7's 102-sample leading span. This supports a shared-window/overlap
interpretation for the repeated ID, but it does not establish the sound or
prosody encoded by that source row.

Separate row-5-only, row-6-only, and rows-5–6-only outputs are available for
listening. The requester hears row 5 alone as somewhat like “uh,” row 6 alone
as somewhat like “ah,” and the pair as “low.” These are listening impressions,
not recovered phone labels. The first row-6-only attempt used an incorrect
source offset and produced an empty WAV; it remains under
`natural-row6-only/` and is excluded. The corrected row-6-only output is under
`natural-row6-only-corrected/`.

PCM comparison of the three valid captures narrows the pair effect to the
join: its first 1,638 frames are identical to row 5 alone, and its final 1,800
frames are identical to row 6 alone starting at frame 86. The pair is 86
frames shorter than concatenating the isolated clips, exactly row 5's
trailing-span length. A derived hard-splice control preserves those same
prefix and suffix samples but substitutes row 6's unblended onset for the
86-frame engine blend. It is available at
`tools/revkit/work/corpus-parity/stage20/natural-rows5-6-no-crossfade/adapted-pah0.wav`
(3,524 frames; SHA-256
`864e2ad6f2832e885cc0e4ace7d39fb668336f2f8f6abe5893f2a84db85d39ed`).
The requester hears the hard-splice and blended pair as the same and considers
“la—” another plausible interpretation. The 86-frame blend is therefore not
required for the reported percept; the ordered sequence remains the stronger
lead, while its sound identity remains unresolved.

A reversed no-crossfade control is available at
`tools/revkit/work/corpus-parity/stage20/natural-rows6-5-reversed-no-crossfade/adapted-pah0.wav`
(3,524 frames; SHA-256
`a0fef9784064e899e1db4fab0b506ae130a914f3d04db81256dbfcec4e41ad34`). It
places the first 1,800 samples of row 6 before all 1,724 samples of row 5,
using the same 86-sample cut at the join as the forward hard-splice control.
This is a derived PCM control, not a native engine render. The reversed-order
control is heard by the requester as “ah-uh—,” consistent with an
order-dependent sequence percept. This is not a phone identification.

A forward-order gap control is available at
`tools/revkit/work/corpus-parity/stage20/natural-rows5-6-100ms-gap/adapted-pah0.wav`
(5,124 frames; SHA-256
`e4bd5c71234173e65bd8c6348f55d4a7477385457a871ec3f3258976c0cf9f56`). It
uses the same 1,638-sample row-5 prefix and all 1,886 row-6 samples as the
forward hard splice, with 1,600 zero samples (100 ms at 16 kHz) inserted
between them. This derived PCM control tests whether immediate adjacency is
needed for the “low/la” impression. The requester hears the 100-ms-gap clip
simply as “uh ah” and is skeptical of assigning it a word impression. The
20-ms and 50-ms versions sound the same in perceptual content, with only the
audible distance changing; possible word impressions remain too uncertain to
label. Across 20–100 ms, the pause separates the sounds without changing the
basic “uh ah” impression. This supports an adjacency-dependent percept for
the no-gap sequence, without establishing a word or phone identity.

Derived forward-order controls with 5-ms, 10-ms, and 15-ms gaps are also
available under `tools/revkit/work/corpus-parity/stage20/`. They preserve the
same row audio and insert silence only between the units. The requester hears
the fused “low/la” impression at 10 ms and below; the 20-, 50-, and 100-ms
versions remain separated as “uh ah,” with the audible spacing changing but
no stable word interpretation. This brackets that perceptual transition
between 10 and 20 ms. The 15-ms control remains available but has no separate
listening report. These derived controls address adjacency only; they do not
test the 2006-versus-2013 selector divergence.

A matched old-engine trace now reaches `FUN_100232d0` (`0x100232d0`) for the
same plain `Hello.` input. Its per-token table has four populated records,
each 20 bytes; the decompiler shows a 20-byte iteration stride and copies
fields from these rows plus companion metadata into synthesis arrays. The
parser consumes `Hello.` as one record, and the first four scalar tree returns
remain `1207, 555, 1740, 4832`. This identifies a useful legacy handoff for
field-by-field comparison, but the row fields are not semantically named and
cannot yet be mapped to the 2013 candidate IDs. It does not localize the
remaining divergence. Captures and raw row bytes are documented in the
[Stage 20 record](../../tools/revkit/work/stage20/README.md).

The trace was extended to `FUN_10013cf0`, the producer that fills those rows.
For this token it receives indices `272822, 272823, 272824, 272825` and returns
the four captured 20-byte rows in order. A runtime trace of the old engine's
unit manager confirms its first bank has 179,995 units and its next boundary
is `[179995, 278128)`, matching the `unit-gen2` range from `dblist.idx`. Thus
all four legacy selections are consecutive local rows 92,827–92,830 of
`unit-gen2`. The 2013 `signature-last-copy` path starts with IDs
`280411, 172786, 270092, 80917, 224948, 40384, 56352, 56352`, which resolve
to a mixture of `unit-etc`, `unit-gen`, and `unit-gen2`. For each sampled ID,
the 19-byte DAT/UPM span record in the adapted index exactly matches its
legacy source row, and every corresponding DAT payload decodes. The legacy
audio-preparation function receives those same four rows with sample counts
`1210, 578, 1180, 4854`. A trace of the old append function shows it emits
`1138, 506, 1112, 4854` samples: it subtracts `72, 72, 68` from the first
three rows and appends the last whole. The resulting 7,610 frames exactly
match the old `Hello.` WAV. The matched 2013 copy-path trace reports eight
timeline rows for `Hello.`. This establishes a concrete selection and
segmentation difference after the matching scalar lookups, while the exact
phone positions and the 2013 UPM/timing operations remain unresolved. It
does not establish that segmentation alone causes the distortion. Per-row
spans and trace details are in the [Stage 20 record](../../tools/revkit/work/stage20/README.md).

A fresh no-assembly-breakpoint control reproduced the same 2013 selected IDs,
eight timeline rows, and 10,650-byte `signature-last-copy` WAV (5,303 PCM
frames; SHA-256 `e86d005f2be73d6c93b9018d833cf38fb8bb3619763b17e60cf01de3b3adf1c8`).
The eight timeline `word3` values sum to 5,972 samples. Pseudocode for
`FUN_1002c8b0` copies that field's count times two bytes from each selected
segment into scratch audio. All rows have 100/100 timing-scale fields, so
`FUN_1002afb0` uses its direct `FUN_1002aac0` path. The 2006
`FUN_10023b40` and 2013 `FUN_1002aac0` routines each carry a row's trailing
span into the next row and advance by `count - trailing_span` for non-final
rows. The shared 8,192-entry coefficient tables are byte-identical between
the DLLs (SHA-256
`6a71e6a37aaa3157d6edaae839a14986da4701112887c2c0b689d9affdbff6a6`), but
equal table contents do not mean the join schedules match. In the 2006
`FUN_10023b40` disassembly (`0x10023be9`–`0x10023c09`), the two table indexes
for overlap sample `k` in a carry span of `m` are computed as
`floor((k + 1) * 8191 / (2*m + 1))` and
`floor((m + k + 1) * 8191 / (2*m + 1))`, where `m` is the pending carry span.
In the 2013 direct path used by this capture, the x86 loops in `FUN_1002aac0`
(`0x1002ab6f` onward) scale their loop counter in steps of 4096 over an active
span `m`, indexing the lower and upper halves at
`floor(k * 4096 / m)` and
`4096 + floor(k * 4096 / m)`. The two routines therefore use non-identical
integer schedules, with the 2013 routine also taking separate branches when
the pending carry and current leading span differ. The 2013 routine then adds
paired scratch-buffer samples 1,000 shorts apart into the output. These are
concrete implementation differences in overlap handling; their audible
contribution has not been isolated.

The 2013 direct routine computes the interior copy length as
`max(count - trailing_span - leading_span, 0)`, carries the computed trailing
span for non-final rows, and advances the PCM cursor by interior plus leading
span. On the final row it also emits the remaining tail. For the captured
rows this arithmetic predicts 5,310 frames after 662 samples of carried
overlap; the stable WAV contains 5,303 frames. The seven-frame residual
therefore lies beyond the current row-count/cursor accounting, or indicates
that one of those captured inputs differs from the values used in the stable
run. The old path emits 7,610 frames from 7,822 raw samples using directly
traced 72/72/68-sample trims. This shows both a segmentation/window difference
and a distinct overlap-weight schedule, but does not establish which causes
the audible defect. The last two 2013 rows both reference unit `56352`: their
1,166 and 956 timeline samples, less the shared 90-sample overlap, equal that
unit's 2,032 decoded DAT samples. This is an exact length consistency for the
repeated-unit case, not proof that the timeline rows partition the original
payload. The third
scale field on every new row is 200, and `FUN_1002c8b0` applies it as a sample
gain before assembly. The seven unique selected DAT payloads peak at
4,224–13,056, so the ×2 step does not saturate their source samples before
joining. The final new WAV's peak/RMS are 25,088/7,030, versus 14,336/3,945
for the old control; neither has full-scale samples. The output level
difference is not a controlled gain comparison because the unit sequences
differ. The 2006 join combines its weighted values and clamps the combined
result before conversion. In the 2013 direct path, the x86 code converts the
weighted lanes separately, then adds two 16-bit samples 1,000 shorts apart
and stores the 16-bit sum without a visible saturation check. That sum can
wrap if its inputs exceed the signed-16-bit range. Source peaks after the 2x
gain are below that range. The table's observed sine-squared shape suggests
the paired lower/upper weights nearly sum to one, which may bound the result
by the lane peaks; exact paired sums and the runtime lane values have not been
checked. The unchecked 16-bit addition is therefore a concrete trace target,
not evidence that this WAV overflowed. The row type byte
alternates 1/2; in `FUN_1002c8b0`, type 2 changes the decoded-buffer copy
start using the row's leading span, while type 1 starts at the buffer base.
That confirms the output is built from selected windows of decoded units. The
type semantics and effect of those windows on the perceived defect remain
open. The per-row arithmetic and captured span table are in the [Stage 20
record](../../tools/revkit/work/stage20/README.md).

Attempts to inspect `FUN_1002bd90` directly changed the result: a minimal
entry breakpoint produced 5,518 bytes, and a detailed entry/return probe
produced 16,134 bytes. Both differ from the stable no-breakpoint control,
while the Stage 5 input and output fixture hashes matched before and after.
The breakpoints therefore perturb this runtime path for an unresolved reason;
their captured post-entry row values are excluded from the behavioral
comparison. Tracing the 2013 join now needs a method that does not stop at
that function. This leaves the 2013 per-row UPM/timing contribution unresolved.

Stage 19 shows that the supplied standard DLL can traverse an adapted Kate
model and write non-silent audio. Its output is sensitive to the MSI embedded
dictionary family, while the MSI exception-table swap did not affect these
three fixtures. The original 2006 DLL also loads its own MSI dictionary and
voice package and successfully synthesizes audio. Taken together, these facts
show that the legacy package is usable by its original engine and that the
phone-ID incompatibility is directly visible in a matched word (`Hello`) and
has an observable synthesis consequence: ID conversion makes the standard
engine's address output byte-identical to the current-corpus control. This
confirms one concrete format incompatibility with a measurable audio effect.
It does not establish that the incompatibility explains the unintelligibility.
The remaining uncertainty is what explains the still-different prose and number
outputs, and whether these conversions improve human intelligibility. Since
the address remap only restores the current control's output, and that control
has also been reported unintelligible, the ID mismatch is not a complete fix.

The forced P trace establishes one concrete 500-point mismatch at signature
byte offset 1. The four downstream probes rule out a simple unit-ID-to-DAT
offset/span displacement for these examples: the timeline consumes the exact
first or second samples predicted by each original record. The remaining
question is whether those selected spans themselves contain the unexpected
sounds, whether context scoring selects poor examples, or whether the later
timing/coarticulation path changes their percept. The exact acoustic role of
the injected byte window and metric columns remains unresolved.
This remains a 2013 scorer-specific clue. The matched Stage 20 tree trace now
shows a separate concrete tree-output contract gap between the engine
generations. The selected-side listening result makes a simple wrong-half or
DAT-offset explanation less likely for the four tested controls, while the
effect of restoring the legacy leaf vectors and the remaining index semantics
on full-sentence intelligibility are still untested.
Separately, the 6,990
direct payloads that still differ after conversion, 106 alternative branches
with residual changes, four branch layouts the recovered parser does not
accept, 17,902 other-flag payloads, and old ID 184 remain unresolved.
`ENGBI_comp` and the old `atmt.tree2` layout also remain open. The absent class
indexes remain runtime facts not yet tied to intelligibility. The current
468-byte verification record causes the old engine's demo path, while the
matched MSI's 315-byte record does not.

## Comparison-only 2006 verification patch

The 2006 MSI DLL does not export the peer team's 2013
`VT_CheckLicense_ENG`; its `VT_GetTTSInfo_ENG` selector 1 and load path call
internal `FUN_10022430` instead. We made a hash-guarded copy of the MSI DLL and
left the extracted DLL and Kate resources byte-identical. The copy changes
one branch:

- At VA `0x10022542` (file offset `0x22542`), the final `jge` in
  `FUN_10022430` becomes an unconditional jump to its zero-return epilogue.

With the newer 468-byte workspace verification record, the unpatched DLL
followed the demo path and injected changing text. The exact MSI record is
315 bytes. With that file mounted read-only, the unpatched DLL returned `0`
from `VT_LOADTTS_ENG`, returned `1` for each synthesis call, and parsed only
the requested sentences. The GDB trace is retained at
[`unpatched-parser-trace.log`](../../tools/revkit/work/stage19/license-control/unpatched-parser-trace.log).
The matching API results and exact verification-file read are in
[`license-control`](../../tools/revkit/work/stage19/license-control/README.md);
the file trace read exactly 315 bytes from the MSI record.

The one-byte patched copy returned the same API statuses and produced
byte-identical WAVE files for all three fixtures. This shows that the patch
works as a forced-success comparison control, but is unnecessary when the
original engine is paired with its own MSI verification file. The earlier
status-only test that returned `-2` from synthesis used the newer record; it
does not show a helper patch is required with the MSI record. The patcher and
reproduction commands are in
[`license-control`](../../tools/revkit/work/stage19/license-control/README.md).
The generated DLL and WAVs remain in ignored local work storage.

A separate one-branch comparison copy bypasses only the nag-prefix insertion
at VA `0x10015957`, without changing the license checker. With the newer
468-byte record it produced the same three WAVE hashes as the valid MSI
record controls. The requester reports clean, intelligible output both with
this demo-text insertion bypassed and with the original insertion present.
This establishes intelligible output for the tested 2006 MSI engine and
resources in both verification paths; the license mismatch prepends speech
without demonstrated degradation of the requested text. It does not establish
intelligibility for the Stage 19 standard-engine `.tree2` adapter. The
unpatched MSI DLL hash remains
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.

## Full static reference dump of the 2006 DLL

The extracted 2006 DLL has now been run through the full local static toolkit.
GNU `objdump` and pefile captured the PE layout, complete 44-name export and
import tables, and raw-byte x86 disassembly. Ghidra 12.1.4 auto-analysis
discovered 947 functions; the all-function pseudocode export completed for all
947 with zero decompilation failures. The bundle also contains `file`,
`strings`, FLOSS 3.1.1, and capa 9.4.0 output. The input SHA-256 was checked
before and after and remains
`00fc9375d08bd8cec303845c992481d79c8f616d1bfb6390a5322cd06f69273d`.

These are analysis results, not recovered source: Ghidra's function discovery
can miss code or misclassify data, and decompiler pseudocode is approximate.
FLOSS and capa results are heuristic triage. The large disassembly, pseudocode,
and Ghidra project remain ignored local artifacts under
[`full-analysis-2006`](../../tools/revkit/work/corpus-parity/full-analysis-2006/README.md).

## Natural Hello: continuity conversion and split threshold are separate

A matched natural `Hello.` experiment now tests the continuity-marker mapping
on a path where the target classes are present. In the no-marker-to-marker A/B,
both variants use the matched-context `attr_b`-transfer index; the only
additional index change was setting bit `0x80` in signature byte `+6` for
`unit-gen2` local rows 92827–92829 from their legacy `attr_48` values. With the
2013 native cutoff, a low-intrusion hardware trace changed the selected source
rows from `273369, 273370, 273370, 273371, 264074, 264074` to
`272822, 272823, 272823, 272824, 272825, 272825`. This shows the marker
conversion moves selection onto the four legacy rows, while the 2013 `>9`
gate still splits positions and repeats handoffs.

To isolate the gate from the marker, both A/B runs then used the same runtime
override that accepted a nonempty whole-position metric above 2. The metric
returns were 9 at slot 1 and 6 at slot 3 in both runs. The no-marker control
selected `273369, 273370, 273371, 255282`; the marker overlay selected exactly
`272822, 272823, 272824, 272825`. This establishes, for this natural `Hello.`
case, that mapping legacy `attr_48` into the 2013 continuity bit fixes the
source-row chain when the cutoff is equalized. The native 2013 cutoff remains
an independent reason for split/repeated positions after that mapping.

The four matched WAVE captures and logs are retained under
`tools/revkit/work/corpus-parity/stage20/hello-continuity-bit7/`. With the
native cutoff, no-marker and marker WAVE hashes are
`8d6b4f63906ad339f5b9b4dde977a495d2c7dfc7a7c9a0818e2a6237d81b3dfb` and
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`.
With the cutoff override, they are
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318` and
`215ff65cf15e37f4dcce58527511aa31f143b58a1dca865ef082252a3a83e5f3`.
All runs exited normally and their restore records show the Stage 5 input and
output fixture hashes were unchanged by each run. More heavily instrumented
traces in this evidence set varied and are excluded from the behavioral
comparison. The result demonstrates two selector mechanisms on one word; it
does not show that these two changes account for the full-sentence defect or
generalize to other contexts. The full method and capture filenames are in
the [Stage 20 record](../../tools/revkit/work/stage20/README.md).

### Apple full-bank marker conversion

A second natural-text control applied the `attr_48` to signature-byte-6 bit
`0x80` conversion across all five Kate indexes. The generated overlay validates
all 283,696 rows and sets the bit on 240,417, using the matched-context
`attr_b`-transfer index as its base. The per-bank hashes and row counts are in
the overlay manifest under
`tools/revkit/work/stage20/index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks/`.

The old `Apple.` engine trace shows five unit IDs repeated for each of three
synthesis slots: `177774, 177774, 100180, 59558, 59559`. With no markers, the
2013 selected-unit handoff trace records nine IDs; with the full-bank overlay,
it records thirteen. The first three IDs are unchanged, the later sequence
changes, and none of the new IDs matches the old five-unit sequence. The
WAVs also differ (12,354 versus 13,178 bytes), but that does not establish an
intelligibility change.

This extends the direct evidence that the marker conversion affects
continuity traversal beyond the single Hello chain. It also shows that a
full-bank conversion alone does not restore old Apple selection. Since the
class-key builder discards bit `0x80`, the conversion cannot fill the three
empty Apple whole-position classes already observed in the no-marker trace.
That leaves query-key conversion and candidate membership upstream of
traversal, and later ranking among the surviving candidates, as separate
compatibility questions. Exact captures, hashes, and restore records are in
the [Stage 20 record](../../tools/revkit/work/stage20/README.md).

The selected-row crosswalk also separates missing classes from ranking. The
19-byte index records for old selected IDs `177774`, `100180`, `59558`, and
`59559` are unchanged in the adapted `unit-gen` index. Their transformed
row keys are respectively `5a 05 35 1e 00`, `04 47 07 00 00`,
`47 07 2b 00 00`, and `07 2b 5a 04 00`. Compared with the corresponding
Apple 2013 target keys (`5a 05 35 a0 00`, `04 47 07 00 00`,
`47 07 2b 05 00`, `07 2b 5a 45 00`), only ID `100180` has an exact key
match. Its class has 50 members, but 2013 selects other IDs. Three old rows
therefore fail at exact class-key lookup; the matched row remains a concrete
ranking or transition-cost case. The continuity conversion cannot address the
three key suffix differences because bit `0x80` is discarded during key
construction. Per-row bytes and context details are in the Stage 20 record.

### Apple transition cost for the surviving matched row

A conditional runtime trace compared old selected unit `100180` with 2013
candidate `128863` at Apple context 2 against predecessor `142230`. Weight
row 2 applies coefficients `10`, `2`, and `5` to the raw/A/B distances;
divisor is 1, categorical penalty is zero, and the scorer adds fixed `2.0`.
`100180` has raw/A/B distances `0.429843`, `0.115108`, `0.116129`; `128863`
has `0`, `0.26087`, `0.116129`. Their duration/prosody terms are `0` and
`0.0592593`, yielding edge totals `604.226` and `600.279` from the same
predecessor cost `597.117`. The edge to `100180` is about `3.94765` more
expensive.

Separate old and new unit-scorer traces both give `100180` score `0` and
`128863` score `0.118519` for target row 2. Both local rankers therefore
prefer `100180`. The old neighbor trace retains both in its position-2 set
and links them to predecessor `177774`; old synthesis selects `100180`. The
old and new position-2 candidate sets share 29 of 30 IDs. The 2013 path trace
links `128863` to `142230`, and its transition metric group reads code `744`
for `100180` versus `884` for `128863`, while the predecessor's group-2 code
is `884`. This pinpoints the immediate 2013 reversal to path scoring on
shared candidates rather than local ranker preference or class membership.

The component arithmetic reproduced across captures, but the heaviest GDB
probe changed the overall backtrace and WAVE, so only its individual score
components are treated as stable. A disposable one-word intervention changed
`100180`'s group-0 code from `744` to `884`, preserving its class key and all
other row bytes. The scorer read the replacement: the raw distance became
zero, and the edge to `100180` became about `0.3515` cheaper than the edge to
`128863` from the same predecessor. The intervention did not restore the
global path: paired low-intrusion runs had identical handoff IDs and WAVE
hashes, and the detailed intervention backtrace chose another context-2
candidate (`64728`). The raw code mismatch contributes to the local edge
ordering, but it is not sufficient to explain or repair the complete Apple
path. Other context paths and transition features remain active hypotheses.
The overlay builder, Compose override, component capture, handoffs, and restore
hashes are documented in the Stage 20 record.

An offline crosswalk resolves whether this is caused by metric-code repacking.
For all 31 distinct IDs in the old and 2013 context-2 candidate lists, all
three 16-bit metric codes in the adapted index match the corresponding 2005
source-index codes. The index adapter preserved them. The old neighbor trace
links all 29 shared candidates to predecessor `177774`; the captured 2013
transition list links those rows to five different predecessors. On the shared
29 rows, recomputed raw-table distances have four zeros under old links and
ten under 2013 links. For `100180` and `128863`, the raw distances against
`177774` are `0.429843` and `0`; in the 2013 candidate trace their respective
predecessors are `188825` and `64727`, producing `0` and `0.450601`. The
current-row codes remain unchanged, while predecessor selection changes the
matrix lookup. These predecessor IDs are specific to that capture; a separate
component capture linked both rows to `142230`.
These are offline lookups using the 2013 triangular table and captured links,
not transition costs measured inside the 2006 engine.

The metric-column mismatch hypothesis is therefore rejected for these
candidates. Different predecessor links are directly observed and can reverse
the raw-term preference, but they do not yet explain the global path or the
audible difference. The raw-term-only crosswalk is reproducible with
`tools/revkit/scripts/compare_apple_metric_edges.py`; implementation details
and capture boundaries are in the Stage 20 record.

The 2006 candidate builder supplies a concrete algorithmic difference for
this shortlist size. Ghidra pseudocode for `FUN_1001dfb0` copies a deduplicated
candidate pool directly when it has fewer than 31 rows. At 31 or more, it
reduces the pool to 30 with `FUN_1001d5f0` and the weighted signature-byte
comparator `FUN_1001d530`. The old Apple context-2 trace has 30 rows, so it
takes the copy-all branch. The 2013 `FUN_10018c80` path instead scores
adjacent candidates with the raw and derived `cepdist` terms. This explains
why identical metric-code bytes do not produce an equivalent old/new path
objective at this position; it does not by itself account for key misses or
the full global selection.

A controlled 2-by-2-by-2 weight ablation on the context-2 call to
`FUN_10018c80` further tests whether one table term explains the global
selection. The unmodified raw/A/B coefficients are `10/2/5`; each run restored
them at function return. Zeroing raw alone, derived A alone, derived B alone,
or any two together left the same 11 selected handoffs and identical
14,766-byte WAVE (SHA-256
`75da51d1f964aa02bfc05bfb57a4775e3c091f02b452c5de6ab436d16e7223f3`). Only
zeroing all three together changed the path to 13 handoffs and a 15,416-byte
WAVE (SHA-256
`159cf7612470a23845b96986e43c48b463af83f36a8a6cf7ec3b0a7f15ba7673`). The
intervention trace confirms the coefficients were `0/0/0` during the selected
transition call and restored to `10/2/5` afterward; all runs exited normally
and restored the Stage 5 fixtures.

For this Apple input, the raw distance by itself does not control the global
path. The weighted table terms affect selection jointly, consistent with
their combined subtotal competing against other path costs. These ablations
do not yet isolate the full-score margin or explain the remaining class-key
mismatches; the adjacent-row continuity gate is tested below. The scripts,
output hashes, and restoration records
are documented in the Stage 20 record.

### Apple adjacent-row continuity gate

The `FUN_10018c80` disassembly and pseudocode read the `+6 & 0x80` marker from
the previous unit's signature. When the current unit ID is the previous ID
plus one, a set predecessor marker makes the combined raw/A/B table subtotal
zero. In the stable Apple path, `64729 <- 64728` at context 3 has measured
raw/A/B distances `0.648320`, `1.30502`, and `0.318471`, and a duration/prosody
term of `75.9138`; the table subtotal is nevertheless zero while row `64728`
has the marker.

A disposable one-byte intervention cleared that marker on row `64728`. The
path changed from 11 handoffs and a 14,766-byte WAVE (SHA-256
`75da51d1f964aa02bfc05bfb57a4775e3c091f02b452c5de6ab436d16e7223f3`) to 10
handoffs and a 14,724-byte WAVE (SHA-256
`5352979756a9f9fdd0923f5517864e9f054fa4dac67537428d015428de0e7fc4`). In the
treatment, row `64729` no longer appears in scorer calls across all contexts,
even when tracing all predecessors, so the nonzero counterfactual edge cost
was not directly observed. The paired result establishes that the predecessor
marker participates in path selection, while the source condition and control
trace establish the zero-subtotal shortcut itself. Both runs restored and
verified the Stage 5 fixtures. Detailed captures and reproduction steps are
in the Stage 20 record.

Entry dumps of the context-3 list show the candidate change occurs before
transition scoring. Control and marker-cleared runs each pass 30 rows to
`FUN_10018c80` and share 29 IDs: the control includes `64729` and omits
`28631`, while the treatment includes `28631` and omits `64729`. The caller
`FUN_100248b0` invokes candidate builder `FUN_10023350` before the transition
scorer. That builder calls `FUN_100230a0` for continuity metadata, then scores
and caps lists above 30. The helper's disassembly tests signature byte `+6`
bit `0x80` during adjacent-context candidate checks at `0x1002316e` and
`0x10023216`. This ties the controlled marker change to a continuity-aware
shortlist substitution before edge scoring. The treatment edge cost for
`64729` remains unmeasured because that row is excluded from the context-3
shortlist. Candidate-array dumps are in the Stage 20 record.

The paired pre-cutoff trace identifies how the candidate swap occurs. Before
`FUN_100230a0`, row `64729` is at candidate-list index 68 with metadata
`+0xc,+0xe,+0x10 = 0,1,0` in both runs. After the helper, it is at index 2
with `1,4,0` in control and at index 165 with `1,2,0` in treatment; row
`28631` stays at index 32 with `2,3,0`. The score stored for `64729` is
`1039.64` immediately after the helper in both runs.

The following local-score pass handles a 47-row prefix in control, where
`64729` is included at index 2, and a 46-row prefix in treatment, where index
165 falls outside it. Row `28631` is included in both prefixes. Control calls
`FUN_100182e0` for `64729`; it returns a categorical penalty of `100` plus
`51.8277` weighted distance, replacing the stored `1039.64` with `151.828`.
Treatment does not call `FUN_100182e0` for `64729`, so `1039.64` remains.
Therefore the observed `+887.812` is a stale-versus-recomputed score
difference caused by marker-dependent ordering and score-pass inclusion. It
does not show a change to the local ranker's categorical or weighted-distance
terms. The treatment's hypothetical local score for `64729` remains unknown.
As a mediation check, forcing only its computed `+0xe` field from 2 to 4 at
the helper's write restored index 2, restored the 47-row scoring prefix, and
reproduced the control WAV byte-for-byte while the source marker stayed clear.
This establishes that `+0xe` is sufficient for the observed reorder and
rescoring change on this fixture.
The paired logs and fixture restore records are in the Stage 20 evidence
directory.

### Candidate coverage fields differ by generation

Decompilation resolves the field roles in the continuity helpers and their
consumers. The 2006 `FUN_1001cd60` stores left-match count at candidate `+0x0c`,
right-match count at `+0x0e`, and total span at `+0x10`. Its sort key rewards
twice the total span, a one-point bilateral-coverage bonus, and a two-point
direction bonus from the current six-byte context record. The old caller
`FUN_1001cff0` derives a side-coverage value from the left/right counts and
total span and uses it with total-span normalization when scaling local costs.

The 2013 `FUN_100230a0` uses those offsets differently: `+0x0c` is the
right-match count, `+0x0e` is total span, and `+0x10` is weighted side
coverage. Context-record modes 0, 1, and 2 contribute weights 2, 1, and 1 to
each matched side. The helper halves the side sums, then records the opposite
side when one side reaches the context edge, the smaller side value when
neither reaches an edge, or both sides' values plus one for full-span coverage.
Its ordering expression is `-(+0x0e + 100 * +0x10)`. The 2013 caller
`FUN_10023350` uses `+0x0e` as total span and combines `+0x10` with normalized
total span to scale the distance component returned by `FUN_100182e0`.
The scale is `DAT_1006d170 / (+0x10 + max(+0x0e / DAT_1006d174, 1))`;
the local `vt_pau.dll` stores float32 `1.0` and `2.0` at those addresses
(file offsets `0x6d170` and `0x6d174`). This multiplier scales the feature
distance component returned by `FUN_100182e0`, not its categorical penalty.
Before that scoring pass, `FUN_10023350` scores the leading full-span group
when the first continuity-ranked candidate covers the complete phone span.
Otherwise it scores every candidate in lists shorter than 31; larger lists
stop at the first strict increase in the initial key after 30 candidates, so a
tie at the boundary remains in the scored prefix. The pseudocode then handles
the remaining list and 30-entry output cap in a separate path.
For lists of at least 31, it multiplies the remaining nodes' local costs by
5.0 when their mapped class byte is ASCII `d`, by the float32 at
`DAT_1006d1a4` when the class byte is `0x0c` and the voice-state byte differs,
and by 1.0 otherwise. The local DLL stores raw bits `0x00002000` at file
offset `0x6d1a4` (a float32 subnormal). It sorts only that suffix in ascending
node-score order with `FUN_1001b5f0`, then sets the shortlist length to 30.
The larger-list sort uses partition swaps and can permute equal-score entries,
so the exact tie order is not inferred from its score comparator alone. The
category-map and voice-state producers remain context-dependent inputs.

Thus the two engines retain the same broad mechanism—adjacent-context coverage
influences candidate order and local-score scaling—but 2013 weights that
coverage and makes it dominate ordinary span differences. The candidate-node
offsets were repurposed consistently within the 2013 producer and consumers;
the evidence does not indicate a stale 2006 struct layout in one 2013 caller.
The marker representation changed as well: the 2006 helper reads a separate
byte-per-unit flag at model `+0x68`, whereas the 2013 helper tests bit `0x80` in
signature byte `+6`.

For the Apple context-3 intervention, marker clearing changed row `64729`'s
2013 total span from 4 to 2 while its weighted side-coverage term remained
zero. Thus the `100 * +0x10` rank component contributes nothing to this row's
observed two-point order-key change; total span alone changes its key. Forcing
only the computed total span back to 4 restored candidate order, score-prefix
inclusion, and the control WAV. This proves span mediation for that path,
while leaving ranking and score-normalization effects entangled because
`+0x0e` feeds both. The 2013 decompilation is preserved in
`tools/revkit/work/reports/vt-pau-2013-span-builder-consumers.txt`; the 2006
pseudocode is in `tools/revkit/work/reports/vt_eng-2006-all-functions-pseudocode.c`.

The source-index bytes confirm the marker mapping on this exact edge. Legacy
2005 `unit-gen` row `64728` has `attr_48=1`; the matched-context adapter has
signature byte `+6 = 0x00`, and the all-bank overlay has `0x80`. Rows `64727`
and `64729` also have `attr_48=1`, with adapter/overlay bytes `0x0a/0x8a` and
`0x00/0x80`. For candidate `64729`, captured 2013 fields (right count, total
span, weighted side coverage) change from `(1,4,0)` to `(1,2,0)` when only
row `64728`'s marker is cleared. The implied left-match count therefore falls
from 2 to 0, while its right-match count stays at 1. The marker removes two
left-side adjacent-context matches from this candidate's span. The all-bank
overlay manifest and runtime captures are under
`tools/revkit/work/stage20/index-adapter-key-repacked-matched-context-attrb-transfer-bit7-all-banks/`
and `tools/revkit/work/corpus-parity/stage20/apple-continuity-bit7/`.

The all-position Apple capture found natural nonzero weighted coverage at
positions 5 and 6. Unit `266885` has total span 4, `+0x10 = 1`, and signature
byte `+6 = 0x00` at both positions. At position 5, two left matches contribute
a weighted sum of 2; at position 6, three left matches contribute a weighted
sum of 3. In both cases the right side reaches the utterance boundary, so the
producer selects the left sum divided by two. Integer division yields 1 for
both rows. The rank expression gives this candidate key `-104`, while the
other observed rows at those positions have span at most 3.

The paired zero-field intervention changed only this candidate's `+0x10`
field before `FUN_1001b5f0` sorted the node list, and updated the corresponding
rank float. Its key changed from `-104` to `-4`, but it remained first at both
positions because span 4 was unique there. The trace confirms that the
weighted field contributes 100 to the numeric ordering key; the Apple case
does not show a different relative order when that contribution is removed.
All-position, producer-intermediate, and intervention captures are under
`tools/revkit/work/corpus-parity/stage20/apple-continuity-bit7/` in
`weighted-coverage-full-marker-correct-stack/`, `weighted-coverage-producer/`,
and `weighted-coverage-zero-266885/`. The corresponding reproducible scripts
are `tools/revkit/work/stage20/trace-apple-weighted-coverage.gdb`,
`trace-apple-weighted-coverage-producer.gdb`, and
`trace-apple-weighted-coverage-zero-266885.gdb`.

For this Apple edge, the proposed `attr_48` to bit-`0x80` conversion is verified
at the source and adapted row bytes, and the runtime intervention shows the
exact two-match span change. The remaining corpus-level question is whether
adjacent class membership and these coverage changes align across other
affected Kate contexts.

The natural Hello position-0 capture gives a case where `+0x10` changes the
pre-score candidate order. Unit `272822` has `(right count, span, weighted
coverage) = (5,6,4)` and key `-406`; unit `273369` has `(4,5,2)` and key
`-205`. Its producer trace records five right matches with weighted sum 6,
zero left matches, and the full-span calculation `0/2 + 6/2 + 1 = 4`.
Zeroing only unit `272822`'s `+0x10` before the sort changes its key to `-6`
and moves it from first to second behind `273369`. The output WAV hash stays
unchanged. A follow-up order-only intervention restored the original metric
and rank after sorting while preserving the altered node order. The selected
six-unit chain exactly matches the same-stack control
(`272822, 272823, 272823, 272824, 272825, 272825`), and both WAV hashes are
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`.

The follow-up traces locate the recovery. In the natural control, `272822` is
first and spans the complete six-phone position, so `FUN_10023350` takes its
leading-full-span branch and scores only that candidate; its local cost from
`FUN_100182e0` is `0.00121581`. In the order-only treatment, `273369` is first
with span 5, so the full-span branch is bypassed and the 30-candidate list is
scored. `273369` receives local cost `0.0156863`, while `272822` retains the
lower `0.00121581`. At the next transition, `FUN_10018c80` has one current
candidate (`272823`) and records predecessor `272822` at index 0 in control
and index 1 in treatment; both paths have the same cumulative cost `2.11883`.
Contexts 2–5 each have one candidate and preserve the same predecessor chain
(`272823`, `272824`, then `272825` twice), with matching cumulative costs
`4.30964`, `12.1469`, `15.4763`, and `19.7343`. The order change therefore
widens only the first scored list; its lower local cost still makes `272822`
the winning predecessor, after which the remaining one-candidate transitions
are forced. Traces, WAVs, and restoration records are under
`tools/revkit/work/corpus-parity/stage20/hello-weighted-local-costs-control/`
and `hello-weighted-order-only-local-costs/`; the corresponding scripts are
`trace-hello-weighted-local-costs-control.gdb` and
`trace-hello-weighted-order-only-local-costs.gdb` under `stage20/`.

This differs structurally from the 2006 selector's route to the shared rows.
On the same plain `Hello.` comparison, the old engine reduces pools of
75, 9, 5, and 6 candidates to the four full-span singletons
`272822 → 272823 → 272824 → 272825` before calling its local scorer. The 2013
trace has six phone positions and duplicate handoffs, but the order-only
intervention still retains the shared candidate at the first position through
local cost and backpointer selection. The measured 2006 local costs
(`0.00141844`, `1.64526`, `6.19391`, `3.9525`) are generation-specific and are
not directly comparable numerically with 2013's costs. Together, the traces
show how the two selectors reach the same four underlying row IDs through
different shortlist and path mechanisms; they do not establish general
cross-generation parity beyond this input and the tested overlays.

The Apple all-position, producer, and zero-field captures are under
`tools/revkit/work/corpus-parity/stage20/apple-continuity-bit7/` in
`weighted-coverage-full-marker-correct-stack/`, `weighted-coverage-producer/`,
and `weighted-coverage-zero-266885/`. The Hello counterparts are under
`tools/revkit/work/corpus-parity/stage20/hello-weighted-coverage-full-marker/`,
`hello-weighted-coverage-producer/`, and `hello-weighted-coverage-zero-272822/`.
The order-only and same-stack selected-unit control captures are under
`hello-weighted-order-only-272822/` and `hello-weighted-coverage-control-selected/`.
Reproduction scripts are in `tools/revkit/work/stage20/` with the matching
`trace-apple-weighted-coverage*` and `trace-hello-weighted-coverage*` names.
The `Hello.` order-only intervention establishes that 2013's `+0x10` field can
change the pre-score order without changing the selected chain or PCM. A
second controlled input, `Hello hello.`, was used to follow the same decision
through crowded later positions and compare the 2006 and 2013 transition
calculations.

For 2013 with the all-bank bit-`0x80` overlay, the intervention moves
`272822` behind `273369` at position 0 (`272822` changes from key `-205` to
`-5` before the original metric is restored after sorting). The control and
intervention preserve all 20 selected unit IDs, all 286 keyed transition rows,
their cumulative costs, and predecessor IDs; only `272823`'s predecessor
array index changes from 0 to 1. The 2006 runtime was then traced at its
backpointer write for candidate `272823`. That write selects predecessor
index 1, which maps through the live previous-position list to `272822`. This
corrects the earlier repeated-Hello link-only interpretation that read a
post-backtracking state table as a per-candidate link. For this edge, both
engines choose `272822`; the earlier apparent 2006 choice of `273369` was an
instrumentation mapping error.

The pair-cost traces show why that predecessor wins. In 2006, edge
`273369 → 272823` has transition feature `30.8357`, previous path cost
`0.0201681`, and current local cost `1.27656`, totaling `32.1324`. Edge
`272822 → 272823` has zero transition feature, previous cost `0.00243161`,
and the same local cost, totaling `1.27899`. The 2013 all-bank overlay run
gives the corresponding edges feature values `21.9468` and `0`, previous
costs `0.00882353` and `0.000945626`, and local cost `3.29502`; totals are
`25.2507` and `3.29596`. Both implementations select the same predecessor.
The complete position-1 predecessor row was captured for this current
candidate in both engines. Its next-lowest total is `21.2669` in 2006 and
`18.0691` in 2013, leaving margins of `19.9879` and `14.7731` over the
`272822` edge. The predecessor candidate sets are not identical, so these
margins establish a robust winner in each observed list, not equal search
spaces.

The raw feature lookups on the penalized edge are the same three values in
both builds, `{1.04473, 1.77778, 0.305296}`, but their slots and coefficients
differ. 2006 category 1 applies weights `{10, 10, 2}` and bias 2 to that
ordering. 2013 category 1 reads the latter two values in reverse slot order
and applies effective weights `{10, 2, 5}`, with the same bias 2. Therefore
the feature subtotal changes from `30.8357` to `21.9468`; this is a verified
transition-objective difference, but it does not flip the predecessor on
this edge. The relevant pseudocode is 2006 `FUN_1001e470` and 2013
`FUN_10018c80`. Captures are under `tools/revkit/work/corpus-parity/stage20/`
in `hello-repeat-2006-paircosts/`, `hello-repeat-2006-pairmatrix/`,
`hello-repeat-2006-backpointer/`, `hello-repeat-2006-feature-components-v2/`,
`hello-repeat-2013-paircosts/`, `hello-repeat-2013-pairmatrix/`, and
`hello-repeat-2013-feature-components/`.

The plain-package controls also confirm that overlay and engine effects must
remain separate. With identical `Hello hello.` text, the original 2006 row
builder emits nine rows (`273369, 273370, 273371, 232670, 266023, 264071,
264072, 264073, 264074`) and a 26,372-byte WAV. The standard adapted 2013 run
reports 16 candidate positions, all with zero weighted-coverage fields, and
produces a 17,334-byte WAV. The all-bank marker-overlay 2013 run reports a
different 20-unit stream and a 26,660-byte WAV. Those outputs establish a
real segmentation and synthesis-path difference for this fixture, but the
2013 overlay stream cannot be attributed to native engine changes alone.

## Late transition comparison and the Hi suffix-map intervention

The 2013 full legacy weight-and-slot remap changes one natural near-tie in
`Hello hello.`. At 2013 context 9 for current row `273370`, the baseline totals
are `888.132` through predecessor `181914` and `888.168` through `76688`, so
the baseline selects `181914`. With the full remap, those totals become
`889.611` and `888.396`, respectively, so the selected predecessor changes
to `76688`. The first `272822 → 272823` edge remains unchanged, as do the
position count and the rest of the post-prune chain. This confirms the metric
slot/weight difference can cause a late path change, but it does not explain
the extra segmentation.

The corresponding 2006 capture evaluates current row `273370` at context 1
against its own 30-row predecessor pool. It selects `273369` at total `2.3531`;
the next-lowest observed total is `16.809` through `182880`. Neither 2013
alternative (`181914`, `76688`) occurs in that old pool. The old and new
decisions therefore cannot be compared as the same pair of candidate edges:
their path contexts and predecessor lists differ. This is a concrete search
space difference, not evidence that the 2006 weights would choose one of the
two 2013 alternatives.

The second matched input, `Hi.`, isolates an earlier candidate-production
failure. In 2006, the two accepted queries are `90 34 17 30 30` and
`34 17 90 04 34`; the selector returns ten and nine candidates and selects
`272820 → 272821`. In the corrected-repack 2013 baseline, six exact class
lookups return zero and `272821` is absent from every candidate list. The
2013 target prefixes nevertheless match the repacked rows after the three
class-key tables:

| Item | First three key bytes | Trailing bytes |
| --- | --- | --- |
| Repacked row `272820` | `5a 22 01` | `1e 00` |
| 2013 whole-position query for its context | `5a 22 01` | `a0 00` |
| Repacked row `272821` | `22 11 5a` | `04 20` |
| 2013 whole-position query for its context | `22 11 5a` | `65 00` |

The row keys come from the repacked legacy index and the query keys from live
`FUN_10023dc0` returns. An isolated index treatment mapped only the two
observed prefix/suffix pairs: `1e 1e → a0 00` for prefix `5a 22 01`, and
`04 22 → 65 00` for prefix `22 11 5a`. In that treatment, the matching
whole-position lookups each return one class, while the other four queried
variants still return zero. 2013 candidate lists contain
`272820` at position 0 and `272821` at positions 1 and 2, and its selected
handoffs are `272820, 272821, 272821`. The WAV is 14,560 bytes (SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`). The
input and output fixtures were restored to identical hashes.

This intervention establishes that the two suffix mismatches are sufficient
to remove the observed `Hi.` candidate miss and recover the old selected row
IDs. It does not establish a general suffix-conversion rule: the map was
conditioned on these two observed key prefixes and suffixes. The remaining
2013 fallback and scoring behavior also has not been shown to match 2006.

The two candidate rows that matter at the first transition are now proven to
be the same source records in the legacy and adapted indexes. A read-only
payload comparison found byte-identical 19-byte records for IDs `272820`
(`915fb2ae…db99c8`) and `65331` (`9280ba13…64481b`); ID `272821` is also an
identical source row (`cd8c71d9…3e1348ba`). This removes numeric ID overlap as
the basis for the correspondence. The old-engine span-filter bypass admits
`65331` with its own coverage-1 weight and still selects `272820`: their
position-0 old scores are `0.0285714` and `0`, respectively.

The nine old position-1 candidates resolve to two distinct seven-byte feature
records: four rows (`272821`, `279498`, `279499`, `280449`) have
`22 11 5a 04 22 00 22`, and five (`4087`, `4102`, `84094`, `177718`,
`243540`) have `22 11 5a 04 04 00 22`. Three isolated repacks mapped the
second record's final pair `04 04` to each of the three remaining generated
2013 suffixes (`5d 00`, `55 00`, or `4d 00`), while preserving the existing
`04 22 → 65 00` and first-position map. Each mapping made its corresponding
exact 2013 lookup return one class. All three runs retained candidate counts
`67, 19, 19`, selected `272820, 272821, 272821`, and emitted the same
14,560-byte WAV with SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.
Therefore these outputs do not distinguish which generated key is the
intended destination for the five-row source class; they only show that each
one-class placement leaves this fixture's ranked stream unchanged. The exact
producer-to-class correspondence remains to be recovered from the field
conversion, not inferred from the unchanged audio.

The new 2013 transition trace shows why that same first row wins when both
shared records are available. At context 1, `272820` and `65331` enter with
the same preceding path cost (`5.4`). For current row `272821`, the measured
transition term is `2` from `272820` and `10.1671` from `65331`; with the same
current local cost (`19.309`), the resulting totals are `26.709` and
`34.8761`. The difference is entirely in the adjacent-context term for this
pair. At context 2, repeating `272821` has zero transition term and zero
penalty, so its cumulative total is `54.709`; backtracking selects
`272820, 272821, 272821`. Thus the shared first edge is not being decided by
the local-score tie or by the 2006 span gate alone: 2013's measured transition
objective also prefers the `272820 → 272821` edge over `65331 → 272821`.

The current suffix-map capture has three candidate positions, while the 2006
capture has two. The earlier unadapted 2013 candidate capture has a different
four-position stream and WAV hash; it is a separate repack state and is not
combined with these transition totals. This transition trace applies only to
the two-entry suffix-map overlay, whose other generated lookup variants still
miss. Reproduction uses `trace-2013-hi-pair-transition.gdb` and the matching
`compose-key-repacked-attr48-key3-attr40-key2-hi-tail-map.yaml`; logs and the
captured WAV are under `corpus-parity/stage20/hi-2013-pair-transition/`.
The low-instrumentation control is in `hi-2013-pair-control/` and has the same
WAV hash and selected IDs, confirming the transition breakpoints did not
change this run's path. Both captures verify exact fixture restoration.
The row-level check is reproducible with
`tools/revkit/work/stage20/crosscheck-hi-source-records.py`; the three
single-suffix overlays are built with `repack-index-key-fields.py` variants
`attr48-key3-attr40-key2-hi-0404-to-{5d00,5500,4d00}`. Their runtime evidence
is under `corpus-parity/stage20/hi-relax-0404-to-{5d00,5500,4d00}/`.

### Query producer comparison: runtime evidence

The producer paths are now traced at both ends for `Hi.`. In the 2006 runtime,
`FUN_1001da00` receives three unique key operands from `FUN_1001da80` at call
site `0x1001dcb3`: `90 34 17 30 30`, `34 17 90 04 34`, and
`34 17 90 04 04`. The first two are the two queries accepted by the 2006
splitter; the third appears in the expansion trace but is not one of those
accepted split queries. The three operands repeat across the runner's three
render passes. This trace records the actual query bytes, rather than only
the inputs to the expansion routine.

In 2013, `FUN_10018770` receives one target-signature source for each
phone-tree leaf and different six-byte context rows. The slot-0 and slot-1
calls use source signature `5a 5a 22 11 5a a0 00`; the slot-2 and slot-3
calls use `5a 22 11 5a 5a 65 00`. The
slot-0 rows differ at byte 4 (`00` then `01`); the slot-1 row has byte 4
`02`; the slot-2 rows have byte 4 `00` then `01`; and the slot-3 row has byte
4 `02`. The full rows and call order are recorded in
`hi-2013-producer/adapted-forced-pah0-winedbg.log`.

The 2013 pseudocode explains the generated alternatives: `FUN_10018770`
selects two relaxation vectors from the target signature's packed byte 5,
filters them using context-row bytes and flag bits, then composes those
choices into signature byte 5. `FUN_10016ea0` copies that generated byte into
key byte 3 and masks signature byte 6 with `0x20` for key byte 4. The live
class lookups confirm the results: slot 0 emits suffixes `a0 00` and
`98 00`; the next context emits `65 00`, `5d 00`, `55 00`, and `4d 00`.
Those bytes are generated relaxation codes. They are not a direct copy of the
legacy query's decimal-pair suffix bytes (`30 30`, `04 34`, and `04 04`).

This identifies the concrete mechanism behind the `Hi.` key misses: the
repacked rows still carry legacy suffix values, while the 2013 producer emits
suffix values assembled from its packed signature and context-specific
relaxation choices. The two-entry overlay proves that translating the two
successful old rows into the corresponding 2013 codes restores those row
lookups. It does not yet supply a general row conversion: the generated code
depends on target-signature and context fields, and different variants of a
single context produce different suffixes. The mismatched suffix domain is a
verified cause of the `Hi.` misses; whether the same mechanism explains the
extra positions and wrong selections across all matched fixtures remains to
be tested.

The expansion trace is `legacy-expansion-hi/original-forced-pah0-winedbg.log`
and its script is `trace-legacy-query-expansion.gdb`. The paired 2013 producer
trace is `hi-2013-producer/adapted-forced-pah0-winedbg.log`, produced by
`trace-2013-query-producer.gdb`; the six transformed lookup keys are in
`hi-2013-keyclass-lookup/adapted-keyclass-lookup-gdb.log`. Both runners
reported exit 0 and identical before/after fixture hashes.

The same producer capture was added to plain `Hello.`. The old engine accepts
four keys with candidate counts 75, 9, 5, and 6. The corrected-repack 2013
trace emits ten exact key lookups as its contexts expand: nine return zero;
only `17 2b 2f 00 00` returns one class. Its source signatures carry byte-5
codes `a0`, `20`, `00`, and `69` across the observed leaf groups, while the
corresponding old query suffix pairs are `30 30`, `00 30`, `00 00`, and
`04 04`. The runtime rows and source signatures are in
`hello-2013-producer/adapted-forced-pah0-winedbg.log`; exact query keys and
counts are in `hello-2013-producer/adapted-query-keyclasses-gdb.log`. Both
2013 runs exited normally and verified fixture restoration.

This extends the observed failure from `Hi.` to the plain `Hello.` repack:
most generated exact keys do not find classes, so the producer hands sparse
candidate populations to the fallback path. The Hello comparison also changes
the transformed prefixes, so it does not establish that suffix conversion
alone explains every miss. The next conversion must account for both the
three table-transformed key bytes and the generated relaxation suffixes.

The same 2013 trace on `Hello hello.` records 20 producer entries and 24 exact
lookups. Only five lookups return a class; the other 19 return zero. The
utterance includes context-specific code families such as `a0/98`,
`20/18/10/08`, and `45/4d`; the repeated `17 2b 2f 00 00` key returns one
class in each word instance. This shows that the sparse exact-hit pattern
persists when the word is repeated and that the producer emits
context-specific alternatives across the repeated utterance. The complete trace is
`hello-repeat-2013-producer/adapted-forced-pah0-winedbg.log`; the runtime
returned exit 0 and restored both input and output hashes. Existing 2006
selected-chain evidence records nine output rows for this fixture. Together,
these measurements connect sparse 2013 lookup coverage with a concrete
segmentation-difference input, but they do not yet prove the exact conversion
that would make the two candidate populations equal.

### Why a suffix map cannot restore the Hello candidate sets

The read-only class-member crosswalk now identifies two independent losses in
the corrected Hello repack. For slot 0, the 2006 feature builder maps 75 rows
to classes with prefix `5a 22 04` and suffix `1e 1e`; the 2013 target query
uses prefix `5a 22 17` and suffix `a0 00`. Thirty rows shared by the engines
have raw source prefixes `5a 22 18`, `5b 22 18`, `5a 22 1e`, or `5a 22 17`,
which the 2013 lookup tables map to `5a 22 17`. The 45 legacy-only rows have
raw prefixes `5a 22 05` or `5b 22 05`, which those tables map to `5a 22 04`.
All 45 have the same legacy seven-byte feature prefix `5a 22 04`; the
third-byte lookup therefore splits one old class across two 2013 prefixes.

Slot 4 shows the opposite problem. Five legacy rows with raw prefix
`17 2b 30` map to the 2013 target prefix `17 2b 2f`, and all five are shared.
Twenty further rows with raw prefixes `18 2b 2f` and `19 2b 2f` also map to
`17 2b 2f`, so the 2013 class has 25 rows. Their legacy seven-byte features
end in `2f 18` or `2f 19`, while the five legacy members end in `30 17`.
The 2013 five-byte key builder does not retain those final two legacy feature
bytes, so it merges rows that 2006 distinguishes. Slots 2 and 5 have exact
member parity in the same crosswalk (9/9 and 6/6).

This is more than a suffix conversion defect: the 2013 class tables remap the
third key field differently at slot 0, and the 2013 key projection loses two
legacy feature distinctions at slot 4. The results explain the measured
45-row subset and 20-row superset without invoking local ranking. A
suffix-only adapter cannot make these populations equal. The member-level
output is reproducible with
`tools/revkit/work/stage20/analyze-key-class-prefix-crosswalk.py`; it compares
the live target prefixes and old feature records against all source index rows.

### Remaining experiments, in order

1. Extend the now reproducible Hello member crosswalk to the Hi classes and
   the repeated-Hello query contexts. Preserve both the transformed prefixes
   and the complete legacy seven-byte features; suffix-only summaries conceal
   the slot-0 and slot-4 failures.
2. Test whether the two lost legacy feature dimensions can be reconstructed
   from other 2013 signature fields or require a change to the 2013 lookup
   path. In parallel, test a field-based treatment of the slot-0 third-byte
   mapping and the Hi suffix misses, with no utterance-specific overrides.
   Exact class membership should be checked against the source-row inventory
   before any output quality claim.
3. The `>2` versus `>9` experiment already shows the strict 2013 check
   rejects an exact six-row class that 2006 accepts. Extend this comparison
   only after identical class member lists are recovered for the remaining
   contexts; then trace post-prune paths and selected rows.
4. Finally compare the resulting selected streams and waveform construction,
   then validate the conversion on another 2005 voice package. Completion
   requires reproducible producer, class-population, and selector behavior
   that accounts for both the extra positions and wrong row choices across
   the matched fixtures.

Do not use listening impressions alone as the causal test. The 2006 key
expansion and 2013 target-signature trace together establish the next concrete
test: reconstruct the semantic suffix fields, then verify candidate
membership before comparing selection thresholds or audio.

Completion criterion: explain the extra positions and wrong row choices with
reproducible producer or selector behavior across the matched fixtures, and
state the remaining boundary if any difference survives the controlled
repack. Do not use listening impressions alone as the causal test.

The old query capture is in `corpus-parity/stage20/hi-2006-query-classes/`;
the 2013 baseline and suffix-map query/list captures are in
`hi-2013-keyclass-lookup/`, `hi-2013-repacked-candidates/`, and
`hi-2013-suffix-map/`. The late old-edge trace is in
`hello-repeat-2006-late-edge/`; the paired 2013 captures are in
`hello-repeat-2013-repacked-late-edge-control/` and
`hello-repeat-2013-repacked-late-edge-full-legacy-map/`. Reproduction scripts
are `trace-legacy-hi-query-classes.gdb`,
`trace-2013-hi-suffix-map.gdb`, and
`trace-legacy-repeat-late-edge.gdb` under `tools/revkit/work/stage20/`.

### Hello legacy-feature crosswalk and controlled cutoff A/B

A diagnostic index overlay maps every source row whose measured 2006
seven-byte feature record belongs to one of the four observed `Hello.`
shortlists to that position's primary 2013 key. Rows that would collide with
those keys without belonging to the measured legacy class are redirected to
a sentinel suffix. The resulting four exact 2013 keys contain 75, 9, 5, and 6
rows, respectively. This is a population crosswalk, not a general conversion:
it maps multiple old feature classes at a position to one 2013 key and does
not populate the additional relaxation keys emitted by the 2013 producer.

The native-cutoff runtime trace confirms the split decision in that controlled
overlay. The four 2013 metric inputs are one class each with weights `75`,
`9`, `5`, and `6`. `FUN_10024060` accepts 75, rejects 9, 5, and 6 under its
`>9` rule, and takes three two-position fallbacks, yielding seven positions.
The paired intervention changes only a failed split-helper result when its
measured sum is above 2. It accepts 9, 5, and 6 and yields four positions.
Both captures use identical detailed GDB instrumentation; the native-cutoff
WAV is 15,004 bytes with SHA-256
`d4878c9814ef2b9fe9186a2c1c6f1e7128519ed0f222bd5c3ef6ddeaa791fab0`, and
the `>2` WAV is also 15,004 bytes with SHA-256
`6be52f6b953a43d286ff8e23d0d707b36f60cf9e024686cf94e33f71108d4d35`.
Both runner records confirm the Stage 5 input and output fixtures were
restored. They do not record a hash of the source text fixture, so these WAV
hashes are not attached to the explicit `Hello.` control below.

In that capture, accepting sums `9`, `5`, and `6` above 2 changes the split
result from seven positions to four after the mapped populations are present.
It does not establish equivalent candidate composition or ranking. The
diagnostic overlay collapses the four old slot-0 class groups and the two old
slot-2 class groups into one target class each; its mapped canonical keys
leave the other emitted 2013 variants empty. The source fixture for the
15,004-byte capture was not fingerprinted, so its selected chain
(`52542, 264072, 264073, 264074`) is retained as an observed overlay result,
not as a verified `Hello.` comparison.

The reproducible analysis and overlay are
`tools/revkit/work/stage20/analyze-key-class-prefix-crosswalk.py` and
`repack-index-legacy-feature-class-crosswalk.py`. Native and overridden
traces are `trace-2013-hello-legacy-feature-crosswalk-native-cutoff.gdb` and
`trace-2013-hello-legacy-feature-crosswalk-old-cutoff.gdb`. Captures are in
`tools/revkit/work/corpus-parity/stage20/hello-legacy-feature-crosswalk-native-cutoff/`
and
`hello-legacy-feature-crosswalk-old-cutoff-v2/`.

### Explicit `Hello.` score and transition replay

The earlier four-position capture did not record the source fixture hash. It
therefore cannot serve as a verified `Hello.` capture. The following
comparison uses the explicit `tools/revkit/work/stage20/hello-plain.txt`
fixture (SHA-256
`a2c064616af4c66c576821616646bdfad5556a263b4b007847605118971f4389`) for both
engines, plus the same legacy-feature crosswalk and the controlled `>2`
intervention for 2013.

The 2006 engine selects `272822, 272823, 272824, 272825`; its local scores are
`0.00141844, 1.64526, 6.19391, 3.9525`. The scores come from
`FUN_1001c860` and are not numerically comparable to 2013's
`FUN_100182e0`. The 2013 output is 18,956 bytes
(`31c1c43d8aea7500f41f9e3e8578651a1d523adbc9d7567221ca003949f38d78`) and
selects `52542, 273370, 273371, 255282`.

At position 0, the raw 75-row list contains old-selected `272822` and `52542`;
both receive a 2013 local cost of `9.6`. The ranker returns only 30 rows to
the transition pass. That list contains `52542` but not `272822`. The 30
retained rows have the same observed local score (`9.6`) and rank field
(`-1`), yet their order differs from the raw list. Disassembly shows
`FUN_1001b5f0` sorting node pointers by the float at node `+4`; the equal-rank
reordering and 30-row cap explain why the legacy row is absent before
backtracking, though the precise tie permutation is an implementation detail
of the sort. Thus the first-position difference is not a backtrack choice
between two tied rows: `272822` has already been removed from the transition
pool.

For positions 1–3, the post-rank lists still contain the old rows. The direct
2013 local scores are:

| Position | Old row and cost | 2013 row and cost | Other observed candidate |
| ---: | ---: | ---: | --- |
| 1 | `272823`: 26.5764 | `273370`: 24.9582 | `264072`: 27.1761 |
| 2 | `272824`: 173.792 | `273371`: 136.48 | `264073`: 142.332 |
| 3 | `272825`: 71.25 | `255282`: 74 | `264074`: 75.25 |

The local scorer favors the selected 2013 row at positions 1 and 2, but favors
old row `272825` at position 3. The 2013 cumulative path costs explain the
last choice:
the old path through positions 1–3 reaches `151.75`, while the selected path
reaches `137.431`. At position 2, the old-to-old path reaches `114.125`; the
new-to-new path reaches `94.0361`. At position 3, the selected path's current
edge has a higher measured distance term (`6.39509` versus `2`) and local
term (`37` versus `35.625`), but starts from a cumulative cost 20.089 lower.
The selector therefore keeps the lower-cost cumulative path instead of
minimizing each position independently.

A per-call counterfactual swaps only the six scalar feature bytes between
`272823`/`273370`, `272824`/`273371`, and `272825`/`255282` while
`FUN_100182e0` scores those rows, then restores them before returning. It
leaves class keys, signatures, `attr_b`, and transition metric codes unchanged.
Each pair's two local costs exchange exactly. The resulting selected chain is
`179388, 272823, 272824, 272825`; its 15,956-byte WAV hash is
`abb49098891b965f6974388485432e3488c92c5f8f1365a1af44447b75ddde3b`. This
shows that the six feature bytes drive the measured pairwise local preference
under the 2013 scorer. It does not identify their intended cross-generation
mapping: the intervention does not produce the exact old chain, and the first
row still changes through the broader path calculation.

The rank, score-input, feature-swap, and transition traces are
`trace-2013-hello-crosswalk-score-fields.gdb`,
`trace-2013-hello-crosswalk-feature-swap.gdb`, and
`trace-2013-hello-crosswalk-transition-path.gdb` under
`tools/revkit/work/stage20/`. Evidence is in
`hello-plain-crosswalk-old-scores/`,
`hello-plain-crosswalk-score-fields/`,
`hello-plain-crosswalk-feature-swap/`, and
`hello-plain-crosswalk-transition-path-v3/`. Fixture restoration hashes match
before and after each runtime capture.

### Slot-2 key-partition intervention

The next controlled overlay keeps the measured old slot-2 classes separate:
the weight-5 feature group maps to `22 17 2b 20 00`, and the weight-4 group
maps to `22 17 2b 18 00`. The latter assignment is an experiment using an
observed 2013 relaxation key, not a recovered semantic conversion. The exact
`Hello.` fixture was hashed before each run (`a2c064616af4c66c576821616646bdfad5556a263b4b007847605118971f4389`),
and the runner restored the Stage 5 input and output hashes after both runs.

Direct 2013 lookup returns one class for each populated key, with weights 5
and 4. The selector combines both returned class IDs before evaluating its
split metric: `13480:5` plus `13479:4` returns sum 9. Under the native `>9`
rule this is rejected and slot 1 expands to two rows. The matched 2006 trace
for the same input returns the two corresponding old groups, sum 9, and
accepts them under `>2`. This confirms the cutoff difference for this
successful, nine-member lookup. It also shows that splitting those members
between queried 2013 keys does not change the accumulated metric.

With the controlled `>2` intervention, the split overlay returns four
positions and selects `52542, 273370, 273371, 255282`. Its WAV is 18,956 bytes
with SHA-256
`31c1c43d8aea7500f41f9e3e8578651a1d523adbc9d7567221ca003949f38d78`,
byte-identical to the unsplit crosswalk intervention. Candidate counts remain
75, 9, 5, and 6 in the four positions; the selected chain is unchanged. The
native-cutoff split overlay returns seven positions and selects
`52542, 273370, 273370, 273371, 273371, 255282, 255282`; its 18,956-byte WAV
hash is `66d692e95be283065147bf578fb089a48db55e667bca141a3d9e840bb3c109b7`.
Thus the partition changes key-level lookup results and class IDs, but does
not explain the later wrong-row choices when both classes are collected. The
cutoff remains the cause of the extra two-row expansion at this populated
slot. Other native fallbacks include empty whole-position lookups and are not
explained by the cutoff alone.

The overlay builder is
`tools/revkit/work/stage20/repack-index-legacy-feature-class-slot2-split.py`;
the isolated mount is declared by
`compose-key-repacked-legacy-feature-class-slot2-split.yaml`. Captures are in
`hello-slot2-split-native/` and `hello-slot2-split-override/` under
`tools/revkit/work/corpus-parity/stage20/`; traces are
`trace-2013-hello-slot2-split-native-cutoff.gdb` and
`trace-2013-hello-slot2-split-candidate-path.gdb` under `stage20/`.

### 2006 span-filter bypass and path-cost replay

The earlier unit-ID-only scorer substitution is superseded. It retained the
old full-span candidate's `param5` score weight while substituting a partial
candidate ID, so its reported local costs were not the costs those partial
candidates receive in the old engine's actual candidate lists. The corrected
trace bypasses only the full-span reduction immediately before old scoring;
the original coverage values and each candidate's own score weight remain in
place. The bypassed run exits normally, keeps the original selected chain
`272822, 272823, 272824, 272825`, and emits the byte-identical old WAV
(`eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`).

With candidates admitted and scored at their native weights, the old scorer
returns:

| Position | Old full-span candidate: coverage / weight / cost | Partial candidate: coverage / weight / cost |
| ---: | --- | --- |
| 0 | `272822`: 4 / `0.166667` / `0.00141844` | `52542`: 1 / `1` / `14.5684`; `273369`: 3 / `0.285714` / `0.0201681` |
| 1 | `272823`: 4 / `0.166667` / `1.64526` | `273370`: 3 / `0.4` / `2.3784` |
| 2 | `272824`: 4 / `0.166667` / `6.19391` | `273371`: 3 / `0.666667` / `4.61288` |
| 3 | `272825`: 4 / `0.166667` / `3.9525` | `255282`: 1 / `1` / `6.04246` |

These are direct scorer returns from the bypass run, not cross-engine score
comparisons. For the direct 2013-selected row chain, the old scorer prefers
its selected row at positions 0, 1, and 3, and prefers `273371` locally at
position 2. The bypassed path's strongest partial-chain competitor starts
with `273369` (not `52542`); its local cost is also above old row `272822`.
The path
trace resolves that local exception: the old chain's cumulative costs are
`1.64668`, `7.84059`, and `11.7931` through positions 1–3, with zero edge
transition terms. The partial alternative
`273369 -> 273370 -> 273371 -> 255282` reaches `2.39857`, `7.01144`, and
`19.9869`; its final edge contributes `6.93303`. Thus the 2006 path objective
still prefers its original complete chain after the coverage gate is bypassed.
The gate explains ordinary eligibility, but does not by itself explain the
old-versus-new row choice. Under the 2013 objective, the corresponding partial
chain totals `137.431` versus `151.75` for the old-row alternative; those raw
totals have a different scale and only establish ordering within 2013.

The score-entry/return and edge traces are in
`hello-plain-crosswalk-old-span-filter-bypass-v3/` and
`hello-plain-crosswalk-old-span-filter-bypass-v4/` under
`tools/revkit/work/corpus-parity/stage20/`. The interventions are
`trace-legacy-hello-bypass-full-span-filter-v3.gdb` and
`trace-legacy-hello-bypass-full-span-filter-v4.gdb`. The runner's Stage 17
input/output restoration hashes match. The unit-ID-only capture remains in
`hello-plain-crosswalk-old-candidate-scores/` as superseded evidence.

The weight's producer is visible in the 2006 pseudocode for
`FUN_1001cff0`, immediately before its `FUN_1001c860` call. Its preceding
`FUN_1001cd60` call stores candidate left reach at node `+0x0c`, right reach at
`+0x0e`, and total span at `+0x10` (`left + right + 1`). The caller selects a
directional reach for the current position, computes a floored span term from
`total_span / _DAT_1006116c`, and passes
`_DAT_10061168 / (directional_reach + span_term)` as scorer `param5`. The
captured weights fit this formula: isolated one-position rows receive 1;
coverage-3 rows vary by position and side; full-span-4 rows receive
`0.166667`. This is the direct connection between the old continuity/span
metadata and the local-feature contribution to its score. The decompiler's
casts and the names of the two constants remain unverified; the prose above
describes the pseudocode dataflow, not recovered source-level intent.

For this exact `Hello.` fixture, the combined explanation is now more precise:
2006 removes partial-span candidates before scoring and, when that gate is
bypassed, its own weighted local-plus-transition objective still selects the
old chain; the adapted 2013 candidate population and objective rank the
partial chain lower. The causal role of individual score arrays is supported
by the one-field interventions below, but their producer and semantics remain
unidentified. The existing `Hi.` runs do not yet form a matched score
comparison over the complete utterance: the old path has two positions and
the suffix-map 2013 path has three, while an earlier unadapted repack produced
four. The source rows are now crosswalked by byte-identical 19-byte payloads,
and the 2013 pair trace directly compares shared rows `272820` and `65331`.
It shows the shared `272820 → 272821` edge winning on the measured transition
term. The old bypass trace independently shows that 2006 still prefers
`272820` when `65331` is admitted and scored at its own weight. The remaining
Hi comparison is the extra 2013 position and its key production, plus whether
the 2013 feature-byte interventions generalize to this proven row pair; raw
cross-engine score magnitudes are still not comparable.

### Per-field scalar-score intervention

The six scalar arrays were then swapped one at a time between the same
old-selected and new-selected row pairs while `FUN_100182e0` scored them.
Their pointers are read from the runtime model object at offsets `+0x48` through
`+0x5c`; these are retained as field indices because their semantic labels
remain unknown. On all three target rows, the runtime target record's byte
`+4` is zero. The decompiled 2013 scorer consequently reads the first and
third arrays of its three-array contextual block at model offsets `+0x54` and
`+0x5c` for this path. This matches the measured sensitivity:

| Array slot | Pointer offset | Pair costs changed? | Four selected IDs | WAV SHA-256 |
| ---: | ---: | --- | --- | --- |
| 0 | `+0x48` | No | `52542, 273370, 273371, 255282` | `31c1c43d8aea7500f41f9e3e8578651a1d523adbc9d7567221ca003949f38d78` |
| 1 | `+0x4c` | No | same as baseline | same as baseline |
| 2 | `+0x50` | No | same as baseline | same as baseline |
| 3 | `+0x54` | Yes | `52542, 264072, 264073, 264074` | `6be52f6b953a43d286ff8e23d0d707b36f60cf9e024686cf94e33f71108d4d35` |
| 4 | `+0x58` | No | same as baseline | same as baseline |
| 5 | `+0x5c` | Yes | `179388, 272823, 272824, 272825` | `abb49098891b965f6974388485432e3488c92c5f8f1365a1af44447b75ddde3b` |

The baseline for unchanged slots is the same four-ID chain and hash listed in
the slot-2 experiment. Swapping slot 3 changed the three measured pair costs
to `29.3758/22.1588`, `156.121/154.15`, and `68.5/76.75` (old/new order).
Swapping slot 5 changed them to `22.1588/29.3758`, `154.15/156.121`, and
`76.75/68.5`, exactly exchanging the slot-3 intervention's pair costs in
these traces. Slots 0, 1, 2, and 4 left all six costs unchanged. A slot-5-only
swap reproduces the earlier all-six swap's selected chain and WAV hash, so
slot 5 alone is sufficient for that particular path change. The first
position still differs from old unit `272822`, which has already fallen out
of the 30-row rank list before backtracking.

These interventions isolate which arrays affect this scorer for these
contexts; they do not supply names or establish a general cross-generation
field conversion. The pseudocode does recover the immediate dataflow: in
`FUN_100182e0`, the target record's byte `+4` chooses two pointers from a
three-array block beginning at model `+0x54`. For byte value zero, the scorer
reads arrays at `+0x54` and `+0x5c`, indexes each by candidate unit ID, and
uses those two byte values with target bytes `param4[1]` and `param4[2]` to
fetch two entries from the pair-distance matrix at model `+0x18`. The `+0x58`
array is the middle member of that block and is not selected for these target
rows. This explains structurally why swapping `+0x54` or `+0x5c` changes the
measured local cost while swapping `+0x58` does not; it does not identify what
the byte dimensions mean in the proprietary model schema. The changed path
also shows that a single scalar field can alter which middle candidates the
cumulative transition pass selects. The six trace scripts are
`trace-2013-hello-crosswalk-feature-field-{0..5}-swap.gdb`; their outputs are
in `hello-plain-crosswalk-feature-field-{0..5}/` under
`tools/revkit/work/corpus-parity/stage20/`. All six runs restored the shared
Stage 5 fixture hashes. The consumer-side explanation is established here;
the source-column crosswalk below identifies the six arrays for these adapted
Kate rows. The Hi scalar-field comparison is now complete for two candidate
pairs under the matched `5d 00` suffix overlay.
For the first-position pair `272820`/`65331`, the context record ends in byte
`+5 = 1`; `FUN_100182e0` takes its nonzero-flag branch and skips the six scalar
arrays. Each field differs between the rows, but six separate swaps leave
both local returns at `10.8`, preserve the selected two-row chain, and
reproduce the same 14,560-byte WAV. The old engine instead gives these rows
local costs `0` and `0.0285714`; in the 2013 path, the next transition favors
`272820 -> 272821` over `65331 -> 272821` by `8.1671`. Thus these two paths
reach the same predecessor through different measured terms.

At the second position the context byte `+5` is zero, so the scorer reads the
scalar arrays. Rows `272821` and `280449` have identical legacy seven-byte
feature records, but the controlled 2013 return scores are `57.118` and
`66.5484`. Swapping each scalar field independently changes the pair scores
only for model offsets `+0x54` and `+0x5c`: the selected/alternative returns
become `58.0484/65.618` and `55.618/68.0484`, respectively. The other four
swaps leave both scores unchanged. All six treatments retain selected IDs
`272820, 272821` and the identical WAV hash
`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`.
In the old span-filter-bypass capture, the corresponding local returns are
`10.3087` for `272821` at weight `0.333333` and `64.2757` for `280449` at
weight `1`; the old path also favors `272821`. These are within-engine scores
from different formulas and normalizations, so their magnitudes are not
cross-comparable. The new fields change local cost in this second context,
but do not explain a wrong selected row in this controlled Hi chain. The
remaining position-count difference still lies in query expansion and
candidate production.

The twelve field-swap logs and WAVs are under
`tools/revkit/work/corpus-parity/stage20/hi-2013-feature-field-{0..5}/` and
`hi-2013-current-field-{0..5}/`. The context-byte trace is in
`hi-2013-score-contexts/`, produced by
`trace-2013-hi-score-contexts.gdb`; the first-pair confirmation is in
`hi-2013-feature-context-gate/`. The per-field scripts are
`trace-2013-hi-feature-field-{0..5}-swap.gdb` and
`trace-2013-hi-current-field-{0..5}-swap.gdb`. Each runner restored identical
Stage 5 fixture hashes. The source-column mapping for these rows is established
below; semantic labels and generalization across other voice packages remain
open. The remaining position-count difference is in query expansion and
candidate production.

### Source mapping for the six scalar arrays

The six arrays are now traced back to the legacy index payload. The 2006
reader pseudocode loads three metric groups, each containing a two-byte value
and two one-byte values. In the 2013 versioned index, the same three groups
are represented column-major as `(word16, byte_a, byte_b)` columns. The Stage
20 repacker preserves the source `metrics[12N]` tail unchanged. For a unit at
local index `i`, its six one-byte lanes are therefore metric columns
`2, 6, 10, 3, 7, 11` in the source tail. The first three populate model
pointer slots `+0x48`, `+0x4c`, and `+0x50`; the last three populate
`+0x54`, `+0x58`, and `+0x5c`.

| Global row ID | Source bank and local row | Slots `+0x48/+0x4c/+0x50` | Slots `+0x54/+0x58/+0x5c` |
| ---: | --- | --- | --- |
| 272820 | `gen2`, 92825 | `85, 126, 168` | `78, 69, 77` |
| 65331 | `gen`, 65331 | `71, 108, 156` | `64, 80, 81` |
| 272821 | `gen2`, 92826 | `168, 171, 127` | `77, 75, 117` |
| 280449 | `etc`, 1324 | `161, 169, 124` | `81, 86, 111` |

The values in this table were read from Kate's original `ver.2005` indexes,
then matched against the live 2013 scorer arrays for all four rows. A separate
native 2006 capture reads `85, 126, 168, 78, 69, 77` for row 272820, matching
the same source record. This establishes a faithful byte transfer for these
rows; it does not name the three group dimensions or prove the mapping for
every package. In the Hi local-cost branch, the selected `+0x54` and `+0x5c`
arrays are the second one-byte lanes of metric groups 1 and 3. The swap
experiments therefore changed original metric bytes, not bytes synthesized
from the repacked signature.

The new trace scripts are
`trace-2013-hi-scalar-array-state.gdb` and
`trace-legacy-hi-scalar-array-state.gdb`. Their logs and runner restoration
records are in `corpus-parity/stage20/hi-scalar-array-state/` and
`corpus-parity/stage20/hi-2006-scalar-array-state/`. Both runs exited normally;
the 2013 state-only trace emitted the native three-position 14,560-byte WAV
with SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`. Each
runner restored the shared input/output fixtures byte-for-byte. This removes
the six-array source-column mapping as a likely cause of the 2013 failure on
the tested rows. The remaining generation difference is in how the engines
build query keys and candidate sets and how their distinct scoring and
transition paths consume them.

### Hi split threshold and PCM gain, isolated

The `5d 00` relaxation overlay restores the five-row class beside the original
four-row class at the matched second Hi query. The 2013 aggregator receives
class weights `13468:4` and `13467:5`, for a sum of 9. Its strict `>9` test
rejects that complete class set and expands the slot into two positions; with
the already accepted first position, the builder returns three. A debugger
intervention that accepts only this sum-9 result reduces the builder to two
positions, with candidate counts 67 and 9 and selected IDs `272820, 272821`.
That is the 2006 selected two-row chain for the same fixture. The native
three-position `5d 00` capture has IDs `272820, 272821, 272821`. This directly
confirms that the strict cutoff, after the missing five-row class is made
visible, is the mechanism adding the third position in this Hi case. The old
2006 `>2` rule accepts the same nine-member total.

The three-position WAV is 7,258 PCM frames and hashes to
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`. The
two-position override is also 7,258 frames and hashes to
`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`; a
repeat produced the same hash. Only 67 of 7,258 aligned samples differ between
those two 2013 outputs, all within samples 2790–2861 and by at most one PCM
count. Thus the extra position is a real selector decision, while its effect
on this rendered fixture is very small at the sample level.

With the two-position chain fixed, the 2013 decoded-unit routine
`FUN_1002c8b0` receives PCM scale 200 for both row calls (1,420 and 5,910
samples). The decompiled routine multiplies each decoded sample by that field
and divides by 100. Changing only those two row fields to 100 halves the
2013 output level: relative to the 2006 WAV, its fitted sample slope changes
from 1.99959 to 0.99979, and its RMS ratio changes from 1.99965 to 0.99982.
The gain-adjusted output is not byte-identical to 2006: 283 samples differ,
with maximum absolute difference 508 and mean absolute difference 3.73; the
two files retain the same 7,258-frame duration. The differences cluster at
samples 0–79, 1348–1419, and 7125–7257, with one additional sample at 7159.
The 72-sample cluster at 1348–1419 ends exactly where the 1,420-sample first
decoded row reaches its 72-sample overlap with the following row. This locates
the residual differences at the utterance edges and the join, rather than
through the shared interior samples. The scale field therefore accounts for
the near-exact 2x level difference. A paired PCM-window capture below now
isolates how the two join routines create those remaining differences for this
Hi. fixture.

A paired hardware-breakpoint trace of the unmodified 2006 run resolves the
corresponding runtime fields. Across all three render passes, its two selected
rows have sample counts 1,420 and 5,910, matching 2013, and effective PCM
scales 100/100; its pitch and duration scales are also 100/100. The 2013
decoded-row calls have the same sample counts and pitch/duration values, but
their sample scale is 200 on both rows. Thus the mismatch is present in the
per-row values entering the two generation-specific decoded-sample paths, not
in selected unit payload or row length. The old producer copies its scale
from state at `param5+0x3838` in `FUN_100232d0`; the new producer copies it
from `param5+0x4cf8` in `FUN_1002c220`.

Setter-entry traces identify the live source of this particular difference.
The 2013 Kate app calls `VT_SetPitchSpeedVolumePause_ENG` once with
pitch/speed/volume/pause/speaker values `100/100/200/925/1`. The paired 2006
probe run makes no call to its corresponding setter and retains the observed
100/100/100 row scales. Rewriting only the 2013 setter's volume argument from
200 to 100 changes all three decoded-row calls to scale 100; their source
sample counts are 1,420, 1,514, and 4,470. That three-position output retains
7,258 frames and hashes to
`f904e6c2718fc1be99c5a1b7ca5b2eace4df3223dab4800cae2791bb466693ef`. In the
two-position cutoff case,
setting this API argument to 100 produces a WAV byte-identical to the earlier
row-field 200-to-100 intervention (SHA-256
`87310485389615cdfae8f8779372d9125eb3543d809551cc3b292542433bce53`). The
measured 2x level gap is therefore caused by the different parameter state in
these two callers, rather than a gain reinterpretation of the shared Kate
source rows. This setting difference is separate from the 2013 `>9` selector
cutoff that adds the third position.

The 2006 capture's WAV hash matches the earlier selected-chain control exactly,
so its hardware trace did not alter the audio. Setter traces are
`trace-legacy-hi-volume-setter.gdb` and
`trace-2013-hi-volume-setter.gdb`; the all-row 2013 intervention uses
`trace-2013-hi-volume100-rows.gdb`; the two-position API-level control is
`trace-2013-hi-cutoff9-volume100.gdb`. Evidence is in
`hi-2006-volume-setter/`, `hi-2013-volume-setter/`, and `hi-volume100-rows/`
under `tools/revkit/work/corpus-parity/stage20/`. The two-position API-level
capture is in `hi-query-cutoff9-volume100/` and hashes identically to the
row-field override.

Reproduce the aggregation trace with
`trace-2013-hi-query-aggregation.gdb` and
`compose-key-repacked-attr48-key3-attr40-key2-hi-0404-to-5d00.yaml`. The
sum-9 intervention uses `trace-2013-hi-query-cutoff9-override.gdb`; the
combined cutoff/gain intervention uses
`trace-2013-hi-cutoff9-gain100.gdb`. Captures are in
`hi-query-aggregation-0404-to-5d00/`, `hi-query-cutoff9-override/`,
`hi-query-cutoff9-override-repeat/`, and `hi-query-cutoff9-gain100-v2/` under
`tools/revkit/work/corpus-parity/stage20/`. Each completed runner restored the
Stage 5 input and output hashes. The first gain-probe attempt used invalid
entry-register offsets and made no edits; its output exactly matches the
two-position control and is retained separately in
`hi-query-cutoff9-gain100/`. The paired old-engine row trace uses
`trace-legacy-hi-row-gain.gdb`; its normal WAV hash matches the earlier
`hi-2006-selected-chain/` control, and the capture is in `hi-2006-row-gain/`.
Setter-entry evidence records the 2013 caller passing volume 200 and no
corresponding 2006 setter call; setting the 2013 argument to 100 drives every
decoded row to scale 100. The scripts and captures are
`trace-legacy-hi-volume-setter.gdb`, `trace-2013-hi-volume-setter.gdb`,
`trace-2013-hi-volume100-rows.gdb`, `hi-2006-volume-setter/`,
`hi-2013-volume-setter/`, and `hi-volume100-rows/`.

### 2026-09-29: paired Hi PCM windows and join arithmetic

Read-only hardware-breakpoint captures at 2013 `FUN_1002aac0` and 2006
`FUN_10023b40` saved the two decoded PCM windows immediately before joining
them. The capture used the accepted two-position chain `272820, 272821` and
scale 100 in both engines. Both engines supplied 1,420 samples for row 0 and
5,910 for row 1. Their byte-level PCM inputs are identical across generations:
row 0 SHA-256 `47ba1f42a8bd102e87fad9d66ed42c5589c86ec2fd21a7dada99acc3c7a15823`;
row 1 SHA-256 `c9d05b1fef65ba6b96d2335d9129b2560b45dcf239e66177a0fa78f8592a7130`.
The instrumented WAV hashes also match their uninstrumented controls: 2006
`331fe0de611f49ded95845376adf93c308d15ad3e98abce6968f1d2d2cbb9c2b` and 2013
`87310485389615cdfae8f8779372d9125eb3543d809551cc3b292542433bce53`. Both
outputs contain 7,258 frames. Each runner restored the exact pre-run fixture
hashes.

The old output is a direct copy of row 0 samples `[0, 1348)` followed by a
72-sample join, then a direct copy of row 1 samples `[72, 5910)`. The new
output uses the same interior source samples but also applies boundary
envelopes. Let `T[i]` be the shared 8,192-entry coefficient table; observed
integer results use truncation toward zero:

- New onset, 80 samples: `trunc(row0[k] * T[floor(k * 4096 / 80)])`.
- New overlap, 72 samples: `trunc(row0[1348+k] * T[4096 + floor(k * 4096 / 72)]) + trunc(row1[k] * T[floor(k * 4096 / 72)])`.
- New final envelope, 134 samples: `trunc(row1[5776+k] * T[4096 + floor(k * 4096 / 134)])`.

These formulas reproduce the corresponding 2013 output samples exactly.
The interior copies are also exact: output `[80, 1348)` equals row 0
`[80, 1348)`, and output `[1420, 7124)` equals row 1 `[72, 5776)`. In the old
join, each of the 72 output samples is exactly reproduced by truncating the
sum of two weighted lanes:
`trunc(row0[1348+k] * T[floor((72+k+1)*8191/145)] + row1[k] * T[floor((k+1)*8191/145)])`.
The 2006 routine clamps that combined result; the 2013 routine clamps each
weighted lane before adding the two 16-bit values. No saturation occurs in
this fixture's overlap.

The paired outputs differ at 283 of 7,258 samples. All differences fall in
the 80-sample onset, the 72-sample join, and the 134-sample final envelope;
the largest absolute difference is 508 PCM counts and occurs in the final
envelope. For this matched Hi chain, selector output, decoded source samples,
sample gain, and duration are held equal. The remaining waveform difference
is therefore caused by the generation-specific edge and overlap arithmetic.
These measurements isolate the synthesis contribution in the Hi. control.
They do not establish how much these mechanics contribute to the separate
multi-row Hello defect.

The two trace scripts are `trace-legacy-hi-window-inputs.gdb` and
`trace-2013-hi-window-inputs.gdb`. The raw row buffers, WAVs, runtime logs,
and restoration records are in `hi-window-inputs-2006/` and
`hi-window-inputs/` under `tools/revkit/work/corpus-parity/stage20/`.

### 2026-09-29: matched four-row Hello PCM comparison

The same input-window capture was repeated for natural `Hello.` after applying
the measured continuity-marker overlay and the diagnostic `>2` cutoff on the
2013 side. That run selects the old chain `272822, 272823, 272824, 272825`;
the 2006 engine selects the same rows natively. The 2013 setter volume was
changed to 100 to match the 2006 decoded rows. Counts and decoded PCM bytes
match across all four rows:

| Row | Samples | Leading span (2013) | Trailing span (both) | PCM SHA-256 |
| ---: | ---: | ---: | ---: | --- |
| 0 | 1,210 | 78 | 72 | `8612a32fc5c91453afe09b9f38c5a1a2a996b4d53629d67de4304f927e6e359b` |
| 1 | 578 | 72 | 72 | `7d019679912d21cd2e847f130d6645ee187ea9199fa322b702ead957512ea7df` |
| 2 | 1,180 | 72 | 68 | `6cfc46ad8678768263c296316d8980fac9eafa9401ceb8896809eb84069d5d39` |
| 3 | 4,854 | 68 | 126 | `1bb1f78c80a60368045d258f73f5a604be5881d7fe9a4f6595265e9c4bee17e9` |

Both output WAVs have 7,610 frames. The 2006 hash is the existing baseline
`eb5e399ff403bd86d82cf02a7adf123379b268e24a8fcb40578d763facd58233`; the
2013 matched-chain, volume-100 hash is
`a8ec8d359acf9a2475c3698bc6b5e7ae4d6b41397cd22e842f9d2300c7564fba`.
Their 405 differing samples are confined to the 78-sample onset, the three
overlaps of 72, 72, and 68 samples, and the 126-sample final envelope. The
maximum absolute difference is 1,375 PCM counts. All intervening row interiors
are exact copies of the same source samples in both WAVs.

The 2013 onset and final-envelope equations match the two-row Hi control. Each
of its three internal joins exactly matches the separately truncated 2013
lanes, using the shared table with indices `floor(k*4096/m)` and
`4096+floor(k*4096/m)` for spans `m = 72, 72, 68`. Each 2006 join exactly
matches the combined-lane equation using indices
`floor((m+k+1)*8191/(2*m+1))` and
`floor((k+1)*8191/(2*m+1))`. The copied interiors and all five edge/join
equations were checked against each captured WAV sample by sample.

This matched Hello comparison separates three causes now measured in this
workspace: the native 2013 selector can return a different chain and extra
positions; the paired application runs supply different volume state; and,
even with the same selected rows, row PCM, volume, and duration, the synthesis
paths apply different utterance-edge and overlap arithmetic. The third cause
is local to 405 of 7,610 samples in this four-row control. It cannot by itself
explain duplicated unit content from the native selector path. The capture
therefore identifies synthesis parity work still needed while preserving the
separate selector and caller-state findings.

The paired captures are `hello-window-inputs-2006/` and
`hello-window-inputs/` under `tools/revkit/work/corpus-parity/stage20/`;
traces are `trace-legacy-hello-window-inputs.gdb` and
`trace-2013-hello-window-inputs.gdb`. Both runtime WAV hashes and restoration
records are included with the captures.

### Native cutoff split on the matched Hello rows

To isolate the native `>9` gate from the corrected row mapping, the same
continuity-marker overlay and volume-100 setter were used without the cutoff
override. Native selection returned six IDs:
`272822, 272823, 272823, 272824, 272825, 272825`. The first and fourth
windows match the corresponding four-row control exactly. The two repeated
IDs are overlapping slices of the same PCM rows, not full units appended a
second time:

- `272823`: the native 290-sample row is four-row-control samples `[0, 290)`;
  the native 360-sample row is `[218, 578)`. Their 72-sample overlap is
  byte-identical.
- `272825`: the native 2,112-sample row is `[0, 2112)`; the native
  2,844-sample row is `[2010, 4854)`. Their 102-sample overlap is
  byte-identical.

The native-cutoff and accepted-cutoff 2013 WAVs both contain 7,610 frames.
Their hashes are `21f8c6629bd7e42c7a86033ffe0d36d7f273f0f6d2a92afb42bb7e15b0e14ecf`
and `a8ec8d359acf9a2475c3698bc6b5e7ae4d6b41397cd22e842f9d2300c7564fba`.
Only 165 output samples differ, exactly at the two newly introduced split
joins, and each differs by one PCM count. On this mapped Hello control, the
extra positions are overlapping splits of existing source PCM; the `>9`
cutoff alone does not append a repeated full-unit sound. The native-cutoff
trace and raw slices are in `hello-native-window-inputs/` under
`tools/revkit/work/corpus-parity/stage20/`, with reproduction script
`trace-2013-hello-native-window-inputs.gdb`.

### Equal-gain native Hello comparison with and without continuity mapping

A fresh native-cutoff capture was run on the no-marker matched-context index,
with the 2013 volume setter changed from 200 to 100. This matches the gain in
the marker-corrected native-cutoff capture. The no-marker selector returned
`273369, 273370, 273370, 273371, 264074, 264074`; its six join inputs were:

| Slot | Unit ID | Samples | Leading span | Trailing span | PCM SHA-256 |
| ---: | ---: | ---: | ---: | ---: | --- |
| 0 | 273369 | 1,416 | 78 | 84 | `281d0986fbf67b7f2375c8a1d7f93b387ca3bfc2ce50a00d9eba60cd2249f34e` |
| 1 | 273370 | 648 | 84 | 84 | `b4ed75ee51a009ddcdce2cf0a52d98aae28968c812305763dc250131e24263d9` |
| 2 | 273370 | 520 | 84 | 90 | `89e594ba14632c58ecde5b50e74edec245e91979ce373655a0a0e6f69b90615d` |
| 3 | 273371 | 1,724 | 90 | 86 | `c3aa34b56c14b29df6f31114546be0246307f39fe3ca33abace4cd546b456624` |
| 4 | 264074 | 1,576 | 86 | 106 | `decae6221523d601e776b07b3e4d8f25948b08e412e532ec7a596b0943fbbccf` |
| 5 | 264074 | 2,482 | 106 | 118 | `1c2b300878ae27a8c683c88962e4091775510518d8de0d386c77851bbc287a28` |

The no-marker PCM inputs total 8,366 samples. Subtracting the five prior-row
trailing spans (`84 + 84 + 90 + 86 + 106 = 450`) gives the captured 7,916
frames. Under the marker-corrected native cutoff, the six rows total 7,996
samples and their five prior-row trailing spans total 386, yielding 7,610
frames. Thus the changed chain and its decoded windows add 306 output frames
even though both runs select six positions. Row-by-row raw PCM comparison
found no exact shared substring of 16 or more samples between the no-marker
and marker-corrected inputs; the near-equal output from the matched four-row
control therefore cannot be explained by merely different segmentation of
the same four source buffers.

The equal-volume no-marker WAV hashes to
`09d5463315c63e51acb2b7722b583921c5412c01b6522587ac89890223b54697`; the
marker-corrected native-cutoff WAV hashes to
`21f8c6629bd7e42c7a86033ffe0d36d7f273f0f6d2a92afb42bb7e15b0e14ecf`.
Across their 7,610 common frames, 7,504 samples differ (98.6%); their maximum
absolute sample difference is 24,192. This is a same-engine, same-input-text,
same-volume, same-native-cutoff comparison. Since the marker A/B changes the
selected IDs and the captured source buffers before synthesis, the broad
waveform delta is evidence for a candidate/source-window change, not a
measurement of join arithmetic. The earlier marker A/B identifies setting
bit `0x80` in signature byte `+6` for three legacy-marked `unit-gen2` rows as
the intervention that redirects this selection toward IDs 272822–272825.
That result remains specific to this `Hello.` fixture; the capture does not
establish the semantic identity of every selected unit.

The capture and six raw row buffers are in
`tools/revkit/work/corpus-parity/stage20/hello-no-marker-window-inputs/`.
Its reproduction trace is
`tools/revkit/work/stage20/trace-2013-hello-no-marker-window-inputs.gdb`.
The runner restoration record confirms identical before/after hashes for the
Stage 5 input and output fixtures.

### 2026-09-29: Hi candidate-member crosswalk

The previously open Hi member comparison is now closed for the two 2006
candidate pools captured by `trace-legacy-hi-span-filter-bypass.gdb`. The
original 2006 runtime returns 10 candidates at slot 0 and 9 at slot 1. The
source-index rows for those exact IDs were passed through the `Hi.` suffix-map
repack and the 2013 `FUN_10016ea0` key transform:

| 2006 pool | Captured legacy feature subgroup | Rows | 2013 key under the Hi tail map |
| --- | --- | ---: | --- |
| Slot 0 | one observed subgroup | 10 | `5a 22 01 a0 00` |
| Slot 1 | `22 11 5a 04 22 00 22` | 4 | `22 11 5a 65 00` |
| Slot 1 | `22 11 5a 04 04 00 22` | 5 | `22 11 5a 04 00` before a relaxation-specific remap |

The ten slot-0 rows therefore remain together under the mapped key. The
slot-1 old pool is partitioned into the same observed groups of four and five
rows, but only the four-row subgroup maps to a generated 2013 query key in
the primary Hi overlay. Three controlled repacks separately map the other
group's legacy tail `04 04` to generated target suffixes `5d 00`, `55 00`,
and `4d 00`. Each makes the corresponding exact lookup return one class
entry. All three runs retain candidate-list counts `67, 19, 19`, select
`272820, 272821, 272821`, and produce the same 14,560-byte WAV hash
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.
The repeated output across three target suffixes does not identify which
generated code is the semantic match for the five-row subgroup.

The old span-filter trace keeps only the full-span row at each pool:
`272820` from the ten-row pool and `272821` from the nine-row pool. Their
measured 2006 local costs are `0` and `10.3087`. Under the Hi suffix-map
2013 capture, the candidate lists contain 67, 19, and 15 rows, and the
selected stream is `272820, 272821, 272821`. The transition trace measures
the edge from `272820` to `272821` at `26.709`, versus `34.8761` from
`65331`; the third repeated `272821` adds no transition term in that path.
Thus the mapped old rows are available and the observed first transition
chooses the old pair, while the 2013 path still contains a third position.
Four of the six generated exact query keys still miss in this overlay; those
early misses return before the `>9` score check and route into fallback
expansion. The extra position is therefore upstream of the final path score,
while the ranker decides which row occupies that expanded path. This separates
key recovery from the remaining path-width/position-builder difference for
this prepared Hi case.

The crosswalk is reproducible with
`tools/revkit/work/stage20/crosswalk-hi-class-members.py`; its captured IDs
come from `corpus-parity/stage20/hi-2006-span-filter-bypass/`, and its 2013
target keys and candidate lists come from `hi-2013-suffix-map/`. The
pre-existing fixture caveat remains: the forced P AH0 2013 runtime input is
not a matched old-engine text comparison because the 2006 parser treats its
VTML tag as text. The result closes the source-row membership question for
these captured pools; it does not establish a general conversion for every
2013 relaxation query or repeated-Hello context.

### 2026-09-29: repeated-Hello query-hit membership

The 2013 `Hello hello.` producer capture contains 24 exact lookups. Five
return one class entry; the other 19 return none. An offline scan applied the
captured 2013 key transform to every Kate source row, then enumerated members
for each successful key:

| Lookup calls | 2013 exact key | Source rows | IDs also in the captured 2006 selected chain |
| --- | --- | ---: | --- |
| 7 and 21 | `17 2b 2f 00 00` | 25 | `264073`, `273371` |
| 10 | `2b 30 22 00 00` | 1 | none |
| 16 | `0d 22 17 00 00` | 2 | none |
| 19 | `22 17 2b 00 00` | 9 | `264072`, `273370` |

The first key occurs once in each word instance. Its 25-row source set is the
known 2013 superset for the corresponding legacy five-row class: all five old
members are present, with 20 additional rows. The nine-row key includes the
two old-selected rows shown and has the same nine source members as the
previously measured legacy nine-row class. These successful lookups therefore
show both outcomes: one repeated-Hello class preserves the known legacy
member set, while another keeps all five old members but broadens the pool by
20. The two other hit classes contain one and two rows and none of the nine
captured old-selected IDs.

The row inventory comes from
`tools/revkit/work/stage20/crosswalk-repeat-hello-key-members.py`; the lookup
calls and counts are in
`corpus-parity/stage20/hello-repeat-2013-producer/adapted-producer-keyclasses-gdb.log`.
At this point the 2006 query pools had not yet been captured; the subsequent
capture and crosswalk below close that member-set gap for the recorded
repeated-Hello contexts.

### 2026-09-29: repeated-Hello 2006 pool crosswalk and threshold confirmation

A dedicated original-engine capture records all 24 2006 split invocations
across its three render passes for `Hello hello.`. Each pass repeats eight
distinct query contexts. At each lookup return the trace records the class
shortlist IDs, their per-ID weights, and the corresponding seven-byte feature
records read from `model+0x98`. The feature-to-source crosswalk applies the
2006 lookup tables to every Kate source row and checks that each feature's
source-row count equals its runtime weight before forming a pool.

The resulting member comparison covers all eight distinct old queries:

| 2006 slot | 2006 pool | Matching 2013 source key(s) from the captured exact hits | Member relation |
| ---: | ---: | --- | --- |
| 0 | 75 rows | The old members form generated source keys `5a22041e00` (45 rows) and `5a22171e00` (30 rows); the live target suffix queries miss | No successful-hit class contains these rows |
| 1 | 9 rows | `22172b0000` (call 19, 9 rows) | Exact nine-row match |
| 2 and 7 | 5 rows each | `172b2f0000` (calls 7 and 21, 25 rows) | Both old pools are the same five rows; the 2013 class adds 20 |
| 3 | 1 row | Generated source key `2b30220200`; live hit `2b30220000` has one different row | Equal counts, different member |
| 5 | 73 rows | Eight generated keys across the `0d`, `10`, `1a`, `23`, and `3b` prefixes | No successful-hit class contains these rows |
| 6 | 4 rows | `22172b0000` (call 19, 9 rows) | Old four-row subset; the 2013 class adds five |
| 8 | 6 rows | Generated effective key `2b305a0400`; live targets are `2b305a4500` and `2b305a4d00` | No exact hit before the controlled suffix map |

These comparisons expose both exact parity and two different forms of
candidate mismatch. The 9-row old slot-1 pool exactly matches the 2013
call-19 class. The old slot-6 pool is a four-row subset of that same class.
The old slot-2/7 pools are exact five-row subsets of the repeated 25-row hit
class. Counts alone would hide the slot-3 equal-count mismatch and the
slot-0/slot-5 misses. The reproducible source-row lists and per-key
intersections are printed by
`tools/revkit/work/stage20/crosswalk-repeat-hello-old-pools.py`; the native
feature bytes and weights are in
`corpus-parity/stage20/repeat-hello-old-query-pools/legacy-query-pools-gdb.log`.

The slot-8 suffix was tested with a prefix-conditioned candidate-only
repack. The raw source pair is `04 04`; `04 00` is its effective 2013 key
after byte 6 is masked with `0x20`. A first trial mistakenly used the
effective pair as the source pair and produced no hit. The corrected
`2b305a:0404→4500` treatment maps exactly six source rows:
`1869`, `5170`, `246677`, `255282`, `264074`, and `272825`. The overlay scan
finds six rows under `2b305a4500`, and live lookup call 23 returns one class
entry for that key. The next relaxed key, call 24 `2b305a4d00`, remains empty.
The candidate-key treatment and six-row count are reproducible with
`repack-index-key-fields.py --variant
attr48-key3-attr40-key2-repeat-hello-slot8-0404-to-4500` and
`verify-repeat-hello-slot8-overlay.py`.

A stable 2013 split-helper trace on this overlay records class ID `18857`
with weight 6 at slot 12, metric sum 6, rejection, and a two-position
fallback. It also records the exact nine-row class (`13478`, sum 9) at slot 9;
that check also rejects and falls back. On the old run, the corresponding
six-row and nine-row pools produce no fallback entries, while the one-row
slot-3 pool does. This is direct repeated-input evidence for the cutoff
difference: the same verified member sets of 6 and 9 pass the old selector's
`>2` sufficiency rule and fail the 2013 strict `>9` rule. The captured 2013
builder returns 14 positions with the corrected slot-8 key, so recovering
this lookup alone does not remove the repeated-input split. The runtime
capture, split trace, and restoration records are under
`corpus-parity/stage20/hello-repeat-2013-slot8-0404-to-4500/` and
`hello-repeat-2013-slot8-split-threshold/`; both report identical before/after
Stage 5 input and output hashes.

This result closes the previously missing repeated-Hello query-pool member
comparison for the captured 2006 contexts and confirms the threshold effect
on two shared sets. It does not recover one universal old-to-new suffix
conversion: slot 0 has two old classes under distinct transformed prefixes,
slot 3 has a one-row identity mismatch, slot 5 fans out across eight generated
keys, and the other exact hits are supersets or subsets. Those contexts still
need producer and member-set interventions before any general repack rule can
be claimed.

### 2026-09-29: repeated-Hello 75-row pool and first-chain divergence

The slot-0 pool was used to test whether the 2013 split threshold alone
explains its path difference. Its 75 rows form two 2006 feature classes: 45
rows transform to the 2013 source key `5a22041e00`, and 30 transform to
`5a22171e00`. The first live 2013 query asks for `5a2217a000`. In the raw
source index, the 45-row group is stored under prefixes `5a2205` and
`5b2205`; the 2013 class tables reduce both to `5a2204`. The corrected
candidate-only overlay maps those stored prefixes to `5a2217` and the stored
tail `1e1e` to `a000`. Together with the native 30-row group, this yields the
full 75-row member set at the target lookup.

The split-helper trace reports weight and metric sum 75 at slot 0 and accepts
that slot. The builder returns 13 positions in this overlay; other calls still
fall back at sums 1, 2, and 9. The separate slot-8 overlay returned 14
positions, but these are individual trials, not a cumulative or one-variable
comparison. The candidate pool confirms the cutoff effect while showing that
it is not enough to recover the old selected chain. Runtime output is 29,752
bytes with SHA-256
`0d3a956bb565bd69ae1af49c0118f751cf6edd436b3f826b20b9e6f8af943889`; Stage 5
input and output fixture hashes match before and after.

The first selected unit still differs under native 2013 ranking. The 2006
chain selects `273369`; the corrected 2013 alias selects `149799`. Both are in
the 75-row pre-score list. A high-precision local-score capture gives each
`9.60000038147`, with initial rank key `-1`. The final 2013 transition row
contains 30 candidates, excludes `273369`, and selects `149799` at index 6.
The old query-pool crosswalk assigns `273369` to 2006 class ID `52165` (weight
2) and `149799` to class ID `52164` (weight 45). This shows that the first
divergence occurs after lookup and the strict threshold: exact feature-score
equality does not guarantee equal retention in the capped 2013 candidate
list.

A causal ordering intervention changed only `273369`'s pre-sort node key from
`-1` to `-2`. The 2013 selected index then becomes 0 and the selected unit is
`273369`; the output changes to 30,800 bytes with SHA-256
`b572e5a7d45a9a245f0a76f116c9a980ff4638bfe5477f5eab26c431cfb207a7`. This
proves that candidate ordering and retention control this first selection
under the tested overlay. It does not establish why the unmodified sort
permutes equal-key rows in that exact order, or whether the 2006 shortlist
ordering policy explains the original preference. The established 2013
sorter uses partition swaps and can permute equal keys; a broader ordering
comparison is still needed before treating this as a general generation
rule.

One initial alias attempt mistakenly treated the effective 2013 class prefix
`5a2204` as stored source bytes, so it mapped only the 30-row native-prefix
class. Its split trace was then overwritten by the corrected capture because
the GDB script retained the first run's absolute log path. The first run's
result artifacts remain in its directory, but that directory is not a complete
capture bundle; the dedicated corrected split trace is under
`corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-final/`. The
corrected candidate indexes and candidate/backtrack traces are under
`tools/revkit/work/stage20/index-adapter-key-repacked-attr48-key3-attr40-key2-repeat-hello-slot0-pool-to-a000-v2/`
and the sibling `hello-repeat-2013-slot0-pool-to-a000-v2/` capture directory.
Exact local scores are in
`hello-repeat-2013-slot0-pool-to-a000-v4/adapted-score-precision-gdb.log`.

### 2026-09-29: repeated-Hello 73-row slot-5 pool alias

The old slot-5 pool's 73 source rows split across ten raw-prefix/tail
combinations and eight effective 2013 keys. The live 2013 call-16 target
`0d22170000` has only two native rows, IDs `4515` and `114916`, stored under
raw prefixes `2f221e` and `2f2218`. A candidate-only overlay maps all ten old
raw-prefix/tail combinations to raw prefix `2f221e` and maps tails `0a0a` and
`1414` to `0000`. The offline crosswalk and exact mapping are reproducible with
`analyze-repeat-hello-slot5-keys.py` and the
`attr48-key3-attr40-key2-repeat-hello-slot5-pool-to-0000` repack variant.

Under that overlay, 2013 reports class weight 79 at split slot 7 and accepts
it. The old 73 source rows have joined the two native call-16 rows; the higher
weight reflects the preserved 2013 metric weights, not a 79-row old pool.
Other split slots still fall back at the observed sums 1, 9, and on queries
without a logged metric class. The position builder returns 13 positions.
Output is 30,872 bytes with SHA-256
`a3078cb23fcd7392dfa599e19fd0031bd67b84e312d365ff48977ea8968d3fe8`, and the
Stage 5 fixture hashes match before and after.

The selected chain capture reports unit `273369` at backtracked context 0,
matching the first unit in the 2006 chain, but this does not establish full
chain parity. Slot 0 still falls back in this isolated overlay, and fallback
changes the position segmentation; context numbers after it cannot be aligned
one-for-one with old positions. The slot-5 test establishes that key
aggregation and the `>9` check can be restored for this pool, while leaving
the effect on later selected units and edges open. Its split and backtrack
captures are under `corpus-parity/stage20/hello-repeat-2013-slot5-pool-to-0000/`
and `hello-repeat-2013-slot5-pool-to-0000-backtrack/`.

### 2026-09-29: combined repeated-Hello pool overlay and cutoff result

A combined candidate-only index variant now applies the observed slot-0,
slot-5, and slot-8 key aliases together. Runtime split instrumentation reports
the slot-0 class at sum 75 and accepts it, and the mapped slot-5 class at sum
79 and accepts it. The builder returns 12 positions. The slot-8 six-row
candidate class is present at split slot 10, but the native 2013 `>9` check
rejects sum 6 and takes the two-position fallback. Other observed failures in
this trace are sum 1 at slot 4 and sum 9 at slot 7. This separates the effects:
the aliases recover those lookup populations, while the 2013 threshold still
splits the six- and nine-weight cases. The differing call-slot numbers reflect
the accepted and expanded positions accumulated earlier in the same builder.

The backtracked 12-position chain begins `149799, 273370, 273370, 273371`;
later rows include `208339, 34312, 49874` before the repeated tail rows. The
first choice is still not old-selected `273369`: the combined overlay leaves
the observed 30-row cap and 2013 scoring/ranking behavior in place. Thus
restoring lookup membership and accepting the 75-weight class do not restore
the 2006 path. Output is 26,400 bytes, SHA-256
`602ce927b25b6e511513285148617feb76c08102fb9f7f797e1d30a90450cd07`; both
Stage 5 fixture hashes match before and after the run.

The first-position trace locates the remaining divergence after candidate
lookup. `273369` and `149799` both receive exact local cost
`9.60000038147`, and both begin with the same `-1` continuity key. In the
75-entry candidate list, `273369` is at index 21 and `149799` at index 63;
after scoring and the 30-row reduction, only `149799` remains. An intervention
that changes only `273369`'s continuity key to `-2` makes it the first selected
unit, proving that the retained order is causal for this first choice. Static
disassembly of 2013 `FUN_1001b5f0` explains why equal-key order is not stable:
the large-range path uses pivot partitioning, stops its scans at equal values,
and swaps crossing entries; small ranges use an insertion-like pass. The
exact deterministic permutation depends on the input order and partition path.
This is direct evidence for 2013's ordering mechanism, but the captured 2006
shortlist is a sequence of class IDs and weights rather than a same-format flat
candidate array, so a one-to-one old/new sort-order claim remains unsupported.

The combined overlay, split trace, and backtrack trace are reproducible with
`attr48-key3-attr40-key2-repeat-hello-slots0-5-8-pools` and
`trace-repeat-slots0-5-8-split.gdb` / `trace-repeat-slots0-5-8-backtrack.gdb`.
The capture directory is
`corpus-parity/stage20/hello-repeat-2013-slots0-5-8-pools/`. Remaining work is
to alias the slot-3 singleton and the slot-1/2/6/7 classes without unintended
collisions, then compare complete matched member sets and post-prune edge
costs. The present result establishes the independent cutoff and rank-cap
effects; it does not yet explain the full 2006 versus 2013 chain divergence.

### 2026-09-29: slot-3 identity mismatch control

The old slot-3 singleton maps through the candidate conversion to effective
key `2b30220200`, while the corresponding 2013 query class uses
`2b30220000`. A further disposable variant changes only the old singleton's
`0202` suffix to `0000`, joining it to the native row. The runtime metric sum
therefore rises from 1 to 2. The 2013 split helper still rejects it at slot 4
under its strict `>9` condition. The builder remains at 12 positions, the
backtracked 12-unit chain is unchanged, and the WAV remains byte-identical at
26,400 bytes / SHA-256
`602ce927b25b6e511513285148617feb76c08102fb9f7f797e1d30a90450cd07`.
Stage 5 input and output fixture hashes again match before and after.

This eliminates the slot-3 suffix mismatch as a cause of this fixture's
position-count difference: mapping it restores the expected two-row lookup
population but does not cross either engine's observed threshold (`>2` in the
2006 selector, `>9` in the 2013 selector). It also leaves the selected 2013
path unchanged in this overlay. The 2006 trace already takes its fallback at
this one-row position, so the mismatch is real but not explanatory for the
old/new split divergence here. Evidence is in
`corpus-parity/stage20/hello-repeat-2013-slots0-3-5-8-pools/`; the variant is
`attr48-key3-attr40-key2-repeat-hello-slots0-3-5-8-pools`.

An order-only intervention was also applied to this combined overlay. It sets
the first-position node key for `273369` from `-1` to `-2`, making that row the
first selected unit. The rest of the selected sequence is
`273370, 273370, 273371, 208339, 34312, 49874, 264072, 264072, 264073, 264074,
264074`, while the 2006 chain continues `273370, 273371, 232670, 266023,
264071, ...`. The builders return 12 and 9 positions respectively, so this is
a sequence comparison and does not assume index-for-index phone alignment.
The intervention output is 27,448 bytes, SHA-256
`184a7d1df5aa92b7822479b7a23612e633f382abf1e02a6e76902c2d14d19b7b`; fixture
hashes match before and after. This bounds the equal-rank/cap effect to the
first choice in this run; later candidate and edge differences remain. The
capture is under
`corpus-parity/stage20/hello-repeat-2013-slots0-3-5-8-order-override/`.

### 2026-09-29: completed path rows and the late shortlist intervention

The completed `FUN_10018c80` rows resolve the apparent conflict between the
matched-edge trace and the backtracked path. The edge trace showed the
individual minimum-cost path ending at each candidate; selection backtracks
from later rows, so the cheapest partial path at one context need not belong
to the globally selected chain. In the controlled 2013 overlay, context 4
stores `232670` at 79.1433 and `208339` at 80.6948, both through `273371`.
Context 5 stores `266023` at 98.4015 through `232670`, and `34312` at 99.3981
through `208339`. At context 6, `264071` is present in the 79-row ranked pool
but absent from the 30-row transition set passed to the path scorer; `49874`
is retained at 132.377 through `34312`. The following `264072` row is 158.995
through `49874`. Thus the old `232670 → 266023 → 264071` route is cut off
before path scoring, even though its partial costs at contexts 4 and 5 are
lower. The captured 2013 baseline path follows
`208339 → 34312 → 49874`.

A correction to the earlier field interpretation is necessary: `81.4368` was
read from the pre-score candidate node, before the final `FUN_100182e0` unit
scorer. The scorer returns `59.8` for both `264071` and `49874`. In the native
order, the second sort receives 79 candidates with no preserved prefix; both
rows have final sort key `59.8`. The bounded 30-row result includes `49874`
and excludes `264071`, showing that equal-score ordering affects shortlist
membership.

The earlier ranking-key intervention at position 6 changes `264071`'s
pre-sort key from `-1` to `-2`. It moves to the head of the initial sort and
forms a one-row preserved prefix; after local scoring, the untouched key
remains `-2`, so it survives the final cutoff. The completed rows then record
`264071` at 140.651 through `266023`, followed by `264072` at 158.802 through
`264071`. Although the path ending at `264071` is 8.274 higher than the
baseline path ending at `49874`, the next transition makes the 264071 route
lower by 0.193 at `264072`. Backtracking consequently restores the legacy
subsequence `232670 → 266023 → 264071` at contexts 4–6. Later 2013 contexts
still contain expanded handoffs, so this is not a full sequence or WAV match.

A cleaner control leaves every candidate field unchanged and swaps only the
`264071` node pointer from index 22 to index 0 immediately before the local
scorer. The preserved prefix remains zero; both target rows still score
`59.8`, and the final 30-row list now retains `264071` at index 11. Its WAV is
byte-identical to the key-promotion run (26,374 bytes, SHA-256
`41d793131f8e57596ad123718bcd17249a781c6e955b0ffce37c46f206a58b8f`), while
the unmodified-order baseline is 27,448 bytes. This isolates the immediate
cause more cleanly: equal-score tie ordering plus the 30-row cutoff can remove
the legacy-selected unit before path scoring; downstream continuity then
selects the old subsequence when that row is retained. It does not yet explain
why the native 2013 candidate ordering differs from the 2006 path's candidate
reduction.

The promoted run exits normally and produces 26,374 bytes (SHA-256
`41d793131f8e57596ad123718bcd17249a781c6e955b0ffce37c46f206a58b8f`); the
controlled order-only baseline produces 27,448 bytes. Both runs restore the
Stage 5 input and output fixtures to their recorded hashes. Completed-row
captures and scripts are `corpus-parity/stage20/hello-repeat-2013-completed-
path-rows/` with `trace-repeat-slots0-3-5-8-completed-path-rows.gdb`, and
`corpus-parity/stage20/hello-repeat-2013-promote-264071/` with
`trace-repeat-slots0-3-5-8-promote-264071.gdb`. The sorter boundary captures
are `corpus-parity/stage20/hello-repeat-2013-shortlist-sort-boundary/` and
`hello-repeat-2013-promoted-shortlist-sort-boundary/`, with scripts
`trace-repeat-native-shortlist-sort-boundary.gdb` and
`trace-repeat-promoted-shortlist-sort-boundary.gdb`. The field-preserving
pointer-swap capture is `corpus-parity/stage20/hello-repeat-2013-order-only-
264071-shortlist/`, from `trace-repeat-order-only-264071-shortlist.gdb`.

### 2026-09-29: matched-unit span and rank comparison

A native 2006 trace at `FUN_1001cff0` captures the same candidate IDs at old
position 5. Candidate `264071` has left reach 0, right reach 3, total span 4,
and pre-score rank key -8; its local score is 7.304. Candidate `49874` has
span 1 and score 23.5254. The 2006 shortlist has 73 candidates and a three-row
preserved prefix; both target units survive its 30-row cap, with `264071` at
index 0. This confirms that the old selector also caps rows at 30. The
difference is the evidence already attached to this target before that cap.

At the aligned repeated-Hello 2013 context 6, both IDs are in the 79-candidate
pool, but both have span 1, weighted coverage 0, and final local score 59.8.
The native second sort has no preserved prefix; it retains `49874` and drops
`264071`. Reordering only the `264071` pointer before that sort admits it,
which proves the 2013 cutoff is sensitive to equal-score ordering. These
traces locate this fixture's immediate divergence: the old span calculation
gives `264071` a coverage-derived rank advantage, while the 2013 calculation
leaves it tied in a large span-one group where the 30-row cutoff removes it.

The marker representation difference was then tested on the exact old unit
chain. In the original 2006 runtime, units `264071`, `264072`, and `264073`
map to phone IDs 3434, 32400, and 25383, and each has continuation byte `1`
in the separate model `+0x68` array. Unit `264074` has byte `0`. In the
matched 2013 runtime, unit `264071` has signature
`1,47,34,30,13,0,0`; the final byte consumed by `FUN_100230a0` has bit `0x80`
clear. Setting only that bit in process memory before span construction
changes the 2013 candidate from span 1 and score 59.8 to span 3 and score
54.8667. Its pre-sort key changes from -1 to -3, the preserved prefix becomes
one row, and `264071` is retained at the head of the 30-row list. The output
is 26,374 bytes with SHA-256
`41d793131f8e57596ad123718bcd17249a781c6e955b0ffce37c46f206a58b8f`, exactly
matching the earlier key-promotion and pointer-order interventions. Fixture
hashes match before and after.

This confirms that the missing 2013 continuation marker is causal for two
additional positions of coverage and for restoring `264071` to the shortlist
in this run. Setting bit `0x80` on all three consecutive 2013 records
`264071`–`264073` raises the target span to 6, lowers its score to 49.9333,
changes its key to -6, and retains it first. The output remains byte-identical
to the one-bit treatment. Thus the direct selection failure is established:
the 2013 model rows used here have continuation markers clear where the
corresponding 2006 chain has them set. Span 6 does not numerically match old
span 4 because the builders emit 12 versus 9 positions. The converter source
resolves why the marker is missing in this adapted run: its declared 2013
signature order is `attr48, key5[:3], attr40, key5[3:]`. That puts old
`attr_48=1` in signature byte 0, while `FUN_100230a0` checks bit `0x80` in
signature byte 6. The active repack therefore preserves the old field value
in a slot that 2013 does not use for continuation.

A disposable disk overlay applied the field conversion
`signature[row*7+6] |= 0x80` for unit IDs `264071`–`264073`, whose source
`attr_48` values are all 1. Runtime memory then showed signature byte 6 `128`
for `264071`, span 6, score 49.9333, and first place in the shortlist. The
completed path selected the old subsequence at aligned contexts 4–6:
`232670 → 266023 → 264071`; context 7 then selected `264072` through
`264071`. The controlled output changed from 27,448 bytes and hash
`184a7d1df5aa92b7822479b7a23612e633f382abf1e02a6e76902c2d14d19b7b` to
25,096 bytes and hash
`7c5137c896e850f067d2f9686ef496b2a7ecf2bfd1406c6ee4cb3040e44d6b35`.
Both runs restore the Stage 5 input and output fixtures. This confirms the
repacker field placement as the immediate cause of this 2013 continuity and
shortlist failure. It does not prove that marker correction alone restores
full audio parity: the 2013 builder still emits 12 positions against the old
builder's 9, and other lookup-key and split-threshold differences remain.

The captures are `corpus-parity/stage20/hello-repeat-2006-shortlist-sort-
boundary/legacy-shortlist-sort-gdb.log`,
`corpus-parity/stage20/hello-repeat-2006-unit-boundary/`,
`corpus-parity/stage20/hello-repeat-2013-span-boundary/`,
`corpus-parity/stage20/hello-repeat-2013-shortlist-unit-signatures/`, and
`corpus-parity/stage20/hello-repeat-2013-bit7-264071-span/`,
`corpus-parity/stage20/hello-repeat-2013-bit7-chain-span/`, and
`corpus-parity/stage20/hello-repeat-2013-bit7-disk-overlay/`. The corresponding
trace scripts are `trace-legacy-repeat-shortlist-sort-boundary.gdb`,
`trace-legacy-repeat-unit-boundary.gdb`,
`trace-repeat-2013-span-boundary.gdb`,
`trace-repeat-2013-shortlist-unit-signatures.gdb`, and
`trace-repeat-2013-bit7-264071-span.gdb`,
`trace-repeat-2013-bit7-chain-span.gdb`, and
`trace-repeat-2013-bit7-disk-shortlist.gdb` with
`trace-repeat-bit7-disk-completed-path-rows.gdb`.

### 2026-09-29: selector-side legacy continuation patch

A disposable copy of the Kate 2013 DLL was patched at both continuation
scans in `FUN_100230a0`. Each scan now begins at signature byte `+0` and tests
bit `0x01`, matching the adapted legacy `attr_48` representation, instead of
beginning at byte `+6` and testing `0x80`. The original Kate DLL remains
unchanged. With the unmarked matched-context index overlay and the natural
`Hello.` fixture, this DLL selected the same six source IDs as the prior
index-side marker treatment:
`272822, 272823, 272823, 272824, 272825, 272825`. Its 15,264-byte WAV has
SHA-256 `936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`,
matching that marker treatment exactly.

A separate runtime intervention at the return from `FUN_10024060` accepted
the two observed nonempty sums, 9 and 6, under the old `>2` rule. Combined
with the selector-side DLL patch, it selected exactly
`272822, 272823, 272824, 272825`; the 15,264-byte WAV hash was
`215ff65cf15e37f4dcce58527511aa31f143b58a1dca865ef082252a3a83e5f3`,
matching the earlier marker-plus-cutoff control. Both runs restored Stage 5
input and output fixtures byte-for-byte.

We also changed two `FUN_10023e70` comparisons from `>=10` to `>=3` in a
separate DLL copy. That direct binary edit had no effect on this fixture: it
still selected six handoffs and produced the marker-only hash. That edit
targeted the wrong path. The active ordinary cutoff is now located and patched
directly below. These results support a targeted Kate compatibility patch,
not general old-voice compatibility. Other class-key and candidate-membership
mismatches remain.

The hash-guarded builder is
`tools/revkit/work/stage20/make-kate-selector-legacy-continuity-dll.py`; the
native and cutoff-intervention captures are under
`tools/revkit/work/corpus-parity/stage20/hello-selector-byte0-patch-native/`
and `hello-selector-byte0-patch-cutoff-override/`. The ineffective direct
cutoff candidate is captured under
`hello-selector-byte0-cutoff3-dll/`.

### 2026-09-29: active cutoff, query-key, and member-set map

The ordinary 2013 whole-position path is now located through its active
machine-code decision. `FUN_10024060` first asks `FUN_10018770` to build the
query variants and collect class IDs. An empty class list returns false before
metric scoring, immediately causing `FUN_100242a0` to expand the position into
two halves. For a nonempty list, it prunes the list, calls `FUN_10023060` to
sum the selected class metrics, then applies this decision:

| Condition | Native 2013 result |
| --- | --- |
| target class byte is 8 and sum > 0 | accept |
| otherwise, sum >= 10 | accept |
| otherwise | split into two positions |

In Kate `vt_kat.dll`, the ordinary test is at VA `0x1002428d`: `cmp eax,0xa`
followed by `setge`. With the established nonnegative integer metric sums,
this is `>=10`, equivalent to the previously described `>9`. The class-byte-8
positive-sum branch is at `0x10024277`–`0x10024280` and remains separate.
`FUN_10023e70` is a later relaxed/half-position helper; changing its two
comparisons did not change the natural Hello path because that is not the
ordinary whole-position decision.

A hash-guarded DLL copy changed the active immediate from 10 to 3 in addition
to the already tested byte-0 continuation reads. No GDB return-value hook was
used. On the matched-context Hello fixture it selected
`272822 → 272823 → 272824 → 272825`; the 15,264-byte WAV hash
`215ff65cf15e37f4dcce58527511aa31f143b58a1dca865ef082252a3a83e5f3` exactly
matches the earlier GDB cutoff intervention. Stage 5 input and output fixture
hashes were identical before and after. This proves that the located cutoff
instruction is responsible for the observed 9- and 6-sum fallback in this
controlled run. It does not repair empty-key lookups or guarantee legacy
candidate ranking in other contexts. The disposable DLL SHA-256 is
`5037eb66665fa7e91c5e96967686617080343531c2697bef46785e372b567e00`.

The data path and currently measured incompatibilities map as follows:

```text
2013 tree leaf / seven-byte signature
  -> FUN_10018770 emits context-dependent variants
  -> FUN_10016ea0 projects a signature into a five-byte class key
  -> FUN_10019450 exact-searches the sorted class table
  -> returned class IDs and metric weights are pruned/combined
  -> FUN_10023060 computes the selected metric sum
  -> FUN_10024060 accepts (special class 8 and sum > 0) or (sum >= 10)
  -> empty key or insufficient sum: FUN_100242a0 two-position expansion

2006 tree leaf / feature record
  -> FUN_1001b2d0 translates feature fields through legacy tables
  -> FUN_1001da80 emits context/suffix key variants
  -> legacy class IDs and weights are accumulated
  -> old sufficiency gate accepts a selected sum > 2
```

The five-byte class lookup is exact; key relaxation is performed by emitting
more keys, not by fuzzy comparison inside `FUN_10019450`. Therefore lookup
failure and candidate-member mismatch are separate cases. For the eight
repeated-Hello contexts, the 2006 pool crosswalk currently shows:

| Old query context | Key/member status in the measured 2013 table |
| ---: | --- |
| 0 | 75 old rows form two transformed keys (45 + 30); the live target key misses until prefix and suffix aliases are applied |
| 1 | 9-row old pool exactly matches a 2013 class |
| 2 and 7 | each old 5-row pool is a subset of a 25-row 2013 class |
| 3 | both sides have one row, but the 2013 hit is a different source row (`2b30220200` versus `2b30220000`) |
| 5 | 73 old rows span eight generated keys; the live target has only two native rows before aliases |
| 6 | old 4-row pool is a subset of the 9-row class also used by context 1 |
| 8 | six old rows require a context-specific `0404 → 4500` tail map to hit the six-row class |

These are not reducible to one suffix substitution. The 2013 producer copies
signature byte `+5` and masks byte `+6` with `0x20`; the 2006 producer divides
its terminal fields by 10 and emits table-directed variants. The observed
stable Hello prefix conversion includes `48 → 47` in one context. Candidate
classes can be absent, exact, supersets, or equal-sized but disjoint, so a
successful lookup alone does not establish equivalent candidates. Full
eight-context member evidence and reproducible crosswalk scripts are in the
preceding “repeated-Hello 2006 pool crosswalk and threshold confirmation”
section.

The active-cutoff DLL builder is
`tools/revkit/work/stage20/make-kate-selector-legacy-continuity-dll.py` with
`--legacy-active-cutoff`; its output is mounted by
`compose-selector-legacy-continuity-active-cutoff3-dll.yaml`. Runtime capture
and restoration records are under
`tools/revkit/work/corpus-parity/stage20/hello-selector-byte0-active-cutoff3-dll/`.

### 2026-09-29: natural Hi query and member-set control

The same key-miss mechanism was checked using the ordinary `Hi.` fixture
(SHA-256 `c7e13072a15d4bec0f88822de02bc825d5c6055ac8cf928913d301e336a28c1e`),
so this comparison does not depend on the forced P AH0 markup. The original
2006 engine emits two accepted lookups: `90 34 17 30 30` returns 10
candidates, and `34 17 90 04 34` returns 9. Its selected pair is
`272820 → 272821`, repeated across the three render passes.

On the corrected-repack 2013 baseline, all six generated exact keys miss:
`5a 22 01 a0 00`, `5a 22 01 98 00`, and
`22 11 5a {65,5d,55,4d} 00`. The existing prefix-conditioned suffix overlay
maps the two source populations to the primary keys. The resulting trace has
one hit for `5a 22 01 a0 00`, one for `22 11 5a 65 00`, and misses for the
other four variants. 2013 then selects
`272820 → 272821 → 272821`: it recovers the first two old row identities but
retains one extra handoff. The output is 14,560 bytes with SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`.
The native-cutoff 2013 baseline is 13,640 bytes with SHA-256
`8cedb6fcdf0e5948c6cca81e785e3a6592a4346671217bac3d399e27d4c3b685`. The
original 2006 output is 14,560 bytes with SHA-256
`331fe0de611f49ded95845376adf93c308d15ad3e98abce6968f1d2d2cbb9c2b`.
Both runtime records show unchanged before/after shared fixture hashes.

The extra handoff is now causally tied to the active whole-position cutoff.
For the mapped `22 11 5a 65 00` key, the matched class is ID `13468` with
metric weight 4. The active `FUN_10024060` gate compares against 10, rejects
that sum, and invokes the two-position fallback; the resulting candidate
lists have weights 67, 19, and 15, and the builder returns three positions.
This is distinct from the other `5d 00` experiment, where weights 4 and 5
combine to 9. A cutoff-only disposable DLL changes the active immediate from
10 to 3 while preserving native continuity handling. The builder then returns
two positions, with candidate weights 67 and 4, and selects exactly
`272820 → 272821`, matching the old chain. Its 14,560-byte output has SHA-256
`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`.
An independent GDB intervention that forces only the sum-4 decision to pass
produces the identical WAV hash and selected IDs. The DLL intervention
therefore confirms the cutoff cause without changing the continuity code; the
runtime intervention independently confirms the branch decision.

The 67-member first class contains all ten old slot-0 members, plus 57
additional rows. The second-position membership mismatch had a narrower cause:
the five old `feature_b` rows map to suffix `04 00`, while the current natural
Hi overlay only aliases suffix `04 22` to `65 00`. A separate
`04 04 → 65 00` alias brings those five rows into the same second class as the
four already recovered `feature_a` rows. Under the cutoff-3 DLL, the resulting
2013 second candidate list is exactly the old nine IDs:
`4087, 4102, 84094, 177718, 243540, 272821, 279498, 279499, 280449`.
The first list remains 67, with the old ten included among the 57 extras. The
two selected IDs and WAV remain unchanged (SHA-256
`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`). This
alias restores the second pool's exact membership but does not yet establish
first-pool parity or general query-key compatibility.

A candidate-prefix control then isolated the first-pool collision. In the
unfiltered 67-row class, the ten old members use raw keys
`5a22111e1e` or `5b22111e1e`; the 57 extras use seven other raw prefixes:
`5a2202`, `5a2203`, `5a220e`, `5a220f`, `5b2202`, `5b220e`, and `5b2212`,
all with tail `1e1e`. The 2013 key tables collapse those prefixes onto the
same `5a2201a000` query class. A disposable candidate-key overlay redirects
those seven non-target prefix/tail pairs to `5a2213`, while preserving the
exact old ten-member key group. Combined with the `04 04 → 65 00` second-pool
alias and cutoff 3, the live 2013 candidate lists are exactly the old 10 and 9
unit IDs, and the selector chooses `272820 → 272821`. The result remains
14,560 bytes with SHA-256
`c902797173c44a136c352b099b6c8eeac59a5ff7b09da3caa73b8a9457c64c4f`.
Stage 5 input/output hashes also match before and after. This shows that the
observed Hi pool mismatch can be corrected with prefix-conditioned candidate
rewrites; it does not establish that this hand-authored map is the correct
general transformation for other input contexts or words.

The seven exact prefix/tail patterns match 57 rows total across the scanned
Kate `gen`, `gen2`, and `etc` indexes; every one is among the 57 extra rows in
this Hi class. No matching rows occur in `num` or `alp`. This bounds the
observed collision to those records in the inspected model snapshot, though
other query contexts could still use their rewritten signatures.

The active gain path was then checked on the same fixture with that sum-4
intervention and an entry-only rewrite of the 2013 volume setter argument from
200 to 100. It still selects `272820 → 272821`; the 14,560-byte output hash is
`87310485389615cdfae8f8779372d9125eb3543d809551cc3b292542433bce53`, matching
the earlier equal-gain capture. Against the original 2006 WAV, both files are
mono PCM16 at 16 kHz with 7,258 frames; only 283 samples differ, with changes
localized to the 80-sample onset envelope, the 72-sample overlap, and the
134-sample tail envelope. The matched source rows themselves decode to
byte-identical PCM buffers across generations. The detailed paired-window
capture already documents the respective join arithmetic in
[the isolated Hi split and PCM gain analysis](#hi-split-threshold-and-pcm-gain-isolated).
Thus the remaining natural-Hi waveform delta after matching IDs is in gain
handling and synthesis edge/join behavior, rather than a different selected
source row.

The logs and scripts are
`tools/revkit/work/corpus-parity/stage20/hi-plain-2006-query-evidence-v2/`,
`hi-plain-2013-query-evidence/`, and
`hi-plain-2013-tailmap-query-evidence-v2/`,
`hi-plain-2013-cutoff9-override/`,
`hi-plain-2013-tailmap-cutoff3-only/`,
`hi-plain-2013-tailmap-cutoff3-featureb-alias/`,
`hi-plain-2013-tailmap-cutoff3-exact-pools/`,
`hi-plain-2013-metric4-override/`, and
`hi-plain-2013-metric4-volume100/`, with traces
`trace-legacy-hi-plain-query-classes.gdb`,
`trace-2013-hi-plain-keyclass-lookup.gdb`, and
`trace-2013-hi-plain-tailmap-keyclass-lookup.gdb`,
`trace-2013-hi-plain-metric-cutoff.gdb`,
`trace-2013-hi-plain-metric4-override.gdb`, and
`trace-2013-hi-plain-metric4-volume100.gdb`. The cutoff-only patch is built by
`tools/revkit/work/stage20/make-kate-tree2-active-cutoff-dll.py` and mounted by
`compose-kate-tree2-active-cutoff-dll.yaml`.

The exact second-pool membership check uses the generated
`attr48-key3-attr40-key2-hi-0404-to-6500` variant from
`repack-index-key-fields.py`, mounted by
`compose-key-repacked-attr48-key3-attr40-key2-hi-0404-to-6500.yaml`.
The exact-pool candidate-prefix control is the
`attr48-key3-attr40-key2-hi-exact-first-pool-filter` variant, mounted by
`compose-key-repacked-attr48-key3-attr40-key2-hi-exact-first-pool-filter.yaml`;
`inspect-hi-first-pool-members.py` lists the original 67-row class's ten old
members and seven extra key-prefix families.

### 2026-09-29: Hi candidate-prefix rule on held-out Hello

The seven-prefix candidate rewrite was then tested without the Hi-specific
suffix aliases on ordinary `Hello.`. Both runs use the same Kate repack layout,
the same selector DLL with the active cutoff and legacy continuation reads,
and the same fixture. The only change is the seven-entry
`candidate_prefix_tail_map`. The base candidate counts at positions 0, 1, and
6 are 465, 111, and 23; with the Hi prefix rewrite they are 408, 117, and 19.
The selected sequence is unchanged:
`273369, 273369, 273370, 273370, 273371, 282025, 282025`. Both WAVs are
17,836 bytes with SHA-256
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`.
Stage 5 fixture hashes match before and after each capture.

The two GDB query traces contain the same ten generated signature/key/count
lines byte-for-byte. Those traces establish equal query inputs, exact keys, and
returned counts; they did not record the returned class identities. A paired
trace at `FUN_10023350` dumps the class shortlist and each class's unit members
before unit expansion. A further trace at `FUN_10023f90` records the result
classes for each query variant. The returned list sizes remain equal, but the
class identities and unit membership change:

| Position context | Same query key | Baseline class result | Prefix-map treatment |
| ---: | --- | --- | --- |
| 0 | `5a2217a000` | Class key `5a22011e00`, 67 units | The same class key has the intended 10 units; expanded count 465→408 |
| 1 | `5a2217a000` | Class key `2922171400`, 1 unit (`243343`) | Class key `3722171400`, 7 units; class key `5a22171e00` keeps its 30 units while its numeric ID changes `32928`→`32929` |
| 6 | `2b305a4500` | Class key `16305e0500`, 5 units | Class key `41305e0500`, 1 unit (`28799`) |

All three use the relaxed feature-scope path through `FUN_10023e70` and
`FUN_10023c70`, called with 10-class and 10,000-unit limits. An earlier reading
mistakenly treated the ten-class result as the input to `FUN_10023c70`. The
paired `FUN_1002df50` range traces show the actual inputs: 84 classes for the
context-1 query `5a2217a000` and 66 for context 6 query `2b305a4500`. Each
baseline/treatment pair has the same set of `(five-byte class key, member
count)` entries; only traversal order differs. A per-class `FUN_10023a70`
trace also shows identical integer distance scores for every key in both
scopes.

The changed output classes are tied at the ten-class cutoff. In context 1,
keys `2922171400` and `3722171400` both score 620. Their input positions
reverse from 74/75 in the baseline to 75/74 in the treatment; the scope has
one class scoring 560 and 35 scoring 620, while only nine 620-score classes
fit after the best-scoring class. The paired post-sort traces directly show
rank 1 changing from `2922171400` to `3722171400`. In context 6, keys
`16305e0500` and `41305e0500` both score 1210, and their input positions move
from 46/44 to 45/46. The scope has four classes scoring below 1210 and nine tied at 1210,
so six of those nine can occupy the remaining shortlist positions. A paired
post-sort trace directly captures context 6: rank 9 is `16305e0500` in the
baseline and `41305e0500` in the treatment, exactly matching the changed
10-class result. Their member counts are 5 and 1, respectively, well below
the 10,000-unit cap.

The static `FUN_1001b5f0` disassembly shows an unstable partition sort: equal
scores can be swapped as the partition scans cross. The runtime trace confirms
that this order dependence decides both changed winners: context 1 at rank 1
and context 6 at rank 9. Thus the held-out cross-context substitutions are
not caused by different query keys,
missing classes, or changed per-class distance scores. Repacking changes the
feature-range traversal order, and the 2013 top-10 sort can select a different
member of a tied score group. Position 0 remains a separate, intended effect:
its matching class shrinks from 67 members to 10 under the seven-prefix map. These differences
are present before `FUN_10023350` expands units and before `FUN_100230a0`
computes continuity metadata or local scoring. The selected Hello units and
rendered waveform remain stable in this capture, but the rewrite is not
context-local and should not be treated as a general adapter rule. The same
virtual address `0x1001b5f0` in the 2006 DLL contains unrelated code, so the old
engine's tie/cutoff path must be identified through its own call graph rather
than by matching the 2013 sorter address. Paired class and unit traces are under
`hello-hi-builder-classes-baseline-repacked/` and
`hello-hi-builder-classes-treatment-repacked/` in
`tools/revkit/work/corpus-parity/stage20/`; both runs exited normally, restored
the Stage 5 fixtures, and produced the same 17,836-byte WAV hash
`da2b244d649242574621d7239ad3c42b283879eae6341a3ee1a5236f1d6ac13e`. The
trace scripts are `trace-2013-hello-hi-prefix-builder-classes-baseline.gdb`,
`trace-2013-hello-hi-prefix-builder-classes-treatment.gdb`, and their shared
`trace-2013-hello-hi-prefix-builder-common.gdb`. The per-query scripts are
`trace-2013-hello-hi-prefix-query-membership-baseline.gdb`,
`trace-2013-hello-hi-prefix-query-membership-treatment.gdb`, and their shared
`trace-2013-hello-hi-prefix-query-membership-common.gdb`; captures are under
`hello-hi-query-membership-baseline-v2/` and
`hello-hi-query-membership-treatment-v2/`.

The full feature-range and per-class distance logs are under
`hello-hi-range-scope-baseline-v2/`, `hello-hi-range-scope-treatment-v2/`,
`hello-hi-scope-distances-baseline/`, and
`hello-hi-scope-distances-treatment/`. The paired context-1/context-6 post-sort
captures are under `hello-hi-scope-sorted-baseline-v6/` and
`hello-hi-scope-sorted-treatment-v6/`; their GDB entry points are
`trace-2013-hello-hi-prefix-scope-sorted-common.gdb` and its baseline/treatment
wrappers. The distance trace entry points are
`trace-2013-hello-hi-prefix-scope-distances-common.gdb` and its baseline and
treatment wrappers.

A diagnostic offline ordering by `(integer distance score, class key)` yields
the same ten class keys for both repacks in each context. This shows that a
stable tie-break can remove this repack-dependent difference, but it is not yet
a compatibility fix: in context 1 it would change the baseline shortlist too.
The old selector helper path is now identified. In the old `FUN_1001e470`
builder, each ordinary position calls `FUN_1001dfb0`; a false result alone
invokes the two-position `FUN_1001e110` fallback. The ordinary helper expands
legacy key variants through `FUN_1001da80`, deduplicates the returned class
IDs, and accepts when their accumulated model weights exceed 2. At 31 or more
unique IDs it first reduces the list to 30 using `FUN_1001d5f0`/`FUN_1001d530`;
below 31 it copies every ID directly. The matched 2006 `Hello.` trace returns
pools of 4, 2, 1, and 1 IDs, with weight sums 75, 9, 5, and 6. All four take
the copy-all arm and are accepted, so old ranking does not participate in this
threshold result. The runtime record is `split-threshold-old-gdb.log` under
`tools/revkit/work/corpus-parity/stage20/`.

The old sorter is `FUN_100140c0`, reached through `FUN_1001d5f0` for large
candidate pools and also used by later unit ranking. Its pseudocode has a
short-range insertion pass, a partition path with strict score comparisons
and crossing swaps, and a heap fallback for highly imbalanced partitions. It
has no class-ID tie-break. Therefore “stable 2006 sort versus unstable 2013
sort” is not supported: both generations can reorder equal-score entries on
their partition paths. The existing old runtime sort-boundary capture at row
5 records 73 unit candidates and the resulting top-30 order; pseudocode
establishes the partition behavior for class ranking as well.

This narrows the generational difference. The demonstrated repack-sensitive
2013 top-10 substitutions happen when equal-distance classes straddle the
cutoff of `FUN_10023c70` feature-scope ranking. The old `Hello.` whole-row
decision instead receives small, key-expanded pools and accepts them before
the 31-entry ranking boundary. For the measured extra 2013 positions, empty
ordinary key lookups trigger fallback before the `>=10` score test, so lowering
that gate alone cannot recover those positions. Candidate-key production,
class membership, and selector-specific marker conversion are confirmed gaps
in the tested contexts; unstable sorting is a secondary sensitivity when a
tied pool reaches a rank cutoff, not the cause of empty lookups.

The archived 2006 query traces inspected here report at most 20 IDs from any
single class lookup (the maximum is in the P AH0 capture); this does not prove
that every possible context stays below 31, but the old 31-ID ranking arm is
not established by these runtime samples. The directly comparable ranking
boundary in the repeated-Hello case is instead the later unit shortlist: the
2006 path ranks 73 rows to 30, and the 2013 path ranks 79 rows to 30. For unit
`264071`, the old runtime supplies span 4 and rank key -8, placing it first;
2013 supplies span 1, zero weighted coverage, and a tie at local score 59.8,
then drops it at the 30-row cutoff. A pointer-order-only intervention admits
the row and restores the old selected subsequence, proving that the cutoff is
causal once the score inputs are tied. More decisively, translating the legacy
continuation flag to the byte and bit consumed by the 2013 selector raises the
row's span and restores that subsequence. For this repeated-Hello path, the
primary failure is an upstream continuation-field interpretation that changes
coverage and rank; the unstable sort determines which row survives only after
that bad input has made it tie with the competing unit.

Static sources are `full-analysis-2006/all-functions-pseudocode.c` at
`FUN_1001e470`, `FUN_1001dfb0`, `FUN_1001e110`, and `FUN_100140c0`. Function
names are Ghidra pseudocode labels. A future probe of the old 31-ID branch
still needs a naturally occurring context with at least 31 accumulated
deduplicated IDs. That branch is separate from the already demonstrated
73/79-row unit shortlist mismatch and should not be used to explain it.

The original audio/query captures are under
`hello-hi-filter-unmapped-baseline/` and
`hello-hi-prefix-filter-isolated/` in
`tools/revkit/work/corpus-parity/stage20/`; the treatment index is generated
by `repack-index-key-fields.py --variant
attr48-key3-attr40-key2-hi-first-candidate-prefix-filter`. The matching query
traces are `hello-hi-prefix-filter-query-keys-baseline.log` and
`hello-hi-prefix-filter-query-keys-treatment.log`; all ten generated
signature/key/count lines are identical. The paired scripts are
`trace-2013-hello-hi-prefix-filter-query-keys.gdb` and
`trace-2013-hello-hi-prefix-filter-query-keys-treatment.gdb`.

### 2026-09-29: projected-key expansion and membership inflation

The source-row crosswalk was extended to measure what happens if a 2013
selector emits every exact class key containing a member of each captured
2006 pool. For each old class feature, the crosswalk first validates its
source-row count against the captured 2006 weight, then projects all source
rows through the current 2013 signature tables and counts rows by resulting
five-byte key. Results are `projected bucket size / old-pool rows in bucket`:

| 2006 context | Old rows | Projected 2013 key buckets | Extra rows from unioning those buckets |
| ---: | ---: | --- | ---: |
| 0 | 75 | `5a22041e00` 45/45; `5a22171e00` 30/30 | 0 |
| 1 | 9 | `22172b0000` 9/9 | 0 |
| 2 and 7 | 5 each | `172b2f0000` 25/5 | 20 each |
| 3 | 1 | `2b30220200` 1/1 | 0 |
| 5 | 73 | Eight keys: `0d22040a00` 1/1; `1022040a00` 27/27; `1022171400` 12/2; `1a22040a00` 1/1; `1a22171400` 25/1; `2322040a00` 19/19; `3b22040a00` 21/21; `3b22170a00` 1/1 | 34 |
| 6 | 4 | `22172b0000` 9/4 | 5 |
| 8 | 6 | `2b305a0400` 6/6 | 0 |

This separates two failure modes. Some old pools can be reconstructed by
querying and unioning their projected keys without adding rows, but the live
2013 producer does not emit those keys for every target context (notably
contexts 0, 3, and 8). Other old pools are strict subsets of broader 2013
classes. Most significantly, old context 5 spans eight projected keys, and
unioning those current buckets adds 34 rows outside the old 73-row pool. The
targeted slot-5 rekey experiment shows that this pool can still be assembled by
mapping its ten raw-prefix/tail combinations to one query key; live 2013 then
returns 75 rows (the 73 legacy members plus two native rows) with weight sum
79. The intervention accepts the position and selects `273369` at its first
backtracked context, but other positions still split and the complete chain
does not match. So the mismatch is recoverable for selected contexts with a
candidate-specific rekey, while generic projection-key union can overfetch
enough to alter ranking. A reusable adapter still needs to preserve legacy
membership boundaries across contexts, not just make the exact lookup hit.

The calculations are an offline crosswalk over the captured 2006 feature
records, original mapping tables, 2005 source indexes, and current 2013 key
projection. They establish projected bucket membership and overfetch, not the
semantic meaning of the feature bytes or the effect of every overfetched row
on the final path. Captured live lookups independently confirm the projected
membership for keys `172b2f0000` and `22172b0000`; the other projected buckets
remain offline results. Reproduce with:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage20/crosswalk-repeat-hello-old-pools.py
```

### 2026-09-29: isolate `attr_b` from signature byte `+4`

The previous corrected transfer test both moved legacy `attr_40` into 2013
`attr_b` and zeroed the reconstructed signature byte `+4`. A new paired
overlay instead duplicates `attr_40` into `attr_b` and preserves all seven
reconstructed signature bytes. Its builder checks that the only changed bytes
in each adapted index are the `attr_b` column. The captured target keys and
first three ranker candidate lists are unchanged from the zero-fill control.

The first eight `FUN_100182e0` returns are exactly the same in the duplicate
and move runs: `16.032, 30.032, 0.0705882, 0.00851064, 15.0082, 16.0348,
16.96, 15.1455`. The zero-fill control returns `25.6, 39.6, 9.6, 9.6,
24.6, 25.6, 25.6, 24.6`. All three controls select the same final-rank lists
and produce byte-identical four-position output, SHA-256
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`.
Fixture input/output hashes match before and after. This isolates the changed
local costs to the value supplied in `attr_b`; removing signature byte `+4`
was not necessary to produce those cost changes, and the cost change alone
does not fix the selected chain in this cutoff-controlled case.

Static access agrees with that paired result on the measured path. The 2013
class-key producer consumes signature offsets `+1,+2,+3,+5,+6`; the traced
`FUN_100182e0` local scorer reads candidate signature offsets `+1,+2,+3,+5`
and flag bits at `+5`, but not candidate byte `+4`. This bounds the earlier
reconstruction: byte `+4` can be retained without affecting these key and
local-score decisions. Whether other 2013 paths consume it remains open, as
does the semantic mapping between legacy `attr_40` and new `attr_b`. The
experiment supports neither treating `attr_40` as a proven `attr_b` field nor
discarding it; it shows only that `attr_b` values causally alter local costs
and that the zero/move/duplicate variants did not alter this final chain.

Reproduce the duplicate overlay and runtime trace with:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 -B tools/revkit/work/stage20/duplicate-legacy-attr40-to-attrb.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-matched-context-attrb-duplicate.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-matched-context-score-zero.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/matched-context-score-duplicate \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

### 2026-09-29: neutralize only the `attr_b` score term

The 2013 scorer reads the candidate `attr_b` array through `model + 0x60`,
looks up the pair cost against the target byte, and retains that float in
`[ebp+0x1c]`. Disassembly shows two paths. At `0x1001864a`, the flagged path
normalizes the cached value directly; at `0x10018736`, the ordinary path adds
it after combining the other feature distances. A GDB overlay zeros only that
cached float at the corresponding point in either branch. Index bytes, lookup
tables, signatures, query keys, candidate membership, and cutoff behavior are
unchanged.

For matched `Hello.`, the first eight local costs change from `25.6, 39.6,
9.6, 9.6, 24.6, 25.6, 25.6, 24.6` to `16, 30, 0, 0, 15, 16, 16, 15`. The
three captured final-rank candidate lists remain identical, as do the four
selected units `273369, 273370, 273371, 255282`. Baseline and neutralized WAVs
are byte-identical: 18,364 bytes, SHA-256
`dbb567894041a543b42593b838f36ea8db060c786a6e806ecbbb87ea02918318`.
Fixture restoration hashes match. Thus `attr_b` changes local costs but does
not cause the remaining selected-chain difference on this controlled input.
Its effect on other contexts, especially ranking ties and candidate cutoffs,
remains untested.

Reproduce with the disposable trace script:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/compose.yaml \
  -f tools/revkit/work/stage19/compose-versioned.yaml \
  -f tools/revkit/work/stage20/compose-vector-overlay.yaml \
  -f tools/revkit/work/stage20/compose-key-repacked-attr48-key3-attr40-key2-matched-context-tail-map.yaml run --rm \
  -e INPUT_FIXTURE=/work/stage20/hello-plain.txt \
  -e GDB_SCRIPT=/work/stage20/trace-matched-context-score-zero-attrb-neutralized.gdb \
  -e EVIDENCE_DIR=/work/corpus-parity/stage20/matched-context-score-zero-attrb-neutralized \
  runtime /bin/bash /work/stage20/run-adapted-forced-pah0.sh
```

### 2026-09-29: `attr_b` changes the repeated-Hello path

The same isolated score-term hook was applied to the repeated-Hello run with
the pool aliases, continuation-bit overlay, vector-tree repair, and first-row
rank intervention held fixed. This time suppressing `attr_b` changes the
selected path. The baseline selects unit `232670` at output rows 4 and 15;
with only the cached `attr_b` score set to zero, those rows become `205629`.
Rows 5 and 6 remain `266023` and `264071`. The WAV changes from 25,096 bytes
(SHA-256
`7c5137c896e850f067d2f9686ef496b2a7ecf2bfd1406c6ee4cb3040e44d6b35`) to
24,470 bytes (SHA-256
`31d866f64d7fdac99252c34ae38dea0dac33b9761fe4199e3fd778d23770b749`).

The completed-path row orders match through contexts 1–3. Context 4 is the
first divergence in the captured 30-row transition list: eight baseline IDs
drop out and eight others enter after neutralization. Backtracking then picks
`205629` instead of `232670`. This confirms that the `attr_b` score term can
change shortlist membership and final selection on repeated speech even when
the stored `attr_b` column is all zero. The function scales the looked-up
distance by a candidate-dependent duration factor, so a constant attribute
does not produce a constant per-unit score. In this case the zero-filled term
helps preserve the native-selected row; removing it moves away from that
choice. That observed benefit does not establish that zero is the intended
legacy mapping.

A second overlay copied legacy `attr_40` into `attr_b` while preserving the
seven-byte signature and continuation marker. This selected `53333` at rows 4
and 15, while retaining `266023` and `264071` at rows 5 and 6; its WAV is
24,154 bytes (SHA-256
`b501f256fdad3be777011356a503e3bfb917913c4d55c0f2cd6065e47a16dd1f`). Thus
copying `attr_40` also changes the path but does not recover the baseline
choice. The three-way result rejects the simple claims that `attr_b` is
irrelevant or that `attr_40` is already proven to be its replacement. The
proper version-2005 behavior remains unresolved; its adapter must avoid
assuming that zero fill is score-neutral.

The neutralized and `attr_40`-duplicate captures both restore Stage 5 fixtures
byte-for-byte. Reproduce the neutralized run with
`trace-repeat-bit7-disk-attrb-neutralized.gdb` and the existing bit-7 compose
overlay. The `attr_40`-duplicate overlay is built by
`duplicate-repeat-hello-attr40-to-attrb.py`, mounted by
`compose-repeat-hello-attrb40-bit7.yaml`, and traced with
`trace-repeat-bit7-disk-attrb40-duplicate.gdb`. Captures are under
`corpus-parity/stage20/hello-repeat-2013-bit7-disk-overlay-attrb-neutralized/`
and `hello-repeat-2013-bit7-disk-attrb40-duplicate/`.

### 2026-09-29: later transition rows on shared unit IDs

A native 2006 trace was added for current units `232670`, `266023`, and
`264071`. The first 2013 matched-row trace used the earlier overlay without
the disk-level continuation-bit correction. Under that condition it recorded
15 predecessor evaluations for `232670`, 25 for `266023`, and none for
`264071`; the absent `264071` edges were therefore conditional on that
candidate set, not a general property of 2013 pair scoring.

The matched-row trace was repeated with the disk-level bit-7 continuity
overlay, while retaining the GDB intervention that promotes `273369` at the
first 2013 position. Its `Hello.` output is the same 25,096-byte WAV as the
bit-7 disk-overlay baseline (SHA-256
`7c5137c896e850f067d2f9686ef496b2a7ecf2bfd1406c6ee4cb3040e44d6b35`), and the
Stage 5 input and output fixtures were restored byte-for-byte. In this
corrected trace, the 2013 scorer evaluates 20 predecessors for `232670`, 28
for `266023`, and 19 for `264071`. The 2006 counts are 5, 22, and 20.

Joining by unit ID gives five shared predecessors for `232670`; 2006 ranks
`273371` first in that intersection, while corrected 2013 ranks `264073`
first (transition terms `7.45827` and `6.96275`, respectively). For
`266023`, 12 of the old and corrected-new predecessor IDs are shared, and
both rank `232670` first (terms `18.977` and `8.49863`). For `264071`, 16
predecessors are shared and both rank `266023` first; the transition term is
`12.3499` in both captures.

The corrected 2013 completed-path backpointers independently confirm those
three locally best rows: `264073 → 232670 → 266023 → 264071 → 264072` at
contexts 3–7. For the traced middle rows, their stored predecessor indices
point to `264073`, `232670`, `266023`, and `264071`, respectively. Thus the
shared-edge minima for `232670`, `266023`, and `264071` agree with the actual
selected chain in this corrected run; this does not assert that the preceding
2006 state or the entire path is identical. The native 2006 selected-chain
trace gives `273371 → 232670 → 266023 → 264071 → 264072`: the path agrees
from `232670` through the following rows, but the predecessor entering
`232670` differs (`273371` versus `264073`). This directly confirms the
corrected score comparison's upstream divergence and later shared chain.

The comparison extends one row beyond `264071`, to current `264072`. Native
2006 scores 30 predecessors at context 6; corrected 2013 scores 25 at context
7, with 10 predecessor IDs shared. Both rank `264071` first in that
intersection, with transition term `0` in both. The path totals (`57.7644`
and `127.809`) remain on engine-specific scales. This confirms the shared
`264071 → 264072` handoff in both selected paths, while the earlier predecessor
of `232670` remains the observed divergence. The additional captures use
`trace-repeat-2006-next-row-edge.gdb` and
`trace-repeat-bit7-disk-next-row-edge.gdb`; their evidence directories are
`corpus-parity/stage20/hello-repeat-2006-next-row-edge/` and
`hello-repeat-2013-bit7-next-row-edge/`.

The context ordinals differ (`232670`: 3 versus 4; `266023`: 4 versus 5), and
the engines accumulate different prior costs. These are shared-unit and
shared-predecessor comparisons, not identical target-state replays; cumulative
totals must not be compared as a common scale. The corrected overlay restores
`264071` to the scored set and preserves the 2006 best shared predecessor for
that row. It also changes the best shared predecessor for `232670`, showing
that continuation-bit handling changes upstream path costs and candidate
edges, not only the target row's shortlist membership. This still does not
identify the meaning of the continuation bit outside the demonstrated
repeated-Hello rows or establish full-engine parity.

### Why the corrected 2013 path enters `232670` through `264073`

Paired 2013 pair-score captures hold the first-position `273369` rank
intervention and all pool aliases fixed, changing only the disk-level
continuation-bit overlay. For current `264073` at context 3, the `264072`
predecessor edge changes as follows:

The prior `6.0899` reduction is itself accounted for by two earlier rows.
At context 1, current `264072` from `252497` has the same predecessor cost
(`4.8`) and edge feature (`8.50575`) in both runs, while its local term falls
from `11.1502` to `8.07509`; total cost falls from `24.4559` to `21.3808`.
At context 2, current `264072` from `264072` has feature `0` in both, while
its local term falls from `11.0296` to `8.01478`; with the prior reduction,
the cumulative cost falls from `35.4855` to `29.3956`. The marker correction
therefore changes the path cost before the `264073` row, rather than the
`232670` pair scorer inventing a new preference at the end.

| Overlay | Previous path | Edge feature | Penalty | Local term | Total |
| --- | ---: | ---: | ---: | ---: | ---: |
| Marker uncorrected | 35.4855 | 2 | 0 | 35.1225 | 72.608 |
| Disk bit-7 corrected | 29.3956 | 0 | 0 | 17.5612 | 46.9568 |

The 25.6512 cost reduction decomposes into 6.0899 from the predecessor's
prior path, 2 from the edge feature, and 17.5613 from the local term (subject
to displayed rounding). The same correction lowers the `273370 → 264073` edge
from 89.8269 to 72.2656: its prior path, feature (15.7059), and penalty (10)
stay fixed while the local term falls from 35.1225 to 17.5612. It therefore
remains more expensive than the corrected `264072 → 264073` route.

At the next row, current `232670`, the edge feature and local term are nearly
unchanged for both candidate predecessors. The corrected 2013 engine selects
`264073` because its accumulated cost is now 46.9568, versus 66.3776 for
`273371`; without the marker correction, the predecessor-path costs are
72.608 for `264073` and 66.3776 for `273371`, yielding incoming edge totals
85.3475 and 79.1433. This directly links the on-disk continuation marker to
the upstream backpointer change. It establishes the mechanism for this
controlled repeated-Hello path, not universal old/new score parity or a full
semantic definition of the marker.

The paired captures are under
`corpus-parity/stage20/hello-repeat-2013-unmarked-prefix-232670/` and
`hello-repeat-2013-bit7-prefix-232670/`. Reproduction uses
`trace-repeat-predecessors-of-232670-unmarked.gdb` with the unmarked
`...-slots0-3-5-8-pools.yaml` overlay and
`trace-repeat-predecessors-of-232670-bit7.gdb` with the corresponding
`...-pools-bit7-continuity.yaml` overlay. Both retain the same first-row rank
intervention and restore the Stage 5 fixtures.

Reproduce the 2006 capture with
`trace-repeat-2006-matched-late-edges.gdb` and the original-MSI Compose
overlay; its Stage 5 input and output fixture hashes match before and after.
The corrected 2013 capture uses
`trace-repeat-bit7-disk-matched-edges.gdb` and the disk-level continuity
Compose overlay. Run `compare-common-transition-edges.py` to join the
captures by unit ID and report intersection winners. The captures are under
`corpus-parity/stage20/hello-repeat-2006-matched-late-edges/`,
`hello-repeat-2013-matched-old-row-edges/`, and
`hello-repeat-2013-bit7-matched-late-edges/` (the middle path is the
uncorrected control).

### 2026-09-29: selector-side continuation read on Hi and Apple

The selector-only DLL patch was checked on natural `Hi.` and `Apple.` with
matched tree2 DLL controls. The patched DLL and the Stage 19 tree2 control DLL
have the same 688,128-byte size and differ at exactly four file offsets:
`0x23155`, `0x23170`, `0x2320a`, and `0x23218`. Those are the two span scans'
signature-byte offset and marker-mask immediates. The tree2 path bytes are
identical. The patch therefore changes only the source byte and mask used for
the two continuation reads in this comparison.

For `Hi.`, both runs used the Hi tail-map index, versioned trees, vector
overlay, and native cutoff. The control and selector-patched DLL both selected
`272820 → 272821 → 272821`, produced 14,560 bytes, and had SHA-256
`3deb8c81f2151c7fc9bdf605965aaf9d03df1782eb6d0fe2e08af9d34dd60547`. The
three-position result and candidate counts `67, 19, 15` were unchanged. The
old source rows have `attr_48` values 1 for `272820` and 0 for `272821`. Thus
the mapped flag's presence on the first selected row is not by itself enough
to change this short Hi output; the helper's neighbor matches and candidate
context also determine whether the flag affects coverage.

For `Apple.`, both runs used the same Apple prefix/tail-map index, versioned
trees, vector overlay, input fixture, trace script, and native cutoff. The
tree2 control DLL selected five callback rows (`177774, 177774, 64719, 82036,
82037`) and produced a 13,212-byte WAV with SHA-256
`167a388e30f74e83e6d955b39208fb9d62acd6063a3d833c12e0e12c6fef1e23`. Changing
only the two continuation reads initially produced a 16-callback trace and a
25,090-byte WAV. That result did not reproduce: later runs with the same
hardware-only trace returned five callbacks and a 14,248-byte WAV. The cause
of this run-to-run discrepancy is not established, so the 16-callback result
is retained as an anomalous capture and is not used to characterize the
patch.

### Apple marker-patch coverage crosswalk and reproducibility correction

A matched coverage trace was run with hardware breakpoints in both DLL
variants. Both runs report five positions and the same candidate counts by
context: `113, 34, 50, 20, 158`. The control selected
`177774 → 177774 → 64719 → 82036 → 82037` and produced 13,212 bytes (SHA-256
`167a388e30f74e83e6d955b39208fb9d62acd6063a3d833c12e0e12c6fef1e23`). The
selector patch selected `177774 → 177774 → 64719 → 2554 → 2555` and produced
14,248 bytes (SHA-256
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`). The
input and output fixture hashes were restored in both runs.

The first context's selected-row coverage is unchanged. From context 1 onward,
the patched candidate metadata differs: among the captured rows, candidates
with span greater than one rise from 1 to 6 in context 1, 0 to 5 in context 2,
0 to 20 in context 3, and 0 to 20 among the first 60 of 158 candidates in
context 4. Context 3's complete 20-row list has span 1 in the control and span
2 after the patch. The position count and candidate-list sizes stay fixed, so
the reproducible effect is a coverage change and different late selected
rows; it does not explain the unreproduced 16-callback observation as extra
positions.

The old 2006 neighbor-link trace contains both late alternatives: candidate
`2554` and candidate `82036` at position 3 link to predecessor `189455`, while
`2555` and `82037` at position 4 link to `2554`. Thus the patched `2554 →
2555` pair is supported by observed old candidate links, but it differs from
the unpatched 2013 selection and is not by itself evidence of old selected-path
parity. The hardware-only crosswalk captures are under
`apple-selector-byte0-crosswalk-hw-control/` and
`apple-selector-byte0-crosswalk-hw-patch/`; reproduce both with
`trace-apple-continuity-crosswalk.gdb`. The earlier software-breakpoint
crosswalk is retained separately and is not used for the hardware-only
comparison. The simple hardware-only patch trace's repeated capture is under
`apple-selector-byte0-patch-repeat/`.

Hardware-breakpoint DP traces locate the late switch after coverage metadata
is attached. Both runs enter the final position with 30 candidates, but the
preceding context's retained transition pool is 16 rows in the control and 19
after the patch. In the control, backtracking selects `82036` at context 3
(cumulative cost `32.1952`) and `82037` at context 4 (`40.2424`), with
`82037` pointing to `82036`. In the patched run, it selects `2554` at context
3 (`31.0598`) and `2555` at context 4 (`36.0851`), with `2555` pointing to
`2554`. The patched context-4 transition list contains `2555` at index 16 and
`82037` at index 0; the former's predecessor is `2554`, the latter's is
`82036`. The control's 30-row context-4 list contains `82037` but not `2555`.
These values explain the selected 2013 alternative within each run; they do
not make cumulative costs directly comparable across the two runs.

The new link is consistent with the old neighbor-link trace's observed
`2554 → 2555` edge. The controlled DLL change therefore affects coverage,
which in turn changes the transition candidate pool and late backpointer
choice without changing the five-position builder result or initial query
pool sizes. This explains the selector behavior for this Apple fixture, not
why the 2006 engine ultimately prefers its own complete path or whether the
patched WAV is more correct. The DP captures are under
`apple-selector-byte0-path-costs-hw-control/` and
`apple-selector-byte0-path-costs-hw-patch/`; reproduce with
`trace-apple-path-costs-hardware.gdb`.

The Hi control remains unchanged: patched and control runs selected
`272820 → 272821 → 272821`, returned candidate counts `67, 19, 15`, and
produced the same 14,560-byte WAV. Across these inputs, the byte-read patch
changes some coverage and selection decisions but does not establish that the
mapping is universally correct. Query-key production, class membership,
cutoff behavior, ranker differences, and synthesis remain separate parity
questions. The selector DLL is generated by
`make-kate-selector-legacy-continuity-dll.py` and mounted by
`compose-selector-legacy-continuity-dll.yaml`.

### Native 2006 selected records versus patched 2013 backpointers

A fresh native 2006 run used the same plain `Apple.` fixture and the original
Kate MSI overlay. The hardware breakpoint at `0x10013cf0` captured the
selected-record callback sequence three times, once per render pass:

`177774 → 177774 → 100180 → 59558 → 59559`.

The runner exited normally and restored the Stage 5 input and output hashes.
The old WAV is 14,340 bytes (SHA-256
`1ccabfc08b9ddbdf8cda2bd2a1f5d79b716649dbffe7c857572cb903bdbaa119`). The
capture is under
`corpus-parity/stage20/apple-selector-byte0-old-selected-records-2026-09-29/`.

This direct selection trace settles the comparison that the earlier
neighbor-link capture could not: native 2006 selects `100180, 59558, 59559`
in the last three positions; the patched 2013 DLL selects
`64719, 2554, 2555`. The 2006 and patched 2013 runs agree on the first two
IDs only. The old selected rows are nevertheless present in the patched
2013 transition lists. In that 2013 run, context 2 gives `100180` cumulative
cost `20.4021` from predecessor `177774`, while selected `64719` costs
`18.8205` from the same predecessor. Context 3 retains `59558` at `34.0749`
from `100180`; selected `2554` costs `31.0598` from `64719` (and `82036`
costs `32.1952` from `64719`). Context 4 retains `59559` at `39.1103` from
`59558`; selected `2555` costs `36.0851` from `2554`. These are within-run
2013 cumulative values, not comparable old-engine scores.

For this fixture, the late divergence therefore cannot be explained solely
by the three old-selected unit IDs being absent from the patched 2013
transition lists. It occurs at or downstream of candidate scoring and
backpointer choice in the 2013 path. The native DP cost and predecessor trace
now identifies the old path directly, and a paired component trace locates
why the two objectives rank the shared alternatives differently.

### Legacy local-score component on the same Apple candidates

A second hardware trace sampled `FUN_1001c860` for selected old rows and the
same competing IDs. Across all three render passes, the component returns
were stable. On context 2, it scores old-selected `100180` at `0` and
`64719` at `0.0765957`. On context 3, it scores old-selected `59558` at
`0.323902`, versus `82036` at `1.06144` and `2554` at `1.3899`. These
context-2/3 local comparisons favor the old-selected rows, consistent with
the 2006 callback chain and different from the patched 2013 backpointer
choice.

Context 4 prevents treating that component as the selector: the same old
score is `0.247923` for `82037`, `0.307186` for selected `59559`, and
`0.377395` for `2555`. The 2006 engine still emits `59559`; the path-cost and
backpointer evidence below shows that its accumulated transition objective
selects `59558 → 59559` despite those local-score values. The new capture
and restore record are under
`corpus-parity/stage20/apple-selector-byte0-old-candidate-scores-2026-09-29/`;
the trace script is `trace-apple-old-selected-candidate-scores.gdb`. The
runner exited normally and restored the recorded Stage 5 fixture hashes.

The context-4 shortlist boundary trace sees 157 candidates with a 20-row
presorted prefix. `59559`, `82037`, and `2555` all survive the returned
top-30 list at indices 5, 16, and 19 respectively. Thus these three rows are
present before the dynamic-programming pass; their later cost/backpointer
ordering accounts for the selected record, not a missing shortlist member.
The capture is under
`corpus-parity/stage20/apple-selector-byte0-old-shortlist-2026-09-29/`.

The old scorer's consumer and selected-record handoff are now traced through
`FUN_1001e470`, its pair-cost updates, and predecessor writes. The remaining
implementation target is to map the observed legacy context weights and
context-4 fast path into a production-compatible 2013 scoring adaptation.

### Apple selector objective comparison and runtime mediation

The 2006 trace at `FUN_1001ec07` captures the dynamic-programming update;
`FUN_1001ec44` records each candidate's best predecessor, and the selected
indices reach the callback at `0x1001ed7f`. For the old chain, context 2
selects `100180` at cumulative cost `15.3267` after predecessor `177774`.
Context 3 selects `59558` through `100180` at `27.5946`; alternatives
`2554` through `64719` and `82036` through `64719` cost `29.9268` and
`31.6901`. Context 4 then selects `59559` through `59558` at `27.9018`.
The competing `2555` through `2554` and `82037` through `82036` paths cost
`30.3042` and `31.9381`. These are measured within-engine path costs and
backpointers, not estimates from the local scorer.

The patched 2013 path uses the opposite late ordering. It selects
`64719 → 2554 → 2555` at costs `18.8205`, `31.0598`, and `36.0851`; the
shared old rows `100180 → 59558 → 59559` cost `20.4021`, `34.0749`, and
`39.1103`. Thus the final 2013 margin is `3.0252` in favor of `2555`, while
the 2006 margin is `2.4024` in favor of `59559`.

The paired edge trace shows that several distance inputs match numerically
for the same unit pairs, but the context weights and local terms differ. In
canonical `(raw, derived-A, derived-B)` order, measured 2006 weights are
`(10,1,5)` at context 2 and `(10,2,10)` at context 3. The corresponding
patched 2013 weights are `(10,2,5)` at both contexts. At context 2 this
raises the 2013 feature subtotal for `100180 → 177774` from `6.84252` to
`7.32373`, chiefly from the doubled derived-A weight. At context 3 it lowers
the subtotal for `64719 → 2554` from `14.2417` to `10.477`, chiefly from
halving derived-B's weight, while the subtotal for `100180 → 59558` changes
from `11.9439` to `11.7516`. The local terms also differ between generations.
At context 4, the old engine returns transition term zero for the three
measured adjacent-ID pairs; the 2013 scorer adds subtotal `2` even though
their measured distance columns are zero. The legacy disassembly contains an
adjacent-unit zero-cost branch, consistent with these returns. These
observations identify concrete selector-objective mismatches for this fixture;
the units and their metric inputs remain present in both transition lists.

A runtime-only intervention replaced the 2013 pair-feature subtotal on the
eight aligned context-2, context-3, and context-4 edges with the corresponding
measured 2006 values. It left 2013 local scores and every other transition
unchanged. The selected sequence changed to
`177774 → 177774 → 100180 → 59558 → 59559`, matching the native 2006
selected IDs. The intervention's WAVE is 14,340 bytes, SHA-256
`4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8`; the
native 2006 WAVE has the same byte length but SHA-256
`1ccabfc08b9ddbdf8cda2bd2a1f5d79b716649dbffe7c857572cb903bdbaa119`.
This isolates the transition subtotals as sufficient to recover the selected
chain on this input, not to recover identical synthesis output.

Partial interventions that changed only the three measured context-3 edges,
or those edges plus the two context-2 edges, did not select the old chain.
Both instead selected `177774 → 177774 → 64719 → 86327 → 86328` and produced
13,412-byte WAVs. This confirms that candidate competition across the
remaining path and the context-4 updates still matters; the complete
context-2-through-4 substitution, not the context-3 change alone, was
sufficient in the tested intervention.

The native path-cost and feature-component captures are under
`corpus-parity/stage20/apple-selector-byte0-old-path-costs-2026-09-29/` and
`apple-selector-byte0-old-shared-path-components-2026-09-29/`. The patched
2013 component and intervention captures are under
`apple-selector-byte0-shared-path-components-2026-09-29/`,
`apple-selector-byte0-legacy-transition-weight-intervention-2026-09-29/`,
`apple-selector-byte0-context3-legacy-transition-intervention-2026-09-29/`,
and `apple-selector-byte0-context2-3-legacy-transition-intervention-2026-09-29/`.
The paired traces and runtime-only interventions are reproducible with
`trace-apple-old-path-costs-and-backpointers.gdb`,
`trace-apple-old-shared-path-components.gdb`,
`trace-apple-shared-path-components.gdb`,
`trace-apple-patch-2013-to-legacy-transition-weights.gdb`,
`trace-apple-patch-2013-context3-legacy-transition.gdb`, and
`trace-apple-patch-2013-context2-3-legacy-transition.gdb` respectively.
The native and adapted runners exited normally and recorded restoration of
their Stage 5 fixtures. The direct stack-float intervention is a causal
selector experiment, not a production DLL patch. A production-compatible
mapping from legacy voice/category state to 2013 scoring terms and full audio
parity remain to be validated.

### Generalized weight and consecutive-unit controls

The previous edge-specific override could have succeeded because it changed
only the eight measured competitors. A broader intervention replaced the
2013 pair-feature subtotal on every edge at contexts 2 and 3 with the measured
2006 raw/A/B coefficient schedules and bias, then set context-4 consecutive
global-ID edges to zero. It did not special-case candidate IDs. The selected
chain became `177774 → 177774 → 188826 → 129559 → 129560`; the WAVE was 13,040
bytes with SHA-256
`08561d76ee2924dba85f697075e204829ab4403e6b4a766bf36d74e96127f55f`.

Factor controls show which part moved the path. Applying the context-2/3
weight schedule alone produced that same chain and identical WAVE hash.
Applying only the context-4 consecutive-ID zero rule left the patched 2013
chain at `177774 → 177774 → 64719 → 2554 → 2555` and the WAVE unchanged at
14,248 bytes, SHA-256
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`.
Thus the tested context-4 shortcut alone does not change this path, while
context-wide weight replacement changes which context-2 candidate leads the
path. The result also shows that the eight-edge override was not equivalent
to replacing the old coefficient schedule for every candidate.

A native 2006 transition trace tested the newly selected IDs. It records no
context-2 scoring call for `188826`, but a boundary trace now locates why:
the old engine does generate that unit before its local-score cutoff. Its
context-2 class input is three classes (`5341`, `5342`, `10816`), which expand
to 48 units. `188826` appears at input index 14 with local score `0.222`; the
old ranker sorts all 48 candidates with no preserved prefix and keeps the
lowest 30. Its retained top-30 ends at score `0.181818`, and `188826` is
excluded before transition scoring.

The 2013 context-2 builder expands 50 units and computes a five-row
full-span/coverage prefix. `188826` is fourth in that prefix. The builder
sorts only the remaining 45 rows and retains the first 30 overall, so
`188826` stays at shortlist rank 4 and reaches the transition pass. The
complete runtime ID sets contain 43 shared candidates, five 2006-only IDs
(`19870, 45873, 225775, 251672, 276016`), and seven 2013-only IDs
(`42412, 92769, 127405, 145945, 258032, 265201, 271727`). The existing
source-row crosswalk independently reports the same 43/5/7 split for this
matched class. The two extra 2013 rows are the net of a real class-membership
mismatch, not a difference in trace boundary or duplicate counting.

The immediate `188826` divergence is still pre-transition retention: 2006
ranks all 48 rows by its local score and excludes this shared candidate at
the top-30 cutoff; 2013 protects its first five coverage-ranked rows before
sorting the tail, so the same ID reaches transition scoring. This explains
how the legacy-weight intervention can select `188826` even though the native
2006 transition scorer never sees it; it is not evidence of a missing old
query hit. The unmatched IDs establish a separate input-population defect for
this context's full pool, but all of them are removed before the respective
transition passes in the captured run.

A direct in-memory removal test now isolates the seven 2013-only candidates.
The list was reduced from 50 to the 43 shared IDs before scoring; coverage,
the remaining local ranking, and transition code were unchanged. The patched
engine still selected `177774 → 177774 → 64719 → 2554 → 2555` and reproduced
the baseline 14,248-byte WAVE hash
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`. The
2006-only five IDs are outside its retained top 30, and the seven 2013-only
IDs are outside the adapted engine's retained top 30. Both engines therefore
hand off only shared context-2 IDs to the transition pass in this run.

The two retained top-30 lists themselves share 22 IDs; each has eight IDs
absent from the other list, all drawn from the common 43-row candidate set.
This rejects the unmatched class members as a necessary cause of the observed
context-2 path difference. The active cause at this boundary is the
generation-specific shortlist ordering/retention over shared rows, followed
by the already measured differences in transition costs. This conclusion is
specific to this Apple context and fixture; other contexts still need the
same active-shortlist test.

The old scorer does also evaluate `129559` at context 3 and `129560` at
context 4. The minimum observed total for `129559` is `31.2409` through
`213256`, and the `129560 ← 129559` total is `31.7638`. Both exceed the old
selected-path costs at those positions (`27.5946` for `59558` and `27.9018`
for `59559`). The later rows therefore lose under the measured old path even
when present; context-2 retention is the earlier decisive boundary in this
counterfactual.

The generalized, weights-only, and consecutive-zero-only captures are under
`apple-selector-byte0-legacy-context-objective-2026-09-29/`,
`apple-selector-byte0-legacy-context-weights-only-2026-09-29/`, and
`apple-selector-byte0-legacy-consecutive-zero-only-2026-09-29/`.
The native comparison is under
`apple-selector-byte0-old-weight-counterfactual-candidates-2026-09-29/` and
`apple-selector-byte0-old-context2-candidate-boundary-2026-09-29/`; the paired
2013 shortlist and complete-candidate captures are under
`apple-selector-byte0-2013-context2-candidate-boundary-2026-09-29/` and
`apple-selector-byte0-2013-context2-candidate-members-2026-09-29/`. The full
candidate-list trace is `trace-apple-2013-context2-candidate-members.gdb`. The
seven-row removal and selected-ID capture are under
`apple-selector-byte0-2013-context2-prune-new-only-selected-2026-09-29/`,
reproduced with `trace-apple-2013-context2-prune-new-only.gdb`.
Reproduction uses the same adapted Apple Compose stack above with, in order,
`trace-apple-2013-legacy-context-objective.gdb`,
`trace-apple-2013-legacy-context-weights-only.gdb`, and
`trace-apple-2013-legacy-consecutive-zero-only.gdb`; the native capture uses
`trace-apple-old-weight-counterfactual-candidates.gdb` with
`compose-original-msi.yaml`, `trace-apple-old-context2-candidate-boundary.gdb`
with the same native stack, and `trace-apple-2013-context2-candidate-boundary.gdb`
with the adapted Apple stack. Every runner exited normally and restored its
Stage 5 fixture hashes. This resolves the context-2 membership delta and shows
that it is inactive after the measured caps. The corresponding context-3 and
context-4 comparisons now show that candidate availability is not sufficient
to explain the Apple path divergence:

| Context | Candidate evidence | Controlled result |
| ---: | --- | --- |
| 2 | The 2006 48-row pool and 2013 50-row pool share 43 rows. The top-30 lists share 22 IDs. | Restricting 2013 to the exact 2006 top-30 IDs leaves its chain and WAVE unchanged. |
| 3 | Both class crosswalks contain the same 20 source rows; neither engine truncates this set. | 2006 selects `59558`; 2013 selects `2554`, although both candidates reach transition scoring. |
| 4 | The 2006 pool has 157 rows; the adapted runtime has 158 entries with 158 distinct IDs. Their active top-30 lists share 20 IDs; `59559`, `82037`, and `2555` are present on both sides. | Restricting 2013 to the exact 2006 top-30 IDs leaves its selected `2555` chain and WAVE unchanged. |

At contexts 3 and 4, the within-engine cumulative costs reverse the ordering
between the compared alternatives. For context 3, 2006 favors `59558` at
`27.5946` over `2554` at `29.9268`, while 2013 favors `2554` at `31.0598` over
`59558` at `34.0749`. At context 4, 2006 favors `59559` at `27.9018` over
`2555` at `30.3042`, while 2013 favors `2555` at `36.0851` over `59559` at
`39.1103`. These are within-engine scores, not cross-generation comparable
absolute values. Together with the context-2 shortlist control, the evidence
locates the remaining Apple chain difference in ranker/transition objectives
over available rows, rather than the mere absence of the compared old-selected
units. A general scoring adaptation and full audio parity remain unverified.

The isolated context-2 and context-4 top-30 membership controls are under
`apple-selector-byte0-2013-context2-legacy-top30-2026-09-29/` and
`apple-selector-byte0-2013-context4-legacy-top30-2026-09-29/`; reproduce them
with `trace-apple-2013-context2-legacy-top30.gdb` and
`trace-apple-2013-context4-legacy-top30.gdb` on the adapted Apple Compose
stack. Both runs exited normally, retained the baseline patched chain and WAVE
hash, and restored the Stage 5 fixtures byte-for-byte.

### Follow-up: context-factor controls for the Apple path

Additional isolated overlays test which context groups are sufficient to
recover the native 2006 selected IDs. Each treatment changes only the measured
2013 pair-feature subtotal for the aligned edges listed below; candidate
construction, all other edge scores, local unit costs, and synthesis remain
2013-controlled. All four runs exited normally and restored both Stage 5
fixtures byte-for-byte.

| 2013 edge groups replaced with measured 2006 subtotals | Selected late chain | WAVE result |
| --- | --- | --- |
| Context 2: `100180←177774`, `64719←177774` | `64719 → 2554 → 2555` | Baseline 14,248-byte SHA-256 `520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865` |
| Context 4: `59559←59558`, `82037←82036`, `2555←2554` | `64719 → 2554 → 2555` | Same baseline WAVE SHA-256 |
| Contexts 2 and 4: the five edges above | `64719 → 2554 → 2555` | Same baseline WAVE SHA-256 |
| Contexts 3 and 4: `59558←100180`, `2554←64719`, `82036←64719`, plus the context-4 edges above | `100180 → 59558 → 59559` | 14,340-byte SHA-256 `4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8` |

Together with the earlier controls, context 3 alone and contexts 2+3 do not
recover the old chain; contexts 2+3 selected `64719 → 86327 → 86328`. The
contexts-3+4 treatment recovers it without changing the context-2 edge terms.
In this specific aligned-edge experiment, restoring the context-3 ranking
differences and context-4 adjacent-unit zero branch is sufficient; the
measured context-2 subtotal differences are not necessary for the final
chain. Context 4 alone does not recover it, so the tested context-3 change is
also necessary in combination with context 4. This is a local sufficiency
result over the sampled edges, not proof that those substitutions reproduce
the complete 2006 scoring function or generalize to other text.

A broader context-3 control replaces the 2013 feature formula on every edge
scored at context 3 with the measured old schedule
`raw×10 + derived-A×2 + derived-B×10 + 2`, while leaving context 2's formula
unchanged. It also zeros the same three verified context-4 continuation edges.
This recovered `100180 → 59558 → 59559` and produced the same 14,340-byte
WAVE hash as the sampled context-3+4 edge substitution. The three reported
context-3 edge values exactly match the native measurements. This supports a
context-wide coefficient correction for context 3 on this fixture; context 4
is still modeled with the three known row pairs, so this run does not test a
general continuation rule in 2013.

The 2006 pseudocode for `FUN_1001e470` sets the pair-feature term to zero when
either the current and previous candidate indices match while the context
selector byte equals 2, or the current index is exactly one greater than the
previous index and the previous candidate's continuation byte equals 1. All
three measured context-4 Apple edges have consecutive unit IDs and return
zero in 2006, while their measured distance columns and coefficients are
nonzero; the 2013 path adds 2. A native runtime breakpoint at `0x1001e952`
captured selector byte `0` and previous-candidate continuation byte `1` for
each exact pair: `59559←59558`, `82037←82036`, and `2555←2554`. The second
pseudocode disjunct therefore accounts for the observed legacy zero on these
rows. The capture repeats each pair in all three render passes and is under
`apple-selector-byte0-old-adjacent-zero-branch-2026-09-29/`; reproduce it with
`trace-apple-old-adjacent-zero-branch.gdb` on the native 2006 Apple Compose
stack. Ghidra output remains pseudocode, but the branch condition, raw marker,
and zero return are now linked by static and runtime evidence.

The context-2-only, context-4-only, contexts-2+4, and contexts-3+4 logs and
outputs are under `apple-selector-byte0-context2-legacy-subtotals-2026-09-29/`,
`apple-selector-byte0-context4-legacy-subtotals-2026-09-29/`,
`apple-selector-byte0-context2-4-legacy-subtotals-2026-09-29/`, and
`apple-selector-byte0-context3-4-legacy-subtotals-2026-09-29/`. Reproduce them
with `trace-apple-patch-2013-context2-legacy-subtotals.gdb`,
`trace-apple-patch-2013-context4-legacy-subtotals.gdb`,
`trace-apple-patch-2013-context2-4-legacy-subtotals.gdb`, and
`trace-apple-patch-2013-context3-4-legacy-subtotals.gdb`, respectively, using
the adapted Apple Compose stack described above.
The broader context-3 schedule plus context-4 continuation-edge control is
captured under
`apple-selector-byte0-context3-weights-context4-continuity-2026-09-29/` and
reproduced with
`trace-apple-patch-2013-context3-weights-and-context4-continuity.gdb`.

### Targeted 2013 continuation-marker control

The 2013 disassembly contains the corresponding adjacent-row shortcut at
`0x10018ed3–0x10018ed9`: after confirming that the current candidate index is
one greater than the previous index, it tests bit 7 of the previous
candidate's signature byte at offset `+6`; if set, the pair-feature subtotal
is zeroed. The prior adapted Apple index had byte value `0` for rows 59558,
82036, and 2554, even though the matched 2006 rows enter the equivalent
branch with continuation byte `1`.

A disposable overlay changed only those three bytes in `unit-gen.idx`; its
manifest verifies that every other byte in that index is unchanged. The
matching 2005 source rows each have `attr_48=1`.
2013 runtime then read `128` for each predecessor and took its native branch
for `59559←59558`, `82037←82036`, and `2555←2554`. With the independently
measured old context-3 schedule
`raw×10 + derived-A×2 + derived-B×10 + 2`, it selected the complete native
chain `177774 → 177774 → 100180 → 59558 → 59559`. The resulting 14,340-byte
WAVE hash, `4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8`,
matches the previous runtime edge-subtotal control exactly. Both Stage 5
fixture hashes were restored.

The factorial counterpart applied the same three on-disk marker bytes without
changing any 2013 scoring weights. The 2013 branch took all three shortcuts,
but the selected chain remained `177774 → 177774 → 64719 → 2554 → 2555` and
the WAVE remained at the patched baseline hash
`520597d82a969c4d03723e884a5d517ccb5795db21781a5374702b2a6eb03865`. A
matched run applied the broad context-3 schedule to the unmarked Apple index;
it selected a different chain beginning `273369 → 273370 → 273370 → 273371`
and produced a 24,720-byte WAVE with SHA-256
`ff36d8aa0b5487a9940f4bb422181d463c8fda0a737f88d3500d4873f6960d04`. Both
factorial runs restored their pre-run fixture hashes. The four treatments
therefore show an interaction: baseline and marker-only retain the patched
path, schedule-only moves to another path, and only marker plus schedule
recovers the old selected chain and the same WAVE hash as the previous
selector intervention. These results apply to this input and tested index
adapter.

This resolves the earlier broad-overlay ambiguity: setting the marker from
`attr_48` on all rows activated the branch but also changed other candidate
behavior, yielding a different chain. Restricting the marker to the three
measured predecessors recovers the expected chain and hash. The result shows
that the 2013 branch is functional and that these three index markers plus
the context-3 coefficient schedule are sufficient for this Apple path. It
does not yet establish that `attr_48` universally maps to signature byte 6
bit 7, or that these changes recover native audio bytes or generalize to
other contexts and voices.

Reproduce the targeted run with
`make-apple-selected-predecessors-bit7-overlay.py`,
`compose-apple-selected-predecessors-bit7.yaml`, and
`trace-apple-patch-2013-context3-weights-selected-bit7-index-overlay.gdb`.
The overlay manifest is
`index-adapter-key-repacked-attr48-key3-attr40-key2-apple-selected-predecessors-bit7/selected-predecessors-bit7-manifest.txt`;
the runtime capture is under
`apple-selector-byte0-context3-weights-selected-bit7-2026-09-29/`.
The marker-only factorial capture is under
`apple-selector-byte0-selected-bit7-only-2026-09-29/` and uses
`trace-apple-patch-2013-selected-bit7-index-overlay-only.gdb`. The matched
schedule-only control is under
`apple-selector-byte0-context3-weights-only-2026-09-29/`; reproduce it with
`trace-apple-patch-2013-context3-weights-only.gdb` and the unmarked adapted
Apple index.

### Context-3 pool marker crosswalk

The complete patched-2013 context-3 pool contains 20 candidates:
`2554, 256581, 232462, 218670, 158830, 131118, 129559, 103238, 97884,
86327, 82036, 59558, 41057, 40214, 28631, 27966, 27794, 27440, 11135,
258822`. The corresponding 2005 source rows each have `attr_48=1`; the pool
spans 16 rows in `unit-gen` and four in `unit-gen2` under the recovered global
ID ranges.

A second disposable overlay set signature byte `+6` bit 7 for all 20 pool
rows according to those source markers. Its manifest verifies exactly 20
changed bytes and structural validity of all five index files. Combined with
the broad context-3 schedule, the runtime again selected
`177774 → 177774 → 100180 → 59558 → 59559` and produced the identical
14,340-byte WAVE hash
`4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8`. This
matches the three-predecessor overlay, so adding markers to the other 17
active context-3 candidates does not change the measured result. This is
stronger evidence for the `attr_48` to signature-bit mapping on this pool,
but it remains a controlled one-input, one-context conversion rather than a
corpus-wide validation.

The pool capture is under `apple-selector-byte0-context3-pool-2026-09-29/`;
reproduce it with `trace-apple-2013-context3-candidate-pool.gdb`. Build the
pool-scoped marker overlay with `make-apple-context3-pool-bit7-overlay.py`
and mount it with `compose-apple-context3-pool-bit7.yaml`. The combined
runtime capture is under
`apple-selector-byte0-context3-pool-bit7-weights-2026-09-29/` and uses
`trace-apple-patch-2013-context3-weights-context3-pool-bit7.gdb`.

### Static coefficient table and native DLL patch control

A live breakpoint at `0x10018fbb` confirms that the 2013 scorer selects its
coefficient triple with a signed selector byte at stack offset `-1`, using a
12-byte stride. The static table reads raw, derived-A, and derived-B weights
from VAs `0x1007c288`, `0x1007c290`, and `0x1007c28c`, then adds the constant
at `0x1006d174`. The measured Apple edges use these slots:

| Apple context | Table slot | Patched 2013 weights | Measured 2006 weights |
| ---: | ---: | --- | --- |
| 2 | 2 | `(10,2,5)` | `(10,1,5)` |
| 3 | 1 | `(10,2,5)` | `(10,2,10)` |
| 4 | 0 | `(10,5,10)` | pair shortcut returns zero for the three measured adjacent rows |

The context-3 discrepancy is one table value: slot 1 derived-B is `5.0` at
VA `0x1007c298`; the legacy schedule requires `10.0`. A hash-guarded
disposable DLL changed only that float from little-endian `00 00 a0 40` to
`00 00 20 41`. No scoring stack values were edited in the runtime. With the
20 context-3-pool markers derived from legacy `attr_48`, the patched DLL
selected `177774 → 177774 → 100180 → 59558 → 59559` and produced the same
14,340-byte WAVE hash as the in-memory controls:
`4d03392e63af2949f31a2d26fe2af6f59759fc54bf8d7870a06194c1a7940bc8`. The
runtime exited normally and restored the Stage 5 fixture hashes.

The context-2 table difference is independently present at slot 2: derived-A
is `2.0` where the measured old schedule uses `1.0`. A second disposable DLL
patched both that float and the context-3 derived-B float. With the same
20-row marker overlay, it selected the same five IDs and produced the
identical 14,340-byte WAVE hash. Thus the context-2 correction adds no audible
or selected-unit change on this fixture after the context-3 and marker
corrections; the observed context-2 table difference is not required for this
recovered output.

### Plain-Hello coefficient transfer control

The slot-1 derived-B patch was then checked on `Hello.` using the same
vector-pitch trees and matched-context `attr_b`-transfer index. In the
unmarked-index A/B pair, the stock and patched DLLs each selected
`273369, 273370, 273370, 273371, 264074, 264074` and produced the same 15,876-
byte WAVE, SHA-256
`8d6b4f63906ad339f5b9b4dde977a495d2c7dfc7a7c9a0818e2a6237d81b3dfb`.
The same breakpoint at `0x10018fbb` recorded 1,776 coefficient lookups in
each run. Three hundred calls used table slot 1; their live triple changed
from `(10,2,5)` to `(10,2,10)`. No other traced slot/mode triple changed.
Thus the table patch is active on plain Hello and changes a consumed scoring
coefficient, but does not change the selected chain or rendered bytes for
this input. A second stock/patched pair with identical path-cost tracing
returned 111 cumulative candidate rows per run. Twenty-nine rows changed
cost and/or best-predecessor link: 8 at context 1, 11 at context 2, 9 at
context 3, and 1 at context 4; context 5 had no changes. The changed context-2
through context-4 rows carry forward alternatives affected by earlier
predecessor costs. The selected rows kept the same cumulative scores in both
runs: `273370` at 2.98051, `273370` at 3.92356, `273371` at 6.45379, `264074`
at 14.7554, and `264074` at 15.7349. The largest changed cumulative cost was
for candidate `13942` at context 1, rising from 16.6608 to 23.2013. This
locates the coefficient effect in losing alternatives for this Hello path,
rather than in a silent or inactive patch.

With the three-row continuity-marker overlay, the same patched-DLL capture
selected `272822, 272823, 272823, 272824, 272825, 272825` and produced the
same 15,264-byte WAVE as the existing stock-DLL marker control, SHA-256
`936e0f8f3c62399858ee21777a7f6ef89fa37c98134fcc4e1263dc78f59168cc`. The
weight-table breakpoint was installed and the process reached the ready
breakpoint, but no lookup at `0x10018fbb` occurred in this marker run. This
suggests that the marker-driven path bypasses that scorer entry for this
query; it does not prove that every scoring routine is bypassed.

Together, these controls show the Apple-derived coefficient correction does
not generally alter plain-Hello output: on the unmarked path it changes 300
live coefficient reads without changing the result, and on the marked path
the traced scorer entry is not visited. This does not invalidate the Apple
result, where the corrected slot-1 term participates in the winning-versus-
competing path comparison. Nor does it show that the complete 2006 and 2013
scoring schedules match for Hello. Reproduce with
`trace-hello-active-weight-table.gdb` and
`compose-apple-context3-bweight10-dll.yaml`; the four stock/patched evidence
sets are under `hello-context3-weight-control-unmarked/`,
`hello-context3-weight-patch-unmarked/`, `hello-context3-weight-patch/`, and
`hello-context3-weight-patch-trace/` within the ignored Stage 20 capture tree.
The paired path-cost captures use `trace-hello-weighted-local-costs-control.gdb`
and are under `hello-context3-weight-path-control/` and
`hello-context3-weight-path-patch/`.

### Cross-generation interpretation for plain Hello

The paired path trace resolves how the coefficient result relates to the
native 2006 route. The old engine's `FUN_1001cd60`/`FUN_1001cff0` path reduces
the four position pools (75, 9, 5, 6 candidates) to their only full-span
rows, `272822 → 272823 → 272824 → 272825`, before the local scorer runs. In
the unmarked 2013 trace, that legacy subsequence remains available and is
cheaper through the first two returned contexts: `272823` costs 2.32276 then
2.85916, compared with `273370` at 2.98051 then 3.92356. The cumulative
ordering crosses at the next context: `272824` costs 12.5038 versus `273371`
at 6.45379. Later, `272825` costs 17.4731 then 19.5923, compared with the
selected `264074` route at 14.7554 then 15.7349. The slot-1 DLL patch leaves
these legacy and selected-route costs unchanged; the two 2013 A/B runs choose
the same route and produce the same WAVE.

The 2006 span-filter-bypass diagnostic provides one direct edge comparison:
`272823` totals 1.64668 through `272822`, versus 32.5011 through `273369`.
That old-engine counterfactual strongly favors the legacy predecessor, but
the native 2006 run normally removes partial-span alternatives before this
scoring stage. This is not an active native candidate comparison.

The corrected three-row marker overlay changes 2013's span traversal and
selects `272822 → 272823 → 272823 → 272824 → 272825 → 272825` under the
unchanged native cutoff. This links the old/new Hello difference to missing
continuity metadata and candidate-pool width before the coefficient term is
considered. The old and new cost values must not be compared numerically as a
shared scale: the old singleton filter and 2013 multi-candidate dynamic
program operate on different candidate sets and scoring paths. The within-
engine A/B and marker interventions establish the causal effects here; they
do not establish a corpus-wide mapping or general parity.

The old full-span evidence and counterfactual costs are under
`hello-plain-crosswalk-old-span-filter-bypass-v4/`; the low-intrusion marker
comparison is under `hello-continuity-bit7/{hardware-base,hardware-marker}/`.
See the corresponding Stage 20 sections for capture controls and restoration
hashes.

This is the first code/data-level correction exercised in the actual
versioned 2013 DLL path, rather than a GDB score substitution. It confirms a
concrete fix for the measured context-3 term and the context-3 marker pool on
this Apple fixture. The context-2 slot difference was not patched in this
control, the context-4 branch still depends on correctly supplied marker
bits, and the result does not establish corpus-wide or audio-byte parity.
Reproduce with `patch-apple-context3-weight-dll.py`,
`make-apple-context3-pool-bit7-overlay.py`,
`compose-apple-context3-pool-bit7.yaml`,
`compose-apple-context3-bweight10-dll.yaml`, and
`trace-apple-context3-pool-markers-native-weight-patch.gdb`. The DLL patch
manifest is `vt_kat-context3-bweight10.manifest.txt`; runtime evidence is
under `apple-selector-byte0-context3-native-weight-patch-2026-09-29/`.
The context-2-plus-context-3 variant is built with the same helper's
`--include-context2` option, mounted by
`compose-apple-context2-context3-weight-dll.yaml`, and captured with
`trace-apple-context3-pool-markers-context2-context3-dll.gdb` under
`apple-selector-byte0-context2-context3-native-weight-patch-2026-09-29/`.

### Held-out 2005 voice key-partition audit

The proposed field layout and the actual 2013 class-key lookup tables were
applied offline to every row in the Kate and Bridget 2005 indexes. For both
voices, the old feature projection uses tables captured from the local 2006
Kate DLL; Bridget has no paired 2006 runtime capture in this checkout. The
comparison then derives the candidate 2013 five-byte key from the repacked
signature and the three 2013 lookup tables. The Bridget result is therefore a
held-out structural projection, not an observed Bridget selector result.
The three lookup-table slices used for this projection are byte-identical in
`vt_pau.dll` and the runtime-tested `vt_kat.dll`.

| Voice | Index rows | Distinct old features | Distinct 2013 keys | 2013 keys containing multiple old features | Rows outside the largest old feature in their 2013 key | Old features split across 2013 keys |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Kate | 283,696 | 52,627 | 33,206 | 10,380 (31.3%) | 63,037 (22.2%) | 0 (0.0%) |
| Bridget | 780,882 | 76,893 | 35,198 | 14,119 (40.1%) | 282,904 (36.2%) | 3,764 (4.9%) |

Under this field conversion, the partition loss is not limited to the selected
Kate rows. Bridget also shows some old feature classes split across multiple
projected keys, while Kate does not in this full-index pass. This establishes
that the tested five-byte projection cannot preserve the complete partition
induced by the reused 2006 transform for either model. It does not establish
that Bridget's actual 2006 runtime uses identical lookup behavior, or that
every merged row changes a live query result: the 2013 selector may
intentionally use broader classes, and no held-out runtime queries were
replayed.

A second comparison used the structural order `[attr_48,key[0:5],attr_40]`
against the runtime-tested signature order `[attr_48,key[0:3],attr_40,key[3:5]]`.
For each layout, a row-pair merge is a pair with different derived 2006
features but the same 2013 exact key; a row-pair split is the inverse. Counts
are full-index structural metrics, not counts of live query errors.

| Voice | Layout | Merged row pairs | Split row pairs | Total partition disagreements |
| --- | --- | ---: | ---: | ---: |
| Kate | Structural | 2,426,245 | 910,383 | 3,336,628 |
| Kate | Signature repack | 4,465,016 | 0 | 4,465,016 |
| Bridget | Structural | 33,123,671 | 10,961,551 | 44,085,222 |
| Bridget | Signature repack | 87,388,861 | 5,193,169 | 92,582,030 |

The structural order has fewer total pair disagreements in both packages,
while the signature repack avoids more old-class splits and merges more
different old features together. This is a useful ranking signal for the two
candidate conversions, not proof that the structural order gives better
speech: the runtime A/B shows different unit streams, and the Bridget old
engine has not been run.

The reproducible analyzer is
`tools/revkit/work/stage20/audit-heldout-voice-class-projection.py`. It reads
the read-only Kate and Bridget `ver.2005` indexes and local DLL lookup tables;
it writes no model files. The key-partition result narrows the remaining
question: which lost distinctions are required by the live old query classes,
and which are deliberately absorbed by the 2013 class and ranking stages.

### Bridget 2005 runtime transfer control

The same signature-repack hypothesis was applied to Bridget's `ver.2005`
indexes and run through the unmodified 2013 `vt_kat.dll` with Bridget's
`tree3`, DAT, UPM, and distance resources mounted read-only under the Kate
runtime paths. Because Bridget's index header says `VoiceText-Bre` and the
2013 reader accepts `VoiceText-Eng`, the disposable index header was
normalized to `VoiceText-Eng`. This is an adapter requirement observed in the
current reader, not a recovered Bridget package conversion.

The `Hello.` run completed and wrote a valid mono 16 kHz 16-bit WAV. The
signature-repacked variant is 7,547 frames (15,138 bytes; SHA-256
`bcb536ddeeff92730603d26b2d924bdd27e605118503a71d2246e5b5cfa0092e`). A
same-input structural control that retained the old seven-byte field order
completed at 7,991 frames (16,026 bytes; SHA-256
`9023035e14daa515c2a35499f6e5d4d73301884bf484b5351f69cb36e441ce10`). The
hardware selected-unit trace recorded 10 callbacks for the repack and 12 for
the structural control, with different IDs. Thus signature placement changes
the active candidate/selection path on a second 2005 voice under the 2013
runtime; the difference is not only a parser acceptance or WAV-header effect.

The paired GDB traces locate the first divergence after exact query-key lookup.
Both layouts produce the same ten query signatures and projected keys for
`Hello.`; all ten exact lookups return zero in both runs. The post-builder
candidate sets then diverge:

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

This isolates the tested mapping's effect to index-driven fallback candidate
construction after exact-key misses, before final backtracking. The tree query
producer is not the variable in this A/B. Positions 4 and 5 show the largest
membership inflation and least overlap, matching the materially different
selected callbacks. These remain Bridget/2013 results; they do not identify
which set the native 2006 engine would have selected.

### Bridget candidate membership to continuity coverage

The candidate-node metadata shows how the position-4/5 membership change
reaches the 2013 scoring path. The `FUN_100230a0` trace records the 2013
continuity fields at candidate offsets `+0x0c` (right-match count), `+0x0e`
(total span), and `+0x10` (weighted side coverage). The full candidate-list
comparison is reproducible with
`tools/revkit/work/stage20/summarize-bridget-span-parity.py`:

| Position | Signature total / span 2 | Structural total / span 2 | Shared IDs | Shared IDs with different span |
| ---: | ---: | ---: | ---: | ---: |
| 2 | 47 / 46 | 27 / 6 | 3 | 1 |
| 3 | 49 / 30 | 23 / 6 | 14 | 11 |
| 4 | 1,848 / 1,343 | 162 / 5 | 132 | 127 |
| 5 | 2,015 / 30 | 49 / 5 | 14 | 4 |

All 127 shared position-4 IDs whose span changes from 1 in the structural
layout to 2 in the signature layout are present in the signature layout's
position-5 candidate pool and absent from the structural layout's position-5
pool. The other five shared position-4 IDs occur in both position-5 pools and
retain the same span. This ties the dominant position-4 metadata change to
adjacent-position candidate membership before transition scoring.

Unit `16223` is a direct row-level example. Its converted signatures are
`01 30 2b 32 43 00 0b` (signature layout) and
`01 30 2b 32 00 0b 43` (structural layout). It appears in positions 4 and 5
only under the signature layout. At `FUN_100230a0` address `0x100232ea`, the
runtime records `(right-match, total-span, weighted-coverage)` as `(1,2,0)`
for the signature layout and `(0,1,0)` for the structural layout. Both runs
retain their original Bridget A/B WAV hashes.

Four shared position-5 IDs change in the opposite direction: span 2 in the
structural layout and span 1 in the signature layout. A direct trace at
`FUN_100230a0` shows that all four are in position 4's retained shortlist only
in the structural run. Row `150652` has the same local score (`163.89`) in both
runs, but its post-coverage rank is 606 of 1,848 in the signature run and 1 of
162 in the structural run. The 2013 sort retains only the first 30 here, so
the signature run drops it while the structural run keeps it. The resulting
position-5 metadata is span 1 versus span 2; the other three rows follow the
same previous-shortlist membership rule.

The class-key populations identify which converted bytes create that rank
pressure. At `FUN_10023350`, each returned class ID expands through the current
model's class member count/list (`model+0x8c` and `model+0x94`). The three
classes containing the four position-5 rows have the same five-byte class keys
in the two layouts, but radically different member counts:

| Projected class key | Target rows | Signature count | Structural count |
| --- | --- | ---: | ---: |
| `0d 2b 2f 00 00` | `365789`, `382197` | 412 | 2 |
| `10 2b 2f 00 00` | `548226` | 444 | 1 |
| `23 2b 2f 00 00` | `150652` | 38 | 1 |

The full position-4 pool is 1,848 candidates in the signature layout and 162
in the structural layout. This is not a score-table difference: the current
class projection reads signature bytes `+1,+2,+3,+5,+6`, and the adapter
layouts place different legacy fields at `+5` and `+6`.

Two disposable field-isolation layouts tested those offsets while preserving
the legacy `attr_40` value at `+4`. With `key[3]` at `+5` and `attr_40` copied
to `+6`, the position-4 pool falls to 827; class populations for the three
target keys fall to 115, 325, and 19. Row `150652` then ranks 81 of 827 and
still falls outside the top 30, so all four position-5 spans remain 1. With
`key[4]` at `+5` and `attr_40` at `+6`, the position-4 pool is 162, the same
as the structural control; row `150652` ranks 1 and is retained. That layout
also reproduces the structural WAV hash. These interventions isolate both
trailing selector inputs as contributors to class membership and ranking;
they do not establish which experimental placement matches a native Bridget
2006 selector.

The exact target-key lookups remain identical and empty in the paired runs.
Thus the demonstrated divergence is after those misses, when relaxed
class-key projection selects class populations and the per-position top-30
shortlist carries that result into continuity coverage. This is direct
Bridget/2013 adapter evidence, not a comparison with a native Bridget 2006
selection.

This remains a conversion experiment, not proof that either output is
intelligible or matches a native Bridget 2006 output. No paired Bridget 2006
runtime or listening comparison is present here. The build and capture helpers
are `build-bridget-2005-index-adapter.py`,
`compose-bridget-2005-adapter.yaml`,
`compose-bridget-2005-structural-adapter.yaml`,
`compose-bridget-2005-selector-key3-attr40.yaml`,
`compose-bridget-2005-selector-key4-attr40.yaml`,
`run-bridget-2005-adapter.sh`, `trace-bridget-context-classes.gdb`,
`trace-bridget-rank150652.gdb`, and `trace-bridget-selected-units.gdb`.
`compare-bridget-selected-index-rows.py` maps selected IDs to their source
bank-local fields. Both runs restored the Stage 5 input and output fixture
hashes.

### 2026-09-30 Kate key5 signature penalty and path replay

The matched Kate `Hello.` run with `attr48-key5-attr40` exposes a second
failure after candidate expansion: old-engine rows can be present in the 2013
candidate list yet receive a large local categorical penalty. The baseline
2013 WAVE is 14,142 bytes (SHA-256
`478c155126f39e07f013cd5202a41275ad77035a7adee0d7f8d1265fb752b317`) and
selects `280485, 21433, 264072, 264072, 264073, 264074, 264074, 21433,
280485`. The native 2006 chain for the same plain `Hello.` input is
`272822, 272823, 272824, 272825`.

Direct local-score and predecessor traces show the following relevant pairs:

| 2013 context | Native 2006 row / local cost | 2013 row / local cost | 2013 cumulative path costs |
| ---: | --- | --- | --- |
| 0 | `272822` / 505.8 | `280485` / 34.8 | — |
| 1 | `272822` / 505.8 | `21433` / 39.8 | 318.8 / 151.131 |
| 2 | `272823` / 1003.89 | `264072` / 22.3004 | 657.395 / 181.853 |
| 3 | `272824` absent from candidates | `264072` / 22.1471 | — |
| 4 | `272824` / 79.3595 | `264073` / 71.5588 | 260.67 / 230.706 |
| 5 | `272825` / 34.0766 | `264074` / 32.5632 | 279.709 / 248.988 |
| 6 | `272825` / 32.3766 | `264074` / 32.2771 | 281.801 / 265.126 |

At context 2, rows `272823` and `264072` share signature bytes `+0..+4`
(`01 22 17 2b 00`), while byte `+5` is `1e` on the old row and `00` on the
2013 row. A runtime-only intervention swapped only signature byte `+5` for
these two rows while `FUN_100182e0` scored them. Their local costs changed
from `1003.89` and `22.3004` to `23.5481` and `1003.58`, respectively. This
directly identifies the old row's roughly 1,000-point local penalty as a
comparison involving the `+5` signature field under this layout. The two
earlier cost changes are not explained by the two scalar feature arrays at
model offsets `+0x54` and `+0x5c`: swapping only those arrays leaves the
context-2 old-row cost at `1003.8`.

The `+5` intervention changes the selected 2013 sequence to
`280485, 21433, 272152, 272152, 272153, 246677, 246677, 21433, 280485` and
the WAVE hash to
`c3946893f170ebdbcb2f5e6f18713b3be479ba2d3f3808fef5b8254a79766741`.
The corrected `272823` path reaches cumulative cost 167.225 at context 2,
below `272152` at 182.371. It still loses the completed path: at context 4,
`264073` reaches 250.344 through `272152`, versus `272824` at 257.216 through
`45237`; at context 5, `264074` reaches 256.63 versus `272825` at 276.254.
Context 3 has no `272824` candidate. Thus the field mismatch explains the
large local penalty on this old row, while the altered seven-context chain
and its later transition costs remain independently incompatible.

The supported diagnosis is specific to this experimental repack: placing
legacy `key[4]` at byte `+5` feeds that byte to a 2013 scorer comparison whose
categorical behavior is incompatible with the old row's value. The intervention
proves the byte's causal effect on local cost; it does not establish the
correct general mapping for 2005 voices. Class membership, the absent context-3
row, and the later path costs still require a mapping that matches the native
2006 behavior. Reproduction uses
`trace-kate-key5-score-path.gdb`, `trace-kate-key5-score-fields.gdb`,
`trace-kate-key5-feature35-swap.gdb`, and
`trace-kate-key5-signature5-swap.gdb`; captures are in the ignored
`corpus-parity/stage20/kate-key5-score-path-hello/` directory. All runs used
the same `hello-plain.txt` fixture and restored the Stage 5 input/output
hashes.
