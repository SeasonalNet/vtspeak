# vtspeak

`vtspeak` is a research and reimplementation workspace for understanding the legacy VoiceText text-to-speech engine and its voice data.

- The current work combines static analysis of the Windows executable and DLLs with controlled runtime traces of the model readers and synthesis path.
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

> [!NOTE]
> No license is asserted for the vendor engine or voice assets. Treat them as proprietary.

## Start the analysis toolkit

Follow [the revkit instructions](tools/revkit/README.md). Run Docker Compose from `tools/revkit`; its bind mounts expect the repository layout shown above. Portable trace evidence, small fixtures, and cited decompiler reports under `tools/revkit/work/` are tracked. Ghidra project databases, caches, Wine prefixes, and full local copies of vendor inputs remain ignored. Narrative findings intended for review live under `docs/reverse-engineering/`.

## Current findings

Please see [docs/INDEX.md](./docs/INDEX.md) for a list of research documentation files.
The [follow-up investigation leads](./docs/reverse-engineering/follow-up-investigation-leads-2026-09-24.md)
summarize proposed next probes and their evidence boundaries.
The current [controlled neighbor-context results](./docs/reverse-engineering/tree-context-neighbor-contrasts-2026-09-24.md)
record the Stage 10 runtime traces and candidate-scoring findings.
The [lead 1 closure](./docs/reverse-engineering/lead1-tree-input-closure-2026-09-24.md)
summarizes the verified tree inputs, downstream effects, and limits for the local
Paul package. Its [peer checkpoint](./docs/reverse-engineering/wag-lead1-checkpoint-2026-09-24.md)
lists the numeric results for Wag's independent vtpaul agents to compare.
The [lead 2 WAV parity matrix](./docs/reverse-engineering/lead2-wav-parity-2026-09-24.md)
records the local text-category captures, peer-manifest mismatches, and API
control evidence.
The [lead 3 TPP findings](./docs/reverse-engineering/lead3-tpp-typed-code-findings-2026-09-24.md)
record the corpus associations, verified runtime writes, and remaining label
limits for typed pronunciation codes.
The [reply to Wag's latest peer response](./docs/reverse-engineering/wag-peer-response-new-2-2026-09-24.md)
records the Lead 1 cross-check, Lead 2 parity scope, and new E-record result.

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
