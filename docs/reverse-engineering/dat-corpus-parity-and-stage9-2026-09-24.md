# Local DAT corpus parity and additional synthesis boundaries (2026-09-24)

These checks use the local 2013 M16 Paul package. The DLL, indexes, DAT and
UPM banks were read as supplied. Their SHA-256 values, along with the probe
source, PE32 build, and ignored result stream, are recorded in the
[input/result manifest](../../tools/revkit/work/reports/dat-corpus-parity-2026-09-24.sha256).
The decompiler addresses and signatures below are pseudocode-level analysis;
the runtime capture is the independent behavior check.

## DAT decoder corpus comparison

The [PE32 probe](../../tools/revkit/scripts/probe_dll_decoder.c) loads the
unchanged `vt_pau.dll`, calls its public `VT_LOADTTS_ENG` entry with the
speaker and database-path values observed in the original executable, and
requires a successful public `VT_TextToFile_ENG` WAVE synthesis before
probing the decoder. That canary produced a valid 16 kHz mono 16-bit WAVE.
The probe then calls the loaded DLL's `FUN_10001b30` at module offset
`0x1b30` using the four arguments established by the Stage 2 GDB capture:
compressed DAT pointer, PCM output pointer, compressed byte count, and output
byte-count pointer. Its first captured-payload call returned 3,036 PCM bytes
identical to `stage2-copy/candidate-00.pcm`.

For the full run, the probe reads each versioned `unit-*.idx` 19-byte payload
record at header offset `0x2d` and its matching `merged-*.dat` span. It checks
the DLL's byte count against four times the sum of the combined UPM bytes,
which is the 16-bit PCM byte count implied by the doubled UPM periods. It
records a SHA-256 of each returned PCM buffer in bank/record order. The
[independent comparator](../../tools/revkit/scripts/compare_dat_corpus.py)
decodes each DAT span with `decode_dat.py`, checks the byte count, and compares
the Python SHA-256 with the DLL result. The complete result stream is kept
under ignored `tools/revkit/work/corpus-parity/`; the small hash manifest is
tracked.

| Bank | Units | DLL/local PCM byte count and SHA-256 pairs |
| --- | ---: | ---: |
| `gen` | 440,124 | 440,124 matched |
| `num` | 24,508 | 24,508 matched |
| `etc` | 115,723 | 115,723 matched |
| `alp` | 119 | 119 matched |
| **Total** | **580,474** | **580,474 matched; zero mismatches** |

The DLL result stream was 20,897,073 bytes: a nine-byte marker and 36 bytes
per unit. Its SHA-256 is
`8a9d212a0b0a0f4d7d3820f533a7970bb31ec9cef68c440b2f2e5898c62eb084`.
This is local full-corpus PCM value parity at the decoder boundary, evidenced
by equal byte counts and SHA-256 digests. The 27 earlier captured payloads
remain exact byte-for-byte comparisons. A separate full stream scan found no
mode-8 frame in this package, so mode-8 runtime behavior is still not tested.
These checks do not establish final WAVE parity for every unit combination or
behavior for another package version or voice.

From the repository root, with local vendor inputs present:

```sh
bash tools/revkit/work/scripts/build_decoder_probe.sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/scripts/run_decoder_probe.sh corpus
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/scripts/compare_dat_corpus.py
```

The PE32 builder is isolated in `tools/revkit/Dockerfile.pe32-builder`; the
existing Wine runtime image, read-only vendor mounts, and disabled container
networking are unchanged.

## Additional synthesis boundary checks

[Stage 9](../../tools/revkit/work/stage9/README.md) reused the Stage 6 number
and abbreviation input files. Under the original executable, GDB counted
`FUN_1002c220` timeline rows and appended every PCM buffer returned by
`FUN_10026870`. The concatenated buffers exactly equal the WAVE data chunk
for both inputs. Each resulting WAVE is also byte-identical to its existing
Stage 6 reference.

| Input | Timeline row counts by call | PCM blocks | WAVE frames | Captured PCM bytes |
| --- | --- | ---: | ---: | ---: |
| Time, decimal, and units | 22, 78 | 4 | 108,601 | 217,202 |
| Honorific, address, and telephone | 41, 23, 17 | 5 | 116,689 | 233,378 |

This extends the sample-buffer-to-WAVE comparison to two different text paths
and multiple chunk boundaries. It does not exercise every context-gate flag
combination or prove general whole-synthesis parity.

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage9/run.sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/scripts/verify_stage9.py
```
