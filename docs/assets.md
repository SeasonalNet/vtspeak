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
- `data-common/` contains shared English dictionary and pronunciation resources, plus a verification file that has not been inspected as part of this analysis.
- `include/` contains `vt_eng.h` and `slicense_eng.h`, headers distributed with the sample package.

Only the Paul model is currently represented in `data-paul/`. Other voices may need their matching model directories for comparative analysis.

## Git handling

The root `.gitignore` keeps `binary/`, `data-common/`, and `data-paul/` local. This avoids accidentally treating the vendor payloads as project-owned, redistributable source. `include/` is not ignored; its origin and licensing should be reviewed before publication.
