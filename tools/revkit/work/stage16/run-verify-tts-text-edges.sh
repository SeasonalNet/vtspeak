#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage16
xpid=
restore() {
  cp /tmp/verify-tts-text-edges-input1.txt "$work/input1.txt"
  cp /tmp/verify-tts-text-edges-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" /tmp/verify-tts-text-edges-input1.txt
cp "$work/output.wav" /tmp/verify-tts-text-edges-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-verify-tts-text-edges.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-verify-tts-text-edges.gdb" > "$probe/verify-tts-text-edges-api.log" 2>&1
