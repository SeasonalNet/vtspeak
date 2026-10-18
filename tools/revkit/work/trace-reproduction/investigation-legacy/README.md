# Retained legacy WAV recovery, 2026-10-18

Every recovery below was synthesized from local DLL/model inputs in isolated,
no-network Docker scratch storage. Original WAVs were read only for comparison.
No stored audio payload is used as a generation source. Full-file comparisons
include RIFF headers and bytes beyond the declared data chunk.

## Stage 19 concurrent fixture captures

The nominal input names do not describe these historical recordings. Replaying
with the actual concurrent Stage 5 input reproduces the complete bytes:

| Retained Stage 19 path | Actual fixture | Bytes | SHA-256 |
| --- | --- | ---: | --- |
| `pih0-context.wav` | `Hello world.\n` | 22084 | `016b5489fafc55ad4d2355884321eae6a81733a877e613ad2eab362e08c6a1ff` |
| `selection-hello-copy-assembly.wav` | `Hello world.\n` | 16134 | `b1627ad8141d22f7458c60e5b33d917e02b5e97bf1726c0ea1b76cf7b3651520` |
| `selection-hello-assembly-entry-only.wav` | Stage 16 `input-ex-records-vtml-mark-boundary.txt` | 5518 | `61bdf41f283af0c9a048b6628b76ae1bc4a292000b77a03e3f83a045eec2e714` |

The corresponding proposed batch JSONs preserve the original runner, trace,
DLL, tree, index and pitch-overlay parameters. `hello-world.txt` pins the first
two input bytes independently of the mutable Stage 5 fixture.

The historical minimal assembly trace ran at 2026-09-27 19:57 UTC. A concurrent
session launched Stage 16 `run-buffer-ex-mark-options.sh 1` at 19:56:25 UTC,
substituting the marker-boundary input in the shared Stage 5 directory. The
isolated replay with this exact input proves that the earlier explanation
that the assembly breakpoint itself changed synthesis output was unsupported.
The differing fixtures explain these three retained recordings.

Verification scratch outputs are in
`/tmp/legacy-stage19-helloworld-v2/{pih0-context,hello-copy-assembly}/work/stage19/`
and `/tmp/legacy-stage19-mark-boundary-v2/hello-assembly-entry-only/work/stage19/`.
Disposable copied models and Wine prefixes were removed after verification;
small traces, generated outputs, and proposed parameter JSONs remain.

## Original MSI demo utterances

The unpatched 2006 DLL chooses among several demo messages. The parser trace
stored under the original name describes one execution; it does not pin the
banner in the retained WAV. Explicitly synthesizing the discovered banner with
the existing `vt_eng-banner-free-2006.dll` removes the extra runtime banner
choice. The text fixtures here, plus the original Stage 17 fixtures, reproduce
these full WAVs exactly:

| Retained path below corpus-parity | Fixture here | Bytes | SHA-256 |
| --- | --- | ---: | --- |
| `runtime-original-msi-2006-prose.wav` | `demo-prose.txt` | 372876 | `03b266a9270fb206b900defd9a80c02e8ab17d9dfdf7bd85419fa5c9c9e39003` |
| `runtime-original-msi-2006-numbers.wav` | `demo-numbers.txt` | 373188 | `fe820fa817a53cbd1312e51b405aea64d4b42471ba774fcd62aa35c2f79f76ab` |
| `stage20/apple-plain-features/old-pah0-address.wav` | `demo-address.txt` | 320008 | `0bf0f679b9bb7bfb0c210e4fad9173efcaabbfe43eb7068ca742e24bf8c67e0c` |

`explicit-demo-batch.json` is the proposed runtime batch. It uses the existing
MSI voice and dictionaries and the existing direct API probe at its original
volume/rate/pitch settings; only the explicit text and existing analysis DLL
overlay differ. No vendor files were edited. The generated text includes the
banner's trailing space plus the separator space before the requested text.

The original `runtime-original-msi-2006-address.wav` has a clean-address header
(data length 109020), followed by demo-address PCM (319964 bytes), followed by
clean-address PCM (109020 bytes). It is reproduced exactly as:

```python
full_historical = clean_address[:44] + demo_address[44:] + clean_address[44:]
```

Generate `clean_address` with the same banner-free DLL and the unchanged
`stage17/address.txt`. The clean WAV is 109064 bytes, SHA-256
`7d1bf4f94baffb45ddbb573fc9670b6492030fb037dbc47acb0117e250893b5a`.
The resulting complete historical file is 429028 bytes, SHA-256
`e81b01d7398976d676214b4e572ec1e961e99e9bb91b2827283192efe6368db6`;
its full byte comparison passed. This establishes the byte recipe, without
claiming that the exact historical append/overwrite sequence has been traced.

Fresh generated demo outputs remain under
`/tmp/legacy-msi-explicit-banner/work/variants-free/{2,5,7}/`.
The clean address remains under that scratch work's `corpus-parity/` and the
verified complete splice is `address-splice.wav` beside that directory.
