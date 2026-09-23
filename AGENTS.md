# Repository instructions

## Scope

This repository is a reverse-engineering research and reimplementation workspace for the legacy VoiceText engine and its model data. Keep changes focused on understanding the binaries, documenting evidence, and building independent analysis tooling.

## Preserve analysis inputs

- Treat `binary/`, `data-common/`, and `data-paul/` as read-only vendor inputs. Do not patch, normalize, rename, repack, or overwrite their files.
- Keep generated Ghidra projects, caches, and raw exports under `tools/revkit/work/`; that directory is intentionally ignored.
- Put curated, reviewable explanations under `docs/` and link them from the root README.
- Do not relicense or attempt to relicense vendor work (e.g. model data, voice corpus), it should currently be assumed and treated as proprietary.

## Reverse-engineering evidence

- Distinguish direct observations from inferences.
- For binary layouts, record the reader function/address, widths, endianness when established, and a sample file/offset that supports the claim.
- Describe Ghidra output as pseudocode. Decompiled names, types, and control flow may be imperfect.
- Validate proposed parsers and decoder behavior against multiple files or known engine output before describing a format as fully recovered.

## Repository safety

- Inspect `git status` before and after edits; preserve existing user work.
- Do not commit, push, tag, publish, or open a pull request unless explicitly asked.
- Do not add large vendor inputs or generated projects to Git. The root and toolkit ignore rules reflect this boundary.
- Do not weaken security controls or silence quality tools to make an analysis command succeed.
