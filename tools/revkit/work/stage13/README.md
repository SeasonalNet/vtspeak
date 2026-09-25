# Stage 13: host-side Intel PT capture

This probe attaches Linux `perf` from the host to the host PID of the existing
Stage 8 Wine runtime container. The target container keeps its configured
non-root UID/GID, `cap_drop: ALL`, read-only root filesystem, no-network mode,
and `no-new-privileges` setting. The host collector needs `sudo` because this
machine restricts unprivileged perf events.

The container pauses at a gate before starting one existing Stage 12
`run-effects.sh` fixture. The host script attaches `perf record` to the
container's init PID; perf's default child-event inheritance follows the
runner and Wine descendants. It then opens the gate. A temporary copy of the
Stage 5 directory is mounted at the same `/work/stage5` path used by the
successful Stage 12 runs. This preserves the expected Wine current directory
while keeping the original Stage 5 files untouched. The Stage 12 runner stores
the fresh GDB log and WAV beside the requested trace.

Run one bounded probe from the repository root:

```sh
tools/revkit/work/stage13/capture-wine-intel-pt.sh \
  /tmp/vtspeak-g83-control.perf.data \
  acoustic-control numeric-context-g83-can-anybody-help
```

The output path must not already exist. `perf.data` may become large; keep the
first capture to one fixture. The run emits its `perf` startup log and points
to the sidecar GDB log and WAV. To decode only the VoiceText DLL events, make a
temporary symbol root that maps the container's `/samples/` paths to the local,
read-only vendor inputs. Find the host PID of `voicetext_paul.exe` in the perf
mmap records first; it changes on each run. `-f` lets `perf script` read a file
owned by root.

```sh
pt_symfs=$(mktemp -d /tmp/vtspeak-pt-symfs.XXXXXX)
mkdir -p "$pt_symfs/samples"
ln -s "$PWD/binary/voicetext_paul.exe" \
  "$pt_symfs/samples/voicetext_paul.exe"
ln -s "$PWD/binary/vt_pau.dll" "$pt_symfs/samples/vt_pau.dll"
sudo perf script -f -i /tmp/vtspeak-g83-control.perf.data \
  --symfs="$pt_symfs" --pid=112937 --itrace=b \
  --dsos=/samples/voicetext_paul.exe,/samples/vt_pau.dll \
  > /tmp/vtspeak-voice-branches.txt \
  2> /tmp/vtspeak-voice-decode.log
```

Replace `112937` with the VoiceText process PID found in this run's mmap/comm
records. The `--symfs` links let perf read the same PE images loaded in the
container; without them, this trace's VoiceText addresses decode as instruction
mismatches. `--itrace=b` synthesizes branch records only; adding `e` also emits
large volumes of decoder error events. Keep the raw `perf.data` local because
it can be large. GDB remains the source for register and memory values at
selected addresses. The trace alone does not label TPP fields or prove
unexecuted behavior.

## Capture prerequisites and limits

- The active Docker daemon must permit the existing `docker compose run`
  workflow. The script does not change the Compose security settings.
- Host `perf` must support `intel_pt`, and the user must be able to run
  `sudo perf record`.
- The script waits for `perf -vv` to report a `perf_event_open` attempt before
  releasing the container gate. If attachment fails, it stops before running
  VoiceText.
- The target is one existing Stage 12 fixture and inherits the runner's
  existing per-invocation 180-second limit. Intel PT output and decode time
  depend on the execution path.
- Some decoded destinations may appear as `[unknown]`; inspect the available
  image mappings and correlate addresses with the Ghidra image base before
  interpreting those edges.

The first operator capture and its bounded address correlation are summarized
in [the Stage 13 trace note](../../../docs/reverse-engineering/intel-pt-wine-capture-2026-09-25.md).
