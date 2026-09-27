#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
xpid=
restore() {
  cp "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-flags-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-flags-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-preprocess-flag2-isolated.log" 2>&1 &
xpid=$!
sleep 1
DISPLAY=:99 WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-preprocess-flag2-isolated.gdb" \
  > "$probe/preprocess-flag2-isolated-api.log" 2>&1
find "$probe" -maxdepth 1 -type f -mmin -2 -name '*preprocess*' \
  -printf '%f %s bytes\n' > "$probe/preprocess-flag2-isolated-files.txt"
