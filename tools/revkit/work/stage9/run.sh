#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage9
scratch=/work/corpus-parity
xpid=

mkdir -p "$scratch"
cp "$work/input1.txt" "$scratch/stage9-original-input1.txt"
cp "$work/output.wav" "$scratch/stage9-original-output.wav"
restore() {
  cp "$scratch/stage9-original-input1.txt" "$work/input1.txt"
  cp "$scratch/stage9-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$scratch/stage9-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for name in numbers abbreviations; do
  cp "$work/probes/stage6/$name/input1.txt" "$work/input1.txt"
  : > "$scratch/current-blocks.pcm"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe/trace-boundary.gdb" > "$probe/$name-run.log" 2>&1
  cp "$work/output.wav" "$probe/$name.wav"
  cp "$scratch/current-blocks.pcm" "$probe/$name-blocks.pcm"
done
