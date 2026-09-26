#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=$probe/.stage5-debug-input.original.$$
output_backup=$probe/.stage5-debug-output.original.$$
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-index-debug.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 30s \
  winedbg --gdb /samples/voicetext_kate.exe \
  < "$probe/trace-index-progress.gdb" >"$probe/index-progress.log" 2>&1
rc=$?
set -e
printf 'exit=%s\n' "$rc" > "$probe/index-progress-result.txt"
