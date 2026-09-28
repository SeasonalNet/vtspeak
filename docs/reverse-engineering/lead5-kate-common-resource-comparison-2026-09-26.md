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

The original and adapted WAVs are 13,648 bytes with SHA-256
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
