#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
fixture=${FIXTURE:?set FIXTURE to prose, numbers, or address}
prefix=${TRACE_PREFIX:?set TRACE_PREFIX to a result label}
trace=${TRACE_GDB:-/probe/trace-cnas-path.gdb}
case "$fixture" in
  prose) source_fixture=input ;;
  numbers|address) source_fixture=$fixture ;;
  *) printf 'unsupported fixture: %s\n' "$fixture" >&2; exit 2 ;;
esac
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

cp "/work/stage17/$source_fixture.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-$prefix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 300s \
  winedbg --gdb /samples/voicetext_kate.exe < "$trace" \
  >"$probe/$prefix.log" 2>&1
rc=$?
set -e
bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
if [ -f "$probe/cnas-path-gdb.log" ]; then
  cp "$probe/cnas-path-gdb.log" "$probe/$prefix-path-gdb.log"
fi
if [ "$bytes" -gt 44 ]; then cp "$work/output.wav" "$probe/$prefix.wav"; fi
printf 'fixture=%s\nexit=%s\noutput_bytes=%s\noutput_sha256=%s\n' \
  "$fixture" "$rc" "$bytes" "$hash" > "$probe/$prefix-result.txt"
if [ "$rc" -ne 0 ] || [ "$bytes" -le 44 ]; then exit 1; fi
