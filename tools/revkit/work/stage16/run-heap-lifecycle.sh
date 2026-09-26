#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
xpid=
restore() {
  cp "$probe/heap-lifecycle-original-input1.txt" "$work/input1.txt"
  cp "$probe/heap-lifecycle-original-output.wav" "$work/output.wav"
  cmp -s "$probe/heap-lifecycle-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/heap-lifecycle-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$probe/heap-lifecycle-original-input1.txt"
cp "$work/output.wav" "$probe/heap-lifecycle-original-output.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-heap-lifecycle.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-heap-lifecycle.gdb" > "$probe/heap-lifecycle-api.log" 2>&1
