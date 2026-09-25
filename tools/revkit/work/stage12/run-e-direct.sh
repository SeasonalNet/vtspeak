#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage12
scratch=/work/corpus-parity
xpid=

mkdir -p "$scratch"
cp "$work/input1.txt" "$scratch/stage12-e-original-input1.txt"
cp "$work/output.wav" "$scratch/stage12-e-original-output.wav"
restore() {
  cp "$scratch/stage12-e-original-input1.txt" "$work/input1.txt"
  cp "$scratch/stage12-e-original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/e-direct-xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

for input in "$@"; do
  cp "$probe/inputs/$input.txt" "$work/input1.txt"
  WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
    winedbg --gdb /samples/voicetext_paul.exe \
    < "$probe/trace-e-direct.gdb" > "$probe/$input-e-direct.log" 2>&1 || true
  if ! grep -q '^FORCED_E_RESULT' "$probe/$input-e-direct.log"; then
    printf 'forced E lookup was not reached for %s\n' "$input" >&2
    exit 1
  fi
done
