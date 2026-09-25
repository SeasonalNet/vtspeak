#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage12
scratch=/work/corpus-parity
xpid=

mkdir -p "$scratch"
cp "$work/input1.txt" "$scratch/stage12-original-input1.txt"
cp "$work/output.wav" "$scratch/stage12-original-output.wav"
restore() {
  cp "$scratch/stage12-original-input1.txt" "$work/input1.txt"
  cp "$scratch/stage12-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for input in "$@"; do
  cp "$probe/inputs/$input.txt" "$work/input1.txt"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe/trace.gdb" > "$probe/$input-run.log" 2>&1
  cp "$work/output.wav" "$probe/$input.wav"
done
