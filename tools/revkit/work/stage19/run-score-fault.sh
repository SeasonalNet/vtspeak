#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/probe
input_backup=$(mktemp "$probe/.stage5-score-input.XXXXXX")
output_backup=$(mktemp "$probe/.stage5-score-output.XXXXXX")
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

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp >"$probe/xvfb-score-fault.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

set +e
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 60s \
  winedbg --gdb /samples/voicetext_kate.exe \
  < "$probe/trace-score-fault.gdb" >"$probe/score-fault.log" 2>&1
rc=$?
set -e
printf 'exit=%s\n' "$rc" > "$probe/score-fault-result.txt"
