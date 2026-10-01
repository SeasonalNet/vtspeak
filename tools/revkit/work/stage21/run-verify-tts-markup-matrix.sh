#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
restore() {
  cp /tmp/verify-tts-markup-matrix-input1.txt "$work/input1.txt"
  cp /tmp/verify-tts-markup-matrix-output.wav "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" /tmp/verify-tts-markup-matrix-input1.txt
cp "$work/output.wav" /tmp/verify-tts-markup-matrix-output.wav
trap restore EXIT
cp "$probe/input.txt" "$work/input1.txt"
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-verify-tts-markup-matrix.log" 2>&1 &
xpid=$!
export DISPLAY=:99
sleep 1
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-verify-tts-markup-matrix.gdb" > "$probe/verify-tts-markup-matrix-api.log" 2>&1
grep '^VERIFY_TTS_MARKUP ' "$probe/verify-tts-markup-matrix-api.log"
