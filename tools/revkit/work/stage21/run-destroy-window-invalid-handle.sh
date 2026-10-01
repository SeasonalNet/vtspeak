#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
xpid=
restore() {
  cp /tmp/destroy-invalid-input1.txt "$work/input1.txt"
  cp /tmp/destroy-invalid-output.wav "$work/output.wav"
  cmp -s /tmp/destroy-invalid-input1.txt "$work/input1.txt"
  cmp -s /tmp/destroy-invalid-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
test -f "$work/input1.txt"
test -f "$work/output.wav"
cp "$work/input1.txt" /tmp/destroy-invalid-input1.txt
cp "$work/output.wav" /tmp/destroy-invalid-output.wav
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-destroy-window-invalid-handle.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-destroy-window-invalid-handle.gdb" > "$probe/destroy-window-invalid-handle-api.log" 2>&1
grep '^DESTROY_WINDOW_' "$probe/destroy-window-invalid-handle-api.log"
