# Stage 17: Lead 5 voice and package comparison

This stage compares the four supplied host/DLL pairs against their local model
packages with three shared text inputs (ordinary prose, a number/time phrase,
and an address/abbreviation phrase). The Paul and James packages carry the
2013 index and `tree3` layouts. Julie and Kate carry older index and `tree2`
layouts; the supplied hosts and DLLs request `tree3` files during load. Bridget
has no matching host or DLL in `binary/`, so it receives structural package
coverage only.

The runtime matrix uses the Stage 8 Wine image and keeps every model mount
read-only. It runs from the established Stage 5 working directory because the
hosts resolve fixed relative model paths there. The James and Julie hosts
request `data-jame` and `data-juli` respectively; the compose overrides mount
the repository's `data-james` and `data-julie` directories at those literal
paths. The runner backs up and restores Stage 5's `input1.txt` and `output.wav`
on exit.

Run from the repository root:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml \
  -f tools/revkit/work/stage17/compose.yaml run --rm runtime \
  /bin/bash /probe/run-voice-matrix.sh
```

`runtime-matrix.tsv` records process exit and fresh output size. Julie and Kate
use Wine file tracing to preserve the requested tree path and failed-open
status. WAV outputs are retained only for successful runs. The traces establish
the first load boundary reached by the prose sample; they are not a broad
compatibility or full-synthesis parity test.

## Additional James package DLL

A second James runtime probe used the `vt_eng_james16.dll` supplied beside the
locally extracted James package, mounted read-only as `/samples/vt_jam.dll`.
That DLL asks for `data-james/`, so the overlay mounts the same read-only voice
directory at that exact path as well as the `data-jame/` path used by the
checked-in host/DLL pair. It ran the same three fixtures through the supplied
James host and produced valid WAVE files. All three files have the same format,
byte length, and frame count as the Stage 17 run with `binary/vt_jam.dll`, but
their complete-file hashes and PCM payloads differ. This establishes a second
successful engine/package path for these inputs, not output parity.

The package DLL is not copied into the repository. Set `JAMES_PACKAGE_DLL` to
its local path before running this overlay:

```sh
JAMES_PACKAGE_DLL="/path/to/vt_eng_james16.dll" \
  docker compose -f tools/revkit/work/stage8/compose.yaml \
    -f tools/revkit/work/stage17/compose.yaml \
    -f tools/revkit/work/stage17/compose-package-james.yaml run --rm runtime \
    /bin/bash /probe/run-james-package.sh
```

The probe restores Stage 5's input and output fixtures on exit. Logs, WAVE
captures, and the hash manifest are retained here. The PCM/WAVE comparison is
in [`james-package-comparison.tsv`](james-package-comparison.tsv). The local package DLL's
SHA-256 is `7644aef98f08c5929283c15a719deb3a8395eb9b5b6fd80862f9aace16522580`.
Its source is the extracted James package under Downloads; no claim is made
about installer state or redistribution rights.

Other historical candidates were inspected but do not add accepted runtime
coverage. A Julie DLL from a directory named `crack` reaches an old common
dictionary request for the absent `dict-eng/dic.pos.txt2` and exits before its
voice tree loads; its integrity and licensing state are not established. A
Bridget DLL candidate also comes from a `crack` directory, and no Bridget host
is present. Neither candidate is treated as evidence of package compatibility.
The Julie file-open trace is retained as `julie-old-engine.log` only to show
that blocker.
