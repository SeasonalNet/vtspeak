# vtspeak

`vtspeak` is a research and reimplementation workspace for understanding the legacy VoiceText text-to-speech engine and its voice data.

- The current work combines static analysis of the Windows executable and DLLs with controlled runtime traces of the model readers and synthesis path.
- A native POSIX reimplementation is the ultimate goal; this repository does not provide one yet.

## Repository map

| Path | Contents |
| --- | --- |
| `binary/` | Local VoiceText executables, DLLs, runtime dependencies, and a WAV sample |
| `data-bridget/` | Bridget voice model, unit indexes, waveform banks, and prosody trees |
| `data-common/` | Shared pronunciation and dictionary resources |
| `data-james/` | James voice model, unit indexes, waveform banks, and prosody trees |
| `data-kate/` | Kate voice model, unit indexes, waveform banks, and prosody trees |
| `data-julie/` | Julie voice model, unit indexes, waveform banks, and prosody trees |
| `data-paul/` | Paul voice model, unit indexes, waveform banks, and prosody trees |
| `include/` | Headers recovered with the VoiceText package |
| `tools/revkit/` | Docker-based Ghidra, FLOSS, capa, and PE inspection toolkit |
| `docs/` | Project notes, asset map, and reverse-engineering findings |

The binaries and voice data remain in the working tree for local analysis. Their redistribution terms have not been established, so the root `.gitignore` excludes those input directories from Git.

> [!NOTE]
> No license is asserted for the vendor engine or voice assets. Treat them as proprietary.

## Start the analysis toolkit

Follow [the revkit instructions](tools/revkit/README.md). Run Docker Compose from `tools/revkit`; its bind mounts expect the repository layout shown above. Portable trace evidence, small fixtures, and cited decompiler reports under `tools/revkit/work/` are tracked. Ghidra project databases, caches, Wine prefixes, and full local copies of vendor inputs remain ignored. Narrative findings intended for review live under `docs/reverse-engineering/`.

## Current findings

Please see [docs/INDEX.md](./docs/INDEX.md) for the list of documentation files.

## Working principles

- Treat vendor binaries and model files as read-only analysis inputs.
- Label observations separately from hypotheses. Ghidra output is decompiler pseudocode, not recovered original source.
- Preserve function addresses, filenames, and byte offsets when recording format claims.
- Don't claim a standard codec or a compatible reimplementation until the behavior has been checked against known outputs.

## Find your role

- Agent or other autonomous system -> [AGENTS.md](./AGENTS.md)
- Person -> [README.md](./README.md)

## License

All reimplementation work (excluding vendor data) in this repository is licensed under the GNU AGPLv3 license. See [LICENSE](./LICENSE).
