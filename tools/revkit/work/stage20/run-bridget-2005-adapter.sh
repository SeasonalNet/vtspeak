#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
prefix=${EVIDENCE_PREFIX:-bridget-hello}
input_backup=$(mktemp "$probe/.bridget-input.XXXXXX")
output_backup=$(mktemp "$probe/.bridget-output.XXXXXX")
xpid=
cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
  rm -f "$input_backup" "$output_backup"
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/$prefix-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

cp "$probe/bridget-hello.txt" "$work/input1.txt"
: > "$work/output.wav"
set +e
if [ "${TRACE_SELECTED:-0}" = 1 ]; then
  gdb_script=${GDB_SCRIPT:-/probe/trace-bridget-selected-units.gdb}
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_kate.exe \
    < "$gdb_script" \
    > "$probe/$prefix-runtime.log" 2>&1
else
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    wine /samples/voicetext_kate.exe >"$probe/$prefix-runtime.log" 2>&1
fi
rc=$?
set -e
bytes=$(stat -c %s "$work/output.wav")
hash=$(sha256sum "$work/output.wav" | cut -d ' ' -f 1)
printf 'exit=%s\nbytes=%s\nsha256=%s\n' "$rc" "$bytes" "$hash" \
  > "$probe/$prefix-result.txt"
if [ "$rc" -eq 0 ] && [ "$bytes" -gt 44 ]; then
  cp "$work/output.wav" "$probe/$prefix.wav"
else
  exit 1
fi
