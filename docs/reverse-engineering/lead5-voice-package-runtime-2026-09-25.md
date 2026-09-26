# Lead 5: voice and package compatibility (2026-09-25)

**Status: runtime paths checked; intelligible Kate output remains unresolved.**
This report combines existing structural package checks with controlled Wine
runs of each supplied host/DLL pair, plus a package-supplied James DLL check.
It does not claim full synthesis parity, support for every engine release, or
compatibility of an untested package.

## Scope and method

The local inputs contain Paul and James `ver.2013` English models, Julie's
`ver.2009` English model, Kate's `ver.2005` English model, and Bridget's
`ver.2005` `VoiceText-Bre` model. The `binary/` directory contains host/DLL
pairs for Paul, James, Julie, and Kate; it has no Bridget host in the
repository's `binary/` set. The supplied four DLLs share executable `.text` bytes according
to the recorded static comparison. Their differing `.data` sections and
voice/model packages remain distinct evidence; a shared code section does not
establish package compatibility.

The read-only inventory pass covers all 21 local unit indexes and 2,658,093
indexed units across the five packages. The existing index inspectors accepted
each package's declared index layout and all indexed DAT/UPM spans; configured
sample decodes matched UPM-derived sample counts. The tree parsers accepted all
85 voice trees: 51 `tree3` files and 34 `tree2` files. Eight `tree2` pitch
files have unconsumed suffixes, so acceptance covers the parsed recursive
prefix rather than complete file recovery. This structural pass does not
semantically decode index columns or tree outputs.

The Stage 17 run used the existing Stage 8 Wine container and the same three
inputs for Paul and James: ordinary prose, a number/time phrase, and an
address/abbreviation phrase. Julie and Kate each received the prose input,
which is sufficient to identify their first load failure. Model and common
data were mounted read-only. The hosts resolve package paths relative to the
Stage 5 working directory; the James and Julie executables request the literal
path components `data-jame` and `data-juli`, so their existing model directories
were mounted under those names without changing the source packages. Stage 5
input and output fixtures were backed up and restored by the runner.

## Compatibility matrix

| Voice/package | Existing structural evidence | Supplied host/DLL runtime result | Boundary |
| --- | --- | --- | --- |
| Paul, `ver.2013` English, 580,474 units, `tree3` | All local DAT records were checked against the Python decoder at the decoder-output boundary; selected synthesis paths were also cross-checked. | All three inputs returned 0 and produced valid 16 kHz mono 16-bit PCM WAVE files. | Strongest local coverage; still not whole-synthesis corpus parity. |
| James, `ver.2013` English, 284,349 units, `tree3` | All index records and DAT/UPM extents passed the inventory checks; sampled decoded counts matched UPM counts. | All three inputs returned 0 and produced valid 16 kHz mono 16-bit PCM WAVE files. | Demonstrates a short successful synthesis path for James, not full decoder or corpus parity. |
| Julie, `ver.2009` English, 728,692 units, `tree2` | All index records and ordered DAT/UPM spans were checked; first, middle, and last records per bank matched sampled UPM counts. | Prose returned 1 and produced no WAVE. Wine records a missing `../data-juli/M16/ttsdata/tree3/pitch/nbt.tree3` (`0xc000003a`). | The supplied DLL requests a `tree3` resource absent from Julie's `tree`/`tree2` package. This does not test Julie under a matching older engine. |
| Bridget, `ver.2005` `VoiceText-Bre`, 780,882 units, `tree3` | All index records and ordered DAT/UPM spans were checked; sampled payload counts matched UPM counts; all 17 trees parse as `tree3`. | Not run: no Bridget host or DLL is present in the local binary set. | Structural compatibility only; the `Bre` resource label does not independently establish locale or runtime support. |
| Kate, `ver.2005` English, 283,696 units, `tree2` | All index records and ordered DAT/UPM spans were checked; sampled payload counts matched UPM counts. | The supplied pair requests absent `tree3/pitch/nbt.tree3` and returns 1. A disposable patched standard DLL plus converted tree/index overlays completed three runs and wrote non-silent WAVs (Stage 19). The requester listened and reports gibberish, with only slight sentence-like fragments in prose. | Resource loading and audio generation work; intelligible voice compatibility is not demonstrated. The adapter converts recursive trees to the current indexed format and adapts legacy index columns; it does not restore the native recursive parser. |

The structural counts, index generations, and tree-family notes are detailed
in the [voice package generation comparison](voice-package-generation-comparison.md).
The `.text` match applies only to the four supplied DLLs; it does not include
Bridget's CAB-labeled `vt_bre_bridget16.dll` or the historical Kate MSI DLL
used for the static `tree2` reader analysis. Additional package-DLL evidence
and candidate limitations are recorded below.

## Additional package-engine evidence

The James package under Downloads contains a second PE32 DLL,
`vt_eng_james16.dll` (SHA-256
`7644aef98f08c5929283c15a719deb3a8395eb9b5b6fd80862f9aace16522580`). Its
`_ENG` exports satisfy the supplied James host's named imports. Mounted
read-only as `vt_jam.dll`, with the model also mounted under the path requested
by that DLL, it produced valid WAVs for all three Stage 17 fixtures. The
outputs have the same byte lengths and frame counts as the supplied 2013 DLL
outputs, but different complete-file hashes and PCM payloads:

| Input | Package DLL bytes / frames | `binary/vt_jam.dll` bytes / frames | PCM equal |
| --- | ---: | ---: | --- |
| Prose | 97,196 / 48,576 | 97,196 / 48,576 | No |
| Number/time | 110,548 / 55,252 | 110,548 / 55,252 | No |
| Address/abbreviation | 107,638 / 53,797 | 107,638 / 53,797 | No |

This confirms a second successful short-synthesis path for the James model.
Equal frame counts do not establish sample parity or explain the output
differences. The captured files, traces, hashes, and reproducible mount/runner
are in [Stage 17](../../tools/revkit/work/stage17/README.md).

Two other historical DLL candidates were found in Downloads. The Julie DLL is
stored in a directory named `crack`; with the supplied Julie host it opens the
older shared dictionary files, then exits when `dict-eng/dic.pos.txt2` is
missing, before opening Julie's voice trees. Its modified/license state is
unknown, so this exploratory attempt is not accepted as an older-engine
compatibility result. A Bridget DLL candidate is also stored in a `crack`
directory; it exports the BRE-named API, and no matching Bridget host or full
legacy common-data set is available for a controlled run. These candidates do
not change the supported matrix.

## Successful output observations

The six successful Paul and James files are valid RIFF/WAVE, PCM format 1,
mono, 16 kHz, 16-bit samples. Their shared inputs produced different sample
counts:

| Input | Paul bytes / frames | James bytes / frames |
| --- | ---: | ---: |
| Prose | 89,478 / 44,717 | 97,196 / 48,576 |
| Number/time | 89,716 / 44,836 | 110,548 / 55,252 |
| Address/abbreviation | 108,358 / 54,157 | 107,638 / 53,797 |

The differing outputs are observed voice/package behavior for these three
inputs. They are not a quality comparison and do not identify which model
fields account for the differences. Per-file hashes are recorded in
[`wav-sha256.txt`](../../tools/revkit/work/stage17/wav-sha256.txt); the runtime
matrix and captures are under
[`Stage 17`](../../tools/revkit/work/stage17/README.md).

## Conclusion and remaining boundary

The checked local `ver.2013` Paul and James packages both load and synthesize
the short shared input set with their checked-in host/DLL pairs. The separate
James package DLL also synthesizes all three fixtures. The supplied DLL set
does not load the older Julie and Kate packages past the first voice-tree
request because it asks for `tree3` while those packages expose `tree2`. A
disposable patch to the standard Kate DLL, paired with converted tree and
index resources, now completes all three fixture runs and writes non-silent
audio. The requester reports the audio is gibberish, so this does not
establish intelligible Kate synthesis. Bridget's model structures are
validated, but there is no Bridget host or complete matching runtime/data set.
Broader text/output parity remains unestablished.

All five checked-in model packages have structural checks, and all four
checked-in host/DLL pairs were attempted. The package-supplied James DLL adds
three more successful runs, and the Stage 19 patched-standard adapter adds
three completed Kate audio-generation runs, although listening review found
the output unintelligible. Julie's external candidate is blocked at a
missing legacy common dictionary file and has unknown modified/license state;
the Bridget candidate has the same provenance concern and lacks a matching
host and common-data set. Kate's adapter demonstrates a bounded resource-load
and audio-generation path; intelligibility and voice compatibility remain
unresolved.

A filename-only patch experiment is recorded in
[Stage 18](../../tools/revkit/work/stage18/README.md). Changing the DLL's
global `.tree2` to `.tree3` suffix rewrite also redirects shared dictionary
trees and fails before it reaches Kate's model tree. The binary's
`FUN_100016e0` tree reader also remains the `tree3` reader, so this experiment
does not establish that a patched DLL can load or synthesize a `tree2` voice.

A second disposable-copy experiment is documented in
[Stage 19](../../tools/revkit/work/stage19/README.md). A small converter
translated Kate's recursive tree2 nodes into the indexed layout expected by
`FUN_100016e0`; deterministic feature probes matched the legacy parser for all
non-leaf trees. The standard-engine DLL copy redirects common and voice tree
paths to `.tree2` names, with converted indexed payloads under those paths.
The eight root-leaf pitch files have opaque suffix bytes the converter does
not preserve, although the tested runtime fixtures synthesize successfully.

The first index adapter inserted a byte inside each 20-byte unit slice. GDB
stopped in adjacent-context scoring at `0x10018f17`: a resulting metric code
of 14,849 was used against Kate's 1,024-row distance table. The trace
identified a format-conversion error. The corrected adapter treats the legacy
20-byte tail as column-major, preserves all 19-byte unit records, changes
`ver.2005` to `ver.2013`, and inserts one zero-filled column between the
1-byte/7-byte columns and the three 4-byte feature groups. It checks all three
masked metric columns for all five indexes against the declared table
dimension. The added column's semantic meaning remains unknown.

With the corrected overlay, the patched standard DLL and standard Kate host
returned zero and wrote new non-silent PCM WAVs for prose, number/time, and
address/abbreviation fixtures. Each file is mono, 16 kHz, 16-bit PCM; exact
lengths and hashes are recorded in [Stage 19](../../tools/revkit/work/stage19/README.md#versioned-index-adapter-follow-up).
The requester listened to the output and reports gibberish, so these runtime
results do not establish intelligible Kate speech or compatibility. Remaining
work is to identify the semantic mismatch in adapted trees, index features,
or the front-end/model pairing. Julie's older-runtime boundary and Bridget's
missing matching runtime remain open.

After that listening report, two controlled Stage 19 comparisons narrowed the
search. Appending the opaque suffixes to the four root-leaf pitch trees leaves
all three WAVs byte-identical, so the suffix bytes do not affect the tested
runs. An alternate index mapping, inferred from the legacy reader's seven-byte
key and separate one-byte copy widths, also completes the three runs but
produces different WAV hashes and substantially different durations. The
post-load GDB trace covers all five indexes for both mappings: the original
preserves the source 0/1 values in the current reader's `+0x44` buffer and
leaves `+0x40` zero, while the reordered candidate puts zeros in `+0x44` and
the source's varying column-7 values in `+0x40`. This supports the source
column order used by the original adapter and identifies the reorder as
shifting the key bytes. The requester compared the prose samples and found
the original mapping somewhat more speech-like; the reordered sample sounded
as if it struggled to form syllables and word parts. Both remain gibberish.
A focused GDB comparison also confirms that the reorder changes the seven-byte
signatures, class candidates, and first three unit candidate pool sizes
(11/19/21 versus 30/30/30); it does not identify linguistically correct units.
The original adapter's exact feature semantics and intelligible synthesis
remain unresolved. Captures and reproduction commands are recorded in the
[Stage 19 evidence]
(../../tools/revkit/work/stage19/README.md#index-column-mapping-and-tree-suffix-controls).

A separate branch-order experiment generated copies with only `C` comparison
edges reversed, only `D` membership edges reversed, or both, leaving the
original index adapter fixed. All 17 converted trees in each variant passed
the current `tree3` structural parser. `C`-only reversal left the prose and
number/time WAVs byte-identical to the original adapter and changed address;
`D`-only reversal changed all three. Reversing both edge types matched the
`D`-only prose and address outputs but changed number/time again. These are
distinct runtime paths, not proof of a recovered branch convention. The
requester listened to the C variant and said it sounds a bit better, while
remaining garbled/gibberish. The fixture was not specified. A family and
per-file split localized the address-only C effect to `duration/cnas.tree2`;
all other duration files and the pitch family were byte-identical to baseline
in their respective controls. A GDB trace at `FUN_10013380`'s return address
`0x1001350b` recorded the 281-node tree's feature vectors, traversed nodes,
comparisons, and returned values. An equality-only variant that reverses the
three `cnas` nodes observed at value-equals-threshold cases (39, 111, and 221)
matches the full C-reversed address WAV byte-for-byte and leaves prose and
number/time identical to baseline. This is a fixture-specific observation;
it does not identify the legacy branch convention or yield intelligible
speech. Matrices, hashes, trace logs, manifests, and reproduction commands
are recorded in the [Stage 19 evidence]
(../../tools/revkit/work/stage19/README.md#c-edge-localization-and-runtime-path-trace).

Further controlled index experiments examined the 2013 reader's inserted
`attr_b` byte. Copying the preserved legacy one-byte column into it leaves all
three WAVs byte-identical to zero-fill baseline. Copying legacy signature byte
7 into `attr_b` while preserving the seven-byte signature changes all three
WAVs; moving that byte into `attr_b` and zero-padding the signature slot also
changes all three, with different hashes. A constant fill of 8 changes prose
and address but leaves number/time byte-identical. GDB confirms these fills
change per-unit costs and, for the transfer candidate, class and unit
shortlists. The legacy gen-bank column 7 distribution is close to modern
Paul's `attr_b` distribution, which is a lead for audition rather than proof
that the columns share semantics. These candidate WAVs have not been judged
for intelligibility. The Stage 19 report records full hashes, traces,
distribution comparisons, and reproduction steps.
