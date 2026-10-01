#!/bin/bash
set -euo pipefail

probe=/work/stage21
work=/work/stage21/sandbox/stage5
xpid=
restore() {
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
trap restore EXIT

test -f "$work/input1.txt"
test -f "$work/output.wav"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-unit-history-byte-matrix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-unit-history-byte-matrix.gdb" > "$probe/unit-history-byte-matrix-api.log" 2>&1
