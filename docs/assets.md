# Asset map

The asset directories contain local samples supplied for analysis. Their source and redistribution terms are not documented here; verify those terms before adding the files to version control or sharing them.

## Executables and libraries

`binary/` contains four VoiceText host executables and four voice DLLs:

| Host executable | Voice DLL |
| --- | --- |
| `voicetext_paul.exe` | `vt_pau.dll` |
| `voicetext_julie.exe` | `vt_jul.dll` |
| `voicetext_james.exe` | `vt_jam.dll` |
| `voicetext_kate.exe` | `vt_kat.dll` |

The directory also contains Cygwin runtime DLLs, command-line utilities, SoX, and `nsh2.wav`. The role and provenance of every auxiliary file have not yet been established. Static comparison so far found matching `.text` bytes across the four voice DLLs; their voice data and small data-section differences still need separate analysis.

## Voice and shared data

- `data-paul/M16/` contains the Paul voice's model data, including four merged unit banks, their unit indexes, and duration/pitch trees.
- `data-julie/M16/` contains Julie's model banks and `tree2` duration/pitch files. Its indexes identify `ver.2009` / `VoiceText-Eng`; its four banks are `gen`, `alp`, `etc`, and `exp`.
- `data-bridget/M16/` contains Bridget's model banks and `tree3` duration/pitch files. Its indexes identify `ver.2005` / `VoiceText-Bre`; its four banks are `gen`, `alp`, `etc`, and `exp`.
- `data-kate/M16/` contains the Kate model banks, indexes, `tree2` duration/pitch files, and distance table extracted from the MSI embedded in `Kate3(1).exe`. Its indexes identify `ver.2005` / `VoiceText-Eng`; the five banks are `gen`, `gen2`, `num`, `etc`, and `alp`.
- `data-james/M16/` contains the James voice's model data, including four
  merged unit banks, their unit indexes, and duration/pitch trees. Its
  alphabetical bank is named `abc` (Paul's is `alp`).
- `data-common/` contains shared English dictionary and pronunciation resources, plus a verification file that has not been inspected as part of this analysis.

The extracted James M16 tree passed the existing DAT/UPM inventory parser for
all 284,349 unit records across its `gen`, `num`, `etc`, and `abc` banks. The
indexes have the expected 2013 English header and 19-byte records with 21
feature-column bytes; every DAT span terminated and its decoded sample count
matched the corresponding UPM entry. This confirms structural fit for the
checked parser, not whole-voice synthesis parity.

Julie and Bridget use an older index layout than that parser accepts: their
unit records are 19 bytes with 20 feature-column bytes, and their index headers
identify the 2009 English and 2005 `VoiceText-Bre` formats, respectively. A
read-only inventory checked all 1,509,574 index records and DAT/UPM spans; first,
middle, and last units in each bank reached a DAT terminator and matched their
UPM sample counts. Bridget's CAB labels its resources `Bre` (`dict-bre` and
`vt_bre_bridget16.dll`); the CAB did not establish an `en-CA` locale. Julie's
and Kate's `tree2` files are not accepted by the repository's `tree3` parser
and need a dedicated `tree2` format check. These checks do not establish
whole-voice synthesis compatibility.

- `include/` contains `vt_eng.h` and `slicense_eng.h`, headers distributed with the sample package.

The local model trees currently represent Paul, James, Julie, Bridget, and
Kate. The James tree was copied from the extracted Neospeech James package in
Downloads; its notes identify it as the 32-bit VW James voice. Julie and
Bridget's model trees were extracted from their CABs without running the installers. Kate's
model files were extracted from the embedded MSI without running either the
outer setup or MSI installer. The Kate indexes use the older 2005 English
layout: five banks totaling 283,696 units, 19-byte records, and 20 feature
bytes per unit. All indexed DAT/UPM spans are in bounds; first, middle, and last
units in each bank reached a DAT terminator and matched their UPM sample counts.
This is a structural sample check, not a full-corpus decode or proof that the
Kate DLL can synthesize from this model. The MSI also carried common dictionary
and tree resources; those were not copied into this voice-specific directory.

## Git handling

The root `.gitignore` keeps `binary/`, `data-common/`, and all per-voice model
directories local. This avoids accidentally treating the vendor payloads
as project-owned, redistributable source. `include/` is not ignored; its origin
and licensing should be reviewed before publication.
