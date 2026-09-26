# Stage 18: patched DLL tree2 path probe

This experiment tests whether the supplied Kate DLL can get past its first
missing-tree failure when pointed at Kate's older `tree2` files. It creates an
ignored copy of `vt_kat.dll` and changes two guarded bytes in that copy:

- the tree-directory template `tree3/` to `tree2/` at file offset `0x7c8f0`;
- the last-character write in `FUN_10001900` from ASCII `3` to `2` at file
  offset `0x191f`, so the generated filenames end in `.tree2`.

The patch does not port the legacy tree reader. The supplied DLL calls
`FUN_100016e0` after loading each file; that routine expects the 7-byte
`tree3` header and variable-length node records. Julie/Kate `tree2` files use
a recursive format. This probe therefore measures whether path selection and
the next loader boundary can be reached; it is not a compatibility fix or a
supported way to synthesize with old voices.

## Observed result

The patched copy loaded as `/samples/vt_kat.dll`, then the host exited 1 with
no WAVE output. The patched `FUN_10001900` is shared by common and per-voice
trees, so it changed the first common-tree request from
`data-common/dict-eng/poly.tree3` to the absent `poly.tree2` before the loader
reached Kate's model tree. The trace records file-not-found status
`0xc0000034`. This confirms that changing the global extension byte is too
broad. It does not test how `FUN_100016e0` handles a Kate `tree2` payload.

The original `binary/vt_kat.dll` SHA-256 remains
`f93c60a7ab0baa3f4ec2271cdc46997f48982143f49088349c32ff1755388461`; the
patched copy SHA-256 is
`9f006d1474cc912535699bd94c04d5d7649a4489431059fa8b816ebcd10421b5`.
`patch-manifest.txt`, `runtime-result.txt`, and `tree2-patched-runtime.log`
record the reproduction evidence. No output WAVE was generated.

The original DLL and all model files are mounted read-only. Only the patched
copy replaces `/samples/vt_kat.dll` inside the disposable runtime container.
The runner backs up and restores the Stage 5 input and output fixtures.

Reproduce from the repository root:

```sh
PYTHONDONTWRITEBYTECODE=1 python3 tools/revkit/work/stage18/patch_tree2.py
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage18/compose.yaml run --rm runtime \
  /bin/bash /probe/run-probe.sh
```

The generated DLL copy is ignored. `patch-manifest.txt`, runtime result,
Wine trace, and any output WAVE are local probe evidence; the original
`binary/vt_kat.dll` remains byte-identical.
