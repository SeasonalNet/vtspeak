#!/bin/bash
set -euo pipefail

work=/work/stage5
scratch=/work/corpus-parity
xpid=

mkdir -p "$scratch"
cp "$work/output.wav" "$scratch/original-output.wav"
restore() {
  cp "$scratch/original-output.wav" "$work/output.wav"
  if [ -n "$xpid" ]; then
    kill "$xpid" 2>/dev/null || true
  fi
}
trap restore EXIT

Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$scratch/xvfb.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1

WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < /work/scripts/probe_init_path.gdb > "$scratch/init-path.log" 2>&1
