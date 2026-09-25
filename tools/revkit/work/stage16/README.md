# Stage 16: Lead 6 API compatibility probe

This stage captures the local Paul's file, buffer, information,
configuration, and selected playback API behavior. It uses the existing Stage
8 Wine runtime container and calls APIs in the loaded process at the
`VT_TextToFile_ENG` entry breakpoint (`0x1001da50`). Vendor binary and model
mounts remain read-only.

The format runner accepts one selector from 0 through 10:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-file-format.sh 4
```

The error runner accepts `null-text`, `empty-text`, or `null-path`:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-file-error.sh null-text
```

The buffer runner accepts a format, flag, and thread ID:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-buffer.sh 0 0 0
```

It calls `VT_TextToBuffer_ENG` in the already loaded process with a 1 MiB
buffer. It writes the captured buffer only when the call returns success, so
failed thread-creation cases do not expose uninitialized process memory.

`run-buffer-errors.sh` captures the buffer API's invalid-format, null-text,
empty-text, and null-buffer returns. `run-info.sh` queries all declared
`VT_GetTTSInfo_ENG` requests plus invalid-request, null-value, and string-size
errors in one loaded process:

```sh
docker compose -f tools/revkit/work/stage8/compose.yaml run --rm runtime \
  /bin/bash /work/stage16/run-info.sh
```

`run-config.sh` queries the current speaker settings, applies upper-bound
values through the configuration setters, checks the resulting getters, then
captures file output with the synthesis arguments left at `-1`.
`run-play-errors.sh` records null/empty/valid-text playback returns and calls
stop, pause, and restart controls. Valid playback fails initialization in this
Wine runtime, so this does not exercise active audio controls.
`run-buffer-capacity.sh` passes several `output_len` values with a physically
oversized output allocation and boundary sentinels; it captures only the
API-reported output bytes.
`run-texttypes.sh` sends the same ASCII text with the default and each declared
text-format value, capturing one selector-4 WAVE for each call.

Runners that call through the Stage 5 executable back up and restore its input
and output fixtures on exit. Each selector result is saved as
`file-format-N.out`; GDB logs record
the API entry and signed-short return. The error runner mutates only the
in-process API arguments or first text byte. The captured matrix and its
limits are documented in
[`Lead 6 API behavior`](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md).

`run-file-repeat.sh` issues two identical selector 4 calls through the DLL in
one loaded process, saves each result, and logs both return values. This
single repeat control is not a general persistence guarantee.

The buffer runner covers formats 0–3 at flag 0/thread ID 0, plus three format
0 flag/thread combinations that return a create-thread error. This does not
cover successful threaded processing or polling. The capacity probe shows
that a smaller incoming `output_len` does not limit writes for the tested
synchronous format-0 path. Configuration and selected playback calls are
documented in the [Lead 6 report](../../../../docs/reverse-engineering/lead6-file-api-behavior-2026-09-25.md).
Successful playback, successful threaded processing, alternate text formats,
and other voice packages remain untested.
