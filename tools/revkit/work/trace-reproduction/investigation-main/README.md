# Recovery of the final historical WAVs

The final 18 recordings have exact recipes in the main runtime catalogue and
artifact manifest. `recovered-verification.json` maps each original path to its
verified runtime batch or exact sample composition. No original recording is
read as a generation input, and no original audio is preserved under another
encoding.

## Actual inputs and isolation

Historical runners reused `stage5/input1.txt` and `output.wav`, including
across simultaneous investigations. Named fixture arguments therefore did not
always describe the text actually read by a process. For example, the original
split-helper and Bridget logs contain the full `Hello world.` phone stream,
although their commands name a `Hello.` fixture. Replays with the former text
match the original files exactly.

Recovered cases include:

- `split-threshold-new`, `apple-generic-tail-map-selector`, the patched Apple
  selector file, two Bridget captures, and two Stage 19 captures: actual
  `Hello world.` input.
- The Hello-labelled marker file and decimal-labelled parser file: actual
  `Apple.` input. Their original complete WAV hashes match isolated Apple
  replays, including the original parser/selection instrumentation.
- The minimal assembly capture and one key-repack capture: concurrent Stage 16
  VTML marker fixtures. The preserved literal inputs and timestamp evidence
  are documented in the sibling investigation directories.
- `hi-weighted-coverage`: the intended `Hi.` input with the complete generated
  model parent and matching index overlay. The previous attempted recipe had
  failed at container mount setup, before synthesis.

The catalogue pins the actual literal texts, DLLs, model overlays, trace
versions, and API controls. Every new batch runs in a separate scratch tree;
parallel verification does not share synthesis fixtures or a Wine prefix.

## Complete historical file bytes

Three recordings combine data from independently generated outputs. Their
RIFF/data lengths alone do not describe their complete physical files. Exact
frame recipes preserve those historical bytes without depending on a timing
race:

- **Apple metric file, 30,750 bytes:** Kate Apple header, then 15,984 bytes of
  Paul Apple PCM, then 14,722 bytes of Kate Apple PCM. Both complete inputs are
  independently synthesized from `Apple.`. Concatenation reproduces the full
  original SHA-256 `596c800226a4f021ed1b886e106f86b169fcb001146831f5777b325c5a90711f`.
- **Original MSI address, 429,028 bytes:** clean-address header, explicit
  demo-address PCM, then clean-address PCM. The explicit banner text removes
  the unpatched DLL's variable demo-message choice. The legacy investigation
  records the exact strings and component hashes.
- **Context-tail file, 516,264 bytes:** its original header, two spans from the
  independently generated Paul reading-rule baseline, then spans from Kate
  `Hello world.` with eight repeated gaps. The manifest records every frame
  offset/count and each silence length. Reconstructed SHA-256:
  `148693ada20cca0b0cb0d5926eae3b511dfd778115d3d5714c11b9b9fe164a15`.

These compositions prove exact reconstruction of the recorded artifacts.
They do not claim that a deterministic engine invocation by itself emits the
mixed files, or that the exact historical scheduling of competing writers has
been reproduced. Runtime sources, not recorded audio dumps or residual-byte
patches, supply every PCM sample.

## Validation

All 16 integrated recovery batches passed fresh scratch execution through
`reproduce_wavs.py run`, and all 149 catalogue batches passed byte verification.
After pruning, all 1,600 WAVE manifest recipes still passed verification.
The Python unit suite passed 24 tests; Go tests and Go vet passed. Stage 8 and
Stage 9 offline capture checks passed after restoring their 16 required WAVs,
which were then hash-checked and cleaned again. Runtime, artifact, and offline
check summaries are retained beside the main catalogue.
