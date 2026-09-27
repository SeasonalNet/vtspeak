# 2006 Kate comparison-only license-status control

This control uses an ignored copy of the extracted MSI DLL. It changes no
vendor binary or model file. The original DLL hash is checked before each copy
is made; the generated DLL remains under the ignored `work/corpus-parity/`
directory.

The old DLL does not export the 2013 `VT_CheckLicense_ENG` entry used by the
peer patch. Its internal `FUN_10022430` checker returns the status exposed by
`VT_GetTTSInfo_ENG` selector 1 and consumed by `VT_LOADTTS_ENG`. The one-byte
status change lets this old DLL return success. With the exact MSI
verification file mounted, both the original and patched DLL synthesize all
three fixtures successfully and produce byte-identical WAVs. The patch does
not establish behavior for another verification record or general use.
Captured API results and parser traces for both controls are in this folder;
the matched runs are byte-identical for prose, numbers, and address.

The two guarded changes in [`patch_copy.py`](patch_copy.py) are:

- At VA `0x10022542`, change `jge` to an unconditional jump to the checker's
  zero-return epilogue.

## Reproduction

Extract the 2006 Kate MSI locally so these files exist under
`tools/revkit/work/corpus-parity/kate-msi/Program Files/NeoSpeech/Kate16/`:
`lib/vt_eng.dll`, `data-kate/`, `data-common/dict-eng/`, and
`data-common/verify/verification.txt`. Then run from the repository root:

```sh
python3 tools/revkit/work/stage19/license-control/patch_copy.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/license-control/compose.yaml \
  run --rm runtime /bin/bash /work/stage19/license-control/run-api.sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage19/license-control/compose.yaml \
  run --rm runtime /bin/bash /work/stage19/license-control/run-parser-trace.sh
```

The runners preserve the existing `runtime-original-msi-2006-*.wav` files.
The patched outputs and traces are written under ignored
`tools/revkit/work/corpus-parity/`.

For the matched unpatched control, use `compose-control.yaml` and the
`run-control-api.sh` and `run-control-parser-trace.sh` runners in place of
`compose.yaml` and the patched runners. The exact MSI verification file is
315 bytes. A prior run that mounted the newer 468-byte local verification
file followed the nag-text path; the correct MSI record did not. In the
matched control, the original DLL returns load status `0`, synthesis status
`1` for all three fixtures, and parses only the requested sentences. The
patched and unpatched WAVE SHA-256 values are identical:

| Fixture | SHA-256 | Frames |
| --- | --- | ---: |
| Prose | `77a4f7e8e39e12fed2723709575609d5b8dec90d27c7e1c6977368cda0a3ebff` | 51,512 |
| Numbers | `5b24220b078e12b0bc018a259e7358e799a86bc260f74abbc95fb8147deccd25` | 46,566 |
| Address | `7d1bf4f94baffb45ddbb573fc9670b6492030fb037dbc47acb0117e250893b5a` | 54,510 |
