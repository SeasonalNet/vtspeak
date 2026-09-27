#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
trace="$probe/trace-preprocess-input-errors.gdb"
log="$probe/preprocess-input-errors-api.log"
xvfb_log="$probe/xvfb-preprocess-input-errors.log"
xpid=
restore() {
  cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

for output in "$log" "$xvfb_log"; do
  if [ -e "$output" ]; then
    echo "refusing to overwrite $output" >&2
    exit 3
  fi
done

cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$xvfb_log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$trace" \
  > "$log" 2>&1
