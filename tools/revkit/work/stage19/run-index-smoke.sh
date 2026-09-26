#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=$probe/.stage5-smoke-input.original.$$
output_backup=$probe/.stage5-smoke-output.original.$$
xpid=

cp "$work/input1.txt" "$input_backup"
cp "$work/output.wav" "$output_backup"
restore() {
  cp "$input_backup" "$work/input1.txt"
  cp "$output_backup" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

cp "$probe/input.txt" "$work/input1.txt"
: > "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-index-smoke.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=+file timeout 180s \
  wine /samples/voicetext_kate.exe >"$probe/index-smoke-runtime.log" 2>&1
rc=$?
set -e
bytes=$(stat -c %s "$work/output.wav")
printf 'exit=%s\noutput_bytes=%s\n' "$rc" "$bytes" \
  > "$probe/index-smoke-result.txt"
