#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
prefix=${TRACE_PREFIX:?set TRACE_PREFIX to a variant label}
trace_log=${TRACE_LOG:-$probe/candidates-gdb.log}
input_backup=$(mktemp "$probe/.stage5-trace-input.XXXXXX")
output_backup=$(mktemp "$probe/.stage5-trace-output.XXXXXX")
xpid=

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  status=$?
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$input_backup" "$output_backup"
  exit "$status"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

cp "$probe/input.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-selection-$prefix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

trace_gdb=${TRACE_GDB:-/work/stage7/trace-candidates.gdb}
if [ "$trace_gdb" = "/work/stage7/trace-candidates.gdb" ]; then
  sed 's#/work/stage7/candidates-gdb.log#/probe/candidates-gdb.log#' \
    "$trace_gdb" > "$probe/trace-selection-$prefix.gdb"
else
  cp "$trace_gdb" "$probe/trace-selection-$prefix.gdb"
fi
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_kate.exe \
  < "$probe/trace-selection-$prefix.gdb" >"$probe/trace-selection-$prefix.log" 2>&1
rc=$?
set -e
if [ -f "$trace_log" ]; then
  cp "$trace_log" "$probe/selection-$prefix.log"
fi
bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
if [ "$bytes" -gt 44 ]; then
  cp "$work/output.wav" "$probe/selection-$prefix.wav"
fi
printf 'exit=%s\noutput_bytes=%s\noutput_sha256=%s\n' "$rc" "$bytes" "$hash" \
  > "$probe/selection-$prefix-result.txt"
