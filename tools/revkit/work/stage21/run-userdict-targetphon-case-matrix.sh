#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/targetphon-case-matrix-input1.txt "$work/input1.txt"
  cp /tmp/targetphon-case-matrix-output.wav "$work/output.wav"
  cmp -s /tmp/targetphon-case-matrix-input1.txt "$work/input1.txt"
  cmp -s /tmp/targetphon-case-matrix-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}
cp "$work/input1.txt" /tmp/targetphon-case-matrix-input1.txt
cp "$work/output.wav" /tmp/targetphon-case-matrix-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-targetphon-case-matrix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-userdict-targetphon-case-matrix.gdb" > "$probe/targetphon-case-matrix-api.log" 2>&1
grep '^TARGET_PHON_CASE ' "$probe/targetphon-case-matrix-api.log" | tail -n 8

