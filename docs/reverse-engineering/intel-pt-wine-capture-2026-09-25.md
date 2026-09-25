# Intel PT trace of the VoiceText Wine run (2026-09-25)

## Capture and decode

### Corrected paired rerun

After correcting the Stage 13 harness to keep the Stage 5 fixture at its
expected `/work/stage5` path, both captures completed with the GDB inferior
exiting normally. The control logged 14 `ACOUSTIC_DESC` records; the
`acoustic-83-zero` run logged `TPP_NUMERIC_ABLATE83 code=83` and 14 descriptor
records. Both sidecar WAVs are valid 16 kHz mono 16-bit PCM and differ from the
Stage 5 seed WAV. Their SHA-256 values match the previously validated Stage 12
outputs: control `29052da3a11213f4616c30b2da31ccd9c4400d657afa1c163a51135653ddad09`,
zero `f8aa3133ef5701e78f46db0bd55cdcf3674430c988d7b6f5834bd804a0683305`.

The corrected control trace was 800.862 MB and the intervention trace was
811.204 MB. Their VoiceText PIDs were 135604 and 135970; both mapped
`vt_pau.dll` at `0x10000000`. A branch-only decode was attempted for each
VoiceText runtime window, but the AUX scan did not finish during this review.
No paired branch-edge counts or path-difference claims are made from these
captures. The prior 591-row branch sample below belongs to the earlier control
capture and is not paired with the corrected intervention run.

The earlier host-side Stage 13 capture ran the existing Stage 12
`acoustic-control` probe for `numeric-context-g83-can-anybody-help`. The
resulting `/tmp/vtspeak-g83-control.perf.data` was 683.171 MB. Its sidecar WAV
is a valid 16 kHz, mono, 16-bit PCM file. The sidecar GDB log says the remote
inferior exited with code 1, so the WAV establishes that output was produced;
it does not establish a clean debugger exit.

The captured mapping names the vendor DLL `/samples/vt_pau.dll`, which is a
container path and is absent on the host. Pointing `perf script` at the local
binary through `--symfs` allowed perf to decode branches. The small filtered
trace is preserved as
[`g83-control-vt-pau-branches.txt`](../../tools/revkit/work/stage13/g83-control-vt-pau-branches.txt).
It contains 591 decoded branch records from VoiceText PID 112937, with 96
distinct DLL source addresses, over 3.764 seconds. The DLL was mapped at
`0x10000000`, so trace virtual addresses convert to RVAs by subtracting that
base. The PID is specific to this run.

## Cross-checks against static analysis

The highest-frequency source addresses cluster around the output helper family
already decompiled in Stage 6 ([formatted-output pseudocode](../../tools/revkit/work/reports/stage6-final-expansion-helpers.txt),
[buffer helper pseudocode](../../tools/revkit/work/reports/stage6-tpp-format-parser.txt)):

| DLL RVA | Branch records | Static correlation |
| --- | ---: | --- |
| `0x65b69`, `0x65b7c`, `0x65b86` | 57 each | Within `FUN_10065b52` at VA `0x10065b52`; its pseudocode appends one byte to a caller buffer or refills the buffer through `FUN_1006529c` |
| `0x65853` | 42 | Within formatted-output routine `FUN_100653b4` at VA `0x100653b4` |
| `0x65bd8`, `0x65bea` | 42, 39 | Inside `FUN_10065bb8`, the byte-by-byte string output helper called by `FUN_100653b4` |
| `0x688f1` | 32 | Within `FUN_10068899`; its pseudocode initializes a descriptor table and imports startup handles |
| `0x2a6f4` | 4 | Inside exported `VT_GetTTSInfo_ENG`; all four observed edges go to `FUN_10024cc0` |

These are address and control-flow correlations, not recovered source names.
The hot output-helper activity shows that the formatted-output path executes in
this run. It does not by itself show that those bytes are speech samples,
identify the formatted strings, or prove a semantic role for G83.
This filtered trace did not yield a decoded branch row whose source was in
`voicetext_paul.exe`; the observed VoiceText rows were in `vt_pau.dll`.

## Decode quality and limits

The first decode without the host-side image mapping produced a 1.6 GB text
file dominated by instruction-trace decode errors. The corrected
filtered decode reports one AUX-loss warning and 2,174 trace errors across the
captured process tree; only 22 were attributed to the VoiceText process. The
capture also includes GDB and Wine helper processes, so aggregate error counts
do not describe the VoiceText process alone. The `Trace doesn't match
instruction` diagnostics mean some portions could not be reconstructed against
the mapped code image. Treat the 591 branch rows as observed decoded coverage,
not as a complete execution trace.

The corrected pair now supplies the control and intervention captures, but
their paired branch-edge comparison remains open. A future bounded decode
should compare the same VoiceText runtime windows. GDB memory/register captures
and the existing audio comparison remain necessary to explain any path
differences. The intervention capture command is:

```sh
tools/revkit/work/stage13/capture-wine-intel-pt.sh \
  /tmp/vtspeak-g83-zero.perf.data \
  acoustic-83-zero numeric-context-g83-can-anybody-help
```
