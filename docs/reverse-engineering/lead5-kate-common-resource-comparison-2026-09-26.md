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
because it imports `VT_SetCommaPause_ENG`, absent from the MSI's 2006 DLL. I
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

To remove dictionary lookup from the experiment, We sent single visible words
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
unused half or from later timing/coarticulation; listening to these
side-selected composites is the next discriminating check.

The fresh downstream runs for P AH0 and T AH0 reproduced their prior WAVE
hashes. The IH0 rows above use the earlier valid captures
`trace-selection-valid-pih0.log` and `trace-selection-valid-tih0.log`, whose
WAVEs are 6,366 and 7,690 bytes. A separate rerun of both IH0 inputs with the
new generalized runner produced the same 22,084-byte output for both and
selected unrelated IDs; We excluded those two reruns from the table. Every
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

## What this changes for Lead 5

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
the copied legacy signature fields and metric columns remains unresolved.
This moves the leading mismatch investigation toward the old-to-new
signature, feature-column, and tree/class mappings; the observed 500-point
context penalty is one concrete sign of that boundary. The side-selected
composites can now test whether unused record halves caused the initial
listening mismatch.
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
