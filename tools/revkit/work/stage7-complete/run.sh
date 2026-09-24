#!/bin/bash
set -euo pipefail
work=/work/stage5
probe=/work/stage7-complete
xpid=
restore() {
  cp "$probe/original-input1.txt" "$work/input1.txt"
  cp "$probe/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" "$probe/original-input1.txt"
cp "$work/output.wav" "$probe/original-output.wav"
trap restore EXIT
cp "$probe/input1.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-stage7.gdb" > "$probe/run.log" 2>&1
cp "$work/output.wav" "$probe/output.wav"
