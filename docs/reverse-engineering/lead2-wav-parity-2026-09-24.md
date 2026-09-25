# Lead 2: bounded WAV parity matrix (2026-09-24)

## Result

The bounded local matrix is complete for eight selected texts and the three
individual API controls already traced in Stage 8. The local runtime produced
valid 16 kHz mono 16-bit PCM WAVs for every text. For all eight cases, the
concatenated buffers returned by `FUN_10026870` exactly equal the WAV data
chunk, and each whole WAV has the same byte length as Wag's peer manifest row.
Every local whole-file SHA-256 differs from the peer SHA-256. Thus this probe
does not reproduce Wag's byte-exact results, and the equal lengths alone do
not establish audio parity.

The peer bundle contains hashes and byte lengths, but not the peer WAVs,
engine DLL, or voice data. Because the peer bytes are unavailable, a hash
mismatch cannot be localized to a particular WAV header field or PCM sample.
Our trace localizes the local side through selected unit IDs, timeline row
counts, returned PCM blocks, and WAV data assembly; it cannot show which of
those stages first differs from Wag's run. Wag's DLL license patch is
user-relayed peer information, not independently inspected evidence.

## Text matrix

Inputs are the exact CP1252 byte strings in the selected rows of Wag's
`golden_speak_sha256.tsv` from
`vtpaul_codex_status_2026-09-24-new-1.zip`. Local hashes, sizes, and peer
comparisons are recorded in
[`local-golden.tsv`](../../tools/revkit/work/stage11/local-golden.tsv); source
peer values and exact input hex are in
[`peer-golden.tsv`](../../tools/revkit/work/stage11/peer-golden.tsv). Every
row below had `local WAV bytes == peer WAV bytes`, `whole-WAV SHA-256 match =
no`, and `captured PCM blocks == local WAV data = yes`.

| Coverage | Fixture | Frames | Timeline rows by call | Selected-unit callbacks | PCM blocks |
| --- | --- | ---: | --- | ---: | ---: |
| Ordinary prose | `prose-hello-world` | 11,803 | 10 | 16 | 1 |
| Numbers and quantities | `numbers-mixed` | 98,719 | 14, 15, 49 | 122 | 4 |
| Calendar date and decimal | `date-calendar` | 73,990 | 73 | 123 | 3 |
| Time range and zone abbreviation | `time-range` | 73,506 | 58 | 110 | 3 |
| Honorific, street abbreviation, place name | `abbrev-name-street` | 69,872 | 66 | 112 | 3 |
| Address, apartment, postal code | `address-street` | 118,548 | 51, 17, 18, 18 | 153 | 5 |
| Short exclamation | `edge-exclamation` | 15,571 | 10 | 18 | 1 |
| Unfinished trailing phrase | `edge-trailing-partial` | 54,020 | 44 | 82 | 2 |

The `prose-hello-world` fixture was run twice in separate runtime invocations;
both local WAV SHA-256 values were
`a9bb244d9d0cdb664a7a64d14eeb2acd0c45b22d19383d88ff157ba337dd1a69`.
This confirms repeatability for that local fixture and invocation setup, not
for other systems or engine builds.

## API controls

Stage 8 independently varies one API control at a time on the same longer
prose input (`Hello from the VoiceText Stage 1 runtime check.`), using
sentinels for the other controls. All outputs are valid mono, 16-bit, 16 kHz
PCM WAVs. These controls demonstrate local behavior, not peer parity: Wag's
manifest does not include corresponding control WAV hashes.

| Local setting | Frames | WAV bytes | Whole-WAV SHA-256 | Observed effect |
| --- | ---: | ---: | --- | --- |
| Defaults (`-1, -1, -1`) | 44,717 | 89,478 | `3bd8bbfde92f1d645715de40a03a6a68e3acedf0a3361158b869b8a07d803187` | Reference for these controls |
| Pitch `120` | 47,156 | 94,356 | `3e7dec025338b047b957c6466fba07d42b83ff734a7b1a67b0cead40962e2079` | Changes frame count and PCM |
| Speed `120` | 38,326 | 76,696 | `cf3af2a088bdb8e6a80e5faf574bd25304e416d6cfdc3d74a0026cb503c46379` | Changes frame count and PCM |
| Volume `120` | 44,717 | 89,478 | `1a4e192094b725de996e7208bb1e6d3af8246e33adeeb3b8a240fd874f6a7c99` | Preserves frame count; PCM peak changes from 30,720 to 18,432 |

The default and pitch WAV data were previously compared byte-for-byte with
the Stage 3 and Stage 4 reference PCM respectively. Stage 8 recorded the
control arguments and timeline/PCM path; its verifier checks WAV format and
frame counts for all four settings. The Stage 11 verifier checks complete
PCM-block-to-WAV-data equality for all eight text fixtures.

## Reproduction and evidence

Run from the repository root with local vendor inputs and the existing
`vtspeak-stage2-wine` image:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage11/run.sh \
  prose-hello-world numbers-mixed date-calendar time-range \
  abbrev-name-street address-street edge-exclamation edge-trailing-partial
python3 tools/revkit/work/scripts/verify_stage11.py
python3 tools/revkit/work/scripts/verify_stage8.py
```

`run.sh` saves WAVs, PCM blocks, and GDB logs under `tools/revkit/work/stage11/`.
Its exit trap restores the prior Stage 5 input and output files. The final
run's restoration was checked byte-for-byte against its saved pre-run copies.
`trace.gdb` records callbacks at `0x1001b200` (selected-unit IDs),
`0x1002c220` (timeline rows), and `0x10026870` (PCM blocks).

These results close the bounded lead-2 probe with a negative whole-WAV parity
result against the supplied peer manifest. Full corpus parity, peer/local
divergence localization, patched-license behavior, and equivalence across
different engine DLLs remain unverified.

Wag's later peer reply reports a separate 186-text Stage 10 set with
byte-exact engine-to-Rust WAV parity. Those are not the eight Stage 11 texts
in this matrix, and the reply bundle does not include their engine WAV bytes
or the DLL/data needed to reproduce the peer run. That peer result therefore
does not change the negative hashes recorded here. It is consistent with
their independent unit-selection field comparison for the Stage 10
neighbor-context fixtures, while this Lead 2 matrix remains an unreproduced
cross-build comparison on its own eight inputs.
