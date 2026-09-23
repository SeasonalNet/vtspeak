# vtspeak

`vtspeak` is a research and reimplementation workspace for understanding the legacy VoiceText text-to-speech engine and its voice data.

- The current work focuses on static analysis of the Windows executable and DLLs, the model file readers, and the synthesis path.
- A native POSIX reimplementation is the ultimate goal; this repository does not provide one yet.

## Repository map

| Path | Contents |
| --- | --- |
| `binary/` | Local VoiceText executables, DLLs, runtime dependencies, and a WAV sample |
| `data-common/` | Shared pronunciation and dictionary resources |
| `data-paul/` | Paul voice model, unit indexes, waveform banks, and prosody trees |
| `include/` | Headers recovered with the VoiceText package |
| `tools/revkit/` | Docker-based Ghidra, FLOSS, capa, and PE inspection toolkit |
| `docs/` | Project notes, asset map, and reverse-engineering findings |

The binaries and voice data remain in the working tree for local analysis. Their redistribution terms have not been established, so the root `.gitignore` excludes those input directories from Git.

[!NOTE]
No license is asserted for the vendor engine or voice assets. Treat them as proprietary.

## Start the analysis toolkit

Follow [the revkit instructions](tools/revkit/README.md). Run Docker Compose from `tools/revkit`; its bind mounts expect the repository layout shown above. Ghidra projects, decompiler exports, and other generated analysis files are stored under `tools/revkit/work/`, which is ignored by Git. Narrative findings intended for review live under `docs/reverse-engineering/`.

## Current findings

Please see [docs/INDEX.md](./docs/INDEX.md) for a list of research documentation files.

## Working principles

- Treat vendor binaries and model files as read-only analysis inputs.
- Label observations separately from hypotheses. Ghidra output is decompiler pseudocode, not recovered original source.
- Preserve function addresses, filenames, and byte offsets when recording format claims.
- Don't claim a standard codec or a compatible reimplementation until the behavior has been checked against known outputs.

## Find your role

- Agent or other autonomous system -> [AGENTS.md](./AGENTS.md)
- Person -> [README.md](./README.md)

## License

All reimplementation work (excluding vendor data) in this repository is licensed under the GNU AGPLv3 license.
