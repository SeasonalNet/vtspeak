#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
xpid=
restore() {
  cp "$probe/lipsync-output-original-input1.txt" "$work/input1.txt"
  cp "$probe/lipsync-output-original-output.wav" "$work/output.wav"
  cmp -s "$probe/lipsync-output-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/lipsync-output-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

if [ -e "$probe/lipsync-output-probe.txt" ]; then
  echo "refusing to overwrite existing lip-sync output" >&2
  exit 1
fi
cp "$work/input1.txt" "$probe/lipsync-output-original-input1.txt"
cp "$work/output.wav" "$probe/lipsync-output-original-output.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-lipsync-output.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-lipsync-output.gdb" > "$probe/lipsync-output.log" 2>&1
