# Voice package generation comparison (2026-09-25)

This note compares the local Paul and James M16 models with the older Julie,
Bridget, and Kate models. “Generation” here refers to observed package and
index markers such as `ver.2013`, `ver.2009`, and `ver.2005`; those markers do
not by themselves establish an engine release date or a supported
voice-to-engine compatibility matrix.

## Model data observed

| Voice | Index header | Banks and indexed units | Prosody trees | Current validation |
| --- | --- | --- | --- | --- |
| Paul | `ver.2013`, `VoiceText-Eng` | `gen`, `num`, `etc`, `alp`; 580,474 units | 17 `tree3` files | All four indexes match the 2013 layout. The local DAT decoder was compared with the DLL across all indexed records at the decoder-output boundary. Selected synthesis paths were also checked. |
| James | `ver.2013`, `VoiceText-Eng` | `gen`, `num`, `etc`, `abc`; 284,349 units | Duration and pitch trees are present | All index records and DAT/UPM extents passed the available inventory checks; decoded sample counts matched UPM entries. Stage 17 also produced valid WAVs for three shared text fixtures with the supplied James host/DLL. |
| Julie | `ver.2009`, `VoiceText-Eng` | `gen`, `alp`, `etc`, `exp`; 728,692 units | 17 `tree2` files | All indexes and ordered DAT/UPM spans were checked. First, middle, and last units in each bank were decoded and matched their UPM sample counts. Thirteen trees consume EOF; `bf`, `nbf`, `qbf`, and `sbf` pitch trees have opaque suffixes. |
| Bridget | `ver.2005`, `VoiceText-Bre` | `gen`, `alp`, `etc`, `exp`; 780,882 units | 17 `tree3` files | All indexes and ordered DAT/UPM spans were checked. First, middle, and last units in each bank were decoded and matched their UPM sample counts. All 17 trees parse with the existing `tree3` parser. |
| Kate | `ver.2005`, `VoiceText-Eng` | `gen`, `gen2`, `num`, `etc`, `alp`; 283,696 units | 17 `tree2` files | All indexes and ordered DAT/UPM spans were checked. First, middle, and last units in each bank were decoded and matched their UPM sample counts. Thirteen trees consume EOF; the same four pitch tree names have opaque suffixes. |

The 2013 Paul and James indexes have 19-byte per-unit records followed by 21
feature-column bytes per unit. The 2005/2009 Julie, Bridget, and Kate indexes
have 19-byte records followed by 20 feature-column bytes per unit. The legacy
index inspector checks complete index extents, all DAT/UPM spans for bounds
and contiguity through EOF, and sampled payload counts. Those checks establish
structural fit and sampled consistency, not whole-voice synthesis support.

The tree formats do not line up one-to-one with index years: Julie and Kate
have `tree2`, Bridget has `tree3` despite its `ver.2005` indexes, and Paul and
James have `tree3`. The Kate MSI's extracted 2006 `vt_eng.dll` contains a
recursive tree reader at `0x100015d0` and child loader at `0x10001740`. This
supports the serialized `tree2` node layout documented in the [format
findings](voice-engine-and-model-formats.md#legacy-tree2-decision-trees). Eight
Julie/Kate pitch files currently parse only as a 9-byte root leaf plus
unexplained suffix; they are not treated as fully recovered files.

## Engine evidence and compatibility boundary

The supplied `binary/` set contains Paul, James, Julie, and Kate voice DLLs.
Their executable `.text` bytes match in the recorded comparison, while their
`.data` sections and voice/model inputs differ. This is evidence that those
four supplied DLLs share the same executable code image. It is not proof that
separately extracted old CAB/MSI model packages load correctly in those DLLs,
or that their voice behavior is identical. Bridget's CAB resource labels
include `Bre` and `vt_bre_bridget16.dll`; the CAB does not establish an
`en-CA` locale, and that DLL is not part of the four-DLL comparison above.

The model sources also differ: James came from an extracted Neospeech James
package whose notes identify a 32-bit VW voice; Julie and Bridget were
extracted from CABs without running installers; Kate's voice files came from
the MSI embedded in `Kate3(1).exe`. Kate's MSI also carried common dictionary
and tree resources that were not copied into `data-kate/`. These source and
layout differences make it important to keep package-specific resources
together when testing a voice.

Paul still has the broadest engine evidence: full local DAT decoder-output
comparison and selected synthesis-path checks. Stage 17 adds three successful
shared-input runs for James, so it now has direct short-synthesis evidence in
addition to broad data-structure validation. Julie and Kate reach the first
voice-tree request with the supplied DLLs but fail because those DLLs request
`tree3` while the local packages provide `tree2`. Bridget has structural
coverage; no host is present, and the only external DLL candidate found is
stored in a directory named `crack` and was not accepted for runtime evidence.
These results do not establish older-engine compatibility or whole-synthesis
parity for any non-Paul package.

A second DLL supplied with the James package also synthesizes the same three
fixtures with the supplied James host. It produces the same byte and frame
counts as `binary/vt_jam.dll`, but different PCM and whole-file hashes. This
supports successful short synthesis across those two DLLs; it does not
establish waveform parity. The extra DLL and results are documented in Stage
17.

Stage 17 is an exhaustive runtime attempt over the four host/DLL pairs actually
present in `binary/`: Paul, James, Julie, and Kate. The matrix used three
shared text classes for Paul and James; Julie and Kate were stopped at the
same first missing voice-tree resource on a prose fixture, so further text
classes would not reach synthesis in this configuration. The report records
the six successful WAVs' formats, frame counts, hashes, and the exact
Julie/Kate failed-open paths.
See the [Lead 5 compatibility report](lead5-voice-package-runtime-2026-09-25.md)
and [Stage 17 evidence](../../tools/revkit/work/stage17/README.md).

Stages 18 and 19 are isolated Kate compatibility experiments. Stage 18's
global suffix change redirects common dictionary tree requests. Stage 19
converts observed `tree2` structures and indexes in a disposable package copy;
the patched standard DLL completes three fixture runs and writes valid PCM,
but the requester reports the Kate output is gibberish. A post-load trace over
all five indexes supports the original adapter's preserved 1-byte/7-byte
source-column order over the reordered candidate, but does not identify the
remaining feature semantics or prove source-to-engine equivalence. No matching
older runtime pair or Bridget runtime pair is available in this checkout.

## Sources and tools

- [Asset map and handling notes](../assets.md) records local package contents
  and extraction provenance.
- [VoiceText engine and model format findings](voice-engine-and-model-formats.md)
  records the observed index, DAT/UPM, and tree layouts and their boundaries.
- `tools/revkit/scripts/inspect_unit_idx.py` checks the 2013 index layout;
  `tools/revkit/scripts/inspect_legacy_unit_idx.py` checks the 2005/2009
  layout and sampled DAT/UPM consistency.
- `tools/revkit/scripts/tree3.py` parses `tree3` files;
  `tools/revkit/scripts/tree2.py` parses observed recursive `tree2` prefixes
  and reports unconsumed suffixes.
- `tools/revkit/work/stage17/README.md` documents the controlled Wine matrix;
  `tools/revkit/work/stage18/README.md` and
  `tools/revkit/work/stage19/README.md` document the isolated Kate adapter
  experiments and their limits.
