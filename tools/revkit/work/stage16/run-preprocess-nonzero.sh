#!/bin/bash
set -euo pipefail

probe=/work/stage16
work=/work/stage5
xpid=
restore() {
  cp "$probe/preprocess-nonzero-original-input1.txt" "$work/input1.txt"
  cp "$probe/preprocess-nonzero-original-output.wav" "$work/output.wav"
  cmp -s "$probe/preprocess-nonzero-original-input1.txt" "$work/input1.txt"
  cmp -s "$probe/preprocess-nonzero-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

if [ -e "$probe/preprocess-flag1-probe.out" ]; then
  echo "refusing to overwrite existing preprocess output" >&2
  exit 1
fi
cp "$work/input1.txt" "$probe/preprocess-nonzero-original-input1.txt"
cp "$work/output.wav" "$probe/preprocess-nonzero-original-output.wav"
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-preprocess-nonzero.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-preprocess-nonzero.gdb" > "$probe/preprocess-nonzero.log" 2>&1
